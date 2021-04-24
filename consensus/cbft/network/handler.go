// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package network

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/hashicorp/golang-lru"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (

	// CbftProtocolName is protocol name of CBFT.
	CbftProtocolName = "cbft"

	// CbftProtocolVersion is protocol version of CBFT.
	CbftProtocolVersion = 1

	// CbftProtocolLength are the number of implemented message corresponding to cbft protocol versions.
	CbftProtocolLength = 40

	// sendQueueSize is maximum threshold for the queue of messages waiting to be sent.
	sendQueueSize = 10240

	// QCBnMonitorInterval is Qc block synchronization detection interval.
	QCBnMonitorInterval = 1

	// SyncViewChangeInterval is ViewChange synchronization detection interval.
	SyncViewChangeInterval = 1

	// SyncPrepareVoteInterval is PrepareVote synchronization detection interval.
	SyncPrepareVoteInterval = 1

	// removeBlacklistInterval is remove blacklist detection interval.
	removeBlacklistInterval = 20

	// TypeForQCBn is the type for QC sync.
	TypeForQCBn = 1

	// TypeForLockedBn is the type for Locked sync.
	TypeForLockedBn = 2

	// TypeForCommitBn is the type for Commit sync.
	TypeForCommitBn = 3

	// The maximum number of queues for message packets
	// that are communicated by peers.
	maxHistoryMessageHash = 5000

	// If the number of blacklists reaches the threshold,
	// the oldest blacklisted node will be re-trusted.
	maxBlacklist = 300
)

// EngineManager responsibles for processing the messages in the network.
type EngineManager struct {
	engine             Cbft
	router             *router
	peers              *PeerSet
	sendQueue          chan *types.MsgPackage
	quitSend           chan struct{}
	sendQueueHook      func(*types.MsgPackage)
	historyMessageHash *lru.ARCCache // Consensus message record that has been processed successfully.
	blacklist          *lru.Cache    // Save node blacklist.
}

// NewEngineManger returns a new handler and do some initialization.
func NewEngineManger(engine Cbft) *EngineManager {
	cache, err := lru.NewARC(maxHistoryMessageHash)
	if err != nil {
		return nil
	}
	handler := &EngineManager{
		engine:             engine,
		peers:              NewPeerSet(),
		sendQueue:          make(chan *types.MsgPackage, sendQueueSize),
		quitSend:           make(chan struct{}),
		historyMessageHash: cache,
	}
	handler.blacklist, _ = lru.New(maxBlacklist)
	// init router
	handler.router = newRouter(handler.Unregister, handler.getPeer, handler.ConsensusNodes, handler.peerList)
	return handler
}

// Start the loop to send message.
func (h *EngineManager) Start() {
	// Launch goroutine loop release separately.
	go h.sendLoop()
	go h.synchronize()
}

// Close turns off the handler for sending messages.
func (h *EngineManager) Close() {
	close(h.quitSend)
}

// The loop reads data from the message queue and sends it.
// If the message specifies the peerId then sends it directionally,
// and if the message does not specify peerId then broadcasts the message.
func (h *EngineManager) sendLoop() {
	for {
		select {
		case m := <-h.sendQueue:
			if h.sendQueueHook != nil {
				h.sendQueueHook(m)
			}
			if len(m.PeerID()) == 0 {
				h.broadcast(m)
			} else {
				h.sendMessage(m)
			}
		case <-h.quitSend:
			log.Warn("Terminate sending message")
			return
		}
	}
}

// Broadcast forwards the message to the router for distribution.
func (h *EngineManager) broadcast(m *types.MsgPackage) {
	h.router.Gossip(m)
}

// Send message to a known peerId. Determine if the peerId has established
// a connection before sending.
func (h *EngineManager) sendMessage(m *types.MsgPackage) {
	h.router.SendMessage(m)
}

// PeerSetting sets the block height of the node related type.
func (h *EngineManager) PeerSetting(peerID string, bType uint64, blockNumber uint64) error {
	p, err := h.peers.get(peerID)
	if err != nil || p == nil {
		return err
	}
	switch bType {
	case TypeForQCBn:
		p.SetQcBn(new(big.Int).SetUint64(blockNumber))
	case TypeForLockedBn:
		p.SetLockedBn(new(big.Int).SetUint64(blockNumber))
	case TypeForCommitBn:
		p.SetCommitdBn(new(big.Int).SetUint64(blockNumber))
	default:
	}
	return nil
}

// GetPeer returns the peer with the specified peerID.
func (h *EngineManager) getPeer(peerID string) (*peer, error) {
	if peerID == "" {
		return nil, fmt.Errorf("invalid peerID parameter - %v", peerID)
	}
	return h.peers.get(peerID)
}

// Send imports messages into the send queue and send it according to the specified ID.
func (h *EngineManager) Send(peerID string, msg types.Message) {
	msgPkg := types.NewMsgPackage(peerID, msg, types.NoneMode)
	select {
	case h.sendQueue <- msgPkg:
		log.Trace("Send message to sendQueue", "msgHash", msg.MsgHash(), "BHash", msg.BHash(), "msg", msg.String())
	default:
		log.Error("Send message failed, message queue blocking", "msgHash", msg.MsgHash(), "BHash", msg.BHash().TerminalString())
	}
}

// Broadcast imports messages into the send queue and send it according to broadcast.
//
// Note: The broadcast of this method defaults to FULL mode.
func (h *EngineManager) Broadcast(msg types.Message) {
	msgPkg := types.NewMsgPackage("", msg, types.FullMode)
	select {
	case h.sendQueue <- msgPkg:
		log.Trace("Broadcast message to sendQueue", "msgHash", msg.MsgHash(), "BHash", msg.BHash().TerminalString(), "msg", msg.String())
	default:
		log.Error("Broadcast message failed, message queue blocking", "msgHash", msg.MsgHash(), "BHash", msg.BHash().TerminalString())
	}
}

// PartBroadcast imports messages into the send queue.
//
// Note: The broadcast of this method defaults to PartMode.
func (h *EngineManager) PartBroadcast(msg types.Message) {
	msgPkg := types.NewMsgPackage("", msg, types.PartMode)
	select {
	case h.sendQueue <- msgPkg:
		log.Debug("PartBroadcast message to sendQueue", "msgHash", msg.MsgHash().TerminalString(), "BHash", msg.BHash().TerminalString(), "msg", msg.String())
	default:
		log.Error("PartBroadcast message failed, message queue blocking", "msgHash", msg.MsgHash().TerminalString(), "BHash", msg.BHash().TerminalString())
	}
}

// Forwarding is used to forward the messages and determine
// whether forwarding is required according to the message type.
//
// Note:
// 1. message type that need to be forwarded:
//    PrepareBlockMsg/PrepareVoteMsg/ViewChangeMsg/BlockQuorumCertMsg
// 2. message type that need not to be forwarded:
//    (Except for the above types, the rest are not forwarded).
func (h *EngineManager) Forwarding(nodeID string, msg types.Message) error {
	msgHash := msg.MsgHash()
	msgType := protocols.MessageType(msg)

	// the logic to forward message.
	forward := func() error {
		peers, err := h.peerList()
		if err != nil || len(peers) == 0 {
			return fmt.Errorf("peers is none, msgHash:%s", msgHash.TerminalString())
		}
		// Check all neighbor node lists and see if the specified message has been processed.
		for _, peer := range peers {
			// exclude currently send peer.
			if peer.id == nodeID {
				continue
			}
			if peer.ContainsMessageHash(msgHash) {
				messageRepeatMeter.Mark(1)
				log.Trace("Needn't to broadcast", "type", reflect.TypeOf(msg), "hash", msgHash.TerminalString(), "BHash", msg.BHash().TerminalString())
				return fmt.Errorf("contain message and not formard, msgHash:%s", msgHash.TerminalString())
			}
		}
		log.Debug("Need to broadcast", "type", reflect.TypeOf(msg), "hash", msgHash.TerminalString(), "BHash", msg.BHash().TerminalString())
		// Need to broadcast the message.
		// For PrepareBlockMsg messages, here are some differences:
		// 1.PrepareBlock does not forward directly but sends its hash (PrepareBlockHash).
		if msgType == protocols.PrepareBlockMsg {
			// Special treatment.
			if v, ok := msg.(*protocols.PrepareBlock); ok {
				pbh := &protocols.PrepareBlockHash{
					Epoch:       v.Epoch,
					ViewNumber:  v.ViewNumber,
					BlockIndex:  v.BlockIndex,
					BlockHash:   v.Block.Hash(),
					BlockNumber: v.Block.NumberU64(),
				}
				h.Broadcast(pbh)
				log.Debug("PrepareBlockHash is forwarded instead of PrepareBlock", "msgHash", pbh.MsgHash())
			}
		} else {
			// Direct forwarding.
			h.Broadcast(msg)
		}
		return nil
	}
	// PrepareBlockMsg does not forward, the message will be forwarded using PrepareBlockHash.
	switch msgType {
	case protocols.PrepareBlockMsg, protocols.PrepareVoteMsg, protocols.ViewChangeMsg:
		err := forward()
		if err != nil {
			messageGossipMeter.Mark(1)
		}
		return err
	default:
		log.Trace("Unmatched message type, need not to be forwarded", "type", reflect.TypeOf(msg), "msgHash", msgHash.TerminalString(), "BHash", msg.BHash().TerminalString())
	}
	return nil
}

// Protocols implemented the Protocols method and returned basic information about the CBFT protocol.
func (h *EngineManager) Protocols() []p2p.Protocol {
	return []p2p.Protocol{
		{
			Name:    CbftProtocolName,
			Version: CbftProtocolVersion,
			Length:  CbftProtocolLength,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				return h.handler(p, rw)
			},
			NodeInfo: func() interface{} {
				return h.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p, err := h.peers.get(fmt.Sprintf("%x", id[:8])); err == nil {
					return p.Info()
				}
				return nil
			},
		},
	}
}

// AliveConsensusNodeIDs returns all NodeID to alive peer.
func (h *EngineManager) AliveConsensusNodeIDs() ([]string, error) {
	cNodes, _ := h.engine.ConsensusNodes()
	peers := h.peers.allPeers()
	target := make([]string, 0, len(peers))
	for _, pNode := range peers {
		for _, cNode := range cNodes {
			if pNode.PeerID() == cNode.TerminalString() {
				target = append(target, pNode.PeerID())
			}
		}
	}
	return target, nil
}

// Return all neighbor node lists.
func (h *EngineManager) peerList() ([]*peer, error) {
	return h.peers.allPeers(), nil
}

// Unregister removes the peer with the specified ID.
func (h *EngineManager) Unregister(id string) error {
	return h.peers.Unregister(id)
}

// ConsensusNodes returns a list of all consensus nodes.
func (h *EngineManager) ConsensusNodes() ([]discover.NodeID, error) {
	return h.engine.ConsensusNodes()
}

// NodeInfo representatives node configuration information.
type NodeInfo struct {
	Config types.Config `json:"config"`
}

// NodeInfo returns the information of Node.
func (h *EngineManager) NodeInfo() *NodeInfo {
	cfg := h.engine.Config()
	return &NodeInfo{
		Config: *cfg,
	}
}

// After the node is successfully connected and the message belongs
// to the cbft protocol message, the method is called.
func (h *EngineManager) handler(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	// Further confirm if the version number needs to be read from the configuration.
	peer := newPeer(CbftProtocolVersion, p, newMeteredMsgWriter(rw))

	// execute handshake
	// 1.need qcBn/qcHash/lockedBn/lockedHash/commitBn/commitHash from cbft.
	var (
		qcBn, qcHash         = h.engine.HighestQCBlockBn()
		lockedBn, lockedHash = h.engine.HighestLockBlockBn()
		commitBn, commitHash = h.engine.HighestCommitBlockBn()
	)
	p.Log().Debug("CBFT peer connected, do handshake", "name", peer.Name())

	// Execution handshake function.
	handshake := func() error {
		// Build a new CbftStatusData object as a handshake parameter
		cbftStatus := &protocols.CbftStatusData{
			ProtocolVersion: CbftProtocolVersion,
			QCBn:            new(big.Int).SetUint64(uint64(qcBn)),
			QCBlock:         qcHash,
			LockBn:          new(big.Int).SetUint64(uint64(lockedBn)),
			LockBlock:       lockedHash,
			CmtBn:           new(big.Int).SetUint64(uint64(commitBn)),
			CmtBlock:        commitHash,
		}
		// do handshake
		remoteStatus, err := peer.Handshake(cbftStatus)
		if err != nil {
			p.Log().Debug("CBFT handshake failed", "err", err)
			return err
		}

		// Blacklist check.
		if h.ContainsBlacklist(peer.PeerID()) {
			p.Log().Error("CBFT handshake, peer that are forbidden to connect")
			return fmt.Errorf("illegal node: {%s}", peer.PeerID())
		}

		// If blockNumber in the local is better than the remote
		// then determine if there is a fork.
		if cbftStatus.QCBn.Uint64() > remoteStatus.QCBn.Uint64() {
			err = h.engine.BlockExists(remoteStatus.QCBn.Uint64(), remoteStatus.QCBlock)
		}
		if cbftStatus.LockBn.Uint64() > remoteStatus.LockBn.Uint64() {
			err = h.engine.BlockExists(remoteStatus.LockBn.Uint64(), remoteStatus.LockBlock)
		}
		if cbftStatus.CmtBn.Uint64() > remoteStatus.CmtBn.Uint64() {
			err = h.engine.BlockExists(remoteStatus.CmtBn.Uint64(), remoteStatus.CmtBlock)
		}
		if err != nil {
			p.Log().Error("CBFT handshake, verify block failed", "err", err)
			return err
		}
		p.Log().Debug("CBFT consensus handshake success", "msgHash", cbftStatus.MsgHash().TerminalString())
		return nil
	}
	handErr := handshake()
	if handErr != nil {
		p.Log().Error("CBFT handshake failed", "err", handErr)
		return handErr
	}

	// The newly established node is registered to the neighbor node list.
	if err := h.peers.Register(peer); err != nil {
		p.Log().Error("Cbft peer registration failed", "err", err)
		return err
	}
	defer h.RemovePeer(peer.PeerID())

	// start ping loop.
	go peer.Run()

	// main loop. handle incoming message.
	// Exit the loop and disconnect if the message
	// is processing abnormally.
	for {
		if err := h.handleMsg(peer); err != nil {
			p.Log().Error("CBFT message handling failed", "peerID", peer.PeerID(), "err", err)
			return err
		}
	}
}

// Main logic: Distribute according to message type and
// transfer message to CBFT layer
func (h *EngineManager) handleMsg(p *peer) error {
	msg, err := p.ReadWriter().ReadMsg()
	if err != nil {
		p.Log().Error("Read peer message error", "err", err)
		return err
	}

	// All messages cannot exceed the maximum specified by the agreement.
	if msg.Size > protocols.CbftProtocolMaxMsgSize {
		return types.ErrResp(types.ErrMsgTooLarge, "%v > %v", msg.Size, protocols.CbftProtocolMaxMsgSize)
	}
	defer msg.Discard()

	// Handle the message depending on msgType and it's content.
	switch {
	case msg.Code == protocols.CBFTStatusMsg:
		// CBFTStatusMsg belongs to the type of handshake message and will not appear here.
		return types.ErrResp(types.ErrExtraStatusMsg, "uncontrolled status message")

	case msg.Code == protocols.PrepareBlockMsg:
		var request protocols.PrepareBlock
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		p.MarkMessageHash((&request).MsgHash())
		request.Block.ReceivedAt = msg.ReceivedAt
		request.Block.ReceivedFrom = p
		// Message transfer to cbft message queue.
		return h.engine.ReceiveMessage(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.PrepareVoteMsg:
		var request protocols.PrepareVote
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		p.MarkMessageHash((&request).MsgHash())
		return h.engine.ReceiveMessage(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.ViewChangeMsg:
		var request protocols.ViewChange
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		p.MarkMessageHash((&request).MsgHash())
		return h.engine.ReceiveMessage(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.GetPrepareBlockMsg:
		var request protocols.GetPrepareBlock
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.GetBlockQuorumCertMsg:
		var request protocols.GetBlockQuorumCert
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.BlockQuorumCertMsg:
		var request protocols.BlockQuorumCert
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		p.MarkMessageHash((&request).MsgHash())
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.GetQCBlockListMsg:
		var request protocols.GetQCBlockList
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.GetPrepareVoteMsg:
		var request protocols.GetPrepareVote
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.PrepareBlockHashMsg:
		var request protocols.PrepareBlockHash
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		p.MarkMessageHash((&request).MsgHash())
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.PrepareVotesMsg:
		var request protocols.PrepareVotes
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.QCBlockListMsg:
		var request protocols.QCBlockList
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.GetLatestStatusMsg:
		var request protocols.GetLatestStatus
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.LatestStatusMsg:
		var request protocols.LatestStatus
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.GetViewChangeMsg:
		var request protocols.GetViewChange
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	case msg.Code == protocols.PingMsg:
		var pingTime protocols.Ping
		if err := msg.Decode(&pingTime); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		// Directly respond to the response message after receiving the ping message.
		go p2p.SendItems(p.ReadWriter(), protocols.PongMsg, pingTime[0])
		p.Log().Trace("Respond to ping message done")
		return nil

	case msg.Code == protocols.PongMsg:
		// Processed after receiving the pong message.
		curTime := time.Now().UnixNano()
		log.Debug("Handle a eth Pong message", "curTime", curTime)
		var pongTime protocols.Pong
		if err := msg.Decode(&pongTime); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		for {
			// Return the first element of list l or nil if the list is empty.
			frontPing := p.ListFront()
			if frontPing == nil {
				log.Trace("End of p.pingList")
				break
			}
			log.Trace("Front element of p.pingList", "element", frontPing)
			if t, ok := p.ListRemove(frontPing).(string); ok {
				if t == pongTime[0] {
					tInt64, err := strconv.ParseInt(t, 10, 64)
					if err != nil {
						return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
					}
					log.Trace("Calculate net latency", "sendPingTime", tInt64, "receivePongTime", curTime)
					latency := (curTime - tInt64) / 2 / 1000000
					// Record the latency in metrics and output it. unit: second.
					log.Trace("Latency", "time", latency)
					h.engine.OnPong(p.id, latency)
					propPeerLatencyMeter.Mark(latency)
					break
				}
			}
		}
		return nil
	case msg.Code == protocols.ViewChangeQuorumCertMsg:
		var request protocols.ViewChangeQuorumCert
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))
	case msg.Code == protocols.ViewChangesMsg:
		var request protocols.ViewChanges
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return h.engine.ReceiveSyncMsg(types.NewMsgInfo(&request, p.PeerID()))

	default:
		return types.ErrResp(types.ErrInvalidMsgCode, "%v", msg.Code)
	}
}

// MarkHistoryMessageHash is used to record the hash value of each message from the peer node.
// If the queue is full, remove the bottom element and add a new one.
func (h *EngineManager) MarkHistoryMessageHash(hash common.Hash) {
	h.historyMessageHash.Add(hash, struct{}{})
}

// ContainsMessageHash returns whether the specified hash exists.
func (h *EngineManager) ContainsHistoryMessageHash(hash common.Hash) bool {
	return h.historyMessageHash.Contains(hash)
}

// MarkBlacklist marks the specified node as a blacklist.
// If the number of recorded blacklists reaches the threshold,
// the node that was first set to blacklist will be removed from the blacklist.
func (h *EngineManager) MarkBlacklist(peerID string) {
	deadline := time.Duration(h.engine.Config().Option.BlacklistDeadline) * time.Minute
	h.blacklist.Add(peerID, time.Now().Add(deadline))
}

// ContainsBlacklist returns whether the specified node is blacklisted.
func (h *EngineManager) ContainsBlacklist(peerID string) bool {
	return h.blacklist.Contains(peerID)
}

// RemoveMessageHash removes the specified hash from the peer.
func (h *EngineManager) RemoveMessageHash(id string, msgHash common.Hash) {
	peer, err := h.peers.get(id)
	if err != nil {
		log.Error("Removing messageHash from peer failed", "err", err)
		return
	}
	peer.RemoveMessageHash(msgHash)
}

// RemovePeer removes and disconnects a node from
// a neighbor node.
func (h *EngineManager) RemovePeer(id string) {
	// Short circuit if the peer was already removed
	peer, err := h.peers.get(id)
	if err != nil {
		log.Error("Removing CBFT peer failed", "err", err)
		return
	}
	log.Debug("Removing CBFT peer", "peer", id)

	if err := h.peers.Unregister(id); err != nil {
		log.Error("CBFT Peer removal failed", "peer", id, "err", err)
	}
	// Hard disconnect at the networking layer
	if peer != nil {
		peer.Peer.Disconnect(p2p.DiscUselessPeer)
	}
}

// Select a node with a height higher than the local node block from
// the neighbor node list, and then synchronize the block data of
// the height difference to the node.
//
// Note:
// 1. Synchronous blocks with inconsistent QC height.
// 2. Synchronous blocks with inconsistent locking block height.
// 3. Synchronous blocks with inconsistent commit block height.
func (h *EngineManager) synchronize() {
	log.Debug("~ Start synchronize in the handler")
	blockNumberTimer := time.NewTimer(QCBnMonitorInterval * time.Second)
	viewTicker := time.NewTicker(SyncViewChangeInterval * time.Second)
	pureBlacklistTicker := time.NewTicker(removeBlacklistInterval * time.Second)
	voteTicker := time.NewTicker(SyncPrepareVoteInterval * time.Second)

	// Logic used to synchronize QC.
	syncQCBnFunc := func() {
		latestStatus := h.engine.LatestStatus()
		log.Debug("Synchronize for qc block send message", "latestStatus", latestStatus.String())
		latestStatus.LogicType = TypeForQCBn
		h.PartBroadcast(latestStatus)
	}

	for {
		select {
		case <-voteTicker.C:
			msg, err := h.engine.MissingPrepareVote()
			if err != nil {
				log.Debug("Request missing prepareVote failed", "err", err)
				break
			}
			log.Debug("Had new prepareVote sync request", "msg", msg.String())
			// Only broadcasts without forwarding.
			h.PartBroadcast(msg)

		case <-blockNumberTimer.C:
			// Sent at random.
			syncQCBnFunc()
			rd := rand.Intn(5)
			if rd == 0 || rd < QCBnMonitorInterval/2 {
				rd = (rd + 1) * 2
			}
			resetTime := time.Duration(rd) * time.Second
			blockNumberTimer.Reset(resetTime)

		case <-viewTicker.C:
			// If the local viewChange has insufficient votes,
			// the GetViewChange message is sent from the missing node.
			msg, err := h.engine.MissingViewChangeNodes()
			if err != nil {
				log.Debug("Request missing viewchange failed", "err", err)
				break
			}
			log.Debug("Had new viewchange sync request", "msg", msg.String())
			// Only broadcasts without forwarding.
			h.PartBroadcast(msg)

		case <-pureBlacklistTicker.C:
			// Iterate over the blacklist and remove
			// the nodes that have expired.
			keys := h.blacklist.Keys()
			log.Debug("Blacklist pure start", "len", len(keys))
			for _, k := range keys {
				v, exists := h.blacklist.Get(k)
				if !exists {
					continue
				}
				if t, ok := v.(time.Time); ok {
					if t.Before(time.Now()) {
						h.blacklist.Remove(k)
						log.Debug("Remove blacklist success", "peerID", k)
					}
				}
			}

		case <-h.quitSend:
			log.Warn("Synchronize quit")
			return
		}
	}
}

// Select a node from the list of nodes that is larger than the specified value.
//
// bType: 1 -> qcBlock, 2 -> lockedBlock, 3 -> CommitBlock
func largerPeer(bType uint64, peers []*peer, number uint64) (*peer, uint64) {
	if len(peers) == 0 {
		return nil, 0
	}
	largerNum := number
	largerIndex := -1
	for index, v := range peers {
		var pNumber uint64
		switch bType {
		case TypeForQCBn:
			pNumber = v.QCBn()
		case TypeForLockedBn:
			pNumber = v.LockedBn()
		case TypeForCommitBn:
			pNumber = v.CommitBn()
		default:
			return nil, 0
		}
		if pNumber > largerNum {
			largerNum, largerIndex = pNumber, index
		}
	}
	if largerIndex != -1 {
		return peers[largerIndex], largerNum
	}
	return nil, 0
}

// Testing is only used for unit testing.
func (h *EngineManager) Testing() {
	peers, _ := h.peerList()
	for _, v := range peers {
		v.Run()
		go func(p *peer) {
			for {
				if err := h.handleMsg(p); err != nil {
					p.Log().Error("In the testing, CBFT message handling failed", "err", err)
					break
				}
			}
		}(v)
	}
}
