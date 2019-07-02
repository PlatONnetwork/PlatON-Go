package xcom

import "math/big"

// plugin rule key
const (
	DefualtRule = iota
	StakingRule
	SlashingRule
	RestrictingRule
	AwardmgrRule

	// ......
)

// config  TODO Configuration is all here
type EconomicModel struct {
	Staking  	    StakingConfig
	Slashing 	    SlashingConfig
	EpochsPerYear   uint32
}

var DefaultConfig = EconomicModel{
	Staking:  		DefaultStakingConfig,
	Slashing: 		DefaultSlashingConfig,
	EpochsPerYear:	1,
}

type StakingConfig struct {
	StakeThreshold	  				*big.Int
	DelegateThreshold 				*big.Int
	ConsValidatorNum  				uint64
	EpochValidatorNum 				uint64
	ShiftValidatorNum 				uint64
	EpochSize		  				uint64
	HesitateRatio	  				uint64
	EffectiveRatio	 				uint64
	ElectionDistance  				uint64
	ConsensusSize	  				uint64
	UnStakeFreezeRatio 				uint64
	PassiveUnDelegateFreezeRatio	uint64
	ActiveUnDelegateFreezeRatio		uint64
}

type SlashingConfig struct {
	BlockAmountLow				uint32
	BlockAmountHigh				uint32
	BlockAmountLowSlashing		uint32
	BlockAmountHighSlashing		uint32
	DuplicateSignNum			uint32
	DuplicateSignLowSlashing	uint32
	DuplicateSignHighSlashing	uint32
}

/**
Staking config
**/
var (
	// The staking minimum threshold allowed (100,0000 LAT)
	StakeThreshold, _ = new(big.Int).SetString("1000000000000000000000000", 10)
	// The delegate minimum threshold allowed (10 LAT)
	DelegateThreshold, _ = new(big.Int).SetString("10", 10)
	// The consensus validators count
	ConsValidatorNum = uint64(25)
	// The epoch (billing cycle) validators count
	EpochValidatorNum = uint64(101)
	// The number of elections and replacements for each of the consensus rounds
	ShiftValidatorNum = uint64(8)
	// Each epoch (billing cycle) is a multiple of the consensus rounds
	EpochSize = uint64(88)
	// Each hesitation period is a multiple of the epoch
	HesitateRatio = uint64(1)
	// Each effective period is a multiple of the epoch
	EffectiveRatio = uint64(1)
	// The interval of the last block of the high-distance
	// consensus round of the election block for each consensus round
	ElectionDistance = uint64(20)
	// Number of blocks per consensus round
	ConsensusSize = uint64(250)

	// The freeze period of the withdrew staking (unit is  epochs)
	UnStakeFreezeRatio = uint64(1)

	// The freeze period of the delegate was invalidated
	// due to the withdrawal of the Stake (unit is  epochs)
	PassiveUnDelegateFreezeRatio = uint64(0)

	// The freeze period of the delegate was invalidated
	// due to active withdrew delegate (unit is  epochs)
	ActiveUnDelegateFreezeRatio = uint64(0)
)

var DefaultStakingConfig = StakingConfig{
	StakeThreshold: 	StakeThreshold,
	DelegateThreshold:  DelegateThreshold,
	ConsValidatorNum: 	uint64(25),
	EpochValidatorNum:  uint64(101),
	ShiftValidatorNum:  uint64(8),
	EpochSize: 			uint64(88),
	HesitateRatio: 		uint64(1),
	EffectiveRatio: 	uint64(1),
	ElectionDistance: 	uint64(20),
	ConsensusSize: 		uint64(250),
	UnStakeFreezeRatio: uint64(1),
	PassiveUnDelegateFreezeRatio: uint64(0),
	ActiveUnDelegateFreezeRatio: uint64(0),
}

var DefaultSlashingConfig = SlashingConfig{
	BlockAmountLow:				0,
	BlockAmountHigh:			0,
	BlockAmountLowSlashing:		0,
	BlockAmountHighSlashing:	0,
	DuplicateSignNum:			0,
	DuplicateSignLowSlashing:	0,
	DuplicateSignHighSlashing:	0,
}