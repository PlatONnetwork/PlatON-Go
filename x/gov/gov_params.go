package gov

import (
	"fmt"
	"math/big"
	"strconv"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var (
	initGovParam sync.Once
)

var governParam []*xcom.GovernParam

func queryInitParam() []*xcom.GovernParam {
	initGovParam.Do(func() {
		log.Info("Init Govern parameters ...")
		governParam = initParam()
	})
	return governParam
}

func initParam() []*xcom.GovernParam {
	return []*xcom.GovernParam{

		/**
		About Staking module
		*/
		{

			ParamItem:  &xcom.ParamItem{xcom.ModuleStaking, xcom.KeyStakeThreshold, fmt.Sprintf("minimum amount of stake, range：[%d, %s) LAT", xcom.MillionLAT, xcom.PositiveInfinity)},
			ParamValue: &xcom.ParamValue{"", xcom.StakeThreshold().String(), 0},
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
			ParamItem:  &xcom.ParamItem{xcom.ModuleStaking, xcom.KeyOperatingThreshold, fmt.Sprintf("minimum amount of stake increasing funds, delegation funds, or delegation withdrawing funds, range：[%d, %s) LAT", xcom.TenLAT, xcom.PositiveInfinity)},
			ParamValue: &xcom.ParamValue{"", xcom.OperatingThreshold().String(), 0},
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
			ParamItem:  &xcom.ParamItem{xcom.ModuleStaking, xcom.KeyMaxValidators, fmt.Sprintf("maximum amount of validator, range：[%d, %d]", xcom.CeilMaxConsensusVals, xcom.CeilMaxValidators)},
			ParamValue: &xcom.ParamValue{"", strconv.Itoa(int(xcom.MaxValidators())), 0},
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
			ParamItem:  &xcom.ParamItem{xcom.ModuleStaking, xcom.KeyUnStakeFreezeDuration, fmt.Sprintf("quantity of epoch for skake withdrawal, range：(MaxEvidenceAge, %d]", xcom.CeilUnStakeFreezeDuration)},
			ParamValue: &xcom.ParamValue{"", strconv.Itoa(int(xcom.UnStakeFreezeDuration())), 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				num, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed UnStakeFreezeDuration is failed: %v", err)
				}

				ageStr, err := xcom.GetGovernParamValue(xcom.ModuleSlashing, xcom.KeyMaxEvidenceAge, blockNumber, blockHash)
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
			ParamItem:  &xcom.ParamItem{xcom.ModuleSlashing, xcom.KeySlashFractionDuplicateSign, fmt.Sprintf("quantity of base point(1BP=1‱). Node's stake will be deducted(BPs*staking amount*1‱) it the node sign block duplicatlly, range：(%d, %d]", xcom.Zero, xcom.TenThousand)},
			ParamValue: &xcom.ParamValue{"", strconv.Itoa(int(xcom.SlashFractionDuplicateSign())), 0},
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
			ParamItem:  &xcom.ParamItem{xcom.ModuleSlashing, xcom.KeyDuplicateSignReportReward, fmt.Sprintf("quantity of base point(1bp=1%%). Bonus(BPs*deduction amount for sign block duplicatlly*%%) to the node who reported another's duplicated-signature, range：(%d, %d]", xcom.Zero, xcom.Eighty)},
			ParamValue: &xcom.ParamValue{"", strconv.Itoa(int(xcom.DuplicateSignReportReward())), 0},
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
			ParamItem:  &xcom.ParamItem{xcom.ModuleSlashing, xcom.KeyMaxEvidenceAge, fmt.Sprintf("quantity of epoch. During these epochs after a node duplicated-sign, others can report it, range：(%d, UnStakeFreezeDuration)", xcom.Zero)},
			ParamValue: &xcom.ParamValue{"", strconv.Itoa(int(xcom.MaxEvidenceAge())), 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				age, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxEvidenceAge is failed: %v", err)
				}

				durationStr, err := xcom.GetGovernParamValue(xcom.ModuleStaking, xcom.KeyUnStakeFreezeDuration, blockNumber, blockHash)
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
			ParamItem:  &xcom.ParamItem{xcom.ModuleSlashing, xcom.KeySlashBlocksReward, fmt.Sprintf("quantity of block, the total bonus amount for these blocks will be deducted from a inefficient node's stake, range：[%d, %d)", xcom.Zero, xcom.CeilBlocksReward)},
			ParamValue: &xcom.ParamValue{"", strconv.Itoa(int(xcom.SlashBlocksReward())), 0},
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
			ParamItem:  &xcom.ParamItem{xcom.ModuleBlock, xcom.KeyMaxBlockGasLimit, fmt.Sprintf("maximum gas limit per block, range：(%d, %s)", xcom.Zero, xcom.PositiveInfinity)},
			ParamValue: &xcom.ParamValue{"", strconv.Itoa(int(params.GenesisGasLimit)), 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				gasLimit, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxBlockGasLimit is failed: %v", err)
				}

				// (0, +∞)
				if gasLimit <= 0 {
					return common.InvalidParameter.Wrap(fmt.Sprintf("The MaxBlockGasLimit must be (%d, %s)", xcom.Zero, xcom.PositiveInfinity))
				}

				return nil
			},
		},

		/**
		About TxPool module
		*/
		{
			ParamItem:  &xcom.ParamItem{xcom.ModuleTxPool, xcom.KeyMaxTxDataLimit, fmt.Sprintf("maximum data length per transaction, range：(%d, %d]", xcom.Zero, xcom.CeilTxSize)},
			ParamValue: &xcom.ParamValue{"", strconv.Itoa(xcom.GenesisTxSize), 0},
			ParamVerifier: func(blockNumber uint64, blockHash common.Hash, value string) error {

				txSize, err := strconv.Atoi(value)
				if nil != err {
					return fmt.Errorf("Parsed MaxTxDataLimit is failed: %v", err)
				}

				// (0, 10MB]
				if txSize > xcom.CeilTxSize {
					return common.InvalidParameter.Wrap(fmt.Sprintf("The MaxTxDataLimit must be (%d, %d]", xcom.Zero, xcom.CeilTxSize))
				}

				return nil
			},
		},
	}
}

var ParamVerifierMap = make(map[string]xcom.ParamVerifier)

func InitGenesisGovernParam(snapDB snapshotdb.DB) error {
	var paramItemList []*xcom.ParamItem
	for _, param := range queryInitParam() {
		paramItemList = append(paramItemList, param.ParamItem)

		key := xcom.KeyParamValue(param.ParamItem.Module, param.ParamItem.Name)
		value := common.MustRlpEncode(param.ParamValue)
		if err := snapDB.PutBaseDB(key, value); err != nil {
			return err
		}
	}

	key := xcom.KeyParamItems()
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

func RegGovernParamVerifier(module, name string, callback xcom.ParamVerifier) {
	ParamVerifierMap[module+"/"+name] = callback
}
