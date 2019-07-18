package xcom

import (
	"math/big"
	"sync"
)

// plugin rule key
const (
	DefualtRule = iota
	StakingRule
	SlashingRule
	RestrictingRule
	RewardRule
	GovernanceRule
)

type commonConfig struct {
	ExpectedMinutes uint64 // expected minutes every epoch
	Interval        uint64 // each block interval (uint: seconds)
	PerRoundBlocks  uint64 // blocks each validator will create per consensus epoch
	ValidatorCount  uint64 // The consensus validators count
}

type stakingConfig struct {
	StakeThreshold               *big.Int // The Staking minimum threshold allowed
	DelegateThreshold            *big.Int // The delegate minimum threshold allowed
	EpochValidatorNum            uint64   // The epoch (billing cycle) validators count
	ShiftValidatorNum            uint64   // The number of elections and replacements for each of the consensus rounds
	HesitateRatio                uint64   // Each hesitation period is a multiple of the epoch
	EffectiveRatio               uint64   // Each effective period is a multiple of the epoch
	ElectionDistance             uint64   // The interval of the last block of the high-distance consensus round of the election block for each consensus round
	UnStakeFreezeRatio           uint64   // The freeze period of the withdrew Staking (unit is  epochs)
	PassiveUnDelegateFreezeRatio uint64   // The freeze period of the delegate was invalidated due to the withdrawal of the Stake (unit is  epochs)
	ActiveUnDelegateFreezeRatio  uint64   // The freeze period of the delegate was invalidated due to active withdrew delegate (unit is  epochs)
}

type slashingConfig struct {
	BlockAmountLow            uint32 // The number of low exceptions per consensus round
	BlockAmountHigh           uint32 // Number of blocks per high consensus exception
	BlockAmountLowSlashing    uint32 // Penalty quota for each consensus round with a low number of abnormal blocks, percentage
	BlockAmountHighSlashing   uint32 // The penalty amount for each consensus round high abnormal number of blocks, percentage
	DuplicateSignNum          uint32 // The conditions for the highest penalty, double signing
	DuplicateSignLowSlashing  uint32 // Double sign low penalty amount, percentage
	DuplicateSignHighSlashing uint32 // DuplicateSignHighSlashing
}

type rewardConfig struct {
	// initial issuance:
	// 2% used for Reward
	// 0.5% used for developer foundation
	// 4.5% used for allowance
	// 2.5% almost used for Staking
	GenesisIssuance *big.Int // first year increase issuance at genesis block
}

type governanceConfig struct {
	SupportRateThreshold float64
}

// total
type EconomicModel struct {
	Common   commonConfig
	Staking  stakingConfig
	Slashing slashingConfig
	Reward   rewardConfig
	Gov      governanceConfig
}

var (
	modelOnce sync.Once
	ec        *EconomicModel
)

// Getting the global EconomicModel single instance
func GetEc(netId int8) *EconomicModel {
	modelOnce.Do(func() {
		ec = getDefaultEMConfig(netId)
	})
	return ec
}

func SetEconomicModel(ecParams *EconomicModel) {
	ec = ecParams
}

const (
	DefaultMainNet      = iota // PlatON default main net flag
	DefaultAlphaTestNet        // PlatON default Alpha test net flag
	DefaultBetaTestNet         // PlatON default Beta test net flag
	DefaultInnerTestNet        // PlatON default inner test net flag
	DefaultInnerDevNet         // PlatON default inner development net flag
	DefaultDeveloperNet        // PlatON default developer net flag
)

func getDefaultEMConfig(netId int8) *EconomicModel {
	var (
		success                bool
		stakeThresholdCount    string
		delegateThresholdCount string
		genesisIssuanceCount   string
		stakeThreshold         *big.Int
		delegateThreshold      *big.Int
		genesisIssuance        *big.Int
	)

	switch netId {
	case DefaultMainNet, DefaultDeveloperNet:
		stakeThresholdCount = "1000000000000000000000000"     // 100W von
		delegateThresholdCount = "10000000000000000000"       // 10 von
		genesisIssuanceCount = "1000000000000000000000000000" // 1,000,000,000 von
	case DefaultAlphaTestNet:
		stakeThresholdCount = "1000000000000000000000000"
		delegateThresholdCount = "10000000000000000000"
		genesisIssuanceCount = "1000000000000000000000000000"
	case DefaultBetaTestNet:
		stakeThresholdCount = "1000000000000000000000000"
		delegateThresholdCount = "10000000000000000000"
		genesisIssuanceCount = "1000000000000000000000000000"
	case DefaultInnerTestNet:
		stakeThresholdCount = "1000000000000000000000000"
		delegateThresholdCount = "10000000000000000000"
		genesisIssuanceCount = "1000000000000000000000000000"
	case DefaultInnerDevNet:
		stakeThresholdCount = "1000000000000000000000000"
		delegateThresholdCount = "10000000000000000000"
		genesisIssuanceCount = "1000000000000000000000000000"
	}

	if stakeThreshold, success = new(big.Int).SetString(stakeThresholdCount, 10); !success {
		return nil
	}
	if delegateThreshold, success = new(big.Int).SetString(delegateThresholdCount, 10); !success {
		return nil
	}
	if genesisIssuance, success = new(big.Int).SetString(genesisIssuanceCount, 10); !success {
		return nil
	}

	switch netId {
	case DefaultMainNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes: uint64(360), // 6 hours
				Interval:        uint64(1),
				PerRoundBlocks:  uint64(10),
				ValidatorCount:  uint64(25),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				DelegateThreshold:            delegateThreshold,
				EpochValidatorNum:            uint64(101),
				ShiftValidatorNum:            uint64(8),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				ElectionDistance:             uint64(20),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				BlockAmountLow:            uint32(8),
				BlockAmountHigh:           uint32(5),
				BlockAmountLowSlashing:    uint32(10),
				BlockAmountHighSlashing:   uint32(20),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(10),
			},
			Reward: rewardConfig{
				GenesisIssuance: genesisIssuance,
			},
			Gov: governanceConfig{
				SupportRateThreshold: float64(0.85),
			},
		}

	case DefaultAlphaTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes: uint64(10), // 10 minutes
				Interval:        uint64(1),
				PerRoundBlocks:  uint64(15),
				ValidatorCount:  uint64(4),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				DelegateThreshold:            delegateThreshold,
				EpochValidatorNum:            uint64(21),
				ShiftValidatorNum:            uint64(1),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				ElectionDistance:             uint64(10),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				BlockAmountLow:            uint32(8),
				BlockAmountHigh:           uint32(5),
				BlockAmountLowSlashing:    uint32(10),
				BlockAmountHighSlashing:   uint32(20),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(10),
			},
			Reward: rewardConfig{
				GenesisIssuance: genesisIssuance,
			},
			Gov: governanceConfig{
				SupportRateThreshold: float64(0.85),
			},
		}

	case DefaultBetaTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes: uint64(10), // 10 minutes
				Interval:        uint64(1),
				PerRoundBlocks:  uint64(15),
				ValidatorCount:  uint64(4),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				DelegateThreshold:            delegateThreshold,
				EpochValidatorNum:            uint64(21),
				ShiftValidatorNum:            uint64(1),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				ElectionDistance:             uint64(10),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				BlockAmountLow:            uint32(8),
				BlockAmountHigh:           uint32(5),
				BlockAmountLowSlashing:    uint32(10),
				BlockAmountHighSlashing:   uint32(20),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(10),
			},
			Reward: rewardConfig{
				GenesisIssuance: genesisIssuance,
			},
			Gov: governanceConfig{
				SupportRateThreshold: float64(0.85),
			},
		}

	case DefaultInnerTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes: uint64(666), // 11 hours
				Interval:        uint64(1),
				PerRoundBlocks:  uint64(25),
				ValidatorCount:  uint64(10),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				DelegateThreshold:            delegateThreshold,
				EpochValidatorNum:            uint64(51),
				ShiftValidatorNum:            uint64(3),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				ElectionDistance:             uint64(20),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				BlockAmountLow:            uint32(8),
				BlockAmountHigh:           uint32(5),
				BlockAmountLowSlashing:    uint32(10),
				BlockAmountHighSlashing:   uint32(20),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(10),
			},
			Reward: rewardConfig{
				GenesisIssuance: genesisIssuance,
			},
			Gov: governanceConfig{
				SupportRateThreshold: float64(0.85),
			},
		}

	case DefaultInnerDevNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes: uint64(10), // 10 minutes
				Interval:        uint64(1),
				PerRoundBlocks:  uint64(15),
				ValidatorCount:  uint64(4),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				DelegateThreshold:            delegateThreshold,
				EpochValidatorNum:            uint64(21),
				ShiftValidatorNum:            uint64(1),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				ElectionDistance:             uint64(10),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				BlockAmountLow:            uint32(8),
				BlockAmountHigh:           uint32(5),
				BlockAmountLowSlashing:    uint32(10),
				BlockAmountHighSlashing:   uint32(20),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(10),
			},
			Reward: rewardConfig{
				GenesisIssuance: genesisIssuance,
			},
			Gov: governanceConfig{
				SupportRateThreshold: float64(0.85),
			},
		}

	default:
		// Default is inner develop net config
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes: uint64(10), // 10 minutes
				Interval:        uint64(1),
				PerRoundBlocks:  uint64(15),
				ValidatorCount:  uint64(4),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				DelegateThreshold:            delegateThreshold,
				EpochValidatorNum:            uint64(21),
				ShiftValidatorNum:            uint64(1),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				ElectionDistance:             uint64(10),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				BlockAmountLow:            uint32(8),
				BlockAmountHigh:           uint32(5),
				BlockAmountLowSlashing:    uint32(10),
				BlockAmountHighSlashing:   uint32(20),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(10),
			},
			Reward: rewardConfig{
				GenesisIssuance: genesisIssuance,
			},
			Gov: governanceConfig{
				SupportRateThreshold: float64(0.85),
			},
		}
	}

	return ec
}

/******
 * Common configure
 ******/
func ExpectedMinutes() uint64 {
	return ec.Common.ExpectedMinutes
}
func Interval() uint64 {
	return ec.Common.Interval
}
func BlocksWillCreate() uint64 {
	return ec.Common.PerRoundBlocks
}
func ConsValidatorNum() uint64 {
	return ec.Common.ValidatorCount
}

/******
 * Staking configure
 ******/
func StakeThreshold() *big.Int {
	return ec.Staking.StakeThreshold
}

func DelegateThreshold() *big.Int {
	return ec.Staking.DelegateThreshold
}

func EpochValidatorNum() uint64 {
	return ec.Staking.EpochValidatorNum
}

func ShiftValidatorNum() uint64 {
	return ec.Staking.ShiftValidatorNum
}

func HesitateRatio() uint64 {
	return ec.Staking.HesitateRatio
}

func EffectiveRatio() uint64 {
	return ec.Staking.EffectiveRatio
}

func ElectionDistance() uint64 {
	return ec.Staking.ElectionDistance
}

func UnStakeFreezeRatio() uint64 {
	return ec.Staking.UnStakeFreezeRatio
}

func PassiveUnDelFreezeRatio() uint64 {
	return ec.Staking.PassiveUnDelegateFreezeRatio
}

func ActiveUnDelFreezeRatio() uint64 {
	return ec.Staking.ActiveUnDelegateFreezeRatio
}

/******
 * Slashing config
 ******/
func BlockAmountLow() uint32 {
	return ec.Slashing.BlockAmountLow
}

func BlockAmountHigh() uint32 {
	return ec.Slashing.BlockAmountHigh
}

func BlockAmountLowSlash() uint32 {
	return ec.Slashing.BlockAmountLowSlashing
}

func BlockAmountHighSlash() uint32 {
	return ec.Slashing.BlockAmountHighSlashing
}

func DuplicateSignNum() uint32 {
	return ec.Slashing.DuplicateSignNum
}

func DuplicateSignLowSlash() uint32 {
	return ec.Slashing.DuplicateSignLowSlashing
}

func DuplicateSignHighSlash() uint32 {
	return ec.Slashing.DuplicateSignHighSlashing
}

/******
 * Reward config
 ******/
func GenesisIssuance() *big.Int {
	return ec.Reward.GenesisIssuance
}

/******
 * Governance config
 ******/
func SupportRateThreshold() float64 {
	return ec.Gov.SupportRateThreshold
}
