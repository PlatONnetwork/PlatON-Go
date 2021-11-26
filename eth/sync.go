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

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (
	forceSyncCycle      = 5 * time.Second // Time interval to force syncs, even if few peers are available
	defaultMinSyncPeers = 5               // Amount of peers desired to start syncing

	// This is the target size for the packs of transactions sent by txsyncLoop.
	// A pack can get larger than this if a single transactions exceeds this size.
	txsyncPackSize = 100 * 1024
)

type txsync struct {
	p   *peer
	txs []*types.Transaction
}

// syncTransactions starts sending all currently pending transactions to the given peer.
func (pm *ProtocolManager) syncTransactions(p *peer) {
	var txs types.Transactions
	pending, _ := pm.txpool.Pending()
	for _, batch := range pending {
		txs = append(txs, batch...)
	}
	if len(txs) == 0 {
		return
	}
	select {
	case pm.txsyncCh <- &txsync{p, txs}:
	case <-pm.quitSync:
	}
}

// txsyncLoop takes care of the initial transaction sync for each new
// connection. When a new peer appears, we relay all currently pending
// transactions. In order to minimise egress bandwidth usage, we send
// the transactions in small packs to one peer at a time.
func (pm *ProtocolManager) txsyncLoop() {
	defer pm.wg.Done()
	var (
		pending = make(map[discover.NodeID]*txsync)
		sending = false               // whether a send is active
		pack    = new(txsync)         // the pack that is being sent
		done    = make(chan error, 1) // result of the send
	)

	// send starts a sending a pack of transactions from the sync.
	send := func(s *txsync) {
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
			delete(pending, s.p.ID())
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
		case s := <-pm.txsyncCh:
			pending[s.p.ID()] = s
			if !sending {
				send(s)
			}
		case err := <-done:
			sending = false
			// Stop tracking peers that cause send failures.
			if err != nil {
				pack.p.Log().Debug("Transaction send failed", "err", err)
				delete(pending, pack.p.ID())
			}
			// Schedule the next send.
			if s := pick(); s != nil {
				send(s)
			}
		case <-pm.quitSync:
			return
		}
	}
}

// chainSyncer coordinates blockchain sync components.
type chainSyncer struct {
	pm          *ProtocolManager
	force       *time.Timer
	forced      bool // true when force timer fired
	peerEventCh chan struct{}
	doneCh      chan error // non-nil when sync is running
}

// chainSyncOp is a scheduled sync operation.
type chainSyncOp struct {
	mode downloader.SyncMode
	peer *peer
	bn   *big.Int
	head common.Hash
	diff *big.Int
}

// newChainSyncer creates a chainSyncer.
func newChainSyncer(pm *ProtocolManager) *chainSyncer {
	return &chainSyncer{
		pm:          pm,
		peerEventCh: make(chan struct{}),
	}
}

// handlePeerEvent notifies the syncer about a change in the peer set.
// This is called for new peers and every time a peer announces a new
// chain head.
func (cs *chainSyncer) handlePeerEvent(p *peer) bool {
	select {
	case cs.peerEventCh <- struct{}{}:
		return true
	case <-cs.pm.quitSync:
		return false
	}
}

// loop runs in its own goroutine and launches the sync when necessary.
func (cs *chainSyncer) loop() {
	defer cs.pm.wg.Done()

	cs.pm.fetcher.Start()
	cs.pm.txFetcher.Start()
	defer cs.pm.fetcher.Stop()
	defer cs.pm.txFetcher.Stop()
	defer cs.pm.downloader.Terminate()

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

		case <-cs.pm.quitSync:
			if cs.doneCh != nil {
				cs.pm.downloader.Cancel()
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

	// Ensure we're at mininum peer count.
	minPeers := defaultMinSyncPeers
	if cs.forced {
		minPeers = 1
	} else if minPeers > cs.pm.maxPeers {
		minPeers = cs.pm.maxPeers
	}
	if cs.pm.peers.Len() < minPeers {
		return nil
	}

	// We have enough peers, check TD.
	peer := cs.pm.peers.BestPeer()
	if peer == nil {
		return nil
	}

	currentBlock := cs.pm.engine.CurrentBlock()

	peerHead, pBn := peer.Head()
	//modified by alaya
	diff := new(big.Int).Sub(pBn, currentBlock.Number())
	if diff.Cmp(big.NewInt(2)) < 0 {
		return nil
	}

	mode := downloader.FullSync
	if currentBlock.NumberU64() > 0 {
		log.Info("Blockchain not empty, auto disabling fast sync")
		atomic.StoreUint32(&cs.pm.fastSync, 0)
	}

	if atomic.LoadUint32(&cs.pm.fastSync) == 1 {
		// Fast sync was explicitly requested, and explicitly granted
		mode = downloader.FastSync
	} else if currentBlock.NumberU64() == 0 && cs.pm.blockchain.CurrentFastBlock().NumberU64() > 0 {
		// The database seems empty as the current block is the genesis. Yet the fast
		// block is ahead, so fast sync was enabled for this node at a certain point.
		// The only scenario where this can happen is if the user manually (or via a
		// bad block) rolled back a fast sync node below the sync point. In this case
		// however it's safe to reenable fast sync.
		atomic.StoreUint32(&cs.pm.fastSync, 1)
		mode = downloader.FastSync
	}

	if mode == downloader.FastSync && cs.pm.blockchain.CurrentFastBlock().Number().Cmp(pBn) >= 0 {
		return nil
	}

	return &chainSyncOp{mode: mode, peer: peer, bn: pBn, head: peerHead, diff: new(big.Int).Sub(pBn, currentBlock.Number())}
}

// startSync launches doSync in a new goroutine.
func (cs *chainSyncer) startSync(op *chainSyncOp) {
	cs.doneCh = make(chan error, 1)
	go func() { cs.doneCh <- cs.pm.doSync(op) }()
}

// doSync synchronizes the local blockchain with a remote peer.
func (pm *ProtocolManager) doSync(op *chainSyncOp) error {
	//wn chain is syncing,keep the chain not receive txs
	if op.diff.Cmp(big.NewInt(5)) > 0 {
		atomic.StoreUint32(&pm.acceptTxs, 0)
	}
	// Run the sync cycle, and disable fast sync if we're past the pivot block
	err := pm.downloader.Synchronise(op.peer.id, op.head, op.bn, op.mode)
	if err != nil {
		return err
	}
	if atomic.LoadUint32(&pm.fastSync) == 1 {
		log.Info("Fast sync complete, auto disabling")
		atomic.StoreUint32(&pm.fastSync, 0)
	}

	// If we've successfully finished a sync cycle and passed any required checkpoint,
	// enable accepting transactions from the network.
	head := pm.blockchain.CurrentBlock()
	atomic.StoreUint32(&pm.acceptTxs, 1) // Mark initial sync done
	if head.NumberU64() > 0 {
		// We've completed a sync cycle, notify all peers of new state. This path is
		// essential in star-topology networks where a gateway node needs to notify
		// all its out-of-date peers of the availability of a new block. This failure
		// scenario will most often crop up in private and hackathon networks with
		// degenerate connectivity, but it should be healthy for the mainnet too to
		// more reliably update peers or the local TD state.
		pm.BroadcastBlock(head, false)
	}

	return nil
}
