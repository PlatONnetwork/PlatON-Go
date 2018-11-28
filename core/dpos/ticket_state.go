package depos

import (
	"Platon-go/common"
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
	"sync"
)

var (
	// Remaining number key
	SurplusQuantityKey			= []byte("sq")
	// Expire ticket prefix
	ExpireTicketPrefix			= []byte("et")

	CandidateNotFindErr			= errors.New("The node has lost its candidacy")
	TicketNilErr				= errors.New("Ticket Insufficient quantity")
	EncodeTicketErr				= errors.New("Encode Ticket error")
	EncodePoolNumberErr			= errors.New("Encode SurplusQuantity error")
	DecodeTicketErr				= errors.New("Decode Ticket error")
	DecodePoolNumberErr			= errors.New("Decode SurplusQuantity error")
	RecordExpireTicketErr		= errors.New("Record Expire Ticket error")
	CandidateNotFindTicketErr	= errors.New("The candidate no longer has this ticket")
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
func NewTicketPool(configs *params.DposConfig, candidatePool *CandidatePool) *TicketPool {
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
	candidate, err := candidatePool.GetCandidate(stateDB, nodeId)
	if err != nil {
		log.Error("GetCandidate error", err)
		return err
	}
	if candidate == nil {
		log.Error("The node has lost its candidacy", err)
		return CandidateNotFindErr
	}
	ticketId, err := generateTicket()
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
	candidate.TicketPool = append(candidate.TicketPool, ticketId)
	candidate.TCount++
	candidate.Epoch = candidate.Epoch.Add(candidate.Epoch, blockNumber)
	if err := t.setTicket(stateDB, ticketId, ticket); err != nil {
		return err
	}
	if err := t.recordExpireTicket(stateDB, blockNumber, ticketId); err != nil {
		return err
	}
	if err := t.candidatePool.UpdateCandidateTicket(stateDB, ticket.CandidateId, candidate); err != nil {
		return err
	}
	if err := t.subPoolNumber(stateDB); err != nil {
		return err
	}
	return nil
}

//
func (t *TicketPool) GetExpireTicket(stateDB vm.StateDB, blockNumber *big.Int) ([]common.Hash, error) {
	var expireTickets []common.Hash
	/*if val := stateDB.GetState(common.TicketAddr, ExpireTicketKey((*blockNumber).Bytes())); len(val) > 0 {
		if err := rlp.DecodeBytes(val, &expireTickets); nil != err {
			log.Error("Decode ExpireTicket error", "key", *blockNumber, "err", err)
			return nil, err
		}
	}*/
	if err := getState(stateDB, ExpireTicketKey((*blockNumber).Bytes()), &expireTickets); nil != err {
		return nil, err
	}
	return expireTickets, nil
}

// 在当前区块投入的票，则把票id放入，以块高为key的value切片中，以便查找过期票
func (t *TicketPool) recordExpireTicket(stateDB vm.StateDB, blockNumber *big.Int, ticketId common.Hash) error {
	expireTickets, err := t.GetExpireTicket(stateDB, blockNumber)
	if err != nil {
		log.Error("recordExpireTicket error", "key", *blockNumber, "err", err)
		return RecordExpireTicketErr
	}
	expireTickets = append(expireTickets, ticketId)
	if value, err := rlp.EncodeToBytes(expireTickets); nil != err {
		log.Error("Failed to encode ticketid object on recordExpireTicket", "key", *blockNumber, "value", ticketId.String(), "err", err)
		return EncodeTicketErr
	} else {
		setState(stateDB, ExpireTicketKey((*blockNumber).Bytes()), value)
	}
	return nil
}

func (t *TicketPool) HandleExpireTicket(stateDB vm.StateDB, blockNumber *big.Int) error {
	ticketIdList, err := t.GetExpireTicket(stateDB, blockNumber)
	if err != nil {
		return err
	}
	for _, ticketId := range ticketIdList {
		ticket, err := t.GetTicket(stateDB, ticketId)
		if err != nil {
			return err
		}
		if err := t.ReleaseTicket(stateDB, ticket.CandidateId, ticketId, blockNumber); err != nil {
			//return err
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
	/*if err := rlp.DecodeBytes(stateDB.GetState(common.TicketAddr, ticketId.Bytes()), &ticket); nil != err {
		log.Error("Decode Ticket error", "key", ticketId, "err", err)
		return nil, DecodeTicketErr
	}*/
	if err := getState(stateDB, ticketId.Bytes(), ticket); nil != err {
		return nil, DecodeTicketErr
	}
	return ticket, nil
}

func (t *TicketPool) setTicket(stateDB vm.StateDB, ticketId common.Hash, ticket *types.Ticket) error {
	if value, err := rlp.EncodeToBytes(ticket); nil != err {
		log.Error("Failed to encode ticket object on setTicket", "key", ticketId.String(), "err", err)
		return EncodeTicketErr
	} else {
		setState(stateDB, ticketId.Bytes(), value)
	}
	return nil
}

func (t *TicketPool) ReleaseTicket(stateDB vm.StateDB, nodeId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) error {
	t.lock.Lock()
	defer t.lock.Unlock()
	candidate, err := candidatePool.GetCandidate(stateDB, nodeId)
	if err != nil {
		log.Error("GetCandidate error", err)
		return err
	}
	if candidate == nil {
		log.Error("The node has lost its candidacy", "err", err)
		return CandidateNotFindErr
	}
	exists := false
	for index, tempTicketId := range candidate.TicketPool {
		if tempTicketId == ticketId {
			exists = true
			candidate.TicketPool = removeCandidateTicket(index, candidate.TicketPool)
			break
		}
	}
	if !exists {
		log.Error("The candidate no longer has this ticket", "err", err)
		return CandidateNotFindTicketErr
	}
	ticket, err := t.GetTicket(stateDB, ticketId)
	if nil != err {
		return err
	}
	if err := t.addPoolNumber(stateDB); err != nil {
		return err
	}
	num := blockNumber.Sub(blockNumber, ticket.BlockNumber)
	candidate.Epoch = candidate.Epoch.Sub(candidate.Epoch, num.Add(num, ticket.BlockNumber))
	candidate.TCount--
	t.candidatePool.UpdateCandidateTicket(stateDB, candidate.CandidateId, candidate)
	return nil
}

func removeCandidateTicket(index int, ticketPool []common.Hash) []common.Hash {
	start := ticketPool[:index]
	end := ticketPool[index+1:]
	return append(start, end...)
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
		setState(stateDB, SurplusQuantityKey, value)
	}
	return nil
}

func (t *TicketPool) GetPoolNumber(stateDB vm.StateDB) (uint64, error) {
	var surplusQuantity uint64
	/*if err := rlp.DecodeBytes(stateDB.GetState(common.TicketAddr, SurplusQuantityKey), &surplusQuantity); nil != err {
		log.Error("Decode ticket pool SurplusQuantity error", "key", string(SurplusQuantityKey), "err", err)
		return 0, DecodePoolNumberErr
	}*/
	if err := getState(stateDB, SurplusQuantityKey, &surplusQuantity); nil != err {
		return 0, DecodePoolNumberErr
	}
	return surplusQuantity, nil
}

func getState(stateDB vm.StateDB, key []byte, result interface{}) error {
	if val := stateDB.GetState(common.TicketAddr, key); len(val) > 0 {
		if err := rlp.DecodeBytes(val, result); nil != err {
			log.Error("Decode Data error", "key", string(key), "err", err)
			return err
		}
	}
	return nil
}

func setState(stateDB vm.StateDB, key []byte, val []byte) {
	stateDB.SetState(common.TicketAddr, key, val)
}

func generateTicket() (common.Hash, error) {
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
