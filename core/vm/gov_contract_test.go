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
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/log"

	//"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/node"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/common"
	commonvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var (
	govPlugin   *plugin.GovPlugin
	gc          *GovContract
	versionSign common.VersionSign
	chandler    *node.CryptoHandler

	paramModule       = gov.ModuleStaking
	paramName         = gov.KeyMaxValidators
	defaultProposalID = txHashArr[1]
)

func commit_sndb(chain *mock.Chain) {
	/*
		//Flush() signs a Hash to the current block which has no hash yet. Flush() do not write the data to database.
		//in this file, all blocks in each test case has a hash already, so, do not call Flush()
				if err := chain.SnapDB.Flush(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number); err != nil {
					fmt.Println("commit_sndb error:", err)
				}
	*/
	if err := chain.SnapDB.Commit(chain.CurrentHeader().Hash()); err != nil {
		fmt.Println("commit_sndb error:", err)
	}
}

func prepair_sndb(chain *mock.Chain, txHash common.Hash) {
	if txHash == common.ZeroHash {
		chain.AddBlock()
	} else {
		chain.AddBlockWithTxHash(txHash)
	}

	//fmt.Println("prepair_sndb::::::", chain.CurrentHeader().ParentHash.Hex())
	if err := chain.SnapDB.NewBlock(chain.CurrentHeader().Number, chain.CurrentHeader().ParentHash, chain.CurrentHeader().Hash()); err != nil {
		fmt.Println("prepair_sndb error:", err)
	}

	//prepare gc to run contract
	gc.Evm = newEvm(chain.CurrentHeader().Number, chain.CurrentHeader().Hash(), chain.StateDB)
}

func skip_emptyBlock(chain *mock.Chain, blockNumber uint64) {
	cnt := blockNumber - chain.CurrentHeader().Number.Uint64()
	for i := uint64(0); i < cnt; i++ {
		prepair_sndb(chain, common.ZeroHash)
		commit_sndb(chain)
	}
}

func init() {
	chandler = node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])
	versionSign.SetBytes(chandler.MustSign(promoteVersion))
}

func buildSubmitTextInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2000))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[1])) // param 1 ...
	input = append(input, common.MustRlpEncode("textUrl"))

	return common.MustRlpEncode(input)
}
func buildSubmitText(nodeID discover.NodeID, pipID string) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2000))) // func type code
	input = append(input, common.MustRlpEncode(nodeID))       // param 1 ...
	input = append(input, common.MustRlpEncode(pipID))

	return common.MustRlpEncode(input)
}

func buildSubmitParam(nodeID discover.NodeID, pipID string, module, name, newValue string) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2002))) // func type code
	input = append(input, common.MustRlpEncode(nodeID))       // param 1 ...
	input = append(input, common.MustRlpEncode(pipID))
	input = append(input, common.MustRlpEncode(module))
	input = append(input, common.MustRlpEncode(name))
	input = append(input, common.MustRlpEncode(newValue))

	return common.MustRlpEncode(input)
}

func buildSubmitVersionInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2001))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode("verionPIPID"))
	input = append(input, common.MustRlpEncode(promoteVersion)) //new version : 1.1.1
	input = append(input, common.MustRlpEncode(xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())))

	return common.MustRlpEncode(input)
}

func buildSubmitVersion(nodeID discover.NodeID, pipID string, newVersion uint32, endVotingRounds uint64) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2001))) // func type code
	input = append(input, common.MustRlpEncode(nodeID))       // param 1 ...
	input = append(input, common.MustRlpEncode(pipID))
	input = append(input, common.MustRlpEncode(newVersion)) //new version : 1.1.1
	input = append(input, common.MustRlpEncode(endVotingRounds))

	return common.MustRlpEncode(input)
}

func buildSubmitCancelInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2005))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ..
	input = append(input, common.MustRlpEncode("cancelPIPID"))
	input = append(input, common.MustRlpEncode(xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())-1))
	input = append(input, common.MustRlpEncode(defaultProposalID))
	return common.MustRlpEncode(input)
}

func buildSubmitCancel(nodeID discover.NodeID, pipID string, endVotingRounds uint64, tobeCanceledProposalID common.Hash) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2005))) // func type code
	input = append(input, common.MustRlpEncode(nodeID))       // param 1 ..
	input = append(input, common.MustRlpEncode(pipID))
	input = append(input, common.MustRlpEncode(endVotingRounds))
	input = append(input, common.MustRlpEncode(tobeCanceledProposalID))
	return common.MustRlpEncode(input)
}

func buildVoteInput(nodeIdx int, proposalID common.Hash) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2003)))       // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[nodeIdx])) // param 1 ...
	input = append(input, common.MustRlpEncode(proposalID))
	input = append(input, common.MustRlpEncode(uint8(1)))
	input = append(input, common.MustRlpEncode(promoteVersion))
	input = append(input, common.MustRlpEncode(versionSign))

	return common.MustRlpEncode(input)
}

func buildVote(nodeIdx int, proposalID common.Hash, option gov.VoteOption, programVersion uint32, sign common.VersionSign) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2003)))       // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[nodeIdx])) // param 1 ...
	input = append(input, common.MustRlpEncode(proposalID))
	input = append(input, common.MustRlpEncode(option))
	input = append(input, common.MustRlpEncode(programVersion))
	input = append(input, common.MustRlpEncode(sign))

	return common.MustRlpEncode(input)
}

func buildDeclareInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2004))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode(promoteVersion))
	input = append(input, common.MustRlpEncode(versionSign))
	return common.MustRlpEncode(input)
}

func buildDeclare(nodeID discover.NodeID, declaredVersion uint32, sign common.VersionSign) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2004))) // func type code
	input = append(input, common.MustRlpEncode(nodeID))       // param 1 ...
	input = append(input, common.MustRlpEncode(declaredVersion))
	input = append(input, common.MustRlpEncode(sign))
	return common.MustRlpEncode(input)
}

func buildGetProposalInput(proposalID common.Hash) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2100))) // func type code
	input = append(input, common.MustRlpEncode(proposalID))   // param 1 ...

	return common.MustRlpEncode(input)
}

func buildGetTallyResultInput(idx int) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2101))) // func type code
	input = append(input, common.MustRlpEncode(txHashArr[0])) // param 1 ...

	return common.MustRlpEncode(input)
}

func buildListProposalInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2102))) // func type code
	return common.MustRlpEncode(input)
}

func buildGetActiveVersionInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2103))) // func type code
	return common.MustRlpEncode(input)
}

func buildGetProgramVersionInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2104))) // func type code
	return common.MustRlpEncode(input)
}

func buildGetAccuVerifiersCountInput(proposalID, blockHash common.Hash) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2105))) // func type code
	input = append(input, common.MustRlpEncode(proposalID))
	input = append(input, common.MustRlpEncode(blockHash))
	return common.MustRlpEncode(input)
}

func buildListGovernParam(module string) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2106))) // func type code
	input = append(input, common.MustRlpEncode(module))
	return common.MustRlpEncode(input)
}

func buildGetGovernParamValueInput(module, name string) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2104)))
	input = append(input, common.MustRlpEncode(module))
	input = append(input, common.MustRlpEncode(name))
	return common.MustRlpEncode(input)
}

func setup(t *testing.T) *mock.Chain {
	t.Log("setup()......")
	//to turn on log's debug level
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	precompiledContract := PlatONPrecompiledContracts[commonvm.GovContractAddr]
	gc, _ = precompiledContract.(*GovContract)
	//default sender of tx, this could be changed in different test case if necessary
	gc.Contract = newContract(common.Big0, sender)

	chain, _ := newChain()
	newPlugins()
	govPlugin = plugin.GovPluginInstance()
	gc.Plugin = govPlugin
	build_staking_data_new(chain)

	if _, err := gov.InitGenesisGovernParam(common.ZeroHash, chain.SnapDB, 2048); err != nil {
		t.Error("error", err)
	}
	gov.RegisterGovernParamVerifiers()

	commit_sndb(chain)

	//the contract will retrieve this txHash as ProposalID
	prepair_sndb(chain, defaultProposalID)
	return chain
}

func clear(chain *mock.Chain, t *testing.T) {
	t.Log("tear down()......")
	if err := chain.SnapDB.Clear(); err != nil {
		t.Error("clear chain.SnapDB error", err)
	}

}

func TestGovContract_SubmitText(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(false, gc, buildSubmitTextInput(), t)
}

func TestGovContract_GetTextProposal(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal and get it, this tx hash in gc is = txHashArr[1]
	runGovContract(false, gc, buildSubmitTextInput(), t)

	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])

	// get the Proposal by txHashArr[1]
	runGovContract(true, gc, buildGetProposalInput(defaultProposalID), t)
}

func TestGovContract_SubmitText_Sender_wrong(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	gc.Contract.CallerAddress = anotherSender

	runGovContract(false, gc, buildSubmitTextInput(), t, gov.TxSenderDifferFromStaking)
}

func TestGovContract_SubmitText_PIPID_empty(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	runGovContract(false, gc, buildSubmitText(nodeIdArr[1], ""), t, gov.PIPIDEmpty)
}

func TestGovContract_SubmitText_ProposalID_exist(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	runGovContract(false, gc, buildSubmitText(nodeIdArr[1], "pipid1"), t)

	runGovContract(false, gc, buildSubmitText(nodeIdArr[1], "pipid33"), t, gov.ProposalIDExist)
}

func TestGovContract_SubmitText_Proposal_Empty(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	runGovContract(false, gc, buildSubmitText(discover.ZeroNodeID, "pipid1"), t, gov.ProposerEmpty)
}

func TestGovContract_ListGovernParam(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(true, gc, buildListGovernParam(paramModule), t)
}

func TestGovContract_ListGovernParam_all(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(true, gc, buildListGovernParam(""), t)
}

func TestGovContract_GetGovernParamValue(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(true, gc, buildGetGovernParamValueInput(paramModule, paramName), t)
}

func TestGovContract_GetGovernParamValue_NotFound(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(true, gc, buildGetGovernParamValueInput(paramModule, "stakeThreshold_Err"), t, gov.UnsupportedGovernParam)
}

func TestGovContract_SubmitParam(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	runGovContract(false, gc, buildSubmitParam(nodeIdArr[1], "pipid3", paramModule, paramName, "30"), t)

	p, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	} else {
		if p == nil {
			t.Fatal("not find proposal error")
		} else {
			pp := p.(*gov.ParamProposal)
			assert.Equal(t, "30", pp.NewValue)
		}
	}
}

func TestGovContract_SubmitParam_thenSubmitParamFailed(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	runGovContract(false, gc, buildSubmitParam(nodeIdArr[1], "pipid3", paramModule, paramName, "30"), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	runGovContract(false, gc, buildSubmitParam(nodeIdArr[2], "pipid4", paramModule, paramName, "35"), t, gov.VotingParamProposalExist)
}

func TestGovContract_SubmitParam_thenSubmitVersionFailed(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	runGovContract(false, gc, buildSubmitParam(nodeIdArr[1], "pipid3", paramModule, paramName, "30"), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	runGovContract(false, gc, buildSubmitVersionInput(), t, gov.VotingParamProposalExist)
}

func TestGovContract_SubmitParam_GetAccuVerifiers(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	value, err := gov.GetGovernParamValue(paramModule, paramName, chain.CurrentHeader().Number.Uint64(), chain.CurrentHeader().Hash())
	if err != nil {
		t.Errorf("%s", err)
	} else {
		assert.Equal(t, "25", value)
	}

	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitParam(nodeIdArr[1], "pipid3", paramModule, paramName, "30"), t)
	//runGovContract(false, gc, buildSubmitTextInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, txHashArr[1], gov.Yes)
	commit_sndb(chain)

	runGovContract(true, gc, buildGetAccuVerifiersCountInput(defaultProposalID, chain.CurrentHeader().Hash()), t)

}

func TestGovContract_voteTwoProposal_punished(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	value, err := gov.GetGovernParamValue(paramModule, paramName, chain.CurrentHeader().Number.Uint64(), chain.CurrentHeader().Hash())
	if err != nil {
		t.Errorf("%s", err)
	} else {
		assert.Equal(t, "25", value)
	}

	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitParam(nodeIdArr[1], "pipid3", paramModule, paramName, "30"), t)
	//runGovContract(false, gc, buildSubmitTextInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.Yes)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[3])
	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitCancel(nodeIdArr[1], "pipid2", 4, defaultProposalID), t)
	//runGovContract(false, gc, buildSubmitTextInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[4])
	allVote(chain, t, txHashArr[3], gov.Yes)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[5])
	punished := make(map[discover.NodeID]struct{})
	currentValidatorList, _ := plugin.StakingInstance().ListCurrentValidatorID(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64())

	// punish last one
	lastOne := currentValidatorList[len(currentValidatorList)-1]
	punished[lastOne] = struct{}{}

	gov.NotifyPunishedVerifiers(chain.CurrentHeader().Hash(), punished, chain.StateDB)

	runGovContract(true, gc, buildGetAccuVerifiersCountInput(defaultProposalID, chain.CurrentHeader().Hash()), t)

	runGovContract(true, gc, buildGetAccuVerifiersCountInput(txHashArr[3], chain.CurrentHeader().Hash()), t)

}

func TestGovContract_SubmitParam_Pass(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	value, err := gov.GetGovernParamValue(paramModule, paramName, chain.CurrentHeader().Number.Uint64(), chain.CurrentHeader().Hash())
	if err != nil {
		t.Errorf("%s", err)
	} else {
		assert.Equal(t, "25", value)
	}

	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitParam(nodeIdArr[1], "pipid3", paramModule, paramName, "30"), t)
	//runGovContract(false, gc, buildSubmitTextInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, txHashArr[1], gov.Yes)
	commit_sndb(chain)

	runGovContract(true, gc, buildGetAccuVerifiersCountInput(defaultProposalID, chain.CurrentHeader().Hash()), t)

	p, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}

	//skip empty block
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	endBlock(chain, t)
	commit_sndb(chain)

	//at the end of voting block, the status=pass;
	result, err := gov.GetTallyResult(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.Pass {
		t.Log("the result status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	} else {
		t.Fatal("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	}

	//from the next to voting block, the parameter value will be the new value
	skip_emptyBlock(chain, p.GetEndVotingBlock()+1)
	value, err = gov.GetGovernParamValue(paramModule, paramName, chain.CurrentHeader().Number.Uint64(), chain.CurrentHeader().Hash())
	if err != nil {
		t.Errorf("%s", err)
	} else {
		assert.Equal(t, "30", value)
	}
}

func TestGovContract_SubmitVersion(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	runGovContract(false, gc, buildSubmitVersionInput(), t)
}

func TestGovContract_SubmitVersion_thenSubmitParamFailed(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	runGovContract(false, gc, buildSubmitParam(nodeIdArr[1], "pipid3", paramModule, paramName, "30"), t, gov.VotingVersionProposalExist)
}

func TestGovContract_GetVersionProposal(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and get it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	runGovContract(true, gc, buildGetProposalInput(defaultProposalID), t)
}

func TestGovContract_SubmitVersion_AnotherVoting(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal
	runGovContract(false, gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", promoteVersion, xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	//submit a proposal
	runGovContract(false, gc, buildSubmitVersion(nodeIdArr[2], "versionPIPID2", promoteVersion, xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())), t, gov.VotingVersionProposalExist)
}

func TestGovContract_SubmitVersion_Passed(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it. the proposalID = txHashArr[1]
	runGovContract(false, gc, buildSubmitVersionInput(), t)

	commit_sndb(chain)

	build_staking_data_more(chain)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.Yes)
	commit_sndb(chain)

	pTemp, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}
	p := pTemp.(*gov.VersionProposal)

	//skip empty blocks
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	endBlock(chain, t)
	commit_sndb(chain)

	result, err := gov.GetTallyResult(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.PreActive {
		t.Log("the result status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	} else {
		t.Fatal("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	}

	//skip empty blocks, this version proposal is pre-active
	skip_emptyBlock(chain, p.GetActiveBlock()-1)
}

func TestGovContract_SubmitVersion_AnotherPreActive(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it. the proposalID = txHashArr[1]
	runGovContract(false, gc, buildSubmitVersionInput(), t)

	commit_sndb(chain)

	build_staking_data_more(chain)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.Yes)
	commit_sndb(chain)

	pTemp, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}
	p := pTemp.(*gov.VersionProposal)

	//skip empty blocks
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	endBlock(chain, t)
	commit_sndb(chain)

	result, err := gov.GetTallyResult(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.PreActive {
		t.Log("the result status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	} else {
		t.Fatal("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	}

	//skip empty blocks, this version proposal is pre-active
	skip_emptyBlock(chain, p.GetActiveBlock()-1)
	//submit another version proposal
	runGovContract(false, gc, buildSubmitVersion(nodeIdArr[2], "versionPIPID2", promoteVersion, xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())), t, gov.PreActiveVersionProposalExist)
}

func TestGovContract_SubmitVersion_Passed_Clear(t *testing.T) {
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it. the proposalID = txHashArr[1]
	runGovContract(false, gc, buildSubmitVersionInput(), t)

	commit_sndb(chain)

	build_staking_data_more(chain)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.Yes)
	commit_sndb(chain)

	pTemp, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}
	p := pTemp.(*gov.VersionProposal)

	//skip empty blocks
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	endBlock(chain, t)
	commit_sndb(chain)

	result, err := gov.GetTallyResult(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.PreActive {
		t.Log("the result status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	} else {
		t.Fatal("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	}

	//skip empty blocks, this version proposal is pre-active
	skip_emptyBlock(chain, p.GetActiveBlock()-1)

	prepair_sndb(chain, common.ZeroHash)

	if preactiveID, err := gov.GetPreActiveProposalID(chain.CurrentHeader().Hash()); err != nil {
		t.Error("GetPreActiveProposalID error", err)
	} else {
		assert.Equal(t, preactiveID, defaultProposalID)
	}

	/*if err := gov.MovePreActiveProposalIDToEnd(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("move pre-active proposal ID to end list failed", "proposalID", defaultProposalID, "blockHash", blockHash)
	}*/

	//----clear all data of this pre-active proposal
	if err := gov.ClearProcessingProposals(chain.CurrentHeader().Hash(), chain.StateDB); err != nil {
		t.Error("ClearProcessingProposals error", err)
	} else {
		if votinglist, err := gov.ListVotingProposalID(chain.CurrentHeader().Hash()); err != nil {
			t.Error("ListVotingProposalID, err", err)
		} else {
			assert.Equal(t, 0, len(votinglist))
		}

		if endList, err := gov.ListEndProposalID(chain.CurrentHeader().Hash()); err != nil {
			t.Error("ListEndProposalID, err", err)
		} else {
			assert.Equal(t, 1, len(endList))
		}

		//-------
		if vvList, err := gov.ListVoteValue(defaultProposalID, chain.CurrentHeader().Hash()); err != nil {
			t.Error("ListVoteValue, err", err)
		} else {
			assert.Equal(t, 0, len(vvList))
		}

		if avList, err := gov.ListAccuVerifier(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
			t.Error("ListVoteValue, err", err)
		} else {
			assert.Equal(t, 0, len(avList))
		}
	}
}

func TestGovContract_SubmitVersion_NewVersionError(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(false, gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", uint32(32), xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())), t, gov.NewVersionError)
}

func TestGovContract_SubmitVersion_EndVotingRoundsTooSmall(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(false, gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", promoteVersion, 0), t, gov.EndVotingRoundsTooSmall)
}

func TestGovContract_SubmitVersion_EndVotingRoundsTooLarge(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//the default rounds is 6 for developer test net
	runGovContract(false, gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", promoteVersion, xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())+1), t, gov.EndVotingRoundsTooLarge)
}

func TestGovContract_DeclareVersion_VotingStage_NotVoted_DeclareActiveVersion(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and get it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))

	//fmt.Println("hash:::", chain.CurrentHeader().Hash().Hex())

	runGovContract(false, gc, buildDeclare(nodeIdArr[0], initProgramVersion, sign), t)

	//if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
	//	t.Error("cannot list ActiveNode")
	//} else if len(nodeList) == 0 {
	//	t.Log("in this case, Gov will notify Staking immediately, so, there's no active node list")
	//} else {
	//	t.Fatal("cannot list ActiveNode")
	//}
}

func TestGovContract_DeclareVersion_VotingStage_NotVoted_DeclareNewVersion(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and get it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])
	runGovContract(false, gc, buildDeclareInput(), t)

	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("in this case, Gov will save the declared node, and notify Staking if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_DeclareVersion_VotingStage_NotVoted_DeclareOtherVersion_Error(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and get it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)

	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	otherVersion := uint32(1<<16 | 3<<8 | 0)
	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(otherVersion))

	runGovContract(false, gc, buildDeclare(nodeIdArr[0], otherVersion, sign), t, gov.DeclareVersionError)

}

func TestGovContract_DeclareVersion_VotingStage_Voted_DeclareNewVersion(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and get it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	//vote new version
	runGovContract(false, gc, buildVoteInput(0, defaultProposalID), t)

	//declare new version
	runGovContract(false, gc, buildDeclare(nodeIdArr[0], promoteVersion, versionSign), t)

	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("voted, Gov will save the declared node, and notify Staking if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}

}

func TestGovContract_DeclareVersion_VotingStage_Voted_DeclareOldVersion_ERROR(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and get it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	//vote new version
	runGovContract(false, gc, buildVoteInput(0, defaultProposalID), t)

	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))

	//vote new version, but declare old version
	runGovContract(false, gc, buildDeclare(nodeIdArr[0], initProgramVersion, sign), t, gov.DeclareVersionError)

	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("voted, Gov will save the declared node, and notify Staking if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_DeclareVersion_VotingStage_Voted_DeclareOtherVersion_ERROR(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and get it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])
	//vote new version
	runGovContract(false, gc, buildVoteInput(0, defaultProposalID), t)

	otherVersion := uint32(1<<16 | 3<<8 | 0)
	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(otherVersion))

	//vote new version, but declare other version
	runGovContract(false, gc, buildDeclare(nodeIdArr[0], otherVersion, sign), t, gov.DeclareVersionError)

	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("voted, Gov will save the declared node, and notify Staking if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_SubmitCancel(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	runGovContract(false, gc, buildSubmitCancelInput(), t)

}

func TestGovContract_SubmitCancel_AnotherVoting(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal
	runGovContract(false, gc, buildSubmitVersion(nodeIdArr[0], "versionPIPID", promoteVersion, xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	runGovContract(false, gc, buildSubmitCancel(nodeIdArr[1], "cancelPIPID", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())-1, defaultProposalID), t)

	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[3])
	runGovContract(false, gc, buildSubmitCancel(nodeIdArr[2], "cancelPIPIDAnother", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())-1, defaultProposalID), t, gov.VotingCancelProposalExist)
}

func TestGovContract_SubmitCancel_EndVotingRounds_TooLarge(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	runGovContract(false, gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()), defaultProposalID), t, gov.EndVotingRoundsTooLarge)
}

func TestGovContract_SubmitCancel_EndVotingRounds_TobeCanceledNotExist(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(false, gc, buildSubmitVersionInput(), t)

	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	//the version proposal's endVotingRounds=5
	runGovContract(false, gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())-1, txHashArr[3]), t, gov.TobeCanceledProposalNotFound)
}

func TestGovContract_SubmitCancel_EndVotingRounds_TobeCanceledNotVersionProposal(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//txHash = txHashArr[2] is a text proposal
	runGovContract(false, gc, buildSubmitTextInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	//try to cancel a text proposal
	runGovContract(false, gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())-1, defaultProposalID), t, gov.TobeCanceledProposalTypeError)
}

func TestGovContract_SubmitCancel_EndVotingRounds_TobeCanceledNotAtVotingStage(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	//move the proposal ID from voting-list to end-list
	err := gov.MoveVotingProposalIDToEnd(defaultProposalID, chain.CurrentHeader().Hash())
	if err != nil {
		t.Fatal("err", err)
	}
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[3])
	//try to cancel a closed version proposal
	runGovContract(false, gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())-1, defaultProposalID), t, gov.TobeCanceledProposalNotAtVoting)
}

func TestGovContract_GetCancelProposal(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	//submit a proposal and get it.
	runGovContract(false, gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds())-1, defaultProposalID), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[3])
	runGovContract(true, gc, buildGetProposalInput(txHashArr[2]), t)
}

func TestGovContract_Vote_VersionProposal(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])

	runGovContract(false, gc, buildVoteInput(0, defaultProposalID), t)

	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("cannot list ActiveNode", "err", err)
	} else if len(nodeList) == 1 {
		t.Log("voted, Gov will save the declared node, and notify Staking if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}
func TestGovContract_Vote_Duplicated(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	runGovContract(false, gc, buildVoteInput(0, defaultProposalID), t)
	runGovContract(false, gc, buildVoteInput(0, defaultProposalID), t, gov.VoteDuplicated)
	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("voted duplicated, Gov will count this node once in active node list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_Vote_OptionError(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	// vote option = 0, it's wrong
	runGovContract(false, gc, buildVote(0, defaultProposalID, 0, promoteVersion, versionSign), t, gov.VoteOptionError)

	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 0 {
		t.Log("option error, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_Vote_ProposalNotExist(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))

	//verify vote new version, but it has not upgraded
	// txIdx=4, not a proposalID
	prepair_sndb(chain, txHashArr[2])
	runGovContract(false, gc, buildVote(0, txHashArr[4], gov.Yes, initProgramVersion, sign), t, gov.ProposalNotFound)

	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("list ActiveNode error", "err", err)
	} else if len(nodeList) == 0 {
		t.Log("proposal not found, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_Vote_TextProposalPassed(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitTextInput(), t)
	commit_sndb(chain)

	build_staking_data_more(chain)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.Yes)
	commit_sndb(chain)

	p, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}

	//skip empty block
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	//tally vote
	endBlock(chain, t)
	commit_sndb(chain)

	runGovContract(false, gc, buildVote(0, defaultProposalID, gov.No, promoteVersion, versionSign), t, gov.ProposalNotAtVoting)

	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 0 {
		t.Log("option error, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_SubmitText_passed_PIPID_exist(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitText(nodeIdArr[1], "pipid1"), t)
	commit_sndb(chain)

	build_staking_data_more(chain)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.Yes)
	commit_sndb(chain)

	p, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}

	//skip empty block
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	//tally vote
	endBlock(chain, t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[3])

	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	runGovContract(false, gc, buildSubmitText(nodeIdArr[2], "pipid1"), t, gov.PIPIDExist)
}

func TestGovContract_SubmitText_voting_PIPID_exist(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitText(nodeIdArr[1], "pipid1"), t)
	commit_sndb(chain)

	build_staking_data_more(chain)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.Yes)
	commit_sndb(chain)

	//submit another proposal
	prepair_sndb(chain, txHashArr[3])
	runGovContract(false, gc, buildSubmitText(nodeIdArr[2], "pipid1"), t, gov.PIPIDExist)
}

func TestGovContract_SubmitText_NotPassed_SamePIPID_Allowed(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitText(nodeIdArr[1], "pipid1"), t)
	commit_sndb(chain)

	build_staking_data_more(chain)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.No)
	commit_sndb(chain)

	p, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}

	//skip empty block
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	//tally vote
	endBlock(chain, t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[3])
	runGovContract(false, gc, buildSubmitText(nodeIdArr[2], "pipid1"), t)

	p, err = gov.GetProposal(txHashArr[3], chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	} else {
		assert.Equal(t, nodeIdArr[2], p.GetProposer())
	}

}

func TestGovContract_Vote_VerifierNotUpgraded(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//the proposalID will be txHashArr[1]
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))

	prepair_sndb(chain, txHashArr[2])
	//verify vote new version, but it has not upgraded
	//txIdx should figure out the proposalID
	runGovContract(false, gc, buildVote(0, defaultProposalID, gov.Yes, initProgramVersion, sign), t, gov.VerifierNotUpgraded)
	commit_sndb(chain)

	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), txHashArr[1]); err != nil {
		t.Error("list ActiveNode error", "err", err)
	} else if len(nodeList) == 0 {
		t.Log("verifier not upgraded, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_Vote_ProgramVersionError(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it. proposalID= txHashArr[1]
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	otherVersion := uint32(1<<16 | 3<<8 | 0)
	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(otherVersion))

	prepair_sndb(chain, txHashArr[2])
	//verify vote new version, but it has not upgraded
	runGovContract(false, gc, buildVote(0, defaultProposalID, gov.Yes, otherVersion, sign), t, gov.VerifierNotUpgraded)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[3])
	if nodeList, err := gov.GetActiveNodeList(chain.CurrentHeader().Hash(), defaultProposalID); err != nil {
		t.Error("list ActiveNode error", "err", err)
	} else if len(nodeList) == 0 {
		t.Log("verifier program version error, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_AllNodeVoteVersionProposal(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it. proposalID= txHashArr[1]
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	chandler := node.GetCryptoHandler()
	prepair_sndb(chain, txHashArr[2])
	for i := 0; i < 3; i++ {
		chandler.SetPrivateKey(priKeyArr[i])
		var sign common.VersionSign
		sign.SetBytes(chandler.MustSign(promoteVersion))
		//verify vote new version, but it has not upgraded
		runGovContract(false, gc, buildVote(i, defaultProposalID, gov.Yes, promoteVersion, sign), t)
	}
}

func TestGovContract_TextProposal_pass(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitTextInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, txHashArr[1], gov.Yes)
	commit_sndb(chain)

	p, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}

	//skip empty block
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	endBlock(chain, t)
	commit_sndb(chain)

	result, err := gov.GetTallyResult(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.Pass {
		t.Log("the result status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	} else {
		t.Fatal("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	}
}

func TestGovContract_VersionProposal_Active(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal and vote for it. proposalID= txHashArr[1]
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.Yes)
	commit_sndb(chain)

	pTemp, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}
	p := pTemp.(*gov.VersionProposal)

	//skip empty block
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	endBlock(chain, t)

	commit_sndb(chain)

	result, err := gov.GetTallyResult(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.PreActive {
		t.Log("the result status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	} else {
		t.Fatal("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	}

	//skip empty block
	skip_emptyBlock(chain, p.GetActiveBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	beginBlock(chain, t)
	commit_sndb(chain)

	result, err = gov.GetTallyResult(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.Active {
		t.Log("the result status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	} else {
		t.Fatal("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	}
}

func TestGovContract_VersionProposal_Active_GetExtraParam_V0_11_0(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal and vote for it. proposalID= txHashArr[1]
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	allVote(chain, t, defaultProposalID, gov.Yes)
	commit_sndb(chain)

	pTemp, err := gov.GetProposal(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}
	p := pTemp.(*gov.VersionProposal)

	//skip empty block
	skip_emptyBlock(chain, p.GetEndVotingBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)
	endBlock(chain, t)

	commit_sndb(chain)

	result, err := gov.GetTallyResult(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.PreActive {
		t.Log("the result status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	} else {
		t.Fatal("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	}

	//skip empty block
	skip_emptyBlock(chain, p.GetActiveBlock()-1)

	// build_staking_data_more will build a new block base on chain.SnapDB.Current
	build_staking_data_more(chain)

	// the version proposal is not be active yet, so the new extra gov parameters do not exist
	govParam, err := gov.FindGovernParam(gov.ModuleStaking, gov.KeyZeroProduceNumberThreshold, chain.CurrentHeader().Hash())
	if err != nil {
		t.Fatal("find govern param err", err)
	}
	assert.Nil(t, govParam)

	govParam, err = gov.FindGovernParam(gov.ModuleStaking, gov.KeyZeroProduceCumulativeTime, chain.CurrentHeader().Hash())
	if err != nil {
		t.Fatal("find govern param err", err)
	}
	assert.Nil(t, govParam)

	beginBlock(chain, t)
	commit_sndb(chain)

	//prepair_sndb(chain, common.ZeroHash)

	result, err = gov.GetTallyResult(defaultProposalID, chain.StateDB)
	if err != nil {
		t.Fatal("get tally result err", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	}
	assert.Equal(t, gov.Active, result.Status)

	// the version proposal is active, so the new extra gov parameters are existing also
	govParam, err = gov.FindGovernParam(gov.ModuleSlashing, gov.KeyZeroProduceNumberThreshold, chain.CurrentHeader().Hash())
	if err != nil {
		t.Fatal("find govern param err", err)
	}
	if govParam == nil {
		t.Fatal("cannot find the extra param: slashing.zeroProduceNumberThreshold")
	} else {
		assert.Equal(t, "2", govParam.ParamValue.Value)
	}

	govParam, err = gov.FindGovernParam(gov.ModuleSlashing, gov.KeyZeroProduceCumulativeTime, chain.CurrentHeader().Hash())
	if err != nil {
		t.Fatal("find govern param err", err)
	}
	if govParam == nil {
		t.Fatal("cannot find the extra param: slashing.zeroProduceCumulativeTime")
	} else {
		assert.Equal(t, "3", govParam.ParamValue.Value)
	}

}

func TestGovContract_ListProposal(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	runGovContract(false, gc, buildSubmitTextInput(), t)
	commit_sndb(chain)

	runGovContract(true, gc, buildListProposalInput(), t)

}

func TestGovContract_GetActiveVersion(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	runGovContract(true, gc, buildGetActiveVersionInput(), t)
}

func TestGovContract_getAccuVerifiersCount(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	//get accu verifiers
	runGovContract(true, gc, buildGetAccuVerifiersCountInput(txHashArr[1], chain.CurrentHeader().Hash()), t)
}

func TestGovContract_getAccuVerifiersCount_wrongProposalID(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//submit a proposal and vote for it.
	runGovContract(false, gc, buildSubmitVersionInput(), t)
	commit_sndb(chain)

	prepair_sndb(chain, txHashArr[2])
	////get accu verifiers
	runGovContract(true, gc, buildGetAccuVerifiersCountInput(txHashArr[2], chain.CurrentHeader().Hash()), t, gov.ProposalNotFound)
}

func runGovContract(callType bool, contract *GovContract, buf []byte, t *testing.T, expectedErrors ...error) {
	res, err := contract.Run(buf)
	assert.True(t, nil == err)

	var result xcom.Result
	if callType {
		err = json.Unmarshal(res, &result)
		assert.True(t, nil == err)
	} else {
		var retCode uint32
		err = json.Unmarshal(res, &retCode)
		assert.True(t, nil == err)
		result.Code = retCode
	}

	if expectedErrors != nil {
		assert.NotEqual(t, common.OkCode, result.Code)
		var expected = false
		for _, expectedError := range expectedErrors {
			expectedCode, _ := common.DecodeError(expectedError)
			expected = expected || result.Code == expectedCode
		}
		assert.True(t, expected)
		t.Log("the expected errCode:", result.Code, "errMsg:", expectedErrors)
	} else {
		assert.Equal(t, common.OkCode, result.Code)
		t.Log("the expected resultCode:", result.Code, "resultData:", result.Ret)
	}
}

func Test_ResetVoteOption(t *testing.T) {
	v := gov.VoteInfo{}
	v.ProposalID = common.ZeroHash
	v.VoteNodeID = discover.NodeID{}
	v.VoteOption = gov.Abstention
	t.Log(v)

	v.VoteOption = gov.Yes
	t.Log(v)
}

func logResult(t *testing.T, resultValue interface{}) {
	if IsBlank(resultValue) {
		resultBytes := xcom.NewResult(common.NotFound, nil)
		t.Log("result  json：", string(resultBytes))
	} else {
		resultBytes := xcom.NewResult(nil, resultValue)
		t.Log("result  json：", string(resultBytes))
	}

}
func Test_Json_Marshal_nil(t *testing.T) {
	// slice
	var vList []gov.GovernParam
	logResult(t, vList)

	vList = []gov.GovernParam{}
	logResult(t, vList)

	// struct
	var vStruct gov.GovernParam
	logResult(t, &vStruct)

	// struct refer
	var vStructRef *gov.GovernParam
	logResult(t, vStructRef)

	// map
	var vMap map[string]gov.GovernParam
	logResult(t, vMap)

	vMap = make(map[string]gov.GovernParam)
	logResult(t, vMap)

	// string
	var vString string
	logResult(t, vString)

	var vUint32 uint32
	logResult(t, vUint32)

	var vUintList []uint32
	logResult(t, vUintList)

	var vProposal gov.Proposal
	logResult(t, vProposal)

	var str string
	str = "20"
	//jsonByte, _ := json.Marshal(str)
	resultBytes := xcom.NewResult(nil, str)
	t.Log("result string", string(resultBytes))

}

func allVote(chain *mock.Chain, t *testing.T, pid common.Hash, option gov.VoteOption) {
	//for _, nodeID := range nodeIdArr {
	currentValidatorList, _ := plugin.StakingInstance().ListCurrentValidatorID(chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64())
	voteCount := len(currentValidatorList)
	chandler := node.GetCryptoHandler()
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	for i := 0; i < voteCount; i++ {
		vote := gov.VoteInfo{
			ProposalID: pid,
			VoteNodeID: nodeIdArr[i],
			VoteOption: option,
		}

		chandler.SetPrivateKey(priKeyArr[i])
		versionSign := common.VersionSign{}
		versionSign.SetBytes(chandler.MustSign(promoteVersion))

		err := gov.Vote(sender, vote, chain.CurrentHeader().Hash(), 1, promoteVersion, versionSign, plugin.StakingInstance(), chain.StateDB)
		if err != nil {
			t.Fatalf("vote err: %s.", err)
		}
	}
}

func beginBlock(chain *mock.Chain, t *testing.T) {
	err := govPlugin.BeginBlock(chain.CurrentHeader().Hash(), chain.CurrentHeader(), chain.StateDB)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
}

func endBlock(chain *mock.Chain, t *testing.T) {
	err := govPlugin.EndBlock(chain.CurrentHeader().Hash(), chain.CurrentHeader(), chain.StateDB)
	if err != nil {
		t.Fatalf("end block err... %s", err)
	}
}
