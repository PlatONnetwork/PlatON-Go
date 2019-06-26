package cbft

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/deckarep/golang-set"
	"math/big"
	"sync"
	"time"
)

var (
	//errClosed            = errors.New("peer set is closed")
	//errAlreadyRegistered = errors.New("peer is already registered")
	errNotRegistered = errors.New("peer is not registered")
)

const (
	maxKnownMessageHash = 20000

	handshakeTimeout = 5 * time.Second
)

type PeerInfo struct {
	ConfirmedHigBN *big.Int `json:"confirmed_hig_bn"`
	LogicHigBn     *big.Int `json:"logic_hig_bn"`
	MsgCount       uint64   `json:"msg_count"`
}

type peer struct {
	id   string
	term chan struct{} // Termination channel to stop the broadcaster
	*p2p.Peer
	rw p2p.MsgReadWriter

	confirmedHigBn *big.Int
	cLock          sync.RWMutex

	logicHigBn *big.Int
	lLock      sync.RWMutex

	knownMessageHash mapset.Set
}

func newPeer(p *p2p.Peer, rw p2p.MsgReadWriter) *peer {
	return &peer{
		Peer:             p,
		rw:               rw,
		id:               fmt.Sprintf("%x", p.ID().Bytes()[:8]),
		term:             make(chan struct{}),
		knownMessageHash: mapset.NewSet(),
		confirmedHigBn:   new(big.Int),
		logicHigBn:       new(big.Int),
	}
}

func (p *peer) close() {
	close(p.term)
}

func (p *peer) MarkMessageHash(hash common.Hash) {
	for p.knownMessageHash.Cardinality() >= maxKnownMessageHash {
		p.knownMessageHash.Pop()
	}
	p.knownMessageHash.Add(hash)
}

// exchange node information with each other.
func (p *peer) Handshake(confirmedBn *big.Int, logicBn *big.Int, head common.Hash) error {
	errc := make(chan error, 2)
	var status cbftStatusData

	go func() {
		errc <- p2p.Send(p.rw, CBFTStatusMsg, &cbftStatusData{
			ConfirmedBn:  confirmedBn,
			LogicBn:      logicBn,
			CurrentBlock: head,
		})
	}()
	go func() {
		errc <- p.readStatus(&status)
		if status.ConfirmedBn != nil {
			p.Log().Debug("Receive the cbftStatusData message", "blockHash", status.CurrentBlock.TerminalString(), "ConfirmedBn", status.ConfirmedBn.Int64(), "LogicBn", status.LogicBn.Int64())
			p.SetConfirmedHighestBn(status.ConfirmedBn)
			p.SetLogicHighestBn(status.LogicBn)
		}
	}()
	timeout := time.NewTicker(handshakeTimeout)
	defer timeout.Stop()
	for i := 0; i < 2; i++ {
		select {
		case err := <-errc:
			if err != nil {
				return err
			}
		case <-timeout.C:
			return p2p.DiscReadTimeout
		}
	}
	// todo: Maybe there is something to be done.
	return nil
}

// SetHighest updates the highest number of the peer.
func (p *peer) SetConfirmedHighestBn(highestBn *big.Int) {
	if highestBn != nil {
		p.cLock.Lock()
		defer p.cLock.Unlock()
		log.Debug("Set confirmed highest number", "peerID", p.id, "oldConfirmedHighest", p.confirmedHigBn.Uint64(), "newConfirmedHighest", highestBn.Uint64())
		p.confirmedHigBn.Set(highestBn)
	}
}

func (p *peer) ConfirmedHighestBn() (bn *big.Int) {
	p.cLock.RLock()
	defer p.cLock.RUnlock()
	return new(big.Int).Set(p.confirmedHigBn)
}

func (p *peer) SetLogicHighestBn(highestBn *big.Int) {
	if highestBn != nil {
		p.cLock.Lock()
		defer p.cLock.Unlock()
		log.Debug("Set logic highest number", "peerID", p.id, "oldLogicHighest", p.logicHigBn.Uint64(), "newLogicHighest", highestBn.Uint64())
		p.logicHigBn.Set(highestBn)
	}
}

func (p *peer) LogicHighestBn() (bn *big.Int) {
	p.lLock.RLock()
	defer p.lLock.RUnlock()
	return new(big.Int).Set(p.logicHigBn)
}

func (p *peer) readStatus(status *cbftStatusData) error {
	msg, err := p.rw.ReadMsg()

	if err != nil {
		return err
	}
	if msg.Code != CBFTStatusMsg {
		return errResp(ErrNoStatusMsg, "first msg has code %x (!= %x)", msg.Code, CBFTStatusMsg)
	}
	if msg.Size > CbftProtocolMaxMsgSize {
		return errResp(ErrMsgTooLarge, "%v > %v", msg.Size, CbftProtocolMaxMsgSize)
	}
	if err := msg.Decode(&status); err != nil {

		return errResp(ErrDecode, "msg %v: %v", msg, err)
	}
	// todo: additional judgment.
	return nil
}

func (p *peer) Info() *PeerInfo {
	c, l, size := p.confirmedHigBn, p.logicHigBn, p.knownMessageHash.Cardinality()

	return &PeerInfo{
		ConfirmedHigBN: c,
		LogicHigBn:     l,
		MsgCount:       uint64(size),
	}
}

type peerSet struct {
	peers  map[string]*peer
	lock   sync.RWMutex
	closed bool
}

func newPeerSet() *peerSet {
	// Monitor output node list
	ps := &peerSet{
		peers: make(map[string]*peer),
	}
	go ps.printPeers()
	return ps
}

func (ps *peerSet) Register(p *peer) {
	ps.lock.Lock()
	defer ps.lock.Unlock()
	ps.peers[p.id] = p
}

func (ps *peerSet) Unregister(id string) error {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	p, ok := ps.peers[id]
	if !ok {
		return errNotRegistered
	}
	delete(ps.peers, id)
	p.close()

	return nil
}

func (ps *peerSet) Get(id string) (*peer, error) {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	p, ok := ps.peers[id]
	if !ok {
		return nil, errNotRegistered
	}

	return p, nil
}

func (ps *peerSet) AllConsensusPeer() []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		list = append(list, p)
	}
	return list
}

func (ps *peerSet) Close() {
	ps.lock.Lock()
	defer ps.lock.Unlock()

	for _, p := range ps.peers {
		p.Disconnect(p2p.DiscQuitting)
	}
	ps.closed = true
}

// Return all peer.
func (ps *peerSet) Peers() []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()

	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		list = append(list, p)
	}
	return list
}

// Returns a list of nodes that are larger than the height of the highest confirmed block.
func (ps *peerSet) ConfirmedHighestBnPeers(highest *big.Int) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		//log.Debug("ConfirmedHighestBnPeers", "pHighest", p.ConfirmedHighestBn().Uint64(), "highest", highest.Uint64(), "peer", p.id)
		if p.ConfirmedHighestBn().Cmp(highest) > 0 {
			list = append(list, p)
		}
	}
	log.Debug("Get confirmed highest peers", "count", len(list), "peers", formatPeers(list))
	return list
}

func (ps *peerSet) LogicHighestBnPeers(highest *big.Int) []*peer {
	ps.lock.RLock()
	defer ps.lock.RUnlock()
	list := make([]*peer, 0, len(ps.peers))
	for _, p := range ps.peers {
		//log.Debug("LogicHighestBnPeers", "pHighest", p.LogicHighestBn().Uint64(), "highest", highest.Uint64(), "peer", p.id)
		if p.LogicHighestBn().Cmp(highest) > 0 {
			list = append(list, p)
		}
	}
	log.Debug("Get logic highest peers", "count", len(list), "peers", formatPeers(list))
	return list
}

func (ps *peerSet) printPeers() {
	// Output in 2 seconds
	outTimer := time.NewTicker(time.Second * 5)
	for {
		if ps.closed {
			break
		}
		select {
		case <-outTimer.C:
			peers := ps.Peers()
			var bf bytes.Buffer
			for idx, peer := range peers {
				bf.WriteString(peer.id)

				if idx < len(peers)-1 {
					bf.WriteString(",")
				}
			}
			pInfo := bf.String()
			log.Info(fmt.Sprintf("The neighbor node owned by the current peer is : {%v}, size: {%d}", pInfo, len(peers)))
		}
	}
}
