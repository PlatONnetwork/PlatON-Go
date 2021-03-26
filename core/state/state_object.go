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

	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"

	"github.com/PlatONnetwork/PlatON-Go/common"
	cvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var emptyCodeHash = crypto.Keccak256(nil)

type Code []byte

func (self Code) String() string {
	return string(self)
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
// Finally, call CommitTrie to write the modified storage trie into a database.
type stateObject struct {
	address  common.Address
	addrHash common.Hash // hash of ethereum address of the account
	data     Account
	db       *StateDB

	// DB error.
	// State objects are used by the consensus core and VM which are
	// unable to deal with database-level errors. Any error that occurs
	// during a database read is memoized here and will eventually be returned
	// by StateDB.Commit.
	dbErr error

	// Write caches.
	trie Trie
	// storage trie, which becomes non-nil on first access
	code Code // contract bytecode, which gets set when code is loaded

	//originStorage      Storage      // Storage cache of original entries to dedup rewrites
	//originValueStorage ValueStorage // Storage cache of original entries to dedup rewrites
	//
	//dirtyStorage      Storage      // Storage entries that need to be flushed to disk
	//dirtyValueStorage ReferenceValueStorage // Storage entries that need to be flushed to disk

	originStorage ValueStorage // Storage cache of original entries to dedup rewrites

	dirtyStorage ValueStorage // Storage entries that need to be flushed to disk

	// Cache flags.
	// When an object is marked suicided it will be delete from the trie
	// during the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
}

// empty returns whether the account is considered empty.
func (s *stateObject) empty() bool {
	if cvm.PrecompiledContractCheckInstance.IsPlatONPrecompiledContract(s.address) {
		return false
	}
	return s.data.Nonce == 0 && s.data.Balance.Sign() == 0 && bytes.Equal(s.data.CodeHash, emptyCodeHash)
}

// Account is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type Account struct {
	Nonce            uint64
	Balance          *big.Int
	Root             common.Hash // merkle root of the storage trie
	CodeHash         []byte
	StorageKeyPrefix []byte // A prefix added to the `key` to ensure that data between different accounts are not shared
}

func (self *Account) empty() bool {
	if self.Nonce != 0 {
		return false
	}
	if self.Balance.Cmp(common.Big0) != 0 {
		return false
	}
	if self.Root != common.ZeroHash {
		return false
	}
	if len(self.CodeHash) != 0 {
		return false
	}
	if len(self.StorageKeyPrefix) != 0 {
		return false
	}
	return true
}

// newObject creates a state object.
func newObject(db *StateDB, address common.Address, data Account) *stateObject {
	if data.Balance == nil {
		data.Balance = new(big.Int)
	}
	if data.CodeHash == nil {
		data.CodeHash = emptyCodeHash
	}
	return &stateObject{
		db:            db,
		address:       address,
		addrHash:      crypto.Keccak256Hash(address[:]),
		data:          data,
		originStorage: make(ValueStorage),
		dirtyStorage:  make(ValueStorage),
	}
}

// EncodeRLP implements rlp.Encoder.
func (c *stateObject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, c.data)
}

// setError remembers the first non-nil error it is called with.
func (self *stateObject) setError(err error) {
	if self.dbErr == nil {
		self.dbErr = err
	}
}

func (self *stateObject) markSuicided() {
	self.suicided = true
}

func (c *stateObject) touch() {
	c.db.journal.append(touchChange{
		account: &c.address,
	})
	if c.address == ripemd {
		// Explicitly put it in the dirty-cache, which is otherwise generated from
		// flattened journals.
		c.db.journal.dirty(c.address)
	}
}

func (c *stateObject) getTrie(db Database) Trie {
	if c.trie == nil {
		var err error
		c.trie, err = db.OpenStorageTrie(c.addrHash, c.data.Root)
		if err != nil {
			c.trie, _ = db.OpenStorageTrie(c.addrHash, common.Hash{})
			c.setError(fmt.Errorf("can't create storage trie: %v", err))
		}
	}
	return c.trie
}

// GetState retrieves a value from the account storage trie.
func (self *stateObject) GetState(db Database, key []byte) []byte {
	// If we have a dirty value for this state entry, return it
	value, dirty := self.dirtyStorage[string(key)]
	if dirty {
		return value
	}
	// Otherwise return the entry's original value
	return self.GetCommittedState(db, key)
}

func (self *stateObject) getCommittedStateCache(key []byte) []byte {
	value, cached := self.originStorage[string(key)]
	if cached {
		return value
	}

	self.db.refLock.Lock()
	parentDB := self.db.parent
	parentCommitted := self.db.parentCommitted
	refLock := &self.db.refLock

	for parentDB != nil {
		value := parentDB.getStateObjectSnapshot(self.address, key)
		if value != nil {
			self.originStorage[string(key)] = value
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
func (self *stateObject) GetCommittedState(db Database, key []byte) []byte {
	// If we have the original value cached, return that
	if value := self.getCommittedStateCache(key); len(value) != 0 {
		//log.Trace("GetCommittedState cache", "key", hex.EncodeToString(key), "value", len(value))
		return value
	}

	// Otherwise load the valueKey from trie
	enc, err := self.getTrie(db).TryGet(key[:])
	if err != nil {
		self.setError(err)
		return []byte{}
	}
	value := make([]byte, 0)
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			self.setError(err)
		}
		value = content
	}

	//log.Trace("GetCommittedState trie", "key", hex.EncodeToString(key), "value", len(value))
	self.originStorage[string(key)] = value
	return value
}

// SetState updates a value in account storage.
// set [prefixKey,value] to storage
func (self *stateObject) SetState(db Database, key, value []byte) {
	//if the new value is the same as old,don't set
	preValue := self.GetState(db, key)
	if bytes.Equal(preValue, value) {
		return
	}

	//New value is different, update and journal the change
	self.db.journal.append(storageChange{
		account:  &self.address,
		key:      key,
		preValue: preValue,
	})

	self.setState(key, value)
}

func (self *stateObject) setState(key []byte, value []byte) {
	cpy := make([]byte, len(value))
	copy(cpy, value)
	self.dirtyStorage[string(key)] = cpy
}

func (self *stateObject) getPrefixValue(pack, key, value []byte) []byte {
	// Empty value deleted on updateTrie
	if len(value) == 0 {
		return []byte{}
	}
	// Ensure the same Value, unique in the same trie and different trie values
	//prefix := append(self.data.StorageKeyPrefix, key...)
	prefix := append(self.data.StorageKeyPrefix, pack...)
	prefix = append(prefix, key...)
	prefixHash := common.Hash{}
	keccak := sha3.NewKeccak256()
	keccak.Write(prefix)
	keccak.Sum(prefixHash[:0])
	return append(prefixHash[:], value...)
}

func (self *stateObject) removePrefixValue(value []byte) []byte {
	if len(value) > common.HashLength {
		return value[common.HashLength:]
	}
	return []byte{}
}

// updateTrie writes cached storage modifications into the object's storage trie.
func (self *stateObject) updateTrie(db Database) Trie {
	tr := self.getTrie(db)
	for key, value := range self.dirtyStorage {
		delete(self.dirtyStorage, key)

		// Skip noop changes, persist actual changes
		oldValue := self.originStorage[key]
		if bytes.Equal(value, oldValue) {
			continue
		}

		self.originStorage[key] = value

		if len(value) == 0 {
			self.setError(tr.TryDelete([]byte(key)))
			continue
		}

		// Encoding []byte cannot fail, ok to ignore the error.
		v, _ := rlp.EncodeToBytes(value)
		self.setError(tr.TryUpdate([]byte(key), v))
	}

	return tr
}

// UpdateRoot sets the trie root to the current root hash of
func (self *stateObject) updateRoot(db Database) {
	self.updateTrie(db)
	//self.data.Root = self.trie.Hash()
	self.data.Root = self.trie.ParallelHash()
}

// CommitTrie the storage trie of the object to db.
// This updates the trie root.
func (self *stateObject) CommitTrie(db Database) error {
	self.updateTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}

	root, err := self.trie.Commit(nil)

	if err == nil {
		self.data.Root = root
	}
	return err
}

// AddBalance removes amount from c's balance.
// It is used to add funds to the destination account of a transfer.
func (c *stateObject) AddBalance(amount *big.Int) {
	// EIP158: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if c.empty() {
			c.touch()
		}

		return
	}
	c.SetBalance(new(big.Int).Add(c.Balance(), amount))
}

// SubBalance removes amount from c's balance.
// It is used to remove funds from the origin account of a transfer.
func (c *stateObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	c.SetBalance(new(big.Int).Sub(c.Balance(), amount))
}

func (self *stateObject) SetBalance(amount *big.Int) {
	self.db.journal.append(balanceChange{
		account: &self.address,
		prev:    new(big.Int).Set(self.data.Balance),
	})
	self.setBalance(amount)
}

func (self *stateObject) setBalance(amount *big.Int) {
	self.data.Balance = amount
}

// Return the gas back to the origin. Used by the Virtual machine or Closures
func (c *stateObject) ReturnGas(gas *big.Int) {}

func (self *stateObject) deepCopy(db *StateDB) *stateObject {
	stateObject := newObject(db, self.address, self.data)
	if self.trie != nil {
		stateObject.trie = db.db.CopyTrie(self.trie)
	}
	stateObject.code = self.code
	stateObject.dirtyStorage = self.dirtyStorage.Copy()
	stateObject.originStorage = self.originStorage.Copy()
	stateObject.suicided = self.suicided
	stateObject.dirtyCode = self.dirtyCode
	stateObject.deleted = self.deleted
	return stateObject
}

// Copy account status, recreate trie
func (self *stateObject) copy(db *StateDB) *stateObject {
	stateObject := newObject(db, self.address, self.data)
	if self.trie != nil {
		stateObject.trie = db.db.NewTrie(self.trie)
	}
	stateObject.code = self.code
	stateObject.suicided = self.suicided
	stateObject.dirtyCode = self.dirtyCode
	stateObject.deleted = self.deleted
	return stateObject
}

//
// Attribute accessors
//

// Returns the address of the contract/account
func (c *stateObject) Address() common.Address {
	return c.address
}

// Code returns the contract code associated with this object, if any.
func (self *stateObject) Code(db Database) []byte {
	if self.code != nil {
		return self.code
	}
	if bytes.Equal(self.CodeHash(), emptyCodeHash) {
		return nil
	}
	code, err := db.ContractCode(self.addrHash, common.BytesToHash(self.CodeHash()))
	if err != nil {
		self.setError(fmt.Errorf("can't load code hash %x: %v", self.CodeHash(), err))
	}
	self.code = code
	return code
}

func (self *stateObject) SetCode(codeHash common.Hash, code []byte) {
	prevcode := self.Code(self.db.db)
	self.db.journal.append(codeChange{
		account:  &self.address,
		prevhash: self.CodeHash(),
		prevcode: prevcode,
	})
	self.setCode(codeHash, code)
}

func (self *stateObject) setCode(codeHash common.Hash, code []byte) {
	self.code = code
	self.data.CodeHash = codeHash[:]
	self.dirtyCode = true
}

func (self *stateObject) SetNonce(nonce uint64) {
	self.db.journal.append(nonceChange{
		account: &self.address,
		prev:    self.data.Nonce,
	})
	self.setNonce(nonce)
}

func (self *stateObject) setNonce(nonce uint64) {
	self.data.Nonce = nonce
}

func (self *stateObject) CodeHash() []byte {
	return self.data.CodeHash
}

func (self *stateObject) Balance() *big.Int {
	return self.data.Balance
}

func (self *stateObject) Nonce() uint64 {
	return self.data.Nonce
}

// Never called, but must be present to allow stateObject to be used
// as a vm.Account interface that also satisfies the vm.ContractRef
// interface. Interfaces are awesome.
func (self *stateObject) Value() *big.Int {
	panic("Value on stateObject should never be called")
}
