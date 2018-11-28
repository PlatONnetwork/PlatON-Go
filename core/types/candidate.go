package types

import (
	"math/big"
	"Platon-go/p2p/discover"
	"Platon-go/common"
)

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
	Fee 			uint32

	// Voted ticket'id set
	//ticketPool		[]common.Hash
	// Voted ticket count
	//TCount    		uint64				`json:"tcount"`
	// Ticket age
	//Epoch			*big.Int			`json:"epoch"`
	// brokerage
	//Brokerage		uint64				`json:"brokerage"`
}

