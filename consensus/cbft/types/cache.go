package types

import (
	"sync"
	"time"
)

type SyncCache struct {
	lock    sync.RWMutex
	items   map[interface{}]time.Time
	timeout time.Duration
}

func NewSyncCache(timeout time.Duration) *SyncCache {
	cache := &SyncCache{
		items:   make(map[interface{}]time.Time),
		timeout: timeout,
	}
	return cache
}

func (s *SyncCache) Add(v interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.items[v] = time.Now()
}

func (s *SyncCache) AddOrReplace(v interface{}) bool {
	s.lock.RLock()
	defer s.lock.RUnlock()
	if t, ok := s.items[v]; ok {
		if time.Since(t) < s.timeout {
			return false
		}
	}
	s.items[v] = time.Now()
	return true
}

func (s *SyncCache) Remove(v interface{}) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	delete(s.items, v)
}

func (s *SyncCache) Purge() {
	s.lock.RLock()
	defer s.lock.RUnlock()
	s.items = make(map[interface{}]time.Time)
}
func (s *SyncCache) Len() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.items)
}
