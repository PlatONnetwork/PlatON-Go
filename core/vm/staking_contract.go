package vm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
	"reflect"
)






var (
	QueryCanErr = errors.New("query candidate info is err")
	CanHasExistErr = errors.New("this candidate is exist")
)

const (
	CreateCandidateEvent       	= "CreateCandidateEvent"
	EditorCandidateEvent 		= "EditorCandidateEvent"
	CandidateWithdrawEvent      = "CandidateWithdrawEvent"
	SetCandidateExtraEvent      = "SetCandidateExtraEvent"

)


type stakingContract struct {
	plugin 		*plugin.StakingPlugin
	Contract 	*Contract
	Evm      	*EVM
}


func (stkc *stakingContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (stkc *stakingContract) Run(input []byte) ([]byte, error) {
	return stkc.execute(input)
}

func (stkc *stakingContract) FnSigns () map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		1000: stkc.CreateCandidate,
		1001: stkc.EditorCandidate,
		1002: stkc.WithdrewCandidate,
		1003: stkc.Delegate,
		1004: stkc.WithdrewDelegate,

		// Get
		2000: stkc.GetVerifierList,
		2001: stkc.GetValidatorList,
		2002: stkc.GetCandidateList,
		2003: stkc.GetDelegateListByAddr,
		2004: stkc.GetDelegateInfo,
		2005: stkc.GetCandidateInfo,
	}
}


func (stkc *stakingContract) execute (input []byte) (ret []byte, err error) {

	// verify the tx data by contracts method
	fn, params, err := plugin.Verify_tx_data(input, stkc.FnSigns())
	if nil != err {
		return nil, err
	}

	// execute contracts method
	result := reflect.ValueOf(fn).Call(params)
	if _, err := result[1].Interface().(error); !err {
		return result[0].Bytes(), nil
	}
	return nil, nil
}



func (stkc *stakingContract) CreateCandidate (typ uint16, benifitAddress common.Address, nodeId discover.NodeID, externalId, nodeName, website, details string) ([]byte, error) {
	//deposit := stkc.Contract.value
	txHash := stkc.Evm.StateDB.TxHash()
	//txindex := stkc.Evm.StateDB.TxIdx()
	blockNumber := stkc.Evm.BlockNumber
	currentHash := stkc.Evm.CurrentBlockHash

	state := stkc.Evm.StateDB

	log.Info("Call CreateCandidate of stakingContract", "txHash", txHash.Hex(), "blockNumber", blockNumber.Uint64(), "nodeId", nodeId.String())

	can, err := stkc.plugin.GetCandidateInfo(state, currentHash, nodeId)
	if nil != err {
		er := fmt.Errorf("Failed to CreateCandidate the reason is %s : %s", QueryCanErr.Error(), err.Error())
		log.Error(er.Error(), "txHash", txHash.Hex(), "blockNumber", blockNumber.Uint64())
		return nil, er
	}

	if nil != can {
		res := xcom.Result{false, "", CanHasExistErr.Error()}
		event, _ := json.Marshal(res)
		stkc.badLog(state, blockNumber.Uint64(), txHash.Hex(), CreateCandidateEvent, string(event))
		return nil, nil
	}





	return nil, nil
}


func (stkc *stakingContract) EditorCandidate (typ, amountType uint16, benifitAddress common.Address, nodeId discover.NodeID, externalId, nodeName, website, details string, amount *big.Int) {
	//deposit := stkc.Contract.value
	//txHash := stkc.Evm.StateDB.TxHash()
	//txindex := stkc.Evm.StateDB.TxIdx()
	//blockNumber := stkc.Evm.Context.BlockNumber



	//state := stkc.Evm.StateDB
	//
	//
	//stkc.plugin.GetVal(state)

}


func (stkc *stakingContract) WithdrewCandidate (nodeId discover.NodeID) {
	//txHash := stkc.Evm.StateDB.TxHash()
	//txindex := stkc.Evm.StateDB.TxIdx()
	//blockNumber := stkc.Evm.Context.BlockNumber

}


func (stkc *stakingContract) Delegate (typ uint16, stakingBlockNum uint64, nodeId discover.NodeID, amout *big.Int) {

}

func (stkc *stakingContract) WithdrewDelegate (stakingBlockNum uint64, nodeId discover.NodeID, amout *big.Int)  {

}


func (stkc *stakingContract) GetVerifierList () {

}

func (stkc *stakingContract) GetValidatorList () {

}

func (stkc *stakingContract) GetCandidateList () {

}

// todo Maybe will implement
func (stkc *stakingContract) GetDelegateListByAddr (addr common.Address) {

}

func (stkc *stakingContract) GetDelegateInfo (stakingBlockNum uint64, addr common.Address, nodeId discover.NodeID) {

}

func (stkc *stakingContract) GetCandidateInfo (nodeId discover.NodeID) {

}


func (stkc *stakingContract) goodLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData string){
	xcom.AddLog(state, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Info("Successed to CreateCandidate", "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}

func (stkc *stakingContract) badLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData string){
	xcom.AddLog(state, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Error("Failed to CreateCandidate", "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}



