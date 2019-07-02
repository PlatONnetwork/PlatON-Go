package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"testing"
	"fmt"
)


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
			uint64(10000000),
			discover.NodeID{},
			gov.TallyResult{},
		},
		32,
		uint64(562222),
	}

}


func getGovDB() (*gov.GovDB, *state.StateDB) {
	db := ethdb.NewMemDatabase()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	return gov.GovDBInstance(), statedb
}

//func (govPlugin *GovPlugin) Submit(curBlockNum uint64, from common.Address, proposal gov.Proposal, blockHash common.Hash, state xcom.StateDB) error {

func TestGovPlugin_BeginBlock(t *testing.T) {
	_, statedb := getGovDB()
	sender := common.HexToAddress("0x11")
	blockhash := common.HexToHash("0x11")

	GovPluginInstance().Submit(10,sender,getVerProposal(),blockhash,statedb)
	fmt.Println("TestHello")
}

func TestGovPlugin_EndBlock(t *testing.T) {

}

func TestGovPlugin_Submit(t *testing.T) {


}

func TestGovPlugin_Vote(t *testing.T) {

}


func TestWorld(t *testing.T) {
	fmt.Println("TestWorld")
}
