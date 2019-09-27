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
