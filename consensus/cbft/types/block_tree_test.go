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

package types

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/assert"
)

func newBlock(parent common.Hash, number uint64) *types.Block {
	header := &types.Header{
		Number:      big.NewInt(int64(number)),
		ParentHash:  parent,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       nil,
		ReceiptHash: common.BytesToHash(utils.Rand32Bytes(32)),
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

	newRoot := fork1[len(fork1)-1]
	tree.Reset(newRoot, &QuorumCert{ViewNumber: 1})
	assert.Equal(t, tree.root.Block.Hash(), newRoot.Hash())
	assert.Equal(t, 1, len(tree.blocks))
}

func Test_blockTree_newRoot(t *testing.T) {
	root := newBlock(common.Hash{}, 0)
	tree := NewBlockTree(root, nil)
	parent := root
	for i := 0; i < 10; i++ {
		block := newBlock(parent.Hash(), uint64(i+1))
		tree.InsertQCBlock(block, &QuorumCert{ViewNumber: 1})
		parent = block
	}
	assert.Equal(t, 11, len(tree.blocks))
	tree.NewRoot(parent)
	assert.Equal(t, 1, len(tree.blocks))
	_, err := tree.root.MarshalJSON()
	assert.Nil(t, err)
	_, err = tree.MarshalJSON()
	assert.Nil(t, err)
	assert.Equal(t, parent.Hash(), tree.root.Block.Hash())
}
