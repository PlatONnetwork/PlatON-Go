package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"math/big"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

type roundCache map[uint64]map[common.Hash]*nodeCache

type nodeCache struct {
	former            *pposRound 	// the previous round of witnesses nodeId
	current           *pposRound 	// the current round of witnesses nodeId
	next              *pposRound 	// the next round of witnesses nodeId
}

type pposRound struct {
	nodeIds []discover.NodeID
	nodes 	[]*discover.Node
	start *big.Int
	end   *big.Int
}


func (r roundCache) getFormerRound(blockNumber *big.Int, blockHash common.Hash) *pposRound {
	num := blockNumber.Uint64()
	if round, ok := r[num]; ok {
		if node, has := round[blockHash]; has {
			if nil != node.former {
				return node.former
			}
		}
	}
	return nil
}


func (r roundCache) getCurrentRound (blockNumber *big.Int, blockHash common.Hash) *pposRound {
	num := blockNumber.Uint64()
	if round, ok := r[num]; ok {
		if node, has := round[blockHash]; has {
			if nil != node.current {
				return node.current
			}
		}
	}
	return nil
}

func (r roundCache) getNextRound (blockNumber *big.Int, blockHash common.Hash) *pposRound {
	num := blockNumber.Uint64()
	if round, ok := r[num]; ok {
		if node, has := round[blockHash]; has {
			if nil != node.next {
				return node.next
			}
		}
	}
	return nil
}

func (r roundCache) getNodeCache (blockNumber *big.Int, blockHash common.Hash) *nodeCache {
	num := blockNumber.Uint64()
	if round, ok := r[num]; ok {
		if node, has := round[blockHash]; has {
			return node
		}
	}
	return nil
}


func (r roundCache) setNodeCache (blockNumber *big.Int, blockHash common.Hash, cache *nodeCache) {
	num := blockNumber.Uint64()
	var node map[common.Hash]*nodeCache
	if _, ok := r[num]; ok {
		node = r[num]
	}else {
		node = make(map[common.Hash]*nodeCache)
	}
	if _, ok := node[blockHash]; !ok {
		node[blockHash] = cache
		r[num] = node
	}
	// delete index out of BaseIrrCount
	baseIrrCountNum := num - BaseIrrCount + 1
	if baseIrrCountNum > 0 {
		delete(r, baseIrrCountNum)
	}
}

