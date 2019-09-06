package state

import (
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/stretchr/testify/assert"
)

func TestNewViewState(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	viewState.ResetView(1, 1)
	viewState.SetViewTimer(1)

	select {
	case <-viewState.ViewTimeout():
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

var (
	BaseMs = uint64(10000)
)

func TestViewVotes(t *testing.T) {
	viewState := NewViewState(BaseMs, nil)
	votes := viewState.viewVotes
	prepareVotes := []*protocols.PrepareVote{
		&protocols.PrepareVote{BlockIndex: uint32(0)},
		&protocols.PrepareVote{BlockIndex: uint32(1)},
		&protocols.PrepareVote{BlockIndex: uint32(2)},
	}

	for i, p := range prepareVotes {
		votes.addVote(uint32(i), p)
	}
	assert.Len(t, viewState.AllPrepareVoteByIndex(0), 1)
	assert.Nil(t, votes.index(uint32(len(prepareVotes))))
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
