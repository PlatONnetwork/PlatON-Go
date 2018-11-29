package types

import (
	"Platon-go/common"
	"Platon-go/p2p/discover"
	"math/big"
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
}

func (t *Ticket) CalcEpoch(blockNumber *big.Int) *big.Int {
	result := new(big.Int).SetUint64(0)
	result.Sub(blockNumber, t.BlockNumber)
	result.Add(result, t.BlockNumber)
	return result
}
