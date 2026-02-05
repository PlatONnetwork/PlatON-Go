// Copyright 2021 The PlatON Network Authors
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

package gov

import (
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

var (
	initGovParam sync.Once
)

var governParam []*GovernParam

func queryInitParam() []*GovernParam {
	initGovParam.Do(func() {
		log.Info("Init Govern parameters ...")
		governParam = initParam()
	})
	return governParam
}

func initParam() []*GovernParam {
	return []*GovernParam{

		/**
		About Staking module
		*/
		{
			ParamItem: &ParamItem{ModuleStaking, KeyStakeThreshold,
				fmt.Sprintf("minimum amount of stake, range: [%d, %d]", params.StakeLowerLimit, params.StakeUpperLimit)},
			ParamValue: &ParamValue{"", params.StakeThreshold().String(), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				threshold, ok := new(big.Int).SetString(value, 10)
				if !ok {
					return fmt.Errorf("Parsed StakeThreshold is failed")
				}

				if err := params.CheckStakeThreshold(threshold); nil != err {
					return err
				}
				return nil
			},
		},

		{
			ParamItem: &ParamItem{ModuleStaking, KeyOperatingThreshold,
				fmt.Sprintf("minimum amount of stake increasing funds, delegation funds, or delegation withdrawing funds, range: [%d, %d]", params.DelegateLowerLimit, params.DelegateUpperLimit)},
			ParamValue: &ParamValue{"", params.OperatingThreshold().String(), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				threshold, ok := new(big.Int).SetString(value, 10)
				if !ok {
					return fmt.Errorf("Parsed OperatingThreshold is failed")
				}

				if err := params.CheckOperatingThreshold(threshold); nil != err {
					return err
				}

				return nil
			},
		},

		{
			ParamItem: &ParamItem{ModuleStaking, KeyMaxValidators,
				fmt.Sprintf("maximum amount of validator, range: [%d, %d]", params.MaxConsensusVals(), params.CeilMaxValidators)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.MaxValidators())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				num, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxValidators is failed: %v", err)
				}

				if err := params.CheckMaxValidators(num); nil != err {
					return err
				}
				return nil
			},
		},

		{
			ParamItem: &ParamItem{ModuleStaking, KeyUnStakeFreezeDuration,
				fmt.Sprintf("quantity of epoch for skake withdrawal, range: (MaxEvidenceAge, %d]", params.CeilUnStakeFreezeDuration)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.UnStakeFreezeDuration())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				num, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed UnStakeFreezeDuration is failed: %v", err)
				}

				age, err := gov.GovernMaxEvidenceAge(blockNumber, blockHash)
				if nil != err {
					return err
				}
				epochNumber, err := gov.GovernZeroProduceFreezeDuration(blockNumber, blockHash)
				if nil != err {
					return err
				}
				if err := params.CheckUnStakeFreezeDuration(num, int(age), int(epochNumber)); nil != err {
					return err
				}

				return nil
			},
		},

		/**
		About Slashing module
		*/

		{
			ParamItem: &ParamItem{ModuleSlashing, KeySlashFractionDuplicateSign,
				fmt.Sprintf("quantity of base point(1BP=1‱). Node's stake will be deducted(BPs*staking amount*1‱) it the node sign block duplicatlly, range: (%d, %d]", params.Zero, params.TenThousand)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.SlashFractionDuplicateSign())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				fraction, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed SlashFractionDuplicateSign is failed: %v", err)
				}

				if err := params.CheckSlashFractionDuplicateSign(fraction); nil != err {
					return err
				}
				return nil
			},
		},

		{
			ParamItem: &ParamItem{ModuleSlashing, KeyDuplicateSignReportReward,
				fmt.Sprintf("quantity of base point(1bp=1%%). Bonus(BPs*deduction amount for sign block duplicatlly*%%) to the node who reported another's duplicated-signature, range: (%d, %d]", params.Zero, params.Eighty)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.DuplicateSignReportReward())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				fraction, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed DuplicateSignReportReward is failed: %v", err)
				}

				if err := params.CheckDuplicateSignReportReward(fraction); nil != err {
					return err
				}
				return nil
			},
		},

		{
			ParamItem: &ParamItem{ModuleSlashing, KeyMaxEvidenceAge,
				fmt.Sprintf("quantity of epoch. During these epochs after a node duplicated-sign, others can report it, range: (%d, UnStakeFreezeDuration)", params.Zero)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.MaxEvidenceAge())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				age, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxEvidenceAge is failed: %v", err)
				}

				duration, err := gov.GovernUnStakeFreezeDuration(blockNumber, blockHash)
				if nil != err {
					return err
				}
				if err := params.CheckMaxEvidenceAge(age, int(duration)); nil != err {
					return err
				}
				return nil
			},
		},
		{
			ParamItem: &ParamItem{ModuleSlashing, KeySlashBlocksReward,
				fmt.Sprintf("quantity of block, the total bonus amount for these blocks will be deducted from a inefficient node's stake, range: [%d, %d)", params.Zero, params.CeilBlocksReward)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.SlashBlocksReward())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				rewards, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed SlashBlocksReward is failed: %v", err)
				}

				if err := params.CheckSlashBlocksReward(rewards); nil != err {
					return err
				}
				return nil
			},
		},

		/**
		About Block module
		*/
		{
			ParamItem:  &ParamItem{ModuleBlock, KeyMaxBlockGasLimit, fmt.Sprintf("maximum gas limit per block, range: [%d, %d]", int(params.GenesisGasLimit), int(params.MaxGasCeil))},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.DefaultMinerGasCeil)), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				gasLimit, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxBlockGasLimit is failed: %v", err)
				}

				// (4712388<= x <=21000 0000)
				if gasLimit < int(params.GenesisGasLimit) || gasLimit > int(params.MaxGasCeil) {
					return common.InvalidParameter.Wrap(fmt.Sprintf("The MaxBlockGasLimit must be [%d, %d]", int(params.GenesisGasLimit), int(params.MaxGasCeil)))
				}
				return nil
			},
		},
		{

			ParamItem: &ParamItem{ModuleSlashing, KeyZeroProduceCumulativeTime,
				fmt.Sprintf("Time range for recording the number of behaviors of zero production blocks, range: [ZeroProduceNumberThreshold, %d]", params.MaxZeroProduceCumulativeTime)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.ZeroProduceCumulativeTime())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				roundNumber, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("parsed ZeroProduceCumulativeTime is failed")
				}

				numberThreshold, err := gov.GovernZeroProduceNumberThreshold(blockNumber, blockHash)
				if nil != err {
					return err
				}
				if err := params.CheckZeroProduceCumulativeTime(uint16(roundNumber), numberThreshold); nil != err {
					return err
				}
				return nil
			},
		},
		{

			ParamItem: &ParamItem{ModuleSlashing, KeyZeroProduceNumberThreshold,
				"Number of zero production blocks, range: [1, ZeroProduceCumulativeTime]"},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.ZeroProduceNumberThreshold())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				number, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("parsed ZeroProduceNumberThreshold is failed")
				}

				roundNumber, err := gov.GovernZeroProduceCumulativeTime(blockNumber, blockHash)
				if nil != err {
					return err
				}
				if err := params.CheckZeroProduceNumberThreshold(roundNumber, uint16(number)); nil != err {
					return err
				}
				return nil
			},
		},
		{

			ParamItem: &ParamItem{ModuleStaking, KeyRewardPerMaxChangeRange,
				fmt.Sprintf("Delegated Reward Ratio The maximum adjustable range of each modification, range: [%d, %d]", params.RewardPerMaxChangeRangeLowerLimit, params.RewardPerMaxChangeRangeUpperLimit)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.RewardPerMaxChangeRange())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				number, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("parsed RewardPerMaxChangeRange is failed")
				}

				if err := params.CheckRewardPerMaxChangeRange(uint16(number)); nil != err {
					return err
				}
				return nil
			},
		},
		{

			ParamItem: &ParamItem{ModuleStaking, KeyRewardPerChangeInterval,
				fmt.Sprintf("The interval for each modification of the commission reward ratio, range: [%d, %d]", params.RewardPerChangeIntervalLowerLimit, params.RewardPerChangeIntervalUpperLimit)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.RewardPerChangeInterval())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				number, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("parsed RewardPerChangeInterval is failed")
				}

				if err := params.CheckRewardPerChangeInterval(uint16(number)); nil != err {
					return err
				}
				return nil
			},
		},
		{

			ParamItem: &ParamItem{ModuleReward, KeyIncreaseIssuanceRatio,
				fmt.Sprintf("Increase the ratio of issuance, range: [%d, %d]", params.IncreaseIssuanceRatioLowerLimit, params.IncreaseIssuanceRatioUpperLimit)},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.IncreaseIssuanceRatio())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				number, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("parsed IncreaseIssuanceRatio is failed")
				}

				if err := params.CheckIncreaseIssuanceRatio(uint16(number)); nil != err {
					return err
				}
				return nil
			},
		},
		{

			ParamItem: &ParamItem{ModuleSlashing, KeyZeroProduceFreezeDuration,
				"Zero production frozen time, range: [1, UnStakeFreezeDuration)"},
			ParamValue: &ParamValue{"", strconv.Itoa(int(params.ZeroProduceFreezeDuration())), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				number, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("parsed KeyZeroProduceFreezeDuration is failed")
				}

				epochNumber, err := NewGov(snapshotdb.Instance()).GovernUnStakeFreezeDuration(blockNumber, blockHash)
				if nil != err {
					return err
				}

				if err := params.CheckZeroProduceFreezeDuration(uint64(number), epochNumber); nil != err {
					return err
				}
				return nil
			},
		},
		{

			ParamItem: &ParamItem{ModuleRestricting, KeyRestrictingMinimumAmount,
				fmt.Sprintf("minimum restricting amount to be released in each epoch, range: [%d, %d]",
					params.FloorMinimumRelease, params.CeilMinimumRelease)},
			ParamValue: &ParamValue{"", params.RestrictingMinimumRelease().String(), 0},
			ParamVerifier: func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
				v, ok := new(big.Int).SetString(value, 10)
				if !ok {
					return fmt.Errorf("parsed KeyRestrictingMinimumAmount is failed")
				}
				if err := params.CheckMinimumRelease(v); err != nil {
					return err
				}
				return nil
			},
		},
	}
}

var ParamVerifierMap = make(map[string]ParamVerifier)

func InitGenesisGovernParam(prevHash common.Hash, snapDB snapshotdb.BaseDB, genesisVersion uint32) (common.Hash, error) {
	var paramItemList []*ParamItem

	initParamList := queryInitParam()

	if genesisVersion >= params.FORKVERSION_1_3_0 {
		log.Info("init 1.3.0 params")
		initParamList = append(initParamList, initUnDelegateFreezeDurationParamGenesis())
	}

	putBasedb_genKVHash_Fn := func(key, val []byte, hash common.Hash) (common.Hash, error) {
		if err := snapDB.PutBaseDB(key, val); nil != err {
			return common.ZeroHash, err
		}
		newHash := common.GenerateKVHash(key, val, hash)
		return newHash, nil
	}

	var lastHash = prevHash
	var err error
	for _, param := range initParamList {
		paramItemList = append(paramItemList, param.ParamItem)

		key := KeyParamValue(param.ParamItem.Module, param.ParamItem.Name)
		value := common.MustRlpEncode(param.ParamValue)
		lastHash, err = putBasedb_genKVHash_Fn(key, value, lastHash)
		if nil != err {
			return lastHash, fmt.Errorf("failed to Store govern parameter: PutBaseDB failed. ParamItem:%s, ParamValue:%s, error:%s", param.ParamItem.Module, param.ParamItem.Name, err.Error())
		}
	}

	key := KeyParamItems()
	value := common.MustRlpEncode(paramItemList)
	lastHash, err = putBasedb_genKVHash_Fn(key, value, lastHash)
	if nil != err {
		return lastHash, fmt.Errorf("failed to Store govern parameter list: PutBaseDB failed. error:%s", err.Error())
	}

	//stateDB.SetState(vm.GovContractAddr, KeyGovernHASHKey(), lastHash.Bytes())
	return lastHash, nil
}

func initUnDelegateFreezeDurationParamGenesis() *GovernParam {
	return &GovernParam{
		ParamItem: &ParamItem{ModuleStaking, KeyUnDelegateFreezeDuration,
			fmt.Sprintf("quantity of epoch for delegate withdrawal, range:  [%d, UnStakeFreezeDuration]", 1)},
		ParamValue:    &ParamValue{"", strconv.Itoa(int(params.UnDelegateFreezeDuration())), 0},
		ParamVerifier: UnDelegateFreezeDurationVerifier,
	}
}

func initUnDelegateFreezeDurationParamVersionUpdate(blockNumber uint64, blockHash common.Hash) (*GovernParam, error) {
	Duration, err := NewGov(snapshotdb.Instance()).GovernUnStakeFreezeDuration(blockNumber, blockHash)
	if nil != err {
		return nil, err
	}
	unDelegateFreezeDuration := 0
	if Duration/3 == 0 {
		unDelegateFreezeDuration = 1
	} else {
		unDelegateFreezeDuration = int(Duration / 3)
	}
	return &GovernParam{
		ParamItem: &ParamItem{ModuleStaking, KeyUnDelegateFreezeDuration,
			fmt.Sprintf("quantity of epoch for delegate withdrawal, range:  [%d, UnStakeFreezeDuration]", 1)},
		ParamValue:    &ParamValue{"", strconv.Itoa(unDelegateFreezeDuration), blockNumber},
		ParamVerifier: UnDelegateFreezeDurationVerifier,
	}, nil
}

var UnDelegateFreezeDurationVerifier = func(gov *Gov, blockNumber uint64, blockHash common.Hash, value string) error {
	num, err := strconv.Atoi(value)
	if nil != err {
		return fmt.Errorf("Parsed UnDelegateFreezeDuration is failed: %v", err)
	}

	Duration, err := gov.GovernUnStakeFreezeDuration(blockNumber, blockHash)
	if nil != err {
		return err
	}
	if err := params.CheckUnDelegateFreezeDuration(num, int(Duration)); nil != err {
		return err
	}
	return nil
}

func RegisterGovernParamVerifiers() {
	for _, param := range queryInitParam() {
		RegGovernParamVerifier(param.ParamItem.Module, param.ParamItem.Name, param.ParamVerifier)
	}

	RegGovernParamVerifier(ModuleStaking, KeyUnDelegateFreezeDuration, UnDelegateFreezeDurationVerifier)
}

func RegGovernParamVerifier(module, name string, callback ParamVerifier) {
	ParamVerifierMap[module+"/"+name] = callback
}
