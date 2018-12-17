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
	stateDBCache		map[common.Hash]*stateDBCache		// key is header stateRoot
	receiptsCache		map[common.Hash]*receiptsCache		// key is header hash
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
func (c *Cache) ReadReceipts(blockHash common.Hash) []*types.Receipt {
	c.receiptsMu.RLock()
	defer c.receiptsMu.RUnlock()
	if obj, exist := c.receiptsCache[blockHash]; exist {
		return obj.receipts
	}
	return nil
}

// Read the StateDB instance from the cache map.
func (c *Cache) ReadStateDB(stateRoot common.Hash) *state.StateDB {
	c.stateDBMu.RLock()
	defer c.stateDBMu.RUnlock()
	if obj, exist := c.stateDBCache[stateRoot]; exist {
		return obj.stateDB
	}
	return nil
}

// Write Receipt to the cache.
func (c *Cache) WriteReceipts(blockHash common.Hash, receipts []*types.Receipt, blockNum uint64) {
	c.receiptsMu.Lock()
	defer c.receiptsMu.Unlock()
	obj, exist := c.receiptsCache[blockHash]
	if exist && obj.blockNum == blockNum {
		obj.receipts = append(obj.receipts, receipts...)
	} else if !exist {
		c.receiptsCache[blockHash] = &receiptsCache{receipts: receipts, blockNum: blockNum}
	}
}

// Write StateDB instance to the cache.
func (c *Cache) WriteStateDB(stateRoot common.Hash, stateDB *state.StateDB, blockNum uint64) {
	c.stateDBMu.Lock()
	defer c.stateDBMu.Unlock()
	if _, exist := c.stateDBCache[stateRoot]; !exist {
		c.stateDBCache[stateRoot] = &stateDBCache{stateDB: stateDB, blockNum: blockNum}
	}
}

// Read the Receipt collection from the cache map.
func (c *Cache) clearReceipts(blockHash common.Hash) {
	c.receiptsMu.Lock()
	defer c.receiptsMu.Unlock()

	var blockNum uint64
	if obj, exist := c.receiptsCache[blockHash]; exist {
		blockNum = obj.blockNum
		//delete(c.receiptsCache, blockHash)
	}
	for hash, obj := range c.receiptsCache {
		if obj.blockNum <= blockNum {
			delete(c.receiptsCache, hash)
		}
	}
}

// Read the statedb instance from the cache map.
func (c *Cache) clearStateDB(stateRoot common.Hash) {
	c.stateDBMu.Lock()
	defer c.stateDBMu.Unlock()

	var blockNum uint64
	if obj, exist := c.stateDBCache[stateRoot]; exist {
		blockNum = obj.blockNum
		//delete(c.stateDBCache, stateRoot)
	}
	for hash, obj := range c.stateDBCache {
		if obj.blockNum <= blockNum {
			delete(c.stateDBCache, hash)
		}
	}
}

// Get the StateDB instance of the corresponding block.
func (c *Cache) MakeStateDB(block *types.Block) (*state.StateDB, error) {
	// Create a StateDB instance from the blockchain based on stateRoot.
	if state, err := c.chain.StateAt(block.Root()); err == nil && state != nil {
		return state, nil
	}
	// Read and copy the stateDB instance in the cache.
	log.Info("~ Read and copy the stateDB instance in the cache.", "stateRoot", block.Root())
	if state := c.ReadStateDB(block.Root()); state != nil {
		return state.Copy(), nil
	} else {
		return nil, errMakeStateDB
	}
}

// Get the StateDB instance of the corresponding block.
func (c *Cache) ClearCache(block *types.Block) {
	c.clearReceipts(block.Hash())
	c.clearStateDB(block.Root())
}