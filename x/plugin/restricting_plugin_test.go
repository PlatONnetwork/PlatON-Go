package plugin

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
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

// showRestrictingAccountInfo prints restricting info of restricting account in stateDB
func showRestrictingAccountInfo(t *testing.T, state xcom.StateDB, account common.Address) {
	restrictingKey := restricting.GetRestrictingKey(account)
	bAccInfo := state.GetState(vm.RestrictingContractAddr, restrictingKey)

	if len(bAccInfo) == 0 {
		t.Logf("Restricting account not found, account: %v", account.String())
		return
	}

	var info restricting.RestrictingInfo
	if err := rlp.Decode(bytes.NewBuffer(bAccInfo), &info); err != nil {
		t.Fatalf("rlp decode info failed, info bytes: %+v", bAccInfo)
	}

	t.Log("actually balance of restrict account: ", info.Balance)
	t.Log("actually debt    of restrict account: ", info.Debt)
	t.Log("actually symbol  of restrict account: ", info.DebtSymbol)
	t.Log("actually list    of restrict account: ", info.ReleaseList)
}

// showReleaseEpoch prints the number of restricting plan at target epoch in stateDB
func showReleaseEpoch(t *testing.T, state xcom.StateDB, epoch uint64) {
	releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
	bAccNumbers := state.GetState(vm.RestrictingContractAddr, releaseEpochKey)

	if len(bAccNumbers) == 0 {
		t.Logf("release Epoch record not found, epoch: %d", epoch)
		return
	} else {
		t.Logf("actually account numbers of release epoch %d: %d", epoch, common.BytesToUint32(bAccNumbers))
	}

	for i := uint32(0); i < common.BytesToUint32(bAccNumbers); i++ {

		index := i + 1
		releaseAccountKey := restricting.GetReleaseAccountKey(epoch, index)
		bReleaseAcc := state.GetState(vm.RestrictingContractAddr, releaseAccountKey)

		if len(bReleaseAcc) == 0 {
			panic("system error, release account can't empty")
		} else {
			t.Logf("actually release accounts of epoch %d: %v", epoch, common.BytesToAddress(bReleaseAcc).String())
		}
	}
}

// showReleaseAmount prints release amount of the restricting account at target epoch
func showReleaseAmount(t *testing.T, state xcom.StateDB, account common.Address, epoch uint64) {
	releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)
	bAmount := state.GetState(vm.RestrictingContractAddr, releaseAmountKey)

	amount := new(big.Int)
	if len(bAmount) == 0 {
		t.Logf("record of restricting account amount not found, account: %v, epoch: %d", account.String(), epoch)
	} else {
		t.Logf("actually release amount of account [%s]: %v", account.String(), amount.SetBytes(bAmount))
	}
}

func TestRestrictingPlugin_EndBlock(t *testing.T) {

	// case1: blockChain not arrived settle block height
	{
		stateDb := buildStateDB(t)
		buildDbRestrictingPlan(addrArr[0], t, stateDb)
		head := types.Header{Number: big.NewInt(1)}

		err := RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb)

		// show expected result
		t.Log("expected do nothing")

		if err != nil {
			t.Fatalf("The case1 of EndBlock failed. function returns error: %s", err.Error())
		} else {
			t.Logf("actually do nothing")
			t.Log("=====================")
			t.Log("case1 pass")
		}
	}

	// case2: blockChain arrived settle block height, restricting plan not exist
	{
		stateDb := buildStateDB(t)
		blockNumber := uint64(1) * xutil.CalcBlocksEachEpoch()

		head := types.Header{Number: big.NewInt(int64(blockNumber))}
		err := RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb)

		// show expected result
		t.Logf("expected do nothing")

		if err != nil {
			t.Fatalf("The case2 of EndBlock failed. function returns error: %s", err.Error())
		} else {
			t.Logf("actually do nothing")
			t.Log("=====================")
			t.Log("case2 pass")
		}

	}

	// case3: blockChain arrived settle block height, restricting plan exist, debt symbol is false,
	// and debt more or equal than release amount
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		blockNumber := uint64(1) * xutil.CalcBlocksEachEpoch()

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(1E18)
		info.Debt = big.NewInt(2E18)
		info.DebtSymbol = false
		info.ReleaseList = []uint64{1, 2}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		for _, epoch := range info.ReleaseList {
			// store epoch
			releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
			stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

			// store release account
			releaseAccountKey := restricting.GetReleaseAccountKey(epoch, uint32(1))
			stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())
		}

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(1, restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, big.NewInt(1E18).Bytes())
		releaseAmountKey = restricting.GetReleaseAmountKey(2, restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, big.NewInt(2E18).Bytes())

		stateDb.AddBalance(vm.RestrictingContractAddr, big.NewInt(1E18))

		// do EndBlock
		head := types.Header{Number: big.NewInt(int64(blockNumber))}
		err = RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb)

		t.Log("=====================")
		t.Log("expected case3 of EndBlock success")
		t.Log("expected balance of restricting account:", big.NewInt(0))
		t.Log("expected balance of restricting contract account:", big.NewInt(1E18))
		t.Log("expected balance of restrict account: ", big.NewInt(1E18))
		t.Log("expected debt    of restrict account: ", big.NewInt(1E18))
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{2})
		t.Log("=====================")
		t.Log("expected [release Epoch record not found, epoch: 1]")
		t.Log("expected [record of restricting account amount not found, account: 0x740cE31B3fAc20Dac379dB243021A51E80AaDd24, epoch: 1]")
		t.Log("expected account numbers of release epoch 2:", 1)
		t.Logf("expected release accounts of epoch 2: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(2E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case3 of EndBlock failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case3 return success!")
			t.Log("actually balance of restricting account:", stateDb.GetBalance(restrictingAcc))
			t.Log("actually balance of restricting contract account:", stateDb.GetBalance(vm.RestrictingContractAddr))

			showRestrictingAccountInfo(t, stateDb, restrictingAcc)
			for _, epoch := range info.ReleaseList {
				showReleaseEpoch(t, stateDb, epoch)
				showReleaseAmount(t, stateDb, restrictingAcc, epoch)
			}
			t.Log("=====================")
			t.Log("case3 pass")
		}
	}

	// case4: blockChain arrived settle block height, restricting plan exist, debt symbol is false,
	// and total debt and restricting balance is more than release amount
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		blockNumber := uint64(1) * xutil.CalcBlocksEachEpoch()
		SetLatestEpoch(stateDb, 0)

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(2E18)
		info.Debt = big.NewInt(1E18)
		info.DebtSymbol = false
		info.ReleaseList = []uint64{1, 2}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		for _, epoch := range info.ReleaseList {
			// store epoch
			releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
			stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

			// store release account
			releaseAccountKey := restricting.GetReleaseAccountKey(epoch, uint32(1))
			stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())
		}

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(1, restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, big.NewInt(2E18).Bytes())
		releaseAmountKey = restricting.GetReleaseAmountKey(2, restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, big.NewInt(1E18).Bytes())

		stateDb.AddBalance(vm.RestrictingContractAddr, big.NewInt(2E18))

		// do EndBlock
		head := types.Header{Number: big.NewInt(int64(blockNumber))}
		err = RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb)

		t.Log("=====================")
		t.Log("expected case4 of EndBlock success")
		t.Log("expected balance of restricting account:", big.NewInt(1E18))
		t.Log("expected balance of restricting contract account:", big.NewInt(1E18))
		t.Log("expected balance of restrict account: ", big.NewInt(1E18))
		t.Log("expected debt    of restrict account: ", big.NewInt(0))
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{2})
		t.Log("=====================")
		t.Log("expected [release Epoch record not found, epoch: 1]")
		t.Log("expected [record of restricting account amount not found, account: 0x740cE31B3fAc20Dac379dB243021A51E80AaDd24, epoch: 1]")
		t.Log("expected account numbers of release epoch 2:", 1)
		t.Logf("expected release accounts of epoch 2: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(1E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case4 of EndBlock failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case4 return success!")
			t.Log("actually balance of restricting account:", stateDb.GetBalance(restrictingAcc))
			t.Log("actually balance of restricting contract account:", stateDb.GetBalance(vm.RestrictingContractAddr))

			showRestrictingAccountInfo(t, stateDb, restrictingAcc)
			for _, epoch := range info.ReleaseList {
				showReleaseEpoch(t, stateDb, epoch)
				showReleaseAmount(t, stateDb, restrictingAcc, epoch)
			}
			t.Log("=====================")
			t.Log("case4 pass")
		}
	}

	// case5: blockChain arrived settle block height, restricting plan exist, debt symbol is false,
	// and total debt and restricting balance is less than release amount
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		blockNumber := uint64(1) * xutil.CalcBlocksEachEpoch()

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(2E18)
		info.Debt = big.NewInt(1E18)
		info.DebtSymbol = false
		info.ReleaseList = []uint64{1, 2}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		for _, epoch := range info.ReleaseList {
			// store epoch
			releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
			stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

			// store release account
			releaseAccountKey := restricting.GetReleaseAccountKey(epoch, uint32(1))
			stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())
		}

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(1, restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, big.NewInt(4E18).Bytes())
		releaseAmountKey = restricting.GetReleaseAmountKey(2, restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, big.NewInt(1E18).Bytes())

		stateDb.AddBalance(vm.RestrictingContractAddr, big.NewInt(2E18))

		// do EndBlock
		head := types.Header{Number: big.NewInt(int64(blockNumber))}
		err = RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb)

		t.Log("=====================")
		t.Log("expected case5 of EndBlock success")
		t.Log("expected balance of restricting account:", big.NewInt(2E18))
		t.Log("expected balance of restricting contract account:", big.NewInt(0))
		t.Log("expected balance of restrict account: ", big.NewInt(0))
		t.Log("expected debt    of restrict account: ", big.NewInt(1E18))
		t.Log("expected symbol  of restrict account: ", true)
		t.Log("expected list    of restrict account: ", []uint64{2})
		t.Log("=====================")
		t.Log("expected [release Epoch record not found, epoch: 1]")
		t.Log("expected [record of restricting account amount not found, account: 0x740cE31B3fAc20Dac379dB243021A51E80AaDd24, epoch: 1]")
		t.Log("expected account numbers of release epoch 2:", 1)
		t.Logf("expected release accounts of epoch 2: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(1E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case5 of EndBlock failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case5 return success!")
			t.Log("actually balance of restricting account:", stateDb.GetBalance(restrictingAcc))
			t.Log("actually balance of restricting contract account:", stateDb.GetBalance(vm.RestrictingContractAddr))

			showRestrictingAccountInfo(t, stateDb, restrictingAcc)
			for _, epoch := range info.ReleaseList {
				showReleaseEpoch(t, stateDb, epoch)
				showReleaseAmount(t, stateDb, restrictingAcc, epoch)
			}
			t.Log("=====================")
			t.Log("case5 pass")
		}
	}

	// case6: blockChain arrived settle block height, restricting plan exist, debt symbol is true,
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		blockNumber := uint64(2) * xutil.CalcBlocksEachEpoch()
		SetLatestEpoch(stateDb, 1)

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(0)
		info.Debt = big.NewInt(1E18)
		info.DebtSymbol = true
		info.ReleaseList = []uint64{2, 3}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		for _, epoch := range info.ReleaseList {
			// store epoch
			releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
			stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

			// store release account
			releaseAccountKey := restricting.GetReleaseAccountKey(epoch, uint32(1))
			stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())
		}

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(2, restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, big.NewInt(2E18).Bytes())
		releaseAmountKey = restricting.GetReleaseAmountKey(3, restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, big.NewInt(1E18).Bytes())

		// do EndBlock
		head := types.Header{Number: big.NewInt(int64(blockNumber))}
		err = RestrictingInstance().EndBlock(common.Hash{}, &head, stateDb)

		t.Log("=====================")
		t.Log("expected case6 of EndBlock success")
		t.Log("expected balance of restricting account:", big.NewInt(0))
		t.Log("expected balance of restricting contract account:", big.NewInt(0))
		t.Log("expected balance of restrict account: ", big.NewInt(0))
		t.Log("expected debt    of restrict account: ", big.NewInt(3E18))
		t.Log("expected symbol  of restrict account: ", true)
		t.Log("expected list    of restrict account: ", []uint64{3})
		t.Log("=====================")
		t.Log("expected [release Epoch record not found, epoch: 2]")
		t.Log("expected [record of restricting account amount not found, account: 0x740cE31B3fAc20Dac379dB243021A51E80AaDd24, epoch: 2]")
		t.Log("expected account numbers of release epoch 3:", 1)
		t.Logf("expected release accounts of epoch 3: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(1E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case6 of EndBlock failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case6 return success!")
			t.Log("actually balance of restricting account:", stateDb.GetBalance(restrictingAcc))
			t.Log("actually balance of restricting contract account:", stateDb.GetBalance(vm.RestrictingContractAddr))

			showRestrictingAccountInfo(t, stateDb, restrictingAcc)
			for _, epoch := range info.ReleaseList {
				showReleaseEpoch(t, stateDb, epoch)
				showReleaseAmount(t, stateDb, restrictingAcc, epoch)
			}
			t.Log("=====================")
			t.Log("case6 pass")
		}
	}

	/*
	 * path coverage
	 */
	// case7: test release genesis allowance
	{
		stateDb, _, err := newChainState()
		if err != nil {
			t.Fatalf("new a chain state failed, err: %s", err.Error())
		}

		// calculate data after do EndBlock at the first year end
		genesisIssue, ok := new(big.Int).SetString("1000000000000000000000000000", 10)
		if !ok {
			t.Fatal("case7 failed as get genesis failed")
		}

		blocks := xutil.CalcBlocksEachYear()
		epochs := xutil.EpochsPerYear()
		rewardPoolB := stateDb.GetBalance(vm.RewardManagerPoolAddr)
		restrictingContractB := stateDb.GetBalance(vm.RestrictingContractAddr)
		for yearEnd := uint64(1); yearEnd <= 2; yearEnd++ {
			var rate int64
			if yearEnd == 1 {
				rate = 15 // rate of the twice year allowance
			} else {
				rate = 5 // rate of the third year allowance
			}

			allowance := new(big.Int).Mul(genesisIssue, big.NewInt(rate))
			allowance = allowance.Div(allowance, big.NewInt(1000))
			rewardPoolB = new(big.Int).Add(rewardPoolB, allowance)
			restrictingContractB = new(big.Int).Sub(restrictingContractB, allowance)

			// do EndBlock
			SetLatestEpoch(stateDb, epochs*yearEnd-1)
			head := types.Header{Number: big.NewInt(int64(blocks * yearEnd))}
			err = RestrictingInstance().EndBlock(head.Hash(), &head, stateDb)

			t.Log("==============================")
			t.Log("expected case7 of EndBlock success")
			t.Log("expected balance of restricting account:", rewardPoolB)
			t.Log("expected balance of restricting contract account:", restrictingContractB)
			t.Log("expected balance of restrict account: ", restrictingContractB)
			t.Log("expected debt    of restrict account: ", 0)
			t.Log("expected symbol  of restrict account: ", false)
			t.Log("expected list    of restrict account: ", []uint64{epochs * 2})
			t.Log("=====================")
			t.Logf("expected [release Epoch record not found, epoch: %d]", epochs)
			t.Logf("expected [record of restricting account amount not found, account:%s, epoch: %d]", vm.RewardManagerPoolAddr.String(), epochs)

			if yearEnd == 1 {
				t.Logf("expected account numbers of release epoch 105120: %d", 1)
				t.Logf("expected release accounts of epoch 105120: %s", vm.RewardManagerPoolAddr.String())
				t.Logf("expected release amount of account [0x1000000000000000000000000000000000000003]: %v", restrictingContractB)

			} else {
				t.Logf("expected [release Epoch record not found, epoch: %d]", epochs*2)
				t.Logf("expected [record of restricting account amount not found, account:%s, epoch: %d]", vm.RewardManagerPoolAddr.String(), epochs*2)
			}

			if err != nil {
				t.Fatalf("case7 of EndBlock failed. Actually returns error: %s", err.Error())
			} else {
				t.Log("=====================")
				t.Log("case7 return success!")
				t.Log("actually balance of restricting account:", stateDb.GetBalance(vm.RewardManagerPoolAddr))
				t.Log("actually balance of restricting contract account:", stateDb.GetBalance(vm.RestrictingContractAddr))

				showRestrictingAccountInfo(t, stateDb, vm.RewardManagerPoolAddr)
				t.Log("=====================")
				showReleaseEpoch(t, stateDb, epochs)
				showReleaseAmount(t, stateDb, vm.RewardManagerPoolAddr, epochs)
				showReleaseEpoch(t, stateDb, epochs*2)
				showReleaseAmount(t, stateDb, vm.RewardManagerPoolAddr, epochs*2)
			}
		}
	}
}

func TestRestrictingPlugin_AddRestrictingRecord(t *testing.T) {

	var err error
	var plan restricting.RestrictingPlan

	// case1: release epoch must more than zero
	{
		stateDb := buildStateDB(t)
		SetLatestEpoch(stateDb, 5)

		var plans = make([]restricting.RestrictingPlan, 1)
		plans[0].Epoch = 0
		plans[0].Amount = big.NewInt(int64(1E18))

		err = RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errParamEpochInvalid)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errParamEpochInvalid.Error() {
			t.Log("case1 of AddRestrictingRecord pass")
		} else {
			t.Error("case1 of AddRestrictingRecord failed.")
		}
		t.Log("=====================")
		t.Log("case1 pass")
	}

	// case2: balance of sender not enough
	{
		stateDb := buildStateDB(t)
		stateDb.AddBalance(sender, big.NewInt(1))

		var plans = make([]restricting.RestrictingPlan, 1)
		plans[0].Epoch = 1
		plans[0].Amount = big.NewInt(int64(1E18))

		err = RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errBalanceNotEnough)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errBalanceNotEnough.Error() {
			t.Log("case2 of AddRestrictingRecord pass")
		} else {
			t.Error("case2 of AddRestrictingRecord failed.")
		}
		t.Log("=====================")
		t.Log("case2 pass")
	}

	// case3: rlp decode failed
	{
		stateDb := buildStateDB(t)
		stateDb.AddBalance(sender, big.NewInt(1E18))
		restrictingAcc := addrArr[0]

		testData := "this is test data"
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, []byte(testData))

		var plans = make([]restricting.RestrictingPlan, 1)
		plans[0].Epoch = 1
		plans[0].Amount = big.NewInt(int64(1E18))

		err = RestrictingInstance().AddRestrictingRecord(sender, restrictingAcc, plans, stateDb)

		// show expected result
		t.Logf("expecetd error is [rlp: expected input list for restricting.RestrictingInfo]")
		t.Logf("actually error is [%v]", err)

		if err != nil {
			if _, ok := err.(*common.SysError); ok {
				t.Log("case3 of AddRestricting pass")
			} else {
				t.Error("case3 of AddRestrictingRecord failed.")
			}
		} else {
			t.Error("case3 of AddRestrictingRecord failed.")
		}

		t.Log("=====================")
		t.Log("case3 pass")
	}

	// case4: account is new user to restricting
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

		err := RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case4 of AddRestrictingRecord success")
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
			t.Logf("expected release accounts of epoch %d: %v", epoch, addrArr[0].String())
			t.Logf("expected release amount of account [%s]: %v", addrArr[0].String(), big.NewInt(int64(1E18)))
		}
		t.Log("=====================")

		if err != nil {
			t.Errorf("case4 of AddRestrictingRecord failed. Actually returns error: %s", err.Error())

		} else {

			t.Log("=====================")
			t.Log("case4 return success!")
			t.Log("actually balance of sender:", stateDb.GetBalance(sender))
			t.Log("actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
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

	// case5: restricting account exist, but restricting epoch not intersect
	{
		stateDb := buildStateDB(t)

		// preset sender balance
		restrictingAmount := big.NewInt(int64(1E18))
		stateDb.AddBalance(sender, restrictingAmount)

		// build db info
		buildDbRestrictingPlan(addrArr[0], t, stateDb)

		// build plans for case3
		var plans = make([]restricting.RestrictingPlan, 1)
		plan.Epoch = uint64(6)
		plan.Amount = restrictingAmount
		plans[0] = plan

		// Deduct a portion of the money to contract in advance
		stateDb.SubBalance(sender, restrictingAmount)
		stateDb.AddBalance(vm.RestrictingContractAddr, restrictingAmount)

		err := RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case5 of AddRestrictingRecord success")
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
			t.Errorf("case5 of AddRestrictingRecord failed. Actually returns error: %s", err.Error())

		} else {

			t.Log("=====================")
			t.Log("case5 return success!")
			t.Log("actually balance of sender:", stateDb.GetBalance(sender))
			t.Log("actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))

			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			for i := 0; i < 6; i++ {
				epoch := i + 1

				t.Log("=====================")
				showReleaseEpoch(t, stateDb, uint64(epoch))
				showReleaseAmount(t, stateDb, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case5 pass")
		}
	}

	// case6: restricting account exist, and restricting epoch intersect
	{
		stateDb := buildStateDB(t)

		// preset sender balance
		restrictingAmount := big.NewInt(int64(1E18))
		stateDb.AddBalance(sender, restrictingAmount)

		// build db info
		buildDbRestrictingPlan(addrArr[0], t, stateDb)

		// build plans for case3
		var plans = make([]restricting.RestrictingPlan, 1)
		plan.Epoch = uint64(5)
		plan.Amount = restrictingAmount
		plans[0] = plan

		// Deduct a portion of the money to contract in advance
		stateDb.SubBalance(sender, restrictingAmount)
		stateDb.AddBalance(vm.RestrictingContractAddr, restrictingAmount)

		err := RestrictingInstance().AddRestrictingRecord(sender, addrArr[0], plans, stateDb)

		t.Log("=====================")
		t.Log("expected case6 of AddRestrictingRecord success")
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
			t.Errorf("case6 of AddRestrictingRecord failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case6 return success!")
			t.Log("actually balance of sender:", stateDb.GetBalance(sender))
			t.Log("actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			for i := 0; i < 5; i++ {
				epoch := i + 1
				t.Log("=====================")
				showReleaseEpoch(t, stateDb, uint64(epoch))
				showReleaseAmount(t, stateDb, addrArr[0], uint64(epoch))
			}
			t.Log("=====================")
			t.Log("case6 pass")
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

		err = RestrictingInstance().PledgeLockFunds(notFoundAccount, lockFunds, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errAccountNotFound)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errAccountNotFound.Error() {
			t.Log("case1 of PledgeLockFunds pass")
			t.Log("=====================")
			t.Log("case1 pass")
		} else {
			t.Error("case1 of PledgeLockFunds failed.")
		}
	}

	// case2: restricting account exist, but Balance not enough
	{
		stateDb := buildStateDB(t)

		// build data in stateDB for case2
		buildDbRestrictingPlan(addrArr[0], t, stateDb)

		lockFunds := big.NewInt(int64(6E18))

		err = RestrictingInstance().PledgeLockFunds(addrArr[0], lockFunds, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errBalanceNotEnough)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errBalanceNotEnough.Error() {
			t.Log("case2 of PledgeLockFunds pass")
			showRestrictingAccountInfo(t, stateDb, addrArr[0])
			t.Log("=====================")
			t.Log("case2 pass")
		} else {
			t.Error("case2 of PledgeLockFunds failed.")
		}
	}

	// case3: rlp decode failed
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		lockFunds := big.NewInt(1)

		testData := "this is test data"
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, []byte(testData))

		err = RestrictingInstance().PledgeLockFunds(addrArr[0], lockFunds, stateDb)

		// show expected result
		t.Logf("expecetd error is [rlp: expected input list for restricting.RestrictingInfo]")
		t.Logf("actually error is [%v]", err)

		if err != nil {
			if _, ok := err.(*common.SysError); ok {
				t.Log("case3 of PledgeLockFunds pass")
				t.Log("=====================")
				t.Log("case3 pass")
			} else {
				t.Error("case3 of PledgeLockFunds failed.")
			}
		} else {
			t.Error("case3 of PledgeLockFunds failed.")
		}
	}

	// case4: restricting account exist, and Balance is enough
	{
		stateDb := buildStateDB(t)

		// build data in stateDB for case4
		buildDbRestrictingPlan(addrArr[0], t, stateDb)

		lockFunds := big.NewInt(int64(2E18))

		err = RestrictingInstance().PledgeLockFunds(addrArr[0], lockFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case4 of PledgeLockFunds success")
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
			t.Errorf("case4 of PledgeLockFunds failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case4 return success!")
			t.Log("actually balance of contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
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

func TestRestrictingPlugin_ReturnLockFunds(t *testing.T) {

	// case1: restricting account not exist
	{
		stateDb := buildStateDB(t)

		returnFunds := big.NewInt(int64(1E18))
		notFoundAccount := common.HexToAddress("0x11")

		err := RestrictingInstance().ReturnLockFunds(notFoundAccount, returnFunds, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errAccountNotFound)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errAccountNotFound.Error() {
			t.Log("case1 of ReturnLockFunds pass")
			t.Log("=====================")
			t.Log("case1 pass")
		} else {
			t.Error("case1 of ReturnLockFunds failed.")
		}
	}

	// case2: rlp decode failed
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		lockFunds := big.NewInt(1)

		testData := "this is test data"
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, []byte(testData))

		err := RestrictingInstance().ReturnLockFunds(addrArr[0], lockFunds, stateDb)

		// show expected result
		t.Logf("expecetd error is [rlp: expected input list for restricting.RestrictingInfo]")
		t.Logf("actually error is [%v]", err)

		if err != nil {
			if _, ok := err.(*common.SysError); ok {
				t.Log("case2 of ReturnLockFunds pass")
				t.Log("=====================")
				t.Log("case2 pass")
			} else {
				t.Error("case2 of ReturnLockFunds failed.")
			}
		} else {
			t.Error("case2 of ReturnLockFunds failed.")
		}
	}

	// case3: restricting account exist, debt symbol is false
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		returnFunds := big.NewInt(1E18)

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(1E18)
		info.Debt = big.NewInt(0)
		info.DebtSymbol = false
		info.ReleaseList = []uint64{5}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		// store epoch
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(5))
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		// store release account
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(5), uint32(1))
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(5), restrictingAcc)
		amount := big.NewInt(2E18)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, amount.Bytes())

		stateDb.AddBalance(vm.StakingContractAddr, big.NewInt(1E18))
		stateDb.AddBalance(vm.RestrictingContractAddr, big.NewInt(1E18))

		// do ReturnLockFunds
		err = RestrictingInstance().ReturnLockFunds(addrArr[0], returnFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case3 of ReturnLockFunds success")
		t.Log("expected balance of restricting contract:", big.NewInt(int64(2E18)))
		t.Log("expected balance of staking contract:", big.NewInt(0))
		t.Log("expected balance of restrict account: ", big.NewInt(int64(2E18)))
		t.Log("expected debt    of restrict account: ", 0)
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{5})
		t.Log("=====================")
		t.Log("expected account numbers of release epoch 5: 1")
		t.Logf("expected release accounts of epoch 5: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(2E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case3 of ReturnLockFunds failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case3 return success!")
			t.Log("actually balance of restricting contract:", stateDb.GetBalance(vm.RestrictingContractAddr))
			t.Log("actually balance of staking contract:", stateDb.GetBalance(vm.StakingContractAddr))
			showRestrictingAccountInfo(t, stateDb, restrictingAcc)

			showReleaseEpoch(t, stateDb, uint64(5))
			showReleaseAmount(t, stateDb, restrictingAcc, uint64(5))
			t.Log("=====================")
			t.Log("case3 pass")
		}
	}

	// case4: restricting account exist, and debt symbol is true, and amount is less than debt
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		returnFunds := big.NewInt(1E18)

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(0)
		info.Debt = big.NewInt(2E18)
		info.DebtSymbol = true
		info.ReleaseList = []uint64{5}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		// store epoch
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(5))
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		// store release account
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(5), uint32(1))
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(5), restrictingAcc)
		amount := big.NewInt(3E18)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, amount.Bytes())

		stateDb.AddBalance(vm.StakingContractAddr, big.NewInt(1E18))

		// do ReturnLockFunds
		err = RestrictingInstance().ReturnLockFunds(addrArr[0], returnFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case4 of ReturnLockFunds success")
		t.Log("expected balance of restricting account:", big.NewInt(int64(1E18)))
		t.Log("expected balance of staking contract:", big.NewInt(0))
		t.Log("expected balance of restrict account: ", big.NewInt(0))
		t.Log("expected debt    of restrict account: ", big.NewInt(int64(1E18)))
		t.Log("expected symbol  of restrict account: ", true)
		t.Log("expected list    of restrict account: ", []uint64{5})
		t.Log("=====================")
		t.Log("expected account numbers of release epoch 5: 1")
		t.Logf("expected release accounts of epoch 5: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(3E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case4 of ReturnLockFunds failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case4 return success!")
			t.Log("actually balance of restricting account:", stateDb.GetBalance(restrictingAcc))
			t.Log("actually balance of staking contract:", stateDb.GetBalance(vm.StakingContractAddr))
			showRestrictingAccountInfo(t, stateDb, restrictingAcc)

			showReleaseEpoch(t, stateDb, uint64(5))
			showReleaseAmount(t, stateDb, restrictingAcc, uint64(5))
			t.Log("=====================")
			t.Log("case4 pass")
		}
	}

	// case5: restricting account exist, and debt symbol is true, and amount is not less than debt
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		returnFunds := big.NewInt(3E18)

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(0)
		info.Debt = big.NewInt(2E18)
		info.DebtSymbol = true
		info.ReleaseList = []uint64{5}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		// store epoch
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(5))
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		// store release account
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(5), uint32(1))
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(5), restrictingAcc)
		amount := big.NewInt(3E18)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, amount.Bytes())

		stateDb.AddBalance(vm.StakingContractAddr, big.NewInt(3E18))

		// do ReturnLockFunds
		err = RestrictingInstance().ReturnLockFunds(addrArr[0], returnFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case5 of ReturnLockFunds success")
		t.Log("expected balance of restricting account:", big.NewInt(int64(2E18)))
		t.Log("expected balance of staking contract:", big.NewInt(0))
		t.Log("expected balance of restrict account: ", big.NewInt(1E18))
		t.Log("expected debt    of restrict account: ", big.NewInt(int64(0)))
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{5})
		t.Log("=====================")
		t.Log("expected account numbers of release epoch 5: 1")
		t.Logf("expected release accounts of epoch 5: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(3E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case5 of ReturnLockFunds failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case5 return success!")
			t.Log("actually balance of restricting account:", stateDb.GetBalance(restrictingAcc))
			t.Log("actually balance of staking contract:", stateDb.GetBalance(vm.StakingContractAddr))
			showRestrictingAccountInfo(t, stateDb, restrictingAcc)

			showReleaseEpoch(t, stateDb, uint64(5))
			showReleaseAmount(t, stateDb, restrictingAcc, uint64(5))
			t.Log("=====================")
			t.Log("case5 pass")
		}
	}
}

func TestRestrictingPlugin_SlashingNotify(t *testing.T) {

	// case1: restricting account not exist
	{
		stateDb := buildStateDB(t)

		slashingFunds := big.NewInt(int64(1E18))
		notFoundAccount := common.HexToAddress("0x11")

		err := RestrictingInstance().SlashingNotify(notFoundAccount, slashingFunds, stateDb)

		// show expected result
		t.Logf("expected error is [%s]", errAccountNotFound)
		t.Logf("actually error is [%v]", err)

		if err != nil && err.Error() == errAccountNotFound.Error() {
			t.Log("case1 of SlashingNotify pass")
			t.Log("=====================")
			t.Log("case1 pass")
		} else {
			t.Error("case1 of SlashingNotify failed.")
		}
	}

	// case2: rlp decode failed
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		lockFunds := big.NewInt(1)

		testData := "this is test data"
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, []byte(testData))

		err := RestrictingInstance().SlashingNotify(restrictingAcc, lockFunds, stateDb)

		// show expected result
		t.Logf("expecetd error is [rlp: expected input list for restricting.RestrictingInfo]")
		t.Logf("actually error is [%v]", err)

		if err != nil {
			if _, ok := err.(*common.SysError); ok {
				t.Log("case2 of SlashingNotify pass")
				t.Log("=====================")
				t.Log("case2 pass")
			} else {
				t.Error("case2 of SlashingNotify failed.")
			}
		} else {
			t.Error("case2 of SlashingNotify failed.")
		}
	}

	// case3: restricting account exist, and debt symbol is false
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		slashFunds := big.NewInt(1E18)

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(1E18)
		info.Debt = big.NewInt(2E18)
		info.DebtSymbol = false
		info.ReleaseList = []uint64{5}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		// store epoch
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(5))
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		// store release account
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(5), uint32(1))
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(5), restrictingAcc)
		amount := big.NewInt(3E18)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, amount.Bytes())

		stateDb.AddBalance(vm.StakingContractAddr, big.NewInt(1E18))

		// do SlashingNotify
		err = RestrictingInstance().SlashingNotify(addrArr[0], slashFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case3 of SlashingNotify success")
		t.Log("expected balance of restricting account:", big.NewInt(0))
		t.Log("expected balance of restricting contract account:", big.NewInt(0))
		t.Log("expected balance of staking contract:", big.NewInt(1E18))
		t.Log("expected balance of slashing contract account:", big.NewInt(0))
		t.Log("expected balance of restrict account: ", big.NewInt(1E18))
		t.Log("expected debt    of restrict account: ", big.NewInt(3E18))
		t.Log("expected symbol  of restrict account: ", false)
		t.Log("expected list    of restrict account: ", []uint64{5})
		t.Log("=====================")
		t.Log("expected account numbers of release epoch 5: 1")
		t.Logf("expected release accounts of epoch 5: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(3E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case3 of SlashingNotify failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case3 return success!")
			t.Log("actually balance of restricting account:", stateDb.GetBalance(restrictingAcc))
			t.Log("expected balance of restricting contract account:", stateDb.GetBalance(vm.RestrictingContractAddr))
			t.Log("expected balance of staking contract:", stateDb.GetBalance(vm.StakingContractAddr))
			t.Log("expected balance of slashing contract account:", stateDb.GetBalance(vm.SlashingContractAddr))
			showRestrictingAccountInfo(t, stateDb, restrictingAcc)
			showReleaseEpoch(t, stateDb, uint64(5))
			showReleaseAmount(t, stateDb, restrictingAcc, uint64(5))
			t.Log("=====================")
			t.Log("case3 pass")
		}
	}

	// case4: restricting account exist, and debt symbol is true, and amount is less than debt
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		slashFunds := big.NewInt(1E18)

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(1E18)
		info.Debt = big.NewInt(2E18)
		info.DebtSymbol = true
		info.ReleaseList = []uint64{5}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		// store epoch
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(5))
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		// store release account
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(5), uint32(1))
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(5), restrictingAcc)
		amount := big.NewInt(3E18)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, amount.Bytes())

		stateDb.AddBalance(vm.StakingContractAddr, big.NewInt(1E18))

		// do SlashingNotify
		err = RestrictingInstance().SlashingNotify(addrArr[0], slashFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case4 of SlashingNotify success")
		t.Log("expected balance of restricting account:", big.NewInt(0))
		t.Log("expected balance of restricting contract account:", big.NewInt(0))
		t.Log("expected balance of staking contract:", big.NewInt(1E18))
		t.Log("expected balance of slashing contract account:", big.NewInt(0))
		t.Log("expected balance of restrict account: ", big.NewInt(1E18))
		t.Log("expected debt    of restrict account: ", big.NewInt(1E18))
		t.Log("expected symbol  of restrict account: ", true)
		t.Log("expected list    of restrict account: ", []uint64{5})
		t.Log("=====================")
		t.Log("expected account numbers of release epoch 5: 1")
		t.Logf("expected release accounts of epoch 5: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(3E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case4 of SlashingNotify failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case4 return success!")
			t.Log("actually balance of restricting account:", stateDb.GetBalance(restrictingAcc))
			t.Log("expected balance of restricting contract account:", stateDb.GetBalance(vm.RestrictingContractAddr))
			t.Log("expected balance of staking contract:", stateDb.GetBalance(vm.StakingContractAddr))
			t.Log("expected balance of slashing contract account:", stateDb.GetBalance(vm.SlashingContractAddr))
			showRestrictingAccountInfo(t, stateDb, restrictingAcc)
			showReleaseEpoch(t, stateDb, uint64(5))
			showReleaseAmount(t, stateDb, restrictingAcc, uint64(5))
			t.Log("=====================")
			t.Log("case4 pass")
		}
	}

	// case5: restricting account exist, and debt symbol is true, and amount is more than debt
	{
		stateDb := buildStateDB(t)
		restrictingAcc := addrArr[0]
		slashFunds := big.NewInt(3E18)

		var info restricting.RestrictingInfo
		info.Balance = big.NewInt(1E18)
		info.Debt = big.NewInt(2E18)
		info.DebtSymbol = true
		info.ReleaseList = []uint64{5}

		bInfo, err := rlp.EncodeToBytes(info)
		if err != nil {
			t.Fatal("rlp encode test data failed")
		}

		// store restricting info
		restrictingKey := restricting.GetRestrictingKey(restrictingAcc)
		stateDb.SetState(vm.RestrictingContractAddr, restrictingKey, bInfo)

		// store epoch
		releaseEpochKey := restricting.GetReleaseEpochKey(uint64(5))
		stateDb.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(1))

		// store release account
		releaseAccountKey := restricting.GetReleaseAccountKey(uint64(5), uint32(1))
		stateDb.SetState(vm.RestrictingContractAddr, releaseAccountKey, restrictingAcc.Bytes())

		// store release amount
		releaseAmountKey := restricting.GetReleaseAmountKey(uint64(5), restrictingAcc)
		amount := big.NewInt(3E18)
		stateDb.SetState(vm.RestrictingContractAddr, releaseAmountKey, amount.Bytes())

		stateDb.AddBalance(vm.StakingContractAddr, big.NewInt(1E18))

		// do SlashingNotify
		err = RestrictingInstance().SlashingNotify(addrArr[0], slashFunds, stateDb)

		// show expected result
		t.Log("=====================")
		t.Log("expected case5 of SlashingNotify success")
		t.Log("expected balance of restricting account:", big.NewInt(0))
		t.Log("expected balance of restricting contract account:", big.NewInt(0))
		t.Log("expected balance of staking contract:", big.NewInt(1E18))
		t.Log("expected balance of slashing contract account:", big.NewInt(0))
		t.Log("expected balance of restrict account: ", big.NewInt(1E18))
		t.Log("expected debt    of restrict account: ", big.NewInt(1E18))
		t.Log("expected symbol  of restrict account: ", true)
		t.Log("expected list    of restrict account: ", []uint64{5})
		t.Log("=====================")
		t.Log("expected account numbers of release epoch 5: 1")
		t.Logf("expected release accounts of epoch 5: %s", restrictingAcc.String())
		t.Logf("expected release amount of account [%s]: %v", restrictingAcc.String(), big.NewInt(3E18))
		t.Log("=====================")

		if err != nil {
			t.Errorf("case5 of SlashingNotify failed. Actually returns error: %s", err.Error())
		} else {
			t.Log("=====================")
			t.Log("case5 return success!")
			t.Log("actually balance of restricting account:", stateDb.GetBalance(restrictingAcc))
			t.Log("expected balance of restricting contract account:", stateDb.GetBalance(vm.RestrictingContractAddr))
			t.Log("expected balance of staking contract:", stateDb.GetBalance(vm.StakingContractAddr))
			t.Log("expected balance of slashing contract account:", stateDb.GetBalance(vm.SlashingContractAddr))
			showRestrictingAccountInfo(t, stateDb, restrictingAcc)
			showReleaseEpoch(t, stateDb, uint64(5))
			showReleaseAmount(t, stateDb, restrictingAcc, uint64(5))
			t.Log("=====================")
			t.Log("case5 pass")
		}
	}
}

func TestRestrictingPlugin_GetRestrictingInfo(t *testing.T) {

	t.Run("restricting account not exist", func(t *testing.T) {
		stateDb := buildStateDB(t)
		notFoundAccount := common.HexToAddress("0x11")
		_, err := RestrictingInstance().GetRestrictingInfo(notFoundAccount, stateDb)
		if err != errAccountNotFound {
			t.Errorf("restricting account not exist ,want err %v,have err %v", errAccountNotFound, err)
		}
	})

	t.Run("restricting account exist", func(t *testing.T) {

		stateDb := buildStateDB(t)
		buildDbRestrictingPlan(addrArr[0], t, stateDb)
		result, err := RestrictingInstance().GetRestrictingInfo(addrArr[0], stateDb)
		if err != nil {
			t.Errorf("case2 of GetRestrictingInfo failed. Actually returns error: %s", err.Error())
		}
		t.Log("expected case2 of GetRestrictingInfo success")
		t.Log("expected balance of restrict account: ", big.NewInt(int64(5E18)))
		t.Log("expected slash   of restrict account: ", big.NewInt(0))
		t.Log("expected debt    of restrict account: ", big.NewInt(0))
		t.Log("expected staking of restrict account: ", big.NewInt(0))
		for i := 0; i < 5; i++ {
			expectedBlocks := uint64(i+1) * xutil.CalcBlocksEachEpoch()
			t.Logf("expected release amount at blockNumber [%d] is: %v", expectedBlocks, big.NewInt(int64(1E18)))
		}

		if len(result) == 0 {
			t.Log("case2 of GetRestrictingInfo failed. Actually result is empty")
		}

		var res restricting.Result
		if err = json.Unmarshal(result, &res); err != nil {
			t.Fatalf("failed to elp decode result, result: %s", result)
		}

		t.Log("actually balance of restrict account: ", res.Balance)
		t.Log("actually debt    of restrict account: ", res.Debt)
		t.Log("actually symbol  of restrict account: ", res.Symbol)

		for _, info := range res.Entry {
			t.Logf("actually release amount at blockNumber [%d] is: %v", info.Height, info.Amount)
		}
	})

	// case3: get genesis restricting info
	{
		stateDb, _, _ := newChainState()
		result, err := RestrictingInstance().GetRestrictingInfo(vm.RewardManagerPoolAddr, stateDb)

		// show expected result
		t.Logf("expected result is [123 34 98 97 108 97 110 99 101 34 58 50 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 44 34 100 101 98 116 34 58 48 44 34 115 121 109 98 111 108 34 58 102 97 108 115 101 44 34 69 110 116 114 121 34 58 91 123 34 98 108 111 99 107 78 117 109 98 101 114 34 58 51 49 53 51 54 48 48 48 44 34 97 109 111 117 110 116 34 58 49 56 51 55 53 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 125 44 123 34 98 108 111 99 107 78 117 109 98 101 114 34 58 54 51 48 55 50 48 48 48 44 34 97 109 111 117 110 116 34 58 54 49 50 53 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 48 125 93 125]")
		t.Logf("actually result is %v", result)

		var res restricting.Result
		if err != nil {
			t.Error(err.Error())
		} else {
			if err := json.Unmarshal(result, &res); err != nil {
				t.Error(err.Error())
			} else {
				t.Log(res)
				t.Log("=====================")
				t.Log("case3 pass")
			}
		}
	}
}
