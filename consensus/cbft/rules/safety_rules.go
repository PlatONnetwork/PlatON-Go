package rules

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
	"time"
)

type SafetyError interface {
	error
	Fetch() bool //Is the error need fetch
}

type safetyError struct {
	text  string
	fetch bool
}

func (s safetyError) Error() string {
	return s.text
}

func (s safetyError) Fetch() bool {
	return s.fetch
}

func newSafetyError(text string, fetch bool) SafetyError {
	return &safetyError{
		text:  text,
		fetch: fetch,
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
		return newSafetyError(fmt.Sprintf("viewNumber too low(local:%d, msg:%d)", r.viewState.ViewNumber(), block.ViewNumber), false)
	}

	if r.viewState.ViewNumber() < block.ViewNumber {
		return newSafetyError(fmt.Sprintf("viewNumber higher then local(local:%d, msg:%d)", r.viewState.ViewNumber(), block.ViewNumber), true)
	}

	if r.viewState.IsDeadline() {
		return newSafetyError(fmt.Sprintf("view's deadline is expire(over:%d)", time.Since(r.viewState.Deadline())), false)
	}
	return nil
}

// PrepareVote rules
// 1.Less than local viewNumber drop
// 2.Synchronization greater than local viewNumber
// 3.Lost more than the time window
func (r *baseSafetyRules) PrepareVoteRules(vote *protocols.PrepareVote) SafetyError {
	if r.viewState.ViewNumber() > vote.ViewNumber {
		return newSafetyError(fmt.Sprintf("viewNumber too low(local:%d, msg:%d)", r.viewState.ViewNumber(), vote.ViewNumber), false)
	}

	if r.viewState.ViewNumber() < vote.ViewNumber {
		return newSafetyError(fmt.Sprintf("viewNumber higher then local(local:%d, msg:%d)", r.viewState.ViewNumber(), vote.ViewNumber), true)
	}

	if r.viewState.IsDeadline() {
		return newSafetyError(fmt.Sprintf("view's deadline is expire(over:%d)", time.Since(r.viewState.Deadline())), false)
	}
	return nil

}

// ViewChange rules
// 1.Less than local viewNumber drop
// 2.Synchronization greater than local viewNumber
func (r *baseSafetyRules) ViewChangeRules(viewChange *protocols.ViewChange) SafetyError {

	if r.viewState.ViewNumber() > viewChange.ViewNumber {
		return newSafetyError(fmt.Sprintf("viewNumber too low(local:%d, msg:%d)", r.viewState.ViewNumber(), viewChange.ViewNumber), false)
	}

	if r.viewState.ViewNumber() < viewChange.ViewNumber {
		return newSafetyError(fmt.Sprintf("viewNumber higher then local(local:%d, msg:%d)", r.viewState.ViewNumber(), viewChange.ViewNumber), true)
	}
	return nil
}

func NewSafetyRules(viewState *state.ViewState) SafetyRules {
	return &baseSafetyRules{
		viewState: viewState,
	}
}
