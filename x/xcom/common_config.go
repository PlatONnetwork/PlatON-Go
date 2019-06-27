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
)
