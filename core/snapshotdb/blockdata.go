package snapshotdb

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/syndtr/goleveldb/leveldb/memdb"
	"math/big"
	"sync"
)

type blockData struct {
	mu         sync.Mutex
	BlockHash  *common.Hash
	ParentHash common.Hash
	Number     *big.Int
	data       *memdb.DB
	readOnly   bool
	kvHash     common.Hash
}
