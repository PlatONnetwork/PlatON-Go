package xcom

import (
	"bytes"
	"fmt"
	"math/big"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

const (
	ModuleStaking  = "Staking"
	ModuleSlashing = "Slashing"
	ModuleBlock    = "Block"
	ModuleTxPool   = "TxPool"
)

const (
	KeyStakeThreshold             = "StakeThreshold"
	KeyOperatingThreshold         = "OperatingThreshold"
	KeyMaxValidators              = "MaxValidators"
	KeyUnStakeFreezeDuration      = "UnStakeFreezeDuration"
	KeySlashFractionDuplicateSign = "SlashFractionDuplicateSign"
	KeyDuplicateSignReportReward  = "DuplicateSignReportReward"
	KeyMaxEvidenceAge             = "MaxEvidenceAge"
	KeySlashBlocksReward          = "SlashBlocksReward"
	KeyMaxBlockGasLimit           = "MaxBlockGasLimit"
	KeyMaxTxDataLimit             = "MaxTxDataLimit"
)

const (
	GenesisTxSize = 1024 * 1024        //  1 MB
	CeilTxSize    = GenesisTxSize * 10 // 10 MB

)

var (
	KeyDelimiter        = []byte(":")
	keyPrefixParamItems = []byte("ParamItems")
	keyPrefixParamValue = []byte("ParamValue")
)

func KeyParamItems() []byte {
	return keyPrefixParamItems
}
func KeyParamValue(module, name string) []byte {
	return bytes.Join([][]byte{
		keyPrefixParamValue,
		[]byte(module + "/" + name),
	}, KeyDelimiter)
}

func GetGovernParamValue(module, name string, blockNumber uint64, blockHash common.Hash) (string, error) {
	paramValue, err := FindGovernParamValue(module, name, blockHash)
	if err != nil {
		return "", err
	}
	if paramValue == nil {
		return "", common.InternalError
	} else {
		if blockNumber >= paramValue.ActiveBlock {
			return paramValue.Value, nil
		} else {
			return paramValue.StaleValue, nil
		}
	}
}

type ParamVerifier func(blockNumber uint64, blockHash common.Hash, value string) error

type GovernParam struct {
	ParamItem     *ParamItem
	ParamValue    *ParamValue
	ParamVerifier ParamVerifier
}

type ParamItem struct {
	Module string `json:"Module"`
	Name   string `json:"Name"`
	Desc   string `json:"Desc"`
}

type ParamValue struct {
	StaleValue  string `json:"StaleValue"`
	Value       string `json:"Value"`
	ActiveBlock uint64 `json:"ActiveBlock"`
}

var paramVerifier = func(blockNumber uint64, blockHash common.Hash, value string) error {
	return nil
}

func GovernStakeThreshold(blockNumber uint64, blockHash common.Hash) (*big.Int, error) {

	thresholdStr, err := GetGovernParamValue(ModuleStaking, KeyStakeThreshold, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to GovernStakeThreshold, query governParams is failed", "err", err)
		return new(big.Int).SetInt64(0), err
	}

	threshold, ok := new(big.Int).SetString(thresholdStr, 10)
	if !ok {
		return new(big.Int).SetInt64(0), fmt.Errorf("Failed to parse the govern stakethreshold")
	}

	return threshold, nil
}

func GovernOperatingThreshold(blockNumber uint64, blockHash common.Hash) (*big.Int, error) {

	thresholdStr, err := GetGovernParamValue(ModuleStaking, KeyOperatingThreshold, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to GovernOperatingThreshold, query governParams is failed", "err", err)
		return new(big.Int).SetInt64(0), err
	}

	threshold, ok := new(big.Int).SetString(thresholdStr, 10)
	if !ok {
		return new(big.Int).SetInt64(0), fmt.Errorf("Failed to parse the govern operatingthreshold")
	}

	return threshold, nil
}

func GovernMaxValidators(blockNumber uint64, blockHash common.Hash) (uint64, error) {
	maxvalidatorsStr, err := GetGovernParamValue(ModuleStaking, KeyMaxValidators, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to CheckOperatingThreshold, query governParams is failed", "err", err)
		return 0, err
	}

	maxvalidators, err := strconv.Atoi(maxvalidatorsStr)
	if nil != err {
		return 0, err
	}

	return uint64(maxvalidators), nil
}

func GovernUnStakeFreezeDuration(blockNumber uint64, blockHash common.Hash) (uint64, error) {
	durationStr, err := GetGovernParamValue(ModuleStaking, KeyUnStakeFreezeDuration, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to GovernUnStakeFreezeDuration, query governParams is failed", "err", err)
		return 0, err
	}

	duration, err := strconv.Atoi(durationStr)
	if nil != err {
		return 0, err
	}

	return uint64(duration), nil
}

func GovernSlashFractionDuplicateSign(blockNumber uint64, blockHash common.Hash) (uint32, error) {
	fractionStr, err := GetGovernParamValue(ModuleSlashing, KeySlashFractionDuplicateSign, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to GovernSlashFractionDuplicateSign, query governParams is failed", "err", err)
		return 0, err
	}

	fraction, err := strconv.Atoi(fractionStr)
	if nil != err {
		return 0, err
	}

	return uint32(fraction), nil
}

func GovernDuplicateSignReportReward(blockNumber uint64, blockHash common.Hash) (uint32, error) {
	rewardStr, err := GetGovernParamValue(ModuleSlashing, KeyDuplicateSignReportReward, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to GovernDuplicateSignReportReward, query governParams is failed", "err", err)
		return 0, err
	}

	reward, err := strconv.Atoi(rewardStr)
	if nil != err {
		return 0, err
	}

	return uint32(reward), nil
}

func GovernMaxEvidenceAge(blockNumber uint64, blockHash common.Hash) (uint32, error) {
	ageStr, err := GetGovernParamValue(ModuleSlashing, KeyMaxEvidenceAge, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to GovernMaxEvidenceAge, query governParams is failed", "err", err)
		return 0, err
	}

	age, err := strconv.Atoi(ageStr)
	if nil != err {
		return 0, err
	}

	return uint32(age), nil
}

func GovernSlashBlocksReward(blockNumber uint64, blockHash common.Hash) (uint32, error) {
	rewardStr, err := GetGovernParamValue(ModuleSlashing, KeySlashBlocksReward, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to GovernSlashBlocksReward, query governParams is failed", "err", err)
		return 0, err
	}

	reward, err := strconv.Atoi(rewardStr)
	if nil != err {
		return 0, err
	}

	return uint32(reward), nil
}

func GovernMaxBlockGasLimit(blockNumber uint64, blockHash common.Hash) (int, error) {
	gasLimitStr, err := GetGovernParamValue(ModuleBlock, KeyMaxBlockGasLimit, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to GovernMaxBlockGasLimit, query governParams is failed", "err", err)
		return 0, err
	}

	gasLimit, err := strconv.Atoi(gasLimitStr)
	if nil != err {
		return 0, err
	}

	return gasLimit, nil
}

func GovernMaxTxDataLimit(blockNumber uint64, blockHash common.Hash) (int, error) {
	sizeStr, err := GetGovernParamValue(ModuleTxPool, KeyMaxTxDataLimit, blockNumber, blockHash)
	if nil != err {
		log.Error("Failed to GovernMaxTxDataLimit, query governParams is failed", "err", err)
		return 0, err
	}

	size, err := strconv.Atoi(sizeStr)
	if nil != err {
		return 0, err
	}

	return size, nil
}
