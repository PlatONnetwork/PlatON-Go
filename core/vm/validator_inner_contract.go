package vm

import (
	"errors"
	"fmt"

	"bytes"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"encoding/json"
)

const (
	currentValidatorKey = "current_validator"
	nextValidatorKey    = "next_validator"

	txTypeUpdate  = 2000
	txTypeCurrent = 2001
	txTypeNext    = 2002
)

type ValidateNode struct {
	Index   int            `json:"index"`
	Address common.Address `json:"address"`
}

type VNode map[discover.NodeID]*ValidateNode
type Validators struct {
	ValidatNodes     VNode  `json:"validateNodes"`
	StartTimeOfEpoch uint64 `json:"startTimeOfEpoch"`
}

type ValidatorInnerContractBase interface {
	UpdateValidators(validators *Validators) error
	CurrentValidators() (*Validators, error)
	NextValidators() (*Validators, error)
}

type validatorInnerContract struct {
	ValidatorInnerContractBase

	Contract *Contract
	Evm      *EVM
}

func (vic *validatorInnerContract) RequiredGas(input []byte) uint64 {
	return params.EcrecoverGas
}

func (vic *validatorInnerContract) Run(input []byte) ([]byte, error) {
	var cmd = map[string]interface{}{
		"UpdateValidators":  vic.UpdateValidators,
		"CurrentValidatros": vic.CurrentValidators,
		"NextValidators":    vic.NextValidators,
	}
	return vic.execute(input, cmd)
}

func (vic *validatorInnerContract) UpdateValidators(validators *Validators) error {
	if len(validators.ValidatNodes) <= 0 {
		log.Error("Empty validator nodes")
		return errors.New("Empty validator nodes")
	}

	vs, err := rlp.EncodeToBytes(validators)
	if err != nil {
		log.Error("RLP encode error", "validators", validators, "error", err)
		return err
	}
	vic.Evm.StateDB.SetState(vic.Contract.Address(), []byte(nextValidatorKey), vs)
	return nil
}

func (vic *validatorInnerContract) CurrentValidators() (*Validators, error) {
	state := vic.Evm.StateDB
	b := state.GetState(vic.Contract.Address(), []byte(currentValidatorKey))

	var vds Validators
	err := rlp.DecodeBytes(b, &vds)
	return &vds, err
}

func (vic *validatorInnerContract) NextValidators() (*Validators, error) {
	state := vic.Evm.StateDB
	b := state.GetState(vic.Contract.Address(), []byte(nextValidatorKey))

	var vds Validators
	err := rlp.DecodeBytes(b, &vds)
	return &vds, err
}

func (vic *validatorInnerContract) execute(input []byte, cmd map[string]interface{}) (ret []byte, err error) {
	defer func() {
		if er := recover(); er != nil {
			ret, err = nil, fmt.Errorf("Validator inner contract execute fail: %v", er)
			log.Error("Validator inner contract execute fail", "error", err)
		}
	}()

	var source [][]byte
	if err = rlp.Decode(bytes.NewReader(input), &source); err != nil {
		log.Error("Validator inner contract execute fail", "error", err)
		return nil, errors.New("RLP decode fail")
	}

	if len(source) < 2 {
		log.Error("Params base length not match")
		return nil, errors.New("Params base length not match")
	}

	funcName := string(source[1])
	if _, ok := cmd[funcName]; !ok {
		log.Error("Function undefined", "function", funcName)
		return nil, errors.New("Function undefined")
	}

	txType := common.BytesToInt64(source[0])
	switch (txType) {
	case txTypeUpdate:
		var vds Validators
		err = json.Unmarshal(source[2], &vds)
		if err != nil {
			log.Error("Parse params fail", "params", string(source[2]), "error", err)
			return nil, err
		}
		err = vic.UpdateValidators(&vds)
		return nil, err

	case txTypeCurrent:
		var vds *Validators = nil
		vds, err = vic.CurrentValidators()
		if err != nil {
			log.Error("Get current validators fail", "error", err)
			return nil, err
		}
		b, _ := json.Marshal(&vds)
		return b, nil

	case txTypeNext:
		var vds *Validators = nil
		vds, err = vic.NextValidators()
		if err != nil {
			log.Error("Get next validators fail", "error", err)
			return nil, err
		}
		b, _ := json.Marshal(&vds)
		return b, nil
	default:
		log.Error("Unexpected transaction type", "txType", txType)
		return nil, errors.New("unexpected transaction type")
	}
}
