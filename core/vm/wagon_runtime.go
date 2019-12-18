package vm

import (
	"github.com/PlatONnetwork/wagon/exec"
	"github.com/PlatONnetwork/wagon/wasm"
	"reflect"
)

type VMContext struct {
	evm    *EVM
	config Config
	db     StateDB
	Input  []byte
	Output []byte
}

func NewHostModule() *wasm.Module {
	m := wasm.NewModule()
	paramTypes := make([]wasm.ValueType, 14)
	for i := 0; i < len(paramTypes); i++ {
		paramTypes[i] = wasm.ValueTypeI32
	}

	m.Types = &wasm.SectionTypes{
		Entries: []wasm.FunctionSig{
			{
				Form: 0,
			},
			{
				Form:        0,
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:       0,
				ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			{
				Form:       0,
				ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
			},
		},
	}
	m.FunctionIndexSpace = []wasm.Function{
		//{
		//	Sig:  &m.Types.Entries[0],
		//	Host: reflect.ValueOf(Abort),
		//	Body: &wasm.FunctionBody{}, // create a dummy wasm body (the actual value will be taken from Host.)
		//},
		//{
		//	Sig:  &m.Types.Entries[2],
		//	Host: reflect.ValueOf(Debug),
		//	Body: &wasm.FunctionBody{}, // create a dummy wasm body (the actual value will be taken from Host.)
		//},
		{
			Sig:  &m.Types.Entries[1],
			Host: reflect.ValueOf(InputLength),
			Body: &wasm.FunctionBody{}, // create a dummy wasm body (the actual value will be taken from Host.)
		},
		{
			Sig:  &m.Types.Entries[2],
			Host: reflect.ValueOf(GetInput),
			Body: &wasm.FunctionBody{}, // create a dummy wasm body (the actual value will be taken from Host.)
		},
		//{
		//	Sig:  &m.Types.Entries[3],
		//	Host: reflect.ValueOf(Print),
		//	Body: &wasm.FunctionBody{}, // create a dummy wasm body (the actual value will be taken from Host.)
		//},
	}
	m.Export = &wasm.SectionExports{
		Entries: map[string]wasm.ExportEntry{
			//"abort": {
			//	FieldStr: "abort",
			//	Kind:     wasm.ExternalFunction,
			//	Index:    0,
			//},
			//"platon_debug": {
			//	FieldStr: "platon_debug",
			//	Kind:     wasm.ExternalFunction,
			//	Index:    1,
			//},
			"platon_input_length": {
				FieldStr: "platon_input_length",
				Kind:     wasm.ExternalFunction,
				Index:    2,
			},
			"platon_get_input": {
				FieldStr: "platon_get_input",
				Kind:     wasm.ExternalFunction,
				Index:    3,
			},
			//"prints_l": {
			//	FieldStr: "prints_l",
			//	Kind:     wasm.ExternalFunction,
			//	Index:    4,
			//},
		},
	}

	return m
}

//func Abort(proc *exec.Process) uint64 {
//	return 0
//}
//
//func Debug(proc *exec.Process, dst uint32) uint64 {
//	fmt.Print("debug:", dst)
//	return 0
//}
//
//func Print(proc *exec.Process, dst uint32, dst2 uint32) uint64 {
//	return 0
//}

func InputLength(proc *exec.Process) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	return uint32(len(ctx.Input))
}

func GetInput(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	_, err := proc.WriteAt(ctx.Input, int64(dst))
	if err != nil {
		panic(err)
	}
}
