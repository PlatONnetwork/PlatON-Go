package vm

import (
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"reflect"
)

type slashingContract struct {
	plugin		*plugin.SlashingPlugin
	Contract 	*Contract
	Evm  		*EVM
}

func (sc *slashingContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (sc *slashingContract) Run(input []byte) ([]byte, error) {
	return sc.execute(input)
}

func (sc *slashingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		3000: sc.ReportMutiSign,
		3001: sc.CheckMutiSign,
	}
}

func (sc *slashingContract) execute(input []byte) ([]byte, error) {
	// verify the tx data by contracts method
	var fn, params, err = plugin.Verify_tx_data(input, sc.FnSigns())
	if nil != err {
		return nil, err
	}

	// execute contracts function
	result := reflect.ValueOf(fn).Call(params)
	err, ok := result[1].Interface().(error)
	if !ok {
		return result[0].Bytes(), nil
	} else {
		return nil, err
	}
}

// Report the double signing behavior of the node
func (sc *slashingContract) ReportMutiSign(mutiSignType uint8, evidence xcom.Evidence) ([]byte, error) {
	if err := sc.plugin.Slash(mutiSignType, evidence, sc.Evm.StateDB); nil != err {

	}
	return nil, nil
}

// Check if the node has double sign behavior at a certain block height
func (sc *slashingContract) CheckMutiSign(nodeId discover.NodeID, blockNumber uint64) ([]byte, error) {

	return nil, nil
}
