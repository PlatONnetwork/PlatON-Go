// Copyright 2016 The go-ethereum Authors
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
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"reflect"
	"strings"
	"testing"
	"testing/quick"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/stretchr/testify/assert"

	"gopkg.in/check.v1"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// Tests that updating a state trie does not leak any database writes prior to
// actually committing the state.
func TestUpdateLeaks(t *testing.T) {
	// Create an empty state database
	db := rawdb.NewMemoryDatabase()
	//dir, _ := ioutil.TempDir("", "eth-core-bench")
	//ethdb,err:= ethdb.NewLDBDatabase(dir,128,128)
	state, _ := New(common.Hash{}, NewDatabase(db))

	// Update it with some accounts
	for i := byte(0); i < 255; i++ {
		addr := common.BytesToAddress([]byte{i})
		state.AddBalance(addr, big.NewInt(int64(11*i)))
		state.SetNonce(addr, uint64(42*i))
		if i%2 == 0 {
			state.SetState(addr, []byte{i, i, i}, []byte{i, i, i, i})
		}
		if i%3 == 0 {
			state.SetCode(addr, []byte{i, i, i, i, i})
		}
		state.IntermediateRoot(false)
	}
	// Ensure that no data was leaked into the database
	it := db.NewIterator()
	for it.Next() {
		t.Errorf("State leaked into database: %x -> %x", it.Key(), it.Value())
	}
	it.Release()
}

func TestClearParentReference(t *testing.T) {
	db := rawdb.NewMemoryDatabase()
	state, _ := New(common.Hash{}, NewDatabase(db))
	count := 10000

	for i := 0; i < count; i++ {
		st := state.NewStateDB()
		if i%2 == 0 {
			go func() {
				st.ClearParentReference()
			}()
		}
	}
	time.Sleep(2 * time.Second)
	assert.Len(t, state.clearReferenceFunc, count)

	clear := 0
	for i := 0; i < count; i++ {
		if i%2 == 0 {
			fn := state.clearReferenceFunc[i]
			if fn == nil {
				clear++
			}
		}
	}
	assert.Equal(t, clear, count/2)
}

func TestNewStateDBAndCopy(t *testing.T) {
	storages := make(map[common.Address]map[string]string)
	db := rawdb.NewMemoryDatabase()
	dbc := rawdb.NewMemoryDatabase()

	s1, _ := New(common.Hash{}, NewDatabase(db))
	s1c, _ := New(common.Hash{}, NewDatabase(dbc))

	modify := func(s1 *StateDB, s2 *StateDB, addr common.Address, i int) {
		s1.AddBalance(addr, big.NewInt(int64(i)))
		s2.AddBalance(addr, big.NewInt(int64(i)))

		s1.SetNonce(addr, uint64(i*4))
		s2.SetNonce(addr, uint64(i*4))
		s1.SetCode(addr, []byte{byte(i), byte(i)})
		s2.SetCode(addr, []byte{byte(i), byte(i)})

		key := randString(i)
		value := randString(i)
		s1.SetState(addr, []byte(key), []byte(value))
		s2.SetState(addr, []byte(key), []byte(value))
		if k, ok := storages[addr]; ok {
			k[key] = value
		} else {
			maps := make(map[string]string)
			maps[key] = value
			storages[addr] = maps
		}
		//if bytes.Equal([]byte{}, []byte(value)) || bytes.Equal([]byte(nil), []byte(value)) {
		//	fmt.Println("xaxaxaxaxaxaxaxa", "addr", addr.String(),  "key", key, "value", value, "[]byte{}",  bytes.Equal([]byte{}, []byte(value)), "[]byte(nil)",  bytes.Equal([]byte(nil), []byte(value)))
		//	fmt.Println("xaxaxaxaxaxaxaxa", bytes.Equal([]byte{}, []byte(nil)))
		//}
		//fmt.Println("addr", addr.String(), "key", key, "value", value)
	}
	for i := 1; i < 255; i++ {
		modify(s1, s1c, common.Address{byte(i)}, i)
	}

	st := s1.NewStateDB()
	for addr, storage := range storages {
		for k, v := range storage {
			value := st.GetState(addr, []byte(k))
			//fmt.Println("storage :::: ", "k", k, "v", v, "state v", string(value))
			//if bytes.Equal([]byte{}, []byte(v)) || bytes.Equal(value, []byte(nil)) {
			//	fmt.Println("lalalalala", "k", k, "v", v, "state v", string(value))
			//}
			assert.Equal(t, []byte(v), value)
		}
	}

	if _, err := s1.Commit(false); err != nil {
		t.Fatalf("failed to commit s1 state: %v", err)
	}
	if _, err := s1c.Commit(false); err != nil {
		t.Fatalf("failed to commit s1c state: %v", err)
	}
	assert.Nil(t, s1.db.TrieDB().Commit(s1.Root(), false, true))
	assert.Nil(t, s1c.db.TrieDB().Commit(s1c.Root(), false, true))

	// test new statedb
	st2 := s1.NewStateDB()
	s1db, _ := New(s1.Root(), NewDatabase(dbc))

	st2.clearParentRef()
	for addr, storage := range storages {
		for k, v := range storage {
			value := st2.GetState(addr, []byte(k))
			value2 := s1c.GetState(addr, []byte(k))
			value3 := s1db.GetState(addr, []byte(k))
			//fmt.Println("v", hex.EncodeToString([]byte(v)), "value", hex.EncodeToString([]byte(value)), "value2", hex.EncodeToString(value2))
			assert.Equal(t, []byte(v), value)
			assert.Equal(t, []byte(v), value2)
			assert.Equal(t, []byte(v), value3)
		}
	}
	it := db.NewIterator()
	for it.Next() {
		if _, err := dbc.Get(it.Key()); err != nil {
			v, _ := db.Get(it.Key())
			t.Fatalf("db get error, key:%s, value:%s", hex.EncodeToString(it.Key()), hex.EncodeToString(v))
		}
	}
	it.Release()

	// s3->s2->s1, s1c is copy of s1. insert random kv, db is the same as dbc
	s2 := s1.NewStateDB()
	s3 := s2.NewStateDB()

	assert.Len(t, s1.clearReferenceFunc, 3)
	assert.Len(t, s2.clearReferenceFunc, 1)
	assert.Len(t, s3.clearReferenceFunc, 0)

	assert.NotNil(t, s2.parent)
	assert.False(t, s2.parentCommitted)
	assert.NotNil(t, s3.parent)
	assert.False(t, s3.parentCommitted)

	s1.ClearReference()
	s2.ClearReference()
	assert.Len(t, s1.clearReferenceFunc, 0)
	assert.Len(t, s2.clearReferenceFunc, 0)
	assert.Len(t, s3.clearReferenceFunc, 0)

	assert.Nil(t, s2.parent)
	assert.True(t, s2.parentCommitted)
	assert.Nil(t, s3.parent)
	assert.True(t, s3.parentCommitted)

	s1cc := s1c.Copy()
	assert.Len(t, s1c.clearReferenceFunc, 0)
	assert.Len(t, s1cc.clearReferenceFunc, 0)
	assert.Nil(t, s1cc.parent)

	s1c.ClearReference()

	assert.Len(t, s1c.clearReferenceFunc, 0)
	assert.Len(t, s1cc.clearReferenceFunc, 0)
	assert.Nil(t, s1cc.parent)

	for i := 1; i < 255; i++ {
		modify(s3, s1cc, common.Address{byte(i)}, i)
	}

	if _, err := s3.Commit(false); err != nil {
		t.Fatalf("failed to commit s1 state: %v", err)
	}
	if _, err := s1cc.Commit(false); err != nil {
		t.Fatalf("failed to commit s1c state: %v", err)
	}

	// test copy statedb
	st4 := s3.Copy()

	st4.clearParentRef()
	for addr, storage := range storages {
		for k, v := range storage {
			value := st4.GetState(addr, []byte(k))
			value2 := s1cc.GetState(addr, []byte(k))
			//fmt.Println("v", hex.EncodeToString([]byte(v)), "value", hex.EncodeToString([]byte(value)), "value2", hex.EncodeToString(value2))
			assert.Equal(t, []byte(v), value)
			assert.Equal(t, []byte(v), value2)
		}
	}

	assert.Nil(t, s3.db.TrieDB().Commit(s1.Root(), false, true))
	assert.Nil(t, s1cc.db.TrieDB().Commit(s1cc.Root(), false, true))

	it2 := db.NewIterator()
	for it2.Next() {
		if _, err := dbc.Get(it2.Key()); err != nil {
			v, _ := db.Get(it2.Key())
			t.Fatalf("db get error, key:%s, value:%s", hex.EncodeToString(it2.Key()), hex.EncodeToString(v))
		}
	}
	it2.Release()

}

func TestStateStorageValueCommit(t *testing.T) {
	storages := make(map[common.Address]map[string]string)
	db := rawdb.NewMemoryDatabase()

	s1, _ := New(common.Hash{}, NewDatabase(db))

	modify := func(s1 *StateDB, addr common.Address, i int) {
		s1.AddBalance(addr, big.NewInt(int64(i)))

		s1.SetNonce(addr, uint64(i*4))
		s1.SetCode(addr, []byte{byte(9), byte(9)})

		key := randString(i + 20)
		value := randString(i + 20)
		s1.SetState(addr, []byte(key), []byte(value))
		if k, ok := storages[addr]; ok {
			k[key] = value
		} else {
			maps := make(map[string]string)
			maps[key] = value
			storages[addr] = maps
		}
	}

	for i := 0; i < 255; i++ {
		modify(s1, common.Address{byte(i)}, i)
	}

	root, err := s1.Commit(true)
	if err != nil {
		t.Fatal(err)
	}
	if err := s1.db.TrieDB().Commit(root, true, true); err != nil {
		t.Fatal(err)
	}
	s2, err := New(root, NewDatabase(db))
	for addr, storage := range storages {
		for key, value := range storage {
			exp := s2.GetState(addr, []byte(key))
			assert.Equal(t, []byte(value[:]), exp)
		}
	}
}

func TestStateStorageValueDelete(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	s1, _ := New(common.Hash{}, NewDatabase(db))

	key1, value1, key2, value2 := []byte("key1"), []byte("value1"), []byte("key2"), []byte("value2")

	addr := common.Address{byte(1)}
	s1.SetCode(addr, []byte{byte(9), byte(9)})

	s1.SetState(addr, key1, value1)
	s1.SetState(addr, key2, value2)

	s1.SetState(addr, key2[:], []byte{})

	root, err := s1.Commit(true)
	if err != nil {
		t.Fatal(err)
	}
	if err := s1.db.TrieDB().Commit(root, true, true); err != nil {
		t.Fatal(err)
	}

	s2, err := New(root, NewDatabase(db))
	exp := s2.GetState(addr, key1)
	assert.Equal(t, exp, value1)

	exp = s2.GetState(addr, key2)
	assert.NotEqual(t, exp, value2)

}

func TestStateStorageRevert(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	s1, _ := New(common.Hash{}, NewDatabase(db))
	s1.Snapshot()
	key1, value1, value2 := []byte("key1"), []byte("value1"), []byte("value2")

	// insert value
	addr := common.Address{byte(1)}
	s1.SetCode(addr, []byte{byte(9), byte(9)})
	s1.SetState(addr, key1, value1)
	s1.IntermediateRoot(false)
	assert.Equal(t, value1, s1.GetState(addr, key1))

	storage := s1.Snapshot()

	// twice
	s1.SetState(addr, key1, value2)
	assert.Equal(t, value2, s1.GetState(addr, key1))

	// revert
	s1.RevertToSnapshot(storage)
	assert.Equal(t, value1, s1.GetState(addr, key1))

	root, err := s1.Commit(true)
	if err != nil {
		t.Fatal(err)
	}

	s1.db.TrieDB().Commit(root, true, true)

	assert.Equal(t, value1, s1.GetState(addr, key1))
	obj := s1.getStateObject(addr)
	assert.Equal(t, 0, len(obj.dirtyStorage))
}

func TestStateStorageValueUpdate(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	s1, _ := New(common.Hash{}, NewDatabase(db))
	s1.Snapshot()
	key1, value1, value2 := []byte("key1"), []byte("value1"), []byte("value2")

	// insert value
	addr := common.Address{byte(1)}
	s1.SetCode(addr, []byte{byte(9), byte(9)})
	s1.SetState(addr, key1, value1)
	s1.SetState(addr, key1, value2)
	obj := s1.getStateObject(addr)
	assert.Equal(t, 1, len(obj.dirtyStorage))
}

// Tests that no intermediate state of an object is stored into the database,
// only the one right before the commit.
func TestIntermediateLeaks(t *testing.T) {
	// Create two state databases, one transitioning to the final state, the other final from the beginning
	transDb := rawdb.NewMemoryDatabase()
	finalDb := rawdb.NewMemoryDatabase()
	transState, _ := New(common.Hash{}, NewDatabase(transDb))
	finalState, _ := New(common.Hash{}, NewDatabase(finalDb))

	modify := func(state *StateDB, addr common.Address, i, tweak byte) {
		state.SetBalance(addr, big.NewInt(int64(11*i)+int64(tweak)))
		state.SetNonce(addr, uint64(42*i+tweak))
		if i%2 == 0 {
			//state.SetState(addr, common.Hash{i, i, i, 0}, common.Hash{})
			//state.SetState(addr, common.Hash{i, i, i, tweak}, common.Hash{i, i, i, i, tweak})
		}
		if i%3 == 0 {
			state.SetCode(addr, []byte{i, i, i, i, i, tweak})
		}
	}

	// Modify the transient state.
	for i := byte(0); i < 255; i++ {
		modify(transState, common.Address{byte(i)}, i, 0)
	}
	// Write modifications to trie.
	transState.IntermediateRoot(false)

	// Overwrite all the data with new values in the transient database.
	for i := byte(0); i < 255; i++ {
		modify(transState, common.Address{byte(i)}, i, 99)
		modify(finalState, common.Address{byte(i)}, i, 99)
	}

	// Commit and cross check the databases.
	if _, err := transState.Commit(false); err != nil {
		t.Fatalf("failed to commit transition state: %v", err)
	}
	if _, err := finalState.Commit(false); err != nil {
		t.Fatalf("failed to commit final state: %v", err)
	}
	it := finalDb.NewIterator()
	for it.Next() {
		key := it.Key()
		if _, err := transDb.Get(key); err != nil {
			t.Errorf("entry missing from the transition database: %x -> %x", key, it.Value())
		}
	}
	it.Release()

	it = transDb.NewIterator()
	for it.Next() {
		key := it.Key()
		if _, err := finalDb.Get(key); err != nil {
			t.Errorf("extra entry in the transition database: %x -> %x", key, it.Value())
		}
	}
}

// TestCopy tests that copying a statedb object indeed makes the original and
// the copy independent of each other. This test is a regression test against
// https://github.com/ethereum/go-ethereum/pull/15549.
func TestCopy(t *testing.T) {
	// Create a random state test to copy and modify "independently"
	orig, _ := New(common.Hash{}, NewDatabase(rawdb.NewMemoryDatabase()))

	for i := byte(0); i < 255; i++ {
		obj := orig.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		obj.AddBalance(big.NewInt(int64(i)))
		orig.updateStateObject(obj)
	}
	orig.Finalise(false)

	// Copy the state, modify both in-memory
	copy := orig.Copy()

	for i := byte(0); i < 255; i++ {
		origObj := orig.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		copyObj := copy.GetOrNewStateObject(common.BytesToAddress([]byte{i}))

		origObj.AddBalance(big.NewInt(2 * int64(i)))
		copyObj.AddBalance(big.NewInt(3 * int64(i)))

		orig.updateStateObject(origObj)
		copy.updateStateObject(copyObj)
	}
	// Finalise the changes on both concurrently
	done := make(chan struct{})
	go func() {
		orig.Finalise(true)
		close(done)
	}()
	copy.Finalise(true)
	<-done

	// Verify that the two states have been updated independently
	for i := byte(0); i < 255; i++ {
		origObj := orig.GetOrNewStateObject(common.BytesToAddress([]byte{i}))
		copyObj := copy.GetOrNewStateObject(common.BytesToAddress([]byte{i}))

		if want := big.NewInt(3 * int64(i)); origObj.Balance().Cmp(want) != 0 {
			t.Errorf("orig obj %d: balance mismatch: have %v, want %v", i, origObj.Balance(), want)
		}
		if want := big.NewInt(4 * int64(i)); copyObj.Balance().Cmp(want) != 0 {
			t.Errorf("copy obj %d: balance mismatch: have %v, want %v", i, copyObj.Balance(), want)
		}
	}
}

func TestSnapshotRandom(t *testing.T) {
	config := &quick.Config{MaxCount: 1000}
	err := quick.Check((*snapshotTest).run, config)
	if cerr, ok := err.(*quick.CheckError); ok {
		test := cerr.In[0].(*snapshotTest)
		t.Errorf("%v:\n%s", test.err, test)
	} else if err != nil {
		t.Error(err)
	}
}

// A snapshotTest checks that reverting StateDB snapshots properly undoes all changes
// captured by the snapshot. Instances of this test with pseudorandom content are created
// by Generate.
//
// The test works as follows:
//
// A new state is created and all actions are applied to it. Several snapshots are taken
// in between actions. The test then reverts each snapshot. For each snapshot the actions
// leading up to it are replayed on a fresh, empty state. The behaviour of all public
// accessor methods on the reverted state must match the return value of the equivalent
// methods on the replayed state.
type snapshotTest struct {
	addrs     []common.Address // all account addresses
	actions   []testAction     // modifications to the state
	snapshots []int            // actions indexes at which snapshot is taken
	err       error            // failure details are reported through this field
}

type testAction struct {
	name   string
	fn     func(testAction, *StateDB)
	args   []int64
	noAddr bool
}

// newTestAction creates a random action that changes state.
func newTestAction(addr common.Address, r *rand.Rand) testAction {
	actions := []testAction{
		{
			name: "SetBalance",
			fn: func(a testAction, s *StateDB) {
				s.SetBalance(addr, big.NewInt(a.args[0]))
			},
			args: make([]int64, 1),
		},
		{
			name: "AddBalance",
			fn: func(a testAction, s *StateDB) {
				s.AddBalance(addr, big.NewInt(a.args[0]))
			},
			args: make([]int64, 1),
		},
		{
			name: "SetNonce",
			fn: func(a testAction, s *StateDB) {
				s.SetNonce(addr, uint64(a.args[0]))
			},
			args: make([]int64, 1),
		},
		{
			name: "SetState",
			fn: func(a testAction, s *StateDB) {
				var key, val common.Hash
				binary.BigEndian.PutUint16(key[:], uint16(a.args[0]))
				binary.BigEndian.PutUint16(val[:], uint16(a.args[1]))
				//s.SetState(addr, key, val)
			},
			args: make([]int64, 2),
		},
		{
			name: "SetCode",
			fn: func(a testAction, s *StateDB) {
				code := make([]byte, 16)
				binary.BigEndian.PutUint64(code, uint64(a.args[0]))
				binary.BigEndian.PutUint64(code[8:], uint64(a.args[1]))
				s.SetCode(addr, code)
			},
			args: make([]int64, 2),
		},
		{
			name: "CreateAccount",
			fn: func(a testAction, s *StateDB) {
				s.CreateAccount(addr)
			},
		},
		{
			name: "Suicide",
			fn: func(a testAction, s *StateDB) {
				s.Suicide(addr)
			},
		},
		{
			name: "AddRefund",
			fn: func(a testAction, s *StateDB) {
				s.AddRefund(uint64(a.args[0]))
			},
			args:   make([]int64, 1),
			noAddr: true,
		},
		{
			name: "AddLog",
			fn: func(a testAction, s *StateDB) {
				data := make([]byte, 2)
				binary.BigEndian.PutUint16(data, uint16(a.args[0]))
				s.AddLog(&types.Log{Address: addr, Data: data})
			},
			args: make([]int64, 1),
		},
	}
	action := actions[r.Intn(len(actions))]
	var nameargs []string
	if !action.noAddr {
		nameargs = append(nameargs, addr.String())
	}
	for _, i := range action.args {
		action.args[i] = rand.Int63n(100)
		nameargs = append(nameargs, fmt.Sprint(action.args[i]))
	}
	action.name += strings.Join(nameargs, ", ")
	return action
}

// Generate returns a new snapshot test of the given size. All randomness is
// derived from r.
func (*snapshotTest) Generate(r *rand.Rand, size int) reflect.Value {
	// Generate random actions.
	addrs := make([]common.Address, 50)
	for i := range addrs {
		addrs[i][0] = byte(i)
	}
	actions := make([]testAction, size)
	for i := range actions {
		addr := addrs[r.Intn(len(addrs))]
		actions[i] = newTestAction(addr, r)
	}
	// Generate snapshot indexes.
	nsnapshots := int(math.Sqrt(float64(size)))
	if size > 0 && nsnapshots == 0 {
		nsnapshots = 1
	}
	snapshots := make([]int, nsnapshots)
	snaplen := len(actions) / nsnapshots
	for i := range snapshots {
		// Try to place the snapshots some number of actions apart from each other.
		snapshots[i] = (i * snaplen) + r.Intn(snaplen)
	}
	return reflect.ValueOf(&snapshotTest{addrs, actions, snapshots, nil})
}

func (test *snapshotTest) String() string {
	out := new(bytes.Buffer)
	sindex := 0
	for i, action := range test.actions {
		if len(test.snapshots) > sindex && i == test.snapshots[sindex] {
			fmt.Fprintf(out, "---- snapshot %d ----\n", sindex)
			sindex++
		}
		fmt.Fprintf(out, "%4d: %s\n", i, action.name)
	}
	return out.String()
}

func (test *snapshotTest) run() bool {
	// Run all actions and create snapshots.
	var (
		state, _     = New(common.Hash{}, NewDatabase(rawdb.NewMemoryDatabase()))
		snapshotRevs = make([]int, len(test.snapshots))
		sindex       = 0
	)
	for i, action := range test.actions {
		if len(test.snapshots) > sindex && i == test.snapshots[sindex] {
			snapshotRevs[sindex] = state.Snapshot()
			sindex++
		}
		action.fn(action, state)
	}
	// Revert all snapshots in reverse order. Each revert must yield a state
	// that is equivalent to fresh state with all actions up the snapshot applied.
	for sindex--; sindex >= 0; sindex-- {
		checkstate, _ := New(common.Hash{}, state.Database())
		for _, action := range test.actions[:test.snapshots[sindex]] {
			action.fn(action, checkstate)
		}
		state.RevertToSnapshot(snapshotRevs[sindex])
		if err := test.checkEqual(state, checkstate); err != nil {
			test.err = fmt.Errorf("state mismatch after revert to snapshot %d\n%v", sindex, err)
			return false
		}
	}
	return true
}

// checkEqual checks that methods of state and checkstate return the same values.
func (test *snapshotTest) checkEqual(state, checkstate *StateDB) error {
	for _, addr := range test.addrs {
		var err error
		checkeq := func(op string, a, b interface{}) bool {
			if err == nil && !reflect.DeepEqual(a, b) {
				err = fmt.Errorf("got %s(%s) == %v, want %v", op, addr.String(), a, b)
				return false
			}
			return true
		}
		// Check basic accessor methods.
		checkeq("Exist", state.Exist(addr), checkstate.Exist(addr))
		checkeq("HasSuicided", state.HasSuicided(addr), checkstate.HasSuicided(addr))
		checkeq("GetBalance", state.GetBalance(addr), checkstate.GetBalance(addr))
		checkeq("GetNonce", state.GetNonce(addr), checkstate.GetNonce(addr))
		checkeq("GetCode", state.GetCode(addr), checkstate.GetCode(addr))
		checkeq("GetCodeHash", state.GetCodeHash(addr), checkstate.GetCodeHash(addr))
		checkeq("GetCodeSize", state.GetCodeSize(addr), checkstate.GetCodeSize(addr))
		// Check storage.
		if obj := state.getStateObject(addr); obj != nil {
			state.ForEachStorage(addr, func(key []byte, value []byte) bool {
				cobj := checkstate.getStateObject(addr)
				return checkeq("GetState("+hex.EncodeToString(key)+")", cobj.GetState(checkstate.db, key), value)
			})
			checkstate.ForEachStorage(addr, func(key []byte, value []byte) bool {
				cobj := checkstate.getStateObject(addr)
				return checkeq("GetState("+hex.EncodeToString(key)+")", cobj.GetState(checkstate.db, key), value)
			})
		}
		if err != nil {
			return err
		}
	}

	if state.GetRefund() != checkstate.GetRefund() {
		return fmt.Errorf("got GetRefund() == %d, want GetRefund() == %d",
			state.GetRefund(), checkstate.GetRefund())
	}
	if !reflect.DeepEqual(state.GetLogs(common.Hash{}), checkstate.GetLogs(common.Hash{})) {
		return fmt.Errorf("got GetLogs(common.Hash{}) == %v, want GetLogs(common.Hash{}) == %v",
			state.GetLogs(common.Hash{}), checkstate.GetLogs(common.Hash{}))
	}
	return nil
}

func (s *StateSuite) TestTouchDelete(c *check.C) {
	s.state.GetOrNewStateObject(common.Address{})
	root, _ := s.state.Commit(false)
	s.state.Reset(root)

	snapshot := s.state.Snapshot()
	s.state.AddBalance(common.Address{}, new(big.Int))

	if len(s.state.journal.dirties) != 1 {
		c.Fatal("expected one dirty state object")
	}
	s.state.RevertToSnapshot(snapshot)
	if len(s.state.journal.dirties) != 0 {
		c.Fatal("expected no dirty state object")
	}
}

// TestCopyOfCopy tests that modified objects are carried over to the copy, and the copy of the copy.
// See https://github.com/ethereum/go-ethereum/pull/15225#issuecomment-380191512
func TestCopyOfCopy(t *testing.T) {
	sdb, _ := New(common.Hash{}, NewDatabase(rawdb.NewMemoryDatabase()))
	addr := common.HexToAddress("aaaa")
	sdb.SetBalance(addr, big.NewInt(42))

	if got := sdb.Copy().GetBalance(addr).Uint64(); got != 42 {
		t.Fatalf("1st copy fail, expected 42, got %v", got)
	}
	if got := sdb.Copy().Copy().GetBalance(addr).Uint64(); got != 42 {
		t.Fatalf("2nd copy fail, expected 42, got %v", got)
	}
}

func TestGetAfterDelete(t *testing.T) {
	db := rawdb.NewMemoryDatabase()

	addr := common.BigToAddress(big.NewInt(1))

	s1, _ := New(common.Hash{}, NewDatabase(db))
	s1.SetNonce(addr, 1)
	s1.SetState(addr, []byte("test"), []byte("value"))
	_, err := s1.Commit(true)
	assert.Nil(t, err)

	s2 := s1.NewStateDB()
	s2.SetState(addr, []byte("test"), []byte{})
	_, err = s2.Commit(true)
	assert.Nil(t, err)

	s3 := s2.NewStateDB()
	buf := s3.GetState(addr, []byte("test"))
	s3.Commit(true)
	assert.True(t, len(buf) == 0, "Expect value is not nil")

	s4 := s3.NewStateDB()
	s4.SetState(addr, []byte("test"), []byte("value"))
	s4.SetState(addr, []byte("test1"), []byte("value1"))
	s4.Commit(true)

	s5 := s4.NewStateDB()
	buf = s5.GetState(addr, []byte("test"))
	buf1 := s5.GetState(addr, []byte("test1"))
	assert.Equal(t, buf, []byte("value"))
	assert.Equal(t, buf1, []byte("value1"))
}
