// Copyright 2014 The go-ethereum Authors
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

package state

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/ethdb/leveldb"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/trie"

	checker "gopkg.in/check.v1"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
)

type StateSuite struct {
	db    ethdb.Database
	state *StateDB
}

var _ = checker.Suite(&StateSuite{})

var toAddr = common.BytesToAddress

func (s *StateSuite) TestDump(c *checker.C) {
	// generate a few entries
	obj1 := s.state.GetOrNewStateObject(toAddr([]byte{0x01}))
	obj1.AddBalance(big.NewInt(22))
	obj2 := s.state.GetOrNewStateObject(toAddr([]byte{0x01, 0x02}))
	obj2.SetCode(crypto.Keccak256Hash([]byte{3, 3, 3, 3, 3, 3, 3}), []byte{3, 3, 3, 3, 3, 3, 3})
	obj3 := s.state.GetOrNewStateObject(toAddr([]byte{0x02}))
	obj3.SetBalance(big.NewInt(44))

	// write some of them to the trie
	s.state.updateStateObject(obj1)
	s.state.updateStateObject(obj2)
	s.state.Commit(false)

	// check that dump contains the state objects that are in trie
	got := string(s.state.Dump())
	want := `{
    "root": "32d937466d6678befa41bcd94571dde0c612392ee2d2fa21a0d420b8f2b803bc",
    "accounts": {
        "0000000000000000000000000000000000000001": {
            "balance": "22",
            "nonce": 0,
            "root": "56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
            "codeHash": "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
            "code": "",
            "storage": {}
        },
        "0000000000000000000000000000000000000002": {
            "balance": "44",
            "nonce": 0,
            "root": "56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
            "codeHash": "c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
            "code": "",
            "storage": {}
        },
        "0000000000000000000000000000000000000102": {
            "balance": "0",
            "nonce": 0,
            "root": "56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
            "codeHash": "87874902497a5bb968da31a2998d8f22e949d1ef6214bcdedd8bae24cca4b9e3",
            "code": "03030303030303",
            "storage": {}
        }
    }
}`
	if got != want {
		c.Errorf("dump mismatch:\ngot: %s\nwant: %s\n", got, want)
	}
}

func (s *StateSuite) SetUpTest(c *checker.C) {
	s.db = rawdb.NewMemoryDatabase()
	s.state, _ = New(common.Hash{}, NewDatabase(s.db))
}

func (s *StateSuite) TestNull(c *checker.C) {
	address := common.MustBech32ToAddress("lax1qqqqqqyzx9q8zzl38xgwg5qpxeexmz64ex89tk")
	s.state.CreateAccount(address)
	value := common.FromHex("0x823140710bf13990e4500136726d8b55")
	//value := nil
	key := common.FromHex("0x823140710bf13990e4500136726d8b55")

	//s.state.SetState(address, common.Hash{}, value)
	s.state.SetState(address, key, value)
	s.state.Commit(false)

	if value := s.state.GetState(address, key); bytes.Compare(value, value) != 0 {
		c.Error("expected empty current value")
	}
}

func (s *StateSuite) TestSnapshot(c *checker.C) {
	stateobjaddr := toAddr([]byte("aa"))
	var storageaddr common.Hash
	data1 := common.BytesToHash([]byte{42})
	data2 := common.BytesToHash([]byte{43})

	// snapshot the genesis state
	genesis := s.state.Snapshot()

	// set initial state object value
	s.state.SetState(stateobjaddr, storageaddr.Bytes(), data1.Bytes())
	snapshot := s.state.Snapshot()

	// set a new state object value, revert it and ensure correct content
	s.state.SetState(stateobjaddr, storageaddr.Bytes(), data2.Bytes())
	s.state.RevertToSnapshot(snapshot)

	c.Assert(common.BytesToHash(s.state.GetState(stateobjaddr, storageaddr.Bytes())), checker.DeepEquals, data1)

	// revert up to the genesis state and ensure correct content
	s.state.RevertToSnapshot(genesis)
	c.Assert(common.BytesToHash(s.state.GetState(stateobjaddr, storageaddr.Bytes())), checker.DeepEquals, common.Hash{})
}

func (s *StateSuite) TestSnapshotEmpty(c *checker.C) {
	s.state.RevertToSnapshot(s.state.Snapshot())
}

// use testing instead of checker because checker does not support
// printing/logging in tests (-check.vv does not work)
func TestSnapshot2(t *testing.T) {
	state, _ := New(common.Hash{}, NewDatabase(rawdb.NewMemoryDatabase()))

	stateobjaddr0 := toAddr([]byte("so0"))
	stateobjaddr1 := toAddr([]byte("so1"))
	var storageaddr common.Address

	data0 := common.BytesToHash([]byte{17})
	data1 := common.BytesToHash([]byte{18})

	state.SetState(stateobjaddr0, storageaddr.Bytes(), data0.Bytes())
	state.SetState(stateobjaddr1, storageaddr.Bytes(), data1.Bytes())

	// db, trie are already non-empty values
	so0 := state.getStateObject(stateobjaddr0)
	so0.SetBalance(big.NewInt(42))
	so0.SetNonce(43)
	so0.SetCode(crypto.Keccak256Hash([]byte{'c', 'a', 'f', 'e'}), []byte{'c', 'a', 'f', 'e'})
	so0.suicided = false
	so0.deleted = false
	state.setStateObject(so0)

	root, _ := state.Commit(false)
	state.Reset(root)

	// and one with deleted == true
	so1 := state.getStateObject(stateobjaddr1)
	so1.SetBalance(big.NewInt(52))
	so1.SetNonce(53)
	so1.SetCode(crypto.Keccak256Hash([]byte{'c', 'a', 'f', 'e', '2'}), []byte{'c', 'a', 'f', 'e', '2'})
	so1.suicided = true
	so1.deleted = true
	state.setStateObject(so1)

	so1 = state.getStateObject(stateobjaddr1)
	if so1 != nil {
		t.Fatalf("deleted object not nil when getting")
	}

	snapshot := state.Snapshot()
	state.RevertToSnapshot(snapshot)

	so0Restored := state.getStateObject(stateobjaddr0)
	// Update lazily-loaded values before comparing.
	//key, _, _ := getKeyValue(stateobjaddr0, storageaddr.Bytes(), nil)
	key := storageaddr.Bytes()
	so0Restored.GetState(state.db, key)
	so0Restored.Code(state.db)
	// non-deleted is equal (restored)
	compareStateObjects(so0Restored, so0, t)

	// deleted should be nil, both before and after restore of state copy
	so1Restored := state.getStateObject(stateobjaddr1)
	if so1Restored != nil {
		t.Fatalf("deleted object not nil after restoring snapshot: %+v", so1Restored)
	}
}

func compareStateObjects(so0, so1 *stateObject, t *testing.T) {
	if so0.Address() != so1.Address() {
		t.Fatalf("Address mismatch: have %v, want %v", so0.address, so1.address)
	}
	if so0.Balance().Cmp(so1.Balance()) != 0 {
		t.Fatalf("Balance mismatch: have %v, want %v", so0.Balance(), so1.Balance())
	}
	if so0.Nonce() != so1.Nonce() {
		t.Fatalf("Nonce mismatch: have %v, want %v", so0.Nonce(), so1.Nonce())
	}
	if so0.data.Root != so1.data.Root {
		t.Errorf("Root mismatch: have %x, want %x", so0.data.Root[:], so1.data.Root[:])
	}
	if !bytes.Equal(so0.CodeHash(), so1.CodeHash()) {
		t.Fatalf("CodeHash mismatch: have %v, want %v", so0.CodeHash(), so1.CodeHash())
	}
	if !bytes.Equal(so0.code, so1.code) {
		t.Fatalf("Code mismatch: have %v, want %v", so0.code, so1.code)
	}

	if len(so1.dirtyStorage) != len(so0.dirtyStorage) {
		t.Errorf("Dirty storage size mismatch: have %d, want %d", len(so1.dirtyStorage), len(so0.dirtyStorage))
	}
	for k, v := range so1.dirtyStorage {
		if !bytes.Equal(so0.dirtyStorage[k], v) {
			t.Errorf("Dirty storage key %x mismatch: have %v, want %v", k, so0.dirtyStorage[k], v)
		}
	}
	for k, v := range so0.dirtyStorage {
		if !bytes.Equal(so1.dirtyStorage[k], v) {
			t.Errorf("Dirty storage key %x mismatch: have %v, want none.", k, v)
		}
	}
	if len(so1.originStorage) != len(so0.originStorage) {
		t.Errorf("Origin storage size mismatch: have %d, want %d", len(so1.originStorage), len(so0.originStorage))
	}
	for k, v := range so1.originStorage {
		if !bytes.Equal(so0.originStorage[k], v) {
			t.Errorf("Origin storage key %x mismatch: have %v, want %v", k, so0.originStorage[k], v)
		}
	}
	for k, v := range so0.originStorage {
		if !bytes.Equal(so1.originStorage[k], v) {
			t.Errorf("Origin storage key %x mismatch: have %v, want none.", k, v)
		}
	}
}

func TestEmptyByte(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "platon")
	defer os.Remove(tmpDir)

	db, _ := leveldb.New(tmpDir, 0, 0, "")
	state, _ := New(common.Hash{}, NewDatabase(db))

	address := common.MustBech32ToAddress("lax1qqqqqqyzx9q8zzl38xgwg5qpxeexmz64ex89tk")
	state.CreateAccount(address)
	so := state.getStateObject(address)

	//value := common.FromHex("0x823140710bf13990e4500136726d8b55")
	pvalue := []byte{'a'}
	key := []byte{'a'}

	//s.state.SetState(address, common.Hash{}, value)
	state.SetState(address, key, pvalue)
	state.Commit(false)

	if value := state.GetState(address, key); !bytes.Equal(value, pvalue) {
		t.Errorf("expected empty current value, got %x", value)
	}
	if value := state.GetCommittedState(address, key); !bytes.Equal(value, pvalue) {
		t.Errorf("expected empty committed value, got %x", value)
	}

	state.trie.NodeIterator(nil)
	it := trie.NewIterator(so.trie.NodeIterator(nil))
	for it.Next() {
		fmt.Println(it.Key, it.Value)
	}

	pvalue = []byte{}
	state.SetState(address, key, []byte{})
	state.Commit(false)

	if value := state.GetState(address, key); !bytes.Equal(value, pvalue) {
		t.Errorf("expected empty current value, got %x", value)
	}
	if value := state.GetCommittedState(address, key); !bytes.Equal(value, pvalue) {
		t.Errorf("expected empty committed value, got %x", value)
	}

	state.trie.NodeIterator(nil)
	it = trie.NewIterator(so.trie.NodeIterator(nil))
	for it.Next() {
		fmt.Println(it.Key, it.Value)
	}

	pvalue = []byte("bbb")
	state.SetState(address, key, pvalue)
	state.Commit(false)
	state.trie.NodeIterator(nil)
	it = trie.NewIterator(so.trie.NodeIterator(nil))
	for it.Next() {
		fmt.Println(it.Key, it.Value)
		fmt.Println(so.db.trie.GetKey(it.Value))
	}

}

func TestForEachStorage(t *testing.T) {
	tmpDir, _ := ioutil.TempDir("", "platon")
	defer os.Remove(tmpDir)
	db, _ := leveldb.New(tmpDir, 0, 0, "")
	state, _ := New(common.Hash{}, NewDatabase(db))

	address := common.MustBech32ToAddress("lax1qqqqqqyzx9q8zzl38xgwg5qpxeexmz64ex89tk")
	state.CreateAccount(address)

	key := []byte("a")
	fvalue := []byte("A")

	fmt.Printf("before Commit, key: %v, value: %v \n", key, fvalue)

	//s.state.SetState(address, common.Hash{}, value)
	state.SetState(address, key, fvalue)
	state.Commit(false)

	svalue := []byte("B")

	fmt.Printf("after Commit, key: %v, value: %v \n", key, svalue)
	state.SetState(address, key, svalue)

	state.ForEachStorage(address, func(key []byte, value []byte) bool {
		fmt.Println("load out, key:", hex.EncodeToString(key), "value:", string(value))
		fmt.Printf("load out, key: %v, value: %v \n", key, value /*Bytes2Bits(key), Bytes2Bits(value)*/)
		return true
	})
}

func TestMigrateStorage(t *testing.T) {

	tmpDir, _ := ioutil.TempDir("", "platon")
	defer os.Remove(tmpDir)
	db, _ := leveldb.New(tmpDir, 0, 0, "")
	state, _ := New(common.Hash{}, NewDatabase(db))

	from := common.MustBech32ToAddress("lax1qqqqqqyzx9q8zzl38xgwg5qpxeexmz64ex89tk")
	state.CreateAccount(from)

	to := common.MustBech32ToAddress("lax1qqqqqqrjxpq8zzl38xgwg5qpxeex6mnxwyzlxv")
	state.CreateAccount(to)

	state.SetState(from, []byte("a"), []byte("fromA"))
	state.SetState(from, []byte("b"), []byte("fromB"))
	state.SetState(from, []byte("c"), []byte("fromC"))

	state.SetState(to, []byte("a"), []byte("I am  A of to"))
	state.SetState(to, []byte("b"), []byte("I am  B of to"))
	state.SetState(to, []byte("d"), []byte("I am  D of to"))
	state.SetState(to, []byte("e"), []byte("I am  E of to"))

	state.Commit(false)

	state.SetState(from, []byte("c"), []byte("fromC2"))
	state.SetState(from, []byte("d"), []byte("fromD2"))
	state.SetState(to, []byte("e"), []byte("I am  E2 of to"))
	state.SetState(to, []byte("f"), []byte("I am  F of to"))

	// test MigrateStorage
	//
	// expect:
	//
	// {
	//		"a": "fromA",
	//		"b": "fromB",
	//		"c": "fromC2",
	// 		"d": "fromD2",
	//		"e": "",
	// 		"f": "",
	// }
	//
	state.MigrateStorage(from, to)

	for _, key := range [][]byte{[]byte("a"), []byte("b"), []byte("c"), []byte("d"), []byte("e"), []byte("f")} {
		value := state.GetState(to, key)

		switch string(key) {
		case "a":
			assert.Equal(t, "fromA", string(value))
		case "b":
			assert.Equal(t, "fromB", string(value))
		case "c":
			assert.Equal(t, "fromC2", string(value))
		case "d":
			assert.Equal(t, "fromD2", string(value))
		case "e":
			assert.Equal(t, "", string(value))
		case "f":
			assert.Equal(t, "", string(value))
		}

		//fmt.Println("key:", string(key), "value:", string(value))
	}
}

func Bytes2Bits(data []byte) []int {
	dst := make([]int, 0)
	for _, v := range data {
		for i := 0; i < 8; i++ {
			move := uint(7 - i)
			dst = append(dst, int((v>>move)&1))
		}
	}
	return dst
}
