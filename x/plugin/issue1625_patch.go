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

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/params"

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

func (a *FixIssue1625Plugin) fix(blockHash common.Hash, head *types.Header, state xcom.StateDB, chainID *big.Int) error {
	if chainID.Cmp(params.AlayaChainConfig.ChainID) != 0 {
		return nil
	}
	issue1625, err := NewIssue1625Accounts()
	if err != nil {
		return err
	}
	for _, issue1625Account := range issue1625 {
		restrictingKey, restrictInfo, err := rt.mustGetRestrictingInfoByDecode(state, issue1625Account.addr)
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
			if wrongNoUseAmount.Cmp(common.Big0) > 0 {
				restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, wrongNoUseAmount)
				rt.storeRestrictingInfo(state, restrictingKey, restrictInfo)
				log.Debug("fix issue 1625  at no use", "no use", wrongNoUseAmount)
			}

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
		if can.IsNotEmpty() {
			//确保节点是委托的时的那个节点
			if can.StakingBlockNum == stakingBlock {
				delInfo.candidate = can
			}
		} else {
			delInfo.candidate = can
		}
		delInfo.nodeID = nodeID
		delInfo.stakingBlock = stakingBlock
		delInfo.originRestrictingAmount = new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
		delInfo.originFreeAmount = new(big.Int).Add(del.Released, del.ReleasedHes)
		delInfo.canAddr = canAddr
		//如果该委托没有用锁仓，无需回滚
		if delInfo.originRestrictingAmount.Cmp(common.Big0) == 0 {
			continue
		}
		dels = append(dels, delInfo)
	}
	sort.Sort(dels)
	epoch := xutil.CalculateEpoch(blockNumber.Uint64())
	stakingdb := staking.NewStakingDBWithDB(a.sdb)
	for i := 0; i < len(dels); i++ {
		if err := dels[i].handleDelegate(hash, blockNumber, epoch, account, amount, state, stakingdb); err != nil {
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
			//如果该质押没有用锁仓，无需回滚
			if candidate.IsNotEmpty() {
				restrictingAmount := new(big.Int).Add(candidate.RestrictingPlan, candidate.RestrictingPlanHes)
				if restrictingAmount.Cmp(common.Big0) == 0 {
					continue
				}
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
func (a *issue1625AccountStakingInfo) handleExistStaking(hash common.Hash, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
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
		a.candidate.RestrictingPlanHes = new(big.Int).SetInt64(0)
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
	left := new(big.Int)
	if a.originRestrictingAmount.Cmp(rollBackAmount) >= 0 {
		left = new(big.Int).Add(a.originFreeAmount, new(big.Int).Sub(a.originRestrictingAmount, rollBackAmount))
	} else {
		left = new(big.Int).Set(a.originFreeAmount)
	}
	if ok, _ := CheckStakeThreshold(blockNumber.Uint64(), hash, left); !ok {
		return true
	}

	return false
}

func (a *issue1625AccountStakingInfo) calImproperRestrictingAmount(rollBackAmount *big.Int) *big.Int {
	//计算此次需要回退的钱
	improperRestrictingAmount := new(big.Int)
	if rollBackAmount.Cmp(a.originRestrictingAmount) >= 0 {
		improperRestrictingAmount = new(big.Int).Set(a.originRestrictingAmount)
	} else {
		improperRestrictingAmount = new(big.Int).Set(rollBackAmount)
	}
	return improperRestrictingAmount
}

func (a *issue1625AccountStakingInfo) fixCandidateInfo(improperRestrictingAmount *big.Int) {
	//修正质押信息
	if a.candidate.RestrictingPlanHes.Cmp(improperRestrictingAmount) >= 0 {
		a.candidate.RestrictingPlanHes.Sub(a.candidate.RestrictingPlanHes, improperRestrictingAmount)
	} else {
		hes := new(big.Int).Set(a.candidate.RestrictingPlanHes)
		a.candidate.RestrictingPlanHes = new(big.Int)
		a.candidate.RestrictingPlan = new(big.Int).Sub(a.candidate.RestrictingPlan, new(big.Int).Sub(improperRestrictingAmount, hes))
	}
	a.candidate.SubShares(improperRestrictingAmount)
}

//减持质押
func (a *issue1625AccountStakingInfo) decreaseStaking(hash common.Hash, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB) error {
	if err := stk.db.DelCanPowerStore(hash, a.candidate); nil != err {
		return err
	}

	//计算此次需要回退的钱
	improperRestrictingAmount := a.calImproperRestrictingAmount(rollBackAmount)

	//回退因为漏洞产生的金额
	if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, improperRestrictingAmount, state); nil != err {
		return err
	}

	//修正质押信息
	a.fixCandidateInfo(improperRestrictingAmount)

	a.candidate.StakingEpoch = uint32(epoch)

	if err := stk.db.SetCanPowerStore(hash, a.canAddr, a.candidate); nil != err {
		return err
	}
	if err := stk.db.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
		return err
	}

	rollBackAmount.Sub(rollBackAmount, improperRestrictingAmount)
	return nil
}

//撤销质押
func (a *issue1625AccountStakingInfo) withdrewStaking(hash common.Hash, epoch uint64, blockNumber *big.Int, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	if err := stdb.DelCanPowerStore(hash, a.candidate); nil != err {
		return err
	}

	//计算此次需要回退的钱
	improperRestrictingAmount := a.calImproperRestrictingAmount(rollBackAmount)

	//回退因为漏洞产生的金额
	if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, improperRestrictingAmount, state); nil != err {
		return err
	}

	//修正质押信息
	a.fixCandidateInfo(improperRestrictingAmount)

	//开始解质押
	//回退犹豫期的自由金
	if a.candidate.ReleasedHes.Cmp(common.Big0) > 0 {
		rt.transferAmount(state, vm.StakingContractAddr, a.candidate.StakingAddress, a.candidate.ReleasedHes)
		a.candidate.ReleasedHes = new(big.Int).SetInt64(0)
	}

	//回退犹豫期的锁仓
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

	rollBackAmount.Sub(rollBackAmount, improperRestrictingAmount)
	return nil
}

func (a *issue1625AccountStakingInfo) handleStaking(hash common.Hash, blockNumber *big.Int, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	lazyCalcStakeAmount(epoch, a.candidate.CandidateMutable)
	log.Debug("fix issue 1625 for staking begin", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "return", rollBackAmount, "restrictingPlan",
		a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes, "released", a.candidate.Released, "releasedHes", a.candidate.ReleasedHes, "share", a.candidate.Shares)
	if a.candidate.Status.IsWithdrew() {
		//已经解质押,节点处于退出锁定期
		if err := a.handleExistStaking(hash, epoch, rollBackAmount, state, stdb); err != nil {
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
	nodeID       discover.NodeID

	originRestrictingAmount, originFreeAmount *big.Int
}

func (a *issue1625AccountDelInfo) shouldWithdrewDel(hash common.Hash, blockNumber *big.Int, rollBackAmount *big.Int) bool {
	leftTotalDelgateAmount := new(big.Int)
	if rollBackAmount.Cmp(a.originRestrictingAmount) >= 0 {
		leftTotalDelgateAmount.Set(a.originFreeAmount)
	} else {
		leftTotalDelgateAmount.Add(a.originFreeAmount, new(big.Int).Sub(a.originRestrictingAmount, rollBackAmount))
	}
	if ok, _ := CheckOperatingThreshold(blockNumber.Uint64(), hash, leftTotalDelgateAmount); ok {
		return false
	}
	return true
}

func (a *issue1625AccountDelInfo) handleDelegate(hash common.Hash, blockNumber *big.Int, epoch uint64, delAddr common.Address, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	improperRestrictingAmount := new(big.Int)
	if rollBackAmount.Cmp(a.originRestrictingAmount) >= 0 {
		improperRestrictingAmount = new(big.Int).Set(a.originRestrictingAmount)
	} else {
		improperRestrictingAmount = new(big.Int).Set(rollBackAmount)
	}
	log.Debug("fix issue 1625 for delegate begin", "account", delAddr, "candidate", a.nodeID.String(), "currentReturn", improperRestrictingAmount, "leftReturn", rollBackAmount, "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
		"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "CumulativeIncome", a.del.CumulativeIncome)
	if a.candidate.IsNotEmpty() {
		log.Debug("fix issue 1625 for delegate ,can begin info", "account", delAddr, "candidate", a.nodeID.String(), "share", a.candidate.Shares, "candidate.del", a.candidate.DelegateTotal, "candidate.delhes", a.candidate.DelegateTotalHes, "canValid", a.candidate.IsValid())
	}
	//先计算委托收益
	delegateRewardPerList, err := RewardMgrInstance().GetDelegateRewardPerList(hash, a.nodeID, a.stakingBlock, uint64(a.del.DelegateEpoch), xutil.CalculateEpoch(blockNumber.Uint64())-1)
	if snapshotdb.NonDbNotFoundErr(err) {
		return err
	}

	rewardsReceive := calcDelegateIncome(epoch, a.del, delegateRewardPerList)
	if err := UpdateDelegateRewardPer(hash, a.nodeID, a.stakingBlock, rewardsReceive, rm.db); err != nil {
		return err
	}
	if a.candidate.IsNotEmpty() {
		lazyCalcNodeTotalDelegateAmount(epoch, a.candidate.CandidateMutable)
	}

	a.del.DelegateEpoch = uint32(epoch)

	withdrewDel := a.shouldWithdrewDel(hash, blockNumber, rollBackAmount)
	if withdrewDel {
		//回滚错误金额
		if err := a.fixImproperRestrictingAmountByDel(delAddr, improperRestrictingAmount, state); err != nil {
			return err
		}

		//开始撤销委托
		//更新candidate中由于撤销委托导致的记录变动
		if a.candidate.IsNotEmpty() {
			hes := new(big.Int).Add(a.del.ReleasedHes, a.del.RestrictingPlanHes)
			lock := new(big.Int).Add(a.del.Released, a.del.RestrictingPlan)
			a.candidate.DelegateTotalHes.Sub(a.candidate.DelegateTotalHes, hes)
			a.candidate.DelegateTotal.Sub(a.candidate.DelegateTotal, lock)
			if a.candidate.Shares.Cmp(new(big.Int).Add(hes, lock)) >= 0 {
				a.candidate.Shares.Sub(a.candidate.Shares, new(big.Int).Add(hes, lock))
			}
		}

		//回退锁仓
		if err := rt.ReturnLockFunds(delAddr, new(big.Int).Add(a.del.RestrictingPlan, a.del.RestrictingPlanHes), state); err != nil {
			return err
		}

		//回退自由
		rt.transferAmount(state, vm.StakingContractAddr, delAddr, new(big.Int).Add(a.del.ReleasedHes, a.del.Released))

		//领取收益
		if err := rm.ReturnDelegateReward(delAddr, a.del.CumulativeIncome, state); err != nil {
			return common.InternalError
		}

		//删除委托
		if err := stdb.DelDelegateStore(hash, delAddr, a.candidate.NodeId, a.stakingBlock); nil != err {
			return err
		}

		log.Debug("fix issue 1625 for delegate,withdrew del", "account", delAddr, "candidate", a.nodeID.String(), "income", a.del.CumulativeIncome)
	} else {
		//不需要解除委托
		if err := a.fixImproperRestrictingAmountByDel(delAddr, improperRestrictingAmount, state); err != nil {
			return err
		}
		if err := stdb.SetDelegateStore(hash, delAddr, a.candidate.NodeId, a.stakingBlock, a.del); nil != err {
			return err
		}
		log.Debug("fix issue 1625 for delegate,decrease del", "account", delAddr, "candidate", a.nodeID.String(), "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
			"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "income", a.del.CumulativeIncome)
	}

	if a.candidate.IsNotEmpty() {
		if a.candidate.IsValid() {
			if err := stdb.DelCanPowerStore(hash, a.candidate); nil != err {
				return err
			}
			if err := stdb.SetCanPowerStore(hash, a.canAddr, a.candidate); nil != err {
				return err
			}
		}
		if err := stdb.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
			return err
		}
	}

	rollBackAmount.Sub(rollBackAmount, improperRestrictingAmount)

	if !a.candidate.IsEmpty() {
		log.Debug("fix issue 1625 for delegate,can last info", "account", delAddr, "candidate", a.nodeID.String(), "share", a.candidate.Shares, "candidate.del", a.candidate.DelegateTotal, "candidate.delhes", a.candidate.DelegateTotalHes)
	}
	log.Debug("fix issue 1625 for delegate end", "account", delAddr, "candidate", a.nodeID.String(), "leftReturn", rollBackAmount, "withdrewDel", withdrewDel)
	return nil
}

//修正委托以及验证人的锁仓信息
func (a *issue1625AccountDelInfo) fixImproperRestrictingAmountByDel(delAddr common.Address, improperRestrictingAmount *big.Int, state xcom.StateDB) error {
	if err := rt.ReturnWrongLockFunds(delAddr, improperRestrictingAmount, state); nil != err {
		return err
	}
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
	if a.candidate.IsNotEmpty() {
		if a.candidate.Shares.Cmp(improperRestrictingAmount) >= 0 {
			a.candidate.SubShares(improperRestrictingAmount)
		}
	}
	return nil
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
