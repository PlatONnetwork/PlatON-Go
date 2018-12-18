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

// 从缓存map中读取Receipt集合
func (c *Cache) ReadReceipts(sealHash common.Hash) []*types.Receipt {
	c.receiptsMu.RLock()
	defer c.receiptsMu.RUnlock()
	if obj, exist := c.receiptsCache[sealHash]; exist {
		return obj.receipts
	}
	return nil
}

// 从缓存map中读取StateDB实例
func (c *Cache) ReadStateDB(sealHash common.Hash) *state.StateDB {
	c.stateDBMu.RLock()
	defer c.stateDBMu.RUnlock()
	log.Info("从缓存map中读取StateDB实例","sealHash", sealHash)
	if obj, exist := c.stateDBCache[sealHash]; exist {
		//return obj.stateDB.Copy()
		return obj.stateDB
	}
	return nil
}

// 将Receipt写入缓存
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

// 将StateDB实例写入缓存
func (c *Cache) WriteStateDB(sealHash common.Hash, stateDB *state.StateDB, blockNum uint64) {
	c.stateDBMu.Lock()
	defer c.stateDBMu.Unlock()
	log.Info("将StateDB实例写入缓存", "sealHash", sealHash, "blockNum", blockNum)
	if _, exist := c.stateDBCache[sealHash]; !exist {
		c.stateDBCache[sealHash] = &stateDBCache{stateDB: stateDB, blockNum: blockNum}
	}
}

// 从缓存map中读取Receipt集合
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

// 从缓存map中读取StateDB实例
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

// 获取相应block的StateDB实例
func (c *Cache) MakeStateDB(block *types.Block) (*state.StateDB, error) {
	// 基于stateRoot从blockchain中创建StateDB实例
	if state, err := c.chain.StateAt(block.Root()); err == nil && state != nil {
		return state, nil
	}
	// 读取并拷贝缓存中StateDB实例
	sealHash := c.chain.Engine().SealHash(block.Header())
	log.Info("读取并拷贝缓存中StateDB实例","sealHash", sealHash, "blockHash", block.Hash(), "blockNum", block.NumberU64())
	if state := c.ReadStateDB(sealHash); state != nil {
		return state.Copy(), nil
		//return state, nil
	} else {
		return nil, errMakeStateDB
	}
}

// 获取相应block的StateDB实例
func (c *Cache) ClearCache(block *types.Block) {
	sealHash := c.chain.Engine().SealHash(block.Header())
	c.clearReceipts(sealHash)
	c.clearStateDB(sealHash)
}