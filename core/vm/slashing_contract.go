package vm

import (
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

type SlashingContract struct {
	Plugin   *plugin.SlashingPlugin
	Contract *Contract
	Evm      *EVM
}

func (sc *SlashingContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (sc *SlashingContract) Run(input []byte) ([]byte, error) {
	return exec_platon_contract(input, sc.FnSigns())
}

func (sc *SlashingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		3000: sc.ReportMutiSign,
		3001: sc.CheckMutiSign,
	}
}

// Report the double signing behavior of the node
func (sc *SlashingContract) ReportMutiSign(data string) ([]byte, error) {
	sender := sc.Contract.CallerAddress
	if err := sc.Plugin.Slash(data, sc.Evm.BlockHash, sc.Evm.BlockNumber.Uint64(), sc.Evm.StateDB, sender); nil != err {
		return xcom.FailResult("", "failed"), err
	}
	return xcom.SuccessResult("", ""), nil
}

// Check if the node has double sign behavior at a certain block height
func (sc *SlashingContract) CheckMutiSign(etype uint32, addr common.Address, blockNumber uint64) ([]byte, error) {
	txHash, err := sc.Plugin.CheckMutiSign(addr, blockNumber, etype, sc.Evm.StateDB)
	data := ""
	if nil != err {
		return xcom.FailResult("", "failed"), err
	}
	if len(txHash) > 0 {
		data = hexutil.Encode(txHash)
	}
	return xcom.SuccessResult(data, ""), nil
}
