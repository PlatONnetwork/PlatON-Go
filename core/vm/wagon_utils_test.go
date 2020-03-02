package vm

import (
	"hash/fnv"
	"io/ioutil"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
)

func TestReadWasmModule(t *testing.T) {
	buf, err := ioutil.ReadFile("./testdata/contract1.wasm")
	assert.Nil(t, err)
	module, err := ReadWasmModule(buf, true)
	assert.Nil(t, err)
	assert.NotNil(t, module)

	buf, err = ioutil.ReadFile("./testdata/bad.wasm")
	assert.Nil(t, err)
	module, err = ReadWasmModule(buf, true)
	assert.NotNil(t, err)
	assert.Nil(t, module)
}

func TestDecodeFuncAndParams(t *testing.T) {

	hash := fnv.New64()
	hash.Write([]byte("init"))
	initUint64 := hash.Sum64()

	params1 := struct {
		FuncName uint64
		Age      uint64
	}{
		FuncName: initUint64,
		Age:      16,
	}

	b1, _ := rlp.EncodeToBytes(params1)
	name1, _, err := decodeFuncAndParams(b1)
	assert.Nil(t, err)


	assert.Equal(t, initUint64, name1)

	type m struct {
		Content string
	}
	params2 := struct {
		M   m
		Age uint64
	}{
		M: m{
			Content: "init",
		},
		Age: 16,
	}

	b2, _ := rlp.EncodeToBytes(params2)

	name2, _, err := decodeFuncAndParams(b2)
	assert.NotNil(t, err)
	assert.NotEqual(t, initUint64, name2)

}
