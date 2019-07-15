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
	newChainState()
	xcom.SetEconomicModel(&xcom.DefaultConfig)

	stateDB := buildStateDB(t)

	// case1: current is common block
	head := types.Header{ Number: big.NewInt(1),}
	if err := plugin.RewardMgrInstance().EndBlock(blockHash, &head, stateDB); err != nil {
		t.Error("The case1 of EndBlock failed.\n expected err is nil")
		t.Errorf("Actually returns err. blockNumber:%d . errors: %s", head.Number.Uint64(), err.Error())
	} else {
		t.Log("Success")
	}

	// case2: current is settle block
	head = types.Header{ Number: big.NewInt(int64(1*xcom.ConsensusSize()*xcom.EpochSize())),}
	if err := plugin.RewardMgrInstance().EndBlock(blockHash, &head, stateDB); err != nil {
		t.Error("The case2 of EndBlock failed.\n expected err is nil")
		t.Errorf("Actually returns err. blockNumber:%d . errors: %s", head.Number.Uint64(), err.Error())
	} else {
		t.Log("Success")
	}

	// case3: current is end of year
	head = types.Header{ Number: big.NewInt(int64(365*24*3600)),}
	if err := plugin.RewardMgrInstance().EndBlock(blockHash, &head, stateDB); err != nil {
		t.Error("The case2 of EndBlock failed.\n expected err is nil")
		t.Errorf("Actually returns err. blockNumber:%d . errors: %s", head.Number.Uint64(), err.Error())
	} else {
		t.Log("Success")
	}
}
