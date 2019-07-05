package vm

import (
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
	"reflect"
)

const (
	AmountIllegalErrStr      = "This amount is illege"
	CanAlreadyExistsErrStr   = "This candidate is already exists"
	CanNotExistErrStr        = "This candidate is not exist"
	CreateCanErrStr          = "Create candidate failed"
	CanStatusInvalidErrStr   = "This candidate status was invalided"
	DelegateNotExistErrStr   = "This is delegate is not exist"
	DelegateErrStr           = "Delegate failed"
	DelegateVonTooLowStr     = "Delegate deposit too low"
	EditCanErrStr            = "Edit candidate failed"
	GetVerifierListErrStr    = "Getting verifierList is failed"
	GetValidatorListErrStr   = "Getting validatorList is failed"
	GetCandidateListErrStr   = "Getting candidateList is failed"
	GetDelegateRelatedErrStr = "Getting related of delegate is failed"
	IncreaseStakingErrStr    = "IncreaseStaking failed"
	QueryCanErrStr           = "Query candidate info failed"
	QueryDelErrSTr           = "Query delegate info failed"
	StakeVonTooLowStr        = "Staking deposit too low"
	StakingAddrNoSomeErrStr  = "Address must be the same as initiated staking"
	WithdrewCanErrStr        = "Withdrew candidate failed"
)

const (
	CreateStakingEvent     = "1000"
	EditorCandidateEvent   = "1001"
	IncreaseStakingEvent   = "1002"
	WithdrewCandidateEvent = "1003"
	DelegateEvent          = "1004"
	WithdrewDelegateEvent  = "1005"
)

type StakingContract struct {
	Plugin   *plugin.StakingPlugin
	Contract *Contract
	Evm      *EVM
}

func (stkc *StakingContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (stkc *StakingContract) Run(input []byte) ([]byte, error) {
	return stkc.execute(input)
}

func (stkc *StakingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		1000: stkc.createStaking,
		1001: stkc.editorCandidate,
		1002: stkc.increaseStaking,
		1003: stkc.withdrewCandidate,
		1004: stkc.delegate,
		1005: stkc.withdrewDelegate,

		// Get
		1100: stkc.getVerifierList,
		1101: stkc.getValidatorList,
		1102: stkc.getCandidateList,
		1103: stkc.getRelatedListByDelAddr,
		1104: stkc.getDelegateInfo,
		1105: stkc.getCandidateInfo,
	}
}

func (stkc *StakingContract) execute(input []byte) (ret []byte, err error) {

	// verify the tx data by contracts method
	fn, params, err := plugin.Verify_tx_data(input, stkc.FnSigns())
	if nil != err {
		return nil, err
	}

	// execute contracts method
	result := reflect.ValueOf(fn).Call(params)
	if err, ok := result[1].Interface().(error); ok {
		return nil, err
	}
	return result[0].Bytes(), nil
}

func (stkc *StakingContract) createStaking(typ uint16, benifitAddress common.Address, nodeId discover.NodeID,
	externalId, nodeName, website, details string, amount *big.Int, processVersion uint32) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	txIndex := stkc.Evm.StateDB.TxIdx()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call createStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String())

	//  TODO  MOCK
	//return stkc.createMock(state, blockNumber.Uint64(), txHash, typ, benifitAddress, nodeId,
	//	externalId, nodeName, website, details, amount, processVersion)

	if !plugin.CheckStakeThreshold(amount) {
		res := xcom.Result{false, "", StakeVonTooLowStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
		return event, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to createStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to createStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil != canOld {
		res := xcom.Result{false, "", CanAlreadyExistsErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
		return event, nil
	}



	/**
	init candidate info
	*/
	canTmp := &staking.Candidate{
		NodeId:          nodeId,
		StakingAddress:  from,
		BenifitAddress:  benifitAddress,
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  txIndex,
		Shares:          amount,

		Description: staking.Description{
			NodeName:   nodeName,
			ExternalId: externalId,
			Website:    website,
			Details:    details,
		},
	}

	// TODO  test
	canJson, _ := json.Marshal(canTmp)
	fmt.Println("Create Candidate canJson is:", string(canJson))

	err = stkc.Plugin.CreateCandidate(state, blockHash, blockNumber, amount, processVersion, typ, canAddr, canTmp)
	if nil != err {
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", CreateCanErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
			return event, nil
		} else {
			log.Error("Failed to createStaking by CreateCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, CreateStakingEvent, string(event), "createStaking")
	return event, nil
}

func (stkc *StakingContract) editorCandidate(benifitAddress common.Address, nodeId discover.NodeID,
	externalId, nodeName, website, details string) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call editorCandidate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String())

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to editorCandidate by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to editorCandidate by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editorCandidate")
		return event, nil
	}

	if !staking.Is_Valid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editorCandidate")
		return event, nil
	}

	if from != canOld.StakingAddress {
		res := xcom.Result{false, "", StakingAddrNoSomeErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editorCandidate")
		return event, nil
	}

	canOld.BenifitAddress = benifitAddress

	canOld.NodeName = nodeName
	canOld.ExternalId = externalId
	canOld.Website = website
	canOld.Details = details

	// TODO test
	canJson, _ := json.Marshal(canOld)
	fmt.Println("Edit Candidate canJson is:", string(canJson))

	err = stkc.Plugin.EditorCandidate(blockHash, blockNumber, canOld)

	if nil != err {

		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", EditCanErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editorCandidate")
			return event, nil
		} else {
			log.Error("Failed to editorCandidate by EditorCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}
	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, EditorCandidateEvent, string(event), "editorCandidate")
	return event, nil
}

func (stkc *StakingContract) increaseStaking(nodeId discover.NodeID, typ uint16, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call increaseStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String())

	if amount.Cmp(common.Big0) <= 0 {
		res := xcom.Result{false, "", AmountIllegalErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
		return event, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to increaseStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to increaseStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
		return event, nil
	}

	if !staking.Is_Valid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
		return event, nil
	}

	if from != canOld.StakingAddress {
		res := xcom.Result{false, "", StakingAddrNoSomeErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
		return event, nil
	}

	err = stkc.Plugin.IncreaseStaking(state, blockHash, blockNumber, amount, typ, canOld)

	if nil != err {

		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", IncreaseStakingErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
			return event, nil
		} else {
			log.Error("Failed to increaseStaking by EditorCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}
	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, IncreaseStakingEvent, string(event), "increaseStaking")
	return event, nil
}

func (stkc *StakingContract) withdrewCandidate(nodeId discover.NodeID) ([]byte, error) {
	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call withdrewCandidate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String())

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to withdrewCandidate by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to withdrewCandidate by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent, string(event), "withdrewCandidate")
		return event, nil
	}

	if !staking.Is_Valid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent, string(event), "withdrewCandidate")
		return event, nil
	}

	if from != canOld.StakingAddress {
		res := xcom.Result{false, "", StakingAddrNoSomeErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent, string(event), "withdrewCandidate")
		return event, nil
	}

	err = stkc.Plugin.WithdrewCandidate(state, blockHash, blockNumber, canOld)
	if nil != err {

		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", WithdrewCanErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent,
				string(event), "withdrewCandidate")
			return event, nil
		} else {
			log.Error("Failed to withdrewCandidate by WithdrewCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}

	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, WithdrewCandidateEvent,
		string(event), "withdrewCandidate")
	return event, nil
}

func (stkc *StakingContract) delegate(typ uint16, nodeId discover.NodeID, amount *big.Int) ([]byte, error) {
	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call delegate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "delAddr", from.Hex(), "nodeId", nodeId.String())

	if amount.Cmp(common.Big0) <= 0 {
		res := xcom.Result{false, "", AmountIllegalErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
		return event, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to delegate by parse nodeId", "txHash", txHash, "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.Plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to delegate by GetCandidateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
		return event, nil
	}

	if !staking.Is_Valid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
		return event, nil
	}

	// todo the delegate caller is candidate stake addr ?? How do that ??

	del, err := stkc.Plugin.GetDelegateInfo(blockHash, from, nodeId, canOld.StakingBlockNum)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to delegate by GetDelegateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == del {

		if !plugin.CheckDelegateThreshold(amount) {
			res := xcom.Result{false, "", DelegateVonTooLowStr}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
			return event, nil
		}

		del = new(staking.Delegation)
	}

	err = stkc.Plugin.Delegate(state, blockHash, blockNumber, from, del, canOld, typ, amount)
	if nil != err {
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", DelegateErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
			return event, nil
		} else {
			log.Error("Failed to delegate by Delegate", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, DelegateEvent, string(event), "delegate")
	return event, nil
}

func (stkc *StakingContract) withdrewDelegate(stakingBlockNum uint64, nodeId discover.NodeID, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call withdrewDelegate of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "delAddr", from.Hex(), "nodeId", nodeId.String())

	if amount.Cmp(common.Big0) <= 0 {
		res := xcom.Result{false, "", AmountIllegalErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewDelegateEvent, string(event), "withdrewDelegate")
		return event, nil
	}

	del, err := stkc.Plugin.GetDelegateInfo(blockHash, from, nodeId, stakingBlockNum)
	if nil != err && err != snapshotdb.ErrNotFound {
		log.Error("Failed to withdrewDelegate by GetDelegateInfo",
			"txHash", txHash.Hex(), "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == del {
		res := xcom.Result{false, "", DelegateNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewDelegateEvent, string(event), "withdrewDelegate")
		return event, nil
	}

	err = stkc.Plugin.WithdrewDelegate(state, blockHash, blockNumber, amount, from, nodeId, stakingBlockNum, del)
	if nil != err {
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", WithdrewCanErrStr + ": " + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash, WithdrewDelegateEvent, string(event), "withdrewDelegate")
			return event, nil
		} else {
			log.Error("Failed to withdrewDelegate by WithdrewDelegate", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash, WithdrewDelegateEvent, string(event), "withdrewDelegate")
	return event, nil
}

func (stkc *StakingContract) getVerifierList() ([]byte, error) {

	//  TODO  MOCK
	return stkc.getVerifierListMock()

	arr, err := stkc.Plugin.GetVerifierList(common.ZeroHash, common.Big0.Uint64(), plugin.QueryStartIrr)

	if nil != err {
		res := xcom.Result{false, "", GetVerifierListErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	jsonByte, err := json.Marshal(arr)
	if nil != err {
		res := xcom.Result{false, "", GetVerifierListErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	res := xcom.Result{true, string(jsonByte), "ok"}
	data, _ := rlp.EncodeToBytes(res)
	return data, nil
}

func (stkc *StakingContract) getValidatorList() ([]byte, error) {

	arr, err := stkc.Plugin.GetValidatorList(common.ZeroHash, common.Big0.Uint64(), plugin.CurrentRound, plugin.QueryStartIrr)
	if nil != err {
		res := xcom.Result{false, "", GetValidatorListErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	arrByte, _ := json.Marshal(arr)
	res := xcom.Result{true, string(arrByte), "ok"}
	data, _ := rlp.EncodeToBytes(res)
	return data, nil
}

func (stkc *StakingContract) getCandidateList() ([]byte, error) {

	arr, err := stkc.Plugin.GetCandidateList(common.ZeroHash, plugin.QueryStartIrr)
	if nil != err {
		res := xcom.Result{false, "", GetCandidateListErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}

	jsonByte, err := json.Marshal(arr)
	if nil != err {
		res := xcom.Result{false, "", GetCandidateListErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	res := xcom.Result{true, string(jsonByte), "ok"}
	data, _ := rlp.EncodeToBytes(res)
	return data, nil
}

// todo Maybe will implement
func (stkc *StakingContract) getRelatedListByDelAddr(addr common.Address) ([]byte, error) {

	arr, err := stkc.Plugin.GetRelatedListByDelAddr(common.ZeroHash, addr, plugin.QueryStartIrr)
	if nil != err {
		res := xcom.Result{false, "", GetDelegateRelatedErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	jsonByte, err := json.Marshal(arr)
	if nil != err {
		res := xcom.Result{false, "", GetDelegateRelatedErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	res := xcom.Result{true, string(jsonByte), "ok"}
	data, _ := rlp.EncodeToBytes(res)
	return data, nil
}

func (stkc *StakingContract) getDelegateInfo(stakingBlockNum uint64, addr common.Address,
	nodeId discover.NodeID) ([]byte, error) {

	addr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		res := xcom.Result{false, "", QueryDelErrSTr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	del, err := stkc.Plugin.GetDelegateInfoByIrr(addr, nodeId, stakingBlockNum)
	if nil != err {
		res := xcom.Result{false, "", QueryDelErrSTr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	jsonByte, err := json.Marshal(del)
	if nil != err {
		res := xcom.Result{false, "", QueryDelErrSTr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	res := xcom.Result{true, string(jsonByte), "ok"}
	data, _ := rlp.EncodeToBytes(res)
	return data, nil
}

func (stkc *StakingContract) getCandidateInfo(nodeId discover.NodeID) ([]byte, error) {

	////  TODO  MOCK
	//return stkc.getCandidateInfoMock()

	addr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		res := xcom.Result{false, "", QueryCanErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	can, err := stkc.Plugin.GetCandidateInfoByIrr(addr)
	if nil != err {
		res := xcom.Result{false, "", QueryCanErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	jsonByte, err := json.Marshal(can)
	if nil != err {
		res := xcom.Result{false, "", QueryDelErrSTr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	res := xcom.Result{true, string(jsonByte), "ok"}
	data, _ := rlp.EncodeToBytes(res)

	return data, nil

}

func (stkc *StakingContract) goodLog(state xcom.StateDB, blockNumber uint64, txHash common.Hash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.StakingContractAddr, txHash, eventType, eventData)
	log.Info("flaged to "+callFn+" of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber, "json: ", eventData)
}

func (stkc *StakingContract) badLog(state xcom.StateDB, blockNumber uint64, txHash common.Hash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.StakingContractAddr, txHash, eventType, eventData)
	log.Debug("Failed to "+callFn+" of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber, "json: ", eventData)
}

//  TODO MOCK
func (stkc *StakingContract) createMock(state xcom.StateDB, blockNumber uint64, txHash common.Hash, typ uint16, benifitAddress common.Address, nodeId discover.NodeID,
	externalId, nodeName, website, details string, amount *big.Int, processVersion uint32) ([]byte, error) {

	fmt.Println("Call createStaking ~~~~~~~~~~~~~~")

	fmt.Println("typ:", typ)
	fmt.Println("benifitAddress:", benifitAddress.Hex())
	fmt.Println("nodeId:", nodeId.String())
	fmt.Println("externalId:", externalId)
	fmt.Println("nodeName:", nodeName)
	fmt.Println("website:", website)
	fmt.Println("details:", details)
	fmt.Println("amount:", amount)
	fmt.Println("processVersion:", processVersion)

	res := xcom.Result{true, "", "ok"}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber, txHash, CreateStakingEvent, string(event), "createStaking")
	return event, nil

}

func (stkc *StakingContract) getCandidateInfoMock() ([]byte, error) {

	return nil, nil
}

func (stkc *StakingContract) getVerifierListMock() ([]byte, error) {

	fmt.Println("Call getVerifierList ~~~~~~~~~~~~~~")

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

	queue := make(staking.ValidatorExQueue, 0)
	for i := 0; i < 4; i++ {

		valEx := &staking.ValidatorEx{
			NodeId:          discover.MustHexID(nodeIdArr[i]),
			StakingAddress:  common.HexToAddress(addrArr[i]),
			BenifitAddress:  vm.StakingContractAddr,
			StakingTxIndex:  uint32(i),
			ProcessVersion:  uint32(i * i),
			StakingBlockNum: uint64(i + 2),
			Shares:          common.Big256,
			Description: staking.Description{
				ExternalId: "xxccccdddddddd",
				NodeName:   "I Am " + fmt.Sprint(i),
				Website:    "www.baidu.com",
				Details:    "this is  baidu ~~",
			},
			ValidatorTerm:   uint32(2),
		}

		queue = append(queue, valEx)
	}

	jsonByte, err := json.Marshal(queue)
	if nil != err {
		res := xcom.Result{false, "", GetVerifierListErrStr + ": " + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	res := xcom.Result{true, string(jsonByte), "ok"}
	data, _ := rlp.EncodeToBytes(res)
	return data, nil
}
