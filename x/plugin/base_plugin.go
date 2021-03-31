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

package plugin

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	gerr "github.com/go-errors/errors"
)

type BasePlugin interface {
	BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error
	EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) error
	Confirmed(nodeId discover.NodeID, block *types.Block) error
}

var (
	DecodeTxDataErr = errors.New("decode tx data is err")
	FuncNotExistErr = errors.New("the func is not exist")
	FnParamsLenErr  = errors.New("the params len and func params len is not equal")
)

func VerifyTxData(input []byte, command map[uint16]interface{}) (cnCode uint16, fn interface{}, FnParams []reflect.Value, err error) {

	defer func() {
		if er := recover(); nil != er {
			fn, FnParams, err = nil, nil, fmt.Errorf("parse tx data is failed: %s", er)
			log.Error("Failed to Verify PlatON inner contract tx data", "error",
				fmt.Errorf("the err stack: %s", gerr.Wrap(er, 2).ErrorStack()))
		}
	}()

	var args [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &args); nil != err {
		log.Error("Failed to Verify PlatON inner contract tx data, Decode rlp input failed", "err", err)
		return 0, nil, nil, fmt.Errorf("%v: %v", DecodeTxDataErr, err)
	}

	//fmt.Println("the Function Type:", byteutil.BytesToUint16(args[0]))

	fnCode := byteutil.BytesToUint16(args[0])
	if fn, ok := command[fnCode]; !ok {
		return 0, nil, nil, FuncNotExistErr
	} else {

		//funcName := runtime.FuncForPC(reflect.ValueOf(fn).Pointer()).Name()
		//fmt.Println("The FuncName is", funcName)

		// the func params type list
		paramList := reflect.TypeOf(fn)
		// the func params len
		paramNum := paramList.NumIn()

		if paramNum != len(args)-1 {
			return 0, nil, nil, FnParamsLenErr
		}

		params := make([]reflect.Value, paramNum)

		for i := 0; i < paramNum; i++ {
			//fmt.Println("byte:", args[i+1])

			targetType := paramList.In(i).String()
			inputByte := []reflect.Value{reflect.ValueOf(args[i+1])}
			params[i] = reflect.ValueOf(byteutil.Bytes2X_CMD[targetType]).Call(inputByte)[0]
			//fmt.Println("num", i+1, "type", targetType)
		}
		return fnCode, fn, params, nil
	}
}
