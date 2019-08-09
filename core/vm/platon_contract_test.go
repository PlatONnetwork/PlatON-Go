package vm_test

import (
	"errors"
	"fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
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

const initGas = 10000000

var (
	nodeIdArr = []discover.NodeID{
		discover.MustHexID("1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429"),
		discover.MustHexID("a6ef31a2006f55f5039e23ccccef343e735d56699bde947cfe253d441f5f291561640a8e2bbaf8a85a8a367b939efcef6f80ae28d2bd3d0b21bdac01c3aa6f2f"),
		discover.MustHexID("c7fc34d6d8b3d894a35895aaf2f788ed445e03b7673f7ce820aa6fdc02908eeab6982b7eb97e983cc708bcec093b3bc512b0b1fbf668e6ab94cd91f2d642e591"),
		discover.MustHexID("97e424be5e58bfd4533303f8f515211599fd4ffe208646f7bfdf27885e50b6dd85d957587180988e76ae77b4b6563820a27b16885419e5ba6f575f19f6cb36b0"),
		discover.MustHexID("3b53564afbc3aef1f6e0678171811f65a7caa27a927ddd036a46f817d075ef0a5198cd7f480829b53fe62bdb063bc6a17f800d2eebf7481b091225aabac2428d"),
		discover.MustHexID("858d6f6ae871e291d3b7b2b91f7369f46deb6334e9dacb66fa8ba6746ee1f025bd4c090b17d17e0d9d5c19fdf81eb8bde3d40a383c9eecbe7ebda9ca95a3fb94"),
		discover.MustHexID("e4556b211eb6712ab94d743990d995c0d3cd15e9d78ec0096bba24c48d34f9f79a52ca1f835cec589c5e7daff30620871ba37d6f5f722678af4b2554a24dd75c"),
		discover.MustHexID("114e48f21d4d83ec9ac39a62062a804a0566742d80b191de5ba23a4dc25f7beda0e78dd169352a7ad3b11584d06a01a09ce047ad88de9bdcb63885e81de00a4d"),
		discover.MustHexID("64ba18ce01172da6a95b0d5b0a93aee727d77e5b2f04255a532a9566edaee7808383812a860acf5e43efeca3d9321547bfcdefd89e9d0c605dcdb65ce0bbb617"),
		discover.MustHexID("d31b3a7714610bd8e03b2c74aca4be16de7fcc319a1e577d50e5e8796680221b4b679bf1c37966d1a158902b8686f3ca2f41a89a7176e538141082540c4f6d66"),
		discover.MustHexID("805b617b9d321a65d8936e758b5c60cd6e8c873b9f1e7c793ad5f887d26ce9667d0db2fe55a9aeb1cc81f9cf9a1e7c54473203473e3ebda89e63c03cbcfe5347"),
		discover.MustHexID("fa147bc3625acc846a9f0e1e89172ca7470baa0f86516994f70860c6fb904ddbb1849e3cf2b40c58255e38401f40d2c3e4a3bd5c2f2849b98465a5bdb80ed6a0"),
		discover.MustHexID("d8c4b58ae052ea9480577264bc1b2c09619757015849a4c92b71a4e4c8b5ede94f35d24107b1181d0711013ed7fdc068f21e6e6084b3e96750a571669715c0b1"),
	}
	addrArr = []common.Address{
		common.HexToAddress("740ce31b3fac20dac379db243021a51e80aadd24"),
		common.HexToAddress("740ce31b3fac20dac379db243021a51e80444555"),
		common.HexToAddress("740ce31b3fac20dac379db243021a51e80eeda12"),
		common.HexToAddress("740ce31b3fac20dac379db243021a51e80522345"),
		common.HexToAddress("ef678007d18427e6022059dbc264f27507cd1ffc"),
		common.HexToAddress("5a5c4368e2692746b286cee36ab0710af3efa6cf"),
		common.HexToAddress("493301712671ada506ba6ca7891f436d29185821"),
		common.HexToAddress("60ceca9c1290ee56b98d4e160ef0453f7c40d219"),
		common.HexToAddress("43355c787c50b647c425f594b441d4bd751951c1"),
		common.HexToAddress("7ef5a6135f1fd6a02593eedc869c6d41d934aef8"),
		common.HexToAddress("f466859ead1932d743d622cb74fc058882e8648a"),
		common.HexToAddress("289d485d9771714cce91d3393d764e1311907acc"),
		common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87"),
	}

	blockNumber = big.NewInt(1)
	blockHash   = common.HexToHash("9d4fb5346abcf593ad80a0d3d5a371b22c962418ad34189d5b1b39065668d663")

	blockNumber2 = big.NewInt(2)
	blockHash2   = common.HexToHash("c95876b92443d652d7eb7d7a9c0e2c58a95e934c0c1197978c5445180cc60980")

	sender = common.HexToAddress("0xeef233120ce31b3fac20dac379db243021a5234")

	sndb = snapshotdb.Instance()

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

	initProgramVersion = uint32(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch)

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
		"44488850000000000000000000000",
		"650073899000000000000000000",
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
	}

	chaList = []string{"A", "a", "B", "b", "C", "c", "D", "d", "E", "e", "F", "f", "G", "g", "H", "h", "J", "j", "K", "k", "M", "m",
		"N", "n", "P", "p", "Q", "q", "R", "r", "S", "s", "T", "t", "U", "u", "V", "v", "W", "w", "X", "x", "Y", "y", "Z", "z"}

	specialCharList = []string{
		"☄", "★", "☎", "☻", "♨", "✠", "❝", "♚", "♘", "✎", "♞", "✩", "✪", "❦", "❥", "❣", "웃", "卍", "Ⓞ", "▶", "◙", "⊕", "◌", "⅓", "∭",
		"∮", "╳", "㏒", "㏕", "‱", "㎏", "❶", "Ň", "🅱", "🅾", "𝖋", "𝕻", "𝕼", "𝕽", "お", "な", "ぬ", "㊎", "㊞", "㊮", "✘"}
)

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

	gen := &core.Genesis{
		Config: &params.ChainConfig{
			Cbft: &params.CbftConfig{
				InitialNodes:  []discover.Node{*node},
				ValidatorMode: "ppos",
			},
		},
	}

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
	govDB.SetActiveVersion(initProgramVersion, state)

	return evm
}

func newContract(value *big.Int) *vm.Contract {
	callerAddress := vm.AccountRef(sender)
	fmt.Println("newContract sender :", callerAddress.Address().Hex())
	contract := vm.NewContract(callerAddress, callerAddress, value, uint64(initGas))
	return contract
}

func build_staking_data(genesisHash common.Hash) {

	stakingDB := staking.NewStakingDB()
	sndb.NewBlock(big.NewInt(1), genesisHash, blockHash)
	// MOCK

	nodeId_A := nodeIdArr[0]
	addr_A, _ := xutil.NodeId2Addr(nodeId_A)

	nodeId_B := nodeIdArr[1]
	addr_B, _ := xutil.NodeId2Addr(nodeId_B)

	nodeId_C := nodeIdArr[2]
	addr_C, _ := xutil.NodeId2Addr(nodeId_C)

	//canArr := make(staking.CandidateQueue, 0)

	c1 := &staking.Candidate{
		NodeId:             nodeId_A,
		StakingAddress:     sender,
		BenefitAddress:     addrArr[1],
		StakingTxIndex:     uint32(2),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "xxccccdddddddd",
			NodeName:   "I Am " + fmt.Sprint(1),
			Website:    "www.baidu.com",
			Details:    "this is  baidu ~~",
		},
	}

	c2 := &staking.Candidate{
		NodeId:             nodeId_B,
		StakingAddress:     sender,
		BenefitAddress:     addrArr[2],
		StakingTxIndex:     uint32(3),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "SFSFSFSFSFSFSSFS",
			NodeName:   "I Am " + fmt.Sprint(2),
			Website:    "www.JD.com",
			Details:    "this is  JD ~~",
		},
	}

	c3 := &staking.Candidate{
		NodeId:             nodeId_C,
		StakingAddress:     sender,
		BenefitAddress:     addrArr[3],
		StakingTxIndex:     uint32(4),
		ProgramVersion:     uint32(1),
		Status:             staking.Valided,
		StakingEpoch:       uint32(1),
		StakingBlockNum:    uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "FWAGGDGDGG",
			NodeName:   "I Am " + fmt.Sprint(3),
			Website:    "www.alibaba.com",
			Details:    "this is  alibaba ~~",
		},
	}

	//canArr = append(canArr, c1)
	//canArr = append(canArr, c2)
	//canArr = append(canArr, c3)

	stakingDB.SetCanPowerStore(blockHash, addr_A, c1)
	stakingDB.SetCanPowerStore(blockHash, addr_B, c2)
	stakingDB.SetCanPowerStore(blockHash, addr_C, c3)

	stakingDB.SetCandidateStore(blockHash, addr_A, c1)
	stakingDB.SetCandidateStore(blockHash, addr_B, c2)
	stakingDB.SetCandidateStore(blockHash, addr_C, c3)

	queue := make(staking.ValidatorQueue, 0)

	v1 := &staking.Validator{
		NodeAddress:   addr_A,
		NodeId:        c1.NodeId,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c1.StakingBlockNum), fmt.Sprint(c1.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v2 := &staking.Validator{
		NodeAddress:   addr_B,
		NodeId:        c2.NodeId,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c2.StakingBlockNum), fmt.Sprint(c2.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v3 := &staking.Validator{
		NodeAddress:   addr_C,
		NodeId:        c3.NodeId,
		StakingWeight: [staking.SWeightItem]string{"1", common.Big256.String(), fmt.Sprint(c3.StakingBlockNum), fmt.Sprint(c3.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	queue = append(queue, v1)
	queue = append(queue, v2)
	queue = append(queue, v3)

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
}

func buildDbRestrictingPlan(t *testing.T, account common.Address, balance *big.Int, epochs int, stateDB xcom.StateDB) {
	//account := addrArr[0]
	fmt.Println("buildDbRestrictingPlan, store addr", account.Hex())
	//const Epochs = 5
	var list = make([]uint64, 0)

	for epoch := 1; epoch <= epochs; epoch++ {
		// build release account record
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(epoch), 1)
		stateDB.SetState(cvm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

		// build release amount record 1eth
		releaseAmount := balance
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(epoch), account)
		stateDB.SetState(cvm.RestrictingContractAddr, releaseAmountKey, releaseAmount.Bytes())

		// build release epoch list record
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(epoch))
		stateDB.SetState(cvm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		list = append(list, uint64(epoch))
	}

	// build restricting user info
	var user restricting.RestrictingInfo
	user.Balance = new(big.Int).Mul(balance, big.NewInt(int64(epochs)))
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

	von, _ := new(big.Int).SetString("100000000000000000000000000", 10)
	stateDB.AddBalance(cvm.RestrictingContractAddr, von)
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
