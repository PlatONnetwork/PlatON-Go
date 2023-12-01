// Copyright 2015 The go-ethereum Authors
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

// Package downloader contains the manual full chain synchronisation.
package downloader

import (
	"errors"
	"fmt"
	"math/big"
	"sync"
	"sync/atomic"
	"time"

	ethereum "github.com/PlatONnetwork/PlatON-Go"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state/snapshot"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/eth/protocols/eth"
	"github.com/PlatONnetwork/PlatON-Go/eth/protocols/snap"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

const (
	KeyFastSyncStatus = "FastSyncStatus"
	FastSyncBegin     = 0
	FastSyncFail      = 1
	FastSyncDel       = 2
)

var (
	MaxBlockFetch   = 128 // Amount of blocks to be fetched per retrieval request
	MaxHeaderFetch  = 192 // Amount of block headers to be fetched per retrieval request
	MaxSkeletonSize = 128 // Number of header fetches to need for a skeleton assembly
	MaxReceiptFetch = 256 // Amount of transaction receipts to allow fetching per request

	maxQueuedHeaders         = 32 * 1024                         // [eth/62] Maximum number of headers to queue for import (DOS protection)
	maxHeadersProcess        = 2048                              // Number of header download results to import at once into the chain
	maxResultsProcess        = 2048                              // Number of content download results to import at once into the chain
	maxForkAncestry   uint64 = params.LightImmutabilityThreshold // Maximum chain reorganisation (locally redeclared so tests can reduce it)

	reorgProtThreshold   = 48 // Threshold number of recent blocks to disable mini reorg protection
	reorgProtHeaderDelay = 2  // Number of headers to delay delivering to cover mini reorgs

	fsHeaderCheckFrequency = 100             // Verification frequency of the downloaded headers during snap sync
	fsHeaderSafetyNet      = 0               // PlatON use PoS and safe distance is 0
	fsHeaderForceVerify    = 24              // Number of headers to verify before and after the pivot to accept it
	fsHeaderContCheck      = 3 * time.Second // Time interval to check for header continuations during state download
	fsMinFullBlocks        = 64              // Number of blocks to retrieve fully even in snap sync
)

var (
	errBusy                    = errors.New("busy")
	errUnknownPeer             = errors.New("peer is unknown or unhealthy")
	errBadPeer                 = errors.New("action from bad peer ignored")
	errStallingPeer            = errors.New("peer is stalling")
	errUnsyncedPeer            = errors.New("unsynced peer")
	errNoPeers                 = errors.New("no peers to keep download active")
	errTimeout                 = errors.New("timeout")
	errEmptyHeaderSet          = errors.New("empty header set by peer")
	errPeersUnavailable        = errors.New("no peers available or all tried for download")
	errInvalidAncestor         = errors.New("retrieved ancestor is invalid")
	errInvalidChain            = errors.New("retrieved hash chain is invalid")
	errInvalidBody             = errors.New("retrieved block body is invalid")
	errInvalidReceipt          = errors.New("retrieved receipt is invalid")
	errCancelStateFetch        = errors.New("state data download canceled (requested)")
	errCancelContentProcessing = errors.New("content processing canceled (requested)")
	errCanceled                = errors.New("syncing canceled (requested)")
	errNoSyncActive            = errors.New("no sync active")
	errTooOld                  = errors.New("peer's protocol version too old")
	//errTooLess                 = errors.New("peer's hight less than 64")
)

// peerDropFn is a callback type for dropping a peer detected as malicious.
type peerDropFn func(id string)

// headerTask is a set of downloaded headers to queue along with their precomputed
// hashes to avoid constant rehashing.
type headerTask struct {
	headers []*types.Header
	hashes  []common.Hash
}

type Downloader struct {
	mode uint32         // Synchronisation mode defining the strategy used (per sync cycle), use d.getMode() to get the SyncMode
	mux  *event.TypeMux // Event multiplexer to announce sync operation events

	queue *queue   // Scheduler for selecting the hashes to download
	peers *peerSet // Set of active peers from which download can proceed

	stateDB    ethdb.Database // Database to state sync into (and deduplicate via)
	snapshotDB snapshotdb.DB

	// Statistics
	syncStatsChainOrigin uint64       // Origin block number where syncing started at
	syncStatsChainHeight uint64       // Highest block number known when syncing started
	syncStatsLock        sync.RWMutex // Lock protecting the sync stats fields

	lightchain LightChain
	blockchain BlockChain

	// Callbacks
	dropPeer peerDropFn // Drops a peer for misbehaving

	// Status
	synchroniseMock func(id string, hash common.Hash) error // Replacement for synchronise during testing
	synchronising   int32
	notified        int32
	committed       int32
	ancientLimit    uint64 // The maximum block number which can be regarded as ancient data.

	// Channels
	headerProcCh     chan *headerTask // Channel to feed the header processor new tasks
	pposInfoCh       chan dataPack    // Channel receiving inbound ppos storage
	pposStorageCh    chan dataPack    // Channel receiving inbound ppos storage
	originAndPivotCh chan dataPack    // Channel receiving origin and pivot block

	// State sync
	pivotHeader *types.Header // Pivot block header to dynamically push the syncing state root
	pivotLock   sync.RWMutex  // Lock protecting pivot header reads from updates

	SnapSyncer     *snap.Syncer // TODO(karalabe): make private! hack for now
	stateSyncStart chan *stateSync

	// Cancellation and termination
	cancelPeer string         // Identifier of the peer currently being used as the master (cancel on drop)
	cancelCh   chan struct{}  // Channel to cancel mid-flight syncs
	cancelLock sync.RWMutex   // Lock to protect the cancel channel and peer in delivers
	cancelWg   sync.WaitGroup // Make sure all fetcher goroutines have exited.

	quitCh   chan struct{} // Quit channel to signal termination
	quitLock sync.Mutex    // Lock to prevent double closes

	// Testing hooks
	syncInitHook     func(uint64, uint64)  // Method to call upon initiating a new sync run
	bodyFetchHook    func([]*types.Header) // Method to call upon starting a block body fetch
	receiptFetchHook func([]*types.Header) // Method to call upon starting a receipt fetch
	chainInsertHook  func([]*fetchResult)  // Method to call upon inserting a chain of blocks (possibly in multiple invocations)
}

// LightChain encapsulates functions required to synchronise a light chain.
type LightChain interface {
	// HasHeader verifies a header's presence in the local chain.
	HasHeader(common.Hash, uint64) bool

	// GetHeaderByHash retrieves a header from the local chain.
	GetHeaderByHash(common.Hash) *types.Header

	// CurrentHeader retrieves the head header from the local chain.
	CurrentHeader() *types.Header

	// InsertHeaderChain inserts a batch of headers into the local chain.
	InsertHeaderChain([]*types.Header, int) (int, error)

	// SetHead rewinds the local chain to a new head.
	SetHead(uint64) error
}

// BlockChain encapsulates functions required to sync a (full or snap) blockchain.
type BlockChain interface {
	LightChain

	// HasBlock verifies a block's presence in the local chain.
	HasBlock(common.Hash, uint64) bool

	// GetBlockByHash retrieves a block from the local chain.
	GetBlockByHash(common.Hash) *types.Block

	// CurrentBlock retrieves the head block from the local chain.
	CurrentBlock() *types.Block

	// CurrentFastBlock retrieves the head snap block from the local chain.
	CurrentFastBlock() *types.Block

	// SnapSyncCommitHead directly commits the head block to a certain entity.
	SnapSyncCommitHead(common.Hash) error

	// InsertChain inserts a batch of blocks into the local chain.
	InsertChain(types.Blocks) (int, error)

	// InsertReceiptChain inserts a batch of receipts into the local chain.
	InsertReceiptChain(types.Blocks, []types.Receipts, uint64) (int, error)

	// Snapshots returns the blockchain snapshot tree to paused it during sync.
	Snapshots() *snapshot.Tree
}

// New creates a new downloader to fetch hashes and blocks from remote peers.
func New(stateDb ethdb.Database, snapshotDB snapshotdb.DB, mux *event.TypeMux, chain BlockChain, lightchain LightChain, dropPeer peerDropFn, decodeExtra decodeExtraFn) *Downloader {
	if lightchain == nil {
		lightchain = chain
	}
	dl := &Downloader{
		stateDB:          stateDb,
		mux:              mux,
		queue:            newQueue(blockCacheMaxItems, blockCacheInitialItems, decodeExtra),
		peers:            newPeerSet(),
		blockchain:       chain,
		lightchain:       lightchain,
		dropPeer:         dropPeer,
		headerProcCh:     make(chan *headerTask, 1),
		pposStorageCh:    make(chan dataPack, 1),
		pposInfoCh:       make(chan dataPack, 1),
		originAndPivotCh: make(chan dataPack, 1),
		quitCh:           make(chan struct{}),
		SnapSyncer:       snap.NewSyncer(stateDb),
		stateSyncStart:   make(chan *stateSync),
		snapshotDB:       snapshotDB,
	}
	go dl.stateFetcher()
	return dl
}

// Progress retrieves the synchronisation boundaries, specifically the origin
// block where synchronisation started at (may have failed/suspended); the block
// or header sync is currently at; and the latest known block which the sync targets.
//
// In addition, during the state download phase of snap synchronisation the number
// of processed and the total number of known states are also returned. Otherwise
// these are zero.
func (d *Downloader) Progress() ethereum.SyncProgress {
	// Lock the current stats and return the progress
	d.syncStatsLock.RLock()
	defer d.syncStatsLock.RUnlock()

	current := uint64(0)
	mode := d.getMode()
	switch {
	case d.blockchain != nil && mode == FullSync:
		current = d.blockchain.CurrentBlock().NumberU64()
	case d.blockchain != nil && mode == SnapSync:
		current = d.blockchain.CurrentFastBlock().NumberU64()
	case d.lightchain != nil:
		current = d.lightchain.CurrentHeader().Number.Uint64()
	default:
		log.Error("Unknown downloader chain/mode combo", "light", d.lightchain != nil, "full", d.blockchain != nil, "mode", mode)
	}
	progress, pending := d.SnapSyncer.Progress()

	return ethereum.SyncProgress{
		StartingBlock:       d.syncStatsChainOrigin,
		CurrentBlock:        current,
		HighestBlock:        d.syncStatsChainHeight,
		SyncedAccounts:      progress.AccountSynced,
		SyncedAccountBytes:  uint64(progress.AccountBytes),
		SyncedBytecodes:     progress.BytecodeSynced,
		SyncedBytecodeBytes: uint64(progress.BytecodeBytes),
		SyncedStorage:       progress.StorageSynced,
		SyncedStorageBytes:  uint64(progress.StorageBytes),
		HealedTrienodes:     progress.TrienodeHealSynced,
		HealedTrienodeBytes: uint64(progress.TrienodeHealBytes),
		HealedBytecodes:     progress.BytecodeHealSynced,
		HealedBytecodeBytes: uint64(progress.BytecodeHealBytes),
		HealingTrienodes:    pending.TrienodeHeal,
		HealingBytecode:     pending.BytecodeHeal,
	}
}

// Synchronising returns whether the downloader is currently retrieving blocks.
func (d *Downloader) Synchronising() bool {
	return atomic.LoadInt32(&d.synchronising) > 0
}

// RegisterPeer injects a new download peer into the set of block source to be
// used for fetching hashes and blocks from.
func (d *Downloader) RegisterPeer(id string, version uint, peer Peer) error {
	var logger log.Logger
	if len(id) < 16 {
		// Tests use short IDs, don't choke on them
		logger = log.New("peer", id)
	} else {
		logger = log.New("peer", id[:8])
	}
	logger.Trace("Registering sync peer")
	if err := d.peers.Register(newPeerConnection(id, version, peer, logger)); err != nil {
		logger.Error("Failed to register sync peer", "err", err)
		return err
	}
	return nil
}

// RegisterLightPeer injects a light client peer, wrapping it so it appears as a regular peer.
func (d *Downloader) RegisterLightPeer(id string, version uint, peer LightPeer) error {
	return d.RegisterPeer(id, version, &lightPeerWrapper{peer})
}

// UnregisterPeer remove a peer from the known list, preventing any action from
// the specified peer. An effort is also made to return any pending fetches into
// the queue.
func (d *Downloader) UnregisterPeer(id string) error {
	// Unregister the peer from the active peer set and revoke any fetch tasks
	var logger log.Logger
	if len(id) < 16 {
		// Tests use short IDs, don't choke on them
		logger = log.New("peer", id)
	} else {
		logger = log.New("peer", id[:8])
	}
	logger.Trace("Unregistering sync peer")
	if err := d.peers.Unregister(id); err != nil {
		logger.Error("Failed to unregister sync peer", "err", err)
		return err
	}
	d.queue.Revoke(id)

	return nil
}

// Synchronise tries to sync up our local block chain with a remote peer, both
// adding various sanity checks as well as wrapping it with various log entries.
func (d *Downloader) Synchronise(id string, head common.Hash, bn *big.Int, mode SyncMode) error {
	log.Debug("Synchronise from other peer", "peerID", id, "head", head, "bn", bn, "mode", mode)
	err := d.synchronise(id, head, bn, mode)
	switch err {
	case nil, errBusy, errCanceled:
		return err
	}
	if errors.Is(err, errInvalidChain) || errors.Is(err, errBadPeer) || errors.Is(err, errTimeout) ||
		errors.Is(err, errStallingPeer) || errors.Is(err, errUnsyncedPeer) || errors.Is(err, errEmptyHeaderSet) ||
		errors.Is(err, errPeersUnavailable) || errors.Is(err, errTooOld) || errors.Is(err, errInvalidAncestor) {
		log.Warn("Synchronisation failed, dropping peer", "peer", id, "err", err)
		if d.dropPeer == nil {
			// The dropPeer method is nil when `--copydb` is used for a local copy.
			// Timeouts can occur if e.g. compaction hits at the wrong time, and can be ignored
			log.Warn("Downloader wants to drop peer, but peerdrop-function is not set", "peer", id)
		} else {
			d.dropPeer(id)
		}
		return err
	}
	log.Warn("Synchronisation failed, retrying", "err", err)
	return err
}

// synchronise will select the peer and use it for synchronising. If an empty string is given
// it will use the best peer possible and synchronize if its TD is higher than our own. If any of the
// checks fail an error will be returned. This method is synchronous
func (d *Downloader) synchronise(id string, hash common.Hash, bn *big.Int, mode SyncMode) error {
	// Mock out the synchronisation if testing
	if d.synchroniseMock != nil {
		return d.synchroniseMock(id, hash)
	}
	// Make sure only one goroutine is ever allowed past this point at once
	if !atomic.CompareAndSwapInt32(&d.synchronising, 0, 1) {
		return errBusy
	}
	defer atomic.StoreInt32(&d.synchronising, 0)

	// Post a user notification of the sync (only once per session)
	if atomic.CompareAndSwapInt32(&d.notified, 0, 1) {
		log.Info("Block synchronisation started")
	}
	// If snap sync was requested, create the snap scheduler and switch to snap
	// sync mode. Long term we could drop snap sync or merge the two together,
	// but until snap becomes prevalent, we should support both. TODO(karalabe).
	if mode == SnapSync {
		// Snap sync uses the snapshot namespace to store potentially flakey data until
		// sync completely heals and finishes. Pause snapshot maintenance in the mean-
		// time to prevent access.
		if snapshots := d.blockchain.Snapshots(); snapshots != nil { // Only nil in tests
			snapshots.Disable()
		}
	}
	// Reset the queue, peer set and wake channels to clean any internal leftover state
	d.queue.Reset(blockCacheMaxItems, blockCacheInitialItems)
	d.peers.Reset()

	for _, ch := range []chan bool{d.queue.blockWakeCh, d.queue.receiptWakeCh} {
		select {
		case <-ch:
		default:
		}
	}
	for empty := false; !empty; {
		select {
		case <-d.headerProcCh:
		default:
			empty = true
		}
	}
	for _, ch := range []chan dataPack{d.pposInfoCh, d.pposStorageCh, d.originAndPivotCh} {
		for empty := false; !empty; {
			select {
			case <-ch:
			default:
				empty = true
			}
		}
	}

	// Create cancel channel for aborting mid-flight and mark the master peer
	d.cancelLock.Lock()
	d.cancelCh = make(chan struct{})
	d.cancelPeer = id
	d.cancelLock.Unlock()

	defer d.Cancel() // No matter what, we can't leave the cancel channel open

	// Atomically set the requested sync mode
	atomic.StoreUint32(&d.mode, uint32(mode))

	// Retrieve the origin peer and initiate the downloading process
	p := d.peers.Peer(id)
	if p == nil {
		return errUnknownPeer
	}
	return d.syncWithPeer(p, hash, bn)
}

func (d *Downloader) getMode() SyncMode {
	return SyncMode(atomic.LoadUint32(&d.mode))
}

// syncWithPeer starts a block synchronization based on the hash chain from the
// specified peer and head hash.
func (d *Downloader) syncWithPeer(p *peerConnection, hash common.Hash, bn *big.Int) (err error) {
	d.mux.Post(StartEvent{})
	defer func() {
		// reset on error
		if err != nil {
			d.mux.Post(FailedEvent{err})
		} else {
			d.mux.Post(DoneEvent{})
		}
	}()
	if p.version < eth.ETH66 {
		return fmt.Errorf("%w: advertized %d < required %d", errTooOld, p.version, eth.ETH66)
	}
	mode := d.getMode()

	log.Debug("Synchronising with the network", "peer", p.id, "eth", p.version, "head", hash, "bn", bn, "mode", mode)
	defer func(start time.Time) {
		log.Debug("Synchronisation terminated", "elapsed", common.PrettyDuration(time.Since(start)), "err", err)
	}(time.Now())

	// 先检查对端与我们是否是同一个链
	originh, snapshotBase, err := d.findOrigin(p)
	if err != nil {
		log.Error("findOrigin error", "err", err.Error())
		return err
	}
	origin := originh.Number.Uint64()
	log.Info("synchronising findOrigin", "peer", p.id, "origin", originh.Number, "base", snapshotBase.Number)

	// 查询对端节点的高度,p点
	// Look up the sync boundaries: the common ancestor and the target block
	latest, pivot, err := d.fetchHead(p)
	if err != nil {
		return err
	}
	if pivot == nil {
		log.Info("synchronising fetchHead", "peer", p.id, "latest", latest.Number)
	} else {
		log.Info("synchronising fetchHead", "peer", p.id, "latest", latest.Number, "pivot", pivot.Number)
	}

	if mode == SnapSync && pivot == nil {
		// If no pivot block was returned, the head is below the min full block
		// threshold (i.e. new chain). In that case we won't really snap sync
		// anyway, but still need a valid pivot block to avoid some code hitting
		// nil panics on an access.
		pivot = d.blockchain.CurrentBlock().Header()
	}
	height := latest.Number.Uint64()

	d.syncStatsLock.Lock()
	if d.syncStatsChainHeight <= origin || d.syncStatsChainOrigin > origin {
		d.syncStatsChainOrigin = origin
	}
	d.syncStatsChainHeight = height
	d.syncStatsLock.Unlock()

	d.committed = 1
	if mode == SnapSync && pivot.Number.Uint64() != 0 {
		d.committed = 0
	}
	if mode == SnapSync {
		// Set the ancient data limitation.
		// If we are running snap sync, all block data older than ancientLimit will be
		// written to the ancient store. More recent data will be written to the active
		// database and will wait for the freezer to migrate.
		//
		// If there is a checkpoint available, then calculate the ancientLimit through
		// that. Otherwise calculate the ancient limit through the advertised height
		// of the remote peer.
		//
		// The reason for picking checkpoint first is that a malicious peer can give us
		// a fake (very high) height, forcing the ancient limit to also be very high.
		// The peer would start to feed us valid blocks until head, resulting in all of
		// the blocks might be written into the ancient store. A following mini-reorg
		// could cause issues.
		//todo checkpoint mod need add
		if height > maxForkAncestry+1 {
			d.ancientLimit = height - maxForkAncestry - 1
		} else {
			d.ancientLimit = 0
		}
		frozen, _ := d.stateDB.Ancients() // Ignore the error here since light client can also hit here.

		// If a part of blockchain data has already been written into active store,
		// disable the ancient style insertion explicitly.
		if origin >= frozen && frozen != 0 {
			d.ancientLimit = 0
			log.Info("Disabling direct-ancient mode", "origin", origin, "ancient", frozen-1)
		} else if d.ancientLimit > 0 {
			log.Debug("Enabling direct-ancient mode", "ancient", d.ancientLimit)
		}
		// Rewind the ancient store and blockchain if reorg happens.
		// we don't have lightchain,so should not Rollback hash
		/*if origin+1 < frozen {
			var hashes []common.Hash
			if err := d.lightchain.SetHead(origin + 1); err != nil {
				return err
			}
			d.lightchain.Rollback(hashes)
		}*/
	}
	// Initiate the sync using a concurrent header and content retrieval algorithm
	d.queue.Prepare(origin+1, mode)
	if d.syncInitHook != nil {
		d.syncInitHook(origin, height)
	}
	fetchers := []func() error{
		func() error { return d.fetchHeaders(p, origin+1, latest.Number.Uint64()) }, // Headers are always retrieved
		func() error { return d.fetchBodies(origin + 1) },                           // Bodies are retrieved during normal and fast sync
		func() error { return d.fetchReceipts(origin + 1) },                         // Receipts are retrieved during fast sync
		func() error { return d.processHeaders(origin+1, bn) },
	}
	if mode == SnapSync {
		if err := d.setFastSyncStatus(FastSyncBegin); err != nil {
			return err
		}

		d.pivotLock.Lock()
		d.pivotHeader = pivot
		d.pivotLock.Unlock()

		fetchers = append(fetchers, func() error { return d.processSnapSyncContent() })
	} else if mode == FullSync {
		fetchers = append(fetchers, d.processFullSyncContent)
	}
	return d.spawnSync(fetchers)
}

// setFastSyncStatus set status to snapshot db when fast sync begin
// if  the user close platon when sync not finish,set status fail
// if the sync is complete,will del the key
func (d *Downloader) setFastSyncStatus(status uint16) error {
	key := []byte(KeyFastSyncStatus)
	log.Debug("set fast sync status", "status", status)
	switch status {
	case FastSyncDel:
		if err := d.snapshotDB.DelBaseDB(key); err != nil {
			log.Error("del fast sync status from snapshotdb fail", "err", err)
			return err
		}
	case FastSyncBegin, FastSyncFail:
		syncStatus := [2][]byte{
			key,
			common.Uint16ToBytes(status),
		}
		if err := d.snapshotDB.WriteBaseDB([][2][]byte{syncStatus}); err != nil {
			log.Error("save fast sync status to snapshotdb fail", "err", err)
			return err
		}
	default:
		return errors.New("status is not supports")
	}
	return nil
}

func (d *Downloader) fetchPPOSStorageV2(p *peerConnection, pivot *types.Header) (err error) {
	p.log.Info("begin  fetchPPOSStorageV2  from remote peer", "pivot number", pivot.Number)

	timeout := time.NewTimer(0) // timer to dump a non-responsive active peer
	<-timeout.C                 // timeout channel should be initially empty
	defer timeout.Stop()
	var ttl time.Duration
	ttl = time.Second * 15
	timeout.Reset(ttl)

	// 先判断对端 base num 是否在我们p点之下
	p.log.Info("begin  fetch  base num  from remote peer", "pivot number", pivot.Number)
	go p.peer.RequestOriginAndPivotByCurrent(pivot.Number.Uint64())
	select {
	case <-d.cancelCh:
		return errCanceled

	case packet := <-d.originAndPivotCh:
		// Discard anything not from the origin peer
		if packet.PeerId() != p.id {
			p.log.Error("Received headers from incorrect peer", "peer", packet.PeerId())
			return errUnknownPeer
		}
		// Make sure the peer actually gave something valid
		headers := packet.(*headerPack).headers
		if len(headers) == 0 {
			p.log.Error("Empty head header set")
			return errEmptyHeaderSet
		}
		if len(headers) != 2 {
			p.log.Error("header length wrong", "len", len(headers))
			return errors.New("header length wrong")
		}
		if headers[0] == nil {
			p.log.Error("not find  current block")
			return errors.New("not find  current block")
		}
		if headers[0].Number.Uint64() != pivot.Number.Uint64() || headers[0].Hash() != pivot.Hash() {
			p.log.Error("retrieved hash chain is invalid", "pivot num", pivot.Number, "remote current num", headers[0].Number.Uint64(), "pivot hash", pivot.Hash(), "remote current hash", headers[0].Hash())
			return errInvalidChain
		}
		if headers[1] == nil {
			return errors.New("not find pivot")
		}

		if headers[1].Number.Uint64() >= pivot.Number.Uint64() {
			return fmt.Errorf("the remote peer base num %d is greater than our request pivot %d", headers[1].Number.Uint64(), pivot.Number.Uint64())
		}
		p.log.Info("finish  fetch  base num  from remote peer", "pivot number", pivot.Number, "base", headers[1].Number)
	case <-timeout.C:
		p.log.Error("Waiting for head header timed out", "elapsed", ttl)
		return errTimeout
	}

	// 获取快照数据
	go p.peer.RequestPPOSStorage(pivot.Number.Uint64())

	var (
		count  int64
		blocks []snapshotdb.BlockData
	)

	ttl = time.Second * 15
	timeout.Reset(ttl)
	if err := d.snapshotDB.SetEmpty(); err != nil {
		return errors.New("set snapshotDB empty fail:" + err.Error())
	}
	if err := d.setFastSyncStatus(FastSyncBegin); err != nil {
		return err
	}

	for {
		select {
		case <-d.cancelCh:
			return errCanceled
		case <-timeout.C:
			p.log.Error("failed to fetch PPOSStorageV2,waiting for ppos storage V2 timed out,elapsed", "ttl", ttl.String())
			return errTimeout
		case packet := <-d.pposStorageCh:
			// Discard anything not from the origin peer
			if packet.PeerId() != p.id {
				return fmt.Errorf("received ppos storage V2 from incorrect peer,%s", packet.PeerId())
			}
			pposDada := packet.(*pposStoragePack)
			if pposDada.base {
				count += int64(len(pposDada.kvs))
				if uint64(count) != pposDada.kvNum {
					return fmt.Errorf("received ppos storage v2 from incorrect kvNum %v,count %v", pposDada.kvNum, count)
				}

				if err := d.snapshotDB.WriteBaseDB(pposDada.KVs()); err != nil {
					return fmt.Errorf("write to base db fail,%v", err)
				}
				if pposDada.last {
					p.log.Info("fetch PPOSStorageV2 base db part has finish")
				}
			} else {
				if pposDada.baseBlock == pivot.Number.Uint64() {
					p.log.Info("fetch PPOSStorageV2 block   has finish", "pivot", pivot.Number, "baseBlock", pposDada.baseBlock, "receive", len(blocks))
					if err := d.snapshotDB.WriteBaseDBWithBlock(pivot, pposDada.blocks); err != nil {
						return fmt.Errorf("updataBaseDBWithBlock fail, %v", err)
					}
					p.log.Info("fetch PPOSStorageV2 block  write success", "pivot", pivot.Number, "baseBlock", pposDada.baseBlock, "receive", len(blocks))
					return nil
				} else {
					if len(pposDada.blocks) > 0 {
						p.log.Info("begin  write PPOSStorageV2  part 2", "pivot number", pivot.Number, "len", len(pposDada.blocks), "begin", pposDada.blocks[0].Number.Uint64(), "end", pposDada.blocks[len(pposDada.blocks)-1].Number.Uint64())
						blocks = append(blocks, pposDada.blocks...)
						if pposDada.blocks[len(pposDada.blocks)-1].Number.Uint64() == pivot.Number.Uint64() {
							if pivot.Number.Uint64()-pposDada.baseBlock != uint64(len(blocks)) {
								return fmt.Errorf("pposDada is less than get,pivot %v,baseBlock %v,get %v", pivot.Number, pposDada.baseBlock, len(blocks))
							}
							p.log.Info("fetch PPOSStorageV2 block   has finish", "pivot", pivot.Number, "baseBlock", pposDada.baseBlock, "receive", len(blocks))
							if err := d.snapshotDB.WriteBaseDBWithBlock(pivot, pposDada.blocks); err != nil {
								return fmt.Errorf("updataBaseDBWithBlock fail, %v", err)
							}
							p.log.Info("fetch PPOSStorageV2 block  write success", "pivot", pivot.Number, "baseBlock", pposDada.baseBlock, "receive", len(blocks))
							return nil
						}
					}
				}
			}
			ttl = d.peers.rates.TargetTimeout()
			timeout.Reset(ttl)
		}
	}
}

// spawnSync runs d.process and all given fetcher functions to completion in
// separate goroutines, returning the first error that appears.
func (d *Downloader) spawnSync(fetchers []func() error) error {
	errc := make(chan error, len(fetchers))
	d.cancelWg.Add(len(fetchers))
	for _, fn := range fetchers {
		fn := fn
		go func() { defer d.cancelWg.Done(); errc <- fn() }()
	}
	// Wait for the first error, then terminate the others.
	var (
		err    error
		failed bool
	)
	for i := 0; i < len(fetchers); i++ {
		if i == len(fetchers)-1 {
			// Close the queue when all fetchers have exited.
			// This will cause the block processor to end when
			// it has processed the queue.
			d.queue.Close()
		}
		if err = <-errc; err != nil {
			failed = true
			if err != errCanceled {
				break
			}
		}
	}
	mode := d.getMode()
	current := d.blockchain.CurrentBlock().NumberU64()
	if mode == SnapSync {
		if failed && current == 0 {
			if err := d.setFastSyncStatus(FastSyncFail); err != nil {
				return err
			}
		} else {
			if err := d.setFastSyncStatus(FastSyncDel); err != nil {
				return err
			}
		}
	}
	if mode == FullSync {
		if err := d.setFastSyncStatus(FastSyncDel); err != nil {
			return err
		}
	}

	d.queue.Close()
	d.Cancel()
	return err
}

// cancel aborts all of the operations and resets the queue. However, cancel does
// not wait for the running download goroutines to finish. This method should be
// used when cancelling the downloads from inside the downloader.
func (d *Downloader) cancel() {
	// Close the current cancel channel
	d.cancelLock.Lock()
	defer d.cancelLock.Unlock()

	if d.cancelCh != nil {
		select {
		case <-d.cancelCh:
			// Channel was already closed
		default:
			close(d.cancelCh)
		}
	}
}

// Cancel aborts all of the operations and waits for all download goroutines to
// finish before returning.
func (d *Downloader) Cancel() {
	d.cancel()
	d.cancelWg.Wait()

	d.ancientLimit = 0
	log.Debug("Reset ancient limit to zero")
}

// Terminate interrupts the downloader, canceling all pending operations.
// The downloader cannot be reused after calling Terminate.
func (d *Downloader) Terminate() {
	// Close the termination channel (make sure double close is allowed)
	d.quitLock.Lock()
	select {
	case <-d.quitCh:
	default:
		close(d.quitCh)
	}
	d.quitLock.Unlock()

	// Cancel any pending download requests
	d.Cancel()
}

// origin is the this chain  current block header ,compare remote  the same num of header in remote,
// the pivot header is the remote peer  snapshotDB base num
func (d *Downloader) findOrigin(p *peerConnection) (*types.Header, *types.Header, error) {
	var current *types.Header
	mode := d.getMode()
	if mode == FullSync {
		current = d.blockchain.CurrentBlock().Header()
	} else if mode == SnapSync {
		current = d.blockchain.CurrentFastBlock().Header()
	} else {
		current = d.lightchain.CurrentHeader()
	}
	currentNumber := current.Number.Uint64()
	go p.peer.RequestOriginAndPivotByCurrent(currentNumber)
	ttl := time.Second * 10
	timeout := time.After(ttl)
	for {
		select {
		case <-d.cancelCh:
			return nil, nil, errCanceled

		case packet := <-d.originAndPivotCh:
			// Discard anything not from the origin peer
			if packet.PeerId() != p.id {
				p.log.Error("Received headers from incorrect peer", "peer", packet.PeerId())
				return nil, nil, errUnknownPeer
			}
			// Make sure the peer actually gave something valid
			headers := packet.(*headerPack).headers
			if len(headers) == 0 {
				p.log.Error("Empty head header set")
				return nil, nil, errEmptyHeaderSet
			}
			if len(headers) != 2 {
				p.log.Error("header length wrong", "len", len(headers))
				return nil, nil, errors.New("header length wrong")
			}
			if headers[0] == nil {
				p.log.Error("not find  current block")
				return nil, nil, errors.New("not find  current block")
			}
			if headers[0].Number.Uint64() != currentNumber || headers[0].Hash() != current.Hash() {
				p.log.Error("retrieved hash chain is invalid", "current num", currentNumber, "remote current num", headers[0].Number.Uint64(), "current hash", current.Hash(), "remote current hash", headers[0].Hash())
				return nil, nil, errInvalidChain
			}
			if headers[1] == nil {
				return nil, nil, errors.New("not find pivot")
			}
			return headers[0], headers[1], nil
		case <-timeout:
			p.log.Error("Waiting for head header timed out", "elapsed", ttl)
			return nil, nil, errTimeout
		case <-d.pposStorageCh:
			// Out of bounds delivery, ignore
		}
	}
}

// fetchHead retrieves the head header and prior pivot block (if available) from
// a remote peer.
func (d *Downloader) fetchHead(p *peerConnection) (head *types.Header, pivot *types.Header, err error) {
	p.log.Debug("Retrieving remote chain head")
	mode := d.getMode()

	// Request the advertised remote head block and wait for the response
	latest, _ := p.peer.Head()
	fetch := 1
	if mode == SnapSync {
		fetch = 2 // head + pivot headers
	}
	headers, hashes, err := d.fetchHeadersByHash(p, latest, fetch, fsMinFullBlocks-1, true)
	if err != nil {
		return nil, nil, err
	}
	// Make sure the peer gave us at least one and at most the requested headers
	if len(headers) == 0 || len(headers) > fetch {
		return nil, nil, fmt.Errorf("%w: returned headers %d != requested %d", errBadPeer, len(headers), fetch)
	}
	// The first header needs to be the head, validate against the checkpoint
	// and request. If only 1 header was returned, make sure there's no pivot
	// or there was not one requested.
	head = headers[0]
	if len(headers) == 1 {
		if mode == SnapSync && head.Number.Uint64() > uint64(fsMinFullBlocks) {
			return nil, nil, fmt.Errorf("%w: no pivot included along head header", errBadPeer)
		}
		p.log.Debug("Remote head identified, no pivot", "number", head.Number, "hash", hashes[0])
		return head, nil, nil
	}
	// At this point we have 2 headers in total and the first is the
	// validated head of the chain. Check the pivot number and return,
	pivot = headers[1]
	if pivot.Number.Uint64() != head.Number.Uint64()-uint64(fsMinFullBlocks) {
		return nil, nil, fmt.Errorf("%w: remote pivot %d != requested %d", errInvalidChain, pivot.Number, head.Number.Uint64()-uint64(fsMinFullBlocks))
	}
	return head, pivot, nil
}

// fetchHeaders keeps retrieving headers concurrently from the number
// requested, until no more are returned, potentially throttling on the way. To
// facilitate concurrency but still protect against malicious nodes sending bad
// headers, we construct a header chain skeleton using the "origin" peer we are
// syncing with, and fill in the missing headers using anyone else. Headers from
// other peers are only accepted if they map cleanly to the skeleton. If no one
// can fill in the skeleton - not even the origin peer - it's assumed invalid and
// the origin is dropped.
func (d *Downloader) fetchHeaders(p *peerConnection, from uint64, head uint64) error {
	p.log.Debug("Directing header downloads", "origin", from)
	defer p.log.Debug("Header download terminated")

	// Start pulling the header chain skeleton until all is done
	var (
		skeleton = true  // Skeleton assembly phase or finishing up
		pivoting = false // Whether the next request is pivot verification
		ancestor = from
		mode     = d.getMode()
	)
	for {
		// Pull the next batch of headers, it either:
		//   - Pivot check to see if the chain moved too far
		//   - Skeleton retrieval to permit concurrent header fetches
		//   - Full header retrieval if we're near the chain head
		var (
			headers []*types.Header
			hashes  []common.Hash
			err     error
		)
		switch {
		case pivoting:
			d.pivotLock.RLock()
			pivot := d.pivotHeader.Number.Uint64()
			d.pivotLock.RUnlock()

			p.log.Trace("Fetching next pivot header", "number", pivot+uint64(fsMinFullBlocks))
			headers, hashes, err = d.fetchHeadersByNumber(p, pivot+uint64(fsMinFullBlocks), 2, fsMinFullBlocks-9, false) // move +64 when it's 2x64-8 deep

		case skeleton:
			p.log.Trace("Fetching skeleton headers", "count", MaxHeaderFetch, "from", from)
			headers, hashes, err = d.fetchHeadersByNumber(p, from+uint64(MaxHeaderFetch)-1, MaxSkeletonSize, MaxHeaderFetch-1, false)

		default:
			p.log.Trace("Fetching full headers", "count", MaxHeaderFetch, "from", from)
			headers, hashes, err = d.fetchHeadersByNumber(p, from, MaxHeaderFetch, 0, false)
		}
		switch err {
		case nil:
			// Headers retrieved, continue with processing

		case errCanceled:
			// Sync cancelled, no issue, propagate up
			return err

		default:
			// Header retrieval either timed out, or the peer failed in some strange way
			// (e.g. disconnect). Consider the master peer bad and drop
			d.dropPeer(p.id)

			// Finish the sync gracefully instead of dumping the gathered data though
			for _, ch := range []chan bool{d.queue.blockWakeCh, d.queue.receiptWakeCh} {
				select {
				case ch <- false:
				case <-d.cancelCh:
				}
			}
			select {
			case d.headerProcCh <- nil:
			case <-d.cancelCh:
			}
			return fmt.Errorf("%w: header request failed: %v", errBadPeer, err)
		}
		// If the pivot is being checked, move if it became stale and run the real retrieval
		var pivot uint64

		d.pivotLock.RLock()
		if d.pivotHeader != nil {
			pivot = d.pivotHeader.Number.Uint64()
		}
		d.pivotLock.RUnlock()

		if pivoting {
			if len(headers) == 2 {
				if have, want := headers[0].Number.Uint64(), pivot+uint64(fsMinFullBlocks); have != want {
					log.Warn("Peer sent invalid next pivot", "have", have, "want", want)
					return fmt.Errorf("%w: next pivot number %d != requested %d", errInvalidChain, have, want)
				}
				if have, want := headers[1].Number.Uint64(), pivot+2*uint64(fsMinFullBlocks)-8; have != want {
					log.Warn("Peer sent invalid pivot confirmer", "have", have, "want", want)
					return fmt.Errorf("%w: next pivot confirmer number %d != requested %d", errInvalidChain, have, want)
				}
				log.Warn("Pivot seemingly stale, moving", "old", pivot, "new", headers[0].Number)
				pivot = headers[0].Number.Uint64()

				d.pivotLock.Lock()
				d.pivotHeader = headers[0]
				d.pivotLock.Unlock()

				// Write out the pivot into the database so a rollback beyond
				// it will reenable snap sync and update the state root that
				// the state syncer will be downloading.
				rawdb.WriteLastPivotNumber(d.stateDB, pivot)
			}
			// Disable the pivot check and fetch the next batch of headers
			pivoting = false
			continue
		}
		// If the skeleton's finished, pull any remaining head headers directly from the origin
		if skeleton && len(headers) == 0 {
			// A malicious node might withhold advertised headers indefinitely
			if from+uint64(MaxHeaderFetch)-1 <= head {
				p.log.Warn("Peer withheld skeleton headers", "advertised", head, "withheld", from+uint64(MaxHeaderFetch)-1)
				return fmt.Errorf("%w: withheld skeleton headers: advertised %d, withheld #%d", errStallingPeer, head, from+uint64(MaxHeaderFetch)-1)
			}
			p.log.Debug("No skeleton, fetching headers directly")
			skeleton = false
			continue
		}
		// If no more headers are inbound, notify the content fetchers and return
		if len(headers) == 0 {
			// Don't abort header fetches while the pivot is downloading
			if atomic.LoadInt32(&d.committed) == 0 && pivot <= from {
				p.log.Debug("No headers, waiting for pivot commit")
				select {
				case <-time.After(fsHeaderContCheck):
					continue
				case <-d.cancelCh:
					return errCanceled
				}
			}
			// Pivot done (or not in snap sync) and no more headers, terminate the process
			p.log.Debug("No more headers available")
			select {
			case d.headerProcCh <- nil:
				return nil
			case <-d.cancelCh:
				return errCanceled
			}
		}
		// If we received a skeleton batch, resolve internals concurrently
		var progressed bool
		if skeleton {
			filled, hashset, proced, err := d.fillHeaderSkeleton(from, headers)
			if err != nil {
				p.log.Debug("Skeleton chain invalid", "err", err)
				return fmt.Errorf("%w: %v", errInvalidChain, err)
			}
			headers = filled[proced:]
			hashes = hashset[proced:]

			progressed = proced > 0
			from += uint64(proced)
		} else {
			// A malicious node might withhold advertised headers indefinitely
			if n := len(headers); n < MaxHeaderFetch && headers[n-1].Number.Uint64() < head {
				p.log.Warn("Peer withheld headers", "advertised", head, "delivered", headers[n-1].Number.Uint64())
				return fmt.Errorf("%w: withheld headers: advertised %d, delivered %d", errStallingPeer, head, headers[n-1].Number.Uint64())
			}
			// If we're closing in on the chain head, but haven't yet reached it, delay
			// the last few headers so mini reorgs on the head don't cause invalid hash
			// chain errors.
			if n := len(headers); n > 0 {
				// Retrieve the current head we're at
				var head uint64
				if mode == LightSync {
					head = d.lightchain.CurrentHeader().Number.Uint64()
				} else {
					head = d.blockchain.CurrentFastBlock().NumberU64()
					if full := d.blockchain.CurrentBlock().NumberU64(); head < full {
						head = full
					}
				}
				// If the head is below the common ancestor, we're actually deduplicating
				// already existing chain segments, so use the ancestor as the fake head.
				// Otherwise, we might end up delaying header deliveries pointlessly.
				if head < ancestor {
					head = ancestor
				}
				// If the head is way older than this batch, delay the last few headers
				if head+uint64(reorgProtThreshold) < headers[n-1].Number.Uint64() {
					delay := reorgProtHeaderDelay
					if delay > n {
						delay = n
					}
					log.Error("fetchHeaders reorgProtThreshold", "head", head, "last", headers[n-1].Number.Uint64(), "delay", delay)
					headers = headers[:n-delay]
					hashes = hashes[:n-delay]
				}
			}
		}
		// If no headers have bene delivered, or all of them have been delayed,
		// sleep a bit and retry. Take care with headers already consumed during
		// skeleton filling
		if len(headers) == 0 && !progressed {
			p.log.Trace("All headers delayed, waiting")
			select {
			case <-time.After(fsHeaderContCheck):
				continue
			case <-d.cancelCh:
				return errCanceled
			}
		}
		// Insert any remaining new headers and fetch the next batch
		if len(headers) > 0 {
			p.log.Trace("Scheduling new headers", "count", len(headers), "from", from)
			select {
			case d.headerProcCh <- &headerTask{
				headers: headers,
				hashes:  hashes,
			}:
			case <-d.cancelCh:
				return errCanceled
			}
			from += uint64(len(headers))
		}
		// If we're still skeleton filling snap sync, check pivot staleness
		// before continuing to the next skeleton filling
		if skeleton && pivot > 0 {
			pivoting = true
		}
	}
}

// fillHeaderSkeleton concurrently retrieves headers from all our available peers
// and maps them to the provided skeleton header chain.
//
// Any partial results from the beginning of the skeleton is (if possible) forwarded
// immediately to the header processor to keep the rest of the pipeline full even
// in the case of header stalls.
//
// The method returns the entire filled skeleton and also the number of headers
// already forwarded for processing.
func (d *Downloader) fillHeaderSkeleton(from uint64, skeleton []*types.Header) ([]*types.Header, []common.Hash, int, error) {
	log.Debug("Filling up skeleton", "from", from)
	d.queue.ScheduleSkeleton(from, skeleton)

	err := d.concurrentFetch((*headerQueue)(d))
	if err != nil {
		log.Debug("Skeleton fill failed", "err", err)
	}
	filled, hashes, proced := d.queue.RetrieveHeaders()
	if err == nil {
		log.Debug("Skeleton fill succeeded", "filled", len(filled), "processed", proced)
	}
	return filled, hashes, proced, err
}

// fetchBodies iteratively downloads the scheduled block bodies, taking any
// available peers, reserving a chunk of blocks for each, waiting for delivery
// and also periodically checking for timeouts.
func (d *Downloader) fetchBodies(from uint64) error {
	log.Debug("Downloading block bodies", "origin", from)
	err := d.concurrentFetch((*bodyQueue)(d))

	log.Debug("Block body download terminated", "err", err)
	return err
}

// fetchReceipts iteratively downloads the scheduled block receipts, taking any
// available peers, reserving a chunk of receipts for each, waiting for delivery
// and also periodically checking for timeouts.
func (d *Downloader) fetchReceipts(from uint64) error {
	log.Debug("Downloading receipts", "origin", from)
	err := d.concurrentFetch((*receiptQueue)(d))

	log.Debug("Receipt download terminated", "err", err)
	return err
}

// processHeaders takes batches of retrieved headers from an input channel and
// keeps processing and scheduling them into the header chain and downloader's
// queue until the stream ends or a failure occurs.
func (d *Downloader) processHeaders(origin uint64, bn *big.Int) error {
	// Keep a count of uncertain headers to roll back
	var (
		rollback    uint64 // Zero means no rollback (fine as you can't unroll the genesis)
		rollbackErr error
		mode        = d.getMode()
	)
	defer func() {
		// PlatON do not support rollback
		if rollback > 0 {
			lastHeader, lastFastBlock, lastBlock := d.lightchain.CurrentHeader().Number, common.Big0, common.Big0
			if mode != LightSync {
				lastFastBlock = d.blockchain.CurrentFastBlock().Number()
				lastBlock = d.blockchain.CurrentBlock().Number()
			}
			if err := d.lightchain.SetHead(rollback - 1); err != nil { // -1 to target the parent of the first uncertain block
				// We're already unwinding the stack, only print the error to make it more visible
				log.Error("Failed to roll back chain segment", "head", rollback-1, "err", err)
			}
			curFastBlock, curBlock := common.Big0, common.Big0
			if mode != LightSync {
				curFastBlock = d.blockchain.CurrentFastBlock().Number()
				curBlock = d.blockchain.CurrentBlock().Number()
			}
			log.Warn("Rolled back chain segment",
				"header", fmt.Sprintf("%d->%d", lastHeader, d.lightchain.CurrentHeader().Number),
				"snap", fmt.Sprintf("%d->%d", lastFastBlock, curFastBlock),
				"block", fmt.Sprintf("%d->%d", lastBlock, curBlock), "reason", rollbackErr)
		}
	}()
	// Wait for batches of headers to process
	gotHeaders := false

	for {
		select {
		case <-d.cancelCh:
			rollbackErr = errCanceled
			return errCanceled

		case task := <-d.headerProcCh:
			// Terminate header processing if we synced up
			if task == nil || len(task.headers) == 0 {
				// Notify everyone that headers are fully processed
				for _, ch := range []chan bool{d.queue.blockWakeCh, d.queue.receiptWakeCh} {
					select {
					case ch <- false:
					case <-d.cancelCh:
					}
				}
				// If no headers were retrieved at all, the peer violated its TD promise that it had a
				// better chain compared to ours. The only exception is if its promised blocks were
				// already imported by other means (e.g. fetcher):
				//
				// R <remote peer>, L <local node>: Both at block 10
				// R: Mine block 11, and propagate it to L
				// L: Queue block 11 for import
				// L: Notice that R's head and TD increased compared to ours, start sync
				// L: Import of block 11 finishes
				// L: Sync begins, and finds common ancestor at 11
				// L: Request new headers up from 11 (R's TD was higher, it must have something)
				// R: Nothing to give
				if mode != LightSync {
					head := d.blockchain.CurrentBlock()
					if !gotHeaders && bn.Cmp(head.Number()) > 0 {
						return errStallingPeer
					}
				}
				// If snap or light syncing, ensure promised headers are indeed delivered. This is
				// needed to detect scenarios where an attacker feeds a bad pivot and then bails out
				// of delivering the post-pivot blocks that would flag the invalid content.
				//
				// This check cannot be executed "as is" for full imports, since blocks may still be
				// queued for processing when the header download completes. However, as long as the
				// peer gave us something useful, we're already happy/progressed (above check).
				if mode == SnapSync || mode == LightSync {
					head := d.lightchain.CurrentHeader()
					if bn.Cmp(head.Number) > 0 {
						return errStallingPeer
					}
				}
				// Disable any rollback and return
				rollback = 0
				return nil
			}
			// Otherwise split the chunk of headers into batches and process them
			headers, hashes := task.headers, task.hashes

			gotHeaders = true
			for len(headers) > 0 {
				// Terminate if something failed in between processing chunks
				select {
				case <-d.cancelCh:
					rollbackErr = errCanceled
					return errCanceled
				default:
				}
				// Select the next chunk of headers to import
				limit := maxHeadersProcess
				if limit > len(headers) {
					limit = len(headers)
				}
				chunkHeaders := headers[:limit]
				chunkHashes := hashes[:limit]

				// In case of header only syncing, validate the chunk immediately
				if mode == SnapSync || mode == LightSync {
					// If we're importing pure headers, verify based on their recentness
					var pivot uint64

					d.pivotLock.RLock()
					if d.pivotHeader != nil {
						pivot = d.pivotHeader.Number.Uint64()
					}
					d.pivotLock.RUnlock()

					frequency := fsHeaderCheckFrequency
					if chunkHeaders[len(chunkHeaders)-1].Number.Uint64()+uint64(fsHeaderForceVerify) > pivot {
						frequency = 1
					}
					if n, err := d.lightchain.InsertHeaderChain(chunkHeaders, frequency); err != nil {
						rollbackErr = err

						// If some headers were inserted, track them as uncertain
						if (mode == SnapSync || frequency > 1) && n > 0 && rollback == 0 {
							rollback = chunkHeaders[0].Number.Uint64()
						}
						log.Warn("Invalid header encountered", "number", chunkHeaders[n].Number, "hash", chunkHashes[n], "parent", chunkHeaders[n].ParentHash, "err", err)
						return fmt.Errorf("%w: %v", errInvalidChain, err)
					}
					// All verifications passed, track all headers within the alloted limits
					if mode == SnapSync {
						head := chunkHeaders[len(chunkHeaders)-1].Number.Uint64()
						if head-rollback > uint64(fsHeaderSafetyNet) {
							rollback = head - uint64(fsHeaderSafetyNet)
						} else {
							rollback = 1
						}
					}
				}
				// Unless we're doing light chains, schedule the headers for associated content retrieval
				if mode == FullSync || mode == SnapSync {
					// If we've reached the allowed number of pending headers, stall a bit
					for d.queue.PendingBodies() >= maxQueuedHeaders || d.queue.PendingReceipts() >= maxQueuedHeaders {
						select {
						case <-d.cancelCh:
							rollbackErr = errCanceled
							return errCanceled
						case <-time.After(time.Second):
						}
					}
					// Otherwise insert the headers for content retrieval
					inserts := d.queue.Schedule(chunkHeaders, chunkHashes, origin)
					if len(inserts) != len(chunkHeaders) {
						rollbackErr = fmt.Errorf("stale headers: len inserts %v len(chunk) %v", len(inserts), len(chunkHeaders))
						return fmt.Errorf("%w: stale headers", errBadPeer)
					}
				}
				headers = headers[limit:]
				hashes = hashes[limit:]
				origin += uint64(limit)
			}
			// Update the highest block number we know if a higher one is found.
			d.syncStatsLock.Lock()
			if d.syncStatsChainHeight < origin {
				d.syncStatsChainHeight = origin - 1
			}
			d.syncStatsLock.Unlock()

			// Signal the content downloaders of the availablility of new tasks
			for _, ch := range []chan bool{d.queue.blockWakeCh, d.queue.receiptWakeCh} {
				select {
				case ch <- true:
				default:
				}
			}
		}
	}
}

// processFullSyncContent takes fetch results from the queue and imports them into the chain.
func (d *Downloader) processFullSyncContent() error {
	for {
		results := d.queue.Results(true)
		if len(results) == 0 {
			return nil
		}
		if d.chainInsertHook != nil {
			d.chainInsertHook(results)
		}
		if err := d.importBlockResults(results); err != nil {
			return err
		}
	}
}

func (d *Downloader) importBlockResults(results []*fetchResult) error {
	// Check for any early termination requests
	if len(results) == 0 {
		return nil
	}
	select {
	case <-d.quitCh:
		return errCancelContentProcessing
	default:
	}
	// Retrieve the a batch of results to import
	first, last := results[0].Header, results[len(results)-1].Header
	log.Debug("Inserting downloaded chain", "items", len(results),
		"firstnum", first.Number, "firsthash", first.Hash(),
		"lastnum", last.Number, "lasthash", last.Hash(),
	)
	blocks := make([]*types.Block, len(results))
	for i, result := range results {
		blocks[i] = types.NewBlockWithHeader(result.Header).WithBody(result.Transactions, result.ExtraData)
	}
	// Downloaded blocks are always regarded as trusted after the
	// transition. Because the downloaded chain is guided by the
	// consensus-layer.
	if index, err := d.blockchain.InsertChain(blocks); err != nil {
		if index < len(results) {
			log.Debug("Downloaded item processing failed", "number", results[index].Header.Number, "hash", results[index].Header.Hash(), "err", err)
		} else {
			// The InsertChain method in blockchain.go will sometimes return an out-of-bounds index,
			// when it needs to preprocess blocks to import a sidechain.
			// The importer will put together a new list of blocks to import, which is a superset
			// of the blocks delivered from the downloader, and the indexing will be off.
			log.Debug("Downloaded item processing failed on sidechain import", "index", index, "err", err)
		}
		return fmt.Errorf("%w: %v", errInvalidChain, err)
	}
	return nil
}

// processSnapSyncContent takes fetch results from the queue and writes them to the
// database. It also controls the synchronisation of state nodes of the pivot block.
func (d *Downloader) processSnapSyncContent() error {
	// Start syncing state of the reported head block. This should get us most of
	// the state of the pivot block.
	d.pivotLock.RLock()
	sync := d.syncState(d.pivotHeader.Root)
	d.pivotLock.RUnlock()

	defer func() {
		// The `sync` object is replaced every time the pivot moves. We need to
		// defer close the very last active one, hence the lazy evaluation vs.
		// calling defer sync.Cancel() !!!
		sync.Cancel()
	}()

	closeOnErr := func(s *stateSync) {
		if err := s.Wait(); err != nil && err != errCancelStateFetch && err != errCanceled && err != snap.ErrCancelled {
			d.queue.Close() // wake up Results
		}
	}
	go closeOnErr(sync)

	// To cater for moving pivot points, track the pivot block and subsequently
	// accumulated download results separately.
	var (
		oldPivot *fetchResult   // Locked in pivot block, might change eventually
		oldTail  []*fetchResult // Downloaded content after the pivot
	)
	for {
		// Wait for the next batch of downloaded data to be available, and if the pivot
		// block became stale, move the goalpost
		results := d.queue.Results(oldPivot == nil) // Block if we're not monitoring pivot staleness
		if len(results) == 0 {
			// If pivot sync is done, stop
			if oldPivot == nil {
				return sync.Cancel()
			}
			// If sync failed, stop
			select {
			case <-d.cancelCh:
				sync.Cancel()
				return errCanceled
			default:
			}
		}
		if d.chainInsertHook != nil {
			d.chainInsertHook(results)
		}
		// If we haven't downloaded the pivot block yet, check pivot staleness
		// notifications from the header downloader
		d.pivotLock.RLock()
		pivot := d.pivotHeader
		d.pivotLock.RUnlock()

		if oldPivot == nil {
			if pivot.Root != sync.root {
				sync.Cancel()
				sync = d.syncState(pivot.Root)

				go closeOnErr(sync)
			}
		} else {
			results = append(append([]*fetchResult{oldPivot}, oldTail...), results...)
		}
		// Split around the pivot block and process the two sides via snap/full sync
		if atomic.LoadInt32(&d.committed) == 0 {
			latest := results[len(results)-1].Header
			// If the height is above the pivot block by 2 sets, it means the pivot
			// become stale in the network and it was garbage collected, move to a
			// new pivot.
			//
			// Note, we have `reorgProtHeaderDelay` number of blocks withheld, Those
			// need to be taken into account, otherwise we're detecting the pivot move
			// late and will drop peers due to unavailable state!!!
			if height := latest.Number.Uint64(); height >= pivot.Number.Uint64()+2*uint64(fsMinFullBlocks)-uint64(reorgProtHeaderDelay) {
				log.Warn("Pivot became stale, moving", "old", pivot.Number.Uint64(), "new", height-uint64(fsMinFullBlocks)+uint64(reorgProtHeaderDelay))
				pivot = results[len(results)-1-fsMinFullBlocks+reorgProtHeaderDelay].Header // must exist as lower old pivot is uncommitted

				d.pivotLock.Lock()
				d.pivotHeader = pivot
				d.pivotLock.Unlock()

				// Write out the pivot into the database so a rollback beyond it will
				// reenable snap sync
				rawdb.WriteLastPivotNumber(d.stateDB, pivot.Number.Uint64())
			}
		}
		P, beforeP, afterP := splitAroundPivot(pivot.Number.Uint64(), results)
		if err := d.commitSnapSyncData(beforeP, sync); err != nil {
			return err
		}
		if P != nil {
			// If new pivot block found, cancel old state retrieval and restart
			if oldPivot != P {
				sync.Cancel()
				sync = d.syncState(P.Header.Root)

				go closeOnErr(sync)
				oldPivot = P
			}
			// Wait for completion, occasionally checking for pivot staleness
			select {
			case <-sync.done:
				if sync.err != nil {
					return sync.err
				}
				// sync PPOS
				//todo when eth67,将pos消息全部启用id并且同步改造
				peers := d.peers.AllPeers()
				if len(peers) == 0 {
					return errors.New("no idle peers to fetch PPOSStorageV2")
				}
				fetchPPOSStorageSuccess := false
				for _, peer := range peers {
					_, remote := peer.peer.Head()
					if remote.Cmp(P.Header.Number) > 0 {
						if err := d.fetchPPOSStorageV2(peer, P.Header); err != nil {
							peer.log.Error("failed to fetch PPOSStorageV2", "err", err)
							if errors.Is(err, errCanceled) {
								return err
							}
						} else {
							fetchPPOSStorageSuccess = true
							break
						}
					}
				}
				if !fetchPPOSStorageSuccess {
					return errors.New("failed to fetch PPOSStorageV2")
				}

				if err := d.commitPivotBlock(P); err != nil {
					return err
				}
				oldPivot = nil

			case <-time.After(time.Second):
				oldTail = afterP
				continue
			}
		}
		// Fast sync done, pivot commit done, full import
		if err := d.importBlockResults(afterP); err != nil {
			return err
		}
	}
}

func splitAroundPivot(pivot uint64, results []*fetchResult) (p *fetchResult, before, after []*fetchResult) {
	if len(results) == 0 {
		return nil, nil, nil
	}
	if lastNum := results[len(results)-1].Header.Number.Uint64(); lastNum < pivot {
		// the pivot is somewhere in the future
		return nil, results, nil
	}
	// This can also be optimized, but only happens very seldom
	for _, result := range results {
		num := result.Header.Number.Uint64()
		switch {
		case num < pivot:
			before = append(before, result)
		case num == pivot:
			p = result
		default:
			after = append(after, result)
		}
	}
	return p, before, after
}

func (d *Downloader) commitSnapSyncData(results []*fetchResult, stateSync *stateSync) error {
	// Check for any early termination requests
	if len(results) == 0 {
		return nil
	}
	select {
	case <-d.quitCh:
		return errCancelContentProcessing
	case <-stateSync.done:
		if err := stateSync.Wait(); err != nil {
			return err
		}
	default:
	}
	// Retrieve the a batch of results to import
	first, last := results[0].Header, results[len(results)-1].Header
	log.Debug("Inserting snap-sync blocks", "items", len(results),
		"firstnum", first.Number, "firsthash", first.Hash(),
		"lastnumn", last.Number, "lasthash", last.Hash(),
	)
	blocks := make([]*types.Block, len(results))
	receipts := make([]types.Receipts, len(results))
	for i, result := range results {
		blocks[i] = types.NewBlockWithHeader(result.Header).WithBody(result.Transactions, result.ExtraData)
		receipts[i] = result.Receipts
	}
	if index, err := d.blockchain.InsertReceiptChain(blocks, receipts, d.ancientLimit); err != nil {
		log.Debug("Downloaded item processing failed", "number", results[index].Header.Number, "hash", results[index].Header.Hash(), "err", err)
		return fmt.Errorf("%w: %v", errInvalidChain, err)
	}
	return nil
}

func (d *Downloader) commitPivotBlock(result *fetchResult) error {
	block := types.NewBlockWithHeader(result.Header).WithBody(result.Transactions, result.ExtraData)
	log.Debug("Committing snap sync pivot as new head", "number", block.Number(), "hash", block.Hash())

	// Commit the pivot block as the new head, will require full sync from here on
	if _, err := d.blockchain.InsertReceiptChain([]*types.Block{block}, []types.Receipts{result.Receipts}, d.ancientLimit); err != nil {
		return err
	}
	if err := d.blockchain.SnapSyncCommitHead(block.Hash()); err != nil {
		return err
	}
	atomic.StoreInt32(&d.committed, 1)
	return nil
}

// DeliverPposStorage injects a new batch of ppos storage received from a remote node.
func (d *Downloader) DeliverPposStorage(id string, kvs [][2][]byte, last bool, kvNum uint64, base bool, blocks []snapshotdb.BlockData, baseBlock uint64) (err error) {
	return d.deliver(d.pposStorageCh, &pposStoragePack{id, kvs, last, kvNum, base, blocks, baseBlock}, pposStorageInMeter, pposStorageDropMeter)
}

// DeliverPposStorage injects a new batch of ppos storage received from a remote node.
func (d *Downloader) DeliverPposInfo(id string, latest, pivot *types.Header) (err error) {
	return d.deliver(d.pposInfoCh, &pposInfoPack{id, latest, pivot}, pposStorageInMeter, pposStorageDropMeter)
}

func (d *Downloader) DeliverOriginAndPivot(id string, headers []*types.Header) (err error) {
	return d.deliver(d.originAndPivotCh, &headerPack{id, headers}, headerInMeter, headerDropMeter)
}

// DeliverSnapPacket is invoked from a peer's message handler when it transmits a
// data packet for the local node to consume.
func (d *Downloader) DeliverSnapPacket(peer *snap.Peer, packet snap.Packet) error {
	switch packet := packet.(type) {
	case *snap.AccountRangePacket:
		hashes, accounts, err := packet.Unpack()
		if err != nil {
			return err
		}
		return d.SnapSyncer.OnAccounts(peer, packet.ID, hashes, accounts, packet.Proof)

	case *snap.StorageRangesPacket:
		hashset, slotset := packet.Unpack()
		return d.SnapSyncer.OnStorage(peer, packet.ID, hashset, slotset, packet.Proof)

	case *snap.ByteCodesPacket:
		return d.SnapSyncer.OnByteCodes(peer, packet.ID, packet.Codes)

	case *snap.TrieNodesPacket:
		return d.SnapSyncer.OnTrieNodes(peer, packet.ID, packet.Nodes)

	default:
		return fmt.Errorf("unexpected snap packet type: %T", packet)
	}
}

// deliver injects a new batch of data received from a remote node.
func (d *Downloader) deliver(destCh chan dataPack, packet dataPack, inMeter, dropMeter metrics.Meter) (err error) {
	// Update the delivery metrics for both good and failed deliveries
	inMeter.Mark(int64(packet.Items()))
	defer func() {
		if err != nil {
			dropMeter.Mark(int64(packet.Items()))
		}
	}()
	// Deliver or abort if the sync is canceled while queuing
	d.cancelLock.RLock()
	cancel := d.cancelCh
	d.cancelLock.RUnlock()
	if cancel == nil {
		return errNoSyncActive
	}
	select {
	case destCh <- packet:
		return nil
	case <-cancel:
		return errNoSyncActive
	}
}
