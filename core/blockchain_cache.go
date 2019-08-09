package core

import (
	"errors"
	"fmt"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
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
	block := pbc.Engine().CurrentBlock()
	if block != nil {
		return block
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

func (pbc *BlockChainCache) GetBlockInMemory(hash common.Hash, number uint64) *types.Block {
	var block *types.Block
	if cbft, ok := pbc.Engine().(consensus.Bft); ok {
		log.Trace("find block in cbft", "RoutineID", common.CurrentGoRoutineID(), "hash", hash, "number", number)
		block = cbft.GetBlockWithoutLock(hash, number)
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
		return state, nil
	} else {
		return bcc.StateAt(header.Root)
	}
}

// Read the StateDB instance from the cache map
func (pbc *BlockChainCache) ReadStateDB(sealHash common.Hash) *state.StateDB {
	pbc.stateDBMu.RLock()
	defer pbc.stateDBMu.RUnlock()
	if obj, exist := pbc.stateDBCache[sealHash]; exist {
		log.Debug("Read the StateDB instance from the cache map", "sealHash", sealHash)
		return obj.stateDB.Copy()
	}
	return nil
}

// Write Receipt to the cache
func (pbc *BlockChainCache) WriteReceipts(sealHash common.Hash, receipts []*types.Receipt, blockNum uint64) {
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
func (bcc *BlockChainCache) WriteStateDB(sealHash common.Hash, stateDB *state.StateDB, blockNum uint64) {
	bcc.stateDBMu.Lock()
	defer bcc.stateDBMu.Unlock()
	log.Info("Write a StateDB instance to the cache", "sealHash", sealHash, "blockNum", blockNum)
	if _, exist := bcc.stateDBCache[sealHash]; !exist {
		bcc.stateDBCache[sealHash] = &stateDBCache{stateDB: stateDB, blockNum: blockNum}
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
		if obj.blockNum < blockNum {
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
		if obj.blockNum < blockNum {
			log.Debug("Clear StateDB", "sealHash", hash, "number", obj.blockNum)
			delete(bcc.stateDBCache, hash)
		}
	}
}

// Get the StateDB instance of the corresponding block
func (bcc *BlockChainCache) MakeStateDB(block *types.Block) (*state.StateDB, error) {
	log.Info("make stateDB", "hash", block.Hash(), "number", block.NumberU64(), "root", block.Root())
	// Create a StateDB instance from the blockchain based on stateRoot
	if state, err := bcc.StateAt(block.Root()); err == nil && state != nil {
		return state, nil
	}
	// Read and copy the stateDB instance in the cache
	sealHash := block.Header().SealHash()
	if state := bcc.ReadStateDB(sealHash); state != nil {
		//return state.Copy(), nil
		return state, nil
	} else {
		return nil, errMakeStateDB
	}
}

// Get the StateDB instance of the corresponding block
func (bcc *BlockChainCache) ClearCache(block *types.Block) {
	log.Debug("Clear Cache block", "sealHash", block.Header().SealHash(), "blockHash", block.Hash())
	sealHash := block.Header().SealHash()
	bcc.clearReceipts(sealHash)
	bcc.clearStateDB(sealHash)
}

func (bcc *BlockChainCache) StateDBString() string {
	status := fmt.Sprintf("[")
	for hash, obj := range bcc.stateDBCache {
		status += fmt.Sprintf("[%s, %d]", hash, obj.blockNum)
	}
	status += fmt.Sprintf("]")
	return status
}

func (bcc *BlockChainCache) Execute(block *types.Block, parent *types.Block) error {
	state, err := bcc.MakeStateDB(parent)
	if err != nil {
		return errors.New("execute block error")
	}

	//to execute
	receipts, err := bcc.ProcessDirectly(block, state, parent)
	log.Debug("Execute block", "number", block.Number(), "hash", block.Hash(),
		"parentNumber", parent.Number(), "parentHash", parent.Hash(), "err", err)
	if err == nil {
		//save the receipts and state to consensusCache
		sealHash := block.Header().SealHash()
		bcc.WriteReceipts(sealHash, receipts, block.NumberU64())
		bcc.WriteStateDB(sealHash, state, block.NumberU64())

	} else {
		return fmt.Errorf("execute block error, err:%s", err.Error())
	}
	return nil
}

func (bcc *BlockChainCache) WriteBlock(block *types.Block) error {
	sealHash := block.Header().SealHash()
	state := bcc.ReadStateDB(sealHash)
	receipts := bcc.ReadReceipts(sealHash)

	if state == nil {
		log.Error("Write Block error, state is nil", "number", block.NumberU64(), "hash", block.Hash())
		return fmt.Errorf("write Block error, state is nil, number:%d, hash:%s", block.NumberU64(), block.Hash().String())
	} else if len(block.Transactions()) > 0 && len(receipts) == 0 {
		log.Error("Write Block error, block has transactions but receipts is nil", "number", block.NumberU64(), "hash", block.Hash())
		return fmt.Errorf("write Block error, block has transactions but receipts is nil, number:%d, hash:%s", block.NumberU64(), block.Hash().String())
	}

	// Different block could share same sealhash, deep copy here to prevent write-write conflict.
	var _receipts = make([]*types.Receipt, len(receipts))
	for i, receipt := range receipts {
		_receipts[i] = new(types.Receipt)
		*_receipts[i] = *receipt
	}
	// Commit block and state to database.
	//block.SetExtraData(extraData)
	log.Debug("Write extra data", "txs", len(block.Transactions()), "extra", len(block.ExtraData()))
	_, err := bcc.WriteBlockWithState(block, _receipts, state)
	if err != nil {
		log.Error("Failed writing block to chain", "hash", block.Hash(), "number", block.NumberU64(), "err", err)
		return fmt.Errorf("Failed writing block to chain, number:%d, hash:%s, err:%s", block.NumberU64(), block.Hash().String(), err.Error())
	}

	log.Info("Successfully write new block", "hash", block.Hash(), "number", block.NumberU64(), "coinbase", block.Coinbase(), "time", block.Time())
	return nil
}
