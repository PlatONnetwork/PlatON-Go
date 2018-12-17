package pposm

import (
	"Platon-go/common"
	"Platon-go/core/state"
	"Platon-go/core/types"
	"Platon-go/core/vm"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"Platon-go/rlp"
	"encoding/json"
	"errors"
	_ "fmt"
	"math/big"
	"net"
	"strconv"
	"sync"
)

var (
	CandidateEncodeErr          = errors.New("Candidate encoding err")
	CandidateDecodeErr          = errors.New("Candidate decoding err")
	WithdrawPriceErr            = errors.New("Withdraw Price err")
	CandidateEmptyErr           = errors.New("Candidate is empty")
	ContractBalanceNotEnoughErr = errors.New("Contract's balance is not enough")
	CandidateOwnerErr           = errors.New("CandidateOwner Addr is illegal")
)

type candidateStorage map[discover.NodeID]*types.Candidate
type refundStorage map[discover.NodeID][]*types.Candidate

//type returnTicket struct {
//	NodeId 		discover.NodeID
//	TicketId	common.Hash
//}

type CandidatePool struct {
	// allow immediate elected max count
	maxCount uint64
	// allow witness max count
	maxChair uint64
	// allow block interval for refunds
	RefundBlockNumber uint64

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
	immediateCacheArr []*types.Candidate
	reserveCacheArr   []*types.Candidate
	lock              *sync.RWMutex
}

var candidatePool *CandidatePool

// Initialize the global candidate pool object
func NewCandidatePool(configs *params.PposConfig) *CandidatePool {
	candidatePool = &CandidatePool{
		maxCount:             configs.Candidate.MaxCount,
		maxChair:             configs.Candidate.MaxChair,
		RefundBlockNumber:    configs.Candidate.RefundBlockNumber,
		preOriginCandidates:  make(candidateStorage, 0),
		originCandidates:     make(candidateStorage, 0),
		nextOriginCandidates: make(candidateStorage, 0),
		immediateCandidates:  make(candidateStorage, 0),
		defeatCandidates:     make(refundStorage, 0),
		immediateCacheArr:    make([]*types.Candidate, 0),
		reserveCacheArr:      make([]*types.Candidate, 0),
		lock:                 &sync.RWMutex{},
	}
	return candidatePool
}

// flag:
// 0: only init previous witness and current witness and next witness
// 1：init previous witness and current witness and next witness and immediate
// 2: init all information
func (c *CandidatePool) initDataByState(state vm.StateDB, flag int) error {
	log.Info("init data by stateDB...")
	// loading previous witness
	var prewitnessIds []discover.NodeID
	c.preOriginCandidates = make(candidateStorage, 0)
	if ids, err := getPreviousWitnessIdsState(state); nil != err {
		log.Error("Failed to decode previous witnessIds on initDataByState", " err", err)
		return err
	} else {
		prewitnessIds = ids
	}
	PrintObject("prewitnessIds", prewitnessIds)
	for _, witnessId := range prewitnessIds {
		//fmt.Println("prewitnessId = ", witnessId.String())
		//var can *types.Candidate
		if ca, err := getPreviousWitnessByState(state, witnessId); nil != err {
			log.Error("Failed to decode previous witness Candidate on initDataByState", "err", err)
			return CandidateDecodeErr
		} else {
			if nil != ca {
				c.preOriginCandidates[witnessId] = ca
			} else {
				delete(c.preOriginCandidates, witnessId)
			}
		}
	}

	// loading current witnesses
	var witnessIds []discover.NodeID
	c.originCandidates = make(candidateStorage, 0)
	if ids, err := getWitnessIdsByState(state); nil != err {
		log.Error("Failed to decode witnessIds on initDataByState", "err", err)
		return err
	} else {
		witnessIds = ids
	}
	PrintObject("witnessIds", witnessIds)
	for _, witnessId := range witnessIds {
		//fmt.Println("witnessId = ", witnessId.String())
		//var can *types.Candidate
		if ca, err := getWitnessByState(state, witnessId); nil != err {
			log.Error("Failed to decode current witness Candidate on initDataByState", "err", err)
			return CandidateDecodeErr
		} else {
			if nil != ca {
				c.originCandidates[witnessId] = ca
			} else {
				delete(c.originCandidates, witnessId)
			}
		}
	}

	// loading next witnesses
	var nextWitnessIds []discover.NodeID
	c.nextOriginCandidates = make(candidateStorage, 0)
	if ids, err := getNextWitnessIdsByState(state); nil != err {
		log.Error("Failed to decode nextWitnessIds on initDataByState", "err", err)
		return err
	} else {
		nextWitnessIds = ids
	}
	PrintObject("nextWitnessIds", nextWitnessIds)
	for _, witnessId := range nextWitnessIds {
		//fmt.Println("nextwitnessId = ", witnessId.String())
		//var can *types.Candidate
		if ca, err := getNextWitnessByState(state, witnessId); nil != err {
			log.Error("Failed to decode next witness Candidate on initDataByState", "err", err)
			return CandidateDecodeErr
		} else {
			if nil != ca {
				c.nextOriginCandidates[witnessId] = ca
			} else {
				delete(c.nextOriginCandidates, witnessId)
			}
		}
	}

	if flag == 1 || flag == 2 {
		// loading immediate elected candidates
		var immediateIds []discover.NodeID
		c.immediateCandidates = make(candidateStorage, 0)
		if ids, err := getImmediateIdsByState(state); nil != err {
			log.Error("Failed to decode immediateIds on initDataByState", "err", err)
			return err
		} else {
			immediateIds = ids
		}

		// cache
		canCache := make([]*types.Candidate, 0)

		PrintObject("immediateIds", immediateIds)
		for _, immediateId := range immediateIds {
			//fmt.Println("immediateId = ", immediateId.String())
			//var can *types.Candidate
			if ca, err := getImmediateByState(state, immediateId); nil != err {
				log.Error("Failed to decode immediate Candidate on initDataByState", "err", err)
				return CandidateDecodeErr
			} else {
				if nil != ca {
					c.immediateCandidates[immediateId] = ca
					canCache = append(canCache, ca)
				} else {
					delete(c.immediateCandidates, immediateId)
				}
			}
		}
		c.immediateCacheArr = canCache

		// loading reserve elected candidates
		var reserveIds []discover.NodeID
		c.reserveCandidates = make(candidateStorage, 0)
		if ids, err := getReserveIdsByState(state); nil != err {
			log.Error("Failed to decode reserveIds on initDataByState", "err", err)
			return err
		} else {
			reserveIds = ids
		}

		canCache = make([]*types.Candidate, 0)

		PrintObject("reserveIds", reserveIds)
		for _, reserveId := range reserveIds {
			if ca, err := getReserveByState(state, reserveId); nil != err {
				log.Error("Failed to decode reserve Candidate on initDataByState", "err", err)
				return CandidateDecodeErr
			} else {
				if nil != ca {
					c.reserveCandidates[reserveId] = ca
					canCache = append(canCache, ca)
				} else {
					delete(c.reserveCandidates, reserveId)
				}
			}
		}
		c.reserveCacheArr = canCache
	}

	if flag == 2 {
		// load refunds
		var defeatIds []discover.NodeID
		c.defeatCandidates = make(refundStorage, 0)
		if ids, err := getDefeatIdsByState(state); nil != err {
			log.Error("Failed to decode defeatIds on initDataByState", "err", err)
			return err
		} else {
			defeatIds = ids
		}
		PrintObject("defeatIds", defeatIds)
		for _, defeatId := range defeatIds {
			//fmt.Println("defeatId = ", defeatId.String())
			//var canArr []*types.Candidate
			if arr, err := getDefeatsByState(state, defeatId); nil != err {
				log.Error("Failed to decode defeat CandidateArr on initDataByState", "err", err)
				return CandidateDecodeErr
			} else {
				if nil != arr && len(arr) != 0 {
					c.defeatCandidates[defeatId] = arr
				} else {
					delete(c.defeatCandidates, defeatId)
				}
			}
		}
	}
	return nil
}

// pledge Candidate // TODO 需要抽小方法
func (c *CandidatePool) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate, parentNumber *big.Int, parentHash common.Hash) error {
	var nodeIds []discover.NodeID
	c.lock.Lock()
	if err := c.initDataByState(state, 2); nil != err {
		c.lock.Unlock()
		log.Error("Failed to initDataByState on SetCandidate", "nodeId", nodeId.String(), " err", err)
		return err
	}
	if arr, err := c.setCandidateInfo(state, nodeId, can); nil != err {
		c.lock.Unlock()
		log.Error("Failed to setCandidateInfo on SetCandidate", "nodeId", nodeId.String(), "err", err)
		return err
	} else {
		nodeIds = arr
	}
	c.lock.Unlock()
	return ticketPool.DropReturnTicket(state, parentNumber, parentHash, nodeIds...)
}

// 还需要补上 如果TCout 小了，要先移到 reserves 中，不然才算落榜
func (c *CandidatePool) setCandidateInfo(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) ([]discover.NodeID, error) {
	PrintObject("发生质押 SetCandidate", *can)
	//c.lock.Lock()
	//defer c.lock.Unlock()
	//if err := c.initDataByState(state, 2); nil != err {
	//	log.Error("Failed to initDataByState on SetCandidate", " err", err)
	//	return nil, err
	//}

	var flag, delimmediate, delreserve bool
	// check ticket count
	if checkTicket(40/*ticketidsCache.TCount()*/, c.maxCount) {
		flag = true
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

	// cache
	cacheArr := make([]*types.Candidate, 0)
	if flag {
		for _, v := range c.immediateCandidates {
			cacheArr = append(cacheArr, v)
		}
	} else {
		for _, v := range c.reserveCandidates {
			cacheArr = append(cacheArr, v)
		}
	}

	// sort cache array
	candidateSort(cacheArr)

	// nodeIds cache for lost elected
	nodeIds := make([]discover.NodeID, 0)

	if len(cacheArr) > int(c.maxCount) {
		// Intercepting the lost candidates to tmpArr
		tmpArr := (cacheArr)[c.maxCount:]
		// qualified elected candidates
		cacheArr = (cacheArr)[:c.maxCount]

		// handle tmpArr
		for _, tmpCan := range tmpArr {
			// return ticket call
			// TODO
			//ticketPool.ReturnTicket(state, )
			if flag {
				// delete the lost candidates from immediate elected candidates of trie
				c.delImmediate(state, tmpCan.CandidateId)
			} else {
				// delete the lost candidates from reserve elected candidates of trie
				c.delReserve(state, tmpCan.CandidateId)
			}
			// append to refunds (defeat) trie
			if err := c.setDefeat(state, tmpCan.CandidateId, tmpCan); nil != err {
				return nil, err
			}
			nodeIds = append(nodeIds, tmpCan.CandidateId)
		}

		// update index of refund (defeat) on trie
		if err := c.setDefeatIndex(state); nil != err {
			return nil, err
		}
	}

	handle := func(candidateArr []*types.Candidate, setInfoFn func(state vm.StateDB, candidateId discover.NodeID, can *types.Candidate) error,
					setIndexFn func(state vm.StateDB, nodeIds []discover.NodeID) error) error {

					// cache id
					sortIds := make([]discover.NodeID, 0)

					/**  doing ... */
					candidateArr = cacheArr
					// insert elected candidate to tire
					for _, can := range candidateArr {
						if err := setInfoFn(state, can.CandidateId, can); nil != err {
							return err
						}
						sortIds = append(sortIds, can.CandidateId)
					}
					// update index of elected candidates on trie
					if err := setIndexFn(state, sortIds); nil != err {
						return err
					}
					return nil
				}

	delFunc := func (title string, delInfoFn func (state vm.StateDB, candidateId discover.NodeID),
					getIndexFn func (state vm.StateDB) ([]discover.NodeID, error),
					setIndexFn func (state vm.StateDB, nodeIds []discover.NodeID) error) error{

					delInfoFn(state, can.CandidateId)
					// update reserve id index
					if ids, err := getIndexFn(state); nil != err {
						log.Error("withdraw failed get"+title+"Index on full withdrawerr", "err", err)
						return err
					} else {
						//for i, id := range ids {
						for i := 0; i < len(ids); i ++ {
							id := ids[i]
							if id == can.CandidateId {
								ids = append(ids[:i], ids[i+1:]...)
								i --
							}
						}
						if err := setIndexFn(state, ids); nil != err {
							log.Error("withdraw failed set"+title+"Index on full withdrawerr", "err", err)
							return err
						}
					}
					return nil
				}

	// cache id
	//sortIds := make([]discover.NodeID, 0)
	if flag {
		/** first delete this can on reserves */
		if delreserve {
			if err := delFunc ("Reserve", c.delReserve, c.getReserveIndex, c.setReserveIndex); nil != err {
				return nil, err
			}
		}
		return nodeIds, handle(c.immediateCacheArr, c.setImmediate, c.setImmediateIndex)
	} else {
		/** first delete this can on immediates */
		if delimmediate {
			if err := delFunc ("Immediate", c.delImmediate, c.getImmediateIndex, c.setImmediateIndex); nil != err {
				return nil, err
			}
		}
		return nodeIds, handle(c.reserveCacheArr, c.setReserve, c.setReserveIndex)
	}
}

// Getting immediate or reserve candidate info by nodeId
func (c *CandidatePool) GetCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error) {
	return c.getCandidate(state, nodeId)
}

// candidate withdraw from immediates or reserve elected candidates // TODO 需要抽小方法
func (c *CandidatePool) WithdrawCandidate(state vm.StateDB, nodeId discover.NodeID, price, blockNumber, parentNumber *big.Int, parentHash common.Hash) error {
	var nodeIds []discover.NodeID
	if arr, err := c.withdrawCandidate(state, nodeId, price, blockNumber); nil != err {
		return err
	} else {
		nodeIds = arr
	}
	return ticketPool.DropReturnTicket(state, parentNumber, parentHash, nodeIds...)
}

func (c *CandidatePool) withdrawCandidate(state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) ([]discover.NodeID, error) {
	log.Info("WithdrawCandidate...")
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 2); nil != err {
		log.Error("Failed to initDataByState on WithdrawCandidate", " err", err)
		return nil, err
	}

	if price.Cmp(new(big.Int).SetUint64(0)) <= 0 {
		log.Error("withdraw failed price invalid", " price", price.String())
		return nil, WithdrawPriceErr
	}
	// cache
	var can *types.Candidate
	var flag bool

	if imCan, ok := c.immediateCandidates[nodeId]; (!ok || nil == imCan) {
		reCan, ok := c.reserveCandidates[nodeId]
		if !ok || nil == reCan {
			log.Error("withdraw failed current Candidate is empty")
			return nil, CandidateEmptyErr
		} else {
			can = reCan
		}
	} else {
		can = imCan
		flag = true
	}

	// check withdraw price
	if can.Deposit.Cmp(price) < 0 {
		log.Error("withdraw failed refund price must less or equal deposit", "key", nodeId.String())
		return nil, WithdrawPriceErr
	} else if can.Deposit.Cmp(price) == 0 { // full withdraw

		// return ticket call
		// TODO
		//ticketPool.ReturnTicket(state, )

		handle := func(tiltle string, delInfoFn func(state vm.StateDB, candidateId discover.NodeID),
			getIndexFn func(state vm.StateDB) ([]discover.NodeID, error),
			setIndexFn func(state vm.StateDB, nodeIds []discover.NodeID) error) error {

			// delete current candidate from this elected candidates
			delInfoFn(state, nodeId)
			// update this id index
			if ids, err := getIndexFn(state); nil != err {
				log.Error("withdraw failed get"+tiltle+"Index on full withdrawerr", "err", err)
				return err
			} else {
				for i := 0; i < len(ids); i ++ {
					id := ids[i]
					if id == nodeId {
						ids = append(ids[:i], ids[i+1:]...)
						i --
					}
				}
				if err := setIndexFn(state, ids); nil != err {
					log.Error("withdraw failed set"+tiltle+"Index on full withdrawerr", "err", err)
					return err
				}
			}
			return nil
		}

		if flag {
			if err := handle("Immediate", c.delImmediate, c.getImmediateIndex, c.setImmediateIndex); nil != err {
				return nil, err
			}

		} else {
			if err := handle("Reserve", c.delReserve, c.getReserveIndex, c.setReserveIndex); nil != err {
				return nil, err
			}
		}

		// append to refund (defeat) trie
		if err := c.setDefeat(state, nodeId, can); nil != err {
			log.Error("withdraw failed setDefeat on full withdrawerr", "err", err)
			return nil, err
		}
		// update index of defeat on trie
		if err := c.setDefeatIndex(state); nil != err {
			log.Error("withdraw failed setDefeatIndex on full withdrawerr", "err", err)
			return nil, err
		}

		return []discover.NodeID{nodeId}, nil

	} else { // withdraw a few ...
		// Only withdraw part of the refunds, need to reorder the immediate elected candidates
		// The remaining candiate price to update current candidate info
		canNew := &types.Candidate{
			Deposit:     new(big.Int).Sub(can.Deposit, price),
			BlockNumber: can.BlockNumber,
			TxIndex:     can.TxIndex,
			CandidateId: can.CandidateId,
			Host:        can.Host,
			Port:        can.Port,
			Owner:       can.Owner,
			From:        can.From,
			Extra:       can.Extra,
			Fee:         can.Fee,
		}

		handle := func(title string, candidateArr []*types.Candidate, candidateMap candidateStorage,
			setInfoFn func(state vm.StateDB, candidateId discover.NodeID, can *types.Candidate) error,
			setIndexFn func(state vm.StateDB, nodeIds []discover.NodeID) error) error {
			// update current candidate
			if err := setInfoFn(state, nodeId, canNew); nil != err {
				log.Error("withdraw failed set"+title+" on a few of withdrawerr", "err", err)
				return err
			}

			// sort current candidates
			candidateArr = make([]*types.Candidate, 0)
			for _, can := range candidateMap {
				candidateArr = append(candidateArr, can)
			}
			candidateSort(candidateArr)
			ids := make([]discover.NodeID, 0)
			for _, can := range candidateArr {
				ids = append(ids, can.CandidateId)
			}
			// update new index
			if err := setIndexFn(state, ids); nil != err {
				log.Error("withdraw failed set"+title+"Index on a few of withdrawerr", "err", err)
				return err
			}
			return nil
		}

		if flag {
			if err := handle("Immediate", c.immediateCacheArr, c.immediateCandidates, c.setImmediate, c.setImmediateIndex); nil != err {
				return nil, err
			}
		} else {
			if err := handle("Reserve", c.reserveCacheArr, c.reserveCandidates, c.setReserve, c.setReserveIndex); nil != err {
				return nil, err
			}
		}

		// the withdraw price to build a new refund into defeat on trie
		canDefeat := &types.Candidate{
			Deposit:     price,
			BlockNumber: blockNumber,
			TxIndex:     can.TxIndex,
			CandidateId: can.CandidateId,
			Host:        can.Host,
			Port:        can.Port,
			Owner:       can.Owner,
			From:        can.From,
			Extra:       can.Extra,
			Fee:         can.Fee,
		}
		// the withdraw
		if err := c.setDefeat(state, nodeId, canDefeat); nil != err {
			log.Error("withdraw failed setDefeat on a few of withdrawerr", "err", err)
			return nil, err
		}
		// update index of defeat on trie
		if err := c.setDefeatIndex(state); nil != err {
			log.Error("withdraw failed setDefeatIndex on a few of withdrawerr", "err", err)
			return nil, err
		}
	}
	return nil, nil
}

// Getting elected candidates array
// flag:
// 0:  Getting all elected candidates array
// 1:  Getting all immediate elected candidates array
// 2:  Getting all reserve elected candidates array
func (c *CandidatePool) GetChosens(state vm.StateDB, flag int) []*types.Candidate {
	log.Info("获取入围列表...")
	c.lock.RLock()
	defer c.lock.RUnlock()
	arr := make([]*types.Candidate, 0)
	if err := c.initDataByState(state, 1); nil != err {
		log.Error("Failed to initDataByState on GetChosens", "err", err)
		return arr
	}
	if flag == 0 || flag == 1 {
		immediateIds, err := c.getImmediateIndex(state)
		if nil != err {
			log.Error("Failed to getImmediateIndex on GetChosens", "err", err)
			return arr
		}
		for _, id := range immediateIds {
			arr = append(arr, c.immediateCandidates[id])
		}
	}
	if flag == 0 || flag == 2 {
		reserveIds, err := c.getReserveIndex(state)
		if nil != err {
			log.Error("Failed to getReserveIndex on GetChosens", "err", err)
			return make([]*types.Candidate, 0)
		}
		for _, id := range reserveIds {
			arr = append(arr, c.reserveCandidates[id])
		}
	}
	return arr
}

// Getting all witness array
func (c *CandidatePool) GetChairpersons(state vm.StateDB) []*types.Candidate {
	log.Info("获取本轮见证人列表...")
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 0); nil != err {
		log.Error("Failed to initDataByState on GetChairpersons", "err", err)
		return nil
	}
	witnessIds, err := c.getWitnessIndex(state)
	if nil != err {
		log.Error("Failed to getWitnessIndex on GetChairpersonserr", "err", err)
		return nil
	}
	arr := make([]*types.Candidate, 0)
	for _, id := range witnessIds {
		arr = append(arr, c.originCandidates[id])
	}
	return arr
}

// Getting all refund array by nodeId
func (c *CandidatePool) GetDefeat(state vm.StateDB, nodeId discover.NodeID) ([]*types.Candidate, error) {
	log.Info("获取退款列表: nodeId = " + nodeId.String())
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 2); nil != err {
		log.Error("Failed to initDataByState on GetDefeat", "err", err)
		return nil, err
	}

	defeat, ok := c.defeatCandidates[nodeId]
	if !ok {
		log.Error("Candidate is empty")
		return nil, nil
	}
	return defeat, nil
}

// Checked current candidate was defeat by nodeId
func (c *CandidatePool) IsDefeat(state vm.StateDB, nodeId discover.NodeID) (bool, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 1); nil != err {
		log.Error("Failed to initDataByState on IsDefeat", "err", err)
		return false, err
	}
	if _, ok := c.immediateCandidates[nodeId]; ok {
		return false, nil
	}
	if _, ok := c.reserveCandidates[nodeId]; ok {
		return false, nil
	}
	if arr, ok := c.defeatCandidates[nodeId]; ok && len(arr) != 0 {
		return true, nil
	}
	return false, nil
}

func (c *CandidatePool) IsChosens(state vm.StateDB, nodeId discover.NodeID) (bool, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 1); nil != err {
		log.Error("Failed to initDataByState on IsDefeat", "err", err)
		return false, err
	}
	if _, ok := c.immediateCandidates[nodeId]; ok {
		return true, nil
	}
	if _, ok := c.reserveCandidates[nodeId]; ok {
		return true, nil
	}
	return false, nil
}

// Getting owner's address of candidate info by nodeId
func (c *CandidatePool) GetOwner(state vm.StateDB, nodeId discover.NodeID) common.Address {
	log.Info("获取收益者地址: nodeId = " + nodeId.String())
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 2); nil != err {
		log.Error("Failed to initDataByState on GetOwner", "err", err)
		return common.Address{}
	}
	pre_can, pre_ok := c.preOriginCandidates[nodeId]
	or_can, or_ok := c.originCandidates[nodeId]
	ne_can, ne_ok := c.nextOriginCandidates[nodeId]
	im_can, im_ok := c.immediateCandidates[nodeId]
	re_can, re_ok := c.reserveCandidates[nodeId]
	canArr, de_ok := c.defeatCandidates[nodeId]

	if pre_ok {
		return pre_can.Owner
	}
	if or_ok {
		return or_can.Owner
	}
	if ne_ok {
		return ne_can.Owner
	}
	if im_ok {
		return im_can.Owner
	}
	if re_ok {
		return re_can.Owner
	}
	if de_ok {
		if len(canArr) != 0 {
			return canArr[0].Owner
		}
	}
	return common.Address{}
}

// refund once
func (c *CandidatePool) RefundBalance(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) error {
	log.Info("一键退款: nodeId = " + nodeId.String() + ",当前块高:" + blockNumber.String())
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 2); nil != err {
		log.Error("Failed to initDataByState on RefundBalance", "err", err)
		return err
	}

	var canArr []*types.Candidate
	if defeatArr, ok := c.defeatCandidates[nodeId]; ok {
		canArr = defeatArr
	} else {
		log.Error("Failed to refundbalance candidate is empty")
		return CandidateDecodeErr
	}
	// cache
	// Used for verification purposes, that is, the beneficiary in the pledge refund information of each nodeId should be the same
	var addr common.Address
	// Grand total refund amount for one-time
	amount := big.NewInt(0)
	// Transfer refund information that needs to be deleted
	delCanArr := make([]*types.Candidate, 0)

	contractBalance := state.GetBalance(common.CandidatePoolAddr)
	//currentNum := new(big.Int).SetUint64(blockNumber)

	// Traverse all refund information belong to this nodeId
	for index := 0; index < len(canArr); index ++ {
		can := canArr[index]
		sub := new(big.Int).Sub(blockNumber, can.BlockNumber)
		log.Info("检查退款信息", "当前块高:", blockNumber.String(), "质押块高:", can.BlockNumber.String(), "相差:", sub.String())
		log.Info("检查退款信息", "当前nodeId:", can.CandidateId.String())
		if sub.Cmp(new(big.Int).SetUint64(c.RefundBlockNumber)) >= 0 { // allow refund
			delCanArr = append(delCanArr, can)
			canArr = append(canArr[:index], canArr[index+1:]...)
			index --
			// add up the refund price
			amount = new(big.Int).Add(amount, can.Deposit)
			//amount += can.Deposit.Uint64()
		} else {
			log.Error("block height number had mismatch, No refunds allowed", "current block height", blockNumber.String(), "deposit block height", can.BlockNumber.String(), "allowed block interval", c.RefundBlockNumber)
			log.Info("块高不匹配，不给予退款...")
			continue
		}

		if addr == common.ZeroAddr {
			addr = can.Owner
		} else {
			if addr != can.Owner {
				log.Info("Failed to refundbalance couse current nodeId had bind different owner address ", "nodeId", nodeId.String(), "addr1", addr.String(), "addr2", can.Owner)
				if len(canArr) != 0 {
					canArr = append(delCanArr, canArr...)
				} else {
					canArr = delCanArr
				}
				c.defeatCandidates[nodeId] = canArr
				log.Info("Failed to refundbalance 发现抵押节点nodeId下有不同受益者地址", "nodeId", nodeId.String(), "addr1", addr.String(), "addr2", can.Owner)
				return CandidateOwnerErr
			}
		}

		// check contract account balance
		//if (contractBalance.Cmp(new(big.Int).SetUint64(amount))) < 0 {
		if (contractBalance.Cmp(amount)) < 0 {
			log.Error("Failed to refundbalance constract account insufficient balance ", "contract's balance", contractBalance.String(), "amount", amount.String())
			if len(canArr) != 0 {
				canArr = append(delCanArr, canArr...)
			} else {
				canArr = delCanArr
			}
			c.defeatCandidates[nodeId] = canArr
			return ContractBalanceNotEnoughErr
		}
	}

	// update the tire
	if len(canArr) == 0 {
		c.delDefeat(state, nodeId)
		if ids, err := getDefeatIdsByState(state); nil != err {
			for i := 0; i < len(ids); i ++ {
				id := ids[i]
				if id == nodeId {
					ids = append(ids[:i], ids[i+1:]...)
					i --
				}
			}
			if len(ids) != 0 {
				if value, err := rlp.EncodeToBytes(&ids); nil != err {
					log.Error("Failed to encode candidate ids on RefundBalance", "err", err)
					return CandidateEncodeErr
				} else {
					setDefeatIdsState(state, value)
				}
			} else {
				setDefeatIdsState(state, []byte{})
			}

		}
	} else {
		// If have some remaining, update that
		if arrVal, err := rlp.EncodeToBytes(canArr); nil != err {
			log.Error("Failed to encode candidate object on RefundBalance", "key", nodeId.String(), "err", err)
			canArr = append(delCanArr, canArr...)
			c.defeatCandidates[nodeId] = canArr
			return CandidateDecodeErr
		} else {
			// update the refund information
			setDefeatState(state, nodeId, arrVal)
			// remaining set back to defeat map
			c.defeatCandidates[nodeId] = canArr
		}
	}
	log.Info("退款转账操作：", "nodeId", nodeId.String(), "contractAddr", common.CandidatePoolAddr.String(), "受益人addr", addr.String(), "退还需转账金额:", amount.String())

	// sub contract account balance
	state.SubBalance(common.CandidatePoolAddr, amount)
	// add owner balace
	state.AddBalance(addr, amount)
	log.Info("一键退款完成...")
	return nil
}

// set elected candidate extra value
func (c *CandidatePool) SetCandidateExtra(state vm.StateDB, nodeId discover.NodeID, extra string) error {
	log.Info("设置扩展信息: nodeId = " + nodeId.String())
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 1); nil != err {
		log.Error("Failed to initDataByState on SetCandidateExtra", "err", err)
		return err
	}

	if can, ok := c.immediateCandidates[nodeId]; ok {
		// update current candidate info and update to tire
		can.Extra = extra
		if err := c.setImmediate(state, nodeId, can); nil != err {
			log.Error("Failed to setImmediate on SetCandidateExtra", "err", err)
			return err
		}
	} else {
		if can, ok := c.reserveCandidates[nodeId]; ok {
			can.Extra = extra
			if err := c.setReserve(state, nodeId, can); nil != err {
				log.Error("Failed to setReserve on SetCandidateExtra", "err", err)
				return err
			}
		} else {
			return CandidateEmptyErr
		}
	}
	return nil
}

// Announce witness // TODO 需要抽出小方法
func (c *CandidatePool) Election(state *state.StateDB, parentNumber *big.Int, parentHash common.Hash) ([]*discover.Node, error) {
	var nodes []*discover.Node
	var cans []*types.Candidate
	if nodeArr, canArr, err := c.election(state, parentHash); nil != err {
		return nil, err
	}else {
		nodes, cans = nodeArr, canArr
	}


	nodeIds := make([]discover.NodeID, 0)
	for _, can := range cans {
		// 释放幸运票 TODO
		if err := ticketPool.ReturnTicket(state, can.CandidateId, can.TicketId, big.NewInt(0), common.HexToHash("")); nil != err {
			log.Error("Failed to ReturnTicket on Election", "nodeId", can.CandidateId.String(), "ticketId", can.TicketId.String(), "err", err)
			continue
		}


		/**
		获取TCount TODO
		然后需要对入围候选人再次做排序 TODO
		 */
		c.lock.Lock()
		if err := c.initDataByState(state, 2); nil != err {
			c.lock.Unlock()
			log.Error("Failed to initDataByState on Election", "nodeId", can.CandidateId.String(), " err", err)
			return nil, err
		}
		// 这个和 setCandidate  还不一样。。。
		// 因为需要先判断是否之前在 immediates 中，如果是则转移到 reserves 中
		if _, err := c.setCandidateInfo(state, can.CandidateId, can); nil != err {
			c.lock.Unlock()
			log.Error("Failed to setCandidateInfo on Election", "nodeId", can.CandidateId.String(), "err", err)
			return nil, err
		}

		//if _, ok := c.immediateCandidates[can.CandidateId]; ok {
		//
		//}
		c.lock.Unlock()
		nodeIds = append(nodeIds, can.CandidateId)
	}
	// 释放落榜的
	if err := ticketPool.DropReturnTicket(state, parentNumber, parentHash, nodeIds...); nil != err {
		log.Error("Failed to DropReturnTicket on Election", "err", err)
		return nil, err
	}
	return nodes, nil
}

func (c *CandidatePool) election (state *state.StateDB, parentHash common.Hash) ([]*discover.Node, []*types.Candidate, error) {
	log.Info("揭榜...", "maxChair", c.maxChair, "maxCount", c.maxCount, "RefundBlockNumber", c.RefundBlockNumber)
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 1); nil != err {
		log.Error("Failed to initDataByState on Election", "err", err)
		return nil, nil, err
	}

	// sort immediate candidates
	candidateSort(c.immediateCacheArr)
	log.Info("揭榜时，排序的数组长度:", "len", len(c.immediateCacheArr))
	PrintObject("揭榜时，排序的数组:", c.immediateCacheArr)
	// cache ids
	immediateIds := make([]discover.NodeID, 0)
	for _, can := range c.immediateCacheArr {
		immediateIds = append(immediateIds, can.CandidateId)
	}
	log.Info("当前入围者ids 长度：", "len", len(immediateIds))
	PrintObject("当前入围者ids：", immediateIds)
	// a certain number of witnesses in front of the cache
	var nextWitIds []discover.NodeID
	// If the number of candidate selected does not exceed the number of witnesses
	if len(immediateIds) <= int(c.maxChair) {
		nextWitIds = make([]discover.NodeID, len(immediateIds))
		copy(nextWitIds, immediateIds)

	} else {
		// If the number of candidate selected exceeds the number of witnesses, the top N is extracted.
		nextWitIds = make([]discover.NodeID, c.maxChair)
		copy(nextWitIds, immediateIds)
	}
	log.Info("选出来的下一轮见证人Ids 个数:", "len", len(nextWitIds))
	PrintObject("选出来的下一轮见证人Ids:", nextWitIds)
	// cache map
	nextWits := make(candidateStorage, 0)

	// copy witnesses information
	copyCandidateMapByIds(nextWits, c.immediateCandidates, nextWitIds)
	log.Info("从入围信息copy过来的见证人个数;", "len", len(nextWits))
	PrintObject("从入围信息copy过来的见证人;", nextWits)
	// clear all old nextwitnesses information （If it is forked, the next round is no empty.）
	for nodeId, _ := range c.nextOriginCandidates {
		c.delNextWitness(state, nodeId)
	}


	arr := make([]*discover.Node, 0)
	caches := make([]*types.Candidate, 0)
	// set up all new nextwitnesses information
	for nodeId, can := range nextWits {

		// 揭榜后回去调 获取幸运票逻辑 TODO
	    luckyId, err := ticketPool.SelectionLuckyTicket(state, nodeId, big.NewInt(0), parentHash)
	    if nil != err {
	    	log.Error("Failed to take luckyId on Election", "nodeId", nodeId.String(), "err", err)
	    	return nil, nil, err
		}


		// 将幸运票ID 置入 next witness 详情中
		can.TicketId = luckyId
		if err := c.setNextWitness(state, nodeId, can); nil != err {
			log.Error("failed to setNextWitness on election", "err", err)
			return nil, nil, err
		}
		caches = append(caches, can)
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build Node on GetWitness", "err", err, "nodeId", can.CandidateId.String())
			continue
		} else {
			arr = append(arr, node)
		}
	}
	// update new nextwitnesses index
	if err := c.setNextWitnessIndex(state, nextWitIds); nil != err {
		log.Error("failed to setNextWitnessIndex on election", "err", err)
		return nil, nil, err
	}
	// replace the next round of witnesses
	c.nextOriginCandidates = nextWits

	log.Info("下一轮见证人node 个数:", "len", len(arr))
	PrintObject("下一轮见证人node信息:", arr)
	log.Info("揭榜完成...")
	return arr, caches, nil
}


// switch next witnesses to current witnesses
func (c *CandidatePool) Switch(state *state.StateDB) bool {
	log.Info("替换见证人...")
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 0); nil != err {
		log.Error("Failed to initDataByState on Switch", "err", err)
		return false
	}
	// clear all old previous witness on trie
	for nodeId, _ := range c.preOriginCandidates {
		c.delPreviousWitness(state, nodeId)
	}
	// set up new witnesses to previous witnesses on trie by current witnesses
	for nodeId, can := range c.originCandidates {
		if err := c.setPreviousWitness(state, nodeId, can); nil != err {
			log.Error("Failed to setPreviousWitness on Switch", "err", err)
			return false
		}
	}
	// update previous witness index by current witness index
	if ids, err := c.getWitnessIndex(state); nil != err {
		log.Error("Failed to getWitnessIndex on Switch", "err", err)
		return false
	} else {
		// replace witnesses index
		if err := c.setPreviousWitnessindex(state, ids); nil != err {
			log.Error("Failed to setPreviousWitnessindex on Switch", "err", err)
			return false
		}
	}

	// clear all old witnesses on trie
	for nodeId, _ := range c.originCandidates {
		c.delWitness(state, nodeId)
	}
	// set up new witnesses to current witnesses on trie by next witnesses
	for nodeId, can := range c.nextOriginCandidates {
		if err := c.setWitness(state, nodeId, can); nil != err {
			log.Error("Failed to setWitness on Switch", "err", err)
			return false
		}
	}
	// update current witness index by next witness index
	if ids, err := c.getNextWitnessIndex(state); nil != err {
		log.Error("Failed to getNextWitnessIndex on Switch", "err", err)
		return false
	} else {
		// replace witnesses index
		if err := c.setWitnessindex(state, ids); nil != err {
			log.Error("Failed to setWitnessindex on Switch", "err", err)
			return false
		}
	}
	// clear all old nextwitnesses information
	for nodeId, _ := range c.nextOriginCandidates {
		c.delNextWitness(state, nodeId)
	}
	// clear next witness index
	c.setNextWitnessIndex(state, make([]discover.NodeID, 0))
	log.Info("替换完成...")
	return true
}

// Getting nodes of witnesses
// flag：-1: the previous round of witnesses  0: the current round of witnesses   1: the next round of witnesses
func (c *CandidatePool) GetWitness(state *state.StateDB, flag int) ([]*discover.Node, error) {
	log.Info("获取见证人: flag = " + strconv.Itoa(flag))
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 0); nil != err {
		log.Error("Failed to initDataByState on GetWitness", "err", err)
		return nil, err
	}
	//var ids []discover.NodeID
	var witness candidateStorage
	if flag == -1 {
		witness = c.preOriginCandidates
	} else if flag == 0 {
		witness = c.originCandidates
	} else if flag == 1 {
		witness = c.nextOriginCandidates
	}

	arr := make([]*discover.Node, 0)
	for _, can := range witness {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build Node on GetWitness", "err", err, "nodeId", can.CandidateId.String())
			return nil, err
		} else {
			arr = append(arr, node)
		}
	}
	return arr, nil
}

// Getting previous and current and next witnesses
func (c *CandidatePool) GetAllWitness(state *state.StateDB) ([]*discover.Node, []*discover.Node, []*discover.Node, error) {
	log.Info("获取所有见证人...")
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 0); nil != err {
		log.Error("Failed to initDataByState on GetAllWitness", "err", err)
		return nil, nil, nil, err
	}
	//var ids []discover.NodeID
	var prewitness, witness, nextwitness candidateStorage
	prewitness = c.preOriginCandidates
	witness = c.originCandidates
	nextwitness = c.nextOriginCandidates

	preArr, curArr, nextArr := make([]*discover.Node, 0), make([]*discover.Node, 0), make([]*discover.Node, 0)
	for _, can := range prewitness {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build pre Node on GetAllWitness", "err", err, "nodeId", can.CandidateId.String())
			//continue
			return nil, nil, nil, err
		} else {
			preArr = append(preArr, node)
		}
	}
	for _, can := range witness {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build cur Node on GetAllWitness", "err", err, "nodeId", can.CandidateId.String())
			//continue
			return nil, nil, nil, err
		} else {
			curArr = append(curArr, node)
		}
	}
	for _, can := range nextwitness {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build next Node on GetAllWitness", "err", err, "nodeId", can.CandidateId.String())
			//continue
			return nil, nil, nil, err
		} else {
			nextArr = append(nextArr, node)
		}
	}
	return preArr, curArr, nextArr, nil
}

func (c *CandidatePool) GetRefundInterval() uint64 {
	return c.RefundBlockNumber
}

//// update candidate's tickets  // TODO 需要抽小方法
//func (c *CandidatePool) UpdateCandidateTicket(state vm.StateDB, blockNumber *big.Int, blockHash common.Hash, nodeId discover.NodeID, can *types.Candidate) error {
//	var nodeIds []discover.NodeID
//	if arr, err := c.updateCandidateTicket(state, nodeId, can); nil != err {
//		return err
//	} else {
//		nodeIds = arr
//	}
//	return ticketPool.DropReturnTicket(state, blockNumber, blockHash, nodeIds...)
//}
//
//func (c *CandidatePool) updateCandidateTicket(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) ([]discover.NodeID, error) {
//	c.lock.Lock()
//	defer c.lock.Unlock()
//	if err := c.initDataByState(state, 1); nil != err {
//		log.Error("Failed to initDataByState on UpdateCandidateTicket err", err)
//		return nil, err
//	}
//
//	// nodeIds cache for lost elected
//	nodeIds := make([]discover.NodeID, 0)
//
//
//	/** // 还需要补上 如果TCout 小了，要先移到 reserves 中，不然才算落榜 TODO  */
//	// 如果本来就在 候选池中，直接修改即可
//	if _, ok := c.immediateCandidates[nodeId]; ok {
//		if err := c.setImmediate(state, nodeId, can); nil != err {
//			return nil, err
//		}
//	} else { // 否则，去备选池中查
//		if _, ok := c.reserveCandidates[nodeId]; !ok {
//			return nil, CandidateEmptyErr
//		} else {
//			// 则需要先判断当前入参的 can的得票数
//			// 如果还是不满足移到 候选池中，则直接修改 备选池的 can 信息
//			if !checkTicket(ticketidsCache.TCount(), c.maxCount) {
//				if err := c.setReserve(state, nodeId, can); nil != err {
//					return nil, err
//				}
//			} else {
//				// 否则，往候选池中追加该 can
//				// 且， 备选池中的该 can 的信息删除
//				/** 先删除 备选池 can 信息 */
//				c.delReserve(state, can.CandidateId)
//				// update reserve id index
//				if ids, err := c.getReserveIndex(state); nil != err {
//					log.Error("withdraw failed getReserveIndex on full withdrawerr", "err", err)
//					return nil, err
//				} else {
//					//for i, id := range ids {
//					for i := 0; i < len(ids); i ++ {
//						id := ids[i]
//						if id == can.CandidateId {
//							ids = append(ids[:i], ids[i+1:]...)
//							i --
//						}
//					}
//					if err := c.setReserveIndex(state, ids); nil != err {
//						log.Error("withdraw failed setReserveIndex on full withdrawerr", "err", err)
//						return nil, err
//					}
//				}
//
//				/** 追加到 候选池 */
//				c.immediateCandidates[can.CandidateId] = can
//				c.immediateCacheArr = make([]*types.Candidate, 0)
//				for _, v := range c.immediateCandidates {
//					c.immediateCacheArr = append(c.immediateCacheArr, v)
//				}
//
//				// sort cache array
//				candidateSort(c.immediateCacheArr)
//
//				if len(c.immediateCacheArr) > int(c.maxCount) {
//					// Intercepting the lost candidates to tmpArr
//					tmpArr := (c.immediateCacheArr)[c.maxCount:]
//					// qualified elected candidates
//					c.immediateCacheArr = (c.immediateCacheArr)[:c.maxCount]
//
//					// handle tmpArr
//					for _, tmpCan := range tmpArr {
//
//						// return ticket call
//						// TODO
//						//ticketPool.ReturnTicket(state, )
//
//						c.delImmediate(state, tmpCan.CandidateId)
//						// append to refunds (defeat) trie
//						if err := c.setDefeat(state, tmpCan.CandidateId, tmpCan); nil != err {
//							return nil, err
//						}
//						nodeIds = append(nodeIds, tmpCan.CandidateId)
//					}
//
//					// update index of refund (defeat) on trie
//					if err := c.setDefeatIndex(state); nil != err {
//						return nil, err
//					}
//				}
//
//				// cache id
//				sortIds := make([]discover.NodeID, 0)
//
//				/** 设置新的 候选人数组 */
//				// insert elected candidate to tire
//				for _, can := range c.immediateCacheArr {
//					if err := c.setImmediate(state, can.CandidateId, can); nil != err {
//						return nil, err
//					}
//					sortIds = append(sortIds, can.CandidateId)
//				}
//				// update index of immediate elected candidates on trie
//				if err := c.setImmediateIndex(state, sortIds); nil != err {
//					return nil, err
//				}
//			}
//		}
//	}
//	return nodeIds, nil
//}

// 根据nodeId 去重新决定当前候选人的去留
func (c *CandidatePool) UpdateElectedQueue(stateDB vm.StateDB, nodeId ... discover.NodeID) error {

	return nil
}

//func (c *CandidatePool) UpdateWitnessInfo () error {
//
//}

// TODO
func (c *CandidatePool) AddCandidateEpoch() error {
	// 遍历出所有 候选人和备选人 逐个叠加票龄

	return nil
}

func checkTicket(t_count, c_count uint64) bool {
	if t_count >= c_count {
		return true
	}
	return false
}

func (c *CandidatePool) setImmediate(state vm.StateDB, candidateId discover.NodeID, can *types.Candidate) error {
	c.immediateCandidates[candidateId] = can
	if value, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode candidate object on setImmediate", "key", candidateId.String(), "err", err)
		return CandidateEncodeErr
	} else {
		// set immediate candidate input the trie
		setImmediateState(state, candidateId, value)
	}
	return nil
}

func (c *CandidatePool) getImmediateIndex(state vm.StateDB) ([]discover.NodeID, error) {
	return getImmediateIdsByState(state)
}

// deleted immediate candidate by nodeId (Automatically update the index)
func (c *CandidatePool) delImmediate(state vm.StateDB, candidateId discover.NodeID) {
	// deleted immediate candidate by id on trie
	setImmediateState(state, candidateId, []byte{})
	// deleted immedidate candidate by id on map
	delete(c.immediateCandidates, candidateId)
}

func (c *CandidatePool) setImmediateIndex(state vm.StateDB, nodeIds []discover.NodeID) error {
	if len(nodeIds) == 0 {
		setImmediateIdsState(state, []byte{})
		return nil
	}
	if val, err := rlp.EncodeToBytes(nodeIds); nil != err {
		log.Error("Failed to encode ImmediateIds", "err", err)
		return err
	} else {
		setImmediateIdsState(state, val)
	}
	return nil
}

//////////////

func (c *CandidatePool) setReserve(state vm.StateDB, candidateId discover.NodeID, can *types.Candidate) error {
	c.reserveCandidates[candidateId] = can
	if value, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode candidate object on setReserve", "key", candidateId.String(), "err", err)
		return CandidateEncodeErr
	} else {
		// set setReserve candidate input the trie
		setReserveState(state, candidateId, value)
	}
	return nil
}

func (c *CandidatePool) getReserveIndex(state vm.StateDB) ([]discover.NodeID, error) {
	return getReserveIdsByState(state)
}

// deleted reserve candidate by nodeId (Automatically update the index)
func (c *CandidatePool) delReserve(state vm.StateDB, candidateId discover.NodeID) {
	// deleted reserve candidate by id on trie
	setReserveState(state, candidateId, []byte{})
	// deleted reserve candidate by id on map
	delete(c.reserveCandidates, candidateId)
}

func (c *CandidatePool) setReserveIndex (state vm.StateDB, nodeIds []discover.NodeID) error {
	if len(nodeIds) == 0 {
		setReserveIdsState(state, []byte{})
		return nil
	}
	if val, err := rlp.EncodeToBytes(nodeIds); nil != err {
		log.Error("Failed to encode ReserveIds", "err", err)
		return err
	} else {
		setReserveIdsState(state, val)
	}
	return nil
}

// setting refund information
func (c *CandidatePool) setDefeat(state vm.StateDB, candidateId discover.NodeID, can *types.Candidate) error {

	var defeatArr []*types.Candidate
	// append refund information
	if defeatArrTmp, ok := c.defeatCandidates[can.CandidateId]; ok {
		defeatArrTmp = append(defeatArrTmp, can)
		//c.defeatCandidates[can.CandidateId] = defeatArrTmp
		defeatArr = defeatArrTmp
	} else {
		defeatArrTmp = make([]*types.Candidate, 0)
		defeatArrTmp = append(defeatArr, can)
		//c.defeatCandidates[can.CandidateId] = defeatArrTmp
		defeatArr = defeatArrTmp
	}
	// setting refund information on trie
	if value, err := rlp.EncodeToBytes(&defeatArr); nil != err {
		log.Error("Failed to encode candidate object on setDefeat", "key", candidateId.String(), "err", err)
		return CandidateEncodeErr
	} else {
		setDefeatState(state, candidateId, value)
		c.defeatCandidates[can.CandidateId] = defeatArr
	}
	return nil
}

func (c *CandidatePool) delDefeat(state vm.StateDB, nodeId discover.NodeID) {
	delete(c.defeatCandidates, nodeId)
	setDefeatState(state, nodeId, []byte{})
}

// update refund index
func (c *CandidatePool) setDefeatIndex(state vm.StateDB) error {
	newdefeatIds := make([]discover.NodeID, 0)
	for id, _ := range c.defeatCandidates {
		newdefeatIds = append(newdefeatIds, id)
	}
	if len(newdefeatIds) == 0 {
		setDefeatIdsState(state, []byte{})
		return nil
	}
	if value, err := rlp.EncodeToBytes(&newdefeatIds); nil != err {
		log.Error("Failed to encode candidate object on setDefeatIds", "err", err)
		return CandidateEncodeErr
	} else {
		setDefeatIdsState(state, value)
	}
	return nil
}

func (c *CandidatePool) delPreviousWitness(state vm.StateDB, candidateId discover.NodeID) {
	// deleted previous witness by id on map
	delete(c.preOriginCandidates, candidateId)
	// delete previous witness by id on trie
	setPreviousWitnessState(state, candidateId, []byte{})
}

func (c *CandidatePool) setPreviousWitness(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	c.preOriginCandidates[nodeId] = can
	if val, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode Candidate on setPreviousWitness", "err", err)
		return err
	} else {
		setPreviousWitnessState(state, nodeId, val)
	}
	return nil
}

func (c *CandidatePool) setPreviousWitnessindex(state vm.StateDB, nodeIds []discover.NodeID) error {
	if len(nodeIds) == 0 {
		setPreviosWitnessIdsState(state, []byte{})
		return nil
	}
	if val, err := rlp.EncodeToBytes(nodeIds); nil != err {
		log.Error("Failed to encode Previous WitnessIds", "err", err)
		return err
	} else {
		setPreviosWitnessIdsState(state, val)
	}
	return nil
}

func (c *CandidatePool) getPreviousWitnessIndex(state vm.StateDB) ([]discover.NodeID, error) {
	return getPreviousWitnessIdsState(state)
}

func (c *CandidatePool) setWitness(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	c.originCandidates[nodeId] = can
	if val, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode Candidate on setWitness", "err", err)
		return err
	} else {
		setWitnessState(state, nodeId, val)
		//PrintObject("设置 setWitness ", *can)
	}
	return nil
}

func (c *CandidatePool) setWitnessindex(state vm.StateDB, nodeIds []discover.NodeID) error {
	if len(nodeIds) == 0 {
		setWitnessIdsState(state, []byte{})
		return nil
	}
	if val, err := rlp.EncodeToBytes(nodeIds); nil != err {
		log.Error("Failed to encode WitnessIds", "err", err)
		return err
	} else {
		setWitnessIdsState(state, val)
	}
	return nil
}

func (c *CandidatePool) delWitness(state vm.StateDB, candidateId discover.NodeID) {
	// deleted witness by id on map
	delete(c.originCandidates, candidateId)
	// delete witness by id on trie
	setWitnessState(state, candidateId, []byte{})
}

func (c *CandidatePool) getWitnessIndex(state vm.StateDB) ([]discover.NodeID, error) {
	return getWitnessIdsByState(state)
}

func (c *CandidatePool) setNextWitness(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	c.nextOriginCandidates[nodeId] = can
	if value, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode candidate object on setImmediate", "key", nodeId.String(), "err", err)
		return CandidateEncodeErr
	} else {
		// setting next witness information on trie
		setNextWitnessState(state, nodeId, value)
		PrintObject("设置 setNextWitness", *can)
	}
	return nil
}

func (c *CandidatePool) delNextWitness(state vm.StateDB, candidateId discover.NodeID) {
	// deleted next witness by id on map
	delete(c.nextOriginCandidates, candidateId)
	// deleted next witness by id on trie
	setNextWitnessState(state, candidateId, []byte{})
}

func (c *CandidatePool) setNextWitnessIndex(state vm.StateDB, nodeIds []discover.NodeID) error {
	if len(nodeIds) == 0 {
		setNextWitnessIdsState(state, []byte{})
		return nil
	}
	if value, err := rlp.EncodeToBytes(&nodeIds); nil != err {
		log.Error("Failed to encode candidate object on setDefeatIds", "err", err)
		return CandidateEncodeErr
	} else {
		setNextWitnessIdsState(state, value)
		//PrintObject("设置 setNextWitnessIndex:", nodeIds)
	}
	return nil
}

func (c *CandidatePool) getNextWitnessIndex(state vm.StateDB) ([]discover.NodeID, error) {
	return getNextWitnessIdsByState(state)
}

func (c *CandidatePool) getCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 1); nil != err {
		log.Error("Failed to initDataByState on getCandidate", "err", err)
		return nil, err
	}
	if candidatePtr, ok := c.immediateCandidates[nodeId]; ok {
		PrintObject("GetCandidate 返回 immediate：", *candidatePtr)
		return candidatePtr, nil
	}
	if candidatePtr, ok := c.reserveCandidates[nodeId]; ok {
		PrintObject("GetCandidate 返回 reserve：", *candidatePtr)
		return candidatePtr, nil
	}
	return nil, nil
}

func (c *CandidatePool) MaxChair() uint64 {
	return c.maxChair
}

func getPreviousWitnessIdsState(state vm.StateDB) ([]discover.NodeID, error) {
	var witnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, PreviousWitnessListKey()); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &witnessIds); nil != err {
			return nil, err
		}
	}
	return witnessIds, nil
}

func setPreviosWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, PreviousWitnessListKey(), arrVal)
}

func getPreviousWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, PreviousWitnessKey(id)); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	}
	return &can, nil
}

func setPreviousWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, PreviousWitnessKey(id), val)
}

func getWitnessIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var witnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, WitnessListKey()); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &witnessIds); nil != err {
			return nil, err
		}
	}
	return witnessIds, nil
}

func setWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, WitnessListKey(), arrVal)
}

func getWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, WitnessKey(id)); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	}
	return &can, nil
}

func setWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, WitnessKey(id), val)
}

func getNextWitnessIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var nextWitnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, NextWitnessListKey()); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &nextWitnessIds); nil != err {
			return nil, err
		}
	}
	return nextWitnessIds, nil
}

func setNextWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, NextWitnessListKey(), arrVal)
}

func getNextWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, NextWitnessKey(id)); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	}
	return &can, nil
}

func setNextWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, NextWitnessKey(id), val)
}

func getImmediateIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var immediateIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, ImmediateListKey()); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &immediateIds); nil != err {
			return nil, err
		}
	}
	return immediateIds, nil
}

func setImmediateIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, ImmediateListKey(), arrVal)
}

func getImmediateByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, ImmediateKey(id)); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	}
	return &can, nil
}

func setImmediateState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, ImmediateKey(id), val)
}

func getReserveIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var reserveIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, ReserveListKey()); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &reserveIds); nil != err {
			return nil, err
		}
	}
	return reserveIds, nil
}

func setReserveIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, ReserveListKey(), arrVal)
}

func getReserveByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, ReserveKey(id)); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &can); nil != err {
			return nil, err
		}
	}
	return &can, nil
}

func setReserveState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, ReserveKey(id), val)
}

func getDefeatIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var defeatIds []discover.NodeID
	if valByte := state.GetState(common.CandidatePoolAddr, DefeatListKey()); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &defeatIds); nil != err {
			return nil, err
		}
	}
	return defeatIds, nil
}

func setDefeatIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidatePoolAddr, DefeatListKey(), arrVal)
}

func getDefeatsByState(state vm.StateDB, id discover.NodeID) ([]*types.Candidate, error) {
	var canArr []*types.Candidate
	if valByte := state.GetState(common.CandidatePoolAddr, DefeatKey(id)); nil != valByte && len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &canArr); nil != err {
			return nil, err
		}
	}
	return canArr, nil
}

func setDefeatState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidatePoolAddr, DefeatKey(id), val)
}

func copyCandidateMapByIds(target, source candidateStorage, ids []discover.NodeID) {
	for _, id := range ids {
		target[id] = source[id]
	}
}

func GetCandidatePtr() *CandidatePool {
	return candidatePool
}

func PrintObject(s string, obj interface{}) {
	objs, _ := json.Marshal(obj)

	log.Info(s, "==", string(objs))
	//fmt.Println(s, string(objs))
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

func compare(c, can *types.Candidate) int {
	// put the larger deposit in front
	if c.Deposit.Cmp(can.Deposit) > 0 {
		return 1
	} else if c.Deposit.Cmp(can.Deposit) == 0 {
		// put the smaller blocknumber in front
		if c.BlockNumber.Cmp(can.BlockNumber) > 0 {
			return -1
		} else if c.BlockNumber.Cmp(can.BlockNumber) == 0 {
			// put the smaller tx'index in front
			if c.TxIndex > can.TxIndex {
				return -1
			} else if c.TxIndex == can.TxIndex {
				return 0
			} else {
				return 1
			}
		} else {
			return 1
		}
	} else {
		return -1
	}
}

// sorted candidates
func candidateSort(arr []*types.Candidate) {
	if len(arr) <= 1 {
		return
	}
	quickSort(arr, 0, len(arr)-1)
}
func quickSort(arr []*types.Candidate, left, right int) {
	if left < right {
		pivot := partition(arr, left, right)
		quickSort(arr, left, pivot-1)
		quickSort(arr, pivot+1, right)
	}
}
func partition(arr []*types.Candidate, left, right int) int {
	for left < right {
		for left < right && compare(arr[left], arr[right]) >= 0 {
			right--
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			left++
		}
		for left < right && compare(arr[left], arr[right]) >= 0 {
			left++
		}
		if left < right {
			arr[left], arr[right] = arr[right], arr[left]
			right--
		}
	}
	return left
}

func ImmediateKey(nodeId discover.NodeID) []byte {
	return immediateKey(nodeId.Bytes())
}
func immediateKey(key []byte) []byte {
	return append(ImmediateBytePrefix, key...)
}

func ReserveKey(nodeId discover.NodeID) []byte {
	return reserveKey(nodeId.Bytes())
}

func reserveKey(key []byte) []byte {
	return append(ReserveBytePrefix, key...)
}

func PreviousWitnessKey(nodeId discover.NodeID) []byte {
	return prewitnessKey(nodeId.Bytes())
}

func prewitnessKey(key []byte) []byte {
	return append(PreWitnessBytePrefix, key...)
}

func WitnessKey(nodeId discover.NodeID) []byte {
	return witnessKey(nodeId.Bytes())
}
func witnessKey(key []byte) []byte {
	return append(WitnessBytePrefix, key...)
}

func NextWitnessKey(nodeId discover.NodeID) []byte {
	return nextWitnessKey(nodeId.Bytes())
}
func nextWitnessKey(key []byte) []byte {
	return append(NextWitnessBytePrefix, key...)
}

func DefeatKey(nodeId discover.NodeID) []byte {
	return defeatKey(nodeId.Bytes())
}
func defeatKey(key []byte) []byte {
	return append(DefeatBytePrefix, key...)
}

func ImmediateListKey() []byte {
	return ImmediateListBytePrefix
}

func ReserveListKey() []byte {
	return ReserveListBytePrefix
}

func PreviousWitnessListKey() []byte {
	return PreWitnessListBytePrefix
}

func WitnessListKey() []byte {
	return WitnessListBytePrefix
}

func NextWitnessListKey() []byte {
	return NextWitnessListBytePrefix
}

func DefeatListKey() []byte {
	return DefeatListBytePrefix
}
