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
	"math/big"
)

func getHeader() types.Header {
	return types.Header{
		Number: big.NewInt(100),
	}
}

func getVerProposal() gov.VersionProposal {
	return gov.VersionProposal{
		gov.TextProposal{
			common.Hash{0x01},
			"p#01",
			gov.Version,
			"up,up,up....",
			"哈哈哈哈哈哈",
			"em。。。。",
			uint64(1000),
			uint64(10000),
			discover.NodeID{},
			gov.TallyResult{},
		},
		3200000,
		uint64(11250),
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

//func TestGovPlugin_BeginBlock(t *testing.T) {
//	_, statedb := GetGovDB()
//	blockhash := common.HexToHash("11")
//	header := getHeader()
//	plugin.GovPluginInstance().BeginBlock(blockhash, &header, statedb)
//
//}
//
//func TestGovPlugin_EndBlock(t *testing.T) {
//	_, statedb := GetGovDB()
//	blockhash := common.HexToHash("11")
//	header := getHeader()
//	plugin.GovPluginInstance().EndBlock(blockhash, &header, statedb)
//}

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

//
//func TestGovPlugin_DeclareVersion(t *testing.T) {
//
//}
