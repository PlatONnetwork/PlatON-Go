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

package trie

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"io/ioutil"
	"math/big"
	"math/rand"
	"os"
	"reflect"
	"testing"
	"testing/quick"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/ethdb/leveldb"

	"github.com/PlatONnetwork/PlatON-Go/ethdb/memorydb"

	"github.com/stretchr/testify/assert"

	"github.com/davecgh/go-spew/spew"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func init() {
	spew.Config.Indent = "    "
	spew.Config.DisableMethods = false
}

// Used for testing
func newEmpty() *Trie {
	trie, _ := New(common.Hash{}, NewDatabase(memorydb.New()))
	return trie
}

func TestEmptyTrie(t *testing.T) {
	var trie Trie
	res := trie.Hash()
	exp := emptyRoot
	if res != common.Hash(exp) {
		t.Errorf("expected %x got %x", exp, res)
	}
}

func TestNull(t *testing.T) {
	var trie Trie
	key := make([]byte, 32)
	value := []byte("test")
	trie.Update(key, value)
	if !bytes.Equal(trie.Get(key), value) {
		t.Fatal("wrong value")
	}
}

func TestMissingRoot(t *testing.T) {
	trie, err := New(common.HexToHash("0beec7b5ea3f0fdbc95d0dd47f3c5bc275da8a33"), NewDatabase(memorydb.New()))
	if trie != nil {
		t.Error("New returned non-nil trie for invalid root")
	}
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("New returned wrong error: %v", err)
	}
}

func TestMissingNodeDisk(t *testing.T)    { testMissingNode(t, false) }
func TestMissingNodeMemonly(t *testing.T) { testMissingNode(t, true) }

func testMissingNode(t *testing.T, memonly bool) {
	diskdb := memorydb.New()
	triedb := NewDatabase(diskdb)

	trie, _ := New(common.Hash{}, triedb)
	updateString(trie, "120000", "qwerqwerqwerqwerqwerqwerqwerqwer")
	updateString(trie, "123456", "asdfasdfasdfasdfasdfasdfasdfasdf")
	root, _ := trie.Commit(nil)
	if !memonly {
		triedb.Commit(root, true, true)
	}

	trie, _ = New(root, triedb)
	_, err := trie.TryGet([]byte("120000"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("120099"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("123456"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	err = trie.TryUpdate([]byte("120099"), []byte("zxcvzxcvzxcvzxcvzxcvzxcvzxcvzxcv"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	err = trie.TryDelete([]byte("123456"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	hash := common.HexToHash("0xe1d943cc8f061a0c0b98162830b970395ac9315654824bf21b73b891365262f9")
	if memonly {
		delete(triedb.dirties, hash)
	} else {
		diskdb.Delete(hash[:])
	}

	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("120000"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}
	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("120099"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}
	trie, _ = New(root, triedb)
	_, err = trie.TryGet([]byte("123456"))
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	trie, _ = New(root, triedb)
	err = trie.TryUpdate([]byte("120099"), []byte("zxcv"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}
	trie, _ = New(root, triedb)
	err = trie.TryDelete([]byte("123456"))
	if _, ok := err.(*MissingNodeError); !ok {
		t.Errorf("Wrong error: %v", err)
	}
}

func TestInsert(t *testing.T) {
	trie := newEmpty()

	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	exp := common.HexToHash("8aad789dff2f538bca5d8ea56e8abe10f4c7ba3a5dea95fea4cd6e7c3a1168d3")
	root := trie.Hash()
	if root != exp {
		t.Errorf("exp %x got %x", exp, root)
	}

	trie = newEmpty()
	updateString(trie, "A", "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")

	exp = common.HexToHash("d23786fb4a010da3ce639d66d5e904a11dbc02746d1ce25029e53290cabf28ab")
	root, err := trie.Commit(nil)
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}
	if root != exp {
		t.Errorf("exp %x got %x", exp, root)
	}
}

func TestGet(t *testing.T) {
	trie := newEmpty()
	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	for i := 0; i < 2; i++ {
		res := getString(trie, "dog")
		if !bytes.Equal(res, []byte("puppy")) {
			t.Errorf("expected puppy got %x", res)
		}

		unknown := getString(trie, "unknown")
		if unknown != nil {
			t.Errorf("expected nil got %x", unknown)
		}

		if i == 1 {
			return
		}
		trie.Commit(nil)
	}
}

func TestDelete(t *testing.T) {
	trie := newEmpty()
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
			updateString(trie, val.k, val.v)
		} else {
			deleteString(trie, val.k)
		}
	}

	hash := trie.Hash()
	exp := common.HexToHash("5991bb8c6514148a29db676a14ac506cd2cd5775ace63c30a4fe457715e9ac84")
	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}

func TestEmptyValues(t *testing.T) {
	trie := newEmpty()

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
		updateString(trie, val.k, val.v)
	}

	hash := trie.Hash()
	exp := common.HexToHash("5991bb8c6514148a29db676a14ac506cd2cd5775ace63c30a4fe457715e9ac84")
	if hash != exp {
		t.Errorf("expected %x got %x", exp, hash)
	}
}

func TestReplication(t *testing.T) {
	trie := newEmpty()
	vals := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		{"shaman", "horse"},
		{"doge", "coin"},
		{"dog", "puppy"},
		{"somethingveryoddindeedthis is", "myothernodedata"},
	}
	for _, val := range vals {
		updateString(trie, val.k, val.v)
	}
	exp, err := trie.Commit(nil)
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}

	// create a new trie on top of the database and check that lookups work.
	trie2, err := New(exp, trie.db)
	if err != nil {
		t.Fatalf("can't recreate trie at %x: %v", exp, err)
	}
	for _, kv := range vals {
		if string(getString(trie2, kv.k)) != kv.v {
			t.Errorf("trie2 doesn't have %q => %q", kv.k, kv.v)
		}
	}
	hash, err := trie2.Commit(nil)
	if err != nil {
		t.Fatalf("commit error: %v", err)
	}
	if hash != exp {
		t.Errorf("root failure. expected %x got %x", exp, hash)
	}

	// perform some insertions on the new trie.
	vals2 := []struct{ k, v string }{
		{"do", "verb"},
		{"ether", "wookiedoo"},
		{"horse", "stallion"},
		// {"shaman", "horse"},
		// {"doge", "coin"},
		// {"ether", ""},
		// {"dog", "puppy"},
		// {"somethingveryoddindeedthis is", "myothernodedata"},
		// {"shaman", ""},
	}
	for _, val := range vals2 {
		updateString(trie2, val.k, val.v)
	}
	if hash := trie2.Hash(); hash != exp {
		t.Errorf("root failure. expected %x got %x", exp, hash)
	}
}

func TestLargeValue(t *testing.T) {
	trie := newEmpty()
	trie.Update([]byte("key1"), []byte{99, 99, 99, 99})
	trie.Update([]byte("key2"), bytes.Repeat([]byte{1}, 32))
	trie.Hash()
}

// TestRandomCases tests som cases that were found via random fuzzing
func TestRandomCases(t *testing.T) {
	var rt []randTestStep = []randTestStep{
		{op: 6, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 0
		{op: 6, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 1
		{op: 0, key: common.Hex2Bytes("d51b182b95d677e5f1c82508c0228de96b73092d78ce78b2230cd948674f66fd1483bd"), value: common.Hex2Bytes("0000000000000002")},           // step 2
		{op: 2, key: common.Hex2Bytes("c2a38512b83107d665c65235b0250002882ac2022eb00711552354832c5f1d030d0e408e"), value: common.Hex2Bytes("")},                         // step 3
		{op: 3, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 4
		{op: 3, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 5
		{op: 6, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 6
		{op: 3, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 7
		{op: 0, key: common.Hex2Bytes("c2a38512b83107d665c65235b0250002882ac2022eb00711552354832c5f1d030d0e408e"), value: common.Hex2Bytes("0000000000000008")},         // step 8
		{op: 0, key: common.Hex2Bytes("d51b182b95d677e5f1c82508c0228de96b73092d78ce78b2230cd948674f66fd1483bd"), value: common.Hex2Bytes("0000000000000009")},           // step 9
		{op: 2, key: common.Hex2Bytes("fd"), value: common.Hex2Bytes("")},                                                                                               // step 10
		{op: 6, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 11
		{op: 6, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 12
		{op: 0, key: common.Hex2Bytes("fd"), value: common.Hex2Bytes("000000000000000d")},                                                                               // step 13
		{op: 6, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 14
		{op: 1, key: common.Hex2Bytes("c2a38512b83107d665c65235b0250002882ac2022eb00711552354832c5f1d030d0e408e"), value: common.Hex2Bytes("")},                         // step 15
		{op: 3, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 16
		{op: 0, key: common.Hex2Bytes("c2a38512b83107d665c65235b0250002882ac2022eb00711552354832c5f1d030d0e408e"), value: common.Hex2Bytes("0000000000000011")},         // step 17
		{op: 5, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 18
		{op: 3, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 19
		{op: 0, key: common.Hex2Bytes("d51b182b95d677e5f1c82508c0228de96b73092d78ce78b2230cd948674f66fd1483bd"), value: common.Hex2Bytes("0000000000000014")},           // step 20
		{op: 0, key: common.Hex2Bytes("d51b182b95d677e5f1c82508c0228de96b73092d78ce78b2230cd948674f66fd1483bd"), value: common.Hex2Bytes("0000000000000015")},           // step 21
		{op: 0, key: common.Hex2Bytes("c2a38512b83107d665c65235b0250002882ac2022eb00711552354832c5f1d030d0e408e"), value: common.Hex2Bytes("0000000000000016")},         // step 22
		{op: 5, key: common.Hex2Bytes(""), value: common.Hex2Bytes("")},                                                                                                 // step 23
		{op: 1, key: common.Hex2Bytes("980c393656413a15c8da01978ed9f89feb80b502f58f2d640e3a2f5f7a99a7018f1b573befd92053ac6f78fca4a87268"), value: common.Hex2Bytes("")}, // step 24
		{op: 1, key: common.Hex2Bytes("fd"), value: common.Hex2Bytes("")},                                                                                               // step 25
	}
	runRandTest(rt)

}

// randTest performs random trie operations.
// Instances of this test are created by Generate.
type randTest []randTestStep

type randTestStep struct {
	op    int
	key   []byte // for opUpdate, opDelete, opGet
	value []byte // for opUpdate
	err   error  // for debugging
}

const (
	opUpdate = iota
	opDelete
	opGet
	opCommit
	opHash
	opReset
	opItercheckhash
	opMax // boundary value, not an actual op
)

func (randTest) Generate(r *rand.Rand, size int) reflect.Value {
	var allKeys [][]byte
	genKey := func() []byte {
		if len(allKeys) < 2 || r.Intn(100) < 10 {
			// new key
			key := make([]byte, r.Intn(50))
			r.Read(key)
			allKeys = append(allKeys, key)
			return key
		}
		// use existing key
		return allKeys[r.Intn(len(allKeys))]
	}

	var steps randTest
	for i := 0; i < size; i++ {
		step := randTestStep{op: r.Intn(opMax)}
		switch step.op {
		case opUpdate:
			step.key = genKey()
			step.value = make([]byte, 8)
			binary.BigEndian.PutUint64(step.value, uint64(i))
		case opGet, opDelete:
			step.key = genKey()
		}
		steps = append(steps, step)
	}
	return reflect.ValueOf(steps)
}

func runRandTest(rt randTest) bool {
	triedb := NewDatabase(memorydb.New())

	tr, _ := New(common.Hash{}, triedb)
	values := make(map[string]string) // tracks content of the trie

	for i, step := range rt {
		switch step.op {
		case opUpdate:
			tr.Update(step.key, step.value)
			values[string(step.key)] = string(step.value)
		case opDelete:
			tr.Delete(step.key)
			delete(values, string(step.key))
		case opGet:
			v := tr.Get(step.key)
			want := values[string(step.key)]
			if string(v) != want {
				rt[i].err = fmt.Errorf("mismatch for key 0x%x, got 0x%x want 0x%x", step.key, v, want)
			}
		case opCommit:
			_, rt[i].err = tr.Commit(nil)
		case opHash:
			tr.Hash()
		case opReset:
			hash, err := tr.Commit(nil)
			if err != nil {
				rt[i].err = err
				return false
			}
			newtr, err := New(hash, triedb)
			if err != nil {
				rt[i].err = err
				return false
			}
			tr = newtr
		case opItercheckhash:
			checktr, _ := New(common.Hash{}, triedb)
			it := NewIterator(tr.NodeIterator(nil))
			for it.Next() {
				checktr.Update(it.Key, it.Value)
			}
			if tr.Hash() != checktr.Hash() {
				//fmt.Printf("phash: %x, chash: %x\n", tr.Hash(), checktr.Hash())
				rt[i].err = fmt.Errorf("hash mismatch in opItercheckhash")
			}
		}
		// Abort the test on error.
		if rt[i].err != nil {
			//fmt.Printf("i: %d, err: %v, i-1_op: %d\n", i, rt[i].err, rt[i-1].op)
			return false
		}
	}
	return true
}

func runRandParallelTest(rt randTest) bool {
	triedb := NewDatabase(memorydb.New())

	tr, _ := New(common.Hash{}, triedb)
	values := make(map[string]string) // tracks content of the trie
	tmpVals := make(map[string][]byte)

	for i, step := range rt {
		switch step.op {
		case opUpdate:
			//fmt.Printf("%d: opUpdate, len: %d\n", i, len(values))
			tr.Update(step.key, step.value)
			values[string(step.key)] = string(step.value)
			tmpVals[string(step.key)] = step.value
		case opDelete:
			//fmt.Printf("%d: opDelete, len: %d\n", i, len(values))
			tr.Delete(step.key)
			//fmt.Printf("del -> %x\n", step.key)
			delete(values, string(step.key))
			delete(tmpVals, string(step.key))
		case opGet:
			//fmt.Printf("%d: opGet, len: %d\n", i, len(values))
			v := tr.Get(step.key)
			want := values[string(step.key)]
			if string(v) != want {
				rt[i].err = fmt.Errorf("mismatch for key 0x%x, got 0x%x want 0x%x", step.key, v, want)
				tr.Get(step.key)
			}
		case opCommit:
			//fmt.Printf("%d: opGet, len: %d\n", i, len(values))
			_, rt[i].err = tr.ParallelCommit(nil)
		case opHash:
			//fmt.Printf("%d: opHash, len: %d\n", i, len(values))
			tr.ParallelHash()
		case opReset:
			//fmt.Printf("%d: opReset, len: %d\n", i, len(values))
			hash, err := tr.ParallelCommit(nil)
			if err != nil {
				rt[i].err = err
				return false
			}
			newtr, err := New(hash, triedb)
			if err != nil {
				rt[i].err = err
				return false
			}
			tr = newtr
		case opItercheckhash:
			//fmt.Printf("%d: opItercheckhash, len: %d\n", i, len(values))
			checktr, _ := New(common.Hash{}, triedb)
			it := NewIterator(tr.NodeIterator(nil))
			for it.Next() {
				checktr.Update(it.Key, it.Value)
			}
			if tr.ParallelHash() != checktr.Hash() {
				//fmt.Printf("phash: %x, chash: %x\n", tr.ParallelHash2(), checktr.Hash())
				rt[i].err = fmt.Errorf("hash mismatch in opItercheckhash")

				nt, _ := New(common.Hash{}, triedb)
				it := NewIterator(tr.NodeIterator(nil))
				for it.Next() {
					nt.Update(it.Key, it.Value)
				}
			}
		}
		// Abort the test on error.
		if rt[i].err != nil {
			//fmt.Printf("i: %d, err: %v, i-1_op: %d, i_op: %d\n", i, rt[i].err, rt[i-1].op, rt[i].op)
			return false
		}
	}
	return true
}

func TestNewFlag(t *testing.T) {
	trie := &Trie{}
	trie.newFlag()
	trie.newFlag()

}

func TestRandom(t *testing.T) {
	if err := quick.Check(runRandTest, nil); err != nil {
		if cerr, ok := err.(*quick.CheckError); ok {
			t.Fatalf("random test iteration %d failed: %s", cerr.Count, spew.Sdump(cerr.In))
		}
		t.Fatal(err)
	}
}

func TestRandomParalle(t *testing.T) {
	if err := quick.Check(runRandParallelTest, nil); err != nil {
		if cerr, ok := err.(*quick.CheckError); ok {
			t.Fatalf("random test iteration %d failed: %s", cerr.Count, spew.Sdump(cerr.In))
		}
		t.Fatal(err)
	}
}

func BenchmarkGet(b *testing.B)      { benchGet(b, false) }
func BenchmarkGetDB(b *testing.B)    { benchGet(b, true) }
func BenchmarkUpdateBE(b *testing.B) { benchUpdate(b, binary.BigEndian) }
func BenchmarkUpdateLE(b *testing.B) { benchUpdate(b, binary.LittleEndian) }

const benchElemCount = 20000

func benchGet(b *testing.B, commit bool) {
	trie := new(Trie)
	if commit {
		_, tmpdb := tempDB()
		trie, _ = New(common.Hash{}, tmpdb)
	}
	k := make([]byte, 32)
	for i := 0; i < benchElemCount; i++ {
		binary.LittleEndian.PutUint64(k, uint64(i))
		trie.Update(k, k)
	}
	binary.LittleEndian.PutUint64(k, benchElemCount/2)
	if commit {
		trie.Commit(nil)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		trie.Get(k)
	}
	b.StopTimer()

	if commit {
		ldb := trie.db.diskdb.(*leveldb.Database)
		ldb.Close()
		os.RemoveAll(ldb.Path())
	}
}

func benchUpdate(b *testing.B, e binary.ByteOrder) *Trie {
	trie := newEmpty()
	k := make([]byte, 32)
	for i := 0; i < b.N; i++ {
		e.PutUint64(k, uint64(i))
		trie.Update(k, k)
	}
	return trie
}

// Benchmarks the trie hashing. Since the trie caches the result of any operation,
// we cannot use b.N as the number of hashing rouns, since all rounds apart from
// the first one will be NOOP. As such, we'll use b.N as the number of account to
// insert into the trie before measuring the hashing.
func BenchmarkHash(b *testing.B) {
	// Make the random benchmark deterministic
	random := rand.New(rand.NewSource(0))

	// Create a realistic account trie to hash
	addresses := make([][20]byte, b.N)
	for i := 0; i < len(addresses); i++ {
		for j := 0; j < len(addresses[i]); j++ {
			addresses[i][j] = byte(random.Intn(256))
		}
	}
	accounts := make([][]byte, len(addresses))
	for i := 0; i < len(accounts); i++ {
		var (
			nonce   = uint64(random.Int63())
			balance = new(big.Int).Rand(random, new(big.Int).Exp(common.Big2, common.Big256, nil))
			root    = emptyRoot
			code    = crypto.Keccak256(nil)
		)
		accounts[i], _ = rlp.EncodeToBytes([]interface{}{nonce, balance, root, code})
	}
	// Insert the accounts into the trie and hash it
	trie := newEmpty()
	trie.dag = nil
	for i := 0; i < len(addresses); i++ {
		trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
	}
	b.ResetTimer()
	b.ReportAllocs()
	trie.Hash()
}

func BenchmarkParallelHash2(b *testing.B) {
	// Make the random benchmark deterministic
	random := rand.New(rand.NewSource(0))

	// Create a realistic account trie to hash
	addresses := make([][20]byte, b.N)
	for i := 0; i < len(addresses); i++ {
		for j := 0; j < len(addresses[i]); j++ {
			addresses[i][j] = byte(random.Intn(256))
		}
	}
	accounts := make([][]byte, len(addresses))
	for i := 0; i < len(accounts); i++ {
		var (
			nonce   = uint64(random.Int63())
			balance = new(big.Int).Rand(random, new(big.Int).Exp(common.Big2, common.Big256, nil))
			root    = emptyRoot
			code    = crypto.Keccak256(nil)
		)
		accounts[i], _ = rlp.EncodeToBytes([]interface{}{nonce, balance, root, code})
	}
	// Insert the accounts into the trie and hash it
	trie := newEmpty()
	for i := 0; i < len(addresses); i++ {
		trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
	}
	b.ResetTimer()
	b.ReportAllocs()
	trie.ParallelHash()
}

func tempDB() (string, *Database) {
	dir, err := ioutil.TempDir("", "trie-bench")
	if err != nil {
		panic(fmt.Sprintf("can't create temporary directory: %v", err))
	}
	diskdb, err := leveldb.New(dir, 256, 0, "")
	if err != nil {
		panic(fmt.Sprintf("can't create temporary database: %v", err))
	}
	return dir, NewDatabase(diskdb)
}

func getString(trie *Trie, k string) []byte {
	return trie.Get([]byte(k))
}

func updateString(trie *Trie, k, v string) {
	trie.Update([]byte(k), []byte(v))
}

func deleteString(trie *Trie, k string) {
	trie.Delete([]byte(k))
}

func TestDeepCopy(t *testing.T) {
	memdb := memorydb.New()
	triedb := NewDatabase(memdb)
	root := common.Hash{}
	tr, _ := NewSecure(root, triedb)
	kv := make(map[common.Hash][]byte)
	leafCB := func(leaf []byte, parent common.Hash) error {
		var valueKey common.Hash
		_, content, _, err := rlp.Split(leaf)
		assert.Nil(t, err)
		valueKey.SetBytes(content)
		if value, ok := kv[valueKey]; ok {
			tr.trie.db.InsertBlob(valueKey, value)
		}

		tr.trie.db.Reference(valueKey, parent)
		return nil
	}
	k, v := randBytes(32), randBytes(32)
	parent := root
	for j := 0; j < 1; j++ {
		for i := 1; i < 100; i++ {
			binary.BigEndian.PutUint32(k, uint32(i))
			binary.BigEndian.PutUint32(v, uint32(i))
			tr.Update(k, v)
			kv[common.BytesToHash(tr.hashKey(k))] = v
		}

		root, _ = tr.Commit(leafCB)
		parent = root
		triedb.Reference(root, common.Hash{})
		triedb.Commit(root, false, false)
		//fmt.Println("commit db", "count", j, time.Since(start))
	}

	tr2, _ := NewSecure(root, triedb)
	for i := 100; i < 200; i++ {
		binary.BigEndian.PutUint32(k, uint32(i))
		binary.BigEndian.PutUint32(v, uint32(i))
		tr2.Update(k, v)
		kv[common.BytesToHash(tr.hashKey(k))] = v
	}

	//root, _ = tr2.Commit(nil)
	root = tr2.Hash()

	cpy := tr2.New().New()

	iter := tr2.NodeIterator(nil)
	cpyIter := cpy.NodeIterator(nil)
	count := 0
	keys := 0
	for iter.Next(true) {
		if !cpyIter.Next(true) {
			t.Fatal("cpy iter failed, next error")
		}
		if !bytes.Equal(iter.Path(), cpyIter.Path()) {
			t.Fatal("iter path failed")
		}
		if !bytes.Equal(iter.Parent().Bytes(), cpyIter.Parent().Bytes()) {
			t.Fatal("iter parent failed")
		}
		if iter.Leaf() {
			if !bytes.Equal(iter.LeafBlob(), iter.LeafBlob()) {
				t.Fatal("iter leaf blob failed")
			}
			if !bytes.Equal(iter.LeafKey(), iter.LeafKey()) {
				t.Fatal("iter leaf key failed")
			}
			if _, ok := kv[common.BytesToHash(iter.LeafKey())]; !ok {
				t.Fatal("find none key")
			}
			keys++
			//fmt.Println(hexutil.Encode(iter.LeafKey()))
			//delete(kv, common.BytesToHash(iter.LeafKey()))
		}
		if iter.Hash() != cpyIter.Hash() {
			t.Fatal("cpy iter failed", iter.Hash(), cpyIter.Hash())
		}
		count++
	}
	assert.Equal(t, len(kv), keys)
	root, _ = tr2.Commit(leafCB)
	triedb.Reference(root, common.Hash{})
	assert.Nil(t, triedb.Commit(root, false, false))
	triedb.DereferenceDB(parent)
	cpyRoot, _ := cpy.Commit(leafCB)
	if root != cpyRoot {
		t.Fatal("cpyroot failed")
	}
	triedb.Reference(cpyRoot, common.Hash{})

	assert.Nil(t, triedb.Commit(cpyRoot, false, false))
	triedb.DereferenceDB(cpyRoot)
}

type Case struct {
	hash  []byte
	value []byte
}

func TestOneTrieCollision(t *testing.T) {
	trieData1 := []Case{
		{common.BytesToHash([]byte{2, 1, 1, 2}).Bytes(), []byte{2, 2, 2, 2, 2, 2}},
		{common.BytesToHash([]byte{2, 1, 1, 3}).Bytes(), []byte{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}},
		{common.BytesToHash([]byte{1, 1, 2, 0}).Bytes()[:31], []byte{2, 2, 2, 2, 2, 2}},
		{common.BytesToHash([]byte{1, 1, 3, 0}).Bytes()[:31], []byte{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}},
	}

	checkTrie := func(trie *Trie) error {
		iter := trie.NodeIterator(nil)

		for iter.Next(true) {
		}
		if iter.Error() != nil {
			return iter.Error()
		}
		return nil
	}

	mem := memorydb.New()
	memdb := NewDatabase(mem)
	trie, _ := New(common.Hash{}, memdb)
	for _, d := range trieData1 {
		trie.Update(d.hash, d.value)
	}
	root, _ := trie.Commit(nil)
	memdb.Commit(root, false, false)

	assert.Nil(t, checkTrie(trie))
	reopenMemdb := NewDatabase(mem)

	reopenTrie, _ := New(root, reopenMemdb)
	reopenTrie.Delete(trieData1[0].hash)

	reopenRoot, _ := reopenTrie.Commit(nil)
	reopenMemdb.Commit(reopenRoot, false, false)
	reopenTrie.Update(trieData1[0].hash, trieData1[0].value)
	reopenRoot, _ = reopenTrie.Commit(nil)
	reopenMemdb.Commit(reopenRoot, false, false)
	reopenMemdb.Reference(root, common.Hash{})

	reopenTrie.Delete(trieData1[0].hash)
	reopenRoot, _ = reopenTrie.Commit(nil)
	reopenMemdb.Commit(reopenRoot, false, false)
	reopenMemdb.Reference(reopenRoot, common.Hash{})
	reopenMemdb.DereferenceDB(root)
	reopenMemdb.UselessGC(1)

	assert.NotNil(t, checkTrie(reopenTrie))

}

func TestTwoTrieCollision(t *testing.T) {
	trieData1 := []Case{
		{common.BytesToHash([]byte{1, 2}).Bytes(), []byte{1, 1, 1, 1, 1}},
		{common.BytesToHash([]byte{1, 3}).Bytes(), []byte{2, 2, 2, 2, 2, 2}},
		{common.BytesToHash([]byte{1, 0}).Bytes()[:31], []byte{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}},
	}
	trieData2 := []Case{
		{common.BytesToHash([]byte{2, 2}).Bytes(), []byte{1, 1, 1, 1, 1}},
		{common.BytesToHash([]byte{2, 3}).Bytes(), []byte{2, 2, 2, 2, 2, 2}},
		{common.BytesToHash([]byte{2, 0}).Bytes()[:31], []byte{3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3}},
	}

	mem1 := memorydb.New()
	memdb1 := NewDatabase(mem1)
	mem2 := memorydb.New()
	memdb2 := NewDatabase(mem2)
	trie1, _ := New(common.Hash{}, memdb1)
	trie2, _ := New(common.Hash{}, memdb2)
	for _, d := range trieData1 {
		trie1.Update(d.hash, d.value)
	}
	for _, d := range trieData2 {
		trie2.Update(d.hash, d.value)
	}

	root1, _ := trie1.Commit(nil)
	root2, _ := trie2.Commit(nil)

	memdb1.Commit(root1, false, false)
	memdb2.Commit(root2, false, false)

	dup := 0
	itr := mem1.NewIterator()
	for itr.Next() {
		_, err := mem2.Get(itr.Key())
		if err != nil {
			dup++
		}
	}
	assert.NotZero(t, dup)
}

func TestCommitAfterHash(t *testing.T) {
	// Create a realistic account trie to hash
	addresses, accounts := makeAccounts(1000)
	trie := newEmpty()
	for i := 0; i < len(addresses); i++ {
		trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
	}
	// Insert the accounts into the trie and hash it
	trie.Hash()
	trie.Commit(nil)
	root := trie.Hash()
	exp := common.HexToHash("e5e9c29bb50446a4081e6d1d748d2892c6101c1e883a1f77cf21d4094b697822")
	if exp != root {
		t.Errorf("got %x, exp %x", root, exp)
	}
	root, _ = trie.Commit(nil)
	if exp != root {
		t.Errorf("got %x, exp %x", root, exp)
	}
}

func makeAccounts(size int) (addresses [][20]byte, accounts [][]byte) {
	// Make the random benchmark deterministic
	random := rand.New(rand.NewSource(0))
	// Create a realistic account trie to hash
	addresses = make([][20]byte, size)
	for i := 0; i < len(addresses); i++ {
		for j := 0; j < len(addresses[i]); j++ {
			addresses[i][j] = byte(random.Intn(256))
		}
	}
	accounts = make([][]byte, len(addresses))
	for i := 0; i < len(accounts); i++ {
		var (
			nonce   = uint64(random.Int63())
			balance = new(big.Int).Rand(random, new(big.Int).Exp(common.Big2, common.Big256, nil))
			root    = emptyRoot
			code    = crypto.Keccak256(nil)
		)
		accounts[i], _ = rlp.EncodeToBytes(&account{nonce, balance, root, code})
	}
	return addresses, accounts
}

type account struct {
	Nonce   uint64
	Balance *big.Int
	Root    common.Hash
	Code    []byte
}

type TriekvPair struct {
	k []byte
	v []byte
}

func genTriekvPairs(n int) []*TriekvPair {
	var keyPrefix = [][]byte{[]byte("kk1"), []byte("kk1kk1"), []byte("kk2"), []byte("kk2k2"), []byte("kk2kk2")}
	var randomKey = []byte{'1', '2', 'k', '1', '2', 'k', 'a', 'b', 'c', 'd'}

	appendKey := func(n int) []byte {
		rk := make([]byte, 0, n)
		for i := 0; i < n*2; i++ {
			rk = append(rk, randomKey[rand.Intn(len(randomKey))])
		}
		return rk
	}

	triekvPairs := make([]*TriekvPair, 0, n)
	for i := 0; i < n; i++ {
		k := keyPrefix[i%5]
		k = append(append(k, appendKey(i%10)...), byte(i))
		v := append([]byte("test"), byte(i))
		triekvPairs = append(triekvPairs, &TriekvPair{
			k: k,
			v: v,
		})
	}
	return triekvPairs
}

func orderDisrupted(triekvPairs []*TriekvPair) []*TriekvPair {
	for i := 0; i < len(triekvPairs); i++ {
		rand := rand.Intn(len(triekvPairs))
		//fmt.Println("rand", "rand", rand)
		temp := triekvPairs[rand]
		triekvPairs[rand] = triekvPairs[i]
		triekvPairs[i] = temp
	}
	return triekvPairs
}

func TestTrieHashByDisorderedData(t *testing.T) {
	//triekvPairs := genTriekvPairs(1000000)
	triekvPairs := genTriekvPairs(10000)

	start := time.Now()
	var trie Trie
	for i := 0; i < len(triekvPairs); i++ {
		err := trie.TryUpdate(triekvPairs[i].k, triekvPairs[i].v)
		if err != nil {
			t.Errorf("TryUpdate Error")
		}
	}
	rootHash := trie.Hash()
	t.Log("Update trie success", "root", rootHash.String(), "duration", time.Since(start))

	// Disrupted order
	for i := 0; i < 1; i++ {
		start = time.Now()
		trie2, _ := New(common.Hash{}, NewDatabase(memorydb.New()))
		triekvPairs2 := orderDisrupted(triekvPairs)
		for i := 0; i < len(triekvPairs2); i++ {
			err := trie2.TryUpdate(triekvPairs2[i].k, triekvPairs2[i].v)
			if err != nil {
				t.Errorf("TryUpdate Error")
			}
		}
		rootHash2 := trie2.ParallelHash()
		t.Log("Update trie success", "root", rootHash2.String(), "duration", time.Since(start))

		assert.Equal(t, rootHash, rootHash2)
	}
}

func TestTrieHashByUpdate(t *testing.T) {
	//triekvPairs := genTriekvPairs(1000000)
	triekvPairs := genTriekvPairs(10000)

	start := time.Now()
	var trie Trie
	for i := 0; i < len(triekvPairs); i++ {
		err := trie.TryUpdate(triekvPairs[i].k, triekvPairs[i].v)
		if err != nil {
			t.Errorf("TryUpdate Error")
		}
	}
	// Randomly update key or delete key
	for i := 0; i < len(triekvPairs); i++ {
		if i%2 == 0 {
			// update key
			trie.TryUpdate(triekvPairs[i].k, byteutil.Concat(triekvPairs[i].v, []byte("update")...))
		} else {
			// delete key
			trie.TryDelete(triekvPairs[i].k)
		}
	}
	rootHash := trie.Hash()
	t.Log("Update trie success", "root", rootHash.String(), "duration", time.Since(start))

	// Dag Trie
	for i := 0; i < 1; i++ {
		start = time.Now()
		trie2, _ := New(common.Hash{}, NewDatabase(memorydb.New()))
		for i := 0; i < len(triekvPairs); i++ {
			err := trie2.TryUpdate(triekvPairs[i].k, triekvPairs[i].v)
			if err != nil {
				t.Errorf("TryUpdate Error")
			}
		}
		// Randomly update key or delete key
		for i := 0; i < len(triekvPairs); i++ {
			if i%2 == 0 {
				// update key
				trie2.TryUpdate(triekvPairs[i].k, byteutil.Concat(triekvPairs[i].v, []byte("update")...))
			} else {
				// delete key
				trie2.TryDelete(triekvPairs[i].k)
			}
		}
		rootHash2 := trie2.ParallelHash()
		t.Log("Update trie success", "root", rootHash2.String(), "duration", time.Since(start))
		assert.Equal(t, rootHash, rootHash2)
	}
}
