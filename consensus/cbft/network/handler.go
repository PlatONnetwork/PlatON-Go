package network

import (
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strconv"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (

	// Protocol name of CBFT
	CbftProtocolName = "cbft"

	// Protocol version of CBFT
	CbftProtocolVersion = 1

	// CbftProtocolLength are the number of implemented message corresponding to cbft protocol versions.
	CbftProtocolLength = 20

	// Maximum threshold for the queue of messages waiting to be sent.
	sendQueueSize = 10240

	QCBnMonitorInterval = 10 // Qc block synchronization detection interval
	//LockedBnMonitorInterval = 4 // Locked block synchronization detection interval
	//CommitBnMonitorInterval = 4 // Commit block synchronization detection interval
	SyncViewChangeInterval = 15

	//
	TypeForQCBn     = 1
	TypeForLockedBn = 2
	TypeForCommitBn = 3
)

// Responsible for processing the messages in the network.
type EngineManager struct {
	engine        Cbft
	router        *router
	peers         *PeerSet
	sendQueue     chan *types.MsgPackage
	quitSend      chan struct{}
	sendQueueHook func(*types.MsgPackage)
}

// Create a new handler and do some initialization.
func NewEngineManger(engine Cbft) *EngineManager {
	handler := &EngineManager{
		engine:    engine,
		peers:     NewPeerSet(),
		sendQueue: make(chan *types.MsgPackage, sendQueueSize),
		quitSend:  make(chan struct{}, 0),
	}
	// init router
	handler.router = NewRouter(handler.Unregister, handler.GetPeer, handler.ConsensusNodes, handler.Peers)
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
			// todo: Need to add to the processing judgment of wal
			if len(m.PeerID()) == 0 {
				h.broadcast(m)
			} else {
				h.sendMessage(m)
			}
		case <-h.quitSend:
			log.Error("Terminate sending message")
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

// Return the peer with the specified peerID.
func (h *EngineManager) GetPeer(peerID string) (*peer, error) {
	if peerID == "" {
		return nil, fmt.Errorf("invalid peerID parameter - %v", peerID)
	}
	return h.peers.Get(peerID)
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
	// todo: Whether to consider the problem of blocking
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

// Broadcast imports messages into the send queue.
//
// Note: The broadcast of this method defaults to PartMode.
func (h *EngineManager) PartBroadcast(msg types.Message) {
	msgPkg := types.NewMsgPackage("", msg, types.PartMode)
	select {
	case h.sendQueue <- msgPkg:
		log.Debug("PartBroadcast message to sendQueue", "msgHash", msg.MsgHash().TerminalString(), "BHash", msg.BHash().TerminalString())
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
func (h *EngineManager) Forwarding(nodeId string, msg types.Message) error {
	msgHash := msg.MsgHash()
	msgType := protocols.MessageType(msg)

	// the logic to forward message.
	forward := func() error {
		peers, err := h.Peers()
		if err != nil || len(peers) == 0 {
			return fmt.Errorf("peers is none, msgHash:%s", msgHash.TerminalString())
		}
		// Check all neighbor node lists and see if the specified message has been processed.
		for _, peer := range peers {
			// exclude currently send peer.
			if peer.id == nodeId {
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
				go h.Broadcast(&protocols.PrepareBlockHash{
					Epoch:       v.Epoch,
					ViewNumber:  v.ViewNumber,
					BlockIndex:  v.BlockIndex,
					BlockHash:   v.Block.Hash(),
					BlockNumber: v.Block.NumberU64(),
				})
				log.Trace("PrepareBlockHash is forwarded instead of PrepareBlock done")
			}
		} else {
			// Direct forwarding.
			go h.Broadcast(msg)
		}
		return nil
	}
	// PrepareBlockMsg does not forward, the message will be forwarded using PrepareBlockHash.
	switch msgType {
	case protocols.PrepareBlockMsg, protocols.PrepareVoteMsg, protocols.ViewChangeMsg,
		protocols.BlockQuorumCertMsg, protocols.PrepareBlockHashMsg:
		err := forward()
		if err != nil {
			messageGossipMeter.Mark(1)
		}
		return err
	default:
		log.Warn("Unmatched message type, need not to be forwarded", "type", reflect.TypeOf(msg), "msgHash", msgHash.TerminalString(), "BHash", msg.BHash().TerminalString())
	}
	return nil
}

// Protocols implemented the Protocols method and returned basic information about the CBFT protocol.
func (h *EngineManager) Protocols() []p2p.Protocol {
	// todo: version and ProtocolLengths need to confirm.
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
				if p, err := h.peers.Get(fmt.Sprintf("%x", id[:8])); err == nil {
					return p.Info()
				}
				return nil
			},
		},
	}
}

// Return all neighbor node lists.
func (h *EngineManager) Peers() ([]*peer, error) {
	return h.peers.Peers(), nil
}

// Remove the peer with the specified ID
func (h *EngineManager) Unregister(id string) error {
	return h.peers.Unregister(id)
}

// Return a list of all consensus nodes
func (h *EngineManager) ConsensusNodes() ([]discover.NodeID, error) {
	return h.engine.ConsensusNodes()
}

// Representative node configuration information.
type NodeInfo struct {
	Config types.Config `json:"config"`
}

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
	peer := NewPeer(CbftProtocolVersion, p, newMeteredMsgWriter(rw))

	// execute handshake
	// todo:
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
		// If blockNumber in the local is better than the remote
		// then determine if there is a fork.
		if cbftStatus.QCBn.Uint64() > remoteStatus.QCBn.Uint64() {
			// todo: to be added
		}
		if cbftStatus.LockBn.Uint64() > remoteStatus.LockBn.Uint64() {
			// todo: to be added
		}
		if cbftStatus.CmtBn.Uint64() > remoteStatus.CmtBn.Uint64() {
			// todo: to be added
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
	defer h.peers.Unregister(peer.PeerID())

	// start ping loop.
	go peer.Run()

	// main loop. handle incoming message.
	// Exit the loop and disconnect if the message
	// is processing abnormally.
	for {
		if err := h.handleMsg(peer); err != nil {
			p.Log().Error("CBFT message handling failed", "err", err)
			return err
		}
	}
}

// Main logic: Distribute according to message type and
// transfer message to CBFT layer
func (h *EngineManager) handleMsg(p *peer) error {
	msg, err := p.ReadWriter().ReadMsg()
	if err != nil {
		p.Log().Error("read peer message error", "err", err)
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
		log.Debug("handle a eth Pong message", "curTime", curTime)
		var pingTime protocols.Pong
		if err := msg.Decode(&pingTime); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		for {
			// Return the first element of list l or nil if the list is empty.
			frontPing := p.PingList.Front()
			if frontPing == nil {
				log.Trace("end of p.PingList")
				break
			}
			log.Trace("Front element of p.PingList", "element", frontPing)
			if t, ok := p.PingList.Remove(frontPing).(string); ok {
				if t == pingTime[0] {
					tInt64, err := strconv.ParseInt(t, 10, 64)
					if err != nil {
						return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
					}
					log.Trace("calculate net latency", "sendPingTime", tInt64, "receivePongTime", curTime)
					latency := (curTime - tInt64) / 2 / 1000000
					// todo: need confirm
					// Record the latency in metrics and output it. unit: second.
					log.Trace("latency", "time", latency)
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

	default:
		return types.ErrResp(types.ErrInvalidMsgCode, "%v", msg.Code)
	}

	return nil
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
	blockNumberTicker := time.NewTicker(QCBnMonitorInterval * time.Second)
	viewTicker := time.NewTicker(SyncViewChangeInterval * time.Second)

	// Logic used to synchronize QC.
	syncQCBnFunc := func() {
		qcBn, _ := h.engine.HighestQCBlockBn()
		log.Debug("Synchronize for qc block send message", "localQCBn", qcBn)
		h.PartBroadcast(&protocols.GetLatestStatus{
			BlockNumber: qcBn,
			LogicType:   TypeForQCBn,
		})
	}

	// Update if it is the same state within 5 seconds
	var (
		lastEpoch      uint64 = 0
		lastViewNumber uint64 = 0
	)

	for {
		select {
		case <-blockNumberTicker.C:
			// Sent at random.
			randomSend(syncQCBnFunc)

		case <-viewTicker.C:
			// If the local viewChange has insufficient votes,
			// the GetViewChange message is sent from the missing node.
			missingViewNodes, msg, err := h.engine.MissingViewChangeNodes()
			if err != nil {
				log.Error("Get consensus nodes failed", "err", err)
				break
			}
			// Initi.al situation.
			if lastEpoch == msg.Epoch && lastViewNumber == msg.ViewNumber {
				log.Debug("Will send GetViewChange", "missingNodes", FormatNodes(missingViewNodes))
				// Only broadcasts without forwarding.
				h.Broadcast(msg)
			} else {
				log.Debug("Waiting for the next round")
			}
			lastEpoch, lastViewNumber = msg.Epoch, msg.ViewNumber

		case <-h.quitSend:
			log.Error("synchronize quit")
			return
		}
	}
}

// Randomly sent during the timer period.
func randomSend(exec func()) {
	sleepTime := rand.Intn(QCBnMonitorInterval)
	for {
		if sleepTime > QCBnMonitorInterval/2 {
			break
		}
		sleepTime *= 2
	}
	time.AfterFunc(time.Duration(int64(sleepTime)), func() {
		exec()
	})
}

// Select a node from the list of nodes that is larger than the specified value.
//
// bType:
//  1 -> qcBlock, 2 -> lockedBlock, 3 -> CommitBlock
func largerPeer(bType uint64, peers []*peer, number uint64) (*peer, uint64) {
	if peers == nil || len(peers) == 0 {
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
	peers, _ := h.Peers()
	for _, v := range peers {
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
