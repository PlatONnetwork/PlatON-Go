package xcom

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
)

const (
	/**
	######   ######   ######   ######
	#	  THE CANDIDATE  STATUS     #
	######   ######   ######   ######
	*/
	Invalided = 1 << iota // 0001: The current candidate withdraws from the staking qualification
	Slashed               // 0010: The candidate was slashed
	NotEnough             // 0100: The current candidate's von does not meet the minimum staking threshold
	Valided   = 0         // 0000: The current candidate is in force
)

func IsCan_Valid(status uint32) bool {
	return status&Valided == Valided
}

func IsCan_Invalid(status uint32) bool {
	return status&Invalided == Invalided
}

func IsCan_PureInvalid(status uint32) bool {
	return status&Invalided == status|Invalided
}

func IsCan_Slashed(status uint32) bool {
	return status&Slashed == Slashed
}

func IsCan_PureSlashed(status uint32) bool {
	return status&Slashed == status|Slashed
}

func IsCan_NotEnough(status uint32) bool {
	return status&NotEnough == NotEnough
}

func IsCand_idatePureNotEnough(status uint32) bool {
	return status&NotEnough == status|NotEnough
}

func IsCan_Invalid_Slashed(status uint32) bool {
	return status&(Invalided|Slashed) == (Invalided | Slashed)
}

func IsCan_Invalid_NotEnough(status uint32) bool {
	return status&(Invalided|NotEnough) == (Invalided | NotEnough)
}

func IsCan_Invalid_Slashed_NotEnough(status uint32) bool {
	return status&(Invalided|Slashed|NotEnough) == (Invalided | Slashed | NotEnough)
}

func IsCan_Slashed_NotEnough(status uint32) bool {
	return status&(Slashed|NotEnough) == (Slashed | NotEnough)
}

// The Candidate info
type Candidate struct {
	NodeId discover.NodeID
	// The account used to initiate the staking
	StakingAddress common.Address
	// The account recieve the block rewards and the staking rewards
	BenifitAddress common.Address
	// Block height at the time of staking
	StakingBlockNum uint64
	// The tx index at the time of staking
	StakingTxIndex uint32
	// The epoch number at staking or edit
	StakingEpoch uint64
	// All vons of staking and delegated
	Shares *big.Int
	// The staking vons  is circulating for effective epoch (in effect)
	Released *big.Int
	// The staking vons  is circulating for hesitant epoch (in hesitation)
	ReleasedTmp *big.Int
	// The staking vons  is locked for effective epoch (in effect)
	LockRepo *big.Int
	// The staking vons  is locked for hesitant epoch (in hesitation)
	LockRepoTmp *big.Int
	// Positive and negative signs:
	// Is it an increase or a decrease? 0: increase; 1: decrease
	Signs uint8

	// The candiate status
	Status uint32

	// Node desc
	Description
}

type Description struct {
	// The Candidate Node's Name  (with a length limit)
	NodeName string
	// External Id for the third party to pull the node description (with length limit)
	ExternalId string
	// The third-party home page of the node (with a length limit)
	Website string
	// Description of the node (with a length limit)
	Details string
}

// the Validator info
// They are Simplified Candidate
// They are consensus nodes and Epoch nodes snapshot
type Validator struct {
	NodeAddress common.Address
	NodeId      discover.NodeID
	// The weight
	// NOTE:
	// converted from the weight of Candidate is: (Int.Max - candidate.shares) + blocknum + txindex
	StakingWeight string
	// Validator's term in the consensus round
	ValidatorTerm uint32
}

// some consensus round validators or current epoch validators
type Validator_array struct {
	// the round start blocknumber or epoch start blocknumber
	Start uint64
	// the round end blocknumber or epoch blocknumber
	End uint64
	// the round validators or epoch validators
	Arr []*Validator
}

// the Delegate information
type Delegation struct {
	// The epoch number at delegate or edit
	DelegateEpoch uint64

	// Positive and negative signs:
	// Is it an increase or a decrease? 0: increase; 1: decrease
	Signs uint8

	// The delegate vons  is circulating for effective epoch (in effect)
	Released *big.Int
	// The delegate vons  is circulating for hesitant epoch (in hesitation)
	ReleasedTmp *big.Int
	// The delegate vons  is locked for effective epoch (in effect)
	LockRepo *big.Int
	// The delegate vons  is locked for hesitant epoch (in hesitation)
	LockRepoTmp *big.Int
}
