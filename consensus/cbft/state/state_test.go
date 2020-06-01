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

package state

import (
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/assert"
)

func TestNewViewState(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	viewState.ResetView(1, 1)

	assert.Equal(t, uint64(1), viewState.Epoch())
	assert.Equal(t, uint64(1), viewState.ViewNumber())
	assert.Equal(t, 0, viewState.ViewBlockSize())
	assert.Equal(t, uint32(0), viewState.NextViewBlockIndex())
	assert.Equal(t, uint32(math.MaxUint32), viewState.MaxQCIndex())
	assert.Equal(t, 0, viewState.ViewVoteSize())

	assert.Equal(t, uint64(1), viewState.view.Epoch())
	assert.Equal(t, uint64(1), viewState.view.ViewNumber())
	_, err := viewState.view.MarshalJSON()
	assert.Nil(t, err)

	assert.Equal(t, 0, viewState.HadSendPrepareVote().Len())
	assert.Equal(t, 0, viewState.PendingPrepareVote().Len())

	viewState.SetExecuting(uint32(1), true)
	_, finish := viewState.Executing()
	assert.True(t, finish)

	viewState.SetLastViewChangeQC(&ctypes.ViewChangeQC{})
	assert.NotNil(t, viewState.LastViewChangeQC())

	viewState.SetHighestCommitBlock(newBlock(1))
	viewState.SetHighestLockBlock(newBlock(2))
	viewState.SetHighestQCBlock(newBlock(3))
	assert.NotNil(t, viewState.HighestCommitBlock())
	assert.NotNil(t, viewState.HighestLockBlock())
	assert.NotNil(t, viewState.HighestQCBlock())
	assert.NotNil(t, viewState.HighestBlockString())

	_, err = viewState.MarshalJSON()
	assert.Nil(t, err)

	viewState.SetViewTimer(1)

	select {
	case <-viewState.ViewTimeout():
		assert.True(t, viewState.IsDeadline())
		assert.True(t, viewState.IsDeadline())
	}
}

func TestPrepareVoteQueue(t *testing.T) {
	queue := newPrepareVoteQueue()

	for i := 0; i < 10; i++ {
		b := &protocols.PrepareVote{BlockIndex: uint32(i)}
		queue.Push(b)
	}

	expect := uint32(0)
	for !queue.Empty() {
		assert.Equal(t, queue.Top().BlockIndex, expect)
		assert.True(t, queue.Had(expect))
		assert.False(t, queue.Had(expect+10))
		queue.Pop()
		expect++
	}

	assert.Equal(t, 0, queue.Len())
	assert.Len(t, queue.Peek(), 0)
}

func TestPrepareVotes(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)

	var b *protocols.PrepareVote
	for i := 0; i < 10; i++ {
		b = &protocols.PrepareVote{BlockIndex: uint32(i)}
		viewState.viewVotes.addVote(uint32(i), b)
	}

	assert.Equal(t, 1, viewState.PrepareVoteLenByIndex(uint32(0)))
	assert.NotNil(t, viewState.FindPrepareVote(uint32(1), uint32(1)))
	assert.True(t, viewState.viewVotes.Votes[uint32(9)].hadVote(b))

	viewState.viewVotes.clear()
}

func TestViewBlocks(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)

	var viewBlock *prepareViewBlock
	for i := 0; i < 10; i++ {
		viewBlock = &prepareViewBlock{
			pb: &protocols.PrepareBlock{BlockIndex: uint32(i), Block: newBlock(0)},
		}
		viewState.viewBlocks.addBlock(viewBlock)
	}
	assert.Equal(t, 10, viewState.viewBlocks.len())
	assert.Equal(t, uint32(9), viewState.viewBlocks.MaxIndex())
	assert.Equal(t, viewBlock.hash(), viewState.viewBlocks.Blocks[9].hash())
	assert.Equal(t, viewBlock.number(), viewState.viewBlocks.Blocks[9].number())

	assert.NotNil(t, viewState.ViewBlockByIndex(9))
	assert.NotNil(t, viewState.PrepareBlockByIndex(9))
	assert.Equal(t, 10, viewState.ViewBlockSize())
}

var (
	BaseMs = uint64(10000)
)

func TestViewVotes(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	votes := viewState.viewVotes
	prepareVotes := []*protocols.PrepareVote{
		{BlockIndex: uint32(5)},
		{BlockIndex: uint32(6)},
		{BlockIndex: uint32(7)},
	}

	for i, p := range prepareVotes {
		viewState.AddPrepareVote(uint32(i), p)
		votes.addVote(uint32(i), p)
	}
	assert.Equal(t, 3, len(viewState.viewVotes.Votes))
	assert.Equal(t, uint32(7), viewState.MaxViewVoteIndex())
	assert.Len(t, viewState.AllPrepareVoteByIndex(5), 1)
	assert.Equal(t, viewState.PrepareVoteLenByIndex(uint32(len(prepareVotes))), 0)
	assert.Len(t, viewState.AllPrepareVoteByIndex(uint32(len(prepareVotes))), 0)

	votes.clear()
	assert.Len(t, votes.Votes, 0)
}

func TestNewViewQC(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	viewQCs := viewState.viewQCs

	for i := uint32(0); i < 10; i++ {
		viewState.AddQC(&ctypes.QuorumCert{BlockIndex: i})
	}

	assert.Equal(t, viewState.MaxQCIndex(), uint32(9))
	for i := uint32(0); i < 10; i++ {
		assert.NotNil(t, viewQCs.index(i))
	}

	for i := uint32(10); i < 20; i++ {
		assert.Nil(t, viewQCs.index(i))
	}

	viewQCs.clear()
	assert.Equal(t, viewQCs.len(), 0)
}

func newBlock(number uint64) *types.Block {
	header := &types.Header{
		Number:     big.NewInt(int64(number)),
		ParentHash: common.Hash{},
		Time:       big.NewInt(time.Now().UnixNano()),
		Extra:      nil,
	}
	block := types.NewBlockWithHeader(header)
	return block
}

func TestNewViewBlock(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	for i := uint64(0); i < 10; i++ {
		viewState.AddQCBlock(newBlock(i), &ctypes.QuorumCert{BlockNumber: i, BlockIndex: uint32(i)})
	}

	for i := uint32(0); i < 10; i++ {
		block, _ := viewState.ViewBlockAndQC(i)
		assert.NotNil(t, block)
	}

	block, qc := viewState.ViewBlockAndQC(11)
	assert.Nil(t, block)
	assert.Nil(t, qc)
}

func TestNewViewChanges(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)

	var v *protocols.ViewChange
	for i := uint32(0); i < 10; i++ {
		v = &protocols.ViewChange{ValidatorIndex: i}
		viewState.AddViewChange(i, v)
	}

	assert.Equal(t, 10, viewState.ViewChangeLen())
	assert.Equal(t, 10, len(viewState.AllViewChange()))
	assert.Equal(t, uint32(9), viewState.ViewChangeByIndex(9).ValidatorIndex)

}
