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
	MainnetGenesisHash = common.HexToHash("0x7c76759a423eadfe7513055efb726679be07875a81e10ea6445238caf0f34a23")
	TestnetGenesisHash = common.HexToHash("0x47c7f3dea852b28c45d5fee6c950321486780065745179c7a296934e493f18e1")
)

var TrustedCheckpoints = map[common.Hash]*TrustedCheckpoint{
	MainnetGenesisHash: MainnetTrustedCheckpoint,
	TestnetGenesisHash: TestnetTrustedCheckpoint,
}

var (
	initialMainNetConsensusNodes = []initNode{
		{
			"enode://d2d670c64375d958ae15030d2e7979a369a1142a8981f41cb6aa31727c90a6af79ea7b8d07284736eec4c690e501d5e638a7dc87a646b0245631afc84f1d0c1f@mf1.a4f9.platon.network:16789",
			"a776f354346f1363c40afa949822a3831798cba0950a71e8cbbc1c396abcf0d65ac534518fbb96338a51708bf773af0cfd3f285374259c3edc22eaed9c270b9255596be5f6dc75c4df62defd6587259c8df123635706a8eddfeca9616f012881",
		},
		{
			"enode://8dae13a4b3417a3443ce8e552974cc5119ed2a0796b81773d2a1457fdae767f2c0698c0f71849fe708e9b2a2b2e326944f56bce11f3d6048660bd8aae7ff0fb3@mf2.8cd2.platon.network:16789",
			"7c157fa27e7491d166115580bece75aeb7b1efc5968f495cfb7dbfc32701462c7d0cbe7010a79bc2cf3861612d58c40957bc46bf255acf177f13efe2d5de76b32442baa902617dfd56294bb768e65b29f4525018972cc99eb80eba258310cd11",
		},
		{
			"enode://0e9b9887bec136c6eeb551072c4ce6c7caf5d68051b6a7c12f39637468e288e2bf4258cc907296f14dce1fc7cc01ab5cef51f7ce9ac4de3a9d4c6d4aa8dd166d@mf3.f9ec.platon.network:16789",
			"3795f3da2ea48578040becfa208b13b7a7d35ea110f1bf092c99486d893028170a531b530df6c7a8e53ceb21f7114a03a2cc31b6528daba17dce9c3c7262d19d0ea0c131353213172c41f81f3d65ad0404b23d9bb0763e3b90c4de697888df90",
		},
		{
			"enode://8568e4fa62b17c80befe49b27c1a4ba61d4b683d7ec59c568c5023cc8f14c87c756439d52b9a15af2750407e006e5fda515e0b68c5092eab297194e3699e936d@mf4.add3.platon.network:16789",
			"67a10e82707b65824333eab1a234b9ff4a171c51be0063acdc5a7cc79708ab7853f03cf13daff0bab6de17908dc64e132e5112e6273b55a1733cac31823279179a5e717b2ec0003e74356deb01c2d28f5ab451bc7dfdd5bc39f80a3cd1423485",
		},
		{
			"enode://e5a06147d5e33cb935fff7d31a78422583d3cf0438c045da0ef6581c48e2ca1597687400adc35f2d79a93125d3d249e6aa3b881d163331e4c405e6e2b5293036@mf5.0354.platon.network:16789",
			"12849c388b4dc0e94dcc240484d157203403df06dfbd5912554fc5c0b06cfd18eef4dec0cfd66f0144b2cca658d18208f002d546fb9c8d02bbfcd64d5d5bbf00d80867241329e4b3609319c0165109f0cdedaeb2d1246857ef3ea29951948098",
		},
		{
			"enode://e88cea24d93efe84f34a0edbe62a844aa4a15c560f1fc7c832170625aa93a9320c9c906bdf579de4c605ffe680734f7a26133f8196eff297089905e42955571b@mf6.9e2c.platon.network:16789",
			"1d0adc463850d70a5f664fad1107bb8e646422e2e93e9dd4c5438c96e14aa44e3c36e9fd762b5ad2fbefbf0bfb088612d4b2daa3a15ddc6a895ea2c1d768b395379c629768a54f27d53a0e517662fa2753607d6a59c4f520e8ef020c558f2d93",
		},
		{
			"enode://b6bb715a59ecc00220b0b4db038f51a586b43542d2063c0096b5c864da0057a47899d6e594094df428d52a52f68d1db9e9c172e8c0f2c4ed659ab74dcd5cfd6c@mf7.e1dd.platon.network:16789",
			"d3635a29e275a570049431953f7693a8c8ad444c3282267042f837f1b31403dd079e952f52d7f5ab5c48a6a1052530035b6a31a91114458676fbb925c20c34e73cb5cffb40d0ef80e48abf5359ecdfa7b69815fb8beb9c46d364004ba2d28500",
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
		ChainID:     big.NewInt(210421),
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
