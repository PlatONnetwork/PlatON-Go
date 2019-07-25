package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

// Get the block from the specified connection, get the block into the fetcher, and execute the block CBFT update state machine
func (cbft *Cbft) fetchBlock(id string, hash common.Hash, number uint64) {
	if cbft.state.HighestQCBlock().NumberU64() < number {

		parent := cbft.state.HighestQCBlock()

		match := func(msg ctypes.Message) bool {
			_, ok := msg.(*protocols.QCBlockList)
			return ok
		}

		executor := func(msg ctypes.Message) {
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
				cbft.fetching = false
			}
		}

		expire := func() {
			cbft.fetching = false
		}
		cbft.fetching = true

		cbft.fetcher.AddTask(id, match, executor, expire)
		cbft.network.Send(id, &protocols.GetQCBlockList{BlockHash: cbft.state.HighestQCBlock().Hash(), BlockNumber: cbft.state.HighestQCBlock().NumberU64()})
	}
}

// Obtain blocks that are not in the local according to the proposed block
func (cbft *Cbft) prepareBlockFetchRules(id string, pb *protocols.PrepareBlock) {
	if pb.Block.NumberU64() > cbft.state.HighestQCBlock().NumberU64() {
		for i := uint32(0); i < pb.BlockIndex; i++ {
			b, _ := cbft.state.ViewBlockAndQC(i)
			if b == nil {
				msg := &protocols.GetPrepareBlock{Epoch: cbft.state.Epoch(), ViewNumber: cbft.state.ViewNumber(), BlockIndex: i}
				cbft.network.Send(id, msg)
				cbft.log.Debug("Send GetPrepareBlock", "peer", id, "msg", msg.String())
			}
		}
	}
}

// Get votes and blocks that are not available locally based on the height of the vote
func (cbft *Cbft) prepareVoteFetchRules(id string, vote *protocols.PrepareVote) {
	// Greater than QC+1 means the vote is behind
	if vote.BlockNumber > cbft.state.HighestQCBlock().NumberU64()+1 {
		for i := uint32(0); i < vote.BlockIndex; i++ {
			b, q := cbft.state.ViewBlockAndQC(i)
			if b == nil {
				msg := &protocols.GetPrepareBlock{Epoch: cbft.state.Epoch(), ViewNumber: cbft.state.ViewNumber(), BlockIndex: i}
				cbft.network.Send(id, msg)
				cbft.log.Debug("Send GetPrepareBlock", "peer", id, "msg", msg.String())
			} else if q != nil {
				msg := &protocols.GetBlockQuorumCert{BlockHash: b.Hash(), BlockNumber: b.NumberU64()}
				cbft.network.Send(id, msg)
				cbft.log.Debug("Send GetBlockQuorumCert", "peer", id, "msg", msg.String())
			}
		}
	}
}

func (cbft *Cbft) OnGetPrepareBlock(id string, msg *protocols.GetPrepareBlock) {
	if msg.Epoch == cbft.state.Epoch() && msg.ViewNumber == cbft.state.ViewNumber() {
		prepareBlock := cbft.state.PrepareBlockByIndex(msg.BlockIndex)
		if prepareBlock != nil {
			cbft.log.Debug("Send PrepareBlock", "prepareBlock", prepareBlock.String())
			cbft.network.Send(id, prepareBlock)
		}
	}
}

func (cbft *Cbft) OnGetBlockQuorumCert(id string, msg *protocols.GetBlockQuorumCert) {
	_, qc := cbft.blockTree.FindBlockAndQC(msg.BlockHash, msg.BlockNumber)
	if qc != nil {
		cbft.network.Send(id, &protocols.BlockQuorumCert{BlockQC: qc})
	}
}

func (cbft *Cbft) OnBlockQuorumCert(id string, msg *protocols.BlockQuorumCert) {
	if msg.BlockQC.Epoch != cbft.state.Epoch() || msg.BlockQC.ViewNumber != cbft.state.ViewNumber() {
		cbft.log.Debug("Receive BlockQuorumCert response failed", "local.epoch", cbft.state.Epoch(), "local.viewNumber", cbft.state.ViewNumber(), "msg", msg.String())
		return
	}

	block := cbft.state.ViewBlockByIndex(msg.BlockQC.BlockIndex)
	if block != nil {
		cbft.insertQCBlock(block, msg.BlockQC)
		cbft.log.Debug("Receive BlockQuorumCert success", "msg", msg.String())
	}
}

func (cbft *Cbft) OnGetQCBlockList(id string, msg *protocols.GetQCBlockList) {
	highestQC := cbft.state.HighestQCBlock()

	if highestQC.NumberU64() > msg.BlockNumber+3 ||
		(highestQC.Hash() == msg.BlockHash && highestQC.NumberU64() == msg.BlockNumber) {
		cbft.log.Debug(fmt.Sprintf("Receive GetQCBlockList failed, local.highestQC:%s,%d, msg:%s", highestQC.Hash().String(), highestQC.NumberU64(), msg.String()))
		return
	}

	lock := cbft.state.HighestLockBlock()
	commit := cbft.state.HighestCommitBlock()

	qcs := make([]*ctypes.QuorumCert, 0)
	blocks := make([]*types.Block, 0)

	if commit.ParentHash() == msg.BlockHash {
		block, qc := cbft.blockTree.FindBlockAndQC(commit.Hash(), commit.NumberU64())
		qcs = append(qcs, qc)
		blocks = append(blocks, block)
	}

	if lock.ParentHash() == msg.BlockHash || commit.ParentHash() == msg.BlockHash {
		block, qc := cbft.blockTree.FindBlockAndQC(lock.Hash(), lock.NumberU64())
		qcs = append(qcs, qc)
		blocks = append(blocks, block)
	}
	if highestQC.ParentHash() == msg.BlockHash || lock.ParentHash() == msg.BlockHash || commit.ParentHash() == msg.BlockHash {
		block, qc := cbft.blockTree.FindBlockAndQC(highestQC.Hash(), highestQC.NumberU64())
		qcs = append(qcs, qc)
		blocks = append(blocks, block)
	}

	if len(qcs) != 0 {
		cbft.network.Send(id, &protocols.QCBlockList{QC: qcs, Blocks: blocks})
		cbft.log.Debug("Send QCBlockList", "len", len(qcs))
	}

}
