// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package core

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"os"
	"path"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

//go:generate gencodec -type Genesis -field-override genesisSpecMarshaling -out gen_genesis.go
//go:generate gencodec -type GenesisAccount -field-override genesisAccountMarshaling -out gen_genesis_account.go

var errGenesisNoConfig = errors.New("genesis has no chain configuration")

// Genesis specifies the header fields, state of a genesis block. It also defines hard
// fork switch-over blocks through the chain configuration.
type Genesis struct {
	Config        *params.ChainConfig `json:"config"`
	EconomicModel *xcom.EconomicModel `json:"economicModel"`
	Nonce         []byte              `json:"nonce"`
	Timestamp     uint64              `json:"timestamp"`
	ExtraData     []byte              `json:"extraData"`
	GasLimit      uint64              `json:"gasLimit"   gencodec:"required"`
	Coinbase      common.Address      `json:"coinbase"`
	Alloc         GenesisAlloc        `json:"alloc"      gencodec:"required"`

	// These fields are used for consensus tests. Please don't use them
	// in actual genesis blocks.
	Number     uint64      `json:"number"`
	GasUsed    uint64      `json:"gasUsed"`
	ParentHash common.Hash `json:"parentHash"`
}

// GenesisAlloc specifies the initial state that is part of the genesis block.
type GenesisAlloc map[common.Address]GenesisAccount

func (ga *GenesisAlloc) UnmarshalJSON(data []byte) error {
	m := make(map[common.UnprefixedAddress]GenesisAccount)
	if err := json.Unmarshal(data, &m); err != nil {
		return err
	}
	*ga = make(GenesisAlloc)
	for addr, a := range m {
		(*ga)[common.Address(addr)] = a
	}
	return nil
}

// GenesisAccount is an account in the state of the genesis block.
type GenesisAccount struct {
	Code       []byte                      `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *big.Int                    `json:"balance" gencodec:"required"`
	Nonce      uint64                      `json:"nonce,omitempty"`
	PrivateKey []byte                      `json:"secretKey,omitempty"` // for tests
}

// field type overrides for gencodec
type genesisSpecMarshaling struct {
	Nonce     hexutil.Bytes
	Timestamp math.HexOrDecimal64
	ExtraData hexutil.Bytes
	GasLimit  math.HexOrDecimal64
	GasUsed   math.HexOrDecimal64
	Number    math.HexOrDecimal64
	Alloc     map[common.UnprefixedAddress]GenesisAccount
}

type genesisAccountMarshaling struct {
	Code       hexutil.Bytes
	Balance    *math.HexOrDecimal256
	Nonce      math.HexOrDecimal64
	Storage    map[storageJSON]storageJSON
	PrivateKey hexutil.Bytes
}

// storageJSON represents a 256 bit byte array, but allows less than 256 bits when
// unmarshaling from hex.
type storageJSON common.Hash

func (h *storageJSON) UnmarshalText(text []byte) error {
	text = bytes.TrimPrefix(text, []byte("0x"))
	if len(text) > 64 {
		return fmt.Errorf("too many hex characters in storage key/value %q", text)
	}
	offset := len(h) - len(text)/2 // pad on the left
	if _, err := hex.Decode(h[offset:], text); err != nil {
		fmt.Println(err)
		return fmt.Errorf("invalid hex storage key/value %q", text)
	}
	return nil
}

func (h storageJSON) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// GenesisMismatchError is raised when trying to overwrite an existing
// genesis block with an incompatible one.
type GenesisMismatchError struct {
	Stored, New common.Hash
}

func (e *GenesisMismatchError) Error() string {
	return fmt.Sprintf("database already contains an incompatible genesis block (have %x, new %x)", e.Stored[:8], e.New[:8])
}

// SetupGenesisBlock writes or updates the genesis block in db.
// The block that will be used is:
//
//                          genesis == nil       genesis != nil
//                       +------------------------------------------
//     db has no genesis |  main-net default  |  genesis
//     db has genesis    |  from DB           |  genesis (if compatible)
//
// The stored chain configuration will be updated if it is compatible (i.e. does not
// specify a fork block below the local head block). In case of a conflict, the
// error is a *params.ConfigCompatError and the new, unwritten config is returned.
//
// The returned chain configuration is never nil.
func SetupGenesisBlock(db ethdb.Database, snapshotPath string, genesis *Genesis) (*params.ChainConfig, common.Hash, error) {

	if genesis != nil && genesis.Config == nil {
		log.Error("Failed to SetupGenesisBlock, the config of genesis is nil")
		return params.AllEthashProtocolChanges, common.Hash{}, errGenesisNoConfig
	}

	// Just commit the new block if there is no stored genesis block.
	stored := rawdb.ReadCanonicalHash(db, 0)

	if (stored == common.Hash{}) {
		if genesis == nil {
			log.Info("Writing default main-net genesis block")
			genesis = DefaultGenesisBlock()
		} else {
			log.Info("Writing custom genesis block")
		}

		// check genesis version
		if genesis.Config == nil || genesis.Config.Version <= 0 {
			log.Error("genesis version is missed")
			return nil, common.Hash{}, errors.New("genesis version is missed")
		}

		// check EconomicModel configuration
		if err := xcom.CheckEconomicModel(); nil != err {
			log.Error("Failed to check economic config", "err", err)
			return nil, common.Hash{}, err
		}
		var sdb snapshotdb.DB
		if snapshotPath != "" {
			os.RemoveAll(snapshotPath)
			sdb = snapshotdb.Instance()
			defer sdb.Close()
		}
		block, err := genesis.Commit(db, sdb)
		log.Debug("SetupGenesisBlock Hash", "Hash", block.Hash().Hex())
		return genesis.Config, block.Hash(), err
	}

	// Check whether the genesis block is already written.
	if genesis != nil {
		hash := genesis.ToBlock(nil, nil).Hash()
		if hash != stored {
			log.Error("Failed to compare the genesisHash and stored", "genesisHash", hash, "stored", stored)
			return genesis.Config, hash, &GenesisMismatchError{stored, hash}
		}
	}

	// Get the existing EconomicModel configuration.
	ecCfg := rawdb.ReadEconomicModel(db, stored)
	if nil == ecCfg {
		log.Warn("Found genesis block without EconomicModel config")
		rawdb.WriteEconomicModel(db, stored, xcom.GetEc(xcom.DefaultMainNet))
	}

	xcom.ResetEconomicDefaultConfig(ecCfg)

	// Get the existing chain configuration.
	newcfg := genesis.configOrDefault(stored) // TODO this line Maybe delete
	storedcfg := rawdb.ReadChainConfig(db, stored)
	if storedcfg == nil {
		log.Warn("Found genesis block without chain config")
		rawdb.WriteChainConfig(db, stored, newcfg)
		return newcfg, stored, nil
	}

	// Sp ecial case: don't change the existing config of a non-mainnet chain if no new
	// config is supplied. These chains would get AllProtocolChanges (and a compat error)
	// if we just continued here.
	if genesis == nil && stored != params.MainnetGenesisHash {
		return storedcfg, stored, nil
	}

	// Check config compatibility and write the config. Compatibility errors
	// are returned to the caller unless we're already at block zero.
	height := rawdb.ReadHeaderNumber(db, rawdb.ReadHeadHeaderHash(db))
	if height == nil {
		log.Error("Failed to query header number by header hash", "headerHash", rawdb.ReadHeadHeaderHash(db))
		return newcfg, stored, fmt.Errorf("missing block number for head header hash")
	}
	compatErr := storedcfg.CheckCompatible(newcfg, *height)
	if compatErr != nil && *height != 0 && compatErr.RewindTo != 0 {
		log.Error("Failed to CheckCompatible", "height", *height, "err", compatErr)
		return newcfg, stored, compatErr
	}
	rawdb.WriteChainConfig(db, stored, newcfg)
	return newcfg, stored, nil
}

func (g *Genesis) configOrDefault(ghash common.Hash) *params.ChainConfig {
	switch {
	case g != nil:
		return g.Config
	case ghash == params.MainnetGenesisHash:
		return params.MainnetChainConfig
	case ghash == params.TestnetGenesisHash:
		return params.TestnetChainConfig
	default:
		return params.AllEthashProtocolChanges
	}
}

// ToBlock creates the genesis block and writes state of a genesis specification
// to the given database (or discards it if nil).
func (g *Genesis) ToBlock(db ethdb.Database, sdb snapshotdb.DB) *types.Block {
	if db == nil {
		db = ethdb.NewMemDatabase()
	}
	var snapDB snapshotdb.DB
	if sdb == nil {
		var err error
		log.Info("begin open snapshotDB in tmp")
		snapDB, err = snapshotdb.Open(path.Join(os.TempDir(), snapshotdb.DBPath), 0, 0)
		if err != nil {
			panic(err)
		}
		defer snapDB.Clear()
	} else {
		snapDB = sdb
	}

	genesisIssuance := new(big.Int)

	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	// First, Store the PlatONFoundation and CommunityDeveloperFoundation
	statedb.AddBalance(xcom.PlatONFundAccount(), xcom.PlatONFundBalance())
	statedb.AddBalance(xcom.CDFAccount(), xcom.CDFBalance())

	genesisIssuance = genesisIssuance.Add(genesisIssuance, xcom.PlatONFundBalance())
	genesisIssuance = genesisIssuance.Add(genesisIssuance, xcom.CDFBalance())

	for addr, account := range g.Alloc {
		statedb.AddBalance(addr, account.Balance)
		statedb.SetCode(addr, account.Code)
		statedb.SetNonce(addr, account.Nonce)
		for key, value := range account.Storage {

			statedb.SetState(addr, key.Bytes(), value.Bytes())
		}

		genesisIssuance = genesisIssuance.Add(genesisIssuance, account.Balance)
	}
	log.Debug("genesisIssuance", "amount", genesisIssuance)

	// Initialized Govern Parameters
	if err := gov.InitGenesisGovernParam(snapDB); err != nil {
		log.Error("Failed to init govern parameter in snapshotdb", "err", err)
		panic("Failed to init govern parameter in snapshotdb")
	}

	// Store genesis version into governance data
	if err := genesisPluginState(g, statedb, genesisIssuance, g.Config.Version); nil != err {
		panic("Failed to Store xxPlugin genesis statedb: " + err.Error())
	}

	// Store genesis staking data
	if err := genesisStakingData(snapDB, g, statedb, g.Config.Version); nil != err {
		panic("Failed Store staking: " + err.Error())
	}

	root := statedb.IntermediateRoot(false)

	log.Debug("ToBlock IntermediateRoot", "root", root.Hex())

	head := &types.Header{
		Number:     new(big.Int).SetUint64(g.Number),
		Nonce:      types.EncodeNonce(g.Nonce),
		Time:       new(big.Int).SetUint64(g.Timestamp),
		ParentHash: g.ParentHash,
		Extra:      g.ExtraData,
		GasLimit:   g.GasLimit,
		GasUsed:    g.GasUsed,
		Coinbase:   vm.RewardManagerPoolAddr,
		Root:       root,
	}
	if g.GasLimit == 0 {
		head.GasLimit = params.GenesisGasLimit
	}

	if _, err := statedb.Commit(false); nil != err {
		panic("Failed to commit genesis stateDB: " + err.Error())
	}
	if err := statedb.Database().TrieDB().Commit(root, true, true); nil != err {
		panic("Failed to trieDB commit by genesis: " + err.Error())
	}

	block := types.NewBlock(head, nil, nil)

	if err := snapDB.SetCurrent(block.Hash(), *common.Big0, *common.Big0); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. genesisHash: %s, error:%s", block.Hash().Hex(), err.Error()))
	}

	log.Debug("Call ToBlock finished", "genesisHash", block.Hash().Hex())
	return block
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(db ethdb.Database, sdb snapshotdb.DB) (*types.Block, error) {
	block := g.ToBlock(db, sdb)
	if block.Number().Sign() != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}

	log.Debug("Commit Hash", "hash", block.Hash().Hex(), "number", block.NumberU64())

	rawdb.WriteBlock(db, block)
	rawdb.WriteReceipts(db, block.Hash(), block.NumberU64(), nil)
	rawdb.WriteCanonicalHash(db, block.Hash(), block.NumberU64())
	rawdb.WriteHeadBlockHash(db, block.Hash())
	rawdb.WriteHeadHeaderHash(db, block.Hash())

	config := g.Config
	if config == nil {
		config = params.AllEthashProtocolChanges
	}
	rawdb.WriteChainConfig(db, block.Hash(), config)
	rawdb.WriteEconomicModel(db, block.Hash(), g.EconomicModel)
	return block, nil
}

// MustCommit writes the genesis block and state to db, panicking on error.
// The block is committed as the canonical head block.
func (g *Genesis) MustCommit(db ethdb.Database) *types.Block {
	block, err := g.Commit(db, snapshotdb.Instance())
	if err != nil {
		panic(err)
	}
	return block
}

// GenesisBlockForTesting creates and writes a block in which addr has the given wei balance.
func GenesisBlockForTesting(db ethdb.Database, addr common.Address, balance *big.Int) *types.Block {
	g := Genesis{Alloc: GenesisAlloc{addr: {Balance: balance}}}
	return g.MustCommit(db)
}

// DefaultGenesisBlock returns the PlatON main net genesis block.
func DefaultGenesisBlock() *Genesis {

	// TODO this should change
	generalAddr := common.HexToAddress("0x5437959B69eD1014cf6Aa8B4a2c77e7Ba2341955")
	generalBalance, _ := new(big.Int).SetString("9718188019000000000000000000", 10)

	rewardMgrPoolIssue, _ := new(big.Int).SetString("200000000000000000000000000", 10)

	genesis := Genesis{
		Config:    params.MainnetChainConfig,
		Nonce:     hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23"),
		Timestamp: 0,
		ExtraData: hexutil.MustDecode("0xd782070186706c61746f6e86676f312e3131856c696e757800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:  params.GenesisGasLimit,
		Alloc: map[common.Address]GenesisAccount{
			vm.RewardManagerPoolAddr: {Balance: rewardMgrPoolIssue},
			generalAddr:              {Balance: generalBalance},
		},
		EconomicModel: xcom.GetEc(xcom.DefaultMainNet),
	}
	xcom.SetNodeBlockTimeWindow(genesis.Config.Cbft.Period / 1000)
	xcom.SetPerRoundBlocks(uint64(genesis.Config.Cbft.Amount))
	return &genesis
}

// DefaultTestnetGenesisBlock returns the PlatON test net genesis block.
func DefaultTestnetGenesisBlock() *Genesis {

	// TODO this should change
	generalAddr := common.HexToAddress("0x9bbac0df99f269af1473fd384cb0970b95311001")
	generalBalance, _ := new(big.Int).SetString("9718188019000000000000000000", 10)

	rewardMgrPoolIssue, _ := new(big.Int).SetString("200000000000000000000000000", 10)

	genesis := Genesis{
		Config:    params.TestnetChainConfig,
		Nonce:     hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23"),
		ExtraData: hexutil.MustDecode("0xd782070186706c61746f6e86676f312e3131856c696e757800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		GasLimit:  params.GenesisGasLimit,
		Timestamp: 1546300800000,
		Alloc: map[common.Address]GenesisAccount{
			vm.RewardManagerPoolAddr: {Balance: rewardMgrPoolIssue},
			generalAddr:              {Balance: generalBalance},
		},
		EconomicModel: xcom.GetEc(xcom.DefaultTestNet),
	}
	xcom.SetNodeBlockTimeWindow(genesis.Config.Cbft.Period / 1000)
	xcom.SetPerRoundBlocks(uint64(genesis.Config.Cbft.Amount))
	return &genesis
}

func DefaultGrapeGenesisBlock() *Genesis {
	return &Genesis{
		Config:    params.GrapeChainConfig,
		Timestamp: 1492009146,
		ExtraData: hexutil.MustDecode("0xd782070186706c61746f6e86676f312e3131856c696e757800000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000"),
		//GasLimit:  3150000000,
		GasLimit: params.GenesisGasLimit,
		Alloc:    decodePrealloc(testnetAllocData),
	}
}

func decodePrealloc(data string) GenesisAlloc {
	var p []struct{ Addr, Balance *big.Int }
	if err := rlp.NewStream(strings.NewReader(data), 0).Decode(&p); err != nil {
		panic(err)
	}
	ga := make(GenesisAlloc, len(p))
	for _, account := range p {
		ga[common.BigToAddress(account.Addr)] = GenesisAccount{Balance: account.Balance}
	}
	return ga
}

func Alloc() GenesisAlloc {
	return decodePrealloc(testnetAllocData)
}
