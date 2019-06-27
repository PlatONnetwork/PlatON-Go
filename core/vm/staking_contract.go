package vm

import (
	"encoding/json"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
	"reflect"
)

const (
	CanAlreadyExistsErrStr = "this candidate is already exists"
	CanNotExistErrStr      = "this candidate is not exist"
	CreateCanErrStr        = "create candidate failed"

	CanStatusInvalidErrStr = "this candidate status was invalided"

	DelegateNotExistErrStr = "this is delegate is not exist"
	DelegateErrStr         = "Delegate failed"

	StakingAddrNoSomeErrStr = "address must be the same as initiated staking"
	EditCanErrStr           = "edit candidate failed"

	IncreaseStakingErrStr   = "increaseStaking failed"

	StakeVonTooLowStr = "Staking deposit too low"
	DelegateVonTooLowStr = "Delegate deposit too low"

	WithdrewCanErrStr = "withdrew candidate failed"

	QueryCanErrStr = "query candidate info err"


	AmountIllegalErrStr = "this amount is illege"



	GetVerifierListErrStr = "getting verifierList is failed"


)

const (
	CreateStakingEvent   = "1000"
	EditorCandidateEvent   = "1001"
	IncreaseStakingEvent   = "1002"
	WithdrewCandidateEvent = "1003"
	DelegateEvent          = "1004"
	WithdrewDelegateEvent  = "1005"
)

type stakingContract struct {
	plugin   *plugin.StakingPlugin
	Contract *Contract
	Evm      *EVM
}

func (stkc *stakingContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (stkc *stakingContract) Run(input []byte) ([]byte, error) {
	return stkc.execute(input)
}

func (stkc *stakingContract) FnSigns() map[uint16]interface{} {
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
		1103: stkc.getDelegateListByAddr,
		1104: stkc.getDelegateInfo,
		1105: stkc.getCandidateInfo,
	}
}

func (stkc *stakingContract) execute(input []byte) (ret []byte, err error) {

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


func (stkc *stakingContract) createStaking(typ uint16, benifitAddress common.Address, nodeId discover.NodeID,
	externalId, nodeName, website, details string, amount *big.Int, processVersion uint32) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	txIndex := stkc.Evm.StateDB.TxIdx()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call createStaking of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String())


	if amount.Cmp(common.Big0) <= 0 {
		res := xcom.Result{false, "", AmountIllegalErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), CreateStakingEvent, string(event), "createStaking")
		return nil, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to createStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err {
		log.Error("Failed to createStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil != canOld {
		res := xcom.Result{false, "", CanAlreadyExistsErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), CreateStakingEvent, string(event), "createStaking")
		return nil, nil
	}

	if !plugin.CheckStakeThreshold(amount) {
		res := xcom.Result{false, "", StakeVonTooLowStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), CreateStakingEvent, string(event), "createStaking")
		return nil, nil
	}

	/**
	init candidate info
	*/
	canTmp := &xcom.Candidate{
		NodeId:          nodeId,
		StakingAddress:  from,
		BenifitAddress:  benifitAddress,
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  txIndex,
		Shares:          amount,

		Description: xcom.Description{
			NodeName:   nodeName,
			ExternalId: externalId,
			Website:    website,
			Details:    details,
		},
	}

	flag, err := stkc.plugin.CreateCandidate(state, blockHash, blockNumber, amount, processVersion, typ, canAddr, canTmp)
	if nil != err {
		if flag {
			res := xcom.Result{false, "", CreateCanErrStr + ":" + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), CreateStakingEvent, string(event), "createStaking")
			return nil, nil
		} else {
			log.Error("Failed to createStaking by CreateCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}

	res := xcom.Result{true, "", ""}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash.Hex(), CreateStakingEvent, string(event), "createStaking")
	return nil, nil
}

func (stkc *stakingContract) editorCandidate(benifitAddress common.Address, nodeId discover.NodeID,
	externalId, nodeName, website, details string, amount *big.Int) ([]byte, error) {

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

	canOld, err := stkc.plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err {
		log.Error("Failed to editorCandidate by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), EditorCandidateEvent, string(event), "editorCandidate")
		return nil, nil
	}

	if !xcom.Is_Valid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), EditorCandidateEvent, string(event), "editorCandidate")
		return nil, nil
	}

	if from != canOld.StakingAddress {
		res := xcom.Result{false, "", StakingAddrNoSomeErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), EditorCandidateEvent, string(event), "editorCandidate")
		return nil, nil
	}

	canOld.BenifitAddress = benifitAddress

	canOld.NodeName = nodeName
	canOld.ExternalId = externalId
	canOld.Website = website
	canOld.Details = details

	flag, err := stkc.plugin.EditorCandidate(blockHash, blockNumber, canOld)

	if nil != err {

		if flag {
			res := xcom.Result{false, "", EditCanErrStr + ":" + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), EditorCandidateEvent, string(event), "editorCandidate")
			return nil, nil
		} else {
			log.Error("Failed to editorCandidate by EditorCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}
	res := xcom.Result{true, "", ""}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash.Hex(), EditorCandidateEvent, string(event), "editorCandidate")
	return nil, nil
}


func (stkc *stakingContract) increaseStaking (nodeId discover.NodeID, typ uint16, amount *big.Int) ([]byte, error) {

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
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), IncreaseStakingEvent, string(event), "increaseStaking")
		return nil, nil
	}


	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to increaseStaking by parse nodeId", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err {
		log.Error("Failed to increaseStaking by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), IncreaseStakingEvent, string(event), "increaseStaking")
		return nil, nil
	}

	if !xcom.Is_Valid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), IncreaseStakingEvent, string(event), "increaseStaking")
		return nil, nil
	}

	if from != canOld.StakingAddress {
		res := xcom.Result{false, "", StakingAddrNoSomeErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), IncreaseStakingEvent, string(event), "increaseStaking")
		return nil, nil
	}

	flag, err := stkc.plugin.IncreaseStaking(state, blockHash, blockNumber, amount, typ, canOld)

	if nil != err {

		if flag {
			res := xcom.Result{false, "", IncreaseStakingErrStr + ":" + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), IncreaseStakingEvent, string(event), "increaseStaking")
			return nil, nil
		} else {
			log.Error("Failed to increaseStaking by EditorCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}
	res := xcom.Result{true, "", ""}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash.Hex(), IncreaseStakingEvent, string(event), "increaseStaking")
	return nil, nil
}

func (stkc *stakingContract) withdrewCandidate(nodeId discover.NodeID) ([]byte, error) {
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

	canOld, err := stkc.plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err {
		log.Error("Failed to withdrewCandidate by GetCandidateInfo", "txHash", txHash,
			"blockNumber", blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewCandidateEvent, string(event), "withdrewCandidate")
		return nil, nil
	}

	if !xcom.Is_Valid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewCandidateEvent, string(event), "withdrewCandidate")
		return nil, nil
	}

	if from != canOld.StakingAddress {
		res := xcom.Result{false, "", StakingAddrNoSomeErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewCandidateEvent, string(event), "withdrewCandidate")
		return nil, nil
	}

	flag, err := stkc.plugin.WithdrewCandidate(state, blockHash, blockNumber, canOld)
	if nil != err {

		if flag {
			res := xcom.Result{false, "", WithdrewCanErrStr + ":" + err.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewCandidateEvent,
				string(event), "withdrewCandidate")
			return nil, nil
		} else {
			log.Error("Failed to withdrewCandidate by WithdrewCandidate", "txHash", txHash,
				"blockNumber", blockNumber, "err", err)
			return nil, err
		}

	}

	res := xcom.Result{true, "", ""}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewCandidateEvent,
		string(event), "withdrewCandidate")
	return nil, nil
}

func (stkc *stakingContract) delegate(typ uint16, nodeId discover.NodeID, amount *big.Int) ([]byte, error) {
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
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), DelegateEvent, string(event), "delegate")
		return nil, nil
	}

	canAddr, err := xutil.NodeId2Addr(nodeId)
	if nil != err {
		log.Error("Failed to delegate by parse nodeId", "txHash", txHash, "blockNumber",
			blockNumber, "blockHash", blockHash.Hex(), "nodeId", nodeId.String(), "err", err)
		return nil, err
	}

	canOld, err := stkc.plugin.GetCandidateInfo(blockHash, canAddr)
	if nil != err {
		log.Error("Failed to delegate by GetCandidateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), DelegateEvent, string(event), "delegate")
		return nil, nil
	}

	if !xcom.Is_Valid(canOld.Status) {
		res := xcom.Result{false, "", CanStatusInvalidErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), DelegateEvent, string(event), "delegate")
		return nil, nil
	}

	// todo the delegate caller is candidate stake addr ?? How do that ??

	del, err := stkc.plugin.GetDelegateInfo(blockHash, from, nodeId, canOld.StakingBlockNum)
	if nil != err {
		log.Error("Failed to delegate by GetDelegateInfo", "txHash", txHash, "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == del {

		if !plugin.CheckDelegateThreshold(amount) {
			res := xcom.Result{false, "", DelegateVonTooLowStr}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), DelegateEvent, string(event), "delegate")
			return nil, nil
		}

		del = new(xcom.Delegation)
	}

	flag, er := stkc.plugin.Delegate(state, blockHash, blockNumber, from, del, canOld, typ, amount)
	if nil != er {
		if flag {
			res := xcom.Result{false, "", DelegateErrStr + ":" + er.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), DelegateEvent, string(event), "delegate")
			return nil, nil
		} else {
			log.Error("Failed to delegate by Delegate", "txHash", txHash, "blockNumber", blockNumber, "err", er)
			return nil, er
		}
	}

	res := xcom.Result{true, "", ""}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash.Hex(), DelegateEvent, string(event), "delegate")
	return nil, nil
}

func (stkc *stakingContract) withdrewDelegate(stakingBlockNum uint64, nodeId discover.NodeID, amount *big.Int) ([]byte, error) {

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
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewDelegateEvent, string(event), "withdrewDelegate")
		return nil, nil
	}


	del, err := stkc.plugin.GetDelegateInfo(blockHash, from, nodeId, stakingBlockNum)
	if nil != err {
		log.Error("Failed to withdrewDelegate by GetDelegateInfo",
			"txHash", txHash.Hex(), "blockNumber", blockNumber, "err", err)
		return nil, err
	}

	if nil == del {
		res := xcom.Result{false, "", DelegateNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewDelegateEvent, string(event), "withdrewDelegate")
		return nil, nil
	}

	flag, er := stkc.plugin.WithdrewDelegate(state, blockHash, blockNumber, amount, from, nodeId, stakingBlockNum, del)
	if nil != er {
		if flag {
			res := xcom.Result{false, "", WithdrewCanErrStr + ":" + er.Error()}
			event, _ := json.Marshal(res)
			stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewDelegateEvent, string(event), "withdrewDelegate")
			return nil, nil
		} else {
			log.Error("Failed to withdrewDelegate by WithdrewDelegate", "txHash", txHash, "blockNumber", blockNumber, "err", er)
			return nil, er
		}
	}


	res := xcom.Result{true, "", ""}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewDelegateEvent, string(event), "withdrewDelegate")
	return nil, nil
}

func (stkc *stakingContract) getVerifierList() ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	blockHash := stkc.Evm.BlockHash


	log.Info("Call getVerifierList of stakingContract", "txHash", txHash.Hex(),
		"blockNumber", blockNumber.Uint64())
	arr, flag, err := stkc.plugin.GetVerifierList(blockHash, blockNumber.Uint64(), plugin.QueryStartIrr)

	if nil != err {
		if flag {
			res := xcom.Result{false, "", GetVerifierListErrStr + ":" + err.Error()}
			data, _ := rlp.EncodeToBytes(res)
			return data, nil

		} else {
			log.Error("Failed to getVerifierList", "txHash", txHash, "blockNumber", blockNumber, "err", err)
			return nil, err
		}
	}
	arrByte, _ := json.Marshal(arr)
	res := xcom.Result{true, string(arrByte), ""}
	data, _ := rlp.EncodeToBytes(res)
	return data, nil
}

func (stkc *stakingContract) getValidatorList() ([]byte, error) {

	//txHash := stkc.Evm.StateDB.TxHash()
	//blockNumber := stkc.Evm.BlockNumber
	//blockHash := stkc.Evm.BlockHash
	//
	//validatorExArr, flag, err := stkc.plugin.GetValidatorList(blockHash, blockNumber.Uint64(), plugin.CurrentRound, plugin.QueryStartIrr)

	return nil, nil
}

func (stkc *stakingContract) getCandidateList() {



}

// todo Maybe will implement
func (stkc *stakingContract) getDelegateListByAddr(addr common.Address) {

}

func (stkc *stakingContract) getDelegateInfo(stakingBlockNum uint64, addr common.Address, nodeId discover.NodeID) {

}

func (stkc *stakingContract) getCandidateInfo(nodeId discover.NodeID) {

}

func (stkc *stakingContract) goodLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Info("flaged to "+callFn+" of stakingContract", "txHash", txHash,
		"blockNumber", blockNumber, "json: ", eventData)
}

func (stkc *stakingContract) badLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Error("Failed to "+callFn+" of stakingContract", "txHash", txHash,
		"blockNumber", blockNumber, "json: ", eventData)
}
