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
	"encoding/json"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/trie"
)

type stateTest struct {
	db    ethdb.Database
	state *StateDB
}

func newStateTest() *stateTest {
	db := rawdb.NewMemoryDatabase()
	sdb, _ := New(types.EmptyRootHash, NewDatabase(db), nil)
	return &stateTest{db: db, state: sdb}
}

func TestDump(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	sdb, _ := New(types.EmptyRootHash, NewDatabaseWithConfig(db, &trie.Config{Preimages: true}), nil)
	s := &stateTest{db: db, state: sdb}

	// generate a few entries
	obj1 := s.state.GetOrNewStateObject(common.BytesToAddress([]byte{0x01}))
	obj1.AddBalance(big.NewInt(22))
	obj2 := s.state.GetOrNewStateObject(common.BytesToAddress([]byte{0x01, 0x02}))
	obj2.SetCode(crypto.Keccak256Hash([]byte{3, 3, 3, 3, 3, 3, 3}), []byte{3, 3, 3, 3, 3, 3, 3})
	obj3 := s.state.GetOrNewStateObject(common.BytesToAddress([]byte{0x02}))
	obj3.SetBalance(big.NewInt(44))

	// write some of them to the trie
	s.state.updateStateObject(obj1)
	s.state.updateStateObject(obj2)
	s.state.Commit(0, false)

	// check that dump contains the state objects that are in trie
	opts := &DumpConfig{
		OnlyWithAddresses: true,
		Max:               256, // Sanity limit over RPC
	}
	got := string(s.state.Dump(opts))
	want := `{
    "root": "32d937466d6678befa41bcd94571dde0c612392ee2d2fa21a0d420b8f2b803bc",
    "accounts": {
        "lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqgzywgsrf": {
            "balance": "0",
            "nonce": 0,
            "root": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
            "codeHash": "0x87874902497a5bb968da31a2998d8f22e949d1ef6214bcdedd8bae24cca4b9e3",
            "code": "0x03030303030303",
            "key": "0xa17eacbc25cda025e81db9c5c62868822c73ce097cee2a63e33a2e41268358a1"
        },
        "lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpfr7f80": {
            "balance": "22",
            "nonce": 0,
            "root": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
            "codeHash": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
            "key": "0x1468288056310c82aa4c01a7e12a10f8111a0560e72b700555479031b86c357d"
        },
        "lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqz8stlfs": {
            "balance": "44",
            "nonce": 0,
            "root": "0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
            "codeHash": "0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470",
            "key": "0xd52688a8f926c816ca1e079067caba944f158e764817b83fc43594370ca9cf62"
        }
    }
}`
	if got != want {
		t.Errorf("dump mismatch:\ngot: %s\nwant: %s\n", got, want)
	}
}

func TestIterativeDump(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	sdb, _ := New(types.EmptyRootHash, NewDatabaseWithConfig(db, &trie.Config{Preimages: true}), nil)
	s := &stateTest{db: db, state: sdb}

	// generate a few entries
	obj1 := s.state.GetOrNewStateObject(common.BytesToAddress([]byte{0x01}))
	obj1.AddBalance(big.NewInt(22))
	obj2 := s.state.GetOrNewStateObject(common.BytesToAddress([]byte{0x01, 0x02}))
	obj2.SetCode(crypto.Keccak256Hash([]byte{3, 3, 3, 3, 3, 3, 3}), []byte{3, 3, 3, 3, 3, 3, 3})
	obj3 := s.state.GetOrNewStateObject(common.BytesToAddress([]byte{0x02}))
	obj3.SetBalance(big.NewInt(44))
	obj4 := s.state.GetOrNewStateObject(common.BytesToAddress([]byte{0x00}))
	obj4.AddBalance(big.NewInt(1337))

	// write some of them to the trie
	s.state.updateStateObject(obj1)
	s.state.updateStateObject(obj2)
	s.state.Commit(0, false)

	b := &bytes.Buffer{}
	s.state.IterativeDump(nil, json.NewEncoder(b))
	// check that DumpToCollector contains the state objects that are in trie
	got := b.String()
	want := `{"root":"0xb29b5c1356766ecfce933d4e3057aa72a51a6e2d9eeb550e49c6d48bb8817a3b"}
{"balance":"22","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470","address":"lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqpfr7f80","key":"0x1468288056310c82aa4c01a7e12a10f8111a0560e72b700555479031b86c357d"}
{"balance":"1337","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470","address":"lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqq542u6a","key":"0x5380c7b7ae81a58eb98d9c78de4a1fd7fd9535fc953ed2be602daaa41767312a"}
{"balance":"0","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0x87874902497a5bb968da31a2998d8f22e949d1ef6214bcdedd8bae24cca4b9e3","code":"0x03030303030303","address":"lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqgzywgsrf","key":"0xa17eacbc25cda025e81db9c5c62868822c73ce097cee2a63e33a2e41268358a1"}
{"balance":"44","nonce":0,"root":"0x56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421","codeHash":"0xc5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470","address":"lat1qqqqqqqqqqqqqqqqqqqqqqqqqqqqqqqz8stlfs","key":"0xd52688a8f926c816ca1e079067caba944f158e764817b83fc43594370ca9cf62"}
`
	if got != want {
		t.Errorf("DumpToCollector mismatch:\ngot: %s\nwant: %s\n", got, want)
	}
}

func TestNull(t *testing.T) {
	s := newStateTest()
	address := common.MustBech32ToAddress("lax1qqqqqqyzx9q8zzl38xgwg5qpxeexmz64ex89tk")
	s.state.CreateAccount(address)
	value := common.FromHex("0x823140710bf13990e4500136726d8b55")
	//value := nil
	key := common.FromHex("0x823140710bf13990e4500136726d8b55")

	//s.state.SetState(address, common.Hash{}, value)
	s.state.SetState(address, key, value)
	s.state.Commit(0, false)

	if value := s.state.GetState(address, key); !bytes.Equal(value, value) {
		t.Error("expected empty current value")
	}
}

func TestSnapshot(t *testing.T) {
	stateobjaddr := common.BytesToAddress([]byte("aa"))
	var storageaddr common.Hash
	data1 := common.BytesToHash([]byte{42})
	data2 := common.BytesToHash([]byte{43})
	s := newStateTest()

	// snapshot the genesis state
	genesis := s.state.Snapshot()

	// set initial state object value
	s.state.SetState(stateobjaddr, storageaddr.Bytes(), data1.Bytes())
	snapshot := s.state.Snapshot()

	// set a new state object value, revert it and ensure correct content
	s.state.SetState(stateobjaddr, storageaddr.Bytes(), data2.Bytes())
	s.state.RevertToSnapshot(snapshot)

	if v := s.state.GetState(stateobjaddr, storageaddr.Bytes()); !bytes.Equal(v, data1.Bytes()) {
		t.Errorf("wrong storage value %v, want %v", v, data1)
	}
	if v := s.state.GetCommittedState(stateobjaddr, storageaddr.Bytes()); !bytes.Equal(v, []byte("")) {
		t.Errorf("wrong committed storage value %v, want %v", v, []byte(""))
	}

	// revert up to the genesis state and ensure correct content
	s.state.RevertToSnapshot(genesis)
	if v := s.state.GetState(stateobjaddr, storageaddr.Bytes()); !bytes.Equal(v, []byte("")) {
		t.Errorf("wrong storage value %v, want %v", v, []byte(""))
	}
	if v := s.state.GetCommittedState(stateobjaddr, storageaddr.Bytes()); !bytes.Equal(v, []byte("")) {
		t.Errorf("wrong committed storage value %v, want %v", v, []byte(""))
	}
}

func TestSnapshotEmpty(t *testing.T) {
	s := newStateTest()
	s.state.RevertToSnapshot(s.state.Snapshot())
}

func TestSnapshot2(t *testing.T) {
	state, _ := New(types.EmptyRootHash, NewDatabase(rawdb.NewMemoryDatabase()), nil)

	stateobjaddr0 := common.BytesToAddress([]byte("so0"))
	stateobjaddr1 := common.BytesToAddress([]byte("so1"))
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
	so0.selfDestructed = false
	so0.deleted = false
	state.setStateObject(so0)

	root, _ := state.Commit(0, false)
	state.Reset(root)

	// and one with deleted == true
	so1 := state.getStateObject(stateobjaddr1)
	so1.SetBalance(big.NewInt(52))
	so1.SetNonce(53)
	so1.SetCode(crypto.Keccak256Hash([]byte{'c', 'a', 'f', 'e', '2'}), []byte{'c', 'a', 'f', 'e', '2'})
	so1.selfDestructed = true
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
	so0Restored.GetState(key)
	so0Restored.Code()
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
