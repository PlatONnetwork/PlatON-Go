package types

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestSyncCache(t *testing.T) {
	cache := NewSyncCache(2 * time.Millisecond)
	cache.Add(1)
	assert.False(t, cache.AddOrReplace(1))
	time.Sleep(3 * time.Millisecond)
	assert.True(t, cache.AddOrReplace(1))
	cache.Remove(1)
	assert.True(t, cache.AddOrReplace(1))
	cache.AddOrReplace(2)
	cache.AddOrReplace(3)
	assert.Equal(t, 3, cache.Len())
	cache.Purge()
	assert.Equal(t, 0, cache.Len())

}
