// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package vm

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

const (
	TxReportDuplicateSign = 3000
	CheckDuplicateSign    = 3001
)

type SlashingContract struct {
	Plugin   *plugin.SlashingPlugin
	Contract *Contract
	Evm      *EVM
}

func (sc *SlashingContract) RequiredGas(input []byte) uint64 {
	if checkInputEmpty(input) {
		return 0
	}
	return params.SlashingGas
}

func (sc *SlashingContract) Run(input []byte) ([]byte, error) {
	if checkInputEmpty(input) {
		return nil, nil
	}
	return execPlatonContract(input, sc.FnSigns())
}

func (sc *SlashingContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		TxReportDuplicateSign: sc.reportDuplicateSign,
		// Get
		CheckDuplicateSign: sc.checkDuplicateSign,
	}
}

func (sc *SlashingContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	return nil
}

// Report the double signing behavior of the node
func (sc *SlashingContract) reportDuplicateSign(dupType uint8, data string) ([]byte, error) {

	txHash := sc.Evm.StateDB.TxHash()
	blockNumber := sc.Evm.BlockNumber
	blockHash := sc.Evm.BlockHash
	from := sc.Contract.CallerAddress

	if !sc.Contract.UseGas(params.ReportDuplicateSignGas) {
		return nil, ErrOutOfGas
	}

	if !sc.Contract.UseGas(params.DuplicateEvidencesGas) {
		return nil, ErrOutOfGas
	}
	if txHash == common.ZeroHash {
		return nil, nil
	}

	log.Debug("Call reportDuplicateSign", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
		"TxHash", txHash.Hex(), "from", from.String())
	evidence, err := sc.Plugin.DecodeEvidence(consensus.EvidenceType(dupType), data)
	if nil != err {
		return txResultHandler(vm.SlashingContractAddr, sc.Evm, "reportDuplicateSign",
			common.InvalidParameter.Wrap(err.Error()).Error(),
			TxReportDuplicateSign, int(common.InvalidParameter.Code)), nil
	}
	if err := sc.Plugin.Slash(evidence, blockHash, blockNumber.Uint64(), sc.Evm.StateDB, from); nil != err {
		if bizErr, ok := err.(*common.BizError); ok {
			return txResultHandler(vm.SlashingContractAddr, sc.Evm, "reportDuplicateSign",
				bizErr.Error(), TxReportDuplicateSign, int(bizErr.Code)), nil
		} else {
			return nil, err
		}
	}
	return txResultHandler(vm.SlashingContractAddr, sc.Evm, "",
		"", TxReportDuplicateSign, int(common.NoErr.Code)), nil
}

// Check if the node has double sign behavior at a certain block height
func (sc *SlashingContract) checkDuplicateSign(dupType uint8, nodeId discover.NodeID, blockNumber uint64) ([]byte, error) {
	log.Info("checkDuplicateSign exist", "blockNumber", blockNumber, "nodeId", nodeId.TerminalString(), "dupType", dupType)
	txHash, err := sc.Plugin.CheckDuplicateSign(nodeId, blockNumber, consensus.EvidenceType(dupType), sc.Evm.StateDB)
	var data string

	if nil != err {
		return callResultHandler(sc.Evm, fmt.Sprintf("checkDuplicateSign, duplicateSignBlockNum: %d, nodeId: %s, dupType: %d",
			blockNumber, nodeId, dupType), data, common.InternalError.Wrap(err.Error())), nil
	}
	if len(txHash) > 0 {
		data = hexutil.Encode(txHash)
	}
	return callResultHandler(sc.Evm, fmt.Sprintf("checkDuplicateSign, duplicateSignBlockNum: %d, nodeId: %s, dupType: %d",
		blockNumber, nodeId, dupType), data, nil), nil
}
