package cbft

import (
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/router"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// Maximum threshold for the queue of messages waiting to be sent.
const sendQueueSize = 10240

const (
	CbftProtocolVersion = 1  // Protocol version of CBFT
	CbftProtocolLength  = 15 // CbftProtocolLength are the number of implemented message corresponding to cbft protocol versions.
)

// Responsible for processing the messages in the network.
type EngineManager struct {
	engine    *Cbft
	peers     *router.PeerSet
	router    Router
	sendQueue chan *types.MsgPackage
	quitSend  chan struct{}
}

// Create a new handler and do some initialization.
func NewEngineManger(engine *Cbft, r Router) *EngineManager {
	return &EngineManager{
		engine:    engine,
		peers:     router.NewPeerSet(),
		router:    r,
		sendQueue: make(chan *types.MsgPackage, sendQueueSize),
		quitSend:  make(chan struct{}, 0),
	}
}

// Start the loop to send message.
func (h *EngineManager) Start() {
	// Launch goroutine loop release separately.
	go h.sendLoop()
}

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
			if len(m.PeerID()) == 0 {
				h.broadcast(m)
			} else {
				h.sendMessage(m)
			}
		case <-h.quitSend:
			log.Error("Terminate sending message")
			break
		}
	}
}

// Broadcast forwards the message to the router for distribution.
func (h *EngineManager) broadcast(m *types.MsgPackage) {
	h.router.gossip(m)
}

// Send message to a known peerId. Determine if the peerId has established
// a connection before sending.
func (h *EngineManager) sendMessage(m *types.MsgPackage) {
	h.router.sendMessage(m)
}

// Return the peer with the specified peerID.
func (h *EngineManager) GetPeer(peerID string) (*router.Peer, error) {
	if peerID == "" {
		return nil, fmt.Errorf("invalid peerID parameter - %v", peerID)
	}
	return h.peers.Get(peerID)
}

// Send imports messages into the send queue and send it according to the specified ID.
func (h *EngineManager) Send(peerID discover.NodeID, msg types.Message) {
	msgPkg := types.NewMsgPackage(peerID.TerminalString(), msg, types.FullMode)
	select {
	case h.sendQueue <- msgPkg:
		log.Debug("Send message to sendQueue", "msgHash", msg.MsgHash().TerminalString(), "BHash", msg.BHash().TerminalString())
	default:
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
		log.Debug("Broadcast message to sendQueue", "msgHash", msg.MsgHash().TerminalString(), "BHash", msg.BHash().TerminalString())
	default:
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
	}
}

// Protocols implemented the Protocols method and returned basic information about the CBFT protocol.
func (h *EngineManager) Protocols() []p2p.Protocol {
	// todo: version and ProtocolLengths need to confirm.
	return []p2p.Protocol{
		{
			Name:    "cbft",
			Version: CbftProtocolVersion,
			Length:  CbftProtocolLength,
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				return h.handler(p, rw)
			},
			NodeInfo: func() interface{} {
				return h.NodeInfo()
			},
			PeerInfo: func(id discover.NodeID) interface{} {
				if p, err := h.peers.Get(fmt.Sprintf("%5x", id[:8])); err == nil {
					return p.Info()
				}
				return nil
			},
		},
	}
}

// Representative node configuration information.
type NodeInfo struct {
	Config Config `json:"config"`
}

func (h *EngineManager) NodeInfo() *NodeInfo {
	// todo: Use methods instead of properties to access directly
	return nil
	/*cfg := h.engine.config
	return &NodeInfo{
		Config: cfg,
	}*/
}

// After the node is successfully connected and the message belongs
// to the cbft protocol message, the method is called.
func (h *EngineManager) handler(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	// Further confirm if the version number needs to be read from the configuration.
	peer := router.NewPeer(CbftProtocolVersion, p, newMeteredMsgWriter(rw))

	// Execute the CBFT handshake
	// todo:
	// 1.need qcBn/qcHash/lockedBn/lockedHash/commitBn/commitHash from cbft.
	var (
		qcBn       = 1
		qcHash     = common.Hash{}
		lockedBn   = 1
		lockedHash = common.Hash{}
		commitBn   = 1
		commitHash = common.Hash{}
	)
	p.Log().Debug("CBFT peer connected, do handshake", "name", peer.Name())

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
	if err := peer.Handshake(cbftStatus); err != nil {
		p.Log().Debug("CBFT handshake failed", "err", err)
		return err
	} else {
		p.Log().Debug("CBFT consensus handshake success", "msgHash", cbftStatus.MsgHash().TerminalString())
	}

	// The newly established node is registered to the neighbor node list.
	if err := h.peers.Register(peer); err != nil {
		p.Log().Error("Cbft peer registration failed", "err", err)
		return err
	}
	defer h.peers.Unregister(peer.PeerID())

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

func (h *EngineManager) handleMsg(p *router.Peer) error {
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
		return types.ErrResp(types.ErrExtraStatusMsg, "uncontrolled status message")
	case msg.Code == protocols.GetPrepareBlockMsg:
		// todo: GetPrepareBlockMsg need to process.
		return nil
	case msg.Code == protocols.PrepareBlockMsg:
		var request protocols.PrepareBlock
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		p.MarkMessageHash((&request).MsgHash())
		request.Block.ReceivedAt = msg.ReceivedAt
		request.Block.ReceivedFrom = p
		// Message transfer to cbft message queue.
		// todo: need to process.
		// h.engine.ReceivePeerMsg(types.NewMessageInfo(&request, p.ID()))
		return nil
	case msg.Code == protocols.PrepareVoteMsg:
		var request protocols.PrepareVote
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return nil

	case msg.Code == protocols.ViewChangeMsg:
		/*var request viewChange
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		if h.cbft.isForwarded(p.ID(), &request) {
			return nil
		}
		p.MarkMessageHash((&request).MsgHash())
		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})*/
		return nil
	case msg.Code == protocols.PrepareBlockHashMsg:
		/*var request protocols.PrepareBlockHash
		if err := msg.Decode(&request); err != nil {
			return types.ErrResp(types.ErrDecode, "%v: %v", msg, err)
		}
		return nil*/
	default:
	}

	return nil
}
