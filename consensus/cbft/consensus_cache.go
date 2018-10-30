package cbft

import (
	"Platon-go/common"
	"Platon-go/core"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"sync"
	"errors"
)

var (
	errMakeStateDB = errors.New("make StateDB error")
)

type Cache struct {
	stateDBCache	map[common.Hash]*state.StateDB		// key is header stateRoot
	receiptsCache	map[common.Hash][]*types.Receipt	// key is header hash
	chain  			*core.BlockChain
	stateDBMu       sync.RWMutex
	receiptsMu		sync.RWMutex
}

func NewCache(blockChain *core.BlockChain) *Cache {
	cache := &Cache{
		stateDBCache:        make(map[common.Hash]*state.StateDB),
		receiptsCache:       make(map[common.Hash][]*types.Receipt),
		chain: blockChain,
	}
	return cache
}

// 从缓存map中读取Receipt集合
func (c *Cache) ReadReceipts(sealhash common.Hash) []*types.Receipt {
	c.receiptsMu.RLock()
	defer c.receiptsMu.RUnlock()
	if receipts, exist := c.receiptsCache[sealhash]; exist {
		return receipts
	}
	return nil
}

// 从缓存map中读取StateDB实例
func (c *Cache) ReadStateDB(stateRoot common.Hash) *state.StateDB {
	c.stateDBMu.RLock()
	defer c.stateDBMu.RUnlock()
	if stateDB, exist := c.stateDBCache[stateRoot]; exist {
		return stateDB
	}
	return nil
}

// 将Receipt写入缓存
func (c *Cache) WriteReceipts(blockHash common.Hash, receipts []*types.Receipt) {
	c.receiptsMu.Lock()
	defer c.receiptsMu.Unlock()
	if _receipts, exist := c.receiptsCache[blockHash]; exist {
		_receipts := append(_receipts, receipts...)
		c.receiptsCache[blockHash] = _receipts
	} else {
		c.receiptsCache[blockHash] = receipts
	}
}

// 将StateDB实例写入缓存
func (c *Cache) WriteStateDB(stateRoot common.Hash, stateDB *state.StateDB) {
	c.stateDBMu.Lock()
	defer c.stateDBMu.Unlock()
	c.stateDBCache[stateRoot] = stateDB
}

// 从缓存map中读取Receipt集合
func (c *Cache) clearReceipts(blockHash common.Hash) {
	c.receiptsMu.Lock()
	defer c.receiptsMu.Unlock()
	delete(c.receiptsCache, blockHash)
}

// 从缓存map中读取StateDB实例
func (c *Cache) clearStateDB(stateRoot common.Hash) {
	c.stateDBMu.Lock()
	defer c.stateDBMu.Unlock()
	delete(c.stateDBCache, stateRoot)
}

// 获取相应block的StateDB实例
func (c *Cache) MakeStateDB(block *types.Block) (*state.StateDB, error) {
	// 基于stateRoot从blockchain中创建StateDB实例
	if state, err := c.chain.StateAt(block.Root()); err == nil && state != nil {
		return state, nil
	}
	// 读取并拷贝缓存中StateDB实例
	if state := c.ReadStateDB(block.Root()); state != nil {
		return state.Copy(), nil
	} else {
		return nil, errMakeStateDB
	}
}

// 获取相应block的StateDB实例
func (c *Cache) ClearCache(block *types.Block) {
	c.clearReceipts(block.Hash())
	c.clearStateDB(block.Root())
}