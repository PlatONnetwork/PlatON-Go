package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"math/big"
	"reflect"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (
	MixMode  = iota // all consensus node
	PartMode        // partial node
	FullMode        // all node
)

const (
	sendQueueSize = 10240
)

type MsgPackage struct {
	peerID string
	msg    Message
	mode   uint64 // forwarding mode.
}

type handler interface {
	Start()
	SendAllConsensusPeer(msg Message)
	Send(peerID discover.NodeID, msg Message)
	SendBroadcast(msg Message)
	SendPartBroadcast(msg Message)
	Protocols() []p2p.Protocol
	PeerSet() *peerSet
	GetPeer(peerID string) (*peer, error)
}

type baseHandler struct {
	cbft      *Cbft
	peers     *peerSet
	sendQueue chan *MsgPackage

	quit chan struct{}
}

func NewHandler(cbft *Cbft) *baseHandler {
	return &baseHandler{
		cbft:      cbft,
		peers:     newPeerSet(),
		sendQueue: make(chan *MsgPackage, sendQueueSize),
		quit:      make(chan struct{}, 0),
	}
}

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

func (h *baseHandler) Start() {
	go h.sendLoop()
	go h.syncHighestStatus()
}

func (h *baseHandler) Close() {
	h.quit <- struct{}{}
	close(h.quit)
}

func (h *baseHandler) sendLoop() {
	for {
		select {
		case m := <-h.sendQueue:
			if m == nil || h.cbft.isLoading() {
				log.Warn("read msg from sendQueue, message is nil or isLoading is true", "isLoading", h.cbft.isLoading())
				return
			}
			log.Debug("send msg to queue", "mode", m.mode, "msgHash", m.msg.MsgHash().TerminalString(), "isLoading", h.cbft.isLoading())
			if len(m.peerID) == 0 {
				h.broadcast(m)
			} else {
				h.sendPeer(m)
			}
		}
	}
}

func (h *baseHandler) broadcast(m *MsgPackage) {
	h.cbft.router.gossip(m)
}

func (h *baseHandler) sendPeer(m *MsgPackage) {
	if peer, err := h.peers.Get(m.peerID); err == nil {
		log.Debug("Send message", "targetPeer", m.peerID, "type", reflect.TypeOf(m.msg), "msgHash", m.msg.MsgHash().TerminalString(), "BHash", m.msg.BHash().TerminalString())

		if err := p2p.Send(peer.rw, MessageType(m.msg), m.msg); err != nil {
			log.Error("Send Peer error")
			h.peers.Unregister(m.peerID)
		}
	}
}

func (h *baseHandler) GetPeer(peerID string) (*peer, error) {
	if peerID == "" {
		return nil, fmt.Errorf("Invalid peer id : %v", peerID)
	}
	return h.peers.Get(peerID)
}

func (h *baseHandler) SendAllConsensusPeer(msg Message) {
	log.Debug("SendAllConsensusPeer Invoke", "hash", msg.MsgHash(), "type", reflect.TypeOf(msg), "BHash", msg.BHash().TerminalString())
	select {
	case h.sendQueue <- &MsgPackage{
		msg:  msg,
		mode: FullMode,
	}:
	default:
	}
}

func (h *baseHandler) Send(peerID discover.NodeID, msg Message) {

	select {
	case h.sendQueue <- &MsgPackage{
		peerID: fmt.Sprintf("%x", peerID.Bytes()[:8]),
		msg:    msg,
	}:
	}
}

func (h *baseHandler) SendBroadcast(msg Message) {
	msgPkg := &MsgPackage{
		msg:  msg,
		mode: MixMode,
	}
	select {
	case h.sendQueue <- msgPkg:
		h.cbft.log.Debug("Send message to broadcast queue", "msgHash", msg.MsgHash().TerminalString(), "BHash", msg.BHash().TerminalString())
	default:
	}
}

func (h *baseHandler) SendPartBroadcast(msg Message) {
	msgPkg := &MsgPackage{
		msg:  msg,
		mode: PartMode,
	}
	select {
	case h.sendQueue <- msgPkg:
		h.cbft.log.Debug("Send message to broadcast queue to partial", "msgHash", msg.MsgHash().TerminalString(), "BHash", msg.BHash().TerminalString())
	default:
	}
}

/*func (h *baseHandler) sendViewChangeVote(id *discover.NodeID, msg *viewChangeVote) error {
	if peer, err := h.peers.Get(fmt.Sprintf("%x", id.Bytes()[:8])); err != nil {
		return err
	} else {
		return p2p.Send(peer.rw, ViewChangeVoteMsg, msg)
	}
}*/

func (h *baseHandler) Protocols() []p2p.Protocol {
	return []p2p.Protocol{
		{
			Name:    "cbft",
			Version: 1,
			Length:  uint64(len(messages)),
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

func (h *baseHandler) PeerSet() *peerSet {
	return h.peers
}

func (h *baseHandler) handler(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	peer := newPeer(p, newMeteredMsgWriter(rw))

	// Execute the CBFT handshake
	var (
		head = h.cbft.blockChain.CurrentHeader()
		hash = head.Hash()
	)
	p.Log().Debug("CBFT peer connected, do handshake", "name", peer.Name())
	confirmedBn := h.cbft.getHighestConfirmed().number
	logicBn := h.cbft.getHighestLogical().number
	if err := peer.Handshake(new(big.Int).SetUint64(confirmedBn), new(big.Int).SetUint64(logicBn), hash); err != nil {
		p.Log().Debug("CBFT handshake failed", "err", err)
		return err
	} else {
		p.Log().Debug("CBFT consensus handshake success", "hash", hash.TerminalString(), "number", head.Number)
	}

	// todo: there is something to be done.

	h.peers.Register(peer)
	defer h.peers.Unregister(peer.id)
	for {
		if err := h.handleMsg(peer); err != nil {
			p.Log().Error("CBFT message handling failed", "err", err)
			return err
		}
	}
}

func (h *baseHandler) handleMsg(p *peer) error {
	msg, err := p.rw.ReadMsg()
	if err != nil {
		p.Log().Error("read peer message error", "err", err)
		return err
	}
	if msg.Size > CbftProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, CbftProtocolMaxMsgSize)
	}
	switch {
	case msg.Code == CBFTStatusMsg:
		return errResp(ErrExtraStatusMsg, "uncontrolled status message")
	case msg.Code == GetPrepareBlockMsg:
		var request getPrepareBlock
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		p.MarkMessageHash((&request).MsgHash())
		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == PrepareBlockMsg:
		// Retrieve and decode the propagated block
		var request prepareBlock
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		p.MarkMessageHash((&request).MsgHash())
		request.Block.ReceivedAt = msg.ReceivedAt
		request.Block.ReceivedFrom = p

		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == PrepareVoteMsg:
		// Retrieve and decode the propagated block
		var request prepareVote
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		if h.cbft.isForwarded(p.ID(), &request) {
			return nil
		}
		p.MarkMessageHash((&request).MsgHash())

		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil

	case msg.Code == ViewChangeMsg:
		var request viewChange
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		if h.cbft.isForwarded(p.ID(), &request) {
			return nil
		}
		p.MarkMessageHash((&request).MsgHash())
		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil

	case msg.Code == ViewChangeVoteMsg:
		var request viewChangeVote
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		if h.cbft.isForwarded(p.ID(), &request) {
			return nil
		}
		p.MarkMessageHash((&request).MsgHash())
		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil

	case msg.Code == ConfirmedPrepareBlockMsg:
		var request confirmedPrepareBlock
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		if h.cbft.isForwarded(p.ID(), &request) {
			return nil
		}
		p.MarkMessageHash((&request).MsgHash())
		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == GetPrepareVoteMsg:
		var request getPrepareVote
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}

		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == PrepareVotesMsg:
		var request prepareVotes
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}

		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == GetHighestPrepareBlockMsg:
		var request getHighestPrepareBlock
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}

		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == HighestPrepareBlockMsg:
		var request highestPrepareBlock
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}

		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == PrepareBlockHashMsg:
		var request prepareBlockHash
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		if h.cbft.isForwarded(p.ID(), &request) {
			return nil
		}
		p.MarkMessageHash((&request).MsgHash())
		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == GetLatestStatusMsg:
		var request getLatestStatus
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == LatestStatusMsg:
		var request latestStatus
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}
		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	default:
	}

	return nil
}

// syncHighestStatus is responsible for HighestPrepareBlock synchronization
func (h *baseHandler) syncHighestStatus() {
	confirmedTicker := time.NewTicker(5 * time.Second)
	logicTicker := time.NewTicker(4 * time.Second)
	for {
		select {
		case <-confirmedTicker.C:
			curHighestNum := h.cbft.getHighestConfirmed().number
			peers := h.PeerSet().ConfirmedHighestBnPeers(new(big.Int).SetUint64(curHighestNum))
			if peers != nil && len(peers) != 0 {
				log.Debug("Sync confirmed highest status", "curHighestNum", curHighestNum, "peers", len(peers))
				largerNum := curHighestNum
				largerIndex := -1
				for index, v := range peers {
					pHighest := v.ConfirmedHighestBn().Uint64()
					if pHighest > largerNum {
						largerNum, largerIndex = pHighest, index
					}
				}
				if largerIndex != -1 {
					largerPeer := peers[largerIndex]
					log.Debug("ConfirmedTicker, send getHighestConfirmedStatus message", "currentHighestBn", curHighestNum, "maxHighestPeer", largerPeer.id, "maxHighestBn", largerNum)
					msg := &getLatestStatus{
						Highest: curHighestNum,
						Type:    HIGHEST_CONFIRMED_BLOCK,
					}
					log.Debug("Send getHighestConfirmedStatus message for confirmed number", "msg", msg.String())
					h.Send(largerPeer.ID(), msg)
				}
			}
		case <-logicTicker.C:
			curLogicHighestNum := h.cbft.getHighestLogical().number
			peers := h.PeerSet().LogicHighestBnPeers(new(big.Int).SetUint64(curLogicHighestNum))
			if peers != nil && len(peers) != 0 {
				log.Debug("Sync logic highest status", "curLogicHighestNum", curLogicHighestNum, "peers", len(peers))
				largerNum := curLogicHighestNum
				largerIndex := -1
				for index, v := range peers {
					pHighest := v.LogicHighestBn().Uint64()
					if pHighest > largerNum {
						largerNum, largerIndex = pHighest, index
					}
				}
				if largerIndex != -1 {
					largerPeer := peers[largerIndex]
					log.Debug("LogicTicker, send getHighestConfirmedStatus message", "currentHighestBn", curLogicHighestNum, "maxLogicHighestPeer", largerPeer.id, "maxLogicHighestBn", largerNum)
					msg := &getLatestStatus{
						Highest: curLogicHighestNum,
						Type:    HIGHEST_LOGIC_BLOCK,
					}
					log.Debug("Send getHighestConfirmedStatus message for logic number", "msg", msg.String())
					h.Send(largerPeer.ID(), msg)
				}
			}
		case <-h.quit:
			log.Warn("Handler quit")
			return
		}
	}
}

type NodeInfo struct {
	Config *params.CbftConfig `json:"config"`
}

func (h *baseHandler) NodeInfo() *NodeInfo {
	cbftConfig := h.cbft.config
	return &NodeInfo{
		Config: cbftConfig,
	}
}
