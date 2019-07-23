package vm

import (
	"encoding/hex"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"

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
	return params.SlashingGas
}

func (sc *SlashingContract) Run(input []byte) ([]byte, error) {
	return exec_platon_contract(input, sc.FnSigns())
}

func (sc *SlashingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		3000: sc.ReportDuplicateSign,
		// Get
		3001: sc.CheckDuplicateSign,
	}
}

// Report the double signing behavior of the node
func (sc *SlashingContract) ReportDuplicateSign(data string) ([]byte, error) {

	if !sc.Contract.UseGas(params.ReportDuplicateSignGas) {
		return nil, ErrOutOfGas
	}

	sender := sc.Contract.CallerAddress

	evidences, err := sc.Plugin.DecodeEvidence(data)
	if nil != err {
		log.Error("slashingContract DecodeEvidence fail", "data", data, "err", err)
		return xcom.FailResult("", "failed"), err
	}
	if len(evidences) == 0 {
		log.Error("slashing failed decodeEvidence len 0", "blockNumber", sc.Evm.BlockNumber.Uint64(), "blockHash", hex.EncodeToString(sc.Evm.BlockHash.Bytes()), "data", data)
		return xcom.FailResult("", "failed"), common.NewBizError("evidences is nil")
	}
	if !sc.Contract.UseGas(params.DuplicateEvidencesGas * uint64(len(evidences))) {
		return nil, ErrOutOfGas
	}
	if err := sc.Plugin.Slash(evidences, sc.Evm.BlockHash, sc.Evm.BlockNumber.Uint64(), sc.Evm.StateDB, sender); nil != err {
		return xcom.FailResult("", "failed"), err
	}
	return xcom.SuccessResult("", ""), nil
}

// Check if the node has double sign behavior at a certain block height
func (sc *SlashingContract) CheckDuplicateSign(etype uint32, addr common.Address, blockNumber uint64) ([]byte, error) {
	txHash, err := sc.Plugin.CheckDuplicateSign(addr, blockNumber, etype, sc.Evm.StateDB)
	data := ""
	if nil != err {
		return xcom.FailResult("", "failed"), err
	}
	if len(txHash) > 0 {
		data = hexutil.Encode(txHash)
	}
	return xcom.SuccessResult(data, ""), nil
}
