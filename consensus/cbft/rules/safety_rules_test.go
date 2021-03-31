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

package rules

import (
	"fmt"
	"math/big"
	"strings"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/stretchr/testify/assert"
)

const (
	// prepareBlock
	overIndexLimit       = "blockIndex higher than amount"
	existIndex           = "blockIndex already exists"
	firstBlockNotQCChild = "the first index block is not contiguous by local highestQC or highestLock"
	notExistPreIndex     = "previous index block not exists"
	diffPreIndexBlock    = "non contiguous index block"
	backwardPrepare      = "prepareBlock time is before parent"
	advancePrepare       = "prepareBlock time is advance"

	viewNumberTooLow = "viewNumber too low"
	needChangeView   = "need change view"
	needFetchBlock   = "epoch higher than local"

	// prepareVote
	overVoteIndexLimit = "voteIndex higher than amount"
	existVote          = "prepare vote has exist"
	noExistPrepare     = "current index block not existed"
	viewNumberTooHigh  = "viewNumber higher than local"
	alreadyQCBlock     = "current index block is already qc block"
	epochTooLow        = "epoch too low"
	epochTooHigh       = "new epoch, need fetch blocks"
)

func NewBlock(parent common.Hash, number uint64, blockTime *big.Int) *types.Block {
	header := &types.Header{
		Number:     big.NewInt(int64(number)),
		ParentHash: parent,
		Time:       big.NewInt(time.Now().UnixNano() / 1e6),
		Coinbase:   common.BytesToAddress(utils.Rand32Bytes(32)),
	}
	if blockTime != nil {
		header.Time = blockTime
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
	viewState := state.NewViewState(Period, nil)
	viewState.ResetView(epoch, viewNumber)
	viewState.SetViewTimer(2)

	parent := NewBlock(common.Hash{}, 0, nil)

	viewState.SetHighestQCBlock(parent)
	viewState.SetHighestLockBlock(parent)
	viewState.SetHighestCommitBlock(parent)

	blockTree := ctypes.NewBlockTree(parent, nil)

	for i := uint32(0); i < amount; i++ {
		block := NewBlock(parent.Hash(), parent.NumberU64()+1, nil)
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
		parent = block
	}
	return viewState, blockTree
}

type testCase struct {
	err          error
	common       bool
	fetch        bool
	newView      bool
	fetchPrepare bool
	pb           *protocols.PrepareBlock
	pv           *protocols.PrepareVote
	vc           *protocols.ViewChange
}

func invokePrepareBlockRules(t *testing.T, rules SafetyRules, tests []testCase) {
	for i, c := range tests {
		if c.err == nil {
			assert.Nil(t, rules.PrepareBlockRules(c.pb), fmt.Sprintf("case:%d failed", i))
		} else {
			err := rules.PrepareBlockRules(c.pb)
			assert.NotNil(t, err, fmt.Sprintf("case:%d failed", i), err)
			assert.Equal(t, c.common, err.Common(), fmt.Sprintf("case:%d failed %s", i, c.pb.String()), err)
			assert.Equal(t, c.fetch, err.Fetch(), fmt.Sprintf("case:%d failed %s", i, c.pb.String()), err)
			assert.Equal(t, c.newView, err.NewView(), fmt.Sprintf("case:%d failed %s", i, c.pb.String()), err)
			assert.Equal(t, c.fetchPrepare, err.FetchPrepare(), fmt.Sprintf("case:%d failed %s", i, c.pb.String()), err)
			assert.True(t, strings.HasPrefix(err.Error(), c.err.Error()))
		}
	}
}

func testBaseSafetyRulesPrepareBlockRules(t *testing.T, viewState *state.ViewState, blockTree *ctypes.BlockTree, rules SafetyRules, amount uint32) {
	newPrepareBlock := func(epoch, viewNumber uint64, parentHash common.Hash, blockNumber uint64, blockIndex uint32, viewChangeQC *ctypes.ViewChangeQC, blockTime *big.Int) *protocols.PrepareBlock {
		return &protocols.PrepareBlock{
			Epoch:        epoch,
			ViewNumber:   viewNumber,
			Block:        NewBlock(parentHash, blockNumber, blockTime),
			BlockIndex:   blockIndex,
			ViewChangeQC: viewChangeQC,
		}
	}
	newViewChangeQC := func(epoch, viewNumber uint64, blockHash common.Hash, blockNumber uint64) *ctypes.ViewChangeQC {
		return &ctypes.ViewChangeQC{
			QCs: []*ctypes.ViewChangeQuorumCert{
				{
					Epoch:       epoch,
					ViewNumber:  viewNumber,
					BlockHash:   blockHash,
					BlockNumber: blockNumber,
				},
			},
		}
	}

	qcBlock := viewState.HighestQCBlock()
	commitBlock := viewState.HighestCommitBlock()
	nextIndex := viewState.NextViewBlockIndex()

	tests := []testCase{
		{newCommonError(overIndexLimit), true, false, false, false, newPrepareBlock(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64()+1, amount+1, nil, nil), nil, nil},
		{newCommonError(existIndex), true, false, false, false, newPrepareBlock(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64()+1, 1, nil, nil), nil, nil},
		{newCommonError(notExistPreIndex), false, false, false, true, newPrepareBlock(Epoch, ViewNumber, common.BytesToHash(utils.Rand32Bytes(32)), qcBlock.NumberU64()+1, nextIndex+1, nil, nil), nil, nil},
		{newCommonError(diffPreIndexBlock), true, false, false, false, newPrepareBlock(Epoch, ViewNumber, common.BytesToHash(utils.Rand32Bytes(32)), qcBlock.NumberU64()+1, nextIndex, nil, nil), nil, nil},
		{newCommonError(backwardPrepare), true, false, false, false, newPrepareBlock(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64()+1, nextIndex, nil, new(big.Int).Sub(qcBlock.Time(), big.NewInt(1))), nil, nil},
		{newCommonError(advancePrepare), true, false, false, false, newPrepareBlock(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64()+1, nextIndex, nil, big.NewInt(common.Millis(time.Now().Add(riseTimeLimit*10000)))), nil, nil},
		{nil, false, false, false, false, newPrepareBlock(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), nextIndex, nil, nil), nil, nil},
	}
	invokePrepareBlockRules(t, rules, tests)
	// change the view
	viewState.ResetView(Epoch, ViewNumber+1)
	tests = []testCase{
		{newCommonError(firstBlockNotQCChild), true, false, false, false, newPrepareBlock(Epoch, ViewNumber+1, commitBlock.Hash(), commitBlock.NumberU64()+1, 0, nil, nil), nil, nil},
		{newCommonError(viewNumberTooLow), true, false, false, false, newPrepareBlock(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64()+1, 1, nil, nil), nil, nil},
		{newViewError(needChangeView), false, false, true, false, newPrepareBlock(Epoch, ViewNumber+2, qcBlock.Hash(), qcBlock.NumberU64()+1, 0, newViewChangeQC(Epoch, ViewNumber+1, qcBlock.Hash(), qcBlock.NumberU64()), nil), nil, nil},
		{newFetchError(needFetchBlock), false, true, false, false, newPrepareBlock(Epoch+1, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64()+1, 1, nil, nil), nil, nil},
	}
	invokePrepareBlockRules(t, rules, tests)
	// change the view
	viewState.ResetView(Epoch, ViewNumber)
}

func testBaseSafetyRulesPrepareVoteRules(t *testing.T, viewState *state.ViewState, blockTree *ctypes.BlockTree, rules SafetyRules, amount uint32) {
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
	pv := newPrepareVote(Epoch, ViewNumber, common.BytesToHash(utils.Rand32Bytes(32)), qcBlock.NumberU64()+1, amount-2)
	viewState.AddPrepareVote(0, pv)

	tests := []testCase{
		{newCommonError(overVoteIndexLimit), true, false, false, false, nil, newPrepareVote(Epoch, ViewNumber, common.BytesToHash(utils.Rand32Bytes(32)), qcBlock.NumberU64()+1, amount+1), nil},
		{newFetchError(existVote), true, false, false, false, nil, pv, nil},
		{newFetchError(noExistPrepare), false, false, false, true, nil, newPrepareVote(Epoch, ViewNumber, common.BytesToHash(utils.Rand32Bytes(32)), qcBlock.NumberU64()+2, amount-1), nil},
		{newCommonError(viewNumberTooLow), true, false, false, false, nil, newPrepareVote(Epoch, ViewNumber-1, common.BytesToHash(utils.Rand32Bytes(32)), qcBlock.NumberU64()+1, 0), nil},
		{newCommonError(viewNumberTooHigh), false, true, false, false, nil, newPrepareVote(Epoch, ViewNumber+1, common.BytesToHash(utils.Rand32Bytes(32)), qcBlock.NumberU64()+1, 0), nil},
		{newCommonError(alreadyQCBlock), true, false, false, false, nil, newPrepareVote(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), 0), nil},
	}
	for i, c := range tests {
		if c.err == nil {
			assert.Nil(t, rules.PrepareVoteRules(c.pv), fmt.Sprintf("case:%d failed", i))
		} else {
			err := rules.PrepareVoteRules(c.pv)
			assert.NotNil(t, err, fmt.Sprintf("case:%d failed", i), err)
			assert.Equal(t, c.common, err.Common(), fmt.Sprintf("case:%d failed %s", i, c.pv.String()), err)
			assert.Equal(t, c.fetch, err.Fetch(), fmt.Sprintf("case:%d failed %s", i, c.pv.String()), err)
			assert.Equal(t, c.newView, err.NewView(), fmt.Sprintf("case:%d failed %s", i, c.pv.String()), err)
			assert.Equal(t, c.fetchPrepare, err.FetchPrepare(), fmt.Sprintf("case:%d failed %s", i, c.pv.String()), err)
			assert.True(t, strings.HasPrefix(err.Error(), c.err.Error()))
		}
	}
}

func testBaseSafetyRulesViewChangeRules(t *testing.T, viewState *state.ViewState, blockTree *ctypes.BlockTree, rules SafetyRules, amount uint32) {
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
		{nil, false, false, false, false, nil, nil, newViewChange(Epoch, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), amount+1)},
		{newError(viewNumberTooHigh), false, true, false, false, nil, nil, newViewChange(Epoch, ViewNumber+1, qcBlock.Hash(), qcBlock.NumberU64(), amount+1)},
		{newError(epochTooHigh), false, true, false, false, nil, nil, newViewChange(Epoch+1, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), 1)},
		{newError(epochTooLow), true, false, false, false, nil, nil, newViewChange(Epoch-1, ViewNumber, qcBlock.Hash(), qcBlock.NumberU64(), 1)},
	}
	for i, c := range tests {
		if c.err == nil {
			assert.Nil(t, rules.ViewChangeRules(c.vc), fmt.Sprintf("case:%d failed", i))
		} else {
			err := rules.ViewChangeRules(c.vc)
			assert.NotNil(t, err, fmt.Sprintf("case:%d failed", i), err)
			assert.Equal(t, c.common, err.Common(), fmt.Sprintf("case:%d failed %s", i, c.vc.String()), err)
			assert.Equal(t, c.fetch, err.Fetch(), fmt.Sprintf("case:%d failed %s", i, c.vc.String()), err)
			assert.Equal(t, c.newView, err.NewView(), fmt.Sprintf("case:%d failed %s", i, c.vc.String()), err)
			assert.Equal(t, c.fetchPrepare, err.FetchPrepare(), fmt.Sprintf("case:%d failed %s", i, c.vc.String()), err)
			assert.True(t, strings.HasPrefix(err.Error(), c.err.Error()))
		}
	}
}

func TestSafetyError(t *testing.T) {
	viewState, blockTree := newEpochViewNumberState(Epoch, ViewNumber, 6)
	amount := uint32(10)
	rules := NewSafetyRules(viewState, blockTree, &ctypes.Config{Sys: &params.CbftConfig{Amount: amount}}, nil)
	testBaseSafetyRulesPrepareBlockRules(t, viewState, blockTree, rules, amount)
	testBaseSafetyRulesPrepareVoteRules(t, viewState, blockTree, rules, amount)
	testBaseSafetyRulesViewChangeRules(t, viewState, blockTree, rules, amount)
}
