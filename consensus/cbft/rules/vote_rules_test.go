package rules

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestVote(t *testing.T) {
	viewState, blockTree := newEpochViewNumberState(Epoch, ViewNumber, 10)
	amount := uint32(10)
	rules := NewVoteRules(viewState)
	testVotes(t, viewState, blockTree, rules, amount)
}

func testVotes(t *testing.T, viewState *state.ViewState, blockTree *ctypes.BlockTree, rules VoteRules, amount uint32) {
	type testCase struct {
		err     error
		discard bool
		pb      *protocols.PrepareVote
	}

	newPrepareVote := func(epoch, viewNumber uint64, hash common.Hash, number uint64, index uint32) *protocols.PrepareVote {
		return &protocols.PrepareVote{
			Epoch:       epoch,
			ViewNumber:  viewNumber,
			BlockHash:   hash,
			BlockNumber: number,
			BlockIndex:  index,
		}
	}
	qcBlock := viewState.HighestQCBlock()

	viewState.HadSendPrepareVote().Push(newPrepareVote(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), amount-1))

	tests := []testCase{
		{nil, false, newPrepareVote(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), amount)},
		{newError(""), true, newPrepareVote(Epoch, ViewNumber+1, qcBlock.Hash(), qcBlock.NumberU64(), amount+1)},
	}
	for i, c := range tests {
		if c.err == nil {
			assert.Nil(t, rules.AllowVote(c.pb), fmt.Sprintf("case:%d failed", i))
		} else {
			err := rules.AllowVote(c.pb)
			assert.NotNil(t, err, fmt.Sprintf("case:%d failed", i))
			assert.Equal(t, c.discard, err.Discard(), fmt.Sprintf("case:%d failed", i))
		}
	}
}
