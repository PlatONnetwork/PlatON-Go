package plugin

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common/math"

	"github.com/PlatONnetwork/PlatON-Go/x/handler"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/vrf"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type StakingPlugin struct {
	db       *staking.StakingDB
	eventMux *event.TypeMux
}

var (
	stakePlnOnce sync.Once
	stk          *StakingPlugin
)

const (
	FreeOrigin            = uint16(0)
	RestrictingPlanOrigin = uint16(1)

	PreviousRound = uint(0)
	CurrentRound  = uint(1)
	NextRound     = uint(2)

	QueryStartIrr    = true
	QueryStartNotIrr = false

	EpochValIndexSize = 2
	RoundValIndexSize = 6
)

// Instance a global StakingPlugin
func StakingInstance() *StakingPlugin {
	stakePlnOnce.Do(func() {
		log.Info("Init Staking plugin ...")
		stk = &StakingPlugin{
			db: staking.NewStakingDB(),
		}
	})
	return stk
}

func (sk *StakingPlugin) SetEventMux(eventMux *event.TypeMux) {
	sk.eventMux = eventMux
}

func (sk *StakingPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	// Do nothings
	return nil
}

func (sk *StakingPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {

	epoch := xutil.CalculateEpoch(header.Number.Uint64())

	if xutil.IsEndOfEpoch(header.Number.Uint64()) {

		// handle UnStaking Item
		err := sk.HandleUnCandidateItem(state, header.Number.Uint64(), blockHash, epoch)
		if nil != err {
			log.Error("Failed to call HandleUnCandidateItem on stakingPlugin EndBlock",
				"blockNumber", header.Number.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		// Election next epoch validators
		if err := sk.ElectNextVerifierList(blockHash, header.Number.Uint64(), state); nil != err {
			log.Error("Failed to call ElectNextVerifierList on stakingPlugin EndBlock",
				"blockNumber", header.Number.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}

	if xutil.IsElection(header.Number.Uint64()) {

		// ELection next round validators
		err := sk.Election(blockHash, header, state)
		if nil != err {
			log.Error("Failed to call Election on stakingPlugin EndBlock",
				"blockNumber", header.Number.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}

	}
	return nil
}

func (sk *StakingPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {

	if xutil.IsElection(block.NumberU64()) {

		next, err := sk.getNextValList(block.Hash(), block.NumberU64(), QueryStartNotIrr)
		if nil != err {
			log.Error("Failed to Query Next validators on stakingPlugin Confirmed When Election block",
				"blockNumber", block.Number().Uint64(), "blockHash", block.Hash().TerminalString(), "err", err)
			return err
		}

		current, err := sk.getCurrValList(block.Hash(), block.NumberU64(), QueryStartNotIrr)
		if nil != err {
			log.Error("Failed to Query Current Round validators on stakingPlugin Confirmed When Election block",
				"blockNumber", block.Number().Uint64(), "blockHash", block.Hash().TerminalString(), "err", err)
			return err
		}

		diff := make(staking.ValidatorQueue, 0)
		var isCurr, isNext bool

		currMap := make(map[discover.NodeID]struct{})
		for _, v := range current.Arr {
			currMap[v.NodeId] = struct{}{}
			if nodeId == v.NodeId {
				isCurr = true
			}
		}

		for _, v := range next.Arr {
			if _, ok := currMap[v.NodeId]; !ok {
				diff = append(diff, v)
			}

			if nodeId == v.NodeId {
				isNext = true
			}
		}

		// This node will only initiating a pre-connection,
		// When the node is one of the next round of validators.
		if isCurr && isNext {
			sk.addConsensusNode(diff)
			log.Debug("Call addConsensusNode finished on stakingPlugin, node is curr validator AND next validator",
				"blockNumber", block.NumberU64(), "blockHash", block.Hash().TerminalString(), "diff size", len(diff))
		} else if !isCurr && isNext {
			sk.addConsensusNode(next.Arr)
			log.Debug("Call addConsensusNode finished on stakingPlugin, node is new validator",
				"blockNumber", block.NumberU64(), "blockHash", block.Hash().TerminalString(), "diff size", len(next.Arr))
		} else {
			return nil
		}
	}

	return nil
}

//func distinct(list, target staking.ValidatorQueue) staking.ValidatorQueue {
//	currentMap := make(map[discover.NodeID]bool)
//	for _, v := range target {
//		currentMap[v.NodeId] = true
//	}
//	result := make(staking.ValidatorQueue, 0)
//	for _, v := range list {
//		if _, ok := currentMap[v.NodeId]; !ok {
//			result = append(result, v)
//		}
//	}
//	return result
//}

func (sk *StakingPlugin) addConsensusNode(nodes staking.ValidatorQueue) {
	for _, node := range nodes {
		sk.eventMux.Post(cbfttypes.AddValidatorEvent{NodeID: node.NodeId})
	}
}

func (sk *StakingPlugin) GetCandidateInfo(blockHash common.Hash, addr common.Address) (*staking.Candidate, error) {
	return sk.db.GetCandidateStore(blockHash, addr)
}

func (sk *StakingPlugin) GetCandidateCompactInfo(blockHash common.Hash, blockNumber uint64, addr common.Address) (*staking.CandidateHex, error) {
	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if nil != err {
		return nil, err
	}

	epoch := xutil.CalculateEpoch(blockNumber)

	lazyCalcStakeAmount(epoch, can)
	canHex := buildCanHex(can)

	return canHex, nil
}

func (sk *StakingPlugin) GetCandidateInfoByIrr(addr common.Address) (*staking.Candidate, error) {
	return sk.db.GetCandidateStoreByIrr(addr)
}

func (sk *StakingPlugin) CreateCandidate(state xcom.StateDB, blockHash common.Hash, blockNumber,
	amount *big.Int, typ uint16, addr common.Address, can *staking.Candidate) error {

	log.Debug("Call CreateCandidate", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
		"nodeId", can.NodeId.String())

	// from account free von
	if typ == FreeOrigin {

		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to CreateCandidate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
				"stakeAddr", can.StakingAddress.Hex(), "originVon", origin, "stakingVon", amount)
			return staking.ErrAccountVonNoEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)
		can.ReleasedHes = amount

	} else if typ == RestrictingPlanOrigin { //  from account RestrictingPlan von

		err := rt.PledgeLockFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to CreateCandidate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
				"stakeAddr", can.StakingAddress.Hex(), "stakingVon", amount, "err", err)
			return err
		}
		can.RestrictingPlanHes = amount
	} else {

		log.Error("Failed to CreateCandidate on stakingPlugin", "err", staking.ErrWrongVonOptType,
			"got type", typ, "need type", fmt.Sprintf("%d or %d", FreeOrigin, RestrictingPlanOrigin))
		return staking.ErrWrongVonOptType
	}

	can.StakingEpoch = uint32(xutil.CalculateEpoch(blockNumber.Uint64()))

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Store Candidate info is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}

	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Store Candidate power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}

	// add the account staking Reference Count
	if err := sk.db.AddAccountStakeRc(blockHash, can.StakingAddress); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Store Staking Account Reference Count (add) is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "NodeID", can.NodeId.String(),
			"staking Account", can.StakingAddress.String(), "err", err)
		return err
	}

	return nil
}

/// This method may only be called when creatStaking
func (sk *StakingPlugin) RollBackStaking(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	addr common.Address, typ uint16) error {

	log.Debug("Call RollBackStaking", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeAddr", addr.String())

	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if nil != err {
		return err
	}

	if blockNumber.Uint64() != can.StakingBlockNum {

		log.Error("Failed to RollBackStaking on stakingPlugin: current blockNumber is not equal stakingBlockNumber",
			"blockNumber", blockNumber, "stakingBlockNumber", can.StakingBlockNum)
		return staking.ErrBlockNumberDisordered
	}

	// RollBack Staking
	if typ == FreeOrigin {

		state.AddBalance(can.StakingAddress, can.ReleasedHes)
		state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)

	} else if typ == RestrictingPlanOrigin {

		err := rt.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
		if nil != err {
			log.Error("Failed to RollBackStaking on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
				"stakeAddr", can.StakingAddress.Hex(), "RollBack stakingVon", can.RestrictingPlanHes, "err", err)
			return err
		}
	} else {

		log.Error("Failed to RollBackStaking on stakingPlugin", "err", staking.ErrWrongVonOptType,
			"got type", typ, "need type", fmt.Sprintf("%d or %d", FreeOrigin, RestrictingPlanOrigin))
		return staking.ErrWrongVonOptType
	}

	if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
		log.Error("Failed to RollBackStaking on stakingPlugin: Delete Candidate info is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}

	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to RollBackStaking on stakingPlugin: Delete Candidate power failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}

	// sub the account staking Reference Count
	if err := sk.db.SubAccountStakeRc(blockHash, can.StakingAddress); nil != err {
		log.Error("Failed to RollBackStaking on stakingPlugin: Store Staking Account Reference Count (sub) is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
			"staking Account", can.StakingAddress.String(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) EditCandidate(blockHash common.Hash, blockNumber *big.Int, can *staking.Candidate) error {

	log.Debug("Call EditCandidate", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
		"nodeId", can.NodeId.String())

	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to EditCandidate on stakingPlugin: Store Candidate info is failed",
			"nodeId", can.NodeId.String(), "blockNumber", blockNumber.Uint64(),
			"blockHash", blockHash.Hex(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) IncreaseStaking(state xcom.StateDB, blockHash common.Hash, blockNumber,
	amount *big.Int, typ uint16, can *staking.Candidate) error {

	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	log.Debug("Call IncreaseStaking", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"epoch", epoch, "nodeId", can.NodeId.String(), "account", can.StakingAddress.Hex(), "typ", typ, "amount", amount)

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	if typ == FreeOrigin {
		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to IncreaseStaking on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
				"nodeId", can.NodeId.String(), "account", can.StakingAddress.Hex(),
				"originVon", origin, "stakingVon", amount)
			return staking.ErrAccountVonNoEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)

		can.ReleasedHes = new(big.Int).Add(can.ReleasedHes, amount)

	} else if typ == RestrictingPlanOrigin {

		err := rt.PledgeLockFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to IncreaseStaking on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
				"nodeId", can.NodeId.String(), "account", can.StakingAddress.Hex(), "amount", amount, "err", err)
			return err
		}

		can.RestrictingPlanHes = new(big.Int).Add(can.RestrictingPlanHes, amount)
	} else {

		log.Error("Failed to IncreaseStaking on stakingPlugin", "err", staking.ErrWrongVonOptType,
			"got type", typ, "need type", fmt.Sprintf("%d or %d", FreeOrigin, RestrictingPlanOrigin))
		return staking.ErrWrongVonOptType
	}

	can.StakingEpoch = uint32(epoch)

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to IncreaseStaking on stakingPlugin: Delete Candidate old power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
			"nodeId", can.NodeId.String(), "err", err)
		return err
	}

	can.Shares = new(big.Int).Add(can.Shares, amount)

	// set new power of can
	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to IncreaseStaking on stakingPlugin: Store Candidate new power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
			"nodeId", can.NodeId.String(), "err", err)
		return err
	}

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to IncreaseStaking on stakingPlugin: Store Candidate info is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
			"nodeId", can.NodeId.String(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) WithdrewStaking(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	can *staking.Candidate) error {
	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	log.Debug("Call WithdrewStaking", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"epoch", epoch, "nodeId", can.NodeId.String())

	lazyCalcStakeAmount(epoch, can)

	canAddr := crypto.PubkeyToAddress(*pubKey)

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to WithdrewStaking on stakingPlugin: Delete Candidate old power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}

	if err := sk.withdrewStakeAmount(state, blockHash, blockNumber.Uint64(), epoch, canAddr, can); nil != err {
		return err
	}

	can.StakingEpoch = uint32(epoch)

	if can.Released.Cmp(common.Big0) > 0 || can.RestrictingPlan.Cmp(common.Big0) > 0 {

		if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: Store Candidate info is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return err
		}
	} else {

		// Clean candidate info
		if err := sk.db.DelCandidateStore(blockHash, canAddr); nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: Delete Candidate info is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return err
		}
	}

	// sub the account staking Reference Count
	if err := sk.db.SubAccountStakeRc(blockHash, can.StakingAddress); nil != err {
		log.Error("Failed to WithdrewStaking on stakingPlugin: Store Staking Account Reference Count (sub) is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
			"staking Account", can.StakingAddress.String(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) withdrewStakeAmount(state xcom.StateDB, blockHash common.Hash, blockNumber, epoch uint64,
	canAddr common.Address, can *staking.Candidate) error {

	// Direct return of money during the hesitation period
	// Return according to the way of coming
	if can.ReleasedHes.Cmp(common.Big0) > 0 {

		state.AddBalance(can.StakingAddress, can.ReleasedHes)
		state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)
		can.ReleasedHes = new(big.Int).SetInt64(0)
	}

	if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {

		err := rt.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
		if nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
				"stakingAddr", can.StakingAddress.Hex(), "restrictingPlanHes", can.RestrictingPlanHes, "err", err)
			return err
		}

		can.RestrictingPlanHes = new(big.Int).SetInt64(0)
	}

	if can.Released.Cmp(common.Big0) > 0 || can.RestrictingPlan.Cmp(common.Big0) > 0 {

		if err := sk.addUnStakeItem(state, blockNumber, blockHash, epoch, can.NodeId, canAddr, can.StakingBlockNum); nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: Add UnStakeItemStore failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return err
		}
	}

	can.Shares = new(big.Int).SetInt64(0)
	can.Status |= staking.Invalided | staking.Withdrew

	return nil
}

func (sk *StakingPlugin) HandleUnCandidateItem(state xcom.StateDB, blockNumber uint64, blockHash common.Hash, epoch uint64) error {

	log.Debug("Call HandleUnCandidateItem start", "blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch)

	unStakeCount, err := sk.db.GetUnStakeCountStore(blockHash, epoch)
	switch {
	case nil != err && err != snapshotdb.ErrNotFound:
		return err
	case nil != err && err == snapshotdb.ErrNotFound:
		unStakeCount = 0
	}

	if unStakeCount == 0 {
		return nil
	}

	filterAddr := make(map[common.Address]struct{})

	for index := 1; index <= int(unStakeCount); index++ {

		stakeItem, err := sk.db.GetUnStakeItemStore(blockHash, epoch, uint64(index))
		if nil != err {
			log.Error("Failed to HandleUnCandidateItem: Query the unStakeItem node addr is failed",
				"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		canAddr := stakeItem.NodeAddress

		log.Debug("Call HandleUnCandidateItem: the candidate Addr",
			"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "addr", canAddr.Hex())

		if _, ok := filterAddr[canAddr]; ok {
			if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
				log.Error("Failed to HandleUnCandidateItem: Delete already handle unstakeItem failed",
					"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
				return err
			}
			continue
		}

		can, err := sk.db.GetCandidateStore(blockHash, canAddr)
		if nil != err && err != snapshotdb.ErrNotFound {
			log.Error("Failed to HandleUnCandidateItem: Query candidate failed",
				"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "canAddr", canAddr.Hex(), "err", err)
			return err
		}

		// This should not be nil
		if (nil != err && err == snapshotdb.ErrNotFound) || nil == can {

			if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
				log.Error("Failed to HandleUnCandidateItem: Candidate is no exist, Delete unstakeItem failed",
					"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
				return err
			}

			continue
		}

		// if the item stakingBlockNum is not enough the stakingBlockNum of candidate info
		if stakeItem.StakingBlockNum != can.StakingBlockNum {

			log.Warn("Call HandleUnCandidateItem: the item stakingBlockNum no equal current candidate stakingBlockNum",
				"item stakingBlockNum", stakeItem.StakingBlockNum, "candidate stakingBlockNum", can.StakingBlockNum)

			if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
				log.Error("Failed to HandleUnCandidateItem: The Item is invilad, cause the stakingBlockNum is less "+
					"than stakingBlockNum of curr candidate, Delete unstakeItem failed",
					"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
				return err
			}

			continue

		}

		// Second handle balabala ...
		if err := sk.handleUnStake(state, blockNumber, blockHash, epoch, canAddr, can); nil != err {
			return err
		}

		if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
			log.Error("Failed to HandleUnCandidateItem: Delete unstakeItem failed",
				"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		filterAddr[canAddr] = struct{}{}
	}

	if err := sk.db.DelUnStakeCountStore(blockHash, epoch); nil != err {
		log.Error("Failed to HandleUnCandidateItem: Delete unstakeCount failed",
			"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) handleUnStake(state xcom.StateDB, blockNumber uint64, blockHash common.Hash, epoch uint64,
	addr common.Address, can *staking.Candidate) error {

	log.Debug("Call handleUnStake Start", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"epoch", epoch, "nodeId", can.NodeId.String())

	lazyCalcStakeAmount(epoch, can)

	refundReleaseFn := func(balance *big.Int) *big.Int {
		if balance.Cmp(common.Big0) > 0 {

			state.AddBalance(can.StakingAddress, balance)
			state.SubBalance(vm.StakingContractAddr, balance)
			return new(big.Int).SetInt64(0)
		}
		return balance
	}

	can.ReleasedHes = refundReleaseFn(can.ReleasedHes)
	can.Released = refundReleaseFn(can.Released)

	refundRestrictingPlanFn := func(title string, balance *big.Int) (*big.Int, error) {

		if balance.Cmp(common.Big0) > 0 {
			err := rt.ReturnLockFunds(can.StakingAddress, balance, state)
			if nil != err {
				log.Error("Failed to HandleUnCandidateItem on stakingPlugin: call Restricting ReturnLockFunds() is failed",
					title, balance, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
					"stakingAddr", can.StakingAddress.Hex(), "err", err)
				return new(big.Int).SetInt64(0), err
			}
			return new(big.Int).SetInt64(0), nil
		}

		return balance, nil
	}

	if balance, err := refundRestrictingPlanFn("RestrictingPlanHes", can.RestrictingPlanHes); nil != err {
		return err
	} else {
		can.RestrictingPlanHes = balance
	}

	if balance, err := refundRestrictingPlanFn("RestrictingPlan", can.RestrictingPlan); nil != err {
		return err
	} else {
		can.RestrictingPlan = balance
	}

	// delete can info
	if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
		log.Error("Failed to HandleUnCandidateItem: Delete candidate info failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"nodeId", can.NodeId.String(), "err", err)
		return err
	}

	log.Debug("Call handleUnStake end", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"epoch", epoch, "nodeId", can.NodeId.String())
	return nil
}

func (sk *StakingPlugin) GetDelegateInfo(blockHash common.Hash, delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*staking.Delegation, error) {
	return sk.db.GetDelegateStore(blockHash, delAddr, nodeId, stakeBlockNumber)
}

func (sk *StakingPlugin) GetDelegateExInfo(blockHash common.Hash, delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*staking.DelegationEx, error) {

	del, err := sk.db.GetDelegateStore(blockHash, delAddr, nodeId, stakeBlockNumber)
	if nil != err {
		return nil, err
	}
	return &staking.DelegationEx{
		Addr:            delAddr,
		NodeId:          nodeId,
		StakingBlockNum: stakeBlockNumber,
		DelegationHex: staking.DelegationHex{
			DelegateEpoch:      del.DelegateEpoch,
			Released:           (*hexutil.Big)(del.Released),
			ReleasedHes:        (*hexutil.Big)(del.ReleasedHes),
			RestrictingPlan:    (*hexutil.Big)(del.RestrictingPlan),
			RestrictingPlanHes: (*hexutil.Big)(del.RestrictingPlanHes),
		},
	}, nil
}

func (sk *StakingPlugin) GetDelegateExCompactInfo(blockHash common.Hash, blockNumber uint64, delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*staking.DelegationEx, error) {

	del, err := sk.db.GetDelegateStore(blockHash, delAddr, nodeId, stakeBlockNumber)
	if nil != err {
		return nil, err
	}

	epoch := xutil.CalculateEpoch(blockNumber)

	lazyCalcDelegateAmount(epoch, del)

	return &staking.DelegationEx{
		Addr:            delAddr,
		NodeId:          nodeId,
		StakingBlockNum: stakeBlockNumber,
		DelegationHex: staking.DelegationHex{
			DelegateEpoch:      del.DelegateEpoch,
			Released:           (*hexutil.Big)(del.Released),
			ReleasedHes:        (*hexutil.Big)(del.ReleasedHes),
			RestrictingPlan:    (*hexutil.Big)(del.RestrictingPlan),
			RestrictingPlanHes: (*hexutil.Big)(del.RestrictingPlanHes),
		},
	}, nil
}

func (sk *StakingPlugin) GetDelegateInfoByIrr(delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*staking.Delegation, error) {

	return sk.db.GetDelegateStoreByIrr(delAddr, nodeId, stakeBlockNumber)
}

func (sk *StakingPlugin) GetDelegateExInfoByIrr(delAddr common.Address,
	nodeId discover.NodeID, stakeBlockNumber uint64) (*staking.DelegationEx, error) {

	del, err := sk.db.GetDelegateStoreByIrr(delAddr, nodeId, stakeBlockNumber)
	if nil != err {
		return nil, err
	}
	return &staking.DelegationEx{
		Addr:            delAddr,
		NodeId:          nodeId,
		StakingBlockNum: stakeBlockNumber,
		DelegationHex: staking.DelegationHex{
			DelegateEpoch:      del.DelegateEpoch,
			Released:           (*hexutil.Big)(del.Released),
			ReleasedHes:        (*hexutil.Big)(del.ReleasedHes),
			RestrictingPlan:    (*hexutil.Big)(del.RestrictingPlan),
			RestrictingPlanHes: (*hexutil.Big)(del.RestrictingPlanHes),
		},
	}, nil
}

func (sk *StakingPlugin) Delegate(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	delAddr common.Address, del *staking.Delegation, can *staking.Candidate, typ uint16, amount *big.Int) error {

	pubKey, _ := can.NodeId.Pubkey()
	canAddr := crypto.PubkeyToAddress(*pubKey)

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	log.Debug("Call Delegate", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch,
		"delAddr", delAddr.String(), "nodeId", can.NodeId.String(), "StakingNum", can.StakingBlockNum, "typ", typ,
		"amount", amount)

	lazyCalcDelegateAmount(epoch, del)

	if typ == FreeOrigin { // from account free von

		origin := state.GetBalance(delAddr)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to Delegate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.String(),
				"originVon", origin, "delegateVon", amount)
			return staking.ErrAccountVonNoEnough
		}
		state.SubBalance(delAddr, amount)
		state.AddBalance(vm.StakingContractAddr, amount)

		del.ReleasedHes = new(big.Int).Add(del.ReleasedHes, amount)

	} else if typ == RestrictingPlanOrigin { //  from account RestrictingPlan von

		err := rt.PledgeLockFunds(delAddr, amount, state)
		if nil != err {
			log.Error("Failed to Delegate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch,
				"delAddr", delAddr.String(), "nodeId", can.NodeId.String(), "StakingNum", can.StakingBlockNum,
				"amount", amount, "err", err)
			return err
		}

		del.RestrictingPlanHes = new(big.Int).Add(del.RestrictingPlanHes, amount)

	} else {

		log.Error("Failed to Delegate on stakingPlugin", "err", staking.ErrWrongVonOptType,
			"got type", typ, "need type", fmt.Sprintf("%d or %d", FreeOrigin, RestrictingPlanOrigin))
		return staking.ErrWrongVonOptType
	}

	del.DelegateEpoch = uint32(epoch)

	// set new delegate info
	if err := sk.db.SetDelegateStore(blockHash, delAddr, can.NodeId, can.StakingBlockNum, del); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Store Delegate info is failed",
			"delAddr", delAddr.String(), "nodeId", can.NodeId.String(), "StakingNum",
			can.StakingBlockNum, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Delete Candidate old power is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}

	// add the candidate power
	can.Shares = new(big.Int).Add(can.Shares, amount)

	// set new power of can
	if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Store Candidate new power is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}

	// update can info about Shares
	if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Store Candidate info is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) WithdrewDelegate(state xcom.StateDB, blockHash common.Hash, blockNumber, amount *big.Int,
	delAddr common.Address, nodeId discover.NodeID, stakingBlockNum uint64, del *staking.Delegation) error {

	log.Debug("Call WithdrewDelegate", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"delAddr", delAddr.String(), "nodeId", nodeId.String(), "StakingNum", stakingBlockNum, "amount", amount)

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to WithdrewDelegate on stakingPlugin: nodeId parse addr failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
		return err
	}

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	can, err := sk.db.GetCandidateStore(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to WithdrewDelegate on stakingPlugin: Query candidate info failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
		return err
	}

	total := calcDelegateTotalAmount(del)
	// First need to deduct the von that is being refunded
	if total.Cmp(amount) < 0 {
		log.Error("Failed to WithdrewDelegate on stakingPlugin: the amount of valid delegate is not enough",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "delegate amount", total,
			"withdrew amount", amount)
		return staking.ErrDelegateVonNoEnough
	}
	refundAmount := calcRealRefund(total, amount)
	realSub := refundAmount
	lazyCalcDelegateAmount(epoch, del)
	del.DelegateEpoch = uint32(epoch)

	switch {
	// Illegal parameter
	case nil != can && stakingBlockNum > can.StakingBlockNum:
		log.Error("Failed to WithdrewDelegate on stakingPlugin: the stakeBlockNum invalid",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "fn.stakeBlockNum", stakingBlockNum,
			"can.stakeBlockNum", can.StakingBlockNum)
		return staking.ErrBlockNumberDisordered
	default:
		log.Info("Call WithdrewDelegate", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.String(), "nodeId", nodeId.String(), "StakingNum", stakingBlockNum,
			"total", total, "amount", amount, "realSub", realSub)

		// handle delegate on Hesitate period
		if refundAmount.Cmp(common.Big0) > 0 {
			rm, rbalance, lbalance, err := rufundDelegateFn(refundAmount, del.ReleasedHes, del.RestrictingPlanHes, delAddr, state)
			if nil != err {
				log.Error("Failed  to WithdrewDelegate, refund the hesitate balance is failed", "blockNumber", blockNumber,
					"blockHash", blockHash.Hex(), "delAddr", delAddr.String(), "nodeId", nodeId.String(), "StakingNum", stakingBlockNum,
					"refund balance", refundAmount, "releaseHes", del.ReleasedHes, "restrictingPlanHes", del.RestrictingPlanHes, "err", err)
				return err
			}
			refundAmount, del.ReleasedHes, del.RestrictingPlanHes = rm, rbalance, lbalance
		}

		// handle delegate on Effective period
		if refundAmount.Cmp(common.Big0) > 0 {
			rm, rbalance, lbalance, err := rufundDelegateFn(refundAmount, del.Released, del.RestrictingPlan, delAddr, state)
			if nil != err {
				log.Error("Failed  to WithdrewDelegate, refund the no hesitate balance is failed", "blockNumber", blockNumber,
					"blockHash", blockHash.Hex(), "delAddr", delAddr.String(), "nodeId", nodeId.String(), "StakingNum", stakingBlockNum,
					"refund balance", refundAmount, "release", del.Released, "restrictingPlan", del.RestrictingPlan, "err", err)
				return err
			}
			refundAmount, del.Released, del.RestrictingPlan = rm, rbalance, lbalance
		}

		if refundAmount.Cmp(common.Big0) != 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: the withdrew ramain is not zero",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
				"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "del balance", total,
				"withdrew balance", amount, "realSub amount", realSub, "withdrew remain", refundAmount)
			return staking.ErrWrongWithdrewDelVonCalc
		}

		// If tatol had full sub,
		// then clean the delegate info
		if total.Cmp(realSub) == 0 {
			if err := sk.db.DelDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Delete detegate is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
					"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
				return err
			}
		} else {
			if err := sk.db.SetDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum, del); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Store detegate is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
					"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
				return err
			}
		}
	}

	if nil != can && stakingBlockNum == can.StakingBlockNum && staking.Is_Valid(can.Status) {
		if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: Delete candidate old power is failed", "blockNumber",
				blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
				"stakingBlockNum", stakingBlockNum, "err", err)
			return err
		}

		// change candidate shares
		if can.Shares.Cmp(realSub) > 0 {
			can.Shares = new(big.Int).Sub(can.Shares, realSub)
		} else {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: the candidate shares is no enough", "blockNumber",
				blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakingBlockNum",
				stakingBlockNum, "can shares", can.Shares, "real withdrew delegate amount", realSub)
			panic("the candidate shares is no enough")
		}

		if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: Store candidate info is failed", "blockNumber",
				blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
				"stakingBlockNum", stakingBlockNum, "err", err)
			return err
		}

		if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: Store candidate old power is failed", "blockNumber",
				blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
				"stakingBlockNum", stakingBlockNum, "err", err)
			return err
		}
	}
	return nil
}

func rufundDelegateFn(refundBalance, aboutRelease, aboutRestrictingPlan *big.Int, delAddr common.Address, state xcom.StateDB) (*big.Int, *big.Int, *big.Int, error) {

	refundTmp := refundBalance
	releaseTmp := aboutRelease
	restrictingPlanTmp := aboutRestrictingPlan

	subDelegateFn := func(source, sub *big.Int) (*big.Int, *big.Int) {
		state.AddBalance(delAddr, sub)
		state.SubBalance(vm.StakingContractAddr, sub)
		return new(big.Int).Sub(source, sub), new(big.Int).SetInt64(0)
	}

	// When remain is greater than or equal to del.ReleasedHes/del.Released
	if refundTmp.Cmp(common.Big0) > 0 {
		if refundTmp.Cmp(releaseTmp) >= 0 && releaseTmp.Cmp(common.Big0) > 0 {
			refundTmp, releaseTmp = subDelegateFn(refundTmp, releaseTmp)
		} else if refundTmp.Cmp(releaseTmp) < 0 {
			// When remain is less than or equal to del.ReleasedHes/del.Released
			releaseTmp, refundTmp = subDelegateFn(releaseTmp, refundTmp)
		}
	}

	if refundTmp.Cmp(common.Big0) > 0 {
		// When remain is greater than or equal to del.RestrictingPlanHes/del.RestrictingPlan
		if refundTmp.Cmp(restrictingPlanTmp) >= 0 && restrictingPlanTmp.Cmp(common.Big0) > 0 {
			err := rt.ReturnLockFunds(delAddr, restrictingPlanTmp, state)
			if nil != err {
				return refundTmp, releaseTmp, restrictingPlanTmp, err
			}
			refundTmp = new(big.Int).Sub(refundTmp, restrictingPlanTmp)
			restrictingPlanTmp = new(big.Int).SetInt64(0)
		} else if refundTmp.Cmp(restrictingPlanTmp) < 0 {
			// When remain is less than or equal to del.RestrictingPlanHes/del.RestrictingPlan
			err := rt.ReturnLockFunds(delAddr, refundTmp, state)
			if nil != err {
				return refundTmp, releaseTmp, restrictingPlanTmp, err
			}
			restrictingPlanTmp = new(big.Int).Sub(restrictingPlanTmp, refundTmp)
			refundTmp = new(big.Int).SetInt64(0)
		}
	}
	return refundTmp, releaseTmp, restrictingPlanTmp, nil
}

func (sk *StakingPlugin) ElectNextVerifierList(blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {

	log.Info("Call ElectNextVerifierList Start", "blockNumber", blockNumber, "blockHash", blockHash.Hex())

	oldVerifierArr, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to ElectNextVerifierList: No found the VerifierLIst", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// todo test
	xcom.PrintObject("Call ElectNextVerifierList old verifier list", oldVerifierArr)

	if oldVerifierArr.End != blockNumber {
		log.Error("Failed to ElectNextVerifierList: this blockNumber invalid", "Old Epoch End blockNumber",
			oldVerifierArr.End, "Current blockNumber", blockNumber)
		return staking.ErrBlockNumberDisordered
	}

	// caculate the new epoch start and end
	start := oldVerifierArr.End + 1
	end := oldVerifierArr.End + xutil.CalcBlocksEachEpoch()

	newVerifierArr := &staking.ValidatorArray{
		Start: start,
		End:   end,
	}

	currOriginVersion := gov.GetVersionForStaking(state)
	currVersion := xutil.CalcVersion(currOriginVersion)

	iter := sk.db.IteratorCandidatePowerByBlockHash(blockHash, int(xcom.EpochValidatorNum()))
	if err := iter.Error(); nil != err {
		log.Error("Failed to ElectNextVerifierList: take iter by candidate power is failed", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}
	defer iter.Release()

	queue := make(staking.ValidatorQueue, 0)

	// todo test
	count := 0

	for iter.Valid(); iter.Next(); {

		count++

		// todo test
		log.Debug("ElectNextVerifierList: iter", "key", hex.EncodeToString(iter.Key()))

		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := sk.db.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {

			log.Error("Failed to ElectNextVerifierList: Query Candidate info is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)

			return err
		}

		if can.ProgramVersion < currVersion {

			log.Warn("Warn ElectNextVerifierList: the can ProgramVersion is less than currVersion",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "canVersion",
				"nodeId", can.NodeId.String(), "canAddr", common.BytesToAddress(addrSuffix).Hex(),
				can.ProgramVersion, "currVersion", currVersion)

			// Low program version cannot be elected for epoch validator
			continue
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProgramVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			BlsPubKey:     can.BlsPubKey,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
	}

	if len(queue) == 0 {
		panic("Failed to ElectNextVerifierList: Select zero size validators~")
	}

	newVerifierArr.Arr = queue

	err = sk.setVerifierListAndIndex(blockNumber, blockHash, newVerifierArr)
	if nil != err {
		log.Error("Failed to ElectNextVerifierList: Set Next Epoch VerifierList is failed", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// todo test
	xcom.PrintObject("Call ElectNextVerifierList new verifier list", newVerifierArr)

	log.Info("Call ElectNextVerifierList end", "new epoch validators length", len(queue), "loopNum", count)
	return nil
}

func (sk *StakingPlugin) GetVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.ValidatorExQueue, error) {

	verifierList, err := sk.getVerifierList(blockHash, blockNumber, isCommit)
	if nil != err {
		return nil, err
	}

	if !isCommit && (blockNumber < verifierList.Start || blockNumber > verifierList.End) {

		log.Error("Failed to GetVerifierList", "start", verifierList.Start,
			"end", verifierList.End, "currentNumer", blockNumber)

		return nil, staking.ErrBlockNumberDisordered
	}

	queue := make(staking.ValidatorExQueue, len(verifierList.Arr))

	for i, v := range verifierList.Arr {

		var can *staking.Candidate
		if !isCommit {
			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetVerifierList, Quey candidate info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		} else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetVerifierList, Quey candidate info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		}

		shares, _ := new(big.Int).SetString(v.StakingWeight[1], 10)

		valEx := &staking.ValidatorEx{
			NodeId:          can.NodeId,
			BlsPubKey:       can.BlsPubKey,
			StakingAddress:  can.StakingAddress,
			BenefitAddress:  can.BenefitAddress,
			StakingTxIndex:  can.StakingTxIndex,
			ProgramVersion:  can.ProgramVersion,
			StakingBlockNum: can.StakingBlockNum,
			Shares:          (*hexutil.Big)(shares),
			Description:     can.Description,
			ValidatorTerm:   v.ValidatorTerm,
		}
		queue[i] = valEx
	}

	return queue, nil
}

func (sk *StakingPlugin) IsCurrVerifier(blockHash common.Hash, blockNumber uint64, nodeId discover.NodeID, isCommit bool) (bool, error) {

	verifierList, err := sk.getVerifierList(blockHash, blockNumber, isCommit)
	if nil != err {
		return false, err
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

	verifierList, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		return nil, err
	}

	if blockNumber < verifierList.Start || blockNumber > verifierList.End {

		log.Error("Failed to ListVerifierNodeID", "start", verifierList.Start,
			"end", verifierList.End, "currentNumer", blockNumber)

		return nil, staking.ErrBlockNumberDisordered
	}

	queue := make([]discover.NodeID, len(verifierList.Arr))

	for i, v := range verifierList.Arr {
		queue[i] = v.NodeId
	}
	return queue, nil
}

func (sk *StakingPlugin) GetCandidateONEpoch(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.CandidateQueue, error) {

	verifierList, err := sk.getVerifierList(blockHash, blockNumber, isCommit)
	if nil != err {
		return nil, err
	}

	queue := make(staking.CandidateQueue, len(verifierList.Arr))

	for i, v := range verifierList.Arr {
		var can *staking.Candidate
		if !isCommit {
			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetCandidateONEpoch, Quey candidate info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		} else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetCandidateONEpoch, Quey candidate info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		}
		queue[i] = can
	}

	return queue, nil
}

// flag:NOTE
// 0: 	Query previous round consensus validator
// 1:  	Query current round consensus validaor
// 2:  	Query next round consensus validator
func (sk *StakingPlugin) GetValidatorList(blockHash common.Hash, blockNumber uint64, flag uint, isCommit bool) (
	staking.ValidatorExQueue, error) {

	var validatorArr *staking.ValidatorArray

	switch flag {
	case PreviousRound:

		arr, err := sk.getPreValList(blockHash, blockNumber, isCommit)
		if nil != err {
			return nil, err
		}
		validatorArr = arr

	case CurrentRound:
		arr, err := sk.getCurrValList(blockHash, blockNumber, isCommit)
		if nil != err {
			return nil, err
		}
		validatorArr = arr
	case NextRound:
		arr, err := sk.getNextValList(blockHash, blockNumber, isCommit)
		if nil != err {
			return nil, err
		}
		validatorArr = arr
	default:
		log.Error("Failed to call GetValidatorList", "err", staking.ErrWrongFuncParams, "flag", flag)

		return nil, staking.ErrWrongFuncParams
	}

	queue := make(staking.ValidatorExQueue, len(validatorArr.Arr))

	for i, v := range validatorArr.Arr {

		var can *staking.Candidate

		if !isCommit {
			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetValidatorList, Quey candidate info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		} else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetValidatorList, Quey candidate info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		}

		shares, _ := new(big.Int).SetString(v.StakingWeight[1], 10)

		valEx := &staking.ValidatorEx{
			NodeId:          can.NodeId,
			BlsPubKey:       can.BlsPubKey,
			StakingAddress:  can.StakingAddress,
			BenefitAddress:  can.BenefitAddress,
			StakingTxIndex:  can.StakingTxIndex,
			ProgramVersion:  can.ProgramVersion,
			StakingBlockNum: can.StakingBlockNum,
			Shares:          (*hexutil.Big)(shares),
			Description:     can.Description,
			ValidatorTerm:   v.ValidatorTerm,
		}
		queue[i] = valEx
	}
	return queue, nil
}

func (sk *StakingPlugin) GetCandidateONRound(blockHash common.Hash, blockNumber uint64,
	flag uint, isCommit bool) (staking.CandidateQueue, error) {

	var validatorArr *staking.ValidatorArray

	switch flag {
	case PreviousRound:
		arr, err := sk.getPreValList(blockHash, blockNumber, isCommit)
		if nil != err {
			return nil, err
		}
		validatorArr = arr
	case CurrentRound:
		arr, err := sk.getCurrValList(blockHash, blockNumber, isCommit)
		if nil != err {
			return nil, err
		}
		validatorArr = arr
	case NextRound:
		arr, err := sk.getNextValList(blockHash, blockNumber, isCommit)
		if nil != err {
			return nil, err
		}
		validatorArr = arr
	default:
		log.Error("Failed to call GetCandidateONRound", "err", staking.ErrWrongFuncParams, "flag", flag)

		return nil, staking.ErrWrongFuncParams

	}

	queue := make(staking.CandidateQueue, len(validatorArr.Arr))

	for i, v := range validatorArr.Arr {

		var can *staking.Candidate

		if !isCommit {

			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetCandidateONRound, Quey candidate info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}

			can = c
		} else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetCandidateONRound, Quey candidate info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		}
		queue[i] = can
	}
	return queue, nil
}

func (sk *StakingPlugin) ListCurrentValidatorID(blockHash common.Hash, blockNumber uint64) ([]discover.NodeID, error) {

	arr, err := sk.getCurrValList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		return nil, err
	}

	queue := make([]discover.NodeID, len(arr.Arr))

	for i, candidate := range arr.Arr {
		queue[i] = candidate.NodeId
	}
	return queue, err
}

func (sk *StakingPlugin) IsCurrValidator(blockHash common.Hash, blockNumber uint64, nodeId discover.NodeID, isCommit bool) (bool, error) {

	validatorArr, err := sk.getCurrValList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		return false, err
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

func (sk *StakingPlugin) GetCandidateList(blockHash common.Hash, blockNumber uint64) (staking.CandidateHexQueue, error) {

	epoch := xutil.CalculateEpoch(blockNumber)

	iter := sk.db.IteratorCandidatePowerByBlockHash(blockHash, 0)
	if err := iter.Error(); nil != err {
		return nil, err
	}
	defer iter.Release()

	queue := make(staking.CandidateHexQueue, 0)

	count := 0

	for iter.Valid(); iter.Next(); {

		count++

		// todo test
		log.Debug("GetCandidateList: iter", "key", hex.EncodeToString(iter.Key()))

		addrSuffix := iter.Value()
		can, err := sk.db.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			return nil, err
		}

		lazyCalcStakeAmount(epoch, can)
		canHex := buildCanHex(can)
		queue = append(queue, canHex)
	}

	// todo test
	log.Debug("GetCandidateList: loop count", "count", count)

	return queue, nil
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
	} else {
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

func (sk *StakingPlugin) GetRelatedListByDelAddr(blockHash common.Hash, addr common.Address) (staking.DelRelatedQueue, error) {

	iter := sk.db.IteratorDelegateByBlockHashWithAddr(blockHash, addr, 0)
	if err := iter.Error(); nil != err {
		return nil, err
	}
	defer iter.Release()

	queue := make(staking.DelRelatedQueue, 0)

	for iter.Valid(); iter.Next(); {
		key := iter.Key()

		prefixLen := len(staking.DelegateKeyPrefix)

		nodeIdLen := discover.NodeIDBits / 8

		// delAddr
		delAddrByte := key[prefixLen : prefixLen+common.AddressLength]
		delAddr := common.BytesToAddress(delAddrByte)

		// nodeId
		nodeIdByte := key[prefixLen+common.AddressLength : prefixLen+common.AddressLength+nodeIdLen]
		nodeId := discover.MustBytesID(nodeIdByte)

		// stakenum
		stakeNumByte := key[prefixLen+common.AddressLength+nodeIdLen:]

		num := common.BytesToUint64(stakeNumByte)

		// related
		related := &staking.DelegateRelated{
			Addr:            delAddr,
			NodeId:          nodeId,
			StakingBlockNum: num,
		}
		queue = append(queue, related)
	}
	return queue, nil
}

func (sk *StakingPlugin) Election(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {

	log.Info("Call Election Start", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.Hex())

	blockNumber := header.Number.Uint64()

	// the validators of Current Epoch
	verifiers, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to call Election: No found current epoch validators", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return staking.ErrValidatorNoExist
	}

	// the validators of Current Round
	curr, err := sk.getCurrValList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to Election: No found the current round validators", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return staking.ErrValidatorNoExist
	}

	log.Info("Call Election start", "Curr epoch validators length", len(verifiers.Arr), "Curr round validators length", len(curr.Arr))
	xcom.PrintObject("Call Election Curr validators", curr)

	if blockNumber != (curr.End - xcom.ElectionDistance()) {
		log.Error("Failed to Election: Current blockNumber invalid", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"Target blockNumber", curr.End-xcom.ElectionDistance())
		return staking.ErrBlockNumberDisordered
	}

	// Never match, maybe!!!
	if nil == verifiers || len(verifiers.Arr) == 0 {
		panic("The Current Epoch VerifierList is empty, blockNumber: " + fmt.Sprint(blockNumber))
	}

	// caculate the next round start and end
	start := curr.End + 1
	end := curr.End + xutil.ConsensusSize()

	hasSlashLen := 0 // duplicateSign And lowRatio No enough von
	needRMwithdrewLen := 0
	needRMLowVersionLen := 0
	invalidLen := 0 // the num that the can need to remove

	removeCans := make(staking.NeedRemoveCans) // the candidates need to remove
	withdrewCans := make(staking.CandidateMap) // the candidates had withdrew
	// TODO test
	slashAddrQueue := make([]discover.NodeID, 0)
	withdrewQueue := make([]discover.NodeID, 0)
	lowVersionQueue := make([]discover.NodeID, 0)
	// need to clean lowRatio status
	lowRatioValidAddrs := make([]common.Address, 0)                 // The addr of candidate that need to clean lowRatio status
	lowRatioValidMap := make(map[common.Address]*staking.Candidate) // The map collect candidate info that need to clean lowRatio status

	// Query Valid programVersion
	originVersion := gov.GetVersionForStaking(state)
	currVersion := xutil.CalcVersion(originVersion)

	// Collecting removed as a result of being slashed
	// That is not withdrew to invalid
	//
	// eg. (lowRatio and must delete) OR (lowRatio and balance no enough) OR duplicateSign
	//
	checkHaveSlash := func(status uint32) bool {
		return staking.Is_Invalid_LowRatioDel(status) ||
			staking.Is_Invalid_LowRatio_NotEnough(status) ||
			staking.Is_Invalid_DuplicateSign(status)
	}

	currMap := make(map[discover.NodeID]struct{}, len(curr.Arr))
	for _, v := range curr.Arr {

		canAddr, _ := xutil.NodeId2Addr(v.NodeId)
		can, err := sk.db.GetCandidateStore(blockHash, canAddr)
		if nil != err {
			log.Error("Failed to Query Candidate Info on Election", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(), "err", err)
			return err
		}

		var isSlash bool

		if checkHaveSlash(can.Status) {
			// -- if staking.Is_Invalid(can.Status) && !staking.Is_Withdrew(can.Status) {

			removeCans[v.NodeId] = can
			slashAddrQueue = append(slashAddrQueue, v.NodeId)
			hasSlashLen++
			isSlash = true
		}

		// Collecting candidate information that active withdrawal
		if staking.Is_Invalid_Withdrew(can.Status) && !isSlash {

			withdrewCans[v.NodeId] = can
			withdrewQueue = append(withdrewQueue, v.NodeId)
		}

		// valid AND lowRatio status, that candidate need to clean the lowRatio status
		if staking.Is_Valid(can.Status) && staking.Is_LowRatio(can.Status) {
			lowRatioValidAddrs = append(lowRatioValidAddrs, canAddr)
			lowRatioValidMap[canAddr] = can
		}

		// Collect candidate who need to be removed
		// from the validators because the version is too low
		if can.ProgramVersion < currVersion {
			removeCans[v.NodeId] = can
			lowVersionQueue = append(lowVersionQueue, v.NodeId)
		}

		currMap[v.NodeId] = struct{}{}
	}
	needRMLowVersionLen = len(lowVersionQueue)

	// Exclude the current consensus round validators from the validators of the Epoch
	diffQueue := make(staking.ValidatorQueue, 0)
	for _, v := range verifiers.Arr {

		if _, ok := withdrewCans[v.NodeId]; ok {

			delete(withdrewCans, v.NodeId)

		}

		if _, ok := currMap[v.NodeId]; ok {
			continue
		}

		addr, _ := xutil.NodeId2Addr(v.NodeId)
		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			log.Error("Failed to Get Candidate on Election", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(), "err", err)
			return err
		}

		// Jump the slashed candidate
		if checkHaveSlash(can.Status) {
			continue
		}

		// Ignore the low version
		if can.ProgramVersion < currVersion {
			continue
		}

		diffQueue = append(diffQueue, v)
	}

	for i := 0; i < len(withdrewQueue); i++ {

		nodeId := withdrewQueue[i]

		if can, ok := withdrewCans[nodeId]; !ok {
			// remove the can on withdrewqueue
			withdrewQueue = append(withdrewQueue[:i], withdrewQueue[i+1:]...)
			i--
		} else {
			// append to the collection that needs to be removed
			removeCans[nodeId] = can
		}

	}
	needRMwithdrewLen = len(withdrewQueue)

	invalidLen = hasSlashLen + needRMwithdrewLen + needRMLowVersionLen

	shuffle := func(invalidLen int, currQueue, vrfQueue staking.ValidatorQueue) staking.ValidatorQueue {
		currQueue.ValidatorSort(removeCans, staking.CompareForDel)
		// Increase term of validator
		copyCurrQueue := make(staking.ValidatorQueue, len(currQueue))
		copy(copyCurrQueue, currQueue)
		for i, v := range copyCurrQueue {
			v.ValidatorTerm++
			copyCurrQueue[i] = v
		}
		// Remove the invalid validators
		copyCurrQueue = copyCurrQueue[invalidLen:]
		return shuffleQueue(copyCurrQueue, vrfQueue)
	}

	var vrfQueue staking.ValidatorQueue
	var vrfLen int
	if len(diffQueue) > int(xcom.ConsValidatorNum()) {
		vrfLen = int(xcom.ConsValidatorNum())
	} else {
		vrfLen = len(diffQueue)
	}

	if vrfLen != 0 {
		if queue, err := vrfElection(diffQueue, vrfLen, header.Nonce.Bytes(), header.ParentHash); nil != err {
			log.Error("Failed to VrfElection on Election",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		} else {
			vrfQueue = queue
		}
	}

	log.Info("Call Election, statistics need to remove node num",
		"has slash count", hasSlashLen, "withdrew and need remove count",
		needRMwithdrewLen, "low version need remove count", needRMLowVersionLen,
		"total remove count", invalidLen, "remove map size", len(removeCans),
		"current validators Size", len(curr.Arr), "ConsValidatorNum", xcom.ConsValidatorNum(),
		"ShiftValidatorNum", xcom.ShiftValidatorNum(), "diffQueueLen", len(diffQueue),
		"vrfQueueLen", len(vrfQueue))

	nextQueue := shuffle(invalidLen, curr.Arr, vrfQueue)

	if len(nextQueue) == 0 {
		panic("The Next Round Validator is empty, blockNumber: " + fmt.Sprint(blockNumber))
	}

	next := &staking.ValidatorArray{
		Start: start,
		End:   end,
		Arr:   nextQueue,
	}

	if err := sk.setRoundValListAndIndex(blockNumber, blockHash, next); nil != err {
		log.Error("Failed to SetNextValidatorList on Election", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// todo test
	if len(slashAddrQueue) != 0 {
		xcom.PrintObject("Election Remove Slashing nodeId", slashAddrQueue)
	}

	// todo test
	if len(withdrewQueue) != 0 {
		xcom.PrintObject("Election Remove Withdrew nodeId", withdrewQueue)
	}

	// todo test
	if len(lowVersionQueue) != 0 {
		xcom.PrintObject("Election Remove Low version nodeId", lowVersionQueue)
	}

	// update candidate status
	// Must sort
	for _, canAddr := range lowRatioValidAddrs {

		can := lowRatioValidMap[canAddr]

		// clean the low package ratio status
		can.Status &^= staking.LowRatio

		// TODO test
		xcom.PrintObject("Call Election, clean lowratio, nodeId:"+can.NodeId.String()+", can Info:", can)

		addr, _ := xutil.NodeId2Addr(can.NodeId)
		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Failed to Store Candidate on Election", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return err
		}
	}

	if err := sk.storeRoundValidatorAddrs(blockHash, start, nextQueue); nil != err {
		log.Error("Failed to storeRoundValidatorAddrs on Election", "blockNumber", blockNumber,
			"blockHash", blockHash.TerminalString(), "err", err)
	}

	log.Info("Call Election end", "next round validators length", len(nextQueue))

	// todo test
	xcom.PrintObject("Call Election Next validators", next)
	return nil
}

func shuffleQueue(remainCurrQueue, vrfQueue staking.ValidatorQueue) staking.ValidatorQueue {

	remainLen := len(remainCurrQueue)
	totalQueue := append(remainCurrQueue, vrfQueue...)

	for remainLen > int(xcom.ConsValidatorNum()-xcom.ShiftValidatorNum()) && len(totalQueue) > int(xcom.ConsValidatorNum()) {
		totalQueue = totalQueue[1:]
		remainLen--
	}

	if len(totalQueue) > int(xcom.ConsValidatorNum()) {
		totalQueue = totalQueue[:xcom.ConsValidatorNum()]
	}

	next := make(staking.ValidatorQueue, len(totalQueue))

	copy(next, totalQueue)

	// re-sort before store next validators
	next.ValidatorSort(nil, staking.CompareForStore)
	return next
}

// NotifyPunishedVerifiers
func (sk *StakingPlugin) SlashCandidates(state xcom.StateDB, blockHash common.Hash, blockNumber uint64, queue ...*staking.SlashNodeItem) error {

	invalidNodeIdMap := make(map[discover.NodeID]struct{}, 0)

	for _, slashItem := range queue {
		needRemove, err := sk.toSlash(state, blockNumber, blockHash, slashItem)
		if nil != err {
			return err
		}

		if needRemove {
			invalidNodeIdMap[slashItem.NodeId] = struct{}{}
		}
	}

	// remove the validator from epoch verifierList
	if err := sk.removeFromVerifiers(blockNumber, blockHash, invalidNodeIdMap); nil != err {
		return err
	}

	// notify gov to do somethings
	if err := gov.NotifyPunishedVerifiers(blockHash, invalidNodeIdMap, state); nil != err {
		log.Error("Failed to SlashCandidates: call NotifyPunishedVerifiers of gov is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "invalidNodeId Size", len(invalidNodeIdMap), "err", err)
		return err
	}
	return nil
}

func (sk *StakingPlugin) toSlash(state xcom.StateDB, blockNumber uint64, blockHash common.Hash, slashItem *staking.SlashNodeItem) (bool, error) {

	log.Info("Call SlashCandidates: call toSlash", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"nodeId", slashItem.NodeId.String(), "amount", slashItem.Amount, "slashType", slashItem.SlashType,
		"benefitAddr", slashItem.BenefitAddr.Hex())

	var needRemove bool

	// check slash type is right
	slashTypeIsWrong := func() bool {
		return uint32(slashItem.SlashType) != staking.LowRatio &&
			uint32(slashItem.SlashType) != staking.LowRatioDel &&
			uint32(slashItem.SlashType) != staking.DuplicateSign
	}
	if slashTypeIsWrong() {
		log.Error("Failed to SlashCandidates: the slashType is wrong", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "slashType", slashItem.SlashType, "benefitAddr", slashItem.BenefitAddr.Hex())

		return needRemove, staking.ErrWrongSlashType
	}

	canAddr, _ := xutil.NodeId2Addr(slashItem.NodeId)
	can, err := sk.db.GetCandidateStore(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to SlashCandidates: Query can is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
		return needRemove, err
	}

	if nil == can {
		log.Error("Failed to SlashCandidates: the can is empty", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String())
		return needRemove, staking.ErrCanNoExist
	}

	epoch := xutil.CalculateEpoch(blockNumber)

	lazyCalcStakeAmount(epoch, can)

	// Balance that can only be effective for Slash
	total := new(big.Int).Add(can.Released, can.RestrictingPlan)

	if total.Cmp(slashItem.Amount) < 0 {

		log.Warn("Warned to SlashCandidates: the candidate total staking amount is not enough",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(),
			"candidate total amount", total, "slashing amount", slashItem.Amount)

		return needRemove, staking.ErrSlashVonOverflow
	}

	// clean the candidate power, first
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to SlashCandidates: Delete candidate old power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String())
		return needRemove, err
	}

	slashBalance := slashItem.Amount

	// slash the balance
	if slashBalance.Cmp(common.Big0) > 0 && can.Released.Cmp(common.Big0) > 0 {
		val, rval, err := slashBalanceFn(slashBalance, can.Released, false, uint32(slashItem.SlashType),
			slashItem.BenefitAddr, can.StakingAddress, state)
		if nil != err {
			log.Error("Failed to SlashCandidates: slash Released", "slashed amount", slashBalance,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
			return needRemove, err
		}
		slashBalance, can.Released = val, rval
	}
	if slashBalance.Cmp(common.Big0) > 0 && can.RestrictingPlan.Cmp(common.Big0) > 0 {
		val, rval, err := slashBalanceFn(slashBalance, can.RestrictingPlan, true, uint32(slashItem.SlashType),
			slashItem.BenefitAddr, can.StakingAddress, state)
		if nil != err {
			log.Error("Failed to SlashCandidates: slash RestrictingPlan", "slashed amount", slashBalance,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
			return needRemove, err
		}
		slashBalance, can.RestrictingPlan = val, rval
	}

	// check slash remain balance
	if slashBalance.Cmp(common.Big0) != 0 {
		log.Error("Failed to SlashCandidates: the ramain is not zero",
			"slashAmount", slashItem.Amount, "slashed remain", slashBalance,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String())
		return needRemove, staking.ErrWrongSlashVonCalc
	}

	sharesHaveBeenClean := func() bool {
		return (staking.Is_Invalid_LowRatio_NotEnough(can.Status) ||
			staking.Is_Invalid_LowRatioDel(can.Status) ||
			staking.Is_Invalid_DuplicateSign(can.Status) ||
			staking.Is_Invalid_Withdrew(can.Status))
	}

	// If the shares is zero, don't need to sub shares
	if !sharesHaveBeenClean() {

		// first slash and no withdrew
		// sub Shares to effect power
		if can.Shares.Cmp(slashItem.Amount) >= 0 {
			can.Shares = new(big.Int).Sub(can.Shares, slashItem.Amount)
		} else {
			log.Error("Failed to SlashCandidates: the candidate shares is no enough", "slashType", slashItem.SlashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "candidate shares",
				can.Shares, "slash amount", slashItem.Amount)
			panic("the candidate shares is no enough")
		}
	}

	// need invalid candidate status
	// need remove from verifierList
	needInvalid, needRemove, changeStatus := handleSlashTypeFn(slashItem.SlashType, calcCandidateTotalAmount(can))

	log.Info("Call SlashCandidates: the status", "needInvalid", needInvalid,
		"needRemove", needRemove, "current can.Status", can.Status, "need to superpose status", changeStatus)

	if needInvalid && staking.Is_Valid(can.Status) {

		if can.ReleasedHes.Cmp(common.Big0) > 0 {

			state.AddBalance(can.StakingAddress, can.ReleasedHes)
			state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)
			can.ReleasedHes = new(big.Int).SetInt64(0)
		}
		if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {

			err := rt.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
			if nil != err {
				log.Error("Failed to SlashCandidates on stakingPlugin: call Restricting ReturnLockFunds() is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "stakingAddr", can.StakingAddress.Hex(),
					"restrictingPlanHes", can.RestrictingPlanHes, "err", err)
				return needRemove, err
			}
			can.RestrictingPlanHes = new(big.Int).SetInt64(0)
		}

		// need to sub account rc
		if err := sk.db.SubAccountStakeRc(blockHash, can.StakingAddress); nil != err {
			log.Error("Failed to SlashCandidates: Sub Account staking Reference Count is failed", "slashType", slashItem.SlashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
			return needRemove, err
		}

		// Must be guaranteed to be the first slash to invalid can status and no active withdrewStake
		if err := sk.addUnStakeItem(state, blockNumber, blockHash, epoch, can.NodeId, canAddr, can.StakingBlockNum); nil != err {
			log.Error("Failed to SlashCandidates on stakingPlugin: Add UnStakeItemStore failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return needRemove, err
		}

		//because of deleted candidate info ,clean Shares
		can.Shares = new(big.Int).SetInt64(0)
		can.Status |= changeStatus
		if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to SlashCandidates on stakingPlugin: Store Candidate info is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return needRemove, err
		}

	} else if !needInvalid && staking.Is_Valid(can.Status) {

		// update the candidate power, If do not need to delete power (the candidate status still be valid)
		if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to SlashCandidates: Store candidate power is failed", "slashType", slashItem.SlashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
			return needRemove, err
		}

		can.Status |= changeStatus
		if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to SlashCandidates: Store candidate is failed", "slashType", slashItem.SlashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
			return needRemove, err
		}

	} else {

		can.Status |= changeStatus
		if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to SlashCandidates: Store candidate is failed", "slashType", slashItem.SlashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
			return needRemove, err
		}
	}
	return needRemove, nil
}

func (sk *StakingPlugin) removeFromVerifiers(blockNumber uint64, blockHash common.Hash, slashNodeIdMap map[discover.NodeID]struct{}) error {
	verifier, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to SlashCandidates: Query Verifier List is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeIdQueue Size", len(slashNodeIdMap), "err", err)
		return err
	}

	// remove the val from epoch validators,
	// because the candidate status is invalid after slashed
	orginLen := len(verifier.Arr)
	for i := 0; i < len(verifier.Arr); i++ {

		val := verifier.Arr[i]

		if _, ok := slashNodeIdMap[val.NodeId]; ok {

			log.Info("Call SlashCandidates, Delete the validator", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", val.NodeId.String())

			verifier.Arr = append(verifier.Arr[:i], verifier.Arr[i+1:]...)
			i--
			break
		}
	}

	dirtyLen := len(verifier.Arr)

	if dirtyLen != orginLen {

		if err := sk.setVerifierListByIndex(blockNumber, blockHash, verifier); nil != err {
			log.Error("Failed to SlashCandidates: Store Verifier List is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}
	return nil
}

func handleSlashTypeFn(slashType int, remain *big.Int) (bool, bool, uint32) {

	var needInvalid, needRemove bool // need invalid candidate status And need remove from verifierList
	var changeStatus uint32          // need to add this status

	switch slashType {
	case staking.LowRatio:

		if !xutil.CheckStakeThreshold(remain) {
			changeStatus |= staking.NotEnough
			changeStatus |= staking.Invalided
			needInvalid = true
			needRemove = true
		}
	case staking.LowRatioDel:
		changeStatus |= staking.Invalided
		needInvalid = true
		needRemove = true
	case staking.DuplicateSign:
		changeStatus |= staking.Invalided
		needInvalid = true
		needRemove = true
	}
	changeStatus |= uint32(slashType)

	return needInvalid, needRemove, changeStatus
}

func slashBalanceFn(slashAmount, canBalance *big.Int, isNotify bool,
	slashType uint32, benefitAddr, stakingAddr common.Address, state xcom.StateDB) (*big.Int, *big.Int, error) {

	// check zero value
	// If there is a zero value, no logic is done.
	if canBalance.Cmp(common.Big0) == 0 || slashAmount.Cmp(common.Big0) == 0 {
		return slashAmount, canBalance, nil
	}

	slashAmountTmp := new(big.Int).SetInt64(0)
	balanceTmp := new(big.Int).SetInt64(0)

	if slashAmount.Cmp(canBalance) >= 0 {

		state.SubBalance(vm.StakingContractAddr, canBalance)

		if staking.Is_DuplicateSign(uint32(slashType)) {
			state.AddBalance(benefitAddr, canBalance)
		} else {
			state.AddBalance(vm.RewardManagerPoolAddr, canBalance)
		}

		if isNotify {
			err := rt.SlashingNotify(stakingAddr, canBalance, state)
			if nil != err {
				return slashAmountTmp, balanceTmp, err
			}
		}

		slashAmountTmp = new(big.Int).Sub(slashAmount, canBalance)
		balanceTmp = new(big.Int).SetInt64(0)

	} else {
		state.SubBalance(vm.StakingContractAddr, slashAmount)
		if staking.Is_DuplicateSign(uint32(slashType)) {
			state.AddBalance(benefitAddr, slashAmount)
		} else {
			state.AddBalance(vm.RewardManagerPoolAddr, slashAmount)
		}

		if isNotify {
			err := rt.SlashingNotify(stakingAddr, slashAmount, state)
			if nil != err {
				return slashAmountTmp, balanceTmp, err
			}
		}

		slashAmountTmp = new(big.Int).SetInt64(0)
		balanceTmp = new(big.Int).Sub(canBalance, slashAmount)
	}

	return slashAmountTmp, balanceTmp, nil
}

func (sk *StakingPlugin) ProposalPassedNotify(blockHash common.Hash, blockNumber uint64, nodeIds []discover.NodeID,
	programVersion uint32) error {

	log.Debug("Call ProposalPassedNotify to promote candidate programVersion", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "version", programVersion, "nodeIdQueueSize", len(nodeIds))

	version := xutil.CalcVersion(programVersion)

	for _, nodeId := range nodeIds {

		log.Info("Call ProposalPassedNotify: promote candidate start", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "real version", programVersion, "calc version", version, "nodeId", nodeId.String())

		addr, _ := xutil.NodeId2Addr(nodeId)
		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err && err != snapshotdb.ErrNotFound {
			log.Error("Failed to ProposalPassedNotify: Query Candidate is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		if (nil != err && err == snapshotdb.ErrNotFound) || nil == can {
			log.Error("Failed to ProposalPassedNotify: Promote candidate programVersion failed, the can is empty",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
			continue
		}

		if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
			log.Error("Failed to ProposalPassedNotify: Delete Candidate old power is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		can.ProgramVersion = version

		if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
			log.Error("Failed to ProposalPassedNotify: Store Candidate new power is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Failed to ProposalPassedNotify: Store Candidate info is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

	}

	return nil
}

func (sk *StakingPlugin) DeclarePromoteNotify(blockHash common.Hash, blockNumber uint64, nodeId discover.NodeID,
	programVersion uint32) error {

	version := xutil.CalcVersion(programVersion)

	log.Info("Call DeclarePromoteNotify to promote candidate programVersion", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "real version", programVersion, "calc version", version, "nodeId", nodeId.String())

	addr, _ := xutil.NodeId2Addr(nodeId)
	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to DeclarePromoteNotify: Query Candidate is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	if (nil != err && err == snapshotdb.ErrNotFound) || nil == can {

		log.Error("Failed to DeclarePromoteNotify: Promote candidate programVersion failed, the can is empty",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(),
			"version", programVersion)
		return nil
	}

	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to DeclarePromoteNotify: Delete Candidate old power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	can.ProgramVersion = version

	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to DeclarePromoteNotify: Store Candidate new power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to DeclarePromoteNotify: Store Candidate info is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) GetLastNumber(blockNumber uint64) uint64 {

	valArr, err := sk.getCurrValList(common.ZeroHash, blockNumber, QueryStartIrr)
	if nil != err {
		log.Error("Failed to GetLastNumber", "blockNumber", blockNumber, "err", err)
		return 0
	}

	if nil == err && nil != valArr {
		return valArr.End
	}
	return 0
}

func (sk *StakingPlugin) GetValidator(blockNumber uint64) (*cbfttypes.Validators, error) {

	valArr, err := sk.getCurrValList(common.ZeroHash, blockNumber, QueryStartIrr)
	if nil != err && err != snapshotdb.ErrNotFound {
		return nil, err
	}

	if nil == err && nil != valArr {
		return buildCbftValidators(valArr.Start, valArr.Arr), nil
	}
	return nil, fmt.Errorf("No Found Validators by blockNumber: %d", blockNumber)
}

// NOTE: Verify that it is the validator of the current Epoch
func (sk *StakingPlugin) IsCandidateNode(nodeID discover.NodeID) bool {

	indexs, err := sk.db.GetEpochValIndexByIrr()
	if nil != err {
		log.Error("Failed to IsCandidateNode: query epoch validators indexArr is failed", "err", err)
		return false
	}

	isCandidate := false

	for i, indexInfo := range indexs {
		queue, err := sk.db.GetEpochValListByIrr(indexInfo.Start, indexInfo.End)
		if nil != err {
			log.Error("Failed to IsCandidateNode: Query epoch validators is failed",
				"index length", len(indexs), "the loop number", i+1, "Start", indexInfo.Start, "End", indexInfo.End, "err", err)
			continue
		} else {
			for _, val := range queue {
				if val.NodeId == nodeID {
					isCandidate = true
					goto label
				}
			}
		}
	}
label:
	log.Info("Call IsCandidateNode", "nodeId", nodeID.String(), "isCandidate", isCandidate)
	return isCandidate
}

func buildCbftValidators(start uint64, arr staking.ValidatorQueue) *cbfttypes.Validators {
	valMap := make(cbfttypes.ValidateNodeMap, len(arr))

	for i, v := range arr {

		pubKey, _ := v.NodeId.Pubkey()

		vn := &cbfttypes.ValidateNode{
			Index:     uint32(i),
			Address:   v.NodeAddress,
			PubKey:    pubKey,
			NodeID:    v.NodeId,
			BlsPubKey: &v.BlsPubKey,
		}

		valMap[v.NodeId] = vn
	}

	res := &cbfttypes.Validators{
		Nodes:            valMap,
		ValidBlockNumber: start,
	}
	return res
}

func lazyCalcStakeAmount(epoch uint64, can *staking.Candidate) {

	changeAmountEpoch := can.StakingEpoch

	sub := epoch - uint64(changeAmountEpoch)

	// todo test
	xcom.PrintObject("lazyCalcStakeAmount before, epoch:"+fmt.Sprint(epoch)+", can", can)

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio() {
		return
	}

	if can.ReleasedHes.Cmp(common.Big0) > 0 {
		can.Released = new(big.Int).Add(can.Released, can.ReleasedHes)
		can.ReleasedHes = new(big.Int).SetInt64(0)
	}

	if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		can.RestrictingPlan = new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
		can.RestrictingPlanHes = new(big.Int).SetInt64(0)
	}

	// todo test
	xcom.PrintObject("lazyCalcStakeAmount end, epoch:"+fmt.Sprint(epoch)+", can", can)

}

func lazyCalcDelegateAmount(epoch uint64, del *staking.Delegation) {

	// When the first time, there was no previous changeAmountEpoch
	if del.DelegateEpoch == 0 {
		return
	}

	changeAmountEpoch := del.DelegateEpoch

	sub := epoch - uint64(changeAmountEpoch)

	// todo test
	xcom.PrintObject("lazyCalcDelegateAmount before, epoch:"+fmt.Sprint(epoch)+", del", del)

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio() {
		return
	}

	if del.ReleasedHes.Cmp(common.Big0) > 0 {
		del.Released = new(big.Int).Add(del.Released, del.ReleasedHes)
		del.ReleasedHes = new(big.Int).SetInt64(0)
	}

	if del.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		del.RestrictingPlan = new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
		del.RestrictingPlanHes = new(big.Int).SetInt64(0)
	}

	// todo test
	xcom.PrintObject("lazyCalcDelegateAmount end, epoch:"+fmt.Sprint(epoch)+", del", del)

}

type sortValidator struct {
	v           *staking.Validator
	x           int64
	weights     int64
	version     uint32
	blockNumber uint64
	txIndex     uint32
}

type sortValidatorQueue []*sortValidator

func (svs sortValidatorQueue) Len() int {
	return len(svs)
}

func (svs sortValidatorQueue) Less(i, j int) bool {
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

func (svs sortValidatorQueue) Swap(i, j int) {
	svs[i], svs[j] = svs[j], svs[i]
}

// Elected verifier by vrf random election
// validatorList：Waiting for the elected node
// nonce：Vrf proof of the current block
// parentHash：Parent block hash
func vrfElection(validatorList staking.ValidatorQueue, shiftLen int, nonce []byte, parentHash common.Hash) (staking.ValidatorQueue, error) {
	preNonces, err := handler.GetVrfHandlerInstance().Load(parentHash)
	if nil != err {
		return nil, err
	}
	if len(preNonces) < len(validatorList) {
		log.Error("Failed to vrfElection on Election", "validatorListSize", len(validatorList),
			"nonceSize", len(nonce), "preNoncesSize", len(preNonces), "parentHash", hex.EncodeToString(parentHash.Bytes()))
		return nil, staking.ErrWrongFuncParams
	}
	if len(preNonces) > len(validatorList) {
		preNonces = preNonces[len(preNonces)-len(validatorList):]
	}
	return probabilityElection(validatorList, shiftLen, vrf.ProofToHash(nonce), preNonces)
}

func probabilityElection(validatorList staking.ValidatorQueue, shiftLen int, currentNonce []byte, preNonces [][]byte) (staking.ValidatorQueue, error) {
	if len(currentNonce) == 0 || len(preNonces) == 0 || len(validatorList) != len(preNonces) {
		log.Error("Failed to probabilityElection on Election", "validatorListSize", len(validatorList),
			"currentNonceSize", len(currentNonce), "preNoncesSize", len(preNonces), "EpochValidatorNum", xcom.EpochValidatorNum())
		return nil, staking.ErrWrongFuncParams
	}
	sumWeights := new(big.Int)
	svList := make(sortValidatorQueue, 0)
	for _, validator := range validatorList {
		weights, err := validator.GetShares()
		if nil != err {
			return nil, err
		}
		weights.Div(weights, new(big.Int).SetUint64(1e18))
		sumWeights.Add(sumWeights, weights)
		version, err := validator.GetProgramVersion()
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
			v:           validator,
			weights:     int64(weights.Uint64()),
			version:     version,
			blockNumber: blockNumber,
			txIndex:     txIndex,
		}
		svList = append(svList, sv)
	}
	var maxValue float64 = (1 << 256) - 1
	sumWeightsFloat, err := strconv.ParseFloat(sumWeights.Text(10), 64)
	if nil != err {
		return nil, err
	}
	// todo This is an empirical formula, and the follow-up will make a better determination.
	p := float64(xcom.ShiftValidatorNum()) * float64(xcom.ConsValidatorNum()) / sumWeightsFloat
	log.Info("probabilityElection Basic parameter on Election", "validatorListSize", len(validatorList),
		"p", p, "sumWeights", sumWeightsFloat, "shiftValidatorNum", shiftLen, "epochValidatorNum", xcom.EpochValidatorNum())
	for index, sv := range svList {
		resultStr := new(big.Int).Xor(new(big.Int).SetBytes(currentNonce), new(big.Int).SetBytes(preNonces[index])).Text(10)
		target, err := strconv.ParseFloat(resultStr, 64)
		if nil != err {
			return nil, err
		}
		targetP := target / maxValue
		bd := math.NewBinomialDistribution(sv.weights, p)
		x, err := bd.InverseCumulativeProbability(targetP)
		if nil != err {
			return nil, err
		}
		sv.x = x
		log.Debug("calculated probability on Election", "nodeId", sv.v.NodeId.TerminalString(),
			"addr", sv.v.NodeAddress.Hex(), "index", index, "currentNonce",
			hex.EncodeToString(currentNonce), "preNonce", hex.EncodeToString(preNonces[index]),
			"target", target, "targetP", targetP, "weight", sv.weights, "x", x, "version", sv.version,
			"blockNumber", sv.blockNumber, "txIndex", sv.txIndex)
	}
	sort.Sort(svList)
	resultValidatorList := make(staking.ValidatorQueue, shiftLen)
	for index, sv := range svList {
		if index == shiftLen {
			break
		}
		resultValidatorList[index] = sv.v
		log.Debug("sort validator on Election", "nodeId", sv.v.NodeId.TerminalString(),
			"addr", sv.v.NodeAddress.Hex(),
			"index", index, "weight", sv.weights, "x", sv.x, "version", sv.version,
			"blockNumber", sv.blockNumber, "txIndex", sv.txIndex)
	}
	return resultValidatorList, nil
}

/**
Internal expansion function
*/

// previous round validators
func (sk *StakingPlugin) getPreValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValidatorArray, error) {

	var targetIndex *staking.ValArrIndex

	var preTargetNumber uint64
	if blockNumber > xutil.ConsensusSize() {
		preTargetNumber = blockNumber - xutil.ConsensusSize()
	}

	var indexArr staking.ValArrIndexQueue

	if !isCommit {
		indexs, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		indexArr = indexs

		for i, index := range indexs {
			if index.Start <= preTargetNumber && index.End >= preTargetNumber {
				targetIndex = indexs[i]
				break
			}
		}
	} else {
		indexs, err := sk.db.GetRoundValIndexByIrr()
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		indexArr = indexs

		for i, index := range indexs {
			if index.Start <= preTargetNumber && index.End >= preTargetNumber {
				targetIndex = indexs[i]
				break
			}
		}
	}

	if nil == targetIndex {
		log.Error("No Found previous validators index", "isCommit", isCommit,
			"current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		xcom.PrintObjForErr("the round indexs arr is", indexArr)
		return nil, staking.ErrValidatorNoExist
	}

	var queue staking.ValidatorQueue

	if !isCommit {
		arr, err := sk.db.GetRoundValListByBlockHash(blockHash, targetIndex.Start, targetIndex.End)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}
		queue = arr

	} else {
		arr, err := sk.db.GetRoundValListByIrr(targetIndex.Start, targetIndex.End)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}
		queue = arr

	}

	if len(queue) == 0 {
		log.Error("No Found previous validators, the queue length is zero", "isCommit", isCommit, "start", targetIndex.Start,
			"end", targetIndex.End, "current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		return nil, staking.ErrValidatorNoExist
	}

	return &staking.ValidatorArray{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) getCurrValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValidatorArray, error) {

	var targetIndex *staking.ValArrIndex

	var indexArr staking.ValArrIndexQueue

	if !isCommit {
		indexs, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		indexArr = indexs

		for i, index := range indexs {
			if index.Start <= blockNumber && index.End >= blockNumber {
				targetIndex = indexs[i]
				break
			}
		}
	} else {
		indexs, err := sk.db.GetRoundValIndexByIrr()
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		indexArr = indexs

		for i, index := range indexs {
			if index.Start <= blockNumber && index.End >= blockNumber {
				targetIndex = indexs[i]
				break
			}
		}
	}

	if nil == targetIndex {
		log.Error("No Found current validators index", "isCommit", isCommit,
			"current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		xcom.PrintObjForErr("the round indexs arr is", indexArr)
		return nil, staking.ErrValidatorNoExist
	}

	var queue staking.ValidatorQueue

	if !isCommit {
		arr, err := sk.db.GetRoundValListByBlockHash(blockHash, targetIndex.Start, targetIndex.End)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}
		queue = arr

	} else {
		arr, err := sk.db.GetRoundValListByIrr(targetIndex.Start, targetIndex.End)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}
		queue = arr

	}

	if len(queue) == 0 {
		log.Error("No Found current validators, the queue length is zero", "isCommit", isCommit, "start", targetIndex.Start,
			"end", targetIndex.End, "current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		return nil, staking.ErrValidatorNoExist
	}

	return &staking.ValidatorArray{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) getNextValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValidatorArray, error) {

	var targetIndex *staking.ValArrIndex

	var indexArr staking.ValArrIndexQueue

	if !isCommit {
		indexs, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		indexArr = indexs

		for i, index := range indexs {
			if index.Start <= blockNumber && index.End >= blockNumber && i < len(indexs)-1 {
				targetIndex = indexs[i+1]
				break
			}
		}
	} else {
		indexs, err := sk.db.GetRoundValIndexByIrr()
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		indexArr = indexs

		for i, index := range indexs {
			if index.Start <= blockNumber && index.End >= blockNumber && i < len(indexs)-1 {
				targetIndex = indexs[i+1]
				break
			}
		}
	}

	if nil == targetIndex {
		log.Error("No Found next validators index", "isCommit", isCommit,
			"current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		xcom.PrintObjForErr("the round indexs arr is", indexArr)
		return nil, staking.ErrValidatorNoExist
	}

	var queue staking.ValidatorQueue

	if !isCommit {
		arr, err := sk.db.GetRoundValListByBlockHash(blockHash, targetIndex.Start, targetIndex.End)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}
		queue = arr

	} else {
		arr, err := sk.db.GetRoundValListByIrr(targetIndex.Start, targetIndex.End)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}
		queue = arr
	}

	if len(queue) == 0 {
		log.Error("No Found next validators, the queue length is zero", "isCommit", isCommit, "start", targetIndex.Start,
			"end", targetIndex.End, "current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		return nil, staking.ErrValidatorNoExist
	}

	return &staking.ValidatorArray{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) setRoundValListAndIndex(blockNumber uint64, blockHash common.Hash, valArr *staking.ValidatorArray) error {

	log.Debug("Call setRoundValListAndIndex", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"Start", valArr.Start, "End", valArr.End, "arr size", len(valArr.Arr))

	queue, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to setRoundValListAndIndex: Query round valIndex is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"Start", valArr.Start, "End", valArr.End, "err", err)
		return err
	}

	index := &staking.ValArrIndex{
		Start: valArr.Start,
		End:   valArr.End,
	}

	shabby, queue := queue.ConstantAppend(index, RoundValIndexSize)

	// delete the shabby validators
	if nil != shabby {

		log.Debug("Call setRoundValListAndIndex, DelEpochValListByBlockHash",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"shabby.Start", shabby.Start, "shabby.End", shabby.End)

		if err := sk.db.DelRoundValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
			log.Error("Failed to setRoundValListAndIndex: delete shabby validators is failed",
				"shabby start", shabby.Start, "shabby end", shabby.End,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}

	// Store new index Arr
	if err := sk.db.SetRoundValIndex(blockHash, queue); nil != err {
		log.Error("Failed to setRoundValListAndIndex: store round validators new indexArr is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "indexs length", len(queue), "err", err)
		return err
	}

	// Store new round validator Item
	if err := sk.db.SetRoundValList(blockHash, index.Start, index.End, valArr.Arr); nil != err {
		log.Error("Failed to setRoundValListAndIndex: store new round validators is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"start", index.Start, "end", index.End, "val arr length", len(valArr.Arr), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) setRoundValListByIndex(blockNumber uint64, blockHash common.Hash, valArr *staking.ValidatorArray) error {

	log.Debug("Call setRoundValListByIndex", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"Start", valArr.Start, "End", valArr.End, "arr size", len(valArr.Arr))

	queue, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to setRoundValListByIndex: Query round valIndex is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"Start", valArr.Start, "End", valArr.End, "err", err)
		return err
	}

	var hasIndex bool
	// check the Round Index
	for _, indexInfo := range queue {
		if valArr.Start == indexInfo.Start && valArr.End == indexInfo.End {
			hasIndex = true
			break
		}
	}

	if !hasIndex {
		log.Error("No Found current validatorList index", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "input Start", valArr.Start, "input End", valArr.End)
		xcom.PrintObjForErr("the history round indexs arr is", queue)
		return staking.ErrValidatorNoExist
	}

	// Store new round validator Item
	if err := sk.db.SetRoundValList(blockHash, valArr.Start, valArr.End, valArr.Arr); nil != err {
		log.Error("Failed to setRoundValListByIndex: store new round validators is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"start", valArr.Start, "end", valArr.End, "val arr length", len(valArr.Arr), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) getVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValidatorArray, error) {

	var targetIndex *staking.ValArrIndex

	var indexArr staking.ValArrIndexQueue

	if !isCommit {
		indexs, err := sk.db.GetEpochValIndexByBlockHash(blockHash)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		indexArr = indexs

		for i, index := range indexs {
			if index.Start <= blockNumber && index.End >= blockNumber {
				targetIndex = indexs[i]
				break
			}
		}
	} else {
		indexs, err := sk.db.GetEpochValIndexByIrr()
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		indexArr = indexs

		for i, index := range indexs {
			if index.Start <= blockNumber && index.End >= blockNumber {
				targetIndex = indexs[i]
				break
			}
		}
	}

	if nil == targetIndex {
		log.Error("No Found epoch validators index", "isCommit", isCommit,
			"current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		xcom.PrintObjForErr("the epoch indexs arr is", indexArr)
		return nil, staking.ErrValidatorNoExist
	}

	var queue staking.ValidatorQueue

	if !isCommit {
		arr, err := sk.db.GetEpochValListByBlockHash(blockHash, targetIndex.Start, targetIndex.End)

		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}
		queue = arr
	} else {
		arr, err := sk.db.GetEpochValListByIrr(targetIndex.Start, targetIndex.End)

		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}
		queue = arr
	}

	if len(queue) == 0 {
		log.Error("No Found epoch validators, the queue is zero", "isCommit", isCommit, "start", targetIndex.Start,
			"end", targetIndex.End, "current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		return nil, staking.ErrValidatorNoExist
	}

	return &staking.ValidatorArray{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) setVerifierListAndIndex(blockNumber uint64, blockHash common.Hash, valArr *staking.ValidatorArray) error {

	log.Debug("Call setVerifierListAndIndex", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"Start", valArr.Start, "End", valArr.End, "arr size", len(valArr.Arr))

	queue, err := sk.db.GetEpochValIndexByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to setVerifierListAndIndex: Query epoch valIndex is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"Start", valArr.Start, "End", valArr.End, "err", err)
		return err
	}

	index := &staking.ValArrIndex{
		Start: valArr.Start,
		End:   valArr.End,
	}

	shabby, queue := queue.ConstantAppend(index, EpochValIndexSize)

	// delete the shabby validators
	if nil != shabby {
		log.Debug("Call setVerifierListAndIndex, DelEpochValListByBlockHash",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"shabby.Start", shabby.Start, "shabby.End", shabby.End)
		if err := sk.db.DelEpochValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
			log.Error("Failed to setVerifierList: delete shabby validators is failed",
				"shabby start", shabby.Start, "shabby end", shabby.End,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}

	// Store new index Arr
	if err := sk.db.SetEpochValIndex(blockHash, queue); nil != err {
		log.Error("Failed to setVerifierListAndIndex: store epoch validators new indexArr is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "indexs length", len(queue), "err", err)
		return err
	}

	// Store new epoch validator Item
	if err := sk.db.SetEpochValList(blockHash, index.Start, index.End, valArr.Arr); nil != err {
		log.Error("Failed to setVerifierListAndIndex: store new epoch validators is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"start", index.Start, "end", index.End, "val arr length", len(valArr.Arr), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) setVerifierListByIndex(blockNumber uint64, blockHash common.Hash, valArr *staking.ValidatorArray) error {

	log.Debug("Call setVerifierListByIndex", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"Start", valArr.Start, "End", valArr.End, "arr size", len(valArr.Arr))

	queue, err := sk.db.GetEpochValIndexByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to setVerifierListByIndex: Query epoch valIndex is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"Start", valArr.Start, "End", valArr.End, "err", err)
		return err
	}

	var hasIndex bool
	// check the Epoch Index
	for _, indexInfo := range queue {
		if valArr.Start == indexInfo.Start && valArr.End == indexInfo.End {
			hasIndex = true
			break
		}
	}

	if !hasIndex {

		log.Error("No Found current verifierList index", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "input Start", valArr.Start, "input End", valArr.End)
		xcom.PrintObjForErr("the history epoch indexs arr is", queue)

		return staking.ErrValidatorNoExist
	}

	// Store new epoch validator Item
	if err := sk.db.SetEpochValList(blockHash, valArr.Start, valArr.End, valArr.Arr); nil != err {
		log.Error("Failed to setVerifierListByIndex: store new epoch validators is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"start", valArr.Start, "end", valArr.End, "val arr length", len(valArr.Arr), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) addUnStakeItem(state xcom.StateDB, blockNumber uint64, blockHash common.Hash, epoch uint64,
	nodeId discover.NodeID, canAddr common.Address, stakingBlockNum uint64) error {

	endVoteNum, err := gov.GetMaxEndVotingBlock(nodeId, blockHash, state)
	if nil != err {
		return err
	}
	var refundEpoch, maxEndVoteEpoch, targetEpoch uint64
	if endVoteNum != 0 {
		maxEndVoteEpoch = xutil.CalculateEpoch(endVoteNum)
	}

	refundEpoch = xutil.CalculateEpoch(blockNumber) + xcom.UnStakeFreezeRatio()

	if maxEndVoteEpoch <= refundEpoch {
		targetEpoch = refundEpoch
	} else {
		targetEpoch = maxEndVoteEpoch
	}

	log.Info("Call addUnStakeItem, AddUnStakeItemStore start", "current blockNumber", blockNumber,
		"govenance max end vote blokNumber", endVoteNum, "unStakeFreeze Epoch", refundEpoch,
		"govenance max end vote epoch", maxEndVoteEpoch, "unstake item target Epoch", targetEpoch,
		"nodeId", nodeId.String())

	if err := sk.db.AddUnStakeItemStore(blockHash, targetEpoch, canAddr, stakingBlockNum); nil != err {
		return err
	}
	return nil
}

// Record the address of the verification node for each consensus round within a certain block range.
func (sk *StakingPlugin) storeRoundValidatorAddrs(blockHash common.Hash, nextRoundBlockNumber uint64, array staking.ValidatorQueue) error {
	curRound := xutil.CalculateRound(nextRoundBlockNumber)
	curEpoch := xutil.CalculateEpoch(nextRoundBlockNumber)
	validEpoch := uint64(xcom.EvidenceValidEpoch() + 1)
	validRound := xutil.EpochSize() * validEpoch
	if curEpoch > validEpoch {
		invalidRound := curRound - validRound
		key := staking.GetRoundValAddrArrKey(invalidRound)
		if err := sk.db.DelRoundValidatorAddrs(blockHash, key); nil != err {
			log.Error("Failed to DelRoundValidatorAddrs", "blockHash", blockHash.TerminalString(), "nextRoundBlockNumber", nextRoundBlockNumber,
				"curRound", curRound, "curEpoch", curEpoch, "validEpoch", validEpoch, "validRound", validRound, "invalidRound", invalidRound, "key", hex.EncodeToString(key), "err", err)
			return err
		}
		log.Debug("delete RoundValidatorAddrs success", "blockHash", blockHash.TerminalString(), "nextRoundBlockNumber", nextRoundBlockNumber, "invalidRound", invalidRound)
	}
	newKey := staking.GetRoundValAddrArrKey(curRound)
	newValue := make([]common.Address, 0, len(array))
	for _, v := range array {
		newValue = append(newValue, v.NodeAddress)
	}
	if err := sk.db.StoreRoundValidatorAddrs(blockHash, newKey, newValue); nil != err {
		log.Error("Failed to StoreRoundValidatorAddrs", "blockHash", blockHash.TerminalString(), "nextRoundBlockNumber", nextRoundBlockNumber,
			"curRound", curRound, "curEpoch", curEpoch, "validEpoch", validEpoch, "validRound", validRound, "validatorLen", len(array), "newKey", hex.EncodeToString(newKey), "err", err)
		return err
	}
	log.Info("store RoundValidatorAddrs success", "blockHash", blockHash.TerminalString(), "nextRoundBlockNumber", nextRoundBlockNumber,
		"curRound", curRound, "curEpoch", curEpoch, "validEpoch", validEpoch, "validRound", validRound, "validatorLen", len(array))
	return nil
}

func (sk *StakingPlugin) checkRoundValidatorAddr(blockHash common.Hash, targetBlockNumber uint64, addr common.Address) (bool, error) {
	targetRound := xutil.CalculateRound(targetBlockNumber)
	addrList, err := sk.db.LoadRoundValidatorAddrs(blockHash, staking.GetRoundValAddrArrKey(targetRound))
	if nil != err {
		log.Error("Failed to checkRoundValidatorAddr", "blockHash", blockHash.TerminalString(), "targetBlockNumber", targetBlockNumber,
			"addr", addr.Hex(), "targetRound", targetRound, "addrListLen", len(addrList), "err", err)
		return false, err
	}
	if len(addrList) > 0 {
		for _, v := range addrList {
			if bytes.Equal(v.Bytes(), addr.Bytes()) {
				return true, nil
			}
		}
	}
	return false, nil
}

func (sk *StakingPlugin) HasStake(blockHash common.Hash, addr common.Address) (bool, error) {
	return sk.db.HasAccountStakeRc(blockHash, addr)
}

func calcCandidateTotalAmount(can *staking.Candidate) *big.Int {
	release := new(big.Int).Add(can.Released, can.ReleasedHes)
	restrictingPlan := new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
	return new(big.Int).Add(release, restrictingPlan)
}

func calcDelegateTotalAmount(del *staking.Delegation) *big.Int {
	release := new(big.Int).Add(del.Released, del.ReleasedHes)
	restrictingPlan := new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
	return new(big.Int).Add(release, restrictingPlan)
}

func calcRealRefund(realtotal, amount *big.Int) *big.Int {
	refundAmount := new(big.Int).SetInt64(0)
	sub := new(big.Int).Sub(realtotal, amount)
	// When the sub less than threshold
	if !xutil.CheckMinimumThreshold(sub) {
		refundAmount = realtotal
	} else {
		refundAmount = amount
	}
	return refundAmount
}

func buildCanHex(can *staking.Candidate) *staking.CandidateHex {
	return &staking.CandidateHex{
		NodeId:             can.NodeId,
		BlsPubKey:          can.BlsPubKey,
		StakingAddress:     can.StakingAddress,
		BenefitAddress:     can.BenefitAddress,
		StakingTxIndex:     can.StakingTxIndex,
		ProgramVersion:     can.ProgramVersion,
		Status:             can.Status,
		StakingEpoch:       can.StakingEpoch,
		StakingBlockNum:    can.StakingBlockNum,
		Shares:             (*hexutil.Big)(can.Shares),
		Released:           (*hexutil.Big)(can.Released),
		ReleasedHes:        (*hexutil.Big)(can.ReleasedHes),
		RestrictingPlan:    (*hexutil.Big)(can.RestrictingPlan),
		RestrictingPlanHes: (*hexutil.Big)(can.RestrictingPlanHes),
		Description:        can.Description,
	}
}
