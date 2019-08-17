package vm

import (
	"encoding/json"
	"errors"
	"math/big"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

const (
	SubmitTextProposalErrorMsg    = "Submit a text proposal error"
	SubmitVersionProposalErrorMsg = "Submit a version proposal error"
	SubmitParamProposalErrorMsg   = "Submit a param proposal error"
	SubmitCancelProposalErrorMsg  = "Submit a cancel proposal error"
	VoteErrorMsg                  = "Vote error"
	DeclareErrorMsg               = "Declare version error"
	GetProposalErrorMsg           = "Find a specified proposal error"
	GetTallyResultErrorMsg        = "Find a specified proposal's tally result error"
	ListProposalErrorMsg          = "List all proposals error"
	GetActiveVersionErrorMsg      = "Get active version error"
	GetProgramVersionErrorMsg     = "Get program version error"
	ListParamErrorMsg             = "List all parameters and values"
)

const (
	SubmitTextEvent        = "2000"
	SubmitVersionEvent     = "2001"
	SubmitParamEvent       = "2002"
	VoteEvent              = "2003"
	DeclareEvent           = "2004"
	SubmitCancelEvent      = "2005"
	GetProposalEvent       = "2100"
	GetResultEvent         = "2101"
	ListProposalEvent      = "2102"
	GetActiveVersionEvent  = "2103"
	GetProgramVersionEvent = "2104"
)

var (
	Delimiter = []byte("")
)

type GovContract struct {
	Plugin   *plugin.GovPlugin
	Contract *Contract
	Evm      *EVM
}

func (gc *GovContract) RequiredGas(input []byte) uint64 {
	return params.GovGas
}

func (gc *GovContract) Run(input []byte) ([]byte, error) {
	return exec_platon_contract(input, gc.FnSigns())
}

func (gc *GovContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		2000: gc.submitText,
		2001: gc.submitVersion,
		2003: gc.vote,
		2004: gc.declareVersion,
		2005: gc.submitCancel,

		// Get
		2100: gc.getProposal,
		2101: gc.getTallyResult,
		2102: gc.listProposal,
		2103: gc.getActiveVersion,
		2104: gc.getProgramVersion,
	}
}

func (gc *GovContract) CheckGasPrice(gasPrice *big.Int, api interface{}) error {
	switch reflect.TypeOf(api).Name() {
	case "submitText":
		if gasPrice.Cmp(params.SubmitTextProposalGasPrice) < 0 {
			return errors.New("Gas price under the min gas price.")
		}
	case "submitVersion":
		if gasPrice.Cmp(params.SubmitVersionProposalGasPrice) < 0 {
			return errors.New("Gas price under the min gas price.")
		}
	case "submitCancel":
		if gasPrice.Cmp(params.SubmitCancelProposalGasPrice) < 0 {
			return errors.New("Gas price under the min gas price.")
		}
	}
	return nil
}

func (gc *GovContract) submitText(verifier discover.NodeID, pipID string) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()

	log.Debug("Call submitText of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"PIPID", pipID,
		"verifierID", verifier.TerminalString())

	if txHash == common.ZeroHash {
		log.Warn("current txHash is empty!!")
		return nil, nil
	}

	if !gc.Contract.UseGas(params.SubmitTextProposalGas) {
		return nil, ErrOutOfGas
	}
	p := &gov.TextProposal{
		PIPID:        pipID,
		ProposalType: gov.Text,
		SubmitBlock:  blockNumber,
		ProposalID:   txHash,
		Proposer:     verifier,
	}
	err := gov.Submit(from, p, blockHash, blockNumber, plugin.StakingInstance(), gc.Evm.StateDB)
	return gc.errHandler("submitText", SubmitTextEvent, err, SubmitTextProposalErrorMsg)
}

func (gc *GovContract) submitVersion(verifier discover.NodeID, pipID string, newVersion uint32, endVotingRounds uint64) ([]byte, error) {
	from := gc.Contract.CallerAddress

	blockNumber := gc.Evm.BlockNumber.Uint64()
	blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()

	log.Debug("Call submitVersion of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"PIPID", pipID,
		"verifierID", verifier.TerminalString(),
		"newVersion", newVersion,
		"endVotingRounds", endVotingRounds)

	if txHash == common.ZeroHash {
		log.Warn("current txHash is empty!!")
		return nil, nil
	}

	if !gc.Contract.UseGas(params.SubmitVersionProposalGas) {
		return nil, ErrOutOfGas
	}

	p := &gov.VersionProposal{
		PIPID:           pipID,
		ProposalType:    gov.Version,
		EndVotingRounds: endVotingRounds,
		SubmitBlock:     blockNumber,
		ProposalID:      txHash,
		Proposer:        verifier,
		NewVersion:      newVersion,
	}
	err := gov.Submit(from, p, blockHash, blockNumber, plugin.StakingInstance(), gc.Evm.StateDB)
	return gc.errHandler("submitVersion", SubmitVersionEvent, err, SubmitVersionProposalErrorMsg)
}

func (gc *GovContract) submitCancel(verifier discover.NodeID, pipID string, endVotingRounds uint64, tobeCanceledProposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress

	blockNumber := gc.Evm.BlockNumber.Uint64()
	blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()

	log.Debug("Call submitCancel of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"PIPID", pipID,
		"verifierID", verifier.TerminalString(),
		"endVotingRounds", endVotingRounds,
		"tobeCanceled", tobeCanceledProposalID)

	if txHash == common.ZeroHash {
		log.Warn("current txHash is empty!!")
		return nil, nil
	}

	if !gc.Contract.UseGas(params.SubmitCancelProposalGas) {
		return nil, ErrOutOfGas
	}

	p := &gov.CancelProposal{
		PIPID:        pipID,
		ProposalType: gov.Cancel,
		SubmitBlock:  blockNumber,
		ProposalID:   txHash,
		Proposer:     verifier,
		TobeCanceled: tobeCanceledProposalID,
	}
	err := gov.Submit(from, p, blockHash, blockNumber, plugin.StakingInstance(), gc.Evm.StateDB)
	return gc.errHandler("submitCancel", SubmitCancelEvent, err, SubmitCancelProposalErrorMsg)
}

func (gc *GovContract) vote(verifier discover.NodeID, proposalID common.Hash, op uint8, programVersion uint32, programVersionSign common.VersionSign) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()

	log.Debug("Call vote of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"verifierID", verifier.TerminalString(),
		"option", op,
		"programVersion", programVersion,
		"programVersionSign", programVersionSign)

	if txHash == common.ZeroHash {
		log.Warn("current txHash is empty!!")
		return nil, nil
	}

	if !gc.Contract.UseGas(params.VoteGas) {
		return nil, ErrOutOfGas
	}

	option := gov.ParseVoteOption(op)

	v := gov.VoteInfo{}
	v.ProposalID = proposalID
	v.VoteNodeID = verifier
	v.VoteOption = option

	err := gov.Vote(from, v, blockHash, blockNumber, programVersion, programVersionSign, plugin.StakingInstance(), gc.Evm.StateDB)

	return gc.errHandler("vote", VoteEvent, err, VoteErrorMsg)
}

func (gc *GovContract) declareVersion(activeNode discover.NodeID, programVersion uint32, programVersionSign common.VersionSign) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("Call declareVersion of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"activeNode", activeNode.TerminalString(),
		"programVersion", programVersion)

	if txHash == common.ZeroHash {
		log.Warn("current txHash is empty!!")
		return nil, nil
	}

	if !gc.Contract.UseGas(params.DeclareVersionGas) {
		return nil, ErrOutOfGas
	}

	err := gov.DeclareVersion(from, activeNode, programVersion, programVersionSign, blockHash, blockNumber, plugin.StakingInstance(), gc.Evm.StateDB)

	return gc.errHandler("declareVersion", DeclareEvent, err, DeclareErrorMsg)
}

func (gc *GovContract) getProposal(proposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	//blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("Call getProposal of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"proposalID", proposalID)

	proposal, err := gov.GetProposal(proposalID, gc.Evm.StateDB)

	return gc.returnHandler("getProposal", proposal, err, GetProposalErrorMsg)
}

func (gc *GovContract) getTallyResult(proposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	//blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("Call getTallyResult of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"proposalID", proposalID)

	rallyResult, err := gov.GetTallyResult(proposalID, gc.Evm.StateDB)

	return gc.returnHandler("getTallyResult", rallyResult, err, GetTallyResultErrorMsg)
}

func (gc *GovContract) listProposal() ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	//blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("Call listProposal of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber)

	proposalList, err := gov.ListProposal(gc.Evm.BlockHash, gc.Evm.StateDB)

	return gc.returnHandler("listProposal", proposalList, err, ListProposalErrorMsg)
}

func (gc *GovContract) getActiveVersion() ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	//blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("Call getActiveVersion of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber)

	activeVersion := gov.GetCurrentActiveVersion(gc.Evm.StateDB)

	return gc.returnHandler("getActiveVersion", activeVersion, nil, GetActiveVersionErrorMsg)
}

func (gc *GovContract) getProgramVersion() ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	//blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("Call getProgramVersion of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber)

	versionValue, err := gov.GetProgramVersion()

	return gc.returnHandler("getProgramVersion", versionValue, err, GetProgramVersionErrorMsg)
}

func (gc *GovContract) errHandler(funcName string, event string, err error, errorMsg string) ([]byte, error) {
	if err != nil {
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", errorMsg + ":" + err.Error()}
			resultBytes, _ := json.Marshal(res)
			xcom.AddLog(gc.Evm.StateDB, gc.Evm.BlockNumber.Uint64(), vm.GovContractAddr, event, string(resultBytes))
			log.Warn("Execute GovContract failed.(Business error)", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "result", string(resultBytes))
			return resultBytes, nil
		} else {
			log.Error("Execute GovContract failed.(System error)", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "err", err)
			return nil, err
		}
	}
	log.Debug("Execute GovContract success.", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash())
	res := xcom.Result{true, "", ""}
	resultBytes, _ := json.Marshal(res)
	xcom.AddLog(gc.Evm.StateDB, gc.Evm.BlockNumber.Uint64(), vm.GovContractAddr, event, string(resultBytes))
	return resultBytes, nil
}

func (gc *GovContract) returnHandler(funcName string, resultValue interface{}, err error, errorMsg string) ([]byte, error) {
	if nil != err {
		log.Error("Call GovContract failed", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "err", err)
		res := xcom.Result{false, "", errorMsg + ":" + err.Error()}
		resultBytes, _ := json.Marshal(res)
		return resultBytes, nil
	}
	jsonByte, err := json.Marshal(resultValue)
	if nil != err {
		log.Debug("Call GovContract failed", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "err", err)
		res := xcom.Result{false, "", err.Error()}
		resultBytes, _ := json.Marshal(res)
		return resultBytes, nil
	}
	log.Debug("Call GovContract success", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", "returnValue", string(jsonByte))
	res := xcom.Result{true, string(jsonByte), ""}
	resultBytes, _ := json.Marshal(res)
	return resultBytes, nil
}
