package vm

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	inner "github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/wagon/exec"
	"github.com/PlatONnetwork/wagon/wasm"

	"math"
	"math/big"
	"reflect"
)

type VMContext struct {
	evm      *EVM
	contract *Contract
	config   Config
	db       StateDB
	Input    []byte
	CallOut  []byte
	Output   []byte
	readOnly bool // Whether to throw on stateful modifications
	Revert   bool
	log      *WasmLogger
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

	// uint64_t platon_gas_price()
	// func $platon_gas_price(result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI64},
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

	// int64_t timestamp()
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

	// void platon_coinbase(uint8_t hash[20])
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

	// void platon_balance(uint8_t addr[32], uin8_t balance[32])
	// func $platon_balance (param $0 i32) (param $1 i32)
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

	// void platon_origin(uint8_t hash[20])
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

	// void platon_caller(uint8_t hash[20])
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

	// void platon_call_value(uint8_t val[32])
	// func $platon_call_value (param $0 i32)
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

	// void platon_address(uint8_t hash[20])
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

	// int32_t platon_transfer(const uint8_t* to, size_t toLen, uint8_t *amount, size_t len)
	// func $platon_transfer  (param $0 i32) (param $1 i32) (param $2 i32) (param $3 i32) (result i64)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
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

	// void platon_debug(uint8_t *dst, size_t len)
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

	// todo
	// int32_t platon_call(const uint8_t* to, const uint8_t* args, size_t argsLen, const uint8_t* amount, size_t amountLen)
	// func $platon_call  (param $0 i32) (param $1 i32) (param $2 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
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

	// todo
	// int32_t platon_delegatecall(const uint8_t* to, const uint8_t* args, size_t argsLen)
	// func $platon_delegatecall (param $0 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(DelegateCallContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_delegatecall",
			Kind:     wasm.ExternalFunction,
		},
	)

	// todo
	// int32_t platon_staticcall(const uint8_t* to, const uint8_t* args, size_t argsLen)
	// func $platon_staticcall (param $0 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
			ReturnTypes: []wasm.ValueType{wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(StaticCallContract),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_staticcall",
			Kind:     wasm.ExternalFunction,
		},
	)

	// todo
	// int32_t platon_destroy()
	// func $platon_destroy (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
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

	// todo
	// int32_t platon_migrate(const uint8_t* old, uint8_t* new, const uint8_t* args, size_t argsLen, const uint8_t* value, size_t valueLen)
	// func $platon_migrate (param $0 i32) (param $1 i32) (param $2 i32) (param $0 i32) (param $1 i32) (param $2 i32) (result i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes:  []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
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

	// todo
	// void platon_event(const uint8_t* args, size_t argsLen)
	// func $platon_event (param $0 i32) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32},
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

	// todo
	// void platon_event1(const uint8_t* topic, size_t topicLen, const uint8_t* args, size_t argsLen)
	// func $platon_event1 (param $0 i32) (param $1 i32) (param $0 i32) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(EmitEvent1),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_event1",
			Kind:     wasm.ExternalFunction,
		},
	)

	// todo
	// void platon_event2(const uint8_t* topic1, size_t topic1Len, const uint8_t* topic2, size_t topic2Len, const uint8_t* args, size_t argsLen)
	// func $platon_event2 (param $0 i32) (param $1 i32) (param $0 i32) (param $1 i32) (param $0 i32) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(EmitEvent2),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_event2",
			Kind:     wasm.ExternalFunction,
		},
	)

	// todo
	// void platon_event3(const uint8_t* topic1, size_t topic1Len, const uint8_t* topic2, size_t topic2Len, const uint8_t* topic3, size_t topic3Len, uint8_t* args, size_t argsLen)
	// func $platon_event3 (param $0 i32) (param $1 i32) (param $0 i32) (param $1 i32) (param $0 i32) (param $1 i32) (param $0 i32) (param $1 i32)
	addFuncExport(m,
		wasm.FunctionSig{
			ParamTypes: []wasm.ValueType{wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32, wasm.ValueTypeI32},
		},
		wasm.Function{
			Host: reflect.ValueOf(EmitEvent3),
			Body: &wasm.FunctionBody{},
		},
		wasm.ExportEntry{
			FieldStr: "platon_event3",
			Kind:     wasm.ExternalFunction,
		},
	)

	return m
}

func GasPrice(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	price := ctx.evm.GasPrice.Uint64()
	return price
}

func BlockHash(proc *exec.Process, num uint64, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	blockHash := ctx.evm.GetHash(num)
	proc.WriteAt(blockHash.Bytes(), int64(dst))
}

func BlockNumber(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	return ctx.evm.BlockNumber.Uint64()
}

func GasLimit(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	return ctx.evm.GasLimit
}

func Timestamp(proc *exec.Process) int64 {
	ctx := proc.HostCtx().(*VMContext)
	return ctx.evm.Time.Int64()
}

func Coinbase(proc *exec.Process, dst uint32) int64 {
	ctx := proc.HostCtx().(*VMContext)
	coinBase := ctx.evm.Coinbase
	proc.WriteAt(coinBase.Bytes(), int64(dst))
	return 0
}

func Balance(proc *exec.Process, dst uint32, balance uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	var addr common.Address
	proc.ReadAt(addr[:], int64(dst))
	value := ctx.evm.StateDB.GetBalance(addr).Bytes()
	proc.WriteAt(value, int64(balance))
	return uint32(len(value))
}

func Origin(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	proc.WriteAt(ctx.evm.Origin.Bytes(), int64(dst))
}

func Caller(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	proc.WriteAt(ctx.contract.caller.Address().Bytes(), int64(dst))
}

// define: uint8_t callValue();
func CallValue(proc *exec.Process, dst uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	value := ctx.contract.value.Bytes()
	proc.WriteAt(value, int64(dst))
	return uint32(len(value))
}

// define: void address(char hash[20]);
func Address(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	proc.WriteAt(ctx.contract.Address().Bytes(), int64(dst))
}

// define: void sha3(char *src, size_t srcLen, char *dest, size_t destLen);
func Sha3(proc *exec.Process, src uint32, srcLen uint32, dst uint32, dstLen uint32) {
	data := make([]byte, srcLen)
	proc.ReadAt(data, int64(src))
	hash := crypto.Keccak256(data)
	if int(dst) < len(hash) {
		// todo terminate
	}
	proc.WriteAt(hash, int64(dst))
}

func CallerNonce(proc *exec.Process) uint64 {
	ctx := proc.HostCtx().(*VMContext)
	addr := ctx.contract.Caller()
	return ctx.evm.StateDB.GetNonce(addr)
}

func Transfer(proc *exec.Process, dst uint32, dstLen uint32, amount uint32, len uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)
	address := make([]byte, dstLen)

	proc.ReadAt(address, int64(dst))

	value := make([]byte, len)
	proc.ReadAt(value, int64(amount))
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(value)
	bValue = inner.U256(bValue)
	addr := common.BytesToAddress(address)

	gas := ctx.evm.callGasTemp
	if bValue.Sign() != 0 {
		gas += params.CallStipend
	}
	_, returnGas, err := ctx.evm.Call(ctx.contract, addr, nil, gas, bValue)
	if err != nil {
		return 1
	}
	if !ctx.contract.UseGas(returnGas) {
		return 1
	}
	return 0
}

// storage external function

func SetState(proc *exec.Process, key uint32, keyLen uint32, val uint32, valLen uint32) {
	ctx := proc.HostCtx().(*VMContext)
	if ctx.readOnly {
		panic(errWASMWriteProtection)
	}
	keyBuf := make([]byte, keyLen)
	proc.ReadAt(keyBuf, int64(key))
	valBuf := make([]byte, valLen)
	proc.ReadAt(valBuf, int64(val))
	ctx.evm.StateDB.SetState(ctx.contract.Address(), keyBuf, valBuf)
}

func GetStateLength(proc *exec.Process, key uint32, keyLen uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	keyBuf := make([]byte, keyLen)
	proc.ReadAt(keyBuf, int64(key))
	val := ctx.evm.StateDB.GetState(ctx.contract.Address(), keyBuf)

	return uint32(len(val))
}

func GetState(proc *exec.Process, key uint32, keyLen uint32, val uint32, valLen uint32) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	keyBuf := make([]byte, keyLen)
	proc.ReadAt(keyBuf, int64(key))
	valBuf := ctx.evm.StateDB.GetState(ctx.contract.Address(), keyBuf)

	if uint32(len(valBuf)) > valLen {
		return math.MaxUint32
	}

	proc.WriteAt(valBuf, int64(val))
	return 0
}

func GetInputLength(proc *exec.Process) uint32 {
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

func GetCallOutputLength(proc *exec.Process) uint32 {
	ctx := proc.HostCtx().(*VMContext)
	return uint32(len(ctx.CallOut))
}

func GetCallOutput(proc *exec.Process, dst uint32) {
	ctx := proc.HostCtx().(*VMContext)
	_, err := proc.WriteAt(ctx.CallOut, int64(dst))
	if err != nil {
		panic(err)
	}
}

func ReturnContract(proc *exec.Process, dst uint32, len uint32) {
	ctx := proc.HostCtx().(*VMContext)
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
	panic("transaction panic")
}

func Debug(proc *exec.Process, dst uint32, len uint32) {
	ctx := proc.HostCtx().(*VMContext)
	buf := make([]byte, len)
	proc.ReadAt(buf, int64(dst))
	ctx.log.Debug(string(buf))
}

func CallContract(proc *exec.Process, addrPtr, args, argsLen, val, valLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	address := make([]byte, common.AddressLength)
	proc.ReadAt(address, int64(addrPtr))
	addr := common.BytesToAddress(address)

	input := make([]byte, argsLen)
	proc.ReadAt(input, int64(args))

	value := make([]byte, valLen)
	proc.ReadAt(value, int64(val))
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(value)
	bValue = inner.U256(bValue)

	gas := ctx.evm.callGasTemp
	if bValue.Sign() != 0 {
		gas += params.CallStipend
	}

	ret, _, err := ctx.evm.Call(ctx.contract, addr, input, gas, bValue)
	if err != nil {
		return 0
	}

	ctx.CallOut = ret
	return int32(len(ctx.CallOut))
}

func DelegateCallContract(proc *exec.Process, addrPtr, params, paramsLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	address := make([]byte, common.AddressLength)
	proc.ReadAt(address, int64(addrPtr))
	addr := common.BytesToAddress(address)

	input := make([]byte, paramsLen)
	proc.ReadAt(input, int64(params))

	//fmt.Println("Addr:", addr.String(), "Data:", input)
	ret, _, err := ctx.evm.DelegateCall(ctx.contract, addr, input, ctx.evm.callGasTemp)
	if err != nil {
		return 0
	}

	ctx.CallOut = ret
	return int32(len(ctx.CallOut))
}

func StaticCallContract(proc *exec.Process, addrPtr, params, paramsLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	address := make([]byte, common.AddressLength)
	proc.ReadAt(address, int64(addrPtr))
	addr := common.BytesToAddress(address)

	input := make([]byte, paramsLen)
	proc.ReadAt(input, int64(params))

	//fmt.Println("Addr:", addr.String(), "Data:", input)
	ret, _, err := ctx.evm.StaticCall(ctx.contract, addr, input, ctx.evm.callGasTemp)
	if err != nil {
		return 0
	}

	ctx.CallOut = ret
	return int32(len(ctx.CallOut))
}

func DestroyContract(proc *exec.Process) int32 {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(errWASMWriteProtection)
	}

	balance := ctx.evm.StateDB.GetBalance(ctx.contract.Address())
	//fmt.Println("sender:", ctx.contract.Caller().String(), "to:", ctx.contract.Address().String(), "value:", balance)
	ctx.evm.StateDB.AddBalance(ctx.contract.Caller(), balance)

	ctx.evm.StateDB.Suicide(ctx.contract.Address())
	return 0
}

func MigrateContract(proc *exec.Process, oldAddr, newAddr, args, argsLen, val, valLen uint32) int32 {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(errWASMWriteProtection)
	}

	// check call depth
	if ctx.evm.depth > int(params.CallCreateDepth) {
		panic(ErrDepth)
	}

	addr1 := make([]byte, common.AddressLength)
	proc.ReadAt(addr1, int64(oldAddr))
	oldContract := common.BytesToAddress(addr1)

	input := make([]byte, argsLen)
	proc.ReadAt(input, int64(args))

	value := make([]byte, valLen)
	proc.ReadAt(value, int64(val))
	bValue := new(big.Int)
	// 256 bits
	bValue.SetBytes(value)
	bValue = inner.U256(bValue)

	gas := ctx.contract.Gas

	sender := ctx.contract.CallerAddress

	// check code of old contract
	oldCode := ctx.evm.StateDB.GetCode(oldContract)
	if len(oldCode) == 0 {
		panic("old target contract is illegal, no contract code exists")
	}

	// create new contract address
	newContract := crypto.CreateAddress(sender, ctx.evm.StateDB.GetNonce(sender))

	// check balance of sender
	if !ctx.evm.CanTransfer(ctx.evm.StateDB, sender, bValue) {
		panic(ErrInsufficientBalance)
	}
	// Create a new account on the state
	snapshot := ctx.evm.StateDB.Snapshot()
	ctx.evm.StateDB.CreateAccount(newContract)
	ctx.evm.StateDB.SetNonce(newContract, 1)

	oldBalance := new(big.Int).Set(ctx.evm.StateDB.GetBalance(oldContract))

	//fmt.Println("sender:", sender.String(),
	//	"\noldCOntract:", oldContract.String(),
	//	"\nnewContract:", newContract.String(),
	//	"\ninput:", input, "\nvalue:", bValue,
	//	"\noldBalance:", oldBalance.String())

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

	// deploy new contract
	ret, err := run(ctx.evm, contract, nil, false)
	//var err error = nil

	// check whether the max code size has been exceeded
	maxCodeSizeExceeded := len(ret) > params.MaxCodeSize
	// if the contract creation ran successfully and no errors were returned
	// calculate the gas required to store the code. If the code could not
	// be stored due to not enough gas set an error and let it be handled
	// by the error checking condition below.
	if err == nil && !maxCodeSizeExceeded {
		createDataGas := uint64(len(ret)) * params.CreateDataGas
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
		ctx.evm.StateDB.RevertToSnapshot(snapshot)
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

	proc.WriteAt(newContract.Bytes(), int64(newAddr))

	return 0
}

func EmitEvent(proc *exec.Process, args, argsLen uint32) {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(errWASMWriteProtection)
	}

	topics := make([]common.Hash, 0)

	input := make([]byte, argsLen)
	proc.ReadAt(input, int64(args))

	bn := ctx.evm.BlockNumber.Uint64()

	//fmt.Println("input:", string(input), "blockNUm:", bn)
	addLog(ctx.evm.StateDB, ctx.contract.Address(), topics, input, bn)
}

func EmitEvent1(proc *exec.Process, t, tLen, args, argsLen uint32) {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(errWASMWriteProtection)
	}

	topic := make([]byte, tLen)
	proc.ReadAt(topic, int64(t))
	topics := []common.Hash{common.BytesToHash(crypto.Keccak256(topic))}

	input := make([]byte, argsLen)
	proc.ReadAt(input, int64(args))

	bn := ctx.evm.BlockNumber.Uint64()

	//fmt.Println("topic:", string(topic), "input:", string(input), "blockNUm:", bn)
	addLog(ctx.evm.StateDB, ctx.contract.Address(), topics, input, bn)
}

func EmitEvent2(proc *exec.Process, t1, t1Len, t2, t2Len, args, argsLen uint32) {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(errWASMWriteProtection)
	}

	topic1 := make([]byte, t1Len)
	proc.ReadAt(topic1, int64(t1))
	topic2 := make([]byte, t2Len)
	proc.ReadAt(topic2, int64(t2))

	arr := [][]byte{topic1, topic2}
	topics := make([]common.Hash, len(arr))
	for i, t := range arr {
		topics[i] = common.BytesToHash(crypto.Keccak256(t))
	}

	input := make([]byte, argsLen)
	proc.ReadAt(input, int64(args))

	bn := ctx.evm.BlockNumber.Uint64()

	//fmt.Println("topic1:", string(topic1), "topic2:", string(topic2), "input:", string(input), "blockNUm:", bn)
	addLog(ctx.evm.StateDB, ctx.contract.Address(), topics, input, bn)
}

func EmitEvent3(proc *exec.Process, t1, t1Len, t2, t2Len, t3, t3Len, args, argsLen uint32) {
	ctx := proc.HostCtx().(*VMContext)

	if ctx.readOnly {
		panic(errWASMWriteProtection)
	}

	topic1 := make([]byte, t1Len)
	proc.ReadAt(topic1, int64(t1))
	topic2 := make([]byte, t2Len)
	proc.ReadAt(topic2, int64(t2))

	topic3 := make([]byte, t3Len)
	proc.ReadAt(topic3, int64(t3))

	arr := [][]byte{topic1, topic2, topic3}
	topics := make([]common.Hash, len(arr))
	for i, t := range arr {
		topics[i] = common.BytesToHash(crypto.Keccak256(t))
	}

	input := make([]byte, argsLen)
	proc.ReadAt(input, int64(args))

	bn := ctx.evm.BlockNumber.Uint64()

	//fmt.Println("topic1:", string(topic1), "topic2:", string(topic2), "topic3:", string(topic3), "input:", string(input), "blockNUm:", bn)
	addLog(ctx.evm.StateDB, ctx.contract.Address(), topics, input, bn)
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
