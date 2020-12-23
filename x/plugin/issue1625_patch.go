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
	"bytes"
	"fmt"
	"math/big"
	"sort"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

func NewFixIssue1625Plugin(sdb snapshotdb.DB) *FixIssue1625Plugin {
	fix := new(FixIssue1625Plugin)
	fix.sdb = sdb
	return fix
}

type FixIssue1625Plugin struct {
	sdb snapshotdb.DB
}

type issue1625Accounts struct {
	addr   common.Address
	amount *big.Int
}

func (a *FixIssue1625Plugin) fix(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	for _, issue1625Account := range []issue1625Accounts{} {
		restrictingKey, restrictInfo, err := rt.getRestrictingInfoByDecode(state, issue1625Account.addr)
		if err != nil {
			return err
		}
		log.Debug("fix issue 1625 begin", "account", issue1625Account.addr, "fix amount", issue1625Account.amount, "info", restrictInfo)

		actualRestrictingAmount := new(big.Int).Sub(restrictInfo.CachePlanAmount, issue1625Account.amount)
		if actualRestrictingAmount.Cmp(common.Big0) < 0 {
			log.Error("seems not good here", "info", restrictInfo, "amount", issue1625Account.amount, "account", issue1625Account.addr)
			return fmt.Errorf("the account restrictInfo seems not right")
		}

		wrongStakingAmount := new(big.Int).Sub(restrictInfo.StakingAmount, actualRestrictingAmount)
		if wrongStakingAmount.Cmp(common.Big0) > 0 {
			//If the user uses the wrong amount,Roll back the unused part first
			//优先回滚没有使用的那部分锁仓余额
			wrongNoUseAmount := new(big.Int).Sub(issue1625Account.amount, wrongStakingAmount)
			restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, wrongNoUseAmount)
			rt.storeRestrictingInfo(state, restrictingKey, restrictInfo)
			log.Debug("fix issue 1625  at no use", "no use", wrongNoUseAmount)
			//roll back del,回滚委托
			if err := a.rollBackDel(blockHash, head.Number, issue1625Account.addr, wrongStakingAmount, state); err != nil {
				return err
			}
			//roll back staking,回滚质押
			if wrongStakingAmount.Cmp(common.Big0) > 0 {
				if err := a.rollBackStaking(blockHash, head.Number, issue1625Account.addr, wrongStakingAmount, state); err != nil {
					return err
				}
			}
		} else {
			//当用户没有使用因为漏洞产生的钱，直接减去漏洞的钱就是正确的余额
			restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, issue1625Account.amount)
			if restrictInfo.StakingAmount.Cmp(common.Big0) == 0 &&
				len(restrictInfo.ReleaseList) == 0 && restrictInfo.CachePlanAmount.Cmp(common.Big0) == 0 {
				state.SetState(vm.RestrictingContractAddr, restrictingKey, []byte{})
				log.Debug("fix issue 1625 finished,set info empty", "account", issue1625Account.addr, "fix amount", issue1625Account.amount)
			} else {
				rt.storeRestrictingInfo(state, restrictingKey, restrictInfo)
				log.Debug("fix issue 1625 finished", "account", issue1625Account.addr, "info", restrictInfo, "fix amount", issue1625Account.amount)
			}
		}
	}
	return nil
}

func (a *FixIssue1625Plugin) rollBackDel(hash common.Hash, blockNumber *big.Int, account common.Address, amount *big.Int, state xcom.StateDB) error {

	delAddrByte := account.Bytes()

	markPre := len(staking.DelegateKeyPrefix)
	markDelAddr := markPre + len(delAddrByte)

	key := make([]byte, markDelAddr)
	copy(key[:markPre], staking.DelegateKeyPrefix)
	copy(key[markPre:markDelAddr], delAddrByte)

	iter := a.sdb.Ranking(hash, key, 0)
	if err := iter.Error(); nil != err {
		return err
	}
	defer iter.Release()

	var dels issue1625AccountDelInfos

	for iter.Valid(); iter.Next(); {
		var del staking.Delegation
		if err := rlp.DecodeBytes(iter.Value(), &del); nil != err {
			return err
		}
		_, nodeID, stakingBlock := staking.DecodeDelegateKey(iter.Key())
		canAddr, err := xutil.NodeId2Addr(nodeID)
		if nil != err {
			return err
		}
		can, err := stk.db.GetCandidateStore(hash, canAddr)
		if snapshotdb.NonDbNotFoundErr(err) {
			return err
		}

		delInfo := new(issue1625AccountDelInfo)
		delInfo.del = &del
		delInfo.candidate = can
		delInfo.stakingBlock = stakingBlock
		delInfo.originRestrictingAmount = new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
		delInfo.originFreeAmount = new(big.Int).Add(del.Released, del.ReleasedHes)
		delInfo.canAddr = canAddr
		dels = append(dels, delInfo)
	}
	sort.Sort(dels)
	epoch := xutil.CalculateEpoch(blockNumber.Uint64())
	stakingdb := staking.NewStakingDBWithDB(a.sdb)
	for i := 0; i < len(dels); i++ {
		if _, err := dels[i].handleDelegate(hash, blockNumber, epoch, account, amount, state, stakingdb); err != nil {
			return err
		}
		if amount.Cmp(common.Big0) <= 0 {
			break
		}
	}
	return nil
}

func (a *FixIssue1625Plugin) rollBackStaking(hash common.Hash, blockNumber *big.Int, account common.Address, amount *big.Int, state xcom.StateDB) error {

	iter := a.sdb.Ranking(hash, staking.CanBaseKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		return err
	}
	defer iter.Release()

	stakingdb := staking.NewStakingDBWithDB(a.sdb)

	var stakings issue1625AccountStakingInfos
	for iter.Valid(); iter.Next(); {
		var canbase staking.CandidateBase
		if err := rlp.DecodeBytes(iter.Value(), &canbase); nil != err {
			return err
		}

		if canbase.StakingAddress == account {
			canAddr, err := xutil.NodeId2Addr(canbase.NodeId)
			if nil != err {
				return err
			}
			canmu, err := stakingdb.GetCanMutableStore(hash, canAddr)
			if nil != err {
				return err
			}
			candidate := staking.Candidate{
				&canbase, canmu,
			}
			stakings = append(stakings, newIssue1625AccountStakingInfo(&candidate, canAddr))
		}
	}
	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	sort.Sort(stakings)
	for i := 0; i < len(stakings); i++ {
		if err := stakings[i].handleStaking(hash, blockNumber, epoch, amount, state, stakingdb); err != nil {
			return err
		}
		if amount.Cmp(common.Big0) <= 0 {
			break
		}
	}
	return nil
}

func newIssue1625AccountStakingInfo(candidate *staking.Candidate, canAddr common.NodeAddress) *issue1625AccountStakingInfo {
	w := &issue1625AccountStakingInfo{
		candidate: candidate,
		canAddr:   canAddr,
	}
	w.originRestrictingAmount = new(big.Int).Add(w.candidate.RestrictingPlan, w.candidate.RestrictingPlanHes)
	w.originFreeAmount = new(big.Int).Add(w.candidate.Released, w.candidate.ReleasedHes)
	return w
}

type issue1625AccountStakingInfo struct {
	candidate                                 *staking.Candidate
	canAddr                                   common.NodeAddress
	originRestrictingAmount, originFreeAmount *big.Int
}

//回退处于退出期的质押信息
func (a *issue1625AccountStakingInfo) handelExistStaking(hash common.Hash, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	if a.originRestrictingAmount.Cmp(rollBackAmount) >= 0 {
		if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, rollBackAmount, state); nil != err {
			return err
		}
		a.candidate.RestrictingPlan = new(big.Int).Sub(a.candidate.RestrictingPlan, rollBackAmount)
		rollBackAmount.SetInt64(0)
	} else {
		if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, a.originRestrictingAmount, state); nil != err {
			return err
		}
		a.candidate.RestrictingPlan = new(big.Int).SetInt64(0)
		rollBackAmount.Sub(rollBackAmount, a.originRestrictingAmount)
	}

	a.candidate.StakingEpoch = uint32(epoch)

	if err := stdb.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
		return err
	}

	return nil
}

//检查是否达到质押门槛
func (a *issue1625AccountStakingInfo) shouldWithdrewStaking(hash common.Hash, blockNumber *big.Int, rollBackAmount *big.Int) bool {
	if a.originRestrictingAmount.Cmp(rollBackAmount) >= 0 {
		left := new(big.Int).Add(a.originFreeAmount, new(big.Int).Sub(a.originRestrictingAmount, rollBackAmount))
		if ok, _ := CheckStakeThreshold(blockNumber.Uint64(), hash, left); !ok {
			return true
		}
	} else {
		left := new(big.Int).Sub(a.originFreeAmount, a.originRestrictingAmount)
		if ok, _ := CheckStakeThreshold(blockNumber.Uint64(), hash, left); !ok {
			return true
		}
	}
	return false
}

//减持质押
func (a *issue1625AccountStakingInfo) decreaseStaking(hash common.Hash, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB) error {
	realSub, err := a.refundWrongLockFunds(rollBackAmount, state)
	if err != nil {
		return err
	}
	if err := stk.db.DelCanPowerStore(hash, a.candidate); nil != err {
		return err
	}

	a.candidate.StakingEpoch = uint32(epoch)
	a.candidate.SubShares(realSub)

	if err := stk.db.SetCanPowerStore(hash, a.canAddr, a.candidate); nil != err {
		return err
	}

	if err := stk.db.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
		return err
	}
	return nil
}

//撤销质押
func (a *issue1625AccountStakingInfo) withdrewStaking(hash common.Hash, epoch uint64, blockNumber *big.Int, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	if err := stdb.DelCanPowerStore(hash, a.candidate); nil != err {
		return err
	}

	//回退自由金
	if a.candidate.ReleasedHes.Cmp(common.Big0) > 0 {
		state.AddBalance(a.candidate.StakingAddress, a.candidate.ReleasedHes)
		state.SubBalance(vm.StakingContractAddr, a.candidate.ReleasedHes)
		a.candidate.ReleasedHes = new(big.Int).SetInt64(0)
	}

	//回退因为漏洞产生的金额
	if _, err := a.refundWrongLockFunds(rollBackAmount, state); err != nil {
		return err
	}

	//回退锁仓
	if a.candidate.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		err := rt.ReturnLockFunds(a.candidate.StakingAddress, a.candidate.RestrictingPlanHes, state)
		if nil != err {
			return err
		}
		a.candidate.RestrictingPlanHes = new(big.Int).SetInt64(0)
	}

	a.candidate.CleanShares()
	a.candidate.Status |= staking.Invalided | staking.Withdrew

	a.candidate.StakingEpoch = uint32(epoch)

	if a.candidate.Released.Cmp(common.Big0) > 0 || a.candidate.RestrictingPlan.Cmp(common.Big0) > 0 {
		//如果质押处于生效期，需要锁定
		if err := stk.addErrorAccountUnStakeItem(blockNumber.Uint64(), hash, a.candidate.NodeId, a.canAddr, a.candidate.StakingBlockNum); nil != err {
			return err
		}
		// sub the account staking Reference Count
		if err := stdb.SubAccountStakeRc(hash, a.candidate.StakingAddress); nil != err {
			return err
		}
		if err := stdb.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
			return err
		}
	} else {
		//如果质押还处于犹豫期，不用锁定
		if err := stdb.DelCandidateStore(hash, a.canAddr); nil != err {
			return err
		}
	}
	return nil
}

func (a *issue1625AccountStakingInfo) handleStaking(hash common.Hash, blockNumber *big.Int, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	lazyCalcStakeAmount(epoch, a.candidate.CandidateMutable)
	log.Debug("fix issue 1625 for staking begin", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "return", rollBackAmount, "restrictingPlan",
		a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes, "released", a.candidate.Released, "releasedHes", a.candidate.ReleasedHes, "share", a.candidate.Shares)
	if a.candidate.Status.IsWithdrew() {
		//已经解质押,节点处于退出锁定期
		if err := a.handelExistStaking(hash, epoch, rollBackAmount, state, stdb); err != nil {
			return err
		}
		log.Debug("fix issue 1625 for staking end", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "status", a.candidate.Status, "return",
			rollBackAmount, "restrictingPlan", a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes, "released", a.candidate.Released, "releasedHes", a.candidate.ReleasedHes, "share", a.candidate.Shares)
	} else {
		//如果质押中,根据回退后的剩余质押金额是否达到质押门槛来判断时候需要撤销或者减持质押
		shouldWithdrewStaking := a.shouldWithdrewStaking(hash, blockNumber, rollBackAmount)
		if shouldWithdrewStaking {
			//撤销质押
			if err := a.withdrewStaking(hash, epoch, blockNumber, rollBackAmount, state, stdb); err != nil {
				return err
			}
		} else {
			//减持质押
			if err := a.decreaseStaking(hash, epoch, rollBackAmount, state); err != nil {
				return err
			}
		}
		log.Debug("fix issue 1625 for staking end", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "status", a.candidate.Status, "withdrewStaking",
			shouldWithdrewStaking, "return", rollBackAmount, "restrictingPlan", a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes, "released", a.candidate.Released,
			"releasedHes", a.candidate.ReleasedHes, "share", a.candidate.Shares)
	}
	return nil
}

//回退因漏洞产生的锁仓金额
func (a *issue1625AccountStakingInfo) refundWrongLockFunds(rollBackAmount *big.Int, state xcom.StateDB) (*big.Int, error) {
	realSub := new(big.Int)
	if a.originRestrictingAmount.Cmp(rollBackAmount) >= 0 {
		if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, rollBackAmount, state); nil != err {
			return nil, err
		}
		realSub.Set(rollBackAmount)
		if a.candidate.RestrictingPlanHes.Cmp(rollBackAmount) >= 0 {
			a.candidate.RestrictingPlanHes = new(big.Int).Sub(a.candidate.RestrictingPlanHes, rollBackAmount)
		} else {
			a.candidate.RestrictingPlan = new(big.Int).Sub(a.candidate.RestrictingPlan, new(big.Int).Sub(rollBackAmount, a.candidate.RestrictingPlanHes))
			a.candidate.RestrictingPlanHes = new(big.Int).SetInt64(0)
		}
		rollBackAmount.SetInt64(0)
	} else {
		if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, a.originRestrictingAmount, state); nil != err {
			return nil, err
		}
		realSub.Set(a.originRestrictingAmount)
		a.candidate.RestrictingPlanHes = new(big.Int)
		a.candidate.RestrictingPlan = new(big.Int)
		rollBackAmount.Sub(rollBackAmount, a.originRestrictingAmount)
	}
	return realSub, nil
}

type issue1625AccountStakingInfos []*issue1625AccountStakingInfo

func (d issue1625AccountStakingInfos) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d issue1625AccountStakingInfos) Len() int      { return len(d) }
func (d issue1625AccountStakingInfos) Less(i, j int) bool {
	//1.节点处于解质押状态,并且质押的时间靠后的排在前面
	//2.按照节点质押的时间远近进行排序，并且质押的时间靠后的排在前面
	if d[i].candidate.IsWithdrew() {
		if d[j].candidate.IsWithdrew() {
			return d.LessByStakingBlockNum(i, j)
		} else {
			return false
		}
	} else {
		if d[j].candidate.IsWithdrew() {
			return false
		} else {
			return d.LessByStakingBlockNum(i, j)
		}
	}
}

func (d issue1625AccountStakingInfos) LessByStakingBlockNum(i, j int) bool {
	if d[i].candidate.StakingBlockNum > d[j].candidate.StakingBlockNum {
		return true
	} else {
		return false
	}
}

type issue1625AccountDelInfo struct {
	del       *staking.Delegation
	candidate *staking.Candidate
	//use for get staking
	stakingBlock uint64
	canAddr      common.NodeAddress

	originRestrictingAmount, originFreeAmount *big.Int
}

func (a *issue1625AccountDelInfo) handleDelegate(hash common.Hash, blockNumber *big.Int, epoch uint64, delAddr common.Address, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) (*big.Int, error) {

	refundAmount := new(big.Int)
	improperRestrictingAmount := new(big.Int)

	//优先回滚此次委托中锁仓的部分，如果剩余的金额没有达到委托门槛，就撤销委托
	if rollBackAmount.Cmp(a.originRestrictingAmount) >= 0 {
		if ok, _ := CheckOperatingThreshold(blockNumber.Uint64(), hash, a.originFreeAmount); !ok {
			refundAmount = new(big.Int).Add(a.originFreeAmount, a.originRestrictingAmount)
		} else {
			refundAmount = new(big.Int).Set(a.originRestrictingAmount)
		}
		improperRestrictingAmount = new(big.Int).Set(a.originRestrictingAmount)
	} else {
		resrtict := new(big.Int).Sub(a.originRestrictingAmount, rollBackAmount)
		if ok, _ := CheckOperatingThreshold(blockNumber.Uint64(), hash, new(big.Int).Add(a.originFreeAmount, resrtict)); !ok {
			refundAmount = new(big.Int).Add(a.originFreeAmount, a.originRestrictingAmount)
		} else {
			refundAmount = new(big.Int).Set(rollBackAmount)
		}
		improperRestrictingAmount = new(big.Int).Set(rollBackAmount)
	}

	log.Debug("fix issue 1625 for delegate begin", "account", delAddr, "currentReturn", improperRestrictingAmount, "leftReturn", rollBackAmount, "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
		"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "share", a.candidate.Shares)

	rollBackAmount.Sub(rollBackAmount, improperRestrictingAmount)

	//节点总共撤回了委托的钱
	realSub := new(big.Int).Set(refundAmount)

	delegateRewardPerList, err := RewardMgrInstance().GetDelegateRewardPerList(hash, a.candidate.NodeId, a.stakingBlock, uint64(a.del.DelegateEpoch), xutil.CalculateEpoch(blockNumber.Uint64())-1)
	if snapshotdb.NonDbNotFoundErr(err) {
		return nil, err
	}

	rewardsReceive := calcDelegateIncome(epoch, a.del, delegateRewardPerList)
	if err := UpdateDelegateRewardPer(hash, a.candidate.NodeId, a.stakingBlock, rewardsReceive, rm.db); err != nil {
		return nil, err
	}
	if a.candidate.IsNotEmpty() {
		lazyCalcNodeTotalDelegateAmount(epoch, a.candidate.CandidateMutable)
	}

	a.del.DelegateEpoch = uint32(epoch)
	switch {
	// Illegal parameter
	case a.candidate.IsNotEmpty() && a.stakingBlock > a.candidate.StakingBlockNum:
		return nil, staking.ErrBlockNumberDisordered
	default:
		// handle delegate on Hesitate period
		if refundAmount.Cmp(common.Big0) > 0 {
			// When remain is greater than or equal to del.ReleasedHes/del.Releaseds
			//回退因为漏洞产生的金额
			if err := rt.ReturnWrongLockFunds(delAddr, improperRestrictingAmount, state); nil != err {
				return nil, err
			}

			//更新锁仓信息
			refundAmount = new(big.Int).Sub(refundAmount, improperRestrictingAmount)
			if a.del.RestrictingPlanHes.Cmp(improperRestrictingAmount) >= 0 {
				a.del.RestrictingPlanHes.Sub(a.del.RestrictingPlanHes, improperRestrictingAmount)
				if a.candidate.IsNotEmpty() {
					a.candidate.DelegateTotalHes.Sub(a.candidate.DelegateTotalHes, improperRestrictingAmount)
				}
			} else {
				hes := new(big.Int).Set(a.del.RestrictingPlanHes)
				a.del.RestrictingPlanHes = new(big.Int)
				a.del.RestrictingPlan = new(big.Int).Sub(a.del.RestrictingPlan, new(big.Int).Sub(improperRestrictingAmount, hes))
				if a.candidate.IsNotEmpty() {
					a.candidate.DelegateTotalHes.Sub(a.candidate.DelegateTotalHes, hes)
					a.candidate.DelegateTotal.Sub(a.candidate.DelegateTotal, new(big.Int).Sub(improperRestrictingAmount, hes))
				}
			}
		}

		//如果还有回退金额，处理犹豫期的钱
		if refundAmount.Cmp(common.Big0) > 0 {
			rm, rbalance, lbalance, err := rufundDelegateFn(refundAmount, a.del.ReleasedHes, a.del.RestrictingPlanHes, delAddr, state)
			if nil != err {
				return nil, err
			}
			if a.candidate.IsNotEmpty() {
				a.candidate.DelegateTotalHes = new(big.Int).Sub(a.candidate.DelegateTotalHes, new(big.Int).Sub(refundAmount, rm))
			}
			refundAmount, a.del.ReleasedHes, a.del.RestrictingPlanHes = rm, rbalance, lbalance
		}

		//如果还有回退金额，处理锁定期的钱
		// handle delegate on Effective period
		if refundAmount.Cmp(common.Big0) > 0 {
			rm, rbalance, lbalance, err := rufundDelegateFn(refundAmount, a.del.Released, a.del.RestrictingPlan, delAddr, state)
			if nil != err {
				return nil, err
			}
			if a.candidate.IsNotEmpty() {
				a.candidate.DelegateTotal = new(big.Int).Sub(a.candidate.DelegateTotal, new(big.Int).Sub(refundAmount, rm))
			}
			refundAmount, a.del.Released, a.del.RestrictingPlan = rm, rbalance, lbalance
		}

		if refundAmount.Cmp(common.Big0) != 0 {
			return nil, staking.ErrWrongWithdrewDelVonCalc
		}
		total := new(big.Int).Add(a.originFreeAmount, a.originFreeAmount)
		// If total had full sub,
		// then clean the delegate info
		issueIncome := new(big.Int)
		if total.Cmp(realSub) == 0 {
			// When the entrusted information is deleted, the entrusted proceeds need to be issued automatically
			issueIncome = issueIncome.Add(issueIncome, a.del.CumulativeIncome)
			if err := rm.ReturnDelegateReward(delAddr, a.del.CumulativeIncome, state); err != nil {
				return nil, common.InternalError
			}
			if err := stdb.DelDelegateStore(hash, delAddr, a.candidate.NodeId, a.stakingBlock); nil != err {
				return nil, err
			}
		} else {
			if err := stdb.SetDelegateStore(hash, delAddr, a.candidate.NodeId, a.stakingBlock, a.del); nil != err {
				return nil, err
			}
		}

	}

	if a.candidate.IsNotEmpty() && a.stakingBlock == a.candidate.StakingBlockNum {
		if a.candidate.IsValid() {
			if err := stdb.DelCanPowerStore(hash, a.candidate); nil != err {
				return nil, err
			}

			// change candidate shares
			if a.candidate.Shares.Cmp(realSub) > 0 {
				a.candidate.SubShares(realSub)
			} else {
				panic("the candidate shares is no enough")
			}

			if err := stdb.SetCanPowerStore(hash, a.canAddr, a.candidate); nil != err {
				return nil, err
			}
		}

		if err := stdb.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
			return nil, err
		}
	}
	log.Debug("fix issue 1625 for delegate end", "account", delAddr, "currentReturn", improperRestrictingAmount, "leftReturn", rollBackAmount, "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
		"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "share", a.candidate.Shares)
	return nil, nil
}

type issue1625AccountDelInfos []*issue1625AccountDelInfo

func (d issue1625AccountDelInfos) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d issue1625AccountDelInfos) Len() int      { return len(d) }
func (d issue1625AccountDelInfos) Less(i, j int) bool {
	//排序顺序
	//1.委托的节点已经完全退出,并且委托的时间靠后
	//2.委托的节点处于解质押状态,并且委托的时间靠后
	//3.根据委托节点的分红比例从小到大排序，如果委托比例相同，根据节点id从小到大排序
	if d[i].candidate.IsEmpty() {
		if d[j].candidate.IsEmpty() {
			return d.LessDelByEpoch(i, j)
		} else {
			return true
		}
	} else {
		if d[j].candidate.IsEmpty() {
			return false
		} else {
			if d[i].candidate.IsWithdrew() {
				if d[j].candidate.IsWithdrew() {
					return d.LessDelByEpoch(i, j)
				} else {
					return true
				}
			} else {
				if d[j].candidate.IsWithdrew() {
					return false
				} else {
					return d.LessDelByRewardPer(i, j)
				}
			}
		}
	}

}

func (d issue1625AccountDelInfos) LessDelByEpoch(i, j int) bool {
	if d[i].del.DelegateEpoch > d[j].del.DelegateEpoch {
		return true
	} else if d[i].del.DelegateEpoch == d[j].del.DelegateEpoch {
		if bytes.Compare(d[i].candidate.NodeId.Bytes(), d[j].candidate.NodeId.Bytes()) < 0 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}

func (d issue1625AccountDelInfos) LessDelByRewardPer(i, j int) bool {
	if d[i].candidate.RewardPer < d[j].candidate.RewardPer {
		return true
	} else if d[i].candidate.RewardPer == d[j].candidate.RewardPer {
		if bytes.Compare(d[i].candidate.NodeId.Bytes(), d[j].candidate.NodeId.Bytes()) < 0 {
			return true
		} else {
			return false
		}
	} else {
		return false
	}
}
