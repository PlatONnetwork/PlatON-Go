package gov_test

import (
	"bytes"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"math/big"
	"testing"
)

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

	proposal := getVerProposal(common.Hash{0x1})
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

func commitBlock(snapdb snapshotdb.DB, blockhash common.Hash) error {
	return snapdb.Commit(blockhash)
}

func TestGovDB_SetProposalT2Snapdb(t *testing.T) {
	db, statedb := getGovDB()

	var proposalIds []common.Hash
	var proposalIdsEnd []common.Hash
	var proposalIdsPre common.Hash

	snapdb := snapshotdb.Instance()
	defer snapdb.Clear()
	//create block
	blockhash, e := newblock(snapdb, big.NewInt(1))
	if e != nil {
		t.Fatalf("create block error ...%s", e)
	}

	totalLen := 10
	for i := 1; i <= totalLen; i++ {
		proposal := getVerProposal(common.Hash{byte(i)})
		if err := db.AddVotingProposalID(blockhash, proposal.ProposalID, statedb); err != nil {
			t.Fatalf("add voting proposal failed...%s", err)
		}
		proposalIds = append(proposalIds, proposal.ProposalID)
		fmt.Printf("add %d into voting..", i)
	}

	for i := 0; i < 2; i++ {
		if err := db.MoveVotingProposalIDToEnd(blockhash, proposalIds[i], statedb); err != nil {
			t.Fatalf("move voting proposal to end failed...%s", err)
		} else {
			proposalIdsEnd = append(proposalIdsEnd, proposalIds[i])
			proposalIds = append(proposalIds[:i], proposalIds[i+1:]...)
		}
	}

	if err := db.MoveVotingProposalIDToPreActive(blockhash, proposalIds[1], statedb); err != nil {
		t.Fatalf("move voting proposal to pre active failed...%s", err)
	} else {
		proposalIdsPre = proposalIds[1]
		proposalIds = append(proposalIds[:1], proposalIds[2:]...)
	}

	//if proposals, e := db.GetProposalList(blockhash, statedb);e != nil{
	//	t.Fatalf("get proposal list error ,%s",e)
	//}else{
	//	if len(proposals)!= totalLen{
	//		t.Fatalf("get proposal list error ,expect len:%d,get len: %d",totalLen, len(proposals))
	//	}
	//}

	if plist, e := db.ListEndProposalID(blockhash, statedb); e != nil {
		t.Fatalf("list end propsal error,%s", e)
	} else {
		if len(plist) != len(proposalIdsEnd) {
			t.Fatalf("get end proposal list error ,expect len:%d,get len: %d", len(proposalIdsEnd), len(plist))
		}
	}
	if plist, e := db.ListVotingProposal(blockhash, statedb); e != nil {
		t.Fatalf("list end propsal error,%s", e)
	} else {
		if len(plist) != len(proposalIds) {
			t.Fatalf("get voting proposal list error ,expect len:%d,get len: %d", len(proposalIds), len(plist))
		}
	}
	if p, e := db.GetPreActiveProposalID(blockhash, statedb); e != nil {
		t.Fatalf("list end propsal error,%s", e)
	} else {
		if p != proposalIdsPre {
			t.Fatalf("get pre-active proposal error ,expect:%d,get: %d", proposalIdsPre, p)
		}
	}

	if err := commitBlock(snapdb, blockhash); err != nil {
		t.Fatalf("commit block error..%s", err)
	}
}

func TestGovDB_SetVote(t *testing.T) {
	db, statedb := getGovDB()
	proposal := getVerProposal(common.Hash{0x1})

	db.SetProposal(proposal, statedb)

	for _, node := range nodeIdTests {
		if nil != db.SetVote(proposal.ProposalID, node.VoteNodeID, node.VoteOption, statedb) {
			t.Fatalf("set vote error...")
		}
	}

	voteList, err := db.ListVoteValue(proposal.GetProposalID(), statedb)
	if err != nil {
		t.Fatalf("get vote list error, expect count：%d,get count:%d", len(nodeIdTests), len(voteList))
	}

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

	if err := db.SetTallyResult(tallyResult, statedb); err != nil {
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
		uint64(1000),
		uint64(10000000),
		discover.NodeID{},
		gov.TallyResult{},
	}
}

func getVerProposal(proposalId common.Hash) gov.VersionProposal {
	return gov.VersionProposal{
		gov.TextProposal{
			proposalId,
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

var nodeIdTests = []gov.VoteValue{
	{
		VoteNodeID: discover.MustHexID("0x1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: gov.Yes,
	},
	{
		VoteNodeID: discover.MustHexID("0x1dd8d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: gov.Yes,
	},
	{
		VoteNodeID: discover.MustHexID("0x1dd7d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: gov.Yes,
	},
	{
		VoteNodeID: discover.MustHexID("0x1dd6d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: gov.Yes,
	},
	{
		VoteNodeID: discover.MustHexID("0x1dd5d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: gov.Yes,
	},
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
