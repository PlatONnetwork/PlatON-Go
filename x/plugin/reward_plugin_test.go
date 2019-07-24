package plugin_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/go-errors/errors"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

var coinbase = common.HexToAddress("0xc7d92a06e9824955f1504b55114b0daad434e79")

// buildDBSimpleCandidate will create a candidate with simple info,
// then store it into snapshot database.
func buildDBSimpleCandidate(t *testing.T, snapDb snapshotdb.DB) error {
	url := "enode://0x7bae841405067598bf65e7260ca693a964316e752249c4970085c805dbee738fdb41fc434e96e2b65e8bf1db2f52f05d9300d04c1e6129c26cb5d0f214b49968@platon.network:16791"
	node, _ := discover.ParseNode(url)
	can := &staking.Candidate{
		NodeId:         node.ID,
		BenefitAddress: addrArr[1],
	}

	nodeAddr, err := xutil.NodeId2Addr(can.NodeId)
	if err != nil {
		t.Errorf("Failed to convert nodeID to address. ID:%v, error:%s", can.NodeId, err)
		return errors.New(" failed to exchange nodeID to addr")
	}

	key := staking.CandidateKeyByAddr(nodeAddr)

	if val, err := rlp.EncodeToBytes(can); nil != err {
		t.Errorf("Failed to Store Candidate Info: rlp encodeing failed. ID:%v, error:%s", can.NodeId, err)
		return errors.New("failed to rlp encode candidate")
	} else {
		if err := snapDb.PutBaseDB(key, val); nil != err {
			t.Errorf("Failed to Store Candidate Info: PutBaseDB failed. ID:%v, error:%s", can.NodeId, err)
			return errors.New("Failed to Store Candidate Info")
		}
	}
	return nil
}

func buildVerifierList(t *testing.T, blockNumber uint64) error {
	url := "enode://0x7bae841405067598bf65e7260ca693a964316e752249c4970085c805dbee738fdb41fc434e96e2b65e8bf1db2f52f05d9300d04c1e6129c26cb5d0f214b49968@platon.network:16791"
	node, _ := discover.ParseNode(url)
	addr, err := xutil.NodeId2Addr(node.ID)
	if err != nil {
		t.Error("exchange NodeID to Address failed")
		return err
	}

	v := &staking.Validator{
		NodeAddress: addr,
		NodeId:      node.ID,
		StakingWeight: [staking.SWeightItem]string{fmt.Sprint(xutil.CalcVersion(initProgramVersion)), string(100),
			string(blockNumber), string(0)},
		ValidatorTerm: 0,
	}
	validatorArr := make(staking.ValidatorQueue, 0)
	queue := append(validatorArr, v)
	oneEpochBlocks := xutil.CalcBlocksEachEpoch()
	indexInfo := &staking.ValArrIndex{
		Start: (blockNumber-1)/oneEpochBlocks*oneEpochBlocks + 1,
		End:   (blockNumber-1)/oneEpochBlocks*oneEpochBlocks + oneEpochBlocks,
	}
	indexQueue := make(staking.ValArrIndexQueue, 0)
	indexQueue = append(indexQueue, indexInfo)
	if indexArr, err := rlp.EncodeToBytes(indexQueue); nil != err {
		t.Errorf("Failed to Store Epoch Validators indexQueue: rlp encodeing failed. error:%s", err)
		return err
	} else {
		if err := snapdb.PutBaseDB(staking.GetEpochIndexKey(), indexArr); err != nil {
			t.Errorf("Failed to Store Epoch Validators indexQueue: PutBaseDB failed. error:%s", err)
			return err
		}
	}

	if verifiers, err := rlp.EncodeToBytes(queue); err != nil {
		t.Errorf("Failed to Store Epoch Validators: rlp encodeing failed. error:%s", err)
		return err

	} else {
		snapdb := snapshotdb.Instance()
		if err := snapdb.PutBaseDB(staking.GetEpochValArrKey(indexInfo.Start, indexInfo.End), verifiers); err != nil {
			t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err)
			return err
		}
	}

	return nil
}

// buildDbRewardData accord to blockNumber calculate reward data and store into level levelDb.
// We assume this is only used for testing the fist three years increase issuance.
// We also don't want the block is the last block of the year
func buildDbRewardData(t *testing.T, state xcom.StateDB, blockNumber uint64) {

	oneYearBlocks := xutil.CalcBlocksEachYear()
	oneYearEpochs := xutil.EpochsPerYear()
	t.Logf("one year blocks are: %v", oneYearBlocks)
	t.Logf("one year epochs are: %v", oneYearEpochs)

	if xutil.IsYearEnd(blockNumber) {
		t.Fatalf("unexpected block, blockNumber: %d", blockNumber)
	}

	year := xutil.CalculateYear(blockNumber)
	switch year {
	case 1:
		t.Log("------")
		t.Log("current blockNumber is ", blockNumber)
		totalRewardBalance := plugin.GetYearEndBalance(state, 0)
		t.Logf("genesis reward balance is: %v", totalRewardBalance)

		// calculate new block reward
		totalNBReward := new(big.Int).Div(totalRewardBalance, big.NewInt(plugin.RewardNewBlockRate))
		oneBlockReward := new(big.Int).Div(totalNBReward, big.NewInt(int64(oneYearBlocks)))
		rewardNB := new(big.Int).Mul(oneBlockReward, big.NewInt(int64(blockNumber)))
		t.Logf("one block reward at the first year is: %v, total is: %v", oneBlockReward, rewardNB)

		// calculate staking reward
		totalSTKReward := new(big.Int).Sub(totalRewardBalance, totalNBReward)
		oneEpochSTKReward := new(big.Int).Div(totalSTKReward, big.NewInt(int64(oneYearEpochs)))
		epoch := xutil.CalculateEpoch(blockNumber)
		rewardSTK := new(big.Int).Mul(oneEpochSTKReward, big.NewInt(int64(epoch-1)))
		t.Logf("one epoch staking reward is %v, current epoch is %d", oneEpochSTKReward, epoch)
		t.Logf("staking reward is sent to RewardManagerPoolAddr, total is: %v", rewardSTK)

		// store balance of coinbase
		state.AddBalance(coinbase, rewardNB)
		t.Logf("balance of coinbase at block %d is: %v", blockNumber, state.GetBalance(coinbase))

		// store balance of reward pool
		state.SubBalance(vm.RewardManagerPoolAddr, rewardNB)
		t.Logf("balance of reward pool at block %d is: %v", blockNumber, state.GetBalance(vm.RewardManagerPoolAddr))

		// store snapshot data
		if err := buildVerifierList(t, blockNumber); err != nil {
			t.Fatal(err.Error())
		}

	case 2:
		t.Log("------")
		t.Log("current blockNumber is ", blockNumber)
		totalRewardBalance := plugin.GetYearEndBalance(state, 1)
		t.Logf("reward balance at the first year end is: %v", totalRewardBalance)

		// calculate new block reward
		totalNBReward := new(big.Int).Div(totalRewardBalance, big.NewInt(plugin.RewardNewBlockRate))
		oneBlockReward := new(big.Int).Div(totalNBReward, big.NewInt(int64(oneYearBlocks)))
		rewardNB := new(big.Int).Mul(oneBlockReward, big.NewInt(int64(blockNumber-oneYearBlocks)))
		t.Logf("one block reward at the twice year is: %v", oneBlockReward)

		// calculate staking reward
		totalSTKReward := new(big.Int).Sub(totalRewardBalance, totalNBReward)
		oneEpochSTKReward := new(big.Int).Div(totalSTKReward, big.NewInt(int64(oneYearEpochs)))
		epoch := xutil.CalculateEpoch(blockNumber - oneYearBlocks)
		rewardSTK := new(big.Int).Mul(oneEpochSTKReward, big.NewInt(int64(epoch-1)))
		t.Logf("one epoch staking reward at the secend year is %v, current epoch is %d", oneEpochSTKReward, epoch)
		t.Logf("staking reward is sent to RewardManagerPoolAddr, total is: %v", rewardSTK)

		// store balance of coinbase
		state.AddBalance(coinbase, rewardNB)
		t.Logf("balance of coinbase at block %d is: %v", blockNumber, state.GetBalance(coinbase))

		// store balance of reward pool
		state.SubBalance(vm.RewardManagerPoolAddr, rewardNB)
		t.Logf("balance of reward pool at block %d is: %v", blockNumber, state.GetBalance(vm.RewardManagerPoolAddr))

		// store snapshot data
		if err := buildVerifierList(t, blockNumber); err != nil {
			t.Fatal(err.Error())
		}

	default:
		t.Fatalf("unexpected block, blockNumber: %d, year: %d", blockNumber, year)
	}

	t.Log("end of buildDbRewardData")
	t.Log("==============")
}

func testEndBlockWithOneBlock(t *testing.T) error {
	state, _, err := newChainState()
	if err != nil {
		t.Error("create genesis block failed")
		return err
	}

	// do end block
	headOne := &types.Header{Number: big.NewInt(1), Coinbase: coinbase}
	if err = plugin.RewardMgrInstance().EndBlock(common.ZeroHash, headOne, state); err != nil {
		t.Error("test end block failed")
		return err
	}

	totalRewardBalance := plugin.GetYearEndBalance(state, 0)
	if totalRewardBalance == nil {
		t.Fatal("system error, GetYearEndBalance is empty")
	}

	totalNBReward := new(big.Int).Div(totalRewardBalance, big.NewInt(plugin.RewardNewBlockRate))

	totalBlocks := xutil.CalcBlocksEachYear()
	oneBlockNBReward := new(big.Int).Div(totalNBReward, big.NewInt(int64(totalBlocks)))

	t.Logf("expected balance of coinbase is: %v", totalRewardBalance.Sub(totalRewardBalance, oneBlockNBReward))
	t.Logf("expected balance of reward pool is: %v", oneBlockNBReward)

	t.Logf("actually balance of coinbase is: %v", state.GetBalance(coinbase))
	t.Logf("actually balance of reward pool is: %v", state.GetBalance(vm.RewardManagerPoolAddr))

	t.Log("test one block pass")
	t.Log("=====================")

	return nil
}

func testEndBlockWithOneEpochBlock(t *testing.T) error {
	state, _, err := newChainState()
	if err != nil {
		t.Error("create genesis block failed")
		return err
	}

	// do end block
	oneEpochBlocks := xutil.CalcBlocksEachEpoch()
	t.Logf("one epoch blocks are:%v", oneEpochBlocks)
	for blockNumber := int64(1); blockNumber <= int64(oneEpochBlocks); blockNumber++ {
		head := &types.Header{Number: big.NewInt(blockNumber), Coinbase: coinbase}
		if err = plugin.RewardMgrInstance().EndBlock(common.ZeroHash, head, state); err != nil {
			t.Errorf("test end block failed, err:%s", err.Error())
			return err
		}
	}

	totalRewardBalance := plugin.GetYearEndBalance(state, 0)
	if totalRewardBalance == nil {
		t.Fatal("system error, GetYearEndBalance is empty")
	}

	oneYearBlocks := xutil.CalcBlocksEachYear()
	oneYearEpochs := xutil.EpochsPerYear()
	t.Logf("one year blocks are: %v", oneYearBlocks)
	t.Logf("one year epochs are: %v", oneYearEpochs)

	totalNBReward := new(big.Int).Div(totalRewardBalance, big.NewInt(plugin.RewardNewBlockRate))
	oneBlockNBReward := new(big.Int).Div(totalNBReward, big.NewInt(int64(oneYearBlocks)))
	oneEpochNBReward := new(big.Int).Mul(oneBlockNBReward, big.NewInt(int64(oneEpochBlocks)))

	totalSTKReward := new(big.Int).Sub(totalRewardBalance, totalNBReward)
	oneEpochSTKReward := new(big.Int).Mul(totalSTKReward, big.NewInt(int64(oneYearEpochs)))
	t.Log("expected staking reward is sent to RewardManagerPoolAddr, amount is: ", oneEpochSTKReward)

	t.Logf("expected balance of coinbase is: %v", oneEpochNBReward)
	t.Logf("expected balance of reward pool is: %v", totalRewardBalance.Sub(totalRewardBalance, oneEpochNBReward))

	t.Logf("actually balance of coinbase is: %v", state.GetBalance(coinbase))
	t.Logf("actually balance of reward pool is: %v", state.GetBalance(vm.RewardManagerPoolAddr))
	t.Log("test one year block pass")
	t.Log("=====================")

	return nil
}

func testEndBlockWithOneYearBlock(t *testing.T) error {
	state, _, err := newChainState()
	if err != nil {
		t.Error("create genesis block failed")
		return err
	}

	firstYearEndBlocks := xutil.CalcBlocksEachYear()
	buildDataBlock := firstYearEndBlocks - 1
	buildDbRewardData(t, state, buildDataBlock)

	// calculate expected balance
	totalRewardBalance := plugin.GetYearEndBalance(state, 0)
	totalNBReward := new(big.Int).Div(totalRewardBalance, big.NewInt(plugin.RewardNewBlockRate))
	oneBlockReward := new(big.Int).Div(totalNBReward, big.NewInt(int64(firstYearEndBlocks)))

	// calculate expected balance after reward at the year end block
	balancePool := state.GetBalance(vm.RewardManagerPoolAddr)
	balanceYearEndBlock := new(big.Int).Sub(balancePool, oneBlockReward)

	histIssuance := plugin.GetHistoryCumulativeIssue(state, 0)
	currIssuance := new(big.Int).Div(histIssuance, big.NewInt(40))
	devIssuance := new(big.Int).Div(currIssuance, big.NewInt(5))
	rewardIssuance := currIssuance.Sub(currIssuance, devIssuance)

	totalRewardBalance = totalRewardBalance.Add(balanceYearEndBlock, rewardIssuance)
	t.Logf("expected balance of reward Pool is %v", totalRewardBalance)

	coinBaseBalance := state.GetBalance(coinbase)
	coinBaseBalance = new(big.Int).Add(coinBaseBalance, oneBlockReward)
	t.Logf("expected balance of coinbase is %v", coinBaseBalance)

	head := &types.Header{
		Number:   big.NewInt(int64(firstYearEndBlocks)),
		Coinbase: coinbase,
	}
	if err = plugin.RewardMgrInstance().EndBlock(head.Hash(), head, state); err != nil {
		t.Errorf("test end block failed, err:%s", err.Error())
		return err
	}

	t.Logf("actually balance of reward pool is: %v", state.GetBalance(vm.RewardManagerPoolAddr))
	t.Logf("actually balance of coinbase is: %v", state.GetBalance(coinbase))
	t.Log("test one year block pass")
	t.Log("=====================")

	return nil
}

func testEndBlockWithTwoYearBlock(t *testing.T) error {
	state, _, err := newChainState()
	if err != nil {
		t.Error("create genesis block failed")
		return err
	}

	// build db data of the first year end block
	firstYearEndBlocks := xutil.CalcBlocksEachYear()
	buildDataBlock := firstYearEndBlocks - 1
	buildDbRewardData(t, state, buildDataBlock)

	head := &types.Header{
		Number:   big.NewInt(int64(firstYearEndBlocks)),
		Coinbase: coinbase,
	}
	if err = plugin.RewardMgrInstance().EndBlock(head.Hash(), head, state); err != nil {
		t.Errorf("test end block failed, err:%s", err.Error())
		return err
	}

	// build db data
	buildDataBlock = 2*firstYearEndBlocks - 1
	buildDbRewardData(t, state, buildDataBlock)

	// calculate expected balance
	totalRewardBalance := plugin.GetYearEndBalance(state, 1)
	totalNBReward := new(big.Int).Div(totalRewardBalance, big.NewInt(plugin.RewardNewBlockRate))
	oneBlockReward := new(big.Int).Div(totalNBReward, big.NewInt(int64(firstYearEndBlocks)))

	// calculate expected balance after reward at the year end block
	balancePool := state.GetBalance(vm.RewardManagerPoolAddr)
	balanceYearEndBlock := new(big.Int).Sub(balancePool, oneBlockReward)

	histIssuance := plugin.GetHistoryCumulativeIssue(state, 1)
	currIssuance := new(big.Int).Div(histIssuance, big.NewInt(40))
	devIssuance := new(big.Int).Div(currIssuance, big.NewInt(5))
	rewardIssuance := currIssuance.Sub(currIssuance, devIssuance)

	totalRewardBalance = totalRewardBalance.Add(balanceYearEndBlock, rewardIssuance)
	t.Logf("expected balance of reward Pool is %v", totalRewardBalance)

	coinBaseBalance := state.GetBalance(coinbase)
	coinBaseBalance = new(big.Int).Add(coinBaseBalance, oneBlockReward)
	t.Logf("expected balance of coinbase is %v", coinBaseBalance)

	head = &types.Header{
		Number:   big.NewInt(int64(2 * firstYearEndBlocks)),
		Coinbase: coinbase,
	}
	if err = plugin.RewardMgrInstance().EndBlock(head.Hash(), head, state); err != nil {
		t.Errorf("test end block failed, err:%s", err.Error())
		return err
	}

	t.Logf("actually balance of reward pool is: %v", state.GetBalance(vm.RewardManagerPoolAddr))
	t.Logf("actually balance of coinbase is: %v", state.GetBalance(coinbase))
	t.Log("test one year block pass")
	t.Log("=====================")
	return nil
}

func TestRewardMgrPlugin_EndBlock(t *testing.T) {
	var err error

	/*
	 * Branch coverage
	 */
	// case1: current is common block
	{
		xcom.GetEc(xcom.DefaultDeveloperNet)
		stateDb := buildStateDB(t)

		totalReward, _ := new(big.Int).SetString("65000000000000000000000000", 10)
		stateDb.AddBalance(vm.RewardManagerPoolAddr, totalReward)
		plugin.SetYearEndBalance(stateDb, 0, totalReward)

		// show expected result
		totalNBReward := new(big.Int).Div(totalReward, big.NewInt(plugin.RewardNewBlockRate))
		oneNBReward := new(big.Int).Div(totalNBReward, big.NewInt(int64(xutil.CalcBlocksEachYear())))
		t.Log("expected case1 of EndBlock only reward new block success")
		t.Logf("expected balance of coinbase is: %v", oneNBReward)
		t.Logf("expected balance of reward pool is: %v", totalReward.Sub(totalReward, oneNBReward))

		currBlockNumber := 1
		head := types.Header{Number: big.NewInt(int64(currBlockNumber)), Coinbase: addrArr[0]}
		if err = plugin.RewardMgrInstance().EndBlock(blockHash, &head, stateDb); err != nil {
			t.Errorf("case1 of EndBlock failed. actually returns error: %v", err.Error())

		} else {

			t.Log("case1 returns Success")
			t.Logf("actually balance of coinbase is %v", stateDb.GetBalance(addrArr[0]))
			t.Logf("actually balance of reward pool is %v", stateDb.GetBalance(vm.RewardManagerPoolAddr))
			t.Log("case1 pass")
			t.Log("=====================")
		}
	}

	// case2: current is settle block
	{

		xcom.GetEc(xcom.DefaultDeveloperNet)
		stateDb := buildStateDB(t)
		snapDb := snapshotdb.Instance()

		if err := buildDBSimpleCandidate(t, snapDb); err != nil {
			t.Fatalf(err.Error())
		}
		if err := buildVerifierList(t, 1); err != nil {
			t.Fatal(err.Error())
		}

		// restore data in levelDB
		totalReward, _ := new(big.Int).SetString("65000000000000000000000000", 10)
		stateDb.AddBalance(vm.RewardManagerPoolAddr, totalReward)
		plugin.SetYearEndBalance(stateDb, 0, totalReward)

		// show expected result
		totalNBReward := new(big.Int).Div(totalReward, big.NewInt(plugin.RewardNewBlockRate))
		oneNBReward := new(big.Int).Div(totalNBReward, big.NewInt(int64(xutil.CalcBlocksEachYear())))
		t.Log("expected case2 of EndBlock reward staking and reward new block success")
		t.Logf("expected balance of coinbase is: %v", oneNBReward)

		totalSTKReward := new(big.Int).Sub(totalReward, totalNBReward)
		oneEpochSTKReward := new(big.Int).Div(totalSTKReward, big.NewInt(int64(xutil.EpochsPerYear())))
		t.Logf("expected balance of staking reward address is: %v", oneEpochSTKReward)

		totalReward = totalReward.Sub(totalReward, oneNBReward)
		totalReward = totalReward.Sub(totalReward, oneEpochSTKReward)
		t.Logf("expected balance of reward pool is: %v", totalReward)

		currBlockNumber := uint64(1) * xutil.CalcBlocksEachEpoch()
		head := types.Header{Number: big.NewInt(int64(currBlockNumber)), Coinbase: addrArr[0]}
		if err = plugin.RewardMgrInstance().EndBlock(blockHash, &head, stateDb); err != nil {
			t.Errorf("case2 of EndBlock failed. error: %s", err.Error())
		} else {
			t.Log("case2 returns Success")
			t.Logf("actually balance of coinbase is: %v", stateDb.GetBalance(addrArr[0]))
			t.Logf("actually balance of staking reward address is: %v", stateDb.GetBalance(addrArr[1]))
			t.Logf("actually balance of reward pool is: %v", stateDb.GetBalance(vm.RewardManagerPoolAddr))
			t.Log("case2 pass")
			t.Log("=====================")
		}
	}

	// case3: current is end of year
	{
		xcom.GetEc(xcom.DefaultDeveloperNet)
		stateDb := buildStateDB(t)
		snapDb := snapshotdb.Instance()
		currBlockNumber := uint64(1) * xutil.CalcBlocksEachYear()

		if err := buildDBSimpleCandidate(t, snapDb); err != nil {
			t.Fatalf(err.Error())
		}
		if err := buildVerifierList(t, currBlockNumber-1); err != nil {
			t.Fatal(err.Error())
		}

		// restore data in levelDB
		histIssuance, _ := new(big.Int).SetString("1000000000000000000000000000", 10)
		totalReward, _ := new(big.Int).SetString("65000000000000000000000000", 10)
		stateDb.AddBalance(vm.RewardManagerPoolAddr, totalReward)
		plugin.SetYearEndBalance(stateDb, 0, totalReward)
		plugin.SetYearEndCumulativeIssue(stateDb, 0, histIssuance)

		// calculate expected balance
		currIssuance := new(big.Int).Div(histIssuance, big.NewInt(40))
		devIssuance := new(big.Int).Div(currIssuance, big.NewInt(5))
		rewardIssuance := new(big.Int).Sub(currIssuance, devIssuance)

		// show expected balance
		totalNBReward := new(big.Int).Div(totalReward, big.NewInt(plugin.RewardNewBlockRate))
		oneNBReward := new(big.Int).Div(totalNBReward, big.NewInt(int64(xutil.CalcBlocksEachYear())))
		t.Log("expected case3 of EndBlock returns success")
		t.Logf("expected balance of coinbase is: %v", oneNBReward)

		totalSTKReward := new(big.Int).Sub(totalReward, totalNBReward)
		oneEpochSTKReward := new(big.Int).Div(totalSTKReward, big.NewInt(int64(xutil.EpochsPerYear())))
		t.Logf("expected balance of staking reward address is: %v", oneEpochSTKReward)

		totalReward = totalReward.Sub(totalReward, oneNBReward)
		totalReward = totalReward.Sub(totalReward, oneEpochSTKReward)
		totalReward = totalReward.Add(totalReward, rewardIssuance)
		t.Logf("expected balance of reward pool is: %v", totalReward)
		t.Logf("expected balance if developer fundtion is: %v", devIssuance)

		head := types.Header{Number: big.NewInt(int64(currBlockNumber)), Coinbase: addrArr[0]}
		if err = plugin.RewardMgrInstance().EndBlock(blockHash, &head, stateDb); err != nil {
			t.Errorf("case3 of EndBlock failed. %v", err.Error())
		} else {
			t.Log("case3 returns Success")
			t.Logf("actually balance of coinbase is: %v", stateDb.GetBalance(addrArr[0]))
			t.Logf("actually balance of staking reward address is: %v", stateDb.GetBalance(addrArr[1]))
			t.Logf("actually balance of reward pool is: %v", stateDb.GetBalance(vm.RewardManagerPoolAddr))
			t.Logf("actually balance of developer fundtion is: %v", stateDb.GetBalance(vm.CommunityDeveloperFoundation))
			t.Log("case3 pass")
			t.Log("=====================")
		}
	}

	/*
	 * Path coverage
	 */
	// case4: test one block
	t.Log("test one block, testing reward new block")
	if err := testEndBlockWithOneBlock(t); err != nil {
		t.Fatal(err.Error())
	}

	// case5: test one settle epoch blocks
	t.Log("test one epoch block, testing reward new block and staking")
	if err := testEndBlockWithOneEpochBlock(t); err != nil {
		t.Fatalf(err.Error())
	}

	// case6: test one year blocks
	t.Log("test one year block, testing reward and increase issuance")
	if err := testEndBlockWithOneYearBlock(t); err != nil {
		t.Fatalf(err.Error())
	}

	// case7: test two year blocks
	t.Log("test two year block, testing reward and increase issuance")
	if err := testEndBlockWithTwoYearBlock(t); err != nil {
		t.Fatalf(err.Error())
	}
}
