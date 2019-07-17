package cbft

import "github.com/PlatONnetwork/PlatON-Go/common"

func (cbft *Cbft) fetchBlock(hash common.Hash, number uint64) {
	if cbft.state.HighestExecutedBlock().NumberU64() < number {
		//todo close receive consensus msg
		//todo add fetch task
		//todo run task & waiting block
		//todo return result
	}
}
