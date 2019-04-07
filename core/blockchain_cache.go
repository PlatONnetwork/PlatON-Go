package core

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"math/big"
	"sync"
)

var (
	errMakeStateDB = errors.New("make StateDB error")
)

type BlockChainCache struct {
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

func (pbc *BlockChainCache) CurrentBlock() *types.Block {
	if cbft, ok := pbc.Engine().(consensus.Bft); ok {
		if block := cbft.HighestLogicalBlock(); block != nil {
			return block
		}
	}
	return pbc.currentBlock.Load().(*types.Block)
}

func (pbc *BlockChainCache) GetBlock(hash common.Hash, number uint64) *types.Block {
	var block *types.Block
	if cbft, ok := pbc.Engine().(consensus.Bft); ok {
		log.Trace("find block in cbft", "RoutineID", common.CurrentGoRoutineID(), "hash", hash, "number", number)
		block = cbft.GetBlock(hash, number)
	}
	if block == nil {
		log.Trace("cannot find block in cbft, try to find it in chain", "RoutineID", common.CurrentGoRoutineID(), "hash", hash, "number", number)
		block = pbc.getBlock(hash, number)
		if block == nil {
			log.Trace("cannot find block in chain", "RoutineID", common.CurrentGoRoutineID(), "hash", hash, "number", number)
		}
	}
	return block
}

func NewBlockChainCache(blockChain *BlockChain) *BlockChainCache {
	pbc := &BlockChainCache{}
	pbc.BlockChain = blockChain
	pbc.stateDBCache = make(map[common.Hash]*stateDBCache)
	pbc.receiptsCache = make(map[common.Hash]*receiptsCache)

	return pbc
}

// Read the Receipt collection from the cache map.
func (bcc *BlockChainCache) ReadReceipts(sealHash common.Hash) []*types.Receipt {
	bcc.receiptsMu.RLock()
	defer bcc.receiptsMu.RUnlock()
	if obj, exist := bcc.receiptsCache[sealHash]; exist {
		return obj.receipts
	}
	return nil
}

// GetState returns a new mutable state based on a particular point in time.
func (bcc *BlockChainCache) GetState(header *types.Header) (*state.StateDB, error) {
	state := bcc.ReadStateDB(header.SealHash())
	if state != nil {
		log.Info("BlockChainCache GetState", "addr", fmt.Sprintf("%p", state), "root", header.Root)
		return state, nil
	} else {
		log.Info("BlockChainCache GetState", "root", header.Root)
		return bcc.StateAt(header.Root, header.Number, header.Hash())
	}
}

// Read the StateDB instance from the cache map
func (pbc *BlockChainCache) ReadStateDB(sealHash common.Hash) *state.StateDB {
	pbc.stateDBMu.RLock()
	defer pbc.stateDBMu.RUnlock()
	log.Info("Read the StateDB instance from the cache map", "sealHash", sealHash)
	if obj, exist := pbc.stateDBCache[sealHash]; exist {
		return obj.stateDB.Copy()
	}
	return nil
}

// Write Receipt to the cache
func (pbc *BlockChainCache) WriteReceipts(sealHash common.Hash, receipts []*types.Receipt, blockNum uint64) {
	pbc.receiptsMu.Lock()
	defer pbc.receiptsMu.Unlock()
	obj, exist := pbc.receiptsCache[sealHash]
	if exist {
		if obj.blockNum == blockNum && len(obj.receipts) == len(receipts) {
			log.Info("the receipts already in cache")
		} else {
			log.Warn("there maybe an error!", "blockNum", blockNum, "obj.blockNum", obj.blockNum, "len(obj.receipts)", len(obj.receipts), "len(receipts)", len(receipts))
			obj.receipts = append(obj.receipts, receipts...)
		}
	} else if !exist {
		pbc.receiptsCache[sealHash] = &receiptsCache{receipts: receipts, blockNum: blockNum}
	}
}

// Write a StateDB instance to the cache
func (bcc *BlockChainCache) WriteStateDB(sealHash common.Hash, stateDB *state.StateDB, blockNum uint64) {
	bcc.stateDBMu.Lock()
	defer bcc.stateDBMu.Unlock()
	log.Info("Write a StateDB instance to the cache", "sealHash", sealHash, "blockNum", blockNum)
	if _, exist := bcc.stateDBCache[sealHash]; !exist {
		stateDBCpy := stateDB.Copy()
		bcc.stateDBCache[sealHash] = &stateDBCache{stateDB: stateDBCpy, blockNum: blockNum}
	}
}

// Read the Receipt collection from the cache map
func (bcc *BlockChainCache) clearReceipts(sealHash common.Hash) {
	bcc.receiptsMu.Lock()
	defer bcc.receiptsMu.Unlock()

	var blockNum uint64
	if obj, exist := bcc.receiptsCache[sealHash]; exist {
		blockNum = obj.blockNum
		//delete(pbc.receiptsCache, sealHash)
	}
	for hash, obj := range bcc.receiptsCache {
		if obj.blockNum <= blockNum {
			delete(bcc.receiptsCache, hash)
		}
	}
}

// Read the StateDB instance from the cache map
func (bcc *BlockChainCache) clearStateDB(sealHash common.Hash) {
	bcc.stateDBMu.Lock()
	defer bcc.stateDBMu.Unlock()
	var blockNum uint64
	if obj, exist := bcc.stateDBCache[sealHash]; exist {
		blockNum = obj.blockNum
		//delete(pbc.stateDBCache, sealHash)
	}
	for hash, obj := range bcc.stateDBCache {
		if obj.blockNum <= blockNum {
			root := obj.stateDB.IntermediateRoot(bcc.chainConfig.IsEIP158(big.NewInt(int64(obj.blockNum))))
			log.Info("Delete StateDB Cache", "blockNumber", obj.blockNum, "sealHash", sealHash.String(), "stateDB root", root.String())
			delete(bcc.stateDBCache, hash)
		}
	}
}

// Get the StateDB instance of the corresponding block
func (bcc *BlockChainCache) MakeStateDB(block *types.Block) (*state.StateDB, error) {
	// Create a StateDB instance from the blockchain based on stateRoot
	log.Info("------make StateDB------", "GoRoutineID", common.CurrentGoRoutineID(), "number", block.NumberU64(), "hash", block.Hash(), "stateRoot", block.Root())
	curBlock := bcc.BlockChain.CurrentBlock()
	if curBlock != nil {
		log.Info("------current block------", "GoRoutineID", common.CurrentGoRoutineID(), "number", curBlock.NumberU64(), "hash", curBlock.Hash(), "stateRoot", curBlock.Root())
	}
	log.Info("---------recheck Block", "number", block.NumberU64(), "hash", block.Hash(), "root", block.Root())

	// Read and copy the stateDB instance in the cache
	sealHash := bcc.Engine().SealHash(block.Header())
	log.Info("Read and copy the stateDB instance in the cache", "sealHash", sealHash, "blockHash", block.Hash(), "blockNum", block.NumberU64(), "stateRoot", block.Root())
	if state := bcc.ReadStateDB(sealHash); state != nil {
		log.Debug("MakeStateDB", "addr", fmt.Sprintf("%p", state))
		//return state.Copy(), nil
		return state, nil
	}
	if state, err := bcc.StateAt(block.Root(), block.Number(), block.Hash()); err == nil && state != nil {
		log.Info("---------recheck check Block", "addr", fmt.Sprintf("%p", state), "number", block.NumberU64(), "hash", block.Hash(), "root", block.Root())
		return state, nil
	} else {
		return nil, errMakeStateDB
	}
}

// Get the StateDB instance of the corresponding block
func (bcc *BlockChainCache) ClearCache(block *types.Block) {
	sealHash := bcc.Engine().SealHash(block.Header())
	bcc.clearReceipts(sealHash)
	bcc.clearStateDB(sealHash)
}
