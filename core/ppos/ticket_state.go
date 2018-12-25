package pposm

import (
	"Platon-go/common"
	"Platon-go/common/hexutil"
	"Platon-go/core/ticketcache"
	"Platon-go/core/types"
	"Platon-go/core/vm"
	"Platon-go/crypto/sha3"
	"Platon-go/log"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"Platon-go/rlp"
	"errors"
	"math/big"
	"sort"
	"strconv"
	"sync"
)

var (
	TicketPoolNilErr          = errors.New("Ticket Insufficient quantity")
	EncodeTicketErr           = errors.New("Encode Ticket error")
	EncodePoolNumberErr       = errors.New("Encode SurplusQuantity error")
	DecodeTicketErr           = errors.New("Decode Ticket error")
	DecodePoolNumberErr       = errors.New("Decode SurplusQuantity error")
	RecordExpireTicketErr     = errors.New("Record Expire Ticket error")
	CandidateNotFindTicketErr = errors.New("The candidate no longer has this ticket")
	CandidateNilTicketErr     = errors.New("This candidate has no ticket")
	GetCandidateTicketIdErr   = errors.New("Get Candidate TicketIds error")
	SetCandidateTicketIdErr   = errors.New("Update Candidate TicketIds error")
	TicketPoolBalanceErr      = errors.New("TicketPool not sufficient funds")
	GetOwnerTicketIdsErr      = errors.New("Get Owner TicketIds error")
	SetOwnerTicketIdsErr      = errors.New("Update Owner TicketIds error")
	TicketIdNotFindErr        = errors.New("TicketId not find")
	HandleExpireTicketErr     = errors.New("Failure to deal with expired tickets")
	GetCandidateAttachErr     = errors.New("Get CandidateAttach error")
	SetCandidateAttachErr     = errors.New("Update CandidateAttach error")
)

type TicketPool struct {
	// Maximum number of ticket pool
	MaxCount uint64
	// Remaining number of ticket pool
	SurplusQuantity uint64
	// Reach expired quantity
	ExpireBlockNumber uint64
	lock              *sync.RWMutex
}

var ticketPool *TicketPool

// initialize the global ticket pool object
func NewTicketPool(configs *params.PposConfig) *TicketPool {
	if nil != ticketPool {
		return ticketPool
	}
	ticketPool = &TicketPool{
		MaxCount:          configs.TicketConfig.MaxCount,
		SurplusQuantity:   configs.TicketConfig.MaxCount,
		ExpireBlockNumber: configs.TicketConfig.ExpireBlockNumber,
		lock:              &sync.RWMutex{},
	}
	return ticketPool
}

func (t *TicketPool) VoteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error) {
	log.Info("开始投票", "购票人：", owner.Hex(), "购票数量：", voteNumber, "购票单价：", deposit.Uint64(), "所投节点：", nodeId.String(), "块高：", blockNumber.Uint64())
	voteTicketIdList, err := t.voteTicket(stateDB, owner, voteNumber, deposit, nodeId, blockNumber)
	if nil != err {
		log.Error("投票失败", "所投节点：", nodeId.String(), "购票数量：", voteNumber, "成功数量：", len(voteTicketIdList), "err", err)
		return voteTicketIdList, err
	}
	// 调用候选人重新排序接口
	log.Info("投票成功，开始更新候选人榜单", "成功票数", len(voteTicketIdList))
	candidatePool.UpdateElectedQueue(stateDB, blockNumber, nodeId)
	log.Info("投票成功，候选人榜单更新成功", "成功票数", len(voteTicketIdList))
	return voteTicketIdList, nil
}

func (t *TicketPool) voteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	voteTicketIdList := make([]common.Hash, 0)
	// check ticket pool count
	t.GetPoolNumber(stateDB)
	log.Info("票池", "剩余数量：", t.SurplusQuantity, "购票数量：", voteNumber, "块高：", blockNumber.Uint64())
	if t.SurplusQuantity == 0 {
		log.Error("Ticket Insufficient quantity")
		return voteTicketIdList, TicketPoolNilErr
	}
	if t.SurplusQuantity < voteNumber {
		voteNumber -= t.SurplusQuantity
	}
	log.Info("开始循环投票", "候选人：", nodeId.String())
	var i uint64 = 0
	for ; i < voteNumber; i++ {
		ticketId, err := generateTicketId(stateDB.TxHash(), i)
		if err != nil {
			return voteTicketIdList, err
		}
		ticket := &types.Ticket{
			TicketId:    ticketId,
			Owner:       owner,
			Deposit:     deposit,
			CandidateId: nodeId,
			BlockNumber: blockNumber,
		}
		ticket.SetNormal()
		voteTicketIdList = append(voteTicketIdList, ticketId)
		if err := t.setTicket(stateDB, ticketId, ticket); err != nil {
			return voteTicketIdList, err
		}
		log.Info("setTicket成功，开始记录待过期票", "块高：", blockNumber.Uint64(), "票Id: ", ticketId.String())
		if err := t.recordExpireTicket(stateDB, blockNumber, ticketId); err != nil {
			return voteTicketIdList, err
		}
		log.Info("记录待过期票成功，开始减少票池数量", "剩余量：", t.SurplusQuantity)
		if err := t.subPoolNumber(stateDB); err != nil {
			return voteTicketIdList, err
		}
		log.Info("减少票池剩余量成功", "剩余量：", t.SurplusQuantity)
	}
	log.Info("结束循环投票", "候选人：", nodeId.String())
	stateDB.AppendTicketCache(nodeId, voteTicketIdList)
	return voteTicketIdList, nil
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
		log.Error("Failed to encode ticketId object on setExpireTicket", "key", *blockNumber, "err", err)
		return EncodeTicketErr
	} else {
		setTicketPoolState(stateDB, ExpireTicketKey((*blockNumber).Bytes()), value)
	}
	return nil
}

func (t *TicketPool) removeExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, ticketId common.Hash) error {
	log.Info("从待过期票记录中删除", "块高：", blockNumber.Uint64(), "票Id：", ticketId.Hex())
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

func (t *TicketPool) handleExpireTicket(stateDB vm.StateDB, expireBlockNumber *big.Int, currentBlockNumber *big.Int) ([]discover.NodeID, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	ticketIdList, err := t.GetExpireTicketIds(stateDB, expireBlockNumber)
	if err != nil {
		return nil, err
	}
	log.Info("待处理的过期票", "数量：", len(ticketIdList), "块高：", expireBlockNumber.Uint64())
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
				return changeNodeIdList, err
			}
			candidateAttachMap[ticket.CandidateId] = candidateAttach
			changeNodeIdList = append(changeNodeIdList, ticket.CandidateId)
		}
		if _, err := t.releaseTicket(stateDB, ticket.CandidateId, candidateAttach, ticketId, currentBlockNumber); nil != err {
			return changeNodeIdList, err
		}
		// Set ticket state to expired
		ticket.SetExpired(currentBlockNumber)
		if err := t.setTicket(stateDB, ticketId, ticket); nil != err {
			return changeNodeIdList, err
		}
	}
	log.Info("处理完过期票，更新候选人总票龄", "候选人数量：", len(changeNodeIdList), "当前块高：", currentBlockNumber.Uint64())
	// Update CandidateAttach
	for nodeId, ca := range candidateAttachMap {
		t.setCandidateAttach(stateDB, nodeId, ca)
	}
	return changeNodeIdList, nil
}

// Get ticket list
func (t *TicketPool) GetTicketList(stateDB vm.StateDB, ticketIds []common.Hash) ([]*types.Ticket, error) {
	var tickets []*types.Ticket
	for _, ticketId := range ticketIds {
		ticket, err := t.GetTicket(stateDB, ticketId)
		if nil != err || ticket.TicketId == (common.Hash{}) {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

// Get ticket details based on TicketId
func (t *TicketPool) GetTicket(stateDB vm.StateDB, ticketId common.Hash) (*types.Ticket, error) {
	var ticket = new(types.Ticket)
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
	t.lock.Lock()
	defer t.lock.Unlock()
	log.Info("开始处理掉榜票", "掉榜候选人数量：", len(nodeIds))
	for _, nodeId := range nodeIds {
		candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
		if nil != err {
			return err
		}
		candidateAttach, err := t.GetCandidateAttach(stateDB, nodeId)
		if nil != err {
			return err
		}
		candidateAttach.Epoch = new(big.Int).SetUint64(0)
		log.Info("更新掉榜信息", "候选人：", nodeId.String(), "票龄：", candidateAttach.Epoch)
		if err := t.setCandidateAttach(stateDB, nodeId, candidateAttach); nil != err {
			return err
		}
		log.Info("删除掉榜信息", "候选人：", nodeId.String(), "所得票：", len(candidateTicketIds))
		if err := stateDB.RemoveTicketCache(nodeId, candidateTicketIds); nil != err {
			return err
		}
		log.Info("开始处理掉榜的票", "候选人：", nodeId.String(), "总票数：", len(candidateTicketIds))
		for _, ticketId := range candidateTicketIds {
			ticket, err := t.GetTicket(stateDB, ticketId)
			if nil != err {
				return err
			}
			ticket.SetInvalid(blockNumber)
			if err := t.setTicket(stateDB, ticketId, ticket); nil != err {
				return err
			}
			t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
		}
	}
	log.Info("结束处理掉榜票")
	return nil
}

func (t *TicketPool) ReturnTicket(stateDB vm.StateDB, nodeId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) error {
	log.Info("释放选中票", "候选人：", nodeId.String(), "票Id：", ticketId.Hex(), "块高：", blockNumber.Uint64())
	t.lock.Lock()
	defer t.lock.Unlock()
	if ticketId == (common.Hash{}) {
		return TicketIdNotFindErr
	}
	candidateAttach, err := t.GetCandidateAttach(stateDB, nodeId)
	if nil != err {
		return err
	}
	ticket, err := t.releaseTicket(stateDB, nodeId, candidateAttach, ticketId, blockNumber)
	ticket.SetSelected(blockNumber)
	log.Info("更新票", "状态为：", ticket.State, "释放块高", blockNumber.Uint64())
	if err := t.setTicket(stateDB, ticketId, ticket); nil != err {
		return err
	}
	log.Info("更新候选人总票龄", "候选人：", nodeId.String(), "票龄：", candidateAttach.Epoch)
	if err := t.setCandidateAttach(stateDB, nodeId, candidateAttach); nil != err {
		return err
	}
	// 从待过期票中删除
	return t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
}

func (t *TicketPool) releaseTicket(stateDB vm.StateDB, candidateId discover.NodeID, candidateAttach *types.CandidateAttach, ticketId common.Hash, blockNumber *big.Int) (*types.Ticket, error) {
	log.Info("开始执行releaseTicket", "候选人：", candidateId.String(), "Epoch：", candidateAttach.Epoch, "块高：", blockNumber.Uint64())
	ticket, err := t.GetTicket(stateDB, ticketId)
	if nil != err {
		return ticket, err
	}
	log.Info("releaseTicket,开始更新", "候选人：", candidateId.String())
	candidateTicketIds := make([]common.Hash, 0)
	candidateTicketIds = append(candidateTicketIds, ticketId)
	if err := stateDB.RemoveTicketCache(candidateId, candidateTicketIds); err != nil {
		return ticket, err
	}
	log.Info("releaseTicket,结束更新", "候选人：", candidateId.String())
	log.Info("releaseTicket,开始更新票池", "剩余量：", t.SurplusQuantity)
	if err := t.addPoolNumber(stateDB); err != nil {
		return ticket, err
	}
	log.Info("releaseTicket,结束更新票池", "剩余量：", t.SurplusQuantity)
	log.Info("releaseTicket,开始更新候选人总票龄", "候选人：", candidateId.String(), "总票龄：", candidateAttach.Epoch, "当前块高：", blockNumber.Uint64(), "票块高：", ticket.BlockNumber.Uint64())
	candidateAttach.SubEpoch(ticket.CalcEpoch(blockNumber))
	log.Info("releaseTicket,结束更新候选人总票龄", "候选人：", candidateId.String(), "总票龄：", candidateAttach.Epoch, "当前块高：", blockNumber.Uint64(), "票块高：", ticket.BlockNumber.Uint64())
	return ticket, nil
}

// 1.给幸运票发放奖励
// 2.检查过期票
// 3.增加总票龄
func (t *TicketPool) Notify(stateDB vm.StateDB, blockNumber *big.Int) error {
	// 发放奖励

	// 检查过期票
	expireBlockNumber, ok := t.calcExpireBlockNumber(stateDB, blockNumber)
	log.Info("检查过期票", "是否需要处理：", ok, "过期票所在块高：", expireBlockNumber.Uint64())
	if ok {
		if nodeIdList, err := t.handleExpireTicket(stateDB, expireBlockNumber, blockNumber); nil != err {
			log.Error("OutBlockNotice method handleExpireTicket error", "blockNumber", *blockNumber, "err", err)
			return HandleExpireTicketErr
		} else {
			// 处理完过期票之后，通知候选人更新榜单信息
			log.Info("处理完过期票，开始更新候选人榜单", "块高：", blockNumber.Uint64(), "变动候选人数量：", len(nodeIdList))
			candidatePool.UpdateElectedQueue(stateDB, blockNumber, nodeIdList...)
		}
	}
	// 每个候选人增加总票龄
	log.Info("开始为所有候选人增加总票龄", "块高：", blockNumber.Uint64())
	if err := t.calcCandidateEpoch(stateDB, blockNumber); nil != err {
		return err
	}
	return nil
}

func (t *TicketPool) calcCandidateEpoch(stateDB vm.StateDB, blockNumber *big.Int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	candidateList := candidatePool.GetChosens(stateDB, 0)
	for _, candidate := range candidateList {
		candidateAttach, err := t.GetCandidateAttach(stateDB, candidate.CandidateId)
		if nil != err {
			return err
		}
		// 获取总票数，增加总票龄
		ticketCount := stateDB.TCount(candidate.CandidateId)
		if ticketCount > 0 {
			candidateAttach.AddEpoch(new(big.Int).SetUint64(ticketCount))
			if err := t.setCandidateAttach(stateDB, candidate.CandidateId, candidateAttach); nil != err {
				return err
			}
		}
	}
	return nil
}

// 简版幸运票算法 --> 根据上一个区块Hash找到第一个比该Hash大的票Id，找不到则取最后一个票Id
func (t *TicketPool) SelectionLuckyTicket(stateDB vm.StateDB, nodeId discover.NodeID, blockHash common.Hash) (common.Hash, error) {
	candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
	log.Info("开始选取幸运票", "候选人", nodeId.String(), "区块Hash", blockHash.Hex(), "候选人票数", len(candidateTicketIds))
	luckyTicketId := common.Hash{}
	if nil != err {
		return luckyTicketId, err
	}
	if len(candidateTicketIds) == 0 {
		return luckyTicketId, CandidateNilTicketErr
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
	log.Info("选出幸运票", "下标", index)
	luckyTicketId = decMap[decList[index]]
	log.Info("结束选取幸运票", "候选人", nodeId.String(), "区块Hash", blockHash.Hex(), "幸运票", luckyTicketId.Hex(), "候选人票数", len(candidateTicketIds))
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
	candidateTicketIds, err := stateDB.GetTicketCache(nodeId)
	if nil != err {
		return nil, err
	}
	return candidateTicketIds, nil
}

func (t *TicketPool) GetCandidateAttach(stateDB vm.StateDB, nodeId discover.NodeID) (*types.CandidateAttach, error) {
	candidateAttach := new(types.CandidateAttach)
	candidateAttach.Epoch = new(big.Int)
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

func (t *TicketPool) GetCandidateEpoch(stateDB vm.StateDB, nodeId discover.NodeID) (uint64, error) {
	candidateAttach, err := t.GetCandidateAttach(stateDB, nodeId)
	if nil != err {
		return 0, err
	}
	return candidateAttach.Epoch.Uint64(), nil
}

func (t *TicketPool) GetTicketPrice(stateDB vm.StateDB) (*big.Int, error) {
	return new(big.Int).SetUint64(1), nil
}

// Save the hash value of the current state of the ticket pool
func (t *TicketPool) CommitHash(stateDB vm.StateDB) error {
	hash, err := ticketcache.GetTicketidsCachePtr().Hash(stateDB.TicketCaceheSnapshot())
	if nil != err {
		return err
	}
	setTicketPoolState(stateDB, TicketPoolHashKey, hash.Bytes())
	return nil
}

func GetTicketPtr() *TicketPool {
	return ticketPool
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

func generateTicketId(txHash common.Hash, index uint64) (common.Hash, error) {
	// generate ticket id
	value := append(txHash.Bytes(), []byte(strconv.Itoa(int(index)))...)
	ticketId := sha3.Sum256(value[:])
	return ticketId, nil
}

func ExpireTicketKey(key []byte) []byte {
	return append(ExpireTicketPrefix, key...)
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
	if left >= len(list) {
		return len(list) - 1
	}
	return left
}
