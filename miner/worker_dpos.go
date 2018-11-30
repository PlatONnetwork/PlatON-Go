package miner

import (
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

func (w *worker) shouldAddNext(blockNumber *big.Int) bool {
	d := new(big.Int).Sub(blockNumber, big.NewInt(cbft.BaseAddNextPeers))
	_, m := new(big.Int).DivMod(d, big.NewInt(cbft.BaseSwitchWitness), new(big.Int))
	return m.Cmp(big.NewInt(0)) == 0
}

func (w *worker) shouldRemoveFormer(blockNumber *big.Int) bool {
	d := new(big.Int).Sub(blockNumber, big.NewInt(cbft.BaseRemoveFormerPeers))
	_, m := new(big.Int).DivMod(d, big.NewInt(cbft.BaseSwitchWitness), new(big.Int))
	return m.Cmp(big.NewInt(0)) == 0
}

func (w *worker) election(blockNumber *big.Int) error {
	if cbftEngine, ok := w.engine.(consensus.Bft); ok {
		if should := w.shouldElection(blockNumber); should {
			log.Info("请求揭榜", "blockNumber", blockNumber)
			_, err := cbftEngine.Election(w.current.state, blockNumber)
			if err != nil {
				log.Error("Failed to election", "blockNumber", blockNumber, "error", err)
				return errors.New("Failed to Election")
			}
			log.Info("Success to election", "blockNumber", blockNumber)
		}
	}
	return nil
}

func (w *worker) switchWitness(blockNumber *big.Int) error {
	if cbftEngine, ok := w.engine.(consensus.Bft); ok {
		if should := w.shouldSwitch(blockNumber); should {
			log.Info("触发替换下轮见证人列表", "blockNumber", blockNumber)
			success := cbftEngine.Switch(w.current.state)
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
	if should := w.shouldElection(blockNumber); should {
		log.Info("尝试连接下一轮共识节点", "blockNumber", blockNumber)
		// 此时只能从statedb中获取下一轮共识节点，因为还没有调用switch进行切换，内存中next没有刷新
		nextNodes, err := w.getWitness(blockNumber, state, 1)	// flag：-1: 上一轮	  0: 本轮见证人   1: 下一轮见证人
		log.Info("下一轮共识节点列表","number", blockNumber, "nextNodes", nextNodes, "nextNodes length", len(nextNodes), "err", err)
		if err == nil && len(nextNodes) > 0 && existsNode(w.engine.(consensus.Bft).GetOwnNodeID(), nextNodes) {
			w.addConsensusPeerFn(nextNodes)
		}
	}
}

func (w *worker) attemptRemoveConsensusPeer(blockNumber *big.Int, state *state.StateDB) {
	if should := w.shouldRemoveFormer(blockNumber); should {
		log.Info("尝试断开上一轮共识节点", "blockNumber", blockNumber)
		formerNodes := w.engine.(consensus.Bft).FormerNodes()
		log.Info("上一轮共识节点列表","number", blockNumber, "formerNodes", formerNodes, "formerNodes length", len(formerNodes))
		currentNodes := w.engine.(consensus.Bft).CurrentNodes()
		log.Info("当前轮共识节点列表","number", blockNumber, "currentNodes", currentNodes, "currentNodes length", len(currentNodes))

		removeNodes := make([]*discover.Node, 0, len(formerNodes))
		if len(formerNodes) > 0 && len(currentNodes) > 0 && existsNode(w.engine.(consensus.Bft).GetOwnNodeID(), formerNodes) {
			if !existsNode(w.engine.(consensus.Bft).GetOwnNodeID(), currentNodes) {
				removeNodes = formerNodes
			} else {
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

func existsNode(nodeID discover.NodeID, nodeIDList []discover.NodeID) bool {
	for _,v := range nodeIDList {
		if nodeID == v {
			return true
		}
	}
	return false
}

func (w *worker) getWitness(blockNumber *big.Int, state *state.StateDB, flag int) ([]*discover.Node, error) {
	log.Info("getWitness", "blockNumber", blockNumber, "flag", flag)
	consensusNodes, err := w.engine.(consensus.Bft).GetWitness(state, flag)
	if err != nil {
		log.Error("Failed to GetWitness", "blockNumber", blockNumber, "state", state, "error", err)
		return nil, errors.New("Failed to GetWitness")
	}
	return consensusNodes, nil
}