package vm

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
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


func (gc *govContract) submitText(funcType uint16, verifier discover.NodeID, githubID, topic, desc, url string, endVotingBlock uint64) ([]byte, error) {
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

	from := gc.Contract.CallerAddress

	gov.GovInstance().Submit(from, p, gc.Evm.StateDB)
	return nil, nil
}

func (gc *govContract) submitVersion(funcType uint16, verifier discover.NodeID, githubID, topic, desc, url string, newVersion uint, endVotingBlock, activeBlock uint64) ([]byte, error) {
	return nil, nil
}

func (gc *govContract) vote(funcType uint16, verifier discover.NodeID, proposalID common.Hash, option gov.VoteOption) ([]byte, error) {
	return nil, nil
}

func (gc *govContract) declareVersion(funcType uint16, activeNode discover.NodeID, version uint) ([]byte, error) {
	return nil, nil
}

func (gc *govContract) getProposal(funcType uint16, proposalID common.Hash) ([]byte, error) {
	return nil, nil
}

func (gc *govContract) getTallyResult(funcType uint16, proposalID common.Hash) ([]byte, error) {
	return nil, nil
}

func (gc *govContract) listProposal(funcType uint16) ([]byte, error) {
	return nil, nil
}


