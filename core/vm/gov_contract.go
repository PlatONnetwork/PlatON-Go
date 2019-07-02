package vm

import (
	"encoding/hex"
	"encoding/json"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
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


type govContract struct {
	plugin 	   *plugin.GovPlugin
	Contract   *Contract
	Evm        *EVM
}

func (gc *govContract) RequiredGas(input []byte) uint64 {
	return params.GovGas
}

func (gc *govContract) Run(input []byte) ([]byte, error) {
	return gc.execute(input)
}

func (gc *govContract) FnSigns() map[uint16]interface{} {
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

func (gc *govContract) execute(input []byte) (ret []byte, err error) {

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


func (gc *govContract) submitText(verifier discover.NodeID, githubID, topic, desc, url string, endVotingBlock uint64) ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call submitText of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]))

	p := gov.TextProposal{}
	p.SetGithubID(githubID)
	p.SetTopic(topic)
	p.SetDesc(desc)
	p.SetUrl(url)
	p.SetProposalType(gov.Text)

	p.SetEndVotingBlock(endVotingBlock)
	p.SetSubmitBlock(gc.Evm.BlockNumber.Uint64())
	p.SetProposalID(gc.Evm.StateDB.TxHash())
	p.SetProposer(verifier)



	err := gc.plugin.Submit(gc.Evm.BlockNumber.Uint64(), from, p, gc.Evm.BlockHash, gc.Evm.StateDB)

	return gc.errHandler("submitText", SubmitTextEvent, err, SubmitTextProposalErrorMsg)
}

func (gc *govContract) submitVersion(verifier discover.NodeID, githubID, topic, desc, url string, newVersion uint32, endVotingBlock, activeBlock uint64) ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call submitVersion of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]))

	p := gov.VersionProposal{}
	p.SetGithubID(githubID)
	p.SetTopic(topic)
	p.SetDesc(desc)
	p.SetUrl(url)
	p.SetProposalType(gov.Text)
	p.SetEndVotingBlock(endVotingBlock)
	p.SetSubmitBlock(gc.Evm.BlockNumber.Uint64())
	p.SetProposalID(gc.Evm.StateDB.TxHash())
	p.SetProposer(verifier)

	p.SetNewVersion(newVersion)
	p.SetActiveBlock(activeBlock)

	err := gc.plugin.Submit(gc.Evm.BlockNumber.Uint64(), from, p, gc.Evm.BlockHash, gc.Evm.StateDB)

	return gc.errHandler("submitVersion", SubmitVersionEvent, err, SubmitVersionProposalErrorMsg)
}

func (gc *govContract) vote(verifier discover.NodeID, proposalID common.Hash, option gov.VoteOption) ([]byte, error) {

	from := gc.Contract.CallerAddress

	log.Info("Call vote of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"verifierID", hex.EncodeToString(verifier.Bytes()[:8]))

	v := gov.Vote{}
	v.ProposalID = proposalID
	v.VoteNodeID = verifier
	v.VoteOption = option

	err := gc.plugin.Vote(from, v, gc.Evm.BlockHash, gc.Evm.BlockNumber.Uint64(), gc.Evm.StateDB)

	return gc.errHandler("vote", VoteEvent, err, VoteErrorMsg)
}

func (gc *govContract) declareVersion(activeNode discover.NodeID, version uint32) ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call declareVersion of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"activeNode", hex.EncodeToString(activeNode.Bytes()[:8]))

	err := gc.plugin.DeclareVersion(from, activeNode, version, gc.Evm.BlockHash, gc.Evm.BlockNumber.Uint64(),  gc.Evm.StateDB)

	return gc.errHandler("declareVersion", DeclareEvent, err, DeclareErrorMsg)
}

func (gc *govContract) getProposal(proposalID common.Hash) ([]byte, error) {

	from := gc.Contract.CallerAddress

	log.Info("Call getProposal of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	proposal, err := gc.plugin.GetProposal(proposalID, gc.Evm.StateDB)

	return gc.returnHandler(proposal, err, GetProposalErrorMsg)
}

func (gc *govContract) getTallyResult(proposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call getTallyResult of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	rallyResult, err := gc.plugin.GetTallyResult(proposalID, gc.Evm.StateDB)

	return gc.returnHandler(rallyResult, err, GetTallyResultErrorMsg)
}

func (gc *govContract) listProposal() ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call listProposal of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	proposalList, err := gc.plugin.ListProposal(gc.Evm.BlockHash, gc.Evm.StateDB)

	return gc.returnHandler(proposalList, err, ListProposalErrorMsg)
}


func  (gc *govContract) errHandler(funcName string, event string, err error, errorMsg string) ([]byte, error) {
	if err != nil {
		log.Error("Call govContract failed.", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(),  "txHash", gc.Evm.StateDB.TxHash().Hex(), "errMsg", err)
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

	log.Info("Call govContract success.", "method", funcName, "blockNumber", gc.Evm.BlockNumber.Uint64(), "txHash", gc.Evm.StateDB.TxHash().Hex(), "json: ", string(result))

	return nil, nil
}

func  (gc *govContract) returnHandler(resultValue interface{}, err error, errorMsg string) ([]byte, error) {
	if nil != err {
		res := xcom.Result{false, "", errorMsg + ":" + err.Error()}
		data, _ := rlp.EncodeToBytes(res)
		return data, nil
	}
	bytes, _ := json.Marshal(resultValue)
	res := xcom.Result{true, string(bytes), ""}
	data, _ := rlp.EncodeToBytes(res)
	return data, nil
}
