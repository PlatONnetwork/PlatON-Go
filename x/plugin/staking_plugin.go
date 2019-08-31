package plugin

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"sync"

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

var (
	AccountVonNotEnough        = common.NewBizError("The von of account is not enough")
	BalanceOperationTypeErr    = common.NewBizError("Balance OperationType is wrong")
	DelegateVonNotEnough       = common.NewBizError("The von of delegate is not enough")
	WithdrewDelegateVonCalcErr = common.NewBizError("Withdrew delegate von calculate err")
	ParamsErr                  = common.NewBizError("The fn params err")
	BlockNumberDisordered      = common.NewBizError("The blockNumber is disordered")
	VonAmountNotRight          = common.NewBizError("The amount of von is not right")
	CandidateNotExist          = common.NewBizError("The candidate is not exist")
	ValidatorNotExist          = common.NewBizError("The validator is not exist")
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

	if xutil.IsSettlementPeriod(header.Number.Uint64()) {

		log.Info("Call EndBlock on staking plugin, IsSettlementPeriod", "blockNumber", header.Number, "blockHash", blockHash.String(), "epoch", epoch)

		// handle UnStaking Item
		err := sk.HandleUnCandidateItem(state, header.Number.Uint64(), blockHash, epoch)
		if nil != err {
			log.Error("Failed to call HandleUnCandidateItem on stakingPlugin EndBlock",
				"blockNumber", header.Number.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		// hanlde UnDelegate Item
		err = sk.HandleUnDelegateItem(state, header.Number.Uint64(), blockHash, epoch)
		if nil != err {
			log.Error("Failed to call HandleUnDelegateItem on stakingPlugin EndBlock",
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

		log.Info("Call EndBlock on staking plugin, IsElection",
			"blockNumber", header.Number, "blockHash", blockHash.String(), "epoch", epoch)

		// ELection next round validators
		err := sk.Election(blockHash, header, state)
		if nil != err {
			log.Error("Failed to call Election on stakingPlugin EndBlock",
				"blockNumber", header.Number.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}

	}
	//log.Info("Finished EndBlock on staking plugin", "blockNumber", header.Number, "blockHash", blockHash.String(), "epoch", epoch)
	return nil
}

func (sk *StakingPlugin) Confirmed(block *types.Block) error {

	log.Info("Call Confirmed on staking plugin", "blockNumber", block.Number(), "blockHash", block.Hash().String())

	if xutil.IsElection(block.NumberU64()) {

		next, err := sk.getNextValList(block.Hash(), block.NumberU64(), QueryStartNotIrr)
		if nil != err {
			log.Error("Failed to Query Next validators on stakingPlugin Confirmed When Election block",
				"blockNumber", block.Number().Uint64(), "blockHash", block.Hash().Hex(), "err", err)
			return err
		}

		current, err := sk.getCurrValList(block.Hash(), block.NumberU64(), QueryStartNotIrr)
		if nil != err {
			log.Error("Failed to Query Current Round validators on stakingPlugin Confirmed When Election block",
				"blockNumber", block.Number().Uint64(), "blockHash", block.Hash().Hex(), "err", err)
			return err
		}
		result := distinct(next.Arr, current.Arr)
		if len(result) > 0 {
			sk.addConsensusNode(result)
			log.Debug("stakingPlugin addConsensusNode success",
				"blockNumber", block.NumberU64(), "blockHash", block.Hash().Hex(), "size", len(result))
		}
	}

	log.Info("Finished Confirmed on staking plugin", "blockNumber", block.Number(), "blockHash", block.Hash().String())
	return nil
}

func distinct(list, target staking.ValidatorQueue) staking.ValidatorQueue {
	currentMap := make(map[discover.NodeID]bool)
	for _, v := range target {
		currentMap[v.NodeId] = true
	}
	result := make(staking.ValidatorQueue, 0)
	for _, v := range list {
		if _, ok := currentMap[v.NodeId]; !ok {
			result = append(result, v)
		}
	}
	return result
}

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
			return AccountVonNotEnough
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
		return common.BizErrorf("%s, got type is: %d, need type: %d or %d", BalanceOperationTypeErr.Error(),
			typ, FreeOrigin, RestrictingPlanOrigin)
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

	amount := common.Big0
	if typ == FreeOrigin {
		amount = can.ReleasedHes
	} else if typ == RestrictingPlanOrigin {
		amount = can.RestrictingPlanHes
	} else {
		// this is never be
		return nil
	}

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) == 0 || contract_balance.Cmp(amount) < 0 {
		log.Error("Failed to RollBackStaking: the balance is invalid of stakingContracr Account",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "nodeAddr", addr.String(),
			"contract_balance", contract_balance, "rollback amount", amount)
		panic("the balance is invalid of stakingContracr Account")
	}

	if blockNumber.Uint64() != can.StakingBlockNum {
		return common.BizErrorf("%v: current blockNumber is not equal stakingBlockNumber, can not rollback staking, current blockNumber: %d, can.stakingNumber: %d", ParamsErr, blockNumber.Uint64(), can.StakingBlockNum)
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
		return common.BizErrorf("%s, got type is: %d, need type: %d or %d", BalanceOperationTypeErr.Error(),
			typ, FreeOrigin, RestrictingPlanOrigin)
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
			return AccountVonNotEnough
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
		return common.BizErrorf("%s, got type is: %d, need type: %d or %d", BalanceOperationTypeErr.Error(),
			typ, FreeOrigin, RestrictingPlanOrigin)
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

	total := calCanTotalAmount(can)

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) == 0 || contract_balance.Cmp(total) < 0 {
		log.Error("Failed to withdrewStakeAmount: the balance is invalid of stakingContracr Account",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "contract_balance",
			contract_balance, "withdrewStake amount", total)
		panic("the balance is invalid of stakingContracr Account")
	}

	// Direct return of money during the hesitation period
	// Return according to the way of coming
	if can.ReleasedHes.Cmp(common.Big0) > 0 {

		state.AddBalance(can.StakingAddress, can.ReleasedHes)
		state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)
		can.ReleasedHes = common.Big0
	}

	if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {

		err := rt.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
		if nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
				"stakingAddr", can.StakingAddress.Hex(), "restrictingPlanHes", can.RestrictingPlanHes, "err", err)
			return err
		}

		can.RestrictingPlanHes = common.Big0
	}

	if can.Released.Cmp(common.Big0) > 0 || can.RestrictingPlan.Cmp(common.Big0) > 0 {

		if err := sk.addUnStakeItem(state, blockNumber, blockHash, epoch, can.NodeId, canAddr); nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: Add UnStakeItemStore failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return err
		}
	}

	can.Shares = common.Big0
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
		addr, err := sk.db.GetUnStakeItemStore(blockHash, epoch, uint64(index))
		if nil != err {
			log.Error("Failed to HandleUnCandidateItem: Query the unStakeItem node addr is failed",
				"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		if _, ok := filterAddr[addr]; ok {
			if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
				log.Error("Failed to HandleUnCandidateItem: Delete already handle unstakeItem failed",
					"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
				return err
			}
			continue
		}

		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err && err != snapshotdb.ErrNotFound {
			log.Error("Failed to HandleUnCandidateItem: Query candidate failed",
				"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "canAddr", addr.Hex(), "err", err)
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

		// Second handle balabala ...
		if err := sk.handleUnStake(state, blockNumber, blockHash, epoch, addr, can); nil != err {
			return err
		}

		if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
			log.Error("Failed to HandleUnCandidateItem: Delete unstakeItem failed",
				"blockNUmber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		filterAddr[addr] = struct{}{}
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

	total := calCanTotalAmount(can)

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) == 0 || contract_balance.Cmp(total) < 0 {
		log.Error("Failed to handleUnStake: the balance is invalid of stakingContracr Account",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "contract_balance",
			contract_balance, "handle unstake amount", total)
		panic("the balance is invalid of stakingContracr Account")
	}

	refundReleaseFn := func(balance *big.Int) *big.Int {
		if balance.Cmp(common.Big0) > 0 {

			state.AddBalance(can.StakingAddress, balance)
			state.SubBalance(vm.StakingContractAddr, balance)
			return common.Big0
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
				return common.Big0, err
			}
			return common.Big0, nil
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
			Reduction:          (*hexutil.Big)(del.Reduction),
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
			Reduction:          (*hexutil.Big)(del.Reduction),
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
			Reduction:          (*hexutil.Big)(del.Reduction),
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
			return AccountVonNotEnough
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
		return common.BizErrorf("%s, got type is: %d, need type: %d or %d", BalanceOperationTypeErr.Error(),
			typ, FreeOrigin, RestrictingPlanOrigin)
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

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) == 0 || contract_balance.Cmp(amount) < 0 {
		log.Error("Failed to WithdrewDelegate: the balance is invalid of stakingContracr Account",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"contract_balance", contract_balance, "withdrew del amount", amount)
		panic("the balance is invalid of stakingContracr Account")
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to WithdrewDelegate on stakingPlugin: nodeId parse addr failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
		return err
	}

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	can, err := sk.db.GetCandidateStore(blockHash, canAddr)
	if nil != err {
		log.Error("Failed to WithdrewDelegate on stakingPlugin: Query candidate info failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
		return err
	}

	aboutRelease := new(big.Int).Add(del.Released, del.ReleasedHes)
	aboutRestrictingPlan := new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
	total := new(big.Int).Add(aboutRelease, aboutRestrictingPlan)

	lazyCalcDelegateAmount(epoch, del)

	// inner Fn
	subDelegateFn := func(source, sub *big.Int) (*big.Int, *big.Int) {

		state.AddBalance(delAddr, sub)
		state.SubBalance(vm.StakingContractAddr, sub)
		return new(big.Int).Sub(source, sub), common.Big0
	}

	refundFn := func(refund, aboutRelease, aboutRestrictingPlan *big.Int) (*big.Int, *big.Int, *big.Int, error) {

		refundTmp := refund
		releaseTmp := aboutRelease
		restrictingPlanTmp := aboutRestrictingPlan

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
					log.Error("Failed to WithdrewDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
						"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "balance", restrictingPlanTmp,
						"err", err)
					return refundTmp, releaseTmp, restrictingPlanTmp, err
				}

				refundTmp = new(big.Int).Sub(refundTmp, restrictingPlanTmp)
				restrictingPlanTmp = common.Big0

			} else if refundTmp.Cmp(restrictingPlanTmp) < 0 {
				// When remain is less than or equal to del.RestrictingPlanHes/del.RestrictingPlan
				err := rt.ReturnLockFunds(delAddr, refundTmp, state)
				if nil != err {
					log.Error("Failed to WithdrewDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
						"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "balance", refundTmp,
						"err", err)
					return refundTmp, releaseTmp, restrictingPlanTmp, err
				}

				restrictingPlanTmp = new(big.Int).Sub(restrictingPlanTmp, refundTmp)
				refundTmp = common.Big0
			}
		}

		return refundTmp, releaseTmp, restrictingPlanTmp, nil
	}

	del.DelegateEpoch = uint32(epoch)

	switch {

	// When the related candidate info does not exist
	case nil == can, nil != can && stakingBlockNum < can.StakingBlockNum,
		nil != can && stakingBlockNum == can.StakingBlockNum && staking.Is_Invalid(can.Status):

		// First need to deduct the von that is being refunded
		realtotal := new(big.Int).Sub(total, del.Reduction)

		log.Info("Call WithdrewDelegate, the candidate is invalid or no exist", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.String(), "nodeId", nodeId.String(), "StakingNum", stakingBlockNum, "amount", amount, "realtotal", realtotal,
			"total", total, "redution", del.Reduction)

		if realtotal.Cmp(amount) < 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: the amount of valid delegate is not enough",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
				"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "delegate amount", realtotal,
				"withdrew amount", amount)
			return common.BizErrorf("withdrewDelegate err: %s, delegate von: %s, withdrew von: %s",
				DelegateVonNotEnough.Error(), realtotal.String(), amount.String())
		}

		refundAmount := common.Big0
		sub := new(big.Int).Sub(realtotal, amount)

		// When the sub less than threshold
		if !xutil.CheckMinimumThreshold(sub) {
			refundAmount = realtotal
		} else {
			refundAmount = amount
		}

		realSub := refundAmount

		log.Debug("Call WithdrewDelegate, the candidate is invalid or no exist", "realSub", realSub, "withdrew amount", amount)

		// handle delegate on Hesitate period
		if refundAmount.Cmp(common.Big0) > 0 {
			rm, rbalance, lbalance, err := refundFn(refundAmount, del.ReleasedHes, del.RestrictingPlanHes)
			if nil != err {
				return err
			}
			refundAmount, del.ReleasedHes, del.RestrictingPlanHes = rm, rbalance, lbalance
		}

		// handle delegate on Effective period
		if refundAmount.Cmp(common.Big0) > 0 {
			rm, rbalance, lbalance, err := refundFn(refundAmount, del.Released, del.RestrictingPlan)
			if nil != err {
				return err
			}
			refundAmount, del.Released, del.RestrictingPlan = rm, rbalance, lbalance
		}

		if refundAmount.Cmp(common.Big0) != 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: the withdrew ramain is not zero",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
				"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "del balance", total,
				"realtatol", realtotal, "redution", del.Reduction, "withdrew balance", amount,
				"realSub amount", realSub, "withdrew remain", refundAmount)
			return WithdrewDelegateVonCalcErr
		}

		// If realtatol had full sub
		// AND redution is zero
		// clean the delegate info
		if realtotal.Cmp(realSub) == 0 && del.Reduction.Cmp(common.Big0) == 0 {

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

	// Illegal parameter
	case nil != can && stakingBlockNum > can.StakingBlockNum:
		log.Error("Failed to WithdrewDelegate on stakingPlugin: the stakeBlockNum invalid",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "fn.stakeBlockNum", stakingBlockNum,
			"can.stakeBlockNum", can.StakingBlockNum)
		return ParamsErr

	// When the delegate is normally revoked
	case nil != can && stakingBlockNum == can.StakingBlockNum && !staking.Is_Invalid(can.Status):

		// First need to deduct the von that is being refunded
		realtotal := new(big.Int).Sub(total, del.Reduction)

		log.Info("Call WithdrewDelegate, the candidate is valid", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"delAddr", delAddr.String(), "nodeId", nodeId.String(), "StakingNum", stakingBlockNum, "amount", amount, "realtotal", realtotal,
			"total", total, "redution", del.Reduction)

		if realtotal.Cmp(amount) < 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: the amount of valid delegate is not enough",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
				"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "delegate amount", realtotal,
				"withdrew amount", amount)
			return common.BizErrorf("withdrewDelegate err: %s, delegate von: %s, withdrew von: %s",
				DelegateVonNotEnough.Error(), realtotal.String(), amount.String())
		}

		refundAmount := common.Big0
		sub := new(big.Int).Sub(realtotal, amount)

		// When the sub less than threshold
		if !xutil.CheckMinimumThreshold(sub) {
			refundAmount = realtotal
		} else {
			refundAmount = amount
		}

		realSub := refundAmount

		log.Debug("Call WithdrewDelegate, the candidate is valid", "realSub", realSub, "withdrew amount", amount)

		// handle delegate on Hesitate period
		if refundAmount.Cmp(common.Big0) > 0 {
			rm, rbalance, lbalance, err := refundFn(refundAmount, del.ReleasedHes, del.RestrictingPlanHes)
			if nil != err {
				return err
			}
			refundAmount, del.ReleasedHes, del.RestrictingPlanHes = rm, rbalance, lbalance
		}

		save_or_del := false // false: save, true: delete

		// handle delegate on Effective period
		if refundAmount.Cmp(common.Big0) > 0 {

			// add a UnDelegateItem
			if err := sk.addUnDelegateItem(blockNumber.Uint64(), blockHash, delAddr, nodeId, epoch, stakingBlockNum, refundAmount); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin：add a UnDelegateItem failed", "blockNumber",
					blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
					"stakingBlockNum", stakingBlockNum, "current epoch", epoch, "refundAmount", refundAmount, "err", err)
				return err
			}

			del.Reduction = new(big.Int).Add(del.Reduction, refundAmount)

		} else {

			hes := new(big.Int).Add(del.ReleasedHes, del.RestrictingPlanHes)
			noHes := new(big.Int).Add(del.Released, del.RestrictingPlan)
			add := new(big.Int).Add(hes, noHes)

			// Need to clean delegate info
			if del.Reduction.Cmp(common.Big0) == 0 && add.Cmp(common.Big0) == 0 {
				save_or_del = true
			}
		}

		if !save_or_del {

			if err := sk.db.SetDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum, del); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Store delegate info is failed", "blockNumber",
					blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
					"stakingBlockNum", stakingBlockNum, "err", err)
				return err
			}
		} else {

			// Clean delegate info
			if err := sk.db.DelDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Delete delegate info is failed", "blockNumber",
					blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
					"stakingBlockNum", stakingBlockNum, "err", err)
				return err
			}
		}

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

func (sk *StakingPlugin) HandleUnDelegateItem(state xcom.StateDB, blockNumber uint64, blockHash common.Hash, epoch uint64) error {

	log.Debug("Call HandleUnDelegateItem start", "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch)

	unDelegateCount, err := sk.db.GetUnDelegateCountStore(blockHash, epoch)
	switch {
	case nil != err && err != snapshotdb.ErrNotFound:
		return err
	case nil != err && err == snapshotdb.ErrNotFound:
		unDelegateCount = 0
	}

	if unDelegateCount == 0 {
		return nil
	}

	for index := 1; index <= int(unDelegateCount); index++ {
		unDelegateItem, err := sk.db.GetUnDelegateItemStore(blockHash, epoch, uint64(index))

		if nil != err {
			log.Error("Failed to HandleUnDelegateItem: Query the unStakeItem is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch, "err", err)
			return err
		}

		del, err := sk.db.GetDelegateStoreBySuffix(blockHash, unDelegateItem.KeySuffix)
		if nil != err && err != snapshotdb.ErrNotFound {
			log.Error("Failed to HandleUnDelegateItem: Query delegate info is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch, "err", err)
			return err
		}

		// This maybe be nil
		if (nil != err && err == snapshotdb.ErrNotFound) || nil == del {
			if err := sk.db.DelUnDelegateItemStore(blockHash, epoch, uint64(index)); nil != err {
				log.Error("Failed to HandleUnDelegateItem: Delegate is no exist, Delete unDelegateItem failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
				return err
			}
			continue
		}

		if err := sk.handleUnDelegate(state, blockNumber, blockHash, epoch, unDelegateItem, del); nil != err {
			return err
		}

		// clean item
		if err := sk.db.DelUnDelegateItemStore(blockHash, epoch, uint64(index)); nil != err {
			log.Error("Failed to HandleUnDelegateItem: Delete unDelegateItem failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}

	// clean count
	if err := sk.db.DelUnDelegateCountStore(blockHash, epoch); nil != err {
		log.Error("Failed to HandleUnDelegateItem: Delete unDelegateCount failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) handleUnDelegate(state xcom.StateDB, blockNumber uint64,
	blockHash common.Hash, epoch uint64, unDel *staking.UnDelegateItem, del *staking.Delegation) error {

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	// Maybe equal zero (maybe slashed)
	// must compare the undelegate amount and contract's balance
	if contract_balance.Cmp(common.Big0) == 0 || contract_balance.Cmp(unDel.Amount) < 0 {
		log.Error("Failed to handleUnDelegate: the balance is invalid of stakingContracr Account",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "contract_balance", contract_balance,
			"unDel.Amount", unDel.Amount)
		panic("the balance is invalid of stakingContracr Account")
	}

	// del addr
	delAddrByte := unDel.KeySuffix[0:common.AddressLength]
	delAddr := common.BytesToAddress(delAddrByte)

	nodeIdLen := discover.NodeIDBits / 8

	nodeIdByte := unDel.KeySuffix[common.AddressLength : common.AddressLength+nodeIdLen]
	nodeId := discover.MustBytesID(nodeIdByte)

	stakeBlockNum := unDel.KeySuffix[common.AddressLength+nodeIdLen:]
	num := common.BytesToUint64(stakeBlockNum)

	lazyCalcDelegateAmount(epoch, del)

	// undelegate amount
	amount := unDel.Amount

	aboutRelease := new(big.Int).Add(del.Released, del.ReleasedHes)
	aboutRestrictingPlan := new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
	total := new(big.Int).Add(aboutRelease, aboutRestrictingPlan)

	if amount.Cmp(del.Reduction) == 0 && del.Reduction.Cmp(total) == 0 { // full withdrawal

		log.Info("Call handleUnDelegate, full withdraw", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "epoch", epoch, "delAddr", delAddr, "full refund", total)

		refundReleaseFn := func(balance *big.Int) *big.Int {
			if balance.Cmp(common.Big0) > 0 {
				state.AddBalance(delAddr, balance)
				state.SubBalance(vm.StakingContractAddr, balance)
				return common.Big0
			}
			return balance
		}

		//del.ReleasedHes = refundReleaseFn(del.ReleasedHes)
		del.Released = refundReleaseFn(del.Released)

		refundRestrictingPlanFn := func(title string, balance *big.Int) (*big.Int, error) {

			if balance.Cmp(common.Big0) > 0 {
				err := rt.ReturnLockFunds(delAddr, balance, state)
				if nil != err {
					log.Error("Failed to handleUnDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						title, balance, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch,
						"delAddr", delAddr.Hex(), "err", err)
					return common.Big0, err
				}
				return common.Big0, nil
			}

			return balance, nil
		}

		//if balance, err := refundRestrictingPlanFn("RestrictingPlanHes", del.RestrictingPlanHes); nil != err {
		//	return err
		//} else {
		//	del.RestrictingPlanHes = balance
		//}

		if balance, err := refundRestrictingPlanFn("RestrictingPlan", del.RestrictingPlan); nil != err {
			return err
		} else {
			del.RestrictingPlan = balance
		}

		// clean the delegate
		if err := sk.db.DelDelegateStoreBySuffix(blockHash, unDel.KeySuffix); nil != err {
			log.Error("Failed to handleUnDelegate on stakingPlugin: Delete delegate info is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch, "err", err)
			return err
		}

	} else { // few withdrawal

		refund_remain := amount

		log.Info("Call handleUnDelegate, few withdraw", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "epoch", epoch, "delAddr", delAddr, "few refund", refund_remain)

		refundReleaseFn := func(balance, refund *big.Int) (*big.Int, *big.Int) {

			if balance.Cmp(common.Big0) > 0 && refund.Cmp(common.Big0) > 0 {

				if refund.Cmp(balance) >= 0 {

					state.SubBalance(vm.StakingContractAddr, balance)
					state.AddBalance(delAddr, balance)

					return common.Big0, new(big.Int).Sub(refund, balance)
				} else {
					state.SubBalance(vm.StakingContractAddr, refund)
					state.AddBalance(delAddr, refund)

					return new(big.Int).Sub(balance, refund), common.Big0
				}
			}

			return balance, refund
		}

		//hes, rm := refundReleaseFn(del.ReleasedHes, refund_remain)
		//del.ReleasedHes, refund_remain = hes, rm
		noHes, rm := refundReleaseFn(del.Released, refund_remain)
		del.Released, refund_remain = noHes, rm

		refundRestrictingPlanFn := func(title string, balance, refund *big.Int) (*big.Int, *big.Int, error) {

			if balance.Cmp(common.Big0) > 0 && refund.Cmp(common.Big0) > 0 {

				if refund.Cmp(balance) >= 0 {

					err := rt.ReturnLockFunds(delAddr, balance, state)
					if nil != err {
						log.Error("Failed to handleUnDelegate on stakingPlugin: call Restricting ReturnLockFunds() return "+title+" is failed",
							title, balance, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch, "delAddr", delAddr.Hex(), "err", err)
						return common.Big0, common.Big0, err
					}
					return common.Big0, new(big.Int).Sub(refund, balance), nil
				} else {
					err := rt.ReturnLockFunds(delAddr, refund, state)
					if nil != err {
						log.Error("Failed to handleUnDelegate on stakingPlugin: call Restricting ReturnLockFunds() return "+title+" is failed",
							"refund amount", refund, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch, "delAddr", delAddr.Hex(), "err", err)
						return common.Big0, common.Big0, err
					}

					return new(big.Int).Sub(balance, refund), common.Big0, nil
				}
			}

			return balance, refund, nil
		}

		//if balance, re, err := refundRestrictingPlanFn("RestrictingPlanHes", del.RestrictingPlanHes, refund_remain); nil != err {
		//	return err
		//} else {
		//	del.RestrictingPlanHes, refund_remain = balance, re
		//}

		if balance, re, err := refundRestrictingPlanFn("RestrictingPlan", del.RestrictingPlan, refund_remain); nil != err {
			return err
		} else {
			del.RestrictingPlan, refund_remain = balance, re
		}

		if refund_remain.Cmp(common.Big0) > 0 {
			log.Error("Failed to call handleUnDelegate: remain is not zero", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "epoch", epoch, "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
				"stakeBlockNumber", num, "refund amount", amount, "refund remain", refund_remain)
			return VonAmountNotRight
		}

		del.Reduction = new(big.Int).Sub(del.Reduction, amount)

		del.DelegateEpoch = uint32(epoch)

		if err := sk.db.SetDelegateStoreBySuffix(blockHash, unDel.KeySuffix, del); nil != err {
			log.Error("Failed to handleUnDelegate on stakingPlugin: Store delegate info is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "epoch", epoch, "err", err)
			return err
		}
	}

	return nil
}

func (sk *StakingPlugin) ElectNextVerifierList(blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {

	log.Info("Call ElectNextVerifierList Start", "blockNumber", blockNumber, "blockHash", blockHash.Hex())

	old_verifierArr, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to ElectNextVerifierList: No found the VerifierLIst", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// todo test
	xcom.PrintObject("Call ElectNextVerifierList old verifier list", old_verifierArr)

	if old_verifierArr.End != blockNumber {
		log.Error("Failed to ElectNextVerifierList: this blockNumber invalid", "Old Epoch End blockNumber",
			old_verifierArr.End, "Current blockNumber", blockNumber)
		return common.BizErrorf("The BlockNumber invalid, Old Epoch End blockNumber: %d, Current blockNumber: %d",
			old_verifierArr.End, blockNumber)
	}

	// caculate the new epoch start and end
	start := old_verifierArr.End + 1
	end := old_verifierArr.End + xutil.CalcBlocksEachEpoch()

	new_verifierArr := &staking.Validator_array{
		Start: start,
		End:   end,
	}

	curr_version := gov.GetVersionForStaking(state)
	currVersion := xutil.CalcVersion(curr_version)

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

			log.Debug("Call ElectNextVerifierList: the can ProgramVersion is less than currVersion",
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
		panic(common.BizErrorf("Failed to ElectNextVerifierList: Select zero size validators~"))
	}

	new_verifierArr.Arr = queue

	err = sk.setVerifierList(blockHash, new_verifierArr)
	if nil != err {
		log.Error("Failed to ElectNextVerifierList: Set Next Epoch VerifierList is failed", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// todo test
	xcom.PrintObject("Call ElectNextVerifierList new verifier list", new_verifierArr)

	log.Info("Call ElectNextVerifierList end", "new epoch validators length", len(queue), "loopNum", count)
	return nil
}

func (sk *StakingPlugin) GetVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.ValidatorExQueue, error) {

	verifierList, err := sk.getVerifierList(blockHash, blockNumber, isCommit)
	if nil != err {
		return nil, err
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
		} else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
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
		return nil, common.BizErrorf("ListVerifierNodeID failed: %s, start: %d, end: %d, currentNumer: %d",
			BlockNumberDisordered.Error(), verifierList.Start, verifierList.End, blockNumber)
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
				return nil, err
			}
			can = c
		} else {
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
// 0: 	Query previous round consensus validator
// 1:  	Query current round consensus validaor
// 2:  	Query next round consensus validator
func (sk *StakingPlugin) GetValidatorList(blockHash common.Hash, blockNumber uint64, flag uint, isCommit bool) (
	staking.ValidatorExQueue, error) {

	var validatorArr *staking.Validator_array

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
		} else {
			c, err := sk.db.GetCandidateStoreByIrr(v.NodeAddress)
			if nil != err {
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

	var validatorArr *staking.Validator_array

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
		} else {
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

	//var iter iterator.Iterator
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

	log.Info("Call Election Start", "blockHash", blockHash.Hex(), "blockNumber", header.Number.Uint64())

	blockNumber := header.Number.Uint64()

	// the validators of Current Epoch
	verifiers, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to call Election: No found current epoch validators", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return ValidatorNotExist
	}

	// the validators of Current Round
	curr, err := sk.getCurrValList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to Election: No found the current round validators", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return ValidatorNotExist
	}

	log.Info("Call Election start", "Curr epoch validators length", len(verifiers.Arr), "Curr round validators length", len(curr.Arr))
	xcom.PrintObject("Call Election Curr validators", curr)

	if blockNumber != (curr.End - xcom.ElectionDistance()) {
		log.Error("Failed to Election: Current blockNumber invalid", "Target blockNumber",
			curr.End-xcom.ElectionDistance(), "blockNumber", blockNumber, "blockHash", blockHash.Hex())
		return common.BizErrorf("The BlockNumber invalid, Target blockNumber: %d, Current blockNumber: %d",
			curr.End-xcom.ElectionDistance(), blockNumber)
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
		// Collecting removed as a result of being slashed
		// That is not withdrew to invalid
		//
		// eg. (lowRatio and must delete) OR (lowRatio and balance no enough) OR duplicateSign
		//
		checkHaveSlash := func() bool {
			return staking.Is_Invalid_LowRatioDel(can.Status) ||
				staking.Is_Invalid_LowRatio_NotEnough(can.Status) ||
				staking.Is_Invalid_DuplicateSign(can.Status)
		}

		if checkHaveSlash() {
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
		if !staking.Is_Invalid(can.Status) && staking.Is_LowRatio(can.Status) {
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
		"diffQueueLen", len(diffQueue), "vrfQueueLen", len(vrfQueue))

	nextQueue := shuffle(invalidLen, curr.Arr, vrfQueue)

	if len(nextQueue) == 0 {
		panic("The Next Round Validator is empty, blockNumber: " + fmt.Sprint(blockNumber))
	}

	next := &staking.Validator_array{
		Start: start,
		End:   end,
		Arr:   nextQueue,
	}

	if err := sk.setRoundValList(blockHash, next); nil != err {
		log.Error("Failed to SetNextValidatorList on Election", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// todo test
	if len(slashAddrQueue) != 0 {
		log.Debug("Election Remove Slashing nodeId", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex())
		xcom.PrintObject("Election Remove Slashing nodeId", slashAddrQueue)
	}

	// todo test
	if len(withdrewQueue) != 0 {
		log.Debug("Election Remove Withdrew nodeId", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex())
		xcom.PrintObject("Election Remove Withdrew nodeId", withdrewQueue)
	}

	// todo test
	if len(lowVersionQueue) != 0 {
		log.Debug("Election Remove Low version nodeId", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex())
		xcom.PrintObject("Election Remove Low version nodeId", lowVersionQueue)
	}

	// update candidate status
	// Must sort
	for _, canAddr := range lowRatioValidAddrs {

		can := lowRatioValidMap[canAddr]

		// clean the low package ratio status
		can.Status &^= staking.LowRatio

		// TODO test
		log.Debug("Call Election, clean lowratio", "nodeId", can.NodeId.String())
		xcom.PrintObject("Call Election, clean lowratio", can)

		addr, _ := xutil.NodeId2Addr(can.NodeId)
		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Failed to Store Candidate on Election", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
			return err
		}
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
	return next
}

func (sk *StakingPlugin) SlashCandidates(state xcom.StateDB, blockHash common.Hash, blockNumber uint64,
	nodeId discover.NodeID, amount *big.Int, slashType int, caller common.Address) error {

	log.Info("Call SlashCandidates start", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"slashType", slashType, "nodeId", nodeId.String(), "amount", amount,
		"reporter", caller.Hex())

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) == 0 || contract_balance.Cmp(amount) < 0 {
		log.Error("Failed to SlashCandidates: the balance is invalid of stakingContracr Account", "contract_balance",
			contract_balance, "slash amount", amount)
		panic("the balance is invalid of stakingContracr Account")
	}

	slashTypeIsWrong := func() bool {
		return uint32(slashType) != staking.LowRatio &&
			uint32(slashType) != staking.LowRatioDel &&
			uint32(slashType) != staking.DuplicateSign
	}
	if slashTypeIsWrong() {
		log.Error("Failed to SlashCandidates: the slashType is wrong", "slashType", slashType,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return common.BizErrorf("Failed to SlashCandidates: the slashType is wrong, slashType: %d", slashType)
	}

	canAddr, _ := xutil.NodeId2Addr(nodeId)
	can, err := sk.db.GetCandidateStore(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Call SlashCandidates: Query can is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	if nil == can {
		log.Error("Call SlashCandidates: the can is empty", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return CandidateNotExist
	}

	epoch := xutil.CalculateEpoch(blockNumber)

	lazyCalcStakeAmount(epoch, can)

	/**
	Balance that can only be effective for Slash
	*/
	total := new(big.Int).Add(can.Released, can.RestrictingPlan)

	if total.Cmp(amount) < 0 {
		log.Error("Failed to SlashCandidates: the candidate total staking amount is not enough",
			"candidate total amount", total, "slashing amount", amount, "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return common.BizErrorf("Failed to SlashCandidates: the candidate total effective staking amount is not enough"+
			", candidate total amount:%s, slashing amount: %s", total, amount)
	}

	// clean the candidate power, first
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Call SlashCandidates: Delete candidate old power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return err
	}

	slash := amount

	slashFunc := func(title string, slashAmount, can_balance *big.Int, isNotify bool) (*big.Int, *big.Int, error) {

		// check zero value
		// If there is a zero value, no logic is done.
		if can_balance.Cmp(common.Big0) == 0 || slashAmount.Cmp(common.Big0) == 0 {
			return slashAmount, can_balance, nil
		}

		slashAmountTmp := common.Big0
		balanceTmp := common.Big0

		if slashAmount.Cmp(can_balance) >= 0 {

			state.SubBalance(vm.StakingContractAddr, can_balance)

			if staking.Is_DuplicateSign(uint32(slashType)) {
				state.AddBalance(caller, can_balance)
			} else {
				state.AddBalance(vm.RewardManagerPoolAddr, can_balance)
			}

			if isNotify {
				err := rt.SlashingNotify(can.StakingAddress, can_balance, state)
				if nil != err {
					log.Error("Failed to SlashCandidates: call restrictingPlugin SlashingNotify() failed", "slashed amount", can_balance,
						"slash:", title, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
					return slashAmountTmp, balanceTmp, err
				}
			}

			slashAmountTmp = new(big.Int).Sub(slashAmount, can_balance)
			balanceTmp = common.Big0

		} else {
			state.SubBalance(vm.StakingContractAddr, slashAmount)
			if staking.Is_DuplicateSign(uint32(slashType)) {
				state.AddBalance(caller, slashAmount)
			} else {
				state.AddBalance(vm.RewardManagerPoolAddr, slashAmount)
			}

			if isNotify {
				err := rt.SlashingNotify(can.StakingAddress, slashAmount, state)
				if nil != err {
					log.Error("Failed to SlashCandidates: call restrictingPlugin SlashingNotify() failed", "slashed amount", slashAmount,
						"slash:", title, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
					return slashAmountTmp, balanceTmp, err
				}
			}

			slashAmountTmp = common.Big0
			balanceTmp = new(big.Int).Sub(can_balance, slashAmount)
		}

		return slashAmountTmp, balanceTmp, nil
	}

	/**
	Balance that can only be effective for Slash
	*/

	if slash.Cmp(common.Big0) > 0 && can.Released.Cmp(common.Big0) > 0 {
		val, rval, err := slashFunc("Released", slash, can.Released, false)
		if nil != err {
			return err
		}
		slash, can.Released = val, rval
	}

	if slash.Cmp(common.Big0) > 0 && can.RestrictingPlan.Cmp(common.Big0) > 0 {
		val, rval, err := slashFunc("RestrictingPlan", slash, can.RestrictingPlan, true)
		if nil != err {
			return err
		}
		slash, can.RestrictingPlan = val, rval
	}

	if slash.Cmp(common.Big0) != 0 {
		log.Error("Failed to SlashCandidates: the ramain is not zero", "slashed remain", slash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return common.BizErrorf("Failed to SlashCandidates: the slashed ramain is not zero, slashAmount:%d, slash remain:%d", amount, slash)
	}

	sharesHaveBeenClean := func() bool {
		return staking.Is_Invalid_LowRatio_NotEnough(can.Status) ||
			staking.Is_Invalid_LowRatioDel(can.Status) ||
			staking.Is_Invalid_DuplicateSign(can.Status) ||
			staking.Is_Invalid_Withdrew(can.Status)
	}

	// If the status has been modified before, this time it will not be modified.
	if sharesHaveBeenClean() {

		log.Info("Call SlashCandidates end, the candidate shares have been clean",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(),
			"can.Status", can.Status)

		// change status

		can.Status |= uint32(slashType)

		if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to SlashCandidates: Store candidate is failed", "slashType", slashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

	} else {
		// sub Shares to effect power
		if can.Shares.Cmp(amount) >= 0 {
			can.Shares = new(big.Int).Sub(can.Shares, amount)
		} else {
			log.Error("Failed to SlashCandidates: the candidate shares is no enough", "slashType", slashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "candidate shares",
				can.Shares, "slash amount", amount)
			panic("the candidate shares is no enough")
		}

		remainRelease := new(big.Int).Add(can.Released, can.ReleasedHes)
		remainRestrictingPlan := new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
		canRemain := new(big.Int).Add(remainRelease, remainRestrictingPlan)

		var needDelete bool

		switch slashType {
		case staking.LowRatio:
			if !xutil.CheckStakeThreshold(canRemain) {
				can.Status |= staking.NotEnough
				needDelete = true
			}
		case staking.LowRatioDel:
			needDelete = true
		case staking.DuplicateSign:
			needDelete = true
		}
		can.Status |= uint32(slashType)

		// the can status is valid and do not need delete
		// update the can power
		if !needDelete && !staking.Is_Invalid(can.Status) {

			// update the candidate power, If do not need to delete power (the candidate status still be valid)
			if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
				log.Error("Failed to SlashCandidates: Store candidate power is failed", "slashType", slashType,
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
				return err
			}

			if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
				log.Error("Failed to SlashCandidates: Store candidate is failed", "slashType", slashType,
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
				return err
			}

		} else {
			//because of deleted candidate info ,clean Shares
			can.Shares = common.Big0
			can.Status |= staking.Invalided

			// need to sub account rc
			if err := sk.db.SubAccountStakeRc(blockHash, can.StakingAddress); nil != err {
				log.Error("Failed to SlashCandidates: Sub Account staking Reference Count is failed", "slashType", slashType,
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
				return err
			}

			// withdrew Stake if the candidate status is invalid
			if can.ReleasedHes.Cmp(common.Big0) > 0 {

				state.AddBalance(can.StakingAddress, can.ReleasedHes)
				state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)
				can.ReleasedHes = common.Big0
			}
			if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {

				err := rt.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
				if nil != err {
					log.Error("Failed to SlashCandidates on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "stakingAddr", can.StakingAddress.Hex(),
						"restrictingPlanHes", can.RestrictingPlanHes, "err", err)
					return err
				}

				can.RestrictingPlanHes = common.Big0
			}

			if can.Released.Cmp(common.Big0) > 0 || can.RestrictingPlan.Cmp(common.Big0) > 0 {

				if err := sk.addUnStakeItem(state, blockNumber, blockHash, epoch, can.NodeId, canAddr); nil != err {
					log.Error("Failed to SlashCandidates on stakingPlugin: Add UnStakeItemStore failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
					return err
				}

				if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
					log.Error("Failed to SlashCandidates on stakingPlugin: Store Candidate info is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(), "err", err)
					return err
				}
			} else {

				// Clean candidate info
				if err := sk.db.DelCandidateStore(blockHash, canAddr); nil != err {
					log.Error("Failed to SlashCandidates on stakingPlugin: Delete Candidate info is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", can.NodeId.String(),
						"canAddr", canAddr.String(), "err", err)
					return err
				}
			}

			validators, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
			if nil != err {
				log.Error("Failed to SlashCandidates: Query Verifier List is failed", "slashType", slashType,
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
				return err
			}

			// remove the val from epoch validators,
			// because the candidate status is invalid after slashed
			orginLen := len(validators.Arr)
			for i := 0; i < len(validators.Arr); i++ {

				val := validators.Arr[i]

				if val.NodeId == nodeId {

					log.Info("Call SlashCandidates, Delete the validator", "slashType", slashType,
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())

					validators.Arr = append(validators.Arr[:i], validators.Arr[i+1:]...)
					i--
					break
				}
			}
			dirtyLen := len(validators.Arr)

			if dirtyLen != orginLen {

				if err := sk.setVerifierList(blockHash, validators); nil != err {
					log.Error("Failed to SlashCandidates: Store Verifier List is failed", "slashType", slashType,
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
					return err
				}
			}
		}

	}

	log.Info("Call SlashCandidates end ...")
	return nil
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
			log.Error("Call ProposalPassedNotify: Delete Candidate old power is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		can.ProgramVersion = version

		if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
			log.Error("Call ProposalPassedNotify: Store Candidate new power is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Call ProposalPassedNotify: Store Candidate info is failed", "blockNumber", blockNumber,
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
		log.Error("Call DeclarePromoteNotify: Query Candidate is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	if (nil != err && err == snapshotdb.ErrNotFound) || nil == can {

		log.Error("Call DeclarePromoteNotify: Promote candidate programVersion failed, the can is empty",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(),
			"version", programVersion)
		return nil
	}

	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Call DeclarePromoteNotify: Delete Candidate old power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	can.ProgramVersion = version

	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Call DeclarePromoteNotify: Store Candidate new power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Call DeclarePromoteNotify: Store Candidate info is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) GetLastNumber(blockNumber uint64) uint64 {

	val_arr, err := sk.getCurrValList(common.ZeroHash, blockNumber, QueryStartIrr)
	if nil != err {
		log.Error("Failed to GetLastNumber", "blockNumber", blockNumber, "err", err)
		return 0
	}

	if nil == err && nil != val_arr {
		return val_arr.End
	}
	return 0
}

func (sk *StakingPlugin) GetValidator(blockNumber uint64) (*cbfttypes.Validators, error) {

	val_arr, err := sk.getCurrValList(common.ZeroHash, blockNumber, QueryStartIrr)
	if nil != err && err != snapshotdb.ErrNotFound {
		return nil, err
	}

	if nil == err && nil != val_arr {
		return buildCbftValidators(val_arr.Start, val_arr.Arr), nil
	}
	return nil, common.BizErrorf("No Found Validators by blockNumber: %d", blockNumber)
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
		can.ReleasedHes = common.Big0
	}

	if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		can.RestrictingPlan = new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
		can.RestrictingPlanHes = common.Big0
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
		del.ReleasedHes = common.Big0
	}

	if del.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		del.RestrictingPlan = new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
		del.RestrictingPlanHes = common.Big0
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
		log.Error("vrfElection failed", "validatorListSize", len(validatorList),
			"nonceSize", len(nonce), "preNoncesSize", len(preNonces), "parentHash", hex.EncodeToString(parentHash.Bytes()))
		return nil, ParamsErr
	}
	if len(preNonces) > len(validatorList) {
		preNonces = preNonces[len(preNonces)-len(validatorList):]
	}
	return ProbabilityElection(validatorList, shiftLen, vrf.ProofToHash(nonce), preNonces)
}

func ProbabilityElection(validatorList staking.ValidatorQueue, shiftLen int, currentNonce []byte, preNonces [][]byte) (staking.ValidatorQueue, error) {
	if len(currentNonce) == 0 || len(preNonces) == 0 || len(validatorList) != len(preNonces) {
		log.Error("probabilityElection failed", "validatorListSize", len(validatorList),
			"currentNonceSize", len(currentNonce), "preNoncesSize", len(preNonces), "EpochValidatorNum", xcom.EpochValidatorNum)
		return nil, ParamsErr
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
	p := float64(xcom.ShiftValidatorNum()) * float64(xcom.ConsValidatorNum()) / sumWeightsFloat
	log.Info("probabilityElection Basic parameter", "validatorListSize", len(validatorList),
		"p", p, "sumWeights", sumWeightsFloat, "shiftValidatorNum", shiftLen, "epochValidatorNum", xcom.EpochValidatorNum())
	for index, sv := range svList {
		resultStr := new(big.Int).Xor(new(big.Int).SetBytes(currentNonce), new(big.Int).SetBytes(preNonces[index])).Text(10)
		target, err := strconv.ParseFloat(resultStr, 64)
		if nil != err {
			return nil, err
		}
		targetP := target / maxValue
		bd := xcom.NewBinomialDistribution(sv.weights, p)
		x, err := bd.InverseCumulativeProbability(targetP)
		if nil != err {
			return nil, err
		}
		sv.x = x
		log.Debug("calculated probability", "nodeId", hex.EncodeToString(sv.v.NodeId.Bytes()),
			"addr", hex.EncodeToString(sv.v.NodeAddress.Bytes()), "index", index, "currentNonce",
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
		log.Debug("sort validator", "addr", hex.EncodeToString(sv.v.NodeAddress.Bytes()),
			"index", index, "weight", sv.weights, "x", sv.x, "version", sv.version,
			"blockNumber", sv.blockNumber, "txIndex", sv.txIndex)
	}
	return resultValidatorList, nil
}

/**
Internal expansion function
*/

// previous round validators
func (sk *StakingPlugin) getPreValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.Validator_array, error) {

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
		return nil, ValidatorNotExist
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
		return nil, ValidatorNotExist
	}

	return &staking.Validator_array{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) getCurrValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.Validator_array, error) {

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
		return nil, ValidatorNotExist
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
		return nil, ValidatorNotExist
	}

	return &staking.Validator_array{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) getNextValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.Validator_array, error) {

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
		return nil, ValidatorNotExist
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
		return nil, ValidatorNotExist
	}

	return &staking.Validator_array{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) setRoundValList(blockHash common.Hash, val_Arr *staking.Validator_array) error {

	queue, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to setRoundValList: Query round valIndex is failed", "blockHash",
			blockHash.Hex(), "Start", val_Arr.Start, "End", val_Arr.End, "err", err)
		return err
	}

	index := &staking.ValArrIndex{
		Start: val_Arr.Start,
		End:   val_Arr.End,
	}

	shabby, queue := queue.ConstantAppend(index, RoundValIndexSize)

	// delete the shabby validators
	if nil != shabby {
		if err := sk.db.DelRoundValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
			log.Error("Failed to setRoundValList: delete shabby validators is failed",
				"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}

	// Store new index Arr
	if err := sk.db.SetRoundValIndex(blockHash, queue); nil != err {
		log.Error("Failed to setRoundValList: store round validators new indexArr is failed",
			"blockHash", blockHash.Hex(), "indexs length", len(queue), "err", err)
		return err
	}

	// Store new round validator Item
	if err := sk.db.SetRoundValList(blockHash, index.Start, index.End, val_Arr.Arr); nil != err {
		log.Error("Failed to setRoundValList: store new round validators is failed",
			"blockHash", blockHash.Hex(), "start", index.Start, "end", index.End,
			"val arr length", len(val_Arr.Arr), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) getVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.Validator_array, error) {

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
		return nil, ValidatorNotExist
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
		return nil, ValidatorNotExist
	}

	return &staking.Validator_array{
		Start: targetIndex.Start,
		End:   targetIndex.End,
		Arr:   queue,
	}, nil
}

func (sk *StakingPlugin) setVerifierList(blockHash common.Hash, val_Arr *staking.Validator_array) error {

	queue, err := sk.db.GetEpochValIndexByBlockHash(blockHash)
	if nil != err {
		log.Error("Failed to setVerifierList: Query epoch valIndex is failed", "blockHash",
			blockHash.Hex(), "Start", val_Arr.Start, "End", val_Arr.End, "err", err)
		return err
	}

	index := &staking.ValArrIndex{
		Start: val_Arr.Start,
		End:   val_Arr.End,
	}

	shabby, queue := queue.ConstantAppend(index, EpochValIndexSize)

	// delete the shabby validators
	if nil != shabby {
		if err := sk.db.DelEpochValListByBlockHash(blockHash, shabby.Start, shabby.End); nil != err {
			log.Error("Failed to setVerifierList: delete shabby validators is failed",
				"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}

	// Store new index Arr
	if err := sk.db.SetEpochValIndex(blockHash, queue); nil != err {
		log.Error("Failed to setVerifierList: store epoch validators new indexArr is failed",
			"blockHash", blockHash.Hex(), "indexs length", len(queue), "err", err)
		return err
	}

	// Store new epoch validator Item
	if err := sk.db.SetEpochValList(blockHash, index.Start, index.End, val_Arr.Arr); nil != err {
		log.Error("Failed to setVerifierList: store new epoch validators is failed",
			"blockHash", blockHash.Hex(), "start", index.Start, "end", index.End,
			"val arr length", len(val_Arr.Arr), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) addUnStakeItem(state xcom.StateDB, blockNumber uint64, blockHash common.Hash, epoch uint64,
	nodeId discover.NodeID, canAddr common.Address) error {

	endVoteNum, err := gov.GetMaxEndVotingBlock(nodeId, blockHash, state)
	if nil != err {
		return err
	}

	refundEpoch := xutil.CalculateEpoch(blockNumber) + xcom.UnStakeFreezeRatio()
	maxEndVoteEpoch := xutil.CalculateEpoch(endVoteNum)

	var targetEpoch uint64

	if maxEndVoteEpoch <= refundEpoch {
		targetEpoch = refundEpoch
	} else {
		targetEpoch = maxEndVoteEpoch
	}

	log.Info("Call addUnStakeItem, AddUnStakeItemStore start", "current blockNumber", blockNumber,
		"govenance max end vote blokNumber", endVoteNum, "unStakeFreeze Epoch", refundEpoch,
		"govenance max end vote epoch", maxEndVoteEpoch, "unstake item target Epoch", targetEpoch,
		"nodeId", nodeId.String())

	if err := sk.db.AddUnStakeItemStore(blockHash, targetEpoch, canAddr); nil != err {
		return err
	}
	return nil
}

func (sk *StakingPlugin) addUnDelegateItem(blockNumber uint64, blockHash common.Hash, delAddr common.Address, nodeId discover.NodeID,
	epoch, stakeBlockNumber uint64, amount *big.Int) error {

	targetEpoch := epoch + xcom.ActiveUnDelFreezeRatio()

	log.Info("Call addUnDelegateItem, AddUnDelegateItemStore start", "current blockNumber",
		blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "current epoch", epoch,
		"target epoch", targetEpoch, "stakingBlockNum", stakeBlockNumber, "refundAmount", amount)

	if err := sk.db.AddUnDelegateItemStore(blockHash, delAddr, nodeId, targetEpoch, stakeBlockNumber, amount); nil != err {
		return err
	}
	return nil
}

func (sk *StakingPlugin) HasStake(blockHash common.Hash, addr common.Address) (bool, error) {
	return sk.db.HasAccountStakeRc(blockHash, addr)
}

func calCanTotalAmount(can *staking.Candidate) *big.Int {
	remainRelease := new(big.Int).Add(can.Released, can.ReleasedHes)
	remainRestrictingPlan := new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
	return new(big.Int).Add(remainRelease, remainRestrictingPlan)
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
