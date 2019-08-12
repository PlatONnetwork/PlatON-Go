package cbft

import (
	"container/list"
	"fmt"
	"math/big"
	"sort"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/network"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// Get the block from the specified connection, get the block into the fetcher, and execute the block CBFT update state machine
func (cbft *Cbft) fetchBlock(id string, hash common.Hash, number uint64) {
	if cbft.fetcher.Len() != 0 {
		cbft.log.Trace("Had fetching block")
		return
	}

	baseBlockHash, baseBlockNumber := common.Hash{}, uint64(0)
	var parentBlock *types.Block

	if cbft.state.HighestQCBlock().NumberU64() < number {
		parentBlock = cbft.state.HighestQCBlock()
		baseBlockHash, baseBlockNumber = parentBlock.Hash(), parentBlock.NumberU64()
	} else if cbft.state.HighestQCBlock().NumberU64() == number {
		parentBlock = cbft.state.HighestLockBlock()
		baseBlockHash, baseBlockNumber = parentBlock.Hash(), parentBlock.NumberU64()
	} else {
		cbft.log.Trace("No suitable block need to request")
		return
	}

	match := func(msg ctypes.Message) bool {
		_, ok := msg.(*protocols.QCBlockList)
		return ok
	}

	executor := func(msg ctypes.Message) {
		defer func() {
			cbft.log.Debug("Close fetching")
			utils.SetFalse(&cbft.fetching)
		}()
		if blockList, ok := msg.(*protocols.QCBlockList); ok {
			// Execution block
			for i, block := range blockList.Blocks {
				if err := cbft.verifyPrepareQC(blockList.QC[i]); err != nil {
					cbft.log.Error("Verify block prepare qc failed", "hash", block.Hash(), "number", block.NumberU64(), "error", err)
					return
				}
				start := time.Now()
				if err := cbft.blockCacheWriter.Execute(block, parentBlock); err != nil {
					cbft.log.Error("Execute block failed", "hash", block.Hash(), "number", block.NumberU64(), "error", err)
					return
				}
				blockExecutedTimer.UpdateSince(start)
				parentBlock = block
			}

			// Update the results to the CBFT state machine
			cbft.asyncCallCh <- func() {
				if err := cbft.OnInsertQCBlock(blockList.Blocks, blockList.QC); err != nil {
					cbft.log.Error("Insert block failed", "error", err)
				}
			}
		}
	}

	expire := func() {
		cbft.log.Debug("Fetch timeout, close fetching")
		utils.SetFalse(&cbft.fetching)
	}

	cbft.log.Debug("Start fetching")

	utils.SetTrue(&cbft.fetching)
	cbft.fetcher.AddTask(id, match, executor, expire)
	cbft.network.Send(id, &protocols.GetQCBlockList{BlockHash: baseBlockHash, BlockNumber: baseBlockNumber})
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

	cbft.insertPrepareQC(msg.BlockQC)
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
	cbft.log.Debug("Received message on OnGetPrepareVote", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
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
	cbft.log.Debug("Received message on OnPrepareVotes", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	for _, vote := range msg.Votes {
		if err := cbft.OnPrepareVote(id, vote); err != nil {
			cbft.log.Error("OnPrepareVotes failed", "peer", id, "err", err)
			return err
		}
	}
	return nil
}

// OnGetLatestStatus hands GetLatestStatus messages.
//
// main logic:
// 1.Compare the blockNumber of the sending node with the local node,
// and if the blockNumber of local node is larger then reply LatestStatus message,
// the message contains the status information of the local node.
func (cbft *Cbft) OnGetLatestStatus(id string, msg *protocols.GetLatestStatus) error {
	cbft.log.Debug("Received message on OnGetLatestStatus", "from", id, "logicType", msg.LogicType, "msgHash", msg.MsgHash(), "message", msg.String())
	// Define a function that performs the send action.
	launcher := func(bType uint64, targetId string, blockNumber uint64, blockHash common.Hash) error {
		p, err := cbft.network.GetPeer(targetId)
		if err != nil {
			cbft.log.Error("GetPeer failed", "err", err, "peerId", targetId)
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
		// Synchronize block data with fetchBlock.
		cbft.fetchBlock(targetId, blockHash, blockNumber)
		return nil
	}
	//
	if msg.LogicType == network.TypeForQCBn {
		localQCNum, localQCHash := cbft.state.HighestQCBlock().NumberU64(), cbft.state.HighestQCBlock().Hash()
		if localQCNum == msg.BlockNumber && localQCHash == msg.BlockHash {
			cbft.log.Debug("Local qcBn is equal the sender's qcBn", "remoteBn", msg.BlockNumber, "localBn", localQCNum, "remoteHash", msg.BlockHash, "localHash", localQCHash)
			return nil
		}
		if localQCNum < msg.BlockNumber || (localQCNum == msg.BlockNumber && localQCHash != msg.BlockHash) {
			cbft.log.Debug("Local qcBn is larger than the sender's qcBn", "remoteBn", msg.BlockNumber, "localBn", localQCNum)
			localLockedNum, localLockedHash := cbft.state.HighestLockBlock().NumberU64(), cbft.state.HighestLockBlock().Hash()
			return launcher(msg.LogicType, id, localLockedNum, localLockedHash)
		} else {
			cbft.log.Debug("Local qcBn is less than the sender's qcBn", "remoteBn", msg.BlockNumber, "localBn", localQCNum)
			cbft.network.Send(id, &protocols.LatestStatus{BlockNumber: localQCNum, BlockHash: localQCHash, LogicType: msg.LogicType})
		}
	}
	// Deprecated
	if msg.LogicType == network.TypeForLockedBn {
		localLockedNum, localLockedHash := cbft.state.HighestLockBlock().NumberU64(), cbft.state.HighestLockBlock().Hash()
		if localLockedNum == msg.BlockNumber {
			cbft.log.Debug("Local lockedBn is equal the sender's lockedBn", "remoteBn", msg.BlockNumber, "localBn", localLockedNum)
			return nil
		}
		if localLockedNum < msg.BlockNumber {
			cbft.log.Debug("Local lockedBn is larger than the sender's lockedBn", "remoteBn", msg.BlockNumber, "localBn", localLockedNum)
			return launcher(msg.LogicType, id, localLockedNum, localLockedHash)
		} else {
			cbft.log.Debug("Local lockedBn is less than the sender's lockedBn", "remoteBn", msg.BlockNumber, "localBn", localLockedNum)
			cbft.network.Send(id, &protocols.LatestStatus{BlockNumber: localLockedNum, LogicType: msg.LogicType})
		}
	}
	// Deprecated
	if msg.LogicType == network.TypeForCommitBn {
		localCommitNum, localCommitHash := cbft.state.HighestCommitBlock().NumberU64(), cbft.state.HighestCommitBlock().Hash()
		if localCommitNum == msg.BlockNumber {
			cbft.log.Debug("Local commitBn is equal the sender's commitBn", "remoteBn", msg.BlockNumber, "localBn", localCommitNum)
			return nil
		}
		if localCommitNum < msg.BlockNumber {
			cbft.log.Debug("Local commitBn is larger than the sender's commitBn", "remoteBn", msg.BlockNumber, "localBn", localCommitNum)
			return launcher(msg.LogicType, id, localCommitNum, localCommitHash)
		} else {
			cbft.log.Debug("Local commitBn is less than the sender's commitBn", "remoteBn", msg.BlockNumber, "localBn", localCommitNum)
			cbft.network.Send(id, &protocols.LatestStatus{BlockNumber: localCommitNum, LogicType: msg.LogicType})
		}
	}
	return nil
}

// OnLatestStatus is used to process LatestStatus messages that received from peer.
func (cbft *Cbft) OnLatestStatus(id string, msg *protocols.LatestStatus) error {
	cbft.log.Debug("Received message on OnLatestStatus", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	switch msg.LogicType {
	case network.TypeForQCBn:
		localQCBn, localQCHash := cbft.state.HighestQCBlock().NumberU64(), cbft.state.HighestQCBlock().Hash()
		localLockedBn, localLockedHash := cbft.state.HighestLockBlock().NumberU64(), cbft.state.HighestLockBlock().Hash()
		if localQCBn < msg.BlockNumber || (localQCBn == msg.BlockNumber && localQCHash != msg.BlockHash) {
			p, err := cbft.network.GetPeer(id)
			if err != nil {
				cbft.log.Error("GetPeer failed", "err", err)
				return err
			}
			p.SetQcBn(new(big.Int).SetUint64(msg.BlockNumber))
			cbft.log.Debug("LocalQCBn is lower than sender's", "localBn", localQCBn, "remoteBn", msg.BlockNumber)
			cbft.fetchBlock(id, localLockedHash, localLockedBn)
		}

	case network.TypeForLockedBn:
		localLockedBn, localLockedHash := cbft.state.HighestLockBlock().NumberU64(), cbft.state.HighestLockBlock().Hash()
		if localLockedBn < msg.BlockNumber {
			p, err := cbft.network.GetPeer(id)
			if err != nil {
				cbft.log.Error("GetPeer failed", "err", err)
				return err
			}
			p.SetLockedBn(new(big.Int).SetUint64(msg.BlockNumber))
			cbft.log.Debug("LocalLockedBn is lower than sender's", "localBn", localLockedBn, "remoteBn", msg.BlockNumber)
			cbft.fetchBlock(id, localLockedHash, localLockedBn)
		}

	case network.TypeForCommitBn:
		localCommitBn, localCommitHash := cbft.state.HighestCommitBlock().NumberU64(), cbft.state.HighestCommitBlock().Hash()
		if localCommitBn < msg.BlockNumber {
			p, err := cbft.network.GetPeer(id)
			if err != nil {
				cbft.log.Error("GetPeer failed", "err", err)
				return err
			}
			p.SetCommitdBn(new(big.Int).SetUint64(msg.BlockNumber))
			cbft.log.Debug("LocalCommitBn is lower than sender's", "localBn", localCommitBn, "remoteBn", msg.BlockNumber)
			cbft.fetchBlock(id, localCommitHash, localCommitBn)
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
	cbft.log.Debug("Received message on OnPrepareBlockHash", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
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

// OnGetViewChange responds to nodes that require viewChange.
//
// The Epoch and viewNumber of viewChange must be consistent
// with the state of the current node.
func (cbft *Cbft) OnGetViewChange(id string, msg *protocols.GetViewChange) error {
	cbft.log.Debug("Received message on OnGetViewChange", "from", id, "msgHash", msg.MsgHash(), "message", msg.String(), "local", cbft.state.ViewString())

	localEpoch, localViewNumber := cbft.state.Epoch(), cbft.state.ViewNumber()

	isEqualLocalView := func() bool {
		return msg.ViewNumber == localViewNumber && msg.Epoch == localEpoch
	}

	isNextView := func() bool {
		return msg.ViewNumber+1 == localViewNumber || (msg.Epoch+1 == localEpoch && localViewNumber == state.DefaultViewNumber)
	}

	if isEqualLocalView() {
		// Get the viewChange belong to local node.
		node, err := cbft.validatorPool.GetValidatorByNodeID(cbft.state.HighestQCBlock().NumberU64(), cbft.config.Option.NodeID)
		if err != nil {
			cbft.log.Error("Get validator error, get view change failed", "err", err)
			return fmt.Errorf("get validator failed")
		}
		viewChanges := cbft.state.AllViewChange()
		if v, ok := viewChanges[uint32(node.Index)]; ok {
			cbft.network.Send(id, v)
			// Return if it contains missing.
			for _, nodeIndex := range msg.NodeIndexes {
				if v2, exists := viewChanges[nodeIndex]; exists {
					cbft.network.Send(id, v2)
				}
			}
		} else {
			cbft.log.Warn("No ViewChange found in current node")
		}
		return nil
	}
	// Return view QC in the case of less than 1.
	if isNextView() {
		lastViewChangeQC := cbft.state.LastViewChangeQC()
		if lastViewChangeQC == nil {
			cbft.log.Error("Not found lastViewChangeQC")
			return nil
		}
		err := lastViewChangeQC.EqualAll(msg.Epoch, msg.ViewNumber)
		if err != nil {
			cbft.log.Error("Last view change is not equal msg.viewNumber", "err", err)
			return err
		}
		cbft.network.Send(id, &protocols.ViewChangeQuorumCert{
			ViewChangeQC: lastViewChangeQC,
		})
		return nil
	}

	return fmt.Errorf("request is not match local view, local:%s,msg:%s", cbft.state.ViewString(), msg.String())
}

func (cbft *Cbft) OnViewChangeQuorumCert(id string, msg *protocols.ViewChangeQuorumCert) {
	cbft.log.Debug("Received message on OnViewChangeQuorumCert", "from", id, "msgHash", msg.MsgHash(), "message", msg.String())
	viewChangeQC := msg.ViewChangeQC
	epoch, viewNumber, _, _ := viewChangeQC.MaxBlock()
	if cbft.state.Epoch() == epoch && cbft.state.ViewNumber() == viewNumber {
		if err := cbft.verifyViewChangeQC(msg.ViewChangeQC); err == nil {
			cbft.tryChangeViewByViewChange(msg.ViewChangeQC)
		} else {
			cbft.log.Debug("Verify ViewChangeQC failed", "err", err)
		}
	}
}

// Returns the node ID of the missing vote.
func (cbft *Cbft) MissingViewChangeNodes() ([]discover.NodeID, *protocols.GetViewChange, error) {
	allViewChange := cbft.state.AllViewChange()
	nodeIds := make([]discover.NodeID, 0, len(allViewChange))
	qcBlockBn := cbft.state.HighestQCBlock().NumberU64()
	for k, _ := range allViewChange {
		nodeId := cbft.validatorPool.GetNodeIDByIndex(qcBlockBn, int(k))
		nodeIds = append(nodeIds, nodeId)
	}
	// all consensus
	consensusNodes, err := cbft.ConsensusNodes()
	if err != nil {
		return nil, nil, err
	}
	//consensusNodesLen := len(consensusNodes)
	missingNodes := make([]discover.NodeID, 0, len(consensusNodes)-len(nodeIds))
	for _, cv := range consensusNodes {
		isExists := false
		for _, v := range nodeIds {
			if cv == v {
				isExists = true
				break
			}
		}
		if !isExists {
			missingNodes = append(missingNodes, cv)
		}
	}

	log.Debug("Missing nodes on MissingViewChangeNodes", "nodes", network.FormatNodes(missingNodes))
	// The node of missingNodes must be in the list of neighbor nodes.
	peers, err := cbft.network.Peers()
	target := missingNodes[:0]
	for _, node := range missingNodes {
		for _, peer := range peers {
			if peer.ID() == node {
				target = append(target, node)
				break
			}
		}
	}
	log.Debug("Missing nodes exists in the peers", "nodes", network.FormatNodes(target))
	nodeIndexes := make([]uint32, 0, len(target))
	for _, v := range target {
		index, err := cbft.validatorPool.GetIndexByNodeID(qcBlockBn, v)
		if err != nil {
			continue
		}
		nodeIndexes = append(nodeIndexes, index)
	}
	cbft.log.Debug("Return missing node", "nodeIndexes", nodeIndexes)
	return target, &protocols.GetViewChange{
		Epoch:       cbft.state.Epoch(),
		ViewNumber:  cbft.state.ViewNumber(),
		NodeIndexes: nodeIndexes,
	}, nil
}

// OnPong is used to receive the average delay time.
func (cbft *Cbft) OnPong(nodeID string, netLatency int64) error {
	cbft.log.Trace("OnPong", "nodeID", nodeID, "netLatency", netLatency)
	cbft.netLatencyLock.Lock()
	defer cbft.netLatencyLock.Unlock()
	latencyList, exist := cbft.netLatencyMap[nodeID]
	if !exist {
		cbft.netLatencyMap[nodeID] = list.New()
		cbft.netLatencyMap[nodeID].PushBack(netLatency)
	} else {
		if latencyList.Len() > 5 {
			e := latencyList.Front()
			cbft.netLatencyMap[nodeID].Remove(e)
		}
		cbft.netLatencyMap[nodeID].PushBack(netLatency)
	}
	return nil
}

// AvgLatency returns the average delay time of the specified node.
//
// The average is the average delay between the current
// node and all consensus nodes.
// Return value unit: milliseconds.
func (cbft *Cbft) AvgLatency() time.Duration {
	cbft.netLatencyLock.Lock()
	defer cbft.netLatencyLock.Unlock()
	// The intersection of peerSets and consensusNodes.
	cNodes, _ := cbft.ConsensusNodes()
	peers, _ := cbft.network.Peers()
	target := make([]string, 0, len(peers))
	for _, pNode := range peers {
		for _, cNode := range cNodes {
			if pNode.PeerID() == cNode.TerminalString() {
				target = append(target, pNode.PeerID())
			}
		}
	}
	var (
		avgSum     int64 = 0
		result     int64 = 0
		validCount int64 = 0
	)
	// Take 2/3 nodes from the target.
	var pair utils.KeyValuePairList
	for _, v := range target {
		if latencyList, exist := cbft.netLatencyMap[v]; exist {
			avg := calAverage(latencyList)
			pair.Push(utils.KeyValuePair{Key: v, Value: avg})
		}
	}
	sort.Sort(pair)
	if pair.Len() == 0 {
		return time.Duration(0)
	}
	validCount = int64(pair.Len() * 2 / 3)
	if validCount == 0 {
		validCount = 1
	}
	for _, v := range pair[:validCount] {
		avgSum += v.Value
	}

	result = avgSum / validCount
	cbft.log.Debug("Get avg latency", "avg", result)
	return time.Duration(result) * time.Millisecond
}

func (cbft *Cbft) DefaultAvgLatency() time.Duration {
	return time.Duration(protocols.DefaultAvgLatency) * time.Millisecond
}

func calAverage(latencyList *list.List) int64 {
	var (
		sum    int64 = 0
		counts int64 = 0
	)
	for e := latencyList.Front(); e != nil; e = e.Next() {
		if latency, ok := e.Value.(int64); ok {
			counts++
			sum += latency
		}
	}
	if counts > 0 {
		return sum / counts
	}
	return 0
}
