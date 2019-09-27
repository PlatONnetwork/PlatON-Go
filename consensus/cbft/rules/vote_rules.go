package rules

import (
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"
)

type VoteError interface {
	error
	Discard() bool //Is the error need discard
}

type voteError struct {
	s       string
	discard bool
}

func (v *voteError) Error() string {
	return v.s
}

func (v *voteError) Discard() bool {
	return v.discard
}

func newVoteError(s string, discard bool) VoteError {
	return &voteError{
		s:       s,
		discard: discard,
	}
}

type VoteRules interface {
	// Determine if the resulting vote is allowed to be sent
	AllowVote(vote *protocols.PrepareVote) VoteError
}

type baseVoteRules struct {
	viewState *state.ViewState
}

// Determine if voting is possible
// viewNumber should be equal
// have you voted
func (v *baseVoteRules) AllowVote(vote *protocols.PrepareVote) VoteError {
	if v.viewState.ViewNumber() != vote.ViewNumber {
		return newVoteError("", true)
	}

	if v.viewState.HadSendPrepareVote().Had(vote.BlockIndex) {
		return newVoteError("", true)
	}
	return nil
}

func NewVoteRules(viewState *state.ViewState) VoteRules {
	return &baseVoteRules{
		viewState: viewState,
	}

}
