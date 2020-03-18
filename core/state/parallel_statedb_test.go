package state

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
)

var (
	accountAddr = common.HexToAddress("0xeB3eb44a60d935DfE53D224648E1a51851c6f3Ae")
)

func TestParallelStateDB_justCreateObject(t *testing.T) {
	db := ethdb.NewMemDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))

	statedb.justCreateObject(accountAddr)

	if _, ok := statedb.stateObjects[accountAddr]; ok {
		t.Fatalf("The state object just created cannot exist in the cache")
	}
}

func TestParallelStateDB_GetOrNewParallelStateObject(t *testing.T) {
	db := ethdb.NewMemDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))
	statedb.GetOrNewParallelStateObject(accountAddr)
	if _, ok := statedb.stateObjects[accountAddr]; ok {
		t.Fatalf("The state object just created cannot exist in the cache")
	}
}

func TestParallelStateDB_justGetStateObject(t *testing.T) {
	db := ethdb.NewMemDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))
	stateObj := statedb.justGetStateObject(accountAddr)
	assert.Nil(t, stateObj)
	if _, ok := statedb.stateObjects[accountAddr]; ok {
		t.Fatalf("The state object just created cannot exist in the cache")
	}
}

func TestParallelStateDB_justGetStateObjectCache(t *testing.T) {
	db := ethdb.NewMemDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))
	stateObj := statedb.justGetStateObjectCache(accountAddr)
	assert.Nil(t, stateObj)
	if _, ok := statedb.stateObjects[accountAddr]; ok {
		t.Fatalf("The state object just created cannot exist in the cache")
	}
}
