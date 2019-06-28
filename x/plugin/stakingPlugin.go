package plugin

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
	"sync"
)

type StakingPlugin struct {
	db   *staking.StakingDB
	once sync.Once
}

var stk *StakingPlugin

var (
	AccountVonNotEnough        = errors.New("The von of account is not enough")
	DelegateVonNotEnough       = errors.New("The von of delegate is not enough")
	WithdrewDelegateVonCalcErr = errors.New("withdrew delegate von calculate err")
	ParamsErr                  = errors.New("the fn params err")
	ProcessVersionErr          = errors.New("The version of the relates node's process is too low")
	BlockNumberDisordered 	   = errors.New("The blockNumber is disordered")

	VonAmountNotRight		   = errors.New("The amount of von is not right")
)

const (
	FreeOrigin     = 0
	LockRepoOrigin = 1

	PriviosRound = uint(0)
	CurrentRound = uint(1)
	NextRound = uint(2)

	QueryStartIrr = true
	QueryStartNotIrr = false

)

// Instance a global StakingPlugin
func StakingInstance(db snapshotdb.DB) *StakingPlugin {
	if nil == stk && nil != db {
		stk = &StakingPlugin{
			db: staking.NewStakingDB(db),
		}
	}
	return stk
}

func (sk *StakingPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error) {

	return true, nil
}

func (sk *StakingPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error) {

	/*epoch := xutil.CalculateEpoch(header.Number.Uint64())

	if xutil.IsSettlementPeriod(header.Number.Uint64()) {
		success, err := sk.HandleUnCandidateReq(state, blockHash, epoch)
		if nil != err {
			log.Error("Failed to call HandleUnCandidateReq on stakingPlugin EndBlock", "blockHash",
	blockHash.Hex(), "blockNumber", header.Number.Uint64(), "err", err)
			return success, err
		}
	}

	if xutil.IsElection(header.Number.Uint64()) {
		success, err := sk.Election(blockHash, header.Number.Uint64())
		if nil != err {
			log.Error("Failed to call Election on stakingPlugin EndBlock", "blockHash", blockHash.Hex(),
	"blockNumber", header.Number.Uint64(), "err", err)
			return success, err
		}
	}

	if xutil.IsSwitch(header.Number.Uint64()) {
		success, err := sk.Switch(blockHash, header.Number.Uint64())
		if nil != err {
			log.Error("Failed to call Switch on stakingPlugin EndBlock", "blockHash", blockHash.Hex(),
	"blockNumber", header.Number.Uint64(), "err", err)
			return success, err
		}
	}
	*/
	return true, nil
}

func (sk *StakingPlugin) Confirmed(block *types.Block) error {

	return nil
}

func (sk *StakingPlugin) GetCandidateInfo (blockHash common.Hash, addr common.Address) (*xcom.Candidate, error) {

	return sk.db.GetCandidateStore(blockHash, addr)
}

func (sk *StakingPlugin) GetCandidateInfoByIrr (addr common.Address) (*xcom.Candidate, error) {
	return sk.db.GetCandidateStoreByIrr(addr)
}

func (sk *StakingPlugin) CreateCandidate (state xcom.StateDB, blockHash common.Hash, blockNumber,
	amount *big.Int, processVersion uint32, typ uint16, addr common.Address, can *xcom.Candidate) (bool, error) {

	// TODO Call gov Plugin

	if processVersion < 10101010 {
		return true, ProcessVersionErr
	} else if processVersion > 100000 {

		// TODO Call gov dclare ?
	} else {
		can.ProcessVersion = processVersion
	}

	// from account free von
	if typ == FreeOrigin {

		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to CreateCandidate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(),
				"blockHash", blockHash.Hex(), "originVon", origin, "stakingVon", amount)
			return true, AccountVonNotEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)
		can.ReleasedTmp = amount

	} else if typ == LockRepoOrigin { //  from account lockRepo von

		flag, err := RestrictingPtr.PledgeLockFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to CreateCandidate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"err", err)
			return flag, err
		}
		can.LockRepoTmp = amount
	}

	can.StakingEpoch = xutil.CalculateEpoch(blockNumber.Uint64())

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Put Can info 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}

	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Put Can power 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}
	return true, nil
}

func (sk *StakingPlugin) EditorCandidate (blockHash common.Hash, blockNumber *big.Int, can *xcom.Candidate) (bool, error) {
	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to EditorCandidate on stakingPlugin: Put Can info 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}

	return true, nil
}

func (sk *StakingPlugin) IncreaseStaking (state xcom.StateDB, blockHash common.Hash, blockNumber,
	amount *big.Int, typ uint16, can *xcom.Candidate) (bool, error) {

	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	if typ == FreeOrigin {
		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to EditorCandidate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "account", can.StakingAddress.Hex(),
				"originVon", origin, "stakingVon", can.ReleasedTmp)
			return true, AccountVonNotEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)

		can.ReleasedTmp = new(big.Int).Add(can.ReleasedTmp, amount)
	} else {

		flag, err := RestrictingPtr.PledgeLockFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to EditorCandidate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"err", err)
			return flag, err
		}

		can.LockRepoTmp = new(big.Int).Add(can.LockRepoTmp, amount)
	}

	can.StakingEpoch = epoch

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to EditorCandidate on stakingPlugin: Del Can old power failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}

	can.Shares = new(big.Int).Add(can.Shares, amount)

	// set new power of can
	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to EditorCandidate on stakingPlugin: Put Can power 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to EditorCandidate on stakingPlugin: Put Can info 2 db failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}

	return true, nil
}

func (sk *StakingPlugin) WithdrewCandidate(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	can *xcom.Candidate) (bool, error) {
	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to WithdrewCandidate on stakingPlugin: Del Can old power failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}

	if success, err := sk.withdrewStakeAmount(state, blockHash, blockNumber.Uint64(), epoch, addr, can); nil != err {
		return success, err
	}

	can.StakingEpoch = epoch

	if can.Released.Cmp(common.Big0) > 0 || can.LockRepo.Cmp(common.Big0) > 0 {

		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Failed to WithdrewCandidate on stakingPlugin: Put Can info 2 db failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return false, err
		}
	} else {
		if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
			log.Error("Failed to WithdrewCandidate on stakingPlugin: Del Can info failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return false, err
		}
	}
	return true, nil
}

func (sk *StakingPlugin) withdrewStakeAmount(state xcom.StateDB, blockHash common.Hash, blockNumber, epoch uint64,
	addr common.Address, can *xcom.Candidate) (bool, error) {

	// Direct return of money during the hesitation period
	// Return according to the way of coming
	if can.ReleasedTmp.Cmp(common.Big0) > 0 {
		state.AddBalance(can.StakingAddress, can.ReleasedTmp)
		state.SubBalance(vm.StakingContractAddr, can.ReleasedTmp)
		can.Shares = new(big.Int).Sub(can.Shares, can.ReleasedTmp)
		can.ReleasedTmp = common.Big0
	}

	if can.LockRepoTmp.Cmp(common.Big0) > 0 {

		flag, err := RestrictingPtr.ReturnLockFunds(can.StakingAddress, can.LockRepoTmp, state)
		if nil != err {
			log.Error("Failed to WithdrewCandidate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"err", err)
			return flag, err
		}

		can.Shares = new(big.Int).Sub(can.Shares, can.LockRepoTmp)
		can.LockRepoTmp = common.Big0
	}

	if can.Released.Cmp(common.Big0) > 0 || can.LockRepo.Cmp(common.Big0) > 0 {
		if err := sk.db.AddUnStakeItemStore(blockHash, epoch, addr); nil != err {
			return false, err
		}
	}
	can.Status |= xcom.Invalided

	return true, nil
}

func (sk *StakingPlugin) HandleUnCandidateReq(state xcom.StateDB, blockHash common.Hash, epoch uint64) (bool, error) {

	releaseEpoch := epoch - xcom.UnStakeFreezeRatio

	unStakeCount, err := sk.db.GetUnStakeCountStore(blockHash, releaseEpoch)
	if nil != err {
		return false, err
	}

	if unStakeCount == 0 {
		return true, nil
	}

	filterAddr := make(map[common.Address]struct{})

	for index := 1; index <= int(unStakeCount); index++ {
		addr, err := sk.db.GetUnStakeItemStore(blockHash, releaseEpoch, uint64(index))
		if nil != err {
			return false, err
		}

		if _, ok := filterAddr[addr]; ok {
			continue
		}

		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			return false, err
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
		if flag, err := sk.handleUnStake(state, blockHash, epoch, addr, can); nil != err {
			return flag, err
		}

		if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
			return false, err
		}

		filterAddr[addr] = struct{}{}
	}

	if err := sk.db.DelUnStakeCountStore(blockHash, releaseEpoch); nil != err {
		return false, err
	}

	return true, nil
}

func (sk *StakingPlugin) handleUnStake(state xcom.StateDB, blockHash common.Hash, epoch uint64,
	addr common.Address, can *xcom.Candidate) (bool, error) {

	lazyCalcStakeAmount(epoch, can)

	// Direct return of money during the hesitation period
	// Return according to the way of coming
	if can.Released.Cmp(common.Big0) > 0 {
		state.AddBalance(can.StakingAddress, can.Released)
		state.SubBalance(vm.StakingContractAddr, can.Released)
	}

	if can.LockRepo.Cmp(common.Big0) > 0 {
		flag, err := RestrictingPtr.ReturnLockFunds(can.StakingAddress, can.LockRepo, state)
		if nil != err {
			log.Error("Failed to HandleUnCandidateReq on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"err", err)
			return flag, err
		}
	}

	// delete can info
	if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
		return false, err
	}

	return true, nil
}

func (sk *StakingPlugin) GetDelegateInfo(blockHash common.Hash, delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*xcom.Delegation, error) {

	return sk.db.GetDelegateStore(blockHash, delAddr, nodeId, stakeBlockNumber)
}

func (sk *StakingPlugin) GetDelegateInfoByIrr (delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*xcom.Delegation, error) {
	return sk.GetDelegateInfoByIrr(delAddr, nodeId, stakeBlockNumber)
}


func (sk *StakingPlugin) Delegate(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	delAddr common.Address, del *xcom.Delegation, can *xcom.Candidate, typ uint16, amount *big.Int) (bool, error) {

	pubKey, _ := can.NodeId.Pubkey()
	canAddr := crypto.PubkeyToAddress(*pubKey)

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcDelegateAmount(epoch, del)

	if typ == FreeOrigin { // from account free von

		origin := state.GetBalance(delAddr)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to Delegate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "originVon", origin,
				"stakingVon", can.ReleasedTmp)
			return true, AccountVonNotEnough
		}
		state.SubBalance(delAddr, amount)
		state.AddBalance(vm.StakingContractAddr, amount)

		del.ReleasedTmp = new(big.Int).Add(del.ReleasedTmp, amount)

	} else if typ == LockRepoOrigin { //  from account lockRepo von

		flag, err := RestrictingPtr.PledgeLockFunds(delAddr, amount, state)
		if nil != err {
			log.Error("Failed to Delegate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"err", err)
			return flag, err
		}

		del.LockRepoTmp = new(big.Int).Add(del.LockRepoTmp, amount)

	}

	del.DelegateEpoch = epoch

	// set new delegate info
	if err := sk.db.SetDelegateStore(blockHash, delAddr, can.NodeId, can.StakingBlockNum, del); nil != err {
		return false, err
	}

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		return false, err
	}

	// add the candidate power
	can.Shares = new(big.Int).Add(can.Shares, amount)

	// set new power of can
	if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
		return false, err
	}

	// update can info about Shares
	if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
		return false, err
	}
	return true, nil
}

func (sk *StakingPlugin) WithdrewDelegate(state xcom.StateDB, blockHash common.Hash, blockNumber, amount *big.Int,
	delAddr common.Address, nodeId discover.NodeID, stakingBlockNum uint64, del *xcom.Delegation) (bool, error) {

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to WithdrewDelegate on stakingPlugin", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return false, err
	}

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	can, err := sk.db.GetCandidateStore(blockHash, canAddr)
	if nil != err {
		return false, err
	}

	aboutRelease := new(big.Int).Add(del.Released, del.ReleasedTmp)
	aboutLockRepo := new(big.Int).Add(del.LockRepo, del.LockRepoTmp)
	total := new(big.Int).Add(aboutRelease, aboutLockRepo)


	lazyCalcDelegateAmount(epoch, del)



	/**
	inner Fn
	*/
	subDelegateFn := func(source, sub *big.Int) (*big.Int, *big.Int) {
		state.AddBalance(delAddr, sub)
		state.SubBalance(vm.StakingContractAddr, sub)
		return new(big.Int).Sub(source, sub), common.Big0
	}

	refundFn := func(remain, aboutRelease, aboutLockRepo *big.Int) (*big.Int, *big.Int, *big.Int, bool, error) {
		// When remain is greater than or equal to del.ReleasedTmp/del.Released
		if remain.Cmp(common.Big0) > 0 {
			if remain.Cmp(aboutRelease) >= 0 && aboutRelease.Cmp(common.Big0) > 0 {

				remain, aboutRelease = subDelegateFn(remain, aboutRelease)

			} else if remain.Cmp(aboutRelease) < 0 {
				// When remain is less than or equal to del.ReleasedTmp/del.Released
				aboutRelease, remain = subDelegateFn(aboutRelease, remain)
			}
		}

		if remain.Cmp(common.Big0) > 0 {

			// When remain is greater than or equal to del.LockRepoTmp/del.LockRepo
			if remain.Cmp(aboutLockRepo) >= 0 && aboutLockRepo.Cmp(common.Big0) > 0 {

				flag, err := RestrictingPtr.ReturnLockFunds(can.StakingAddress, aboutLockRepo, state)
				if nil != err {
					log.Error("Failed to WithdrewDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"err", err)
					return remain, aboutRelease, aboutLockRepo, flag, err
				}


				remain = new(big.Int).Sub(remain, aboutLockRepo)
				aboutLockRepo = common.Big0
			} else if remain.Cmp(aboutLockRepo) < 0 {
				// When remain is less than or equal to del.LockRepoTmp/del.LockRepo


				flag, err := RestrictingPtr.ReturnLockFunds(can.StakingAddress, remain, state)
				if nil != err {
					log.Error("Failed to WithdrewDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"err", err)
					return remain, aboutRelease, aboutLockRepo, flag, err
				}

				aboutLockRepo = new(big.Int).Sub(aboutLockRepo, remain)
				remain = common.Big0
			}
		}

		return remain, aboutRelease, aboutLockRepo, true, nil
	}

	del.DelegateEpoch = epoch

	switch {

	// When the related candidate info does not exist
	case nil == can, nil != can && stakingBlockNum < can.StakingBlockNum,
	nil != can && stakingBlockNum == can.StakingBlockNum && xcom.Is_Invalid(can.Status):

		if total.Cmp(amount) < 0 {
			return true, fmt.Errorf("withdrewDelegate err: %s, delegate von: %s, withdrew von: %s",
				DelegateVonNotEnough.Error(), total.String(), amount.String())
		}

		remain := amount

		/**
		handle delegate on HesitateRatio
		*/
		remain, rtmp, ltmp, flag, err := refundFn(remain, del.ReleasedTmp, del.LockRepoTmp)
		if nil != err {
			return flag, err
		}
		del.ReleasedTmp, del.LockRepoTmp =  rtmp, ltmp
		/**
		handle delegate on EffectiveRatio
		*/
		if remain.Cmp(common.Big0) > 0 {
			remain, rtmp, ltmp, flag, err = refundFn(remain, del.Released, del.LockRepo)
			if nil != err {
				return flag, err
			}
			del.Released, del.LockRepo = rtmp, ltmp
		}

		if remain.Cmp(common.Big0) != 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: sub delegate von calculation error",
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String())
			return true, WithdrewDelegateVonCalcErr
		}

		if total.Cmp(amount) == 0 {
			if err := sk.db.DelDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum); nil != err {
				return false, err
			}

		}else {
			sub := new(big.Int).Sub(total, del.Reduction)

			if sub.Cmp(amount) < 0 {
				tmp := new(big.Int).Sub(amount, sub)
				del.Reduction = new(big.Int).Sub(del.Reduction, tmp)
			}

			if err := sk.db.SetDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum, del); nil != err {
				return false, err
			}
		}

	// Illegal parameter
	case nil != can && stakingBlockNum > can.StakingBlockNum:
		log.Error("Failed to WithdrewDelegate on stakingPlugin: the stakeBlockNum err",
			"blockHash", blockHash.Hex(), "fn.stakeBlockNum", stakingBlockNum, "can.stakeBlockNum", can.StakingBlockNum)
		return true, ParamsErr

	// When the delegate is normally revoked
	case nil != can && stakingBlockNum == can.StakingBlockNum && xcom.Is_Valid(can.Status):

		total = new(big.Int).Sub(total, del.Reduction)

		if total.Cmp(amount) < 0 {
			return true, fmt.Errorf("withdrewDelegate err: %s, delegate von: %s, withdrew von: %s",
				DelegateVonNotEnough.Error(), total.String(), amount.String())
		}

		remain := amount

		/**
		handle delegate on HesitateRatio
		*/
		//var flag bool
		//var er error
		remain, rtmp, ltmp, flag, err := refundFn(remain, del.ReleasedTmp, del.LockRepoTmp)
		if nil != err {
			return flag, err
		}
		del.ReleasedTmp, del.LockRepoTmp = rtmp, ltmp
		/**
		handle delegate on EffectiveRatio
		*/
		if remain.Cmp(common.Big0) > 0 {

			// add a UnDelegateItem
			sk.db.AddUnDelegateItemStore(blockHash, delAddr, nodeId, epoch, stakingBlockNum, remain)
			del.Reduction = new(big.Int).Add(del.Reduction, remain)
		}

		if err := sk.db.SetDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum, del); nil != err {
			return false, err
		}
	}

	// delete old can power
	if nil != can && stakingBlockNum == can.StakingBlockNum && xcom.Is_Valid(can.Status) {
		if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
			return false, err
		}

		can.Shares = new(big.Int).Sub(can.Shares, amount)

		if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
			return false, err
		}

		if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
			return false, err
		}
	}

	return true, nil
}

func (sk *StakingPlugin) HandleUnDelegateReq(state xcom.StateDB, blockHash common.Hash, epoch uint64) (bool, error) {
	releaseEpoch := epoch - xcom.ActiveUnDelegateFreezeRatio

	unDelegateCount, err := sk.db.GetUnDelegateCountStore(blockHash, releaseEpoch)
	if nil != err {
		return false, err
	}

	if unDelegateCount == 0 {
		return true, nil
	}

	//filterAddr := make(map[string]struct{})

	for index := 1; index <= int(unDelegateCount); index++ {
		unDelegateItem, err := sk.db.GetUnDelegateItemStore(blockHash, releaseEpoch, uint64(index))
		if nil != err {
			return false, err
		}

		//if _, ok := filterAddr[fmt.Sprint(unDelegateItem.KeySuffix)]; ok {
		//	continue
		//}

		del, err := sk.db.GetDelegateStoreBySuffix(blockHash, unDelegateItem.KeySuffix)
		if nil != err {
			return false, err
		}

		if nil == del {
			// todo This maybe be nil
			continue
		}

		if flag, err := sk.handleUnDelegate(state, blockHash, epoch, unDelegateItem, del); nil != err {
			return flag, err
		}

		//filterAddr[fmt.Sprint(unDelegateItem.KeySuffix)] = struct{}{}
	}

	return true, nil
}

func (sk *StakingPlugin) handleUnDelegate(state xcom.StateDB, blockHash common.Hash, epoch uint64,
	unDel *xcom.UnDelegateItem, del *xcom.Delegation) (bool, error) {

	// del addr
	delAddrByte := unDel.KeySuffix[0:common.AddressLength]
	delAddr := common.BytesToAddress(delAddrByte)

	nodeIdLen := discover.NodeIDBits / 8

	canAddrByte := unDel.KeySuffix[common.AddressLength : common.AddressLength + nodeIdLen]
	canAddr := common.BytesToAddress(canAddrByte)

	//
	stakeBlockNum := unDel.KeySuffix[common.AddressLength + nodeIdLen:]
	num_int, _ := strconv.Atoi(string(stakeBlockNum))
	num := uint64(num_int)


	lazyCalcDelegateAmount(epoch, del)


	amount := unDel.Amount


	if amount.Cmp(del.Reduction) >= 0 { // full withdrawal
		state.SubBalance(vm.StakingContractAddr, del.Released)
		state.AddBalance(delAddr, del.Released)

		// todo call Restricting for flush lockRepo

		if err := sk.db.DelDelegateStoreBySuffix(blockHash, unDel.KeySuffix); nil != err {
			return false, err
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

			if remain.Cmp(del.LockRepo) >= 0 {
				// todo call Restricting for flush lockRepo

				del.LockRepo = common.Big0; remain = new(big.Int).Sub(remain, del.LockRepo)
			}else {
				// todo call Restricting for flush remain

				del.LockRepo = new(big.Int).Sub(del.LockRepo, remain); remain = common.Big0
			}
		}

		if remain.Cmp(common.Big0) > 0 {
			log.Error("Failed to call handleUnDelegate", "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "canAddr", canAddr.Hex(), "stakeBlockNumber", num)
			return false, VonAmountNotRight
		}

		del.Reduction = new(big.Int).Sub(del.Reduction, amount)

		del.DelegateEpoch = epoch

		if err := sk.db.SetDelegateStoreBySuffix(blockHash, unDel.KeySuffix, del); nil != err {
			return false, err
		}
	}

	return true, nil
}




func (sk *StakingPlugin) ElectNextVerifierList(blockHash common.Hash, blockNumber uint64) (bool, error) {

	log.Info("Call ElectNextVerifierList", "blockNumber", blockNumber, "blockHash", blockHash.Hex())

	old_verifierArr, err := sk.db.GetVerifierListByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to ElectNextVerifierList", "blockNumber", blockNumber, "blockHash",
			blockHash.Hex(), "err", err)
		return false, err
	}

	if old_verifierArr.End != blockNumber {
		log.Error("Failed to ElectNextVerifierList: this blockNumber invalid", "Old Epoch End blockNumber",
			old_verifierArr.End, "Current blockNumber", blockNumber)
		return true, fmt.Errorf("The BlockNumber invalid, Old Epoch End blockNumber: %d, Current blockNumber: %d",
			old_verifierArr.End, blockNumber)
	}



	iter := sk.db.IteratorCandidatePowerByBlockHash(blockHash, int(xcom.EpochValidatorNum))

	start := old_verifierArr.End + 1
	end := old_verifierArr.End + xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &xcom.Validator_array{
		Start: 	start,
		End: 	end,
	}

	queue := make(xcom.ValidatorQueue, 0)

	for count := 0; iter.Valid() && count < int(xcom.EpochValidatorNum); iter.Next() {
		addrSuffix := iter.Value()
		var can *xcom.Candidate

		can, err := sk.db.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			log.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return false, err
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [4]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &xcom.Validator{
			NodeAddress: addr,
			NodeId: 	 can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
	}

	if len(queue) == 0 {
		panic(fmt.Errorf("Failed to ElectNextVerifierList: Select zero validators~"))
		//return true, fmt.Errorf("Failed to ElectNextVerifierList: Select zero validators~")
	}

	new_verifierArr.Arr = queue

	err = sk.db.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		return false, err
	}
	return true, nil
}

func (sk *StakingPlugin) GetVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (xcom.CandidateQueue, error) {

	var verifierList *xcom.Validator_array
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
		return nil, fmt.Errorf("GetVerifierList failed: %s, start: %d, end: %d, currentNumer: %d",
			BlockNumberDisordered.Error(), verifierList.Start, verifierList.End, blockNumber)
	}

	resultArr := make(xcom.CandidateQueue, len(verifierList.Arr))

	for _, v := range verifierList.Arr {

		var can *xcom.Candidate
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

		resultArr = append(resultArr, can)
	}

	return resultArr, nil
}


func (sk *StakingPlugin) IsCurrVerifier(blockHash common.Hash, nodeId discover.NodeID, isCommit bool) (bool, error) {

	var verifierList *xcom.Validator_array

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
	var verifierNodeIDList []discover.NodeID

	candidateQueue, err := sk.GetVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err{
		return nil, err
	}
	for _, candidate := range candidateQueue{
		verifierNodeIDList = append(verifierNodeIDList, candidate.NodeId)
	}
	return verifierNodeIDList, nil
}

// flag:NOTE
// 1: Query previous round consensus validator
// 2:  Query current round consensus validaor
// 3:  Query next round consensus validator
func (sk *StakingPlugin) GetValidatorList(blockHash common.Hash, blockNumber uint64, flag uint, isCommit bool) (
	xcom.ValidatorExQueue, error) {

	var validatorArr *xcom.Validator_array

	switch flag {
	case PriviosRound:
		if !isCommit {
			arr, err := sk.db.GetPreValidatorListByBlockHash(blockHash)
			if nil != err {
				return nil, err
			}

			if blockNumber < arr.Start || blockNumber > arr.End {
				return nil, fmt.Errorf("Get Previous ValidatorList failed: %s, start: %d, end: %d, currentNumer: %d",
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
				return nil, fmt.Errorf("Get Current ValidatorList failed: %s, start: %d, end: %d, currentNumer: %d",
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
				return nil, fmt.Errorf("Get Next ValidatorList failed: %s, start: %d, end: %d, currentNumer: %d",
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

		return nil, fmt.Errorf(ParamsErr.Error() + ", flag:=" + fmt.Sprint(flag))
	}

	result := make(xcom.ValidatorExQueue, len(validatorArr.Arr))

	for _, v := range validatorArr.Arr {

		var can *xcom.Candidate

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

		canEx := &xcom.ValidatorEx{
			Candidate: can,
			ValidatorTerm: v.ValidatorTerm,
		}
		result = append(result, canEx)
	}
	return result, nil
}

func (sk *StakingPlugin) IsCurrValidate(blockHash common.Hash, nodeId discover.NodeID, isCommit bool) (bool, error) {

	var validatorArr *xcom.Validator_array

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


func (sk *StakingPlugin) GetCandidateList(blockHash common.Hash, isCommit bool) (xcom.CandidateQueue, error) {


	var iter iterator.Iterator

	if !isCommit {

		itr := sk.db.IteratorCandidatePowerByBlockHash(blockHash, 0)
		iter = itr
	}else {
		itr := sk.db.IteratorCandidatePowerByIrr(0)
		iter = itr
	}

	queue := make(xcom.CandidateQueue, 0)


	for iter.Valid(); iter.Next(); {
		addrSuffix := iter.Value()
		var can *xcom.Candidate

		if !isCommit {
			c, err := sk.db.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
			if nil != err {
				return nil, err
			}
			can = c
		}else {
			c, err := sk.db.GetCandidateStoreByIrrWithSuffix(addrSuffix)
			if nil != err {
				return nil, err
			}
			can = c
		}
		queue = append(queue, can)
	}
	return nil, nil
}

func (sk *StakingPlugin) IsCandidate(blockHash common.Hash, nodeId discover.NodeID, isCommit bool) (bool, error) {

	var can *xcom.Candidate

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

	if nil == can {
		return false, nil
	}
	return true, nil
}

func (sk *StakingPlugin) GetRelatedListByDelAddr (blockHash common.Hash, addr common.Address,
	isCommit bool) (xcom.DelRelatedQueue, error) {

	var iter iterator.Iterator

	if !isCommit {

		itr := sk.db.IteratorDelegateByBlockHashWithAddr(blockHash, addr, 0)
		iter = itr
	}else {
		itr := sk.db.IteratorDelegateByIrrWithAddr(addr, 0)
		iter = itr
	}

	queue := make(xcom.DelRelatedQueue, 0)

	for iter.Valid(); iter.Next(); {
		key := iter.Key()

		prefixLen := len(xcom.DelegateKeyPrefix)

		nodeIdLen := discover.NodeIDBits / 8

		// delAddr
		delAddrByte := key[prefixLen: prefixLen+common.AddressLength]
		delAddr := common.BytesToAddress(delAddrByte)

		// nodeId
		nodeIdByte := key[prefixLen+common.AddressLength: prefixLen+common.AddressLength+nodeIdLen]
		nodeId := discover.MustBytesID(nodeIdByte)

		// stakenum
		stakeNumByte := key[prefixLen+common.AddressLength+nodeIdLen:]
		numInt, err := strconv.Atoi(string(stakeNumByte))
		if nil != err {
			return nil, err
		}
		num := uint64(numInt)

		// related
		related := &xcom.DelegateRelated{
			Addr: 				delAddr,
			NodeId: 			nodeId,
			StakingBlockNum: 	num,
		}
		queue = append(queue, related)
	}
	return queue, nil
}


func (sk *StakingPlugin) Election(blockHash common.Hash, blockNumber uint64) (bool, error) {

	round := xutil.CalculateRound(blockNumber)
	_ = round

	return true, nil
}

func (sk *StakingPlugin) Switch(blockHash common.Hash, blockNumber uint64) (bool, error) {

	current, err := sk.db.GetCurrentValidatorListByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to Switch: Query Current round validator arr is failed",
			"blockNumber", blockNumber, "blockHash", blockHash)
		return false, err
	}

	if blockNumber != current.End {
		log.Error("Failed to Switch: this blockNumber invalid", "Current Round End blockNumber",
			current.End, "Current blockNumber", blockNumber)
		return true, fmt.Errorf("The BlockNumber invalid, Current Round End blockNumber: " +
			"%d, Current blockNumber: %d", current.End, blockNumber)
	}

	next, err := sk.db.GetNextValidatorListByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to Switch: Query Next round validator arr is failed", "blockNumber",
			blockNumber, "blockHash", blockHash)
		return false, err
	}

	if len(next.Arr) == 0 {
		panic(fmt.Errorf("Failed to Switch: next round validators is empty~"))
	}

	if err := sk.db.SetPreValidatorList(blockHash, current); nil != err {
		// TODO log.Error("// todo")
		return false, err
	}

	if err := sk.db.SetCurrentValidatorList(blockHash, current); nil != err {
		// TODO log.Error("// todo")
		return false, err
	}

	return true, nil
}

func (sk *StakingPlugin) GetAllPackageRatio() {

}

func (sk *StakingPlugin) GetPackageRatio() {

}

func (sk *StakingPlugin) SlashCandidates() {

}



func lazyCalcStakeAmount(epoch uint64, can *xcom.Candidate) {

	changeAmountEpoch := can.StakingEpoch

	sub := epoch - changeAmountEpoch

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio {
		return
	}

	if can.ReleasedTmp.Cmp(common.Big0) > 0 {
		can.Released = new(big.Int).Add(can.Released, can.ReleasedTmp)
	}

	if can.LockRepoTmp.Cmp(common.Big0) > 0 {
		can.LockRepo = new(big.Int).Add(can.LockRepo, can.LockRepoTmp)
	}
}

func lazyCalcDelegateAmount(epoch uint64, del *xcom.Delegation) {

	changeAmountEpoch := del.DelegateEpoch

	sub := epoch - changeAmountEpoch

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio {
		return
	}

	if del.ReleasedTmp.Cmp(common.Big0) > 0 {
		del.Released = new(big.Int).Add(del.Released, del.ReleasedTmp)
	}

	if del.LockRepoTmp.Cmp(common.Big0) > 0 {
		del.LockRepo = new(big.Int).Add(del.LockRepo, del.LockRepoTmp)
	}

	/*switch  {
	case canStatus&xcom.NotExist == xcom.NotExist:

		mergeFn()

	case xcom.Is_Invalid_Slashed(canStatus),
		xcom.Is_Invalid_NotEnough(canStatus):

		if epoch - xcom.PassiveUnDelegateFreezeRatio <= changeAmountEpoch {
			return
		}
		mergeFn()
	case xcom.Is_Valid(canStatus):

		if epoch - changeAmountEpoch < xcom.HesitateRatio || epoch - xcom.ActiveUnDelegateFreezeRatio <= changeAmountEpoch {
			return
		}

		mergeFn()

	}*/
}

/*func mergeAmount(mark uint8, target, tmp *big.Int) *big.Int {
	if mark == increase {
		return new(big.Int).Add(target, tmp)
	} else if mark == decrease {
		return new(big.Int).Sub(target, tmp)
	}
	return target
}*/

//func (sk *StakingPlugin) sumStakeAmount (can *xcom.Candidate) *big.Int {
//
//	aoubt_release := new(big.Int).Add(can.Released, can.ReleasedTmp)
//
//	about_locked := new(big.Int).Add(can.L)
//
//	return
//}

func CheckStakeThreshold(stake *big.Int) bool {
	return stake.Cmp(xcom.StakeThreshold) >= 0
}

func CheckDelegateThreshold(delegate *big.Int) bool {
	return delegate.Cmp(xcom.DelegateThreshold) >= 0
}
