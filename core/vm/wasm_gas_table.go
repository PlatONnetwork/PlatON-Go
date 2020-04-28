package vm

import "math"

const (
	MigrateContractGas = uint64(68000)
)

var WasmGasCostTable [255]uint64

func init() {
	WasmGasCostTable[Unreachable] = 0
	WasmGasCostTable[Nop] = 0
	WasmGasCostTable[Block] = 0
	WasmGasCostTable[Loop] = 0
	WasmGasCostTable[If] = 0
	WasmGasCostTable[Else] = 0
	WasmGasCostTable[End] = 0
	WasmGasCostTable[Br] = 2
	WasmGasCostTable[BrIf] = 3
	WasmGasCostTable[BrTable] = 2
	WasmGasCostTable[Return] = 2
	WasmGasCostTable[Call] = 2
	WasmGasCostTable[CallIndirect] = 3

	WasmGasCostTable[Drop] = 3
	WasmGasCostTable[Select] = 3

	WasmGasCostTable[GetLocal] = 3
	WasmGasCostTable[SetLocal] = 3
	WasmGasCostTable[TeeLocal] = 3
	WasmGasCostTable[GetGlobal] = 3
	WasmGasCostTable[SetGlobal] = 3

	WasmGasCostTable[I32Load] = 3
	WasmGasCostTable[I64Load] = 3
	WasmGasCostTable[F32Load] = 3
	WasmGasCostTable[F64Load] = 3
	WasmGasCostTable[I32Load8s] = 3
	WasmGasCostTable[I32Load8u] = 3
	WasmGasCostTable[I32Load16s] = 3
	WasmGasCostTable[I32Load16u] = 3
	WasmGasCostTable[I64Load8s] = 3
	WasmGasCostTable[I64Load8u] = 3
	WasmGasCostTable[I64Load16s] = 3
	WasmGasCostTable[I64Load16u] = 3
	WasmGasCostTable[I64Load32s] = 3
	WasmGasCostTable[I64Load32u] = 3

	WasmGasCostTable[I32Store] = 3
	WasmGasCostTable[I64Store] = 3
	WasmGasCostTable[F32Store] = 3
	WasmGasCostTable[F64Store] = 3
	WasmGasCostTable[I32Store8] = 3
	WasmGasCostTable[I32Store16] = 3
	WasmGasCostTable[I64Store8] = 3
	WasmGasCostTable[I64Store16] = 3
	WasmGasCostTable[I64Store32] = 3

	WasmGasCostTable[CurrentMemory] = 3
	WasmGasCostTable[GrowMemory] = 1024

	WasmGasCostTable[I32Const] = 0
	WasmGasCostTable[I64Const] = 0
	WasmGasCostTable[F32Const] = 0
	WasmGasCostTable[F64Const] = 0

	WasmGasCostTable[I32Eqz] = 1
	WasmGasCostTable[I32Eq] = 1
	WasmGasCostTable[I32Ne] = 1
	WasmGasCostTable[I32LtS] = 1
	WasmGasCostTable[I32LtU] = 1
	WasmGasCostTable[I32GtS] = 1
	WasmGasCostTable[I32GtU] = 1
	WasmGasCostTable[I32LeS] = 1
	WasmGasCostTable[I32LeU] = 1
	WasmGasCostTable[I32GeS] = 1
	WasmGasCostTable[I32GeU] = 1
	WasmGasCostTable[I64Eqz] = 1
	WasmGasCostTable[I64Eq] = 1
	WasmGasCostTable[I64Ne] = 1
	WasmGasCostTable[I64LtS] = 1
	WasmGasCostTable[I64LtU] = 1
	WasmGasCostTable[I64GtS] = 1
	WasmGasCostTable[I64GtU] = 1
	WasmGasCostTable[I64LeS] = 1
	WasmGasCostTable[I64LeU] = 1
	WasmGasCostTable[I64GeS] = 1
	WasmGasCostTable[I64GeU] = 1
	WasmGasCostTable[F32Eq] = math.MaxInt64
	WasmGasCostTable[F32Ne] = math.MaxInt64
	WasmGasCostTable[F32Lt] = math.MaxInt64
	WasmGasCostTable[F32Gt] = math.MaxInt64
	WasmGasCostTable[F32Le] = math.MaxInt64
	WasmGasCostTable[F32Ge] = math.MaxInt64
	WasmGasCostTable[F64Eq] = math.MaxInt64
	WasmGasCostTable[F64Ne] = math.MaxInt64
	WasmGasCostTable[F64Lt] = math.MaxInt64
	WasmGasCostTable[F64Gt] = math.MaxInt64
	WasmGasCostTable[F64Le] = math.MaxInt64
	WasmGasCostTable[F64Ge] = math.MaxInt64

	WasmGasCostTable[I32Clz] = 105
	WasmGasCostTable[I32Ctz] = 105
	WasmGasCostTable[I32Popcnt] = 1
	WasmGasCostTable[I32Add] = 1
	WasmGasCostTable[I32Sub] = 1
	WasmGasCostTable[I32Mul] = 3
	WasmGasCostTable[I32DivS] = 80
	WasmGasCostTable[I32DivU] = 80
	WasmGasCostTable[I32RemS] = 80
	WasmGasCostTable[I32RemU] = 80
	WasmGasCostTable[I32And] = 1
	WasmGasCostTable[I32Or] = 1
	WasmGasCostTable[I32Xor] = 1
	WasmGasCostTable[I32Shl] = 2
	WasmGasCostTable[I32ShrS] = 2
	WasmGasCostTable[I32ShrU] = 2
	WasmGasCostTable[I32Rotl] = 2
	WasmGasCostTable[I32Rotr] = 2
	WasmGasCostTable[I64Clz] = 105
	WasmGasCostTable[I64Ctz] = 105
	WasmGasCostTable[I64Popcnt] = 1
	WasmGasCostTable[I64Add] = 1
	WasmGasCostTable[I64Sub] = 1
	WasmGasCostTable[I64Mul] = 3
	WasmGasCostTable[I64DivS] = 80
	WasmGasCostTable[I64DivU] = 80
	WasmGasCostTable[I64RemS] = 80
	WasmGasCostTable[I64RemU] = 80
	WasmGasCostTable[I64And] = 1
	WasmGasCostTable[I64Or] = 1
	WasmGasCostTable[I64Xor] = 1
	WasmGasCostTable[I64Shl] = 2
	WasmGasCostTable[I64ShrS] = 2
	WasmGasCostTable[I64ShrU] = 2
	WasmGasCostTable[I64Rotl] = 2
	WasmGasCostTable[I64Rotr] = 2
	WasmGasCostTable[F32Abs] = math.MaxInt64
	WasmGasCostTable[F32Neg] = math.MaxInt64
	WasmGasCostTable[F32Ceil] = math.MaxInt64
	WasmGasCostTable[F32Floor] = math.MaxInt64
	WasmGasCostTable[F32Trunc] = math.MaxInt64
	WasmGasCostTable[F32Nearest] = math.MaxInt64
	WasmGasCostTable[F32Sqrt] = math.MaxInt64
	WasmGasCostTable[F32Add] = math.MaxInt64
	WasmGasCostTable[F32Sub] = math.MaxInt64
	WasmGasCostTable[F32Mul] = math.MaxInt64
	WasmGasCostTable[F32Div] = math.MaxInt64
	WasmGasCostTable[F32Min] = math.MaxInt64
	WasmGasCostTable[F32Max] = math.MaxInt64
	WasmGasCostTable[F32Copysign] = math.MaxInt64
	WasmGasCostTable[F64Abs] = math.MaxInt64
	WasmGasCostTable[F64Neg] = math.MaxInt64
	WasmGasCostTable[F64Ceil] = math.MaxInt64
	WasmGasCostTable[F64Floor] = math.MaxInt64
	WasmGasCostTable[F64Trunc] = math.MaxInt64
	WasmGasCostTable[F64Nearest] = math.MaxInt64
	WasmGasCostTable[F64Sqrt] = math.MaxInt64
	WasmGasCostTable[F64Add] = math.MaxInt64
	WasmGasCostTable[F64Sub] = math.MaxInt64
	WasmGasCostTable[F64Mul] = math.MaxInt64
	WasmGasCostTable[F64Div] = math.MaxInt64
	WasmGasCostTable[F64Min] = math.MaxInt64
	WasmGasCostTable[F64Max] = math.MaxInt64
	WasmGasCostTable[F64Copysign] = math.MaxInt64

	WasmGasCostTable[I32WrapI64] = 3
	WasmGasCostTable[I32TruncSF32] = 3
	WasmGasCostTable[I32TruncUF32] = 3
	WasmGasCostTable[I32TruncSF64] = 3
	WasmGasCostTable[I32TruncUF64] = 3
	WasmGasCostTable[I64ExtendSI32] = 3
	WasmGasCostTable[I64ExtendUI32] = 3
	WasmGasCostTable[I64TruncSF32] = 3
	WasmGasCostTable[I64TruncUF32] = 3
	WasmGasCostTable[I64TruncSF64] = 3
	WasmGasCostTable[I64TruncUF64] = 3
	WasmGasCostTable[F32ConvertSI32] = math.MaxInt64
	WasmGasCostTable[F32ConvertUI32] = math.MaxInt64
	WasmGasCostTable[F32ConvertSI64] = math.MaxInt64
	WasmGasCostTable[F32ConvertUI64] = math.MaxInt64
	WasmGasCostTable[F32DemoteF64] = math.MaxInt64
	WasmGasCostTable[F64ConvertSI32] = math.MaxInt64
	WasmGasCostTable[F64ConvertUI32] = math.MaxInt64
	WasmGasCostTable[F64ConvertSI64] = math.MaxInt64
	WasmGasCostTable[F64ConvertUI64] = math.MaxInt64
	WasmGasCostTable[F64PromoteF32] = math.MaxInt64

	WasmGasCostTable[I32ReinterpretF32] = 3
	WasmGasCostTable[I64ReinterpretF64] = 3
	WasmGasCostTable[F32ReinterpretI32] = math.MaxInt64
	WasmGasCostTable[F64ReinterpretI64] = math.MaxInt64
}

var WasmInstrString = map[byte]string{
	Unreachable:  "Unreachable",
	Nop:          "Nop",
	Block:        "Block",
	Loop:         "Loop",
	If:           "If",
	Else:         "Else",
	End:          "End",
	Br:           "Br",
	BrIf:         "BrIf",
	BrTable:      "BrTable",
	Return:       "Return",
	Call:         "Call",
	CallIndirect: "CallIndirect",

	Drop:   "Drop",
	Select: "Select",

	GetLocal:  "GetLocal",
	SetLocal:  "SetLocal",
	TeeLocal:  "TeeLocal",
	GetGlobal: "GetGlobal",
	SetGlobal: "SetGlobal",

	I32Load:    "I32Load",
	I64Load:    "I64Load",
	F32Load:    "F32Load",
	F64Load:    "F64Load",
	I32Load8s:  "I32Load8s",
	I32Load8u:  "I32Load8u",
	I32Load16s: "I32Load16s",
	I32Load16u: "I32Load16u",
	I64Load8s:  "I64Load8s",
	I64Load8u:  "I64Load8u",
	I64Load16s: "I64Load16s",
	I64Load16u: "I64Load16u",
	I64Load32s: "I64Load32s",
	I64Load32u: "I64Load32u",

	I32Store:   "I32Store",
	I64Store:   "I64Store",
	F32Store:   "F32Store",
	F64Store:   "F64Store",
	I32Store8:  "I32Store8",
	I32Store16: "I32Store16",
	I64Store8:  "I64Store8",
	I64Store16: "I64Store16",
	I64Store32: "I64Store32",

	CurrentMemory: "CurrentMemory",
	GrowMemory:    "GrowMemory",

	I32Const: "I32Const",
	I64Const: "I64Const",
	F32Const: "F32Const",
	F64Const: "F64Const",

	I32Eqz: "I32Eqz",
	I32Eq:  "I32Eq",
	I32Ne:  "I32Ne",
	I32LtS: "I32LtS",
	I32LtU: "I32LtU",
	I32GtS: "I32GtS",
	I32GtU: "I32GtU",
	I32LeS: "I32LeS",
	I32LeU: "I32LeU",
	I32GeS: "I32GeS",
	I32GeU: "I32GeU",
	I64Eqz: "I64Eqz",
	I64Eq:  "I64Eq",
	I64Ne:  "I64Ne",
	I64LtS: "I64LtS",
	I64LtU: "I64LtU",
	I64GtS: "I64GtS",
	I64GtU: "I64GtU",
	I64LeS: "I64LeS",
	I64LeU: "I64LeU",
	I64GeS: "I64GeS",
	I64GeU: "I64GeU",
	F32Eq:  "F32Eq",
	F32Ne:  "F32Ne",
	F32Lt:  "F32Lt",
	F32Gt:  "F32Gt",
	F32Le:  "F32Le",
	F32Ge:  "F32Ge",
	F64Eq:  "F64Eq",
	F64Ne:  "F64Ne",
	F64Lt:  "F64Lt",
	F64Gt:  "F64Gt",
	F64Le:  "F64Le",
	F64Ge:  "F64Ge",

	I32Clz:      "I32Clz",
	I32Ctz:      "I32Ctz",
	I32Popcnt:   "I32Popcnt",
	I32Add:      "I32Add",
	I32Sub:      "I32Sub",
	I32Mul:      "I32Mul",
	I32DivS:     "I32DivS",
	I32DivU:     "I32DivU",
	I32RemS:     "I32RemS",
	I32RemU:     "I32RemU",
	I32And:      "I32And",
	I32Or:       "I32Or",
	I32Xor:      "I32Xor",
	I32Shl:      "I32Shl",
	I32ShrS:     "I32ShrS",
	I32ShrU:     "I32ShrU",
	I32Rotl:     "I32Rotl",
	I32Rotr:     "I32Rotr",
	I64Clz:      "I64Clz",
	I64Ctz:      "I64Ctz",
	I64Popcnt:   "I64Popcnt",
	I64Add:      "I64Add",
	I64Sub:      "I64Sub",
	I64Mul:      "I64Mul",
	I64DivS:     "I64DivS",
	I64DivU:     "I64DivU",
	I64RemS:     "I64RemS",
	I64RemU:     "I64RemU",
	I64And:      "I64And",
	I64Or:       "I64Or",
	I64Xor:      "I64Xor",
	I64Shl:      "I64Shl",
	I64ShrS:     "I64ShrS",
	I64ShrU:     "I64ShrU",
	I64Rotl:     "I64Rotl",
	I64Rotr:     "I64Rotr",
	F32Abs:      "F32Abs",
	F32Neg:      "F32Neg",
	F32Ceil:     "F32Ceil",
	F32Floor:    "F32Floor",
	F32Trunc:    "F32Trunc",
	F32Nearest:  "F32Nearest",
	F32Sqrt:     "F32Sqrt",
	F32Add:      "F32Add",
	F32Sub:      "F32Sub",
	F32Mul:      "F32Mul",
	F32Div:      "F32Div",
	F32Min:      "F32Min",
	F32Max:      "F32Max",
	F32Copysign: "F32Copysign",
	F64Abs:      "F64Abs",
	F64Neg:      "F64Neg",
	F64Ceil:     "F64Ceil",
	F64Floor:    "F64Floor",
	F64Trunc:    "F64Trunc",
	F64Nearest:  "F64Nearest",
	F64Sqrt:     "F64Sqrt",
	F64Add:      "F64Add",
	F64Sub:      "F64Sub",
	F64Mul:      "F64Mul",
	F64Div:      "F64Div",
	F64Min:      "F64Min",
	F64Max:      "F64Max",
	F64Copysign: "F64Copysign",

	I32WrapI64:     "I32WrapI64",
	I32TruncSF32:   "I32TruncSF32",
	I32TruncUF32:   "I32TruncUF32",
	I32TruncSF64:   "I32TruncSF64",
	I32TruncUF64:   "I32TruncUF64",
	I64ExtendSI32:  "I64ExtendSI32",
	I64ExtendUI32:  "I64ExtendUI32",
	I64TruncSF32:   "I64TruncSF32",
	I64TruncUF32:   "I64TruncUF32",
	I64TruncSF64:   "I64TruncSF64",
	I64TruncUF64:   "I64TruncUF64",
	F32ConvertSI32: "F32ConvertSI32",
	F32ConvertUI32: "F32ConvertUI32",
	F32ConvertSI64: "F32ConvertSI64",
	F32ConvertUI64: "F32ConvertUI64",
	F32DemoteF64:   "F32DemoteF64",
	F64ConvertSI32: "F64ConvertSI32",
	F64ConvertUI32: "F64ConvertUI32",
	F64ConvertSI64: "F64ConvertSI64",
	F64ConvertUI64: "F64ConvertUI64",
	F64PromoteF32:  "F64PromoteF32",

	I32ReinterpretF32: "I32ReinterpretF32",
	I64ReinterpretF64: "I64ReinterpretF64",
	F32ReinterpretI32: "F32ReinterpretI32",
	F64ReinterpretI64: "F64ReinterpretI64",
}

const (
	Unreachable  = 0x00
	Nop          = 0x01
	Block        = 0x02
	Loop         = 0x03
	If           = 0x04
	Else         = 0x05
	End          = 0x0b
	Br           = 0x0c
	BrIf         = 0x0d
	BrTable      = 0x0e
	Return       = 0x0f
	Call         = 0x10
	CallIndirect = 0x11

	Drop   = 0x1a
	Select = 0x1b

	GetLocal  = 0x20
	SetLocal  = 0x21
	TeeLocal  = 0x22
	GetGlobal = 0x23
	SetGlobal = 0x24

	I32Load    = 0x28
	I64Load    = 0x29
	F32Load    = 0x2a
	F64Load    = 0x2b
	I32Load8s  = 0x2c
	I32Load8u  = 0x2d
	I32Load16s = 0x2e
	I32Load16u = 0x2f
	I64Load8s  = 0x30
	I64Load8u  = 0x31
	I64Load16s = 0x32
	I64Load16u = 0x33
	I64Load32s = 0x34
	I64Load32u = 0x35

	I32Store   = 0x36
	I64Store   = 0x37
	F32Store   = 0x38
	F64Store   = 0x39
	I32Store8  = 0x3a
	I32Store16 = 0x3b
	I64Store8  = 0x3c
	I64Store16 = 0x3d
	I64Store32 = 0x3e

	// TODO: rename operations accordingly

	CurrentMemory = 0x3f
	GrowMemory    = 0x40

	I32Const = 0x41
	I64Const = 0x42
	F32Const = 0x43
	F64Const = 0x44

	I32Eqz = 0x45
	I32Eq  = 0x46
	I32Ne  = 0x47
	I32LtS = 0x48
	I32LtU = 0x49
	I32GtS = 0x4a
	I32GtU = 0x4b
	I32LeS = 0x4c
	I32LeU = 0x4d
	I32GeS = 0x4e
	I32GeU = 0x4f
	I64Eqz = 0x50
	I64Eq  = 0x51
	I64Ne  = 0x52
	I64LtS = 0x53
	I64LtU = 0x54
	I64GtS = 0x55
	I64GtU = 0x56
	I64LeS = 0x57
	I64LeU = 0x58
	I64GeS = 0x59
	I64GeU = 0x5a
	F32Eq  = 0x5b
	F32Ne  = 0x5c
	F32Lt  = 0x5d
	F32Gt  = 0x5e
	F32Le  = 0x5f
	F32Ge  = 0x60
	F64Eq  = 0x61
	F64Ne  = 0x62
	F64Lt  = 0x63
	F64Gt  = 0x64
	F64Le  = 0x65
	F64Ge  = 0x66

	I32Clz      = 0x67
	I32Ctz      = 0x68
	I32Popcnt   = 0x69
	I32Add      = 0x6a
	I32Sub      = 0x6b
	I32Mul      = 0x6c
	I32DivS     = 0x6d
	I32DivU     = 0x6e
	I32RemS     = 0x6f
	I32RemU     = 0x70
	I32And      = 0x71
	I32Or       = 0x72
	I32Xor      = 0x73
	I32Shl      = 0x74
	I32ShrS     = 0x75
	I32ShrU     = 0x76
	I32Rotl     = 0x77
	I32Rotr     = 0x78
	I64Clz      = 0x79
	I64Ctz      = 0x7a
	I64Popcnt   = 0x7b
	I64Add      = 0x7c
	I64Sub      = 0x7d
	I64Mul      = 0x7e
	I64DivS     = 0x7f
	I64DivU     = 0x80
	I64RemS     = 0x81
	I64RemU     = 0x82
	I64And      = 0x83
	I64Or       = 0x84
	I64Xor      = 0x85
	I64Shl      = 0x86
	I64ShrS     = 0x87
	I64ShrU     = 0x88
	I64Rotl     = 0x89
	I64Rotr     = 0x8a
	F32Abs      = 0x8b
	F32Neg      = 0x8c
	F32Ceil     = 0x8d
	F32Floor    = 0x8e
	F32Trunc    = 0x8f
	F32Nearest  = 0x90
	F32Sqrt     = 0x91
	F32Add      = 0x92
	F32Sub      = 0x93
	F32Mul      = 0x94
	F32Div      = 0x95
	F32Min      = 0x96
	F32Max      = 0x97
	F32Copysign = 0x98
	F64Abs      = 0x99
	F64Neg      = 0x9a
	F64Ceil     = 0x9b
	F64Floor    = 0x9c
	F64Trunc    = 0x9d
	F64Nearest  = 0x9e
	F64Sqrt     = 0x9f
	F64Add      = 0xa0
	F64Sub      = 0xa1
	F64Mul      = 0xa2
	F64Div      = 0xa3
	F64Min      = 0xa4
	F64Max      = 0xa5
	F64Copysign = 0xa6

	I32WrapI64     = 0xa7
	I32TruncSF32   = 0xa8
	I32TruncUF32   = 0xa9
	I32TruncSF64   = 0xaa
	I32TruncUF64   = 0xab
	I64ExtendSI32  = 0xac
	I64ExtendUI32  = 0xad
	I64TruncSF32   = 0xae
	I64TruncUF32   = 0xaf
	I64TruncSF64   = 0xb0
	I64TruncUF64   = 0xb1
	F32ConvertSI32 = 0xb2
	F32ConvertUI32 = 0xb3
	F32ConvertSI64 = 0xb4
	F32ConvertUI64 = 0xb5
	F32DemoteF64   = 0xb6
	F64ConvertSI32 = 0xb7
	F64ConvertUI32 = 0xb8
	F64ConvertSI64 = 0xb9
	F64ConvertUI64 = 0xba
	F64PromoteF32  = 0xbb

	I32ReinterpretF32 = 0xbc
	I64ReinterpretF64 = 0xbd
	F32ReinterpretI32 = 0xbe
	F64ReinterpretI64 = 0xbf
)
