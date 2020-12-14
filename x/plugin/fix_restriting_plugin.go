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

type wrongAddr struct {
	addr   common.Address
	amount *big.Int
}

func (a *FixRestrictingPlugin) fix(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	for _, wrong := range []wrongAddr{} {
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

	var wrongDels wrongDelInfos

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

		wrongDel := new(wrongDelInfo)
		wrongDel.del = &del
		wrongDel.candidate = can
		wrongDel.stakingBlock = stakingBlock
		wrongDels = append(wrongDels, wrongDel)
	}
	sort.Sort(wrongDels)
	epoch := xutil.CalculateEpoch(blockNumber.Uint64())
	stakingdb := staking.NewStakingDBWithDB(a.sdb)
	for i := 0; i < len(wrongDels); i++ {
		if _, err := wrongDels[i].WithdrewDelegate(hash, blockNumber, epoch, account, amount, state, stakingdb); err != nil {
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

	var wrongStakings wrongStakingInfos
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
			wrongStakings = append(wrongStakings, &wrongStakingInfo{
				candidate: &candidate,
				canAddr:   canAddr,
			})
		}
	}
	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	sort.Sort(wrongStakings)
	for i := 0; i < len(wrongStakings); i++ {
		if err := wrongStakings[i].withdrewStaking(hash, blockNumber, epoch, amount, state, stakingdb); err != nil {
			return err
		}
		if amount.Cmp(common.Big0) <= 0 {
			break
		}
	}
	return nil
}

type wrongStakingInfo struct {
	candidate *staking.Candidate
	canAddr   common.NodeAddress
}

func (a *wrongStakingInfo) withdrewStaking(hash common.Hash, blockNumber *big.Int, epoch uint64, amount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) error {
	lazyCalcStakeAmount(epoch, a.candidate.CandidateMutable)
	res := new(big.Int).Add(a.candidate.RestrictingPlan, a.candidate.RestrictingPlanHes)
	release := new(big.Int).Add(a.candidate.Released, a.candidate.ReleasedHes)
	shouldWithdrewStaking := false

	if res.Cmp(amount) >= 0 {
		left := new(big.Int).Add(release, new(big.Int).Sub(res, amount))
		if ok, _ := CheckStakeThreshold(blockNumber.Uint64(), hash, left); !ok {
			shouldWithdrewStaking = true
		}
	} else {
		left := new(big.Int).Sub(release, res)
		if ok, _ := CheckStakeThreshold(blockNumber.Uint64(), hash, left); !ok {
			shouldWithdrewStaking = true
		}
	}
	log.Debug("fix restricting for staking begin", "account", a.candidate.StakingAddress, "withdrewStaking", shouldWithdrewStaking, "return", amount, "restrictingPlan", a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes,
		"release", release, "share", a.candidate.Shares)
	if shouldWithdrewStaking {
		if err := stdb.DelCanPowerStore(hash, a.candidate); nil != err {
			return err
		}

		// Direct return of money during the hesitation period
		// Return according to the way of coming
		if a.candidate.ReleasedHes.Cmp(common.Big0) > 0 {
			state.AddBalance(a.candidate.StakingAddress, a.candidate.ReleasedHes)
			state.SubBalance(vm.StakingContractAddr, a.candidate.ReleasedHes)
			a.candidate.ReleasedHes = new(big.Int).SetInt64(0)
		}

		if res.Cmp(amount) >= 0 {
			if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, amount, state); nil != err {
				return err
			}
			if a.candidate.RestrictingPlanHes.Cmp(amount) >= 0 {
				a.candidate.RestrictingPlanHes = new(big.Int).Sub(a.candidate.RestrictingPlanHes, amount)
			} else {
				a.candidate.RestrictingPlan = new(big.Int).Sub(a.candidate.RestrictingPlan, new(big.Int).Sub(amount, a.candidate.RestrictingPlanHes))
				a.candidate.RestrictingPlanHes = new(big.Int).SetInt64(0)
			}
			amount.SetInt64(0)
		} else {
			if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, res, state); nil != err {
				return err
			}
			a.candidate.RestrictingPlanHes = new(big.Int).SetInt64(0)
			a.candidate.RestrictingPlan = new(big.Int).SetInt64(0)
			amount.Sub(amount, res)
		}

		if a.candidate.RestrictingPlanHes.Cmp(common.Big0) > 0 {
			err := rt.ReturnLockFunds(a.candidate.StakingAddress, a.candidate.RestrictingPlanHes, state)
			if nil != err {
				return err
			}
			a.candidate.RestrictingPlanHes = new(big.Int).SetInt64(0)
		}

		//todo add a.candidate.Status check
		if a.candidate.Status.IsValid() {
			if a.candidate.Released.Cmp(common.Big0) > 0 || a.candidate.RestrictingPlan.Cmp(common.Big0) > 0 {
				if err := stk.addErrorAccountUnStakeItem(blockNumber.Uint64(), hash, a.candidate.NodeId, a.canAddr, a.candidate.StakingBlockNum); nil != err {
					return err
				}
				// sub the account staking Reference Count
				if err := stdb.SubAccountStakeRc(hash, a.candidate.StakingAddress); nil != err {
					return err
				}
			}
		}

		a.candidate.CleanShares()
		a.candidate.Status |= staking.Invalided | staking.Withdrew

		a.candidate.StakingEpoch = uint32(epoch)

		if a.candidate.Released.Cmp(common.Big0) > 0 || a.candidate.RestrictingPlan.Cmp(common.Big0) > 0 {
			if err := stdb.SetCanMutableStore(hash, a.canAddr, a.candidate.CandidateMutable); nil != err {
				return err
			}
		} else {
			if err := stdb.DelCandidateStore(hash, a.canAddr); nil != err {
				return err
			}
		}

	} else {

		realSub := new(big.Int)
		if res.Cmp(amount) >= 0 {
			if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, amount, state); nil != err {
				return err
			}
			realSub.Set(amount)
			if a.candidate.RestrictingPlanHes.Cmp(amount) >= 0 {
				a.candidate.RestrictingPlanHes = new(big.Int).Sub(a.candidate.RestrictingPlanHes, amount)
			} else {
				a.candidate.RestrictingPlan = new(big.Int).Sub(a.candidate.RestrictingPlan, new(big.Int).Sub(amount, a.candidate.RestrictingPlanHes))
				a.candidate.RestrictingPlanHes = new(big.Int).SetInt64(0)
			}
			amount = new(big.Int).SetInt64(0)
		} else {
			if err := rt.ReturnWrongLockFunds(a.candidate.StakingAddress, res, state); nil != err {
				return err
			}
			realSub.Set(res)
			a.candidate.RestrictingPlanHes = new(big.Int)
			a.candidate.RestrictingPlan = new(big.Int)
			amount = new(big.Int).Sub(amount, res)
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
	}
	log.Debug("fix restricting for staking end", "account", a.candidate.StakingAddress, "return", amount, "restrictingPlan", a.candidate.RestrictingPlan, "restrictingPlanRes", a.candidate.RestrictingPlanHes,
		"release", release, "share", a.candidate.Shares)
	return nil
}

type wrongStakingInfos []*wrongStakingInfo

func (d wrongStakingInfos) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d wrongStakingInfos) Len() int      { return len(d) }
func (d wrongStakingInfos) Less(i, j int) bool {
	if d[i].candidate.StakingBlockNum > d[j].candidate.StakingBlockNum {
		return true
	}
	return false
}

type wrongDelInfo struct {
	del       *staking.Delegation
	candidate *staking.Candidate
	//use for get staking
	stakingBlock uint64
	canAddr      common.NodeAddress
}

func (a *wrongDelInfo) WithdrewDelegate(hash common.Hash, blockNumber *big.Int, epoch uint64, delAddr common.Address, amount *big.Int, state xcom.StateDB, stdb *staking.StakingDB) (*big.Int, error) {
	restrictingPlan := new(big.Int).Add(a.del.RestrictingPlan, a.del.RestrictingPlanHes)
	released := new(big.Int).Add(a.del.Released, a.del.ReleasedHes)

	refundAmount := new(big.Int)
	wrongRestrictingAmount := new(big.Int)

	if amount.Cmp(restrictingPlan) >= 0 {
		if ok, _ := CheckOperatingThreshold(blockNumber.Uint64(), hash, released); !ok {
			refundAmount = new(big.Int).Add(released, restrictingPlan)
		} else {
			refundAmount = restrictingPlan
		}
		wrongRestrictingAmount = new(big.Int).Set(restrictingPlan)
	} else {
		resrtict := new(big.Int).Sub(restrictingPlan, amount)
		if ok, _ := CheckOperatingThreshold(blockNumber.Uint64(), hash, new(big.Int).Add(released, resrtict)); !ok {
			refundAmount = new(big.Int).Add(released, restrictingPlan)
		} else {
			refundAmount = new(big.Int).Set(amount)
		}
		wrongRestrictingAmount = new(big.Int).Set(amount)
	}

	log.Debug("fix restricting for delegate begin", "account", delAddr, "currentReturn", wrongRestrictingAmount, "leftReturn", amount, "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
		"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "share", a.candidate.Shares)

	amount.Sub(amount, wrongRestrictingAmount)

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
			if err := rt.ReturnWrongLockFunds(delAddr, wrongRestrictingAmount, state); nil != err {
				return nil, err
			}
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
		total := new(big.Int).Add(restrictingPlan, released)
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
	log.Debug("fix restricting for delegate end", "account", delAddr, "currentReturn", wrongRestrictingAmount, "leftReturn", amount, "restrictingPlan", a.del.RestrictingPlan, "restrictingPlanRes", a.del.RestrictingPlanHes,
		"release", a.del.Released, "releaseHes", a.del.ReleasedHes, "share", a.candidate.Shares)
	return nil, nil
}

type wrongDelInfos []*wrongDelInfo

func (d wrongDelInfos) Swap(i, j int) { d[i], d[j] = d[j], d[i] }
func (d wrongDelInfos) Len() int      { return len(d) }
func (d wrongDelInfos) Less(i, j int) bool {
	//todo add quote
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
