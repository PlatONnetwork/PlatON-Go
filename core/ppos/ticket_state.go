package pposm

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/core/types"
	"Platon-go/core/vm"
	"Platon-go/crypto/sha3"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"Platon-go/rlp"
	"errors"
	"github.com/satori/go.uuid"
	"math/big"
	"sort"
	"sync"
)

var (

	CandidateNotFindErr			= errors.New("The node has lost its candidacy")
	TicketNilErr				= errors.New("Ticket Insufficient quantity")
	EncodeTicketErr				= errors.New("Encode Ticket error")
	EncodePoolNumberErr			= errors.New("Encode SurplusQuantity error")
	DecodeTicketErr				= errors.New("Decode Ticket error")
	DecodePoolNumberErr			= errors.New("Decode SurplusQuantity error")
	RecordExpireTicketErr		= errors.New("Record Expire Ticket error")
	CandidateNotFindTicketErr	= errors.New("The candidate no longer has this ticket")
	GetCandidateTicketIdErr		= errors.New("Get Candidate TicketIds error")
	SetCandidateTicketIdErr		= errors.New("Update Candidate TicketIds error")
	TicketPoolBalanceErr		= errors.New("TicketPool not sufficient funds")
	GetOwnerTicketIdsErr		= errors.New("Get Owner TicketIds error")
	SetOwnerTicketIdsErr		= errors.New("Update Owner TicketIds error")
	TicketIdNotFindErr			= errors.New("TicketId not find")
	HandleExpireTicketErr		= errors.New("Failure to deal with expired tickets")
	GetCandidateAttachErr		= errors.New("Get CandidateAttach error")
	SetCandidateAttachErr		= errors.New("Update CandidateAttach error")
)

type TicketPool struct {
	// Maximum number of ticket pool
	MaxCount			uint64
	// Remaining number of ticket pool
	SurplusQuantity		uint64
	// Overdue
	ExpireBlockNumber	uint64
	lock				*sync.RWMutex
}

var ticketPool *TicketPool

// initialize the global ticket pool object
func NewTicketPool(configs *params.PposConfig) *TicketPool {
	ticketPool = &TicketPool{
		MaxCount:				configs.TicketConfig.MaxCount,
		SurplusQuantity:		configs.TicketConfig.MaxCount,
		ExpireBlockNumber:		configs.TicketConfig.ExpireBlockNumber,
		lock:					&sync.RWMutex{},
	}
	return ticketPool
}

func(t *TicketPool) VoteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error) {
	voteTicketIdList, err := t.voteTicket(stateDB, owner, voteNumber, deposit, nodeId, blockNumber)
	if nil != err {
		return voteTicketIdList, err
	}
	// 调用候选人重新排序接口
	candidatePool.UpdateElectedQueue(stateDB, nodeId)
	return voteTicketIdList, nil
}

func(t *TicketPool) voteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	voteTicketIdList := make([]common.Hash, 0)
	// check ticket pool count
	t.GetPoolNumber(stateDB)
	if t.SurplusQuantity == 0 {
		log.Error("Ticket Insufficient quantity", TicketNilErr)
		return voteTicketIdList, TicketNilErr
	}
	if t.SurplusQuantity < voteNumber {
		voteNumber -= t.SurplusQuantity
	}
	candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
	if nil != err {
		return voteTicketIdList, err
	}
	ownerTicketIds, err := t.GetOwnerNormalTicketIds(stateDB, owner)
	if nil != err {
		return voteTicketIdList, err
	}
	var i uint64 = 0
	for ; i < voteNumber; i++ {
		ticketId, err := generateTicketId()
		if err != nil {
			return voteTicketIdList, err
		}
		ticket := &types.Ticket{
			TicketId:		ticketId,
			Owner:			owner,
			Deposit:		deposit,
			CandidateId:	nodeId,
			BlockNumber:	blockNumber,
		}
		voteTicketIdList = append(voteTicketIdList, ticketId)
		candidateTicketIds = append(candidateTicketIds, ticketId)
		ownerTicketIds = append(ownerTicketIds, ticketId)
		if err := t.setTicket(stateDB, ticketId, ticket); err != nil {
			return voteTicketIdList, err
		}
		if err := t.recordExpireTicket(stateDB, blockNumber, ticketId); err != nil {
			return voteTicketIdList, err
		}
		if err := t.subPoolNumber(stateDB); err != nil {
			return voteTicketIdList, err
		}
	}
	// 设置购票人的购票关联信息
	if err := t.setOwnerNormalTicketIds(stateDB, owner, ownerTicketIds); err != nil {
		return voteTicketIdList, err
	}
	if err := t.setCandidateTicketIds(stateDB, nodeId, candidateTicketIds); err != nil {
		return voteTicketIdList, err
	}
	candidateAttach := new(types.CandidateAttach)
	if err := t.setCandidateAttach(stateDB, nodeId, candidateAttach); nil != err {
		return voteTicketIdList, err
	}
	return voteTicketIdList, nil
}

func setOwnerTicketIds(stateDB vm.StateDB, key []byte, ticketIds []common.Hash) error {
	if value, err := rlp.EncodeToBytes(ticketIds); nil != err {
		log.Error("setOwnerTicketIds error", "key", string(key), "err", err)
		return SetOwnerTicketIdsErr
	} else {
		setTicketPoolState(stateDB, key, value)
	}
	return nil
}

func (t *TicketPool) setOwnerExpireTicketIds(stateDB vm.StateDB, owner common.Address, ticketIds []common.Hash) error {
	return setOwnerTicketIds(stateDB, AccountExpireTicketIdsKey(owner.Bytes()), ticketIds)
}

func (t *TicketPool) setOwnerNormalTicketIds(stateDB vm.StateDB, owner common.Address, ticketIds []common.Hash) error {
	return setOwnerTicketIds(stateDB, AccountNormalTicketIdsKey(owner.Bytes()), ticketIds)
}

func getOwnerTicketIds(stateDB vm.StateDB, key []byte) ([]common.Hash, error) {
	var ticketIds []common.Hash
	if err := getTicketPoolState(stateDB, key, &ticketIds); nil != err {
		log.Error("getOwnerTicketIds error", "key", string(key))
		return nil, GetOwnerTicketIdsErr
	}
	return ticketIds, nil
}

func (t *TicketPool) GetOwnerExpireTicketIds(stateDB vm.StateDB, owner common.Address) ([]common.Hash, error) {
	return getOwnerTicketIds(stateDB, AccountExpireTicketIdsKey(owner.Bytes()))
}

func (t *TicketPool) GetOwnerNormalTicketIds(stateDB vm.StateDB, owner common.Address) ([]common.Hash, error) {
	return getOwnerTicketIds(stateDB, AccountNormalTicketIdsKey(owner.Bytes()))
}

func (t *TicketPool) calcExpireBlockNumber(stateDB vm.StateDB, blockNumber *big.Int) (*big.Int, bool) {
	num := new(big.Int).SetUint64(0)
	if blockNumber.Cmp(new(big.Int).SetUint64(t.ExpireBlockNumber)) >= 0 {
		num.Sub(blockNumber, new(big.Int).SetUint64(t.ExpireBlockNumber))
		return num, true
	}
	return num, false
}

func (t *TicketPool) GetExpireTicketIds(stateDB vm.StateDB, blockNumber *big.Int) ([]common.Hash, error) {
	var expireTicketIds []common.Hash
	if err := getTicketPoolState(stateDB, ExpireTicketKey((*blockNumber).Bytes()), &expireTicketIds); nil != err {
		return nil, err
	}
	return expireTicketIds, nil
}

// 在当前区块投入的票，则把票id放入，以块高为key的value切片中，以便查找过期票
func (t *TicketPool) recordExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, ticketId common.Hash) error {
	expireTickets, err := t.GetExpireTicketIds(stateDB, blockNumber)
	if err != nil {
		log.Error("recordExpireTicket error", "key", *blockNumber, "ticketId", ticketId.String(), "err", err)
		return RecordExpireTicketErr
	}
	expireTickets = append(expireTickets, ticketId)
	return t.setExpireTicket(stateDB, blockNumber, expireTickets)
}

func (t *TicketPool) setExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, expireTickets []common.Hash) error {
	if value, err := rlp.EncodeToBytes(expireTickets); nil != err {
		log.Error("Failed to encode ticketid object on setExpireTicket", "key", *blockNumber, "err", err)
		return EncodeTicketErr
	} else {
		setTicketPoolState(stateDB, ExpireTicketKey((*blockNumber).Bytes()), value)
	}
	return nil
}

func (t *TicketPool) removeExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, ticketId common.Hash) error {
	ticketIdList, err := t.GetExpireTicketIds(stateDB, blockNumber)
	if err != nil {
		return err
	}
	ticketIdList, success := removeTicketId(ticketId, ticketIdList)
	if !success {
		return TicketIdNotFindErr
	}
	return t.setExpireTicket(stateDB, blockNumber, ticketIdList)
}

func (t *TicketPool) handleExpireTicket(stateDB vm.StateDB, blockNumber *big.Int) ([]discover.NodeID, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	ticketIdList, err := t.GetExpireTicketIds(stateDB, blockNumber)
	if err != nil {
		return nil, err
	}
	candidateAttachMap := make(map[discover.NodeID]*types.CandidateAttach)
	changeNodeIdList := make([]discover.NodeID, 0)
	for _, ticketId := range ticketIdList {
		ticket, err := t.GetTicket(stateDB, ticketId)
		if err != nil {
			return nil, err
		}
		candidateAttach, ok := candidateAttachMap[ticket.CandidateId]
		if !ok {
			tempCandidateAttach, err := t.GetCandidateAttach(stateDB, ticket.CandidateId)
			candidateAttach = tempCandidateAttach
			if nil != err {
				return nil, err
			}
			candidateAttachMap[ticket.CandidateId] = candidateAttach
			changeNodeIdList = append(changeNodeIdList, ticket.CandidateId)
		}
		if _, err := t.releaseTicket(stateDB, ticket.CandidateId, candidateAttach, ticketId, blockNumber); nil != err {
			continue
		}
		// 把过期票放入，购买者的过期票列表中
		ticketIds, err := t.GetOwnerNormalTicketIds(stateDB, ticket.Owner)
		if err != nil {
			return nil, err
		}
		ticketIds, success := removeTicketId(ticketId, ticketIds)
		if success {
			if err := t.setOwnerNormalTicketIds(stateDB, ticket.Owner, ticketIds); nil != err {
				return nil, err
			}
			ticketIds, err := t.GetOwnerExpireTicketIds(stateDB, ticket.Owner)
			if err != nil {
				return nil, err
			}
			ticketIds = append(ticketIds, ticketId)
			if err := t.setOwnerExpireTicketIds(stateDB, ticket.Owner, ticketIds); nil != err {
				return nil, err
			}
		}
	}
	// Update CandidateAttach
	for nodeId, ca := range candidateAttachMap{
		t.setCandidateAttach(stateDB, nodeId, ca)
	}
	return changeNodeIdList, nil
}

// Get ticket list
func (t *TicketPool) GetTicketList(stateDB vm.StateDB, ticketIds []common.Hash) ([]*types.Ticket, error) {
	var tickets []*types.Ticket
	for _, ticketId := range ticketIds {
		ticket, err := t.GetTicket(stateDB, ticketId)
		if nil != err {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

// Get ticket details based on TicketId
func (t *TicketPool) GetTicket(stateDB vm.StateDB, ticketId common.Hash) (*types.Ticket, error) {
	var ticket= new(types.Ticket)
	if err := getTicketPoolState(stateDB, ticketId.Bytes(), ticket); nil != err {
		return nil, DecodeTicketErr
	}
	return ticket, nil
}

func (t *TicketPool) setTicket(stateDB vm.StateDB, ticketId common.Hash, ticket *types.Ticket) error {
	if value, err := rlp.EncodeToBytes(ticket); nil != err {
		log.Error("Failed to encode ticket object on setTicket", "key", ticketId.String(), "err", err)
		return EncodeTicketErr
	} else {
		setTicketPoolState(stateDB, ticketId.Bytes(), value)
	}
	return nil
}

func (t *TicketPool) DropReturnTicket(stateDB vm.StateDB, blockNumber *big.Int, nodeIds ...discover.NodeID) error {
	for _, nodeId := range nodeIds {
		candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
		if nil != err {
			return err
		}
		for _, ticketId := range candidateTicketIds {
			t.ReturnTicket(stateDB, nodeId, ticketId, blockNumber)
		}
	}
	return nil
}

func (t *TicketPool) ReturnTicket(stateDB vm.StateDB, nodeId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	candidateAttach, err := t.GetCandidateAttach(stateDB, nodeId)
	if nil != err {
		return err
	}
	if nil == candidateAttach {
		candidateAttach = new(types.CandidateAttach)
	}
	ticket, err := t.releaseTicket(stateDB, nodeId, candidateAttach, ticketId, blockNumber)
	if err := t.setCandidateAttach(stateDB, nodeId, candidateAttach); nil != err {
		return err
	}
	// 从待过期票中删除
 	return t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
}

func (t *TicketPool) releaseTicket(stateDB vm.StateDB, candidateId discover.NodeID, candidateAttach *types.CandidateAttach, ticketId common.Hash, blockNumber *big.Int) (*types.Ticket, error) {
	ticket, err := t.GetTicket(stateDB, ticketId)
	if nil != err {
		return ticket, err
	}
	candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, candidateId)
	if nil != err {
		return ticket, err
	}
	candidateTicketIds, success := removeTicketId(ticketId, candidateTicketIds)
	if !success {
		log.Error("The candidate no longer has this ticket", "err", err)
		return ticket, CandidateNotFindTicketErr
	}
	if err := t.setCandidateTicketIds(stateDB, candidateId, candidateTicketIds); err != nil {
		return ticket, err
	}
	if err := t.addPoolNumber(stateDB); err != nil {
		return ticket, err
	}
	candidateAttach.SubEpoch(ticket.CalcEpoch(blockNumber))
	return ticket, transfer(stateDB, common.TicketPoolAddr, ticket.Owner, ticket.Deposit)
}

// 1.给幸运票发放奖励
// 2.检查过期票
// 3.增加总票龄
func (t *TicketPool) Notify(stateDB vm.StateDB, blockNumber *big.Int, nodeId discover.NodeID) error {
	// 发放奖励

	// 检查过期票
	expireBlockNumber, ok := t.calcExpireBlockNumber(stateDB, blockNumber)
	if ok {
		if nodeIdList, err := t.handleExpireTicket(stateDB, expireBlockNumber); nil != err {
			log.Error("OutBlockNotice method handleExpireTicket error", "blockNumber", *blockNumber, "err", err)
			return HandleExpireTicketErr
		} else {
			// 处理完过期票之后，通知候选人更新榜单信息
			candidatePool.UpdateElectedQueue(stateDB, nodeIdList...)
		}
	}
	// 每个候选人增加总票龄
	if err := t.calcCandidateEpoch(stateDB); nil != err {
		return err
	}
	return nil
}

func (t *TicketPool) calcCandidateEpoch(stateDB vm.StateDB) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	candidateList := candidatePool.GetChosens(stateDB, 0)
	for _, candidate := range candidateList {
		candidateAttach, err := t.GetCandidateAttach(stateDB, candidate.CandidateId)
		if nil != err {
			return err
		}
		// 获取总票数，增加总票龄
		//candidateAttach.AddEpoch()
		if err := t.setCandidateAttach(stateDB, candidate.CandidateId, candidateAttach); nil != err {
			return err
		}
	}
	return nil
}

func (t *TicketPool) SelectionLuckyTicket(stateDB vm.StateDB, nodeId discover.NodeID, blockHash common.Hash) (common.Hash, error) {
	candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
	luckyTicketId := common.Hash{}
	if nil != err {
		return luckyTicketId, err
	}
	decList := make([]float64, 0)
	decMap := make(map[float64]common.Hash, 0)
	for _, ticketId := range candidateTicketIds {
		decNumber := hexutil.HexDec(ticketId.Hex()[2:])
		decList = append(decList, decNumber)
		decMap[decNumber] = ticketId
	}
	sort.Float64s(decList)
	index := findFirstMatch(decList, hexutil.HexDec(blockHash.Hex()[2:]))
	luckyTicketId = decMap[decList[index]]
	return luckyTicketId, nil
}

func removeTicketId(ticketId common.Hash, ticketIds []common.Hash) ([]common.Hash, bool) {
	for index, tempTicketId := range ticketIds {
		if tempTicketId == ticketId {
			start := ticketIds[:index]
			end := ticketIds[index+1:]
			return append(start, end...), true
		}
	}
	return ticketIds, false
}

func (t *TicketPool) addPoolNumber(stateDB vm.StateDB) error {
	t.SurplusQuantity++
	return t.setPoolNumber(stateDB, t.SurplusQuantity)
}

func (t *TicketPool) subPoolNumber(stateDB vm.StateDB) error {
	t.SurplusQuantity--
	return t.setPoolNumber(stateDB, t.SurplusQuantity)
}

func (t *TicketPool) setPoolNumber(stateDB vm.StateDB, surplusQuantity uint64) error {
	if value, err := rlp.EncodeToBytes(surplusQuantity); nil != err {
		log.Error("Failed to encode surplusQuantity object on setPoolNumber", "key", string(SurplusQuantityKey), "err", err)
		return EncodePoolNumberErr
	} else {
		setTicketPoolState(stateDB, SurplusQuantityKey, value)
	}
	return nil
}

func (t *TicketPool) GetPoolNumber(stateDB vm.StateDB) (uint64, error) {
	var surplusQuantity uint64
	if val := stateDB.GetState(common.TicketPoolAddr, SurplusQuantityKey); len(val) > 0 {
		if err := rlp.DecodeBytes(val, &surplusQuantity); nil != err {
			log.Error("Decode PoolNumber error", "key", string(SurplusQuantityKey), "err", err)
			return surplusQuantity, DecodePoolNumberErr
		}
		t.SurplusQuantity = surplusQuantity
	} else {
		// Default initialization values
		surplusQuantity = t.SurplusQuantity
	}
	return surplusQuantity, nil
}

func (t *TicketPool) GetCandidateTicketIds(stateDB vm.StateDB, nodeId discover.NodeID) ([]common.Hash, error) {
	var ticketIds []common.Hash
	if err := getTicketPoolState(stateDB, nodeId.Bytes(), &ticketIds); nil != err {
		log.Error("get Candidate ticketIds error", "key", string(nodeId.Bytes()), "err", err)
		return nil, GetCandidateTicketIdErr
	}
	return ticketIds, nil
}

func (t *TicketPool) setCandidateTicketIds(stateDB vm.StateDB, nodeId discover.NodeID, ticketIds []common.Hash) error {
	if value, err := rlp.EncodeToBytes(ticketIds); nil != err {
		log.Error("Failed to encode ticketIds object on setCandidateTicketIds", "key", string(nodeId.Bytes()), "err", err)
		return SetCandidateTicketIdErr
	} else {
		setTicketPoolState(stateDB, nodeId.Bytes(), value)
	}
	return nil
}

func (t *TicketPool) GetCandidateAttach(stateDB vm.StateDB, nodeId discover.NodeID) (*types.CandidateAttach, error) {
	candidateAttach := new(types.CandidateAttach)
	if err := getTicketPoolState(stateDB, CandidateAttachKey(nodeId.Bytes()), candidateAttach); nil != err {
		log.Error("GetCandidateAttach error", "key", string(nodeId.Bytes()), "err", err)
		return nil, GetCandidateAttachErr
	}
	return candidateAttach, nil
}

func (t *TicketPool) setCandidateAttach(stateDB vm.StateDB, nodeId discover.NodeID, candidateAttach *types.CandidateAttach) error {
	if value, err := rlp.EncodeToBytes(candidateAttach); nil != err {
		log.Error("Failed to encode candidateAttach object on setCandidateAttach", "key", string(nodeId.Bytes()), "value", *candidateAttach, "err", err)
		return SetCandidateAttachErr
	} else {
		setTicketPoolState(stateDB, CandidateAttachKey(nodeId.Bytes()), value)
	}
	return nil
}

func checkBalance(stateDB vm.StateDB, addr common.Address, amount *big.Int) bool {
	if stateDB.GetBalance(addr).Cmp(amount) >= 0 {
		return true
	}
	return false
}

func transfer(stateDB vm.StateDB, from common.Address, to common.Address, amount *big.Int) error {
	if !checkBalance(stateDB, from, amount) {
		log.Error("TicketPool not sufficient funds", "from", from.Hex(), "to", to.Hex(), "money", *amount)
		return TicketPoolBalanceErr
	}
	stateDB.SubBalance(from, amount)
	stateDB.AddBalance(to, amount)
	return nil
}

func getTicketPoolState(stateDB vm.StateDB, key []byte, result interface{}) error {
	return getState(common.TicketPoolAddr, stateDB, key, result)
}

func getState(addr common.Address, stateDB vm.StateDB, key []byte, result interface{}) error {
	if val := stateDB.GetState(addr, key); len(val) > 0 {
		if err := rlp.DecodeBytes(val, result); nil != err {
			log.Error("Decode Data error", "key", string(key), "err", err)
			return err
		}
	}
	return nil
}

func setTicketPoolState(stateDB vm.StateDB, key []byte, val []byte) {
	stateDB.SetState(common.TicketPoolAddr, key, val)
}

func generateTicketId() (common.Hash, error) {
	// generate ticket id
	uuid, err := uuid.NewV4()
	if err != nil {
		log.Error("generate ticket error", err)
		return common.Hash{}, err
	}
	ticketId := sha3.Sum256(uuid[:])
	return ticketId, nil
}

func ExpireTicketKey(key []byte) []byte {
	return append(ExpireTicketPrefix, key...)
}

func AccountNormalTicketIdsKey(key []byte) []byte {
	return append(AccountNormalTicketPrefix, key...)
}

func AccountExpireTicketIdsKey(key []byte) []byte {
	return append(AccountExpireTicketPrefix, key...)
}

func CandidateAttachKey(key []byte) []byte {
	return append(CandidateAttachPrefix, key...)
}

func findFirstMatch(list []float64, key float64) int {
	left := 0
	right := len(list) - 1
	for left <= right {
		mid := (left + right) / 2
		if list[mid] >= key {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	// 如果找不到匹配的，默认返回最后一个下标
	if left >= len(list)  {
		return len(list) - 1
	}
	return left
}
