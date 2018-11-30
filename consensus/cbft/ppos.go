package cbft

import (
	"Platon-go/common"
	"Platon-go/core"
	"Platon-go/core/ppos"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"Platon-go/core/vm"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"math/big"
	"sync"
)

type ppos struct {
	former            *pposRound // the previous round of witnesses nodeId
	current           *pposRound // the current round of witnesses nodeId
	next              *pposRound // the next round of witnesses nodeId
	chain             *core.BlockChain
	lastCycleBlockNum uint64
	startTimeOfEpoch  int64 // 一轮共识开始时间，通常是上一轮共识结束时最后一个区块的出块时间；如果是第一轮，则从1970.1.1.0.0.0.0开始。单位：秒
	config            *params.PposConfig
	//initialNodes      []discover.Node
	// added by candidatepool module

	lock sync.RWMutex
	// the candidate pool object pointer
	candidatePool *pposm.CandidatePool
}

type pposRound struct {
	nodeIds []discover.NodeID
	nodes 	[]*discover.Node
	start *big.Int
	end   *big.Int
}

func newPpos(initialNodes []discover.Node, config *params.CbftConfig) *ppos {
	initNodeArr := make([]*discover.Node, 0, len(initialNodes))
	initialNodesIDs := make([]discover.NodeID, 0, len(initialNodes))
	for _, n := range config.InitialNodes {
		node := n
		initialNodesIDs = append(initialNodesIDs, node.ID)
		initNodeArr = append(initNodeArr, &node)
	}

	formerRound := &pposRound{
		nodeIds: make([]discover.NodeID, 0),
		nodes: 	make([]*discover.Node, 0),
		start: big.NewInt(0),
		end:   big.NewInt(0),
	}
	currentRound := &pposRound{
		nodeIds: initialNodesIDs,
		//nodes: 	initNodeArr,
		start: big.NewInt(1),
		end:   big.NewInt(BaseSwitchWitness),
	}
	currentRound.nodes = make([]*discover.Node, len(initNodeArr))
	copy(currentRound.nodes, initNodeArr)

	log.Info("初始化 ppos 当前轮配置节点:", "start", currentRound.start, "end", currentRound.end)
	pposm.PrintObject("初始化 ppos 当前轮 nodeIds:", initialNodesIDs)
	pposm.PrintObject("初始化 ppos 当前轮 nodes:", initNodeArr)
	return &ppos{
		former:            formerRound,
		current:           currentRound,
		lastCycleBlockNum: 0,
		config:            config.PposConfig,
		//initialNodes: 	   config.InitialNodes,
		candidatePool:     pposm.NewCandidatePool(config.PposConfig),
	}
	//return pposPtr
}

func (d *ppos) AnyIndex(nodeID discover.NodeID) int64 {
	d.lock.RLock()
	defer d.lock.RUnlock()
	nodeList := make([]discover.NodeID, 0)
	if d.former != nil && d.former.nodes != nil && len(d.former.nodes) > 0 {
		nodeList = append(nodeList, d.former.nodeIds...)
	}
	if d.current != nil && d.current.nodes != nil && len(d.current.nodes) > 0 {
		nodeList = append(nodeList, d.current.nodeIds...)
	}
	if d.next != nil && d.next.nodes != nil && len(d.next.nodes) > 0 {
		nodeList = append(nodeList, d.next.nodeIds...)
	}
	for idx, node := range nodeList {
		if node == nodeID {
			return int64(idx)
		}
	}
	return int64(-1)
}

func (d *ppos) BlockProducerIndex(number uint64, nodeID discover.NodeID) int64 {
	d.lock.RLock()
	defer d.lock.RUnlock()

	log.Warn("BlockProducerIndex", "number", number, "nodeID", nodeID)
	pposm.PrintObject("BlockProducerIndex nodeID", nodeID)

	pposm.PrintObject("former nodes", d.former.nodes)
	pposm.PrintObject("former start", d.former.start)
	pposm.PrintObject("former end", d.former.end)

	pposm.PrintObject("current nodes", d.current.nodes)
	pposm.PrintObject("current start", d.current.start)
	pposm.PrintObject("current end", d.current.end)

	if d.next != nil {
		pposm.PrintObject("next nodes", d.next.nodes)
		pposm.PrintObject("next start", d.next.start)
		pposm.PrintObject("next end", d.next.end)
	}

	if number == 0 {
		for idx, node := range d.current.nodeIds {
			if node == nodeID {
				return int64(idx)
			}
		}
		return -1
	}

	if number >= d.former.start.Uint64() && number <= d.former.end.Uint64() {
		for idx, node := range d.former.nodeIds {
			if node == nodeID {
				return int64(idx)
			}
		}
		return -1
	}

	if number >= d.current.start.Uint64() && number <= d.current.end.Uint64() {
		for idx, node := range d.current.nodeIds {
			if node == nodeID {
				return int64(idx)
			}
		}
		return -1
	}

	if d.next != nil && number >= d.next.start.Uint64() && number <= d.next.end.Uint64() {
		for idx, node := range d.next.nodeIds {
			if node == nodeID {
				return int64(idx)
			}
		}
		return -1
	}
	return -1

}

func (d *ppos) NodeIndexInFuture(nodeID discover.NodeID) int64 {
	d.lock.RLock()
	defer d.lock.RUnlock()
	nodeList := append(d.current.nodeIds, d.next.nodeIds...)
	for idx, node := range nodeList {
		if node == nodeID {
			return int64(idx)
		}
	}
	return int64(-1)
}

func (d *ppos) getFormerNodeID () []discover.NodeID {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.former.nodeIds
}

func (d *ppos) getCurrentNodeID() []discover.NodeID {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.current.nodeIds
}

func (d *ppos) getNextNodeID () []discover.NodeID {
	d.lock.RLock()
	defer d.lock.RUnlock()
	if nil != d.next {
		return d.next.nodeIds
	}else {
		return make([]discover.NodeID, 0)
	}
}

func (d *ppos) getFormerNodes () []*discover.Node {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.former.nodes
}

func (d *ppos) getCurrentNodes () []*discover.Node {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.current.nodes
}

func (d *ppos) consensusNodes(blockNumber *big.Int) []discover.NodeID {
	d.lock.RLock()
	defer d.lock.RUnlock()

	if d.former != nil && blockNumber.Cmp(d.former.start) >= 0 && blockNumber.Cmp(d.former.end) <= 0 {
		return d.former.nodeIds
	} else if d.current != nil && blockNumber.Cmp(d.current.start) >= 0 && blockNumber.Cmp(d.current.end) <= 0 {
		return d.current.nodeIds
	} else if d.next != nil && blockNumber.Cmp(d.next.start) >= 0 && blockNumber.Cmp(d.next.end) <= 0 {
		return d.next.nodeIds
	}
	return nil
}

func (d *ppos) LastCycleBlockNum() uint64 {
	// 获取最后一轮共识结束时的区块高度
	return d.lastCycleBlockNum
}

func (d *ppos) SetLastCycleBlockNum(blockNumber uint64) {
	// 设置最后一轮共识结束时的区块高度
	d.lastCycleBlockNum = blockNumber
}

// modify by platon
// 返回当前共识节点地址列表
/*func (b *ppos) ConsensusNodes() []discover.Node {
	return b.primaryNodeList
}
*/
// 判断某个节点是否本轮或上一轮选举共识节点
/*func (b *ppos) CheckConsensusNode(id discover.NodeID) bool {
	nodes := b.ConsensusNodes()
	for _, node := range nodes {
		if node.ID == id {
			return true
		}
	}
	return false
}*/

// 判断当前节点是否本轮或上一轮选举共识节点
/*func (b *ppos) IsConsensusNode() (bool, error) {
	return true, nil
}
*/

func (d *ppos) StartTimeOfEpoch() int64 {
	return d.startTimeOfEpoch
}

func (d *ppos) SetStartTimeOfEpoch(startTimeOfEpoch int64) {
	// 设置最后一轮共识结束时的出块时间
	d.startTimeOfEpoch = startTimeOfEpoch
	log.Info("设置最后一轮共识结束时的出块时间", "startTimeOfEpoch", startTimeOfEpoch)
}

/** ppos was added func */
/** Method provided to the cbft module call */
// Announce witness
func (d *ppos) Election(state *state.StateDB, blocknumber *big.Int) ([]*discover.Node, error) {
	if nextNodes, err := d.candidatePool.Election(state); nil != err {
		log.Error("ppos election next witness err", err)
		panic("Election error " + err.Error())
	} else {
		log.Info("揭榜完成，再次查看stateDB信息...")
		d.candidatePool.GetAllWitness(state)
		// current round
		round := calcurround(blocknumber)

		d.lock.Lock()
		log.Info("揭榜维护", "blockNumber:", blocknumber.Uint64(), "round:", round)
		nextStart := big.NewInt(int64(BaseSwitchWitness*round) + 1)
		nextEnd := new(big.Int).Add(nextStart, big.NewInt(int64(BaseSwitchWitness-1)))
		d.next = &pposRound{
			nodeIds: convertNodeID(nextNodes),
			//nodes:	nextNodes,
			start: nextStart,
			end:   nextEnd,
		}
		d.next.nodes = make([]*discover.Node, len(nextNodes))
		copy(d.next.nodes, nextNodes)

		log.Info("揭榜维护:下一轮", "start", d.next.start, "end", d.next.end)
		log.Info("揭榜维护下一轮的nodeIds长度:", "len", len(nextNodes))
		pposm.PrintObject("揭榜维护下一轮的nodeIds:", nextNodes)
		pposm.PrintObject("揭榜的上轮pposRound：", d.former.nodes)
		pposm.PrintObject("揭榜的当前轮pposRound：", d.current.nodes)
		pposm.PrintObject("揭榜维护下一轮pposRound：", d.next.nodes)
		d.lock.Unlock()
		return nextNodes, nil
	}
}

// switch next witnesses to current witnesses
func (d *ppos) Switch(state *state.StateDB) bool {
	log.Info("Switch begin...")
	if !d.candidatePool.Switch(state) {
		return false
	}
	log.Info("Switch success...")
	_, curArr, _, err := d.candidatePool.GetAllWitness(state)
	if nil != err {
		return false
	}
	d.lock.Lock()

	cur_start := d.current.start
	cur_end :=  d.current.end
	d.former.start = cur_start
	d.former.end = cur_end

	next_start :=  d.next.start
	next_end := d.next.end
	d.current.start = next_start
	d.current.end = next_end
	d.former.nodeIds = convertNodeID(d.current.nodes)
	d.former.nodes = make([]*discover.Node, len(d.current.nodes))
	copy(d.former.nodes, d.current.nodes)
	if len(curArr) != 0 {
		d.current.nodeIds = convertNodeID(curArr)
		d.current.nodes = make([]*discover.Node, len(curArr))
		copy(d.current.nodes, curArr)
	}

	d.next = nil
	log.Info("Switch获取:上一轮", "start", d.former.start, "end", d.former.end)
	log.Info("Switch获取:当前轮", "start", d.current.start, "end", d.current.end)
	//log.Info("Switch获取:下一轮", "start", d.next.start, "end", d.next.end)
	//pposm.PrintObject("Switch获取上一轮nodes：", preArr)

	pposm.PrintObject("Switch获取当前轮nodes：", curArr)
	pposm.PrintObject("Switch的上轮pposRound：", d.former.nodes)
	pposm.PrintObject("Switch的当前轮pposRound：", d.current.nodes)

	d.lock.Unlock()
	return true
}

// Getting nodes of witnesses
// flag：-1: the previous round of witnesses  0: the current round of witnesses   1: the next round of witnesses
func (d *ppos) GetWitness(state *state.StateDB, flag int) ([]*discover.Node, error) {
	return d.candidatePool.GetWitness(state, flag)
}

func (d *ppos) GetAllWitness(state *state.StateDB) ([]*discover.Node, []*discover.Node, []*discover.Node, error) {
	return d.candidatePool.GetAllWitness(state)
}

// setting candidate pool of ppos module
func (d *ppos) SetCandidatePool(blockChain *core.BlockChain) {
	// When the highest block in the chain is not a genesis block, Need to load witness nodeIdList from the stateDB.
	if blockChain.Genesis().NumberU64() != blockChain.CurrentBlock().NumberU64() {
		state, err := blockChain.State()
		log.Warn("---重新启动节点，更新formerlyNodeList、primaryNodeList和nextNodeList---", "state", state)
		if nil != err {
			log.Error("Load state from chain failed on SetCandidatePool err", err)
			return
		}
		if preArr, curArr, nextArr, err := d.candidatePool.GetAllWitness(state); nil != err {
			log.Error("Load Witness from state failed on SetCandidatePool err", err)
		} else {
			d.lock.Lock()
			// current round
			round := calcurround(blockChain.CurrentBlock().Number())
			log.Info("重新加载", "blockNumber:", blockChain.CurrentBlock().Number(), "round:", round)
			// is'nt first round
			if round != 1 {
				d.former.start = big.NewInt(int64(BaseSwitchWitness*(round-2)) + 1)
				d.former.end = new(big.Int).Add(d.former.start, big.NewInt(int64(BaseSwitchWitness-1)))
			}
			log.Info("重新加载:上一轮", "start", d.former.start, "end", d.former.end)
			if len(preArr) != 0 {
				d.former.nodeIds = convertNodeID(preArr)
				d.former.nodes = make([]*discover.Node, len(preArr))
				copy(d.former.nodes, preArr)
				//d.former.nodes = preArr
			}else {
				if round == 1 {
					d.former.nodeIds = convertNodeID(d.current.nodes)
					d.former.nodes = make([]*discover.Node, len(d.current.nodes))
					copy(d.former.nodes, d.current.nodes)
				}
				//d.former.nodes = d.current.nodes
			}

			d.current.start = big.NewInt(int64(BaseSwitchWitness*(round-1)) + 1)
			d.current.end = new(big.Int).Add(d.current.start, big.NewInt(int64(BaseSwitchWitness-1)))
			log.Info("重新加载:当前轮", "start", d.current.start, "end", d.current.end)
			if len(curArr) != 0 {
				d.current.nodeIds = convertNodeID(curArr)
				d.current.nodes = make([]*discover.Node, len(curArr))
				copy(d.current.nodes, curArr)
				//d.current.nodes = curArr
			}
			if len(nextArr) != 0 {
				start := big.NewInt(int64(BaseSwitchWitness*round) + 1)
				end := new(big.Int).Add(start, big.NewInt(int64(BaseSwitchWitness-1)))

				d.next = &pposRound{
					nodeIds: 	convertNodeID(nextArr),
					//nodes: 		nextArr,
					start: 		start,
					end: 		end,
				}
				d.next.nodes = make([]*discover.Node, len(nextArr))
				copy(d.next.nodes, nextArr)

				log.Info("重新加载:下一轮", "start", d.next.start, "end", d.next.end)
				pposm.PrintObject("重新加载获取下一轮nodes：", nextArr)
				pposm.PrintObject("重新加载的下一轮pposRound：", d.next.nodes)
			}
			pposm.PrintObject("重新加载获取上一轮nodes：", preArr)
			pposm.PrintObject("重新加载获取当前轮nodes：", curArr)
			pposm.PrintObject("重新加载的上轮pposRound：", d.former.nodes)
			pposm.PrintObject("重新加载的当前轮pposRound：", d.current.nodes)

			d.lock.Unlock()
		}
	}/*else { // if current block is genesis , loading config nodes into stateDB


	}*/
}

/** Method provided to the built-in contract call */
// pledge Candidate
func (d *ppos) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	return d.candidatePool.SetCandidate(state, nodeId, can)
}

// Getting immediate candidate info by nodeId
func (d *ppos) GetCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error) {
	return d.candidatePool.GetCandidate(state, nodeId)
}

// candidate withdraw from immediates elected candidates
func (d *ppos) WithdrawCandidate(state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error {
	return d.candidatePool.WithdrawCandidate(state, nodeId, price, blockNumber)
}

// Getting all immediate elected candidates array
func (d *ppos) GetChosens(state vm.StateDB) []*types.Candidate {
	return d.candidatePool.GetChosens(state)
}

// Getting all witness array
func (d *ppos) GetChairpersons(state vm.StateDB) []*types.Candidate {
	return d.candidatePool.GetChairpersons(state)
}

// Getting all refund array by nodeId
func (d *ppos) GetDefeat(state vm.StateDB, nodeId discover.NodeID) ([]*types.Candidate, error) {
	return d.candidatePool.GetDefeat(state, nodeId)
}

// Checked current candidate was defeat by nodeId
func (d *ppos) IsDefeat(state vm.StateDB, nodeId discover.NodeID) (bool, error) {
	return d.candidatePool.IsDefeat(state, nodeId)
}

// refund once
func (d *ppos) RefundBalance(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) error {
	return d.candidatePool.RefundBalance(state, nodeId, blockNumber)
}

// Getting owner's address of candidate info by nodeId
func (d *ppos) GetOwner(state vm.StateDB, nodeId discover.NodeID) common.Address {
	return d.candidatePool.GetOwner(state, nodeId)
}

// Getting allow block interval for refunds
func (d *ppos) GetRefundInterval() uint64 {
	return d.candidatePool.GetRefundInterval()
}

// cbft共识区块产生分叉后需要更新primaryNodeList和formerlyNodeList
func (d *ppos) UpdateNodeList(state *state.StateDB, blocknumber *big.Int) {
	log.Warn("---cbft共识区块产生分叉，更新formerlyNodeList、primaryNodeList和nextNodeList---", "state", state)
	if preArr, curArr, _, err := d.candidatePool.GetAllWitness(state); nil != err {
		log.Error("Load Witness from state failed on UpdateNodeList err", err)
		panic("UpdateNodeList error")
	} else {
		d.lock.Lock()

		// current round
		round := calcurround(blocknumber)
		log.Info("分叉获取", "blockNumber:", blocknumber.Uint64(), "round:", round)
		start := big.NewInt(int64(BaseSwitchWitness*(round-1)) + 1)
		end := new(big.Int).Add(d.current.start, big.NewInt(int64(BaseSwitchWitness-1)))

		// is'nt first round
		if round != 1 {
			d.former.start = new(big.Int).Sub(start, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
			d.former.end = new(big.Int).Sub(end, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
		}
		log.Info("分叉获取:上一轮", "start", d.former.start, "end", d.former.end)
		if len(preArr) != 0 {
			//d.former = &pposRound{
			//	nodeIds: convertNodeID(preArr),
			//	//nodes: 	preArr,
			//}
			d.former.nodeIds = convertNodeID(preArr)
			d.former.nodes = make([]*discover.Node, len(preArr))
			copy(d.former.nodes, preArr)
		}

		d.current.start = start
		d.current.end = end
		log.Info("分叉获取:当前轮", "start", d.current.start, "end", d.current.end)
		if len(curArr) != 0 {
			//d.current = &pposRound{
			//	nodeIds: convertNodeID(curArr),
			//	//nodes: 	curArr,
			//}
			d.current.nodeIds = convertNodeID(curArr)
			d.current.nodes = make([]*discover.Node, len(curArr))
			copy(d.current.nodes, curArr)
		}
		d.next = nil
		pposm.PrintObject("分叉获取上一轮nodes：", preArr)
		pposm.PrintObject("分叉获取当前轮nodes：", curArr)
		pposm.PrintObject("分叉的上轮pposRound：", d.former.nodes)
		pposm.PrintObject("分叉的当前轮pposRound：", d.current.nodes)
		d.lock.Unlock()
	}
}

func convertNodeID(nodes []*discover.Node) []discover.NodeID {
	nodesID := make([]discover.NodeID, 0, len(nodes))
	for _, n := range nodes {
		nodesID = append(nodesID, n.ID)
	}
	return nodesID
}

// calculate current round number by current blocknumber
func calcurround(blocknumber *big.Int) uint64 {
	// current num
	var round uint64
	div := blocknumber.Uint64() / BaseSwitchWitness
	mod := blocknumber.Uint64() % BaseSwitchWitness
	if (div == 0 && mod == 0) || (div == 0 && mod > 0 && mod < BaseSwitchWitness) { // first round
		round = 1
	} else if div > 0 && mod == 0 {
		round = div
	} else if div > 0 && mod > 0 && mod < BaseSwitchWitness {
		round = div + 1
	}
	return round
}

//func (d *ppos) MaxChair() int64 {
//	return int64(d.candidatePool.MaxChair())
//}
