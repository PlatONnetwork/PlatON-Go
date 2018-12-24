package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

type rotating struct {
	ppos         *ppos
	rotaList     []common.Address // 本轮循环出块节点顺序列表
	startTime    int64            // 本轮循环开始时间戳，单位毫秒
	endTime      int64            // 本轮循环结束时间戳，单位毫秒
	timeInterval int64            // 每个节点出块时间，单位毫秒
}

func newRotating(ppos *ppos, timeInterval int64) *rotating {
	rotating := &rotating{
		ppos:         ppos,
		timeInterval: timeInterval,
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

/*func (r *rotating) inturn(number uint64, signer common.Address) bool {
	sort.Sort(signerOrderingRule(r.rotaList))
	offset :=  0
	for offset < len(r.rotaList) && r.rotaList[offset] != signer {
		offset++
	}
	return (number % uint64(len(r.rotaList))) == uint64(offset)
}

type signerOrderingRule []common.Address
func (s signerOrderingRule) Len() int           { return len(s) }
func (s signerOrderingRule) Less(i, j int) bool { return bytes.Compare(s[i][:], s[j][:]) < 0 }
func (s signerOrderingRule) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }*/
