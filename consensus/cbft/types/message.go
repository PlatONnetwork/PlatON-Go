package types

import (
	"fmt"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

const (
	NoneMode = iota // none consensus node
	PartMode        // partial node
	FullMode        // all node
)

// Define error enumeration values related to messages.
const (
	ErrMsgTooLarge = iota
	ErrExtraStatusMsg
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
	ErrForkedBlock:                 "Forked Block",
}

// Build an error object based on the error code.
func ErrResp(code ErrCode, format string, v ...interface{}) error {
	return fmt.Errorf("%v - %v", code, fmt.Sprintf(format, v...))
}

type ConsensusMsg interface {
	EpochNum() uint64
	ViewNum() uint64
	BlockNum() uint64
	NodeIndex() uint32
	CannibalizeBytes() ([]byte, error)
	Sign() []byte
	SetSign([]byte)
}

type Message interface {
	String() string
	MsgHash() common.Hash
	BHash() common.Hash
}

type MsgInfo struct {
	Msg    Message
	PeerID string
}

func (m MsgInfo) String() string {
	return fmt.Sprintf("{peer:%s,type:%s,msg:%s}", m.PeerID, reflect.TypeOf(m.Msg), m.Msg.String())
}

// Create a new MsgInfo object.
func NewMsgInfo(message Message, id string) *MsgInfo {
	return &MsgInfo{
		Msg:    message,
		PeerID: id,
	}
}

// MsgPackage represents a specific message package.
// It contains the node ID, the message body, and
// the forwarding mode from the sender.
type MsgPackage struct {
	peerID string  // from the sender of the message
	msg    Message // message body
	mode   uint64  // forwarding mode.
}

// Create a new MsgPackage based on params.
func NewMsgPackage(pid string, msg Message, mode uint64) *MsgPackage {
	return &MsgPackage{
		peerID: pid,
		msg:    msg,
		mode:   mode,
	}
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
