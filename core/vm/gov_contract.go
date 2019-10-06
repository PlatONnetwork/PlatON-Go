package vm

import (
	"encoding/json"
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
	SubmitText            = uint16(2000)
	SubmitVersion         = uint16(2001)
	Vote                  = uint16(2003)
	Declare               = uint16(2004)
	SubmitCancel          = uint16(2005)
	GetProposal           = uint16(2100)
	GetResult             = uint16(2101)
	ListProposal          = uint16(2102)
	GetActiveVersion      = uint16(2103)
	GetAccuVerifiersCount = uint16(2105)
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
	return execPlatonContract(input, gc.FnSigns())
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
		GetProposal:           gc.getProposal,
		GetResult:             gc.getTallyResult,
		ListProposal:          gc.listProposal,
		GetActiveVersion:      gc.getActiveVersion,
		GetAccuVerifiersCount: gc.getAccuVerifiersCount,
	}
}

func (gc *GovContract) CheckGasPrice(gasPrice *big.Int, fcode uint16) error {
	switch fcode {
	case SubmitText:
		if gasPrice.Cmp(params.SubmitTextProposalGasPrice) < 0 {
			return common.InvalidParameter.Wrap("Gas price under the min gas price.")
		}
	case SubmitVersion:
		if gasPrice.Cmp(params.SubmitVersionProposalGasPrice) < 0 {
			return common.InvalidParameter.Wrap("Gas price under the min gas price.")
		}
	case SubmitCancel:
		if gasPrice.Cmp(params.SubmitCancelProposalGasPrice) < 0 {
			return common.InvalidParameter.Wrap("Gas price under the min gas price.")
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
	return gc.nonCallHandler("submitText", SubmitText, err)
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
	return gc.nonCallHandler("submitVersion", SubmitVersion, err)
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
	return gc.nonCallHandler("submitCancel", SubmitCancel, err)
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

	return gc.nonCallHandler("vote", Vote, err)
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

	return gc.nonCallHandler("declareVersion", Declare, err)
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

	return gc.callHandler("getProposal", proposal, err)
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
		err = gov.TallyResultNotFound
	}
	return gc.callHandler("getTallyResult", tallyResult, err)
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

	return gc.callHandler("listProposal", proposalList, err)
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

	return gc.callHandler("getActiveVersion", activeVersion, nil)
}

func (gc *GovContract) getAccuVerifiersCount(proposalID, blockHash common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress
	blockNumber := gc.Evm.BlockNumber.Uint64()
	//blockHash := gc.Evm.BlockHash
	txHash := gc.Evm.StateDB.TxHash()
	log.Debug("call getAccuVerifiesCount of GovContract",
		"from", from.Hex(),
		"txHash", txHash,
		"blockNumber", blockNumber,
		"blockHash", blockHash,
		"proposalID", proposalID)

	proposal, err := gov.GetProposal(proposalID, gc.Evm.StateDB)
	if err != nil {
		return gc.callHandler("getAccuVerifiesCount", nil, common.InternalError.Wrap(err.Error()))
	} else if proposal == nil {
		return gc.callHandler("getAccuVerifiesCount", nil, gov.ProposalNotFound)
	}

	list, err := gov.ListAccuVerifier(blockHash, proposalID)
	if err != nil {
		return gc.callHandler("getAccuVerifiesCount", nil, common.InternalError.Wrap(err.Error()))
	}

	yeas, nays, abstentions, err := gov.TallyVoteValue(proposalID, gc.Evm.StateDB)
	if err != nil {
		return gc.callHandler("getAccuVerifiesCount", nil, common.InternalError.Wrap(err.Error()))
	}

	returnValue := []uint16{uint16(len(list)), yeas, nays, abstentions}
	return gc.callHandler("getAccuVerifiesCount", returnValue, nil)
}

func (gc *GovContract) nonCallHandler(funcName string, fcode uint16, err error) ([]byte, error) {
	var event = strconv.Itoa(int(fcode))
	if err != nil {
		if _, ok := err.(*common.BizError); ok {
			resultBytes := xcom.NewFailResult(err)
			xcom.AddLog(gc.Evm.StateDB, gc.Evm.BlockNumber.Uint64(), vm.GovContractAddr, event, string(resultBytes))
			log.Warn("Execute GovContract failed.(Business error)", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "result", string(resultBytes))
			return resultBytes, nil
		} else {
			log.Error("Execute GovContract failed.(System error)", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "err", err)
			return nil, err
		}
	} else {
		log.Debug("Execute GovContract success.", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash())
		xcom.AddLog(gc.Evm.StateDB, gc.Evm.BlockNumber.Uint64(), vm.GovContractAddr, event, string(xcom.OkResultByte))
		return xcom.OkResultByte, nil
	}
}

func (gc *GovContract) callHandler(funcName string, resultValue interface{}, err error) ([]byte, error) {
	if nil != err {
		log.Error("call GovContract failed", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "err", err)
		resultBytes := xcom.NewFailResult(err)
		return resultBytes, nil
	}
	jsonByte, e := json.Marshal(resultValue)
	if nil != e {
		log.Debug("call GovContract failed", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "err", err)
		resultBytes := xcom.NewFailResult(e)
		return resultBytes, nil
	} else {
		log.Debug("call GovContract success", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "returnValue", string(jsonByte))
		resultBytes := xcom.NewSuccessResult(string(jsonByte))
		return resultBytes, nil
	}
}
