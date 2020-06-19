// Copyright 2018-2020 The PlatON Network Authors
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

package plugin

import (
	"encoding/hex"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/node"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	//	"github.com/PlatONnetwork/PlatON-Go/core/state"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	//	"github.com/PlatONnetwork/PlatON-Go/core/state"
	//	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var (
	snapdb    snapshotdb.DB
	govPlugin *GovPlugin
	//	evm       *vm.EVM
	//newVersion     = uint32(2<<16 | 0<<8 | 0)
	/*endVotingBlock uint64
	activeBlock    uint64*/
	stateDB xcom.StateDB
	chainID = big.NewInt(100)

//	stk            *StakingPlugin
)

func setup(t *testing.T) func() {
	t.Log("setup()......")

	state, genesis, _ := newChainState()
	newEvm(blockNumber, blockHash, state)
	stateDB = state
	newPlugins()

	GovPluginInstance().SetChainID(chainID)
	govPlugin = GovPluginInstance()
	stk = StakingInstance()

	lastBlockHash = genesis.Hash()

	build_staking_data(genesis.Hash())

	snapdb = snapshotdb.Instance()

	// init data
	if _, err := gov.InitGenesisGovernParam(common.ZeroHash, snapdb, 2048); err != nil {
		t.Fatalf("cannot init genesis govern param...")
	}

	if freezeDuration, err := gov.GovernUnStakeFreezeDuration(lastBlockNumber, lastBlockHash); err != nil {
		t.Fatalf("cannot find init gov param (FreezeDuration)")
	} else {
		t.Logf("freezeDuration:: %d", freezeDuration)
	}

	if maxEvidenceAge, err := gov.GovernMaxEvidenceAge(lastBlockNumber, lastBlockHash); err != nil {
		t.Fatalf("cannot find init gov param(EvidenceAge)")
	} else {
		t.Logf("maxEvidenceAge:: %d", maxEvidenceAge)
	}

	return func() {
		t.Log("tear down()......")
		snapdb.Clear()
	}
}

func submitText(t *testing.T, pid common.Hash) {
	vp := &gov.TextProposal{
		ProposalID:   pid,
		ProposalType: gov.Text,
		PIPID:        "textPIPID",
		SubmitBlock:  1,
		Proposer:     nodeIdArr[0],
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], lastBlockHash, 0)
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))()
	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	if err != nil {
		t.Fatalf("submit text proposal err: %s", err)
	}
}

func buildTextProposal(proposalID common.Hash, pipID string) *gov.TextProposal {
	return &gov.TextProposal{
		ProposalID:   proposalID,
		ProposalType: gov.Text,
		PIPID:        pipID,
		SubmitBlock:  1,
		Proposer:     nodeIdArr[0],
	}
}

func submitVersion(t *testing.T, pid common.Hash) {
	vp := &gov.VersionProposal{
		ProposalID:      pid,
		ProposalType:    gov.Version,
		PIPID:           "versionIPID",
		SubmitBlock:     1,
		EndVotingRounds: xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()),
		Proposer:        nodeIdArr[0],
		NewVersion:      promoteVersion,
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], blockHash, 0)

	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		t.Fatalf("submit version proposal err: %s", err)
	}
}

func buildVersionProposal(proposalID common.Hash, pipID string, endVotingRounds uint64, newVersion uint32) *gov.VersionProposal {
	return &gov.VersionProposal{
		ProposalID:      proposalID,
		ProposalType:    gov.Version,
		PIPID:           pipID,
		SubmitBlock:     1,
		EndVotingRounds: endVotingRounds,
		Proposer:        nodeIdArr[0],
		NewVersion:      newVersion,
	}
}

func submitCancel(t *testing.T, pid, tobeCanceled common.Hash) {
	pp := &gov.CancelProposal{
		ProposalID:      pid,
		ProposalType:    gov.Cancel,
		PIPID:           "CancelPIPID",
		SubmitBlock:     1,
		EndVotingRounds: xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()) - 1,
		Proposer:        nodeIdArr[0],
		TobeCanceled:    tobeCanceled,
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(sender, pp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		t.Fatalf("submit cancel proposal err: %s", err)
	}
}

func allVote(t *testing.T, pid common.Hash) {
	//for _, nodeID := range nodeIdArr {
	currentValidatorList, _ := stk.ListCurrentValidatorID(lastBlockHash, lastBlockNumber)
	voteCount := len(currentValidatorList)
	chandler := node.GetCryptoHandler()

	for i := 0; i < voteCount; i++ {
		vote := gov.VoteInfo{
			ProposalID: pid,
			VoteNodeID: nodeIdArr[i],
			VoteOption: gov.Yes,
		}

		chandler.SetPrivateKey(priKeyArr[i])
		versionSign := common.VersionSign{}
		versionSign.SetBytes(chandler.MustSign(promoteVersion))

		err := gov.Vote(sender, vote, lastBlockHash, 1, promoteVersion, versionSign, stk, stateDB)
		if err != nil {
			t.Fatalf("vote err: %s.", err)
		}
	}
}

func halfVote(t *testing.T, pid common.Hash) {
	currentValidatorList, _ := stk.ListCurrentValidatorID(lastBlockHash, lastBlockNumber)
	voteCount := len(currentValidatorList)
	chandler := node.GetCryptoHandler()
	for i := 0; i < voteCount/2; i++ {
		vote := gov.VoteInfo{
			ProposalID: pid,
			VoteNodeID: nodeIdArr[i],
			VoteOption: gov.Yes,
		}

		chandler.SetPrivateKey(priKeyArr[i])
		versionSign := common.VersionSign{}
		versionSign.SetBytes(chandler.MustSign(promoteVersion))

		err := gov.Vote(sender, vote, lastBlockHash, 1, promoteVersion, versionSign, stk, stateDB)
		if err != nil {
			t.Fatalf("vote err: %s.", err)
		}
	}

}

func beginBlock(t *testing.T) {
	err := govPlugin.BeginBlock(lastBlockHash, &lastHeader, stateDB)
	if err != nil {
		t.Fatalf("begin block err... %s", err)
	}
}

func endBlock(t *testing.T) {
	err := govPlugin.EndBlock(lastBlockHash, &lastHeader, stateDB)
	if err != nil {
		t.Fatalf("end block err... %s", err)
	}
}

func TestGovPlugin_SubmitText(t *testing.T) {
	defer setup(t)()
	submitText(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted text proposal error:", err)
	} else {
		t.Log("Get the submitted text proposal success:", p)
	}
}

func TestGovPlugin_GetProposal(t *testing.T) {
	defer setup(t)()
	submitText(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	buildBlockNoCommit(2)

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatalf("Get proposal error: %s", err)
	} else {
		t.Logf("Get proposal success: %x", p.GetProposalID())
	}
}

func TestGovPlugin_SubmitText_PIPID_empty(t *testing.T) {
	defer setup(t)()

	tp := buildTextProposal(txHashArr[0], "")
	err := gov.Submit(sender, tp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		if err == gov.PIPIDEmpty {
			t.Logf("detected empty PIPID.")
		} else {
			t.Fatal("didn't detect empty PIPID.")
		}
	}
}

func TestGovPlugin_SubmitText_PIPID_duplicated(t *testing.T) {
	defer setup(t)()

	t.Log("CurrentActiveVersion", "version", gov.GetCurrentActiveVersion(stateDB))

	tp := buildTextProposal(txHashArr[0], "pipID")

	err := gov.Submit(sender, tp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		t.Fatalf("submit proposal err: %s", err)
	}

	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	buildBlockNoCommit(2)

	if p, err := gov.ListPIPID(stateDB); err == nil {
		t.Log("ListPIPID", "p", p)
	}

	tp2 := buildTextProposal(txHashArr[1], "pipID")

	err = gov.Submit(sender, tp2, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		if err == gov.PIPIDExist {
			t.Logf("detected duplicated PIPID.")
		} else {
			t.Fatal("didn't detect duplicated PIPID.")
		}
	}
}

func TestGovPlugin_SubmitText_invalidSender(t *testing.T) {
	defer setup(t)()

	vp := &gov.TextProposal{
		ProposalID:   txHashArr[0],
		ProposalType: gov.Text,
		PIPID:        "textPIPID",
		SubmitBlock:  1,
		Proposer:     nodeIdArr[0],
	}

	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(anotherSender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID) //sender error
	if err != nil {
		if err == gov.TxSenderDifferFromStaking || err == gov.TxSenderIsNotVerifier {
			t.Log("detected invalid sender.", err)
		} else {
			t.Fatal("didn't detect invalid sender.")
		}
	}
}

func TestGovPlugin_SubmitText_invalidType(t *testing.T) {
	defer setup(t)()

	vp := &gov.TextProposal{
		ProposalID:   txHashArr[0],
		ProposalType: gov.Version, //error type
		PIPID:        "textPIPID",
		SubmitBlock:  1,
		Proposer:     nodeIdArr[0],
	}

	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(anotherSender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID) //sender error
	if err != nil {
		if err == gov.ProposalTypeError {
			t.Log("detected invalid type.", err)
		} else {
			t.Fatal("didn't detect invalid type.")
		}
	}
}

func TestGovPlugin_SubmitText_Proposer_empty(t *testing.T) {
	defer setup(t)()

	vp := &gov.TextProposal{
		ProposalID:   txHashArr[0],
		ProposalType: gov.Text,
		PIPID:        "textPIPID",
		SubmitBlock:  1,
		Proposer:     discover.ZeroNodeID,
	}

	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID) //empty proposal
	if err != nil {
		if err == gov.ProposerEmpty {
			t.Log("detected invalid proposer.", err)
		} else {
			t.Fatal("didn't detect invalid proposer.")
		}
	}
}

func TestGovPlugin_SubmitVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}
}

func TestGovPlugin_SubmitVersion_PIPID_empty(t *testing.T) {
	defer setup(t)()

	vp := buildVersionProposal(txHashArr[0], "", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()), uint32(1<<16|2<<8|0))
	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		if err == gov.PIPIDEmpty {
			t.Logf("detected empty PIPID.")
		} else {
			t.Fatal("didn't detect empty PIPID.")
		}
	}
}

func TestGovPlugin_SubmitVersion_PIPID_duplicated(t *testing.T) {
	defer setup(t)()

	t.Log("CurrentActiveVersion", "version", gov.GetCurrentActiveVersion(stateDB))

	vp := buildVersionProposal(txHashArr[0], "pipID", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()), uint32(1<<16|2<<8|0))

	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		t.Fatalf("submit proposal err: %s", err)
	}

	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	buildBlockNoCommit(2)

	if p, err := gov.ListPIPID(stateDB); err == nil {
		t.Log("ListPIPID", "p", p)
	}

	vp2 := buildVersionProposal(txHashArr[1], "pipID", xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()), uint32(1<<16|3<<8|0))

	err = gov.Submit(sender, vp2, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		if err == gov.PIPIDExist {
			t.Logf("detected duplicated PIPID.")
		} else {
			t.Fatal("didn't detect duplicated PIPID.")
		}
	}
}

func TestGovPlugin_SubmitVersion_invalidEndVotingRounds(t *testing.T) {
	defer setup(t)()

	vp := &gov.VersionProposal{
		ProposalID:      txHashArr[0],
		ProposalType:    gov.Version,
		PIPID:           "versionPIPID",
		SubmitBlock:     1,
		EndVotingRounds: xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()) + 1, //error
		Proposer:        nodeIdArr[0],
		NewVersion:      promoteVersion,
	}
	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		if err == gov.EndVotingRoundsTooLarge {
			t.Logf("detected invalid end-voting-rounds.")
		} else {
			t.Fatal("didn't detect invalid end-voting-rounds.")
		}
	}
}

func TestGovPlugin_SubmitVersion_ZeroEndVotingRounds(t *testing.T) {
	defer setup(t)()

	vp := &gov.VersionProposal{
		ProposalID:      txHashArr[0],
		ProposalType:    gov.Version,
		PIPID:           "versionPIPID",
		SubmitBlock:     1,
		EndVotingRounds: 0, //error
		Proposer:        nodeIdArr[0],
		NewVersion:      promoteVersion,
	}
	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		if err == gov.EndVotingRoundsTooSmall {
			t.Logf("detected zero end-voting-rounds.")
		} else {
			t.Fatal("didn't detect zero end-voting-rounds.")
		}
	}
}

func TestGovPlugin_SubmitVersion_NewVersionError(t *testing.T) {
	defer setup(t)()

	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	version := uint32(1<<16 | 2<<8 | 0)
	newVersionErr := uint32(1<<16 | 2<<8 | 4)

	if err := gov.AddActiveVersion(version, 10000, state); err != nil {
		t.Fatalf("add active version error...%s", err)
	}

	t.Log("CurrentActiveVersion", "version", gov.GetCurrentActiveVersion(state))

	vp := &gov.VersionProposal{
		ProposalID:      txHashArr[0],
		ProposalType:    gov.Version,
		PIPID:           "versionPIPID",
		SubmitBlock:     1,
		EndVotingRounds: xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()),
		Proposer:        nodeIdArr[0],
		NewVersion:      newVersionErr, //error, less than activeVersion
	}

	err := gov.Submit(sender, vp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		if err == gov.NewVersionError {
			t.Logf("detected invalid NewVersioin.")
		} else {
			t.Fatal("didn't detect invalid NewVersion.", "err", err)
		}
	}
}

func TestGovPlugin_SubmitCancel(t *testing.T) {
	defer setup(t)()

	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}

	submitCancel(t, txHashArr[1], txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err = gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted cancel proposal error:", err)
	} else {
		t.Log("Get the submitted cancel proposal success:", p)
	}
}

func TestGovPlugin_SubmitCancel_invalidEndVotingRounds(t *testing.T) {
	defer setup(t)()

	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("Get the submitted version proposal error:", err)
	} else {
		t.Log("Get the submitted version proposal success:", p)
	}

	pp := &gov.CancelProposal{
		ProposalID:      txHashArr[1],
		ProposalType:    gov.Cancel,
		PIPID:           "CancelPIPID",
		SubmitBlock:     1,
		EndVotingRounds: xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()),
		Proposer:        nodeIdArr[1],
		TobeCanceled:    txHashArr[0],
	}

	//state := stateDB.(*state.StateDB)
	//state.Prepare(txHashArr[0], lastBlockHash, 0)

	err = gov.Submit(sender, pp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		if err == gov.EndVotingRoundsTooLarge {
			t.Logf("detected invalid end-voting-rounds.")
		} else {
			t.Fatal("didn't detect invalid end-voting-rounds.")
		}
	}
}

func TestGovPlugin_SubmitCancel_noVersionProposal(t *testing.T) {
	defer setup(t)()

	pp := &gov.CancelProposal{
		ProposalID:      txHashArr[1],
		ProposalType:    gov.Cancel,
		PIPID:           "cancelPIPID",
		SubmitBlock:     1,
		EndVotingRounds: xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()) - 1,
		Proposer:        nodeIdArr[0],
		TobeCanceled:    txHashArr[0],
	}
	state := stateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], lastBlockHash, 0)

	err := gov.Submit(sender, pp, lastBlockHash, lastBlockNumber, stk, stateDB, chainID)
	if err != nil {
		if err == gov.TobeCanceledProposalNotFound {
			t.Logf("detected this case.")
		} else {
			t.Fatal("didn't detect is case.")
		}
	}
}

func TestGovPlugin_VoteSuccess(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	nodeIdx := 3
	v := gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	nodeIdx = 1
	v = gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler = node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign = common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err = gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	votedValue, err := gov.ListVoteValue(txHashArr[0], lastBlockHash)
	if err != nil {
		t.Fatal("vote err:", err)
	} else {
		t.Log("voted count:", len(votedValue))
	}

	votedMap, err := gov.GetVotedVerifierMap(txHashArr[0], lastBlockHash)
	if err != nil {
		t.Fatal("vote failed, cannot list voted verifiers", err)
	} else {
		t.Log("voted count:", len(votedMap))
	}
}

func TestGovPlugin_Vote_Repeat(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	nodeIdx := 3
	v := gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	v = gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx], //repeated
		gov.Yes,
	}

	err = gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		if err == gov.VoteDuplicated {
			t.Log("detected repeated vote", err)
		} else {
			t.Fatal("didn't detect repeated vote")
		}
	}
}

func TestGovPlugin_Vote_invalidSender(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	nodeIdx := 3
	v := gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.Vote(anotherSender, v, lastBlockHash, 2, initProgramVersion, versionSign, stk, stateDB)
	if err != nil {
		if err == gov.TxSenderIsNotVerifier || err == gov.TxSenderDifferFromStaking {
			t.Log("detected invalid sender", err)
		} else {
			t.Fatal("didn't detect invalid sender")
		}
	}
}

func TestGovPlugin_DeclareVersion_rightVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)
	nodeIdx := 0
	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.DeclareVersion(sender, nodeIdArr[nodeIdx], promoteVersion, versionSign, lastBlockHash, 2, stk, stateDB)
	if err != nil {
		t.Fatalf("Declare Version err ...%s", err)
	}

	activeNodeList, err := gov.GetActiveNodeList(lastBlockHash, txHashArr[0])
	if err != nil {
		t.Fatalf("List actived nodes error: %s", err)
	} else {
		t.Logf("List actived nodes success: %d", len(activeNodeList))
	}
}

func TestGovPlugin_DeclareVersion_wrongSign(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	wrongVersion := uint32(1<<16 | 1<<8 | 1)

	nodeIdx := 0
	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(wrongVersion))

	err := gov.DeclareVersion(sender, nodeIdArr[nodeIdx], promoteVersion, versionSign, lastBlockHash, 2, stk, stateDB)

	if err != nil {
		if err == gov.VersionSignError {
			t.Log("detected incorrect version declaration.", err)
		} else {
			t.Fatal("not detected incorrect version declaration.", err)
		}
	}
}

func TestGovPlugin_DeclareVersion_wrongVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	wrongVersion := uint32(1<<16 | 1<<8 | 1)

	nodeIdx := 0
	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(wrongVersion))

	err := gov.DeclareVersion(sender, nodeIdArr[nodeIdx], wrongVersion, versionSign, lastBlockHash, 2, stk, stateDB)

	if err != nil {
		if err == gov.DeclareVersionError {
			t.Log("detected incorrect version declaration.", err)
		} else {
			t.Fatal("not detected incorrect version declaration.", err)
		}
	}
}

func TestGovPlugin_VotedNew_DeclareOld(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	nodeIdx := 3
	v := gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	nodeIdx = 1
	v = gov.VoteInfo{
		txHashArr[0],
		nodeIdArr[nodeIdx],
		gov.Yes,
	}

	chandler = node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign = common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err = gov.Vote(sender, v, lastBlockHash, 2, promoteVersion, versionSign, stk, stateDB)
	if err != nil {
		t.Fatal("vote err:", err)
	}

	votedValue, err := gov.ListVoteValue(txHashArr[0], lastBlockHash)
	if err != nil {
		t.Fatal("vote err:", err)
	} else {
		t.Log("voted count:", len(votedValue))
	}

	//declare
	versionSign = common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(initProgramVersion))

	err = gov.DeclareVersion(sender, nodeIdArr[nodeIdx], initProgramVersion, versionSign, lastBlockHash, 2, stk, stateDB)

	if err != nil {
		if err == gov.DeclareVersionError {
			t.Log("detected incorrect version declaration.", err)
		} else {
			t.Fatal("not detected incorrect version declaration.", err)
		}
	}
}

func TestGovPlugin_DeclareVersion_invalidSender(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	nodeIdx := 0
	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	versionSign := common.VersionSign{}
	versionSign.SetBytes(chandler.MustSign(promoteVersion))

	err := gov.DeclareVersion(anotherSender, nodeIdArr[nodeIdx], promoteVersion, versionSign, lastBlockHash, 2, stk, stateDB)
	if err != nil {
		if err == gov.TxSenderDifferFromStaking || err == gov.TxSenderIsNotCandidate {
			t.Log("detected an incorrect version declaration.", err)
		} else {
			t.Fatal("didn't detected an incorrect version declaration.", err)
		}
	}
}

func TestGovPlugin_ListProposal(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	pList, err := gov.ListProposal(lastBlockHash, stateDB)
	if err != nil {
		t.Fatalf("List all proposals error: %s", err)
	} else {
		t.Logf("List all proposals success: %d", len(pList))
	}

}

func TestGovPlugin_textProposalPassed(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	buildBlockNoCommit(2)

	allVote(t, txHashArr[0])
	sndb.Commit(lastBlockHash) //commit
	sndb.Compaction()          //write to level db

	p, err := gov.GetProposal(txHashArr[0], stateDB)
	if err != nil {
		t.Fatal("find proposal error", "err", err)
	}

	lastBlockNumber = uint64(xutil.CalcBlocksEachEpoch() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(uint64(xutil.CalcBlocksEachEpoch()))
	beginBlock(t)
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(p.GetEndVotingBlock() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(p.GetEndVotingBlock())

	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	endBlock(t)

	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(3), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	sndb.Commit(lastBlockHash)

	result, err := gov.GetTallyResult(txHashArr[0], stateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.Pass {
		t.Log("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	} else {
		t.Fatal("tallyResult", "status", result.Status, "yeas", result.Yeas, "accuVerifiers", result.AccuVerifiers)
	}
}

func TestGovPlugin_textProposalFailed(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	endVotingBlock := xutil.CalEndVotingBlock(1, xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()))
	//	actvieBlock := xutil.CalActiveBlock(endVotingBlock)

	buildBlockNoCommit(2)

	halfVote(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(xutil.CalcBlocksEachEpoch() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(uint64(xutil.CalcBlocksEachEpoch()))
	beginBlock(t)
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(endVotingBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(endVotingBlock)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := gov.GetTallyResult(txHashArr[0], stateDB)
	if err != nil {
		t.Errorf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.Failed {
		t.Logf("the result status, %s", result.Status.ToString())
	} else {
		t.Fatalf("the result status error, %s", result.Status.ToString())
	}
}

func TestGovPlugin_versionProposalPreActive(t *testing.T) {

	defer setup(t)()

	submitText(t, txHashArr[0])
	submitVersion(t, txHashArr[1])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	endVotingBlock := xutil.CalEndVotingBlock(1, xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()))
	//	actvieBlock := xutil.CalActiveBlock(endVotingBlock)

	buildBlockNoCommit(2)

	allVote(t, txHashArr[0])
	allVote(t, txHashArr[1])
	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(xutil.CalcBlocksEachEpoch() - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(uint64(xutil.CalcBlocksEachEpoch()))

	beginBlock(t)

	sndb.Commit(lastBlockHash)

	//buildSnapDBDataCommitted(20001, 22229)
	sndb.Compaction()
	lastBlockNumber = uint64(endVotingBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(endVotingBlock)
	endBlock(t)
	sndb.Commit(lastBlockHash)

	result, err := gov.GetTallyResult(txHashArr[0], stateDB)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else {
		t.Logf("the result status, %s", result.Status.ToString())
	}

	result, err = gov.GetTallyResult(txHashArr[1], stateDB)
	if err != nil {
		t.Fatalf("%s", err)
	}
	if result == nil {
		t.Fatal("cannot find the tally result")
	} else if result.Status == gov.PreActive {
		t.Logf("the result status, %s", result.Status.ToString())
	} else {
		t.Logf("the result status error, %s", result.Status.ToString())
	}
}

func TestGovPlugin_GetPreActiveVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	buildBlockNoCommit(2)

	if err := gov.SetPreActiveVersion(lastBlockHash, uint32(10)); err != nil {
		t.Error("SetPreActiveVersion error", err)
	} else {
		ver := gov.GetPreActiveVersion(lastBlockHash)
		assert.Equal(t, uint32(10), ver)
	}

}

func TestGovPlugin_GetActiveVersion(t *testing.T) {
	defer setup(t)()
	submitVersion(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	buildBlockNoCommit(2)

	ver := gov.GetCurrentActiveVersion(stateDB)
	assert.Equal(t, initProgramVersion, ver)
}

func TestGovPlugin_versionProposalActive(t *testing.T) {

	defer setup(t)()

	//submit version proposal
	submitVersion(t, txHashArr[0])
	sndb.Commit(lastBlockHash)
	sndb.Compaction() //flush to LevelDB

	endVotingBlock := xutil.CalEndVotingBlock(1, xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()))
	actvieBlock := xutil.CalActiveBlock(endVotingBlock)

	buildBlockNoCommit(2)
	//voting
	allVote(t, txHashArr[0])

	sndb.Commit(lastBlockHash)
	sndb.Compaction()

	lastBlockNumber = uint64(endVotingBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	build_staking_data_more(endVotingBlock)

	//tally result
	endBlock(t)
	sndb.Commit(lastBlockHash)
	sndb.Compaction()
	lastBlockNumber = uint64(actvieBlock - 1)
	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	lastBlockHash = lastHeader.Hash()
	sndb.SetCurrent(lastBlockHash, *big.NewInt(int64(lastBlockNumber)), *big.NewInt(int64(lastBlockNumber)))

	//buildBlockNoCommit(23480)
	build_staking_data_more(actvieBlock)
	//active
	beginBlock(t)
	sndb.Commit(lastBlockHash)

	activeVersion := gov.GetCurrentActiveVersion(stateDB)
	if activeVersion == promoteVersion {
		t.Logf("active SUCCESS, %d", activeVersion)
	} else {
		t.Fatalf("active FALSE, %d", activeVersion)
	}
}

func TestGovPlugin_printVersion(t *testing.T) {
	defer setup(t)()

	t.Logf("ver.1.2.0, %d", uint32(1<<16|2<<8|0))

}

func TestGovPlugin_TestNodeID(t *testing.T) {
	var nodeID discover.NodeID
	nodeID = [64]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	t.Logf("nodeID is empty, %t", nodeID == discover.ZeroNodeID)

}

/*func TestNodeID1(t *testing.T) {
	var nodeID discover.NodeID
	nodeID = [64]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	t.Error("nodeID is empty", "nodeID", nodeID)

	var proposalID common.Hash
	proposalID = [32]byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x01}

	t.Error("proposalID is empty", "proposalID", proposalID)
}*/

func TestGovPlugin_Test_MakeExtraData(t *testing.T) {
	defer setup(t)()

	lastHeader = types.Header{
		Number: big.NewInt(int64(lastBlockNumber)),
	}
	t.Log(lastHeader.Extra)
	beginBlock(t)
	t.Log(lastHeader.Extra)

	if len(lastHeader.Extra) > 0 {
		var tobeDecoded []byte
		tobeDecoded = lastHeader.Extra
		if len(lastHeader.Extra) <= 32 {
			tobeDecoded = lastHeader.Extra
		} else {
			tobeDecoded = lastHeader.Extra[:32]
		}

		var extraData []interface{}
		err := rlp.DecodeBytes(tobeDecoded, &extraData)
		if err != nil {
			t.Error("rlp decode header extra error")
		}
		//reference to makeExtraData() in gov_plugin.go
		if len(extraData) == 4 {
			versionBytes := extraData[0].([]byte)
			versionInHeader := common.BytesToUint32(versionBytes)

			activeVersion := gov.GetCurrentActiveVersion(stateDB)
			t.Log("verify header version", "headerVersion", versionInHeader, "activeVersion", activeVersion, "blockNumber", lastHeader.Number.Uint64())
			assert.Equal(t, activeVersion, versionInHeader)
		} else {
			t.Fatalf("unknown header extra data, elementCount= %d", len(extraData))
		}
	}

}

func TestGovPlugin_Test_version(t *testing.T) {
	ver := uint32(66048) //1.2.0
	t.Log(common.Uint32ToBytes(ver))
}

func TestGovPlugin_Test_genVersionSign(t *testing.T) {

	ver := uint32(66048) //1.2.0
	chandler := node.GetCryptoHandler()

	for i := 0; i < 4; i++ {
		chandler.SetPrivateKey(priKeyArr[i])
		t.Log("0x" + hex.EncodeToString(chandler.MustSign(ver)))
	}

}

var (
	chandler *node.CryptoHandler
	priKey   = crypto.HexMustToECDSA("8e1477549bea04b97ea15911e2e9b3041b7a9921f80bd6ddbe4c2b080473de22")
	nodeID   = discover.MustHexID("3e7864716b671c4de0dc2d7fd86215e0dcb8419e66430a770294eb2f37b714a07b6a3493055bb2d733dee9bfcc995e1c8e7885f338a69bf6c28930f3cf341819")
)

func initChandlerHandler() {
	chandler = node.GetCryptoHandler()
	chandler.SetPrivateKey(priKey)
}
