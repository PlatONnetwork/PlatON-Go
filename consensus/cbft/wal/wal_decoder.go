// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package wal

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// struct SendPrepareBlock for rlp decode
type MessageSendPrepareBlock struct {
	Timestamp uint64
	Data      *protocols.SendPrepareBlock
}

// struct SendPrepareVote for rlp decode
type MessageSendPrepareVote struct {
	Timestamp uint64
	Data      *protocols.SendPrepareVote
}

// struct SendViewChange for rlp decode
type MessageSendViewChange struct {
	Timestamp uint64
	Data      *protocols.SendViewChange
}

// struct ConfirmedViewChange for rlp decode
type MessageConfirmedViewChange struct {
	Timestamp uint64
	Data      *protocols.ConfirmedViewChange
}

func WALDecode(pack []byte, msgType uint16) (interface{}, error) {
	switch msgType {
	case protocols.ConfirmedViewChangeMsg:
		var j MessageConfirmedViewChange
		if err := rlp.DecodeBytes(pack, &j); err != nil {
			return nil, err

		}
		return j.Data, nil

	case protocols.SendViewChangeMsg:
		var j MessageSendViewChange
		if err := rlp.DecodeBytes(pack, &j); err != nil {
			return nil, err

		}
		return j.Data, nil

	case protocols.SendPrepareBlockMsg:
		var j MessageSendPrepareBlock
		if err := rlp.DecodeBytes(pack, &j); err != nil {
			return nil, err
		}
		return j.Data, nil

	case protocols.SendPrepareVoteMsg:
		var j MessageSendPrepareVote
		if err := rlp.DecodeBytes(pack, &j); err != nil {
			return nil, err

		}
		return j.Data, nil
	}
	panic(fmt.Sprintf("invalid msg type %d", msgType))
}
