package network

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
)

var (

	// for the message
	propPrepareBlockInPacketsMeter  = metrics.NewRegisteredMeter("cbft/prop/prepare_blocks/in/packets", nil)
	propPrepareBlockInTrafficMeter  = metrics.NewRegisteredMeter("cbft/prop/prepare_blocks/in/traffic", nil)
	propPrepareBlockOutPacketsMeter = metrics.NewRegisteredMeter("cbft/prop/prepare_blocks/out/packets", nil)
	propPrepareBlockOutTrafficMeter = metrics.NewRegisteredMeter("cbft/prop/prepare_blocks/out/traffic", nil)

	propViewChangeInPacketsMeter  = metrics.NewRegisteredMeter("cbft/prop/view_changes/in/packets", nil)
	propViewChangeInTrafficMeter  = metrics.NewRegisteredMeter("cbft/prop/view_changes/in/traffic", nil)
	propViewChangeOutPacketsMeter = metrics.NewRegisteredMeter("cbft/prop/view_changes/out/packets", nil)
	propViewChangeOutTrafficMeter = metrics.NewRegisteredMeter("cbft/prop/view_changes/out/traffic", nil)

	propViewChangeVoteInPacketsMeter  = metrics.NewRegisteredMeter("cbft/prop/view_change_votes/in/packets", nil)
	propViewChangeVoteInTrafficMeter  = metrics.NewRegisteredMeter("cbft/prop/view_change_votes/in/traffic", nil)
	propViewChangeVoteOutPacketsMeter = metrics.NewRegisteredMeter("cbft/prop/view_change_votes/out/packets", nil)
	propViewChangeVoteOutTrafficMeter = metrics.NewRegisteredMeter("cbft/prop/view_change_votes/out/traffic", nil)

	propPrepareVoteInPacketsMeter  = metrics.NewRegisteredMeter("cbft/prop/prepare_votes/in/packets", nil)
	propPrepareVoteInTrafficMeter  = metrics.NewRegisteredMeter("cbft/prop/prepare_votes/in/traffic", nil)
	propPrepareVoteOutPacketsMeter = metrics.NewRegisteredMeter("cbft/prop/prepare_votes/out/packets", nil)
	propPrepareVoteOutTrafficMeter = metrics.NewRegisteredMeter("cbft/prop/prepare_votes/out/traffic", nil)

	propConfirmedPrepareBlockInPacketsMeter  = metrics.NewRegisteredMeter("cbft/prop/prepare_votes/in/packets", nil)
	propConfirmedPrepareBlockInTrafficMeter  = metrics.NewRegisteredMeter("cbft/prop/prepare_votes/in/traffic", nil)
	propConfirmedPrepareBlockOutPacketsMeter = metrics.NewRegisteredMeter("cbft/prop/prepare_votes/out/packets", nil)
	propConfirmedPrepareBlockOutTrafficMeter = metrics.NewRegisteredMeter("cbft/prop/prepare_votes/out/traffic", nil)

	reqPrepareVotesInPacketsMeter  = metrics.NewRegisteredMeter("cbft/req/prepare_votes/in/packets", nil)
	reqPrepareVotesInTrafficMeter  = metrics.NewRegisteredMeter("cbft/req/prepare_votes/in/traffic", nil)
	reqPrepareVotesOutPacketsMeter = metrics.NewRegisteredMeter("cbft/req/prepare_votes/out/packets", nil)
	reqPrepareVotesOutTrafficMeter = metrics.NewRegisteredMeter("cbft/req/prepare_votes/out/traffic", nil)

	reqHighestPrepareBlockInPacketsMeter  = metrics.NewRegisteredMeter("cbft/req/highest_prepare_blocks/in/packets", nil)
	reqHighestPrepareBlockInTrafficMeter  = metrics.NewRegisteredMeter("cbft/req/highest_prepare_blocks/in/traffic", nil)
	reqHighestPrepareBlockOutPacketsMeter = metrics.NewRegisteredMeter("cbft/req/highest_prepare_blocks/out/packets", nil)
	reqHighestPrepareBlockOutTrafficMeter = metrics.NewRegisteredMeter("cbft/req/highest_prepare_blocks/out/traffic", nil)

	propPrepareBlockHashInPacketsMeter  = metrics.NewRegisteredMeter("cbft/prop/prepare_block_hashes/in/packets", nil)
	propPrepareBlockHashInTrafficMeter  = metrics.NewRegisteredMeter("cbft/prop/prepare_block_hashes/in/traffic", nil)
	propPrepareBlockHashOutPacketsMeter = metrics.NewRegisteredMeter("cbft/prop/prepare_block_hashes/out/packets", nil)
	propPrepareBlockHashOutTrafficMeter = metrics.NewRegisteredMeter("cbft/prop/prepare_block_hashes/out/traffic", nil)

	miscInPacketsMeter  = metrics.NewRegisteredMeter("cbft/misc/in/packets", nil)
	miscInTrafficMeter  = metrics.NewRegisteredMeter("cbft/misc/in/traffic", nil)
	miscOutPacketsMeter = metrics.NewRegisteredMeter("cbft/misc/out/packets", nil)
	miscOutTrafficMeter = metrics.NewRegisteredMeter("cbft/misc/out/traffic", nil)

	// for the consensus of cbft
	//blockConfirmedCountFulfillTimer		= metrics.NewRegisteredTimer("cbft/blocks/count/fulfill", nil)
	blockMinedTimer = metrics.NewRegisteredTimer("cbft/blocks/mined", nil)
	//blockImportTimer					= metrics.NewRegisteredTimer("cbft/blocks/imported", nil)
	blockExecuteTimer          = metrics.NewRegisteredTimer("cbft/blocks/execute", nil)
	blockConfirmedTimer        = metrics.NewRegisteredTimer("cbft/blocks/confirmed", nil)
	viewChangeTimer            = metrics.NewRegisteredTimer("cbft/views/change", nil)
	viewChangeConfirmedTimer   = metrics.NewRegisteredTimer("cbft/views/confirm", nil)
	viewChangeVoteFulfillTimer = metrics.NewRegisteredTimer("cbft/views/count/fulfill", nil)

	messageGossipMeter            = metrics.NewRegisteredMeter("cbft/meter/message/gossip", nil)
	messageRepeatMeter            = metrics.NewRegisteredMeter("cbft/meter/message/repeat", nil)
	blockMinedMeter               = metrics.NewRegisteredMeter("cbft/meter/blocks/mined", nil)
	blockVerifyFailMeter          = metrics.NewRegisteredMeter("cbft/meter/blocks/verify/fail", nil)
	signatureVerifyFailMeter      = metrics.NewRegisteredMeter("cbft/meter/signature/verify/fail", nil)
	viewChangeVoteVerifyFailMeter = metrics.NewRegisteredMeter("cbft/meter/view_change_votes/verify/fail", nil)
	blockConfirmedMeter           = metrics.NewRegisteredMeter("cbft/meter/blocks/confirmed", nil)
	//blockMissMeter						= metrics.NewRegisteredMeter("cbft/meter/blocks/miss", nil)
	viewChangeTimeoutMeter     = metrics.NewRegisteredMeter("cbft/meter/view/view_changes/timeout", nil)
	viewChangeVoteTimeoutMeter = metrics.NewRegisteredMeter("cbft/meter/view/view_change_votes/timeout", nil)

	viewChangeCounter = metrics.NewRegisteredCounter("cbft/counter/view_changes/count", nil)
	//blockMinedCountCounter				= metrics.NewRegisteredCounter("cbft/counter/blocks/mined", nil)		//  The number of blocks in a round of views
	consensusJoinCounter = metrics.NewRegisteredCounter("cbft/counter/consensus/join", nil) //

	//blockUnconfirmedGauage				= metrics.NewRegisteredGauge("cbft/gauage/blocks/unconfirmed", nil)
	blockHighNumConfirmedGauage = metrics.NewRegisteredGauge("cbft/gauage/blocks/height_num/confirmed", nil)
	blockHighNumLogicGauage     = metrics.NewRegisteredGauge("cbft/gauage/blocks/height_num/logic", nil)
	viewChangeGauage            = metrics.NewRegisteredGauge("cbft/gauage/views/viewchange", nil)
)

// meteredMsgReadWriter is a wrapper around a p2p.MsgReadWriter, capable of
// accumulating the above defined metrics based on the data stream contents.
type meteredMsgReadWriter struct {
	p2p.MsgReadWriter     // Wrapped message stream to meter
	version           int // Protocol version to select correct meters
}

// newMeteredMsgWriter wraps a p2p MsgReadWriter with metering support. If the
// metrics system is disabled, this function returns the original object.
func newMeteredMsgWriter(rw p2p.MsgReadWriter) p2p.MsgReadWriter {
	if !metrics.Enabled {
		return rw
	}
	return &meteredMsgReadWriter{MsgReadWriter: rw}
}

// Init sets the protocol version used by the stream to know which meters to
// increment in case of overlapping message ids between protocol versions.
func (rw *meteredMsgReadWriter) Init(version int) {
	rw.version = version
}

func (rw *meteredMsgReadWriter) ReadMsg() (p2p.Msg, error) {
	// Read the message and short circuit in case of an error
	msg, err := rw.MsgReadWriter.ReadMsg()
	if err != nil {
		return msg, err
	}
	packets, traffic := miscInPacketsMeter, miscInTrafficMeter
	switch {
	case msg.Code == protocols.PrepareBlockMsg:
		packets, traffic = propPrepareBlockInPacketsMeter, propPrepareBlockInTrafficMeter
	case msg.Code == protocols.ViewChangeMsg:
		packets, traffic = propViewChangeInPacketsMeter, propViewChangeInTrafficMeter
	case msg.Code == protocols.PrepareVoteMsg:
		packets, traffic = propPrepareVoteInPacketsMeter, propPrepareVoteInTrafficMeter
	case msg.Code == protocols.PrepareBlockHashMsg:
		packets, traffic = propPrepareBlockHashInPacketsMeter, propPrepareBlockHashInTrafficMeter
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	return msg, err
}

func (rw *meteredMsgReadWriter) WriteMsg(msg p2p.Msg) error {
	// Account for the data traffic
	packets, traffic := miscOutPacketsMeter, miscOutTrafficMeter
	switch {
	case msg.Code == protocols.PrepareBlockMsg:
		packets, traffic = propPrepareBlockOutPacketsMeter, propPrepareBlockOutTrafficMeter
		common.PrepareBlockCBFTEgressTrafficMeter.Mark(int64(msg.Size))
	case msg.Code == protocols.ViewChangeMsg:
		packets, traffic = propViewChangeOutPacketsMeter, propViewChangeOutTrafficMeter
		common.ViewChangeEgressTrafficMeter.Mark(int64(msg.Size))
	case msg.Code == protocols.PrepareVoteMsg:
		packets, traffic = propPrepareVoteOutPacketsMeter, propPrepareVoteOutTrafficMeter
		common.PrepareVoteEgressTrafficMeter.Mark(int64(msg.Size))
	case msg.Code == protocols.PrepareBlockHashMsg:
		packets, traffic = propPrepareBlockHashOutPacketsMeter, propPrepareBlockHashOutTrafficMeter
		common.PrepareBlockHashEgressTrafficMeter.Mark(int64(msg.Size))

	case msg.Code == protocols.GetPrepareVoteMsg:
		common.GetPrepareVoteEgressTrafficMeter.Mark(int64(msg.Size))
	case msg.Code == protocols.GetPrepareBlockMsg:
		common.GetPrepareBlockEgressTrafficMeter.Mark(int64(msg.Size))
	case msg.Code == protocols.CBFTStatusMsg:
		common.CBFTStatusEgressTrafficMeter.Mark(int64(msg.Size))
	}
	packets.Mark(1)
	traffic.Mark(int64(msg.Size))

	// Send the packet to the p2p layer
	return rw.MsgReadWriter.WriteMsg(msg)
}
