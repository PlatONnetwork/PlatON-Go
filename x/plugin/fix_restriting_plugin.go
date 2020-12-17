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

func NewFixRestrictingPlugin(sdb snapshotdb.DB) *FixRestrictingPlugin {
	fix := new(FixRestrictingPlugin)
	fix.sdb = sdb
	return fix
}

type FixRestrictingPlugin struct {
	sdb snapshotdb.DB
}

type rollBackInfo struct {
	addr   common.Address
	amount *big.Int
}

func (a *FixRestrictingPlugin) fix(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	for _, wrong := range []rollBackInfo{} {
		restrictingKey, restrictInfo, err := rt.getRestrictingInfoByDecode(state, wrong.addr)
		if err != nil {
			return err
		}
		log.Debug("fix restricting begin", "account", wrong.addr, "fix amount", wrong.amount, "info", restrictInfo)

		realAmount := new(big.Int).Sub(restrictInfo.CachePlanAmount, wrong.amount)
		if realAmount.Cmp(common.Big0) < 0 {
			log.Error("seems not good here", "info", restrictInfo, "amount", wrong.amount, "account", wrong.addr)
			return fmt.Errorf("the account restrictInfo seems not right")
		}

		wrongStakingAmount := new(big.Int).Sub(restrictInfo.StakingAmount, realAmount)
		if wrongStakingAmount.Cmp(common.Big0) > 0 {
			//If the user uses the wrong amount,Roll back the unused part first
			wrongNoUseAmount := new(big.Int).Sub(wrong.amount, wrongStakingAmount)
			restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, wrongNoUseAmount)
			rt.storeRestrictingInfo(state, restrictingKey, restrictInfo)
			log.Debug("fix restricting  in wrongStakingAmount", "no use", wrongNoUseAmount)
			//roll back del
			if err := a.rollBackDel(blockHash, head.Number, wrong.addr, wrongStakingAmount, state); err != nil {
				return err
			}
			//roll back staking
			if wrongStakingAmount.Cmp(common.Big0) > 0 {
				if err := a.rollBackStaking(blockHash, head.Number, wrong.addr, wrongStakingAmount, state); err != nil {
					return err
				}
			}
		} else {
			//犹豫期解质押，此时查不到质押信息
			restrictInfo.CachePlanAmount.Sub(restrictInfo.CachePlanAmount, wrong.amount)
			if restrictInfo.StakingAmount.Cmp(common.Big0) == 0 &&
				len(restrictInfo.ReleaseList) == 0 && restrictInfo.CachePlanAmount.Cmp(common.Big0) == 0 {
				state.SetState(vm.RestrictingContractAddr, restrictingKey, []byte{})
				log.Debug("fix restricting finished,set info empty", "account", wrong.addr, "fix amount", wrong.amount)
			} else {
				rt.storeRestrictingInfo(state, restrictingKey, restrictInfo)
				log.Debug("fix restricting finished", "account", wrong.addr, "info", restrictInfo, "fix amount", wrong.amount)
			}
		}
	}
	return nil
}

func (a *FixRestrictingPlugin) rollBackDel(hash common.Hash, blockNumber *big.Int, account common.Address, amount *big.Int, state xcom.StateDB) error {

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

	var wrongDels rollBackDelInfos

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

		wrongDel := new(rollBackDelInfo)
		wrongDel.del = &del
		wrongDel.candidate = can
		wrongDel.stakingBlock = stakingBlock
		wrongDel.originRestrictingAmount = new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
		wrongDel.originFreeAmount = new(big.Int).Add(del.Released, del.ReleasedHes)
		wrongDels = append(wrongDels, wrongDel)
	}
	sort.Sort(wrongDels)
	epoch := xutil.CalculateEpoch(blockNumber.Uint64())
	stakingdb := staking.NewStakingDBWithDB(a.sdb)
	for i := 0; i < len(wrongDels); i++ {
		if _, err := wrongDels[i].HandleDelegate(hash, blockNumber, epoch, account, amount, state, stakingdb); err != nil {
			return err
		}
		if amount.Cmp(common.Big0) <= 0 {
			break
		}
	}
	return nil
}

func (a *FixRestrictingPlugin) rollBackStaking(hash common.Hash, blockNumber *big.Int, account common.Address, amount *big.Int, state xcom.StateDB) error {

	iter := a.sdb.Ranking(hash, staking.CanBaseKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		return err
	}
	defer iter.Release()

	stakingdb := staking.NewStakingDBWithDB(a.sdb)

	var wrongStakings rollBackStakingInfos
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
			wrongStakings = append(wrongStakings, newWrongStakingInfo(&candidate, canAddr))
		}
	}
	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	sort.Sort(wrongStakings)
	for i := 0; i < len(wrongStakings); i++ {
		if err := wrongStakings[i].handleStaking(hash, blockNumber, epoch, amount, state, stakingdb); err != nil {
			return err
		}
		if amount.Cmp(common.Big0) <= 0 {
			break
		}
	}
	return nil
}

func newWrongStakingInfo(candidate *staking.Candidate, canAddr common.NodeAddress) *rollBackStakingInfo {
	w := &rollBackStakingInfo{
		candidate: candidate,
		canAddr:   canAddr,
	}
	w.originRestrictingAmount = new(big.Int).Add(w.candidate.RestrictingPlan, w.candidate.RestrictingPlanHes)
	w.originFreeAmount = new(big.Int).Add(w.candidate.Released, w.candidate.ReleasedHes)
	return w
}

type rollBackStakingInfo struct {
	candidate                                 *staking.Candidate
	canAddr                                   common.NodeAddress
	originRestrictingAmount, originFreeAmount *big.Int
}

func (a *rollBackStakingInfo) HandelExistStaking(hash common.Hash, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
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

func (a *rollBackStakingInfo) shouldWithdrewStaking(hash common.Hash, blockNumber *big.Int, rollBackAmount *big.Int) bool {
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

func (a *rollBackStakingInfo) decreaseStaking(hash common.Hash, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB) error {
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

func (a *rollBackStakingInfo) withdrewStaking(hash common.Hash, epoch uint64, blockNumber *big.Int, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
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

func (a *rollBackStakingInfo) handleStaking(hash common.Hash, blockNumber *big.Int, epoch uint64, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	lazyCalcStakeAmount(epoch, a.candidate.CandidateMutable)
	log.Debug("fix restricting for staking begin", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "return", rollBackAmount, "restrictingPlan",
		a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes, "released", a.candidate.Released, "releasedHes", a.candidate.ReleasedHes, "share", a.candidate.Shares)
	if a.candidate.Status.IsWithdrew() {
		//已经解质押,节点处于退出锁定期
		if err := a.HandelExistStaking(hash, epoch, rollBackAmount, state, stdb); err != nil {
			return err
		}
		log.Debug("fix restricting for staking end", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "status", a.candidate.Status, "return",
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
		log.Debug("fix restricting for staking end", "account", a.candidate.StakingAddress, "nodeID", a.candidate.NodeId.TerminalString(), "status", a.candidate.Status, "withdrewStaking",
			shouldWithdrewStaking, "return", rollBackAmount, "restrictingPlan", a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes, "released", a.candidate.Released,
			"releasedHes", a.candidate.ReleasedHes, "share", a.candidate.Shares)
	}
	return nil
}

func (a *rollBackStakingInfo) refundWrongLockFunds(rollBackAmount *big.Int, state xcom.StateDB) (*big.Int, error) {
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

type rollBackStakingInfos []*rollBackStakingInfo

func (d rollBackStakingInfos) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d rollBackStakingInfos) Len() int      { return len(d) }
func (d rollBackStakingInfos) Less(i, j int) bool {
	//按照节点质押的时间远近进行排序，最近的质押排在前面
	if d[i].candidate.StakingBlockNum > d[j].candidate.StakingBlockNum {
		return true
	}
	return false
}

type rollBackDelInfo struct {
	del       *staking.Delegation
	candidate *staking.Candidate
	//use for get staking
	stakingBlock uint64
	canAddr      common.NodeAddress

	originRestrictingAmount, originFreeAmount *big.Int
}

func (a *rollBackDelInfo) HandleDelegate(hash common.Hash, blockNumber *big.Int, epoch uint64, delAddr common.Address, rollBackAmount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) (*big.Int, error) {

	refundAmount := new(big.Int)
	wrongRestrictingAmount := new(big.Int)

	if rollBackAmount.Cmp(a.originRestrictingAmount) >= 0 {
		if ok, _ := CheckOperatingThreshold(blockNumber.Uint64(), hash, a.originFreeAmount); !ok {
			refundAmount = new(big.Int).Add(a.originFreeAmount, a.originRestrictingAmount)
		} else {
			refundAmount = new(big.Int).Set(a.originRestrictingAmount)
		}
		wrongRestrictingAmount = new(big.Int).Set(a.originRestrictingAmount)
	} else {
		resrtict := new(big.Int).Sub(a.originRestrictingAmount, rollBackAmount)
		if ok, _ := CheckOperatingThreshold(blockNumber.Uint64(), hash, new(big.Int).Add(a.originFreeAmount, resrtict)); !ok {
			refundAmount = new(big.Int).Add(a.originFreeAmount, a.originRestrictingAmount)
		} else {
			refundAmount = new(big.Int).Set(rollBackAmount)
		}
		wrongRestrictingAmount = new(big.Int).Set(rollBackAmount)
	}

	log.Debug("fix restricting for delegate begin", "account", delAddr, "currentReturn", wrongRestrictingAmount, "leftReturn", rollBackAmount, "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
		"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "share", a.candidate.Shares)

	rollBackAmount.Sub(rollBackAmount, wrongRestrictingAmount)

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
			if err := rt.ReturnWrongLockFunds(delAddr, wrongRestrictingAmount, state); nil != err {
				return nil, err
			}

			//更新锁仓信息
			refundAmount = new(big.Int).Sub(refundAmount, wrongRestrictingAmount)
			if a.del.RestrictingPlanHes.Cmp(wrongRestrictingAmount) >= 0 {
				a.del.RestrictingPlanHes.Sub(a.del.RestrictingPlanHes, wrongRestrictingAmount)
				if a.candidate.IsNotEmpty() {
					a.candidate.DelegateTotalHes.Sub(a.candidate.DelegateTotalHes, wrongRestrictingAmount)
				}
			} else {
				hes := new(big.Int).Set(a.del.RestrictingPlanHes)
				a.del.RestrictingPlanHes = new(big.Int)
				a.del.RestrictingPlan = new(big.Int).Sub(a.del.RestrictingPlan, new(big.Int).Sub(wrongRestrictingAmount, hes))
				if a.candidate.IsNotEmpty() {
					a.candidate.DelegateTotalHes.Sub(a.candidate.DelegateTotalHes, hes)
					a.candidate.DelegateTotal.Sub(a.candidate.DelegateTotal, new(big.Int).Sub(wrongRestrictingAmount, hes))
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
	log.Debug("fix restricting for delegate end", "account", delAddr, "currentReturn", wrongRestrictingAmount, "leftReturn", rollBackAmount, "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
		"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "share", a.candidate.Shares)
	return nil, nil
}

type rollBackDelInfos []*rollBackDelInfo

func (d rollBackDelInfos) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d rollBackDelInfos) Len() int      { return len(d) }
func (d rollBackDelInfos) Less(i, j int) bool {
	//根据委托节点的分红比例从小到大排序，如果委托比例相同，根据节点id从小到大排序
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
