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
	"fmt"
)

type ppos struct {
	//former            *pposRound // the previous round of witnesses nodeId
	//current           *pposRound // the current round of witnesses nodeId
	//next              *pposRound // the next round of witnesses nodeId

	nodeRound 		  roundCache
	//chain             *core.BlockChain
	lastCycleBlockNum uint64
	startTimeOfEpoch  int64 // 一轮共识开始时间，通常是上一轮共识结束时最后一个区块的出块时间；如果是第一轮，则从1970.1.1.0.0.0.0开始。单位：秒
	config            *params.PposConfig
	//initialNodes      []discover.Node
	// added by candidatepool module

	lock sync.RWMutex
	// the candidate pool object pointer
	candidatePool *pposm.CandidatePool
}


func newPpos(/*initialNodes []discover.Node, */config *params.CbftConfig) *ppos {
	//initNodeArr := make([]*discover.Node, 0, len(initialNodes))
	//initialNodesIDs := make([]discover.NodeID, 0, len(initialNodes))
	//for _, n := range config.InitialNodes {
	//	node := n
	//	initialNodesIDs = append(initialNodesIDs, node.ID)
	//	initNodeArr = append(initNodeArr, &node)
	//}
	//
	//formerRound := &pposRound{
	//	nodeIds: make([]discover.NodeID, 0),
	//	nodes: 	make([]*discover.Node, 0),
	//	start: big.NewInt(0),
	//	end:   big.NewInt(0),
	//}
	//currentRound := &pposRound{
	//	nodeIds: initialNodesIDs,
	//	//nodes: 	initNodeArr,
	//	start: big.NewInt(1),
	//	end:   big.NewInt(BaseSwitchWitness),
	//}
	//currentRound.nodes = make([]*discover.Node, len(initNodeArr))
	//copy(currentRound.nodes, initNodeArr)
	//
	//log.Info("初始化 ppos 当前轮配置节点:", "start", currentRound.start, "end", currentRound.end)
	//pposm.PrintObject("初始化 ppos 当前轮 nodeIds:", initialNodesIDs)
	//pposm.PrintObject("初始化 ppos 当前轮 nodes:", initNodeArr)
	return &ppos{
		//former:            formerRound,
		//current:           currentRound,
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

func (d *ppos) BlockProducerIndex(parentNumber *big.Int, parentHash common.Hash, commitNumber *big.Int, nodeID discover.NodeID) int64 {
	d.lock.RLock()
	defer d.lock.RUnlock()

	log.Warn("BlockProducerIndex", "parentNumber", parentNumber, "parentHash", parentHash, "commitNumber", commitNumber, "nodeID", nodeID)
	pposm.PrintObject("BlockProducerIndex nodeID", nodeID)

	currentRound := d.nodeRound.getCurrentRound(parentNumber, parentHash)
	if currentRound != nil {
		if commitNumber.Cmp(currentRound.start) >= 0 && commitNumber.Cmp(currentRound.end) <= 0 {
			for idx, nid := range currentRound.nodeIds {
				if nid == nodeID {
					return int64(idx)
				}
			}
		}
	}
	return -1


	//if number == 0 {
	//	for idx, node := range d.current.nodeIds {
	//		if node == nodeID {
	//			return int64(idx)
	//		}
	//	}
	//	return -1
	//}

	//if number >= d.former.start.Uint64() && number <= d.former.end.Uint64() {
	//	for idx, node := range d.former.nodeIds {
	//		if node == nodeID {
	//			return int64(idx)
	//		}
	//	}
	//	return -1
	//}
	//
	//if number >= d.current.start.Uint64() && number <= d.current.end.Uint64() {
	//	for idx, node := range d.current.nodeIds {
	//		if node == nodeID {
	//			return int64(idx)
	//		}
	//	}
	//	return -1
	//}
	//
	//if d.next != nil && number >= d.next.start.Uint64() && number <= d.next.end.Uint64() {
	//	for idx, node := range d.next.nodeIds {
	//		if node == nodeID {
	//			return int64(idx)
	//		}
	//	}
	//	return -1
	//}
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

//func (d *ppos) consensusNodes(blockNumber *big.Int) []discover.NodeID {
//	d.lock.RLock()
//	defer d.lock.RUnlock()
//
//	if d.former != nil && blockNumber.Cmp(d.former.start) >= 0 && blockNumber.Cmp(d.former.end) <= 0 {
//		return d.former.nodeIds
//	} else if d.current != nil && blockNumber.Cmp(d.current.start) >= 0 && blockNumber.Cmp(d.current.end) <= 0 {
//		return d.current.nodeIds
//	} else if d.next != nil && blockNumber.Cmp(d.next.start) >= 0 && blockNumber.Cmp(d.next.end) <= 0 {
//		return d.next.nodeIds
//	}
//	return nil
//}

func (d *ppos) consensusNodes(parentNumber *big.Int, parentHash common.Hash, commitNumber *big.Int) []discover.NodeID {
	d.lock.RLock()
	defer d.lock.RUnlock()

	nodeCache := d.nodeRound.getNodeCache(parentNumber, parentHash)
	if nodeCache != nil {
		if nodeCache.former != nil && commitNumber.Cmp(nodeCache.former.start) >= 0 && commitNumber.Cmp(nodeCache.former.end) <= 0 {
			return nodeCache.former.nodeIds
		} else if nodeCache.current != nil && commitNumber.Cmp(nodeCache.current.start) >= 0 && commitNumber.Cmp(nodeCache.current.end) <= 0 {
			return nodeCache.current.nodeIds
		} else if nodeCache.next != nil && commitNumber.Cmp(nodeCache.next.start) >= 0 && commitNumber.Cmp(nodeCache.next.end) <= 0 {
			return nodeCache.next.nodeIds
		}
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
		//round := calcurround(blocknumber)

		//d.lock.Lock()
		//log.Info("揭榜维护", "blockNumber:", blocknumber.Uint64(), "round:", round)
		//nextStart := big.NewInt(int64(BaseSwitchWitness*round) + 1)
		//nextEnd := new(big.Int).Add(nextStart, big.NewInt(int64(BaseSwitchWitness-1)))
		//d.next = &pposRound{
		//	nodeIds: convertNodeID(nextNodes),
		//	//nodes:	nextNodes,
		//	start: nextStart,
		//	end:   nextEnd,
		//}
		//d.next.nodes = make([]*discover.Node, len(nextNodes))
		//copy(d.next.nodes, nextNodes)
		//
		//log.Info("揭榜维护:下一轮", "start", d.next.start, "end", d.next.end)
		//log.Info("揭榜维护下一轮的nodeIds长度:", "len", len(nextNodes))
		//pposm.PrintObject("揭榜维护下一轮的nodeIds:", nextNodes)
		//pposm.PrintObject("揭榜的上轮pposRound：", d.former.nodes)
		//pposm.PrintObject("揭榜的当前轮pposRound：", d.current.nodes)
		//pposm.PrintObject("揭榜维护下一轮pposRound：", d.next.nodes)
		//d.lock.Unlock()
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
	/*_, curArr, _, err := */d.candidatePool.GetAllWitness(state)
	//if nil != err {
	//	return false
	//}
	//d.lock.Lock()

	//cur_start := d.current.start
	//cur_end :=  d.current.end
	//d.former.start = cur_start
	//d.former.end = cur_end
	//
	//next_start :=  d.next.start
	//next_end := d.next.end
	//d.current.start = next_start
	//d.current.end = next_end
	//d.former.nodeIds = convertNodeID(d.current.nodes)
	//d.former.nodes = make([]*discover.Node, len(d.current.nodes))
	//copy(d.former.nodes, d.current.nodes)
	//if len(curArr) != 0 {
	//	d.current.nodeIds = convertNodeID(curArr)
	//	d.current.nodes = make([]*discover.Node, len(curArr))
	//	copy(d.current.nodes, curArr)
	//}
	//
	//d.next = nil
	//log.Info("Switch获取:上一轮", "start", d.former.start, "end", d.former.end)
	//log.Info("Switch获取:当前轮", "start", d.current.start, "end", d.current.end)
	////log.Info("Switch获取:下一轮", "start", d.next.start, "end", d.next.end)
	////pposm.PrintObject("Switch获取上一轮nodes：", preArr)
	//
	//pposm.PrintObject("Switch获取当前轮nodes：", curArr)
	//pposm.PrintObject("Switch的上轮pposRound：", d.former.nodes)
	//pposm.PrintObject("Switch的当前轮pposRound：", d.current.nodes)

	//d.lock.Unlock()
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
func (d *ppos) SetCandidatePool(blockChain *core.BlockChain, initialNodes []discover.Node) {
	log.Info("---start node，to update nodeRound---")
	genesis := blockChain.Genesis()
	// init roundCache by config
	d.nodeRound = buildGenesisRound(genesis.NumberU64(), genesis.Hash(), initialNodes)
	// When the highest block in the chain is not a genesis block, Need to load witness nodeIdList from the stateDB.
	if genesis.NumberU64() != blockChain.CurrentBlock().NumberU64() {

		currentBlock := blockChain.CurrentBlock()
		var currBlockNumber uint64
		var currBlockHash common.Hash

		currBlockNumber = blockChain.CurrentBlock().NumberU64()
		currBlockHash = blockChain.CurrentBlock().Hash()

		//d.lock.Lock()
		//defer d.lock.Unlock()


		count := 0
		blockArr := make([]*types.Block, 0)
		for {
			if currBlockNumber == genesis.NumberU64() || count == BaseIrrCount {
				break
			}
			parentNum := currBlockNumber - 1
			parentHash := currentBlock.ParentHash()
			blockArr = append(blockArr, currentBlock)

			currBlockNumber = parentNum
			currBlockHash = parentHash
			currentBlock = blockChain.GetBlock(currBlockHash, currBlockNumber)
			count ++

		}

		for i := len(blockArr) - 1; 0 <= i; i-- {
			currentBlock := blockArr[i]
			currentNum := currentBlock.NumberU64()
			currentHash := currentBlock.Hash()

			parentNum := currentNum - 1
			parentHash := currentBlock.ParentHash()

			// 特殊处理数组最后一个块, 也就是最高块往前推第20个块
			if i == len(blockArr) - 1 && currentNum > 1  {

				var parent, current *state.StateDB

				// parentStateDB by block
				parentStateRoot := blockChain.GetBlock(parentHash, parentNum).Root()
				if parentState, err := blockChain.StateAt(parentStateRoot); nil != err {
					log.Error("Failed to load parentStateDB by block", "parentNum", parentNum, "Hash", parentHash.String(), "err", err)
					panic("Failed to load parentStateDB by block parentNum" + fmt.Sprint(parentNum) + ", Hash" + parentHash.String() + "err" + err.Error())
				}else {
					parent = parentState
				}

				// currentStateDB by block
				stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
				if currntState, err := blockChain.StateAt(stateRoot); nil != err {
					log.Error("Failed to load currentStateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
					panic("Failed to load currentStateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
				}else {
					current = currntState
				}

				if err := d.setEarliestIrrNodeCache(parent, current, genesis.NumberU64(), currentNum, genesis.Hash(), currentHash); nil != err {
					log.Error("Failed to setEarliestIrrNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
					panic("Failed to setEarliestIrrNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
				}
				continue
			}

			// stateDB by block
			stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
			if currntState, err := blockChain.StateAt(stateRoot); nil != err {
				log.Error("Failed to load stateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to load stateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}else {
				if err := d.setGeneralNodeCache(currntState, parentNum, currentNum, parentHash, currentHash); nil != err {
					log.Error("Failed to setGeneralNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
					panic("Failed to setGeneralNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
				}
			}
		}

		//for {
		//
		//	if blockNumber == genesis.NumberU64() || count == BaseIrrCount {
		//		break
		//	}
		//
		//	parentNum := blockNumber - 1
		//	parentHash := currentBlock.ParentHash()
		//
		//	// stateDB by block
		//	stateRoot := blockChain.GetBlock(blockHash, blockNumber).Root()
		//	if currntState, err := blockChain.StateAt(stateRoot); nil != err {
		//		log.Error("Failed to load stateDB by block", "blockNumber", blockNumber, "Hash", blockHash.String(), "err", err)
		//		panic("Failed to load stateDB by block blockNumber" + fmt.Sprint(blockNumber) + ", Hash" + blockHash.String() + "err" + err.Error())
		//	}else {
		//		if err := d.setNodeCache(currntState, parentNum, blockNumber, parentHash, blockHash); nil != err {
		//			log.Error("Failed to load stateDB by block", "blockNumber", blockNumber, "Hash", blockHash.String(), "err", err)
		//			panic("Failed to load stateDB by block blockNumber" + fmt.Sprint(blockNumber) + ", Hash" + blockHash.String() + "err" + err.Error())
		//		}
		//	}
		//	blockNumber = parentNum
		//	blockHash = parentHash
		//	currentBlock = blockChain.GetBlock(blockHash, blockNumber)
		//	count ++
		//}
	}
	pposm.PrintObject("启动node时, nodeRound:", d.nodeRound)
}


func buildGenesisRound(blockNumber uint64, blockHash common.Hash, initialNodes []discover.Node) roundCache {
	initNodeArr := make([]*discover.Node, 0, len(initialNodes))
	initialNodesIDs := make([]discover.NodeID, 0, len(initialNodes))
	for _, n := range initialNodes {
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


	log.Info("根据配置文件初始化 ppos 当前轮配置节点:", "blockNumber", blockNumber, "blockHash", blockHash.String(), "start", currentRound.start, "end", currentRound.end)
	pposm.PrintObject("初始化 ppos 当前轮 nodeIds:", initialNodesIDs)
	pposm.PrintObject("初始化 ppos 当前轮 nodes:", initNodeArr)

	node := &nodeCache{
		former: 	formerRound,
		current: 	currentRound,
	}

	nodeRound := make(roundCache, 0)
	hashRound := make(map[common.Hash]*nodeCache, 0)
	hashRound[blockHash] = node
	nodeRound[blockNumber] = hashRound
	return nodeRound
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

// cbft consensus fork need to update  nodeRound
func (d *ppos) UpdateNodeList(blockChain *core.BlockChain, blocknumber *big.Int, blockHash common.Hash) {
	log.Info("---cbft consensus fork，update nodeRound---")
	// clean nodeCache
	d.cleanNodeRound()


	var curBlockNumber uint64 = blocknumber.Uint64()
	var curBlockHash common.Hash = blockHash

	currentBlock := blockChain.GetBlock(curBlockHash, curBlockNumber)
	genesis := blockChain.Genesis()
	d.lock.Lock()
	defer d.lock.Unlock()

	count := 0
	blockArr := make([]*types.Block, 0)
	for {
		if curBlockNumber == genesis.NumberU64() || count == BaseIrrCount {
			break
		}
		parentNum := curBlockNumber - 1
		parentHash := currentBlock.ParentHash()
		blockArr = append(blockArr, currentBlock)

		curBlockNumber = parentNum
		curBlockHash = parentHash
		currentBlock = blockChain.GetBlock(curBlockHash, curBlockNumber)
		count ++

	}

	for i := len(blockArr) - 1; 0 <= i; i-- {
		currentBlock := blockArr[i]
		currentNum := currentBlock.NumberU64()
		currentHash := currentBlock.Hash()

		parentNum := currentNum - 1
		parentHash := currentBlock.ParentHash()


		// 特殊处理数组最后一个块, 也就是最高块往前推第20个块
		if i == len(blockArr) - 1 && currentNum > 1  {

			var parent, current *state.StateDB

			// parentStateDB by block
			parentStateRoot := blockChain.GetBlock(parentHash, parentNum).Root()
			if parentState, err := blockChain.StateAt(parentStateRoot); nil != err {
				log.Error("Failed to load parentStateDB by block", "parentNum", parentNum, "Hash", parentHash.String(), "err", err)
				panic("Failed to load parentStateDB by block parentNum" + fmt.Sprint(parentNum) + ", Hash" + parentHash.String() + "err" + err.Error())
			}else {
				parent = parentState
			}

			// currentStateDB by block
			stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
			if currntState, err := blockChain.StateAt(stateRoot); nil != err {
				log.Error("Failed to load currentStateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to load currentStateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}else {
				current = currntState
			}

			if err := d.setEarliestIrrNodeCache(parent, current, genesis.NumberU64(), currentNum, genesis.Hash(), currentHash); nil != err {
				log.Error("Failed to setEarliestIrrNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to setEarliestIrrNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}
			continue
		}

		// stateDB by block
		stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
		if currntState, err := blockChain.StateAt(stateRoot); nil != err {
			log.Error("Failed to load stateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
			panic("Failed to load stateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
		}else {
			if err := d.setGeneralNodeCache(currntState, parentNum, currentNum, parentHash, currentHash); nil != err {
				log.Error("Failed to setGeneralNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to setGeneralNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}
		}
	}


	//for {
	//
	//	if curBlockNumber == genesis.NumberU64() || count == BaseIrrCount {
	//		break
	//	}
	//
	//	parentNum := curBlockNumber - 1
	//	parentHash := currentBlock.ParentHash()
	//
	//	// stateDB by block
	//	stateRoot := blockChain.GetBlock(curBlockHash, curBlockNumber).Root()
	//	if currntState, err := blockChain.StateAt(stateRoot); nil != err {
	//		log.Error("Failed to load stateDB by block", "curBlockNumber", curBlockNumber, "Hash", curBlockHash.String(), "err", err)
	//		panic("Failed to load stateDB by block curBlockNumber" + fmt.Sprint(curBlockNumber) + ", Hash" + curBlockHash.String() + "err" + err.Error())
	//	}else {
	//		if err := d.setNodeCache(currntState, parentNum, curBlockNumber, parentHash, curBlockHash); nil != err {
	//			log.Error("Failed to load stateDB by block", "curBlockNumber", curBlockNumber, "Hash", curBlockHash.String(), "err", err)
	//			panic("Failed to load stateDB by block curBlockNumber" + fmt.Sprint(curBlockNumber) + ", Hash" + curBlockHash.String() + "err" + err.Error())
	//		}
	//	}
	//	currentBlock = blockChain.GetBlock(curBlockHash, curBlockNumber)
	//	curBlockNumber = parentNum
	//	curBlockHash = parentHash
	//	count ++
	//}
	pposm.PrintObject("分叉重载时, nodeRound:", d.nodeRound)
}

func convertNodeID(nodes []*discover.Node) []discover.NodeID {
	nodesID := make([]discover.NodeID, 0, len(nodes))
	for _, n := range nodes {
		nodesID = append(nodesID, n.ID)
	}
	return nodesID
}

// calculate current round number by current blocknumber
func calcurround(blocknumber uint64) uint64 {
	// current num
	var round uint64
	div := blocknumber / BaseSwitchWitness
	mod := blocknumber % BaseSwitchWitness
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


func (d *ppos) GetFormerRound(blockNumber *big.Int, blockHash common.Hash) *pposRound {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.nodeRound.getFormerRound(blockNumber, blockHash)
}

func (d *ppos) GetCurrentRound (blockNumber *big.Int, blockHash common.Hash) *pposRound {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.nodeRound.getCurrentRound(blockNumber, blockHash)
}

func (d *ppos)  GetNextRound (blockNumber *big.Int, blockHash common.Hash) *pposRound {
	d.lock.RLock()
	defer d.lock.RUnlock()
	return d.nodeRound.getNextRound(blockNumber, blockHash)
}

func (d *ppos) SetNodeCache (state *state.StateDB, parentNumber, currentNumber *big.Int, parentHash, currentHash common.Hash) error {
	d.lock.Lock()
	defer d.lock.Unlock()
	return d.setGeneralNodeCache(state, parentNumber.Uint64(), currentNumber.Uint64(), parentHash, currentHash)
}
func (d *ppos) setGeneralNodeCache (state *state.StateDB, parentNumber, currentNumber uint64, parentHash, currentHash common.Hash) error {
	parentNumBigInt := big.NewInt(int64(parentNumber))
	// current round
	round := calcurround(currentNumber)
	log.Info("设置当前区块", "currentNumber:", currentNumber, "round:", round)

	preNodes, curNodes, nextNodes, err := d.candidatePool.GetAllWitness(state)

	if nil != err {
		log.Error("Failed to setting nodeCache on setGeneralNodeCache", "err", err)
		return err
	}


	var start, end *big.Int

	// 判断是否是 本轮的最后一个块，如果是，则start 为下一轮的 start， end 为下一轮的 end
	if (currentNumber  % BaseSwitchWitness) == 0 {
		start = big.NewInt(int64(BaseSwitchWitness*round) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(BaseSwitchWitness-1)))
	}else {
		start = big.NewInt(int64(BaseSwitchWitness*(round-1)) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(BaseSwitchWitness-1)))
	}

	// former
	formerRound := &pposRound{}
	// former start, end
	if round != 1 {
		formerRound.start = new(big.Int).Sub(start, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
		formerRound.end = new(big.Int).Sub(end, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
	}
	log.Info("设置当前区块:上一轮", "start",formerRound.start, "end", formerRound.end)
	if len(preNodes) != 0 {
		formerRound.nodeIds = convertNodeID(preNodes)
		formerRound.nodes = make([]*discover.Node, len(preNodes))
		copy(formerRound.nodes, preNodes)
	}else { // Reference parent
		// if last block of round
		if (currentNumber % BaseSwitchWitness) == 0 {
			parentCurRound := d.nodeRound.getCurrentRound(parentNumBigInt, parentHash)
			if nil != parentCurRound {
				formerRound.nodeIds = make([]discover.NodeID, len(parentCurRound.nodeIds))
				copy(formerRound.nodeIds, parentCurRound.nodeIds)
				formerRound.nodes = make([]*discover.Node, len(parentCurRound.nodes))
				copy(formerRound.nodes, parentCurRound.nodes)
			}
		}else { // Is'nt last block of round
			parentFormerRound := d.nodeRound.getFormerRound(parentNumBigInt, parentHash)
			if nil != parentFormerRound {
				formerRound.nodeIds = make([]discover.NodeID, len(parentFormerRound.nodeIds))
				copy(formerRound.nodeIds, parentFormerRound.nodeIds)
				formerRound.nodes = make([]*discover.Node, len(parentFormerRound.nodes))
				copy(formerRound.nodes, parentFormerRound.nodes)
			}
		}
	}

	// current
	currentRound := &pposRound{}
	// current start, end
	currentRound.start = start
	currentRound.end = end
	log.Info("设置当前区块:当前轮", "start", currentRound.start, "end",currentRound.end)
	if len(curNodes) != 0 {
		currentRound.nodeIds = convertNodeID(curNodes)
		currentRound.nodes = make([]*discover.Node, len(curNodes))
		copy(currentRound.nodes, curNodes)
	}else { // Reference parent
		// if last block of round
		if (currentNumber % BaseSwitchWitness) == 0 {
			parentNextRound := d.nodeRound.getNextRound(parentNumBigInt, parentHash)
			if nil != parentNextRound {
				currentRound.nodeIds = make([]discover.NodeID, len(parentNextRound.nodeIds))
				copy(currentRound.nodeIds, parentNextRound.nodeIds)
				currentRound.nodes = make([]*discover.Node, len(parentNextRound.nodes))
				copy(currentRound.nodes, parentNextRound.nodes)
			}
		}else { // Is'nt last block of round
			parentCurRound := d.nodeRound.getCurrentRound(parentNumBigInt, parentHash)
			if nil != parentCurRound {
				currentRound.nodeIds = make([]discover.NodeID, len(parentCurRound.nodeIds))
				copy(currentRound.nodeIds, parentCurRound.nodeIds)
				currentRound.nodes = make([]*discover.Node, len(parentCurRound.nodes))
				copy(currentRound.nodes, parentCurRound.nodes)
			}
		}
	}


	// next
	nextRound := &pposRound{}
	// next start, end
	nextRound.start = new(big.Int).Add(start, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
	nextRound.end = new(big.Int).Add(end, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
	log.Info("设置当前区块:下一轮", "start", nextRound.start, "end",nextRound.end)
	if len(nextNodes) != 0 {
		nextRound.nodeIds = convertNodeID(nextNodes)
		nextRound.nodes = make([]*discover.Node, len(nextNodes))
		copy(nextRound.nodes, nextNodes)
	}else { // Reference parent

		if (currentNumber % BaseElection) == 0 { // election index == cur index
			parentCurRound := d.nodeRound.getCurrentRound(parentNumBigInt, parentHash)
			if nil != parentCurRound {
				nextRound.nodeIds = make([]discover.NodeID, len(parentCurRound.nodeIds))
				copy(nextRound.nodeIds, parentCurRound.nodeIds)
				nextRound.nodes = make([]*discover.Node, len(parentCurRound.nodes))
				copy(nextRound.nodes, parentCurRound.nodes)
			}
		}else if (currentNumber % BaseElection) != 0 && (currentNumber / BaseElection) == round {  // election index < cur index < switch index
			parentNextRound := d.nodeRound.getNextRound(parentNumBigInt, parentHash)
			if nil != parentNextRound {
				nextRound.nodeIds = make([]discover.NodeID, len(parentNextRound.nodeIds))
				copy(nextRound.nodeIds, parentNextRound.nodeIds)
				nextRound.nodes = make([]*discover.Node, len(parentNextRound.nodes))
				copy(nextRound.nodes, parentNextRound.nodes)
			}
		}else { // switch index <= cur index < next election index
			nextRound.nodeIds = make([]discover.NodeID, 0)
			nextRound.nodes = make([]*discover.Node, 0)
		}
	}

	pposm.PrintObject("设置当前区块 stateDB 上一轮nodes：", preNodes)
	pposm.PrintObject("设置当前区块 stateDB 当前轮nodes：", curNodes)
	pposm.PrintObject("设置当前区块 stateDB 下一轮nodes：", nextNodes)
	pposm.PrintObject("设置当前区块的上轮pposRound：", formerRound.nodes)
	pposm.PrintObject("设置当前区块的当前轮pposRound：", currentRound.nodes)
	pposm.PrintObject("设置当前区块的下一轮pposRound：", nextRound.nodes)

	cache := &nodeCache{
		former: 	formerRound,
		current: 	currentRound,
		next: 		nextRound,
	}
	d.nodeRound.setNodeCache(big.NewInt(int64(currentNumber)), currentHash, cache)
	log.Info("设置当前区块的信息时", "currentBlockNum", currentNumber, "parentNum", parentNumber, "currentHash", currentHash.String(), "parentHash", parentHash.String())
	pposm.PrintObject("设置当前区块的信息时, nodeRound:", d.nodeRound)
	return nil
}

func (d *ppos) setEarliestIrrNodeCache (parentState, currentState *state.StateDB, genesisNumber, currentNumber uint64, genesisHash, currentHash common.Hash) error {
	genesisNumBigInt := big.NewInt(int64(genesisNumber))
	// current round
	round := calcurround(currentNumber)
	log.Info("设置最远允许缓存保留区块", "currentNumber:", currentNumber, "round:", round)

	curr_preNodes, curr_curNodes, curr_nextNodes, err := d.candidatePool.GetAllWitness(currentState)

	if nil != err {
		log.Error("Failed to setting nodeCache by currentStateDB on setEarliestIrrNodeCache", "err", err)
		return err
	}

	parent_preNodes, parent_curNodes, parent_nextNodes, err := d.candidatePool.GetAllWitness(parentState)
	if nil != err {
		log.Error("Failed to setting nodeCache by parentStateDB on setEarliestIrrNodeCache", "err", err)
		return err
	}


	var start, end *big.Int

	// 判断是否是 本轮的最后一个块，如果是，则start 为下一轮的 start， end 为下一轮的 end
	if (currentNumber  % BaseSwitchWitness) == 0 {
		start = big.NewInt(int64(BaseSwitchWitness*round) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(BaseSwitchWitness-1)))
	}else {
		start = big.NewInt(int64(BaseSwitchWitness*(round-1)) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(BaseSwitchWitness-1)))
	}

	// former
	formerRound := &pposRound{}
	// former start, end
	if round != 1 {
		formerRound.start = new(big.Int).Sub(start, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
		formerRound.end = new(big.Int).Sub(end, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
	}
	log.Info("设置最远允许缓存保留区块:上一轮", "start",formerRound.start, "end", formerRound.end)
	if len(curr_preNodes) != 0 {
		formerRound.nodeIds = convertNodeID(curr_preNodes)
		formerRound.nodes = make([]*discover.Node, len(curr_preNodes))
		copy(formerRound.nodes, curr_preNodes)
	}else { // Reference parent
		// if last block of round
		if (currentNumber % BaseSwitchWitness) == 0 {
			// 先从上一个块的stateDB拿, 上一个块的stateDB 也没有，就从对应着创世块的 nodeCache拿
			if len(parent_curNodes) != 0 {
				//formerRound.nodeIds = make([]discover.NodeID, len(parent_curNodes))
				formerRound.nodeIds = convertNodeID(parent_curNodes)
				formerRound.nodes = make([]*discover.Node, len(parent_curNodes))
				copy(formerRound.nodes, parent_curNodes)
			}else {
				genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash)
				if nil != genesisCurRound {
					formerRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
					copy(formerRound.nodeIds, genesisCurRound.nodeIds)
					formerRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
					copy(formerRound.nodes, genesisCurRound.nodes)
				}
			}
		}else { // Is'nt last block of round

			if len(parent_preNodes) != 0 {
				//formerRound.nodeIds = make([]discover.NodeID, len(parent_preNodes))
				formerRound.nodeIds = convertNodeID(parent_preNodes)
				formerRound.nodes = make([]*discover.Node, len(parent_preNodes))
				copy(formerRound.nodes, parent_preNodes)
			}else {
				genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash)
				if nil != genesisCurRound {
					formerRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
					copy(formerRound.nodeIds, genesisCurRound.nodeIds)
					formerRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
					copy(formerRound.nodes, genesisCurRound.nodes)
				}
			}
		}
	}

	// current
	currentRound := &pposRound{}
	// current start, end
	currentRound.start = start
	currentRound.end = end
	log.Info("设置最远允许缓存保留区块:当前轮", "start", currentRound.start, "end",currentRound.end)
	if len(curr_curNodes) != 0 {
		currentRound.nodeIds = convertNodeID(curr_curNodes)
		currentRound.nodes = make([]*discover.Node, len(curr_curNodes))
		copy(currentRound.nodes, curr_curNodes)
	}else { // Reference parent
		// if last block of round
		if (currentNumber % BaseSwitchWitness) == 0 {
			if len(parent_nextNodes) != 0  {
				currentRound.nodeIds = convertNodeID(parent_nextNodes)
				currentRound.nodes = make([]*discover.Node, len(parent_nextNodes))
				copy(currentRound.nodes, parent_nextNodes)
			}else {
				genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash)
				if nil != genesisCurRound {
					currentRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
					copy(currentRound.nodeIds, genesisCurRound.nodeIds)
					currentRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
					copy(currentRound.nodes, genesisCurRound.nodes)
				}
			}
		}else { // Is'nt last block of round

			if len(parent_curNodes) != 0 {
				currentRound.nodeIds = convertNodeID(parent_curNodes)
				currentRound.nodes = make([]*discover.Node, len(parent_curNodes))
				copy(currentRound.nodes, parent_curNodes)
			}else {
				genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash)
				if nil != genesisCurRound {
					currentRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
					copy(currentRound.nodeIds, genesisCurRound.nodeIds)
					currentRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
					copy(currentRound.nodes, genesisCurRound.nodes)
				}
			}
		}
	}


	// next
	nextRound := &pposRound{}
	// next start, end
	nextRound.start = new(big.Int).Add(start, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
	nextRound.end = new(big.Int).Add(end, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
	log.Info("设置最远允许缓存保留区块:下一轮", "start", nextRound.start, "end",nextRound.end)
	if len(curr_nextNodes) != 0 {
		nextRound.nodeIds = convertNodeID(curr_nextNodes)
		nextRound.nodes = make([]*discover.Node, len(curr_nextNodes))
		copy(nextRound.nodes, curr_nextNodes)
	}else { // Reference parent

		if (currentNumber % BaseElection) == 0 || ((currentNumber % BaseElection) != 0 && (currentNumber / BaseElection) == round) { // election index == cur index || election index < cur index < switch index

			genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash)
			if nil != genesisCurRound {
				nextRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
				copy(nextRound.nodeIds, genesisCurRound.nodeIds)
				nextRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
				copy(nextRound.nodes, genesisCurRound.nodes)
			}
		}else { // switch index <= cur index < next election index
			nextRound.nodeIds = make([]discover.NodeID, 0)
			nextRound.nodes = make([]*discover.Node, 0)
		}
	}

	pposm.PrintObject("设置最远允许缓存保留区块 stateDB 上一轮nodes：", curr_preNodes)
	pposm.PrintObject("设置最远允许缓存保留区块 stateDB 当前轮nodes：", curr_curNodes)
	pposm.PrintObject("设置最远允许缓存保留区块 stateDB 下一轮nodes：", curr_nextNodes)

	pposm.PrintObject("设置最远允许缓存保留区块  parentStateDB 上一轮nodes：", curr_preNodes)
	pposm.PrintObject("设置最远允许缓存保留区块 parentStateDB 当前轮nodes：", curr_curNodes)
	pposm.PrintObject("设置最远允许缓存保留区块 parentStateDB 下一轮nodes：", curr_nextNodes)

	pposm.PrintObject("设置最远允许缓存保留区块的上轮pposRound：", formerRound.nodes)
	pposm.PrintObject("设置最远允许缓存保留区块的当前轮pposRound：", currentRound.nodes)
	pposm.PrintObject("设置最远允许缓存保留区块的下一轮pposRound：", nextRound.nodes)

	cache := &nodeCache{
		former: 	formerRound,
		current: 	currentRound,
		next: 		nextRound,
	}
	d.nodeRound.setNodeCache(big.NewInt(int64(currentNumber)), currentHash, cache)
	log.Info("设置最远允许缓存保留区块的信息时", "currentBlockNum", currentNumber, "currentHash", currentHash.String())
	pposm.PrintObject("设置最远允许缓存保留区块的信息时, nodeRound:", d.nodeRound)
	return nil
}


func (d *ppos) cleanNodeRound () {
	d.lock.Lock()
	d.nodeRound =  make(roundCache, 0)
	d.lock.Unlock()
}