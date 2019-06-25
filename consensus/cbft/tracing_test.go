package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestReceiveRecord(t *testing.T) {
	testCases := []*receiveRecord{
		{
			Time: time.Now(), SelfId: "001", FromId: "002", MsgHash: common.BytesToHash([]byte("Hash")).TerminalString(), Type: "prepareBlock",
		},
	}
	for _, v := range testCases {
		if json, err := v.ToJSON(); err != nil {
			t.Error("marshalJSON fail", "err", err)
		} else {
			t.Log(string(json))
			var tmp receiveRecord
			err := tmp.ParseFromJSON(json)
			assert.Nil(t, err)
		}
	}
}

func TestSendRecord(t *testing.T) {
	testCases := []*sendRecord{
		{
			Time: time.Now(), SelfId: "001", TargetIds: []string{"002"}, MsgHash: common.BytesToHash([]byte("Hash")).TerminalString(), Type: "prepareBlock",
		},
	}
	for _, v := range testCases {
		if json, err := v.ToJSON(); err != nil {
			t.Error("marshalJSON fail", "err", err)
		} else {
			t.Log(string(json))
			var tmp sendRecord
			err := tmp.ParseFromJSON(json)
			assert.Nil(t, err)
		}
	}
}

func TestNewTracing(t *testing.T) {
	trac := NewTracing()
	receiveUseCases := []*receiveRecord{
		{ Time: time.Now(), Order: 1, MsgHash: "01-hash",},
		{ Time: time.Now(), Order: 3, MsgHash: "03-hash",},
		{ Time: time.Now(), Order: 2, MsgHash: "02-hash",},
	}
	sendUseCases := []*sendRecord {
		{ Time:time.Now(), Order: 3, MsgHash: "03-hash", TargetIds: []string{"1","2","3"}, },
		{ Time:time.Now(), Order: 1, MsgHash: "01-hash", TargetIds: []string{"1","2","3"},},
		{ Time:time.Now(), Order: 2, MsgHash: "02-hash", TargetIds: []string{"1","2","3"},},
	}
	trac.On()
	for _, v := range receiveUseCases {
		trac.recordReceiveAction(v)
	}
	for _, v := range sendUseCases {
		trac.recordSendAction(v)
	}
	t.Log(trac)
}

func TestStart(t *testing.T) {
	tracing := NewTracing()
	tracing.On()
	time.Sleep(time.Second * 2)
	tracing.Off()
}