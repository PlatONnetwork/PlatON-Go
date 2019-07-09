package vm

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
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
	return exec_platon_contract(input, sc.FnSigns())
}

func (sc *slashingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		3000: sc.ReportMutiSign,
		3001: sc.CheckMutiSign,
	}
}


// Report the double signing behavior of the node
func (sc *slashingContract) ReportMutiSign(data string) ([]byte, error) {
	if err := sc.plugin.Slash(data, sc.Evm.BlockHash, sc.Evm.BlockNumber.Uint64(), sc.Evm.StateDB); nil != err {
		return nil, err
	}
	return nil, nil
}

// Check if the node has double sign behavior at a certain block height
func (sc *slashingContract) CheckMutiSign(etype int32, addr common.Address, blockNumber uint64) ([]byte, error) {
	if success, txHash, _ := sc.plugin.CheckMutiSign(addr, blockNumber, etype, sc.Evm.StateDB); success {
		return txHash, nil
	}
	return nil, nil
}
