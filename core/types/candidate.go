package types

import (
	"math/big"
	"Platon-go/p2p/discover"
	"Platon-go/common"
)

type CandidateQueue []*Candidate

// candiate info
type Candidate struct {

	// Mortgage amount (margin)
	Deposit			*big.Int
	// Current block height number at the time of the mortgage
	BlockNumber 	*big.Int
	// Current tx'index at the time of the mortgage
	TxIndex 		uint32
	// candidate's server nodeId
	CandidateId 	discover.NodeID
	Host 			string
	Port 			string
	// Mortgage beneficiary's account address
	Owner 			common.Address
	// The account address of initiating a mortgaged
	From 			common.Address
	Extra 			string
	// brokerage   example: (fee/10000) * 100% == x%
	Fee 			uint64
	// Selected TicketId
	TicketId		common.Hash
}

type CandidateAttach struct {
	// Sum Ticket age
	Epoch			*big.Int			`json:"epoch"`
}

func (ca *CandidateAttach) AddEpoch(number *big.Int) {
	ca.Epoch.Add(ca.Epoch, number)
}

func (ca *CandidateAttach) SubEpoch(number *big.Int) {
	if ca.Epoch.Cmp(number) >= 0 && number.Uint64() > 0 {
		ca.Epoch.Sub(ca.Epoch, number)
	}
}

