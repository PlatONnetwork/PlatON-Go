package trie

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/stretchr/testify/assert"
)

func TestTrieDag(t *testing.T) {
	triedb := NewDatabase(ethdb.NewMemDatabase())
	tr, _ := New(common.Hash{}, triedb)

	tr.Update([]byte("doe"), []byte("reindeer"))
	tr.Update([]byte("dog"), []byte("puppy"))
	tr.Update([]byte("dogglesworth"), []byte("cat"))

	hashed, _, err := tr.parallelHashRoot(nil, nil)
	assert.Nil(t, err)

	checkr, _ := New(common.Hash{}, NewDatabase(ethdb.NewMemDatabase()))
	checkr.Update([]byte("doe"), []byte("reindeer"))
	checkr.Update([]byte("dog"), []byte("puppy"))
	checkr.Update([]byte("dogglesworth"), []byte("cat"))
	ch, _, _ := checkr.hashRoot(nil, nil)

	assert.Equal(t, hashed, ch)
}
