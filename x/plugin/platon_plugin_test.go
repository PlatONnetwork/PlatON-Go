package plugin_test

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

//func init() {
//	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
//}

var (
	nodeIdArr = []discover.NodeID{
		discover.MustHexID("a20aef0b2c6baeaa34be2848e7dfc04c899b5985adf6fa0e98b38f754f2bb0c47974506a8de13f2a2ae97c08bcb12b438b3dcbf237b7be58f6d6d8beb36dd235"),
		discover.MustHexID("49ddf47ba0463eb4b40e66365660990ab035daf9cafd75c3c1cf7aa98f3021f157933381fe4cdc4496c40bf3fb0f732711c0ba8459367a8310a9a7ffeb33cb71"),
		discover.MustHexID("8332fd21a9da740ad795ce50a322075beec16333daea6e5cf8049fd9b1cdd0277e16e7ba3bd29a9edb6a357349b6cdb01892095ab94d2d2a4f84a31717df68d7"),
		discover.MustHexID("02e5274c02b0142b6e59d74c68729afc307f32ebf6ad3ac15fb2c69c537928869a1b02a054421d3160a48f724250818e7dbb7ced3408b9f7a073eca62674500f"),
		discover.MustHexID("becb9e9c467b742f41a92dd6629cb107c2159421a1241b1176062c7f652ca5e090401aa379a016efea2e65691da17be51409c04f904ef3fa461830e6cc1a85ce"),
		discover.MustHexID("0beed0b7bdfd674e1b51f3dc1d5cfbafba3def1ecdcdcef2640d3fb842aef3bb702490729096611bd74989297633416321578ddd5514aab90a359f6428a3cc2e"),
		discover.MustHexID("f79d932933e952758e568ec5d21c71c7eb3fc7554370e6bd11cac9f1eda6a99cff64a48991df48bb5a94198c65745767b6c4368dd183be9962ec7a17bb597fbe"),
		discover.MustHexID("5fa08d295fd276da2d351f0ff96d415cbe85eb2704df00ea7f83405038164bef3dc1c3f160f8149524379412f5a18e1f61b99ab51e213f075afe3f1ab95b8ba5"),
		discover.MustHexID("5a4ad751151b7a8cce27cb60b780ae04bfb380c51132df7a9795c2626ad026912111eec579e49d594ebc560db56f12d8d1239dd52f6d23b6039e16bef18a549c"),
		discover.MustHexID("79df7a746bc4c298cb53848435927004bd3e48f301b32d7b68b0b79f414ffdedd30eb6532d64b56393f59fcc74acb4d730408d017697d3d76f9e146e250ea4e3"),
		discover.MustHexID("5fc56f8faf8eb42bcaff11184ba73c81dbcf503bf07772b794b6e7c13523f401d5bc5efd877cd244afca3dc8d335c9f3062f777bd1e631f077d55bc5153dc665"),
		discover.MustHexID("7f8f1828a89acae44ba101ff07f78ebaac40e2cbf6609256482c75097475b273348e0a2acb4bcfd1d2cba4f2bdf72640a327faf50cef1b86ba1cad7baf6fe818"),
		discover.MustHexID("fe104c8152f823c0e9b105bca203ff6f0f09a7feae2c4949afe8b00cbea67ff39bd48aa6b232e53baa1312fe1051408364e3dfba6438afa0ef7a27fd0a068e0d"),
		discover.MustHexID("6c7cdedda488a335b8b0851decd09e512763ca4411021c7171a4cf5d535f3084057ac6b0e98876bce5dbaec9a966e8f7555828fca2a73385eb446f3fed287582"),
		discover.MustHexID("cd40af30fe8da249a4d02dc51a89cc6f1008d281f8364c5777eda1168084d26f4e2553d56b6246ba70fe5a331d111c5d62bc833fe527cb82c2052a8c003eef24"),
		discover.MustHexID("bef1f8017ee14c58a8bcc9288a28acb38cf52bbbd1fbe7d7eea9b52c11456ea50ade33531ea45a68b8aad71c2408fa19e7ba61b0c0ac2a2bfb6c5864b54392ab"),
		discover.MustHexID("58430dacdba6823deee7af5c4d24df65e9678fc1b0b85dce49e59ccadb5bc220425563d86baaa93194023cb99e6dfeae58229966e9058d83ccba087fa4edd636"),
		discover.MustHexID("6eaf01e5fd1b953a71a60fcd79317e1f69c212e5673d5670a8980b5a3729a4b8e7c6d6ca16719d81e4d77f0db399ef1163601179b0ac1d88b8c998600a407a9d"),
		discover.MustHexID("4f1f036e5e18cc812347d5073cbec2a8da7930de323063c39b0d4413a396e088bfa90e8c28174313d8d82e9a14bc0884b13a48fc28e619e44c48a49b4fd9f107"),
		discover.MustHexID("f18c596232d637409c6295abb1e720db99ffc12363a1eb8123d6f54af80423a5edd06f91115115a1dca1377e97b9031e2ddb864d34d9b3491d6fa07e8d9b951b"),
		discover.MustHexID("7a8f7a28ac1c4eaf98b2be890f372e5abc58ebe6d3aab47aedcb0076e34eb42882e926676ebab327a4ef4e2ea5c4296e9c7bc0991360cb44f52672631012db1b"),
		discover.MustHexID("9eeb448babf9e93449e831b91f98d9cbc0c2324fe8c43baac69d090717454f3f930713084713fe3a9f01e4ca59b80a0f2b41dbd6d531f414650bab0363e3691a"),
		discover.MustHexID("cc1d7314c15e30dc5587f675eb5f803b1a2d88bfe76cec591cec1ff678bc6abce98f40054325bdcb44fb83174f27d38a54fbce4846af8f027b333868bc5144a4"),
		discover.MustHexID("e4d99694be2fc8a53d8c2446f947aec1c7de3ee26f7cd43f4f6f77371f56f11156218dec32b51ddce470e97127624d330bb7a3237ba5f0d87d2d3166faf1035e"),
		discover.MustHexID("9c61f59f70296b6d494e7230888e58f19b13c5c6c85562e57e1fe02d0ff872b4957238c73559d017c8770b999891056aa6329dbf628bc19028d8f4d35ec35823"),
	}

	addrArr = []common.Address{
		common.HexToAddress("0xc9E1C2B330Cf7e759F2493c5C754b34d98B07f93"),
		common.HexToAddress("0xd87E10F8efd2C32f5e88b7C279953aEF6EE58902"),
		common.HexToAddress("0xeAEc60C738eeD9468e6AcCc1d403faCF1A670F6D"),
		common.HexToAddress("0x5c5994165265Ac31AAFE874a231f2C5d0eF29C3a"),
		common.HexToAddress("0xB9449Eb226cb93c3BF5FeCA16c85a737538e24f0"),
		common.HexToAddress("0x908bad1823BddA66cc65E788b9d0194b7975976A"),
		common.HexToAddress("0x3DfC64A87db521662675DffEa48d0c208414D4f8"),
		common.HexToAddress("0xad8adf35068Cdf572c9eFb5a069dA48D2E165Aa1"),
		common.HexToAddress("0xf33b5Da47c6ECbC61cF07C7387Afc6ef0EA2f866"),
		common.HexToAddress("0x2E5FB4F78E3FB9b1898DE7d7D8dB3d44C62040be"),
		common.HexToAddress("0x285CF84ea3E177E1fC9F396aEbc9329a08f51bb5"),
		common.HexToAddress("0x91BffdC88329AfDD97DF6fe92cfd4FcB7927Aecd"),
		common.HexToAddress("0x58b62FfF5046aF2252F1F8Ecb5c3342ada394F72"),
		common.HexToAddress("0x8ec116c11d8515e8222Cabc4BEc06A880C51D929"),
		common.HexToAddress("0x364eCBade4c35beE2F8a04F8209BaB236B48A35a"),
		common.HexToAddress("0x26896c394A1E12095e822e5b080e8EfA050c738C"),
		common.HexToAddress("0x314253824CD6b7BCF1613CAB00126D6076F7a389"),
		common.HexToAddress("0x5544F05D51E45fa6497AFEC0F1A5d64531B21be0"),
		common.HexToAddress("0x3da830FAd2A6983d948d7262B2AdF7eA53b953be"),
		common.HexToAddress("0x815A7910C035F2FB9451cDA349969788449c2288"),
		common.HexToAddress("0x4Cdd49e08587c824c7629e7d124390B70d105740"),
		common.HexToAddress("0xD041b5fAaa4B721241A55107FE8F19ce1ba3E7fD"),
		common.HexToAddress("0xcbc583DEdbbE6b51B86036C040596bB0a0299a73"),
		common.HexToAddress("0x1c0A4509Ba46deA47775Ad8B20A19f398B820642"),
		common.HexToAddress("0xEEE10Fc4A3AB339f5a788f3b82Eb57738F075EcE"),
	}

	priKeyArr = []*ecdsa.PrivateKey{
		crypto.HexMustToECDSA("d30b490011d2a08053d46506ae533ff96f2cf6a37f73be740f52ad24243c4958"),
		crypto.HexMustToECDSA("3dbace449229be40a641e056c250fee39a2b6077f3f40a47512e3559cb7b6174"),
		crypto.HexMustToECDSA("a31f894228699550e0d53142d1b50c01d991d0e195d7f6e98d720fb4b14cbb84"),
		crypto.HexMustToECDSA("7a6a73919df68e4c36c2b38b40fe6d154c8453d1c439262bd3912a0537ce4e51"),
		crypto.HexMustToECDSA("351f8010623c4b7b1c358ce4dce9e1d1a6df4df947c0426d8dd50efdb9e5e2bf"),
		crypto.HexMustToECDSA("17014f74092b63bd72f55d56f9b511532de7607620acbdfe9a55c3359e0a05d8"),
		crypto.HexMustToECDSA("62dcf763f4bd71a22bfacdee247aa827858303f9e2f48d06be287c10b76c392e"),
		crypto.HexMustToECDSA("b4b61440fe55277f47bcdefe93f849e10371bceabc26aa7fdcd049a80b8eaf98"),
		crypto.HexMustToECDSA("d704124f3623ae25465fcad35405f632abb01d592f57bb201d729a8955c31647"),
		crypto.HexMustToECDSA("82e9691ff8853c656e743515e606e38617edd7ecbcabb101a6796e7c53c107e5"),
		crypto.HexMustToECDSA("7eab72a552af80835c994f265e9fbef3d043d779e943efeaff29e228199bbf6a"),
		crypto.HexMustToECDSA("8df3b87b8b25c7e17dc6b46146b6581c965d11cd8d50db43c60b9a7566bc0f18"),
		crypto.HexMustToECDSA("ada826fbded7ee2a06d3cec6e5d1567359f9af0ffcb52cf1701c686e6884e01a"),
		crypto.HexMustToECDSA("5f97915d61ecdb38a048d8fa5a7bf7e16a46f3ecb92d8d7e48506cf1c6705ea4"),
		crypto.HexMustToECDSA("625308f7e5a3de990f7297e20d0ea10a6b36a2d374dc1a68b051053cfd739313"),
		crypto.HexMustToECDSA("207b68f34e8432163013003934da51889c3ac8c00fa5fa6f80299678882d2f50"),
		crypto.HexMustToECDSA("aba5c7e329207d7cdccc1d3e1657d935d2f6029372bbafdb953c6d0b3a81d73b"),
		crypto.HexMustToECDSA("77576a78f1fe8018bb0ae27076c86980c56e62c6014317796a2ca3ce59b1e7c0"),
		crypto.HexMustToECDSA("548ceef29a39093e48ef65bc98b210320dedd79ca40acebeb573f8eb72018aac"),
		crypto.HexMustToECDSA("73a2bd8694f883ff5f11551c04303ff7180ae6ef1b89170a67ace10d04c7c3e2"),
		crypto.HexMustToECDSA("996e2bb9c1371e50125fb8b1d0e6f9c46148dfb8b01d9edd6e8b5ec1a6241316"),
		crypto.HexMustToECDSA("51c977a01d5517406fcce2bf7bbb44c67e6b876641a5dac6d2fc26b2f6a97001"),
		crypto.HexMustToECDSA("41d4ce3f8b18fc7ccb4bb0e9514e0863d0c0bd4bb26e9fba3c2a384189c2000b"),
		crypto.HexMustToECDSA("3653b25ba39e59d12a3f45f0fb324b8588db839de4bafd9b938315c356a37051"),
		crypto.HexMustToECDSA("e066f9c4daabcc354162165f8aa161c0bc1cede1b0d14a269f63f6d6bdb1ec5d"),
	}

	blockNumber = big.NewInt(1)
	blockHash   = common.HexToHash("9d4fb5346abcf593ad80a0d3d5a371b22c962418ad34189d5b1b39065668d663")

	blockNumber2 = big.NewInt(2)
	blockHash2   = common.HexToHash("c95876b92443d652d7eb7d7a9c0e2c58a95e934c0c1197978c5445180cc60980")

	blockNumber3 = big.NewInt(3)
	blockHash3   = common.HexToHash("3b198bfd5d2907285af009e9ae84a0ecd63677110d89d7e030251acb87f6487e")

	lastBlockNumber uint64
	lastBlockHash   common.Hash
	lastHeader      types.Header

	sender        = common.HexToAddress("0xeef233120ce31b3fac20dac379db243021a5234")
	anotherSender = common.HexToAddress("0xeef233120ce31b3fac20dac379db243021a5235")
	sndb          = snapshotdb.Instance()

	// serial use only
	sender_balance, _ = new(big.Int).SetString("9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)

	txHashArr = []common.Hash{
		common.HexToHash("0x00000000000000000000000000000000000000886d5ba2d3dfb2e2f6a1814f22"),
		common.HexToHash("0x000000000000000000000000000000005249b59609286f2fa91a2abc8555e887"),
		common.HexToHash("0x000000008dba388834e2515c4d9ccb02a48bae177e73959330e55067211c2456"),
		common.HexToHash("0x0000000000000000000000000000000000009a715a765a72b8a289156f9543c9"),
		common.HexToHash("0x0000e1b4a5508c11772b61f463657585c33b577019e4a23bd359c018a4e306d1"),
		common.HexToHash("0x00fd854f940e2d2af8e74c33e640ea6f75c1d9ee49b816b8a4647611d0c91863"),
		common.HexToHash("0x0000000000001038575739a53385cfe42321585a56050e18f8ea2b3e8dc21966"),
		common.HexToHash("0x0000000000000000000000000000000000000048f3b312dc8d081e1186abe8c2"),
		common.HexToHash("0x000000000000000000000000f5bd37579e7ca954eba8fbe7a65646250e92ab7d"),
		common.HexToHash("0x00000000000000000000000000000000000000001d65a5a69fed6ddb0cb58dff"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000d2"),
		common.HexToHash("0x0000000000000000000000000000000000000000000000000000f2e8b2706c9e"),
		common.HexToHash("0x00000000000000000000000000e22a393898aac376b079e0894e8e2be6024d03"),
		common.HexToHash("0x000000000000000000000000000000000000000000000000483570dd0679860a"),
		common.HexToHash("0x000000000000000000000000000000000000007fc9e1dc435b5d0064ac50fd4e"),
		common.HexToHash("0x00000000000000000000000000cbeb8f4d51969d7eb70a4f6e8505950d870df7"),
		common.HexToHash("0x00000000000000000000000000000000000000000000000000000000000000b4"),
		common.HexToHash("0x000000008fd2abdf28d87efb2c7fa2d37618c8dba97059376d6a58007bee3d8b"),
		common.HexToHash("0x0000000000000000000000003566f3a0adf49d90e610ef3d3548b5a72b1fe199"),
		common.HexToHash("0x00000000000054fa3d19eb57e98aa1dd69d216722054d8539ede4b89c5b77ee9"),
	}

	initProgramVersion = uint32(1<<16 | 0<<8 | 0) // 65536, version: 1.0.0
	promoteVersion     = uint32(2<<16 | 0<<8 | 0) // 131072, version: 2.0.0

	balanceStr = []string{
		"9000000000000000000000000",
		"60000000000000000000000000",
		"1300000000000000000000000",
		"1100000000000000000000000",
		"1000000000000000000000000",
		"4879000000000000000000000",
		"1800000000000000000000000",
		"1000000000000000000000000",
		"1000000000000000000000000",
		"70000000000000000000000000",
		"5550000000000000000000000",
		"9000000000000000000000000",
		"60000000000000000000000000",
		"1300000000000000000000000",
		"1100000000000000000000000",
		"1000000000000000000000000",
		"4879000000000000000000000",
		"1800000000000000000000000",
		"1000000000000000000000000",
		"1000000000000000000000000",
		"70000000000000000000000000",
		"5550000000000000000000000",
		"1000000000000000000000000",
		"70000000000000000000000000",
		"5550000000000000000000000",
	}

	nodeNameArr = []string{
		"PlatON",
		"Gavin",
		"Emma",
		"Kally",
		"Juzhen",
		"Baidu",
		"Alibaba",
		"Tencent",
		"ming",
		"hong",
		"gang",
		"guang",
		"hua",
		"PlatON_2",
		"Gavin_2",
		"Emma_2",
		"Kally_2",
		"Juzhen_2",
		"Baidu_2",
		"Alibaba_2",
		"Tencent_2",
		"ming_2",
		"hong_2",
		"gang_2",
		"guang_2",
	}

	chaList = []string{"A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g", "H", "h", "J", "j", "K", "k", "M", "m",
		"N", "n", "P", "p", "Q", "q", "R", "r", "S", "s", "T", "t", "U", "u", "V", "v", "W", "w", "X", "x", "Y", "y", "Z", "z"}

	specialCharList = []string{
		"☄", "★", "☎", "☻", "♨", "✠", "❝", "♚", "♘", "✎", "♞", "✩", "✪", "❦", "❥", "❣", "웃", "卍", "Ⓞ", "▶", "◙", "⊕", "◌", "⅓", "∭",
		"∮", "╳", "㏒", "㏕", "‱", "㎏", "❶", "Ň", "🅱", "🅾", "𝖋", "𝕻", "𝕼", "𝕽", "お", "な", "ぬ", "㊎", "㊞", "㊮", "✘"}
)

func TestVersion(t *testing.T) {

	t.Log("the version is:", promoteVersion)
}

func newEvm(blockNumber *big.Int, blockHash common.Hash, state *state.StateDB) *vm.EVM {
	if nil == state {
		state, _, _ = newChainState()
	}
	evm := &vm.EVM{
		StateDB: state,
	}
	context := vm.Context{
		BlockNumber: blockNumber,
		BlockHash:   blockHash,
	}
	evm.Context = context

	//set a default active version
	govDB := gov.GovDBInstance()
	govDB.AddActiveVersion(initProgramVersion, 0, state)

	return evm
}

func newPlugins() {
	plugin.GovPluginInstance()
	plugin.StakingInstance()
	plugin.SlashInstance()
	plugin.RestrictingInstance()
	plugin.RewardMgrInstance()

	snapshotdb.Instance()
}

func newChainState() (*state.StateDB, *types.Block, error) {

	url := "enode://0x7bae841405067598bf65e7260ca693a964316e752249c4970085c805dbee738fdb41fc434e96e2b65e8bf1db2f52f05d9300d04c1e6129c26cb5d0f214b49968@platon.network:16791"
	node, _ := discover.ParseNode(url)
	gen := core.DefaultGenesisBlock()
	gen.Config.Cbft.InitialNodes = []discover.Node{*node}
	gen.Config.Cbft.ValidatorMode = "ppos"

	var (
		db      = ethdb.NewMemDatabase()
		genesis = gen.MustCommit(db)
	)

	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, nil, vm.Config{}, nil)

	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		return nil, nil, errors.New("reference statedb failed" + err.Error())
	} else {
		state = statedb
	}
	state.AddBalance(sender, sender_balance)
	for i, addr := range addrArr {

		amount, _ := new(big.Int).SetString(balanceStr[len(addrArr)-1-i], 10)
		amount = new(big.Int).Mul(common.Big257, amount)
		state.AddBalance(addr, amount)
	}

	return state, genesis, nil
}

func build_staking_data_more(block uint64) {

	no := int64(block)
	header := types.Header{
		Number: big.NewInt(no),
	}
	hash := header.Hash()

	stakingDB := staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(int64(block)), lastBlockHash, hash)
	// MOCK

	validatorArr := make(staking.ValidatorQueue, 0)

	// build  more data
	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		rand.Seed(time.Now().UnixNano())

		weight := rand.Intn(1000000000)

		ii := rand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		randBuildFunc := func() (discover.NodeID, common.Address, error) {
			privateKey, err := crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random NodeId private key: %v", err)
				return discover.NodeID{}, common.ZeroAddr, err
			}

			nodeId := discover.PubkeyID(&privateKey.PublicKey)

			privateKey, err = crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random Address private key: %v", err)
				return discover.NodeID{}, common.ZeroAddr, err
			}

			addr := crypto.PubkeyToAddress(privateKey.PublicKey)

			return nodeId, addr, nil
		}

		var nodeId discover.NodeID
		var addr common.Address

		if i < 25 {
			nodeId = nodeIdArr[i]
			ar, _ := xutil.NodeId2Addr(nodeId)
			addr = ar
		} else {
			id, ar, err := randBuildFunc()
			if nil != err {
				return
			}
			nodeId = id
			addr = ar
		}

		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(i + 1),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		stakingDB.SetCanPowerStore(blockHash, canAddr, canTmp)
		stakingDB.SetCandidateStore(blockHash, canAddr, canTmp)

		v := &staking.Validator{
			NodeAddress: canAddr,
			NodeId:      canTmp.NodeId,
			StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
				fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
			ValidatorTerm: 0,
		}
		validatorArr = append(validatorArr, v)
	}

	queue := validatorArr[:25]

	epoch_Arr := &staking.Validator_array{
		//Start: ((block-1)/22000)*22000 + 1,
		//End:   ((block-1)/22000)*22000 + 22000,
		Start: ((block-1)/uint64(xutil.CalcBlocksEachEpoch()))*uint64(xutil.CalcBlocksEachEpoch()) + 1,
		End:   ((block-1)/uint64(xutil.CalcBlocksEachEpoch()))*uint64(xutil.CalcBlocksEachEpoch()) + uint64(xutil.CalcBlocksEachEpoch()),
		Arr:   queue,
	}

	pre_Arr := &staking.Validator_array{
		Start: 0,
		End:   0,
		Arr:   queue,
	}

	curr_Arr := &staking.Validator_array{
		//Start: ((block-1)/250)*250 + 1,
		//End:   ((block-1)/250)*250 + 250,
		Start: ((block-1)/uint64(xutil.ConsensusSize()))*uint64(xutil.ConsensusSize()) + 1,
		End:   ((block-1)/uint64(xutil.ConsensusSize()))*uint64(xutil.ConsensusSize()) + uint64(xutil.ConsensusSize()),
		Arr:   queue,
	}

	setVerifierList(hash, epoch_Arr)
	setRoundValList(hash, pre_Arr)
	setRoundValList(hash, curr_Arr)

	lastBlockHash = hash
	lastBlockNumber = block
	lastHeader = header
}

func build_staking_data(genesisHash common.Hash) {
	stakingDB := staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(1), genesisHash, blockHash)
	// MOCK

	validatorArr := make(staking.ValidatorQueue, 0)

	count := 0
	// build  more data
	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		rand.Seed(time.Now().UnixNano())

		weight := rand.Intn(1000000000)

		ii := rand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		randBuildFunc := func() (discover.NodeID, common.Address, error) {
			privateKey, err := crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random NodeId private key: %v", err)
				return discover.NodeID{}, common.ZeroAddr, err
			}

			nodeId := discover.PubkeyID(&privateKey.PublicKey)

			privateKey, err = crypto.GenerateKey()
			if nil != err {
				fmt.Printf("Failed to generate random Address private key: %v", err)
				return discover.NodeID{}, common.ZeroAddr, err
			}

			addr := crypto.PubkeyToAddress(privateKey.PublicKey)

			return nodeId, addr, nil
		}

		var nodeId discover.NodeID
		var addr common.Address

		if i < 25 {
			nodeId = nodeIdArr[i]
			ar, _ := xutil.NodeId2Addr(nodeId)
			addr = ar
		} else {
			id, ar, err := randBuildFunc()
			if nil != err {
				return
			}
			nodeId = id
			addr = ar
		}

		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(i + 1),
			Shares:          balance,
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),
			// Prevent null pointer initialization
			Released:           common.Big256,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		err := stakingDB.SetCanPowerStore(blockHash, canAddr, canTmp)
		if nil != err {
			fmt.Printf("Failed to SetCanPowerStore: %v", err)
			return
		}
		err = stakingDB.SetCandidateStore(blockHash, canAddr, canTmp)
		if nil != err {
			fmt.Printf("Failed to SetCandidateStore: %v", err)
			return
		}

		v := &staking.Validator{
			NodeAddress: canAddr,
			NodeId:      canTmp.NodeId,
			StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), canTmp.Shares.String(),
				fmt.Sprint(canTmp.StakingBlockNum), fmt.Sprint(canTmp.StakingTxIndex)},
			ValidatorTerm: 0,
		}
		validatorArr = append(validatorArr, v)
		count++
	}

	fmt.Printf("build staking  data count: %d \n", count)
	queue := validatorArr[:25]

	epoch_Arr := &staking.Validator_array{
		Start: 1,
		End:   uint64(xutil.CalcBlocksEachEpoch()),
		Arr:   queue,
	}

	pre_Arr := &staking.Validator_array{
		Start: 0,
		End:   0,
		Arr:   queue,
	}

	curr_Arr := &staking.Validator_array{
		Start: 1,
		End:   uint64(xutil.ConsensusSize()),
		Arr:   queue,
	}

	setVerifierList(blockHash, epoch_Arr)
	setRoundValList(blockHash, pre_Arr)
	setRoundValList(blockHash, curr_Arr)

	lastBlockHash = blockHash
	lastBlockNumber = blockNumber.Uint64()
	lastHeader = types.Header{
		Number: blockNumber,
	}

}

func buildBlockNoCommit(blockNum int) {

	no := int64(blockNum)
	header := types.Header{
		Number: big.NewInt(no),
	}
	hash := header.Hash()

	staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(int64(blockNum)), lastBlockHash, hash)

	lastBlockHash = hash
	lastBlockNumber = uint64(blockNum)
	lastHeader = header
}

func build_gov_data(state *state.StateDB) {

	//set a default active version
	govDB := gov.GovDBInstance()
	govDB.AddActiveVersion(initProgramVersion, 0, state)
}

func buildStateDB(t *testing.T) xcom.StateDB {
	db := ethdb.NewMemDatabase()
	stateDb, err := state.New(common.Hash{}, state.NewDatabase(db))

	if err != nil {
		t.Errorf("new state db failed: %s", err.Error())
	}

	return stateDb
}

func buildDbRestrictingPlan(account common.Address, t *testing.T, stateDB xcom.StateDB) {

	const Epochs = 5
	var list = make([]uint64, 0)

	for epoch := 1; epoch <= Epochs; epoch++ {
		// build release account record
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(epoch), 1)
		stateDB.SetState(cvm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

		// build release amount record
		releaseAmount := big.NewInt(int64(1E18))
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(epoch), account)
		stateDB.SetState(cvm.RestrictingContractAddr, releaseAmountKey, releaseAmount.Bytes())

		// build release epoch record
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(epoch))
		stateDB.SetState(cvm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		list = append(list, uint64(epoch))
	}

	// build restricting user info
	var user restricting.RestrictingInfo
	user.Balance = big.NewInt(int64(5E18))
	user.Debt = big.NewInt(0)
	user.DebtSymbol = false
	user.ReleaseList = list

	bUser, err := rlp.EncodeToBytes(user)
	if err != nil {
		t.Fatalf("failed to rlp encode restricting info: %s", err.Error())
	}

	// build restricting account info record
	restrictingKey := restricting.GetRestrictingKey(account)
	stateDB.SetState(cvm.RestrictingContractAddr, restrictingKey, bUser)

	stateDB.AddBalance(sender, sender_balance)
	stateDB.AddBalance(cvm.RestrictingContractAddr, big.NewInt(int64(5E18)))
}

func buildDBStakingRestrictingFunds(t *testing.T, stateDB xcom.StateDB) {

	account := addrArr[0]

	// build release account record
	releaseAccountKey := restricting.GetReleaseAccountKey(1, 1)
	stateDB.SetState(cvm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

	// build release amount record
	releaseAmount := big.NewInt(int64(2E18))
	releaseAmountKey := restricting.GetReleaseAmountKey(1, account)
	stateDB.SetState(cvm.RestrictingContractAddr, releaseAmountKey, releaseAmount.Bytes())

	// build release epoch record
	releaseEpochKey := restricting.GetReleaseEpochKey(1)
	stateDB.SetState(cvm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

	var releaseEpochList = []uint64{1}

	// build restricting user info
	var user restricting.RestrictingInfo
	user.Balance = big.NewInt(int64(1E18))
	user.Debt = big.NewInt(0)
	user.DebtSymbol = false
	user.ReleaseList = releaseEpochList

	bUser, err := rlp.EncodeToBytes(user)
	if err != nil {
		t.Fatalf("failed to rlp encode restricting info: %s", err.Error())
	}

	// build restricting account info record
	restrictingKey := restricting.GetRestrictingKey(account)
	stateDB.SetState(cvm.RestrictingContractAddr, restrictingKey, bUser)

	stateDB.AddBalance(cvm.RestrictingContractAddr, big.NewInt(int64(1E18)))
}

func setRoundValList(blockHash common.Hash, val_Arr *staking.Validator_array) error {

	stakeDB := staking.NewStakingDB()

	queue, err := stakeDB.GetRoundValIndexByBlockHash(blockHash)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to setRoundValList: Query round valIndex is failed", "blockHash",
			blockHash.Hex(), "Start", val_Arr.Start, "End", val_Arr.End, "err", err)
		return err
	}

	var indexQueue staking.ValArrIndexQueue

	index := &staking.ValArrIndex{
		Start: val_Arr.Start,
		End:   val_Arr.End,
	}

	if len(queue) == 0 {
		indexQueue = make(staking.ValArrIndexQueue, 0)
		_, indexQueue = indexQueue.ConstantAppend(index, plugin.RoundValIndexSize)
	} else {

		has := false
		for _, indexInfo := range queue {
			if indexInfo.Start == val_Arr.Start && indexInfo.End == val_Arr.End {
				has = true
				break
			}
		}
		indexQueue = queue
		if !has {

			shabby, queue := queue.ConstantAppend(index, plugin.RoundValIndexSize)
			indexQueue = queue
			// delete the shabby validators
			if nil != shabby {
				if err := stakeDB.DelRoundValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
					log.Error("Failed to setRoundValList: delete shabby validators is failed",
						"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex())
					return err
				}
			}
		}
	}

	// Store new index Arr
	if err := stakeDB.SetRoundValIndex(blockHash, indexQueue); nil != err {
		log.Error("Failed to setRoundValList: store round validators new indexArr is failed", "blockHash", blockHash.Hex())
		return err
	}

	// Store new round validator Item
	if err := stakeDB.SetRoundValList(blockHash, index.Start, index.End, val_Arr.Arr); nil != err {
		log.Error("Failed to setRoundValList: store new round validators is failed", "blockHash", blockHash.Hex())
		return err
	}

	return nil
}

func setVerifierList(blockHash common.Hash, val_Arr *staking.Validator_array) error {

	stakeDB := staking.NewStakingDB()

	queue, err := stakeDB.GetEpochValIndexByBlockHash(blockHash)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to setVerifierList: Query epoch valIndex is failed", "blockHash",
			blockHash.Hex(), "Start", val_Arr.Start, "End", val_Arr.End, "err", err)
		return err
	}

	var indexQueue staking.ValArrIndexQueue

	index := &staking.ValArrIndex{
		Start: val_Arr.Start,
		End:   val_Arr.End,
	}

	if len(queue) == 0 {
		indexQueue = make(staking.ValArrIndexQueue, 0)
		_, indexQueue = indexQueue.ConstantAppend(index, plugin.EpochValIndexSize)
	} else {

		has := false
		for _, indexInfo := range queue {
			if indexInfo.Start == val_Arr.Start && indexInfo.End == val_Arr.End {
				has = true
				break
			}
		}
		indexQueue = queue
		if !has {

			shabby, queue := queue.ConstantAppend(index, plugin.EpochValIndexSize)
			indexQueue = queue
			// delete the shabby validators
			if nil != shabby {
				if err := stakeDB.DelEpochValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
					log.Error("Failed to setVerifierList: delete shabby validators is failed",
						"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex())
					return err
				}
			}
		}
	}

	// Store new index Arr
	if err := stakeDB.SetEpochValIndex(blockHash, indexQueue); nil != err {
		log.Error("Failed to setVerifierList: store epoch validators new indexArr is failed", "blockHash", blockHash.Hex())
		return err
	}

	// Store new epoch validator Item
	if err := stakeDB.SetEpochValList(blockHash, index.Start, index.End, val_Arr.Arr); nil != err {
		log.Error("Failed to setVerifierList: store new epoch validators is failed", "blockHash", blockHash.Hex())
		return err
	}

	return nil
}
