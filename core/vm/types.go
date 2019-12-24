package vm

// vm interpreter types
type InterpType byte

const (
	Illegal    InterpType = 0x00 // Illegal input in transaction
	PPOS       InterpType = 0x01
	EvmInterp  InterpType = 0x02
	WasmInterp InterpType = 0x03
)

func (t InterpType) String() string {

	switch t {
	case Illegal:
		return "illegal"
	case PPOS:
		return "ppos"
	case EvmInterp:
		return "evm"
	case WasmInterp:
		return "wasm"
	default:
		return "unknown"
	}
}

func (t InterpType) Byte() byte {
	return byte(t)
}

func Byte2Interp(typ byte) InterpType {
	switch typ {
	case PPOS.Byte():
		return PPOS
	case EvmInterp.Byte():
		return EvmInterp
	case WasmInterp.Byte():
		return WasmInterp
	default:
		return Illegal
	}
}

func Str2Interp(str string) InterpType {
	switch str {
	case PPOS.String():
		return PPOS
	case EvmInterp.String():
		return EvmInterp
	case WasmInterp.String():
		return WasmInterp
	default:
		return Illegal
	}
}
