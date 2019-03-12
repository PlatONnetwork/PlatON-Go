package pposm

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
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
	CandidateNotFindErr 	  = errors.New("The Candidate not find")
	CandidateNilTicketErr     = errors.New("This candidate has no ticket")
	TicketPoolBalanceErr      = errors.New("TicketPool not sufficient funds")
	TicketIdNotFindErr        = errors.New("TicketId not find")
	HandleExpireTicketErr     = errors.New("Failure to deal with expired tickets")
	GetCandidateAttachErr     = errors.New("Get CandidateAttach error")
	SetCandidateAttachErr     = errors.New("Update CandidateAttach error")
	VoteTicketErr        	  = errors.New("Voting failed")
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
	var ticketPrice *big.Int
	if price, ok := new(big.Int).SetString(configs.TicketConfig.TicketPrice, 10); !ok {
		ticketPrice, _ = new(big.Int).SetString("100000000000000000000", 10)
	} else {
		ticketPrice = price
	}

	ticketPool := &TicketPool{
		TicketPrice:       ticketPrice,
		MaxCount:          configs.TicketConfig.MaxCount,
		ExpireBlockNumber: configs.TicketConfig.ExpireBlockNumber,
		lock:              &sync.Mutex{},
	}
	return ticketPool
}

func (t *TicketPool) VoteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) (uint64, error) {
	log.Debug("Call Voting", "statedb addr", fmt.Sprintf("%p", stateDB))
	log.Info("Start Voting,VoteTicket", "owner", owner.Hex(), "voteNumber", voteNumber, "price", deposit.Uint64(), "nodeId", nodeId.String(), "blockNumber", blockNumber.Uint64())
	successCount, err := t.voteTicket(stateDB, owner, voteNumber, deposit, nodeId, blockNumber)
	if nil != err {
		log.Error("Voting failed", "nodeId", nodeId.String(), "voteNumber", voteNumber, "successNum", successCount, "err", err)
		return successCount, err
	}
	// Voting completed, candidates reordered
	log.Debug("Successfully voted to start updating the list of candidates,VoteTicket", "successNum", successCount)
	if err := cContext.UpdateElectedQueue(stateDB, blockNumber, nodeId); nil != err {
		log.Error("Failed to Update candidate when voteTicket success", "err", err)
	}
	log.Debug("Successful vote, candidate list updated successfully,VoteTicket", "successNum", successCount)
	return successCount, nil
}

func (t *TicketPool) voteTicket(stateDB vm.StateDB, owner common.Address, voteNumber uint64, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) (uint64, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	// check ticket pool count
	surplusQuantity := t.GetPoolNumber(stateDB)
	log.Debug("Execute voteTicket", "surplusQuantity", surplusQuantity, "voteNumber", voteNumber, "blockNumber", blockNumber.Uint64())
	if surplusQuantity == 0 {
		log.Error("Ticket Insufficient quantity")
		return 0, TicketPoolNilErr
	}
	if surplusQuantity < voteNumber {
		voteNumber = surplusQuantity
	}
	log.Debug("Start circular voting", "nodeId", nodeId.String(), "voteNumber", voteNumber)

	ticketId := stateDB.TxHash()
	ticket := &types.Ticket{
		owner,
		deposit,
		nodeId,
		blockNumber,
		voteNumber,
	}
	if err := t.setTicket(stateDB, ticketId, ticket); nil != err {
		return voteNumber, err
	}
	log.Debug("SetTicket succeeds, start recording tickets to expire", "blockNumber", blockNumber.Uint64(), "ticketId", ticketId.Hex())
	if err := t.recordExpireTicket(stateDB, blockNumber, ticketId); err != nil {
		return voteNumber, err
	}
	log.Debug("Record the success of the ticket to expire, and start reducing the number of tickets", "surplusQuantity", surplusQuantity)
	if err := t.setPoolNumber(stateDB, surplusQuantity-voteNumber); err != nil {
		return voteNumber, err
	}
	stateDB.GetPPOSCache().AppendTicket(nodeId, ticketId, ticket)
	log.Debug("Voting SUCCUESS !!!!!!  Reduce the remaining amount of the ticket pool successfully", "surplusQuantity", t.GetPoolNumber(stateDB), "nodeId", nodeId.String())
	return voteNumber, nil
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
	if value, err := stateDB.GetPPOSCache().GetExpireTicket(blockNumber); nil != err {
		return nil, err
	} else {
		return value, nil
	}
}

// In the current block,
// the ticket id is placed in the value slice with the block height as the key to find the expired ticket.
func (t *TicketPool) recordExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, ticketId common.Hash) error {
	return stateDB.GetPPOSCache().SetExpireTicket(blockNumber, ticketId)
}

func (t *TicketPool) removeExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, ticketId common.Hash) error {
	return stateDB.GetPPOSCache().RemoveExpireTicket(blockNumber, ticketId)
}

func (t *TicketPool) handleExpireTicket(stateDB vm.StateDB, expireBlockNumber *big.Int, currentBlockNumber *big.Int) ([]discover.NodeID, error) {
	t.lock.Lock()
	defer t.lock.Unlock()
	ticketIdList, err := t.GetExpireTicketIds(stateDB, expireBlockNumber)
	if err != nil {
		return nil, err
	}
	log.Info("Pending ticket to be processed", "amount", len(ticketIdList), "expireBlockNumber", expireBlockNumber.Uint64(), "currentBlockNumber", currentBlockNumber.Uint64())
	candidateAttachMap := make(map[discover.NodeID]bool)
	changeNodeIdList := make([]discover.NodeID, 0)
	for _, ticketId := range ticketIdList {
		ticket, err := t.GetTicket(stateDB, ticketId)
		if err != nil {
			return nil, err
		}
		if ticket == nil {
			continue
		}
		_, ok := candidateAttachMap[ticket.CandidateId]
		if !ok {
			candidateAttachMap[ticket.CandidateId] = true
			changeNodeIdList = append(changeNodeIdList, ticket.CandidateId)
		}
		if _, err := t.releaseTxTicket(stateDB, ticket.CandidateId, ticketId, currentBlockNumber); nil != err {
			return changeNodeIdList, err
		}
	}
	return changeNodeIdList, nil
}

// Get ticket list
func (t *TicketPool) GetTicketList(stateDB vm.StateDB, ticketIds []common.Hash) ([]*types.Ticket, error) {
	log.Debug("Call GetTickList", "statedb addr", fmt.Sprintf("%p", stateDB))
	var tickets []*types.Ticket
	for _, ticketId := range ticketIds {
		ticket, err := t.GetTicket(stateDB, ticketId)
		if nil != err || ticket == nil {
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
	if value, err := stateDB.GetPPOSCache().GetTicketInfo(ticketId); nil != err {
		return nil, err
	} else {
		return value, nil
	}
}

func (t *TicketPool) setTicket(stateDB vm.StateDB, ticketId common.Hash, ticket *types.Ticket) error {
	if err := stateDB.GetPPOSCache().SetTicketInfo(ticketId, ticket); nil != err {
		log.Error("Failed to setTicket", "key", ticketId.Hex(), "err", err)
		return err
	}
	return nil
}

func (t *TicketPool) DropReturnTicket(stateDB vm.StateDB, blockNumber *big.Int, nodeIds ...discover.NodeID) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	log.Debug("Call DropReturnTicket", "statedb addr", fmt.Sprintf("%p", stateDB))
	log.Info("Start processing tickets for the drop list on DropReturnTicket", "candidateNum", len(nodeIds), "blockNumber", blockNumber.Uint64())
	for _, nodeId := range nodeIds {
		if  nodeId == (discover.NodeID{}) {
			continue
		}
		candidateTicketIds, err := t.GetCandidateTicketIds(stateDB, nodeId)
		if nil != err {
			return err
		}
		if len(candidateTicketIds) == 0 {
			continue
		}
		epoch, err := t.GetCandidateEpoch(stateDB, nodeId)
		if nil != err {
			return err
		}
		ticketCount := t.GetCandidateTicketCount(stateDB, nodeId)
		log.Debug("Delete candidate ticket collection on DropReturnTicket", "nodeId", nodeId.String(), "ticketSize", ticketCount, "epoch", epoch)
		if err := stateDB.GetPPOSCache().RemoveTicketDependency(nodeId); err != nil {
			return err
		}
		surplusQuantity := t.GetPoolNumber(stateDB)
		log.Debug("Start reducing the number of tickets on DropReturnTicket", "surplusQuantity", surplusQuantity, "candidateTicketIds", ticketCount)
		if err := t.setPoolNumber(stateDB, surplusQuantity + ticketCount); err != nil {
			return err
		}
		log.Debug("Start processing each invalid ticket on DropReturnTicket", "nodeId", nodeId.String(), "ticketSize", ticketCount)
		for _, ticketId := range candidateTicketIds {
			ticket, err := t.GetTicket(stateDB, ticketId)
			if nil != err {
				return err
			}
			if ticket == nil {
				continue
			}
			log.Debug("Start transfer on DropReturnTicket", "nodeId", nodeId.String(), "ticketId", ticketId.Hex(), "deposit", ticket.Deposit, "remaining", ticket.Remaining)
			if err := transfer(stateDB, common.TicketPoolAddr, ticket.Owner, new(big.Int).Mul(ticket.Deposit, new(big.Int).SetUint64(ticket.Remaining))); nil != err {
				return err
			}
			stateDB.GetPPOSCache().RemoveTicket(nodeId, ticketId)
			if err := t.removeExpireTicket(stateDB, blockNumber, ticketId); nil != err {
				return err
			}
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
	ticket, err := t.releaseTicket(stateDB, nodeId, ticketId, blockNumber)
	if nil != err {
		return err
	}
	if ticket == nil {
		return TicketIdNotFindErr
	}
	if ticket.Remaining == 0 {
		// Remove from pending expire tickets
		return t.removeExpireTicket(stateDB, ticket.BlockNumber, ticketId)
	}
	return nil
}

func (t *TicketPool) releaseTicket(stateDB vm.StateDB, candidateId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) (*types.Ticket, error) {
	log.Debug("Start executing releaseTicket", "nodeId", candidateId.String(), "ticketId", ticketId.Hex(), "blockNumber", blockNumber.Uint64())
	ticket, err := t.GetTicket(stateDB, ticketId)
	if nil != err {
		return ticket, err
	}
	if ticket == nil {
		return nil, nil
	}
	log.Debug("releaseTicket,Start Update", "nodeId", candidateId.String(), "ticketId", ticketId.Hex())
	if err := stateDB.GetPPOSCache().SubTicket(candidateId, ticketId); err != nil {
		return ticket, err
	}
	log.Debug("releaseTicket, end update", "nodeId", candidateId.String())
	surplusQuantity := t.GetPoolNumber(stateDB)
	log.Debug("releaseTicket, start to update the ticket pool", "surplusQuantity", surplusQuantity)
	if err := t.addPoolNumber(stateDB); err != nil {
		return ticket, err
	}
	surplusQuantity = t.GetPoolNumber(stateDB)
	log.Debug("releaseTicket, end the update ticket pool", "surplusQuantity", surplusQuantity)
	epoch, err := t.GetCandidateEpoch(stateDB, candidateId)
	if nil != err {
		return ticket, err
	}
	log.Debug("releaseTicket, start updating the total epoch of candidates", "nodeId", candidateId.String(), "totalEpoch", epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	dependency, err := stateDB.GetPPOSCache().GetTicketDependency(candidateId)
	if nil != err {
		return ticket, err
	}
	dependency.SubAge(ticket.CalcEpoch(blockNumber))
	epoch, err = t.GetCandidateEpoch(stateDB, candidateId)
	if nil != err {
		return ticket, err
	}
	log.Debug("releaseTicket, the end of the update candidate total epoch", "nodeId", candidateId.String(), "totalEpoch", epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	return ticket, transfer(stateDB, common.TicketPoolAddr, ticket.Owner, ticket.Deposit)
}

func (t *TicketPool) releaseTxTicket(stateDB vm.StateDB, candidateId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) (*types.Ticket, error) {
	log.Debug("Start executing releaseTicket", "nodeId", candidateId.String(), "ticketId", ticketId.Hex(), "blockNumber", blockNumber.Uint64())
	ticket, err := t.GetTicket(stateDB, ticketId)
	if nil != err {
		return ticket, err
	}
	if ticket == nil {
		return nil, nil
	}
	log.Debug("releaseTicket,Start Update", "nodeId", candidateId.String(), "ticketId", ticketId.Hex())
	if err := stateDB.GetPPOSCache().RemoveTicket(candidateId, ticketId); err != nil {
		return ticket, err
	}
	log.Debug("releaseTicket, end update", "nodeId", candidateId.String())
	if err := t.removeExpireTicket(stateDB, blockNumber, ticketId); nil != err {
		return ticket, err
	}
	surplusQuantity := t.GetPoolNumber(stateDB)
	log.Debug("releaseTicket, start to update the ticket pool", "surplusQuantity", surplusQuantity)
	if err := t.setPoolNumber(stateDB, surplusQuantity + ticket.Remaining); err != nil {
		return ticket, err
	}
	surplusQuantity = t.GetPoolNumber(stateDB)
	log.Debug("releaseTicket, end the update ticket pool", "surplusQuantity", surplusQuantity)
	epoch, err := t.GetCandidateEpoch(stateDB, candidateId)
	if nil != err {
		return ticket, err
	}
	log.Debug("releaseTicket, start updating the total epoch of candidates", "nodeId", candidateId.String(), "totalEpoch", epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	dependency, err := stateDB.GetPPOSCache().GetTicketDependency(candidateId)
	if nil != err {
		return ticket, err
	}
	dependency.SubAge(new(big.Int).Mul(ticket.CalcEpoch(blockNumber), new(big.Int).SetUint64(ticket.Remaining)))
	epoch, err = t.GetCandidateEpoch(stateDB, candidateId)
	if nil != err {
		return ticket, err
	}
	log.Debug("releaseTicket, the end of the update candidate total epoch", "nodeId", candidateId.String(), "totalEpoch", epoch, "blockNumber", blockNumber.Uint64(), "ticketBlockNumber", ticket.BlockNumber.Uint64())
	return ticket, transfer(stateDB, common.TicketPoolAddr, ticket.Owner, new(big.Int).Mul(ticket.Deposit, new(big.Int).SetUint64(ticket.Remaining)))
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
			if len(nodeIdList) > 0 {
				if err := cContext.UpdateElectedQueue(stateDB, blockNumber, nodeIdList...); nil != err {
					log.Error("Failed to Update candidate when handleExpireTicket success on Notify", "err", err)
				}
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
	candidateList := cContext.GetChosens(stateDB, 0)
	for _, candidate := range candidateList {
		epoch, err := t.GetCandidateEpoch(stateDB, candidate.CandidateId)
		if nil != err {
			return err
		}
		// Get the total number of votes, increase the total epoch
		ticketCount := stateDB.GetPPOSCache().GetCandidateTicketCount(candidate.CandidateId)
		log.Debug("increase the total epoch", "candidateId", candidate.CandidateId.String(), "blockNumber", blockNumber.Uint64(), "ticketCount", ticketCount, "epoch", epoch)
		if ticketCount > 0 {
			dependency, err := stateDB.GetPPOSCache().GetTicketDependency(candidate.CandidateId)
			if nil != err {
				return err
			}
			dependency.AddAge(new(big.Int).SetUint64(ticketCount))
			epoch, err = t.GetCandidateEpoch(stateDB, candidate.CandidateId)
			if nil != err {
				return err
			}
			log.Debug("increase the total epoch success", "candidateId", candidate.CandidateId.String(), "blockNumber", blockNumber.Uint64(), "ticketCount", ticketCount, "epoch", epoch)
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

func (t *TicketPool) addPoolNumber(stateDB vm.StateDB) error {
	surplusQuantity := t.GetPoolNumber(stateDB)
	if surplusQuantity == t.MaxCount {
		return TicketPoolOverflowErr
	}
	surplusQuantity++
	return t.setPoolNumber(stateDB, surplusQuantity)
}

func (t *TicketPool) subPoolNumber(stateDB vm.StateDB) error {
	surplusQuantity := t.GetPoolNumber(stateDB)
	if surplusQuantity == 0 {
		return TicketPoolNilErr
	}
	surplusQuantity--
	return t.setPoolNumber(stateDB, surplusQuantity)
}

func (t *TicketPool) setPoolNumber(stateDB vm.StateDB, surplusQuantity uint64) error {
	return stateDB.GetPPOSCache().SetTotalRemain(int(surplusQuantity))
}

func (t *TicketPool) GetPoolNumber(stateDB vm.StateDB) uint64 {
	if val := stateDB.GetPPOSCache().GetTotalRemian(); val >= 0 {
		return uint64(val)
	} else {
		// Default initialization values
		return t.MaxCount
	}
}

func (t *TicketPool) GetCandidateTicketIds(stateDB vm.StateDB, nodeId discover.NodeID) ([]common.Hash, error) {
	log.Debug("Call GetCandidaieTicketIds", "statedb addr", fmt.Sprintf("%p", stateDB))
	candidateTicketIds, err := stateDB.GetPPOSCache().GetCandidateTxHashs(nodeId)
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
			if nil != err || nil == ticketIds {
				continue
			}
			result[nodeId] = ticketIds
		}
	}
	return result, nil
}

func (t *TicketPool) GetCandidateTicketCount(stateDB vm.StateDB, nodeId discover.NodeID) uint64 {
	return stateDB.GetPPOSCache().GetCandidateTicketCount(nodeId)
}

func (t *TicketPool) GetCandidatesTicketCount(stateDB vm.StateDB, nodeIds []discover.NodeID) (map[discover.NodeID]int, error) {
	log.Debug("Call GetCandidatesTicketCount", "statedb addr", fmt.Sprintf("%p", stateDB))
	result := make(map[discover.NodeID]int)
	if nil != nodeIds {
		for _, nodeId := range nodeIds {
			result[nodeId] = int(stateDB.GetPPOSCache().GetCandidateTicketCount(nodeId))
		}
	}
	return result, nil
}

func (t *TicketPool) setCandidateEpoch(stateDB vm.StateDB, nodeId discover.NodeID, age *big.Int) error {
	return stateDB.GetPPOSCache().SetCandidateTicketAge(nodeId, age)
}

func (t *TicketPool) GetCandidateEpoch(stateDB vm.StateDB, nodeId discover.NodeID) (uint64, error) {
	log.Debug("Call GetCandidateEpoch", "statedb addr", fmt.Sprintf("%p", stateDB))
	if value, err := stateDB.GetPPOSCache().GetCandidateTicketAge(nodeId); nil != err {
		log.Error("GetCandidateEpoch error", "key", nodeId.String(), "err", err)
		return 0, err
	} else {
		if value != nil {
			return value.Uint64(), nil
		} else {
			return 0, nil
		}
	}
}

func (t *TicketPool) GetTicketPrice(stateDB vm.StateDB) (*big.Int, error) {
	return t.TicketPrice, nil
}

// Save the hash value of the current state of the ticket pool
func (t *TicketPool) CommitHash(stateDB vm.StateDB) error {
	//hash, err := ticketcache.Hash(stateDB.TicketCaceheSnapshot())
	//if nil != err {
	//	return err
	//}
	//setTicketPoolState(stateDB, addCommonPrefix(TicketPoolHashKey), hash.Bytes())
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
