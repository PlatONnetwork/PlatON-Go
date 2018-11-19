package cbft

import (
	"Platon-go/common"
	"Platon-go/core"
	"Platon-go/crypto"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"bytes"
	"sync"

	"Platon-go/core/dpos"
	"Platon-go/params"
	"Platon-go/core/state"
)

type dpos struct {
	primaryNodeList   	[]discover.NodeID
	chain             	*core.BlockChain
	lastCycleBlockNum 	uint64
	startTimeOfEpoch  	int64 // 一轮共识开始时间，通常是上一轮共识结束时最后一个区块的出块时间；如果是第一轮，则从1970.1.1.0.0.0.0开始。单位：秒
	config              *params.DposConfig

	// dpos

	// 是否正在 揭榜
	publishing 			bool

	// 读写锁
	lock 				sync.RWMutex
	// 当前轮见证人
	chairperson			common.Hash
	// 下一轮见证人
	nextChairperson  	common.Hash

	// dpos 候选人池
	candidatePool		*depos.CandidatePool
}

func newDpos(initialNodes []discover.NodeID, config *params.CbftConfig) *dpos {
	dpos := &dpos{
		primaryNodeList:   initialNodes,
		lastCycleBlockNum: 0,
		config: 			config.DposConfig,
	}
	return dpos
}

func (d *dpos) IsPrimary(addr common.Address) bool {
	// 判断当前节点是否是共识节点
	for _, node := range d.primaryNodeList {
		pub, err := node.Pubkey()
		if err != nil || pub == nil {
			log.Error("nodeID.ID.Pubkey error!")
		}
		address := crypto.PubkeyToAddress(*pub)
		return bytes.Equal(address[:], addr[:])
	}
	return false
}

func (d *dpos) NodeIndex(nodeID discover.NodeID) int64 {
	for idx, node := range d.primaryNodeList {
		if node == nodeID {
			return int64(idx)
		}
	}
	return int64(-1)
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
/*func (b *dpos) ConsensusNodes() []discover.Node {
	return b.primaryNodeList
}
*/
// 判断某个节点是否本轮或上一轮选举共识节点
/*func (b *dpos) CheckConsensusNode(id discover.NodeID) bool {
	nodes := b.ConsensusNodes()
	for _, node := range nodes {
		if node.ID == id {
			return true
		}
	}
	return false
}*/

// 判断当前节点是否本轮或上一轮选举共识节点
/*func (b *dpos) IsConsensusNode() (bool, error) {
	return true, nil
}
*/

func (d *dpos) StartTimeOfEpoch() int64 {
	return d.startTimeOfEpoch
}

func (d *dpos) SetStartTimeOfEpoch(startTimeOfEpoch int64) {
	// 设置最后一轮共识结束时的出块时间
	d.startTimeOfEpoch = startTimeOfEpoch
	log.Info("设置最后一轮共识结束时的出块时间", "startTimeOfEpoch", startTimeOfEpoch)
}
// dpos 新增func
// 设置 dpos 竞选池
func (d *dpos) SetCandidatePool(blockChain *core.BlockChain) {
//func (d *dpos) SetCandidatePool(state *state.StateDB, isgenesis bool){
	if canPool, err := depos.NewCandidatePool(blockChain, d.config); nil != err {
		log.Error("Failed to init CandidatePool", err)
	}else {
		d.candidatePool = canPool
	}
}

// 质押竞选人
func (d *dpos) SetCandidate(nodeId discover.NodeID, can *depos.Candidate){
	d.candidatePool.SetCandidate(nodeId, can)
}
// 查询入围者信息
func(d *dpos) GetCandidate(nodeId discover.NodeID) *depos.Candidate {
	return d.candidatePool.GetCandidate(nodeId)
}
// 入围者退出质押
func (d *dpos) WithdrawCandidate (nodeId discover.NodeID, price int) bool {
	return d.candidatePool.WithdrawCandidate (nodeId, price)
}
// 获取当前实时的入围者列表
func (d *dpos) GetChosens () []*depos.Candidate {
	return d.candidatePool.GetChosens()
}
// 获取当前见证人列表
func (d *dpos) GetChairpersons () []*depos.Candidate {
	return d.candidatePool.GetChairpersons()
}
// 获取某竞选者所有可提款信息
func (d *dpos) GetDefeat(nodeId discover.NodeID) []*depos.Candidate{
	return d.candidatePool.GetDefeat(nodeId)
}
// 判断某个竞选人是否入围
func (d *dpos) IsDefeat(nodeId discover.NodeID) bool {
	return d.candidatePool.IsDefeat(nodeId)
}
// 揭榜
func (d *dpos)  Election() bool {
	return d.candidatePool.Election()
}
// 提款
func (d *dpos) RefundBalance (nodeId discover.NodeID, index int) bool{
	return d.candidatePool.RefundBalance (nodeId, index)
}
// 根据nodeId查询 质押信息中的 受益者地址
func (d *dpos) GetOwner (nodeId discover.NodeID) common.Address {
	return d.candidatePool.GetOwner(nodeId)
}
// 触发替换下轮见证人列表
func (d *dpos)  Switch() bool {
	return d.candidatePool.Switch()
}
// 根据块高重置 state
func (d *dpos) ResetStateByBlockNumber  (blockNumber uint64) bool {
	return d.candidatePool.ResetStateByBlockNumber(blockNumber)
}