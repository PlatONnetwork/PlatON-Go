package vm

import (
	"encoding/hex"
	"encoding/json"

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
		2002: gc.submitParam,
		2003: gc.vote,
		2004: gc.declareVersion,

		// Get
		2100: gc.getProposal,
		2101: gc.getTallyResult,
		2102: gc.listProposal,
		2103: gc.getActiveVersion,
		2104: gc.getProgramVersion,
		2105: gc.listParam,
	}
}

func (gc *GovContract) submitText(verifier discover.NodeID, url string, endVotingBlock uint64) ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call submitText of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]),
		"endVotingBlock", endVotingBlock)

	if !gc.Contract.UseGas(params.SubmitTextProposalGas) {
		return nil, ErrOutOfGas
	}

	p := gov.TextProposal{
		///GithubID: githubID,
		//Topic:          topic,
		//Desc:           desc,
		Url:            url,
		ProposalType:   gov.Text,
		EndVotingBlock: endVotingBlock,
		SubmitBlock:    gc.Evm.BlockNumber.Uint64(),
		ProposalID:     gc.Evm.StateDB.TxHash(),
		Proposer:       verifier,
	}
	err := gc.Plugin.Submit(from, p, gc.Evm.BlockHash, gc.Evm.BlockNumber.Uint64(), gc.Evm.StateDB)
	return gc.errHandler("submitText", SubmitTextEvent, err, SubmitTextProposalErrorMsg)
}

func (gc *GovContract) submitVersion(verifier discover.NodeID, url string, newVersion uint32, endVotingBlock, activeBlock uint64) ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call submitVersion of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]),
		"newVersion", newVersion,
		"endVotingBlock", endVotingBlock,
		"activeBlock", activeBlock)

	if !gc.Contract.UseGas(params.SubmitVersionProposalGas) {
		return nil, ErrOutOfGas
	}

	p := gov.VersionProposal{
		//GithubID:       githubID,
		//Topic:          topic,
		//Desc:           desc,
		Url:            url,
		ProposalType:   gov.Version,
		EndVotingBlock: endVotingBlock,
		SubmitBlock:    gc.Evm.BlockNumber.Uint64(),
		ProposalID:     gc.Evm.StateDB.TxHash(),
		Proposer:       verifier,
		NewVersion:     newVersion,
		ActiveBlock:    activeBlock,
	}
	err := gc.Plugin.Submit(from, p, gc.Evm.BlockHash, gc.Evm.BlockNumber.Uint64(), gc.Evm.StateDB)
	return gc.errHandler("submitVersion", SubmitVersionEvent, err, SubmitVersionProposalErrorMsg)
}

func (gc *GovContract) submitParam(verifier discover.NodeID, url string, paramName string, currentValue, newValue string, endVotingBlock uint64) ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call submitVersion of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]),
		"endVotingBlock", endVotingBlock,
		"ParamName", paramName,
		"CurrentValue", currentValue,
		"NewValue", newValue)

	if !gc.Contract.UseGas(params.SubmitParamProposalGas) {
		return nil, ErrOutOfGas
	}

	txHash := gc.Evm.StateDB.TxHash()
	if txHash == common.ZeroHash {
		log.Warn("Call submitParam current txHash is empty!!")
		return nil, nil
	}

	p := gov.ParamProposal{
		//GithubID:       githubID,
		//Topic:          topic,
		//Desc:           desc,
		Url:            url,
		ProposalType:   gov.Param,
		EndVotingBlock: endVotingBlock,
		SubmitBlock:    gc.Evm.BlockNumber.Uint64(),
		ProposalID:     gc.Evm.StateDB.TxHash(),
		Proposer:       verifier,
		ParamName:      paramName,
		CurrentValue:   currentValue,
		NewValue:       newValue,
	}
	err := gc.Plugin.Submit(from, p, gc.Evm.BlockHash, gc.Evm.BlockNumber.Uint64(), gc.Evm.StateDB)
	return gc.errHandler("submitParam", SubmitParamEvent, err, SubmitParamProposalErrorMsg)
}

func (gc *GovContract) vote(verifier discover.NodeID, proposalID common.Hash, op uint8, programVersion uint32, programVersionSign common.VersionSign) ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call vote of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]),
		"option", op,
		"programVersion", programVersion,
		"programVersionSign", programVersionSign)

	if !gc.Contract.UseGas(params.VoteGas) {
		return nil, ErrOutOfGas
	}

	txHash := gc.Evm.StateDB.TxHash()
	if txHash == common.ZeroHash {
		log.Warn("Call vote current txHash is empty!!")
		return nil, nil
	}

	option := gov.ParseVoteOption(op)

	v := gov.Vote{}
	v.ProposalID = proposalID
	v.VoteNodeID = verifier
	v.VoteOption = option

	err := gc.Plugin.Vote(from, v, gc.Evm.BlockHash, gc.Evm.BlockNumber.Uint64(), programVersion, programVersionSign, gc.Evm.StateDB)

	return gc.errHandler("vote", VoteEvent, err, VoteErrorMsg)
}

func (gc *GovContract) declareVersion(activeNode discover.NodeID, programVersion uint32, programVersionSign common.VersionSign) ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call declareVersion of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"activeNode", hex.EncodeToString(activeNode.Bytes()[:8]),
		"programVersion", programVersion)

	if !gc.Contract.UseGas(params.DeclareVersionGas) {
		return nil, ErrOutOfGas
	}

	txHash := gc.Evm.StateDB.TxHash()
	if txHash == common.ZeroHash {
		log.Warn("Call declareVersion current txHash is empty!!")
		return nil, nil
	}

	err := gc.Plugin.DeclareVersion(from, activeNode, programVersion, programVersionSign, gc.Evm.BlockHash, gc.Evm.BlockNumber.Uint64(), gc.Evm.StateDB)

	return gc.errHandler("declareVersion", DeclareEvent, err, DeclareErrorMsg)
}

func (gc *GovContract) getProposal(proposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call getProposal of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"proposalID", proposalID)

	proposal, err := gc.Plugin.GetProposal(proposalID, gc.Evm.StateDB)

	return gc.returnHandler("getProposal", proposal, err, GetProposalErrorMsg)
}

func (gc *GovContract) getTallyResult(proposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call getTallyResult of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"proposalID", proposalID)

	rallyResult, err := gc.Plugin.GetTallyResult(proposalID, gc.Evm.StateDB)

	return gc.returnHandler("getTallyResult", rallyResult, err, GetTallyResultErrorMsg)
}

func (gc *GovContract) listProposal() ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call listProposal of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	proposalList, err := gc.Plugin.ListProposal(gc.Evm.BlockHash, gc.Evm.StateDB)

	return gc.returnHandler("listProposal", proposalList, err, ListProposalErrorMsg)
}

func (gc *GovContract) getActiveVersion() ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call getActiveVersion of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	activeVersion := gc.Plugin.GetCurrentActiveVersion(gc.Evm.StateDB)

	return gc.returnHandler("getActiveVersion", activeVersion, nil, GetActiveVersionErrorMsg)
}

func (gc *GovContract) getProgramVersion() ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call getProgramVersion of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	versionValue, err := gc.Plugin.GetProgramVersion()

	return gc.returnHandler("getProgramVersion", versionValue, err, GetProgramVersionErrorMsg)
}

func (gc *GovContract) listParam() ([]byte, error) {
	from := gc.Contract.CallerAddress
	log.Debug("Call listParam of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	paramList, err := gc.Plugin.ListParam(gc.Evm.StateDB)

	return gc.returnHandler("listParam", paramList, err, ListParamErrorMsg)
}

func (gc *GovContract) errHandler(funcName string, event string, err error, errorMsg string) ([]byte, error) {
	if err != nil {
		log.Error("Process GovContract failed.", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "errMsg", err)
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", errorMsg + ":" + err.Error()}
			resultBytes, _ := json.Marshal(res)
			xcom.AddLog(gc.Evm.StateDB, gc.Evm.BlockNumber.Uint64(), vm.GovContractAddr, event, string(resultBytes))
		} else {
			return nil, err
		}
	}
	log.Debug("Process GovContract success.", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash())
	res := xcom.Result{true, "", ""}
	resultBytes, _ := json.Marshal(res)
	xcom.AddLog(gc.Evm.StateDB, gc.Evm.BlockNumber.Uint64(), vm.GovContractAddr, event, string(resultBytes))
	return resultBytes, nil
}

func (gc *GovContract) returnHandler(funcName string, resultValue interface{}, err error, errorMsg string) ([]byte, error) {
	if nil != err {
		log.Debug("Process GovContract failed", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "errMsg", err)
		res := xcom.Result{false, "", errorMsg + ":" + err.Error()}
		resultBytes, _ := json.Marshal(res)
		return resultBytes, nil
	}
	jsonByte, err := json.Marshal(resultValue)
	if nil != err {
		log.Debug("Process GovContract failed", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash(), "errMsg", err)
		res := xcom.Result{false, "", err.Error()}
		resultBytes, _ := json.Marshal(res)
		return resultBytes, nil
	}
	log.Debug("Process GovContract success", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", "returnValue", string(jsonByte))
	res := xcom.Result{true, string(jsonByte), ""}
	resultBytes, _ := json.Marshal(res)
	return resultBytes, nil
}

/*func rlpEncodeProposalList(proposals []gov.Proposal) []byte{
	if len(proposals) == 0 {
		return nil
	}
	var data []byte
	for _, p := range proposals {
		eachData := rlpEncodeProposal(p)
		data = bytes.Join([][]byte{data, eachData}, Delimiter)
	}
	return data
}*/

/*func rlpEncodeProposalList(proposals []gov.Proposal) [][]byte {
	if len(proposals) == 0 {
		return nil
	}
	var data [][]byte
	for _, p := range proposals {
		eachData := rlpEncodeProposal(p)
		data = append(data, eachData)
	}
	return data
}

func rlpEncodeProposal(proposal gov.Proposal) []byte {
	if proposal == nil {
		return nil
	}

	var data []byte
	if proposal.GetProposalType() == gov.Text {
		txt, _ := proposal.(gov.TextProposal)
		encode := common.MustRlpEncode(txt)

		data = []byte{byte(gov.Text)}
		data = bytes.Join([][]byte{data, encode}, Delimiter)
	} else if proposal.GetProposalType() == gov.Version {
		version, _ := proposal.(gov.VersionProposal)
		encode := common.MustRlpEncode(version)

		data = []byte{byte(gov.Version)}
		data = bytes.Join([][]byte{data, encode}, Delimiter)
	}
	return data
}*/
