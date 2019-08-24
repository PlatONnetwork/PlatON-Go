package network

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// Cbft defines the network layer to use the relevant interface
// to the consensus layer.
type Cbft interface {

	// Returns the ID value of the current node.
	NodeID() discover.NodeID

	// Return a list of all consensus nodes.
	ConsensusNodes() ([]discover.NodeID, error)

	// Return configuration information of CBFT consensus.
	Config() *types.Config

	// Entrance: The messages related to the consensus are entered from here.
	// The message sent from the peer node is sent to the CBFT message queue and
	// there is a loop that will distribute the incoming message.
	ReceiveMessage(msg *types.MsgInfo) error

	// ReceiveSyncMsg is used to receive messages that are synchronized from other nodes.
	ReceiveSyncMsg(msg *types.MsgInfo) error

	// Return the highest QC block number of the current node.
	HighestQCBlockBn() (uint64, common.Hash)

	// Return the highest locked block number of the current node.
	HighestLockBlockBn() (uint64, common.Hash)

	// Return the highest commit block number of the current node.
	HighestCommitBlockBn() (uint64, common.Hash)

	// Returns the node ID of the missing vote.
	MissingViewChangeNodes() (*protocols.GetViewChange, error)

	// Returns the missing vote.
	MissingPrepareVote() (*protocols.GetPrepareVote, error)

	// OnPong records net delay time.
	OnPong(nodeID string, netLatency int64) error

	// BlockExists determines if a block exists.
	BlockExists(blockNumber uint64, blockHash common.Hash) error
}
