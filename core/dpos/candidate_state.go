package depos

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

const (
	// immediate elected candidate
	ImmediatePrefix     = "id"
	ImmediateListPrefix = "iL"
	// previous witness
	PreWitnessPrefix     = "Pwn"
	PreWitnessListPrefix = "PwL"
	// witness
	WitnessPrefix     = "wn"
	WitnessListPrefix = "wL"
	// next witness
	NextWitnessPrefix     = "Nwn"
	NextWitnessListPrefix = "NwL"
	// need refund
	DefeatPrefix     = "df"
	DefeatListPrefix = "dL"
)

var (
	// immediate elected candidate
	ImmediateBtyePrefix     = []byte(ImmediatePrefix)
	ImmediateListBtyePrefix = []byte(ImmediateListPrefix)
	// previous witness
	PreWitnessBytePrefix     = []byte(PreWitnessPrefix)
	PreWitnessListBytePrefix = []byte(PreWitnessListPrefix)
	// witness
	WitnessBtyePrefix     = []byte(WitnessPrefix)
	WitnessListBtyePrefix = []byte(WitnessListPrefix)
	// next witness
	NextWitnessBtyePrefix     = []byte(NextWitnessPrefix)
	NextWitnessListBytePrefix = []byte(NextWitnessListPrefix)
	// need refund
	DefeatBtyePrefix     = []byte(DefeatPrefix)
	DefeatListBtyePrefix = []byte(DefeatListPrefix)

	CandidateEncodeErr          = errors.New("Candidate encoding err")
	CandidateDecodeErr          = errors.New("Candidate decoding err")
	WithdrawPriceErr            = errors.New("Withdraw Price err")
	CandidateEmptyErr           = errors.New("Candidate is empty")
	ContractBalanceNotEnoughErr = errors.New("Contract's balance is not enough")
	CandidateOwnerErr           = errors.New("CandidateOwner Addr is illegal")
)

type CandidatePool struct {
	// allow immediate elected max count
	maxCount uint64
	// allow witness max count
	maxChair uint64
	// allow block interval for refunds
	RefundBlockNumber uint64

	// previous witness
	preOriginCandidates map[discover.NodeID]*types.Candidate
	// current witnesses
	originCandidates map[discover.NodeID]*types.Candidate
	// next witnesses
	nextOriginCandidates map[discover.NodeID]*types.Candidate
	// immediates
	immediateCandates map[discover.NodeID]*types.Candidate
	// refunds
	defeatCandidates map[discover.NodeID][]*types.Candidate

	// cache
	candidateCacheArr []*types.Candidate
	lock              *sync.RWMutex
}

var candidatePool *CandidatePool

// Initialize the global candidate pool object
func NewCandidatePool(configs *params.DposConfig) *CandidatePool {
	PrintObject("初始化配置文件", *configs)
	candidatePool = &CandidatePool{
		maxCount:             configs.Candidate.MaxCount,
		maxChair:             configs.Candidate.MaxChair,
		RefundBlockNumber:    configs.Candidate.RefundBlockNumber,
		preOriginCandidates:  make(map[discover.NodeID]*types.Candidate, 0),
		originCandidates:     make(map[discover.NodeID]*types.Candidate, 0),
		nextOriginCandidates: make(map[discover.NodeID]*types.Candidate, 0),
		immediateCandates:    make(map[discover.NodeID]*types.Candidate, 0),
		defeatCandidates:     make(map[discover.NodeID][]*types.Candidate, 0),
		candidateCacheArr:    make([]*types.Candidate, 0),
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
	c.preOriginCandidates = make(map[discover.NodeID]*types.Candidate, 0)
	if ids, err := getPreviousWitnessIdsState(state); nil != err {
		log.Error("Failed to decode previous witnessIds on initDataByState err", err)
		return err
	} else {
		prewitnessIds = ids
	}
	PrintObject("prewitnessIds = ", prewitnessIds)
	for _, witnessId := range prewitnessIds {
		//fmt.Println("prewitnessId = ", witnessId.String())
		var can *types.Candidate
		if ca, err := getPreviousWitnessByState(state, witnessId); nil != err {
			log.Error("Failed to decode Candidate on initDataByState", "err", err)
			return CandidateDecodeErr
		} else {
			can = ca
		}
		c.preOriginCandidates[witnessId] = can
	}

	// loading current witnesses
	var witnessIds []discover.NodeID
	c.originCandidates = make(map[discover.NodeID]*types.Candidate, 0)
	if ids, err := getWitnessIdsByState(state); nil != err {
		log.Error("Failed to decode witnessIds on initDataByState", "err", err)
		return err
	} else {
		witnessIds = ids
	}
	PrintObject("witnessIds = ", witnessIds)
	for _, witnessId := range witnessIds {
		//fmt.Println("witnessId = ", witnessId.String())
		var can *types.Candidate
		if ca, err := getWitnessByState(state, witnessId); nil != err {
			log.Error("Failed to decode Candidate on initDataByState", "err", err)
			return CandidateDecodeErr
		} else {
			can = ca
		}
		c.originCandidates[witnessId] = can
	}

	// loading next witnesses
	var nextWitnessIds []discover.NodeID
	c.nextOriginCandidates = make(map[discover.NodeID]*types.Candidate, 0)
	if ids, err := getNextWitnessIdsByState(state); nil != err {
		log.Error("Failed to decode nextWitnessIds on initDataByState", "err", err)
		return err
	} else {
		nextWitnessIds = ids
	}
	PrintObject("nextWitnessIds = ", nextWitnessIds)
	for _, witnessId := range nextWitnessIds {
		//fmt.Println("nextwitnessId = ", witnessId.String())
		var can *types.Candidate
		if ca, err := getNextWitnessByState(state, witnessId); nil != err {
			log.Error("Failed to decode Candidate on initDataByState", "err", err)
			return CandidateDecodeErr
		} else {
			can = ca
		}
		c.nextOriginCandidates[witnessId] = can
	}

	if flag == 1 || flag == 2 {
		// loading immediate elected candidates
		var immediateIds []discover.NodeID
		c.immediateCandates = make(map[discover.NodeID]*types.Candidate, 0)
		if ids, err := getImmediateIdsByState(state); nil != err {
			log.Error("Failed to decode immediateIds on initDataByState", "err", err)
			return err
		} else {
			immediateIds = ids
		}

		// cache
		canCache := make([]*types.Candidate, 0)

		PrintObject("immediateIds = ", immediateIds)
		for _, immediateId := range immediateIds {
			//fmt.Println("immediateId = ", immediateId.String())
			var can *types.Candidate
			if ca, err := getImmediateByState(state, immediateId); nil != err {
				log.Error("Failed to decode Candidate on initDataByState", "err", err)
				return CandidateDecodeErr
			} else {
				can = ca
			}
			c.immediateCandates[immediateId] = can
			canCache = append(canCache, can)
		}
		c.candidateCacheArr = canCache
	}

	if flag == 2 {
		// load refunds
		var defeatIds []discover.NodeID
		c.defeatCandidates = make(map[discover.NodeID][]*types.Candidate, 0)
		if ids, err := getDefeatIdsByState(state); nil != err {
			log.Error("Failed to decode defeatIds on initDataByState", "err", err)
			return err
		} else {
			defeatIds = ids
		}
		PrintObject("defeatIds = ", defeatIds)
		for _, defeatId := range defeatIds {
			//fmt.Println("defeatId = ", defeatId.String())
			var canArr []*types.Candidate
			if arr, err := getDefeatsByState(state, defeatId); nil != err {
				log.Error("Failed to decode CandidateArr on initDataByState", "err", err)
				return CandidateDecodeErr
			} else {
				canArr = arr
			}
			c.defeatCandidates[defeatId] = canArr
		}
	}
	return nil
}

// pledge Candidate
func (c *CandidatePool) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	PrintObject("SetCandidate:", can)
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 2); nil != err {
		log.Error("Failed to initDataByState on SetCandidate err", err)
		return err
	}
	c.immediateCandates[can.CandidateId] = can
	c.candidateCacheArr = make([]*types.Candidate, 0)
	// append to the cache array and then sort
	if len(c.immediateCandates) != 0 && len(c.candidateCacheArr) == 0 {
		for _, v := range c.immediateCandates {
			c.candidateCacheArr = append(c.candidateCacheArr, v)
		}
	}

	// Whether the current candidate is new
	// then append to cache array
	//var needSort bool
	//if _, ok := c.immediateCandates[can.CandidateId]; !ok {
	//	c.candidateCacheArr = append(c.candidateCacheArr, can)
	//	//needSort = true
	//}

	PrintObject("SetCandidate immediateMap:", c.immediateCandates)
	// sort cache array
	candidateSort(c.candidateCacheArr)
	//if needSort {
	//	candidateSort(c.candidateCacheArr)
	//}
	// move the excessive of immediate elected candidate to refunds
	if len(c.candidateCacheArr) > int(c.maxCount) {
		// Intercepting the lost candidates to tmpArr
		tmpArr := (c.candidateCacheArr)[c.maxCount:]
		// Reserve elected candidates
		c.candidateCacheArr = (c.candidateCacheArr)[:c.maxCount]

		newimmediateIds := make([]discover.NodeID, 0)
		for _, can := range c.candidateCacheArr {
			newimmediateIds = append(newimmediateIds, can.CandidateId)
		}
		// handle tmpArr
		for _, tmpCan := range tmpArr {
			// delete the lost candidates from immediate elected candidates of trie
			if err := c.delImmediate(state, tmpCan.CandidateId); nil != err {
				return err
			}
			// append to refunds (defeat) trie
			if err := c.setDefeat(state, tmpCan.CandidateId, tmpCan); nil != err {
				return err
			}
		}

		// update immediate index
		if err := c.setImmediateIndex(state, newimmediateIds); nil != err {
			log.Error("Failed to encode immediate ids on SetCandidate err", err)
			return err
		}

		// update index of refund (defeat) on trie
		if err := c.setDefeatIndex(state); nil != err {
			return err
		}
	}

	// cache id
	sortIds := make([]discover.NodeID, 0)

	// insert elected candidate to tire
	for _, can := range c.candidateCacheArr {
		c.setImmediate(state, can.CandidateId, can)
		sortIds = append(sortIds, can.CandidateId)
	}
	// update index of immediate elected candidates on trie
	c.setImmediateIndex(state, sortIds)
	return nil
}

// Getting immediate candidate info by nodeId
func (c *CandidatePool) GetCandidate(state vm.StateDB, nodeId discover.NodeID) (*types.Candidate, error) {
	return c.getCandidate(state, nodeId)
}

// candidate withdraw from immediates elected candidates
func (c *CandidatePool) WithdrawCandidate(state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error {
	log.Info("WithdrawCandidate...")
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 2); nil != err {
		log.Error("Failed to initDataByState on WithdrawCandidate err", err)
		return err
	}

	if price.Cmp(new(big.Int).SetUint64(0)) <= 0 {
		log.Error("withdraw failed price invalid, price", price.String())
		return WithdrawPriceErr
	}
	can, ok := c.immediateCandates[nodeId]
	if !ok || nil == can {
		log.Error("withdraw failed current Candidate is empty")
		return CandidateEmptyErr
	}
	//if nil == can {
	//	log.Error("withdraw failed current Candidate is empty")
	//	return CandidateEmptyErr
	//}

	// check withdraw price
	if can.Deposit.Cmp(price) < 0 {
		log.Error("withdraw failed refund price must less or equal deposit", "key", nodeId.String())
		return WithdrawPriceErr
	} else if can.Deposit.Cmp(price) == 0 { // full withdraw
		// delete current candidate from immediate elected candidates
		if err := c.delImmediate(state, nodeId); nil != err {
			return err
		}
		// update immediate id index
		if ids, err := c.getImmediateIndex(state); nil != err {
			return err
		} else {
			for i, id := range ids {
				if id == nodeId {
					ids = append(ids[:i], ids[i+1:]...)
				}
			}
			if err := c.setImmediateIndex(state, ids); nil != err {
				return err
			}
		}

		// append to refund (defeat) trie
		if err := c.setDefeat(state, nodeId, can); nil != err {
			return err
		}
		// update index of defeat on trie
		if err := c.setDefeatIndex(state); nil != err {
			return err
		}
	} else {
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
		}

		// update current candidate
		if err := c.setImmediate(state, nodeId, canNew); nil != err {
			return err
		}

		// sort immediate
		c.candidateCacheArr = make([]*types.Candidate, 0)
		for _, can := range c.immediateCandates {
			c.candidateCacheArr = append(c.candidateCacheArr, can)
		}
		candidateSort(c.candidateCacheArr)
		ids := make([]discover.NodeID, 0)
		for _, can := range c.candidateCacheArr {
			ids = append(ids, can.CandidateId)
		}
		// update new index
		if err := c.setImmediateIndex(state, ids); nil != err {
			return err
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
		}
		// the withdraw
		if err := c.setDefeat(state, nodeId, canDefeat); nil != err {
			return err
		}
		// update index of defeat on trie
		if err := c.setDefeatIndex(state); nil != err {
			return err
		}
	}
	return nil
}

// Getting all immediate elected candidates array
func (c *CandidatePool) GetChosens(state vm.StateDB) []*types.Candidate {
	log.Info("获取实时入围列表...")
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 1); nil != err {
		log.Error("Failed to initDataByState on WithdrawCandidate err", err)
		return nil
	}
	immediateIds, err := c.getImmediateIndex(state)
	if nil != err {
		log.Error("Failed to getImmediateIndex err", err)
		return nil
	}
	arr := make([]*types.Candidate, 0)
	for _, id := range immediateIds {
		arr = append(arr, c.immediateCandates[id])
	}
	return arr
}

// Getting all witness array
func (c *CandidatePool) GetChairpersons(state vm.StateDB) []*types.Candidate {
	log.Info("获取本轮见证人列表...")
	c.lock.RLock()
	defer c.lock.RUnlock()
	if err := c.initDataByState(state, 0); nil != err {
		log.Error("Failed to initDataByState on GetChairpersons err", err)
		return nil
	}
	witnessIds, err := c.getWitnessIndex(state)
	if nil != err {
		log.Error("Failed to getWitnessIndex on GetChairpersonserr", err)
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
		log.Error("Failed to initDataByState on GetDefeat err", err)
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
		log.Error("Failed to initDataByState on IsDefeat err", err)
		return false, err
	}

	if _, ok := c.immediateCandates[nodeId]; ok {
		log.Error("Candidate is empty")
		return false, nil
	}

	if arr, ok := c.defeatCandidates[nodeId]; ok && len(arr) != 0 {
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
		log.Error("Failed to initDataByState on GetOwner err", err)
		return common.Address{}
	}
	pre_can, pre_ok := c.preOriginCandidates[nodeId]
	or_can, or_ok := c.originCandidates[nodeId]
	ne_can, ne_ok := c.nextOriginCandidates[nodeId]
	im_can, im_ok := c.immediateCandates[nodeId]
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
		log.Error("Failed to initDataByState on RefundBalance err", err)
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
	var amount uint64
	// Transfer refund information that needs to be deleted
	delCanArr := make([]*types.Candidate, 0)

	contractBalance := state.GetBalance(common.CandidateAddr)
	//currentNum := new(big.Int).SetUint64(blockNumber)

	// Traverse all refund information belong to this nodeId
	for index, can := range canArr {
		sub := new(big.Int).Sub(blockNumber, can.BlockNumber)
		log.Info("当前块高:", blockNumber.String(), "质押块高:", can.BlockNumber.String(), "相差:", sub.String())
		log.Info("当前nodeId:", can.CandidateId.String())
		if sub.Cmp(new(big.Int).SetUint64(c.RefundBlockNumber)) >= 0 { // allow refund
			delCanArr = append(delCanArr, can)
			canArr = append(canArr[:index], canArr[index+1:]...)
			// add up the refund price
			amount += can.Deposit.Uint64()
		} else {
			log.Error("block height number had mismatch, No refunds allowed, current block height", blockNumber.String(), "deposit block height", can.BlockNumber.String(), "allowed block interval", c.RefundBlockNumber)
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
		if (contractBalance.Cmp(new(big.Int).SetUint64(amount))) < 0 {
			log.Error("Failed to refundbalance constract account insufficient balance ", state.GetBalance(common.CandidateAddr).String(), "amount", amount)
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
		//delete(c.defeatCandidates, nodeId)
		if err := c.delDefeat(state, nodeId); nil != err {
			log.Error("RefundBalance failed to delDefeat err", err)
			return err
		}
		if ids, err := getDefeatIdsByState(state); nil != err {
			for i, id := range ids {
				if id == nodeId {
					ids = append(ids[:i], ids[i+1:]...)
				}
			}
			if value, err := rlp.EncodeToBytes(&ids); nil != err {
				log.Error("Failed to encode candidate ids on RefundBalance err", err)
				return CandidateEncodeErr
			} else {
				setDefeatIdsState(state, value)
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
	// sub contract account balance
	state.SubBalance(common.CandidateAddr, new(big.Int).SetUint64(amount))
	// add owner balace
	state.AddBalance(addr, new(big.Int).SetUint64(amount))
	log.Info("一键退款完成...")
	return nil
}

// set immediate candidate extra value
func (c *CandidatePool) SetCandidateExtra(state vm.StateDB, nodeId discover.NodeID, extra string) error {
	log.Info("设置推展信息: nodeId = " + nodeId.String())
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 1); nil != err {
		log.Error("Failed to initDataByState on SetCandidateExtra err", err)
		return err
	}
	if can, ok := c.immediateCandates[nodeId]; ok {
		// update current candidate info and update to tire
		can.Extra = extra
		if err := c.setImmediate(state, nodeId, can); nil != err {
			return err
		}
	} else {
		return CandidateEmptyErr
	}
	return nil
}

// Announce witness
func (c *CandidatePool) Election(state *state.StateDB) ([]*discover.Node, error) {
	log.Info("揭榜...")
	log.Info("揭榜 candidate：%+v", *c)
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 1); nil != err {
		log.Error("Failed to initDataByState on Election err", err)
		return nil, err
	}

	// sort immediate candidates
	candidateSort(c.candidateCacheArr)
	PrintObject("揭榜时，排序的数组:", c.candidateCacheArr)
	// cache ids
	immediateIds := make([]discover.NodeID, 0)
	for _, can := range c.candidateCacheArr {
		immediateIds = append(immediateIds, can.CandidateId)
	}
	PrintObject("当前入围者ids：", immediateIds)
	log.Info("当前配置的允许见证人个数:", c.maxChair)
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
	PrintObject("选出来的下一轮见证人Ids:", nextWitIds)
	// cache map
	nextWits := make(map[discover.NodeID]*types.Candidate, 0)

	// copy witnesses information
	copyCandidateMapByIds(nextWits, c.immediateCandates, nextWitIds)
	PrintObject("从入围信息copy过来的见证人;", nextWits)
	// clear all old nextwitnesses information （If it is forked, the next round is no empty.）
	for nodeId, _ := range c.nextOriginCandidates {
		if err := c.delNextWitness(state, nodeId); nil != err {
			log.Error("failed to delNextWitness on election err", err)
			return nil, err
		}
	}

	// set up all new nextwitnesses information
	for nodeId, can := range nextWits {
		if err := c.setNextWitness(state, nodeId, can); nil != err {
			log.Error("failed to setNextWitness on election err", err)
			return nil, err
		}
	}
	// update new nextwitnesses index
	if err := c.setNextWitnessIndex(state, nextWitIds); nil != err {
		log.Error("failed to setNextWitnessIndex on election err", err)
		return nil, err
	}
	// replace the next round of witnesses
	c.nextOriginCandidates = nextWits
	arr := make([]*discover.Node, 0)
	for _, can := range nextWits {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build Node on GetWitness err", err, "nodeId", can.CandidateId.String())
			continue
		} else {
			arr = append(arr, node)
		}
	}
	PrintObject("下一轮见证人node信息:", arr)
	log.Info("揭榜完成...")
	return arr, nil
}

// switch next witnesses to current witnesses
func (c *CandidatePool) Switch(state *state.StateDB) bool {
	log.Info("替换见证人...")
	c.lock.Lock()
	defer c.lock.Unlock()
	if err := c.initDataByState(state, 0); nil != err {
		log.Error("Failed to initDataByState on Switch err", err)
		return false
	}
	// clear all old previous witness on trie
	for nodeId, _ := range c.preOriginCandidates {
		if err := c.delPreviousWitness(state, nodeId); nil != err {
			log.Error("Failed to delPreviousWitness on Switch err", err)
			return false
		}
	}
	// set up new witnesses to previous witnesses on trie by current witnesses
	for nodeId, can := range c.originCandidates {
		if err := c.setPreviousWitness(state, nodeId, can); nil != err {
			log.Error("Failed to setPreviousWitness on Switch err", err)
			return false
		}

	}
	// update previous witness index by current witness index
	if ids, err := c.getWitnessIndex(state); nil != err {
		log.Error("Failed to getWitnessIndex on Switch err", err)
		return false
	} else {
		// replace witnesses index
		if err := c.setPreviousWitnessindex(state, ids); nil != err {
			log.Error("Failed to setPreviousWitnessindex on Switch err", err)
			return false
		}
	}

	// clear all old witnesses on trie
	for nodeId, _ := range c.originCandidates {
		if err := c.delWitness(state, nodeId); nil != err {
			log.Error("Failed to delWitness on Switch err", err)
			return false
		}
	}
	// set up new witnesses to current witnesses on trie by next witnesses
	for nodeId, can := range c.nextOriginCandidates {
		if err := c.setWitness(state, nodeId, can); nil != err {
			log.Error("Failed to setWitness on Switch err", err)
			return false
		}
	}
	// update current witness index by next witness index
	if ids, err := c.getNextWitnessIndex(state); nil != err {
		log.Error("Failed to getNextWitnessIndex on Switch err", err)
		return false
	} else {
		// replace witnesses index
		if err := c.setWitnessindex(state, ids); nil != err {
			log.Error("Failed to setWitnessindex on Switch err", err)
			return false
		}
	}
	// clear all old nextwitnesses information
	for nodeId, _ := range c.nextOriginCandidates {
		if err := c.delNextWitness(state, nodeId); nil != err {
			log.Error("failed to delNextWitness on election err", err)
			return false
		}
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
		log.Error("Failed to initDataByState on GetWitness err", err)
		return nil, err
	}
	//var ids []discover.NodeID
	var witness map[discover.NodeID]*types.Candidate
	if flag == -1 {
		//prewitnessIds, err := c.getPreviousWitnessIndex(state)
		//if nil != err {
		//	log.Error("Failed to getPreviousWitnessIndex on GetWitness err", err)
		//	return nil, err
		//}
		//ids = prewitnessIds
		witness = c.preOriginCandidates
	} else if flag == 0 {
		//witnessIds, err := c.getWitnessIndex(state)
		//if nil != err {
		//	log.Error("Failed to getWitnessIndex on GetWitness err", err)
		//	return nil, err
		//}
		//ids = witnessIds
		witness = c.originCandidates
	} else if flag == 1 {
		//nextWitnessIds, err := c.getNextWitnessIndex(state)
		//if nil != err {
		//	log.Error("Failed to getNextWitnessIndex on GetWitness err", err)
		//	return nil, err
		//}
		//ids = nextWitnessIds
		witness = c.nextOriginCandidates
	}

	arr := make([]*discover.Node, 0)
	for _, can := range witness {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build Node on GetWitness err", err, "nodeId", can.CandidateId.String())
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
		log.Error("Failed to initDataByState on GetAllWitness err", err)
		return nil, nil, nil, err
	}
	//var ids []discover.NodeID
	var prewitness, witness, nextwitness map[discover.NodeID]*types.Candidate
	prewitness = c.preOriginCandidates
	witness = c.originCandidates
	nextwitness = c.nextOriginCandidates

	preArr, curArr, nextArr := make([]*discover.Node, 0), make([]*discover.Node, 0), make([]*discover.Node, 0)
	for _, can := range prewitness {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build pre Node on GetAllWitness err", err, "nodeId", can.CandidateId.String())
			//continue
			return nil, nil, nil, err
		} else {
			preArr = append(preArr, node)
		}
	}
	for _, can := range witness {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build cur Node on GetAllWitness err", err, "nodeId", can.CandidateId.String())
			//continue
			return nil, nil, nil, err
		} else {
			curArr = append(curArr, node)
		}
	}
	for _, can := range nextwitness {
		if node, err := buildWitnessNode(can); nil != err {
			log.Error("Failed to build next Node on GetAllWitness err", err, "nodeId", can.CandidateId.String())
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

func (c *CandidatePool) setImmediate(state vm.StateDB, candidateId discover.NodeID, can *types.Candidate /*, isADD bool*/) error {
	c.immediateCandates[candidateId] = can
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
func (c *CandidatePool) delImmediate(state vm.StateDB, candidateId discover.NodeID /*, isDEL bool*/) error {

	// deleted immediate candidate by id on trie
	setImmediateState(state, candidateId, []byte{})
	// deleted immedidate candidate by id on map
	delete(c.immediateCandates, candidateId)
	//// delete the corresponding id in the index
	//var canIds []discover.NodeID
	//if ids, err := getImmediateIdsByState(state); nil != err {
	//	log.Error("Failed to decode ImmediateIds err", err)
	//	return err
	//}else {
	//	canIds = ids
	//}
	//
	//var flag bool
	//for i, id := range canIds {
	//	if id == candidateId {
	//		flag = true
	//		canIds = append(canIds[:i], canIds[i+1:]...)
	//	}
	//}
	//if flag {
	//	if val, err := rlp.EncodeToBytes(canIds); nil != err {
	//		log.Error("Failed to encode ImmediateIds err", err)
	//		return err
	//	}else {
	//		setImmediateIdsState(state, val)
	//	}
	//}
	return nil
}

func (c *CandidatePool) setImmediateIndex(state vm.StateDB, nodeIds []discover.NodeID) error {
	if val, err := rlp.EncodeToBytes(nodeIds); nil != err {
		log.Error("Failed to encode ImmediateIds err", err)
		return err
	} else {
		setImmediateIdsState(state, val)
	}
	return nil
}

// setting refund information
func (c *CandidatePool) setDefeat(state vm.StateDB, candidateId discover.NodeID, can *types.Candidate) error {

	var defeatArr []*types.Candidate
	// append refund information
	if defeatArrTmp, ok := c.defeatCandidates[can.CandidateId]; ok {
		defeatArrTmp = append(defeatArrTmp, can)
		c.defeatCandidates[can.CandidateId] = defeatArrTmp
		defeatArr = defeatArrTmp
	} else {
		defeatArrTmp = make([]*types.Candidate, 0)
		defeatArrTmp = append(defeatArr, can)
		c.defeatCandidates[can.CandidateId] = defeatArrTmp
		defeatArr = defeatArrTmp
	}
	// setting refund information on trie
	if value, err := rlp.EncodeToBytes(&defeatArr); nil != err {
		log.Error("Failed to encode candidate object on setDefeat", "key", candidateId.String(), "err", err)
		return CandidateEncodeErr
	} else {
		setDefeatState(state, candidateId, value)
	}
	return nil
}

func (c *CandidatePool) delDefeat(state vm.StateDB, nodeId discover.NodeID) error {
	delete(c.defeatCandidates, nodeId)
	setDefeatState(state, nodeId, []byte{})

	//// delete the corresponding id in the index
	//var canIds []discover.NodeID
	//if ids, err := getDefeatIdsByState(state); nil != err {
	//	log.Error("Failed to decode DefeatIds err", err)
	//	return err
	//}else {
	//	canIds = ids
	//}
	//
	//var flag bool
	//for i, id := range canIds {
	//	if id == nodeId {
	//		flag = true
	//		canIds = append(canIds[:i], canIds[i+1:]...)
	//	}
	//}
	//if flag {
	//	if val, err := rlp.EncodeToBytes(canIds); nil != err {
	//		log.Error("Failed to encode ImmediateIds err", err)
	//		return err
	//	}else {
	//		setDefeatIdsState(state, val)
	//	}
	//}
	return nil
}

// update refund index
func (c *CandidatePool) setDefeatIndex(state vm.StateDB) error {
	newdefeatIds := make([]discover.NodeID, 0)
	for id, _ := range c.defeatCandidates {
		newdefeatIds = append(newdefeatIds, id)
	}
	if value, err := rlp.EncodeToBytes(&newdefeatIds); nil != err {
		log.Error("Failed to encode candidate object on setDefeatIds err", err)
		return CandidateEncodeErr
	} else {
		setDefeatIdsState(state, value)
	}
	return nil
}

func (c *CandidatePool) delPreviousWitness(state vm.StateDB, candidateId discover.NodeID) error {
	// deleted previous witness by id on map
	delete(c.preOriginCandidates, candidateId)
	// delete previous witness by id on trie
	setPreviousWitnessState(state, candidateId, []byte{})
	//// delete the corresponding id in the index
	//var canIds []discover.NodeID
	//if ids, err := getPreviousWitnessIdsState(state); nil != err {
	//	log.Error("Failed to decode PreviousWitnessIds err", err)
	//	return err
	//}else {
	//	canIds = ids
	//}
	//
	//var flag bool
	//for i, id := range canIds {
	//	if id == candidateId {
	//		flag = true
	//		canIds = append(canIds[:i], canIds[i+1:]...)
	//	}
	//}
	//if flag {
	//	if arrVal, err := rlp.EncodeToBytes(canIds); nil != err {
	//		log.Error("Failed to encode PreviousWitnessIds err", err)
	//		return err
	//	}else {
	//		setPreviosWitnessIdsState(state, arrVal)
	//	}
	//}
	return nil
}

func (c *CandidatePool) setPreviousWitness(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	c.preOriginCandidates[nodeId] = can
	if val, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode Candidate on setPreviousWitness err", err)
		return err
	} else {
		setPreviousWitnessState(state, nodeId, val)
	}
	return nil
}

func (c *CandidatePool) setPreviousWitnessindex(state vm.StateDB, nodeIds []discover.NodeID) error {
	if val, err := rlp.EncodeToBytes(nodeIds); nil != err {
		log.Error("Failed to encode Previous WitnessIds err", err)
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
	PrintObject("把下一轮设置本轮", can)
	c.originCandidates[nodeId] = can
	if val, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode Candidate on setWitness err", err)
		return err
	} else {
		setWitnessState(state, nodeId, val)
	}
	return nil
}

func (c *CandidatePool) setWitnessindex(state vm.StateDB, nodeIds []discover.NodeID) error {
	if val, err := rlp.EncodeToBytes(nodeIds); nil != err {
		log.Error("Failed to encode WitnessIds err", err)
		return err
	} else {
		setWitnessIdsState(state, val)
	}
	return nil
}

func (c *CandidatePool) delWitness(state vm.StateDB, candidateId discover.NodeID) error {
	// deleted witness by id on map
	delete(c.originCandidates, candidateId)
	// delete witness by id on trie
	setWitnessState(state, candidateId, []byte{})
	//// delete the corresponding id in the index
	//var canIds []discover.NodeID
	//if ids, err := getWitnessIdsByState(state); nil != err {
	//	log.Error("Failed to decode WitnessIds err", err)
	//	return err
	//}else {
	//	canIds = ids
	//}
	//
	//var flag bool
	//for i, id := range canIds {
	//	if id == candidateId {
	//		flag = true
	//		canIds = append(canIds[:i], canIds[i+1:]...)
	//	}
	//}
	//if flag {
	//	if arrVal, err := rlp.EncodeToBytes(canIds); nil != err {
	//		log.Error("Failed to encode WitnessIds err", err)
	//		return err
	//	}else {
	//		setWitnessIdsState(state, arrVal)
	//	}
	//}
	return nil
}

func (c *CandidatePool) getWitnessIndex(state vm.StateDB) ([]discover.NodeID, error) {
	return getWitnessIdsByState(state)
}

func (c *CandidatePool) setNextWitness(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	PrintObject("设置 setNextWitness", *can)
	c.nextOriginCandidates[nodeId] = can
	if value, err := rlp.EncodeToBytes(can); nil != err {
		log.Error("Failed to encode candidate object on setImmediate", "key", nodeId.String(), "err", err)
		return CandidateEncodeErr
	} else {
		// setting next witness information on trie
		setNextWitnessState(state, nodeId, value)
	}
	return nil
}

func (c *CandidatePool) delNextWitness(state vm.StateDB, candidateId discover.NodeID) error {
	// deleted next witness by id on map
	delete(c.nextOriginCandidates, candidateId)
	// deleted next witness by id on trie
	setNextWitnessState(state, candidateId, []byte{})

	//// getting origin index
	//var canIds []discover.NodeID
	//if ids, err := getNextWitnessIdsByState(state); nil != err {
	//	log.Error("Failed to decode NextWitnessIds err", err)
	//	return err
	//}else {
	//	canIds = ids
	//}
	//
	//// delete the corresponding id in the index
	//var flag bool
	//for i, id := range canIds {
	//	if id == candidateId {
	//		flag = true
	//		canIds = append(canIds[:i], canIds[i+1:]...)
	//	}
	//}
	//if flag {
	//	if arrVal, err := rlp.EncodeToBytes(canIds); nil != err {
	//		log.Error("Failed to encode NextWitnessIds err", err)
	//		return err
	//	}else {
	//		setNextWitnessIdsState(state, arrVal)
	//	}
	//}
	return nil
}

func (c *CandidatePool) setNextWitnessIndex(state vm.StateDB, nodeIds []discover.NodeID) error {
	PrintObject("设置 setNextWitnessIndex:", nodeIds)
	if value, err := rlp.EncodeToBytes(&nodeIds); nil != err {
		log.Error("Failed to encode candidate object on setDefeatIds err", err)
		return CandidateEncodeErr
	} else {
		setNextWitnessIdsState(state, value)
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
		log.Error("Failed to initDataByState on getCandidate err", err)
		return nil, err
	}
	if candidatePtr, ok := c.immediateCandates[nodeId]; ok {
		PrintObject("GetCandidate 返回：", candidatePtr)
		return candidatePtr, nil
	}
	return nil, nil
}

func (c *CandidatePool) MaxChair() uint64 {
	return c.maxChair
}

func getPreviousWitnessIdsState(state vm.StateDB) ([]discover.NodeID, error) {
	var witnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidateAddr, PreviousWitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &witnessIds); nil != err {
			return nil, err
		}
	}
	return witnessIds, nil
}

func setPreviosWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidateAddr, PreviousWitnessListKey(), arrVal)
}

func getPreviousWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if err := rlp.DecodeBytes(state.GetState(common.CandidateAddr, PreviousWitnessKey(id)), &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func setPreviousWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidateAddr, PreviousWitnessKey(id), val)
}

func getWitnessIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var witnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidateAddr, WitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &witnessIds); nil != err {
			return nil, err
		}
	}
	return witnessIds, nil
}

func setWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidateAddr, WitnessListKey(), arrVal)
}

func getWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if err := rlp.DecodeBytes(state.GetState(common.CandidateAddr, WitnessKey(id)), &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func setWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidateAddr, WitnessKey(id), val)
}

func getNextWitnessIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var nextWitnessIds []discover.NodeID
	if valByte := state.GetState(common.CandidateAddr, NextWitnessListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &nextWitnessIds); nil != err {
			return nil, err
		}
	}
	return nextWitnessIds, nil
}

func setNextWitnessIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidateAddr, NextWitnessListKey(), arrVal)
}

func getNextWitnessByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if err := rlp.DecodeBytes(state.GetState(common.CandidateAddr, NextWitnessKey(id)), &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func setNextWitnessState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidateAddr, NextWitnessKey(id), val)
}

func getImmediateIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var immediateIds []discover.NodeID
	if valByte := state.GetState(common.CandidateAddr, ImmediateListKey()); len(valByte) != 0 {
		if err := rlp.DecodeBytes(valByte, &immediateIds); nil != err {
			return nil, err
		}
	}
	return immediateIds, nil
}

func setImmediateIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidateAddr, ImmediateListKey(), arrVal)
}

func getImmediateByState(state vm.StateDB, id discover.NodeID) (*types.Candidate, error) {
	var can types.Candidate
	if err := rlp.DecodeBytes(state.GetState(common.CandidateAddr, ImmediateKey(id)), &can); nil != err {
		return nil, err
	}
	return &can, nil
}

func setImmediateState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidateAddr, ImmediateKey(id), val)
}

func getDefeatIdsByState(state vm.StateDB) ([]discover.NodeID, error) {
	var defeatIds []discover.NodeID
	if valBtye := state.GetState(common.CandidateAddr, DefeatListKey()); len(valBtye) != 0 {
		if err := rlp.DecodeBytes(valBtye, &defeatIds); nil != err {
			return nil, err
		}
	}
	return defeatIds, nil
}

func setDefeatIdsState(state vm.StateDB, arrVal []byte) {
	state.SetState(common.CandidateAddr, DefeatListKey(), arrVal)
}

func getDefeatsByState(state vm.StateDB, id discover.NodeID) ([]*types.Candidate, error) {
	var canArr []*types.Candidate
	if err := rlp.DecodeBytes(state.GetState(common.CandidateAddr, DefeatKey(id)), &canArr); nil != err {
		return nil, err
	}
	return canArr, nil
}

func setDefeatState(state vm.StateDB, id discover.NodeID, val []byte) {
	state.SetState(common.CandidateAddr, DefeatKey(id), val)
}

func copyCandidateMapByIds(target, source map[discover.NodeID]*types.Candidate, ids []discover.NodeID) {
	for _, id := range ids {
		target[id] = source[id]
	}
}

func GetCandidatePtr() *CandidatePool {
	return candidatePool
}

func PrintObject(s string, obj interface{}) {
	objs, _ := json.Marshal(obj)
	log.Info(s, string(objs), "\n")
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
	//key, _ := rlp.EncodeToBytes(nodeId)
	return immediateKey(nodeId.Bytes())
}
func immediateKey(key []byte) []byte {
	return append(ImmediateBtyePrefix, key...)
}

func PreviousWitnessKey(nodeId discover.NodeID) []byte {
	return prewitnessKey(nodeId.Bytes())
}

func prewitnessKey(key []byte) []byte {
	return append(PreWitnessBytePrefix, key...)
}

func WitnessKey(nodeId discover.NodeID) []byte {
	//key, _ := rlp.EncodeToBytes(nodeId)
	return witnessKey(nodeId.Bytes())
}
func witnessKey(key []byte) []byte {
	return append(WitnessBtyePrefix, key...)
}

func NextWitnessKey(nodeId discover.NodeID) []byte {
	return nextWitnessKey(nodeId.Bytes())
}
func nextWitnessKey(key []byte) []byte {
	return append(NextWitnessBtyePrefix, key...)
}

func DefeatKey(nodeId discover.NodeID) []byte {
	//key, _ := rlp.EncodeToBytes(nodeId)
	return defeatKey(nodeId.Bytes())
}
func defeatKey(key []byte) []byte {
	return append(DefeatBtyePrefix, key...)
}

func ImmediateListKey() []byte {
	return ImmediateListBtyePrefix
}

func PreviousWitnessListKey() []byte {
	return PreWitnessListBytePrefix
}

func WitnessListKey() []byte {
	return WitnessListBtyePrefix
}

func NextWitnessListKey() []byte {
	return NextWitnessListBytePrefix
}

func DefeatListKey() []byte {
	return DefeatListBtyePrefix
}
