package pposm

import (
	"errors"
	"fmt"
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
	"runtime"
	"sort"
	"strconv"
	"strings"
	"sync"
)

var (
	TicketPoolNilErr      = errors.New("Ticket Insufficient quantity")
	TicketPoolOverflowErr = errors.New("Number of ticket pool overflow")
	EncodeTicketErr       = errors.New("Encode Ticket error")
	EncodePoolNumberErr   = errors.New("Encode SurplusQuantity error")
	DecodeTicketErr       = errors.New("Decode Ticket error")
	DecodePoolNumberErr   = errors.New("Decode SurplusQuantity error")
	RecordExpireTicketErr = errors.New("Record Expire Ticket error")
	CandidateNotFindErr   = errors.New("The Candidate not find")
	CandidateNilTicketErr = errors.New("This candidate has no ticket")
	TicketPoolBalanceErr  = errors.New("TicketPool not sufficient funds")
	TicketIdNotFindErr    = errors.New("TicketId not find")
	HandleExpireTicketErr = errors.New("Failure to deal with expired tickets")
	GetCandidateAttachErr = errors.New("Get CandidateAttach error")
	SetCandidateAttachErr = errors.New("Update CandidateAttach error")
	VoteTicketErr         = errors.New("Voting failed")
)

type TicketPool struct {
	// Ticket price
	LowestTicketPrice *big.Int
	//adjust cycle
	AdjustCycle *big.Int
	// Maximum number of ticket pool
	MaxCount uint64
	// Reach expired quantity
	ExpireBlockNumber uint64
	lock              *sync.Mutex
}

//var ticketPool *TicketPool

// initialize the global ticket pool object
func NewTicketPool(configs *params.PposConfig) *TicketPool {
	//if nil != ticketPool {
	//	return ticketPool
	//}
	log.Debug("Build a New TicketPool Info ...")
	if "" == strings.TrimSpace(configs.TicketConfig.TicketPrice) {
		configs.TicketConfig.TicketPrice = "100000000000000000000"
	}
	if "" == strings.TrimSpace(configs.TicketConfig.AdjustCycle) {
		configs.TicketConfig.AdjustCycle = "10000"
	}
	var ticketPrice, adjustCycle *big.Int
	if price, ok := new(big.Int).SetString(configs.TicketConfig.TicketPrice, 10); !ok {
		ticketPrice, _ = new(big.Int).SetString("100000000000000000000", 10)
	} else {
		ticketPrice = price
	}
	if cycle, ok := new(big.Int).SetString(configs.TicketConfig.AdjustCycle, 10); !ok {
		adjustCycle, _ = new(big.Int).SetString("10000", 10)
	} else {
		adjustCycle = cycle
	}

	ticketPool := &TicketPool{
		LowestTicketPrice:       ticketPrice,
		AdjustCycle:			adjustCycle,
		MaxCount:          configs.TicketConfig.MaxCount,
		ExpireBlockNumber: configs.TicketConfig.ExpireBlockNumber,
		lock:              &sync.Mutex{},
	}
	return ticketPool
}

func (t *TicketPool) VoteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) ([]common.Hash, error) {
	log.Debug("Call Voting", "statedb addr", fmt.Sprintf("%p", stateDB))
	log.Info("Start Voting,VoteTicket", "owner", owner.Hex(), "voteNumber", voteNumber, "price", deposit.Uint64(), "nodeId", nodeId.String(), "blockNumber", blockNumber.Uint64())
	voteTicketIdList, err := t.voteTicket(stateDB, owner, voteNumber, deposit, nodeId, blockNumber)
	if nil != err {
		log.Error("Voting failed", "nodeId", nodeId.String(), "voteNumber", voteNumber, "successNum", len(voteTicketIdList), "err", err)
		return voteTicketIdList, err
	}
	// Voting completed, candidates reordered
	log.Debug("Successfully voted to start updating the list of candidates,VoteTicket", "successNum", len(voteTicketIdList))
	if err := cContext.UpdateElectedQueue(stateDB, blockNumber, nodeId); nil != err {
		log.Error("Failed to Update candidate when voteTicket success", "err", err)
	}
	log.Debug("Successful vote, candidate list updated successfully,VoteTicket", "successNum", len(voteTicketIdList))
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
	log.Debug("Execute voteTicket", "surplusQuantity", surplusQuantity, "voteNumber", voteNumber, "blockNumber", blockNumber.Uint64())
	if surplusQuantity == 0 {
		log.Error("Ticket Insufficient quantity")
		return voteTicketIdList, TicketPoolNilErr
	}
	if surplusQuantity < voteNumber {
		voteNumber = surplusQuantity
	}
	log.Debug("Start circular voting", "nodeId", nodeId.String(), "voteNumber", voteNumber)

	parentRoutineID := fmt.Sprintf("%s", common.CurrentGoRoutineID())

	runtime.GOMAXPROCS(runtime.NumCPU())
	var wg sync.WaitGroup
	wg.Add(int(voteNumber))
	type genTicket struct {
		index    uint64
		ticketId common.Hash
	}
	resultCh := make(chan genTicket, voteNumber)
	var i uint64 = 0
	for ; i < voteNumber; i++ {
		go func(i uint64) {

			ticketId, _ := generateTicketId(stateDB.TxHash(), i)

			log.Debug("Call Voting parent routine "+parentRoutineID, "statedb addr", fmt.Sprintf("%p", stateDB), "ticketId", ticketId.String())

			ticket := &types.Ticket{
				TicketId:    ticketId,
				Owner:       owner,
				Deposit:     deposit,
				CandidateId: nodeId,
				BlockNumber: blockNumber,
			}
			ticket.SetNormal()
			if err := t.setTicket(stateDB, ticketId, ticket); err != nil {
				log.Error("Ticket information record failed", "ticketId", ticketId.Hex(), "err", err)
				wg.Done()
				return
			}
			genTicket := genTicket{
				index:    i,
				ticketId: ticketId,
			}
			resultCh <- genTicket
			wg.Done()
		}(i)
	}
	wg.Wait()
	close(resultCh)
	resultSize := len(resultCh)
	if resultSize == 0 {
		log.Error("Voting failed resultSize = 0", "err", err)
		return voteTicketIdList, VoteTicketErr
	}
	voteTicketIdList = make([]common.Hash, voteNumber)
	for result := range resultCh {
		voteTicketIdList[result.index] = result.ticketId
	}
	if resultSize != int(voteNumber) {
		list := make([]common.Hash, 0)
		for i := 0; i < len(voteTicketIdList); i++ {
			tid := voteTicketIdList[i]
			if tid != (common.Hash{}) {
				list = append(list, tid)
			}
		}
		voteTicketIdList = list
	}
	log.Debug("SetTicket succeeds, start recording tickets to expire", "blockNumber", blockNumber.Uint64(), "ticketAmount", len(voteTicketIdList))
	if err := t.recordExpireTicket(stateDB, blockNumber, voteTicketIdList); err != nil {
		return voteTicketIdList, err
	}
	log.Debug("Record the success of the ticket to expire, and start reducing the number of tickets", "surplusQuantity", surplusQuantity)
	if err := t.setPoolNumber(stateDB, surplusQuantity-uint64(len(voteTicketIdList))); err != nil {
		return voteTicketIdList, err
	}
	surplusQuantity, err = t.GetPoolNumber(stateDB)
	if nil != err {
		return voteTicketIdList, err
	}

	stateDB.AppendTicketCache(nodeId, voteTicketIdList[:])

	log.Debug("Voting SUCCUESS !!!!!!  Reduce the remaining amount of the ticket pool successfully", "surplusQuantity", surplusQuantity, "nodeId", nodeId.String())
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
	log.Debug("Call GetExpireTicketIds", "statedb addr", fmt.Sprintf("%p", stateDB))
	var expireTicketIds []common.Hash
	if err := getTicketPoolState(stateDB, ExpireTicketKey((*blockNumber).Bytes()), &expireTicketIds); nil != err {
		return nil, err
	}
	return expireTicketIds, nil
}

// In the current block,
// the ticket id is placed in the value slice with the block height as the key to find the expired ticket.
func (t *TicketPool) recordExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, ticketIds []common.Hash) error {
	expireTickets, err := t.GetExpireTicketIds(stateDB, blockNumber)
	if err != nil {
		log.Error("recordExpireTicket error", "key", blockNumber.Uint64(), "ticketIdSize", len(ticketIds), "err", err)
		return RecordExpireTicketErr
	}
	expireTickets = append(expireTickets, ticketIds...)
	return t.setExpireTicket(stateDB, blockNumber, expireTickets)
}

func (t *TicketPool) setExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, expireTickets []common.Hash) error {
	if value, err := rlp.EncodeToBytes(expireTickets); nil != err {
		log.Error("Failed to encode ticketId object on setExpireTicket", "key", blockNumber.Uint64(), "err", err)
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
	if len(ticketIdList) == 0 {
		return nil
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
	log.Info("Pending ticket to be processed", "amount", len(ticketIdList), "expireBlockNumber", expireBlockNumber.Uint64(), "currentBlockNumber", currentBlockNumber.Uint64())
	candidateAttachMap := make(map[discover.NodeID]*types.CandidateAttach)
	changeNodeIdList := make([]discover.NodeID, 0)
	for _, ticketId := range ticketIdList {
		ticket, err := t.GetTicket(stateDB, ticketId)
		if err != nil {
			return nil, err
		}
		if ticket.TicketId == (common.Hash{}) {
			continue
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
	log.Debug("Call GetTickList", "statedb addr", fmt.Sprintf("%p", stateDB))
	var tickets []*types.Ticket
	for _, ticketId := range ticketIds {
		ticket, err := t.GetTicket(stateDB, ticketId)
		if nil != err || ticket.TicketId == (common.Hash{}) {
			log.Error("find this ticket fail", "ticketId", ticketId.Hex())
			continue
		}
		tickets = append(tickets, ticket)
	}
	return tickets, nil
}

// Get ticket details based on TicketId
func (t *TicketPool) GetTicket(stateDB vm.StateDB, ticketId common.Hash) (*types.Ticket, error) {
	log.Debug("Call GetTicket", "statedb addr", fmt.Sprintf("%p", stateDB))
	var ticket = new(types.Ticket)
	if err := getTicketPoolState(stateDB, ticketId.Bytes(), ticket); nil != err {
		return nil, DecodeTicketErr
	}
	return ticket, nil
}

func (t *TicketPool) setTicket(stateDB vm.StateDB, ticketId common.Hash, ticket *types.Ticket) error {
	if value, err := rlp.EncodeToBytes(ticket); nil != err {
		log.Error("Failed to encode ticket object on setTicket", "key", ticketId.Hex(), "err", err)
		return EncodeTicketErr
	} else {
		setTicketPoolState(stateDB, ticketId.Bytes(), value)
	}
	return nil
}

func (t *TicketPool) DropReturnTicket(stateDB vm.StateDB, blockNumber *big.Int, nodeIds ...discover.NodeID) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	log.Debug("Call DropReturnTicket", "statedb addr", fmt.Sprintf("%p", stateDB))
	log.Info("Start processing tickets for the drop list on DropReturnTicket", "candidateNum", len(nodeIds), "blockNumber", blockNumber.Uint64())
	for _, nodeId := range nodeIds {
		if nodeId == (discover.NodeID{}) {
			continue
		}
		candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
		if nil != err {
			return err
		}
		if len(candidateTicketIds) == 0 {
			continue
		}
		candidateAttach, err := t.GetCandidateAttach(stateDB, nodeId)
		if nil != err {
			return err
		}
		candidateAttach.Epoch = new(big.Int).SetUint64(0)
		log.Debug("Update candidate information on DropReturnTicket", "nodeId", nodeId.String(), "epoch", candidateAttach.Epoch.Uint64())
		if err := t.setCandidateAttach(stateDB, nodeId, candidateAttach); nil != err {
			return err
		}
		log.Debug("Delete candidate ticket collection on DropReturnTicket", "nodeId", nodeId.String(), "ticketSize", len(candidateTicketIds))
		if err := stateDB.RemoveTicketCache(nodeId, candidateTicketIds[:]); nil != err {
			return err
		}
		surplusQuantity, err := t.GetPoolNumber(stateDB)
		if nil != err {
			return err
		}
		log.Debug("Start reducing the number of tickets on DropReturnTicket", "surplusQuantity", surplusQuantity, "candidateTicketIds", len(candidateTicketIds))
		if err := t.setPoolNumber(stateDB, surplusQuantity+uint64(len(candidateTicketIds))); err != nil {
			return err
		}
		log.Debug("Start processing each invalid ticket on DropReturnTicket", "nodeId", nodeId.String(), "ticketSize", len(candidateTicketIds))
		for _, ticketId := range candidateTicketIds {
			ticket, err := t.GetTicket(stateDB, ticketId)
			if nil != err {
				return err
			}
			if ticket.TicketId == (common.Hash{}) {
				continue
			}
			ticket.SetInvalid(blockNumber)
			if err := t.setTicket(stateDB, ticketId, ticket); nil != err {
				return err
			}
			if err := transfer(stateDB, common.TicketPoolAddr, ticket.Owner, ticket.Deposit); nil != err {
				return err
			}
			t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
		}
	}
	log.Debug("End processing the list on DropReturnTicket")
	return nil
}

func (t *TicketPool) ReturnTicket(stateDB vm.StateDB, nodeId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) error {
	log.Debug("Call ReturnTicket", "statedb addr", fmt.Sprintf("%p", stateDB))
	log.Info("Release the selected ticket on ReturnTicket", "nodeId", nodeId.String(), "ticketId", ticketId.Hex(), "blockNumber", blockNumber.Uint64())
	t.lock.Lock()
	defer t.lock.Unlock()
	if ticketId == (common.Hash{}) {
		return TicketIdNotFindErr
	}
	if nodeId == (discover.NodeID{}) {
		return CandidateNotFindErr
	}
	candidateAttach, err := t.GetCandidateAttach(stateDB, nodeId)
	if nil != err {
		return err
	}
	ticket, err := t.releaseTicket(stateDB, nodeId, candidateAttach, ticketId, blockNumber)
	if ticket == nil {
		return TicketIdNotFindErr
	}
	ticket.SetSelected(blockNumber)
	log.Debug("Update ticket on ReturnTicket", "state", ticket.State, "blockNumber", blockNumber.Uint64())
	if err := t.setTicket(stateDB, ticketId, ticket); nil != err {
		return err
	}
	log.Debug("Update candidate total epoch on ReturnTicket", "nodeId", nodeId.String(), "epoch", candidateAttach.Epoch.Uint64())
	if err := t.setCandidateAttach(stateDB, nodeId, candidateAttach); nil != err {
		return err
	}
	// Remove from pending expire tickets
	return t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
}

func (t *TicketPool) releaseTicket(stateDB vm.StateDB, candidateId discover.NodeID, candidateAttach *types.CandidateAttach, ticketId common.Hash, blockNumber *big.Int) (*types.Ticket, error) {
	log.Debug("Start executing releaseTicket", "nodeId", candidateId.String(), "ticketId", ticketId.Hex(), "epoch", candidateAttach.Epoch.Uint64(), "blockNumber", blockNumber.Uint64())
	ticket, err := t.GetTicket(stateDB, ticketId)
	if nil != err {
		return ticket, err
	}
	if ticket.TicketId == (common.Hash{}) {
		return nil, nil
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
	log.Debug("releaseTicket, start updating the total epoch of candidates", "nodeId", candidateId.String(), "totalEpoch", candidateAttach.Epoch.Uint64(), "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	candidateAttach.SubEpoch(ticket.CalcEpoch(blockNumber))
	log.Debug("releaseTicket, the end of the update candidate total epoch", "nodeId", candidateId.String(), "totalEpoch", candidateAttach.Epoch.Uint64(), "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	return ticket, transfer(stateDB, common.TicketPoolAddr, ticket.Owner, ticket.Deposit)
}

func (t *TicketPool) Notify(stateDB vm.StateDB, blockNumber *big.Int) error {
	log.Debug("Call Notify", "statedb addr", fmt.Sprintf("%p", stateDB))
	// Check expired tickets
	expireBlockNumber, ok := t.calcExpireBlockNumber(stateDB, blockNumber)
	log.Debug("Check expired tickets on Notify", "isOk", ok, "expireBlockNumber", expireBlockNumber.Uint64())
	if ok {
		if nodeIdList, err := t.handleExpireTicket(stateDB, expireBlockNumber, blockNumber); nil != err {
			log.Error("OutBlockNotice method handleExpireTicket error", "blockNumber", blockNumber.Uint64(), "err", err)
			return HandleExpireTicketErr
		} else {
			// Notify the candidate to update the list information after processing the expired ticket
			log.Debug("After processing the expired ticket, start updating the candidate list on Notify", "blockNumber", blockNumber.Uint64(), "nodeIdList", len(nodeIdList))
			if err := cContext.UpdateElectedQueue(stateDB, blockNumber, nodeIdList...); nil != err {
				log.Error("Failed to Update candidate when handleExpireTicket success on Notify", "err", err)
			}
		}
	}
	// Increase the total number of epoch for each candidate
	log.Debug("Increase the total number of epoch for each candidate on Notify", "blockNumber", blockNumber.Uint64())
	if err := t.calcCandidateEpoch(stateDB, blockNumber); nil != err {
		return err
	}
	return nil
}

func (t *TicketPool) calcCandidateEpoch(stateDB vm.StateDB, blockNumber *big.Int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	kindCandidateList := cContext.GetChosens(stateDB, 0)
	for _, candidateList := range kindCandidateList {
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
	}
	return nil
}

// Simple version of the lucky ticket algorithm
// According to the previous block Hash,
// find the first ticket Id which is larger than the Hash. If not found, the last ticket Id is taken.
func (t *TicketPool) SelectionLuckyTicket(stateDB vm.StateDB, nodeId discover.NodeID, blockHash common.Hash) (common.Hash, error) {
	log.Debug("Call SelectionLuckyTicket", "statedb addr", fmt.Sprintf("%p", stateDB))
	candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
	log.Debug("Start picking lucky tickets on SelectionLuckyTicket", "nodeId", nodeId.String(), "blockHash", blockHash.Hex(), "candidateTicketIds", len(candidateTicketIds))
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
	log.Debug("Pick out a lucky ticket on SelectionLuckyTicket", "index", index)
	luckyTicketId = decMap[decList[index]]
	log.Debug("End the selection of lucky tickets on SelectionLuckyTicket", "nodeId", nodeId.String(), "blockHash", blockHash.Hex(), "luckyTicketId", luckyTicketId.Hex(), "candidateTicketIds", len(candidateTicketIds))
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
	log.Debug("Call GetCandidaieTicketIds", "statedb addr", fmt.Sprintf("%p", stateDB))
	candidateTicketIds, err := stateDB.GetTicketCache(nodeId)
	if nil != err {
		return nil, err
	}
	return candidateTicketIds, nil
}

func (t *TicketPool) GetCandidatesTicketIds(stateDB vm.StateDB, nodeIds []discover.NodeID) (map[discover.NodeID][]common.Hash, error) {
	log.Debug("Call GetCandidateArrTicketIds", "statedb addr", fmt.Sprintf("%p", stateDB))
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
	log.Debug("Call GetCandidatesTicketCount", "statedb addr", fmt.Sprintf("%p", stateDB))
	result := make(map[discover.NodeID]int)
	if nil != nodeIds {
		for _, nodeId := range nodeIds {
			result[nodeId] = int(stateDB.TCount(nodeId))
		}
	}
	return result, nil
}

func (t *TicketPool) GetCandidateAttach(stateDB vm.StateDB, nodeId discover.NodeID) (*types.CandidateAttach, error) {
	log.Debug("Call GetCandidateAttach", "statedb addr", fmt.Sprintf("%p", stateDB))
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
		log.Error("Failed to encode candidateAttach object on setCandidateAttach", "key", string(nodeId.Bytes()), "value", candidateAttach.Epoch.Uint64(), "err", err)
		return SetCandidateAttachErr
	} else {
		setTicketPoolState(stateDB, CandidateAttachKey(nodeId.Bytes()), value)
	}
	return nil
}

func (t *TicketPool) GetCandidateEpoch(stateDB vm.StateDB, nodeId discover.NodeID) (uint64, error) {
	log.Debug("Call GetCandidateEpoch", "statedb addr", fmt.Sprintf("%p", stateDB))
	candidateAttach, err := t.GetCandidateAttach(stateDB, nodeId)
	if nil != err {
		return 0, err
	}
	return candidateAttach.Epoch.Uint64(), nil
}

func (t *TicketPool) GetTicketPrice(stateDB vm.StateDB) (*big.Int, error) {
	price := new(big.Int)
	if err := getTicketPoolState(stateDB, GetTicketPriceKey(), price); err!=nil {
		return nil, err
	}
	return price, nil
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

//func GetTicketPtr() *TicketPool {
//	return ticketPool
//}

func checkBalance(stateDB vm.StateDB, addr common.Address, amount *big.Int) bool {
	if stateDB.GetBalance(addr).Cmp(amount) >= 0 {
		return true
	}
	return false
}

func transfer(stateDB vm.StateDB, from common.Address, to common.Address, amount *big.Int) error {
	if !checkBalance(stateDB, from, amount) {
		log.Error("TicketPool not sufficient funds", "from", from.Hex(), "to", to.Hex(), "money", amount.Uint64())
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
	return addCommonPrefix(append(ExpireTicketPrefix, key...))
}

func CandidateAttachKey(key []byte) []byte {
	return addCommonPrefix(append(CandidateAttachPrefix, key...))
}

func GetSurplusQuantityKey() []byte {
	return addCommonPrefix(SurplusQuantityKey)
}

func GetTicketPriceKey() []byte {
	return addCommonPrefix(TicketPriceKey)
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


func SetTicketPrice(stateDB vm.StateDB, price *big.Int) {
	setTicketPoolState(stateDB, GetTicketPriceKey(), price.Bytes())
}

func (t *TicketPool) GetLowestTicketPrice() (*big.Int) {
	//get from cbft.json if not 10000*1e7
	return t.LowestTicketPrice
}

func (t *TicketPool) GetAdjustPriceCycle() (*big.Int){
	//get from cbft.json if not 10000
	return t.AdjustCycle
}

func (t *TicketPool) GetMaxPoolNumber() uint64 {
	return t.MaxCount
}
