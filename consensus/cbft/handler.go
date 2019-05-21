package cbft

import (
	"fmt"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (
	sendQueueSize = 10240
)

type MsgPackage struct {
	peerID string
	msg    Message
}

type handler struct {
	cbft      *Cbft
	peers     *peerSet
	sendQueue chan *MsgPackage
}

func NewHandler(cbft *Cbft) *handler {
	return &handler{
		cbft:      cbft,
		peers:     newPeerSet(),
		sendQueue: make(chan *MsgPackage, sendQueueSize),
	}
}

func errResp(code errCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

func (h *handler) Start() {
	go h.sendLoop()
}

func (h *handler) sendLoop() {
	for {
		select {
		case m := <-h.sendQueue:
			if len(m.peerID) == 0 {
				h.broadcast(m)
			} else {
				h.sendPeer(m)
			}
		}
	}
}

func (h *handler) broadcast(m *MsgPackage) {
	h.cbft.router.gossip(m)
}

func (h *handler) sendPeer(m *MsgPackage) {
	if peer, err := h.peers.Get(m.peerID); err == nil {
		log.Debug("Send message", "targetPeer", m.peerID, "type", reflect.TypeOf(m.msg), "msgHash", m.msg.MsgHash().TerminalString(), "BHash", m.msg.BHash().TerminalString())

		if err := p2p.Send(peer.rw, MessageType(m.msg), m.msg); err != nil {
			log.Error("Send Peer error")
			h.peers.Unregister(m.peerID)
		}
	}
}

func (h *handler) SendAllConsensusPeer(msg Message) {
	log.Debug("SendAllConsensusPeer Invoke", "hash", msg.MsgHash(), "type", reflect.TypeOf(msg), "BHash", msg.BHash().TerminalString())
	select {
	case h.sendQueue <- &MsgPackage{
		msg: msg,
	}:
	default:
	}
}

func (h *handler) Send(peerID discover.NodeID, msg Message) {
	select {
	case h.sendQueue <- &MsgPackage{
		peerID: fmt.Sprintf("%x", peerID.Bytes()[:8]),
		msg:    msg,
	}:
	default:
	}
}

func (h *handler) SendBroadcast(msg Message) {
	msgPkg := &MsgPackage{
		msg: msg,
	}
	select {
	case h.sendQueue <- msgPkg:
		h.cbft.log.Debug("Send message to broadcast queue", "msgHash", msg.MsgHash().TerminalString(), "BHash", msg.BHash().TerminalString())
	default:
	}
}

func (h *handler) sendViewChangeVote(id *discover.NodeID, msg *viewChangeVote) error {
	if peer, err := h.peers.Get(fmt.Sprintf("%x", id.Bytes()[:8])); err != nil {
		return err
	} else {
		return p2p.Send(peer.rw, ViewChangeVoteMsg, msg)
	}
}

func (h *handler) Protocols() []p2p.Protocol {
	return []p2p.Protocol{
		{
			Name:    "cbft",
			Version: 1,
			Length:  uint64(len(messages)),
			Run: func(p *p2p.Peer, rw p2p.MsgReadWriter) error {
				return h.handler(p, rw)
			},
		},
	}
}

func (h *handler) handler(p *p2p.Peer, rw p2p.MsgReadWriter) error {
	peer := newPeer(p, newMeteredMsgWriter(rw))

	// Execute the CBFT handshake
	var (
		head = h.cbft.blockChain.CurrentHeader()
		hash = head.Hash()
	)
	p.Log().Debug("CBFT peer connected, do handshake", "name", peer.Name())
	if err := peer.Handshake(head.Number, hash); err != nil {
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

func (h *handler) handleMsg(p *peer) error {
	msg, err := p.rw.ReadMsg()
	if err != nil {
		p.Log().Error("read peer message error", "err", err)
		return err
	}

	switch {
	case msg.Code == CBFTStatusMsg:
		return errResp(ErrExtraStatusMsg, "uncontrolled status message")
	case msg.Code == GetPrepareBlockMsg:
		var request getPrepareBlock
		if err := msg.Decode(&request); err != nil {
			return errResp(ErrDecode, "%v: %v", msg, err)
		}

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
		log.Info("Received a PrepareBlockMsg", "peer", p.id, "prepare", request.String())

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
		//p.MarkMessageHash((&request).MsgHash())

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
		//p.MarkMessageHash((&request).MsgHash())

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
		p.MarkMessageHash((&request).MsgHash())

		h.cbft.ReceivePeerMsg(&MsgInfo{
			Msg:    &request,
			PeerID: p.ID(),
		})
		return nil
	case msg.Code == GetPrepareBlockMsg:
		var request getPrepareBlock
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
