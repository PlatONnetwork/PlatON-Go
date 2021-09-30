package vm

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/binary"
	"encoding/json"
	"hash/fnv"
	"io/ioutil"
	"math/big"
	"strings"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

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
		check: func(self *Case, err error) bool {
			return big.NewInt(12).Cmp(new(big.Int).SetBytes(self.ctx.Output)) == 0
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				BlockNumber: big.NewInt(99),
				GetHash: func(u uint64) common.Hash {
					return common.Hash{1, 2, 3}
				}},
			},
		},
		funcName: "platon_block_hash_test",
		check: func(self *Case, err error) bool {
			hash := common.Hash{1, 2, 3}
			return bytes.Equal(hash[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				BlockNumber: big.NewInt(99),
			}}},
		funcName: "platon_block_number_test",
		check: func(self *Case, err error) bool {
			var res [8]byte
			binary.LittleEndian.PutUint64(res[:], self.ctx.evm.BlockNumber.Uint64())
			return bytes.Equal(res[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				GasLimit: 99,
			}}},
		funcName: "platon_gas_limit_test",
		check: func(self *Case, err error) bool {
			var res [8]byte
			binary.LittleEndian.PutUint64(res[:], self.ctx.evm.GasLimit)
			return bytes.Equal(res[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{Gas: 99}},
		funcName: "platon_gas_test",
		check: func(self *Case, err error) bool {
			gas := binary.LittleEndian.Uint64(self.ctx.Output)
			return gas >= self.ctx.contract.Gas
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				Time: big.NewInt(93),
			}}},
		funcName: "platon_timestamp_test",
		check: func(self *Case, err error) bool {
			var res [8]byte
			binary.LittleEndian.PutUint64(res[:], self.ctx.evm.Time.Uint64())
			return bytes.Equal(res[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				Coinbase: addr1,
			}}},
		funcName: "platon_coinbase_test",
		check: func(self *Case, err error) bool {
			return bytes.Equal(self.ctx.evm.Coinbase[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{
				Context: Context{
					Coinbase: addr1,
				},
				StateDB: &mock.MockStateDB{
					Balance: map[common.Address]*big.Int{common.Address{1}: big.NewInt(99)},
					Journal: mock.NewJournal(),
				},
			},
		},
		funcName: "platon_balance_test",
		check: func(self *Case, err error) bool {
			return big.NewInt(99).Cmp(new(big.Int).SetBytes(self.ctx.Output)) == 0
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{
				Context: Context{
					Coinbase: addr1,
				},
				StateDB: &mock.MockStateDB{
					Balance: map[common.Address]*big.Int{
						common.Address{1}: big.NewInt(99),
					},
					Journal: mock.NewJournal(),
				},
			},
		},
		funcName: "platon_balance_test",
		check: func(self *Case, err error) bool {
			return big.NewInt(99).Cmp(new(big.Int).SetBytes(self.ctx.Output)) == 0
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{Context: Context{
				Origin: addr1,
			}}},
		funcName: "platon_origin_test",
		check: func(self *Case, err error) bool {
			addr := addr1
			return bytes.Equal(addr[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{caller: &AccountRef{1, 2, 3}},
			evm: &EVM{Context: Context{
				BlockNumber: big.NewInt(99),
				GetHash: func(u uint64) common.Hash {
					return common.Hash{1, 2, 3}
				}},
				StateDB: &mock.MockStateDB{
					Balance: make(map[common.Address]*big.Int),
					State:   make(map[common.Address]map[string][]byte),
					Journal: mock.NewJournal(),
				},
			},
		},
		funcName: "platon_caller_test",
		init: func(self *Case, t *testing.T) {
			curAv := gov.ActiveVersionValue{ActiveVersion: params.FORKVERSION_0_11_0, ActiveBlock: 1}
			avList := []gov.ActiveVersionValue{curAv}
			avListBytes, _ := json.Marshal(avList)
			self.ctx.evm.StateDB.SetState(vm.GovContractAddr, gov.KeyActiveVersions(), avListBytes)
		},
		check: func(self *Case, err error) bool {
			addr := addr1
			return bytes.Equal(addr[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{value: big.NewInt(99)}},
		funcName: "platon_call_value_test",
		check: func(self *Case, err error) bool {
			return big.NewInt(99).Cmp(new(big.Int).SetBytes(self.ctx.Output)) == 0
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{self: &AccountRef{1, 2, 3}}},
		funcName: "platon_address_test",
		check: func(self *Case, err error) bool {
			addr := addr1
			return bytes.Equal(addr[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			contract: &Contract{self: &AccountRef{1, 2, 3}},
			evm: &EVM{
				StateDB: &mock.MockStateDB{
					Balance: map[common.Address]*big.Int{
						common.Address{1}: big.NewInt(99),
					},
					Journal: mock.NewJournal(),
				}}},
		funcName: "platon_caller_nonce_test",
		check: func(self *Case, err error) bool {
			var res [8]byte
			binary.LittleEndian.PutUint64(res[:], 0)
			return bytes.Equal(res[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			Input: []byte{1},
		},
		funcName: "platon_sha3_test",
		check: func(self *Case, err error) bool {
			value := crypto.Keccak256([]byte{1})
			return bytes.Equal(value[:], self.ctx.Output)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{StateDB: &mock.MockStateDB{
				Balance: make(map[common.Address]*big.Int),
				State:   make(map[common.Address]map[string][]byte),
				Journal: mock.NewJournal(),
			}},
			contract: &Contract{self: &AccountRef{1, 2, 3}},
			Input:    []byte{1, 2, 3},
		},
		funcName: "platon_set_state_test",
		check: func(self *Case, err error) bool {
			value := self.ctx.evm.StateDB.GetState(addr1, []byte("key"))
			return bytes.Equal(value[:], self.ctx.Input)
		},
	},
	{
		ctx: &VMContext{
			evm: &EVM{StateDB: &mock.MockStateDB{
				Balance: make(map[common.Address]*big.Int),
				State:   make(map[common.Address]map[string][]byte),
				Journal: mock.NewJournal(),
			}},
			contract: &Contract{self: &AccountRef{1, 2, 3}},
			Input:    []byte{1, 2, 3},
		},
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.StateDB.SetState(self.ctx.contract.Address(), []byte("key"), self.ctx.Input)
		},
		funcName: "platon_get_state_test",
		check: func(self *Case, err error) bool {
			return bytes.Equal(self.ctx.Output, self.ctx.Input)
		},
	},
	{
		ctx: &VMContext{
			CallOut: []byte{1, 2, 3},
		},
		funcName: "platon_get_call_output_test",
		check: func(self *Case, err error) bool {
			return bytes.Equal(self.ctx.Output, self.ctx.CallOut)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_revert_test",
		check: func(self *Case, err error) bool {
			return self.ctx.Revert == true
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_panic_test",
		check: func(self *Case, err error) bool {
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
					State:   make(map[common.Address]map[string][]byte),
					Journal: mock.NewJournal(),
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: []byte{1}},
			Input:    addr2.Bytes(),
		},
		funcName: "platon_transfer_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
		},
		check: func(self *Case, err error) bool {
			value := new(big.Int).SetBytes(self.ctx.Output)
			value = value.Add(value, big.NewInt(1000))
			return self.ctx.evm.StateDB.GetBalance(addr2).Cmp(value) == 0
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
					Journal:  mock.NewJournal(),
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: []byte{1}},
			Input:    callContractInput(),
		},
		funcName: "platon_call_contract_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
			deployContract(self.ctx, addr1, addr2, readContractCode())
		},
		check: func(self *Case, err error) bool {
			value := new(big.Int).SetBytes(self.ctx.Output)
			value = value.Add(value, big.NewInt(1002))
			valueFlag := self.ctx.evm.StateDB.GetBalance(addr2).Cmp(value) == 0
			return valueFlag && checkContractRet(self.ctx.CallOut)
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
					Journal:  mock.NewJournal(),
				}},
			contract: &Contract{self: &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input:    callContractInput(),
		},
		funcName: "platon_delegate_call_contract_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
			deployContract(self.ctx, addr1, addr2, readContractCode())
		},
		check: func(self *Case, err error) bool {
			value := new(big.Int)
			value = value.Add(value, big.NewInt(1000))

			flag := self.ctx.evm.StateDB.GetBalance(addr2).Cmp(value) == 0

			return flag && checkContractRet(self.ctx.CallOut)
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
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
			deployContract(ctx, readContractCode())
		},
		check: func(self *Case, err error) bool {
			value := new(big.Int)
			value = value.Add(value, big.NewInt(1000))
			valueFlag := self.ctx.evm.StateDB.GetBalance(addr2).Cmp(value) == 0
			return valueFlag && checkContractRet(self.ctx.CallOut)
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
					Journal:  mock.NewJournal(),
				}},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
		},
		funcName: "platon_destroy_contract_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
		},
		check: func(self *Case, err error) bool {
			to := addr1
			flag := self.ctx.evm.StateDB.GetBalance(addr3).Cmp(big.NewInt(2000)) == 0
			suicided := self.ctx.evm.StateDB.HasSuicided(to)
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
					Journal:  mock.NewJournal(),
				},
			},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 4},
				Gas:           1000000,
				Code:          WasmInterp.Bytes(),
			},
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
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
			deployContract(self.ctx, addr1, addr2, readContractCode())
		},
		check: func(self *Case, err error) bool {

			newContract := common.BytesToAddress(self.ctx.Output)
			//fmt.Println("new ContractAddr", newContract.String())
			newBalance := self.ctx.evm.StateDB.GetBalance(newContract)
			oldBalance := self.ctx.evm.StateDB.GetBalance(addr2)
			code := self.ctx.evm.StateDB.GetCode(newContract)
			codeBytes := readContractCode()

			if oldBalance.Cmp(common.Big0) != 0 {
				return false
			}
			if newBalance.Cmp(big.NewInt(int64(1002))) != 0 {
				return false
			}

			if bytes.Compare(code, codeBytes) != 0 {
				return false
			}

			count := 0
			storage := make(map[string][]byte)
			// check storage of newcontract
			self.ctx.evm.StateDB.ForEachStorage(newContract, func(key []byte, value []byte) bool {
				storage[string(key)] = value
				count++
				return false
			})

			if len(storage) != 4 || count != 4 {
				return false
			}

			for _, key := range []string{"A", "B", "C"} {

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

	// MIGRATE CLONE
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
					Code: map[common.Address][]byte{
						addr1: readContractCode(),
					},
					CodeHash: map[common.Address][]byte{},
					Nonce:    map[common.Address]uint64{},
					Suicided: map[common.Address]bool{},
					Journal:  mock.NewJournal(),
				}},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 4},
				Gas:           1000000,
				Code:          WasmInterp.Bytes(),
			},
			Input: func() []byte {

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
				return input
			}(),
		},
		funcName: "platon_clone_migrate_contract_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
			deployContract(self.ctx, addr1, addr2, readContractCode())
		},
		check: func(self *Case, err error) bool {

			newContract := common.BytesToAddress(self.ctx.Output)
			//fmt.Println("new ContractAddr", newContract.String())
			newBalance := self.ctx.evm.StateDB.GetBalance(newContract)
			oldBalance := self.ctx.evm.StateDB.GetBalance(addr2)
			code := self.ctx.evm.StateDB.GetCode(newContract)
			codeBytes := readContractCode()

			if oldBalance.Cmp(common.Big0) != 0 {
				return false
			}
			if newBalance.Cmp(big.NewInt(int64(1002))) != 0 {
				return false
			}

			if bytes.Compare(code, codeBytes) != 0 {
				return false
			}

			count := 0
			storage := make(map[string][]byte)
			// check storage of newcontract
			self.ctx.evm.StateDB.ForEachStorage(newContract, func(key []byte, value []byte) bool {
				storage[string(key)] = value
				count++
				return false
			})

			if len(storage) != 4 || count != 4 {
				return false
			}

			for _, key := range []string{"A", "B", "C"} {

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

	// MIGRATE CLONE ERROR
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
					Code: map[common.Address][]byte{
						addr1: readContractCode(),
					},
					CodeHash: map[common.Address][]byte{},
					Nonce:    map[common.Address]uint64{},
					Suicided: map[common.Address]bool{},
					Journal:  mock.NewJournal(),
				}},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 4},
				Gas:           1000000,
				Code:          WasmInterp.Bytes(),
			},
			Input: func() []byte {

				hash := fnv.New64()
				hash.Write([]byte("init_error"))
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
				return input
			}(),
		},
		funcName: "platon_clone_migrate_contract_error_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
			deployContract(self.ctx, addr1, addr2, readContractCode())
		},
		check: func(self *Case, err error) bool {
			newContract := common.BytesToAddress(self.ctx.Output)
			if got := new(big.Int).SetBytes(newContract.Bytes()); got.Cmp(common.Big0) != 0 {
				return false
			}
			newState := self.ctx.evm.StateDB.(*mock.MockStateDB)
			oldState := self.stateDb.(*mock.MockStateDB)
			if !oldState.Equal(newState) {
				return false
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
					Logs:    make(map[common.Hash][]*types.Log),
					Journal: mock.NewJournal(),
				}},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: []byte("I am wagon"),
		},
		funcName: "platon_event0_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
		},
		check: func(self *Case, err error) bool {
			logs := self.ctx.evm.StateDB.GetLogs(common.Hash{1, 1, 1, 1})
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
					Logs:    make(map[common.Hash][]*types.Log),
					Journal: mock.NewJournal(),
				}},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 3}, Gas: 1000000, Code: WasmInterp.Bytes()},
			Input: []byte("I am wagon"),
		},
		funcName: "platon_event3_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
		},
		check: func(self *Case, err error) bool {
			logs := self.ctx.evm.StateDB.GetLogs(common.Hash{1, 1, 1, 1})
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
		check: func(self *Case, err error) bool {
			input := []byte{1, 2, 3}
			h := sha256.Sum256(input)
			return bytes.Equal(h[:], self.ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_ripemd160_test",
		check: func(self *Case, err error) bool {
			input := []byte{1, 2, 3}
			rip := ripemd160.New()
			rip.Write(input)
			h160 := rip.Sum(nil)
			return bytes.Equal(h160[:], self.ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_ecrecover_test",
		check: func(self *Case, err error) bool {
			var testPrivHex = "289c2857d4598e37fb9647507e47a309d6133539bf21a8b9cb6df88fd5232032"
			key, _ := crypto.HexToECDSA(testPrivHex)
			//addr := common.HexToAddress(testAddrHex)

			msg := crypto.Keccak256([]byte("foo"))
			sig, err := crypto.Sign(msg, key)
			pubKey, _ := crypto.Ecrecover(msg, sig)

			addr := crypto.Keccak256(pubKey[1:])[12:]
			return bytes.Equal(addr[:], self.ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "rlp_u128_size_test",
		check: func(self *Case, err error) bool {
			res := []byte{17, 0, 0, 0, 0, 0, 0, 0}
			return bytes.Equal(res, self.ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_rlp_u128_test",
		check: func(self *Case, err error) bool {
			res := []byte{0x90, 0x01, 0x23, 0x45, 0x67, 0x89, 0xab, 0xcd, 0xef, 0xfe, 0xdc, 0xba, 0x98, 0x76, 0x54, 0x32, 0x10}
			return bytes.Equal(res, self.ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "rlp_bytes_size_test",
		check: func(self *Case, err error) bool {
			res := []byte{17, 0, 0, 0, 0, 0, 0, 0}
			return bytes.Equal(res, self.ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_rlp_bytes_test",
		check: func(self *Case, err error) bool {
			res := []byte{0x90, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
			return bytes.Equal(res, self.ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "rlp_list_size_test",
		check: func(self *Case, err error) bool {
			res := []byte{17, 0, 0, 0, 0, 0, 0, 0}
			return bytes.Equal(res, self.ctx.Output)
		},
	},
	{
		ctx:      &VMContext{},
		funcName: "platon_rlp_list_test",
		check: func(self *Case, err error) bool {
			res := []byte{0xd0, 0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
			return bytes.Equal(res, self.ctx.Output)
		},
	},

	// platon_contract_code_length_test
	{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
			evm: &EVM{
				StateDB: &mock.MockStateDB{
					Balance: map[common.Address]*big.Int{},
					State:   map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						addr1: readContractCode(),
					},
					CodeHash: map[common.Address][]byte{},
					Nonce:    map[common.Address]uint64{},
					Suicided: map[common.Address]bool{},
					Journal:  mock.NewJournal(),
				},
			},
			contract: &Contract{
				Gas: 1000000,
			},
		},
		funcName: "platon_contract_code_length_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
		},
		check: func(self *Case, err error) bool {
			code := readContractCode()
			var length uint32 = uint32(len(code))
			rlpBytes, err := rlp.EncodeToBytes(length)
			if nil != err {
				return false
			}
			if bytes.Compare(rlpBytes, self.ctx.Output) != 0 {
				return false
			}

			return true
		},
	},

	// platon_contract_code_test
	{
		ctx: &VMContext{
			config: Config{WasmType: Wagon},
			evm: &EVM{
				StateDB: &mock.MockStateDB{
					Balance: map[common.Address]*big.Int{},
					State:   map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						addr1: []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f},
					},
					CodeHash: map[common.Address][]byte{},
					Nonce:    map[common.Address]uint64{},
					Suicided: map[common.Address]bool{},
					Journal:  mock.NewJournal(),
				},
			},
			contract: &Contract{
				Gas: 1000000,
			},
		},
		funcName: "platon_contract_code_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
		},
		check: func(self *Case, err error) bool {
			codeBytes := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x00a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f}
			if bytes.Compare(codeBytes, self.ctx.Output) != 0 {
				return false
			}

			return true
		},
	},

	// platon_deploy
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
						addr2: big.NewInt(1000),
					},
					State:    map[common.Address]map[string][]byte{},
					Code:     map[common.Address][]byte{},
					CodeHash: map[common.Address][]byte{},
					Nonce:    map[common.Address]uint64{},
					Suicided: map[common.Address]bool{},
					Journal:  mock.NewJournal(),
				},
			},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 4},
				Gas:           10000000000,
				Code:          WasmInterp.Bytes(),
			},
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
		funcName: "platon_deploy_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
		},
		check: func(self *Case, err error) bool {

			newContract := common.BytesToAddress(self.ctx.Output)
			newBalance := self.ctx.evm.StateDB.GetBalance(newContract)
			oldBalance := self.ctx.evm.StateDB.GetBalance(addr2)
			code := self.ctx.evm.StateDB.GetCode(newContract)
			codeBytes := readContractCode()

			if oldBalance.Cmp(big.NewInt(int64(998))) != 0 {
				return false
			}

			if newBalance.Cmp(big.NewInt(int64(2))) != 0 {
				return false
			}

			if bytes.Compare(code, codeBytes) != 0 {
				return false
			}

			return true
		},
	},

	// platon_clone
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
						addr2: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						addr1: readContractCode(),
					},
					CodeHash: map[common.Address][]byte{},
					Nonce:    map[common.Address]uint64{},
					Suicided: map[common.Address]bool{},
					Journal:  mock.NewJournal(),
				},
			},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 4},
				Gas:           10000000000,
				Code:          WasmInterp.Bytes(),
			},
			Input: func() []byte {
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
				return input
			}(),
		},
		funcName: "platon_clone_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
		},
		check: func(self *Case, err error) bool {

			newContract := common.BytesToAddress(self.ctx.Output)
			newBalance := self.ctx.evm.StateDB.GetBalance(newContract)
			oldBalance := self.ctx.evm.StateDB.GetBalance(addr2)
			code := self.ctx.evm.StateDB.GetCode(newContract)
			codeBytes := readContractCode()

			if oldBalance.Cmp(big.NewInt(int64(998))) != 0 {
				return false
			}

			if newBalance.Cmp(big.NewInt(int64(2))) != 0 {
				return false
			}

			if bytes.Compare(code, codeBytes) != 0 {
				return false
			}

			return true
		},
	},

	// platon clone error
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
						addr2: big.NewInt(1000),
					},
					State: map[common.Address]map[string][]byte{},
					Code: map[common.Address][]byte{
						addr1: readContractCode(),
					},
					CodeHash: map[common.Address][]byte{},
					Nonce:    map[common.Address]uint64{},
					Suicided: map[common.Address]bool{},
					Journal:  mock.NewJournal(),
				},
			},
			contract: &Contract{
				CallerAddress: addr2,
				caller:        &AccountRef{1, 2, 4},
				self:          &AccountRef{1, 2, 4},
				Gas:           10000000000,
				Code:          WasmInterp.Bytes(),
			},
			Input: func() []byte {
				hash := fnv.New64()
				hash.Write([]byte("init_error"))
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
				return input
			}(),
		},
		funcName: "platon_clone_error_test",
		init: func(self *Case, t *testing.T) {
			self.ctx.evm.interpreters = append(self.ctx.evm.interpreters, NewWASMInterpreter(self.ctx.evm, self.ctx.config))
		},
		check: func(self *Case, err error) bool {
			newContract := common.BytesToAddress(self.ctx.Output)
			if got := new(big.Int).SetBytes(newContract.Bytes()); got.Cmp(common.Big0) != 0 {
				return false
			}
			oldBalance := self.ctx.evm.StateDB.GetBalance(addr2)
			if oldBalance.Cmp(big.NewInt(int64(1000))) != 0 {
				return false
			}

			originState := mock.MockStateDB{
				Balance: map[common.Address]*big.Int{
					addr2: big.NewInt(1000),
				},
				State: map[common.Address]map[string][]byte{},
				Code: map[common.Address][]byte{
					addr1: readContractCode(),
				},
				CodeHash: map[common.Address][]byte{},
				Nonce:    map[common.Address]uint64{},
				Suicided: map[common.Address]bool{},
				Journal:  mock.NewJournal(),
			}

			newState := self.ctx.evm.StateDB.(*mock.MockStateDB)
			if !originState.Equal(newState) {
				return false
			}

			return true
		},
	},
}

var initExternalGas = uint64(10000000)

type testContract struct{}

func (testContract) Address() common.Address {
	return common.Address{}
}

func newTestVM(evm *EVM) *exec.VM {
	code := "0x0061736d010000000108026000006000017f03030200010405017001010105030100020615037f01418088040b7f00418088040b7f004180080b072c04066d656d6f727902000b5f5f686561705f6261736503010a5f5f646174615f656e640302046d61696e00010a090202000b0400412a0b004d0b2e64656275675f696e666f3d0000000400000000000401000000000c0023000000000000004300000005000000040000000205000000040000005c000000010439000000036100000005040000100e2e64656275675f6d6163696e666f0000400d2e64656275675f616262726576011101250e1305030e10171b0e110112060000022e0011011206030e3a0b3b0b49133f190000032400030e3e0b0b0b000000005e0b2e64656275675f6c696e654e000000040037000000010101fb0e0d0001010101000000010000012f746d702f6275696c645f7664717864336f336f316c2e24000066696c652e630001000000000502050000001505030a3d020100010100700a2e64656275675f737472636c616e672076657273696f6e20382e302e3020287472756e6b2033343139363029002f746d702f6275696c645f7664717864336f336f316c2e242f66696c652e63002f746d702f6275696c645f7664717864336f336f316c2e24006d61696e00696e74000021046e616d65011a0200115f5f7761736d5f63616c6c5f63746f727301046d61696e"
	module, _ := ReadWasmModule(hexutil.MustDecode(code), false)

	vm, _ := exec.NewVM(module.RawModule)
	vm.SetHostCtx(&VMContext{evm: evm, contract: NewContract(&testContract{}, &testContract{}, big.NewInt(0), initExternalGas)})
	return vm
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
	init     func(*Case, *testing.T)
	funcName string
	check    func(*Case, error) bool
	stateDb  StateDB
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
		c.init(c, t)
	}

	var (
		snapshot    int
		canRevert   bool = false
		mockStateDb *mock.MockStateDB
	)
	if nil != c.ctx && nil != c.ctx.evm && nil != c.ctx.evm.StateDB {
		if state, ok := c.ctx.evm.StateDB.(*mock.MockStateDB); ok {
			canRevert = true
			mockStateDb = state
			snapshot = state.Snapshot()

			stateTemp := &mock.MockStateDB{}
			stateTemp.DeepCopy(state)
			c.stateDb = stateTemp
		}
	}

	vm, err := exec.NewVMWithCompiled(module, 1024*1024)
	assert.Nil(t, err)

	vm.SetHostCtx(c.ctx)

	entry, ok := module.RawModule.Export.Entries[c.funcName]

	assert.True(t, ok, c.funcName)

	index := int64(entry.Index)
	vm.RecoverPanic = true

	_, err = vm.ExecCode(index)

	if c.funcName == "platon_panic_test" {
		assert.NotNil(t, err)
	} else if strings.Contains(c.funcName, "_error_test") {
		assert.NotNil(t, err)
		if canRevert {
			mockStateDb.RevertToSnapshot(snapshot)
		}
	} else {
		if !assert.Nil(t, err) {
			t.Log("funcName:", c.funcName)
		}
	}

	assert.True(t, c.check(c, err), "test failed "+c.funcName)
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

func TestGetBlockHash(t *testing.T) {

	testBlockHash := common.BytesToHash([]byte{1, 2, 3, 4})
	newProc := func(blockNumber int64) *exec.Process {
		return exec.NewProcess(newTestVM(&EVM{
			Context: Context{
				GetHash: func(u uint64) common.Hash {
					return testBlockHash
				},
				BlockNumber: big.NewInt(blockNumber),
			},
		}))
	}

	type TestCase struct {
		blockNumber int64
		getNumber   uint64
		expect      common.Hash
	}
	cases := []TestCase{
		{1, 123, common.Hash{}},
		{123, 123, common.Hash{}},
		{123, 122, testBlockHash},
		{1024, 122, common.Hash{}},
		{1024, 1024 - 256, testBlockHash},
	}
	for _, c := range cases {
		proc := newProc(c.blockNumber)
		BlockHash(proc, c.getNumber, 1024)
		res := common.Hash{}
		proc.ReadAt(res[:], 1024)
		assert.Equal(t, c.expect, res)
		assert.Equal(t, initExternalGas-GasExtStep, proc.HostCtx().(*VMContext).contract.Gas)
	}
}
