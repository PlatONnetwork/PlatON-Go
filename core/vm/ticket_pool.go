package vm

import (
	"Platon-go/common"
	"Platon-go/p2p/discover"
	"Platon-go/params"
	"errors"
	"math/big"
)

// error def
var (
	ErrIllegalDeposit = errors.New("Deposit balance not match")
)

const (
	VoteTicketEvent = "VoteTicketEvent"
)


// type ticketPool interface { }

type ticketContract struct {
	contract *Contract
	evm *EVM
}

func (t *ticketContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (t *ticketContract) Run(input []byte) ([]byte, error) {
	var command = map[string] interface{}{
		"VoteTicket" : t.VoteTicket,
		"GetTicketDetail" : t.GetTicketDetail,
		"GetCandidateTicketIds" : t.GetCandidateTicketIds,
		"GetEpoch" : t.GetEpoch,
		"GetPoolRemainder" : t.GetPoolRemainder,
		"GetTicketPrice": t.GetTicketPrice,
	}
	return execute(input, command)
}

// VoteTicket let a account buy tickets and vote to the chosen candidate.
func (t *ticketContract) VoteTicket(count *big.Int, price *big.Int, nodeId discover.NodeID) ([]byte, error) {
	totalPrice := new(big.Int).Mul(count, price)
	deposit := t.contract.value
	if totalPrice != deposit {
		return nil, ErrIllegalDeposit
	}
	return nil, nil
}

// GetTicketDetail returns the ticket info.
func (t *ticketContract) GetTicketDetail(ticketId common.Hash) ([]byte, error) {
	return nil, nil
}

// GetCandidateTicketIds returns the list of ticketId for the candidate.
func (t *ticketContract) GetCandidateTicketIds(nodeId discover.NodeID, blockNumber *big.Int) ([]byte, error) {
	return nil, nil
}

// GetEpoch returns the current ticket age for the candidate.
func (t *ticketContract) GetEpoch(nodeId discover.NodeID) ([]byte, error) {
	return nil, nil
}

// GetPoolRemainder returns the amount of remaining tikcets in the ticket pool.
func (t *ticketContract) GetPoolRemainder() ([]byte, error) {
	return nil, nil
}

// GetTicketPrice returns the current ticket price for the ticket pool.
func (t *ticketContract) GetTicketPrice() ([]byte, error) {
	return nil, nil
}