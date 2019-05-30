package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestPeerSet_all(t *testing.T) {

	peerSet := newPeerSet()

	var nodeIds []discover.NodeID

	for i := 0; i < 10; i++ {
		node, _ := discover.BytesID(Rand32Bytes(64))
		nodeIds = append(nodeIds, node)
	}
	for _, n := range nodeIds {
		_, mr, peer, _ := p2p.NewMockPeerNodeID(n, nil)

		peerSet.Register(newPeer(peer, mr))
	}
	assert.Equal(t, len(nodeIds), len(peerSet.Peers()))
	assert.Len(t, peerSet.AllConsensusPeer(), len(nodeIds))
	for _, n := range nodeIds {
		id := fmt.Sprintf("%x", n.Bytes()[:8])
		p, err := peerSet.Get(id)
		assert.Nil(t, err)
		assert.Equal(t, p.id, id)
		assert.Nil(t, peerSet.Unregister(id))

	}

	assert.Empty(t, peerSet.AllConsensusPeer())

	time.Sleep(6 * time.Second)

	peerSet.Close()
}

func TestPeer(t *testing.T) {
	node, _ := discover.BytesID(Rand32Bytes(64))
	_, mr, p, _ := p2p.NewMockPeerNodeID(node, nil)
	peer := newPeer(p, mr)
	for i := 0; i < 10; i++ {
		peer.MarkMessageHash(common.BytesToHash(Rand32Bytes(32)))
	}
}
