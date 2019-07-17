package types

import (
	"fmt"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const (
	MixMode  = iota // all consensus node
	PartMode        // partial node
	FullMode        // all node
)

// Define error enumeration values related to messages.
const (
	ErrMsgTooLarge = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrCbftProtocolVersionMismatch
	ErrNoStatusMsg
	ErrForkedBlock
)

type ErrCode int

func (e ErrCode) String() string {
	return errorToString[int(e)]
}

// Error code mapping error message.
var errorToString = map[int]string{
	ErrMsgTooLarge:                 "Message too long",
	ErrDecode:                      "Invalid message",
	ErrInvalidMsgCode:              "Invalid message code",
	ErrCbftProtocolVersionMismatch: "CBFT Protocol version mismatch",
	ErrNoStatusMsg:                 "No status message",
	ErrForkedBlock:                 "Forked block",
}

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

// MsgPackage represents a specific message package.
// It contains the node ID, the message body, and
// the forwarding mode from the sender.
type MsgPackage struct {
	peerID string  // from the sender of the message
	msg    Message // message body
	mode   uint64  // forwarding mode.
}

func (m *MsgPackage) Message() Message {
	return m.msg
}

func (m *MsgPackage) PeerID() string {
	return m.peerID
}

func (m *MsgPackage) Mode() uint64 {
	return m.mode
}

func (m *MsgPackage) MessageType() uint64 {
	return messageType(m.msg)
}

// A is used to convert specific message types according to the message body.
// The program is forcibly terminated if there is an unmatched message type and
// all types must exist in the match list.
func messageType(msg interface{}) uint64 {
	switch msg.(type) {
	default:
		return protocols.PrepareBlockHashMsg
	}
	panic(fmt.Sprintf("unknown message type [%v]", reflect.TypeOf(msg)))
}
