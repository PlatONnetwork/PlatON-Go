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

package cbft

import (
	"errors"
	"fmt"
	"reflect"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/state"

	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/node"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var (
	errNonContiguous = errors.New("non contiguous chain block state")
)

var (
	viewChangeQCPrefix = []byte("qc") // viewChangeQCPrefix + epoch (uint64 big endian) + viewNumber (uint64 big endian) -> viewChangeQC
	epochPrefix        = []byte("e")
)

// Bridge encapsulates functions required to update consensus state and consensus msg.
// As a bridge layer for cbft and wal.
type Bridge interface {
	UpdateChainState(qcState, lockState, commitState *protocols.State)
	ConfirmViewChange(epoch, viewNumber uint64, block *types.Block, qc *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC, preEpoch, preViewNumber uint64)
	SendViewChange(view *protocols.ViewChange)
	SendPrepareBlock(pb *protocols.PrepareBlock)
	SendPrepareVote(block *types.Block, vote *protocols.PrepareVote)
	GetViewChangeQC(epoch uint64, viewNumber uint64) (*ctypes.ViewChangeQC, error)

	Close()
}

// emptyBridge is a empty implementation for Bridge
type emptyBridge struct {
}

func (b *emptyBridge) UpdateChainState(qcState, lockState, commitState *protocols.State) {
}

func (b *emptyBridge) ConfirmViewChange(epoch, viewNumber uint64, block *types.Block, qc *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC, preEpoch, preViewNumber uint64) {
}

func (b *emptyBridge) SendViewChange(view *protocols.ViewChange) {
}

func (b *emptyBridge) SendPrepareBlock(pb *protocols.PrepareBlock) {
}

func (b *emptyBridge) SendPrepareVote(block *types.Block, vote *protocols.PrepareVote) {
}

func (b *emptyBridge) GetViewChangeQC(epoch uint64, viewNumber uint64) (*ctypes.ViewChangeQC, error) {
	return nil, nil
}

func (b *emptyBridge) Close() {

}

// baseBridge is a default implementation for Bridge
type baseBridge struct {
	cbft *Cbft
}

// NewBridge creates a new Bridge to update consensus state and consensus msg.
func NewBridge(ctx *node.ServiceContext, cbft *Cbft) (Bridge, error) {
	if ctx == nil {
		return &emptyBridge{}, nil
	}
	baseBridge := &baseBridge{
		cbft: cbft,
	}
	return baseBridge, nil
}

// UpdateChainState tries to update consensus state to wal
// If the write fails, the process will stop
// lockChainState or commitChainState may be nil, if it is nil, we only append qc to the qcChain array
func (b *baseBridge) UpdateChainState(qcState, lockState, commitState *protocols.State) {
	tStart := time.Now()
	var chainState *protocols.ChainState
	var err error
	// load current consensus state
	b.cbft.wal.LoadChainState(func(cs *protocols.ChainState) error {
		chainState = cs
		return nil
	})

	if chainState == nil {
		log.Trace("ChainState is empty,may be the first time to update chainState")
		if err = b.newChainState(commitState, lockState, qcState); err != nil {
			panic(fmt.Sprintf("update chain state error: %s", err.Error()))
		}
		return
	}
	if !chainState.ValidChainState() {
		panic(fmt.Sprintf("invalid chain state from wal"))
	}
	walCommitNumber := chainState.Commit.QuorumCert.BlockNumber
	commitNumber := commitState.QuorumCert.BlockNumber
	if walCommitNumber != commitNumber && walCommitNumber+1 != commitNumber {
		log.Warn("The chainState of wal and updating are discontinuous,ignore this updating", "walCommit", chainState.Commit.String(), "updateCommit", commitState.String())
		return
	}

	if chainState.Commit.EqualState(commitState) && chainState.Lock.EqualState(lockState) && chainState.QC[0].QuorumCert.BlockNumber == qcState.QuorumCert.BlockNumber {
		err = b.addQCState(qcState, chainState)
	} else {
		err = b.newChainState(commitState, lockState, qcState)
	}

	if err != nil {
		panic(fmt.Sprintf("update chain state error: %s", err.Error()))
	}
	log.Info("Success to update chainState", "commitState", commitState.String(), "lockState", lockState.String(), "qcState", qcState.String(), "elapsed", time.Since(tStart))
}

// newChainState tries to update consensus state to wal
// Need to do continuous block check before writing.
func (b *baseBridge) newChainState(commit *protocols.State, lock *protocols.State, qc *protocols.State) error {
	log.Debug("New chainState", "commitState", commit.String(), "lockState", lock.String(), "qcState", qc.String())
	if !commit.ValidState() || !lock.ValidState() || !qc.ValidState() {
		return errNonContiguous
	}
	// check continuous block chain
	if !contiguousChainBlock(commit.Block, lock.Block) || !contiguousChainBlock(lock.Block, qc.Block) {
		return errNonContiguous
	}
	chainState := &protocols.ChainState{
		Commit: commit,
		Lock:   lock,
		QC:     []*protocols.State{qc},
	}
	return b.cbft.wal.UpdateChainState(chainState)
}

// addQCState tries to add consensus qc state to wal
// Need to do continuous block check before writing.
func (b *baseBridge) addQCState(qc *protocols.State, chainState *protocols.ChainState) error {
	log.Debug("Add qcState", "qcState", qc.String())
	lock := chainState.Lock
	// check continuous block chain
	if !contiguousChainBlock(lock.Block, qc.Block) {
		return errNonContiguous
	}
	chainState.QC = append(chainState.QC, qc)
	return b.cbft.wal.UpdateChainState(chainState)
}

// ConfirmViewChange tries to update ConfirmedViewChange consensus msg to wal.
// at the same time we will record the current fileID and fileSequence.
// the next time the platon node restart, we will recovery the msg from this check point.
func (b *baseBridge) ConfirmViewChange(epoch, viewNumber uint64, block *types.Block, qc *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC, preEpoch, preViewNumber uint64) {
	tStart := time.Now()
	// save the identity location of the wal message in the file system
	meta := &wal.ViewChangeMessage{
		Epoch:      epoch,
		ViewNumber: viewNumber,
	}
	if err := b.cbft.wal.UpdateViewChange(meta); err != nil {
		panic(fmt.Sprintf("update viewChange meta error, err:%s", err.Error()))
	}
	// save ConfirmedViewChange message, the viewChangeQC is last viewChangeQC
	vc := &protocols.ConfirmedViewChange{
		Epoch:        epoch,
		ViewNumber:   viewNumber,
		Block:        block,
		QC:           qc,
		ViewChangeQC: viewChangeQC,
	}
	if err := b.cbft.wal.WriteSync(vc); err != nil {
		panic(fmt.Sprintf("write confirmed viewChange error, err:%s", err.Error()))
	}
	// save last viewChangeQC, for viewChangeQC synchronization
	if viewChangeQC != nil {
		b.cbft.wal.UpdateViewChangeQC(preEpoch, preViewNumber, viewChangeQC)
	}
	log.Debug("Success to confirm viewChange", "confirmedViewChange", vc.String(), "elapsed", time.Since(tStart))
}

// SendViewChange tries to update SendViewChange consensus msg to wal.
func (b *baseBridge) SendViewChange(view *protocols.ViewChange) {
	tStart := time.Now()
	s := &protocols.SendViewChange{
		ViewChange: view,
	}
	if err := b.cbft.wal.WriteSync(s); err != nil {
		panic(fmt.Sprintf("write send viewChange error, err:%s", err.Error()))
	}
	log.Debug("Success to send viewChange", "view", view.String(), "elapsed", time.Since(tStart))
}

// SendPrepareBlock tries to update SendPrepareBlock consensus msg to wal.
func (b *baseBridge) SendPrepareBlock(pb *protocols.PrepareBlock) {
	tStart := time.Now()
	s := &protocols.SendPrepareBlock{
		Prepare: pb,
	}
	if err := b.cbft.wal.WriteSync(s); err != nil {
		panic(fmt.Sprintf("write send prepareBlock error, err:%s", err.Error()))
	}
	log.Debug("Success to send prepareBlock", "prepareBlock", pb.String(), "elapsed", time.Since(tStart))
}

// SendPrepareVote tries to update SendPrepareVote consensus msg to wal.
func (b *baseBridge) SendPrepareVote(block *types.Block, vote *protocols.PrepareVote) {
	tStart := time.Now()
	s := &protocols.SendPrepareVote{
		Block: block,
		Vote:  vote,
	}
	if err := b.cbft.wal.WriteSync(s); err != nil {
		panic(fmt.Sprintf("write send prepareVote error, err:%s", err.Error()))
	}
	log.Debug("Success to send prepareVote", "prepareVote", vote.String(), "elapsed", time.Since(tStart))
}

func (b *baseBridge) GetViewChangeQC(epoch uint64, viewNumber uint64) (*ctypes.ViewChangeQC, error) {
	return b.cbft.wal.GetViewChangeQC(epoch, viewNumber)
}

func (b *baseBridge) Close() {
	b.cbft.wal.Close()
}

// recoveryChainState tries to recovery consensus chainState from wal when the platon node restart.
// need to do some necessary checks based on the latest blockchain block.
// execute commit/lock/qcs block and load the corresponding state to cbft consensus.
func (cbft *Cbft) recoveryChainState(chainState *protocols.ChainState) error {
	cbft.log.Info("Recover chain state from wal", "chainState", chainState.String())
	commit, lock, qcs := chainState.Commit, chainState.Lock, chainState.QC
	// The highest block that has been written to disk

	//rootBlock := cbft.blockChain.GetBlock(cbft.blockChain.CurrentHeader().Hash(), cbft.blockChain.CurrentHeader().Number.Uint64())
	rootBlock := cbft.blockChain.CurrentBlock()

	isCurrent := rootBlock.NumberU64() == commit.Block.NumberU64() && rootBlock.Hash() == commit.Block.Hash()
	isParent := contiguousChainBlock(rootBlock, commit.Block)

	if !isCurrent && !isParent {
		return fmt.Errorf("recovery chain state errror,non contiguous chain block state, curNum:%d, curHash:%s, commitNum:%d, commitHash:%s", rootBlock.NumberU64(), rootBlock.Hash().String(), commit.Block.NumberU64(), commit.Block.Hash().String())
	}
	if isParent {
		// recovery commit state
		if err := cbft.recoveryCommitState(commit, rootBlock); err != nil {
			return err
		}
	}
	// recovery lock state
	if err := cbft.recoveryLockState(lock, commit.Block); err != nil {
		return err
	}
	// recovery qc state
	if err := cbft.recoveryQCState(qcs, lock.Block); err != nil {
		return err
	}
	return nil
}

func (cbft *Cbft) recoveryCommitState(commit *protocols.State, parent *types.Block) error {
	log.Info("Recover commit state", "commitNumber", commit.Block.NumberU64(), "commitHash", commit.Block.Hash(), "parentNumber", parent.NumberU64(), "parentHash", parent.Hash())
	// recovery commit state
	if err := cbft.executeBlock(commit.Block, parent, math.MaxUint32); err != nil {
		return err
	}
	// write commit block to chain
	extra, err := ctypes.EncodeExtra(byte(cbftVersion), commit.QuorumCert)
	if err != nil {
		return err
	}
	commit.Block.SetExtraData(extra)
	if err := cbft.blockCacheWriter.WriteBlock(commit.Block); err != nil {
		return err
	}
	if err := cbft.validatorPool.Commit(commit.Block); err != nil {
		return err
	}
	cbft.recoveryChainStateProcess(protocols.CommitState, commit)
	cbft.blockTree.NewRoot(commit.Block)
	return nil
}

func (cbft *Cbft) recoveryLockState(lock *protocols.State, parent *types.Block) error {
	// recovery lock state
	if err := cbft.executeBlock(lock.Block, parent, math.MaxUint32); err != nil {
		return err
	}
	cbft.recoveryChainStateProcess(protocols.LockState, lock)
	return nil
}

func (cbft *Cbft) recoveryQCState(qcs []*protocols.State, parent *types.Block) error {
	// recovery qc state
	for _, qc := range qcs {
		if err := cbft.executeBlock(qc.Block, parent, math.MaxUint32); err != nil {
			return err
		}
		cbft.recoveryChainStateProcess(protocols.QCState, qc)
	}
	return nil
}

// recoveryChainStateProcess tries to recovery the corresponding state to cbft consensus.
func (cbft *Cbft) recoveryChainStateProcess(stateType uint16, s *protocols.State) {
	cbft.trySwitchValidator(s.Block.NumberU64())
	cbft.tryWalChangeView(s.QuorumCert.Epoch, s.QuorumCert.ViewNumber, s.Block, s.QuorumCert, nil)
	cbft.state.AddQCBlock(s.Block, s.QuorumCert)
	cbft.state.AddQC(s.QuorumCert)
	cbft.blockTree.InsertQCBlock(s.Block, s.QuorumCert)
	cbft.state.SetExecuting(s.QuorumCert.BlockIndex, true)

	switch stateType {
	case protocols.CommitState:
		cbft.state.SetHighestCommitBlock(s.Block)
	case protocols.LockState:
		cbft.state.SetHighestLockBlock(s.Block)
	case protocols.QCState:
		cbft.TrySetHighestQCBlock(s.Block)
	}

	// The state may have reached the automatic switch point, then advance to the next view
	if cbft.validatorPool.EqualSwitchPoint(s.Block.NumberU64()) {
		cbft.log.Info("QCBlock is equal to switchPoint, change epoch", "state", s.String(), "view", cbft.state.ViewString())
		cbft.tryWalChangeView(cbft.state.Epoch()+1, state.DefaultViewNumber, s.Block, s.QuorumCert, nil)
		return
	}
	if s.QuorumCert.BlockIndex+1 == cbft.config.Sys.Amount {
		cbft.log.Info("QCBlock is the last index on the view, change view", "state", s.String(), "view", cbft.state.ViewString())
		cbft.tryWalChangeView(cbft.state.Epoch(), cbft.state.ViewNumber()+1, s.Block, s.QuorumCert, nil)
		return
	}
}

// trySwitch tries to switch next validator.
func (cbft *Cbft) trySwitchValidator(blockNumber uint64) {
	if cbft.validatorPool.ShouldSwitch(blockNumber) {
		if err := cbft.validatorPool.Update(blockNumber, cbft.state.Epoch()+1, cbft.eventMux); err != nil {
			cbft.log.Debug("Update validator error", "err", err.Error())
		}
	}
}

// tryWalChangeView tries to change view.
func (cbft *Cbft) tryWalChangeView(epoch, viewNumber uint64, block *types.Block, qc *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC) {
	if epoch > cbft.state.Epoch() || epoch == cbft.state.Epoch() && viewNumber > cbft.state.ViewNumber() {
		cbft.changeView(epoch, viewNumber, block, qc, viewChangeQC)
	}
}

// recoveryMsg tries to recovery consensus msg from wal when the platon node restart.
func (cbft *Cbft) recoveryMsg(msg interface{}) error {
	cbft.log.Info("Recover journal message from wal", "msgType", reflect.TypeOf(msg))

	switch m := msg.(type) {
	case *protocols.ConfirmedViewChange:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "confirmedViewChange", m.String())
		cbft.tryWalChangeView(m.Epoch, m.ViewNumber, m.Block, m.QC, m.ViewChangeQC)

	case *protocols.SendViewChange:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendViewChange", m.String())

		should, err := cbft.shouldRecovery(m)
		if err != nil {
			return err
		}
		if should {
			node, err := cbft.validatorPool.GetValidatorByNodeID(m.ViewChange.Epoch, cbft.config.Option.NodeID)
			if err != nil {
				return err
			}
			cbft.state.AddViewChange(uint32(node.Index), m.ViewChange)
		}

	case *protocols.SendPrepareBlock:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareBlock", m.String())

		should, err := cbft.shouldRecovery(m)
		if err != nil {
			return err
		}
		if should {
			// execute block
			block := m.Prepare.Block
			if cbft.state.ViewBlockByIndex(m.Prepare.BlockIndex) == nil {
				if err := cbft.executeBlock(block, nil, m.Prepare.BlockIndex); err != nil {
					return err
				}
				cbft.state.SetExecuting(m.Prepare.BlockIndex, true)
			}
			cbft.signMsgByBls(m.Prepare)
			cbft.state.AddPrepareBlock(m.Prepare)
		}

	case *protocols.SendPrepareVote:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareVote", m.String())

		should, err := cbft.shouldRecovery(m)
		if err != nil {
			return err
		}
		if should {
			// execute block
			block := m.Block
			if cbft.state.ViewBlockByIndex(m.Vote.BlockIndex) == nil {
				if err := cbft.executeBlock(block, nil, m.Vote.BlockIndex); err != nil {
					return err
				}
				cbft.state.SetExecuting(m.Vote.BlockIndex, true)
				cbft.state.AddPrepareBlock(&protocols.PrepareBlock{
					Epoch:      m.Vote.Epoch,
					ViewNumber: m.Vote.ViewNumber,
					Block:      block,
					BlockIndex: m.Vote.BlockIndex,
				})
			}

			cbft.state.HadSendPrepareVote().Push(m.Vote)
			node, _ := cbft.validatorPool.GetValidatorByNodeID(m.Vote.Epoch, cbft.config.Option.NodeID)
			cbft.state.AddPrepareVote(uint32(node.Index), m.Vote)
		}
	}
	return nil
}

// contiguousChainBlock check if the two incoming blocks are continuous.
func contiguousChainBlock(p *types.Block, s *types.Block) bool {
	contiguous := p.NumberU64()+1 == s.NumberU64() && p.Hash() == s.ParentHash()
	if !contiguous {
		log.Info("Non contiguous block", "sNumber", s.NumberU64(), "sParentHash", s.ParentHash(), "pNumber", p.NumberU64(), "pHash", p.Hash())
	}
	return contiguous
}

// executeBlock call blockCacheWriter to execute block.
func (cbft *Cbft) executeBlock(block *types.Block, parent *types.Block, index uint32) error {
	if parent == nil {
		if parent, _ = cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1); parent == nil {
			if parent = cbft.state.ViewBlockByIndex(index - 1); parent == nil {
				return fmt.Errorf("find executable block's parent failed, blockNum:%d, blockHash:%s", block.NumberU64(), block.Hash().String())
			}
		}
	}
	if err := cbft.blockCacheWriter.Execute(block, parent); err != nil {
		return fmt.Errorf("execute block failed, blockNum:%d, blockHash:%s, parentNum:%d, parentHash:%s, err:%s", block.NumberU64(), block.Hash().String(), parent.NumberU64(), parent.Hash().String(), err.Error())
	}
	return nil
}

// shouldRecovery check if the consensus msg needs to be recovery.
// if the msg does not belong to the current view or the msg number is smaller than the qc number discard it.
func (cbft *Cbft) shouldRecovery(msg protocols.WalMsg) (bool, error) {
	if cbft.higherViewState(msg) {
		return false, fmt.Errorf("higher view state, curEpoch:%d, curViewNum:%d, msgEpoch:%d, msgViewNum:%d", cbft.state.Epoch(), cbft.state.ViewNumber(), msg.Epoch(), msg.ViewNumber())
	}
	if cbft.lowerViewState(msg) {
		// The state may have reached the automatic switch point, so advance to the next view
		return false, nil
	}
	// equalViewState
	highestQCBlockBn, _ := cbft.HighestQCBlockBn()
	return msg.BlockNumber() > highestQCBlockBn, nil
}

// equalViewState check if the msg view is equal with current.
func (cbft *Cbft) equalViewState(msg protocols.WalMsg) bool {
	return msg.Epoch() == cbft.state.Epoch() && msg.ViewNumber() == cbft.state.ViewNumber()
}

// lowerViewState check if the msg view is lower than current.
func (cbft *Cbft) lowerViewState(msg protocols.WalMsg) bool {
	return msg.Epoch() < cbft.state.Epoch() || msg.Epoch() == cbft.state.Epoch() && msg.ViewNumber() < cbft.state.ViewNumber()
}

// higherViewState check if the msg view is higher than current.
func (cbft *Cbft) higherViewState(msg protocols.WalMsg) bool {
	return msg.Epoch() > cbft.state.Epoch() || msg.Epoch() == cbft.state.Epoch() && msg.ViewNumber() > cbft.state.ViewNumber()
}
