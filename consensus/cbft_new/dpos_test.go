package cbft_new

import (
	"Platon-go/crypto"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"testing"
)

func newTesterAccountPool() ([]discover.Node, error) {
	var accounts []discover.Node
	for _, url := range params.MainnetBootnodes {
		node, err := discover.ParseNode(url)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, *node)
	}
	return accounts, nil
}

func TestDpos(t *testing.T) {
	nodes, _ := newTesterAccountPool()
	dpos := newDpos(nodes)
	dpos.SetLastCycleBlockNum(100)
	if dpos.LastCycleBlockNum() != 100 {
		t.Errorf("dpos.SetLastCycleBlockNum failed!")
	}

	nodeID := dpos.primaryNodeList[0]
	if dpos.NodeIndex(nodeID.ID) <= 0 {
		t.Errorf("dpos.CheckConsensusNode failed!")
	}

	addr, err := nodeID.ID.Pubkey()
	if err != nil || addr == nil {
		t.Errorf("nodeID.ID.Pubkey error!")
	}
	if !dpos.IsPrimary(crypto.PubkeyToAddress(*addr)) {
		t.Errorf("dpos.IsPrimary!")
	}
}
