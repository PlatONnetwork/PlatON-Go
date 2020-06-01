package state

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func TestParallelStateObject(t *testing.T) {
	account := Account{
		Root:    common.HexToHash("0x1000000000000000000000000000000000000001"),
		Nonce:   1,
		Balance: common.Big100,
	}
	stateObject := newObject(nil, accountAddr, account)
	parallelStateObj := NewParallelStateObject(stateObject, false)
	assert.Equal(t, uint64(1), parallelStateObj.GetNonce())
	assert.Equal(t, uint64(100), parallelStateObj.GetBalance().Uint64())

	parallelStateObj.AddBalance(common.Big32)
	assert.Equal(t, uint64(132), parallelStateObj.GetBalance().Uint64())

	parallelStateObj.SubBalance(common.Big3)
	assert.Equal(t, uint64(129), parallelStateObj.GetBalance().Uint64())
}
