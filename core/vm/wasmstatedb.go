package vm

import (
	"Platon-go/common"
	"Platon-go/core/types"
	"math/big"
)

type WasmStateDB struct {
	StateDB  StateDB
	evm      *EVM
	cfg      *Config
	contract *Contract
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

func (self *WasmStateDB) CallValue() int64 {
	return self.contract.Value().Int64()
}

func (self *WasmStateDB) AddLog(log *types.Log)  {
	self.evm.StateDB.AddLog(log)
}

func (self *WasmStateDB) SetState(key []byte, value []byte)  {
	self.evm.StateDB.SetState(self.Address(), key, value)
}

func (self *WasmStateDB) GetState(key []byte) []byte {
	return self.evm.StateDB.GetState(self.Address(), key)
}


