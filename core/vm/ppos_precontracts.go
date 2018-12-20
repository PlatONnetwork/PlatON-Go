package vm

import (
	"Platon-go/common"
	"Platon-go/common/byteutil"
	"Platon-go/log"
	"Platon-go/rlp"
	"bytes"
	"encoding/hex"
	"errors"
	"reflect"
)

// Ppos pre-compiled contract address
var PrecompiledContractsPpos = map[common.Address]PrecompiledContract{
	common.CandidatePoolAddr: &candidateContract{},
	common.TicketPoolAddr:    &ticketContract{},
}

// error def
var (
	ErrParamsRlpDecode = errors.New("Rlp decode fail")
	ErrParamsBaselen   = errors.New("Params Base length does not match")
	ErrParamsLen       = errors.New("Params length does not match")
	ErrUndefFunction   = errors.New("Undefined function")
	ErrCallRecode      = errors.New("Call recode error, panic...")
)

// execute decode input data and call the function
func execute(input []byte, command map[string]interface{}) ([]byte, error) {
	// debug
	log.Error("Run==> ", "input: ", hex.EncodeToString(input))
	defer func() {
		if err := recover(); nil != err {
			// catch call panic
			log.Error("Run==> ", "ErrCallRecode: ", ErrCallRecode.Error())
		}
	}()
	var source [][]byte
	if err := rlp.Decode(bytes.NewReader(input), &source); nil != err {
		log.Error("Run==> ", err.Error())
		return nil, ErrParamsRlpDecode
	}
	// check
	if len(source) < 2 {
		log.Error("Run==> ", "ErrParamsBaselen: ", ErrParamsBaselen.Error())
		return nil, ErrParamsBaselen
	}
	// get func and param list
	if _, ok := command[byteutil.BytesToString(source[1])]; !ok {
		log.Error("Run==> ", "ErrUndefFunction: ", ErrUndefFunction.Error())
		return nil, ErrUndefFunction
	}
	funcValue := command[byteutil.BytesToString(source[1])]
	paramList := reflect.TypeOf(funcValue)
	paramNum := paramList.NumIn()
	// var param []interface{}
	params := make([]reflect.Value, paramNum)
	if paramNum != len(source)-2 {
		log.Error("Run==> ", "ErrParamsLen: ", ErrParamsLen.Error())
		return nil, ErrParamsLen
	}
	for i := 0; i < paramNum; i++ {
		targetType := paramList.In(i).String()
		originByte := []reflect.Value{reflect.ValueOf(source[i+2])}
		params[i] = reflect.ValueOf(byteutil.Command[targetType]).Call(originByte)[0]
	}
	// call func
	result := reflect.ValueOf(funcValue).Call(params)
	log.Info("Run==> ", "result[0]: ", result[0].Bytes())
	if _, err := result[1].Interface().(error); !err {
		return result[0].Bytes(), nil
	}
	log.Info(result[1].Interface().(error).Error())
	return result[0].Bytes(), result[1].Interface().(error)
}

type ResultCommon struct {
	Ret    bool
	ErrMsg string
}

// return string format
func DecodeResultStr(result string) []byte {
	// 0x0000000000000000000000000000000000000020
	// 00000000000000000000000000000000000000000d
	// 00000000000000000000000000000000000000000

	resultBytes := []byte(result)
	strHash := common.BytesToHash(common.Int32ToBytes(32))
	sizeHash := common.BytesToHash(common.Int64ToBytes(int64((len(resultBytes)))))
	var dataRealSize = len(resultBytes)
	if (dataRealSize % 32) != 0 {
		dataRealSize = dataRealSize + (32 - (dataRealSize % 32))
	}
	dataByt := make([]byte, dataRealSize)
	copy(dataByt[0:], resultBytes)

	finalData := make([]byte, 0)
	finalData = append(finalData, strHash.Bytes()...)
	finalData = append(finalData, sizeHash.Bytes()...)
	finalData = append(finalData, dataByt...)
	//encodedStr := hex.EncodeToString(finalData)
	//fmt.Println("finalData: ", encodedStr)
	return finalData
}
