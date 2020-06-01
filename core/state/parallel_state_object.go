package state

import (
	"math/big"
)

type ParallelStateObject struct {
	stateObject *stateObject
	prevAmount  *big.Int
	createFlag  bool
}

func NewParallelStateObject(stateObject *stateObject, createFlag bool) *ParallelStateObject {
	return &ParallelStateObject{
		stateObject: stateObject,
		prevAmount:  new(big.Int).Set(stateObject.Balance()),
		createFlag:  createFlag,
	}
}

func (parallelObject *ParallelStateObject) GetNonce() uint64 {
	return parallelObject.stateObject.Nonce()
}

func (parallelObject *ParallelStateObject) SetNonce(nonce uint64) {
	parallelObject.stateObject.setNonce(nonce)
}

func (parallelObject *ParallelStateObject) GetBalance() *big.Int {
	return parallelObject.stateObject.Balance()
}

func (parallelObject *ParallelStateObject) AddBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	parallelObject.stateObject.setBalance(new(big.Int).Add(parallelObject.stateObject.Balance(), amount))
}

func (parallelObject *ParallelStateObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	parallelObject.stateObject.setBalance(new(big.Int).Sub(parallelObject.stateObject.Balance(), amount))
}

func (parallelObject *ParallelStateObject) UpdateRoot() {
	parallelObject.stateObject.updateRoot(parallelObject.stateObject.db.db)
}
