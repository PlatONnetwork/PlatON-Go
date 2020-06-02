package vm

import (
	"crypto/sha256"

	"golang.org/x/crypto/ripemd160"

	"github.com/PlatONnetwork/PlatON-Go/common"
	imath "github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/wagon/exec"
	"github.com/PlatONnetwork/wagon/wasm"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/params"

	"math/big"
	"reflect"
)

type VMContext struct {
	evm      *EVM
	contract *Contract
	config   Config
	gasTable params.GasTable
	db       StateDB
	Input    []byte
	CallOut  []byte
	Output   []byte
	readOnly bool // Whether to throw on stateful modifications
	Revert   bool
	Log      *WasmLogger
}

func addFuncExport(m *wasm.Module, sig wasm.FunctionSig, function wasm.Function, export wasm.ExportEntry) {
	typesLen := len(m.Types.Entries)
	m.Types.Entries = append(m.Types.Entries, sig)
	function.Sig = &m.Types.Entries[typesLen]
	funcLen := len(m.FunctionIndexSpace)
	m.FunctionIndexSpace = append(m.FunctionIndexSpace, function)
	export.Index = uint32(funcLen)
	m.Export.Entries[export.FieldStr] = export
}
func NewHostModule() *wasm.Module {
	m := wasm.NewModule()
	m.Export.Entries = make(map[string]wasm.ExportEntry)

	// void platon_gas_price(uint8_t gas_price)
	// func $platon_gas_price(param $0 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GasPrice),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_gas_price",
			Kind:     wasm.ExternalFunction,
		},
	)
	// platon_block_hash(int64_t num,  uint8_t hash[32])
	// func $platon_block_hash(param $0 i64) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI64, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(BlockHash),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_block_hash",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint64_t platon_block_number()
	// func $platon_block_number (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(BlockNumber),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_block_number",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint64_t platon_gas_limit()
	// func $platon_gas_limit (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(GasLimit),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_gas_limit",
			Kind:     wasm.ExternalFunction,
		},
	)
	// uint64_t platon_gas()
	// func $platon_gas (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(Gas),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_gas",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int64_t platon_timestamp()
	// func $timestamp (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(Timestamp),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_timestamp",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_coinbase(uint8_t addr[20])
	// func $platon_coinbase (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Coinbase),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_coinbase",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint8_t platon_balance(const uint8_t addr[20], uint8_t balance[32])
	// func $platon_balance (param $0 i32) (param $1 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Balance),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_balance",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_origin(uint8_t addr[20])
	// func $platon_origin (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Origin),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_origin",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_caller(uint8_t addr[20])
	// func $platon_caller (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Caller),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_caller",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint8_t platon_call_value(uint8_t val[32]);
	// func $platon_call_value (param $0 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(CallValue),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_call_value",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_address(uint8_t addr[20])
	// func $platon_address  (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Address),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_address",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_sha3(const uint8_t *src, size_t srcLen, uint8_t *dest, size_t destLen)
	// func $platon_sha3  (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Sha3),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_sha3",
			Kind:     wasm.ExternalFunction,
		},
	)

	// uint64_t platon_caller_nonce()
	// func $platon_caller_nonce  (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
		},
		wasm.Function{
			Host: reflect.ValueOf(CallerNonce),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_caller_nonce",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_transfer(const uint8_t to[20], const uint8_t *amount, size_t len)
	// func $platon_transfer  (param $1 i32) (param $2 i32) (param $3 i32) (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Transfer),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_transfer",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_set_state(const uint8_t* key, size_t klen, const uint8_t *value, size_t vlen)
	// func $platon_set_state (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(SetState),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_set_state",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t platon_get_state_length (const uint8_t* key, size_t klen)
	// func $platon_get_state_length (param $0 i32) (param $1 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{

			Host: reflect.ValueOf(GetStateLength),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_state_length",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t platon_get_state(const uint8_t* key, size_t klen, uint8_t *value, size_t vlen)
	// func $platon_get_state (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetState),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_state",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t platon_get_input_length()
	// func $platon_get_input_length  (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetInputLength),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_input_length",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_get_input(const uint8_t *value)
	// func $platon_get_input (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetInput),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_input",
			Kind:     wasm.ExternalFunction,
		},
	)

	// size_t platon_get_call_output_length()
	// func $platon_get_call_output_length  (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetCallOutputLength),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_call_output_length",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_get_call_output(const uint8_t *value)
	// func $platon_get_call_output (param $0 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(GetCallOutput),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_get_call_output",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_return(const uint8_t *value, size_t len)
	// func $platon_return(param $0 i32) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(ReturnContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_return",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_revert()
	// func $platon_return()
	addFuncExport(m,
		wasm.FunctionSig{},
		wasm.Function{
			Host: reflect.ValueOf(Revert),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_revert",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_panic()
	// func $platon_panic()
	addFuncExport(m,
		wasm.FunctionSig{},
		wasm.Function{
			Host: reflect.ValueOf(Panic),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_panic",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_debug(const uint8_t *dst, size_t len)
	// func $platon_debug (param i32 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Debug),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_debug",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_call(const uint8_t to[20], const uint8_t *args, size_t argsLen, const uint8_t *value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
	// func $platon_call  (param $0 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(CallContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_call",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_delegate_call(const uint8_t to[20], const uint8_t* args, size_t argsLen, const uint8_t* callCost, size_t callCostLen);
	// func $platon_delegate_call (param $0 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(DelegateCallContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_delegate_call",
			Kind:     wasm.ExternalFunction,
		},
	)

	/*	// int32_t platon_static_call(const uint8_t to[20], const uint8_t* args, size_t argsLen, const uint8_t* callCost, size_t callCostLen);
		// func $platon_static_call (param $0 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (result i32)
		addFuncExport(m,
			wasm.FunctionSig{
				ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
				ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
			},
			wasm.Function{
				Host: reflect.ValueOf(StaticCallContract),
				Body: &wasm.FunctionBody{},
			},
			wasm.ExportEntry{
				FieldStr: "platon_static_call",
				Kind:     wasm.ExternalFunction,
			},
		)*/

	// int32_t platon_destroy(const uint8_t to[20])
	// func $platon_destroy (param $0 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(DestroyContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_destroy",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_migrate(uint8_t newAddr[20], const uint8_t* args, size_t argsLen, const uint8_t* value, size_t valueLen, const uint8_t* callCost, size_t callCostLen);
	// func $platon_migrate  (param $1 i32) (param $2 i32) (param $0 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32,
				wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(MigrateContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_migrate",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_event(const uint8_t* indexes, size_t indexesLen, const uint8_t* args, size_t argsLen)
	// func $platon_event (param $0 i32) (param $1 i32) (param $0 i32) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(EmitEvent),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_event",
			Kind:     wasm.ExternalFunction,
		},
	)

	// int32_t platon_ecrecover(const uint8_t hash[32], const uint8_t* sig, const uint8_t sig_len, uint8_t addr[20])
	// func platon_ecrecover (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Ecrecover),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_ecrecover",
			Kind:     wasm.ExternalFunction,
		},
	)
	// void platon_ripemd160(const uint8_t *input, uint32_t input_len, uint8_t addr[20])
	// func platon_ripemd160 (param $0 i32) (param $1 i32) (param $2 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Ripemd160),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_ripemd160",
			Kind:     wasm.ExternalFunction,
		},
	)

	// void platon_sha256(const uint8_t *input, uint32_t input_len, uint8_t hash[32])
	// func platon_sha256 (param $0 i32) (param $1 i32) (param $2 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(Sha256),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_sha256",
			Kind:     wasm.ExternalFunction,
		},
	)
	return m
}

func checkGas(ctx *VMContext, gas uint64) {
	if !ctx.contract.UseGas(gas) {
		panic(ErrOutOfGas)
	}
}
func GasPrice(proc *exec.Process, gasPrice uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	value := ctx.evm.GasPrice.Bytes()
	_, err := proc.WriteAt(value, int64(gasPrice))
	if err != nil {
		panic(err)
	}

	return uint32(len(value))
}

func BlockHash(proc *exec.Process, num uint64, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasExtStep)
	blockHash := ctx.evm.GetHash(num)
	_, err := proc.WriteAt(blockHash.Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

func BlockNumber(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return ctx.evm.BlockNumber.Uint64()
}

func GasLimit(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return ctx.evm.GasLimit
}

func Gas(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return ctx.contract.Gas
}

func Timestamp(proc *exec.Process) int64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return ctx.evm.Time.Int64()
}

func Coinbase(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	coinBase := ctx.evm.Coinbase
	_, err := proc.WriteAt(coinBase.Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

func Balance(proc *exec.Process, dst uint32, balance uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.Balance)
	var addr common.Address
	_, err := proc.ReadAt(addr[:], int64(dst))
	if nil != err {
		panic(err)
	}
	value := ctx.evm.StateDB.GetBalance(addr).Bytes()
	_, err = proc.WriteAt(value, int64(balance))
	if nil != err {
		panic(err)
	}
	return uint32(len(value))
}

func Origin(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.evm.Origin.Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

func Caller(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.contract.caller.Address().Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

// define: uint8_t callValue();
func CallValue(proc *exec.Process, dst uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	value := ctx.contract.value.Bytes()
	_, err := proc.WriteAt(value, int64(dst))
	if nil != err {
		panic(err)
	}
	return uint32(len(value))
}

// define: void address(char hash[20]);
func Address(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.contract.Address().Bytes(), int64(dst))
	if nil != err {
		panic(err)
	}
}

// define: void sha3(char *src, size_t srcLen, char *dest, size_t destLen);
func Sha3(proc *exec.Process, src uint32, srcLen uint32, dst uint32, dstLen uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		wordGas  uint64
		overflow bool
	)

	if wordGas, overflow = imath.SafeMul(toWordSize(uint64(srcLen)), params.Sha3WordGas); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(wordGas, params.Sha3Gas); overflow {
		panic(errGasUintOverflow)
	}

	checkGas(ctx, gas)

	data := make([]byte, srcLen)
	_, err := proc.ReadAt(data, int64(src))
	if nil != err {
		panic(err)
	}
	hash := crypto.Keccak256(data)
	if int(dstLen) < len(hash) {
		panic(ErrWASMSha3DstToShort)
	}
	_, err = proc.WriteAt(hash, int64(dst))
	if nil != err {
		panic(err)
	}
}

func CallerNonce(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	addr := ctx.contract.Caller()
	return ctx.evm.StateDB.GetNonce(addr)
}

func Transfer(proc *exec.Process, dst uint32, amount uint32, len uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)
	address := make([]byte, common.AddressLength)

	_, err := proc.ReadAt(address, int64(dst))
	if nil != err {
		panic(err)
	}

	value := make([]byte, len)
	_, err = proc.ReadAt(value, int64(amount))
	if nil != err {
		panic(err)
	}
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(value)
	bValue = imath.U256(bValue)
	addr := common.BytesToAddress(address)

	transfersValue := bValue.Sign() != 0
	gas := ctx.gasTable.Calls
	if transfersValue {
		gas += params.CallValueTransferGas
	}
	gasTemp, err := callGasWasm(ctx.contract.Gas, params.TxGas, new(big.Int).SetUint64(ctx.contract.Gas))
	if nil != err {
		panic(err)
	}
	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(gas, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	gas = ctx.evm.callGasTemp
	if transfersValue {
		if gas, overflow = imath.SafeAdd(gas, params.CallStipend); overflow {
			panic(errGasUintOverflow)
		}
	}

	_, returnGas, err := ctx.evm.Call(ctx.contract, addr, nil, gas, bValue)

	var status int32

	if err != nil {
		status = -1
	} else {
		status = 0
	}
	ctx.contract.Gas += returnGas

	return status
}

// storage external function

func SetState(proc *exec.Process, key uint32, keyLen uint32, val uint32, valLen uint32) {
	ctx := proc.HostCtx().(*VMContext)
	if ctx.readOnly {
		panic(ErrWASMWriteProtection)
	}

	switch {
	case valLen == 0:
		checkGas(ctx, params.SstoreClearGas)
	default:
		checkGas(ctx, (toWordSize(uint64(keyLen)+(uint64(valLen)))/32)*params.SstoreSetGas)
	}

	keyBuf := make([]byte, keyLen)
	_, err := proc.ReadAt(keyBuf, int64(key))
	if nil != err {
		panic(err)
	}
	valBuf := make([]byte, valLen)
	_, err = proc.ReadAt(valBuf, int64(val))
	if nil != err {
		panic(err)
	}
	ctx.evm.StateDB.SetState(ctx.contract.Address(), keyBuf, valBuf)
}

func GetStateLength(proc *exec.Process, key uint32, keyLen uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	keyBuf := make([]byte, keyLen)
	_, err := proc.ReadAt(keyBuf, int64(key))
	if nil != err {
		panic(err)
	}
	val := ctx.evm.StateDB.GetState(ctx.contract.Address(), keyBuf)

	checkGas(ctx, ctx.gasTable.SLoad)

	return uint32(len(val))
}

func GetState(proc *exec.Process, key uint32, keyLen uint32, val uint32, valLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, ctx.gasTable.SLoad)

	keyBuf := make([]byte, keyLen)
	_, err := proc.ReadAt(keyBuf, int64(key))
	if nil != err {
		panic(err)
	}
	valBuf := ctx.evm.StateDB.GetState(ctx.contract.Address(), keyBuf)
	vlen := len(valBuf)
	if uint32(vlen) > valLen {
		return -1
	}

	_, err = proc.WriteAt(valBuf, int64(val))
	if nil != err {
		panic(err)
	}
	return int32(vlen)
}

func GetInputLength(proc *exec.Process) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return uint32(len(ctx.Input))
}

func GetInput(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.Input, int64(dst))
	if err != nil {
		panic(err)
	}
}

func GetCallOutputLength(proc *exec.Process) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	return uint32(len(ctx.CallOut))
}

func GetCallOutput(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	checkGas(ctx, GasQuickStep)
	_, err := proc.WriteAt(ctx.CallOut, int64(dst))
	if err != nil {
		panic(err)
	}
}

func ReturnContract(proc *exec.Process, dst uint32, len uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)
	if gas, overflow = imath.SafeAdd(params.MemoryGas, uint64(len)); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, GasQuickStep); overflow {
		panic(errGasUintOverflow)
	}

	checkGas(ctx, gas)
	ctx.Output = make([]byte, len)
	_, err := proc.ReadAt(ctx.Output, int64(dst))
	if err != nil {
		panic(err)
	}
}

func Revert(proc *exec.Process) {
	ctx := proc.HostCtx().(*VMContext)
	ctx.Revert = true
	proc.Terminate()
}

func Panic(proc *exec.Process) {
	panic(ErrWASMPanicOp)
}

func Debug(proc *exec.Process, dst uint32, len uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)

	if gas, overflow = imath.SafeAdd(params.MemoryGas, toWordSize(uint64(len))); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, GasSlowStep); overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	buf := make([]byte, len)
	_, err := proc.ReadAt(buf, int64(dst))
	if nil != err {
		panic(err)
	}
	ctx.Log.Debug("WASM:" + string(buf) + "\n")
	ctx.Log.Flush()
}

func CallContract(proc *exec.Process, addrPtr, args, argsLen, val, valLen, callCost, callCostLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	addr := common.BytesToAddress(address)

	input := make([]byte, argsLen)
	_, err = proc.ReadAt(input, int64(args))
	if nil != err {
		panic(err)
	}

	value := make([]byte, valLen)
	_, err = proc.ReadAt(value, int64(val))
	if nil != err {
		panic(err)
	}
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(value)
	bValue = imath.U256(bValue)

	cost := make([]byte, callCostLen)
	_, err = proc.ReadAt(cost, int64(callCost))
	if nil != err {
		panic(err)
	}
	bCost := new(big.Int)
	// 256 bits
	bCost.SetBytes(cost)
	bCost = imath.U256(bCost)
	if bCost.Cmp(common.Big0) == 0 {
		bCost = new(big.Int).SetUint64(ctx.contract.Gas)
	}

	gas := ctx.gasTable.Calls
	transfersValue := bValue.Sign() != 0
	if transfersValue && ctx.evm.StateDB.Empty(addr) {
		gas += params.CallNewAccountGas
	}

	if transfersValue {
		gas += params.CallValueTransferGas
	}

	gasTemp, err := callGasWasm(ctx.contract.Gas, gas, bCost)
	if nil != err {
		panic(err)
	}
	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(gas, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	gas = ctx.evm.callGasTemp
	if bValue.Sign() != 0 {
		if gas, overflow = imath.SafeAdd(gas, params.CallStipend); overflow {
			panic(errGasUintOverflow)
		}
	}

	ret, returnGas, err := ctx.evm.Call(ctx.contract, addr, input, gas, bValue)

	var status int32

	if err != nil {
		status = -1
	} else {
		status = 0
	}
	if err == nil || err == errExecutionReverted {
		ctx.CallOut = ret
	}
	ctx.contract.Gas += returnGas

	return status
}

func DelegateCallContract(proc *exec.Process, addrPtr, params, paramsLen, callCost, callCostLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	addr := common.BytesToAddress(address)

	input := make([]byte, paramsLen)
	_, err = proc.ReadAt(input, int64(params))
	if nil != err {
		panic(err)
	}

	cost := make([]byte, callCostLen)
	_, err = proc.ReadAt(cost, int64(callCost))
	if nil != err {
		panic(err)
	}
	bCost := new(big.Int)
	// 256 bits
	bCost.SetBytes(cost)
	bCost = imath.U256(bCost)
	if bCost.Cmp(common.Big0) == 0 {
		bCost = new(big.Int).SetUint64(ctx.contract.Gas)
	}

	gasTemp, err := callGasWasm(ctx.contract.Gas, ctx.gasTable.Calls, bCost)
	if nil != err {
		panic(err)
	}
	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(ctx.gasTable.Calls, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	gas = ctx.evm.callGasTemp

	ret, returnGas, err := ctx.evm.DelegateCall(ctx.contract, addr, input, gas)

	var status int32

	if err != nil {
		status = -1
	} else {
		status = 0
	}
	if err == nil || err == errExecutionReverted {
		ctx.CallOut = ret
	}
	ctx.contract.Gas += returnGas

	return status
}

func StaticCallContract(proc *exec.Process, addrPtr, params, paramsLen, callCost, callCostLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	addr := common.BytesToAddress(address)

	input := make([]byte, paramsLen)
	_, err = proc.ReadAt(input, int64(params))
	if nil != err {
		panic(err)
	}

	cost := make([]byte, callCostLen)
	_, err = proc.ReadAt(cost, int64(callCost))
	if nil != err {
		panic(err)
	}
	bCost := new(big.Int)
	// 256 bits
	bCost.SetBytes(cost)
	bCost = imath.U256(bCost)
	if bCost.Cmp(common.Big0) == 0 {
		bCost = new(big.Int).SetUint64(ctx.contract.Gas)
	}

	gasTemp, err := callGasWasm(ctx.contract.Gas, ctx.gasTable.Calls, bCost)
	if nil != err {
		panic(err)
	}

	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(ctx.gasTable.Calls, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	gas = ctx.evm.callGasTemp

	ret, returnGas, err := ctx.evm.StaticCall(ctx.contract, addr, input, gas)

	var status int32

	if err != nil {
		status = -1
	} else {
		status = 0
	}
	if err == nil || err == errExecutionReverted {
		ctx.CallOut = ret
	}
	ctx.contract.Gas += returnGas

	return status
}

func DestroyContract(proc *exec.Process, addrPtr uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(ErrWASMWriteProtection)
	}

	address := make([]byte, common.AddressLength)
	_, err := proc.ReadAt(address, int64(addrPtr))
	if nil != err {
		panic(err)
	}
	addr := common.BytesToAddress(address)

	contractAddr := ctx.contract.Address()

	gas := ctx.gasTable.Suicide
	if ctx.evm.StateDB.Empty(addr) && ctx.evm.StateDB.GetBalance(contractAddr).Sign() != 0 {
		gas += ctx.gasTable.CreateBySuicide
	}

	if !ctx.evm.StateDB.HasSuicided(ctx.contract.Address()) {
		ctx.evm.StateDB.AddRefund(params.SuicideRefundGas)
	}
	checkGas(ctx, gas)

	balance := ctx.evm.StateDB.GetBalance(contractAddr)

	ctx.evm.StateDB.AddBalance(addr, balance)

	ctx.evm.StateDB.Suicide(contractAddr)

	return 0
}

func MigrateContract(proc *exec.Process, newAddr, args, argsLen, val, valLen, callCost, callCostLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(ErrWASMWriteProtection)
	}

	// check call depth
	if ctx.evm.depth > int(params.CallCreateDepth) {
		panic(ErrDepth)
	}

	oldContract := ctx.contract.Address()

	input := make([]byte, argsLen)
	_, err := proc.ReadAt(input, int64(args))
	if nil != err {
		panic(err)
	}

	if len(input) == 0 {
		panic(ErrWASMMigrate)
	}

	value := make([]byte, valLen)
	_, err = proc.ReadAt(value, int64(val))
	if nil != err {
		panic(err)
	}
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(value)
	bValue = imath.U256(bValue)

	cost := make([]byte, callCostLen)
	_, err = proc.ReadAt(cost, int64(callCost))
	if nil != err {
		panic(err)
	}
	bCost := new(big.Int)
	// 256 bits
	bCost.SetBytes(cost)
	bCost = imath.U256(bCost)
	if bCost.Cmp(common.Big0) == 0 {
		bCost = new(big.Int).SetUint64(ctx.contract.Gas)
	}

	gas := MigrateContractGas
	if bValue.Sign() != 0 {
		gas += params.CallNewAccountGas
	}
	gasTemp, err := callGasWasm(ctx.contract.Gas, gas, bCost)
	if nil != err {
		panic(err)
	}

	ctx.evm.callGasTemp = gasTemp
	gas, overflow := imath.SafeAdd(gas, ctx.evm.callGasTemp)
	if overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)
	gas = ctx.evm.callGasTemp

	sender := ctx.contract.CallerAddress

	// check code of old contract
	oldCode := ctx.evm.StateDB.GetCode(oldContract)
	if len(oldCode) == 0 {
		panic(ErrWASMOldContractCodeNotExists)
	}

	// check balance of sender
	if !ctx.evm.CanTransfer(ctx.evm.StateDB, sender, bValue) {
		panic(ErrInsufficientBalance)
	}

	senderNonce := ctx.evm.StateDB.GetNonce(sender)

	// create new contract address
	newContract := crypto.CreateAddress(sender, senderNonce)
	ctx.evm.StateDB.SetNonce(sender, senderNonce+1)

	// Ensure there's no existing contract already at the designated address
	contractHash := ctx.evm.StateDB.GetCodeHash(newContract)
	if ctx.evm.StateDB.GetNonce(newContract) != 0 || (contractHash != (common.Hash{}) && contractHash != emptyCodeHash) {
		panic(ErrContractAddressCollision)
	}

	// Create a new account on the state
	snapshotForSnapshotDB, snapshotForStateDB := ctx.evm.DBSnapshot()
	ctx.evm.StateDB.CreateAccount(newContract)
	ctx.evm.StateDB.SetNonce(newContract, 1)

	oldBalance := new(big.Int).Set(ctx.evm.StateDB.GetBalance(oldContract))

	// migrate balance from old contract to new contract
	ctx.evm.Transfer(ctx.evm.StateDB, oldContract, newContract, oldBalance)
	// transfer balance from sender to new contract
	ctx.evm.Transfer(ctx.evm.StateDB, sender, newContract, bValue)

	// migrate stateObject storage from old contract to new contract
	ctx.evm.StateDB.MigrateStorage(oldContract, newContract)

	// suicided the old contract
	ctx.evm.StateDB.Suicide(oldContract)

	balance := new(big.Int).Add(bValue, oldBalance)

	// init new contract context
	contract := NewContract(AccountRef(sender), AccountRef(newContract), balance, gas)
	contract.SetCallCode(&newContract, crypto.Keccak256Hash(input), input)
	contract.DeployContract = true

	// deploy new contract
	ret, err := run(ctx.evm, contract, nil, false)

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := len(ret) > params.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateWasmDataGas
		if contract.UseGas(createDataGas) {
			ctx.evm.StateDB.SetCode(newContract, ret)
		} else {
			err = ErrCodeStoreOutOfGas
		}
	}

	// When an error was returned by the VM or when setting the creation code
	// above we revert to the snapshot and consume any gas remaining. Additionally
	// when we're in homestead this also counts for code storage gas errors.
	if maxCodeSizeExceeded || (err != nil && err != ErrCodeStoreOutOfGas) {
		ctx.evm.RevertToDBSnapshot(snapshotForSnapshotDB, snapshotForStateDB)
		if err != errExecutionReverted {
			contract.UseGas(contract.Gas)
		}
	}

	// Assign err if contract code size exceeds the max while the err is still empty.
	if maxCodeSizeExceeded && err == nil {
		err = errMaxCodeSizeExceeded
	}

	if nil != err {
		panic(err)
	}

	ctx.contract.Gas += contract.Gas

	_, err = proc.WriteAt(newContract.Bytes(), int64(newAddr))
	if nil != err {
		panic(err)
	}

	return 0
}

func EmitEvent(proc *exec.Process, indexesPtr, indexesLen, args, argsLen uint32) {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(ErrWASMWriteProtection)
	}

	topics := make([]common.Hash, 0)

	if indexesLen != 0 {

		indexes := make([]byte, indexesLen)
		_, err := proc.ReadAt(indexes, int64(indexesPtr))
		if nil != err {
			panic(err)
		}

		content, _, err := rlp.SplitList(indexes)
		if nil != err {
			panic(err)
		}

		topicCount, err := rlp.CountValues(content)
		if nil != err {
			panic(err)
		}
		if topicCount > WasmTopicNum {
			panic(ErrWASMEventCountToLarge)
		}

		decodeTopics := func(b []byte) ([]byte, []byte, error) {
			member, rest, err := rlp.SplitString(b)
			if nil != err {
				panic(err)
			}
			return member, rest, nil
		}

		for len(content) > 0 {
			mem, tail, err := decodeTopics(content)
			if nil != err {
				panic(err)
			}

			if len(mem) > common.HashLength {
				panic(ErrWASMEventContentToLong)
			}

			topics = append(topics, common.BytesToHash(mem))
			content = tail
		}

	}

	input := make([]byte, argsLen)
	_, err := proc.ReadAt(input, int64(args))
	if nil != err {
		panic(err)
	}

	gas, err := logGas(uint64(len(topics)), uint64(argsLen))
	if nil != err {
		panic(err)
	}
	checkGas(ctx, gas)

	bn := ctx.evm.BlockNumber.Uint64()

	addLog(ctx.evm.StateDB, ctx.contract.Address(), topics, input, bn)
}

func Ecrecover(proc *exec.Process, hashPtr, sigPtr, sigLen, addrPtr uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	checkGas(ctx, params.EcrecoverGas)
	hash := make([]byte, 32)
	_, err := proc.ReadAt(hash, int64(hashPtr))
	if err != nil {
		panic(err)
	}

	sig := make([]byte, sigLen)
	_, err = proc.ReadAt(sig, int64(sigPtr))
	if err != nil {
		panic(err)
	}

	pubKey, err := crypto.Ecrecover(hash, sig)
	if err != nil {
		return -1
	}

	if _, err = proc.WriteAt(crypto.Keccak256(pubKey[1:])[12:], int64(addrPtr)); err != nil {
		return -1
	}
	return 0
}

func Ripemd160(proc *exec.Process, inputPtr, inputLen uint32, outputPtr uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)
	if gas, overflow = imath.SafeMul(toWordSize(uint64(inputLen)), params.Ripemd160PerWordGas); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, params.Ripemd160BaseGas); overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	input := make([]byte, inputLen)
	_, err := proc.ReadAt(input, int64(inputPtr))
	if err != nil {
		panic(err)
	}
	ripemd := ripemd160.New()
	ripemd.Write(input)
	output := ripemd.Sum(nil)
	proc.WriteAt(output, int64(outputPtr))
}

func Sha256(proc *exec.Process, inputPtr, inputLen uint32, outputPtr uint32) {
	ctx := proc.HostCtx().(*VMContext)
	var (
		gas      uint64
		overflow bool
	)

	if gas, overflow = imath.SafeMul(toWordSize(uint64(inputLen)), params.Sha256PerWordGas); overflow {
		panic(errGasUintOverflow)
	}
	if gas, overflow = imath.SafeAdd(gas, params.Sha256BaseGas); overflow {
		panic(errGasUintOverflow)
	}
	checkGas(ctx, gas)

	input := make([]byte, inputLen)
	_, err := proc.ReadAt(input, int64(inputPtr))
	if err != nil {
		panic(err)
	}
	h := sha256.Sum256(input)

	proc.WriteAt(h[:], int64(outputPtr))
}

func addLog(state StateDB, address common.Address, topics []common.Hash, data []byte, bn uint64) {
	log := &types.Log{
		Address:     address,
		Topics:      topics,
		Data:        data,
		BlockNumber: bn,
	}
	state.AddLog(log)
}

func logGas(topicNum, dataSize uint64) (uint64, error) {
	gas := params.LogGas
	var overflow bool
	if gas, overflow = imath.SafeAdd(gas, topicNum*params.LogTopicGas); overflow {
		return 0, errGasUintOverflow
	}

	var logSizeGas uint64
	if logSizeGas, overflow = imath.SafeMul(dataSize, params.LogDataGas); overflow {
		return 0, errGasUintOverflow
	}
	if gas, overflow = imath.SafeAdd(gas, logSizeGas); overflow {
		return 0, errGasUintOverflow
	}
	return gas, nil
}
