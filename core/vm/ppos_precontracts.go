package vm

import (
	"Platon-go/common"
	"Platon-go/core/types"
	"Platon-go/crypto"
	"Platon-go/log"
	"Platon-go/rlp"
	"bytes"
	"fmt"
)

// Ppos pre-compiled contract address
var PrecompiledContractsPpos = map[common.Address]PrecompiledContract{
	common.CandidatePoolAddr : &candidateContract{},
	common.TicketPoolAddr : &ticketContract{},
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