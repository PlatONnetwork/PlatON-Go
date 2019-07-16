package plugin_test

import (
	"os"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

var (
	snapdb    snapshotdb.DB
	govPlugin *plugin.GovPlugin
	evm       *vm.EVM
	govDB     *gov.GovDB
)

func setup(t *testing.T) func() {

	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	t.Log("setup()......")

	state, genesis, _ := newChainState()
	evm = newEvm(blockNumber, blockHash, state)

	newPlugins()

	govPlugin = plugin.GovPluginInstance()

	lastBlockHash = genesis.Hash()

	build_staking_data(genesis.Hash())

	snapdb = snapshotdb.Instance()

	govDB = gov.GovDBInstance()

	return func() {
		t.Log("tear down()......")
		snapdb.Clear()
	}
}

func getVerProposal() gov.VersionProposal {
	return gov.VersionProposal{
		common.Hash{0x01},
		"p#01",
		gov.Version,
		"up,up,up....",
		"version proposal example",
		"http://url",
		uint64(1000),
		uint64(10000),
		discover.NodeID{0x11},
		gov.TallyResult{},
		uint32(1<<16 | 1<<8 | 1),
		uint64(11250),
	}
}

func submitText(t *testing.T, pid common.Hash) {
	vp := gov.TextProposal{
		ProposalID:     pid,
		GithubID:       "githubID",
		ProposalType:   gov.Text,
		Topic:          "textTopic",
		Desc:           "textDesc",
		Url:            "textUrl",
		SubmitBlock:    1,
		EndVotingBlock: uint64(22230),
		Proposer:       nodeIdArr[0],
	}

	state := evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := govPlugin.Submit(sender, vp, lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("submit text proposal err: %s", err)
	}
}

func submitVersion(t *testing.T, pid common.Hash) {
	vp := gov.VersionProposal{
		ProposalID:     pid,
		GithubID:       "githubID",
		ProposalType:   gov.Version,
		Topic:          "versionTopic",
		Desc:           "versionDesc",
		Url:            "versionUrl",
		SubmitBlock:    1,
		EndVotingBlock: uint64(22230),
		Proposer:       nodeIdArr[0],
		NewVersion:     uint32(1<<16 | 1<<8 | 1),
		ActiveBlock:    uint64(23480),
	}

	state := evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], blockHash, 0)

	err := govPlugin.Submit(sender, vp, lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("submit version proposal err: %s", err)
	}
}

func submitParam(t *testing.T, pid common.Hash) {
	vp := gov.ParamProposal{
		ProposalID:     pid,
		GithubID:       "githubID",
		ProposalType:   gov.Param,
		Topic:          "paramTopic",
		Desc:           "paramDesc",
		Url:            "paramUrl",
		SubmitBlock:    1,
		EndVotingBlock: uint64(22230),
		Proposer:       nodeIdArr[0],

		ParamName:    "param3",
		CurrentValue: 12.5,
		NewValue:     0.85,
	}

	state := evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := govPlugin.Submit(sender, vp, lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("submit param proposal err: %s", err)
	}
}

func allVote(t *testing.T, pid common.Hash) {
	for _, nodeID := range nodeIdArr {
		vote := gov.Vote{
			ProposalID: pid,
			VoteNodeID: nodeID,
			VoteOption: gov.Yes,
		}
		err := govPlugin.Vote(sender, vote, lastBlockHash, 1, evm.StateDB)
		if err != nil {
			t.Fatalf("vote err: %s.", err)
		}
	}
}

func halfVote(t *testing.T, pid common.Hash) {
	for i := 0; i < len(nodeIdArr)/2; i++ {
		vote := gov.Vote{
			ProposalID: pid,
			VoteNodeID: nodeIdArr[i],
			VoteOption: gov.Yes,
		}
		err := govPlugin.Vote(sender, vote, lastBlockHash, 1, evm.StateDB)
		if err != nil {
			t.Fatalf("vote err: %s.", err)
		}
	}

}

func beginBlock(t *testing.T) {
	err := govPlugin.BeginBlock(lastBlockHash, &lastHeader, evm.StateDB)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
}

func endBlock(t *testing.T) {
	err := govPlugin.EndBlock(lastBlockHash, &lastHeader, evm.StateDB)
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

	p, err := govDB.GetProposal(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Fatal("Get the submitted text proposal error:", err)
	} else {
		t.Log("Get the submitted text proposal success:", p)
	}
}

func TestGovPlugin_SubmitVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := govDB.GetProposal(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}
}

func TestGovPlugin_SubmitParam(t *testing.T) {
	defer setup(t)()
	submitParam(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := govDB.GetProposal(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}
}

func TestGovPlugin_VoteSuccess(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	v := gov.Vote{
		txHashArr[0],
		nodeIdArr[3],
		gov.Yes,
	}

	err := govPlugin.Vote(sender, v, lastBlockHash, 2, evm.StateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	v = gov.Vote{
		txHashArr[0],
		nodeIdArr[1],
		gov.Yes,
	}

	err = govPlugin.Vote(sender, v, lastBlockHash, 2, evm.StateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	votedValue, err := govDB.ListVoteValue(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	} else {
		t.Log("voted count:", len(votedValue))
	}
}

func TestGovPlugin_Vote_senderError(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	v := gov.Vote{
		txHashArr[0],
		nodeIdArr[3],
		gov.Yes,
	}

	err := govPlugin.Vote(anotherSender, v, lastBlockHash, 2, evm.StateDB)
	if err != nil && err.Error() == "tx sender is not a verifier, or mismatch the verifier's nodeID" {
		t.Log("vote err:", err)
	}
	votedValue, err := govDB.ListVoteValue(txHashArr[0], evm.StateDB)
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

	err := govPlugin.DeclareVersion(sender, nodeIdArr[0], getVerProposal().NewVersion, lastBlockHash, 2, evm.StateDB)
	if err != nil {
		t.Fatalf("Declare Version err ...%s", err)
	}

	activeNodeList, err := govDB.GetActiveNodeList(lastBlockHash, txHashArr[0])
	if err != nil {
		t.Fatalf("List actived nodes error: %s", err)
	} else {
		t.Logf("List actived nodes success: %d", len(activeNodeList))
	}
}

func TestGovPlugin_DeclareVersion_wrongVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	err := govPlugin.DeclareVersion(sender, nodeIdArr[0], uint32(1<<16|2<<8|1), lastBlockHash, 2, evm.StateDB)
	if err != nil && err.Error() == "declared version neither equals active version nor new version." {
		t.Log("system has detected an incorrect version declaration.", err)
	} else {
		t.Fatal("system has not detected an incorrect version declaration.", err)
	}
}

func TestGovPlugin_ListProposal(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	pList, err := govPlugin.ListProposal(lastBlockHash, evm.StateDB)
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

	lastBlockNumber = uint64(21999)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22000)
	beginBlock(t)
	sndb.Commit(lastBlockHash)

	//buildSnapDBDataCommitted(20001, 22229)
	sndb.Compaction()
	lastBlockNumber = uint64(22229)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govPlugin.GetTallyResult(txHashArr[0], evm.StateDB)
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

	buildBlockNoCommit(2)

	halfVote(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(21999)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22000)
	beginBlock(t)
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(22229)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govPlugin.GetTallyResult(txHashArr[0], evm.StateDB)
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

func TestGovPlugin_twoProposalsSuccess(t *testing.T) {

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

	lastBlockNumber = uint64(21999)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22000)

	beginBlock(t)

	sndb.Commit(lastBlockHash)

	//buildSnapDBDataCommitted(20001, 22229)
	sndb.Compaction()
	lastBlockNumber = uint64(22229)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govPlugin.GetTallyResult(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else {
		t.Logf("the result status, %s", result.Status.ToString())
	}

	result, err = govPlugin.GetTallyResult(txHashArr[1], evm.StateDB)
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

func TestGovPlugin_versionProposalSuccess(t *testing.T) {

	defer setup(t)()

	submitVersion(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction() //flush to LevelDB

	buildBlockNoCommit(2)
	allVote(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(21999)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22000)

	beginBlock(t)

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(22229)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	lastBlockNumber = uint64(23479)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	//buildBlockNoCommit(23480)
	build_staking_data_more(23480)

	endBlock(t)
	sndb.Commit(lastBlockHash)

	activeVersion := govPlugin.GetActiveVersion(evm.StateDB)
	if activeVersion == uint32(1<<16|1<<8|1) {
		t.Logf("active SUCCESS, %d", activeVersion)
	} else {
		t.Fatalf("active FALSE, %d", activeVersion)
	}
}

func TestGovPlugin_Param(t *testing.T) {
	defer setup(t)()

	paraMap := make(map[string]interface{})
	paraMap["param1"] = 12
	paraMap["param2"] = "stringValue"
	paraMap["param3"] = 12.5

	if err := govPlugin.SetParam(paraMap, evm.StateDB); err != nil {
		t.Errorf("set param failed, %s", err.Error())
		return
	}

	list, err := govPlugin.ListParam(evm.StateDB)
	if err != nil {
		t.Fatalf("list param failed, %s", err)
	} else {
		t.Logf("list size: %d", len(list))
	}

	value, err := govPlugin.GetParamValue("param3", evm.StateDB)
	if err != nil {
		t.Fatalf("get param failed, %s", err)
	} else {
		t.Logf("param name: %s, value: %2.2f", "param3", value.(float64))
	}

}

func TestGovPlugin_ParamProposalSuccess(t *testing.T) {
	defer setup(t)()

	paraMap := make(map[string]interface{})
	paraMap["param1"] = 12
	paraMap["param2"] = "stringValue"
	paraMap["param3"] = 12.5

	if err := govPlugin.SetParam(paraMap, evm.StateDB); err != nil {
		t.Fatalf("set param failed, %s", err)
	}
	submitParam(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction() //flush to LevelDB

	buildBlockNoCommit(2)

	allVote(t, txHashArr[0])

	sndb.Commit(blockHash)
	//buildSnapDBDataCommitted(2, 19999)
	sndb.Compaction()
	lastBlockNumber = uint64(21999)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22000)
	beginBlock(t)
	sndb.Commit(lastBlockHash)

	//buildSnapDBDataCommitted(20001, 22229)
	sndb.Compaction()
	lastBlockNumber = uint64(22229)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govPlugin.GetTallyResult(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.Pass {
		t.Logf("the result status, %d", result.Status)

		value, err := govPlugin.GetParamValue("param3", evm.StateDB)
		if err != nil {
			t.Fatalf("cannot find the param value, %s", err.Error())
		} else {
			t.Logf("the param value, %2.2f", value.(float64))
		}

	} else {
		t.Logf("the result status error, %d", result.Status)
	}
}
