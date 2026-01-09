// Copyright 2021 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package gov

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/sha3"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var (
	// snapdbTest snapshotdb.DB
	txHash = common.HexToHash("0x00000000000000000000000000000000000000886d5ba2d3dfb2e2f6a1814f22")
)

func TestGovDB_SetProposal_GetProposal_text(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	proposal := getTxtProposal()
	if e := gdb.SetProposal(proposal, chain.StateDB); e != nil {
		t.Errorf("set proposal error,%s", e)
	}

	if proposalGet, e := gdb.GetProposal(proposal.ProposalID, chain.StateDB); e != nil {
		t.Errorf("get proposal error,%s", e)
	} else {
		if proposalGet.GetPIPID() != proposal.GetPIPID() {
			t.Fatalf("get proposal error,expect %s,get %s", proposal.GetPIPID(), proposalGet.GetPIPID())
		}
	}
}

func TestGovDB_SetProposal_GetProposal_version(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	proposal := getVerProposal(common.Hash{0x1})
	if e := gdb.SetProposal(proposal, chain.StateDB); e != nil {
		t.Errorf("set proposal error,%s", e)
	}

	//var proposalGet  Proposal
	if proposalGet, e := gdb.GetProposal(proposal.ProposalID, chain.StateDB); e != nil {
		t.Errorf("get proposal error,%s", e)
	} else {
		if proposalGet.GetPIPID() != proposal.GetPIPID() {
			t.Fatalf("get proposal error,expect %s,get %s", proposal.GetPIPID(), proposalGet.GetPIPID())
		}
	}
}

func TestGovDB_SetProposal_GetProposal_Cancel(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	proposal := getCancelProposal()
	if e := gdb.SetProposal(proposal, chain.StateDB); e != nil {
		t.Errorf("set proposal error,%s", e)
	}

	//var proposalGet  Proposal
	if proposalGet, e := gdb.GetProposal(proposal.ProposalID, chain.StateDB); e != nil {
		t.Errorf("get proposal error,%s", e)
	} else {
		if proposalGet.GetPIPID() != proposal.GetPIPID() {
			t.Fatalf("get proposal error,expect %s,get %s", proposal.GetPIPID(), proposalGet.GetPIPID())
		}
	}
}

func TestGovDB_addGovernParam(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)

	value := &ParamValue{"", "initValue", 0}
	if err := gdb.addGovernParam("PPOS", "testName1", "for testing", value, blockHash); err != nil {
		t.Fatalf("addGovernParam error, %s", err)
	}
	if paramList, err := gdb.listGovernParam("PPOS", blockHash); err != nil {
		t.Fatalf("listGovernParam error, %s", err)
	} else {
		assert.Equal(t, 1, len(paramList))
		assert.Equal(t, "initValue", paramList[0].ParamValue.Value)
	}
}

func TestGovDB_findGovernParamValue(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)

	value := &ParamValue{"", "initValue", 0}
	if err := gdb.addGovernParam("PPOS", "testName1", "for testing", value, blockHash); err != nil {
		t.Fatalf("addGovernParam error, %s", err)
	}
	if value, err := gdb.findGovernParamValue("PPOS", "testName1", blockHash); err != nil {
		t.Fatalf("findGovernParamValue error, %s", err)
	} else {
		assert.Equal(t, "initValue", value.Value)
	}
}

func TestGovDB_updateGovernParamValue(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	value := &ParamValue{"", "initValue", 0}
	if err := chain.AddBlockWithSnapDB(false, func(hash common.Hash, header *types.Header, sdb snapshotdb.DB) error {
		if err := gdb.addGovernParam("PPOS", "testName1", "for testing", value, hash); err != nil {
			return err
		}
		if err := gdb.updateGovernParamValue("PPOS", "testName1", "newValue", uint64(10000), hash); err != nil {
			return err
		} else {
			if value, err := gdb.findGovernParamValue("PPOS", "testName1", hash); err != nil {
				return err
			} else {
				assert.Equal(t, "newValue", value.Value)
				return nil
			}
		}
	}, nil, nil); err != nil {
		t.Error(err)
	}
}

func TestGovDB_SetProposal_GetProposal_param(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)

	chain.StateDB.SetTxContext(txHash, 0)

	value := &ParamValue{"", "initValue", 0}
	if err := gdb.addGovernParam("PPOS", "testName1", "for testing", value, blockHash); err != nil {
		t.Fatalf("addGovernParam error, %s", err)
	}

	proposal := getParamProposal()
	if e := gdb.SetProposal(proposal, chain.StateDB); e != nil {
		t.Errorf("set proposal error,%s", e)
	}

	if proposalGet, e := gdb.GetProposal(proposal.ProposalID, chain.StateDB); e != nil {
		t.Errorf("get proposal error,%s", e)
	} else {
		if proposalGet.GetPIPID() != proposal.GetPIPID() {
			t.Fatalf("get proposal error,expect %s,get %s", proposal.GetPIPID(), proposalGet.GetPIPID())
		}
	}
}

func TestGovDB_GetProposalList(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)

	tp := getTxtProposal()
	if err := gdb.SetProposal(tp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, tp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	vp := getVerProposal(common.Hash{0x2})
	if err := gdb.SetProposal(vp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, vp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	cp := getCancelProposal()
	if err := gdb.SetProposal(cp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, cp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	if proposalList, err := gdb.GetProposalList(blockHash, chain.StateDB); err != nil {
		t.Errorf("list proposal error,%s", err)
	} else {
		assert.Equal(t, 3, len(proposalList), "list proposals errors")
	}
}

func TestGovDB_ListVotingProposal(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)

	if err := gdb.AddVotingProposalID(blockHash, common.Hash{0x01}); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	if err := gdb.AddVotingProposalID(blockHash, common.Hash{0x02}); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	if err := gdb.AddVotingProposalID(blockHash, common.Hash{0x04}); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	if idList, err := gdb.ListVotingProposal(blockHash); err != nil {
		t.Errorf("list proposal error,%s", err)
	} else {
		if len(idList) != 3 {
			t.Fatalf("list voting proposal ID error,expect %d,get %d", 3, len(idList))
		}
	}
}

func TestGovDB_ListEndProposalID(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)

	if err := gdb.AddVotingProposalID(blockHash, common.Hash{0x01}); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}
	if err := gdb.MoveVotingProposalIDToEnd(common.Hash{0x01}, blockHash); err != nil {
		t.Errorf("MoveVotingProposalIDToEnd error,%s", err)
	}

	if err := gdb.AddVotingProposalID(blockHash, common.Hash{0x02}); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}
	if err := gdb.MoveVotingProposalIDToEnd(common.Hash{0x02}, blockHash); err != nil {
		t.Errorf("MoveVotingProposalIDToEnd error,%s", err)
	}

	if err := gdb.AddVotingProposalID(blockHash, common.Hash{0x04}); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}
	if err := gdb.MoveVotingProposalIDToEnd(common.Hash{0x04}, blockHash); err != nil {
		t.Errorf("MoveVotingProposalIDToEnd error,%s", err)
	}

	if idList, err := gdb.ListEndProposalID(blockHash); err != nil {
		t.Errorf("list end proposal error,%s", err)
	} else {
		if len(idList) != 3 {
			t.Fatalf("list end proposal ID error,expect %d,get %d", 3, len(idList))
		}
	}
}

func TestGovDB_SetVote_ListVoteValue(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	proposalID := common.Hash{0x03}
	blockHash := newBlock(chain)

	for _, nodeId := range NodeIDList {
		if err := gdb.AddVoteValue(proposalID, nodeId, Yes, blockHash); err != nil {
			t.Errorf("set vote error,%s", err)
		}
	}
	if voteValueList, err := gdb.ListVoteValue(proposalID, blockHash); err != nil {
		t.Errorf("list proposal's vote value error,%s", err)
	} else {
		if len(voteValueList) != len(NodeIDList) {
			t.Fatalf("list proposal error,expect %d,get %d", len(NodeIDList), len(voteValueList))
		}
	}
}

/*
func TestGovDB_ListVotedVerifier(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	proposalID := common.Hash{0x03}

	for _, nodeId := range NodeIDList {
		if err := AddVoteValue(proposalID, nodeId, Yes, statedb); err != nil {
			t.Errorf("set vote error,%s", err)
		}
	}

	if voteValueList, err := ListVotedVerifier(proposalID, statedb); err != nil {
		t.Errorf("list proposal's vote value error,%s", err)
	} else {
		if len(voteValueList) != len(NodeIDList) {
			t.Fatalf("list proposal error,expect %d,get %d", len(NodeIDList), len(voteValueList))
		}
	}
}
*/

func TestGovDB_GetVotedVerifierMap(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	proposalID := common.Hash{0x03}
	blockHash := newBlock(chain)

	for _, nodeId := range NodeIDList {
		if err := gdb.AddVoteValue(proposalID, nodeId, Yes, blockHash); err != nil {
			t.Errorf("set vote error,%s", err)
		}
	}

	if votedMap, err := gdb.GetVotedVerifierMap(proposalID, blockHash); err != nil {
		t.Errorf("get proposal's voted verifier map error,%s", err)
	} else {
		if len(votedMap) != len(NodeIDList) {
			t.Fatalf("list proposal error,expect %d,get %d", len(NodeIDList), len(votedMap))
		}
	}
}

func TestGovDB_SetProposalT2Snapdb(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	var proposalIdsVoting []common.Hash
	var proposalIdsEnd []common.Hash
	var proposalIdPre common.Hash

	snapdbTest := chain.SnapDB

	blockHash := newBlock(chain)

	totalLen := 10
	newVersion := uint32(0)
	for i := 1; i <= totalLen; i++ {
		proposal := getVerProposal(common.Hash{byte(i)})
		newVersion = proposal.NewVersion
		if err := gdb.AddVotingProposalID(blockHash, proposal.ProposalID); err != nil {
			t.Fatalf("add voting proposal failed...%s", err)
		}
		proposalIdsVoting = append(proposalIdsVoting, proposal.ProposalID)

		gdb.SetProposal(proposal, chain.StateDB)
	}

	//-- move voting to end
	for i := 0; i < 2; i++ {
		if err := gdb.MoveVotingProposalIDToEnd(proposalIdsVoting[i], blockHash); err != nil {
			t.Fatalf("move voting proposal to end failed...%s", err)
		} else {
			proposalIdsEnd = append(proposalIdsEnd, proposalIdsVoting[i])
			proposalIdsVoting = append(proposalIdsVoting[:i], proposalIdsVoting[i+1:]...)
		}
	}
	if proposals, err := gdb.GetProposalList(blockHash, chain.StateDB); err != nil {
		t.Fatalf("get proposal list error ,%s", err)
	} else {
		assert.Equal(t, totalLen, len(proposals))
	}
	if plist, err := gdb.ListEndProposalID(blockHash); err != nil {
		t.Fatalf("list end propsal error,%s", err)
	} else {
		assert.Equal(t, len(proposalIdsEnd), len(plist))
	}

	//-- move preactive to end
	proposalIdPre = proposalIdsVoting[0]
	proposalIdsVoting = proposalIdsVoting[1:]
	if err := gdb.MoveVotingProposalIDToPreActive(blockHash, proposalIdPre, 32); err != nil {
		t.Fatalf("move voting proposal to pre active failed...%s", err)
	}

	if plist, err := gdb.GetProposalList(blockHash, chain.StateDB); err != nil {
		t.Fatalf("get proposal list error ,%s", err)
	} else {
		assert.Equal(t, totalLen, len(plist))
	}
	if plist, err := gdb.ListEndProposalID(blockHash); err != nil {
		t.Fatalf("list end propsal error,%s", err)
	} else {
		assert.Equal(t, len(proposalIdsEnd), len(plist))
	}
	if prePID, err := gdb.GetPreActiveProposalID(blockHash); err != nil {
		t.Fatalf("list end propsal error,%s", err)
	} else {
		assert.Equal(t, proposalIdPre, prePID)
	}
	preVersion := gdb.GetPreActiveVersion(blockHash)
	assert.Equal(t, newVersion, preVersion)

	if plist, err := gdb.ListVotingProposal(blockHash); err != nil {
		t.Fatalf("list end propsal error,%s", err)
	} else {
		assert.Equal(t, len(proposalIdsVoting), len(plist))
	}

	snapdbTest.Commit(blockHash)
	blockHash = newBlock(chain)

	// move preactive to end
	proposalIdsEnd = append(proposalIdsEnd, proposalIdPre)
	if err := gdb.MovePreActiveProposalIDToEnd(blockHash, proposalIdPre); err != nil {
		t.Fatalf("move preactive proposal id to end failed...%s", err)
	}
	if prePID, err := gdb.GetPreActiveProposalID(blockHash); err != nil {
		t.Fatalf("list end propsal error,%s", err)
	} else {
		assert.Equal(t, common.ZeroHash, prePID)
	}
	preVersion = gdb.GetPreActiveVersion(blockHash)
	assert.Equal(t, uint32(0), preVersion)

	if plist, err := gdb.GetProposalList(blockHash, chain.StateDB); err != nil {
		t.Fatalf("get proposal list error ,%s", err)
	} else {
		assert.Equal(t, totalLen, len(plist))
	}
	if plist, err := gdb.ListEndProposalID(blockHash); err != nil {
		t.Fatalf("list end propsal error,%s", err)
	} else {
		assert.Equal(t, len(proposalIdsEnd), len(plist))
	}
	if plist, err := gdb.ListVotingProposal(blockHash); err != nil {
		t.Fatalf("list end propsal error,%s", err)
	} else {
		assert.Equal(t, len(proposalIdsVoting), len(plist))
	}
}

func TestGovDB_SetPreActiveVersion(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	blockHash := newBlock(chain)
	version := uint32(32)
	if err := gdb.SetPreActiveVersion(blockHash, version); err != nil {
		t.Fatalf("set pre-active version error...%s", err)
	}
	vget := gdb.GetPreActiveVersion(blockHash)
	if vget != version {
		t.Fatalf("get pre-active version error,expect version:%d,get version:%d", version, vget)
	}
}

func TestGovDB_GetPreActiveVersionNotExist(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()
	gdb := NewGovDB(chain.SnapDB)
	blockHash := newBlock(chain)
	vget := gdb.GetPreActiveVersion(blockHash)
	t.Logf("get pre-active version error,get version:%d", vget)
}

func TestGovDB_AddActiveVersion(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	gov := NewGov(chain.SnapDB)

	version := uint32(32)
	//proposal := getVerProposal(common.Hash{0x1})
	if err := gdb.AddActiveVersion(version, 10000, chain.StateDB); err != nil {
		t.Fatalf("add active version error...%s", err)
	}

	version = uint32(33)
	//proposal := getVerProposal(common.Hash{0x1})
	if err := gdb.AddActiveVersion(version, 20000, chain.StateDB); err != nil {
		t.Fatalf("add active version error...%s", err)
	}

	vget := gov.GetCurrentActiveVersion(chain.StateDB)
	if vget != version {
		t.Fatalf("get current active version error,expect version:%d,get version:%d", version, vget)
	}
}

func TestGovDB_TallyResult(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)

	proposal := getVerProposal(common.Hash{0x03})
	if e := gdb.SetProposal(proposal, chain.StateDB); e != nil {
		t.Errorf("set proposal error,%s", e)
	}

	proposalID := common.Hash{0x03}

	tallyResult := TallyResult{
		ProposalID:    proposalID,
		Yeas:          15,
		Nays:          0,
		Abstentions:   0,
		AccuVerifiers: 1000,
		Status:        Pass,
	}

	if err := gdb.SetTallyResult(tallyResult, chain.StateDB); err != nil {
		t.Fatalf("set vote result error")
	}

	if result, err := gdb.GetTallyResult(proposalID, chain.StateDB); err != nil {
		t.Fatalf("get vote result error,%s", err)
	} else {
		if result.Status != tallyResult.Status {
			t.Fatalf("get vote result error")
		}
	}
}

func TestGovDB_GetTallyResult_ProposalNotFound(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)

	proposalID := common.Hash{0x03}

	tallyResult := TallyResult{
		ProposalID:    proposalID,
		Yeas:          15,
		Nays:          0,
		Abstentions:   0,
		AccuVerifiers: 1000,
		Status:        Pass,
	}

	if err := gdb.SetTallyResult(tallyResult, chain.StateDB); err != nil {
		t.Fatalf("set vote result error")
	}

	if result, err := gdb.GetTallyResult(proposalID, chain.StateDB); err != nil {
		if err == ProposalNotFound {
			t.Log("get expected error")
		} else {
			t.Fatalf("get vote result error,%s", err)
		}
	} else {
		if result.Status != tallyResult.Status {
			t.Fatalf("get vote result error")
		}
	}
}

func TestGovDB_AddActiveNode(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)

	proposal := getTxtProposal()

	for _, node := range voteValueList {
		if err := gdb.AddActiveNode(blockHash, proposal.ProposalID, node.VoteNodeID); err != nil {
			t.Fatalf("add active node error...%s", err)
		}
	}

	if ids, err := gdb.GetActiveNodeList(blockHash, proposal.ProposalID); err != nil {
		t.Fatalf("get active node list error...%s", err)
	} else {
		if len(ids) != len(voteValueList) {
			t.Fatalf(" get active node list error, expect len:%d,get len:%d", len(voteValueList), len(ids))
		}
	}

	if err := gdb.ClearActiveNodes(blockHash, proposal.ProposalID); err != nil {
		t.Fatalf("clear active node list error...%s", err)
	} else {
		if ids, err := gdb.GetActiveNodeList(blockHash, proposal.ProposalID); err != nil {
			t.Fatalf("get active node list after clear error...%s", err)
		} else {
			if len(ids) != 0 {
				t.Fatalf(" get active node list after clear error, expect len:0,get len:%d", len(ids))
			}
		}
	}
}

func TestGovDB_addAccuVerifiers(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	proposalID := generateHash("pipID")

	blockHash := newBlock(chain)

	if err := gdb.addAccuVerifiers(blockHash, proposalID, NodeIDList); err != nil {
		t.Fatalf("addAccuVerifiers error...%s", err)
	} else {
		if nodeList, err := gdb.ListAccuVerifier(blockHash, proposalID); err != nil {
			t.Fatalf("ListAccuVerifier error...%s", err)
		} else {
			if len(nodeList) != 4 {
				t.Fatalf("node count error")
			}
		}
	}
}

func TestGovDB_AddPIPID(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)

	if err := gdb.AddPIPID("pip_1", chain.StateDB); err != nil {
		t.Fatalf("add PIPID error ...%s", err)
	}
}

func TestGovDB_ListPIPID(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)

	if err := gdb.AddPIPID("pip_1", chain.StateDB); err != nil {
		t.Fatalf("add PIPID error ...%s", err)
	}
	if err := gdb.AddPIPID("pip_2", chain.StateDB); err != nil {
		t.Fatalf("add PIPID error ...%s", err)
	}

	if idList, err := gdb.ListPIPID(chain.StateDB); err != nil {
		t.Fatalf("list PIPID error ...%s", err)
	} else {
		if len(idList) != 2 {
			t.Fatalf("list PIPID count error")
		} else {
			t.Log("list PIPID", "idList", idList)
		}
	}
}

func TestGovDB_GetExistProposal(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)

	proposal := getTxtProposal()
	if err := gdb.SetProposal(proposal, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}

	if existing, err := gdb.GetExistProposal(proposal.ProposalID, chain.StateDB); err != nil {
		t.Errorf("get exist proposal error,%s", err)
	} else {
		if existing.GetPIPID() != proposal.GetPIPID() {
			t.Fatalf("get exist proposal error,expect %s,get %s", proposal.GetPIPID(), existing.GetPIPID())
		}
	}

	if _, err := gdb.GetExistProposal(common.Hash{0x10}, chain.StateDB); err != nil {
		if err == ProposalNotFound {
			t.Log("throw exception correctly if not found the proposal")
		} else {
			t.Fatal("do not throw exception correctly if not found the proposal")
		}
	} else {
		t.Fatalf("do not throw exception correctly if not found the proposal")
	}
}

func TestGovDB_FindVotingVersionProposal_success(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	gov := NewGov(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)

	tp := getTxtProposal()
	if err := gdb.SetProposal(tp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, tp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	vp := getVerProposal(common.Hash{0x2})
	if err := gdb.SetProposal(vp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, vp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	cp := getCancelProposal()
	if err := gdb.SetProposal(cp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, cp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}
	if p, err := gov.FindVotingProposal(blockHash, chain.StateDB, Version); err != nil {
		t.Fatalf("find voting proposal ID error,%s", err)
	} else if p == nil {
		t.Log("not find voting proposal ID")
	} else {
		vp := p.(*VersionProposal)
		t.Log("find voting proposal ID success", "proposalID", vp.ProposalID)
	}
}

func TestGovDB_FindVotingVersionProposal_NoVersionProposalID(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	gov := NewGov(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)
	tp := getTxtProposal()
	if err := gdb.SetProposal(tp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, tp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	cp := getCancelProposal()
	if err := gdb.SetProposal(cp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, cp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}
	if p, err := gov.FindVotingProposal(blockHash, chain.StateDB, Version); err != nil {
		t.Fatalf("find voting proposal ID error,%s", err)
	} else if p == nil {
		t.Log("not find voting proposal ID")
	} else {
		vp := p.(*VersionProposal)
		t.Log("find voting proposal ID success", "proposalID", vp.ProposalID)
	}
}

func TestGovDB_FindVotingVersionProposal_DataError(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	gov := NewGov(chain.SnapDB)
	//create block
	blockHash := newBlock(chain)

	tp := getTxtProposal()
	if err := gdb.SetProposal(tp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, tp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	/*vp := getVerProposal(common.Hash{0x2})
	if e := gdb.SetProposal(vp, statedb); e != nil {
		t.Errorf("set proposal error,%s", e)
	}*/
	if err := gdb.AddVotingProposalID(blockHash, common.Hash{0x2}); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}

	cp := getCancelProposal()
	if err := gdb.SetProposal(cp, chain.StateDB); err != nil {
		t.Errorf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, cp.ProposalID); err != nil {
		t.Errorf("add voting proposal ID error,%s", err)
	}
	if p, err := gov.FindVotingProposal(blockHash, chain.StateDB, Version); err != nil {
		if err == ProposalNotFound {
			t.Log("throw a exception correctly if data error")
		} else {
			t.Fatalf("find voting proposal ID error,%s", err)
		}
	} else if p == nil {
		t.Log("not find voting proposal ID")
	} else {
		vp := p.(*VersionProposal)
		t.Log("find voting proposal ID success", "proposalID", vp.ProposalID)
	}
}

func TestGovDB_ListActiveVersion(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	gov := NewGov(chain.SnapDB)

	version := uint32(32)
	//proposal := getVerProposal(common.Hash{0x1})
	if err := gdb.AddActiveVersion(version, 10000, chain.StateDB); err != nil {
		t.Fatalf("add active version error...%s", err)
	}

	version = uint32(33)
	//proposal := getVerProposal(common.Hash{0x1})
	if err := gdb.AddActiveVersion(version, 20000, chain.StateDB); err != nil {
		t.Fatalf("add active version error...%s", err)
	}

	if avList, err := gov.GetCurrentActiveVersionList(chain.StateDB); err != nil {
		t.Fatal("list active version error")
	} else if len(avList) != 2 {
		t.Fatal("count of active version error")
	}
}

func TestGovDB_AddGovernParam(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	blockHash := newBlock(chain)

	if err := gdb.addGovernParam("PPOS", "testName1", "desc1", &ParamValue{"", "initValue", 0}, blockHash); err != nil {
		t.Fatalf("add govern param error...%s", err)
	}
}

func TestGovDB_FindGovernParamValue(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	blockHash := newBlock(chain)
	if err := gdb.addGovernParam("PPOS", "testName1", "desc1", &ParamValue{"", "initValue", 0}, blockHash); err != nil {
		t.Fatalf("add govern param error...%s", err)
	}

	if item, err := gdb.findGovernParamValue("PPOS", "testName1", blockHash); err != nil {
		t.Fatalf("add govern param error...%s", err)
	} else if item == nil {
		t.Logf("govern param not found")
	} else if item.Value == "initValue" {
		assert.Equal(t, "initValue", item.Value, "error")
	}
}

func TestGovDB_setPreactiveVersion(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)
	blockHash := newBlock(chain)
	version := uint32(2561)
	if err := gdb.SetPreActiveVersion(blockHash, version); err != nil {
		t.Fatalf("add govern param error...%s", err)
	}

	ver := gdb.GetPreActiveVersion(blockHash)
	assert.Equal(t, version, ver)
}

func TestGovDB_setPreactiveProposalID(t *testing.T) {
	chain := mock.NewChain()
	defer chain.SnapDB.Clear()

	gdb := NewGovDB(chain.SnapDB)

	pid := common.Hash{0x1}
	//create block
	blockHash := newBlock(chain)
	vp := getVerProposal(pid)
	if err := gdb.SetProposal(vp, chain.StateDB); err != nil {
		t.Fatalf("set proposal error,%s", err)
	}
	if err := gdb.AddVotingProposalID(blockHash, vp.ProposalID); err != nil {
		t.Fatalf("add voting proposal ID error,%s", err)
	}

	if err := gdb.MoveVotingProposalIDToPreActive(blockHash, vp.ProposalID, vp.NewVersion); err != nil {
		t.Fatalf("move voting ID to pre-active ID error,%s", err)
	}
	if p, err := gdb.GetPreActiveProposalID(blockHash); err != nil {
		t.Fatalf("find  pre-active ID error,%s", err)
	} else {
		t.Log("p:", p)
		assert.Equal(t, vp.ProposalID, p)
	}

	version := gdb.GetPreActiveVersion(blockHash)
	assert.Equal(t, vp.NewVersion, version)

	if err := gdb.MovePreActiveProposalIDToEnd(blockHash, vp.ProposalID); err != nil {
		t.Fatalf("reset pre-active ID error,%s", err)
	}

	if p, err := gdb.GetPreActiveProposalID(blockHash); err != nil {
		t.Fatalf("find  pre-active ID error,%s", err)
	} else {
		t.Log("p:", p)
		assert.Equal(t, common.ZeroHash, p)
	}
}

func TestGovDB_Version(t *testing.T) {
	version := uint32(0<<16 | 7<<8 | 4)
	fmt.Println(version)
	log.Warn("Store version for gov into genesis statedb", "genesis version", fmt.Sprintf("%d/%s", version, params.FormatVersion(version)))

	var hash common.Hash
	assert.Equal(t, common.ZeroHash, hash)
	assert.Equal(t, common.Hash{}, hash)
	assert.Equal(t, common.Hash{0x00}, hash)
}

func newBlock(chain *mock.Chain) common.Hash {
	chain.AddBlock()
	chain.SnapDB.NewBlock(chain.CurrentHeader().Number, chain.CurrentHeader().ParentHash, chain.CurrentHeader().Hash())
	return chain.CurrentHeader().Hash()
}

func getTxtProposal() *TextProposal {
	return &TextProposal{
		ProposalID:   common.Hash{0x01},
		ProposalType: Text,
		PIPID:        "em1",
		SubmitBlock:  uint64(1000),
		Proposer:     enode.IDv0{},
	}
}

func getVerProposal(proposalId common.Hash) *VersionProposal {
	return &VersionProposal{
		ProposalID:      proposalId,
		ProposalType:    Version,
		PIPID:           "em2",
		SubmitBlock:     uint64(1000),
		EndVotingRounds: uint64(8),
		Proposer:        enode.IDv0{},
		NewVersion:      32,
	}
}

func getCancelProposal() *CancelProposal {
	return &CancelProposal{
		ProposalID:      common.Hash{0x03},
		ProposalType:    Cancel,
		PIPID:           "em3",
		SubmitBlock:     uint64(1000),
		EndVotingRounds: uint64(5),
		Proposer:        enode.IDv0{},
		TobeCanceled:    common.Hash{0x02},
	}
}

func getParamProposal() *ParamProposal {
	return &ParamProposal{
		ProposalID:   common.Hash{0x05},
		ProposalType: Param,
		PIPID:        "em5",
		SubmitBlock:  uint64(1000),
		Proposer:     enode.IDv0{},
		Module:       "PPOS",
		Name:         "testName1",
		NewValue:     "newValue1",
	}
}

var voteValueList = []VoteValue{
	{
		VoteNodeID: enode.MustHexIDv0("0x1dd9d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: Yes,
	},
	{
		VoteNodeID: enode.MustHexIDv0("0x1dd8d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: Yes,
	},
	{
		VoteNodeID: enode.MustHexIDv0("0x1dd7d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: Yes,
	},
	{
		VoteNodeID: enode.MustHexIDv0("0x1dd6d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: Yes,
	},
	{
		VoteNodeID: enode.MustHexIDv0("0x1dd5d65c4552b5eb43d5ad55a2ee3f56c6cbc1c64a5c8d659f51fcd51bace24351232b8d7821617d2b29b54b81cdefb9b3e9c37d7fd5f63270bcc9e1a6f6a439"),
		VoteOption: Yes,
	},
}

var NodeIDList = []enode.IDv0{
	enode.MustHexIDv0("5a942bc607d970259e203f5110887d6105cc787f7433c16ce28390fb39f1e67897b0fb445710cc836b89ed7f951c57a1f26a0940ca308d630448b5bd391a8aa6"),
	enode.MustHexIDv0("c453d29394e613e85999129b8fb93146d584d5a0be16f7d13fd1f44de2d01bae104878eba8e8f6b8d2c162b5a35d5939d38851f856e56186471dd7de57e9bfa9"),
	enode.MustHexIDv0("2c1733caf5c23086612a309f5ee8e76ca45455351f7cf069bcde59c07175607325cf2bf2485daa0fbf1f9cdee6eea246e5e00b9a0d0bfed0f02b37f3b0c70490"),
	enode.MustHexIDv0("e7edfb4f9c3e1fe0288ddcf0894535214fa03acea941c7360ccf90e86460aefa118ba9f2573921349c392cd1b5d4db90b4795ab353df3c915b2e8481d241ec57"),
}

func generateHash(n string) common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(n))
	return rlpHash(buf.Bytes())
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewLegacyKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
