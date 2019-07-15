package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type blockTree struct {
	root   *blockTree
	blocks map[uint64]map[common.Hash]*blockTree
}

type blockExt struct {
	viewNumber   uint64
	block        *types.Block
	inTree       bool
	executing    bool
	isExecuted   bool
	isSigned     bool
	isConfirmed  bool
	rcvTime      int64
	prepareVotes *prepareVoteSet //all prepareVotes for block
	prepareBlock *prepareBlock
	parent       *blockExt
	children     map[common.Hash]*blockExt
}
