package types

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

type ConsensusMsg interface {
	CannibalizeBytes() ([]byte, error)
	Sign() []byte
}

type Message interface {
	String() string
	MsgHash() common.Hash
	BHash() common.Hash
}

type MsgInfo struct {
	Msg    Message
	PeerID discover.NodeID
}
