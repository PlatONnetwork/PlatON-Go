package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
)

// Get the block from the specified connection, get the block into the fetcher, and execute the block CBFT update state machine
func (cbft *Cbft) fetchBlock(id string, hash common.Hash, number uint64) {
	if cbft.state.HighestQCBlock().NumberU64() < number {
		//todo close receive consensus msg

		parent := cbft.state.HighestQCBlock()

		match := func(msg types.Message) bool {
			_, ok := msg.(*protocols.QCBlockList)
			return ok
		}

		executor := func(msg types.Message) {
			if blockList, ok := msg.(*protocols.QCBlockList); ok {
				// Execution block
				for _, block := range blockList.Blocks {
					if err := cbft.blockCacheWriter.Execute(block, parent); err != nil {
						cbft.log.Error("Execute block failed", "hash", block.Hash(), "number", block.NumberU64(), "error", err)
						return
					}
				}

				// Update the results to the CBFT state machine
				cbft.asyncCallCh <- func() {
					if err := cbft.OnInsertQCBlock(blockList.Blocks, blockList.QC); err != nil {
						cbft.log.Error("Insert block failed", "error", err)
					}
				}
			}
		}

		cbft.fetcher.AddTask("", match, executor, nil)
	}
}

func (cbft *Cbft) prepareVoteFetchRules(vote *protocols.PrepareVote) {
	if vote.BlockNumber > cbft.state.HighestQCBlock().NumberU64()+1 {
		//todo fetch qc
	}
}

func (cbft *Cbft) prepareBlockFetchRules(block *protocols.PrepareBlock) {
	if cbft.state.
}
