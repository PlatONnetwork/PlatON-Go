package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTryViewChange(t *testing.T) {

	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
	}

	result := make(chan *types.Block, 1)

	parent := nodes[0].chain.Genesis()
	for i := 0; i < 4; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil)

		_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i-1), nodes[0].engine.state.MaxQCIndex())
			for j := 1; j < 3; j++ {
				msg := &protocols.PrepareVote{
					Epoch:          nodes[0].engine.state.Epoch(),
					ViewNumber:     nodes[0].engine.state.ViewNumber(),
					BlockIndex:     uint32(i),
					BlockHash:      b.Hash(),
					BlockNumber:    b.NumberU64(),
					ValidatorIndex: uint32(j),
					ParentQC:       qc,
				}
				assert.Nil(t, nodes[j].engine.signMsgByBls(msg))
				assert.Nil(t, nodes[0].engine.OnPrepareVote("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
			}
			parent = b
		}
	}
	time.Sleep(10 * time.Second)

	block := nodes[0].engine.state.HighestQCBlock()
	block, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())

	for i := 0; i < 4; i++ {
		epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()
		viewchange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      block.Hash(),
			BlockNumber:    block.NumberU64(),
			ValidatorIndex: uint32(i),
			PrepareQC:      qc,
		}
		assert.Nil(t, nodes[i].engine.signMsgByBls(viewchange))
		assert.Nil(t, nodes[0].engine.OnViewChange("id", viewchange))
	}
	assert.NotNil(t, nodes[0].engine.state.LastViewChangeQC())

	assert.Equal(t, uint64(1), nodes[0].engine.state.ViewNumber())

}
