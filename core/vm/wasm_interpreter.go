package vm

import (
	"errors"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/params"
)

var (
	ErrWASMWriteProtection          = errors.New("WASM: write protection")
	ErrWASMMigrate                  = errors.New("WASM: failed to migrate contract")
	ErrWASMEventCountToLarge        = errors.New("WASM: event indexed count too large")
	ErrWASMEventContentToLong       = errors.New("WASM: event indexed content too long")
	ErrWASMSha3DstToShort           = errors.New("WASM: sha3 dst len too short")
	ErrWASMPanicOp                  = errors.New("WASM: transaction err op")
	ErrWASMOldContractCodeNotExists = errors.New("WASM: old contract code is not exists")
	ErrWASMUndefinedPanic           = errors.New("WASM: vm undefined err")
	ErrWASMRlpItemCountTooLarge     = errors.New("WASM: itemCount too large for RLP")
)

// WASMInterpreter represents an WASM interpreter
type WASMInterpreter struct {
	evm      *EVM
	cfg      Config
	gasTable params.GasTable
}

// NewWASMInterpreter returns a new instance of the Interpreter
func NewWASMInterpreter(evm *EVM, cfg Config) *WASMInterpreter {

	return &WASMInterpreter{
		evm:      evm,
		cfg:      cfg,
		gasTable: evm.ChainConfig().GasTable(evm.BlockNumber),
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
		if r := recover(); r != nil {
			switch e := r.(type) {
			case error:
				ret, err = nil, e
			default:
				ret, err = nil, fmt.Errorf("WASM: execute fail, %v", e)
			}
		}
	}()

	in.evm.depth++
	defer func() {
		in.evm.depth--
	}()

	//// Don't bother with the execution if there's no code.
	//if len(contract.Code) == 0 {
	//	return nil, nil
	//}

	creator, err := NewWasmEngineCreator(in.cfg.WasmType)
	if err != nil {
		return nil, err
	}

	engine, err := creator.Create(in.evm, in.cfg, in.gasTable, contract)
	if err != nil {
		return nil, err
	}

	ret, err = engine.Run(input, readOnly)
	if err != nil {
		return ret, err
	}

	return ret, nil
}

// CanRun tells if the contract, passed as an argument, can be run
// by the current interpreter
func (in *WASMInterpreter) CanRun(code []byte) bool {
	return CanUseWASMInterp(code)
}

type WasmInsType byte

const (
	Unknown WasmInsType = 0x00
	Wagon   WasmInsType = 0x01
)

func (t WasmInsType) String() string {
	switch t {
	case Wagon:
		return "wagon"
	default:
		return "unknown"
	}
}

func Str2WasmType(str string) WasmInsType {
	switch str {
	case "wagon":
		return Wagon
	default:
		return Unknown
	}
}

var creators = map[WasmInsType]wasmEngineCreator{
	Wagon: &wagonEngineCreator{},
}

func NewWasmEngineCreator(vm WasmInsType) (wasmEngineCreator, error) {

	if creator, ok := creators[vm]; ok {
		return creator, nil
	}
	return nil, fmt.Errorf("unsupport wasm type: %d", vm)
}

type wasmEngineCreator interface {
	Create(evm *EVM, config Config, gasTable params.GasTable, contract *Contract) (wasmEngine, error)
}
