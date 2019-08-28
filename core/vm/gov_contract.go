package vm

import (
	"encoding/json"
	"errors"
	"math/big"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

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
	GetProposalErrorMsg           = "Find proposal error"
	GetTallyResultErrorMsg        = "Find tally result error"
	ListProposalErrorMsg          = "List all proposals error"
	GetActiveVersionErrorMsg      = "Get active version error"
	GetProgramVersionErrorMsg     = "Get program version error"
	ListParamErrorMsg             = "List all parameters and values"
)

const (
	SubmitText        = uint16(2000)
	SubmitVersion     = uint16(2001)
	Vote              = uint16(2003)
	Declare           = uint16(2004)
	SubmitCancel      = uint16(2005)
	GetProposal       = uint16(2100)
	GetResult         = uint16(2101)
	ListProposal      = uint16(2102)
	GetActiveVersion  = uint16(2103)
	GetProgramVersion = uint16(2104)
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
		SubmitText:    gc.submitText,
		SubmitVersion: gc.submitVersion,
		Vote:          gc.vote,
		Declare:       gc.declareVersion,
		SubmitCancel:  gc.submitCancel,

		// Get
		GetProposal:       gc.getProposal,
		GetResult:         gc.getTallyResult,
		ListProposal:      gc.listProposal,
		GetActiveVersion:  gc.getActiveVersion,
		GetProgramVersion: gc.getProgramVersion,
	}
}

func (gc *GovContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	switch fcode {
	case SubmitText:
		if gasPrice.Cmp(params.SubmitTextProposalGasPrice) < 0 {
			return errors.New("Gas price under the min gas price.")
		}
	case SubmitVersion:
		if gasPrice.Cmp(params.SubmitVersionProposalGasPrice) < 0 {
			return errors.New("Gas price under the min gas price.")
		}
	case SubmitCancel:
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

	log.Debug("call submitText of GovContract",
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
	return gc.errHandler("submitText", SubmitText, err, SubmitTextProposalErrorMsg)
}

func (gc *GovContract) submitVersion(verifier discover.NodeID, pipID string, newVersion uint32, endVotingRounds uint64) ([]byte, error) {
	from := gc.Contract.CallerAddress

	blockNumber := gc.Evm.BlockNumber.Uint64()
	blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()

	log.Debug("call submitVersion of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"PIPID", pipID,
		"verifierID", verifier.TerminalString(),
		"newVersion", newVersion,
		"newVersionString", xutil.ProgramVersion2Str(newVersion),
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
	return gc.errHandler("submitVersion", SubmitVersion, err, SubmitVersionProposalErrorMsg)
}

func (gc *GovContract) submitCancel(verifier discover.NodeID, pipID string, endVotingRounds uint64, tobeCanceledProposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress

	blockNumber := gc.Evm.BlockNumber.Uint64()
	blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()

	log.Debug("call submitCancel of GovContract",
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
		PIPID:           pipID,
		EndVotingRounds: endVotingRounds,
		ProposalType:    gov.Cancel,
		SubmitBlock:     blockNumber,
		ProposalID:      txHash,
		Proposer:        verifier,
		TobeCanceled:    tobeCanceledProposalID,
	}
	err := gov.Submit(from, p, blockHash, blockNumber, plugin.StakingInstance(), gc.Evm.StateDB)
	return gc.errHandler("submitCancel", SubmitCancel, err, SubmitCancelProposalErrorMsg)
}

func (gc *GovContract) vote(verifier discover.NodeID, proposalID common.Hash, op uint8, programVersion uint32, programVersionSign common.VersionSign) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()

	log.Debug("call vote of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"verifierID", verifier.TerminalString(),
		"option", op,
		"programVersion", programVersion,
		"programVersionString", xutil.ProgramVersion2Str(programVersion),
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

	return gc.errHandler("vote", Vote, err, VoteErrorMsg)
}

func (gc *GovContract) declareVersion(activeNode discover.NodeID, programVersion uint32, programVersionSign common.VersionSign) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("call declareVersion of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"activeNode", activeNode.TerminalString(),
		"programVersion", programVersion,
		"programVersionString", xutil.ProgramVersion2Str(programVersion))

	if txHash == common.ZeroHash {
		log.Warn("current txHash is empty!!")
		return nil, nil
	}

	if !gc.Contract.UseGas(params.DeclareVersionGas) {
		return nil, ErrOutOfGas
	}

	err := gov.DeclareVersion(from, activeNode, programVersion, programVersionSign, blockHash, blockNumber, plugin.StakingInstance(), gc.Evm.StateDB)

	return gc.errHandler("declareVersion", Declare, err, DeclareErrorMsg)
}

func (gc *GovContract) getProposal(proposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	//blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("call getProposal of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"proposalID", proposalID)

	proposal, err := gov.GetExistProposal(proposalID, gc.Evm.StateDB)

	return gc.returnHandler("getProposal", proposal, err, GetProposalErrorMsg)
}

func (gc *GovContract) getTallyResult(proposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	//blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("call getTallyResult of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"proposalID", proposalID)

	tallyResult, err := gov.GetTallyResult(proposalID, gc.Evm.StateDB)

	if tallyResult == nil {
		err = common.BizErrorf("tally result not found")
	}
	return gc.returnHandler("getTallyResult", tallyResult, err, GetTallyResultErrorMsg)
}

func (gc *GovContract) listProposal() ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	//blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("call listProposal of GovContract",
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
	log.Debug("call getActiveVersion of GovContract",
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
	log.Debug("call getProgramVersion of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber)

	versionValue, err := gov.GetProgramVersion()

	return gc.returnHandler("getProgramVersion", versionValue, err, GetProgramVersionErrorMsg)
}

func (gc *GovContract) errHandler(funcName string, fcode uint16, err error, errorMsg string) ([]byte, error) {
	var event = strconv.Itoa(int(fcode))
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
		log.Error("call GovContract failed", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "err", err)
		res := xcom.Result{false, "", errorMsg + ":" + err.Error()}
		resultBytes, _ := json.Marshal(res)
		return resultBytes, nil
	}
	jsonByte, err := json.Marshal(resultValue)
	if nil != err {
		log.Debug("call GovContract failed", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "err", err)
		res := xcom.Result{false, "", errorMsg + ":" + err.Error()}
		resultBytes, _ := json.Marshal(res)
		return resultBytes, nil
	}
	log.Debug("call GovContract success", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "returnValue", string(jsonByte))
	res := xcom.Result{true, string(jsonByte), ""}
	resultBytes, _ := json.Marshal(res)
	return resultBytes, nil
}
