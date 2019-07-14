package plugin_test

import (
	"bytes"
	"encoding/json"
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
	errBalanceNotEnough = common.NewBizError("balance not enough to restrict")
	errAccountNotFound  = common.NewBizError("account is not found")
)

func showRestrictingAccountInfo(t *testing.T, state xcom.StateDB, account common.Address) {
	restrictingKey := restricting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		t.Logf("Restricting account not found, account: %v", account.String())
		return
	}

	var info restricting.RestrictingInfo
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

		index := i + 1
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
		t.Logf("record of restricting account amount not found, account: %v, epoch: %d", account.String(), epoch)
	} else {
		t.Logf("expect release amount of account [%s]: %v", account.String(), amount.SetBytes(bAmount))
	}

	return
}

func TestRestrictingPlugin_EndBlock(t *testing.T) {
	var err error

	// case1: blockNumber not arrived settle block height
	{
		stateDb := buildStateDB(t)
		buildDbRestrictingPlan(t, stateDb)
		head := types.Header{Number: big.NewInt(1)}

		err = plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb)

		t.Logf("expected do nothing")

		if err != nil {
			t.Error("The case1 of EndBlock failed. expected err is nil")
		}
	}

	// case2: blockNumber arrived settle block height, restricting plan not exist
	{
		stateDb := buildStateDB(t)

		plugin.SetLatestEpoch(stateDb, 0)

		head := types.Header{Number: big.NewInt(int64(1 * xcom.ConsensusSize() * xcom.EpochSize()))}
		err = plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb)

		// show expected result
		t.Log("expected case2 of EndBlock success.")

		if err == nil {
			t.Log("Actually case2 of EndBlock success.")
			t.Log("case2 pass")
		} else {
			t.Error("case2 of EndBlock failed.")
		}
	}

	// case3: blockNumber arrived settle block height, restricting plan exist
	{
		stateDb := buildStateDB(t)

		plugin.SetLatestEpoch(stateDb, 0)
		buildDbRestrictingPlan(t, stateDb)
		head := types.Header{Number: big.NewInt(int64(1 * xcom.ConsensusSize() * xcom.EpochSize()))}

		err = plugin.RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb)

		t.Log("=====================")
		t.Log("expected case3 of EndBlock success.")
		t.Log("expected balance of restricting account:", big.NewInt(int64(1E18)))
		t.Log("expected balance of contract:", big.NewInt(int64(4E18)))
		t.Log("expected balance of restrict account: ", big.NewInt(int64(4E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{2, 3, 4, 5})
		for i := 1; i < 5; i++ {
			epoch := i + 1
			t.Log("=====================")
			t.Logf("expected account numbers of release epoch %d: 1", epoch)
			t.Logf("expect release accounts of epoch %d: %v", epoch, addrArr[0].String())
			t.Logf("expect release amount of account [%s]: %v", addrArr[0].String(), big.NewInt(int64(1E18)))
		}
		t.Log("=====================")

		if err != nil {
			t.Errorf("case3 of EndBlock failed. Actually returns error: %s", err.Error())

		} else {

			t.Log("=====================")
			t.Log("case3 return success!")
			t.Log("expected balance of restricting account:", stateDb.GetBalance(addrArr[0]))
			t.Log("Actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			for i := 0; i < 5; i++ {
				epoch := i + 1

				t.Log("=====================")
				showReleaseEpoch(t, stateDb, uint64(epoch))
				showReleaseAmount(t, stateDb, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case3 pass")
		}
	}
}

func TestRestrictingPlugin_AddRestrictingRecord(t *testing.T) {

	var err error
	var plan restricting.RestrictingPlan

	// case1: balance of sender not enough
	{
		stateDb := buildStateDB(t)
		stateDb.AddBalance(sender, big.NewInt(1))

		var plans = make([]restricting.RestrictingPlan, 5)
		for i := 0; i < 5; i++ {
			v := reflect.ValueOf(&plans[i]).Elem()

			epoch := i + 1
			amount := big.NewInt(int64(1E18))
			v.FieldByName("Epoch").SetUint(uint64(epoch))
			v.FieldByName("Amount").Set(reflect.ValueOf(amount))
		}

		err = plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errBalanceNotEnough)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errBalanceNotEnough.Error() {
			t.Log("case1 of AddRestrictingRecord pass")
		} else {
			t.Error("case1 of AddRestrictingRecord failed.")
		}
	}

	// case2: account is new user to restricting
	{
		stateDb := buildStateDB(t)

		// preset sender balance
		restrictingAmount := big.NewInt(int64(5E18))
		senderBalance := new(big.Int).Add(sender_balance, restrictingAmount)
		stateDb.AddBalance(sender, senderBalance)

		// build input plans for case1
		var plans = make([]restricting.RestrictingPlan, 5)
		for i := 0; i < 5; i++ {
			v := reflect.ValueOf(&plans[i]).Elem()

			epoch := i + 1
			amount := big.NewInt(int64(1E18))
			v.FieldByName("Epoch").SetUint(uint64(epoch))
			v.FieldByName("Amount").Set(reflect.ValueOf(amount))
		}

		// Deduct a portion of the money to contract in advance
		stateDb.SubBalance(sender, restrictingAmount)
		stateDb.AddBalance(vm.RestrictingContractAddr, restrictingAmount)

		err := plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case2 of AddRestrictingRecord success")
		t.Log("expected balance of sender:", sender_balance)
		t.Log("expected balance of contract:", restrictingAmount)
		t.Log("expected balance of restrict account: ", big.NewInt(int64(5E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{1, 2, 3, 4, 5})
		for i := 0; i < 5; i++ {
			epoch := i + 1
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
			t.Log("Actually balance of sender:", stateDb.GetBalance(sender))
			t.Log("Actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			for i := 0; i < 5; i++ {
				epoch := i + 1

				t.Log("=====================")
				showReleaseEpoch(t, stateDb, uint64(epoch))
				showReleaseAmount(t, stateDb, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case2 pass")
		}
	}

	// case3: restricting account exist, but restricting epoch not intersect
	{
		stateDb := buildStateDB(t)

		// preset sender balance
		restrictingAmount := big.NewInt(int64(1E18))
		stateDb.AddBalance(sender, restrictingAmount)

		// build db info
		buildDbRestrictingPlan(t, stateDb)

		// build plans for case3
		var plans = make([]restricting.RestrictingPlan, 1)
		plan.Epoch = uint64(6)
		plan.Amount = restrictingAmount
		plans[0] = plan

		// Deduct a portion of the money to contract in advance
		stateDb.SubBalance(sender, restrictingAmount)
		stateDb.AddBalance(vm.RestrictingContractAddr, restrictingAmount)

		err := plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case3 of AddRestrictingRecord success")
		t.Log("expected balance of sender:", sender_balance)
		t.Log("expected balance of contract:", big.NewInt(int64(6E18)))
		t.Log("expected balance of restrict account: ", big.NewInt(int64(6E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{1, 2, 3, 4, 5, 6})
		for i := 0; i < 6; i++ {
			epoch := i + 1
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
			t.Log("Actually balance of sender:", stateDb.GetBalance(sender))
			t.Log("Actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))

			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			for i := 0; i < 6; i++ {
				epoch := i + 1

				t.Log("=====================")
				showReleaseEpoch(t, stateDb, uint64(epoch))
				showReleaseAmount(t, stateDb, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case3 pass")
		}
	}

	// case4: restricting account exist, and restricting epoch intersect
	{
		stateDb := buildStateDB(t)

		// preset sender balance
		restrictingAmount := big.NewInt(int64(1E18))
		stateDb.AddBalance(sender, restrictingAmount)

		// build db info
		buildDbRestrictingPlan(t, stateDb)

		// build plans for case3
		var plans = make([]restricting.RestrictingPlan, 1)
		plan.Epoch = uint64(5)
		plan.Amount = restrictingAmount
		plans[0] = plan

		// Deduct a portion of the money to contract in advance
		stateDb.SubBalance(sender, restrictingAmount)
		stateDb.AddBalance(vm.RestrictingContractAddr, restrictingAmount)

		err := plugin.RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDb)

		t.Log("=====================")
		t.Log("expected case4 of AddRestrictingRecord success")
		t.Log("expected balance of sender:", sender_balance)
		t.Log("expected balance of contract:", big.NewInt(int64(6E18)))
		t.Log("expected balance of restrict account: ", big.NewInt(int64(6E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{1, 2, 3, 4, 5})
		for i := 0; i < 5; i++ {
			epoch := i + 1
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
			t.Errorf("case4 of AddRestrictingRecord failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case4 return success!")
			t.Log("Actually balance of sender:", stateDb.GetBalance(sender))
			t.Log("Actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			for i := 0; i < 5; i++ {
				epoch := i + 1
				t.Log("=====================")
				showReleaseEpoch(t, stateDb, uint64(epoch))
				showReleaseAmount(t, stateDb, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case4 pass")
		}
	}
}

func TestRestrictingPlugin_PledgeLockFunds(t *testing.T) {

	var err error

	// case1: restricting account not exist
	{
		stateDb := buildStateDB(t)

		lockFunds := big.NewInt(int64(2E18))
		notFoundAccount := common.HexToAddress("0x11")

		err = plugin.RestrictingInstance().PledgeLockFunds(notFoundAccount, lockFunds, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errAccountNotFound)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errAccountNotFound.Error() {
			t.Log("case1 of PledgeLockFunds pass")
		} else {
			t.Error("case1 of PledgeLockFunds failed.")
		}
	}

	// case2: restricting account exist, but Balance not enough
	{
		stateDb := buildStateDB(t)

		// build data in stateDB for case2
		buildDbRestrictingPlan(t, stateDb)

		lockFunds := big.NewInt(int64(6E18))

		err = plugin.RestrictingInstance().PledgeLockFunds(addrArr[0], lockFunds, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errBalanceNotEnough)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errBalanceNotEnough.Error() {
			t.Log("case2 of PledgeLockFunds pass")
			showRestrictingAccountInfo(t, stateDb, addrArr[0])
		} else {
			t.Error("case2 of PledgeLockFunds failed.")
		}
	}

	// case3: restricting account exist, and Balance is enough
	{
		stateDb := buildStateDB(t)

		// build data in stateDB for case3
		buildDbRestrictingPlan(t, stateDb)

		lockFunds := big.NewInt(int64(2E18))

		err = plugin.RestrictingInstance().PledgeLockFunds(addrArr[0], lockFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case3 of PledgeLockFunds success")
		t.Log("expected balance of contract:", big.NewInt(int64(3E18)))
		t.Log("expected balance of restrict account: ", big.NewInt(int64(3E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{1, 2, 3, 4, 5})
		for i := 0; i < 5; i++ {
			epoch := i + 1
			t.Log("=====================")
			t.Logf("expected account numbers of release epoch %d: 1", epoch)
			t.Logf("expected release accounts of epoch %d: %v", epoch, addrArr[0].String())
			t.Logf("expected release amount of account [%s]: %v", addrArr[0].String(), big.NewInt(int64(1E18)))
		}
		t.Log("=====================")

		if err != nil {
			t.Errorf("case3 of PledgeLockFunds failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case3 return success!")
			t.Log("Actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			for i := 0; i < 5; i++ {
				epoch := i + 1
				t.Log("=====================")
				showReleaseEpoch(t, stateDb, uint64(epoch))
				showReleaseAmount(t, stateDb, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case3 pass")
		}
	}
}

func TestRestrictingPlugin_ReturnLockFunds(t *testing.T) {

	var err error

	// case1: restricting account not exist
	{
		stateDb := buildStateDB(t)

		returnFunds := big.NewInt(int64(1E18))
		notFoundAccount := common.HexToAddress("0x11")

		err = plugin.RestrictingInstance().PledgeLockFunds(notFoundAccount, returnFunds, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errAccountNotFound)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errAccountNotFound.Error() {
			t.Log("case1 of ReturnLockFunds pass")
		} else {
			t.Error("case1 of ReturnLockFunds failed.")
		}
	}

	// case2: restricting account exist
	{
		stateDb := buildStateDB(t)

		returnFunds := big.NewInt(int64(1E18))
		stateDb.AddBalance(vm.StakingContractAddr, returnFunds)

		// build date of restricting account for case2
		buildDBStakingRestrictingFunds(t, stateDb)

		err = plugin.RestrictingInstance().ReturnLockFunds(addrArr[0], returnFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case2 of ReturnLockFunds success")
		t.Log("expected balance of staking contract:", big.NewInt(0))
		t.Log("expected balance of restricting contract:", big.NewInt(int64(2E18)))
		t.Log("expected balance of restrict account: ", big.NewInt(int64(2E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{1, 2, 3, 4, 5})
		for i := 0; i < 5; i++ {
			epoch := i + 1
			t.Log("=====================")
			t.Logf("expected account numbers of release epoch %d: 1", epoch)
			t.Logf("expected release accounts of epoch %d: %v", epoch, addrArr[0].String())
			t.Logf("expected release amount of account [%s]: %v", addrArr[0].String(), big.NewInt(int64(1E18)))
		}
		t.Log("=====================")

		if err != nil {
			t.Errorf("case2 of ReturnLockFunds failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case2 return success!")
			t.Log("Actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			for i := 0; i < 5; i++ {
				epoch := i + 1
				t.Log("=====================")
				showReleaseEpoch(t, stateDb, uint64(epoch))
				showReleaseAmount(t, stateDb, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case2 pass")
		}
	}
}

func TestRestrictingPlugin_SlashingNotify(t *testing.T) {

	var err error

	// case1: restricting account not exist
	{
		stateDb := buildStateDB(t)

		slashingFunds := big.NewInt(int64(1E18))
		notFoundAccount := common.HexToAddress("0x11")

		err = plugin.RestrictingInstance().SlashingNotify(notFoundAccount, slashingFunds, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errAccountNotFound)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errAccountNotFound.Error() {
			t.Log("case1 of SlashingNotify pass")
		} else {
			t.Error("case1 of SlashingNotify failed.")
		}
	}

	// case2: restricting account exist
	{
		stateDb := buildStateDB(t)

		slashingFunds := big.NewInt(int64(1E18))

		// build date of restricting account for case2
		buildDbRestrictingPlan(t, stateDb)

		// preset stating condition
		pledgeFunds := big.NewInt(int64(2E18))
		err = plugin.RestrictingInstance().PledgeLockFunds(addrArr[0], pledgeFunds, stateDb)
		if err != nil {
			t.Fatalf("failed to preset for SlashingNotify")
		}

		err = plugin.RestrictingInstance().SlashingNotify(addrArr[0], slashingFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case2 of SlashingNotify success")
		t.Log("expected balance of restricting contract:", big.NewInt(int64(3E18)))
		t.Log("expected balance of restrict account: ", big.NewInt(int64(3E18)))
		t.Log("expected debt    of restrict account: ", big.NewInt(int64(1E18)))
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{1, 2, 3, 4, 5})
		for i := 0; i < 5; i++ {
			epoch := i + 1
			t.Log("=====================")
			t.Logf("expected account numbers of release epoch %d: 1", epoch)
			t.Logf("expected release accounts of epoch %d: %v", epoch, addrArr[0].String())
			t.Logf("expected release amount of account [%s]: %v", addrArr[0].String(), big.NewInt(int64(1E18)))
		}
		t.Log("=====================")

		if err != nil {
			t.Errorf("case2 of SlashingNotify failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case2 return success!")
			t.Log("Actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			for i := 0; i < 5; i++ {
				epoch := i + 1
				t.Log("=====================")
				showReleaseEpoch(t, stateDb, uint64(epoch))
				showReleaseAmount(t, stateDb, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case2 pass")
		}

	}
}

func TestRestrictingPlugin_GetRestrictingInfo(t *testing.T) {

	// case1: restricting account not exist
	{
		stateDb := buildStateDB(t)

		notFoundAccount := common.HexToAddress("0x11")
		_, err := plugin.RestrictingInstance().GetRestrictingInfo(notFoundAccount, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errAccountNotFound)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errAccountNotFound.Error() {
			t.Log("case1 of GetRestrictingInfo pass")
		} else {
			t.Error("case1 of GetRestrictingInfo failed.")
		}
	}

	// case2: restricting account exist
	{
		stateDb := buildStateDB(t)

		buildDbRestrictingPlan(t, stateDb)

		result, err := plugin.RestrictingInstance().GetRestrictingInfo(addrArr[0], stateDb)

		t.Log("=====================")
		t.Log("expected case2 of GetRestrictingInfo success")
		t.Log("expected balance of restrict account: ", big.NewInt(int64(5E18)))
		t.Log("expected slash   of restrict account: ", big.NewInt(0))
		t.Log("expected debt    of restrict account: ", big.NewInt(0))
		t.Log("expected staking of restrict account: ", big.NewInt(0))
		for i := 0; i < 5; i++ {
			expectedBlocks := uint64(i+1) * xcom.EpochSize() * xcom.ConsensusSize()
			t.Logf("expected release amount at blockNumber [%d] is: %v", expectedBlocks, big.NewInt(int64(1E18)))
		}
		t.Log("=====================")

		if err != nil {
			t.Errorf("case2 of GetRestrictingInfo failed. Actually returns error: %s", err.Error())
		} else {

			if len(result) == 0 {
				t.Log("case2 of GetRestrictingInfo failed. Actually result is empty")
			}

			var res restricting.Result
			if err = rlp.Decode(bytes.NewBuffer(result), &res); err != nil {
				t.Fatalf("failed to elp decode result, result: %s", result)
			}

			t.Log("Actually balance of restrict account: ", res.Balance)
			t.Log("Actually debt    of restrict account: ", res.Debt)
			t.Log("Actually slash   of restrict account: ", res.Slash)
			t.Log("Actually staking of restrict account: ", res.Staking)

			var infos []restricting.ReleaseAmountInfo
			if err = json.Unmarshal(res.Entry, &infos); err != nil {
				t.Fatalf("unmarshal release amout info failed, err:%s", err.Error())
			}

			for _, info := range infos {
				t.Logf("Actually release amount at blockNumber [%d] is: %v", info.Height, info.Amount)
			}
		}
	}
}
