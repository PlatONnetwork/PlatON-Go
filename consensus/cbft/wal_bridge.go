package cbft

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/node"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var (
	errNonContiguous = errors.New("non contiguous chain block state")
)

// Bridge encapsulates functions required to update consensus state and consensus msg.
// As a bridge layer for cbft and wal.
type Bridge interface {
	UpdateChainState(qcState, lockState, commitState *protocols.State)
	ConfirmViewChange(epoch, viewNumber uint64, block *types.Block, qc *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC)
	SendViewChange(view *protocols.ViewChange)
	SendPrepareBlock(pb *protocols.PrepareBlock)
	SendPrepareVote(block *types.Block, vote *protocols.PrepareVote)
}

// emptyBridge is a empty implementation for Bridge
type emptyBridge struct {
}

func (b *emptyBridge) UpdateChainState(qcState, lockState, commitState *protocols.State) {
}

func (b *emptyBridge) ConfirmViewChange(epoch, viewNumber uint64, block *types.Block, qc *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC) {
}

func (b *emptyBridge) SendViewChange(view *protocols.ViewChange) {
}

func (b *emptyBridge) SendPrepareBlock(pb *protocols.PrepareBlock) {
}

func (b *emptyBridge) SendPrepareVote(block *types.Block, vote *protocols.PrepareVote) {
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
	if lockState == nil || commitState == nil {
		if err := b.cbft.addQCState(qcState); err != nil {
			panic(fmt.Sprintf("update chain state error: %s", err.Error()))
		}
	} else {
		if err := b.cbft.newChainState(commitState, lockState, qcState); err != nil {
			panic(fmt.Sprintf("update chain state error: %s", err.Error()))
		}
	}
}

// ConfirmViewChange tries to update ConfirmedViewChange consensus msg to wal.
// at the same time we will record the current fileID and fileSequence.
// the next time the platon node restart, we will recovery the msg from this check point.
func (b *baseBridge) ConfirmViewChange(epoch, viewNumber uint64, block *types.Block, qc *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC) {
	meta := &wal.ViewChangeMessage{
		Epoch:      epoch,
		ViewNumber: viewNumber,
	}
	if err := b.cbft.wal.UpdateViewChange(meta); err != nil {
		panic(fmt.Sprintf("update viewChange meta error, err:%s", err.Error()))
	}
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
}

// SendViewChange tries to update SendViewChange consensus msg to wal.
func (b *baseBridge) SendViewChange(view *protocols.ViewChange) {
	s := &protocols.SendViewChange{
		ViewChange: view,
	}
	if err := b.cbft.wal.WriteSync(s); err != nil {
		panic(fmt.Sprintf("write send viewChange error, err:%s", err.Error()))
	}
}

// SendPrepareBlock tries to update SendPrepareBlock consensus msg to wal.
func (b *baseBridge) SendPrepareBlock(pb *protocols.PrepareBlock) {
	s := &protocols.SendPrepareBlock{
		Prepare: pb,
	}
	if err := b.cbft.wal.WriteSync(s); err != nil {
		panic(fmt.Sprintf("write send prepareBlock error, err:%s", err.Error()))
	}
}

// SendPrepareVote tries to update SendPrepareVote consensus msg to wal.
func (b *baseBridge) SendPrepareVote(block *types.Block, vote *protocols.PrepareVote) {
	s := &protocols.SendPrepareVote{
		Block: block,
		Vote:  vote,
	}
	if err := b.cbft.wal.WriteSync(s); err != nil {
		panic(fmt.Sprintf("write send prepareVote error, err:%s", err.Error()))
	}
}

// newChainState tries to update consensus state to wal
// Need to do continuous block check before writing.
func (cbft *Cbft) newChainState(commit *protocols.State, lock *protocols.State, qc *protocols.State) error {
	if commit == nil || commit.Block == nil || lock == nil || lock.Block == nil || qc == nil || qc.Block == nil {
		return errNonContiguous
	}
	// check continuous block chain
	if !cbft.contiguousChainBlock(commit.Block, lock.Block) || !cbft.contiguousChainBlock(lock.Block, qc.Block) {
		return errNonContiguous
	}
	chainState := &protocols.ChainState{
		Commit: commit,
		Lock:   lock,
		QC:     []*protocols.State{qc},
	}
	return cbft.wal.UpdateChainState(chainState)
}

// addQCState tries to add consensus qc state to wal
// Need to do continuous block check before writing.
func (cbft *Cbft) addQCState(qc *protocols.State) error {
	var chainState *protocols.ChainState
	// load current consensus state
	cbft.wal.LoadChainState(func(cs *protocols.ChainState) error {
		chainState = cs
		return nil
	})
	if chainState == nil || chainState.Commit == nil || chainState.Lock == nil || len(chainState.QC) <= 0 {
		return nil
	}
	lock := chainState.Lock
	// check continuous block chain
	if !cbft.contiguousChainBlock(lock.Block, qc.Block) {
		return errNonContiguous
	}
	chainState.QC = append(chainState.QC, qc)
	return cbft.wal.UpdateChainState(chainState)
}

// recoveryChainState tries to recovery consensus chainState from wal when the platon node restart.
// need to do some necessary checks based on the latest blockchain block.
// execute commit/lock/qcs block and load the corresponding state to cbft consensus.
func (cbft *Cbft) recoveryChainState(chainState *protocols.ChainState) error {
	cbft.log.Debug("Recover chain state from wal", "chainState", chainState.String())
	commit, lock, qcs := chainState.Commit, chainState.Lock, chainState.QC
	// The highest block that has been written to disk
	rootBlock := cbft.blockChain.GetBlock(cbft.blockChain.CurrentHeader().Hash(), cbft.blockChain.CurrentHeader().Number.Uint64())

	isCurrent := rootBlock.NumberU64() == commit.Block.NumberU64() && rootBlock.Hash() == commit.Block.Hash()
	isParent := cbft.contiguousChainBlock(rootBlock, commit.Block)

	if !isCurrent && !isParent {
		return fmt.Errorf("recovery chain state errror,non contiguous chain block state, curNum:%d, curHash:%s, commitNum:%d, commitHash:%s", rootBlock.NumberU64(), rootBlock.Hash().String(), commit.Block.NumberU64(), commit.Block.Hash().String())
	}
	if isParent {
		// recovery commit state
		if err := cbft.executeBlock(commit.Block, rootBlock); err != nil {
			return err
		}
		cbft.recoveryChainStateProcess(protocols.CommitState, commit)
	}
	// recovery lock state
	if err := cbft.executeBlock(lock.Block, commit.Block); err != nil {
		return err
	}
	cbft.recoveryChainStateProcess(protocols.LockState, lock)
	// recovery qc state
	for _, qc := range qcs {
		if err := cbft.executeBlock(qc.Block, lock.Block); err != nil {
			return err
		}
		cbft.recoveryChainStateProcess(protocols.QCState, qc)
	}
	return nil
}

// recoveryChainStateProcess tries to recovery the corresponding state to cbft consensus.
func (cbft *Cbft) recoveryChainStateProcess(stateType uint16, state *protocols.State) {
	cbft.tryWalChangeView(state.QuorumCert.Epoch, state.QuorumCert.ViewNumber, state.Block, state.QuorumCert, nil)
	cbft.state.AddQCBlock(state.Block, state.QuorumCert)
	cbft.state.AddQC(state.QuorumCert)
	cbft.blockTree.InsertQCBlock(state.Block, state.QuorumCert)
	//cbft.state.SetHighestExecutedBlock(state.Block)
	cbft.state.SetExecuting(state.QuorumCert.BlockIndex, true)

	switch stateType {
	case protocols.CommitState:
		cbft.state.SetHighestCommitBlock(state.Block)
	case protocols.LockState:
		cbft.state.SetHighestLockBlock(state.Block)
	case protocols.QCState:
		cbft.state.SetHighestQCBlock(state.Block)
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
	cbft.log.Debug("Recover journal message from wal", "msgType", reflect.TypeOf(msg))

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
			node, _ := cbft.validatorPool.GetValidatorByNodeID(cbft.state.HighestQCBlock().NumberU64(), cbft.config.Option.NodeID)
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
			if err := cbft.executeBlock(block, nil); err != nil {
				return err
			}

			cbft.signMsgByBls(m.Prepare)
			cbft.state.SetExecuting(m.Prepare.BlockIndex, true)
			cbft.state.AddPrepareBlock(m.Prepare)
			//cbft.OnPrepareBlock("", m.Prepare)
			//cbft.signBlock(block.Hash(), block.NumberU64(), m.Prepare.BlockIndex)
			//cbft.findQCBlock()
			//cbft.state.SetHighestExecutedBlock(block)
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
			if err := cbft.executeBlock(block, nil); err != nil {
				return err
			}

			cbft.state.AddPrepareBlock(&protocols.PrepareBlock{
				Epoch:      m.Vote.Epoch,
				ViewNumber: m.Vote.ViewNumber,
				Block:      block,
				BlockIndex: m.Vote.BlockIndex,
			})
			cbft.state.SetExecuting(m.Vote.BlockIndex, true)
			cbft.state.HadSendPrepareVote().Push(m.Vote)
			node, _ := cbft.validatorPool.GetValidatorByNodeID(m.Vote.BlockNum(), cbft.config.Option.NodeID)
			cbft.state.AddPrepareVote(uint32(node.Index), m.Vote)
			//cbft.state.SetHighestExecutedBlock(block)
		}
	}
	return nil
}

// contiguousChainBlock check if the two incoming blocks are continuous.
func (cbft *Cbft) contiguousChainBlock(p *types.Block, s *types.Block) bool {
	return p.NumberU64() == s.NumberU64()-1 && p.Hash() == s.ParentHash()
}

// executeBlock call blockCacheWriter to execute block.
func (cbft *Cbft) executeBlock(block *types.Block, parent *types.Block) error {
	if parent == nil {
		if parent, _ = cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1); parent == nil {
			if parent = cbft.state.HighestExecutedBlock(); parent == nil {
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
	if !cbft.equalViewState(msg) {
		return false, fmt.Errorf("non equal view state, curEpoch:%d, curViewNum:%d, preEpoch:%d, preViewNum:%d", cbft.state.Epoch(), cbft.state.ViewNumber(), msg.Epoch(), msg.ViewNumber())
	}
	highestQCBlockBn, _ := cbft.HighestQCBlockBn()
	return msg.BlockNumber() > highestQCBlockBn, nil
}

// equalViewState check if the view is equal.
func (cbft *Cbft) equalViewState(msg protocols.WalMsg) bool {
	return msg.Epoch() == cbft.state.Epoch() && msg.ViewNumber() == cbft.state.ViewNumber()
}
