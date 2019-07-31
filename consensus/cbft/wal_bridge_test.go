package cbft

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/stretchr/testify/assert"
)

func TestUpdateChainState(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	pk, sk, cbftnodes := GenerateCbftNode(1)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 1; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		node.engine.wal, _ = wal.NewWal(nil, tempDir)
		node.engine.bridge, _ = NewBridge(node.engine.nodeServiceContext, node.engine)
		nodes = append(nodes, node)
		fmt.Println(i, node.engine.config.Option.NodeID.TerminalString())
	}

	result := make(chan *types.Block, 1)
	var commit, lock, qc *types.Block

	parent := nodes[0].chain.Genesis()
	for i := 0; i < 3; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil)

		// test UpdateChainState
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i), nodes[0].engine.state.MaxQCIndex())
			switch i {
			case 0:
				commit = block
			case 1:
				lock = block
			case 2:
				qc = block
			}
			parent = b
		}
	}

	// test recoveryChainState
	restartNode := MockNode(pk[0], sk[0], cbftnodes, 10000, 10)
	assert.Nil(t, restartNode.Start())

	restartNode.engine.wal = nodes[0].engine.wal
	restartNode.engine.bridge, _ = NewBridge(restartNode.engine.nodeServiceContext, restartNode.engine)

	assert.Nil(t, restartNode.engine.wal.LoadChainState(restartNode.engine.recoveryChainState))

	// check viewBlocks and viewQCs
	commitBlock, commitQC := restartNode.engine.state.ViewBlockAndQC(0)
	assert.Equal(t, commit.NumberU64(), commitBlock.NumberU64())
	assert.Equal(t, commit.NumberU64(), commitQC.BlockNumber)
	lockBlock, lockQC := restartNode.engine.state.ViewBlockAndQC(1)
	assert.Equal(t, lock.NumberU64(), lockBlock.NumberU64())
	assert.Equal(t, lock.NumberU64(), lockQC.BlockNumber)
	qcBlock, qcQC := restartNode.engine.state.ViewBlockAndQC(2)
	assert.Equal(t, qc.NumberU64(), qcBlock.NumberU64())
	assert.Equal(t, qc.NumberU64(), qcQC.BlockNumber)
	assert.Equal(t, qc.Hash().TerminalString(), qcBlock.Hash().TerminalString())
	assert.Equal(t, qc.Hash().TerminalString(), qcQC.BlockHash.TerminalString())

	// check blockTree
	commitBlock, commitQC = restartNode.engine.blockTree.FindBlockAndQC(commit.Hash(), commit.NumberU64())
	assert.Equal(t, commit.NumberU64(), commitBlock.NumberU64())
	assert.Equal(t, commit.Hash(), commitBlock.Hash())
	lockBlock, lockQC = restartNode.engine.blockTree.FindBlockAndQC(lock.Hash(), lock.NumberU64())
	assert.Equal(t, lock.NumberU64(), lockBlock.NumberU64())
	qcBlock, qcQC = restartNode.engine.blockTree.FindBlockAndQC(qc.Hash(), qc.NumberU64())
	assert.Equal(t, qc.NumberU64(), qcBlock.NumberU64())

	// check highest
	highestCommitBlockNumber, _ := restartNode.engine.HighestCommitBlockBn()
	assert.Equal(t, commit.NumberU64(), highestCommitBlockNumber)
	highestLockBlockNumber, _ := restartNode.engine.HighestLockBlockBn()
	assert.Equal(t, lock.NumberU64(), highestLockBlockNumber)
	highestQCBlockNumber, _ := restartNode.engine.HighestQCBlockBn()
	assert.Equal(t, qc.NumberU64(), highestQCBlockNumber)

	// test addQCState
	testAddQCState(t, lock, qc, restartNode)
}

func testAddQCState(t *testing.T, lock, qc *types.Block, node *TestCBFT) {
	result := make(chan *types.Block, 1)
	var appendQC *types.Block
	node.engine.state.SetExecuting(1, true) // lockBlock

	block := NewBlock(lock.Hash(), lock.NumberU64()+1)
	assert.True(t, node.engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	node.engine.OnSeal(block, result, nil)

	// test addQCState
	select {
	case b := <-result:
		assert.NotNil(t, b)
		appendQC = block
	}

	// test recoveryChainState
	var chainState *protocols.ChainState
	assert.Nil(t, node.engine.wal.LoadChainState(func(cs *protocols.ChainState) error {
		chainState = cs
		return nil
	}))

	qcs := chainState.QC
	assert.Equal(t, 2, len(qcs))

	assert.Equal(t, qc.Hash().TerminalString(), qcs[0].Block.Hash().TerminalString())
	assert.Equal(t, appendQC.Hash().TerminalString(), qcs[1].Block.Hash().TerminalString())
	assert.Equal(t, lock.Hash().TerminalString(), qcs[1].Block.ParentHash().TerminalString())
}

func TestRecordCbftMsg(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	pk, sk, cbftnodes := GenerateCbftNode(1)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 1; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		node.engine.wal, _ = wal.NewWal(nil, tempDir)
		node.engine.bridge, _ = NewBridge(node.engine.nodeServiceContext, node.engine)
		nodes = append(nodes, node)
		fmt.Println(i, node.engine.config.Option.NodeID.TerminalString())
	}

	result := make(chan *types.Block, 1)
	parent := nodes[0].chain.Genesis()
	var recoveryPoint *types.Block
	// test changeView, when i=10 it will change view and ConfirmViewChange
	for i := 0; i < 19; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil)

		select {
		case b := <-result:
			assert.NotNil(t, b)
			if block.NumberU64() == 10 {
				recoveryPoint = block
			}
			parent = b
		}
	}

	// test recoveryChainState
	nodes[0].engine.state.SetHighestQCBlock(recoveryPoint)
	assert.Nil(t, nodes[0].engine.wal.Load(nodes[0].engine.recoveryMsg))

	block := nodes[0].engine.state.ViewBlockByIndex(8)
	assert.Equal(t, parent.NumberU64(), block.NumberU64())
	assert.Equal(t, parent.Hash().TerminalString(), block.Hash().TerminalString())

	assert.True(t, nodes[0].engine.state.HadSendPrepareVote().Had(8))

	assert.Equal(t, 1, nodes[0].engine.state.PrepareVoteLenByIndex(8))
}
