package plugin_test

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"os"
	"testing"
	//"fmt"
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
)

var (
	snapdb 		snapshotdb.DB
	govPlugin	*plugin.GovPlugin
	evm			*vm.EVM
)


func setup(t *testing.T) func() {

	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	t.Log("setup()......")

	state, _ := newChainState()
	evm = newEvm(blockNumber, blockHash, state)

	newPlugins()

	govPlugin = govPlugin

	build_staking_data()

	snapdb = snapshotdb.Instance()

	return func() {
		t.Log("tear down()......")
		snapdb.Clear()
		ClearPlatONPluginTestData()
	}
}

func getHeader() types.Header {
	size := int64(xcom.ConsensusSize * xcom.EpochSize)
	return types.Header{
		Number: big.NewInt(size),
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

func getTxtProposal() gov.TextProposal {
	return gov.TextProposal{
		common.Hash{0x02},
		"p#02",
		gov.Text,
		"up,up,up....",
		"text proposal example",
		"http://url",
		uint64(1000),
		uint64(10000000),
		discover.NodeID{},
		gov.TallyResult{},
	}
}

func newblock(snapdb snapshotdb.DB, blockNumber *big.Int) (common.Hash, error) {

	recognizedHash := generateHash("recognizedHash")

	commitHash := recognizedHash
	if err := snapdb.NewBlock(blockNumber, common.Hash{}, commitHash); err != nil {
		return common.Hash{}, err
	}

	if err := snapdb.Put(commitHash, []byte("wu"), []byte("wei")); err != nil {
		return common.Hash{}, err
	}

	get, err := snapdb.Get(commitHash, []byte("wu"))
	if err != nil {
		return common.Hash{}, err
	}
	fmt.Printf("get result :%s", get)

	return commitHash, nil
}

func generateHash(n string) common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(n))
	return rlpHash(buf.Bytes())
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func commitBlock(snapdb snapshotdb.DB, blockhash common.Hash) error {
	return snapdb.Commit(blockhash)
}

func GetGovDB() (*gov.GovDB, *state.StateDB) {
	db := ethdb.NewMemDatabase()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	return gov.GovDBInstance(), statedb
}


func TestGovPlugin_Submit(t *testing.T) {
	defer setup(t)()
	submitText(t, txHashArr[0])
	//submitVersion(t)
}


func TestGovPlugin_Vote(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(blockHash)

	InitPlatONPluginTestData()

	buildSnapDBDataNoCommit(2)

	v := gov.Vote{
		txHashArr[0],
		nodeIdArr[0],
		gov.Yes,
	}

	err := govPlugin.Vote(sender, v, lastBlockHash, 2, evm.StateDB)
	if err != nil {
		t.Fatalf("vote err: %s.", err)
	}
}

func TestGovPlugin_DeclareVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(blockHash)

	InitPlatONPluginTestData()

	buildSnapDBDataNoCommit(2)

	err := govPlugin.DeclareVersion(sender, nodeIdArr[0], getVerProposal().NewVersion, lastBlockHash, 2, evm.StateDB)
	if err != nil {
		t.Fatalf("Declare Version err ...%s", err)
	}
}

func TestGovPlugin_ListProposal(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])
	//submitVersion(t)

	pList, err := govPlugin.ListProposal(lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("List Proposal err ...%s", err)
	}else {
		t.Logf("proposal count, %d", len(pList))
	}

}


func submitText(t *testing.T, pid common.Hash) {
	vp := gov.TextProposal{
		ProposalID:		pid,
		GithubID:		"githubID",
		ProposalType:	gov.Text,
		Topic:			"versionTopic",
		Desc: 			"versionDesc",
		Url:			"versionUrl",
		EndVotingBlock:	uint64(22230),
		Proposer:		nodeIdArr[0],
	}

	state := evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := govPlugin.Submit(blockNumber.Uint64(), sender, vp, lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}
}


func submitVersion(t *testing.T, pid common.Hash) {
	vp := gov.VersionProposal{
		ProposalID:		pid,
		GithubID:		"githubID",
		ProposalType:	gov.Version,
		Topic:			"versionTopic",
		Desc: 			"versionDesc",
		Url:			"versionUrl",
		EndVotingBlock:	uint64(22230),
		Proposer:		nodeIdArr[0],
		NewVersion:		uint32(1<<16 | 1<<8 | 1),
		ActiveBlock:	uint64(23480),
	}

	state := evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], blockHash, 0)

	err := govPlugin.Submit(blockNumber.Uint64(), sender, vp, lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}
}

func submitParam(t *testing.T, pid common.Hash) {
	vp := gov.ParamProposal{
		ProposalID:		pid,
		GithubID:		"githubID",
		ProposalType:	gov.Text,
		Topic:			"versionTopic",
		Desc: 			"versionDesc",
		Url:			"versionUrl",
		EndVotingBlock:	uint64(22230),
		Proposer:		nodeIdArr[0],

		ParamName: 		"param3",
		CurrentValue:   12.5,
		NewValue:		0.85,
	}

	state := evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := govPlugin.Submit(blockNumber.Uint64(), sender, vp, lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}
}

func allVote(t *testing.T, pid common.Hash) {
	for _, nodeID := range nodeIdArr {
		vote := gov.Vote{
			ProposalID:		pid,
			VoteNodeID:		nodeID,
			VoteOption:		gov.Yes,
		}
		err := govPlugin.Vote(sender, vote, blockHash, 1, evm.StateDB)
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

func TestGovPlugin_textProposalSuccess(t *testing.T) {

	defer setup(t)()

	InitPlatONPluginTestData()

	submitText(t, txHashArr[0])

	allVote(t, txHashArr[0])
	sndb.Commit(blockHash)

	//buildSnapDBDataCommitted(2, 19999)
	sndb.Compaction()
	lastBlockNumber = uint64(19999)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))


	buildSnapDBDataNoCommit(20000)
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


	buildSnapDBDataNoCommit(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govPlugin.GetTallyResult(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Error("cannot find the tally result")
	}else if result.Status == gov.Pass{
		t.Logf("the result status, %d", result.Status)
	}else {
		t.Logf("the result status error, %d", result.Status )
	}
}

func TestGovPlugin_twoProposalsSuccess(t *testing.T) {

	defer setup(t)()

	InitPlatONPluginTestData()

	submitText(t, txHashArr[0])
	allVote(t, txHashArr[0])

	submitVersion(t, txHashArr[1])
	allVote(t, txHashArr[1])


	sndb.Commit(blockHash)

	//buildSnapDBDataCommitted(2, 19999)
	sndb.Compaction()
	lastBlockNumber = uint64(19999)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	buildSnapDBDataNoCommit(20000)

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


	buildSnapDBDataNoCommit(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govPlugin.GetTallyResult(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Error("cannot find the tally result")
	}else{
		t.Logf("the result status, %d", result.Status)
	}

	result, err = govPlugin.GetTallyResult(txHashArr[1], evm.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Error("cannot find the tally result")
	}else if result.Status == gov.PreActive{
		t.Logf("the result status, %d", result.Status)
	}else {
		t.Logf("the result status error, %d", result.Status )
	}
}

func TestGovPlugin_versionProposalSuccess(t *testing.T) {

	defer setup(t)()

	InitPlatONPluginTestData()

	submitVersion(t, txHashArr[0])

	allVote(t, txHashArr[0])

	sndb.Commit(blockHash)

	//buildSnapDBDataCommitted(2, 19999)
	sndb.Compaction()
	lastBlockNumber = uint64(19999)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	buildSnapDBDataNoCommit(20000)

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


	buildSnapDBDataNoCommit(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	//buildSnapDBDataCommitted(22231, 23479)
	sndb.Compaction()
	lastBlockNumber = uint64(23479)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	buildSnapDBDataNoCommit(23480)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	activeVersion := govPlugin.GetActiveVersion(evm.StateDB)
	t.Logf("active version, %d", activeVersion)
	if activeVersion == uint32(1<<16 | 1<<8 | 1) {
		t.Logf("active SUCCESS, %d", activeVersion)
	}else{
		t.Errorf("active FALSE, %d", activeVersion)
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
		t.Errorf("list param failed, %s", err.Error())
		return
	}else {
		t.Logf("list size: %d", len(list))
	}


	value, err := govPlugin.GetParamValue("param3", evm.StateDB )
	if err != nil {
		t.Errorf("get param failed, %s", err.Error())
		return
	}else {
		t.Logf("param name: %s, value: %2.2f", "param3", value.(float64))
	}

}


func TestGovPlugin_ParamProposalSuccess(t *testing.T) {

	defer setup(t)()

	InitPlatONPluginTestData()

	paraMap := make(map[string]interface{})
	paraMap["param1"] = 12
	paraMap["param2"] = "stringValue"
	paraMap["param3"] = 12.5

	if err := govPlugin.SetParam(paraMap, evm.StateDB); err != nil {
		t.Errorf("set param failed, %s", err.Error())
		return
	}

	submitParam(t, txHashArr[0])

	allVote(t, txHashArr[0])
	sndb.Commit(blockHash)

	//buildSnapDBDataCommitted(2, 19999)
	sndb.Compaction()
	lastBlockNumber = uint64(19999)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))


	buildSnapDBDataNoCommit(20000)
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


	buildSnapDBDataNoCommit(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := govPlugin.GetTallyResult(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Error("cannot find the tally result")
	}else if result.Status == gov.Pass{
		t.Logf("the result status, %d", result.Status)

		value, err := govPlugin.GetParamValue("param3", evm.StateDB)
		if err != nil {
			t.Error("cannot find the param value, %s", err.Error())
			return
		}else{
			t.Logf("the param value, %2.2f", value.(float64))
		}

	}else {
		t.Logf("the result status error, %d", result.Status )
	}
}
