package cbft

type safetyError interface {
	error
	Discard() bool //Is the error need discard
	Fetch() bool   //Is the error need fetch
}

type safetyRules interface {
	//Security rules for proposed blocks
	prepareBlockRules(block *prepareBlock) safetyError

	//Security rules for proposed votes
	prepareVoteRules(vote *prepareVote) safetyError

	//Security rules for viewChange
	viewChangeRules(vote *viewChange) safetyError
}
