package plugin

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"sort"
	"strconv"
	"sync"

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
	once     sync.Once
}

var (
	stakePlnOnce sync.Once
	stk          *StakingPlugin
)

var (
	AccountVonNotEnough        = common.NewBizError("The von of account is not enough")
	DelegateVonNotEnough       = common.NewBizError("The von of delegate is not enough")
	WithdrewDelegateVonCalcErr = common.NewBizError("withdrew delegate von calculate err")
	ParamsErr                  = common.NewBizError("the fn params err")
	BlockNumberDisordered      = common.NewBizError("The blockNumber is disordered")
	VonAmountNotRight          = common.NewBizError("The amount of von is not right")
	CandidateNotExist          = common.NewBizError("The candidate is not exist")
	ValidatorNotExist          = common.NewBizError("The validator is not exist")
)

const (
	FreeOrigin            = 0
	RestrictingPlanOrigin = 1

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

	log.Info("Call EndBlock on staking plugin", "blockNumber", header.Number, "blockHash", blockHash.String(), "epoch", epoch)

	if xutil.IsSettlementPeriod(header.Number.Uint64()) {
		// handle UnStaking Item
		err := sk.HandleUnCandidateItem(state, blockHash, epoch)
		if nil != err {
			log.Error("Failed to call HandleUnCandidateItem on stakingPlugin EndBlock", "blockHash",
				blockHash.Hex(), "blockNumber", header.Number.Uint64(), "err", err)
			return err //  TODO common.NewSysError(err.Error())
		}

		// hanlde UnDelegate Item
		err = sk.HandleUnDelegateItem(state, blockHash, epoch)
		if nil != err {
			log.Error("Failed to call HandleUnDelegateItem on stakingPlugin EndBlock", "blockHash",
				blockHash.Hex(), "blockNumber", header.Number.Uint64(), "err", err)
			return err
		}

		// Election next epoch validators
		if err := sk.ElectNextVerifierList(blockHash, header.Number.Uint64(), state); nil != err {
			log.Error("Failed to call ElectNextVerifierList on stakingPlugin EndBlock", "blockHash",
				blockHash.Hex(), "blockNumber", header.Number.Uint64(), "err", err)
			return err
		}
	}

	if xutil.IsElection(header.Number.Uint64()) {
		// ELection next round validators
		err := sk.Election(blockHash, header)
		if nil != err {
			log.Error("Failed to call Election on stakingPlugin EndBlock", "blockHash", blockHash.Hex(),
				"blockNumber", header.Number.Uint64(), "err", err)
			return err
		}

		// todo test
		pposHash := snapshotdb.Instance().GetLastKVHash(blockHash)
		log.Info("Query ppos hash, After Election", "blockHash", header.Hash().Hex(), "blockNumber", header.Number.Uint64(), "pposHash", hex.EncodeToString(pposHash))

	}
	log.Info("Finished EndBlock on staking plugin", "blockNumber", header.Number, "blockHash", blockHash.String(), "epoch", epoch)
	return nil
}

func (sk *StakingPlugin) Confirmed(block *types.Block) error {

	log.Info("Call Confirmed on staking plugin", "blockNumber", block.Number(), "blockHash", block.Hash().String())

	if xutil.IsElection(block.NumberU64()) {

		next, err := sk.getNextValList(block.Hash(), block.NumberU64(), QueryStartNotIrr)
		if nil != err {
			log.Error("Failed to Query Next validators on stakingPlugin Confirmed When Election block",
				"blockHash", block.Hash().Hex(), "blockNumber", block.Number().Uint64(), "err", err)
			return err
		}

		current, err := sk.getCurrValList(block.Hash(), block.NumberU64(), QueryStartNotIrr)
		if nil != err {
			log.Error("Failed to Query Current Round validators on stakingPlugin Confirmed When Election block",
				"blockHash", block.Hash().Hex(), "blockNumber", block.Number().Uint64(), "err", err)
			return err
		}
		result := distinct(next.Arr, current.Arr)
		if len(result) > 0 {
			sk.addConsensusNode(result)
			log.Debug("stakingPlugin addConsensusNode success", "blockNumber", block.NumberU64(), "size", len(result))
		}
	}

	if xutil.IsSwitch(block.NumberU64()) {
		pre, err := sk.getPreValList(block.Hash(), block.NumberU64(), QueryStartNotIrr)
		if nil != err {
			log.Error("Failed to Query Previous Round validators on stakingPlugin Confirmed When Switch block",
				"blockHash", block.Hash().Hex(), "blockNumber", block.Number().Uint64(), "err", err)
			return err
		}
		current, err := sk.getCurrValList(block.Hash(), block.NumberU64(), QueryStartNotIrr)
		if nil != err {
			log.Error("Failed to Query Current Round validators on stakingPlugin Confirmed When Switch block",
				"blockHash", block.Hash().Hex(), "blockNumber", block.Number().Uint64(), "err", err)
			return err
		}
		result := distinct(pre.Arr, current.Arr)
		if len(result) > 0 {
			sk.removeConsensusNode(result)
			log.Debug("stakingPlugin removeConsensusNode success", "blockNumber", block.NumberU64(), "size", len(result))
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

func (sk *StakingPlugin) removeConsensusNode(nodes staking.ValidatorQueue) {
	for _, node := range nodes {
		sk.eventMux.Post(cbfttypes.RemoveValidatorEvent{NodeID: node.NodeId})
	}
}

func (sk *StakingPlugin) GetCandidateInfo(blockHash common.Hash, addr common.Address) (*staking.Candidate, error) {
	return sk.db.GetCandidateStore(blockHash, addr)
}

func (sk *StakingPlugin) GetCandidateCompactInfo(blockHash common.Hash, blockNumber uint64, addr common.Address) (*staking.Candidate, error) {
	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if nil != err {
		return nil, err
	}

	epoch := xutil.CalculateEpoch(blockNumber)

	lazyCalcStakeAmount(epoch, can)

	return can, nil
}

func (sk *StakingPlugin) GetCandidateInfoByIrr(addr common.Address) (*staking.Candidate, error) {
	return sk.db.GetCandidateStoreByIrr(addr)
}

func (sk *StakingPlugin) CreateCandidate(state xcom.StateDB, blockHash common.Hash, blockNumber,
	amount *big.Int, typ uint16, addr common.Address, can *staking.Candidate) error {

	log.Debug("Call CreateCandidate", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
		"nodeId", can.NodeId.String())

	// todo test
	xcom.PrintObject("CreateCandidate, Method Start can", can)

	// from account free von
	if typ == FreeOrigin {

		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to CreateCandidate on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "addr", addr.String(),
				"originVon", origin, "stakingVon", amount)
			return AccountVonNotEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)
		can.ReleasedHes = amount

	} else if typ == RestrictingPlanOrigin { //  from account RestrictingPlan von

		err := rt.PledgeLockFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to CreateCandidate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "addr", addr.String(),
				"stakingVon", amount, "err", err)
			return err
		}
		can.RestrictingPlanHes = amount
	}

	can.StakingEpoch = uint32(xutil.CalculateEpoch(blockNumber.Uint64()))

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("CreateCandidate, Method Before SetCandidateStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("CreateCandidate, Method Before SetCandidateStore", can)

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Store Candidate info is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "addr", addr.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("CreateCandidate, Method Before SetCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("CreateCandidate, Method Before SetCanPowerStore", can)

	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Store Candidate power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "addr", addr.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("CreateCandidate, Method Before AddAccountStakeRc", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("CreateCandidate, Method Before AddAccountStakeRc", can)

	// add the account staking Reference Count
	if err := sk.db.AddAccountStakeRc(blockHash, can.StakingAddress); nil != err {
		log.Error("Failed to CreateCandidate on stakingPlugin: Store Staking Account Reference Count (add) is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "canAddr", addr.String(),
			"staking Account", can.StakingAddress.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("CreateCandidate, Method end", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	return nil
}

/// This method may only be called when creatStaking
func (sk *StakingPlugin) RollBackStaking(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	addr common.Address, typ uint16) error {

	log.Debug("Call RollBackStaking", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "addr", addr.String())

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) <= 0 {
		log.Error("Failed to RollBackStaking: the balance is invalid of stakingContracr Account", "contract_balance", contract_balance)
		panic("the balance is invalid of stakingContracr Account")
	}

	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if nil != err {
		return err
	}

	// todo test
	xcom.PrintObject("RollBackStaking, Method Start can", can)

	if blockNumber.Uint64() != can.StakingBlockNum {
		return common.BizErrorf("%v: current blockNumber is not equal stakingBlockNumber, can not rollback staking ...", ParamsErr)
	}

	// RollBack Staking

	if typ == FreeOrigin {

		state.AddBalance(can.StakingAddress, can.ReleasedHes)
		state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)

	} else if typ == RestrictingPlanOrigin {

		err := rt.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
		if nil != err {
			log.Error("Failed to RollBackStaking on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "addr", addr.String(),
				"RollBack stakingVon", can.RestrictingPlanHes, "err", err)
			return err
		}
	}

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("RollBackStaking, Method Before DelCandidateStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("RollBackStaking, Method Before DelCandidateStore", can)

	if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
		log.Error("Failed to RollBackStaking on stakingPlugin: Delete Candidate info is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "addr", addr.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("RollBackStaking, Method Before DelCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("RollBackStaking, Method Before DelCanPowerStore", can)

	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to RollBackStaking on stakingPlugin: Delete Candidate power failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "addr", addr.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("RollBackStaking, Method Before SubAccountStakeRc", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("RollBackStaking, Method Before SubAccountStakeRc", can)

	// sub the account staking Reference Count
	if err := sk.db.SubAccountStakeRc(blockHash, can.StakingAddress); nil != err {
		log.Error("Failed to RollBackStaking on stakingPlugin: Store Staking Account Reference Count (sub) is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "canAddr", addr.String(),
			"staking Account", can.StakingAddress.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("RollBackStaking, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	return nil
}

func (sk *StakingPlugin) EditCandidate(blockHash common.Hash, blockNumber *big.Int, can *staking.Candidate) error {

	log.Debug("Call EditCandidate", "blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(),
		"nodeId", can.NodeId.String())

	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("EditCandidate, Method Before SetCandidateStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("EditCandidate, Method Before SetCandidateStore", can)

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to EditCandidate on stakingPlugin: Store Candidate info is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("EditCandidate, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	return nil
}

func (sk *StakingPlugin) IncreaseStaking(state xcom.StateDB, blockHash common.Hash, blockNumber,
	amount *big.Int, typ uint16, can *staking.Candidate) error {

	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	log.Debug("Call IncreaseStaking", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"epoch", epoch, "nodeId", can.NodeId.String(), "typ", typ, "amount", amount)

	// todo test
	xcom.PrintObject("IncreaseStaking, Method Start can", can)

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	if typ == FreeOrigin {
		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(amount) < 0 {
			log.Error("Failed to IncreaseStaking on stakingPlugin: the account free von is not Enough",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "account", can.StakingAddress.Hex(),
				"originVon", origin, "stakingVon", can.ReleasedHes)
			return AccountVonNotEnough
		}
		state.SubBalance(can.StakingAddress, amount)
		state.AddBalance(vm.StakingContractAddr, amount)

		can.ReleasedHes = new(big.Int).Add(can.ReleasedHes, amount)

	} else {

		err := rt.PledgeLockFunds(can.StakingAddress, amount, state)
		if nil != err {
			log.Error("Failed to IncreaseStaking on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		can.RestrictingPlanHes = new(big.Int).Add(can.RestrictingPlanHes, amount)
	}

	can.StakingEpoch = uint32(epoch)

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("IncreaseStaking pposHash, Before DelCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("IncreaseStaking, Method Before DelCanPowerStore", can)

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to IncreaseStaking on stakingPlugin: Delete Candidate old power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	can.Shares = new(big.Int).Add(can.Shares, amount)

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("IncreaseStaking, Method Before SetCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("IncreaseStaking, Method Before SetCanPowerStore", can)

	// set new power of can
	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Failed to IncreaseStaking on stakingPlugin: Store Candidate new power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("IncreaseStaking, Method Before SetCandidateStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("IncreaseStaking, Method Before SetCandidateStore", can)

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to IncreaseStaking on stakingPlugin: Store Candidate info is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("IncreaseStaking pposHash, End SetCandidateStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	return nil
}

func (sk *StakingPlugin) WithdrewStaking(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	can *staking.Candidate) error {
	pubKey, _ := can.NodeId.Pubkey()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())

	log.Debug("Call WithdrewStaking", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"epoch", epoch, "nodeId", can.NodeId.String())

	lazyCalcStakeAmount(epoch, can)

	addr := crypto.PubkeyToAddress(*pubKey)

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("WithdrewStaking, Method Before DelCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("WithdrewStaking, Method Before DelCanPowerStore", can)

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to WithdrewStaking on stakingPlugin: Delete Candidate old power is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	if err := sk.withdrewStakeAmount(state, blockHash, blockNumber.Uint64(), epoch, addr, can); nil != err {
		return err
	}

	can.StakingEpoch = uint32(epoch)

	if can.Released.Cmp(common.Big0) > 0 || can.RestrictingPlan.Cmp(common.Big0) > 0 {

		// TODO test
		pposHash := sk.db.GetLastKVHash(blockHash)
		log.Debug("WithdrewStaking, Method Before SetCandidateStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("WithdrewStaking, Method Before SetCandidateStore", can)

		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: Store Candidate info is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	} else {

		// TODO test
		pposHash := sk.db.GetLastKVHash(blockHash)
		log.Debug("WithdrewStaking, Method Before DelCandidateStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("WithdrewStaking, Method Before DelCandidateStore", can)

		// Clean candidate info
		if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: Delete Candidate info is failed",
				"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("WithdrewStaking, Method Before SubAccountStakeRc", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("WithdrewStaking, Method Before SubAccountStakeRc", can)

	// sub the account staking Reference Count
	if err := sk.db.SubAccountStakeRc(blockHash, can.StakingAddress); nil != err {
		log.Error("Failed to WithdrewStaking on stakingPlugin: Store Staking Account Reference Count (sub) is failed",
			"blockNumber", blockNumber.Uint64(), "blockHash", blockHash.Hex(), "canAddr", addr.String(),
			"staking Account", can.StakingAddress.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("WithdrewStaking, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	return nil
}

func (sk *StakingPlugin) withdrewStakeAmount(state xcom.StateDB, blockHash common.Hash, blockNumber, epoch uint64,
	addr common.Address, can *staking.Candidate) error {

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) <= 0 {
		log.Error("Failed to withdrewStakeAmount: the balance is invalid of stakingContracr Account", "contract_balance", contract_balance)
		panic("the balance is invalid of stakingContracr Account")
	}

	// Direct return of money during the hesitation period
	// Return according to the way of coming
	if can.ReleasedHes.Cmp(common.Big0) > 0 {

		state.AddBalance(can.StakingAddress, can.ReleasedHes)
		state.SubBalance(vm.StakingContractAddr, can.ReleasedHes)
		//can.Shares = new(big.Int).Sub(can.Shares, can.ReleasedHes)
		can.ReleasedHes = common.Big0
	}

	if can.RestrictingPlanHes.Cmp(common.Big0) > 0 {

		err := rt.ReturnLockFunds(can.StakingAddress, can.RestrictingPlanHes, state)
		if nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: call Restricting ReturnLockFunds() is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		//can.Shares = new(big.Int).Sub(can.Shares, can.RestrictingPlanHes)
		can.RestrictingPlanHes = common.Big0
	}
	//addItem := false
	//
	//if can.Released.Cmp(common.Big0) > 0 {
	//	can.Shares = new(big.Int).Sub(can.Shares, can.Released)
	//	addItem = true
	//}
	//
	//if can.RestrictingPlan.Cmp(common.Big0) > 0 {
	//	can.Shares = new(big.Int).Sub(can.Shares, can.RestrictingPlan)
	//	addItem = true
	//}
	if can.Released.Cmp(common.Big0) > 0 || can.RestrictingPlan.Cmp(common.Big0) > 0 {

		// TODO test
		pposHash := sk.db.GetLastKVHash(blockHash)
		log.Debug("WithdrewStaking, Method Before AddUnStakeItemStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("WithdrewStaking, Method Before AddUnStakeItemStore", can)

		if err := sk.db.AddUnStakeItemStore(blockHash, epoch, addr); nil != err {
			log.Error("Failed to WithdrewStaking on stakingPlugin: Add UnStakeItemStore failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}
	}

	can.Shares = common.Big0
	can.Status |= staking.Invalided

	return nil
}

func (sk *StakingPlugin) HandleUnCandidateItem(state xcom.StateDB, blockHash common.Hash, epoch uint64) error {

	log.Info("Call HandleUnCandidateItem", "blockHash", blockHash.Hex(), "epoch", epoch)

	releaseEpoch := epoch - xcom.UnStakeFreezeRatio()

	unStakeCount, err := sk.db.GetUnStakeCountStore(blockHash, releaseEpoch)
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
		addr, err := sk.db.GetUnStakeItemStore(blockHash, releaseEpoch, uint64(index))
		if nil != err {
			log.Error("Failed to HandleUnCandidateItem: Query the unStakeItem node addr is failed",
				"blockHash", blockHash.Hex(), "err", err)
			return err
		}

		if _, ok := filterAddr[addr]; ok {
			if err := sk.db.DelUnStakeItemStore(blockHash, epoch, uint64(index)); nil != err {
				log.Error("Failed to HandleUnCandidateItem: Delete already handle unstakeItem failed",
					"blockHash", blockHash.Hex(), "err", err)
				return err
			}
			continue
		}

		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			log.Error("Failed to HandleUnCandidateItem: Query candidate failed", "blockHash", blockHash.Hex(), "err", err)
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
			log.Error("Failed to HandleUnCandidateItem: Delete unstakeItem failed", "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		filterAddr[addr] = struct{}{}
	}

	if err := sk.db.DelUnStakeCountStore(blockHash, releaseEpoch); nil != err {
		log.Error("Failed to HandleUnCandidateItem: Delete unstakeCount failed", "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	return nil
}

func (sk *StakingPlugin) handleUnStake(state xcom.StateDB, blockHash common.Hash, epoch uint64,
	addr common.Address, can *staking.Candidate) error {

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) <= 0 {
		log.Error("Failed to handleUnStake: the balance is invalid of stakingContracr Account", "contract_balance", contract_balance)
		panic("the balance is invalid of stakingContracr Account")
	}

	lazyCalcStakeAmount(epoch, can)

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
					title, balance, "blockHash", blockHash.Hex(), "err", err)
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

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("handleUnStake, Method Before DelCandidateStore",
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("handleUnStake, Method Before DelCandidateStore", can)

	// delete can info
	if err := sk.db.DelCandidateStore(blockHash, addr); nil != err {
		log.Error("Failed to HandleUnCandidateItem: Delete candidate info failed", "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("handleUnStake, Method End",
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

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
		Delegation:      *del,
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
		Delegation:      *del,
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
		Delegation:      *del,
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
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "originVon", origin,
				"stakingVon", can.ReleasedHes)
			return AccountVonNotEnough
		}
		state.SubBalance(delAddr, amount)
		state.AddBalance(vm.StakingContractAddr, amount)

		del.ReleasedHes = new(big.Int).Add(del.ReleasedHes, amount)

	} else if typ == RestrictingPlanOrigin { //  from account RestrictingPlan von

		err := rt.PledgeLockFunds(delAddr, amount, state)
		if nil != err {
			log.Error("Failed to Delegate on stakingPlugin: call Restricting PledgeLockFunds() is failed",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
			return err
		}

		del.RestrictingPlanHes = new(big.Int).Add(del.RestrictingPlanHes, amount)

	}

	del.DelegateEpoch = uint32(epoch)

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("Delegate, Method Before SetDelegateStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("Delegate, Method Before SetDelegateStore, del", del)

	// set new delegate info
	if err := sk.db.SetDelegateStore(blockHash, delAddr, can.NodeId, can.StakingBlockNum, del); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Store Delegate info is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("Delegate, Method Before DelCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("Delegate, Method Before DelCanPowerStore, can", can)

	// delete old power of can
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Delete Candidate old power is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// add the candidate power
	can.Shares = new(big.Int).Add(can.Shares, amount)

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("Delegate, Method Before SetCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("Delegate, Method Before SetCanPowerStore, can", can)

	// set new power of can
	if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Store Candidate new power is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("Delegate, Method Before SetCandidateStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("Delegate, Method Before SetCandidateStore, can", can)

	// update can info about Shares
	if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
		log.Error("Failed to Delegate on stakingPlugin: Store Candidate info is failed",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("Delegate, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	return nil
}

func (sk *StakingPlugin) WithdrewDelegate(state xcom.StateDB, blockHash common.Hash, blockNumber, amount *big.Int,
	delAddr common.Address, nodeId discover.NodeID, stakingBlockNum uint64, del *staking.Delegation) error {

	log.Debug("Call WithdrewDelegate", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"delAddr", delAddr.String(), "nodeId", nodeId.String(), "StakingNum", stakingBlockNum, "amount", amount)
	// todo test
	xcom.PrintObject("Call WithdrewDelegate, the delegate info", del)

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) <= 0 {
		log.Error("Failed to WithdrewDelegate: the balance is invalid of stakingContracr Account", "contract_balance", contract_balance)
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

	/**
	inner Fn
	*/
	subDelegateFn := func(source, sub *big.Int) (*big.Int, *big.Int) {

		state.AddBalance(delAddr, sub)
		state.SubBalance(vm.StakingContractAddr, sub)
		return new(big.Int).Sub(source, sub), common.Big0
	}

	refundFn := func(remain, aboutRelease, aboutRestrictingPlan *big.Int) (*big.Int, *big.Int, *big.Int, error) {

		remainTmp := remain
		releaseTmp := aboutRelease
		restrictingPlanTmp := aboutRestrictingPlan

		// When remain is greater than or equal to del.ReleasedHes/del.Released
		if remainTmp.Cmp(common.Big0) > 0 {
			if remainTmp.Cmp(releaseTmp) >= 0 && releaseTmp.Cmp(common.Big0) > 0 {

				remainTmp, releaseTmp = subDelegateFn(remainTmp, releaseTmp)

			} else if remainTmp.Cmp(releaseTmp) < 0 {
				// When remain is less than or equal to del.ReleasedHes/del.Released
				releaseTmp, remainTmp = subDelegateFn(releaseTmp, remainTmp)
			}
		}

		if remainTmp.Cmp(common.Big0) > 0 {

			// When remain is greater than or equal to del.RestrictingPlanHes/del.RestrictingPlan
			if remainTmp.Cmp(restrictingPlanTmp) >= 0 && restrictingPlanTmp.Cmp(common.Big0) > 0 {

				err := rt.ReturnLockFunds(can.StakingAddress, restrictingPlanTmp, state)
				if nil != err {
					log.Error("Failed to WithdrewDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
						"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
					return remainTmp, releaseTmp, restrictingPlanTmp, err
				}

				remainTmp = new(big.Int).Sub(remainTmp, restrictingPlanTmp)
				restrictingPlanTmp = common.Big0

			} else if remainTmp.Cmp(restrictingPlanTmp) < 0 {
				// When remain is less than or equal to del.RestrictingPlanHes/del.RestrictingPlan

				err := rt.ReturnLockFunds(can.StakingAddress, remainTmp, state)
				if nil != err {
					log.Error("Failed to WithdrewDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
						"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
					return remainTmp, releaseTmp, restrictingPlanTmp, err
				}

				restrictingPlanTmp = new(big.Int).Sub(restrictingPlanTmp, remainTmp)
				remainTmp = common.Big0
			}
		}

		return remainTmp, releaseTmp, restrictingPlanTmp, nil
	}

	del.DelegateEpoch = uint32(epoch)

	switch {

	/**
	**
	**
	When the related candidate info does not exist
	**
	**
	*/
	case nil == can, nil != can && stakingBlockNum < can.StakingBlockNum,
		nil != can && stakingBlockNum == can.StakingBlockNum && staking.Is_Invalid(can.Status):

		if total.Cmp(amount) < 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: the amount of invalid delegate is not enough",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
				"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "delegate amount", total,
				"withdrew amount", amount)
			return common.BizErrorf("withdrewDelegate err: %s, delegate von: %s, withdrew von: %s",
				DelegateVonNotEnough.Error(), total.String(), amount.String())
		}

		remain := common.Big0
		sub := new(big.Int).Sub(total, amount)

		// When the sub less than threshold
		if !xutil.CheckMinimumThreshold(sub) {
			remain = total
		} else {
			remain = amount
		}

		realSub := remain

		/**
		handle delegate on Hesitate period
		*/
		rm, rbalance, lbalance, err := refundFn(remain, del.ReleasedHes, del.RestrictingPlanHes)
		if nil != err {
			return err
		}
		remain, del.ReleasedHes, del.RestrictingPlanHes = rm, rbalance, lbalance
		/**
		handle delegate on Effective period
		*/
		if remain.Cmp(common.Big0) > 0 {
			rm, rbalance, lbalance, err = refundFn(remain, del.Released, del.RestrictingPlan)
			if nil != err {
				return err
			}
			remain, del.Released, del.RestrictingPlan = rm, rbalance, lbalance
		}

		if remain.Cmp(common.Big0) != 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: the ramain is not zero",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
				"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "del balance", total,
				"withdrew balance", amount, "realSub", realSub, "withdrewdiff", remain)
			return WithdrewDelegateVonCalcErr
		}

		if total.Cmp(realSub) == 0 {

			// TODO test
			pposHash := sk.db.GetLastKVHash(blockHash)
			log.Debug("WithdrewDelegate, Method Before DelDelegateStore", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

			// todo test
			xcom.PrintObject("WithdrewDelegate, Method Before DelDelegateStore, del", del)

			// clean delegate
			if err := sk.db.DelDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Delete detegate is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
					"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
				return err
			}

		} else {
			sub := new(big.Int).Sub(total, del.Reduction)

			if sub.Cmp(realSub) < 0 {
				diff := new(big.Int).Sub(realSub, sub)
				del.Reduction = new(big.Int).Sub(del.Reduction, diff)
			}

			// TODO test
			pposHash := sk.db.GetLastKVHash(blockHash)
			log.Debug("WithdrewDelegate, Method Before SetDelegateStore", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

			// todo test
			xcom.PrintObject("WithdrewDelegate, Method Before SetDelegateStore, del", del)

			if err := sk.db.SetDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum, del); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Store detegate is failed",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
					"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "err", err)
				return err
			}
		}

	/**
	**
	**
	Illegal parameter
	**
	**
	*/
	case nil != can && stakingBlockNum > can.StakingBlockNum:
		log.Error("Failed to WithdrewDelegate on stakingPlugin: the stakeBlockNum invalid",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
			"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "fn.stakeBlockNum", stakingBlockNum,
			"can.stakeBlockNum", can.StakingBlockNum)
		return ParamsErr

	/**
	**
	**
	When the delegate is normally revoked
	**
	**
	*/
	case nil != can && stakingBlockNum == can.StakingBlockNum && !staking.Is_Invalid(can.Status):

		// First need to deduct the von that is being refunded
		realtotal := new(big.Int).Sub(total, del.Reduction)

		if realtotal.Cmp(amount) < 0 {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: the amount of valid delegate is not enough",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(),
				"nodeId", nodeId.String(), "stakingBlockNum", stakingBlockNum, "delegate amount", realtotal,
				"withdrew amount", amount)
			return common.BizErrorf("withdrewDelegate err: %s, delegate von: %s, withdrew von: %s",
				DelegateVonNotEnough.Error(), realtotal.String(), amount.String())
		}

		remain := common.Big0
		realSub := common.Big0
		sub := new(big.Int).Sub(realtotal, amount)

		// When the sub less than threshold
		if !xutil.CheckMinimumThreshold(sub) {
			remain = realtotal
		} else {
			remain = amount
		}

		realSub = remain
		/**
		handle delegate on Hesitate period
		*/
		//var flag bool
		//var er error
		rm, rbalance, lbalance, err := refundFn(remain, del.ReleasedHes, del.RestrictingPlanHes)
		if nil != err {
			return err
		}
		remain, del.ReleasedHes, del.RestrictingPlanHes = rm, rbalance, lbalance

		save_or_del := false // false: save, true: delete

		/**
		handle delegate on Effective period
		*/
		if remain.Cmp(common.Big0) > 0 {

			// TODO test
			pposHash := sk.db.GetLastKVHash(blockHash)
			log.Debug("WithdrewDelegate, Method Before AddUnDelegateItemStore", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

			// todo test
			xcom.PrintObject("WithdrewDelegate, Method Before AddUnDelegateItemStore, del", del)

			// add a UnDelegateItem
			sk.db.AddUnDelegateItemStore(blockHash, delAddr, nodeId, epoch, stakingBlockNum, remain)
			del.Reduction = new(big.Int).Add(del.Reduction, remain)

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

			// TODO test
			pposHash := sk.db.GetLastKVHash(blockHash)
			log.Debug("WithdrewDelegate, Method Before SetDelegateStore", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

			// todo test
			xcom.PrintObject("WithdrewDelegate, Method Before SetDelegateStore, del", del)

			if err := sk.db.SetDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum, del); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Store delegate info is failed", "blockNumber",
					blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
					"stakingBlockNum", stakingBlockNum, "err", err)
				return err
			}
		} else {

			// TODO test
			pposHash := sk.db.GetLastKVHash(blockHash)
			log.Debug("WithdrewDelegate, Method Before DelDelegateStore", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

			// todo test
			xcom.PrintObject("WithdrewDelegate, Method Before DelDelegateStore, del", del)

			// Clean delegate info
			if err := sk.db.DelDelegateStore(blockHash, delAddr, nodeId, stakingBlockNum); nil != err {
				log.Error("Failed to WithdrewDelegate on stakingPlugin: Delete delegate info is failed", "blockNumber",
					blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
					"stakingBlockNum", stakingBlockNum, "err", err)
				return err
			}
		}

		// TODO test
		pposHash := sk.db.GetLastKVHash(blockHash)
		log.Debug("WithdrewDelegate, Method Before DelCanPowerStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("WithdrewDelegate, Method Before DelCanPowerStore, can", can)

		if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: Delete candidate old power is failed", "blockNumber",
				blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
				"stakingBlockNum", stakingBlockNum, "err", err)
			return err
		}

		/**
		**
		change candidate shares
		**
		*/
		can.Shares = new(big.Int).Sub(can.Shares, realSub)

		// TODO test
		pposHash = sk.db.GetLastKVHash(blockHash)
		log.Debug("WithdrewDelegate, Method Before SetCandidateStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("WithdrewDelegate, Method Before SetCandidateStore, can", can)

		if err := sk.db.SetCandidateStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: Store candidate info is failed", "blockNumber",
				blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
				"stakingBlockNum", stakingBlockNum, "err", err)
			return err
		}

		// TODO test
		pposHash = sk.db.GetLastKVHash(blockHash)
		log.Debug("WithdrewDelegate, Method Before SetCanPowerStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("WithdrewDelegate, Method Before SetCanPowerStore, can", can)

		if err := sk.db.SetCanPowerStore(blockHash, canAddr, can); nil != err {
			log.Error("Failed to WithdrewDelegate on stakingPlugin: Store candidate old power is failed", "blockNumber",
				blockNumber, "blockHash", blockHash.Hex(), "delAddr", delAddr.Hex(), "nodeId", nodeId.String(),
				"stakingBlockNum", stakingBlockNum, "err", err)
			return err
		}

	}

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("WithdrewDelegate, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	return nil
}

func (sk *StakingPlugin) HandleUnDelegateItem(state xcom.StateDB, blockHash common.Hash, epoch uint64) error {

	log.Info("Call HandleUnDelegateItem", "blockHash", blockHash.Hex(), "epoch", epoch)

	releaseEpoch := epoch - xcom.ActiveUnDelFreezeRatio()

	unDelegateCount, err := sk.db.GetUnDelegateCountStore(blockHash, releaseEpoch)
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
		unDelegateItem, err := sk.db.GetUnDelegateItemStore(blockHash, releaseEpoch, uint64(index))

		if nil != err {
			log.Error("Failed to HandleUnCandidateItem: Query the unStakeItem is failed", "blockHash",
				blockHash.Hex(), "epoch", epoch, "err", err)
			return err
		}

		del, err := sk.db.GetDelegateStoreBySuffix(blockHash, unDelegateItem.KeySuffix)
		if nil != err {
			log.Error("Failed to HandleUnCandidateItem: Query delegate info is failed", "blockHash",
				blockHash.Hex(), "epoch", epoch, "err", err)
			return err
		}

		if nil == del {
			// This maybe be nil
			continue
		}

		if err := sk.handleUnDelegate(state, blockHash, epoch, unDelegateItem, del); nil != err {
			return err
		}

	}

	return nil
}

func (sk *StakingPlugin) handleUnDelegate(state xcom.StateDB, blockHash common.Hash, epoch uint64,
	unDel *staking.UnDelegateItem, del *staking.Delegation) error {

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) <= 0 {
		log.Error("Failed to handleUnDelegate: the balance is invalid of stakingContracr Account", "contract_balance", contract_balance)
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

	amount := unDel.Amount

	aboutRelease := new(big.Int).Add(del.Released, del.ReleasedHes)
	aboutRestrictingPlan := new(big.Int).Add(del.RestrictingPlan, del.RestrictingPlanHes)
	total := new(big.Int).Add(aboutRelease, aboutRestrictingPlan)

	if amount.Cmp(del.Reduction) >= 0 && del.Reduction.Cmp(total) == 0 { // full withdrawal

		refundReleaseFn := func(balance *big.Int) *big.Int {
			if balance.Cmp(common.Big0) > 0 {

				state.AddBalance(delAddr, balance)
				state.SubBalance(vm.StakingContractAddr, balance)
				return common.Big0
			}
			return balance
		}

		del.ReleasedHes = refundReleaseFn(del.ReleasedHes)
		del.Released = refundReleaseFn(del.Released)

		refundRestrictingPlanFn := func(title string, balance *big.Int) (*big.Int, error) {

			if balance.Cmp(common.Big0) > 0 {
				err := rt.ReturnLockFunds(delAddr, balance, state)
				if nil != err {
					log.Error("Failed to handleUnDelegate on stakingPlugin: call Restricting ReturnLockFunds() is failed",
						title, balance, "blockHash", blockHash.Hex(), "epoch", epoch, "err", err)
					return common.Big0, err
				}
				return common.Big0, nil
			}

			return balance, nil
		}

		if balance, err := refundRestrictingPlanFn("RestrictingPlanHes", del.RestrictingPlanHes); nil != err {
			return err
		} else {
			del.RestrictingPlanHes = balance
		}

		if balance, err := refundRestrictingPlanFn("RestrictingPlanHes", del.RestrictingPlan); nil != err {
			return err
		} else {
			del.RestrictingPlan = balance
		}

		// TODO test
		pposHash := sk.db.GetLastKVHash(blockHash)
		log.Debug("handleUnDelegate, Method Before DelDelegateStoreBySuffix",
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("handleUnDelegate, Method Before DelDelegateStoreBySuffix, del", del)

		// clean the delegate
		if err := sk.db.DelDelegateStoreBySuffix(blockHash, unDel.KeySuffix); nil != err {
			log.Error("Failed to handleUnDelegate on stakingPlugin: Delete delegate info is failed",
				"blockHash", blockHash.Hex(), "epoch", epoch, "err", err)
			return err
		}

	} else { //few withdrawal

		remain := amount

		refundReleaseFn := func(balance, remain *big.Int) (*big.Int, *big.Int) {
			if remain.Cmp(common.Big0) > 0 {

				if remain.Cmp(balance) >= 0 {
					state.SubBalance(vm.StakingContractAddr, balance)
					state.AddBalance(delAddr, balance)
					return common.Big0, new(big.Int).Sub(remain, balance)
				} else {
					state.SubBalance(vm.StakingContractAddr, remain)
					state.AddBalance(delAddr, remain)
					return new(big.Int).Sub(balance, remain), common.Big0
				}
			}
			return balance, remain
		}

		hes, rm := refundReleaseFn(del.ReleasedHes, remain)
		del.ReleasedHes, remain = hes, rm
		noHes, rm := refundReleaseFn(del.Released, remain)
		del.Released, remain = noHes, rm

		refundRestrictingPlanFn := func(title string, balance, remain *big.Int) (*big.Int, *big.Int, error) {
			if remain.Cmp(common.Big0) > 0 {

				if remain.Cmp(balance) >= 0 {

					err := rt.ReturnLockFunds(delAddr, balance, state)
					if nil != err {
						log.Error("Failed to handleUnDelegate on stakingPlugin: call Restricting ReturnLockFunds() return "+title+" is failed",
							title, balance, "blockHash", blockHash.Hex(), "epoch", epoch, "err", err)
						return common.Big0, common.Big0, err
					}
					return common.Big0, new(big.Int).Sub(remain, balance), nil
				} else {

					err := rt.ReturnLockFunds(delAddr, remain, state)
					if nil != err {
						log.Error("Failed to handleUnDelegate on stakingPlugin: call Restricting ReturnLockFunds() return "+title+" is failed",
							"remain", remain, "blockHash", blockHash.Hex(), "epoch", epoch, "err", err)
						return common.Big0, common.Big0, err
					}

					return new(big.Int).Sub(balance, remain), common.Big0, nil
				}
			}

			return balance, remain, nil
		}

		if balance, re, err := refundRestrictingPlanFn("RestrictingPlanHes", del.RestrictingPlanHes, remain); nil != err {
			return err
		} else {
			del.RestrictingPlanHes, remain = balance, re
		}

		if balance, re, err := refundRestrictingPlanFn("RestrictingPlan", del.RestrictingPlan, remain); nil != err {
			return err
		} else {
			del.RestrictingPlan, remain = balance, re
		}

		if remain.Cmp(common.Big0) > 0 {
			log.Error("Failed to call handleUnDelegate: remain is not zero", "blockHash", blockHash.Hex(), "epoch", epoch,
				"delAddr", delAddr.Hex(), "nodeId", nodeId.String(), "stakeBlockNumber", num)
			return VonAmountNotRight
		}

		del.Reduction = new(big.Int).Sub(del.Reduction, amount)

		del.DelegateEpoch = uint32(epoch)

		// TODO test
		pposHash := sk.db.GetLastKVHash(blockHash)
		log.Debug("handleUnDelegate, Method Before SetDelegateStoreBySuffix",
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("handleUnDelegate, Method Before SetDelegateStoreBySuffix, del", del)

		if err := sk.db.SetDelegateStoreBySuffix(blockHash, unDel.KeySuffix, del); nil != err {
			log.Error("Failed to handleUnDelegate on stakingPlugin: Store delegate info is failed",
				"blockHash", blockHash.Hex(), "epoch", epoch, "err", err)
			return err
		}
	}

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("handleUnDelegate, Method End",
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

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

	/*if nil != old_verifierArr {

	}*/

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

	curr_version := govp.GetActiveVersion(state)
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
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
	}

	if len(queue) == 0 {
		panic(common.BizErrorf("Failed to ElectNextVerifierList: Select zero size validators~"))
	}

	new_verifierArr.Arr = queue

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("ElectNextVerifierList pposHash, Method Before setVerifierList", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("ElectNextVerifierList, Method Before setVerifierList, arr", new_verifierArr)

	err = sk.setVerifierList(blockHash, new_verifierArr)
	if nil != err {
		log.Error("Failed to ElectNextVerifierList: Set Next Epoch VerifierList is failed", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("ElectNextVerifierList pposHash, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

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
			StakingAddress:  can.StakingAddress,
			BenefitAddress:  can.BenefitAddress,
			StakingTxIndex:  can.StakingTxIndex,
			ProgramVersion:  can.ProgramVersion,
			StakingBlockNum: can.StakingBlockNum,
			Shares:          shares,
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
// 3:  	Query next round consensus validator
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
			StakingAddress:  can.StakingAddress,
			BenefitAddress:  can.BenefitAddress,
			StakingTxIndex:  can.StakingTxIndex,
			ProgramVersion:  can.ProgramVersion,
			StakingBlockNum: can.StakingBlockNum,
			Shares:          shares,
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

	//if blockNumber < arr.Start || blockNumber > arr.End {
	//	return nil, common.BizErrorf("Get Current ValidatorList failed: %s, start: %d, end: %d, currentNumer: %d",
	//		BlockNumberDisordered.Error(), arr.Start, arr.End, blockNumber)
	//}

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

func (sk *StakingPlugin) GetCandidateList(blockHash common.Hash, blockNumber uint64) (staking.CandidateQueue, error) {

	epoch := xutil.CalculateEpoch(blockNumber)

	iter := sk.db.IteratorCandidatePowerByBlockHash(blockHash, 0)
	if err := iter.Error(); nil != err {
		return nil, err
	}
	defer iter.Release()

	queue := make(staking.CandidateQueue, 0)

	for iter.Valid(); iter.Next(); {
		addrSuffix := iter.Value()
		can, err := sk.db.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			return nil, err
		}

		lazyCalcStakeAmount(epoch, can)

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

func (sk *StakingPlugin) Election(blockHash common.Hash, header *types.Header) error {

	log.Info("Call Election Start", "blockHash", blockHash.Hex(), "blockNumber", header.Number.Uint64())

	blockNumber := header.Number.Uint64()

	// the validators of Current Epoch
	verifiers, err := sk.getVerifierList(blockHash, blockNumber, QueryStartIrr)
	if nil != err {
		log.Error("Failed to call Election: No found current epoch validators", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return ValidatorNotExist
	}

	// the validators of Current Round
	curr, err := sk.getCurrValList(blockHash, blockNumber, QueryStartIrr)
	if nil != err {
		log.Error("Failed to Election: No found the current round validators", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return ValidatorNotExist
	}

	// todo test
	log.Info("Call Election start", "Curr epoch validators length", len(verifiers.Arr))
	log.Info("Call Election start", "Curr round validators length", len(curr.Arr))
	xcom.PrintObject("Call Election Curr validators", curr)

	if blockNumber != (curr.End - xcom.ElectionDistance()) {
		log.Error("Failed to Election: this blockNumber invalid", "Target blockNumber",
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

	currMap := make(map[discover.NodeID]struct{}, len(curr.Arr))
	for _, v := range curr.Arr {
		currMap[v.NodeId] = struct{}{}
	}

	// Exclude the current consensus round validators from the validators of the Epoch
	diffQueue := make(staking.ValidatorQueue, 0)
	for _, v := range verifiers.Arr {
		if _, ok := currMap[v.NodeId]; ok {
			continue
		}
		diffQueue = append(diffQueue, v)
	}

	mbn := 1 // Minimum allowed total number of consensus nodes
	diffQueueLen := len(diffQueue)
	invalidLen := 0 // duplicateSign And lowRatio No enough von
	currLen := len(curr.Arr)

	slashCans := make(staking.SlashCandidate, 0)
	slashAddrQueue := make([]discover.NodeID, 0)

	for _, v := range curr.Arr {

		addr, _ := xutil.NodeId2Addr(v.NodeId)
		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			log.Error("Failed to Get Candidate on Election", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", v.NodeId.String(), "err", err)
			return err
		}

		if staking.Is_LowRatio(can.Status) || staking.Is_DuplicateSign(can.Status) {
			//addr, _ := xutil.NodeId2Addr(v.NodeId)
			slashCans[v.NodeId] = can
			slashAddrQueue = append(slashAddrQueue, v.NodeId)

			if staking.Is_Invalid(can.Status) {
				invalidLen++
			}
		}
	}

	shuffle := func(deleteLen int, shiftQueue staking.ValidatorQueue) staking.ValidatorQueue {

		realDeleteLen := 0

		if currLen-deleteLen+len(shiftQueue) < mbn {

			log.Warn("Warn Election, finally the next round validators num less than Minimum allowed", "blockNumber",
				blockNumber, "blockHash", blockHash.Hex(), "next round num will be", currLen-deleteLen+len(shiftQueue),
				"Minimum allowed", mbn)

			// Must remain one validator TODO (Normally, this should not be the case.)
			realDeleteLen = deleteLen - 1
		} else {

			// Maybe this diffQueue length large than eight,
			// But it must less than current validator size.
			realDeleteLen = deleteLen
		}

		// Sort before removal
		if realDeleteLen != 0 {
			curr.Arr.ValidatorSort(slashCans, staking.CompareForDel)
		}

		// Increase term of validator
		nextValidators := make(staking.ValidatorQueue, len(curr.Arr))
		copy(nextValidators, curr.Arr)

		for i, v := range nextValidators {
			v.ValidatorTerm++
			nextValidators[i] = v
		}

		// Replace the validators that can be replaced
		nextValidators = nextValidators[deleteLen:]

		if len(shiftQueue) != 0 {
			nextValidators = append(nextValidators, shiftQueue...)
		}
		// Sort before storage
		nextValidators.ValidatorSort(nil, staking.CompareForStore)
		return nextValidators
	}

	var nextQueue staking.ValidatorQueue

	if invalidLen > 0 && invalidLen >= diffQueueLen {

		log.Warn("Call Election, the invalidLen large than or equal diffQueueLen", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "invalidLen", invalidLen, "diffQueueLen", diffQueueLen)

		nextQueue = shuffle(invalidLen, diffQueue)
	} else {

		if diffQueueLen <= int(xcom.ShiftValidatorNum()) {

			log.Info("Call Election, the invalidLen less than  diffQueueLen AND diffQueueLen less than ShiftValidatorNum",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "invalidLen", invalidLen,
				"diffQueueLen", diffQueueLen, "ShiftValidatorNum", xcom.ShiftValidatorNum())

			nextQueue = shuffle(diffQueueLen, diffQueue)
		} else {
			/**
			elect ShiftValidatorNum (default is 8) validators by vrf
			*/
			if queue, err := sk.VrfElection(diffQueue, header.Nonce.Bytes(), header.ParentHash); nil != err {
				log.Error("Failed to VrfElection on Election",
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
				return err
			} else {

				if invalidLen >= len(queue) {

					log.Info("Call Election, the invalidLen large than or equal vrf queue", "blockNumber",
						blockNumber, "blockHash", blockHash.Hex(), "invalidLen", invalidLen, "vrfQueueLen", len(queue))

					nextQueue = shuffle(invalidLen, queue)

				} else {

					log.Info("Call Election, the next is curr shift vrf queue", "blockNumber",
						blockNumber, "blockHash", blockHash.Hex(), "vrfQueueLen", len(queue))

					nextQueue = shuffle(len(queue), queue)
				}
			}
		}

	}

	next := &staking.Validator_array{
		Start: start,
		End:   end,
		Arr:   nextQueue,
	}

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("Election pposHash, Method Before setRoundValList", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("Election, Method Before setRoundValList, arr", next)

	if err := sk.setRoundValList(blockHash, next); nil != err {
		log.Error("Failed to SetNextValidatorList on Election", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// todo test
	if len(slashAddrQueue) != 0 {
		log.Debug("Election Slashing nodeId", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex())
		xcom.PrintObject("Election Slashing nodeId", slashAddrQueue)
	}

	// update candidate status
	// Must sort
	for _, nodeId := range slashAddrQueue {
		can := slashCans[nodeId]
		if !staking.Is_Invalid(can.Status) && staking.Is_LowRatio(can.Status) {

			// clean the low package ratio status
			can.Status &^= staking.LowRatio

			// TODO test
			log.Debug("Election slashed nodeId, clean lowratio", "nodeId", nodeId.String())
			xcom.PrintObject("Election slashed nodeId, clean lowratio", can)

			addr, _ := xutil.NodeId2Addr(nodeId)

			// TODO test
			pposHash := sk.db.GetLastKVHash(blockHash)
			log.Debug("Election pposHash, Method Before SetCandidateStore", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

			// todo test
			xcom.PrintObject("Election, Method Before SetCandidateStore, can", can)

			if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
				log.Error("Failed to Store Candidate on Election", "blockNumber", blockNumber,
					"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
				return err
			}
		}
	}

	// todo test
	if len(slashAddrQueue) == 0 {
		log.Debug("Election slashed nodeId Done, the slashAddrQueue len is zero ...")
	}

	log.Info("Call Election end", "next round validators length", len(nextQueue))

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("Election pposHash, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("Call Election Next validators", next)
	return nil
}

func (sk *StakingPlugin) SlashCandidates(state xcom.StateDB, blockHash common.Hash, blockNumber uint64,
	nodeId discover.NodeID, amount *big.Int, needDelete bool, slashType int, caller common.Address) error {

	log.Info("Call SlashCandidates start", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"slashType", slashType, "nodeId", nodeId.String(), "amount", amount, "needDelete", needDelete,
		"reporter", caller.Hex())

	contract_balance := state.GetBalance(vm.StakingContractAddr)
	if contract_balance.Cmp(common.Big0) <= 0 {
		log.Error("Failed to SlashCandidates: the balance is invalid of stakingContracr Account", "contract_balance", contract_balance)
		panic("the balance is invalid of stakingContracr Account")
	}

	addr, _ := xutil.NodeId2Addr(nodeId)
	can, err := sk.db.GetCandidateStore(blockHash, addr)
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

	aboutRelease := new(big.Int).Add(can.Released, can.ReleasedHes)
	aboutRestrictingPlan := new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
	total := new(big.Int).Add(aboutRelease, aboutRestrictingPlan)

	if total.Cmp(amount) < 0 {
		log.Error("Failed to SlashCandidates: the candidate total staking amount is not enough",
			"candidate total amount", total, "slashing amount", amount, "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return common.BizErrorf("Failed to SlashCandidates: the candidate total staking amount is not enough"+
			", candidate total amount:%s, slashing amount: %s", total, amount)
	}

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("SlashCandidates pposHash, Method Before DelCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("SlashCandidates, Method Before DelCanPowerStore, can", can)

	// clean the candidate power, first
	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Call SlashCandidates: Delete candidate old power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return err
	}

	remain := amount

	slashFunc := func(title string, remain, balance *big.Int, isNotify bool) (*big.Int, *big.Int, error) {

		remainTmp := common.Big0
		balanceTmp := common.Big0

		if remain.Cmp(balance) >= 0 {

			state.SubBalance(vm.StakingContractAddr, balance)
			if staking.Is_DuplicateSign(uint32(slashType)) {
				state.AddBalance(caller, balance)
			} else {
				state.AddBalance(vm.RewardManagerPoolAddr, balance)
			}

			if isNotify {
				err := rt.SlashingNotify(can.StakingAddress, balance, state)
				if nil != err {
					log.Error("Failed to SlashCandidates: call restrictingPlugin SlashingNotify() failed", "amount", balance,
						"slash:", title, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
					return remainTmp, balanceTmp, err
				}
			}

			remainTmp = new(big.Int).Sub(remain, balance)
			balanceTmp = common.Big0

		} else {
			state.SubBalance(vm.StakingContractAddr, remain)
			if staking.Is_DuplicateSign(uint32(slashType)) {
				state.AddBalance(caller, balance)
			} else {
				state.AddBalance(vm.RewardManagerPoolAddr, balance)
			}

			if isNotify {
				err := rt.SlashingNotify(can.StakingAddress, remain, state)
				if nil != err {
					log.Error("Failed to SlashCandidates: call restrictingPlugin SlashingNotify() failed", "amount", remain,
						"slash:", title, "blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
					return remainTmp, balanceTmp, err
				}
			}

			remainTmp = common.Big0
			balanceTmp = new(big.Int).Sub(balance, remain)
		}

		return remainTmp, balanceTmp, nil
	}

	if can.ReleasedHes.Cmp(common.Big0) > 0 {

		val, rval, err := slashFunc("ReleasedHes", remain, can.ReleasedHes, false)
		if nil != err {
			return err
		}
		remain, can.ReleasedHes = val, rval

	}

	if remain.Cmp(common.Big0) > 0 && can.RestrictingPlanHes.Cmp(common.Big0) > 0 {
		val, rval, err := slashFunc("RestrictingPlanHes", remain, can.RestrictingPlanHes, true)
		if nil != err {
			return err
		}
		remain, can.RestrictingPlanHes = val, rval
	}

	if remain.Cmp(common.Big0) > 0 && can.Released.Cmp(common.Big0) > 0 {
		val, rval, err := slashFunc("Released", remain, can.Released, false)
		if nil != err {
			return err
		}
		remain, can.Released = val, rval
	}

	if remain.Cmp(common.Big0) > 0 && can.RestrictingPlan.Cmp(common.Big0) > 0 {
		val, rval, err := slashFunc("RestrictingPlan", remain, can.RestrictingPlan, true)
		if nil != err {
			return err
		}
		remain, can.RestrictingPlan = val, rval
	}

	if remain.Cmp(common.Big0) != 0 {
		log.Error("Failed to SlashCandidates: the ramain is not zero", "remain", remain,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return common.BizErrorf("Failed to SlashCandidates: the ramain is not zero, remain:%s", remain)
	}

	// If the status has been modified before, this time it will not be modified.
	if staking.Is_Invalid_LowRatio_NotEnough(can.Status) || staking.Is_Invalid_DuplicateSign(can.Status) {
		log.Info("Call SlashCandidates end, the can status has been modified",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return nil
	}

	// sub Shares to effect power
	can.Shares = new(big.Int).Sub(can.Shares, amount)

	remainRelease := new(big.Int).Add(can.Released, can.ReleasedHes)
	remainRestrictingPlan := new(big.Int).Add(can.RestrictingPlan, can.RestrictingPlanHes)
	canRemain := new(big.Int).Add(remainRelease, remainRestrictingPlan)

	if slashType == staking.LowRatio {
		can.Status |= staking.LowRatio
		if !xutil.CheckStakeThreshold(canRemain) {
			can.Status |= staking.NotEnough
			needDelete = true
		}
	} else if slashType == staking.DuplicateSign {
		can.Status |= staking.DuplicateSign
		needDelete = true
	} else {
		log.Error("Failed to SlashCandidates: the slashType is wrong", "slashType", slashType,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
		return common.BizErrorf("Failed to SlashCandidates: the slashType is wrong, slashType: %d", slashType)
	}

	// the can status is valid and do not need delete
	// update the can power
	if !needDelete && !staking.Is_Invalid(can.Status) {

		// TODO test
		pposHash := sk.db.GetLastKVHash(blockHash)
		log.Debug("SlashCandidates pposHash, Method Before SetCanPowerStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("SlashCandidates, Method Before SetCanPowerStore, can", can)

		// update the candidate power, If do not need to delete power (the candidate status still be valid)
		sk.db.SetCanPowerStore(blockHash, addr, can)
	} else {
		//because of deleted candidate info ,clean Shares
		can.Shares = common.Big0
		can.Status |= staking.Invalided

		validators, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
		if nil != err {
			log.Error("Failed to SlashCandidates: Query Verifier List is failed", "slashType", slashType,
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		// remove the val from epoch validators,
		// becouse the candidate status is invalid after slashed
		orginLen := len(validators.Arr)
		for i := 0; i < len(validators.Arr); i++ {

			val := validators.Arr[i]

			if val.NodeId == nodeId {

				log.Info("Delete the validator when slash candidate on SlashCandidates", "slashType", slashType,
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())

				validators.Arr = append(validators.Arr[:i], validators.Arr[i+1:]...)
				i--
				break
			}
		}
		dirtyLen := len(validators.Arr)

		if dirtyLen != orginLen {

			// TODO test
			pposHash := sk.db.GetLastKVHash(blockHash)
			log.Debug("SlashCandidates pposHash, Method Before setVerifierList", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

			// todo test
			xcom.PrintObject("SlashCandidates, Method Before setVerifierList, validators", validators)

			if err := sk.setVerifierList(blockHash, validators); nil != err {
				log.Error("Failed to SlashCandidates: Store Verifier List is failed", "slashType", slashType,
					"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "err", err)
				return err
			}
		}
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("SlashCandidates pposHash, Method Before SetCandidateStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("SlashCandidates, Method Before SetCandidateStore, can", can)

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Failed to SlashCandidates: Store candidate is failed", "slashType", slashType,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("SlashCandidates pposHash, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	log.Info("Call SlashCandidates end ...")
	return nil
}

func (sk *StakingPlugin) ProposalPassedNotify(blockHash common.Hash, blockNumber uint64, nodeIds []discover.NodeID,
	programVersion uint32) error {

	log.Debug("Call ProposalPassedNotify to promote candidate programVersion", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "version", programVersion, "nodeIdQueueSize", len(nodeIds))

	// delete low version validator of epoch
	epochValidators, err := sk.getVerifierList(blockHash, blockNumber, QueryStartNotIrr)
	if nil != err {
		log.Error("Failed to ProposalPassedNotify: No found the VerifierLIst", "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "err", err)
		return err
	}

	version := xutil.CalcVersion(programVersion)

	epochNodeIds := make(map[discover.NodeID]struct{})

	for _, val := range epochValidators.Arr {
		epochNodeIds[val.NodeId] = struct{}{}
	}

	for _, nodeId := range nodeIds {

		addr, _ := xutil.NodeId2Addr(nodeId)
		can, err := sk.db.GetCandidateStore(blockHash, addr)
		if nil != err {
			log.Error("Call ProposalPassedNotify: Query Candidate is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		if nil == can {

			log.Error("Call ProposalPassedNotify: Promote candidate programVersion failed, the can is empty",
				"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String())
			continue
		}

		// TODO test
		pposHash := sk.db.GetLastKVHash(blockHash)
		log.Debug("ProposalPassedNotify pposHash, Method Before DelCanPowerStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("ProposalPassedNotify, Method Before DelCanPowerStore, can", can)

		if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
			log.Error("Call ProposalPassedNotify: Delete Candidate old power is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		can.ProgramVersion = version

		// TODO test
		pposHash = sk.db.GetLastKVHash(blockHash)
		log.Debug("ProposalPassedNotify pposHash, Method Before SetCanPowerStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("ProposalPassedNotify, Method Before SetCanPowerStore, can", can)

		if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
			log.Error("Call ProposalPassedNotify: Store Candidate new power is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}

		// TODO test
		pposHash = sk.db.GetLastKVHash(blockHash)
		log.Debug("ProposalPassedNotify pposHash, Method Before SetCandidateStore", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

		// todo test
		xcom.PrintObject("ProposalPassedNotify, Method Before SetCandidateStore, can", can)

		if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
			log.Error("Call ProposalPassedNotify: Store Candidate info is failed", "blockNumber", blockNumber,
				"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
			return err
		}
		delete(epochNodeIds, nodeId)
	}

	arr := make(staking.ValidatorQueue, len(epochValidators.Arr))
	copy(arr, epochValidators.Arr)

	for i := 0; i < len(arr); i++ {
		val := arr[i]
		if _, ok := epochNodeIds[val.NodeId]; ok {
			arr = append(arr[:i], arr[i+1:]...)
			i--
		}
	}
	epochValidators.Arr = arr

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("ProposalPassedNotify pposHash, Method Before setVerifierList", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("ProposalPassedNotify, Method Before setVerifierList, epochValidators", epochValidators)

	// update epoch validators
	if err := sk.setVerifierList(blockHash, epochValidators); nil != err {
		log.Error("Call ProposalPassedNotify: Store epoch validators after update validators is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("ProposalPassedNotify pposHash, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	return nil
}

func (sk *StakingPlugin) DeclarePromoteNotify(blockHash common.Hash, blockNumber uint64, nodeId discover.NodeID,
	programVersion uint32) error {

	log.Debug("Call DeclarePromoteNotify to promote candidate programVersion", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "version", programVersion, "nodeId", nodeId.String())

	addr, _ := xutil.NodeId2Addr(nodeId)
	can, err := sk.db.GetCandidateStore(blockHash, addr)
	if nil != err {
		log.Error("Call DeclarePromoteNotify: Query Candidate is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	if nil == can {

		log.Error("Call DeclarePromoteNotify: Promote candidate programVersion failed, the can is empty",
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(),
			"version", programVersion)
		return nil
	}

	// TODO test
	pposHash := sk.db.GetLastKVHash(blockHash)
	log.Debug("DeclarePromoteNotify pposHash, Method Before DelCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("DeclarePromoteNotify, Method Before DelCanPowerStore, can", can)

	if err := sk.db.DelCanPowerStore(blockHash, can); nil != err {
		log.Error("Call DeclarePromoteNotify: Delete Candidate old power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	can.ProgramVersion = xutil.CalcVersion(programVersion)

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("DeclarePromoteNotify pposHash, Method Before SetCanPowerStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("DeclarePromoteNotify, Method Before SetCanPowerStore, can", can)

	if err := sk.db.SetCanPowerStore(blockHash, addr, can); nil != err {
		log.Error("Call DeclarePromoteNotify: Store Candidate new power is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("DeclarePromoteNotify pposHash, Method Before SetCandidateStore", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	// todo test
	xcom.PrintObject("DeclarePromoteNotify, Method Before SetCandidateStore, can", can)

	if err := sk.db.SetCandidateStore(blockHash, addr, can); nil != err {
		log.Error("Call DeclarePromoteNotify: Store Candidate info is failed", "blockNumber", blockNumber,
			"blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return err
	}

	// TODO test
	pposHash = sk.db.GetLastKVHash(blockHash)
	log.Debug("DeclarePromoteNotify pposHash, Method End", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "pposHash", hex.EncodeToString(pposHash))

	return nil
}

func (sk *StakingPlugin) GetLastNumber(blockNumber uint64) uint64 {

	val_arr, err := sk.getCurrValList(common.ZeroHash, blockNumber, QueryStartIrr)
	if nil != err && err != snapshotdb.ErrNotFound {
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
		return build_CBFT_Validators(val_arr.Arr), nil
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
				"index length", len(indexs), "the number", i+1, "Start", indexInfo.Start, "End", indexInfo.End, "err", err)
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

func build_CBFT_Validators(arr staking.ValidatorQueue) *cbfttypes.Validators {

	valMap := make(cbfttypes.ValidateNodeMap, len(arr))

	for i, v := range arr {

		pubKey, _ := v.NodeId.Pubkey()

		vn := &cbfttypes.ValidateNode{
			Index:   i,
			Address: v.NodeAddress,
			PubKey:  pubKey,
		}

		valMap[v.NodeId] = vn
	}

	res := &cbfttypes.Validators{
		Nodes: valMap,
	}
	return res
}

func lazyCalcStakeAmount(epoch uint64, can *staking.Candidate) {

	changeAmountEpoch := can.StakingEpoch

	sub := epoch - uint64(changeAmountEpoch)

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio() {
		return
	}

	// todo test
	xcom.PrintObject("lazyCalcStakeAmount before, epoch:"+fmt.Sprint(epoch)+", can", can)

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

	// If it is during the same hesitation period, short circuit
	if sub < xcom.HesitateRatio() {
		return
	}

	// todo test
	xcom.PrintObject("lazyCalcDelegateAmount before, epoch:"+fmt.Sprint(epoch)+", del", del)

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
func (sk *StakingPlugin) VrfElection(validatorList staking.ValidatorQueue, nonce []byte, parentHash common.Hash) (staking.ValidatorQueue, error) {
	preNonces, err := xcom.GetVrfHandlerInstance().Load(parentHash)
	if nil != err {
		return nil, err
	}
	if len(preNonces) < len(validatorList) {
		log.Error("vrfElection failed", "validatorListSize", len(validatorList), "nonceSize", len(nonce), "preNoncesSize", len(preNonces), "parentHash", hex.EncodeToString(parentHash.Bytes()))
		return nil, ParamsErr
	}
	if len(preNonces) > len(validatorList) {
		preNonces = preNonces[len(preNonces)-len(validatorList):]
	}
	return sk.ProbabilityElection(validatorList, vrf.ProofToHash(nonce), preNonces)
}

func (sk *StakingPlugin) ProbabilityElection(validatorList staking.ValidatorQueue, currentNonce []byte, preNonces [][]byte) (staking.ValidatorQueue, error) {
	if len(currentNonce) == 0 || len(preNonces) == 0 || len(validatorList) != len(preNonces) {
		log.Error("probabilityElection failed", "validatorListSize", len(validatorList), "currentNonceSize", len(currentNonce), "preNoncesSize", len(preNonces), "EpochValidatorNum", xcom.EpochValidatorNum)
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
	p := (sumWeightsFloat / float64(len(validatorList))) * float64(xcom.ShiftValidatorNum()) / sumWeightsFloat
	log.Info("probabilityElection Basic parameter", "validatorListSize", len(validatorList), "p", p, "sumWeights", sumWeightsFloat, "shiftValidatorNum", xcom.ShiftValidatorNum, "epochValidatorNum", xcom.EpochValidatorNum)
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
		log.Debug("calculated probability", "nodeId", hex.EncodeToString(sv.v.NodeId.Bytes()), "addr", hex.EncodeToString(sv.v.NodeAddress.Bytes()), "index", index, "currentNonce", hex.EncodeToString(currentNonce), "preNonce", hex.EncodeToString(preNonces[index]), "target", target, "targetP", targetP, "weight", sv.weights, "x", x, "version", sv.version, "blockNumber", sv.blockNumber, "txIndex", sv.txIndex)
	}
	sort.Sort(svList)
	resultValidatorList := make(staking.ValidatorQueue, 0)
	for index, sv := range svList {
		if index == int(xcom.ShiftValidatorNum()) {
			break
		}
		resultValidatorList = append(resultValidatorList, sv.v)
		log.Debug("sort validator", "addr", hex.EncodeToString(sv.v.NodeAddress.Bytes()), "index", index, "weight", sv.weights, "x", sv.x, "version", sv.version, "blockNumber", sv.blockNumber, "txIndex", sv.txIndex)
	}
	return resultValidatorList, nil
}

/**
Internal expansion function
*/

// previous round validators
func (sk *StakingPlugin) getPreValList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.Validator_array, error) {

	var targetIndex *staking.ValArrIndex

	if !isCommit {
		indexs, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		for i, index := range indexs {
			if index.Start <= blockNumber && index.End >= blockNumber && 0 < i {
				targetIndex = indexs[i-1]
				break
			}
		}
	} else {
		indexs, err := sk.db.GetRoundValIndexByIrr()
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

		for i, index := range indexs {
			if index.Start <= blockNumber && index.End >= blockNumber && 0 < i {
				targetIndex = indexs[i-1]
				break
			}
		}
	}

	if nil == targetIndex {
		log.Error("No Found previous validators index", "isCommit", isCommit,
			"current blockNumber", blockNumber, "current blockHash", blockHash.Hex())
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
		log.Error("No Found previous validators", "isCommit", isCommit, "start", targetIndex.Start,
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

	if !isCommit {
		indexs, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

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
		log.Error("No Found current validators", "isCommit", isCommit, "start", targetIndex.Start,
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

	if !isCommit {
		indexs, err := sk.db.GetRoundValIndexByBlockHash(blockHash)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

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
		log.Error("No Found next validators", "isCommit", isCommit, "start", targetIndex.Start,
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
				"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex())
			return err
		}
	}

	// Store new index Arr
	if err := sk.db.SetRoundValIndex(blockHash, queue); nil != err {
		log.Error("Failed to setRoundValList: store round validators new indexArr is failed", "blockHash", blockHash.Hex())
		return err
	}

	// Store new round validator Item
	if err := sk.db.SetRoundValList(blockHash, index.Start, index.End, val_Arr.Arr); nil != err {
		log.Error("Failed to setRoundValList: store new round validators is failed", "blockHash", blockHash.Hex())
		return err
	}

	return nil
}

func (sk *StakingPlugin) getVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (*staking.Validator_array, error) {

	var targetIndex *staking.ValArrIndex

	if !isCommit {
		indexs, err := sk.db.GetEpochValIndexByBlockHash(blockHash)
		if nil != err && err != snapshotdb.ErrNotFound {
			return nil, err
		}

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
		log.Error("No Found epoch validators", "isCommit", isCommit, "start", targetIndex.Start,
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
				"shabby start", shabby.Start, "shabby end", shabby.End, "blockHash", blockHash.Hex())
			return err
		}
	}

	// Store new index Arr
	if err := sk.db.SetEpochValIndex(blockHash, queue); nil != err {
		log.Error("Failed to setVerifierList: store epoch validators new indexArr is failed", "blockHash", blockHash.Hex())
		return err
	}

	// Store new epoch validator Item
	if err := sk.db.SetEpochValList(blockHash, index.Start, index.End, val_Arr.Arr); nil != err {
		log.Error("Failed to setVerifierList: store new epoch validators is failed", "blockHash", blockHash.Hex())
		return err
	}

	return nil
}

func (sk *StakingPlugin) HasStake(blockHash common.Hash, addr common.Address) (bool, error) {
	return sk.db.HasAccountStakeRc(blockHash, addr)
}
