package cbft

import (
	"errors"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/wal"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
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

func (cbft *Cbft) confirmViewChange(epoch uint64, viewNumber uint64) error {
	meta := &wal.ViewChangeMessage{
		Epoch:      epoch,
		ViewNumber: viewNumber,
	}
	if err := cbft.wal.UpdateViewChange(meta); err != nil {
		return err
	}
	vc := &protocols.ConfirmedViewChange{
		Epoch:      epoch,
		ViewNumber: viewNumber,
	}
	return cbft.wal.WriteSync(vc)
}

func (cbft *Cbft) sendViewChange(view *protocols.ViewChange) error {
	s := &protocols.SendViewChange{
		ViewChange: view,
	}
	return cbft.wal.WriteSync(s)
}

func (cbft *Cbft) sendPrepareBlock(block *types.Block) error {
	s := &protocols.SendPrepareBlock{
		Block: block,
	}
	return cbft.wal.WriteSync(s)
}

func (cbft *Cbft) sendPrepareVote(block *types.Block, vote *protocols.PrepareVote) error {
	s := &protocols.SendPrepareVote{
		Block: block,
		Vote:  vote,
	}
	return cbft.wal.WriteSync(s)
}

func (cbft *Cbft) recoveryMsg(msg interface{}) {
	cbft.log.Debug("Recover journal message from wal", "msgType", reflect.TypeOf(msg))

	switch m := msg.(type) {
	case *protocols.ConfirmedViewChange:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "confirmedViewChange,epoch", m.Epoch, "viewNumber", m.ViewNumber)
		// 去掉number小于highestQC的
		// TODO

	case *protocols.SendViewChange:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendViewChange", m.ViewChange.String())
		// TODO

	case *protocols.SendPrepareBlock:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareBlock,number", m.Block.NumberU64(), "hash", m.Block.Hash())
		// TODO

	case *protocols.SendPrepareVote:
		log.Debug("Load journal message from wal", "msgType", reflect.TypeOf(msg), "sendPrepareVote", m.Vote.String())
		// TODO
	}
}

func (cbft *Cbft) recoveryChainState(chainState *protocols.ChainState) {
	cbft.log.Debug("Recover chain state from wal", "chainState", chainState.String())
	// TODO
}
