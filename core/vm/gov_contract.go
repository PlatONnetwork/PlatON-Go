package vm

import (
	"encoding/hex"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"reflect"
)

const (
	SubtmitProposalError        = "submit proposal error"
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
	Contract   *Contract
	Evm        *EVM
}

func (stkc *govContract) RequiredGas(input []byte) uint64 {
	return 0
}

func (stkc *govContract) Run(input []byte) ([]byte, error) {
	return stkc.execute(input)
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


	gov.GovInstance().Submit(from, p, gc.Evm.BlockHash, gc.Evm.StateDB)
	return nil, nil
}

func (gc *govContract) submitVersion(verifier discover.NodeID, githubID, topic, desc, url string, newVersion uint, endVotingBlock, activeBlock uint64) ([]byte, error) {
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


	gov.GovInstance().Submit(from, p, gc.Evm.BlockHash, gc.Evm.StateDB)
	return nil, nil
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


	gov.GovInstance().Vote(from, v, gc.Evm.BlockHash, gc.Evm.StateDB)
	return nil, nil
}

func (gc *govContract) declareVersion(activeNode discover.NodeID, version uint) ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call declareVersion of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64(),
		"activeNode", hex.EncodeToString(activeNode.Bytes()[:8]))

	gov.GovInstance().DeclareVersion(from, &activeNode, version, gc.Evm.BlockHash)
	return nil, nil
}

func (gc *govContract) getProposal(proposalID common.Hash) ([]byte, error) {

	from := gc.Contract.CallerAddress

	log.Info("Call getProposal of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	gov.GovInstance().GetProposal(proposalID, gc.Evm.StateDB)
	return nil, nil
}

func (gc *govContract) getTallyResult(proposalID common.Hash) ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call getTallyResult of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	gov.GovInstance().GetTallyResult(proposalID, gc.Evm.StateDB)
	return nil, nil
}

func (gc *govContract) listProposal() ([]byte, error) {
	from := gc.Contract.CallerAddress

	log.Info("Call listProposal of govContract",
		"from", from.Hex(),
		"txHash", gc.Evm.StateDB.TxHash(),
		"blockNumber", gc.Evm.BlockNumber.Uint64())

	gov.GovInstance().ListProposal(gc.Evm.BlockHash, gc.Evm.StateDB)

	return nil, nil
}


