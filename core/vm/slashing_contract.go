package vm

import (
	"math/big"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

const (
	ReportDuplicateSignEvent = 3000
	CheckDuplicateSignEvent  = 3001
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
	return execPlatonContract(input, sc.FnSigns())
}

func (sc *SlashingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		ReportDuplicateSignEvent: sc.ReportDuplicateSign,
		// Get
		CheckDuplicateSignEvent: sc.CheckDuplicateSign,
	}
}

func (sc *SlashingContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	return nil
}

// Report the double signing behavior of the node
func (sc *SlashingContract) ReportDuplicateSign(dupType uint8, data string) ([]byte, error) {

	txHash := sc.Evm.StateDB.TxHash()
	blockNumber := sc.Evm.BlockNumber
	blockHash := sc.Evm.BlockHash

	from := sc.Contract.CallerAddress

	log.Info("Call ReportDuplicateSign", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"TxHash", txHash.Hex(), "from", from.Hex())

	if !sc.Contract.UseGas(params.ReportDuplicateSignGas) {
		return nil, ErrOutOfGas
	}

	sender := sc.Contract.CallerAddress
	evidence, err := sc.Plugin.DecodeEvidence(consensus.EvidenceType(dupType), data)

	if nil != err {
		log.Error("slashingContract DecodeEvidence fail", "data", data, "err", err)
		return sc.buildResult(ReportDuplicateSignEvent, "ReportDuplicateSign", "", false, common.InvalidParameter.Wrap(err.Error())), nil
	}
	if !sc.Contract.UseGas(params.DuplicateEvidencesGas) {
		return nil, ErrOutOfGas
	}
	if txHash == common.ZeroHash {
		log.Warn("Call ReportDuplicateSign current txHash is empty!!")
		return nil, nil
	}
	if err := sc.Plugin.Slash(evidence, sc.Evm.BlockHash, sc.Evm.BlockNumber.Uint64(), sc.Evm.StateDB, sender); nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			return sc.buildResult(ReportDuplicateSignEvent, "ReportDuplicateSign", "", false, bizErr), nil
		} else {
			return nil, err
		}
	}
	return sc.buildResult(ReportDuplicateSignEvent, "ReportDuplicateSign", "", true, nil), nil
}

// Check if the node has double sign behavior at a certain block height
func (sc *SlashingContract) CheckDuplicateSign(dupType uint8, addr common.Address, blockNumber uint64) ([]byte, error) {
	txHash, err := sc.Plugin.CheckDuplicateSign(addr, blockNumber, consensus.EvidenceType(dupType), sc.Evm.StateDB)
	data := ""
	if nil != err {
		return sc.buildResult(CheckDuplicateSignEvent, "CheckDuplicateSign", data, false, common.InternalError.Wrap(err.Error())), nil
	}
	if len(txHash) > 0 {
		data = hexutil.Encode(txHash)
	}
	return sc.buildResult(CheckDuplicateSignEvent, "CheckDuplicateSign", data, true, nil), nil
}

func (sc *SlashingContract) buildResult(eventType int, callFn, data string, success bool, err *common.BizError) []byte {
	var result []byte = nil
	if success {
		result = xcom.SuccessResult(data)
	} else {
		result = xcom.FailResult(data, err)
	}
	blockNumber := sc.Evm.BlockNumber.Uint64()
	xcom.AddLog(sc.Evm.StateDB, blockNumber, vm.SlashingContractAddr, strconv.Itoa(eventType), string(result))
	log.Info("flaged to "+callFn+" of slashingContract", "txHash", sc.Evm.StateDB.TxHash().Hex(),
		"blockNumber", blockNumber, "json: ", string(result))
	return result
}
