package rules

import "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

type SafetyError interface {
	error
	Discard() bool //Is the error need discard
	Fetch() bool   //Is the error need fetch
}

type SafetyRules interface {
	//Security rules for proposed blocks
	prepareBlockRules(block *protocols.PrepareBlock) SafetyError

	//Security rules for proposed votes
	prepareVoteRules(vote *protocols.PrepareVote) SafetyError

	//Security rules for viewChange
	viewChangeRules(vote *protocols.ViewChange) SafetyError
}
