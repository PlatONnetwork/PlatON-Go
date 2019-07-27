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
	PrepareBlockCBFTEgressTrafficMeter   = metrics.NewRegisteredMeter("p2p/flow/cbft/PrepareBlock/OutboundTraffic", nil)
	PrepareVoteEgressTrafficMeter        = metrics.NewRegisteredMeter("p2p/flow/cbft/PrepareVote/OutboundTraffic", nil)
	ViewChangeEgressTrafficMeter         = metrics.NewRegisteredMeter("p2p/flow/cbft/ViewChange/OutboundTraffic", nil)
	GetPrepareBlockEgressTrafficMeter    = metrics.NewRegisteredMeter("p2p/flow/cbft/GetPrepareBlock/OutboundTraffic", nil)
	PrepareBlockHashEgressTrafficMeter   = metrics.NewRegisteredMeter("p2p/flow/cbft/PrepareBlockHash/OutboundTraffic", nil)
	GetPrepareVoteEgressTrafficMeter     = metrics.NewRegisteredMeter("p2p/flow/cbft/GetPrepareVote/OutboundTraffic", nil)
	PrepareVotesEgressTrafficMeter       = metrics.NewRegisteredMeter("p2p/flow/cbft/PrepareVotes/OutboundTraffic", nil)
	GetBlockQuorumCertEgressTrafficMeter = metrics.NewRegisteredMeter("p2p/flow/cbft/GetBlockQuorumCert/OutboundTraffic", nil)
	BlockQuorumCertEgressTrafficMeter    = metrics.NewRegisteredMeter("p2p/flow/cbft/BlockQuorumCert/OutboundTraffic", nil)
	CBFTStatusEgressTrafficMeter         = metrics.NewRegisteredMeter("p2p/flow/cbft/CBFTStatus/OutboundTraffic", nil)
	GetQCBlockListEgressTrafficMeter     = metrics.NewRegisteredMeter("p2p/flow/cbft/GetQCBlockList/OutboundTraffic", nil)
	QCBlockListEgressTrafficMeter        = metrics.NewRegisteredMeter("p2p/flow/cbft/QCBlockList/OutboundTraffic", nil)
)
