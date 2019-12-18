package vm

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/wagon/exec"
	"github.com/pkg/errors"
)

const callEntryName = "invoke"
const memoryLimit = 16 * 1014 * 1024

func decodeInput(input []byte) (byte, []byte, error) {
	kind, content, _, err := rlp.Split(input)

	switch {
	case err != nil:
		return 0, nil, err
	case kind != rlp.List:
		return 0, nil, fmt.Errorf("input type error")
	}

	_, vmType, rest, err := rlp.Split(content)
	switch {
	case err != nil:
		return 0, nil, err
	case len(vmType) != 1:
		return 0, nil, fmt.Errorf("vm type error")
	}
	return vmType[0], rest, nil

}

var engines = map[byte]wasmEngineCreator{
	0x1: &wagonEngineCreator{},
}

func NewWasmEngineCreator(vm byte) (wasmEngineCreator, error) {
	if engine, ok := engines[vm]; ok {
		return engine, nil
	}
	return nil, fmt.Errorf("unsupport vm type")
}

type wasmEngineCreator interface {
	Create(evm *EVM, config Config, db StateDB) (wasmEngine, error)
}

type wagonEngineCreator struct {
}

func (w *wagonEngineCreator) Create(evm *EVM, config Config, db StateDB) (wasmEngine, error) {
	return &wagonEngine{
		evm:    evm,
		config: config,
		db:     db,
	}, nil
}

type wasmEngine interface {
	Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error)
}

type wagonEngine struct {
	evm    *EVM
	config Config
	db     StateDB
	vm     *exec.VM
}

func (w *wagonEngine) Run(contract *Contract, input []byte, readOnly bool) (ret []byte, err error) {
	//parse input function, params
	_, data, err := decodeInput(input)
	//load module from contract
	module, err := ReadWasmModule(contract.Code, false)

	vm, err := exec.NewVMWithCompiled(module, memoryLimit)

	//set input bytes
	vm.SetHostCtx(&VMContext{
		evm:    w.evm,
		config: w.config,
		db:     w.db,
		Input:  data,
	})
	//verify function name in module
	entry, ok := module.RawModule.Export.Entries[callEntryName]
	if !ok {
		return nil, nil
	}

	index := int64(entry.Index)

	fidx := module.RawModule.Function.Types[int(index)]

	ftype := module.RawModule.Types.Entries[int(fidx)]

	if len(ftype.ReturnTypes) > 0 {
		return nil, fmt.Errorf("function sig error")
	}

	//exec vm
	_, err = vm.ExecCode(index)
	if err != nil {
		return nil, errors.Wrap(err, "execute function code")
	}

	return ret, err
}
