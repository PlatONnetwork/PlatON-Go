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
	"encoding/binary"
	"math"
	"math/big"
	"math/rand"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/eth/downloader"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

// Tests that protocol versions and modes of operations are matched up properly.
func TestProtocolCompatibility(t *testing.T) {
	// Define the compatibility chart
	tests := []struct {
		version    uint
		mode       downloader.SyncMode
		compatible bool
	}{
		{61, downloader.FullSync, true}, {62, downloader.FullSync, true}, {63, downloader.FullSync, true},
		{61, downloader.FastSync, false}, {62, downloader.FastSync, false}, {63, downloader.FastSync, true},
	}
	// Make sure anything we screw up is restored
	backup := ProtocolVersions
	defer func() { ProtocolVersions = backup }()
	// Try all available compatibility configs and check for errors
	for i, tt := range tests {
		ProtocolVersions = []uint{tt.version}
		//db := ethdb.NewMemDatabase()
		pm, _, err := newTestProtocolManager(tt.mode, 0, nil, nil)
		if pm != nil {
			defer pm.Stop()
		}
		if (err == nil && !tt.compatible) || (err != nil && tt.compatible) {
			t.Errorf("test %d: compatibility mismatch: have error %v, want compatibility %v", i, err, tt.compatible)
		}
	}
}

// Tests that block headers can be retrieved from a remote chain based on user queries.
func TestGetBlockHeaders62(t *testing.T) { testGetBlockHeaders(t, 62) }

func TestGetBlockHeaders63(t *testing.T) { testGetBlockHeaders(t, 63) }

func testGetBlockHeaders(t *testing.T, protocol int) {
	//db := ethdb.NewMemDatabase()
	pm, _ := newTestProtocolManagerMust(t, downloader.FullSync, downloader.MaxHashFetch+15, nil, nil)
	peer, _ := newTestPeer("peer", protocol, pm, true)
	defer peer.close()

	// Create a "random" unknown hash for testing
	var unknown common.Hash
	for i := range unknown {
		unknown[i] = byte(i)
	}
	// Create a batch of tests for various scenarios
	limit := uint64(downloader.MaxHeaderFetch)
	tests := []struct {
		query  *getBlockHeadersData // The query to execute for header retrieval
		expect []common.Hash        // The hashes of the block whose headers are expected
	}{
		// A single random block should be retrievable by hash and number too
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: pm.blockchain.GetBlockByNumber(limit / 2).Hash()}, Amount: 1},
			[]common.Hash{pm.blockchain.GetBlockByNumber(limit / 2).Hash()},
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Amount: 1},
			[]common.Hash{pm.blockchain.GetBlockByNumber(limit / 2).Hash()},
		},
		// Multiple headers should be retrievable in both directions
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Amount: 3},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 1).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 2).Hash(),
			},
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Amount: 3, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 1).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 2).Hash(),
			},
		},
		// Multiple headers with skip lists should be retrievable
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 4).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 + 8).Hash(),
			},
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: limit / 2}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(limit / 2).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 4).Hash(),
				pm.blockchain.GetBlockByNumber(limit/2 - 8).Hash(),
			},
		},
		// The chain endpoints should be retrievable
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 0}, Amount: 1},
			[]common.Hash{pm.blockchain.GetBlockByNumber(0).Hash()},
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64()}, Amount: 1},
			[]common.Hash{pm.blockchain.CurrentBlock().Hash()},
		},
		// Ensure protocol limits are honored
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() - 1}, Amount: limit + 10, Reverse: true},
			pm.blockchain.GetBlockHashesFromHash(pm.blockchain.CurrentBlock().Hash(), limit),
		},
		// Check that requesting more than available is handled gracefully
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() - 4}, Skip: 3, Amount: 3},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64() - 4).Hash(),
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64()).Hash(),
			},
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: 4}, Skip: 3, Amount: 3, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(4).Hash(),
				pm.blockchain.GetBlockByNumber(0).Hash(),
			},
		},
		// Check that requesting more than available is handled gracefully, even if mid skip
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() - 4}, Skip: 2, Amount: 3},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64() - 4).Hash(),
				pm.blockchain.GetBlockByNumber(pm.blockchain.CurrentBlock().NumberU64() - 1).Hash(),
			},
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: 4}, Skip: 2, Amount: 3, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(4).Hash(),
				pm.blockchain.GetBlockByNumber(1).Hash(),
			},
		},
		// Check a corner case where requesting more can iterate past the endpoints
		{
			&getBlockHeadersData{Origin: hashOrNumber{Number: 2}, Amount: 5, Reverse: true},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(2).Hash(),
				pm.blockchain.GetBlockByNumber(1).Hash(),
				pm.blockchain.GetBlockByNumber(0).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back into the chain start
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: pm.blockchain.GetBlockByNumber(3).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64 - 1},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(3).Hash(),
			},
		},
		// Check a corner case where skipping overflow loops back to the same header
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: pm.blockchain.GetBlockByNumber(1).Hash()}, Amount: 2, Reverse: false, Skip: math.MaxUint64},
			[]common.Hash{
				pm.blockchain.GetBlockByNumber(1).Hash(),
			},
		},
		// Check that non existing headers aren't returned
		{
			&getBlockHeadersData{Origin: hashOrNumber{Hash: unknown}, Amount: 1},
			[]common.Hash{},
		}, {
			&getBlockHeadersData{Origin: hashOrNumber{Number: pm.blockchain.CurrentBlock().NumberU64() + 1}, Amount: 1},
			[]common.Hash{},
		},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the headers to expect in the response
		headers := []*types.Header{}
		for _, hash := range tt.expect {
			headers = append(headers, pm.blockchain.GetBlockByHash(hash).Header())
		}
		// Send the hash request and verify the response
		p2p.Send(peer.app, 0x03, tt.query)
		if err := p2p.ExpectMsg(peer.app, 0x04, headers); err != nil {
			t.Errorf("test %d: headers mismatch: %v", i, err)
		}
		// If the test used number origins, repeat with hashes as the too
		if tt.query.Origin.Hash == (common.Hash{}) {
			if origin := pm.blockchain.GetBlockByNumber(tt.query.Origin.Number); origin != nil {
				tt.query.Origin.Hash, tt.query.Origin.Number = origin.Hash(), 0

				p2p.Send(peer.app, 0x03, tt.query)
				if err := p2p.ExpectMsg(peer.app, 0x04, headers); err != nil {
					t.Errorf("test %d: headers mismatch: %v", i, err)
				}
			}
		}
	}
}

// Tests that block contents can be retrieved from a remote chain based on their hashes.
func TestGetBlockBodies62(t *testing.T) { testGetBlockBodies(t, 62) }

func TestGetBlockBodies63(t *testing.T) { testGetBlockBodies(t, 63) }

func testGetBlockBodies(t *testing.T, protocol int) {
	pm, _ := newTestProtocolManagerMust(t, downloader.FullSync, downloader.MaxBlockFetch+15, nil, nil)
	peer, _ := newTestPeer("peer", protocol, pm, true)
	defer peer.close()

	// Create a batch of tests for various scenarios
	limit := downloader.MaxBlockFetch
	tests := []struct {
		random    int           // Number of blocks to fetch randomly from the chain
		explicit  []common.Hash // Explicitly requested blocks
		available []bool        // Availability of explicitly requested blocks
		expected  int           // Total number of existing blocks to expect
	}{
		{1, nil, nil, 1},             // A single random block should be retrievable
		{10, nil, nil, 10},           // Multiple random blocks should be retrievable
		{limit, nil, nil, limit},     // The maximum possible blocks should be retrievable
		{limit + 1, nil, nil, limit}, // No more than the possible block count should be returned
		{0, []common.Hash{pm.blockchain.Genesis().Hash()}, []bool{true}, 1},      // The genesis block should be retrievable
		{0, []common.Hash{pm.blockchain.CurrentBlock().Hash()}, []bool{true}, 1}, // The chains head block should be retrievable
		{0, []common.Hash{{}}, []bool{false}, 0},                                 // A non existent block should not be returned

		// Existing and non-existing blocks interleaved should not cause problems
		{0, []common.Hash{
			{},
			pm.blockchain.GetBlockByNumber(1).Hash(),
			{},
			pm.blockchain.GetBlockByNumber(10).Hash(),
			{},
			pm.blockchain.GetBlockByNumber(100).Hash(),
			{},
		}, []bool{false, true, false, true, false, true, false}, 3},
	}
	// Run each of the tests and verify the results against the chain
	for i, tt := range tests {
		// Collect the hashes to request, and the response to expect
		hashes, seen := []common.Hash{}, make(map[int64]bool)
		bodies := []*blockBody{}

		for j := 0; j < tt.random; j++ {
			for {
				num := rand.Int63n(int64(pm.blockchain.CurrentBlock().NumberU64()))
				if !seen[num] {
					seen[num] = true

					block := pm.blockchain.GetBlockByNumber(uint64(num))
					hashes = append(hashes, block.Hash())
					if len(bodies) < tt.expected {
						bodies = append(bodies, &blockBody{Transactions: block.Transactions()})
					}
					break
				}
			}
		}
		for j, hash := range tt.explicit {
			hashes = append(hashes, hash)
			if tt.available[j] && len(bodies) < tt.expected {
				block := pm.blockchain.GetBlockByHash(hash)
				bodies = append(bodies, &blockBody{Transactions: block.Transactions()})
			}
		}
		// Send the hash request and verify the response
		p2p.Send(peer.app, 0x05, hashes)
		if err := p2p.ExpectMsg(peer.app, 0x06, bodies); err != nil {
			t.Errorf("test %d: bodies mismatch: %v", i, err)
		}
	}
}

func newSnapshotdb() (snapshotdb.DB, error) {
	db := snapshotdb.Instance()
	highth := big.NewInt(100)
	if err := db.SetCurrent(common.ZeroHash, *highth, *highth); err != nil {
		return nil, err
	}
	var m [][2][]byte
	for i := 0; i < 1028; i++ {
		b := make([]byte, 8)
		binary.LittleEndian.PutUint64(b, rand.Uint64())
		m = append(m, [2][]byte{b, b})
	}
	if err := db.WriteBaseDB(m); err != nil {
		return nil, err
	}
	return db, nil
}

func TestGetOriginAndPivotMsg(t *testing.T) {
	//	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	pm, _ := newTestProtocolManagerMust(t, downloader.FastSync, downloader.MaxBlockFetch+15, nil, nil)
	peer, _ := newTestPeer("peer", 63, pm, true)
	db, err := newSnapshotdb()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		peer.close()
		db.Clear()
	}()

	if err := p2p.Send(peer.app, GetOriginAndPivotMsg, pm.blockchain.CurrentBlock().Header().Number.Uint64()); err != nil {
		t.Error(err)
		return
	}
	data := make([]*types.Header, 0)

	oHead := pm.blockchain.GetHeaderByNumber(pm.blockchain.CurrentBlock().Header().Number.Uint64())
	pivot, err := db.BaseNum()
	if err != nil {
		t.Error(err)
		return
	}
	pHead := pm.blockchain.GetHeaderByNumber(pivot.Uint64())
	data = append(data, oHead, pHead)
	if err := p2p.ExpectMsg(peer.app, OriginAndPivotMsg, &data); err != nil {
		t.Errorf("TestGetOriginAndPivotMsg : data mismatch: %v", err)
	}
}

func TestGetPPOSStorageMsg(t *testing.T) {
	//	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	pm, _ := newTestProtocolManagerMust(t, downloader.FullSync, downloader.MaxBlockFetch+15, nil, nil)
	peer, _ := newTestPeer("peer", 63, pm, true)
	db, err := newSnapshotdb()
	if err != nil {
		t.Error(err)
		return
	}
	defer func() {
		peer.close()
		db.Clear()
	}()
	if err := p2p.Send(peer.app, GetPPOSStorageMsg, []interface{}{}); err != nil {
		t.Error(err)
		return
	}
	var data PPOSInfo
	msg, err := peer.app.ReadMsg()
	if err != nil {
		t.Error(err)
		return
	}
	if err := msg.Decode(&data); err != nil {
		t.Error(err)
		return
	}
	if data.Pivot.Number.Int64() != 100 {
		t.Error("Pivot num is wrong")
	}
	if data.Latest.Number.Int64() != int64(downloader.MaxBlockFetch+15) {
		t.Error("latest num is wrong")
	}
	var data2 PPOSStorage

	for {
		msg, err := peer.app.ReadMsg()
		if err != nil {
			t.Error(err)
			return
		}
		if err := msg.Decode(&data2); err != nil {
			t.Error(err)
			return
		}
		if data2.Last {
			break
		}
	}
}

// Tests that the node state database can be retrieved based on hashes.
func TestGetNodeData63(t *testing.T) { testGetNodeData(t, 63) }

func testGetNodeData(t *testing.T, protocol int) {
	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr := crypto.PubkeyToAddress(acc2Key.PublicKey)
	//chanid := new(big.Int)
	signer := types.NewEIP155Signer(big.NewInt(1))
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *core.BlockGen) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some ether.
			tx, _ := types.SignTx(types.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(10000), params.TxGas, nil, nil), signer, testBankKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more ether to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := types.SignTx(types.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, testBankKey)
			tx2, _ := types.SignTx(types.NewTransaction(block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was mined by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetExtra([]byte("yeehaw"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			b2 := block.PrevBlock(1).Header()
			b2.Extra = []byte("foo")
			b3 := block.PrevBlock(2).Header()
			b3.Extra = []byte("foo")
		}
	}
	// Assemble the test environment
	pm, db := newTestProtocolManagerMust(t, downloader.FullSync, 4, generator, nil)
	peer, _ := newTestPeer("peer", protocol, pm, true)
	defer peer.close()

	// Fetch for now the entire chain db
	hashes := []common.Hash{}
	it := db.NewIterator()
	for it.Next() {
		if key := it.Key(); len(key) == common.HashLength {
			hashes = append(hashes, common.BytesToHash(key))
		}
	}
	it.Release()
	p2p.Send(peer.app, 0x0d, hashes)
	msg, err := peer.app.ReadMsg()
	if err != nil {
		t.Fatalf("failed to read node data response: %v", err)
	}
	if msg.Code != 0x0e {
		t.Fatalf("response packet code mismatch: have %x, want %x", msg.Code, 0x0c)
	}
	var data [][]byte
	if err := msg.Decode(&data); err != nil {
		t.Fatalf("failed to decode response node data: %v", err)
	}
	// Verify that all hashes correspond to the requested data, and reconstruct a state tree
	for i, want := range hashes {
		if hash := crypto.Keccak256Hash(data[i]); hash != want {
			t.Errorf("data hash mismatch: have %x, want %x", hash, want)
		}
	}
	statedb := rawdb.NewMemoryDatabase()
	for i := 0; i < len(data); i++ {
		statedb.Put(hashes[i].Bytes(), data[i])
	}
	accounts := []common.Address{testBank, acc1Addr, acc2Addr}
	for i := uint64(0); i <= pm.blockchain.CurrentBlock().NumberU64(); i++ {
		trie, _ := state.New(pm.blockchain.GetBlockByNumber(i).Root(), state.NewDatabase(statedb))

		for j, acc := range accounts {
			state, _ := pm.blockchain.State()
			bw := state.GetBalance(acc)
			bh := trie.GetBalance(acc)

			if (bw != nil && bh == nil) || (bw == nil && bh != nil) {
				t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
			}
			if bw != nil && bh != nil && bw.Cmp(bw) != 0 {
				t.Errorf("test %d, account %d: balance mismatch: have %v, want %v", i, j, bh, bw)
			}
		}
	}
}

// Tests that the transaction receipts can be retrieved based on hashes.
func TestGetReceipt63(t *testing.T) { testGetReceipt(t, 63) }

func testGetReceipt(t *testing.T, protocol int) {
	// Define three accounts to simulate transactions with
	acc1Key, _ := crypto.HexToECDSA("8a1f9a8f95be41cd7ccb6168179afb4504aefe388d1e14474d32c45c72ce7b7a")
	acc2Key, _ := crypto.HexToECDSA("49a7b37aa6f6645917e7b807e9d1c00d4fa71f18343b0d4122a4d2df64dd6fee")
	acc1Addr := crypto.PubkeyToAddress(acc1Key.PublicKey)
	acc2Addr := crypto.PubkeyToAddress(acc2Key.PublicKey)

	signer := types.NewEIP155Signer(big.NewInt(1))
	// Create a chain generator with some simple transactions (blatantly stolen from @fjl/chain_markets_test)
	generator := func(i int, block *core.BlockGen) {
		switch i {
		case 0:
			// In block 1, the test bank sends account #1 some ether.
			tx, _ := types.SignTx(types.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(10000), params.TxGas, nil, nil), signer, testBankKey)
			block.AddTx(tx)
		case 1:
			// In block 2, the test bank sends some more ether to account #1.
			// acc1Addr passes it on to account #2.
			tx1, _ := types.SignTx(types.NewTransaction(block.TxNonce(testBank), acc1Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, testBankKey)
			tx2, _ := types.SignTx(types.NewTransaction(block.TxNonce(acc1Addr), acc2Addr, big.NewInt(1000), params.TxGas, nil, nil), signer, acc1Key)
			block.AddTx(tx1)
			block.AddTx(tx2)
		case 2:
			// Block 3 is empty but was mined by account #2.
			block.SetCoinbase(acc2Addr)
			block.SetExtra([]byte("yeehaw"))
		case 3:
			// Block 4 includes blocks 2 and 3 as uncle headers (with modified extra data).
			b2 := block.PrevBlock(1).Header()
			b2.Extra = []byte("foo")
			b3 := block.PrevBlock(2).Header()
			b3.Extra = []byte("foo")
		}
	}
	// Assemble the test environment
	pm, _ := newTestProtocolManagerMust(t, downloader.FullSync, 4, generator, nil)
	peer, _ := newTestPeer("peer", protocol, pm, true)
	defer peer.close()

	// Collect the hashes to request, and the response to expect
	hashes, receipts := []common.Hash{}, []types.Receipts{}
	for i := uint64(0); i <= pm.blockchain.CurrentBlock().NumberU64(); i++ {
		block := pm.blockchain.GetBlockByNumber(i)

		hashes = append(hashes, block.Hash())
		receipts = append(receipts, pm.blockchain.GetReceiptsByHash(block.Hash()))
	}
	// Send the hash request and verify the response
	p2p.Send(peer.app, 0x0f, hashes)
	if err := p2p.ExpectMsg(peer.app, 0x10, receipts); err != nil {
		t.Errorf("receipts mismatch: %v", err)
	}
}
