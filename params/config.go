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
			"enode://a84f5047c13a01b01b1f5c8be63c360d1d9d49d4d0f0c55deacade344bd3f1d12033b1334d8e42646f7eaf4a46ab325e6e02d020f4ac157ac71a4d202cd5ccea@13.69.9.245:16789", //TEST-SEA
			"6af3f69120080ea55fc803fbe5fc659a2274d413b8c6162fbd8ab6b1bb88de7e01c395977a3a64f91f04132104041c0035b13fa36251672e10820c043d2fdd09798120527682df31f57a4cd57030ecb084f84fa1c5b78d1272c162b001d20280",
		},
		{
			"enode://01fd4af4c8cdf5d7b934c0ad955d0b7b44114555457f49032d2eccf838abe18eff5a336be7d0fd93cb50359192d1a4c063eaad4b395f6c5b6b6dc3fda64ccfbb@207.46.233.122:16789", //TEST-SG
			"ac7f88010a1838b23c60c075d83f1cbf412e833187dd838d18c560e4e6f1cd17672695f31262578b580501804177e70c9c0fff88f7dcd257848f7af8e9a39921db378547e6f8e7111d1de2651261ce03d2a4a1821e7774bae25acd28d5f09a0b",
		},
		{
			"enode://a2b436b09818b678c38983c51e80de7e564a42c5f875ad0a7e6dddddbd6b3b0396ba02fa1f19542765b03c74ffb33774f67fb13615e59bb31c451a44e9678117@52.142.166.24:16789", //TEST-NA
			"4e76bff3d7487cc8ee2f47ba09c997df386d3a06e9771067023ecb1f7dbcc8123ff59af9b6fd7814cc793e4e11d4420890532ae6d63fddb5a0e56b019567a5f6181d4ad14bc3e5771a285738c8def864526246cccb7fa6ee3cbdef3083e08009",
		},
		{
			"enode://f32160f10265d2273e4a625c9795564689ac80aa40a75f7f2e93aa31f214a0459a47c6813079355504928f1f521f86a70f1a31d4db607ee128df72825f967dd8@52.143.129.153:16789", //TEST-US
			"b96ddc93537ef7ec25626d272be4f2ba445fdedb5b40fbf4b8ea519c9144df309e862e8e64f6871cdbfb6cb1b8583d1960f4781358999b283c66b096fc16b596c81392f912a9905f0b6df3eb7cdf7a044eba2ef3efd046cd039b7bc85553130b",
		},
		{
			"enode://4f4946d376233ed5a0440c0a68fbf9a0dd926e276a7c88715e889ce2029f1f8c25d5a7abb9dfcb155cf0a73888e2e2b79e07ac10f3c4a810ec57201f2fbd2389@52.63.239.155:16789", //TEST-US
			"81e4458fe917b277378b6de409244973c9d1a27d8bc0b18e7230edee30dc69e4575d2d732fa662cc43021ce8f3c20901be367e7b3596405d28181beea68be83d3a4e9a627b9ffb6cbd098859e62ac3fc778a2cd4b3e70d1eafb9edef06dc1d09",
		},
		{
			"enode://53c2582db1d5397dab9c079e9d03a3fd81b5332eda145f3c10e62df365f3fea828a880f8d39803e6555f9df84453cbf27d2bf80b27c09727f0909044da5d560c@35.182.73.204:16789", //TEST-US
			"455582d7d02884c3aa9a8d8b3173b88a76248853cc4ee1e973b3bd0a192f28834d8121f99cdba1dcee06e8437b9fc60b6d88f5e5c662d2dd24be4469be209c20a07f0b9cb210ff5a8b99d8a2b4a90d7d3ed8b60c73b8727e109be2f803abe811",
		},
		{
			"enode://61241ed97c11ffb8843b6f62263b63eb2d47b73e1ccc7ca474fd3e8f212c56d3e1687f4fe380ce122a478773ba70b5f893ee4baece836513d3b68a88a2977631@3.122.68.59:16789", //TEST-US
			"0c722bbee6314348218648cb63244364c490ca8e04e8cfbc6f38d2ec90a8b33be65c987692e84aae5fe47926f4e0a00c7e08f3a1ccc36e5576b684b7fda2b34646723ffdd8a841d4ad9e8c1b228073fba4ddf6e03620ecdd8f6c07b158103289",
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

	// TestnetChainConfig contains the chain parameters to run a node on the test network.
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
	}

	// AllEthashProtocolChanges contains every protocol change (EIPs) introduced
	//
	// This configuration is intentionally not using keyed fields to force anyone
	// adding flags to the config to also have to set these fields.
	AllEthashProtocolChanges = &ChainConfig{big.NewInt(1337), "", big.NewInt(0), big.NewInt(0), nil, nil, ""}

	TestChainConfig = &ChainConfig{big.NewInt(1), "", big.NewInt(0), big.NewInt(0), nil, new(CbftConfig), ""}
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
