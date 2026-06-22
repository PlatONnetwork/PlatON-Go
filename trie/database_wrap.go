// Copyright 2022 The go-ethereum Authors
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
	"runtime"
	"time"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/trie/triedb/hashdb"
	"github.com/PlatONnetwork/PlatON-Go/trie/trienode"
)

// Config defines all necessary options for database.
type Config struct {
	Cache     int    // Memory allowance (MB) to use for caching trie nodes in memory
	Journal   string // Journal of clean cache to survive node restarts
	Preimages bool   // Flag whether the preimage of trie key is recorded
}

// backend defines the methods needed to access/update trie nodes in different
// state scheme.
type backend interface {
	Scheme() string
	Initialized(genesisRoot common.Hash) bool
	Size() common.StorageSize
	Update(root common.Hash, parent common.Hash, block uint64, nodes *trienode.MergedNodeSet) error
	Commit(root common.Hash, report bool, uncache bool) error
	Close() error
}

// Database is the wrapper of the underlying backend which is shared by different
// types of node backend as an entrypoint.
type Database struct {
	config    *Config
	diskdb    ethdb.Database
	cleans    *fastcache.Cache
	preimages *preimageStore
	backend   backend
	hashdb    *hashdb.Database
}

func prepare(diskdb ethdb.Database, config *Config) *Database {
	var cleans *fastcache.Cache
	if config != nil && config.Cache > 0 {
		if config.Journal == "" {
			cleans = fastcache.New(config.Cache * 1024 * 1024)
		} else {
			cleans = fastcache.LoadFromFileOrNew(config.Journal, config.Cache*1024*1024)
		}
	}
	var preimages *preimageStore
	if config != nil && config.Preimages {
		preimages = newPreimageStore(diskdb)
	}
	return &Database{
		config:    config,
		diskdb:    diskdb,
		cleans:    cleans,
		preimages: preimages,
	}
}

// NewDatabase creates a new trie database to store ephemeral trie content before
// its written out to disk or garbage collected.
func NewDatabase(diskdb ethdb.Database) *Database {
	return NewDatabaseWithConfig(diskdb, nil)
}

// NewDatabaseWithConfig creates a new trie database with provided configs.
func NewDatabaseWithConfig(diskdb ethdb.Database, config *Config) *Database {
	db := prepare(diskdb, config)
	hdb := hashdb.New(diskdb, db.cleans, mptChildResolver{})
	db.hashdb = hdb
	db.backend = hdb
	return db
}

// mptChildResolver adapts mptResolver to hashdb.ChildResolver.
type mptChildResolver struct{}

func (mptChildResolver) ForEach(node []byte, onChild func(common.Hash)) {
	mptResolver{}.forEach(node, onChild)
}

// Reader returns a reader for accessing all trie nodes with provided state root.
func (db *Database) Reader(blockRoot common.Hash) Reader {
	return db.hashdb.Reader(blockRoot)
}

// GetReader is an alias of Reader for backward compatibility.
func (db *Database) GetReader(blockRoot common.Hash) Reader {
	return db.Reader(blockRoot)
}

// Update performs a state transition by committing dirty nodes contained in the
// given set in order to update state from the specified parent to the specified root.
func (db *Database) Update(root common.Hash, parent common.Hash, block uint64, nodes *trienode.MergedNodeSet) error {
	if db.preimages != nil {
		db.preimages.commit(false)
	}
	return db.backend.Update(root, parent, block, nodes)
}

// Commit iterates over all the children of a particular node, writes them out to disk.
func (db *Database) Commit(root common.Hash, report bool, uncache bool) error {
	if db.preimages != nil {
		db.preimages.commit(true)
	}
	return db.backend.Commit(root, report, uncache)
}

// Size returns the storage size of dirty trie nodes and cached preimages.
func (db *Database) Size() (common.StorageSize, common.StorageSize) {
	var (
		storages  common.StorageSize
		preimages common.StorageSize
	)
	storages = db.backend.Size()
	if db.preimages != nil {
		preimages = db.preimages.size()
	}
	return storages, preimages
}

// Initialized returns an indicator if the state data is already initialized.
func (db *Database) Initialized(genesisRoot common.Hash) bool {
	return db.backend.Initialized(genesisRoot)
}

// Scheme returns the node scheme used in the database.
func (db *Database) Scheme() string {
	return db.backend.Scheme()
}

// Close flushes dangling preimages and closes the trie database.
func (db *Database) Close() error {
	if db.preimages != nil {
		db.preimages.commit(true)
	}
	return db.backend.Close()
}

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

// SaveCache atomically saves fast cache data to the given dir.
func (db *Database) SaveCache(dir string) error {
	return db.saveCache(dir, runtime.GOMAXPROCS(0))
}

// SaveCachePeriodically atomically saves fast cache data periodically.
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

// CommitPreimages flushes dangling preimages to disk.
func (db *Database) CommitPreimages() error {
	if db.preimages == nil {
		return nil
	}
	return db.preimages.commit(true)
}

// Cap iteratively flushes old but still referenced trie nodes.
func (db *Database) Cap(limit common.StorageSize) error {
	if db.preimages != nil {
		db.preimages.commit(false)
	}
	return db.hashdb.Cap(limit)
}

// Reference adds a new reference from a parent node to a child node.
func (db *Database) Reference(root common.Hash, parent common.Hash) error {
	db.hashdb.Reference(root, parent)
	return nil
}

// Dereference removes an existing reference from a root node.
func (db *Database) Dereference(root common.Hash) error {
	db.hashdb.Dereference(root)
	return nil
}

// Node retrieves the rlp-encoded node blob with provided node hash.
func (db *Database) Node(hash common.Hash) ([]byte, error) {
	return db.hashdb.Node(hash)
}

// Nodes retrieves the hashes of all nodes cached within the memory database.
func (db *Database) Nodes() []common.Hash {
	return db.hashdb.Nodes()
}

// NodeExistsInMemory reports whether the node exists in memory.
func (db *Database) NodeExistsInMemory(hash common.Hash) bool {
	return db.hashdb.NodeExistsInMemory(hash)
}

// NodeVersion returns the current node version.
func (db *Database) NodeVersion() uint64 {
	return db.hashdb.NodeVersion()
}

// IncrVersion increments the node version.
func (db *Database) IncrVersion() {
	db.hashdb.IncrVersion()
}

// ReferenceVersion traverses down from the root node with version tagging.
func (db *Database) ReferenceVersion(root common.Hash) {
	db.hashdb.ReferenceVersion(root)
}

// DereferenceDB removes reference and records useless nodes for GC.
func (db *Database) DereferenceDB(root common.Hash) {
	db.hashdb.DereferenceDB(root)
}

// UselessSize returns the number of useless node batches.
func (db *Database) UselessSize() int {
	return db.hashdb.UselessSize()
}

// ResetUseless clears useless node tracking.
func (db *Database) ResetUseless() {
	db.hashdb.ResetUseless()
}

// UselessGC garbage-collects useless nodes from disk.
func (db *Database) UselessGC(num int) {
	db.hashdb.UselessGC(num)
}

// CapNode caps in-memory nodes without flushing to disk.
func (db *Database) CapNode(limit common.StorageSize) {
	db.hashdb.CapNode(limit)
}

// DiskDB returns the persistent database handle.
func (db *Database) DiskDB() ethdb.Database {
	return db.diskdb
}

// Dirties returns dirty trie nodes. For testing only.
func (db *Database) Dirties() map[common.Hash]*hashdb.CachedNode {
	return db.hashdb.Dirties()
}

// GetDirty returns a dirty trie node. For testing only.
func (db *Database) GetDirty(hash common.Hash) *hashdb.CachedNode {
	return db.hashdb.GetDirty(hash)
}

// DeleteDirty removes a dirty trie node. For testing only.
func (db *Database) DeleteDirty(hash common.Hash) {
	db.hashdb.DeleteDirty(hash)
}

// SetDirty restores a dirty trie node. For testing only.
func (db *Database) SetDirty(hash common.Hash, node *hashdb.CachedNode) {
	db.hashdb.SetDirty(hash, node)
}
