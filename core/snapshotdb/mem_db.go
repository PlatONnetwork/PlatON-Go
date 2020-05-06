// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

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
