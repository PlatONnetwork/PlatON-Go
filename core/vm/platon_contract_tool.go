package vm

import (
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

func execPlatonContract(input []byte, command map[uint16]interface{}) (ret []byte, err error) {

	// verify the tx data by contracts method
	_, fn, params, err := plugin.VerifyTxData(input, command)
	if nil != err {
		log.Error("Failed to verify contract tx", "err", err)
		return nil, err
	}

	// execute contracts method
	result := reflect.ValueOf(fn).Call(params)
	if err, ok := result[1].Interface().(error); ok {
		return nil, err
	}
	return result[0].Bytes(), nil
}
