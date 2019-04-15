package pposm

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/ppos_storage"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"math/big"
	"net"
	"strconv"
	"strings"
	"sync"
)

const (
	GET_WITNESS   = 1
	GET_IM_RE     = 2
	GET_WIT_IM_RE = 3
)

var (
	//CandidateEncodeErr          = errors.New("Candidate encoding err")
	//CandidateDecodeErr          = errors.New("Candidate decoding err")
	CandidateEmptyErr           = errors.New("Candidate is empty")
	ContractBalanceNotEnoughErr = errors.New("Contract's balance is not enough")
	CandidateOwnerErr           = errors.New("CandidateOwner Addr is illegal")
	DepositLowErr               = errors.New("Candidate deposit too low")
	WithdrawPriceErr            = errors.New("Withdraw Price err")
	WithdrawLowErr              = errors.New("Withdraw Price too low")
	RefundEmptyErr              = errors.New("Refund is empty")
)

type candidateStorage map[discover.NodeID]*types.Candidate
type refundStorage map[discover.NodeID]types.CandidateQueue

type CandidatePool struct {
	// min deposit allow threshold
	threshold *big.Int
	// min deposit limit percentage
	depositLimit uint32
	// allow put into immedidate condition
	allowed uint32
	// allow immediate elected max count
	maxCount uint32
	// allow witness max count
	maxChair uint32
	// allow block interval for refunds
	refundBlockNumber uint32

	// previous witness
	preOriginCandidates candidateStorage
	// current witnesses
	originCandidates candidateStorage
	// next witnesses
	nextOriginCandidates candidateStorage
	// immediates
	immediateCandidates candidateStorage
	// reserves
	reserveCandidates candidateStorage
	// refunds
	defeatCandidates refundStorage

	// cache
	//immediateCacheArr types.CandidateQueue
	//reserveCacheArr   types.CandidateQueue

	storage *ppos_storage.Ppos_storage
}

// Initialize the global candidate pool object
func NewCandidatePool(configs *params.PposConfig) *CandidatePool {

	log.Debug("Build a New CandidatePool Info ...")
	if "" == strings.TrimSpace(configs.CandidateConfig.Threshold) {
		configs.CandidateConfig.Threshold = "1000000000000000000000000"
	}
	var threshold *big.Int
	if thd, ok := new(big.Int).SetString(configs.CandidateConfig.Threshold, 10); !ok {
		threshold, _ = new(big.Int).SetString("1000000000000000000000000", 10)
	} else {
		threshold = thd
	}
	return &CandidatePool{
		threshold:            threshold,
		depositLimit:         configs.CandidateConfig.DepositLimit,
		allowed:              configs.CandidateConfig.Allowed,
		maxCount:             configs.CandidateConfig.MaxCount,
		maxChair:             configs.CandidateConfig.MaxChair,
		refundBlockNumber:    configs.CandidateConfig.RefundBlockNumber,
		preOriginCandidates:  make(candidateStorage, 0),
		originCandidates:     make(candidateStorage, 0),
		nextOriginCandidates: make(candidateStorage, 0),
		immediateCandidates:  make(candidateStorage, 0),
		reserveCandidates:    make(candidateStorage, 0),
		defeatCandidates:     make(refundStorage, 0),
		//immediateCacheArr:    make(types.CandidateQueue, 0),
		//reserveCacheArr:      make(types.CandidateQueue, 0),
	}
}

// flag:
// 0: only init previous witness and current witness and next witness
// 1：init previous witness and current witness and next witness and immediate and reserve
// 2: init all information
//func (c *CandidatePool) initDataByState(state vm.StateDB, flag int) error {
//	log.Info("init data by stateDB...", "statedb addr", fmt.Sprintf("%p", state))
//
//	parentRoutineID := fmt.Sprintf("%s", common.CurrentGoRoutineID())
//
//	//loading  candidates func
//	loadWitFunc := func(title string, canMap candidateStorage,
//		getIndexFn func(state vm.StateDB) ([]discover.NodeID, error),
//		getInfoFn func(state vm.StateDB, id discover.NodeID) (*types.Candidate, error)) error {
//
//		log.Debug("initDataByState by Getting "+title+" parent routine "+parentRoutineID, "statedb addr", fmt.Sprintf("%p", state))
//		var witnessIds []discover.NodeID
//		if ids, err := getIndexFn(state); nil != err {
//			log.Error("Failed to decode "+title+" witnessIds on initDataByState", " err", err)
//			return err
//		} else {
//			witnessIds = ids
//		}
//
//		PrintObject(title+" witnessIds", witnessIds)
//		for _, witnessId := range witnessIds {
//
//			if ca, err := getInfoFn(state, witnessId); nil != err {
//				log.Error("Failed to decode "+title+" witness Candidate on initDataByState", "err", err)
//				return CandidateDecodeErr
//			} else {
//				if nil != ca {
//					PrintObject(title+"Id:"+witnessId.String()+", can", ca)
//					canMap[witnessId] = ca
//				} else {
//					delete(canMap, witnessId)
//				}
//			}
//		}
//		return nil
//	}
//
//	witErrCh := make(chan error, 3)
//	var wg sync.WaitGroup
//	wg.Add(3)
//
//	// loading witnesses
//	go func() {
//		c.preOriginCandidates = make(candidateStorage, 0)
//		witErrCh <- loadWitFunc("previous", c.preOriginCandidates, getPreviousWitnessIdsState, getPreviousWitnessByState)
//		wg.Done()
//	}()
//	go func() {
//		c.originCandidates = make(candidateStorage, 0)
//		witErrCh <- loadWitFunc("current", c.originCandidates, getWitnessIdsByState, getWitnessByState)
//		wg.Done()
//	}()
//	go func() {
//		c.nextOriginCandidates = make(candidateStorage, 0)
//		witErrCh <- loadWitFunc("next", c.nextOriginCandidates, getNextWitnessIdsByState, getNextWitnessByState)
//		wg.Done()
//	}()
//	var err error
//	for i := 1; i <= 3; i++ {
//		if err = <-witErrCh; nil != err {
//			break
//		}
//	}
//	wg.Wait()
//	close(witErrCh)
//	if nil != err {
//		return err
//	}
//
//	// loading elected candidates
//	if flag == 1 || flag == 2 {
//
//		loadElectedFunc := func(title string, canMap candidateStorage,
//			getIndexFn func(state vm.StateDB) ([]discover.NodeID, error),
//			getInfoFn func(state vm.StateDB, id discover.NodeID) (*types.Candidate, error)) (types.CandidateQueue, error) {
//			var witnessIds []discover.NodeID
//
//			log.Debug("initDataByState by Getting "+title+" parent routine "+parentRoutineID, "statedb addr", fmt.Sprintf("%p", state))
//			if ids, err := getIndexFn(state); nil != err {
//				log.Error("Failed to decode "+title+"Ids on initDataByState", " err", err)
//				return nil, err
//			} else {
//				witnessIds = ids
//			}
//			// cache
//			canCache := make(types.CandidateQueue, 0)
//
//			PrintObject(title+" Ids", witnessIds)
//			for _, witnessId := range witnessIds {
//
//				if ca, err := getInfoFn(state, witnessId); nil != err {
//					log.Error("Failed to decode "+title+" Candidate on initDataByState", "err", err)
//					return nil, CandidateDecodeErr
//				} else {
//					if nil != ca {
//						PrintObject(title+"Id:"+witnessId.String()+", can", ca)
//						canMap[witnessId] = ca
//						canCache = append(canCache, ca)
//					} else {
//						delete(canMap, witnessId)
//					}
//				}
//			}
//			return canCache, nil
//		}
//		type result struct {
//			Type int // 1: immediate; 2: reserve
//			Arr  types.CandidateQueue
//			Err  error
//		}
//		resCh := make(chan *result, 2)
//		wg.Add(2)
//		go func() {
//			res := new(result)
//			res.Type = IS_IMMEDIATE
//			c.immediateCandidates = make(candidateStorage, 0)
//			if arr, err := loadElectedFunc("immediate", c.immediateCandidates, getImmediateIdsByState, getImmediateByState); nil != err {
//				res.Err = err
//				resCh <- res
//			} else {
//				res.Arr = arr
//				resCh <- res
//			}
//			wg.Done()
//		}()
//		go func() {
//			res := new(result)
//			res.Type = IS_RESERVE
//			c.reserveCandidates = make(candidateStorage, 0)
//			if arr, err := loadElectedFunc("reserve", c.reserveCandidates, getReserveIdsByState, getReserveByState); nil != err {
//				res.Err = err
//				resCh <- res
//			} else {
//				res.Arr = arr
//				resCh <- res
//			}
//			wg.Done()
//		}()
//		wg.Wait()
//		close(resCh)
//		for res := range resCh {
//			if nil != res.Err {
//				return res.Err
//			}
//			switch res.Type {
//			case IS_IMMEDIATE:
//				c.immediateCacheArr = res.Arr
//			case IS_RESERVE:
//				c.reserveCacheArr = res.Arr
//			default:
//				continue
//			}
//		}
//
//	}
//
//	// load refunds
//	if flag == 2 {
//
//		var defeatIds []discover.NodeID
//		c.defeatCandidates = make(refundStorage, 0)
//		if ids, err := getDefeatIdsByState(state); nil != err {
//			log.Error("Failed to decode defeatIds on initDataByState", "err", err)
//			return err
//		} else {
//			defeatIds = ids
//		}
//		PrintObject("defeatIds", defeatIds)
//		for _, defeatId := range defeatIds {
//			if arr, err := getDefeatsByState(state, defeatId); nil != err {
//				log.Error("Failed to decode defeat CandidateArr on initDataByState", "err", err)
//				return CandidateDecodeErr
//			} else {
//				if nil != arr && len(arr) != 0 {
//					PrintObject("defeatId:"+defeatId.String()+", arr", arr)
//					c.defeatCandidates[defeatId] = arr
//				} else {
//					delete(c.defeatCandidates, defeatId)
//				}
//			}
//		}
//	}
//	return nil
//}

func (c *CandidatePool) initDataByState(state vm.StateDB) {
	c.storage = state.GetPPOSCache()

	log.Debug("initDataByState", "state addr", fmt.Sprintf("%p", state))
	log.Debug("initDataByState", "ppos storage addr", fmt.Sprintf("%p", c.storage))
}

// flag:
// 1: witness
// 2: im and re (im/re)
// 3 : wit + im/re
func (c *CandidatePool) initData2Cache(state vm.StateDB, flag int) {
	c.initDataByState(state)

	loadQueueFunc := func(arr types.CandidateQueue, canMap candidateStorage) {
		for _, can := range arr {
			canMap[can.CandidateId] = can
		}
	}
	var wg sync.WaitGroup

	switch flag {
	case GET_WITNESS:
		wg.Add(3)
		c.getWitnessMap(&wg, loadQueueFunc)
		wg.Wait()
	case GET_IM_RE:
		wg.Add(2)
		c.getImAndReMap(&wg, loadQueueFunc)
		wg.Wait()
	case GET_WIT_IM_RE:
		wg.Add(5)
		c.getWitnessMap(&wg, loadQueueFunc)
		c.getImAndReMap(&wg, loadQueueFunc)
		wg.Wait()
	default:
		return
	}
}

func (c *CandidatePool) getWitnessMap(wg *sync.WaitGroup, loadQueueFunc func(arr types.CandidateQueue, canMap candidateStorage)) {
	go func() {
		loadQueueFunc(c.storage.GetCandidateQueue(ppos_storage.PREVIOUS), c.preOriginCandidates)
		wg.Done()
	}()
	go func() {
		loadQueueFunc(c.storage.GetCandidateQueue(ppos_storage.CURRENT), c.originCandidates)
		wg.Done()
	}()
	go func() {
		loadQueueFunc(c.storage.GetCandidateQueue(ppos_storage.NEXT), c.nextOriginCandidates)
		wg.Done()
	}()
}
func (c *CandidatePool) getImAndReMap(wg *sync.WaitGroup, loadQueueFunc func(arr types.CandidateQueue, canMap candidateStorage)) {
	go func() {
		loadQueueFunc(c.storage.GetCandidateQueue(ppos_storage.IMMEDIATE), c.immediateCandidates)
		wg.Done()
	}()
	go func() {
		loadQueueFunc(c.storage.GetCandidateQueue(ppos_storage.RESERVE), c.reserveCandidates)
		wg.Done()
	}()
}

// pledge Candidate
func (c *CandidatePool) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	log.Debug("Call SetCandidate start ...", "threshold", c.threshold.String(), "depositLimit", c.depositLimit, "allowed", c.allowed, "maxCount", c.maxCount, "maxChair", c.maxChair, "refundBlockNumber", c.refundBlockNumber)

	PrintObject("Call SetCandidate start ...", *can)

	c.initData2Cache(state, GET_IM_RE)
	var nodeIds []discover.NodeID

	// If it is the first pledge, judge the pledge threshold
	if !c.checkFirstThreshold(can) {
		log.Warn("Failed to checkFirstThreshold on SetCandidate", "Deposit", can.Deposit.String(), "threshold", c.threshold)
		return errors.New(DepositLowErr.Error() + ", Current Deposit:" + can.Deposit.String() + ", target threshold:" + fmt.Sprint(c.threshold))
	}

	// Before each pledge, we need to check whether the current can deposit is not less
	// than the minimum can deposit when the corresponding queue to be placed is full.
	//if _, ok := c.checkDeposit(state, can, false); !ok {
	if ok := c.checkDeposit(can); !ok {
		log.Warn("Failed to checkDeposit on SetCandidate", "nodeId", nodeId.String(), " err", DepositLowErr)
		return DepositLowErr
	}
	nodeIds = c.setCandidateInfo(state, nodeId, can, can.BlockNumber, nil)
	//go ticketPool.DropReturnTicket(state, nodeIds...)
	if len(nodeIds) > 0 {
		if err := tContext.DropReturnTicket(state, can.BlockNumber, nodeIds...); nil != err {
			log.Error("Failed to DropReturnTicket on SetCandidate ...",  "current blockNumber", can.BlockNumber, "err", err)
			//return err
		}
	}
	log.Debug("Call SetCandidate successfully...")
	return nil
}

// If TCout is small, you must first move to reserves, otherwise it will be counted.
func (c *CandidatePool) setCandidateInfo(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate, currentBlockNumber *big.Int, promoteReserveFunc func(state vm.StateDB, currentBlockNumber *big.Int)/* []discover.NodeID*/) []discover.NodeID {

	var allowed, delimmediate, delreserve bool
	// check ticket count
	if c.checkTicket(tContext.GetCandidateTicketCount(state, nodeId)) {
		allowed = true
		if _, ok := c.reserveCandidates[can.CandidateId]; ok {
			delreserve = true
		}
		c.immediateCandidates[can.CandidateId] = can
	} else {
		if _, ok := c.immediateCandidates[can.CandidateId]; ok {
			delimmediate = true
		}
		c.reserveCandidates[can.CandidateId] = can
	}

	// delete Func
	delCandidateFunc := func(nodeId discover.NodeID, flag int) {
		queue := c.getCandidateQueue(flag)
		/*//for i, id := range ids {
		for i := 0; i < len(ids); i++ {
			id := ids[i]
			if id == can.CandidateId {
				ids = append(ids[:i], ids[i+1:]...)
				i--
			}
		}*/
		for i, can := range queue {
			if can.CandidateId == nodeId {
				queue = append(queue[:i], queue[i+1:]...)
				break
			}
		}
		c.setCandidateQueue(queue, flag)
	}


	/**
	handle the reserve queue func
	*/
	handleReserveFunc := func(re_queue types.CandidateQueue) []discover.NodeID {

		re_queueCopy := make(types.CandidateQueue, len(re_queue))
		copy(re_queueCopy, re_queue)

		str := "Call setCandidateInfo to handleReserveFunc to sort the reserve queue ..."

		// sort reserve array
		makeCandidateSort(str, state, re_queueCopy)

		nodeIds := make([]discover.NodeID, 0)


		if len(re_queueCopy) > int(c.maxCount) {
			// Intercepting the lost candidates to tmpArr
			tempArr := (re_queueCopy)[c.maxCount:]
			// qualified elected candidates
			re_queueCopy = (re_queueCopy)[:c.maxCount]

			// handle tmpArr
			for _, tmpCan := range tempArr {
				deposit, _ := new(big.Int).SetString(tmpCan.Deposit.String(), 10)
				refund := &types.CandidateRefund{
					Deposit:     deposit,
					BlockNumber: big.NewInt(currentBlockNumber.Int64()),
					Owner:       tmpCan.Owner,
				}
				c.setRefund(tmpCan.CandidateId, refund)
				nodeIds = append(nodeIds, tmpCan.CandidateId)
			}
		}

		c.setCandidateQueue(re_queueCopy, ppos_storage.RESERVE)

		return nodeIds
	}


	var str string
	// using the cache handle current queue
	cacheArr := make(types.CandidateQueue, 0)
	if allowed {

		/** first delete this can on reserves */
		if delreserve {
			delCandidateFunc(can.CandidateId, ppos_storage.RESERVE)
		}

		str = "Call setCandidateInfo to sort the immediate queue ..."
		for _, v := range c.immediateCandidates {
			cacheArr = append(cacheArr, v)
		}
	} else {

		/** first delete this can on immediates */
		if delimmediate {
			delCandidateFunc(can.CandidateId, ppos_storage.IMMEDIATE)
		}


		str = "Call setCandidateInfo to sort the reserve queue ..."
		for _, v := range c.reserveCandidates {
			cacheArr = append(cacheArr, v)
		}
	}


	// sort cache array
	makeCandidateSort(str, state, cacheArr)

	nodeIds := make([]discover.NodeID, 0)

	if len(cacheArr) > int(c.maxCount) {
		// Intercepting the lost candidates to tmpArr
		tempArr := (cacheArr)[c.maxCount:]
		// qualified elected candidates
		cacheArr = (cacheArr)[:c.maxCount]

		// add reserve queue cache
		addreserveQueue := make(types.CandidateQueue, 0)

		// handle tmpArr
		for _, tmpCan := range tempArr {


			// if ticket count great allowed && no need delete reserve
			// so this can move to reserve from immediate now
			if allowed {
				addreserveQueue = append(addreserveQueue, tmpCan)
			} else {

				deposit, _ := new(big.Int).SetString(tmpCan.Deposit.String(), 10)
				refund := &types.CandidateRefund{
					Deposit:     deposit,
					BlockNumber: big.NewInt(currentBlockNumber.Int64()),
					Owner:       tmpCan.Owner,
				}
				c.setRefund(tmpCan.CandidateId, refund)
				nodeIds = append(nodeIds, tmpCan.CandidateId)
			}

		}

		if len(addreserveQueue) != 0 {
			re_queue := c.getCandidateQueue(ppos_storage.RESERVE)
			re_queue = append(re_queue, addreserveQueue...)
			if ids := handleReserveFunc(re_queue); len(ids) != 0 {
				nodeIds = append(nodeIds, ids...)
			}
			//c.setCandidateQueue(re_queue, ppos_storage.RESERVE)
		}

	}

	if allowed {
		c.setCandidateQueue(cacheArr, ppos_storage.IMMEDIATE)
	} else {
		c.setCandidateQueue(cacheArr, ppos_storage.RESERVE)
	}

	if nil != promoteReserveFunc {
		promoteReserveFunc(state, currentBlockNumber)
	}

	return nodeIds
}



// If TCout is small, you must first move to reserves, otherwise it will be counted.
func (c *CandidatePool) electionUpdateCanById(state vm.StateDB, currentBlockNumber *big.Int, nodeIds ... discover.NodeID) []discover.NodeID {

	im_del_temp := make(candidateStorage, 0)
	re_del_temp := make(candidateStorage, 0)

	for _, nodeId := range nodeIds {
		// check ticket count
		if c.checkTicket(tContext.GetCandidateTicketCount(state, nodeId)) {
			if can, ok := c.reserveCandidates[nodeId]; ok {
				re_del_temp[nodeId] = can
			}

		} else {
			if can, ok := c.immediateCandidates[nodeId]; ok {
				im_del_temp[nodeId] = can
			}
		}
	}

	if len(im_del_temp) == 0 && len(re_del_temp) == 0 {
		log.Debug("Call Election to electionUpdateCanById, had not change on double queue ...")
		return nil
	}

	im_queue := c.getCandidateQueue(ppos_storage.IMMEDIATE)

	re_queue := c.getCandidateQueue(ppos_storage.RESERVE)


	PrintObject("Call Election to electionUpdateCanById, Before shuffle double queue the old immediate queue len:=" + fmt.Sprint(len(im_queue)) + " ,queue is", im_queue)

	PrintObject("Call Election to electionUpdateCanById, Before shuffle double queue the old reserve queue len:=" + fmt.Sprint(len(re_queue)) + " ,queue is", re_queue)

	// immediate queue
	for _, can := range re_del_temp{
		im_queue = append(im_queue, can)
	}
	// for i := 0; i < len(ids); i++ {
	for i := 0; i < len(im_queue); i++ {
		im := im_queue[i]
		if _, ok := im_del_temp[im.CandidateId]; ok {
			im_queue = append(im_queue[:i], im_queue[i+1:]...)
			i--
		}
	}

	// reserve  queue
	for i := 0; i < len(re_queue); i++ {
		re := re_queue[i]
		if _, ok := re_del_temp[re.CandidateId]; ok {
			re_queue = append(re_queue[:i], re_queue[i+1:]...)
			i--
		}
	}
	for _, can := range im_del_temp {
		re_queue = append(re_queue, can)
	}


	PrintObject("Call Election to electionUpdateCanById, After shuffle double queue the old immediate queue len:=" + fmt.Sprint(len(im_queue)) + " ,queue is", im_queue)

	PrintObject("Call Election to electionUpdateCanById, After shuffle double queue the old reserve queue len:=" + fmt.Sprint(len(re_queue)) + " ,queue is", re_queue)


	/**
	 handle the reserve queue func
	 */
	handleReserveFunc := func(re_queue types.CandidateQueue) []discover.NodeID {

		re_queueCopy := make(types.CandidateQueue, len(re_queue))
		copy(re_queueCopy, re_queue)

		str := "Call Election to handleReserveFunc to sort the reserve queue ..."

		// sort reserve array
		makeCandidateSort(str, state, re_queueCopy)

		nodeIds := make([]discover.NodeID, 0)


		if len(re_queueCopy) > int(c.maxCount) {
			// Intercepting the lost candidates to tmpArr
			tempArr := (re_queueCopy)[c.maxCount:]
			// qualified elected candidates
			re_queueCopy = (re_queueCopy)[:c.maxCount]

			// handle tmpArr
			for _, tmpCan := range tempArr {
				deposit, _ := new(big.Int).SetString(tmpCan.Deposit.String(), 10)
				refund := &types.CandidateRefund{
					Deposit:     deposit,
					BlockNumber: big.NewInt(currentBlockNumber.Int64()),
					Owner:       tmpCan.Owner,
				}
				c.setRefund(tmpCan.CandidateId, refund)
				nodeIds = append(nodeIds, tmpCan.CandidateId)
			}
		}

		PrintObject("Call Election to electionUpdateCanById to handleReserveFunc, Finally shuffle double queue the old reserve queue len:=" + fmt.Sprint(len(re_queueCopy)) + " ,queue is", re_queueCopy)

		c.setCandidateQueue(re_queueCopy, ppos_storage.RESERVE)

		return nodeIds
	}





	str := "Call Election to start sort immediate queue ..."

	// sort immediate array
	makeCandidateSort(str, state, im_queue)

	nodeIdArr := make([]discover.NodeID, 0)

	if len(im_queue) > int(c.maxCount) {
		// Intercepting the lost candidates to tmpArr
		tempArr := (im_queue)[c.maxCount:]
		// qualified elected candidates
		im_queue = (im_queue)[:c.maxCount]

		// add reserve queue cache
		addreserveQueue := make(types.CandidateQueue, 0)

		// handle tmpArr
		for _, tmpCan := range tempArr {
			addreserveQueue = append(addreserveQueue, tmpCan)
		}

		if len(addreserveQueue) != 0 {
			// append into reserve queue
			re_queue = append(re_queue, addreserveQueue...)
		}

	}/*else {
		PrintObject("Call Election to electionUpdateCanById to direct, Finally shuffle double queue the old reserve queue is", re_queue)
		c.setCandidateQueue(re_queue, ppos_storage.RESERVE)
	}*/

	// hanle and sets reserve queue
	if ids := handleReserveFunc(re_queue); len(ids) != 0 {
		nodeIdArr = append(nodeIdArr, ids...)
	}

	PrintObject("Call Election to electionUpdateCanById, Finally shuffle double queue the old immediate queue len:=" + fmt.Sprint(len(im_queue)) + " ,queue is", im_queue)
	// set immediate queue
	c.setCandidateQueue(im_queue, ppos_storage.IMMEDIATE)

	c.promoteReserveQueue(state, currentBlockNumber)

	return nodeIdArr
}



// Getting immediate or reserve candidate info by nodeId
func (c *CandidatePool) GetCandidate(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) *types.Candidate {
	return c.getCandidate(state, nodeId, blockNumber)
}

// Getting immediate or reserve candidate info arr by nodeIds
func (c *CandidatePool) GetCandidateArr(state vm.StateDB, blockNumber *big.Int, nodeIds ...discover.NodeID) types.CandidateQueue {
	return c.getCandidates(state, blockNumber, nodeIds...)
}

// candidate withdraw from immediates or reserve elected candidates
func (c *CandidatePool) WithdrawCandidate(state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error {
	log.Info("WithdrawCandidate...", "nodeId", nodeId.String(), "price", price.String(), "blockNumber", blockNumber.String(), "threshold", c.threshold.String(), "depositLimit", c.depositLimit, "allowed", c.allowed, "maxCount", c.maxCount, "maxChair", c.maxChair, "refundBlockNumber", c.refundBlockNumber)

	c.initData2Cache(state, GET_IM_RE)
	var nodeIds []discover.NodeID

	if arr, err := c.withdrawCandidate(state, nodeId, price, blockNumber); nil != err {
		return err
	} else {
		nodeIds = arr
	}
	//go ticketPool.DropReturnTicket(state, nodeIds...)
	if len(nodeIds) > 0 {
		if err := tContext.DropReturnTicket(state, blockNumber, nodeIds...); nil != err {
			log.Error("Failed to DropReturnTicket on WithdrawCandidate ...", "blockNumber", blockNumber.String(), "err", err)
		}
	}
	return nil
}

func (c *CandidatePool) withdrawCandidate(state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) ([]discover.NodeID, error) {

	if price.Cmp(new(big.Int).SetUint64(0)) <= 0 {
		log.Error("Failed to WithdrawCandidate price invalid", "blockNumber", blockNumber.String(), "nodeId", nodeId.String(), " price", price.String())
		return nil, WithdrawPriceErr
	}

	// cache
	var can *types.Candidate
	var isImmediate bool

	if imCan, ok := c.immediateCandidates[nodeId]; !ok {
		reCan, ok := c.reserveCandidates[nodeId]
		if !ok {
			log.Error("Failed to WithdrawCandidate current Candidate is empty", "blockNumber", blockNumber.String(), "nodeId", nodeId.String(), "price", price.String())
			return nil, CandidateEmptyErr
		} else {
			can = reCan
		}
	} else {
		can = imCan
		isImmediate = true
	}


	// delete Func
	delCandidateFunc := func(nodeId discover.NodeID, flag int) {
		queue := c.getCandidateQueue(flag)

		for i, can := range queue {
			if can.CandidateId == nodeId {
				queue = append(queue[:i], queue[i+1:]...)
				break
			}
		}
		c.setCandidateQueue(queue, flag)
	}

	var nodeIdArr []discover.NodeID
	deposit, _ := new(big.Int).SetString(can.Deposit.String(), 10)
	// check withdraw price
	if can.Deposit.Cmp(price) < 0 {
		log.Error("Failed to WithdrawCandidate refund price must less or equal deposit", "blockNumber", blockNumber.String(), "nodeId", nodeId.String(), " price", price.String())
		return nil, WithdrawPriceErr
	} else if can.Deposit.Cmp(price) == 0 { // full withdraw

		log.Info("WithdrawCandidate into full withdraw", "blockNumber", blockNumber.String(), "canId", can.CandidateId.String(), "current can deposit", can.Deposit.String(), "withdraw price is", price.String())

		if isImmediate {
			delCandidateFunc(can.CandidateId, ppos_storage.IMMEDIATE)
		} else {
			delCandidateFunc(can.CandidateId, ppos_storage.RESERVE)
		}

		refund := &types.CandidateRefund{
			Deposit:     deposit,
			BlockNumber: big.NewInt(blockNumber.Int64()),
			Owner:       can.Owner,
		}

		c.setRefund(can.CandidateId, refund)

		c.promoteReserveQueue(state, blockNumber)

		nodeIds := []discover.NodeID{nodeId}

		nodeIdArr = nodeIds

	} else { // withdraw a few ...
		/*// Only withdraw part of the refunds, need to reorder the immediate elected candidates
		// The remaining candiate price to update current candidate info

		log.Info("WithdrawCandidate into withdraw a few", "canId", can.CandidateId.String(), "current can deposit", can.Deposit.String(), "withdraw price is", price.String())

		if err := c.checkWithdraw(can.Deposit, price); nil != err {
			log.Error("Failed to price invalid on WithdrawCandidate", " price", price.String(), "err", err)
			return nil, err
		}

		remainMoney := new(big.Int).Sub(can.Deposit, price)
		half_threshold := new(big.Int).Mul(c.threshold, big.NewInt(2))

		var isFullWithdraw bool
		var canNew *types.Candidate
		var refund *types.CandidateRefund

		if remainMoney.Cmp(half_threshold) < 0 { // full withdraw
			isFullWithdraw = true

			refund = &types.CandidateRefund{
				Deposit:     deposit,
				BlockNumber: big.NewInt(blockNumber.Int64()),
				Owner:       can.Owner,
			}

		} else {
			// remain info
			canNew = &types.Candidate{
				Deposit:     remainMoney,
				BlockNumber: big.NewInt(can.BlockNumber.Int64()),
				TxIndex:     can.TxIndex,
				CandidateId: can.CandidateId,
				Host:        can.Host,
				Port:        can.Port,
				Owner:       can.Owner,
				Extra:       can.Extra,
				Fee:         can.Fee,
			}

			refund = &types.CandidateRefund{
				Deposit:     price,
				BlockNumber: big.NewInt(blockNumber.Int64()),
				Owner:       can.Owner,
			}

		}

		if isFullWithdraw {
			if isImmediate {
				delCandidateFunc(can.CandidateId, ppos_storage.IMMEDIATE)
			} else {
				delCandidateFunc(can.CandidateId, ppos_storage.RESERVE)
			}

			c.setRefund(can.CandidateId, refund)

			nIds := c.shuffleQueue(state, blockNumber)
			nodeIds := []discover.NodeID{nodeId}

			if len(nIds) != 0 {
				nodeIds = append(nodeIds, nIds...)
			}

			nodeIdArr = nodeIds
		} else {

			handleFunc := func(nodeId discover.NodeID, flag int) {
				queue := c.getCandidateQueue(flag)

				for i, can := range queue {
					if can.CandidateId == nodeId {
						queue[i] = canNew
						break
					}
				}
				c.setCandidateQueue(queue, flag)
				c.setRefund(nodeId, refund)
			}

			if isImmediate {
				handleFunc(can.CandidateId, ppos_storage.IMMEDIATE)
			} else {
				handleFunc(can.CandidateId, ppos_storage.RESERVE)
			}
			nodeIdArr = append(nodeIdArr, nodeId)
			if arr := c.shuffleQueue(state, blockNumber); len(arr) != 0 {
				nodeIdArr = append(nodeIdArr, arr...)
			}
		}*/
		log.Error("Failed to WithdrawCandidate, must full withdraw", "blockNumber", blockNumber.String(), "nodeId", nodeId.String(), "the can deposit", can.Deposit.String(), "current will withdraw price", price.String())
		return nil, WithdrawLowErr
	}
	log.Info("Call WithdrawCandidate SUCCESS !!!!!!!!!!!!")
	return nodeIdArr, nil
}

// Getting elected candidates array
// flag:
// 0:  Getting all elected candidates array
// 1:  Getting all immediate elected candidates array
// 2:  Getting all reserve elected candidates array
func (c *CandidatePool) GetChosens(state vm.StateDB, flag int, blockNumber *big.Int) types.KindCanQueue {
	log.Debug("Call GetChosens getting immediate or reserve candidates ...", "blockNumber", blockNumber.String(), "flag", flag)
	c.initDataByState(state)
	im := make(types.CandidateQueue, 0)
	re := make(types.CandidateQueue, 0)
	arr := make(types.KindCanQueue, 0)
	if flag == 0 || flag == 1 {
		im = c.getCandidateQueue(ppos_storage.IMMEDIATE)
	}
	if flag == 0 || flag == 2 {
		re = c.getCandidateQueue(ppos_storage.RESERVE)

	}
	arr = append(arr, im, re)
	PrintObject("GetChosens return", arr)
	return arr
}

// Getting elected candidates array
// flag:
// 0:  Getting all elected candidates array
// 1:  Getting all immediate elected candidates array
// 2:  Getting all reserve elected candidates array
func (c *CandidatePool) GetCandidatePendArr(state vm.StateDB, flag int, blockNumber *big.Int) types.CandidateQueue {
	log.Debug("Call GetCandidatePendArr getting immediate candidates ...", "blockNumber", blockNumber.String(), "flag", flag)
	c.initDataByState(state)
	arr := make(types.CandidateQueue, 0)
	if flag == 0 || flag == 1 {
		if queue := c.getCandidateQueue(ppos_storage.IMMEDIATE); len(queue) != 0 {
			arr = append(arr, queue...)
		}
	}
	if flag == 0 || flag == 2 {
		if queue := c.getCandidateQueue(ppos_storage.RESERVE); len(queue) != 0 {
			arr = append(arr, queue...)
		}
	}
	PrintObject("GetChosens ==>", arr)
	return arr
}

// Getting current witness array
func (c *CandidatePool) GetChairpersons(state vm.StateDB, blockNumber *big.Int) types.CandidateQueue {
	log.Debug("Call GetChairpersons getting witnesses ...", "blockNumber", blockNumber.String())
	c.initDataByState(state)
	return c.getCandidateQueue(ppos_storage.CURRENT)
}

// Getting all refund array by nodeId
func (c *CandidatePool) GetDefeat(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) types.RefundQueue {
	log.Debug("Call GetDefeat getting defeat arr", "blockNumber", blockNumber, "nodeId", nodeId.String())
	c.initDataByState(state)
	return c.getRefunds(nodeId)
}

// Checked current candidate was defeat by nodeId
func (c *CandidatePool) IsDefeat(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) bool {
	log.Debug("Call IsDefeat", "blockNumber", blockNumber.String())

	c.initData2Cache(state, GET_IM_RE)

	if _, ok := c.immediateCandidates[nodeId]; ok {
		return false
	}
	if _, ok := c.reserveCandidates[nodeId]; ok {
		return false
	}
	if queue := c.getRefunds(nodeId); len(queue) != 0 {
		return true
	}
	return false
}

func (c *CandidatePool) IsChosens(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) bool {
	log.Debug("Call IsChosens", "blockNumber", blockNumber.String())
	c.initData2Cache(state, GET_IM_RE)
	if _, ok := c.immediateCandidates[nodeId]; ok {
		return true
	}
	if _, ok := c.reserveCandidates[nodeId]; ok {
		return true
	}
	return false
}

// Getting owner's address of candidate info by nodeId
func (c *CandidatePool) GetOwner(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) common.Address {
	log.Debug("Call GetOwner", "blockNumber", blockNumber.String(), "curr nodeId", nodeId.String())

	//c.initData2Cache(state, GET_WIT_IM_RE)
	c.initData2Cache(state, GET_IM_RE)

	/*pre_can, pre_ok := c.preOriginCandidates[nodeId]
	or_can, or_ok := c.originCandidates[nodeId]
	ne_can, ne_ok := c.nextOriginCandidates[nodeId]*/

	im_can, im_ok := c.immediateCandidates[nodeId]
	re_can, re_ok := c.reserveCandidates[nodeId]

	queue := c.getRefunds(nodeId)

	de_ok := len(queue) != 0

	/*
	if pre_ok {
		return pre_can.Owner
	}
	if or_ok {
		return or_can.Owner
	}
	if ne_ok {
		return ne_can.Owner
	}*/
	if im_ok {
		return im_can.Owner
	}
	if re_ok {
		return re_can.Owner
	}
	if de_ok {
		return queue[0].Owner
	}
	return common.Address{}
}

// refund once
func (c *CandidatePool) RefundBalance(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) error {

	log.Info("Call RefundBalance",  "curr blocknumber", blockNumber.String(), "curr nodeId", nodeId.String(), "threshold", c.threshold.String(), "depositLimit", c.depositLimit, "allowed", c.allowed, "maxCount", c.maxCount, "maxChair", c.maxChair, "refundBlockNumber", c.refundBlockNumber)

	c.initDataByState(state)
	queueCopy := c.getRefunds(nodeId)

	if len(queueCopy) == 0 {
		log.Warn("Warning Call RefundBalance the refund is empty")
		return RefundEmptyErr
	}
	// cache
	// Used for verification purposes, that is, the beneficiary in the pledge refund information of each nodeId should be the same
	var addr common.Address
	// Grand total refund amount for one-time
	amount := big.NewInt(0)

	// cantract balance
	contractBalance := state.GetBalance(common.CandidatePoolAddr)


	PrintObject("Call RefundBalance Into a few RefundBlockNumber Remain Refund Arr ,Before  Calculate  curr blocknumber:" + blockNumber.String() + " ,len:=" + fmt.Sprint(len(queueCopy)) + " ,refunds is", queueCopy)

	// Traverse all refund information belong to this nodeId
	for index := 0; index < len(queueCopy); index++ {
		refund := queueCopy[index]
		sub := new(big.Int).Sub(blockNumber, refund.BlockNumber)
		log.Info("Check defeat detail on RefundBalance", "nodeId:", nodeId.String(), "curr blocknumber:", blockNumber.String(), "withdraw candidate blocknumber:", refund.BlockNumber.String(), " diff:", sub.String(), "config.RefundBlockNumber", c.refundBlockNumber)
		if sub.Cmp(new(big.Int).SetUint64(uint64(c.refundBlockNumber))) >= 0 { // allow refund

			queueCopy = append(queueCopy[:index], queueCopy[index+1:]...)
			index--
			// add up the refund price
			amount = new(big.Int).Add(amount, refund.Deposit)

		} else {
			log.Warn("block height number had mismatch, No refunds allowed on RefundBalance", "curr blocknumber:", blockNumber.String(), "deposit block height", refund.BlockNumber.String(), "nodeId", nodeId.String(), "allowed block interval", c.refundBlockNumber)
			continue
		}

		if addr == common.ZeroAddr {
			addr = refund.Owner
		} else {
			if addr != refund.Owner {
				log.Error("Failed to RefundBalance Different beneficiary addresses under the same node", "curr blocknumber:", blockNumber.String(), "nodeId", nodeId.String(), "addr1", addr.String(), "addr2", refund.Owner)
				return CandidateOwnerErr
			}
		}

		// check contract account balance
		if (contractBalance.Cmp(amount)) < 0 {
			PrintObject("Failed to RefundBalance constract account insufficient balance ,curr blocknumber:" + blockNumber.String() + ",len:=" + fmt.Sprint(len(queueCopy)) + " ,remain refunds is", queueCopy)
			log.Error("Failed to RefundBalance constract account insufficient balance ", "curr blocknumber:", blockNumber.String(), "nodeId", nodeId.String(), "contract's balance", state.GetBalance(common.CandidatePoolAddr).String(), "amount", amount.String())
			return ContractBalanceNotEnoughErr
		}
	}

	PrintObject("Call RefundBalance Into a few RefundBlockNumber Remain Refund Arr , After Calculate curr blocknumber:" + blockNumber.String() + " ,len:=" + fmt.Sprint(len(queueCopy)) + " ,refunds is", queueCopy)

	// update the tire
	if len(queueCopy) == 0 { // full RefundBlockNumber
		log.Info("Call RefundBalance Into full RefundBlockNumber ...", "curr blocknumber:", blockNumber.String(), "nodeId", nodeId.String())
		c.delRefunds(nodeId)
	} else {
		log.Info("Call RefundBalance Into a few RefundBlockNumber ...", "curr blocknumber:", blockNumber.String(), "nodeId", nodeId.String())
		// If have some remaining, update that
		c.setRefunds(nodeId, queueCopy)
	}
	log.Info("Call RefundBalance to tansfer value：", "curr blocknumber:", blockNumber.String(), "nodeId", nodeId.String(), "contractAddr", common.CandidatePoolAddr.String(),
		"owner's addr", addr.String(), "Return the amount to be transferred:", amount.String())

	// sub contract account balance
	state.SubBalance(common.CandidatePoolAddr, amount)
	// add owner balace
	state.AddBalance(addr, amount)
	log.Debug("Call RefundBalance success ...")
	return nil
}

// set elected candidate extra value
func (c *CandidatePool) SetCandidateExtra(state vm.StateDB, nodeId discover.NodeID, extra string) error {

	log.Info("Call SetCandidateExtra:", "nodeId", nodeId.String(), "extra", extra)

	c.initDataByState(state)
	im_queue := c.getCandidateQueue(ppos_storage.IMMEDIATE)

	for i, can := range im_queue {
		if can.CandidateId == nodeId {
			can.Extra = extra
			im_queue[i] = can
			c.setCandidateQueue(im_queue, ppos_storage.IMMEDIATE)
			log.Debug("Call SetCandidateExtra SUCCESS !!!!!! ")
			return nil
		}
	}

	re_queue := c.getCandidateQueue(ppos_storage.RESERVE)

	for i, can := range re_queue {
		if can.CandidateId == nodeId {
			can.Extra = extra
			re_queue[i] = can
			c.setCandidateQueue(re_queue, ppos_storage.RESERVE)
			log.Debug("Call SetCandidateExtra SUCCESS !!!!!! ")
			return nil
		}
	}
	return CandidateEmptyErr
}

// Announce witness
func (c *CandidatePool) Election(state *state.StateDB, parentHash common.Hash, currBlockNumber *big.Int) ([]*discover.Node, error) {
	log.Info("Call Election start ...", "current blockNumber", currBlockNumber.String(), "threshold", c.threshold.String(), "depositLimit", c.depositLimit, "allowed", c.allowed, "maxCount", c.maxCount, "maxChair", c.maxChair, "refundBlockNumber", c.refundBlockNumber)
	c.initData2Cache(state, GET_IM_RE)

	var nodes []*discover.Node
	var nextQueue types.CandidateQueue
	var isEmptyElection bool

	if nodeArr, canArr, flag, err := c.election(state, parentHash, currBlockNumber); nil != err {
		return nil, err
	} else {
		nodes, nextQueue, isEmptyElection = nodeArr, canArr, flag
	}


	nodeIds := make([]discover.NodeID, 0)

	for _, can := range nextQueue {
		// Release lucky ticket TODO
		if (common.Hash{}) != can.TxHash {
			if err := tContext.ReturnTicket(state, can.CandidateId, can.TxHash, currBlockNumber); nil != err {
				log.Error("Failed to ReturnTicket on Election", "current blockNumber", currBlockNumber.String(), "nodeId", can.CandidateId.String(), "ticketId", can.TxHash.String(), "err", err)
				continue
			}

			if !isEmptyElection {
				nodeIds = append(nodeIds, can.CandidateId)
			}
		}

		/**
		handing before  Re-pledging
		*/
		/*if flag, nIds := c.repledgCheck(state, can, currBlockNumber); !flag {
			nodeIds = append(nodeIds, nIds...)
			// continue handle next one
			continue
		}*/

		/*if !isEmptyElection {
			PrintObject("Election Update Candidate to SetCandidate again ...", *can)
			// Because you need to first ensure if you are in immediates, and if so, move to reserves
			if ids := c.setCandidateInfo(state, can.CandidateId, can, currBlockNumber, c.promoteReserveQueue); len(ids) != 0 {
				nodeIds = append(nodeIds, ids...)
			}
		}*/

	}


	//var dropTick_nodeIds []discover.NodeID

	// finally update the double queue once
	if !isEmptyElection && len(nodeIds) != 0 {

		nodeIds = c.electionUpdateCanById(state, currBlockNumber, nodeIds...)
	}

	// Release the lost list
	//go ticketPool.DropReturnTicket(state, nodeIds...)
	if len(nodeIds) > 0 {
		if err := tContext.DropReturnTicket(state, currBlockNumber, nodeIds...); nil != err {
			log.Error("Failed to DropReturnTicket on Election ...", "current blockNumber", currBlockNumber.String(), "err", err)
		}
	}
	return nodes, nil
}

//return params
// []*discover.Node:  		the cleaned nodeId
// types.CandidateQueue:	the next witness
// bool:					is empty election
// error:					err
func (c *CandidatePool) election(state *state.StateDB, parentHash common.Hash, blockNumber *big.Int) ([]*discover.Node, types.CandidateQueue, bool, error) {

	imm_queue := c.getCandidateQueue(ppos_storage.IMMEDIATE)

	str := "When Election, to election to sort immediate queue ..."
	// sort immediate candidates
	makeCandidateSort(str, state, imm_queue)

	log.Info("When Election, Sorted the immediate array length:", "current blockNumber", blockNumber.String(), "len", len(imm_queue))
	PrintObject("When Election, Sorted the immediate array: current blockNumber:" + blockNumber.String() + ":", imm_queue)
	// cache ids
	immediateIds := make([]discover.NodeID, len(imm_queue))
	for i, can := range imm_queue {
		immediateIds[i] = can.CandidateId
	}
	PrintObject("When Election, current immediate is: current blockNumber:" + blockNumber.String() + " ,len:=" + fmt.Sprint(len(immediateIds)) + " ,arr is:", immediateIds)

	// a certain number of witnesses in front of the cache
	var nextIdArr []discover.NodeID
	// If the number of candidate selected does not exceed the number of witnesses
	if len(immediateIds) <= int(c.maxChair) {
		nextIdArr = make([]discover.NodeID, len(immediateIds))
		copy(nextIdArr, immediateIds)

	} else {
		// If the number of candidate selected exceeds the number of witnesses, the top N is extracted.
		nextIdArr = make([]discover.NodeID, c.maxChair)
		copy(nextIdArr, immediateIds)
	}

	log.Info("When Election, Selected next round of witnesses Ids's count:", "current blockNumber", blockNumber.String(), "len", len(nextIdArr))
	PrintObject("When Election, Selected next round of witnesses Ids: current blockNumber:" + blockNumber.String() + ":", nextIdArr)

	nextQueue := make(types.CandidateQueue, len(nextIdArr))


	for i, next_canId := range nextIdArr {
		im_can := c.immediateCandidates[next_canId]
		// deepCopy
		can := *im_can
		nextQueue[i] = &can
	}


	log.Info("When Election, the count of the copy the witness info from immediate:", "current blockNumber", blockNumber.String(), "len", len(nextQueue))
	PrintObject("When Election, the information of the copy the witness info from immediate: current blockNumber:" + blockNumber.String() + ":", nextQueue)

	// clear all old nextwitnesses information （If it is forked, the next round is no empty.）
	c.delCandidateQueue(ppos_storage.NEXT)

	nodeArr := make([]*discover.Node, 0)
	// Check election whether it's empty or not
	nextQueue, isEmptyElection := c.handleEmptyElection(nextQueue)

	// handle all next witness information
	for i, can := range nextQueue {
		// After election to call Selected LuckyTicket TODO
		luckyId, err := tContext.SelectionLuckyTicket(state, can.CandidateId, parentHash)
		if nil != err {
			log.Error("Failed to take luckyId on Election", "current blockNumber", blockNumber.String(), "nodeId", can.CandidateId.String(), "err", err)
			return nil, nil, false, errors.New(err.Error() + ", nodeId: " + can.CandidateId.String())
		}
		log.Debug("Call Election, select lucky ticket Id is", "current blockNumber", blockNumber.String(), "lucky ticket Id", luckyId.Hex())
		if can.TxHash != luckyId {
			can.TxHash = luckyId
			if luckyId == (common.Hash{}) {
				can.TOwner = common.Address{}
			} else {
				if tick := tContext.GetTicket(state, luckyId); nil != tick {
					can.TOwner = tick.Owner
				}else {
					can.TOwner = common.Address{}
					log.Error("Failed to Gets lucky ticketInfo on Election", "current blockNumber", blockNumber.String(), "nodeId", can.CandidateId.String(), "lucky ticketId", luckyId.Hex())
				}

			}
		}

		nextQueue[i] = can

		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build Node on Election", "current blockNumber", blockNumber.String(), "nodeId", can.CandidateId.String(), "err", err)
			continue
		} else {
			nodeArr = append(nodeArr, node)
		}
	}

	// set next witness
	c.setCandidateQueue(nextQueue, ppos_storage.NEXT)

	log.Info("When Election,next round witness node count is:", "current blockNumber", blockNumber.String(), "len", len(nodeArr))
	PrintObject("When Election,next round witness node information is: current blockNumber:" + blockNumber.String() + ":", nodeArr)
	log.Info("Call Election SUCCESS !!!!!!!", "current blockNumber", blockNumber.String())
	return nodeArr, nextQueue, isEmptyElection, nil
}

// return params
// types.CandidateQueue：the next witness
// bool: is empty election
func (c *CandidatePool) handleEmptyElection(nextQueue types.CandidateQueue) (types.CandidateQueue, bool) {
	// There is'nt empty election
	if len(nextQueue) != 0 {
		return nextQueue, false
	}
	log.Info("Call Election, current is emptyElection, we take current witness become next witness ...")
	// empty election
	// Otherwise, it means that the next witness is nil, then we need to check whether the current round has a witness.
	// If had, use the current round of witness as the next witness,
	// [remark]: The pool of rewards need to use the witness can info
	return c.getCandidateQueue(ppos_storage.CURRENT).DeepCopy(), true
}

/*// The operation before re-setcandiate after election
// false: direct get out
// true: pass
func (c *CandidatePool) repledgCheck(state vm.StateDB, can *types.Candidate, currentBlockNumber *big.Int) (bool, []discover.NodeID) {

	nodeIds := make([]discover.NodeID, 0)
	// If the verification does not pass,
	// it will drop the list directly,
	// but before the list is dropped,
	// it needs to determine which queue was in the queue.
	if _, ok := c.checkDeposit(state, can, true); !ok {
		log.Warn("Failed to checkDeposit on preElectionReset", "nodeId", can.CandidateId.String(), " err", DepositLowErr)
		var del int // del: 1 del immiedate; 2  del reserve
		if _, ok := c.immediateCandidates[can.CandidateId]; ok {
			del = IS_IMMEDIATE
		}
		if _, ok := c.reserveCandidates[can.CandidateId]; ok {
			del = IS_RESERVE
		}

		delCanFunc := func(flag int) {
			queue := c.getCandidateQueue(flag)

			for i := 0; i < len(queue); i++ {
				ca := queue[i]
				if ca.CandidateId == can.CandidateId {
					queue = append(queue[:i], queue[i+1:]...)
					break
				}
			}

			c.setCandidateQueue(queue, flag)
			deposit, _ := new(big.Int).SetString(can.Deposit.String(), 10)
			refund := &types.CandidateRefund{
				Deposit:     deposit,
				BlockNumber: big.NewInt(currentBlockNumber.Int64()),
				Owner:       can.Owner,
			}
			c.setRefund(can.CandidateId, refund)
			nodeIds = append(nodeIds, can.CandidateId)

			*//*nIds := c.promoteReserveQueue(state, currentBlockNumber)

			if len(nIds) != 0 {
				nodeIds = append(nodeIds, nIds...)
			}*//*

			c.promoteReserveQueue(state, currentBlockNumber)
		}
		if del == IS_IMMEDIATE {
			*//** first delete this can on immediates *//*
			delCanFunc(ppos_storage.IMMEDIATE)
			return false, nodeIds
		} else {
			*//** first delete this can on reserves *//*
			delCanFunc(ppos_storage.RESERVE)
			return false, nodeIds
		}
	}
	return true, nodeIds
}*/

// switch next witnesses to current witnesses
func (c *CandidatePool) Switch(state *state.StateDB, blockNumber *big.Int) bool {
	log.Info("Call Switch start ...", "blockNumber", blockNumber.String())
	c.initDataByState(state)
	curr_queue := c.getCandidateQueue(ppos_storage.CURRENT)
	next_queue := c.getCandidateQueue(ppos_storage.NEXT)

	// set previous witness
	c.setCandidateQueue(curr_queue, ppos_storage.PREVIOUS)

	// set current witness
	c.setCandidateQueue(next_queue, ppos_storage.CURRENT)

	// set next witness
	c.delCandidateQueue(ppos_storage.NEXT)

	log.Info("Call Switch SUCCESS !!!!!!!")
	return true
}

// Getting nodes of witnesses
// flag：
// -1: the previous round of witnesses
// 0: the current round of witnesses
// 1: the next round of witnesses
func (c *CandidatePool) GetWitness(state *state.StateDB, flag int, blockNumber *big.Int) ([]*discover.Node, error) {
	log.Debug("Call GetWitness", "blockNumber", blockNumber.String(), "flag", strconv.Itoa(flag))
	c.initDataByState(state)
	var queue types.CandidateQueue

	if flag == PREVIOUS_C {
		queue = c.getCandidateQueue(ppos_storage.PREVIOUS)
	} else if flag == CURRENT_C {
		queue = c.getCandidateQueue(ppos_storage.CURRENT)
	} else if flag == NEXT_C {
		queue = c.getCandidateQueue(ppos_storage.NEXT)
	}

	arr := make([]*discover.Node, 0)
	for _, can := range queue {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build Node on GetWitness", "nodeId", can.CandidateId.String(), "err", err)
			return nil, err
		} else {
			arr = append(arr, node)
		}
	}
	return arr, nil
}

// Getting previous and current and next witnesses
func (c *CandidatePool) GetAllWitness(state *state.StateDB, blockNumber *big.Int) ([]*discover.Node, []*discover.Node, []*discover.Node, error) {
	log.Debug("Call GetAllWitness ...", "blockNumber", blockNumber.String())
	c.initDataByState(state)
	loadFunc := func(title string, flag int) ([]*discover.Node, error) {
		queue := c.getCandidateQueue(flag)
		arr := make([]*discover.Node, 0)
		for _, can := range queue {
			if node, err := buildWitnessNode(can); nil != err {
				log.Error("Failed to build Node on Get "+title+" Witness", "nodeId", can.CandidateId.String(), "err", err)
				return nil, err
			} else {
				arr = append(arr, node)
			}
		}
		return arr, nil
	}

	preArr, curArr, nextArr := make([]*discover.Node, 0), make([]*discover.Node, 0), make([]*discover.Node, 0)

	type result struct {
		Type  int // -1: previous; 0: current; 1: next
		Err   error
		nodes []*discover.Node
	}
	var wg sync.WaitGroup
	wg.Add(3)
	resCh := make(chan *result, 3)

	go func() {
		res := new(result)
		res.Type = PREVIOUS_C
		if nodes, err := loadFunc("previous", ppos_storage.PREVIOUS); nil != err {
			res.Err = err
		} else {
			res.nodes = nodes
		}
		resCh <- res
		wg.Done()
	}()
	go func() {
		res := new(result)
		res.Type = CURRENT_C
		if nodes, err := loadFunc("current", ppos_storage.CURRENT); nil != err {
			res.Err = err
		} else {
			res.nodes = nodes
		}
		resCh <- res
		wg.Done()
	}()
	go func() {
		res := new(result)
		res.Type = NEXT_C
		if nodes, err := loadFunc("next", ppos_storage.NEXT); nil != err {
			res.Err = err
		} else {
			res.nodes = nodes
		}
		resCh <- res
		wg.Done()
	}()
	wg.Wait()
	close(resCh)
	for res := range resCh {
		if nil != res.Err {
			return nil, nil, nil, res.Err
		}
		switch res.Type {
		case PREVIOUS_C:
			preArr = res.nodes
		case CURRENT_C:
			curArr = res.nodes
		case NEXT_C:
			nextArr = res.nodes
		default:
			continue
		}
	}
	return preArr, curArr, nextArr, nil
}

// Getting can by witnesses
// flag:
// -1: 		previous round
// 0:		current round
// 1: 		next round
func (c *CandidatePool) GetWitnessCandidate(state vm.StateDB, nodeId discover.NodeID, flag int, blockNumber *big.Int) *types.Candidate {

	log.Debug("Call GetWitnessCandidate", "blockNumber", blockNumber.String(), "nodeId", nodeId.String(), "flag", flag)

	c.initDataByState(state)
	switch flag {
	case PREVIOUS_C:

		for _, can := range c.getCandidateQueue(ppos_storage.PREVIOUS) {
			if can.CandidateId == nodeId {
				return can
			}
		}
		log.Warn("Call GetWitnessCandidate, can no exist in previous witnesses ", "nodeId", nodeId.String())
		return nil

	case CURRENT_C:

		for _, can := range c.getCandidateQueue(ppos_storage.CURRENT) {
			if can.CandidateId == nodeId {
				return can
			}
		}
		log.Warn("Call GetWitnessCandidate, can no exist in current witnesses ", "nodeId", nodeId.String())
		return nil

	case NEXT_C:

		for _, can := range c.getCandidateQueue(ppos_storage.NEXT) {
			if can.CandidateId == nodeId {
				return can
			}
		}
		log.Warn("Call GetWitnessCandidate, can no exist in next witnesses ", "nodeId", nodeId.String())
		return nil

	default:
		log.Error("Failed to found can on GetWitnessCandidate, flag is invalid", "flag", flag)
		return nil
	}
}

func (c *CandidatePool) GetRefundInterval(blockNumber *big.Int) uint32 {
	log.Info("Call GetRefundInterval", "blockNumber", blockNumber.String(), "RefundBlockNumber", c.refundBlockNumber)
	return c.refundBlockNumber
}

// According to the nodeId to ensure the current candidate's stay
func (c *CandidatePool) UpdateElectedQueue(state vm.StateDB, currBlockNumber *big.Int, nodeIds ...discover.NodeID) error {
	log.Info("Call UpdateElectedQueue start ...", "threshold", c.threshold.String(), "depositLimit", c.depositLimit, "allowed", c.allowed, "maxCount", c.maxCount, "maxChair", c.maxChair, "refundBlockNumber", c.refundBlockNumber)
	arr := c.updateQueue(state, currBlockNumber, nodeIds...)
	log.Info("Call UpdateElectedQueue SUCCESS !!!!!!!!! ")
	//go ticketPool.DropReturnTicket(state, ids...)
	if len(arr) > 0 {
		return tContext.DropReturnTicket(state, currBlockNumber, arr...)
	}
	return nil
}

func (c *CandidatePool) updateQueue(state vm.StateDB, currentBlockNumber *big.Int, nodeIds ...discover.NodeID) []discover.NodeID {
	log.Info("Call UpdateElectedQueue Update the Campaign queue Start ...")
	PrintObject("Call UpdateElectedQueue input param's nodeIds len:= " + fmt.Sprint(len(nodeIds)) + " ,is:", nodeIds)
	if len(nodeIds) == 0 {
		log.Debug("UpdateElectedQueue FINISH, input param's nodeIds is empty !!!!!!!!!!")
		return nil
	}

	c.initData2Cache(state, GET_IM_RE)


	/**
	delete can by queue
	 */
	delCanFromQueueFunc := func(title string, nodeId discover.NodeID, queue types.CandidateQueue, flag int)  types.CandidateQueue {


		queueCopy := make(types.CandidateQueue, len(queue))
		copy(queueCopy, queue)

		log.Debug("Call UpdateElectedQueue, Before delete the can by old queue, the queue is " + title, "queue len", len(queueCopy), "nodeId", nodeId.String())

		for i, can := range queueCopy {
			if nodeId == can.CandidateId {
				queueCopy = append(queueCopy[:i], queueCopy[i+1:]...)
				break
			}
		}

		PrintObject("Call UpdateElectedQueue, After delete the can by old queue, the queue is " + title + ", queue len:" + fmt.Sprint(len(queueCopy)) + " ,the remain queue is", queueCopy)

		if len(queueCopy) != 0 {
			c.setCandidateQueue(queueCopy, flag)
		}else {
			c.delCandidateQueue(flag)
		}

		return queueCopy
	}


	/**
	handler the Reserve queue
	 */
	handleReserveFunc := func(re_queue types.CandidateQueue) []discover.NodeID {

		queueCopy := make(types.CandidateQueue, len(re_queue))
		copy(queueCopy, re_queue)

		PrintObject("Call UpdateElectedQueue, handleReserveFunc, Before update the reserve len is " + fmt.Sprint(len(queueCopy)) + ", config.maxCount:" + fmt.Sprint(c.maxCount) + " , the queue is", queueCopy)

		str := "Call UpdateElectedQueue, handleReserveFunc, o sort reserve queue ..."
		makeCandidateSort(str, state, queueCopy)
		if len(queueCopy) > int(c.maxCount) {
			// Intercepting the lost candidates to tmpArr
			tmpArr := (queueCopy)[c.maxCount:]
			// qualified elected candidates
			queueCopy = (queueCopy)[:c.maxCount]

			// cache
			nodeIdQueue := make([]discover.NodeID, 0)

			// handle tmpArr
			for _, tmpCan := range tmpArr {
				deposit, _ := new(big.Int).SetString(tmpCan.Deposit.String(), 10)
				refund := &types.CandidateRefund{
					Deposit:     deposit,
					BlockNumber: big.NewInt(currentBlockNumber.Int64()),
					Owner:       tmpCan.Owner,
				}

				c.setRefund(tmpCan.CandidateId, refund)
				nodeIdQueue = append(nodeIdQueue, tmpCan.CandidateId)

				delete(c.reserveCandidates, tmpCan.CandidateId)

			}



			PrintObject("Call UpdateElectedQueue, handleReserveFunc, After update the reserve len is " + fmt.Sprint(len(queueCopy)) + " , the queue is", queueCopy)

			c.setCandidateQueue(queueCopy, ppos_storage.RESERVE)

			return nodeIdQueue
		}else {

			PrintObject("Call UpdateElectedQueue, handleReserveFunc, After update the reserve len is " + fmt.Sprint(len(queueCopy)) + " , the queue is", queueCopy)

			c.setCandidateQueue(queueCopy, ppos_storage.RESERVE)
			return nil
		}
	}

	/**
	This function handles Immediate queues and Reserve queues
	for moving into the opposing queue
	*/
	workFunc := func(oldQueueFlag, newQueueFlag int, can *types.Candidate) []discover.NodeID {
		old_queue := c.getCandidateQueue(oldQueueFlag)
		new_queue := c.getCandidateQueue(newQueueFlag)

		//old_queueCopy := make(types.CandidateQueue, len(old_queue))
		//new_queueCopy := make(types.CandidateQueue, len(new_queue))

		nodeIdQueue := make([]discover.NodeID, 0)

		/**
		immediate move to reserve
		*/
		if oldQueueFlag == ppos_storage.IMMEDIATE {

			log.Debug("Call UpdateElectedQueue, workFunc the old queue is immediate", "nodeId", can.CandidateId.String())

			// delete immediate
			delCanFromQueueFunc("imms", can.CandidateId, old_queue, oldQueueFlag)
			delete(c.immediateCandidates, can.CandidateId)


			// input reserve
			new_queue = append(new_queue, can)
			if nodeIdArr := handleReserveFunc(new_queue); len(nodeIdArr) != 0 {
				nodeIdQueue = append(nodeIdQueue, nodeIdArr...)
			}
			return nodeIdQueue
		} else {

			log.Debug("Call UpdateElectedQueue, workFunc the old queue is reserve", "nodeId", can.CandidateId.String())

			/**
			reserve move to immediate
			*/
			// delete reserve
			old_queue = delCanFromQueueFunc("res", can.CandidateId, old_queue, oldQueueFlag)
			delete(c.reserveCandidates, can.CandidateId)

			str := "Call UpdateElectedQueue, workFunc to sort immediate queue ..."
			// input immediate
			new_queue = append(new_queue, can)
			makeCandidateSort(str, state, new_queue)

			var inRes bool
			if len(new_queue) > int(c.maxCount) {
				// Intercepting the lost candidates to tmpArr
				tmpArr := (new_queue)[c.maxCount:]
				// qualified elected candidates
				new_queue = (new_queue)[:c.maxCount]

				inRes = true
				// reenter into reserve
				old_queue = append(old_queue, tmpArr...)
			}

			log.Debug("Call UpdateElectedQueue, workFunc the old queue is reserve, the new immediate", "len", len(new_queue), "config.maxCount", c.maxCount)
			PrintObject("Call UpdateElectedQueue, workFunc the old queue is reserve, the new immediate", new_queue)

			// update immediate
			c.setCandidateQueue(new_queue, ppos_storage.IMMEDIATE)
			if inRes {

				log.Debug("Call UpdateElectedQueue, workFunc the old queue is reserve, the new immediate, had add reserve", "reserver len", len(old_queue))

				if nodeIdArr := handleReserveFunc(old_queue); len(nodeIdArr) != 0 {
					nodeIdQueue = append(nodeIdQueue, nodeIdArr...)
				}
			}
			return nodeIdQueue
		}

	}


	resNodeIds := make([]discover.NodeID, 0)

	/**
	########
	Real Handle Can queue
	########
	*/
	for _, nodeId := range nodeIds {

		tcount := c.checkTicket(tContext.GetCandidateTicketCount(state, nodeId))

		switch c.checkExist(nodeId) {

		case IS_IMMEDIATE:
			log.Debug("Call UpdateElectedQueue The current nodeId was originally in the Immediate ...")
			can := c.immediateCandidates[nodeId]

			if !tcount {

				log.Debug("Call UpdateElectedQueue, the node will from immediate to reserve ...", "nodeId", nodeId.String())
				if ids := workFunc(ppos_storage.IMMEDIATE, ppos_storage.RESERVE, can); len(ids) != 0 {
					resNodeIds = append(resNodeIds, ids...)
				}
			}

		case IS_RESERVE:
			log.Debug("Call UpdateElectedQueue The current nodeId was originally in the Reserve ...")
			can := c.reserveCandidates[nodeId]

			if tcount {

				log.Debug("Call UpdateElectedQueue, the node will from reserve to immediate ...", "nodeId", nodeId.String())
				if ids := workFunc(ppos_storage.RESERVE, ppos_storage.IMMEDIATE, can); len(ids) != 0 {
					resNodeIds = append(resNodeIds, ids...)
				}
			}

		default:
			continue
		}
	}

	// promoteReserve queues
	c.promoteReserveQueue(state, currentBlockNumber)

	return resNodeIds
}

func (c *CandidatePool) checkFirstThreshold(can *types.Candidate) bool {
	var exist bool
	if _, ok := c.immediateCandidates[can.CandidateId]; ok {
		exist = true
	}

	if _, ok := c.reserveCandidates[can.CandidateId]; ok {
		exist = true
	}

	if !exist && can.Deposit.Cmp(c.threshold) < 0 {
		return false
	}
	return true
}

// false: invalid deposit
// true:  pass
//func (c *CandidatePool) checkDeposit(state vm.StateDB, can *types.Candidate, holdself bool) (bool, bool) {
//
//
//
//
//
//	// Number of vote ticket by nodes
//	tcount := c.checkTicket(tContext.GetCandidateTicketCount(state, can.CandidateId))
//
//	// if self have already exist:
//	// b、no first pledge: x >= self * 110%
//	//
//	// if the pool is full:(Only reserve pool)
//	// c、x > last * 110 %
//
//
//	compareFunc := func(target, current *types.Candidate, logA, logB string) (bool, bool) {
//		lastDeposit := target.Deposit
//
//		// y = 100 + x
//		percentage := new(big.Int).Add(big.NewInt(100), big.NewInt(int64(c.depositLimit)))
//		// z = old * y
//		tmp := new(big.Int).Mul(lastDeposit, percentage)
//		// z/100 == old * (100 + x) / 100 == old * (y%)
//		tmp = new(big.Int).Div(tmp, big.NewInt(100))
//		if can.Deposit.Cmp(tmp) < 0 {
//			// If last is self and holdslef flag is true
//			// we must return true (Keep self on staying in the original queue)
//			if holdself && can.CandidateId == target.CandidateId {
//				log.Debug(logA + tmp.String())
//				return tcount, true
//			}
//			log.Debug(logB + tmp.String())
//			return tcount, false
//		}
//		return tcount, true
//	}
//
//
//	//var queueFlag int
//	var storageMap candidateStorage
//
//	if _, ok := c.immediateCandidates[can.CandidateId]; ok {
//		//queueFlag = IS_IMMEDIATE
//		storageMap = c.immediateCandidates
//	}
//	if _, ok := c.reserveCandidates[can.CandidateId]; ok {
//		//queueFlag = IS_RESERVE
//		storageMap = c.reserveCandidates
//	}
//
//	/*
//	If the current number of votes achieving the conditions for entering the immediate queue
//	*/
//	if tcount {
//
//		// if it already exist immediate
//		// compare old self deposit
//		if old_self, ok := storageMap[can.CandidateId]; ok {
//
//			logA := `The can already exist in queue, and holdslef is true, Keep self on staying in the original queue, current can nodeId:` + can.CandidateId.String() +
//				` the length of current queue:` + fmt.Sprint(len(storageMap)) + ` the limit of Configuration:` + fmt.Sprint(c.maxCount) +
//				` current can's Deposit: ` + can.Deposit.String() + ` 110% of it old_self in the queue: `
//
//			logB := `The can already exist in queue, and holdslef is true,and the current can's Deposit is less than 110% of the it old self in the queue.` +
//				`the length of current queue:` + fmt.Sprint(len(storageMap)) + ` the limit of Configuration:` + fmt.Sprint(c.maxCount) + ` current can's Deposit:` +
//				can.Deposit.String() + ` 110% of it old_self in the queue: `
//
//			return compareFunc(old_self, can, logA, logB)
//		}
//
//	}
//
//	/**
//	If the can will enter the reserve pool
//	*/
//
//	if !tcount {
//
//		reserveArr := c.storage.GetCandidateQueue(ppos_storage.RESERVE)
//
//		// Is first pledge and the reserve pool is full
//		if _, ok := storageMap[can.CandidateId]; !ok && uint32(len(reserveArr)) == c.maxCount {
//			last := reserveArr[len(reserveArr)-1]
//
//			logA := `The reserve pool is full, and last is self and holdslef is true, Keep self on staying in the original queue, the length of current reserve pool: ` +
//				fmt.Sprint(len(reserveArr)) + ` the limit of Configuration:` + fmt.Sprint(c.maxCount) + ` current can nodeId: ` + can.CandidateId.String() + ` last can nodeId: ` +
//				last.CandidateId.String() + ` current can's Deposit: ` + can.Deposit.String() + ` 110% of the last can in the reserve pool:`
//
//			logB := `The reserve pool is full,and the current can's Deposit is less than 110% of the last can in the reserve pool., the length of current reserve pool:` +
//				fmt.Sprint(len(reserveArr)) + ` the limit of Configuration:` + fmt.Sprint(c.maxCount) + ` current can's Deposit:` + can.Deposit.String() +
//				`110% of the last can in the reserve pool:`
//
//			return compareFunc(last, can, logA, logB)
//
//		} else if old_self, ok := storageMap[can.CandidateId]; ok { // Is'nt first pledge
//			logA := `The can already exist in reserve, and holdslef is true, Keep self on staying in the original queue, current can nodeId:` + can.CandidateId.String() +
//				` the length of current reserve pool:` + fmt.Sprint(len(reserveArr)) + ` the limit of Configuration:` + fmt.Sprint(c.maxCount) +
//				` current can's Deposit: ` + can.Deposit.String() + ` 110% of it old_self in the reserve pool: `
//
//			logB := `The can already exist in reserve, and holdslef is true,and the current can's Deposit is less than 110% of the it old self in the reserve pool.` +
//				`the length of current reserve pool:` + fmt.Sprint(len(reserveArr)) + ` the limit of Configuration:` + fmt.Sprint(c.maxCount) + ` current can's Deposit:` +
//				can.Deposit.String() + ` 110% of it old_self in the reserve pool: `
//
//			return compareFunc(old_self, can, logA, logB)
//		}
//	}
//	return tcount, true
//}

func (c *CandidatePool) checkDeposit (can *types.Candidate) bool {

	im_queue := c.getCandidateQueue(ppos_storage.IMMEDIATE)

	re_queue := c.getCandidateQueue(ppos_storage.RESERVE)

	if len(im_queue) == int(c.maxCount) && len(re_queue) == int(c.maxCount) {
		last := re_queue[len(re_queue)-1]

		lastDeposit := last.Deposit

		// y = 100 + x
		percentage := new(big.Int).Add(big.NewInt(100), big.NewInt(int64(c.depositLimit)))
		// z = old * y
		tmp := new(big.Int).Mul(lastDeposit, percentage)
		// z/100 == old * (100 + x) / 100 == old * (y%)
		tmp = new(big.Int).Div(tmp, big.NewInt(100))
		if can.Deposit.Cmp(tmp) < 0 {
			log.Debug("The current can's Deposit is less than 110% of the last can in the reserve queue.", "current can Deposit", can.Deposit.String(), "last can Deposit", last.Deposit.String(), "the target Deposit", tmp.String())
			return false
		}
	}
	return true
}


func (c *CandidatePool) checkWithdraw(source, price *big.Int) error {
	// y = old * x
	percentage := new(big.Int).Mul(source, big.NewInt(int64(c.depositLimit)))
	// y/100 == old * (x/100) == old * x%
	tmp := new(big.Int).Div(percentage, big.NewInt(100))
	if price.Cmp(tmp) < 0 {
		log.Debug("When withdrawing the refund, the amount of the current can't be refunded is less than 10% of the remaining amount of self:", "The current amount of can want to refund:", price.String(), "current Can own deposit remaining amount:", tmp.String())
		return WithdrawLowErr
	}
	return nil
}

// 0: empty
// 1: in immediates
// 2: in reserves
func (c *CandidatePool) checkExist(nodeId discover.NodeID) int {
	log.Info("Check which queue the current nodeId originally belongs to:", "nodeId", nodeId.String())
	if _, ok := c.immediateCandidates[nodeId]; ok {
		log.Info("The current nodeId originally belonged to the immediate queue ...", "nodeId", nodeId.String())
		return IS_IMMEDIATE
	}
	if _, ok := c.reserveCandidates[nodeId]; ok {
		log.Info("The current nodeId originally belonged to the reserve queue ...", "nodeId", nodeId.String())
		return IS_RESERVE
	}
	log.Info("The current nodeId does not belong to any queue ...", "nodeId", nodeId.String())
	return IS_LOST
}

func (c *CandidatePool) checkTicket(t_count uint32) bool {
	log.Debug("Compare the current candidate’s votes to:", "t_count", t_count, "the allowed limit of config:", c.allowed)
	if t_count >= c.allowed {
		log.Debug(" The current candidate’s votes are in line with the immediate pool....")
		return true
	}
	log.Debug("Not eligible to enter the immediate pool ...")
	return false
}

func (c *CandidatePool) getCandidate(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) *types.Candidate {
	log.Debug("Call GetCandidate", "blockNumber", blockNumber.String(), "nodeId", nodeId.String())
	c.initData2Cache(state, GET_IM_RE)
	if candidatePtr, ok := c.immediateCandidates[nodeId]; ok {
		PrintObject("Call GetCandidate return immediate：", *candidatePtr)
		return candidatePtr
	}
	if candidatePtr, ok := c.reserveCandidates[nodeId]; ok {
		PrintObject("Call GetCandidate return reserve：", *candidatePtr)
		return candidatePtr
	}
	return nil
}

func (c *CandidatePool) getCandidates(state vm.StateDB, blockNumber *big.Int, nodeIds ...discover.NodeID) types.CandidateQueue {
	log.Debug("Call GetCandidates...", "blockNumber", blockNumber.String())
	c.initData2Cache(state, GET_IM_RE)
	canArr := make(types.CandidateQueue, 0)
	tem := make(map[discover.NodeID]struct{}, 0)
	for _, nodeId := range nodeIds {
		if _, ok := tem[nodeId]; ok {
			continue
		}
		if candidatePtr, ok := c.immediateCandidates[nodeId]; ok {
			canArr = append(canArr, candidatePtr)
			tem[nodeId] = struct{}{}
		}
		if _, ok := tem[nodeId]; ok {
			continue
		}
		if candidatePtr, ok := c.reserveCandidates[nodeId]; ok {
			canArr = append(canArr, candidatePtr)
			tem[nodeId] = struct{}{}
		}
	}
	return canArr
}

func (c *CandidatePool) MaxChair() uint32 {
	return c.maxChair
}

func (c *CandidatePool) MaxCount() uint32 {
	return c.maxCount
}

// TODO
func (c *CandidatePool) promoteReserveQueue(state vm.StateDB, currentBlockNumber *big.Int) /*[]discover.NodeID */ {

	log.Debug("Call promoteReserveQueue Start ...", "blockNumber", currentBlockNumber.String())

	// Violence traverses the pools
	im_queue := c.storage.GetCandidateQueue(ppos_storage.IMMEDIATE)
	re_queue := c.storage.GetCandidateQueue(ppos_storage.RESERVE)

	PrintObject("Call promoteReserveQueue Before the immediate queue len:=" + fmt.Sprint(len(im_queue)) + " ,queue is", im_queue)
	PrintObject("Call promoteReserveQueue Before the reserve queue len:=" + fmt.Sprint(len(re_queue)) + " ,queue is", re_queue)


	// current ticket price
	ticket_price := tContext.GetTicketPrice(state)


	//for i, can := range re_queue {
	for i := 0; i < len(re_queue); i++ {

		// current re can
		can := re_queue[i]

		log.Debug("Call promoteReserveQueue for range , Current reserve id", "blockNumber", currentBlockNumber.String(), "nodeId", can.CandidateId.String())

		re_tCount := tContext.GetCandidateTicketCount(state, can.CandidateId)
		if re_tCount < c.allowed {
			continue
		}

		// check immediate
		// if the immediate queue is full
		// We can take the immediate last can and current reserve can compare
		if uint32(len(im_queue)) >= c.maxCount {

			// last im can
			last := im_queue[len(im_queue)-1]
			last_tCount := tContext.GetCandidateTicketCount(state, last.CandidateId)

			last_tmoney := new(big.Int).Mul(big.NewInt(int64(last_tCount)), ticket_price)
			last_money := new(big.Int).Add(last.Deposit, last_tmoney)

			// current re can
			re_tmoney := new(big.Int).Mul(big.NewInt(int64(re_tCount)), ticket_price)
			re_money := new(big.Int).Add(can.Deposit, re_tmoney)

			// Terminate the cycle comparison.
			// When current re can's money is
			// less than the last one in the im queue
			if types.CompareCan(can, last, re_money, last_money) < 0 {
				break
			}


			// into immediate queue If match condition
			im_queue = append(im_queue, can)
			re_queue = append(re_queue[:i], re_queue[i+1:]...)
			i--
		} else {
			// direct - into immediate queue
			im_queue = append(im_queue, can)
			re_queue = append(re_queue[:i], re_queue[i+1:]...)
			i--
		}
	}


	str := "Call promoteReserveQueue to sort the new immediate queue ..."
	// sort immediate
	makeCandidateSort(str, state, im_queue)

	//nodeIds := make([]discover.NodeID, 0)

	var addRe_queue types.CandidateQueue


	// Cut off the immediate queue
	if len(im_queue) > int(c.maxCount) {
		// Intercepting the lost candidates to tmpArr
		tempArr := (im_queue)[c.maxCount:]
		// qualified elected candidates
		im_queue = (im_queue)[:c.maxCount]




		//newRe_queue = make(types.CandidateQueue, 0)
		addRe_queue = make(types.CandidateQueue, len(tempArr))

		// handle tmpArr
		for i, tmpCan := range tempArr {

			addRe_queue[i] = tmpCan
		}
	}

	// Sets the new immediate queue
	c.setCandidateQueue(im_queue, ppos_storage.IMMEDIATE)

	// sort new reserve queue
	if len(addRe_queue) > 0 {
		re_queue = append(re_queue, addRe_queue...)

		str := "Call promoteReserveQueue to sort the new reserve queue ..."
		makeCandidateSort(str, state, re_queue)

	}

	// Sets the new reserve queue
	c.setCandidateQueue(re_queue, ppos_storage.RESERVE)

	PrintObject("Call promoteReserveQueue Finish the immediate queue len:=" + fmt.Sprint(len(im_queue)) + " ,queue is", im_queue)
	PrintObject("Call promoteReserveQueue Finish the reserve queue len:=" + fmt.Sprint(len(re_queue)) + " ,queue is", re_queue)
	log.Debug("Call promoteReserveQueue Finish !!! ...", "blockNumber", currentBlockNumber.String())
}

/** builin function */

func (c *CandidatePool) setCandidateQueue(queue types.CandidateQueue, flag int) {
	c.storage.SetCandidateQueue(queue, flag)
}

func (c *CandidatePool) getCandidateQueue(flag int) types.CandidateQueue {
	return c.storage.GetCandidateQueue(flag)
}

func (c *CandidatePool) delCandidateQueue(flag int) {
	c.storage.DelCandidateQueue(flag)
}

func (c *CandidatePool) setRefund(nodeId discover.NodeID, refund *types.CandidateRefund) {
	c.storage.SetRefund(nodeId, refund)
}

func (c *CandidatePool) setRefunds(nodeId discover.NodeID, refundArr types.RefundQueue) {
	c.storage.SetRefunds(nodeId, refundArr)
}

func (c *CandidatePool) appendRefunds(nodeId discover.NodeID, refundArr types.RefundQueue) {
	c.storage.AppendRefunds(nodeId, refundArr)
}

func (c *CandidatePool) getRefunds(nodeId discover.NodeID) types.RefundQueue {
	return c.storage.GetRefunds(nodeId)
}

func (c *CandidatePool) delRefunds(nodeId discover.NodeID) {
	c.storage.DelRefunds(nodeId)
}

/*func getPreviousWitnessIdsState(state vm.StateDB) ([]discover.NodeID, error) {
	var witnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, PreviousWitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &witnessIds); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return witnessIds, nil
}

func setPreviosWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, PreviousWitnessListKey(), arrVal)
}

func getPreviousWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, PreviousWitnessKey(id)); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &can, nil
}

func setPreviousWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, PreviousWitnessKey(id), val)
}

func getWitnessIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var witnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, WitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &witnessIds); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return witnessIds, nil
}

func setWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, WitnessListKey(), arrVal)
}

func getWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, WitnessKey(id)); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &can, nil
}

func setWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, WitnessKey(id), val)
}

func getNextWitnessIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var nextWitnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, NextWitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &nextWitnessIds); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return nextWitnessIds, nil
}

func setNextWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, NextWitnessListKey(), arrVal)
}

func getNextWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, NextWitnessKey(id)); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &can, nil
}

func setNextWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, NextWitnessKey(id), val)
}

func getImmediateIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var immediateIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, ImmediateListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &immediateIds); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return immediateIds, nil
}

func setImmediateIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, ImmediateListKey(), arrVal)
}

func getImmediateByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, ImmediateKey(id)); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &can, nil
}

func setImmediateState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, ImmediateKey(id), val)
}

func getReserveIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var reserveIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, ReserveListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &reserveIds); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return reserveIds, nil
}

func setReserveIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, ReserveListKey(), arrVal)
}

func getReserveByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, ReserveKey(id)); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return &can, nil
}

func setReserveState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, ReserveKey(id), val)
}

func getDefeatIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var defeatIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, DefeatListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &defeatIds); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return defeatIds, nil
}

func setDefeatIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, DefeatListKey(), arrVal)
}

func getDefeatsByState(state vm.StateDB, id discover.NodeID) (types.CandidateQueue, error) {
	var canArr types.CandidateQueue
	if valByte := state.GetState(common.CandidatePoolAddr, DefeatKey(id)); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &canArr); nil != err {
			return nil, err
		}
	} else {
		return nil, nil
	}
	return canArr, nil
}

func setDefeatState(state vm.StateDB, id discover.NodeID, val []byte) {
	log.Debug("SetDefeatArr ... ...", "nodeId:", id.String(), "keyTrie:", buildKeyTrie(DefeatKey(id)))
	state.SetState(common.CandidatePoolAddr, DefeatKey(id), val)
}*/

func copyCandidateMapByIds(target, source candidateStorage, ids []discover.NodeID) {
	for _, id := range ids {
		if v, ok := source[id]; ok {
			target[id] = v
		}
	}
}

//func GetCandidatePtr() *CandidatePool {
//	return candidatePool
//}

func PrintObject(s string, obj interface{}) {
	objs, _ := json.Marshal(obj)
	log.Debug(s, "==", string(objs))
}

func buildWitnessNode(can *types.Candidate) (*discover.Node, error) {
	if nil == can {
		return nil, CandidateEmptyErr
	}
	ip := net.ParseIP(can.Host)
	// uint16
	var port uint16
	if portInt, err := strconv.Atoi(can.Port); nil != err {
		return nil, err
	} else {
		port = uint16(portInt)
	}
	return discover.NewNode(can.CandidateId, ip, port, port), nil
}

func makeCandidateSort(logStr string, state vm.StateDB, arr types.CandidateQueue) {

	log.Debug(logStr)

	cand := make(types.CanConditions, 0)
	for _, can := range arr {
		tCount := tContext.GetCandidateTicketCount(state, can.CandidateId)
		price := tContext.GetTicketPrice(state)
		tprice := new(big.Int).Mul(big.NewInt(int64(tCount)), price)

		money := new(big.Int).Add(can.Deposit, tprice)

		cand[can.CandidateId] = money
	}
	arr.CandidateSort(cand)
}

/*
func ImmediateKey(nodeId discover.NodeID) []byte {
	return immediateKey(nodeId.Bytes())
}
func immediateKey(key []byte) []byte {
	return append(append(common.CandidatePoolAddr.Bytes(), ImmediateBytePrefix...), key...)
}

func ReserveKey(nodeId discover.NodeID) []byte {
	return reserveKey(nodeId.Bytes())
}

func reserveKey(key []byte) []byte {
	return append(append(common.CandidatePoolAddr.Bytes(), ReserveBytePrefix...), key...)
}

func PreviousWitnessKey(nodeId discover.NodeID) []byte {
	return prewitnessKey(nodeId.Bytes())
}

func prewitnessKey(key []byte) []byte {
	return append(append(common.CandidatePoolAddr.Bytes(), PreWitnessBytePrefix...), key...)
}

func WitnessKey(nodeId discover.NodeID) []byte {
	return witnessKey(nodeId.Bytes())
}
func witnessKey(key []byte) []byte {
	return append(append(common.CandidatePoolAddr.Bytes(), WitnessBytePrefix...), key...)
}

func NextWitnessKey(nodeId discover.NodeID) []byte {
	return nextWitnessKey(nodeId.Bytes())
}
func nextWitnessKey(key []byte) []byte {
	return append(append(common.CandidatePoolAddr.Bytes(), NextWitnessBytePrefix...), key...)
}

func DefeatKey(nodeId discover.NodeID) []byte {
	return defeatKey(nodeId.Bytes())
}
func defeatKey(key []byte) []byte {
	return append(append(common.CandidatePoolAddr.Bytes(), DefeatBytePrefix...), key...)
}

func ImmediateListKey() []byte {
	return append(common.CandidatePoolAddr.Bytes(), ImmediateListBytePrefix...)
}

func ReserveListKey() []byte {
	return append(common.CandidatePoolAddr.Bytes(), ReserveListBytePrefix...)
}

func PreviousWitnessListKey() []byte {
	return append(common.CandidatePoolAddr.Bytes(), PreWitnessListBytePrefix...)
}

func WitnessListKey() []byte {
	return append(common.CandidatePoolAddr.Bytes(), WitnessListBytePrefix...)
}

func NextWitnessListKey() []byte {
	return append(common.CandidatePoolAddr.Bytes(), NextWitnessListBytePrefix...)
}

func DefeatListKey() []byte {
	return append(common.CandidatePoolAddr.Bytes(), DefeatListBytePrefix...)
}

func (c *CandidatePool) ForEachStorage(state vm.StateDB, title string) {
	c.lock.Lock()
	log.Debug(title + ":Full view of data in the candidate pool ...")
	c.initDataByState(state, 2)
	c.lock.Unlock()
}
*/
