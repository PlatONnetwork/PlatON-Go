package exec

import (
	"Platon-go/common"
	"math/big"
)

// 定义库所需要的所有接口, 具体由wasmstatedb.go/WasmStateDB 来实现
// StateDB is an EVM database for full state querying.
type StateDB interface {
	GasPrice() int64
	BlockHash(num uint64) common.Hash
	BlockNumber() *big.Int
	GasLimimt() uint64
	Time() *big.Int
	Coinbase() common.Address
	GetBalance(addr common.Address) *big.Int
	Origin() common.Address
	Caller() common.Address
	Address() common.Address
}
