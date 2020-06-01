package state

import (
	"math/big"
	"sync"
	"time"

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
		return self.justCreateObject(addr)
	}
	return NewParallelStateObject(stateObject, false)
}

func (self *StateDB) justGetStateObject(addr common.Address) (stateObject *stateObject) {
	if obj := self.justGetStateObjectCache(addr); obj != nil {
		if obj.deleted {
			return nil
		}
		return obj
	}
	// Load the object from the database.
	start := time.Now()
	parallelLocker.Lock()
	if start.Add(20 * time.Millisecond).Before(time.Now()) {
		log.Trace("Get parallelLocker overtime", "address", addr.String(), "duration", time.Since(start))
	}
	start = time.Now()
	enc, err := self.trie.TryGet(addr[:])
	if start.Add(20 * time.Millisecond).Before(time.Now()) {
		log.Trace("Trie tryGet overtime", "address", addr.String(), "duration", time.Since(start))
	}
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

func (self *StateDB) justCreateObject(addr common.Address) *ParallelStateObject {
	//newobj := newObject(self, addr, Account{})
	newobj := newObject(self, addr, Account{StorageKeyPrefix: addr.Bytes()})
	//self.journal.append(createObjectChange{account: &addr})
	newobj.setNonce(0)
	return &ParallelStateObject{
		stateObject: newobj,
		prevAmount:  big.NewInt(0),
		createFlag:  true,
	}
	//return self.createObject(addr)
}
