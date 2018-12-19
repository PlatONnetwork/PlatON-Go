package vm

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"encoding/hex"
	"fmt"
	"math/big"
)

type WasmStateDB struct {
	StateDB  StateDB
	evm      *EVM
	cfg      *Config
	contract *Contract
}

func NewWasmStateDB(db *WasmStateDB, contract ContractRef) *WasmStateDB {
	stateDb := &WasmStateDB{
		StateDB: db.StateDB,
		evm: db.evm,
		cfg : db.cfg,
	}
	if c, ok := contract.(*Contract); ok {
		stateDb.contract = c
	}
	return stateDb
}


func (self *WasmStateDB) GasPrice() int64 {
	return self.evm.Context.GasPrice.Int64()
}

func (self *WasmStateDB) BlockHash(num uint64) common.Hash {
	return self.evm.GetHash(num)
}

func (self *WasmStateDB) BlockNumber() *big.Int {
	return self.evm.BlockNumber
}

func (self *WasmStateDB) GasLimimt() uint64 {
	return self.evm.GasLimit
}

func (self *WasmStateDB) Time() *big.Int {
	return self.evm.Time
}

func (self *WasmStateDB) Coinbase() common.Address {
	return self.evm.Coinbase
}

func (self *WasmStateDB) GetBalance(addr common.Address) *big.Int {
	return self.StateDB.GetBalance(addr)
}

func (self *WasmStateDB) Origin() common.Address {
	return self.evm.Origin
}

func (self *WasmStateDB) Caller() common.Address {
	return self.contract.Caller()
}

func (self *WasmStateDB) Address() common.Address {
	return self.contract.Address()
}

func (self *WasmStateDB) CallValue() *big.Int {
	return self.contract.Value()
}

/*func (self *WasmStateDB) AddLog(log *types.Log)  {
	self.evm.StateDB.AddLog(log)
}*/

func (self *WasmStateDB) AddLog(address common.Address, topics []common.Hash, data []byte, bn uint64)  {
	log := &types.Log {
		Address: address,
		Topics: topics,
		Data: data,
		BlockNumber: bn,
	}
	self.evm.StateDB.AddLog(log)
}


func (self *WasmStateDB) SetState(key []byte, value []byte)  {
	self.evm.StateDB.SetState(self.Address(), key, value)
}

func (self *WasmStateDB) GetState(key []byte) []byte {
	return self.evm.StateDB.GetState(self.Address(), key)
}

func (self *WasmStateDB) GetCallerNonce() int64 {
	addr := self.contract.Caller()
	return int64(self.evm.StateDB.GetNonce(addr))
}

func (self *WasmStateDB) Transfer(toAddr common.Address, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	caller := self.contract
	
	gas := self.evm.callGasTemp
	if value.Sign() != 0 {
		gas += params.CallStipend
	}
	fmt.Println("Transfer to:", toAddr.String())
	fmt.Println("Transfer caller:", caller.self.Address().Hex())
	ret, returnGas, err := self.evm.Call(caller, toAddr, nil, gas, value)
	return ret, returnGas, err
}

func (self *WasmStateDB) Call(addr, param []byte) ([]byte, error) {
	
	ret, _, err := self.evm.Call(self.contract, common.HexToAddress(hex.EncodeToString(addr)), param, self.contract.Gas, self.contract.value)
	return ret, err
}

func (self *WasmStateDB) DelegateCall(addr, param []byte) ([]byte, error) {
	
	ret, _, err := self.evm.DelegateCall(self.contract, common.HexToAddress(hex.EncodeToString(addr)), param, self.contract.Gas)
	return ret, err
}


