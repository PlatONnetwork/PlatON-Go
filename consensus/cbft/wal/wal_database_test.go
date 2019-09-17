package wal

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/stretchr/testify/assert"
)

var (
	testKey   = []byte("test-key")
	testValue = []byte("test-value")
)

func TestWALDatabase(t *testing.T) {

	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	// empty path
	_, err := createWalDB("")
	assert.NotNil(t, err)

	waldb, err := createWalDB(tempDir)
	assert.Nil(t, err)

	// Put
	assert.Nil(t, waldb.Put(testKey, testValue, &opt.WriteOptions{Sync: true}))
	// Has
	exist, err := waldb.Has(testKey)
	assert.Nil(t, err)
	assert.True(t, exist)

	// Get
	val, err := waldb.Get(testKey)
	assert.Nil(t, err)
	assert.Equal(t, testValue, val)
	_, err = waldb.Get(testValue)
	assert.NotNil(t, err)

	// Delete
	assert.Nil(t, waldb.Delete(testKey))
	exist, err = waldb.Has(testKey)
	assert.False(t, exist)

	// NewIterator
	in := 100
	out := 0
	for i := 0; i < in; i++ {
		key := append(testKey, byte(i))
		value := append(testValue, byte(i))
		waldb.Put(key, value, &opt.WriteOptions{Sync: true})
	}
	it := waldb.NewIterator(testKey, nil)
	for it.Next() {
		out += 1
	}
	assert.Equal(t, in, out)

	waldb.Close()
}
