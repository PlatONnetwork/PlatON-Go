package exec

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/wagon/disasm"
	"github.com/PlatONnetwork/wagon/exec/internal/compile"
	"github.com/PlatONnetwork/wagon/wasm"
)

type CompileVM struct {
	VM
}

type CompiledModule struct {
	RawModule *wasm.Module
	globals   []uint64
	memory    []byte
	funcs     []function
}

func CompileModule(module *wasm.Module) (*CompiledModule, error) {
	var compiled CompiledModule

	if module.Memory != nil && len(module.Memory.Entries) != 0 {
		if len(module.Memory.Entries) > 1 {
			return nil, ErrMultipleLinearMemories
		}

		memsize := uint(module.Memory.Entries[0].Limits.Initial) * wasmPageSize
		compiled.memory = make([]byte, memsize)
		copy(compiled.memory, module.LinearMemoryIndexSpace[0])
	}

	compiled.funcs = make([]function, len(module.FunctionIndexSpace))
	compiled.globals = make([]uint64, len(module.GlobalIndexSpace))
	compiled.RawModule = module

	nNatives := 0
	for i, fn := range module.FunctionIndexSpace {
		// Skip native methods as they need not be
		// disassembled; simply add them at the end
		// of the `funcs` array as is, as specified
		// in the spec. See the "host functions"
		// section of:
		// https://webassembly.github.io/spec/core/exec/modules.html#allocation
		if fn.IsHost() {
			compiled.funcs[i] = goFunction{
				typ: fn.Host.Type(),
				val: fn.Host,
			}
			nNatives++
			continue
		}

		disassembly, err := disasm.NewDisassembly(fn, module)
		if err != nil {
			return nil, err
		}

		totalLocalVars := 0
		totalLocalVars += len(fn.Sig.ParamTypes)
		for _, entry := range fn.Body.Locals {
			totalLocalVars += int(entry.Count)
		}
		code, meta := compile.Compile(disassembly.Code)
		compiled.funcs[i] = compiledFunction{
			code:           code,
			branchTables:   meta.BranchTables,
			maxDepth:       disassembly.MaxDepth,
			totalLocalVars: totalLocalVars,
			args:           len(fn.Sig.ParamTypes),
			returns:        len(fn.Sig.ReturnTypes) != 0,
		}
	}

	for i, global := range module.GlobalIndexSpace {
		val, err := module.ExecInitExpr(global.Init)
		if err != nil {
			return nil, err
		}
		switch v := val.(type) {
		case int32:
			compiled.globals[i] = uint64(v)
		case int64:
			compiled.globals[i] = uint64(v)
			//case float32:
			//	compiled.globals[i] = uint64(math.Float32bits(v))
			//case float64:
			//	compiled.globals[i] = uint64(math.Float64bits(v))
		}
	}

	if module.Start != nil {
		//_, err := compiled.ExecCode(int64(module.Start.Index))
		//if err != nil {
		//	return nil, err
		//}
		return nil, errors.New("start entry is not supported in smart contract")
	}

	return &compiled, nil
}

func NewVMWithCompiled(module *CompiledModule, memLimit uint64) (*CompileVM, error) {
	var vm CompileVM

	memsize := len(module.memory)
	if uint64(memsize) > memLimit {
		return nil, fmt.Errorf("memory is exceed the limitation of %d", memLimit)
	}
	vm.MemoryLimitation = memLimit
	vm.memory = make([]byte, memsize)
	copy(vm.memory, module.memory)

	vm.funcs = module.funcs
	vm.globals = make([]uint64, len(module.RawModule.GlobalIndexSpace))
	copy(vm.globals, module.globals)
	vm.newFuncTable()
	vm.module = module.RawModule

	return &vm, nil
}

