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

const (
	Zero                      = 0
	Eighty                    = 80
	Hundred                   = 100
	TenThousand               = 10000
	CeilBlocksReward          = 60101
	CeilMaxValidators         = 201
	CeilMaxConsensusVals      = 25
	PositiveInfinity          = "+∞"
	CeilUnStakeFreezeDuration = 28 * 4
	CeilMaxEvidenceAge        = CeilUnStakeFreezeDuration - 1
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
	MaxEpochMinutes     uint64 // expected minutes every epoch
	NodeBlockTimeWindow uint64 // Node block time window (uint: seconds)
	PerRoundBlocks      uint64 // blocks each validator will create per consensus epoch
	MaxConsensusVals    uint64 // The consensus validators count
	AdditionalCycleTime uint64 // Additional cycle time (uint: minutes)
}

type stakingConfig struct {
	StakeThreshold        *big.Int // The Staking minimum threshold allowed
	OperatingThreshold    *big.Int // The (incr, decr) delegate or incr staking minimum threshold allowed
	MaxValidators         uint64   // The epoch (billing cycle) validators count
	HesitateRatio         uint64   // Each hesitation period is a multiple of the epoch
	UnStakeFreezeDuration uint64   // The freeze period of the withdrew Staking (unit is  epochs)
}

type slashingConfig struct {
	SlashFractionDuplicateSign uint32 // Proportion of fines when double signing occurs
	DuplicateSignReportReward  uint32 // The percentage of rewards for whistleblowers, calculated from the penalty
	MaxEvidenceAge             uint32 // Validity period of evidence (unit is  epochs)
	SlashBlocksReward          uint32 // the number of blockReward to slashing per round

}

type governanceConfig struct {
	VersionProposalVote_DurationSeconds uint64 // voting duration, it will count into Consensus-Round.
	//VersionProposalActive_ConsensusRounds uint64  // default Consensus-Round counts for version proposal's active duration.
	VersionProposal_SupportRate       float64 // the version proposal will pass if the support rate exceeds this value.
	TextProposalVote_DurationSeconds  uint64  // voting duration, it will count into Consensus-Round.
	TextProposal_VoteRate             float64 // the text proposal will pass if the vote rate exceeds this value.
	TextProposal_SupportRate          float64 // the text proposal will pass if the vote support reaches this value.
	CancelProposal_VoteRate           float64 // the cancel proposal will pass if the vote rate exceeds this value.
	CancelProposal_SupportRate        float64 // the cancel proposal will pass if the vote support reaches this value.
	ParamProposalVote_DurationSeconds uint64  // voting duration, it will count into Epoch Round.
	ParamProposal_VoteRate            float64 // the param proposal will pass if the vote rate exceeds this value.
	ParamProposal_SupportRate         float64 // the param proposal will pass if the vote support reaches this value.
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

func ResetEconomicDefaultConfig(newEc *EconomicModel) {
	ec = newEc
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
				MaxEpochMinutes:     uint64(360), // 6 hours
				NodeBlockTimeWindow: uint64(20),  // 20 seconds
				PerRoundBlocks:      uint64(10),
				MaxConsensusVals:    uint64(25),
				AdditionalCycleTime: uint64(525600),
			},
			Staking: stakingConfig{
				StakeThreshold:        new(big.Int).Set(MillionLAT),
				OperatingThreshold:    new(big.Int).Set(TenLAT),
				MaxValidators:         uint64(101),
				HesitateRatio:         uint64(1),
				UnStakeFreezeDuration: uint64(28), // freezing 28 epoch
			},
			Slashing: slashingConfig{
				SlashFractionDuplicateSign: uint32(10),
				DuplicateSignReportReward:  uint32(50),
				MaxEvidenceAge:             uint32(27),
				SlashBlocksReward:          uint32(0),
			},
			Gov: governanceConfig{
				VersionProposalVote_DurationSeconds: uint64(14 * 24 * 3600),
				//VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:       float64(0.667),
				TextProposalVote_DurationSeconds:  uint64(14 * 24 * 3600),
				TextProposal_VoteRate:             float64(0.50),
				TextProposal_SupportRate:          float64(0.667),
				CancelProposal_VoteRate:           float64(0.50),
				CancelProposal_SupportRate:        float64(0.667),
				ParamProposalVote_DurationSeconds: uint64(14 * 24 * 3600),
				ParamProposal_VoteRate:            float64(0.50),
				ParamProposal_SupportRate:         float64(0.667),
			},
			Reward: rewardConfig{
				NewBlockRate:         50,
				PlatONFoundationYear: 10,
			},
			InnerAcc: innerAccount{
				PlatONFundAccount: common.HexToAddress("0x72188da050f4B3dD9a991b209221DBFE0A0fdC42"),
				PlatONFundBalance: new(big.Int).SetInt64(0),
				CDFAccount:        common.HexToAddress("0x8BAb06a9706F7613188d4Fb6310b1E5117dfd914"),
				CDFBalance:        new(big.Int).Set(cdfundBalance),
			},
		}

	case DefaultTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				MaxEpochMinutes:     uint64(6),  // 6 minutes
				NodeBlockTimeWindow: uint64(10), // 10 seconds
				PerRoundBlocks:      uint64(10),
				MaxConsensusVals:    uint64(4),
				AdditionalCycleTime: uint64(28),
			},
			Staking: stakingConfig{
				StakeThreshold:        new(big.Int).Set(MillionLAT),
				OperatingThreshold:    new(big.Int).Set(TenLAT),
				MaxValidators:         uint64(25),
				HesitateRatio:         uint64(1),
				UnStakeFreezeDuration: uint64(2),
			},
			Slashing: slashingConfig{
				SlashFractionDuplicateSign: uint32(10),
				DuplicateSignReportReward:  uint32(50),
				MaxEvidenceAge:             uint32(1),
				SlashBlocksReward:          uint32(0),
			},
			Gov: governanceConfig{
				VersionProposalVote_DurationSeconds: uint64(160),
				//VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:       float64(0.667),
				TextProposalVote_DurationSeconds:  uint64(160),
				TextProposal_VoteRate:             float64(0.50),
				TextProposal_SupportRate:          float64(0.667),
				CancelProposal_VoteRate:           float64(0.50),
				CancelProposal_SupportRate:        float64(0.667),
				ParamProposalVote_DurationSeconds: uint64(160),
				ParamProposal_VoteRate:            float64(0.50),
				ParamProposal_SupportRate:         float64(0.667),
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
				MaxEpochMinutes:     uint64(3),  // 3 minutes
				NodeBlockTimeWindow: uint64(10), // 10 seconds
				PerRoundBlocks:      uint64(10),
				MaxConsensusVals:    uint64(4),
				AdditionalCycleTime: uint64(28),
			},
			Staking: stakingConfig{
				StakeThreshold:        new(big.Int).Set(MillionLAT),
				OperatingThreshold:    new(big.Int).Set(TenLAT),
				MaxValidators:         uint64(25),
				HesitateRatio:         uint64(1),
				UnStakeFreezeDuration: uint64(2),
			},
			Slashing: slashingConfig{
				SlashFractionDuplicateSign: uint32(10),
				DuplicateSignReportReward:  uint32(50),
				MaxEvidenceAge:             uint32(1),
				SlashBlocksReward:          uint32(0),
			},
			Gov: governanceConfig{
				VersionProposalVote_DurationSeconds: uint64(160),
				//VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposal_SupportRate:       float64(0.667),
				TextProposalVote_DurationSeconds:  uint64(160),
				TextProposal_VoteRate:             float64(0.50),
				TextProposal_SupportRate:          float64(0.667),
				CancelProposal_VoteRate:           float64(0.50),
				CancelProposal_SupportRate:        float64(0.667),
				ParamProposalVote_DurationSeconds: uint64(160),
				ParamProposal_VoteRate:            float64(0.50),
				ParamProposal_SupportRate:         float64(0.667),
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

func CheckStakeThreshold(threshold *big.Int) error {

	if threshold.Cmp(MillionLAT) < 0 {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The StakeThreshold must be [%d, %s) LAT", MillionLAT, PositiveInfinity))
	}
	return nil
}

func CheckOperatingThreshold(threshold *big.Int) error {
	if threshold.Cmp(TenLAT) < 0 {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The OperatingThreshold must be [%d, %s) LAT", TenLAT, PositiveInfinity))
	}
	return nil
}

func CheckMaxValidators(num int) error {
	if num < CeilMaxConsensusVals || num > CeilMaxValidators {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The MaxValidators must be [%d, %d]", CeilMaxConsensusVals, CeilMaxValidators))
	}
	return nil
}

func CheckUnStakeFreezeDuration(duration, maxEvidenceAge int) error {
	if duration <= maxEvidenceAge || duration > CeilUnStakeFreezeDuration {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The UnStakeFreezeDuration must be (%d, %d]", maxEvidenceAge, CeilUnStakeFreezeDuration))
	}

	return nil
}

func CheckSlashFractionDuplicateSign(fraction int) error {
	if fraction <= Zero || fraction > TenThousand {
		return common.InvalidParameter.Wrap(fmt.Sprintf("SlashFractionDuplicateSign must be  (%d, %d]", Zero, TenThousand))
	}
	return nil
}

func CheckDuplicateSignReportReward(fraction int) error {
	if fraction <= Zero || fraction > Eighty {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The DuplicateSignReportReward must be (%d, %d]", Zero, Eighty))
	}
	return nil
}

func CheckMaxEvidenceAge(age, unStakeFreezeDuration int) error {
	if age <= Zero || age >= unStakeFreezeDuration {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The MaxEvidenceAge must be (%d, %d)", Zero, unStakeFreezeDuration))
	}
	return nil
}

func CheckSlashBlocksReward(rewards int) error {
	if rewards < Zero || rewards >= CeilBlocksReward {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The SlashBlocksReward must be [%d, %d)", Zero, CeilBlocksReward))
	}

	return nil
}

func CheckEconomicModel() error {
	if nil == ec {
		return errors.New("EconomicModel config is nil")
	}

	// epoch duration of config
	epochDuration := ec.Common.MaxEpochMinutes * 60
	// package perblock duration
	blockDuration := ec.Common.NodeBlockTimeWindow / ec.Common.PerRoundBlocks
	// round duration
	roundDuration := ec.Common.MaxConsensusVals * ec.Common.PerRoundBlocks * blockDuration
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

	if err := CheckMaxValidators(int(ec.Staking.MaxValidators)); nil != err {
		return err
	}

	if err := CheckOperatingThreshold(ec.Staking.OperatingThreshold); nil != err {
		return err
	}

	if ec.Staking.HesitateRatio < 1 {
		return errors.New("The HesitateRatio must be greater than or equal to 1")
	}

	if err := CheckStakeThreshold(ec.Staking.StakeThreshold); nil != err {
		return err
	}

	if err := CheckUnStakeFreezeDuration(int(ec.Staking.UnStakeFreezeDuration), int(ec.Slashing.MaxEvidenceAge)); nil != err {
		return err
	}

	if ec.Reward.PlatONFoundationYear < 1 {
		return errors.New("The PlatONFoundationYear must be greater than or equal to 1")
	}

	if ec.Reward.NewBlockRate < 0 || ec.Reward.NewBlockRate > 100 {
		return errors.New("The NewBlockRate must be greater than or equal to 0 and less than or equal to 100")
	}

	if err := CheckSlashFractionDuplicateSign(int(ec.Slashing.SlashFractionDuplicateSign)); nil != err {
		return err
	}

	if err := CheckDuplicateSignReportReward(int(ec.Slashing.DuplicateSignReportReward)); nil != err {
		return err
	}

	if err := CheckMaxEvidenceAge(int(ec.Slashing.MaxEvidenceAge), int(ec.Staking.UnStakeFreezeDuration)); nil != err {
		return err
	}

	if err := CheckSlashBlocksReward(int(ec.Slashing.SlashBlocksReward)); nil != err {
		return err
	}

	return nil
}

/******
 * Common configure
 ******/
func MaxEpochMinutes() uint64 {
	return ec.Common.MaxEpochMinutes
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
func MaxConsensusVals() uint64 {
	return ec.Common.MaxConsensusVals
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

func OperatingThreshold() *big.Int {
	return ec.Staking.OperatingThreshold
}

func MaxValidators() uint64 {
	return ec.Staking.MaxValidators
}

func ShiftValidatorNum() uint64 {
	return (ec.Common.MaxConsensusVals - 1) / 3
}

func HesitateRatio() uint64 {
	return ec.Staking.HesitateRatio
}

func ElectionDistance() uint64 {
	// min need two view
	return 2 * ec.Common.PerRoundBlocks
}

func UnStakeFreezeDuration() uint64 {
	return ec.Staking.UnStakeFreezeDuration
}

/******
 * Slashing config
 ******/
func SlashFractionDuplicateSign() uint32 {
	return ec.Slashing.SlashFractionDuplicateSign
}

func DuplicateSignReportReward() uint32 {
	return ec.Slashing.DuplicateSignReportReward
}

func MaxEvidenceAge() uint32 {
	return ec.Slashing.MaxEvidenceAge
}

func SlashBlocksReward() uint32 {
	return ec.Slashing.SlashBlocksReward
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
/*func VersionProposalVote_ConsensusRounds() uint64 {
	return ec.Gov.VersionProposalVote_DurationSeconds / (Interval() * ec.Common.PerRoundBlocks * ec.Common.MaxConsensusVals)
}*/

func VersionProposalVote_DurationSeconds() uint64 {
	return ec.Gov.VersionProposalVote_DurationSeconds
}

/*func VersionProposalActive_ConsensusRounds() uint64 {
	return ec.Gov.VersionProposalActive_ConsensusRounds
}*/

func VersionProposal_SupportRate() float64 {
	return ec.Gov.VersionProposal_SupportRate
}

/*func TextProposalVote_ConsensusRounds() uint64 {
	return ec.Gov.TextProposalVote_DurationSeconds / (Interval() * ec.Common.PerRoundBlocks * ec.Common.MaxConsensusVals)
}*/
func TextProposalVote_DurationSeconds() uint64 {
	return ec.Gov.TextProposalVote_DurationSeconds
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

func ParamProposalVote_DurationSeconds() uint64 {
	return ec.Gov.ParamProposalVote_DurationSeconds
}

func ParamProposal_VoteRate() float64 {
	return ec.Gov.ParamProposal_VoteRate
}

func ParamProposal_SupportRate() float64 {
	return ec.Gov.ParamProposal_SupportRate
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
