package plugin_test

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"testing"
	//"fmt"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
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

func GetGovDB() (*gov.GovDB, *state.StateDB) {
	db := ethdb.NewMemDatabase()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	return gov.GovDBInstance(), statedb
}

func TestGovPlugin_BeginBlock(t *testing.T) {
	_, statedb := GetGovDB()
	blockhash := common.HexToHash("11")
	header := getHeader()
	plugin.GovPluginInstance().BeginBlock(blockhash, &header, statedb)

}

func TestGovPlugin_EndBlock(t *testing.T) {
	_, statedb := GetGovDB()
	blockhash := common.HexToHash("11")
	header := getHeader()
	plugin.GovPluginInstance().EndBlock(blockhash, &header, statedb)
}

func TestGovPlugin_Submit(t *testing.T) {

	_, statedb := GetGovDB()

	sender := common.HexToAddress("0x11")
	blockhash := common.HexToHash("0x11")

	err := plugin.GovPluginInstance().Submit(99, sender, getVerProposal(), blockhash, statedb)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}
}

func TestGovPlugin_Vote(t *testing.T) {


	proposalID := common.Hash{0x01}
	node := discover.NodeID{0x11}

	v := gov.Vote{
		proposalID,
		node,
		gov.Yes,
	}
	sender := common.HexToAddress("0x11")
	blockhash := common.HexToHash("0x11")
	_, statedb := GetGovDB()

	//submit
	err := plugin.GovPluginInstance().Submit(99, sender, getVerProposal(), blockhash, statedb)
	if err != nil {
		t.Fatalf("submit err: %s", err)
	}

	err = plugin.GovPluginInstance().Vote(sender, v, blockhash, 12, statedb)
	if err != nil {
		t.Fatalf("vote err: %s.", err)
	}
}

func TestGovPlugin_DeclareVersion(t *testing.T) {

}
