package vm

import (
	"Platon-go/common/byteutil"
	"Platon-go/params"
	"Platon-go/rlp"
	"bytes"
	"encoding/hex"
	"reflect"
)

type ticketContract struct {
	contract *Contract
	evm *EVM
}

func (t *ticketContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (t *ticketContract) Run(input []byte) ([]byte, error) {

	// debug
	logError("Run==> ", "input: ", hex.EncodeToString(input))

	defer func() {
		if err := recover(); nil != err {
			// catch call panic
			logError("Run==> ", "ErrCallRecode: ", ErrCallRecode.Error())
		}
	}()
	var command = map[string] interface{}{
		// 接口列表
	}
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &source); err != nil {
		logError("Run==> ", err.Error())
		return nil, ErrParamsRlpDecode
	}
	// check
	if len(source)<2 {
		logError("Run==> ", "ErrParamsBaselen: ", ErrParamsBaselen.Error())
		return nil, ErrParamsBaselen
	}
	// get func and param list
	if _, ok := command[byteutil.BytesToString(source[1])]; !ok {
		logError("Run==> ", "ErrUndefFunction: ", ErrUndefFunction.Error())
		return nil, ErrUndefFunction
	}
	funcValue := command[byteutil.BytesToString(source[1])]
	paramList := reflect.TypeOf(funcValue)
	paramNum := paramList.NumIn()
	// var param []interface{}
	params := make([]reflect.Value, paramNum)
	if paramNum!=len(source)-2 {
		logError("Run==> ", "ErrParamsLen: ",ErrParamsLen.Error())
		return nil, ErrParamsLen
	}
	for i := 0; i < paramNum; i++ {
		targetType := paramList.In(i).String()
		originByte := []reflect.Value{reflect.ValueOf(source[i+2])}
		params[i] = reflect.ValueOf(byteutil.Command[targetType]).Call(originByte)[0]
	}
	// call func
	result := reflect.ValueOf(funcValue).Call(params)
	logInfo("Run==> ", "result[0]: ", result[0].Bytes())
	if _, err := result[1].Interface().(error); !err {
		return result[0].Bytes(), nil
	}
	logInfo(result[1].Interface().(error).Error())
	return result[0].Bytes(), result[1].Interface().(error)
}