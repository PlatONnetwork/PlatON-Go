package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"sync/atomic"
)

type view struct {
	epoch      uint64
	viewNumber uint64

	//viewchange received by the current view
	viewChanges viewChanges

	//This view has been sent to other verifiers for voting
	hadSendPrepareVote prepareVotes

	//Pending votes of current view, parent block need receive N-f prepareVotes
	pendingVote prepareVotes

	//Current view of the proposed block by the proposer
	viewBlocks viewBlocks
}

//The block of current view, there two types, prepareBlock and block
type viewBlock interface {
	hash() common.Hash
	number() uint64
	blockIndex() uint32
	//If prepareBlock is an implementation of viewBlock, return prepareBlock, otherwise nil
	prepareBlock() *prepareBlock
}

type viewState struct {

	//Include ViewNumber, viewChanges, prepareVote , proposal block of current view
	currentView view

	highestQCBlock     atomic.Value
	highestLockBlock   atomic.Value
	highestCommitBlock atomic.Value

	//Set the timer of the view time window
	viewTimer viewTimer
}
