package vm

import (
	"context"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
)

func TestWASMInterpreterRun(t *testing.T) {
	// good deploy
	interp := &WASMInterpreter{
		evm: &EVM{Context: Context{
			CanTransfer: func(db StateDB, addr common.Address, amount *big.Int) bool {
				return db.GetBalance(addr).Cmp(amount) >= 0
			},
			Transfer: func(db StateDB, sender, recipient common.Address, amount *big.Int) {
				db.SubBalance(sender, amount)
				db.AddBalance(recipient, amount)
			},
			Ctx: context.TODO(),
		},
			StateDB: &mock.MockStateDB{
				Balance: map[common.Address]*big.Int{
					addr1: big.NewInt(2000),
					addr2: big.NewInt(1000),
				},
				State:    map[common.Address]map[string][]byte{},
				Code:     map[common.Address][]byte{},
				CodeHash: map[common.Address][]byte{},
				Journal:  mock.NewJournal(),
			}},
		cfg: Config{WasmType: Wagon},
	}
	contract := &Contract{
		self:           &AccountRef{1, 2, 3},
		Gas:            1000000,
		Code:           deployData(t, "init", "./testdata/contract_hello.wasm"),
		CodeAddr:       &addr2,
		CodeHash:       common.ZeroHash,
		DeployContract: true,
	}

	ret, err := interp.Run(contract, nil, false)
	assert.Nil(t, err)
	assert.NotNil(t, ret)
	interp.evm.StateDB.SetCode(addr2, ret)
	// good exec
	contract.DeployContract = false
	buf := callData(t, "add_message")
	ret, err = interp.Run(contract, buf, false)
	assert.Nil(t, err)
	assert.NotNil(t, ret)

	// bad exec
	contract.DeployContract = true
	buf = callData(t, "add_message")
	ret, err = interp.Run(contract, buf, false)
	assert.NotNil(t, err)
	assert.Nil(t, ret)

}

func TestWASMInterpreterCanRun(t *testing.T) {
	buf := append(EvmInterpOld.Bytes(), []byte{1, 2, 3}...)
	interp := &WASMInterpreter{}
	assert.Equal(t, false, interp.CanRun(buf))
}

func TestWasmInsTypeString(t *testing.T) {
	assert.Equal(t, "wagon", Wagon.String())
	assert.Equal(t, "unknown", Unknown.String())
}

func TestStr2WasmType(t *testing.T) {
	assert.EqualValues(t, Wagon, Str2WasmType("wagon"))
	assert.EqualValues(t, Unknown, Str2WasmType("ssss"))
}
