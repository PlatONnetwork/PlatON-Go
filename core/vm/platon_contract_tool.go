// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

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
		log.Error("Failed to verify contract tx before exec", "err", err)
		return xcom.NewResult(common.InvalidParameter, nil), err
	}

	// execute contracts method
	result := reflect.ValueOf(fn).Call(params)
	switch errtyp := result[1].Interface().(type) {
	case *common.BizError:
		log.Error("Failed to execute contract tx", "err", err)
		return xcom.NewResult(errtyp, nil), errtyp
	case error:
		log.Error("Failed to execute contract tx", "err", err)
		return xcom.NewResult(common.InternalError, nil), errtyp
	default:
	}
	return result[0].Bytes(), nil
}

func txResultHandler(contractAddr common.Address, evm *EVM, title, reason string, fncode int, errCode *common.BizError) ([]byte, error) {
	event := strconv.Itoa(fncode)
	receipt := strconv.Itoa(int(errCode.Code))
	blockNumber := evm.BlockNumber.Uint64()
	if errCode.Code != 0 {
		txHash := evm.StateDB.TxHash()
		log.Error("Failed to "+title, "txHash", txHash.Hex(),
			"blockNumber", blockNumber, "receipt: ", receipt, "the reason", reason)
	}
	xcom.AddLog(evm.StateDB, blockNumber, contractAddr, event, receipt)
	if errCode.Code == common.NoErr.Code {
		return []byte(receipt), nil
	}
	return []byte(receipt), errCode
}

func txResultHandlerWithRes(contractAddr common.Address, evm *EVM, title, reason string, fncode, errCode int, res interface{}) []byte {
	event := strconv.Itoa(fncode)
	receipt := strconv.Itoa(errCode)
	blockNumber := evm.BlockNumber.Uint64()
	if errCode != 0 {
		txHash := evm.StateDB.TxHash()
		log.Error("Failed to "+title, "txHash", txHash.Hex(),
			"blockNumber", blockNumber, "receipt: ", receipt, "the reason", reason)
	}
	xcom.AddLogWithRes(evm.StateDB, blockNumber, contractAddr, event, receipt, res)
	return []byte(receipt)
}

func callResultHandler(evm *EVM, title string, resultValue interface{}, err *common.BizError) []byte {
	txHash := evm.StateDB.TxHash()
	blockNumber := evm.BlockNumber.Uint64()

	if nil != err {
		log.Error("Failed to "+title, "txHash", txHash.Hex(),
			"blockNumber", blockNumber, "the reason", err.Error())
		return xcom.NewResult(err, nil)
	}

	if IsBlank(resultValue) {
		return xcom.NewResult(common.NotFound, nil)
	}

	log.Debug("Call "+title+" finished", "blockNumber", blockNumber,
		"txHash", txHash, "result", resultValue)
	return xcom.NewResult(nil, resultValue)
}

func IsBlank(i interface{}) bool {
	defer func() {
		recover()
	}()

	typ := reflect.TypeOf(i)
	val := reflect.ValueOf(i)
	if typ == nil {
		return true
	} else {
		if typ.Kind() == reflect.Slice {
			return val.Len() == 0
		}
		if typ.Kind() == reflect.Map {
			return val.Len() == 0
		}
	}
	return val.IsNil()
}

func checkInputEmpty(input []byte) bool {
	if len(input) == 0 {
		return true
	} else {
		return false
	}
}
