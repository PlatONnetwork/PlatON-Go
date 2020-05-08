package rawdb

import (
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/ethdb/memorydb"
)

// NewMemoryDatabase creates an ephemeral in-memory key-value database without a
// freezer moving immutable chain segments into cold storage.
func NewMemoryDatabase() ethdb.Database {
	return NewDatabase(memorydb.New())
}

// NewDatabase creates a high level database on top of a given key-value data
// store without a freezer moving immutable chain segments into cold storage.
func NewDatabase(db ethdb.KeyValueStore) ethdb.Database {
	return &nofreezedb{
		KeyValueStore: db,
	}
}

// nofreezedb is a database wrapper that disables freezer data retrievals.
type nofreezedb struct {
	ethdb.KeyValueStore
}

func (n *nofreezedb) Close() {
	n.KeyValueStore.Close()
}
