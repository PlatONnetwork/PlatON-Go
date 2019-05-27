package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"reflect"
	"testing"
)

var peersData = map[string]*peer{
	"test0": &peer{
		id: "test0",
	},
	"test1": &peer{
		id: "test1",
	},
	"test2": &peer{
		id: "test2",
	},
	"test3": &peer{
		id: "test3",
	},
	"test4": &peer{
		id: "test4",
	},
}

var TestPeerSet = &peerSet{
	peers: peersData,
}

func newTestRouter() *router {
	return &router{
		msgHandler: &handler{
			peers: TestPeerSet,
			cbft:  newCbft(),
		},
	}
}

func newCbft() *Cbft {
	config := &params.CbftConfig{
		InitialNodes: []discover.Node{
			{ID: StringID("test1")},
			{ID: StringID("test2")},
			{ID: StringID("test5")},
		},
	}

	return New(config, &event.TypeMux{}, nil)
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
			if !exist {
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

func StringID(s string) discover.NodeID {
	var id discover.NodeID
	b := []byte(s)
	copy(id[:], b)
	return id
}

func TestFormatPeers(t *testing.T) {
	peers := []*peer{
		&peer{
			id: "id01",
		},
		&peer{
			id: "id02",
		},
	}
	peersStr := formatPeers(peers)
	if peersStr != "id01,id02" {
		t.Error("error")
	} else {
		t.Log("test success")
	}
}