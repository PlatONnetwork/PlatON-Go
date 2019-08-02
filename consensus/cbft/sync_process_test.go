package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	types2 "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

//func init() {
//	log.Root().SetHandler(log.StdoutHandler)
//
//}

func TestFetch(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 200000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
		fmt.Println(i, node.engine.config.Option.NodeID.TerminalString())
	}

	result := make(chan *types.Block, 1)

	var fetchBlock *types.Block
	qcBlocks := &protocols.QCBlockList{}
	parent := nodes[0].chain.Genesis()
	for i := 0; i < 3; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil)
		fetchBlock = block
		qcBlocks.Blocks = append(qcBlocks.Blocks, block)
		_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i-1), nodes[0].engine.state.MaxQCIndex())
			for j := 1; j < 4; j++ {
				msg := &protocols.PrepareVote{
					Epoch:          nodes[0].engine.state.Epoch(),
					ViewNumber:     nodes[0].engine.state.ViewNumber(),
					BlockIndex:     uint32(i),
					BlockHash:      b.Hash(),
					BlockNumber:    b.NumberU64(),
					ValidatorIndex: uint32(j),
					ParentQC:       qc,
				}
				pb := nodes[0].engine.state.PrepareBlockByIndex(uint32(i))
				assert.NotNil(t, pb)
				assert.Nil(t, nodes[j].engine.OnPrepareBlock("id", pb))
				time.Sleep(1000 * time.Millisecond)
				index, finish := nodes[j].engine.state.Executing()
				assert.True(t, index == uint32(i) && finish, fmt.Sprintf("%d,%v", index, finish))
				assert.Nil(t, nodes[j].engine.signMsgByBls(msg))
				assert.Nil(t, nodes[0].engine.OnPrepareVote("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
			}
			_, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
			assert.NotNil(t, qc)
			qcBlocks.QC = append(qcBlocks.QC, qc)
			parent = b
		}
	}
	assert.Equal(t, uint64(3), nodes[0].engine.state.HighestQCBlock().NumberU64())
	assert.Equal(t, uint64(0), nodes[1].engine.state.HighestQCBlock().NumberU64())

	nodes[1].engine.fetchBlock("id", fetchBlock.Hash(), fetchBlock.NumberU64())
	nodes[1].engine.syncMsgCh <- &types2.MsgInfo{PeerID: "id", Msg: qcBlocks}
	time.Sleep(100 * time.Millisecond)
	assert.Equal(t, uint64(3), nodes[1].engine.state.HighestQCBlock().NumberU64())
}

func TestSyncBlock(t *testing.T) {

	nodes := Mock4NodePipe(false)

	result := make(chan *types.Block, 1)

	var fetchBlock *types.Block
	qcBlocks := &protocols.QCBlockList{}
	parent := nodes[0].chain.Genesis()
	for i := 0; i < 3; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil)
		fetchBlock = block
		qcBlocks.Blocks = append(qcBlocks.Blocks, block)
		_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i-1), nodes[0].engine.state.MaxQCIndex())
			for j := 1; j < 4; j++ {
				msg := &protocols.PrepareVote{
					Epoch:          nodes[0].engine.state.Epoch(),
					ViewNumber:     nodes[0].engine.state.ViewNumber(),
					BlockIndex:     uint32(i),
					BlockHash:      b.Hash(),
					BlockNumber:    b.NumberU64(),
					ValidatorIndex: uint32(j),
					ParentQC:       qc,
				}
				pb := nodes[0].engine.state.PrepareBlockByIndex(uint32(i))
				assert.NotNil(t, pb)
				assert.Nil(t, nodes[j].engine.OnPrepareBlock(nodes[0].engine.config.Option.NodeID.TerminalString(), pb))
				time.Sleep(50 * time.Millisecond)
				index, finish := nodes[j].engine.state.Executing()
				assert.True(t, index == uint32(i) && finish, fmt.Sprintf("%d,%v", index, finish))
				assert.Nil(t, nodes[j].engine.signMsgByBls(msg))
				assert.Nil(t, nodes[0].engine.OnPrepareVote(nodes[j].engine.config.Option.NodeID.TerminalString(), msg), fmt.Sprintf("number:%d", b.NumberU64()))
			}
			_, qc := nodes[0].engine.blockTree.FindBlockAndQC(block.Hash(), block.NumberU64())
			assert.NotNil(t, qc)
			qcBlocks.QC = append(qcBlocks.QC, qc)
			parent = b
		}
	}
	assert.Equal(t, uint64(3), nodes[0].engine.state.HighestQCBlock().NumberU64())
	assert.Equal(t, uint64(0), nodes[1].engine.state.HighestQCBlock().NumberU64())
	assert.Equal(t, uint64(0), nodes[2].engine.state.HighestQCBlock().NumberU64())

	assert.Equal(t, 0, nodes[1].engine.fetcher.Len())
	assert.Equal(t, uint64(0), nodes[1].engine.state.HighestQCBlock().NumberU64())

	for i := 0; i < 4; i++ {
		nodes[i].engine.network.Testing()
	}

	nodes[1].engine.fetchBlock(nodes[0].engine.config.Option.NodeID.TerminalString(), fetchBlock.Hash(), fetchBlock.NumberU64())
	//nodes[1].engine.syncMsgCh <- &types2.MsgInfo{PeerID: "id", Msg: qcBlocks}
	time.Sleep(1000 * time.Millisecond)
	assert.Equal(t, uint64(3), nodes[1].engine.state.HighestQCBlock().NumberU64())

}
