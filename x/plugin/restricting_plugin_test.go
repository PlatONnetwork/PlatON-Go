// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package plugin

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

var RestrictingTxHash = common.HexToHash("abc")

func TestRestrictingPlugin_EndBlock(t *testing.T) {
	plugin := new(RestrictingPlugin)
	plugin.log = log.Root().New("package", "RestrictingPlugin")

	t.Run("blockChain not arrived settle block height", func(t *testing.T) {
		chain := mock.NewChain()
		buildDbRestrictingPlan(addrArr[0], t, chain.StateDB)
		head := types.Header{Number: big.NewInt(1)}

		if err := RestrictingInstance().EndBlock(common.Hash{}, &head, chain.StateDB); err != nil {
			t.Error(err)
			return
		}
		res, err := RestrictingInstance().GetRestrictingInfo(addrArr[0], chain.StateDB)
		if err != nil {
			t.Error(err)
			return
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
		chain := mock.NewChain()
		blockNumber := uint64(1) * xutil.CalcBlocksEachEpoch()
		head := types.Header{Number: big.NewInt(int64(blockNumber))}
		err := RestrictingInstance().EndBlock(common.Hash{}, &head, chain.StateDB)
		if err != nil {
			t.Error(err)
			return
		}
		if _, err := RestrictingInstance().GetRestrictingInfo(addrArr[0], chain.StateDB); err != restricting.ErrAccountNotFound {
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
	sdb := snapshotdb.Instance()
	defer sdb.Clear()
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{"", new(big.Int).SetInt64(0).String(), 0})
	if err := sdb.PutBaseDB(key, value); nil != err {
		t.Error(err)
		return
	}

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
				expect: restricting.ErrCountRestrictPlansInvalid,
				des:    "0 plan",
			},
			{
				input:  nil,
				expect: restricting.ErrCountRestrictPlansInvalid,
				des:    "nil plan",
			},
			{
				input:  []restricting.RestrictingPlan{{0, big.NewInt(1e15)}},
				expect: restricting.ErrParamEpochInvalid,
				des:    "epoch is zero",
			},
			{
				input:  []restricting.RestrictingPlan{{1, big.NewInt(0)}},
				expect: restricting.ErrCreatePlanAmountLessThanZero,
				des:    "amount is 0",
			},
			{
				input:  largePlans,
				expect: restricting.ErrCountRestrictPlansInvalid,
				des:    fmt.Sprintf("must less than %d", restricting.RestrictTxPlanSize),
			},
			{
				input:  largeMountPlans,
				expect: restricting.ErrBalanceNotEnough,
				des:    "amount not enough",
			},
			{
				input:  notEnough,
				expect: restricting.ErrLockedAmountTooLess,
				des:    "amount too small",
			},
		}
		for _, value := range x {
			if err := plugin.AddRestrictingRecord(sender, addrArr[0], 20, common.ZeroHash, value.input, mockDB, RestrictingTxHash); err != value.expect {
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

		if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, mockDB, RestrictingTxHash); err != nil {
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
		info.AdvanceAmount = big.NewInt(1e18)
		info.CachePlanAmount = big.NewInt(1e18 + 2e18)
		info.ReleaseList = []uint64{1, 2}
		plugin.storeRestrictingInfo(mockDB, restricting.GetRestrictingKey(to), info)
		info2.NeedRelease = big.NewInt(0)
		info2.AdvanceAmount = big.NewInt(1e18)
		info2.CachePlanAmount = big.NewInt(1e18)
		plugin.storeRestrictingInfo(mockDB, restricting.GetRestrictingKey(account2), info2)
		mockDB.AddBalance(vm.RestrictingContractAddr, big.NewInt(2e18))

		plans := make([]restricting.RestrictingPlan, 0)
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e17)})
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e17)})
		plans = append(plans, restricting.RestrictingPlan{2, big.NewInt(1e18)})
		plans = append(plans, restricting.RestrictingPlan{3, big.NewInt(1e18)})
		if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, mockDB, RestrictingTxHash); err != nil {
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
		assert.Equal(t, big.NewInt(1e18), info2.AdvanceAmount)
		assert.Equal(t, big.NewInt(0), info2.NeedRelease)
		assert.Equal(t, 3, len(info2.ReleaseList))

		balance := mockDB.GetBalance(vm.RestrictingContractAddr)
		assert.Equal(t, big.NewInt(2e18+2e17+2e18), balance)

	})

}

type TestRestrictingPlugin struct {
	RestrictingPlugin
	from, to common.Address
	mockDB   *mock.MockStateDB
}

func NewTestRestrictingPlugin() *TestRestrictingPlugin {
	tp := new(TestRestrictingPlugin)
	tp.log = log.Root()
	tp.from, tp.to = common.MustBech32ToAddress("lax1avltgjnqmy6alefayfry3cd9rpguduawcph8ja"), common.MustBech32ToAddress("lax1rkdnqnnsl5shqm7e00897dpey33h3pcntluqar")
	tp.mockDB = mock.NewChain().StateDB
	tp.mockDB.AddBalance(tp.from, big.NewInt(9e18))
	return tp
}

//the plan is AdvanceLockedFunds,then release, then ReturnLockFunds,the info will delete
func TestRestrictingPlugin_Compose3(t *testing.T) {
	plugin := NewTestRestrictingPlugin()

	sdb := snapshotdb.Instance()
	defer sdb.Clear()
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{"", new(big.Int).SetInt64(0).String(), 0})
	if err := sdb.PutBaseDB(key, value); nil != err {
		t.Error(err)
		return
	}

	plans := make([]restricting.RestrictingPlan, 0)
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e18)})
	if err := plugin.AddRestrictingRecord(plugin.from, plugin.to, xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, plugin.mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}
	if err := plugin.AdvanceLockedFunds(plugin.to, big.NewInt(1e18), plugin.mockDB); err != nil {
		t.Error()
	}
	if err := plugin.releaseRestricting(1, plugin.mockDB); err != nil {
		t.Error(err)
	}
	if err := plugin.ReturnLockFunds(plugin.to, big.NewInt(1e18), plugin.mockDB); err != nil {
		t.Error(err)
	}
	assert.Equal(t, plugin.mockDB.GetBalance(plugin.to), big.NewInt(1e18))
	assert.Equal(t, plugin.mockDB.GetBalance(vm.RestrictingContractAddr).Uint64(), uint64(0))

	_, info := plugin.getRestrictingInfo(plugin.mockDB, plugin.to)
	if len(info) != 0 {
		t.Error("info must del")
	}
}

//the record  exist,have NeedRelease,the NeedRelease amount is less than add plan amount
func TestRestrictingPlugin_Compose2(t *testing.T) {
	plugin := new(RestrictingPlugin)

	sdb := snapshotdb.Instance()
	defer sdb.Clear()
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{"", new(big.Int).SetInt64(0).String(), 0})
	if err := sdb.PutBaseDB(key, value); nil != err {
		t.Error(err)
		return
	}

	plugin.log = log.Root()
	from, to := addrArr[0], addrArr[1]
	//	plugin.log.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	mockDB := buildStateDB(t)
	mockDB.AddBalance(from, big.NewInt(9e18))
	plans := make([]restricting.RestrictingPlan, 0)
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e18)})
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e18)})
	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}
	if err := plugin.AdvanceLockedFunds(to, big.NewInt(2e18), mockDB); err != nil {
		t.Error(err)
	}
	if err := plugin.releaseRestricting(1, mockDB); err != nil {
		t.Error(err)
	}

	plans2 := []restricting.RestrictingPlan{restricting.RestrictingPlan{1, big.NewInt(3e18)}}
	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()+10, common.ZeroHash, plans2, mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}
	infoAssertF := func(CachePlanAmount *big.Int, ReleaseList []uint64, StakingAmount *big.Int, NeedRelease *big.Int) {
		_, info, err := plugin.mustGetRestrictingInfoByDecode(mockDB, to)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, info.CachePlanAmount, CachePlanAmount)
		assert.Equal(t, info.ReleaseList, ReleaseList)
		assert.Equal(t, info.AdvanceAmount, StakingAmount)
		assert.Equal(t, info.NeedRelease, NeedRelease)
	}

	assert.Equal(t, mockDB.GetBalance(from), big.NewInt(4e18))
	assert.Equal(t, mockDB.GetBalance(to), big.NewInt(2e18))
	assert.Equal(t, mockDB.GetBalance(vm.RestrictingContractAddr), big.NewInt(1e18))
	infoAssertF(big.NewInt(3e18), []uint64{2}, big.NewInt(2e18), big.NewInt(0))
}

//the record  exist,have NeedRelease,the NeedRelease amount is grate or equal than  add plan amount
func TestRestrictingPlugin_Compose(t *testing.T) {
	plugin := new(RestrictingPlugin)
	plugin.log = log.Root()
	//	plugin.log.SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	sdb := snapshotdb.Instance()
	defer sdb.Clear()
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{"", new(big.Int).SetInt64(0).String(), 0})
	if err := sdb.PutBaseDB(key, value); nil != err {
		t.Error(err)
		return
	}

	from, to := addrArr[0], addrArr[1]
	mockDB := buildStateDB(t)
	infoAssertF := func(CachePlanAmount *big.Int, ReleaseList []uint64, StakingAmount *big.Int, NeedRelease *big.Int) {
		_, info, err := plugin.mustGetRestrictingInfoByDecode(mockDB, to)
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, info.CachePlanAmount, CachePlanAmount)
		assert.Equal(t, info.ReleaseList, ReleaseList)
		assert.Equal(t, info.AdvanceAmount, StakingAmount)
		assert.Equal(t, info.NeedRelease, NeedRelease)
	}
	mockDB.AddBalance(from, big.NewInt(9e18))
	plans := make([]restricting.RestrictingPlan, 0)
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e18)})
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e18)})

	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}
	assert.Equal(t, mockDB.GetBalance(from), big.NewInt(7e18))
	assert.Equal(t, mockDB.GetBalance(to), big.NewInt(0))
	assert.Equal(t, mockDB.GetBalance(vm.RestrictingContractAddr), big.NewInt(2e18))
	infoAssertF(big.NewInt(2e18), []uint64{1}, big.NewInt(0), big.NewInt(0))

	if err := plugin.AdvanceLockedFunds(to, big.NewInt(2e18), mockDB); err != nil {
		t.Error(err)
	}
	assert.Equal(t, mockDB.GetBalance(to), big.NewInt(0))
	assert.Equal(t, mockDB.GetBalance(vm.RestrictingContractAddr).Uint64(), uint64(0))
	assert.Equal(t, mockDB.GetBalance(vm.StakingContractAddr), big.NewInt(2e18))
	infoAssertF(big.NewInt(2e18), []uint64{1}, big.NewInt(2e18), big.NewInt(0))

	if err := plugin.releaseRestricting(1, mockDB); err != nil {
		t.Error(err)
	}
	assert.Equal(t, mockDB.GetBalance(to).Uint64(), uint64(0))
	assert.Equal(t, mockDB.GetBalance(vm.RestrictingContractAddr).Uint64(), uint64(0))
	assert.Equal(t, mockDB.GetBalance(vm.StakingContractAddr), big.NewInt(2e18))
	infoAssertF(big.NewInt(2e18), []uint64{}, big.NewInt(2e18), big.NewInt(2e18))

	plans2 := []restricting.RestrictingPlan{restricting.RestrictingPlan{1, big.NewInt(1e18)}}
	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()+10, common.ZeroHash, plans2, mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}
	assert.Equal(t, mockDB.GetBalance(from), big.NewInt(6e18))
	assert.Equal(t, mockDB.GetBalance(to), big.NewInt(1e18))
	assert.Equal(t, mockDB.GetBalance(vm.RestrictingContractAddr).Uint64(), uint64(0))
	infoAssertF(big.NewInt(2e18), []uint64{2}, big.NewInt(2e18), big.NewInt(1e18))
}

func TestRestrictingPlugin_GetRestrictingInfo(t *testing.T) {

	sdb := snapshotdb.Instance()
	defer sdb.Clear()
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{"", new(big.Int).SetInt64(0).String(), 0})
	if err := sdb.PutBaseDB(key, value); nil != err {
		t.Error(err)
		return
	}

	t.Run("restricting account not exist", func(t *testing.T) {
		chain := mock.NewChain()
		notFoundAccount := common.HexToAddress("0x11")
		_, err := RestrictingInstance().GetRestrictingInfo(notFoundAccount, chain.StateDB)
		if err != restricting.ErrAccountNotFound {
			t.Errorf("restricting account not exist ,want err %v,have err %v", restricting.ErrAccountNotFound, err)
		}
	})

	t.Run("restricting account exist", func(t *testing.T) {

		chain := mock.NewChain()
		chain.StateDB.AddBalance(addrArr[1], big.NewInt(8e18))

		plans := make([]restricting.RestrictingPlan, 0)
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e18)})
		plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(1e18)})
		plans = append(plans, restricting.RestrictingPlan{2, big.NewInt(1e18)})
		total := new(big.Int)
		for _, value := range plans {
			total.Add(total, value.Amount)
		}
		if err := RestrictingInstance().AddRestrictingRecord(addrArr[1], addrArr[0], xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, chain.StateDB, RestrictingTxHash); err != nil {
			t.Error(err)
		}

		res, err := RestrictingInstance().GetRestrictingInfo(addrArr[0], chain.StateDB)
		if err != nil {
			t.Errorf("get restrictingInfo fail  error: %s", err.Error())
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

	sdb := snapshotdb.Instance()
	defer sdb.Clear()
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{"", new(big.Int).SetInt64(0).String(), 0})
	if err := sdb.PutBaseDB(key, value); nil != err {
		t.Error(err)
		return
	}

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
	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}
	if err := plugin.releaseRestricting(1, mockDB); err != nil {
		t.Error(err)
	}
	//	SetLatestEpoch(mockDB, 1)
	if err := plugin.AdvanceLockedFunds(to, big.NewInt(5e18), mockDB); err != nil {
		t.Error(err)
	}
	if err := plugin.releaseRestricting(2, mockDB); err != nil {
		t.Error(err)
	}
	//	SetLatestEpoch(mockDB, 2)
	if err := plugin.releaseRestricting(3, mockDB); err != nil {
		t.Error(err)
	}
	//	SetLatestEpoch(mockDB, 3)
	plans2 := make([]restricting.RestrictingPlan, 0)
	plans2 = append(plans2, restricting.RestrictingPlan{1, big.NewInt(1e18)})
	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()*3+10, common.ZeroHash, plans2, mockDB, RestrictingTxHash); err != nil {
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
	//	SetLatestEpoch(mockDB, 4)

	assert.Equal(t, big.NewInt(9e18).Add(big.NewInt(9e18), big.NewInt(1e18)), mockDB.GetBalance(to))
	assert.Equal(t, true, mockDB.GetBalance(vm.RestrictingContractAddr).Cmp(big.NewInt(0)) == 0)
	assert.Equal(t, true, mockDB.GetBalance(vm.StakingContractAddr).Cmp(big.NewInt(0)) == 0)
}

func TestNewRestrictingPlugin_MixPledgeLockFunds(t *testing.T) {
	sdb := snapshotdb.Instance()
	defer sdb.Clear()
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{"", new(big.Int).SetInt64(0).String(), 0})
	if err := sdb.PutBaseDB(key, value); nil != err {
		t.Error(err)
		return
	}

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
	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}
	mockDB.AddBalance(to, big.NewInt(2e18))

	res, free, err := plugin.MixAdvanceLockedFunds(to, new(big.Int).Mul(big.NewInt(1e18), big.NewInt(10)), mockDB)
	if err != nil {
		t.Error(err)
	}

	if res.Cmp(big.NewInt(9e18)) != 0 {
		t.Errorf("restricting von cost wrong,%v", res)
	}
	if free.Cmp(big.NewInt(1e18)) != 0 {
		t.Errorf("free von cost wrong,%v", free)
	}

	if mockDB.GetBalance(to).Cmp(big.NewInt(1e18)) != 0 {
		t.Errorf("to balance von cost wrong")
	}

	if _, _, err := plugin.MixAdvanceLockedFunds(to, new(big.Int).Mul(big.NewInt(1e18), big.NewInt(10)), mockDB); err == nil {
		t.Error("should not success")
	}

}

func TestRestrictingInstanceWithSlashing(t *testing.T) {

	sdb := snapshotdb.Instance()
	defer sdb.Clear()
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{"", new(big.Int).SetInt64(0).String(), 0})
	if err := sdb.PutBaseDB(key, value); nil != err {
		t.Error(err)
		return
	}
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
	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}

	if err := plugin.releaseRestricting(1, mockDB); err != nil {
		t.Error(err)
	}
	//	SetLatestEpoch(mockDB, 1)

	if err := plugin.AdvanceLockedFunds(to, big.NewInt(5e18), mockDB); err != nil {
		t.Error(err)
	}

	if err := plugin.releaseRestricting(2, mockDB); err != nil {
		t.Error(err)
	}
	//	SetLatestEpoch(mockDB, 2)

	if err := plugin.releaseRestricting(3, mockDB); err != nil {
		t.Error(err)
	}
	//	SetLatestEpoch(mockDB, 3)

	mockDB.SubBalance(vm.StakingContractAddr, big.NewInt(1e18))
	if err := plugin.SlashingNotify(to, big.NewInt(1e18), mockDB); err != nil {
		t.Error(err)
	}

	plans2 := make([]restricting.RestrictingPlan, 0)
	plans2 = append(plans2, restricting.RestrictingPlan{1, big.NewInt(1e18)})
	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()*3+10, common.ZeroHash, plans2, mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}
	if err := plugin.ReturnLockFunds(to, big.NewInt(4e18), mockDB); err != nil {
		t.Error(err)
	}

	assert.Equal(t, big.NewInt(9e18), mockDB.GetBalance(to))

	if err := plugin.releaseRestricting(4, mockDB); err != nil {
		t.Error(err)
	}
	//	SetLatestEpoch(mockDB, 4)

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
	//	SetLatestEpoch(mockDB, 5)

}

func TestRestrictingGetRestrictingInfo(t *testing.T) {

	sdb := snapshotdb.Instance()
	defer sdb.Clear()
	key := gov.KeyParamValue(gov.ModuleRestricting, gov.KeyRestrictingMinimumAmount)
	value := common.MustRlpEncode(&gov.ParamValue{"", new(big.Int).SetInt64(0).String(), 0})
	if err := sdb.PutBaseDB(key, value); nil != err {
		t.Error(err)
		return
	}
	mockDB := buildStateDB(t)
	plugin := new(RestrictingPlugin)
	plugin.log = log.Root()
	from, to := addrArr[0], addrArr[1]
	mockDB.AddBalance(from, big.NewInt(9e18).Add(big.NewInt(9e18), big.NewInt(9e18)))
	plans := make([]restricting.RestrictingPlan, 0)
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(3e18)})
	plans = append(plans, restricting.RestrictingPlan{1, big.NewInt(3e18)})

	if err := plugin.AddRestrictingRecord(from, to, xutil.CalcBlocksEachEpoch()-10, common.ZeroHash, plans, mockDB, RestrictingTxHash); err != nil {
		t.Error(err)
	}
	res, err := plugin.getRestrictingInfoToReturn(to, mockDB)
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, res.Balance.ToInt(), big.NewInt(6e18))

}
