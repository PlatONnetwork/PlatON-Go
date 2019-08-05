package vm

import (
	"encoding/json"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/restricting"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

const (
	CreateRestrictingPlanEvent = "4000"
)

type RestrictingContract struct {
	Plugin   *plugin.RestrictingPlugin
	Contract *Contract
	Evm      *EVM
}

func (rc *RestrictingContract) RequiredGas(input []byte) uint64 {
	return params.RestrictingPlanGas
}

func (rc *RestrictingContract) Run(input []byte) ([]byte, error) {
	return exec_platon_contract(input, rc.FnSigns())
}

func (rc *RestrictingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		4000: rc.createRestrictingPlan,

		// Get
		4100: rc.getRestrictingInfo,
	}
}

// createRestrictingPlan is a PlatON precompiled contract function, used for create a restricting plan
func (rc *RestrictingContract) createRestrictingPlan(account common.Address, plans []restricting.RestrictingPlan) ([]byte, error) {

	sender := rc.Contract.Caller()
	txHash := rc.Evm.StateDB.TxHash()
	blockNum := rc.Evm.BlockNumber
	state := rc.Evm.StateDB

	log.Info("Call createRestrictingPlan of RestrictingContract", "txHash", txHash.Hex(), "blockNumber", blockNum.Uint64())

	if !rc.Contract.UseGas(params.CreateRestrictingPlanGas) {
		return nil, ErrOutOfGas
	}
	if !rc.Contract.UseGas(params.ReleasePlanGas * uint64(len(plans))) {
		return nil, ErrOutOfGas
	}

	if err := rc.Plugin.AddRestrictingRecord(sender, account, plans, state); err != nil {
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{Status: false, Data: "", ErrMsg: "create restricting plan:" + err.Error()}
			event, _ := json.Marshal(res)
			rc.badLog(state, blockNum.Uint64(), txHash.Hex(), CreateRestrictingPlanEvent, string(event), "createRestrictingPlan")
			return event, nil

		} else {
			log.Debug("AddRestrictingRecord failed to createRestrictingPlan", "txHash", txHash.Hex(), "blockNumber", blockNum.Uint64(), "error", err)
			return nil, err
		}
	}

	res := xcom.Result{Status: true, Data: "", ErrMsg: ""}
	event, _ := json.Marshal(res)
	rc.goodLog(state, blockNum.Uint64(), txHash.Hex(), CreateRestrictingPlanEvent, string(event), "createRestrictingPlan")

	return event, nil
}

// createRestrictingPlan is a PlatON precompiled contract function, used for getting restricting info.
// first output param is a slice of byte of restricting info;
// the secend output param is the result what plugin executed GetRestrictingInfo returns.
func (rc *RestrictingContract) getRestrictingInfo(account common.Address) ([]byte, error) {
	txHash := rc.Evm.StateDB.TxHash()
	currNumber := rc.Evm.BlockNumber
	state := rc.Evm.StateDB

	log.Info("Call getRestrictingInfo of RestrictingContract", "txHash", txHash.Hex(), "blockNumber", currNumber.Uint64())

	result, err := rc.Plugin.GetRestrictingInfo(account, state)
	var res xcom.Result
	if err != nil {
		res.Status = false
		res.Data = ""
		res.ErrMsg = "get restricting info:" + err.Error()
	} else {
		res.Status = true
		res.Data = string(result)
		res.ErrMsg = "ok"
	}
	return json.Marshal(res)
}

func (rc *RestrictingContract) goodLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.RestrictingContractAddr, eventType, eventData)
	log.Info("Successed to "+callFn, "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}

func (rc *RestrictingContract) badLog(state xcom.StateDB, blockNumber uint64, txHash, eventType, eventData, callFn string) {
	xcom.AddLog(state, blockNumber, vm.RestrictingContractAddr, eventType, eventData)
	log.Debug("Failed to "+callFn, "txHash", txHash, "blockNumber", blockNumber, "json: ", eventData)
}
