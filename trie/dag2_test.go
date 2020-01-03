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

func TestTrieDAG2(t *testing.T) {
	triedb := NewDatabase(ethdb.NewMemDatabase())

	tr, _ := New(common.Hash{}, triedb)

	updateString(tr, "doe", "reindeer")
	updateString(tr, "dog", "puppy")
	updateString(tr, "dogglesworth", "cat")

	//h, _ := tr.ParallelCommit(nil)
	h := tr.ParallelHash2()

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

	//h0 := tr.ParallelHash2()
	h1 := checkr.Hash()
	assert.Equal(t, h, h1)

	tr.dag.clear()
	deleteString(tr, "dog")
	h = tr.ParallelHash2()

	deleteString(tr0, "dog")
	hh = tr0.Hash()

	assert.Equal(t, h, hh)
}

func TestRnd2(t *testing.T) {
	testTrieDAGRnd2(t, 1)
	testTrieDAGRnd2(t, 10)
	testTrieDAGRnd2(t, 100)
	testTrieDAGRnd2(t, 248) // special point, error
	testTrieDAGRnd2(t, 500)
	testTrieDAGRnd2(t, 1000)
	testTrieDAGRnd2(t, 2045) // special point, error
	testTrieDAGRnd2(t, 10000)
	//testTrieDAGRnd2(t, 513294) //513295
	testTrieDAGRnd2(t, 1000000)
	//testTrieDAGRnd2(t, 5000000)
}

func testTrieDAGRnd2(t *testing.T, n int) {
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
	fmt.Println("gen accounts end")
	// Insert the accounts into the trie and hash it
	trie := newEmpty()
	cpyTrie := newEmpty()
	for i := 0; i < len(addresses); i++ {
		if i == 247 {
			trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
			cpyTrie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
			continue
		}
		if i == 2044 {
			trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
			cpyTrie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
			continue
		}
		if i == 200000 {
			trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
			cpyTrie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
			continue
		}
		trie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
		cpyTrie.Update(crypto.Keccak256(addresses[i][:]), accounts[i])
	}

	tm := time.Now()
	trie.dag.init(trie.root)
	hashed, _, err := trie.dag.hash(nil, false, nil)
	fmt.Printf("parallel hash duration: %s\n", time.Since(tm))
	assert.Nil(t, err)
	tm = time.Now()
	h, _, e := cpyTrie.hashRoot(nil, nil)
	fmt.Printf("serial hash duration: %s\n", time.Since(tm))
	assert.Nil(t, e)
	assert.Equal(t, hashed, h)
}
