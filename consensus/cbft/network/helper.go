package network

import (
	"math/rand"

	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// ============================ simulation network ============================
func randomID() []discover.NodeID {
	ids := make([]discover.NodeID, 0)
	for i := 0; i < 4; i++ {
		var id discover.NodeID
		rand.Read(id[:])
		ids = append(ids, id)
	}
	return ids
}

// EnhanceEngineManager is used to register a batch of handlers to
// simulate the test environment.
//
// The number of simulated network nodes is fixed at four.
func EnhanceEngineManager(ids []discover.NodeID, handlers []*EngineManager) {

	// node 1 => 1 <--> 2 association.
	rw1_node_1_2, rw2_node_1_2 := p2p.MsgPipe()
	rw1_node_1_3, rw2_node_1_3 := p2p.MsgPipe()
	rw1_node_1_4, rw2_node_1_4 := p2p.MsgPipe()
	rw1_node_3_2, rw2_node_3_2 := p2p.MsgPipe()
	rw1_node_3_4, rw2_node_3_4 := p2p.MsgPipe()
	rw1_node_2_4, rw2_node_2_4 := p2p.MsgPipe()

	peer1_1_2 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[0], ids[0].TerminalString(), nil), rw1_node_1_2)
	peer2_1_2 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[1], ids[1].TerminalString(), nil), rw2_node_1_2)

	// 1 <--> 3 association.
	peer1_1_3 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[0], ids[0].TerminalString(), nil), rw1_node_1_3)
	peer3_1_3 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[2], ids[2].TerminalString(), nil), rw2_node_1_3)

	// 1 <--> 4 association.
	peer1_1_4 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[0], ids[0].TerminalString(), nil), rw1_node_1_4)
	peer4_1_4 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[3], ids[3].TerminalString(), nil), rw2_node_1_4)

	peer3_3_2 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[2], ids[2].TerminalString(), nil), rw1_node_3_2)
	peer2_3_2 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[1], ids[1].TerminalString(), nil), rw2_node_3_2)

	peer3_3_4 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[2], ids[2].TerminalString(), nil), rw1_node_3_4)
	peer4_3_4 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[3], ids[3].TerminalString(), nil), rw2_node_3_4)

	peer2_2_4 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[1], ids[1].TerminalString(), nil), rw1_node_2_4)
	peer4_2_4 := NewPeer(CbftProtocolVersion, p2p.NewPeer(ids[3], ids[3].TerminalString(), nil), rw2_node_2_4)

	// register for h1
	handlers[0].peers.Register(peer2_1_2)
	handlers[0].peers.Register(peer3_1_3)
	handlers[0].peers.Register(peer4_1_4)

	// register for h2
	handlers[1].peers.Register(peer1_1_2)
	handlers[1].peers.Register(peer3_3_2)
	handlers[1].peers.Register(peer4_2_4)

	// register for h3
	handlers[2].peers.Register(peer1_1_3)
	handlers[2].peers.Register(peer2_3_2)
	handlers[2].peers.Register(peer4_3_4)

	// register for h4
	handlers[3].peers.Register(peer1_1_4)
	handlers[3].peers.Register(peer3_3_4)
	handlers[3].peers.Register(peer2_2_4)
}
