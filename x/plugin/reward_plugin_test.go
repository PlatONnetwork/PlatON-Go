package plugin_test

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)


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


func TestRewardMgrPlugin_EndBlock(t *testing.T) {
	var err error

	// case1: current is common block
	{
		stateDB := buildStateDB(t)

		currBlockNumber := 1
		plugin.SetYearEndCumulativeIssue(stateDB, 0, big.NewInt(1000000000000))
		head := types.Header{ Number: big.NewInt(int64(currBlockNumber)), Coinbase:addrArr[0]}
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

		currBlockNumber := uint64(1) * xcom.ConsensusSize() * xcom.EpochSize()
		head := types.Header{ Number: big.NewInt(int64(currBlockNumber)), Coinbase:addrArr[0]}
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
		head := types.Header{ Number: big.NewInt(int64(currBlockNumber)), Coinbase:addrArr[0]}
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
}
