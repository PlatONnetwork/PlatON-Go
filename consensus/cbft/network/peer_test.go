package network

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

func newPeer(version int, name string) (*peer, discover.NodeID) {
	_, net := p2p.MsgPipe()

	// Generate a random id and create the peer.
	var id discover.NodeID
	rand.Read(id[:])

	// Create a peer that belonging to cbft.
	peer := NewPeer(version, p2p.NewPeer(id, name, nil), net)
	return peer, id
}

func Test_NewPeer(t *testing.T) {
	version := 1
	name := "test"
	p, id := newPeer(version, name)
	if p.version != 1 {
		t.Fatalf("version not equal. expect:{1}, actual:{%d}", p.version)
	}
	if p.Name() != name {
		t.Fatalf("name not equal. expect:{1}, actual:{%d}", p.version)
	}
	assert.Equal(t, id.TerminalString(), p.PeerID())

	// test markMessageHash
	for i := 0; i < maxKnownMessageHash+2; i++ {
		p.MarkMessageHash(common.BytesToHash(common.Uint64ToBytes(uint64(i))))
	}
	if !p.ContainsMessageHash(common.BytesToHash(common.Uint64ToBytes(1))) {
		t.Fatalf("does not contain a specified hash")
	}
	if p.ContainsMessageHash(common.BytesToHash(common.Uint64ToBytes(maxKnownMessageHash + 2))) {
		t.Fatalf("should not contain a specified hash")
	}

	// test SetQcBn/QCBn/SetLockedBn/LockedBN/SetCommitBn/CommitBn
	qcBn := new(big.Int).SetUint64(100)
	p.SetQcBn(qcBn)
	assert.Equal(t, qcBn.Uint64(), p.QCBn())

	lockedBn := new(big.Int).SetUint64(200)
	p.SetLockedBn(lockedBn)
	assert.Equal(t, lockedBn.Uint64(), p.LockedBn())

	commitBn := new(big.Int).SetUint64(300)
	p.SetCommitdBn(commitBn)
	assert.Equal(t, commitBn.Uint64(), p.CommitBn())

	// test PeerInfo
	peerInfo := p.Info()
	json, err := json.Marshal(peerInfo)
	if err != nil {
		t.Error(err)
	}
	t.Log(string(json))
	assert.Contains(t, string(json), "{")

}

func Test_PeerSet_Register(t *testing.T) {
	ps := NewPeerSet()
	p1, _ := newPeer(1, "ps1")
	p2, _ := newPeer(1, "ps2")
	//p3, _ := newPeer(1, "ps3")

	// for the function of Register.
	err := ps.Register(p1)
	if err != nil {
		t.Error("err should not be nil")
	}
	err = ps.Register(p1)
	assert.Equal(t, err.Error(), errAlreadyRegistered.Error())
	ps.Close()
	err = ps.Register(p2)
	assert.Equal(t, err.Error(), errClosed.Error())
}

func Test_PeerSet_Unregister(t *testing.T) {
	// Create new peerSet and do some initialization.
	ps := NewPeerSet()
	p1, _ := newPeer(1, "ps1")
	p2, _ := newPeer(1, "ps2")
	p3, _ := newPeer(1, "ps3")
	ps.Register(p1)
	ps.Register(p2)

	// unregister
	err := ps.Unregister(p1.id)
	if err != nil {
		t.Error("err should not be nil")
	}
	// Try to destroy a peer that does not exist,
	// match the expected error.
	err = ps.Unregister(p3.id)
	assert.Equal(t, err.Error(), errNotRegistered.Error())

	//
	rp, _ := ps.Get(p1.id)
	assert.Equal(t, p1, rp)
}
