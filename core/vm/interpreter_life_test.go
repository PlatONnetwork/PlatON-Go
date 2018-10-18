
package vm

import (
	"Platon-go/common"
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"
)

var abi = `{
	"version": "0.01",
	"abi": [{
			"method": "transfer",
			"args": [{
					"name": "from",
					"typeName": "address",
					"realTypeName": "string"
				}, {
					"name": "to",
					"typeName": "address",
					"realTypeName": "string"
				}, {
					"name": "asset",
					"typeName": "",
					"realTypeName": "int64"
				}
			]
		}
	]
}`

func TestAddressUtil(t *testing.T) {
	ref := ContractRefSelf{}
	addr := ref.Address()
	fmt.Println(addr.Hex())
}

func TestWasmInterpreter(t *testing.T) {

	evm := &EVM{
		Context: Context{
			GasLimit: 1000000,
		},
		StateDB: nil,
	}
	cfg := Config{}

	wasmInterpreter := NewWASMInterpreter(evm, cfg)

	code, _ := ioutil.ReadFile("C:\\sunzone\\MyDocument\\goworkspace\\src\\Platon-go\\life\\contract\\hello.wasm")
	contract := &Contract{
		CallerAddress: common.BigToAddress(big.NewInt(88888)),
		caller: ContractRefCaller{},
		self: ContractRefSelf{},
		Code : code,
		Gas: 1000000,
		ABI: []byte(abi),
	}
	// 构建input, {1}{transfer}{from}{to}{asset}
	input := genInput()
	wasmInterpreter.Run(contract, input, true)




}