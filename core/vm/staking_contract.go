package vm

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"reflect"
)




type stakingContract struct {
	plugin 		*plugin.StakingPlugin
	Contract 	*Contract
	Evm      	*EVM
}


var (
	DecodeStakingTxDataErr = errors.New("decode staking tx data is err")
)




func (stkc *stakingContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (stkc *stakingContract) Run(input []byte) ([]byte, error) {
	return stkc.execute(input)
}

func (stkc *stakingContract) FnSigns () map[uint16]interface{} {
	return map[uint16]interface{}{
		1000: stkc.CreateCandidate,
		1001: stkc.EditorCandidate,
	}
}


func (stkc *stakingContract) execute (input []byte) (ret []byte, err error) {

	// verify the tx data by contracts method
	fn, params, err := plugin.Verify_tx_data(input, stkc.FnSigns())
	if nil != err {
		return nil, err
	}

	// execute contracts method
	result := reflect.ValueOf(fn).Call(params)
	if _, err := result[1].Interface().(error); !err {
		return result[0].Bytes(), nil
	}
	return nil, nil
}



func (stkc *stakingContract) CreateCandidate(name string){

}


func (stkc *stakingContract) EditorCandidate () {

}