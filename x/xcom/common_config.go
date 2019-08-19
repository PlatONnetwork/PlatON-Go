package xcom

import (
	"encoding/json"
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

var SecondsPerYear = uint64(365 * 24 * 3600)

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
	ExpectedMinutes       uint64 // expected minutes every epoch
	NodeBlockTimeWindow   uint64 // Node block time window (uint: seconds)
	PerRoundBlocks        uint64 // blocks each validator will create per consensus epoch
	OptValidatorCount     uint64 // Alternative number of validation nodes
	ValidatorCount        uint64 // The consensus validators count
	AdditionalCycleTime   uint64 // Additional cycle time (uint: minutes)
	PleRetLockCycle       uint64 // Pledge return lock cycle
	TextVotCycle          uint64 // Voting cycles for text proposals
	UpgradeMaxVotCycle    uint64 // The upgrade proposal has a maximum voting period
	UpgradeEffectiveCycle uint64 // Upgrade proposal effective period
}

type stakingConfig struct {
	StakeThreshold               *big.Int // The Staking minimum threshold allowed
	MinimumThreshold             *big.Int // The (incr, decr) delegate or incr staking minimum threshold allowed
	EpochValidatorNum            uint64   // The epoch (billing cycle) validators count
	HesitateRatio                uint64   // Each hesitation period is a multiple of the epoch
	EffectiveRatio               uint64   // Each effective period is a multiple of the epoch
	UnStakeFreezeRatio           uint64   // The freeze period of the withdrew Staking (unit is  epochs)
	PassiveUnDelegateFreezeRatio uint64   // The freeze period of the delegate was invalidated due to the withdrawal of the Stake (unit is  epochs)
	ActiveUnDelegateFreezeRatio  uint64   // The freeze period of the delegate was invalidated due to active withdrew delegate (unit is  epochs)
}

type slashingConfig struct {
	PackAmountAbnormal        uint32 // The number of blocks packed per round, reaching this value is abnormal
	PackAmountHighAbnormal    uint32 // The number of blocks packed per round, reaching this value is a high degree of abnormality
	PackAmountLowSlashRate    uint32 // Proportion of deducted quality deposit (when the number of packing blocks is abnormal); 10% -> 10
	PackAmountHighSlashRate   uint32 // Proportion of quality deposits deducted (when the number of packing blocks is high degree of abnormality); 20% -> 20
	DuplicateSignNum          uint32 // Number of multiple signatures
	DuplicateSignLowSlashing  uint32 // Deduction ratio when the number of multi-signs is lower than DuplicateSignNum; 10% -> 10
	DuplicateSignHighSlashing uint32 // Deduction ratio when the number of multi-signs is higher than DuplicateSignNum; 20% -> 20
}

type governanceConfig struct {
	VersionProposalVote_ConsensusRounds   uint64  // max Consensus-Round counts for version proposal's vote duration.
	VersionProposalActive_ConsensusRounds uint64  // default Consensus-Round counts for version proposal's active duration.
	VersionProposal_SupportRate           float64 // the version proposal will pass if the support rate exceeds this value.
	TextProposalVote_ConsensusRounds      uint64  // default Consensus-Round counts for text proposal's vote duration.
	TextProposal_VoteRate                 float64 // the text proposal will pass if the vote rate exceeds this value.
	TextProposal_SupportRate              float64 // the text proposal will pass if the vote support reaches this value.
	CancelProposal_VoteRate               float64 // the cancel proposal will pass if the vote rate exceeds this value.
	CancelProposal_SupportRate            float64 // the cancel proposal will pass if the vote support reaches this value.
}

// total
type EconomicModel struct {
	Common   commonConfig
	Staking  stakingConfig
	Slashing slashingConfig
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
		success               bool
		stakeThresholdCount   string
		minimumThresholdCount string
		stakeThreshold        *big.Int
		minimumThreshold      *big.Int
	)

	switch netId {
	case DefaultMainNet:
		stakeThresholdCount = "10000000000000000000000000" // 1000W von
		minimumThresholdCount = "10000000000000000000"     // 10 von
	case DefaultAlphaTestNet:
		stakeThresholdCount = "10000000000000000000000000"
		minimumThresholdCount = "10000000000000000000"
	case DefaultBetaTestNet:
		stakeThresholdCount = "10000000000000000000000000"
		minimumThresholdCount = "10000000000000000000"
	case DefaultInnerTestNet:
		stakeThresholdCount = "10000000000000000000000000"
		minimumThresholdCount = "10000000000000000000"
	case DefaultInnerDevNet:
		stakeThresholdCount = "10000000000000000000000000"
		minimumThresholdCount = "10000000000000000000"
	default: // DefaultDeveloperNet
		stakeThresholdCount = "10000000000000000000000000"
		minimumThresholdCount = "10000000000000000000"
	}

	if stakeThreshold, success = new(big.Int).SetString(stakeThresholdCount, 10); !success {
		return nil
	}
	if minimumThreshold, success = new(big.Int).SetString(minimumThresholdCount, 10); !success {
		return nil
	}

	switch netId {
	case DefaultMainNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes:       uint64(360), // 6 hours
				NodeBlockTimeWindow:   uint64(20),  // 20 seconds
				PerRoundBlocks:        uint64(10),
				OptValidatorCount:     uint64(101),
				ValidatorCount:        uint64(25),
				AdditionalCycleTime:   uint64(525600),
				PleRetLockCycle:       uint64(28),
				TextVotCycle:          uint64(2419),
				UpgradeMaxVotCycle:    uint64(2419),
				UpgradeEffectiveCycle: uint64(5),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				MinimumThreshold:             minimumThreshold,
				EpochValidatorNum:            uint64(101),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				UnStakeFreezeRatio:           uint64(28), // freezing 28 epoch
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				PackAmountAbnormal:        uint32(6),
				PackAmountHighAbnormal:    uint32(2),
				PackAmountLowSlashRate:    uint32(10),
				PackAmountHighSlashRate:   uint32(50),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(100),
			},
			Gov: governanceConfig{
				VersionProposalVote_ConsensusRounds:   uint64(2419),
				VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:           float64(0.667),
				TextProposalVote_ConsensusRounds:      uint64(2419),
				TextProposal_VoteRate:                 float64(0.50),
				TextProposal_SupportRate:              float64(0.667),
				CancelProposal_VoteRate:               float64(0.50),
				CancelProposal_SupportRate:            float64(0.667),
			},
		}

	case DefaultAlphaTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes:       uint64(10), // 10 minutes
				NodeBlockTimeWindow:   uint64(30), // 30 seconds
				PerRoundBlocks:        uint64(15),
				OptValidatorCount:     uint64(24),
				ValidatorCount:        uint64(4),
				AdditionalCycleTime:   uint64(525600),
				PleRetLockCycle:       uint64(28),
				TextVotCycle:          uint64(2419),
				UpgradeMaxVotCycle:    uint64(2419),
				UpgradeEffectiveCycle: uint64(5),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				MinimumThreshold:             minimumThreshold,
				EpochValidatorNum:            uint64(21),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				PackAmountAbnormal:        uint32(6),
				PackAmountHighAbnormal:    uint32(2),
				PackAmountLowSlashRate:    uint32(10),
				PackAmountHighSlashRate:   uint32(50),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(100),
			},
			Gov: governanceConfig{
				VersionProposalVote_ConsensusRounds:   uint64(2419),
				VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:           float64(0.667),
				TextProposalVote_ConsensusRounds:      uint64(2419),
				TextProposal_VoteRate:                 float64(0.50),
				TextProposal_SupportRate:              float64(0.667),
				CancelProposal_VoteRate:               float64(0.50),
				CancelProposal_SupportRate:            float64(0.667),
			},
		}

	case DefaultBetaTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes:       uint64(10), // 10 minutes
				NodeBlockTimeWindow:   uint64(30), // 30 seconds
				PerRoundBlocks:        uint64(15),
				OptValidatorCount:     uint64(24),
				ValidatorCount:        uint64(4),
				AdditionalCycleTime:   uint64(525600),
				PleRetLockCycle:       uint64(28),
				TextVotCycle:          uint64(2419),
				UpgradeMaxVotCycle:    uint64(2419),
				UpgradeEffectiveCycle: uint64(5),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				MinimumThreshold:             minimumThreshold,
				EpochValidatorNum:            uint64(21),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				PackAmountAbnormal:        uint32(6),
				PackAmountHighAbnormal:    uint32(2),
				PackAmountLowSlashRate:    uint32(10),
				PackAmountHighSlashRate:   uint32(50),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(100),
			},
			Gov: governanceConfig{
				VersionProposalVote_ConsensusRounds:   uint64(2419),
				VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:           float64(0.667),
				TextProposalVote_ConsensusRounds:      uint64(2419),
				TextProposal_VoteRate:                 float64(0.50),
				TextProposal_SupportRate:              float64(0.667),
				CancelProposal_VoteRate:               float64(0.50),
				CancelProposal_SupportRate:            float64(0.667),
			},
		}

	case DefaultInnerTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes:       uint64(666), // 11 hours
				NodeBlockTimeWindow:   uint64(50),  // 50 seconds
				PerRoundBlocks:        uint64(25),
				OptValidatorCount:     uint64(24),
				ValidatorCount:        uint64(10),
				AdditionalCycleTime:   uint64(525600),
				PleRetLockCycle:       uint64(28),
				TextVotCycle:          uint64(2419),
				UpgradeMaxVotCycle:    uint64(2419),
				UpgradeEffectiveCycle: uint64(5),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				MinimumThreshold:             minimumThreshold,
				EpochValidatorNum:            uint64(51),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				PackAmountAbnormal:        uint32(6),
				PackAmountHighAbnormal:    uint32(2),
				PackAmountLowSlashRate:    uint32(10),
				PackAmountHighSlashRate:   uint32(50),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(100),
			},
			Gov: governanceConfig{
				VersionProposalVote_ConsensusRounds:   uint64(2419),
				VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:           float64(0.667),
				TextProposalVote_ConsensusRounds:      uint64(2419),
				TextProposal_VoteRate:                 float64(0.50),
				TextProposal_SupportRate:              float64(0.667),
				CancelProposal_VoteRate:               float64(0.50),
				CancelProposal_SupportRate:            float64(0.667),
			},
		}

	case DefaultInnerDevNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes:       uint64(10), // 10 minutes
				NodeBlockTimeWindow:   uint64(30), // 30 seconds
				PerRoundBlocks:        uint64(15),
				OptValidatorCount:     uint64(24),
				ValidatorCount:        uint64(4),
				AdditionalCycleTime:   uint64(525600),
				PleRetLockCycle:       uint64(28),
				TextVotCycle:          uint64(2419),
				UpgradeMaxVotCycle:    uint64(2419),
				UpgradeEffectiveCycle: uint64(5),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				MinimumThreshold:             minimumThreshold,
				EpochValidatorNum:            uint64(21),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				PackAmountAbnormal:        uint32(6),
				PackAmountHighAbnormal:    uint32(2),
				PackAmountLowSlashRate:    uint32(10),
				PackAmountHighSlashRate:   uint32(50),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(100),
			},
			Gov: governanceConfig{
				VersionProposalVote_ConsensusRounds:   uint64(2419),
				VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:           float64(0.667),
				TextProposalVote_ConsensusRounds:      uint64(2419),
				TextProposal_VoteRate:                 float64(0.50),
				TextProposal_SupportRate:              float64(0.667),
				CancelProposal_VoteRate:               float64(0.50),
				CancelProposal_SupportRate:            float64(0.667),
			},
		}

	default: // DefaultDeveloperNet
		// Default is inner develop net config
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes:       uint64(10), // 10 minutes
				NodeBlockTimeWindow:   uint64(30), // 30 seconds
				PerRoundBlocks:        uint64(15),
				OptValidatorCount:     uint64(24),
				ValidatorCount:        uint64(4),
				AdditionalCycleTime:   uint64(525600),
				PleRetLockCycle:       uint64(28),
				TextVotCycle:          uint64(2419),
				UpgradeMaxVotCycle:    uint64(2419),
				UpgradeEffectiveCycle: uint64(5),
			},
			Staking: stakingConfig{
				StakeThreshold:               stakeThreshold,
				MinimumThreshold:             minimumThreshold,
				EpochValidatorNum:            uint64(21),
				HesitateRatio:                uint64(1),
				EffectiveRatio:               uint64(1),
				UnStakeFreezeRatio:           uint64(1),
				PassiveUnDelegateFreezeRatio: uint64(0),
				ActiveUnDelegateFreezeRatio:  uint64(0),
			},
			Slashing: slashingConfig{
				PackAmountAbnormal:        uint32(6),
				PackAmountHighAbnormal:    uint32(2),
				PackAmountLowSlashRate:    uint32(10),
				PackAmountHighSlashRate:   uint32(50),
				DuplicateSignNum:          uint32(2),
				DuplicateSignLowSlashing:  uint32(10),
				DuplicateSignHighSlashing: uint32(100),
			},
			Gov: governanceConfig{
				VersionProposalVote_ConsensusRounds:   uint64(2419),
				VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:           float64(0.667),
				TextProposalVote_ConsensusRounds:      uint64(2419),
				TextProposal_VoteRate:                 float64(0.50),
				TextProposal_SupportRate:              float64(0.667),
				CancelProposal_VoteRate:               float64(0.50),
				CancelProposal_SupportRate:            float64(0.667),
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
	return ec.Common.NodeBlockTimeWindow / ec.Common.PerRoundBlocks
}
func BlocksWillCreate() uint64 {
	return ec.Common.PerRoundBlocks
}
func ConsValidatorNum() uint64 {
	return ec.Common.ValidatorCount
}

func AdditionalCycleTime() uint64 {
	return ec.Common.AdditionalCycleTime
}

/******
 * Staking configure
 ******/
func StakeThreshold() *big.Int {
	return ec.Staking.StakeThreshold
}

func MinimumThreshold() *big.Int {
	return ec.Staking.MinimumThreshold
}

func EpochValidatorNum() uint64 {
	return ec.Staking.EpochValidatorNum
}

func ShiftValidatorNum() uint64 {
	return (ec.Common.ValidatorCount - 1) / 3
}

func HesitateRatio() uint64 {
	return ec.Staking.HesitateRatio
}

func EffectiveRatio() uint64 {
	return ec.Staking.EffectiveRatio
}

func ElectionDistance() uint64 {
	return 20
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
func PackAmountAbnormal() uint32 {
	return ec.Slashing.PackAmountAbnormal
}

func PackAmountHighAbnormal() uint32 {
	return ec.Slashing.PackAmountHighAbnormal
}

func PackAmountLowSlashRate() uint32 {
	return ec.Slashing.PackAmountLowSlashRate
}

func PackAmountHighSlashRate() uint32 {
	return ec.Slashing.PackAmountHighSlashRate
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

/******
 * Governance config
 ******/
func VersionProposalVote_ConsensusRounds() uint64 {
	return ec.Gov.VersionProposalVote_ConsensusRounds
}

func VersionProposalActive_ConsensusRounds() uint64 {
	return ec.Gov.VersionProposalActive_ConsensusRounds
}

func VersionProposal_SupportRate() float64 {
	return ec.Gov.VersionProposal_SupportRate
}

func TextProposalVote_ConsensusRounds() uint64 {
	return ec.Gov.TextProposalVote_ConsensusRounds
}

func TextProposal_VoteRate() float64 {
	return ec.Gov.TextProposal_VoteRate
}

func TextProposal_SupportRate() float64 {
	return ec.Gov.TextProposal_SupportRate
}

func CancelProposal_VoteRate() float64 {
	return ec.Gov.CancelProposal_VoteRate
}

func CancelProposal_SupportRate() float64 {
	return ec.Gov.CancelProposal_SupportRate
}

func PrintEc(blockNUmber *big.Int, blockHash common.Hash) {
	ecByte, _ := json.Marshal(ec)
	log.Debug("Current EconomicModel config", "blockNumber", blockNUmber, "blockHash", blockHash.Hex(), "ec", string(ecByte))
	//fmt.Println("Current EconomicModel config", "blockNumber", blockNUmber, "blockHash", blockHash.Hex(), "ec", string(ecByte))
}
