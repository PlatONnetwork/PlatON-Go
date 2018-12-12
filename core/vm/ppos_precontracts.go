package vm

import (
	"Platon-go/common"
	"Platon-go/common/byteutil"
	"Platon-go/core/types"
	"Platon-go/crypto"
	"Platon-go/log"
	"Platon-go/rlp"
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
)

// Ppos pre-compiled contract address
var PrecompiledContractsPpos = map[common.Address]PrecompiledContract{
	common.CandidatePoolAddr : &candidateContract{},
	common.TicketPoolAddr : &ticketContract{},
}

// error def
var (
	ErrParamsRlpDecode = errors.New("Rlp decode faile")
	ErrParamsBaselen = errors.New("Params Base length does not match")
	ErrParamsLen = errors.New("Params length does not match")
	ErrUndefFunction = errors.New("Undefined function")
	ErrCallRecode = errors.New("Call recode error, panic...")
)

// execute decode input data and call the function
func execute(input []byte, command map[string]interface{}) ([]byte, error) {
	// debug
	logError("Run==> ", "input: ", hex.EncodeToString(input))
	defer func() {
		if err := recover(); nil != err {
			// catch call panic
			logError("Run==> ", "ErrCallRecode: ", ErrCallRecode.Error())
		}
	}()
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

type ResultCommon struct {
	Ret bool
	ErrMsg string
}

// transaction add event
func (c *candidateContract) addLog(event, data string) {
	var logdata [][]byte
	logdata = make([][]byte, 0)
	logdata = append(logdata, []byte(data))
	buf := new(bytes.Buffer)
	if err := rlp.Encode(buf, logdata); err!=nil {
		logError("addlog==> ","rlp encode fail: ", err.Error())
	}
	c.evm.StateDB.AddLog(&types.Log{
		Address:common.CandidatePoolAddr,
		Topics: []common.Hash{common.BytesToHash(crypto.Keccak256([]byte(event)))},
		Data: buf.Bytes(),
		BlockNumber: c.evm.Context.BlockNumber.Uint64(),
	})
}

// debug log
func logInfo(msg string, ctx ...interface{})  {
	log.Info(msg, ctx...)
	//args := []interface{}{msg}
	//args = append(args, ctx...)
	//fmt.Println(args...)
	/*if c.evm.vmConfig.ConsoleOutput {
		//console output
		args := []interface{}{msg}
		args = append(args, ctx...)
		fmt.Println(args...)
	}else {
		//log output
		log.Info(msg, ctx...)
	}*/
}
func logError(msg string, ctx ...interface{})  {
	log.Error(msg, ctx...)
	//args := []interface{}{msg}
	//args = append(args, ctx...)
	//fmt.Println(args...)
	/*if c.evm.vmConfig.ConsoleOutput {
		//console output
		args := []interface{}{msg}
		args = append(args, ctx...)
		fmt.Println(args...)
	}else {
		//log output
		log.Error(msg, ctx...)
	}*/
}
func (c *candidateContract) logPrint(level log.Lvl, msg string, ctx ...interface{})  {
	if 	c.evm.vmConfig.ConsoleOutput {
		//console output
		args := make([]interface{}, len(ctx)+1)
		args[0] = msg
		for i, v := range ctx{
			args[i+1] = v
		}
		fmt.Println(args...)
	}else {
		//log output
		switch level {
		case log.LvlCrit:
			log.Crit(msg, ctx...)
		case log.LvlError:
			log.Error(msg, ctx...)
		case log.LvlWarn:
			log.Warn(msg, ctx...)
		case log.LvlInfo:
			log.Info(msg, ctx...)
		case log.LvlDebug:
			log.Debug(msg, ctx...)
		case log.LvlTrace:
			log.Trace(msg, ctx...)
		}
	}
}

// return string format
func DecodeResultStr (result string) []byte {
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