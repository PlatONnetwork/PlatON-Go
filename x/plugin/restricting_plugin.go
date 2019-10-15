package plugin

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

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
	restrictTxPlanSize                   = 36 // The plans num size of restricting tx
	errParamEpochInvalid                 = common.NewBizError(304001, "param epoch can't be zero")
	errCountRestrictPlansInvalid         = common.NewBizError(304002, "the number of the restricting plan can't be zero or more than 36")
	errLockedAmountTooLess               = common.NewBizError(304003, "total restricting amount need more than 1 LAT")
	errBalanceNotEnough                  = common.NewBizError(304004, "create plan,the sender balance is not enough in restrict")
	errAccountNotFound                   = common.NewBizError(304005, "account is not found on restricting contract")
	errSlashingTooMuch                   = common.NewBizError(304006, "slashing amount is larger than staking amount")
	errStakingAmountEmpty                = common.NewBizError(304007, "staking amount is 0")
	errPledgeLockFundsAmountLessThanZero = common.NewBizError(304008, "pledge lock funds amount can't less than 0")
	errReturnLockFundsAmountLessThanZero = common.NewBizError(304009, "return lock funds amount can't less than 0")
	errSlashingAmountLessThanZero        = common.NewBizError(304010, "slashing amount can't less than 0")
	errCreatePlanAmountLessThanZero      = common.NewBizError(304011, "create plan each amount can't less than 0")
	errStakingAmountInvalid              = common.NewBizError(304012, "staking return amount is wrong")
	errRestrictBalanceNotEnough          = common.NewBizError(304013, "the user restricting balance is not enough for pledge lock funds")
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
		restrictLog := log.Root().New("package", "RestrictingPlugin")
		restrictLog.Info("Init Restricting plugin ...")
		rt = &RestrictingPlugin{restrictLog}
	})
	return rt
}

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

	rp.log.Info("begin to release restricting plan", "currentHash", blockHash, "currBlock", head.Number, "expectBlock", expectBlock, "expectEpoch", expect)
	if err := rp.releaseRestricting(expect, state); err != nil {
		return err
	}
	SetLatestEpoch(state, expect)
	return nil
}

// Confirmed is empty function
func (rp *RestrictingPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}

func (rp *RestrictingPlugin) mergeAmount(state xcom.StateDB, plans []restricting.RestrictingPlan) (*big.Int, map[uint64]*big.Int, error) {
	// latest is the epoch of a settlement block closest to current block
	latestEpoch := GetLatestEpoch(state)

	totalAmount := new(big.Int)

	planMap := make(map[uint64]*big.Int, restrictTxPlanSize)
	for _, plan := range plans {
		epoch, amount := plan.Epoch, new(big.Int).Set(plan.Amount)
		if epoch == 0 {
			rp.log.Error(errParamEpochInvalid.Error())
			return nil, nil, errParamEpochInvalid
		}
		if amount.Cmp(common.Big0) <= 0 {
			rp.log.Error("Failed to mergeAmount for plans on restricting RestrictingPlugin: the amount must be more than zero", "epoch", epoch, "amount", amount)
			return nil, nil, errCreatePlanAmountLessThanZero
		}
		totalAmount.Add(totalAmount, amount)
		newEpoch := epoch + latestEpoch
		if value, ok := planMap[newEpoch]; ok {
			planMap[newEpoch] = value.Add(amount, value)
		} else {
			planMap[newEpoch] = amount
		}
	}
	return totalAmount, planMap, nil
}

func (rp *RestrictingPlugin) initEpochInfo(state xcom.StateDB, epoch uint64, account common.Address, amount *big.Int) {
	// step1: get account numbers at target epoch
	releaseEpochKey, lastEpochAccountIndex := rp.getReleaseEpochNumber(state, epoch)
	newEpochAccountIndex := lastEpochAccountIndex + 1

	// step2: save the numbers of restricting record at target epoch
	rp.storeNumber2ReleaseEpoch(state, releaseEpochKey, newEpochAccountIndex)

	// step3: save account at target index

	rp.storeAccount2ReleaseAccount(state, epoch, newEpochAccountIndex, account)

	// step4: save restricting amount at target epoch
	rp.storeAmount2ReleaseAmount(state, epoch, account, amount)
}

func (rp *RestrictingPlugin) transferAmount(state xcom.StateDB, from, to common.Address, mount *big.Int) {
	state.SubBalance(from, mount)
	state.AddBalance(to, mount)
}

// AddRestrictingRecord stores four K-V record in StateDB:
// RestrictingInfo: the account info to be released
// ReleaseEpoch:   the number of accounts to be released on the epoch corresponding to the target block height
// ReleaseAccount: the account on the index on the target epoch
// ReleaseAmount: the amount of the account to be released on the target epoch
func (rp *RestrictingPlugin) AddRestrictingRecord(from, account common.Address, plans []restricting.RestrictingPlan, state xcom.StateDB) error {

	rp.log.Info("Call AddRestrictingRecord begin", "sender", from.String(), "account", account.String(), "plans", fmt.Sprintf("%+v", plans))

	if len(plans) == 0 || len(plans) > restrictTxPlanSize {
		rp.log.Error(fmt.Sprintf("Failed to AddRestrictingRecord: the number of restricting plan %d can't be zero or more than %d", len(plans), restrictTxPlanSize))
		return errCountRestrictPlansInvalid
	}
	// totalAmount is total restricting amount
	totalAmount, mPlans, err := rp.mergeAmount(state, plans)
	if err != nil {
		return err
	}
	// pre-check
	{

		if totalAmount.Cmp(big.NewInt(1e18)) < 0 {
			rp.log.Error("Failed to AddRestrictingRecord: total restricting amount need more than 1 LAT", "from", from, "amount", totalAmount)
			return errLockedAmountTooLess
		}

		if state.GetBalance(from).Cmp(totalAmount) < 0 {
			rp.log.Error("Failed to AddRestrictingRecord: balance of the sender is not enough", "total", totalAmount, "balance", state.GetBalance(from))
			return errBalanceNotEnough
		}
	}

	var (
		epochList []uint64
		info      restricting.RestrictingInfo
	)

	rp.transferAmount(state, from, vm.RestrictingContractAddr, totalAmount)

	restrictingKey, bAccInfo := rp.getRestrictingInfo(state, account)
	if len(bAccInfo) == 0 {
		rp.log.Info("restricting record not exist", "account", account.String())
		for epoch, amount := range mPlans {
			rp.initEpochInfo(state, epoch, account, amount)
			epochList = append(epochList, epoch)
		}
		info.CachePlanAmount = totalAmount
		info.NeedRelease = big.NewInt(0)
		info.StakingAmount = big.NewInt(0)
		info.ReleaseList = epochList
	} else {
		rp.log.Info("restricting record exist", "account", account.String())
		if err = rlp.DecodeBytes(bAccInfo, &info); err != nil {
			rp.log.Error("failed to rlp decode the restricting account", "err", err.Error())
			return common.InternalError.Wrap(err.Error())
		}
		if info.NeedRelease.Cmp(common.Big0) > 0 {
			if info.NeedRelease.Cmp(totalAmount) >= 0 {
				info.NeedRelease.Sub(info.NeedRelease, totalAmount)
				rp.transferAmount(state, vm.RestrictingContractAddr, account, totalAmount)
			} else {
				rp.transferAmount(state, vm.RestrictingContractAddr, account, info.NeedRelease)
				totalAmount.Sub(totalAmount, info.NeedRelease)
				info.CachePlanAmount.Add(info.CachePlanAmount, totalAmount)
				info.NeedRelease = new(big.Int).SetInt64(0)
			}
		} else {
			info.CachePlanAmount.Add(info.CachePlanAmount, totalAmount)
		}
		for epoch, releaseAmount := range mPlans {
			// step1: get restricting amount at target epoch
			_, currentAmount := rp.getReleaseAmount(state, epoch, account)
			if currentAmount.Cmp(common.Big0) == 0 {
				rp.log.Trace("release record not exist on curr epoch ", "account", account.String(), "epoch", epoch)
				rp.initEpochInfo(state, epoch, account, releaseAmount)
				info.ReleaseList = append(info.ReleaseList, epoch)
			} else {
				rp.log.Trace("release record exist at curr epoch", "account", account.String(), "epoch", epoch)
				currentAmount.Add(currentAmount, releaseAmount)
				// step4: save restricting amount at target epoch
				rp.storeAmount2ReleaseAmount(state, epoch, account, currentAmount)
			}
		}
	}

	// sort release list
	sort.Slice(info.ReleaseList, func(i, j int) bool {
		return info.ReleaseList[i] < info.ReleaseList[j]
	})
	rp.storeRestrictingInfo(state, restrictingKey, info)
	rp.log.Debug("Call AddRestrictingRecord finished", "account", account.String(), "restrictingInfo", fmt.Sprintf("%+v", info))

	return nil
}

// PledgeLockFunds transfer the money from the restricting contract account to the staking contract account
func (rp *RestrictingPlugin) PledgeLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {

	restrictingKey, info, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return err
	}
	rp.log.Debug("Call PledgeLockFunds begin", "account", account.String(), "amount", amount, "old info", fmt.Sprintf("%+v", info))

	if amount.Cmp(common.Big0) < 0 {
		return errPledgeLockFundsAmountLessThanZero
	} else if amount.Cmp(common.Big0) == 0 {
		return nil
	}

	canStaking := new(big.Int).Sub(info.CachePlanAmount, info.StakingAmount)
	if canStaking.Cmp(amount) < 0 {
		rp.log.Warn("Balance of restricting account not enough", "total", info.CachePlanAmount, "stanking", info.StakingAmount, "funds", amount)
		return errRestrictBalanceNotEnough
	}

	// sub Balance
	info.StakingAmount.Add(info.StakingAmount, amount)

	// save restricting account info
	rp.storeRestrictingInfo(state, restrictingKey, info)
	rp.transferAmount(state, vm.RestrictingContractAddr, vm.StakingContractAddr, amount)

	rp.log.Info("Call PledgeLockFunds finished", "RestrictingContractBalance", state.GetBalance(vm.RestrictingContractAddr), "StakingContractBalance", state.GetBalance(vm.StakingContractAddr), "new info", fmt.Sprintf("%+v", info))
	return nil
}

// ReturnLockFunds transfer the money from the staking contract account to the restricting contract account
func (rp *RestrictingPlugin) ReturnLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {
	amountCompareWithZero := amount.Cmp(common.Big0)
	if amountCompareWithZero == 0 {
		return nil
	} else if amountCompareWithZero < 0 {
		return errReturnLockFundsAmountLessThanZero
	}
	restrictingKey, info, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return err
	}
	rp.log.Info("Call ReturnLockFunds begin", "account", account.String(), "amount", amount, "info", fmt.Sprintf("%+v", info))

	if info.StakingAmount.Cmp(amount) < 0 {
		return errStakingAmountInvalid
	}

	rp.transferAmount(state, vm.StakingContractAddr, vm.RestrictingContractAddr, amount)
	if info.NeedRelease.Cmp(common.Big0) > 0 {
		if info.NeedRelease.Cmp(amount) >= 0 {
			info.NeedRelease.Sub(info.NeedRelease, amount)
			info.CachePlanAmount.Sub(info.CachePlanAmount, amount)
			rp.transferAmount(state, vm.RestrictingContractAddr, account, amount)
		} else {
			rp.transferAmount(state, vm.RestrictingContractAddr, account, info.NeedRelease)
			tmp := new(big.Int).Sub(amount, info.NeedRelease)
			info.CachePlanAmount.Add(info.CachePlanAmount, tmp)
			info.NeedRelease = big.NewInt(0)
		}
	}
	info.StakingAmount.Sub(info.StakingAmount, amount)
	// save restricting account info
	if info.NeedRelease.Cmp(common.Big0) == 0 && info.StakingAmount.Cmp(common.Big0) == 0 && len(info.ReleaseList) == 0 && info.CachePlanAmount.Cmp(common.Big0) == 0 {
		state.SetState(vm.RestrictingContractAddr, restrictingKey, []byte{})
	} else {
		rp.storeRestrictingInfo(state, restrictingKey, info)
	}
	rp.log.Info("Call ReturnLockFunds finished", "RCContractBalance", state.GetBalance(vm.RestrictingContractAddr), "info", fmt.Sprintf("%+v", info))
	return nil
}

// SlashingNotify modify Debt of restricting account
func (rp *RestrictingPlugin) SlashingNotify(account common.Address, amount *big.Int, state xcom.StateDB) error {
	rp.log.Info("Call SlashingNotify begin", "account", account.String(), "amount", amount)

	restrictingKey, info, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return err
	}
	if amount.Cmp(common.Big0) < 0 {
		return errSlashingAmountLessThanZero
	} else if amount.Cmp(common.Big0) == 0 {
		return nil
	}
	if info.StakingAmount.Cmp(common.Big0) <= 0 {
		rp.log.Error("Failed to SlashingNotify", "account", account.String(), "Debt", info.StakingAmount, "slashing", amount, "err", errStakingAmountEmpty.Error())
		return errStakingAmountEmpty
	}

	if info.StakingAmount.Cmp(amount) < 0 {
		return errSlashingTooMuch
	}
	info.StakingAmount.Sub(info.StakingAmount, amount)
	info.CachePlanAmount.Sub(info.CachePlanAmount, amount)

	rp.storeRestrictingInfo(state, restrictingKey, info)

	// save restricting account info
	rp.log.Info("Call SlashingNotify finished", "restrictingInfo", fmt.Sprintf("%+v", info))

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
	restrictingKey, accInfoByte := rp.getRestrictingInfo(state, account)
	if len(accInfoByte) == 0 {
		rp.log.Warn("record not found in GetRestrictingInfo", "account", account.String())
		return []byte{}, info, errAccountNotFound
	}
	if err := rlp.DecodeBytes(accInfoByte, &info); err != nil {
		rp.log.Error("Failed to rlp decode restricting account", "error", err.Error())
		return restrictingKey, info, common.InternalError.Wrap(err.Error())
	}
	return restrictingKey, info, nil
}

func (rp *RestrictingPlugin) getRestrictingInfoByDecode(state xcom.StateDB, account common.Address) ([]byte, restricting.RestrictingInfo, error) {
	restrictingKey, bAccInfo := rp.getRestrictingInfo(state, account)
	var info restricting.RestrictingInfo
	if err := rlp.DecodeBytes(bAccInfo, &info); err != nil {
		rp.log.Error("Failed to rlp decode restricting account", "error", err.Error(), "account", account.String())
		return restrictingKey, info, common.InternalError.Wrap(err.Error())
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

func (rp *RestrictingPlugin) storeRestrictingInfo(state xcom.StateDB, restrictingKey []byte, info restricting.RestrictingInfo) {
	bNewInfo, err := rlp.EncodeToBytes(info)
	if err != nil {
		rp.log.Error("Failed to rlp encode restricting info", "error", err, "info", fmt.Sprintf("%+v", info))
		panic(err)
	}
	state.SetState(vm.RestrictingContractAddr, restrictingKey, bNewInfo)
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

// releaseRestricting will release restricting plans on target epoch
func (rp *RestrictingPlugin) releaseRestricting(epoch uint64, state xcom.StateDB) error {

	rp.log.Info("Call releaseRestricting begin", "epoch", epoch)
	releaseEpochKey, numbers := rp.getReleaseEpochNumber(state, epoch)
	if numbers == 0 {
		rp.log.Info("Call releaseRestricting: there is no release record on curr epoch", "epoch", epoch)
		return nil
	}

	rp.log.Info("Call releaseRestricting: many restricting records need release", "epoch", epoch, "records", numbers)

	for index := numbers; index > 0; index-- {
		releaseAccountKey, account := rp.getReleaseAccount(state, epoch, index)

		restrictingKey, info, err := rp.getRestrictingInfoByDecode(state, account)
		if err != nil {
			return err
		}

		releaseAmountKey, releaseAmount := rp.getReleaseAmount(state, epoch, account)
		rp.log.Debug("Call releaseRestricting: begin to release record", "index", index, "account", account.String(), "info", fmt.Sprintf("%+v", info), "releaseAmount", releaseAmount)

		if info.NeedRelease.Cmp(common.Big0) > 0 {
			//info.CachePlanAmount.Sub(info.CachePlanAmount, releaseAmount)
			if info.CachePlanAmount.Cmp(common.Big0) == 0 {
				info.NeedRelease.Sub(info.NeedRelease, releaseAmount)
			} else {
				info.NeedRelease.Add(info.NeedRelease, releaseAmount)
			}
		} else {
			canRelease := new(big.Int).Sub(info.CachePlanAmount, info.StakingAmount)
			if canRelease.Cmp(releaseAmount) >= 0 {
				rp.transferAmount(state, vm.RestrictingContractAddr, account, releaseAmount)
				info.CachePlanAmount.Sub(info.CachePlanAmount, releaseAmount)
			} else {
				needRelease := new(big.Int).Sub(releaseAmount, canRelease)
				rp.transferAmount(state, vm.RestrictingContractAddr, account, canRelease)
				info.NeedRelease.Add(info.NeedRelease, needRelease)
				info.CachePlanAmount.Sub(info.CachePlanAmount, canRelease)
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

		if info.CachePlanAmount.Cmp(common.Big0) == 0 {
			if info.NeedRelease.Cmp(common.Big0) == 0 || len(info.ReleaseList) == 0 {
				//if all is release,remove info
				state.SetState(vm.RestrictingContractAddr, restrictingKey, []byte{})
			} else {
				rp.storeRestrictingInfo(state, restrictingKey, info)
			}
		} else {
			rp.storeRestrictingInfo(state, restrictingKey, info)
		}
	}

	// delete ReleaseEpoch
	state.SetState(vm.RestrictingContractAddr, releaseEpochKey, []byte{})

	rp.log.Info("Call releaseRestricting finished")

	return nil
}

func (rp *RestrictingPlugin) getRestrictingInfoToReturn(account common.Address, state xcom.StateDB) (restricting.Result, error) {
	_, info, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return restricting.Result{}, err
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
		plan.Amount = (*hexutil.Big)(bAmount)
		plans = append(plans, plan)
	}

	result.Balance = (*hexutil.Big)(info.CachePlanAmount)
	result.Debt = (*hexutil.Big)(info.NeedRelease)
	result.Entry = plans
	result.Pledge = (*hexutil.Big)(info.StakingAmount)
	rp.log.Info("Call releaseRestricting: query restricting result", "account", account.String(), "result", result)
	return result, nil
}

func (rp *RestrictingPlugin) GetRestrictingInfo(account common.Address, state xcom.StateDB) ([]byte, error) {
	result, err := rp.getRestrictingInfoToReturn(account, state)
	if err != nil {
		return nil, err
	}
	bResult, err := json.Marshal(result)
	if err != nil {
		rp.log.Error("Failed to Marshal restricting result")
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
	return common.BytesToUint64(bEpoch)
}

func GetBlockNumberByEpoch(epoch uint64) uint64 {
	return epoch * xutil.CalcBlocksEachEpoch()
}
