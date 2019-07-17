package protocols

import (
	"fmt"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
)

const (
	SendPrepareBlockMsg    = 0x00
	SendPrepareVoteMsg     = 0x01
	SendViewChangeMsg      = 0x02
	ConfirmedViewChangeMsg = 0x03
)

type State struct {
	PrepareBlock *PrepareBlock
	QuorumCert   *types.QuorumCert
}

type ChainState struct {
	Commit *State
	Lock   *State
	QC     []*State
}

func (cs *ChainState) String() string {
	if cs == nil {
		return ""
	}
	return fmt.Sprintf("[commitNum:%d,commitHash:%s,lockNum:%d,lockHash:%s,qcNum:%d,qcHash:%s]",
		cs.Commit.PrepareBlock.Block.NumberU64(), cs.Commit.PrepareBlock.Block.Hash(),
		cs.Lock.PrepareBlock.Block.NumberU64(), cs.Lock.PrepareBlock.Block.Hash(),
		cs.QC[0].PrepareBlock.Block.NumberU64(), cs.QC[0].PrepareBlock.Block.Hash())
}

// SendPrepareBlock
type SendPrepareBlock struct {
	PrepareBlock *PrepareBlock
}

// SendPrepareVote
type SendPrepareVote struct {
	PrepareBlock *PrepareBlock
	PrepareVote  *PrepareVote
}

// SendViewChange
type SendViewChange struct {
	ViewChange *ViewChange
}

// ConfirmedViewChange
type ConfirmedViewChange struct {
	ViewChange *ViewChange
	Master     bool
}

var (
	WalMessages = []interface{}{
		SendPrepareBlock{},
		SendPrepareVote{},
		SendViewChange{},
		ConfirmedViewChange{},
	}
)

func WalMessageType(msg interface{}) uint64 {
	switch msg.(type) {
	case *SendPrepareBlock:
		return SendPrepareBlockMsg
	case *SendPrepareVote:
		return SendPrepareVoteMsg
	case *SendViewChange:
		return SendViewChangeMsg
	case *ConfirmedViewChange:
		return ConfirmedViewChangeMsg
	}
	panic(fmt.Sprintf("invalid wal msg type %v", reflect.TypeOf(msg)))
}
