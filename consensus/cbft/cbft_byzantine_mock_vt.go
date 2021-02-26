package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
)

// 验证人伪造远大于当前index 的vote并广播，试图让其他节点FetchPrepare
// 预期结果：其他节点收到vote，校验签名通过，触发FetchPrepare调用
func (cbft *Cbft) MockVT01(proposalIndex uint32) {
	if !cbft.isProposer(cbft.state.Epoch(), cbft.state.ViewNumber(), proposalIndex) {
		block := cbft.state.ViewBlockByIndex(cbft.state.MaxViewBlockIndex())
		if block != nil {
			prepareVote := &protocols.PrepareVote{
				Epoch:          cbft.state.Epoch(),
				ViewNumber:     cbft.state.ViewNumber(),
				BlockHash:      common.BytesToHash(utils.Rand32Bytes(32)),
				BlockNumber:    block.NumberU64() + 5,
				BlockIndex:     cbft.config.Sys.Amount - 1,
				ValidatorIndex: proposalIndex,
			}
			cbft.signMsgByBls(prepareVote)
			cbft.network.Broadcast(prepareVote)
			cbft.log.Warn("[Mock-VT01]Broadcast future index prepareVote by validator", "nodeId", cbft.NodeID(), "prepareVote", prepareVote.String(), "currentIndex", cbft.state.NextViewBlockIndex()-1)
		}
	}
}

// 提议人伪造远大于当前index 的vote并广播，试图让其他节点FetchPrepare
// 预期结果：其他节点收到vote，校验签名通过，触发FetchPrepare调用
func (cbft *Cbft) MockVT02(proposalIndex uint32) {
	if cbft.isProposer(cbft.state.Epoch(), cbft.state.ViewNumber(), proposalIndex) {
		block := cbft.state.ViewBlockByIndex(cbft.state.MaxViewBlockIndex())
		if block != nil {
			prepareVote := &protocols.PrepareVote{
				Epoch:          cbft.state.Epoch(),
				ViewNumber:     cbft.state.ViewNumber(),
				BlockHash:      common.BytesToHash(utils.Rand32Bytes(32)),
				BlockNumber:    block.NumberU64() + 5,
				BlockIndex:     cbft.config.Sys.Amount - 1,
				ValidatorIndex: proposalIndex,
			}
			cbft.signMsgByBls(prepareVote)
			cbft.network.Broadcast(prepareVote)
			cbft.log.Warn("[Mock-VT01]Broadcast future index prepareVote by proposer", "nodeId", cbft.NodeID(), "prepareVote", prepareVote.String(), "currentIndex", cbft.state.NextViewBlockIndex()-1)
		}
	}
}

// 验证人伪造上一个区块签名，对下一个区块进行投票
// 预期结果：其他节点收到vote，校验parentQC不成功
func (cbft *Cbft) MockVT03(proposalIndex uint32) {
	if !cbft.isProposer(cbft.state.Epoch(), cbft.state.ViewNumber(), proposalIndex) {
		pending := cbft.state.PendingPrepareVote()
		if pending.Empty() {
			return
		}

		p := pending.Top()
		if p != nil {
			block := cbft.state.ViewBlockByIndex(p.BlockIndex)
			if block != nil {
				// 前一个区块未达到qc状态
				if b, qc := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1); b == nil || qc == nil {
					// 伪造上一个区块qc
					qcBlock := cbft.state.HighestQCBlock()
					_, qc := cbft.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
					p.ParentQC = qc
					cbft.network.Broadcast(p)
					cbft.log.Warn("[Mock-VT03]Broadcast mock prepareVote base fake parentQC by validator", "nodeId", cbft.NodeID(), "prepareVote", p.String(), "parentQC", p.ParentQC.String())
				}
			}
		}
	}
}

// 提议人伪造上一个区块签名，对下一个区块进行投票
// 预期结果：其他节点收到vote，校验parentQC不成功
func (cbft *Cbft) MockVT04(proposalIndex uint32) {
	if cbft.isProposer(cbft.state.Epoch(), cbft.state.ViewNumber(), proposalIndex) {
		pending := cbft.state.PendingPrepareVote()
		if pending.Empty() {
			return
		}

		p := pending.Top()
		if p != nil {
			block := cbft.state.ViewBlockByIndex(p.BlockIndex)
			if block != nil {
				// 前一个区块未达到qc状态
				if b, qc := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1); b == nil || qc == nil {
					// 伪造上一个区块qc
					qcBlock := cbft.state.HighestQCBlock()
					_, qc := cbft.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
					p.ParentQC = qc
					cbft.network.Broadcast(p)
					cbft.log.Warn("[Mock-VT04]Broadcast mock prepareVote base fake parentQC by proposer", "nodeId", cbft.NodeID(), "prepareVote", p.String(), "parentQC", p.ParentQC.String())
				}
			}
		}
	}
}

// 验证人发出不携带parentQC 的prepareVote
// 预期结果：其他节点收到vote，校验parentQC不成功
func (cbft *Cbft) MockVT05(proposalIndex uint32) {
	if !cbft.isProposer(cbft.state.Epoch(), cbft.state.ViewNumber(), proposalIndex) {
		pending := cbft.state.PendingPrepareVote()
		if pending.Empty() {
			return
		}

		p := pending.Top()
		if p != nil {
			p.ParentQC = nil
			cbft.network.Broadcast(p)
			cbft.log.Warn("[Mock-VT05]Broadcast mock prepareVote base empty parentQC by validator", "nodeId", cbft.NodeID(), "prepareVote", p.String())
		}
	}
}

// 验证人双签
// 预期结果：其他节点收到此消息，记录节点双签证据
func (cbft *Cbft) MockVT06(proposalIndex uint32) {
	if !cbft.isProposer(cbft.state.Epoch(), cbft.state.ViewNumber(), proposalIndex) {
		hadSend := cbft.state.HadSendPrepareVote()
		if hadSend.Empty() {
			return
		}

		p := hadSend.Top()
		if p != nil {
			prepareVote := &protocols.PrepareVote{
				Epoch:          p.Epoch,
				ViewNumber:     p.ViewNumber,
				BlockHash:      common.BytesToHash(utils.Rand32Bytes(32)),
				BlockNumber:    p.BlockNumber,
				BlockIndex:     p.BlockIndex,
				ValidatorIndex: p.ValidatorIndex,
			}
			prepareVote.ParentQC = p.ParentQC
			cbft.signMsgByBls(prepareVote)
			cbft.network.Broadcast(prepareVote)
			cbft.log.Warn("[Mock-VT06]Broadcast duplicate prepareVote", "nodeId", cbft.NodeID(), "prepareVote", prepareVote.String())
		}
	}
}

// 验证人伪造下一轮 view 的 vote 发出，试图让其他节点Fetch
// 预期结果：其他节点收到此消息，校验签名通过，但校验parentQC不通过，不会触发Fetch调用
func (cbft *Cbft) MockVT07(proposalIndex uint32) {
	qcBlock := cbft.state.HighestQCBlock()
	_, qc := cbft.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	if qc != nil {
		prepareVote := &protocols.PrepareVote{
			Epoch:          cbft.state.Epoch(),
			ViewNumber:     cbft.state.ViewNumber() + 1,
			BlockHash:      common.BytesToHash(utils.Rand32Bytes(32)),
			BlockNumber:    qc.BlockNumber + 3,
			BlockIndex:     cbft.config.Sys.Amount - 1,
			ValidatorIndex: proposalIndex,
		}
		prepareVote.ParentQC = qc
		cbft.signMsgByBls(prepareVote)
		cbft.network.Broadcast(prepareVote)
		cbft.log.Warn("[Mock-VT07]Broadcast next view prepareVote", "nodeId", cbft.NodeID(), "prepareVote", prepareVote.String())
	}
}

// 验证人伪造下一轮 epoch 的 vote 发出，试图让其他节点Fetch
// 预期结果：其他节点收到此消息，校验签名通过，但校验parentQC不通过，不会触发Fetch调用
func (cbft *Cbft) MockVT08(proposalIndex uint32) {
	qcBlock := cbft.state.HighestQCBlock()
	_, qc := cbft.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())

	if qc != nil {
		prepareVote := &protocols.PrepareVote{
			Epoch:          cbft.state.Epoch() + 1,
			ViewNumber:     state.DefaultViewNumber,
			BlockHash:      common.BytesToHash(utils.Rand32Bytes(32)),
			BlockNumber:    qc.BlockNumber + 1,
			BlockIndex:     cbft.config.Sys.Amount - 1,
			ValidatorIndex: proposalIndex,
		}
		prepareVote.ParentQC = qc
		cbft.signMsgByBls(prepareVote)
		cbft.network.Broadcast(prepareVote)
		cbft.log.Warn("[Mock-VT08]Broadcast next epoch prepareVote", "nodeId", cbft.NodeID(), "prepareVote", prepareVote.String())
	}
}
