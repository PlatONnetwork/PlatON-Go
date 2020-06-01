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

package state

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

const DefaultEpoch = 1
const DefaultViewNumber = 0

type PrepareVoteQueue struct {
	Votes []*protocols.PrepareVote `json:"votes"`
}

func newPrepareVoteQueue() *PrepareVoteQueue {
	return &PrepareVoteQueue{
		Votes: make([]*protocols.PrepareVote, 0),
	}
}

func (p *PrepareVoteQueue) Top() *protocols.PrepareVote {
	return p.Votes[0]
}

func (p *PrepareVoteQueue) Pop() *protocols.PrepareVote {
	v := p.Votes[0]
	p.Votes = p.Votes[1:]
	return v
}

func (p *PrepareVoteQueue) Push(vote *protocols.PrepareVote) {
	p.Votes = append(p.Votes, vote)
}

func (p *PrepareVoteQueue) Peek() []*protocols.PrepareVote {
	return p.Votes
}

func (p *PrepareVoteQueue) Empty() bool {
	return len(p.Votes) == 0
}

func (p *PrepareVoteQueue) Len() int {
	return len(p.Votes)
}

func (p *PrepareVoteQueue) reset() {
	p.Votes = make([]*protocols.PrepareVote, 0)
}

func (p *PrepareVoteQueue) Had(index uint32) bool {
	for _, p := range p.Votes {
		if p.BlockIndex == index {
			return true
		}
	}
	return false
}

type prepareVotes struct {
	Votes map[uint32]*protocols.PrepareVote `json:"votes"`
}

func newPrepareVotes() *prepareVotes {
	return &prepareVotes{
		Votes: make(map[uint32]*protocols.PrepareVote),
	}
}

func (p *prepareVotes) hadVote(vote *protocols.PrepareVote) bool {
	for _, v := range p.Votes {
		if v.MsgHash() == vote.MsgHash() {
			return true
		}
	}
	return false
}

func (p *prepareVotes) len() int {
	return len(p.Votes)
}

func (p *prepareVotes) clear() {
	p.Votes = make(map[uint32]*protocols.PrepareVote)
}

type viewBlocks struct {
	Blocks map[uint32]viewBlock `json:"blocks"`
}

func (v *viewBlocks) MarshalJSON() ([]byte, error) {
	type viewBlocks struct {
		Hash   common.Hash `json:"hash"`
		Number uint64      `json:"number"`
		Index  uint32      `json:"blockIndex"`
	}

	vv := make(map[uint32]viewBlocks)
	for index, block := range v.Blocks {
		vv[index] = viewBlocks{
			Hash:   block.hash(),
			Number: block.number(),
			Index:  block.blockIndex(),
		}
	}

	return json.Marshal(vv)
}

func newViewBlocks() *viewBlocks {
	return &viewBlocks{
		Blocks: make(map[uint32]viewBlock),
	}
}

func (v *viewBlocks) index(i uint32) viewBlock {
	return v.Blocks[i]
}

func (v *viewBlocks) addBlock(block viewBlock) {
	v.Blocks[block.blockIndex()] = block
}

func (v *viewBlocks) clear() {
	v.Blocks = make(map[uint32]viewBlock)
}

func (v *viewBlocks) len() int {
	return len(v.Blocks)
}

func (v *viewBlocks) MaxIndex() uint32 {
	max := uint32(math.MaxUint32)
	for _, b := range v.Blocks {
		if max == math.MaxUint32 || b.blockIndex() > max {
			max = b.blockIndex()
		}
	}
	return max
}

type viewQCs struct {
	MaxIndex uint32                        `json:"maxIndex"`
	QCs      map[uint32]*ctypes.QuorumCert `json:"qcs"`
}

func newViewQCs() *viewQCs {
	return &viewQCs{
		MaxIndex: math.MaxUint32,
		QCs:      make(map[uint32]*ctypes.QuorumCert),
	}
}

func (v *viewQCs) index(i uint32) *ctypes.QuorumCert {
	return v.QCs[i]
}

func (v *viewQCs) addQC(qc *ctypes.QuorumCert) {
	v.QCs[qc.BlockIndex] = qc
	if v.MaxIndex == math.MaxUint32 {
		v.MaxIndex = qc.BlockIndex
	}
	if v.MaxIndex < qc.BlockIndex {
		v.MaxIndex = qc.BlockIndex
	}
}

func (v *viewQCs) maxQCIndex() uint32 {
	return v.MaxIndex
}

func (v *viewQCs) clear() {
	v.QCs = make(map[uint32]*ctypes.QuorumCert)
	v.MaxIndex = math.MaxUint32
}

func (v *viewQCs) len() int {
	return len(v.QCs)
}

type viewVotes struct {
	Votes map[uint32]*prepareVotes `json:"votes"`
}

func newViewVotes() *viewVotes {
	return &viewVotes{
		Votes: make(map[uint32]*prepareVotes),
	}
}

func (v *viewVotes) addVote(id uint32, vote *protocols.PrepareVote) {
	if ps, ok := v.Votes[vote.BlockIndex]; ok {
		ps.Votes[id] = vote
	} else {
		ps := newPrepareVotes()
		ps.Votes[id] = vote
		v.Votes[vote.BlockIndex] = ps
	}
}

func (v *viewVotes) index(i uint32) *prepareVotes {
	return v.Votes[i]
}

func (v *viewVotes) MaxIndex() uint32 {
	max := uint32(math.MaxUint32)
	for index, _ := range v.Votes {
		if max == math.MaxUint32 || index > max {
			max = index
		}
	}
	return max
}

func (v *viewVotes) clear() {
	v.Votes = make(map[uint32]*prepareVotes)
}

type viewChanges struct {
	ViewChanges map[uint32]*protocols.ViewChange `json:"viewchanges"`
}

func newViewChanges() *viewChanges {
	return &viewChanges{
		ViewChanges: make(map[uint32]*protocols.ViewChange),
	}
}

func (v *viewChanges) addViewChange(id uint32, viewChange *protocols.ViewChange) {
	v.ViewChanges[id] = viewChange
}

func (v *viewChanges) len() int {
	return len(v.ViewChanges)
}

func (v *viewChanges) clear() {
	v.ViewChanges = make(map[uint32]*protocols.ViewChange)
}

type executing struct {
	// Block index of current view
	BlockIndex uint32 `json:"blockIndex"`
	// Whether to complete
	Finish bool `json:"finish"`
}

type view struct {
	epoch      uint64
	viewNumber uint64

	// The status of the block is currently being executed,
	// Finish indicates whether the execution is complete,
	// and the next block can be executed asynchronously after the execution is completed.
	executing executing

	// viewchange received by the current view
	viewChanges *viewChanges

	// QC of the previous view
	lastViewChangeQC *ctypes.ViewChangeQC

	// This view has been sent to other verifiers for voting
	hadSendPrepareVote *PrepareVoteQueue

	//Pending Votes of current view, parent block need receive N-f prepareVotes
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
	atomic.StoreUint64(&v.epoch, 0)
	atomic.StoreUint64(&v.viewNumber, 0)
	v.executing.BlockIndex = math.MaxUint32
	v.executing.Finish = false
	v.viewChanges.clear()
	v.hadSendPrepareVote.reset()
	v.pendingVote.reset()
	v.viewBlocks.clear()
	v.viewQCs.clear()
	v.viewVotes.clear()
}

func (v *view) ViewNumber() uint64 {
	return atomic.LoadUint64(&v.viewNumber)
}

func (v *view) Epoch() uint64 {
	return atomic.LoadUint64(&v.epoch)
}

func (v *view) MarshalJSON() ([]byte, error) {
	type view struct {
		Epoch              uint64               `json:"epoch"`
		ViewNumber         uint64               `json:"viewNumber"`
		Executing          executing            `json:"executing"`
		ViewChanges        *viewChanges         `json:"viewchange"`
		LastViewChangeQC   *ctypes.ViewChangeQC `json:"lastViewchange"`
		HadSendPrepareVote *PrepareVoteQueue    `json:"hadSendPrepareVote"`
		PendingVote        *PrepareVoteQueue    `json:"pendingPrepareVote"`
		ViewBlocks         *viewBlocks          `json:"viewBlocks"`
		ViewQCs            *viewQCs             `json:"viewQcs"`
		ViewVotes          *viewVotes           `json:"viewVotes"`
	}
	vv := &view{
		Epoch:              atomic.LoadUint64(&v.epoch),
		ViewNumber:         atomic.LoadUint64(&v.viewNumber),
		Executing:          v.executing,
		ViewChanges:        v.viewChanges,
		LastViewChangeQC:   v.lastViewChangeQC,
		HadSendPrepareVote: v.hadSendPrepareVote,
		PendingVote:        v.pendingVote,
		ViewBlocks:         v.viewBlocks,
		ViewQCs:            v.viewQCs,
		ViewVotes:          v.viewVotes,
	}

	return json.Marshal(vv)
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

	//Include ViewNumber, ViewChanges, prepareVote , proposal block of current view
	*view

	highestQCBlock     atomic.Value
	highestLockBlock   atomic.Value
	highestCommitBlock atomic.Value

	//Set the timer of the view time window
	viewTimer *viewTimer

	blockTree *ctypes.BlockTree
}

func NewViewState(period uint64, blockTree *ctypes.BlockTree) *ViewState {
	return &ViewState{
		view:      newView(),
		viewTimer: newViewTimer(period),
		blockTree: blockTree,
	}
}

func (vs *ViewState) ResetView(epoch uint64, viewNumber uint64) {
	vs.view.Reset()
	atomic.StoreUint64(&vs.view.epoch, epoch)
	atomic.StoreUint64(&vs.view.viewNumber, viewNumber)
}

func (vs *ViewState) Epoch() uint64 {
	return vs.view.Epoch()
}

func (vs *ViewState) ViewNumber() uint64 {
	return vs.view.ViewNumber()
}

func (vs *ViewState) ViewString() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d}", atomic.LoadUint64(&vs.view.epoch), atomic.LoadUint64(&vs.view.viewNumber))
}

func (vs *ViewState) Deadline() time.Time {
	return vs.viewTimer.deadline
}

func (vs *ViewState) NextViewBlockIndex() uint32 {
	return vs.viewBlocks.MaxIndex() + 1
}

func (vs *ViewState) MaxViewBlockIndex() uint32 {
	max := vs.viewBlocks.MaxIndex()
	if max == math.MaxUint32 {
		return 0
	}
	return max
}

func (vs *ViewState) MaxQCIndex() uint32 {
	return vs.view.viewQCs.maxQCIndex()
}

func (vs *ViewState) ViewVoteSize() int {
	return len(vs.viewVotes.Votes)
}

func (vs *ViewState) MaxViewVoteIndex() uint32 {
	max := vs.viewVotes.MaxIndex()
	if max == math.MaxUint32 {
		return 0
	}
	return max
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

func (vs *ViewState) ViewBlockSize() int {
	return len(vs.viewBlocks.Blocks)
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
		return ps.Votes
	}
	return nil
}

func (vs *ViewState) FindPrepareVote(blockIndex, validatorIndex uint32) *protocols.PrepareVote {
	ps := vs.viewVotes.index(blockIndex)
	if ps != nil {
		if v, ok := ps.Votes[validatorIndex]; ok {
			return v
		}
	}
	return nil
}

func (vs *ViewState) AllViewChange() map[uint32]*protocols.ViewChange {
	return vs.viewChanges.ViewChanges
}

// Returns the block index being executed, has it been completed
func (vs *ViewState) Executing() (uint32, bool) {
	return vs.view.executing.BlockIndex, vs.view.executing.Finish
}

func (vs *ViewState) SetLastViewChangeQC(qc *ctypes.ViewChangeQC) {
	vs.view.lastViewChangeQC = qc
}

func (vs *ViewState) LastViewChangeQC() *ctypes.ViewChangeQC {
	return vs.view.lastViewChangeQC
}

// Set Executing block status
func (vs *ViewState) SetExecuting(index uint32, finish bool) {
	vs.view.executing.BlockIndex, vs.view.executing.Finish = index, finish
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

func (vs *ViewState) AddViewChange(id uint32, viewChange *protocols.ViewChange) {
	vs.view.viewChanges.addViewChange(id, viewChange)
}

func (vs *ViewState) ViewChangeByIndex(index uint32) *protocols.ViewChange {
	return vs.view.viewChanges.ViewChanges[index]
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
	if vs.executing.BlockIndex == math.MaxUint32 || (vs.executing.BlockIndex == 0 && !vs.executing.Finish) {
		block := vs.HighestQCBlock()
		if vs.lastViewChangeQC != nil {
			_, _, _, _, hash, _ := vs.lastViewChangeQC.MaxBlock()
			// fixme insertQCBlock should also change the state of executing
			if b := vs.blockTree.FindBlockByHash(hash); b != nil {
				block = b
			}
		}
		return block
	}

	var block *types.Block
	if vs.executing.Finish {
		block = vs.viewBlocks.index(vs.executing.BlockIndex).block()
	} else {
		block = vs.viewBlocks.index(vs.executing.BlockIndex - 1).block()
	}
	return block
}

func (vs *ViewState) FindBlock(hash common.Hash, number uint64) *types.Block {
	for _, b := range vs.viewBlocks.Blocks {
		if b.hash() == hash && b.number() == number {
			return b.block()
		}
	}
	return nil
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

func (vs *ViewState) MarshalJSON() ([]byte, error) {
	type hashNumber struct {
		Hash   common.Hash `json:"hash"`
		Number uint64      `json:"number"`
	}
	type state struct {
		View               *view      `json:"view"`
		HighestQCBlock     hashNumber `json:"highestQCBlock"`
		HighestLockBlock   hashNumber `json:"highestLockBlock"`
		HighestCommitBlock hashNumber `json:"highestCommitBlock"`
	}

	s := &state{
		View:               vs.view,
		HighestQCBlock:     hashNumber{Hash: vs.HighestQCBlock().Hash(), Number: vs.HighestQCBlock().NumberU64()},
		HighestLockBlock:   hashNumber{Hash: vs.HighestLockBlock().Hash(), Number: vs.HighestLockBlock().NumberU64()},
		HighestCommitBlock: hashNumber{Hash: vs.HighestCommitBlock().Hash(), Number: vs.HighestCommitBlock().NumberU64()},
	}
	return json.Marshal(s)
}
