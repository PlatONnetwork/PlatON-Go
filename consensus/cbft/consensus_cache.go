package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"sync"
	"errors"
)

var (
	errMakeStateDB = errors.New("make StateDB error")
)

type Cache struct {
	stateDBCache		map[common.Hash]*stateDBCache		// key is header SealHash
	receiptsCache		map[common.Hash]*receiptsCache		// key is header SealHash
	chain  				*core.BlockChain
	stateDBMu       	sync.RWMutex
	receiptsMu			sync.RWMutex
}

type stateDBCache struct {
	stateDB *state.StateDB
	blockNum uint64
}

type receiptsCache struct {
	receipts []*types.Receipt
	blockNum uint64
}

func NewCache(blockChain *core.BlockChain) *Cache {
	cache := &Cache{
		stateDBCache:        make(map[common.Hash]*stateDBCache),
		receiptsCache:       make(map[common.Hash]*receiptsCache),
		chain: blockChain,
	}
	return cache
}

// Read the Receipt collection from the cache map.
func (c *Cache) ReadReceipts(sealHash common.Hash) []*types.Receipt {
	c.receiptsMu.RLock()
	defer c.receiptsMu.RUnlock()
	if obj, exist := c.receiptsCache[sealHash]; exist {
		return obj.receipts
	}
	return nil
}

// Read the StateDB instance from the cache map
func (c *Cache) ReadStateDB(sealHash common.Hash) *state.StateDB {
	c.stateDBMu.RLock()
	defer c.stateDBMu.RUnlock()
	log.Info("Read the StateDB instance from the cache map","sealHash", sealHash)
	if obj, exist := c.stateDBCache[sealHash]; exist {
		return obj.stateDB.Copy()
	}
	return nil
}

// Write Receipt to the cache
func (c *Cache) WriteReceipts(sealHash common.Hash, receipts []*types.Receipt, blockNum uint64) {
	c.receiptsMu.Lock()
	defer c.receiptsMu.Unlock()
	obj, exist := c.receiptsCache[sealHash]
	if exist && obj.blockNum == blockNum {
		obj.receipts = append(obj.receipts, receipts...)
	} else if !exist {
		c.receiptsCache[sealHash] = &receiptsCache{receipts: receipts, blockNum: blockNum}
	}
}

// Write a StateDB instance to the cache
func (c *Cache) WriteStateDB(sealHash common.Hash, stateDB *state.StateDB, blockNum uint64) {
	c.stateDBMu.Lock()
	defer c.stateDBMu.Unlock()
	log.Info("Write a StateDB instance to the cache", "sealHash", sealHash, "blockNum", blockNum)
	if _, exist := c.stateDBCache[sealHash]; !exist {
		c.stateDBCache[sealHash] = &stateDBCache{stateDB: stateDB, blockNum: blockNum}
	}
}

// Read the Receipt collection from the cache map
func (c *Cache) clearReceipts(sealHash common.Hash) {
	c.receiptsMu.Lock()
	defer c.receiptsMu.Unlock()

	var blockNum uint64
	if obj, exist := c.receiptsCache[sealHash]; exist {
		blockNum = obj.blockNum
		//delete(c.receiptsCache, sealHash)
	}
	for hash, obj := range c.receiptsCache {
		if obj.blockNum <= blockNum {
			delete(c.receiptsCache, hash)
		}
	}
}

// Read the StateDB instance from the cache map
func (c *Cache) clearStateDB(sealHash common.Hash) {
	c.stateDBMu.Lock()
	defer c.stateDBMu.Unlock()

	var blockNum uint64
	if obj, exist := c.stateDBCache[sealHash]; exist {
		blockNum = obj.blockNum
		//delete(c.stateDBCache, sealHash)
	}
	for hash, obj := range c.stateDBCache {
		if obj.blockNum <= blockNum {
			delete(c.stateDBCache, hash)
		}
	}
}

// Get the StateDB instance of the corresponding block
func (c *Cache) MakeStateDB(block *types.Block) (*state.StateDB, error) {
	// Create a StateDB instance from the blockchain based on stateRoot
	if state, err := c.chain.StateAt(block.Root()); err == nil && state != nil {
		return state, nil
	}
	// Read and copy the stateDB instance in the cache
	sealHash := c.chain.Engine().SealHash(block.Header())
	log.Info("Read and copy the stateDB instance in the cache","sealHash", sealHash, "blockHash", block.Hash(), "blockNum", block.NumberU64())
	if state := c.ReadStateDB(sealHash); state != nil {
		//return state.Copy(), nil
		return state, nil
	} else {
		return nil, errMakeStateDB
	}
}

// Get the StateDB instance of the corresponding block
func (c *Cache) ClearCache(block *types.Block) {
	sealHash := c.chain.Engine().SealHash(block.Header())
	c.clearReceipts(sealHash)
	c.clearStateDB(sealHash)
}