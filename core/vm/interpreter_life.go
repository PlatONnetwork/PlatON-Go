package vm

import (
	"Platon-go/life/utils"
	"Platon-go/rlp"
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"

	"Platon-go/life/exec"
	"Platon-go/life/resolver"
)

// WASM解释器，用于负责解析WASM指令集，具体执行将委托至Life虚拟机完成
// 实现Interpreter的接口 run/canRun.
// WASMInterpreter represents an WASM interpreter
type WASMInterpreter struct {
	vmContext 	*exec.VMContext
	lvm 		*exec.VirtualMachine

	resolver 	exec.ImportResolver
	returnData	[]byte
}

// NewWASMInterpreter returns a new instance of the Interpreter
func NewWASMInterpreter(evm *EVM, cfg Config) *WASMInterpreter {

	// 初始化WASM解释器，保存WASM虚拟机需要的配置及上下文信息
	return &WASMInterpreter{
		vmContext: &exec.VMContext{
			Config: exec.VMConfig{
				EnableJIT: false,
				DefaultMemoryPages: 128,
				DynamicMemoryPages: 1,
			},
			Addr: [20]byte{},
			GasUsed : 0,
			GasLimit: evm.Context.GasLimit,
			// 验证此处是否可行
			StateDB: evm.StateDB,
		},
		resolver : resolver.NewResolver(0x01),
	}
}

// Run loops and evaluates the contract's code with the given input data and returns.
// the return byte-slice and an error if one occurred
//
// It's important to note that any errors returned by the interpreter should be
// considered a revert-and-consume-all-gas operations except for
// errExecutionReverted which means revert-and-keep-gas-lfet.
func (in *WASMInterpreter) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {

	//in.vmContext.Evm.depth++
	//defer func(){ in.vmContext.Evm.depth-- }()

	if len(contract.Code) == 0 {
		return nil, nil
	}

	if len(contract.ABI) == 0 {
		return nil,nil
	}

	in.vmContext.Addr = contract.Address()
	in.vmContext.GasLimit = contract.Gas		// 可使用的即为受限制的
	//in.vmContext.Contract = contract

	// 获取执行器对象
	in.lvm, err = exec.NewVirtualMachine(contract.Code, *in.vmContext, in.resolver,nil)
	if err != nil {
		return nil, err
	}

	// input 代表着交易的data, 需要从中解析出entryPoint.
	contract.Input = input
	var (
		funcName string
		//txType	int		// 交易类型：合约创建、交易、投票等类型
		params 	[]int64
	)

	if input == nil {
		funcName = "init"	// init function.
	} else {
		// parse input.
		_, funcName, params, err = parseInputFromAbi(in.lvm, input, contract.ABI)
		if err != nil {
			return nil, err
		}
	}
	entryID, ok := in.lvm.GetFunctionExport(funcName)
	if !ok {
		return nil, fmt.Errorf("entryId not found.")
	}
	res, err := in.lvm.RunWithGasLimit(entryID,int(in.vmContext.GasLimit), params...)
	if err != nil {
		in.lvm.PrintStackTrace()
		return nil, err
	}
	if contract.Gas > in.vmContext.GasUsed {
		contract.Gas = contract.Gas - in.vmContext.GasUsed
	} else {
		return nil, fmt.Errorf("out of gas.")
	}

	//todo: 问题点，需解决
	in.returnData = Int64ToBytes(res)
	return Int64ToBytes(res),nil
}

// CanRun tells if the contract, passed as an argument, can be run
// by the current interpreter
func (in *WASMInterpreter) CanRun(code []byte) bool {
	return true
}

func Int64ToBytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf,uint64(i))
	return buf
}

func BytesToInt64(bys []byte) int64 {
	buf := bytes.NewBuffer(bys)
	var res int64
	binary.Read(buf, binary.BigEndian, &res)
	return res
}

// parse input(payload)
func parseInputFromAbi(vm *exec.VirtualMachine, input []byte, abi []byte) (txType int, funcName string, params []int64, err error) {
	if input == nil || len(input) <= 1 {
		return -1,"",nil, fmt.Errorf("invalid input.")
	}
	// [txType][msg.to][funcName][args1][args2]
	// rlp decode
	ptr := new(interface{})
	err = rlp.Decode(bytes.NewReader(input), &ptr)
	if err != nil {
		return -1, "", nil, err
	}
	rlpList := reflect.ValueOf(ptr).Elem().Interface();

	if _, ok := rlpList.([]interface{}); !ok {
		return -1, "", nil, fmt.Errorf("invalid rlp format.")
	}

	iRlpList := rlpList.([]interface{})
	if len(iRlpList) <= 2 {
		return -1, "", nil, fmt.Errorf("invalid input. ele must greater than 2")
	}

	wasmabi := new(utils.WasmAbi)
	err = wasmabi.FromJson(abi)
	if err != nil {
		return -1, "", nil, fmt.Errorf("invalid abi, encoded fail.")
	}

	params = make([]int64, 0)
	if v, ok := iRlpList[0].([]byte); ok {
		txType = int(v[0])
	}
	if v, ok := iRlpList[1].([]byte); ok {
		funcName = string(v)
	}

	// 查找方法名对应的args
	var args []utils.Args
	for _, v := range wasmabi.Abi {
		if strings.EqualFold(funcName, v.Method) {
			args = v.Args
			break
		}
	}
	if len(args) == 0 {
		return -1, "", nil, fmt.Errorf("no match abi args by funcName.")
	}

	argsRlp := iRlpList[2:]
	if len(args) != len(argsRlp) {
		return -1, "", nil, fmt.Errorf("invalid input or invalid abi.")
	}
	// todo: abi类型解析，需要继续添加
	// uint64 uint32  uint16 uint8 int64 int32  int16 int8 float32 float64 string void
	// 此处参数是否替换为uint64
	for i, v := range args {
		bts := argsRlp[i].([]byte)
		switch v.RealTypeName {
		case "string":
			pos := resolver.MallocString(vm, string(bts))
			params = append(params, pos)
		case "int8":
			params = append(params, int64(bts[0]))
		case "int16":
			params = append(params, int64(binary.BigEndian.Uint16(bts)))
		case "int32":
			params = append(params, int64(binary.BigEndian.Uint32(bts)))
		case "int64":
			params = append(params, int64(binary.BigEndian.Uint64(bts)))
		case "uint8":
			params = append(params, int64(bts[0]))
		case "uint32":
			params = append(params, int64(binary.BigEndian.Uint32(bts)))
		case "uint64":
			params = append(params, int64(binary.BigEndian.Uint64(bts)))
		case "bool":
			params = append(params, int64(bts[0]))
		}
	}

	return txType, funcName, params, nil
}















