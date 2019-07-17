package wal

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

//
type MsgInfoSendPrepareBlock struct {
	Msg *sendPrepareBlock
}

type JournalMessageSendPrepareBlock struct {
	Timestamp uint64
	Data      *MsgInfoSendPrepareBlock
}

//
type MsgInfoSendPrepareVote struct {
	Msg *sendPrepareVote
}

type JournalMessageSendPrepareVote struct {
	Timestamp uint64
	Data      *MsgInfoSendPrepareVote
}

//
type MsgInfoSendViewChange struct {
	Msg *sendViewChange
}

type JournalMessageSendViewChange struct {
	Timestamp uint64
	Data      *MsgInfoSendViewChange
}

//
type MsgInfoConfirmedViewChange struct {
	Msg *confirmedViewChange
}

type JournalMessageConfirmedViewChange struct {
	Timestamp uint64
	Data      *MsgInfoConfirmedViewChange
}

func WALDecode(pack []byte, msgType uint16) (*WalMsg, error) {
	switch msgType {
	case SendPrepareBlockMsg:
		var j JournalMessageSendPrepareBlock
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &WalMsg{
				Msg: j.Data.Msg,
			}, nil
		} else {
			return nil, err
		}
	case SendPrepareVoteMsg:
		var j JournalMessageSendPrepareVote
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &WalMsg{
				Msg: j.Data.Msg,
			}, nil
		} else {
			return nil, err
		}
	case SendViewChangeMsg:
		var j JournalMessageSendViewChange
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &WalMsg{
				Msg: j.Data.Msg,
			}, nil
		} else {
			return nil, err
		}
	case ConfirmedViewChangeMsg:
		var j JournalMessageConfirmedViewChange
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &WalMsg{
				Msg: j.Data.Msg,
			}, nil
		} else {
			return nil, err
		}
	}
	panic(fmt.Sprintf("invalid msg type %d", msgType))
}
