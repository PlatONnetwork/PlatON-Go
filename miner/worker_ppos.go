package miner

import (
	"Platon-go/common"
	"Platon-go/consensus"
	"Platon-go/consensus/cbft"
	"Platon-go/core/state"
	"Platon-go/log"
	"Platon-go/p2p/discover"
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
			log.Info("请求揭榜", "blockNumber", blockNumber)
			_, err := cbftEngine.Election(state, blockNumber)
			if err != nil {
				log.Error("Failed to election", "blockNumber", blockNumber, "error", err)
				return errors.New("Failed to Election")
			}
			log.Info("Success to election", "blockNumber", blockNumber)
		}
	}
	return nil
}

func (w *worker) switchWitness(state *state.StateDB, blockNumber *big.Int) error {
	if cbftEngine, ok := w.engine.(consensus.Bft); ok {
		if should := w.shouldSwitch(blockNumber); should {
			log.Info("触发替换下轮见证人列表", "blockNumber", blockNumber)
			success := cbftEngine.Switch(state)
			if !success {
				log.Error("Failed to switchWitness", "blockNumber", blockNumber)
				return errors.New("Failed to switchWitness")
			}
			log.Info("Success to switchWitness", "blockNumber", blockNumber)
		}
	}
	return nil
}

func (w *worker) attemptAddConsensusPeer(blockNumber *big.Int, state *state.StateDB) {
	if should := w.shouldAddNextPeers(blockNumber); should {
		log.Info("尝试连接下一轮共识节点", "blockNumber", blockNumber)
		// 此时只能从statedb中获取下一轮共识节点，因为还没有调用switch进行切换，内存中next没有刷新
		nextNodes, err := w.getWitness(blockNumber, state, 1)	// flag：-1: 上一轮	  0: 本轮见证人   1: 下一轮见证人
		log.Info("下一轮共识节点列表","number", blockNumber, "nextNodes", nextNodes, "nextNodes length", len(nextNodes), "err", err)
		if err == nil && len(nextNodes) > 0 && existsNode(w.engine.(consensus.Bft).GetOwnNodeID(), nextNodes) {
			w.addConsensusPeerFn(nextNodes)
		}
	}
}
/*
func (w *worker) attemptRemoveConsensusPeer(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, state *state.StateDB) {
	if should := w.shouldRemoveFormerPeers(blockNumber); should {
		log.Info("尝试断开上一轮共识节点","parentNumber", parentNumber, "parentHash", parentHash, "blockNumber", blockNumber)
		formerNodes := w.engine.(consensus.Bft).FormerNodes(parentNumber, parentHash, blockNumber)
		log.Info("上一轮共识节点列表","parentNumber", parentNumber, "parentHash", parentHash, "number", blockNumber, "formerNodes", formerNodes, "formerNodes length", len(formerNodes))
		currentNodes := w.engine.(consensus.Bft).CurrentNodes(parentNumber, parentHash, blockNumber)
		log.Info("当前轮共识节点列表","parentNumber", parentNumber, "parentHash", parentHash, "number", blockNumber, "currentNodes", currentNodes, "currentNodes length", len(currentNodes))

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
	log.Info("getWitness begin", "blockNumber", blockNumber, "flag", flag)
	consensusNodes, err := w.engine.(consensus.Bft).GetWitness(state, flag)
	if err != nil {
		log.Error("Failed to GetWitness", "blockNumber", blockNumber, "state", state, "error", err)
		return nil, errors.New("Failed to GetWitness")
	}
	log.Info("getWitness end", "blockNumber", blockNumber, "flag", flag, "consensusNodes", consensusNodes, "consensusNodes length", len(consensusNodes))
	return consensusNodes, nil
}

func (w *worker) setNodeCache(state *state.StateDB, parentNumber, currentNumber *big.Int, parentHash, currentHash common.Hash) error {
	return w.engine.(consensus.Bft).SetNodeCache(state, parentNumber, currentNumber, parentHash, currentHash)
}