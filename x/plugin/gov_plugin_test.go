package plugin_test

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
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
	t.Log("setup()......")

	state, _ := newChainState()
	evm = newEvm(blockNumber, blockHash, state)

	newPlugins()

	govPlugin = plugin.GovPluginInstance()

	build_staking_data()

	snapdb = snapshotdb.Instance()

	return func() {
		t.Log("tear down()......")
		snapdb.Clear()
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
	submitText(t)
	//submitVersion(t)
}


func TestGovPlugin_Vote(t *testing.T) {
	defer setup(t)()
	submitVersion(t)

	sndb.Commit(blockHash)

	InitPlatONPluginTestData()

	buildSnapDBDataNoCommit(2)

	v := gov.Vote{
		txHashArr[0],
		nodeIdArr[0],
		gov.Yes,
	}

	err := plugin.GovPluginInstance().Vote(sender, v, lastBlockHash, 2, evm.StateDB)
	if err != nil {
		t.Fatalf("vote err: %s.", err)
	}
}

func TestGovPlugin_DeclareVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t)

	sndb.Commit(blockHash)

	InitPlatONPluginTestData()

	buildSnapDBDataNoCommit(2)

	err := plugin.GovPluginInstance().DeclareVersion(sender, nodeIdArr[0], getVerProposal().NewVersion, lastBlockHash, 2, evm.StateDB)
	if err != nil {
		t.Fatalf("Declare Version err ...%s", err)
	}
}

func TestGovPlugin_ListProposal(t *testing.T) {

	defer setup(t)()

	submitText(t)
	//submitVersion(t)

	pList, err := plugin.GovPluginInstance().ListProposal(lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("List Proposal err ...%s", err)
	}else {
		t.Logf("proposal count, %d", len(pList))
	}

}


func submitText(t *testing.T) {
	vp := gov.TextProposal{
		ProposalID:		txHashArr[0],
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

	err := plugin.GovPluginInstance().Submit(blockNumber.Uint64(), sender, vp, lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}
}


func submitVersion(t *testing.T) {
	vp := gov.VersionProposal{
		ProposalID:		txHashArr[0],
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

	err := plugin.GovPluginInstance().Submit(blockNumber.Uint64(), sender, vp, lastBlockHash, evm.StateDB)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}
}

func allVote(t *testing.T) {
	for _, nodeID := range nodeIdArr {
		vote := gov.Vote{
			ProposalID:		txHashArr[0],
			VoteNodeID:		nodeID,
			VoteOption:		gov.Yes,
		}
		err := plugin.GovPluginInstance().Vote(sender, vote, blockHash, 1, evm.StateDB)
		if err != nil {
			t.Fatalf("vote err: %s.", err)
		}
	}
}

func beginBlock(t *testing.T) {
	err := plugin.GovPluginInstance().BeginBlock(lastBlockHash, &lastHeader, evm.StateDB)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
}

func endBlock(t *testing.T) {
	err := plugin.GovPluginInstance().EndBlock(lastBlockHash, &lastHeader, evm.StateDB)
	if err != nil {
		t.Fatalf("end block err... %s", err)
	}
}

func TestGovPlugin_versionProposalSuccess(t *testing.T) {

	defer setup(t)()

	InitPlatONPluginTestData()

	submitVersion(t)

	allVote(t)

	sndb.Commit(blockHash)

	buildSnapDBDataCommitted(2, 19999)

	buildSnapDBDataNoCommit(20000)

	beginBlock(t)

	sndb.Commit(lastBlockHash)

	buildSnapDBDataCommitted(20001, 22229)

	buildSnapDBDataNoCommit(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	buildSnapDBDataCommitted(22231, 23479)

	buildSnapDBDataNoCommit(23480)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	activeVersion := gov.GovDBInstance().GetActiveVersion(evm.StateDB)

	if activeVersion == uint32(1<<16 | 1<<8 | 1) {
		t.Log("SUCCESS")
	}else{
		t.Log("FALSE")
	}
}


func TestGovPlugin_textProposalSuccess(t *testing.T) {

	defer setup(t)()

	InitPlatONPluginTestData()

	submitText(t)

	allVote(t)

	sndb.Commit(blockHash)

	buildSnapDBDataCommitted(2, 19999)

	buildSnapDBDataNoCommit(20000)

	beginBlock(t)

	sndb.Commit(lastBlockHash)

	buildSnapDBDataCommitted(20001, 22229)

	buildSnapDBDataNoCommit(22230)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := gov.GovDBInstance().GetTallyResult(txHashArr[0], evm.StateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Error("cannot find the tally result")
	}else{
		t.Logf("the result status, %d", result.Status)
	}
}