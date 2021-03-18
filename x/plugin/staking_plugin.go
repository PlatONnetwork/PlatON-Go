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
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/x/reward"

	"github.com/PlatONnetwork/PlatON-Go/common/math"

	"github.com/PlatONnetwork/PlatON-Go/x/handler"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
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
	FreeVon            = uint16(0)
	RestrictVon        = uint16(1)
	RestrictAndFreeVon = uint16(2)

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

func NewStakingPlugin(db snapshotdb.DB) *StakingPlugin {
	stakePlnOnce.Do(func() {
		log.Info("Init Staking plugin ...")
		stk = &StakingPlugin{
			db: staking.NewStakingDBWithDB(db),
		}
	})
	return stk
}

func (sk *StakingPlugin) SetEventMux(eventMux *event.TypeMux) {
	sk.eventMux = eventMux
}

func (sk *StakingPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error {
	// adjust rewardPer and nextRewardPer
	blockNumber := header.Number.Uint64()
	if xutil.IsBeginOfEpoch(blockNumber) {
		current, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
		if err != nil {
			log.Error("Failed to query current round validators on stakingPlugin BeginBlock",
				"blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "err", err)
			return err
		}
		for _, v := range current.Arr {
			canOld, err := sk.GetCanMutable(blockHash, v.NodeAddress)
			if snapshotdb.NonDbNotFoundErr(err) || canOld.IsEmpty() {
				log.Error("Failed to get candidate info on stakingPlugin BeginBlock", "nodeAddress", v.NodeAddress.String(),
					"blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "err", err)
				return fmt.Errorf("Failed to get candidate info on stakingPlugin BeginBlock, nodeAddress:%s, blockNumber:%d, blockHash:%s", v.NodeAddress.String(), blockNumber, blockHash.TerminalString())
			}
			if canOld.IsInvalid() {
				continue
			}
			var changed bool
			changed = lazyCalcNodeTotalDelegateAmount(xutil.CalculateEpoch(blockNumber), canOld)
			if canOld.RewardPer != canOld.NextRewardPer {
				canOld.RewardPer = canOld.NextRewardPer
				changed = true
			}
			if canOld.CurrentEpochDelegateReward.Cmp(common.Big0) > 0 {
				canOld.CleanCurrentEpochDelegateReward()
				changed = true
			}
			if changed {
				if err = sk.db.SetCanMutableStore(blockHash, v.NodeAddress, canOld); err != nil {
					log.Error("Failed to editCandidate on stakingPlugin BeginBlock", "nodeAddress", v.NodeAddress.String(),
						"blockNumber", blockNumber, "blockHash", blockHash.TerminalString(), "err", err)
					return err
				}
			}

		}
	}
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

func (sk *StakingPlugin) addConsensusNode(nodes staking.ValidatorQueue) {
	for _, node := range nodes {
		if err := sk.eventMux.Post(cbfttypes.AddValidatorEvent{NodeID: node.NodeId}); nil != err {
			log.Error("post AddValidatorEvent failed", "nodeId", node.NodeId.TerminalString(), "err", err)
		}
	}
}

func (sk *StakingPlugin) GetCandidateInfo(blockHash common.Hash, addr common.NodeAddress) (*staking.Candidate, error) {
	return sk.db.GetCandidateStore(blockHash, addr)
}

func (sk *StakingPlugin) GetCanBase(blockHash common.Hash, addr common.NodeAddress) (*staking.CandidateBase, error) {
	return sk.db.GetCanBaseStore(blockHash, addr)
}

func (sk *StakingPlugin) GetCanMutable(blockHash common.Hash, addr common.NodeAddress) (*staking.CandidateMutable, error) {
	return sk.db.GetCanMutableStore(blockHash, addr)
}

func (sk *StakingPlugin) GetCandidateCompactInfo(blockHash common.Hash, blockNumber uint64, addr common.NodeAddress) (*staking.CandidateHex, error) {
	can, err := sk.GetCandidateInfo(blockHash, addr)
	if nil != err {
		return nil, err
	}

	epoch := xutil.CalculateEpoch(blockNumber)
	lazyCalcStakeAmount(epoch, can.CandidateMutable)
	canHex := buildCanHex(can)
	return canHex, nil
}

func (sk *StakingPlugin) GetCandidateInfoByIrr(addr common.NodeAddress) (*staking.Candidate, error) {
	return sk.db.GetCandidateStoreByIrr(addr)
}

func (sk *StakingPlugin) GetCanBaseByIrr(addr common.NodeAddress) (*staking.CandidateBase, error) {
	return sk.db.GetCanBaseStoreByIrr(addr)
}
func (sk *StakingPlugin) GetCanMutableByIrr(addr common.NodeAddress) (*staking.CandidateMutable, error) {
	return sk.db.GetCanMutableStoreByIrr(addr)
}

func (sk *StakingPlugin) CreateCandidate(state xcom.StateDB, blockHash common.Hash, blockNumber, amount *big.Int,
	typ uint16, addr common.NodeAddress, can *staking.Candidate) error {

	if typ == FreeVon { // from account free von

		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to CreateCandidate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
				"stakeAddr", can.StakingAddress, "originVon", origin, "stakingVon", amount)
			return staking.ErrAccountVonNoEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)
		can.ReleasedHes = amount

	} else if typ == RestrictVon { //  from account RestrictingPlan von

		err := rt.AdvanceLockedFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to CreateCandidate on stakingPlugin: call Restricting AdvanceLockedFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
				"stakeAddr", can.StakingAddress, "stakingVon", amount, "err", err)
			return err
		}
		can.RestrictingPlanHes = amount
	} else if typ == RestrictAndFreeVon { //  use Restricting and free von
		restrictingPlanHes, releasedHes, err := rt.MixAdvanceLockedFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to CreateCandidate on stakingPlugin: call Restricting MixAdvanceLockedFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
				"stakeAddr", can.StakingAddress, "stakingVon", amount, "err", err)
			return err
		}
		can.RestrictingPlanHes = restrictingPlanHes
		can.ReleasedHes = releasedHes
	} else {

		log.Error("Failed to CreateCandidate on stakingPlugin", "err", staking.ErrWrongVonOptType,
			"got type", typ, "need type", fmt.Sprintf("%d or %d", FreeVon, RestrictVon))
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
	addr common.NodeAddress, typ uint16) error {

	log.Debug("Call RollBackStaking", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeAddr", addr.Hex())

	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if nil != err {
		return err
	}

	if blockNumber.Uint64() != can.StakingBlockNum {

		log.Error("Failed to RollBackStaking on stakingPlugin: current blockNumber is not equal stakingBlockNumber",
			"blockNumber", blockNumber, "stakingBlockNumber", can.StakingBlockNum)
		return staking.ErrBlockNumberDisordered
	}

	if typ == FreeVon {

		state.AddBalance(can.StakingAddress, can.ReleasedHes)
		state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)

	} else if typ == RestrictVon {

		err := rt.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
		if nil != err {
			log.Error("Failed to RollBackStaking on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
				"stakeAddr", can.StakingAddress, "RollBack stakingVon", can.RestrictingPlanHes, "err", err)
			return err
		}
	} else {

		log.Error("Failed to RollBackStaking on stakingPlugin", "err", staking.ErrWrongVonOptType,
			"got type", typ, "need type", fmt.Sprintf("%d or %d", FreeVon, RestrictVon))
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

func (sk *StakingPlugin) EditCandidate(blockHash common.Hash, blockNumber *big.Int, canAddr common.NodeAddress, can *staking.Candidate) error {
	if err := sk.db.SetCanBaseStore(blockHash, canAddr, can.CandidateBase); nil != err {
		log.Error("Failed to EditCandidate on stakingPlugin: Store CandidateBase info is failed",
			"nodeId", can.NodeId.String(), "blockNumber", blockNumber.Uint64(),
			"blockHash", blockHash.Hex(), "err", err)
		return err
	}
	if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
		log.Error("Failed to EditCandidate on stakingPlugin: Store CandidateMutable info is failed",
			"nodeId", can.NodeId.String(), "blockNumber", blockNumber.Uint64(),
			"blockHash", blockHash.Hex(), "err", err)
		return err
	}
	return nil
}

func (sk *StakingPlugin) IncreaseStaking(state xcom.StateDB, blockHash common.Hash, blockNumber,
	amount *big.Int, typ uint16, canAddr common.NodeAddress, can *staking.Candidate) error {

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can.CandidateMutable)

	if typ == FreeVon {
		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to IncreaseStaking on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
				"nodeId", can.NodeId.String(), "account", can.StakingAddress,
				"originVon", origin, "stakingVon", amount)
			return staking.ErrAccountVonNoEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)
		can.ReleasedHes = new(big.Int).Add(can.ReleasedHes, amount)

	} else if typ == RestrictVon {

		err := rt.AdvanceLockedFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to IncreaseStaking on stakingPlugin: call Restricting AdvanceLockedFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
				"nodeId", can.NodeId.String(), "account", can.StakingAddress, "amount", amount, "err", err)
			return err
		}

		can.RestrictingPlanHes = new(big.Int).Add(can.RestrictingPlanHes, amount)
	} else {

		log.Error("Failed to IncreaseStaking on stakingPlugin", "err", staking.ErrWrongVonOptType,
			"got type", typ, "need type", fmt.Sprintf("%d or %d", FreeVon, RestrictVon))
		return staking.ErrWrongVonOptType
	}

	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to IncreaseStaking on stakingPlugin: Delete Candidate old power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
			"nodeId", can.NodeId.String(), "err", err)
		return err
	}

	can.StakingEpoch = uint32(epoch)
	can.AddShares(amount)

	if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
		log.Error("Failed to IncreaseStaking on stakingPlugin: Store Candidate new power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
			"nodeId", can.NodeId.String(), "err", err)
		return err
	}

	if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
		log.Error("Failed to IncreaseStaking on stakingPlugin: Store CandidateMutable info is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
			"nodeId", can.NodeId.String(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) WithdrewStaking(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	canAddr common.NodeAddress, can *staking.Candidate) error {

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can.CandidateMutable)

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

		if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: Store CandidateMutable info is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return err
		}
	} else {

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
	canAddr common.NodeAddress, can *staking.Candidate) error {

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
				"stakingAddr", can.StakingAddress, "restrictingPlanHes", can.RestrictingPlanHes, "err", err)
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

	can.CleanShares()
	can.Status |= staking.Invalided | staking.Withdrew

	return nil
}

func (sk *StakingPlugin) HandleUnCandidateItem(state xcom.StateDB, blockNumber uint64, blockHash common.Hash, epoch uint64) error {

	unStakeCount, err := sk.db.GetUnStakeCountStore(blockHash, epoch)
	switch {
	case snapshotdb.NonDbNotFoundErr(err):
		return err
	case snapshotdb.IsDbNotFoundErr(err):
		unStakeCount = 0
	}

	if unStakeCount == 0 {
		return nil
	}

	filterAddr := make(map[common.NodeAddress]struct{})

	for index := 1; index <= int(unStakeCount); index++ {

		stakeItem, err := sk.db.GetUnStakeItemStore(blockHash, epoch, uint64(index))
		if nil != err {
			log.Error("Failed to HandleUnCandidateItem: Query the unStakeItem node addr is failed",
				"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		canAddr := stakeItem.NodeAddress

		//log.Debug("Call HandleUnCandidateItem: the candidate Addr",
		//	"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "addr", canAddr.Hex())

		if _, ok := filterAddr[canAddr]; ok {
			if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
				log.Error("Failed to HandleUnCandidateItem: Delete already handle unstakeItem failed",
					"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
				return err
			}
			continue
		}

		can, err := sk.db.GetCandidateStore(blockHash, canAddr)
		if snapshotdb.NonDbNotFoundErr(err) {
			log.Error("Failed to HandleUnCandidateItem: Query candidate failed",
				"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "canAddr", canAddr.Hex(), "err", err)
			return err
		}

		// This should not be nil
		if snapshotdb.IsDbNotFoundErr(err) || can.IsEmpty() {

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

		// The state of the node needs to be restored
		if stakeItem.Recovery {
			// If the node is reported double-signed during the lock-up period，
			// Then you need to enter the double-signed lock-up period after the lock-up period expires and release the pledge after the expiration
			// Otherwise, the state of the node is restored to the normal staking state
			if can.IsDuplicateSign() {

				// Because there is no need to release the staking when the zero-out block is locked, "SubAccountStakeRc" is not executed
				if err := sk.db.SubAccountStakeRc(blockHash, can.StakingAddress); nil != err {
					log.Error("Failed to HandleUnCandidateItem: Sub Account staking Reference Count is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
					return err
				}

				// Lock the node again and will release the staking
				if err := sk.addUnStakeItem(state, blockNumber, blockHash, epoch, can.NodeId, canAddr, can.StakingBlockNum); nil != err {
					log.Error("Failed to SlashCandidates on stakingPlugin: Add UnStakeItemStore failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
					return err
				}
				can.CleanLowRatioStatus()
				if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
					log.Error("Failed to HandleUnCandidateItem on stakingPlugin: Store CandidateMutable info is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
					return err
				}
				log.Debug("Call HandleUnCandidateItem: Node double sign", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
					"status", can.Status, "shares", can.Shares)
			} else {

				can.SetValided()
				if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
					log.Error("Failed to HandleUnCandidateItem on stakingPlugin: Store Candidate power is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
					return err
				}
				if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
					log.Error("Failed to HandleUnCandidateItem on stakingPlugin: Store CandidateMutable info is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
					return err
				}
				log.Debug("Call HandleUnCandidateItem: Node state recovery", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
					"status", can.Status, "shares", can.Shares)
			}
		} else {
			// Second handle balabala ...
			if err := sk.handleUnStake(state, blockNumber, blockHash, epoch, canAddr, can); nil != err {
				return err
			}
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
	addr common.NodeAddress, can *staking.Candidate) error {

	log.Debug("Call handleUnStake", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"epoch", epoch, "nodeId", can.NodeId.String())

	lazyCalcStakeAmount(epoch, can.CandidateMutable)

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

	refundRestrictFn := func(title string, balance *big.Int) (*big.Int, error) {
		if balance.Cmp(common.Big0) > 0 {
			err := rt.ReturnLockFunds(can.StakingAddress, balance, state)
			if nil != err {
				log.Error("Failed to HandleUnCandidateItem on stakingPlugin: call Restricting ReturnLockFunds() is failed",
					title, balance, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
					"stakingAddr", can.StakingAddress, "err", err)
				return new(big.Int).SetInt64(0), err
			}
			return new(big.Int).SetInt64(0), nil
		}
		return balance, nil
	}

	if balance, err := refundRestrictFn("RestrictingPlanHes", can.RestrictingPlanHes); nil != err {
		return err
	} else {
		can.RestrictingPlanHes = balance
	}

	if balance, err := refundRestrictFn("RestrictingPlan", can.RestrictingPlan); nil != err {
		return err
	} else {
		can.RestrictingPlan = balance
	}

	if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
		log.Error("Failed to HandleUnCandidateItem: Delete candidate info failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"nodeId", can.NodeId.String(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) GetDelegatesInfo(blockHash common.Hash, delAddr common.Address) ([]*staking.DelegationInfo, error) {
	return sk.db.GetDelegatesInfo(blockHash, delAddr)
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
			DelegateEpoch:    del.DelegateEpoch,
			Released:         (*hexutil.Big)(del.Released),
			ReleasedHes:      (*hexutil.Big)(del.ReleasedHes),
			RestrictingPlan:  (*hexutil.Big)(del.RestrictingPlan),
			CumulativeIncome: (*hexutil.Big)(del.CumulativeIncome),
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
			CumulativeIncome:   (*hexutil.Big)(del.CumulativeIncome),
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
	delAddr common.Address, del *staking.Delegation, canAddr common.NodeAddress, can *staking.Candidate,
	typ uint16, amount *big.Int, delegateRewardPerList []*reward.DelegateRewardPer) error {

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	rewardsReceive := calcDelegateIncome(epoch, del, delegateRewardPerList)

	if err := UpdateDelegateRewardPer(blockHash, can.NodeId, can.StakingBlockNum, rewardsReceive, sk.db.GetDB()); err != nil {
		return err
	}

	if typ == FreeVon { // from account free von
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

	} else if typ == RestrictVon { //  from account RestrictingPlan von
		err := rt.AdvanceLockedFunds(delAddr, amount, state)
		if nil != err {
			log.Error("Failed to Delegate on stakingPlugin: call Restricting AdvanceLockedFunds() is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch,
				"delAddr", delAddr.String(), "nodeId", can.NodeId.String(), "StakingNum", can.StakingBlockNum,
				"amount", amount, "err", err)
			return err
		}
		del.RestrictingPlanHes = new(big.Int).Add(del.RestrictingPlanHes, amount)

	} else {
		log.Error("Failed to Delegate on stakingPlugin", "err", staking.ErrWrongVonOptType,
			"got type", typ, "need type", fmt.Sprintf("%d or %d", FreeVon, RestrictVon))
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
	can.AddShares(amount)
	// Update total delegate
	lazyCalcNodeTotalDelegateAmount(epoch, can.CandidateMutable)
	can.DelegateTotalHes = new(big.Int).Add(can.DelegateTotalHes, amount)
	can.DelegateEpoch = uint32(epoch)

	// set new power of can
	if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Store Candidate new power is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}

	// update can info about Shares
	if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Store CandidateMutable info is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
		return err
	}
	return nil
}

func (sk *StakingPlugin) WithdrewDelegate(state xcom.StateDB, blockHash common.Hash, blockNumber, amount *big.Int,
	delAddr common.Address, nodeId discover.NodeID, stakingBlockNum uint64, del *staking.Delegation, delegateRewardPerList []*reward.DelegateRewardPer) (*big.Int, error) {
	issueIncome := new(big.Int)
	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to WithdrewDelegate on stakingPlugin: nodeId parse addr failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr,
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
		return nil, err
	}

	can, err := sk.db.GetCandidateStore(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to WithdrewDelegate on stakingPlugin: Query candidate info failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr,
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
		return nil, err
	}

	total := calcDelegateTotalAmount(del)
	// First need to deduct the von that is being refunded
	if total.Cmp(amount) < 0 {
		log.Error("Failed to WithdrewDelegate on stakingPlugin: the amount of valid delegate is not enough",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr,
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "delegate amount", total,
			"withdrew amount", amount)
		return nil, staking.ErrDelegateVonNoEnough
	}

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())
	refundAmount := calcRealRefund(blockNumber.Uint64(), blockHash, total, amount)
	realSub := refundAmount

	rewardsReceive := calcDelegateIncome(epoch, del, delegateRewardPerList)

	if err := UpdateDelegateRewardPer(blockHash, nodeId, stakingBlockNum, rewardsReceive, rm.db); err != nil {
		return nil, err
	}

	// Update total delegate
	if can.IsNotEmpty() {
		lazyCalcNodeTotalDelegateAmount(epoch, can.CandidateMutable)
	}

	del.DelegateEpoch = uint32(epoch)

	switch {
	// Illegal parameter
	case can.IsNotEmpty() && stakingBlockNum > can.StakingBlockNum:
		log.Error("Failed to WithdrewDelegate on stakingPlugin: the stakeBlockNum invalid",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr,
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "fn.stakeBlockNum", stakingBlockNum,
			"can.stakeBlockNum", can.StakingBlockNum)
		return nil, staking.ErrBlockNumberDisordered
	default:
		log.Debug("Call WithdrewDelegate", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.String(), "nodeId", nodeId.String(), "StakingNum", stakingBlockNum,
			"total", total, "amount", amount, "realSub", realSub)

		// handle delegate on Hesitate period
		if refundAmount.Cmp(common.Big0) > 0 {
			rm, rbalance, lbalance, err := rufundDelegateFn(refundAmount, del.ReleasedHes, del.RestrictingPlanHes, delAddr, state)
			if nil != err {
				log.Error("Failed  to WithdrewDelegate, refund the hesitate balance is failed", "blockNumber", blockNumber,
					"blockHash", blockHash.Hex(), "delAddr", delAddr.String(), "nodeId", nodeId.String(), "StakingNum", stakingBlockNum,
					"refund balance", refundAmount, "releaseHes", del.ReleasedHes, "restrictingPlanHes", del.RestrictingPlanHes, "err", err)
				return nil, err
			}
			if can.IsNotEmpty() {
				can.DelegateTotalHes = new(big.Int).Sub(can.DelegateTotalHes, new(big.Int).Sub(refundAmount, rm))
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
				return nil, err
			}
			if can.IsNotEmpty() {
				can.DelegateTotal = new(big.Int).Sub(can.DelegateTotal, new(big.Int).Sub(refundAmount, rm))
			}
			refundAmount, del.Released, del.RestrictingPlan = rm, rbalance, lbalance
		}

		if refundAmount.Cmp(common.Big0) != 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: the withdrew ramain is not zero",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr,
				"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "del balance", total,
				"withdrew balance", amount, "realSub amount", realSub, "withdrew remain", refundAmount)
			return nil, staking.ErrWrongWithdrewDelVonCalc
		}

		// If tatol had full sub,
		// then clean the delegate info
		if total.Cmp(realSub) == 0 {
			// When the entrusted information is deleted, the entrusted proceeds need to be issued automatically
			issueIncome = issueIncome.Add(issueIncome, del.CumulativeIncome)
			if err := rm.ReturnDelegateReward(delAddr, del.CumulativeIncome, state); err != nil {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: return delegate reward is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr,
					"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
				return nil, common.InternalError
			}
			log.Debug("Successful ReturnDelegateReward", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.TerminalString(),
				"delAddr", delAddr, "cumulativeIncome", issueIncome)
			if err := sk.db.DelDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Delete detegate is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr,
					"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
				return nil, err
			}
		} else {
			if err := sk.db.SetDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum, del); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Store detegate is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr,
					"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
				return nil, err
			}
		}
	}

	if can.IsNotEmpty() && stakingBlockNum == can.StakingBlockNum {
		if can.IsValid() {
			if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Delete candidate old power is failed", "blockNumber",
					blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr, "nodeId", nodeId.String(),
					"stakingBlockNum", stakingBlockNum, "err", err)
				return nil, err
			}

			// change candidate shares
			if can.Shares.Cmp(realSub) > 0 {
				can.SubShares(realSub)
			} else {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: the candidate shares is no enough", "blockNumber",
					blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr, "nodeId", nodeId.String(), "stakingBlockNum",
					stakingBlockNum, "can shares", can.Shares, "real withdrew delegate amount", realSub)
				panic("the candidate shares is no enough")
			}

			if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Store candidate old power is failed", "blockNumber",
					blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr, "nodeId", nodeId.String(),
					"stakingBlockNum", stakingBlockNum, "err", err)
				return nil, err
			}
		} else {
			if can.Shares != nil && can.Shares.Cmp(realSub) > 0 {
				can.SubShares(realSub)
			}
		}

		if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: Store CandidateMutable info is failed", "blockNumber",
				blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr, "nodeId", nodeId.String(),
				"stakingBlockNum", stakingBlockNum, "candidateMutable", can.CandidateMutable, "err", err)
			return nil, err
		}
	}
	return issueIncome, nil
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

	oldIndex, err := sk.getVeriferIndex(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to ElectNextVerifierList: Not found the VerifierIndex", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	log.Debug("Call ElectNextVerifierList old verifiers Index", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "index", oldIndex)

	if oldIndex.End != blockNumber {
		log.Error("Failed to ElectNextVerifierList: this blockNumber invalid", "Old Epoch End blockNumber",
			oldIndex.End, "Current blockNumber", blockNumber)
		return staking.ErrBlockNumberDisordered
	}

	// caculate the new epoch start and end
	newVerifierArr := &staking.ValidatorArray{
		Start: oldIndex.End + 1,
		End:   oldIndex.End + xutil.CalcBlocksEachEpoch(),
	}

	currOriginVersion := gov.GetVersionForStaking(blockHash, state)
	currVersion := xutil.CalcVersion(currOriginVersion)

	currentVersion := gov.GetCurrentActiveVersion(state)
	if currentVersion == 0 {
		log.Error("Failed to ElectNextVerifierList, GetCurrentActiveVersion is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString())
		return errors.New("Failed to get CurrentActiveVersion")
	}

	maxvalidators, err := gov.GovernMaxValidators(blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to ElectNextVerifierList: query govern params `maxvalidators` is failed", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	iter := sk.db.IteratorCandidatePowerByBlockHash(blockHash, int(maxvalidators))
	if err := iter.Error(); nil != err {
		log.Error("Failed to ElectNextVerifierList: take iter by candidate power is failed", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}
	defer iter.Release()

	queue := make(staking.ValidatorQueue, 0)

	for iter.Valid(); iter.Next(); {

		addrSuffix := iter.Value()
		canBase, err := sk.db.GetCanBaseStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			log.Error("Failed to ElectNextVerifierList: Query CandidateBase info is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "canAddr", common.BytesToNodeAddress(addrSuffix).Hex(), "err", err)
			if err == snapshotdb.ErrNotFound {
				if err := sk.db.Del(blockHash, iter.Key()); err != nil {
					return err
				}
				// for fix bug Power exist, bug Base is del
				continue
			}
			return err
		}

		if xutil.CalcVersion(canBase.ProgramVersion) < currVersion {
			log.Warn("Warn ElectNextVerifierList: the can ProgramVersion is less than currVersion",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "canVersion",
				"nodeId", canBase.NodeId.String(), "canAddr", common.BytesToNodeAddress(addrSuffix).Hex(),
				canBase.ProgramVersion, "currVersion", currVersion)

			// Low program version cannot be elected for epoch validator
			continue
		}

		addr := common.BytesToNodeAddress(addrSuffix)

		canMutable, err := sk.db.GetCanMutableStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			log.Error("Failed to ElectNextVerifierList: Query CandidateMutable info is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "canAddr", common.BytesToNodeAddress(addrSuffix).Hex(), "err", err)
			return err
		}

		val := &staking.Validator{
			NodeAddress:     addr,
			NodeId:          canBase.NodeId,
			BlsPubKey:       canBase.BlsPubKey,
			ProgramVersion:  canBase.ProgramVersion,
			Shares:          canMutable.Shares,
			StakingBlockNum: canBase.StakingBlockNum,
			StakingTxIndex:  canBase.StakingTxIndex,
			ValidatorTerm:   0,
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

	log.Debug("Call ElectNextVerifierList  Next verifiers", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"list length", len(queue), "list", newVerifierArr)

	return nil
}

func (sk *StakingPlugin) GetVerifierCandidateInfo(blockHash common.Hash, blockNumber uint64) ([]*staking.Candidate, error) {
	verifierList, err := sk.getVerifierList(blockHash, blockNumber, false)
	if nil != err {
		return nil, err
	}

	if blockNumber < verifierList.Start || blockNumber > verifierList.End {
		log.Error("Failed to GetVerifierList", "start", verifierList.Start,
			"end", verifierList.End, "currentNumer", blockNumber)
		return nil, staking.ErrBlockNumberDisordered
	}

	queue := make([]*staking.Candidate, len(verifierList.Arr))

	for i, v := range verifierList.Arr {
		//var can *staking.CandidateBase
		c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
		if nil != err {
			log.Error("Failed to call GetVerifierList, Query Candidate Store info is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
				"canAddr", v.NodeAddress.Hex(), "err", err.Error())
			return nil, err
		}
		queue[i] = c
	}
	return queue, nil
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
		//can, err := sk.GetCandidateInfo(blockHash, v.NodeAddress)

		//var can *staking.CandidateBase
		var can *staking.Candidate
		if !isCommit {
			//c, err := sk.db.GetCanBaseStore(blockHash, v.NodeAddress)
			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetVerifierList, Query Candidate Store info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		} else {
			//c, err := sk.db.GetCanBaseStoreByIrr(v.NodeAddress)
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetVerifierList, Query Candidate Store info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		}

		//shares, _ := new(big.Int).SetString(v.StakingWeight[1], 10)

		valEx := &staking.ValidatorEx{
			NodeId:               can.NodeId,
			BlsPubKey:            can.BlsPubKey,
			StakingAddress:       can.StakingAddress,
			BenefitAddress:       can.BenefitAddress,
			RewardPer:            can.RewardPer,
			NextRewardPer:        can.NextRewardPer,
			RewardPerChangeEpoch: can.RewardPerChangeEpoch,
			StakingTxIndex:       can.StakingTxIndex,
			ProgramVersion:       can.ProgramVersion,
			StakingBlockNum:      can.StakingBlockNum,
			Shares:               (*hexutil.Big)(v.Shares),
			Description:          can.Description,
			ValidatorTerm:        v.ValidatorTerm,
			DelegateTotal:        (*hexutil.Big)(can.DelegateTotal),
			DelegateRewardTotal:  (*hexutil.Big)(can.DelegateRewardTotal),
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

		//var can *staking.CandidateBase
		var can *staking.Candidate
		if !isCommit {
			//c, err := sk.db.GetCanBaseStore(blockHash, v.NodeAddress)
			c, err := sk.db.GetCandidateStore(blockHash, v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetValidatorList, Quey Candidate Store info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		} else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
				log.Error("Failed to call GetValidatorList, Quey Candidate Store info is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(),
					"canAddr", v.NodeAddress.Hex(), "isCommit", isCommit, "err", err.Error())
				return nil, err
			}
			can = c
		}

		valEx := &staking.ValidatorEx{
			NodeId:               can.NodeId,
			BlsPubKey:            can.BlsPubKey,
			StakingAddress:       can.StakingAddress,
			BenefitAddress:       can.BenefitAddress,
			RewardPer:            can.RewardPer,
			NextRewardPer:        can.NextRewardPer,
			RewardPerChangeEpoch: can.RewardPerChangeEpoch,
			StakingTxIndex:       can.StakingTxIndex,
			ProgramVersion:       can.ProgramVersion,
			StakingBlockNum:      can.StakingBlockNum,
			Shares:               (*hexutil.Big)(v.Shares),
			Description:          can.Description,
			ValidatorTerm:        v.ValidatorTerm,
			DelegateTotal:        (*hexutil.Big)(can.DelegateTotal),
			DelegateRewardTotal:  (*hexutil.Big)(can.DelegateRewardTotal),
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

	for iter.Valid(); iter.Next(); {

		addrSuffix := iter.Value()
		can, err := sk.db.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			return nil, err
		}

		lazyCalcStakeAmount(epoch, can.CandidateMutable)
		canHex := buildCanHex(can)
		queue = append(queue, canHex)
	}

	return queue, nil
}

func (sk *StakingPlugin) GetCanBaseList(blockHash common.Hash, blockNumber uint64) (staking.CandidateBaseQueue, error) {

	iter := sk.db.IteratorCandidatePowerByBlockHash(blockHash, 0)
	if err := iter.Error(); nil != err {
		return nil, err
	}
	defer iter.Release()

	queue := make(staking.CandidateBaseQueue, 0)

	for iter.Valid(); iter.Next(); {

		addrSuffix := iter.Value()
		can, err := sk.db.GetCanBaseStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			return nil, err
		}
		queue = append(queue, can)
	}

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

	if can.IsEmpty() || can.IsInvalid() {
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

	blockNumber := header.Number.Uint64()

	// the validators of Current Epoch
	verifiers, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to call Election: Not found current epoch validators", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return staking.ErrValidatorNoExist
	}

	// the validators of Current Round
	curr, err := sk.getCurrValList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to Election: Not found the current round validators", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return staking.ErrValidatorNoExist
	}

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

	invalidCan := make(map[discover.NodeID]struct{})
	removeCans := make(staking.NeedRemoveCans) // the candidates need to remove
	withdrewCans := make(staking.CandidateMap) // the candidates had withdrew

	withdrewQueue := make([]discover.NodeID, 0)
	lowRatioValidAddrs := make([]common.NodeAddress, 0)                 // The addr of candidate that need to clean lowRatio status
	lowRatioValidMap := make(map[common.NodeAddress]*staking.Candidate) // The map collect candidate info that need to clean lowRatio status

	// Query Valid programVersion
	originVersion := gov.GetVersionForStaking(blockHash, state)
	currVersion := xutil.CalcVersion(originVersion)

	// Collecting removed as a result of being slashed
	// That is not withdrew to invalid
	//
	// eg. (lowRatio and must delete) OR (lowRatio and balance no enough) OR duplicateSign
	//
	checkHaveSlash := func(status staking.CandidateStatus) bool {
		return status.IsInvalidLowRatioDel() ||
			status.IsInvalidLowRatio() ||
			status.IsInvalidLowRatioNotEnough() ||
			status.IsInvalidDuplicateSign()
	}

	currentVersion := gov.GetCurrentActiveVersion(state)
	if currentVersion == 0 {
		log.Error("Failed to Election, GetCurrentActiveVersion is failed", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.TerminalString())
		return errors.New("Failed to get CurrentActiveVersion")
	}

	currMap := make(map[discover.NodeID]*big.Int, len(curr.Arr))
	currqueen := make([]*staking.Validator, 0)
	for _, v := range curr.Arr {

		canAddr, _ := xutil.NodeId2Addr(v.NodeId)
		can, err := sk.db.GetCandidateStore(blockHash, canAddr)
		if nil != err {
			log.Error("Failed to Query Candidate Info on Election", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(), "err", err)
			if err == snapshotdb.ErrNotFound {
				// for fix bug Power exist, bug Base is del
				continue
			}
			return err
		}

		var isSlash bool

		if checkHaveSlash(can.Status) {
			removeCans[v.NodeId] = can
			hasSlashLen++
			invalidCan[can.NodeId] = struct{}{}
			isSlash = true
		}

		// Collecting candidate information that active withdrawal
		if can.IsInvalidWithdrew() && !isSlash {
			withdrewCans[v.NodeId] = can
			withdrewQueue = append(withdrewQueue, v.NodeId)
		}

		// valid AND lowRatio status, that candidate need to clean the lowRatio status
		if can.IsValid() && can.IsLowRatio() {
			lowRatioValidAddrs = append(lowRatioValidAddrs, canAddr)
			lowRatioValidMap[canAddr] = can
		}

		// Collect candidate who need to be removed
		// from the validators because the version is too low

		if xutil.CalcVersion(can.ProgramVersion) < currVersion {
			removeCans[v.NodeId] = can
			invalidCan[can.NodeId] = struct{}{}
			needRMLowVersionLen++
		}

		currMap[v.NodeId] = v.Shares
		currqueen = append(currqueen, v)
	}

	// Exclude the current consensus round validators from the validators of the Epoch
	diffQueue := make(staking.ValidatorQueue, 0)
	for _, v := range verifiers.Arr {

		if _, ok := withdrewCans[v.NodeId]; ok {
			delete(withdrewCans, v.NodeId)
		}

		if _, ok := currMap[v.NodeId]; ok {
			currMap[v.NodeId] = new(big.Int).Set(v.Shares)
			continue
		}

		addr, _ := xutil.NodeId2Addr(v.NodeId)
		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			log.Error("Failed to Get Candidate on Election", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(), "err", err)
			if err == snapshotdb.ErrNotFound {
				// for fix bug Power exist, bug Base is del
				continue
			}
			return err
		}

		// Jump the slashed candidate
		if checkHaveSlash(can.Status) {
			continue
		}

		// Ignore the low version
		if xutil.CalcVersion(can.ProgramVersion) < currVersion {
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
	for _, nodeID := range withdrewQueue {
		invalidCan[nodeID] = struct{}{}
	}

	// some validators that meets the following conditions must be replaced first.
	// eg.
	// 1. Be reported as evil
	// 2. The package ratio is low and the remaining deposit balance is less than the minimum staking threshold
	// 3. The version number in the validator's real-time details
	// 	  is lower than the version of the governance module on the current chain.
	// 4. withdrew staking and not in the current epoch validator list
	//
	//invalidLen = hasSlashLen + needRMwithdrewLen + needRMLowVersionLen
	invalidLen = len(invalidCan)

	shuffle := func(invalidLen int, currQueue, vrfQueue staking.ValidatorQueue, blockNumber uint64, parentHash common.Hash) (staking.ValidatorQueue, error) {

		// increase term and use new shares  one by one
		for i, v := range currQueue {
			v.ValidatorTerm++
			v.Shares = currMap[v.NodeId]
			currQueue[i] = v
		}

		// sort the validator by del rule
		currQueue.ValidatorSort(removeCans, staking.CompareForDel)
		// Increase term of validator
		copyCurrQueue := make(staking.ValidatorQueue, len(currQueue)-invalidLen)
		// Remove the invalid validators
		copy(copyCurrQueue, currQueue[invalidLen:])
		return shuffleQueue(copyCurrQueue, vrfQueue, blockNumber, parentHash)
	}

	var vrfQueue staking.ValidatorQueue
	var vrfLen int
	if len(diffQueue) > int(xcom.MaxConsensusVals()) {
		vrfLen = int(xcom.MaxConsensusVals())
	} else {
		vrfLen = len(diffQueue)
	}

	if vrfLen != 0 {
		if queue, err := vrfElection(diffQueue, vrfLen, header.Nonce.Bytes(), header.ParentHash, blockNumber, currentVersion); nil != err {
			log.Error("Failed to VrfElection on Election",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		} else {
			vrfQueue = queue
		}
	}

	log.Debug("Call Election, statistics need to remove node num",
		"has slash count", hasSlashLen, "withdrew and need remove count",
		needRMwithdrewLen, "low version need remove count", needRMLowVersionLen,
		"total remove count", invalidLen, "remove map size", len(removeCans),
		"current validators Size", len(curr.Arr), "MaxConsensusVals", xcom.MaxConsensusVals(),
		"ShiftValidatorNum", xcom.ShiftValidatorNum(), "diffQueueLen", len(diffQueue),
		"vrfQueueLen", len(vrfQueue))

	nextQueue, err := shuffle(invalidLen, currqueen, vrfQueue, blockNumber, header.ParentHash)
	if nil != err {
		return err
	}

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

	// update candidate status
	// Must sort
	for _, canAddr := range lowRatioValidAddrs {

		can := lowRatioValidMap[canAddr]
		// clean the low package ratio status
		can.CleanLowRatioStatus()

		addr, _ := xutil.NodeId2Addr(can.NodeId)
		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Failed to Store Candidate on Election", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return err
		}
	}

	if err := sk.storeRoundValidatorAddrs(blockNumber, blockHash, start, nextQueue); nil != err {
		log.Error("Failed to storeRoundValidatorAddrs on Election", "blockNumber", blockNumber,
			"blockHash", blockHash.TerminalString(), "err", err)
		return err
	}

	log.Debug("Call Election Next validators", "blockNumber", header.Number.Uint64(), "blockHash", blockHash.Hex(),
		"list length", len(next.Arr), "list", next)

	return nil
}

func shuffleQueue(remainCurrQueue, vrfQueue staking.ValidatorQueue, blockNumber uint64, parentHash common.Hash) (staking.ValidatorQueue, error) {

	remainLen := len(remainCurrQueue)
	totalQueue := append(remainCurrQueue, vrfQueue...)

	for remainLen > int(xcom.MaxConsensusVals()-xcom.ShiftValidatorNum()) && len(totalQueue) > int(xcom.MaxConsensusVals()) {
		totalQueue = totalQueue[1:]
		remainLen--
	}

	if len(totalQueue) > int(xcom.MaxConsensusVals()) {
		totalQueue = totalQueue[:xcom.MaxConsensusVals()]
	}

	next := make(staking.ValidatorQueue, len(totalQueue))

	copy(next, totalQueue)

	// Divide all consensus nodes into two groups, the front and back positions of each group are not changed,
	// but random ordering is performed in each group
	// The first group: the first f nodes
	// The second group: the last 2f + 1 nodes
	next, err := randomOrderValidatorQueue(blockNumber, parentHash, next)
	if nil != err {
		return nil, err
	}
	return next, nil
}

type randomOrderValidator struct {
	validator *staking.Validator
	value     *big.Int
}
type randomOrderValidatorList []*randomOrderValidator

func (r randomOrderValidatorList) Len() int {
	return len(r)
}

func (r randomOrderValidatorList) Less(i, j int) bool {
	return r[i].value.Cmp(r[j].value) > 0
}

func (r randomOrderValidatorList) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

// Randomly sort nodes
func randomOrderValidatorQueue(blockNumber uint64, parentHash common.Hash, queue staking.ValidatorQueue) (staking.ValidatorQueue, error) {
	preNonces, err := handler.GetVrfHandlerInstance().Load(parentHash)
	if nil != err {
		return nil, err
	}
	if len(preNonces) < len(queue) {
		log.Error("Failed to randomOrderValidatorQueue on Election", "blockNumber", blockNumber, "validatorListSize", len(queue),
			"preNoncesSize", len(preNonces), "parentHash", parentHash.TerminalString())
		return nil, staking.ErrWrongFuncParams
	}
	if len(preNonces) > len(queue) {
		preNonces = preNonces[len(preNonces)-len(queue):]
	}

	if len(queue) <= int(xcom.ShiftValidatorNum()) {
		return queue, nil
	}

	orderList := make(randomOrderValidatorList, len(queue))
	for i, v := range queue {
		value := new(big.Int).Xor(new(big.Int).SetBytes(v.NodeAddress.Bytes()), new(big.Int).SetBytes(preNonces[i][:common.AddressLength]))
		orderList[i] = &randomOrderValidator{
			validator: v,
			value:     value,
		}
		log.Debug("Call randomOrderValidatorQueue xor", "nodeId", v.NodeId.TerminalString(), "nodeAddress", v.NodeAddress.Hex(), "nonce", hexutil.Encode(preNonces[i]), "xorValue", value)
	}

	frontPart := orderList[:xcom.ShiftValidatorNum()]
	backPart := orderList[xcom.ShiftValidatorNum():]

	sort.Sort(frontPart)
	sort.Sort(backPart)

	orderList = make(randomOrderValidatorList, 0)
	orderList = append(orderList, frontPart...)
	orderList = append(orderList, backPart...)

	resultQueue := make(staking.ValidatorQueue, len(orderList))
	for i, v := range orderList {
		resultQueue[i] = v.validator
	}
	log.Debug("Call randomOrderValidatorQueue success", "blockNumber", blockNumber, "parentHash", parentHash.TerminalString(), "resultQueueSize", len(resultQueue))
	return resultQueue, nil
}

// NotifyPunishedVerifiers
func (sk *StakingPlugin) SlashCandidates(state xcom.StateDB, blockHash common.Hash, blockNumber uint64, queue ...*staking.SlashNodeItem) error {

	// Nodes that need to be deleted from the candidate list
	// Keep governance votes that have been voted
	invalidNodeIdMap := make(map[discover.NodeID]struct{}, 0)
	// Need to remove eligibility to govern voting
	invalidRemoveGovNodeIdMap := make(map[discover.NodeID]struct{}, 0)

	for _, slashItem := range queue {
		needRemove, err := sk.toSlash(state, blockNumber, blockHash, slashItem)
		if nil != err {
			return err
		}
		if needRemove {
			invalidNodeIdMap[slashItem.NodeId] = struct{}{}
			if slashItem.SlashType != staking.LowRatio {
				invalidRemoveGovNodeIdMap[slashItem.NodeId] = struct{}{}
			}
		}
	}

	if len(invalidNodeIdMap) != 0 {
		// remove the validator from epoch verifierList
		if err := sk.removeFromVerifiers(blockNumber, blockHash, invalidNodeIdMap); nil != err {
			return err
		}
		if len(invalidRemoveGovNodeIdMap) > 0 {
			// notify gov to do somethings
			if err := gov.NotifyPunishedVerifiers(blockHash, invalidRemoveGovNodeIdMap, state); nil != err {
				log.Error("Failed to SlashCandidates: call NotifyPunishedVerifiers of gov is failed", "blockNumber", blockNumber,
					"blockHash", blockHash.Hex(), "invalidNodeId Size", len(invalidRemoveGovNodeIdMap), "err", err)
				return err
			}
		}
	}

	return nil
}

func (sk *StakingPlugin) toSlash(state xcom.StateDB, blockNumber uint64, blockHash common.Hash, slashItem *staking.SlashNodeItem) (bool, error) {

	log.Debug("Call SlashCandidates: call toSlash", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"nodeId", slashItem.NodeId.String(), "amount", slashItem.Amount, "slashType", slashItem.SlashType,
		"benefitAddr", slashItem.BenefitAddr)

	var needRemove bool

	// check slash type is right
	slashTypeIsWrong := func() bool {
		return !slashItem.SlashType.IsLowRatio() &&
			!slashItem.SlashType.IsLowRatioDel() &&
			!slashItem.SlashType.IsDuplicateSign()
	}
	if slashTypeIsWrong() {
		log.Error("Failed to SlashCandidates: the slashType is wrong", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "slashType", slashItem.SlashType, "benefitAddr", slashItem.BenefitAddr)
		return needRemove, staking.ErrWrongSlashType
	}

	canAddr, _ := xutil.NodeId2Addr(slashItem.NodeId)
	can, err := sk.db.GetCandidateStore(blockHash, canAddr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to SlashCandidates: Query can is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
		return needRemove, err
	}

	if can.IsEmpty() {
		log.Error("Failed to SlashCandidates: the can is empty", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String())
		return needRemove, staking.ErrCanNoExist
	}

	epoch := xutil.CalculateEpoch(blockNumber)
	lazyCalcStakeAmount(epoch, can.CandidateMutable)

	// Balance that can only be effective for Slash
	total := new(big.Int).Add(can.Released, can.RestrictingPlan)

	if slashItem.Amount != nil && total.Cmp(slashItem.Amount) < 0 {
		log.Error("Warned to SlashCandidates: the candidate total staking amount is not enough",
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

	// If the node is already in a state of low block rate,
	// it will not punish the behavior of low block rate again
	// If the penalty is imposed again,
	// the deposit may be lower than the minimum deposit and may be forced to release the staking during the lock-in period
	if can.IsLowRatio() && slashItem.SlashType.IsLowRatio() {
		log.Info("Call SlashCandidates: node has already been punished", "nodeId", slashItem.NodeId.String(), "nodeStatus", can.Status,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "slashType", slashItem.SlashType, "slashAmount", slashItem.Amount)
	} else {
		slashBalance := slashItem.Amount
		// slash the balance
		if slashBalance.Cmp(common.Big0) > 0 && can.Released.Cmp(common.Big0) > 0 {
			val, rval, err := slashBalanceFn(slashBalance, can.Released, false, slashItem.SlashType,
				slashItem.BenefitAddr, can.StakingAddress, state)
			if nil != err {
				log.Error("Failed to SlashCandidates: slash Released", "slashed amount", slashBalance,
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
				return needRemove, err
			}
			slashBalance, can.Released = val, rval
		}
		if slashBalance.Cmp(common.Big0) > 0 && can.RestrictingPlan.Cmp(common.Big0) > 0 {
			val, rval, err := slashBalanceFn(slashBalance, can.RestrictingPlan, true, slashItem.SlashType,
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
			return (can.IsInvalidLowRatioNotEnough() ||
				can.IsInvalidLowRatioDel() ||
				can.IsInvalidDuplicateSign() ||
				can.IsInvalidWithdrew())
		}

		// If the shares is zero, don't need to sub shares
		if !sharesHaveBeenClean() {

			// first slash and no withdrew
			// sub Shares to effect power
			if can.Shares.Cmp(slashItem.Amount) >= 0 {
				can.SubShares(slashItem.Amount)
			} else {
				log.Error("Failed to SlashCandidates: the candidate shares is no enough", "slashType", slashItem.SlashType,
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "candidate shares",
					can.Shares, "slash amount", slashItem.Amount)
				panic("the candidate shares is no enough")
			}
		}
	}

	// need invalid candidate status
	// need remove from verifierList
	needInvalid, needRemove, needReturnHes, changeStatus := handleSlashTypeFn(blockNumber, blockHash, slashItem.SlashType, calcCandidateTotalAmount(can))

	log.Debug("Call SlashCandidates: the status", "needInvalid", needInvalid,
		"needRemove", needRemove, "needReturnHes", needReturnHes, "current can.Status", can.Status, "need to superpose status", changeStatus)

	if needRemove {
		if err := rm.ReturnDelegateReward(can.BenefitAddress, can.CurrentEpochDelegateReward, state); err != nil {
			log.Error("Call SlashCandidates:return delegateReward", "err", err)
		}
		can.CleanCurrentEpochDelegateReward()
	}

	// Only when the staking is released, the staking-related information needs to be emptied.
	// When penalizing the low block rate first, and then report double signing, the pledged deposit in the period of hesitation should be returned
	if needReturnHes {
		// Return the pledged deposit during the hesitation period
		if can.ReleasedHes.Cmp(common.Big0) > 0 {
			state.AddBalance(can.StakingAddress, can.ReleasedHes)
			state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)
			can.ReleasedHes = new(big.Int).SetInt64(0)
		}
		if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {
			err := rt.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
			if nil != err {
				log.Error("Failed to SlashCandidates on stakingPlugin: call Restricting ReturnLockFunds() is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "stakingAddr", can.StakingAddress,
					"restrictingPlanHes", can.RestrictingPlanHes, "err", err)
				return needRemove, err
			}
			can.RestrictingPlanHes = new(big.Int).SetInt64(0)
		}

		//because of deleted candidate info ,clean Shares
		can.CleanShares()
	}

	if needInvalid && can.IsValid() {
		can.AppendStatus(changeStatus)
		// Only when the staking is released, the staking-related information needs to be emptied.
		if needReturnHes {
			// need to sub account rc
			// Only need to be executed if the pledge is released
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
		} else {
			// Add a freeze message, after the freeze is over, it can return to normal state
			if err := sk.addRecoveryUnStakeItem(blockNumber, blockHash, can.NodeId, canAddr, can.StakingBlockNum); nil != err {
				log.Error("Failed to SlashCandidates on stakingPlugin: addRecoveryUnStakeItem failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
				return needRemove, err
			}
		}

		if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
			log.Error("Failed to SlashCandidates on stakingPlugin: Store CandidateMutable info is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return needRemove, err
		}

	} else if !needInvalid && can.IsValid() {

		// update the candidate power, If do not need to delete power (the candidate status still be valid)
		if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to SlashCandidates: Store candidate power is failed", "slashType", slashItem.SlashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
			return needRemove, err
		}

		can.AppendStatus(changeStatus)
		if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
			log.Error("Failed to SlashCandidates: Store CandidateMutable is failed", "slashType", slashItem.SlashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", slashItem.NodeId.String(), "err", err)
			return needRemove, err
		}

	} else {
		can.AppendStatus(changeStatus)
		if err := sk.db.SetCanMutableStore(blockHash, canAddr, can.CandidateMutable); nil != err {
			log.Error("Failed to SlashCandidates: Store CandidateMutable is failed", "slashType", slashItem.SlashType,
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

func handleSlashTypeFn(blockNumber uint64, blockHash common.Hash, slashType staking.CandidateStatus, remain *big.Int) (bool, bool, bool, staking.CandidateStatus) {

	var needInvalid, needRemove, needReturnHes bool // need invalid candidate status And need remove from verifierList,Refund of pledged deposits during hesitation period
	var changeStatus staking.CandidateStatus        // need to add this status

	switch slashType {
	case staking.LowRatio:
		if ok, _ := CheckStakeThreshold(blockNumber, blockHash, remain); !ok {
			changeStatus |= staking.NotEnough
			needReturnHes = true
		}
		changeStatus |= staking.Invalided
		needInvalid = true
		needRemove = true
	case staking.LowRatioDel:
		changeStatus |= staking.Invalided
		needInvalid = true
		needRemove = true
		needReturnHes = true
	case staking.DuplicateSign:
		changeStatus |= staking.Invalided
		needInvalid = true
		needRemove = true
		needReturnHes = true
	}
	changeStatus |= slashType

	return needInvalid, needRemove, needReturnHes, changeStatus
}

func slashBalanceFn(slashAmount, canBalance *big.Int, isNotify bool,
	slashType staking.CandidateStatus, benefitAddr, stakingAddr common.Address, state xcom.StateDB) (*big.Int, *big.Int, error) {

	// check zero value
	// If there is a zero value, no logic is done.
	if canBalance.Cmp(common.Big0) == 0 || slashAmount.Cmp(common.Big0) == 0 {
		return slashAmount, canBalance, nil
	}

	slashAmountTmp := new(big.Int).SetInt64(0)
	balanceTmp := new(big.Int).SetInt64(0)

	if slashAmount.Cmp(canBalance) >= 0 {

		state.SubBalance(vm.StakingContractAddr, canBalance)

		if slashType.IsDuplicateSign() {
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
		if slashType.IsDuplicateSign() {
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

	log.Info("Call ProposalPassedNotify to promote candidate programVersion", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "version", programVersion, "nodeIdQueueSize", len(nodeIds))

	for _, nodeId := range nodeIds {
		log.Info("Call ProposalPassedNotify itr nodeId", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeid", nodeId)
		addr, _ := xutil.NodeId2Addr(nodeId)
		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if snapshotdb.NonDbNotFoundErr(err) {
			log.Error("Failed to ProposalPassedNotify: Query Candidate is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		if snapshotdb.IsDbNotFoundErr(err) || can.IsEmpty() {
			log.Error("Failed to ProposalPassedNotify: Promote candidate programVersion failed, the can is empty",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
			continue
		}

		if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
			log.Error("Failed to ProposalPassedNotify: Delete Candidate old power is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}
		can.ProgramVersion = programVersion
		//Store full version
		if err := sk.db.SetCanBaseStore(blockHash, addr, can.CandidateBase); nil != err {
			log.Error("Failed to ProposalPassedNotify: Store CandidateBase info is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		if can.IsInvalid() {
			log.Warn(" can status is invalid,no need set can power", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "status", can.Status)
			continue
		}
		if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
			log.Error("Failed to ProposalPassedNotify: Store Candidate new power is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

	}

	return nil
}

func (sk *StakingPlugin) DeclarePromoteNotify(blockHash common.Hash, blockNumber uint64, nodeId discover.NodeID,
	programVersion uint32) error {

	log.Info("Call DeclarePromoteNotify to promote candidate programVersion", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "real version", programVersion, "calc version", xutil.CalcVersion(programVersion), "nodeId", nodeId.String())

	addr, _ := xutil.NodeId2Addr(nodeId)
	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if snapshotdb.NonDbNotFoundErr(err) {
		log.Error("Failed to DeclarePromoteNotify: Query Candidate is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	if snapshotdb.IsDbNotFoundErr(err) || can.IsEmpty() {
		log.Error("Failed to DeclarePromoteNotify: Promote candidate programVersion failed, the can is empty",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(),
			"version", programVersion)
		return nil
	}

	if can.IsInvalid() {
		log.Warn(" can status is invalid,no need set can power", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "status", can.Status)
		return nil
	}

	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to DeclarePromoteNotify: Delete Candidate old power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	can.ProgramVersion = programVersion

	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to DeclarePromoteNotify: Store Candidate new power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}
	//Store full version
	if err := sk.db.SetCanBaseStore(blockHash, addr, can.CandidateBase); nil != err {
		log.Error("Failed to DeclarePromoteNotify: Store CandidateBase info is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) GetLastNumber(blockNumber uint64) uint64 {

	valIndex, err := sk.getCurrValIndex(common.ZeroHash, blockNumber, QueryStartIrr)
	if nil != err {
		log.Error("Failed to GetLastNumber", "blockNumber", blockNumber, "err", err)
		return 0
	}

	if nil == err && nil != valIndex {
		return valIndex.End
	}
	return 0
}

func (sk *StakingPlugin) GetValidator(blockNumber uint64) (*cbfttypes.Validators, error) {

	valArr, err := sk.getCurrValList(common.ZeroHash, blockNumber, QueryStartIrr)
	if snapshotdb.NonDbNotFoundErr(err) {
		return nil, err
	}

	if nil == err && nil != valArr {
		return buildCbftValidators(valArr.Start, valArr.Arr), nil
	}
	return nil, fmt.Errorf("Not Found Validators by blockNumber: %d", blockNumber)
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
	return isCandidate
}

func buildCbftValidators(start uint64, arr staking.ValidatorQueue) *cbfttypes.Validators {
	valMap := make(cbfttypes.ValidateNodeMap, len(arr))

	for i, v := range arr {

		pubKey, _ := v.NodeId.Pubkey()
		blsPk, _ := v.BlsPubKey.ParseBlsPubKey()

		vn := &cbfttypes.ValidateNode{
			Index:     uint32(i),
			Address:   v.NodeAddress,
			PubKey:    pubKey,
			NodeID:    v.NodeId,
			BlsPubKey: blsPk,
		}

		valMap[v.NodeId] = vn
	}

	res := &cbfttypes.Validators{
		Nodes:            valMap,
		ValidBlockNumber: start,
	}
	return res
}

func lazyCalcStakeAmount(epoch uint64, can *staking.CandidateMutable) {
	if can.IsEmpty() {
		return
	}

	changeAmountEpoch := can.StakingEpoch

	sub := epoch - uint64(changeAmountEpoch)

	log.Debug("lazyCalcStakeAmount before", "current epoch", epoch, "canMutable", can)

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

	log.Debug("lazyCalcStakeAmount end", "current epoch", epoch, "canMutable", can)

}

// The total delegate amount of the compute node
func lazyCalcNodeTotalDelegateAmount(epoch uint64, can *staking.CandidateMutable) bool {
	if can.IsEmpty() {
		return false
	}
	changeAmountEpoch := can.DelegateEpoch
	sub := epoch - uint64(changeAmountEpoch)
	log.Debug("lazyCalcNodeTotalDelegateAmount before", "current epoch", epoch, "canMutable", can)

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio() {
		return false
	}
	if can.DelegateTotalHes.Cmp(common.Big0) > 0 {
		can.DelegateTotal = new(big.Int).Add(can.DelegateTotal, can.DelegateTotalHes)
		can.DelegateTotalHes = new(big.Int).SetInt64(0)
		return true
	}
	return false
}

func lazyCalcDelegateAmount(epoch uint64, del *staking.Delegation) {
	if del.IsEmpty() {
		return
	}

	// When the first time, there was no previous changeAmountEpoch
	if del.DelegateEpoch == 0 {
		return
	}

	changeAmountEpoch := del.DelegateEpoch

	sub := epoch - uint64(changeAmountEpoch)

	log.Debug("lazyCalcDelegateAmount before", "epoch", epoch, "del", del)

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

	log.Debug("lazyCalcDelegateAmount end", "epoch", epoch, "del", del)
}

// Calculating Total Entrusted Income
func calcDelegateIncome(epoch uint64, del *staking.Delegation, per []*reward.DelegateRewardPer) []reward.DelegateRewardReceipt {
	if del.IsEmpty() {
		return nil
	}
	// Triggered again in the same cycle, no need to calculate revenue
	if uint64(del.DelegateEpoch) == epoch {
		return nil
	}
	// When the settlement period when the first delegation is the same as the settlement period when the node is staking,
	// and when the delegation is performed again in the next settlement cycle,
	// the "per" value of the node is not available at this time, so you need to directly lazy
	if len(per) == 0 {
		lazyCalcDelegateAmount(epoch, del)
		return nil
	}

	delegateRewardReceives := make([]reward.DelegateRewardReceipt, 0)
	// When the settlement period at the first delegation is the same as the settlement period at the node's staking,
	// , And when a second delegation is made after multiple billing cycles,
	// For example: the node staking when the settlement period = 1, and the delegation is also in this settlement period. At this time,
	// the node's "per" value only starts when the settlement period = 2; now it is the settlement period = 5, and this time it needs to be calculated When entrusting the income,
	// you need to convert the entrustment from the hesitation period to the lock-in period before calculating.
	if per[0].Epoch > uint64(del.DelegateEpoch) {
		lazyCalcDelegateAmount(epoch, del)
	}
	totalReleased := new(big.Int).Add(del.Released, del.RestrictingPlan)
	for i, rewardPer := range per {
		if totalReleased.Cmp(common.Big0) > 0 {
			if nil == del.CumulativeIncome {
				del.CumulativeIncome = new(big.Int)
			}
			delegateRewardReceive := reward.DelegateRewardReceipt{
				Epoch:    rewardPer.Epoch,
				Delegate: new(big.Int).Set(totalReleased),
			}
			delegateRewardReceives = append(delegateRewardReceives, delegateRewardReceive)
			del.CumulativeIncome = new(big.Int).Add(del.CumulativeIncome, rewardPer.CalDelegateReward(totalReleased))
		}
		if i == 0 {
			lazyCalcDelegateAmount(epoch, del)
			totalReleased = new(big.Int).Add(del.Released, del.RestrictingPlan)
		}
	}
	log.Debug("Call calcDelegateIncome end", "currEpoch", epoch, "perLen", len(per), "delegateRewardReceivesLen", len(delegateRewardReceives),
		"totalDelegate", totalReleased, "totalHes", new(big.Int).Add(del.ReleasedHes, del.RestrictingPlanHes), "income", del.CumulativeIncome)
	return delegateRewardReceives
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
	if xutil.CalcVersion(svs[i].version) == xutil.CalcVersion(svs[j].version) {
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
		return xutil.CalcVersion(svs[i].version) > xutil.CalcVersion(svs[j].version)
	}
}

func (svs sortValidatorQueue) Swap(i, j int) {
	svs[i], svs[j] = svs[j], svs[i]
}

// Elected verifier by vrf random election
// validatorList：Waiting for the elected node
// nonce：Vrf proof of the current block
// parentHash：Parent block hash
func vrfElection(validatorList staking.ValidatorQueue, shiftLen int, nonce []byte, parentHash common.Hash, blockNumber uint64, currentVersion uint32) (staking.ValidatorQueue, error) {
	preNonces, err := handler.GetVrfHandlerInstance().Load(parentHash)
	if nil != err {
		return nil, err
	}
	if len(preNonces) < len(validatorList) {
		log.Error("Failed to vrfElection on Election", "blockNumber", blockNumber, "validatorListSize", len(validatorList),
			"nonceSize", len(nonce), "preNoncesSize", len(preNonces), "parentHash", hex.EncodeToString(parentHash.Bytes()))
		return nil, staking.ErrWrongFuncParams
	}
	if len(preNonces) > len(validatorList) {
		preNonces = preNonces[len(preNonces)-len(validatorList):]
	}
	return probabilityElection(validatorList, shiftLen, vrf.ProofToHash(nonce), preNonces, blockNumber, currentVersion)
}

func probabilityElection(validatorList staking.ValidatorQueue, shiftLen int, currentNonce []byte, preNonces [][]byte, blockNumber uint64, currentVersion uint32) (staking.ValidatorQueue, error) {
	if len(currentNonce) == 0 || len(preNonces) == 0 || len(validatorList) != len(preNonces) {
		log.Error("Failed to probabilityElection", "blockNumber", blockNumber, "currentVersion", currentVersion, "validators Size", len(validatorList),
			"currentNonceSize", len(currentNonce), "preNoncesSize", len(preNonces))
		return nil, staking.ErrWrongFuncParams
	}
	totalWeights := new(big.Int)
	totalSqrtWeights := new(big.Int)
	svList := make(sortValidatorQueue, 0)
	for _, val := range validatorList {

		weights := new(big.Int).Div(val.Shares, new(big.Int).SetUint64(1e18))
		totalWeights.Add(totalWeights, weights)
		weights = new(big.Int).Sqrt(weights)
		totalSqrtWeights.Add(totalSqrtWeights, weights)

		sv := &sortValidator{
			v:           val,
			weights:     int64(weights.Uint64()),
			version:     val.ProgramVersion,
			blockNumber: val.StakingBlockNum,
			txIndex:     val.StakingTxIndex,
		}
		svList = append(svList, sv)
	}
	var maxValue float64 = (1 << 256) - 1
	totalWeightsFloat, err := strconv.ParseFloat(totalWeights.Text(10), 64)
	if nil != err {
		return nil, err
	}
	totalSqrtWeightsFloat, err := strconv.ParseFloat(totalSqrtWeights.Text(10), 64)
	if nil != err {
		return nil, err
	}

	p := xcom.CalcP(totalWeightsFloat, totalSqrtWeightsFloat)

	log.Debug("Call probabilityElection Basic parameter on Election", "blockNumber", blockNumber, "currentVersion", currentVersion, "validatorListSize", len(validatorList),
		"p", p, "totalWeights", totalWeightsFloat, "totalSqrtWeightsFloat", totalSqrtWeightsFloat, "shiftValidatorNum", shiftLen)

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

		log.Debug("Call probabilityElection, calculated probability on Election", "nodeId", sv.v.NodeId.TerminalString(),
			"addr", sv.v.NodeAddress.Hex(), "index", index, "currentNonce",
			hex.EncodeToString(currentNonce), "preNonce", hex.EncodeToString(preNonces[index]),
			"target", target, "targetP", targetP, "weight", sv.weights, "x", x, "version", sv.version,
			"blockNumber", sv.blockNumber, "txIndex", sv.txIndex)
	}

	vrfQueue := make(staking.ValidatorQueue, shiftLen)

	log.Debug("Call probabilityElection, sort probability queue", "blockNumber", blockNumber, "currentVersion", currentVersion, "list", svList)

	sort.Sort(svList)
	for index, sv := range svList {
		if index == shiftLen {
			break
		}
		vrfQueue[index] = sv.v
	}

	log.Debug("Call probabilityElection finished", "blockNumber", blockNumber, "currentVersion", currentVersion, "vrfQueue", vrfQueue)

	return vrfQueue, nil
}

/**
Internal expansion function
*/

// previous round validators
func (sk *StakingPlugin) getPreValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValidatorArray, error) {

	targetIndex, err := sk.getPreValIndex(blockHash, blockNumber, isCommit)
	if nil != err {
		return nil, err
	}

	var queue staking.ValidatorQueue

	if !isCommit {
		arr, err := sk.db.GetRoundValListByBlockHash(blockHash, targetIndex.Start, targetIndex.End)
		if snapshotdb.NonDbNotFoundErr(err) {
			return nil, err
		}
		queue = arr

	} else {
		arr, err := sk.db.GetRoundValListByIrr(targetIndex.Start, targetIndex.End)
		if snapshotdb.NonDbNotFoundErr(err) {
			return nil, err
		}
		queue = arr

	}

	if len(queue) == 0 {
		log.Error("Not Found previous validators, the queue length is zero", "isCommit", isCommit, "start", targetIndex.Start,
			"end", targetIndex.End, "current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		return nil, staking.ErrValidatorNoExist
	}

	return &staking.ValidatorArray{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) getPreValIndex(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValArrIndex, error) {
	var targetIndex *staking.ValArrIndex

	var preTargetNumber uint64
	if blockNumber > xutil.ConsensusSize() {
		preTargetNumber = blockNumber - xutil.ConsensusSize()
	}

	var indexArr staking.ValArrIndexQueue

	if !isCommit {
		indexs, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
		if snapshotdb.NonDbNotFoundErr(err) {
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
		if snapshotdb.NonDbNotFoundErr(err) {
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
		log.Error("Not Found previous validators index", "isCommit", isCommit,
			"current blockNumber", blockNumber, "current blockHash", blockHash.Hex(),
			"\nThe round indexs arr", indexArr)
		return nil, staking.ErrValidatorNoExist
	}
	return targetIndex, nil
}

func (sk *StakingPlugin) getCurrValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValidatorArray, error) {

	targetIndex, err := sk.getCurrValIndex(blockHash, blockNumber, isCommit)
	if nil != err {
		return nil, err
	}

	var queue staking.ValidatorQueue

	if !isCommit {
		arr, err := sk.db.GetRoundValListByBlockHash(blockHash, targetIndex.Start, targetIndex.End)
		if snapshotdb.NonDbNotFoundErr(err) {
			return nil, err
		}
		queue = arr

	} else {
		arr, err := sk.db.GetRoundValListByIrr(targetIndex.Start, targetIndex.End)
		if snapshotdb.NonDbNotFoundErr(err) {
			return nil, err
		}
		queue = arr

	}

	if len(queue) == 0 {
		log.Error("Not Found current validators, the queue length is zero", "isCommit", isCommit, "start", targetIndex.Start,
			"end", targetIndex.End, "current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		return nil, staking.ErrValidatorNoExist
	}

	return &staking.ValidatorArray{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) getCurrValIndex(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValArrIndex, error) {
	var targetIndex *staking.ValArrIndex

	var indexArr staking.ValArrIndexQueue

	if !isCommit {
		indexs, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
		if snapshotdb.NonDbNotFoundErr(err) {
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
		if snapshotdb.NonDbNotFoundErr(err) {
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
		log.Error("Not Found current validators index", "isCommit", isCommit,
			"current blockNumber", blockNumber, "current blockHash", blockHash.Hex(),
			"\nThe round indexs arr", indexArr)
		return nil, staking.ErrValidatorNoExist
	}

	return targetIndex, nil
}

func (sk *StakingPlugin) getNextValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValidatorArray, error) {

	targetIndex, err := sk.getNextValIndex(blockHash, blockNumber, isCommit)
	if nil != err {
		return nil, err
	}

	var queue staking.ValidatorQueue

	if !isCommit {
		arr, err := sk.db.GetRoundValListByBlockHash(blockHash, targetIndex.Start, targetIndex.End)
		if snapshotdb.NonDbNotFoundErr(err) {
			return nil, err
		}
		queue = arr

	} else {
		arr, err := sk.db.GetRoundValListByIrr(targetIndex.Start, targetIndex.End)
		if snapshotdb.NonDbNotFoundErr(err) {
			return nil, err
		}
		queue = arr
	}

	if len(queue) == 0 {
		log.Error("Not Found next validators, the queue length is zero", "isCommit", isCommit, "start", targetIndex.Start,
			"end", targetIndex.End, "current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		return nil, staking.ErrValidatorNoExist
	}

	return &staking.ValidatorArray{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) getNextValIndex(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValArrIndex, error) {
	var targetIndex *staking.ValArrIndex

	var indexArr staking.ValArrIndexQueue

	if !isCommit {
		indexs, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
		if snapshotdb.NonDbNotFoundErr(err) {
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
		if snapshotdb.NonDbNotFoundErr(err) {
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
		log.Error("Not Found next validators index", "isCommit", isCommit,
			"current blockNumber", blockNumber, "current blockHash", blockHash.Hex(),
			"\nThe round indexs arr", indexArr)
		return nil, staking.ErrValidatorNoExist
	}

	return targetIndex, nil
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
		log.Error("Not Found current validatorList index", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "input Start", valArr.Start, "input End", valArr.End, "queue", queue)
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

	targetIndex, err := sk.getVeriferIndex(blockHash, blockNumber, isCommit)
	if nil != err {
		return nil, err
	}

	var queue staking.ValidatorQueue

	if !isCommit {
		arr, err := sk.db.GetEpochValListByBlockHash(blockHash, targetIndex.Start, targetIndex.End)
		if snapshotdb.NonDbNotFoundErr(err) {
			return nil, err
		}
		queue = arr
	} else {
		arr, err := sk.db.GetEpochValListByIrr(targetIndex.Start, targetIndex.End)
		if snapshotdb.NonDbNotFoundErr(err) {
			return nil, err
		}
		queue = arr
	}

	if len(queue) == 0 {
		log.Error("Not Found epoch validators, the queue is zero", "isCommit", isCommit, "start", targetIndex.Start,
			"end", targetIndex.End, "current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
		return nil, staking.ErrValidatorNoExist
	}

	return &staking.ValidatorArray{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) getVeriferIndex(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.ValArrIndex, error) {
	var targetIndex *staking.ValArrIndex

	var indexArr staking.ValArrIndexQueue

	if !isCommit {
		indexs, err := sk.db.GetEpochValIndexByBlockHash(blockHash)
		if snapshotdb.NonDbNotFoundErr(err) {
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
		if snapshotdb.NonDbNotFoundErr(err) {
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
		log.Error("Not Found epoch validators index", "isCommit", isCommit,
			"current blockNumber", blockNumber, "current blockHash", blockHash.Hex(),
			"\nThe epoch indexs arr", indexArr)
		return nil, staking.ErrValidatorNoExist
	}
	return targetIndex, nil
}

func (sk *StakingPlugin) setVerifierListAndIndex(blockNumber uint64, blockHash common.Hash, valArr *staking.ValidatorArray) error {

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
		log.Error("Not Found current verifierList index", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "input Start", valArr.Start, "input End", valArr.End,
			"\nThe history epoch indexs arr", queue)
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
	nodeId discover.NodeID, canAddr common.NodeAddress, stakingBlockNum uint64) error {

	endVoteNum, err := gov.GetMaxEndVotingBlock(nodeId, blockHash, state)
	if nil != err {
		return err
	}
	var refundEpoch, maxEndVoteEpoch, targetEpoch uint64
	if endVoteNum != 0 {
		maxEndVoteEpoch = xutil.CalculateEpoch(endVoteNum)
	}

	duration, err := gov.GovernUnStakeFreezeDuration(blockNumber, blockHash)
	if nil != err {
		return err
	}

	refundEpoch = xutil.CalculateEpoch(blockNumber) + duration

	if maxEndVoteEpoch <= refundEpoch {
		targetEpoch = refundEpoch
	} else {
		targetEpoch = maxEndVoteEpoch
	}

	log.Debug("Call addUnStakeItem, AddUnStakeItemStore start", "current blockNumber", blockNumber,
		"govenance max end vote blokNumber", endVoteNum, "unStakeFreeze Epoch", refundEpoch,
		"govenance max end vote epoch", maxEndVoteEpoch, "unstake item target Epoch", targetEpoch,
		"nodeId", nodeId.String())

	if err := sk.db.AddUnStakeItemStore(blockHash, targetEpoch, canAddr, stakingBlockNum, false); nil != err {
		return err
	}
	return nil
}

func (sk *StakingPlugin) addRecoveryUnStakeItem(blockNumber uint64, blockHash common.Hash, nodeId discover.NodeID,
	canAddr common.NodeAddress, stakingBlockNum uint64) error {

	duration, err := gov.GovernZeroProduceFreezeDuration(blockNumber, blockHash)
	if nil != err {
		return err
	}

	targetEpoch := xutil.CalculateEpoch(blockNumber) + duration

	log.Debug("Call addRecoveryUnStakeItem, AddUnStakeItemStore start", "current blockNumber", blockNumber,
		"duration", duration, "unstake item target Epoch", targetEpoch,
		"nodeId", nodeId.String())

	if err := sk.db.AddUnStakeItemStore(blockHash, targetEpoch, canAddr, stakingBlockNum, true); nil != err {
		return err
	}
	return nil
}

// Record the address of the verification node for each consensus round within a certain block range.
func (sk *StakingPlugin) storeRoundValidatorAddrs(blockNumber uint64, blockHash common.Hash, nextStart uint64, array staking.ValidatorQueue) error {
	nextRound := xutil.CalculateRound(nextStart)
	nextEpoch := xutil.CalculateEpoch(nextStart)

	evidenceAge, err := gov.GovernMaxEvidenceAge(blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to storeRoundValidatorAddrs, query Gov SlashFractionDuplicateSign is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
			"err", err)
		return err
	}

	validEpochCount := uint64(evidenceAge + 1)
	validRoundCount := xutil.EpochSize() * validEpochCount

	// Only store the address of last consensus rounds on `validEpochCount` epochs
	if nextEpoch > validEpochCount {
		invalidRound := nextRound - validRoundCount

		boundary, er := sk.db.GetRoundAddrBoundary(blockHash)
		if snapshotdb.NonDbNotFoundErr(er) {
			return er
		}
		if boundary == 0 && (invalidRound-1) >= 0 {
			boundary = invalidRound - 1
		}

		// Clean all outside the boundarys of previous valAddrs
		var count int
		for invalidRound > boundary {
			key := staking.GetRoundValAddrArrKey(invalidRound)
			if err := sk.db.DelRoundValidatorAddrs(blockHash, key); nil != err {
				log.Error("Failed to DelRoundValidatorAddrs", "blockHash", blockHash.TerminalString(), "nextStart", nextStart,
					"nextRound", nextRound, "nextEpoch", nextEpoch, "validEpochCount", validEpochCount, "validRoundCount", validRoundCount, "invalidRound", invalidRound, "key", hex.EncodeToString(key), "err", err)
				return err
			}

			if count == 0 {
				if err := sk.db.SetRoundAddrBoundary(blockHash, nextRound-validRoundCount); nil != err {
					return err
				}
			}
			count++
			invalidRound--
		}

	}
	newKey := staking.GetRoundValAddrArrKey(nextRound)
	newValue := make([]common.NodeAddress, 0, len(array))
	for _, v := range array {
		newValue = append(newValue, v.NodeAddress)
	}
	if err := sk.db.StoreRoundValidatorAddrs(blockHash, newKey, newValue); nil != err {
		log.Error("Failed to StoreRoundValidatorAddrs", "blockHash", blockHash.TerminalString(), "nextStart", nextStart,
			"nextRound", nextRound, "nextEpoch", nextEpoch, "validEpochCount", validEpochCount, "validRoundCount", validRoundCount,
			"validatorLen", len(array), "newKey", hex.EncodeToString(newKey), "err", err)
		return err
	}
	return nil
}

func (sk *StakingPlugin) checkRoundValidatorAddr(blockHash common.Hash, targetBlockNumber uint64, addr common.NodeAddress) (bool, error) {
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

func calcRealRefund(blockNumber uint64, blockHash common.Hash, realtotal, amount *big.Int) *big.Int {
	refundAmount := new(big.Int).SetInt64(0)
	sub := new(big.Int).Sub(realtotal, amount)
	// When the sub less than threshold
	if ok, _ := CheckOperatingThreshold(blockNumber, blockHash, sub); !ok {
		refundAmount = realtotal
	} else {
		refundAmount = amount
	}
	return refundAmount
}

func buildCanHex(can *staking.Candidate) *staking.CandidateHex {
	return &staking.CandidateHex{
		NodeId:               can.NodeId,
		BlsPubKey:            can.BlsPubKey,
		StakingAddress:       can.StakingAddress,
		BenefitAddress:       can.BenefitAddress,
		RewardPer:            can.RewardPer,
		NextRewardPer:        can.NextRewardPer,
		RewardPerChangeEpoch: can.RewardPerChangeEpoch,
		StakingTxIndex:       can.StakingTxIndex,
		ProgramVersion:       can.ProgramVersion,
		Status:               can.Status,
		StakingEpoch:         can.StakingEpoch,
		StakingBlockNum:      can.StakingBlockNum,
		Shares:               (*hexutil.Big)(can.Shares),
		Released:             (*hexutil.Big)(can.Released),
		ReleasedHes:          (*hexutil.Big)(can.ReleasedHes),
		RestrictingPlan:      (*hexutil.Big)(can.RestrictingPlan),
		RestrictingPlanHes:   (*hexutil.Big)(can.RestrictingPlanHes),
		DelegateEpoch:        can.DelegateEpoch,
		DelegateTotal:        (*hexutil.Big)(can.DelegateTotal),
		DelegateTotalHes:     (*hexutil.Big)(can.DelegateTotalHes),
		Description:          can.Description,
		DelegateRewardTotal:  (*hexutil.Big)(can.DelegateRewardTotal),
	}
}

func CheckStakeThreshold(blockNumber uint64, blockHash common.Hash, stake *big.Int) (bool, *big.Int) {

	threshold, err := gov.GovernStakeThreshold(blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to CheckStakeThreshold, query governParams is failed", "err", err)
		return false, common.Big0
	}

	return stake.Cmp(threshold) >= 0, threshold
}

func CheckOperatingThreshold(blockNumber uint64, blockHash common.Hash, balance *big.Int) (bool, *big.Int) {

	threshold, err := gov.GovernOperatingThreshold(blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to CheckOperatingThreshold, query governParams is failed", "err", err)
		return false, common.Big0
	}
	return balance.Cmp(threshold) >= 0, threshold
}
