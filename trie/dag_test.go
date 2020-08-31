package trie

//
//import (
//	"container/list"
//	"fmt"
//	"math/big"
//	"math/rand"
//	"testing"
//	"time"
//
//	"github.com/PlatONnetwork/PlatON-Go/ethdb/memorydb"
//
//	"github.com/stretchr/testify/assert"
//
//	"github.com/PlatONnetwork/PlatON-Go/common"
//	"github.com/PlatONnetwork/PlatON-Go/crypto"
//	"github.com/PlatONnetwork/PlatON-Go/rlp"
//)
//
//func TestTrieDAG2(t *testing.T) {
//	triedb := NewDatabase(memorydb.New())
//
//	tr, _ := New(common.Hash{}, triedb)
//
//	updateString(tr, "doe", "reindeer")
//	updateString(tr, "dog", "puppy")
//	updateString(tr, "dogglesworth", "cat")
//
//	h := tr.ParallelHash()
//
//	ntr := tr.DeepCopyTrie()
//
//	h2, _ := tr.ParallelCommit(nil)
//	assert.Equal(t, h, h2)
//
//	fmt.Printf("%x\n", h)
//
//	tr0 := newEmpty()
//	updateString(tr0, "doe", "reindeer")
//	updateString(tr0, "dog", "puppy")
//	updateString(tr0, "dogglesworth", "cat")
//	hh := tr0.Hash()
//	fmt.Printf("%x\n", hh)
//
//	assert.Equal(t, h, hh)
//	nh := ntr.ParallelHash()
//	assert.Equal(t, hh, nh)
//
//	checkr, _ := New(common.Hash{}, triedb)
//
//	it := NewIterator(tr.NodeIterator(nil))
//	for it.Next() {
//		checkr.Update(it.Key, it.Value)
//	}
//
//	//h0 := tr.ParallelHash2()
//	h1 := checkr.Hash()
//	assert.Equal(t, h, h1)
//
//	deleteString(tr, "dog")
//	h = tr.ParallelHash()
//
//	deleteString(tr0, "dog")
//	hh = tr0.Hash()
//
//	assert.Equal(t, h, hh)
//
//	deleteString(ntr, "dog")
//	updateString(ntr, "dob", "dddd")
//	updateString(tr0, "dob", "dddd")
//	assert.Equal(t, ntr.ParallelHash(), tr0.ParallelHash())
//}
//
//func TestList(t *testing.T) {
//	l := list.New()
//	l.PushFront(1)
//	l.PushFront(2)
//	l.PushFront(3)
//
//	e := l.Back()
//	l.Remove(e)
//	e = l.Back()
//	l.Remove(e)
//	assert.NotNil(t, e)
//}
//
//func TestRnd2(t *testing.T) {
//	testTrieDAGRnd2(t, 1)
//	testTrieDAGRnd2(t, 10)
//	testTrieDAGRnd2(t, 100)
//	testTrieDAGRnd2(t, 248) // special point, error
//	testTrieDAGRnd2(t, 500)
//	testTrieDAGRnd2(t, 1000)
//	testTrieDAGRnd2(t, 2045) // special point, error
//	testTrieDAGRnd2(t, 10000)
//	testTrieDAGRnd2(t, 100000)
//	testTrieDAGRnd2(t, 513294) //513295
//	//testTrieDAGRnd2(t, 1000000)
//	//testTrieDAGRnd2(t, 5000000)
//}
//
//func testTrieDAGRnd2(t *testing.T, n int) {
//	// Make the random benchmark deterministic
//	random := rand.New(rand.NewSource(0))
//
//	// Create a realistic account trie to hash
//	addresses := make([][20]byte, n)
//	for i := 0; i < len(addresses); i++ {
//		for j := 0; j < len(addresses[i]); j++ {
//			addresses[i][j] = byte(random.Intn(256))
//		}
//	}
//	accounts := make([][]byte, len(addresses))
//	for i := 0; i < len(accounts); i++ {
//		var (
//			nonce   = uint64(random.Int63())
//			balance = new(big.Int).Rand(random, new(big.Int).Exp(common.Big2, common.Big256, nil))
//			root    = emptyRoot
//			code    = crypto.Keccak256(nil)
//		)
//		accounts[i], _ = rlp.EncodeToBytes([]interface{}{nonce, balance, root, code})
//	}
//	fmt.Println("gen accounts end")
//	// Insert the accounts into the trie and hash it
//	trie := newEmpty()
//	cpyTrie := newEmpty()
//	for i := 0; i < len(addresses); i++ {
//		if i == 247 {
//			trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
//			cpyTrie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
//			continue
//		}
//		if i == 2044 {
//			trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
//			cpyTrie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
//			continue
//		}
//		if i == 200000 {
//			trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
//			cpyTrie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
//			continue
//		}
//		trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
//		cpyTrie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
//	}
//
//	tm := time.Now()
//	hashed, _, err := trie.dag.hash(nil, false, nil)
//	fmt.Printf("n: %d, parallel hash duration: %s\n", n, time.Since(tm))
//	assert.Nil(t, err)
//	tm = time.Now()
//	h, _, e := cpyTrie.hashRoot(nil, nil)
//	fmt.Printf("n: %d, serial hash duration: %s\n", n, time.Since(tm))
//	assert.Nil(t, e)
//	assert.Equal(t, hashed, h)
//}
//
//func TestDAGCommit(t *testing.T) {
//	triedb := NewDatabase(memorydb.New())
//	tr, _ := New(common.Hash{}, triedb)
//
//	tr.Update([]byte("adabce"), []byte("12312312"))
//
//	tr.ParallelHash()
//	tr.ParallelCommit(nil)
//
//	tr.Update([]byte("bcdade"), []byte("12312321"))
//	tr.Update([]byte("quedad"), []byte("asfasf"))
//	tr.ParallelCommit(nil)
//	tr.ParallelCommit(nil)
//	tr.ParallelCommit(nil)
//	tr.Update([]byte("asdfasbvf"), []byte("asfasbe"))
//	tr.Update([]byte("asdfsafasfds"), []byte("fasgdsafa"))
//	tr.Delete([]byte("quedad"))
//	h := tr.ParallelHash()
//
//	tr0 := newEmpty()
//	tr0.Update([]byte("adabce"), []byte("12312312"))
//	tr0.Update([]byte("bcdade"), []byte("12312321"))
//	tr0.Update([]byte("quedad"), []byte("asfasf"))
//	tr0.Update([]byte("asdfasbvf"), []byte("asfasbe"))
//	tr0.Update([]byte("asdfsafasfds"), []byte("fasgdsafa"))
//	tr0.Delete([]byte("quedad"))
//	h0 := tr0.Hash()
//
//	assert.Equal(t, h, h0)
//	fmt.Printf("h: %x\n", h)
//	fmt.Printf("h0: %x\n", h0)
//
//	tr = newEmpty()
//	tr.Update([]byte("abc"), []byte("1111"))
//	tr.ParallelHash()
//	tr.ParallelCommit(nil)
//	tr.Update([]byte("12312"), []byte("12312312"))
//	tr.Delete([]byte("abc"))
//	tr.ParallelHash()
//	tr.ParallelCommit(nil)
//	tr.Update([]byte("12312"), []byte("12312312"))
//	h = tr.ParallelHash()
//
//	tr0 = newEmpty()
//	tr0.Update([]byte("12312"), []byte("12312312"))
//	h0 = tr0.ParallelHash()
//
//	assert.Equal(t, h, h0)
//}
//
//func TestDAGFull(t *testing.T) {
//	tr := newEmpty()
//	tr0 := newEmpty()
//
//	tr.Update([]byte("abc"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	tr0.Update([]byte("abc"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	assert.True(t, len(tr.dag.nodes) == 1)
//	tr.Update([]byte("abc"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	tr0.Update([]byte("abc"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	assert.True(t, len(tr.dag.nodes) == 1)
//
//	tr.Update([]byte("abcd"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	tr0.Update([]byte("abcd"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	assert.True(t, len(tr.dag.nodes) == 3)
//
//	tr.Update([]byte("123"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	tr0.Update([]byte("123"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	assert.True(t, len(tr.dag.nodes) == 5)
//
//	tr.Update([]byte("de"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	tr0.Update([]byte("de"), []byte("abc1231231231231231222222222222222222222222222222222222222222222222222222222222222222222222222"))
//	assert.True(t, len(tr.dag.nodes) == 7)
//
//	assert.True(t, tr.ParallelHash() == tr0.Hash())
//	tr.ParallelCommit(nil)
//
//	tr0.Commit(nil)
//}
