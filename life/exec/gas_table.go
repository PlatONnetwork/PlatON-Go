package exec

import (
	"Platon-go/life/compiler/opcodes"
)

type (
	executionFunc func(vm *VirtualMachine, frame Frame)
	gasFunc       func(vm *VirtualMachine, frame *Frame) (uint64, error)
)

type Instruction struct {
	Execute executionFunc
	GasCost gasFunc
}

func constGasFunc(gas uint64) gasFunc {
	return func(vm *VirtualMachine, frame *Frame) (uint64, error) {
		return gas, nil
	}
}

func NewGasTable() [256]Instruction {
	return [256]Instruction{
		opcodes.Nop: {
			Execute: nil,
			GasCost: constGasFunc(0),
		},
		opcodes.Unreachable: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.Select: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Const: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Add: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Sub: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Mul: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32DivS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32DivU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32RemS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32RemU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32And: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Or: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Xor: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Shl: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32ShrS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32ShrU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Rotl: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Rotr: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Clz: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Ctz: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32PopCnt: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32EqZ: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Eq: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Ne: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32LtS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32LtU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32LeS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32LeU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32GtS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32GtU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32GeS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32GeU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.I64Const: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Add: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Sub: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Mul: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64DivS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64DivU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64RemS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64RemU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Rotl: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Rotr: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Clz: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Ctz: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64PopCnt: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64EqZ: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64And: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Or: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Xor: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Shl: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64ShrS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64ShrU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Eq: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Ne: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64LtS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64LtU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64LeS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64LeU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64GtS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64GtU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64GeS: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64GeU: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.F32Add: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Sub: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Mul: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Div: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Sqrt: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Min: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Max: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Ceil: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Floor: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Trunc: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Nearest: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Abs: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Neg: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32CopySign: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Eq: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Ne: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Lt: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Le: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Gt: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32Ge: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.F64Add: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Sub: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Mul: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Div: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Sqrt: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Min: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Max: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Ceil: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Floor: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Trunc: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Nearest: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Abs: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Neg: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64CopySign: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Eq: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Ne: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Lt: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Le: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Gt: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64Ge: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.I32WrapI64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32TruncUF32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32TruncUF64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32TruncSF32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32TruncSF64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64TruncUF32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64TruncUF64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64TruncSF32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64TruncSF64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64ExtendUI32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64ExtendSI32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.F32DemoteF64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64PromoteF32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32ConvertSI32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32ConvertSI64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32ConvertUI32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F32ConvertUI64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64ConvertSI32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64ConvertSI64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64ConvertUI32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.F64ConvertUI64: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.I32Load: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Load: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.I32Store: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Store: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.I32Load8S: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Load16S: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Load8S: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Load16S: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Load32S: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.I32Load8U: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Load16U: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Load8U: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Load16U: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Load32U: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.I32Store8: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I32Store16: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Store8: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Store16: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.I64Store32: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.Jmp: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.JmpIf: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.JmpEither: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.JmpTable: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.ReturnValue: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.ReturnVoid: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.GetLocal: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.SetLocal: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},

		opcodes.GetGlobal: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.SetGlobal: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.Call: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.CallIndirect: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.InvokeImport: {
			Execute: nil,
			GasCost: ImportGasFunc,
		},
		opcodes.CurrentMemory: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.GrowMemory: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.Phi: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.AddGas: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
		opcodes.Unknown: {
			Execute: nil,
			GasCost: constGasFunc(1),
		},
	}
}
