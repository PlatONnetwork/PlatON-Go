package lru

import (
	"bytes"
	"encoding/gob"
	"path/filepath"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/life/compiler"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/hashicorp/golang-lru/simplelru"
	"github.com/syndtr/goleveldb/leveldb"
)

var (
	DefaultWasmCacheSize = 1024
	wasmCache, _         = NewWasmCache(DefaultWasmCacheSize)
	DefaultWasmCacheDir  = "wasmcache"
)

type WasmLDBCache struct {
	lru  *simplelru.LRU
	db   *leveldb.DB
	lock sync.RWMutex
}

type WasmModule struct {
	Module       *compiler.Module
	FunctionCode []compiler.InterpreterCode
}

func WasmCache() *WasmLDBCache {
	return wasmCache
}

func SetWasmDB(dataDir string) error {
	path := filepath.Join(dataDir, DefaultWasmCacheDir)

	db, err := leveldb.OpenFile(path, nil)
	if err != nil {
		return err
	}
	wasmCache.SetDB(db)
	return nil
}

func NewWasmCache(size int) (*WasmLDBCache, error) {
	w := &WasmLDBCache{}

	onEvicted := func(k interface{}, v interface{}) {
		var addr common.Address
		var module *WasmModule
		var ok bool

		if addr, ok = k.(common.Address); !ok {
			return
		}

		if module, ok = v.(*WasmModule); !ok {
			return
		}
		if w.db != nil {
			if ok, err := w.db.Has(addr.Bytes(), nil); err != nil || !ok {
				buffer := new(bytes.Buffer)
				enc := gob.NewEncoder(buffer)
				if err := enc.Encode(module); err != nil {
					log.Error("encode module err:", err)
					return
				}
				w.db.Put(addr.Bytes(), buffer.Bytes(), nil)
			}
		}
	}

	lru, err := simplelru.NewLRU(size, simplelru.EvictCallback(onEvicted))

	if err != nil {
		return nil, err
	}

	w.lru = lru
	return w, nil
}

func NewWasmLDBCache(size int, db *leveldb.DB) (*WasmLDBCache, error) {
	w, err := NewWasmCache(size)
	if err != nil {
		return nil, err
	}
	w.db = db
	return w, nil
}

func (w *WasmLDBCache) SetDB(db *leveldb.DB) {
	w.db = db
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
		if w.db != nil {
			if value, err := w.db.Get(key.Bytes(), nil); err == nil {
				module := WasmModule{}
				buffer := bytes.NewReader(value)
				dec := gob.NewDecoder(buffer)
				if err := dec.Decode(&module); err != nil {
					log.Error("decode module err:", err)
					return nil, false
				}
				w.lru.Add(key, &module)
				return &module, true
			}
		}
		return nil, false
	}
	return value.(*WasmModule), ok
}

// Check if a key is in the cache, without updating the recent-ness
// or deleting it for being stale.
func (w *WasmLDBCache) Contains(key common.Address) bool {
	w.lock.RLock()
	defer w.lock.RUnlock()
	if !w.lru.Contains(key) {
		ok := false
		if w.db != nil {
			ok, _ = w.db.Has(key.Bytes(), nil)
		}
		return ok
	}
	return true
}

// Returns the key value (or undefined if not found) without updating
// the "recently used"-ness of the key.
func (w *WasmLDBCache) Peek(key common.Address) (*WasmModule, bool) {
	w.lock.Lock()
	defer w.lock.Unlock()
	value, ok := w.lru.Peek(key)
	if !ok {
		if w.db != nil {
			if value, err := w.db.Get(key.Bytes(), nil); err == nil {
				var module WasmModule
				buffer := bytes.NewReader(value)
				dec := gob.NewDecoder(buffer)
				dec.Decode(&module)
				return &module, true
			}
		}
		return nil, false
	}
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
	if w.db != nil {
		w.db.Delete(key.Bytes(), nil)
	}
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
