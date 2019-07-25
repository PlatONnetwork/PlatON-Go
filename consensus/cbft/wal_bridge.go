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
func (cbft *Cbft) newChainState(commit *protocols.State, lock *protocols.State, qc []*protocols.State) error {
	if commit == nil || commit.Block == nil || lock == nil || lock.Block == nil || len(qc) <= 0 {
		return errNonContiguous
	}
	if commit.Block.NumberU64() != lock.Block.NumberU64()+1 || commit.Block.Hash() != lock.Block.ParentHash() {
		return errNonContiguous
	}
	for _, q := range qc {
		if lock.Block.NumberU64() != q.Block.NumberU64()+1 || lock.Block.Hash() != q.Block.ParentHash() {
			return errNonContiguous
		}
	}
	chainState := &protocols.ChainState{
		Commit: commit,
		Lock:   lock,
		QC:     qc,
	}
	return cbft.wal.UpdateChainState(chainState)
}

// addQCState tries to add consensus qc state to wal
// Need to do continuous block check before writing.
func (cbft *Cbft) addQCState(qc *protocols.State) error {
	var chainState *protocols.ChainState
	cbft.wal.LoadChainState(func(cs *protocols.ChainState) {
		chainState = cs
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

func (cbft *Cbft) sendPrepareBlock(block *types.Block) {
	s := &protocols.SendPrepareBlock{
		Block: block,
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

func (cbft *Cbft) recoveryChainState(chainState *protocols.ChainState) {
	cbft.log.Debug("Recover chain state from wal", "chainState", chainState.String())
	// TODO
}

func (cbft *Cbft) recoveryMsg(msg interface{}) {
	cbft.log.Debug("Recover journal message from wal", "msgType", reflect.TypeOf(msg))

	switch m := msg.(type) {
	case *protocols.ConfirmedViewChange:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "confirmedViewChange", m.String())
		// TODO : 去掉number小于highestQC的
		cbft.changeView(m.Epoch, m.ViewNumber, m.Block, m.QC, m.ViewChangeQC)

	case *protocols.SendViewChange:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendViewChange", m.ViewChange.String())
		// TODO : 去掉number小于highestQC的

	case *protocols.SendPrepareBlock:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareBlock,number", m.Block.NumberU64(), "hash", m.Block.Hash())
		// TODO : 去掉number小于highestQC的
		if m.Block.NumberU64() > cbft.HighestQCBlockBn() {
			// 执行区块
			block := m.Block
			parent, _ := cbft.blockTree.FindBlockAndQC(block.ParentHash(), block.NumberU64()-1)
			if parent == nil {
				panic(fmt.Sprintf("Find executable block's parent failed :[%d,%s]", block.NumberU64(), block.Hash()))
			}
			if err := cbft.blockCacheWriter.Execute(block, parent); err != nil {
				panic(fmt.Sprintf("Execute block failed", "hash", block.Hash(), "number", block.NumberU64(), "error", err))
			}
			// 写回内存

		}

	case *protocols.SendPrepareVote:
		cbft.log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareVote", m.Vote.String())
		// TODO : 去掉number小于highestQC的
	}
}
