// Copyright 2014 The go-ethereum Authors
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

// Package core implements the Ethereum consensus protocol.
package core

import (
	"errors"
	"fmt"
	"io"
	mrand "math/rand"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/internal/syncx"

	"github.com/PlatONnetwork/PlatON-Go/core/state/snapshot"

	lru "github.com/hashicorp/golang-lru"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mclock"
	"github.com/PlatONnetwork/PlatON-Go/common/prque"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/trie"
)

var (
	headBlockGauge     = metrics.NewRegisteredGauge("chain/head/block", nil)
	headHeaderGauge    = metrics.NewRegisteredGauge("chain/head/header", nil)
	headFastBlockGauge = metrics.NewRegisteredGauge("chain/head/receipt", nil)

	accountReadTimer   = metrics.NewRegisteredTimer("chain/account/reads", nil)
	accountHashTimer   = metrics.NewRegisteredTimer("chain/account/hashes", nil)
	accountUpdateTimer = metrics.NewRegisteredTimer("chain/account/updates", nil)
	accountCommitTimer = metrics.NewRegisteredTimer("chain/account/commits", nil)

	storageReadTimer   = metrics.NewRegisteredTimer("chain/storage/reads", nil)
	storageHashTimer   = metrics.NewRegisteredTimer("chain/storage/hashes", nil)
	storageUpdateTimer = metrics.NewRegisteredTimer("chain/storage/updates", nil)
	storageCommitTimer = metrics.NewRegisteredTimer("chain/storage/commits", nil)

	snapshotAccountReadTimer = metrics.NewRegisteredTimer("chain/snapshot/account/reads", nil)
	snapshotStorageReadTimer = metrics.NewRegisteredTimer("chain/snapshot/storage/reads", nil)
	snapshotCommitTimer      = metrics.NewRegisteredTimer("chain/snapshot/commits", nil)

	blockInsertTimer     = metrics.NewRegisteredTimer("chain/inserts", nil)
	blockValidationTimer = metrics.NewRegisteredTimer("chain/validation", nil)
	//blockExecutionTimer  = metrics.NewRegisteredTimer("chain/execution", nil)
	//blockWriteTimer = metrics.NewRegisteredTimer("chain/write", nil)

	blockReorgMeter         = metrics.NewRegisteredMeter("chain/reorg/executes", nil)
	blockReorgAddMeter      = metrics.NewRegisteredMeter("chain/reorg/add", nil)
	blockReorgDropMeter     = metrics.NewRegisteredMeter("chain/reorg/drop", nil)
	blockReorgInvalidatedTx = metrics.NewRegisteredMeter("chain/reorg/invalidTx", nil)

	ErrNoGenesis          = errors.New("Genesis not found in chain")
	defaultCapNodePercent = common.StorageSize(1) / 4

	errInsertionInterrupted = errors.New("insertion is interrupted")
	errChainStopped         = errors.New("blockchain is stopped")
)

// CacheConfig contains the configuration values for the trie caching/pruning
// that's resident in a blockchain.
type CacheConfig struct {
	Disabled bool // Whether to disable trie write caching (archive node)

	TrieCleanLimit     int           // Memory allowance (MB) to use for caching trie nodes in memory
	TrieCleanJournal   string        // Disk journal for saving clean cache entries.
	TrieCleanRejournal time.Duration // Time interval to dump clean cache to disk periodically
	TrieDirtyLimit     int           // Memory limit (MB) at which to flush the current in-memory trie to disk
	TrieTimeLimit      time.Duration // Time limit after which to flush the current in-memory trie to disk
	SnapshotLimit      int           // Memory allowance (MB) to use for caching snapshot entries in memory

	SnapshotWait bool // Wait for snapshot construction on startup. TODO(karalabe): This is a dirty hack for testing, nuke it

	Preimages bool // Whether to store preimage of trie key to the disk

	BodyCacheLimit  int
	BlockCacheLimit int
	MaxFutureBlocks int
	BadBlockLimit   int
	TriesInMemory   int

	DBDisabledGC common.AtomicBool // Whether to disable database garbage collection
	DBGCInterval uint64            // Block interval for database garbage collection
	DBGCTimeout  time.Duration
	DBGCMpt      bool
	DBGCBlock    int
}

// defaultCacheConfig are the default caching values if none are specified by the
// user (also used during testing).
var defaultCacheConfig = &CacheConfig{
	TrieCleanLimit:  512,
	TrieDirtyLimit:  256 * 1024 * 1024,
	TrieTimeLimit:   5 * time.Minute,
	SnapshotLimit:   256,
	SnapshotWait:    true,
	BodyCacheLimit:  256,
	BlockCacheLimit: 256,
	MaxFutureBlocks: 256,
	BadBlockLimit:   10,
	TriesInMemory:   128,
	DBGCInterval:    86400,
	DBGCTimeout:     time.Minute,
	Preimages:       true,
}

// mining related configuration
type MiningConfig struct {
	MiningLogAtDepth       uint
	TxChanSize             int
	ChainHeadChanSize      int
	ChainSideChanSize      int
	ResultQueueSize        int
	ResubmitAdjustChanSize int
	MinRecommitInterval    time.Duration
	MaxRecommitInterval    time.Duration
	IntervalAdjustRatio    float64
	IntervalAdjustBias     float64
	StaleThreshold         uint64
	DefaultCommitRatio     float64
}

const (
	receiptsCacheLimit = 32
	txLookupCacheLimit = 1024

	// BlockChainVersion ensures that an incompatible database forces a resync from scratch.
	//
	// Changelog:
	//
	// - Version 4
	//   The following incompatible database changes were added:
	//   * the `BlockNumber`, `TxHash`, `TxIndex`, `BlockHash` and `Index` fields of log are deleted
	//   * the `Bloom` field of receipt is deleted
	//   * the `BlockIndex` and `TxIndex` fields of txlookup are deleted
	// - Version 5
	//  The following incompatible database changes were added:
	//    * the `TxHash`, `GasCost`, and `ContractAddress` fields are no longer stored for a receipt
	//    * the `TxHash`, `GasCost`, and `ContractAddress` fields are computed by looking up the
	//      receipts' corresponding block
	// - Version 6
	//  The following incompatible database changes were added:
	//    * Transaction lookup information stores the corresponding block number instead of block hash
	// - Version 7
	//  The following incompatible database changes were added:
	//    * Use freezer as the ancient database to maintain all ancient data
	// - Version 8
	//  The following incompatible database changes were added:
	//    * New scheme for contract code in order to separate the codes and trie nodes
	BlockChainVersion uint64 = 8
)

// BlockChain represents the canonical chain given a database with a genesis
// block. The Blockchain manages chain imports, reverts, chain reorganisations.
//
// Importing blocks in to the block chain happens according to the set of rules
// defined by the two stage Validator. Processing of blocks is done using the
// Processor which processes the included transaction. The validation of the state
// is done in the second part of the Validator. Failing results in aborting of
// the import.
//
// The BlockChain also helps in returning blocks from **any** chain included
// in the database as well as blocks that represents the canonical chain. It's
// important to note that GetBlock can return any block and does not need to be
// included in the canonical one where as GetBlockByNumber always represents the
// canonical chain.
type BlockChain struct {
	chainConfig *params.ChainConfig // Chain & network configuration
	cacheConfig *CacheConfig        // Cache configuration for pruning

	db     ethdb.Database // Low level persistent database to store final content in
	snaps  *snapshot.Tree // Snapshot tree for fast trie leaf access
	triegc *prque.Prque   // Priority queue mapping block numbers to tries to gc
	gcproc time.Duration  // Accumulates canonical block processing for trie dumping

	// txLookupLimit is the maximum number of blocks from head whose tx indices
	// are reserved:
	//  * 0:   means no limit and regenerate any missing indexes
	//  * N:   means N block limit [HEAD-N+1, HEAD] and delete extra indexes
	//  * nil: disable tx reindexer/deleter, but still index new blocks
	txLookupLimit uint64

	hc            *HeaderChain
	rmLogsFeed    event.Feed
	chainFeed     event.Feed
	chainSideFeed event.Feed
	chainHeadFeed event.Feed
	logsFeed      event.Feed
	scope         event.SubscriptionScope
	genesisBlock  *types.Block

	BlockFeed        event.Feed
	BlockExecuteFeed event.Feed

	// This mutex synchronizes chain write operations.
	// Readers don't need to take it, they can just read the database.
	chainmu *syncx.ClosableMutex
	procmu  sync.RWMutex // block processor lock

	checkpoint       int          // checkpoint counts towards the new checkpoint
	currentBlock     atomic.Value // Current head of the block chain
	currentFastBlock atomic.Value // Current head of the fast-sync chain (may be above the block chain!)

	stateCache    state.Database // State database to reuse between imports (contains state cache)
	bodyCache     *lru.Cache     // Cache for the most recent block bodies
	bodyRLPCache  *lru.Cache     // Cache for the most recent block bodies in RLP encoded format
	receiptsCache *lru.Cache     // Cache for the most recent receipts per block
	blockCache    *lru.Cache     // Cache for the most recent entire blocks
	txLookupCache *lru.Cache     // Cache for the most recent transaction lookup data.
	futureBlocks  *lru.Cache     // future blocks are blocks added for later processing

	wg            sync.WaitGroup //
	quit          chan struct{}  // shutdown signal, closed in Stop.
	running       int32          // 0 if chain is running, 1 when stopped
	procInterrupt int32          // interrupt signaler for block processing
	executeWG     sync.WaitGroup // execute block processing wait group for shutting dow

	engine    consensus.Engine
	processor Processor // block processor interface
	validator Validator // block and state validator interface
	vmConfig  vm.Config

	shouldPreserve func(*types.Block) bool // Function used to determine whether should preserve the given block.

	cleaner *Cleaner
}

// NewBlockChain returns a fully initialised block chain using information
// available in the database. It initialises the default Ethereum Validator and
// Processor.
func NewBlockChain(db ethdb.Database, cacheConfig *CacheConfig, chainConfig *params.ChainConfig, engine consensus.Engine, vmConfig vm.Config, shouldPreserve func(block *types.Block) bool, txLookupLimit *uint64) (*BlockChain, error) {
	if cacheConfig == nil {
		cacheConfig = defaultCacheConfig
	}
	bodyCache, _ := lru.New(cacheConfig.BodyCacheLimit)
	bodyRLPCache, _ := lru.New(cacheConfig.BodyCacheLimit)
	receiptsCache, _ := lru.New(receiptsCacheLimit)
	blockCache, _ := lru.New(cacheConfig.BlockCacheLimit)
	txLookupCache, _ := lru.New(txLookupCacheLimit)
	futureBlocks, _ := lru.New(cacheConfig.MaxFutureBlocks)

	bc := &BlockChain{
		chainConfig: chainConfig,
		cacheConfig: cacheConfig,
		db:          db,
		triegc:      prque.New(nil),
		stateCache: state.NewDatabaseWithConfig(db, &trie.Config{
			Cache:     cacheConfig.TrieCleanLimit,
			Journal:   cacheConfig.TrieCleanJournal,
			Preimages: cacheConfig.Preimages,
		}),
		quit:           make(chan struct{}),
		chainmu:        syncx.NewClosableMutex(),
		shouldPreserve: shouldPreserve,
		bodyCache:      bodyCache,
		bodyRLPCache:   bodyRLPCache,
		receiptsCache:  receiptsCache,
		blockCache:     blockCache,
		txLookupCache:  txLookupCache,
		futureBlocks:   futureBlocks,
		engine:         engine,
		vmConfig:       vmConfig,
	}

	bc.SetValidator(NewBlockValidator(chainConfig, bc, engine))
	bc.SetProcessor(NewParallelStateProcessor(chainConfig, bc, engine))
	//bc.SetProcessor(NewStateProcessor(chainConfig, bc, engine))

	var err error
	bc.hc, err = NewHeaderChain(db, chainConfig, engine, bc.insertStopped)
	if err != nil {
		return nil, err
	}
	bc.genesisBlock = bc.GetBlockByNumber(0)
	if bc.genesisBlock == nil {
		return nil, ErrNoGenesis
	}

	var nilBlock *types.Block
	bc.currentBlock.Store(nilBlock)
	bc.currentFastBlock.Store(nilBlock)

	// Initialize the chain with ancient data if it isn't empty.
	var txIndexBlock uint64

	if bc.empty() {
		rawdb.InitDatabaseFromFreezer(bc.db)
		// If ancient database is not empty, reconstruct all missing
		// indices in the background.
		frozen, _ := bc.db.Ancients()
		if frozen > 0 {
			txIndexBlock = frozen
		}
	}

	if err := bc.loadLastState(); err != nil {
		return nil, err
	}

	// Make sure the state associated with the block is available
	head := bc.CurrentBlock()
	if _, err := state.New(head.Root(), bc.stateCache, bc.snaps); err != nil {
		// Head state is missing, before the state recovery, find out the
		// disk layer point of snapshot(if it's enabled). Make sure the
		// rewound point is lower than disk layer.
		var diskRoot common.Hash
		if bc.cacheConfig.SnapshotLimit > 0 {
			diskRoot = rawdb.ReadSnapshotRoot(bc.db)
		}
		if diskRoot != (common.Hash{}) {
			log.Warn("Head state missing, repairing", "number", head.Number(), "hash", head.Hash(), "snaproot", diskRoot)

			snapDisk, err := bc.setHeadBeyondRoot(head.NumberU64(), diskRoot, true)
			if err != nil {
				return nil, err
			}
			// Chain rewound, persist old snapshot number to indicate recovery procedure
			if snapDisk != 0 {
				rawdb.WriteSnapshotRecoveryNumber(bc.db, snapDisk)
			}
		} else {
			log.Warn("Head state missing, repairing", "number", head.Number(), "hash", head.Hash())
			if _, err := bc.setHeadBeyondRoot(head.NumberU64(), common.Hash{}, true); err != nil {
				return nil, err
			}
		}
	}

	// Ensure that a previous crash in SetHead doesn't leave extra ancients
	if frozen, err := bc.db.Ancients(); err == nil && frozen > 0 {
		var (
			needRewind bool
			low        uint64
		)
		// The head full block may be rolled back to a very low height due to
		// blockchain repair. If the head full block is even lower than the ancient
		// chain, truncate the ancient store.
		fullBlock := bc.CurrentBlock()
		if fullBlock != nil && fullBlock.Hash() != bc.genesisBlock.Hash() && fullBlock.NumberU64() < frozen-1 {
			needRewind = true
			low = fullBlock.NumberU64()
		}
		// In fast sync, it may happen that ancient data has been written to the
		// ancient store, but the LastFastBlock has not been updated, truncate the
		// extra data here.
		fastBlock := bc.CurrentFastBlock()
		if fastBlock != nil && fastBlock.NumberU64() < frozen-1 {
			needRewind = true
			if fastBlock.NumberU64() < low || low == 0 {
				low = fastBlock.NumberU64()
			}
		}
		if needRewind {
			return nil, errors.New("needRewind due to data error,please clean your data")
			/*var hashes []common.Hash
			previous := bc.CurrentHeader().Number.Uint64()
			for i := low + 1; i <= bc.CurrentHeader().Number.Uint64(); i++ {
				hashes = append(hashes, rawdb.ReadCanonicalHash(bc.db, i))
			}
			bc.Rollback(hashes)
			log.Warn("Truncate ancient chain", "from", previous, "to", low)*/
		}
	}

	log.Debug("DB config", "DBDisabledGC", bc.cacheConfig.DBDisabledGC, "DBGCInterval", bc.cacheConfig.DBGCInterval, "DBGCTimeout", bc.cacheConfig.DBGCTimeout, "DBGCMpt", bc.cacheConfig.DBGCMpt)
	bc.cleaner = NewCleaner(bc, bc.cacheConfig.DBGCInterval, bc.cacheConfig.DBGCTimeout, bc.cacheConfig.DBGCMpt)

	// Load any existing snapshot, regenerating it if loading failed
	if bc.cacheConfig.SnapshotLimit > 0 {
		// If the chain was rewound past the snapshot persistent layer (causing
		// a recovery block number to be persisted to disk), check if we're still
		// in recovery mode and in that case, don't invalidate the snapshot on a
		// head mismatch.
		var recover bool

		head := bc.CurrentBlock()
		if layer := rawdb.ReadSnapshotRecoveryNumber(bc.db); layer != nil && *layer > head.NumberU64() {
			log.Warn("Enabling snapshot recovery", "chainhead", head.NumberU64(), "diskbase", *layer)
			recover = true
		}
		bc.snaps, _ = snapshot.New(bc.db, bc.stateCache.TrieDB(), bc.cacheConfig.SnapshotLimit, head.Root(), !bc.cacheConfig.SnapshotWait, true, recover)
	}

	// Start future block processor.
	bc.wg.Add(1)
	go bc.futureBlocksLoop()

	// Start tx indexer/unindexer.
	if txLookupLimit != nil {
		bc.txLookupLimit = *txLookupLimit
		bc.wg.Add(1)
		go bc.maintainTxIndex(txIndexBlock)
	}

	// If periodic cache journal is required, spin it up.
	if bc.cacheConfig.TrieCleanRejournal > 0 {
		if bc.cacheConfig.TrieCleanRejournal < time.Minute {
			log.Warn("Sanitizing invalid trie cache journal time", "provided", bc.cacheConfig.TrieCleanRejournal, "updated", time.Minute)
			bc.cacheConfig.TrieCleanRejournal = time.Minute
		}
		triedb := bc.stateCache.TrieDB()
		bc.wg.Add(1)
		go func() {
			defer bc.wg.Done()
			triedb.SaveCachePeriodically(bc.cacheConfig.TrieCleanJournal, bc.cacheConfig.TrieCleanRejournal, bc.quit)
		}()
	}
	return bc, nil
}

// StopInsert interrupts all insertion methods, causing them to return
// errInsertionInterrupted as soon as possible. Insertion is permanently disabled after
// calling this method.
func (bc *BlockChain) StopInsert() {
	atomic.StoreInt32(&bc.procInterrupt, 1)
}

// insertStopped returns true after StopInsert has been called.
func (bc *BlockChain) insertStopped() bool {
	return atomic.LoadInt32(&bc.procInterrupt) == 1
}

// empty returns an indicator whether the blockchain is empty.
// Note, it's a special case that we connect a non-empty ancient
// database with an empty node, so that we can plugin the ancient
// into node seamlessly.
func (bc *BlockChain) empty() bool {
	genesis := bc.genesisBlock.Hash()
	for _, hash := range []common.Hash{rawdb.ReadHeadBlockHash(bc.db), rawdb.ReadHeadHeaderHash(bc.db), rawdb.ReadHeadFastBlockHash(bc.db)} {
		if hash != genesis {
			return false
		}
	}
	return true
}

// loadLastState loads the last known chain state from the database. This method
// assumes that the chain manager mutex is held.
func (bc *BlockChain) loadLastState() error {
	// Restore the last known head block
	head := rawdb.ReadHeadBlockHash(bc.db)
	if head == (common.Hash{}) {
		// Corrupt or empty database, init from scratch
		log.Warn("Empty database, resetting chain")
		return errors.New("Empty database, resetting chain")
		// return bc.Reset()
	}
	// Make sure the entire head block is available
	currentBlock := bc.GetBlockByHash(head)
	if currentBlock == nil {
		// Corrupt or empty database, init from scratch
		log.Warn("Head block missing, resetting chain", "hash", head)
		return errors.New("Head block missing, resetting chain")
		//return bc.Reset()
	}

	// Everything seems to be fine, set as the head block
	bc.currentBlock.Store(currentBlock)
	headBlockGauge.Update(int64(currentBlock.NumberU64()))

	// Restore the last known head header
	currentHeader := currentBlock.Header()
	if head := rawdb.ReadHeadHeaderHash(bc.db); head != (common.Hash{}) {
		if header := bc.GetHeaderByHash(head); header != nil {
			currentHeader = header
		}
	}
	bc.hc.SetCurrentHeader(currentHeader)

	// Restore the last known head fast block
	bc.currentFastBlock.Store(currentBlock)
	headFastBlockGauge.Update(int64(currentBlock.NumberU64()))

	if head := rawdb.ReadHeadFastBlockHash(bc.db); head != (common.Hash{}) {
		if block := bc.GetBlockByHash(head); block != nil {
			bc.currentFastBlock.Store(block)
			headFastBlockGauge.Update(int64(block.NumberU64()))
		}
	}
	// Issue a status log for the user
	currentFastBlock := bc.CurrentFastBlock()

	log.Info("Loaded most recent local header", "number", currentHeader.Number, "hash", currentHeader.Hash(), "age", common.PrettyAge(time.Unix(int64(currentHeader.Time), 0)))
	log.Info("Loaded most recent local full block", "number", currentBlock.Number(), "hash", currentBlock.Hash(), "age", common.PrettyAge(time.Unix(int64(currentBlock.Time()), 0)))
	log.Info("Loaded most recent local fast block", "number", currentFastBlock.Number(), "hash", currentFastBlock.Hash(), "age", common.PrettyAge(time.Unix(int64(currentFastBlock.Time()), 0)))

	return nil
}

// SetHead rewinds the local chain to a new head. Depending on whether the node
// was fast synced or full synced and in which state, the method will try to
// delete minimal data from disk whilst retaining chain consistency.
func (bc *BlockChain) SetHead(head uint64) error {
	// PlatON do not support rewind
	return nil
}

// setHeadBeyondRoot rewinds the local chain to a new head with the extra condition
// that the rewind must pass the specified state root. This method is meant to be
// used when rewiding with snapshots enabled to ensure that we go back further than
// persistent disk layer. Depending on whether the node was fast synced or full, and
// in which state, the method will try to delete minimal data from disk whilst
// retaining chain consistency.
//
// The method returns the block number where the requested root cap was found.
func (bc *BlockChain) setHeadBeyondRoot(head uint64, root common.Hash, repair bool) (uint64, error) {
	if !bc.chainmu.TryLock() {
		return 0, errChainStopped
	}
	defer bc.chainmu.Unlock()

	// Track the block number of the requested root hash
	var rootNumber uint64 // (no root == always 0)

	// Retrieve the last pivot block to short circuit rollbacks beyond it and the
	// current freezer limit to start nuking id underflown
	pivot := rawdb.ReadLastPivotNumber(bc.db)
	frozen, _ := bc.db.Ancients()

	updateFn := func(db ethdb.KeyValueWriter, header *types.Header) (uint64, bool) {
		// Rewind the blockchain, ensuring we don't end up with a stateless head
		// block. Note, depth equality is permitted to allow using SetHead as a
		// chain reparation mechanism without deleting any data!
		if currentBlock := bc.CurrentBlock(); currentBlock != nil && header.Number.Uint64() <= currentBlock.NumberU64() {
			newHeadBlock := bc.GetBlock(header.Hash(), header.Number.Uint64())
			if newHeadBlock == nil {
				log.Error("Gap in the chain, rewinding to genesis", "number", header.Number, "hash", header.Hash())
				newHeadBlock = bc.genesisBlock
			} else {
				// Block exists, keep rewinding until we find one with state,
				// keeping rewinding until we exceed the optional threshold
				// root hash
				beyondRoot := (root == common.Hash{}) // Flag whether we're beyond the requested root (no root, always true)

				for {
					// If a root threshold was requested but not yet crossed, check
					if root != (common.Hash{}) && !beyondRoot && newHeadBlock.Root() == root {
						beyondRoot, rootNumber = true, newHeadBlock.NumberU64()
					}
					if _, err := state.New(newHeadBlock.Root(), bc.stateCache, bc.snaps); err != nil {
						log.Trace("Block state missing, rewinding further", "number", newHeadBlock.NumberU64(), "hash", newHeadBlock.Hash())
						if pivot == nil || newHeadBlock.NumberU64() > *pivot {
							parent := bc.GetBlock(newHeadBlock.ParentHash(), newHeadBlock.NumberU64()-1)
							if parent != nil {
								newHeadBlock = parent
								continue
							}
							log.Error("Missing block in the middle, aiming genesis", "number", newHeadBlock.NumberU64()-1, "hash", newHeadBlock.ParentHash())
							newHeadBlock = bc.genesisBlock
						} else {
							log.Trace("Rewind passed pivot, aiming genesis", "number", newHeadBlock.NumberU64(), "hash", newHeadBlock.Hash(), "pivot", *pivot)
							newHeadBlock = bc.genesisBlock
						}
					}
					if beyondRoot || newHeadBlock.NumberU64() == 0 {
						log.Debug("Rewound to block with state", "number", newHeadBlock.NumberU64(), "hash", newHeadBlock.Hash())
						break
					}
					log.Debug("Skipping block with threshold state", "number", newHeadBlock.NumberU64(), "hash", newHeadBlock.Hash(), "root", newHeadBlock.Root())
					newHeadBlock = bc.GetBlock(newHeadBlock.ParentHash(), newHeadBlock.NumberU64()-1) // Keep rewinding
				}
			}
			rawdb.WriteHeadBlockHash(db, newHeadBlock.Hash())

			// Degrade the chain markers if they are explicitly reverted.
			// In theory we should update all in-memory markers in the
			// last step, however the direction of SetHead is from high
			// to low, so it's safe the update in-memory markers directly.
			bc.currentBlock.Store(newHeadBlock)
			headBlockGauge.Update(int64(newHeadBlock.NumberU64()))
		}
		// Rewind the fast block in a simpleton way to the target head
		if currentFastBlock := bc.CurrentFastBlock(); currentFastBlock != nil && header.Number.Uint64() < currentFastBlock.NumberU64() {
			newHeadFastBlock := bc.GetBlock(header.Hash(), header.Number.Uint64())
			// If either blocks reached nil, reset to the genesis state
			if newHeadFastBlock == nil {
				newHeadFastBlock = bc.genesisBlock
			}
			rawdb.WriteHeadFastBlockHash(db, newHeadFastBlock.Hash())

			// Degrade the chain markers if they are explicitly reverted.
			// In theory we should update all in-memory markers in the
			// last step, however the direction of SetHead is from high
			// to low, so it's safe the update in-memory markers directly.
			bc.currentFastBlock.Store(newHeadFastBlock)
			headFastBlockGauge.Update(int64(newHeadFastBlock.NumberU64()))
		}
		head := bc.CurrentBlock().NumberU64()

		// If setHead underflown the freezer threshold and the block processing
		// intent afterwards is full block importing, delete the chain segment
		// between the stateful-block and the sethead target.
		var wipe bool
		if head+1 < frozen {
			wipe = pivot == nil || head >= *pivot
		}
		return head, wipe // Only force wipe if full synced
	}
	// Rewind the header chain, deleting all block bodies until then
	delFn := func(db ethdb.KeyValueWriter, hash common.Hash, num uint64) {
		// Ignore the error here since light client won't hit this path
		frozen, _ := bc.db.Ancients()
		if num+1 <= frozen {
			// Truncate all relative data(header, total difficulty, body, receipt
			// and canonical hash) from ancient store.
			if err := bc.db.TruncateAncients(num); err != nil {
				log.Crit("Failed to truncate ancient data", "number", num, "err", err)
			}
			// Remove the hash <-> number mapping from the active store.
			rawdb.DeleteHeaderNumber(db, hash)
		} else {
			// Remove relative body and receipts from the active store.
			// The header, total difficulty and canonical hash will be
			// removed in the hc.SetHead function.
			rawdb.DeleteBody(db, hash, num)
			rawdb.DeleteReceipts(db, hash, num)
		}
		// Todo(rjl493456442) txlookup, bloombits, etc
	}
	// If SetHead was only called as a chain reparation method, try to skip
	// touching the header chain altogether, unless the freezer is broken
	hdnum := bc.CurrentBlock().Header().Number.Uint64()
	log.Debug("hdnum", "head", hdnum)
	if block := bc.CurrentBlock(); block.NumberU64() == head {
		if repair {
			if target, force := updateFn(bc.db, bc.CurrentBlock().Header()); force {
				bc.hc.SetHead(target, updateFn, delFn)
			}
		}
	} else {
		// Rewind the chain to the requested head and keep going backwards until a
		// block with a state is found or fast sync pivot is passed
		log.Warn("Rewinding blockchain", "target", head)
		bc.hc.SetHead(head, updateFn, delFn)
	}
	// Clear out any stale content from the caches
	bc.bodyCache.Purge()
	bc.bodyRLPCache.Purge()
	bc.receiptsCache.Purge()
	bc.blockCache.Purge()
	bc.txLookupCache.Purge()
	bc.futureBlocks.Purge()

	return rootNumber, bc.loadLastState()
}

// FastSyncCommitHead sets the current head block to the one defined by the hash
// irrelevant what the chain contents were prior.
func (bc *BlockChain) FastSyncCommitHead(hash common.Hash) error {

	// Make sure that both the block as well at its state trie exists
	block := bc.GetBlockByHash(hash)
	if block == nil {
		return fmt.Errorf("non existent block [%x..]", hash[:4])
	}
	if _, err := trie.NewSecure(block.Root(), bc.stateCache.TrieDB()); err != nil {
		return err
	}
	// If all checks out, manually set the head block

	// If all checks out, manually set the head block.
	if !bc.chainmu.TryLock() {
		return errChainStopped
	}
	bc.currentBlock.Store(block)
	headBlockGauge.Update(int64(block.NumberU64()))
	bc.chainmu.Unlock()

	// Destroy any existing state snapshot and regenerate it in the background,
	// also resuming the normal maintenance of any previously paused snapshot.
	if bc.snaps != nil {
		bc.snaps.Rebuild(block.Root())
	}
	log.Info("Committed new head block", "number", block.Number(), "hash", hash)
	bc.engine.Pause()
	defer bc.engine.Resume()
	return bc.engine.FastSyncCommitHead(block)
}

// SetProcessor sets the processor required for making state modifications.
func (bc *BlockChain) SetProcessor(processor Processor) {
	bc.procmu.Lock()
	defer bc.procmu.Unlock()
	bc.processor = processor
}

// SetValidator sets the validator which is used to validate incoming blocks.
func (bc *BlockChain) SetValidator(validator Validator) {
	bc.procmu.Lock()
	defer bc.procmu.Unlock()
	bc.validator = validator
}

// Reset purges the entire blockchain, restoring it to its genesis state.
//func (bc *BlockChain) Reset() error {
//	return bc.ResetWithGenesisBlock(bc.genesisBlock)
//}

// ResetWithGenesisBlock purges the entire blockchain, restoring it to the
// specified genesis state.
/*func (bc *BlockChain) ResetWithGenesisBlock(genesis *types.Block) error {
	// Dump the entire block chain and purge the caches
	if err := bc.SetHead(0); err != nil {
		return err
	}
	bc.chainmu.Lock()
	defer bc.chainmu.Unlock()

	// Prepare the genesis block and reinitialise the chain
	rawdb.WriteBlock(bc.db, genesis)

	bc.genesisBlock = genesis
	bc.insert(bc.genesisBlock)
	bc.currentBlock.Store(bc.genesisBlock)
	headBlockGauge.Update(int64(bc.genesisBlock.NumberU64()))
	bc.hc.SetGenesis(bc.genesisBlock.Header())
	bc.hc.SetCurrentHeader(bc.genesisBlock.Header())
	bc.currentFastBlock.Store(bc.genesisBlock)

	return nil
}*/

// repair tries to repair the current blockchain by rolling back the current block
// until one with associated state is found. This is needed to fix incomplete db
// writes caused either by crashes/power outages, or simply non-committed tries.
//
// This method only rolls back the current block. The current header and current
// fast block are left intact.
func (bc *BlockChain) repair(head **types.Block) error {
	for {
		// Abort if we've rewound to a head block that does have associated state
		if _, err := state.New((*head).Root(), bc.stateCache, bc.snaps); err == nil {
			log.Info("Rewound blockchain to past state", "number", (*head).Number(), "hash", (*head).Hash())
			return nil
		}
		// Otherwise rewind one block and recheck state availability there
		block := bc.GetBlock((*head).ParentHash(), (*head).NumberU64()-1)
		if block == nil {
			return fmt.Errorf("missing block %d [%x]", (*head).NumberU64()-1, (*head).ParentHash())
		}
		*head = block
	}
}

// Export writes the active chain to the given writer.
func (bc *BlockChain) Export(w io.Writer) error {
	return bc.ExportN(w, uint64(0), bc.CurrentBlock().NumberU64())
}

// ExportN writes a subset of the active chain to the given writer.
func (bc *BlockChain) ExportN(w io.Writer, first uint64, last uint64) error {
	if !bc.chainmu.TryLock() {
		return errChainStopped
	}
	defer bc.chainmu.Unlock()

	if first > last {
		return fmt.Errorf("export failed: first (%d) is greater than last (%d)", first, last)
	}
	log.Info("Exporting batch of blocks", "count", last-first+1)

	start, reported := time.Now(), time.Now()
	for nr := first; nr <= last; nr++ {
		block := bc.GetBlockByNumber(nr)
		if block == nil {
			return fmt.Errorf("export failed on #%d: not found", nr)
		}
		if err := block.EncodeRLP(w); err != nil {
			return err
		}
		if time.Since(reported) >= statsReportLimit {
			log.Info("Exporting blocks", "exported", block.NumberU64()-first, "elapsed", common.PrettyDuration(time.Since(start)))
			reported = time.Now()
		}
	}

	return nil
}

// writeHeadBlock injects a new head block into the current block chain. This method
// assumes that the block is indeed a true head. It will also reset the head
// header and the head fast sync block to this very same block if they are older
// or if they are on a different side chain.
//
// Note, this function assumes that the `mu` mutex is held!
func (bc *BlockChain) writeHeadBlock(block *types.Block) {
	// If the block is on a side chain or an unknown one, force other heads onto it too
	updateHeads := rawdb.ReadCanonicalHash(bc.db, block.NumberU64()) != block.Hash()

	// Add the block to the canonical chain number scheme and mark as the head
	batch := bc.db.NewBatch()
	rawdb.WriteCanonicalHash(batch, block.Hash(), block.NumberU64())
	rawdb.WriteTxLookupEntriesByBlock(batch, block)
	rawdb.WriteHeadBlockHash(batch, block.Hash())

	// If the block is better than our head or is on a different chain, force update heads
	if updateHeads {
		rawdb.WriteHeadHeaderHash(batch, block.Hash())
		rawdb.WriteHeadFastBlockHash(batch, block.Hash())
	}
	// Flush the whole batch into the disk, exit the node if failed
	if err := batch.Write(); err != nil {
		log.Crit("Failed to update chain indexes and markers", "err", err)
	}
	// Update all in-memory chain markers in the last step
	if updateHeads {
		bc.hc.SetCurrentHeader(block.Header())
		bc.currentFastBlock.Store(block)
		headFastBlockGauge.Update(int64(block.NumberU64()))
	}
	bc.currentBlock.Store(block)
	headBlockGauge.Update(int64(block.NumberU64()))
}

func (bc *BlockChain) getBlock(hash common.Hash, number uint64) *types.Block {
	// Short circuit if the block's already in the cache, retrieve otherwise
	if block, ok := bc.blockCache.Get(hash); ok {
		return block.(*types.Block)
	}
	block := rawdb.ReadBlock(bc.db, hash, number)
	if block == nil {
		return nil
	}
	// Cache the found block for next time and return
	bc.blockCache.Add(block.Hash(), block)
	return block
}

// Stop stops the blockchain service. If any imports are currently in progress
// it will abort them using the procInterrupt.
func (bc *BlockChain) Stop() {
	if !atomic.CompareAndSwapInt32(&bc.running, 0, 1) {
		return
	}

	bc.executeWG.Wait()
	// Unsubscribe all subscriptions registered from blockchain
	bc.scope.Close()
	close(bc.quit)
	bc.StopInsert()

	bc.cleaner.Stop()

	// Now wait for all chain modifications to end and persistent goroutines to exit.
	//
	// Note: Close waits for the mutex to become available, i.e. any running chain
	// modification will have exited when Close returns. Since we also called StopInsert,
	// the mutex should become available quickly. It cannot be taken again after Close has
	// returned.
	bc.chainmu.Close()
	bc.wg.Wait()

	// Ensure that the entirety of the state snapshot is journalled to disk.
	var snapBase common.Hash
	if bc.snaps != nil {
		var err error
		if snapBase, err = bc.snaps.Journal(bc.CurrentBlock().Root()); err != nil {
			log.Error("Failed to journal state snapshot", "err", err)
		}
	}
	// Ensure the state of a recent block is also stored to disk before exiting.
	// We're writing three different states to catch different restart scenarios:
	//  - HEAD:     So we don't need to reprocess any blocks in the general case
	//  - HEAD-1:   So we don't do large reorgs if our HEAD becomes an uncle
	//  - HEAD-127: So we have a hard limit on the number of blocks reexecuted
	if !bc.cacheConfig.Disabled {
		triedb := bc.stateCache.TrieDB()
		if 0 >= bc.cacheConfig.TriesInMemory {
			bc.cacheConfig.TriesInMemory = 128
		}

		for _, offset := range []uint64{0, 1, (uint64)(bc.cacheConfig.TriesInMemory - 1)} {
			if number := bc.CurrentBlock().NumberU64(); number > offset {
				recent := bc.GetBlockByNumber(number - offset)

				log.Info("Writing cached state to disk", "block", recent.Number(), "hash", recent.Hash(), "root", recent.Root())
				if err := triedb.Commit(recent.Root(), true, true); err != nil {
					log.Error("Failed to commit recent state trie", "err", err)
				}
			}
		}
		if snapBase != (common.Hash{}) {
			log.Info("Writing snapshot state to disk", "root", snapBase)
			if err := triedb.Commit(snapBase, true, true); err != nil {
				log.Error("Failed to commit recent state trie", "err", err)
			}
		}
		for !bc.triegc.Empty() {
			triedb.Dereference(bc.triegc.PopItem().(common.Hash))
		}
		if size, _ := triedb.Size(); size != 0 {
			log.Error("Dangling trie nodes after full cleanup")
		}
	}
	// Ensure all live cached entries be saved into disk, so that we can skip
	// cache warmup when node restarts.
	if bc.cacheConfig.TrieCleanJournal != "" {
		triedb := bc.stateCache.TrieDB()
		triedb.SaveCache(bc.cacheConfig.TrieCleanJournal)
	}
	log.Info("Blockchain stopped")
}

func (bc *BlockChain) procFutureBlocks() {
	blocks := make([]*types.Block, 0, bc.futureBlocks.Len())
	for _, hash := range bc.futureBlocks.Keys() {
		if block, exist := bc.futureBlocks.Peek(hash); exist {
			blocks = append(blocks, block.(*types.Block))
		}
	}
	if len(blocks) > 0 {
		sort.Slice(blocks, func(i, j int) bool {
			return blocks[i].NumberU64() < blocks[j].NumberU64()
		})

		// Insert one by one as chain insertion needs contiguous ancestry between blocks
		for i := range blocks {
			bc.InsertChain(blocks[i : i+1])
		}
	}
}

// WriteStatus status of write
type WriteStatus byte

const (
	NonStatTy WriteStatus = iota
	CanonStatTy
	SideStatTy
)

// numberHash is just a container for a number and a hash, to represent a block
type numberHash struct {
	number uint64
	hash   common.Hash
}

// InsertReceiptChain attempts to complete an already existing header chain with
// transaction and receipt data.
func (bc *BlockChain) InsertReceiptChain(blockChain types.Blocks, receiptChain []types.Receipts, ancientLimit uint64) (int, error) {
	// We don't require the chainMu here since we want to maximize the
	// concurrency of header insertion and receipt insertion.
	bc.wg.Add(1)
	defer bc.wg.Done()

	var (
		ancientBlocks, liveBlocks     types.Blocks
		ancientReceipts, liveReceipts []types.Receipts
	)
	// Do a sanity check that the provided chain is actually ordered and linked
	for i := 0; i < len(blockChain); i++ {
		if i != 0 {
			if blockChain[i].NumberU64() != blockChain[i-1].NumberU64()+1 || blockChain[i].ParentHash() != blockChain[i-1].Hash() {
				log.Error("Non contiguous receipt insert", "number", blockChain[i].Number(), "hash", blockChain[i].Hash(), "parent", blockChain[i].ParentHash(),
					"prevnumber", blockChain[i-1].Number(), "prevhash", blockChain[i-1].Hash())
				return 0, fmt.Errorf("non contiguous insert: item %d is #%d [%x..], item %d is #%d [%x..] (parent [%x..])", i-1, blockChain[i-1].NumberU64(),
					blockChain[i-1].Hash().Bytes()[:4], i, blockChain[i].NumberU64(), blockChain[i].Hash().Bytes()[:4], blockChain[i].ParentHash().Bytes()[:4])
			}
		}
		if blockChain[i].NumberU64() <= ancientLimit {
			ancientBlocks, ancientReceipts = append(ancientBlocks, blockChain[i]), append(ancientReceipts, receiptChain[i])
		} else {
			liveBlocks, liveReceipts = append(liveBlocks, blockChain[i]), append(liveReceipts, receiptChain[i])
		}
	}

	var (
		stats = struct{ processed, ignored int32 }{}
		start = time.Now()
		size  = int64(0)
	)

	// updateHead updates the head fast sync block if the inserted blocks are better
	// and returns an indicator whether the inserted blocks are canonical.
	updateHead := func(head *types.Block) bool {
		if !bc.chainmu.TryLock() {
			return false
		}
		defer bc.chainmu.Unlock()

		if bn := head.Number(); bn != nil {
			// Rewind may have occurred, skip in that case.
			if bc.CurrentHeader().Number.Cmp(head.Number()) >= 0 {
				currentFastBlock := bc.CurrentFastBlock()
				if currentFastBlock.Number().Cmp(bn) < 0 {
					rawdb.WriteHeadFastBlockHash(bc.db, head.Hash())
					bc.currentFastBlock.Store(head)
					headFastBlockGauge.Update(int64(head.NumberU64()))
					return true
				}
			}
		}
		return false
	}

	// writeAncient writes blockchain and corresponding receipt chain into ancient store.
	//
	// this function only accepts canonical chain data. All side chain will be reverted
	// eventually.
	writeAncient := func(blockChain types.Blocks, receiptChain []types.Receipts) (int, error) {
		first := blockChain[0]
		last := blockChain[len(blockChain)-1]

		// Ensure genesis is in ancients.
		if first.NumberU64() == 1 {
			if frozen, _ := bc.db.Ancients(); frozen == 0 {
				b := bc.genesisBlock
				writeSize, err := rawdb.WriteAncientBlocks(bc.db, []*types.Block{b}, []types.Receipts{nil})
				size += writeSize
				if err != nil {
					log.Error("Error writing genesis to ancients", "err", err)
					return 0, err
				}
				log.Info("Wrote genesis to ancients")
			}
		}
		// Before writing the blocks to the ancients, we need to ensure that
		// they correspond to the what the headerchain 'expects'.
		// We only check the last block/header, since it's a contiguous chain.
		if !bc.HasHeader(last.Hash(), last.NumberU64()) {
			return 0, fmt.Errorf("containing header #%d [%x..] unknown", last.Number(), last.Hash().Bytes()[:4])
		}

		// Write all chain data to ancients.
		writeSize, err := rawdb.WriteAncientBlocks(bc.db, blockChain, receiptChain)
		size += writeSize
		if err != nil {
			log.Error("Error importing chain data to ancients", "err", err)
			return 0, err
		}

		// Write tx indices if any condition is satisfied:
		// * If user requires to reserve all tx indices(txlookuplimit=0)
		// * If all ancient tx indices are required to be reserved(txlookuplimit is even higher than ancientlimit)
		// * If block number is large enough to be regarded as a recent block
		// It means blocks below the ancientLimit-txlookupLimit won't be indexed.
		//
		// But if the `TxIndexTail` is not nil, e.g. Geth is initialized with
		// an external ancient database, during the setup, blockchain will start
		// a background routine to re-indexed all indices in [ancients - txlookupLimit, ancients)
		// range. In this case, all tx indices of newly imported blocks should be
		// generated.
		var batch = bc.db.NewBatch()
		for _, block := range blockChain {
			if bc.txLookupLimit == 0 || ancientLimit <= bc.txLookupLimit || block.NumberU64() >= ancientLimit-bc.txLookupLimit {
				rawdb.WriteTxLookupEntriesByBlock(batch, block)
			} else if rawdb.ReadTxIndexTail(bc.db) != nil {
				rawdb.WriteTxLookupEntriesByBlock(batch, block)
			}
			stats.processed++
		}

		// Flush all tx-lookup index data.
		size += int64(batch.ValueSize())
		if err := batch.Write(); err != nil {
			// The tx index data could not be written.
			// Roll back the ancient store update.
			fastBlock := bc.CurrentFastBlock().NumberU64()
			if err := bc.db.TruncateAncients(fastBlock + 1); err != nil {
				log.Error("Can't truncate ancient store after failed insert", "err", err)
			}
			return 0, err
		}

		// Sync the ancient store explicitly to ensure all data has been flushed to disk.
		if err := bc.db.Sync(); err != nil {
			return 0, err
		}

		// Update the current fast block because all block data is now present in DB.
		previousFastBlock := bc.CurrentFastBlock().NumberU64()
		if !updateHead(blockChain[len(blockChain)-1]) {
			// We end up here if the header chain has reorg'ed, and the blocks/receipts
			// don't match the canonical chain.
			if err := bc.db.TruncateAncients(previousFastBlock + 1); err != nil {
				log.Error("Can't truncate ancient store after failed insert", "err", err)
			}
			return 0, errSideChainReceipts
		}

		// Delete block data from the main database.
		batch.Reset()
		canonHashes := make(map[common.Hash]struct{})
		for _, block := range blockChain {
			canonHashes[block.Hash()] = struct{}{}
			if block.NumberU64() == 0 {
				continue
			}
			rawdb.DeleteCanonicalHash(batch, block.NumberU64())
			rawdb.DeleteBlockWithoutNumber(batch, block.Hash(), block.NumberU64())
		}
		// Delete side chain hash-to-number mappings.
		for _, nh := range rawdb.ReadAllHashesInRange(bc.db, first.NumberU64(), last.NumberU64()) {
			if _, canon := canonHashes[nh.Hash]; !canon {
				rawdb.DeleteHeader(batch, nh.Hash, nh.Number)
			}
		}
		if err := batch.Write(); err != nil {
			return 0, err
		}
		return 0, nil
	}

	// writeLive writes blockchain and corresponding receipt chain into active store.
	writeLive := func(blockChain types.Blocks, receiptChain []types.Receipts) (int, error) {
		skipPresenceCheck := false
		batch := bc.db.NewBatch()
		for i, block := range blockChain {
			// Short circuit insertion if shutting down or processing failed
			if bc.insertStopped() {
				return 0, errInsertionInterrupted
			}
			// Short circuit if the owner header is unknown
			if !bc.HasHeader(block.Hash(), block.NumberU64()) {
				return i, fmt.Errorf("containing header #%d [%x..] unknown", block.Number(), block.Hash().Bytes()[:4])
			}
			if !skipPresenceCheck {
				// Ignore if the entire data is already known
				if bc.HasBlock(block.Hash(), block.NumberU64()) {
					stats.ignored++
					continue
				} else {
					// If block N is not present, neither are the later blocks.
					// This should be true, but if we are mistaken, the shortcut
					// here will only cause overwriting of some existing data
					skipPresenceCheck = true
				}
			}
			// Write all the data out into the database
			rawdb.WriteBody(batch, block.Hash(), block.NumberU64(), block.Body())
			if receiptChain != nil {
				rawdb.WriteReceipts(batch, block.Hash(), block.NumberU64(), receiptChain[i])
			}
			rawdb.WriteTxLookupEntriesByBlock(batch, block) // Always write tx indices for live blocks, we assume they are needed

			// Write everything belongs to the blocks into the database. So that
			// we can ensure all components of body is completed(body, receipts,
			// tx indexes)
			if batch.ValueSize() >= ethdb.IdealBatchSize {
				if err := batch.Write(); err != nil {
					return 0, err
				}
				size += int64(batch.ValueSize())
				batch.Reset()
			}
			stats.processed++
		}
		// Write everything belongs to the blocks into the database. So that
		// we can ensure all components of body is completed(body, receipts,
		// tx indexes)
		if batch.ValueSize() > 0 {
			size += int64(batch.ValueSize())
			if err := batch.Write(); err != nil {
				return 0, err
			}
		}
		updateHead(blockChain[len(blockChain)-1])
		return 0, nil
	}

	// Write downloaded chain data and corresponding receipt chain data
	if len(ancientBlocks) > 0 {
		// fast同步的时候不会写入回执
		//if n, err := writeAncient(ancientBlocks, ancientReceipts); err != nil {
		if n, err := writeAncient(ancientBlocks, nil); err != nil {
			if err == errInsertionInterrupted {
				return 0, nil
			}
			return n, err
		}
	}
	// Write the tx index tail (block number from where we index) before write any live blocks
	if len(liveBlocks) > 0 && liveBlocks[0].NumberU64() == ancientLimit+1 {
		// The tx index tail can only be one of the following two options:
		// * 0: all ancient blocks have been indexed
		// * ancient-limit: the indices of blocks before ancient-limit are ignored
		if tail := rawdb.ReadTxIndexTail(bc.db); tail == nil {
			if bc.txLookupLimit == 0 || ancientLimit <= bc.txLookupLimit {
				rawdb.WriteTxIndexTail(bc.db, 0)
			} else {
				rawdb.WriteTxIndexTail(bc.db, ancientLimit-bc.txLookupLimit)
			}
		}
	}
	if len(liveBlocks) > 0 {
		// fast同步的时候不会写入回执
		// if n, err := writeLive(liveBlocks, liveReceipts); err != nil {
		if n, err := writeLive(liveBlocks, nil); err != nil {
			if err == errInsertionInterrupted {
				return 0, nil
			}
			return n, err
		}
	}

	head := blockChain[len(blockChain)-1]
	context := []interface{}{
		"count", stats.processed, "elapsed", common.PrettyDuration(time.Since(start)),
		"number", head.Number(), "hash", head.Hash(), "age", common.PrettyAge(time.Unix(int64(head.Time()), 0)),
		"size", common.StorageSize(size),
	}
	if stats.ignored > 0 {
		context = append(context, []interface{}{"ignored", stats.ignored}...)
	}
	log.Info("Imported new block receipts", context...)

	return 0, nil
}

var lastWrite uint64

// WriteBlockWithoutState writes only the block and its metadata to the database,
// but does not write any state. This is used to construct competing side forks
// up to the point where they exceed the canonical total difficulty.
func (bc *BlockChain) WriteBlockWithoutState(block *types.Block) (err error) {
	if bc.insertStopped() {
		return errInsertionInterrupted
	}

	rawdb.WriteBlock(bc.db, block)

	return nil
}

// WriteBlockWithState writes the block and all associated state to the database.
func (bc *BlockChain) WriteBlockWithState(block *types.Block, receipts []*types.Receipt, logs []*types.Log, state *state.StateDB, emitHeadEvent bool) (status WriteStatus, err error) {
	if !bc.chainmu.TryLock() {
		return NonStatTy, errInsertionInterrupted
	}
	defer bc.chainmu.Unlock()
	return bc.writeBlockWithState(block, receipts, logs, state, emitHeadEvent)
}

// WriteBlockWithState writes the block and all associated state to the database.
func (bc *BlockChain) writeBlockWithState(block *types.Block, receipts []*types.Receipt, logs []*types.Log, state *state.StateDB, emitHeadEvent bool) (status WriteStatus, err error) {
	if bc.insertStopped() {
		return NonStatTy, errInsertionInterrupted
	}

	// Make sure no inconsistent state is leaked during insertion
	currentBlock := bc.CurrentBlock()
	if block.NumberU64() <= currentBlock.NumberU64() {
		log.Warn("block lower than current block in chain", "blockHash", block.Hash(), "blockNumber", block.NumberU64(), "currentHash", currentBlock.Hash(), "currentNumber", currentBlock.NumberU64())
		return NonStatTy, nil
	}
	localBn := currentBlock.Number()
	externBn := block.Number()

	// Irrelevant of the canonical status, write the block itself to the database
	rawdb.WriteBlock(bc.db, block)

	triedb := bc.stateCache.TrieDB()
	root, err := state.Commit(true)

	if err != nil {
		log.Error("check block is EIP158 error", "hash", block.Hash(), "number", block.NumberU64())
		return NonStatTy, err
	}

	// If we're running an archive node, always flush
	if bc.cacheConfig.Disabled {
		limit := common.StorageSize(bc.cacheConfig.TrieDirtyLimit) * 1024 * 1024
		oversize := false
		if !(bc.cacheConfig.DBGCMpt && !bc.cacheConfig.DBDisabledGC.IsSet()) {
			//triedb.ReferenceVersion(root)
			if err := triedb.Commit(root, false, false); err != nil {
				log.Error("Commit to triedb error", "root", root)
				return NonStatTy, err
			}
			//triedb.Dereference(currentBlock.Root())
			nodes, _ := triedb.Size()
			oversize = nodes > limit
		} else {
			triedb.ReferenceVersion(root)
			triedb.DereferenceDB(currentBlock.Root())

			if err := triedb.Commit(root, false, false); err != nil {
				log.Error("Commit to triedb error", "root", root)
				return NonStatTy, err
			}

			if triedb.UselessSize() > bc.cacheConfig.DBGCBlock {
				triedb.UselessGC(1)
			}

			nodes, _ := triedb.Size()
			oversize = nodes > limit
		}

		if oversize {
			triedb.CapNode(limit * defaultCapNodePercent)
			triedb.ResetUseless()
		}
		log.Debug("archive node commit stateDB trie", "blockNumber", block.NumberU64(), "blockHash", block.Hash().Hex(), "root", root.String())
	} else {
		log.Debug("non-archive node put stateDB trie", "blockNumber", block.NumberU64(), "blockHash", block.Hash().Hex(), "root", root.String())
		// Full but not archive node, do proper garbage collection
		triedb.Reference(root, common.Hash{}) // metadata reference to keep trie alive
		bc.triegc.Push(root, -int64(block.NumberU64()))

		if 0 >= bc.cacheConfig.TriesInMemory {
			bc.cacheConfig.TriesInMemory = 128
		}

		if current := block.NumberU64(); current > (uint64)(bc.cacheConfig.TriesInMemory) {
			// If we exceeded our memory allowance, flush matured singleton nodes to disk
			var (
				nodes, imgs = triedb.Size()
				limit       = common.StorageSize(bc.cacheConfig.TrieDirtyLimit) * 1024 * 1024
			)
			if nodes > limit || imgs > 4*1024*1024 {
				triedb.Cap(limit - ethdb.IdealBatchSize)
			}
			// Find the next state trie we need to commit
			header := bc.GetHeaderByNumber(current - (uint64)(bc.cacheConfig.TriesInMemory))
			chosen := header.Number.Uint64()

			// If we exceeded out time allowance, flush an entire trie to disk
			if bc.gcproc > bc.cacheConfig.TrieTimeLimit {
				// If we're exceeding limits but haven't reached a large enough memory gap,
				// warn the user that the system is becoming unstable.
				if chosen < lastWrite+(uint64)(bc.cacheConfig.TriesInMemory) && bc.gcproc >= 2*bc.cacheConfig.TrieTimeLimit {
					log.Info("State in memory for too long, committing", "time", bc.gcproc, "allowance", bc.cacheConfig.TrieTimeLimit, "optimum",
						float64(chosen-lastWrite)/(float64)(bc.cacheConfig.TriesInMemory))
				}
				// Flush an entire trie and restart the counters
				triedb.Commit(header.Root, true, true)
				lastWrite = chosen
				bc.gcproc = 0
			}
			// Garbage collect anything below our required write retention
			for !bc.triegc.Empty() {
				root, number := bc.triegc.Pop()
				if uint64(-number) > chosen {
					bc.triegc.Push(root, number)
					break
				}
				triedb.Dereference(root.(common.Hash))
			}
		}
	}

	// Write other block data using a batch.
	batch := bc.db.NewBatch()
	rawdb.WriteReceipts(batch, block.Hash(), block.NumberU64(), receipts)

	// If the total difficulty is higher than our known, add it to the canonical chain
	// Second clause in the if statement reduces the vulnerability to selfish mining.
	// Please refer to http://www.cs.cornell.edu/~ie53/publications/btcProcFC.pdf
	reorg := externBn.Cmp(localBn) > 0
	currentBlock = bc.CurrentBlock()
	if !reorg && externBn.Cmp(localBn) == 0 {
		// Split same-difficulty blocks by number, then preferentially select
		// the block generated by the local miner as the canonical block.
		if block.NumberU64() < currentBlock.NumberU64() {
			reorg = true
		} else if block.NumberU64() == currentBlock.NumberU64() {
			var currentPreserve, blockPreserve bool
			if bc.shouldPreserve != nil {
				currentPreserve, blockPreserve = bc.shouldPreserve(currentBlock), bc.shouldPreserve(block)
			}
			reorg = !currentPreserve && (blockPreserve || mrand.Float64() < 0.5)
		}
	}
	if reorg {
		// Reorganise the chain if the parent is not the head block
		if block.ParentHash() != currentBlock.Hash() {
			if err := bc.reorg(currentBlock, block); err != nil {
				return NonStatTy, err
			}
		}
		// Write the positional metadata for transaction/receipt lookups and preimages
		rawdb.WriteTxLookupEntriesByBlock(batch, block)
		rawdb.WritePreimages(batch, state.Preimages())

		status = CanonStatTy
	} else {
		status = SideStatTy
	}

	// Write the positional metadata for transaction/receipt lookups and preimages
	rawdb.WriteTxLookupEntriesByBlock(batch, block)
	rawdb.WritePreimages(batch, state.Preimages())

	status = CanonStatTy
	if err := batch.Write(); err != nil {
		return NonStatTy, err
	}
	log.Debug("insert into chain", "WriteStatus", status, "hash", block.Hash(), "number", block.NumberU64())

	// Set new head.
	if status == CanonStatTy {
		bc.writeHeadBlock(block)

		// parse block and retrieves txs
		//receipts := bc.GetReceiptsByHash(block.Hash())
		//if MPC_POOL != nil{
		//	MPC_POOL.InjectTxs(block, receipts, bc, state)
		//}

		//if VC_POOL != nil {
		//	VC_POOL.InjectTxs(block, receipts, bc, state)
		//}

	}
	bc.futureBlocks.Remove(block.Hash())

	if status == CanonStatTy {
		bc.chainFeed.Send(ChainEvent{Block: block, Hash: block.Hash(), Logs: logs})
		if len(logs) > 0 {
			bc.logsFeed.Send(logs)
		}
		// In theory we should fire a ChainHeadEvent when we inject
		// a canonical block, but sometimes we can insert a batch of
		// canonicial blocks. Avoid firing too much ChainHeadEvents,
		// we will fire an accumulated ChainHeadEvent and disable fire
		// event here.
		if emitHeadEvent {
			bc.chainHeadFeed.Send(ChainHeadEvent{Block: block})
		}
	} else {
		bc.chainSideFeed.Send(ChainSideEvent{Block: block})
	}

	bc.blockCache.Add(block.Hash(), block)

	bc.hc.SetCurrentHeader(block.Header())
	// Cleanup storage
	if !bc.cacheConfig.DBDisabledGC.IsSet() && bc.cleaner.NeedCleanup() {
		bc.cleaner.Cleanup()
	}

	bc.BlockFeed.Send(block)

	// Update the metrics touched during block processing
	accountReadTimer.Update(state.AccountReads)                 // Account reads are complete, we can mark them
	storageReadTimer.Update(state.StorageReads)                 // Storage reads are complete, we can mark them
	accountUpdateTimer.Update(state.AccountUpdates)             // Account updates are complete, we can mark them
	storageUpdateTimer.Update(state.StorageUpdates)             // Storage updates are complete, we can mark them
	snapshotAccountReadTimer.Update(state.SnapshotAccountReads) // Account reads are complete, we can mark them
	snapshotStorageReadTimer.Update(state.SnapshotStorageReads) // Storage reads are complete, we can mark them
	accountHashTimer.Update(state.AccountHashes)                // Account hashes are complete, we can mark them
	storageHashTimer.Update(state.StorageHashes)                // Storage hashes are complete, we can mark them
	accountCommitTimer.Update(state.AccountCommits)             // Account commits are complete, we can mark them
	storageCommitTimer.Update(state.StorageCommits)             // Storage commits are complete, we can mark them
	snapshotCommitTimer.Update(state.SnapshotCommits)           // Snapshot commits are complete, we can mark them
	return status, nil
}

// InsertChain attempts to insert the given batch of blocks in to the canonical
// chain or, otherwise, create a fork. If an error is returned it will return
// the index number of the failing block as well an error describing what went
// wrong.
//
// After insertion is done, all accumulated events will be fired.
func (bc *BlockChain) InsertChain(chain types.Blocks) (int, error) {
	// Sanity check that we have something meaningful to import
	if len(chain) == 0 {
		return 0, nil
	}

	// Do a sanity check that the provided chain is actually ordered and linked.
	for i := 1; i < len(chain); i++ {
		block, prev := chain[i], chain[i-1]
		if block.NumberU64() != prev.NumberU64()+1 || block.ParentHash() != prev.Hash() {
			log.Error("Non contiguous block insert",
				"number", block.Number(),
				"hash", block.Hash(),
				"parent", block.ParentHash(),
				"prevnumber", prev.Number(),
				"prevhash", prev.Hash(),
			)
			return 0, fmt.Errorf("non contiguous insert: item %d is #%d [%x..], item %d is #%d [%x..] (parent [%x..])", i-1, prev.NumberU64(),
				prev.Hash().Bytes()[:4], i, block.NumberU64(), block.Hash().Bytes()[:4], block.ParentHash().Bytes()[:4])
		}
	}
	// Pre-check passed, start the full block imports.
	return bc.insertChain(chain, true)
}

// insertChain is the internal implementation of insertChain, which assumes that
// 1) chains are contiguous, and 2) The chain mutex is held.
//
// This method is split out so that import batches that require re-injecting
// historical blocks can do so without releasing the lock, which could lead to
// racey behaviour. If a sidechain import is in progress, and the historic state
// is imported, but then new canon-head is added before the actual sidechain
// completes, then the historic state could be pruned again
func (bc *BlockChain) insertChain(chain types.Blocks, verifySeals bool) (int, error) {
	// If the chain is terminating, don't even bother starting up
	if bc.insertStopped() {
		return 0, nil
	}

	// A queued approach to delivering events. This is generally
	// faster than direct delivery and requires much less mutex
	// acquiring.
	var (
		stats = insertStats{startTime: mclock.Now()}
	)
	// Start the parallel header verifier
	headers := make([]*types.Header, len(chain))
	seals := make([]bool, len(chain))

	for i, block := range chain {
		headers[i] = block.Header()
		seals[i] = verifySeals
	}
	abort, results := bc.engine.VerifyHeaders(bc, headers, seals)
	defer close(abort)

	// Pause engine
	bc.engine.Pause()
	defer bc.engine.Resume()

	// Peek the error for the first block to decide the directing import logic
	it := newInsertIterator(chain, results, bc.Validator())
	block, err := it.next()
	switch {
	case bc.skipBlock(err, it):
		// Skip all known blocks that behind us
		current := bc.CurrentBlock().NumberU64()

		for block != nil && bc.skipBlock(err, it) && current >= block.NumberU64() {
			stats.ignored++
			block, err = it.next()
		}
		// Falls through to the block import

	// Some other error(except ErrKnownBlock) occurred, abort.
	// ErrKnownBlock is allowed here since some known blocks
	// still need re-execution to generate snapshots that are missing
	case err != nil && !errors.Is(err, ErrKnownBlock):
		bc.futureBlocks.Remove(block.Hash())
		stats.ignored += len(it.chain)
		bc.reportBlock(block, nil, err)
		return it.index, err
	}
	// No validation errors for the first block (or chain prefix skipped)
	for ; block != nil && err == nil; block, err = it.next() {
		// If the chain is terminating, stop processing blocks
		if bc.insertStopped() {
			log.Debug("Abort during block processing")
			break
		}
		start := time.Now()
		err = bc.engine.InsertChain(block)
		if err != nil {
			return it.index, err
		}

		blockInsertTimer.UpdateSince(start)
		dirty, _ := bc.stateCache.TrieDB().Size()
		stats.report(chain, it.index, dirty)
	}

	stats.ignored += it.remaining()

	return it.index, err
}

// joey.lyu
func (bc *BlockChain) ProcessDirectly(block *types.Block, state *state.StateDB, parent *types.Block) (types.Receipts, error) {
	// Process block using the parent state as reference point.
	start := time.Now()
	receipts, _, usedGas, err := bc.processor.Process(block, state, bc.vmConfig)
	if err != nil {
		log.Error("Failed to ProcessDirectly", "blockNumber", block.Number(), "blockHash", block.Hash(), "err", err)
		bc.reportBlock(block, receipts, err)
		return nil, err
	}
	log.Debug("Execute block time", "blockNumber", block.Number(), "blockHash", block.Hash(), "time", time.Since(start))

	// Validate the state using the default validator
	start = time.Now()
	err = bc.Validator().ValidateState(block, state, receipts, usedGas)
	if err != nil {
		bc.reportBlock(block, receipts, err)
		return nil, err
	}
	blockValidationTimer.UpdateSince(start)

	bc.BlockExecuteFeed.Send(block)

	return receipts, nil
}

func countTransactions(chain []*types.Block) (c int) {
	for _, b := range chain {
		c += len(b.Transactions())
	}
	return c
}

// reorg takes two blocks, an old chain and a new chain and will reconstruct the
// blocks and inserts them to be part of the new canonical chain and accumulates
// potential missing transactions and post an event about them.
func (bc *BlockChain) reorg(oldBlock, newBlock *types.Block) error {
	var (
		newChain    types.Blocks
		oldChain    types.Blocks
		commonBlock *types.Block

		deletedTxs types.Transactions
		addedTxs   types.Transactions

		deletedLogs [][]*types.Log
		rebirthLogs [][]*types.Log

		// collectLogs collects the logs that were generated or removed during
		// the processing of the block that corresponds with the given hash.
		// These logs are later announced as deleted or reborn
		collectLogs = func(hash common.Hash, removed bool) {
			number := bc.hc.GetBlockNumber(hash)
			if number == nil {
				return
			}
			receipts := rawdb.ReadReceipts(bc.db, hash, *number, bc.Config())

			var logs []*types.Log
			for _, receipt := range receipts {
				for _, log := range receipt.Logs {
					l := *log
					if removed {
						l.Removed = true
					} else {
					}
					logs = append(logs, &l)
				}
			}
			if len(logs) > 0 {
				if removed {
					deletedLogs = append(deletedLogs, logs)
				} else {
					rebirthLogs = append(rebirthLogs, logs)
				}
			}
		}
		// mergeLogs returns a merged log slice with specified sort order.
		mergeLogs = func(logs [][]*types.Log, reverse bool) []*types.Log {
			var ret []*types.Log
			if reverse {
				for i := len(logs) - 1; i >= 0; i-- {
					ret = append(ret, logs[i]...)
				}
			} else {
				for i := 0; i < len(logs); i++ {
					ret = append(ret, logs[i]...)
				}
			}
			return ret
		}
	)
	// Reduce the longer chain to the same number as the shorter one
	if oldBlock.NumberU64() > newBlock.NumberU64() {
		// Old chain is longer, gather all transactions and logs as deleted ones
		for ; oldBlock != nil && oldBlock.NumberU64() != newBlock.NumberU64(); oldBlock = bc.GetBlock(oldBlock.ParentHash(), oldBlock.NumberU64()-1) {
			oldChain = append(oldChain, oldBlock)
			deletedTxs = append(deletedTxs, oldBlock.Transactions()...)
			collectLogs(oldBlock.Hash(), true)
		}
	} else {
		// New chain is longer, stash all blocks away for subsequent insertion
		for ; newBlock != nil && newBlock.NumberU64() != oldBlock.NumberU64(); newBlock = bc.GetBlock(newBlock.ParentHash(), newBlock.NumberU64()-1) {
			newChain = append(newChain, newBlock)
		}
	}
	if oldBlock == nil {
		return fmt.Errorf("invalid old chain")
	}
	if newBlock == nil {
		return fmt.Errorf("invalid new chain")
	}
	// Both sides of the reorg are at the same number, reduce both until the common
	// ancestor is found
	for {
		// If the common ancestor was found, bail out
		if oldBlock.Hash() == newBlock.Hash() {
			commonBlock = oldBlock
			break
		}
		// Remove an old block as well as stash away a new block
		oldChain = append(oldChain, oldBlock)
		deletedTxs = append(deletedTxs, oldBlock.Transactions()...)
		collectLogs(oldBlock.Hash(), true)

		newChain = append(newChain, newBlock)

		// Step back with both chains
		oldBlock = bc.GetBlock(oldBlock.ParentHash(), oldBlock.NumberU64()-1)
		if oldBlock == nil {
			return fmt.Errorf("invalid old chain")
		}
		newBlock = bc.GetBlock(newBlock.ParentHash(), newBlock.NumberU64()-1)
		if newBlock == nil {
			return fmt.Errorf("invalid new chain")
		}
	}
	// Ensure the user sees large reorgs
	if len(oldChain) > 0 && len(newChain) > 0 {
		logFn := log.Info
		msg := "Chain reorg detected"
		if len(oldChain) > 63 {
			msg = "Large chain reorg detected"
			logFn = log.Warn
		}
		logFn(msg, "number", commonBlock.Number(), "hash", commonBlock.Hash(),
			"drop", len(oldChain), "dropfrom", oldChain[0].Hash(), "add", len(newChain), "addfrom", newChain[0].Hash())
		blockReorgAddMeter.Mark(int64(len(newChain)))
		blockReorgDropMeter.Mark(int64(len(oldChain)))
	} else {
		log.Error("Impossible reorg, please file an issue", "oldnum", oldBlock.Number(), "oldhash", oldBlock.Hash(), "newnum", newBlock.Number(), "newhash", newBlock.Hash())
	}
	// Insert the new chain(except the head block(reverse order)),
	// taking care of the proper incremental order.
	for i := len(newChain) - 1; i >= 1; i-- {
		// Insert the block in the canonical way, re-writing history
		bc.writeHeadBlock(newChain[i])

		// Collect reborn logs due to chain reorg
		collectLogs(newChain[i].Hash(), false)

		// Write lookup entries for hash based transaction/receipt searches
		rawdb.WriteTxLookupEntriesByBlock(bc.db, newChain[i])
		addedTxs = append(addedTxs, newChain[i].Transactions()...)
	}
	// When transactions get deleted from the database, the receipts that were
	// created in the fork must also be deleted
	batch := bc.db.NewBatch()
	for _, tx := range types.TxDifference(deletedTxs, addedTxs) {
		rawdb.DeleteTxLookupEntry(batch, tx.Hash())
	}
	// Delete any canonical number assignments above the new head
	number := bc.CurrentBlock().NumberU64()
	for i := number + 1; ; i++ {
		hash := rawdb.ReadCanonicalHash(bc.db, i)
		if hash == (common.Hash{}) {
			break
		}
		rawdb.DeleteCanonicalHash(batch, i)
	}
	batch.Write()
	// If any logs need to be fired, do it now. In theory we could avoid creating
	// this goroutine if there are no events to fire, but realistcally that only
	// ever happens if we're reorging empty blocks, which will only happen on idle
	// networks where performance is not an issue either way.
	if len(deletedLogs) > 0 {
		bc.rmLogsFeed.Send(RemovedLogsEvent{mergeLogs(deletedLogs, true)})
	}
	if len(rebirthLogs) > 0 {
		bc.logsFeed.Send(mergeLogs(rebirthLogs, false))
	}
	if len(oldChain) > 0 {
		for i := len(oldChain) - 1; i >= 0; i-- {
			bc.chainSideFeed.Send(ChainSideEvent{Block: oldChain[i]})
		}
	}
	return nil
}

// futureBlocksLoop processes the 'future block' queue.
func (bc *BlockChain) futureBlocksLoop() {
	defer bc.wg.Done()
	futureTimer := time.NewTicker(5 * time.Second)
	defer futureTimer.Stop()
	for {
		select {
		case <-futureTimer.C:
			bc.procFutureBlocks()
		case <-bc.quit:
			return
		}
	}
}

// skipBlock returns 'true', if the block being imported can be skipped over, meaning
// that the block does not need to be processed but can be considered already fully 'done'.
func (bc *BlockChain) skipBlock(err error, it *insertIterator) bool {
	// We can only ever bypass processing if the only error returned by the validator
	// is ErrKnownBlock, which means all checks passed, but we already have the block
	// and state.
	if !errors.Is(err, ErrKnownBlock) {
		return false
	}
	// If we're not using snapshots, we can skip this, since we have both block
	// and (trie-) state
	if bc.snaps == nil {
		return true
	}
	var (
		header     = it.current() // header can't be nil
		parentRoot common.Hash
	)
	// If we also have the snapshot-state, we can skip the processing.
	if bc.snaps.Snapshot(header.Root) != nil {
		return true
	}
	// In this case, we have the trie-state but not snapshot-state. If the parent
	// snapshot-state exists, we need to process this in order to not get a gap
	// in the snapshot layers.
	// Resolve parent block
	if parent := it.previous(); parent != nil {
		parentRoot = parent.Root
	} else if parent = bc.GetHeaderByHash(header.ParentHash); parent != nil {
		parentRoot = parent.Root
	}
	if parentRoot == (common.Hash{}) {
		return false // Theoretically impossible case
	}
	// Parent is also missing snapshot: we can skip this. Otherwise process.
	if bc.snaps.Snapshot(parentRoot) == nil {
		return true
	}
	return false
}

// maintainTxIndex is responsible for the construction and deletion of the
// transaction index.
//
// User can use flag `txlookuplimit` to specify a "recentness" block, below
// which ancient tx indices get deleted. If `txlookuplimit` is 0, it means
// all tx indices will be reserved.
//
// The user can adjust the txlookuplimit value for each launch after fast
// sync, Geth will automatically construct the missing indices and delete
// the extra indices.
func (bc *BlockChain) maintainTxIndex(ancients uint64) {
	defer bc.wg.Done()

	// Before starting the actual maintenance, we need to handle a special case,
	// where user might init Geth with an external ancient database. If so, we
	// need to reindex all necessary transactions before starting to process any
	// pruning requests.
	if ancients > 0 {
		var from = uint64(0)
		if bc.txLookupLimit != 0 && ancients > bc.txLookupLimit {
			from = ancients - bc.txLookupLimit
		}
		rawdb.IndexTransactions(bc.db, from, ancients, bc.quit)
	}
	// indexBlocks reindexes or unindexes transactions depending on user configuration
	indexBlocks := func(tail *uint64, head uint64, done chan struct{}) {
		defer func() { done <- struct{}{} }()

		// If the user just upgraded Geth to a new version which supports transaction
		// index pruning, write the new tail and remove anything older.
		if tail == nil {
			if bc.txLookupLimit == 0 || head < bc.txLookupLimit {
				// Nothing to delete, write the tail and return
				rawdb.WriteTxIndexTail(bc.db, 0)
			} else {
				// Prune all stale tx indices and record the tx index tail
				rawdb.UnindexTransactions(bc.db, 0, head-bc.txLookupLimit+1, bc.quit)
			}
			return
		}
		// If a previous indexing existed, make sure that we fill in any missing entries
		if bc.txLookupLimit == 0 || head < bc.txLookupLimit {
			if *tail > 0 {
				// It can happen when chain is rewound to a historical point which
				// is even lower than the indexes tail, recap the indexing target
				// to new head to avoid reading non-existent block bodies.
				end := *tail
				if end > head+1 {
					end = head + 1
				}
				rawdb.IndexTransactions(bc.db, 0, end, bc.quit)
			}
			return
		}
		// Update the transaction index to the new chain state
		if head-bc.txLookupLimit+1 < *tail {
			// Reindex a part of missing indices and rewind index tail to HEAD-limit
			rawdb.IndexTransactions(bc.db, head-bc.txLookupLimit+1, *tail, bc.quit)
		} else {
			// Unindex a part of stale indices and forward index tail to HEAD-limit
			rawdb.UnindexTransactions(bc.db, *tail, head-bc.txLookupLimit+1, bc.quit)
		}
	}
	// Any reindexing done, start listening to chain events and moving the index window
	var (
		done   chan struct{}                  // Non-nil if background unindexing or reindexing routine is active.
		headCh = make(chan ChainHeadEvent, 1) // Buffered to avoid locking up the event feed
	)
	sub := bc.SubscribeChainHeadEvent(headCh)
	if sub == nil {
		return
	}
	defer sub.Unsubscribe()

	for {
		select {
		case head := <-headCh:
			if done == nil {
				done = make(chan struct{})
				go indexBlocks(rawdb.ReadTxIndexTail(bc.db), head.Block.NumberU64(), done)
			}
		case <-done:
			done = nil
		case <-bc.quit:
			return
		}
	}
}

// reportBlock logs a bad block error.
func (bc *BlockChain) reportBlock(block *types.Block, receipts types.Receipts, err error) {
	rawdb.WriteBadBlock(bc.db, block)

	var receiptString string
	for i, receipt := range receipts {
		receiptString += fmt.Sprintf("\t %d: cumulative: %v gas: %v contract: %v status: %v tx: %v logs: %v bloom: %x state: %x\n",
			i, receipt.CumulativeGasUsed, receipt.GasUsed, receipt.ContractAddress.Bech32(),
			receipt.Status, receipt.TxHash.Hex(), receipt.Logs, receipt.Bloom, receipt.PostState)
	}
	log.Error(fmt.Sprintf(`
########## BAD BLOCK #########
Chain config: %v

Number: %v
Hash: 0x%x
%v

Error: %v
##############################
`, bc.chainConfig, block.Number(), block.Hash(), receiptString, err))
}

// InsertHeaderChain attempts to insert the given header chain in to the local
// chain, possibly creating a reorg. If an error is returned, it will return the
// index number of the failing header as well an error describing what went wrong.
//
// The verify parameter can be used to fine tune whether nonce verification
// should be done or not. The reason behind the optional check is because some
// of the header retrieval mechanisms already need to verify nonces, as well as
// because nonces can be verified sparsely, not needing to check each.
func (bc *BlockChain) InsertHeaderChain(chain []*types.Header, checkFreq int) (int, error) {
	start := time.Now()
	if i, err := bc.hc.ValidateHeaderChain(chain, checkFreq); err != nil {
		return i, err
	}

	if !bc.chainmu.TryLock() {
		return 0, errChainStopped
	}
	defer bc.chainmu.Unlock()

	_, err := bc.hc.InsertHeaderChain(chain, start)
	return 0, err
}

// GetBlockHashesFromHash retrieves a number of block hashes starting at a given
// hash, fetching towards the genesis block.
func (bc *BlockChain) GetBlockHashesFromHash(hash common.Hash, max uint64) []common.Hash {
	return bc.hc.GetBlockHashesFromHash(hash, max)
}

// Config retrieves the blockchain's chain configuration.
func (bc *BlockChain) CacheConfig() *CacheConfig { return bc.cacheConfig }

// SubscribeLogsEvent registers a subscription of *types.Block.
func (bc *BlockChain) SubscribeExecuteBlocksEvent(ch chan<- *types.Block) event.Subscription {
	return bc.scope.Track(bc.BlockExecuteFeed.Subscribe(ch))
}

// SubscribeLogsEvent registers a subscription of *types.Block.
func (bc *BlockChain) SubscribeWriteStateBlocksEvent(ch chan<- *types.Block) event.Subscription {
	return bc.scope.Track(bc.BlockFeed.Subscribe(ch))
}

// EnableDBGC enable database garbage collection.
func (bc *BlockChain) EnableDBGC() {
	bc.cacheConfig.DBDisabledGC.Set(false)
}

// DisableDBGC disable database garbage collection.
func (bc *BlockChain) DisableDBGC() {
	bc.cacheConfig.DBDisabledGC.Set(true)
}
