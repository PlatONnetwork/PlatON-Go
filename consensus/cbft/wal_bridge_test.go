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

package cbft

import (
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

	node := MockNode(pk[0], sk[0], cbftnodes, 10000, 10)
	assert.Nil(t, node.Start())
	node.engine.wal, _ = wal.NewWal(nil, tempDir)
	node.engine.bridge, _ = NewBridge(node.engine.nodeServiceContext, node.engine)
	node.engine.updateChainStateHook = node.engine.bridge.UpdateChainState

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	var commit, lock, qc *types.Block

	parent := node.chain.Genesis()
	for i := 0; i < 3; i++ {
		block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, node)
		assert.True(t, node.engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		node.engine.OnSeal(block, result, nil, complete)
		<-complete

		// test newChainState
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i), node.engine.state.MaxQCIndex())
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

	restartNode.engine.wal = node.engine.wal
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
	complete := make(chan struct{}, 1)
	var appendQC *types.Block
	node.engine.state.SetExecuting(1, true) // lockBlock

	// base lock seal duplicate qc
	block := NewBlockWithSign(lock.Hash(), lock.NumberU64()+1, node)
	assert.True(t, node.engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
	node.engine.OnSeal(block, result, nil, complete)
	<-complete

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

	pk, sk, cbftnodes := GenerateCbftNode(1)

	node := MockNode(pk[0], sk[0], cbftnodes, 10000, 20)
	assert.Nil(t, node.Start())
	node.engine.wal, _ = wal.NewWal(nil, tempDir)
	node.engine.bridge, _ = NewBridge(node.engine.nodeServiceContext, node.engine)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := node.chain.Genesis()

	epoch := node.engine.state.Epoch()
	viewNumber := node.engine.state.ViewNumber()
	_, qc := makePrepareQC(epoch, viewNumber, parent, 0)
	viewChangeQC := makeViewChangeQC(epoch, viewNumber, parent.NumberU64())
	node.engine.bridge.ConfirmViewChange(epoch, viewNumber, parent, qc, viewChangeQC, epoch, viewNumber)
	for i := 0; i < 10; i++ {
		block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, node)
		assert.True(t, node.engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		node.engine.OnSeal(block, result, nil, complete)
		<-complete

		select {
		case b := <-result:
			assert.NotNil(t, b)
			parent = b
		}
	}
	node.engine.bridge.SendViewChange(makeViewChange(epoch, viewNumber, parent, 9, uint32(0)))

	// test recoveryMsg
	restartNode := MockNode(pk[0], sk[0], cbftnodes, 10000, 10)
	assert.Nil(t, restartNode.Start())
	restartNode.engine.wal = node.engine.wal
	restartNode.engine.bridge, _ = NewBridge(restartNode.engine.nodeServiceContext, restartNode.engine)
	assert.Nil(t, restartNode.engine.wal.Load(restartNode.engine.recoveryMsg))

	block := restartNode.engine.state.ViewBlockByIndex(9)
	assert.Equal(t, parent.NumberU64(), block.NumberU64())
	assert.Equal(t, parent.Hash().String(), block.Hash().String())

	lastViewChangeQC, err := restartNode.engine.bridge.GetViewChangeQC(epoch, viewNumber)
	assert.Nil(t, err)
	assert.Equal(t, 3, len(lastViewChangeQC.QCs))
}

func TestInsertQCBlock_fork_priority(t *testing.T) {
	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	pk, sk, cbftnodes := GenerateCbftNode(1)

	node := MockNode(pk[0], sk[0], cbftnodes, 10000, 20)
	assert.Nil(t, node.Start())
	node.engine.wal, _ = wal.NewWal(nil, tempDir)
	node.engine.bridge, _ = NewBridge(node.engine.nodeServiceContext, node.engine)
	node.engine.updateChainStateHook = node.engine.bridge.UpdateChainState

	parent := node.chain.Genesis()

	var forkBlock *types.Block
	var forkQC *ctypes.QuorumCert
	for i := 0; i < 10; i++ {
		if i == 9 {
			node.engine.updateChainStateDelayHook = func(qcState, lockState, commitState *protocols.State) {
				time.Sleep(1 * time.Second)
				node.engine.bridge.UpdateChainState(qcState, lockState, commitState)
			}
		}
		block, qc := makePrepareQC(node.engine.state.Epoch(), node.engine.state.ViewNumber(), parent, uint32(i))
		node.engine.insertQCBlock(block, qc)
		if i == 9 {
			// moke fork block
			forkBlock, forkQC = makePrepareQC(node.engine.state.Epoch(), node.engine.state.ViewNumber()+1, parent, uint32(i))
			node.engine.insertQCBlock(forkBlock, forkQC)
		}
		parent = block
	}

	time.Sleep(2 * time.Second)
	assert.Equal(t, forkBlock.NumberU64(), node.engine.state.HighestQCBlock().NumberU64())
	assert.Equal(t, forkBlock.Hash(), node.engine.state.HighestQCBlock().Hash())

	var chainState *protocols.ChainState
	assert.Nil(t, node.engine.wal.LoadChainState(func(cs *protocols.ChainState) error {
		chainState = cs
		return nil
	}))

	assert.Equal(t, 2, len(chainState.QC))
	assert.Equal(t, forkBlock.NumberU64(), chainState.QC[0].QuorumCert.BlockNumber)
	assert.Equal(t, forkBlock.Hash(), chainState.QC[0].QuorumCert.BlockHash)
}

func makePrepareQC(epoch, viewNumber uint64, parent *types.Block, blockIndex uint32) (*types.Block, *ctypes.QuorumCert) {
	header := &types.Header{
		Number:      big.NewInt(int64(parent.NumberU64() + 1)),
		ParentHash:  parent.Hash(),
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		ReceiptHash: common.BytesToHash(utils.Rand32Bytes(32)),
		Root:        common.BytesToHash(utils.Rand32Bytes(32)),
	}
	block := types.NewBlockWithHeader(header)
	qc := &ctypes.QuorumCert{
		Epoch:        epoch,
		ViewNumber:   viewNumber,
		BlockHash:    block.Hash(),
		BlockNumber:  block.NumberU64(),
		BlockIndex:   blockIndex,
		Signature:    ctypes.BytesToSignature(utils.Rand32Bytes(64)),
		ValidatorSet: utils.NewBitArray(4),
	}
	return block, qc
}

func makeViewChangeQuorumCert(epoch, viewNumber uint64, blockNumber uint64) *ctypes.ViewChangeQuorumCert {
	return &ctypes.ViewChangeQuorumCert{
		Epoch:        epoch,
		ViewNumber:   viewNumber,
		BlockHash:    common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:  blockNumber,
		Signature:    ctypes.BytesToSignature(utils.Rand32Bytes(64)),
		ValidatorSet: utils.NewBitArray(25),
	}
}

func makeViewChangeQC(epoch, viewNumber uint64, blockNumber uint64) *ctypes.ViewChangeQC {
	return &ctypes.ViewChangeQC{
		QCs: []*ctypes.ViewChangeQuorumCert{
			makeViewChangeQuorumCert(epoch, viewNumber, blockNumber),
			makeViewChangeQuorumCert(epoch, viewNumber, blockNumber),
			makeViewChangeQuorumCert(epoch, viewNumber, blockNumber),
		},
	}
}

func makeViewChange(epoch, viewNumber uint64, block *types.Block, blockIndex uint32, validatorIndex uint32) *protocols.ViewChange {
	_, prepareQC := makePrepareQC(epoch, viewNumber, block, blockIndex)
	return &protocols.ViewChange{
		Epoch:          epoch,
		ViewNumber:     viewNumber,
		BlockHash:      block.Hash(),
		BlockNumber:    block.NumberU64(),
		ValidatorIndex: validatorIndex,
		PrepareQC:      prepareQC,
		Signature:      ctypes.BytesToSignature(utils.Rand32Bytes(64)),
	}
}
