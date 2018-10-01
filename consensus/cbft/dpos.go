package cbft

import (
	"Platon-go/common"
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




