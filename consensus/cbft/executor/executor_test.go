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

package executor

import (
	"math/big"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

func NewBlock(parent common.Hash, number uint64) *types.Block {
	header := &types.Header{
		Number:     big.NewInt(int64(number)),
		ParentHash: parent,
		Time:       big.NewInt(time.Now().UnixNano()),
		Extra:      make([]byte, 77),
		Coinbase:   common.Address{},
		GasLimit:   10000000000,
	}
	block := types.NewBlockWithHeader(header)
	return block
}

func TestExecute(t *testing.T) {
	executor := func(block *types.Block, parent *types.Block) error {
		return nil
	}
	asyncExecutor := NewAsyncExecutor(executor)

	executeBlocks := make(map[uint64]*types.Block)
	parent := NewBlock(common.BytesToHash(utils.Rand32Bytes(32)), 1)
	for i := 0; i < 20; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		executeBlocks[block.NumberU64()] = block
		asyncExecutor.Execute(block, parent)
		parent = block
	}

	success := 0
loop:
	for {
		select {
		case result := <-asyncExecutor.ExecuteStatus():
			assert.Nil(t, result.Err)
			b := executeBlocks[result.Number]
			assert.NotNil(t, b)
			assert.Equal(t, b.Hash(), result.Hash)
			success = success + 1
			if b.NumberU64() > uint64(len(executeBlocks)) {
				break loop
			}
		}
	}

	assert.Equal(t, 20, success)
	asyncExecutor.Stop()
}
