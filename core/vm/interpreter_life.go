package vm

import (
	"Platon-go/life/exec"
	"Platon-go/life/resolver"
	"Platon-go/log"
)

const (
	CALL_CANTRACT_FLAG = 9
)

var DEFAULT_VM_CONFIG = exec.VMConfig {
	EnableJIT:          false,
	DefaultMemoryPages: 512,
	DynamicMemoryPages: 5,
}

// WASMInterpreter represents an WASM interpreter
type WASMInterpreter struct {
	evm       	*EVM
	cfg       	Config
	wasmStateDB *WasmStateDB
	WasmLogger log.Logger
	resolver   exec.ImportResolver
	returnData []byte
}

// NewWASMInterpreter returns a new instance of the Interpreter
func NewWASMInterpreter(evm *EVM, cfg Config) *WASMInterpreter {
	ws := &WasmStateDB{
		StateDB: evm.StateDB,
		evm:     evm,
		cfg:     &cfg,
	}
	return &WASMInterpreter{
		evm: evm,
		cfg: cfg,
		WasmLogger: NewWasmLogger(cfg, log.WasmRoot()),
		wasmStateDB: ws,
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
	panic("No supported life's vm interpreter yet.")
}

// CanRun tells if the contract, passed as an argument, can be run
// by the current interpreter
func (in *WASMInterpreter) CanRun(code []byte) bool {
	return true
}

// parse input(payload)
func parseInputFromAbi(vm *exec.VirtualMachine, input []byte, abi []byte) (txType int, funcName string, params []int64, returnType string, err error) {
	//TODO: to retrive params from input
	return 0, "", nil, "", nil
}

func parseRlpData(rlpData []byte) (int64, []byte, []byte, error) {
	//TODO: to retrive params from input
	return 0, nil, nil, nil
}
