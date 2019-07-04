package plugin_test

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
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
		"哈哈哈哈哈哈",
		"em。。。。",
		uint64(1000),
		uint64(10000),
		discover.NodeID{0x11},
		gov.TallyResult{},
		3200000,
		uint64(11250),
	}
}

func getTxtProposal() gov.TextProposal {
	return gov.TextProposal{
		common.Hash{0x01},
		"p#01",
		gov.Text,
		"up,up,up....",
		"This is an example...",
		"em。。。。",
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
	_, statedb := GetGovDB()
	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	blockhash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}
	header := getHeader()
	_, err := plugin.GovPluginInstance().BeginBlock(blockhash, &header, statedb)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
}

func TestGovPlugin_EndBlock(t *testing.T) {
	_, statedb := GetGovDB()
	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	blockhash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}
	header := getHeader()
	_, err := plugin.GovPluginInstance().EndBlock(blockhash, &header, statedb)
	if err != nil {
		t.Fatalf("end block err... %s", err)
	}
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
	_, state := GetGovDB()
	newVersion := uint32(1792)

	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	blockhash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}

	header := getHeader()
	_, err := plugin.GovPluginInstance().BeginBlock(blockhash, &header, state)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}

	err = plugin.GovPluginInstance().DeclareVersion(sender, node, newVersion, blockhash, 1, state)
	if err != nil {
		t.Fatalf("Declare Version err ...%s", e)
	}
}

func TestGovPlugin_ListProposal(t *testing.T) {

	_, statedb := GetGovDB()
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
}

func TestGovPlugin_TestVersionTally(t *testing.T) {

	sender := common.HexToAddress("0x11")

	db, statedb := GetGovDB()
	proposalID := common.Hash{0x01}

	vp := getVerProposal()

	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
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

	_, err = plugin.GovPluginInstance().BeginBlock(blockHash, &header, statedb)
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
}


func TestGovPlugin_TestTextTally(t *testing.T) {

	sender := common.HexToAddress("0x11")

	db, statedb := GetGovDB()
	proposalID := common.Hash{0x01}

	tp := getTxtProposal()

	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
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

	_, err = plugin.GovPluginInstance().BeginBlock(blockHash, &header, statedb)
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
}