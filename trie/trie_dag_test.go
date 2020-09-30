package trie

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/ethdb/memorydb"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func TestTrieDag(t *testing.T) {
	triedb := NewDatabase(memorydb.New())
	tr, _ := New(common.Hash{}, triedb)

	tr.Update([]byte("doe"), []byte("reindeer"))
	tr.Update([]byte("dog"), []byte("puppy"))
	tr.Update([]byte("dogglesworth"), []byte("cat"))

	hashed, _, err := tr.parallelHashRoot(nil, nil)
	assert.Nil(t, err)

	checkr, _ := New(common.Hash{}, NewDatabase(memorydb.New()))
	checkr.Update([]byte("doe"), []byte("reindeer"))
	checkr.Update([]byte("dog"), []byte("puppy"))
	checkr.Update([]byte("dogglesworth"), []byte("cat"))
	ch, _, _ := checkr.hashRoot(nil, nil)

	assert.Equal(t, hashed, ch)
}
