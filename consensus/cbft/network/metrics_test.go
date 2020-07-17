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

package network

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/metrics"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"

	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/stretchr/testify/assert"
)

type fakeMetricsRW struct {
	msg p2p.Msg
}

func (rw *fakeMetricsRW) ReadMsg() (p2p.Msg, error) { return rw.msg, nil }

func (rw *fakeMetricsRW) WriteMsg(msg p2p.Msg) error { return nil }

func Test_MeteredMsgReadWriter_ReadMsg(t *testing.T) {
	// cmd: go test -v  -run Test_MeteredMsgReadWriter_ReadMsg ./ --args "metrics"
	if !metrics.Enabled {
		t.Log("metrics disabled, abort unit test")
		return
	}
	readMsg := func(code uint64, size uint32) {
		metricsRW := newMeteredMsgWriter(&fakeMetricsRW{
			msg: p2p.Msg{Code: code, Size: size},
		})
		metricsRW.ReadMsg()
	}
	testCases := []struct {
		code uint64
		size uint32
		want int64
	}{
		{protocols.PrepareBlockMsg, 100, 100},
		{protocols.ViewChangeMsg, 111, 111},
		{protocols.PrepareVoteMsg, 131, 131},
		{protocols.GetPrepareBlockMsg, 1, 1},
		{protocols.GetBlockQuorumCertMsg, 3, 3},
		{protocols.BlockQuorumCertMsg, 3, 3},
		{protocols.GetPrepareVoteMsg, 3, 3},
		{protocols.PrepareVotesMsg, 3, 3},
		{protocols.GetQCBlockListMsg, 3, 3},
		{protocols.QCBlockListMsg, 3, 3},
		{protocols.PrepareBlockHashMsg, 3, 3},
	}
	for _, v := range testCases {
		readMsg(v.code, v.size)
		switch {
		case v.code == protocols.PrepareBlockMsg:
			assert.NotEqual(t, 0, propPrepareBlockInPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareBlockInTrafficMeter.Count())

		case v.code == protocols.PrepareVoteMsg:
			assert.NotEqual(t, 0, propPrepareVoteInPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareVoteInTrafficMeter.Count())

		case v.code == protocols.ViewChangeMsg:
			assert.NotEqual(t, 0, propViewChangeInPacketsMeter.Count())
			assert.Equal(t, v.want, propViewChangeInTrafficMeter.Count())

		case v.code == protocols.GetPrepareBlockMsg:
			assert.NotEqual(t, 0, reqGetPrepareBlockInPacketsMeter.Count())
			assert.Equal(t, v.want, reqGetPrepareBlockInTrafficMeter.Count())

		case v.code == protocols.GetBlockQuorumCertMsg:
			assert.NotEqual(t, 0, reqGetQuorumCertInPacketsMeter.Count())
			assert.Equal(t, v.want, reqGetQuorumCertInTrafficMeter.Count())

		case v.code == protocols.BlockQuorumCertMsg:
			assert.NotEqual(t, 0, reqBlockQuorumCertInPacketsMeter.Count())
			assert.Equal(t, v.want, reqBlockQuorumCertInTrafficMeter.Count())

		case v.code == protocols.GetPrepareVoteMsg:
			assert.NotEqual(t, 0, reqGetPrepareVoteInPacketsMeter.Count())
			assert.Equal(t, v.want, reqGetPrepareVoteInTrafficMeter.Count())

		case v.code == protocols.PrepareVotesMsg:
			assert.NotEqual(t, 0, reqPrepareVotesInPacketsMeter.Count())
			assert.Equal(t, v.want, reqPrepareVotesInTrafficMeter.Count())

		case v.code == protocols.GetQCBlockListMsg:
			assert.NotEqual(t, 0, reqGetQCBlockListInPacketsMeter.Count())
			assert.Equal(t, v.want, reqGetQCBlockListInTrafficMeter.Count())

		case v.code == protocols.QCBlockListMsg:
			assert.NotEqual(t, 0, reqQCBlockListInPacketsMeter.Count())
			assert.Equal(t, v.want, reqGetQCBlockListInTrafficMeter.Count())
		}
	}
}

func TestMeteredMsgReadWriter_WriteMsg(t *testing.T) {
	// cmd: go test -v  -run TestMeteredMsgReadWriter_WriteMsg ./ --args "metrics"
	if !metrics.Enabled {
		t.Log("metrics disabled, abort unit test")
		return
	}
	writeMsg := func(code uint64, size uint32) {
		metricsRW := newMeteredMsgWriter(&fakeMetricsRW{})
		metricsRW.WriteMsg(p2p.Msg{Code: code, Size: size})
	}
	testCases := []struct {
		code uint64
		size uint32
		want int64
	}{
		{protocols.PrepareBlockMsg, 100, 100},
		{protocols.ViewChangeMsg, 111, 111},
		{protocols.PrepareVoteMsg, 131, 131},
		{protocols.GetPrepareBlockMsg, 1, 1},
		{protocols.GetBlockQuorumCertMsg, 3, 3},
		{protocols.BlockQuorumCertMsg, 3, 3},
		{protocols.GetPrepareVoteMsg, 3, 3},
		{protocols.PrepareVotesMsg, 3, 3},
		{protocols.GetQCBlockListMsg, 3, 3},
		{protocols.QCBlockListMsg, 3, 3},
		{protocols.PrepareBlockHashMsg, 3, 3},
	}

	for _, v := range testCases {
		writeMsg(v.code, v.size)
		switch {
		case v.code == protocols.PrepareBlockMsg:
			assert.NotEqual(t, 0, propPrepareBlockOutPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareBlockOutTrafficMeter.Count())

		case v.code == protocols.PrepareVoteMsg:
			assert.NotEqual(t, 0, propPrepareVoteOutPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareVoteOutTrafficMeter.Count())

		case v.code == protocols.ViewChangeMsg:
			assert.NotEqual(t, 0, propViewChangeOutPacketsMeter.Count())
			assert.Equal(t, v.want, propViewChangeOutTrafficMeter.Count())

		case v.code == protocols.GetPrepareBlockMsg:
			assert.NotEqual(t, 0, reqGetPrepareBlockOutPacketsMeter.Count())
			assert.Equal(t, v.want, reqGetPrepareBlockOutTrafficMeter.Count())

		case v.code == protocols.GetBlockQuorumCertMsg:
			assert.NotEqual(t, 0, reqGetQuorumCertOutPacketsMeter.Count())
			assert.Equal(t, v.want, reqGetQuorumCertOutTrafficMeter.Count())

		case v.code == protocols.BlockQuorumCertMsg:
			assert.NotEqual(t, 0, reqBlockQuorumCertOutPacketsMeter.Count())
			assert.Equal(t, v.want, reqBlockQuorumCertOutTrafficMeter.Count())

		case v.code == protocols.GetPrepareVoteMsg:
			assert.NotEqual(t, 0, reqGetPrepareVoteOutPacketsMeter.Count())
			assert.Equal(t, v.want, reqGetPrepareVoteOutTrafficMeter.Count())

		case v.code == protocols.PrepareVotesMsg:
			assert.NotEqual(t, 0, reqPrepareVotesOutPacketsMeter.Count())
			assert.Equal(t, v.want, reqPrepareVotesOutTrafficMeter.Count())

		case v.code == protocols.GetQCBlockListMsg:
			assert.NotEqual(t, 0, reqGetQCBlockListOutPacketsMeter.Count())
			assert.Equal(t, v.want, reqGetQCBlockListOutTrafficMeter.Count())

		case v.code == protocols.QCBlockListMsg:
			assert.NotEqual(t, 0, reqQCBlockListOutPacketsMeter.Count())
			assert.Equal(t, v.want, reqQCBlockListOutTrafficMeter.Count())
		}
	}
}
