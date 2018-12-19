package core

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"sync"
)

var (
	errMakeStateDB = errors.New("make StateDB error")
)

type PlatonBlockChain struct {
	*BlockChain
	stateDBCache  map[common.Hash]*stateDBCache  // key is header SealHash
	receiptsCache map[common.Hash]*receiptsCache // key is header SealHash
	stateDBMu     sync.RWMutex
	receiptsMu    sync.RWMutex
}

type stateDBCache struct {
	stateDB  *state.StateDB
	blockNum uint64
}

type receiptsCache struct {
	receipts []*types.Receipt
	blockNum uint64
}

func (pbc *PlatonBlockChain) CurrentBlock() *types.Block {
	if cbft, ok := pbc.Engine().(consensus.Bft); ok {
		return cbft.HighestLogicalBlock()
	} else {
		return pbc.currentBlock.Load().(*types.Block)
	}
}

func (pbc *PlatonBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	var block *types.Block
	if cbft, ok := pbc.Engine().(consensus.Bft); ok {
		block = cbft.GetBlock(hash, number)
	}
	if block == nil {
		return pbc.getBlock(hash, number)
	} else {
		return block
	}
}

func NewPlatonBlockChain(blockChain *BlockChain) *PlatonBlockChain {
	pbc := &PlatonBlockChain{}
	pbc.BlockChain = blockChain
	pbc.stateDBCache = make(map[common.Hash]*stateDBCache)
	pbc.receiptsCache = make(map[common.Hash]*receiptsCache)

	return pbc
}

// Read the Receipt collection from the cache map.
func (pbc *PlatonBlockChain) ReadReceipts(sealHash common.Hash) []*types.Receipt {
	pbc.receiptsMu.RLock()
	defer pbc.receiptsMu.RUnlock()
	if obj, exist := pbc.receiptsCache[sealHash]; exist {
		return obj.receipts
	}
	return nil
}

// Read the StateDB instance from the cache map
func (pbc *PlatonBlockChain) ReadStateDB(sealHash common.Hash) *state.StateDB {
	pbc.stateDBMu.RLock()
	defer pbc.stateDBMu.RUnlock()
	log.Info("Read the StateDB instance from the cache map", "sealHash", sealHash)
	if obj, exist := pbc.stateDBCache[sealHash]; exist {
		return obj.stateDB.Copy()
	}
	return nil
}

// Write Receipt to the cache
func (pbc *PlatonBlockChain) WriteReceipts(sealHash common.Hash, receipts []*types.Receipt, blockNum uint64) {
	pbc.receiptsMu.Lock()
	defer pbc.receiptsMu.Unlock()
	obj, exist := pbc.receiptsCache[sealHash]
	if exist && obj.blockNum == blockNum {
		obj.receipts = append(obj.receipts, receipts...)
	} else if !exist {
		pbc.receiptsCache[sealHash] = &receiptsCache{receipts: receipts, blockNum: blockNum}
	}
}

// Write a StateDB instance to the cache
func (pbc *PlatonBlockChain) WriteStateDB(sealHash common.Hash, stateDB *state.StateDB, blockNum uint64) {
	pbc.stateDBMu.Lock()
	defer pbc.stateDBMu.Unlock()
	log.Info("Write a StateDB instance to the cache", "sealHash", sealHash, "blockNum", blockNum)
	if _, exist := pbc.stateDBCache[sealHash]; !exist {
		pbc.stateDBCache[sealHash] = &stateDBCache{stateDB: stateDB, blockNum: blockNum}
	}
}

// Read the Receipt collection from the cache map
func (pbc *PlatonBlockChain) clearReceipts(sealHash common.Hash) {
	pbc.receiptsMu.Lock()
	defer pbc.receiptsMu.Unlock()

	var blockNum uint64
	if obj, exist := pbc.receiptsCache[sealHash]; exist {
		blockNum = obj.blockNum
		//delete(pbc.receiptsCache, sealHash)
	}
	for hash, obj := range pbc.receiptsCache {
		if obj.blockNum <= blockNum {
			delete(pbc.receiptsCache, hash)
		}
	}
}

// Read the StateDB instance from the cache map
func (pbc *PlatonBlockChain) clearStateDB(sealHash common.Hash) {
	pbc.stateDBMu.Lock()
	defer pbc.stateDBMu.Unlock()

	var blockNum uint64
	if obj, exist := pbc.stateDBCache[sealHash]; exist {
		blockNum = obj.blockNum
		//delete(pbc.stateDBCache, sealHash)
	}
	for hash, obj := range pbc.stateDBCache {
		if obj.blockNum <= blockNum {
			delete(pbc.stateDBCache, hash)
		}
	}
}

// Get the StateDB instance of the corresponding block
func (pbc *PlatonBlockChain) MakeStateDB(block *types.Block) (*state.StateDB, error) {
	// Create a StateDB instance from the blockchain based on stateRoot
	if state, err := pbc.StateAt(block.Root()); err == nil && state != nil {
		return state, nil
	}
	// Read and copy the stateDB instance in the cache
	sealHash := pbc.Engine().SealHash(block.Header())
	log.Info("Read and copy the stateDB instance in the cache", "sealHash", sealHash, "blockHash", block.Hash(), "blockNum", block.NumberU64(), "stateRoot", block.Root())
	if state := pbc.ReadStateDB(sealHash); state != nil {
		//return state.Copy(), nil
		return state, nil
	} else {
		return nil, errMakeStateDB
	}
}

// Get the StateDB instance of the corresponding block
func (pbc *PlatonBlockChain) ClearCache(block *types.Block) {
	sealHash := pbc.Engine().SealHash(block.Header())
	pbc.clearReceipts(sealHash)
	pbc.clearStateDB(sealHash)
}
