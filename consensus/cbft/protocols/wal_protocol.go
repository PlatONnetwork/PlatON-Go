package protocols

import (
	"fmt"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
)

const (
	ConfirmedViewChangeMsg = 0x01
	SendViewChangeMsg      = 0x02
	SendPrepareBlockMsg    = 0x03
	SendPrepareVoteMsg     = 0x04
)

type State struct {
	Block      *types.Block
	QuorumCert *ctypes.QuorumCert
}

// ChainState indicates the latest consensus state.
type ChainState struct {
	Commit *State
	Lock   *State
	QC     []*State
}

func (cs *ChainState) String() string {
	if cs == nil {
		return ""
	}
	return fmt.Sprintf("[commitNum:%d, commitHash:%s, lockNum:%d, lockHash:%s, qcNum:%d, qcHash:%s]",
		cs.Commit.Block.NumberU64(), cs.Commit.Block.Hash(),
		cs.Lock.Block.NumberU64(), cs.Lock.Block.Hash(),
		cs.QC[0].Block.NumberU64(), cs.QC[0].Block.Hash())
}

// ConfirmedViewChange indicates the latest confirmed view.
type ConfirmedViewChange struct {
	Epoch        uint64
	ViewNumber   uint64
	Block        *types.Block
	QC           *ctypes.QuorumCert
	ViewChangeQC *ctypes.ViewChangeQC
}

func (c *ConfirmedViewChange) String() string {
	if c == nil {
		return ""
	}
	return fmt.Sprintf("[epoch:%d, viewNumber:%d, blockIndex:%d, blockNumber:%d, blockHash:%s]", c.Epoch, c.ViewNumber, c.QC.BlockIndex, c.QC.BlockNumber, c.QC.BlockHash.String())
}

// SendViewChange indicates the viewChange sent by the local node.
type SendViewChange struct {
	ViewChange *ViewChange
}

// SendPrepareBlock indicates the prepareBlock sent by the local node.
type SendPrepareBlock struct {
	Block *types.Block
}

// SendPrepareVote indicates the prepareVote sent by the local node.
type SendPrepareVote struct {
	Block *types.Block
	Vote  *PrepareVote
}

var (
	WalMessages = []interface{}{
		ConfirmedViewChange{},
		SendViewChange{},
		SendPrepareBlock{},
		SendPrepareVote{},
	}
)

func WalMessageType(msg interface{}) uint64 {
	switch msg.(type) {
	case *ConfirmedViewChange:
		return ConfirmedViewChangeMsg
	case *SendViewChange:
		return SendViewChangeMsg
	case *SendPrepareBlock:
		return SendPrepareBlockMsg
	case *SendPrepareVote:
		return SendPrepareVoteMsg
	}
	panic(fmt.Sprintf("invalid wal msg type %v", reflect.TypeOf(msg)))
}
