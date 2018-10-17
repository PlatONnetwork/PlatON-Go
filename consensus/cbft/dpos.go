package cbft

import (
	"Platon-go/common"
	"Platon-go/core"
	"Platon-go/p2p/discover"
)

type dpos struct {
	primaryNodeList   []discover.Node
	chain             *core.BlockChain
	lastCycleBlockNum uint64
	startTimeOfEpoch  int64 // 一轮共识开始时间，通常是上一轮共识结束时最后一个区块的出块时间；如果是第一轮，则从1970.1.1.0.0.0.0开始

}

func newDpos(initialNodes []discover.Node) *dpos {
	dpos := &dpos{
		primaryNodeList:   initialNodes,
		lastCycleBlockNum: 0,
	}
	return dpos
}

func (d *dpos) IsPrimary(address common.Address) bool {
	// 判断当前节点是否是共识节点
	return false
}

func (d *dpos) NodeIndex(nodeID discover.NodeID) int {
	for idx, node := range d.primaryNodeList {
		if node.ID == nodeID {
			return idx
		}
	}
	return -1
}

func (d *dpos) LastCycleBlockNum() uint64 {
	// 获取最后一轮共识结束时的区块高度
	return d.lastCycleBlockNum
}

func (d *dpos) SetLastCycleBlockNum(blockNumber uint64) {
	// 设置最后一轮共识结束时的区块高度
	d.lastCycleBlockNum = blockNumber
}

// modify by platon
// 返回当前共识节点地址列表
func (b *dpos) ConsensusNodes() []discover.Node {
	return b.primaryNodeList
}

// 判断某个节点是否本轮或上一轮选举共识节点
func (b *dpos) CheckConsensusNode(id discover.NodeID) bool {
	nodes := b.ConsensusNodes()
	for _, node := range nodes {
		if node.ID == id {
			return true
		}
	}
	return false
}

// 判断当前节点是否本轮或上一轮选举共识节点
func (b *dpos) IsConsensusNode() (bool, error) {
	return false, nil
}

// 判断是否轮到当前节点打包交易出块
func (b *dpos) ShouldSeal() (bool, error) {
	return true, nil
}

func (d *dpos) StartTimeOfEpoch() int64 {

	return d.startTimeOfEpoch
}

func (d *dpos) SetStartTimeOfEpoch(startTimeOfEpoch int64) {
	// 设置最后一轮共识结束时的区块高度
	d.startTimeOfEpoch = startTimeOfEpoch
}
