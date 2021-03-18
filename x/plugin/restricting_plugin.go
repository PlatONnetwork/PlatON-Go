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
	"sort"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

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

type RestrictingPlugin struct {
	log log.Logger
	db  snapshotdb.DB
}

var (
	restrictingOnce sync.Once
	rt              *RestrictingPlugin
)

func RestrictingInstance() *RestrictingPlugin {
	restrictingOnce.Do(func() {
		restrictLog := log.Root().New("package", "RestrictingPlugin")
		restrictLog.Info("Init Restricting plugin ...")
		rt = &RestrictingPlugin{restrictLog, snapshotdb.Instance()}
	})
	return rt
}

func NewRestrictingPlugin(snapdb snapshotdb.DB) *RestrictingPlugin {
	restrictLog := log.Root().New("package", "RestrictingPlugin")
	return &RestrictingPlugin{restrictLog, snapdb}
}

// BeginBlock does something like check input params before execute transactions,
// in RestrictingPlugin it does nothing.
func (rp *RestrictingPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	return nil
}

// EndBlock invoke releaseRestricting
func (rp *RestrictingPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	if xutil.IsEndOfEpoch(head.Number.Uint64()) {
		expect := xutil.CalculateEpoch(head.Number.Uint64())
		rp.log.Info("begin to release restricting plan", "currentHash", blockHash, "currBlock", head.Number, "expectBlock", head.Number, "expectEpoch", expect)
		if err := rp.releaseRestricting(expect, state); err != nil {
			return err
		}
		if ok, _ := xcom.IsYearEnd(blockHash, head.Number.Uint64()); ok {
			rp.log.Info(fmt.Sprintf("release genesis restricting plan, blocknumber:%d", head.Number.Uint64()))
			return rp.releaseGenesisRestrictingPlans(blockHash, state)
		}
	}
	return nil
}

// Confirmed is empty function
func (rp *RestrictingPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}

func (rp *RestrictingPlugin) mergeAmount(state xcom.StateDB, blockNum uint64, blockHash common.Hash, plans []restricting.RestrictingPlan) (*big.Int, map[uint64]*big.Int, error) {
	// latest is the epoch of a settlement block closest to current block
	latestEpoch := xutil.CalculateEpoch(blockNum)

	totalAmount := new(big.Int)

	planMap := make(map[uint64]*big.Int, restricting.RestrictTxPlanSize)

	minimumAmount, err := gov.GovernRestrictingMinimumAmount(blockNum, blockHash)
	if err != nil {
		return nil, nil, err
	}

	for _, plan := range plans {
		epoch, amount := plan.Epoch, new(big.Int).Set(plan.Amount)
		if epoch == 0 {
			rp.log.Error(restricting.ErrParamEpochInvalid.Error())
			return nil, nil, restricting.ErrParamEpochInvalid
		}
		if amount.Cmp(common.Big0) <= 0 {
			rp.log.Error("Failed to mergeAmount for plans on restricting RestrictingPlugin: the amount must be more than zero", "epoch", epoch, "amount", amount)
			return nil, nil, restricting.ErrCreatePlanAmountLessThanZero
		}
		if amount.Cmp(minimumAmount) < 0 {
			rp.log.Error("Failed to mergeAmount for plans on restricting RestrictingPlugin: the amount must be more than minimumAmount", "epoch", epoch, "amount", amount, "mini", minimumAmount)
			return nil, nil, restricting.ErrCreatePlanAmountLessThanMiniAmount
		}
		totalAmount.Add(totalAmount, amount)
		newEpoch := epoch + latestEpoch - 1
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

// update genesis restricting plans
func (rp *RestrictingPlugin) updateGenesisRestrictingPlans(plans []*big.Int, stateDB xcom.StateDB) error {

	if val, err := rlp.EncodeToBytes(plans); nil != err {
		return fmt.Errorf("Failed to Store genesisAllowancePlans Info: rlp encodeing failed")
	} else {
		stateDB.SetState(vm.RestrictingContractAddr, restricting.InitialFoundationRestricting, val)
	}
	return nil
}

// init the genesis restricting plans
func (rp *RestrictingPlugin) InitGenesisRestrictingPlans(statedb xcom.StateDB) error {

	genesisAllowancePlans := []*big.Int{
		new(big.Int).Mul(big.NewInt(55965742), big.NewInt(1e18)),
		new(big.Int).Mul(big.NewInt(49559492), big.NewInt(1e18)),
		new(big.Int).Mul(big.NewInt(42993086), big.NewInt(1e18)),
		new(big.Int).Mul(big.NewInt(36262520), big.NewInt(1e18)),
		new(big.Int).Mul(big.NewInt(29363689), big.NewInt(1e18)),
		new(big.Int).Mul(big.NewInt(22292388), big.NewInt(1e18)),
		new(big.Int).Mul(big.NewInt(15044304), big.NewInt(1e18)),
		new(big.Int).Mul(big.NewInt(7615018), big.NewInt(1e18)),
	}

	//initial release from genesis restricting plans(62215742LAT)
	initialRelease := new(big.Int).Mul(big.NewInt(62215742), big.NewInt(1e18))
	statedb.SubBalance(xcom.CDFAccount(), initialRelease)
	statedb.AddBalance(vm.RewardManagerPoolAddr, initialRelease)

	//transfer 259096239LAT from CDFAccount to vm.RestrictingContractAddr
	totalRestrictingPlan := new(big.Int).Mul(big.NewInt(259096239), big.NewInt(1e18))
	statedb.SubBalance(xcom.CDFAccount(), totalRestrictingPlan)
	statedb.AddBalance(vm.RestrictingContractAddr, totalRestrictingPlan)

	if err := rp.updateGenesisRestrictingPlans(genesisAllowancePlans, statedb); nil != err {
		return err
	}
	return nil
}

// release genesis restricting plans
func (rp *RestrictingPlugin) releaseGenesisRestrictingPlans(blockHash common.Hash, statedb xcom.StateDB) error {

	plansBytes := statedb.GetState(vm.RestrictingContractAddr, restricting.InitialFoundationRestricting)
	var genesisAllowancePlans []*big.Int
	if len(plansBytes) > 0 {
		if err := rlp.DecodeBytes(plansBytes, &genesisAllowancePlans); err != nil {
			rp.log.Error("failed to rlp decode the genesis allowance plans", "err", err.Error())
			return common.InternalError.Wrap(err.Error())
		} else {
			remains := len(genesisAllowancePlans)
			if remains > 0 {
				allowance := genesisAllowancePlans[0]
				statedb.SubBalance(vm.RestrictingContractAddr, allowance)
				statedb.AddBalance(vm.RewardManagerPoolAddr, allowance)
				rp.log.Info("Genesis restricting plan release", "remains", remains, "allowance", allowance)
				genesisAllowancePlans = append(genesisAllowancePlans[:0], genesisAllowancePlans[1:]...)
				if err := rp.updateGenesisRestrictingPlans(genesisAllowancePlans, statedb); nil != err {
					return err
				}
			} else {
				statedb.SetState(vm.RestrictingContractAddr, restricting.InitialFoundationRestricting, []byte{})
			}
			rp.log.Info("release genesis restricting plan", "remains:", remains, "left:", len(genesisAllowancePlans))
		}
	} else {
		rp.log.Info("Genesis restricting plan had all been released")
	}

	return nil
}

// AddRestrictingRecord stores four K-V record in StateDB:
// RestrictingInfo: the account info to be released
// ReleaseEpoch:   the number of accounts to be released on the epoch corresponding to the target block height
// ReleaseAccount: the account on the index on the target epoch
// ReleaseAmount: the amount of the account to be released on the target epoch
func (rp *RestrictingPlugin) AddRestrictingRecord(from, account common.Address, blockNum uint64, blockHash common.Hash, plans []restricting.RestrictingPlan, state xcom.StateDB, txhash common.Hash) error {

	rp.log.Debug("Call AddRestrictingRecord begin", "sender", from, "account", account, "plans", plans)

	if len(plans) == 0 || len(plans) > restricting.RestrictTxPlanSize {
		rp.log.Error(fmt.Sprintf("Failed to AddRestrictingRecord: the number of restricting plan %d can't be zero or more than %d",
			len(plans), restricting.RestrictTxPlanSize))
		return restricting.ErrCountRestrictPlansInvalid
	}
	// totalAmount is total restricting amount
	totalAmount, totalPlans, err := rp.mergeAmount(state, blockNum, blockHash, plans)
	if err != nil {
		return err
	}
	// pre-check
	{

		if totalAmount.Cmp(big.NewInt(1e18)) < 0 {
			rp.log.Error("Failed to AddRestrictingRecord: total restricting amount need more than 1 LAT",
				"from", from, "amount", totalAmount)
			return restricting.ErrLockedAmountTooLess
		}

		if state.GetBalance(from).Cmp(totalAmount) < 0 {
			rp.log.Error("Failed to AddRestrictingRecord: balance of the sender is not enough",
				"total", totalAmount, "balance", state.GetBalance(from))
			return restricting.ErrBalanceNotEnough
		}
	}
	if txhash == common.ZeroHash {
		return nil
	}
	var (
		epochArr     []uint64
		restrictInfo restricting.RestrictingInfo
	)

	rp.transferAmount(state, from, vm.RestrictingContractAddr, totalAmount)

	restrictingKey, restrictInfoByte := rp.getRestrictingInfo(state, account)
	if len(restrictInfoByte) == 0 {
		rp.log.Trace("restricting record not exist", "account", account.String())
		for epoch, amount := range totalPlans {
			rp.initEpochInfo(state, epoch, account, amount)
			epochArr = append(epochArr, epoch)
		}
		restrictInfo.CachePlanAmount = totalAmount
		restrictInfo.NeedRelease = big.NewInt(0)
		restrictInfo.AdvanceAmount = big.NewInt(0)
		restrictInfo.ReleaseList = epochArr
	} else {
		rp.log.Trace("restricting record exist", "account", account.String())
		if err = rlp.DecodeBytes(restrictInfoByte, &restrictInfo); err != nil {
			rp.log.Error("failed to rlp decode the restricting account", "err", err.Error())
			return common.InternalError.Wrap(err.Error())
		}
		if restrictInfo.NeedRelease.Cmp(common.Big0) > 0 {
			if restrictInfo.NeedRelease.Cmp(totalAmount) >= 0 {
				restrictInfo.NeedRelease.Sub(restrictInfo.NeedRelease, totalAmount)
				rp.transferAmount(state, vm.RestrictingContractAddr, account, totalAmount)
			} else {
				rp.transferAmount(state, vm.RestrictingContractAddr, account, restrictInfo.NeedRelease)
				totalAmount.Sub(totalAmount, restrictInfo.NeedRelease)
				restrictInfo.CachePlanAmount.Add(restrictInfo.CachePlanAmount, totalAmount)
				restrictInfo.NeedRelease = new(big.Int).SetInt64(0)
			}
		} else {
			restrictInfo.CachePlanAmount.Add(restrictInfo.CachePlanAmount, totalAmount)
		}
		for epoch, releaseAmount := range totalPlans {
			// step1: get restricting amount at target epoch
			_, currentAmount := rp.getReleaseAmount(state, epoch, account)
			if currentAmount.Cmp(common.Big0) == 0 {
				rp.log.Trace("release record not exist on curr epoch ", "account", account, "epoch", epoch)
				rp.initEpochInfo(state, epoch, account, releaseAmount)
				restrictInfo.ReleaseList = append(restrictInfo.ReleaseList, epoch)
			} else {
				rp.log.Trace("release record exist at curr epoch", "account", account, "epoch", epoch)
				currentAmount.Add(currentAmount, releaseAmount)
				// step4: save restricting amount at target epoch
				rp.storeAmount2ReleaseAmount(state, epoch, account, currentAmount)
			}
		}
	}

	// sort release list
	sort.Slice(restrictInfo.ReleaseList, func(i, j int) bool {
		return restrictInfo.ReleaseList[i] < restrictInfo.ReleaseList[j]
	})
	rp.storeRestrictingInfo(state, restrictingKey, restrictInfo)
	rp.log.Debug("Call AddRestrictingRecord finished", "account", account, "restrictingInfo", restrictInfo)

	return nil
}

// AdvanceLockedFunds transfer the money from the restricting contract account to the staking contract account
func (rp *RestrictingPlugin) AdvanceLockedFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {

	restrictingKey, restrictInfo, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return err
	}
	rp.log.Debug("Call AdvanceLockedFunds begin", "account", account, "amount", amount, "old info", restrictInfo)

	if amount.Cmp(common.Big0) < 0 {
		return restricting.ErrPledgeLockFundsAmountLessThanZero
	} else if amount.Cmp(common.Big0) == 0 {
		return nil
	}

	canStaking := new(big.Int).Sub(restrictInfo.CachePlanAmount, restrictInfo.AdvanceAmount)
	if canStaking.Cmp(amount) < 0 {
		rp.log.Warn("Balance of restricting account not enough", "totalAmount",
			restrictInfo.CachePlanAmount, "stankingAmount", restrictInfo.AdvanceAmount, "funds", amount)
		return restricting.ErrRestrictBalanceNotEnough
	}

	// sub Balance
	restrictInfo.AdvanceAmount.Add(restrictInfo.AdvanceAmount, amount)

	// save restricting account info
	rp.storeRestrictingInfo(state, restrictingKey, restrictInfo)
	rp.transferAmount(state, vm.RestrictingContractAddr, vm.StakingContractAddr, amount)

	rp.log.Debug("Call AdvanceLockedFunds finished", "RestrictingContractBalance", state.GetBalance(vm.RestrictingContractAddr), "StakingContractBalance", state.GetBalance(vm.StakingContractAddr), "new info", restrictInfo)
	return nil
}

// MixAdvanceLockedFunds transfer the money from the restricting contract account to the staking contract account,use restricting von first,if restricting not en
func (rp *RestrictingPlugin) MixAdvanceLockedFunds(account common.Address, amount *big.Int, state xcom.StateDB) (*big.Int, *big.Int, error) {

	restrictingKey, restrictInfo, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		if err == restricting.ErrAccountNotFound {
			//if not found restricting,we just use free
			origin := state.GetBalance(account)
			if origin.Cmp(amount) < 0 {
				return nil, nil, staking.ErrAccountVonNoEnough
			}
			rp.transferAmount(state, account, vm.StakingContractAddr, amount)
			return new(big.Int), new(big.Int).Set(amount), nil
		} else {
			return nil, nil, err
		}
	}
	rp.log.Debug("Call MixAdvanceLockedFunds begin", "account", account, "amount", amount, "old info", restrictInfo)

	if amount.Cmp(common.Big0) < 0 {
		return nil, nil, restricting.ErrPledgeLockFundsAmountLessThanZero
	} else if amount.Cmp(common.Big0) == 0 {
		return amount, amount, nil
	}

	canStakingRestricting := new(big.Int).Sub(restrictInfo.CachePlanAmount, restrictInfo.AdvanceAmount)

	canStakingFree := state.GetBalance(account)

	total := new(big.Int).Add(canStakingRestricting, canStakingFree)

	if total.Cmp(amount) < 0 {
		rp.log.Warn("Balance of restricting and free not enough", "totalAmount",
			restrictInfo.CachePlanAmount, "stankingAmount", restrictInfo.AdvanceAmount, "free", canStakingFree, "funds", amount)
		return nil, nil, restricting.ErrRestrictBalanceAndFreeNotEnough
	}

	forRestricting := new(big.Int).Set(amount)
	forFree := new(big.Int)
	if canStakingRestricting.Cmp(amount) < 0 {
		forRestricting.Set(canStakingRestricting)
		forFree = new(big.Int).Sub(amount, canStakingRestricting)
		rp.transferAmount(state, account, vm.StakingContractAddr, forFree)
	}

	restrictInfo.AdvanceAmount.Add(restrictInfo.AdvanceAmount, forRestricting)
	// save restricting account info
	rp.storeRestrictingInfo(state, restrictingKey, restrictInfo)
	rp.transferAmount(state, vm.RestrictingContractAddr, vm.StakingContractAddr, forRestricting)

	rp.log.Debug("Call mixPledgeLockFunds finished", "RestrictingContractBalance", state.GetBalance(vm.RestrictingContractAddr), "StakingContractBalance", state.GetBalance(vm.StakingContractAddr), "new info", restrictInfo, "for free", forFree, "for restricting", forRestricting)
	return forRestricting, forFree, nil
}

// ReturnLockFunds transfer the money from the staking contract account to the restricting contract account
func (rp *RestrictingPlugin) ReturnLockFunds(account common.Address, amount *big.Int, state xcom.StateDB) error {
	amountCompareWithZero := amount.Cmp(common.Big0)
	if amountCompareWithZero == 0 {
		return nil
	} else if amountCompareWithZero < 0 {
		return restricting.ErrReturnLockFundsAmountLessThanZero
	}
	restrictingKey, restrictInfo, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return err
	}
	rp.log.Debug("Call ReturnLockFunds begin", "account", account, "amount", amount, "info", restrictInfo)

	if restrictInfo.AdvanceAmount.Cmp(amount) < 0 {
		return restricting.ErrStakingAmountInvalid
	}

	rp.transferAmount(state, vm.StakingContractAddr, vm.RestrictingContractAddr, amount)
	if restrictInfo.NeedRelease.Cmp(common.Big0) > 0 {
		if restrictInfo.NeedRelease.Cmp(amount) >= 0 {
			restrictInfo.NeedRelease.Sub(restrictInfo.NeedRelease, amount)
			restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, amount)
			rp.transferAmount(state, vm.RestrictingContractAddr, account, amount)
		} else {
			rp.transferAmount(state, vm.RestrictingContractAddr, account, restrictInfo.NeedRelease)
			restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, restrictInfo.NeedRelease)
			restrictInfo.NeedRelease = big.NewInt(0)
		}
	}
	restrictInfo.AdvanceAmount.Sub(restrictInfo.AdvanceAmount, amount)
	// save restricting account info
	if restrictInfo.AdvanceAmount.Cmp(common.Big0) == 0 &&
		len(restrictInfo.ReleaseList) == 0 && restrictInfo.CachePlanAmount.Cmp(common.Big0) == 0 {
		state.SetState(vm.RestrictingContractAddr, restrictingKey, []byte{})
		rp.log.Debug("Call ReturnLockFunds finished,set info empty", "RCContractBalance", state.GetBalance(vm.RestrictingContractAddr))
	} else {
		rp.storeRestrictingInfo(state, restrictingKey, restrictInfo)
		rp.log.Debug("Call ReturnLockFunds finished", "RCContractBalance", state.GetBalance(vm.RestrictingContractAddr), "info", restrictInfo)
	}
	return nil
}

// SlashingNotify modify Debt of restricting account
func (rp *RestrictingPlugin) SlashingNotify(account common.Address, amount *big.Int, state xcom.StateDB) error {

	restrictingKey, restrictInfo, err := rp.mustGetRestrictingInfoByDecode(state, account)
	if err != nil {
		return err
	}
	if amount.Cmp(common.Big0) < 0 {
		return restricting.ErrSlashingAmountLessThanZero
	} else if amount.Cmp(common.Big0) == 0 {
		return nil
	}
	if restrictInfo.AdvanceAmount.Cmp(common.Big0) <= 0 {
		rp.log.Error("Failed to SlashingNotify", "account", account, "Debt", restrictInfo.AdvanceAmount,
			"slashing", amount, "err", restricting.ErrStakingAmountEmpty.Error())
		return restricting.ErrStakingAmountEmpty
	}

	if restrictInfo.AdvanceAmount.Cmp(amount) < 0 {
		return restricting.ErrSlashingTooMuch
	}
	restrictInfo.AdvanceAmount.Sub(restrictInfo.AdvanceAmount, amount)
	restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, amount)

	if restrictInfo.AdvanceAmount.Cmp(common.Big0) == 0 &&
		len(restrictInfo.ReleaseList) == 0 && restrictInfo.CachePlanAmount.Cmp(common.Big0) == 0 {
		state.SetState(vm.RestrictingContractAddr, restrictingKey, []byte{})
		// save restricting account info
		rp.log.Debug("Call SlashingNotify finished,set empty info", "account", account, "amount", amount)
	} else {
		rp.storeRestrictingInfo(state, restrictingKey, restrictInfo)
		// save restricting account info
		rp.log.Debug("Call SlashingNotify finished", "restrictingInfo", restrictInfo, "account", account, "amount", amount)
	}
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
	restrictInfoByte := state.GetState(vm.RestrictingContractAddr, restrictingKey)
	return restrictingKey, restrictInfoByte
}

func (rp *RestrictingPlugin) mustGetRestrictingInfoByDecode(state xcom.StateDB, account common.Address) ([]byte, restricting.RestrictingInfo, *common.BizError) {
	var restrictInfo restricting.RestrictingInfo
	restrictingKey, restrictInfoByte := rp.getRestrictingInfo(state, account)
	if len(restrictInfoByte) == 0 {
		rp.log.Error("record not found in GetRestrictingInfo", "account", account.String())
		return []byte{}, restrictInfo, restricting.ErrAccountNotFound
	}
	if err := rlp.DecodeBytes(restrictInfoByte, &restrictInfo); err != nil {
		rp.log.Error("Failed to rlp decode restricting account", "error", err.Error())
		return restrictingKey, restrictInfo, common.InternalError.Wrap(err.Error())
	}
	return restrictingKey, restrictInfo, nil
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
		rp.log.Error("Failed to rlp encode restricting info", "error", err, "info", info)
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

		restrictingKey, restrictInfo, err := rp.mustGetRestrictingInfoByDecode(state, account)
		if err != nil {
			if err == restricting.ErrAccountNotFound {
				continue
			}
			return err
		}

		releaseAmountKey, releaseAmount := rp.getReleaseAmount(state, epoch, account)
		rp.log.Debug("Call releaseRestricting: begin to release record", "index", index, "account", account,
			"restrictInfo", restrictInfo, "releaseAmount", releaseAmount)

		if restrictInfo.NeedRelease.Cmp(common.Big0) > 0 {
			restrictInfo.NeedRelease.Add(restrictInfo.NeedRelease, releaseAmount)
		} else {
			canRelease := new(big.Int).Sub(restrictInfo.CachePlanAmount, restrictInfo.AdvanceAmount)
			if canRelease.Cmp(releaseAmount) >= 0 {
				rp.transferAmount(state, vm.RestrictingContractAddr, account, releaseAmount)
				restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, releaseAmount)
			} else {
				needRelease := new(big.Int).Sub(releaseAmount, canRelease)
				rp.transferAmount(state, vm.RestrictingContractAddr, account, canRelease)
				restrictInfo.NeedRelease.Add(restrictInfo.NeedRelease, needRelease)
				restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, canRelease)
			}
		}

		// delete ReleaseAmount
		state.SetState(vm.RestrictingContractAddr, releaseAmountKey, []byte{})
		// delete ReleaseAccount
		state.SetState(vm.RestrictingContractAddr, releaseAccountKey, []byte{})

		// delete epoch in ReleaseList
		// In general, the first epoch is released first.
		// info.ReleaseList = info.ReleaseList[1:]
		restrictInfo.RemoveEpoch(epoch)

		if restrictInfo.CachePlanAmount.Cmp(common.Big0) == 0 {
			if restrictInfo.NeedRelease.Cmp(common.Big0) == 0 || len(restrictInfo.ReleaseList) == 0 {
				//if all is release,remove info
				state.SetState(vm.RestrictingContractAddr, restrictingKey, []byte{})
			} else {
				rp.storeRestrictingInfo(state, restrictingKey, restrictInfo)
			}
		} else {
			rp.storeRestrictingInfo(state, restrictingKey, restrictInfo)
		}
	}

	// delete ReleaseEpoch
	state.SetState(vm.RestrictingContractAddr, releaseEpochKey, []byte{})

	rp.log.Info("Call releaseRestricting finish", "epoch", epoch, "records", numbers)

	return nil
}

func (rp *RestrictingPlugin) getRestrictingInfoToReturn(account common.Address, state xcom.StateDB) (*restricting.Result, *common.BizError) {
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
		plan.Amount = (*hexutil.Big)(bAmount)
		plans = append(plans, plan)
	}

	result.Balance = (*hexutil.Big)(info.CachePlanAmount)
	result.Debt = (*hexutil.Big)(info.NeedRelease)
	result.Entry = plans
	result.Pledge = (*hexutil.Big)(info.AdvanceAmount)
	rp.log.Debug("Call releaseRestricting: query restricting result", "account", account, "result", result)
	return &result, nil
}

func (rp *RestrictingPlugin) GetRestrictingInfo(account common.Address, state xcom.StateDB) (*restricting.Result, *common.BizError) {
	return rp.getRestrictingInfoToReturn(account, state)
}

func GetBlockNumberByEpoch(epoch uint64) uint64 {
	return epoch * xutil.CalcBlocksEachEpoch()
}
