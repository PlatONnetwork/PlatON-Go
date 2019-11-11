package vm

import (
	"reflect"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

func execPlatonContract(input []byte, command map[uint16]interface{}) (ret []byte, err error) {

	// verify the tx data by contracts method
	_, fn, params, err := plugin.VerifyTxData(input, command)
	if nil != err {
		log.Error("Failed to verify contract tx", "err", err)
		return nil, err
	}

	// execute contracts method
	result := reflect.ValueOf(fn).Call(params)
	if err, ok := result[1].Interface().(error); ok {
		return nil, err
	}
	return result[0].Bytes(), nil
}

func txResultHandler(contractAddr common.Address, evm *EVM, title, reason string, fncode, errCode int) []byte {
	event := strconv.Itoa(fncode)
	receipt := strconv.Itoa(errCode)

	if errCode == 0 {
		blockNumber := evm.BlockNumber.Uint64()
		xcom.AddLog(evm.StateDB, blockNumber, contractAddr, event, receipt)
	} else {
		txHash := evm.StateDB.TxHash()
		blockNumber := evm.BlockNumber.Uint64()
		xcom.AddLog(evm.StateDB, blockNumber, contractAddr, event, receipt)
		log.Error("Failed to "+title, "txHash", txHash.Hex(),
			"blockNumber", blockNumber, "receipt: ", receipt, "the reason", reason)
	}
	return []byte(receipt)
}

func callResultHandler(evm *EVM, title string, resultType CallResultType, resultValue interface{}, err error) []byte {

	txHash := evm.StateDB.TxHash()
	blockNumber := evm.BlockNumber.Uint64()

	if nil != err {
		log.Error("Failed to "+title, "txHash", txHash.Hex(),
			"blockNumber", blockNumber, "the reason", err.Error())
		resultBytes := xcom.NewFailedResult(err)
		return resultBytes
	}

	if xcom.IsNil(resultValue) {
		if resultType == ResultTypeStruct {
			resultValue = ""
		} else if resultType == ResultTypeSlice {
			resultValue = []string{}
		} else if resultType == ResultTypeMap {
			resultValue = make(map[string]string)
		} else {
			resultValue = ""
		}
	}

	log.Debug("Call "+title+" finished", "blockNumber", blockNumber,
		"txHash", txHash, "result", resultValue)
	resultBytes := xcom.NewOkResult(resultValue)
	return resultBytes
}
