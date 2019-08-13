package plugin

import (
	"encoding/json"
	"math/big"
	"sort"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

var (
	errParamEpochInvalid   = common.NewBizError("param epoch can't be zero")
	errEmptyRestrictPlan   = common.NewBizError("the number of the restricting plan can't be zero")
	errTooMuchPlan         = common.NewBizError("the number of the restricting plan is too much")
	errLockedAmountTooLess = common.NewBizError("total restricting amount need more than 1 LAT")
	errBalanceNotEnough    = common.NewBizError("the balance is not enough in restrict")
	errAccountNotFound     = common.NewBizError("account is not found on restricting contract")
	monthOfThreeYear       = 12 * 3
)

type RestrictingPlugin struct {
	log log.Logger
}

var (
	restrictingOnce sync.Once
	rt              *RestrictingPlugin
)

func RestrictingInstance() *RestrictingPlugin {
	restrictingOnce.Do(func() {
		log2 := log.Root().New("package", "RestrictingPlugin")
		log2.Info("Init Restricting plugin ...")
		rt = &RestrictingPlugin{log2}
	})
	return rt
}

/*func ClearRestricting() error {
	if nil == rt {
		return common.NewSysError("the RestrictingPlugin already be nil")
	}
	rt = nil
	return nil
}*/

// BeginBlock does something like check input params before execute transactions,
// in RestrictingPlugin it does nothing.
func (rp *RestrictingPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	return nil
}

// EndBlock invoke releaseRestricting
func (rp *RestrictingPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	expect := GetLatestEpoch(state) + 1
	expectBlock := GetBlockNumberByEpoch(expect)

	if expectBlock != head.Number.Uint64() {
		rp.log.Debug("not expected block number", "expectEpoch", expect, "expectBlock", expectBlock, "currBlock", head.Number.Uint64())
		return nil
	}

	rp.log.Info("begin to release restricting plan", "curr", head.Number, "epoch", expectBlock)
	return rp.releaseRestricting(expect, state)
}

// Confirmed is empty function
func (rp *RestrictingPlugin) Confirmed(block *types.Block) error {
	return nil
}

func (rp *RestrictingPlugin) mergeAmount(state xcom.StateDB, plans []restricting.RestrictingPlan) (*big.Int, map[uint64]*big.Int, error) {
	// latest is the epoch of a settlement block closest to current block
	latestEpoch := GetLatestEpoch(state)

	totalAmount := new(big.Int)

	mPlans := make(map[uint64]*big.Int, monthOfThreeYear)
	for _, plan := range plans {
		epoch, amount := plan.Epoch, plan.Amount
		if epoch == 0 {
			rp.log.Error(errParamEpochInvalid.Error())
			return nil, nil, errParamEpochInvalid
		}
		totalAmount.Add(totalAmount, amount)
		newEpoch := epoch + latestEpoch
		if value, ok := mPlans[newEpoch]; ok {
			mPlans[newEpoch] = value.Add(amount, value)
		} else {
			mPlans[newEpoch] = amount
		}
	}
	return totalAmount, mPlans, nil
}

func (rp *RestrictingPlugin) initEpochInfo(state xcom.StateDB, epoch uint64, account common.Address, amount *big.Int) {
	// step1: get account numbers at target epoch
	releaseEpochKey, bAccNumbers := rp.getReleaseEpochNumber(state, epoch)
	var accNumbers uint32
	if bAccNumbers == 0 {
		accNumbers = uint32(1)
	} else {
		accNumbers = bAccNumbers + 1
	}
	// step2: save the numbers of restricting record at target epoch
	rp.storeNumber2ReleaseEpoch(state, releaseEpochKey, accNumbers)

	// step3: save account at target index
	rp.storeAccount2ReleaseAccount(state, epoch, accNumbers, account)

	if amount != nil {
		// step4: save restricting amount at target epoch
		rp.storeAmount2ReleaseAmount(state, epoch, account, amount)
	}
}

// AddRestrictingRecord stores four K-V record in StateDB:
// RestrictingInfo: the account info to be released
// ReleaseEpoch:   the number of accounts to be released on the epoch corresponding to the target block height
// ReleaseAccount: the account on the index on the target epoch
// ReleaseAmount: the amount of the account to be released on the target epoch
func (rp *RestrictingPlugin) AddRestrictingRecord(sender common.Address, account common.Address, plans []restricting.RestrictingPlan, state xcom.StateDB) error {

	rp.log.Info("begin to addRestrictingRecord", "sender", sender.String(), "account", account.String(), "plans", plans)

	if len(plans) == 0 {
		rp.log.Error(errEmptyRestrictPlan.Error())
		return errEmptyRestrictPlan
	}
	// totalAmount is total restricting amount
	totalAmount, mPlans, err := rp.mergeAmount(state, plans)
	if err != nil {
		return err
	}
	// pre-check
	{
		if len(mPlans) > monthOfThreeYear {
			rp.log.Error("the number of the restricting plan must less or equal than monthOfThreeYear", "monthOfThreeYear", monthOfThreeYear, "have", len(mPlans))
			return errTooMuchPlan
		}

		if totalAmount.Cmp(big.NewInt(1E18)) < 0 {
			rp.log.Error("total restricting amount need more than 1 LAT", "sender", sender, "amount", totalAmount)
			return errLockedAmountTooLess
		}

		if state.GetBalance(sender).Cmp(totalAmount) < 0 {
			rp.log.Error("balance of the sender is not enough", "total", totalAmount, "balance", state.GetBalance(sender))
			return errBalanceNotEnough
		}
	}

	var (
		epochList []uint64
		info      restricting.RestrictingInfo
		repay     = common.Big0
	)

	restrictingKey, bAccInfo := rp.getRestrictingInfo(state, account)
	if len(bAccInfo) == 0 {
		rp.log.Info("restricting record not exist", "account", account.String())
		for epoch, amount := range mPlans {
			rp.initEpochInfo(state, epoch, account, amount)
			epochList = append(epochList, epoch)
		}
		info.Balance = totalAmount
		info.Debt = big.NewInt(0)
		info.DebtSymbol = false
		info.ReleaseList = epochList
	} else {
		rp.log.Info("restricting record exist", "account", account.String())
		if err = rlp.DecodeBytes(bAccInfo, &info); err != nil {
			rp.log.Error("failed to rlp decode the restricting account", "err", err.Error())
			return common.NewSysError(err.Error())
		}

		if info.DebtSymbol && info.Debt.Cmp(common.Big0) > 0 {
			if totalAmount.Cmp(info.Debt) >= 0 {
				repay = info.Debt
				totalAmount.Sub(totalAmount, info.Debt)
				info.Debt = common.Big0
				info.DebtSymbol = false
			} else {
				repay = totalAmount
				totalAmount = common.Big0
				info.Debt.Sub(info.Debt, totalAmount)
			}
		}
		info.Balance.Add(info.Balance, totalAmount)

		for epoch, releaseAmount := range mPlans {
			// step1: get restricting amount at target epoch
			releaseAmountKey, currentAmount := rp.getReleaseAmount(state, epoch, account)
			if currentAmount.Cmp(common.Big0) == 0 {
				rp.log.Info("release record not exist on curr epoch ", "account", account.String(), "epoch", epoch)
				rp.initEpochInfo(state, epoch, account, nil)
				info.ReleaseList = append(info.ReleaseList, epoch)
			} else {
				rp.log.Info("release record exist at curr epoch", "account", account.String(), "epoch", epoch)
				currentAmount.Add(currentAmount, releaseAmount)
				// step4: save restricting amount at target epoch
				state.SetState(vm.RestrictingContractAddr, releaseAmountKey, currentAmount.Bytes())
			}
		}
	}
	// sort release list
	sort.Slice(info.ReleaseList, func(i, j int) bool {
		return info.ReleaseList[i] < info.ReleaseList[j]
	})
	if err := rp.storeRestrictingInfo(state, restrictingKey, info); err != nil {
		rp.log.Error("failed to rlp encode restricting info", "account", account.String(), "error", err)
		return common.NewSysError(err.Error())
	}
	if repay.Cmp(common.Big0) > 0 {
		state.AddBalance(account, repay)
	}
	state.SubBalance(sender, totalAmount)
	state.AddBalance(vm.RestrictingContractAddr, totalAmount)

	rp.log.Debug("end to addRestrictingRecord", "account", account.String(), "restrictingInfo", bAccInfo)

	return nil
}

// PledgeLockFunds transfer the money from the restricting contract account to the staking contract account
func (rp *RestrictingPlugin) PledgeLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {

	rp.log.Info("begin to PledgeLockFunds", "account", account.String(), "amount", amount)
	restrictingKey, info, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return err
	}
	if info.Balance.Cmp(amount) < 0 {
		rp.log.Warn("Balance of restricting account not enough", "balance", info.Balance, "funds", amount)
		return errBalanceNotEnough
	}

	// sub Balance
	info.Balance.Sub(info.Balance, amount)

	// save restricting account info
	if err := rp.storeRestrictingInfo(state, restrictingKey, info); err != nil {
		rp.log.Error("failed to rlp encode the restricting account", "account", account.String(), "error", err)
		return common.NewSysError(err.Error())
	}

	state.SubBalance(vm.RestrictingContractAddr, amount)
	state.AddBalance(vm.StakingContractAddr, amount)

	rp.log.Info("end to PledgeLockFunds", "RCContractBalance", state.GetBalance(vm.RestrictingContractAddr), "STKContractBalance", state.GetBalance(vm.StakingContractAddr))
	return nil
}

// ReturnLockFunds transfer the money from the staking contract account to the restricting contract account
func (rp *RestrictingPlugin) ReturnLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {

	rp.log.Info("begin to ReturnLockFunds", "account", account.String(), "amount", amount)
	var (
		repay = new(big.Int) // repay the money owed in the past
		left  = new(big.Int) // money left after the repayment
	)

	restrictingKey, info, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return err
	}

	if info.DebtSymbol {
		rp.log.Info("Balance was owed to release in the past", "account", account.String(), "Debt", info.Debt, "funds", amount)

		if amount.Cmp(info.Debt) < 0 {
			// the money returned back is not enough to repay the money owed to release
			repay = amount
			info.Debt.Sub(info.Debt, amount)
		} else {
			// the money returned back is more than the money owed to release
			repay = info.Debt

			left.Sub(amount, info.Debt)
			if left.Cmp(big.NewInt(0)) > 0 {
				info.Balance.Add(info.Balance, left)
			}

			info.Debt = big.NewInt(0)
			info.DebtSymbol = false
		}

	} else {
		rp.log.Info("directly add Balance while symbol is false", "account", account.String(), "Debt", info.Debt)
		repay = big.NewInt(0)
		left = amount
		info.Balance.Add(info.Balance, left)
	}

	// save restricting account info
	if err := rp.storeRestrictingInfo(state, restrictingKey, info); err != nil {
		rp.log.Error("failed to rlp encode the restricting account", "func", "ReturnLockFunds", "account", account.String(), "error", err)
		return common.NewSysError(err.Error())
	}

	state.SubBalance(vm.StakingContractAddr, amount)
	if repay.Cmp(big.NewInt(0)) == 1 {
		state.AddBalance(account, repay)
	}
	state.AddBalance(vm.RestrictingContractAddr, left)

	rp.log.Info("end to ReturnLockFunds", "RCContractBalance", state.GetBalance(vm.RestrictingContractAddr))

	return nil
}

// SlashingNotify modify Debt of restricting account
func (rp *RestrictingPlugin) SlashingNotify(account common.Address, amount *big.Int, state xcom.StateDB) error {
	rp.log.Info("begin to SlashingNotify", "account", account.String(), "amount", amount)

	restrictingKey, info, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return err
	}
	if info.DebtSymbol {
		rp.log.Info("Balance was owed to release in the past", "account", account.String(), "Debt", info.Debt, "funds", amount)
		if amount.Cmp(info.Debt) < 0 {
			info.Debt.Sub(info.Debt, amount)
		} else {
			info.Debt.Sub(amount, info.Debt)
			info.DebtSymbol = false
		}
	} else {
		info.Debt = info.Debt.Add(info.Debt, amount)
	}

	// save restricting account info
	if err := rp.storeRestrictingInfo(state, restrictingKey, info); err != nil {
		rp.log.Error("failed to encode restricting account", "func", "SlashingNotify", "account", account.String(), "error", err)
		return common.NewSysError(err.Error())
	}
	rp.log.Info("begin to SlashingNotify", "restrictingInfo", info)

	return nil
}

func (rp *RestrictingPlugin) getReleaseEpochNumber(state xcom.StateDB, epoch uint64) ([]byte, uint32) {
	releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
	bAccNumbers := state.GetState(vm.RestrictingContractAddr, releaseEpochKey)
	return releaseEpochKey, common.BytesToUint32(bAccNumbers)
}

func (rp *RestrictingPlugin) getReleaseAccount(state xcom.StateDB, epoch uint64, index uint32) ([]byte, common.Address) {
	releaseAccountKey := restricting.GetReleaseAccountKey(epoch, index)
	bAccount := state.GetState(vm.RestrictingContractAddr, releaseAccountKey)
	account := common.BytesToAddress(bAccount)
	return releaseAccountKey, account
}

func (rp *RestrictingPlugin) getRestrictingInfo(state xcom.StateDB, account common.Address) ([]byte, []byte) {
	restrictingKey := restricting.GetRestrictingKey(account)
	bAccInfo := state.GetState(vm.RestrictingContractAddr, restrictingKey)
	return restrictingKey, bAccInfo
}

func (rp *RestrictingPlugin) mustGetRestrictingInfoByDecode(state xcom.StateDB, account common.Address) ([]byte, restricting.RestrictingInfo, error) {
	var info restricting.RestrictingInfo
	restrictingKey, bAccInfo := rp.getRestrictingInfo(state, account)
	if len(bAccInfo) == 0 {
		rp.log.Warn("record not found in GetRestrictingInfo", "account", account.String())
		return []byte{}, info, errAccountNotFound
	}
	if err := rlp.DecodeBytes(bAccInfo, &info); err != nil {
		rp.log.Error("failed to rlp decode restricting account", "error", err.Error(), "info", bAccInfo)
		return restrictingKey, info, common.NewSysError(err.Error())
	}
	return restrictingKey, info, nil
}

func (rp *RestrictingPlugin) getRestrictingInfoByDecode(state xcom.StateDB, account common.Address) ([]byte, restricting.RestrictingInfo, error) {
	restrictingKey, bAccInfo := rp.getRestrictingInfo(state, account)
	var info restricting.RestrictingInfo
	if err := rlp.DecodeBytes(bAccInfo, &info); err != nil {
		rp.log.Error("failed to rlp decode restricting account", "error", err.Error(), "info", bAccInfo)
		return restrictingKey, info, common.NewSysError(err.Error())
	}
	return restrictingKey, info, nil
}

func (rp *RestrictingPlugin) getReleaseAmount(state xcom.StateDB, epoch uint64, account common.Address) ([]byte, *big.Int) {
	releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)
	bRelease := state.GetState(vm.RestrictingContractAddr, releaseAmountKey)
	release := new(big.Int)
	release.SetBytes(bRelease)
	return releaseAmountKey, release
}

func (rp *RestrictingPlugin) storeRestrictingInfo(state xcom.StateDB, restrictingKey []byte, info restricting.RestrictingInfo) error {
	bNewInfo, err := rlp.EncodeToBytes(info)
	if err != nil {
		return err
	}
	state.SetState(vm.RestrictingContractAddr, restrictingKey, bNewInfo)
	return nil
}

func (rp *RestrictingPlugin) storeNumber2ReleaseEpoch(state xcom.StateDB, releaseEpochKey []byte, accNumbers uint32) {
	state.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(accNumbers))
}

func (rp *RestrictingPlugin) storeAccount2ReleaseAccount(state xcom.StateDB, epoch uint64, index uint32, account common.Address) {
	releaseAccountKey := restricting.GetReleaseAccountKey(epoch, index)
	state.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())
}

func (rp *RestrictingPlugin) storeAmount2ReleaseAmount(state xcom.StateDB, epoch uint64, account common.Address, amount *big.Int) {
	releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)
	state.SetState(vm.RestrictingContractAddr, releaseAmountKey, amount.Bytes())
}

//todo 在完成之前不能删除key
// releaseRestricting will release restricting plans on target epoch
func (rp *RestrictingPlugin) releaseRestricting(epoch uint64, state xcom.StateDB) error {

	rp.log.Info("begin to releaseRestricting", "epoch", epoch)

	releaseEpochKey, numbers := rp.getReleaseEpochNumber(state, epoch)
	if numbers == 0 {
		rp.log.Info("there is no release record on curr epoch", "epoch", epoch)
		return nil
	}

	rp.log.Info("many restricting records need release", "epoch", epoch, "records", numbers)

	for index := numbers; index > 0; index-- {
		releaseAccountKey, account := rp.getReleaseAccount(state, epoch, index)

		rp.log.Info("begin to release record", "index", index, "account", account.String())

		restrictingKey, info, err := rp.getRestrictingInfoByDecode(state, account)
		if err != nil {
			return err
		}

		releaseAmountKey, releaseAmount := rp.getReleaseAmount(state, epoch, account)

		if info.DebtSymbol {
			rp.log.Debug("Balance is owed to release in the past", "account", account.String(), "Debt", info.Debt, "symbol", info.DebtSymbol)
			info.Debt.Add(info.Debt, releaseAmount)
		} else {
			// release amount isn't more than debt
			if releaseAmount.Cmp(info.Debt) <= 0 {
				info.Debt.Sub(info.Debt, releaseAmount)
			} else if releaseAmount.Cmp(new(big.Int).Add(info.Debt, info.Balance)) <= 0 {
				// release amount isn't more than the sum of balance and debt
				releaseAmount.Sub(releaseAmount, info.Debt)

				info.Balance.Sub(info.Balance, releaseAmount)
				info.Debt = big.NewInt(0)

				rp.log.Debug("show balance", "balance", info.Balance)

				state.SubBalance(vm.RestrictingContractAddr, releaseAmount)
				state.AddBalance(account, releaseAmount)

			} else {
				// release amount is more than the sum of balance and debt
				origBalance := info.Balance

				releaseAmount.Sub(releaseAmount, info.Balance)
				info.Balance = big.NewInt(0)
				info.Debt.Sub(releaseAmount, info.Debt)
				info.DebtSymbol = true

				state.SubBalance(vm.RestrictingContractAddr, origBalance)
				state.AddBalance(account, origBalance)
			}
		}

		// delete ReleaseAmount
		state.SetState(vm.RestrictingContractAddr, releaseAmountKey, []byte{})

		// delete ReleaseAccount
		state.SetState(vm.RestrictingContractAddr, releaseAccountKey, []byte{})

		// delete epoch in ReleaseList
		// In general, the first epoch is released first.
		// info.ReleaseList = info.ReleaseList[1:]
		info.RemoveEpoch(epoch)

		// just restore restricting info, don't delete
		if err := rp.storeRestrictingInfo(state, restrictingKey, info); err != nil {
			rp.log.Error("failed to rlp encode new info while release", "account", account.String(), "info", info)
			return common.NewSysError(err.Error())
		}
	}

	// delete ReleaseEpoch
	state.SetState(vm.RestrictingContractAddr, releaseEpochKey, []byte{})

	rp.log.Info("end to releaseRestricting")

	return nil
}

func (rp *RestrictingPlugin) GetRestrictingInfo(account common.Address, state xcom.StateDB) ([]byte, error) {

	rp.log.Info("begin to GetRestrictingInfo", "account", account.String())

	_, info, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return nil, err
	}

	var (
		plan   restricting.ReleaseAmountInfo
		plans  []restricting.ReleaseAmountInfo
		result restricting.Result
	)

	for i := 0; i < len(info.ReleaseList); i++ {
		epoch := info.ReleaseList[i]
		_, bAmount := rp.getReleaseAmount(state, epoch, account)
		plan.Height = GetBlockNumberByEpoch(epoch)
		plan.Amount = bAmount
		plans = append(plans, plan)
	}

	result.Balance = info.Balance
	result.Debt = info.Debt
	result.Symbol = info.DebtSymbol
	result.Entry = plans
	rp.log.Info("get restricting result", "account", account.String(), "result", result)

	bResult, err := json.Marshal(result)
	if err != nil {
		rp.log.Error("failed to Marshal restricting result")
		return []byte{}, err
	}

	rp.log.Info("end to GetRestrictingInfo", "restrictingInfo", bResult)

	return bResult, nil
}

// state DB operation
func SetLatestEpoch(stateDb xcom.StateDB, epoch uint64) {
	key := restricting.GetLatestEpochKey()
	stateDb.SetState(vm.RestrictingContractAddr, key, common.Uint64ToBytes(epoch))
}

func GetLatestEpoch(stateDb xcom.StateDB) uint64 {
	key := restricting.GetLatestEpochKey()
	bEpoch := stateDb.GetState(vm.RestrictingContractAddr, key)
	return common.BytesToUint64(bEpoch)
}

func GetBlockNumberByEpoch(epoch uint64) uint64 {
	return epoch * xutil.CalcBlocksEachEpoch()
}
