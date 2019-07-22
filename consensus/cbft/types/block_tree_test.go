package types

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"math/big"
	"testing"
	"time"
)

func newBlock(parent common.Hash, number uint64) *types.Block {
	header := &types.Header{
		Number:     big.NewInt(int64(number)),
		ParentHash: parent,
		Time:       big.NewInt(time.Now().UnixNano()),
		Extra:      nil,
	}
	block := types.NewBlockWithHeader(header)
	return block
}
func TestNewBlockTree(t *testing.T) {

	root := newBlock(common.Hash{}, 0)
	fork1 := types.Blocks{
		root,
	}
	fork2 := types.Blocks{
		root,
	}
	for i := uint64(0); i < 5; i++ {
		fork1 = append(fork1, newBlock(fork1[i].Hash(), i+1))
	}
	for i := uint64(0); i < 5; i++ {
		fork2 = append(fork2, newBlock(fork2[i].Hash(), i+1))
	}
	forks := []types.Blocks{
		fork1,
		fork2,
	}

	tree := NewBlockTree(root, &QuorumCert{ViewNumber: 1})

	for _, f := range forks {
		for i, b := range f {
			t.Log(i)

			tree.InsertQCBlock(b, &QuorumCert{ViewNumber: 1})
		}
	}

	tree.PruneBlock(fork1[1].Hash(), fork1[1].NumberU64(), nil)

	t.Log(len(tree.blocks))
}
