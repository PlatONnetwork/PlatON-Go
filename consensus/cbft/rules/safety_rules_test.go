package rules

import (
	"fmt"
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/stretchr/testify/assert"
)

func NewBlock(parent common.Hash, number uint64) *types.Block {
	header := &types.Header{
		Number:     big.NewInt(int64(number)),
		ParentHash: parent,
		Time:       big.NewInt(time.Now().UnixNano()),
	}
	block := types.NewBlockWithHeader(header)
	return block
}

const (
	Epoch      = uint64(1)
	ViewNumber = uint64(1)
	Period     = uint64(10000)
)

func prepareQC(epoch, viewNumber uint64, hash common.Hash, number uint64, index uint32) *ctypes.QuorumCert {
	return &ctypes.QuorumCert{
		Epoch:       epoch,
		ViewNumber:  viewNumber,
		BlockHash:   hash,
		BlockNumber: number,
		BlockIndex:  index,
	}
}

func newEpochViewNumberState(epoch, viewNumber uint64, amount uint32) (*state.ViewState, *ctypes.BlockTree) {
	viewState := state.NewViewState(Period)
	viewState.ResetView(epoch, viewNumber)
	viewState.SetViewTimer(2)

	parent := NewBlock(common.Hash{}, 0)

	viewState.SetHighestQCBlock(parent)
	viewState.SetHighestLockBlock(parent)
	viewState.SetHighestCommitBlock(parent)

	blockTree := ctypes.NewBlockTree(parent, nil)

	for i := uint32(0); i < amount; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1)
		qc := prepareQC(epoch, epoch, block.Hash(), block.NumberU64(), i)
		viewState.AddQCBlock(block, qc)
		viewState.AddQC(qc)
		blockTree.InsertQCBlock(block, qc)

		if b := viewState.HighestLockBlock(); b != nil {
			viewState.SetHighestCommitBlock(b)
		}
		if b := viewState.HighestQCBlock(); b != nil {
			viewState.SetHighestLockBlock(b)
		}
		viewState.SetHighestQCBlock(block)
	}
	return viewState, blockTree
}

func testBaseSafetyRulesPrepareBlockRules(t *testing.T, viewState *state.ViewState, blockTree *ctypes.BlockTree, rules SafetyRules, amount uint32) {
	type testCase struct {
		err     error
		fetch   bool
		newView bool
		pb      *protocols.PrepareBlock
	}

	newPrepareBlock := func(epoch, viewNumber uint64, block *types.Block, blockIndex uint32) *protocols.PrepareBlock {
		return &protocols.PrepareBlock{
			Epoch:      epoch,
			ViewNumber: viewNumber,
			Block:      block,
			BlockIndex: blockIndex,
		}
	}

	qcBlock := viewState.HighestQCBlock()
	tests := []testCase{
		{nil, false, false, newPrepareBlock(Epoch, ViewNumber, qcBlock, 1)},
		{newError(""), false, false, newPrepareBlock(Epoch, ViewNumber, qcBlock, amount+1)},
		{newError(""), true, false, newPrepareBlock(Epoch+1, ViewNumber, qcBlock, 1)},
		{newError(""), true, false, newPrepareBlock(Epoch, ViewNumber+1, qcBlock, 1)},
		{newError(""), true, false, newPrepareBlock(Epoch, ViewNumber+1, NewBlock(qcBlock.Hash(), qcBlock.NumberU64()), 0)},
		{newError(""), false, true, newPrepareBlock(Epoch+1, ViewNumber, NewBlock(qcBlock.Hash(), qcBlock.NumberU64()+1), 0)},
	}
	for i, c := range tests {
		if c.err == nil {
			assert.Nil(t, rules.PrepareBlockRules(c.pb), fmt.Sprintf("case:%d failed", i))
		} else {
			err := rules.PrepareBlockRules(c.pb)
			assert.NotNil(t, err, fmt.Sprintf("case:%d failed", i), err)
			assert.Equal(t, c.fetch, err.Fetch(), fmt.Sprintf("case:%d failed %s", i, c.pb.String()), err)
			assert.Equal(t, c.newView, err.NewView(), fmt.Sprintf("case:%d failed", i), err)
		}
	}
}

func testBaseSafetyRulesPrepareVoteRules(t *testing.T, viewState *state.ViewState, blockTree *ctypes.BlockTree, rules SafetyRules, amount uint32) {
	type testCase struct {
		err     error
		fetch   bool
		newView bool
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

	tests := []testCase{
		{nil, false, false, newPrepareVote(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), amount+1)},
		{newError(""), true, false, newPrepareVote(Epoch, ViewNumber+1, qcBlock.Hash(), qcBlock.NumberU64(), amount+1)},
		{newError(""), true, false, newPrepareVote(Epoch+1, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), 1)},
		{newError(""), false, false, newPrepareVote(Epoch-1, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), 1)},
	}
	for i, c := range tests {
		if c.err == nil {
			assert.Nil(t, rules.PrepareVoteRules(c.pb), fmt.Sprintf("case:%d failed", i))
		} else {
			err := rules.PrepareVoteRules(c.pb)
			assert.NotNil(t, err, fmt.Sprintf("case:%d failed", i))
			assert.Equal(t, c.fetch, err.Fetch(), fmt.Sprintf("case:%d failed", i))
			assert.Equal(t, c.newView, err.NewView(), fmt.Sprintf("case:%d failed", i))
		}
	}
}

func testBaseSafetyRulesViewChangeRules(t *testing.T, viewState *state.ViewState, blockTree *ctypes.BlockTree, rules SafetyRules, amount uint32) {
	type testCase struct {
		err     error
		fetch   bool
		newView bool
		pb      *protocols.ViewChange
	}
	newViewChange := func(epoch, viewNumber uint64, hash common.Hash, number uint64, index uint32) *protocols.ViewChange {
		return &protocols.ViewChange{
			Epoch:       epoch,
			ViewNumber:  viewNumber,
			BlockHash:   hash,
			BlockNumber: number,
		}
	}

	qcBlock := viewState.HighestQCBlock()

	tests := []testCase{
		{nil, false, false, newViewChange(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), amount+1)},
		{newError(""), true, false, newViewChange(Epoch, ViewNumber+1, qcBlock.Hash(), qcBlock.NumberU64(), amount+1)},
		{newError(""), true, false, newViewChange(Epoch+1, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), 1)},
		{newError(""), false, false, newViewChange(Epoch-1, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), 1)},
	}
	for i, c := range tests {
		if c.err == nil {
			assert.Nil(t, rules.ViewChangeRules(c.pb), fmt.Sprintf("case:%d failed", i))
		} else {
			err := rules.ViewChangeRules(c.pb)
			assert.NotNil(t, err, fmt.Sprintf("case:%d failed", i))
			assert.NotNil(t, c.newView, err.NewView(), fmt.Sprintf("case:%d failed", i))
			assert.NotNil(t, c.fetch, err.Fetch(), fmt.Sprintf("case:%d failed", i))
		}
	}
}

func TestSafetyError(t *testing.T) {
	viewState, blockTree := newEpochViewNumberState(Epoch, ViewNumber, 10)
	amount := uint32(10)
	rules := NewSafetyRules(viewState, blockTree, &ctypes.Config{Sys: &params.CbftConfig{Amount: amount}})
	testBaseSafetyRulesPrepareBlockRules(t, viewState, blockTree, rules, amount)
	testBaseSafetyRulesPrepareVoteRules(t, viewState, blockTree, rules, amount)

	testBaseSafetyRulesViewChangeRules(t, viewState, blockTree, rules, amount)
}
