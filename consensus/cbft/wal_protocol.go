package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

const (
	SendPrepareBlockMsg    = 0x64
	SendViewChangeMsg      = 0x65
	ConfirmedViewChangeMsg = 0x66
)

type sendPrepareBlock struct {
	PrepareBlock *prepareBlock
}

func (s *sendPrepareBlock) String() string {
	return ""
}

func (s *sendPrepareBlock) MsgHash() common.Hash {
	return common.Hash{}
}

func (s *sendPrepareBlock) BHash() common.Hash {
	return common.Hash{}
}

type sendViewChange struct {
	ViewChange *viewChange
	//viewChangeVotes ViewChangeVotes
	Master bool
}

func (s *sendViewChange) String() string {
	return ""
}

func (s *sendViewChange) MsgHash() common.Hash {
	return common.Hash{}
}

func (s *sendViewChange) BHash() common.Hash {
	return common.Hash{}
}

type confirmedViewChange struct {
	ViewChange      *viewChange
	ViewChangeResp  *viewChangeVote `rlp:"nil"`
	ViewChangeVotes []*viewChangeVote
	Master          bool
}

func (c *confirmedViewChange) String() string {
	return ""
}

func (c *confirmedViewChange) MsgHash() common.Hash {
	return common.Hash{}
}

func (c *confirmedViewChange) BHash() common.Hash {
	return common.Hash{}
}

var (
	wal_messages = []interface{}{
		prepareBlock{},
		prepareVote{},
		viewChange{},
		viewChangeVote{},
		confirmedPrepareBlock{},
		getPrepareVote{},
		prepareVotes{},
		getPrepareBlock{},
		getHighestPrepareBlock{},
		highestPrepareBlock{},
		cbftStatusData{},
		prepareBlockHash{},
		sendPrepareBlock{},
		sendViewChange{},
		confirmedViewChange{},
	}
)

func WalMessageType(msg interface{}) uint64 {
	switch msg.(type) {
	case *sendPrepareBlock:
		return SendPrepareBlockMsg
	case *sendViewChange:
		return SendViewChangeMsg
	case *confirmedViewChange:
		return ConfirmedViewChangeMsg
	}
	return MessageType(msg)
}
