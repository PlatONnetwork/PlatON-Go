package types

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/assert"
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

	tree := NewBlockTree(root, nil)

	for _, f := range forks {
		for _, b := range f {
			if b.NumberU64() == 0 {
				tree.InsertQCBlock(b, nil)
			}
			tree.InsertQCBlock(b, &QuorumCert{ViewNumber: 1})
		}
	}

	tree.PruneBlock(fork1[1].Hash(), fork1[1].NumberU64(), func(block *types.Block) {
		for _, b := range fork2 {
			if b.Hash() == block.Hash() {
				return
			}
		}
		t.Error(fmt.Sprintf("Clear Block failed"))
	})

	for _, b := range fork1[1:] {
		assert.NotNil(t, tree.FindBlockByHash(b.Hash()))
		b, q := tree.FindBlockAndQC(b.Hash(), b.NumberU64())
		assert.NotNil(t, b)
		assert.NotNil(t, q)
	}
	for _, b := range fork2[1:] {
		assert.Nil(t, tree.FindBlockByHash(b.Hash()))
		b, q := tree.FindBlockAndQC(b.Hash(), b.NumberU64())
		assert.Nil(t, b)
		assert.Nil(t, q)
	}
}
