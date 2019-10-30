package vm

import (
	"bytes"
	"encoding/hex"
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

	//log.Debug("Call ReportDuplicateSign", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
	//	"TxHash", txHash.Hex(), "from", from.Hex())

	if !sc.Contract.UseGas(params.ReportDuplicateSignGas) {
		return nil, ErrOutOfGas
	}

	if !sc.Contract.UseGas(params.DuplicateEvidencesGas) {
		return nil, ErrOutOfGas
	}
	if txHash == common.ZeroHash {
		return nil, nil
	}

	evidence, err := sc.Plugin.DecodeEvidence(consensus.EvidenceType(dupType), data)
	if nil != err {
		return sc.buildReceipt(ReportDuplicateSignEvent, "ReportDuplicateSign", false, common.InvalidParameter.Wrap(err.Error())), nil
	}

	if err := sc.Plugin.Slash(evidence, blockHash, blockNumber.Uint64(), sc.Evm.StateDB, from); nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			return sc.buildReceipt(ReportDuplicateSignEvent, "ReportDuplicateSign", false, bizErr), nil
		} else {
			return nil, err
		}
	}
	return sc.buildReceipt(ReportDuplicateSignEvent, "ReportDuplicateSign", true, nil), nil
}

// Check if the node has double sign behavior at a certain block height
func (sc *SlashingContract) CheckDuplicateSign(dupType uint8, addr common.Address, blockNumber uint64) ([]byte, error) {
	log.Info("CheckDuplicateSign exist", "blockNumber", blockNumber, "addr", hex.EncodeToString(addr.Bytes()), "dupType", dupType)
	txHash, err := sc.Plugin.CheckDuplicateSign(addr, blockNumber, consensus.EvidenceType(dupType), sc.Evm.StateDB)
	var data string

	if nil != err {
		return sc.buildResult("CheckDuplicateSign", data, false, common.InternalError.Wrap(err.Error())), nil
	}
	if bytes.Equal(txHash, common.ZeroHash.Bytes()) {
		data = hexutil.Encode(txHash)
	}
	return sc.buildResult("CheckDuplicateSign", data, true, nil), nil
}

func (sc *SlashingContract) buildReceipt(eventType int, callFn string, ok bool, err *common.BizError) []byte {
	var receipt string
	blockNumber := sc.Evm.BlockNumber.Uint64()
	if ok {
		receipt = strconv.Itoa(int(common.NoErr.Code))
	} else {
		receipt = strconv.Itoa(int(err.Code))
		log.Error("Failed to "+callFn+" of slashingContract", "txHash", sc.Evm.StateDB.TxHash().Hex(),
			"blockNumber", blockNumber, "receipt: ", receipt, "the reason", err.Msg)
	}
	xcom.AddLog(sc.Evm.StateDB, blockNumber, vm.SlashingContractAddr, strconv.Itoa(eventType), receipt)
	return []byte(receipt)
}

func (sc *SlashingContract) buildResult(callFn, data string, success bool, err *common.BizError) []byte {
	var result []byte = nil
	blockNumber := sc.Evm.BlockNumber.Uint64()
	if success {
		result = xcom.OkResult(data)
		log.Debug("Call "+callFn+" of slashingContract", "txHash", sc.Evm.StateDB.TxHash().Hex(),
			"blockNumber", blockNumber, "json: ", string(result))
	} else {
		result = xcom.FailResult(err)
		log.Error("Failed to "+callFn+" of slashingContract", "txHash", sc.Evm.StateDB.TxHash().Hex(),
			"blockNumber", blockNumber, "json: ", string(result))
	}
	return result
}
