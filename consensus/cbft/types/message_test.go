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
		assert.NotEmpty(t, msgInfo.String())
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
		assert.Equal(t, msgInfo.mode, msgInfo.Mode())
		msgInfo.Mode()
	}
}

func Test_ErrCode(t *testing.T) {
	var errCode ErrCode
	errCode = ErrMsgTooLarge
	assert.Equal(t, errorToString[ErrMsgTooLarge], errCode.String())
	ErrResp(errCode, "%s", errCode.String())
}
