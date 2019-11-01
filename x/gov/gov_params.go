package gov

import (
	"fmt"
	"math/big"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

const (
	ModuleStaking  = "Staking"
	ModuleSlashing = "Slashing"
	ModuleBlock    = "Block"
	ModuleTxPool   = "TxPool"
)

const (
	KeyStakeThreshold          = "StakeThreshold"
	KeyOperatingThreshold      = "OperatingThreshold"
	MaxValidators              = "MaxValidators"
	UnStakeFreezeDuration      = "UnStakeFreezeDuration"
	SlashFractionDuplicateSign = "SlashFractionDuplicateSign"
	DuplicateSignReportReward  = "DuplicateSignReportReward"
	MaxEvidenceAge             = "MaxEvidenceAge"
	SlashBlocksReward          = "SlashBlocksReward"
	MaxBlockGasLimit           = "MaxBlockGasLimit"
	MaxTxDataLimit             = "MaxTxDataLimit"
)

const (
	genesisTxSize = 1024 * 1024        //  1 MB
	ceilTxSize    = genesisTxSize * 10 // 10 MB

)

func queryInitParam() []*GovernParam {
	return []*GovernParam{

		/**
		About Staking module
		*/
		{

			ParamItem:  &ParamItem{ModuleStaking, KeyStakeThreshold, fmt.Sprintf("xxxxx, range：[%d, %s) LAT", xcom.MillionLAT, xcom.PositiveInfinity)},
			ParamValue: &ParamValue{xcom.StakeThreshold().String(), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				threshold, ok := new(big.Int).SetString(value, 10)
				if !ok {
					return fmt.Errorf("Parsed StakeThreshold is failed")
				}

				if err := xcom.CheckStakeThreshold(threshold); nil != err {
					return err
				}
				return nil
			},
		},

		{
			ParamItem:  &ParamItem{ModuleStaking, KeyOperatingThreshold, fmt.Sprintf("xxxxx, range：[%d, %s) LAT", xcom.TenLAT, xcom.PositiveInfinity)},
			ParamValue: &ParamValue{xcom.OperatingThreshold().String(), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				threshold, ok := new(big.Int).SetString(value, 10)
				if !ok {
					return fmt.Errorf("Parsed OperatingThreshold is failed")
				}

				if err := xcom.CheckOperatingThreshold(threshold); nil != err {
					return err
				}

				return nil

			},
		},

		{
			ParamItem:  &ParamItem{ModuleStaking, MaxValidators, fmt.Sprintf("xxxxx, range：[%d, %d]", xcom.CeilMaxConsensusVals, xcom.CeilMaxValidators)},
			ParamValue: &ParamValue{strconv.Itoa(int(xcom.MaxValidators())), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				num, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxValidators is failed: %v", err)
				}

				if err := xcom.CheckMaxValidators(num); nil != err {
					return err
				}

				return nil

			},
		},

		{
			ParamItem:  &ParamItem{ModuleStaking, UnStakeFreezeDuration, fmt.Sprintf("xxxxx, range：(MaxEvidenceAge, %d]", xcom.CeilUnStakeFreezeDuration)},
			ParamValue: &ParamValue{strconv.Itoa(int(xcom.UnStakeFreezeDuration())), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				num, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed UnStakeFreezeDuration is failed: %v", err)
				}

				ageStr, err := GetGovernParamValue(ModuleSlashing, MaxEvidenceAge, blockNumber, blockHash)
				if nil != err {
					return err
				}

				age, _ := strconv.Atoi(ageStr)

				if err := xcom.CheckUnStakeFreezeDuration(num, age); nil != err {
					return err
				}

				return nil

			},
		},

		/**
		About Slashing module
		*/

		{
			ParamItem:  &ParamItem{ModuleSlashing, SlashFractionDuplicateSign, fmt.Sprintf("xxxxx, range：(%d, %d]", xcom.Zero, xcom.TenThousand)},
			ParamValue: &ParamValue{strconv.Itoa(int(xcom.SlashFractionDuplicateSign())), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				fraction, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed SlashFractionDuplicateSign is failed: %v", err)
				}

				if err := xcom.CheckSlashFractionDuplicateSign(fraction); nil != err {
					return err
				}

				return nil
			},
		},

		{
			ParamItem:  &ParamItem{ModuleSlashing, DuplicateSignReportReward, fmt.Sprintf("xxxxx, range：(%d, %d]", xcom.Zero, xcom.Eighty)},
			ParamValue: &ParamValue{strconv.Itoa(int(xcom.DuplicateSignReportReward())), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				fraction, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed DuplicateSignReportReward is failed: %v", err)
				}

				if err := xcom.CheckDuplicateSignReportReward(fraction); nil != err {
					return err
				}

				return nil
			},
		},

		{
			ParamItem:  &ParamItem{ModuleSlashing, MaxEvidenceAge, fmt.Sprintf("xxxxx, range：(%d, UnStakeFreezeDuration)", xcom.Zero)},
			ParamValue: &ParamValue{strconv.Itoa(int(xcom.MaxEvidenceAge())), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				age, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxEvidenceAge is failed: %v", err)
				}

				durationStr, err := GetGovernParamValue(ModuleStaking, UnStakeFreezeDuration, blockNumber, blockHash)
				if nil != err {
					return err
				}

				duration, _ := strconv.Atoi(durationStr)

				if err := xcom.CheckMaxEvidenceAge(age, duration); nil != err {
					return err
				}

				return nil

			},
		},
		{
			ParamItem:  &ParamItem{ModuleSlashing, SlashBlocksReward, fmt.Sprintf("xxxxx, range：[%d, %d)", xcom.Zero, xcom.CeilBlocksReward)},
			ParamValue: &ParamValue{strconv.Itoa(int(xcom.SlashBlocksReward())), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				rewards, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed SlashBlocksReward is failed: %v", err)
				}

				if err := xcom.CheckSlashBlocksReward(rewards); nil != err {
					return err
				}

				return nil

			},
		},

		/**
		About Block module
		*/
		{
			ParamItem: &ParamItem{ModuleBlock, MaxBlockGasLimit, fmt.Sprintf("xxxxx, range：(%d, %s)", xcom.Zero, xcom.PositiveInfinity)},
			//ParamValue: &ParamValue{strconv.Itoa(int(params.GenesisGasLimit)), "", 0},
			ParamValue: &ParamValue{strconv.Itoa(100800000), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				gasLimit, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxBlockGasLimit is failed: %v", err)
				}

				// (0, +∞)
				if gasLimit <= 0 {
					return fmt.Errorf("The MaxBlockGasLimit must be (%d, %s)", xcom.Zero, xcom.PositiveInfinity)
				}

				return nil
			},
		},

		/**
		About TxPool module
		*/
		{
			ParamItem:  &ParamItem{ModuleTxPool, MaxTxDataLimit, fmt.Sprintf("xxxxx, range：(%d, %d]", xcom.Zero, ceilTxSize)},
			ParamValue: &ParamValue{strconv.Itoa(genesisTxSize), "", 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				txSize, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxTxDataLimit is failed: %v", err)
				}

				// (0, 10MB]
				if txSize > ceilTxSize {
					return fmt.Errorf("The MaxBlockGasLimit must be (%d, %d]", xcom.Zero, ceilTxSize)
				}

				return nil
			},
		},
	}
}

var ParamVerifierMap = make(map[string]ParamVerifier)

func InitGenesisGovernParam(snapDB snapshotdb.DB) error {
	var paramItemList []*ParamItem
	for _, param := range queryInitParam() {
		paramItemList = append(paramItemList, param.ParamItem)

		key := KeyParamValue(param.ParamItem.Module, param.ParamItem.Name)
		value := common.MustRlpEncode(param.ParamValue)
		if err := snapDB.PutBaseDB(key, value); err != nil {
			return err
		}
	}

	key := KeyParamItems()
	value := common.MustRlpEncode(paramItemList)
	if err := snapDB.PutBaseDB(key, value); err != nil {
		return err
	}
	return nil
}

func RegisterGovernParamVerifiers() {
	for _, param := range queryInitParam() {
		RegGovernParamVerifier(param.ParamItem.Module, param.ParamItem.Name, param.ParamVerifier)
	}
}

func RegGovernParamVerifier(module, name string, callback ParamVerifier) {
	ParamVerifierMap[module+"/"+name] = callback
}
