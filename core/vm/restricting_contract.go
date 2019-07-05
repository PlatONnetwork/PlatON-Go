package vm

import (
	"encoding/json"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)


type restrictingContract struct {
	plugin      *plugin.RestrictingPlugin
	Contract 	*Contract
	Evm  		*EVM
}


func (rc *restrictingContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (rc *restrictingContract) Run(input []byte) ([]byte, error) {
	return rc.execute(input)
}

func (rc *restrictingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		4000: rc.createRestrictingPlan,

		// Get
		4100: rc.getRestrictingInfo,
	}
}

func (rc *restrictingContract) execute(input []byte) ([]byte, error) {
	// verify the tx data by contracts method
	var fn, params, err = plugin.Verify_tx_data(input, rc.FnSigns())
	if nil != err {
		return nil, err
	}

	// execute contracts function
	result := reflect.ValueOf(fn).Call(params)
	err, ok := result[1].Interface().(error)
	if !ok {
		return result[0].Bytes(), nil
	} else {
		return nil, err
	}
}

func (rc *restrictingContract) createRestrictingPlan(account common.Address, plans []restricting.RestrictingPlan) ([]byte, error) {
	sender := rc.Contract.Caller()
	txHash := rc.Evm.StateDB.TxHash()
	blockNum := rc.Evm.BlockNumber
	state := rc.Evm.StateDB

	log.Info("Call createRestrictingPlan of restrictingContract", "txHash", txHash.Hex(), "blockNumber", blockNum.Uint64())

	if err := rc.plugin.AddRestrictingRecord(sender, account, plans, state); err != nil {
		if _, ok := err.(*common.SysError); ok {
			res := xcom.Result{Status:false, Data:"", ErrMsg:"create lock repo plan:" + err.Error()}
			event, _ := json.Marshal(res)
			rc.badLog(state, blockNum.Uint64(), txHash.Hex(), "4000", string(event), "createRestrictingPlan")
			return nil, nil
		} else {
			log.Error("AddRestrictingRecord failed to createRestrictingPlan", "txHash", txHash.Hex(), "blockNumber", blockNum.Uint64(), "error", err)
			return nil, err
		}
	}

	res := xcom.Result{Status:true, Data:"", ErrMsg:""}
	event, _ := json.Marshal(res)
	rc.goodLog(state, blockNum.Uint64(), txHash.Hex(), "4000", string(event), "createRestrictingPlan")

	return nil, nil
}

func (rc *restrictingContract) getRestrictingInfo(account common.Address) ([]byte, error) {

	txHash := rc.Evm.StateDB.TxHash()
	currNumber := rc.Evm.BlockNumber
	state := rc.Evm.StateDB

	log.Info("Call getRestrictingInfo of restrictingContract", "txHash", txHash.Hex(), "blockNumber", currNumber.Uint64())

	return rc.plugin.GetRestrictingInfo(account, state)
}


func (rc *restrictingContract) goodLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData, callFn string) {
	_ = xcom.AddLog(state, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Info("Successed to " + callFn, "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}


func (rc *restrictingContract) badLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData, callFn string) {
	_ = xcom.AddLog(state, blockNumber, vm.StakingContractAddr, eventType, eventData)
	log.Error("Failed to " + callFn, "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}
