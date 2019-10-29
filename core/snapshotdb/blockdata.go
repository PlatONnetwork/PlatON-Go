package snapshotdb

import (
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb/memdb"
)

type blockData struct {
	BlockHash  common.Hash
	ParentHash common.Hash
	Number     *big.Int
	data       *memdb.DB
	readOnly   bool
	kvHash     common.Hash
}

type unCommitBlocks struct {
	blocks map[common.Hash]*blockData
	sync.RWMutex
}

func (u *unCommitBlocks) Get(key common.Hash) *blockData {
	u.RLock()
	block, ok := u.blocks[key]
	u.RUnlock()
	if !ok {
		return nil
	}
	return block
}

func (u *unCommitBlocks) Set(key common.Hash, block *blockData) {
	u.Lock()
	u.blocks[key] = block
	u.Unlock()
}
