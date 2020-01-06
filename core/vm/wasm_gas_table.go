package vm

const (
	IndirectCallGas = uint64(1)
	Sha3DataGas     = uint64(1)
	StoreLenGas     = uint64(1)
	StoreGas        = uint64(1)
	ExternalDataGas = uint64(1)
)

var wasmGasCostTable = map[byte]uint64{
	Unreachable:  0,
	Nop:          0,
	Block:        0,
	Loop:         0,
	If:           0,
	Else:         0,
	End:          0,
	Br:           0,
	BrIf:         0,
	BrTable:      0,
	Return:       0,
	Call:         0,
	CallIndirect: 0,

	Drop:   0,
	Select: 0,

	GetLocal:  0,
	SetLocal:  0,
	TeeLocal:  0,
	GetGlobal: 0,
	SetGlobal: 0,

	I32Load:    0,
	I64Load:    0,
	F32Load:    0,
	F64Load:    0,
	I32Load8s:  0,
	I32Load8u:  0,
	I32Load16s: 0,
	I32Load16u: 0,
	I64Load8s:  0,
	I64Load8u:  0,
	I64Load16s: 0,
	I64Load16u: 0,
	I64Load32s: 0,
	I64Load32u: 0,

	I32Store:   0,
	I64Store:   0,
	F32Store:   0,
	F64Store:   0,
	I32Store8:  0,
	I32Store16: 0,
	I64Store8:  0,
	I64Store16: 0,
	I64Store32: 0,

	CurrentMemory: 0,
	GrowMemory:    0,

	I32Const: 0,
	I64Const: 0,
	F32Const: 0,
	F64Const: 0,

	I32Eqz: 0,
	I32Eq:  0,
	I32Ne:  0,
	I32LtS: 0,
	I32LtU: 0,
	I32GtS: 0,
	I32GtU: 0,
	I32LeS: 0,
	I32LeU: 0,
	I32GeS: 0,
	I32GeU: 0,
	I64Eqz: 0,
	I64Eq:  0,
	I64Ne:  0,
	I64LtS: 0,
	I64LtU: 0,
	I64GtS: 0,
	I64GtU: 0,
	I64LeS: 0,
	I64LeU: 0,
	I64GeS: 0,
	I64GeU: 0,
	F32Eq:  0,
	F32Ne:  0,
	F32Lt:  0,
	F32Gt:  0,
	F32Le:  0,
	F32Ge:  0,
	F64Eq:  0,
	F64Ne:  0,
	F64Lt:  0,
	F64Gt:  0,
	F64Le:  0,
	F64Ge:  0,

	I32Clz:      0,
	I32Ctz:      0,
	I32Popcnt:   0,
	I32Add:      0,
	I32Sub:      0,
	I32Mul:      0,
	I32DivS:     0,
	I32DivU:     0,
	I32RemS:     0,
	I32RemU:     0,
	I32And:      0,
	I32Or:       0,
	I32Xor:      0,
	I32Shl:      0,
	I32ShrS:     0,
	I32ShrU:     0,
	I32Rotl:     0,
	I32Rotr:     0,
	I64Clz:      0,
	I64Ctz:      0,
	I64Popcnt:   0,
	I64Add:      0,
	I64Sub:      0,
	I64Mul:      0,
	I64DivS:     0,
	I64DivU:     0,
	I64RemS:     0,
	I64RemU:     0,
	I64And:      0,
	I64Or:       0,
	I64Xor:      0,
	I64Shl:      0,
	I64ShrS:     0,
	I64ShrU:     0,
	I64Rotl:     0,
	I64Rotr:     0,
	F32Abs:      0,
	F32Neg:      0,
	F32Ceil:     0,
	F32Floor:    0,
	F32Trunc:    0,
	F32Nearest:  0,
	F32Sqrt:     0,
	F32Add:      0,
	F32Sub:      0,
	F32Mul:      0,
	F32Div:      0,
	F32Min:      0,
	F32Max:      0,
	F32Copysign: 0,
	F64Abs:      0,
	F64Neg:      0,
	F64Ceil:     0,
	F64Floor:    0,
	F64Trunc:    0,
	F64Nearest:  0,
	F64Sqrt:     0,
	F64Add:      0,
	F64Sub:      0,
	F64Mul:      0,
	F64Div:      0,
	F64Min:      0,
	F64Max:      0,
	F64Copysign: 0,

	I32WrapI64:     0,
	I32TruncSF32:   0,
	I32TruncUF32:   0,
	I32TruncSF64:   0,
	I32TruncUF64:   0,
	I64ExtendSI32:  0,
	I64ExtendUI32:  0,
	I64TruncSF32:   0,
	I64TruncUF32:   0,
	I64TruncSF64:   0,
	I64TruncUF64:   0,
	F32ConvertSI32: 0,
	F32ConvertUI32: 0,
	F32ConvertSI64: 0,
	F32ConvertUI64: 0,
	F32DemoteF64:   0,
	F64ConvertSI32: 0,
	F64ConvertUI32: 0,
	F64ConvertSI64: 0,
	F64ConvertUI64: 0,
	F64PromoteF32:  0,

	I32ReinterpretF32: 0,
	I64ReinterpretF64: 0,
	F32ReinterpretI32: 0,
	F64ReinterpretI64: 0,
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
