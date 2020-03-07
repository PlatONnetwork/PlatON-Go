package lru

import (
	"sync"

	"github.com/PlatONnetwork/wagon/exec"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/hashicorp/golang-lru/simplelru"
)

var (
	DefaultWasmCacheSize = 1024
	wasmCache, _         = NewWasmCache(DefaultWasmCacheSize)
	DefaultWasmCacheDir  = "wasmcache"
)

type WasmLDBCache struct {
	lru  *simplelru.LRU
	lock sync.RWMutex
}

type WasmModule struct {
	Module *exec.CompiledModule
}

func WasmCache() *WasmLDBCache {
	return wasmCache
}

func NewWasmCache(size int) (*WasmLDBCache, error) {
	w := &WasmLDBCache{}
	lru, err := simplelru.NewLRU(size, nil)

	if err != nil {
		return nil, err
	}
	w.lru = lru
	return w, nil
}

func NewWasmLDBCache(size int) (*WasmLDBCache, error) {
	return NewWasmCache(size)
}

// Purge is used to completely clear the cache
func (w *WasmLDBCache) Purge() {
	w.lock.Lock()
	w.lru.Purge()
	w.lock.Unlock()
}

// Add adds a value to the cache.  Returns true if an eviction occurred.
func (w *WasmLDBCache) Add(key common.Address, value *WasmModule) bool {
	w.lock.Lock()
	defer w.lock.Unlock()
	return w.lru.Add(key, value)
}

// Get looks up a key's value from the cache.
func (w *WasmLDBCache) Get(key common.Address) (*WasmModule, bool) {
	w.lock.Lock()
	defer w.lock.Unlock()
	value, ok := w.lru.Get(key)
	if !ok {
		return nil, ok
	}
	return value.(*WasmModule), ok
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (w *WasmLDBCache) Contains(key common.Address) bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	if !w.lru.Contains(key) {
		//ok := false
		//if w.db != nil {
		//	ok, _ = w.db.Has(key.Bytes(), nil)
		//}
		//return ok

		return false
	}
	return true
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (w *WasmLDBCache) Peek(key common.Address) (*WasmModule, bool) {
	w.lock.Lock()
	defer w.lock.Unlock()
	value, ok := w.lru.Peek(key)
	return value.(*WasmModule), ok
}

// ContainsOrAdd checks if a key is in the cache  without updating the
// recent-ness or deleting it for being stale,  and if not, adds the value.
// Returns whether found and whether an eviction occurred.
func (w *WasmLDBCache) ContainsOrAdd(key common.Address, value *WasmModule) (ok, evict bool) {
	w.lock.Lock()
	defer w.lock.Unlock()

	if w.lru.Contains(key) {
		return true, false
	} else {
		evict := w.lru.Add(key, value)
		return false, evict
	}
}

// Remove removes the provided key from the cache.
func (w *WasmLDBCache) Remove(key common.Address) {
	w.lock.Lock()
	w.lru.Remove(key)
	w.lock.Unlock()
}

// RemoveOldest removes the oldest item from the cache.
func (w *WasmLDBCache) RemoveOldest() {
	w.lock.Lock()
	w.lru.RemoveOldest()
	w.lock.Unlock()
}

// Keys returns a slice of the keys in the cache, from oldest to newest.
func (w *WasmLDBCache) Keys() []interface{} {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.lru.Keys()
}

// Len returns the number of items in the cache.
func (w *WasmLDBCache) Len() int {
	w.lock.RLock()
	defer w.lock.RUnlock()
	return w.lru.Len()
}
