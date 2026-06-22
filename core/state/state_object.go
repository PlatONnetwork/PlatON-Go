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

package state

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"time"

	"golang.org/x/crypto/sha3"

	"github.com/PlatONnetwork/PlatON-Go/common"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/trie/trienode"
)

type Code []byte

func (c Code) String() string {
	return string(c) //strings.Join(Disassemble(c), " ")
}

type ValueStorage map[string][]byte

func (self ValueStorage) String() (str string) {
	for key, value := range self {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}
	return
}

func (self ValueStorage) Copy() ValueStorage {
	cpy := make(ValueStorage, len(self))
	for key, value := range self {
		cpy[key] = value
	}
	return cpy
}

// stateObject represents an Ethereum account which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// Account values can be accessed and modified through the object.
// Finally, call commitTrie to write the modified storage trie into a database.
type stateObject struct {
	address  common.Address
	addrHash common.Hash // hash of ethereum address of the account
	data     types.StateAccount
	db       *StateDB

	// Write caches.
	trie Trie // storage trie, which becomes non-nil on first access
	code Code // contract bytecode, which gets set when code is loaded

	originStorage  ValueStorage // Storage cache of original entries to dedup rewrites, reset for every transaction
	pendingStorage ValueStorage // Storage entries that need to be flushed to disk, at the end of an entire block
	dirtyStorage   ValueStorage // Storage entries that have been modified in the current transaction execution

	// Cache flags.
	dirtyCode bool // true if the code was updated

	// Flag whether the account was marked as self-destructed. The self-destructed account
	// is still accessible in the scope of same transaction.
	selfDestructed bool

	// Flag whether the account was marked as deleted. A self-destructed account
	// or an account that is considered as empty will be marked as deleted at
	// the end of transaction and no longer accessible anymore.
	deleted bool
}

// empty returns whether the account is considered empty.
func (s *stateObject) empty() bool {
	if cvm.PrecompiledContractCheckInstance.IsPlatONPrecompiledContract(s.address) {
		return false
	}
	return s.data.Nonce == 0 && s.data.Balance.Sign() == 0 && bytes.Equal(s.data.CodeHash, types.EmptyCodeHash.Bytes())
}

// newObject creates a state object.
func newObject(db *StateDB, address common.Address, data types.StateAccount) *stateObject {
	if data.Balance == nil {
		data.Balance = new(big.Int)
	}
	if data.CodeHash == nil {
		data.CodeHash = types.EmptyCodeHash.Bytes()
	}
	if data.Root == (common.Hash{}) {
		data.Root = types.EmptyRootHash
	}
	return &stateObject{
		db:             db,
		address:        address,
		addrHash:       crypto.Keccak256Hash(address[:]),
		data:           data,
		originStorage:  make(ValueStorage),
		pendingStorage: make(ValueStorage),
		dirtyStorage:   make(ValueStorage),
	}
}

// EncodeRLP implements rlp.Encoder.
func (s *stateObject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, &s.data)
}

func (s *stateObject) markSelfdestructed() {
	s.selfDestructed = true
}

func (s *stateObject) touch() {
	s.db.journal.append(touchChange{
		account: &s.address,
	})
	if s.address == ripemd {
		// Explicitly put it in the dirty-cache, which is otherwise generated from
		// flattened journals.
		s.db.journal.dirty(s.address)
	}
}

// getTrie returns the associated storage trie. The trie will be opened
// if it's not loaded previously. An error will be returned if trie can't
// be loaded.
func (s *stateObject) getTrie() (Trie, error) {
	if s.trie == nil {
		tr, err := s.db.db.OpenStorageTrie(s.db.originalRoot, s.addrHash, s.data.Root)
		if err != nil {
			return nil, err
		}
		s.trie = tr
	}
	return s.trie, nil
}

// GetState retrieves a value from the account storage trie.
func (s *stateObject) GetState(key []byte) []byte {
	// If we have a dirty value for this state entry, return it
	value, dirty := s.dirtyStorage[string(key)]
	if dirty {
		return value
	}
	// Otherwise return the entry's original value
	return s.GetCommittedState(key)
}

func (s *stateObject) getCommittedStateCache(key []byte) []byte {
	value, cached := s.originStorage[string(key)]
	if cached {
		return value
	}

	s.db.refLock.Lock()
	parentDB := s.db.parent
	parentCommitted := s.db.parentCommitted
	refLock := &s.db.refLock

	for parentDB != nil {
		value := parentDB.getStateObjectSnapshot(s.address, key)
		if value != nil {
			s.originStorage[string(key)] = value
			refLock.Unlock()
			return value
		} else if parentCommitted {
			refLock.Unlock()
			return nil
		}
		refLock.Unlock()
		parentDB.refLock.Lock()
		refLock = &parentDB.refLock
		if parentDB.parent == nil {
			break
		}
		parentCommitted = parentDB.parentCommitted
		parentDB = parentDB.parent
	}
	refLock.Unlock()

	return nil
}

// GetCommittedState retrieves a value from the committed account storage trie.
func (s *stateObject) GetCommittedState(key []byte) []byte {
	// If we have a pending write or clean cached, return that
	if value, pending := s.pendingStorage[string(key)]; pending {
		return value
	}
	// If we have the original value cached, return that
	if value := s.getCommittedStateCache(key); len(value) != 0 {
		return value
	}

	// If the object was destructed in *this* block (and potentially resurrected),
	// the storage has been cleared out, and we should *not* consult the previous
	// database about any storage values. The only possible alternatives are:
	//   1) resurrect happened, and new slot values were set -- those should
	//      have been handles via pendingStorage above.
	//   2) we don't have new values, and can deliver empty response back
	if _, destructed := s.db.stateObjectsDestruct[s.address]; destructed {
		return []byte{}
	}

	// If no live objects are available, attempt to use snapshots
	var (
		enc []byte
		err error
	)
	if s.db.snap != nil {
		if metrics.EnabledExpensive {
			defer func(start time.Time) { s.db.SnapshotStorageReads += time.Since(start) }(time.Now())
		}
		enc, err = s.db.snap.Storage(s.addrHash, crypto.Keccak256Hash(key[:]))
	}
	// If the snapshot is unavailable or reading from it fails, load from the database.
	if s.db.snap == nil || err != nil {
		if metrics.EnabledExpensive {
			defer func(start time.Time) { s.db.StorageReads += time.Since(start) }(time.Now())
		}
		tr, err := s.getTrie()
		if err != nil {
			s.db.setError(err)
			return []byte{}
		}
		enc, err = tr.GetStorage(s.address, key[:])
	}
	value := make([]byte, 0)
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			s.db.setError(err)
		}
		value = content
	}

	//log.Trace("GetCommittedState trie", "key", hex.EncodeToString(key), "value", len(value))
	s.originStorage[string(key)] = value
	return value
}

// SetState updates a value in account storage.
// set [prefixKey,value] to storage
func (s *stateObject) SetState(key, value []byte) {
	//if the new value is the same as old,don't set
	preValue := s.GetState(key)
	if bytes.Equal(preValue, value) {
		return
	}

	//New value is different, update and journal the change
	s.db.journal.append(storageChange{
		account:  &s.address,
		key:      key,
		preValue: preValue,
	})

	s.setState(key, value)
}

func (s *stateObject) setState(key []byte, value []byte) {
	cpy := make([]byte, len(value))
	copy(cpy, value)
	s.dirtyStorage[string(key)] = cpy
}

func (s *stateObject) getPrefixValue(pack, key, value []byte) []byte {
	// Empty value deleted on updateTrie
	if len(value) == 0 {
		return []byte{}
	}
	// Ensure the same Value, unique in the same trie and different trie values
	//prefix := append(s.data.StorageKeyPrefix, key...)
	prefix := append(s.data.StorageKeyPrefix, pack...)
	prefix = append(prefix, key...)
	prefixHash := common.Hash{}
	keccak := sha3.NewLegacyKeccak256()
	keccak.Write(prefix)
	keccak.Sum(prefixHash[:0])
	return append(prefixHash[:], value...)
}

func (s *stateObject) removePrefixValue(value []byte) []byte {
	if len(value) > common.HashLength {
		return value[common.HashLength:]
	}
	return []byte{}
}

// finalise moves all dirty storage slots into the pending area to be hashed or
// committed later. It is invoked at the end of every transaction.
func (s *stateObject) finalise() {
	for key, value := range s.dirtyStorage {
		s.pendingStorage[key] = value
	}
	if len(s.dirtyStorage) > 0 {
		s.dirtyStorage = make(ValueStorage)
	}
}

// updateTrie writes cached storage modifications into the object's storage trie.
// It will return nil if the trie has not been loaded and no changes have been
// made. An error will be returned if the trie can't be loaded/updated correctly.
func (s *stateObject) updateTrie() (Trie, error) {
	// Make sure all dirty slots are finalized into the pending storage area
	s.finalise()
	if len(s.pendingStorage) == 0 {
		return s.trie, nil
	}
	// Track the amount of time wasted on updating the storage trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.db.StorageUpdates += time.Since(start) }(time.Now())
	}
	// Retrieve the snapshot storage map for the object
	var storage map[common.Hash][]byte
	if s.db.snap != nil {
		// Retrieve the old storage map, if available, create a new one otherwise
		storage = s.db.snapStorage[s.addrHash]
		if storage == nil {
			storage = make(map[common.Hash][]byte)
			s.db.snapStorage[s.addrHash] = storage
		}
	}
	// Insert all the pending updates into the trie
	tr, err := s.getTrie()
	if err != nil {
		s.db.setError(err)
		return nil, err
	}
	for key, value := range s.pendingStorage {
		// Skip noop changes, persist actual changes
		oldValue := s.originStorage[key]
		if bytes.Equal(value, oldValue) {
			continue
		}

		s.originStorage[key] = value

		var v []byte
		if len(value) == 0 {
			if err := tr.DeleteStorage(s.address, []byte(key)); err != nil {
				s.db.setError(err)
				return nil, err
			}
			s.db.StorageDeleted += 1
		} else {
			// Encoding []byte cannot fail, ok to ignore the error.
			v, _ = rlp.EncodeToBytes(value[:])
			if err := tr.UpdateStorage(s.address, []byte(key), v); err != nil {
				s.db.setError(err)
				return nil, err
			}
			s.db.StorageUpdated += 1
		}
		// If state snapshotting is active, cache the data til commit
		if storage != nil {
			storage[crypto.Keccak256Hash([]byte(key))] = v // v will be nil if it's deleted
		}
	}

	if len(s.pendingStorage) > 0 {
		s.pendingStorage = make(ValueStorage)
	}

	return tr, nil
}

// UpdateRoot sets the trie root to the current root hash of. An error
// will be returned if trie root hash is not computed correctly.
func (s *stateObject) updateRoot() {
	tr, err := s.updateTrie()
	if err != nil {
		return
	}
	// If nothing changed, don't bother with hashing anything
	if tr == nil {
		return
	}
	// Track the amount of time wasted on hashing the storage trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.db.StorageHashes += time.Since(start) }(time.Now())
	}
	s.data.Root = tr.Hash()
}

// commitTrie submits the storage changes into the storage trie and re-computes
// the root. Besides, all trie changes will be collected in a nodeset and returned.
func (s *stateObject) commitTrie() (*trienode.NodeSet, error) {
	tr, err := s.updateTrie()
	if err != nil {
		return nil, err
	}
	// If nothing changed, don't bother with committing anything
	if tr == nil {
		return nil, nil
	}
	// Track the amount of time wasted on committing the storage trie
	if metrics.EnabledExpensive {
		defer func(start time.Time) { s.db.StorageCommits += time.Since(start) }(time.Now())
	}
	root, nodes, err := tr.Commit(false)

	if err == nil {
		s.data.Root = root
	}
	return nodes, err
}

// AddBalance adds amount to s's balance.
// It is used to add funds to the destination account of a transfer.
func (s *stateObject) AddBalance(amount *big.Int) {
	// EIP161: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if s.empty() {
			s.touch()
		}

		return
	}
	s.SetBalance(new(big.Int).Add(s.Balance(), amount))
}

// SubBalance removes amount from s's balance.
// It is used to remove funds from the origin account of a transfer.
func (s *stateObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetBalance(new(big.Int).Sub(s.Balance(), amount))
}

func (s *stateObject) SetBalance(amount *big.Int) {
	s.db.journal.append(balanceChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Balance),
	})
	s.setBalance(amount)
}

func (s *stateObject) setBalance(amount *big.Int) {
	s.data.Balance = amount
}

func (s *stateObject) deepCopy(db *StateDB) *stateObject {
	stateObject := newObject(db, s.address, s.data)
	if s.trie != nil {
		stateObject.trie = db.db.CopyTrie(s.trie)
	}
	stateObject.code = s.code
	stateObject.dirtyStorage = s.dirtyStorage.Copy()
	stateObject.originStorage = s.originStorage.Copy()
	stateObject.pendingStorage = s.pendingStorage.Copy()
	stateObject.selfDestructed = s.selfDestructed
	stateObject.dirtyCode = s.dirtyCode
	stateObject.deleted = s.deleted
	return stateObject
}

// Copy account status, recreate trie
func (s *stateObject) copy(db *StateDB) *stateObject {
	stateObject := newObject(db, s.address, s.data)
	if s.trie != nil {
		stateObject.trie = db.db.NewTrie(s.trie)
	}
	stateObject.code = s.code
	stateObject.selfDestructed = s.selfDestructed
	stateObject.dirtyCode = s.dirtyCode
	stateObject.deleted = s.deleted
	return stateObject
}

//
// Attribute accessors
//

// Address returns the address of the contract/account
func (s *stateObject) Address() common.Address {
	return s.address
}

// Code returns the contract code associated with this object, if any.
func (s *stateObject) Code() []byte {
	if s.code != nil {
		return s.code
	}
	if bytes.Equal(s.CodeHash(), types.EmptyCodeHash.Bytes()) {
		return nil
	}
	code, err := s.db.db.ContractCode(s.addrHash, common.BytesToHash(s.CodeHash()))
	if err != nil {
		s.db.setError(fmt.Errorf("can't load code hash %x: %v", s.CodeHash(), err))
	}
	s.code = code
	return code
}

// CodeSize returns the size of the contract code associated with this object,
// or zero if none. This method is an almost mirror of Code, but uses a cache
// inside the database to avoid loading codes seen recently.
func (s *stateObject) CodeSize() int {
	if s.code != nil {
		return len(s.code)
	}
	if bytes.Equal(s.CodeHash(), types.EmptyCodeHash.Bytes()) {
		return 0
	}
	size, err := s.db.db.ContractCodeSize(s.addrHash, common.BytesToHash(s.CodeHash()))
	if err != nil {
		s.db.setError(fmt.Errorf("can't load code size %x: %v", s.CodeHash(), err))
	}
	return size
}

func (s *stateObject) SetCode(codeHash common.Hash, code []byte) {
	prevcode := s.Code()
	s.db.journal.append(codeChange{
		account:  &s.address,
		prevhash: s.CodeHash(),
		prevcode: prevcode,
	})
	s.setCode(codeHash, code)
}

func (s *stateObject) setCode(codeHash common.Hash, code []byte) {
	s.code = code
	s.data.CodeHash = codeHash[:]
	s.dirtyCode = true
}

func (s *stateObject) SetNonce(nonce uint64) {
	s.db.journal.append(nonceChange{
		account: &s.address,
		prev:    s.data.Nonce,
	})
	s.setNonce(nonce)
}

func (s *stateObject) setNonce(nonce uint64) {
	s.data.Nonce = nonce
}

func (s *stateObject) CodeHash() []byte {
	return s.data.CodeHash
}

func (s *stateObject) Balance() *big.Int {
	return s.data.Balance
}

func (s *stateObject) Nonce() uint64 {
	return s.data.Nonce
}
