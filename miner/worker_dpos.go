package miner

import (
	"Platon-go/consensus"
	"Platon-go/core/state"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"errors"
	"math/big"
)

func (w *worker) election(blockNumber *big.Int) error {
	if cbftEngine, ok := w.engine.(consensus.Bft); ok {
		if should := w.shouldElection(blockNumber); should {
			log.Info("请求揭榜", "blockNumber", blockNumber)
			_, err := cbftEngine.Election(w.current.state)
			if err != nil {
				log.Error("Failed to Election", "blockNumber", blockNumber, "error", err)
				return errors.New("Failed to Election")
			}
		}
	}
	return nil
}

func (w *worker) shouldElection(blockNumber *big.Int) bool {
	_, m := new(big.Int).DivMod(blockNumber, big.NewInt(baseElection), new(big.Int))
	return m.Cmp(big.NewInt(0)) == 0
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
		}
	}
	return nil
}

func (w *worker) shouldSwitch(blockNumber *big.Int) bool {
	_, m := new(big.Int).DivMod(blockNumber, big.NewInt(baseSwitchWitness), new(big.Int))
	return m.Cmp(big.NewInt(0)) == 0
}

func (w *worker) attemptAddConsensusPeer(blockNumber *big.Int, state *state.StateDB) {
	if should := w.shouldElection(blockNumber); should {
		consensusNodes, err := w.getWitness(blockNumber, state, 1)	// flag：-1: 上一轮	  0: 本轮见证人   1: 下一轮见证人
		if err == nil && len(consensusNodes) > 0 {
			w.addConsensusPeerFn(consensusNodes)
		}
	}
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

func (w *worker) getAllWitness(blockNumber *big.Int, state *state.StateDB) ([]*discover.Node, []*discover.Node, []*discover.Node, error) {
	log.Info("getAllWitness", "blockNumber", blockNumber)
	preArr, curArr, nextArr, err := w.engine.(consensus.Bft).GetAllWitness(state)
	if err != nil {
		log.Error("Failed to getAllWitness", "blockNumber", blockNumber, "state", state, "error", err)
		return nil, nil, nil, errors.New("Failed to GetWitness")
	}
	return preArr, curArr, nextArr, nil
}

func (w *worker) attemptRemoveConsensusPeer(blockNumber *big.Int, state *state.StateDB) {
	if should := w.shouldRemoveFormer(blockNumber); should {
		formerNodes,currentNodes,_,err := w.getAllWitness(blockNumber, state)	// 上一轮、当前轮
		removeNodes := make([]*discover.Node, 0, len(formerNodes))
		if err == nil && len(formerNodes) > 0 && len(currentNodes) > 0 {
			currentNodesMap := make(map[discover.NodeID]discover.NodeID)
			for _,n := range currentNodes {
				currentNodesMap[n.ID] = n.ID
			}
			for _,n := range formerNodes {
				if _,ok := currentNodesMap[n.ID]; !ok {
					removeNodes = append(removeNodes, n)
				}
			}
			if len(removeNodes) > 0 {
				w.removeConsensusPeerFn(removeNodes)
			}
		}
	}
}

func (w *worker) shouldRemoveFormer(blockNumber *big.Int) bool {
	_, m := new(big.Int).DivMod(blockNumber, big.NewInt(baseRemoveFormerPeers), new(big.Int))
	return m.Cmp(big.NewInt(0)) == 0
}