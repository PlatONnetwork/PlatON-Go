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
)

type TicketPool struct {
	// Maximum number of ticket pool
	MaxCount			uint64
	// Remaining number of ticket pool
	SurplusQuantity		uint64
	// Overdue
	ExpireBlockNumber	uint64
	// the candidate pool object pointer
	candidatePool		*CandidatePool
	lock				*sync.RWMutex
}

var ticketPool *TicketPool

// initialize the global ticket pool object
func NewTicketPool(configs *params.PposConfig, candidatePool *CandidatePool) *TicketPool {
	ticketPool = &TicketPool{
		MaxCount:				configs.TicketConfig.MaxCount,
		SurplusQuantity:		configs.TicketConfig.MaxCount,
		ExpireBlockNumber:		configs.TicketConfig.ExpireBlockNumber,
		candidatePool:			candidatePool,
		lock:					&sync.RWMutex{},
	}
	return ticketPool
}

func(t *TicketPool) VoteTicket(stateDB vm.StateDB, owner common.Address, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	// check ticket pool count
	if t.SurplusQuantity == 0 {
		log.Error("Ticket Insufficient quantity", TicketNilErr)
		return TicketNilErr
	}
	candidate, err := t.candidatePool.GetCandidate(stateDB, nodeId)
	if err != nil {
		log.Error("GetCandidate error", err)
		return err
	}
	if candidate == nil {
		log.Error("The node has lost its candidacy", err)
		return CandidateNotFindErr
	}
	candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
	if nil != err {
		return err
	}
	ownerTicketIds, err := t.GetOwnerNormalTicketIds(stateDB, owner)
	if nil != err {
		return err
	}
	ticketId, err := generateTicketId()
	if err != nil {
		return err
	}
	ticket := &types.Ticket{
		TicketId:		ticketId,
		Owner:			owner,
		Deposit:		deposit,
		CandidateId:	nodeId,
		BlockNumber:	blockNumber,
	}
	candidateTicketIds = append(candidateTicketIds, ticketId)
	ownerTicketIds = append(ownerTicketIds, ticketId)
	candidate.TCount++
	if err := t.setTicket(stateDB, ticketId, ticket); err != nil {
		return err
	}
	// 设置购票人的购票关联信息
	if err := t.setOwnerNormalTicketIds(stateDB, owner, ownerTicketIds); err != nil {
		return err
	}
	if err := t.recordExpireTicket(stateDB, blockNumber, ticketId); err != nil {
		return err
	}
	if err := t.setCandidateTicketIds(stateDB, nodeId, candidateTicketIds); err != nil {
		return err
	}
	if err := t.subPoolNumber(stateDB); err != nil {
		return err
	}
	if err := t.candidatePool.UpdateCandidateTicket(stateDB, ticket.CandidateId, candidate); err != nil {
		return err
	}
	return nil
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

func (t *TicketPool) calcExpireBlockNumber(stateDB vm.StateDB, blockNumber *big.Int) *big.Int {
	num := new(big.Int).SetUint64(-1)
	if blockNumber.Cmp(new(big.Int).SetUint64(t.ExpireBlockNumber)) >= 0 {
		num.Sub(blockNumber, new(big.Int).SetUint64(t.ExpireBlockNumber))
	}
	return num
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

func (t *TicketPool) handleExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, candidateList []*types.Candidate) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ticketIdList, err := t.GetExpireTicketIds(stateDB, blockNumber)
	if err != nil {
		return err
	}
	candidateMap := make(map[discover.NodeID]*types.Candidate)
	for _, c := range candidateList {
		candidateMap[c.CandidateId] = c
	}
	for _, ticketId := range ticketIdList {
		ticket, err := t.GetTicket(stateDB, ticketId)
		if err != nil {
			return err
		}
		candidate := candidateMap[ticket.CandidateId]
		if candidate != nil {
			if _, err := t.releaseTicket(stateDB, candidate, ticketId, blockNumber); nil != err {
				continue
			}
			// 把过期票放入，购买者的过期票列表中
			ticketIds, err := t.GetOwnerNormalTicketIds(stateDB, ticket.Owner)
			if err != nil {
				return err
			}
			ticketIds, success := removeTicketId(ticketId, ticketIds)
			if success {
				if err := t.setOwnerNormalTicketIds(stateDB, ticket.Owner, ticketIds); nil != err {
					return err
				}
				ticketIds, err := t.GetOwnerExpireTicketIds(stateDB, ticket.Owner)
				if err != nil {
					return err
				}
				ticketIds = append(ticketIds, ticketId)
				if err := t.setOwnerExpireTicketIds(stateDB, ticket.Owner, ticketIds); nil != err {
					return  err
				}
			}
		}
	}
	return nil
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

func (t *TicketPool) ReturnTicket(stateDB vm.StateDB, candidate *types.Candidate, ticketId common.Hash, blockNumber *big.Int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	ticket, err := t.releaseTicket(stateDB, candidate, ticketId, blockNumber)
	t.candidatePool.UpdateCandidateTicket(stateDB, candidate.CandidateId, candidate)
	if nil != err {
		return err
	}
	// 从待过期票中删除
 	return t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
}

func (t *TicketPool) releaseTicket(stateDB vm.StateDB, candidate *types.Candidate, ticketId common.Hash, blockNumber *big.Int) (*types.Ticket, error) {
	ticket, err := t.GetTicket(stateDB, ticketId)
	if nil != err {
		return ticket, err
	}
	candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, candidate.CandidateId)
	if nil != err {
		return ticket, err
	}
	candidateTicketIds, success := removeTicketId(ticketId, candidateTicketIds)
	if !success {
		log.Error("The candidate no longer has this ticket", "err", err)
		return ticket, CandidateNotFindTicketErr
	}
	if err := t.setCandidateTicketIds(stateDB, candidate.CandidateId, candidateTicketIds); err != nil {
		return ticket, err
	}
	if err := t.addPoolNumber(stateDB); err != nil {
		return ticket, err
	}
	candidate.Epoch = candidate.Epoch.Sub(candidate.Epoch, ticket.CalcEpoch(blockNumber))
	candidate.TCount--
	return ticket, transfer(stateDB, common.TicketPoolAddr, ticket.Owner, ticket.Deposit)
}

// 1.给幸运票发放奖励
// 2.检查过期票
// 3.增加总票龄
func (t *TicketPool) Notify(stateDB vm.StateDB, blockNumber *big.Int, nodeId discover.NodeID) error {
	// 发放奖励

	candidateList := t.candidatePool.GetChosens(stateDB, 0)
	expireBlockNumber := t.calcExpireBlockNumber(stateDB, blockNumber)
	if expireBlockNumber.Uint64() > -1 {
		if err := t.handleExpireTicket(stateDB, expireBlockNumber, candidateList); nil != err {
			log.Error("OutBlockNotice method handleExpireTicket error", "blockNumber", *blockNumber, "err", err)
			return HandleExpireTicketErr
		}
	}
	// 每个候选人增加总票龄
	for _, candidate := range candidateList {
		candidate.Epoch = candidate.Epoch.Add(candidate.Epoch, new(big.Int).SetUint64(candidate.TCount))
		t.candidatePool.UpdateCandidateTicket(stateDB, candidate.CandidateId, candidate)
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
	if err := getTicketPoolState(stateDB, SurplusQuantityKey, &surplusQuantity); nil != err {
		return 0, DecodePoolNumberErr
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
