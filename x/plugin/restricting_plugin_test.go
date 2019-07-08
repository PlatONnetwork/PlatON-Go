package plugin_test

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)


type restrictingInfo struct {
	balance     *big.Int `json:"balance"` // balance representation all locked amount
	debt        *big.Int `json:"debt"`    // debt representation will released amount. Positive numbers can be used instead of release, 0 means no release, negative numbers indicate not enough to release
	releaseList []uint64 `json:"list"`    // releaseList representation which epoch will release restricting
}

type releaseAmountInfo struct {
	height uint64    `json:"blockNumber"`   // blockNumber representation of the block number at the released epoch
	amount *big.Int	 `json:"amount"`		// amount representation of the released amount
}


func buildStateDB (t *testing.T) xcom.StateDB{
	db := ethdb.NewMemDatabase()
	stateDb, err := state.New(common.Hash{}, state.NewDatabase(db))

	if err != nil {
		t.Errorf("new state db failed: %s", err.Error())
	}

	return stateDb
}


func TestRestrictingPlugin_EndBlock(t *testing.T) {

	stateDB := buildStateDB(t)
	buildDbRestrictingPlan(t, stateDB)

	// case1: blockNumber not arrived settle block height
	head := types.Header{ Number: big.NewInt(1),}
	if err := plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDB); err != nil {
		t.Error("The case1 of EndBlock failed.\n expected err is nil")
		t.Errorf("Actually returns err. blockNumber:%d . errors: %s", head.Number.Uint64(), err.Error())
	}

	// case2: blockNumber arrived settle block height, restricting plan not exist
	head = types.Header{Number: big.NewInt(int64(6*xcom.ConsensusSize*xcom.EpochSize))}
	if err := plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDB); err != nil {
		t.Error("The case2 of EndBlock failed.\n expected: ")
		t.Errorf("Actually returns err. blockNumber:%d . errors: %s", head.Number.Uint64(), err.Error())
	} else {
		t.Log("")
	}

	// case3: blockNumber arrived settle block height, restricting plan exist
	head = types.Header{Number: big.NewInt(int64(1*xcom.ConsensusSize*xcom.EpochSize))}
	if err := plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDB); err != nil {
		t.Errorf("Actually returns err. blockNumber:%d . errors: %s", head.Number.Uint64(), err.Error())
	}
}

func TestRestrictingPlugin_AddRestrictingRecord(t *testing.T) {
	stateDB := buildStateDB(t)

	var plan restricting.RestrictingPlan
	var plans = make([]restricting.RestrictingPlan, 5)
	var temp  = plans

	// case1: account is new user to restricting
	for epoch := 1; epoch < 6; epoch++ {
		plan.Epoch = uint64(epoch)
		plan.Amount = big.NewInt(int64(1E18))
		plans = append(plans, plan)
	}

	err := plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDB)
	if err != nil {
		t.Error("The case1 of AddRestrictingRecord failed.")
	}

	// case2: restricting account exist, but restricting epoch not intersect
	plans = make([]restricting.RestrictingPlan, 1)
	plan.Epoch = uint64(6)
	plan.Amount = big.NewInt(int64(1E18))
	plans = append(plans, plan)

	err = plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDB)
	if err != nil {
		t.Error("The case2 of AddRestrictingRecord failed.")
	}

	// case3: restricting account exist, and restricting epoch intersect
	plans = temp
	err = plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDB)
	if err != nil {
		t.Error("The case3 of AddRestrictingRecord failed.")
	}

}

func TestRestrictingPlugin_PledgeLockFunds(t *testing.T) {

	stateDB := buildStateDB(t)
	buildDbRestrictingPlan(t, stateDB)
	lockFunds := big.NewInt(int64(2E18))

	// case1: restricting account not exist
	notFoundAccount := common.HexToAddress("0x11")
	err := plugin.RestrictingInstance().PledgeLockFunds(notFoundAccount, lockFunds, stateDB)
	if err != nil {
		t.Error("The case1 of PledgeLockFunds failed.")
	}

	// case2: restricting account exist, but balance not enough
	lockFunds = big.NewInt(int64(2E18))
	err = plugin.RestrictingInstance().PledgeLockFunds(addrArr[0], lockFunds, stateDB)
	if err != nil {
		t.Error("The case2 of PledgeLockFunds failed.")
	}

	// case3: restricting account exist, and balance is enough
	lockFunds = big.NewInt(int64(2E18))
	err = plugin.RestrictingInstance().PledgeLockFunds(addrArr[0], lockFunds, stateDB)
	if err != nil {
		t.Error("The case3 of PledgeLockFunds failed.")
	}
}


func TestRestrictingPlugin_ReturnLockFunds(t *testing.T) {

	stateDB := buildStateDB(t)
	buildDbRestrictingPlan(t, stateDB)
	returnFunds := big.NewInt(int64(1E18))

	// case1: restricting account not exist
	notFoundAccount := common.HexToAddress("0x11")
	err := plugin.RestrictingInstance().ReturnLockFunds(notFoundAccount, returnFunds, stateDB)
	if err != nil {
		t.Error("The case1 of ReturnLockFunds failed.")
	}

	// case2: restricting account exist
	err = plugin.RestrictingInstance().ReturnLockFunds(addrArr[0], returnFunds, stateDB)
	if err != nil {
		t.Error("The case2 of ReturnLockFunds failed.")
	}
}

func TestRestrictingPlugin_SlashingNotify(t *testing.T) {

	stateDB := buildStateDB(t)
	buildDbRestrictingPlan(t, stateDB)
	slashingFunds := big.NewInt(int64(1E18))

	// case1: restricting account not exist
	notFoundAccount := common.HexToAddress("0x11")
	err := plugin.RestrictingInstance().SlashingNotify(notFoundAccount, slashingFunds, stateDB)
	if err != nil {
		t.Error("The case1 of SlashingNotify failed.")
	}

	// case2: restricting account exist
	err = plugin.RestrictingInstance().SlashingNotify(addrArr[0], slashingFunds, stateDB)
	if err != nil {
		t.Error("The case2 of SlashingNotify failed.")
	}
}

func TestRestrictingPlugin_GetRestrictingInfo(t *testing.T) {

	stateDB := buildStateDB(t)

	// case1: restricting account not exist
	notFoundAccount := common.HexToAddress("0x11")
	_, err := plugin.RestrictingInstance().GetRestrictingInfo(notFoundAccount, stateDB)
	if err != nil {
		t.Error("The case2 of GetRestrictingInfo failed.")
	}

	// case2: restricting account exist
	if result, err := plugin.RestrictingInstance().GetRestrictingInfo(addrArr[0], stateDB); err != nil {
		t.Errorf("The case2 of GetRestrictingInfo failed.")
	} else {
		t.Log(string(result))
	}
}