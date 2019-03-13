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
	Remaining		uint32
}

func (t *Ticket) TotalDeposit() *big.Int {
	return new(big.Int).Mul(t.Deposit, new(big.Int).SetUint64(uint64(t.Remaining)))
}

func (t *Ticket) TotalEpoch(blockNumber *big.Int) uint64 {
	return t.CalcEpoch(blockNumber) * uint64(t.Remaining)
}

func (t *Ticket) CalcEpoch(blockNumber *big.Int) uint64 {
	result := new(big.Int).SetUint64(0)
	result.Sub(blockNumber, t.BlockNumber)
	return result.Uint64()
}

func (t *Ticket) SubRemaining() {
	if t.Remaining > 0 {
		t.Remaining--
	}
}

func (t *Ticket) DeepCopy() *Ticket {
	newDeposit := new(big.Int)
	newDeposit.Add(t.Deposit, newDeposit)
	newBlockNumber := new(big.Int)
	newBlockNumber.Add(t.BlockNumber, newBlockNumber)
	ticket := &Ticket{
		t.Owner,
		newDeposit,
		t.CandidateId,
		newBlockNumber,
		t.Remaining,
	}
	return ticket
}
