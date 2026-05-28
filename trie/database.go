// Copyright 2018 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package trie

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/trie/trienode"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/VictoriaMetrics/fastcache"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

const (
	DereferenceTimeLimit = 300 * time.Millisecond
)

var (
	memcacheCleanHitMeter   = metrics.NewRegisteredMeter("trie/memcache/clean/hit", nil)
	memcacheCleanMissMeter  = metrics.NewRegisteredMeter("trie/memcache/clean/miss", nil)
	memcacheCleanReadMeter  = metrics.NewRegisteredMeter("trie/memcache/clean/read", nil)
	memcacheCleanWriteMeter = metrics.NewRegisteredMeter("trie/memcache/clean/write", nil)

	memcacheDirtyHitMeter   = metrics.NewRegisteredMeter("trie/memcache/dirty/hit", nil)
	memcacheDirtyMissMeter  = metrics.NewRegisteredMeter("trie/memcache/dirty/miss", nil)
	memcacheDirtyReadMeter  = metrics.NewRegisteredMeter("trie/memcache/dirty/read", nil)
	memcacheDirtyWriteMeter = metrics.NewRegisteredMeter("trie/memcache/dirty/write", nil)

	memcacheFlushTimeTimer  = metrics.NewRegisteredResettingTimer("trie/memcache/flush/time", nil)
	memcacheFlushNodesMeter = metrics.NewRegisteredMeter("trie/memcache/flush/nodes", nil)
	memcacheFlushSizeMeter  = metrics.NewRegisteredMeter("trie/memcache/flush/size", nil)

	memcacheGCTimeTimer  = metrics.NewRegisteredResettingTimer("trie/memcache/gc/time", nil)
	memcacheGCNodesMeter = metrics.NewRegisteredMeter("trie/memcache/gc/nodes", nil)
	memcacheGCSizeMeter  = metrics.NewRegisteredMeter("trie/memcache/gc/size", nil)

	memcacheCommitTimeTimer  = metrics.NewRegisteredResettingTimer("trie/memcache/commit/time", nil)
	memcacheCommitNodesMeter = metrics.NewRegisteredMeter("trie/memcache/commit/nodes", nil)
	memcacheCommitSizeMeter  = metrics.NewRegisteredMeter("trie/memcache/commit/size", nil)
)

// childResolver defines the required method to decode the provided
// trie node and iterate the children on top.
type childResolver interface {
	forEach(node []byte, onChild func(common.Hash))
}

// Database is an intermediate write layer between the trie data structures and
// the disk database. The aim is to accumulate trie writes in-memory and only
// periodically flush a couple tries to disk, garbage collecting the remainder.
//
// Note, the trie Database is **not** thread safe in its mutations, but it **is**
// thread safe in providing individual, independent node access. The rationale
// behind this split design is to provide read access to RPC handlers and sync
// servers even while the trie is executing expensive garbage collection.
type Database struct {
	diskdb   ethdb.Database // Persistent storage for matured trie nodes
	resolver childResolver  // The handler to resolve children of nodes

	freshNodes map[common.Hash]struct{}

	cleans  *fastcache.Cache            // GC friendly memory cache of clean node RLPs
	dirties map[common.Hash]*cachedNode // Data and references relationships of dirty trie nodes
	oldest  common.Hash                 // Oldest tracked node, flush-list head
	newest  common.Hash                 // Newest tracked node, flush-list tail

	nodeVersion uint64
	useless     []map[string]struct{}

	gctime  time.Duration      // Time spent on garbage collection since last commit
	gcnodes uint64             // Nodes garbage collected since last commit
	gcsize  common.StorageSize // Data storage garbage collected since last commit

	flushtime  time.Duration      // Time spent on data flushing since last commit
	flushnodes uint64             // Nodes flushed since last commit
	flushsize  common.StorageSize // Data storage flushed since last commit

	dirtiesSize  common.StorageSize // Storage size of the dirty node cache (exc. metadata)
	childrenSize common.StorageSize // Storage size of the external children tracking
	preimages    *preimageStore     // The store for caching preimages

	lock sync.RWMutex
}

// cachedNode is all the information we know about a single cached trie node
// in the memory database write layer.
type cachedNode struct {
	node      []byte                   // Encoded node blob
	parents   uint32                   // Number of live nodes referencing this one
	external  map[common.Hash]struct{} // The set of external children
	flushPrev common.Hash              // Previous node in the flush-list
	flushNext common.Hash              // Next node in the flush-list

	version uint64 // The version number of this node
}

// cachedNodeSize is the raw size of a cachedNode data structure without any
// node data included. It's an approximate size, but should be a lot better
// than not counting them.
var cachedNodeSize = int(reflect.TypeOf(cachedNode{}).Size())

// forChildren invokes the callback for all the tracked children of this node,
// both the implicit ones from inside the node as well as the explicit ones
// from outside the node.
func (n *cachedNode) forChildren(resolver childResolver, onChild func(hash common.Hash)) {
	for child := range n.external {
		onChild(child)
	}
	resolver.forEach(n.node, onChild)
}

// Config defines all necessary options for database.
type Config struct {
	Cache     int    // Memory allowance (MB) to use for caching trie nodes in memory
	Journal   string // Journal of clean cache to survive node restarts
	Preimages bool   // Flag whether the preimage of trie key is recorded
}

// NewDatabase creates a new trie database to store ephemeral trie content before
// its written out to disk or garbage collected. No read cache is created, so all
// data retrievals will hit the underlying disk database.
func NewDatabase(diskdb ethdb.Database) *Database {
	return NewDatabaseWithConfig(diskdb, nil)
}

// NewDatabaseWithConfig creates a new trie database to store ephemeral trie content
// before its written out to disk or garbage collected. It also acts as a read cache
// for nodes loaded from disk.
func NewDatabaseWithConfig(diskdb ethdb.Database, config *Config) *Database {
	var cleans *fastcache.Cache
	if config != nil && config.Cache > 0 {
		if config.Journal == "" {
			cleans = fastcache.New(config.Cache * 1024 * 1024)
		} else {
			cleans = fastcache.LoadFromFileOrNew(config.Journal, config.Cache*1024*1024)
		}
	}
	var preimage *preimageStore
	if config != nil && config.Preimages {
		preimage = newPreimageStore(diskdb)
	}
	return &Database{
		diskdb:      diskdb,
		resolver:    mptResolver{},
		cleans:      cleans,
		dirties:     make(map[common.Hash]*cachedNode),
		nodeVersion: 0,
		freshNodes:  make(map[common.Hash]struct{}),
		preimages:   preimage,
	}
}

func (db *Database) NodeVersion() uint64 {
	return db.nodeVersion
}

func (db *Database) IncrVersion() {
	db.nodeVersion++
}

func (db *Database) insertFreshNode(hash common.Hash) {
	db.freshNodes[hash] = struct{}{}
}

func (db *Database) resetFreshNode() {
	db.freshNodes = make(map[common.Hash]struct{})
}

// insert inserts a simplified trie node into the memory database.
// All nodes inserted by this function will be reference tracked
// and in theory should only used for **trie nodes** insertion.
func (db *Database) insert(hash common.Hash, node []byte) {
	// If the node's already cached, skip
	if _, ok := db.dirties[hash]; ok {
		return
	}
	memcacheDirtyWriteMeter.Mark(int64(len(node)))

	// Create the cached entry for this node
	entry := &cachedNode{
		//node:      simplifyNode(node),
		node:      node,
		flushPrev: db.newest,
		version:   db.NodeVersion(),
	}
	//entry.forChildren(db.resolver, func(child common.Hash) {
	//	if c := db.dirties[child]; c != nil {
	//		c.parents++
	//	}
	//})
	db.dirties[hash] = entry

	// Update the flush-list endpoints
	if db.oldest == (common.Hash{}) {
		db.oldest, db.newest = hash, hash
	} else {
		db.dirties[db.newest].flushNext, db.newest = hash, hash
	}
	db.dirtiesSize += common.StorageSize(common.HashLength + len(node))
}

// Node retrieves an encoded cached trie node from memory. If it cannot be found
// cached, the method queries the persistent database for the content.
func (db *Database) Node(hash common.Hash) ([]byte, error) {
	// It doesn't make sense to retrieve the metaroot
	if hash == (common.Hash{}) {
		return nil, errors.New("not found")
	}
	// Retrieve the node from the clean cache if available
	if db.cleans != nil {
		if enc := db.cleans.Get(nil, hash[:]); enc != nil {
			memcacheCleanHitMeter.Mark(1)
			memcacheCleanReadMeter.Mark(int64(len(enc)))
			return enc, nil
		}
	}
	// Retrieve the node from the dirty cache if available
	db.lock.RLock()
	dirty := db.dirties[hash]
	db.lock.RUnlock()

	if dirty != nil {
		memcacheDirtyHitMeter.Mark(1)
		memcacheDirtyReadMeter.Mark(int64(len(dirty.node)))
		return dirty.node, nil
	}
	memcacheDirtyMissMeter.Mark(1)

	// Content unavailable in memory, attempt to retrieve from disk
	enc := rawdb.ReadLegacyTrieNode(db.diskdb, hash)
	if len(enc) != 0 {
		if db.cleans != nil {
			db.cleans.Set(hash[:], enc)
			memcacheCleanMissMeter.Mark(1)
			memcacheCleanWriteMeter.Mark(int64(len(enc)))
		}
		return enc, nil
	}
	return nil, errors.New("not found")
}

func (db *Database) NodeExistsInMemory(hash common.Hash) bool {
	db.lock.RLock()
	defer db.lock.RUnlock()

	_, ok := db.dirties[hash]
	return ok
}

// Nodes retrieves the hashes of all the nodes cached within the memory database.
// This method is extremely expensive and should only be used to validate internal
// states in test code.
func (db *Database) Nodes() []common.Hash {
	db.lock.RLock()
	defer db.lock.RUnlock()

	var hashes = make([]common.Hash, 0, len(db.dirties))
	for hash := range db.dirties {
		hashes = append(hashes, hash)
	}
	return hashes
}

// Reference adds a new reference from a parent node to a child node.
// This function is used to add reference between internal trie node
// and external node(e.g. storage trie root), all internal trie nodes
// are referenced together by database itself.
func (db *Database) Reference(child common.Hash, parent common.Hash) {
	db.lock.Lock()
	defer db.lock.Unlock()

	db.reference(child, parent)
}

// ReferenceVersion traverses down from the root node, with a version number for each node.
func (db *Database) ReferenceVersion(root common.Hash) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	start := time.Now()
	db.referenceVersion(root)
	if start.Add(400 * time.Millisecond).Before(time.Now()) {
		log.Warn("ReferenceVersion overtime", "root", root.String(), "duration", time.Since(start))
	}
}

// referenceVersion is the private locked version of referenceVersion.
func (db *Database) referenceVersion(hash common.Hash) {
	node, ok := db.dirties[hash]
	if !ok {
		return
	}
	node.forChildren(db.resolver, func(h common.Hash) {
		db.referenceVersion(h)
	})
	node.version = db.NodeVersion()
}

// reference is the private locked version of Reference.
func (db *Database) reference(child common.Hash, parent common.Hash) {
	// If the node does not exist, it's a node pulled from disk, skip
	node, ok := db.dirties[child]
	if !ok {
		return
	}
	// The reference is for state root, increase the reference counter.
	if parent == (common.Hash{}) {
		node.parents += 1
		return
	}
	// The reference is for external storage trie, don't duplicate if
	// the reference is already existent.
	if db.dirties[parent].external == nil {
		db.dirties[parent].external = make(map[common.Hash]struct{})
	}
	if _, ok := db.dirties[parent].external[child]; ok {
		return
	}
	node.parents++
	db.dirties[parent].external[child] = struct{}{}
	db.childrenSize += common.HashLength
}

// Dereference removes an existing reference from a root node.
func (db *Database) DereferenceDB(root common.Hash) {
	// Sanity check to ensure that the meta-root is not removed
	if root == (common.Hash{}) {
		log.Error("Attempted to dereference the trie cache meta root")
		return
	}

	db.lock.Lock()
	defer db.lock.Unlock()
	nodes, storage, start := len(db.dirties), db.dirtiesSize, time.Now()
	useless := make(map[string]struct{})
	clearFn := func(hash []byte) {
		useless[string(hash)] = struct{}{}
		if db.cleans != nil {
			db.cleans.Del(hash[:])
		}
	}

	db.dereference(root, clearFn, start)

	//if start.Add(400 * time.Millisecond).Before(time.Now()) {
	//	log.Warn("DereferenceDB overtime", "root", root.String(), "duration", time.Since(start))
	//}

	db.useless = append(db.useless, useless)
	db.gcnodes += uint64(nodes - len(db.dirties))
	db.gcsize += storage - db.dirtiesSize
	db.gctime += time.Since(start)

	memcacheGCTimeTimer.Update(time.Since(start))
	memcacheGCSizeMeter.Mark(int64(storage - db.dirtiesSize))
	memcacheGCNodesMeter.Mark(int64(nodes - len(db.dirties)))

	log.Debug("Dereferenced trie from memory database", "nodes", nodes-len(db.dirties), "size", storage-db.dirtiesSize, "time", time.Since(start),
		"gcnodes", db.gcnodes, "gcsize", db.gcsize, "gctime", db.gctime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)
}

//func (db *Database) uselessTotal() int {
//	db.lock.RLock()
//	defer db.lock.RUnlock()
//	sum := 0
//	for _, m := range db.useless {
//		sum += len(m)
//	}
//	return sum
//}

func (db *Database) UselessSize() int {
	return len(db.useless)
}

func (db *Database) ResetUseless() {
	db.lock.Lock()
	defer db.lock.Unlock()
	db.useless = nil
}

func (db *Database) UselessGC(num int) {
	db.lock.Lock()
	defer db.lock.Unlock()
	var (
		start = time.Now()
		total = 0
		batch = db.diskdb.NewBatch()
		size  = 0
	)

	for i, m := range db.useless {
		if total >= num {
			break
		}

		for k := range m {
			if db.dirties[common.BytesToHash([]byte(k))] == nil {
				batch.Delete([]byte(k))
			}
			if batch.ValueSize() > ethdb.IdealBatchSize {
				batch.Write()
				batch.Reset()
			}
			size++
		}

		db.useless[i] = nil
		total++
	}
	db.useless = db.useless[total:]
	batch.Write()
	log.Debug("UselessGC clean node", "size", size, "elapse", time.Since(start))
}

// Dereference removes an existing reference from a root node.
func (db *Database) Dereference(root common.Hash) {
	// Sanity check to ensure that the meta-root is not removed
	// PlatON 出块流程不匹配，本流程不适用
	//if root == (common.Hash{}) {
	//	log.Error("Attempted to dereference the trie cache meta root")
	//	return
	//}
	//db.lock.Lock()
	//defer db.lock.Unlock()
	//
	//cleanFn := func(hash []byte) {
	//	if db.cleans != nil {
	//		db.cleans.Del(hash)
	//	}
	//}
	//
	//nodes, storage, start := len(db.dirties), db.dirtiesSize, time.Now()
	//db.dereference(root, cleanFn, start)
	//
	//db.gcnodes += uint64(nodes - len(db.dirties))
	//db.gcsize += storage - db.dirtiesSize
	//db.gctime += time.Since(start)
	//
	//memcacheGCTimeTimer.Update(time.Since(start))
	//memcacheGCSizeMeter.Mark(int64(storage - db.dirtiesSize))
	//memcacheGCNodesMeter.Mark(int64(nodes - len(db.dirties)))
	//
	//log.Debug("Dereferenced trie from memory database", "nodes", nodes-len(db.dirties), "size", storage-db.dirtiesSize, "time", time.Since(start),
	//	"gcnodes", db.gcnodes, "gcsize", db.gcsize, "gctime", db.gctime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)
}

// dereference is the private locked version of Dereference.
func (db *Database) dereference(hash common.Hash, clearFn func([]byte), start time.Time) {
	if _, ok := db.freshNodes[hash]; ok {
		return
	}
	// If the child does not exist, it's a previously committed node.
	node, ok := db.dirties[hash]
	if !ok {
		return
	}
	if start.Add(DereferenceTimeLimit).Before(time.Now()) {
		log.Warn("DereferenceDB overtime, Interrupt the dereference", "duration", time.Since(start))
		return
	}
	if node.version < db.NodeVersion() {
		// Remove the node from the flush-list
		switch hash {
		case db.oldest:
			db.oldest = node.flushNext
			db.dirties[node.flushNext].flushPrev = common.Hash{}
		case db.newest:
			db.newest = node.flushPrev
			db.dirties[node.flushPrev].flushNext = common.Hash{}
		default:
			db.dirties[node.flushPrev].flushNext = node.flushNext
			db.dirties[node.flushNext].flushPrev = node.flushPrev
		}
		// Dereference all children and delete the node
		node.forChildren(db.resolver, func(child common.Hash) {
			db.dereference(child, clearFn, start)
		})
		delete(db.dirties, hash)

		if clearFn != nil {
			//// rawNode is contract code, only remove trie node
			//if _, ok := node.node.(rawNode); !ok {
			//	clearFn(hash.Bytes())
			//}
			clearFn(hash.Bytes())
		}
		db.dirtiesSize -= common.StorageSize(common.HashLength + len(node.node))
		if node.external != nil {
			db.childrenSize -= common.StorageSize(len(node.external) * common.HashLength)
		}
	}
}

func (db *Database) CapNode(limit common.StorageSize) {
	db.lock.RLock()
	nodes, storage, start := len(db.dirties), db.dirtiesSize, time.Now()
	size := db.dirtiesSize + common.StorageSize((len(db.dirties)-1)*2*common.HashLength)

	oldest := db.oldest
	for size > limit && oldest != (common.Hash{}) {
		// Fetch the oldest referenced node and push into the batch
		node := db.dirties[oldest]
		// Iterate to the next flush item, or abort if the size cap was achieved. Size
		// is the total size, including both the useful cached data (hash -> blob), as
		// well as the flushlist metadata (2*hash). When flushing items from the cache,
		// we need to reduce both.
		size -= common.StorageSize(3*common.HashLength + len(node.node))
		oldest = node.flushNext
	}
	db.lock.RUnlock()

	db.lock.Lock()
	defer db.lock.Unlock()
	for db.oldest != oldest {
		node := db.dirties[db.oldest]
		delete(db.dirties, db.oldest)
		db.oldest = node.flushNext

		db.dirtiesSize -= common.StorageSize(common.HashLength + len(node.node))
	}
	db.flushnodes += uint64(nodes - len(db.dirties))
	db.flushsize += storage - db.dirtiesSize
	db.flushtime += time.Since(start)

	memcacheFlushTimeTimer.Update(time.Since(start))
	memcacheFlushSizeMeter.Mark(int64(storage - db.dirtiesSize))
	memcacheFlushNodesMeter.Mark(int64(nodes - len(db.dirties)))

	log.Debug("Persisted nodes from memory database", "nodes", nodes-len(db.dirties), "size", storage-db.dirtiesSize, "time", time.Since(start),
		"flushnodes", db.flushnodes, "flushsize", db.flushsize, "flushtime", db.flushtime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)
}

// Cap iteratively flushes old but still referenced trie nodes until the total
// memory usage goes below the given threshold.
func (db *Database) Cap(limit common.StorageSize) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	// Create a database batch to flush persistent data out. It is important that
	// outside code doesn't see an inconsistent state (referenced data removed from
	// memory cache during commit but not yet in persistent storage). This is ensured
	// by only uncaching existing data when the database write finalizes.
	nodes, storage, start := len(db.dirties), db.dirtiesSize, time.Now()
	batch := db.diskdb.NewBatch()

	// db.dirtiesSize only contains the useful data in the cache, but when reporting
	// the total memory consumption, the maintenance metadata is also needed to be
	// counted.
	size := db.dirtiesSize + common.StorageSize(len(db.dirties)*cachedNodeSize)
	size += db.childrenSize

	// If the preimage cache got large enough, push to disk. If it's still small
	// leave for later to deduplicate writes.
	if db.preimages != nil {
		if err := db.preimages.commit(false); err != nil {
			return err
		}
	}
	// Keep committing nodes from the flush-list until we're below allowance
	oldest := db.oldest
	for size > limit && oldest != (common.Hash{}) {
		// Fetch the oldest referenced node and push into the batch
		node := db.dirties[oldest]
		rawdb.WriteLegacyTrieNode(batch, oldest, node.node)

		// If we exceeded the ideal batch size, commit and reset
		if batch.ValueSize() >= ethdb.IdealBatchSize {
			if err := batch.Write(); err != nil {
				log.Error("Failed to write flush list to disk", "err", err)
				return err
			}
			batch.Reset()
		}
		// Iterate to the next flush item, or abort if the size cap was achieved. Size
		// is the total size, including the useful cached data (hash -> blob), the
		// cache item metadata, as well as external children mappings.
		size -= common.StorageSize(common.HashLength + len(node.node) + cachedNodeSize)
		if node.external != nil {
			size -= common.StorageSize(len(node.external) * common.HashLength)
		}
		oldest = node.flushNext
	}
	// Flush out any remainder data from the last batch
	if err := batch.Write(); err != nil {
		log.Error("Failed to write flush list to disk", "err", err)
		return err
	}
	// Write successful, clear out the flushed data
	for db.oldest != oldest {
		node := db.dirties[db.oldest]
		delete(db.dirties, db.oldest)
		db.oldest = node.flushNext

		db.dirtiesSize -= common.StorageSize(common.HashLength + len(node.node))
		if node.external != nil {
			db.childrenSize -= common.StorageSize(len(node.external) * common.HashLength)
		}
	}
	if db.oldest != (common.Hash{}) {
		db.dirties[db.oldest].flushPrev = common.Hash{}
	}
	db.flushnodes += uint64(nodes - len(db.dirties))
	db.flushsize += storage - db.dirtiesSize
	db.flushtime += time.Since(start)

	memcacheFlushTimeTimer.Update(time.Since(start))
	memcacheFlushSizeMeter.Mark(int64(storage - db.dirtiesSize))
	memcacheFlushNodesMeter.Mark(int64(nodes - len(db.dirties)))

	log.Debug("Persisted nodes from memory database", "nodes", nodes-len(db.dirties), "size", storage-db.dirtiesSize, "time", time.Since(start),
		"flushnodes", db.flushnodes, "flushsize", db.flushsize, "flushtime", db.flushtime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)

	return nil
}

// Commit iterates over all the children of a particular node, writes them out
// to disk, forcefully tearing down all references in both directions. As a side
// effect, all pre-images accumulated up to this point are also written.
//
// Note, this method is a non-synchronized mutator. It is unsafe to call this
// concurrently with other mutators.
func (db *Database) Commit(node common.Hash, report bool, uncache bool) error {
	db.lock.Lock()
	defer db.lock.Unlock()
	// Create a database batch to flush persistent data out. It is important that
	// outside code doesn't see an inconsistent state (referenced data removed from
	// memory cache during commit but not yet in persistent storage). This is ensured
	// by only uncaching existing data when the database write finalizes.
	start := time.Now()
	batch := db.diskdb.NewBatch()

	// Move all of the accumulated preimages into a write batch
	if db.preimages != nil {
		if err := db.preimages.commit(true); err != nil {
			return err
		}
	}
	// Move the trie itself into the batch, flushing if enough data is accumulated
	nodes, storage := len(db.dirties), db.dirtiesSize

	uncacher := &cleaner{db, uncache}
	if err := db.commit(node, batch, uncacher); err != nil {
		log.Error("Failed to commit trie from trie database", "err", err)
		return err
	}
	// Trie mostly committed to disk, flush any batch leftovers
	if err := batch.Write(); err != nil {
		log.Error("Failed to write trie to disk", "err", err)
		return err
	}

	// Uncache any leftovers in the last batch
	if err := batch.Replay(uncacher); err != nil {
		return err
	}
	batch.Reset()

	// Reset the storage counters and bumpd metrics
	db.resetFreshNode()

	memcacheCommitTimeTimer.Update(time.Since(start))
	memcacheCommitSizeMeter.Mark(int64(storage - db.dirtiesSize))
	memcacheCommitNodesMeter.Mark(int64(nodes - len(db.dirties)))

	logger := log.Info
	if !report {
		logger = log.Debug
	}
	logger("Persisted trie from memory database", "nodes", nodes-len(db.dirties)+int(db.flushnodes), "size", storage-db.dirtiesSize+db.flushsize, "time", time.Since(start)+db.flushtime,
		"gcnodes", db.gcnodes, "gcsize", db.gcsize, "gctime", db.gctime, "livenodes", len(db.dirties), "livesize", db.dirtiesSize)

	// Reset the garbage collection statistics
	db.gcnodes, db.gcsize, db.gctime = 0, 0, 0
	db.flushnodes, db.flushsize, db.flushtime = 0, 0, 0

	return nil
}

// commit is the private locked version of Commit.
func (db *Database) commit(hash common.Hash, batch ethdb.Batch, uncacher *cleaner) error {
	// If the node does not exist, it's a previously committed node
	_, ok := db.freshNodes[hash]
	if !ok {
		return nil
	}
	node, ok := db.dirties[hash]
	if !ok {
		return nil
	}
	var err error
	node.forChildren(db.resolver, func(child common.Hash) {
		if err == nil {
			err = db.commit(child, batch, uncacher)
		}
	})
	if err != nil {
		return err
	}
	// If we've reached an optimal batch size, commit and start over
	rawdb.WriteLegacyTrieNode(batch, hash, node.node)

	if batch.ValueSize() >= ethdb.IdealBatchSize {
		if err := batch.Write(); err != nil {
			return err
		}
		err := batch.Replay(uncacher)
		if err != nil {
			return err
		}
		batch.Reset()
	}
	return nil
}

// cleaner is a database batch replayer that takes a batch of write operations
// and cleans up the trie database from anything written to disk.
type cleaner struct {
	db      *Database
	uncache bool
}

// Put reacts to database writes and implements dirty data uncaching. This is the
// post-processing step of a commit operation where the already persisted trie is
// removed from the dirty cache and moved into the clean cache. The reason behind
// the two-phase commit is to ensure data availability while moving from memory
// to disk.
func (c *cleaner) Put(key []byte, rlp []byte) error {
	hash := common.BytesToHash(key)

	if c.uncache {
		// If the node does not exist, we're done on this path
		node, ok := c.db.dirties[hash]
		if !ok {
			return nil
		}
		// Node still exists, remove it from the flush-list
		switch hash {
		case c.db.oldest:
			c.db.oldest = node.flushNext
			if node.flushNext != (common.Hash{}) {
				c.db.dirties[node.flushNext].flushPrev = common.Hash{}
			}
		case c.db.newest:
			c.db.newest = node.flushPrev
			if node.flushPrev != (common.Hash{}) {
				c.db.dirties[node.flushPrev].flushNext = common.Hash{}
			}
		default:
			c.db.dirties[node.flushPrev].flushNext = node.flushNext
			c.db.dirties[node.flushNext].flushPrev = node.flushPrev
		}
		// Remove the node from the dirty cache
		delete(c.db.dirties, hash)
		c.db.dirtiesSize -= common.StorageSize(common.HashLength + len(node.node))
		if node.external != nil {
			c.db.childrenSize -= common.StorageSize(len(node.external) * common.HashLength)
		}
	}
	// Move the flushed node into the clean cache to prevent insta-reloads
	if c.db.cleans != nil {
		c.db.cleans.Set(hash[:], rlp)
		memcacheCleanWriteMeter.Mark(int64(len(rlp)))
	}
	return nil

}

func (c *cleaner) Delete(key []byte) error {
	panic("not implemented")
}

// Update inserts the dirty nodes in provided nodeset into database and
// link the account trie with multiple storage tries if necessary.
func (db *Database) Update(nodes *MergedNodeSet) error {
	db.lock.Lock()
	defer db.lock.Unlock()

	// Insert dirty nodes into the database. In the same tree, it must be
	// ensured that children are inserted first, then parent so that children
	// can be linked with their parent correctly.
	//
	// Note, the storage tries must be flushed before the account trie to
	// retain the invariant that children go into the dirty cache first.
	var order []common.Hash
	for owner := range nodes.sets {
		if owner == (common.Hash{}) {
			continue
		}
		order = append(order, owner)
	}
	if _, ok := nodes.sets[common.Hash{}]; ok {
		order = append(order, common.Hash{})
	}
	//for _, owner := range order {
	//	subset := nodes.sets[owner]
	//	for _, path := range subset.updates.order {
	//		n, ok := subset.updates.nodes[path]
	//		if !ok {
	//			return fmt.Errorf("missing node %x %v", owner, path)
	//		}
	//		db.insert(n.hash, int(n.size), n.node)
	//		db.insertFreshNode(n.hash)
	//	}
	//}
	for _, owner := range order {
		subset := nodes.sets[owner]
		subset.forEachWithOrder(func(path string, n *trienode.Node) {
			if n.IsDeleted() {
				return // ignore deletion
			}
			db.insert(n.Hash, n.Blob)
			db.insertFreshNode(n.Hash)
		})
	}
	// Link up the account trie and storage trie if the node points
	// to an account trie leaf.
	if set, present := nodes.sets[common.Hash{}]; present {
		for _, n := range set.leaves {
			var account types.StateAccount
			if err := rlp.DecodeBytes(n.blob, &account); err != nil {
				return err
			}
			if account.Root != types.EmptyRootHash {
				//db.Reference(account.Root, n.parent)
				db.reference(account.Root, n.parent)
			}
		}
	}
	return nil
}

// Size returns the current storage size of the memory cache in front of the
// persistent database layer.
func (db *Database) Size() (common.StorageSize, common.StorageSize) {
	db.lock.RLock()
	defer db.lock.RUnlock()

	// db.dirtiesSize only contains the useful data in the cache, but when reporting
	// the total memory consumption, the maintenance metadata is also needed to be
	// counted.
	var metadataSize = common.StorageSize(len(db.dirties) * cachedNodeSize)
	var preimageSize common.StorageSize
	if db.preimages != nil {
		preimageSize = db.preimages.size()
	}
	return db.dirtiesSize + db.childrenSize + metadataSize, preimageSize
}

// GetReader retrieves a node reader belonging to the given state root.
func (db *Database) GetReader(root common.Hash) Reader {
	return newHashReader(db)
}

// hashReader is reader of hashDatabase which implements the Reader interface.
type hashReader struct {
	db *Database
}

// newHashReader initializes the hash reader.
func newHashReader(db *Database) *hashReader {
	return &hashReader{db: db}
}

// Node retrieves the RLP-encoded trie node blob with the given node hash.
// No error will be returned if the node is not found.
func (reader *hashReader) Node(_ common.Hash, _ []byte, hash common.Hash) ([]byte, error) {
	blob, _ := reader.db.Node(hash)
	return blob, nil
}

func (reader *hashReader) NodeExistsInMemory(_ common.Hash, _ []byte, hash common.Hash) bool {
	return reader.db.NodeExistsInMemory(hash)
}

// saveCache saves clean state cache to given directory path
// using specified CPU cores.
func (db *Database) saveCache(dir string, threads int) error {
	if db.cleans == nil {
		return nil
	}
	log.Info("Writing clean trie cache to disk", "path", dir, "threads", threads)

	start := time.Now()
	err := db.cleans.SaveToFileConcurrent(dir, threads)
	if err != nil {
		log.Error("Failed to persist clean trie cache", "error", err)
		return err
	}
	log.Info("Persisted the clean trie cache", "path", dir, "elapsed", common.PrettyDuration(time.Since(start)))
	return nil
}

// SaveCache atomically saves fast cache data to the given dir using all
// available CPU cores.
func (db *Database) SaveCache(dir string) error {
	return db.saveCache(dir, runtime.GOMAXPROCS(0))
}

// SaveCachePeriodically atomically saves fast cache data to the given dir with
// the specified interval. All dump operation will only use a single CPU core.
func (db *Database) SaveCachePeriodically(dir string, interval time.Duration, stopCh <-chan struct{}) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			db.saveCache(dir, 1)
		case <-stopCh:
			return
		}
	}
}

// CommitPreimages flushes the dangling preimages to disk. It is meant to be
// called when closing the blockchain object, so that preimages are persisted
// to the database.
func (db *Database) CommitPreimages() error {
	db.lock.Lock()
	defer db.lock.Unlock()

	if db.preimages == nil {
		return nil
	}
	return db.preimages.commit(true)
}

// Scheme returns the node scheme used in the database.
func (db *Database) Scheme() string {
	return rawdb.HashScheme
}
