package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/ppos"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	//"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"fmt"
	"math/big"
	"sync"
	"github.com/PlatONnetwork/PlatON-Go/core/ticketcache"
)

type ppos struct {

	nodeRound 		  roundCache
	//chain             *core.BlockChain
	lastCycleBlockNum uint64
	// A round of consensus start time is usually the block time of the last block at the end of the last round of consensus;
	// if it is the first round, it starts from 1970.1.1.0.0.0.0. Unit: second
	startTimeOfEpoch  int64
	config            *params.PposConfig

	// added by candidatepool module
	lock 					sync.RWMutex
	// the candidate pool object pointer
	candidatePool 			*pposm.CandidatePool
	// the ticket pool object pointer
	ticketPool				*pposm.TicketPool
	// the ticket id list cache
	ticketidsCache 			*ticketcache.NumBlocks
}



func newPpos(config *params.CbftConfig) *ppos {
	return &ppos{
		lastCycleBlockNum: 	0,
		config:            	config.PposConfig,
		candidatePool:     	pposm.NewCandidatePool(config.PposConfig),
		ticketPool: 		pposm.NewTicketPool(config.PposConfig),
	}
}


func (d *ppos) NodeIndex(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, nodeID discover.NodeID, round int32) int64 {
	d.lock.RLock()
	defer d.lock.RUnlock()

	log.Warn("BlockProducerIndex", "parentNumber", parentNumber, "parentHash", parentHash.String(), "blockNumber", blockNumber.String(), "nodeID", nodeID.String(), "round", round)

	nodeCache := d.nodeRound.getNodeCache(parentNumber, parentHash)
	d.printMapInfo("BlockProducerIndex", parentNumber.Uint64(), parentHash)
	if nodeCache != nil {
		_former := nodeCache.former
		_current := nodeCache.current
		_next := nodeCache.next

		switch round {
			case former:
				if _former != nil && _former.start != nil && _former.end != nil && blockNumber.Cmp(_former.start) >= 0 && blockNumber.Cmp(_former.end) <= 0 {
					return d.roundIndex(nodeID, _former)
				}

			case current:
				if _current != nil && _current.start != nil && _current.end != nil && blockNumber.Cmp(_current.start) >= 0 && blockNumber.Cmp(_current.end) <= 0 {
					return d.roundIndex(nodeID, _current)
				}

			case next:
				if _next != nil && _next.start != nil && _next.end != nil && blockNumber.Cmp(_next.start) >= 0 && blockNumber.Cmp(_next.end) <= 0 {
					return d.roundIndex(nodeID, _next)
				}

			default:
				if _former != nil && _former.start != nil && _former.end != nil && blockNumber.Cmp(_former.start) >= 0 && blockNumber.Cmp(_former.end) <= 0 {
					return d.roundIndex(nodeID, _former)
				} else if _current != nil && _current.start != nil && _current.end != nil && blockNumber.Cmp(_current.start) >= 0 && blockNumber.Cmp(_current.end) <= 0 {
					return d.roundIndex(nodeID, _current)
				} else if _next != nil && _next.start != nil && _next.end != nil && blockNumber.Cmp(_next.start) >= 0 && blockNumber.Cmp(_next.end) <= 0 {
					return d.roundIndex(nodeID, _next)
				}
		}
	}
	return -1
}

func (d *ppos) roundIndex(nodeID discover.NodeID, round *pposRound) int64 {
	for idx, nid := range round.nodeIds {
		if nid == nodeID {
			return int64(idx)
		}
	}
	return -1
}

func (d *ppos) NodeIndexInFuture(nodeID discover.NodeID) int64 {
	return -1
}

func (d *ppos) getFormerNodes (parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) []*discover.Node {
	d.lock.RLock()
	defer d.lock.RUnlock()

	formerRound := d.nodeRound.getFormerRound(parentNumber, parentHash)
	if formerRound != nil && len(formerRound.nodes) > 0 && blockNumber.Cmp(formerRound.start) >= 0 && blockNumber.Cmp(formerRound.end) <= 0{
		return formerRound.nodes
	}
	return nil
}

func (d *ppos) getCurrentNodes (parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) []*discover.Node {
	d.lock.RLock()
	defer d.lock.RUnlock()

	currentRound := d.nodeRound.getCurrentRound(parentNumber, parentHash)
	if currentRound != nil && currentRound.start != nil && currentRound.end != nil && len(currentRound.nodes) > 0 && blockNumber.Cmp(currentRound.start) >= 0 && blockNumber.Cmp(currentRound.end) <= 0{
		return currentRound.nodes
	}
	return nil
}


func (d *ppos) consensusNodes(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) []discover.NodeID {
	d.lock.RLock()
	defer d.lock.RUnlock()

	log.Warn("call consensusNodes", "parentNumber", parentNumber.Uint64(), "parentHash", parentHash, "blockNumber", blockNumber.Uint64())
	nodeCache := d.nodeRound.getNodeCache(parentNumber, parentHash)
	d.printMapInfo("consensusNodes nodeCache", parentNumber.Uint64(), parentHash)
	if nodeCache != nil {
		if nodeCache.former != nil && nodeCache.former.start != nil && nodeCache.former.end != nil && blockNumber.Cmp(nodeCache.former.start) >= 0 && blockNumber.Cmp(nodeCache.former.end) <= 0 {
			return nodeCache.former.nodeIds
		} else if nodeCache.current != nil && nodeCache.current.start != nil && nodeCache.current.end != nil && blockNumber.Cmp(nodeCache.current.start) >= 0 && blockNumber.Cmp(nodeCache.current.end) <= 0 {
			return nodeCache.current.nodeIds
		} else if nodeCache.next != nil && nodeCache.next.start != nil && nodeCache.next.end != nil && blockNumber.Cmp(nodeCache.next.start) >= 0 && blockNumber.Cmp(nodeCache.next.end) <= 0 {
			return nodeCache.next.nodeIds
		}
	}
	return nil
}

func (d *ppos) LastCycleBlockNum() uint64 {
	// Get the block height at the end of the final round of consensus
	return d.lastCycleBlockNum
}

func (d *ppos) SetLastCycleBlockNum(blockNumber uint64) {
	// Set the block height at the end of the last round of consensus
	d.lastCycleBlockNum = blockNumber
}


func (d *ppos) StartTimeOfEpoch() int64 {
	return d.startTimeOfEpoch
}

func (d *ppos) SetStartTimeOfEpoch(startTimeOfEpoch int64) {
	d.startTimeOfEpoch = startTimeOfEpoch
}

/** ppos was added func */
/** Method provided to the cbft module call */
// Announce witness
func (d *ppos) Election(state *state.StateDB, parentHash common.Hash, currBlocknumber *big.Int) ([]*discover.Node, error) {
	// TODO
	if nextNodes, err := d.candidatePool.Election(state, parentHash, currBlocknumber); nil != err {
		log.Error("ppos election next witness", " err: ", err)
		panic("Election error " + err.Error())
	} else {
		log.Info("Election finish，view stateDB content again ...")
		d.candidatePool.GetAllWitness(state)

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
	d.candidatePool.GetAllWitness(state)

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

// Getting can by witnesses
// flag:
// -1: 		previous round
// 0:		current round
// 1: 		next round
func (d *ppos) GetWitnessCandidate (state vm.StateDB, nodeId discover.NodeID, flag int) (*types.Candidate, error) {
	return d.candidatePool.GetWitnessCandidate(state, nodeId, flag)
}

// setting candidate pool of ppos module
func (d *ppos) setCandidatePool(blockChain *core.BlockChain, initialNodes []discover.Node) {
	log.Info("---start node，to update nodeRound---")
	genesis := blockChain.Genesis()
	// init roundCache by config
	d.buildGenesisRound(genesis.NumberU64(), genesis.Hash(), initialNodes)
	d.printMapInfo("启动时读取创世块配置", genesis.NumberU64(), genesis.Hash())
	// When the highest block in the chain is not a genesis block, Need to load witness nodeIdList from the stateDB.
	if genesis.NumberU64() != blockChain.CurrentBlock().NumberU64() {

		currentBlock := blockChain.CurrentBlock()
		var currBlockNumber uint64
		var currBlockHash common.Hash
		var currentBigInt *big.Int

		currBlockNumber = blockChain.CurrentBlock().NumberU64()
		currentBigInt = blockChain.CurrentBlock().Number()
		currBlockHash = blockChain.CurrentBlock().Hash()



		count := 0
		blockArr := make([]*types.Block, 0)
		for {
			if currBlockNumber == genesis.NumberU64() || count == BaseIrrCount {
				break
			}

			parentNum := currBlockNumber - 1
			parentBigInt := new(big.Int).Sub(currentBigInt, big.NewInt(1))
			parentHash := currentBlock.ParentHash()
			blockArr = append(blockArr, currentBlock)

			currBlockNumber = parentNum
			currentBigInt = parentBigInt
			currBlockHash = parentHash
			currentBlock = blockChain.GetBlock(currBlockHash, currBlockNumber)
			count ++

		}

		for i := len(blockArr) - 1; 0 <= i; i-- {
			currentBlock := blockArr[i]
			currentNum := currentBlock.NumberU64()
			currentBigInt := currentBlock.Number()
			currentHash := currentBlock.Hash()

			parentNum := currentNum - 1
			parentBigInt := new(big.Int).Sub(currentBigInt, big.NewInt(1))
			parentHash := currentBlock.ParentHash()

			// Special processing of the last block of the array, that is,
			// the highest block pushes the BaseIrrCount block forward
			if i == len(blockArr) - 1 && currentNum > 1  {

				var parent, current *state.StateDB

				// parentStateDB by block
				parentStateRoot := blockChain.GetBlock(parentHash, parentNum).Root()
				log.Debug("【Reload the oldest block at startup】 ", "parentNum", parentNum, "parentHash", parentHash, "parentStateRoot", parentStateRoot.String())
				if parentState, err := blockChain.StateAt(parentStateRoot, parentBigInt, parentHash); nil != err {
					log.Error("Failed to load parentStateDB by block", "currtenNum", currentNum, "Hash", currentHash.String(), "parentNum", parentNum, "Hash", parentHash.String(), "err", err)
					//panic("Failed to load parentStateDB by block parentNum" + fmt.Sprint(parentNum) + ", Hash" + parentHash.String() + "err" + err.Error())
				}else {
					parent = parentState
				}

				// currentStateDB by block
				stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
				log.Debug("【Reload the oldest block at startup】", "currentNum", currentNum, "currentHash", currentHash, "stateRoot", stateRoot.String())
				if currntState, err := blockChain.StateAt(stateRoot, currentBigInt, currentHash); nil != err {
					log.Error("Failed to load currentStateDB by block", "currtenNum", currentNum, "Hash", currentHash.String(), "err", err)
					//panic("Failed to load currentStateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
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
			log.Debug("【Reload the front normal fast at startup】", "currentNum", currentNum, "currentHash", currentHash, "stateRoot", stateRoot.String())
			if currntState, err := blockChain.StateAt(stateRoot, currentBigInt, currentHash); nil != err {
				log.Error("Failed to load stateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				//panic("Failed to load stateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}else {
				if err := d.setGeneralNodeCache(currntState, parentNum, currentNum, parentHash, currentHash); nil != err {
					log.Error("Failed to setGeneralNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
					panic("Failed to setGeneralNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
				}
			}
		}
	}
}


func (d *ppos)buildGenesisRound(blockNumber uint64, blockHash common.Hash, initialNodes []discover.Node) {
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
		start: big.NewInt(1),
		end:   big.NewInt(BaseSwitchWitness),
	}
	currentRound.nodes = make([]*discover.Node, len(initNodeArr))
	copy(currentRound.nodes, initNodeArr)


	log.Debug("【Initialize ppos according to the configuration file】:", "blockNumber", blockNumber, "blockHash", blockHash.String(), "start", currentRound.start, "end", currentRound.end)

	node := &nodeCache{
		former: 	formerRound,
		current: 	currentRound,
	}
	res := make(roundCache, 0)
	hashRound := make(map[common.Hash]*nodeCache, 0)
	hashRound[blockHash] = node
	res[blockNumber] = hashRound
	/* set nodeRound ... */
	d.nodeRound = res
}

func (d *ppos)printMapInfo(title string, blockNumber uint64, blockHash common.Hash){
	res := d.nodeRound[blockNumber]

	log.Debug(title + ":Traversing out the RoundNodes，num: " + fmt.Sprint(blockNumber) + ", hash: " + blockHash.String())
	if round, ok  := res[blockHash]; ok {
		if nil != round.former{
			pposm.PrintObject(title + ":Traversing out of the round，num: " + fmt.Sprint(blockNumber) + ", hash: " + blockHash.String() + ", previous round: start:" + round.former.start.String() + ", end:" + round.former.end.String() + ", nodes: ", round.former.nodes)
		}
		if nil != round.current {
			pposm.PrintObject(title + ":Traversing out of the round，num: " + fmt.Sprint(blockNumber) + ", hash: " + blockHash.String() + ", current round: start:" + round.current.start.String() + ", end:" + round.current.end.String() + ", nodes: ", round.current.nodes)
		}
		if nil != round.next {
			pposm.PrintObject(title + ":Traversing out of the round，num: " + fmt.Sprint(blockNumber) + ", hash: " + blockHash.String() + ", next round: start:" + round.next.start.String() + ", end:" + round.next.end.String() + ", nodes: ", round.next.nodes)
		}
	}else {
		log.Error(title + ":Traversing out of the round is NOT EXIST !!!!!!!!，num: " + fmt.Sprint(blockNumber) + ", hash: " + blockHash.String())
	}
}

/** Method provided to the built-in contract call */
// pledge Candidate
func (d *ppos) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	return d.candidatePool.SetCandidate(state, nodeId, can)
}

// Getting immediate or reserve candidate info by nodeId
func (d *ppos) GetCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error) {
	return d.candidatePool.GetCandidate(state, nodeId)
}

// Getting immediate or reserve candidate info arr by nodeIds
func (d *ppos) GetCandidateArr (state vm.StateDB, nodeIds ... discover.NodeID) (types.CandidateQueue, error) {
	return d.candidatePool.GetCandidateArr(state, nodeIds...)
}

// candidate withdraw from  elected candidates
func (d *ppos) WithdrawCandidate(state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error {
	return d.candidatePool.WithdrawCandidate(state, nodeId, price, blockNumber)
}

// Getting all  elected candidates array
func (d *ppos) GetChosens(state vm.StateDB, flag int) []*types.Candidate {
	return d.candidatePool.GetChosens(state, flag)
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



/** about ticketpool's method */

func (d *ppos) GetPoolNumber (state vm.StateDB) (uint64, error) {
	return d.ticketPool.GetPoolNumber(state)
}

func (d *ppos) VoteTicket (state vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error) {
	return d.ticketPool.VoteTicket(state, owner, voteNumber, deposit, nodeId, blockNumber)
}

func (d *ppos) GetTicket(state vm.StateDB, ticketId common.Hash) (*types.Ticket, error) {
	return d.ticketPool.GetTicket(state, ticketId)
}

func (d *ppos) GetTicketList (state vm.StateDB, ticketIds []common.Hash) ([]*types.Ticket, error) {
	return d.ticketPool.GetTicketList(state, ticketIds)
}

func (d *ppos) GetCandidateTicketIds (state vm.StateDB, nodeId discover.NodeID) ([]common.Hash, error) {
	return d.ticketPool.GetCandidateTicketIds(state, nodeId)
}

func (d *ppos) GetCandidateEpoch (state vm.StateDB, nodeId discover.NodeID) (uint64, error) {
	return d.ticketPool.GetCandidateEpoch(state, nodeId)
}

func (d *ppos) GetTicketPrice (state vm.StateDB) (*big.Int, error) {
	return d.ticketPool.GetTicketPrice(state)
}

func (d *ppos) GetCandidateAttach (state vm.StateDB, nodeId discover.NodeID) (*types.CandidateAttach, error) {
	return d.ticketPool.GetCandidateAttach(state, nodeId)
}

// TODO 每一个块执行交易之后，揭榜或替换之前 都会调用的方法
func (d *ppos) Notify (state vm.StateDB, blockNumber *big.Int) error {
	return d.ticketPool.Notify(state, blockNumber)
}

// TODO 添加一个方法， 每次finalize 之前，调用求Hash 加入 stateDB
func (d *ppos) StoreHash (state *state.StateDB) {
	if err := d.ticketPool.CommitHash(state); nil != err {
		log.Error("Failed to StoreHash", "err", err)
		panic("Failed to StoreHash err" + err.Error())
	}
}

// TODO 添加一个方法，每 seal 完一个块之后，就调用该 Func
func (d *ppos) Submit2Cache (state *state.StateDB, currBlocknumber,  blockInterval *big.Int, currBlockhash common.Hash) {
	d.ticketidsCache.Submit2Cache(currBlocknumber,  blockInterval, currBlockhash, state.TicketCaceheSnapshot())
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
		currentBigInt := currentBlock.Number()
		currentHash := currentBlock.Hash()

		parentNum := currentNum - 1
		parentBigInt := new(big.Int).Sub(currentBigInt, big.NewInt(1))
		parentHash := currentBlock.ParentHash()


		// Special processing of the last block of the array, that is,
		// the highest block pushes the BaseIrrCount block forward
		if i == len(blockArr) - 1 && currentNum > 1  {

			var parent, current *state.StateDB

			// parentStateDB by block
			parentStateRoot := blockChain.GetBlock(parentHash, parentNum).Root()
			if parentState, err := blockChain.StateAt(parentStateRoot, parentBigInt, parentHash); nil != err {
				log.Error("Failed to load parentStateDB by block", "currtenNum", currentNum, "Hash", currentHash.String(), "parentNum", parentNum, "Hash", parentHash.String(), "err", err)
				panic("Failed to load parentStateDB by block parentNum" + fmt.Sprint(parentNum) + ", Hash" + parentHash.String() + "err" + err.Error())
			}else {
				parent = parentState
			}

			// currentStateDB by block
			stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
			if currntState, err := blockChain.StateAt(stateRoot, currentBigInt, currentHash); nil != err {
				log.Error("Failed to load currentStateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to load currentStateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}else {
				current = currntState
			}

			if err := d.setEarliestIrrNodeCache(parent, current, genesis.NumberU64(), currentNum, genesis.Hash(), currentHash); nil != err {
				log.Error("Failed to setEarliestIrrNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to setEarliestIrrNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}
			d.printMapInfo("【Reload the oldest block when forked】", currentNum, currentHash)
			continue
		}

		// stateDB by block
		stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
		if currntState, err := blockChain.StateAt(stateRoot, currentBigInt, currentHash); nil != err {
			log.Error("Failed to load stateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
			panic("Failed to load stateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
		}else {
			if err := d.setGeneralNodeCache(currntState, parentNum, currentNum, parentHash, currentHash); nil != err {
				log.Error("Failed to setGeneralNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to setGeneralNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}
		}
		d.printMapInfo("【Reload the previous normal block when forking】", currentNum, currentHash)
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

	log.Debug("【Setting current  block Node Cache】", "parentNumber", parentNumber, "ParentHash", parentHash.String(), "currentNumber:", currentNumber, "hash", currentHash.String(), "round:", round)

	preNodes, curNodes, nextNodes, err := d.candidatePool.GetAllWitness(state)

	if nil != err {
		log.Error("Failed to setting nodeCache on setGeneralNodeCache", "err", err)
		return err
	}


	var start, end *big.Int

	// Determine if it is the last block of the current round.
	// If it is, start is the start of the next round,
	// and end is the end of the next round.
	if cmpSwitch(round, currentNumber) == 0 {
		start = big.NewInt(int64(BaseSwitchWitness*round) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(BaseSwitchWitness-1)))
	}else {
		start = big.NewInt(int64(BaseSwitchWitness*(round-1)) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(BaseSwitchWitness-1)))
	}

	// former
	formerRound := &pposRound{}
	// former start, end
	if round != FirstRound {
		formerRound.start = new(big.Int).Sub(start, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
		formerRound.end = new(big.Int).Sub(end, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
	}else {
		formerRound.start = big.NewInt(0)
		formerRound.end = big.NewInt(0)
	}

	log.Debug("【Setting current  block Node Cache】 Previous round ", "start",formerRound.start, "end", formerRound.end)

	if len(preNodes) != 0 {
		formerRound.nodeIds = convertNodeID(preNodes)
		formerRound.nodes = make([]*discover.Node, len(preNodes))
		copy(formerRound.nodes, preNodes)
	}else { // Reference parent
		// if last block of round
		if cmpSwitch(round, currentNumber) == 0 {
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

	log.Debug("【Setting current  block Node Cache】 Current round ", "start", currentRound.start, "end",currentRound.end)

	if len(curNodes) != 0 {
		currentRound.nodeIds = convertNodeID(curNodes)
		currentRound.nodes = make([]*discover.Node, len(curNodes))
		copy(currentRound.nodes, curNodes)
	}else { // Reference parent
		// if last block of round
		if cmpSwitch(round, currentNumber) == 0 {
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

	log.Debug("【Setting current  block Node Cache】 Next round ", "start", nextRound.start, "end",nextRound.end)

	if len(nextNodes) != 0 {
		nextRound.nodeIds = convertNodeID(nextNodes)
		nextRound.nodes = make([]*discover.Node, len(nextNodes))
		copy(nextRound.nodes, nextNodes)
	}else { // Reference parent

		if cmpElection(round, currentNumber) == 0  { // election index == cur index
			parentCurRound := d.nodeRound.getCurrentRound(parentNumBigInt, parentHash)
			if nil != parentCurRound {
				nextRound.nodeIds = make([]discover.NodeID, len(parentCurRound.nodeIds))
				copy(nextRound.nodeIds, parentCurRound.nodeIds)
				nextRound.nodes = make([]*discover.Node, len(parentCurRound.nodes))
				copy(nextRound.nodes, parentCurRound.nodes)
			}
		}else if cmpElection(round, currentNumber) > 0  &&  cmpSwitch(round, currentNumber) < 0 {  // election index < cur index < switch index
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

	cache := &nodeCache{
		former: 	formerRound,
		current: 	currentRound,
		next: 		nextRound,
	}
	d.nodeRound.setNodeCache(big.NewInt(int64(currentNumber)), currentHash, cache)

	log.Debug("【When setting the information of the current block】", "currentBlockNum", currentNumber, "parentNum", parentNumber, "currentHash", currentHash.String(), "parentHash", parentHash.String())
	d.printMapInfo("【When setting the information of the current block】", currentNumber, currentHash)

	return nil
}

func (d *ppos) setEarliestIrrNodeCache (parentState, currentState *state.StateDB, genesisNumber, currentNumber uint64, genesisHash, currentHash common.Hash) error {
	genesisNumBigInt := big.NewInt(int64(genesisNumber))
	// current round
	round := calcurround(currentNumber)
	log.Debug("【Set the farthest allowed cache reserved block】", "currentNumber:", currentNumber, "round:", round)

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

	genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash)

	var start, end *big.Int

	// Determine if it is the last block of the current round.
	// If it is, start is the start of the next round,
	// and end is the end of the next round.
	if cmpSwitch(round, currentNumber) == 0 {
		start = big.NewInt(int64(BaseSwitchWitness*round) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(BaseSwitchWitness-1)))
	}else {
		start = big.NewInt(int64(BaseSwitchWitness*(round-1)) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(BaseSwitchWitness-1)))
	}

	// former
	formerRound := &pposRound{}
	// former start, end
	if round != FirstRound {
		formerRound.start = new(big.Int).Sub(start, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
		formerRound.end = new(big.Int).Sub(end, new(big.Int).SetUint64(uint64(BaseSwitchWitness)))
	}else {
		formerRound.start = big.NewInt(0)
		formerRound.end = big.NewInt(0)
	}
	log.Debug("【Set the farthest allowed cache reserved block】: Previous round ", "start",formerRound.start, "end", formerRound.end)
	if len(curr_preNodes) != 0 {
		formerRound.nodeIds = convertNodeID(curr_preNodes)
		formerRound.nodes = make([]*discover.Node, len(curr_preNodes))
		copy(formerRound.nodes, curr_preNodes)
	}else {
		// Reference parent
		// if last block of round
		if cmpSwitch(round, currentNumber) == 0 {
			// First take the stateDB from the previous block,
			// the stateDB of the previous block is not,
			// just take the nodeCache corresponding to the creation block.
			if len(parent_curNodes) != 0 {
				//formerRound.nodeIds = make([]discover.NodeID, len(parent_curNodes))
				formerRound.nodeIds = convertNodeID(parent_curNodes)
				formerRound.nodes = make([]*discover.Node, len(parent_curNodes))
				copy(formerRound.nodes, parent_curNodes)
			}else {
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
				if /*genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash);*/ nil != genesisCurRound {
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
	log.Debug("【Set the farthest allowed cache reserved block】: Current round", "start", currentRound.start, "end",currentRound.end)
	if len(curr_curNodes) != 0 {
		currentRound.nodeIds = convertNodeID(curr_curNodes)
		currentRound.nodes = make([]*discover.Node, len(curr_curNodes))
		copy(currentRound.nodes, curr_curNodes)
	}else { // Reference parent
		// if last block of round
		if cmpSwitch(round, currentNumber) == 0 {
			if len(parent_nextNodes) != 0  {
				currentRound.nodeIds = convertNodeID(parent_nextNodes)
				currentRound.nodes = make([]*discover.Node, len(parent_nextNodes))
				copy(currentRound.nodes, parent_nextNodes)
			}else {
				if /*genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash); */ nil != genesisCurRound {
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
				if /*genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash);*/ nil != genesisCurRound {
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
	log.Debug("【Set the farthest allowed cache reserved block】: Next round", "start", nextRound.start, "end",nextRound.end)
	if len(curr_nextNodes) != 0 {
		nextRound.nodeIds = convertNodeID(curr_nextNodes)
		nextRound.nodes = make([]*discover.Node, len(curr_nextNodes))
		copy(nextRound.nodes, curr_nextNodes)
	}else { // Reference parent
		// election index == cur index || election index < cur index < switch index
		if cmpElection(round, currentNumber) == 0 || (cmpElection(round, currentNumber) > 0 && cmpSwitch(round, currentNumber) < 0)  {

			if /*genesisCurRound := d.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash); */nil != genesisCurRound {
				nextRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
				copy(nextRound.nodeIds, genesisCurRound.nodeIds)
				nextRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
				copy(nextRound.nodes, genesisCurRound.nodes)
			}
		}else { // parent switch index <= cur index < election index  || switch index <= cur index < next election index
			nextRound.nodeIds = make([]discover.NodeID, 0)
			nextRound.nodes = make([]*discover.Node, 0)
		}
	}

	cache := &nodeCache{
		former: 	formerRound,
		current: 	currentRound,
		next: 		nextRound,
	}
	d.nodeRound.setNodeCache(big.NewInt(int64(currentNumber)), currentHash, cache)
	log.Debug("【Set the farthest allowed to cache the information of the reserved block】", "currentBlockNum", currentNumber, "currentHash", currentHash.String())
	d.printMapInfo("【Set the farthest allowed to cache the information of the reserved block】", currentNumber, currentHash)
	return nil
}


func (d *ppos) cleanNodeRound () {
	d.lock.Lock()
	d.nodeRound =  make(roundCache, 0)
	d.lock.Unlock()
}

// election index == cur       0
// cur < election index       -1
// election index < cur        1
// param invalid              -2
func cmpElection (round, currentNumber uint64) int {
	// last num of round
	last := int(round * BaseSwitchWitness)
	ele_sub := int(BaseSwitchWitness - BaseElection)
	curr_sub := last - int(currentNumber)
	sub := ele_sub - curr_sub
	//fmt.Println("sss ", sub)
	if curr_sub < int(0)  {
		return -2
	}else if sub > int(0) {
		return 1
	}else if sub == int(0) {
		return 0
	}else {
		return -1
	}
}

// switch index == cur       0
// cur < switch index       -1
// switch index < cur        1
// param invalid            -2
func cmpSwitch (round, currentNum uint64) int {
	last := round * BaseSwitchWitness
	if last < currentNum {
		return 1
	}else if last == currentNum {
		return 0
	}else {
		return -1
	}
}

func (d *ppos) setTicketPoolCache () {
	d.ticketidsCache = ticketcache.GetTicketidsCachePtr()
}
