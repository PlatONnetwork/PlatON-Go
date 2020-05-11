package trie

import (
	"container/list"
	"errors"
	"sync"
	"sync/atomic"
)

const (
	// offset64 FNVa offset basis. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	offset64 = 14695981039346656037
	// prime64 FNVa prime value. See https://en.wikipedia.org/wiki/Fowler–Noll–Vo_hash_function#FNV-1a_hash
	prime64 = 1099511628211

	// maxValueSize The maximum size for entry value.
	maxValueSize = 1024
)

// Sum64 gets the string and returns its uint64 hash value.
func Sum64(key string) uint64 {
	var hash uint64 = offset64
	for i := 0; i < len(key); i++ {
		hash ^= uint64(key[i])
		hash *= prime64
	}

	return hash
}

// LRU implements a non-thread safe fixed size LRU cache
type LRU struct {
	capacity    uint64
	maxCapacity uint64
	evictList   *list.List
	items       map[string]*list.Element
}

// entry is used to hold a value in the evictList
type entry struct {
	key   string
	value []byte
}

// NewLRU constructs an LRU of the given size
func NewLRU(size uint64) (*LRU, error) {
	if size <= 0 {
		return nil, errors.New("Must provide a positive size")
	}
	c := &LRU{
		capacity:    0,
		maxCapacity: size,
		evictList:   list.New(),
		items:       make(map[string]*list.Element),
	}
	return c, nil
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (c *LRU) Add(key string, value []byte) bool {
	// Check for existing item
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		c.capacity -= uint64(len(ent.Value.(*entry).value))
		ent.Value.(*entry).value = value
		c.capacity += uint64(len(value))
		return false
	}

	// Add new item
	ent := &entry{key, value}
	entry := c.evictList.PushFront(ent)
	c.items[key] = entry
	c.capacity += uint64(len(value))
	c.capacity += uint64(len(key))

	evict := c.capacity > c.maxCapacity
	// Verify size not exceeded
	if evict {
		c.removeOldest()
	}
	return evict
}

// Get looks up a key's value from the cache.
func (c *LRU) Get(key string) (value []byte, ok bool) {
	if ent, ok := c.items[key]; ok {
		c.evictList.MoveToFront(ent)
		return ent.Value.(*entry).value, true
	}
	return
}

// Remove removes the provided key from the cache, returning if the
// key was contained.
func (c *LRU) Remove(key string) bool {
	if ent, ok := c.items[key]; ok {
		c.removeElement(ent)
		return true
	}
	return false
}

// Len returns the number of items in the cache.
func (c *LRU) Len() int {
	return c.evictList.Len()
}

func (c *LRU) Capacity() uint64 {
	return c.capacity
}

// removeOldest removes the oldest item from the cache.
func (c *LRU) removeOldest() {
	for c.capacity > c.maxCapacity {
		ent := c.evictList.Back()
		if ent != nil {
			c.removeElement(ent)
		}
	}
}

// removeElement is used to remove a given list element from the cache
func (c *LRU) removeElement(e *list.Element) {
	c.evictList.Remove(e)
	kv := e.Value.(*entry)
	delete(c.items, kv.key)
	c.capacity -= uint64(len(kv.value))
	c.capacity -= uint64(len(kv.key))
}

type Stats struct {
	Hits      int64 // Hits is a number of successfully found keys
	Misses    int64 // Misses is a number of not found keys
	DelHits   int64 // DelHits is a number of successfully deleted keys
	DelMisses int64 // DelMisses is a number of not deleted keys
}

type BigCache struct {
	shards    []*shard
	shardMask uint64
	lock      sync.Mutex
	lru       *LRU

	hits int64
}

type shard struct {
	entries *LRU
	lock    sync.Mutex
	stats   Stats
}

func NewBigCache(capacity uint64, size int) *BigCache {
	shards := make([]*shard, size)
	shardCap := (7 * capacity / 8) / uint64(size)
	for i := 0; i < size; i++ {
		shards[i] = newShard(shardCap)
	}
	lru, _ := NewLRU(capacity / 8)

	return &BigCache{
		shards:    shards,
		shardMask: uint64(size) - 1,
		lru:       lru,
	}
}

func (c *BigCache) Set(key string, value []byte) {
	if len(value) > maxValueSize {
		return
	}

	hash := Sum64(key)
	shard := c.getShard(hash)
	shard.set(key, value)
}

func (c *BigCache) SetLru(key string, value []byte) {
	if len(value) > maxValueSize {
		return
	}

	c.lock.Lock()
	defer c.lock.Unlock()
	c.lru.Add(key, value)
}

func (c *BigCache) Get(key string) ([]byte, bool) {
	c.lock.Lock()
	if val, ok := c.lru.Get(key); ok {
		atomic.AddInt64(&c.hits, 1)
		c.lock.Unlock()
		return val, true
	}
	c.lock.Unlock()

	hash := Sum64(key)
	shard := c.getShard(hash)
	return shard.get(key)
}

func (c *BigCache) Delele(key string) {
	hash := Sum64(key)
	shard := c.getShard(hash)
	shard.delete(key)
}

func (c *BigCache) Len() uint64 {
	len := c.lru.Len()
	for _, shard := range c.shards {
		len += shard.entries.Len()
	}
	return uint64(len)
}

func (c *BigCache) Capacity() uint64 {
	var cap uint64 = 0
	c.lock.Lock()
	cap = c.lru.Capacity()
	c.lock.Unlock()

	for _, shard := range c.shards {
		cap += shard.entries.Capacity()
	}
	return cap
}

func (c *BigCache) Stats() Stats {
	var stats Stats
	stats.Hits = atomic.LoadInt64(&c.hits)
	for _, shard := range c.shards {
		s := shard.Stats()
		stats.Hits += s.Hits
		stats.Misses += s.Misses
		stats.DelHits += s.DelHits
		stats.DelMisses += s.DelMisses
	}
	return stats
}

func (c *BigCache) ResetStats() {
	c.lock.Lock()
	c.hits = 0
	c.lock.Unlock()

	for _, shard := range c.shards {
		shard.resetStats()
	}
}

func (c *BigCache) getShard(hash uint64) *shard {
	return c.shards[c.shardMask&hash]
}

func newShard(capacity uint64) *shard {
	entries, _ := NewLRU(capacity)
	return &shard{
		entries: entries,
	}
}

func (s *shard) set(key string, value []byte) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.entries.Add(key, value)
}

func (s *shard) get(key string) ([]byte, bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if val, ok := s.entries.Get(key); ok {
		s.stats.Hits++
		return val, ok
	}
	s.stats.Misses++
	return nil, false
}

func (s *shard) delete(key string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	if s.entries.Remove(key) {
		s.stats.DelHits++
		return
	}
	s.stats.DelMisses++
}

func (s *shard) Stats() Stats {
	s.lock.Lock()
	defer s.lock.Unlock()
	return s.stats
}

func (s *shard) resetStats() {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.stats = Stats{}
}
