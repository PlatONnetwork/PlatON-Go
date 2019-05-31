package cbft

import (
	"fmt"
	_ "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/prepare"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type fakeMetricsRW struct {
	msg p2p.Msg
}

func (rw *fakeMetricsRW) ReadMsg() (p2p.Msg, error) {
	return rw.msg, nil
}

func (rw *fakeMetricsRW) WriteMsg(msg p2p.Msg) error {
	fmt.Println("FakeMetricsRW write msg")
	return nil
}

func TestMeteredMsgReadWriter_ReadMsg(t *testing.T) {
	readMsg := func(code uint64, size uint32) {
		metricsRW := newMeteredMsgWriter(&fakeMetricsRW{
			msg: p2p.Msg{ Code: code, Size: size, },
		})
		metricsRW.ReadMsg()
	}
	testCases := []struct{
		code uint64
		size uint32
		want int64
	}{
		{PrepareBlockMsg, 100, 100},
		{ViewChangeMsg, 111, 111},
		{ViewChangeVoteMsg, 121, 121},
		{PrepareVoteMsg, 131, 131},
		{ConfirmedPrepareBlockMsg, 144, 144},
		{PrepareVotesMsg, 1, 1},
		{HighestPrepareBlockMsg, 2, 2},
		{PrepareBlockHashMsg, 3, 3},
		{CBFTStatusMsg, 3, 3},
	}
	for _, v := range testCases {
		readMsg(v.code, v.size)
		time.Sleep(100)
		switch {
		case v.code == PrepareBlockMsg:
			assert.NotEqual(t, int64(0), propPrepareBlockInPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareBlockInTrafficMeter.Count())
		case v.code == ViewChangeMsg:
			assert.NotEqual(t,int64(0), propViewChangeInPacketsMeter.Count(),)
			assert.Equal(t, v.want, propViewChangeInTrafficMeter.Count())
		case v.code == ViewChangeVoteMsg:
			assert.NotEqual(t,int64(0), propViewChangeVoteInPacketsMeter.Count())
			assert.Equal(t, v.want, propViewChangeVoteInTrafficMeter.Count())
		case v.code == PrepareVoteMsg:
			assert.NotEqual(t,int64(0), propPrepareVoteInPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareVoteInTrafficMeter.Count())
		case v.code == ConfirmedPrepareBlockMsg:
			assert.NotEqual(t,int64(0), propConfirmedPrepareBlockInPacketsMeter.Count())
			assert.Equal(t, v.want, propConfirmedPrepareBlockInTrafficMeter.Count())
		case v.code == PrepareVotesMsg:
			assert.NotEqual(t,int64(0), reqPrepareVotesInPacketsMeter.Count())
			assert.Equal(t, v.want, reqPrepareVotesInTrafficMeter.Count())
		case v.code == HighestPrepareBlockMsg:
			assert.NotEqual(t,int64(0), reqHighestPrepareBlockInPacketsMeter.Count())
			assert.Equal(t, v.want, reqHighestPrepareBlockInTrafficMeter.Count())
		case v.code == PrepareBlockHashMsg:
			assert.NotEqual(t,int64(0), propPrepareBlockHashInPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareBlockHashInTrafficMeter.Count())
		case v.code == CBFTStatusMsg:
			assert.NotEqual(t,int64(0), miscInPacketsMeter.Count())
			assert.Equal(t, v.want, miscInTrafficMeter.Count())
		}
	}
}

func TestMeteredMsgReadWriter_WriteMsg(t *testing.T) {
	writeMsg := func(code uint64, size uint32) {
		metricsRW := newMeteredMsgWriter(&fakeMetricsRW{})
		metricsRW.WriteMsg(p2p.Msg{ Code: code, Size: size, })
	}
	testCases := []struct{
		code uint64
		size uint32
		want int64
	}{
		{PrepareBlockMsg, 100, 100},
		{ViewChangeMsg, 111, 111},
		{ViewChangeVoteMsg, 121, 121},
		{PrepareVoteMsg, 131, 131},
		{ConfirmedPrepareBlockMsg, 144, 144},
		{PrepareVotesMsg, 1, 1},
		{HighestPrepareBlockMsg, 2, 2},
		{PrepareBlockHashMsg, 3, 3},
		{CBFTStatusMsg, 3, 3},
	}
	for _, v := range testCases {
		writeMsg(v.code, v.size)
		switch {
		case v.code == PrepareBlockMsg:
			assert.NotEqual(t, int64(0), propPrepareBlockOutPacketsMeter.Count())
			assert.NotEqual(t, int64(0), propPrepareBlockOutTrafficMeter.Count())
		case v.code == ViewChangeMsg:
			assert.NotEqual(t, int64(0), propViewChangeOutPacketsMeter.Count(),)
			assert.NotEqual(t, int64(0), propViewChangeOutTrafficMeter.Count())
		case v.code == ViewChangeVoteMsg:
			assert.NotEqual(t, int64(0), propViewChangeVoteOutPacketsMeter.Count())
			assert.NotEqual(t, int64(0), propViewChangeVoteOutTrafficMeter.Count())
		case v.code == PrepareVoteMsg:
			assert.NotEqual(t, int64(0), propPrepareVoteOutPacketsMeter.Count())
			assert.NotEqual(t, int64(0), propPrepareVoteOutTrafficMeter.Count())
		case v.code == ConfirmedPrepareBlockMsg:
			assert.NotEqual(t, int64(0), propConfirmedPrepareBlockOutPacketsMeter.Count())
			assert.NotEqual(t, int64(0), propConfirmedPrepareBlockOutTrafficMeter.Count())
		case v.code == PrepareVotesMsg:
			assert.NotEqual(t, int64(0), reqPrepareVotesOutPacketsMeter.Count())
			assert.NotEqual(t, int64(0), reqPrepareVotesOutTrafficMeter.Count())
		case v.code == HighestPrepareBlockMsg:
			assert.NotEqual(t, int64(0), reqHighestPrepareBlockOutPacketsMeter.Count())
			assert.NotEqual(t, int64(0), reqHighestPrepareBlockOutTrafficMeter.Count())
		case v.code == PrepareBlockHashMsg:
			assert.NotEqual(t, int64(0), propPrepareBlockHashOutPacketsMeter.Count())
			assert.NotEqual(t, int64(0), propPrepareBlockHashOutTrafficMeter.Count())
		case v.code == CBFTStatusMsg:
			assert.NotEqual(t, int64(0), miscOutPacketsMeter.Count())
			assert.NotEqual(t, int64(0), miscOutTrafficMeter.Count())
		}

	}
}