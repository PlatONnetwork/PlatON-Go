package vm

const InterpTypeLen = 4

// vm interpreter types
type InterpType [InterpTypeLen]byte

var (
	EvmInterpOld = BytesToInterpType([]byte{0x60, 0x60, 0x60, 0x40}) // 60606040  Match until evm 4.2 version
	EvmInterpNew = BytesToInterpType([]byte{0x60, 0x80, 0x60, 0x40}) // 60806040  Matches after evm 4.2 version
	WasmInterp   = BytesToInterpType([]byte{0x00, 0x61, 0x73, 0x6d}) // 0061736d  Matches wasm
)

func (t InterpType) String() string {

	switch t {
	case EvmInterpOld, EvmInterpNew:
		return "evm"
	case WasmInterp:
		return "wasm"
	default:
		return "unknown"
	}
}

func (t InterpType) Bytes() []byte {
	return t[:]
}

func BytesToInterpType(b []byte) InterpType {
	var a InterpType
	a.SetBytes(b)
	return a
}

func (t *InterpType) SetBytes(b []byte) {
	if len(b) > len(t) {
		b = b[len(b)-InterpTypeLen:]
	}
	copy(t[InterpTypeLen-len(b):], b)
}
