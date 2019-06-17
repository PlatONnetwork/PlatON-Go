package common

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
)


// The Candidate info
type Candidate struct {

	NodeId 			discover.NodeID
	// The account used to initiate the staking
	StakingAddress 	common.Address
	// The account recieve the block rewards and the staking rewards
	BenifitAddress  common.Address
	// Block height at the time of staking
	StakingBlockNum uint64
	// The tx index at the time of staking
	StakingTxIndex  uint32
	// The epoch number at staking or edit
	StakingEpoch 	uint64
	// All vons of staking and delegated
	Shares			*big.Int
	// The staking vons  is circulating for effective epoch (in effect)
	Released 		*big.Int
	// The staking vons  is circulating for hesitant epoch (in hesitation)
	ReleasedTmp 	*big.Int
	// The staking vons  is locked for effective epoch (in effect)
	LockRepo		*big.Int
	// The staking vons  is locked for hesitant epoch (in hesitation)
	LockRepoTmp		*big.Int
	// Positive and negative signs:
	// Is it an increase or a decrease? 0: increase; 1: decrease
	Signs 			uint8

	// The candiate status
	Status 			uint32

	// Node desc
	*Description
}


type Description struct {
	// The Candidate Node's Name  (with a length limit)
	NodeName	string
	// External Id for the third party to pull the node description (with length limit)
	ExternalId	string
	// The third-party home page of the node (with a length limit)
	Website		string
	// Description of the node (with a length limit)
	Details		string
}

// the Validator info
// They are Simplified Candidate
// They are consensus nodes and Epoch nodes snapshot
type Validator struct {
	NodeAddress 	common.Address
	NodeId 			discover.NodeID
	// The weight
	// NOTE:
	// converted from the weight of Candidate is: (Int.Max - candidate.shares) + blocknum + txindex
	StakingWeight	string
	// Validator's term in the consensus round
	ValidatorTerm 	uint32
}

// some consensus round validators or current epoch validators
type Validator_array struct {
	// the round start blocknumber or epoch start blocknumber
	Start 		uint64
	// the round end blocknumber or epoch blocknumber
	End 		uint64
	// the round validators or epoch validators
	Arr			[]*Validator
}


// the Delegate information
type Delegation struct {

	// The epoch number at delegate or edit
	DelegateEpoch 	uint64

	// Positive and negative signs:
	// Is it an increase or a decrease? 0: increase; 1: decrease
	Signs 			uint8

	// The delegate vons  is circulating for effective epoch (in effect)
	Released 		*big.Int
	// The delegate vons  is circulating for hesitant epoch (in hesitation)
	ReleasedTmp 	*big.Int
	// The delegate vons  is locked for effective epoch (in effect)
	LockRepo		*big.Int
	// The delegate vons  is locked for hesitant epoch (in hesitation)
	LockRepoTmp		*big.Int

}