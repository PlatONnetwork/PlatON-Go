package snapshotdb

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"math/big"
	"sync"
)

type MemDatabase struct {
	db      map[string][]byte
	lock    sync.RWMutex
	current *current
}

func NewMemBaseDB() *MemDatabase {
	return &MemDatabase{
		db: make(map[string][]byte),
	}
}

func (db *MemDatabase) PutBaseDB(key, value []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.db[string(key)] = common.CopyBytes(value)
	return nil
}

func (db *MemDatabase) GetBaseDB(key []byte) ([]byte, error) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	if entry, ok := db.db[string(key)]; ok {
		return common.CopyBytes(entry), nil
	}
	return nil, ErrNotFound
}

func (db *MemDatabase) DelBaseDB(key []byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	delete(db.db, string(key))
	return nil
}

// WriteBaseDB apply the given [][2][]byte to the baseDB.
func (db *MemDatabase) WriteBaseDB(kvs [][2][]byte) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	for _, kv := range kvs {
		db.db[string(kv[0])] = common.CopyBytes(kv[1])
	}
	return nil
}

//SetCurrent use for fast sync
func (db *MemDatabase) SetCurrent(highestHash common.Hash, base, height big.Int) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.current = newCurrent(&height, &base, highestHash)
	return nil
}

func (db *MemDatabase) GetCurrent() *current {
	db.lock.Lock()
	defer db.lock.Unlock()
	return db.current
}
