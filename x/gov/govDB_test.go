package gov_test

import (
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
	return gov.GovDBInstance(), statedb
}

func TestGovDB_SetGetTxtProposal(t *testing.T) {
	db, statedb := getGovDB()

	proposal := getTxtProposal()
	if e := db.SetProposal(proposal, statedb); e != nil {
		t.Errorf("set proposal error,%s", e)
	}

	//var proposalGet  Proposal
	if proposalGet, e := db.GetProposal(proposal.ProposalID, statedb); e != nil {
		t.Errorf("get proposal error,%s", e)
	} else {
		if proposalGet.GetUrl() != proposal.GetUrl() {
			t.Fatalf("get proposal error,expect %s,get %s", proposal.GetUrl(), proposalGet.GetUrl())
		}
	}
}

func TestGovDB_SetGetVerProposal(t *testing.T) {
	db, statedb := getGovDB()

	proposal := getVerProposal()
	if e := db.SetProposal(proposal, statedb); e != nil {
		t.Errorf("set proposal error,%s", e)
	}

	//var proposalGet  Proposal
	if proposalGet, e := db.GetProposal(proposal.ProposalID, statedb); e != nil {
		t.Errorf("get proposal error,%s", e)
	} else {
		if proposalGet.GetUrl() != proposal.GetUrl() {
			t.Fatalf("get proposal error,expect %s,get %s", proposal.GetUrl(), proposalGet.GetUrl())
		}
	}
}

func TestGovDB_ListProposalID(t *testing.T) {
	//db, statedb := getGovDB()

}

func TestGovDB_SetVote(t *testing.T) {
	db, statedb := getGovDB()
	proposal := getVerProposal()

	db.SetProposal(proposal, statedb)

	for _, node := range nodeIdTests {
		if !db.SetVote(proposal.ProposalID, node.Voter, node.Option, statedb) {
			t.Fatalf("set vote error...")
		}
	}

	voteList := db.ListVote(proposal.GetProposalID(), statedb)

	if len(voteList) != len(nodeIdTests) {
		t.Fatalf("get vote list error, expect count：%d,get count:%d", len(nodeIdTests), len(voteList))
	}

	tallyResult := gov.TallyResult{
		proposal.GetProposalID(),
		uint16(len(voteList)),
		0,
		0,
		1000,
		gov.Pass,
	}

	if !db.SetTallyResult(tallyResult, statedb) {
		t.Fatalf("set vote result error")
	}

	if result, e := db.GetTallyResult(proposal.ProposalID, statedb); e != nil {
		t.Fatalf("get vote result error,%s", e)
	} else {
		if result.Status != tallyResult.Status {
			t.Fatalf("get vote result error")
		}
	}
}

func getTxtProposal() gov.TextProposal {
	return gov.TextProposal{
		common.Hash{0x01},
		"p#01",
		gov.Text,
		"up,up,up....",
		"哈哈哈哈哈哈",
		"em。。。。",
		big.NewInt(1000),
		big.NewInt(10000000),
		discover.NodeID{},
		gov.TallyResult{},
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
			big.NewInt(1000),
			big.NewInt(10000000),
			discover.NodeID{},
			gov.TallyResult{},
		},
		32,
		big.NewInt(562222),
	}

}

var nodeIdTests = []gov.VoteValue{
	{
		Voter:  discover.MustHexID("0x1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		Option: gov.Yes,
	},
	{
		Voter:  discover.MustHexID("0x1dd8d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		Option: gov.Yes,
	},
	{
		Voter:  discover.MustHexID("0x1dd7d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		Option: gov.Yes,
	},
	{
		Voter:  discover.MustHexID("0x1dd6d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		Option: gov.Yes,
	},
	{
		Voter:  discover.MustHexID("0x1dd5d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		Option: gov.Yes,
	},
}
