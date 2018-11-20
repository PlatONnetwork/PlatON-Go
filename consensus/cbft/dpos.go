package cbft

import (
	"Platon-go/common"
	"Platon-go/core"
	"Platon-go/crypto"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"bytes"
	"sync"

	"Platon-go/core/Dpos"
	"Platon-go/params"
)

type Dpos struct {
	primaryNodeList   	[]discover.NodeID
	chain             	*core.BlockChain
	lastCycleBlockNum 	uint64
	startTimeOfEpoch  	int64 // 一轮共识开始时间，通常是上一轮共识结束时最后一个区块的出块时间；如果是第一轮，则从1970.1.1.0.0.0.0开始。单位：秒
	config              *params.DposConfig

	// Dpos

	// 是否正在 揭榜
	publishing 			bool

	// 读写锁
	lock 				sync.RWMutex
	// 当前轮见证人
	chairperson			common.Hash
	// 下一轮见证人
	nextChairperson  	common.Hash

	// Dpos 候选人池
	candidatePool		*depos.CandidatePool
}
// 定义一个全局的Dpos
var DposPtr *Dpos

func newDpos(initialNodes []discover.NodeID, config *params.CbftConfig) *Dpos {
	DposPtr = &Dpos{
		primaryNodeList:   initialNodes,
		lastCycleBlockNum: 0,
		config: 			config.DposConfig,
	}
	return DposPtr
}

func (d *Dpos) IsPrimary(addr common.Address) bool {
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

func (d *Dpos) NodeIndex(nodeID discover.NodeID) int64 {
	for idx, node := range d.primaryNodeList {
		if node == nodeID {
			return int64(idx)
		}
	}
	return int64(-1)
}

func (d *Dpos) LastCycleBlockNum() uint64 {
	// 获取最后一轮共识结束时的区块高度
	return d.lastCycleBlockNum
}

func (d *Dpos) SetLastCycleBlockNum(blockNumber uint64) {
	// 设置最后一轮共识结束时的区块高度
	d.lastCycleBlockNum = blockNumber
}

// modify by platon
// 返回当前共识节点地址列表
/*func (b *Dpos) ConsensusNodes() []discover.Node {
	return b.primaryNodeList
}
*/
// 判断某个节点是否本轮或上一轮选举共识节点
/*func (b *Dpos) CheckConsensusNode(id discover.NodeID) bool {
	nodes := b.ConsensusNodes()
	for _, node := range nodes {
		if node.ID == id {
			return true
		}
	}
	return false
}*/

// 判断当前节点是否本轮或上一轮选举共识节点
/*func (b *Dpos) IsConsensusNode() (bool, error) {
	return true, nil
}
*/

func (d *Dpos) StartTimeOfEpoch() int64 {
	return d.startTimeOfEpoch
}

func (d *Dpos) SetStartTimeOfEpoch(startTimeOfEpoch int64) {
	// 设置最后一轮共识结束时的出块时间
	d.startTimeOfEpoch = startTimeOfEpoch
	log.Info("设置最后一轮共识结束时的出块时间", "startTimeOfEpoch", startTimeOfEpoch)
}
// Dpos 新增func
// 设置 Dpos 竞选池
func (d *Dpos) SetCandidatePool(blockChain *core.BlockChain) {
//func (d *Dpos) SetCandidatePool(state *state.StateDB, isgenesis bool){
	if canPool, err := depos.NewCandidatePool(blockChain, d.config); nil != err {
		log.Error("Failed to init CandidatePool", err)
	}else {
		d.candidatePool = canPool
	}
}

// 质押竞选人
func (d *Dpos) SetCandidate(nodeId discover.NodeID, can *depos.Candidate) error{
	return d.candidatePool.SetCandidate(nodeId, can)
}
// 查询入围者信息
func(d *Dpos) GetCandidate(nodeId discover.NodeID) (*depos.Candidate, error) {
	return d.candidatePool.GetCandidate(nodeId)
}
// 入围者退出质押
func (d *Dpos) WithdrawCandidate (nodeId discover.NodeID, price int) error {
	return d.candidatePool.WithdrawCandidate (nodeId, price)
}
// 获取当前实时的入围者列表
func (d *Dpos) GetChosens () []*depos.Candidate {
	return d.candidatePool.GetChosens()
}
// 获取当前见证人列表
func (d *Dpos) GetChairpersons () []*depos.Candidate {
	return d.candidatePool.GetChairpersons()
}
// 获取某竞选者所有可提款信息
func (d *Dpos) GetDefeat(nodeId discover.NodeID) ([]*depos.Candidate, error){
	return d.candidatePool.GetDefeat(nodeId)
}
// 判断某个竞选人是否入围
func (d *Dpos) IsDefeat(nodeId discover.NodeID) (bool, error) {
	return d.candidatePool.IsDefeat(nodeId)
}
// 揭榜
func (d *Dpos)  Election() bool {
	return d.candidatePool.Election()
}
// 提款
func (d *Dpos) RefundBalance (nodeId discover.NodeID) error{
	return d.candidatePool.RefundBalance (nodeId)
}
// 根据nodeId查询 质押信息中的 受益者地址
func (d *Dpos) GetOwner (nodeId discover.NodeID) common.Address {
	return d.candidatePool.GetOwner(nodeId)
}
// 触发替换下轮见证人列表
func (d *Dpos)  Switch() bool {
	return d.candidatePool.Switch()
}
// 根据块高重置 state
func (d *Dpos) ResetStateByBlockNumber  (blockNumber uint64) bool {
	return d.candidatePool.ResetStateByBlockNumber(blockNumber)
}

func GetDpos() *Dpos{
	return DposPtr
}