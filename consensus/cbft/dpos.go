package cbft

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	errInvalidatorCandidateAddress = errors.New("invalid address")
)

type dpos struct {
	primaryNodeList   []discover.NodeID
	chain             *core.BlockChain
	lastCycleBlockNum uint64
	startTimeOfEpoch  int64 // A round of consensus start time is usually the block time of the last block at the end of the last round of consensus;
	// if it is the first round, it starts from 1970.1.1.0.0.0.0. Unit: second

}

func newDpos(initialNodes []discover.NodeID) *dpos {
	dpos := &dpos{
		primaryNodeList:   initialNodes,
		lastCycleBlockNum: 0,
	}
	return dpos
}

func (d *dpos) IsPrimary(addr common.Address) bool {
	// Determine whether the current node is a consensus node
	for _, node := range d.primaryNodeList {
		pub, err := node.Pubkey()
		if err != nil || pub == nil {
			log.Error("nodeID.ID.Pubkey error!")
		}
		address := crypto.PubkeyToAddress(*pub)
		return bytes.Equal(address[:], addr[:])
	}
	return false
}

func (d *dpos) NodeID(index int) discover.NodeID {
	return d.primaryNodeList[index]
}
func (d *dpos) AddressIndex(addr common.Address) (int, error) {
	// Determine whether the current node is a consensus node
	for i, node := range d.primaryNodeList {
		pub, err := node.Pubkey()
		if err != nil || pub == nil {
			log.Error(fmt.Sprintf("NodeID Pubkey error!"))
		}
		address := crypto.PubkeyToAddress(*pub)

		if bytes.Equal(address[:], addr[:]) {
			return i, nil
		}
	}
	return -1, errInvalidatorCandidateAddress
}

func (d *dpos) NodeIndex(nodeID discover.NodeID) (int, error) {
	for idx, node := range d.primaryNodeList {
		if node == nodeID {
			return idx, nil
		}
	}
	return -1, errInvalidatorCandidateAddress
}

func (d *dpos) NodeIndexAddress(nodeID discover.NodeID) (int, common.Address, error) {
	for idx, node := range d.primaryNodeList {
		if node == nodeID {
			pubkey, err := nodeID.Pubkey()
			if err != nil {
				break
			}
			return idx, crypto.PubkeyToAddress(*pubkey), nil
		}
	}
	return -1, common.Address{}, errInvalidatorCandidateAddress
}

func (d *dpos) LastCycleBlockNum() uint64 {
	// Get the block height at the end of the final round of consensus
	return d.lastCycleBlockNum
}

func (d *dpos) SetLastCycleBlockNum(blockNumber uint64) {
	// Set the block height at the end of the last round of consensus
	d.lastCycleBlockNum = blockNumber
}

func (d *dpos) Total() int {
	return len(d.primaryNodeList)
}

// Returns the current consensus node address list
/*func (b *dpos) ConsensusNodes() []discover.Node {
	return b.primaryNodeList
}
*/
// Determine whether a node is the current or previous round of election consensus nodes
/*func (b *dpos) CheckConsensusNode(id discover.NodeID) bool {
	nodes := b.ConsensusNodes()
	for _, node := range nodes {
		if node.ID == id {
			return true
		}
	}
	return false
}*/

// Determine whether the current node is the current or previous round of election consensus nodes
/*func (b *dpos) IsConsensusNode() (bool, error) {
	return true, nil
}
*/

func (d *dpos) StartTimeOfEpoch() int64 {
	return d.startTimeOfEpoch
}

func (d *dpos) SetStartTimeOfEpoch(startTimeOfEpoch int64) {
	// Set the block time at the end of the last round of consensus
	d.startTimeOfEpoch = startTimeOfEpoch
	log.Info("Set the block time at the end of the last round of consensus", "startTimeOfEpoch", startTimeOfEpoch)
}
