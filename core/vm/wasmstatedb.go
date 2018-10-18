package vm

import (
	"Platon-go/common"
	"math/big"
)

type WasmStateDB struct {
	StateDB StateDB
	evm     *EVM
	cfg     *Config
	contract *Contract
}

func (self *WasmStateDB) GasPrice() int64 {
	return self.evm.Context.GasPrice.Int64()
}

func (self *WasmStateDB) BlockHash(num uint64) common.Hash {
	// todo: 待细看
	return common.Hash{}
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
