package cbft

import (
	"Platon-go/core/types"
	"Platon-go/p2p/discover"
	"Platon-go/params"
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
