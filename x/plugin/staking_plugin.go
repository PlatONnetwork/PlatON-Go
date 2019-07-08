package plugin

import (
	"encoding/hex"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
	"sort"
	"strconv"
	"sync"
)

type StakingPlugin struct {
	db   *staking.StakingDB
	once sync.Once
}

var stk *StakingPlugin

var (
	AccountVonNotEnough        = common.NewBizError("The von of account is not enough")
	DelegateVonNotEnough       = common.NewBizError("The von of delegate is not enough")
	WithdrewDelegateVonCalcErr = common.NewBizError("withdrew delegate von calculate err")
	ParamsErr                  = common.NewBizError("the fn params err")
	ProcessVersionErr          = common.NewBizError("The version of the relates node's process is too low")
	BlockNumberDisordered 	   = common.NewBizError("The blockNumber is disordered")
	VonAmountNotRight		   = common.NewBizError("The amount of von is not right")
	CandidateNotExist 		   = common.NewBizError("The candidate is not exist")
	ValidatorNotExist 		   = common.NewBizError("The validator is not exist")
	ElectionErr				   = common.NewBizError("Election failure")
)

const (
	FreeOrigin     = 0
	RestrictingPlanOrigin = 1

	PriviosRound = uint(0)
	CurrentRound = uint(1)
	NextRound = uint(2)

	QueryStartIrr = true
	QueryStartNotIrr = false
	


)

// Instance a global StakingPlugin
func StakingInstance() *StakingPlugin {
	if nil == stk {
		stk = &StakingPlugin{
			db: staking.NewStakingDB(),
		}
	}
	return stk
}

func (sk *StakingPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {

	return nil
}

func (sk *StakingPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {

	epoch := xutil.CalculateEpoch(header.Number.Uint64())

	if xutil.IsSettlementPeriod(header.Number.Uint64()) {
		// handle UnStaking Item
		err := sk.HandleUnCandidateItem(state, blockHash, epoch)
		if nil != err {
			log.Error("Failed to call HandleUnCandidateReq on stakingPlugin EndBlock", "blockHash",
	blockHash.Hex(), "blockNumber", header.Number.Uint64(), "err", err)
			return err //  TODO common.NewSysError(err.Error())
		}

		// Election next epoch validators
		if err := sk.ElectNextVerifierList(blockHash, header.Number.Uint64()); nil != err {
			return err
		}
	}

	if xutil.IsElection(header.Number.Uint64()) {
		// ELection next round validators
		err := sk.Election(blockHash, header.Number.Uint64())
		if nil != err {
			log.Error("Failed to call Election on stakingPlugin EndBlock", "blockHash", blockHash.Hex(),
	"blockNumber", header.Number.Uint64(), "err", err)
			return err
		}
	}

	if xutil.IsSwitch(header.Number.Uint64()) {
		// Switch previous, current and next round validators
		err := sk.Switch(blockHash, header.Number.Uint64())
		if nil != err {
			log.Error("Failed to call Switch on stakingPlugin EndBlock", "blockHash", blockHash.Hex(),
	"blockNumber", header.Number.Uint64(), "err", err)
			return err
		}
	}

	return nil
}

func (sk *StakingPlugin) Confirmed(block *types.Block) error {

	if xutil.IsElection(block.NumberU64()) {

		next, err := sk.db.GetNextValidatorListByBlockHash(block.Hash())
		if nil != err {
			return err
		}

		// TODO Notify P2P module
		_ = next

	}

	return nil
}

func (sk *StakingPlugin) GetCandidateInfo (blockHash common.Hash, addr common.Address) (*staking.Candidate, error) {

	return sk.db.GetCandidateStore(blockHash, addr)
}

func (sk *StakingPlugin) GetCandidateInfoByIrr (addr common.Address) (*staking.Candidate, error) {
	return sk.db.GetCandidateStoreByIrr(addr)
}

func (sk *StakingPlugin) CreateCandidate (state xcom.StateDB, blockHash common.Hash, blockNumber,
	amount *big.Int, processVersion uint32, typ uint16, addr common.Address, can *staking.Candidate) error {

	// Query current active version
	curr_version := govPlugin.GetActiveVersion(state)

	var isDeclareVersion bool

	if processVersion < curr_version {
		return ProcessVersionErr
	} else if processVersion > curr_version {
		isDeclareVersion = true
	}
	can.ProcessVersion = curr_version

	// from account free von
	if typ == FreeOrigin {

		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to CreateCandidate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(),
				"blockHash", blockHash.Hex(), "originVon", origin, "stakingVon", amount)
			return AccountVonNotEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)
		can.ReleasedHes = amount

	} else if typ == RestrictingPlanOrigin { //  from account RestrictingPlan von

		err := RestrictingPtr.PledgeLockFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to CreateCandidate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"err", err)
			return err
		}
		can.RestrictingPlanHes = amount
	}

	can.StakingEpoch = uint32(xutil.CalculateEpoch(blockNumber.Uint64()))

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Put Can info 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Put Can power 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	if isDeclareVersion {
		// Declare new Version
		err := govPlugin.DeclareVersion(can.StakingAddress, can.NodeId, processVersion, blockHash, blockNumber.Uint64(), state)
		if nil != err {
			log.Error("Call CreateCandidate with govplugin DelareVersion failed", "err", err)
		}
	}
	return nil
}

func (sk *StakingPlugin) EditorCandidate (blockHash common.Hash, blockNumber *big.Int, can *staking.Candidate) error {
	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to EditorCandidate on stakingPlugin: Put Can info 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) IncreaseStaking (state xcom.StateDB, blockHash common.Hash, blockNumber,
	amount *big.Int, typ uint16, can *staking.Candidate) error {

	xcom.PrintObject("IncreaseStaking 进来的时候can:", can)

	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	if typ == FreeOrigin {
		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to EditorCandidate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "account", can.StakingAddress.Hex(),
				"originVon", origin, "stakingVon", can.ReleasedHes)
			return AccountVonNotEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)

		can.ReleasedHes = new(big.Int).Add(can.ReleasedHes, amount)

	} else {

		err := RestrictingPtr.PledgeLockFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to EditorCandidate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"err", err)
			return err
		}

		can.RestrictingPlanHes = new(big.Int).Add(can.RestrictingPlanHes, amount)
	}

	can.StakingEpoch = uint32(epoch)

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to EditorCandidate on stakingPlugin: Del Can old power failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	can.Shares = new(big.Int).Add(can.Shares, amount)

	// set new power of can
	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to EditorCandidate on stakingPlugin: Put Can power 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}


	xcom.PrintObject("IncreaseStaking 存储之前can:", can)


	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to EditorCandidate on stakingPlugin: Put Can info 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) WithdrewCandidate(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	can *staking.Candidate) error {
	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to WithdrewCandidate on stakingPlugin: Del Can old power failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	if err := sk.withdrewStakeAmount(state, blockHash, blockNumber.Uint64(), epoch, addr, can); nil != err {
		return err
	}

	can.StakingEpoch = uint32(epoch)

	if can.Released.Cmp(common.Big0) > 0 || can.RestrictingPlan.Cmp(common.Big0) > 0 {

		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Failed to WithdrewCandidate on stakingPlugin: Put Can info 2 db failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	} else {
		if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
			log.Error("Failed to WithdrewCandidate on stakingPlugin: Del Can info failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}
	return nil
}

func (sk *StakingPlugin) withdrewStakeAmount(state xcom.StateDB, blockHash common.Hash, blockNumber, epoch uint64,
	addr common.Address, can *staking.Candidate) error {

	// Direct return of money during the hesitation period
	// Return according to the way of coming
	if can.ReleasedHes.Cmp(common.Big0) > 0 {
		state.AddBalance(can.StakingAddress, can.ReleasedHes)
		state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)
		can.Shares = new(big.Int).Sub(can.Shares, can.ReleasedHes)
		can.ReleasedHes = common.Big0
	}

	if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {

		err := RestrictingPtr.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
		if nil != err {
			log.Error("Failed to WithdrewCandidate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"err", err)
			return err
		}

		can.Shares = new(big.Int).Sub(can.Shares, can.RestrictingPlanHes)
		can.RestrictingPlanHes = common.Big0
	}

	if can.Released.Cmp(common.Big0) > 0 || can.RestrictingPlan.Cmp(common.Big0) > 0 {
		if err := sk.db.AddUnStakeItemStore(blockHash, epoch, addr); nil != err {
			return err
		}
	}
	can.Status |= staking.Invalided

	return nil
}

func (sk *StakingPlugin) HandleUnCandidateItem(state xcom.StateDB, blockHash common.Hash, epoch uint64) error {

	releaseEpoch := epoch - xcom.UnStakeFreezeRatio

	unStakeCount, err := sk.db.GetUnStakeCountStore(blockHash, releaseEpoch)
	if nil != err {
		return err
	}

	if unStakeCount == 0 {
		return nil
	}

	filterAddr := make(map[common.Address]struct{})

	for index := 1; index <= int(unStakeCount); index++ {
		addr, err := sk.db.GetUnStakeItemStore(blockHash, releaseEpoch, uint64(index))
		if nil != err {
			return err
		}

		if _, ok := filterAddr[addr]; ok {
			if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
				return err
			}
			continue
		}

		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			return err
		}

		if nil == can {
			// This should not be nil
			continue
		}

		// Already deleted power
		/*// First delete the weight information
		if err := sk.db.delCanPowerStore(blockHash, can); nil != err {
			return false, err
		}*/

		// Second handle balabala ...
		if err := sk.handleUnStake(state, blockHash, epoch, addr, can); nil != err {
			return err
		}

		if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
			return err
		}

		filterAddr[addr] = struct{}{}
	}

	if err := sk.db.DelUnStakeCountStore(blockHash, releaseEpoch); nil != err {
		return err
	}

	return nil
}

func (sk *StakingPlugin) handleUnStake(state xcom.StateDB, blockHash common.Hash, epoch uint64,
	addr common.Address, can *staking.Candidate) error {

	lazyCalcStakeAmount(epoch, can)

	// Direct return of money during the hesitation period
	// Return according to the way of coming
	if can.Released.Cmp(common.Big0) > 0 {
		state.AddBalance(can.StakingAddress, can.Released)
		state.SubBalance(vm.StakingContractAddr, can.Released)
	}

	if can.RestrictingPlan.Cmp(common.Big0) > 0 {
		err := RestrictingPtr.ReturnLockFunds(can.StakingAddress, can.RestrictingPlan, state)
		if nil != err {
			log.Error("Failed to HandleUnCandidateReq on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"err", err)
			return err
		}
	}

	// delete can info
	if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
		return err
	}

	return nil
}

func (sk *StakingPlugin) GetDelegateInfo(blockHash common.Hash, delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*staking.Delegation, error) {
	return sk.db.GetDelegateStore(blockHash, delAddr, nodeId, stakeBlockNumber)
}

func (sk *StakingPlugin) GetDelegateExInfo (blockHash common.Hash, delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*staking.DelegationEx, error) {

	del, err := sk.db.GetDelegateStore(blockHash, delAddr, nodeId, stakeBlockNumber)
	if nil != err {
		return nil, err
	}
	return &staking.DelegationEx{
		Addr: delAddr,
		NodeId: nodeId,
		StakingBlockNum: stakeBlockNumber,
		Delegation: *del,
	}, nil
}

func (sk *StakingPlugin) GetDelegateInfoByIrr (delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*staking.Delegation, error) {

	return sk.db.GetDelegateStoreByIrr(delAddr, nodeId, stakeBlockNumber)
}

func (sk *StakingPlugin) GetDelegateExInfoByIrr (delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*staking.DelegationEx, error) {

	del, err := sk.db.GetDelegateStoreByIrr(delAddr, nodeId, stakeBlockNumber)
	if nil != err {
		return nil, err
	}
	return &staking.DelegationEx{
		Addr: delAddr,
		NodeId: nodeId,
		StakingBlockNum: stakeBlockNumber,
		Delegation: *del,
	}, nil
}



func (sk *StakingPlugin) Delegate(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	delAddr common.Address, del *staking.Delegation, can *staking.Candidate, typ uint16, amount *big.Int) error {


	xcom.PrintObject("进来的can:", can)
	xcom.PrintObject("进来的del:", del)

	pubKey, _ := can.NodeId.Pubkey()
	canAddr := crypto.PubkeyToAddress(*pubKey)

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcDelegateAmount(epoch, del)

	if typ == FreeOrigin { // from account free von

		origin := state.GetBalance(delAddr)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to Delegate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "originVon", origin,
				"stakingVon", can.ReleasedHes)
			return AccountVonNotEnough
		}
		state.SubBalance(delAddr, amount)
		state.AddBalance(vm.StakingContractAddr, amount)

		del.ReleasedHes = new(big.Int).Add(del.ReleasedHes, amount)

	} else if typ == RestrictingPlanOrigin { //  from account RestrictingPlan von

		err := RestrictingPtr.PledgeLockFunds(delAddr, amount, state)
		if nil != err {
			log.Error("Failed to Delegate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"err", err)
			return err
		}

		del.RestrictingPlanHes = new(big.Int).Add(del.RestrictingPlanHes, amount)

	}

	del.DelegateEpoch = uint32(epoch)

	xcom.PrintObject("修改完存储之前的del:", del)

	// set new delegate info
	if err := sk.db.SetDelegateStore(blockHash, delAddr, can.NodeId, can.StakingBlockNum, del); nil != err {
		return err
	}

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		return err
	}

	// add the candidate power
	can.Shares = new(big.Int).Add(can.Shares, amount)

	xcom.PrintObject("修改完存储之前的can:", can)

	// set new power of can
	if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
		return err
	}

	// update can info about Shares
	if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
		return err
	}
	return nil
}

func (sk *StakingPlugin) WithdrewDelegate(state xcom.StateDB, blockHash common.Hash, blockNumber, amount *big.Int,
	delAddr common.Address, nodeId discover.NodeID, stakingBlockNum uint64, del *staking.Delegation) error {

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to WithdrewDelegate on stakingPlugin", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	can, err := sk.db.GetCandidateStore(blockHash, canAddr)
	if nil != err {
		return err
	}

	xcom.PrintObject("进来时的can:", can)
	xcom.PrintObject("进来时的del:", del)

	aboutRelease := new(big.Int).Add(del.Released, del.ReleasedHes)
	aboutRestrictingPlan := new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
	total := new(big.Int).Add(aboutRelease, aboutRestrictingPlan)


	lazyCalcDelegateAmount(epoch, del)



	/**
	inner Fn
	*/
	subDelegateFn := func(source, sub *big.Int) (*big.Int, *big.Int) {
		state.AddBalance(delAddr, sub)
		state.SubBalance(vm.StakingContractAddr, sub)
		return new(big.Int).Sub(source, sub), common.Big0
	}

	refundFn := func(remain, aboutRelease, aboutRestrictingPlan *big.Int) (*big.Int, *big.Int, *big.Int, error) {
		// When remain is greater than or equal to del.ReleasedHes/del.Released
		if remain.Cmp(common.Big0) > 0 {
			if remain.Cmp(aboutRelease) >= 0 && aboutRelease.Cmp(common.Big0) > 0 {

				remain, aboutRelease = subDelegateFn(remain, aboutRelease)

			} else if remain.Cmp(aboutRelease) < 0 {
				// When remain is less than or equal to del.ReleasedHes/del.Released
				aboutRelease, remain = subDelegateFn(aboutRelease, remain)
			}
		}

		if remain.Cmp(common.Big0) > 0 {

			// When remain is greater than or equal to del.RestrictingPlanHes/del.RestrictingPlan
			if remain.Cmp(aboutRestrictingPlan) >= 0 && aboutRestrictingPlan.Cmp(common.Big0) > 0 {

				err := RestrictingPtr.ReturnLockFunds(can.StakingAddress, aboutRestrictingPlan, state)
				if nil != err {
					log.Error("Failed to WithdrewDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"err", err)
					return remain, aboutRelease, aboutRestrictingPlan, err
				}


				remain = new(big.Int).Sub(remain, aboutRestrictingPlan)
				aboutRestrictingPlan = common.Big0
			} else if remain.Cmp(aboutRestrictingPlan) < 0 {
				// When remain is less than or equal to del.RestrictingPlanHes/del.RestrictingPlan


				err := RestrictingPtr.ReturnLockFunds(can.StakingAddress, remain, state)
				if nil != err {
					log.Error("Failed to WithdrewDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"err", err)
					return remain, aboutRelease, aboutRestrictingPlan, err
				}

				aboutRestrictingPlan = new(big.Int).Sub(aboutRestrictingPlan, remain)
				remain = common.Big0
			}
		}

		return remain, aboutRelease, aboutRestrictingPlan, nil
	}

	del.DelegateEpoch = uint32(epoch)

	switch {

	// When the related candidate info does not exist
	case nil == can, nil != can && stakingBlockNum < can.StakingBlockNum,
	nil != can && stakingBlockNum == can.StakingBlockNum && staking.Is_Invalid(can.Status):

		if total.Cmp(amount) < 0 {
			return common.BizErrorf("withdrewDelegate err: %s, delegate von: %s, withdrew von: %s",
				DelegateVonNotEnough.Error(), total.String(), amount.String())
		}

		remain := amount

		/**
		handle delegate on HesitateRatio
		*/
		remain, rtmp, ltmp, err := refundFn(remain, del.ReleasedHes, del.RestrictingPlanHes)
		if nil != err {
			return err
		}
		del.ReleasedHes, del.RestrictingPlanHes =  rtmp, ltmp
		/**
		handle delegate on EffectiveRatio
		*/
		if remain.Cmp(common.Big0) > 0 {
			remain, rtmp, ltmp, err = refundFn(remain, del.Released, del.RestrictingPlan)
			if nil != err {
				return err
			}
			del.Released, del.RestrictingPlan = rtmp, ltmp
		}

		if remain.Cmp(common.Big0) != 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: sub delegate von calculation error",
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String())
			return WithdrewDelegateVonCalcErr
		}

		if total.Cmp(amount) == 0 {
			if err := sk.db.DelDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum); nil != err {
				return err
			}

		}else {
			sub := new(big.Int).Sub(total, del.Reduction)

			if sub.Cmp(amount) < 0 {
				tmp := new(big.Int).Sub(amount, sub)
				del.Reduction = new(big.Int).Sub(del.Reduction, tmp)
			}

			if err := sk.db.SetDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum, del); nil != err {
				return err
			}
		}

	// Illegal parameter
	case nil != can && stakingBlockNum > can.StakingBlockNum:
		log.Error("Failed to WithdrewDelegate on stakingPlugin: the stakeBlockNum err",
			"blockHash", blockHash.Hex(), "fn.stakeBlockNum", stakingBlockNum, "can.stakeBlockNum", can.StakingBlockNum)
		return ParamsErr

	// When the delegate is normally revoked
	case nil != can && stakingBlockNum == can.StakingBlockNum && staking.Is_Valid(can.Status):

		total = new(big.Int).Sub(total, del.Reduction)

		if total.Cmp(amount) < 0 {
			return common.BizErrorf("withdrewDelegate err: %s, delegate von: %s, withdrew von: %s",
				DelegateVonNotEnough.Error(), total.String(), amount.String())
		}

		remain := amount

		/**
		handle delegate on HesitateRatio
		*/
		//var flag bool
		//var er error
		remain, rtmp, ltmp, err := refundFn(remain, del.ReleasedHes, del.RestrictingPlanHes)
		if nil != err {
			return err
		}
		del.ReleasedHes, del.RestrictingPlanHes = rtmp, ltmp
		/**
		handle delegate on EffectiveRatio
		*/
		if remain.Cmp(common.Big0) > 0 {

			// add a UnDelegateItem
			sk.db.AddUnDelegateItemStore(blockHash, delAddr, nodeId, epoch, stakingBlockNum, remain)
			del.Reduction = new(big.Int).Add(del.Reduction, remain)
		}

		if err := sk.db.SetDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum, del); nil != err {
			return err
		}
	}




	// delete old can power
	if nil != can && stakingBlockNum == can.StakingBlockNum && staking.Is_Valid(can.Status) {
		if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
			return err
		}

		can.Shares = new(big.Int).Sub(can.Shares, amount)

		if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
			return err
		}

		if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
			return err
		}
	}
	xcom.PrintObject("撤销之后的can:", can)
	xcom.PrintObject("撤销之后的del:", del)
	return nil
}

func (sk *StakingPlugin) HandleUnDelegateReq(state xcom.StateDB, blockHash common.Hash, epoch uint64) error {
	releaseEpoch := epoch - xcom.ActiveUnDelegateFreezeRatio

	unDelegateCount, err := sk.db.GetUnDelegateCountStore(blockHash, releaseEpoch)
	if nil != err {
		return err
	}

	if unDelegateCount == 0 {
		return nil
	}

	//filterAddr := make(map[string]struct{})

	for index := 1; index <= int(unDelegateCount); index++ {
		unDelegateItem, err := sk.db.GetUnDelegateItemStore(blockHash, releaseEpoch, uint64(index))
		if nil != err {
			return err
		}

		//if _, ok := filterAddr[fmt.Sprint(unDelegateItem.KeySuffix)]; ok {
		//	continue
		//}

		del, err := sk.db.GetDelegateStoreBySuffix(blockHash, unDelegateItem.KeySuffix)
		if nil != err {
			return err
		}

		if nil == del {
			// This maybe be nil
			continue
		}

		if err := sk.handleUnDelegate(state, blockHash, epoch, unDelegateItem, del); nil != err {
			return err
		}

		//filterAddr[fmt.Sprint(unDelegateItem.KeySuffix)] = struct{}{}
	}

	return nil
}

func (sk *StakingPlugin) handleUnDelegate(state xcom.StateDB, blockHash common.Hash, epoch uint64,
	unDel *staking.UnDelegateItem, del *staking.Delegation) error {

	// del addr
	delAddrByte := unDel.KeySuffix[0:common.AddressLength]
	delAddr := common.BytesToAddress(delAddrByte)

	nodeIdLen := discover.NodeIDBits / 8

	nodeIdByte := unDel.KeySuffix[common.AddressLength : common.AddressLength + nodeIdLen]
	nodeId := discover.MustBytesID(nodeIdByte)

	//
	stakeBlockNum := unDel.KeySuffix[common.AddressLength + nodeIdLen:]
	num := common.BytesToUint64(stakeBlockNum)

	lazyCalcDelegateAmount(epoch, del)


	amount := unDel.Amount


	if amount.Cmp(del.Reduction) >= 0 { // full withdrawal
		state.SubBalance(vm.StakingContractAddr, del.Released)
		state.AddBalance(delAddr, del.Released)

		err := RestrictingPtr.ReturnLockFunds(delAddr, del.RestrictingPlan, state)
		if nil != err {
			log.Error("Failed to HandleUnDelegateReq on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"err", err)
			return err
		}

		if err := sk.db.DelDelegateStoreBySuffix(blockHash, unDel.KeySuffix); nil != err {
			return err
		}

	}else { //few withdrawal

		remain := amount

		if remain.Cmp(del.Released) >= 0 {
			state.SubBalance(vm.StakingContractAddr, del.Released)
			state.AddBalance(delAddr, del.Released)
			del.Released = common.Big0; remain = new(big.Int).Sub(remain, del.Released)
		}else {
			state.SubBalance(vm.StakingContractAddr, amount)
			state.AddBalance(delAddr, amount)
			del.Released = new(big.Int).Sub(del.Released, remain); remain = common.Big0
		}

		if remain.Cmp(common.Big0) > 0 {

			if remain.Cmp(del.RestrictingPlan) >= 0 {

				err := RestrictingPtr.ReturnLockFunds(delAddr, del.RestrictingPlan, state)
				if nil != err {
					log.Error("Failed to HandleUnDelegateReq on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"err", err)
					return err
				}


				del.RestrictingPlan = common.Big0; remain = new(big.Int).Sub(remain, del.RestrictingPlan)
			}else {

				err := RestrictingPtr.ReturnLockFunds(delAddr, remain, state)
				if nil != err {
					log.Error("Failed to HandleUnDelegateReq on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"err", err)
					return err
				}

				del.RestrictingPlan = new(big.Int).Sub(del.RestrictingPlan, remain); remain = common.Big0
			}
		}

		if remain.Cmp(common.Big0) > 0 {
			log.Error("Failed to call handleUnDelegate", "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
				"nodeId", nodeId.String(), "stakeBlockNumber", num)
			return VonAmountNotRight
		}

		del.Reduction = new(big.Int).Sub(del.Reduction, amount)

		del.DelegateEpoch = uint32(epoch)

		if err := sk.db.SetDelegateStoreBySuffix(blockHash, unDel.KeySuffix, del); nil != err {
			return err
		}
	}

	return nil
}




func (sk *StakingPlugin) ElectNextVerifierList(blockHash common.Hash, blockNumber uint64) error {

	log.Info("Call ElectNextVerifierList", "blockNumber", blockNumber, "blockHash", blockHash.Hex())

	old_verifierArr, err := sk.db.GetVerifierListByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to ElectNextVerifierList", "blockNumber", blockNumber, "blockHash",
			blockHash.Hex(), "err", err)
		return err
	}

	/*if nil != old_verifierArr {

	}*/

	if old_verifierArr.End != blockNumber {
		log.Error("Failed to ElectNextVerifierList: this blockNumber invalid", "Old Epoch End blockNumber",
			old_verifierArr.End, "Current blockNumber", blockNumber)
		return common.BizErrorf("The BlockNumber invalid, Old Epoch End blockNumber: %d, Current blockNumber: %d",
			old_verifierArr.End, blockNumber)
	}


	iter := sk.db.IteratorCandidatePowerByBlockHash(blockHash, int(xcom.EpochValidatorNum))

	start := old_verifierArr.End + 1
	end := old_verifierArr.End + xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
		Start: 	start,
		End: 	end,
	}

	queue := make(staking.ValidatorQueue, 0)

	for count := 0; iter.Valid() && count < int(xcom.EpochValidatorNum); iter.Next() {
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := sk.db.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			log.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return err
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress: addr,
			NodeId: 	 can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
	}

	if len(queue) == 0 {
		panic(common.BizErrorf("Failed to ElectNextVerifierList: Select zero validators~"))
		//return true, common.BizErrorf("Failed to ElectNextVerifierList: Select zero validators~")
	}

	new_verifierArr.Arr = queue

	err = sk.db.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		return err
	}
	return nil
}

func (sk *StakingPlugin) GetVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.ValidatorExQueue, error) {

	var verifierList *staking.Validator_array
	if !isCommit {
		arr, err := sk.db.GetVerifierListByBlockHash(blockHash)
		if nil != err {
			return nil, err
		}
		verifierList = arr
	}else {
		arr, err := sk.db.GetVerifierListByIrr()
		if nil != err {
			return nil, err
		}
		verifierList = arr
	}


	if !isCommit && (blockNumber < verifierList.Start || blockNumber > verifierList.End) {
		return nil, common.BizErrorf("GetVerifierList failed: %s, start: %d, end: %d, currentNumer: %d",
			BlockNumberDisordered.Error(), verifierList.Start, verifierList.End, blockNumber)
	}

	queue := make(staking.ValidatorExQueue, len(verifierList.Arr))

	for i, v := range verifierList.Arr {

		var can *staking.Candidate
		if !isCommit {
			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				return nil, err
			}
			can = c
		}else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				return nil, err
			}
			can = c
		}

		shares, _ := new(big.Int).SetString(v.StakingWeight[1], 10)

		valEx := &staking.ValidatorEx{
			NodeId: can.NodeId,
			StakingAddress: can.StakingAddress,
			BenifitAddress: can.BenifitAddress,
			StakingTxIndex: can.StakingTxIndex,
			ProcessVersion: can.ProcessVersion,
			StakingBlockNum: can.StakingBlockNum,
			Shares: shares,
			Description: can.Description,
			ValidatorTerm: v.ValidatorTerm,
		}
		queue[i] = valEx
	}

	return queue, nil
}

func (sk *StakingPlugin) GetVerifierListFake(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.CandidateQueue, error) {

	cand := &staking.Candidate{
		discover.NodeID{0x11},
		common.HexToAddress("0x11"),
		common.HexToAddress("0x11"),
		1,
		1,
		1,
		100,
		100,
		big.NewInt(100),
		big.NewInt(100),
		big.NewInt(100),
		big.NewInt(100),
		big.NewInt(100),
		staking.Description{"", "", "",""},
	}
	que := staking.CandidateQueue{cand}

	return que, nil
}

func (sk *StakingPlugin) IsCurrVerifier(blockHash common.Hash, nodeId discover.NodeID, isCommit bool) (bool, error) {

	var verifierList *staking.Validator_array

	if !isCommit {
		arr, err := sk.db.GetVerifierListByBlockHash(blockHash)
		if nil != err {
			return false, err
		}
		verifierList = arr
	}else {
		arr, err := sk.db.GetVerifierListByIrr()
		if nil != err {
			return false, err
		}
		verifierList = arr
	}

	var flag bool
	for _, v := range verifierList.Arr {
		if v.NodeId == nodeId {
			flag = true
			break
		}
	}
	return flag, nil
}


func (sk *StakingPlugin) ListVerifierNodeID(blockHash common.Hash, blockNumber uint64) ([]discover.NodeID, error) {

	verifierList, err := sk.db.GetVerifierListByBlockHash(blockHash)
	if nil != err {
		return nil, err
	}

	if blockNumber < verifierList.Start || blockNumber > verifierList.End  {
		return nil, common.BizErrorf("ListVerifierNodeID failed: %s, start: %d, end: %d, currentNumer: %d",
			BlockNumberDisordered.Error(), verifierList.Start, verifierList.End, blockNumber)
	}

	queue := make([]discover.NodeID, len(verifierList.Arr))

	for i, v := range verifierList.Arr {
		queue[i] = v.NodeId
	}
	return queue, nil
}

func (sk *StakingPlugin) ListVerifierNodeIDFake(blockHash common.Hash, blockNumber uint64) ([]discover.NodeID, error) {

	nodeId := discover.NodeID{0x11}
	queue := make([]discover.NodeID, 0)
	queue = append(queue, nodeId)
	return queue, nil
}


func (sk *StakingPlugin) GetCandidateONEpoch(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.CandidateQueue, error) {

	var verifierList *staking.Validator_array
	if !isCommit {
		arr, err := sk.db.GetVerifierListByBlockHash(blockHash)
		if nil != err {
			return nil, err
		}
		verifierList = arr
	}else {
		arr, err := sk.db.GetVerifierListByIrr()
		if nil != err {
			return nil, err
		}
		verifierList = arr
	}


	if !isCommit && (blockNumber < verifierList.Start || blockNumber > verifierList.End) {
		return nil, common.BizErrorf("GetVerifierList failed: %s, start: %d, end: %d, currentNumer: %d",
			BlockNumberDisordered.Error(), verifierList.Start, verifierList.End, blockNumber)
	}

	queue := make(staking.CandidateQueue, len(verifierList.Arr))

	for i, v := range verifierList.Arr {

		var can *staking.Candidate
		if !isCommit {
			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				return nil, err
			}
			can = c
		}else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				return nil, err
			}
			can = c
		}
		queue[i] = can
	}

	return queue, nil
}






// flag:NOTE
// 1: Query previous round consensus validator
// 2:  Query current round consensus validaor
// 3:  Query next round consensus validator
func (sk *StakingPlugin) GetValidatorList(blockHash common.Hash, blockNumber uint64, flag uint, isCommit bool) (
	staking.ValidatorExQueue, error) {

	var validatorArr *staking.Validator_array

	switch flag {
	case PriviosRound:
		if !isCommit {
			arr, err := sk.db.GetPreValidatorListByBlockHash(blockHash)
			if nil != err {
				return nil, err
			}

			if blockNumber < arr.Start || blockNumber > arr.End {
				return nil, common.BizErrorf("Get Previous ValidatorList failed: %s, start: %d, end: %d, currentNumer: %d",
					BlockNumberDisordered.Error(), arr.Start, arr.End, blockNumber)
			}
			validatorArr = arr

		}else {
			arr, err := sk.db.GetPreValidatorListByIrr()
			if nil != err {
				return nil, err
			}

			validatorArr = arr
		}
	case CurrentRound:
		if !isCommit {
			arr, err := sk.db.GetCurrentValidatorListByBlockHash(blockHash)
			if nil != err {
				return nil, err
			}

			if blockNumber < arr.Start || blockNumber > arr.End {
				return nil, common.BizErrorf("Get Current ValidatorList failed: %s, start: %d, end: %d, currentNumer: %d",
					BlockNumberDisordered.Error(), arr.Start, arr.End, blockNumber)
			}
			validatorArr = arr
		}else {
			arr, err := sk.db.GetCurrentValidatorListByIrr()
			if nil != err {
				return nil, err
			}

			validatorArr = arr
		}
	case NextRound:
		if !isCommit {
			arr, err := sk.db.GetNextValidatorListByBlockHash(blockHash)
			if nil != err {
				return nil, err
			}

			if blockNumber < arr.Start || blockNumber > arr.End {
				return nil, common.BizErrorf("Get Next ValidatorList failed: %s, start: %d, end: %d, currentNumer: %d",
					BlockNumberDisordered.Error(), arr.Start, arr.End, blockNumber)
			}
			validatorArr = arr
		}else {
			arr, err := sk.db.GetNextValidatorListByIrr()
			if nil != err {
				return nil, err
			}
			validatorArr = arr
		}
	default:
		log.Error("Failed to call GetValidatorList", "err", ParamsErr, "flag", flag)

		return nil, common.NewBizError(ParamsErr.Error() + ", flag:=" + fmt.Sprint(flag))
	}

	queue := make(staking.ValidatorExQueue, len(validatorArr.Arr))

	for i, v := range validatorArr.Arr {

		var can *staking.Candidate

		if !isCommit {
			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				return nil, err
			}
			can = c
		}else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				return nil, err
			}
			can = c
		}



		shares, _ := new(big.Int).SetString(v.StakingWeight[1], 10)

		valEx := &staking.ValidatorEx{
			NodeId: can.NodeId,
			StakingAddress: can.StakingAddress,
			BenifitAddress: can.BenifitAddress,
			StakingTxIndex: can.StakingTxIndex,
			ProcessVersion: can.ProcessVersion,
			StakingBlockNum: can.StakingBlockNum,
			Shares: shares,
			Description: can.Description,
			ValidatorTerm: v.ValidatorTerm,
		}
		queue[i] = valEx
	}
	return queue, nil
}


func (sk *StakingPlugin) GetCandidateONRound (blockHash common.Hash, blockNumber uint64, flag uint, isCommit bool) (staking.CandidateQueue, error) {
	var validatorArr *staking.Validator_array

	switch flag {
	case PriviosRound:
		if !isCommit {
			arr, err := sk.db.GetPreValidatorListByBlockHash(blockHash)
			if nil != err {
				return nil, err
			}

			if blockNumber < arr.Start || blockNumber > arr.End {
				return nil, common.BizErrorf("Get Previous ValidatorList on GetCandidateONRound failed: %s, start: %d, end: %d, currentNumer: %d",
					BlockNumberDisordered.Error(), arr.Start, arr.End, blockNumber)
			}
			validatorArr = arr

		}else {
			arr, err := sk.db.GetPreValidatorListByIrr()
			if nil != err {
				return nil, err
			}

			validatorArr = arr
		}
	case CurrentRound:
		if !isCommit {
			arr, err := sk.db.GetCurrentValidatorListByBlockHash(blockHash)
			if nil != err {
				return nil, err
			}

			if blockNumber < arr.Start || blockNumber > arr.End {
				return nil, common.BizErrorf("Get Current ValidatorList on GetCandidateONRound failed: %s, start: %d, end: %d, currentNumer: %d",
					BlockNumberDisordered.Error(), arr.Start, arr.End, blockNumber)
			}
			validatorArr = arr
		}else {
			arr, err := sk.db.GetCurrentValidatorListByIrr()
			if nil != err {
				return nil, err
			}

			validatorArr = arr
		}
	case NextRound:
		if !isCommit {
			arr, err := sk.db.GetNextValidatorListByBlockHash(blockHash)
			if nil != err {
				return nil, err
			}

			if blockNumber < arr.Start || blockNumber > arr.End {
				return nil, common.BizErrorf("Get Next ValidatorList on GetCandidateONRound failed: %s, start: %d, end: %d, currentNumer: %d",
					BlockNumberDisordered.Error(), arr.Start, arr.End, blockNumber)
			}
			validatorArr = arr
		}else {
			arr, err := sk.db.GetNextValidatorListByIrr()
			if nil != err {
				return nil, err
			}
			validatorArr = arr
		}
	default:
		log.Error("Failed to call GetCandidateONRound", "err", ParamsErr, "flag", flag)

		return nil, common.NewBizError(ParamsErr.Error() + ", flag:=" + fmt.Sprint(flag))
	}

	queue := make(staking.CandidateQueue, len(validatorArr.Arr))

	for i, v := range validatorArr.Arr {

		var can *staking.Candidate

		if !isCommit {
			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				return nil, err
			}
			can = c
		}else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				return nil, err
			}
			can = c
		}
		queue[i] = can
	}
	return queue, nil
}



func (sk *StakingPlugin) ListCurrentValidatorID(blockHash common.Hash, blockNumber uint64) ([]discover.NodeID, error) {

	arr, err := sk.db.GetCurrentValidatorListByBlockHash(blockHash)
	if nil != err {
		return nil, err
	}

	if blockNumber < arr.Start || blockNumber > arr.End {
		return nil, common.BizErrorf("Get Current ValidatorList failed: %s, start: %d, end: %d, currentNumer: %d",
			BlockNumberDisordered.Error(), arr.Start, arr.End, blockNumber)
	}

	queue := make([]discover.NodeID, len(arr.Arr))

	for i, candidate := range arr.Arr {
		queue[i] = candidate.NodeId
	}
	return queue, err
}

func (sk *StakingPlugin) IsCurrValidator(blockHash common.Hash, nodeId discover.NodeID, isCommit bool) (bool, error) {

	var validatorArr *staking.Validator_array

	if !isCommit {
		arr, err := sk.db.GetCurrentValidatorListByBlockHash(blockHash)
		if nil != err {
			return false, err
		}
		validatorArr = arr
	}else {
		arr, err := sk.db.GetCurrentValidatorListByIrr()
		if nil != err {
			return false, err
		}
		validatorArr = arr
	}

	var flag bool
	for _, v := range validatorArr.Arr {
		if v.NodeId == nodeId {
			flag = true
			break
		}
	}
	return flag, nil
}


func (sk *StakingPlugin) GetCandidateList(blockHash common.Hash) (staking.CandidateQueue, error) {


	iter := sk.db.IteratorCandidatePowerByBlockHash(blockHash, 0)

	queue := make(staking.CandidateQueue, 0)


	for iter.Valid(); iter.Next(); {
		addrSuffix := iter.Value()
		can, err := sk.db.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			return nil, err
		}
		queue = append(queue, can)
	}


	/*// TODO MOCK
	nodeIdArr := []string{
		"0x1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee28422334",
		"0x2f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee28435466",
		"0x3f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee28544878",
		"0x3f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee28564646",
	}

	addrArr := []string{
		"0x740ce31b3fac20dac379db243021a51e80qeqqee",
		"0x740ce31b3fac20dac379db243021a51e80444555",
		"0x740ce31b3fac20dac379db243021a51e80wrwwwd",
		"0x740ce31b3fac20dac379db243021a51e80vvbbbb",
	}

	queue = make(staking.CandidateQueue, 0)
	for i:= 0; i < 4; i++ {
		can := &staking.Candidate{
			NodeId:             discover.MustHexID(nodeIdArr[i]),
			StakingAddress:     common.HexToAddress(addrArr[i]),
			BenifitAddress:     vm.StakingContractAddr,
			StakingTxIndex:     uint32(i),
			ProcessVersion:     uint32(i*i),
			Status:             staking.LowRatio,
			StakingEpoch:       uint32(1),
			StakingBlockNum:    uint64(i+2),
			Shares:             common.Big256,
			Released:           common.Big2,
			ReleasedHes:        common.Big32,
			RestrictingPlan:    common.Big1,
			RestrictingPlanHes: common.Big257,

			Description: staking.Description{
				ExternalId: "xxccccdddddddd",
				NodeName: "I Am " +fmt.Sprint(i),
				Website: "www.baidu.com",
				Details: "this is  baidu ~~",
			},
		}
		queue = append(queue, can)
	}*/
	return queue, nil
}

func (sk *StakingPlugin) GetCandidateListFake(blockHash common.Hash, isCommit bool) (staking.CandidateQueue, error) {
	return staking.CandidateQueue{}, nil
}

func (sk *StakingPlugin) IsCandidate(blockHash common.Hash, nodeId discover.NodeID, isCommit bool) (bool, error) {

	var can *staking.Candidate

	addr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		return false, err
	}

	if !isCommit {
		c, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			return false, err
		}
		can = c
	}else {
		c, err := sk.db.GetCandidateStoreByIrr(addr)
		if nil != err {
			return false, err
		}
		can = c
	}

	if nil == can || staking.Is_Invalid(can.Status) {
		return false, nil
	}
	return true, nil
}

func (sk *StakingPlugin) GetRelatedListByDelAddr (blockHash common.Hash, addr common.Address) (staking.DelRelatedQueue, error) {

	//var iter iterator.Iterator

	iter := sk.db.IteratorDelegateByBlockHashWithAddr(blockHash, addr, 0)
	queue := make(staking.DelRelatedQueue, 0)

	for iter.Valid(); iter.Next(); {
		key := iter.Key()

		prefixLen := len(staking.DelegateKeyPrefix)

		nodeIdLen := discover.NodeIDBits / 8

		// delAddr
		delAddrByte := key[prefixLen: prefixLen+common.AddressLength]
		delAddr := common.BytesToAddress(delAddrByte)

		// nodeId
		nodeIdByte := key[prefixLen+common.AddressLength: prefixLen+common.AddressLength+nodeIdLen]
		nodeId := discover.MustBytesID(nodeIdByte)

		// stakenum
		stakeNumByte := key[prefixLen+common.AddressLength+nodeIdLen:]

		num := common.BytesToUint64(stakeNumByte)

		// related
		related := &staking.DelegateRelated{
			Addr: 				delAddr,
			NodeId: 			nodeId,
			StakingBlockNum: 	num,
		}
		queue = append(queue, related)
	}
	return queue, nil
}


func (sk *StakingPlugin) Election(blockHash common.Hash, blockNumber uint64) error {


	validators, err := sk.db.GetVerifierListByIrr()
	if nil != err {
		return err
	}

	curr, err := sk.db.GetCurrentValidatorListByIrr()
	if nil != err {
		log.Error("Failed to Election: No found the current round validators", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex())
		return ValidatorNotExist
	}

	if blockNumber != (curr.End - xcom.ElectionDistance) {
		log.Error("Failed to Election: this blockNumber invalid", "Target blockNumber",
			curr.End - xcom.ElectionDistance, "Current blockNumber", blockNumber)
		return common.BizErrorf("The BlockNumber invalid, Target blockNumber: %d, Current blockNumber: %d",
			curr.End - xcom.ElectionDistance, blockNumber)
	}

	// caculate the next round start and end
	start := curr.End + 1
	end := curr.End + xcom.ConsensusSize

	proremoteCurr2NextFunc := func(start, end uint64, validators staking.ValidatorQueue) error {

		// Increase term of validator
		for i, v := range validators {
			v.ValidatorTerm++
			validators[i] = v
		}

		next := &staking.Validator_array{
			Start: start,
			End: end,
			Arr: validators,
		}

		if err := sk.db.SetNextValidatorList(blockHash, next); nil != err {
			return err
		}
		return nil
	}

	// Never match, maybe
	if nil == validators || len(validators.Arr) == 0  {
		arr := make(staking.ValidatorQueue, len(curr.Arr))
		copy(arr, curr.Arr)
		return proremoteCurr2NextFunc(start, end, arr)
	}

	currMap := make(map[discover.NodeID]struct{}, len(curr.Arr))
	for _, v := range curr.Arr {
		currMap[v.NodeId] = struct{}{}
	}

	// Exclude the current consensus round validators from the validators of the Epoch
	tmpQueue := make(staking.ValidatorQueue, 0)
	for _, v := range validators.Arr {
		if _, ok := currMap[v.NodeId]; ok {
			continue
		}
		tmpQueue = append(tmpQueue, v)
	}


	var shiftQueue staking.ValidatorQueue

	switch {
	case len(tmpQueue) == 0:
		arr := make(staking.ValidatorQueue, len(curr.Arr))
		copy(arr, curr.Arr)
		return proremoteCurr2NextFunc(start, end, arr)
	case len(tmpQueue) > 0 &&  len(tmpQueue) <= int(xcom.ShiftValidatorNum):
		shiftQueue = tmpQueue
	default:
		// elect ShiftValidatorNum (default is 8) validators by vrf
		// TODO vrf


	}

	slashCans := make(staking.SlashCandidate, 0)
	for _, v := range curr.Arr {

		addr, _ := xutil.NodeId2Addr(v.NodeId)
		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			return err
		}

		if staking.Is_LowRatio(can.Status) || staking.Is_DoubleSign(can.Status) {
			addr, _ := xutil.NodeId2Addr(v.NodeId)
			slashCans[addr] = can
		}
	}

	// Sort before removal
	curr.Arr.ValidatorSort(slashCans, staking.CompareForDel)

	// Replace the validators that can be replaced
	tmp := curr.Arr[len(shiftQueue):]
	nextValidators := make(staking.ValidatorQueue, len(tmp))
	copy(nextValidators, tmp)

	// Increase term of validator
	for i,v := range nextValidators {
		v.ValidatorTerm++
		nextValidators[i] = v
	}

	nextValidators = append(nextValidators, shiftQueue...)

	// Sort before storage
	nextValidators.ValidatorSort(slashCans, staking.CompareForStore)

	next := &staking.Validator_array{
		Start: start,
		End: end,
		Arr: nextValidators,
	}

	if err := sk.db.SetNextValidatorList(blockHash, next); nil != err {
		return err
	}

	// update candidate status
	for addr, can := range slashCans {
		if staking.Is_Valid(can.Status) && staking.Is_LowRatio(can.Status) {
			// clean the low package ratio status
			can.Status &^= staking.LowRatio
			if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
				return err
			}
		}
	}
	return nil
}

func (sk *StakingPlugin) Switch(blockHash common.Hash, blockNumber uint64) error {

	current, err := sk.db.GetCurrentValidatorListByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to Switch: Query Current round validator arr is failed",
			"blockNumber", blockNumber, "blockHash", blockHash)
		return err
	}

	if blockNumber != current.End {
		log.Error("Failed to Switch: this blockNumber invalid", "Current Round End blockNumber",
			current.End, "Current blockNumber", blockNumber)
		return common.BizErrorf("The BlockNumber invalid, Current Round End blockNumber: " +
			"%d, Current blockNumber: %d", current.End, blockNumber)
	}

	next, err := sk.db.GetNextValidatorListByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to Switch: Query Next round validator arr is failed", "blockNumber",
			blockNumber, "blockHash", blockHash)
		return err
	}

	if len(next.Arr) == 0 {
		panic(common.BizErrorf("Failed to Switch: next round validators is empty~"))
	}

	if err := sk.db.SetPreValidatorList(blockHash, current); nil != err {
		log.Error("Failed to Switch: Set Current become to Previous failed", "err", err)
		return err
	}

	if err := sk.db.SetCurrentValidatorList(blockHash, next); nil != err {
		log.Error("Failed to Switch: Set Next become to Current failed", "err", err)
		return err
	}

	if err := sk.db.DelNextValidatorListByBlockHash(blockHash); nil != err {
		return err
	}

	return nil
}

func (sk *StakingPlugin) SlashCandidates(state xcom.StateDB, blockHash common.Hash, blockNumber uint64,
	nodeId discover.NodeID, amount *big.Int, needDelete bool, slashType int)  error {

	addr, _ := xutil.NodeId2Addr(nodeId)
	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if nil != err {
		return err
	}

	if nil != can {

		log.Error("Call SlashCandidates: the can is empty", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return CandidateNotExist
	}

	epoch := xutil.CalculateEpoch(blockNumber)

	lazyCalcStakeAmount(epoch, can)



	aboutRelease := new(big.Int).Add(can.Released, can.ReleasedHes)
	aboutRestrictingPlan := new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
	total := new(big.Int).Add(aboutRelease, aboutRestrictingPlan)

	if total.Cmp(amount) < 0 {
		log.Error("Failed to SlashCandidates: the candidate total staking amount is not enough",
			"candidate total amount", total, "slashing amount" , amount)
		return common.BizErrorf("Failed to SlashCandidates: the candidate total staking amount is not enough"+
			", candidate total amount:%s, slashing amount: %s", total , amount)
	}


	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		return err
	}


	remain := amount

	slashFunc := func (remian, balance *big.Int, isNotify bool) (*big.Int, *big.Int, error) {
		if remain.Cmp(balance) >= 0 {
			state.SubBalance(vm.StakingContractAddr, balance)
			state.AddBalance(vm.RewardManagerPoolAddr, balance)

			if isNotify {
				err := RestrictingPtr.SlashingNotify(can.StakingAddress, balance, state)
				if nil != err {
					log.Error("Failed to SlashCandidates: call restrictingPlugin SlashingNotify() failed", "amount",
						balance, "err", err)
					return remian, balance, err
				}
			}

			balance = common.Big0; remain = new(big.Int).Sub(remain, balance)
		}else {
			state.SubBalance(vm.StakingContractAddr, remain)
			state.AddBalance(vm.RewardManagerPoolAddr, remain)

			if isNotify {
				err := RestrictingPtr.SlashingNotify(can.StakingAddress, remain, state)
				if nil != err {
					log.Error("Failed to SlashCandidates: call restrictingPlugin SlashingNotify() failed", "amount",
						remain, "err", err)
					return remian, balance, err
				}
			}
			balance = new(big.Int).Sub(balance, remain); remain = common.Big0
		}
		return remian, balance, nil
	}

	if can.ReleasedHes.Cmp(common.Big0) > 0 {

		val, rval, err := slashFunc(remain, can.ReleasedHes, false)
		if nil != err {
			return err
		}
		remain, can.ReleasedHes = val, rval
	}

	if remain.Cmp(common.Big0) > 0 && can.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		val, rval, err := slashFunc(remain, can.RestrictingPlanHes, true)
		if nil != err {
			return err
		}
		remain, can.RestrictingPlanHes = val, rval
	}

	if remain.Cmp(common.Big0) > 0 && can.Released.Cmp(common.Big0) > 0 {
		val, rval, err := slashFunc(remain, can.Released, false)
		if nil != err {
			return err
		}
		remain, can.Released = val, rval
	}

	if remain.Cmp(common.Big0) > 0 && can.RestrictingPlan.Cmp(common.Big0) > 0 {
		val, rval, err := slashFunc(remain, can.RestrictingPlan, true)
		if nil != err {
			return err
		}
		remain, can.RestrictingPlan = val, rval
	}

	if remain.Cmp(common.Big0) != 0 {
		log.Error("Failed to SlashCandidates: the ramain is not zero", "remain", remain)
		return common.BizErrorf("Failed to SlashCandidates: the ramain is not zero, remain:%s", remain)
	}

	remainRelease := new(big.Int).Add(can.Released, can.ReleasedHes)
	remainRestrictingPlan := new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
	canRemain := new(big.Int).Add(remainRelease, remainRestrictingPlan)

	if slashType == staking.LowRatio {
		can.Status |= staking.LowRatio
		if !xutil.CheckStakeThreshold(canRemain) {
			can.Status |= staking.NotEnough
			needDelete = true
		}
	}else if slashType == staking.DoubleSign {
		can.Status |= staking.DoubleSign
		needDelete = true
	}else {
		log.Error("Failed to SlashCandidates: the slashType is wrong", "slashType", slashType)
		return common.BizErrorf("Failed to SlashCandidates: the slashType is wrong, slashType: %d", slashType)
	}


	if !needDelete {
		sk.db.SetCanPowerStore(blockHash, addr, can)
		can.Status |= staking.Invalided
	}else {
		validators, err := sk.db.GetVerifierListByBlockHash(blockHash)
		if nil != err {
			return err
		}

		for i, val := range validators.Arr {
			if val.NodeId == nodeId {
				validators.Arr = append(validators.Arr[:i], validators.Arr[i+1:]...)
				break
			}
		}

		if err := sk.db.SetVerfierList(blockHash, validators); nil != err {
			return err
		}
	}

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		return err
	}

	return nil
}


func (sk *StakingPlugin) ProposalPassedNotify (blockHash common.Hash, blockNumber uint64, nodeIds []discover.NodeID,
	processVersion uint32) error {

	log.Info("Call ProposalPassedNotify to promote candidate processVersion", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(),
		"version", processVersion, "nodeId num", len(nodeIds))
	for _, nodeId := range nodeIds {


		addr, _ := xutil.NodeId2Addr(nodeId)
		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			return err
		}

		if nil != can {

			log.Error("Call ProposalPassedNotify: Proremote candidate processVersion failed, the can is empty",
				"blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "version", processVersion)
			continue
		}

		if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
			return err
		}

		can.ProcessVersion = processVersion

		if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
			return err
		}

		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			return err
		}
	}

	return nil
}


func (sk *StakingPlugin) DeclarePromoteNotify (blockHash common.Hash, blockNumber uint64, nodeId discover.NodeID,
	processVersion uint32) error {
	addr, _ := xutil.NodeId2Addr(nodeId)
	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if nil != err {
		return err
	}

	if nil != can {

		log.Error("Call DeclarePromoteNotify: Proremote candidate processVersion failed, the can is empty",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(),
			"version", processVersion)
		return nil
	}

	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		return err
	}

	can.ProcessVersion = processVersion

	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		return err
	}

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		return err
	}

	return nil
}


func (sk *StakingPlugin) GetLastNumber(blockNumber uint64) uint64 {

	pre, err := sk.db.GetPreValidatorListByIrr()
	if nil != err {
		return 0
	}

	if nil != pre && pre.Start <= blockNumber && pre.End >= blockNumber {
		return pre.End
	}

	curr, err := sk.db.GetCurrentValidatorListByIrr()
	if nil != err {
		return 0
	}

	if nil != curr && curr.Start <= blockNumber && curr.End >= blockNumber {
		return curr.End
	}

	next, err := sk.db.GetNextValidatorListByIrr()
	if nil != err {
		return 0
	}

	if nil != next && next.Start <= blockNumber && next.End >= blockNumber {
		return next.End
	}
	return 0
}


func (sk *StakingPlugin) GetValidator(blockNumber uint64) (*cbfttypes.Validators, error) {
	pre, err := sk.db.GetPreValidatorListByIrr()
	if nil != err {
		return nil, err
	}

	if nil != pre && pre.Start <= blockNumber && pre.End >= blockNumber {
		return build_CBFT_Validators(pre.Arr), nil
	}

	curr, err := sk.db.GetCurrentValidatorListByIrr()
	if nil != err {
		return nil, err
	}

	if nil != curr && curr.Start <= blockNumber && curr.End >= blockNumber {
		return build_CBFT_Validators(curr.Arr), nil
	}


	next, err := sk.db.GetNextValidatorListByIrr()
	if nil != err {
		return nil, err
	}

	if nil != next && next.Start <= blockNumber && next.End >= blockNumber {
		return build_CBFT_Validators(next.Arr), nil
	}

	return nil, common.BizErrorf("No Found Validators by blockNumber: %d", blockNumber)
}


// NOTE: Verify that it is the validator of the current Epoch
func (sk *StakingPlugin) IsCandidateNode(nodeID discover.NodeID) bool {

	val_arr, err := sk.db.GetVerifierListByIrr()
	if nil != err {
		log.Error("Failed to IsCandidateNode", "err", err)
		return false
	}
	for _, v := range val_arr.Arr {
		if v.NodeId == nodeID {
			return true
		}
	}
	return false
}

func build_CBFT_Validators (arr staking.ValidatorQueue) *cbfttypes.Validators {

	valMap := make(cbfttypes.ValidateNodeMap, len(arr))

	for i, v := range arr {

		pubKey, _ := v.NodeId.Pubkey()

		vn := &cbfttypes.ValidateNode {
			Index: i,
			Address: v.NodeAddress,
			PubKey: pubKey,
		}

		valMap[v.NodeId] = vn
	}

	res := &cbfttypes.Validators{
		Nodes: 	valMap,
	}
	return res
}

func lazyCalcStakeAmount(epoch uint64, can *staking.Candidate) {

	changeAmountEpoch := can.StakingEpoch

	sub := epoch - uint64(changeAmountEpoch)

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio {
		return
	}

	if can.ReleasedHes.Cmp(common.Big0) > 0 {
		can.Released = new(big.Int).Add(can.Released, can.ReleasedHes)
	}

	if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		can.RestrictingPlan = new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
	}
}

func lazyCalcDelegateAmount(epoch uint64, del *staking.Delegation) {

	// When the first time, there was no previous changeAmountEpoch
	if del.DelegateEpoch == 0 {
		return
	}

	changeAmountEpoch := del.DelegateEpoch

	sub := epoch - uint64(changeAmountEpoch)

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio {
		return
	}

	if del.ReleasedHes.Cmp(common.Big0) > 0 {
		del.Released = new(big.Int).Add(del.Released, del.ReleasedHes)
	}

	if del.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		del.RestrictingPlan = new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
	}

}


func CheckDelegateThreshold(delegate *big.Int) bool {
	return delegate.Cmp(xcom.DelegateThreshold) >= 0
}

type sortValidator struct {
	v			*staking.Validator
	x			int
	weight		int
	version		uint32
	blockNumber	uint64
	txIndex		uint32
}

type sortValidators []*sortValidator

func (svs sortValidators) Len() int {
	return len(svs)
}

func (svs sortValidators) Less(i, j int) bool {
	if svs[i].version == svs[j].version {
		if svs[i].x == svs[j].x {
			if svs[i].blockNumber == svs[j].blockNumber {
				if svs[i].txIndex == svs[j].txIndex {
					return false
				} else {
					return svs[i].txIndex < svs[j].txIndex
				}
			} else {
				return svs[i].blockNumber < svs[j].blockNumber
			}
		} else {
			return svs[i].x > svs[j].x
		}
	} else {
		return svs[i].version > svs[j].version
	}
}

func (svs sortValidators) Swap(i, j int) {
	svs[i], svs[j] = svs[j], svs[i]
}

func (sk *StakingPlugin) ProbabilityElection(validatorList staking.ValidatorQueue, currentNonce []byte, preNonces [][]byte) (staking.ValidatorQueue, error) {
	if len(currentNonce) == 0 || len(preNonces) == 0 || len(validatorList) != len(preNonces) {
		log.Error("probabilityElection failed", "validatorListSize", len(validatorList), "currentNonceSize", len(preNonces), "preNoncesSize", len(preNonces), "EpochValidatorNum", xcom.EpochValidatorNum)
		return nil, ParamsErr
	}
	sumWeights := new(big.Int)
	svList := make(sortValidators, 0)
	for _, validator := range validatorList {
		weights, err := validator.GetShares()
		if nil != err {
			return nil, ElectionErr
		}
		sumWeights.Add(sumWeights, weights)
		version, err := validator.GetProcessVersion()
		if nil != err {
			return nil, err
		}
		blockNumber, err := validator.GetStakingBlockNumber()
		if nil != err {
			return nil, err
		}
		txIndex, err := validator.GetStakingTxIndex()
		if nil != err {
			return nil, err
		}
		sv := &sortValidator{
			v:validator,
			weight:int(weights.Uint64()),
			version:version,
			blockNumber:blockNumber,
			txIndex:txIndex,
		}
		svList = append(svList, sv)
	}
	var maxValue float64 = (1 << 256) - 1
	sumWeightsFloat, err := strconv.ParseFloat(sumWeights.Text(10), 64)
	if nil != err {
		return nil, err
	}
	p := float64(len(validatorList)) * float64(xcom.ShiftValidatorNum) / sumWeightsFloat
	log.Info("probabilityElection Basic parameter", "validatorListSize", len(validatorList), "p", p, "sumWeights", sumWeightsFloat, "shiftValidatorNum", xcom.ShiftValidatorNum, "epochValidatorNum", xcom.EpochValidatorNum)
	for index, sv := range svList {
		resultStr := new(big.Int).Xor(new(big.Int).SetBytes(currentNonce), new(big.Int).SetBytes(preNonces[index])).Text(10)
		target, err := strconv.ParseFloat(resultStr, 64)
		if nil != err {
			return nil, err
		}
		targetP := target / maxValue
		bd :=  xcom.NewBinomialDistribution(sv.weight, p)
		x, err := bd.InverseCumulativeProbability(targetP)
		if nil != err {
			return nil, err
		}
		sv.x = x
		log.Debug("calculated probability", "nodeId", hex.EncodeToString(sv.v.NodeId.Bytes()), "addr", hex.EncodeToString(sv.v.NodeAddress.Bytes()), "index", index, "currentNonce", hex.EncodeToString(currentNonce), "preNonce", hex.EncodeToString(preNonces[index]), "target", target, "targetP", targetP, "weight", sv.weight, "x", x, "version", sv.version, "blockNumber", sv.blockNumber, "txIndex", sv.txIndex)
	}
	sort.Sort(svList)
	resultValidatorList := make(staking.ValidatorQueue, 0)
	for index, sv := range svList {
		if index == int(xcom.ShiftValidatorNum) {
			break
		}
		resultValidatorList = append(resultValidatorList, sv.v)
		log.Debug("sort validator", "addr", hex.EncodeToString(sv.v.NodeAddress.Bytes()), "index", index, "weight", sv.weight, "x", sv.x, "version", sv.version, "blockNumber", sv.blockNumber, "txIndex", sv.txIndex)
	}
	return resultValidatorList, nil
}
