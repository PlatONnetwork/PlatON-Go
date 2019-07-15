package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"sync/atomic"
	"time"
)

type Signature struct {
}

type quorumCert struct {
	ViewNumber  uint64      `json:"view_number"`
	BlockHash   common.Hash `json:"block_hash"`
	BlockNumber uint64      `json:"block_number"`
	Signature   Signature   `json:"signature"`
}

type blockTree struct {
	root   *blockTree
	blocks map[uint64]map[common.Hash]*blockTree
}

type prepareVotes struct {
	votes []*prepareVote
}

type viewBlocks struct {
	blocks map[uint32]*viewBlock
}

type blockExt struct {
	viewNumber   uint64
	block        *types.Block
	inTree       bool
	executing    bool
	isExecuted   bool
	isSigned     bool
	isConfirmed  bool
	rcvTime      int64
	prepareVotes *prepareVoteSet //all prepareVotes for block
	prepareBlock *prepareBlock
	parent       *blockExt
	children     map[common.Hash]*blockExt
}

type prepareVoteSet struct {
	votes map[uint32]*prepareVote
}

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

type voteError interface {
	error
	Discard() bool //Is the error need discard
}

type voteRules interface {
	//Determine if the resulting vote is allowed to be sent
	allowVote(vote *prepareVote) voteError
}

//The block of current view, there two types, prepareBlock and block
type viewBlock interface {
	hash() common.Hash
	number() uint64
	blockIndex() uint32
	//If prepareBlock is an implementation of viewBlock, return prepareBlock, otherwise nil
	prepareBlock() *prepareBlock
}

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

type viewState struct {

	//Include ViewNumber, viewChanges, prepareVote , proposal block of current view
	currentView view

	highestQCBlock     atomic.Value
	highestLockBlock   atomic.Value
	highestCommitBlock atomic.Value

	//Set the timer of the view time window
	viewTimer viewTimer
}

type Sync interface {
}

type viewTimer struct {
	//Timer last timeout
	deadline time.Time
	timer    *time.Timer

	//Time window length calculation module
	timeInterval viewTimeInterval
}

func (t viewTimer) setupTimer() {

}

// Calculate the time window of each view，time=b*e^m
type viewTimeInterval struct {
	baseMs       uint64
	exponentBase float64
	maxExponent  uint64
}

func (vt viewTimeInterval) getViewTimeInterval(viewInterval uint64) time.Duration {
	return 0
}

type viewChanges struct {
	viewChanges map[common.Address]*viewChange
}

type blockExecutor interface {
	//Execution block, you need to pass in the parent block to find the parent block state
	execute(block *types.Block, parent *types.Block) error
}

//Block execution results, including block hash, block number, error message
type blockExecuteStatus struct {
	hash   common.Hash
	number uint64
	err    error
}

type asyncBlockExecutor interface {
	blockExecutor
	//Asynchronous acquisition block execution results
	executeStatus() chan blockExecuteStatus
}

type EvidencePool interface {
	//Deserialization of evidence
	UnmarshalEvidence([]byte) (consensus.Evidence, error)
	//Get current evidences
	Evidences() []consensus.Evidence
	//Clear all evidences
	Clear()
	Close()
}
