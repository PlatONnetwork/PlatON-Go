package state

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	accountAddr = common.HexToAddress("0xeB3eb44a60d935DfE53D224648E1a51851c6f3Ae")
)

func TestParallelStateDB_justCreateObject(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))

	statedb.justCreateObject(accountAddr)

	if _, ok := statedb.stateObjects[accountAddr]; ok {
		t.Fatalf("The state object just created cannot exist in the cache")
	}
}

func TestParallelStateDB_GetOrNewParallelStateObject(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))
	statedb.GetOrNewParallelStateObject(accountAddr)
	if _, ok := statedb.stateObjects[accountAddr]; ok {
		t.Fatalf("The state object just created cannot exist in the cache")
	}
}

func TestParallelStateDB_justGetStateObject(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))
	stateObj := statedb.justGetStateObject(accountAddr)
	assert.Nil(t, stateObj)
	if _, ok := statedb.stateObjects[accountAddr]; ok {
		t.Fatalf("The state object just created cannot exist in the cache")
	}
}

func TestParallelStateDB_justGetStateObjectCache(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))
	stateObj := statedb.justGetStateObjectCache(accountAddr)
	assert.Nil(t, stateObj)
	if _, ok := statedb.stateObjects[accountAddr]; ok {
		t.Fatalf("The state object just created cannot exist in the cache")
	}
}

func TestParallelStateDB_rlp(t *testing.T) {

	db := rawdb.NewMemoryDatabase()
	statedb, _ := New(common.Hash{}, NewDatabase(db))

	var count = 6000
	var objList = make([]*ParallelStateObject, count)

	for i := 0; i < count; i++ {
		prikey, _ := crypto.GenerateKey()
		address := crypto.PubkeyToAddress(prikey.PublicKey)
		pobj := statedb.GetOrNewParallelStateObject(address)
		pobj.SetNonce(uint64(i))
		pobj.AddBalance(big.NewInt(int64(300000000000)))
		objList[i] = pobj
	}

	start := time.Now()
	for i := 0; i < count; i++ {

		if _, err := rlp.EncodeToBytes(objList[i].stateObject); err != nil {
			t.Fatal("error")
		}
	}

	t.Logf("cost %s", common.PrettyDuration(time.Since(start)))

}

func TestParallelStateDB_random(t *testing.T) {
	fmt.Println(rand.Intn(100))
	fmt.Println(rand.Intn(100))

	for i := 0; i < 200; i++ {
		t.Logf("random number: %d \n", rand.Int31n(int32(1000)))
	}

}
