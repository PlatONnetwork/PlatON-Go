package types

import (
	"Platon-go/common"
	"Platon-go/p2p/discover"
	"math/big"
)

const (
	Normal uint8 = iota + 1
	Selected
	Expired
	Invalid
)

// ticket info
type Ticket struct {
	// ticket id
	TicketId 		common.Hash
	// Current owner of tickets
	Owner 			common.Address
	// Mortgage amount (margin)
	Deposit			*big.Int
	// candidate's server nodeId
	CandidateId 	discover.NodeID
	// current block height number when purchasing tickets
	BlockNumber 	*big.Int
	// Ticket state
	// 1 -> Normal
	// 2 -> Selected
	// 3 -> Expired
	// 4 -> Invalid
	State 			uint8
	// Block height when released
	RBlockNumber	*big.Int
}

func (t *Ticket) CalcEpoch(blockNumber *big.Int) *big.Int {
	result := new(big.Int).SetUint64(0)
	result.Sub(blockNumber, t.BlockNumber)
	return result
}

func (t *Ticket) SetNormal() {
	t.State = Normal
}

func (t *Ticket) SetSelected(blockNumber *big.Int) {
	t.State = Selected
	t.RBlockNumber = blockNumber
}

func (t *Ticket) SetExpired(blockNumber *big.Int) {
	t.State = Expired
	t.RBlockNumber = blockNumber
}

func (t *Ticket) SetInvalid(blockNumber *big.Int) {
	t.State = Invalid
	t.RBlockNumber = blockNumber
}
