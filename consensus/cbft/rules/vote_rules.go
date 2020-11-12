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
