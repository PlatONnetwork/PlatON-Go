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
	"sync/atomic"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/eth/protocols/eth"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
)

// Tests that fast sync is disabled after a successful sync cycle.
func TestFastSyncDisabling65(t *testing.T) { testFastSyncDisabling(t, eth.ETH65) }
func TestFastSyncDisabling66(t *testing.T) { testFastSyncDisabling(t, eth.ETH66) }

// Tests that fast sync gets disabled as soon as a real block is successfully
// imported into the blockchain.
func testFastSyncDisabling(t *testing.T, protocol uint) {
	t.Parallel()

	// Create an empty handler and ensure it's in fast sync mode
	empty := newTestHandler()
	if atomic.LoadUint32(&empty.handler.fastSync) == 0 {
		t.Fatalf("fast sync disabled on pristine blockchain")
	}
	defer empty.close()

	// Create a full handler and ensure fast sync ends up disabled
	full := newTestHandlerWithBlocks2(1024)
	if atomic.LoadUint32(&full.handler.fastSync) == 1 {
		t.Fatalf("fast sync not disabled on non-empty blockchain")
	}
	defer full.close()

	// Sync up the two handlers
	emptyPipe, fullPipe := p2p.MsgPipe()
	defer emptyPipe.Close()
	defer fullPipe.Close()

	emptyPeer := eth.NewPeer(protocol, p2p.NewPeer(enode.ID{1}, "", nil), emptyPipe, empty.txpool, nil)
	fullPeer := eth.NewPeer(protocol, p2p.NewPeer(enode.ID{2}, "", nil), fullPipe, full.txpool, nil)
	defer emptyPeer.Close()
	defer fullPeer.Close()

	go empty.handler.runEthPeer(emptyPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(empty.handler), peer)
	})
	go full.handler.runEthPeer(fullPeer, func(peer *eth.Peer) error {
		return eth.Handle((*ethHandler)(full.handler), peer)
	})
	// Wait a bit for the above handlers to start
	time.Sleep(250 * time.Millisecond)

	// Check that fast sync was disabled
	p, _ := empty.handler.peers.peerWithHighestBlock()
	head, bn := p.Head()
	op := &chainSyncOp{mode: downloader.FastSync, peer: p, head: head, bn: bn, diff: new(big.Int).SetUint64(0)}
	if err := empty.handler.doSync(op); err != nil {
		t.Fatal("sync failed:", err)
	}
	if atomic.LoadUint32(&empty.handler.fastSync) == 1 {
		t.Fatalf("fast sync not disabled after successful synchronisation")
	}
}
