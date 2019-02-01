package pposm

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/ticketcache"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	TicketPoolNilErr          = errors.New("Ticket Insufficient quantity")
	TicketPoolOverflowErr     = errors.New("Number of ticket pool overflow")
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
	// Ticket price
	TicketPrice *big.Int
	// Maximum number of ticket pool
	MaxCount uint64
	// Reach expired quantity
	ExpireBlockNumber uint64
	lock              *sync.Mutex
}

var ticketPool *TicketPool

// initialize the global ticket pool object
func NewTicketPool(configs *params.PposConfig) *TicketPool {
	if nil != ticketPool {
		return ticketPool
	}

	if "" == strings.TrimSpace(configs.TicketConfig.TicketPrice) {
		configs.TicketConfig.TicketPrice = "1000000000000000000"
	}
	var ticketPrice *big.Int
	if price, ok := new(big.Int).SetString(configs.TicketConfig.TicketPrice, 10); !ok {
		ticketPrice, _ = new(big.Int).SetString("1000000000000000000", 10)
	} else {
		ticketPrice = price
	}

	ticketPool = &TicketPool{
		TicketPrice:       ticketPrice,
		MaxCount:          configs.TicketConfig.MaxCount,
		ExpireBlockNumber: configs.TicketConfig.ExpireBlockNumber,
		lock:              &sync.Mutex{},
	}
	return ticketPool
}

func (t *TicketPool) VoteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error) {
	log.Debug("Start voting", "owner", owner.Hex(), "voteNumber", voteNumber, "price", deposit.Uint64(), "nodeId", nodeId.String(), "blockNumber", blockNumber.Uint64())
	voteTicketIdList, err := t.voteTicket(stateDB, owner, voteNumber, deposit, nodeId, blockNumber)
	if nil != err {
		log.Error("Voting failed", "nodeId", nodeId.String(), "voteNumber", voteNumber, "successNum", len(voteTicketIdList), "err", err)
		return voteTicketIdList, err
	}
	// Voting completed, candidates reordered
	log.Debug("Successfully voted to start updating the list of candidates", "successNum", len(voteTicketIdList))
	if err := cContext.UpdateElectedQueue(stateDB, blockNumber, nodeId); nil != err {
		log.Error("Failed to Update candidate when voteTicket success", "err", err)
	}
	log.Debug("Successful vote, candidate list updated successfully", "successNum", len(voteTicketIdList))
	return voteTicketIdList, nil
}

func (t *TicketPool) voteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	voteTicketIdList := make([]common.Hash, 0)
	// check ticket pool count
	surplusQuantity, err := t.GetPoolNumber(stateDB)
	if nil != err {
		return voteTicketIdList, err
	}
	log.Debug("Ticket pool", "surplusQuantity", surplusQuantity, "voteNumber", voteNumber, "blockNumber", blockNumber.Uint64())
	if surplusQuantity == 0 {
		log.Error("Ticket Insufficient quantity")
		return voteTicketIdList, TicketPoolNilErr
	}
	if surplusQuantity < voteNumber {
		voteNumber = surplusQuantity
	}
	log.Debug("Start circular voting", "nodeId", nodeId.String())
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
		log.Debug("setTicket succeeds, start recording tickets to expire", "blockNumber", blockNumber.Uint64(), "ticketId", ticketId.String())
		if err := t.recordExpireTicket(stateDB, blockNumber, ticketId); err != nil {
			return voteTicketIdList, err
		}
		surplusQuantity, err := t.GetPoolNumber(stateDB)
		if nil != err {
			return voteTicketIdList, err
		}
		log.Debug("Record the success of the ticket to expire, and start reducing the number of tickets", "surplusQuantity", surplusQuantity)
		if err := t.subPoolNumber(stateDB); err != nil {
			return voteTicketIdList, err
		}
		surplusQuantity, err = t.GetPoolNumber(stateDB)
		if nil != err {
			return voteTicketIdList, err
		}
		log.Debug("Reduce the remaining amount of the ticket pool successfully", "surplusQuantity", surplusQuantity)
	}
	log.Debug("End loop voting", "nodeId", nodeId.String())
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

// In the current block,
// the ticket id is placed in the value slice with the block height as the key to find the expired ticket.
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
	log.Debug("Remove from pending expired tickets", "blockNumber", blockNumber.Uint64(), "ticketId", ticketId.Hex())
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
	log.Info("Pending ticket to be processed", "amount", len(ticketIdList), "blockNumber", expireBlockNumber.Uint64())
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
	log.Debug("After processing the expired ticket, update the candidate's total ticket age", "candidateNum", len(changeNodeIdList), "currentBlockNumber", currentBlockNumber.Uint64())
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
			log.Error("Did not find this ticket", "ticketId", ticketId)
			continue
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
	log.Info("Start processing tickets for the drop list", "candidateNum", len(nodeIds))
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
		log.Debug("Update candidate information", "nodeId", nodeId.String(), "epoch", candidateAttach.Epoch)
		if err := t.setCandidateAttach(stateDB, nodeId, candidateAttach); nil != err {
			return err
		}
		log.Debug("Delete candidate ticket collection", "nodeId", nodeId.String(), "ticketSize", len(candidateTicketIds))
		if err := stateDB.RemoveTicketCache(nodeId, candidateTicketIds[:]); nil != err {
			return err
		}
		log.Debug("Start processing each invalid ticket", "nodeId", nodeId.String(), "ticketSize", len(candidateTicketIds))
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
	log.Debug("End processing the list")
	return nil
}

func (t *TicketPool) ReturnTicket(stateDB vm.StateDB, nodeId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) error {
	log.Info("Release the selected ticket", "nodeId", nodeId.String(), "ticketId", ticketId.Hex(), "blockNumber", blockNumber.Uint64())
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
	log.Debug("Update ticket", "state", ticket.State, "blockNumber", blockNumber.Uint64())
	if err := t.setTicket(stateDB, ticketId, ticket); nil != err {
		return err
	}
	log.Debug("Update candidate total epoch", "nodeId", nodeId.String(), "epoch", candidateAttach.Epoch)
	if err := t.setCandidateAttach(stateDB, nodeId, candidateAttach); nil != err {
		return err
	}
	// Remove from pending expire tickets
	return t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
}

func (t *TicketPool) releaseTicket(stateDB vm.StateDB, candidateId discover.NodeID, candidateAttach *types.CandidateAttach, ticketId common.Hash, blockNumber *big.Int) (*types.Ticket, error) {
	log.Debug("Start executing releaseTicket", "nodeId", candidateId.String(), "ticketId", ticketId.Hex(), "epoch", candidateAttach.Epoch, "blockNumber", blockNumber.Uint64())
	ticket, err := t.GetTicket(stateDB, ticketId)
	if nil != err {
		return ticket, err
	}

	log.Debug("releaseTicket,Start Update", "nodeId", candidateId.String())
	candidateTicketIds := make([]common.Hash, 0)
	candidateTicketIds = append(candidateTicketIds, ticketId)
	if err := stateDB.RemoveTicketCache(candidateId, candidateTicketIds); err != nil {
		return ticket, err
	}
	surplusQuantity, err := t.GetPoolNumber(stateDB)
	if nil != err {
		return ticket, err
	}
	log.Debug("releaseTicket, end update", "nodeId", candidateId.String())
	log.Debug("releaseTicket, start to update the ticket pool", "surplusQuantity", surplusQuantity)
	if err := t.addPoolNumber(stateDB); err != nil {
		return ticket, err
	}
	surplusQuantity, err = t.GetPoolNumber(stateDB)
	if nil != err {
		return ticket, err
	}
	log.Debug("releaseTicket, end the update ticket pool", "surplusQuantity", surplusQuantity)
	log.Debug("releaseTicket, start updating the total epoch of candidates", "nodeId", candidateId.String(), "totalEpoch", candidateAttach.Epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	candidateAttach.SubEpoch(ticket.CalcEpoch(blockNumber))
	log.Debug("releaseTicket, the end of the update candidate total epoch", "nodeId", candidateId.String(), "totalEpoch", candidateAttach.Epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	return ticket, nil
}

func (t *TicketPool) Notify(stateDB vm.StateDB, blockNumber *big.Int) error {
	// Check expired tickets
	expireBlockNumber, ok := t.calcExpireBlockNumber(stateDB, blockNumber)
	log.Debug("Check expired tickets", "isOk", ok, "expireBlockNumber", expireBlockNumber.Uint64())
	if ok {
		if nodeIdList, err := t.handleExpireTicket(stateDB, expireBlockNumber, blockNumber); nil != err {
			log.Error("OutBlockNotice method handleExpireTicket error", "blockNumber", *blockNumber, "err", err)
			return HandleExpireTicketErr
		} else {
			// Notify the candidate to update the list information after processing the expired ticket
			log.Debug("After processing the expired ticket, start updating the candidate list", "blockNumber", blockNumber.Uint64(), "nodeIdList", len(nodeIdList))
			if err := cContext.UpdateElectedQueue(stateDB, blockNumber, nodeIdList...); nil != err {
				log.Error("Failed to Update candidate when handleExpireTicket success on Notify", "err", err)
			}
		}
	}
	// Increase the total number of epoch for each candidate
	log.Debug("Increase the total number of epoch for each candidate", "blockNumber", blockNumber.Uint64())
	if err := t.calcCandidateEpoch(stateDB, blockNumber); nil != err {
		return err
	}
	return nil
}

func (t *TicketPool) calcCandidateEpoch(stateDB vm.StateDB, blockNumber *big.Int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	candidateList := cContext.GetChosens(stateDB, 0)
	for _, candidate := range candidateList {
		candidateAttach, err := t.GetCandidateAttach(stateDB, candidate.CandidateId)
		if nil != err {
			return err
		}
		// Get the total number of votes, increase the total epoch
		ticketCount := stateDB.TCount(candidate.CandidateId)
		log.Debug("increase the total epoch", "candidateId", candidate.CandidateId.String(), "ticketCount", ticketCount)
		if ticketCount > 0 {
			candidateAttach.AddEpoch(new(big.Int).SetUint64(ticketCount))
			if err := t.setCandidateAttach(stateDB, candidate.CandidateId, candidateAttach); nil != err {
				return err
			}
		}
	}
	return nil
}

// Simple version of the lucky ticket algorithm
// According to the previous block Hash,
// find the first ticket Id which is larger than the Hash. If not found, the last ticket Id is taken.
func (t *TicketPool) SelectionLuckyTicket(stateDB vm.StateDB, nodeId discover.NodeID, blockHash common.Hash) (common.Hash, error) {
	candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
	log.Debug("Start picking lucky tickets", "nodeId", nodeId.String(), "blockHash", blockHash.Hex(), "candidateTicketIds", len(candidateTicketIds))
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
	log.Debug("Pick out a lucky ticket", "index", index)
	luckyTicketId = decMap[decList[index]]
	log.Debug("End the selection of lucky tickets", "nodeId", nodeId.String(), "blockHash", blockHash.Hex(), "luckyTicketId", luckyTicketId.Hex(), "candidateTicketIds", len(candidateTicketIds))
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
	surplusQuantity, err := t.GetPoolNumber(stateDB)
	if nil != err {
		return err
	}
	if surplusQuantity == t.MaxCount {
		return TicketPoolOverflowErr
	}
	surplusQuantity++
	return t.setPoolNumber(stateDB, surplusQuantity)
}

func (t *TicketPool) subPoolNumber(stateDB vm.StateDB) error {
	surplusQuantity, err := t.GetPoolNumber(stateDB)
	if nil != err {
		return err
	}
	if surplusQuantity == 0 {
		return TicketPoolNilErr
	}
	surplusQuantity--
	return t.setPoolNumber(stateDB, surplusQuantity)
}

func (t *TicketPool) setPoolNumber(stateDB vm.StateDB, surplusQuantity uint64) error {
	if value, err := rlp.EncodeToBytes(surplusQuantity); nil != err {
		log.Error("Failed to encode surplusQuantity object on setPoolNumber", "key", GetSurplusQuantityKey(), "err", err)
		return EncodePoolNumberErr
	} else {
		setTicketPoolState(stateDB, GetSurplusQuantityKey(), value)
	}
	return nil
}

func (t *TicketPool) GetPoolNumber(stateDB vm.StateDB) (uint64, error) {
	var surplusQuantity uint64
	if val := stateDB.GetState(common.TicketPoolAddr, GetSurplusQuantityKey()); len(val) > 0 {
		if err := rlp.DecodeBytes(val, &surplusQuantity); nil != err {
			log.Error("Decode PoolNumber error", "key", string(GetSurplusQuantityKey()), "err", err)
			return surplusQuantity, DecodePoolNumberErr
		}
	} else {
		// Default initialization values
		surplusQuantity = t.MaxCount
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

func (t *TicketPool) GetCandidatesTicketIds(stateDB vm.StateDB, nodeIds []discover.NodeID) (map[discover.NodeID][]common.Hash, error) {
	result := make(map[discover.NodeID][]common.Hash)
	if nodeIds != nil {
		for _, nodeId := range nodeIds {
			ticketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
			if nil != err {
				continue
			}
			result[nodeId] = ticketIds
		}
	}
	return result, nil
}

func (t *TicketPool) GetCandidatesTicketCount(stateDB vm.StateDB, nodeIds []discover.NodeID) (map[discover.NodeID]int, error) {
	result := make(map[discover.NodeID]int)
	if nil != nodeIds {
		for _, nodeId := range nodeIds {
			candidateTicketCount, err := stateDB.GetTicketCache(nodeId)
			if nil != err {
				continue
			}
			result[nodeId] = len(candidateTicketCount)
		}
	}
	return result, nil
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
	return t.TicketPrice, nil
}

// Save the hash value of the current state of the ticket pool
func (t *TicketPool) CommitHash(stateDB vm.StateDB) error {
	hash, err := ticketcache.Hash(stateDB.TicketCaceheSnapshot())
	if nil != err {
		return err
	}
	setTicketPoolState(stateDB, addCommonPrefix(TicketPoolHashKey), hash.Bytes())
	return nil
}

func GetTicketPtr() *TicketPool {
	return ticketPool
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
	return addCommonPrefix(append(ExpireTicketPrefix, key...))
}

func CandidateAttachKey(key []byte) []byte {
	return addCommonPrefix(append(CandidateAttachPrefix, key...))
}

func GetSurplusQuantityKey() []byte {
	return addCommonPrefix(SurplusQuantityKey)
}

func addCommonPrefix(key []byte) []byte {
	return append(common.TicketPoolAddr.Bytes(), key...)
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
	// If no match is found, the last subscript is returned by default.
	if left >= len(list) {
		return len(list) - 1
	}
	return left
}
