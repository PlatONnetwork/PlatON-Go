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
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/ethdb/memorydb"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/syndtr/goleveldb/leveldb/iterator"
	"github.com/syndtr/goleveldb/leveldb/util"
)

// setupMemoryArchiveDB creates an archiveDB with memory database for testing
func setupMemoryArchiveDB(t *testing.T) *archiveDB {
	memDB := memorydb.New()
	rawDB := rawdb.NewDatabase(memDB)

	archiveDB, err := NewArchiveDBWithDB(rawDB)
	if err != nil {
		t.Fatalf("Failed to create archiveDB with memory DB: %v", err)
	}

	return archiveDB
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
		t.Fatalf("Failed to create archiveDB with memory DB: %v", err)
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
	walkFunc := func(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
		return nil
	}
	err := archiveDB.init(walkFunc)
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
	walkFunc := func(slice *util.Range, f func(num *big.Int, iter iterator.Iterator) error) error {
		return nil
	}
	err := archiveDB.init(walkFunc)
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
		t.Fatalf("Failed to initialize archiveDB: %v", err)
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
