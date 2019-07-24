package plugin

import (
	"bytes"
	"encoding/json"
	"math/big"
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
	errTooMuchPlan         = common.NewBizError("the number of the restricting plan is too much")
	errLockedAmountTooLess = common.NewBizError("total restricting amount need more than 1 LAT")
	errBalanceNotEnough    = common.NewBizError("balance not enough to restrict")
	errAccountNotFound     = common.NewBizError("account is not found")
	monthOfThreeYear       = 12 * 3
)

type RestrictingPlugin struct {
}

var (
	restrictingOnce sync.Once
	rt              *RestrictingPlugin
)

func RestrictingInstance() *RestrictingPlugin {
	restrictingOnce.Do(func() {
		rt = &RestrictingPlugin{}
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
		return nil
	}

	log.Info("begin to release restricting plan", "curr", head.Number, "epoch", expectBlock)
	return rp.releaseRestricting(expect, state)
}

// Confirmed is empty function
func (rp *RestrictingPlugin) Confirmed(block *types.Block) error {
	return nil
}

// AddRestrictingRecord stores four K-V record in StateDB:
// RestrictingInfo: the account info to be released
// ReleaseEpoch:   the number of accounts to be released on the epoch corresponding to the target block height
// ReleaseAccount: the account on the index on the target epoch
// ReleaseAmount: the amount of the account to be released on the target epoch
func (rp *RestrictingPlugin) AddRestrictingRecord(sender common.Address, account common.Address, plans []restricting.RestrictingPlan,
	state xcom.StateDB) error {

	// latest is the epoch of a settlement block closest to current block
	latest := GetLatestEpoch(state)
	// totalAmount is total restricting amount
	totalAmount := new(big.Int)

	// merge the amount of the same release epoch
	mPlans := make(map[uint64]*big.Int, monthOfThreeYear)
	for _, plan := range plans {

		k, v := plan.Epoch, plan.Amount
		if k == 0 {
			log.Debug("param epoch can't be zero")
			return errParamEpochInvalid
		}

		k += latest
		if mPlans[k] == nil {
			mPlans[k] = v
		} else {
			mPlans[k] = v.Add(v, mPlans[k])
		}
		totalAmount = totalAmount.Add(totalAmount, v)
	}

	// pre-check
	if len(mPlans) > monthOfThreeYear {
		log.Debug("the number of the restricting plan must less or equal than %d", monthOfThreeYear)
		return errTooMuchPlan
	}

	if totalAmount.Cmp(big.NewInt(1E18)) == -1 {
		log.Debug("total restricting amount need more than 1 LAT", "sender", sender, "amount", totalAmount)
		return errLockedAmountTooLess
	}

	if state.GetBalance(sender).Cmp(totalAmount) == -1 {
		log.Debug("balance of the sender is not enough", "total", totalAmount, "balance", state.GetBalance(sender))
		return errBalanceNotEnough
	}

	// TODO
	var (
		err        error
		epochList  []uint64
		index      uint32
		info       restricting.RestrictingInfo
		accNumbers uint32
	)

	restrictingKey := restricting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Debug("restricting record not exist", "account", account.Bytes())

		for epoch, amount := range mPlans {
			// step1: get account numbers at target epoch
			releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
			bAccNumbers := state.GetState(vm.RestrictingContractAddr, releaseEpochKey)

			if len(bAccNumbers) == 0 {
				accNumbers = uint32(1)
			} else {
				accNumbers = common.BytesToUint32(bAccNumbers) + 1
			}
			index = accNumbers

			// step2: save the numbers of restricting record at target epoch
			state.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(accNumbers))

			// step3: save account at target index
			releaseAccountKey := restricting.GetReleaseAccountKey(epoch, index)
			state.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

			// step4: save restricting amount at target epoch
			releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)

			state.SetState(account, releaseAmountKey, amount.Bytes())

			epochList = append(epochList, epoch)
		}

		info.Balance = totalAmount
		info.Debt = big.NewInt(0)
		info.DebtSymbol = false
		info.ReleaseList = epochList

	} else {
		log.Debug("restricting record exist", "account", account.String())

		if err = rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
			log.Error("failed to rlp decode the restricting account", "err", err.Error())
			return common.NewSysError(err.Error())
		}

		for epoch, amount := range mPlans {
			// step1: get restricting amount at target epoch
			releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)
			bAmount := state.GetState(account, releaseAmountKey)

			if len(bAmount) == 0 {
				log.Trace("release record not exist on curr epoch ", "account", account, "epoch", epoch)

				releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
				bAccNumbers := state.GetState(vm.RestrictingContractAddr, releaseEpochKey)

				if len(bAccNumbers) == 0 {
					accNumbers = uint32(1)
				} else {
					accNumbers = common.BytesToUint32(bAccNumbers) + 1
				}
				index = accNumbers

				// step2: save account numbers at target epoch
				state.SetState(vm.RestrictingContractAddr, releaseEpochKey, common.Uint32ToBytes(accNumbers))

				// step3: save account at target index
				releaseAccountKey := restricting.GetReleaseAccountKey(epoch, index)
				state.SetState(vm.RestrictingContractAddr, releaseAccountKey, account.Bytes())

				info.ReleaseList = append(info.ReleaseList, epoch)

			} else {
				log.Trace("release record exist at curr epoch", "account", account, "epoch", epoch)

				origAmount := new(big.Int)
				origAmount = origAmount.SetBytes(bAmount)
				amount = amount.Add(amount, origAmount)
			}

			// step4: save restricting amount at target epoch
			state.SetState(account, releaseAmountKey, amount.Bytes())
		}

		info.Balance = info.Balance.Add(info.Balance, totalAmount)
	}

	// step5: save restricting account info
	bAccInfo, err = rlp.EncodeToBytes(info)
	if err != nil {
		log.Error("failed to rlp encode restricting info", "account", account, "error", err)
		return common.NewSysError(err.Error())
	}

	state.SetState(account, restrictingKey, bAccInfo)
	state.SubBalance(sender, totalAmount)
	state.AddBalance(vm.RestrictingContractAddr, totalAmount)

	return nil
}

// PledgeLockFunds transfer the money from the restricting contract account to the staking contract account
func (rp *RestrictingPlugin) PledgeLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {

	restrictingKey := restricting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Debug("record not found in PledgeLockFunds", "account", account, "funds", amount)
		return errAccountNotFound
	}

	var (
		err  error
		info restricting.RestrictingInfo
	)

	if err := rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
		log.Error("failed to rlp decode the restricting account", "info", bAccInfo, "error", err.Error())
		return common.NewSysError(err.Error())
	}

	if info.Balance.Cmp(amount) == -1 {
		log.Debug("Balance of restricting account not enough", "balance", info.Balance, "funds", amount)
		return errBalanceNotEnough
	}

	// sub Balance
	info.Balance = info.Balance.Sub(info.Balance, amount)

	// save restricting account info
	if bAccInfo, err = rlp.EncodeToBytes(info); err != nil {
		log.Error("failed to rlp encode the restricting account", "account", account, "error", err)
		return common.NewSysError(err.Error())
	}
	state.SetState(account, restrictingKey, bAccInfo)

	state.SubBalance(vm.RestrictingContractAddr, amount)
	state.AddBalance(vm.StakingContractAddr, amount)

	return nil
}

// ReturnLockFunds transfer the money from the staking contract account to the restricting contract account
func (rp *RestrictingPlugin) ReturnLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {

	restrictingKey := restricting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Debug("record not found in ReturnLockFunds", "account", account, "funds", amount)
		return errAccountNotFound
	}

	var (
		err   error
		info  restricting.RestrictingInfo
		repay = new(big.Int) // repay the money owed in the past
		left  = new(big.Int) // money left after the repayment
	)

	if err = rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
		log.Error("failed to rlp encode the restricting account", "error", err.Error())
		return common.NewSysError(err.Error())
	}

	if info.DebtSymbol {
		log.Trace("Balance was owed to release in the past", "account", account, "Debt", info.Debt, "funds", amount)

		if amount.Cmp(info.Debt) == -1 {
			// the money returned back is not enough to repay the money owed to release
			repay = amount
			info.Debt = info.Debt.Sub(info.Debt, amount)

		} else {
			// the money returned back is more than the money owed to release
			repay = info.Debt

			left = left.Sub(amount, info.Debt)
			if left.Cmp(big.NewInt(0)) == 1 {
				info.Balance = info.Balance.Add(info.Balance, left)
			}

			info.Debt = big.NewInt(0)
			info.DebtSymbol = false
		}

	} else {
		log.Trace("directly add Balance while symbol is false", "account", account, "Debt", info.Debt)

		repay = big.NewInt(0)
		left = amount
		info.Balance = info.Balance.Add(info.Balance, left)
	}

	// save restricting account info
	if bAccInfo, err = rlp.EncodeToBytes(info); err != nil {
		log.Error("failed to rlp encode the restricting account", "account", account, "error", err)
		return common.NewSysError(err.Error())
	}
	state.SetState(account, restrictingKey, bAccInfo)

	state.SubBalance(vm.StakingContractAddr, amount)
	if repay.Cmp(big.NewInt(0)) == 1 {
		state.AddBalance(account, repay)
	}
	state.AddBalance(vm.RestrictingContractAddr, left)

	return nil
}

// SlashingNotify modify Debt of restricting account
func (rp *RestrictingPlugin) SlashingNotify(account common.Address, amount *big.Int, state xcom.StateDB) error {

	restrictingKey := restricting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Debug("record not found in SlashingNotify", "account", account, "funds", amount)
		return errAccountNotFound
	}

	var (
		err  error
		info restricting.RestrictingInfo
	)

	if err = rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
		log.Error("failed to rlp decode restricting account", "error", err.Error(), "info", bAccInfo)
		return common.NewSysError(err.Error())
	}

	if info.DebtSymbol {
		log.Trace("Balance was owed to release in the past", "account", account, "Debt", info.Debt, "funds", amount)

		if amount.Cmp(info.Debt) < 0 {
			info.Debt = info.Debt.Sub(info.Debt, amount)

		} else {
			info.Debt = info.Debt.Sub(amount, info.Debt)
			info.DebtSymbol = false
		}

	} else {
		info.Debt = info.Debt.Add(info.Debt, amount)
	}

	// save restricting account info
	if bAccInfo, err = rlp.EncodeToBytes(info); err != nil {
		log.Error("failed to encode restricting account", "account", account, "error", err)
		return common.NewSysError(err.Error())
	}
	state.SetState(account, restrictingKey, bAccInfo)

	return nil
}

// releaseRestricting will release restricting plans on target epoch
func (rp *RestrictingPlugin) releaseRestricting(epoch uint64, state xcom.StateDB) error {

	releaseEpochKey := restricting.GetReleaseEpochKey(epoch)
	bAccNumbers := state.GetState(vm.RestrictingContractAddr, releaseEpochKey)

	if len(bAccNumbers) == 0 {
		log.Debug("there is no release record on curr epoch", "epoch", epoch)
		return nil
	}
	numbers := common.BytesToUint32(bAccNumbers)
	log.Debug("many restricting records need release", "epoch", epoch, "records", numbers)

	// TODO
	var (
		info    restricting.RestrictingInfo
		release = new(big.Int) // amount need released
	)

	for index := numbers; index > 0; index-- {

		releaseAccountKey := restricting.GetReleaseAccountKey(epoch, index)
		bAccount := state.GetState(vm.RestrictingContractAddr, releaseAccountKey)
		account := common.BytesToAddress(bAccount)

		log.Trace("begin to release record", "index", index, "account", account)

		restrictingKey := restricting.GetRestrictingKey(account)
		bAccInfo := state.GetState(account, restrictingKey)

		if err := rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
			log.Error("failed to rlp decode restricting account", "error", err.Error(), "info", bAccInfo)
			return common.NewSysError(err.Error())
		}

		releaseAmountKey := restricting.GetReleaseAmountKey(epoch, account)
		bRelease := state.GetState(account, releaseAmountKey)
		release = release.SetBytes(bRelease)

		if info.DebtSymbol {
			log.Debug("Balance is owed to release in the past", "account", account, "Debt", info.Debt, "symbol", info.DebtSymbol)
			info.Debt = info.Debt.Add(info.Debt, release)

		} else {

			// release amount isn't more than debt
			if release.Cmp(info.Debt) <= 0 {
				info.Debt = info.Debt.Sub(info.Debt, release)

			} else if release.Cmp(new(big.Int).Add(info.Debt, info.Balance)) <= 0 {
				// release amount isn't more than the sum of balance and debt
				release = release.Sub(release, info.Debt)
				info.Balance = info.Balance.Sub(info.Balance, release)
				info.Debt = big.NewInt(0)

				log.Trace("show balance", "balance", info.Balance)

				state.SubBalance(vm.RestrictingContractAddr, release)
				state.AddBalance(account, release)

			} else {
				// release amount is more than the sum of balance and debt
				origBalance := info.Balance

				release = release.Sub(release, info.Balance)
				info.Balance = big.NewInt(0)
				info.Debt = info.Debt.Sub(release, info.Debt)
				info.DebtSymbol = true

				state.SubBalance(vm.RestrictingContractAddr, origBalance)
				state.AddBalance(account, origBalance)
			}
		}

		// delete ReleaseAmount
		state.SetState(account, releaseAmountKey, []byte{})

		// delete ReleaseAccount
		state.SetState(vm.RestrictingContractAddr, releaseAccountKey, []byte{})

		// delete epoch in ReleaseList
		// In general, the first epoch is released first.
		// info.ReleaseList = info.ReleaseList[1:]
		for i, target := range info.ReleaseList {
			if target == epoch {
				info.ReleaseList = append(info.ReleaseList[:i], info.ReleaseList[i+1:]...)
				break
			}
		}

		// just restore restricting info, don't delete
		if bNewInfo, err := rlp.EncodeToBytes(info); err != nil {
			log.Error("failed to rlp encode new info while release", "account", account, "info", info)
			return common.NewSysError(err.Error())
		} else {
			state.SetState(account, restrictingKey, bNewInfo)
		}
	}

	// delete ReleaseEpoch
	state.SetState(vm.RestrictingContractAddr, releaseEpochKey, []byte{})

	return nil
}

func (rp *RestrictingPlugin) GetRestrictingInfo(account common.Address, state xcom.StateDB) ([]byte, error) {

	restrictingKey := restricting.GetRestrictingKey(account)
	bAccInfo := state.GetState(account, restrictingKey)

	if len(bAccInfo) == 0 {
		log.Debug("record not found in GetRestrictingInfo", "account", account)
		return []byte{}, errAccountNotFound
	}

	var (
		bAmount          []byte
		info             restricting.RestrictingInfo
		plan             restricting.ReleaseAmountInfo
		plans            []restricting.ReleaseAmountInfo
		releaseAmountKey []byte
		result           restricting.Result
	)

	if err := rlp.Decode(bytes.NewReader(bAccInfo), &info); err != nil {
		log.Error("failed to rlp encode the restricting account", "error", err.Error(), "info", bAccInfo)
		return []byte{}, common.NewSysError(err.Error())
	}

	var amount = new(big.Int)
	for i := 0; i < len(info.ReleaseList); i++ {
		epoch := info.ReleaseList[i]
		releaseAmountKey = restricting.GetReleaseAmountKey(epoch, account)
		bAmount = state.GetState(account, releaseAmountKey)

		plan.Height = GetBlockNumberByEpoch(epoch)
		plan.Amount = amount.SetBytes(bAmount)
		plans = append(plans, plan)
	}

	result.Balance = info.Balance
	result.Debt = info.Debt
	result.Symbol = info.DebtSymbol
	result.Entry = plans
	log.Trace("get restricting result", "account", account, "result", result)

	bResult, err := json.Marshal(result)
	if err != nil {
		log.Error("failed to Marshal restricting result")
		return []byte{}, err
	}

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

	if len(bEpoch) == 0 {
		return 0
	} else {
		return common.BytesToUint64(bEpoch)
	}
}

func GetBlockNumberByEpoch(epoch uint64) uint64 {
	return epoch * xutil.CalcBlocksEachEpoch()
}
