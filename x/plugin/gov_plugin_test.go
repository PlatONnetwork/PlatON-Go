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
	"github.com/PlatONnetwork/PlatON-Go/params"
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

func TestGovPlugin_BeginBlock(t *testing.T) {
	db, statedb := GetGovDB()
	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	blockhash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}
	header := getHeader()
	err := plugin.GovPluginInstance().BeginBlock(blockhash, &header, statedb)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
	db.Reset()
}

func TestGovPlugin_EndBlock(t *testing.T) {
	db, statedb := GetGovDB()
	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	blockhash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}
	header := getHeader()
	err := plugin.GovPluginInstance().EndBlock(blockhash, &header, statedb)
	if err != nil {
		t.Fatalf("end block err... %s", err)
	}
	db.Reset()
}

func TestGovPlugin_Submit(t *testing.T) {

	db, statedb := GetGovDB()
	sender := common.HexToAddress("0x11")

	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()

	//create block
	blockhash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}

	err := plugin.GovPluginInstance().Submit(99, sender, getVerProposal(), blockhash, statedb)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}
	db.Reset()
}


func TestGovPlugin_Vote(t *testing.T) {

	proposalID := common.Hash{0x01}
	node := discover.NodeID{0x11}

	db, statedb := GetGovDB()

	v := gov.Vote{
		proposalID,
		node,
		gov.Yes,
	}
	sender := common.HexToAddress("0x11")
	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	if plugin.GovPluginInstance().GetActiveVersion(statedb) == 0 {
		defaultVersion := uint32(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch)
		db.SetActiveVersion(defaultVersion, statedb)
	}
	blockhash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}

	//submit
	err := plugin.GovPluginInstance().Submit(99, sender, getVerProposal(), blockhash, statedb)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}
	err = plugin.GovPluginInstance().Vote(sender, v, blockhash, 1, statedb)
	if err != nil {
		t.Fatalf("vote err: %s.", err)
	}

	if err := commitBlock(snapdb, blockhash); err != nil {
		t.Fatalf("commit block error..%s", err)
	}

	db.Reset()
}

func TestGovPlugin_DeclareVersion(t *testing.T) {
	//func (govPlugin *GovPlugin) DeclareVersion(from common.Address, declaredNodeID discover.NodeID, version uint32, blockHash common.Hash, curBlockNum uint64, state xcom.StateDB) error {
	sender := common.HexToAddress("0x11")
	node := discover.NodeID{0x11}
	db, state := GetGovDB()
	//newVersion := uint32(1792)

	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	if plugin.GovPluginInstance().GetActiveVersion(state) == 0 {
		defaultVersion := uint32(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch)
		db.SetActiveVersion(defaultVersion, state)
	}
	blockhash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}

	header := getHeader()
	err := plugin.GovPluginInstance().BeginBlock(blockhash, &header, state)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
	err = plugin.GovPluginInstance().Submit(99, sender, getVerProposal(), blockhash, state)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}

	err = plugin.GovPluginInstance().DeclareVersion(sender, node, getVerProposal().NewVersion, blockhash, 1, state)
	if err != nil {
		t.Fatalf("Declare Version err ...%s", e)
	}
	db.Reset()
}

func TestGovPlugin_ListProposal(t *testing.T) {

	db, statedb := GetGovDB()
	sender := common.HexToAddress("0x11")

	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	blockHash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}

	err := plugin.GovPluginInstance().Submit(99, sender, getVerProposal(), blockHash, statedb)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}

	_, err = plugin.GovPluginInstance().ListProposal(blockHash, statedb)
	if err != nil {
		t.Fatalf("List Proposal err ...%s", e)
	}
	db.Reset()
}

func TestGovPlugin_TestVersionTally(t *testing.T) {

	sender := common.HexToAddress("0x11")

	db, statedb := GetGovDB()
	proposalID := common.Hash{0x01}

	vp := getVerProposal()

	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	if plugin.GovPluginInstance().GetActiveVersion(statedb) == 0 {
		defaultVersion := uint32(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch)
		db.SetActiveVersion(defaultVersion, statedb)
	}
	blockHash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}

	node := discover.NodeID{0x11}
	v := gov.Vote{
		proposalID,
		node,
		gov.Yes,
	}

	header := getHeader()

	//submit
	err := plugin.GovPluginInstance().Submit(99, sender, getVerProposal(), blockHash, statedb)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}

	err = plugin.GovPluginInstance().BeginBlock(blockHash, &header, statedb)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
	err = plugin.GovPluginInstance().Vote(sender, v, blockHash, 1, statedb)
	if err != nil {
		t.Fatalf("vote err: %s.", err)
	}
	votedList, err := db.ListVotedVerifier(proposalID, statedb)
	//tallyForVersionProposal(votedVerifierList []discover.NodeID, accuCnt uint16, proposal gov.VersionProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {
	err = plugin.GovPluginInstance().TestTally(votedList, 1, vp, blockHash, 1, statedb)
	if err != nil {
		t.Fatalf("Test Tally ...%s", err)
	}
	db.Reset()
}


func TestGovPlugin_TestTextTally(t *testing.T) {

	sender := common.HexToAddress("0x11")

	db, statedb := GetGovDB()
	proposalID := common.Hash{0x01}

	tp := getTxtProposal()

	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	if plugin.GovPluginInstance().GetActiveVersion(statedb) == 0 {
		defaultVersion := uint32(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch)
		db.SetActiveVersion(defaultVersion, statedb)
	}
	blockHash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}

	node := discover.NodeID{0x11}
	v := gov.Vote{
		proposalID,
		node,
		gov.No,
	}

	header := getHeader()

	//submit
	err := plugin.GovPluginInstance().Submit(99, sender, getVerProposal(), blockHash, statedb)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}

	err = plugin.GovPluginInstance().BeginBlock(blockHash, &header, statedb)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
	err = plugin.GovPluginInstance().Vote(sender, v, blockHash, 1, statedb)
	if err != nil {
		t.Fatalf("vote err: %s.", err)
	}
	votedList, err := db.ListVotedVerifier(proposalID, statedb)
	//tallyForVersionProposal(votedVerifierList []discover.NodeID, accuCnt uint16, proposal gov.VersionProposal, blockHash common.Hash, blockNumber uint64, state xcom.StateDB) error {
	err = plugin.GovPluginInstance().TestTally(votedList, 1, tp, blockHash, 1, statedb)
	if err != nil {
		t.Fatalf("Test Tally ...%s", err)
	}
	db.Reset()
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

	err := plugin.GovPluginInstance().Submit(blockNumber.Uint64(), sender, vp, blockHash, evm.StateDB)
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

func TestGovPlugin_activeSuccess(t *testing.T) {

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