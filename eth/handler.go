// Copyright 2015 The PlatON-Go Authors
// This file is part of the PlatON-Go library.
//
// The go-PlatON library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package eth

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/eth/fetcher"
	"github.com/PlatONnetwork/PlatON-Go/eth/protocols/eth"
	"github.com/PlatONnetwork/PlatON-Go/eth/protocols/snap"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

const (
	// txChanSize is the size of channel listening to NewTxsEvent.
	// The number is referenced from the size of tx pool.
	txChanSize = 4096

	numBroadcastTxPeers     = 5 // Maximum number of peers for broadcast transactions
	numBroadcastTxHashPeers = 5 // Maximum number of peers for broadcast transactions hash
	numBroadcastBlockPeers  = 5 // Maximum number of peers for broadcast new block

	defaultTxsCacheSize      = 20
	defaultBroadcastInterval = 100 * time.Millisecond
)

var (
	syncChallengeTimeout = 15 * time.Second // Time allowance for a node to reply to the sync progress challenge
)

// txPool defines the methods needed from a transaction pool implementation to
// support all the operations needed by the Ethereum chain protocols.
type txPool interface {
	// Has returns an indicator whether txpool has a transaction
	// cached with the given hash.
	Has(hash common.Hash) bool

	// Get retrieves the transaction from local txpool with given
	// tx hash.
	Get(hash common.Hash) *types.Transaction

	// AddRemotes should add the given transactions to the pool.
	AddRemotes([]*types.Transaction) []error

	// Pending should return pending transactions.
	// The slice should be modifiable by the caller.
	Pending(enforceTips, limited bool) map[common.Address]types.Transactions

	// SubscribeNewTxsEvent should return an event subscription of
	// NewTxsEvent and send events to the given channel.
	SubscribeNewTxsEvent(chan<- core.NewTxsEvent) event.Subscription
}

// handlerConfig is the collection of initialization parameters to create a full
// node network handler.
type handlerConfig struct {
	Database   ethdb.Database            // Database for direct sync insertions
	Chain      *core.BlockChain          // Blockchain to serve data from
	TxPool     txPool                    // Transaction pool to propagate from
	Network    uint64                    // Network identifier to adfvertise
	Sync       downloader.SyncMode       // Whether to fast or full sync
	BloomCache uint64                    // Megabytes to alloc for fast sync bloom
	EventMux   *event.TypeMux            // Legacy event mux, deprecate for `feed`
	Checkpoint *params.TrustedCheckpoint // Hard coded checkpoint for sync challenges
	Whitelist  map[uint64]common.Hash    // Hard coded whitelist for sync challenged
}

type handler struct {
	networkID uint64

	fastSync        uint32 // Flag whether fast sync is enabled (gets disabled if we already have blocks)
	snapSync        uint32 // Flag whether fast sync should operate on top of the snap protocol
	acceptTxs       uint32 // Flag whether we're considered synchronised (enables transaction processing)
	acceptRemoteTxs uint32 // Flag whether we're accept remote txs

	chainconfig *params.ChainConfig
	database    ethdb.Database
	txpool      txPool
	chain       *core.BlockChain
	maxPeers    int

	downloader   *downloader.Downloader
	blockFetcher *fetcher.BlockFetcher
	txFetcher    *fetcher.TxFetcher
	peers        *peerSet

	eventMux             *event.TypeMux
	txsCh                chan core.NewTxsEvent
	txsCache             []*types.Transaction
	txsSub               event.Subscription
	minedBlockSub        *event.TypeMuxSubscription
	prepareMinedBlockSub *event.TypeMuxSubscription
	blockSignatureSub    *event.TypeMuxSubscription

	whitelist map[uint64]common.Hash

	// channels for fetcher, syncer, txsyncLoop
	txsyncCh chan *txsync
	quitSync chan struct{}

	chainSync *chainSyncer
	wg        sync.WaitGroup

	handlerStartCh chan struct{}
	handlerDoneCh  chan struct{}

	engine consensus.Engine
}

// newHandler returns a handler for all Ethereum chain management protocol.
func newHandler(config *handlerConfig) (*handler, error) {
	// Create the protocol manager with the base fields
	if config.EventMux == nil {
		config.EventMux = new(event.TypeMux) // Nicety initialization for tests
	}
	h := &handler{
		networkID:      config.Network,
		eventMux:       config.EventMux,
		database:       config.Database,
		txpool:         config.TxPool,
		chain:          config.Chain,
		peers:          newPeerSet(),
		txsyncCh:       make(chan *txsync),
		quitSync:       make(chan struct{}),
		handlerDoneCh:  make(chan struct{}),
		handlerStartCh: make(chan struct{}),
		engine:         config.Chain.Engine(),
	}
	if config.Sync == downloader.FullSync {
		// The database seems empty as the current block is the genesis. Yet the fast
		// block is ahead, so fast sync was enabled for this node at a certain point.
		// The scenarios where this can happen is
		// * if the user manually (or via a bad block) rolled back a fast sync node
		//   below the sync point.
		// * the last fast sync is not finished while user specifies a full sync this
		//   time. But we don't have any recent state for full sync.
		// In these cases however it's safe to reenable fast sync.
		fullBlock, fastBlock := h.chain.CurrentBlock(), h.chain.CurrentFastBlock()
		if fullBlock.NumberU64() == 0 && fastBlock.NumberU64() > 0 {
			h.fastSync = uint32(1)
			log.Warn("Switch sync mode from full sync to fast sync")
		}
	} else {
		if h.chain.CurrentBlock().NumberU64() > 0 {
			// Print warning log if database is not empty to run fast sync.
			log.Warn("Switch sync mode from fast sync to full sync")
		} else {
			// If fast sync was requested and our database is empty, grant it
			h.fastSync = uint32(1)
			if config.Sync == downloader.SnapSync {
				h.snapSync = uint32(1)
			}
		}
	}

	decodeExtra := func(extra []byte) (common.Hash, uint64, error) {
		return h.engine.DecodeExtra(extra)
	}

	// Construct the downloader (long sync) and its backing state bloom if fast
	// sync is requested. The downloader is responsible for deallocating the state
	// bloom when it's done.
	h.downloader = downloader.New(config.Database, snapshotdb.Instance(), h.eventMux, h.chain, nil, h.removePeer, decodeExtra)

	// Construct the fetcher (short sync)
	validator := func(header *types.Header) error {
		return h.chain.Engine().VerifyHeader(h.chain, header, true)
	}
	heighter := func() uint64 {
		return h.chain.Engine().CurrentBlock().NumberU64() + 1
	}
	inserter := func(blocks types.Blocks) (int, error) {
		// If fast sync is running, deny importing weird blocks. This is a problematic
		// clause when starting up a new network, because fast-syncing miners might not
		// accept each others' blocks until a restart. Unfortunately we haven't figured
		// out a way yet where nodes can decide unilaterally whether the network is new
		// or not. This should be fixed if we figure out a solution.
		if atomic.LoadUint32(&h.fastSync) == 1 {
			log.Warn("Fast syncing, discarded propagated block", "number", blocks[0].Number(), "hash", blocks[0].Hash())
			return 0, nil
		}
		n, err := h.chain.InsertChain(blocks)
		if err == nil {
			atomic.StoreUint32(&h.acceptTxs, 1) // Mark initial sync done on any fetcher import
		}
		return n, err
	}
	getBlockByHash := func(hash common.Hash) *types.Block {
		return h.chain.GetBlockByHash(hash)
	}
	h.blockFetcher = fetcher.NewBlockFetcher(getBlockByHash, validator, h.BroadcastBlock, heighter, inserter, h.removePeer, decodeExtra)

	fetchTx := func(peer string, hashes []common.Hash) error {
		p := h.peers.peer(peer)
		if p == nil {
			return errors.New("unknown peer")
		}
		return p.RequestTxs(hashes)
	}
	h.txFetcher = fetcher.NewTxFetcher(h.txpool.Has, h.txpool.AddRemotes, fetchTx)
	h.chainSync = newChainSyncer(h)
	return h, nil
}

// protoTracker tracks the number of active protocol handlers.
func (h *handler) protoTracker() {
	defer h.wg.Done()
	var active int
	for {
		select {
		case <-h.handlerStartCh:
			active++
		case <-h.handlerDoneCh:
			active--
		case <-h.quitSync:
			// Wait for all active handlers to finish.
			for ; active > 0; active-- {
				<-h.handlerDoneCh
			}
			return
		}
	}
}

// incHandlers signals to increment the number of active handlers if not
// quitting.
func (h *handler) incHandlers() bool {
	select {
	case h.handlerStartCh <- struct{}{}:
		return true
	case <-h.quitSync:
		return false
	}
}

// decHandlers signals to decrement the number of active handlers.
func (h *handler) decHandlers() {
	h.handlerDoneCh <- struct{}{}
}

// runEthPeer registers an eth peer into the joint eth/snap peerset, adds it to
// various subsistems and starts handling messages.
func (h *handler) runEthPeer(peer *eth.Peer, handler eth.Handler) error {
	if !h.incHandlers() {
		return p2p.DiscQuitting
	}
	defer h.decHandlers()

	// If the peer has a `snap` extension, wait for it to connect so we can have
	// a uniform initialization/teardown mechanism
	snap, err := h.peers.waitSnapExtension(peer)
	if err != nil {
		peer.Log().Error("Snapshot extension barrier failed", "err", err)
		return err
	}

	// Execute the Ethereum handshake
	var (
		genesis = h.chain.Genesis()
		head    = h.chain.CurrentHeader()
		hash    = head.Hash()
	)
	verify := func(remoteNum *big.Int, remoteHash common.Hash) error {
		// A simple hash consistency check,but does not prevent malicious node connections
		if head.Number == remoteNum && hash != remoteHash {
			return fmt.Errorf("block %v, hash %s != remote %s", remoteNum, hash, remoteHash)
		} else if head.Number.Uint64() > remoteNum.Uint64() {
			lowHeader := h.chain.GetHeaderByNumber(remoteNum.Uint64())
			if lowHeader.Hash() != remoteHash {
				return fmt.Errorf("block %v, hash %s != remote %s", remoteNum, lowHeader.Hash(), remoteHash)
			}
		}
		return nil
	}
	if err := peer.Handshake(h.networkID, head.Number, hash, genesis.Hash(), verify); err != nil {
		peer.Log().Debug("PlatON handshake failed", "err", err)
		return err
	}

	reject := false // reserved peer slots
	if atomic.LoadUint32(&h.snapSync) == 1 {
		if snap == nil {
			// If we are running snap-sync, we want to reserve roughly half the peer
			// slots for peers supporting the snap protocol.
			// The logic here is; we only allow up to 5 more non-snap peers than snap-peers.
			if all, snp := h.peers.len(), h.peers.snapLen(); all-snp > snp+5 {
				reject = true
			}
		}
	}
	// Ignore maxPeers if this is a trusted peer
	if !peer.Peer.Info().Network.Trusted {
		if reject || h.peers.len() >= h.maxPeers {
			return p2p.DiscTooManyPeers
		}
	}
	peer.Log().Debug("PlatON peer connected", "name", peer.Name())

	// Register the peer locally
	if err := h.peers.registerPeer(peer, snap); err != nil {
		peer.Log().Error("PlatON peer registration failed", "err", err)
		return err
	}
	defer h.unregisterPeer(peer.ID())

	p := h.peers.peer(peer.ID())
	if p == nil {
		return errors.New("peer dropped during handling")
	}
	// Register the peer in the downloader. If the downloader considers it banned, we disconnect
	if err := h.downloader.RegisterPeer(peer.ID(), peer.Version(), peer); err != nil {
		peer.Log().Error("Failed to register peer in eth syncer", "err", err)
		return err
	}
	if snap != nil {
		if err := h.downloader.SnapSyncer.Register(snap); err != nil {
			peer.Log().Error("Failed to register peer in snap syncer", "err", err)
			return err
		}
	}
	h.chainSync.handlePeerEvent()

	// Propagate existing transactions. new transactions appearing
	// after this will be sent via broadcasts.
	h.syncTransactions(peer)

	// If we have any explicit whitelist block hashes, request them
	for number := range h.whitelist {
		if err := peer.RequestHeadersByNumber(number, 1, 0, false); err != nil {
			return err
		}
	}
	// Handle incoming messages until the connection is torn down
	return handler(peer)
}

// runSnapExtension registers a `snap` peer into the joint eth/snap peerset and
// starts handling inbound messages. As `snap` is only a satellite protocol to
// `eth`, all subsystem registrations and lifecycle management will be done by
// the main `eth` handler to prevent strange races.
func (h *handler) runSnapExtension(peer *snap.Peer, handler snap.Handler) error {
	if !h.incHandlers() {
		return p2p.DiscQuitting
	}
	defer h.decHandlers()

	if err := h.peers.registerSnapExtension(peer); err != nil {
		peer.Log().Error("Snapshot extension registration failed", "err", err)
		return err
	}
	return handler(peer)
}

// removePeer requests disconnection of a peer.
func (h *handler) removePeer(id string) {
	peer := h.peers.peer(id)
	if peer != nil {
		peer.Peer.Disconnect(p2p.DiscUselessPeer)
	}
}

// unregisterPeer removes a peer from the downloader, fetchers and main peer set.
func (h *handler) unregisterPeer(id string) {
	// Create a custom logger to avoid printing the entire id
	var logger log.Logger
	if len(id) < 16 {
		// Tests use short IDs, don't choke on them
		logger = log.New("peer", id)
	} else {
		logger = log.New("peer", id[:8])
	}
	// Abort if the peer does not exist
	peer := h.peers.peer(id)
	if peer == nil {
		logger.Error("PlatON peer removal failed", "err", errPeerNotRegistered)
		return
	}
	// Remove the `eth` peer if it exists
	logger.Debug("Removing PlatON peer", "snap", peer.snapExt != nil)

	// Remove the `snap` extension if it exists
	if peer.snapExt != nil {
		h.downloader.SnapSyncer.Unregister(id)
	}
	h.downloader.UnregisterPeer(id)
	h.txFetcher.Drop(id)

	if err := h.peers.unregisterPeer(id); err != nil {
		logger.Error("PlatON peer removal failed", "err", err)
	}
}

func (h *handler) Start(maxPeers int) {
	h.maxPeers = maxPeers

	// broadcast transactions
	h.wg.Add(1)
	h.txsCh = make(chan core.NewTxsEvent, txChanSize)
	h.txsSub = h.txpool.SubscribeNewTxsEvent(h.txsCh)
	go h.txBroadcastLoop()

	// broadcast mined blocks
	h.wg.Add(1)
	h.minedBlockSub = h.eventMux.Subscribe(core.NewMinedBlockEvent{})
	go h.minedBroadcastLoop()

	// start sync handlers
	h.wg.Add(2)
	go h.chainSync.loop()
	go h.txsyncLoop64() // TODO(karalabe): Legacy initial tx echange, drop with eth/64.

	// start peer handler tracker
	h.wg.Add(1)
	go h.protoTracker()
}

func (h *handler) Stop() {
	h.txsSub.Unsubscribe()        // quits txBroadcastLoop
	h.minedBlockSub.Unsubscribe() // quits blockBroadcastLoop

	// Quit chainSync and txsync64.
	// After this is done, no new peers will be accepted.
	close(h.quitSync)

	// Disconnect existing sessions.
	// This also closes the gate for any new registrations on the peer set.
	// sessions which are already established but not added to h.peers yet
	// will exit when they try to register.
	h.peers.close()
	h.wg.Wait()

	log.Info("PlatON protocol stopped")
}

// BroadcastBlock will either propagate a block to a subset of its peers, or
// will only announce its availability (depending what's requested).
func (h *handler) BroadcastBlock(block *types.Block, propagate bool) {
	hash := block.Hash()
	peers := h.peers.peersWithoutBlock(hash)

	// If propagation is requested, send to a subset of the peer
	if propagate {
		if parent := h.chain.GetBlock(block.ParentHash(), block.NumberU64()-1); parent != nil {
		} else {
			log.Warn("Propagating dangling block", "number", block.Number(), "hash", hash)
			return
		}

		var transfer []*ethPeer
		if len(peers) <= numBroadcastBlockPeers {
			// Send the block to all peers
			transfer = peers
		} else {
			// Send the block to a subset of our peers
			rd := rand.New(rand.NewSource(time.Now().UnixNano()))
			indexes := rd.Perm(len(peers))
			maxPeers := int(math.Sqrt(float64(len(peers))))
			transfer = make([]*ethPeer, 0, maxPeers)
			for i := 0; i < maxPeers; i++ {
				transfer = append(transfer, peers[indexes[i]])
			}
		}
		for _, peer := range transfer {
			peer.AsyncSendNewBlock(block)
		}
		log.Trace("Propagated block", "hash", hash, "recipients", len(transfer), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
		return
	}
	// Otherwise if the block is indeed in out own chain, announce it
	if h.chain.HasBlock(hash, block.NumberU64()) {
		for _, peer := range peers {
			peer.AsyncSendNewBlockHash(block)
		}
		log.Trace("Announced block", "hash", hash, "recipients", len(peers), "duration", common.PrettyDuration(time.Since(block.ReceivedAt)))
	}
}

// BroadcastTransactions will propagate a batch of transactions
// - To a square root of all peers
// - And, separately, as announcements to all peers which are not known to
// already have the given transaction.
func (h *handler) BroadcastTransactions(txs types.Transactions) {
	var (
		annoCount   int // Count of announcements made
		annoPeers   int
		directCount int // Count of the txs sent directly to peers
		directPeers int // Count of the peers that were sent transactions directly

		txset = make(map[*ethPeer]types.Transactions) // Set peer->transaction to transfer directly
		annos = make(map[*ethPeer][]common.Hash)      // Set peer->hash to announce

	)
	rd := rand.New(rand.NewSource(time.Now().UnixNano()))

	// Broadcast transactions to a batch of peers not knowing about it
	for _, tx := range txs {
		peers := h.peers.peersWithoutTransaction(tx.Hash())
		if len(peers) <= numBroadcastTxPeers {
			for _, peer := range peers {
				txset[peer] = append(txset[peer], tx)
			}
		} else {
			indexes := rd.Perm(len(peers))
			numAnnos := len(peers) - numBroadcastTxPeers
			countAnnos := 0
			if numAnnos > numBroadcastTxHashPeers {
				numAnnos = numBroadcastTxHashPeers
			}
			for i, c := 0, 0; i < len(peers) && countAnnos < numAnnos; i, c = i+1, c+1 {
				peer := peers[indexes[i]]
				if c < numBroadcastTxPeers {
					txset[peer] = append(txset[peer], tx)
				} else {
					// For the remaining peers, send announcement only
					annos[peer] = append(annos[peer], tx.Hash())
					countAnnos++
				}
			}
		}
	}

	for peer, txs := range txset {
		directPeers++
		directCount += len(txs)
		peer.AsyncSendTransactions(txs)
	}
	for peer, hashes := range annos {
		annoPeers++
		annoCount += len(hashes)
		peer.AsyncSendPooledTransactionHashes(hashes)
	}
	log.Debug("Transaction broadcast", "txs", len(txs),
		"announce packs", annoPeers, "announced hashes", annoCount,
		"tx packs", directPeers, "broadcast txs", directCount)
}

// minedBroadcastLoop sends mined blocks to connected peers.
func (h *handler) minedBroadcastLoop() {
	defer h.wg.Done()

	for obj := range h.minedBlockSub.Chan() {
		if ev, ok := obj.Data.(core.NewMinedBlockEvent); ok {
			h.BroadcastBlock(ev.Block, true)  // First propagate block to peers
			h.BroadcastBlock(ev.Block, false) // Only then announce to the rest
		}
	}
}

// txBroadcastLoop announces new transactions to connected peers.
func (h *handler) txBroadcastLoop() {
	defer h.wg.Done()

	timer := time.NewTimer(defaultBroadcastInterval)

	for {
		select {
		case event := <-h.txsCh:
			h.txsCache = append(h.txsCache, event.Txs...)
			if len(h.txsCache) >= defaultTxsCacheSize {
				//log.Trace("broadcast txs", "count", len(pm.txsCache))
				h.BroadcastTransactions(h.txsCache)
				h.txsCache = make([]*types.Transaction, 0)
				timer.Reset(defaultBroadcastInterval)
			}
		case <-timer.C:
			if len(h.txsCache) > 0 {
				//log.Trace("broadcast txs", "count", len(pm.txsCache))
				h.BroadcastTransactions(h.txsCache)
				h.txsCache = make([]*types.Transaction, 0)
			}
			timer.Reset(defaultBroadcastInterval)

			// Err() channel will be closed when unsubscribing.
		case <-h.txsSub.Err():
			return
		}
	}
}
