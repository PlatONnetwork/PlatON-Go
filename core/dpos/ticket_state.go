package depos

import (
	"Platon-go/common"
	"Platon-go/core/vm"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"math/big"
	"sync"
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
func NewTicketPool(configs *params.DposConfig) *TicketPool {
	ticketPool = &TicketPool{
		MaxCount:				0,
		SurplusQuantity:		0,
		ExpireBlockNumber:		0,
		lock:					&sync.RWMutex{},
	}
	return ticketPool
}

func(t *TicketPool) buyTicket(stateDB vm.StateDB, owner common.Hash, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) error {
	//candidate := candidatePool.immediateCandates[nodeId]

	return nil
}
