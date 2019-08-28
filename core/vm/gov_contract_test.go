package vm

import (
	"encoding/json"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	commonvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var (
	govPlugin   *plugin.GovPlugin
	gc          *GovContract
	versionSign common.VersionSign
	chandler    *xcom.CryptoHandler
)

func init() {
	chandler = xcom.GetCryptoHandler()
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

func buildSubmitVersionInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2001))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode("verionPIPID"))
	input = append(input, common.MustRlpEncode(promoteVersion)) //new version : 1.1.1
	input = append(input, common.MustRlpEncode(uint64(5)))

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
	input = append(input, common.MustRlpEncode(uint64(4)))
	input = append(input, common.MustRlpEncode(txHashArr[2]))
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

func buildVoteInput(nodeIdx, txIdx int) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2003)))       // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[nodeIdx])) // param 1 ...
	input = append(input, common.MustRlpEncode(txHashArr[txIdx]))
	input = append(input, common.MustRlpEncode(uint8(1)))
	input = append(input, common.MustRlpEncode(promoteVersion))
	input = append(input, common.MustRlpEncode(versionSign))

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

func buildGetProposalInput(idx int) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2100)))   // func type code
	input = append(input, common.MustRlpEncode(txHashArr[idx])) // param 1 ...

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

var successExpected = hexutil.Encode(common.MustRlpEncode(xcom.Result{true, "", ""}))

func buildBlock2() {
	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, int(2+1))
	sndb.NewBlock(blockNumber2, blockHash, blockHash2)
	context := Context{
		BlockNumber: blockNumber2,
		BlockHash:   blockHash2,
	}
	gc.Evm.Context = context
}

func buildBlock3() {
	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[3], blockHash2, int(3+1))
	sndb.NewBlock(blockNumber3, blockHash2, blockHash3)
	context := Context{
		BlockNumber: blockNumber3,
		BlockHash:   blockHash3,
	}
	gc.Evm.Context = context
}

func setup(t *testing.T) {
	t.Log("setup()......")
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	precompiledContract := PlatONPrecompiledContracts[commonvm.GovContractAddr]
	gc, _ = precompiledContract.(*GovContract)
	gc.Contract = newContract(common.Big0, sender)

	state, genesis, _ := newChainState()
	gc.Evm = newEvm(blockNumber, blockHash, state)

	newPlugins()
	govPlugin = plugin.GovPluginInstance()
	gc.Plugin = govPlugin

	build_staking_data(genesis.Hash())
}

func clear(t *testing.T) {
	t.Log("tear down()......")
	sndb.Clear()
}

func TestGovContract_SubmitText(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildSubmitTextInput(), t)
}

func TestGovContract_GetTextProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()
	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	//submit a proposal and get it, this tx hash = txHashArr[2]
	runGovContract(gc, buildSubmitTextInput(), t)

	// get the Proposal by txHashArr[2]
	runGovContract(gc, buildGetProposalInput(2), t)
}

func TestGovContract_SubmitText_Sender_wrong(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	gc.Contract.CallerAddress = anotherSender

	runGovContract(gc, buildSubmitTextInput(), t, gov.TxSenderDifferFromStaking)
}

func TestGovContract_SubmitText_PIPID_empty(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildSubmitText(nodeIdArr[1], ""), t, gov.PIPIDEmpty)
}

func TestGovContract_SubmitText_PIPID_exist(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildSubmitText(nodeIdArr[1], "pipid1"), t)

	runGovContract(gc, buildSubmitText(nodeIdArr[1], "pipid1"), t, gov.PIPIDExist)
}

func TestGovContract_SubmitText_Proposal_Empty(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildSubmitText(discover.ZeroNodeID, "pipid1"), t, gov.ProposerEmpty)
}

func TestGovContract_SubmitVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildSubmitVersionInput(), t)
}

func TestGovContract_GetVersionProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildGetProposalInput(2), t)
}

func TestGovContract_SubmitVersion_NewVersionError(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", uint32(32), 5), t, gov.NewVersionError)
}

func TestGovContract_SubmitVersion_EndVotingRoundsTooSmall(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", promoteVersion, 0), t, gov.EndVotingRoundsTooSmall)
}

func TestGovContract_SubmitVersion_EndVotingRoundsTooLarge(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	//the default rounds is 6 for developer test net
	runGovContract(gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", promoteVersion, 7), t, gov.EndVotingRoundsTooLarge)
}

func TestGovContract_Float(t *testing.T) {
	t.Log(int(math.Ceil(0.667 * 1000)))
	t.Log(int(math.Floor(0.5 * 1000)))
}

func TestGovContract_DeclareVersion_VotingStage_NotVoted_DeclareActiveVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)
	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))

	runGovContract(gc, buildDeclare(nodeIdArr[0], initProgramVersion, sign), t)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 0 {
		t.Log("in this case, Gov will notify Stakging immediately, so, there's no active node list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_DeclareVersion_VotingStage_NotVoted_DeclareNewVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])
	runGovContract(gc, buildDeclareInput(), t)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("in this case, Gov will save the declared node, and notify Stakging if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_DeclareVersion_VotingStage_NotVoted_DeclareOtherVersion_Error(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	otherVersion := uint32(1<<16 | 3<<8 | 0)
	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(otherVersion))

	runGovContract(gc, buildDeclare(nodeIdArr[0], otherVersion, sign), t, gov.DeclareVersionError)

}

func TestGovContract_DeclareVersion_VotingStage_Voted_DeclareNewVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	//vote new version
	runGovContract(gc, buildVoteInput(0, 2), t)

	/*var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))*/

	//declare new version
	runGovContract(gc, buildDeclare(nodeIdArr[0], promoteVersion, versionSign), t)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("voted, Gov will save the declared node, and notify Stakging if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}

}

func TestGovContract_DeclareVersion_VotingStage_Voted_DeclareActiveVersion_ERROR(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	//vote new version
	runGovContract(gc, buildVoteInput(0, 2), t)

	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))

	//declare old version
	runGovContract(gc, buildDeclare(nodeIdArr[0], initProgramVersion, sign), t, gov.DeclareVersionError)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("voted, Gov will save the declared node, and notify Stakging if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_DeclareVersion_VotingStage_Voted_DeclareOtherVersion_ERROR(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	//vote new version
	runGovContract(gc, buildVoteInput(0, 2), t)

	otherVersion := uint32(1<<16 | 3<<8 | 0)
	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(otherVersion))

	//declare other version
	runGovContract(gc, buildDeclare(nodeIdArr[0], otherVersion, sign), t, gov.DeclareVersionError)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("voted, Gov will save the declared node, and notify Stakging if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_SubmitCancel(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	buildBlock3()
	runGovContract(gc, buildSubmitCancelInput(), t)
}

func TestGovContract_SubmitCancel_EndVotingRounds_TooLarge(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	buildBlock3()
	//the version proposal's endVotingRounds=5
	runGovContract(gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", 5, txHashArr[2]), t, gov.EndVotingRoundsTooLarge)
}

func TestGovContract_SubmitCancel_EndVotingRounds_TobeCanceledNotExist(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	buildBlock3()
	//the version proposal's endVotingRounds=5
	runGovContract(gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", 4, txHashArr[1]), t, gov.TobeCanceledProposalNotFound)
}

func TestGovContract_SubmitCancel_EndVotingRounds_TobeCanceledNotVersionProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	buildBlock2()
	runGovContract(gc, buildSubmitTextInput(), t)

	buildBlock3()

	runGovContract(gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", 4, txHashArr[2]), t, gov.TobeCanceledProposalTypeError)
}

func TestGovContract_SubmitCancel_EndVotingRounds_TobeCanceledNotAtVotingStage(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	gov.MoveVotingProposalIDToEnd(blockHash2, txHashArr[2])

	buildBlock3()

	runGovContract(gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", 4, txHashArr[2]), t, gov.TobeCanceledProposalNotAtVoting)
}

func TestGovContract_GetCancelProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	buildBlock3()
	//submit a proposal and get it.
	runGovContract(gc, buildSubmitCancelInput(), t)

	runGovContract(gc, buildGetProposalInput(2), t)
}

func TestGovContract_OneNodeVoteVersionProposal(t *testing.T) {
	setup(t)
	defer clear(t)
	//submit a proposal and vote for it.
	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildVoteInput(0, 2), t)
}

func TestGovContract_AllNodeVoteVersionProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	//submit a proposal and vote for it.
	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	for i := 0; i < 3; i++ {
		runGovContract(gc, buildVoteInput(i, 2), t)
	}
}

func TestGovContract_ListProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildSubmitVersionInput(), t)
	runGovContract(gc, buildSubmitTextInput(), t)

	runGovContract(gc, buildListProposalInput(), t)

}

func TestGovContract_GetActiveVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildGetActiveVersionInput(), t)
}

func TestGovContract_GetProgramVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildGetProgramVersionInput(), t)
}

func runGovContract(contract *GovContract, buf []byte, t *testing.T, expectedErrors ...error) {
	res, err := contract.Run(buf)

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	if expectedErrors != nil {
		assert.Equal(t, false, r.Status)
		var expected = false
		for _, expectedError := range expectedErrors {
			expected = expected || strings.Contains(r.ErrMsg, expectedError.Error())
		}
		assert.True(t, true, expected)
		t.Log("the staking result Msg:", r.ErrMsg)
	} else {
		assert.Equal(t, true, r.Status)
		t.Log("the staking result:", r)
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
