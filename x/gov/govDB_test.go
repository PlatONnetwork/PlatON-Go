package gov_test

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"math/big"
	"testing"
)

//var _ = checker.Suite(&StateSuite{})
//
//func TestGovDB_SetProposal(t *testing.T) {
//	proposal:= getTProposal()
//
//
//}

func getGovDB() (*gov.GovDB, *state.StateDB) {
	db := ethdb.NewMemDatabase()
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	return gov.NewGovDB(), statedb
}

func getProposal() gov.TextProposal {
	return gov.TextProposal{
		common.Hash{0x01},
		"p#01",
		gov.Version,
		"up,up,up....",
		"哈哈哈哈哈哈",
		"em。。。。",
		big.NewInt(1000),
		big.NewInt(10000000),
		discover.NodeID{},
		gov.TallyResult{},
	}
}

func TestGovDB_SetProposal(t *testing.T) {
	//db, statedb := getGovDB()
	//
	////db := ethdb.NewMemDatabase()
	////statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	//
	//proposal := getProposal()
	//if e := db.SetProposal(proposal, statedb); e != nil {
	//	t.Errorf("set proposal error,%s", e.Error())
	//}
	//
	////var proposalGet  Proposal
	//if proposalGet, e := db.GetProposal(proposal.ProposalID, statedb); e != nil {
	//	t.Errorf("get proposal error,%s", e.Error())
	//} else {
	//	fmt.Printf(proposalGet.GetUrl())
	//}

}

func TestGovDB_ListVote(t *testing.T) {

}

func Test_Nothing(t *testing.T) {
	proposal := getProposal()
	proposalBytes, _ := json.Marshal(proposal)

	fmt.Printf("%s \n", hex.EncodeToString(proposalBytes))
	proposalBytes = append(proposalBytes, byte(proposal.GetProposalType()))

	fmt.Printf("%s \n", hex.EncodeToString(proposalBytes))

	var txp gov.TextProposal
	json.Unmarshal(proposalBytes, &txp)
	fmt.Println(txp.String())

}
