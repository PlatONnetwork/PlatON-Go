package plugin

import (
	"encoding/json"
	"fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

func TestRestrictingPlugin_EndBlock(t *testing.T) {
	plugin := new(RestrictingPlugin)
	plugin.log = log.Root().New("package", "RestrictingPlugin")

	t.Run("blockChain not arrived settle block height", func(t *testing.T) {
		chain := mock.NewChain(nil)
		buildDbRestrictingPlan(addrArr[0], t, chain.StateDB)
		head := types.Header{Number: big.NewInt(1)}

		err := RestrictingInstance().EndBlock(common.Hash{}, &head, chain.StateDB)

		if err != nil {
			t.Error(err)
			return
		}
		result, err := RestrictingInstance().GetRestrictingInfo(addrArr[0], chain.StateDB)
		if err != nil {
			t.Error(err)
			return
		}
		var res restricting.Result
		if err = json.Unmarshal(result, &res); err != nil {
			t.Fatalf("failed to json decode result, result: %s", result)
		}
		if res.Balance.ToInt().Cmp(big.NewInt(5e18)) != 0 {
			t.Errorf("balance not cmp")
		}
		if res.Debt.ToInt().Cmp(common.Big0) != 0 {
			t.Error("Debt not cmp")
		}
		if len(res.Entry) == 0 {
			t.Error("release entry must not 0")
		}
		var count int = 1
		for _, entry := range res.Entry {
			if entry.Height != uint64(count)*xutil.CalcBlocksEachEpoch() {
				t.Errorf("release block number not  cmp,want %v ,have %v ", uint64(count)*xutil.CalcBlocksEachEpoch(), entry.Height)
			}
			if entry.Amount.ToInt().Cmp(big.NewInt(int64(1e18))) != 0 {
				t.Errorf("release amount  not  cmp,want %v ,have %v ", big.NewInt(int64(1e18)), entry.Amount)
			}
			count++
		}
	})

	t.Run("blockChain arrived settle block height, restricting plan not exist", func(t *testing.T) {
		chain := mock.NewChain(nil)
		blockNumber := uint64(1) * xutil.CalcBlocksEachEpoch()
		head := types.Header{Number: big.NewInt(int64(blockNumber))}
		err := RestrictingInstance().EndBlock(common.Hash{}, &head, chain.StateDB)
		if err != nil {
			t.Error(err)
			return
		}
		if _, err := RestrictingInstance().GetRestrictingInfo(addrArr[0], chain.StateDB); err != errAccountNotFound {
			t.Error("account must not found")
			return
		}
	})
}

func TestRestrictingPlugin_AddRestrictingRecord(t *testing.T) {
	plugin := new(RestrictingPlugin)
	plugin.log = log.Root()
	//	plugin.log.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	from, to := addrArr[0], addrArr[1]

	t.Run("test parameter plans", func(t *testing.T) {
		mockDB := buildStateDB(t)
		mockDB.AddBalance(sender, big.NewInt(1e16))
		type testtmp struct {
			input  []restricting.RestrictingPlan
			expect error
			des    string
		}
		var largePlans, largeMountPlans, notEnough []restricting.RestrictingPlan
		for i := 0; i < 40; i++ {
			largePlans = append(largePlans, restricting.RestrictingPlan{1, big.NewInt(1e15)})
		}
		for i := 0; i < 4; i++ {
			largeMountPlans = append(largeMountPlans, restricting.RestrictingPlan{1, big.NewInt(1e18)})
		}
		for i := 0; i < 4; i++ {
			notEnough = append(notEnough, restricting.RestrictingPlan{1, big.NewInt(1e16)})
		}
		x := []testtmp{
			{
				input:  make([]restricting.RestrictingPlan, 0),
				expect: errCountRestrictPlansInvalid,
				des:    "0 plan",
			},
			{
				input:  nil,
				expect: errCountRestrictPlansInvalid,
				des:    "nil plan",
			},
			{
				input:  []restricting.RestrictingPlan{{0, big.NewInt(1e15)}},
				expect: errParamEpochInvalid,
				des:    "epoch is zero",
			},
			{
				input:  []restricting.RestrictingPlan{{1, big.NewInt(0)}},
				expect: errCreatePlanAmountLessThanZero,
				des:    "amount is 0",
			},
			{
				input:  largePlans,
				expect: errCountRestrictPlansInvalid,
				des:    fmt.Sprintf("must less than %d", restrictTxPlanSize),
			},
			{
				input:  largeMountPlans,
				expect: errBalanceNotEnough,
				des:    "amount not enough",
			},
			{
				input:  notEnough,
				expect: errLockedAmountTooLess,
				des:    "amount too small",
			},
		}
		for _, value := range x {
			if err := plugin.AddRestrictingRecord(sender, addrArr[0], value.input, mockDB); err != value.expect {
				t.Errorf("have %v,want %v", err, value.des)
			}
		}
	})
	t.Run("the record not exist", func(t *testing.T) {
		mockDB := buildStateDB(t)
		mockDB.AddBalance(from, big.NewInt(8e18))
		plans := make([]restricting.RestrictingPlan, 0)
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e17)})
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e17)})
		plans = append(plans, restricting.RestrictingPlan{2, big.NewInt(1e18)})

		if err := plugin.AddRestrictingRecord(from, to, plans, mockDB); err != nil {
			t.Error(err)
		}
		_, rAmount := plugin.getReleaseAmount(mockDB, 1, to)
		assert.Equal(t, big.NewInt(2e17), rAmount)
		_, rAmount2 := plugin.getReleaseAmount(mockDB, 2, to)
		assert.Equal(t, big.NewInt(1e18), rAmount2)

		_, num1 := plugin.getReleaseEpochNumber(mockDB, 1)
		_, num2 := plugin.getReleaseEpochNumber(mockDB, 2)
		_, account1 := plugin.getReleaseAccount(mockDB, 1, num1)

		assert.Equal(t, to, account1)

		_, account2 := plugin.getReleaseAccount(mockDB, 2, num2)
		assert.Equal(t, to, account2)
		res, _ := plugin.getRestrictingInfoToReturn(to, mockDB)
		assert.Equal(t, big.NewInt(1e17+1e17+1e18), res.Balance.ToInt())
		assert.Equal(t, big.NewInt(0), res.Debt.ToInt())

		balance := mockDB.GetBalance(vm.RestrictingContractAddr)
		assert.Equal(t, big.NewInt(1e17+1e17+1e18), balance)
	})

	t.Run("the record  exist,not have NeedRelease", func(t *testing.T) {
		account2 := addrArr[2]
		mockDB := buildStateDB(t)
		mockDB.AddBalance(from, big.NewInt(9e18))
		plugin.storeNumber2ReleaseEpoch(mockDB, restricting.GetReleaseEpochKey(1), 1)
		plugin.storeNumber2ReleaseEpoch(mockDB, restricting.GetReleaseEpochKey(2), 2)
		plugin.storeAmount2ReleaseAmount(mockDB, 1, to, big.NewInt(1e18))
		plugin.storeAmount2ReleaseAmount(mockDB, 2, to, big.NewInt(2e18))
		plugin.storeAmount2ReleaseAmount(mockDB, 2, account2, big.NewInt(1e18))
		plugin.storeAccount2ReleaseAccount(mockDB, 1, 1, to)
		plugin.storeAccount2ReleaseAccount(mockDB, 2, 1, to)
		plugin.storeAccount2ReleaseAccount(mockDB, 2, 2, account2)
		var info, info2 restricting.RestrictingInfo
		info.NeedRelease = big.NewInt(0)
		info.StakingAmount = big.NewInt(1e18)
		info.CachePlanAmount = big.NewInt(1e18 + 2e18)
		info.ReleaseList = []uint64{1, 2}
		plugin.storeRestrictingInfo(mockDB, restricting.GetRestrictingKey(to), info)
		info2.NeedRelease = big.NewInt(0)
		info2.StakingAmount = big.NewInt(1e18)
		info2.CachePlanAmount = big.NewInt(1e18)
		plugin.storeRestrictingInfo(mockDB, restricting.GetRestrictingKey(account2), info2)
		mockDB.AddBalance(vm.RestrictingContractAddr, big.NewInt(2e18))

		plans := make([]restricting.RestrictingPlan, 0)
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e17)})
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e17)})
		plans = append(plans, restricting.RestrictingPlan{2, big.NewInt(1e18)})
		plans = append(plans, restricting.RestrictingPlan{3, big.NewInt(1e18)})
		if err := plugin.AddRestrictingRecord(from, to, plans, mockDB); err != nil {
			t.Error(err)
		}

		_, rAmount := plugin.getReleaseAmount(mockDB, 1, to)
		assert.Equal(t, big.NewInt(1e18+1e17+1e17), rAmount)
		_, rAmount2 := plugin.getReleaseAmount(mockDB, 2, to)
		assert.Equal(t, big.NewInt(1e18+2e18), rAmount2)

		_, account1 := plugin.getReleaseAccount(mockDB, 1, 1)

		assert.Equal(t, to, account1)

		_, account3 := plugin.getReleaseAccount(mockDB, 2, 1)
		assert.Equal(t, to, account3)
		_, info2, err := plugin.mustGetRestrictingInfoByDecode(mockDB, to)
		if err != nil {
			t.Error()
		}
		assert.Equal(t, big.NewInt(3e18+2e17+2e18), info2.CachePlanAmount)
		assert.Equal(t, big.NewInt(1e18), info2.StakingAmount)
		assert.Equal(t, big.NewInt(0), info2.NeedRelease)
		assert.Equal(t, 3, len(info2.ReleaseList))

		balance := mockDB.GetBalance(vm.RestrictingContractAddr)
		assert.Equal(t, big.NewInt(2e18+2e17+2e18), balance)

	})
	t.Run("the record  exist,have NeedRelease", func(t *testing.T) {
		t.Run("the NeedRelease amount is grate or equal than  add plan amount", func(t *testing.T) {
			mockDB := buildStateDB(t)
			mockDB.AddBalance(from, big.NewInt(9e18))
			plugin.storeNumber2ReleaseEpoch(mockDB, restricting.GetReleaseEpochKey(2), 2)
			plugin.storeAmount2ReleaseAmount(mockDB, 2, to, big.NewInt(2e18))
			plugin.storeAccount2ReleaseAccount(mockDB, 2, 1, to)
			var info restricting.RestrictingInfo
			info.NeedRelease = big.NewInt(2e18)
			info.StakingAmount = big.NewInt(4e18)
			info.CachePlanAmount = big.NewInt(4e18)
			info.ReleaseList = []uint64{2}
			plugin.storeRestrictingInfo(mockDB, restricting.GetRestrictingKey(to), info)

			mockDB.AddBalance(to, big.NewInt(1e18))
			mockDB.AddBalance(vm.RestrictingContractAddr, big.NewInt(0))
			mockDB.SetState(vm.RestrictingContractAddr, restricting.GetLatestEpochKey(), common.Uint32ToBytes(1))

			plans := make([]restricting.RestrictingPlan, 0)
			plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e17)})
			plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e17)})
			plans = append(plans, restricting.RestrictingPlan{2, big.NewInt(1e18)})
			if err := plugin.AddRestrictingRecord(from, to, plans, mockDB); err != nil {
				t.Error(err)
			}

			_, rAmount := plugin.getReleaseAmount(mockDB, 2, to)
			assert.Equal(t, big.NewInt(2e18+1e17+1e17), rAmount)
			_, rAmount2 := plugin.getReleaseAmount(mockDB, 3, to)
			assert.Equal(t, big.NewInt(1e18), rAmount2)

			_, account1 := plugin.getReleaseAccount(mockDB, 3, 1)
			assert.Equal(t, to, account1)
			_, account3 := plugin.getReleaseAccount(mockDB, 2, 1)
			assert.Equal(t, to, account3)
			_, info2, err := plugin.mustGetRestrictingInfoByDecode(mockDB, to)
			if err != nil {
				t.Error()
			}
			assert.Equal(t, big.NewInt(4e18), info2.CachePlanAmount)
			assert.Equal(t, big.NewInt(4e18), info2.StakingAmount)
			assert.Equal(t, big.NewInt(2e18-2e17-1e18), info2.NeedRelease)
			assert.Equal(t, 2, len(info2.ReleaseList))

			assert.Equal(t, true, big.NewInt(0).Cmp(mockDB.GetBalance(vm.RestrictingContractAddr)) == 0)
			assert.Equal(t, big.NewInt(1e18+1e18+2e17), mockDB.GetBalance(to))
		})
		t.Run("the NeedRelease amount is less than add plan amount", func(t *testing.T) {
			mockDB := buildStateDB(t)
			mockDB.AddBalance(from, big.NewInt(9e18))
			plugin.storeNumber2ReleaseEpoch(mockDB, restricting.GetReleaseEpochKey(2), 2)
			plugin.storeAmount2ReleaseAmount(mockDB, 2, to, big.NewInt(2e18))
			plugin.storeAccount2ReleaseAccount(mockDB, 2, 1, to)
			var info restricting.RestrictingInfo
			info.NeedRelease = big.NewInt(2e18)
			info.StakingAmount = big.NewInt(4e18)
			info.CachePlanAmount = big.NewInt(4e18)
			info.ReleaseList = []uint64{2}
			plugin.storeRestrictingInfo(mockDB, restricting.GetRestrictingKey(to), info)

			mockDB.AddBalance(to, big.NewInt(1e18))
			mockDB.AddBalance(vm.RestrictingContractAddr, big.NewInt(0))
			mockDB.SetState(vm.RestrictingContractAddr, restricting.GetLatestEpochKey(), common.Uint32ToBytes(1))

			plans := make([]restricting.RestrictingPlan, 0)
			plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(2e18)})
			plans = append(plans, restricting.RestrictingPlan{2, big.NewInt(1e18)})
			if err := plugin.AddRestrictingRecord(from, to, plans, mockDB); err != nil {
				t.Error(err)
			}

			_, rAmount := plugin.getReleaseAmount(mockDB, 2, to)
			assert.Equal(t, big.NewInt(2e18+2e18), rAmount)
			_, rAmount2 := plugin.getReleaseAmount(mockDB, 3, to)
			assert.Equal(t, big.NewInt(1e18), rAmount2)

			_, account1 := plugin.getReleaseAccount(mockDB, 3, 1)
			assert.Equal(t, to, account1)
			_, account3 := plugin.getReleaseAccount(mockDB, 2, 1)
			assert.Equal(t, to, account3)
			_, info2, err := plugin.mustGetRestrictingInfoByDecode(mockDB, to)
			if err != nil {
				t.Error()
			}
			assert.Equal(t, big.NewInt(4e18+3e18-2e18), info2.CachePlanAmount)
			assert.Equal(t, big.NewInt(4e18), info2.StakingAmount)
			assert.Equal(t, big.NewInt(0), info2.NeedRelease)
			assert.Equal(t, 2, len(info2.ReleaseList))

			assert.Equal(t, big.NewInt(1e18), mockDB.GetBalance(vm.RestrictingContractAddr))
			assert.Equal(t, big.NewInt(1e18+2e18), mockDB.GetBalance(to))
		})

	})
}

func TestRestrictingPlugin_GetRestrictingInfo(t *testing.T) {

	t.Run("restricting account not exist", func(t *testing.T) {
		chain := mock.NewChain(nil)
		notFoundAccount := common.HexToAddress("0x11")
		_, err := RestrictingInstance().GetRestrictingInfo(notFoundAccount, chain.StateDB)
		if err != errAccountNotFound {
			t.Errorf("restricting account not exist ,want err %v,have err %v", errAccountNotFound, err)
		}
	})

	t.Run("restricting account exist", func(t *testing.T) {

		chain := mock.NewChain(nil)
		chain.StateDB.AddBalance(addrArr[1], big.NewInt(8e18))

		plans := make([]restricting.RestrictingPlan, 0)
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e18)})
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e18)})
		plans = append(plans, restricting.RestrictingPlan{2, big.NewInt(1e18)})
		total := new(big.Int)
		for _, value := range plans {
			total.Add(total, value.Amount)
		}
		if err := RestrictingInstance().AddRestrictingRecord(addrArr[1], addrArr[0], plans, chain.StateDB); err != nil {
			t.Error(err)
		}

		result, err := RestrictingInstance().GetRestrictingInfo(addrArr[0], chain.StateDB)
		if err != nil {
			t.Errorf("get restrictingInfo fail  error: %s", err.Error())
		}

		var res restricting.Result
		if err = json.Unmarshal(result, &res); err != nil {
			t.Fatalf("failed to json decode result, result: %s", result)
		}

		if res.Balance.ToInt().Cmp(total) != 0 {
			t.Errorf("Balance num is not cmp,should %v have %v", total, res.Balance)
		}
		if res.Debt.ToInt().Cmp(common.Big0) != 0 {
			t.Errorf("Debt num is not cmp,should %v have %v", total, res.Debt)
		}

		if len(res.Entry) != 2 {
			t.Error("wrong num of RestrictingInfo Entry")
		}

		if res.Entry[0].Height != uint64(1)*xutil.CalcBlocksEachEpoch() {
			t.Errorf("release block num is not right,want %v have %v", uint64(1)*xutil.CalcBlocksEachEpoch(), res.Entry[0].Height)
		}
		if res.Entry[0].Amount.ToInt().Cmp(big.NewInt(2e18)) != 0 {
			t.Errorf("release amount not compare ,want %v have %v", big.NewInt(2e18), res.Entry[0].Amount)
		}

		if res.Entry[1].Height != uint64(2)*xutil.CalcBlocksEachEpoch() {
			t.Errorf("release block num is not right,want %v have %v", uint64(2)*xutil.CalcBlocksEachEpoch(), res.Entry[1].Height)
		}
		if res.Entry[1].Amount.ToInt().Cmp(big.NewInt(1e18)) != 0 {
			t.Errorf("release amount not compare ,want %v have %v", big.NewInt(1e18), res.Entry[1].Amount)
		}
	})
}

func TestRestrictingInstance(t *testing.T) {
	mockDB := buildStateDB(t)
	plugin := new(RestrictingPlugin)
	plugin.log = log.Root()
	//	plugin.log.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	from, to := addrArr[0], addrArr[1]
	mockDB.AddBalance(from, big.NewInt(9e18).Add(big.NewInt(9e18), big.NewInt(9e18)))
	plans := make([]restricting.RestrictingPlan, 0)
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(3e18)})
	plans = append(plans, restricting.RestrictingPlan{2, big.NewInt(4e18)})
	plans = append(plans, restricting.RestrictingPlan{3, big.NewInt(2e18)})
	if err := plugin.AddRestrictingRecord(from, to, plans, mockDB); err != nil {
		t.Error(err)
	}
	if err := plugin.releaseRestricting(1, mockDB); err != nil {
		t.Error(err)
	}
	SetLatestEpoch(mockDB, 1)
	if err := plugin.PledgeLockFunds(to, big.NewInt(5e18), mockDB); err != nil {
		t.Error(err)
	}
	if err := plugin.releaseRestricting(2, mockDB); err != nil {
		t.Error(err)
	}
	SetLatestEpoch(mockDB, 2)
	if err := plugin.releaseRestricting(3, mockDB); err != nil {
		t.Error(err)
	}
	SetLatestEpoch(mockDB, 3)
	plans2 := make([]restricting.RestrictingPlan, 0)
	plans2 = append(plans2, restricting.RestrictingPlan{1, big.NewInt(1e18)})
	if err := plugin.AddRestrictingRecord(from, to, plans2, mockDB); err != nil {
		t.Error(err)
	}
	if err := plugin.ReturnLockFunds(to, big.NewInt(5e18), mockDB); err != nil {
		t.Error(err)
	}
	assert.Equal(t, big.NewInt(9e18), mockDB.GetBalance(to))
	assert.Equal(t, big.NewInt(1e18), mockDB.GetBalance(vm.RestrictingContractAddr))

	if err := plugin.releaseRestricting(4, mockDB); err != nil {
		t.Error(err)
	}
	SetLatestEpoch(mockDB, 4)

	assert.Equal(t, big.NewInt(9e18).Add(big.NewInt(9e18), big.NewInt(1e18)), mockDB.GetBalance(to))
	assert.Equal(t, true, mockDB.GetBalance(vm.RestrictingContractAddr).Cmp(big.NewInt(0)) == 0)
	assert.Equal(t, true, mockDB.GetBalance(vm.StakingContractAddr).Cmp(big.NewInt(0)) == 0)
}

func TestRestrictingInstanceWithSlashing(t *testing.T) {
	mockDB := buildStateDB(t)
	plugin := new(RestrictingPlugin)
	plugin.log = log.Root()
	//	plugin.log.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	from, to := addrArr[0], addrArr[1]
	mockDB.AddBalance(from, big.NewInt(9e18).Add(big.NewInt(9e18), big.NewInt(9e18)))
	plans := make([]restricting.RestrictingPlan, 0)
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(3e18)})
	plans = append(plans, restricting.RestrictingPlan{2, big.NewInt(4e18)})
	plans = append(plans, restricting.RestrictingPlan{3, big.NewInt(2e18)})
	if err := plugin.AddRestrictingRecord(from, to, plans, mockDB); err != nil {
		t.Error(err)
	}

	if err := plugin.releaseRestricting(1, mockDB); err != nil {
		t.Error(err)
	}
	SetLatestEpoch(mockDB, 1)

	if err := plugin.PledgeLockFunds(to, big.NewInt(5e18), mockDB); err != nil {
		t.Error(err)
	}

	if err := plugin.releaseRestricting(2, mockDB); err != nil {
		t.Error(err)
	}
	SetLatestEpoch(mockDB, 2)

	if err := plugin.releaseRestricting(3, mockDB); err != nil {
		t.Error(err)
	}
	SetLatestEpoch(mockDB, 3)

	mockDB.SubBalance(vm.StakingContractAddr, big.NewInt(1e18))
	if err := plugin.SlashingNotify(to, big.NewInt(1e18), mockDB); err != nil {
		t.Error(err)
	}

	plans2 := make([]restricting.RestrictingPlan, 0)
	plans2 = append(plans2, restricting.RestrictingPlan{1, big.NewInt(1e18)})
	if err := plugin.AddRestrictingRecord(from, to, plans2, mockDB); err != nil {
		t.Error(err)
	}
	if err := plugin.ReturnLockFunds(to, big.NewInt(4e18), mockDB); err != nil {
		t.Error(err)
	}

	assert.Equal(t, big.NewInt(9e18), mockDB.GetBalance(to))

	if err := plugin.releaseRestricting(4, mockDB); err != nil {
		t.Error(err)
	}
	SetLatestEpoch(mockDB, 4)

	assert.Equal(t, big.NewInt(9e18), mockDB.GetBalance(to))
	if mockDB.GetBalance(vm.RestrictingContractAddr).Cmp(big.NewInt(0)) != 0 {
		t.Error("RestrictingContractAddr should compare", vm.RestrictingContractAddr)
	}
	if mockDB.GetBalance(vm.StakingContractAddr).Cmp(big.NewInt(0)) != 0 {
		t.Error("StakingContractAddr should compare", vm.StakingContractAddr)
	}
	if err := plugin.releaseRestricting(5, mockDB); err != nil {
		t.Error(err)
	}
	SetLatestEpoch(mockDB, 5)

}

func TestRestrictingGetRestrictingInfo(t *testing.T) {
	mockDB := buildStateDB(t)
	plugin := new(RestrictingPlugin)
	plugin.log = log.Root()
	from, to := addrArr[0], addrArr[1]
	mockDB.AddBalance(from, big.NewInt(9e18).Add(big.NewInt(9e18), big.NewInt(9e18)))
	plans := make([]restricting.RestrictingPlan, 0)
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(3e18)})
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(3e18)})

	if err := plugin.AddRestrictingRecord(from, to, plans, mockDB); err != nil {
		t.Error(err)
	}
	res, err := plugin.getRestrictingInfoToReturn(to, mockDB)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.Balance.ToInt(), big.NewInt(6e18))

}

func tmp(des string, plugin *RestrictingPlugin, mockDB xcom.StateDB, to common.Address) {
	_, info, _ := plugin.mustGetRestrictingInfoByDecode(mockDB, to)
	log.Debug("info", "info", info, "des", des)
}
