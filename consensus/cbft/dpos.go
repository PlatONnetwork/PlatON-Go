package cbft

import (
	"Platon-go/common"
	"time"
)

type dpos struct {
	primaryList             []common.Address //共识节点地址
	lastBlockNumOfPreEpoch  uint64           //上个轮回的最后一个块高
	lastBlockTimeOfPreEpoch int64            //上个轮回的最后一个块时间，单位：毫秒
}

var configPrimary []common.Address

func init() {
	// 读取共识节点配置文件
	// 初始化共识节点列表
	configPrimary = []common.Address{}
}

func newDpos() *dpos {
	dpos := &dpos{
		primaryList:             configPrimary,
		lastBlockNumOfPreEpoch:  0,
		lastBlockTimeOfPreEpoch: time.Date(1970, 1, 1, 0, 0, 0, 0, time.UTC).Unix(),
	}
	return dpos
}

//func (d *dpos) getCurrentPrimary() []common.Address {
func (d *dpos) getPrimaryList() []common.Address {
	// 返回共识节点列表
	return d.primaryList
}

func (d *dpos) getPrimaryIndex(address common.Address) int {
	// 判断当前节点是否是共识节点
	if len(d.primaryList) > 0 {
		for i := 0; i < len(d.primaryList); i++ {
			if d.primaryList[i] == address {
				return i
			}
		}
	}
	return -1
}

func (d *dpos) isPrimary(address common.Address) bool {
	// 判断当前节点是否是共识节点
	if len(d.primaryList) > 0 {
		for i := 0; i < len(d.primaryList); i++ {
			if d.primaryList[i] == address {
				return true
			}
		}
	}
	return false
	//return false
}
func (d *dpos) getLastBlockNumOfPreEpoch() uint64 {
	// 获取最后一轮共识结束时的区块高度
	return d.lastBlockNumOfPreEpoch
}

func (d *dpos) setLastBlockNumOfPreEpoch(blockNumber uint64) {
	// 设置最后一轮共识结束时的区块高度
	d.lastBlockNumOfPreEpoch = blockNumber
}

func (d *dpos) getLastBlockTimeOfPreEpoch() int64 {
	// 获取最后一轮共识结束时的出块时间
	return d.lastBlockTimeOfPreEpoch
}

func (d *dpos) setLastBlockTimeOfPreEpoch(lastBlockTimeOfPreEpoch int64) {
	// 设置最后一轮共识结束时的区块高度
	d.lastBlockTimeOfPreEpoch = lastBlockTimeOfPreEpoch
}

// modify by platon
// 返回当前共识节点地址列表
func (b *Cbft) ConsensusNodes() ([]string, error) {
	return nil, nil
}

func (b *Cbft) ShouldSeal() (bool, error) {
	return false, nil
}
