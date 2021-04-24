// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package network

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"sync"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

func Test_NewPeer(t *testing.T) {
	version := 1
	name := "test"
	p, id := newTestPeer(version, name)
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
	p1, _ := newTestPeer(1, "ps1")
	p2, _ := newTestPeer(1, "ps2")
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
	p1, _ := newTestPeer(1, "ps1")
	p2, _ := newTestPeer(1, "ps2")
	p3, _ := newTestPeer(1, "ps3")
	ps.Register(p1)
	ps.Register(p2)

	// Verify the number of successful registrations.
	len := ps.Len()
	assert.Equal(t, 2, len)

	rp, err := ps.get(p1.id)
	if err != nil {
		t.Error("Get peer should not be return nil")
	}
	assert.Equal(t, p1.id, rp.id)

	// unregister
	err = ps.Unregister(p1.id)
	if err != nil {
		t.Error("err should not be nil")
	}
	// Try to destroy a peer that does not exist,
	// match the expected error.
	err = ps.Unregister(p3.id)
	assert.Equal(t, err.Error(), errNotRegistered.Error())

	_, err = ps.get(p3.id)
	assert.Equal(t, err.Error(), errNotRegistered.Error())
}

func Test_PeerSet_Peers(t *testing.T) {
	// Randomly generate a batch of nodes.
	// Init the node of consensus.
	// Init the node of peerSet.
	ps := NewPeerSet()
	var peers []*peer
	var ids []discover.NodeID
	for i := 0; i < 11; i++ {
		p, id := newTestPeer(1, fmt.Sprintf("%d", i))
		peers = append(peers, p)
		// The id is oddly set to the consensus node.
		if i%2 != 0 {
			ids = append(ids, id)
			p.SetQcBn(new(big.Int).SetUint64(uint64(i) * 100))
			p.SetLockedBn(new(big.Int).SetUint64(uint64(i) * 100))
			p.SetCommitdBn(new(big.Int).SetUint64(uint64(i) * 100))
		}
		ps.Register(p)
	}

	// test PeersWithoutConsensus.
	pwoc := ps.peersWithoutConsensus(ids)
	assert.Equal(t, len(peers)-len(ids), len(pwoc))

	// test peersWithConsensus
	pwc := ps.peersWithConsensus(ids)
	assert.Equal(t, len(ids), len(pwc))

	// test peers
	pees := ps.allPeers()
	assert.Equal(t, len(peers), len(pees))

	// test PeersWithHighestQCBn, i(1/3/5/7/9) * 100 (i is an odd number).
	// If qcNumber is 700, then the count of results is 1.
	pwhqb := ps.peersWithHighestQCBn(700)
	assert.Equal(t, 1, len(pwhqb))

	// If lockedNumber is 700, then the count of results is 1.
	pwhlb := ps.peersWithHighestLockedBn(500)
	assert.Equal(t, 2, len(pwhlb))

	// If lockedNumber is 700, then the count of results is 1.
	pwhmb := ps.peersWithHighestCommitBn(300)
	assert.Equal(t, 3, len(pwhmb))

	// Print node information.
	go ps.printPeers()
	var wg sync.WaitGroup
	wg.Add(1)
	time.AfterFunc(time.Second*2, func() {
		ps.Close()
		wg.Done()
	})
	wg.Wait()
}

func Test_Peer_Handshake(t *testing.T) {
	exec := func(close chan<- struct{}, inStatus, outStatus *protocols.CbftStatusData, wantErr error) {
		in, out := p2p.MsgPipe()
		var id discover.NodeID
		rand.Read(id[:])
		me := newPeer(1, p2p.NewPeer(id, "me", nil), in)
		you := newPeer(1, p2p.NewPeer(id, "you", nil), out)
		go func() {
			_, err := me.Handshake(inStatus)
			if err != nil && wantErr != nil {
				assert.Contains(t, err.Error(), wantErr.Error())
				t.Log(err.Error())
				t.Log(wantErr.Error())
			} else {
				assert.Equal(t, outStatus.QCBn.Uint64(), me.QCBn())
				assert.Equal(t, outStatus.LockBn.Uint64(), me.LockedBn())
				assert.Equal(t, outStatus.CmtBn.Uint64(), me.CommitBn())
			}
			close <- struct{}{}
			t.Log("handshake done to me")
		}()
		go func() {
			_, err := you.Handshake(outStatus)
			if err != nil && wantErr != nil {
				t.Log(err.Error())
				t.Log(wantErr.Error())
				assert.Contains(t, err.Error(), wantErr.Error())
			} else {
				assert.Equal(t, inStatus.QCBn.Uint64(), you.QCBn())
				assert.Equal(t, inStatus.LockBn.Uint64(), you.LockedBn())
				assert.Equal(t, inStatus.CmtBn.Uint64(), you.CommitBn())
			}
			close <- struct{}{}
			t.Log("handshake done to you")
		}()
	}
	// test suite
	testCase := []struct {
		in      *protocols.CbftStatusData
		out     *protocols.CbftStatusData
		wantErr error
	}{
		{
			in:      &protocols.CbftStatusData{ProtocolVersion: 1, QCBn: big.NewInt(1), QCBlock: common.Hash{}, LockBn: big.NewInt(2), LockBlock: common.Hash{}, CmtBn: big.NewInt(3), CmtBlock: common.Hash{}},
			out:     &protocols.CbftStatusData{ProtocolVersion: 1, QCBn: big.NewInt(2), QCBlock: common.Hash{}, LockBn: big.NewInt(3), LockBlock: common.Hash{}, CmtBn: big.NewInt(4), CmtBlock: common.Hash{}},
			wantErr: nil,
		},
		{
			in:      &protocols.CbftStatusData{ProtocolVersion: 1, QCBn: big.NewInt(1), QCBlock: common.Hash{}, LockBn: big.NewInt(2), LockBlock: common.Hash{}, CmtBn: big.NewInt(3), CmtBlock: common.Hash{}},
			out:     &protocols.CbftStatusData{ProtocolVersion: 2, QCBn: big.NewInt(9), QCBlock: common.Hash{}, LockBn: big.NewInt(8), LockBlock: common.Hash{}, CmtBn: big.NewInt(7), CmtBlock: common.Hash{}},
			wantErr: types.ErrResp(types.ErrCbftProtocolVersionMismatch, "%s", ""),
		},
	}
	for _, v := range testCase {
		close := make(chan struct{}, 2)
		exec(close, v.in, v.out, v.wantErr)
		timeout := time.NewTicker(handshakeTimeout)
		defer timeout.Stop()
		for i := 0; i < 2; i++ {
			select {
			case <-close:
			case <-timeout.C:
				t.Error("handshake test timeout")
			}
		}
	}
}
