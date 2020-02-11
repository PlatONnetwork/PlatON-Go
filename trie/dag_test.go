package trie

import (
	"fmt"
	"math/big"
	"math/rand"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
)

func TestDAG(t *testing.T) {
	dag := NewDAG()
	assert.True(t, dag.totalVertexs == 0)
	assert.True(t, dag.totalConsumed == 0)

	dag.addVertex(1)
	dag.addVertex(2)
	dag.addVertex(3)
	dag.addVertex(4)
	dag.addVertex(5)
	assert.True(t, dag.totalVertexs == 5)

	dag.addEdge(1, 2)
	assert.True(t, dag.vtxs[1].outEdge[0] == 2)
	//assert.True(t, dag.vtxs[2].inDegree == 1)

	dag.addEdge(2, 3)
	dag.addEdge(3, 4)
	dag.addEdge(4, 5)
	dag.generate()
	assert.True(t, dag.vtxs[2].inDegree == 1)
	assert.True(t, dag.topLevel.Len() == 1)

	id := dag.waitPop()
	assert.True(t, id == 1)
	id = dag.consume(id)
	assert.True(t, id == 2)
	id = dag.consume(id)
	assert.True(t, id == 3)
	assert.True(t, dag.totalConsumed == 2)
	id = dag.consume(id)
	assert.True(t, id == 4)
	id = dag.consume(id)
	assert.True(t, id == 5)
	id = dag.consume(id)
	assert.True(t, id == invalidId)
	id = dag.waitPop()
	assert.True(t, id == invalidId)
	id = dag.waitPop()
	assert.True(t, id == invalidId)

	dag.clear()
	assert.True(t, dag.totalVertexs == 0)
	assert.True(t, dag.totalConsumed == 0)
}

func TestTrieDAG(t *testing.T) {
	dag := NewTrieDAG(0, 0)

	trie := newEmpty()

	updateString(trie, "doe", "reindeer")
	updateString(trie, "dog", "puppy")
	updateString(trie, "dogglesworth", "cat")

	cpyTrie := trie.DeepCopyTrie()

	dag.init(trie.root)
	assert.True(t, len(dag.nodes) > 0)

	hashed, _, err := dag.hash(nil, false, nil)
	assert.Nil(t, err)

	h, _, e := cpyTrie.hashRoot(nil, nil)
	assert.Nil(t, e)
	assert.Equal(t, hashed, h)

	fmt.Printf("%x\n", common.BytesToHash(h.(hashNode)))

	b := getString(trie, "doe")
	assert.Equal(t, b, []byte("reindeer"))

	h1 := trie.Hash()
	fmt.Printf("%x\n", h1)

	hash1 := trie.ParallelHash()
	assert.Equal(t, hash1, h1)
	hash, err := trie.ParallelCommit(nil)
	assert.Nil(t, err)
	assert.Equal(t, hash, h1)

	nt, err := New(hash, trie.db)

	b0 := getString(nt, "doe")
	assert.Equal(t, b, b0)

	assert.Nil(t, err)
	updateString(nt, "doog", "test")

	cpyNt := nt.DeepCopyTrie()

	assert.Equal(t, nt.ParallelHash(), cpyNt.Hash())
}

func TestTrieDAGCommit(t *testing.T) {
	triedb := NewDatabase(ethdb.NewMemDatabase())

	tr, _ := New(common.Hash{}, triedb)

	updateString(tr, "doe", "reindeer")
	updateString(tr, "dog", "puppy")
	updateString(tr, "dogglesworth", "cat")

	//h, _ := tr.ParallelCommit(nil)
	h := tr.ParallelHash()

	tr0 := newEmpty()
	updateString(tr0, "doe", "reindeer")
	updateString(tr0, "dog", "puppy")
	updateString(tr0, "dogglesworth", "cat")
	hh := tr0.Hash()

	assert.Equal(t, h, hh)

	checkr, _ := New(common.Hash{}, triedb)

	it := NewIterator(tr.NodeIterator(nil))
	for it.Next() {
		checkr.Update(it.Key, it.Value)
	}

	h0 := tr.ParallelHash()
	h1 := checkr.Hash()
	assert.Equal(t, h0, h1)
}

func testTrieDAGRnd(t *testing.T, n int) {
	// Make the random benchmark deterministic
	random := rand.New(rand.NewSource(0))

	// Create a realistic account trie to hash
	addresses := make([][20]byte, n)
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
	cpyTrie := newEmpty()
	for i := 0; i < len(addresses); i++ {
		trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
		cpyTrie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
	}

	dag := NewTrieDAG(0, 0)
	tm := time.Now()
	dag.init(trie.root)
	fmt.Printf("n: %d, init duration: %s\n", n, time.Since(tm))
	tm = time.Now()
	hashed, _, err := dag.hash(nil, false, nil)
	fmt.Printf("n: %d, hash duration: %s\n", n, time.Since(tm))
	assert.Nil(t, err)
	tm = time.Now()
	h, _, e := cpyTrie.hashRoot(nil, nil)
	fmt.Printf("n: %d, serial hash duration: %s\n", n, time.Since(tm))
	assert.Nil(t, e)
	assert.Equal(t, hashed, h)
}

func TestRnd(t *testing.T) {
	testTrieDAGRnd(t, 1)
	testTrieDAGRnd(t, 10)
	testTrieDAGRnd(t, 100)
	testTrieDAGRnd(t, 1000)
	testTrieDAGRnd(t, 10000)
	testTrieDAGRnd(t, 100000)
	//testTrieDAGRnd(t, 1000000)
	//testTrieDAGRnd(t, 5000000)
	//testTrieDAGRnd(t, 10000000)
}
