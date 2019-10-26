package xcom

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
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

var (

	// 10 LAT
	TenLAT, _ = new(big.Int).SetString("10000000000000000000", 10)

	// hard code genesis staking balance
	// 150W LAT
	GeneStakingAmount, _ = new(big.Int).SetString("1500000000000000000000000", 10)

	// 100W LAT
	MillionLAT, _ = new(big.Int).SetString("1000000000000000000000000", 10)
	// 1000W LAT
	TenMillionLAT, _ = new(big.Int).SetString("10000000000000000000000000", 10)
)

type commonConfig struct {
	ExpectedMinutes     uint64 // expected minutes every epoch
	NodeBlockTimeWindow uint64 // Node block time window (uint: seconds)
	PerRoundBlocks      uint64 // blocks each validator will create per consensus epoch
	ValidatorCount      uint64 // The consensus validators count
	AdditionalCycleTime uint64 // Additional cycle time (uint: minutes)
}

type stakingConfig struct {
	StakeThreshold              *big.Int // The Staking minimum threshold allowed
	MinimumThreshold            *big.Int // The (incr, decr) delegate or incr staking minimum threshold allowed
	EpochValidatorNum           uint64   // The epoch (billing cycle) validators count
	HesitateRatio               uint64   // Each hesitation period is a multiple of the epoch
	UnStakeFreezeRatio          uint64   // The freeze period of the withdrew Staking (unit is  epochs)
	ActiveUnDelegateFreezeRatio uint64   // The freeze period of the delegate was invalidated due to active withdrew delegate (unit is  epochs)
}

type slashingConfig struct {
	DuplicateSignHighSlashing      uint32 // Proportion of fines when double signing occurs
	DuplicateSignReportReward      uint32 // The percentage of rewards for whistleblowers, calculated from the penalty
	NumberOfBlockRewardForSlashing uint32 // the number of blockReward to slashing per round
	EvidenceValidEpoch             uint32 // Validity period of evidence (unit is  epochs)
}

type governanceConfig struct {
	VersionProposalVote_DurationSeconds   uint64  // max Consensus-Round counts for version proposal's vote duration.
	VersionProposalActive_ConsensusRounds uint64  // default Consensus-Round counts for version proposal's active duration.
	VersionProposal_SupportRate           float64 // the version proposal will pass if the support rate exceeds this value.
	TextProposalVote_DurationSeconds      uint64  // default Consensus-Round counts for text proposal's vote duration.
	TextProposal_VoteRate                 float64 // the text proposal will pass if the vote rate exceeds this value.
	TextProposal_SupportRate              float64 // the text proposal will pass if the vote support reaches this value.
	CancelProposal_VoteRate               float64 // the cancel proposal will pass if the vote rate exceeds this value.
	CancelProposal_SupportRate            float64 // the cancel proposal will pass if the vote support reaches this value.
}

type rewardConfig struct {
	NewBlockRate         uint64 // This is the package block reward AND staking reward  rate, eg: 20 ==> 20%, newblock: 20%, staking: 80%
	PlatONFoundationYear uint32 // Foundation allotment year, representing a percentage of the boundaries of the Foundation each year
}

type innerAccount struct {
	// Account of PlatONFoundation
	PlatONFundAccount common.Address
	PlatONFundBalance *big.Int
	// Account of CommunityDeveloperFoundation
	CDFAccount common.Address
	CDFBalance *big.Int
}

// total
type EconomicModel struct {
	Common   commonConfig
	Staking  stakingConfig
	Slashing slashingConfig
	Gov      governanceConfig
	Reward   rewardConfig
	InnerAcc innerAccount
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
	DefaultMainNet = iota // PlatON default main net flag
	DefaultTestNet        // PlatON default test net flag
)

func getDefaultEMConfig(netId int8) *EconomicModel {
	var (
		ok            bool
		cdfundBalance *big.Int
	)

	// 3.31811981  thousand millions LAT
	if cdfundBalance, ok = new(big.Int).SetString("331811981000000000000000000", 10); !ok {
		return nil
	}

	switch netId {
	case DefaultMainNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes:     uint64(360), // 6 hours
				NodeBlockTimeWindow: uint64(20),  // 20 seconds
				PerRoundBlocks:      uint64(10),
				ValidatorCount:      uint64(25),
				AdditionalCycleTime: uint64(525600),
			},
			Staking: stakingConfig{
				StakeThreshold:              new(big.Int).Set(MillionLAT),
				MinimumThreshold:            new(big.Int).Set(TenLAT),
				EpochValidatorNum:           uint64(101),
				HesitateRatio:               uint64(1),
				UnStakeFreezeRatio:          uint64(28), // freezing 28 epoch
				ActiveUnDelegateFreezeRatio: uint64(0),
			},
			Slashing: slashingConfig{
				DuplicateSignHighSlashing:      uint32(10),
				DuplicateSignReportReward:      uint32(50),
				NumberOfBlockRewardForSlashing: uint32(0),
				EvidenceValidEpoch:             uint32(27),
			},
			Gov: governanceConfig{
				VersionProposalVote_DurationSeconds:   uint64(14 * 24 * 3600),
				VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:           float64(0.667),
				TextProposalVote_DurationSeconds:      uint64(14 * 24 * 3600),
				TextProposal_VoteRate:                 float64(0.50),
				TextProposal_SupportRate:              float64(0.667),
				CancelProposal_VoteRate:               float64(0.50),
				CancelProposal_SupportRate:            float64(0.667),
			},
			Reward: rewardConfig{
				NewBlockRate:         50,
				PlatONFoundationYear: 10,
			},
			InnerAcc: innerAccount{
				PlatONFundAccount: common.HexToAddress("0x55bfd49472fd41211545b01713a9c3a97af78b05"),
				PlatONFundBalance: new(big.Int).SetInt64(0),
				CDFAccount:        common.HexToAddress("0x60ceca9c1290ee56b98d4e160ef0453f7c40d219"),
				CDFBalance:        new(big.Int).Set(cdfundBalance),
			},
		}

	case DefaultTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes:     uint64(6),  // 6 minutes
				NodeBlockTimeWindow: uint64(10), // 10 seconds
				PerRoundBlocks:      uint64(10),
				ValidatorCount:      uint64(4),
				AdditionalCycleTime: uint64(28),
			},
			Staking: stakingConfig{
				StakeThreshold:              new(big.Int).Set(MillionLAT),
				MinimumThreshold:            new(big.Int).Set(TenLAT),
				EpochValidatorNum:           uint64(24),
				HesitateRatio:               uint64(1),
				UnStakeFreezeRatio:          uint64(2),
				ActiveUnDelegateFreezeRatio: uint64(0),
			},
			Slashing: slashingConfig{
				DuplicateSignHighSlashing:      uint32(10),
				DuplicateSignReportReward:      uint32(50),
				NumberOfBlockRewardForSlashing: uint32(0),
				EvidenceValidEpoch:             uint32(1),
			},
			Gov: governanceConfig{
				VersionProposalVote_DurationSeconds:   uint64(160),
				VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:           float64(0.667),
				TextProposalVote_DurationSeconds:      uint64(160),
				TextProposal_VoteRate:                 float64(0.50),
				TextProposal_SupportRate:              float64(0.667),
				CancelProposal_VoteRate:               float64(0.50),
				CancelProposal_SupportRate:            float64(0.667),
			},
			Reward: rewardConfig{
				NewBlockRate:         50,
				PlatONFoundationYear: 10,
			},
			InnerAcc: innerAccount{
				PlatONFundAccount: common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185821"),
				PlatONFundBalance: new(big.Int).SetInt64(0),
				CDFAccount:        common.HexToAddress("0xc1f330b214668beac2e6418dd651b09c759a4bf5"),
				CDFBalance:        new(big.Int).Set(cdfundBalance),
			},
		}

	default: // DefaultTestNet
		// Default is test net config
		ec = &EconomicModel{
			Common: commonConfig{
				ExpectedMinutes:     uint64(3),  // 3 minutes
				NodeBlockTimeWindow: uint64(10), // 10 seconds
				PerRoundBlocks:      uint64(10),
				ValidatorCount:      uint64(4),
				AdditionalCycleTime: uint64(28),
			},
			Staking: stakingConfig{
				StakeThreshold:              new(big.Int).Set(MillionLAT),
				MinimumThreshold:            new(big.Int).Set(TenLAT),
				EpochValidatorNum:           uint64(24),
				HesitateRatio:               uint64(1),
				UnStakeFreezeRatio:          uint64(2),
				ActiveUnDelegateFreezeRatio: uint64(0),
			},
			Slashing: slashingConfig{
				DuplicateSignHighSlashing:      uint32(10),
				DuplicateSignReportReward:      uint32(50),
				NumberOfBlockRewardForSlashing: uint32(0),
				EvidenceValidEpoch:             uint32(1),
			},
			Gov: governanceConfig{
				VersionProposalVote_DurationSeconds:   uint64(160),
				VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:           float64(0.667),
				TextProposalVote_DurationSeconds:      uint64(160),
				TextProposal_VoteRate:                 float64(0.50),
				TextProposal_SupportRate:              float64(0.667),
				CancelProposal_VoteRate:               float64(0.50),
				CancelProposal_SupportRate:            float64(0.667),
			},
			Reward: rewardConfig{
				NewBlockRate:         50,
				PlatONFoundationYear: 1,
			},
			InnerAcc: innerAccount{
				PlatONFundAccount: common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185821"),
				PlatONFundBalance: new(big.Int).SetInt64(0),
				CDFAccount:        common.HexToAddress("0xc1f330b214668beac2e6418dd651b09c759a4bf5"),
				CDFBalance:        new(big.Int).Set(cdfundBalance),
			},
		}
	}

	return ec
}

func CheckEconomicModel() error {
	if nil == ec {
		return errors.New("EconomicModel config is nil")
	}

	// epoch duration of config
	epochDuration := ec.Common.ExpectedMinutes * 60
	// package perblock duration
	blockDuration := ec.Common.NodeBlockTimeWindow / ec.Common.PerRoundBlocks
	// round duration
	roundDuration := ec.Common.ValidatorCount * ec.Common.PerRoundBlocks * blockDuration
	// epoch Size, how many consensus round
	epochSize := epochDuration / roundDuration
	//real epoch duration
	realEpochDuration := epochSize * roundDuration

	log.Info("Call CheckEconomicModel: check epoch and consensus round", "config epoch duration", fmt.Sprintf("%d s", epochDuration),
		"perblock duration", fmt.Sprintf("%d s", blockDuration), "round duration", fmt.Sprintf("%d s", roundDuration),
		"real epoch duration", fmt.Sprintf("%d s", realEpochDuration), "consensus count of epoch", epochSize)

	if epochSize < 4 {
		return errors.New("The settlement period must be more than four times the consensus period")
	}

	// additionalCycle Size, how many epoch duration
	additionalCycleSize := ec.Common.AdditionalCycleTime * 60 / realEpochDuration
	// realAdditionalCycleDuration
	realAdditionalCycleDuration := additionalCycleSize * realEpochDuration / 60

	log.Info("Call CheckEconomicModel: additional cycle and epoch", "config additional cycle duration", fmt.Sprintf("%d min", ec.Common.AdditionalCycleTime),
		"real additional cycle duration", fmt.Sprintf("%d min", realAdditionalCycleDuration), "epoch count of additional cycle", additionalCycleSize)

	if additionalCycleSize < 4 {
		return errors.New("The issuance period must be integer multiples of the settlement period and multiples must be greater than or equal to 4")
	}
	if ec.Staking.EpochValidatorNum < ec.Common.ValidatorCount {
		return errors.New("The EpochValidatorNum must be greater than or equal to the ValidatorCount")
	}

	if ec.Staking.MinimumThreshold.Cmp(TenLAT) < 0 {
		return errors.New(fmt.Sprintf("The MinimumThreshold must be greater than or equal to %s von", TenLAT.String()))
	}

	if ec.Staking.StakeThreshold.Cmp(common.Big0) <= 0 {
		return errors.New(fmt.Sprintf("The StakeThreshold must be greater than %s von", common.Big0.String()))
	}

	// the StakeThreshold must be less than geneStakeAmount
	if ec.Staking.StakeThreshold.Cmp(GeneStakingAmount) > 0 {
		return errors.New(fmt.Sprintf("The StakeThreshold must be less than or equal to %s von", GeneStakingAmount.String()))
	}

	if ec.Staking.HesitateRatio < 1 {
		return errors.New("The HesitateRatio must be greater than or equal to 1")
	}

	if ec.Staking.UnStakeFreezeRatio < 1 {
		return errors.New("The UnStakeFreezeRatio must be greater than or equal to 1")
	}

	if ec.Reward.PlatONFoundationYear < 1 {
		return errors.New("The PlatONFoundationYear must be greater than or equal to 1")
	}

	if ec.Reward.NewBlockRate < 0 || ec.Reward.NewBlockRate > 100 {
		return errors.New("The NewBlockRate must be greater than or equal to 0 and less than or equal to 100")
	}

	if ec.Slashing.DuplicateSignHighSlashing < 0 || ec.Slashing.DuplicateSignHighSlashing > 10000 {
		return errors.New("DuplicateSignHighSlashing must be a floating point value between 0 and 10000")
	}

	if ec.Slashing.DuplicateSignReportReward < 0 || ec.Slashing.DuplicateSignReportReward > 100 {
		return errors.New("The DuplicateSignReportReward must be greater than or equal to 0 and less than or equal to 100")
	}

	if uint64(ec.Slashing.EvidenceValidEpoch) >= ec.Staking.UnStakeFreezeRatio {
		return errors.New("The EvidenceValidEpoch must be less than to the UnStakeFreezeRatio")
	}

	return nil
}

/******
 * Common configure
 ******/
func ExpectedMinutes() uint64 {
	return ec.Common.ExpectedMinutes
}

// set the value by genesis block
func SetNodeBlockTimeWindow(period uint64) {
	if ec != nil {
		ec.Common.NodeBlockTimeWindow = period
	}
}
func SetPerRoundBlocks(amount uint64) {
	if ec != nil {
		ec.Common.PerRoundBlocks = amount
	}
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

func ElectionDistance() uint64 {
	// min need two view
	return 2 * ec.Common.PerRoundBlocks
}

func UnStakeFreezeRatio() uint64 {
	return ec.Staking.UnStakeFreezeRatio
}

func ActiveUnDelFreezeRatio() uint64 {
	return ec.Staking.ActiveUnDelegateFreezeRatio
}

/******
 * Slashing config
 ******/
func DuplicateSignHighSlash() uint32 {
	return ec.Slashing.DuplicateSignHighSlashing
}

func DuplicateSignReportReward() uint32 {
	return ec.Slashing.DuplicateSignReportReward
}

func NumberOfBlockRewardForSlashing() uint32 {
	return ec.Slashing.NumberOfBlockRewardForSlashing
}

func EvidenceValidEpoch() uint32 {
	return ec.Slashing.EvidenceValidEpoch
}

/******
 * Reward config
 ******/
func NewBlockRewardRate() uint64 {
	return ec.Reward.NewBlockRate
}

func PlatONFoundationYear() uint32 {
	return ec.Reward.PlatONFoundationYear
}

/******
 * Governance config
 ******/
func VersionProposalVote_ConsensusRounds() uint64 {
	return ec.Gov.VersionProposalVote_DurationSeconds / (Interval() * ec.Common.PerRoundBlocks * ec.Common.ValidatorCount)
}

func VersionProposalActive_ConsensusRounds() uint64 {
	return ec.Gov.VersionProposalActive_ConsensusRounds
}

func VersionProposal_SupportRate() float64 {
	return ec.Gov.VersionProposal_SupportRate
}

func TextProposalVote_ConsensusRounds() uint64 {
	return ec.Gov.TextProposalVote_DurationSeconds / (Interval() * ec.Common.PerRoundBlocks * ec.Common.ValidatorCount)
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

/******
 * Inner Account Config
 ******/
func PlatONFundAccount() common.Address {
	return ec.InnerAcc.PlatONFundAccount
}

func PlatONFundBalance() *big.Int {
	return ec.InnerAcc.PlatONFundBalance
}

func CDFAccount() common.Address {
	return ec.InnerAcc.CDFAccount
}

func CDFBalance() *big.Int {
	return ec.InnerAcc.CDFBalance
}

func EconomicString() string {
	if nil != ec {
		ecByte, _ := json.Marshal(ec)
		return string(ecByte)
	} else {
		return ""
	}
}
