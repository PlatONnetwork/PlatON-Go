package vm

import (
	"encoding/hex"
	"strconv"

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
	return exec_platon_contract(input, sc.FnSigns())
}

func (sc *SlashingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		ReportDuplicateSignEvent: sc.ReportDuplicateSign,
		// Get
		CheckDuplicateSignEvent: sc.CheckDuplicateSign,
	}
}

// Report the double signing behavior of the node
func (sc *SlashingContract) ReportDuplicateSign(data string) ([]byte, error) {

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
	evidences, err := sc.Plugin.DecodeEvidence(data)
	if nil != err {
		log.Error("slashingContract DecodeEvidence fail", "data", data, "err", err)
		return sc.buildResult(ReportDuplicateSignEvent, "ReportDuplicateSign", "", false), err
	}
	if len(evidences) == 0 {
		log.Error("slashing failed decodeEvidence len 0", "blockNumber", sc.Evm.BlockNumber.Uint64(), "blockHash", hex.EncodeToString(sc.Evm.BlockHash.Bytes()), "data", data)
		return sc.buildResult(ReportDuplicateSignEvent, "ReportDuplicateSign", "", false), common.NewBizError("evidences is nil")
	}
	if !sc.Contract.UseGas(params.DuplicateEvidencesGas * uint64(len(evidences))) {
		return nil, ErrOutOfGas
	}
	if txHash == common.ZeroHash {
		log.Warn("Call ReportDuplicateSign current txHash is empty!!")
		return nil, nil
	}
	if err := sc.Plugin.Slash(evidences, sc.Evm.BlockHash, sc.Evm.BlockNumber.Uint64(), sc.Evm.StateDB, sender); nil != err {
		if _, ok := err.(*common.BizError); ok {
			return sc.buildResult(ReportDuplicateSignEvent, "ReportDuplicateSign", "", false), err
		} else {
			return xcom.FailResult("", "fail"), err
		}
	}
	return sc.buildResult(ReportDuplicateSignEvent, "ReportDuplicateSign", "", true), nil
}

// Check if the node has double sign behavior at a certain block height
func (sc *SlashingContract) CheckDuplicateSign(etype uint32, addr common.Address, blockNumber uint64) ([]byte, error) {
	txHash, err := sc.Plugin.CheckDuplicateSign(addr, blockNumber, etype, sc.Evm.StateDB)
	data := ""
	if nil != err {
		return sc.buildResult(CheckDuplicateSignEvent, "CheckDuplicateSign", data, false), err
	}
	if len(txHash) > 0 {
		data = hexutil.Encode(txHash)
	}
	return sc.buildResult(CheckDuplicateSignEvent, "CheckDuplicateSign", data, true), nil
}

func (sc *SlashingContract) buildResult(eventType int, callFn, data string, success bool) []byte {
	var result []byte = nil
	if success {
		result = xcom.SuccessResult(data, "success")
	} else {
		result = xcom.FailResult(data, "fail")
	}
	blockNumber := sc.Evm.BlockNumber.Uint64()
	xcom.AddLog(sc.Evm.StateDB, blockNumber, vm.SlashingContractAddr, strconv.Itoa(eventType), string(result))
	log.Info("flaged to "+callFn+" of slashingContract", "txHash", sc.Evm.StateDB.TxHash().Hex(),
		"blockNumber", blockNumber, "json: ", string(result))
	return result
}
