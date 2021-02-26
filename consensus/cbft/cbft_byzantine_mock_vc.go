package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
)

// 验证节点基于block（n），prepareQC（n-1）发起viewChange
// 预期结果：其他节点收到viewChange，校验prepareQC不成功
func (cbft *Cbft) MockVC01() {
	node, err := cbft.isCurrentValidator()
	if err != nil {
		return
	}

	qcHash, qcNumber := cbft.state.HighestQCBlock().Hash(), cbft.state.HighestQCBlock().NumberU64()
	lockHash, lockNumber := cbft.state.HighestLockBlock().Hash(), cbft.state.HighestLockBlock().NumberU64()
	_, lockQC := cbft.blockTree.FindBlockAndQC(lockHash, lockNumber)

	if lockQC != nil {
		viewChange := &protocols.ViewChange{
			Epoch:          cbft.state.Epoch(),
			ViewNumber:     cbft.state.ViewNumber(),
			BlockHash:      qcHash,
			BlockNumber:    qcNumber,
			ValidatorIndex: node.Index,
			PrepareQC:      lockQC,
		}
		cbft.signMsgByBls(viewChange)

		cbft.network.Broadcast(viewChange)
		cbft.log.Warn("[Mock-VC01]Broadcast mock viewChange base highest qcBlock,fake prepareQC", "nodeId", cbft.NodeID(), "viewChange", viewChange.String())
	}
}

// 非验证人伪造验证节点发起viewChange
// 预期结果：其他节点收到viewChange，校验消息签名不成功
func (cbft *Cbft) MockVC02() {
	_, err := cbft.isCurrentValidator()
	if err != nil { // current node is not Validator
		qcHash, qcNumber := cbft.state.HighestQCBlock().Hash(), cbft.state.HighestQCBlock().NumberU64()
		_, qc := cbft.blockTree.FindBlockAndQC(qcHash, qcNumber)

		if qc != nil {
			viewChange := &protocols.ViewChange{
				Epoch:          cbft.state.Epoch(),
				ViewNumber:     cbft.state.ViewNumber(),
				BlockHash:      qcHash,
				BlockNumber:    qcNumber,
				ValidatorIndex: 0, // fake nodeIndex
				PrepareQC:      qc,
			}
			cbft.signMsgByBls(viewChange)

			cbft.network.Broadcast(viewChange)
			cbft.log.Warn("[Mock-VC02]Broadcast mock viewChange by non validator", "nodeId", cbft.NodeID(), "viewChange", viewChange.String())
		}
	}
}

// 验证人双出
// 预期结果：其他节点收到此消息，记录节点双出证据
func (cbft *Cbft) MockVC03() {
	node, err := cbft.isCurrentValidator()
	if err != nil {
		return
	}

	hadSend := cbft.state.ViewChangeByIndex(node.Index)
	if hadSend == nil { // had not send viewchange
		return
	}

	viewChange := &protocols.ViewChange{
		Epoch:          hadSend.Epoch,
		ViewNumber:     hadSend.ViewNumber,
		BlockHash:      common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:    hadSend.BlockNumber,
		ValidatorIndex: hadSend.ValidatorIndex,
		PrepareQC:      hadSend.PrepareQC,
	}
	cbft.signMsgByBls(viewChange)

	cbft.network.Broadcast(viewChange)
	cbft.log.Warn("[Mock-VC03]Broadcast duplicate viewChange", "nodeId", cbft.NodeID(), "viewChange", viewChange.String())
}

// 验证人伪造下一轮view 的viewChange发出，试图让其他节点Fetch
// 预期结果：其他节点收到此消息，校验签名通过，校验prepareQC通过，触发Fetch调用
func (cbft *Cbft) MockVC04() {
	node, err := cbft.isCurrentValidator()
	if err != nil {
		return
	}

	qcHash, qcNumber := cbft.state.HighestQCBlock().Hash(), cbft.state.HighestQCBlock().NumberU64()
	_, qc := cbft.blockTree.FindBlockAndQC(qcHash, qcNumber)

	if qc != nil {
		viewChange := &protocols.ViewChange{
			Epoch:          cbft.state.Epoch(),
			ViewNumber:     cbft.state.ViewNumber() + 1,
			BlockHash:      qcHash,
			BlockNumber:    qcNumber,
			ValidatorIndex: node.Index,
			PrepareQC:      qc,
		}
		cbft.signMsgByBls(viewChange)

		cbft.network.Broadcast(viewChange)
		cbft.log.Warn("[Mock-VC04]Broadcast next view viewChange", "nodeId", cbft.NodeID(), "viewChange", viewChange.String())
	}
}

// 验证人伪造下一轮epoch 的viewChange发出，试图让其他节点Fetch
// 预期结果：其他节点收到此消息，校验签名通过，校验prepareQC通过，触发Fetch调用
func (cbft *Cbft) MockVC05() {
	node, err := cbft.isCurrentValidator()
	if err != nil {
		return
	}

	qcHash, qcNumber := cbft.state.HighestQCBlock().Hash(), cbft.state.HighestQCBlock().NumberU64()
	_, qc := cbft.blockTree.FindBlockAndQC(qcHash, qcNumber)

	if qc != nil {
		viewChange := &protocols.ViewChange{
			Epoch:          cbft.state.Epoch() + 1,
			ViewNumber:     state.DefaultViewNumber,
			BlockHash:      qcHash,
			BlockNumber:    qcNumber,
			ValidatorIndex: node.Index,
			PrepareQC:      qc,
		}
		cbft.signMsgByBls(viewChange)

		cbft.network.Broadcast(viewChange)
		cbft.log.Warn("[Mock-VC05]Broadcast next epoch viewChange", "nodeId", cbft.NodeID(), "viewChange", viewChange.String())
	}
}
