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
	MainnetGenesisHash = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	TestnetGenesisHash = common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d")
)

var TrustedCheckpoints = map[common.Hash]*TrustedCheckpoint{
	MainnetGenesisHash: MainnetTrustedCheckpoint,
	TestnetGenesisHash: TestnetTrustedCheckpoint,
}

var (
	initialMainNetConsensusNodes = []initNode{}

	initialTestnetConsensusNodes = []initNode{
		{
			"enode://beb5e7c53bc85c7bc516fe97071bf53d947c9f2889a61944d6770970f869e680da944e08c2ef44fb11416b9afc2447f7253fa891e87a0e01651189713bb60d8d@35.204.79.107:16789",
			"f2aa6dd970e1ea05200e509f3bcfa6b5db120538dad1cc6f4bb345a60426091c43f2beeb7077a4031a4cdbc30f571c18973b9b5275394d0faa82d1f507625aaae0bdccb3e5612b2c5706a97c3f8457f990a0f57633ecb2f7903e3c608685e591",
		},
		{
			"enode://1e522c3f60430bb77819f2fbf0db044d87c19c2acf17c0f4c6ef29b5131511d0150836dc1dea3bd11e23da2d25f1628b2052d764b792f0619437866b94e1875d@34.65.187.110:16789",
			"3e2a1ea482652106d37f75c945ea7399b55207cd88da5f545b4120da80adcc8ff35166fa1a9edc61083b7baebd6dda003a302ecff3ea3c7ba655251d9a6a6130e12160771c0f793caa2cf46c263c1ceee48cbc31244b3cfeaed6f24b8dbe1001",
		},
		{
			"enode://dda92565a6d02af28f98d36ce25417e42a6ad6bb00fff0280300fd2d7fe127ad7ad265d681b35a6d0331fad644a5aabcf070a93c39309dad27e2083253cff1a3@13.211.42.206:16789",
			"192fdd69f21a899038f2fbb8731821ff2a672adc9137b6b0dfe9696d1efd86c0654e72e3cecd9ed8b832e6204869a506a2b5fd0bf4bd467b96c2cd553715e24d07aa9d5957e556d26c17e97ed505957a536cd01af70166707e6cd03fecc8b789",
		},
		{
			"enode://eaea0f8302c9818efff378f4fce1bcca3f9dccac960d882fc448f294f0985266a85f3507769f575221b374021e48611348b944cb5cbe3f2fe1218f898413107f@3.6.118.118:16789",
			"6243fb60d114e1c9d196979ac98c70c97c402eb39481416383b0fc04010721714606f4c808fbd4debe816abb6daeb0171db85b092eb88a1633d6753ae3cfa0b7193473d32caf18b849212a1e42d275658961886eb1769db71b7ca402fa8fc587",
		},
		{
			"enode://6f61b708dbcca26630f4a9811cdf1b54fc91da4e42d09c6e827a67720d0d201c7a7d0c4c410ba4cdf2e02328c0444fe9b7a90e506153ac0bb859cffcd2fe55fd@15.222.177.170:16789",
			"d6d53da1a081fd481647e3c75862c609ed0b1b6d06782e20f3ccad849f448802d7bbb3a7de11d328037f4ee6251520097d1487d7efaa3cd1adfbec9139b98bd77419a7e6a04a0bc5f7ca1272361a0998b5879aa3c4aaf4320513a39999cc8510",
		},
		{
			"enode://40683c74d5f0dbbf3daefacad5a9afb6384d3c9109e48bba8ccc61084e72005ccac984f6a946712fc3c1e77f26072315ab46046c31374b842a3d3949b14af4b7@35.176.104.241:16789",
			"3d07f9ff2b620c19eb2fe302d8dff4ba8f9687c38e923558cb37b7f77198c7ed7ac8405b27aacd45e20918ae588b6c1438e81a1df717f919969e4055633cd0d8da54b0a6d5a60a9fab6b080e7ba1f05a26e47722911b3174318a59b775422d0b",
		},
		{
			"enode://636607cdb83df0c2ed863875aa999ce0adf4a143061ff1f1e99a3200dcc16d55945870c73c641ea9e972105496d1f051789a7115e396e1326835dab6a8acee92@3.123.104.238:16789",
			"77022119b671104057983c7d9102d3ceb435b1c63e421c93077d6538ef5b7eeea6aadb7b943855760d4828b89e34af0035e1cb311f64b5ab84424fb4f8eec3792d39ad3d49e7adf0e4db3666b71ae77edef9409a44d9427059d8529b0ce47194",
		},
	}

	initialDemoNetConsensusNodes = []initNode{}

	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(100),
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
		ChainID:     big.NewInt(101),
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

	// DemonetChainConfig is the chain parameters to run a node on the demo network.
	DemonetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(399),
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(1),
		Cbft: &CbftConfig{
			InitialNodes:  ConvertNodeUrl(initialDemoNetConsensusNodes),
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
	AllEthashProtocolChanges = &ChainConfig{big.NewInt(1337), "", big.NewInt(0), big.NewInt(0), nil, nil, GenesisVersion}

	TestChainConfig = &ChainConfig{big.NewInt(1), "", big.NewInt(0), big.NewInt(0), nil, new(CbftConfig), GenesisVersion}
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

// Rules wraps ChainConfig and is merely syntactic sugar or can be used for functions
// that do not have or require information about the block.
//
// Rules is a one time interface meaning that it shouldn't be used in between transition
// phases.
type Rules struct {
	ChainID  *big.Int
	IsEIP155 bool
}

// Rules ensures c's ChainID is not nil.
func (c *ChainConfig) Rules(num *big.Int) Rules {
	chainID := c.ChainID
	if chainID == nil {
		chainID = new(big.Int)
	}
	return Rules{
		ChainID:  new(big.Int).Set(chainID),
		IsEIP155: c.IsEIP155(num),
	}
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
