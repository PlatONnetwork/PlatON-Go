package vm

import (
	"bytes"
	"encoding/binary"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/wagon/exec"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"math/big"
	"testing"
)

var testCase = []*Case{
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				GasPrice: big.NewInt(12)},
			},
		},
		funcName: "platon_gas_price_test",
		check: func(ctx *VMContext, err error) bool {
			var res [8]byte
			binary.LittleEndian.PutUint64(res[:], ctx.evm.GasPrice.Uint64())
			return bytes.Equal(res[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				GetHash: func(u uint64) common.Hash {
					return common.Hash{1, 2, 3}
				}},
			},
		},
		funcName: "platon_block_hash_test",
		check: func(ctx *VMContext, err error) bool {
			hash := common.Hash{1, 2, 3}
			return bytes.Equal(hash[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				BlockNumber: big.NewInt(99),
			}}},
		funcName: "platon_block_number_test",
		check: func(ctx *VMContext, err error) bool {
			var res [8]byte
			binary.LittleEndian.PutUint64(res[:], ctx.evm.BlockNumber.Uint64())
			return bytes.Equal(res[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				GasLimit: 99,
			}}},
		funcName: "platon_gas_limit_test",
		check: func(ctx *VMContext, err error) bool {
			var res [8]byte
			binary.LittleEndian.PutUint64(res[:], ctx.evm.GasLimit)
			return bytes.Equal(res[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				Time: big.NewInt(93),
			}}},
		funcName: "platon_timestamp_test",
		check: func(ctx *VMContext, err error) bool {
			var res [8]byte
			binary.LittleEndian.PutUint64(res[:], ctx.evm.Time.Uint64())
			return bytes.Equal(res[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				Coinbase: common.Address{1, 2, 3},
			}}},
		funcName: "platon_coinbase_test",
		check: func(ctx *VMContext, err error) bool {
			return bytes.Equal(ctx.evm.Coinbase[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				Coinbase: common.Address{1, 2, 3},
			},
				StateDB: &mock.MockStateDB{Balance: map[common.Address]*big.Int{
					common.Address{1}: big.NewInt(99),
				}}}},
		funcName: "platon_balance_test",
		check: func(ctx *VMContext, err error) bool {
			return big.NewInt(99).Cmp(new(big.Int).SetBytes(ctx.Output)) == 0
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				Coinbase: common.Address{1, 2, 3},
			},
				StateDB: &mock.MockStateDB{Balance: map[common.Address]*big.Int{
					common.Address{1}: big.NewInt(99),
				}}}},
		funcName: "platon_balance_test",
		check: func(ctx *VMContext, err error) bool {
			return big.NewInt(99).Cmp(new(big.Int).SetBytes(ctx.Output)) == 0
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				Origin: common.Address{1, 2, 3},
			}}},
		funcName: "platon_origin_test",
		check: func(ctx *VMContext, err error) bool {
			addr := common.Address{1, 2, 3}
			return bytes.Equal(addr[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{caller: AccountRef{1, 2, 3}}},
		funcName: "platon_caller_test",
		check: func(ctx *VMContext, err error) bool {
			addr := common.Address{1, 2, 3}
			return bytes.Equal(addr[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{value: big.NewInt(99)}},
		funcName: "platon_call_value_test",
		check: func(ctx *VMContext, err error) bool {
			return big.NewInt(99).Cmp(new(big.Int).SetBytes(ctx.Output)) == 0
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{self: &AccountRef{1, 2, 3}}},
		funcName: "platon_address_test",
		check: func(ctx *VMContext, err error) bool {
			addr := common.Address{1, 2, 3}
			return bytes.Equal(addr[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{self: &AccountRef{1, 2, 3}},
			evm: &EVM{
				StateDB: &mock.MockStateDB{Balance: map[common.Address]*big.Int{
					common.Address{1}: big.NewInt(99),
				}}}},
		funcName: "platon_caller_nonce_test",
		check: func(ctx *VMContext, err error) bool {
			var res [8]byte
			binary.LittleEndian.PutUint64(res[:], 0)
			return bytes.Equal(res[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			Input: []byte{1},
		},
		funcName: "platon_sha3_test",
		check: func(ctx *VMContext, err error) bool {
			value := crypto.Keccak256([]byte{1})
			return bytes.Equal(value[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{StateDB: &mock.MockStateDB{
				Balance: make(map[common.Address]*big.Int),
				State:   make(map[common.Address]map[string][]byte),
			}},
			contract: &Contract{self: &AccountRef{1, 2, 3}},
			Input:    []byte{1, 2, 3},
		},
		funcName: "platon_set_state_test",
		check: func(ctx *VMContext, err error) bool {
			value := ctx.evm.StateDB.GetState(common.Address{1, 2, 3}, []byte("key"))
			return bytes.Equal(value[:], ctx.Input)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{StateDB: &mock.MockStateDB{
				Balance: make(map[common.Address]*big.Int),
				State:   make(map[common.Address]map[string][]byte),
			}},
			contract: &Contract{self: &AccountRef{1, 2, 3}},
			Input:    []byte{1, 2, 3},
		},
		init: func(ctx *VMContext) {
			ctx.evm.StateDB.SetState(ctx.contract.Address(), []byte("key"), ctx.Input)
		},
		funcName: "platon_get_state_test",
		check: func(ctx *VMContext, err error) bool {
			return bytes.Equal(ctx.Output, ctx.Input)
		},
	},
	{
		ctx: &VMContext{
			CallOut: []byte{1, 2, 3},
		},
		funcName: "platon_get_call_output_test",
		check: func(ctx *VMContext, err error) bool {
			return bytes.Equal(ctx.Output, ctx.CallOut)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_revert_test",
		check: func(ctx *VMContext, err error) bool {
			return ctx.Revert == true
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_panic_test",
		check: func(ctx *VMContext, err error) bool {
			return true
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{
				Context: Context{
					CanTransfer: func(db StateDB, addr common.Address, amount *big.Int) bool {
						return db.GetBalance(addr).Cmp(amount) >= 0
					},
					Transfer: func(db StateDB, sender, recipient common.Address, amount *big.Int) {
						db.SubBalance(sender, amount)
						db.AddBalance(recipient, amount)
					},
				},
				StateDB: &mock.MockStateDB{
					Balance: map[common.Address]*big.Int{
						common.Address{1, 2, 3}: big.NewInt(2000),
						common.Address{1, 2, 4}: big.NewInt(1000),
					},
					State: make(map[common.Address]map[string][]byte),
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: []byte{1}},
			Input:    common.Address{1, 2, 4}.Bytes(),
		},
		funcName: "platon_transfer_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewEVMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			value := new(big.Int).SetBytes(ctx.Output)
			value = value.Add(value, big.NewInt(1000))
			return ctx.evm.StateDB.GetBalance(common.Address{1, 2, 4}).Cmp(value) == 0
		},
	},
}

func TestExternalFunction(t *testing.T) {
	buf, err := ioutil.ReadFile("./testdata/external.wasm")
	assert.Nil(t, err)
	module, err := ReadWasmModule(buf, false)
	assert.Nil(t, err)

	for _, c := range testCase {
		ExecCase(t, module, c)
	}
}

type Case struct {
	ctx      *VMContext
	init     func(*VMContext)
	funcName string
	check    func(*VMContext, error) bool
}

func ExecCase(t *testing.T, module *exec.CompiledModule, c *Case) {
	if c.init != nil {
		c.init(c.ctx)
	}
	vm, err := exec.NewVMWithCompiled(module, 1024*1024)
	assert.Nil(t, err)

	vm.SetHostCtx(c.ctx)

	entry, ok := module.RawModule.Export.Entries[c.funcName]

	assert.True(t, ok, c.funcName)

	index := int64(entry.Index)
	vm.RecoverPanic = true
	_, err = vm.ExecCode(index)
	assert.True(t, c.check(c.ctx, err), "test failed "+c.funcName)
}
