package rules

import (
	"fmt"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type SafetyError interface {
	error
	Fetch() bool   //Is the error need fetch
	NewView() bool //need change view
}

type safetyError struct {
	text    string
	fetch   bool
	newView bool
}

func (s safetyError) Error() string {
	return s.text
}

func (s safetyError) Fetch() bool {
	return s.fetch
}
func (s safetyError) NewView() bool {
	return s.newView
}

//func newSafetyError(text string, fetch, newView bool) SafetyError {
//	return &safetyError{
//		text:    text,
//		fetch:   fetch,
//		newView: newView,
//	}
//}

func newFetchError(text string) SafetyError {
	return &safetyError{
		text:    text,
		fetch:   true,
		newView: false,
	}
}
func newViewError(text string) SafetyError {
	return &safetyError{
		text:    text,
		fetch:   false,
		newView: true,
	}
}

func newError(text string) SafetyError {
	return &safetyError{
		text:    text,
		fetch:   false,
		newView: false,
	}
}

type SafetyRules interface {
	// Security rules for proposed blocks
	PrepareBlockRules(block *protocols.PrepareBlock) SafetyError

	// Security rules for proposed votes
	PrepareVoteRules(vote *protocols.PrepareVote) SafetyError

	// Security rules for viewChange
	ViewChangeRules(vote *protocols.ViewChange) SafetyError

	// Security rules for qcblock
	QCBlockRules(block *types.Block, qc *ctypes.QuorumCert) SafetyError
}

type baseSafetyRules struct {
	viewState *state.ViewState
	blockTree *ctypes.BlockTree
	config    *ctypes.Config
}

// PrepareBlock rules
// 1.Less than local viewNumber drop
// 2.Synchronization greater than local viewNumber
// 3.Lost more than the time window
func (r *baseSafetyRules) PrepareBlockRules(block *protocols.PrepareBlock) SafetyError {
	if r.viewState.Epoch() != block.Epoch {
		return r.changeEpochBlockRules(block)
	}
	if r.viewState.ViewNumber() > block.ViewNumber {
		return newError(fmt.Sprintf("viewNumber too low(local:%d, msg:%d)", r.viewState.ViewNumber(), block.ViewNumber))
	}

	if r.viewState.ViewNumber() < block.ViewNumber {
		isQCChild := func() bool {
			return block.Block.NumberU64() == r.viewState.HighestQCBlock().NumberU64()+1 &&
				r.blockTree.FindBlockByHash(block.Block.ParentHash()) != nil
		}
		isLockChild := func() bool {
			return block.Block.ParentHash() == r.viewState.HighestLockBlock().Hash()
		}
		isFirstBlock := func() bool {
			return block.BlockIndex == 0
		}
		isNextView := func() bool {
			return r.viewState.ViewNumber()+1 == block.ViewNumber
		}

		acceptViewChangeQC := func() bool {
			if block.ViewChangeQC == nil {
				return r.config.Sys.Amount == r.viewState.MaxQCIndex()+1
			} else {
				_, _, hash, number := block.ViewChangeQC.MaxBlock()
				return number+1 == block.Block.NumberU64() && block.Block.ParentHash() == hash
			}
		}
		if isNextView() && isFirstBlock() && (isQCChild() || isLockChild()) && acceptViewChangeQC() {
			return newViewError("need change view")
		}

		return newFetchError(fmt.Sprintf("viewNumber higher then local(local:%d, msg:%d)", r.viewState.ViewNumber(), block.ViewNumber))
	}

	if r.viewState.IsDeadline() {
		return newError(fmt.Sprintf("view's deadline is expire(over:%s)", time.Since(r.viewState.Deadline())))
	}
	return nil
}

func (r *baseSafetyRules) changeEpochBlockRules(block *protocols.PrepareBlock) SafetyError {
	if r.viewState.Epoch() > block.Epoch {
		return newError(fmt.Sprintf("epoch too low(local:%d, msg:%d)", r.viewState.Epoch(), block.Epoch))
	}
	if block.Block.ParentHash() != r.viewState.HighestQCBlock().Hash() {
		return newFetchError(fmt.Sprintf("epoch higher then local(local:%d, msg:%d)", r.viewState.Epoch(), block.Epoch))
	}
	return newViewError("new epoch, need change view")
}

// PrepareVote rules
// 1.Less than local viewNumber drop
// 2.Synchronization greater than local viewNumber
// 3.Lost more than the time window
func (r *baseSafetyRules) PrepareVoteRules(vote *protocols.PrepareVote) SafetyError {
	if r.viewState.Epoch() != vote.Epoch {
		return r.changeEpochVoteRules(vote)
	}
	if r.viewState.ViewNumber() > vote.ViewNumber {
		return newError(fmt.Sprintf("viewNumber too low(local:%d, msg:%d)", r.viewState.ViewNumber(), vote.ViewNumber))
	}

	if r.viewState.ViewNumber() < vote.ViewNumber {
		return newFetchError(fmt.Sprintf("viewNumber higher then local(local:%d, msg:%d)", r.viewState.ViewNumber(), vote.ViewNumber))
	}

	if r.viewState.IsDeadline() {
		return newError(fmt.Sprintf("view's deadline is expire(over:%d)", time.Since(r.viewState.Deadline())))
	}
	return nil
}

func (r *baseSafetyRules) changeEpochVoteRules(vote *protocols.PrepareVote) SafetyError {
	if r.viewState.Epoch() > vote.Epoch {
		return newError(fmt.Sprintf("epoch too low(local:%d, msg:%d)", r.viewState.Epoch(), vote.Epoch))
	}

	return newFetchError("new epoch, need fetch blocks")
}

// ViewChange rules
// 1.Less than local viewNumber drop
// 2.Synchronization greater than local viewNumber
func (r *baseSafetyRules) ViewChangeRules(viewChange *protocols.ViewChange) SafetyError {

	if r.viewState.Epoch() != viewChange.Epoch {
		return r.changeEpochViewChangeRules(viewChange)
	}
	if r.viewState.ViewNumber() > viewChange.ViewNumber {
		return newError(fmt.Sprintf("viewNumber too low(local:%d, msg:%d)", r.viewState.ViewNumber(), viewChange.ViewNumber))
	}

	if r.viewState.ViewNumber() < viewChange.ViewNumber {
		return newFetchError(fmt.Sprintf("viewNumber higher then local(local:%d, msg:%d)", r.viewState.ViewNumber(), viewChange.ViewNumber))
	}
	return nil
}

func (r *baseSafetyRules) changeEpochViewChangeRules(viewChange *protocols.ViewChange) SafetyError {
	if r.viewState.Epoch() > viewChange.Epoch {
		return newError(fmt.Sprintf("epoch too low(local:%d, msg:%d)", r.viewState.Epoch(), viewChange.Epoch))
	}

	return newFetchError("new epoch, need fetch blocks")
}

func (r *baseSafetyRules) QCBlockRules(block *types.Block, qc *ctypes.QuorumCert) SafetyError {
	if r.viewState.Epoch() > qc.Epoch || r.viewState.ViewNumber() > qc.ViewNumber {
		return newError(fmt.Sprintf("epoch or viewNumber too low(local:%s, msg:{Epoch:%d,ViewNumber:%d})", r.viewState.ViewString(), qc.Epoch, qc.ViewNumber))
	}

	if b := r.blockTree.FindBlockByHash(qc.BlockHash); b == nil {
		return newError(fmt.Sprintf("not find parent qc block"))
	}
	if r.viewState.Epoch() > qc.Epoch || r.viewState.ViewNumber() > qc.ViewNumber {
		return newViewError("need change view")
	}
	return nil
}

func NewSafetyRules(viewState *state.ViewState, blockTree *ctypes.BlockTree, config *ctypes.Config) SafetyRules {
	return &baseSafetyRules{
		viewState: viewState,
		blockTree: blockTree,
		config:    config,
	}
}
