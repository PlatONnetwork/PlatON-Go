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
	"errors"
	"fmt"
	"math/big"
	"sort"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie"
)

type revision struct {
	id           int
	journalIndex int
}

var (
	// emptyState is the known hash of an empty state trie entry.
	emptyState = crypto.Keccak256Hash(nil)

	// emptyCode is the known hash of the empty EVM bytecode.
	emptyCode = crypto.Keccak256Hash(nil)

	emptyStorage = crypto.Keccak256Hash(nil)
)

// StateDBs within the ethereum protocol are used to store anything
// within the merkle trie. StateDBs take care of caching and storing
// nested states. It's the general query interface to retrieve:
// * Contracts
// * Accounts
type StateDB struct {
	db   Database
	trie Trie

	// This map holds 'live' objects, which will get modified while processing a state transition.
	stateObjects      map[common.Address]*stateObject
	stateObjectsDirty map[common.Address]struct{}

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// The refund counter, also used by state transitioning.
	refund uint64

	thash, bhash common.Hash
	txIndex      int
	logs         map[common.Hash][]*types.Log
	logSize      uint

	preimages map[common.Hash][]byte

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
}

// Create a new state from a given trie.
func New(root common.Hash, db Database) (*StateDB, error) {
	tr, err := db.OpenTrie(root)
	if err != nil {
		return nil, err
	}
	state := &StateDB{
		db:                 db,
		trie:               tr,
		stateObjects:       make(map[common.Address]*stateObject),
		stateObjectsDirty:  make(map[common.Address]struct{}),
		logs:               make(map[common.Hash][]*types.Log),
		preimages:          make(map[common.Hash][]byte),
		journal:            newJournal(),
		clearReferenceFunc: make([]func(), 0),
	}
	return state, nil
}

// New StateDB based on the parent StateDB
func (self *StateDB) NewStateDB() *StateDB {
	stateDB := &StateDB{
		db:                 self.db,
		trie:               self.db.NewTrie(self.trie),
		stateObjects:       make(map[common.Address]*stateObject),
		stateObjectsDirty:  make(map[common.Address]struct{}),
		logs:               make(map[common.Hash][]*types.Log),
		preimages:          make(map[common.Hash][]byte),
		journal:            newJournal(),
		parent:             self,
		clearReferenceFunc: make([]func(), 0),
	}

	index := self.AddReferenceFunc(stateDB.clearParentRef)
	stateDB.referenceFuncIndex = index

	//if stateDB.parent != nil {
	//	stateDB.parent.DumpStorage(false)
	//}
	return stateDB
}

func (self *StateDB) HadParent() bool {
	self.refLock.Lock()
	defer self.refLock.Unlock()
	return self.parent != nil
}

/*func (self *StateDB) DumpStorage(check bool) {
	log.Debug("statedb stateobjects", "len", len(self.stateObjects), "root", self.Root())
	disk, err := New(self.Root(), self.db)
	if check && err != nil {
		panic(fmt.Sprintf("new statdb error, root:%s, error:%s", self.Root().String(), err.Error()))
	}
	for addr, obj := range self.stateObjects {
		log.Debug("dump storage", "addr", addr.String())
		for k, v := range obj.originStorage {
			log.Debug(fmt.Sprintf("origin: key:%s, valueKey:%s, value:[%s] len:%d", hexutil.Encode([]byte(k)), v.String(), hexutil.Encode(vk), len(vk)))
			if check {
				vg := disk.GetCommittedState(addr, []byte(k))

				if check && !bytes.Equal(vk, vg) {
					panic(fmt.Sprintf("not equal, key:%s, value:[%s] len:%d", hexutil.Encode([]byte(k)), hexutil.Encode(vg), len(vg)))
				}
			}
		}

		for k, vk := range obj.dirtyStorage {
			v, ok := obj.dirtyValueStorage[vk]
			if ok {
				log.Debug("dirty: key:%s, valueKey:%s, value:%s len:%d", hexutil.Encode([]byte(k)), vk.String(), hexutil.Encode(v.Value), len(v.Value))
				if check {
					vg := disk.GetCommittedState(addr, []byte(k))

					if check && !bytes.Equal(v.Value, vg) {
						panic(fmt.Sprintf("not equal, key:%s, value:%s len:%d", hexutil.Encode([]byte(k)), hexutil.Encode(vg), len(vg)))
					}
				}
			}
		}
	}
}*/

// setError remembers the first non-nil error it is called with.
func (self *StateDB) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *StateDB) Error() error {
	return self.dbErr
}

// Reset clears out all ephemeral state objects from the state db, but keeps
// the underlying state trie to avoid reloading data for the next operations.
func (self *StateDB) Reset(root common.Hash) error {
	tr, err := self.db.OpenTrie(root)
	if err != nil {
		return err
	}
	self.trie = tr
	self.stateObjects = make(map[common.Address]*stateObject)
	self.stateObjectsDirty = make(map[common.Address]struct{})
	self.thash = common.Hash{}
	self.bhash = common.Hash{}
	self.txIndex = 0
	self.logs = make(map[common.Hash][]*types.Log)
	self.logSize = 0
	self.preimages = make(map[common.Hash][]byte)
	self.clearJournalAndRefund()
	return nil
}

func (self *StateDB) AddLog(logInfo *types.Log) {
	self.journal.append(addLogChange{txhash: self.thash})

	logInfo.TxHash = self.thash
	logInfo.BlockHash = self.bhash
	logInfo.TxIndex = uint(self.txIndex)
	logInfo.Index = self.logSize
	self.logs[self.thash] = append(self.logs[self.thash], logInfo)
	self.logSize++
}

func (self *StateDB) GetLogs(hash common.Hash) []*types.Log {
	return self.logs[hash]
}

func (self *StateDB) Logs() []*types.Log {
	var logs []*types.Log
	for _, lgs := range self.logs {
		logs = append(logs, lgs...)
	}
	return logs
}

// AddPreimage records a SHA3 preimage seen by the VM.
func (self *StateDB) AddPreimage(hash common.Hash, preimage []byte) {
	if _, ok := self.preimages[hash]; !ok {
		self.journal.append(addPreimageChange{hash: hash})
		pi := make([]byte, len(preimage))
		copy(pi, preimage)
		self.preimages[hash] = pi
	}
}

// Preimages returns a list of SHA3 preimages that have been submitted.
func (self *StateDB) Preimages() map[common.Hash][]byte {
	return self.preimages
}

// AddRefund adds gas to the refund counter
func (self *StateDB) AddRefund(gas uint64) {
	self.journal.append(refundChange{prev: self.refund})
	self.refund += gas
}

// SubRefund removes gas from the refund counter.
// This method will panic if the refund counter goes below zero
func (self *StateDB) SubRefund(gas uint64) {
	self.journal.append(refundChange{prev: self.refund})
	if gas > self.refund {
		panic("Refund counter below zero")
	}
	self.refund -= gas
}

// Exist reports whether the given account address exists in the state.
// Notably this also returns true for suicided accounts.
func (self *StateDB) Exist(addr common.Address) bool {
	return self.getStateObject(addr) != nil
}

// Empty returns whether the state object is either non-existent
// or empty according to the EIP161 specification (balance = nonce = code = 0)
func (self *StateDB) Empty(addr common.Address) bool {
	so := self.getStateObject(addr)
	return so == nil || so.empty()
}

// Retrieve the balance from the given address or 0 if object not found
func (self *StateDB) GetBalance(addr common.Address) *big.Int {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Balance()
	}
	return common.Big0
}

func (self *StateDB) GetNonce(addr common.Address) uint64 {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Nonce()
	}

	return 0
}

func (self *StateDB) GetCode(addr common.Address) []byte {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.Code(self.db)
	}
	return nil
}

func (self *StateDB) GetCodeSize(addr common.Address) int {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return 0
	}
	if stateObject.code != nil {
		return len(stateObject.code)
	}
	size, err := self.db.ContractCodeSize(stateObject.addrHash, common.BytesToHash(stateObject.CodeHash()))
	if err != nil {
		self.setError(err)
	}
	return size
}

func (self *StateDB) GetCodeHash(addr common.Address) common.Hash {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return common.Hash{}
	}
	return common.BytesToHash(stateObject.CodeHash())
}

// GetState retrieves a value from the given account's storage trie.
func (self *StateDB) GetState(addr common.Address, key []byte) []byte {
	self.lock.Lock()
	defer self.lock.Unlock()
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.removePrefixValue(stateObject.GetState(self.db, key))
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

// GetProof returns the MerkleProof for a given Account
func (s *StateDB) GetProof(a common.Address) ([][]byte, error) {
	var proof proofList
	err := s.trie.Prove(crypto.Keccak256(a.Bytes()), 0, &proof)
	return [][]byte(proof), err
}

// GetProof returns the StorageProof for given key
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
func (self *StateDB) GetCommittedState(addr common.Address, key []byte) []byte {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.removePrefixValue(stateObject.GetCommittedState(self.db, key))
	}
	return []byte{}
}

// Database retrieves the low level database supporting the lower level trie ops.
func (self *StateDB) Database() Database {
	return self.db
}

// StorageTrie returns the storage trie of an account.
// The return value is a copy and is nil for non-existent accounts.
func (self *StateDB) StorageTrie(addr common.Address) Trie {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return nil
	}
	cpy := stateObject.deepCopy(self)
	return cpy.updateTrie(self.db)
}

func (self *StateDB) HasSuicided(addr common.Address) bool {
	stateObject := self.getStateObject(addr)
	if stateObject != nil {
		return stateObject.suicided
	}
	return false
}

/*
 * SETTERS
 */

// AddBalance adds amount to the account associated with addr.
func (self *StateDB) AddBalance(addr common.Address, amount *big.Int) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.AddBalance(amount)
	}
}

// SubBalance subtracts amount from the account associated with addr.
func (self *StateDB) SubBalance(addr common.Address, amount *big.Int) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SubBalance(amount)
	}
}

func (self *StateDB) SetBalance(addr common.Address, amount *big.Int) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetBalance(amount)
	}
}

func (self *StateDB) SetNonce(addr common.Address, nonce uint64) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetNonce(nonce)
	}
}

func (self *StateDB) SetCode(addr common.Address, code []byte) {
	stateObject := self.GetOrNewStateObject(addr)
	if stateObject != nil {
		stateObject.SetCode(crypto.Keccak256Hash(code), code)
	}
}

func (self *StateDB) SetState(address common.Address, key, value []byte) {
	self.lock.Lock()
	stateObject := self.GetOrNewStateObject(address)

	if stateObject != nil {
		//prefixKey := stateObject.getPrefixKey(key)
		stateObject.SetState(self.db, key, stateObject.getPrefixValue(key, value))
	}
	self.lock.Unlock()
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
//	keccak := sha3.NewKeccak256()
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
func (self *StateDB) Suicide(addr common.Address) bool {
	stateObject := self.getStateObject(addr)
	if stateObject == nil {
		return false
	}
	self.journal.append(suicideChange{
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
func (self *StateDB) updateStateObject(stateObject *stateObject) {
	addr := stateObject.Address()
	data, err := rlp.EncodeToBytes(stateObject)
	if err != nil {
		panic(fmt.Errorf("can't encode object at %x: %v", addr[:], err))
	}
	self.setError(self.trie.TryUpdate(addr[:], data))
}

// deleteStateObject removes the given object from the state trie.
func (self *StateDB) deleteStateObject(stateObject *stateObject) {
	stateObject.deleted = true
	addr := stateObject.Address()
	self.setError(self.trie.TryDelete(addr[:]))
}

// Get the current StateDB cache and the parent StateDB cache
func (self *StateDB) getStateObjectCache(addr common.Address) (stateObject *stateObject) {
	// Prefer 'live' objects.
	if obj := self.stateObjects[addr]; obj != nil {
		return obj
	}
	self.refLock.Lock()
	parentDB := self.parent
	parentCommitted := self.parentCommitted
	refLock := &self.refLock

	for parentDB != nil {
		obj := parentDB.getStateObjectLocalCache(addr)
		if obj != nil {
			refLock.Unlock()
			cpy := obj.copy(self)
			self.setStateObject(cpy)
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
func (self *StateDB) getStateObjectLocalCache(addr common.Address) (stateObject *stateObject) {
	if obj := self.stateObjects[addr]; obj != nil {
		if obj.deleted {
			return nil
		}
		return obj
	}
	return nil
}

// Find stateObject storage in cache
func (self *StateDB) getStateObjectSnapshot(addr common.Address, key []byte) []byte {
	if obj := self.stateObjects[addr]; obj != nil {
		if obj.deleted {
			return nil
		}
		value, dirty := obj.dirtyStorage[string(key)]
		if dirty {
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
func (self *StateDB) AddReferenceFunc(fn func()) int {
	self.refLock.Lock()
	defer self.refLock.Unlock()
	// It must not be nil
	if self.clearReferenceFunc == nil {
		panic("statedb had cleared")
	}
	self.clearReferenceFunc = append(self.clearReferenceFunc, fn)
	return len(self.clearReferenceFunc) - 1
}

// Clear reference when StateDB is committed
func (self *StateDB) ClearReference() {
	self.refLock.Lock()
	defer self.refLock.Unlock()
	for _, fn := range self.clearReferenceFunc {
		if nil != fn {
			fn()
		}
	}
	log.Trace("clear all ref", "reflen", len(self.clearReferenceFunc))
	if self.parent != nil {
		if len(self.parent.clearReferenceFunc) > 0 {
			panic("parent ref > 0")
		}
	}
	self.clearReferenceFunc = nil
	self.parent = nil
}

// Clear reference by index
func (self *StateDB) ClearIndexReference(index int) {
	self.refLock.Lock()
	defer self.refLock.Unlock()

	if len(self.clearReferenceFunc) > index && self.clearReferenceFunc[index] != nil {
		//fn := self.clearReferenceFunc[index]
		//fn()
		log.Trace("Before clear index ref", "reflen", len(self.clearReferenceFunc), "index", index)
		//self.clearReferenceFunc = append(self.clearReferenceFunc[:index], self.clearReferenceFunc[index+1:]...)
		self.clearReferenceFunc[index] = nil
		log.Trace("After clear index ref", "reflen", len(self.clearReferenceFunc), "index", index)
	}
}

// Clear Parent reference
func (self *StateDB) ClearParentReference() {
	self.refLock.Lock()
	defer self.refLock.Unlock()

	if self.parent != nil && self.referenceFuncIndex >= 0 {
		self.parent.ClearIndexReference(self.referenceFuncIndex)
		self.parent = nil
		self.referenceFuncIndex = -1
	}
}

// Retrieve a state object given by the address. Returns nil if not found.
func (self *StateDB) getStateObject(addr common.Address) (stateObject *stateObject) {
	if obj := self.getStateObjectCache(addr); obj != nil {
		if obj.deleted {
			return nil
		}
		return obj
	}
	// Load the object from the database.
	enc, err := self.trie.TryGet(addr[:])
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data Account
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state object", "addr", addr, "err", err)
		return nil
	}
	// [NOTE]: set the prefix for storage key
	if data.empty() {
		data.StorageKeyPrefix = addr.Bytes()
	}
	// Insert into the live set.
	obj := newObject(self, addr, data)
	self.setStateObject(obj)
	return obj
}

func (self *StateDB) setStateObject(object *stateObject) {
	if len(self.clearReferenceFunc) > 0 {
		panic("statedb readonly")
	}
	self.stateObjects[object.Address()] = object
}

// Retrieve a state object or create a new state object if nil.
func (self *StateDB) GetOrNewStateObject(addr common.Address) *stateObject {
	stateObject := self.getStateObject(addr)
	if stateObject == nil || stateObject.deleted {
		stateObject, _ = self.createObject(addr)
	}
	return stateObject
}

// createObject creates a new state object. If there is an existing account with
// the given address, it is overwritten and returned as the second return value.
func (self *StateDB) createObject(addr common.Address) (newobj, prev *stateObject) {
	prev = self.getStateObject(addr)
	if prev == nil {
		newobj = newObject(self, addr, Account{StorageKeyPrefix: addr.Bytes()})
		self.journal.append(createObjectChange{account: &addr})
	} else {
		prefix := make([]byte, len(prev.data.StorageKeyPrefix))
		copy(prefix, prev.data.StorageKeyPrefix)
		newobj = newObject(self, addr, Account{StorageKeyPrefix: prefix})
		self.journal.append(resetObjectChange{prev: prev})
	}
	newobj.setNonce(0) // sets the object to dirty
	self.setStateObject(newobj)
	return newobj, prev
}

// CreateAccount explicitly creates a state object. If a state object with the address
// already exists the balance is carried over to the new account.
//
// CreateAccount is called during the EVM CREATE operation. The situation might arise that
// a contract does the following:
//
//   1. sends funds to sha(account ++ (nonce + 1))
//   2. tx_create(sha(account ++ nonce)) (note that this gets the address of 1)
//
// Carrying over the balance ensures that Ether doesn't disappear.
func (self *StateDB) CreateAccount(addr common.Address) {
	newObj, prev := self.createObject(addr)
	if prev != nil {
		newObj.setBalance(prev.data.Balance)
	}
}

func (self *StateDB) TxHash() common.Hash {
	return self.thash
}

func (self *StateDB) TxIdx() uint32 {
	return uint32(self.txIndex)
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
func (self *StateDB) Copy() *StateDB {
	self.lock.Lock()
	defer self.lock.Unlock()

	// Copy all the basic fields, initialize the memory ones
	state := &StateDB{
		db:                 self.db,
		trie:               self.db.CopyTrie(self.trie),
		stateObjects:       make(map[common.Address]*stateObject, len(self.journal.dirties)),
		stateObjectsDirty:  make(map[common.Address]struct{}, len(self.journal.dirties)),
		refund:             self.refund,
		logs:               make(map[common.Hash][]*types.Log, len(self.logs)),
		logSize:            self.logSize,
		preimages:          make(map[common.Hash][]byte),
		journal:            newJournal(),
		clearReferenceFunc: make([]func(), 0),
	}

	// Copy the dirty states, logs, and preimages
	for addr := range self.journal.dirties {
		// As documented [here](https://github.com/ethereum/go-ethereum/pull/16485#issuecomment-380438527),
		// and in the Finalise-method, there is a case where an object is in the journal but not
		// in the stateObjects: OOG after touch on ripeMD prior to Byzantium. Thus, we need to check for
		// nil
		if object, exist := self.stateObjects[addr]; exist {
			state.stateObjects[addr] = object.deepCopy(state)
			state.stateObjectsDirty[addr] = struct{}{}
		}
	}
	// Above, we don't copy the actual journal. This means that if the copy is copied, the
	// loop above will be a no-op, since the copy's journal is empty.
	// Thus, here we iterate over stateObjects, to enable copies of copies
	for addr := range self.stateObjectsDirty {
		if _, exist := state.stateObjects[addr]; !exist {
			state.stateObjects[addr] = self.stateObjects[addr].deepCopy(state)
			state.stateObjectsDirty[addr] = struct{}{}
		}
	}
	for hash, logs := range self.logs {
		cpy := make([]*types.Log, len(logs))
		for i, l := range logs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		state.logs[hash] = cpy
	}
	for hash, preimage := range self.preimages {
		state.preimages[hash] = preimage
	}
	// Copy parent state
	self.refLock.Lock()
	if self.parent != nil {
		if !self.parentCommitted {
			state.parent = self.parent
			state.parent.AddReferenceFunc(state.clearParentRef)
		} else {
			self.parent = nil
		}
	}
	state.parentCommitted = self.parentCommitted
	self.refLock.Unlock()

	return state
}

// Clear parent StateDB reference
func (self *StateDB) clearParentRef() {
	self.refLock.Lock()
	defer self.refLock.Unlock()

	if self.parent != nil {
		self.parentCommitted = true
		log.Trace("clearParentRef", "parent root", self.parent.Root().String())
		// Parent is nil, find the parent state based on current StateDB
		self.parent = nil
	}
}

// Snapshot returns an identifier for the current revision of the state.
func (self *StateDB) Snapshot() int {
	id := self.nextRevisionId
	self.nextRevisionId++
	self.validRevisions = append(self.validRevisions, revision{id, self.journal.length()})
	return id
}

// RevertToSnapshot reverts all state changes made since the given revision.
func (self *StateDB) RevertToSnapshot(revid int) {
	// Find the snapshot in the stack of valid snapshots.
	idx := sort.Search(len(self.validRevisions), func(i int) bool {
		return self.validRevisions[i].id >= revid
	})
	if idx == len(self.validRevisions) || self.validRevisions[idx].id != revid {
		panic(fmt.Errorf("revision id %v cannot be reverted", revid))
	}
	snapshot := self.validRevisions[idx].journalIndex

	// Replay the journal to undo changes and remove invalidated snapshots
	self.journal.revert(self, snapshot)
	self.validRevisions = self.validRevisions[:idx]
}

// GetRefund returns the current value of the refund counter.
func (self *StateDB) GetRefund() uint64 {
	return self.refund
}

// Finalise finalises the state by removing the self destructed objects
// and clears the journal as well as the refunds.
func (s *StateDB) Finalise(deleteEmptyObjects bool) {
	for addr := range s.journal.dirties {
		stateObject, exist := s.stateObjects[addr]
		if !exist {
			// ripeMD is 'touched' at block 1714175, in tx 0x1237f737031e40bcde4a8b7e717b2d15e3ecadfe49bb1bbc71ee9deb09c6fcf2
			// That tx goes out of gas, and although the notion of 'touched' does not exist there, the
			// touch-event will still be recorded in the journal. Since ripeMD is a special snowflake,
			// it will persist in the journal even though the journal is reverted. In this special circumstance,
			// it may exist in `s.journal.dirties` but not in `s.stateObjects`.
			// Thus, we can safely ignore it here
			continue
		}
		if stateObject.suicided || (deleteEmptyObjects && stateObject.empty()) {
			s.deleteStateObject(stateObject)
		} else {
			stateObject.updateRoot(s.db)
			s.updateStateObject(stateObject)
			/*	log.Trace("Finalise single", "address", stateObject.address.String(), "balance", stateObject.Balance().Uint64(), "nonce", stateObject.Nonce(),
				"codeHash", common.Bytes2Hex(stateObject.CodeHash()), "storageRoot", stateObject.data.Root.String(), "storageKeyPrefix", common.Bytes2Hex(stateObject.data.StorageKeyPrefix))*/
		}
		s.stateObjectsDirty[addr] = struct{}{}
	}
	// Invalidate journal because reverting across transactions is not allowed.
	s.clearJournalAndRefund()
}

// IntermediateRoot computes the current root hash of the state trie.
// It is called in between transactions to get the root hash that
// goes into transaction receipts.
func (s *StateDB) IntermediateRoot(deleteEmptyObjects bool) common.Hash {
	s.Finalise(deleteEmptyObjects)
	//return s.trie.Hash()
	return s.trie.ParallelHash()
}

func (s *StateDB) Root() common.Hash {
	//return s.trie.Hash()
	return s.trie.ParallelHash()
}

// Prepare sets the current transaction hash and index and block hash which is
// used when the EVM emits new state logs.
func (self *StateDB) Prepare(thash, bhash common.Hash, ti int) {
	self.thash = thash
	self.bhash = bhash
	self.txIndex = ti
}

func (s *StateDB) clearJournalAndRefund() {
	s.journal = newJournal()
	s.validRevisions = s.validRevisions[:0]
	s.refund = 0
}

// Commit writes the state to the underlying in-memory trie database.
func (s *StateDB) Commit(deleteEmptyObjects bool) (root common.Hash, err error) {
	s.lock.Lock()
	defer s.lock.Unlock()

	defer s.clearJournalAndRefund()

	for addr := range s.journal.dirties {
		s.stateObjectsDirty[addr] = struct{}{}
	}
	// Commit objects to the trie.
	for addr, stateObject := range s.stateObjects {
		_, isDirty := s.stateObjectsDirty[addr]
		switch {
		case stateObject.suicided || (isDirty && deleteEmptyObjects && stateObject.empty()):
			// If the object has been removed, don't bother syncing it
			// and just mark it for deletion in the trie.
			s.deleteStateObject(stateObject)
		case isDirty:
			// Write any contract code associated with the state object
			if stateObject.code != nil && stateObject.dirtyCode {
				s.db.TrieDB().InsertBlob(common.BytesToHash(stateObject.CodeHash()), stateObject.code)
				stateObject.dirtyCode = false
			}
			// Write any storage changes in the state object to its storage trie.
			if err := stateObject.CommitTrie(s.db); err != nil {
				return common.Hash{}, err
			}
			// Update the object in the main account trie.
			s.updateStateObject(stateObject)
		}
		delete(s.stateObjectsDirty, addr)
	}
	// Write trie changes.
	//root, err = s.trie.Commit(func(leaf []byte, parent common.Hash) error {
	root, err = s.trie.ParallelCommit(func(leaf []byte, parent common.Hash) error {
		var account Account
		if err := rlp.DecodeBytes(leaf, &account); err != nil {
			return nil
		}
		if account.Root != emptyState {
			s.db.TrieDB().Reference(account.Root, parent)
		}
		code := common.BytesToHash(account.CodeHash)
		if code != emptyCode {
			s.db.TrieDB().Reference(code, parent)
		}
		return nil
	})

	return root, err
}

func (self *StateDB) SetInt32(addr common.Address, key []byte, value int32) {
	self.SetState(addr, key, common.Int32ToBytes(value))
}
func (self *StateDB) SetInt64(addr common.Address, key []byte, value int64) {
	self.SetState(addr, key, common.Int64ToBytes(value))
}
func (self *StateDB) SetFloat32(addr common.Address, key []byte, value float32) {
	self.SetState(addr, key, common.Float32ToBytes(value))
}
func (self *StateDB) SetFloat64(addr common.Address, key []byte, value float64) {
	self.SetState(addr, key, common.Float64ToBytes(value))
}
func (self *StateDB) SetString(addr common.Address, key []byte, value string) {
	self.SetState(addr, key, []byte(value))
}
func (self *StateDB) SetByte(addr common.Address, key []byte, value byte) {
	self.SetState(addr, key, []byte{value})
}

func (self *StateDB) GetInt32(addr common.Address, key []byte) int32 {
	return common.BytesToInt32(self.GetState(addr, key))
}
func (self *StateDB) GetInt64(addr common.Address, key []byte) int64 {
	return common.BytesToInt64(self.GetState(addr, key))
}
func (self *StateDB) GetFloat32(addr common.Address, key []byte) float32 {
	return common.BytesToFloat32(self.GetState(addr, key))
}
func (self *StateDB) GetFloat64(addr common.Address, key []byte) float64 {
	return common.BytesToFloat64(self.GetState(addr, key))
}
func (self *StateDB) GetString(addr common.Address, key []byte) string {
	return string(self.GetState(addr, key))
}
func (self *StateDB) GetByte(addr common.Address, key []byte) byte {
	ret := self.GetState(addr, key)
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

func (self *StateDB) IncreaseTxIdx() {
	self.txIndex++
}
