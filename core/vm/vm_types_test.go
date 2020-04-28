package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInterpTypeString(t *testing.T) {
	assert.Equal(t, "evm", EvmInterpOld.String())
	assert.Equal(t, "evm", EvmInterpNew.String())
	assert.Equal(t, "wasm", WasmInterp.String())
}

func TestINterpTypeBytes(t *testing.T) {
	assert.EqualValues(t, []byte{0x60, 0x60, 0x60, 0x40}, EvmInterpOld.Bytes())
	assert.EqualValues(t, []byte{0x60, 0x80, 0x60, 0x40}, EvmInterpNew.Bytes())
	assert.EqualValues(t, []byte{0x00, 0x61, 0x73, 0x6d}, WasmInterp.Bytes())
}

func TestBytesToInterpType(t *testing.T) {
	assert.EqualValues(t, EvmInterpOld, BytesToInterpType([]byte{0x60, 0x60, 0x60, 0x40}))
	assert.EqualValues(t, EvmInterpNew, BytesToInterpType([]byte{0x60, 0x80, 0x60, 0x40}))
	assert.EqualValues(t, WasmInterp, BytesToInterpType([]byte{0x00, 0x61, 0x73, 0x6d}))
}

func TestInterpSetBytes(t *testing.T) {
	var evminterpold InterpType
	var evminterpnew InterpType
	var wasminterp InterpType
	evminterpold.SetBytes([]byte{0x60, 0x60, 0x60, 0x40})
	evminterpnew.SetBytes([]byte{0x60, 0x80, 0x60, 0x40})
	wasminterp.SetBytes([]byte{0x00, 0x61, 0x73, 0x6d})

	assert.EqualValues(t, EvmInterpOld, evminterpold)
	assert.EqualValues(t, EvmInterpNew, evminterpnew)
	assert.EqualValues(t, WasmInterp, wasminterp)

}
