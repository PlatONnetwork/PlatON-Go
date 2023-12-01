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

// Package state provides a caching layer atop the Ethereum state trie.
package state

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/state/snapshot"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
)

type revision struct {
	id           int
	journalIndex int
}

var (
	// emptyRoot is the known root hash of an empty trie.
	emptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
)

// StateDB structs within the ethereum protocol are used to store anything
// within the merkle trie. StateDBs take care of caching and storing
// nested states. It's the general query interface to retrieve:
// * Contracts
// * Accounts
type StateDB struct {
	db   Database
	trie Trie

	snaps         *snapshot.Tree
	snap          snapshot.Snapshot
	snapDestructs map[common.Hash]struct{}
	snapAccounts  map[common.Hash][]byte
	snapStorage   map[common.Hash]map[common.Hash][]byte

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateObjects        map[common.Address]*stateObject
	stateObjectsPending map[common.Address]struct{} // State objects finalized but not yet written to the trie
	stateObjectsDirty   map[common.Address]struct{} // State objects modified in the current execution

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// The refund counter, also used by state transitioning.
	refund uint64

	thash   common.Hash
	txIndex int
	logs    map[common.Hash][]*types.Log
	logSize uint

	preimages map[common.Hash][]byte

	// Per-transaction access list
	accessList *accessList

	// Journal of state modifications. This is the backbone of
	// Snapshot and RevertToSnapshot.
	journal        *journal
	validRevisions []revision
	nextRevisionId int

	lock sync.Mutex

	// Prevent concurrent access to parent StateDB and reference function
	refLock sync.Mutex

	// The flag of parent StateDB
	parentCommitted bool
	// children StateDB callback, is called when parent committed
	clearReferenceFunc []func()
	parent             *StateDB

	// The index in clearReferenceFunc of parent StateDB
	referenceFuncIndex int
	// statedb is created based on this root
	originRoot common.Hash

	// Measurements gathered during execution for debugging purposes
	AccountReads         time.Duration
	AccountHashes        time.Duration
	AccountUpdates       time.Duration
	AccountCommits       time.Duration
	StorageReads         time.Duration
	StorageHashes        time.Duration
	StorageUpdates       time.Duration
	StorageCommits       time.Duration
	SnapshotAccountReads time.Duration
	SnapshotStorageReads time.Duration
	SnapshotCommits      time.Duration

	AccountUpdated int
	StorageUpdated int
	AccountDeleted int
	StorageDeleted int
}

// New creates a new state from a given trie.
func New(root common.Hash, db Database, snaps *snapshot.Tree) (*StateDB, error) {
	tr, err := db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	sdb := &StateDB{
		db:                  db,
		trie:                tr,
		snaps:               snaps,
		stateObjects:        make(map[common.Address]*stateObject),
		stateObjectsPending: make(map[common.Address]struct{}),
		stateObjectsDirty:   make(map[common.Address]struct{}),
		logs:                make(map[common.Hash][]*types.Log),
		preimages:           make(map[common.Hash][]byte),
		journal:             newJournal(),
		accessList:          newAccessList(),
		clearReferenceFunc:  make([]func(), 0),
		originRoot:          root,
	}
	if sdb.snaps != nil {
		if sdb.snap = sdb.snaps.Snapshot(root); sdb.snap != nil {
			sdb.snapDestructs = make(map[common.Hash]struct{})
			sdb.snapAccounts = make(map[common.Hash][]byte)
			sdb.snapStorage = make(map[common.Hash]map[common.Hash][]byte)
		}
	}
	return sdb, nil
}

// New StateDB based on the parent StateDB
func (s *StateDB) NewStateDB() *StateDB {
	stateDB := &StateDB{
		db:                  s.db,
		trie:                s.db.NewTrie(s.trie),
		snaps:               s.snaps,
		stateObjects:        make(map[common.Address]*stateObject),
		stateObjectsPending: make(map[common.Address]struct{}),
		stateObjectsDirty:   make(map[common.Address]struct{}),
		logs:                make(map[common.Hash][]*types.Log),
		preimages:           make(map[common.Hash][]byte),
		journal:             newJournal(),
		parent:              s,
		accessList:          newAccessList(),
		clearReferenceFunc:  make([]func(), 0),
		originRoot:          s.Root(),
	}

	index := s.AddReferenceFunc(stateDB.clearParentRef)
	stateDB.referenceFuncIndex = index

	if stateDB.snaps != nil {
		if stateDB.snap = stateDB.snaps.Snapshot(s.Root()); stateDB.snap != nil {
			stateDB.snapDestructs = make(map[common.Hash]struct{})
			stateDB.snapAccounts = make(map[common.Hash][]byte)
			stateDB.snapStorage = make(map[common.Hash]map[common.Hash][]byte)
		} else {
			log.Error("Failed to find snapshot of parentRoot", "parentRoot", s.Root())
		}
	}
	return stateDB
}

// TxIndex returns the current transaction index set by Prepare.
func (s *StateDB) TxIndex() int {
	return s.txIndex
}

func (s *StateDB) HadParent() bool {
	s.refLock.Lock()
	defer s.refLock.Unlock()
	return s.parent != nil
}

// setError remembers the first non-nil error it is called with.
func (s *StateDB) setError(err error) {
	if s.dbErr == nil {
		s.dbErr = err
	}
}

func (s *StateDB) Error() error {
	return s.dbErr
}

// Reset clears out all ephemeral state objects from the state db, but keeps
// the underlying state trie to avoid reloading data for the next operations.
func (s *StateDB) Reset(root common.Hash) error {
	tr, err := s.db.OpenTrie(root)
	if err != nil {
		return err
	}

	s.trie = tr
	s.stateObjects = make(map[common.Address]*stateObject)
	s.stateObjectsPending = make(map[common.Address]struct{})
	s.stateObjectsDirty = make(map[common.Address]struct{})
	s.thash = common.Hash{}
	s.txIndex = 0
	s.logs = make(map[common.Hash][]*types.Log)
	s.logSize = 0
	s.preimages = make(map[common.Hash][]byte)
	s.clearJournalAndRefund()

	if s.snaps != nil {
		s.snapAccounts, s.snapDestructs, s.snapStorage = nil, nil, nil
		if s.snap = s.snaps.Snapshot(root); s.snap != nil {
			s.snapDestructs = make(map[common.Hash]struct{})
			s.snapAccounts = make(map[common.Hash][]byte)
			s.snapStorage = make(map[common.Hash]map[common.Hash][]byte)
		}
	}
	s.accessList = newAccessList()
	return nil
}

func (s *StateDB) AddLog(logInfo *types.Log) {
	s.journal.append(addLogChange{txhash: s.thash})

	logInfo.TxHash = s.thash
	logInfo.TxIndex = uint(s.txIndex)
	logInfo.Index = s.logSize
	s.logs[s.thash] = append(s.logs[s.thash], logInfo)
	s.logSize++
}

func (s *StateDB) GetLogs(hash common.Hash, blockHash common.Hash) []*types.Log {
	logs := s.logs[hash]
	for _, l := range logs {
		l.BlockHash = blockHash
	}
	return logs
}

func (s *StateDB) Logs() []*types.Log {
	var logs []*types.Log
	for _, lgs := range s.logs {
		logs = append(logs, lgs...)
	}
	return logs
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (s *StateDB) AddPreimage(hash common.Hash, preimage []byte) {
	if _, ok := s.preimages[hash]; !ok {
		s.journal.append(addPreimageChange{hash: hash})
		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		s.preimages[hash] = pi
	}
}

// Preimages returns a list of SHA3 preimages that have been submitted.
func (s *StateDB) Preimages() map[common.Hash][]byte {
	return s.preimages
}

// AddRefund adds gas to the refund counter
func (s *StateDB) AddRefund(gas uint64) {
	s.journal.append(refundChange{prev: s.refund})
	s.refund += gas
}

// SubRefund removes gas from the refund counter.
// This method will panic if the refund counter goes below zero
func (s *StateDB) SubRefund(gas uint64) {
	s.journal.append(refundChange{prev: s.refund})
	if gas > s.refund {
		panic(fmt.Sprintf("Refund counter below zero (gas: %d > refund: %d)", gas, s.refund))
	}
	s.refund -= gas
}

// Exist reports whether the given account address exists in the state.
// Notably this also returns true for suicided accounts.
func (s *StateDB) Exist(addr common.Address) bool {
	return s.getStateObject(addr) != nil
}

// Empty returns whether the state object is either non-existent
// or empty according to the EIP161 specification (balance = nonce = code = 0)
func (s *StateDB) Empty(addr common.Address) bool {
	so := s.getStateObject(addr)
	return so == nil || so.empty()
}

// GetBalance retrieves the balance from the given address or 0 if object not found
func (s *StateDB) GetBalance(addr common.Address) *big.Int {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Balance()
	}
	return common.Big0
}

func (s *StateDB) GetNonce(addr common.Address) uint64 {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Nonce()
	}

	return 0
}

func (s *StateDB) GetCode(addr common.Address) []byte {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Code(s.db)
	}
	return nil
}

func (s *StateDB) GetCodeSize(addr common.Address) int {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.CodeSize(s.db)
	}
	return 0
}

func (s *StateDB) GetCodeHash(addr common.Address) common.Hash {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return common.Hash{}
	}
	return common.BytesToHash(stateObject.CodeHash())
}

// GetState retrieves a value from the given account's storage trie.
func (s *StateDB) GetState(addr common.Address, key []byte) []byte {
	s.lock.Lock()
	defer s.lock.Unlock()
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.removePrefixValue(stateObject.GetState(s.db, key))
	}
	return []byte{}
}

type proofList [][]byte

func (n *proofList) Put(key []byte, value []byte) error {
	*n = append(*n, value)
	return nil
}

func (n *proofList) Delete(key []byte) error {
	panic("not supported")
}

// GetProof returns the Merkle proof for a given account.
func (s *StateDB) GetProof(addr common.Address) ([][]byte, error) {
	return s.GetProofByHash(crypto.Keccak256Hash(addr.Bytes()))
}

// GetProofByHash returns the Merkle proof for a given account.
func (s *StateDB) GetProofByHash(addrHash common.Hash) ([][]byte, error) {
	var proof proofList
	err := s.trie.Prove(addrHash[:], 0, &proof)
	return [][]byte(proof), err
}

// GetStorageProof returns the Merkle proof for given storage slot.
func (s *StateDB) GetStorageProof(a common.Address, key common.Hash) ([][]byte, error) {
	var proof proofList
	trie := s.StorageTrie(a)
	if trie == nil {
		return proof, errors.New("storage trie for requested address does not exist")
	}
	err := trie.Prove(crypto.Keccak256(key.Bytes()), 0, &proof)
	return [][]byte(proof), err
}

// GetCommittedState retrieves a value from the given account's committed storage trie.
func (s *StateDB) GetCommittedState(addr common.Address, key []byte) []byte {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.removePrefixValue(stateObject.GetCommittedState(s.db, key))
	}
	return []byte{}
}

// Database retrieves the low level database supporting the lower level trie ops.
func (s *StateDB) Database() Database {
	return s.db
}

// StorageTrie returns the storage trie of an account.
// The return value is a copy and is nil for non-existent accounts.
func (s *StateDB) StorageTrie(addr common.Address) Trie {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return nil
	}
	cpy := stateObject.deepCopy(s)
	return cpy.updateTrie(s.db)
}

func (s *StateDB) HasSuicided(addr common.Address) bool {
	stateObject := s.getStateObject(addr)
	if stateObject != nil {
		return stateObject.suicided
	}
	return false
}

/*
 * SETTERS
 */

// AddBalance adds amount to the account associated with addr.
func (s *StateDB) AddBalance(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddBalance(amount)
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (s *StateDB) SubBalance(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubBalance(amount)
	}
}

func (s *StateDB) SetBalance(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetBalance(amount)
	}
}

func (s *StateDB) SetNonce(addr common.Address, nonce uint64) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetNonce(nonce)
	}
}

func (s *StateDB) SetCode(addr common.Address, code []byte) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetCode(crypto.Keccak256Hash(code), code)
	}
}

func (s *StateDB) SetState(address common.Address, key, value []byte) {
	s.lock.Lock()
	stateObject := s.GetOrNewStateObject(address)

	if stateObject != nil {
		//stateObject.SetState(self.db, key, stateObject.getPrefixValue(key, value))
		stateObject.SetState(s.db, key, stateObject.getPrefixValue(s.originRoot.Bytes(), key, value))
	}
	s.lock.Unlock()
}

// SetStorage replaces the entire storage for the specified account with given
// storage. This function should only be used for debugging.
func (s *StateDB) SetStorage(addr common.Address, storage map[common.Hash]common.Hash) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetStorage(storage)
	}
}

//func getKeyValue(address common.Address, key []byte, value []byte) (string, common.Hash, []byte) {
//	var buffer bytes.Buffer
//	//buffer.Write(address[:])
//	buffer.Write(key)
//	keyTrie := buffer.String()
//
//	//if value != nil && !bytes.Equal(value,[]byte{}){
//	buffer.Reset()
//	buffer.Write(value)
//
//	valueKey := common.Hash{}
//	keccak := sha3.NewLegacyKeccak256()
//	keccak.Write(buffer.Bytes())
//	keccak.Sum(valueKey[:0])
//
//	return keyTrie, valueKey, value
//}

// Suicide marks the given account as suicided.
// This clears the account balance.
//
// The account's state object is still available until the state is committed,
// getStateObject will return a non-nil account after Suicide.
func (s *StateDB) Suicide(addr common.Address) bool {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		return false
	}
	s.journal.append(suicideChange{
		account:     &addr,
		prev:        stateObject.suicided,
		prevbalance: new(big.Int).Set(stateObject.Balance()),
	})
	stateObject.markSuicided()
	stateObject.data.Balance = new(big.Int)

	return true
}

//
// Setting, updating & deleting state object methods.
//

// updateStateObject writes the given object to the trie.
func (s *StateDB) updateStateObject(obj *stateObject) {
	// Track the amount of time wasted on updating the account from the trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.AccountUpdates += time.Since(start) }(time.Now())
	}
	addr := obj.Address()
	if err := s.trie.TryUpdateAccount(addr[:], &obj.data); err != nil {
		s.setError(fmt.Errorf("updateStateObject (%x) error: %v", addr[:], err))
	}

	// If state snapshotting is active, cache the data til commit. Note, this
	// update mechanism is not symmetric to the deletion, because whereas it is
	// enough to track account updates at commit time, deletions need tracking
	// at transaction boundary level to ensure we capture state clearing.
	if s.snap != nil {
		s.snapAccounts[obj.addrHash] = snapshot.SlimAccountRLP(obj.data.Nonce, obj.data.Balance, obj.data.Root, obj.data.CodeHash, obj.data.StorageKeyPrefix)
	}
}

// deleteStateObject removes the given object from the state trie.
func (s *StateDB) deleteStateObject(obj *stateObject) {
	// Track the amount of time wasted on deleting the account from the trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.AccountUpdates += time.Since(start) }(time.Now())
	}

	// Delete the account from the trie
	addr := obj.Address()
	if err := s.trie.TryDeleteAccount(addr[:]); err != nil {
		s.setError(fmt.Errorf("deleteStateObject (%x) error: %v", addr[:], err))
	}
}

// Get the current StateDB cache and the parent StateDB cache
func (s *StateDB) getStateObjectCache(addr common.Address) (stateObject *stateObject) {
	// Prefer 'live' objects.
	if obj := s.stateObjects[addr]; obj != nil {
		return obj
	}
	s.refLock.Lock()
	parentDB := s.parent
	parentCommitted := s.parentCommitted
	refLock := &s.refLock

	for parentDB != nil {
		obj := parentDB.getStateObjectLocalCache(addr)
		if obj != nil {
			refLock.Unlock()
			cpy := obj.copy(s)
			s.setStateObject(cpy)
			return cpy
		} else if parentCommitted {
			refLock.Unlock()
			//if len(parentDB.clearReferenceFunc) > 0 {
			//	panic(fmt.Sprintf("had parentCommitted statedb clearref is not empty:%d, root:%s", len(parentDB.clearReferenceFunc), parentDB.Root().String()))
			//}
			//if len(self.clearReferenceFunc) > 0 {
			//	panic(fmt.Sprintf("executing statedb clearref is not empty:%d, root:%s", len(self.clearReferenceFunc), parentDB.Root().String()))
			//}
			//if parentDB.parent != nil {
			//	panic(fmt.Sprintf("parent is not nil"))
			//}
			//if parentDB.disk != 3 {
			//	panic(fmt.Sprintf("disk change parent error"))
			//}
			//obj := parentDB.getStateObject(addr)
			//if obj != nil {
			//	cpy := obj.copy(self)
			//	self.setStateObject(cpy)
			//	return cpy
			//}
			return nil
		}

		if obj == nil {
			refLock.Unlock()
			parentDB.refLock.Lock()
			refLock = &parentDB.refLock
			if parentDB.parent == nil {
				break
			}
			parentCommitted = parentDB.parentCommitted
			parentDB = parentDB.parent
		}
	}

	refLock.Unlock()
	return nil
}

// Find stateObject in cache
func (s *StateDB) getStateObjectLocalCache(addr common.Address) (stateObject *stateObject) {
	if obj := s.stateObjects[addr]; obj != nil {
		if obj.deleted {
			return nil
		}
		return obj
	}
	return nil
}

// Find stateObject storage in cache
func (s *StateDB) getStateObjectSnapshot(addr common.Address, key []byte) []byte {
	if obj := s.stateObjects[addr]; obj != nil {
		if obj.deleted {
			return nil
		}
		value, dirty := obj.dirtyStorage[string(key)]
		if dirty {
			return value
		}

		value, pending := obj.pendingStorage[string(key)]
		if pending {
			return value
		}

		value, cached := obj.originStorage[string(key)]
		if cached {
			return value
		}
	}
	return nil
}

// Add childrent statedb reference
func (s *StateDB) AddReferenceFunc(fn func()) int {
	s.refLock.Lock()
	defer s.refLock.Unlock()
	// It must not be nil
	if s.clearReferenceFunc == nil {
		panic("statedb had cleared")
	}
	s.clearReferenceFunc = append(s.clearReferenceFunc, fn)
	return len(s.clearReferenceFunc) - 1
}

// Clear reference when StateDB is committed
func (s *StateDB) ClearReference() {
	s.refLock.Lock()
	defer s.refLock.Unlock()
	for _, fn := range s.clearReferenceFunc {
		if nil != fn {
			fn()
		}
	}
	log.Trace("clear all ref", "reflen", len(s.clearReferenceFunc))
	if s.parent != nil {
		if len(s.parent.clearReferenceFunc) > 0 {
			panic("parent ref > 0")
		}
	}
	s.clearReferenceFunc = nil
	s.parent = nil
}

// Clear reference by index
func (s *StateDB) ClearIndexReference(index int) {
	s.refLock.Lock()
	defer s.refLock.Unlock()

	if len(s.clearReferenceFunc) > index && s.clearReferenceFunc[index] != nil {
		//fn := self.clearReferenceFunc[index]
		//fn()
		log.Trace("Before clear index ref", "reflen", len(s.clearReferenceFunc), "index", index)
		//self.clearReferenceFunc = append(self.clearReferenceFunc[:index], self.clearReferenceFunc[index+1:]...)
		s.clearReferenceFunc[index] = nil
		log.Trace("After clear index ref", "reflen", len(s.clearReferenceFunc), "index", index)
	}
}

// Clear Parent reference
func (s *StateDB) ClearParentReference() {
	s.refLock.Lock()
	defer s.refLock.Unlock()

	if s.parent != nil && s.referenceFuncIndex >= 0 {
		s.parent.ClearIndexReference(s.referenceFuncIndex)
		s.parent = nil
		s.referenceFuncIndex = -1
	}
}

// getStateObject retrieves a state object given by the address, returning nil if
// the object is not found or was deleted in this execution context. If you need
// to differentiate between non-existent/just-deleted, use getDeletedStateObject.
func (s *StateDB) getStateObject(addr common.Address) *stateObject {
	if obj := s.getDeletedStateObject(addr); obj != nil && !obj.deleted {
		return obj
	}
	return nil
}

// getDeletedStateObject is similar to getStateObject, but instead of returning
// nil for a deleted state object, it returns the actual object with the deleted
// flag set. This is needed by the state journal to revert to the correct self-
// destructed object instead of wiping all knowledge about the state object.
func (s *StateDB) getDeletedStateObject(addr common.Address) *stateObject {
	// Prefer live objects if any is available
	if obj := s.getStateObjectCache(addr); obj != nil {
		return obj
	}
	// If no live objects are available, attempt to use snapshots
	var data *types.StateAccount
	if s.snap != nil {
		start := time.Now()
		acc, err := s.snap.Account(crypto.Keccak256Hash(addr[:]))
		if metrics.EnabledExpensive {
			s.SnapshotAccountReads += time.Since(start)
		}
		if err == nil {
			if acc == nil {
				return nil
			}
			data = &types.StateAccount{
				Nonce:            acc.Nonce,
				Balance:          acc.Balance,
				CodeHash:         acc.CodeHash,
				Root:             common.BytesToHash(acc.Root),
				StorageKeyPrefix: acc.StorageKeyPrefix,
			}
			if len(data.CodeHash) == 0 {
				data.CodeHash = emptyCodeHash
			}
			if data.Root == (common.Hash{}) {
				data.Root = emptyRoot
			}
		}
	}
	// If snapshot unavailable or reading from it failed, load from the database
	if data == nil {
		start := time.Now()
		enc, err := s.trie.TryGet(addr[:])
		if metrics.EnabledExpensive {
			s.AccountReads += time.Since(start)
		}
		if err != nil {
			s.setError(fmt.Errorf("getDeleteStateObject (%x) error: %v", addr.Bytes(), err))
			return nil
		}
		if len(enc) == 0 {
			return nil
		}
		data = new(types.StateAccount)
		if err := rlp.DecodeBytes(enc, data); err != nil {
			log.Error("Failed to decode state object", "addr", addr, "err", err)
			return nil
		}
	}
	// [NOTE]: set the prefix for storage key
	if data.Empty() {
		data.StorageKeyPrefix = addr.Bytes()
	}
	// Insert into the live set.
	obj := newObject(s, addr, *data)
	s.setStateObject(obj)
	return obj
}

func (s *StateDB) setStateObject(object *stateObject) {
	if len(s.clearReferenceFunc) > 0 {
		panic("statedb readonly")
	}
	s.stateObjects[object.Address()] = object
}

// GetOrNewStateObject retrieves a state object or create a new state object if nil.
func (s *StateDB) GetOrNewStateObject(addr common.Address) *stateObject {
	stateObject := s.getStateObject(addr)
	if stateObject == nil {
		stateObject, _ = s.createObject(addr)
	}
	return stateObject
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (s *StateDB) createObject(addr common.Address) (newobj, prev *stateObject) {
	prev = s.getDeletedStateObject(addr) // Note, prev might have been deleted, we need that!
	var prevdestruct bool
	if s.snap != nil && prev != nil {
		_, prevdestruct = s.snapDestructs[prev.addrHash]
		if !prevdestruct {
			s.snapDestructs[prev.addrHash] = struct{}{}
		}
	}
	if prev == nil {
		newobj = newObject(s, addr, types.StateAccount{StorageKeyPrefix: addr.Bytes()})
		s.journal.append(createObjectChange{account: &addr})
	} else {
		prefix := make([]byte, len(prev.data.StorageKeyPrefix))
		copy(prefix, prev.data.StorageKeyPrefix)
		newobj = newObject(s, addr, types.StateAccount{StorageKeyPrefix: prefix})
		s.journal.append(resetObjectChange{prev: prev, prevdestruct: prevdestruct})
	}
	s.setStateObject(newobj)
	if prev != nil && !prev.deleted {
		return newobj, prev
	}
	return newobj, nil
}

// CreateAccount explicitly creates a state object. If a state object with the address
// already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might arise that
// a contract does the following:
//
//  1. sends funds to sha(account ++ (nonce + 1))
//  2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (s *StateDB) CreateAccount(addr common.Address) {
	newObj, prev := s.createObject(addr)
	if prev != nil {
		newObj.setBalance(prev.data.Balance)
	}
}

func (s *StateDB) TxHash() common.Hash {
	return s.thash
}

func (s *StateDB) TxIdx() uint32 {
	return uint32(s.txIndex)
}

func (db *StateDB) ForEachStorage(addr common.Address, cb func(key, value []byte) bool) {
	so := db.getStateObject(addr)
	if so == nil {
		return
	}

	it := trie.NewIterator(so.getTrie(db.db).NodeIterator(nil))
	for it.Next() {
		key := db.trie.GetKey(it.Key)
		if value, ok := so.dirtyStorage[string(key)]; ok {
			cb(key, value)
			continue
		}

		cb(key, it.Value)
	}
}

func (db *StateDB) MigrateStorage(from, to common.Address) {
	fromObj := db.getStateObject(from)
	toObj := db.getStateObject(to)
	if nil != fromObj && nil != toObj {
		// replace storage key prefix
		toObj.data.StorageKeyPrefix = make([]byte, len(fromObj.data.StorageKeyPrefix))
		copy(toObj.data.StorageKeyPrefix, fromObj.data.StorageKeyPrefix)
		// replace storageRootHash
		toObj.data.Root = fromObj.data.Root
		// replace storageTrie
		if nil != fromObj.trie {
			toObj.trie = db.db.CopyTrie(fromObj.trie)
		}
		// replace storage
		toObj.dirtyStorage = fromObj.dirtyStorage.Copy()
		toObj.originStorage = fromObj.originStorage.Copy()
	}
}

// Copy creates a deep, independent copy of the state.
// Snapshots of the copied state cannot be applied to the copy.
func (s *StateDB) Copy() *StateDB {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Copy all the basic fields, initialize the memory ones
	state := &StateDB{
		db:                  s.db,
		trie:                s.db.CopyTrie(s.trie),
		stateObjects:        make(map[common.Address]*stateObject, len(s.journal.dirties)),
		stateObjectsPending: make(map[common.Address]struct{}, len(s.stateObjectsPending)),
		stateObjectsDirty:   make(map[common.Address]struct{}, len(s.journal.dirties)),
		refund:              s.refund,
		logs:                make(map[common.Hash][]*types.Log, len(s.logs)),
		logSize:             s.logSize,
		preimages:           make(map[common.Hash][]byte, len(s.preimages)),
		journal:             newJournal(),
		clearReferenceFunc:  make([]func(), 0),
		originRoot:          s.originRoot,
	}

	// Copy the dirty states, logs, and preimages
	for addr := range s.journal.dirties {
		// As documented [here](https://github.com/ethereum/go-ethereum/pull/16485#issuecomment-380438527),
		// and in the Finalise-method, there is a case where an object is in the journal but not
		// in the stateObjects: OOG after touch on ripeMD prior to Byzantium. Thus, we need to check for
		// nil
		if object, exist := s.stateObjects[addr]; exist {
			// Even though the original object is dirty, we are not copying the journal,
			// so we need to make sure that any side-effect the journal would have caused
			// during a commit (or similar op) is already applied to the copy.
			state.stateObjects[addr] = object.deepCopy(state)

			state.stateObjectsDirty[addr] = struct{}{}   // Mark the copy dirty to force internal (code/state) commits
			state.stateObjectsPending[addr] = struct{}{} // Mark the copy pending to force external (account) commits
		}
	}
	// Above, we don't copy the actual journal. This means that if the copy
	// is copied, the loop above will be a no-op, since the copy's journal
	// is empty. Thus, here we iterate over stateObjects, to enable copies
	// of copies.
	for addr := range s.stateObjectsPending {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = s.stateObjects[addr].deepCopy(state)
		}
		state.stateObjectsPending[addr] = struct{}{}
	}
	for addr := range s.stateObjectsDirty {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = s.stateObjects[addr].deepCopy(state)
		}
		state.stateObjectsDirty[addr] = struct{}{}
	}
	for hash, logs := range s.logs {
		cpy := make([]*types.Log, len(logs))
		for i, l := range logs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		state.logs[hash] = cpy
	}
	for hash, preimage := range s.preimages {
		state.preimages[hash] = preimage
	}

	// Do we need to copy the access list? In practice: No. At the start of a
	// transaction, the access list is empty. In practice, we only ever copy state
	// _between_ transactions/blocks, never in the middle of a transaction.
	// However, it doesn't cost us much to copy an empty list, so we do it anyway
	// to not blow up if we ever decide copy it in the middle of a transaction
	state.accessList = s.accessList.Copy()

	// Copy parent state
	s.refLock.Lock()
	if s.parent != nil {
		if !s.parentCommitted {
			state.parent = s.parent
			state.parent.AddReferenceFunc(state.clearParentRef)
		} else {
			s.parent = nil
		}
	}
	state.parentCommitted = s.parentCommitted
	s.refLock.Unlock()

	if s.snaps != nil {
		// In order for the miner to be able to use and make additions
		// to the snapshot tree, we need to copy that aswell.
		// Otherwise, any block mined by ourselves will cause gaps in the tree,
		// and force the miner to operate trie-backed only
		state.snaps = s.snaps
		state.snap = s.snap

		// deep copy needed
		state.snapDestructs = make(map[common.Hash]struct{})
		for k, v := range s.snapDestructs {
			state.snapDestructs[k] = v
		}
		state.snapAccounts = make(map[common.Hash][]byte)
		for k, v := range s.snapAccounts {
			state.snapAccounts[k] = v
		}
		state.snapStorage = make(map[common.Hash]map[common.Hash][]byte)
		for k, v := range s.snapStorage {
			temp := make(map[common.Hash][]byte)
			for kk, vv := range v {
				temp[kk] = vv
			}
			state.snapStorage[k] = temp
		}
	}

	return state
}

// Clear parent StateDB reference
func (s *StateDB) clearParentRef() {
	s.refLock.Lock()
	defer s.refLock.Unlock()

	if s.parent != nil {
		s.parentCommitted = true
		log.Trace("clearParentRef", "parent root", s.parent.Root().String())
		// Parent is nil, find the parent state based on current StateDB
		s.parent = nil
	}
}

// Snapshot returns an identifier for the current revision of the state.
func (s *StateDB) Snapshot() int {
	id := s.nextRevisionId
	s.nextRevisionId++
	s.validRevisions = append(s.validRevisions, revision{id, s.journal.length()})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (s *StateDB) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(s.validRevisions), func(i int) bool {
		return s.validRevisions[i].id >= revid
	})
	if idx == len(s.validRevisions) || s.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := s.validRevisions[idx].journalIndex

	// Replay the journal to undo changes and remove invalidated snapshots
	s.journal.revert(s, snapshot)
	s.validRevisions = s.validRevisions[:idx]
}

// GetRefund returns the current value of the refund counter.
func (s *StateDB) GetRefund() uint64 {
	return s.refund
}

// Finalise finalises the state by removing the self destructed objects and clears
// the journal as well as the refunds. Finalise, however, will not push any updates
// into the tries just yet. Only IntermediateRoot or Commit will do that.
func (s *StateDB) Finalise(deleteEmptyObjects bool) {
	for addr := range s.journal.dirties {
		obj, exist := s.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `s.journal.dirties` but not in `s.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}
		if obj.suicided || (deleteEmptyObjects && obj.empty()) {
			obj.deleted = true

			// If state snapshotting is active, also mark the destruction there.
			// Note, we can't do this only at the end of a block because multiple
			// transactions within the same block might self destruct and then
			// ressurrect an account; but the snapshotter needs both events.
			if s.snap != nil {
				s.snapDestructs[obj.addrHash] = struct{}{} // We need to maintain account deletions explicitly (will remain set indefinitely)
				delete(s.snapAccounts, obj.addrHash)       // Clear out any previously updated account data (may be recreated via a ressurrect)
				delete(s.snapStorage, obj.addrHash)        // Clear out any previously updated storage data (may be recreated via a ressurrect)
			}
		} else {
			obj.finalise()
		}
		s.stateObjectsPending[addr] = struct{}{}
		s.stateObjectsDirty[addr] = struct{}{}
	}
	// Invalidate journal because reverting across transactions is not allowed.
	s.clearJournalAndRefund()
}

// IntermediateRoot computes the current root hash of the state trie.
// It is called in between transactions to get the root hash that
// goes into transaction receipts.
func (s *StateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	// Finalise all the dirty storage states and write them into the tries
	s.Finalise(deleteEmptyObjects)

	for addr := range s.stateObjectsPending {
		obj := s.stateObjects[addr]
		if obj.deleted {
			s.deleteStateObject(obj)
			s.AccountDeleted += 1
		} else {
			obj.updateRoot(s.db)
			s.updateStateObject(obj)
			s.AccountUpdated += 1
		}
	}
	if len(s.stateObjectsPending) > 0 {
		s.stateObjectsPending = make(map[common.Address]struct{})
	}
	// Track the amount of time wasted on hashing the account trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.AccountHashes += time.Since(start) }(time.Now())
	}
	return s.trie.Hash()
}

func (s *StateDB) Root() common.Hash {
	return s.trie.Hash()
}

// Prepare sets the current transaction hash and index and block hash which is
// used when the EVM emits new state logs.
func (s *StateDB) Prepare(thash common.Hash, ti int) {
	s.thash = thash
	s.txIndex = ti
}

func (s *StateDB) clearJournalAndRefund() {
	if len(s.journal.entries) > 0 {
		s.journal = newJournal()
		s.refund = 0
	}
	s.validRevisions = s.validRevisions[:0] // Snapshots can be created without journal entires
}

// AddAddressToAccessList adds the given address to the access list
func (s *StateDB) AddAddressToAccessList(addr common.Address) {
	if s.accessList.AddAddress(addr) {
		s.journal.append(accessListAddAccountChange{&addr})
	}
}

// AddSlotToAccessList adds the given (address, slot)-tuple to the access list
func (s *StateDB) AddSlotToAccessList(addr common.Address, slot common.Hash) {
	addrMod, slotMod := s.accessList.AddSlot(addr, slot)
	if addrMod {
		// In practice, this should not happen, since there is no way to enter the
		// scope of 'address' without having the 'address' become already added
		// to the access list (via call-variant, create, etc).
		// Better safe than sorry, though
		s.journal.append(accessListAddAccountChange{&addr})
	}
	if slotMod {
		s.journal.append(accessListAddSlotChange{
			address: &addr,
			slot:    &slot,
		})
	}
}

// PrepareAccessList handles the preparatory steps for executing a state transition with
// regards to both EIP-2929 and EIP-2930:
//
// - Add sender to access list (2929)
// - Add destination to access list (2929)
// - Add precompiles to access list (2929)
// - Add the contents of the optional tx access list (2930)
//
// This method should only be called if Yolov3/Berlin/2929+2930 is applicable at the current number.
func (s *StateDB) PrepareAccessList(sender common.Address, dst *common.Address, precompiles []common.Address, list types.AccessList) {
	// Clear out any leftover from previous executions
	s.accessList = newAccessList()

	s.AddAddressToAccessList(sender)
	if dst != nil {
		s.AddAddressToAccessList(*dst)
		// If it's a create-tx, the destination will be added inside evm.create
	}
	for _, addr := range precompiles {
		s.AddAddressToAccessList(addr)
	}
	for _, el := range list {
		s.AddAddressToAccessList(el.Address)
		for _, key := range el.StorageKeys {
			s.AddSlotToAccessList(el.Address, key)
		}
	}
}

// AddressInAccessList returns true if the given address is in the access list.
func (s *StateDB) AddressInAccessList(addr common.Address) bool {
	return s.accessList.ContainsAddress(addr)
}

// SlotInAccessList returns true if the given (address, slot)-tuple is in the access list.
func (s *StateDB) SlotInAccessList(addr common.Address, slot common.Hash) (addressPresent bool, slotPresent bool) {
	return s.accessList.Contains(addr, slot)
}

// Since the execution and submission of the block are parallel
// the snapshot needs to be updated before the state is committed.
func (s *StateDB) UpdateSnaps() error {
	// If snapshotting is enabled, update the snapshot tree with this new version
	if s.snap != nil {
		s.lock.Lock()
		defer s.lock.Unlock()

		// Finalize any pending changes and merge everything into the tries
		root := s.IntermediateRoot(true)

		if metrics.EnabledExpensive {
			defer func(start time.Time) { s.SnapshotCommits += time.Since(start) }(time.Now())
		}
		// Only update if there's a state transition
		if parent := s.snap.Root(); parent != root {
			if err := s.snaps.Update(root, parent, s.snapDestructs, s.snapAccounts, s.snapStorage); err != nil {
				log.Warn("Failed to update snapshot tree", "from", parent, "to", root, "err", err)
			}
		}
		s.snapDestructs, s.snapAccounts, s.snapStorage = nil, nil, nil
	}
	return nil
}

// Commit writes the state to the underlying in-memory trie database.
func (s *StateDB) Commit(deleteEmptyObjects bool) (common.Hash, error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	// Finalize any pending changes and merge everything into the tries
	s.IntermediateRoot(deleteEmptyObjects)

	// Increasing node version in memory database
	s.db.TrieDB().IncrVersion()

	// Commit objects to the trie, measuring the elapsed time
	var storageCommitted int
	codeWriter := s.db.TrieDB().DiskDB().NewBatch()
	for addr := range s.stateObjectsDirty {
		if obj := s.stateObjects[addr]; !obj.deleted {
			// Write any contract code associated with the state object
			if obj.code != nil && obj.dirtyCode {
				rawdb.WriteCode(codeWriter, common.BytesToHash(obj.CodeHash()), obj.code)
				obj.dirtyCode = false
			}
			// Write any storage changes in the state object to its storage trie
			committed, err := obj.CommitTrie(s.db)
			if err != nil {
				return common.Hash{}, err
			}
			storageCommitted += committed
		}
	}
	if len(s.stateObjectsDirty) > 0 {
		s.stateObjectsDirty = make(map[common.Address]struct{})
	}
	if codeWriter.ValueSize() > 0 {
		if err := codeWriter.Write(); err != nil {
			log.Crit("Failed to commit dirty codes", "error", err)
		}
	}
	// Write the account trie changes, measuing the amount of wasted time
	var start time.Time
	if metrics.EnabledExpensive {
		start = time.Now()
	}
	// Write trie changes.
	root, accountCommitted, err := s.trie.Commit(func(_ [][]byte, _ []byte, leaf []byte, parent common.Hash, _ []byte) error {
		var account types.StateAccount
		if err := rlp.DecodeBytes(leaf, &account); err != nil {
			return nil
		}
		if account.Root != emptyRoot {
			s.db.TrieDB().Reference(account.Root, parent)
		}
		return nil
	})
	if err != nil {
		return common.Hash{}, err
	}
	if metrics.EnabledExpensive {
		s.AccountCommits += time.Since(start)

		accountUpdatedMeter.Mark(int64(s.AccountUpdated))
		storageUpdatedMeter.Mark(int64(s.StorageUpdated))
		accountDeletedMeter.Mark(int64(s.AccountDeleted))
		storageDeletedMeter.Mark(int64(s.StorageDeleted))
		accountCommittedMeter.Mark(int64(accountCommitted))
		storageCommittedMeter.Mark(int64(storageCommitted))
		s.AccountUpdated, s.AccountDeleted = 0, 0
		s.StorageUpdated, s.StorageDeleted = 0, 0
	}
	// If snapshotting is enabled, update the snapshot tree with this new version
	if s.snap != nil {
		if metrics.EnabledExpensive {
			defer func(start time.Time) { s.SnapshotCommits += time.Since(start) }(time.Now())
		}
		// Only update if there's a state transition
		if parent := s.snap.Root(); parent != root {
			// Keep 128 diff layers in the memory, persistent layer is 129th.
			// - head layer is paired with HEAD state
			// - head-1 layer is paired with HEAD-1 state
			// - head-127 layer(bottom-most diff layer) is paired with HEAD-127 state
			if err := s.snaps.Cap(root, snapshot.CapLayers); err != nil {
				log.Warn("Failed to cap snapshot tree", "root", root, "layers", snapshot.CapLayers, "err", err)
			}
		}
		s.snap = nil
	}
	return root, err
}

func (s *StateDB) SetInt32(addr common.Address, key []byte, value int32) {
	s.SetState(addr, key, common.Int32ToBytes(value))
}
func (s *StateDB) SetInt64(addr common.Address, key []byte, value int64) {
	s.SetState(addr, key, common.Int64ToBytes(value))
}
func (s *StateDB) SetFloat32(addr common.Address, key []byte, value float32) {
	s.SetState(addr, key, common.Float32ToBytes(value))
}
func (s *StateDB) SetFloat64(addr common.Address, key []byte, value float64) {
	s.SetState(addr, key, common.Float64ToBytes(value))
}
func (s *StateDB) SetString(addr common.Address, key []byte, value string) {
	s.SetState(addr, key, []byte(value))
}
func (s *StateDB) SetByte(addr common.Address, key []byte, value byte) {
	s.SetState(addr, key, []byte{value})
}

func (s *StateDB) GetInt32(addr common.Address, key []byte) int32 {
	return common.BytesToInt32(s.GetState(addr, key))
}
func (s *StateDB) GetInt64(addr common.Address, key []byte) int64 {
	return common.BytesToInt64(s.GetState(addr, key))
}
func (s *StateDB) GetFloat32(addr common.Address, key []byte) float32 {
	return common.BytesToFloat32(s.GetState(addr, key))
}
func (s *StateDB) GetFloat64(addr common.Address, key []byte) float64 {
	return common.BytesToFloat64(s.GetState(addr, key))
}
func (s *StateDB) GetString(addr common.Address, key []byte) string {
	return string(s.GetState(addr, key))
}
func (s *StateDB) GetByte(addr common.Address, key []byte) byte {
	ret := s.GetState(addr, key)
	return ret[0]
}

func (s *StateDB) AddMinerEarnings(addr common.Address, amount *big.Int) {
	stateObject := s.GetOrNewStateObject(addr)
	if stateObject != nil {
		//stateObject.db = s
		stateObject.AddBalance(amount)
	}
}

func (s *StateDB) Merge(idx int, from, to *ParallelStateObject, deleteEmptyObjects bool) {
	if from.stateObject.address != to.stateObject.address {
		if from.stateObject.suicided || (deleteEmptyObjects && from.stateObject.empty()) {
			log.Warn("deleteStateObject", "from", from.stateObject.address.String(), "suicided", from.stateObject.suicided, "empty", from.stateObject.empty())
			s.deleteStateObject(from.stateObject)
		} else {
			s.stateObjects[from.stateObject.address] = from.stateObject
			s.journal.append(balanceChange{
				account: &from.stateObject.address,
				prev:    from.prevAmount,
			})
			s.stateObjectsDirty[from.stateObject.address] = struct{}{}
		}
	}
	if to.stateObject.suicided || (deleteEmptyObjects && to.stateObject.empty()) {
		log.Warn("deleteStateObject", "to", to.stateObject.address.String(), "suicided", to.stateObject.suicided, "empty", to.stateObject.empty())
		s.deleteStateObject(to.stateObject)
	} else {
		if to.createFlag {
			s.journal.append(createObjectChange{account: &to.stateObject.address})
		}
		s.stateObjects[to.stateObject.address] = to.stateObject
		s.journal.append(balanceChange{
			account: &to.stateObject.address,
			prev:    to.prevAmount,
		})
		s.stateObjectsDirty[to.stateObject.address] = struct{}{}
	}
}

func (s *StateDB) IncreaseTxIdx() {
	s.txIndex++
}

// Obtain version information maintained by governance
func (s *StateDB) ListActiveVersion() ([]gov.ActiveVersionValue, error) {
	//avListBytes := self.GetState(vm.GovContractAddr, gov.KeyActiveVersions())
	var avListBytes []byte
	stateObject := s.getStateObject(vm.GovContractAddr)
	if stateObject != nil {
		avListBytes = stateObject.removePrefixValue(stateObject.GetState(s.db, gov.KeyActiveVersions()))
	}

	if len(avListBytes) == 0 {
		return nil, nil
	}
	var avList []gov.ActiveVersionValue
	if err := json.Unmarshal(avListBytes, &avList); err != nil {
		return nil, err
	}
	return avList, nil
}

func (s *StateDB) GetCurrentActiveVersion() uint32 {
	avList, err := s.ListActiveVersion()
	if err != nil {
		log.Error("Cannot find active version list", "err", err)
		return 0
	}

	var version uint32
	if len(avList) == 0 {
		log.Warn("cannot find current active version, The ActiveVersion List is nil")
		return 0
	} else {
		version = avList[0].ActiveVersion
	}
	return version
}
