package vm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"hash/fnv"
	"io/ioutil"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"golang.org/x/crypto/ripemd160"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/PlatONnetwork/PlatON-Go/common/math"

	"github.com/PlatONnetwork/wagon/exec"
	"github.com/stretchr/testify/assert"
)

var (
	addr1 = common.Address{1, 2, 3}
	addr2 = common.Address{1, 2, 4}
	addr3 = common.Address{1, 2, 6}
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
			return big.NewInt(12).Cmp(new(big.Int).SetBytes(ctx.Output)) == 0
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
			contract: &Contract{Gas: 99}},
		funcName: "platon_gas_test",
		check: func(ctx *VMContext, err error) bool {
			gas := binary.LittleEndian.Uint64(ctx.Output)
			return gas >= ctx.contract.Gas
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
				Coinbase: addr1,
			}}},
		funcName: "platon_coinbase_test",
		check: func(ctx *VMContext, err error) bool {
			return bytes.Equal(ctx.evm.Coinbase[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				Coinbase: addr1,
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
				Coinbase: addr1,
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
				Origin: addr1,
			}}},
		funcName: "platon_origin_test",
		check: func(ctx *VMContext, err error) bool {
			addr := addr1
			return bytes.Equal(addr[:], ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{caller: AccountRef{1, 2, 3}}},
		funcName: "platon_caller_test",
		check: func(ctx *VMContext, err error) bool {
			addr := addr1
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
			addr := addr1
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
			value := ctx.evm.StateDB.GetState(addr1, []byte("key"))
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
		init: func(ctx *VMContext, t *testing.T) {
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
			gasTable: params.GasTableConstantinople,
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
						addr1: big.NewInt(2000),
						addr2: big.NewInt(1000),
					},
					State: make(map[common.Address]map[string][]byte),
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: []byte{1}},
			Input:    addr2.Bytes(),
		},
		funcName: "platon_transfer_test",
		init: func(ctx *VMContext, t *testing.T) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			value := new(big.Int).SetBytes(ctx.Output)
			value = value.Add(value, big.NewInt(1000))
			return ctx.evm.StateDB.GetBalance(addr2).Cmp(value) == 0
		},
	},

	// CALL
	{
		ctx: &VMContext{
			gasTable: params.GasTableConstantinople,
			config:   Config{WasmType: Wagon},
			evm: &EVM{
				Context: Context{
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
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: []byte{1}},
			Input:    callContractInput(),
		},
		funcName: "platon_call_contract_test",
		init: func(ctx *VMContext, t *testing.T) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
			deployContract(ctx, addr1, addr2, readContractCode())
		},
		check: func(ctx *VMContext, err error) bool {
			value := new(big.Int).SetBytes(ctx.Output)
			value = value.Add(value, big.NewInt(1002))
			valueFlag := ctx.evm.StateDB.GetBalance(addr2).Cmp(value) == 0
			return valueFlag && checkContractRet(ctx.CallOut)
		},
	},

	// DELEGATECALL
	{
		ctx: &VMContext{
			gasTable: params.GasTableConstantinople,
			config:   Config{WasmType: Wagon},
			evm: &EVM{
				Context: Context{
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
						addr1: big.NewInt(2000), // from
						addr2: big.NewInt(1000), // to
					},
					State:    map[common.Address]map[string][]byte{},
					Code:     map[common.Address][]byte{},
					CodeHash: map[common.Address][]byte{},
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input:    callContractInput(),
		},
		funcName: "platon_delegate_call_contract_test",
		init: func(ctx *VMContext, t *testing.T) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
			deployContract(ctx, addr1, addr2, readContractCode())
		},
		check: func(ctx *VMContext, err error) bool {
			value := new(big.Int)
			value = value.Add(value, big.NewInt(1000))

			flag := ctx.evm.StateDB.GetBalance(addr2).Cmp(value) == 0

			return flag && checkContractRet(ctx.CallOut)
		},
	},

	// STATICCALL
	/*{
		ctx: &VMContext{
			gasTable: params.GasTableConstantinople,
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
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input:    queryContractInput(),
		},
		funcName: "platon_static_call_contract_test",
		init: func(ctx *VMContext, t *testing.T) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
			deployContract(ctx, readContractCode())
		},
		check: func(ctx *VMContext, err error) bool {
			value := new(big.Int)
			value = value.Add(value, big.NewInt(1000))
			valueFlag := ctx.evm.StateDB.GetBalance(addr2).Cmp(value) == 0
			return valueFlag && checkContractRet(ctx.CallOut)
		},
	},*/

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
						addr1: big.NewInt(2000),
						addr2: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						addr1: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						addr2: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
					Suicided: map[common.Address]bool{},
				}},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
		},
		funcName: "platon_destroy_contract_test",
		init: func(ctx *VMContext, t *testing.T) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
		},
		check: func(ctx *VMContext, err error) bool {
			to := addr1
			flag := ctx.evm.StateDB.GetBalance(addr3).Cmp(big.NewInt(2000)) == 0
			suicided := ctx.evm.StateDB.HasSuicided(to)
			return flag && suicided
		},
	},

	// MIGRATE
	{
		ctx: &VMContext{
			gasTable: params.GasTableConstantinople,
			config:   Config{WasmType: Wagon},
			evm: &EVM{
				Context: Context{
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
					State: map[common.Address]map[string][]byte{
						addr2: {
							"A": []byte("aaa"),
							"B": []byte("bbb"),
							"C": []byte("ccc"),
						},
					},
					Code:     map[common.Address][]byte{},
					CodeHash: map[common.Address][]byte{},
					Nonce:    map[common.Address]uint64{},
					Suicided: map[common.Address]bool{},
				}},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 4}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: func() []byte {

				code := readContractCode()

				hash := fnv.New64()
				hash.Write([]byte("init"))
				initUint64 := hash.Sum64()
				params := struct {
					FuncName uint64
				}{
					FuncName: initUint64,
				}
				input, err := rlp.EncodeToBytes(params)
				if nil != err {
					panic(err)
				}
				arr := [][]byte{code, input}
				barr, err := rlp.EncodeToBytes(arr)
				if nil != err {
					panic(err)
				}
				interp := []byte{0x00, 0x61, 0x73, 0x6d}
				return append(interp, barr...)
			}(),
		},
		funcName: "platon_migrate_contract_test",
		init: func(ctx *VMContext, t *testing.T) {
			ctx.evm.interpreters = append(ctx.evm.interpreters, NewWASMInterpreter(ctx.evm, ctx.config))
			deployContract(ctx, addr1, addr2, readContractCode())
		},
		check: func(ctx *VMContext, err error) bool {

			newContract := common.BytesToAddress(ctx.Output)
			//fmt.Println("new ContractAddr", newContract.String())
			newBalance := ctx.evm.StateDB.GetBalance(newContract)
			oldBalance := ctx.evm.StateDB.GetBalance(addr2)

			if oldBalance.Cmp(common.Big0) != 0 {
				return false
			}
			if newBalance.Cmp(big.NewInt(int64(1002))) != 0 {
				return false
			}

			count := 0
			storage := make(map[string][]byte)
			// check storage of newcontract
			ctx.evm.StateDB.ForEachStorage(newContract, func(key []byte, value []byte) bool {
				storage[string(key)] = []byte("aaa")
				count++
				return false
			})

			if len(storage) != 3 && count != 4 {
				return false
			}

			for _, key := range []string{} {

				value, ok := storage[key]
				if !ok {
					return false
				}

				if key == "A" && bytes.Compare(value, []byte("aaa")) != 0 {
					return false
				}
				if key == "B" && bytes.Compare(value, []byte("bbb")) != 0 {
					return false
				}
				if key == "C" && bytes.Compare(value, []byte("ccc")) != 0 {
					return false
				}
			}
			return true
		},
	},

	// EVENT, empty topic
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
						addr1: big.NewInt(2000),
						addr2: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						addr1: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						addr2: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
					Logs: make(map[common.Hash][]*types.Log),
				}},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: []byte("I am wagon"),
		},
		funcName: "platon_event0_test",
		init: func(ctx *VMContext, t *testing.T) {
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

	// EVENT3, had tree topics
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
						addr1: big.NewInt(2000),
						addr2: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						addr1: append(WasmInterp.Bytes(), []byte{0x00, 0x01}...),
						addr2: append(WasmInterp.Bytes(), []byte{0x00, 0x02}...),
					},
					Logs: make(map[common.Hash][]*types.Log),
				}},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: []byte("I am wagon"),
		},
		funcName: "platon_event3_test",
		init: func(ctx *VMContext, t *testing.T) {
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

			if len(log.Topics) == 0 {
				return false
			}
			if bytes.Compare(log.Data, []byte("I am wagon")) != 0 {
				return false
			}
			return true
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_sha256_test",
		check: func(ctx *VMContext, err error) bool {
			input := []byte{1, 2, 3}
			h := sha256.Sum256(input)
			return bytes.Equal(h[:], ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_ripemd160_test",
		check: func(ctx *VMContext, err error) bool {
			input := []byte{1, 2, 3}
			rip := ripemd160.New()
			rip.Write(input)
			h160 := rip.Sum(nil)
			return bytes.Equal(h160[:], ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_ecrecover_test",
		check: func(ctx *VMContext, err error) bool {
			var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
			key, _ := crypto.HexToECDSA(testPrivHex)
			//addr := common.HexToAddress(testAddrHex)

			msg := crypto.Keccak256([]byte("foo"))
			sig, err := crypto.Sign(msg, key)
			pubKey, _ := crypto.Ecrecover(msg, sig)

			addr := crypto.Keccak256(pubKey[1:])[12:]
			return bytes.Equal(addr[:], ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "rlp_u128_size_test",
		check: func(ctx *VMContext, err error) bool {
			res := []byte{17, 0, 0, 0, 0, 0, 0, 0}
			return bytes.Equal(res, ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_rlp_u128_test",
		check: func(ctx *VMContext, err error) bool {
			res := []byte{0x90, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
			return bytes.Equal(res, ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "rlp_bytes_size_test",
		check: func(ctx *VMContext, err error) bool {
			res := []byte{17, 0, 0, 0, 0, 0, 0, 0}
			return bytes.Equal(res, ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_rlp_bytes_test",
		check: func(ctx *VMContext, err error) bool {
			res := []byte{0x90, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
			return bytes.Equal(res, ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "rlp_list_size_test",
		check: func(ctx *VMContext, err error) bool {
			res := []byte{17, 0, 0, 0, 0, 0, 0, 0}
			return bytes.Equal(res, ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_rlp_list_test",
		check: func(ctx *VMContext, err error) bool {
			res := []byte{0xd0, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
			return bytes.Equal(res, ctx.Output)
		},
	},
}

func TestExternalFunction(t *testing.T) {
	buf, err := ioutil.ReadFile("./testdata/external.wasm")
	assert.Nil(t, err)
	module, err := ReadWasmModule(buf, false)
	assert.Nil(t, err)

	for i, c := range testCase {
		ExecCase(t, module, c, i)
	}
}

type Case struct {
	ctx      *VMContext
	init     func(*VMContext, *testing.T)
	funcName string
	check    func(*VMContext, error) bool
}

func ExecCase(t *testing.T, module *exec.CompiledModule, c *Case, i int) {

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
		c.init(c.ctx, t)
	}
	vm, err := exec.NewVMWithCompiled(module, 1024*1024)
	assert.Nil(t, err)

	vm.SetHostCtx(c.ctx)

	entry, ok := module.RawModule.Export.Entries[c.funcName]

	assert.True(t, ok, c.funcName)

	index := int64(entry.Index)
	vm.RecoverPanic = true

	_, err = vm.ExecCode(index)

	if c.funcName != "platon_panic_test" {
		if !assert.Nil(t, err) {
			t.Log("funcName:", c.funcName)
		}
	} else {
		assert.NotNil(t, err)
	}

	assert.True(t, c.check(c.ctx, err), "test failed "+c.funcName)
}

func readContractCode() []byte {
	buf, err := ioutil.ReadFile("./testdata/contract_hello.wasm")
	if nil != err {
		panic(err)
	}
	return buf
}

func deployContract(ctx *VMContext, addr1, addr2 common.Address, code []byte) {
	hash := fnv.New64()
	hash.Write([]byte("init"))
	initUint64 := hash.Sum64()
	params := struct {
		FuncName uint64
	}{
		FuncName: initUint64,
	}
	input, err := rlp.EncodeToBytes(params)
	if nil != err {
		panic(err)
	}
	var (
		sender  = AccountRef(addr1)
		receive = AccountRef(addr2)
	)

	arr := [][]byte{code, input}
	barr, err := rlp.EncodeToBytes(arr)
	if nil != err {
		panic(err)
	}
	interp := []byte{0x00, 0x61, 0x73, 0x6d}
	code = append(interp, barr...)

	contract := NewContract(sender, receive, common.Big0, 1000000)
	contract.SetCallCode(&addr2, ctx.evm.StateDB.GetCodeHash(addr2), code)
	contract.DeployContract = true
	ret, err := run(ctx.evm, contract, nil, false)
	if nil != err {
		panic(err)
	}
	ctx.evm.StateDB.SetCode(addr2, ret)
}

func callContractInput() []byte {
	type M struct {
		Head string
	}

	type Message struct {
		M
		Body string
		End  string
	}

	hash := fnv.New64()
	hash.Write([]byte("add_message"))
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

func checkContractRet(ret []byte) bool {
	if len(ret) == 0 {
		return false
	}
	type M struct {
		Head string
	}

	type Message struct {
		M
		Body string
		End  string
	}
	var arr []Message
	er := rlp.DecodeBytes(ret, &arr)
	if nil != er {
		return false
	}

	if len(arr) != 1 {
		return false
	}

	if arr[0].Head != "Gavin" || arr[0].Body != "I am gavin" || arr[0].End != "finished" {
		return false
	}
	return true
}
