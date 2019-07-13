package plugin_test

import (
	"bytes"
	"math/big"
	"reflect"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var (
	errParamPeriodInvalid = common.NewBizError("param epoch invalid")
	errBalanceNotEnough   = common.NewBizError("balance not enough to restrict")
	errAccountNotFound    = common.NewBizError("account is not found")
)

type restrictingInfo struct {
	Balance     *big.Int `json:"Balance"` // Balance representation all locked amount
	Debt        *big.Int `json:"Debt"`    // Debt representation will released amount.
	DebtSymbol  bool     `json:"symbol"`  // Debt is owed to release in the past while symbol is true, else Debt can be used instead of release
	ReleaseList []uint64 `json:"list"`    // ReleaseList representation which epoch will release restricting
}

type releaseAmountInfo struct {
	height uint64    `json:"blockNumber"`   // blockNumber representation of the block number at the released epoch
	amount *big.Int	 `json:"amount"`		// amount representation of the released amount
}


func showRestrictingAccountInfo(t *testing.T, state xcom.StateDB, account common.Address) {
	restrictingKey := restricting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		t.Logf("Restricting account not found, account: %v", account.String())
		return
	}

	var info restrictingInfo
	if err := rlp.Decode(bytes.NewBuffer(bAccInfo), &info); err != nil {
		t.Fatalf("rlp decode info failed, info bytes: %+v", bAccInfo)
	}

	t.Log("Actually balance of restrict account: ", info.Balance)
	t.Log("Actually debt    of restrict account: ", info.Debt)
	t.Log("Actually symbol  of restrict account: ", info.DebtSymbol)
	t.Log("Actually list    of restrict account: ", info.ReleaseList)
}

func showReleaseEpoch(t *testing.T, state xcom.StateDB, epoch uint64) {
	releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
	bAccNumbers := state.GetState(vm.RestrictingContractAddr, releaseEpochKey)

	if len(bAccNumbers) == 0 {
		t.Logf("Release Epoch record not found, epoch: %d", epoch)
		return
	} else {
		t.Logf("Actually account numbers of release epoch %d: %d", epoch, common.BytesToUint32(bAccNumbers))
	}

	for i := uint32(0); i < common.BytesToUint32(bAccNumbers); i++ {

		index := i +1
		releaseAccountKey := restricting.GetReleaseAccountKey(epoch, index)
		bReleaseAcc := state.GetState(vm.RestrictingContractAddr, releaseAccountKey)

		if len(bReleaseAcc) == 0 {
			panic("system error, release account can't empty")
		} else {
			t.Logf("Actually release accounts of epoch %d: %v", epoch, common.BytesToAddress(bReleaseAcc).String())
		}
	}
}

func showReleaseAmount(t *testing.T, state xcom.StateDB, account common.Address, epoch uint64) {
	releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)
	bAmount := state.GetState(account, releaseAmountKey)

	amount := new(big.Int)
	if len(bAmount) == 0 {
		t.Logf("record of restricting account amount not found, account: %v, epoch: %d", account, epoch)
	} else {
		t.Logf("expect release amount of account [%s]: %v", account.String(), amount.SetBytes(bAmount))
	}

	return
}

func TestRestrictingPlugin_EndBlock(t *testing.T) {

	newChainState()
	xcom.SetEconomicModel(&xcom.DefaultConfig)

	stateDB := buildStateDB(t)
	buildDbRestrictingPlan(t, stateDB)

	// case1: blockNumber not arrived settle block height
	head := types.Header{ Number: big.NewInt(1),}
	if err := plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDB); err != nil {
		t.Error("The case1 of EndBlock failed.\n expected err is nil")
		t.Errorf("Actually returns err. blockNumber:%d . errors: %s", head.Number.Uint64(), err.Error())
	} else {

	}

	// case2: blockNumber arrived settle block height, restricting plan not exist
	head = types.Header{Number: big.NewInt(int64(6*xcom.ConsensusSize()*xcom.EpochSize()))}
	if err := plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDB); err != nil {
		t.Error("The case2 of EndBlock failed.\n expected success")
		t.Errorf("Actually returns err. blockNumber:%d . errors: %s", head.Number.Uint64(), err.Error())
	}

	// case3: blockNumber arrived settle block height, restricting plan exist
	head = types.Header{Number: big.NewInt(int64(1*xcom.ConsensusSize()*xcom.EpochSize()))}
	if err := plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDB); err != nil {
		t.Error("The case3 of EndBlock failed.\n expected success")
		t.Errorf("Actually returns err. blockNumber:%d . errors: %s", head.Number.Uint64(), err.Error())
	}
}

func TestRestrictingPlugin_AddRestrictingRecord(t *testing.T) {

	var err error
	var plan restricting.RestrictingPlan

	// case1: balance of sender not enough
	{
		stateDB := buildStateDB(t)
		stateDB.AddBalance(sender, big.NewInt(1))

		var plans = make([]restricting.RestrictingPlan, 5)
		for i := 0; i < 5; i++ {
			v := reflect.ValueOf(&plans[i]).Elem()

			epoch := i+1
			amount := big.NewInt(int64(1E18))
			v.FieldByName("Epoch").SetUint(uint64(epoch))
			v.FieldByName("Amount").Set(reflect.ValueOf(amount))
		}

		err = plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDB)

		t.Logf("expected error is [%s]", errBalanceNotEnough)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errBalanceNotEnough.Error(){
			t.Log("case1 of AddRestrictingRecord pass")
		} else {
			t.Error("case1 of AddRestrictingRecord failed.")
		}
	}

	// case2: account is new user to restricting
	{
		stateDB := buildStateDB(t)

		// preset sender balance
		restrictingAmount :=  big.NewInt(int64(5E18))
		senderBalance := new(big.Int).Add(sender_balance, restrictingAmount)
		stateDB.AddBalance(sender, senderBalance)

		// build input plans for case1
		var plans = make([]restricting.RestrictingPlan, 5)
		for i := 0; i < 5; i++ {
			v := reflect.ValueOf(&plans[i]).Elem()

			epoch := i+1
			amount := big.NewInt(int64(1E18))
			v.FieldByName("Epoch").SetUint(uint64(epoch))
			v.FieldByName("Amount").Set(reflect.ValueOf(amount))
		}

		// Deduct a portion of the money to contract in advance
		stateDB.SubBalance(sender, restrictingAmount)
		stateDB.AddBalance(vm.RestrictingContractAddr, restrictingAmount)

		err := plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDB)

		t.Log("=====================")
		t.Log("expected case2 of AddRestrictingRecord success")
		t.Log("expected balance of sender:", sender_balance)
		t.Log("expected balance of contract:", restrictingAmount)
		t.Log("expected balance of restrict account: ", big.NewInt(int64(5E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{1, 2 ,3, 4, 5})
		for i := 0; i < 5; i++ {
			epoch := i+1
			t.Log("=====================")
			t.Logf("expected account numbers of release epoch %d: 1", epoch)
			t.Logf("expect release accounts of epoch %d: %v", epoch, addrArr[0].String())
			t.Logf("expect release amount of account [%s]: %v", addrArr[0].String(), big.NewInt(int64(1E18)))
		}
		t.Log("=====================")


		if err != nil {
			t.Errorf("case2 of AddRestrictingRecord failed. Actually returns error: %s", err.Error())

		} else {

			t.Log("=====================")
			t.Log("case2 return success!")
			t.Log("Actually balance of sender:", stateDB.GetBalance(sender))
			t.Log("Actually balance of contract:", stateDB.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDB, addrArr[0])
			for i := 0; i < 5; i++ {
				epoch := i+1

				t.Log("=====================")
				showReleaseEpoch(t, stateDB, uint64(epoch))
				showReleaseAmount(t, stateDB, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case2 pass")
		}
	}


	// case3: restricting account exist, but restricting epoch not intersect
	{
		stateDB := buildStateDB(t)

		// preset sender balance
		restrictingAmount := big.NewInt(int64(1E18))
		stateDB.AddBalance(sender, restrictingAmount)

		// build db info
		buildDbRestrictingPlan(t, stateDB)

		// build plans for case3
		var plans = make([]restricting.RestrictingPlan, 1)
		plan.Epoch = uint64(6)
		plan.Amount = restrictingAmount
		plans[0] = plan

		// Deduct a portion of the money to contract in advance
		stateDB.SubBalance(sender, restrictingAmount)
		stateDB.AddBalance(vm.RestrictingContractAddr, restrictingAmount)

		err := plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDB)

		t.Log("=====================")
		t.Log("expected case3 of AddRestrictingRecord success")
		t.Log("expected balance of sender:", sender_balance)
		t.Log("expected balance of contract:", big.NewInt(int64(6E18)))
		t.Log("expected balance of restrict account: ", big.NewInt(int64(6E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{1, 2 ,3, 4, 5, 6})
		for i := 0; i < 6; i++ {
			epoch := i+1
			t.Log("=====================")
			t.Logf("expected account numbers of release epoch %d: 1", epoch)
			t.Logf("expect release accounts of epoch %d: %v", epoch, addrArr[0].String())
			t.Logf("expect release amount of account [%s]: %v", addrArr[0].String(), big.NewInt(int64(1E18)))
		}
		t.Log("=====================")


		if err != nil {
			t.Errorf("case3 of AddRestrictingRecord failed. Actually returns error: %s", err.Error())

		} else {

			t.Log("=====================")
			t.Log("case3 return success!")
			t.Log("Actually balance of sender:", stateDB.GetBalance(sender))
			t.Log("Actually balance of contract:", stateDB.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDB, addrArr[0])
			for i := 0; i < 6; i++ {
				epoch := i+1

				t.Log("=====================")
				showReleaseEpoch(t, stateDB, uint64(epoch))
				showReleaseAmount(t, stateDB, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case3 pass")
		}
	}

	// case4: restricting account exist, and restricting epoch intersect
	{
		stateDB := buildStateDB(t)

		// preset sender balance
		restrictingAmount := big.NewInt(int64(1E18))
		stateDB.AddBalance(sender, restrictingAmount)

		// build db info
		buildDbRestrictingPlan(t, stateDB)

		// build plans for case3
		var plans = make([]restricting.RestrictingPlan, 1)
		plan.Epoch = uint64(5)
		plan.Amount = restrictingAmount
		plans[0] = plan

		// Deduct a portion of the money to contract in advance
		stateDB.SubBalance(sender, restrictingAmount)
		stateDB.AddBalance(vm.RestrictingContractAddr, restrictingAmount)

		err := plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDB)

		t.Log("=====================")
		t.Log("expected case4 of AddRestrictingRecord success")
		t.Log("expected balance of sender:", sender_balance)
		t.Log("expected balance of contract:", big.NewInt(int64(6E18)))
		t.Log("expected balance of restrict account: ", big.NewInt(int64(6E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{1, 2 ,3, 4, 5})
		for i := 0; i < 5; i++ {
			epoch := i+1
			t.Log("=====================")
			t.Logf("expected account numbers of release epoch %d: 1", epoch)
			t.Logf("expect release accounts of epoch %d: %v", epoch, addrArr[0].String())
			if epoch == 5 {
				t.Logf("expect release amount of account [%s]: %v", addrArr[0].String(), big.NewInt(int64(2E18)))
			} else {
				t.Logf("expect release amount of account [%s]: %v", addrArr[0].String(), big.NewInt(int64(1E18)))
			}
		}
		t.Log("=====================")

		if err != nil {
			t.Errorf("case3 of AddRestrictingRecord failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case4 return success!")
			t.Log("Actually balance of sender:", stateDB.GetBalance(sender))
			t.Log("Actually balance of contract:", stateDB.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDB, addrArr[0])
			for i := 0; i < 5; i++ {
				epoch := i+1
				t.Log("=====================")
				showReleaseEpoch(t, stateDB, uint64(epoch))
				showReleaseAmount(t, stateDB, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case4 pass")
		}
	}
}

func TestRestrictingPlugin_PledgeLockFunds(t *testing.T) {

	stateDB := buildStateDB(t)
	buildDbRestrictingPlan(t, stateDB)
	var err error

	// case1: restricting account not exist
	{
		lockFunds := big.NewInt(int64(2E18))
		notFoundAccount := common.HexToAddress("0x11")
		err = plugin.RestrictingInstance().PledgeLockFunds(notFoundAccount, lockFunds, stateDB)
		if err == nil || err.Error() != "account is not found" {
			t.Error("The case1 of PledgeLockFunds failed. expected err is errAccountNotFound")
			t.Errorf("Actually returns [%s]", err.Error())
		}
	}


	// case2: restricting account exist, but Balance not enough
	{
		lockFunds := big.NewInt(int64(2E18))
		err = plugin.RestrictingInstance().PledgeLockFunds(addrArr[0], lockFunds, stateDB)
		if err == nil || err.Error() != "Balance not enough to restrict" {
			t.Error("The case2 of PledgeLockFunds failed. expected err is errBalanceNotEnough")
			t.Errorf("Actually returns [%s]", err.Error())
		}
	}



	// case3: restricting account exist, and Balance is enough
	{
		lockFunds := big.NewInt(int64(2E18))
		err = plugin.RestrictingInstance().PledgeLockFunds(addrArr[0], lockFunds, stateDB)
		if err != nil {
			t.Error("The case3 of PledgeLockFunds failed. expected success")
			t.Errorf("Actually returns err. errors: %s", err.Error())
		}
	}


}


func TestRestrictingPlugin_ReturnLockFunds(t *testing.T) {

	stateDB := buildStateDB(t)
	buildDbRestrictingPlan(t, stateDB)
	returnFunds := big.NewInt(int64(1E18))

	// case1: restricting account not exist
	notFoundAccount := common.HexToAddress("0x11")
	err := plugin.RestrictingInstance().ReturnLockFunds(notFoundAccount, returnFunds, stateDB)
	if err == nil || err.Error() != "account is not found" {
		t.Error("The case1 of ReturnLockFunds failed. expected err is errAccountNotFound")
		t.Errorf("Actually returns [%s]", err.Error())
	}

	// case2: restricting account exist
	err = plugin.RestrictingInstance().ReturnLockFunds(addrArr[0], returnFunds, stateDB)
	if err != nil {
		t.Error("The case2 of ReturnLockFunds failed. expected success")
		t.Errorf("Actually returns err. errors: %s", err.Error())
	}
}

func TestRestrictingPlugin_SlashingNotify(t *testing.T) {

	stateDB := buildStateDB(t)
	buildDbRestrictingPlan(t, stateDB)
	slashingFunds := big.NewInt(int64(1E18))

	// case1: restricting account not exist
	notFoundAccount := common.HexToAddress("0x11")
	err := plugin.RestrictingInstance().SlashingNotify(notFoundAccount, slashingFunds, stateDB)
	if err == nil || err.Error() != "account is not found" {
		t.Error("The case1 of SlashingNotify failed. expected err is errAccountNotFound")
		t.Errorf("Actually returns [%s]", err.Error())
	}

	// case2: restricting account exist
	err = plugin.RestrictingInstance().SlashingNotify(addrArr[0], slashingFunds, stateDB)
	if err != nil {
		t.Error("The case2 of SlashingNotify failed. expected success")
		t.Errorf("Actually returns err. errors: %s", err.Error())
	}
}

func TestRestrictingPlugin_GetRestrictingInfo(t *testing.T) {

	stateDB := buildStateDB(t)

	// case1: restricting account not exist
	notFoundAccount := common.HexToAddress("0x11")
	_, err := plugin.RestrictingInstance().GetRestrictingInfo(notFoundAccount, stateDB)
	if err == nil || err.Error() != "account is not found" {
		t.Error("The case1 of GetRestrictingInfo failed. expected err is errAccountNotFound")
		t.Errorf("Actually returns [%s]", err.Error())
	}

	// case2: restricting account exist
	if result, err := plugin.RestrictingInstance().GetRestrictingInfo(addrArr[0], stateDB); err != nil {
		t.Errorf("The case2 of GetRestrictingInfo failed. expected success")
		t.Errorf("Actually returns err. errors: %s", err.Error())
	} else {
		t.Log(string(result))
	}
}