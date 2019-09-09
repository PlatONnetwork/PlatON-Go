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
			"enode://b95630ee516effebd04c40eff7a4ea695947ed10e9e310ce01a12338f90748fcced68ba548a8261f7fce784e33e1ca7e0e6d7c54ef6d9887ed64f01727d72512@137.135.187.49:16789",
			"4d861e581b94b0b3e94a2be8f2e141e67fac978bd62344d424e0ba35113edb97bcb7966e6cae997f9bfdf4a5b229c407cb6879c6cefeff39223bcb2b1ef6b9f520ebfdd147dfa9f96285107be780076aab0a854beac882408beee1ba4deef105",
		},
		{
			"enode://3d0f7b74358a07562e89919ca274ec140ebc8cf6db4a241a6ddda6f896f0cfc37028a433e0efe5e6bd91b02d91883fd2c46f68cf77d08aa8f1f1c30ffcb2f614@52.236.129.246:16789",
			"7b7888e8cf571e297b7c4955386281e934b3136364af1560907eb95f70b212c9d47357e56e284f81014a624a08969f11c4257c10f6b056c3cadf544c649d2f32d5e3846b06f6cd34e9e6930a3dac5aaf41fd05c4b11212314f9f7177e9901114",
		},
		{
			"enode://2b2c71e0f5a5ce2f1fa8b16c429b7dbc2e3c5a9cb890e05dbbaf274959df1937c5c5cff1c613d017c1b06b2c24dc9ff464ee6fa03934a96be8a2faa3c7c08a0f@52.228.38.97:16789",
			"efecd9acb2fb1343e1be8c0bc8e879033bb670b26bee56526a15e9b31cc9ee2eb90ad2ffca175b0c6d047f1f681dbc00bf1c40a5172bd7ddc777438ca7e12442688f6dfc918b5261a6eb2ba43c0c49a5769d9b22f836eec2ebc682fd24245f89",
		},
		{
			"enode://60f682356cdac0d0a9b5cae8b5ffd2ae051e0767a059afd537aa1d47cd407925b1bdba25060d01af21a47089ac38f0b09d816f4ff71fcf3d570e71f1ef1f870b@51.105.52.194:16789",
			"76b1f46d2c488f26f12b2aea82bc84bacf1de2f85acc99780940f3b8845f802ce1dc7b5987acac7d82067221fa6d7106f47d986631899eeacaa5eadf4ba80b34b8dbf020ba26af3118c610f973962cfad7706435b84ae1caf5ac44cafb224b83",
		},
		{
			"enode://4ea4cdbf8c6e1141e4aec3a6ad7dec3a6988fcd31b6ceb9dddd8c23b4f9557e2b771d4c29b301f796f74da2c6f4e7415e1aa7ba3eaa9147b7ab9995751486685@3.120.24.72:16789",
			"988096ea81e13e8c7f4c080f076d9a0583a0fe11b72a5264e2c6b6bb96c2aa5e76c219f56a2cd32d24c0dd2218d7fd0369f0d184315a65c49f075a68fc3c103459401cb4570da7905e69c3bb33df1f0d22bd762bbace3ba90134406a4400b511",
		},
		{
			"enode://5c28c6443fd1bdd3ae02625cb54c0c8b739cde1c137ed50a4f03103a77a5936abe7c1327db2b9439f429568ca0aace7b8da5f47c6b2bf6ede1eb128eb0ea9201@13.126.189.133:16789",
			"7a3fa7022ba9e9b8ada7b0d273f00999c3f3c32a6db13615c9b1a57dd4abeee6814ae0c1ac2601e9839255d793c503060b7a43d9d1b144f0ab69ddd1972d2e80ed29e3e407b5bd9b7860d77a34ba693ff5393714fa1d377da3cba78d5256fc89",
		},
		{
			"enode://f2afba5553cf79eef374e900fec4c3ab3d6223a4fe91ded26872cea6d01caea1a8394f4f09faa6d24b9f3c6105e342c075b82d806b187798cb875381c88e29a2@18.140.124.104:16789",
			"	cb6490c888a683ac2a707c51c733be20621be21192fe5b80b008920ad6fc52cba26969479f567d3ecf1259819691eb0356ca89317176779c1d5002a6a306ba14c284e14ed6c8bbcb3a4b5bc8100168c6c0f7cbe3289337cf47771cb2a1098a95",
		},
	}

	initialTestnetConsensusNodes = []initNode{
		{
			"enode://a6ef31a2006f55f5039e23ccccef343e735d56699bde947cfe253d441f5f291561640a8e2bbaf8a85a8a367b939efcef6f80ae28d2bd3d0b21bdac01c3aa6f2f@test-sea.platon.network:16791", //TEST-SEA
			"",
		},
		{
			"enode://d124e660938dc3fd63d913ff753fafc262764b22294431e760b572b0b58d5e6b813b32ccbacc326c03171542ae0ff8ff6528625a2d612e0c49240f111eba3c22@test-sg.platon.network:16792", //TEST-SG
			"",
		},
		{
			"enode://24b0c456ae5cad46c4fb9bc02c867b997e22f30696e6e330926f785ca2e7410baf1eb34ffd9b5b07b5ba6e02b693faf57afb33f7c66cfbcf4c9186b4bfac737d@test-na.platon.network:16791", //TEST-NA
			"",
		},
		{
			"enode://c7fc34d6d8b3d894a35895aaf2f788ed445e03b7673f7ce820aa6fdc02908eeab6982b7eb97e983cc708bcec093b3bc512b0b1fbf668e6ab94cd91f2d642e591@test-us.platon.network:16792", //TEST-US
			"",
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
		EIP155Block: big.NewInt(3),
		Cbft: &CbftConfig{
			InitialNodes: ConvertNodeUrl(initialTestnetConsensusNodes),
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
