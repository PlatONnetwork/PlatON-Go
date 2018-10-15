package cbft

import (
	"Platon-go/common"
	"Platon-go/core"
	"Platon-go/p2p/discover"
)

type dpos struct {
	primaryNodeList   []discover.Node
	chain *core.BlockChain
	lastCycleBlockNum uint64
}

func newDpos(initialNodes []discover.Node) *dpos {
	dpos := &dpos {
		primaryNodeList: initialNodes,
		lastCycleBlockNum : 0,
	}
	return dpos
}

func (d *dpos) getCurrentPrimary() []discover.Node {
	// 返回共识节点列表
	return d.primaryNodeList
}

func (d *dpos) IsPrimary(address common.Address) bool {
	// 判断当前节点是否是共识节点
	return false
}

func (d *dpos) getLastCycleBlockNum() uint64 {
	// 获取最后一轮共识结束时的区块高度
	return d.lastCycleBlockNum
}

func (d *dpos) setLastCycleBlockNum(blockNumber uint64) {
	// 设置最后一轮共识结束时的区块高度
	d.lastCycleBlockNum = blockNumber
}

// modify by platon
// 返回当前共识节点地址列表
func (b *Cbft) ConsensusNodes() ([]discover.Node, error) {
	return b.dpos.getCurrentPrimary(),nil
}
// 判断某个节点是否本轮或上一轮选举共识节点
func (b *Cbft) CheckConsensusNode(id discover.NodeID) (bool, error) {
	nodes,err := b.ConsensusNodes()
	if err != nil {
		return false, err
	}

	for _, node := range nodes {
		if node.ID == id {
			return true, nil
		}
	}
	return false,nil
}

// 判断当前节点是否本轮或上一轮选举共识节点
func (b *Cbft) IsConsensusNode() (bool, error) {
	return false,nil
}

// 判断是否轮到当前节点打包交易出块
func (b *Cbft) ShouldSeal() (bool, error) {
	return true,nil
}


