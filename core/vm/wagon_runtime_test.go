package vm

import (
	"bytes"
	"encoding/binary"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/wagon/exec"
	"github.com/stretchr/testify/assert"
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
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			value := new(big.Int).SetBytes(ctx.Output)
			value = value.Add(value, big.NewInt(1000))
			return ctx.evm.StateDB.GetBalance(common.Address{1, 2, 4}).Cmp(value) == 0
		},
	},

	// CALL
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
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						common.Address{1, 2, 3}: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						common.Address{1, 2, 4}: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: []byte{1}},
			Input:    common.Address{1, 2, 4}.Bytes(),
		},
		funcName: "platon_call_contract_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			value := new(big.Int).SetBytes(ctx.Output)
			value = value.Add(value, big.NewInt(1000))
			return ctx.evm.StateDB.GetBalance(common.Address{1, 2, 4}).Cmp(value) == 0
		},
	},

	// DELEGATECALL
	{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
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
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						common.Address{1, 2, 3}: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						common.Address{1, 2, 4}: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input:    common.Address{1, 2, 4}.Bytes(), // todo need to add delegatecall input
		},
		funcName: "platon_delegatecall_contract_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			value := new(big.Int) /*.SetBytes(ctx.Output)*/
			value = value.Add(value, big.NewInt(1000))

			flag := ctx.evm.StateDB.GetBalance(common.Address{1, 2, 4}).Cmp(value) == 0

			return flag
		},
	},

	// STATICCALL
	{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
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
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						common.Address{1, 2, 3}: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						common.Address{1, 2, 4}: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input:    common.Address{1, 2, 4}.Bytes(), // todo need to add staticcall input
		},
		funcName: "platon_staticcall_contract_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			value := new(big.Int) /*.SetBytes(ctx.Output)*/
			value = value.Add(value, big.NewInt(1000))

			flag := ctx.evm.StateDB.GetBalance(common.Address{1, 2, 4}).Cmp(value) == 0

			return flag
		},
	},

	// DESTROY
	{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
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
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						common.Address{1, 2, 3}: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						common.Address{1, 2, 4}: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
					Suicided: map[common.Address]bool{},
				}},
			contract: &Contract{
				CallerAddress: common.Address{1, 2, 4},
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
		},
		funcName: "platon_destroy_contract_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			sender := common.Address{1, 2, 4}
			to := common.Address{1, 2, 3}

			flag := ctx.evm.StateDB.GetBalance(sender).Cmp(big.NewInt(3000)) == 0
			suicided := ctx.evm.StateDB.HasSuicided(to)
			return flag && suicided
		},
	},

	// MIGRATE todo
	/*{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
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
						// old contract
						common.Address{1, 2, 3}: big.NewInt(2000),
						// EOA
						common.Address{1, 2, 4}: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{
						common.Address{1, 2, 3}: {
							"A": []byte("aaa"),
							"B": []byte("bbb"),
							"C": []byte("ccc"),
						},
					},
					Code: map[common.Address][]byte{
						common.Address{1, 2, 3}: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
					},
					CodeHash: map[common.Address][]byte{
						common.Address{1, 2, 3}: func() []byte {
							var h common.Hash
							hw := sha3.NewKeccak256()
							rlp.Encode(hw, append(WasmInterp.Bytes(), []byte{0x00, 0x01}...))
							hw.Sum(h[:0])
							return h[:]
						}(),
					},
					Nonce:    map[common.Address]uint64{},
					Suicided: map[common.Address]bool{},
				}},
			contract: &Contract{
				CallerAddress: common.Address{1, 2, 4},
				caller:        &AccountRef{1, 2, 3},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: common.Address{1, 2, 3}.Bytes(),
		},
		funcName: "platon_migrate_contract_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {

			newContract := common.BytesToAddress(ctx.Output)

			newBalance := ctx.evm.StateDB.GetBalance(newContract)
			oldBalance := ctx.evm.StateDB.GetBalance(common.Address{1, 2, 3})

			if oldBalance.Cmp(common.Big0) != 0 {
				return false
			}
			if newBalance.Cmp(big.NewInt(int64(2001))) != 0 {
				return false
			}

			flag := 1
			// check storage of newcontract
			ctx.evm.StateDB.ForEachStorage(newContract, func(key, value []byte) bool {

				if string(key) == "A" && bytes.Compare(value, []byte("aaa")) == 0 {
					flag &= 1
					return true
				}
				if string(key) == "B" && bytes.Compare(value, []byte("bbb")) == 0 {
					flag &= 1
					return true
				}
				if string(key) == "C" && bytes.Compare(value, []byte("ccc")) == 0 {
					flag &= 1
					return true
				}
				flag &= 0
				return false
			})

			return flag == 1
		},
	},*/

	// EVENT
	{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
			evm: &EVM{
				Context: Context{
					CanTransfer: func(db StateDB, addr common.Address, amount *big.Int) bool {
						return db.GetBalance(addr).Cmp(amount) >= 0
					},
					Transfer: func(db StateDB, sender, recipient common.Address, amount *big.Int) {
						db.SubBalance(sender, amount)
						db.AddBalance(recipient, amount)
					},
					BlockNumber: big.NewInt(13),
				},
				StateDB: &mock.MockStateDB{
					Thash:   common.Hash{1, 1, 1, 1},
					TxIndex: 1,
					Bhash:   common.Hash{2, 2, 2, 2},
					Balance: map[common.Address]*big.Int{
						common.Address{1, 2, 3}: big.NewInt(2000),
						common.Address{1, 2, 4}: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						common.Address{1, 2, 3}: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						common.Address{1, 2, 4}: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
					Logs: make(map[common.Hash][]*types.Log),
				}},
			contract: &Contract{
				CallerAddress: common.Address{1, 2, 4},
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: []byte("I am wagon"),
		},
		funcName: "platon_event_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			logs := ctx.evm.StateDB.GetLogs(common.Hash{1, 1, 1, 1})
			if len(logs) != 1 {
				return false
			}
			log := logs[0]
			if log.BlockNumber != big.NewInt(13).Uint64() {
				return false
			}
			if log.TxIndex != 1 {
				return false
			}
			if log.TxHash != (common.Hash{1, 1, 1, 1}) {
				return false
			}
			if log.BlockHash != (common.Hash{2, 2, 2, 2}) {
				return false
			}

			if len(log.Topics) != 0 {
				return false
			}
			if bytes.Compare(log.Data, []byte("I am wagon")) != 0 {
				return false
			}
			return true
		},
	},

	// EVENT1
	{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
			evm: &EVM{
				Context: Context{
					CanTransfer: func(db StateDB, addr common.Address, amount *big.Int) bool {
						return db.GetBalance(addr).Cmp(amount) >= 0
					},
					Transfer: func(db StateDB, sender, recipient common.Address, amount *big.Int) {
						db.SubBalance(sender, amount)
						db.AddBalance(recipient, amount)
					},
					BlockNumber: big.NewInt(13),
				},
				StateDB: &mock.MockStateDB{
					Thash:   common.Hash{1, 1, 1, 1},
					TxIndex: 1,
					Bhash:   common.Hash{2, 2, 2, 2},
					Balance: map[common.Address]*big.Int{
						common.Address{1, 2, 3}: big.NewInt(2000),
						common.Address{1, 2, 4}: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						common.Address{1, 2, 3}: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						common.Address{1, 2, 4}: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
					Logs: make(map[common.Hash][]*types.Log),
				}},
			contract: &Contract{
				CallerAddress: common.Address{1, 2, 4},
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: []byte("I am wagon"),
		},
		funcName: "platon_event1_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			logs := ctx.evm.StateDB.GetLogs(common.Hash{1, 1, 1, 1})
			if len(logs) != 1 {
				return false
			}
			log := logs[0]
			if log.BlockNumber != big.NewInt(13).Uint64() {
				return false
			}
			if log.TxIndex != 1 {
				return false
			}
			if log.TxHash != (common.Hash{1, 1, 1, 1}) {
				return false
			}
			if log.BlockHash != (common.Hash{2, 2, 2, 2}) {
				return false
			}

			if len(log.Topics) != 1 {
				return false
			}
			if bytes.Compare(log.Data, []byte("I am wagon")) != 0 {
				return false
			}
			return true
		},
	},

	// EVENT2
	{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
			evm: &EVM{
				Context: Context{
					CanTransfer: func(db StateDB, addr common.Address, amount *big.Int) bool {
						return db.GetBalance(addr).Cmp(amount) >= 0
					},
					Transfer: func(db StateDB, sender, recipient common.Address, amount *big.Int) {
						db.SubBalance(sender, amount)
						db.AddBalance(recipient, amount)
					},
					BlockNumber: big.NewInt(13),
				},
				StateDB: &mock.MockStateDB{
					Thash:   common.Hash{1, 1, 1, 1},
					TxIndex: 1,
					Bhash:   common.Hash{2, 2, 2, 2},
					Balance: map[common.Address]*big.Int{
						common.Address{1, 2, 3}: big.NewInt(2000),
						common.Address{1, 2, 4}: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						common.Address{1, 2, 3}: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						common.Address{1, 2, 4}: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
					Logs: make(map[common.Hash][]*types.Log),
				}},
			contract: &Contract{
				CallerAddress: common.Address{1, 2, 4},
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: []byte("I am wagon"),
		},
		funcName: "platon_event2_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			logs := ctx.evm.StateDB.GetLogs(common.Hash{1, 1, 1, 1})
			if len(logs) != 1 {
				return false
			}
			log := logs[0]
			if log.BlockNumber != big.NewInt(13).Uint64() {
				return false
			}
			if log.TxIndex != 1 {
				return false
			}
			if log.TxHash != (common.Hash{1, 1, 1, 1}) {
				return false
			}
			if log.BlockHash != (common.Hash{2, 2, 2, 2}) {
				return false
			}

			if len(log.Topics) != 2 {
				return false
			}
			if bytes.Compare(log.Data, []byte("I am wagon")) != 0 {
				return false
			}
			return true
		},
	},

	// EVENT3
	{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
			evm: &EVM{
				Context: Context{
					CanTransfer: func(db StateDB, addr common.Address, amount *big.Int) bool {
						return db.GetBalance(addr).Cmp(amount) >= 0
					},
					Transfer: func(db StateDB, sender, recipient common.Address, amount *big.Int) {
						db.SubBalance(sender, amount)
						db.AddBalance(recipient, amount)
					},
					BlockNumber: big.NewInt(13),
				},
				StateDB: &mock.MockStateDB{
					Thash:   common.Hash{1, 1, 1, 1},
					TxIndex: 1,
					Bhash:   common.Hash{2, 2, 2, 2},
					Balance: map[common.Address]*big.Int{
						common.Address{1, 2, 3}: big.NewInt(2000),
						common.Address{1, 2, 4}: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						common.Address{1, 2, 3}: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						common.Address{1, 2, 4}: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
					Logs: make(map[common.Hash][]*types.Log),
				}},
			contract: &Contract{
				CallerAddress: common.Address{1, 2, 4},
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: []byte("I am wagon"),
		},
		funcName: "platon_event3_test",
		init: func(ctx *VMContext) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			logs := ctx.evm.StateDB.GetLogs(common.Hash{1, 1, 1, 1})
			if len(logs) != 1 {
				return false
			}
			log := logs[0]
			if log.BlockNumber != big.NewInt(13).Uint64() {
				return false
			}
			if log.TxIndex != 1 {
				return false
			}
			if log.TxHash != (common.Hash{1, 1, 1, 1}) {
				return false
			}
			if log.BlockHash != (common.Hash{2, 2, 2, 2}) {
				return false
			}

			if len(log.Topics) != 3 {
				return false
			}
			if bytes.Compare(log.Data, []byte("I am wagon")) != 0 {
				return false
			}
			return true
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
	if c.ctx.contract == nil {
		c.ctx.contract = &Contract{
			Gas: math.MaxUint64,
		}
	} else {
		if c.ctx.contract.Gas == 0 {
			c.ctx.contract.Gas = math.MaxUint64
		}
	}

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
	assert.Nil(t, err)
	assert.True(t, c.check(c.ctx, err), "test failed "+c.funcName)
}
