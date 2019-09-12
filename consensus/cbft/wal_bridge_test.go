package cbft

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
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
		node.engine.updateChainStateHook = node.engine.bridge.UpdateChainState
		nodes = append(nodes, node)
		fmt.Println(i, node.engine.config.Option.NodeID.String())
	}

	result := make(chan *types.Block, 1)
	var commit, lock, qc *types.Block

	parent := nodes[0].chain.Genesis()
	for i := 0; i < 3; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil)

		// test newChainState
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
	assert.Equal(t, qc.Hash().String(), qcBlock.Hash().String())
	assert.Equal(t, qc.Hash().String(), qcQC.BlockHash.String())

	// check blockTree
	commitBlock, _ = restartNode.engine.blockTree.FindBlockAndQC(commit.Hash(), commit.NumberU64())
	assert.Equal(t, commit.NumberU64(), commitBlock.NumberU64())
	assert.Equal(t, commit.Hash(), commitBlock.Hash())
	lockBlock, _ = restartNode.engine.blockTree.FindBlockAndQC(lock.Hash(), lock.NumberU64())
	assert.Equal(t, lock.NumberU64(), lockBlock.NumberU64())
	qcBlock, _ = restartNode.engine.blockTree.FindBlockAndQC(qc.Hash(), qc.NumberU64())
	assert.Equal(t, qc.NumberU64(), qcBlock.NumberU64())

	// check highest
	highestCommitBlockNumber, _ := restartNode.engine.HighestCommitBlockBn()
	assert.Equal(t, commit.NumberU64(), highestCommitBlockNumber)
	highestLockBlockNumber, _ := restartNode.engine.HighestLockBlockBn()
	assert.Equal(t, lock.NumberU64(), highestLockBlockNumber)
	highestQCBlockNumber, _ := restartNode.engine.HighestQCBlockBn()
	assert.Equal(t, qc.NumberU64(), highestQCBlockNumber)

	// check blockChain currentHeader
	assert.Equal(t, commit.Hash(), restartNode.engine.blockChain.CurrentHeader().Hash())
	assert.Equal(t, commit.NumberU64(), restartNode.engine.blockChain.CurrentHeader().Number.Uint64())

	// test addQCState
	testAddQCState(t, lock, qc, restartNode)
}

func testAddQCState(t *testing.T, lock, qc *types.Block, node *TestCBFT) {
	result := make(chan *types.Block, 1)
	var appendQC *types.Block
	node.engine.state.SetExecuting(1, true) // lockBlock

	// base lock seal duplicate qc
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

	assert.Equal(t, qc.Hash().String(), qcs[0].Block.Hash().String())
	assert.Equal(t, appendQC.Hash().String(), qcs[1].Block.Hash().String())
	assert.Equal(t, lock.Hash().String(), qcs[1].Block.ParentHash().String())
}

func TestRecordCbftMsg(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		if i == 0 {
			node.engine.wal, _ = wal.NewWal(nil, tempDir)
			node.engine.bridge, _ = NewBridge(node.engine.nodeServiceContext, node.engine)
		}
		nodes = append(nodes, node)
		fmt.Println(i, node.engine.config.Option.NodeID.String())
	}

	result := make(chan *types.Block, 1)
	parent := nodes[0].chain.Genesis()

	epoch := nodes[0].engine.state.Epoch()
	viewNumber := nodes[0].engine.state.ViewNumber()
	nodes[0].engine.bridge.ConfirmViewChange(epoch, viewNumber, parent, &ctypes.QuorumCert{Signature: ctypes.BytesToSignature(utils.Rand32Bytes(32)), ValidatorSet: utils.NewBitArray(110)}, nil)
	for i := 0; i < 10; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil)

		select {
		case b := <-result:
			assert.NotNil(t, b)
			parent = b
		}
	}

	// test recoveryMsg
	restartNode := MockNode(pk[0], sk[0], cbftnodes, 10000, 10)
	assert.Nil(t, restartNode.Start())
	restartNode.engine.wal = nodes[0].engine.wal
	restartNode.engine.bridge, _ = NewBridge(restartNode.engine.nodeServiceContext, restartNode.engine)

	assert.Nil(t, restartNode.engine.wal.Load(restartNode.engine.recoveryMsg))

	block := restartNode.engine.state.ViewBlockByIndex(9)
	assert.Equal(t, parent.NumberU64(), block.NumberU64())
	assert.Equal(t, parent.Hash().String(), block.Hash().String())
}

func TestInsertQCBlock_fork_priority(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	pk, sk, cbftnodes := GenerateCbftNode(1)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 1; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 20)
		assert.Nil(t, node.Start())

		node.engine.wal, _ = wal.NewWal(nil, tempDir)
		node.engine.bridge, _ = NewBridge(node.engine.nodeServiceContext, node.engine)
		node.engine.updateChainStateHook = node.engine.bridge.UpdateChainState
		nodes = append(nodes, node)
	}

	parent := nodes[0].chain.Genesis()

	var forkBlock *types.Block
	var forkQC *ctypes.QuorumCert
	for i := 0; i < 10; i++ {
		if i == 9 {
			nodes[0].engine.updateChainStateDelayHook = func(qcState, lockState, commitState *protocols.State) {
				time.Sleep(1 * time.Second)
				nodes[0].engine.bridge.UpdateChainState(qcState, lockState, commitState)
			}
		}
		block, qc := makePrepareQC(nodes[0], parent, uint32(i), nodes[0].engine.state.ViewNumber())
		nodes[0].engine.insertQCBlock(block, qc)
		if i == 9 {
			// moke fork block
			forkBlock, forkQC = makePrepareQC(nodes[0], parent, uint32(i), nodes[0].engine.state.ViewNumber()+1)
			nodes[0].engine.insertQCBlock(forkBlock, forkQC)
		}
		parent = block
	}

	time.Sleep(2 * time.Second)
	assert.Equal(t, forkBlock.NumberU64(), nodes[0].engine.state.HighestQCBlock().NumberU64())
	assert.Equal(t, forkBlock.Hash(), nodes[0].engine.state.HighestQCBlock().Hash())

	var chainState *protocols.ChainState
	assert.Nil(t, nodes[0].engine.wal.LoadChainState(func(cs *protocols.ChainState) error {
		chainState = cs
		return nil
	}))

	assert.Equal(t, 2, len(chainState.QC))
	assert.Equal(t, forkBlock.NumberU64(), chainState.QC[0].QuorumCert.BlockNumber)
	assert.Equal(t, forkBlock.Hash(), chainState.QC[0].QuorumCert.BlockHash)
}

func makePrepareQC(node *TestCBFT, parent *types.Block, blockIndex uint32, viewNumber uint64) (*types.Block, *ctypes.QuorumCert) {
	header := &types.Header{
		Number:      big.NewInt(int64(parent.NumberU64() + 1)),
		ParentHash:  parent.Hash(),
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 77),
		ReceiptHash: common.BytesToHash(utils.Rand32Bytes(32)),
		Root:        common.BytesToHash(utils.Rand32Bytes(32)),
	}
	block := types.NewBlockWithHeader(header)
	qc := &ctypes.QuorumCert{
		Epoch:        node.engine.state.Epoch(),
		ViewNumber:   viewNumber,
		BlockHash:    block.Hash(),
		BlockNumber:  block.NumberU64(),
		BlockIndex:   blockIndex,
		Signature:    ctypes.Signature{},
		ValidatorSet: utils.NewBitArray(32),
	}
	return block, qc
}
