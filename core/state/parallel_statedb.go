package state

import (
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var (
	parallelLocker sync.Mutex
)

func (self *StateDB) GetOrNewParallelStateObject(addr common.Address) *ParallelStateObject {
	stateObject := self.justGetStateObject(addr)
	if stateObject == nil || stateObject.deleted {
		log.Debug("Cannot find stateObject in Parallel", "addr", addr.Hex(), "isNil", stateObject == nil)
		stateObject, _ = self.justCreateObject(addr)
	}
	return NewParallelStateObject(stateObject)
}

func (self *StateDB) justGetStateObject(addr common.Address) (stateObject *stateObject) {
	if obj := self.justGetStateObjectCache(addr); obj != nil {
		if obj.deleted {
			return nil
		}
		return obj
	}
	// Load the object from the database.
	parallelLocker.Lock()
	enc, err := self.trie.TryGet(addr[:])
	parallelLocker.Unlock()
	if len(enc) == 0 {
		self.setError(err)
		return nil
	}
	var data Account
	if err := rlp.DecodeBytes(enc, &data); err != nil {
		log.Error("Failed to decode state object", "addr", addr, "err", err)
		return nil
	}
	obj := newObject(self, addr, data)
	//do not set to state.stateObjects.
	//self.setStateObject(obj)
	return obj
}

func (self *StateDB) justGetStateObjectCache(addr common.Address) (stateObject *stateObject) {
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
			//do not set to state.stateObjects.
			//self.setStateObject(cpy)
			return cpy
		} else if parentCommitted {
			refLock.Unlock()
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

func (self *StateDB) justCreateObject(addr common.Address) (newobj, prev *stateObject) {
	prev = self.justGetStateObject(addr)
	newobj = newObject(self, addr, Account{})
	newobj.setNonce(0)
	return newobj, prev
}
