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
			"enode://e0c521bc365d6e4be8e7c57eef4d82c41fa679f3531f4860a0a9bd10f1d03772a600f293f36baa037f2378a32aeed015b8781c17ca80eb93b33d15f0b8d6b0eb@137.135.187.49:16789", //TEST-SEA
			"1fae6a56386e6ca3f9bc3727105adff78c4822c5cdca8594edac3d5f4e145b3ee6b3fe995072cdf3bf40e7decfb35013d330676ac516e09d74695d8a5248a585279172447e1855bc28e04efea1efc7ae54a60ac4445dd168d3be3836d4ecb304",
		},
		{
			"enode://6d53868800442172dad019e11e46f2c311e82fa43d8b6673800d1ba96b4ccdfc79d2ed0e76bc2f3adc79e92ea2d66123bb8ff7229a0a62b3dacd2fe3f292f237@52.236.129.246:16789", //TEST-SG
			"7638475108d18c280516ab0afd8991bb7d300e174cefa568530c7812158a29906009a63a37b590904f1639919bbeb909f5e871efd062aa11444dfe8f28ef3646fe260d91d8b27087884f6b678309bf4c639a4e3f7ce6fbbb2c247f3515c2750e",
		},
		{
			"enode://3241190f8f9cc4daa6ad82b703719a9c374d681e99f74085f86c1ee7f8a8d2ec7ca2a691182441c250a376917d4e41342bb168f7172b2b73d24cd2d1cbaaa13c@52.228.38.97:16789", //TEST-NA
			"2cf585988d27aca716e33f30eeb2b5f7f5d146aca2f1a854e4667bb305bc0b8e73f756c6e81ac5c90690b36ecab89f15621ff49cdcaadcdfa6c78f28aeda0c498cc3502f1a59d0366e4bf767801649a1ca96cfe8eb663147721ca98840a43300",
		},
		{
			"enode://8823be50eeaec8b4091b1e6874d0ca62b339d10e4fffb2bb335910bad1461818b10e121ede9cc76d7ecc666f70037ba8538a893c4c41f031577745cd281ce4cd@51.105.52.194:16789", //TEST-US
			"2db52bc9ec0ab826fb64aafb179284fe5b2a5a2263d95a9a44cf7d5bf81e913aca4e745d109a02e44612f04c633d350f4f3992ddd250901487618f9534178949859c1fbdf40af2959161de73ef48abf3adb7b6ebcdc0890f76c9447ed10d6b8a",
		},
		{
			"enode://3da6f3dd87b6ab69933df5c8c7ce3384f96e31993b31120f6dd60e48f5454b1deeedf65708617b7b7aa3b4d2852379197ece293ae2d37bf80488bbf8305f9538@3.120.24.72:16789", //TEST-US
			"e5c17a1ec1dbbc0805086e7e3b47ec7356726b5afffc827d3ddbbd4021d8b5c7301b9a6316ea99c713bc42edad77ae0e625be6dad67558f638891b093051b963608b58b831732916f9973a65f18fdb65429c39a0ce7b96d6d86ac14d333fdf91",
		},
		{
			"enode://c51941700cda66534a951d54431c81bbc2859a17bf6fd2dd5be9dbc0e1d6b2858146a3919b8ab3ef2eeb7b6c89e632616537d94301f8dd4a3f7d9d0cb15eabf8@13.126.189.133:16789", //TEST-US
			"f9374427202567be87a1397cb018e66b1d695453bd9fba4b7e4ab13d4a730b83251ed48bdc0e35e581715ee1dcc7d00cef04ab4e02acfef826c8c130989461063a324f68f10610dd6c6edd03293c7b082dad41c1d8941b2ad28f85cb70fc5594",
		},
		{
			"enode://93bdcf7dd8e579e61bf1d8d1c6be9ffdfc492b9828321ea41d61c3f00f1f6d957310b63abf4b93f03738b7ce9cf141637ce8eaefcbd12e06a528dd9233ac8bc5@18.140.124.104:16789", //TEST-US
			"5135e6790ed70084f74f8daa8ad3a84980d8571c8b793ee19d634677571716e4cf42af46f3af53e3eea342440c62e60586751c7d2a86d59f3e09adab6b76399db1e3200afe315043ee895f5ad963fa6bfbebd8739b7d3a0f19ad60903a4b9793",
		},
	}

	initialRallyNetConsensusNodes = []initNode{
		{
			"enode://8eb53f40f5444fc620a603d08f5b33ffe518edb7ff71537bd6ddc92f971a6ac11c9c1f5315eb50df837b3a420681b7553d34f12fce224e00a4757b6e1e59260e@13.69.9.245:16789",
			"de70742cdee73b4ca931c85a02ed876e6396faa553655ca5a72125ef54af008aefa62e9dc705a82710790ff5f2a5a500e965e8dc70a2e46b5454f9adc7314d3597b0dc9560c7841202c938a725cb21e1a837b4b0a3a934569351576bf2f58c0b",
		},
		{
			"enode://8c3b29b053b0b1d0c940b1d97dad7e41aa069cd19072aa96ca65ddd75cda8ab7f252a4f527f48dcc9809bbfab4866eb49bba9aa3d47ae70eda25ae0b022494f0@207.46.233.122:16789",
			"7bd67034a569a4b945b91b653c8c62563c78a9fbacc2aea3265762c74c7dd6099c40d0c2bb5c6b61c3feec6f8ca01a0971ce15cb06a99ffca918d0abc1a1fa64ba954dd7da4a9f2e3033d3397a10db0f75eb440f5b2afa2d16d547655e020284",
		},
		{
			"enode://1184b02c88ca5100dde8cb208781f30ee5743e6049271f3678d7b7ee927f8604bf9b1906ebbb5538fdbad9fcf7a07f519ec06923a7c5c3d28b82193a8e05e9e8@52.142.166.24:16789",
			"b49f00ac24357921f3eb1fc62daffb34c5200de13e385bfff18f42b91655cc1533b393dea381a6701154e764bf3487119df64f3a0c7d0f6c7543ef8246e28f1fec076fdb9fde748e9a3748eb2cd46d234a5c93fcc5b5ef8c8bb3a05b4d1f2c17",
		},
		{
			"enode://638c147e82f34b38b17dfee8d7fdf3a4b0cbeeafff0261850bd7ffbdd49cd309710d1220855ff7a5c47837f0bb43b21ef0784fd532f43b5675ea553a4aeb696a@52.143.129.153:16789",
			"55b1468d1c5b46aee64795d04681bbe0d95fb52f616105341e6fb5572ae6530cfc0635ce9c1e2d39f323e7bd0721b609f8f5df8ddd2ba04a0d6e154460b44be83e426d2e546d865930eabc3240ee880c71b3df4db566e58b3fff18dc7ccc9614",
		},
		{
			"enode://9eb996af791008f5302eb9aa96db400dd17e3dfdd81f6ef3e2a9477c43ab736a767e371bbed0aa56a17e443850ea4854436678903214f4d4e1c3baecfffca416@52.63.239.155:16789",
			"acd829b1b645893741776cba1debad8023f166037a593e2fabaae7665b422074d6cf06e048258f0bfaad97c1702f361463b8aece65d43d3b945f92c5ce649eb019c087e634d4cc328c5269350620850298b467befabd9d4e219bc95a6b74c397",
		},
		{
			"enode://33d96f68153ca98ae3ac88f61f9295263962418b48c1dc5d71928d15edf6dbeeefdf2d05c3502ccbb62f1101098ae14fe79e4271c316a535b215f93658a3732f@35.182.73.204:16789",
			"6ec2b2d5ace2c6c858fbcfd228a9cad5a1cdd34be42b7c5deec5aff462308a9acd90b51a987be164ccf5fd78950b49058a370d89b1b6a6011419b2518f061bc491499692e9d627d32a0e026df8cfa73ddbc6dc508275999ec2f3b1b41ce51e97",
		},
		{
			"enode://d24bdf6948435427c2844aa0ac94e2691d52f6d82be7ba8bca657955eb386db633043ccec0cbe4fb98b65d0b4a6266d67195d1506900f3c92e98bb1a3de75039@3.122.68.59:16789",
			"4d8fd606967215b6f5e44e94a7336d9c1cada0d8b84c06e4316b63b42fc0dd0bac48d014f8c02484600d178f5394d219d9c108166ee7011b2d3549f3571207bec5a55fa23911b9778630e3ada7a6b3ae932e3dda5be904202d5fe4c36f57b307",
		},
	}

	initialUatNetConsensusNodes = []initNode{
		{
			"enode://ef4f6f74b0883637c54c257c530e22dae0e55902be086b2e059651b145a2afb549dd71dc59c1c302c32d9c41464aa5bdbd24c77d3c42fed679cf514a66b197f6@23.102.22.162:16789",
			"8b4d0c3b58aedce24ef254ea778247466f8269637bf70d16f811bfdc5274cafd8a28b716ab0bc433b94b4b34fe2fe50f9ce71b4b10f96a7d71ee0d4bb41c1de4d8aaa32004578344dfdc45c45d9956d47a5c7ffd21afe3cf76074d432082ab81",
		},
		{
			"enode://1084a3e3f76a8d15a811d5e3c564125e43d96a2bdc44f4f773f9eb2ac5c13d63ddc2c03ccb8743335575c6d6b5adb62345abe259b9be8298657cd5ce7fdad2fa@54.251.161.158:16789",
			"9df5ceceef3b2c3f601cd84911c92544b8dbbabd31f007324512ee7c88b6dff651f6a2263f50d9beda4335e23c06ac160dec487031e0a3c257f4215fd25336dbcbc40ecb884366bbe65c050343265a82a8e2f2494aa0c906f82d00fda7591e96",
		},
		{
			"enode://e64ff3d35bbcb37a991b16b7c120459ef9cf3c96ff5ab3e5124fc75d60f0ce7c280683377e46fed0a5c39913dfae8a2c398434dba9c364907406fbc0afb13e6f@149.129.136.253:16789",
			"5bf31a518acbf66513c8a5df86f78f2142d1cb9a26fbd10864637e851f1e6bf05848a604f0464ff809f98701e2381e0b99b1e654eb8b68da8fede0a2380571caced833da5f3fb21e56537c34353af89d899717ff7763e0fd93f2ac1a7c37e88d",
		},

		{
			"enode://db3b21aeed13527ad2a510762c85c98bd77f4a10d6c3396db333c0fd69063f198f91d5ed2b19f1d1eb05d3f4184a25689106b6f94b583d9e5d07909768cdefdb@149.129.136.60:16789",
			"12ab9ae8a51746ea7329cc4c0c31eb73ee9b00ca35efcf34762baa3e6b108a95936affaabb198a4db33f56dd037f28028e6466dcfb8569a76e32a119426ea048d4c0310ce1b0d2268c750fc907a268efdf313af5de4c741b2453617fcddf0383",
		},

		{
			"enode://fd22593e5a72d906dcf2df85aec9dc58b7101468493c547d238d17ad20cf1c28308b15b2521bca03f28af4706a0800f9edac5d4ab4fc89e5f109997ae294ef84@149.129.147.224:16789",
			"2f7cde8789ca182c15dd26a7a2fde303c16c9824072b5a20519cb650ad15756289c28c9239eed8f8cec1ab1f8646f20b981dcfd0d9bb9c80979f7c56c6b1d4ed5c94b88bb232ca3c92f69963c17768499a10023add6bab0ed084954bc9886a10",
		},
		{
			"enode://8b01ba20ab81e9753ecae32b0b6d35104957b91ab748fde20627b41361811a9c34e2cdca658fb1602871acdbf091dc8c00a97b87b5b2c33790bd8c57949a7675@149.129.132.215:16789",
			"2dcb3f2059802a10d9d9f918fa20d89ef3a2bbf80a356583e1159de3fcf103b6a2e878aff77519c2eb2c63a2ae0084106797cd267a5d8865e5f76d7e27854255cf04139c7e7d472a441f3f7835a54db633ba3a9299947488b438455cb2c2d603",
		},
		{
			"enode://defa79441cc774ec65c2eb78776cd8e12496d85fad00d9ef61a5454a0b62c5220d8b276f74ac81c8c19a61c6130b2c2a3bcb343a0426e830b2fc29a6cb083431@149.129.176.49:16789",
			"1cf7d7e829f66e826277814f44112019bb107509a26da1e0ae05fa2c6d4e251c01d1b861981382926d6dbe93ffcb7a078c5183553b79137bdb497076baaf01fea8867888860b1e4b6c4e9dcd6a44a321650329f953e69330b8f42886a09b650d",
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

	// TestnetChainConfig contains the chain parameters to run a node on the test network.
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

	// RallynetChainConfig is the chain parameters to run a node on the Rally network.
	RallynetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(95),
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(1),
		Cbft: &CbftConfig{
			InitialNodes:  ConvertNodeUrl(initialRallyNetConsensusNodes),
			Amount:        10,
			ValidatorMode: "ppos",
			Period:        20000,
		},
		VMInterpreter:  "wasm",
		GenesisVersion: GenesisVersion,
	}

	// UatnetChainConfig is the chain parameters to run a node on the Uat network.
	UatnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(299),
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(1),
		Cbft: &CbftConfig{
			InitialNodes:  ConvertNodeUrl(initialUatNetConsensusNodes),
			Amount:        10,
			ValidatorMode: "ppos",
			Period:        20000,
		},
		VMInterpreter:  "wasm",
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
