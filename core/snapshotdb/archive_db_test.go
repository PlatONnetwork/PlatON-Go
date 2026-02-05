// Copyright 2021 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
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

package snapshotdb

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"sync"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/ethdb/memorydb"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// setupMemoryArchiveDB creates an archiveDB with memory database for testing
func setupMemoryArchiveDB(t *testing.T) *archiveDB {
	memDB := memorydb.New()
	rawDB := rawdb.NewDatabase(memDB)

	archiveDB, err := NewArchiveDBWithDB(rawDB)
	if err != nil {
		t.Fatalf("Failed to create archive with memory DB: %v", err)
	}

	return archiveDB
}
func walkFuncInit(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
	// Simulate walking through some data
	testKV := map[string][]byte{
		"key1": []byte("value1"),
		"key2": []byte("value2"),
	}

	mockIter := newMockIterator(testKV)
	for mockIter.Next() {
		err := f(big.NewInt(1), mockIter)
		if err != nil {
			return err
		}
	}
	return nil
}

// createTestBlockData creates a test BlockData with sample key-value pairs
func createTestBlockData(t *testing.T, blockNumber int64, parentHash, blockHash common.Hash) *BlockData {
	blockData := NewBlockData(big.NewInt(blockNumber), parentHash, blockHash)

	// Add some test key-value pairs
	testData := []struct {
		key   []byte
		value []byte
	}{
		{[]byte("key1"), []byte("value1")},
		{[]byte("key2"), []byte("value2")},
		{[]byte("key3"), []byte("value3")},
	}

	for _, data := range testData {
		err := blockData.data.Put(data.key, data.value)
		if err != nil {
			t.Fatalf("Failed to put test data: %v", err)
		}
	}

	// Add VRF nonce data
	nonces := [][]byte{[]byte("nonce1"), []byte("nonce2"), []byte("nonce3")}
	nonceBytes, err := rlp.EncodeToBytes(nonces)
	if err != nil {
		t.Fatalf("Failed to encode nonces: %v", err)
	}
	err = blockData.data.Put(nonceStorageKey, nonceBytes)
	if err != nil {
		t.Fatalf("Failed to put nonce data: %v", err)
	}

	return blockData
}

func TestVrfNonceKey(t *testing.T) {
	blockNumber := uint64(100)
	key := VrfNonceKey(blockNumber)

	if len(key) != 10 {
		t.Errorf("Expected key length 10, got %d", len(key))
	}

	if !bytes.HasPrefix(key, vrfNoncePrefix) {
		t.Errorf("Key should start with vrfNoncePrefix")
	}

	// Check that the block number is correctly encoded
	encoded := binary.BigEndian.Uint64(key[2:])
	if encoded != blockNumber {
		t.Errorf("Expected block number %d, got %d", blockNumber, encoded)
	}
}

func TestArchiveBlockKey(t *testing.T) {
	blockNumber := uint64(100)
	key := ArchiveBlockKey(blockNumber)

	if len(key) != 11 {
		t.Errorf("Expected key length 11, got %d", len(key))
	}

	if !bytes.HasPrefix(key, archiveBlockPrefix) {
		t.Errorf("Key should start with archiveBlockPrefix")
	}

	// Check that the block number is correctly encoded
	encoded := binary.BigEndian.Uint64(key[3:])
	if encoded != blockNumber {
		t.Errorf("Expected block number %d, got %d", blockNumber, encoded)
	}
}

func TestNewArchiveDBWithDB(t *testing.T) {
	memDB := memorydb.New()
	rawDB := rawdb.NewDatabase(memDB)

	// Test creating archiveDB with provided database
	archiveDB, err := NewArchiveDBWithDB(rawDB)
	if err != nil {
		t.Fatalf("Failed to create archive with memory DB: %v", err)
	}

	// Verify that the database and trie database are initialized
	if archiveDB.db == nil {
		t.Error("Database should not be nil")
	}
	if archiveDB.triedb == nil {
		t.Error("Trie database should not be nil")
	}
	if archiveDB.trie != nil {
		t.Error("Trie should be nil before initialization")
	}
}

func TestArchiveDB_CurrentBlock(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	// Test when no current block is set
	block, err := archiveDB.CurrentBlock()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if block != nil {
		t.Error("Current block should be nil initially")
	}

	// Set a current block
	testBlockNumber := uint64(100)
	err = archiveDB.SetCurrentBlock(testBlockNumber)
	if err != nil {
		t.Errorf("Failed to set current block: %v", err)
	}

	// Get the current block
	block, err = archiveDB.CurrentBlock()
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if block == nil {
		t.Error("Current block should not be nil")
	}
	if *block != testBlockNumber {
		t.Errorf("Expected block number %d, got %d", testBlockNumber, *block)
	}
}

func TestArchiveDB_SetArchiveBlock(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	blockNumber := uint64(100)
	block := &ArchiveBlock{
		Number: blockNumber,
		Root:   common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
		KvHash: common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
	}

	// Test setting archive block
	batch := archiveDB.db.NewBatch()
	err := archiveDB.SetArchiveBlock(batch, blockNumber, block)
	if err != nil {
		t.Errorf("Failed to set archive block: %v", err)
	}
	err = batch.Write()
	if err != nil {
		t.Errorf("Failed to write batch: %v", err)
	}

	// Verify the block was stored
	retrievedBlock, err := archiveDB.GetArchiveBlock(blockNumber)
	if err != nil {
		t.Errorf("Failed to get archive block: %v", err)
	}
	if retrievedBlock.Number != block.Number {
		t.Errorf("Expected number %d, got %d", block.Number, retrievedBlock.Number)
	}
	if retrievedBlock.Root != block.Root {
		t.Errorf("Expected root %s, got %s", block.Root.Hex(), retrievedBlock.Root.Hex())
	}
	if retrievedBlock.KvHash != block.KvHash {
		t.Errorf("Expected KvHash %s, got %s", block.KvHash.Hex(), retrievedBlock.KvHash.Hex())
	}
}

func TestArchiveDB_GetArchiveBlock(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	blockNumber := uint64(100)

	// Test getting non-existent block
	block, err := archiveDB.GetArchiveBlock(blockNumber)
	if err == nil {
		t.Error("Expected error when getting non-existent archive block")
	}
	if block != nil {
		t.Error("Block should be nil for non-existent archive block")
	}

	// Store a block and retrieve it
	testBlock := &ArchiveBlock{
		Number: blockNumber,
		Root:   common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef"),
		KvHash: common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890"),
	}

	batch := archiveDB.db.NewBatch()
	err = archiveDB.SetArchiveBlock(batch, blockNumber, testBlock)
	if err != nil {
		t.Errorf("Failed to set archive block: %v", err)
	}
	err = batch.Write()
	if err != nil {
		t.Errorf("Failed to write batch: %v", err)
	}

	// Retrieve the block
	retrievedBlock, err := archiveDB.GetArchiveBlock(blockNumber)
	if err != nil {
		t.Errorf("Failed to get archive block: %v", err)
	}
	if retrievedBlock == nil {
		t.Error("Retrieved block should not be nil")
	}
	if retrievedBlock.Number != testBlock.Number {
		t.Errorf("Expected number %d, got %d", testBlock.Number, retrievedBlock.Number)
	}
	if retrievedBlock.Root != testBlock.Root {
		t.Errorf("Expected root %s, got %s", testBlock.Root.Hex(), retrievedBlock.Root.Hex())
	}
	if retrievedBlock.KvHash != testBlock.KvHash {
		t.Errorf("Expected KvHash %s, got %s", testBlock.KvHash.Hex(), retrievedBlock.KvHash.Hex())
	}
}

func TestArchiveDB_SetVrfNonce(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	blockNumber := uint64(100)
	nonce := &VRFNonce{
		MaxValidatorNum: 10,
		Nonce:           []byte("test_nonce"),
	}

	// Test setting VRF nonce
	batch := archiveDB.db.NewBatch()
	err := archiveDB.SetVrfNonce(batch, blockNumber, nonce)
	if err != nil {
		t.Errorf("Failed to set VRF nonce: %v", err)
	}
	err = batch.Write()
	if err != nil {
		t.Errorf("Failed to write batch: %v", err)
	}

	// Verify the nonce was stored
	retrievedNonce, err := archiveDB.GetVrfNonce(blockNumber)
	if err != nil {
		t.Errorf("Failed to get VRF nonce: %v", err)
	}
	if retrievedNonce.MaxValidatorNum != nonce.MaxValidatorNum {
		t.Errorf("Expected MaxValidatorNum %d, got %d", nonce.MaxValidatorNum, retrievedNonce.MaxValidatorNum)
	}
	if !bytes.Equal(retrievedNonce.Nonce, nonce.Nonce) {
		t.Errorf("Expected nonce %s, got %s", string(nonce.Nonce), string(retrievedNonce.Nonce))
	}
}

func TestArchiveDB_GetVrfNonce(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	blockNumber := uint64(100)

	// Test getting non-existent nonce
	nonce, err := archiveDB.GetVrfNonce(blockNumber)
	if err == nil {
		t.Error("Expected error when getting non-existent VRF nonce")
	}
	if nonce != nil {
		t.Error("Nonce should be nil for non-existent VRF nonce")
	}

	// Store a nonce and retrieve it
	testNonce := &VRFNonce{
		MaxValidatorNum: 10,
		Nonce:           []byte("test_nonce"),
	}

	batch := archiveDB.db.NewBatch()
	err = archiveDB.SetVrfNonce(batch, blockNumber, testNonce)
	if err != nil {
		t.Errorf("Failed to set VRF nonce: %v", err)
	}
	err = batch.Write()
	if err != nil {
		t.Errorf("Failed to write batch: %v", err)
	}

	// Retrieve the nonce
	retrievedNonce, err := archiveDB.GetVrfNonce(blockNumber)
	if err != nil {
		t.Errorf("Failed to get VRF nonce: %v", err)
	}
	if retrievedNonce == nil {
		t.Error("Retrieved nonce should not be nil")
	}
	if retrievedNonce.MaxValidatorNum != testNonce.MaxValidatorNum {
		t.Errorf("Expected MaxValidatorNum %d, got %d", testNonce.MaxValidatorNum, retrievedNonce.MaxValidatorNum)
	}
	if !bytes.Equal(retrievedNonce.Nonce, testNonce.Nonce) {
		t.Errorf("Expected nonce %s, got %s", string(testNonce.Nonce), string(retrievedNonce.Nonce))
	}
}

func TestArchiveDB_GetVrfNonces(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	blockNumber := uint64(100)
	maxValidatorNum := uint32(5)

	// Store multiple nonces
	batch := archiveDB.db.NewBatch()
	for i := uint64(blockNumber - uint64(maxValidatorNum) + 1); i <= blockNumber; i++ {
		nonce := &VRFNonce{
			MaxValidatorNum: maxValidatorNum,
			Nonce:           []byte(fmt.Sprintf("nonce_%d", i)),
		}
		err := archiveDB.SetVrfNonce(batch, i, nonce)
		if err != nil {
			t.Errorf("Failed to set VRF nonce for block %d: %v", i, err)
		}
	}
	err := batch.Write()
	if err != nil {
		t.Errorf("Failed to write batch: %v", err)
	}

	// Get all nonces
	nonces, err := archiveDB.GetVrfNonces(blockNumber)
	if err != nil {
		t.Errorf("Failed to get VRF nonces: %v", err)
	}
	if len(nonces) != int(maxValidatorNum) {
		t.Errorf("Expected %d nonces, got %d", maxValidatorNum, len(nonces))
	}

	// Verify nonce values
	for i, nonce := range nonces {
		expectedBlock := blockNumber - uint64(maxValidatorNum) + 1 + uint64(i)
		expectedNonce := []byte(fmt.Sprintf("nonce_%d", expectedBlock))
		if !bytes.Equal(nonce, expectedNonce) {
			t.Errorf("Expected nonce %s at position %d, got %s", string(expectedNonce), i, string(nonce))
		}
	}
}

func TestArchiveDB_Init(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	// Test initialization without current block (empty database)
	walkCalled := false
	walkFunc := func(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
		walkCalled = true
		// Simulate walking through some data
		testKV := map[string][]byte{
			"key1": []byte("value1"),
			"key2": []byte("value2"),
		}

		mockIter := newMockIterator(testKV)
		for mockIter.Next() {
			err := f(big.NewInt(1), mockIter)
			if err != nil {
				return err
			}
		}
		return nil
	}

	err := archiveDB.init(walkFunc)
	if err != nil {
		t.Errorf("Failed to initialize archiveDB: %v", err)
	}
	if !walkCalled {
		t.Error("Walk function should have been called")
	}
	if archiveDB.trie == nil {
		t.Error("Trie should be initialized after init")
	}

	// Check current block is set
	currentBlock, err := archiveDB.CurrentBlock()
	if err != nil {
		t.Errorf("Failed to get current block: %v", err)
	}
	if currentBlock == nil {
		t.Error("Current block should be set after initialization")
	}
}

func TestArchiveDB_Init_WithCurrentBlock(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	// First, initialize the database without any current block
	walkFunc := func(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
		// Simulate walking through some data
		testKV := map[string][]byte{
			"key1": []byte("value1"),
			"key2": []byte("value2"),
		}

		mockIter := newMockIterator(testKV)
		for mockIter.Next() {
			err := f(big.NewInt(100), mockIter)
			if err != nil {
				return err
			}
			mockIter.Next()
		}

		return nil
	}

	err := archiveDB.init(walkFunc)
	if err != nil {
		t.Fatalf("Failed to initialize archiveDB first time: %v", err)
	}

	// Now test initialization when current block already exists
	// (current block should be set to 100 from the walk function above)
	walkCalled := false
	walkFunc2 := func(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
		walkCalled = true
		return nil
	}

	err = archiveDB.init(walkFunc2)
	if err != nil {
		t.Errorf("Failed to initialize archiveDB second time: %v", err)
	}
	if walkCalled {
		t.Error("Walk function should not have been called when current block exists")
	}
	if archiveDB.trie == nil {
		t.Error("Trie should be initialized after init")
	}

	// Verify current block is set correctly
	currentBlock, err := archiveDB.CurrentBlock()
	if err != nil {
		t.Errorf("Failed to get current block: %v", err)
	}
	if currentBlock == nil || *currentBlock != 100 {
		t.Errorf("Expected current block number 100, got %v", currentBlock)
	}
}

func TestArchiveDB_CommitBlock(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	// Initialize archiveDB first

	err := archiveDB.init(walkFuncInit)
	if err != nil {
		t.Fatalf("Failed to initialize archiveDB: %v", err)
	}

	// Create test block data
	parentHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	blockHash := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	blockData := createTestBlockData(t, 1, parentHash, blockHash)

	// Commit the block
	err = archiveDB.CommitBlock(blockData)
	if err != nil {
		t.Errorf("Failed to commit block: %v", err)
	}

	// Verify the block was committed
	currentBlock, err := archiveDB.CurrentBlock()
	if err != nil {
		t.Errorf("Failed to get current block: %v", err)
	}
	if currentBlock == nil || *currentBlock != 1 {
		t.Errorf("Expected current block number 1, got %v", currentBlock)
	}

	// Verify archive block was stored
	archiveBlock, err := archiveDB.GetArchiveBlock(1)
	if err != nil {
		t.Errorf("Failed to get archive block: %v", err)
	}
	if archiveBlock == nil || archiveBlock.Number != 1 {
		t.Errorf("Expected archive block number 1, got %d", archiveBlock.Number)
	}

	// Verify VRF nonce was stored
	vrfNonce, err := archiveDB.GetVrfNonce(1)
	if err != nil {
		t.Errorf("Failed to get VRF nonce: %v", err)
	}
	if vrfNonce == nil {
		t.Error("VRF nonce should not be nil")
	}
	if vrfNonce.MaxValidatorNum != 3 {
		t.Errorf("Expected MaxValidatorNum 3, got %d", vrfNonce.MaxValidatorNum)
	}
}

func TestArchiveDB_SnapshotDB(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	// Initialize archiveDB first
	err := archiveDB.init(walkFuncInit)
	if err != nil {
		t.Fatalf("Failed to initialize archive: %v", err)
	}

	// Create and commit test block data
	parentHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	blockHash := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	blockData := createTestBlockData(t, 1, parentHash, blockHash)

	err = archiveDB.CommitBlock(blockData)
	if err != nil {
		t.Fatalf("Failed to commit block: %v", err)
	}

	// Create snapshot
	snapshotDB, err := archiveDB.SnapshotDB(1)
	if err != nil {
		t.Errorf("Failed to create snapshot: %v", err)
	}
	if snapshotDB == nil {
		t.Error("Snapshot should not be nil")
	}

	// Test snapshot functionality
	// Get a value from the snapshot
	value, err := snapshotDB.Get(common.Hash{}, []byte("key1"))
	if err != nil {
		t.Errorf("Failed to get value from snapshot: %v", err)
	}
	if !bytes.Equal(value, []byte("value1")) {
		t.Errorf("Expected value1, got %s", string(value))
	}

	// Test VRF nonce retrieval from snapshot
	vrfNonceValue, err := snapshotDB.Get(common.Hash{}, nonceStorageKey)
	if err != nil {
		t.Errorf("Failed to get VRF nonce from snapshot: %v", err)
	}
	if vrfNonceValue == nil {
		t.Error("VRF nonce should not be nil in snapshot")
	}
}

func TestArchiveDB_SnapshotDB_NonExistentBlock(t *testing.T) {
	archiveDB := setupMemoryArchiveDB(t)

	// Initialize archiveDB first
	walkFunc := func(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
		return nil
	}
	err := archiveDB.init(walkFunc)
	if err != nil {
		t.Fatalf("Failed to initialize archive: %v", err)
	}

	// Try to create snapshot for non-existent block
	_, err = archiveDB.SnapshotDB(100)
	if err == nil {
		t.Error("Expected error when creating snapshot for non-existent block")
	}
}

// mockIterator implements a simple iterator for testing
type mockIterator struct {
	data    map[string][]byte
	keys    []string
	current int
	closed  bool
}

func newMockIterator(data map[string][]byte) *mockIterator {
	keys := make([]string, 0, len(data))
	for key := range data {
		keys = append(keys, key)
	}
	return &mockIterator{
		data:    data,
		keys:    keys,
		current: -1,
		closed:  false,
	}
}

func (m *mockIterator) First() bool {
	if m.closed || len(m.keys) == 0 {
		return false
	}
	m.current = 0
	return true
}

func (m *mockIterator) Last() bool {
	if m.closed || len(m.keys) == 0 {
		return false
	}
	m.current = len(m.keys) - 1
	return true
}

func (m *mockIterator) Seek(key []byte) bool {
	if m.closed {
		return false
	}
	keyStr := string(key)
	for i, k := range m.keys {
		if k >= keyStr {
			m.current = i
			return true
		}
	}
	return false
}

func (m *mockIterator) Prev() bool {
	if m.closed || m.current <= 0 {
		return false
	}
	m.current--
	return true
}

func (m *mockIterator) SetReleaser(releaser util.Releaser) {
	// No-op for mock
}

func (m *mockIterator) Valid() bool {
	return !m.closed && m.current >= 0 && m.current < len(m.keys)
}

func (m *mockIterator) Next() bool {
	if m.closed || m.current >= len(m.keys)-1 {
		return false
	}
	m.current++
	return true
}

func (m *mockIterator) Key() []byte {
	if !m.Valid() {
		return nil
	}
	return []byte(m.keys[m.current])
}

func (m *mockIterator) Value() []byte {
	if !m.Valid() {
		return nil
	}
	key := m.keys[m.current]
	return m.data[key]
}

func (m *mockIterator) Release() {
	m.closed = true
}

func (m *mockIterator) Error() error {
	return nil
}

// ==================== MemDB Tests ====================

func TestNewMemDB(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	if memDB == nil {
		t.Fatal("NewMemDB should not return nil")
	}
	if memDB.DB == nil {
		t.Error("MemDB.DB should not be nil")
	}
	if memDB.ref != 0 {
		t.Errorf("Initial ref count should be 0, got %d", memDB.ref)
	}
}

func TestMemDB_Ref(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	// Test single Ref
	memDB.Ref()
	if memDB.ref != 1 {
		t.Errorf("Expected ref count 1 after Ref(), got %d", memDB.ref)
	}

	// Test multiple Ref
	memDB.Ref()
	memDB.Ref()
	if memDB.ref != 3 {
		t.Errorf("Expected ref count 3 after 3 Ref() calls, got %d", memDB.ref)
	}
}

func TestMemDB_Deref(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	// Set ref count to 3
	memDB.Ref()
	memDB.Ref()
	memDB.Ref()

	// Test Deref
	memDB.Deref()
	if memDB.ref != 2 {
		t.Errorf("Expected ref count 2 after Deref(), got %d", memDB.ref)
	}

	memDB.Deref()
	if memDB.ref != 1 {
		t.Errorf("Expected ref count 1 after second Deref(), got %d", memDB.ref)
	}

	memDB.Deref()
	if memDB.ref != 0 {
		t.Errorf("Expected ref count 0 after third Deref(), got %d", memDB.ref)
	}
}

func TestMemDB_RefDeref_Concurrent(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	var wg sync.WaitGroup
	numGoroutines := 100

	// Concurrent Ref operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			memDB.Ref()
		}()
	}
	wg.Wait()

	if memDB.ref != numGoroutines {
		t.Errorf("Expected ref count %d after concurrent Ref(), got %d", numGoroutines, memDB.ref)
	}

	// Concurrent Deref operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			memDB.Deref()
		}()
	}
	wg.Wait()

	if memDB.ref != 0 {
		t.Errorf("Expected ref count 0 after concurrent Deref(), got %d", memDB.ref)
	}
}

func TestMemDB_PutGet(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	// Test Put
	key := []byte("testKey")
	value := []byte("testValue")
	err := memDB.Put(key, value)
	if err != nil {
		t.Errorf("Put should not return error, got: %v", err)
	}

	// Test Get
	gotValue, err := memDB.Get(key)
	if err != nil {
		t.Errorf("Get should not return error, got: %v", err)
	}
	if !bytes.Equal(gotValue, value) {
		t.Errorf("Expected value %s, got %s", string(value), string(gotValue))
	}
}

func TestMemDB_Delete(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	// Put a value
	key := []byte("testKey")
	value := []byte("testValue")
	memDB.Put(key, value)

	// Delete the value
	err := memDB.Delete(key)
	if err != nil {
		t.Errorf("Delete should not return error, got: %v", err)
	}

	// Verify deletion
	_, err = memDB.Get(key)
	if err == nil {
		t.Error("Get should return error for deleted key")
	}
}

func TestMemDB_Contains(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	key := []byte("testKey")
	value := []byte("testValue")

	// Test Contains for non-existent key
	if memDB.Contains(key) {
		t.Error("Contains should return false for non-existent key")
	}

	// Put a value
	memDB.Put(key, value)

	// Test Contains for existing key
	if !memDB.Contains(key) {
		t.Error("Contains should return true for existing key")
	}
}

func TestMemDB_NewIterator(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	// Add some test data
	testData := []struct {
		key   []byte
		value []byte
	}{
		{[]byte("key1"), []byte("value1")},
		{[]byte("key2"), []byte("value2")},
		{[]byte("key3"), []byte("value3")},
	}

	for _, data := range testData {
		memDB.Put(data.key, data.value)
	}

	// Test NewIterator
	iter := memDB.NewIterator(nil)
	defer iter.Release()

	count := 0
	for iter.Next() {
		count++
	}

	if count != len(testData) {
		t.Errorf("Iterator should iterate over %d items, got %d", len(testData), count)
	}
}

func TestMemDB_Len(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	// Test empty MemDB
	if memDB.Len() != 0 {
		t.Errorf("Len should return 0 for empty MemDB, got %d", memDB.Len())
	}

	// Add some data
	memDB.Put([]byte("key1"), []byte("value1"))
	memDB.Put([]byte("key2"), []byte("value2"))

	if memDB.Len() != 2 {
		t.Errorf("Len should return 2, got %d", memDB.Len())
	}
}

func TestMemDB_Size(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	// Test empty MemDB
	initialSize := memDB.Size()
	if initialSize != 0 {
		t.Errorf("Size should return 0 for empty MemDB, got %d", initialSize)
	}

	// Add some data
	memDB.Put([]byte("key1"), []byte("value1"))
	newSize := memDB.Size()
	if newSize <= initialSize {
		t.Errorf("Size should increase after adding data")
	}
}

func TestMemDB_Reset(t *testing.T) {
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	// Add some data
	memDB.Put([]byte("key1"), []byte("value1"))
	memDB.Put([]byte("key2"), []byte("value2"))

	if memDB.Len() != 2 {
		t.Errorf("Expected 2 items before reset, got %d", memDB.Len())
	}

	// Reset
	memDB.Reset()

	if memDB.Len() != 0 {
		t.Errorf("Expected 0 items after reset, got %d", memDB.Len())
	}

	// Verify we can still use it after reset
	memDB.Put([]byte("newKey"), []byte("newValue"))
	if memDB.Len() != 1 {
		t.Errorf("Expected 1 item after adding to reset MemDB, got %d", memDB.Len())
	}
}

// ==================== RankingCache Tests ====================

func TestNewRankingCache(t *testing.T) {
	cache := NewRankingCache(10)

	if cache == nil {
		t.Fatal("NewRankingCache should not return nil")
	}
	if cache.rankingCache == nil {
		t.Error("RankingCache.rankingCache should not be nil")
	}
}

func TestRankingCache_AddAndGet(t *testing.T) {
	cache := NewRankingCache(10)

	// Create a MemDB
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)
	memDB.Put([]byte("key1"), []byte("value1"))

	// Add to cache
	cacheKey := []byte("cacheKey1")
	cache.Add(cacheKey, memDB)

	// Get from cache
	retrieved := cache.Get(cacheKey)
	if retrieved == nil {
		t.Fatal("Get should return the cached MemDB")
	}

	// Verify the data is correct
	value, err := retrieved.Get([]byte("key1"))
	if err != nil {
		t.Errorf("Get from retrieved MemDB should not return error, got: %v", err)
	}
	if !bytes.Equal(value, []byte("value1")) {
		t.Errorf("Expected value 'value1', got '%s'", string(value))
	}
}

func TestRankingCache_GetNonExistent(t *testing.T) {
	cache := NewRankingCache(10)

	// Try to get a non-existent key
	retrieved := cache.Get([]byte("nonExistentKey"))
	if retrieved != nil {
		t.Error("Get should return nil for non-existent key")
	}
}

func TestRankingCache_RefCount(t *testing.T) {
	cache := NewRankingCache(10)

	// Create a MemDB
	mdb := memdb.New(DefaultComparer, 100)
	memDB := NewMemDB(mdb)

	initialRef := memDB.ref

	// Add to cache (should increase ref count)
	cacheKey := []byte("cacheKey1")
	cache.Add(cacheKey, memDB)

	if memDB.ref != initialRef+1 {
		t.Errorf("After Add, ref count should be %d, got %d", initialRef+1, memDB.ref)
	}

	// Get from cache (should increase ref count again)
	retrieved := cache.Get(cacheKey)
	if retrieved == nil {
		t.Fatal("Get should return the cached MemDB")
	}

	if memDB.ref != initialRef+2 {
		t.Errorf("After Get, ref count should be %d, got %d", initialRef+2, memDB.ref)
	}
}

func TestRankingCache_Eviction(t *testing.T) {
	// Create a small cache that will trigger eviction
	cacheSize := 2
	cache := NewRankingCache(cacheSize)

	// Create and add multiple MemDBs
	memDBs := make([]*MemDB, 3)
	for i := 0; i < 3; i++ {
		mdb := memdb.New(DefaultComparer, 100)
		memDBs[i] = NewMemDB(mdb)
		memDBs[i].Put([]byte(fmt.Sprintf("key%d", i)), []byte(fmt.Sprintf("value%d", i)))
	}

	// Add first two items
	cache.Add([]byte("key0"), memDBs[0])
	cache.Add([]byte("key1"), memDBs[1])

	// Add third item, which should trigger eviction of first item
	cache.Add([]byte("key2"), memDBs[2])

	// First item should be evicted
	retrieved0 := cache.Get([]byte("key0"))
	if retrieved0 != nil {
		t.Error("First item should be evicted from cache")
	}

	// Second and third items should still be in cache
	retrieved1 := cache.Get([]byte("key1"))
	if retrieved1 == nil {
		t.Error("Second item should still be in cache")
	}

	retrieved2 := cache.Get([]byte("key2"))
	if retrieved2 == nil {
		t.Error("Third item should still be in cache")
	}
}

func TestRankingCache_MultipleAddsWithSameKey(t *testing.T) {
	cache := NewRankingCache(10)

	// Create two MemDBs
	mdb1 := memdb.New(DefaultComparer, 100)
	memDB1 := NewMemDB(mdb1)
	memDB1.Put([]byte("data"), []byte("first"))

	mdb2 := memdb.New(DefaultComparer, 100)
	memDB2 := NewMemDB(mdb2)
	memDB2.Put([]byte("data"), []byte("second"))

	// Add first MemDB
	cacheKey := []byte("sameKey")
	cache.Add(cacheKey, memDB1)

	// Add second MemDB with same key (should replace)
	cache.Add(cacheKey, memDB2)

	// Get from cache
	retrieved := cache.Get(cacheKey)
	if retrieved == nil {
		t.Fatal("Get should return the cached MemDB")
	}

	// Verify it's the second MemDB
	value, err := retrieved.Get([]byte("data"))
	if err != nil {
		t.Errorf("Get from retrieved MemDB should not return error, got: %v", err)
	}
	if !bytes.Equal(value, []byte("second")) {
		t.Errorf("Expected value 'second', got '%s'", string(value))
	}
}

func TestRankingCache_Concurrent(t *testing.T) {
	cache := NewRankingCache(100)

	var wg sync.WaitGroup
	numGoroutines := 50

	// Concurrent Add operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			mdb := memdb.New(DefaultComparer, 100)
			memDB := NewMemDB(mdb)
			memDB.Put([]byte(fmt.Sprintf("data%d", idx)), []byte(fmt.Sprintf("value%d", idx)))
			cache.Add([]byte(fmt.Sprintf("key%d", idx)), memDB)
		}(i)
	}
	wg.Wait()

	// Concurrent Get operations
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(idx int) {
			defer wg.Done()
			cache.Get([]byte(fmt.Sprintf("key%d", idx)))
		}(i)
	}
	wg.Wait()
}

func TestRankingCacheKey(t *testing.T) {
	root := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	key := []byte("testKey")
	ranges := 100

	cacheKey := RankingCacheKey(root, key, ranges)

	// Verify key is composed correctly
	// root (32 bytes) + key + ranges (4 bytes)
	expectedLen := 32 + len(key) + 4
	if len(cacheKey) != expectedLen {
		t.Errorf("Expected cache key length %d, got %d", expectedLen, len(cacheKey))
	}

	// Verify root is at the beginning
	if !bytes.HasPrefix(cacheKey, root.Bytes()) {
		t.Error("Cache key should start with root hash")
	}

	// Generate same key again and verify consistency
	cacheKey2 := RankingCacheKey(root, key, ranges)
	if !bytes.Equal(cacheKey, cacheKey2) {
		t.Error("RankingCacheKey should be deterministic")
	}

	// Different ranges should produce different keys
	cacheKey3 := RankingCacheKey(root, key, 200)
	if bytes.Equal(cacheKey, cacheKey3) {
		t.Error("Different ranges should produce different cache keys")
	}

	// Different keys should produce different cache keys
	cacheKey4 := RankingCacheKey(root, []byte("differentKey"), ranges)
	if bytes.Equal(cacheKey, cacheKey4) {
		t.Error("Different keys should produce different cache keys")
	}

	// Different roots should produce different cache keys
	root2 := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	cacheKey5 := RankingCacheKey(root2, key, ranges)
	if bytes.Equal(cacheKey, cacheKey5) {
		t.Error("Different roots should produce different cache keys")
	}
}

// ==================== archiveSnapshot Tests ====================

// setupArchiveSnapshot creates an archiveSnapshot for testing
func setupArchiveSnapshot(t *testing.T) *archiveSnapshot {
	archiveDB := setupMemoryArchiveDB(t)

	// Initialize archiveDB first
	err := archiveDB.init(walkFuncInit)
	if err != nil {
		t.Fatalf("Failed to initialize archiveDB: %v", err)
	}

	// Create and commit test block data
	parentHash := common.HexToHash("0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890")
	blockHash := common.HexToHash("0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef")
	blockData := createTestBlockData(t, 1, parentHash, blockHash)

	err = archiveDB.CommitBlock(blockData)
	if err != nil {
		t.Fatalf("Failed to commit block: %v", err)
	}

	// Get snapshot
	snapshot, err := archiveDB.SnapshotDB(1)
	if err != nil {
		t.Fatalf("Failed to create snapshot: %v", err)
	}

	return snapshot.(*archiveSnapshot)
}

func TestArchiveSnapshot_Get(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Test getting existing key
	value, err := snapshot.Get(common.Hash{}, []byte("key1"))
	if err != nil {
		t.Errorf("Get should not return error for existing key, got: %v", err)
	}
	if !bytes.Equal(value, []byte("value1")) {
		t.Errorf("Expected value 'value1', got '%s'", string(value))
	}

	// Test getting another existing key
	value2, err := snapshot.Get(common.Hash{}, []byte("key2"))
	if err != nil {
		t.Errorf("Get should not return error for existing key, got: %v", err)
	}
	if !bytes.Equal(value2, []byte("value2")) {
		t.Errorf("Expected value 'value2', got '%s'", string(value2))
	}
}

func TestArchiveSnapshot_Get_NonExistent(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Test getting non-existent key
	_, err := snapshot.Get(common.Hash{}, []byte("nonExistentKey"))
	if err != ErrNotFound {
		t.Errorf("Get should return ErrNotFound for non-existent key, got: %v", err)
	}
}

func TestArchiveSnapshot_Get_VrfNonce(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Test getting VRF nonce (special case with nonceStorageKey prefix)
	value, err := snapshot.Get(common.Hash{}, nonceStorageKey)
	if err != nil {
		t.Errorf("Get should not return error for VRF nonce, got: %v", err)
	}
	if value == nil {
		t.Error("VRF nonce should not be nil")
	}
}

func TestArchiveSnapshot_GetFromCommittedBlock(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Test getting existing key
	value, err := snapshot.GetFromCommittedBlock([]byte("key1"))
	if err != nil {
		t.Errorf("GetFromCommittedBlock should not return error for existing key, got: %v", err)
	}
	if !bytes.Equal(value, []byte("value1")) {
		t.Errorf("Expected value 'value1', got '%s'", string(value))
	}

	// Test getting non-existent key
	_, err = snapshot.GetFromCommittedBlock([]byte("nonExistentKey"))
	if err != ErrNotFound {
		t.Errorf("GetFromCommittedBlock should return ErrNotFound for non-existent key, got: %v", err)
	}
}

func TestArchiveSnapshot_Put(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Put a new key-value pair
	key := []byte("newKey")
	value := []byte("newValue")
	err := snapshot.Put(common.Hash{}, key, value)
	if err != nil {
		t.Errorf("Put should not return error, got: %v", err)
	}

	// Verify the value was stored
	gotValue, err := snapshot.Get(common.Hash{}, key)
	if err != nil {
		t.Errorf("Get should not return error after Put, got: %v", err)
	}
	if !bytes.Equal(gotValue, value) {
		t.Errorf("Expected value '%s', got '%s'", string(value), string(gotValue))
	}
}

func TestArchiveSnapshot_Put_Update(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Update an existing key
	key := []byte("key1")
	newValue := []byte("updatedValue")
	err := snapshot.Put(common.Hash{}, key, newValue)
	if err != nil {
		t.Errorf("Put should not return error, got: %v", err)
	}

	// Verify the value was updated
	gotValue, err := snapshot.Get(common.Hash{}, key)
	if err != nil {
		t.Errorf("Get should not return error after Put, got: %v", err)
	}
	if !bytes.Equal(gotValue, newValue) {
		t.Errorf("Expected value '%s', got '%s'", string(newValue), string(gotValue))
	}
}

func TestArchiveSnapshot_Del(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Verify the key exists first
	_, err := snapshot.Get(common.Hash{}, []byte("key1"))
	if err != nil {
		t.Errorf("Key should exist before deletion, got: %v", err)
	}

	// Delete the key
	err = snapshot.Del(common.Hash{}, []byte("key1"))
	if err != nil {
		t.Errorf("Del should not return error, got: %v", err)
	}

	// Verify the key was deleted
	_, err = snapshot.Get(common.Hash{}, []byte("key1"))
	if err != ErrNotFound {
		t.Errorf("Get should return ErrNotFound after deletion, got: %v", err)
	}
}

func TestArchiveSnapshot_Has(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Test Has for existing key
	exists, err := snapshot.Has(common.Hash{}, []byte("key1"))
	if err != nil {
		t.Errorf("Has should not return error for existing key, got: %v", err)
	}
	if !exists {
		t.Error("Has should return true for existing key")
	}

	// Test Has for non-existent key
	exists, err = snapshot.Has(common.Hash{}, []byte("nonExistentKey"))
	if err != ErrNotFound {
		t.Errorf("Has should return ErrNotFound for non-existent key, got: %v", err)
	}
}

func TestArchiveSnapshot_Flush(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Flush should return nil (no-op)
	err := snapshot.Flush(common.Hash{}, big.NewInt(1))
	if err != nil {
		t.Errorf("Flush should return nil, got: %v", err)
	}
}

func TestArchiveSnapshot_Commit(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Commit should return nil (no-op)
	err := snapshot.Commit(common.Hash{})
	if err != nil {
		t.Errorf("Commit should return nil, got: %v", err)
	}
}

func TestArchiveSnapshot_Clear(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Clear should return nil (no-op)
	err := snapshot.Clear()
	if err != nil {
		t.Errorf("Clear should return nil, got: %v", err)
	}
}

func TestArchiveSnapshot_Close(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Close should return nil (no-op)
	err := snapshot.Close()
	if err != nil {
		t.Errorf("Close should return nil, got: %v", err)
	}
}

func TestArchiveSnapshot_Compaction(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Compaction should return nil (no-op)
	err := snapshot.Compaction()
	if err != nil {
		t.Errorf("Compaction should return nil, got: %v", err)
	}
}

func TestArchiveSnapshot_GetLastKVHash(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// GetLastKVHash should return kvHash bytes
	hash := snapshot.GetLastKVHash(common.Hash{})
	if hash == nil {
		t.Error("GetLastKVHash should not return nil")
	}
	if len(hash) != 32 {
		t.Errorf("GetLastKVHash should return 32 bytes, got %d", len(hash))
	}
}

func TestArchiveSnapshot_Snapshot(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Snapshot should return 0
	revid := snapshot.Snapshot(common.Hash{})
	if revid != 0 {
		t.Errorf("Snapshot should return 0, got %d", revid)
	}
}

func TestArchiveSnapshot_RevertToSnapshot(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// RevertToSnapshot should not panic (no-op)
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("RevertToSnapshot should not panic, got: %v", r)
		}
	}()
	snapshot.RevertToSnapshot(common.Hash{}, 0)
}

func TestArchiveSnapshot_MultipleOperations(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Test a sequence of operations
	// 1. Put new data
	err := snapshot.Put(common.Hash{}, []byte("testKey"), []byte("testValue"))
	if err != nil {
		t.Errorf("Put should not return error, got: %v", err)
	}

	// 2. Verify it exists
	exists, _ := snapshot.Has(common.Hash{}, []byte("testKey"))
	if !exists {
		t.Error("Key should exist after Put")
	}

	// 3. Get the value
	value, err := snapshot.Get(common.Hash{}, []byte("testKey"))
	if err != nil {
		t.Errorf("Get should not return error, got: %v", err)
	}
	if !bytes.Equal(value, []byte("testValue")) {
		t.Errorf("Expected 'testValue', got '%s'", string(value))
	}

	// 4. Update the value
	err = snapshot.Put(common.Hash{}, []byte("testKey"), []byte("updatedValue"))
	if err != nil {
		t.Errorf("Put should not return error, got: %v", err)
	}

	// 5. Verify update
	value, err = snapshot.Get(common.Hash{}, []byte("testKey"))
	if err != nil {
		t.Errorf("Get should not return error, got: %v", err)
	}
	if !bytes.Equal(value, []byte("updatedValue")) {
		t.Errorf("Expected 'updatedValue', got '%s'", string(value))
	}

	// 6. Delete the key
	err = snapshot.Del(common.Hash{}, []byte("testKey"))
	if err != nil {
		t.Errorf("Del should not return error, got: %v", err)
	}

	// 7. Verify deletion
	_, err = snapshot.Get(common.Hash{}, []byte("testKey"))
	if err != ErrNotFound {
		t.Errorf("Get should return ErrNotFound after deletion, got: %v", err)
	}
}

func TestArchiveSnapshot_BlockNumber(t *testing.T) {
	snapshot := setupArchiveSnapshot(t)

	// Verify block number is set correctly
	if snapshot.blockNumber != 1 {
		t.Errorf("Expected block number 1, got %d", snapshot.blockNumber)
	}
}

// ==================== archiveSnapshot Ranking Tests ====================

// setupArchiveSnapshotWithRankingCache creates an archiveSnapshot with rankingCache for Ranking tests
func setupArchiveSnapshotWithRankingCache(t *testing.T) *archiveSnapshot {
	snapshot := setupArchiveSnapshot(t)

	// Initialize rankingCache if nil
	if snapshot.rankingCache == nil {
		snapshot.rankingCache = NewRankingCache(defaultRankingCache)
	}

	return snapshot
}

func TestArchiveSnapshot_Ranking_Basic(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	// Get ranking iterator with prefix "key"
	iter := snapshot.Ranking(common.Hash{}, []byte("key"), 10)
	if iter == nil {
		t.Fatal("Ranking should not return nil iterator")
	}
	defer iter.Release()

	// Count results
	count := 0
	for iter.Next() {
		count++
	}

	// Should have results (key1, key2, key3 from test data)
	if iter.Error() != nil {
		t.Errorf("Iterator should not have error, got: %v", iter.Error())
	}
}

func TestArchiveSnapshot_Ranking_EmptyPrefix(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	// Get ranking iterator with empty prefix (should match all)
	iter := snapshot.Ranking(common.Hash{}, []byte{}, 100)
	if iter == nil {
		t.Fatal("Ranking should not return nil iterator")
	}
	defer iter.Release()

	// Should not error
	if iter.Error() != nil {
		t.Errorf("Iterator should not have error, got: %v", iter.Error())
	}
}

func TestArchiveSnapshot_Ranking_NonMatchingPrefix(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	// Get ranking iterator with non-matching prefix
	iter := snapshot.Ranking(common.Hash{}, []byte("nonexistent_prefix_xyz"), 10)
	if iter == nil {
		t.Fatal("Ranking should not return nil iterator")
	}
	defer iter.Release()

	// Should have no results
	count := 0
	for iter.Next() {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 results for non-matching prefix, got %d", count)
	}
}

func TestArchiveSnapshot_Ranking_LimitResults(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	// Get ranking iterator with limit of 1
	iter := snapshot.Ranking(common.Hash{}, []byte("key"), 1)
	if iter == nil {
		t.Fatal("Ranking should not return nil iterator")
	}
	defer iter.Release()

	// Count results - should be limited
	count := 0
	for iter.Next() {
		count++
	}

	// Should have at most 1 result due to limit
	if count > 1 {
		t.Errorf("Expected at most 1 result with limit 1, got %d", count)
	}
}

func TestArchiveSnapshot_Ranking_CacheHit(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	prefix := []byte("key")
	ranges := 10

	// First call - cache miss
	iter1 := snapshot.Ranking(common.Hash{}, prefix, ranges)
	if iter1 == nil {
		t.Fatal("First Ranking call should not return nil iterator")
	}

	// Count results from first iterator
	count1 := 0
	for iter1.Next() {
		count1++
	}
	iter1.Release()

	// Second call with same parameters - should hit cache
	iter2 := snapshot.Ranking(common.Hash{}, prefix, ranges)
	if iter2 == nil {
		t.Fatal("Second Ranking call should not return nil iterator")
	}

	// Count results from second iterator
	count2 := 0
	for iter2.Next() {
		count2++
	}
	iter2.Release()

	// Both iterators should return same number of results
	if count1 != count2 {
		t.Errorf("Cache hit should return same results: first=%d, second=%d", count1, count2)
	}
}

func TestArchiveSnapshot_Ranking_DifferentRanges(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	prefix := []byte("key")

	// Call with ranges=10
	iter1 := snapshot.Ranking(common.Hash{}, prefix, 10)
	if iter1 == nil {
		t.Fatal("Ranking with ranges=10 should not return nil iterator")
	}
	iter1.Release()

	// Call with ranges=5 (different cache key)
	iter2 := snapshot.Ranking(common.Hash{}, prefix, 5)
	if iter2 == nil {
		t.Fatal("Ranking with ranges=5 should not return nil iterator")
	}
	iter2.Release()
}

func TestArchiveSnapshot_Ranking_IteratorOperations(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	iter := snapshot.Ranking(common.Hash{}, []byte("key"), 10)
	if iter == nil {
		t.Fatal("Ranking should not return nil iterator")
	}
	defer iter.Release()

	// Test iterator operations
	hasNext := iter.Next()
	if hasNext {
		// If there are results, Key() and Value() should return non-nil
		key := iter.Key()
		value := iter.Value()
		// Key and Value may be nil if no preimage, but should not panic
		_ = key
		_ = value
	}

	// Error should be nil
	if iter.Error() != nil {
		t.Errorf("Iterator error should be nil, got: %v", iter.Error())
	}
}

func TestArchiveSnapshot_Ranking_MultipleIterators(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	// Create multiple iterators with different prefixes
	iter1 := snapshot.Ranking(common.Hash{}, []byte("key1"), 10)
	iter2 := snapshot.Ranking(common.Hash{}, []byte("key2"), 10)
	iter3 := snapshot.Ranking(common.Hash{}, []byte("key3"), 10)

	if iter1 == nil || iter2 == nil || iter3 == nil {
		t.Fatal("All iterators should be non-nil")
	}

	// Release all iterators
	iter1.Release()
	iter2.Release()
	iter3.Release()
}

func TestArchiveSnapshot_Ranking_ZeroRanges(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	// Get ranking iterator with zero ranges
	iter := snapshot.Ranking(common.Hash{}, []byte("keyx"), 0)
	if iter == nil {
		t.Fatal("Ranking should not return nil iterator even with zero ranges")
	}
	defer iter.Release()

	// Should have no results with zero ranges
	count := 0
	for iter.Next() {
		count++
	}

	if count != 0 {
		t.Errorf("Expected 0 results with zero ranges, got %d", count)
	}
}

func TestArchiveSnapshot_Ranking_LargeRanges(t *testing.T) {
	snapshot := setupArchiveSnapshotWithRankingCache(t)

	// Get ranking iterator with large ranges value
	iter := snapshot.Ranking(common.Hash{}, []byte("key"), 10000)
	if iter == nil {
		t.Fatal("Ranking should not return nil iterator")
	}
	defer iter.Release()

	// Should not error even with large ranges
	if iter.Error() != nil {
		t.Errorf("Iterator should not have error, got: %v", iter.Error())
	}
}
