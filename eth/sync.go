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

package eth

import (
	"math/big"
	"math/rand"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/eth/protocols/eth"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

const (
	forceSyncCycle      = 5 * time.Second // Time interval to force syncs, even if few peers are available
	defaultMinSyncPeers = 5               // Amount of peers desired to start syncing

	// This is the target size for the packs of transactions sent by txsyncLoop.
	// A pack can get larger than this if a single transactions exceeds this size.
	txsyncPackSize = 100 * 1024
)

type txsync struct {
	p   *eth.Peer
	txs []*types.Transaction
}

// syncTransactions starts sending all currently pending transactions to the given peer.
func (h *handler) syncTransactions(p *eth.Peer) {
	var txs types.Transactions
	pending := h.txpool.Pending(false, false)
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	if len(txs) == 0 {
		return
	}
	// The eth/65 protocol introduces proper transaction announcements, so instead
	// of dripping transactions across multiple peers, just send the entire list as
	// an announcement and let the remote side decide what they need (likely nothing).
	if p.Version() >= eth.ETH65 {
		hashes := make([]common.Hash, len(txs))
		for i, tx := range txs {
			hashes[i] = tx.Hash()
		}
		p.AsyncSendPooledTransactionHashes(hashes)
		return
	}

	select {
	case h.txsyncCh <- &txsync{p: p, txs: txs}:
	case <-h.quitSync:
	}
}

// txsyncLoop64 takes care of the initial transaction sync for each new
// connection. When a new peer appears, we relay all currently pending
// transactions. In order to minimise egress bandwidth usage, we send
// the transactions in small packs to one peer at a time.
func (h *handler) txsyncLoop64() {
	defer h.wg.Done()

	var (
		pending = make(map[enode.ID]*txsync)
		sending = false               // whether a send is active
		pack    = new(txsync)         // the pack that is being sent
		done    = make(chan error, 1) // result of the send
	)

	// send starts a sending a pack of transactions from the sync.
	send := func(s *txsync) {
		if s.p.Version() >= eth.ETH65 {
			panic("initial transaction syncer running on eth/65+")
		}
		// Fill pack with transactions up to the target size.
		size := common.StorageSize(0)
		pack.p = s.p
		pack.txs = pack.txs[:0]
		for i := 0; i < len(s.txs) && size < txsyncPackSize; i++ {
			pack.txs = append(pack.txs, s.txs[i])
			size += s.txs[i].Size()
		}
		// Remove the transactions that will be sent.
		s.txs = s.txs[:copy(s.txs, s.txs[len(pack.txs):])]
		if len(s.txs) == 0 {
			delete(pending, s.p.Peer.ID())
		}
		// Send the pack in the background.
		s.p.Log().Trace("Sending batch of transactions", "count", len(pack.txs), "bytes", size)
		sending = true
		go func() { done <- pack.p.SendTransactions(pack.txs) }()
	}
	// pick chooses the next pending sync.
	pick := func() *txsync {
		if len(pending) == 0 {
			return nil
		}
		n := rand.Intn(len(pending)) + 1
		for _, s := range pending {
			if n--; n == 0 {
				return s
			}
		}
		return nil
	}

	for {
		select {
		case s := <-h.txsyncCh:
			pending[s.p.Peer.ID()] = s
			if !sending {
				send(s)
			}
		case err := <-done:
			sending = false
			// Stop tracking peers that cause send failures.
			if err != nil {
				pack.p.Log().Debug("Transaction send failed", "err", err)
				delete(pending, pack.p.Peer.ID())
			}
			// Schedule the next send.
			if s := pick(); s != nil {
				send(s)
			}
		case <-h.quitSync:
			return
		}
	}
}

// chainSyncer coordinates blockchain sync components.
type chainSyncer struct {
	handler     *handler
	force       *time.Timer
	forced      bool // true when force timer fired
	peerEventCh chan struct{}
	doneCh      chan error // non-nil when sync is running
}

// chainSyncOp is a scheduled sync operation.
type chainSyncOp struct {
	mode downloader.SyncMode
	peer *eth.Peer
	head common.Hash
	bn   *big.Int
	diff *big.Int
}

// newChainSyncer creates a chainSyncer.
func newChainSyncer(handler *handler) *chainSyncer {
	return &chainSyncer{
		handler:     handler,
		peerEventCh: make(chan struct{}),
	}
}

// handlePeerEvent notifies the syncer about a change in the peer set.
// This is called for new peers and every time a peer announces a new
// chain head.
func (cs *chainSyncer) handlePeerEvent(peer *eth.Peer) bool {
	select {
	case cs.peerEventCh <- struct{}{}:
		return true
	case <-cs.handler.quitSync:
		return false
	}
}

// loop runs in its own goroutine and launches the sync when necessary.
func (cs *chainSyncer) loop() {
	defer cs.handler.wg.Done()

	cs.handler.blockFetcher.Start()
	cs.handler.txFetcher.Start()
	defer cs.handler.blockFetcher.Stop()
	defer cs.handler.txFetcher.Stop()
	defer cs.handler.downloader.Terminate()

	// The force timer lowers the peer count threshold down to one when it fires.
	// This ensures we'll always start sync even if there aren't enough peers.
	cs.force = time.NewTimer(forceSyncCycle)
	defer cs.force.Stop()

	for {
		if op := cs.nextSyncOp(); op != nil {
			cs.startSync(op)
		}
		select {
		case <-cs.peerEventCh:
			// Peer information changed, recheck.
		case <-cs.doneCh:
			cs.doneCh = nil
			cs.force.Reset(forceSyncCycle)
			cs.forced = false
		case <-cs.force.C:
			cs.forced = true

		case <-cs.handler.quitSync:
			// Disable all insertion on the blockchain. This needs to happen before
			// terminating the downloader because the downloader waits for blockchain
			// inserts, and these can take a long time to finish.
			cs.handler.chain.StopInsert()
			cs.handler.downloader.Terminate()
			if cs.doneCh != nil {
				<-cs.doneCh
			}
			return
		}
	}
}

// nextSyncOp determines whether sync is required at this time.
func (cs *chainSyncer) nextSyncOp() *chainSyncOp {
	if cs.doneCh != nil {
		return nil // Sync already running.
	}

	// Ensure we're at minimum peer count.
	minPeers := defaultMinSyncPeers
	if cs.forced {
		minPeers = 1
	} else if minPeers > cs.handler.maxPeers {
		minPeers = cs.handler.maxPeers
	}
	if cs.handler.peers.len() < minPeers {
		return nil
	}
	// We have enough peers, check highest block number
	peer, pBn := cs.handler.peers.peerWithHighestBlock()
	if peer == nil {
		return nil
	}
	mode, ourHighest := cs.modeAndLocalHead()
	if pBn.Uint64() <= ourHighest {
		return nil
	}
	diff := new(big.Int).Sub(pBn, new(big.Int).SetUint64(ourHighest))

	if diff.Cmp(big.NewInt(5)) < 0 {
		return nil
	}

	if mode == downloader.FastSync && atomic.LoadUint32(&cs.handler.snapSync) == 1 {
		// Fast sync via the snap protocol
		mode = downloader.SnapSync
	}
	head, bn := peer.Head()
	return &chainSyncOp{mode: mode, peer: peer, head: head, bn: new(big.Int).Set(bn), diff: diff}
}

func (cs *chainSyncer) modeAndLocalHead() (downloader.SyncMode, uint64) {
	// If we're in fast sync mode, return that directly
	ehead := cs.handler.engine.CurrentBlock()
	if ehead.NumberU64() > 0 {
		log.Info("Blockchain not empty, auto disabling fast sync")
		atomic.StoreUint32(&cs.handler.fastSync, 0)
		return downloader.FullSync, ehead.NumberU64()
	}
	if atomic.LoadUint32(&cs.handler.fastSync) == 1 {
		return downloader.FastSync, ehead.NumberU64()
	}

	// We are probably in full sync, but we might have rewound to before the
	// fast sync pivot, check if we should reenable
	/*if pivot := rawdb.ReadLastPivotNumber(cs.handler.database); pivot != nil {
		if head := cs.handler.chain.CurrentBlock(); head.NumberU64() < *pivot {
			return downloader.FastSync, head.NumberU64()
		}
	}*/
	// Nope, we're really full syncing
	return downloader.FullSync, ehead.NumberU64()
}

// startSync launches doSync in a new goroutine.
func (cs *chainSyncer) startSync(op *chainSyncOp) {
	cs.doneCh = make(chan error, 1)
	go func() { cs.doneCh <- cs.handler.doSync(op) }()
}

// doSync synchronizes the local blockchain with a remote peer.
func (h *handler) doSync(op *chainSyncOp) error {
	if op.mode == downloader.FastSync || op.mode == downloader.SnapSync {
		// Before launch the fast sync, we have to ensure user uses the same
		// txlookup limit.
		// The main concern here is: during the fast sync Geth won't index the
		// block(generate tx indices) before the HEAD-limit. But if user changes
		// the limit in the next fast sync(e.g. user kill Geth manually and
		// restart) then it will be hard for Geth to figure out the oldest block
		// has been indexed. So here for the user-experience wise, it's non-optimal
		// that user can't change limit during the fast sync. If changed, Geth
		// will just blindly use the original one.
		limit := h.chain.TxLookupLimit()
		if stored := rawdb.ReadFastTxLookupLimit(h.database); stored == nil {
			rawdb.WriteFastTxLookupLimit(h.database, limit)
		} else if *stored != limit {
			h.chain.SetTxLookupLimit(*stored)
			log.Warn("Update txLookup limit", "provided", limit, "updated", *stored)
		}
	}
	//wn chain is syncing,keep the chain not receive txs
	if op.diff.Cmp(big.NewInt(5)) > 0 {
		atomic.StoreUint32(&h.acceptTxs, 0)
		defer atomic.StoreUint32(&h.acceptTxs, 1) // Mark initial sync done
	}

	// Run the sync cycle, and disable fast sync if we're past the pivot block
	err := h.downloader.Synchronise(op.peer.ID(), op.head, op.bn, op.mode)
	if err != nil {
		log.Debug("doSync synchronise fail", "err", err)
		return err
	}
	if atomic.LoadUint32(&h.fastSync) == 1 {
		log.Info("Fast sync complete, auto disabling")
		atomic.StoreUint32(&h.fastSync, 0)
	}
	if atomic.LoadUint32(&h.snapSync) == 1 {
		log.Info("Snap sync complete, auto disabling")
		atomic.StoreUint32(&h.snapSync, 0)
	}
	// If we've successfully finished a sync cycle and passed any required checkpoint,
	// enable accepting transactions from the network.
	head := h.chain.CurrentBlock()
	if head.NumberU64() > 0 {
		// We've completed a sync cycle, notify all peers of new state. This path is
		// essential in star-topology networks where a gateway node needs to notify
		// all its out-of-date peers of the availability of a new block. This failure
		// scenario will most often crop up in private and hackathon networks with
		// degenerate connectivity, but it should be healthy for the mainnet too to
		// more reliably update peers or the local TD state.
		h.BroadcastBlock(head, false)
	}
	return nil
}
