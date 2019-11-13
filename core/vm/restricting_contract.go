package vm

import (
	"math/big"
	"strconv"

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
	return execPlatonContract(input, rc.FnSigns())
}

func (rc *RestrictingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		4000: rc.createRestrictingPlan,

		// Get
		4100: rc.getRestrictingInfo,
	}
}

func (rc *RestrictingContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	return nil
}

// createRestrictingPlan is a PlatON precompiled contract function, used for create a restricting plan
func (rc *RestrictingContract) createRestrictingPlan(account common.Address, plans []restricting.RestrictingPlan) ([]byte, error) {

	//sender := rc.Contract.Caller()
	from := rc.Contract.CallerAddress
	txHash := rc.Evm.StateDB.TxHash()
	blockNum := rc.Evm.BlockNumber
	blockHash := rc.Evm.BlockHash
	state := rc.Evm.StateDB

	log.Debug("Call createRestrictingPlan of RestrictingContract", "blockNumber", blockNum.Uint64(),
		"blockHash", blockHash.TerminalString(), "txHash", txHash.Hex(), "from", from.String(), "account", account.String())

	if !rc.Contract.UseGas(params.CreateRestrictingPlanGas) {
		return nil, ErrOutOfGas
	}
	if !rc.Contract.UseGas(params.ReleasePlanGas * uint64(len(plans))) {
		return nil, ErrOutOfGas
	}
	if txHash == common.ZeroHash {
		return nil, nil
	}

	err := rc.Plugin.AddRestrictingRecord(from, account, blockNum.Uint64(), plans, state)
	switch err.(type) {
	case nil:
		receipt := strconv.Itoa(int(common.NoErr.Code))
		rc.goodLog(CreateRestrictingPlanEvent, receipt, "createRestrictingPlan")
		return []byte(receipt), nil
	case *common.BizError:
		bizErr := err.(*common.BizError)
		receipt := strconv.Itoa(int(bizErr.Code))
		rc.badLog(CreateRestrictingPlanEvent, receipt, bizErr.Msg, "createRestrictingPlan")
		return []byte(receipt), nil
	default:
		log.Error("Failed to cal addRestrictingRecord on createRestrictingPlan", "blockNumber", blockNum.Uint64(),
			"blockHash", blockHash.TerminalString(), "txHash", txHash.Hex(), "error", err)
		return nil, err
	}
}

// createRestrictingPlan is a PlatON precompiled contract function, used for getting restricting info.
// first output param is a slice of byte of restricting info;
// the secend output param is the result what plugin executed GetRestrictingInfo returns.
func (rc *RestrictingContract) getRestrictingInfo(account common.Address) ([]byte, error) {
	currNumber := rc.Evm.BlockNumber
	state := rc.Evm.StateDB

	log.Debug("Call getRestrictingInfo of RestrictingContract", "blockNumber", currNumber.Uint64(), "account", account.String())

	result, err := rc.Plugin.GetRestrictingInfo(account, state)
	if err != nil {
		return xcom.NewFailedResult(err), nil
	} else {
		return xcom.NewOkResult(string(result)), nil
	}
}

func (rc *RestrictingContract) goodLog(eventType, eventData, callFn string) {

	blockNumber := rc.Evm.BlockNumber.Uint64()
	xcom.AddLog(rc.Evm.StateDB, blockNumber, vm.RestrictingContractAddr, eventType, eventData)
}

func (rc *RestrictingContract) badLog(eventType, eventData, reason, callFn string) {

	txHash := rc.Evm.StateDB.TxHash()
	blockNumber := rc.Evm.BlockNumber.Uint64()
	xcom.AddLog(rc.Evm.StateDB, blockNumber, vm.RestrictingContractAddr, eventType, eventData)
	log.Error("Failed to "+callFn, "txHash", txHash, "blockNumber", blockNumber, "receipt: ", eventData, "reason", reason)
}
