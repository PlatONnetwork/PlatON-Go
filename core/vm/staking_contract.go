package vm

import (
	"encoding/json"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
	"reflect"
)

const (
	QueryCanErrStr         = "query candidate info err"
	CanAlreadyExistsErrStr = "this candidate is already exists"
	CanNotExistErrStr      = "this candidate is not exist"
	CreateCanErrStr        = "create candidate failed"

	CanStatusInvalidErrStr = "this candidate status was invalided"
	StakingAddrNoSomeErrStr = "address must be the same as initiated staking"
	EditCanErrStr 			= "edit candidate failed"


	StakeVonToLowStr = "Staking deposit too low"

)

const (
	CreateCandidateEvent   = "1000"
	EditorCandidateEvent   = "1001"
	WithdrewCandidateEvent = "1002"
	DelegateEvent          = "1003"
	WithdrewDelegateEvent  = "1004"
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
		1000: stkc.createCandidate,
		1001: stkc.editorCandidate,
		1002: stkc.withdrewCandidate,
		1003: stkc.delegate,
		1004: stkc.withdrewDelegate,

		// Get
		2000: stkc.getVerifierList,
		2001: stkc.getValidatorList,
		2002: stkc.getCandidateList,
		2003: stkc.getDelegateListByAddr,
		2004: stkc.getDelegateInfo,
		2005: stkc.getCandidateInfo,
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
	if _, ok := result[1].Interface().(error); !ok {
		return result[0].Bytes(), nil
	}
	return nil, result[1].Interface().(error)
}

func (stkc *stakingContract) createCandidate(typ uint16, benifitAddress common.Address, nodeId discover.NodeID, externalId, nodeName, website, details string, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	txIndex := stkc.Evm.StateDB.TxIdx()
	blockNumber := stkc.Evm.BlockNumber
	currentHash := stkc.Evm.CurrentBlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call createCandidate of stakingContract", "txHash", txHash.Hex(), "blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String())

	canOld, err := stkc.plugin.GetCandidateInfo(currentHash, nodeId)
	if nil != err {
		res := xcom.Result{false, "", QueryCanErrStr + ":" + err.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), CreateCandidateEvent, string(event), "createCandidate")
		return nil, nil
	}

	if nil != canOld {
		res := xcom.Result{false, "", CanAlreadyExistsErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), CreateCandidateEvent, string(event), "createCandidate")
		return nil, nil
	}

	if !plugin.CheckStakeThreshold(amount) {
		res := xcom.Result{false, "", StakeVonToLowStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), CreateCandidateEvent, string(event), "createCandidate")
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

	if typ == plugin.FreeOrigin {
		canTmp.ReleasedTmp = amount
	} else if typ == plugin.LockRepoOrigin {
		canTmp.LockRepoTmp = amount
	}

	err = stkc.plugin.CreateCandidate(state, currentHash, blockNumber, typ, canTmp)
	if nil != err {
		res := xcom.Result{false, "", CreateCanErrStr + ":" + err.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), CreateCandidateEvent, string(event), "createCandidate")
		return nil, nil
	}

	res := xcom.Result{true, "", ""}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash.Hex(), CreateCandidateEvent, string(event), "createCandidate")
	return nil, nil
}

func (stkc *stakingContract) editorCandidate(typ uint16, benifitAddress common.Address, nodeId discover.NodeID, externalId, nodeName, website, details string, amount *big.Int) ([]byte, error) {

	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	currentHash := stkc.Evm.CurrentBlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call editorCandidate of stakingContract", "txHash", txHash.Hex(), "blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String())

	canOld, err := stkc.plugin.GetCandidateInfo(currentHash, nodeId)
	if nil != err {
		res := xcom.Result{false, "", QueryCanErrStr + ":" + err.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), EditorCandidateEvent, string(event), "editorCandidate")
		return nil, nil
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), EditorCandidateEvent, string(event), "editorCandidate")
		return nil, nil
	}

	if !xcom.IsCan_Valid(canOld.Status) {
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


	err = stkc.plugin.EditorCandidate(state, currentHash, blockNumber, canOld, typ, amount)

	if nil != err {
		res := xcom.Result{false, "", EditCanErrStr + ":" + err.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), EditorCandidateEvent, string(event), "editorCandidate")
		return nil, nil
	}
	res := xcom.Result{true, "", ""}
	event, _ := json.Marshal(res)
	stkc.goodLog(state, blockNumber.Uint64(), txHash.Hex(), EditorCandidateEvent, string(event), "editorCandidate")
	return nil, nil
}


func (stkc *stakingContract) withdrewCandidate(nodeId discover.NodeID)  ([]byte, error) {
	txHash := stkc.Evm.StateDB.TxHash()
	blockNumber := stkc.Evm.BlockNumber
	currentHash := stkc.Evm.CurrentBlockHash

	from := stkc.Contract.CallerAddress

	state := stkc.Evm.StateDB

	log.Info("Call withdrewCandidate of stakingContract", "txHash", txHash.Hex(), "blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String())

	canOld, err := stkc.plugin.GetCandidateInfo(currentHash, nodeId)
	if nil != err {
		res := xcom.Result{false, "", QueryCanErrStr + ":" + err.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewCandidateEvent, string(event), "withdrewCandidate")
		return nil, nil
	}

	if nil == canOld {
		res := xcom.Result{false, "", CanNotExistErrStr}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), WithdrewCandidateEvent, string(event), "withdrewCandidate")
		return nil, nil
	}

	if !xcom.IsCan_Valid(canOld.Status) {
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





}

func (stkc *stakingContract) delegate(typ uint16, stakingBlockNum uint64, nodeId discover.NodeID, amout *big.Int) {

}

func (stkc *stakingContract) withdrewDelegate(stakingBlockNum uint64, nodeId discover.NodeID, amout *big.Int) {

}

func (stkc *stakingContract) getVerifierList() {

}

func (stkc *stakingContract) getValidatorList() {

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
	log.Info("Successed to " + callFn, "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}

func (stkc *stakingContract) badLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Error("Failed to " + callFn, "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}
