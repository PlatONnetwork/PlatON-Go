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
	"encoding/json"
	"fmt"
	"math/big"
	"math/rand"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/p2p"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// fakeCbft is a fake cbft for testing.It implements all
// methods of the Cbft interface.
type fakeCbft struct {
	localPeer      *peer             // Represents a local peer
	consensusNodes []discover.NodeID // All consensus nodes
	writer         p2p.MsgReadWriter // Pipeline for receiving data.
	peers          []*peer           // Pre-initialized node for testing.
}

// Returns the ID of the local node.
func (s *fakeCbft) NodeID() discover.NodeID {
	return s.localPeer.Peer.ID()
}

// Return all consensus nodes.
func (s *fakeCbft) ConsensusNodes() ([]discover.NodeID, error) {
	return s.consensusNodes, nil
}

// Return to simulation test configuration.
func (s *fakeCbft) Config() *types.Config {
	return &types.Config{
		Option: &types.OptionsConfig{
			WalMode:           false,
			PeerMsgQueueSize:  1024,
			EvidenceDir:       "evidencedata",
			MaxPingLatency:    5000,
			MaxQueuesLimit:    4096,
			BlacklistDeadline: 5,
		},
		Sys: &params.CbftConfig{
			Period: 3,
		},
	}
}

// ReceiveMessage receives consensus messages.
func (s *fakeCbft) ReceiveMessage(msg *types.MsgInfo) error {
	fmt.Println(fmt.Sprintf("ReceiveMessage, type: %T", msg.Msg))
	return nil
}

// ReceiveSyncMsg receives synchronization messages.
func (s *fakeCbft) ReceiveSyncMsg(msg *types.MsgInfo) error {
	fmt.Println(fmt.Sprintf("ReceiveSyncMsg, type: %T", msg.Msg))
	return nil
}

// Returns the highest local QC height.
func (s *fakeCbft) HighestQCBlockBn() (uint64, common.Hash) {
	return s.localPeer.QCBn(), common.Hash{}
}

// Returns the highest local Lock height.
func (s *fakeCbft) HighestLockBlockBn() (uint64, common.Hash) {
	return s.localPeer.LockedBn(), common.Hash{}
}

// Returns the highest local Commit height.
func (s *fakeCbft) HighestCommitBlockBn() (uint64, common.Hash) {
	return s.localPeer.CommitBn(), common.Hash{}
}

func (s *fakeCbft) MissingViewChangeNodes() (*protocols.GetViewChange, error) {
	return &protocols.GetViewChange{
		Epoch:      1,
		ViewNumber: 1,
	}, nil
}

func (s *fakeCbft) MissingPrepareVote() (*protocols.GetPrepareVote, error) {
	return &protocols.GetPrepareVote{
		Epoch:      1,
		ViewNumber: 1,
		UnKnownSet: utils.NewBitArray(10),
	}, nil
}

func (s *fakeCbft) LatestStatus() *protocols.GetLatestStatus {
	return &protocols.GetLatestStatus{
		BlockNumber:  s.localPeer.QCBn(),
		BlockHash:    common.Hash{},
		QuorumCert:   nil,
		LBlockNumber: s.localPeer.LockedBn(),
		LBlockHash:   common.Hash{},
		LQuorumCert:  nil,
		LogicType:    TypeForQCBn,
	}
}

func (s *fakeCbft) OnPong(nodeID string, netLatency int64) error {
	return nil
}

func (s *fakeCbft) BlockExists(blockNumber uint64, blockHash common.Hash) error {
	return nil
}

// Create a new EngineManager.
func newHandle(t *testing.T) (*EngineManager, *fakeCbft) {
	// init local peer and engineManager.
	var consensusNodes []discover.NodeID
	var peers []*peer
	writer, reader := p2p.MsgPipe()
	var localID discover.NodeID
	rand.Read(localID[:])
	localPeer := newPeer(1, p2p.NewPeer(localID, "local", nil), reader)

	// Simulation generation test node.
	for i := 0; i < testingPeerCount; i++ {
		p, id := newLinkedPeer(writer, 1, fmt.Sprintf("p%d", i))
		// Set the node whose base is indexed as a consensus node.
		if i%2 != 0 {
			consensusNodes = append(consensusNodes, id)
		}
		peers = append(peers, p)
	}
	// define fake cbft.
	fake := &fakeCbft{
		localPeer:      localPeer,
		consensusNodes: consensusNodes,
		writer:         writer,
		peers:          peers,
	}
	engineManager := NewEngineManger(fake)
	return engineManager, fake
}

func Test_EngineManager_Handle(t *testing.T) {
	h, fake := newHandle(t)
	peers := fake.peers
	fakePeer := peers[0]
	// Local and fake need to start sending and reading messages
	// at the same time before starting handshake.
	go func() {
		for {
			msg, _ := fake.localPeer.ReadWriter().ReadMsg()
			t.Logf("localPeer read msg done. type: %d", msg.Code)
			msg.Discard()
		}
	}()
	//
	pingTime := strconv.FormatInt(time.Now().UnixNano(), 10)
	fakePeer.ListPushFront(pingTime)
	testCases := []struct {
		msg     types.Message
		msgType uint64
	}{
		{newFakePrepareBlock(), protocols.PrepareBlockMsg},
		{newFakePrepareVote(), protocols.PrepareVoteMsg},
		{newFakeViewChange(), protocols.ViewChangeMsg},
		{newFakeGetPrepareBlock(), protocols.GetPrepareBlockMsg},
		{newFakeGetBlockQuorumCert(), protocols.GetBlockQuorumCertMsg},
		{newFakeBlockQuorumCert(), protocols.BlockQuorumCertMsg},
		{newFakeGetPrepareVote(), protocols.GetPrepareVoteMsg},
		{newFakePrepareVotes(), protocols.PrepareVotesMsg},
		{newFakeGetQCBlockList(), protocols.GetQCBlockListMsg},
		{newFakeQCBlockList(), protocols.QCBlockListMsg},
		{newFakeGetLatestStatus(), protocols.GetLatestStatusMsg},
		{newFakeLatestStatus(), protocols.LatestStatusMsg},
		{newFakePrepareBlockHash(), protocols.PrepareBlockHashMsg},
		{newFakePing(pingTime), protocols.PingMsg},
		{newFakePong(pingTime), protocols.PongMsg},
		{newFakeCbftStatusData(), protocols.CBFTStatusMsg},
	}
	// First send a status message and then to
	// send consensus messages for processing.
	go func() {
		status := &protocols.CbftStatusData{ProtocolVersion: 1, QCBn: big.NewInt(1), QCBlock: common.Hash{},
			LockBn: big.NewInt(2), LockBlock: common.Hash{}, CmtBn: big.NewInt(3), CmtBlock: common.Hash{}}
		p2p.Send(fake.localPeer.rw, protocols.CBFTStatusMsg, status)
		t.Log("send status success.")
		// send message that the type of consensus.
		for _, v := range testCases {
			p2p.Send(fake.localPeer.rw, v.msgType, v.msg)
		}
	}()
	//
	protocols := h.Protocols()
	protocols[0].NodeInfo()
	pi := protocols[0].PeerInfo(fake.NodeID())
	assert.Nil(t, pi)
	err := protocols[0].Run(fakePeer.Peer, fakePeer.rw)
	//err := h.handler(fakePeer.Peer, fakePeer.rw)
	// Terminate the handle by finally sending an unsupported message type.
	if err != nil && !strings.Contains(err.Error(), "uncontrolled") {
		t.Error("handle failed - ", err)
	}
}

func Test_EngineManager_Forwarding(t *testing.T) {
	handle, fake := newHandle(t)
	peers := fake.peers
	// scenes:
	// 1.need to broadcast, the count of sendQueues is 1.
	// 2.need not to broadcast, the count of sendQueues is 0.
	fakePeer := peers[1]
	handle.peers.Register(fakePeer)
	fakeMessage := newFakeViewChange()
	forward := func(msg types.Message) error {
		return handle.Forwarding("", fakeMessage)
	}
	err := forward(fakeMessage)
	if err != nil {
		t.Error("forwarding failed.", err)
	}

	// mark the msgHash to queues.
	fakePeer.MarkMessageHash(fakeMessage.MsgHash())
	err = forward(fakeMessage)
	assert.NotNil(t, err)
	assert.Equal(t, 1, len(handle.sendQueue))
}

func Test_EngineManager_Send(t *testing.T) {
	handle, _ := newHandle(t)
	handle.Start()

	// Test the messages sent in three different modes to verify
	// whether the sendQueue successfully records the message.
	handle.Send("I", newFakeGetPrepareVote())
	handle.Broadcast(newFakeViewChange())
	handle.PartBroadcast(newFakePrepareBlockHash())
	var wg sync.WaitGroup
	wg.Add(1)
	time.AfterFunc(1*time.Second, func() {
		handle.Close()
		wg.Done()
	})
	wg.Wait()
	assert.Equal(t, 0, len(handle.sendQueue))
}

func Test_EngineManager_Synchronize(t *testing.T) {
	handle, fake := newHandle(t)
	peers := fake.peers

	// Register the simulation node into EngineManage.
	for idx, v := range peers {
		v.SetQcBn(new(big.Int).SetUint64(uint64(idx) + 1))
		v.SetLockedBn(new(big.Int).SetUint64(uint64(idx) + 2))
		v.SetCommitdBn(new(big.Int).SetUint64(uint64(idx) + 3))
		handle.peers.Register(v)
	}

	// Verify that registration is successful.
	checkedPeer := peers[1]
	p, err := handle.getPeer(checkedPeer.id)
	if err != nil {
		t.Error("register peer failed", err)
	}
	assert.Equal(t, checkedPeer.id, p.id)

	// Should return an error if an empty string is passed in.
	_, err = handle.getPeer("")
	assert.NotNil(t, err)

	// blacklist
	p1 := peers[0].PeerID()
	p2 := peers[1].PeerID()
	handle.MarkBlacklist(p1)
	handle.MarkBlacklist(p2)

	assert.True(t, handle.ContainsBlacklist(p1))
	assert.True(t, handle.ContainsBlacklist(p2))

	// The length of ConsensusNodes not equal to 0.
	ds, _ := handle.ConsensusNodes()
	assert.NotEqual(t, 0, len(ds))
	go func() {
		handle.synchronize()
		t.Log("handle done")
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	time.AfterFunc(23*time.Second, func() {
		assert.True(t, handle.ContainsBlacklist(p1))
		assert.True(t, handle.ContainsBlacklist(p2))
		handle.Close()
		t.Log("handle close")
		wg.Done()
	})
	wg.Wait()
}

func Test_EngineManager_NodeInfo(t *testing.T) {
	handle, _ := newHandle(t)

	nodeInfo := handle.NodeInfo()
	jsbyt, err := json.Marshal(nodeInfo)
	if err != nil {
		t.Error("NodeInfo test failed", err)
	}
	t.Log(string(jsbyt))
	assert.Contains(t, string(jsbyt), "}")
}

func Test_EngineManager_LargerPeer(t *testing.T) {
	_, fake := newHandle(t)
	peers := fake.peers
	peers[1].SetQcBn(big.NewInt(10))
	peers[2].SetLockedBn(big.NewInt(11))
	peers[3].SetCommitdBn(big.NewInt(12))

	p, largeBn := largerPeer(TypeForQCBn, peers, 9)
	assert.Equal(t, uint64(10), largeBn)
	assert.Equal(t, p.id, peers[1].id)

	p, largeBn = largerPeer(TypeForLockedBn, peers, 9)
	assert.Equal(t, uint64(11), largeBn)
	assert.Equal(t, p.id, peers[2].id)

	p, largeBn = largerPeer(TypeForCommitBn, peers, 9)
	assert.Equal(t, uint64(12), largeBn)
	assert.Equal(t, p.id, peers[3].id)

	_, largeBn = largerPeer(TypeForCommitBn, nil, 9)
	assert.Equal(t, uint64(0), largeBn)
}
