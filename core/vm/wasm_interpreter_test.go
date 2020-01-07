package vm

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var abi = `{
	"version": "0.01",
	"abi": [{
			"method": "transfer",
			"args": [{
					"name": "from",
					"typeName": "Address",
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
		StateDB: stateDB{},
		Context: Context{
			GasLimit:    1000000,
			BlockNumber: big.NewInt(10),
		},
	}
	cfg := Config{}

	wasmInterpreter := NewWASMInterpreter(evm, cfg)

	code, _ := ioutil.ReadFile("..\\..\\life\\contract\\hello.wasm")

	contract := &Contract{
		CallerAddress: common.BigToAddress(big.NewInt(88888)),
		caller:        ContractRefCaller{},
		self:          ContractRefSelf{},
		Code:          code,
		Gas:           1000000,
		ABI:           []byte(abi),
	}
	// build input, {1}{transfer}{from}{to}{asset}
	input := genInput()
	wasmInterpreter.Run(contract, input, true)

}

type stateDB struct {
	StateDB
}

func (stateDB) CreateAccount(common.Address) {}

func (stateDB) SubBalance(common.Address, *big.Int) {}
func (stateDB) AddBalance(common.Address, *big.Int) {}
func (stateDB) GetBalance(common.Address) *big.Int  { return nil }

func (stateDB) GetNonce(common.Address) uint64  { return 0 }
func (stateDB) SetNonce(common.Address, uint64) {}

func (stateDB) GetCodeHash(common.Address) common.Hash { return common.Hash{} }
func (stateDB) GetCode(common.Address) []byte          { return nil }
func (stateDB) SetCode(common.Address, []byte)         {}
func (stateDB) GetCodeSize(common.Address) int         { return 0 }

// todo: new func for abi of contract.
func (stateDB) GetAbiHash(common.Address) common.Hash { return common.Hash{} }
func (stateDB) GetAbi(common.Address) []byte          { return nil }
func (stateDB) SetAbi(common.Address, []byte)         {}

func (stateDB) AddRefund(uint64)  {}
func (stateDB) SubRefund(uint64)  {}
func (stateDB) GetRefund() uint64 { return 0 }

func (stateDB) GetCommittedState(common.Address, []byte) []byte { return nil }
func (stateDB) GetState(common.Address, []byte) []byte          { return []byte("world+++++++**") }
func (stateDB) SetState(common.Address, []byte, []byte)         {}
func (stateDB) Suicide(common.Address) bool                     { return true }
func (stateDB) HasSuicided(common.Address) bool                 { return true }

// Exist reports whether the given account exists in state.
// Notably this should also return true for suicided accounts.
func (stateDB) Exist(common.Address) bool { return true }

// Empty returns whether the given account is empty. Empty
// is defined according to EIP161 (balance = nonce = code = 0).
func (stateDB) Empty(common.Address) bool { return true }

func (stateDB) RevertToSnapshot(int) {}
func (stateDB) Snapshot() int        { return 0 }

func (stateDB) AddPreimage(common.Hash, []byte) {}

func (stateDB) ForEachStorage(common.Address, func([]byte, []byte) bool) {}

func (stateDB) MigrateStorage(from, to common.Address) {}

func (stateDB) AddLog(*types.Log) {
	fmt.Println("add log")
}
