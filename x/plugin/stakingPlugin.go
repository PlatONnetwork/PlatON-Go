package plugin

import (
	"crypto/ecdsa"
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
	"sync"
)

type StakingPlugin struct {
	db *StakingDB
	once  sync.Once

}



var stk *StakingPlugin



var (
	AccountVonNotEnough = errors.New("The account von is not enough")

)

const (
	FreeOrigin, increase = 0, uint8(0)
	LockRepoOrigin, decrease = 1, uint8(1)

	valided = uint8(2)
)



// Instance a global StakingPlugin
func  StakingInstance (db xcom.SnapshotDB) *StakingPlugin {
	if nil == stk && nil != db {
		stk = &StakingPlugin{
			db: NewStakingDB(db),
		}
	}
	return stk
}




func (sk *StakingPlugin) BeginBlock (header *types.Header, state xcom.StateDB) (bool, error) {

	return false, nil
}

func (sk *StakingPlugin) EndBlock(header *types.Header, state xcom.StateDB) (bool, error) {


	return false, nil
}

func (sk *StakingPlugin) Confirmed(block *types.Block) error {

	return nil
}






func (sk *StakingPlugin) GetCandidateInfo(blockHash common.Hash,  nodeId discover.NodeID) (*xcom.Candidate, error) {

	canByte, err := sk.db.getCandidate(blockHash, nodeId)
	if nil != err {
		return nil, err
	}

	var can xcom.Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func (sk *StakingPlugin) CreateCandidate(state xcom.StateDB, blockHash common.Hash,  blockNumber *big.Int, can *xcom.Candidate, typ uint16) (bool, error) {

	var pubKey ecdsa.PublicKey

	if pk, err := can.NodeId.Pubkey(); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: nodeId convert pubkey failed", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}else {
		pubKey = *pk
	}

	// from account free von
	if typ == FreeOrigin {

		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(can.ReleasedTmp) < 0 {
			log.Error("Failed to CreateCandidate on stakingPlugin: the account free von is not Enough", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "originVon", origin, "stakingVon", can.ReleasedTmp)
			return true, AccountVonNotEnough
		}
		state.SubBalance(can.StakingAddress, can.ReleasedTmp)
		state.AddBalance(vm.StakingContractAddr, can.ReleasedTmp)

	}else if typ == LockRepoOrigin {  //  from account lockRepo von

		 // TODO call lockRepoPlugin

	}

	can.StakingEpoch = xutil.CalculateEpoch(blockNumber.Uint64())

	addr := crypto.PubkeyToAddress(pubKey)

	if err := sk.db.setCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Put Can info 2 db failed", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}

	if err := sk.db.setCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Put Can power 2 db failed", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}
	return true, nil
}


func (sk *StakingPlugin) EditorCandidate (state xcom.StateDB, blockHash common.Hash,  blockNumber *big.Int, can *xcom.Candidate, typ uint16, amount *big.Int) (bool, error) {

	pubKey, _ := can.NodeId.Pubkey()

	lazyCalcStakeAmount(blockNumber.Uint64(), can)

	can.StakingEpoch = xutil.CalculateEpoch(blockNumber.Uint64())

	addr := crypto.PubkeyToAddress(*pubKey)

	if amount.Cmp(big.NewInt(0)) > 0 {
		if typ == FreeOrigin  {
			origin := state.GetBalance(can.StakingAddress)
			if origin.Cmp(amount) < 0 {
				log.Error("Failed to EditorCandidate on stakingPlugin: the account free von is not Enough", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "originVon", origin, "stakingVon", can.ReleasedTmp)
				return true, AccountVonNotEnough
			}
			state.SubBalance(can.StakingAddress, amount)
			state.AddBalance(vm.StakingContractAddr, amount)

			can.ReleasedTmp = new(big.Int).Add(can.ReleasedTmp, amount)
		}else {

			// TODO call lockRepoPlugin



			can.LockRepoTmp = new(big.Int).Add(can.LockRepoTmp, amount)
		}

		// delete old power of can
		if err := sk.db.delCanPowerStore(blockHash, can); nil != err {
			log.Error("Failed to EditorCandidate on stakingPlugin: Del Can old power failed", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return false, err
		}

		can.Shares = new(big.Int).Add(can.Shares, amount)

		if err := sk.db.setCanPowerStore(blockHash, addr, can); nil != err {
			log.Error("Failed to EditorCandidate on stakingPlugin: Put Can power 2 db failed", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return false, err
		}

	}

	if err := sk.db.setCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to EditorCandidate on stakingPlugin: Put Can info 2 db failed", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}

	return true, nil
}



func (sk *StakingPlugin) WithdrewCandidate (state xcom.StateDB, blockHash common.Hash,  blockNumber *big.Int, can *xcom.Candidate) (bool, error) {
	pubKey, _ := can.NodeId.Pubkey()

	lazyCalcStakeAmount(blockNumber.Uint64(), can)

	can.StakingEpoch = xutil.CalculateEpoch(blockNumber.Uint64())

	addr := crypto.PubkeyToAddress(*pubKey)


	// delete old power of can
	if err := sk.db.delCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to WithdrewCandidate on stakingPlugin: Del Can old power failed", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return false, err
	}

	if success, err := sk.withdrewStakeAmount(state, blockHash, blockNumber.Uint64(), addr, can); nil != err {
		return success, err
	}

	if can.Released.Cmp(common.Big0) > 0 || can.LockRepo.Cmp(common.Big0) > 0 {
		if err := sk.db.setCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Failed to WithdrewCandidate on stakingPlugin: Put Can info 2 db failed", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return false, err
		}
	}else {
		if err := sk.db.delCandidateStore(blockHash, addr); nil != err {
			log.Error("Failed to WithdrewCandidate on stakingPlugin: Del Can info failed", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return false, err
		}
	}
	return true, nil
}






func (sk *StakingPlugin) withdrewStakeAmount (state xcom.StateDB, blockHash common.Hash, blockNumber uint64, addr common.Address, can *xcom.Candidate) (bool, error) {
	curEpoch := xutil.CalculateEpoch(blockNumber)

	// Direct return of money during the hesitation period
	// Return according to the way of coming
	if can.ReleasedTmp.Cmp(common.Big0) > 0 {
		state.AddBalance(can.StakingAddress, can.ReleasedTmp)
		state.SubBalance(vm.StakingContractAddr, can.ReleasedTmp)
		can.Shares = new(big.Int).Sub(can.Shares, can.ReleasedTmp)
	}

	if can.LockRepoTmp.Cmp(common.Big0) > 0 {
		// TODO call lockRepoPlugin

		can.Shares = new(big.Int).Sub(can.Shares, can.LockRepoTmp)
	}

	if can.Released.Cmp(common.Big0) > 0 || can.LockRepo.Cmp(common.Big0) > 0 {
		if err := sk.db.addUnStakeItem(blockHash, curEpoch, addr); nil != err {
			return false, err
		}
	}

	can.Mark = valided
	can.Status|= xcom.Invalided

	return true, nil
}



















func lazyCalcStakeAmount (blockNumber uint64, can *xcom.Candidate) {

	curEpoch := xutil.CalculateEpoch(blockNumber)

	changeAmountEpoch := can.StakingEpoch

	sub := curEpoch - changeAmountEpoch

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio {
		return
	}

	if can.ReleasedTmp.Cmp(common.Big0) > 0 {
		can.Released = mergeAmount(can.Mark, can.Released, can.ReleasedTmp)
	}

	if can.LockRepoTmp.Cmp(common.Big0) > 0 {
		can.LockRepo = mergeAmount(can.Mark, can.LockRepo, can.LockRepoTmp)
	}
}



func (sk *StakingPlugin) lazyCalcDelegateAmount (del  *xcom.Delegation) {



}

func  mergeAmount (mark uint8, target, tmp *big.Int) *big.Int {
	if mark == increase {
		return  new(big.Int).Add(target, tmp)
	}else if mark == decrease {
		return  new(big.Int).Sub(target, tmp)
	}
	return target
}



//func (sk *StakingPlugin) sumStakeAmount (can *xcom.Candidate) *big.Int {
//
//	aoubt_release := new(big.Int).Add(can.Released, can.ReleasedTmp)
//
//	about_locked := new(big.Int).Add(can.L)
//
//	return
//}

func CheckStakeThreshold (stake *big.Int) bool {
	return stake.Cmp(xcom.StakeThreshold) >= 0
}


func CheckDelegateThreshold (delegate *big.Int) bool {
	return delegate.Cmp(xcom.DelegateThreshold) >= 0
}




