// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

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
	CollectDeclareVersionRule
)

const (
	Zero                      = 0
	Eighty                    = 80
	Hundred                   = 100
	TenThousand               = 10000
	CeilBlocksReward          = 50000
	CeilMaxValidators         = 10000
	FloorMaxConsensusVals     = 4
	CeilMaxConsensusVals      = 43
	PositiveInfinity          = "+∞"
	CeilUnStakeFreezeDuration = 168 * 2
	CeilMaxEvidenceAge        = CeilUnStakeFreezeDuration - 1
	// The maximum time range for the cumulative number of zero blocks (No more than 64)
	maxZeroProduceCumulativeTime uint16 = 50

	RewardPerMaxChangeRangeUpperLimit = 2000
	RewardPerMaxChangeRangeLowerLimit = 1
	RewardPerChangeIntervalUpperLimit = 28
	RewardPerChangeIntervalLowerLimit = 2
	IncreaseIssuanceRatioUpperLimit   = 2000
	IncreaseIssuanceRatioLowerLimit   = 0

	// When electing consensus nodes, it is used to calculate the P value of the binomial distribution
	ElectionBase = 30
)

var (
	one, _ = new(big.Int).SetString("1000000000000000000", 10)

	// 10 ATP
	DelegateLowerLimit, _ = new(big.Int).SetString("10000000000000000000", 10)

	// 1W ATP
	DelegateUpperLimit, _ = new(big.Int).SetString("10000000000000000000000", 10)

	// hard code genesis staking balance
	// 15W LAT
	GeneStakingAmount, _ = new(big.Int).SetString("150000000000000000000000", 10)

	// 10W
	StakeLowerLimit, _ = new(big.Int).SetString("100000000000000000000000", 10)
	// 1000W ATP
	StakeUpperLimit, _ = new(big.Int).SetString("10000000000000000000000000", 10)

	FloorMinimumRelease = new(big.Int).Mul(new(big.Int).SetUint64(100), one)
	CeilMinimumRelease  = new(big.Int).Mul(new(big.Int).SetUint64(10000000), one)
)

type commonConfig struct {
	MaxEpochMinutes     uint64 `json:"maxEpochMinutes"`     // expected minutes every epoch
	NodeBlockTimeWindow uint64 `json:"nodeBlockTimeWindow"` // Node block time window (uint: seconds)
	PerRoundBlocks      uint64 `json:"perRoundBlocks"`      // blocks each validator will create per consensus epoch
	MaxConsensusVals    uint64 `json:"maxConsensusVals"`    // The consensus validators count
	AdditionalCycleTime uint64 `json:"additionalCycleTime"` // Additional cycle time (uint: minutes)
}

type stakingConfig struct {
	StakeThreshold          *big.Int `json:"stakeThreshold"`          // The Staking minimum threshold allowed
	OperatingThreshold      *big.Int `json:"operatingThreshold"`      // The (incr, decr) delegate or incr staking minimum threshold allowed
	MaxValidators           uint64   `json:"maxValidators"`           // The epoch (billing cycle) validators count
	UnStakeFreezeDuration   uint64   `json:"unStakeFreezeDuration"`   // The freeze period of the withdrew Staking (unit is  epochs)
	RewardPerMaxChangeRange uint16   `json:"rewardPerMaxChangeRange"` // The maximum amount of commission reward ratio that can be modified each time
	RewardPerChangeInterval uint16   `json:"rewardPerChangeInterval"` // The interval for each modification of the commission reward ratio (unit: epoch)
}

type slashingConfig struct {
	SlashFractionDuplicateSign uint32 `json:"slashFractionDuplicateSign"` // Proportion of fines when double signing occurs
	DuplicateSignReportReward  uint32 `json:"duplicateSignReportReward"`  // The percentage of rewards for whistleblowers, calculated from the penalty
	MaxEvidenceAge             uint32 `json:"maxEvidenceAge"`             // Validity period of evidence (unit is  epochs)
	SlashBlocksReward          uint32 `json:"slashBlocksReward"`          // the number of blockReward to slashing per round
	ZeroProduceCumulativeTime  uint16 `json:"zeroProduceCumulativeTime"`  // Count the number of zero-production blocks in this time range and check it. If it reaches a certain number of times, it can be punished (unit is consensus round)
	ZeroProduceNumberThreshold uint16 `json:"zeroProduceNumberThreshold"` // Threshold for the number of zero production blocks. punishment is reached within the specified time range
	ZeroProduceFreezeDuration  uint64 `json:"zeroProduceFreezeDuration"`  // Number of settlement cycles frozen after zero block penalty (unit is epochs)
}

type governanceConfig struct {
	VersionProposalVoteDurationSeconds uint64 `json:"versionProposalVoteDurationSeconds"` // voting duration, it will count into Consensus-Round.
	VersionProposalSupportRate         uint64 `json:"versionProposalSupportRate"`         // the version proposal will pass if the support rate exceeds this value.
	TextProposalVoteDurationSeconds    uint64 `json:"textProposalVoteDurationSeconds"`    // voting duration, it will count into Consensus-Round.
	TextProposalVoteRate               uint64 `json:"textProposalVoteRate"`               // the text proposal will pass if the vote rate exceeds this value.
	TextProposalSupportRate            uint64 `json:"textProposalSupportRate"`            // the text proposal will pass if the vote support reaches this value.
	CancelProposalVoteRate             uint64 `json:"cancelProposalVoteRate"`             // the cancel proposal will pass if the vote rate exceeds this value.
	CancelProposalSupportRate          uint64 `json:"cancelProposalSupportRate"`          // the cancel proposal will pass if the vote support reaches this value.
	ParamProposalVoteDurationSeconds   uint64 `json:"paramProposalVoteDurationSeconds"`   // voting duration, it will count into Epoch Round.
	ParamProposalVoteRate              uint64 `json:"paramProposalVoteRate"`              // the param proposal will pass if the vote rate exceeds this value.
	ParamProposalSupportRate           uint64 `json:"paramProposalSupportRate"`           // the param proposal will pass if the vote support reaches this value.
}

type rewardConfig struct {
	NewBlockRate                 uint64 `json:"newBlockRate"`                 // This is the package block reward AND staking reward  rate, eg: 20 ==> 20%, newblock: 20%, staking: 80%
	PlatONFoundationYear         uint32 `json:"platonFoundationYear"`         // Foundation allotment year, representing a percentage of the boundaries of the Foundation each year
	IncreaseIssuanceRatio        uint16 `json:"increaseIssuanceRatio"`        // According to the total amount issued in the previous year, increase the proportion of issuance
	TheNumberOfDelegationsReward uint16 `json:"TheNumberOfDelegationsReward"` // The maximum number of delegates that can receive rewards at a time
}

type restrictingConfig struct {
	MinimumRelease *big.Int `json:"minimumRelease"` //The minimum number of Restricting release in one epoch
}

type innerAccount struct {
	// Account of PlatONFoundation
	PlatONFundAccount common.Address `json:"platonFundAccount"`
	PlatONFundBalance *big.Int       `json:"platonFundBalance"`
	// Account of CommunityDeveloperFoundation
	CDFAccount common.Address `json:"cdfAccount"`
	CDFBalance *big.Int       `json:"cdfBalance"`
}

// total
type EconomicModel struct {
	Common      commonConfig      `json:"common"`
	Staking     stakingConfig     `json:"staking"`
	Slashing    slashingConfig    `json:"slashing"`
	Gov         governanceConfig  `json:"gov"`
	Reward      rewardConfig      `json:"reward"`
	Restricting restrictingConfig `json:"restricting"`
	InnerAcc    innerAccount      `json:"innerAcc"`
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
	DefaultMainNet     = iota // PlatON default main net flag
	DefaultTestNet            // PlatON default test net flag
	DefaultUnitTestNet        // PlatON default unit test
)

func getDefaultEMConfig(netId int8) *EconomicModel {
	var (
		ok            bool
		cdfundBalance *big.Int
	)

	// 3.31811981  thousand millions LAT
	if cdfundBalance, ok = new(big.Int).SetString("322361981000000000000000000", 10); !ok {
		return nil
	}

	one, _ := new(big.Int).SetString("1000000000000000000", 10)

	switch netId {
	case DefaultMainNet:
		ec = &EconomicModel{
			Common: commonConfig{
				MaxEpochMinutes:     uint64(360), // 6 hours
				NodeBlockTimeWindow: uint64(20),  // 20 seconds
				PerRoundBlocks:      uint64(10),
				MaxConsensusVals:    uint64(43),
				AdditionalCycleTime: uint64(525960),
			},
			Staking: stakingConfig{
				StakeThreshold:          new(big.Int).Set(StakeLowerLimit),
				OperatingThreshold:      new(big.Int).Set(DelegateLowerLimit),
				MaxValidators:           uint64(201),
				UnStakeFreezeDuration:   uint64(168), // freezing 168 epoch
				RewardPerMaxChangeRange: uint16(500),
				RewardPerChangeInterval: uint16(10),
			},
			Slashing: slashingConfig{
				SlashFractionDuplicateSign: uint32(10),
				DuplicateSignReportReward:  uint32(50),
				MaxEvidenceAge:             uint32(27),
				SlashBlocksReward:          uint32(2500),
				ZeroProduceCumulativeTime:  uint16(20),
				ZeroProduceNumberThreshold: uint16(1),
				ZeroProduceFreezeDuration:  uint64(56),
			},
			Gov: governanceConfig{
				VersionProposalVoteDurationSeconds: uint64(14 * 24 * 3600),
				//VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposalSupportRate:       6670,
				TextProposalVoteDurationSeconds:  uint64(14 * 24 * 3600),
				TextProposalVoteRate:             5000,
				TextProposalSupportRate:          6670,
				CancelProposalVoteRate:           5000,
				CancelProposalSupportRate:        6670,
				ParamProposalVoteDurationSeconds: uint64(14 * 24 * 3600),
				ParamProposalVoteRate:            5000,
				ParamProposalSupportRate:         6670,
			},
			Reward: rewardConfig{
				NewBlockRate:                 50,
				PlatONFoundationYear:         10,
				IncreaseIssuanceRatio:        250,
				TheNumberOfDelegationsReward: 20,
			},
			Restricting: restrictingConfig{
				MinimumRelease: new(big.Int).Mul(one, new(big.Int).SetInt64(100)),
			},
			InnerAcc: innerAccount{
				PlatONFundAccount: common.Bech32ToAddressWithoutCheckHrp("lat1ptchnk2k8nh2uavp9nz8mchmmz0lhxgvztue0j"),
				PlatONFundBalance: new(big.Int).SetInt64(0),
				CDFAccount:        common.Bech32ToAddressWithoutCheckHrp("lat19au5e52762l3ffsa9uzft9hzhxk5x0g9zmcxdq"),
				CDFBalance:        new(big.Int).Set(cdfundBalance),
			},
		}
	case DefaultTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				MaxEpochMinutes:     uint64(360), // 6 hours
				NodeBlockTimeWindow: uint64(20),  // 20 seconds
				PerRoundBlocks:      uint64(10),
				MaxConsensusVals:    uint64(25),
				AdditionalCycleTime: uint64(525960),
			},
			Staking: stakingConfig{
				StakeThreshold:          new(big.Int).Set(StakeLowerLimit),
				OperatingThreshold:      new(big.Int).Set(DelegateLowerLimit),
				MaxValidators:           uint64(101),
				UnStakeFreezeDuration:   uint64(2), // freezing 2 epoch
				RewardPerMaxChangeRange: uint16(500),
				RewardPerChangeInterval: uint16(10),
			},
			Slashing: slashingConfig{
				SlashFractionDuplicateSign: uint32(10),
				DuplicateSignReportReward:  uint32(50),
				MaxEvidenceAge:             uint32(1),
				SlashBlocksReward:          uint32(250),
				ZeroProduceCumulativeTime:  uint16(30),
				ZeroProduceNumberThreshold: uint16(1),
				ZeroProduceFreezeDuration:  uint64(1),
			},
			Gov: governanceConfig{
				VersionProposalVoteDurationSeconds: uint64(14 * 24 * 3600),
				//VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposalSupportRate:       6670,
				TextProposalVoteDurationSeconds:  uint64(14 * 24 * 3600),
				TextProposalVoteRate:             5000,
				TextProposalSupportRate:          6670,
				CancelProposalVoteRate:           5000,
				CancelProposalSupportRate:        6670,
				ParamProposalVoteDurationSeconds: uint64(24 * 3600),
				ParamProposalVoteRate:            5000,
				ParamProposalSupportRate:         6670,
			},
			Reward: rewardConfig{
				NewBlockRate:                 50,
				PlatONFoundationYear:         10,
				IncreaseIssuanceRatio:        250,
				TheNumberOfDelegationsReward: 20,
			},
			Restricting: restrictingConfig{
				MinimumRelease: new(big.Int).Set(FloorMinimumRelease),
			},
			InnerAcc: innerAccount{
				PlatONFundAccount: common.HexToAddress("0x01C71CecaeFF76b78325577E6a74A94D24A86BE2"),
				PlatONFundBalance: new(big.Int).SetInt64(0),
				CDFAccount:        common.HexToAddress("0x02CddA362DCA508709a651fDe1513b22D3C2a4e5"),
				CDFBalance:        new(big.Int).Set(cdfundBalance),
			},
		}
	case DefaultUnitTestNet:
		ec = &EconomicModel{
			Common: commonConfig{
				MaxEpochMinutes:     uint64(6),  // 6 minutes
				NodeBlockTimeWindow: uint64(10), // 10 seconds
				PerRoundBlocks:      uint64(10),
				MaxConsensusVals:    uint64(4),
				AdditionalCycleTime: uint64(28),
			},
			Staking: stakingConfig{
				StakeThreshold:          new(big.Int).Set(StakeLowerLimit),
				OperatingThreshold:      new(big.Int).Set(DelegateLowerLimit),
				MaxValidators:           uint64(25),
				UnStakeFreezeDuration:   uint64(2),
				RewardPerMaxChangeRange: uint16(500),
				RewardPerChangeInterval: uint16(10),
			},
			Slashing: slashingConfig{
				SlashFractionDuplicateSign: uint32(10),
				DuplicateSignReportReward:  uint32(50),
				MaxEvidenceAge:             uint32(1),
				SlashBlocksReward:          uint32(0),
				ZeroProduceCumulativeTime:  uint16(3),
				ZeroProduceNumberThreshold: uint16(2),
				ZeroProduceFreezeDuration:  uint64(1),
			},
			Gov: governanceConfig{
				VersionProposalVoteDurationSeconds: uint64(160),
				//VersionProposalActive_ConsensusRounds: uint64(5),
				VersionProposalSupportRate:       6670,
				TextProposalVoteDurationSeconds:  uint64(160),
				TextProposalVoteRate:             5000,
				TextProposalSupportRate:          6670,
				CancelProposalVoteRate:           5000,
				CancelProposalSupportRate:        6670,
				ParamProposalVoteDurationSeconds: uint64(160),
				ParamProposalVoteRate:            5000,
				ParamProposalSupportRate:         6670,
			},
			Reward: rewardConfig{
				NewBlockRate:                 50,
				PlatONFoundationYear:         10,
				IncreaseIssuanceRatio:        250,
				TheNumberOfDelegationsReward: 2,
			},
			Restricting: restrictingConfig{
				MinimumRelease: new(big.Int).Set(FloorMinimumRelease),
			},
			InnerAcc: innerAccount{
				PlatONFundAccount: common.HexToAddress("0x493301712671Ada506ba6Ca7891F436D29185821"),
				PlatONFundBalance: new(big.Int).SetInt64(0),
				CDFAccount:        common.HexToAddress("0xC1f330B214668beAc2E6418Dd651B09C759a4Bf5"),
				CDFBalance:        new(big.Int).Set(cdfundBalance),
			},
		}
	default: // DefaultTestNet
		log.Error("not support chainID", "netId", netId)
		return nil
	}

	return ec
}

func CheckStakeThreshold(threshold *big.Int) error {

	if threshold.Cmp(StakeLowerLimit) < 0 || threshold.Cmp(StakeUpperLimit) > 0 {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The StakeThreshold must be [%d, %d] LAT", StakeLowerLimit, StakeUpperLimit))
	}
	return nil
}

func CheckOperatingThreshold(threshold *big.Int) error {
	if threshold.Cmp(DelegateLowerLimit) < 0 || threshold.Cmp(DelegateUpperLimit) > 0 {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The OperatingThreshold must be [%d, %d] LAT ", DelegateLowerLimit, DelegateUpperLimit))
	}
	return nil
}

func CheckMaxValidators(num int) error {
	if num < int(ec.Common.MaxConsensusVals) || num > CeilMaxValidators {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The MaxValidators must be [%d, %d]", int(ec.Common.MaxConsensusVals), CeilMaxValidators))
	}
	return nil
}

func CheckUnStakeFreezeDuration(duration, maxEvidenceAge, zeroProduceFreezeDuration int) error {
	if duration <= maxEvidenceAge || duration > CeilUnStakeFreezeDuration {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The UnStakeFreezeDuration must be (%d, %d]", maxEvidenceAge, CeilUnStakeFreezeDuration))
	}
	if duration <= zeroProduceFreezeDuration || duration > CeilUnStakeFreezeDuration {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The UnStakeFreezeDuration must be (%d, %d]", zeroProduceFreezeDuration, CeilUnStakeFreezeDuration))
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

func CheckZeroProduceCumulativeTime(zeroProduceCumulativeTime uint16, zeroProduceNumberThreshold uint16) error {
	if zeroProduceCumulativeTime < zeroProduceNumberThreshold || zeroProduceCumulativeTime > maxZeroProduceCumulativeTime {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The ZeroProduceCumulativeTime must be [%d, %d]", zeroProduceNumberThreshold, maxZeroProduceCumulativeTime))
	}
	return nil
}

func CheckZeroProduceNumberThreshold(zeroProduceCumulativeTime uint16, zeroProduceNumberThreshold uint16) error {
	if zeroProduceNumberThreshold < 1 || zeroProduceNumberThreshold > zeroProduceCumulativeTime {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The ZeroProduceNumberThreshold must be [%d, %d]", 1, zeroProduceCumulativeTime))
	}
	return nil
}

func CheckRewardPerMaxChangeRange(rewardPerMaxChangeRange uint16) error {
	if rewardPerMaxChangeRange < RewardPerMaxChangeRangeLowerLimit || rewardPerMaxChangeRange > RewardPerMaxChangeRangeUpperLimit {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The RewardPerMaxChangeRange must be [%d, %d]", RewardPerMaxChangeRangeLowerLimit, RewardPerMaxChangeRangeUpperLimit))
	}
	return nil
}

func CheckRewardPerChangeInterval(rewardPerChangeInterval uint16) error {
	if rewardPerChangeInterval < RewardPerChangeIntervalLowerLimit || rewardPerChangeInterval > RewardPerChangeIntervalUpperLimit {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The RewardPerChangeInterval must be [%d, %d]", RewardPerChangeIntervalLowerLimit, RewardPerChangeIntervalUpperLimit))
	}
	return nil
}

func CheckIncreaseIssuanceRatio(increaseIssuanceRatio uint16) error {
	if increaseIssuanceRatio < IncreaseIssuanceRatioLowerLimit || increaseIssuanceRatio > IncreaseIssuanceRatioUpperLimit {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The IncreaseIssuanceRatio must be [%d, %d]", IncreaseIssuanceRatioLowerLimit, IncreaseIssuanceRatioUpperLimit))
	}
	return nil
}

func CheckZeroProduceFreezeDuration(zeroProduceFreezeDuration uint64, unStakeFreezeDuration uint64) error {
	if zeroProduceFreezeDuration < 1 || zeroProduceFreezeDuration >= unStakeFreezeDuration {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The ZeroProduceFreezeDuration must be [%d, %d]", 1, unStakeFreezeDuration-1))
	}
	return nil
}

func CheckMinimumRelease(minimumRelease *big.Int) error {
	if minimumRelease.Cmp(FloorMinimumRelease) < 0 || minimumRelease.Cmp(CeilMinimumRelease) > 0 {
		return common.InvalidParameter.Wrap(fmt.Sprintf("The MinimumRelease must be [%d, %d]", FloorMinimumRelease, CeilMinimumRelease))
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

	log.Info("Call CheckEconomicModel: check epoch and consensus round,", "config epoch duration", fmt.Sprintf("%d s", epochDuration),
		"perblock duration", fmt.Sprintf("%d s", blockDuration), "round duration", fmt.Sprintf("%d s", roundDuration),
		"real epoch duration", fmt.Sprintf("%d s", realEpochDuration), "consensus count of epoch", epochSize)

	if epochSize < 4 {
		return errors.New("The settlement period must be more than four times the consensus period")
	}

	// additionalCycle Size, how many epoch duration
	additionalCycleSize := ec.Common.AdditionalCycleTime * 60 / realEpochDuration
	// realAdditionalCycleDuration
	realAdditionalCycleDuration := additionalCycleSize * realEpochDuration / 60

	log.Info("Call CheckEconomicModel: additional cycle and epoch,", "config additional cycle duration", fmt.Sprintf("%d min", ec.Common.AdditionalCycleTime),
		"real additional cycle duration", fmt.Sprintf("%d min", realAdditionalCycleDuration), "epoch count of additional cycle", additionalCycleSize)

	if additionalCycleSize < 4 {
		return errors.New("The issuance period must be integer multiples of the settlement period and multiples must be greater than or equal to 4")
	}

	if ec.Common.MaxConsensusVals < FloorMaxConsensusVals || ec.Common.MaxConsensusVals > CeilMaxConsensusVals {
		return fmt.Errorf("The consensus validator num must be [%d, %d]", FloorMaxConsensusVals, CeilMaxConsensusVals)
	}

	if err := CheckMaxValidators(int(ec.Staking.MaxValidators)); nil != err {
		return err
	}

	if err := CheckOperatingThreshold(ec.Staking.OperatingThreshold); nil != err {
		return err
	}

	if err := CheckStakeThreshold(ec.Staking.StakeThreshold); nil != err {
		return err
	}

	if err := CheckUnStakeFreezeDuration(int(ec.Staking.UnStakeFreezeDuration), int(ec.Slashing.MaxEvidenceAge), int(ec.Slashing.ZeroProduceFreezeDuration)); nil != err {
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

	if err := CheckZeroProduceNumberThreshold(ec.Slashing.ZeroProduceCumulativeTime, ec.Slashing.ZeroProduceNumberThreshold); nil != err {
		return err
	}

	if err := CheckZeroProduceCumulativeTime(ec.Slashing.ZeroProduceCumulativeTime, ec.Slashing.ZeroProduceNumberThreshold); nil != err {
		return err
	}

	if err := CheckRewardPerMaxChangeRange(ec.Staking.RewardPerMaxChangeRange); nil != err {
		return err
	}

	if err := CheckRewardPerChangeInterval(ec.Staking.RewardPerChangeInterval); nil != err {
		return err
	}

	if err := CheckIncreaseIssuanceRatio(ec.Reward.IncreaseIssuanceRatio); nil != err {
		return err
	}

	if err := CheckZeroProduceFreezeDuration(ec.Slashing.ZeroProduceFreezeDuration, ec.Staking.UnStakeFreezeDuration); nil != err {
		return err
	}

	if err := CheckMinimumRelease(ec.Restricting.MinimumRelease); nil != err {
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

func ConsensusSize() uint64 {
	return BlocksWillCreate() * MaxConsensusVals()
}

func EpochSize() uint64 {
	consensusSize := ConsensusSize()
	em := MaxEpochMinutes()
	i := Interval()

	epochSize := em * 60 / (i * consensusSize)
	return epochSize
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
	return 1
}

func ElectionDistance() uint64 {
	// min need two view
	return 2 * ec.Common.PerRoundBlocks
}

func UnStakeFreezeDuration() uint64 {
	return ec.Staking.UnStakeFreezeDuration
}

func RewardPerMaxChangeRange() uint16 {
	return ec.Staking.RewardPerMaxChangeRange
}

func RewardPerChangeInterval() uint16 {
	return ec.Staking.RewardPerChangeInterval
}

/******
 * Restricting config
 ******/

func RestrictingMinimumRelease() *big.Int {
	return ec.Restricting.MinimumRelease
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

func ZeroProduceCumulativeTime() uint16 {
	return ec.Slashing.ZeroProduceCumulativeTime
}

func ZeroProduceNumberThreshold() uint16 {
	return ec.Slashing.ZeroProduceNumberThreshold
}

func ZeroProduceFreezeDuration() uint64 {
	return ec.Slashing.ZeroProduceFreezeDuration
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

func IncreaseIssuanceRatio() uint16 {
	return ec.Reward.IncreaseIssuanceRatio
}

func TheNumberOfDelegationsReward() uint16 {
	return ec.Reward.TheNumberOfDelegationsReward
}

/******
 * Governance config
 ******/
/*func VersionProposalVote_ConsensusRounds() uint64 {
	return ec.Gov.VersionProposalVoteDurationSeconds / (Interval() * ec.Common.PerRoundBlocks * ec.Common.MaxConsensusVals)
}*/

func VersionProposalVote_DurationSeconds() uint64 {
	return ec.Gov.VersionProposalVoteDurationSeconds
}

/*func VersionProposalActive_ConsensusRounds() uint64 {
	return ec.Gov.VersionProposalActive_ConsensusRounds
}*/

func VersionProposal_SupportRate() uint64 {
	return ec.Gov.VersionProposalSupportRate
}

/*func TextProposalVote_ConsensusRounds() uint64 {
	return ec.Gov.TextProposalVoteDurationSeconds / (Interval() * ec.Common.PerRoundBlocks * ec.Common.MaxConsensusVals)
}*/
func TextProposalVote_DurationSeconds() uint64 {
	return ec.Gov.TextProposalVoteDurationSeconds
}
func TextProposal_VoteRate() uint64 {
	return ec.Gov.TextProposalVoteRate
}

func TextProposal_SupportRate() uint64 {
	return ec.Gov.TextProposalSupportRate
}

func CancelProposal_VoteRate() uint64 {
	return ec.Gov.CancelProposalVoteRate
}

func CancelProposal_SupportRate() uint64 {
	return ec.Gov.CancelProposalSupportRate
}

func ParamProposalVote_DurationSeconds() uint64 {
	return ec.Gov.ParamProposalVoteDurationSeconds
}

func ParamProposal_VoteRate() uint64 {
	return ec.Gov.ParamProposalVoteRate
}

func ParamProposal_SupportRate() uint64 {
	return ec.Gov.ParamProposalSupportRate
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

// Calculate the P value of the binomial distribution
// Parameter: The total weight of the election
func CalcP(totalWeight float64, sqrtWeight float64) float64 {
	return float64(ElectionBase) / sqrtWeight
}
