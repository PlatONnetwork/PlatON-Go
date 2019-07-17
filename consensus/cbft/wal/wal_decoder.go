package wal

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

//
type JournalMessageSendPrepareBlock struct {
	Timestamp uint64
	Data      *protocols.SendPrepareBlock
}

//
type JournalMessageSendPrepareVote struct {
	Timestamp uint64
	Data      *protocols.SendPrepareVote
}

//
type JournalMessageSendViewChange struct {
	Timestamp uint64
	Data      *protocols.SendViewChange
}

//
type JournalMessageConfirmedViewChange struct {
	Timestamp uint64
	Data      *protocols.ConfirmedViewChange
}

func WALDecode(pack []byte, msgType uint16) (interface{}, error) {
	switch msgType {
	case protocols.SendPrepareBlockMsg:
		var j JournalMessageSendPrepareBlock
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return j.Data, nil
		} else {
			return nil, err
		}
	case protocols.SendPrepareVoteMsg:
		var j JournalMessageSendPrepareVote
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return j.Data, nil
		} else {
			return nil, err
		}
	case protocols.SendViewChangeMsg:
		var j JournalMessageSendViewChange
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return j.Data, nil
		} else {
			return nil, err
		}
	case protocols.ConfirmedViewChangeMsg:
		var j JournalMessageConfirmedViewChange
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return j.Data, nil
		} else {
			return nil, err
		}
	}
	panic(fmt.Sprintf("invalid msg type %d", msgType))
}
