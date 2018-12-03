package vm

import (
	"Platon-go/common"
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
	panic("Not supported yet.")
}

func (self *WasmStateDB) BlockHash(num uint64) common.Hash {
	panic("Not supported yet.")
}

func (self *WasmStateDB) BlockNumber() *big.Int {
	panic("Not supported yet.")
}

func (self *WasmStateDB) GasLimimt() uint64 {
	panic("Not supported yet.")
}

func (self *WasmStateDB) Time() *big.Int {
	panic("Not supported yet.")
}

func (self *WasmStateDB) Coinbase() common.Address {
	panic("Not supported yet.")
}

func (self *WasmStateDB) GetBalance(addr common.Address) *big.Int {
	panic("Not supported yet.")
}

func (self *WasmStateDB) Origin() common.Address {
	panic("Not supported yet.")
}

func (self *WasmStateDB) Caller() common.Address {
	panic("Not supported yet.")
}

func (self *WasmStateDB) Address() common.Address {
	panic("Not supported yet.")
}

func (self *WasmStateDB) CallValue() int64 {
	panic("Not supported yet.")
}

func (self *WasmStateDB) AddLog(address common.Address, topics []common.Hash, data []byte, bn uint64)  {
	panic("Not supported yet.")
}

func (self *WasmStateDB) SetState(key []byte, value []byte)  {
	panic("Not supported yet.")
}

func (self *WasmStateDB) GetState(key []byte) []byte {
	panic("Not supported yet.")
}

func (self *WasmStateDB) GetCallerNonce() int64 {
	panic("Not supported yet.")
}

func (self *WasmStateDB) Transfer(toAddr common.Address, value *big.Int) (ret []byte, leftOverGas uint64, err error) {
	panic("Not supported yet.")
}

func (self *WasmStateDB) Call(addr, param []byte) ([]byte, error) {
	panic("Not supported yet.")
}

func (self *WasmStateDB) DelegateCall(addr, param []byte) ([]byte, error) {
	panic("Not supported yet.")
}


