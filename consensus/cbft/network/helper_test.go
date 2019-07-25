package network

import (
	"crypto/rand"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/p2p"
)

// fakePeer is a simulated peer to allow testing direct network calls.
type fakePeer struct {
	net   p2p.MsgReadWriter // Network layer reader/writer to simulate remote messaging.
	app   *p2p.MsgPipeRW    // Application layer reader/writer to simulate the local side.
	*peer                   // The peer belonging to CBFT layer.
}

// newFakePeer creates a new peer registered at the given protocol manager.
func newFakePeer(name string, version int, pm *EngineManager, shake bool) (*fakePeer, <-chan error) {
	// Create a message pipe to communicate through.
	app, net := p2p.MsgPipe()

	// Generate a random id and create the peer.
	var id discover.NodeID
	rand.Read(id[:])

	// Create a peer that belonging to cbft.
	peer := NewPeer(version, p2p.NewPeer(id, name, nil), net)

	// Start the peer on a new thread
	errc := make(chan error, 1)
	go func() {
		//
		errc <- pm.handler(peer.Peer, peer.rw)
	}()
	tp := &fakePeer{app: app, net: net, peer: peer}
	return tp, errc
}

// Create a new peer for testing, return peer and ID.
func newPeer(version int, name string) (*peer, discover.NodeID) {
	_, net := p2p.MsgPipe()

	// Generate a random id and create the peer.
	var id discover.NodeID
	rand.Read(id[:])

	// Create a peer that belonging to cbft.
	peer := NewPeer(version, p2p.NewPeer(id, name, nil), net)
	return peer, id
}

func newLinkedPeer(rw p2p.MsgReadWriter, version int, name string) (*peer, discover.NodeID) {
	// Generate a random id and create the peer.
	var id discover.NodeID
	rand.Read(id[:])

	// Create a peer that belonging to cbft.
	peer := NewPeer(version, p2p.NewPeer(id, name, nil), rw)
	return peer, id
}
