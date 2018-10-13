package cbft

import (
	"Platon-go/common"
	"Platon-go/p2p/discover"
)

type dpos struct {
	primaryList []common.Address
	lastCycleBlockNum uint64
}

var configPrimary []common.Address

func init() {
	// 读取共识节点配置文件
	// 初始化共识节点列表
	configPrimary = []common.Address{}
}

func newDpos() *dpos {
	dpos := &dpos {
		primaryList : configPrimary,
		lastCycleBlockNum : 0,
	}
	return dpos
}

func (d *dpos) getCurrentPrimary() []common.Address {
	// 返回共识节点列表
	return nil
}

func (d *dpos) isPrimary(address common.Address) bool {
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
	return nil,nil
}

// 判断某个节点是否本轮或上一轮选举共识节点
func (b *Cbft) CheckConsensusNode(discover.NodeID) (bool, error) {
	return false,nil
}

// 判断当前节点是否本轮或上一轮选举共识节点
func (b *Cbft) IsConsensusNode() (bool, error) {
	return false,nil
}

// 判断是否轮到当前节点打包交易出块
func (b *Cbft) ShouldSeal() (bool, error) {
	return false,nil
}



