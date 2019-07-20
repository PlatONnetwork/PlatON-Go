package rules

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	"time"
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
}

type baseSafetyRules struct {
	viewState *state.ViewState
}

// PrepareBlock rules
// 1.Less than local viewNumber drop
// 2.Synchronization greater than local viewNumber
// 3.Lost more than the time window
func (r *baseSafetyRules) PrepareBlockRules(block *protocols.PrepareBlock) SafetyError {
	if r.viewState.ViewNumber() > block.ViewNumber {
		return newError(fmt.Sprintf("viewNumber too low(local:%d, msg:%d)", r.viewState.ViewNumber(), block.ViewNumber))
	}

	if r.viewState.ViewNumber() < block.ViewNumber {
		isQCChild := func() bool {
			return block.Block.ParentHash() == r.viewState.HighestQCBlock().Hash()
		}
		isLockChild := func() bool {
			return block.Block.ParentHash() == r.viewState.HighestLockBlock().Hash()
		}
		isFirstBlock := func() bool {
			return block.BlockIndex == 0
		}
		if isFirstBlock() && (isQCChild() || isLockChild()) {
			return newViewError("need change view")
		}

		return newFetchError(fmt.Sprintf("viewNumber higher then local(local:%d, msg:%d)", r.viewState.ViewNumber(), block.ViewNumber))
	}

	if r.viewState.IsDeadline() {
		return newError(fmt.Sprintf("view's deadline is expire(over:%d)", time.Since(r.viewState.Deadline())))
	}
	return nil
}

// PrepareVote rules
// 1.Less than local viewNumber drop
// 2.Synchronization greater than local viewNumber
// 3.Lost more than the time window
func (r *baseSafetyRules) PrepareVoteRules(vote *protocols.PrepareVote) SafetyError {
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

// ViewChange rules
// 1.Less than local viewNumber drop
// 2.Synchronization greater than local viewNumber
func (r *baseSafetyRules) ViewChangeRules(viewChange *protocols.ViewChange) SafetyError {

	if r.viewState.ViewNumber() > viewChange.ViewNumber {
		return newError(fmt.Sprintf("viewNumber too low(local:%d, msg:%d)", r.viewState.ViewNumber(), viewChange.ViewNumber))
	}

	if r.viewState.ViewNumber() < viewChange.ViewNumber {
		return newFetchError(fmt.Sprintf("viewNumber higher then local(local:%d, msg:%d)", r.viewState.ViewNumber(), viewChange.ViewNumber))
	}
	return nil
}

func NewSafetyRules(viewState *state.ViewState) SafetyRules {
	return &baseSafetyRules{
		viewState: viewState,
	}
}
