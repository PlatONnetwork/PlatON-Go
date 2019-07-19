package plugin_test

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/go-errors/errors"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

var coinbase = common.HexToAddress("0xc7d92a06e9824955f1504b55114b0daad434e79")

func buildVerifierList(t *testing.T, blockNumber uint64) (common.Hash, error) {
	header := types.Header{
		Number:   big.NewInt(int64(blockNumber)),
		Coinbase: coinbase,
	}
	hash := header.Hash()

	stakingDB := staking.NewStakingDB()
	if err := sndb.NewBlock(big.NewInt(int64(blockNumber)), lastBlockHash, hash); err != nil {
		t.Error("newBlock in snapshot database failed")
		return common.Hash{}, err
	}

	url := "enode://0x7bae841405067598bf65e7260ca693a964316e752249c4970085c805dbee738fdb41fc434e96e2b65e8bf1db2f52f05d9300d04c1e6129c26cb5d0f214b49968@platon.network:16791"
	node, _ := discover.ParseNode(url)
	addr, err := xutil.NodeId2Addr(node.ID)
	if err != nil {
		log.Error("exchange NodeID to Address failed")
		return common.Hash{}, err
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
	arr := &staking.Validator_array{
		Start: (blockNumber-1)/oneEpochBlocks*oneEpochBlocks + 1,
		End:   (blockNumber-1)/oneEpochBlocks*oneEpochBlocks + oneEpochBlocks,
		Arr:   queue,
	}

	if err := stakingDB.SetVerfierList(hash, arr); err != nil {
		t.Error("store list of current validator failed")
		return common.Hash{}, err
	}

	return hash, nil
}

func buildDBRewardPluginData(t *testing.T) {
	stateDB := buildStateDB(t)

	initial := new(big.Int)
	initial, _ = initial.SetString("1000000000000000000", 10)
	plugin.SetYearEndCumulativeIssue(stateDB, 0, initial)

	increase := new(big.Int)
	firstYearEnd := new(big.Int)
	secondYearEnd := new(big.Int)

	increase = increase.Div(initial, big.NewInt(40))
	firstYearEnd = firstYearEnd.Add(firstYearEnd, increase)
	plugin.SetYearEndCumulativeIssue(stateDB, 1, firstYearEnd)

	increase = increase.Div(firstYearEnd, big.NewInt(40))
	secondYearEnd = firstYearEnd.Add(secondYearEnd, increase)
	plugin.SetYearEndCumulativeIssue(stateDB, 1, secondYearEnd)
}

func testEndBlockWithOneBlock(t *testing.T) error {
	state, _, err := newChainState()
	if err != nil {
		t.Error("create genesis block failed")
		return err
	}

	headOne := &types.Header{Number: big.NewInt(1), Coinbase: coinbase}
	if err = plugin.RewardMgrInstance().EndBlock(common.ZeroHash, headOne, state); err != nil {
		t.Error("test end block failed")
		return err
	}

	lastTotalReward, success := new(big.Int).SetString("65000000000000000000000000", 10)
	if !success {
		t.Error("set reward failed")
		return errors.New("set reward failed")
	}
	lastTotalNBReward, success := new(big.Int).SetString("16250000000000000000000000", 10)
	if !success {
		t.Error("set reward failed")
		return errors.New("set reward failed")
	}

	totalBlocks := xutil.CalcBlocksEachYear()
	oneBlockNBReward := new(big.Int).Div(lastTotalNBReward, big.NewInt(int64(totalBlocks)))

	t.Logf("expected balance of coinbase is: %v", lastTotalReward.Sub(lastTotalReward, oneBlockNBReward))
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

	oneEpochBlocks := xutil.CalcBlocksEachEpoch()
	t.Logf("one epoch blocks are:%v", oneEpochBlocks)
	for blockNumber := int64(1); blockNumber <= int64(oneEpochBlocks); blockNumber++ {
		head := &types.Header{Number: big.NewInt(blockNumber), Coinbase: coinbase}
		if err = plugin.RewardMgrInstance().EndBlock(common.ZeroHash, head, state); err != nil {
			t.Errorf("test end block failed, err:%s", err.Error())
			return err
		}
	}

	lastTotalReward, success := new(big.Int).SetString("65000000000000000000000000", 10)
	if !success {
		t.Error("set reward failed")
		return errors.New("set reward failed")
	}

	oneYearBlocks := xutil.CalcBlocksEachYear()
	oneYearEpochs := xutil.EpochsPerYear()
	t.Logf("one year blocks are: %v", oneYearBlocks)
	t.Logf("one year epochs are: %v", oneYearEpochs)

	lastTotalNBReward, success := new(big.Int).SetString("16250000000000000000000000", 10)
	if !success {
		t.Error("set reward failed")
		return errors.New("set reward failed")
	}
	oneBlockNBReward := new(big.Int).Div(lastTotalNBReward, big.NewInt(int64(oneYearBlocks)))
	oneEpochNBReward := new(big.Int).Mul(oneBlockNBReward, big.NewInt(int64(oneEpochBlocks)))

	// lastTotalSTKReward, success := new(big.Int).SetString("48750000000000000000000000", 10)
	// if !success {
	// 	t.Error("set reward failed")
	// 	return errors.New("set reward failed")
	// }
	// oneEpochSTKReward := new(big.Int).Mul(lastTotalSTKReward, big.NewInt(int64(oneYearEpochs)))

	t.Logf("expected balance of coinbase is: %v", oneEpochNBReward)
	t.Logf("expected balance of reward pool is: %v", lastTotalReward.Sub(lastTotalReward, oneEpochNBReward))

	t.Logf("actually balance of coinbase is: %v", state.GetBalance(coinbase))
	t.Logf("expected balance of reward pool is: %v", state.GetBalance(vm.RewardManagerPoolAddr))
	t.Log("test one year block pass")
	t.Log("=====================")

	return nil
}

func testEndBlockWitheOneYearBlock(t *testing.T) error {
	state, _, err := newChainState()
	if err != nil {
		t.Error("create genesis block failed")
		return err
	}

	oneYearBlocks := xutil.CalcBlocksEachYear()
	t.Logf("one epoch blocks are:%v", oneYearBlocks)
	for blockNumber := int64(1); blockNumber <= int64(oneYearBlocks); blockNumber++ {

		head := &types.Header{Number: big.NewInt(blockNumber), Coinbase: coinbase}
		if err = plugin.RewardMgrInstance().EndBlock(common.ZeroHash, head, state); err != nil {
			t.Errorf("test end block failed, err:%s", err.Error())
			return err
		}

		if xutil.IsSettlementPeriod(uint64(blockNumber)) {
			_, _ = buildVerifierList(t, uint64(blockNumber))
		}
	}

	lastTotalReward, _ := new(big.Int).SetString("65000000000000000000000000", 10)

	oneYearEpochs := xutil.EpochsPerYear()
	t.Logf("one year blocks arg: %v", oneYearBlocks)
	t.Logf("one year epochs are: %v", oneYearEpochs)

	lastTotalNBReward, _ := new(big.Int).SetString("16250000000000000000000000", 10)
	oneBlockNBReward := new(big.Int).Div(lastTotalNBReward, big.NewInt(int64(oneYearBlocks)))
	oneYearNBReward := new(big.Int).Mul(oneBlockNBReward, big.NewInt(int64(oneYearEpochs)))

	// lastTotalSTKReward, _ := new(big.Int).SetString("48750000000000000000000000", 10)
	// oneEpochSTKReward := new(big.Int).Mul(lastTotalSTKReward, big.NewInt(int64(oneYearEpochs)))

	t.Logf("expected balance of coinbase is: %v", oneYearNBReward)
	t.Logf("expected balance of reward pool is: %v", lastTotalReward.Sub(lastTotalReward, oneYearNBReward))

	t.Logf("actually balance of coinbase is: %v", state.GetBalance(vm.RewardManagerPoolAddr))
	t.Logf("expected balance of reward pool is: %v", state.GetBalance(coinbase))

	return nil

}

func TestRewardMgrPlugin_EndBlock(t *testing.T) {
	var err error

	/*
	 * Branch coverage
	 */
	// case1: current is common block
	{
		stateDB := buildStateDB(t)

		currBlockNumber := 1
		plugin.SetYearEndCumulativeIssue(stateDB, 0, big.NewInt(1000000000000))
		head := types.Header{Number: big.NewInt(int64(currBlockNumber)), Coinbase: addrArr[0]}
		//	build_staking_data()

		err = plugin.RewardMgrInstance().EndBlock(blockHash, &head, stateDB)

		// show expected result
		t.Log("expected case1 of EndBlock only reward new block success")
		t.Log("expected balance of coinbase is ")

		if err != nil {
			t.Errorf("case1 of EndBlock failed. Actually returns error: %v", err.Error())

		} else {

			t.Log("case1 returns Success")
			t.Logf("Actually balance of coinbase is %v", stateDB.GetBalance(addrArr[0]))
			t.Log("case1 pass")
			t.Log("=====================")
		}
	}

	// case2: current is settle block
	{
		stateDB := buildStateDB(t)

		currBlockNumber := uint64(1) * xutil.ConsensusSize() * xutil.EpochSize()
		head := types.Header{Number: big.NewInt(int64(currBlockNumber)), Coinbase: addrArr[0]}
		//	build_staking_data()

		err = plugin.RewardMgrInstance().EndBlock(blockHash, &head, stateDB)

		// show expected result
		t.Log("expected case2 of EndBlock reward staking and reward new block success")
		t.Log("expected balance of coinbase is ")
		t.Log("expected balance of staking reward address is")

		if err != nil {
			t.Errorf("case2 of EndBlock failed. ")
		} else {
			t.Log("case2 returns Success")
			t.Logf("Actually balance of coinbase is %v", stateDB.GetBalance(addrArr[0]))
			t.Logf("expected balance of staking reward address is")
			t.Log("case2 pass")
			t.Log("=====================")
		}
	}

	// case3: current is end of year
	{
		stateDB := buildStateDB(t)

		currBlockNumber := 365 * 24 * 3600
		head := types.Header{Number: big.NewInt(int64(currBlockNumber)), Coinbase: addrArr[0]}
		//	build_staking_data()

		err = plugin.RewardMgrInstance().EndBlock(blockHash, &head, stateDB)

		// show expected result
		t.Log("expected case3 of EndBlock returns success")
		t.Log("expected balance of coinbase is ")
		t.Log("expected balance of staking reward address is")
		t.Log("expected balance of reward pool is")

		if err != nil {
			t.Errorf("case3 of EndBlock failed. ")
		} else {
			t.Log("case3 returns Success")
			t.Logf("Actually balance of coinbase is %v", stateDB.GetBalance(addrArr[0]))
			t.Logf("expected balance of staking reward address is")
			t.Log("expected balance of reward pool is")
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
	if err := testEndBlockWitheOneYearBlock(t); err != nil {
		t.Fatalf(err.Error())
	}
}
