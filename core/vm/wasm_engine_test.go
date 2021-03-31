package vm

import (
	"context"
	"hash/fnv"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
)

func TestWasmRun(t *testing.T) {

	// good deploy
	engine := &wagonEngine{
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
		config: Config{WasmType: Wagon},
		contract: &Contract{
			self:           &AccountRef{1, 2, 3},
			Gas:            1000000,
			Code:           deployData(t, "init", "./testdata/contract_hello.wasm"),
			CodeAddr:       &addr2,
			CodeHash:       common.ZeroHash,
			DeployContract: true,
		},
	}
	ret, err := engine.Run(nil, false)
	assert.Nil(t, err)
	assert.NotNil(t, ret)

	engine.evm.StateDB.SetCode(addr2, ret)

	// good call
	buf := callData(t, "add_message")
	engine.contract.DeployContract = false
	ret, err = engine.Run(buf, false)
	assert.Nil(t, err)
	assert.NotNil(t, ret)

	// read Only call
	// good call
	engine.contract.DeployContract = false
	ret, err = engine.Run(buf, true)
	assert.NotNil(t, err)
	assert.Nil(t, ret)

	// bad call for validate funcName
	buf = callData(t, "init")
	engine.contract.DeployContract = false
	ret, err = engine.Run(buf, false)
	assert.NotNil(t, err)
	assert.Nil(t, ret)

	// bad call for empty input
	buf = callData(t, "add_message")
	engine.contract.DeployContract = false
	ret, err = engine.Run(nil, false)
	assert.Nil(t, err)
	assert.Nil(t, ret)

	// bad deploy for validate funcName
	badEngine := &wagonEngine{
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
		config: Config{WasmType: Wagon},
		contract: &Contract{
			self:           &AccountRef{1, 2, 3},
			Gas:            1000000,
			Code:           deployData(t, "bad", "./testdata/contract_hello.wasm"),
			CodeAddr:       &addr2,
			CodeHash:       common.ZeroHash,
			DeployContract: true,
		},
	}
	ret, err = badEngine.Run(nil, false)
	assert.NotNil(t, err)
	assert.Nil(t, ret)

	// bad deploy for bad rlp
	badEngine2 := &wagonEngine{
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
		config: Config{WasmType: Wagon},
		contract: &Contract{
			self:           &AccountRef{1, 2, 3},
			Gas:            1000000,
			Code:           append(WasmInterp.Bytes(), []byte{1, 2, 3}...),
			CodeAddr:       &addr2,
			CodeHash:       common.ZeroHash,
			DeployContract: true,
		},
	}
	ret, err = badEngine2.Run(nil, false)
	assert.NotNil(t, err)
	assert.Nil(t, ret)

	// bad deploy for bad code
	badEngine3 := &wagonEngine{
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
		config: Config{WasmType: Wagon},
		contract: &Contract{
			self:           &AccountRef{1, 2, 3},
			Gas:            1000000,
			Code:           deployData(t, "init", "./testdata/bad.wasm"),
			CodeAddr:       &addr2,
			CodeHash:       common.ZeroHash,
			DeployContract: true,
		},
	}
	ret, err = badEngine3.Run(nil, false)
	assert.NotNil(t, err)
	assert.Nil(t, ret)
}

func deployData(t *testing.T, funcName, filePath string) []byte {

	buf, err := ioutil.ReadFile(filePath)
	assert.Nil(t, err)

	hash := fnv.New64()
	hash.Write([]byte(funcName))
	initUint64 := hash.Sum64()
	params := struct {
		FuncName uint64
	}{
		FuncName: initUint64,
	}

	bparams, err := rlp.EncodeToBytes(params)
	assert.Nil(t, err)

	arr := [][]byte{buf, bparams}
	barr, err := rlp.EncodeToBytes(arr)
	assert.Nil(t, err)

	interp := []byte{0x00, 0x61, 0x73, 0x6d}
	input := append(interp, barr...)
	return input
}

func callData(t *testing.T, funcName string) []byte {
	type M struct {
		Head string
	}

	type Message struct {
		M
		Body string
		End  string
	}

	hash := fnv.New64()
	hash.Write([]byte(funcName))
	funcUint64 := hash.Sum64()
	params := struct {
		FuncName uint64
		Msg      Message
	}{
		FuncName: funcUint64,
		Msg: Message{
			M: M{
				Head: "Gavin",
			},
			Body: "I am gavin",
			End:  "finished",
		},
	}

	bparams, err := rlp.EncodeToBytes(params)
	if nil != err {
		panic(err)
	}

	return bparams
}
