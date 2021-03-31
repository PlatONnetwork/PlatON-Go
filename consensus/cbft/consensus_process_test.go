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
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/validator"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/assert"
)

func TestViewChange(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
	}

	// TestTryViewChange
	testTryViewChange(t, nodes)

	// TestTryChangeViewByViewChange
	testTryChangeViewByViewChange(t, nodes)
}

func testTryViewChange(t *testing.T, nodes []*TestCBFT) {
	tempDir, _ := ioutil.TempDir("", "wal")
	defer os.RemoveAll(tempDir)

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()

	nodes[0].engine.wal, _ = wal.NewWal(nil, tempDir)
	nodes[0].engine.bridge, _ = NewBridge(nodes[0].engine.nodeServiceContext, nodes[0].engine)

	for i := 0; i < 4; i++ {
		block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil, complete)
		<-complete

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
		assert.Nil(t, nodes[0].engine.OnViewChanges("id", &protocols.ViewChanges{
			VCs: []*protocols.ViewChange{
				viewchange,
			},
		}))
	}
	assert.NotNil(t, nodes[0].engine.state.LastViewChangeQC())
	assert.Equal(t, uint64(1), nodes[0].engine.state.ViewNumber())
	lastViewChangeQC, _ := nodes[0].engine.bridge.GetViewChangeQC(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber())
	assert.Nil(t, lastViewChangeQC)
	lastViewChangeQC, _ = nodes[0].engine.bridge.GetViewChangeQC(nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()-1)
	assert.NotNil(t, lastViewChangeQC)
	epoch, viewNumber, _, _, _, blockNumber := lastViewChangeQC.MaxBlock()
	assert.Equal(t, nodes[0].engine.state.Epoch(), epoch)
	assert.Equal(t, nodes[0].engine.state.ViewNumber()-1, viewNumber)
	assert.Equal(t, uint64(4), blockNumber)
}

func testTryChangeViewByViewChange(t *testing.T, nodes []*TestCBFT) {
	// note: node-0 has been successfully switched to view-1, HighestQC blockNumber = 4
	// build a duplicate block-4
	number, hash := nodes[0].engine.HighestQCBlockBn()
	block, _ := nodes[0].engine.blockTree.FindBlockAndQC(hash, number)
	dulBlock := NewBlock(block.ParentHash(), block.NumberU64())
	_, preQC := nodes[0].engine.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
	// Vote and generate prepareQC for dulBlock
	votes := make(map[uint32]*protocols.PrepareVote)
	for j := 1; j < 3; j++ {
		vote := &protocols.PrepareVote{
			Epoch:          nodes[0].engine.state.Epoch(),
			ViewNumber:     nodes[0].engine.state.ViewNumber(),
			BlockIndex:     uint32(0),
			BlockHash:      dulBlock.Hash(),
			BlockNumber:    dulBlock.NumberU64(),
			ValidatorIndex: uint32(j),
			ParentQC:       preQC,
		}
		assert.Nil(t, nodes[j].engine.signMsgByBls(vote))
		votes[uint32(j)] = vote
	}
	dulQC := nodes[0].engine.generatePrepareQC(votes)
	// build new viewChange
	viewChanges := make(map[uint32]*protocols.ViewChange)
	for j := 1; j < 3; j++ {
		viewchange := &protocols.ViewChange{
			Epoch:          nodes[0].engine.state.Epoch(),
			ViewNumber:     nodes[0].engine.state.ViewNumber(),
			BlockHash:      dulBlock.Hash(),
			BlockNumber:    dulBlock.NumberU64(),
			ValidatorIndex: uint32(j),
			PrepareQC:      dulQC,
		}
		assert.Nil(t, nodes[j].engine.signMsgByBls(viewchange))
		viewChanges[uint32(j)] = viewchange
	}
	viewChangeQC := nodes[0].engine.generateViewChangeQC(viewChanges)

	// Case1: local highestqc is behind other validators and not exist viewChangeQC.maxBlock, sync qc block
	nodes[0].engine.tryChangeViewByViewChange(viewChangeQC)
	assert.Equal(t, uint64(1), nodes[0].engine.state.ViewNumber())
	assert.Equal(t, hash, nodes[0].engine.state.HighestQCBlock().Hash())

	// Case2: local highestqc is equal other validators and exist viewChangeQC.maxBlock, change the view
	nodes[0].engine.insertQCBlock(dulBlock, dulQC)
	nodes[0].engine.tryChangeViewByViewChange(viewChangeQC)
	assert.Equal(t, uint64(2), nodes[0].engine.state.ViewNumber())
	assert.Equal(t, dulQC.BlockHash, nodes[0].engine.state.HighestQCBlock().Hash())

	// based on the view-2 build a duplicate block-4
	dulBlock = NewBlock(block.ParentHash(), block.NumberU64())
	// Vote and generate prepareQC for dulBlock
	votes = make(map[uint32]*protocols.PrepareVote)
	for j := 1; j < 3; j++ {
		vote := &protocols.PrepareVote{
			Epoch:          nodes[0].engine.state.Epoch(),
			ViewNumber:     nodes[0].engine.state.ViewNumber(),
			BlockIndex:     uint32(0),
			BlockHash:      dulBlock.Hash(),
			BlockNumber:    dulBlock.NumberU64(),
			ValidatorIndex: uint32(j),
			ParentQC:       preQC,
		}
		assert.Nil(t, nodes[j].engine.signMsgByBls(vote))
		votes[uint32(j)] = vote
	}
	dulQC = nodes[0].engine.generatePrepareQC(votes)
	nodes[0].engine.blockTree.InsertQCBlock(dulBlock, dulQC)
	nodes[0].engine.state.SetHighestQCBlock(dulBlock)
	// Case3: local highestqc is ahead other validators, and not send viewChange, generate new viewChange quorumCert and change the view
	nodes[0].engine.tryChangeViewByViewChange(viewChangeQC)
	assert.Equal(t, uint64(3), nodes[0].engine.state.ViewNumber())
	assert.Equal(t, dulQC.BlockHash, nodes[0].engine.state.HighestQCBlock().Hash())
	_, _, _, blockView, _, _ := viewChangeQC.MaxBlock()
	assert.Equal(t, uint64(2), blockView)
}

type testCase struct {
	hadViewTimeout bool
}

func TestRichViewChangeQC(t *testing.T) {
	tests := []testCase{
		{true},
		{false},
	}
	for _, c := range tests {
		testRichViewChangeQCCase(t, c)
	}
}

func testRichViewChangeQCCase(t *testing.T, c testCase) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		assert.Nil(t, node.Start())

		nodes = append(nodes, node)
	}

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()
	for i := 0; i < 4; i++ {
		block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil, complete)
		<-complete

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
	if c.hadViewTimeout {
		time.Sleep(10 * time.Second)
	}

	hadSend := nodes[0].engine.state.ViewChangeByIndex(0)
	if c.hadViewTimeout {
		assert.NotNil(t, hadSend)
	}
	if !c.hadViewTimeout {
		assert.Nil(t, hadSend)
	}

	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qc := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockqc := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())

	viewChanges := make(map[uint32]*protocols.ViewChange, 0)

	for i := 1; i < 4; i++ {
		epoch, view := nodes[0].engine.state.Epoch(), nodes[0].engine.state.ViewNumber()
		v := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      lockBlock.Hash(), // base lock qc
			BlockNumber:    lockBlock.NumberU64(),
			ValidatorIndex: uint32(i),
			PrepareQC:      lockqc,
		}
		assert.Nil(t, nodes[i].engine.signMsgByBls(v))
		viewChanges[v.ValidatorIndex] = v
	}

	viewChangeQC := nodes[0].engine.generateViewChangeQC(viewChanges)
	nodes[0].engine.richViewChangeQC(viewChangeQC)

	epoch, viewNumber, blockEpoch, blockViewNumber, blockHash, blockNumber := viewChangeQC.MaxBlock()
	assert.Equal(t, qc.Epoch, epoch)
	assert.Equal(t, qc.ViewNumber, viewNumber)
	assert.Equal(t, qc.Epoch, blockEpoch)
	assert.Equal(t, qc.ViewNumber, blockViewNumber)
	assert.Equal(t, qc.BlockHash, blockHash)
	assert.Equal(t, qc.BlockNumber, blockNumber)
}

func TestViewChangeBySwitchPoint(t *testing.T) {
	pk, sk, cbftnodes := GenerateCbftNode(4)
	nodes := make([]*TestCBFT, 0)
	for i := 0; i < 4; i++ {
		node := MockNode(pk[i], sk[i], cbftnodes, 10000, 10)
		node.agency = validator.NewMockAgency(cbftnodes, 10)
		assert.Nil(t, node.Start())
		node.engine.validatorPool.MockSwitchPoint(10)
		nodes = append(nodes, node)
	}

	result := make(chan *types.Block, 1)
	complete := make(chan struct{}, 1)
	parent := nodes[0].chain.Genesis()
	for i := 0; i < 10; i++ {
		block := NewBlockWithSign(parent.Hash(), parent.NumberU64()+1, nodes[0])
		assert.True(t, nodes[0].engine.state.HighestExecutedBlock().Hash() == block.ParentHash())
		nodes[0].engine.OnSeal(block, result, nil, complete)
		<-complete

		_, qc := nodes[0].engine.blockTree.FindBlockAndQC(parent.Hash(), parent.NumberU64())
		select {
		case b := <-result:
			assert.NotNil(t, b)
			assert.Equal(t, uint32(i-1), nodes[0].engine.state.MaxQCIndex())
			pb := &protocols.PrepareBlock{
				Epoch:         nodes[0].engine.state.Epoch(),
				ViewNumber:    nodes[0].engine.state.ViewNumber(),
				Block:         b,
				BlockIndex:    uint32(i),
				ProposalIndex: uint32(0),
			}
			nodes[0].engine.signMsgByBls(pb)
			nodes[1].engine.OnPrepareBlock("id", pb)
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
				if j == 1 {
					nodes[1].engine.state.HadSendPrepareVote().Push(msg)
				}
				assert.Nil(t, nodes[j].engine.signMsgByBls(msg))
				nodes[0].engine.OnPrepareVote("id", msg)
				if i < 9 {
					assert.Nil(t, nodes[1].engine.OnPrepareVote("id", msg), fmt.Sprintf("number:%d", b.NumberU64()))
				}
			}
			parent = b
		}
	}
	// node-0 enough 10 block qc,change the epoch
	assert.Equal(t, uint64(2), nodes[0].engine.state.Epoch())
	assert.Equal(t, uint64(0), nodes[0].engine.state.ViewNumber())

	// node-1 change the view base lock block
	lockBlock := nodes[0].engine.state.HighestLockBlock()
	_, lockQC := nodes[0].engine.blockTree.FindBlockAndQC(lockBlock.Hash(), lockBlock.NumberU64())
	for i := 0; i < 4; i++ {
		epoch, view := nodes[1].engine.state.Epoch(), nodes[1].engine.state.ViewNumber()
		viewchange := &protocols.ViewChange{
			Epoch:          epoch,
			ViewNumber:     view,
			BlockHash:      lockBlock.Hash(),
			BlockNumber:    lockBlock.NumberU64(),
			ValidatorIndex: uint32(i),
			PrepareQC:      lockQC,
		}
		assert.Nil(t, nodes[i].engine.signMsgByBls(viewchange))
		assert.Nil(t, nodes[1].engine.OnViewChanges("id", &protocols.ViewChanges{
			VCs: []*protocols.ViewChange{
				viewchange,
			},
		}))
	}
	assert.NotNil(t, nodes[1].engine.state.LastViewChangeQC())
	assert.Equal(t, uint64(1), nodes[1].engine.state.ViewNumber())

	qcBlock := nodes[0].engine.state.HighestQCBlock()
	_, qcQC := nodes[0].engine.blockTree.FindBlockAndQC(qcBlock.Hash(), qcBlock.NumberU64())
	// change view by switchPoint
	nodes[1].engine.insertQCBlock(qcBlock, qcQC)
	assert.Equal(t, uint64(2), nodes[1].engine.state.Epoch())
	assert.Equal(t, uint64(0), nodes[1].engine.state.ViewNumber())
}
