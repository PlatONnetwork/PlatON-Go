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
	"encoding/hex"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"io"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"

	//"runtime/debug"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var emptyCodeHash = crypto.Keccak256(nil)

type Code []byte
type Abi []byte

func (self Code) String() string {
	return string(self) //strings.Join(Disassemble(self), " ")
}

type Storage map[string]common.Hash
type ValueStorage map[common.Hash][]byte
type ReferenceValueStorage map[common.Hash]*ReferenceValue

type ReferenceValue struct {
	Count  uint32
	Value  []byte
}


func (self ReferenceValueStorage) deleleReferenceValue(valueKey common.Hash) {
	if refValue, ok := self[valueKey]; ok {
		// Because the reference of a value may be 0 before there is no fork,
		// `Count--` may be negative after the fork.
		if refValue.Count--; refValue.Count <= 0 {
			delete(self, valueKey)
		}else {
			self[valueKey] = refValue
		}
	}
}

func (self ReferenceValueStorage) pureDelReferenceValue(valueKey common.Hash) {
	delete(self, valueKey)
}


func (self ReferenceValueStorage) setReferenceValue (valueKey common.Hash, val []byte) {

	var ref *ReferenceValue
	if refValue, ok := self[valueKey]; ok {
		refValue.Count++
		ref = refValue
	}else {
		ref = &ReferenceValue{
			Count: 1,
			Value: val,
		}
	}
	self[valueKey] = ref
}

func (self ReferenceValueStorage) pureSetReferenceValue (valueKey common.Hash, val []byte) {
	ref := &ReferenceValue{
		Count: 1,
		Value: val,
	}
	self[valueKey] = ref
}

func (self ReferenceValueStorage) Copy() ReferenceValueStorage {
	cpy := make(ReferenceValueStorage, len(self))
	for key, value := range self {
		// deep copy
		v := make([]byte, len(value.Value))
		copy(v, value.Value)
		cpv := &ReferenceValue{
			Count: value.Count,
			Value: v,
		}
		cpy[key] = cpv
	}

	return cpy
}

// Storage -> hash : hash , common.Hash ([32]byte)
//type Storage map[common.Hash]common.Hash

func (self Storage) String() (str string) {
	for key, value := range self {
		// %X -> Provide hexadecimal
		str += fmt.Sprintf("%X : %X\n", key, value)
	}

	return
}

// Copy a copy of Storage
func (self Storage) Copy() Storage {
	cpy := make(Storage, len(self))
	for key, value := range self {
		cpy[key] = value
	}

	return cpy
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

	abi Abi

	originStorage      Storage      // Storage cache of original entries to dedup rewrites
	originValueStorage ValueStorage // Storage cache of original entries to dedup rewrites

	dirtyStorage      Storage      // Storage entries that need to be flushed to disk
	dirtyValueStorage ReferenceValueStorage // Storage entries that need to be flushed to disk

	// Cache flags.
	// When an object is marked suicided it will be delete from the trie
	// during the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
}

// empty returns whether the account is considered empty.
func (s *stateObject) empty() bool {
	if vm.IsPlatONPrecompiledContract(s.address) {
		return false
	}
	return s.data.Nonce == 0 && s.data.Balance.Sign() == 0 && bytes.Equal(s.data.CodeHash, emptyCodeHash)
}

func (s *stateObject) delReferenceValueStorage(valueKey common.Hash) {
	if checkForkPIP0_11_0(s.db) {
		log.Trace("delReferenceValueStorage Start, deleleReferenceValue", "govVersion", s.db.govVersion, "valueKey", valueKey.String())
		s.dirtyValueStorage.deleleReferenceValue(valueKey)
	}else {
		s.dirtyValueStorage.pureDelReferenceValue(valueKey)
	}
}

func (s *stateObject) setReferenceValueStorage (valueKey common.Hash, val []byte) {
	if checkForkPIP0_11_0(s.db) {
		log.Trace("setReferenceValueStorage Start, setReferenceValue", "govVersion", s.db.govVersion, "valueKey", valueKey.String(), "val", hexutil.Encode(val))
		s.dirtyValueStorage.setReferenceValue(valueKey, val)
	}else {
		s.dirtyValueStorage.pureSetReferenceValue(valueKey, val)
	}
}

func checkForkPIP0_11_0(state *StateDB) bool {
	if state.govVersion >= plugin.FORKVERSION_0_11_0 {
		return true
	} else {
		return false
	}
}

// Account is the Ethereum consensus representation of accounts.
// These objects are stored in the main account trie.
type Account struct {
	Nonce    uint64
	Balance  *big.Int
	Root     common.Hash // merkle root of the storage trie
	CodeHash []byte
	AbiHash  []byte
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
		db:       db,
		address:  address,
		addrHash: crypto.Keccak256Hash(address[:]),
		data:     data,

		originStorage:      make(Storage),
		originValueStorage: make(ValueStorage),

		dirtyStorage:      make(Storage),
		dirtyValueStorage: make(ReferenceValueStorage),
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
//func (self *stateObject) GetState(db Database, key common.Hash) common.Hash {
//	// If we have a dirty value for this state entry, return it
//	value, dirty := self.dirtyStorage[key]
//	if dirty {
//		return value
//	}
//	// Otherwise return the entry's original value
//	return self.GetCommittedState(db, key)
//}

// GetState retrieves a value from the account storage trie.
func (self *stateObject) GetState(db Database, keyTree string) []byte {
	// If we have a dirty value for this state entry, return it
	valueKey, dirty := self.dirtyStorage[keyTree]
	if dirty {
		refValue, ok := self.dirtyValueStorage[valueKey]
		if ok {
			return refValue.Value
		}
	}
	// Otherwise return the entry's original value
	return self.GetCommittedState(db, keyTree)
}

// GetCommittedState retrieves a value from the committed account storage trie.
//func (self *stateObject) GetCommittedState(db Database, key common.Hash) common.Hash {
//	// If we have the original value cached, return that
//	value, cached := self.originStorage[key]
//	if cached {
//		return value
//	}
//	// Otherwise load the value from the database
//	enc, err := self.getTrie(db).TryGet(key[:])
//	if err != nil {
//		self.setError(err)
//		return common.Hash{}
//	}
//	if len(enc) > 0 {
//		_, content, _, err := rlp.Split(enc)
//		if err != nil {
//			self.setError(err)
//		}
//		value.SetBytes(content)
//	}
//	self.originStorage[key] = value
//	return value
//}
func (self *stateObject) getCommittedStateCache(key string) []byte {
	valueKey, cached := self.originStorage[key]
	if cached {
		value, cached2 := self.originValueStorage[valueKey]
		if cached2 {
			return value
		}
	}

	self.db.refLock.Lock()
	parentDB := self.db.parent
	parentCommitted := self.db.parentCommitted
	refLock := &self.db.refLock

	for parentDB != nil {
		valueKey, value := parentDB.getStateObjectSnapshot(self.address, key)
		if value != nil {
			self.originStorage[key] = valueKey
			self.originValueStorage[valueKey] = value
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
func (self *stateObject) GetCommittedState(db Database, key string) []byte {
	value := make([]byte, 0)
	valueKey := common.Hash{}
	// If we have the original value cached, return that
	if value := self.getCommittedStateCache(key); value != nil {
		log.Trace("GetCommittedState cache", "key", hex.EncodeToString([]byte(key)), "value", len(value))
		return value
	}

	// Otherwise load the valueKey from trie
	enc, err := self.getTrie(db).TryGet([]byte(key))
	if err != nil {
		self.setError(err)
		return []byte{}
	}
	if len(enc) > 0 {
		_, content, _, err := rlp.Split(enc)
		if err != nil {
			self.setError(err)
		}
		valueKey.SetBytes(content)

		//load value from db
		value = self.db.trie.GetKey(valueKey.Bytes())
		if err != nil {
			self.setError(err)
		}
	}

	if valueKey != emptyStorage && len(value) == 0 {
		return []byte{}
	}
	if len(value) == 0 && valueKey == emptyStorage {
		log.Trace("Empty storage valuekey", "key", hex.EncodeToString([]byte(key)), "valueKey", valueKey.String())
	}
	log.Trace("GetCommittedState", "key", hex.EncodeToString([]byte(key)), "valueKey", valueKey.String(), "value", len(value))
	self.originStorage[key] = valueKey
	self.originValueStorage[valueKey] = value
	return value
}

// SetState updates a value in account storage.
// set [keyTrie,valueKey] to storage
// set [valueKey,value] to db
func (self *stateObject) SetState(db Database, keyTrie string, valueKey common.Hash, value []byte) {
	//if the new value is the same as old,don't set
	preValue := self.GetState(db, keyTrie) // get value key
	if bytes.Equal(preValue, value) {
		return
	}

	//New value is different, update and journal the change
	self.db.journal.append(storageChange{
		account:  &self.address,
		key:      keyTrie,
		valueKey: self.originStorage[keyTrie],
		preValue: preValue,
	})

	self.setState(keyTrie, valueKey, value)
}

func (self *stateObject) setState(key string, valueKey common.Hash, value []byte) {
	cpy := make([]byte, len(value))
	copy(cpy, value)

	// delete value storage
	if dirtyValue, ok := self.dirtyStorage[key]; ok {
		//delete(self.dirtyValueStorage, dirtyValue)
		self.delReferenceValueStorage(dirtyValue)
	}

	self.dirtyStorage[key] = valueKey
	//self.dirtyValueStorage[valueKey] = cpy
	self.setReferenceValueStorage(valueKey, cpy)
}

// updateTrie writes cached storage modifications into the object's storage trie.
func (self *stateObject) updateTrie(db Database) Trie {
	tr := self.getTrie(db)
	for key, valueKey := range self.dirtyStorage {
		delete(self.dirtyStorage, key)

		if valueKey == self.originStorage[key] {
			//delete(self.dirtyValueStorage, valueKey)
			self.delReferenceValueStorage(valueKey)
			continue
		}

		if valueKey == emptyStorage || (valueKey == common.Hash{}) {
			//delete(self.dirtyValueStorage, valueKey)
			self.delReferenceValueStorage(valueKey)

			if oldValueKey, ok := self.originStorage[key]; ok {
				delete(self.originValueStorage, oldValueKey)
			}
			self.originStorage[key] = valueKey
			self.originValueStorage[valueKey] = []byte{}

			self.setError(tr.TryDelete([]byte(key)))
			continue
		}

		self.originStorage[key] = valueKey

		v, _ := rlp.EncodeToBytes(bytes.TrimLeft(valueKey[:], "\x00"))
		self.setError(tr.TryUpdate([]byte(key), v))

		//flush dirty value
		if refValue, ok := self.dirtyValueStorage[valueKey]; ok {
			//delete(self.dirtyValueStorage, valueKey)
			self.delReferenceValueStorage(valueKey)
			self.originValueStorage[valueKey] = refValue.Value
		}
	}

	return tr
}

// UpdateRoot sets the trie root to the current root hash of
func (self *stateObject) updateRoot(db Database) {
	self.updateTrie(db)
	self.data.Root = self.trie.Hash()
}

// CommitTrie the storage trie of the object to db.
// This updates the trie root.
func (self *stateObject) CommitTrie(db Database) error {
	self.updateTrie(db)
	if self.dbErr != nil {
		return self.dbErr
	}

	root, err := self.trie.Commit(func(leaf []byte, parent common.Hash) error {
		var valueKey common.Hash
		_, content, _, err := rlp.Split(leaf)
		if err != nil {
			self.setError(err)
		}
		valueKey.SetBytes(content)
		if value, ok := self.originValueStorage[valueKey]; ok {
			self.db.db.TrieDB().InsertBlob(valueKey, value)
		}

		self.db.db.TrieDB().Reference(valueKey, parent)
		return nil
	})

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
	stateObject.dirtyValueStorage = self.dirtyValueStorage.Copy()
	stateObject.originStorage = self.originStorage.Copy()
	stateObject.originValueStorage = self.originValueStorage.Copy()
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

// todo: New method
// ======================================= New method ===============================

// todo: new method -> AbiHash
func (self *stateObject) AbiHash() []byte {
	return self.data.AbiHash
}

// ABI returns the contract abi associated with this object, if any.
func (self *stateObject) Abi(db Database) []byte {
	//if self.Abi != nil {
	//	return self.abi
	//}
	if bytes.Equal(self.AbiHash(), emptyCodeHash) {
		return nil
	}
	// Extract the code from the tree, enter the parameters: address and hash, here you need to find the acquisition rules in depth
	abi, err := db.ContractAbi(self.addrHash, common.BytesToHash(self.AbiHash()))
	if err != nil {
		self.setError(fmt.Errorf("can't load abi hash %x: %v", self.AbiHash(), err))
	}
	self.abi = abi
	return abi
}

// todo: new method -> SetAbi.
func (self *stateObject) SetAbi(abiHash common.Hash, abi []byte) {
	prevabi := self.Abi(self.db.db)
	self.db.journal.append(abiChange{
		account:  &self.address,
		prevhash: self.AbiHash(),
		prevabi:  prevabi,
	})
	self.setAbi(abiHash, abi)
}

// todo: new method -> setAbi
func (self *stateObject) setAbi(abiHash common.Hash, abi []byte) {
	self.abi = abi
	self.data.AbiHash = abiHash[:]
}
