// Package cbft implements  a concrete consensus engines.
package router

import (
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

// The fanout value of the gossip protocol, used to indicate
// the number of nodes selected per broadcast.
const DEFAULT_FAN_OUT = 5

// Router implements the message protocol of gossip.
//
// 1.Responsible for picking receiving nodes based on message type.
// 2.Generate a random node list based on fan-out.
// 3.Duplicate verification of messages.
type router struct {
	bft    Cbft                          // Implementation of cbft interface
	filter func(*peer, common.Hash) bool // Used for filtering node
	lock   sync.RWMutex
}

// NewRouter creates a new router. It is mainly used for message forwarding
func NewRouter(bft Cbft) *router {
	return &router{
		bft: bft,
		filter: func(p *peer, condition common.Hash) bool {
			return p.knownMessageHash.Contains(condition)
		},
	}
}

// Randomly return a list of consensus nodes that exist in the peerSet.
//
// 1.Get all consensus node IDs.
// 2.Get a list of neighbor nodes.
// 3.Take the intersection of the above.
// 4.If the random is true then get the number of random nodes and the fan-out is 5.
func (r *router) kConsensusRandomNodes(random bool, condition common.Hash) ([]*peer, error) {
	cNodes, err := r.bft.ConsensusNodes()
	if err != nil {
		return nil, err
	}
	//existsPeers := r.msgHandler.PeerSet().Peers()
	existsPeers := make([]*peer, 0)
	log.Debug("kConsensusRandomNodes select node", "msgHash", condition, "cNodesLen", len(cNodes), "peerSetLen", len(existsPeers))

	// The maximum capacity will not exceed the capacity of existsPeers.
	// Slices of the specified capacity have certain performance advantages.
	consensusPeers := make([]*peer, len(existsPeers))
	for _, peer := range existsPeers {
		for _, node := range cNodes {
			if peer.id == node.TerminalString() {
				consensusPeers = append(consensusPeers, peer)
				break
			}
		}
	}
	if random {
		return kRandomNodes(DEFAULT_FAN_OUT, consensusPeers, condition, r.filter), nil
	}
	return consensusPeers, nil
}

// kMixingRandomNodes returns all consensus nodes and k randomly generated non-consensus nodes.
func (r *router) kMixingRandomNodes(condition common.Hash) ([]*peer, error) {
	// all consensus nodes + a number of k non-consensus nodes
	cNodes, err := r.bft.ConsensusNodes()
	if err != nil {
		return nil, err
	}
	//existsPeers := r.msgHandler.PeerSet().Peers()
	existsPeers := make([]*peer, 0)
	consensusPeers := make([]*peer, len(existsPeers))
	// The length of non-consensus nodes is equal to the default fan-out value.
	nonconsensusPeers := make([]*peer, DEFAULT_FAN_OUT)
	for _, peer := range existsPeers {
		isConsensus := false
		for _, node := range cNodes {
			if peer.id == node.TerminalString() {
				isConsensus = true
				break
			}
		}
		if isConsensus {
			consensusPeers = append(consensusPeers, peer)
		} else {
			nonconsensusPeers = append(nonconsensusPeers, peer)
		}
	}
	// Obtain random nodes from non-consensus nodes.
	kNonconsensusNodes := kRandomNodes(DEFAULT_FAN_OUT, nonconsensusPeers, condition, r.filter)
	// Summary target peers and return.
	consensusPeers = append(consensusPeers, kNonconsensusNodes...)
	return consensusPeers, nil
}

// kRandomNodes is used to select up to k random nodes, excluding any nodes where
// the filter function returns true. It is possible that less than k nodes are returned.
func kRandomNodes(k int, peers []*peer, condition common.Hash, filterFn func(*peer, common.Hash) bool) []*peer {
	n := len(peers)
	kNodes := make([]*peer, 0, k)
OUTER:
	// Probe up to 3*n times, with large n this is not necessary
	// since k << n, but with small n we want search to be
	// exhaustive.
	for i := 0; i < 3*n && len(kNodes) < k; i++ {
		// Get random node
		idx := utils.RandomOffset(n)
		node := peers[idx]

		// Give the filter a shot at it.
		if filterFn != nil && filterFn(node, condition) {
			continue OUTER
		}

		// Check if we have this node already
		for j := 0; j < len(kNodes); j++ {
			if node == kNodes[j] {
				continue OUTER
			}
		}
		// Append the node
		kNodes = append(kNodes, node)
	}
	return kNodes
}
