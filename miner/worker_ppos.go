package miner

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"errors"
	"math/big"
)

func (w *worker) shouldElection(blockNumber *big.Int) bool {
	d := new(big.Int).Sub(blockNumber, big.NewInt(cbft.BaseElection))
	_, m := new(big.Int).DivMod(d, big.NewInt(cbft.BaseSwitchWitness), new(big.Int))
	return m.Cmp(big.NewInt(0)) == 0
}

func (w *worker) shouldSwitch(blockNumber *big.Int) bool {
	_, m := new(big.Int).DivMod(blockNumber, big.NewInt(cbft.BaseSwitchWitness), new(big.Int))
	return m.Cmp(big.NewInt(0)) == 0
}

func (w *worker) shouldAddNextPeers(blockNumber *big.Int) bool {
	d := new(big.Int).Sub(blockNumber, big.NewInt(cbft.BaseAddNextPeers))
	_, m := new(big.Int).DivMod(d, big.NewInt(cbft.BaseSwitchWitness), new(big.Int))
	return m.Cmp(big.NewInt(0)) == 0
}
/*
func (w *worker) shouldRemoveFormerPeers(blockNumber *big.Int) bool {
	d := new(big.Int).Sub(blockNumber, big.NewInt(cbft.BaseRemoveFormerPeers))
	_, m := new(big.Int).DivMod(d, big.NewInt(cbft.BaseSwitchWitness), new(big.Int))
	return m.Cmp(big.NewInt(0)) == 0
}
*/
func (w *worker) election(state *state.StateDB, blockNumber *big.Int) error {
	if cbftEngine, ok := w.engine.(consensus.Bft); ok {
		if should := w.shouldElection(blockNumber); should {
			log.Debug("Election call:", "blockNumber", blockNumber)
			_, err := cbftEngine.Election(state, blockNumber)
			if err != nil {
				log.Error("Failed to election", "blockNumber", blockNumber, "error", err)
				return errors.New("Failed to Election")
			}
			log.Debug("Success to election", "blockNumber", blockNumber)
		}
	}
	return nil
}

func (w *worker) switchWitness(state *state.StateDB, blockNumber *big.Int) error {
	if cbftEngine, ok := w.engine.(consensus.Bft); ok {
		if should := w.shouldSwitch(blockNumber); should {
			log.Debug("SwitchWitness call:", "blockNumber", blockNumber)
			success := cbftEngine.Switch(state)
			if !success {
				log.Error("Failed to switchWitness", "blockNumber", blockNumber)
				return errors.New("Failed to switchWitness")
			}
			log.Debug("Success to switchWitness", "blockNumber", blockNumber)
		}
	}
	return nil
}

func (w *worker) attemptAddConsensusPeer(blockNumber *big.Int, state *state.StateDB) {
	if should := w.shouldAddNextPeers(blockNumber); should {
		log.Debug("Attempt to connect the next round of consensus nodes", "blockNumber", blockNumber)
		// At this point, only the next round of consensus nodes can be obtained from statedb,
		// because switch has not been invoked for switching, and next in memory has not been refreshed.
		nextNodes, err := w.getWitness(blockNumber, state, 1)	// flag：-1: former	  0: current   1: next
		log.Info("Next round consensus node list:","number", blockNumber, "nextNodes", nextNodes, "nextNodes length", len(nextNodes), "err", err)
		if err == nil && len(nextNodes) > 0 && existsNode(w.engine.(consensus.Bft).GetOwnNodeID(), nextNodes) {
			w.addConsensusPeerFn(nextNodes)
		}
	}
}
/*
func (w *worker) attemptRemoveConsensusPeer(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, state *state.StateDB) {
	if should := w.shouldRemoveFormerPeers(blockNumber); should {
		log.Info("Attempt to disconnect the former round of consensus nodes","parentNumber", parentNumber, "parentHash", parentHash, "blockNumber", blockNumber)
		formerNodes := w.engine.(consensus.Bft).FormerNodes(parentNumber, parentHash, blockNumber)
		log.Info("Former round consensus node list:","parentNumber", parentNumber, "parentHash", parentHash, "number", blockNumber, "formerNodes", formerNodes, "formerNodes length", len(formerNodes))
		currentNodes := w.engine.(consensus.Bft).CurrentNodes(parentNumber, parentHash, blockNumber)
		log.Info("Current round consensus node list:","parentNumber", parentNumber, "parentHash", parentHash, "number", blockNumber, "currentNodes", currentNodes, "currentNodes length", len(currentNodes))

		removeNodes := formerNodes
		ownNodeID := w.engine.(consensus.Bft).GetOwnNodeID()
		if len(formerNodes) > 0 && len(currentNodes) > 0 && existsNode(ownNodeID, formerNodes) {
			if existsNode(ownNodeID, currentNodes) {
				currentNodesMap := make(map[discover.NodeID]discover.NodeID)
				for _,n := range currentNodes {
					currentNodesMap[n.ID] = n.ID
				}
				for _,n := range formerNodes {
					if _,ok := currentNodesMap[n.ID]; !ok {
						removeNodes = append(removeNodes, n)
					}
				}
			}

			if len(removeNodes) > 0 {
				w.removeConsensusPeerFn(removeNodes)
			}
		}
	}
}
*/
func existsNode(nodeID discover.NodeID, nodes []*discover.Node) bool {
	for _,n := range nodes {
		if nodeID == n.ID {
			return true
		}
	}
	return false
}

func (w *worker) getWitness(blockNumber *big.Int, state *state.StateDB, flag int) ([]*discover.Node, error) {
	log.Debug("GetWitness begin", "blockNumber", blockNumber, "flag", flag)
	consensusNodes, err := w.engine.(consensus.Bft).GetWitness(state, flag)
	if err != nil {
		log.Error("Failed to GetWitness", "blockNumber", blockNumber, "state", state, "error", err)
		return nil, errors.New("Failed to GetWitness")
	}
	log.Info("GetWitness end", "blockNumber", blockNumber, "flag", flag, "consensusNodes", consensusNodes, "consensusNodes length", len(consensusNodes))
	return consensusNodes, nil
}

func (w *worker) setNodeCache(state *state.StateDB, parentNumber, currentNumber *big.Int, parentHash, currentHash common.Hash) error {
	return w.engine.(consensus.Bft).SetNodeCache(state, parentNumber, currentNumber, parentHash, currentHash)
}