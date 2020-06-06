// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package trie

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/ethdb/memorydb"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
)

func newEmptySecure() *SecureTrie {
	trie, _ := NewSecure(common.Hash{}, NewDatabase(memorydb.New()))
	return trie
}

// makeTestSecureTrie creates a large enough secure trie for testing.
func makeTestSecureTrie() (*Database, *SecureTrie, map[string][]byte) {
	// Create an empty trie
	triedb := NewDatabase(memorydb.New())

	trie, _ := NewSecure(common.Hash{}, triedb)

	// Fill it with some arbitrary data
	content := make(map[string][]byte)
	for i := byte(0); i < 255; i++ {
		// Map the same data under multiple keys
		key, val := common.LeftPadBytes([]byte{1, i}, 32), []byte{i}
		content[string(key)] = val
		trie.Update(key, val)

		key, val = common.LeftPadBytes([]byte{2, i}, 32), []byte{i}
		content[string(key)] = val
		trie.Update(key, val)

		// Add some other data to inflate the trie
		for j := byte(3); j < 13; j++ {
			key, val = common.LeftPadBytes([]byte{j, i}, 32), []byte{j, i}
			content[string(key)] = val
			trie.Update(key, val)
		}
	}
	trie.Commit(nil)

	// Return the generated trie
	return triedb, trie, content
}

func TestInsertBlob(t *testing.T) {
	triedb1, triedb2 := NewDatabase(memorydb.New()), NewDatabase(memorydb.New())
	trie1, _ := NewSecure(common.Hash{}, triedb1)
	trie2, _ := NewSecure(common.Hash{}, triedb2)

	storages := make(map[string][]byte)
	valueKeys := make(map[common.Hash][]byte)

	insert := func(t1 *SecureTrie, t2 *SecureTrie) {
		key := randBytes(20)
		value := randBytes(30)
		valueKey := crypto.Keccak256([]byte(value))

		t1.Update(key, valueKey)
		t2.Update(key, valueKey)
		storages[string(key)] = value
		valueKeys[common.BytesToHash(valueKey)] = value
	}

	leafcallback := func(database *Database, st *SecureTrie) func(leaf []byte, parent common.Hash) error {
		return func(leaf []byte, parent common.Hash) error {
			valuekey := common.BytesToHash(leaf)
			database.InsertBlob(valuekey, valueKeys[valuekey])
			database.Reference(valuekey, parent)
			return nil
		}
	}

	start := time.Now()

	for i := 0; i < 2000; i++ {
		insert(trie1, trie2)
	}
	fmt.Println("duration", time.Since(start))
	start = time.Now()

	root1, _ := trie1.Commit(leafcallback(triedb1, trie1))
	fmt.Println("duration", time.Since(start))

	root2, _ := trie2.Commit(leafcallback(triedb2, trie2))

	assert.Equal(t, root1, root2)
	triedb1.Commit(root1, true, true)
	triedb2.Commit(root2, true, true)
	for k, v := range storages {
		valueKey1 := trie1.Get([]byte(k))
		valueKey2 := trie2.Get([]byte(k))
		assert.Equal(t, valueKey1, valueKey2)
		value1 := trie1.GetKey(valueKey1)
		value2 := trie2.GetKey(valueKey2)
		assert.Equal(t, v, value1)
		assert.Equal(t, v, value2)
	}

	for i := 0; i < 2; i++ {
		insert(trie1, trie2)
	}

	root1, _ = trie1.Commit(leafcallback(triedb1, trie1))
	root2, _ = trie2.Commit(leafcallback(triedb2, trie2))

	assert.Equal(t, root1, root2)
	triedb1.Commit(root1, true, true)
	triedb2.Commit(root2, true, true)
	for k, v := range storages {
		valueKey1 := trie1.Get([]byte(k))
		valueKey2 := trie2.Get([]byte(k))
		assert.Equal(t, valueKey1, valueKey2)
		value1 := trie1.GetKey(valueKey1)
		value2 := trie2.GetKey(valueKey2)
		assert.Equal(t, v, value1)
		assert.Equal(t, v, value2)
	}

}
func TestSecureDelete(t *testing.T) {
	trie := newEmptySecure()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"ether", ""},
		{"dog", "puppy"},
		{"shaman", ""},
	}
	for _, val := range vals {
		if val.v != "" {
			trie.Update([]byte(val.k), []byte(val.v))
		} else {
			trie.Delete([]byte(val.k))
		}
	}
	hash := trie.Hash()
	exp := common.HexToHash("29b235a58c3c25ab83010c327d5932bcf05324b7d6b1185e650798034783ca9d")
	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}

func TestSecureGetKey(t *testing.T) {
	trie := newEmptySecure()
	trie.Update([]byte("foo"), []byte("bar"))

	key := []byte("foo")
	value := []byte("bar")
	seckey := crypto.Keccak256(key)

	if !bytes.Equal(trie.Get(key), value) {
		t.Errorf("Get did not return bar")
	}
	if k := trie.GetKey(seckey); !bytes.Equal(k, key) {
		t.Errorf("GetKey returned %q, want %q", k, key)
	}
}

func TestSecureTrieConcurrency(t *testing.T) {
	// Create an initial trie and copy if for concurrent access
	_, trie, _ := makeTestSecureTrie()

	threads := runtime.NumCPU()
	tries := make([]*SecureTrie, threads)
	for i := 0; i < threads; i++ {
		cpy := *trie
		tries[i] = &cpy
	}
	// Start a batch of goroutines interactng with the trie
	pend := new(sync.WaitGroup)
	pend.Add(threads)
	for i := 0; i < threads; i++ {
		go func(index int) {
			defer pend.Done()

			for j := byte(0); j < 255; j++ {
				// Map the same data under multiple keys
				key, val := common.LeftPadBytes([]byte{byte(index), 1, j}, 32), []byte{j}
				tries[index].Update(key, val)

				key, val = common.LeftPadBytes([]byte{byte(index), 2, j}, 32), []byte{j}
				tries[index].Update(key, val)

				// Add some other data to inflate the trie
				for k := byte(3); k < 13; k++ {
					key, val = common.LeftPadBytes([]byte{byte(index), k, j}, 32), []byte{k, j}
					tries[index].Update(key, val)
				}
			}
			tries[index].Commit(nil)
		}(i)
	}
	// Wait for all threads to finish
	pend.Wait()
}

func TestJ(t *testing.T) {
	_, trie, _ := makeTestSecureTrie()
	fmt.Println(hexutil.Encode(trie.hashKey(hexutil.MustDecode("0x1000000000000000000000000000000000000004"))))
	fmt.Println(hexutil.Encode(trie.hashKey(hexutil.MustDecode("0x1000000000000000000000000000000000000002"))))

}
