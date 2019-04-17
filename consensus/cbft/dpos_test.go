package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/stretchr/testify/assert"
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
	//if dpos.NodeIndex(nodeID) <= 0 {
	//	t.Errorf("dpos.CheckConsensusNode failed!")
	//}

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
	block1 := types.NewBlock(header, nil, nil)

	block2Obj := *block1

	block2Obj.Header().Number = big.NewInt(2)

	block2 := &block2Obj
	block2.Header().Number = big.NewInt(3)

	println(block1.Number().Uint64())
	println(block2.Number().Uint64())

}

func TestHash(t *testing.T) {
	h := common.Hash{}
	t.Log(fmt.Sprintf("%v", h))
}
func TestNodeIndex(t *testing.T) {
	addrs := []string{
		"86b1fb3c765dd90433e1d4028b26fcd19e16b2167f6e87cd02194a393173c265fcf14484415fb1d116c07464dabf944d1dc2e98c05745b9e7cb3b4b8cb49e507",
		"3283bb46341277196744e958d4b7287725c6deaff42eb69d15d1c5702362d53443138b603398a9d278f6e5c212a2ef9d5a1047b98997eb1303e52813a71ea553",
		"27b6e1b4523c26b6aa0feaa8bffc8362de2365367096bd6a4fdd72de3132e6a2fe53f213f81179f089b9543afe8f8f131926c9297d11edc1fc6eeca1707f4ed3",
		"b6d2dcce5963cc512b178489b6ebcce56ba2baefe34e06538d1b958c838b57819cf3bb956bf7aef38cc3d01ace2096051342d6b255099f78427fb4e72f817408",
	}

	nodes := make([]discover.NodeID, 0)

	for _, addr := range addrs {
		id, _ := discover.HexID(addr)
		nodes = append(nodes, id)
	}

	dpos := newDpos(nodes)
	index, addr, _ := dpos.NodeIndexAddress(nodes[1])
	t.Log(addr.String())
	assert.Equal(t, index, 1)
	index, _ = dpos.AddressIndex(common.HexToAddress("0x0C3eacBA94Fda90912798fDaC1a3A9d15C4F3388"))
	assert.Equal(t, index, 1)

}
