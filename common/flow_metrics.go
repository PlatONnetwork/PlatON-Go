package common

import (
	"github.com/PlatONnetwork/PlatON-Go/metrics"
)

var (
	//  the meter for the protocol of eth
	StatusEgressTrafficMeter          = metrics.NewRegisteredMeter("p2p/flow/eth/status/OutboundTraffic", nil)
	NewBlockHashesEgressTrafficMeter  = metrics.NewRegisteredMeter("p2p/flow/eth/newBlockHashes/OutboundTraffic", nil)
	TxTrafficMeter                    = metrics.NewRegisteredMeter("p2p/flow/eth/tx/OutboundTraffic", nil)
	GetBlockHeadersEgressTrafficMeter = metrics.NewRegisteredMeter("p2p/flow/eth/getBlockHeaders/OutboundTraffic", nil)
	BlockHeadersEgressTrafficMeter    = metrics.NewRegisteredMeter("p2p/flow/eth/blockHeaders/OutboundTraffic", nil)
	GetBlockBodiesEgressTrafficMeter  = metrics.NewRegisteredMeter("p2p/flow/eth/getBlockBodies/OutboundTraffic", nil)
	BlockBodiesEgressTrafficMeter     = metrics.NewRegisteredMeter("p2p/flow/eth/blockBodies/OutboundTraffic", nil)
	NewBlockEgressTrafficMeter        = metrics.NewRegisteredMeter("p2p/flow/eth/newBlock/OutboundTraffic", nil)
	PrepareBlockEgressTrafficMeter    = metrics.NewRegisteredMeter("p2p/flow/eth/prepareBlock/OutboundTraffic", nil)
	BlockSignatureEgressTrafficMeter  = metrics.NewRegisteredMeter("p2p/flow/eth/blockSignature/OutboundTraffic", nil)
	PongEgressTrafficMeter            = metrics.NewRegisteredMeter("p2p/flow/eth/pong/OutboundTraffic", nil)
	GetNodeDataEgressTrafficMeter     = metrics.NewRegisteredMeter("p2p/flow/eth/getNodeData/OutboundTraffic", nil)
	NodeDataEgressTrafficMeter        = metrics.NewRegisteredMeter("p2p/flow/eth/nodeData/OutboundTraffic", nil)
	GetReceiptsEgressTrafficMeter     = metrics.NewRegisteredMeter("p2p/flow/eth/getReceipts/OutboundTraffic", nil)
	ReceiptsTrafficMeter              = metrics.NewRegisteredMeter("p2p/flow/eth/receipts/OutboundTraffic", nil)

	// the meter for the protocol of cbft
	PrepareBlockCBFTEgressTrafficMeter       = metrics.NewRegisteredMeter("p2p/flow/cbft/prepareBlock/OutboundTraffic", nil)
	PrepareVoteEgressTrafficMeter            = metrics.NewRegisteredMeter("p2p/flow/cbft/prepareVote/OutboundTraffic", nil)
	ViewChangeEgressTrafficMeter             = metrics.NewRegisteredMeter("p2p/flow/cbft/viewChange/OutboundTraffic", nil)
	ViewChangeVoteEgressTrafficMeter         = metrics.NewRegisteredMeter("p2p/flow/cbft/viewChangeVote/OutboundTraffic", nil)
	ConfirmedPrepareBlockEgressTrafficMeter  = metrics.NewRegisteredMeter("p2p/flow/cbft/confirmedPrepareBlock/OutboundTraffic", nil)
	GetPrepareVoteEgressTrafficMeter         = metrics.NewRegisteredMeter("p2p/flow/cbft/getPrepareVote/OutboundTraffic", nil)
	PrepareVotesEgressTrafficMeter           = metrics.NewRegisteredMeter("p2p/flow/cbft/prepareVotes/OutboundTraffic", nil)
	GetPrepareBlockEgressTrafficMeter        = metrics.NewRegisteredMeter("p2p/flow/cbft/getPrepareBlock/OutboundTraffic", nil)
	GetHighestPrepareBlockEgressTrafficMeter = metrics.NewRegisteredMeter("p2p/flow/cbft/getHighestPrepareBlock/OutboundTraffic", nil)
	HighestPrepareBlockEgressTrafficMeter    = metrics.NewRegisteredMeter("p2p/flow/cbft/highestPrepareBlock/OutboundTraffic", nil)
	CBFTStatusEgressTrafficMeter             = metrics.NewRegisteredMeter("p2p/flow/cbft/CBFTStatus/OutboundTraffic", nil)
	PrepareBlockHashEgressTrafficMeter       = metrics.NewRegisteredMeter("p2p/flow/cbft/prepareBlockHash/OutboundTraffic", nil)
	GetLatestStatusEgressTrafficMeter        = metrics.NewRegisteredMeter("p2p/flow/cbft/getLatestStatus/OutboundTraffic", nil)
	LatestStatusEgressTrafficMeter           = metrics.NewRegisteredMeter("p2p/flow/cbft/latestStatus/OutboundTraffic", nil)
)
