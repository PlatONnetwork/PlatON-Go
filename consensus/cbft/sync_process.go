package cbft

import (
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/network"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
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
				for i, block := range blockList.Blocks {
					if err := cbft.verifyPrepareQC(blockList.QC[i]); err != nil {
						cbft.log.Error("Verify block prepare qc failed", "hash", block.Hash(), "number", block.NumberU64(), "error", err)
						return
					}

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
				utils.SetFalse(&cbft.fetching)
			}
		}

		expire := func() {
			utils.SetFalse(&cbft.fetching)
		}
		utils.SetTrue(&cbft.fetching)

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

	if err := cbft.verifyPrepareQC(msg.BlockQC); err != nil {
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

// OnGetPrepareVote is responsible for processing the business logic
// of the GetPrepareVote message. It will synchronously return a
// PrepareVotes message to the sender.
func (cbft *Cbft) OnGetPrepareVote(id string, msg *protocols.GetPrepareVote) error {
	// Get all the received PrepareVote of the block according to the index
	// position of the block in the view.
	prepareVoteMap := cbft.state.AllPrepareVoteByIndex(msg.BlockIndex)

	// Defining an array for receiving PrepareVote.
	votes := make([]*protocols.PrepareVote, 0, len(prepareVoteMap))
	if prepareVoteMap != nil {
		for k, v := range prepareVoteMap {
			if !msg.VoteBits.GetIndex(k) {
				votes = append(votes, v)
			}
		}
	} else {
		// todo: need to confirm.
		// Is it necessary to obtain the PrepareVotes from the blockchain
		// when it is not in the memory?
	}
	if len(votes) != 0 {
		cbft.network.Send(id, &protocols.PrepareVotes{BlockHash: msg.BlockHash, BlockNumber: msg.BlockNumber, Votes: votes})
		cbft.log.Debug("Send PrepareVotes", "peer", id, "hash", msg.BlockHash, "number", msg.BlockNumber)
	}
	return nil
}

// OnPrepareVotes handling response from GetPrepareVote response.
func (cbft *Cbft) OnPrepareVotes(id string, msg *protocols.PrepareVotes) error {
	for _, vote := range msg.Votes {
		if err := cbft.OnPrepareVote(id, vote); err != nil {
			cbft.log.Error("OnPrepareVotes failed", "peer", id, "err", err)
			return err
		}
	}
	return nil
}

func (cbft *Cbft) OnQCBlockList(id string, msg *protocols.QCBlockList) {
	// todo: Logic is incomplete.
}

// OnGetLatestStatus hands GetLatestStatus messages.
//
// main logic:
// 1.Compare the blockNumber of the sending node with the local node,
// and if the blockNumber of local node is larger then reply LatestStatus message,
// the message contains the status information of the local node.
func (cbft *Cbft) OnGetLatestStatus(id string, msg *protocols.GetLatestStatus) error {
	cbft.log.Debug("Received message on OnGetLatestStatus", "from", id, "logicType", msg.LogicType, "msgHash", msg.MsgHash().TerminalString())
	// Define a function that performs the send action.
	launcher := func(bType uint64, targetId string, message ctypes.Message) error {
		p, err := cbft.network.GetPeer(id)
		if err != nil {
			cbft.log.Error("GetPeer failed", "err", err, "peerId", id)
			return err
		}
		switch bType {
		case network.TypeForQCBn:
			p.SetQcBn(new(big.Int).SetUint64(msg.BlockNumber))
		case network.TypeForLockedBn:
			p.SetLockedBn(new(big.Int).SetUint64(msg.BlockNumber))
		case network.TypeForCommitBn:
			p.SetCommitdBn(new(big.Int).SetUint64(msg.BlockNumber))
		default:
		}
		// response to sender.
		cbft.network.Send(targetId, message)
		return nil
	}
	//
	if msg.LogicType == network.TypeForQCBn {
		localQCNum := cbft.state.HighestQCBlock().NumberU64()
		if localQCNum < msg.BlockNumber {
			cbft.log.Debug("Local qcBn is larger than the sender's qcBn", "remoteBn", msg.BlockNumber, "localBn", localQCNum)
			return launcher(msg.LogicType, id, &protocols.GetQCBlockList{
				BlockNumber: localQCNum,
				BlockHash:   cbft.state.HighestQCBlock().Hash(),
			})
		} else {
			cbft.log.Debug("Local qcBn is less than the sender's qcBn", "remoteBn", msg.BlockNumber, "localBn", localQCNum)
			cbft.network.Send(id, &protocols.LatestStatus{BlockNumber: localQCNum, LogicType: msg.LogicType})
		}
	}
	//
	if msg.LogicType == network.TypeForLockedBn {
		localLockedNum := cbft.state.HighestLockBlock().NumberU64()
		if localLockedNum < msg.BlockNumber {
			cbft.log.Debug("Local lockedBn is larger than the sender's lockedBn", "remoteBn", msg.BlockNumber, "localBn", localLockedNum)
			return launcher(msg.LogicType, id, &protocols.GetQCBlockList{
				BlockNumber: localLockedNum,
				BlockHash:   cbft.state.HighestLockBlock().Hash(),
			})
		} else {
			cbft.log.Debug("Local lockedBn is less than the sender's lockedBn", "remoteBn", msg.BlockNumber, "localBn", localLockedNum)
			cbft.network.Send(id, &protocols.LatestStatus{BlockNumber: localLockedNum, LogicType: msg.LogicType})
		}
	}
	//
	if msg.LogicType == network.TypeForCommitBn {
		localCommitNum := cbft.state.HighestCommitBlock().NumberU64()
		if localCommitNum < msg.BlockNumber {
			cbft.log.Debug("Local commitBn is larger than the sender's commitBn", "remoteBn", msg.BlockNumber, "localBn", localCommitNum)
			return launcher(msg.LogicType, id, &protocols.GetQCBlockList{
				BlockNumber: localCommitNum,
				BlockHash:   cbft.state.HighestCommitBlock().Hash(),
			})
		} else {
			cbft.log.Debug("Local commitBn is less than the sender's commitBn", "remoteBn", msg.BlockNumber, "localBn", localCommitNum)
			cbft.network.Send(id, &protocols.LatestStatus{BlockNumber: localCommitNum, LogicType: msg.LogicType})
		}
	}
	return nil
}

// OnLatestStatus is used to process LatestStatus messages that received from peer.
func (cbft *Cbft) OnLatestStatus(id string, msg *protocols.LatestStatus) error {
	cbft.log.Debug("Received message on OnLatestStatus", "from", id, "msgHash", msg.MsgHash().TerminalString())
	switch msg.LogicType {
	case network.TypeForQCBn:
		localQCBn := cbft.state.HighestQCBlock().NumberU64()
		if localQCBn < msg.BlockNumber {
			p, err := cbft.network.GetPeer(id)
			if err != nil {
				cbft.log.Error("GetPeer failed", "err", err)
				return err
			}
			p.SetQcBn(new(big.Int).SetUint64(msg.BlockNumber))
			cbft.log.Debug("LocalQCBn is lower than sender's", "localBn", localQCBn, "remoteBn", msg.BlockNumber)
			cbft.network.Send(id, &protocols.GetQCBlockList{
				BlockNumber: localQCBn,
				BlockHash:   cbft.state.HighestQCBlock().Hash(),
			})
		}

	case network.TypeForLockedBn:
		localLockedBn := cbft.state.HighestLockBlock().NumberU64()
		if localLockedBn < msg.BlockNumber {
			p, err := cbft.network.GetPeer(id)
			if err != nil {
				cbft.log.Error("GetPeer failed", "err", err)
				return err
			}
			p.SetLockedBn(new(big.Int).SetUint64(msg.BlockNumber))
			cbft.log.Debug("LocalLockedBn is lower than sender's", "localBn", localLockedBn, "remoteBn", msg.BlockNumber)
			cbft.network.Send(id, &protocols.GetQCBlockList{
				BlockNumber: localLockedBn,
				BlockHash:   cbft.state.HighestLockBlock().Hash(),
			})
		}

	case network.TypeForCommitBn:
		localCommitBn := cbft.state.HighestCommitBlock().NumberU64()
		if localCommitBn < msg.BlockNumber {
			p, err := cbft.network.GetPeer(id)
			if err != nil {
				cbft.log.Error("GetPeer failed", "err", err)
				return err
			}
			p.SetCommitdBn(new(big.Int).SetUint64(msg.BlockNumber))
			cbft.log.Debug("LocalCommitBn is lower than sender's", "localBn", localCommitBn, "remoteBn", msg.BlockNumber)
			cbft.network.Send(id, &protocols.GetQCBlockList{
				BlockNumber: localCommitBn,
				BlockHash:   cbft.state.HighestCommitBlock().Hash(),
			})
		}
	}
	return nil
}

// OnPrepareBlockHash responsible for handling PrepareBlockHash message.
//
// Note: After receiving the PrepareBlockHash message, it is determined whether the
// block information exists locally. If not, send a network request to get
// the block data.
func (cbft *Cbft) OnPrepareBlockHash(id string, msg *protocols.PrepareBlockHash) error {
	cbft.log.Debug("Received message on OnPrepareBlockHash", "from", id, "msgHash", msg.MsgHash())
	block := cbft.blockTree.FindBlockByHash(msg.BlockHash)
	if block == nil {
		cbft.network.Send(id, &protocols.GetPrepareBlock{
			Epoch:      msg.Epoch,
			ViewNumber: msg.ViewNumber,
			BlockIndex: msg.BlockIndex,
		})
	}
	return nil
}
