package router

import "github.com/PlatONnetwork/PlatON-Go/p2p/discover"

type Cbft interface {

	// Returns the ID value of the current node
	NodeId() discover.NodeID

	// Return a list of all consensus nodes
	ConsensusNodes() ([]discover.NodeID, error)
}

type Handler interface {

	// Return all neighbor node lists.
	Peers() ([]*peer, error)
}
