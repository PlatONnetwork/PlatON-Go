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

package gov

import (
	"fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/node"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/mock"
)

var (
	sender = common.MustBech32ToAddress("lax1pmhjxvfqeccm87kzpkkr08djgvpp55355nr8j7")
	nodeID = discover.MustHexID("0x362003c50ed3a523cdede37a001803b8f0fed27cb402b3d6127a1a96661ec202318f68f4c76d9b0bfbabfd551a178d4335eaeaa9b7981a4df30dfc8c0bfe3384")

	priKey = crypto.HexMustToECDSA("0c6ccec28e36dc5581ea3d8af1303c774b51523da397f55cdc4acd9d2b988132")

	senderBalance = "9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999"

	tpProposalID = common.HexToHash("0x00000000000000000000000000000000000000886d5ba2d3dfb2e2f6a1814f22")
	tpPIPID      = "tpPIPID"

	vpProposalID      = common.HexToHash("0x000000000000000000000000000000005249b59609286f2fa91a2abc8555e887")
	vpPIPID           = "vpPIPID"
	vpEndVotingRounds = uint64(2)

	tempActiveVersion = params.GenesisVersion + uint32(0<<16|1<<8|0)

	chainID = big.NewInt(100)
)

type MockStaking struct {
	DeclaeredVodes map[discover.NodeID]uint32
}

func (stk *MockStaking) GetVerifierList(blockHash common.Hash, blockNumber uint64, isCommit bool) (staking.ValidatorExQueue, error) {
	valEx := &staking.ValidatorEx{
		NodeId:          nodeID,
		StakingAddress:  sender,
		ProgramVersion:  params.GenesisVersion,
		StakingBlockNum: 0,
	}

	return []*staking.ValidatorEx{valEx}, nil
}

func (stk *MockStaking) ListVerifierNodeID(blockHash common.Hash, blockNumber uint64) ([]discover.NodeID, error) {
	return []discover.NodeID{nodeID}, nil
}

func (stk *MockStaking) GetCanBaseList(blockHash common.Hash, blockNumber uint64) (staking.CandidateBaseQueue, error) {
	candidate := &staking.CandidateBase{
		NodeId:          nodeID,
		StakingAddress:  sender,
		ProgramVersion:  params.GenesisVersion,
		StakingBlockNum: 0,
	}
	return []*staking.CandidateBase{candidate}, nil
}

func (stk *MockStaking) GetCandidateInfo(blockHash common.Hash, addr common.NodeAddress) (*staking.Candidate, error) {
	return nil, nil
}

func (stk *MockStaking) GetCanBase(blockHash common.Hash, addr common.NodeAddress) (*staking.CandidateBase, error) {
	return nil, nil
}

func (stk *MockStaking) GetCanMutable(blockHash common.Hash, addr common.NodeAddress) (*staking.CandidateMutable, error) {
	can := &staking.CandidateMutable{Status: staking.Valided}
	return can, nil
}
func (stk *MockStaking) DeclarePromoteNotify(blockHash common.Hash, blockNumber uint64, nodeId discover.NodeID, programVersion uint32) error {
	if stk.DeclaeredVodes == nil {
		stk.DeclaeredVodes = make(map[discover.NodeID]uint32)
	}
	stk.DeclaeredVodes[nodeID] = programVersion
	return nil
}

func (stk *MockStaking) ListDeclaredNode() map[discover.NodeID]uint32 {
	return stk.DeclaeredVodes
}

func NewMockStaking() *MockStaking {
	return &MockStaking{}
}

func commit_sndb(chain *mock.Chain) {
	if err := chain.SnapDB.Commit(chain.CurrentHeader().Hash()); err != nil {
		fmt.Println("commit_sndb error, blockNumber:", chain.CurrentHeader().Number.Uint64(), err)
	}
}

func prepair_sndb(chain *mock.Chain) {
	/*if txHash == common.ZeroHash {
		chain.AddBlock()
	} else {
		chain.AddBlockWithTxHash(txHash)
	}*/

	chain.AddBlock()
	//fmt.Println("prepair_sndb::::::", chain.CurrentHeader().ParentHash.Hex())
	if err := chain.SnapDB.NewBlock(chain.CurrentHeader().Number, chain.CurrentHeader().ParentHash, chain.CurrentHeader().Hash()); err != nil {
		fmt.Println("prepair_sndb error:", err)
	}
}

func skip_emptyBlock(chain *mock.Chain, blockNumber uint64) {
	if blockNumber > chain.CurrentHeader().Number.Uint64() {
		cnt := blockNumber - chain.CurrentHeader().Number.Uint64()
		for i := uint64(0); i < cnt; i++ {
			prepair_sndb(chain)
			commit_sndb(chain)
		}
	} else {
		fmt.Println("Warning: blockNumber < currentBlockNumber")
	}
}

func getHandler() *node.CryptoHandler {
	chandler := node.GetCryptoHandler()
	chandler.SetPrivateKey(priKey)
	return chandler
}

func sign(version uint32) []byte {
	return getHandler().MustSign(version)
}

func newChain() *mock.Chain {
	chain := mock.NewChain()
	sBalance, _ := new(big.Int).SetString(senderBalance, 10)
	chain.StateDB.AddBalance(sender, sBalance)
	return chain
}

func clear(chain *mock.Chain, t *testing.T) {
	t.Log("tear down()......")
	if err := chain.SnapDB.Clear(); err != nil {
		t.Error("clear chain.SnapDB error", err)
	}

}

func setup(t *testing.T) *mock.Chain {
	t.Log("setup()......")
	chain := newChain()

	chain.AddBlock()
	err := chain.SnapDB.NewBlock(chain.CurrentHeader().Number, chain.CurrentHeader().ParentHash, chain.CurrentHeader().Hash())
	if err != nil {
		fmt.Println("newBlock, %", err)
	}

	if _, err := InitGenesisGovernParam(common.ZeroHash, chain.SnapDB, 2048); err != nil {
		t.Error("InitGenesisGovernParam, error", err)
	}

	RegisterGovernParamVerifiers()

	if err := AddActiveVersion(params.GenesisVersion, 0, chain.StateDB); err != nil {
		t.Error("AddActiveVersion, err", err)
	}
	commit_sndb(chain)

	//the contract will retrieve this txHash as ProposalID
	prepair_sndb(chain)

	return chain
}

func submitText(t *testing.T, chain *mock.Chain) Proposal {
	p := &TextProposal{
		PIPID:        tpPIPID,
		ProposalType: Text,
		SubmitBlock:  chain.CurrentHeader().Number.Uint64(),
		ProposalID:   tpProposalID,
		Proposer:     nodeID,
	}
	if err := Submit(sender, p, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), NewMockStaking(), chain.StateDB, chainID); err != nil {
		t.Error("submitText, err", err)
	}
	return p
}

func submitVersion(t *testing.T, chain *mock.Chain, stk *MockStaking) Proposal {
	p := &VersionProposal{
		PIPID:           vpPIPID,
		ProposalType:    Version,
		SubmitBlock:     chain.CurrentHeader().Number.Uint64(),
		ProposalID:      vpProposalID,
		Proposer:        nodeID,
		NewVersion:      tempActiveVersion,
		EndVotingRounds: vpEndVotingRounds,
	}
	if err := Submit(sender, p, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), stk, chain.StateDB, chainID); err != nil {
		t.Error("submitVersion, err", err)
	}
	return p
}

func Max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}

func TestGov_GetVersionForStaking_No_PreActiveVersion(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	version := GetVersionForStaking(chain.CurrentHeader().Hash(), chain.StateDB)
	assert.Equal(t, params.GenesisVersion, version)
}

func TestGov_GetVersionForStaking_With_PreActiveVersion(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	if err := SetPreActiveVersion(chain.CurrentHeader().Hash(), tempActiveVersion); err != nil {
		t.Error("SetPreActiveVersion, err", err)
	}
	version := GetVersionForStaking(chain.CurrentHeader().Hash(), chain.StateDB)
	assert.Equal(t, tempActiveVersion, version)
}
func TestGov_GetCurrentActiveVersion(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	version := GetCurrentActiveVersion(chain.StateDB)
	t.Log("version", version)
	assert.Equal(t, params.GenesisVersion, version)
}

func TestGov_GetCurrentActiveVersion_NewActiveVersion(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	commit_sndb(chain)
	skip_emptyBlock(chain, 1001)

	if err := AddActiveVersion(tempActiveVersion, 1002, chain.StateDB); err != nil {
		t.Error("AddActiveVersion, err", err)
	}

	version := GetCurrentActiveVersion(chain.StateDB)
	t.Log("1st version", version, "chain.CurrentHeader().Number", chain.CurrentHeader().Number.Uint64())
	assert.Equal(t, tempActiveVersion, version)

	skip_emptyBlock(chain, 1002)
	version = GetCurrentActiveVersion(chain.StateDB)
	t.Log("2nd version", version, "chain.CurrentHeader().Number", chain.CurrentHeader().Number.Uint64())
	assert.Equal(t, tempActiveVersion, version)
}

func TestGov_Submit(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(6), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

	commit_sndb(chain)
	prepair_sndb(chain)
	submitText(t, chain)

	if tp, err := FindVotingProposal(chain.CurrentHeader().Hash(), chain.StateDB, Text); err != nil {
		t.Error("FindVotingProposal, err", err)
	} else {
		assert.Equal(t, tpProposalID, tp.(*TextProposal).ProposalID)
	}
}

func TestGov_Vote(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	submitText(t, chain)

	commit_sndb(chain)
	prepair_sndb(chain)

	versionSign := common.BytesToVersionSign(sign(params.GenesisVersion))

	vi := VoteInfo{ProposalID: tpProposalID, VoteNodeID: nodeID, VoteOption: Yes}
	if err := Vote(sender, vi, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), params.GenesisVersion, versionSign, NewMockStaking(), chain.StateDB); err != nil {
		t.Error("Vote, err", err)
		return
	}

	commit_sndb(chain)
	prepair_sndb(chain)

	if vvList, err := ListVoteValue(tpProposalID, chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 1, len(vvList))
	}

	if avList, err := ListAccuVerifier(chain.CurrentHeader().Hash(), tpProposalID); err != nil {
		t.Error("ListAccuVerifier, err", err)
	} else {
		assert.Equal(t, 1, len(avList))
		assert.Equal(t, nodeID, avList[0])
	}
}

// no voting proposal, no pre-active proposal
func TestGov_DeclareVersion_1(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	stk := NewMockStaking()

	versionSign := common.BytesToVersionSign(sign(params.GenesisVersion))
	if err := DeclareVersion(sender, nodeID, params.GenesisVersion, versionSign, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), stk, chain.StateDB); err != nil {
		t.Error("DeclareVersion, err", err)
	} else {
		anList := stk.ListDeclaredNode()
		assert.Equal(t, 1, len(anList))
		assert.Equal(t, params.GenesisVersion, anList[nodeID])
	}
}

// no voting proposal, there's pre-active proposal
func TestGov_DeclareVersion_2(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	if err := SetPreActiveVersion(chain.CurrentHeader().Hash(), tempActiveVersion); err != nil {
		t.Error("SetPreActiveVersion, err", err)
	}

	stk := NewMockStaking()

	versionSign := common.BytesToVersionSign(sign(params.GenesisVersion))
	err := DeclareVersion(sender, nodeID, params.GenesisVersion, versionSign, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), stk, chain.StateDB)
	assert.Equal(t, DeclareVersionError, err)

	err = DeclareVersion(sender, nodeID, tempActiveVersion, versionSign, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), stk, chain.StateDB)
	assert.Equal(t, VersionSignError, err)

	versionSign = common.BytesToVersionSign(sign(tempActiveVersion))
	if err := DeclareVersion(sender, nodeID, tempActiveVersion, versionSign, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), stk, chain.StateDB); err != nil {
		t.Error("DeclareVersion, err", err)
	} else {
		anList := stk.ListDeclaredNode()
		assert.Equal(t, 1, len(anList))
		assert.Equal(t, tempActiveVersion, anList[nodeID])
	}
}

// thers's voting proposal, no pre-active proposal
func TestGov_DeclareVersion_3(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	stk := NewMockStaking()

	submitVersion(t, chain, stk)

	commit_sndb(chain)
	prepair_sndb(chain)

	versionSign := common.BytesToVersionSign(sign(params.GenesisVersion))
	if err := DeclareVersion(sender, nodeID, params.GenesisVersion, versionSign, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), stk, chain.StateDB); err != nil {
		t.Error("DeclareVersion, err", err)
	} else {
		anList := stk.ListDeclaredNode()
		assert.Equal(t, 1, len(anList))
		assert.Equal(t, params.GenesisVersion, anList[nodeID])
	}

	//declared new version, gov will save all these nodes and notify staking if proposal is passed
	versionSign = common.BytesToVersionSign(sign(tempActiveVersion))
	if err := DeclareVersion(sender, nodeID, tempActiveVersion, versionSign, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), stk, chain.StateDB); err != nil {
		t.Error("DeclareVersion, err", err)
	} else {
		if anList, err := GetActiveNodeList(chain.CurrentHeader().Hash(), vpProposalID); err != nil {
			t.Error("DeclareVersion, err", err)
		} else {
			assert.Equal(t, 1, len(anList))
			assert.Equal(t, nodeID, anList[0])
		}

	}

	err := DeclareVersion(sender, nodeID, params.GenesisVersion, versionSign, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), stk, chain.StateDB)
	assert.Equal(t, VersionSignError, err)
}

func TestGov_ListProposal(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	stk := NewMockStaking()

	submitVersion(t, chain, stk)

	commit_sndb(chain)
	prepair_sndb(chain)
	if pList, err := ListProposal(chain.CurrentHeader().Hash(), chain.StateDB); err != nil {
		t.Error("ListProposal, err", err)
	} else {
		assert.Equal(t, 1, len(pList))
		assert.Equal(t, vpProposalID, pList[0].GetProposalID())
	}
}

func TestGov_ListVotingProposalID(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	stk := NewMockStaking()

	submitVersion(t, chain, stk)

	commit_sndb(chain)
	prepair_sndb(chain)

	if pIDList, err := ListVotingProposalID(chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVotingProposalID, err", err)
	} else {
		assert.Equal(t, 1, len(pIDList))
		assert.Equal(t, vpProposalID, pIDList[0])
	}
}

func TestGov_FindVotingProposal(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	stk := NewMockStaking()

	submitVersion(t, chain, stk)

	commit_sndb(chain)
	prepair_sndb(chain)
	if p, err := FindVotingProposal(chain.CurrentHeader().Hash(), chain.StateDB, Version); err != nil {
		t.Error("FindVotingProposal, err", err)
	} else {
		assert.NotNil(t, p)
		assert.Equal(t, vpProposalID, p.GetProposalID())
	}
}

func TestGov_GetMaxEndVotingBlock(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	stk := NewMockStaking()

	tp := submitText(t, chain)

	commit_sndb(chain)
	prepair_sndb(chain)

	vp := submitVersion(t, chain, stk)

	commit_sndb(chain)
	prepair_sndb(chain)

	versionSign := common.BytesToVersionSign(sign(params.GenesisVersion))

	vi := VoteInfo{ProposalID: tpProposalID, VoteNodeID: nodeID, VoteOption: Yes}
	if err := Vote(sender, vi, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), params.GenesisVersion, versionSign, stk, chain.StateDB); err != nil {
		t.Error("Vote, err", err)
		return
	}

	vi = VoteInfo{ProposalID: vpProposalID, VoteNodeID: nodeID, VoteOption: Yes}
	versionSign = common.BytesToVersionSign(sign(tempActiveVersion))
	if err := Vote(sender, vi, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), tempActiveVersion, versionSign, stk, chain.StateDB); err != nil {
		t.Error("Vote, err", err)
		return
	}

	if maxBlockNumber, err := GetMaxEndVotingBlock(nodeID, chain.CurrentHeader().Hash(), chain.StateDB); err != nil {
		t.Error("FindVotingProposal, err", err)
	} else {
		t.Log("maxBlockNumber", maxBlockNumber, "tp.GetEndVotingBlock()", tp.GetEndVotingBlock(), "vp.GetEndVotingBlock()", vp.GetEndVotingBlock())
		assert.Equal(t, maxBlockNumber, Max(tp.GetEndVotingBlock(), vp.GetEndVotingBlock()))
	}
}

func TestGov_NotifyPunishedVerifiers(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	stk := NewMockStaking()

	submitText(t, chain)

	commit_sndb(chain)
	prepair_sndb(chain)

	submitVersion(t, chain, stk)

	commit_sndb(chain)
	prepair_sndb(chain)

	versionSign := common.BytesToVersionSign(sign(params.GenesisVersion))

	vi := VoteInfo{ProposalID: tpProposalID, VoteNodeID: nodeID, VoteOption: Yes}
	if err := Vote(sender, vi, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), params.GenesisVersion, versionSign, stk, chain.StateDB); err != nil {
		t.Error("Vote, err", err)
		return
	}

	vi = VoteInfo{ProposalID: vpProposalID, VoteNodeID: nodeID, VoteOption: Yes}
	versionSign = common.BytesToVersionSign(sign(tempActiveVersion))
	if err := Vote(sender, vi, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), tempActiveVersion, versionSign, stk, chain.StateDB); err != nil {
		t.Error("Vote, err", err)
		return
	}
	commit_sndb(chain)
	prepair_sndb(chain)

	if vvList, err := ListVoteValue(tpProposalID, chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 1, len(vvList))
	}

	if vvList, err := ListVoteValue(vpProposalID, chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 1, len(vvList))
	}

	punishedVerifierMap := make(map[discover.NodeID]struct{})
	punishedVerifierMap[nodeID] = struct{}{}

	if err := NotifyPunishedVerifiers(chain.CurrentHeader().Hash(), punishedVerifierMap, chain.StateDB); err != nil {
		t.Error("NotifyPunishedVerifiers, err", err)
	}

	if vvList, err := ListVoteValue(tpProposalID, chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 0, len(vvList))
	}

	if vvList, err := ListVoteValue(vpProposalID, chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 0, len(vvList))
	}
}

func TestGov_SetGovernParam(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	commit_sndb(chain)
	prepair_sndb(chain)

	if err := SetGovernParam("myModule", "myName", "myDesc", "initValue", uint64(2), chain.CurrentHeader().Hash()); err != nil {
		t.Error("SetGovernParam, err", err)
	}

	commit_sndb(chain)
	prepair_sndb(chain)

	if gp, err := FindGovernParam("myModule", "myName", chain.CurrentHeader().Hash()); err != nil {
		t.Error("FindGovernParam, err", err)
	} else {
		assert.NotNil(t, gp)
		assert.Equal(t, "myModule", gp.ParamItem.Module)
		assert.Equal(t, "myName", gp.ParamItem.Name)
		assert.Equal(t, "initValue", gp.ParamValue.Value)
		assert.Equal(t, uint64(2), gp.ParamValue.ActiveBlock)
	}
}

func TestGov_UpdateGovernParam(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	commit_sndb(chain)
	prepair_sndb(chain)

	if err := SetGovernParam("myModule", "myName", "myDesc", "initValue", uint64(2), chain.CurrentHeader().Hash()); err != nil {
		t.Error("SetGovernParam, err", err)
	}
	commit_sndb(chain)
	prepair_sndb(chain)

	if err := UpdateGovernParamValue("myModule", "myName", "newValue", uint64(4), chain.CurrentHeader().Hash()); err != nil {
		t.Error("UpdateGovernParamValue, err", err)
	}

	if gp, err := FindGovernParam("myModule", "myName", chain.CurrentHeader().Hash()); err != nil {
		t.Error("FindGovernParam, err", err)
	} else {
		assert.NotNil(t, gp)
		assert.Equal(t, "myModule", gp.ParamItem.Module)
		assert.Equal(t, "myName", gp.ParamItem.Name)
		assert.Equal(t, "newValue", gp.ParamValue.Value)
		assert.Equal(t, uint64(4), gp.ParamValue.ActiveBlock)
	}
}

func TestGov_ListGovernParam(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	commit_sndb(chain)
	prepair_sndb(chain)

	if err := SetGovernParam("myModule", "myName", "myDesc", "initValue", uint64(2), chain.CurrentHeader().Hash()); err != nil {
		t.Error("SetGovernParam, err", err)
	}
	commit_sndb(chain)
	prepair_sndb(chain)

	if err := SetGovernParam("myModule", "myName2", "myDesc2", "initValue2", uint64(3), chain.CurrentHeader().Hash()); err != nil {
		t.Error("SetGovernParam, err", err)
	}

	commit_sndb(chain)
	prepair_sndb(chain)

	if err := SetGovernParam("myModule3", "myName3", "myDesc3", "initValue3", uint64(4), chain.CurrentHeader().Hash()); err != nil {
		t.Error("SetGovernParam, err", err)
	}

	if gpList, err := ListGovernParam("myModule", chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListGovernParam, err", err)
	} else {
		assert.NotNil(t, gpList)
		assert.Equal(t, 2, len(gpList))

		assert.Equal(t, "myModule", gpList[0].ParamItem.Module)
		assert.Equal(t, "myName", gpList[0].ParamItem.Name)
		assert.Equal(t, "initValue", gpList[0].ParamValue.Value)
		assert.Equal(t, uint64(2), gpList[0].ParamValue.ActiveBlock)

		assert.Equal(t, "myModule", gpList[1].ParamItem.Module)
		assert.Equal(t, "myName2", gpList[1].ParamItem.Name)
		assert.Equal(t, "initValue2", gpList[1].ParamValue.Value)
		assert.Equal(t, uint64(3), gpList[1].ParamValue.ActiveBlock)
	}
}

func TestGov_UpdateGovernParamValue(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	commit_sndb(chain)
	prepair_sndb(chain)

	if err := SetGovernParam("myModule", "myName", "myDesc", "initValue", uint64(2), chain.CurrentHeader().Hash()); err != nil {
		t.Error("SetGovernParam, err", err)
	}
	commit_sndb(chain)
	prepair_sndb(chain)

	if err := SetGovernParam("myModule", "myName2", "myDesc2", "initValue2", uint64(3), chain.CurrentHeader().Hash()); err != nil {
		t.Error("SetGovernParam, err", err)
	}

	commit_sndb(chain)
	prepair_sndb(chain)

	if err := SetGovernParam("myModule3", "myName3", "myDesc3", "initValue3", uint64(4), chain.CurrentHeader().Hash()); err != nil {
		t.Error("SetGovernParam, err", err)
	}

	if gpList, err := ListGovernParam("myModule", chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListGovernParam, err", err)
	} else {
		assert.NotNil(t, gpList)
		assert.Equal(t, 2, len(gpList))

		assert.Equal(t, "myModule", gpList[0].ParamItem.Module)
		assert.Equal(t, "myName", gpList[0].ParamItem.Name)
		assert.Equal(t, "initValue", gpList[0].ParamValue.Value)
		assert.Equal(t, uint64(2), gpList[0].ParamValue.ActiveBlock)

		assert.Equal(t, "myModule", gpList[1].ParamItem.Module)
		assert.Equal(t, "myName2", gpList[1].ParamItem.Name)
		assert.Equal(t, "initValue2", gpList[1].ParamValue.Value)
		assert.Equal(t, uint64(3), gpList[1].ParamValue.ActiveBlock)
	}
}

func TestGov_GetGovernParamValue(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)

	commit_sndb(chain)
	prepair_sndb(chain)

	if err := SetGovernParam("myModule", "myName", "myDesc", "initValue", uint64(2), chain.CurrentHeader().Hash()); err != nil {
		t.Error("SetGovernParam, err", err)
	}
	commit_sndb(chain)
	prepair_sndb(chain)

	if gpv, err := GetGovernParamValue("myModule", "myName", uint64(2), chain.CurrentHeader().Hash()); err != nil {
		t.Error("GetGovernParamValue, err", err)
	} else {
		assert.NotNil(t, gpv)
		assert.Equal(t, "initValue", gpv)
	}
}

func TestGov_GovernStakeThreshold(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)
	if threshold, err := GovernStakeThreshold(1, chain.CurrentHeader().Hash()); err != nil {
		t.Error("GovernStakeThreshold, err", err)
	} else {
		assert.NotNil(t, threshold)
		assert.Equal(t, xcom.StakeThreshold().String(), threshold.String())
	}
}

func TestGov_GovernOperatingThreshold(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)
	if threshold, err := GovernOperatingThreshold(1, chain.CurrentHeader().Hash()); err != nil {
		t.Error("GovernOperatingThreshold, err", err)
	} else {
		assert.NotNil(t, threshold)
		assert.Equal(t, xcom.OperatingThreshold().String(), threshold.String())
	}
}

func TestGov_GovernMaxValidators(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)
	if threshold, err := GovernMaxValidators(1, chain.CurrentHeader().Hash()); err != nil {
		t.Error("GovernMaxValidators, err", err)
	} else {
		assert.NotNil(t, threshold)
		assert.Equal(t, xcom.MaxValidators(), threshold)
	}
}

func TestGov_GovernUnStakeFreezeDuration(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)
	if threshold, err := GovernUnStakeFreezeDuration(1, chain.CurrentHeader().Hash()); err != nil {
		t.Error("GovernUnStakeFreezeDuration, err", err)
	} else {
		assert.NotNil(t, threshold)
		assert.Equal(t, xcom.UnStakeFreezeDuration(), threshold)
	}
}

func TestGov_GovernSlashFractionDuplicateSign(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)
	if threshold, err := GovernSlashFractionDuplicateSign(1, chain.CurrentHeader().Hash()); err != nil {
		t.Error("GovernSlashFractionDuplicateSign, err", err)
	} else {
		assert.NotNil(t, threshold)
		assert.Equal(t, xcom.SlashFractionDuplicateSign(), threshold)
	}
}

func TestGov_GovernDuplicateSignReportReward(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)
	if threshold, err := GovernDuplicateSignReportReward(1, chain.CurrentHeader().Hash()); err != nil {
		t.Error("GovernDuplicateSignReportReward, err", err)
	} else {
		assert.NotNil(t, threshold)
		assert.Equal(t, xcom.DuplicateSignReportReward(), threshold)
	}
}

func TestGov_GovernMaxEvidenceAge(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)
	if threshold, err := GovernMaxEvidenceAge(1, chain.CurrentHeader().Hash()); err != nil {
		t.Error("GovernMaxEvidenceAge, err", err)
	} else {
		assert.NotNil(t, threshold)
		assert.Equal(t, xcom.MaxEvidenceAge(), threshold)
	}
}

func TestGov_GovernSlashBlocksReward(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)
	if threshold, err := GovernSlashBlocksReward(1, chain.CurrentHeader().Hash()); err != nil {
		t.Error("GovernSlashBlocksReward, err", err)
	} else {
		assert.NotNil(t, threshold)
		assert.Equal(t, xcom.SlashBlocksReward(), threshold)
	}
}

func TestGov_GovernMaxBlockGasLimit(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)
	if threshold, err := GovernMaxBlockGasLimit(1, chain.CurrentHeader().Hash()); err != nil {
		t.Error("GovernMaxBlockGasLimit, err", err)
	} else {
		assert.NotNil(t, threshold)
		assert.Equal(t, int(params.DefaultMinerGasCeil), threshold)
	}
}

func TestGov_ClearProcessingProposals(t *testing.T) {
	chain := setup(t)
	defer clear(chain, t)
	commit_sndb(chain)
	prepair_sndb(chain)

	stk := NewMockStaking()

	submitText(t, chain)

	commit_sndb(chain)
	prepair_sndb(chain)

	submitVersion(t, chain, stk)

	commit_sndb(chain)
	prepair_sndb(chain)

	versionSign := common.BytesToVersionSign(sign(params.GenesisVersion))

	vi := VoteInfo{ProposalID: tpProposalID, VoteNodeID: nodeID, VoteOption: Yes}
	if err := Vote(sender, vi, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), params.GenesisVersion, versionSign, stk, chain.StateDB); err != nil {
		t.Error("Vote, err", err)
		return
	}

	vi = VoteInfo{ProposalID: vpProposalID, VoteNodeID: nodeID, VoteOption: Yes}
	versionSign = common.BytesToVersionSign(sign(tempActiveVersion))
	if err := Vote(sender, vi, chain.CurrentHeader().Hash(), chain.CurrentHeader().Number.Uint64(), tempActiveVersion, versionSign, stk, chain.StateDB); err != nil {
		t.Error("Vote, err", err)
		return
	}
	commit_sndb(chain)
	prepair_sndb(chain)

	//-------
	//-------
	//-------
	if votinglist, err := ListVotingProposalID(chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVotingProposalID, err", err)
	} else {
		assert.Equal(t, 2, len(votinglist))
	}

	if endList, err := ListEndProposalID(chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListEndProposalID, err", err)
	} else {
		assert.Equal(t, 0, len(endList))
	}

	//-------
	if vvList, err := ListVoteValue(tpProposalID, chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 1, len(vvList))
	}

	if avList, err := ListAccuVerifier(chain.CurrentHeader().Hash(), tpProposalID); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 1, len(avList))
	}

	//-------
	if vvList, err := ListVoteValue(vpProposalID, chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 1, len(vvList))
	}
	if avList, err := ListAccuVerifier(chain.CurrentHeader().Hash(), vpProposalID); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 1, len(avList))
	}

	//------------------------------------------
	//------------------------------------------
	//------------------------------------------
	if err := ClearProcessingProposals(chain.CurrentHeader().Hash(), chain.StateDB); err != nil {
		t.Error("ClearProcessingProposals, err", err)
	}

	//-------
	//-------
	//-------
	if votinglist, err := ListVotingProposalID(chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVotingProposalID, err", err)
	} else {
		assert.Equal(t, 0, len(votinglist))
	}

	if endList, err := ListEndProposalID(chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListEndProposalID, err", err)
	} else {
		assert.Equal(t, 2, len(endList))
	}

	//-------
	if vvList, err := ListVoteValue(tpProposalID, chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 0, len(vvList))
	}

	if avList, err := ListAccuVerifier(chain.CurrentHeader().Hash(), tpProposalID); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 0, len(avList))
	}

	//-------
	if vvList, err := ListVoteValue(vpProposalID, chain.CurrentHeader().Hash()); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 0, len(vvList))
	}
	if avList, err := ListAccuVerifier(chain.CurrentHeader().Hash(), vpProposalID); err != nil {
		t.Error("ListVoteValue, err", err)
	} else {
		assert.Equal(t, 0, len(avList))
	}
}
