package rules

import "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

type VoteError interface {
	error
	Discard() bool //Is the error need discard
}

type VoteRules interface {
	//Determine if the resulting vote is allowed to be sent
	allowVote(vote *protocols.PrepareVote) VoteError
}
