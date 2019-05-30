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
			assert.Equal(t, int64(1), propPrepareBlockInPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareBlockInTrafficMeter.Count())
		case v.code == ViewChangeMsg:
			assert.Equal(t,int64(1), propViewChangeInPacketsMeter.Count(),)
			assert.Equal(t, v.want, propViewChangeInTrafficMeter.Count())
		case v.code == ViewChangeVoteMsg:
			assert.Equal(t,int64(1), propViewChangeVoteInPacketsMeter.Count())
			assert.Equal(t, v.want, propViewChangeVoteInTrafficMeter.Count())
		case v.code == PrepareVoteMsg:
			assert.Equal(t,int64(1), propPrepareVoteInPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareVoteInTrafficMeter.Count())
		case v.code == ConfirmedPrepareBlockMsg:
			assert.Equal(t,int64(1), propConfirmedPrepareBlockInPacketsMeter.Count())
			assert.Equal(t, v.want, propConfirmedPrepareBlockInTrafficMeter.Count())
		case v.code == PrepareVotesMsg:
			assert.Equal(t,int64(1), reqPrepareVotesInPacketsMeter.Count())
			assert.Equal(t, v.want, reqPrepareVotesInTrafficMeter.Count())
		case v.code == HighestPrepareBlockMsg:
			assert.Equal(t,int64(1), reqHighestPrepareBlockInPacketsMeter.Count())
			assert.Equal(t, v.want, reqHighestPrepareBlockInTrafficMeter.Count())
		case v.code == PrepareBlockHashMsg:
			assert.Equal(t,int64(1), propPrepareBlockHashInPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareBlockHashInTrafficMeter.Count())
		case v.code == CBFTStatusMsg:
			assert.Equal(t,int64(1), miscInPacketsMeter.Count())
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
			assert.Equal(t, int64(1), propPrepareBlockOutPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareBlockOutTrafficMeter.Count())
		case v.code == ViewChangeMsg:
			assert.Equal(t,int64(1), propViewChangeOutPacketsMeter.Count(),)
			assert.Equal(t, v.want, propViewChangeOutTrafficMeter.Count())
		case v.code == ViewChangeVoteMsg:
			assert.Equal(t,int64(1), propViewChangeVoteOutPacketsMeter.Count())
			assert.Equal(t, v.want, propViewChangeVoteOutTrafficMeter.Count())
		case v.code == PrepareVoteMsg:
			assert.Equal(t,int64(1), propPrepareVoteOutPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareVoteOutTrafficMeter.Count())
		case v.code == ConfirmedPrepareBlockMsg:
			assert.Equal(t,int64(1), propConfirmedPrepareBlockOutPacketsMeter.Count())
			assert.Equal(t, v.want, propConfirmedPrepareBlockOutTrafficMeter.Count())
		case v.code == PrepareVotesMsg:
			assert.Equal(t,int64(1), reqPrepareVotesOutPacketsMeter.Count())
			assert.Equal(t, v.want, reqPrepareVotesOutTrafficMeter.Count())
		case v.code == HighestPrepareBlockMsg:
			assert.Equal(t,int64(1), reqHighestPrepareBlockOutPacketsMeter.Count())
			assert.Equal(t, v.want, reqHighestPrepareBlockOutTrafficMeter.Count())
		case v.code == PrepareBlockHashMsg:
			assert.Equal(t,int64(1), propPrepareBlockHashOutPacketsMeter.Count())
			assert.Equal(t, v.want, propPrepareBlockHashOutTrafficMeter.Count())
		case v.code == CBFTStatusMsg:
			assert.Equal(t,int64(1), miscOutPacketsMeter.Count())
			assert.Equal(t, v.want, miscOutTrafficMeter.Count())
		}

	}
}