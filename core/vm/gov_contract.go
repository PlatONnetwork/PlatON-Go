package vm

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"reflect"
)

const (
	SubmitTextProposalErrorMsg        	= "Submit a text proposal error"
	SubmitVersionProposalErrorMsg       = "Submit a version proposal error"
	VoteErrorMsg        				= "Vote error"
	DeclareErrorMsg        				= "Declare version error"
	GetProposalErrorMsg        			= "Find a specified proposal error"
	GetTallyResultErrorMsg				= "Find a specified proposal's tally result error"
	ListProposalErrorMsg				= "List all proposals error"
)


const (
	SubmitTextEvent   		= "2000"
	SubmitVersionEvent   	= "2001"
	VoteEvent   			= "2002"
	DeclareEvent 			= "2003"
	GetProposalEvent 		= "2100"
	GetResultEvent 			= "2101"
	ListProposalEvent 		= "2102"
)

var (
	Delimiter               = []byte("")
)

type GovContract struct {
	Plugin 	   *plugin.GovPlugin
	Contract   *Contract
	Evm        *EVM
}

func (gc *GovContract) RequiredGas(input []byte) uint64 {
	return params.GovGas
}

func (gc *GovContract) Run(input []byte) ([]byte, error) {
	return gc.execute(input)
}

func (gc *GovContract) FnSigns() map[uint16]interface{} {
	return map[uint16]interface{}{
		// Set
		2000: gc.submitText,
		2001: gc.submitVersion,
		2002: gc.vote,
		2003: gc.declareVersion,

		// Get
		2100: gc.getProposal,
		2101: gc.getTallyResult,
		2102: gc.listProposal,
	}
}

func (gc *GovContract) execute(input []byte) (ret []byte, err error) {

	// verify the tx data by contracts method
	fn, params, err := plugin.Verify_tx_data(input, gc.FnSigns())
	if nil != err {
		return nil, err
	}

	// execute contracts method
	result := reflect.ValueOf(fn).Call(params)
	if _, ok := result[1].Interface().(error); !ok {
		return result[0].Bytes(), nil
	}
	return nil, result[1].Interface().(error)
}


func (gc *GovContract) submitText(verifier discover.NodeID, githubID, topic, desc, url string, endVotingBlock uint64) ([]byte, error) {
	from := gc.Contract.CallerAddress

	fmt.Printf("endVotingBlock %d", endVotingBlock)
	log.Debug("submitText", "endVotingBlock", endVotingBlock)

	txHash := gc.Evm.StateDB.TxHash().Hex()
	log.Debug("submitText", "txHash", txHash)

	log.Info("Call submitText of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]))

	p := gov.TextProposal{
		GithubID : 			githubID,
		Topic : 			topic,
		Desc : 				desc,
		Url : 				url,
		ProposalType : 		gov.Text,
		EndVotingBlock : 	endVotingBlock,
		SubmitBlock : 		gc.Evm.BlockNumber.Uint64(),
		ProposalID : 		gc.Evm.StateDB.TxHash(),
		Proposer : 			verifier,
	}



	err := gc.Plugin.Submit(gc.Evm.BlockNumber.Uint64(), from, p, gc.Evm.BlockHash, gc.Evm.StateDB)

	return gc.errHandler("submitText", SubmitTextEvent, err, SubmitTextProposalErrorMsg)
}

func (gc *GovContract) submitVersion(verifier discover.NodeID, githubID, topic, desc, url string, newVersion uint32, endVotingBlock, activeBlock uint64) ([]byte, error) {
	from := gc.Contract.CallerAddress

	fmt.Printf("endVotingBlock %d", endVotingBlock)
	log.Debug("submitText", "endVotingBlock", endVotingBlock)

	txHash := gc.Evm.StateDB.TxHash().Hex()
	log.Debug("submitText", "txHash", txHash)

	log.Info("Call submitVersion of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]))

	p := gov.VersionProposal{
		GithubID : 			githubID,
		Topic : 			topic,
		Desc : 				desc,
		Url : 				url,
		ProposalType : 		gov.Version,
		EndVotingBlock : 	endVotingBlock,
		SubmitBlock : 		gc.Evm.BlockNumber.Uint64(),
		ProposalID : 		gc.Evm.StateDB.TxHash(),
		Proposer : 			verifier,
		NewVersion:			newVersion,
		ActiveBlock:		activeBlock,
	}



	err := gc.Plugin.Submit(gc.Evm.BlockNumber.Uint64(), from, p, gc.Evm.BlockHash, gc.Evm.StateDB)

	return gc.errHandler("submitVersion", SubmitVersionEvent, err, SubmitVersionProposalErrorMsg)
}

func (gc *GovContract) vote(verifier discover.NodeID, proposalID common.Hash, op uint8) ([]byte, error) {

	option := gov.ParseVoteOption(op)

	from := gc.Contract.CallerAddress

	log.Info("Call vote of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]))

	v := gov.Vote{}
	v.ProposalID = proposalID
	v.VoteNodeID = verifier
	v.VoteOption = option

	err := gc.Plugin.Vote(from, v, gc.Evm.BlockHash, gc.Evm.BlockNumber.Uint64(), gc.Evm.StateDB)

	return gc.errHandler("vote", VoteEvent, err, VoteErrorMsg)
}

func (gc *GovContract) declareVersion(activeNode discover.NodeID, version uint32) ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call declareVersion of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"activeNode", hex.EncodeToString(activeNode.Bytes()[:8]))

	err := gc.Plugin.DeclareVersion(from, activeNode, version, gc.Evm.BlockHash, gc.Evm.BlockNumber.Uint64(),  gc.Evm.StateDB)

	return gc.errHandler("declareVersion", DeclareEvent, err, DeclareErrorMsg)
}

func (gc *GovContract) getProposal(proposalID common.Hash) ([]byte, error) {

	from := gc.Contract.CallerAddress

	log.Info("Call getProposal of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	proposal, err := gc.Plugin.GetProposal(proposalID, gc.Evm.StateDB)

	return gc.returnHandler(proposal, err, GetProposalErrorMsg)
}



func (gc *GovContract) getTallyResult(proposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call getTallyResult of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	rallyResult, err := gc.Plugin.GetTallyResult(proposalID, gc.Evm.StateDB)

	return gc.returnHandler(rallyResult, err, GetTallyResultErrorMsg)
}

func (gc *GovContract) listProposal() ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call listProposal of GovContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	proposalList, err := gc.Plugin.ListProposal(gc.Evm.BlockHash, gc.Evm.StateDB)

	return gc.returnHandler(proposalList, err, ListProposalErrorMsg)
}


func  (gc *GovContract) errHandler(funcName string, event string, err error, errorMsg string) ([]byte, error) {
	if err != nil {
		log.Error("Call GovContract failed.", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(),  "txHash", gc.Evm.StateDB.TxHash().Hex(), "errMsg", err)
		if _, ok := err.(*common.BizError); ok {
			res := xcom.Result{false, "", errorMsg + ":" + err.Error()}
			result, _ := json.Marshal(res)
			xcom.AddLog(gc.Evm.StateDB, gc.Evm.BlockNumber.Uint64(), vm.GovContractAddr, event, string(result))
			return nil, nil
		}else {
			return nil, err
		}
	}
	res := xcom.Result{true, "", ""}
	result, _ := json.Marshal(res)
	xcom.AddLog(gc.Evm.StateDB, gc.Evm.BlockNumber.Uint64(), vm.GovContractAddr, SubmitTextEvent, string(result))

	log.Info("Call GovContract success.", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash().Hex(), "json: ", string(result))
	return common.MustRlpEncode(res), nil
}

func  (gc *GovContract) returnHandler(resultValue interface{}, err error, errorMsg string) ([]byte, error) {
	if nil != err {
		res := xcom.Result{false, "", errorMsg + ":" + err.Error()}
		return common.MustRlpEncode(res), nil
	}
	jsonByte, err := json.Marshal(resultValue)
	if nil != err {
		res := xcom.Result{false, "", err.Error()}
		return common.MustRlpEncode(res), nil
	}
	res := xcom.Result{true, string(jsonByte), ""}
	return common.MustRlpEncode(res), nil
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

func rlpEncodeProposalList(proposals []gov.Proposal) [][]byte{
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

func rlpEncodeProposal(proposal gov.Proposal) []byte{
	if proposal == nil {
		return nil
	}

	var data []byte
	if proposal.GetProposalType() == gov.Text {
		txt, _ := proposal.(gov.TextProposal)
		encode := common.MustRlpEncode(txt)

		data = []byte{byte(gov.Text)}
		data = bytes.Join([][]byte{data, encode}, Delimiter)
	}else if proposal.GetProposalType() == gov.Version{
		version, _ := proposal.(gov.VersionProposal)
		encode := common.MustRlpEncode(version)

		data = []byte{byte(gov.Version)}
		data = bytes.Join([][]byte{data, encode}, Delimiter)
	}
	return data
}