package vm

import (
	"Platon-go/life/exec"
	"Platon-go/life/resolver"
	"bytes"
	"encoding/binary"
	"fmt"
)

// WASM解释器，用于负责解析WASM指令集，具体执行将委托至Life虚拟机完成
// 实现Interpreter的接口 run/canRun.
// WASMInterpreter represents an WASM interpreter
type WASMInterpreter struct {
	evm 		*EVM
	cfg 		Config
	vmContext 	*exec.VMContext
	lvm 		*exec.VirtualMachine

	resolver 	exec.ImportResolver
	returnData	[]byte
}

// NewWASMInterpreter returns a new instance of the Interpreter
func NewWASMInterpreter(evm *EVM, cfg Config) *WASMInterpreter {

	return &WASMInterpreter{
		evm : evm,
		cfg : cfg,
		vmContext: &exec.VMContext{
			Config: exec.VMConfig{
				EnableJIT: false,
				DefaultMemoryPages: 128,
				DynamicMemoryPages: 1,
			},
			Addr: [20]byte{},
			GasUsed : 0,
			GasLimit: evm.Context.GasLimit,
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

	// 解释器的执行委托给Life虚拟机
	// 执行过程首先创建一个新的VirtualMachine.

	in.evm.depth++
	defer func(){ in.evm.depth-- }()

	if len(contract.Code) == 0 {
		return nil, nil
	}

	if len(contract.ABI) == 0 {
		return nil,nil
	}

	// 获取执行器对象
	in.lvm, err = exec.NewVirtualMachine(contract.Code, *in.vmContext, in.resolver,nil)
	if err != nil {
		return nil, err
	}

	// input 代表着交易的data, 需要从中解析出entryPoint.
	contract.Input = input

	// 1、通过input解析出方法名；
	// 2、根据ABI解析出参数类型，从input中获取对应参数;
	// 3、获取entryID，
	// 4、执行调用

	// for test
	funcName := "transfer"
	entryID, ok := in.lvm.GetFunctionExport(funcName)
	if !ok {
		return nil, fmt.Errorf("entryId not found.")
	}

	// todo: 此处暂时未测试点
	params := []int64{
			resolver.MallocString(in.lvm,"hello "),
			resolver.MallocString(in.lvm,"world"),
			45,
		}

	res, err := in.lvm.Run(entryID, params...)
	if err != nil {
		in.lvm.PrintStackTrace()
		return nil, err
	}
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















