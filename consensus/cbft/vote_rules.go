package cbft

type voteError interface {
	error
	Discard() bool //Is the error need discard
}

type voteRules interface {
	//Determine if the resulting vote is allowed to be sent
	allowVote(vote *prepareVote) voteError
}
