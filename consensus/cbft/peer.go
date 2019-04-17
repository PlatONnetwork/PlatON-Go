package cbft

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"sync"
)

var (
	errClosed            = errors.New("peer set is closed")
	errAlreadyRegistered = errors.New("peer is already registered")
	errNotRegistered     = errors.New("peer is not registered")
)

type peer struct {
	id   string
	term chan struct{} // Termination channel to stop the broadcaster
	*p2p.Peer
	rw p2p.MsgReadWriter
}

func newPeer(p *p2p.Peer, rw p2p.MsgReadWriter) *peer {
	return &peer{
		Peer: p,
		rw:   rw,
		id:   fmt.Sprintf("%x", p.ID().Bytes()[:8]),
		term: make(chan struct{}),
	}
}

func (p *peer) close() {
	close(p.term)
}

type peerSet struct {
	peers  map[string]*peer
	lock   sync.RWMutex
	closed bool
}

func newPeerSet() *peerSet {
	return &peerSet{
		peers: make(map[string]*peer),
	}
}

func (ps *peerSet) Register(p *peer) {
	ps.lock.Lock()
	ps.lock.Unlock()
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
	ps.lock.Lock()
	defer ps.lock.Unlock()

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
