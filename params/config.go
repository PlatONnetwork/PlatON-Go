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
			"enode://b7f1f7757a900cce7ce4caf8663ecf871205763ac201c65f9551d5b841731a9cd9550bc05f3a16fbc2ef589c9faeef74d4500b60d76047939e2ba7fa4a5915aa@35.204.79.107:16789",
			"f1735bac863706b49809a4e635fe0c2e224aef5ad549f18ba3f2f6b61c0c9d0005f12d497a301ba26a8aaf009c90e4198301875002984c5cd9bd614cd2fbcb81c57f6355a8400d56c20804e1dfb34782c1f2eadda82c8b226aa4a71bfa60be8c",
		},
		{
			"enode://3b2acc72f673173a97295728f6d9f93a8d75d1c615455f3ec3fdc4471e707e54935d89ba4736082db4b618b05de203fe828a8182b0f1ba09d495ad9a8ddb418b@34.65.187.110:16789",
			"bd0d378e9d87e552d6f3842b38e30776e45eb14af3462822cbadf5eea492477dd764d10cc73c521e05fa90c6146dd70a6ae1bd25f473147b91ddd65e5077ead12a1010de2714ef7977067df4b519ab2f7d50db7c2e150dc2d5cb2bb6cc30e485",
		},
		{
			"enode://e594018258f640bebe92f650ef745de0da5753d6fc4c91cceb3611960e632a03decfd9712ebca1578695e5396a0f3d1f2d605105fb3bac08541278f00774acfe@13.211.42.206:16789",
			"67a318be5df24ca95e942c286e985e7ad64e72ec13e6117163c6fec47c1048f1b3530dff00d2f2727ad04f1148fb5b0b3fc460b22670fb7cbfc8ca30d7bf9eee8f9499c7fdbb9ad4da8329c371117f9d27a110df0048d4507e20e9a8a29d9614",
		},
		{
			"enode://730ba88e0a779b0601981f59c4c1a01e479d967be5ca5f5847c9d28b4c7c8a6ef652cfab8ca5d618930b3cb5f775f1ada3810780970081eaf61829713c6fc815@3.6.118.118:16789",
			"5a7fe973b25a934dc98d04d0422b11bb0151f223e31f639d104791ec23d11f577c70b41c6c11e918dfd3e58cc3b98404e85d7a684154b2d90708e400772cc0f6cc23175ce8702b9850521e048dfc10288a983882aaebc43c95bbe0dfd2168097",
		},
		{
			"enode://2fb7f8a8a2b0d612a51bb67ebf2a40cfbf543e03e7dff5c1d8f1fe7e27017e2decb907567e32b0da1fcbba06cd4bea3d65334c301c8688f602718fe8f8eaf14c@15.222.177.170:16789",
			"3428febc87735c21adbe763903c999cf8e4ab29d06b5f3eb89171414503df9cafa4f91903139c22616245014e05d8b0a51aee707903e52c311a6901792533fa360cb633b5366137025d6ffaf8b45f418947e387ad892610f8bd371264ad1b898",
		},
		{
			"enode://540b1e973ce8890976d5fbbd62798f5104144c6ed83f37b30b5b55d4605de876b1db7a0d0610c5a537c4df3cf3d30a7377aa3d080f8a2eaed8cafa28a46d5375@35.176.104.241:16789",
			"9f539ecd69bc35581c3fa228840003028961187ff0165853a4feabda303e4c83d0193695e9856f97f59cbd7ffb96a1036c28ac55a97bedd55c7c0b7a2cef384ea04380d8a8cffbf0372f356ff845b244d41314cf158ba91d3d88b0db2a04bd14",
		},
		{
			"enode://1380da3143604da44f1ffb27ce5d9b3383f9d25031a6f6aa99ffe4387a930dbe6ac56230fd0e5722fc91eb1ea0e7569c386d2534eaa7bf62385c99a4c1410a74@3.123.104.238:16789",
			"d0fbf0cfcb2ab4287ae675a832c67b5ef70613f0ce9668c66af11c514c6713fd0923b0e31dd2bb369d9c79f123976206c0c9a5e254fb2a0834d063f9146ca8b856068e87a4db2fa6bd792ae7806fcecd6a3d3e5d359044df734f5138c98bde0e",
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
		VMInterpreter:  "evm",
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
		VMInterpreter:  "evm",
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
		VMInterpreter:  "wasm",
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
	AllEthashProtocolChanges = &ChainConfig{big.NewInt(1337), "", big.NewInt(0), big.NewInt(0), nil, nil, "", GenesisVersion}

	TestChainConfig = &ChainConfig{big.NewInt(1), "", big.NewInt(0), big.NewInt(0), nil, new(CbftConfig), "", GenesisVersion}
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

	// Various vm interpreter
	VMInterpreter  string `json:"interpreter,omitempty"`
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
