package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

// simulated message type.
type FakeMessage struct {
	id string
}

func (s *FakeMessage) String() string       { return s.id }
func (s *FakeMessage) MsgHash() common.Hash { return common.Hash{} }
func (s *FakeMessage) BHash() common.Hash   { return common.Hash{} }

func Test_NewMsgInfo(t *testing.T) {
	testCase := []struct {
		pid string
		msg Message
	}{
		{pid: "p01", msg: &FakeMessage{id: "p01"}},
		{pid: "p02", msg: &FakeMessage{id: "p02"}},
	}
	for _, v := range testCase {
		msgInfo := NewMsgInfo(v.msg, v.pid)
		assert.Equal(t, msgInfo.Msg.String(), v.pid)
	}
}

func Test_NewMsgPackage(t *testing.T) {
	testCase := []struct {
		pid  string
		msg  Message
		mode uint64
	}{
		{pid: "p01", msg: &FakeMessage{}, mode: NoneMode},
		{pid: "p02", msg: &FakeMessage{}, mode: FullMode},
		{pid: "p03", msg: &FakeMessage{}, mode: PartMode},
	}
	for _, v := range testCase {
		msgInfo := NewMsgPackage(v.pid, v.msg, v.mode)
		assert.Equal(t, msgInfo.msg, msgInfo.Message())
		assert.Equal(t, msgInfo.peerID, msgInfo.PeerID())
		assert.Equal(t, msgInfo.mode, msgInfo.mode)
	}
}
