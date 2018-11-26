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
	"Platon-go/core/types"
	"Platon-go/core/vm"
	"math/big"
)

type dpos struct {
	formerlyNodeList 	[]discover.NodeID   // the previous round of witnesses nodeId
	primaryNodeList   	[]discover.NodeID	// the current round of witnesses nodeId
	nextNodeList 		[]discover.NodeID  	// the next round of witnesses nodeId
	chain             	*core.BlockChain
	lastCycleBlockNum 	uint64
	startTimeOfEpoch  	int64 // 一轮共识开始时间，通常是上一轮共识结束时最后一个区块的出块时间；如果是第一轮，则从1970.1.1.0.0.0.0开始。单位：秒
	config              *params.DposConfig

	// added by candidatepool module

	lock 				sync.RWMutex
	// the candidate pool object pointer
	candidatePool		*depos.CandidatePool
}

func newDpos(initialNodes []discover.NodeID, config *params.CbftConfig) *dpos {
	dposPtr := &dpos{
		formerlyNodeList:  initialNodes,
		primaryNodeList:   initialNodes,
		lastCycleBlockNum: 0,
		config: 			config.DposConfig,
		candidatePool:		depos.NewCandidatePool(config.DposConfig),
	}
	return dposPtr
}

func (d *dpos) IsPrimary(addr common.Address) bool {
	d.lock.RLock()
	defer d.lock.RUnlock()
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
	d.lock.RLock()
	defer d.lock.RUnlock()
	nodeList := append(d.primaryNodeList, d.formerlyNodeList...)
	for idx, node := range nodeList {
		if node == nodeID {
			return int64(idx)
		}
	}
	return int64(-1)
}

func (d *dpos) getPrimaryNodes() []discover.NodeID {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.primaryNodeList
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

/** dpos was added func */
/** Method provided to the cbft module call */
// Announce witness
func (d *dpos)  Election(state *state.StateDB) ([]*discover.Node, error) {
	return d.candidatePool.Election(state)
}

// switch next witnesses to current witnesses
func (d *dpos)  Switch(state *state.StateDB) bool {

	if !d.candidatePool.Switch(state) {
		return false
	}
	preArr, curArr, nextArr, err := d.candidatePool.GetAllWitness(state)
	if nil != err {
		return false
	}
	d.lock.Lock()
	defer d.lock.Unlock()
	if len(preArr) != 0 {
		d.formerlyNodeList = convertNodeID(preArr)
	}
	if len(curArr) != 0 {
		d.primaryNodeList = convertNodeID(curArr)
	}
	if len(nextArr) != 0 {
		d.nextNodeList = convertNodeID(nextArr)
	}
	return true
}

// Getting nodes of witnesses
// flag：-1: the previous round of witnesses  0: the current round of witnesses   1: the next round of witnesses
func (d *dpos) GetWitness(state *state.StateDB, flag int) ([]*discover.Node, error) {
	return d.candidatePool.GetWitness(state, flag)
}

func (d *dpos) GetAllWitness(state *state.StateDB) ([]*discover.Node, []*discover.Node, []*discover.Node, error) {
	return d.candidatePool.GetAllWitness(state)
}

// setting candidate pool of dpos module
func (d *dpos) SetCandidatePool(blockChain *core.BlockChain) {
	// When the highest block in the chain is not a genesis block, Need to load witness nodeIdList from the stateDB.
	if  blockChain.Genesis().NumberU64() != blockChain.CurrentBlock().NumberU64() {
		state, err := blockChain.State()
		if nil != err {
			log.Error("Load state from chain failed on SetCandidatePool err", err)
			return
		}
		if preArr, curArr, nextArr, err := d.candidatePool.GetAllWitness(state); nil != err {
			log.Error("Load Witness from state failed on SetCandidatePool err", err)
		}else {
			d.lock.Lock()
			defer d.lock.Unlock()
			if len(preArr) != 0 {
				d.formerlyNodeList = convertNodeID(preArr)
			}
			if len(curArr) != 0 {
				d.primaryNodeList = convertNodeID(curArr)
			}
			if len(nextArr) != 0 {
				d.nextNodeList = convertNodeID(nextArr)
			}
		}
	}
}



/** Method provided to the built-in contract call */
// pledge Candidate
func (d *dpos) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error{
	return d.candidatePool.SetCandidate(state, nodeId, can)
}
// Getting immediate candidate info by nodeId
func(d *dpos) GetCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error) {
	return d.candidatePool.GetCandidate(state, nodeId)
}
// candidate withdraw from immediates elected candidates
func (d *dpos) WithdrawCandidate (state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error {
	return d.candidatePool.WithdrawCandidate (state, nodeId, price, blockNumber)
}
// Getting all immediate elected candidates array
func (d *dpos) GetChosens (state vm.StateDB) []*types.Candidate {
	return d.candidatePool.GetChosens(state)
}
// Getting all witness array
func (d *dpos) GetChairpersons (state vm.StateDB) []*types.Candidate {
	return d.candidatePool.GetChairpersons(state)
}
// Getting all refund array by nodeId
func (d *dpos) GetDefeat(state vm.StateDB, nodeId discover.NodeID) ([]*types.Candidate, error){
	return d.candidatePool.GetDefeat(state, nodeId)
}
// Checked current candidate was defeat by nodeId
func (d *dpos) IsDefeat(state vm.StateDB, nodeId discover.NodeID) (bool, error) {
	return d.candidatePool.IsDefeat(state, nodeId)
}

// refund once
func (d *dpos) RefundBalance (state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) error{
	return d.candidatePool.RefundBalance (state, nodeId, blockNumber)
}
// Getting owner's address of candidate info by nodeId
func (d *dpos) GetOwner (state vm.StateDB, nodeId discover.NodeID) common.Address {
	return d.candidatePool.GetOwner(state, nodeId)
}

// Getting allow block interval for refunds
func (d *dpos) GetRefundInterval () uint64 {
	return d.candidatePool.GetRefundInterval()
}

// cbft共识区块产生分叉后需要更新primaryNodeList和formerlyNodeList
func (d *dpos) UpdateNodeList (state *state.StateDB) {
	log.Warn("---cbft共识区块产生分叉，更新primaryNodeList和formerlyNodeList---", "state", state)
	d.lock.Lock()
	defer d.lock.Unlock()
	//formerNodes, err1 := d.GetWitness(state, -1)	// flag：-1: 上一轮	  0: 本轮见证人   1: 下一轮见证人
	//currentNodes, err2 := d.GetWitness(state, 0)	// flag：-1: 上一轮	  0: 本轮见证人   1: 下一轮见证人
	if preArr, curArr, nextArr, err := d.candidatePool.GetAllWitness(state); nil != err {
		log.Error("Load Witness from state failed on SetCandidatePool err", err)
		panic("UpdateNodeList error")
	}else {
		if len(preArr) != 0 {
			d.formerlyNodeList = convertNodeID(preArr)
		}
		if len(curArr) != 0 {
			d.primaryNodeList = convertNodeID(curArr)
		}
		if len(nextArr) != 0 {
			d.nextNodeList = convertNodeID(nextArr)
		}
	}
	//if err1 == nil && err2 == nil && len(formerNodes) > 0 && len(currentNodes) > 0 {
	//	d.primaryNodeList = convertNodeID(currentNodes)
	//	d.formerlyNodeList = convertNodeID(formerNodes)
	//} else {
	//	panic("UpdateNodeList error")
	//}
}

func convertNodeID(nodes []*discover.Node) []discover.NodeID {
	nodesID := make([]discover.NodeID, len(nodes), len(nodes))
	for _,n := range nodes {
		nodesID = append(nodesID, n.ID)
	}
	return nodesID
}
