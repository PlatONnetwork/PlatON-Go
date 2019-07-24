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
