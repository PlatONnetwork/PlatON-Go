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
	"math/big"
	"math/rand"
	"sync"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/event"

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
		{61, downloader.SnapSync, false}, {62, downloader.SnapSync, false}, {63, downloader.SnapSync, true},
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

var (
	// testKey is a private key to use for funding a tester account.
	testKey, _ = crypto.HexToECDSA("b71c71a67e1177ad4e901695e1b4b9ee17ae16c6668d313eac2f96dbcda3f291")
	// testAddr is the Ethereum address of the tester account.
	testAddr = crypto.PubkeyToAddress(testKey.PublicKey)
)

// testTxPool is a mock transaction pool that blindly accepts all transactions.
// Its goal is to get around setting up a valid statedb for the balance and nonce
// checks.
type testTxPool struct {
	pool   map[common.Hash]*types.Transaction // Hash map of collected transactions
	txFeed event.Feed                         // Notification feed to allow waiting for inclusion
	lock   sync.RWMutex                       // Protects the transaction pool
}

// newTestTxPool creates a mock transaction pool.
func newTestTxPool() *testTxPool {
	return &testTxPool{
		pool: make(map[common.Hash]*types.Transaction),
	}
}

// Has returns an indicator whether txpool has a transaction
// cached with the given hash.
func (p *testTxPool) Has(hash common.Hash) bool {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.pool[hash] != nil
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

	pm, _ := newTestProtocolManagerMust(t, downloader.SnapSync, downloader.MaxBlockFetch+15, nil, nil)
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
	it := db.NewIterator(nil, nil)
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
		trie, _ := state.New(pm.blockchain.GetBlockByNumber(i).Root(), state.NewDatabase(statedb), nil)

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
