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
	MainnetGenesisHash      = common.HexToHash("0xd4e56740f876aef8c010b86a40d5f56745a118d0906a34e69aec8c0db1cb8fa3")
	TestnetGenesisHash      = common.HexToHash("0x41941023680923e0fe4d74a34bdac8141f2540e3ae90623718e47d66d1ca4a2d")
	BeatnetGenesisHash      = common.HexToHash("0x6341fd3daf94b748c72ced5a5b26028f2474f5f00d824504e4fa37a75767e177")
	InnerTestnetGenesisHash = common.HexToHash("0x4f2a0d2caed299677b865f9e4c4420512d1ab9a34b8abbaadc3983668f67c5da")
	InnerDevnetGenesisHash  = common.HexToHash("0xc2d90f565b3c6dc16234c48978da62c2c20e15875644e9843d26a871860fd736")
)

var TrustedCheckpoints = map[common.Hash]*TrustedCheckpoint{
	MainnetGenesisHash:      MainnetTrustedCheckpoint,
	TestnetGenesisHash:      TestnetTrustedCheckpoint,
	BeatnetGenesisHash:      BetanetTrustedCheckpoint,
	InnerTestnetGenesisHash: InnerTestnetTrustedCheckpoint,
	InnerDevnetGenesisHash:  InnerDevnetTrustedCheckpoint,
}

var (
	initialMainNetConsensusNodes = []initNode{
		{
			"enode://94b2a01c3e894ce3b112bfaa5c59b680e90feb68bd036572675023711420e87b06cf3dd2506aa841279a61ff0bd469ac5f32cef8a21111ef57b16b5ff2c8039c@13.69.9.245:16789", //TEST-SEA
			"140843b34bea6182bd8c7ad639b2a83ca52c7fb9dd7c350b67ffbd3191e36bb4907edca3cfbfa57d67f75585cf013700fc5d094273e4460c06e07c933b6d75a67a4ae2d5c0f7d133681bb52ed05121e8f27a874d017b6c6e4043627f9c12e100",
		},
		{
			"enode://625327f3c60ba3688ba4bde25de91f258c4f51befd296d0a1c544ab36dd906bbf6e03a202daed37c4fc607b4dfa2445899a7ecea2bf708d33a2d4a2adcc45ba8@207.46.233.122:16789", //TEST-SG
			"261a6c9d90cc1cabafeb96c71473d532bf884d91a11e8ff48b260296e737a0011ad78c26ad0e805e00b543ccbe989f10350a9a2d113c98ea0dd227c0c53238ec2a5aa52103eb1fc22ae8b9a2e5dd3ef1e4337c4efeed09aec72ace312490b283",
		},
		{
			"enode://74f105b8283d657c4cf454b7abfe515b2357be69a8ac3e45bfacc8316d4a937d4a79cf79794d28b2d9662b448cb8cd0b2b356a35e7ff00f11f5a99da948a8576@52.142.166.24:16789", //TEST-NA
			"97974c4935b6a6f97d3369cf66e7928b7f490a92d962cb3022c50e61b6e1554a14e9d7a7c3353ed7568bf54b62b6f213e8fe9caab04c2dd457691ad047493428fb628356966b66d12ddba5f7b1c68ed7ed833d46471f0da39d831e373021b495",
		},
		{
			"enode://1d4c264111141f90ffcaf0c6d5561db7b99a633656f4a25366621e8296858b58b619269c9e18e0bed23f965ece75c4df918e88c0f39eead669262a24293e981b@52.143.129.153:16789", //TEST-US
			"d9fe561f324b04e71dfb55b602025b5a5154ae9909fafa8677b90b8b710e842a736e5127e1b4ccbff573f514a4229303372bab45a4e3079841de806915d82b8ba88e4f5bd3b3893f14c3e91d013aa3d038bbe063b42cea11c8b0f37a9d4e4618",
		},
		{
			"enode://9d3a8c1758ee4b46644bc1a40f93fcdb95ffb1081f793ee8345ccf58e0ad699efb981b2bcc736cac0319a12a7c1989d4c8ab32bd58fa5648fdb4dd9ffa33441b@52.63.239.155:16789", //TEST-US
			"729ecb8a450e08749956b3f8d7c6902c3395c051e0f9eb9f011aa73ee9506a73eb56691643f1e14fde9c27404443d819973944d99844fcec0322eb7385a082fef50a68f6503cb155b0b577499dd336b1a8da020170ffe80fd98400c93c841f98",
		},
		{
			"enode://daf9f6dfa8400b1e9b67166e032d95ec7b07c9780458006c09656f4ba150558ab3adbc0e6b1ff51954eb42ad62a9d3796d3ad4c5db764f7c1b79203ad10f861d@35.182.73.204:16789", //TEST-US
			"2735edde6102ca9f2f556f5d352c68b4cdb90c8d749a9bdd3ce6e9a0dc9c56b58fc6d7c1161a0f60b884cdba9cd86b17f770fb60306d9ca094a1a3d52325a94481fcae0d0e9a59369d869cbfc29e1aa6beb0397d7091bec9f2f2c9cf19a8db87",
		},
		{
			"enode://182be84e03ad9c8f679bbca734ec3fb94fe34858a7437afca663eee51c733b5b393a569db6949e5ca27fb1a9dc14545fbf30c12818caa8ff9d7cc456003e6a78@3.122.68.59:16789", //TEST-US
			"9ba9df42752da34b17d1106958d22e4d47f336febf780a1c400d79d0fd7d3e3c7dfdcc4fda2590ebc34aaefdd1e02b0c056f8199461c6e7159f685ced9e02dd75f75f5028ce891681798a7a291a3587eabb26c3db61e73734dc7a41b50802795",
		},
	}

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

	initialBetanetConsensusNodes = []initNode{
		{
			"enode://bcb7e49461cdd5f3227bb6cc6c36675cd936c11b69c3fd366c36997d514beabc423f8dfee6f91330a96273988bb68b1785161631181fd738d0f46d263b3ce8b3@54.176.216.82:16791",
			"",
		},
		{
			"enode://5449094bf985a688d378a90cf334d5a1abc55d694d6f2362899494d18048ef6b6bd724f4e51084bfe0563c732c481869c9da05d92e56f29f6880ad15ea851f13@54.176.216.82:16792",
			"",
		},
		{
			"enode://c0f7ae43af0605b80e35a5469adaa142059eaaf41d152613d74d42feffd6871f059f9ac4d596bd134bb1d6bbfbcea5391adff6f005ea9042c21797d51d0b7697@3.1.59.5:16791",
			"",
		},
		{
			"enode://b6883e86e833cec2405fb548405f7a1e693379f77ee8fc6bbf41b5c853d7ad654a2a3fb7ffbe57ae848509d1ed7e11acaf28666f8f81646eab575dafa8d51d0b@3.1.59.5:16792",
			"",
		},
	}

	initialInnerTestnetConsensusNodes = []initNode{
		{
			"enode://97e424be5e58bfd4533303f8f515211599fd4ffe208646f7bfdf27885e50b6dd85d957587180988e76ae77b4b6563820a27b16885419e5ba6f575f19f6cb36b0@192.168.120.81:16789",
			"",
		},
		{
			"enode://3b53564afbc3aef1f6e0678171811f65a7caa27a927ddd036a46f817d075ef0a5198cd7f480829b53fe62bdb063bc6a17f800d2eebf7481b091225aabac2428d@192.168.120.82:16789",
			"",
		},
		{
			"enode://858d6f6ae871e291d3b7b2b91f7369f46deb6334e9dacb66fa8ba6746ee1f025bd4c090b17d17e0d9d5c19fdf81eb8bde3d40a383c9eecbe7ebda9ca95a3fb94@192.168.120.83:16789",
			"",
		},
		{
			"enode://e4556b211eb6712ab94d743990d995c0d3cd15e9d78ec0096bba24c48d34f9f79a52ca1f835cec589c5e7daff30620871ba37d6f5f722678af4b2554a24dd75c@192.168.120.84:16789",
			"",
		},
		{
			"enode://114e48f21d4d83ec9ac39a62062a804a0566742d80b191de5ba23a4dc25f7beda0e78dd169352a7ad3b11584d06a01a09ce047ad88de9bdcb63885e81de00a4d@192.168.120.85:16789",
			"",
		},
		{
			"enode://64ba18ce01172da6a95b0d5b0a93aee727d77e5b2f04255a532a9566edaee7808383812a860acf5e43efeca3d9321547bfcdefd89e9d0c605dcdb65ce0bbb617@192.168.120.86:16789",
			"",
		},
		{
			"enode://d31b3a7714610bd8e03b2c74aca4be16de7fcc319a1e577d50e5e8796680221b4b679bf1c37966d1a158902b8686f3ca2f41a89a7176e538141082540c4f6d66@192.168.120.87:16789",
			"",
		},
		{
			"enode://805b617b9d321a65d8936e758b5c60cd6e8c873b9f1e7c793ad5f887d26ce9667d0db2fe55a9aeb1cc81f9cf9a1e7c54473203473e3ebda89e63c03cbcfe5347@192.168.120.88:16789",
			"",
		},
		{
			"enode://fa147bc3625acc846a9f0e1e89172ca7470baa0f86516994f70860c6fb904ddbb1849e3cf2b40c58255e38401f40d2c3e4a3bd5c2f2849b98465a5bdb80ed6a0@192.168.120.89:16789",
			"",
		},
		{
			"enode://d8c4b58ae052ea9480577264bc1b2c09619757015849a4c92b71a4e4c8b5ede94f35d24107b1181d0711013ed7fdc068f21e6e6084b3e96750a571669715c0b1@192.168.120.90:16789",
			"",
		},
	}

	initialInnerDevnetConsensusNodes = []initNode{
		{
			"enode://0abaf3219f454f3d07b6cbcf3c10b6b4ccf605202868e2043b6f5db12b745df0604ef01ef4cb523adc6d9e14b83a76dd09f862e3fe77205d8ac83df707969b47@192.168.9.76:16789",
			"",
		},

		{
			"enode://e0b6af6cc2e10b2b74540b87098083d48343805a3ff09c655eab0b20dba2b2851aea79ee75b6e150bde58ead0be03ee4a8619ea1dfaf529cbb8ff55ca23531ed@192.168.9.76:16790",
			"",
		},
		{
			"enode://15245d4dceeb7552b52d70e56c53fc86aa030eab6b7b325e430179902884fca3d684b0e896ea421864a160e9c18418e4561e9a72f911e2511c29204a857de71a@192.168.120.76:16789",
			"",
		},
		{
			"enode://fb886b3da4cf875f7d85e820a9b39df2170fd1966ffa0ddbcd738027f6f8e0256204e4873a2569ef299b324da3d0ed1afebb160d8ff401c2f09e20fb699e4005@192.168.120.76:16790",
			"",
		},
	}

	// MainnetChainConfig is the chain parameters to run a node on the main network.
	MainnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(101),
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(1),
		Cbft: &CbftConfig{
			InitialNodes:  ConvertNodeUrl(initialMainNetConsensusNodes),
			Amount:        10,
			ValidatorMode: "ppos",
			Period:        20000,
		},
		VMInterpreter: "wasm",
	}

	// MainnetTrustedCheckpoint contains the light client trusted checkpoint for the main network.
	MainnetTrustedCheckpoint = &TrustedCheckpoint{
		Name:         "mainnet",
		SectionIndex: 193,
		SectionHead:  common.HexToHash("0xc2d574295ecedc4d58530ae24c31a5a98be7d2b3327fba0dd0f4ed3913828a55"),
		CHTRoot:      common.HexToHash("0x5d1027dfae688c77376e842679ceada87fd94738feb9b32ef165473bfbbb317b"),
		BloomRoot:    common.HexToHash("0xd38be1a06aabd568e10957fee4fcc523bc64996bcf31bae3f55f86e0a583919f"),
	}

	// TestnetChainConfig contains the chain parameters to run a node on the Alpha test network.
	TestnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(103),
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(1),
		Cbft: &CbftConfig{
			InitialNodes:  ConvertNodeUrl(initialTestnetConsensusNodes),
			Amount:        10,
			ValidatorMode: "ppos",
			Period:        20000,
		},
		VMInterpreter: "wasm",
	}

	// TestnetTrustedCheckpoint contains the light client trusted checkpoint for the Alpha test network.
	TestnetTrustedCheckpoint = &TrustedCheckpoint{
		Name:         "testnet",
		SectionIndex: 123,
		SectionHead:  common.HexToHash("0xa372a53decb68ce453da12bea1c8ee7b568b276aa2aab94d9060aa7c81fc3dee"),
		CHTRoot:      common.HexToHash("0x6b02e7fada79cd2a80d4b3623df9c44384d6647fc127462e1c188ccd09ece87b"),
		BloomRoot:    common.HexToHash("0xf2d27490914968279d6377d42868928632573e823b5d1d4a944cba6009e16259"),
	}

	// InnerTestnetChainConfig contains the chain parameters to run a node on the inner test network.
	InnerTestnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(203),
		EIP155Block: big.NewInt(3),
		Cbft: &CbftConfig{
			InitialNodes: ConvertNodeUrl(initialInnerTestnetConsensusNodes),
		},
		VMInterpreter: "wasm",
	}

	// InnerTestnetTrustedCheckpoint contains the light client trusted checkpoint for the inner test network.
	InnerTestnetTrustedCheckpoint = &TrustedCheckpoint{
		Name:         "innertestnet",
		SectionIndex: 123,
		SectionHead:  common.HexToHash("0xa372a53decb68ce453da12bea1c8ee7b568b276aa2aab94d9060aa7c81fc3dee"),
		CHTRoot:      common.HexToHash("0x6b02e7fada79cd2a80d4b3623df9c44384d6647fc127462e1c188ccd09ece87b"),
		BloomRoot:    common.HexToHash("0xf2d27490914968279d6377d42868928632573e823b5d1d4a944cba6009e16259"),
	}

	// InnerDevnetChainConfig contains the chain parameters to run a node on the inner test network.
	InnerDevnetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(204),
		EIP155Block: big.NewInt(3),
		Cbft: &CbftConfig{
			InitialNodes: ConvertNodeUrl(initialInnerDevnetConsensusNodes),
		},
		VMInterpreter: "wasm",
	}

	// InnerDevnetTrustedCheckpoint contains the light client trusted checkpoint for the inner test network.
	InnerDevnetTrustedCheckpoint = &TrustedCheckpoint{
		Name:         "innerdevnet",
		SectionIndex: 123,
		SectionHead:  common.HexToHash("0xa372a53decb68ce453da12bea1c8ee7b568b276aa2aab94d9060aa7c81fc3dee"),
		CHTRoot:      common.HexToHash("0x6b02e7fada79cd2a80d4b3623df9c44384d6647fc127462e1c188ccd09ece87b"),
		BloomRoot:    common.HexToHash("0xf2d27490914968279d6377d42868928632573e823b5d1d4a944cba6009e16259"),
	}

	// BetanetChainConfig contains the chain parameters to run a node on the Beta test network.
	BetanetChainConfig = &ChainConfig{
		ChainID:     big.NewInt(104),
		EmptyBlock:  "on",
		EIP155Block: big.NewInt(3),
		Cbft: &CbftConfig{
			InitialNodes: ConvertNodeUrl(initialBetanetConsensusNodes),
		},
		VMInterpreter: "wasm",
	}

	// BetanetTrustedCheckpoint contains the light client trusted checkpoint for the Beta test network.
	BetanetTrustedCheckpoint = &TrustedCheckpoint{
		Name:         "betanet",
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
	}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllEthashProtocolChanges = &ChainConfig{big.NewInt(1337), "", big.NewInt(0), big.NewInt(0), nil, nil, ""}

	// AllCliqueProtocolChanges contains every protocol change (EIPs) introduced
	// and accepted by the Ethereum core developers into the Clique consensus.
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllCliqueProtocolChanges = &ChainConfig{big.NewInt(1337), "", big.NewInt(0), big.NewInt(0), &CliqueConfig{Period: 0, Epoch: 30000}, nil, ""}

	TestChainConfig = &ChainConfig{big.NewInt(1), "", big.NewInt(0), big.NewInt(0), nil, new(CbftConfig), ""}

	AllCbftProtocolChanges = &ChainConfig{big.NewInt(1337), "", big.NewInt(0), nil, nil, new(CbftConfig), ""}
	TestRules              = TestChainConfig.Rules(new(big.Int))
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
	VMInterpreter string `json:"interpreter,omitempty"`
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
	return GasTableHomestead
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
