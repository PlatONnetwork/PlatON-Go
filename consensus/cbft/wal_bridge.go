package cbft

import (
	"errors"
	"fmt"
	"reflect"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var (
	errNonContiguous = errors.New("non contiguous chain block state")
)

// newChainState tries to update consensus state to wal
// Need to do continuous block check before writing.
func (cbft *Cbft) newChainState(commit *protocols.State, lock *protocols.State, qc *protocols.State) error {
	if commit == nil || commit.Block == nil || lock == nil || lock.Block == nil || qc == nil || qc.Block == nil {
		return errNonContiguous
	}
	if commit.Block.NumberU64() != lock.Block.NumberU64()+1 || commit.Block.Hash() != lock.Block.ParentHash() ||
		lock.Block.NumberU64() != qc.Block.NumberU64()+1 || lock.Block.Hash() != qc.Block.ParentHash() {
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
	cbft.wal.LoadChainState(func(cs *protocols.ChainState) error {
		chainState = cs
		return nil
	})
	lock := chainState.Lock
	if lock.Block.NumberU64() != qc.Block.NumberU64()+1 || lock.Block.Hash() != qc.Block.ParentHash() {
		return errNonContiguous
	}
	chainState.QC = append(chainState.QC, qc)
	return cbft.wal.UpdateChainState(chainState)
}

func (cbft *Cbft) confirmViewChange(epoch, viewNumber uint64, block *types.Block, qc *ctypes.QuorumCert, viewChangeQC *ctypes.ViewChangeQC) {
	meta := &wal.ViewChangeMessage{
		Epoch:      epoch,
		ViewNumber: viewNumber,
	}
	if err := cbft.wal.UpdateViewChange(meta); err != nil {
		panic(fmt.Sprintf("update viewChange meta error: %s", err.Error()))
	}
	vc := &protocols.ConfirmedViewChange{
		Epoch:        epoch,
		ViewNumber:   viewNumber,
		Block:        block,
		QC:           qc,
		ViewChangeQC: viewChangeQC,
	}
	if err := cbft.wal.WriteSync(vc); err != nil {
		panic(fmt.Sprintf("write confirmed viewChange error: %s", err.Error()))
	}
}

func (cbft *Cbft) sendViewChange(view *protocols.ViewChange) {
	s := &protocols.SendViewChange{
		ViewChange: view,
	}
	if err := cbft.wal.WriteSync(s); err != nil {
		panic(fmt.Sprintf("write send viewChange error: %s", err.Error()))
	}
}

func (cbft *Cbft) sendPrepareBlock(pb *protocols.PrepareBlock) {
	s := &protocols.SendPrepareBlock{
		Prepare: pb,
	}
	if err := cbft.wal.WriteSync(s); err != nil {
		panic(fmt.Sprintf("write send prepareBlock error: %s", err.Error()))
	}
}

func (cbft *Cbft) sendPrepareVote(block *types.Block, vote *protocols.PrepareVote) {
	s := &protocols.SendPrepareVote{
		Block: block,
		Vote:  vote,
	}
	if err := cbft.wal.WriteSync(s); err != nil {
		panic(fmt.Sprintf("write send prepareVote error: %s", err.Error()))
	}
}

func (cbft *Cbft) recoveryChainState(chainState *protocols.ChainState) error {
	cbft.log.Debug("Recover chain state from wal", "chainState", chainState.String())
	commit, lock, qcs := chainState.Commit, chainState.Lock, chainState.QC
	currentBlock := cbft.blockChain.GetBlock(cbft.blockChain.CurrentHeader().Hash(), cbft.blockChain.CurrentHeader().Number.Uint64())

	isCurrent := currentBlock.NumberU64() == commit.Block.NumberU64() && currentBlock.Hash() == commit.Block.Hash()
	isParent := currentBlock.NumberU64()+1 == commit.Block.NumberU64() && currentBlock.Hash() == commit.Block.ParentHash()

	if !isCurrent && !isParent {
		return errors.New(fmt.Sprintf("recovery chain state errror,non contiguous chain block state", "curNumber", currentBlock.NumberU64(),
			"curHash", currentBlock.Hash(), "commitNumber", commit.Block.NumberU64(), "commitHash", commit.Block.Hash()))
	}
	if isParent {
		// recovery commit state
		if err := cbft.blockCacheWriter.Execute(commit.Block, currentBlock); err != nil {
			return errors.New(fmt.Sprintf("execute commit block failed", "hash", commit.Block.Hash(), "number", commit.Block.NumberU64(), "error", err))
		}
		cbft.blockTree.InsertQCBlock(commit.Block, commit.QuorumCert)
		cbft.state.SetHighestCommitBlock(commit.Block)
	}
	// recovery lock state
	if err := cbft.blockCacheWriter.Execute(lock.Block, commit.Block); err != nil {
		return errors.New(fmt.Sprintf("execute lock block failed", "hash", lock.Block.Hash(), "number", lock.Block.NumberU64(), "error", err))
	}
	cbft.blockTree.InsertQCBlock(lock.Block, lock.QuorumCert)
	cbft.state.SetHighestLockBlock(lock.Block)
	// recovery qc state
	for _, qc := range qcs {
		if err := cbft.blockCacheWriter.Execute(qc.Block, lock.Block); err != nil {
			return errors.New(fmt.Sprintf("execute qc block failed", "hash", qc.Block.Hash(), "number", qc.Block.NumberU64(), "error", err))
		}
		cbft.blockTree.InsertQCBlock(qc.Block, qc.QuorumCert)
		cbft.state.SetHighestExecutedBlock(qc.Block)
		cbft.state.SetHighestQCBlock(qc.Block)
	}
	return nil
}

func (cbft *Cbft) recoveryMsg(msg interface{}) error {
	cbft.log.Debug("Recover journal message from wal", "msgType", reflect.TypeOf(msg))

	switch m := msg.(type) {
	case *protocols.ConfirmedViewChange:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "confirmedViewChange", m.String())
		cbft.changeView(m.Epoch, m.ViewNumber, m.Block, m.QC, m.ViewChangeQC)

	case *protocols.SendViewChange:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendViewChange", m.String())
		// TODO : 去掉number小于highestQC的

	case *protocols.SendPrepareBlock:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareBlock", m.String())
		if m.Prepare.Epoch != cbft.state.Epoch() || m.Prepare.ViewNumber != cbft.state.ViewNumber() {
			return errors.New(fmt.Sprintf("non equal view state", "curEpoch", cbft.state.Epoch(), "curViewNum", cbft.state.ViewNumber(),
				"preEpoch", m.Prepare.Epoch, "preViewNum", m.Prepare.ViewNumber))
		}
		block := m.Prepare.Block
		if block.NumberU64() > cbft.HighestQCBlockBn() {
			// execute block
			parent, _ := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
			if parent == nil {
				return errors.New(fmt.Sprintf("find executable block's parent failed :[%d,%s]", block.NumberU64(), block.Hash()))
			}
			if err := cbft.blockCacheWriter.Execute(block, parent); err != nil {
				return errors.New(fmt.Sprintf("execute block failed :[%d,%s]", block.NumberU64(), block.Hash()))
			}

			cbft.signMsgByBls(m.Prepare)
			cbft.state.SetExecuting(m.Prepare.BlockIndex, true)
			//cbft.OnPrepareBlock("", m.Prepare)
			cbft.state.AddPrepareBlock(m.Prepare)
			cbft.signBlock(block.Hash(), block.NumberU64(), m.Prepare.BlockIndex)
			//cbft.findQCBlock()
			cbft.state.SetHighestExecutedBlock(block)
		}

	case *protocols.SendPrepareVote:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareVote", m.String())
		if m.Vote.Epoch != cbft.state.Epoch() || m.Vote.ViewNumber != cbft.state.ViewNumber() {
			return errors.New(fmt.Sprintf("non equal view state", "curEpoch", cbft.state.Epoch(), "curViewNum", cbft.state.ViewNumber(),
				"preEpoch", m.Vote.Epoch, "preViewNum", m.Vote.ViewNumber))
		}
		block := m.Block
		if block.NumberU64() > cbft.HighestQCBlockBn() {
			// execute block
			parent, _ := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
			if parent == nil {
				return errors.New(fmt.Sprintf("find executable block's parent failed :[%d,%s]", block.NumberU64(), block.Hash()))
			}
			if err := cbft.blockCacheWriter.Execute(block, parent); err != nil {
				return errors.New(fmt.Sprintf("execute block failed :[%d,%s]", block.NumberU64(), block.Hash()))
			}
			cbft.state.HadSendPrepareVote().Push(m.Vote)
			node, _ := cbft.validatorPool.GetValidatorByNodeID(cbft.state.HighestQCBlock().NumberU64(), cbft.config.Option.NodeID)
			cbft.state.AddPrepareVote(uint32(node.Index), m.Vote)
		}
	}
	return nil
}
