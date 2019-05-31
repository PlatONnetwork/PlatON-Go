package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/stretchr/testify/assert"
	"math/big"
	"math/rand"
	"testing"
	"time"
)

func buildPeer() *peer {
	name := "nodename"
	caps := []p2p.Cap{{"foo", 2}, {"bar", 3}}
	id := randomID()
	p := p2p.NewPeer(id, name, caps)
	peer := newPeer(p, &fakeRW{})
	return peer
}

func TestPeer_NewPeer(t *testing.T) {
	// define peer
	peer := buildPeer()
	assert.NotNil(t, peer)

	// close
	peer.close()
	_, ok := <- peer.term
	assert.True(t, !ok)
}

func TestPeer_MarkMessageHash(t *testing.T) {
	peer := buildPeer()
	for i := 0; i < maxKnownMessageHash + 1; i++ {
		peer.MarkMessageHash(common.BytesToHash([]byte(fmt.Sprintf("I'm hash, %v", i))))
	}
	assert.Equal(t, maxKnownMessageHash, peer.knownMessageHash.Cardinality())
}

func randomID() (id discover.NodeID) {
	for i := range id {
		id[i] = byte(rand.Intn(255))
	}
	return id
}

func TestPeer_Handshake(t *testing.T) {
	var id discover.NodeID
	rand.Read(id[:])
	exec := func(close chan<- struct{}, sNum *big.Int, sHash common.Hash, msg interface{}, rCode uint64) {
		in, out := p2p.MsgPipe()
		p2pPeer := p2p.NewPeer(id, "pid", nil)
		p := newPeer(p2pPeer, in)
		go func() {
			p.Handshake(sNum, sHash)
			close <- struct{}{}
			fmt.Println("handshake close")
		}()
		go func() {
			p2p.Send(out, rCode, msg)
			close <- struct{}{}
			fmt.Println("send close")
		}()
	}
	testCases := []struct{
		sNum *big.Int
		sHash common.Hash
		msg interface{}
		rCode uint64
	}{
		{ big.NewInt(1), common.BytesToHash([]byte("I'm hash")), &cbftStatusData{big.NewInt(1), common.BytesToHash([]byte("I'm hash"))}, CBFTStatusMsg },
		{ big.NewInt(1), common.BytesToHash([]byte("I'm hash")), &cbftStatusData{big.NewInt(1), common.BytesToHash([]byte("I'm hash"))}, PrepareBlockMsg },
		{ big.NewInt(1), common.BytesToHash([]byte("I'm hash")), &prepareBlockHash{}, CBFTStatusMsg },
	}
	for _, v := range testCases {
		close := make(chan struct{}, 2)
		exec(close, v.sNum, v.sHash, v.msg, v.rCode)
		if len(close) != 2 {
			time.Sleep(1 * time.Second)
		}
	}
}

func TestPeerSet_All(t *testing.T) {
	id := randomID()
	p := p2p.NewPeer(id, "test", nil)
	peer := newPeer(p, &fakeRW{})
	ps := newPeerSet()
	ps.Register(peer)

	gotPeer, err := ps.Get(id.TerminalString())
	assert.Nil(t, err)
	assert.Equal(t, peer, gotPeer)

	gotPeer, err = ps.Get("invalid id")
	assert.IsType(t, errNotRegistered, err)

	peers := ps.AllConsensusPeer()
	assert.True(t, len(peers) == 1)

	peers = ps.Peers()
	assert.True(t, len(peers) == 1)

	err = ps.Unregister(id.TerminalString())
	assert.Nil(t, err)
	gotPeer, err = ps.Get(id.TerminalString())
	assert.IsType(t, errNotRegistered, err)

	ps.Close()
	assert.True(t, ps.closed)
}

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
