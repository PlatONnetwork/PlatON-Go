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
	SurplusQuantityKey			= []byte("sq")
	ExpireTicketPrefix			= []byte("et")

	CandidateNotFindErr			= errors.New("The node has lost its candidacy")
	TicketNilErr				= errors.New("Ticket Insufficient quantity")

	EncodeTicketErr				= errors.New("Encode Ticket error")
	EncodePoolNumberErr			= errors.New("Encode SurplusQuantity error")
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
	t.setTicketInfo(stateDB, ticketId, ticket)
	// 调用候选池接口更新候选人信息
	return nil
}

func (t *TicketPool) setTicketInfo(stateDB vm.StateDB, ticketId common.Hash, ticket *types.Ticket) error {
	if value, err := rlp.EncodeToBytes(ticket); nil != err {
		log.Error("Failed to encode ticket object on setTicket", "key", ticketId.String(), "err", err)
		return EncodeTicketErr
	} else {
		setState(stateDB, ticketId.Bytes(), value)
	}
	return nil
}

func (t *TicketPool) setPoolNumber(stateDB vm.StateDB) error {
	t.SurplusQuantity--
	if value, err := rlp.EncodeToBytes(t.SurplusQuantity); nil != err {
		log.Error("Failed to encode surplusQuantity object on setPoolNumber", "key", string(SurplusQuantityKey), "err", err)
		return EncodePoolNumberErr
	} else {
		setState(stateDB, SurplusQuantityKey, value)
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
