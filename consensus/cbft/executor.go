package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type blockExecutor interface {
	//Execution block, you need to pass in the parent block to find the parent block state
	execute(block *types.Block, parent *types.Block) error
}

//Block execution results, including block hash, block number, error message
type blockExecuteStatus struct {
	hash   common.Hash
	number uint64
	err    error
}

type asyncBlockExecutor interface {
	blockExecutor
	//Asynchronous acquisition block execution results
	executeStatus() chan blockExecuteStatus
}
