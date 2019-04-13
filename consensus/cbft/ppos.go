package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/ppos"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"fmt"
	"math/big"
	"sync"
	"github.com/PlatONnetwork/PlatON-Go/core/ppos_storage"
)

type ppos struct {

	nodeRound 		  roundCache
	//chain             *core.BlockChain
	lastCycleBlockNum uint64
	// A round of consensus start time is usually the block time of the last block at the end of the last round of consensus;
	// if it is the first round, it starts from 1970.1.1.0.0.0.0. Unit: second
	startTimeOfEpoch  int64
	config            *params.PposConfig

	// added by candidateContext module
	lock 					sync.RWMutex
	// the candidateContext pool object pointer
	candidateContext 	*pposm.CandidatePoolContext
	// the ticket pool object pointer
	ticketContext				*pposm.TicketPoolContext

	pposTemp 			*ppos_storage.PPOS_TEMP
}



func newPpos(config *params.CbftConfig) *ppos {
	return &ppos{
		lastCycleBlockNum: 	0,
		config:            	config.PposConfig,
		candidateContext:   pposm.NewCandidatePoolContext(config.PposConfig),
		ticketContext: 		pposm.NewTicketPoolContext(config.PposConfig),
	}
}


func (p *ppos) BlockProducerIndex(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int, nodeID discover.NodeID, round int32) (int64, []discover.NodeID) {
	p.lock.RLock()
	defer p.lock.RUnlock()

	log.Info("BlockProducerIndex", "parentNumber", parentNumber, "parentHash", parentHash.String(), "blockNumber", blockNumber.String(), "nodeID", nodeID.String(), "round", round)

	nodeCache := p.nodeRound.getNodeCache(parentNumber, parentHash)
	p.printMapInfo("BlockProducerIndex", parentNumber.Uint64(), parentHash)
	if nodeCache != nil {
		_former := nodeCache.former
		_current := nodeCache.current
		_next := nodeCache.next

		switch round {
			case former:
				if _former != nil && _former.start != nil && _former.end != nil && blockNumber.Cmp(_former.start) >= 0 && blockNumber.Cmp(_former.end) <= 0 {
					return p.roundIndex(nodeID, _former), _former.nodeIds
				}

			case current:
				if _current != nil && _current.start != nil && _current.end != nil && blockNumber.Cmp(_current.start) >= 0 && blockNumber.Cmp(_current.end) <= 0 {
					return p.roundIndex(nodeID, _current), _current.nodeIds
				}

			case next:
				if _next != nil && _next.start != nil && _next.end != nil && blockNumber.Cmp(_next.start) >= 0 && blockNumber.Cmp(_next.end) <= 0 {
					return p.roundIndex(nodeID, _next), _next.nodeIds
				}

			default:
				if _former != nil && _former.start != nil && _former.end != nil && blockNumber.Cmp(_former.start) >= 0 && blockNumber.Cmp(_former.end) <= 0 {
					return p.roundIndex(nodeID, _former), _former.nodeIds
				} else if _current != nil && _current.start != nil && _current.end != nil && blockNumber.Cmp(_current.start) >= 0 && blockNumber.Cmp(_current.end) <= 0 {
					return p.roundIndex(nodeID, _current), _current.nodeIds
				} else if _next != nil && _next.start != nil && _next.end != nil && blockNumber.Cmp(_next.start) >= 0 && blockNumber.Cmp(_next.end) <= 0 {
					return p.roundIndex(nodeID, _next), _next.nodeIds
				}
		}
	}
	return -1, nil
}

func (p *ppos) roundIndex(nodeID discover.NodeID, round *pposRound) int64 {
	for idx, nid := range round.nodeIds {
		if nid == nodeID {
			return int64(idx)
		}
	}
	return -1
}

func (p *ppos) NodeIndexInFuture(nodeID discover.NodeID) int64 {
	return -1
}

func (p *ppos) getFormerNodes (parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) []*discover.Node {
	p.lock.RLock()
	defer p.lock.RUnlock()

	formerRound := p.nodeRound.getFormerRound(parentNumber, parentHash)
	if formerRound != nil && len(formerRound.nodes) > 0 && blockNumber.Cmp(formerRound.start) >= 0 && blockNumber.Cmp(formerRound.end) <= 0{
		return formerRound.nodes
	}
	return nil
}

func (p *ppos) getCurrentNodes (parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) []*discover.Node {
	p.lock.RLock()
	defer p.lock.RUnlock()

	currentRound := p.nodeRound.getCurrentRound(parentNumber, parentHash)
	if currentRound != nil && currentRound.start != nil && currentRound.end != nil && len(currentRound.nodes) > 0 && blockNumber.Cmp(currentRound.start) >= 0 && blockNumber.Cmp(currentRound.end) <= 0{
		return currentRound.nodes
	}
	return nil
}


func (p *ppos) consensusNodes(parentNumber *big.Int, parentHash common.Hash, blockNumber *big.Int) []discover.NodeID {
	p.lock.RLock()
	defer p.lock.RUnlock()

	log.Debug("call consensusNodes", "parentNumber", parentNumber.Uint64(), "parentHash", parentHash, "blockNumber", blockNumber.Uint64())
	nodeCache := p.nodeRound.getNodeCache(parentNumber, parentHash)
	p.printMapInfo("consensusNodes nodeCache", parentNumber.Uint64(), parentHash)
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

func (p *ppos) LastCycleBlockNum() uint64 {
	// Get the block height at the end of the final round of consensus
	return p.lastCycleBlockNum
}

func (p *ppos) SetLastCycleBlockNum(blockNumber uint64) {
	// Set the block height at the end of the last round of consensus
	p.lastCycleBlockNum = blockNumber
}


func (p *ppos) StartTimeOfEpoch() int64 {
	return p.startTimeOfEpoch
}

func (p *ppos) SetStartTimeOfEpoch(startTimeOfEpoch int64) {
	p.startTimeOfEpoch = startTimeOfEpoch
}

/** ppos was added func */
/** Method provided to the cbft module call */
// Announce witness
func (p *ppos) Election(state *state.StateDB, parentHash common.Hash, currBlocknumber *big.Int) ([]*discover.Node, error) {
	if nextNodes, err := p.candidateContext.Election(state, parentHash, currBlocknumber); nil != err {
		log.Error("PPOS Election next witness", " err: ", err)
		/*panic("Election error " + err.Error())*/
		return nil, err
	} else {
		//d.candidateContext.ForEachStorage(state, "PPOS Election finish，view stateDB content again ...")
		return nextNodes, nil
	}
}

// switch next witnesses to current witnesses
func (p *ppos) Switch(state *state.StateDB, blockNumber *big.Int) bool {
	log.Info("Switch begin...")
	if !p.candidateContext.Switch(state, blockNumber) {
		return false
	}
	log.Info("Switch success...")
	p.candidateContext.GetAllWitness(state, blockNumber)

	return true
}

// Getting nodes of witnesses
// flag：-1: the previous round of witnesses  0: the current round of witnesses   1: the next round of witnesses
func (p *ppos) GetWitness(state *state.StateDB, flag int, blockNumber *big.Int) ([]*discover.Node, error) {
	return p.candidateContext.GetWitness(state, flag, blockNumber)
}

func (p *ppos) GetAllWitness(state *state.StateDB, blockNumber *big.Int) ([]*discover.Node, []*discover.Node, []*discover.Node, error) {
	return p.candidateContext.GetAllWitness(state, blockNumber)
}

// Getting can by witnesses
// flag:
// -1: 		previous round
// 0:		current round
// 1: 		next round
func (p *ppos) GetWitnessCandidate (state vm.StateDB, nodeId discover.NodeID, flag int, blockNumber *big.Int) *types.Candidate {
	return p.candidateContext.GetWitnessCandidate(state, nodeId, flag, blockNumber)
}

// setting candidate pool of ppos module
func (p *ppos) SetCandidateContextOption(blockChain *core.BlockChain, initialNodes []discover.Node) {
	log.Info("Start node, build the nodeRound ...")

	genesis := blockChain.Genesis()

	// init roundCache by config
	p.buildGenesisRound(genesis.NumberU64(), genesis.Hash(), initialNodes)
	p.printMapInfo("Read Genesis block configuration at startup:", genesis.NumberU64(), genesis.Hash())

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
			if currBlockNumber == genesis.NumberU64() || count == 2*common.BaseIrrCount {
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
			//parentBigInt := new(big.Int).Sub(currentBigInt, big.NewInt(1))
			parentHash := currentBlock.ParentHash()

			// Special processing of the last block of the array, that is,
			// the highest block pushes the BaseIrrCount block forward
			if i == len(blockArr) - 1 && currentNum > 1  {

				var /*parent,*/ current *state.StateDB

				// parentStateDB by block
				/*parentStateRoot := blockChain.GetBlock(parentHash, parentNum).Root()
				log.Debug("Reload the oldest block at startup", "parentNum", parentNum, "parentHash", parentHash, "parentStateRoot", parentStateRoot.String())
				if parentState, err := blockChain.StateAt(parentStateRoot, parentBigInt, parentHash); nil != err {
					log.Error("Failed to load parentStateDB by block", "currtenNum", currentNum, "Hash", currentHash.String(), "parentNum", parentNum, "Hash", parentHash.String(), "err", err)
					panic("Failed to load parentStateDB by block parentNum" + fmt.Sprint(parentNum) + ", Hash" + parentHash.String() + "err" + err.Error())
				}else {
					parent = parentState
				}*/

				// currentStateDB by block
				stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
				log.Debug("Reload the oldest block at startup", "currentNum", currentNum, "currentHash", currentHash, "stateRoot", stateRoot.String())
				if currntState, err := blockChain.StateAt(stateRoot, currentBigInt, currentHash); nil != err {
					log.Error("Failed to load currentStateDB by block", "currtenNum", currentNum, "Hash", currentHash.String(), "err", err)
					panic("Failed to load currentStateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
				}else {
					current = currntState
				}

				if err := p.setEarliestIrrNodeCache(/*parent,*/ current, genesis.NumberU64(), currentNum, genesis.Hash(), currentHash); nil != err {
					log.Error("Failed to setEarliestIrrNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
					panic("Failed to setEarliestIrrNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
				}
				continue
			}

			// stateDB by block
			stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
			log.Debug("Reload the front normal fast at startup", "currentNum", currentNum, "currentHash", currentHash, "stateRoot", stateRoot.String())
			if currntState, err := blockChain.StateAt(stateRoot, currentBigInt, currentHash); nil != err {
				log.Error("Failed to load stateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to load stateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}else {
				if err := p.setGeneralNodeCache(currntState, genesis.NumberU64(), parentNum, currentNum, genesis.Hash(), parentHash, currentHash); nil != err {
					log.Error("Failed to setGeneralNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
					panic("Failed to setGeneralNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
				}
			}
		}
	}
}


func (p *ppos)buildGenesisRound(blockNumber uint64, blockHash common.Hash, initialNodes []discover.Node) {
	initNodeArr := make([]*discover.Node, 0, len(initialNodes))
	initialNodesIDs := make([]discover.NodeID, 0, len(initialNodes))
	for _, n := range initialNodes {
		node := n
		initialNodesIDs = append(initialNodesIDs, node.ID)
		initNodeArr = append(initNodeArr, &node)
	}

	// previous round
	formerRound := &pposRound{
		nodeIds: make([]discover.NodeID, 0),
		nodes: 	make([]*discover.Node, 0),
		start: big.NewInt(0),
		end:   big.NewInt(0),
	}

	// current round
	currentRound := &pposRound{
		nodeIds: initialNodesIDs,
		start: big.NewInt(1),
		end:   big.NewInt(common.BaseSwitchWitness),
	}

	currentRound.nodes = make([]*discover.Node, len(initNodeArr))
	copy(currentRound.nodes, initNodeArr)


	log.Debug("Initialize ppos according to the configuration file:", "blockNumber", blockNumber, "blockHash", blockHash.String(), "start", currentRound.start, "end", currentRound.end)

	node := &nodeCache{
		former: 	formerRound,
		current: 	currentRound,
	}
	res := make(roundCache, 0)
	hashRound := make(map[common.Hash]*nodeCache, 0)
	hashRound[blockHash] = node
	res[blockNumber] = hashRound
	/* set nodeRound ... */
	p.nodeRound = res
}

func (p *ppos)printMapInfo(title string, blockNumber uint64, blockHash common.Hash){
	res := p.nodeRound[blockNumber]

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
func (p *ppos) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	return p.candidateContext.SetCandidate(state, nodeId, can)
}

// Getting immediate or reserve candidate info by nodeId
func (p *ppos) GetCandidate(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) *types.Candidate {
	return p.candidateContext.GetCandidate(state, nodeId, blockNumber)
}

// Getting immediate or reserve candidate info arr by nodeIds
func (p *ppos) GetCandidateArr (state vm.StateDB, blockNumber *big.Int, nodeIds ... discover.NodeID) types.CandidateQueue {
	return p.candidateContext.GetCandidateArr(state, blockNumber, nodeIds...)
}

// candidate withdraw from  elected candidates
func (p *ppos) WithdrawCandidate(state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error {
	return p.candidateContext.WithdrawCandidate(state, nodeId, price, blockNumber)
}

// Getting all  elected candidates array
func (p *ppos) GetChosens(state vm.StateDB, flag int, blockNumber *big.Int) types.KindCanQueue {
	return p.candidateContext.GetChosens(state, flag, blockNumber)
}

func (p *ppos) GetCandidatePendArr (state vm.StateDB, flag int, blockNumber *big.Int) types.CandidateQueue {
	return p.candidateContext.GetCandidatePendArr(state, flag,  blockNumber)
}

// Getting all witness array
func (p *ppos) GetChairpersons(state vm.StateDB, blockNumber *big.Int) []*types.Candidate {
	return p.candidateContext.GetChairpersons(state, blockNumber)
}

// Getting all refund array by nodeId
func (p *ppos) GetDefeat(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) types.RefundQueue {
	return p.candidateContext.GetDefeat(state, nodeId, blockNumber)
}

// Checked current candidate was defeat by nodeId
func (p *ppos) IsDefeat(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) bool {
	return p.candidateContext.IsDefeat(state, nodeId, blockNumber)
}

// refund once
func (p *ppos) RefundBalance(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) error {
	return p.candidateContext.RefundBalance(state, nodeId, blockNumber)
}

// Getting owner's address of candidate info by nodeId
func (p *ppos) GetOwner(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) common.Address {
	return p.candidateContext.GetOwner(state, nodeId, blockNumber)
}

// Getting allow block interval for refunds
func (p *ppos) GetRefundInterval(blockNumber *big.Int) uint32 {
	return p.candidateContext.GetRefundInterval(blockNumber)
}



/** about ticketpool's method */

func (p *ppos) GetPoolNumber (state vm.StateDB) uint32 {
	return p.ticketContext.GetPoolNumber(state)
}

func (p *ppos) VoteTicket (state vm.StateDB, owner common.Address, voteNumber uint32, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) (uint32, error) {
	return p.ticketContext.VoteTicket(state, owner, voteNumber, deposit, nodeId, blockNumber)
}

func (d *ppos) GetTicket(state vm.StateDB, ticketId common.Hash) *types.Ticket {
	return d.ticketContext.GetTicket(state, ticketId)
}

func (p *ppos) GetTicketList (state vm.StateDB, ticketIds []common.Hash) []*types.Ticket {
	return p.ticketContext.GetTicketList(state, ticketIds)
}

func (p *ppos) GetCandidateTicketIds (state vm.StateDB, nodeId discover.NodeID) []common.Hash {
	return p.ticketContext.GetCandidateTicketIds(state, nodeId)
}

func (p *ppos) GetCandidateEpoch (state vm.StateDB, nodeId discover.NodeID) uint64 {
	return p.ticketContext.GetCandidateEpoch(state, nodeId)
}

func (p *ppos) GetTicketPrice (state vm.StateDB) *big.Int {
	return p.ticketContext.GetTicketPrice(state)
}

func (p *ppos) Notify (state vm.StateDB, blockNumber *big.Int) error {
	return p.ticketContext.Notify(state, blockNumber)
}

func (p *ppos) StoreHash (state *state.StateDB, blockNumber *big.Int, blockHash common.Hash) {
	if err := p.ticketContext.StoreHash(state, blockNumber, blockHash); nil != err {
		log.Error("Failed to StoreHash", "err", err)
		panic("Failed to StoreHash err" + err.Error())
	}
}

func (p *ppos) Submit2Cache (state *state.StateDB, currBlocknumber,  blockInterval *big.Int, currBlockhash common.Hash) {
	p.pposTemp.SubmitPposCache2Temp(currBlocknumber,  blockInterval, currBlockhash, state.SnapShotPPOSCache())
}

// cbft consensus fork need to update  nodeRound
func (p *ppos) UpdateNodeList(blockChain *core.BlockChain, blocknumber *big.Int, blockHash common.Hash) {
	log.Info("---cbft consensus fork，update nodeRound---")
	// clean nodeCache
	p.cleanNodeRound()


	var curBlockNumber uint64 = blocknumber.Uint64()
	var curBlockHash common.Hash = blockHash

	currentBlock := blockChain.GetBlock(curBlockHash, curBlockNumber)
	genesis := blockChain.Genesis()
	p.lock.Lock()
	defer p.lock.Unlock()

	count := 0
	blockArr := make([]*types.Block, 0)
	for {
		if curBlockNumber == genesis.NumberU64() || count == 2*common.BaseIrrCount {
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
		//parentBigInt := new(big.Int).Sub(currentBigInt, big.NewInt(1))
		parentHash := currentBlock.ParentHash()


		// Special processing of the last block of the array, that is,
		// the highest block pushes the BaseIrrCount block forward
		if i == len(blockArr) - 1 && currentNum > 1  {

			var/* parent,*/ current *state.StateDB

			/*// parentStateDB by block
			parentStateRoot := blockChain.GetBlock(parentHash, parentNum).Root()
			if parentState, err := blockChain.StateAt(parentStateRoot, parentBigInt, parentHash); nil != err {
				log.Error("Failed to load parentStateDB by block", "currtenNum", currentNum, "Hash", currentHash.String(), "parentNum", parentNum, "Hash", parentHash.String(), "err", err)
				panic("Failed to load parentStateDB by block parentNum" + fmt.Sprint(parentNum) + ", Hash" + parentHash.String() + "err" + err.Error())
			}else {
				parent = parentState
			}*/

			// currentStateDB by block
			stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
			if currntState, err := blockChain.StateAt(stateRoot, currentBigInt, currentHash); nil != err {
				log.Error("Failed to load currentStateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to load currentStateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}else {
				current = currntState
			}

			if err := p.setEarliestIrrNodeCache(/*parent,*/ current, genesis.NumberU64(), currentNum, genesis.Hash(), currentHash); nil != err {
				log.Error("Failed to setEarliestIrrNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to setEarliestIrrNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}
			p.printMapInfo("Reload the oldest block when forked", currentNum, currentHash)
			continue
		}

		// stateDB by block
		stateRoot := blockChain.GetBlock(currentHash, currentNum).Root()
		if currntState, err := blockChain.StateAt(stateRoot, currentBigInt, currentHash); nil != err {
			log.Error("Failed to load stateDB by block", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
			panic("Failed to load stateDB by block currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
		}else {
			if err := p.setGeneralNodeCache(currntState, genesis.NumberU64(), parentNum, currentNum, genesis.Hash(), parentHash, currentHash); nil != err {
				log.Error("Failed to setGeneralNodeCache", "currentNum", currentNum, "Hash", currentHash.String(), "err", err)
				panic("Failed to setGeneralNodeCache currentNum" + fmt.Sprint(currentNum) + ", Hash" + currentHash.String() + "err" + err.Error())
			}
		}
		p.printMapInfo("Reload the previous normal block when forking", currentNum, currentHash)
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
	div := blocknumber / common.BaseSwitchWitness
	mod := blocknumber % common.BaseSwitchWitness
	if (div == 0 && mod == 0) || (div == 0 && mod > 0 && mod < common.BaseSwitchWitness) { // first round
		round = 1
	} else if div > 0 && mod == 0 {
		round = div
	} else if div > 0 && mod > 0 && mod < common.BaseSwitchWitness {
		round = div + 1
	}
	return round
}


func (p *ppos) GetFormerRound(blockNumber *big.Int, blockHash common.Hash) *pposRound {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.nodeRound.getFormerRound(blockNumber, blockHash)
}

func (p *ppos) GetCurrentRound (blockNumber *big.Int, blockHash common.Hash) *pposRound {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.nodeRound.getCurrentRound(blockNumber, blockHash)
}

func (p *ppos)  GetNextRound (blockNumber *big.Int, blockHash common.Hash) *pposRound {
	p.lock.RLock()
	defer p.lock.RUnlock()
	return p.nodeRound.getNextRound(blockNumber, blockHash)
}

func (p *ppos) SetNodeCache (state *state.StateDB, genesisNumber, parentNumber, currentNumber *big.Int, genesisHash, parentHash, currentHash common.Hash) error {
	p.lock.Lock()
	defer p.lock.Unlock()
	return p.setGeneralNodeCache(state, genesisNumber.Uint64(), parentNumber.Uint64(), currentNumber.Uint64(), genesisHash, parentHash, currentHash)
}
func (p *ppos) setGeneralNodeCache (state *state.StateDB, genesisNumber, parentNumber, currentNumber uint64, genesisHash, parentHash, currentHash common.Hash) error {
	parentNumBigInt := big.NewInt(int64(parentNumber))
	// current round
	round := calcurround(currentNumber)

	log.Debug("Setting current  block Node Cache", "parentNumber", parentNumber, "ParentHash", parentHash.String(), "currentNumber:", currentNumber, "hash", currentHash.String(), "round:", round)

	/** ------------------------------ current ppos ------------------------------ **/
	preNodes, curNodes, nextNodes, err := p.candidateContext.GetAllWitness(state, big.NewInt(int64(currentNumber)))

	if nil != err {
		log.Error("Failed to setting nodeCache on setGeneralNodeCache", "err", err)
		return err
	}

	/** ------------------------------    parent    ------------------------------ **/

	parent_Former_Round := p.nodeRound.getFormerRound(parentNumBigInt, parentHash)

	parent_Current_Round := p.nodeRound.getCurrentRound(parentNumBigInt, parentHash)

	parent_Next_Round := p.nodeRound.getNextRound(parentNumBigInt, parentHash)

	/** ------------------------------    genesis    ------------------------------ **/

	genesisNumBigInt := big.NewInt(int64(genesisNumber))
	genesis_Current_Round := p.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash)

	/** ------------------------------       end       ---------------------------- **/

	var start, end *big.Int

	// Determine if it is the last block of the current round.
	// If it is, start is the start of the next round,
	// and end is the end of the next round.
	if cmpSwitch(round, currentNumber) == 0 {
		start = big.NewInt(int64(common.BaseSwitchWitness*round) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(common.BaseSwitchWitness-1)))
	}else {
		start = big.NewInt(int64(common.BaseSwitchWitness*(round-1)) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(common.BaseSwitchWitness-1)))
	}

	// former
	formerRound := &pposRound{}
	// former start, end
	if round != common.FirstRound {
		formerRound.start = new(big.Int).Sub(start, new(big.Int).SetUint64(uint64(common.BaseSwitchWitness)))
		formerRound.end = new(big.Int).Sub(end, new(big.Int).SetUint64(uint64(common.BaseSwitchWitness)))
	}else {
		formerRound.start = big.NewInt(0)
		formerRound.end = big.NewInt(0)
	}

	log.Debug("Setting current  block Node Cache Previous round ", "start",formerRound.start, "end", formerRound.end)

	if len(preNodes) != 0 {
		formerRound.nodeIds = convertNodeID(preNodes)
		formerRound.nodes = make([]*discover.Node, len(preNodes))
		copy(formerRound.nodes, preNodes)
	}else { // Reference parent
		// if last block of round
		if cmpSwitch(round, currentNumber) == 0 {
			//parentCurRound := p.nodeRound.getCurrentRound(parentNumBigInt, parentHash)
			if nil != parent_Current_Round {
				formerRound.nodeIds = make([]discover.NodeID, len(parent_Current_Round.nodeIds))
				copy(formerRound.nodeIds, parent_Current_Round.nodeIds)
				formerRound.nodes = make([]*discover.Node, len(parent_Current_Round.nodes))
				copy(formerRound.nodes, parent_Current_Round.nodes)
			}else {
				formerRound.nodeIds = make([]discover.NodeID, len(genesis_Current_Round.nodeIds))
				copy(formerRound.nodeIds, genesis_Current_Round.nodeIds)
				formerRound.nodes = make([]*discover.Node, len(genesis_Current_Round.nodes))
				copy(formerRound.nodes, genesis_Current_Round.nodes)
			}
		}else { // Is'nt last block of round
			//parentFormerRound := p.nodeRound.getFormerRound(parentNumBigInt, parentHash)
			if nil != parent_Former_Round {
				formerRound.nodeIds = make([]discover.NodeID, len(parent_Former_Round.nodeIds))
				copy(formerRound.nodeIds, parent_Former_Round.nodeIds)
				formerRound.nodes = make([]*discover.Node, len(parent_Former_Round.nodes))
				copy(formerRound.nodes, parent_Former_Round.nodes)
			}else {
				formerRound.nodeIds = make([]discover.NodeID, len(genesis_Current_Round.nodeIds))
				copy(formerRound.nodeIds, genesis_Current_Round.nodeIds)
				formerRound.nodes = make([]*discover.Node, len(genesis_Current_Round.nodes))
				copy(formerRound.nodes, genesis_Current_Round.nodes)
			}
		}
	}

	// current
	currentRound := &pposRound{}
	// current start, end
	currentRound.start = start
	currentRound.end = end

	log.Debug("Setting current  block Node Cache Current round ", "start", currentRound.start, "end",currentRound.end)

	if len(curNodes) != 0 {
		currentRound.nodeIds = convertNodeID(curNodes)
		currentRound.nodes = make([]*discover.Node, len(curNodes))
		copy(currentRound.nodes, curNodes)
	}else { // Reference parent
		// if last block of round
		if cmpSwitch(round, currentNumber) == 0 {
			//parentNextRound := p.nodeRound.getNextRound(parentNumBigInt, parentHash)
			if nil != parent_Next_Round {
				currentRound.nodeIds = make([]discover.NodeID, len(parent_Next_Round.nodeIds))
				copy(currentRound.nodeIds, parent_Next_Round.nodeIds)
				currentRound.nodes = make([]*discover.Node, len(parent_Next_Round.nodes))
				copy(currentRound.nodes, parent_Next_Round.nodes)
			}else {
				currentRound.nodeIds = make([]discover.NodeID, len(genesis_Current_Round.nodeIds))
				copy(currentRound.nodeIds, genesis_Current_Round.nodeIds)
				currentRound.nodes = make([]*discover.Node, len(genesis_Current_Round.nodes))
				copy(currentRound.nodes, genesis_Current_Round.nodes)
			}

		}else { // Is'nt last block of round
			//parentCurRound := p.nodeRound.getCurrentRound(parentNumBigInt, parentHash)
			if nil != parent_Current_Round {
				currentRound.nodeIds = make([]discover.NodeID, len(parent_Current_Round.nodeIds))
				copy(currentRound.nodeIds, parent_Current_Round.nodeIds)
				currentRound.nodes = make([]*discover.Node, len(parent_Current_Round.nodes))
				copy(currentRound.nodes, parent_Current_Round.nodes)
			}else {
				currentRound.nodeIds = make([]discover.NodeID, len(genesis_Current_Round.nodeIds))
				copy(currentRound.nodeIds, genesis_Current_Round.nodeIds)
				currentRound.nodes = make([]*discover.Node, len(genesis_Current_Round.nodes))
				copy(currentRound.nodes, genesis_Current_Round.nodes)
			}
		}
	}


	// next
	nextRound := &pposRound{}
	// next start, end
	nextRound.start = new(big.Int).Add(start, new(big.Int).SetUint64(uint64(common.BaseSwitchWitness)))
	nextRound.end = new(big.Int).Add(end, new(big.Int).SetUint64(uint64(common.BaseSwitchWitness)))

	log.Debug("Setting current  block Node Cache Next round ", "start", nextRound.start, "end",nextRound.end)

	if len(nextNodes) != 0 {
		nextRound.nodeIds = convertNodeID(nextNodes)
		nextRound.nodes = make([]*discover.Node, len(nextNodes))
		copy(nextRound.nodes, nextNodes)
	}else { // Reference parent

		if cmpElection(round, currentNumber) == 0  { // election index == cur index
			//parentCurRound := p.nodeRound.getCurrentRound(parentNumBigInt, parentHash)
			if nil != parent_Current_Round {
				nextRound.nodeIds = make([]discover.NodeID, len(parent_Current_Round.nodeIds))
				copy(nextRound.nodeIds, parent_Current_Round.nodeIds)
				nextRound.nodes = make([]*discover.Node, len(parent_Current_Round.nodes))
				copy(nextRound.nodes, parent_Current_Round.nodes)
			}else {
				nextRound.nodeIds = make([]discover.NodeID, len(genesis_Current_Round.nodeIds))
				copy(nextRound.nodeIds, genesis_Current_Round.nodeIds)
				nextRound.nodes = make([]*discover.Node, len(genesis_Current_Round.nodes))
				copy(nextRound.nodes, genesis_Current_Round.nodes)
			}
		}else if cmpElection(round, currentNumber) > 0  &&  cmpSwitch(round, currentNumber) < 0 {  // election index < cur index < switch index
			//parentNextRound := p.nodeRound.getNextRound(parentNumBigInt, parentHash)
			if nil != parent_Next_Round {
				nextRound.nodeIds = make([]discover.NodeID, len(parent_Next_Round.nodeIds))
				copy(nextRound.nodeIds, parent_Next_Round.nodeIds)
				nextRound.nodes = make([]*discover.Node, len(parent_Next_Round.nodes))
				copy(nextRound.nodes, parent_Next_Round.nodes)
			}else {
				nextRound.nodeIds = make([]discover.NodeID, len(genesis_Current_Round.nodeIds))
				copy(nextRound.nodeIds, genesis_Current_Round.nodeIds)
				nextRound.nodes = make([]*discover.Node, len(genesis_Current_Round.nodes))
				copy(nextRound.nodes, genesis_Current_Round.nodes)
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
	p.nodeRound.setNodeCache(big.NewInt(int64(currentNumber)), currentHash, cache)

	log.Debug("When setting the information of the current block", "currentBlockNum", currentNumber, "parentNum", parentNumber, "currentHash", currentHash.String(), "parentHash", parentHash.String())
	p.printMapInfo("When setting the information of the current block", currentNumber, currentHash)

	return nil
}

func (p *ppos) setEarliestIrrNodeCache (/*parentState, */currentState *state.StateDB, genesisNumber, currentNumber uint64, genesisHash, currentHash common.Hash) error {
	genesisNumBigInt := big.NewInt(int64(genesisNumber))
	// current round
	round := calcurround(currentNumber)
	log.Debug("Set the farthest allowed cache reserved block", "currentNumber:", currentNumber, "round:", round)

	curr_PRE_Nodes, curr_CURR_Nodes, curr_NEXT_Nodes, err := p.candidateContext.GetAllWitness(currentState, big.NewInt(int64(currentNumber)))

	if nil != err {
		log.Error("Failed to setting nodeCache by currentStateDB on setEarliestIrrNodeCache", "err", err)
		return err
	}

	/*parent_preNodes, parent_curNodes, parent_nextNodes, err := p.candidateContext.GetAllWitness(parentState)
	if nil != err {
		log.Error("Failed to setting nodeCache by parentStateDB on setEarliestIrrNodeCache", "err", err)
		return err
	}*/

	genesisCurRound := p.nodeRound.getCurrentRound(genesisNumBigInt, genesisHash)

	var start, end *big.Int

	/**
	Determine if it is the last block of the current round.
	If it is, start is the start of the next round,
	and end is the end of the next round.
	 */
	if cmpSwitch(round, currentNumber) == 0 {
		start = big.NewInt(int64(common.BaseSwitchWitness*round) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(common.BaseSwitchWitness-1)))
	}else {
		start = big.NewInt(int64(common.BaseSwitchWitness*(round-1)) + 1)
		end = new(big.Int).Add(start, big.NewInt(int64(common.BaseSwitchWitness-1)))
	}

	/**
	Sets former info
	 */
	formerRound := &pposRound{}
	// former start, end
	if round != common.FirstRound {
		formerRound.start = new(big.Int).Sub(start, new(big.Int).SetUint64(uint64(common.BaseSwitchWitness)))
		formerRound.end = new(big.Int).Sub(end, new(big.Int).SetUint64(uint64(common.BaseSwitchWitness)))
	}else {
		formerRound.start = big.NewInt(0)
		formerRound.end = big.NewInt(0)
	}

	log.Debug("Set the farthest allowed cache reserved block: Previous round ", "start",formerRound.start, "end", formerRound.end)

	// sets previous
	if len(curr_PRE_Nodes) != 0 {
		formerRound.nodeIds = convertNodeID(curr_PRE_Nodes)
		formerRound.nodes = make([]*discover.Node, len(curr_PRE_Nodes))
		copy(formerRound.nodes, curr_PRE_Nodes)
	}else {
		// Reference parent
		// if last block of roundcurrentState
		if cmpSwitch(round, currentNumber) == 0 {
			// First take the stateDB from the previous block,
			// the stateDB of the previous block is not,
			// just take the nodeCache corresponding to the creation block.
			/*if len(parent_curNodes) != 0 {
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
			}*/

			if nil != genesisCurRound {
				formerRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
				copy(formerRound.nodeIds, genesisCurRound.nodeIds)
				formerRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
				copy(formerRound.nodes, genesisCurRound.nodes)
			}

		}else { // Is'nt last block of round

			/*if len(parent_preNodes) != 0 {
				//formerRound.nodeIds = make([]discover.NodeID, len(parent_preNodes))
				formerRound.nodeIds = convertNodeID(parent_preNodes)
				formerRound.nodes = make([]*discover.Node, len(parent_preNodes))
				copy(formerRound.nodes, parent_preNodes)
			}else {
				if  nil != genesisCurRound {
					formerRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
					copy(formerRound.nodeIds, genesisCurRound.nodeIds)
					formerRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
					copy(formerRound.nodes, genesisCurRound.nodes)
				}
			}*/

			if  nil != genesisCurRound {
				formerRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
				copy(formerRound.nodeIds, genesisCurRound.nodeIds)
				formerRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
				copy(formerRound.nodes, genesisCurRound.nodes)
			}

		}
	}

	/**
	Sets current info
	  */
	currentRound := &pposRound{}
	// current start, end
	currentRound.start = start
	currentRound.end = end

	log.Debug("Set the farthest allowed cache reserved block: Current round", "start", currentRound.start, "end",currentRound.end)

	// sets current
	if len(curr_CURR_Nodes) != 0 {
		currentRound.nodeIds = convertNodeID(curr_CURR_Nodes)
		currentRound.nodes = make([]*discover.Node, len(curr_CURR_Nodes))
		copy(currentRound.nodes, curr_CURR_Nodes)
	}else { // Reference parent
		// if last block of round
		if cmpSwitch(round, currentNumber) == 0 {
			/*if len(parent_nextNodes) != 0  {
				currentRound.nodeIds = convertNodeID(parent_nextNodes)
				currentRound.nodes = make([]*discover.Node, len(parent_nextNodes))
				copy(currentRound.nodes, parent_nextNodes)
			}else {
				if  nil != genesisCurRound {
					currentRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
					copy(currentRound.nodeIds, genesisCurRound.nodeIds)
					currentRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
					copy(currentRound.nodes, genesisCurRound.nodes)
				}
			}*/

			if  nil != genesisCurRound {
				currentRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
				copy(currentRound.nodeIds, genesisCurRound.nodeIds)
				currentRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
				copy(currentRound.nodes, genesisCurRound.nodes)
			}
		}else { // Is'nt last block of round

			/*if len(parent_curNodes) != 0 {
				currentRound.nodeIds = convertNodeID(parent_curNodes)
				currentRound.nodes = make([]*discover.Node, len(parent_curNodes))
				copy(currentRound.nodes, parent_curNodes)
			}else {
				if  nil != genesisCurRound {
					currentRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
					copy(currentRound.nodeIds, genesisCurRound.nodeIds)
					currentRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
					copy(currentRound.nodes, genesisCurRound.nodes)
				}
			}*/

			if  nil != genesisCurRound {
				currentRound.nodeIds = make([]discover.NodeID, len(genesisCurRound.nodeIds))
				copy(currentRound.nodeIds, genesisCurRound.nodeIds)
				currentRound.nodes = make([]*discover.Node, len(genesisCurRound.nodes))
				copy(currentRound.nodes, genesisCurRound.nodes)
			}

		}
	}

	/**
	Sets next info
	  */
	nextRound := &pposRound{}
	// next start, end
	nextRound.start = new(big.Int).Add(start, new(big.Int).SetUint64(uint64(common.BaseSwitchWitness)))
	nextRound.end = new(big.Int).Add(end, new(big.Int).SetUint64(uint64(common.BaseSwitchWitness)))

	log.Debug("Set the farthest allowed cache reserved block: Next round", "start", nextRound.start, "end",nextRound.end)

	// sets next
	if len(curr_NEXT_Nodes) != 0 {
		nextRound.nodeIds = convertNodeID(curr_NEXT_Nodes)
		nextRound.nodes = make([]*discover.Node, len(curr_NEXT_Nodes))
		copy(nextRound.nodes, curr_NEXT_Nodes)
	}else { // Reference parent
		// election index == cur index || election index < cur index < switch index
		if cmpElection(round, currentNumber) == 0 || (cmpElection(round, currentNumber) > 0 && cmpSwitch(round, currentNumber) < 0)  {

			if nil != genesisCurRound {
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

	/*
	Sets nodeCache
	 */
	cache := &nodeCache{
		former: 	formerRound,
		current: 	currentRound,
		next: 		nextRound,
	}
	p.nodeRound.setNodeCache(big.NewInt(int64(currentNumber)), currentHash, cache)
	log.Debug("Set the farthest allowed to cache the information of the reserved block", "currentBlockNum", currentNumber, "currentHash", currentHash.String())
	p.printMapInfo("Set the farthest allowed to cache the information of the reserved block", currentNumber, currentHash)
	return nil
}


func (p *ppos) cleanNodeRound () {
	p.lock.Lock()
	p.nodeRound =  make(roundCache, 0)
	p.lock.Unlock()
}

// election index == cur       0
// cur < election index       -1
// election index < cur        1
// param invalid              -2
func cmpElection (round, currentNumber uint64) int {
	// last num of round
	last := int(round * common.BaseSwitchWitness)
	ele_sub := int(common.BaseSwitchWitness - common.BaseElection)
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
	last := round * common.BaseSwitchWitness
	if last < currentNum {
		return 1
	}else if last == currentNum {
		return 0
	}else {
		return -1
	}
}

func (p *ppos) setPPOS_Temp(){
	p.pposTemp = ppos_storage.GetPPosTempPtr()
}

