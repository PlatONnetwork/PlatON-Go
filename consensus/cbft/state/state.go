package state

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common/math"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type PrepareVoteQueue struct {
	votes []*protocols.PrepareVote
}

func newPrepareVoteQueue() *PrepareVoteQueue {
	return &PrepareVoteQueue{
		votes: make([]*protocols.PrepareVote, 0),
	}
}

func (p *PrepareVoteQueue) Top() *protocols.PrepareVote {
	return p.votes[0]
}

func (p *PrepareVoteQueue) Pop() *protocols.PrepareVote {
	v := p.votes[0]
	p.votes = p.votes[1:]
	return v
}

func (p *PrepareVoteQueue) Push(vote *protocols.PrepareVote) {
	p.votes = append(p.votes, vote)
}

func (p *PrepareVoteQueue) Peek() []*protocols.PrepareVote {
	return p.votes
}

func (p *PrepareVoteQueue) Empty() bool {
	return len(p.votes) == 0
}

func (p *PrepareVoteQueue) Len() int {
	return len(p.votes)
}

func (p *PrepareVoteQueue) reset() {
	p.votes = make([]*protocols.PrepareVote, 0)
}

func (p *PrepareVoteQueue) Had(index uint32) bool {
	for _, p := range p.votes {
		if p.BlockIndex == index {
			return true
		}
	}
	return false
}

type prepareVotes struct {
	votes map[uint32]*protocols.PrepareVote
}

func newPrepareVotes() *prepareVotes {
	return &prepareVotes{
		votes: make(map[uint32]*protocols.PrepareVote),
	}
}

func (p *prepareVotes) hadVote(vote *protocols.PrepareVote) bool {
	for _, v := range p.votes {
		if v.MsgHash() == vote.MsgHash() {
			return true
		}
	}
	return false
}

func (p *prepareVotes) len() int {
	return len(p.votes)
}

func (p *prepareVotes) clear() {
	p.votes = make(map[uint32]*protocols.PrepareVote)
}

type viewBlocks struct {
	blocks map[uint32]viewBlock
}

func newViewBlocks() *viewBlocks {
	return &viewBlocks{
		blocks: make(map[uint32]viewBlock),
	}
}

func (v *viewBlocks) index(i uint32) viewBlock {
	return v.blocks[i]
}

func (v *viewBlocks) addBlock(block viewBlock) {
	v.blocks[block.blockIndex()] = block
}

func (v *viewBlocks) clear() {
	v.blocks = make(map[uint32]viewBlock)
}

func (v *viewBlocks) len() int {
	return len(v.blocks)
}

func (v *viewBlocks) MaxIndex() uint32 {
	max := uint32(math.MaxUint32)
	for _, b := range v.blocks {
		if max == math.MaxUint32 || b.blockIndex() > max {
			max = b.blockIndex()
		}
	}
	return max
}

type viewQCs struct {
	maxIndex uint32
	qcs      map[uint32]*ctypes.QuorumCert
}

func newViewQCs() *viewQCs {
	return &viewQCs{
		maxIndex: math.MaxUint32,
		qcs:      make(map[uint32]*ctypes.QuorumCert),
	}
}

func (v *viewQCs) index(i uint32) *ctypes.QuorumCert {
	return v.qcs[i]
}

func (v *viewQCs) addQC(qc *ctypes.QuorumCert) {
	v.qcs[qc.BlockIndex] = qc
	if v.maxIndex == math.MaxUint32 {
		v.maxIndex = qc.BlockIndex
	}
	if v.maxIndex < qc.BlockIndex {
		v.maxIndex = qc.BlockIndex
	}
}

func (v *viewQCs) maxQCIndex() uint32 {
	return v.maxIndex
}

func (v *viewQCs) clear() {
	v.qcs = make(map[uint32]*ctypes.QuorumCert)
	v.maxIndex = math.MaxUint32
}

func (v *viewQCs) len() int {
	return len(v.qcs)
}

type viewVotes struct {
	votes map[uint32]*prepareVotes
}

func newViewVotes() *viewVotes {
	return &viewVotes{
		votes: make(map[uint32]*prepareVotes),
	}
}

func (v *viewVotes) addVote(id uint32, vote *protocols.PrepareVote) {
	if ps, ok := v.votes[vote.BlockIndex]; ok {
		ps.votes[id] = vote
	} else {
		ps := newPrepareVotes()
		ps.votes[id] = vote
		v.votes[vote.BlockIndex] = ps
	}
}
func (v *viewVotes) index(i uint32) *prepareVotes {
	return v.votes[i]
}

func (v *viewVotes) clear() {
	v.votes = make(map[uint32]*prepareVotes)
}

type viewChanges struct {
	viewChanges map[uint32]*protocols.ViewChange
}

func newViewChanges() *viewChanges {
	return &viewChanges{
		viewChanges: make(map[uint32]*protocols.ViewChange),
	}
}

func (v *viewChanges) addViewChange(id uint32, viewChange *protocols.ViewChange) {
	v.viewChanges[id] = viewChange
}

func (v *viewChanges) len() int {
	return len(v.viewChanges)
}

func (v *viewChanges) clear() {
	v.viewChanges = make(map[uint32]*protocols.ViewChange)
}

type executing struct {
	// Block index of current view
	blockIndex uint32
	// Whether to complete
	finish bool
}

type view struct {
	epoch      uint64
	viewNumber uint64

	// The status of the block is currently being executed,
	// finish indicates whether the execution is complete,
	// and the next block can be executed asynchronously after the execution is completed.
	executing executing

	// viewchange received by the current view
	viewChanges *viewChanges

	// QC of the previous view
	lastViewChangeQC *ctypes.ViewChangeQC

	// This view has been sent to other verifiers for voting
	hadSendPrepareVote *PrepareVoteQueue

	//Pending votes of current view, parent block need receive N-f prepareVotes
	pendingVote *PrepareVoteQueue

	//Current view of the proposed block by the proposer
	viewBlocks *viewBlocks

	viewQCs *viewQCs

	//The current view generated by the vote
	viewVotes *viewVotes
}

func newView() *view {
	return &view{
		executing:          executing{math.MaxUint32, false},
		viewChanges:        newViewChanges(),
		hadSendPrepareVote: newPrepareVoteQueue(),
		pendingVote:        newPrepareVoteQueue(),
		viewBlocks:         newViewBlocks(),
		viewQCs:            newViewQCs(),
		viewVotes:          newViewVotes(),
	}
}
func (v *view) Reset() {
	v.epoch = 0
	v.viewNumber = 0
	v.executing.blockIndex = math.MaxUint32
	v.executing.finish = false
	v.viewChanges.clear()
	v.hadSendPrepareVote.reset()
	v.pendingVote.reset()
	v.viewBlocks.clear()
	v.viewQCs.clear()
	v.viewVotes.clear()
}

func (v *view) ViewNumber() uint64 {
	return v.viewNumber
}

func (v *view) Epoch() uint64 {
	return v.epoch
}

//func (v *view) HadSendPrepareVote(vote *protocols.PrepareVote) bool {
//	return v.hadSendPrepareVote.hadVote(vote)
//}

//The block of current view, there two types, prepareBlock and block
type viewBlock interface {
	hash() common.Hash
	number() uint64
	blockIndex() uint32
	block() *types.Block
	//If prepareBlock is an implementation of viewBlock, return prepareBlock, otherwise nil
	prepareBlock() *protocols.PrepareBlock
}

type prepareViewBlock struct {
	pb *protocols.PrepareBlock
}

func (p prepareViewBlock) hash() common.Hash {
	return p.pb.Block.Hash()
}

func (p prepareViewBlock) number() uint64 {
	return p.pb.Block.NumberU64()
}
func (p prepareViewBlock) blockIndex() uint32 {
	return p.pb.BlockIndex
}
func (p prepareViewBlock) block() *types.Block {
	return p.pb.Block
}
func (p prepareViewBlock) prepareBlock() *protocols.PrepareBlock {
	return p.pb
}

type qcBlock struct {
	b  *types.Block
	qc *ctypes.QuorumCert
}

func (q qcBlock) hash() common.Hash {
	return q.b.Hash()
}

func (q qcBlock) number() uint64 {
	return q.b.NumberU64()
}
func (q qcBlock) blockIndex() uint32 {
	if q.qc == nil {
		return 0
	}
	return q.qc.BlockIndex
}

func (q qcBlock) block() *types.Block {
	return q.b
}

func (q qcBlock) prepareBlock() *protocols.PrepareBlock {
	return nil
}

type ViewState struct {

	//Include ViewNumber, viewChanges, prepareVote , proposal block of current view
	*view

	//Highest executed block height
	highestExecutedBlock atomic.Value

	highestQCBlock     atomic.Value
	highestLockBlock   atomic.Value
	highestCommitBlock atomic.Value

	//Set the timer of the view time window
	viewTimer *viewTimer
}

func NewViewState() *ViewState {
	return &ViewState{
		view:      newView(),
		viewTimer: newViewTimer(),
	}
}

func (vs *ViewState) ResetView(epoch uint64, viewNumber uint64) {
	vs.view.Reset()
	vs.view.epoch = epoch
	vs.view.viewNumber = viewNumber
}

func (vs *ViewState) Epoch() uint64 {
	return vs.view.epoch
}

func (vs *ViewState) ViewNumber() uint64 {
	return vs.view.viewNumber
}

func (vs *ViewState) ViewString() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d}", vs.view.epoch, vs.view.viewNumber)
}
func (vs *ViewState) Deadline() time.Time {
	return vs.viewTimer.deadline
}

func (vs *ViewState) NumViewBlocks() uint32 {
	return uint32(vs.viewBlocks.len())
}

func (vs *ViewState) NextViewBlockIndex() uint32 {
	return vs.viewBlocks.MaxIndex() + 1
}

func (vs *ViewState) MaxQCIndex() uint32 {
	return vs.view.viewQCs.maxQCIndex()
}

func (vs *ViewState) PrepareVoteLenByIndex(index uint32) int {
	ps := vs.viewVotes.index(index)
	if ps != nil {
		return ps.len()
	}
	return 0
}

// Find the block corresponding to the current view according to the index
func (vs *ViewState) ViewBlockByIndex(index uint32) *types.Block {
	if b := vs.view.viewBlocks.index(index); b != nil {
		return b.block()
	}
	return nil
}

func (vs *ViewState) PrepareBlockByIndex(index uint32) *protocols.PrepareBlock {
	if b := vs.view.viewBlocks.index(index); b != nil {
		return b.prepareBlock()
	}
	return nil
}

func (vs *ViewState) HadSendPrepareVote() *PrepareVoteQueue {
	return vs.view.hadSendPrepareVote
}

func (vs *ViewState) PendingPrepareVote() *PrepareVoteQueue {
	return vs.view.pendingVote
}

func (vs *ViewState) AllPrepareVoteByIndex(index uint32) map[uint32]*protocols.PrepareVote {
	ps := vs.viewVotes.index(index)
	if ps != nil {
		return ps.votes
	}
	return nil
}

func (vs *ViewState) AllViewChange() map[uint32]*protocols.ViewChange {
	return vs.viewChanges.viewChanges
}

// Returns the block index being executed, has it been completed
func (vs *ViewState) Executing() (uint32, bool) {
	return vs.view.executing.blockIndex, vs.view.executing.finish
}

func (vs *ViewState) SetLastViewChangeQC(qc *ctypes.ViewChangeQC) {
	vs.view.lastViewChangeQC = qc
}

func (vs *ViewState) LastViewChangeQC() *ctypes.ViewChangeQC {
	return vs.view.lastViewChangeQC
}

// Set Executing block status
func (vs *ViewState) SetExecuting(index uint32, finish bool) {
	vs.view.executing.blockIndex, vs.view.executing.finish = index, finish
}

func (vs *ViewState) ViewBlockAndQC(blockIndex uint32) (*types.Block, *ctypes.QuorumCert) {
	qc := vs.viewQCs.index(blockIndex)
	if b := vs.view.viewBlocks.index(blockIndex); b != nil {
		return b.block(), qc
	}
	return nil, qc
}

func (vs *ViewState) AddPrepareBlock(pb *protocols.PrepareBlock) {
	vs.view.viewBlocks.addBlock(&prepareViewBlock{pb})
}

func (vs *ViewState) AddQCBlock(block *types.Block, qc *ctypes.QuorumCert) {
	vs.view.viewBlocks.addBlock(&qcBlock{b: block, qc: qc})
}

func (vs *ViewState) AddQC(qc *ctypes.QuorumCert) {
	vs.view.viewQCs.addQC(qc)
}

func (vs *ViewState) AddPrepareVote(id uint32, vote *protocols.PrepareVote) {
	vs.view.viewVotes.addVote(id, vote)
}

func (vs *ViewState) AddViewChange(id uint32, vote *protocols.ViewChange) {
	vs.view.viewChanges.addViewChange(id, vote)
}

func (vs *ViewState) ViewChangeLen() int {
	return vs.view.viewChanges.len()
}

func (vs *ViewState) HighestBlockString() string {
	qc := vs.HighestQCBlock()
	lock := vs.HighestLockBlock()
	commit := vs.HighestCommitBlock()
	return fmt.Sprintf("{HighestQC:{hash:%s,number:%d},HighestLock:{hash:%s,number:%d},HighestCommit:{hash:%s,number:%d}}",
		qc.Hash().TerminalString(), qc.NumberU64(),
		lock.Hash().TerminalString(), lock.NumberU64(),
		commit.Hash().TerminalString(), commit.NumberU64())
}

func (vs *ViewState) HighestExecutedBlock() *types.Block {
	if (vs.executing.blockIndex == 0 && vs.executing.finish == false) ||
		vs.executing.blockIndex == math.MaxUint32 {
		return vs.HighestQCBlock()
	}

	var block *types.Block
	if vs.executing.finish {
		block = vs.viewBlocks.index(vs.executing.blockIndex).block()
	} else {
		block = vs.viewBlocks.index(vs.executing.blockIndex - 1).block()
	}
	return block
}

func (vs *ViewState) SetHighestQCBlock(ext *types.Block) {
	vs.highestQCBlock.Store(ext)
}

func (vs *ViewState) HighestQCBlock() *types.Block {
	if v := vs.highestQCBlock.Load(); v == nil {
		panic("Get highest qc block failed")
	} else {
		return v.(*types.Block)
	}
}

func (vs *ViewState) SetHighestLockBlock(ext *types.Block) {
	vs.highestLockBlock.Store(ext)
}

func (vs *ViewState) HighestLockBlock() *types.Block {
	if v := vs.highestLockBlock.Load(); v == nil {
		panic("Get highest lock block failed")
	} else {
		return v.(*types.Block)
	}
}

func (vs *ViewState) SetHighestCommitBlock(ext *types.Block) {
	vs.highestCommitBlock.Store(ext)
}

func (vs *ViewState) HighestCommitBlock() *types.Block {
	if v := vs.highestCommitBlock.Load(); v == nil {
		panic("Get highest commit block failed")
	} else {
		return v.(*types.Block)
	}
}

func (vs *ViewState) IsDeadline() bool {
	return vs.viewTimer.isDeadline()
}

func (vs *ViewState) ViewTimeout() <-chan time.Time {
	return vs.viewTimer.timerChan()
}

func (vs *ViewState) SetViewTimer(viewInterval uint64) {
	vs.viewTimer.setupTimer(viewInterval)
}

func (vs *ViewState) String() string {
	return fmt.Sprintf("")
}
