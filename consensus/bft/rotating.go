package bft

import (
	"Platon-go/common"
)

type rotating struct {
	dpos *dpos
	rotaList []common.Address	// 本轮循环出块节点顺序列表
	startTime uint64			// 本轮循环开始时间戳，单位毫秒
	endTime	uint64				// 本轮循环结束时间戳，单位毫秒
	timeInterval uint64			// 每个节点出块时间，单位毫秒
}

func newRotating(dpos *dpos, timeInterval uint64) *rotating {
	rotating := &rotating {
		dpos : dpos,
		timeInterval : timeInterval,
	}
	return rotating
}

func sort() {
	// 新一轮共识的排序函数
	// xor(上一轮的最后区块hash + 节点公钥地址)
}

func (r *rotating) IsRotating(common.Address) bool {
	// 判断当前节点是否轮值出块
	// 根据共识排序以及时间窗口
	return false
}
