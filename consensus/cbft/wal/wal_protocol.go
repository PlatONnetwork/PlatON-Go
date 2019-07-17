package wal

import (
	"fmt"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
)

const (
	SendPrepareBlockMsg    = 0x00
	SendPrepareVoteMsg     = 0x01
	SendViewChangeMsg      = 0x02
	ConfirmedViewChangeMsg = 0x03
)

// sendPrepareBlock
type sendPrepareBlock struct {
	PrepareBlock *protocols.PrepareBlock
}

// sendPrepareVote
type sendPrepareVote struct {
	PrepareBlock *protocols.PrepareBlock
	PrepareVote  *protocols.PrepareVote
}

// sendViewChange
type sendViewChange struct {
	ViewChange *protocols.ViewChange
}

// confirmedViewChange
type confirmedViewChange struct {
	ViewChange *protocols.ViewChange
	Master     bool
}

var (
	wal_messages = []interface{}{
		sendPrepareBlock{},
		sendPrepareVote{},
		sendViewChange{},
		confirmedViewChange{},
	}
)

func WalMessageType(msg interface{}) uint64 {
	switch msg.(type) {
	case *sendPrepareBlock:
		return SendPrepareBlockMsg
	case *sendPrepareVote:
		return SendPrepareVoteMsg
	case *sendViewChange:
		return SendViewChangeMsg
	case *confirmedViewChange:
		return ConfirmedViewChangeMsg
	}
	panic(fmt.Sprintf("invalid wal msg type %v", reflect.TypeOf(msg)))
}
