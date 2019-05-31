package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

type MsgInfoPrepareBlock struct {
	Msg    *prepareBlock
	PeerID discover.NodeID
}

type JournalMessagePrepareBlock struct {
	Timestamp uint64
	Data      *MsgInfoPrepareBlock
}

type MsgInfoPrepareVote struct {
	Msg    *prepareVote
	PeerID discover.NodeID
}

type JournalMessagePrepareVote struct {
	Timestamp uint64
	Data      *MsgInfoPrepareVote
}

type MsgInfoViewChange struct {
	Msg    *viewChange
	PeerID discover.NodeID
}

type JournalMessageViewChange struct {
	Timestamp uint64
	Data      *MsgInfoViewChange
}

type MsgInfoViewChangeVote struct {
	Msg    *viewChangeVote
	PeerID discover.NodeID
}

type JournalMessageViewChangeVote struct {
	Timestamp uint64
	Data      *MsgInfoViewChangeVote
}

type MsgInfoConfirmedPrepareBlock struct {
	Msg    *confirmedPrepareBlock
	PeerID discover.NodeID
}

type JournalMessageConfirmedPrepareBlock struct {
	Timestamp uint64
	Data      *MsgInfoConfirmedPrepareBlock
}

type MsgInfoGetPrepareVote struct {
	Msg    *getPrepareVote
	PeerID discover.NodeID
}

type JournalMessageGetPrepareVote struct {
	Timestamp uint64
	Data      *MsgInfoGetPrepareVote
}

type MsgInfoPrepareVotes struct {
	Msg    *prepareVotes
	PeerID discover.NodeID
}

type JournalMessagePrepareVotes struct {
	Timestamp uint64
	Data      *MsgInfoPrepareVotes
}

type MsgInfoGetPrepareBlock struct {
	Msg    *getPrepareBlock
	PeerID discover.NodeID
}

type JournalMessageGetPrepareBlock struct {
	Timestamp uint64
	Data      *MsgInfoGetPrepareBlock
}

type MsgInfoGetHighestPrepareBlock struct {
	Msg    *getHighestPrepareBlock
	PeerID discover.NodeID
}

type JournalMessageGetHighestPrepareBlock struct {
	Timestamp uint64
	Data      *MsgInfoGetHighestPrepareBlock
}

type MsgInfoHighestPrepareBlock struct {
	Msg    *highestPrepareBlock
	PeerID discover.NodeID
}

type JournalMessageHighestPrepareBlock struct {
	Timestamp uint64
	Data      *MsgInfoHighestPrepareBlock
}

type MsgInfoCbftStatusData struct {
	Msg    *cbftStatusData
	PeerID discover.NodeID
}

type JournalMessageCbftStatusData struct {
	Timestamp uint64
	Data      *MsgInfoCbftStatusData
}

type MsgInfoPrepareBlockHash struct {
	Msg    *prepareBlockHash
	PeerID discover.NodeID
}

type JournalMessagePrepareBlockHash struct {
	Timestamp uint64
	Data      *MsgInfoPrepareBlockHash
}

func WALDecode(pack []byte, msgType uint16) (*MsgInfo, error) {
	switch msgType {
	case PrepareBlockMsg:
		var j JournalMessagePrepareBlock
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}

	case PrepareVoteMsg:
		var j JournalMessagePrepareVote
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}

	case ViewChangeMsg:
		var j JournalMessageViewChange
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}

	case ViewChangeVoteMsg:
		var j JournalMessageViewChangeVote
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}
	case ConfirmedPrepareBlockMsg:
		var j JournalMessageConfirmedPrepareBlock
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}
	case GetPrepareVoteMsg:
		var j JournalMessageGetPrepareVote
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}
	case PrepareVotesMsg:
		var j JournalMessagePrepareVotes
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}
	case GetPrepareBlockMsg:
		var j JournalMessageGetPrepareBlock
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}
	case GetHighestPrepareBlockMsg:
		var j JournalMessageGetHighestPrepareBlock
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}
	case HighestPrepareBlockMsg:
		var j JournalMessageHighestPrepareBlock
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}
	case CBFTStatusMsg:
		var j JournalMessageCbftStatusData
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}
	case PrepareBlockHashMsg:
		var j JournalMessagePrepareBlockHash
		if err := rlp.DecodeBytes(pack, &j); err == nil {
			return &MsgInfo{
				Msg:    j.Data.Msg,
				PeerID: j.Data.PeerID,
			}, nil
		} else {
			return nil, err
		}
	}
	panic(fmt.Sprintf("invalid msg type %d", msgType))
}
