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

// Consensus message interface, all consensus message
// types must implement this interface.
type ConsensusMsg interface {
	EpochNum() uint64
	ViewNum() uint64
	BlockNum() uint64
	NodeIndex() uint32
	CannibalizeBytes() ([]byte, error)
	Sign() []byte
	SetSign([]byte)
}

// Message interface, all message structures must
// implement this interface.
type Message interface {
	String() string
	MsgHash() common.Hash
	BHash() common.Hash
}

type MsgInfo struct {
	Msg    Message
	PeerID string
	Inner  bool
}

func (m MsgInfo) String() string {
	return fmt.Sprintf("{peer:%s,type:%s,msg:%s}", m.PeerID, reflect.TypeOf(m.Msg), m.Msg.String())
}

// Create a new MsgInfo object.
func NewMsgInfo(message Message, id string) *MsgInfo {
	return &MsgInfo{
		Msg:    message,
		PeerID: id,
		Inner:  false,
	}
}

func NewInnerMsgInfo(message Message, id string) *MsgInfo {
	return &MsgInfo{
		Msg:    message,
		PeerID: id,
		Inner:  true,
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
