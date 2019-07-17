package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
)

func (cbft *Cbft) fetchBlock(hash common.Hash, number uint64) {
	if cbft.state.HighestQCBlock().NumberU64() < number {
		//parent := cbft.state.HighestQCBlock()
		//todo close receive consensus msg
		match := func(msg types.Message) bool {
			_, ok := msg.(*protocols.QCBlockList)
			return ok
		}
		executor := func(msg types.Message) {
			//var exe executor.BlockExecutor
			//if blockList, ok := msg.(*protocols.QCBlockList); ok {
			//	for _, block := range blockList.Blocks {
			//	}
			//}
		}
		cbft.fetcher.AddTask("", match, executor)
		//todo add fetch task
		//todo run task & waiting block
		//todo return result
	}
}
