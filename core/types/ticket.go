package types

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
)

// ticket info
type Ticket struct {
	// Current owner of tickets
	Owner 			common.Address
	// Mortgage amount (margin)
	Deposit			*big.Int
	// candidate's server nodeId
	CandidateId 	discover.NodeID
	// current block height number when purchasing tickets
	BlockNumber 	*big.Int
	// The number of remaining tickets
	remaining		uint64
}

func (t *Ticket) CalcEpoch(blockNumber *big.Int) *big.Int {
	result := new(big.Int).SetUint64(0)
	result.Sub(blockNumber, t.BlockNumber)
	return result
}
