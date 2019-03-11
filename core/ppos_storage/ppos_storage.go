package ppos_storage

import (
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"math/big"
)

type candidateStorage map[discover.NodeID]*types.Candidate
type refundStorage map[discover.NodeID]types.CandidateQueue



type candidate_temp struct {
	// previous witness
	pres 	types.CandidateQueue
	// current witness
	currs 	types.CandidateQueue
	// next witness
	nexts 	types.CandidateQueue
	//immediate
	imms 	types.CandidateQueue
	// reserve
	res 	types.CandidateQueue
	// refund
	refunds refundStorage
}


type ticketDependency struct {
	// ticket age
	Age  uint32
	// ticket count
	Num  uint32
	// ticketIds
	tIds []common.Hash
}

type ticket_temp struct {
	// total remian  k-v
	Sq  uint32
	// ticketInfo  map[txHash]ticketInfo
	Infos map[common.Hash]*types.Ticket
	// ExpireTicket  map[blockNumber]txHash
	Ets map[string]common.Hash
	// ticket's attachment  of node
	Dependencys map[discover.NodeID]*ticketDependency
}

type Ppos_storage struct {
	c_storage *candidate_temp
	t_storage *ticket_temp
}


/** candidate related func */

// Get CandidateQueue
func (p *Ppos_storage) GetCandidateQueue (flag int) (types.CandidateQueue, error){
	return nil, nil
}

// Set CandidateQueue
func (p *Ppos_storage) SetCandidateQueue(queue types.CandidateQueue, flag int) error {
	return nil
}

// Get Refund
func (p *Ppos_storage) GetRefund (nodeId discover.NodeID) (refundStorage, error) {
	return nil, nil
}

// Set Refund
func (p *Ppos_storage) SetRefund (nodeId discover.NodeID, refund refundStorage) error {
	return nil
}

/** ticket related func */

// Get total remian
func (p *Ppos_storage) GetTotalRemian() uint32 {
	return 0
}

// Set total remain
func (p *Ppos_storage) SetTotalRemain (count uint32) error {
	return nil
}


// Get TicketInfo
func(p *Ppos_storage) GetTicketInfo(txHash common.Hash) (*types.Ticket, error) {
	return nil, nil
}

//Set TicketInfo
func(p *Ppos_storage) SetTicketInfo(ticket *types.Ticket) error {
	return nil
}

//GetTiketArr
func (p *Ppos_storage) GetTicketArr(txHashs ... common.Hash) ([]*types.Ticket, error) {
	return nil, nil
}

//Get ExpireTicket
func (p *Ppos_storage) GetExpireTicket(blockNumber *big.Int) ([]common.Hash, error) {
	return nil, nil
}

// Set ExpireTicket
func (p *Ppos_storage) SetExpireTicket (blockNumber *big.Int, txHash common.Hash) error {
	return nil
}

//Get ticket dependency
func (p *Ppos_storage) GetTicketDependency (nodeId discover.NodeID) (*ticketDependency, error) {
	return nil, nil
}

// Set ticket dependency
func (p *Ppos_storage) SetTicketDependency (nodeId discover.NodeID, ependency *ticketDependency) error {
	return nil
}
