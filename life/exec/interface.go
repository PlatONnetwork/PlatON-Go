package exec

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"math/big"
)

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
	CallValue() *big.Int
	AddLog(address common.Address, topics []common.Hash, data []byte, bn uint64)
	SetState(key []byte, value []byte)
	GetState(key []byte) []byte

	GetCallerNonce() int64
	Transfer(addr common.Address, value *big.Int) (ret []byte, leftOverGas uint64, err error)
	DelegateCall(addr, params []byte) ([]byte, error)
	Call(addr, params []byte) ([]byte, error)
}
