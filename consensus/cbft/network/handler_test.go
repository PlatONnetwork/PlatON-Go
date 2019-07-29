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
	localPeer      *peer
	consensusNodes []discover.NodeID
	writer         p2p.MsgReadWriter // Pipeline for receiving data.
	peers          []*peer           // Pre-initialized node for testing.
}

func (s *fakeCbft) NodeId() discover.NodeID {
	return s.localPeer.Peer.ID()
}

func (s *fakeCbft) ConsensusNodes() ([]discover.NodeID, error) {
	return s.consensusNodes, nil
}

func (s *fakeCbft) Config() *types.Config {
	return &types.Config{
		Option: &types.OptionsConfig{
			WalMode:          false,
			PeerMsgQueueSize: 1024,
			EvidenceDir:      "evidencedata",
			MaxPingLatency:   5000,
			MaxQueuesLimit:   4096,
		},
		Sys: &params.CbftConfig{
			Period: 3,
			Epoch:  30000,
		},
	}
}
func (s *fakeCbft) ReceiveMessage(msg *types.MsgInfo) {
	fmt.Println(fmt.Sprintf("ReceiveMessage, type: %T", msg.Msg))
}
func (s *fakeCbft) ReceiveSyncMsg(msg *types.MsgInfo) {
	fmt.Println(fmt.Sprintf("ReceiveSyncMsg, type: %T", msg.Msg))
}
func (s *fakeCbft) HighestQCBlockBn() (uint64, common.Hash) {
	return s.localPeer.QCBn(), common.Hash{}
}
func (s *fakeCbft) HighestLockBlockBn() (uint64, common.Hash) {
	return s.localPeer.LockedBn(), common.Hash{}
}
func (s *fakeCbft) HighestCommitBlockBn() (uint64, common.Hash) {
	return s.localPeer.CommitBn(), common.Hash{}
}

// Create a new EngineManager.
func newHandle(t *testing.T) (*EngineManager, *fakeCbft) {
	// init local peer and engineManager.
	var consensusNodes []discover.NodeID
	var peers []*peer
	writer, reader := p2p.MsgPipe()
	var localId discover.NodeID
	rand.Read(localId[:])
	localPeer := NewPeer(1, p2p.NewPeer(localId, "local", nil), reader)

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
	fakePeer.PingList.PushFront(pingTime)
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
		status := &protocols.CbftStatusData{1, big.NewInt(1), common.Hash{}, big.NewInt(2), common.Hash{}, big.NewInt(3), common.Hash{}}
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
	pi := protocols[0].PeerInfo(fake.NodeId())
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
	// 2. need not to broadcast, the count of sendQueues is 0.
	fakePeer := peers[1]
	handle.peers.Register(fakePeer)
	fakeMessage := &protocols.PrepareBlockHash{
		BlockHash:   common.Hash{},
		BlockNumber: 1,
	}
	forward := func(msg types.Message) error {
		return handle.Forwarding("", fakeMessage)
	}
	err := forward(fakeMessage)
	if err != nil {
		t.Error("forwarding failed.", err)
	}
	assert.Equal(t, 1, len(handle.sendQueue))
	select {
	case <-handle.sendQueue:
	}

	// mark the msgHash to queues.
	fakePeer.MarkMessageHash(fakeMessage.MsgHash())
	err = forward(fakeMessage)
	assert.NotNil(t, err)
	assert.Equal(t, 0, len(handle.sendQueue))
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
	p, err := handle.GetPeer(checkedPeer.id)
	if err != nil {
		t.Error("register peer failed", err)
	}
	assert.Equal(t, checkedPeer.id, p.id)

	// Should return an error if an empty string is passed in.
	p, err = handle.GetPeer("")
	assert.NotNil(t, err)

	// The length of ConsensusNodes not equal to 0.
	ds, err := handle.ConsensusNodes()
	assert.NotEqual(t, 0, len(ds))
	go func() {
		handle.synchronize()
		t.Log("handle done")
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	time.AfterFunc(7*time.Second, func() {
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

	p, largeBn = largerPeer(TypeForCommitBn, nil, 9)
	assert.Equal(t, uint64(0), largeBn)
}
