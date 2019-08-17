package gov

import (
	"encoding/hex"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/x/plugin"

	govdb "github.com/PlatONnetwork/PlatON-Go/x/gov/db"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	//	"github.com/PlatONnetwork/PlatON-Go/core/state"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var (
	snapdb    snapshotdb.DB
	govPlugin *plugin.GovPlugin
	//newVersion     = uint32(2<<16 | 0<<8 | 0)
	endVotingRounds uint64
	endVotingBlock  uint64
	activeBlock     uint64
	stateDB         xcom.StateDB

	stk *plugin.StakingPlugin
)

func setup(t *testing.T) func() {
	t.Log("setup()......")

	state, genesis, _ := newChainState()
	newEvm(blockNumber, blockHash, state)
	stateDB = state
	newPlugins()

	stk = plugin.StakingInstance()

	lastBlockHash = genesis.Hash()

	build_staking_data(genesis.Hash())

	snapdb = snapshotdb.Instance()
	// init data
	endVotingRounds = 5
	endVotingBlock = xutil.CalEndVotingBlock(uint64(1), endVotingRounds)
	activeBlock = xutil.CalActiveBlock(endVotingBlock)
	return func() {
		t.Log("tear down()......")
		snapdb.Clear()
	}
}

func submitText(t *testing.T, pid common.Hash) {
	vp := &TextProposal{
		ProposalID:   pid,
		ProposalType: Text,
		PIPID:        "textPIPID",
		SubmitBlock:  1,
		Proposer:     nodeIdArr[0],
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := Submit(sender, vp, lastBlockHash, lastBlockNumber, stateDB)
	if err != nil {
		t.Fatalf("submit text proposal err: %s", err)
	}
}

func submitVersion(t *testing.T, pid common.Hash) {
	vp := &VersionProposal{
		ProposalID:      pid,
		ProposalType:    Version,
		PIPID:           "versionIPID",
		SubmitBlock:     1,
		EndVotingRounds: endVotingRounds,
		Proposer:        nodeIdArr[0],
		NewVersion:      promoteVersion,
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], blockHash, 0)

	err := Submit(sender, vp, lastBlockHash, lastBlockNumber, stateDB)
	if err != nil {
		t.Fatalf("submit version proposal err: %s", err)
	}
}

func submitCancel(t *testing.T, pid, tobeCanceled common.Hash) {
	pp := &CancelProposal{
		ProposalID:      pid,
		ProposalType:    Cancel,
		PIPID:           "CancelPIPID",
		SubmitBlock:     1,
		EndVotingRounds: endVotingRounds,
		Proposer:        nodeIdArr[0],
		TobeCanceled:    tobeCanceled,
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := Submit(sender, pp, lastBlockHash, lastBlockNumber, stateDB)
	if err != nil {
		t.Fatalf("submit cancel proposal err: %s", err)
	}
}

func allVote(t *testing.T, pid common.Hash) {
	//for _, nodeID := range nodeIdArr {
	currentValidatorList, _ := stk.ListCurrentValidatorID(lastBlockHash, lastBlockNumber)
	voteCount := len(currentValidatorList)
	chandler := xcom.GetCryptoHandler()

	for i := 0; i < voteCount; i++ {
		vote := VoteInfo{
			ProposalID: pid,
			VoteNodeID: nodeIdArr[i],
			VoteOption: Yes,
		}

		chandler.SetPrivateKey(priKeyArr[i])
		versionSign := common.VersionSign{}
		versionSign.SetBytes(chandler.MustSign(promoteVersion))

		err := Vote(sender, vote, lastBlockHash, 1, promoteVersion, versionSign, stateDB)
		if err != nil {
			t.Fatalf("vote err: %s.", err)
		}
	}
}

func halfVote(t *testing.T, pid common.Hash) {
	currentValidatorList, _ := stk.ListCurrentValidatorID(lastBlockHash, lastBlockNumber)
	voteCount := len(currentValidatorList)
	chandler := xcom.GetCryptoHandler()
	for i := 0; i < voteCount/2; i++ {
		vote := VoteInfo{
			ProposalID: pid,
			VoteNodeID: nodeIdArr[i],
			VoteOption: Yes,
		}

		chandler.SetPrivateKey(priKeyArr[i])
		versionSign := common.VersionSign{}
		versionSign.SetBytes(chandler.MustSign(promoteVersion))

		err := Vote(sender, vote, lastBlockHash, 1, promoteVersion, versionSign, stateDB)
		if err != nil {
			t.Fatalf("vote err: %s.", err)
		}
	}

}

func beginBlock(t *testing.T) {
	err := govPlugin.BeginBlock(lastBlockHash, &lastHeader, stateDB)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
}

func endBlock(t *testing.T) {
	err := govPlugin.EndBlock(lastBlockHash, &lastHeader, stateDB)
	if err != nil {
		t.Fatalf("end block err... %s", err)
	}
}

func TestGovPlugin_SubmitText(t *testing.T) {
	defer setup(t)()
	submitText(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := govdb.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted text proposal error:", err)
	} else {
		t.Log("Get the submitted text proposal success:", p)
	}
}

func TestGovPlugin_GetProposal(t *testing.T) {
	defer setup(t)()
	submitText(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	buildBlockNoCommit(2)

	p, err := govdb.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatalf("Get proposal error: %s", err)
	} else {
		t.Logf("Get proposal success: %x", p.GetProposalID())
	}
}

func TestGovPlugin_SubmitText_invalidSender(t *testing.T) {
	defer setup(t)()

	vp := &TextProposal{
		ProposalID:   txHashArr[0],
		ProposalType: Text,
		PIPID:        "textPIPID",
		SubmitBlock:  1,
		Proposer:     nodeIdArr[0],
	}

	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := Submit(anotherSender, vp, lastBlockHash, lastBlockNumber, stateDB) //sender error
	if err != nil && (err.Error() == "tx sender is not verifier." || err.Error() == "tx sender should be node's staking address.") {
		t.Log("detected invalid sender.", err)
	} else {
		t.Fatal("didn't detect invalid sender.")
	}
}

func TestGovPlugin_SubmitVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := govdb.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}
}

func TestGovPlugin_SubmitCancel(t *testing.T) {
	defer setup(t)()

	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := govdb.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}

	submitCancel(t, txHashArr[1], txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err = govdb.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted cancel proposal error:", err)
	} else {
		t.Log("Get the submitted cancel proposal success:", p)
	}
}

func TestGovPlugin_SubmitCancel_invalidEndVotingRounds(t *testing.T) {
	defer setup(t)()

	pp := &CancelProposal{
		ProposalID:      txHashArr[1],
		ProposalType:    Cancel,
		PIPID:           "cancelPIPID",
		SubmitBlock:     1,
		EndVotingRounds: 6, //>5, error
		Proposer:        nodeIdArr[1],
		TobeCanceled:    txHashArr[0],
	}
	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[1], lastBlockHash, 0)

	err := Submit(sender, pp, lastBlockHash, lastBlockNumber, stateDB)
	if err != nil && err.Error() == "end-voting-block invalid." {
		t.Logf("detected invalid end-voting-block.")
	} else {
		t.Fatal("didn't detect invalid end-voting-block.")
	}
}

func TestGovPlugin_SubmitCancel_noVersionProposal(t *testing.T) {
	defer setup(t)()

	pp := &CancelProposal{
		ProposalID:      txHashArr[1],
		ProposalType:    Cancel,
		PIPID:           "cancelPIPID",
		SubmitBlock:     1,
		EndVotingRounds: endVotingRounds,
		Proposer:        nodeIdArr[0],
		TobeCanceled:    txHashArr[0],
	}
	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := Submit(sender, pp, lastBlockHash, lastBlockNumber, stateDB)
	if err != nil && err.Error() == "find to be canceled version proposal error" {
		t.Logf("detected this case.")
	} else {
		t.Fatal("didn't detect is case.")
	}
}

func TestGovPlugin_VoteSuccess(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	nodeIdx := 3
	v := VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		Yes,
	}

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	nodeIdx = 1
	v = VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		Yes,
	}

	chandler = xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign = common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err = Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	votedValue, err := govdb.ListVoteValue(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	} else {
		t.Log("voted count:", len(votedValue))
	}

	nodeList, err := govdb.ListVotedVerifier(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("vote failed, cannot list voted verifiers", err)
	} else {
		t.Log("voted count:", len(nodeList))
	}
}

func TestGovPlugin_Vote_Repeat(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	nodeIdx := 3
	v := VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		Yes,
	}

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	v = VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx], //repeated
		Yes,
	}

	err = Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stateDB)
	if err != nil && err.Error() == "node has voted this proposal." {
		t.Log("detected repeated vote", err)
	} else {
		t.Fatal("didn't detect repeated vote")
	}
}

func TestGovPlugin_Vote_invalidSender(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	nodeIdx := 3
	v := VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		Yes,
	}

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := Vote(anotherSender, v, lastBlockHash, 2, initProgramVersion, versionSign, stateDB)
	if err != nil && err.Error() == "tx sender is not a verifier, or mismatch the verifier's nodeID" {
		t.Log("vote err:", err)
	}
	votedValue, err := govdb.ListVoteValue(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	} else {
		t.Log("voted count:", len(votedValue))
	}
}

func TestGovPlugin_DeclareVersion_rightVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	nodeIdx := 0
	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := DeclareVersion(sender, nodeIdArr[nodeIdx], promoteVersion, versionSign, lastBlockHash, 2, stateDB)
	if err != nil {
		t.Fatalf("Declare Version err ...%s", err)
	}

	activeNodeList, err := govdb.GetActiveNodeList(lastBlockHash, txHashArr[0])
	if err != nil {
		t.Fatalf("List actived nodes error: %s", err)
	} else {
		t.Logf("List actived nodes success: %d", len(activeNodeList))
	}
}

func TestGovPlugin_DeclareVersion_wrongSign(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	wrongVersion := uint32(1<<16 | 1<<8 | 1)

	nodeIdx := 0
	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(wrongVersion))

	err := DeclareVersion(sender, nodeIdArr[nodeIdx], promoteVersion, versionSign, lastBlockHash, 2, stateDB)

	if err != nil && (err.Error() == "version sign error." || err.Error() == "declared version neither equals active version nor new version.") {
		t.Log("system has detected an incorrect version declaration.", err)
	} else {
		t.Fatal("system has not detected an incorrect version declaration.", err)
	}
}

func TestGovPlugin_DeclareVersion_wrongVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	wrongVersion := uint32(1<<16 | 1<<8 | 1)

	nodeIdx := 0
	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(wrongVersion))

	err := DeclareVersion(sender, nodeIdArr[nodeIdx], wrongVersion, versionSign, lastBlockHash, 2, stateDB)

	if err != nil && (err.Error() == "version sign error." || err.Error() == "declared version neither equals active version nor new version.") {
		t.Log("system has detected an incorrect version declaration.", err)
	} else {
		t.Fatal("system has not detected an incorrect version declaration.", err)
	}
}

func TestGovPlugin_VotedNew_DeclareOld(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	nodeIdx := 3
	v := VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		Yes,
	}

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	nodeIdx = 1
	v = VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		Yes,
	}

	chandler = xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign = common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err = Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	votedValue, err := govdb.ListVoteValue(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	} else {
		t.Log("voted count:", len(votedValue))
	}

	//declare
	versionSign = common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(initProgramVersion))

	err = DeclareVersion(sender, nodeIdArr[nodeIdx], initProgramVersion, versionSign, lastBlockHash, 2, stateDB)

	if err != nil && (err.Error() == "version sign error." || err.Error() == "declared version should be same as proposal's version") {
		t.Log("system has detected an incorrect version declaration.", err)
	} else {
		t.Fatal("system has not detected an incorrect version declaration.", err)
	}
}

func TestGovPlugin_DeclareVersion_invalidSender(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	nodeIdx := 0
	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := DeclareVersion(anotherSender, nodeIdArr[nodeIdx], promoteVersion, versionSign, lastBlockHash, 2, stateDB)
	if err != nil && (err.Error() == "tx sender is not candidate." || err.Error() == "tx sender should be node's staking address.") {
		t.Log("detected an incorrect version declaration.", err)
	} else {
		t.Fatal("didn't detected an incorrect version declaration.", err)
	}
}

func TestGovPlugin_ListProposal(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	pList, err := ListProposal(lastBlockHash, stateDB)
	if err != nil {
		t.Fatalf("List all proposals error: %s", err)
	} else {
		t.Logf("List all proposals success: %d", len(pList))
	}

}

func TestGovPlugin_textProposalPassed(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	allVote(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(xutil.CalcBlocksEachEpoch() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(uint64(xutil.CalcBlocksEachEpoch()))
	beginBlock(t)
	sndb.Commit(lastBlockHash)

	//buildSnapDBDataCommitted(20001, 22229)
	sndb.Compaction()
	lastBlockNumber = uint64(uint64(xutil.ConsensusSize()*5-xcom.ElectionDistance()) - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(endVotingBlock)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govdb.GetTallyResult(txHashArr[0], stateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == Pass {
		t.Logf("the result status, %s", result.Status.ToString())
	} else {
		t.Fatalf("the result status error, %s", result.Status.ToString())
	}
}

func TestGovPlugin_textProposalFailed(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	halfVote(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(xutil.CalcBlocksEachEpoch() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(uint64(xutil.CalcBlocksEachEpoch()))
	beginBlock(t)
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(endVotingBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(endVotingBlock)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govdb.GetTallyResult(txHashArr[0], stateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == Failed {
		t.Logf("the result status, %s", result.Status.ToString())
	} else {
		t.Fatalf("the result status error, %s", result.Status.ToString())
	}
}

func TestGovPlugin_versionProposalPreActive(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])
	submitVersion(t, txHashArr[1])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	allVote(t, txHashArr[0])
	allVote(t, txHashArr[1])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(xutil.CalcBlocksEachEpoch() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(uint64(xutil.CalcBlocksEachEpoch()))

	beginBlock(t)

	sndb.Commit(lastBlockHash)

	//buildSnapDBDataCommitted(20001, 22229)
	sndb.Compaction()
	lastBlockNumber = uint64(endVotingBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(endVotingBlock)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govdb.GetTallyResult(txHashArr[0], stateDB)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else {
		t.Logf("the result status, %s", result.Status.ToString())
	}

	result, err = govdb.GetTallyResult(txHashArr[1], stateDB)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == PreActive {
		t.Logf("the result status, %s", result.Status.ToString())
	} else {
		t.Logf("the result status error, %s", result.Status.ToString())
	}
}

func TestGovPlugin_GetPreActiveVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	buildBlockNoCommit(2)

	ver := govdb.GetPreActiveVersion(stateDB)
	t.Logf("Get pre-active version: %d", ver)
}

func TestGovPlugin_GetActiveVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	buildBlockNoCommit(2)

	ver := GetCurrentActiveVersion(stateDB)
	t.Logf("Get active version: %d", ver)
}

func TestGovPlugin_versionProposalActive(t *testing.T) {

	defer setup(t)()

	//submit version proposal
	submitVersion(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction() //flush to LevelDB

	buildBlockNoCommit(2)
	//voting
	allVote(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(endVotingBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(endVotingBlock)

	//tally result
	endBlock(t)
	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	lastBlockNumber = uint64(activeBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	//buildBlockNoCommit(23480)
	build_staking_data_more(uint64(activeBlock))
	//active
	beginBlock(t)
	sndb.Commit(lastBlockHash)

	activeVersion := GetCurrentActiveVersion(stateDB)
	if activeVersion == promoteVersion {
		t.Logf("active SUCCESS, %d", activeVersion)
	} else {
		t.Fatalf("active FALSE, %d", activeVersion)
	}
}

func TestGovPlugin_printVersion(t *testing.T) {
	defer setup(t)()

	t.Logf("ver.1.2.0, %d", uint32(1<<16|2<<8|0))

}

func TestNodeID(t *testing.T) {
	var nodeID discover.NodeID
	nodeID = [64]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	t.Logf("nodeID is empty, %t", nodeID == discover.ZeroNodeID)

}

/*func TestNodeID1(t *testing.T) {
	var nodeID discover.NodeID
	nodeID = [64]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	t.Error("nodeID is empty", "nodeID", nodeID)

	var proposalID common.Hash
	proposalID = [32]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	t.Error("proposalID is empty", "proposalID", proposalID)
}*/

func Test_MakeExtraData(t *testing.T) {
	defer setup(t)()

	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	t.Log(lastHeader.Extra)
	beginBlock(t)
	t.Log(lastHeader.Extra)

	if len(lastHeader.Extra) > 0 {
		var tobeDecoded []byte
		tobeDecoded = lastHeader.Extra
		if len(lastHeader.Extra) <= 32 {
			tobeDecoded = lastHeader.Extra
		} else {
			tobeDecoded = lastHeader.Extra[:32]
		}

		var extraData []interface{}
		err := rlp.DecodeBytes(tobeDecoded, &extraData)
		if err != nil {
			t.Error("rlp decode header extra error")
		}
		//reference to makeExtraData() in gov_plugin.go
		if len(extraData) == 4 {
			versionBytes := extraData[0].([]byte)
			versionInHeader := common.BytesToUint32(versionBytes)

			activeVersion := GetActiveVersion(lastHeader.Number.Uint64(), stateDB)
			t.Log("verify header version", "headerVersion", versionInHeader, "activeVersion", activeVersion, "blockNumber", lastHeader.Number.Uint64())
			if activeVersion == versionInHeader {
				t.Log("OK")
			} else {
				t.Error("header version error")
			}
		} else {
			t.Error("unknown header extra data", "elementCount", len(extraData))
		}
	}

}

func Test_version(t *testing.T) {
	ver := uint32(66048) //1.2.0
	t.Log(common.Uint32ToBytes(ver))
}

func Test_genVersionSign(t *testing.T) {

	ver := uint32(66048) //1.2.0
	chandler := xcom.GetCryptoHandler()

	for i := 0; i < 18; i++ {
		chandler.SetPrivateKey(priKeyArr[i])
		t.Log("0x" + hex.EncodeToString(chandler.MustSign(ver)))
	}

}
