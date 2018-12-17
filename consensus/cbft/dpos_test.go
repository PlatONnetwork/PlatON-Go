package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"math/big"
	"testing"
)

func newTesterAccountPool() ([]discover.NodeID, error) {
	var accounts []discover.NodeID
	for _, url := range params.MainnetBootnodes {
		node, err := discover.ParseNode(url)
		if err != nil {
			return nil, err
		}
		accounts = append(accounts, node.ID)
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
	if dpos.NodeIndex(nodeID) <= 0 {
		t.Errorf("dpos.CheckConsensusNode failed!")
	}

	addr, err := nodeID.Pubkey()
	if err != nil || addr == nil {
		t.Errorf("nodeID.ID.Pubkey error!")
	}
	if !dpos.IsPrimary(crypto.PubkeyToAddress(*addr)) {
		t.Errorf("dpos.IsPrimary!")
	}
}

func TestBlockCopy(t *testing.T) {
	header := &types.Header{
		Number: big.NewInt(1),
	}
	block1 := types.NewBlock(header, nil, nil, nil)

	block2Obj := *block1

	block2Obj.Header().Number = big.NewInt(2)

	block2 := &block2Obj
	block2.Header().Number = big.NewInt(3)

	println(block1.Number().Uint64())
	println(block2.Number().Uint64())

}
