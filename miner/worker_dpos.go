package miner

import (
	"Platon-go/consensus"
	"Platon-go/core/state"
	"Platon-go/p2p/discover"
	"math/big"
	"errors"
	"Platon-go/log"
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

func (w *worker) attemptAddConsensusPeer(blockNumber *big.Int, state *state.StateDB, flag int) {
	if should := w.shouldElection(blockNumber); should {
		consensusNodes, err := w.getWitness(blockNumber, state, flag)
		if err == nil && len(consensusNodes) > 0 {
			w.addConsensusPeerFn(consensusNodes)
		}
	}
}

func (w *worker) getWitness(blockNumber *big.Int, state *state.StateDB, flag int) ([]*discover.Node, error) {
	log.Info("获取当前轮榜单", "blockNumber", blockNumber, "flag", flag)
	consensusNodes, err := w.engine.(consensus.Bft).GetWitness(state, flag)
	if err != nil {
		log.Error("Failed to GetWitness", "blockNumber", blockNumber, "state", state, "error", err)
		return nil, errors.New("Failed to GetWitness")
	}
	return consensusNodes, nil
}


