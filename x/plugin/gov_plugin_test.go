package plugin

import (
	"encoding/hex"
	"encoding/json"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	//	"github.com/PlatONnetwork/PlatON-Go/core/state"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	//	"github.com/PlatONnetwork/PlatON-Go/core/state"
	//	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var (
	snapdb    snapshotdb.DB
	govPlugin *GovPlugin
	//	evm       *vm.EVM
	//newVersion     = uint32(2<<16 | 0<<8 | 0)
	/*endVotingBlock uint64
	activeBlock    uint64*/
	stateDB xcom.StateDB

//	stk            *StakingPlugin
)

func setup(t *testing.T) func() {
	t.Log("setup()......")

	state, genesis, _ := newChainState()
	newEvm(blockNumber, blockHash, state)
	stateDB = state
	newPlugins()

	govPlugin = GovPluginInstance()
	stk = StakingInstance()

	lastBlockHash = genesis.Hash()

	build_staking_data(genesis.Hash())

	snapdb = snapshotdb.Instance()

	// init data

	return func() {
		t.Log("tear down()......")
		snapdb.Clear()
	}
}

func submitText(t *testing.T, pid common.Hash) {
	vp := &gov.TextProposal{
		ProposalID:   pid,
		ProposalType: gov.Text,
		PIPID:        "textPIPID",
		SubmitBlock:  1,
		Proposer:     nodeIdArr[0],
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], lastBlockHash, 0)
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB)
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	if err != nil {
		t.Fatalf("submit text proposal err: %s", err)
	}
}

func submitVersion(t *testing.T, pid common.Hash) {
	vp := &gov.VersionProposal{
		ProposalID:      pid,
		ProposalType:    gov.Version,
		PIPID:           "versionIPID",
		SubmitBlock:     1,
		EndVotingRounds: xcom.VersionProposalVote_ConsensusRounds(),
		Proposer:        nodeIdArr[0],
		NewVersion:      promoteVersion,
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], blockHash, 0)

	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB)
	if err != nil {
		t.Fatalf("submit version proposal err: %s", err)
	}
}

func submitCancel(t *testing.T, pid, tobeCanceled common.Hash) {
	pp := &gov.CancelProposal{
		ProposalID:      pid,
		ProposalType:    gov.Cancel,
		PIPID:           "CancelPIPID",
		SubmitBlock:     1,
		EndVotingRounds: xcom.VersionProposalVote_ConsensusRounds() - 1,
		Proposer:        nodeIdArr[0],
		TobeCanceled:    tobeCanceled,
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(sender, pp, lastBlockHash, lastBlockNumber, stk, stateDB)
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
		vote := gov.VoteInfo{
			ProposalID: pid,
			VoteNodeID: nodeIdArr[i],
			VoteOption: gov.Yes,
		}

		chandler.SetPrivateKey(priKeyArr[i])
		versionSign := common.VersionSign{}
		versionSign.SetBytes(chandler.MustSign(promoteVersion))

		err := gov.Vote(sender, vote, lastBlockHash, 1, promoteVersion, versionSign, stk, stateDB)
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
		vote := gov.VoteInfo{
			ProposalID: pid,
			VoteNodeID: nodeIdArr[i],
			VoteOption: gov.Yes,
		}

		chandler.SetPrivateKey(priKeyArr[i])
		versionSign := common.VersionSign{}
		versionSign.SetBytes(chandler.MustSign(promoteVersion))

		err := gov.Vote(sender, vote, lastBlockHash, 1, promoteVersion, versionSign, stk, stateDB)
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

	p, err := gov.GetProposal(txHashArr[0], stateDB)
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

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatalf("Get proposal error: %s", err)
	} else {
		t.Logf("Get proposal success: %x", p.GetProposalID())
	}
}

func TestGovPlugin_SubmitText_invalidSender(t *testing.T) {
	defer setup(t)()

	vp := &gov.TextProposal{
		ProposalID:   txHashArr[0],
		ProposalType: gov.Text,
		PIPID:        "textPIPID",
		SubmitBlock:  1,
		Proposer:     nodeIdArr[0],
	}

	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(anotherSender, vp, lastBlockHash, lastBlockNumber, stk, stateDB) //sender error
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

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}
}

func TestGovPlugin_SubmitVersion_invalidEndVotingRounds(t *testing.T) {
	defer setup(t)()

	vp := &gov.VersionProposal{
		ProposalID:      txHashArr[0],
		ProposalType:    gov.Version,
		PIPID:           "versionPIPID",
		SubmitBlock:     1,
		EndVotingRounds: xcom.VersionProposalVote_ConsensusRounds() + 1, //error
		Proposer:        nodeIdArr[0],
		NewVersion:      promoteVersion,
	}
	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB)
	if err != nil && err.Error() == "voting consensus rounds too large." {
		t.Logf("detected invalid end-voting-rounds.")
	} else {
		t.Fatal("didn't detect invalid end-voting-rounds.")
	}
}

func TestGovPlugin_SubmitCancel(t *testing.T) {
	defer setup(t)()

	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}

	submitCancel(t, txHashArr[1], txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err = gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted cancel proposal error:", err)
	} else {
		t.Log("Get the submitted cancel proposal success:", p)
	}
}

func TestGovPlugin_SubmitCancel_invalidEndVotingRounds(t *testing.T) {
	defer setup(t)()

	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}

	pp := &gov.CancelProposal{
		ProposalID:      txHashArr[1],
		ProposalType:    gov.Cancel,
		PIPID:           "CancelPIPID",
		SubmitBlock:     1,
		EndVotingRounds: xcom.VersionProposalVote_ConsensusRounds(),
		Proposer:        nodeIdArr[1],
		TobeCanceled:    txHashArr[0],
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], lastBlockHash, 0)

	err = gov.Submit(sender, pp, lastBlockHash, lastBlockNumber, stk, stateDB)
	if err != nil && err.Error() == "voting consensus rounds too large." {
		t.Logf("detected invalid end-voting-rounds.")
	} else {
		t.Fatal("didn't detect invalid end-voting-rounds.")
	}
}

func TestGovPlugin_SubmitCancel_noVersionProposal(t *testing.T) {
	defer setup(t)()

	pp := &gov.CancelProposal{
		ProposalID:      txHashArr[1],
		ProposalType:    gov.Cancel,
		PIPID:           "cancelPIPID",
		SubmitBlock:     1,
		EndVotingRounds: xcom.VersionProposalVote_ConsensusRounds() - 1,
		Proposer:        nodeIdArr[0],
		TobeCanceled:    txHashArr[0],
	}
	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(sender, pp, lastBlockHash, lastBlockNumber, stk, stateDB)
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
	v := gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	nodeIdx = 1
	v = gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler = xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign = common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err = gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	votedValue, err := gov.ListVoteValue(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	} else {
		t.Log("voted count:", len(votedValue))
	}

	nodeList, err := gov.ListVotedVerifier(txHashArr[0], stateDB)
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
	v := gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	v = gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx], //repeated
		gov.Yes,
	}

	err = gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
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
	v := gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.Vote(anotherSender, v, lastBlockHash, 2, initProgramVersion, versionSign, stk, stateDB)
	if err != nil && err.Error() == "tx sender is not a verifier, or mismatch the verifier's nodeID" {
		t.Log("vote err:", err)
	}
	votedValue, err := gov.ListVoteValue(txHashArr[0], stateDB)
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

	err := gov.DeclareVersion(sender, nodeIdArr[nodeIdx], promoteVersion, versionSign, lastBlockHash, 2, stk, stateDB)
	if err != nil {
		t.Fatalf("Declare Version err ...%s", err)
	}

	activeNodeList, err := gov.GetActiveNodeList(lastBlockHash, txHashArr[0])
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

	err := gov.DeclareVersion(sender, nodeIdArr[nodeIdx], promoteVersion, versionSign, lastBlockHash, 2, stk, stateDB)

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

	err := gov.DeclareVersion(sender, nodeIdArr[nodeIdx], wrongVersion, versionSign, lastBlockHash, 2, stk, stateDB)

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
	v := gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	nodeIdx = 1
	v = gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler = xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign = common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err = gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	votedValue, err := gov.ListVoteValue(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	} else {
		t.Log("voted count:", len(votedValue))
	}

	//declare
	versionSign = common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(initProgramVersion))

	err = gov.DeclareVersion(sender, nodeIdArr[nodeIdx], initProgramVersion, versionSign, lastBlockHash, 2, stk, stateDB)

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

	err := gov.DeclareVersion(anotherSender, nodeIdArr[nodeIdx], promoteVersion, versionSign, lastBlockHash, 2, stk, stateDB)
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

	pList, err := gov.ListProposal(lastBlockHash, stateDB)
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

	endVotingBlock := xutil.CalEndVotingBlock(1, xcom.VersionProposalVote_ConsensusRounds())
	//	actvieBlock := xutil.CalActiveBlock(endVotingBlock)

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

	lastBlockNumber = uint64(endVotingBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(endVotingBlock)

	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	endBlock(t)

	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	sndb.Commit(lastBlockHash)

	result, err := gov.GetTallyResult(txHashArr[0], stateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.Pass {
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

	endVotingBlock := xutil.CalEndVotingBlock(1, xcom.VersionProposalVote_ConsensusRounds())
	//	actvieBlock := xutil.CalActiveBlock(endVotingBlock)

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

	result, err := gov.GetTallyResult(txHashArr[0], stateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.Failed {
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

	endVotingBlock := xutil.CalEndVotingBlock(1, xcom.VersionProposalVote_ConsensusRounds())
	//	actvieBlock := xutil.CalActiveBlock(endVotingBlock)

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

	result, err := gov.GetTallyResult(txHashArr[0], stateDB)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else {
		t.Logf("the result status, %s", result.Status.ToString())
	}

	result, err = gov.GetTallyResult(txHashArr[1], stateDB)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.PreActive {
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

	ver := gov.GetPreActiveVersion(stateDB)
	t.Logf("Get pre-active version: %d", ver)
}

func TestGovPlugin_GetActiveVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	buildBlockNoCommit(2)

	ver := gov.GetCurrentActiveVersion(stateDB)
	t.Logf("Get active version: %d", ver)
}

func TestGovPlugin_versionProposalActive(t *testing.T) {

	defer setup(t)()

	//submit version proposal
	submitVersion(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction() //flush to LevelDB

	endVotingBlock := xutil.CalEndVotingBlock(1, xcom.VersionProposalVote_ConsensusRounds())
	actvieBlock := xutil.CalActiveBlock(endVotingBlock)

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
	lastBlockNumber = uint64(actvieBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	//buildBlockNoCommit(23480)
	build_staking_data_more(actvieBlock)
	//active
	beginBlock(t)
	sndb.Commit(lastBlockHash)

	activeVersion := gov.GetCurrentActiveVersion(stateDB)
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

func TestGovPlugin_TestNodeID(t *testing.T) {
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

func TestGovPlugin_Test_MakeExtraData(t *testing.T) {
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

			activeVersion := gov.GetActiveVersion(lastHeader.Number.Uint64(), stateDB)
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

func TestGovPlugin_Test_version(t *testing.T) {
	ver := uint32(66048) //1.2.0
	t.Log(common.Uint32ToBytes(ver))
}

func TestGovPlugin_Test_genVersionSign(t *testing.T) {

	ver := uint32(66048) //1.2.0
	chandler := xcom.GetCryptoHandler()

	for i := 0; i < 4; i++ {
		chandler.SetPrivateKey(priKeyArr[i])
		t.Log("0x" + hex.EncodeToString(chandler.MustSign(ver)))
	}

}

var (
	chandler *xcom.CryptoHandler
	priKey   = crypto.HexMustToECDSA("8e1477549bea04b97ea15911e2e9b3041b7a9921f80bd6ddbe4c2b080473de22")
	nodeID   = discover.MustHexID("3e7864716b671c4de0dc2d7fd86215e0dcb8419e66430a770294eb2f37b714a07b6a3493055bb2d733dee9bfcc995e1c8e7885f338a69bf6c28930f3cf341819")
)

func initChandlerHandler() {
	chandler = xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKey)
}

func TestGovPlugin_Test_Encode(t *testing.T) {
	initChandlerHandler()
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	sig, err := chandler.Sign(uint32(1792))

	value := &gov.ProgramVersionValue{ProgramVersion: uint32(1792), ProgramVersionSign: hexutil.Encode(sig)}

	jsonByte, err := json.Marshal(value)
	if nil != err {
		log.Error("json.Marshal err", "err", err)
	}

	log.Error("encode result", "sig", hex.EncodeToString(jsonByte))
	res := xcom.Result{true, string(jsonByte), ""}
	resultBytes, _ := json.Marshal(res)
	log.Error("encode result", "bytes", hex.EncodeToString(resultBytes))
}
