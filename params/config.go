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
	initialMainNetConsensusNodes = []initNode{
		{
			"enode://6de1c55250563c4eb75f7ad3dd54cf5f6b86150bd8613b171ee28e1abc0c4936cf1d0cbc3939734fde31a147af226f98c33bd14d52be23cc2b34b5ba6e0d8765@mf1.a4f9.platon.network:16789",
			"555d942d5cc88f3891933986d6dcf9b597f9865c1c8443666ceb0b9097b54a31473a036ceac47e714eaaaeeee2775f15a190449b1fa64cec4072b9407d685b196900ed7a30cd428715ac0f5a4f9c688d66b605daa5efbcab4b188f605a27b502",
		},
		{
			"enode://3b416a1fb1e66154087e4d98960a24935ab203da330d2fb7be44a1ffce0861cade982a83ff541cd9eef4f0b0b950449be290fd993279973d21a411d3832a52d4@mf2.8cd2.platon.network:16789",
			"83268442e50daf719eab1c97d5d9928fcf02bbd6912af904d1d6bedc6c83943f8e9165bf38ffda995f4ff7b89fac5d104141feb5babe40428137701eff3a2289b9aeffbcd2d9160dee560378cd2ee46b90c70a2bdf3d8b62c76bf197f33fe407",
		},
		{
			"enode://e0d5a4a8732b411040a525db2780cdd8dd45e434f4e8d80400fefbe9b48079e8490ea9e56327c5ae76646107995a7f44d655821818a11bbc0ba6026ce480ccc9@mf3.f9ec.platon.network:16789",
			"8d28bef573a9df55a622419e48a6ec2a8a7719c1750ce555e19e5f41927677b5d2885522d9186c538e665a3ada26c50fc52cbe95f4889025c65b8300cd713a2e5afbe7ee4d4953c6dad73adf9ecc5a696cf84a6dd1e1cb3cf7f57b7bb9ade813",
		},
		{
			"enode://7ebbb666b11d7d63f278308a4750f3721d55bd2fc5115f529fbae12792053672ae7309b1567cd0403b001fe00d678c17ffbc8dfb2a26a8365bd174d37e473d4b@mf4.add3.platon.network:16789",
			"99c17a47260722ecdb6074aa1f3d318173702daf930aaaff13b659bbd19f4a98eb10f6edd86ddd1aa588c0989676bf18244ed886a7ff02de930d2a85609710a1e26ae2031a8d88c6d5e2dc1add3462a94a6e2590e16986934c0e2e8d2354680c",
		},
		{
			"enode://3d5d15388828830ec67711b6d93785868693ddcf27bb8b02217c3dff5ee99b75c6287245a7dfd268f7d340acde44babe3c7ea0c98d0d8f156b985ec81714e293@mf5.0354.platon.network:16789",
			"aeb6804f9e1561338b2277b833349aebfbf524e164f96218f2b4c06a9c92ff59d06e0ff5866bbed81c43cfb16075e90911186f37b4678b44a11d7a17ede1d0aab223709d4b80ec4c09843d403544e16c865ce92f20f3455bf2b1502c3933b390",
		},
		{
			"enode://0f2a38f778e0617ed42b34d21d7b3d28a526a49f142269b0307de505b332c85b7f89f05ae3e39a5896a65928ea3d112e640853b5bbaa77d7c5763bfecb328f7a@mf6.9e2c.platon.network:16789",
			"bd8778181425e509f7a37604d2778511f1ede950751d785cbf5a9573623d006c1767e7074c7566088592d6d62106e5001e297a3f352159b5a6e78ed7d55b4691b6f0eade7b25240b2a468099e2c1adb3fbf24bc2e849d476574c7c91e597f684",
		},
		{
			"enode://f4ff6f37bddb21c02bdaa5a66f0c6572c6aaa0493fdc44cc341d3c2623379c1210862324439a6ea09e11635fafd535d0c3294ba439065b6aa4af50887bece72f@mf7.e1dd.platon.network:16789",
			"a298740a0cc11d32d1fe5088db78b9492be0bbefabc85e723f43f715c587eb71d5509122d87c63c61f513bd4f436e607547c249e05211e007d18909e7a275fe973b341099c9d8b540aea15fe1dd04eeb5ebaf2395594626ffe85b419e0410a00",
		},
	}

	initialTestnetConsensusNodes = []initNode{
		{
			"enode://b7f1f7757a900cce7ce4caf8663ecf871205763ac201c65f9551d5b841731a9cd9550bc05f3a16fbc2ef589c9faeef74d4500b60d76047939e2ba7fa4a5915aa@tf1.8e1f.platon.network:16789",
			"f1735bac863706b49809a4e635fe0c2e224aef5ad549f18ba3f2f6b61c0c9d0005f12d497a301ba26a8aaf009c90e4198301875002984c5cd9bd614cd2fbcb81c57f6355a8400d56c20804e1dfb34782c1f2eadda82c8b226aa4a71bfa60be8c",
		},
		{
			"enode://3b2acc72f673173a97295728f6d9f93a8d75d1c615455f3ec3fdc4471e707e54935d89ba4736082db4b618b05de203fe828a8182b0f1ba09d495ad9a8ddb418b@tf2.dd0b.platon.network:16789",
			"bd0d378e9d87e552d6f3842b38e30776e45eb14af3462822cbadf5eea492477dd764d10cc73c521e05fa90c6146dd70a6ae1bd25f473147b91ddd65e5077ead12a1010de2714ef7977067df4b519ab2f7d50db7c2e150dc2d5cb2bb6cc30e485",
		},
		{
			"enode://184fb0464cd84a28e6c9aec564c058113b6c93ae80eccf1dc50be0481fd27bbfea3dfef492898987aeabba07fd0b5b7048a88163da348b729b0b56b6184f6e6c@tf3.4984.platon.network:16789",
			"fd89eb74d9277b8a02eb32f2144bf571d7badbca85ea7ea1158f47f268e36119a8e83608b4b972a5fd332bdc4f0306149390ed0afa4f6bae3e18f80a6062feff0a53cd265c38a3bff43afdf93c06dd8cdd0928a57039aaf59712623cd412a38d",
		},
		{
			"enode://31df231e05e089ee517f577a4ded288d210fc9d313c4ed862a7a758e4a74d8ae4b84de69be837371120a521199bce5ecf8665a56cd1efc2c8b665dd563a3f590@tf4.3c1b.platon.network:16789",
			"1c6792a7868006106547d48fb77e2941240352fd5def1f447b38c88585655e02bdf91209b48db698866a92b59914d00cacbb4206f99ece4e78fb82a12c732039d5a17e23b9980cd26421f953c3017199dbd41b62f13227d69123f431fdae3506",
		},
		{
			"enode://50ed65eb0180f771bbef453f0750bb324b626e6283b10b5203cae28e2e50744ba4bae51242f5d5dd17ce8d5bb71c50599f0432fd64b6a64d8026de181d8c5ae3@tf5.6021.platon.network:16789",
			"0e82ee8823646e871d319a5e418c5845af4efebbffc873ace0b27ff1391c720fafb98a118447cca88efa2653cb1d8a1443fd38079bd8385e6291a52776158882360657b3c401c1a1dc2fb2cfc083247f454ea160fc6aadecdbc1b44ee7aa6f8c",
		},
		{
			"enode://487a68847885a7b198b24920c1a5addc82d7c62fb43a503d7a4701a9fd72deed5e4f53116d8cca139d66731a2dc7537d5234d400188710864b83a7ccd65ab74d@tf6.f6a9.platon.network:16789",
			"a41e72db25f7e4d2af738e6170d5d84c72c69a185be75f55a62fb5ad0b931061f8f32c13d851f96d0b819f67d079c205b12d73507fa156967c8e6f385c2ab5b373f6f4a5e3e86bd95fe9857001a1feabcdf9507b45f446a7066213563a8d7a81",
		},
		{
			"enode://081da72f53a7cc4eaeace9f8a2960c1c01bf606153fd5acf47f540e0f2a1ebe5aed36fc8bf10263308b18845fe0079e027c85af31e03f995b6834d88f8461a19@tf7.f361.platon.network:16789",
			"77d174270903955fd8b38de00c05c17b078031722e1fc414188251a5da34aa8025dc3cb75bebedf5cd8eb91d25bee213a5fd78f12ed77892bd2f2e4a21595a413280e37d1d1726bf9fa62e042307de89ee6e243304ef768bac3a9d466dc3a406",
		},
	}

	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(210324),
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
