package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"math/big"
	"os"
	"reflect"
	"testing"
	"time"
)

var TestPeerSet = &peerSet{
	peers: map[string]*peer{
		"test0": &peer{id: "test0"},
		"test1": &peer{id: "test1"},
		"test2": &peer{id: "test2"},
		"test3": &peer{id: "test3"},
		"test4": &peer{id: "test4"},
	},
}

func newTestRouter() *router {
	path := path()
	defer os.RemoveAll(path)
	engine, _, _ := randomCBFT(path, 1)
	peerId := engine.getValidators().NodeList()[0].TerminalString()
	handler := makeHandler(engine, peerId, common.Hash{})
	return NewRouter(engine, handler)
}

func TestGossip(t *testing.T) {
	router := newTestRouter()
	testCases := []struct {
		mode uint64
		msg  Message
	}{
		{mode: FullMode, msg: &prepareBlockHash{}},
		{mode: PartMode, msg: &prepareBlockHash{}},
		{mode: MixMode, msg: &prepareBlockHash{}},
		{mode: FullMode, msg: makeFakePrepareBlock()},
		{mode: FullMode, msg: makeFakePrepareVote()},
		{mode: FullMode, msg: &confirmedPrepareBlock{}},
		{mode: FullMode, msg: makeFakeViewChange()},
		{mode: FullMode, msg: makeFakeGetPrepareBlock()},
		{mode: FullMode, msg: makeFakeGetHighestPrepareBlock()},
		{mode: FullMode, msg: &cbftStatusData{}},
	}
	for _, v := range testCases {
		router.gossip(&MsgPackage{peerID: "peerid", mode: v.mode, msg: v.msg})
	}
}

func makeFakePrepareBlock() *prepareBlock {
	block := types.NewBlockWithHeader(&types.Header{
		GasLimit:  uint64(3141592),
		GasUsed:   uint64(21000),
		Coinbase:  common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"),
		Root:      common.HexToHash("ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017"),
		//Hash: common.HexToHash("0a5843ac1cb04865017cb35a57b50b07084e5fcee39b5acadade33149f4fff9e"),
		Nonce: types.EncodeNonce(RandBytes(81)),
		Time:  big.NewInt(1426516743),
		Extra: make([]byte, 100),
	})
	pb := &prepareBlock{
		Timestamp:     uint64(time.Now().Unix()),
		Block:         block,
		ProposalIndex: 1,
		ProposalAddr:  common.BytesToAddress([]byte("I'm address")),
		View:          &viewChange{},
	}
	return pb
}

func makeFakeGetPrepareBlock() *getPrepareBlock {
	return &getPrepareBlock{
		Hash:   common.BytesToHash([]byte("Empty block")),
		Number: 1,
	}
}

func makeFakeGetHighestPrepareBlock() *getHighestPrepareBlock {
	return &getHighestPrepareBlock{
		Lowest: 1,
	}
}

func makeFakePrepareVote() *prepareVote {
	pv := &prepareVote{
		Timestamp:      uint64(time.Now().Unix()),
		Hash:           common.BytesToHash([]byte("I'm hash")),
		Number:         1,
		ValidatorIndex: 0,
		ValidatorAddr:  common.BytesToAddress([]byte("I'm address")),
	}
	return pv
}

func makeFakeConfirmedPrepareBlock() *confirmedPrepareBlock {
	pv := &confirmedPrepareBlock{
		Hash:     common.BytesToHash([]byte("I'm hash")),
		Number:   1,
		VoteBits: NewBitArray(12),
	}
	return pv
}

func makeFakeGetPrepareVote() *getPrepareVote {
	pv := &getPrepareVote{
		Hash:     common.BytesToHash([]byte("I'm hash")),
		Number:   1,
		VoteBits: NewBitArray(12),
	}
	return pv
}

func makeFakeViewChange() *viewChange {
	privateHex := "e4eb3e58ab7810984a0c77d432b07fe9f9897158dd4bb4f63d0a4366e6d949fa"
	pri, _ := crypto.HexToECDSA(privateHex)
	pv := &viewChange{
		Timestamp:     uint64(time.Now().Unix()),
		ProposalIndex: 0,
		ProposalAddr:  common.BytesToAddress([]byte("I'm address")),
		BaseBlockHash: common.BytesToHash([]byte("I'm hash")),
		BaseBlockNum:  1,
		Extra:         make([]byte, 100),
	}
	var consensusMsg ConsensusMsg = pv
	cb, _ := consensusMsg.CannibalizeBytes()
	sign, _ := crypto.Sign(cb, pri)
	pv.Signature.SetBytes(sign)
	return pv
}

func makeFakeViewChangeVote() *viewChangeVote {
	privateHex := "e4eb3e58ab7810984a0c77d432b07fe9f9897158dd4bb4f63d0a4366e6d949fa"
	pri, _ := crypto.HexToECDSA(privateHex)
	pv := &viewChangeVote{
		Timestamp:     uint64(time.Now().Unix()),
		ProposalIndex: 0,
		ProposalAddr:  common.BytesToAddress([]byte("I'm address")),
		BlockHash:     common.BytesToHash([]byte("I'm hash")),
		BlockNum:      1,
		Extra:         make([]byte, 100),
	}
	var consensusMsg ConsensusMsg = pv
	cb, _ := consensusMsg.CannibalizeBytes()
	sign, _ := crypto.Sign(cb, pri)
	pv.Signature.SetBytes(sign)
	return pv
}

func TestSelectNodesByMsgType(t *testing.T) {
	router := newTestRouter()
	router.filter = func(p *peer, msgType uint64, condition interface{}) bool {
		if p.id == "test0" {
			return true
		}
		return false
	}
	// 1-3 consensus, 0 not exists, 4 non-consensus
	expects := []struct {
		msgType  uint64
		wantedId []string
	}{
		// minxingRandomNodes
		{PrepareBlockMsg, []string{"test1", "test2", "test3", "test4"}},
		{PrepareVoteMsg, []string{"test1", "test2", "test3", "test4"}},
		{ConfirmedPrepareBlockMsg, []string{"test1", "test2", "test3", "test4"}},
		{PrepareBlockHashMsg, []string{"test1", "test2", "test3", "test4"}},
		// consensusRandomNodes
		{ViewChangeMsg, []string{"test1", "test2", "test3"}},
	}
	for _, res := range expects {
		peers, err := router.selectNodesByMsgType(res.msgType, "")
		if err != nil {
			t.Error("Error occur", err)
		}
		for _, p := range peers {
			exist := false
			for _, w := range res.wantedId {
				if p.id == w {
					exist = true
				}
			}
			if exist {
				t.Fatalf("Select fail, result:%v, wanted:%v", p.id, res.wantedId)
			}
		}
	}
}

func TestKRandomNodes(t *testing.T) {
	peers := []*peer{}
	for i := 0; i < 90; i++ {
		// Simulating non-consensus nodes
		peers = append(peers, &peer{
			id: fmt.Sprintf("test%d", i),
		})
	}
	filterFunc := func(p *peer, k uint64, condition interface{}) bool {
		if p.id == "test0" {
			return true
		}
		return false
	}
	s1 := kRandomNodes(3, peers, 1, "", filterFunc)
	s2 := kRandomNodes(3, peers, 1, "", filterFunc)
	s3 := kRandomNodes(3, peers, 1, "", filterFunc)

	if reflect.DeepEqual(s1, s2) {
		t.Fatalf("Unexpected equal.")
	}
	if reflect.DeepEqual(s1, s3) {
		t.Fatalf("Unexpected equal.")
	}
	if reflect.DeepEqual(s2, s3) {
		t.Fatalf("Unexpected equal.")
	}
	for _, s := range [][]*peer{s1, s2, s3} {
		if len(s) != 3 {
			t.Fatalf("Bad length")
		}
		for _, n := range s {
			if n.id == "test0" {
				t.Fatal("Bad name")
			}
		}
	}
}

func TestFormatPeers(t *testing.T) {
	peers := []*peer{
		&peer{id: "id01"},
		&peer{id: "id02"},
	}
	peersStr := formatPeers(peers)
	if peersStr != "id01,id02" {
		t.Error("error")
	} else {
		t.Log("test success")
	}
}
