package cbft

import (
	"Platon-go/common"
	"Platon-go/core"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"Platon-go/log"
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
	stateDB state.StateDB
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

// 从缓存map中读取Receipt集合
func (c *Cache) ReadReceipts(blockHash common.Hash) []*types.Receipt {
	c.receiptsMu.RLock()
	defer c.receiptsMu.RUnlock()
	if obj, exist := c.receiptsCache[blockHash]; exist {
		return obj.receipts
	}
	return nil
}

// 从缓存map中读取StateDB实例
func (c *Cache) ReadStateDB(stateRoot common.Hash) *state.StateDB {
	c.stateDBMu.RLock()
	defer c.stateDBMu.RUnlock()
	log.Info("从缓存map中读取StateDB实例", "stateRoot", stateRoot)
	if obj, exist := c.stateDBCache[stateRoot]; exist {
		state := obj.stateDB
		return &state
	}
	return nil
}

// 将Receipt写入缓存
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

// 将StateDB实例写入缓存
func (c *Cache) WriteStateDB(stateRoot common.Hash, stateDB state.StateDB, blockNum uint64) {
	c.stateDBMu.Lock()
	defer c.stateDBMu.Unlock()
	log.Info("将StateDB实例写入缓存", "stateRoot", stateRoot, "blockNum", blockNum)
	if _, exist := c.stateDBCache[stateRoot]; !exist {
		c.stateDBCache[stateRoot] = &stateDBCache{stateDB: stateDB, blockNum: blockNum}
	}
}

// 从缓存map中读取Receipt集合
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

// 从缓存map中读取StateDB实例
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

// 获取相应block的StateDB实例
func (c *Cache) MakeStateDB(block *types.Block) (*state.StateDB, error) {
	// 基于stateRoot从blockchain中创建StateDB实例
	if state, err := c.chain.StateAt(block.Root()); err == nil && state != nil {
		return state, nil
	}
	// 读取并拷贝缓存中StateDB实例
	log.Info("读取并拷贝缓存中StateDB实例", "stateRoot", block.Root())
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