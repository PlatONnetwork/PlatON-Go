package vm

import (
	"encoding/json"
	"math"
	"math/big"
	"strings"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/node"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

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
	chandler    *node.CryptoHandler
)

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

func buildSubmitVersionInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2001))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode("verionPIPID"))
	input = append(input, common.MustRlpEncode(promoteVersion)) //new version : 1.1.1
	input = append(input, common.MustRlpEncode(xcom.VersionProposalVote_ConsensusRounds()))

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
	input = append(input, common.MustRlpEncode(xcom.VersionProposalVote_ConsensusRounds()-1))
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

func buildVote(nodeIdx, txIdx int, option uint8, programVersion uint32, sign common.VersionSign) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2003)))       // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[nodeIdx])) // param 1 ...
	input = append(input, common.MustRlpEncode(txHashArr[txIdx]))
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

func buildGetAccuVerifiersCountInput(proposalID, blockHash common.Hash) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2105))) // func type code
	input = append(input, common.MustRlpEncode(proposalID))
	input = append(input, common.MustRlpEncode(blockHash))
	return common.MustRlpEncode(input)
}

var successExpected = hexutil.Encode(common.MustRlpEncode(xcom.Result{0, "", ""}))

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

	sndb = snapshotdb.Instance()
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

func TestGovContract_SubmitVersion_AnotherVoting(t *testing.T) {
	setup(t)
	defer clear(t)
	stateDB := gc.Evm.StateDB.(*mock.MockStateDB)
	stateDB.Prepare(txHashArr[0], lastBlockHash, 0)

	//submit a proposal
	runGovContract(gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", promoteVersion, xcom.VersionProposalVote_ConsensusRounds()), t)

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	stateDB.Prepare(txHashArr[1], lastBlockHash, 0)
	//submit a proposal
	runGovContract(gc, buildSubmitVersion(nodeIdArr[2], "versionPIPID2", promoteVersion, 4), t, gov.VotingVersionProposalExist)

	sndb.Commit(lastBlockHash) //commit
	sndb.Compaction()          //write to level db

}

func TestGovContract_SubmitVersion_AnotherPreActive(t *testing.T) {
	setup(t)
	defer clear(t)

	stateDB := gc.Evm.StateDB.(*mock.MockStateDB)
	stateDB.Prepare(txHashArr[0], lastBlockHash, 0)

	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	build_staking_data_more(uint64(2))

	allVote(stateDB, t, txHashArr[0])
	sndb.Commit(lastBlockHash) //commit
	sndb.Compaction()          //write to level db

	pTemp, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}
	p := pTemp.(*gov.VersionProposal)

	// res
	lastBlockNumber = uint64(p.GetEndVotingBlock() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	// build_staking_data_more will build a new block base on sndb.Current
	build_staking_data_more(p.GetEndVotingBlock())
	endBlock(stateDB, t)

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	result, err := gov.GetTallyResult(txHashArr[0], stateDB)
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

	lastBlockNumber = uint64(p.GetActiveBlock() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	stateDB.Prepare(txHashArr[1], lastBlockHash, 0)
	//submit a proposal
	runGovContract(gc, buildSubmitVersion(nodeIdArr[2], "versionPIPID2", promoteVersion, 4), t, gov.PreActiveVersionProposalExist)
}

func TestGovContract_SubmitVersion_NewVersionError(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", uint32(32), 4), t, gov.NewVersionError)
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
	runGovContract(gc, buildSubmitVersion(nodeIdArr[1], "versionPIPID", promoteVersion, 55), t, gov.EndVotingRoundsTooLarge)
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

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))

	runGovContract(gc, buildDeclare(nodeIdArr[0], initProgramVersion, sign), t)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 0 {
		t.Log("in this case, Gov will notify Staking immediately, so, there's no active node list")
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

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])
	runGovContract(gc, buildDeclareInput(), t)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("in this case, Gov will save the declared node, and notify Staking if the proposal is passed later")
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

	chandler := node.GetCryptoHandler()
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

	chandler := node.GetCryptoHandler()
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
		t.Log("voted, Gov will save the declared node, and notify Staking if the proposal is passed later")
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

	chandler := node.GetCryptoHandler()
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
		t.Log("voted, Gov will save the declared node, and notify Staking if the proposal is passed later")
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

	chandler := node.GetCryptoHandler()
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
		t.Log("voted, Gov will save the declared node, and notify Staking if the proposal is passed later")
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

func TestGovContract_SubmitCancel_AnotherVoting(t *testing.T) {
	setup(t)
	defer clear(t)
	stateDB := gc.Evm.StateDB.(*mock.MockStateDB)
	stateDB.Prepare(txHashArr[0], lastBlockHash, 0)

	//submit a proposal
	runGovContract(gc, buildSubmitVersion(nodeIdArr[0], "versionPIPID", promoteVersion, xcom.VersionProposalVote_ConsensusRounds()), t)

	buildBlockNoCommit(1)
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	stateDB.Prepare(txHashArr[1], lastBlockHash, 1)
	//submit a proposal
	runGovContract(gc, buildSubmitCancel(nodeIdArr[1], "cancelPIPID", xcom.VersionProposalVote_ConsensusRounds()-1, txHashArr[0]), t)

	buildBlockNoCommit(3)

	stateDB.Prepare(txHashArr[3], lastBlockHash, 0)
	runGovContract(gc, buildSubmitCancel(nodeIdArr[1], "cancelPIPIDAnother", xcom.VersionProposalVote_ConsensusRounds()-1, txHashArr[0]), t, gov.VotingCancelProposalExist)
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
	//txHash = txHashArr[2] is a text proposal
	runGovContract(gc, buildSubmitTextInput(), t)

	buildBlock3()

	//try to cancel a text proposal
	runGovContract(gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", 4, txHashArr[2]), t, gov.TobeCanceledProposalTypeError)
}

func TestGovContract_SubmitCancel_EndVotingRounds_TobeCanceledNotAtVotingStage(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	//move the proposal ID from voting-list to end-list
	gov.MoveVotingProposalIDToEnd(blockHash2, txHashArr[2])

	buildBlock3()

	//try to cancel a closed version proposal
	runGovContract(gc, buildSubmitCancel(nodeIdArr[0], "cancelPIPID", 4, txHashArr[2]), t, gov.TobeCanceledProposalNotAtVoting)
}

func TestGovContract_GetCancelProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 0)

	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	buildBlock3()
	//submit a proposal and get it.
	runGovContract(gc, buildSubmitCancelInput(), t)

	runGovContract(gc, buildGetProposalInput(2), t)
}

func TestGovContract_Vote_VersionProposal(t *testing.T) {
	setup(t)
	defer clear(t)
	//submit a proposal and vote for it.
	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)
	runGovContract(gc, buildVoteInput(0, 2), t)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("voted, Gov will save the declared node, and notify Staking if the proposal is passed later")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}
func TestGovContract_Vote_Duplicated(t *testing.T) {
	setup(t)
	defer clear(t)
	//submit a proposal and vote for it.
	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)
	runGovContract(gc, buildVoteInput(0, 2), t)
	runGovContract(gc, buildVoteInput(0, 2), t, gov.VoteDuplicated)
	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("voted duplicated, Gov will count this node once in active node list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_Vote_OptionError(t *testing.T) {
	setup(t)
	defer clear(t)
	//submit a proposal and vote for it.
	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)
	// vote option = 0, it's wrong
	runGovContract(gc, buildVote(0, 2, 0, promoteVersion, versionSign), t, gov.VoteOptionError)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 0 {
		t.Log("option error, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_Vote_ProposalNotExist(t *testing.T) {
	setup(t)
	defer clear(t)
	//submit a proposal and vote for it.
	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))

	//verify vote new version, but it has not upgraded
	// txIdx=4, not a proposalID
	runGovContract(gc, buildVote(0, 4, 1, initProgramVersion, sign), t, gov.ProposalNotFound)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("list ActiveNode error", "err", err)
	} else if len(nodeList) == 0 {
		t.Log("proposal not found, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_Vote_TextProposalPassed(t *testing.T) {
	setup(t)
	defer clear(t)

	stateDB := gc.Evm.StateDB.(*mock.MockStateDB)
	stateDB.Prepare(txHashArr[0], lastBlockHash, 0)

	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitTextInput(), t)

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	build_staking_data_more(uint64(2))

	allVote(stateDB, t, txHashArr[0])
	sndb.Commit(lastBlockHash) //commit
	sndb.Compaction()          //write to level db

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}

	lastBlockNumber = uint64(p.GetEndVotingBlock() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	// build_staking_data_more will build a new block base on sndb.Current
	build_staking_data_more(p.GetEndVotingBlock())
	endBlock(stateDB, t)
	sndb.Commit(lastBlockHash)

	// vote option = 0, it's wrong
	runGovContract(gc, buildVote(0, 0, 0, promoteVersion, versionSign), t, gov.ProposalNotAtVoting)

	if nodeList, err := gov.GetActiveNodeList(lastBlockHash, txHashArr[0]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 0 {
		t.Log("option error, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}

}

func TestGovContract_Vote_VerifierNotUpgraded(t *testing.T) {
	setup(t)
	defer clear(t)
	//submit a proposal and vote for it.
	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(initProgramVersion))

	//verify vote new version, but it has not upgraded
	runGovContract(gc, buildVote(0, 2, 1, initProgramVersion, sign), t, gov.VerifierNotUpgraded)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("list ActiveNode error", "err", err)
	} else if len(nodeList) == 0 {
		t.Log("verifier not upgraded, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_Vote_ProgramVersionError(t *testing.T) {
	setup(t)
	defer clear(t)
	//submit a proposal and vote for it.
	buildBlock2()
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])

	otherVersion := uint32(1<<16 | 3<<8 | 0)
	var sign common.VersionSign
	sign.SetBytes(chandler.MustSign(otherVersion))

	//verify vote new version, but it has not upgraded
	runGovContract(gc, buildVote(0, 2, 1, otherVersion, sign), t, gov.VerifierNotUpgraded)

	if nodeList, err := gov.GetActiveNodeList(blockHash2, txHashArr[2]); err != nil {
		t.Error("list ActiveNode error", "err", err)
	} else if len(nodeList) == 0 {
		t.Log("verifier program version error, this node will not be added to active list")
	} else {
		t.Fatal("cannot list ActiveNode")
	}
}

func TestGovContract_AllNodeVoteVersionProposal(t *testing.T) {
	setup(t)
	defer clear(t)
	//submit a proposal and vote for it.
	buildBlock2()

	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := node.GetCryptoHandler()

	for i := 0; i < 3; i++ {
		chandler.SetPrivateKey(priKeyArr[i])
		var sign common.VersionSign
		sign.SetBytes(chandler.MustSign(promoteVersion))
		//verify vote new version, but it has not upgraded
		runGovContract(gc, buildVote(i, 2, 1, promoteVersion, sign), t)
	}
}

func TestGovContract_TextProposal_pass(t *testing.T) {
	setup(t)
	defer clear(t)

	stateDB := gc.Evm.StateDB.(*mock.MockStateDB)
	stateDB.Prepare(txHashArr[0], lastBlockHash, 0)

	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitTextInput(), t)

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	build_staking_data_more(uint64(2))

	allVote(stateDB, t, txHashArr[0])
	sndb.Commit(lastBlockHash) //commit
	sndb.Compaction()          //write to level db

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}

	lastBlockNumber = uint64(p.GetEndVotingBlock() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	// build_staking_data_more will build a new block base on sndb.Current
	build_staking_data_more(p.GetEndVotingBlock())
	endBlock(stateDB, t)
	sndb.Commit(lastBlockHash)

	result, err := gov.GetTallyResult(txHashArr[0], stateDB)
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
	setup(t)
	defer clear(t)

	stateDB := gc.Evm.StateDB.(*mock.MockStateDB)
	stateDB.Prepare(txHashArr[0], lastBlockHash, 0)

	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	build_staking_data_more(uint64(2))

	allVote(stateDB, t, txHashArr[0])
	sndb.Commit(lastBlockHash) //commit
	sndb.Compaction()          //write to level db

	pTemp, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}
	p := pTemp.(*gov.VersionProposal)

	// res
	lastBlockNumber = uint64(p.GetEndVotingBlock() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	// build_staking_data_more will build a new block base on sndb.Current
	build_staking_data_more(p.GetEndVotingBlock())
	endBlock(stateDB, t)

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	result, err := gov.GetTallyResult(txHashArr[0], stateDB)
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

	lastBlockNumber = uint64(p.GetActiveBlock() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	// build_staking_data_more will build a new block base on sndb.Current
	build_staking_data_more(uint64(p.GetActiveBlock()))
	beginBlock(stateDB, t)
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	result, err = gov.GetTallyResult(txHashArr[0], stateDB)
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

func TestGovContract_ListProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	stateDB := gc.Evm.StateDB.(*mock.MockStateDB)
	stateDB.Prepare(txHashArr[0], lastBlockHash, 0)

	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	stateDB.Prepare(txHashArr[1], lastBlockHash, 0)
	runGovContract(gc, buildSubmitTextInput(), t)
	sndb.Commit(lastBlockHash) //commit
	sndb.Compaction()          //write to level db

	runGovContract(gc, buildListProposalInput(), t)

}

func TestGovContract_GetActiveVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	buildBlock2()

	runGovContract(gc, buildGetActiveVersionInput(), t)
}

func TestGovContract_getAccuVerifiersCount(t *testing.T) {
	setup(t)
	defer clear(t)

	stateDB := gc.Evm.StateDB.(*mock.MockStateDB)
	stateDB.Prepare(txHashArr[0], lastBlockHash, 0)
	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)
	//submit a proposal and vote for it.
	runGovContract(gc, buildGetAccuVerifiersCountInput(txHashArr[0], lastBlockHash), t)
}

func TestGovContract_getAccuVerifiersCount_wrongProposalID(t *testing.T) {
	setup(t)
	defer clear(t)

	stateDB := gc.Evm.StateDB.(*mock.MockStateDB)
	stateDB.Prepare(txHashArr[0], lastBlockHash, 0)
	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)
	//submit a proposal and vote for it.
	runGovContract(gc, buildGetAccuVerifiersCountInput(txHashArr[1], lastBlockHash), t, gov.ProposalNotFound)
}

func runGovContract(contract *GovContract, buf []byte, t *testing.T, expectedErrors ...error) {
	res, err := contract.Run(buf)

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	if expectedErrors != nil {
		assert.NotEqual(t, common.OkCode, r.Code)
		var expected = false
		for _, expectedError := range expectedErrors {
			expected = expected || strings.Contains(r.ErrMsg, expectedError.Error())
		}
		assert.True(t, true, expected)
		t.Log("the expected result Msg:", r.ErrMsg)
	} else {
		assert.Equal(t, common.OkCode, r.Code)
		t.Log("the expected result:", r)
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

func allVote(stateDB *mock.MockStateDB, t *testing.T, pid common.Hash) {
	//for _, nodeID := range nodeIdArr {
	currentValidatorList, _ := plugin.StakingInstance().ListCurrentValidatorID(lastBlockHash, lastBlockNumber)
	voteCount := len(currentValidatorList)
	chandler := node.GetCryptoHandler()
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	for i := 0; i < voteCount; i++ {
		vote := gov.VoteInfo{
			ProposalID: pid,
			VoteNodeID: nodeIdArr[i],
			VoteOption: gov.Yes,
		}

		chandler.SetPrivateKey(priKeyArr[i])
		versionSign := common.VersionSign{}
		versionSign.SetBytes(chandler.MustSign(promoteVersion))

		err := gov.Vote(sender, vote, lastBlockHash, 1, promoteVersion, versionSign, plugin.StakingInstance(), stateDB)
		if err != nil {
			t.Fatalf("vote err: %s.", err)
		}
	}
}

func beginBlock(stateDB *mock.MockStateDB, t *testing.T) {
	err := govPlugin.BeginBlock(lastBlockHash, &lastHeader, stateDB)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
}

func endBlock(stateDB *mock.MockStateDB, t *testing.T) {
	err := govPlugin.EndBlock(lastBlockHash, &lastHeader, stateDB)
	if err != nil {
		t.Fatalf("end block err... %s", err)
	}
}
