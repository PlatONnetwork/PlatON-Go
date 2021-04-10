// Copyright 2016 The go-ethereum Authors
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

package params

import (
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// Genesis hashes to enforce below configs on.
var (
	MainnetGenesisHash = common.HexToHash("0xcb6867d07b029fe59e9960f09470a8476e03c73fa861511231855294aed74e55")
	TestnetGenesisHash = common.HexToHash("0x47c7f3dea852b28c45d5fee6c950321486780065745179c7a296934e493f18e1")
)

var TrustedCheckpoints = map[common.Hash]*TrustedCheckpoint{
	MainnetGenesisHash: MainnetTrustedCheckpoint,
	TestnetGenesisHash: TestnetTrustedCheckpoint,
}

var (
	initialMainNetConsensusNodes = []initNode{
		{
			"enode://4b23f32791d7caf6c7e601d6c1d8369f64fe8aa020c7d897e37ae28539b1761f45fce40edb1c429c174b124656de614a117ce6ad44b2c1286ff3da04320c7113@mf1.a4f9.platon.network:16789",
			"bc3aa98cf8ade88d32c5efec9802126c41fb37849741143cb21ecbb8fe31964a451904a65bb3740faa24474a6dd0810437ac324e64c11f5877d146c0464ab4ed0d512f328a263e138f967c3a56fdd849b326a6c9c0c6b60f5ac05776272d1c95",
		},
		{
			"enode://10a29ee607a9f4e7d87669add40ef5ce37fb553d9c39cb6ed413da22c1e646df015198abde3459146d99ee84b3c64396b7ab04b63dcac9b42803ea9700903455@mf2.8cd2.platon.network:16789",
			"636f329f8d7ff423316d0aea154b3ed4619e97cad4d034c368d6523b08e507ef522e9e0375c6073495c5cf843cd10e06b42653fd0b9e5101cfbc9d97061a6fedd5e2dd644feae4ceab6fdf473e63bc3a4875c403e3a2cf50260968777952dd92",
		},
		{
			"enode://923de68a0478a335f455c77758855a80ebc40c4d6693a9d438e8aa340274045614d13d60559560a9ccea888d745df689a042f8c17bb9941e0d26ed8373a535c3@mf3.f9ec.platon.network:16789",
			"59093810db7d504b997f6b05487e9a9da14155de9934581d7ab02dd60348ac0f7612c92a2038f762a7e0983c4e005e05db9ad91660f4add32c331fce92880087b78318deedd87bfa56fa3e6e0a47c546015a1aad44e77ce6b13043a56f43138b",
		},
		{
			"enode://95e7214df1659235be6dbdc8da93b04a16aa49b1f44378d2a14d216d2b68698264f020f0903dcada964e21809a54c9983b35cc309946dfe370002de57050dae4@mf4.add3.platon.network:16789",
			"448778b150b31e8c2ed7141465772b44e6a214c84dea80e0ae2bead52b694b10db93582804cdaf97dd9404b85ae90c0265e2722033fdfef79d28bd40b0f8dd561642a5a787954b7619b2ec038a637b2d95885d892cf3dfec679db5c2c58e6183",
		},
		{
			"enode://3fd5d96881f1f6139aeabb08509ac6c62f1666034d1417038d2781dca9e4a6c06313be53f235cb8cf67e25d7915e2377c4720207c83223720792edfc0ec35782@mf5.0354.platon.network:16789",
			"6d1727eca4926fd81ee61d130562b78b1f189fdb44e33eb1f722df4bf4a20d9b9e587c33f9b89fcd3fd0e811cffffb06391941bcd0c9670846db8b2bedcdf3e8a7b8141bde62f249445f0bf12f731daf44f2b10e37e51d5197ec5bcbb44f8415",
		},
		{
			"enode://4b3f7f01f957245465a955ccaaa34b34c2b8919e376183ddc2ea2ed875e1933c38a09ce97c4c8014b6ac22593d4c3be4f6533bf132f95f2daf9ff01c1c790c49@mf6.9e2c.platon.network:16789",
			"5e20c9c72b0b571157815009c77a8911cce10354bc6e05b06e4c4d56a1ae9f9cff072a5bb0fe5924d34c11cc92230718ad92b9a4a3742624aedadc57f22d49684abd5d106ee7d1311be681fe2b949af821c1a33e215a0e96dc0ad32e87fd3105",
		},
		{
			"enode://e017ef712a9aac70c560995d0e806303f76f8d27774a7458ccb667735fe5ecf4963fb3596c610627d5e184a24710230907c6702d30eef7291c75d76712cb9427@mf7.e1dd.platon.network:16789",
			"5e15780e3fc248715ec8588863043cf9d921afbeb6c7f0ca92aa17d0403f94c28ea76898d8b82ac590511391dd7a3f1675372250355be8be17a80a47aa51f38a04fd2965d7e252625d4fc41ab0b41d004623d8538c50065406dcd6d18036d382",
		},
	}

	initialTestnetConsensusNodes = []initNode{
		{
			"enode://b7f1f7757a900cce7ce4caf8663ecf871205763ac201c65f9551d5b841731a9cd9550bc05f3a16fbc2ef589c9faeef74d4500b60d76047939e2ba7fa4a5915aa@127.0.0.1:16789",
			"f1735bac863706b49809a4e635fe0c2e224aef5ad549f18ba3f2f6b61c0c9d0005f12d497a301ba26a8aaf009c90e4198301875002984c5cd9bd614cd2fbcb81c57f6355a8400d56c20804e1dfb34782c1f2eadda82c8b226aa4a71bfa60be8c",
		},
	}

	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(210411),
		AddressHRP:  "lat",
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(1),
		Cbft: &CbftConfig{
			InitialNodes:  ConvertNodeUrl(initialMainNetConsensusNodes),
			Amount:        10,
			ValidatorMode: "ppos",
			Period:        20000,
		},
		GenesisVersion: GenesisVersion,
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetTrustedCheckpoint = &TrustedCheckpoint{
		Name:         "mainnet",
		SectionIndex: 193,
		SectionHead:  common.HexToHash("0xc2d574295ecedc4d58530ae24c31a5a98be7d2b3327fba0dd0f4ed3913828a55"),
		CHTRoot:      common.HexToHash("0x5d1027dfae688c77376e842679ceada87fd94738feb9b32ef165473bfbbb317b"),
		BloomRoot:    common.HexToHash("0xd38be1a06aabd568e10957fee4fcc523bc64996bcf31bae3f55f86e0a583919f"),
	}

	// TestnetChainConfig is the chain parameters to run a node on the test network.
	TestnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(104),
		AddressHRP:  "lat",
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(1),
		Cbft: &CbftConfig{
			InitialNodes:  ConvertNodeUrl(initialTestnetConsensusNodes),
			Amount:        10,
			ValidatorMode: "ppos",
			Period:        20000,
		},
		GenesisVersion: GenesisVersion,
	}

	// TestnetTrustedCheckpoint contains the light client trusted checkpoint for the test network.
	TestnetTrustedCheckpoint = &TrustedCheckpoint{
		Name:         "testnet",
		SectionIndex: 123,
		SectionHead:  common.HexToHash("0xa372a53decb68ce453da12bea1c8ee7b568b276aa2aab94d9060aa7c81fc3dee"),
		CHTRoot:      common.HexToHash("0x6b02e7fada79cd2a80d4b3623df9c44384d6647fc127462e1c188ccd09ece87b"),
		BloomRoot:    common.HexToHash("0xf2d27490914968279d6377d42868928632573e823b5d1d4a944cba6009e16259"),
	}

	GrapeChainConfig = &ChainConfig{
		AddressHRP:  "lat",
		ChainID:     big.NewInt(304),
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(3),
		Cbft: &CbftConfig{
			Period: 3,
		},
		GenesisVersion: GenesisVersion,
	}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllEthashProtocolChanges = &ChainConfig{big.NewInt(1337), "lat", "", big.NewInt(0), big.NewInt(0), nil, nil, GenesisVersion}

	TestChainConfig = &ChainConfig{big.NewInt(1), "lat", "", big.NewInt(0), big.NewInt(0), nil, new(CbftConfig), GenesisVersion}
)

// TrustedCheckpoint represents a set of post-processed trie roots (CHT and
// BloomTrie) associated with the appropriate section index and head hash. It is
// used to start light syncing from this checkpoint and avoid downloading the
// entire header chain while still being able to securely access old headers/logs.
type TrustedCheckpoint struct {
	Name         string      `json:"-"`
	SectionIndex uint64      `json:"sectionIndex"`
	SectionHead  common.Hash `json:"sectionHead"`
	CHTRoot      common.Hash `json:"chtRoot"`
	BloomRoot    common.Hash `json:"bloomRoot"`
}

// ChainConfig is the core config which determines the blockchain settings.
//
// ChainConfig is stored in the database on a per block basis. This means
// that any network, identified by its genesis block, can have its own
// set of configuration options.
type ChainConfig struct {
	ChainID     *big.Int `json:"chainId"` // chainId identifies the current chain and is used for replay protection
	AddressHRP  string   `json:"addressHRP"`
	EmptyBlock  string   `json:"emptyBlock"`
	EIP155Block *big.Int `json:"eip155Block,omitempty"` // EIP155 HF block
	EWASMBlock  *big.Int `json:"ewasmBlock,omitempty"`  // EWASM switch block (nil = no fork, 0 = already activated)
	// Various consensus engines
	Clique *CliqueConfig `json:"clique,omitempty"`
	Cbft   *CbftConfig   `json:"cbft,omitempty"`

	GenesisVersion uint32 `json:"genesisVersion"`
}

type CbftNode struct {
	Node      discover.Node `json:"node"`
	BlsPubKey bls.PublicKey `json:"blsPubKey"`
}

type initNode struct {
	Enode     string
	BlsPubkey string
}

type CbftConfig struct {
	Period        uint64     `json:"period,omitempty"`        // Number of seconds between blocks to enforce
	Amount        uint32     `json:"amount,omitempty"`        //The maximum number of blocks generated per cycle
	InitialNodes  []CbftNode `json:"initialNodes,omitempty"`  //Genesis consensus node
	ValidatorMode string     `json:"validatorMode,omitempty"` //Validator mode for easy testing
}

// CliqueConfig is the consensus engine configs for proof-of-authority based sealing.
type CliqueConfig struct {
	Period uint64 `json:"period"` // Number of seconds between blocks to enforce
	Epoch  uint64 `json:"epoch"`  // Epoch length to reset votes and checkpoint
}

// String implements the stringer interface, returning the consensus engine details.
func (c *CliqueConfig) String() string {
	return "clique"
}

// String implements the fmt.Stringer interface.
func (c *ChainConfig) String() string {
	var engine interface{}
	switch {
	case c.Clique != nil:
		engine = c.Clique
	case c.Cbft != nil:
		engine = c.Cbft
	default:
		engine = "unknown"
	}
	return fmt.Sprintf("{ChainID: %v EIP155: %v Engine: %v}",
		c.ChainID,
		c.EIP155Block,
		engine,
	)
}

// IsEIP155 returns whether num is either equal to the EIP155 fork block or greater.
func (c *ChainConfig) IsEIP155(num *big.Int) bool {
	//	return isForked(c.EIP155Block, num)
	return true
}

// IsEWASM returns whether num represents a block number after the EWASM fork
func (c *ChainConfig) IsEWASM(num *big.Int) bool {
	return isForked(c.EWASMBlock, num)
}

// GasTable returns the gas table corresponding to the current phase (homestead or homestead reprice).
//
// The returned GasTable's fields shouldn't, under any circumstances, be changed.
func (c *ChainConfig) GasTable(num *big.Int) GasTable {
	return GasTableConstantinople
}

// CheckCompatible checks whether scheduled fork transitions have been imported
// with a mismatching chain configuration.
func (c *ChainConfig) CheckCompatible(newcfg *ChainConfig, height uint64) *ConfigCompatError {
	bhead := new(big.Int).SetUint64(height)

	// Iterate checkCompatible to find the lowest conflict.
	var lasterr *ConfigCompatError
	for {
		err := c.checkCompatible(newcfg, bhead)
		if err == nil || (lasterr != nil && err.RewindTo == lasterr.RewindTo) {
			break
		}
		lasterr = err
		bhead.SetUint64(err.RewindTo)
	}
	return lasterr
}

func (c *ChainConfig) checkCompatible(newcfg *ChainConfig, head *big.Int) *ConfigCompatError {
	if isForkIncompatible(c.EIP155Block, newcfg.EIP155Block, head) {
		return newCompatError("EIP155 fork block", c.EIP155Block, newcfg.EIP155Block)
	}
	if isForkIncompatible(c.EWASMBlock, newcfg.EWASMBlock, head) {
		return newCompatError("ewasm fork block", c.EWASMBlock, newcfg.EWASMBlock)
	}
	return nil
}

// isForkIncompatible returns true if a fork scheduled at s1 cannot be rescheduled to
// block s2 because head is already past the fork.
func isForkIncompatible(s1, s2, head *big.Int) bool {
	return (isForked(s1, head) || isForked(s2, head)) && !configNumEqual(s1, s2)
}

// isForked returns whether a fork scheduled at block s is active at the given head block.
func isForked(s, head *big.Int) bool {
	if s == nil || head == nil {
		return false
	}
	return s.Cmp(head) <= 0
}

func configNumEqual(x, y *big.Int) bool {
	if x == nil {
		return y == nil
	}
	if y == nil {
		return x == nil
	}
	return x.Cmp(y) == 0
}

// ConfigCompatError is raised if the locally-stored blockchain is initialised with a
// ChainConfig that would alter the past.
type ConfigCompatError struct {
	What string
	// block numbers of the stored and new configurations
	StoredConfig, NewConfig *big.Int
	// the block number to which the local chain must be rewound to correct the error
	RewindTo uint64
}

func newCompatError(what string, storedblock, newblock *big.Int) *ConfigCompatError {
	var rew *big.Int
	switch {
	case storedblock == nil:
		rew = newblock
	case newblock == nil || storedblock.Cmp(newblock) < 0:
		rew = storedblock
	default:
		rew = newblock
	}
	err := &ConfigCompatError{what, storedblock, newblock, 0}
	if rew != nil && rew.Sign() > 0 {
		err.RewindTo = rew.Uint64() - 1
	}
	return err
}

func (err *ConfigCompatError) Error() string {
	return fmt.Sprintf("mismatching %s in database (have %d, want %d, rewindto %d)", err.What, err.StoredConfig, err.NewConfig, err.RewindTo)
}

func ConvertNodeUrl(initialNodes []initNode) []CbftNode {
	bls.Init(bls.BLS12_381)
	NodeList := make([]CbftNode, 0, len(initialNodes))
	for _, n := range initialNodes {

		cbftNode := new(CbftNode)

		if node, err := discover.ParseNode(n.Enode); nil == err {
			cbftNode.Node = *node
		}

		if n.BlsPubkey != "" {
			var blsPk bls.PublicKey
			if err := blsPk.UnmarshalText([]byte(n.BlsPubkey)); nil == err {
				cbftNode.BlsPubKey = blsPk
			}
		}

		NodeList = append(NodeList, *cbftNode)
	}
	return NodeList
}
