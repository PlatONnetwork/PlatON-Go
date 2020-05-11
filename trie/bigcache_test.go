package trie

import (
	"fmt"
	"runtime"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBigCache(t *testing.T) {
	bc := NewBigCache(1000000, 6)
	assert.NotNil(t, bc)

	assert.True(t, bc.getShard(0) == bc.shards[0])
	assert.True(t, bc.getShard(5) == bc.shards[5])

	bc.Set("test", []byte("111111"))
	v, ok := bc.Get("test")
	assert.Equal(t, v, []byte("111111"))
	assert.True(t, ok)

	assert.True(t, bc.Len() == 1)
	assert.True(t, bc.Capacity() == 10)

	stats := bc.Stats()
	assert.True(t, stats.Hits == 1)
	assert.True(t, stats.Misses == 0)

	bc.Delele("test")
	_, ok = bc.Get("test")
	assert.False(t, ok)

	stats = bc.Stats()
	assert.True(t, stats.DelHits == 1)
	assert.True(t, stats.Misses == 1)

	bc.SetLru("test", []byte("111111"))
	v, ok = bc.Get("test")
	assert.True(t, ok)
	assert.Equal(t, v, []byte("111111"))

	_, ok = bc.Get("11")
	assert.False(t, ok)

	val := make([]byte, 1024)
	bc.Set("key1024", val)
	v, ok = bc.Get("key1024")
	assert.True(t, ok)
	assert.Equal(t, v, val)

	val = make([]byte, 1025)
	bc.Set("key1025", val)
	_, ok = bc.Get("key1025")
	assert.False(t, ok)

	bc.SetLru("lru1025", val)
	_, ok = bc.Get("lru1025")
	assert.False(t, ok)

	bc.SetLru("test", []byte("123"))
	v, _ = bc.Get("test")
	assert.Equal(t, v, []byte("123"))

	bc = NewBigCache(20, 1)

	bc.Set("test", []byte("111111"))
	v, _ = bc.Get("test")
	assert.Equal(t, v, []byte("111111"))
	bc.Set("test1", []byte("123"))
	_, ok = bc.Get("test")
	assert.False(t, ok)
	_, ok = bc.Get("test1")
	assert.True(t, ok)

}

func TestBigCacheConcurrency(t *testing.T) {
	bc := NewBigCache(32, 10)

	threads := runtime.NumCPU()
	var wg sync.WaitGroup
	wg.Add(threads)
	for i := 0; i < threads; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				bc.SetLru(fmt.Sprintf("%d", j), []byte(fmt.Sprintf("%d", j)))
				bc.Set(fmt.Sprintf("%d", j+1), []byte(fmt.Sprintf("%d", j+1)))
				bc.Get(fmt.Sprintf("%d", j-1))
				bc.Get(fmt.Sprintf("%d", j))
			}
			wg.Done()
		}()
	}
	wg.Wait()
}
