package vm

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/core/lru"
	"github.com/PlatONnetwork/PlatON-Go/life/utils"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"reflect"
	"runtime"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/life/exec"
	"github.com/PlatONnetwork/PlatON-Go/life/resolver"
)

var (
	errReturnInvalidRlpFormat = errors.New("interpreter_life: invalid rlp format.")
	errReturnInsufficientParams = errors.New("interpreter_life: invalid input. ele must greater than 2")
	errReturnInvalidAbi = errors.New("interpreter_life: invalid abi, encoded fail.")
)

const (
	CALL_CANTRACT_FLAG = 9
)

var DEFAULT_VM_CONFIG = exec.VMConfig{
	EnableJIT:          false,
	DefaultMemoryPages: exec.DefaultMemoryPages,
	DynamicMemoryPages: exec.DynamicMemoryPages,
}

// WASMInterpreter represents an WASM interpreter
type WASMInterpreter struct {
	evm         *EVM
	cfg         Config
	wasmStateDB *WasmStateDB
	WasmLogger  log.Logger
	resolver    exec.ImportResolver
	returnData  []byte
}

// NewWASMInterpreter returns a new instance of the Interpreter
func NewWASMInterpreter(evm *EVM, cfg Config) *WASMInterpreter {

	wasmStateDB := &WasmStateDB{
		StateDB: evm.StateDB,
		evm:     evm,
		cfg:     &cfg,
	}
	return &WASMInterpreter{
		evm:         evm,
		cfg:         cfg,
		WasmLogger:  NewWasmLogger(cfg, log.WasmRoot()),
		wasmStateDB: wasmStateDB,
		resolver:    resolver.NewResolver(0x01),
	}
}

// Run loops and evaluates the contract's code with the given input data and returns.
// the return byte-slice and an error if one occurred
//
// It's important to note that any errors returned by the interpreter should be
// considered a revert-and-consume-all-gas operations except for
// errExecutionReverted which means revert-and-keep-gas-lfet.
func (in *WASMInterpreter) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {
	defer func() {
		if er := recover(); er != nil {
			ret, err = nil, fmt.Errorf("VM execute fail：%v", er)
		}
	}()
	in.evm.depth++
	defer func() {
		in.evm.depth--
		if in.evm.depth == 0 {
			logger, ok := in.WasmLogger.(*WasmLogger)
			if ok {
				logger.Flush()
			}
		}
	}()

	if len(contract.Code) == 0 {
		return nil, nil
	}
	_, abi, code, er := parseRlpData(contract.Code)
	if er != nil {
		return nil, er
	}

	context := &exec.VMContext{
		Config:   DEFAULT_VM_CONFIG,
		Addr:     contract.Address(),
		GasLimit: contract.Gas,
		StateDB:  NewWasmStateDB(in.wasmStateDB, contract),
		Log:      in.WasmLogger,
	}

	var lvm *exec.VirtualMachine
	var module *lru.WasmModule
	module, ok := lru.WasmCache().Get(contract.Address())

	if !ok {
		module = &lru.WasmModule{}
		module.Module, module.FunctionCode, err = exec.ParseModuleAndFunc(code, nil)
		if err != nil {
			return nil, err
		}
		lru.WasmCache().Add(contract.Address(), module)
	}

	lvm, err = exec.NewVirtualMachineWithModule(module.Module, module.FunctionCode, context, in.resolver, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		lvm.Stop()
	}()

	contract.Input = input
	var (
		funcName   string
		txType     int
		params     []int64
		returnType string
	)

	if input == nil {
		funcName = "init" // init function.
	} else {
		// parse input.
		txType, funcName, params, returnType, err = parseInputFromAbi(lvm, input, abi)
		if err != nil {
			if err == errReturnInsufficientParams && txType == 0 { // transfer to contract address.
				return nil, nil
			}
			return nil, err
		}
		if txType == 0 {
			return nil, nil
		}
	}
	entryID, ok := lvm.GetFunctionExport(funcName)
	if !ok {
		return nil, fmt.Errorf("entryId not found.")
	}
	res, err := lvm.RunWithGasLimit(entryID, int(context.GasLimit), params...)
	if err != nil {
		fmt.Println("throw exception:", err.Error())
		return nil, err
	}
	if contract.Gas > context.GasUsed {
		contract.Gas = contract.Gas - context.GasUsed
	} else {
		return nil, fmt.Errorf("out of gas.")
	}

	if input == nil {
		return contract.Code, nil
	}

	// todo: more type need to be completed
	switch returnType {
	case "void", "int8", "int", "int32", "int64":
		if txType == CALL_CANTRACT_FLAG {
			return utils.Int64ToBytes(res), nil
		}
		bigRes := new(big.Int)
		bigRes.SetInt64(res)
		finalRes := utils.Align32Bytes(math.U256(bigRes).Bytes())
		return finalRes, nil
	case "uint8", "uint16", "uint32", "uint64":
		if txType == CALL_CANTRACT_FLAG {
			return utils.Uint64ToBytes(uint64(res)), nil
		}
		finalRes := utils.Align32Bytes(utils.Uint64ToBytes((uint64(res))))
		return finalRes, nil
	case "string":
		returnBytes := make([]byte, 0)
		copyData := lvm.Memory.Memory[res:]
		for _, v := range copyData {
			if v == 0 {
				break
			}
			returnBytes = append(returnBytes, v)
		}
		if txType == CALL_CANTRACT_FLAG {
			return returnBytes, nil
		}
		strHash := common.BytesToHash(common.Int32ToBytes(32))
		sizeHash := common.BytesToHash(common.Int64ToBytes(int64((len(returnBytes)))))
		var dataRealSize = len(returnBytes)
		if (dataRealSize % 32) != 0 {
			dataRealSize = dataRealSize + (32 - (dataRealSize % 32))
		}
		dataByt := make([]byte, dataRealSize)
		copy(dataByt[0:], returnBytes)

		finalData := make([]byte, 0)
		finalData = append(finalData, strHash.Bytes()...)
		finalData = append(finalData, sizeHash.Bytes()...)
		finalData = append(finalData, dataByt...)

		//fmt.Println("CallReturn:", string(returnBytes))
		return finalData, nil
	}
	return nil, nil
}

// CanRun tells if the contract, passed as an argument, can be run
// by the current interpreter
func (in *WASMInterpreter) CanRun(code []byte) bool {
	return true
}

// parse input(payload)
func parseInputFromAbi(vm *exec.VirtualMachine, input []byte, abi []byte) (txType int, funcName string, params []int64, returnType string, err error) {
	if input == nil || len(input) <= 1 {
		return -1, "", nil, "", fmt.Errorf("invalid input.")
	}
	// [txType][funcName][args1][args2]
	// rlp decode
	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(input), &ptr)
	if err != nil {
		return -1, "", nil, "", err
	}
	rlpList := reflect.ValueOf(ptr).Elem().Interface()

	if _, ok := rlpList.([]interface{}); !ok {
		return -1, "", nil, "", errReturnInvalidRlpFormat
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) < 2 {
		if len(iRlpList) != 0 {
			if v, ok := iRlpList[0].([]byte); ok {
				txType = int(common.BytesToInt64(v))
			}
		} else {
			txType = -1
		}
		return txType, "", nil, "", errReturnInsufficientParams
	}

	wasmabi := new(utils.WasmAbi)
	err = wasmabi.FromJson(abi)
	if err != nil {
		return -1, "", nil, "", errReturnInvalidAbi
	}

	params = make([]int64, 0)
	if v, ok := iRlpList[0].([]byte); ok {
		txType = int(common.BytesToInt64(v))
	}
	if v, ok := iRlpList[1].([]byte); ok {
		funcName = string(v)
	}

	var args []utils.InputParam
	for _, v := range wasmabi.AbiArr {
		if strings.EqualFold(funcName, v.Name) && strings.EqualFold(v.Type, "function") {
			args = v.Inputs
			if len(v.Outputs) != 0 {
				returnType = v.Outputs[0].Type
			} else {
				returnType = "void"
			}
			break
		}
	}
	argsRlp := iRlpList[2:]
	if len(args) != len(argsRlp) {
		return -1, "", nil, returnType, fmt.Errorf("invalid input or invalid abi.")
	}
	// uint64 uint32  uint16 uint8 int64 int32  int16 int8 float32 float64 string void
	for i, v := range args {
		bts := argsRlp[i].([]byte)
		switch v.Type {
		case "string":
			pos := resolver.MallocString(vm, string(bts))
			params = append(params, pos)
		case "int8":
			params = append(params, int64(bts[0]))
		case "int16":
			params = append(params, int64(binary.BigEndian.Uint16(bts)))
		case "int32", "int":
			params = append(params, int64(binary.BigEndian.Uint32(bts)))
		case "int64":
			params = append(params, int64(binary.BigEndian.Uint64(bts)))
		case "uint8":
			params = append(params, int64(bts[0]))
		case "uint32", "uint":
			params = append(params, int64(binary.BigEndian.Uint32(bts)))
		case "uint64":
			params = append(params, int64(binary.BigEndian.Uint64(bts)))
		case "bool":
			params = append(params, int64(bts[0]))
		}
	}
	return txType, funcName, params, returnType, nil
}

// rlpData=RLP([txType][code][abi])
func parseRlpData(rlpData []byte) (int64, []byte, []byte, error) {
	ptr := new(interface{})
	err := rlp.Decode(bytes.NewReader(rlpData), &ptr)
	if err != nil {
		return -1, nil, nil, err
	}
	rlpList := reflect.ValueOf(ptr).Elem().Interface()

	if _, ok := rlpList.([]interface{}); !ok {
		return -1, nil, nil, fmt.Errorf("invalid rlp format.")
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) <= 2 {
		return -1, nil, nil, fmt.Errorf("invalid input. ele must greater than 2")
	}
	var (
		txType int64
		code   []byte
		abi    []byte
	)
	if v, ok := iRlpList[0].([]byte); ok {
		txType = utils.BytesToInt64(v)
	}
	if v, ok := iRlpList[1].([]byte); ok {
		code = v
		//fmt.Println("dstCode: ", common.Bytes2Hex(code))
	}
	if v, ok := iRlpList[2].([]byte); ok {
		abi = v
		//fmt.Println("dstAbi:", common.Bytes2Hex(abi))
	}
	return txType, abi, code, nil
}

func stack() string {
	var buf [2 << 10]byte
	return string(buf[:runtime.Stack(buf[:], true)])
}
