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

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

// in genesis.json, the config value of the supportRate( and voteRate...) is the real value * 10000.
// the RateCoefficient is used to calculate the supportRate (and voteRate) of a proposal correctly.
const RateCoefficient = uint64(10000)

type ProposalType uint8

const (
	Text    ProposalType = 0x01
	Version ProposalType = 0x02
	Param   ProposalType = 0x03
	Cancel  ProposalType = 0x04
)

type ProposalStatus uint8

const (
	Voting    ProposalStatus = 0x01
	Pass      ProposalStatus = 0x02
	Failed    ProposalStatus = 0x03
	PreActive ProposalStatus = 0x04
	Active    ProposalStatus = 0x05
	Canceled  ProposalStatus = 0x06
)

func (status ProposalStatus) ToString() string {
	switch status {
	case Voting:
		return "Voting"
	case Pass:
		return "Pass"
	case Failed:
		return "Failed"
	case PreActive:
		return "PreActive"
	case Active:
		return "Active"
	case Canceled:
		return "Canceled"
	default: //default case
		return ""
	}
}

type VoteOption uint8

const (
	Yes        VoteOption = 0x01
	No         VoteOption = 0x02
	Abstention VoteOption = 0x03
)

func ParseVoteOption(option uint8) VoteOption {
	switch option {
	case 0x01:
		return Yes
	case 0x02:
		return No
	case 0x03:
		return Abstention
	}
	return Abstention
}

type Proposal interface {
	GetProposalID() common.Hash
	GetProposalType() ProposalType
	GetPIPID() string
	GetSubmitBlock() uint64
	GetEndVotingBlock() uint64
	GetProposer() discover.NodeID
	GetTallyResult() TallyResult
	Verify(blockNumber uint64, blockHash common.Hash, state xcom.StateDB, chainID *big.Int) error
	String() string
}

type TextProposal struct {
	ProposalID     common.Hash
	ProposalType   ProposalType
	PIPID          string
	SubmitBlock    uint64
	EndVotingBlock uint64
	Proposer       discover.NodeID
	Result         TallyResult `json:"-"`
}

func (tp *TextProposal) GetProposalID() common.Hash {
	return tp.ProposalID
}

func (tp *TextProposal) GetProposalType() ProposalType {
	return tp.ProposalType
}

func (tp *TextProposal) GetPIPID() string {
	return tp.PIPID
}

func (tp *TextProposal) GetSubmitBlock() uint64 {
	return tp.SubmitBlock
}

func (tp *TextProposal) GetEndVotingBlock() uint64 {
	return tp.EndVotingBlock
}

func (tp *TextProposal) GetProposer() discover.NodeID {
	return tp.Proposer
}

func (tp *TextProposal) GetTallyResult() TallyResult {
	return tp.Result
}

func (tp *TextProposal) Verify(submitBlock uint64, blockHash common.Hash, state xcom.StateDB, chainID *big.Int) error {
	if tp.ProposalType != Text {
		return ProposalTypeError
	}

	if err := verifyBasic(tp, blockHash, state); err != nil {
		return err
	}

	endVotingBlock := xutil.CalEndVotingBlock(submitBlock, xutil.EstimateConsensusRoundsForGov(xcom.TextProposalVote_DurationSeconds()))
	if endVotingBlock <= submitBlock {
		log.Error("the end-voting-block is lower than submit-block. Please check configuration")
		return common.InternalError
	}
	tp.EndVotingBlock = endVotingBlock

	log.Debug("verify Text Proposal", "PIPID", tp.PIPID, "voteDuration", xcom.TextProposalVote_DurationSeconds(), "endVotingBlock", endVotingBlock, "blockNumber", submitBlock, "blockHash", blockHash)
	return nil
}

func (tp *TextProposal) String() string {
	return fmt.Sprintf(`Proposal %x: 
  Type:               	%x
  PIPID:			    %s
  Proposer:            	%x
  SubmitBlock:        	%d
  EndVotingBlock:   	%d`,
		tp.ProposalID, tp.ProposalType, tp.PIPID, tp.Proposer, tp.SubmitBlock, tp.EndVotingBlock)
}

type VersionProposal struct {
	ProposalID      common.Hash
	ProposalType    ProposalType
	PIPID           string
	SubmitBlock     uint64
	EndVotingRounds uint64
	EndVotingBlock  uint64
	Proposer        discover.NodeID
	Result          TallyResult `json:"-"`
	NewVersion      uint32
	ActiveBlock     uint64
}

func (vp *VersionProposal) GetProposalID() common.Hash {
	return vp.ProposalID
}

func (vp *VersionProposal) GetProposalType() ProposalType {
	return vp.ProposalType
}

func (vp *VersionProposal) GetPIPID() string {
	return vp.PIPID
}

func (vp *VersionProposal) GetSubmitBlock() uint64 {
	return vp.SubmitBlock
}

func (vp *VersionProposal) GetEndVotingBlock() uint64 {
	return vp.EndVotingBlock
}

func (vp *VersionProposal) GetProposer() discover.NodeID {
	return vp.Proposer
}

func (vp *VersionProposal) GetTallyResult() TallyResult {
	return vp.Result
}

func (vp *VersionProposal) GetNewVersion() uint32 {
	return vp.NewVersion
}

func (vp *VersionProposal) GetActiveBlock() uint64 {
	return vp.ActiveBlock
}

func (vp *VersionProposal) Verify(submitBlock uint64, blockHash common.Hash, state xcom.StateDB, chainID *big.Int) error {

	if vp.ProposalType != Version {
		return ProposalTypeError
	}
	if vp.EndVotingRounds <= 0 {
		return EndVotingRoundsTooSmall
	}

	if vp.EndVotingRounds > xutil.EstimateConsensusRoundsForGov(xcom.VersionProposalVote_DurationSeconds()) {
		return EndVotingRoundsTooLarge
	}

	if err := verifyBasic(vp, blockHash, state); err != nil {
		return err
	}

	endVotingBlock := xutil.CalEndVotingBlock(submitBlock, vp.EndVotingRounds)
	if endVotingBlock <= submitBlock {
		log.Error("the end-voting-block is lower than submit-block. Please check configuration")
		return common.InternalError
	}
	activeBlock := xutil.CalActiveBlock(endVotingBlock)

	vp.EndVotingBlock = endVotingBlock
	vp.ActiveBlock = activeBlock

	if vp.NewVersion <= 0 || vp.NewVersion>>8 <= uint32(GetCurrentActiveVersion(state))>>8 {
		return NewVersionError
	}

	if exist, err := FindVotingProposal(blockHash, state, Version, Param); err != nil {
		return err
	} else if exist != nil {
		if exist.GetProposalType() == Version {
			return VotingVersionProposalExist
		} else {
			return VotingParamProposalExist
		}
	}

	//another VersionProposal in Pre-active process，exit
	proposalID, err := GetPreActiveProposalID(blockHash)
	if err != nil {
		log.Error("check pre-active version proposal error", "blockNumber", submitBlock, "blockHash", blockHash)
		return err
	}
	if proposalID != common.ZeroHash {
		return PreActiveVersionProposalExist
	}

	log.Debug("verify Version Proposal", "PIPID", vp.PIPID, "voteDuration", xcom.VersionProposalVote_DurationSeconds(), "endVotingBlock", endVotingBlock, "activeBlock", activeBlock, "blockNumber", submitBlock, "blockHash", blockHash)
	return nil
}

func (vp *VersionProposal) String() string {
	return fmt.Sprintf(`Proposal %x: 
  Type:               	%x
  PIPID:			    %s
  Proposer:            	%x
  SubmitBlock:        	%d
  EndVotingBlock:   	%d
  ActiveBlock:   		%d
  NewVersion:   		%d`,
		vp.ProposalID, vp.ProposalType, vp.PIPID, vp.Proposer, vp.SubmitBlock, vp.EndVotingBlock, vp.ActiveBlock, vp.NewVersion)
}

type CancelProposal struct {
	ProposalID      common.Hash
	ProposalType    ProposalType
	PIPID           string
	SubmitBlock     uint64
	EndVotingRounds uint64
	EndVotingBlock  uint64
	Proposer        discover.NodeID
	TobeCanceled    common.Hash
	Result          TallyResult `json:"-"`
}

func (cp *CancelProposal) GetProposalID() common.Hash {
	return cp.ProposalID
}

func (cp *CancelProposal) GetProposalType() ProposalType {
	return cp.ProposalType
}

func (cp *CancelProposal) GetPIPID() string {
	return cp.PIPID
}

func (cp *CancelProposal) GetSubmitBlock() uint64 {
	return cp.SubmitBlock
}

func (cp *CancelProposal) GetEndVotingBlock() uint64 {
	return cp.EndVotingBlock
}

func (cp *CancelProposal) GetProposer() discover.NodeID {
	return cp.Proposer
}

func (cp *CancelProposal) GetTallyResult() TallyResult {
	return cp.Result
}

func (cp *CancelProposal) Verify(submitBlock uint64, blockHash common.Hash, state xcom.StateDB, chainID *big.Int) error {
	if cp.ProposalType != Cancel {
		return ProposalTypeError
	}

	if err := verifyBasic(cp, blockHash, state); err != nil {
		return err
	}

	if cp.EndVotingRounds <= 0 {
		return EndVotingRoundsTooSmall
	}

	endVotingBlock := xutil.CalEndVotingBlock(submitBlock, cp.EndVotingRounds)
	if endVotingBlock <= submitBlock {
		log.Error("the end-voting-block is lower than submit-block. Please check configuration")
		return common.InternalError
	}
	cp.EndVotingBlock = endVotingBlock

	if exist, err := FindVotingProposal(blockHash, state, Cancel); err != nil {
		log.Error("find voting cancel proposal error", "err", err)
		return err
	} else if exist != nil {
		return VotingCancelProposalExist
	}

	if tobeCanceled, err := GetProposal(cp.TobeCanceled, state); err != nil {
		log.Error("find to be canceled version proposal error", "err", err)
		return err
	} else if tobeCanceled == nil {
		return TobeCanceledProposalNotFound
	} else if tobeCanceled.GetProposalType() != Version && tobeCanceled.GetProposalType() != Param {
		return TobeCanceledProposalTypeError
	} else if votingList, err := ListVotingProposal(blockHash); err != nil {
		log.Error("list voting proposal error", "err", err)
		return err
	} else if !xutil.InHashList(cp.TobeCanceled, votingList) {
		return TobeCanceledProposalNotAtVoting
	} else if cp.EndVotingBlock >= tobeCanceled.GetEndVotingBlock() {
		return EndVotingRoundsTooLarge
	}
	log.Debug("verify Cancel Proposal", "PIPID", cp.PIPID, "endVotingBlock", endVotingBlock, "blockNumber", submitBlock, "blockHash", blockHash)
	return nil
}

func (cp *CancelProposal) String() string {
	return fmt.Sprintf(`Proposal %x: 
  Type:               	%x
  PIPID:			    %s
  Proposer:            	%x
  SubmitBlock:        	%d
  EndVotingBlock:   	%d
  TobeCanceled:   		%s`,
		cp.ProposalID, cp.ProposalType, cp.PIPID, cp.Proposer, cp.SubmitBlock, cp.EndVotingBlock, cp.TobeCanceled.Hex())
}

type ParamProposal struct {
	ProposalID     common.Hash
	ProposalType   ProposalType
	PIPID          string
	SubmitBlock    uint64
	EndVotingBlock uint64
	Proposer       discover.NodeID
	Result         TallyResult `json:"-"`
	Module         string
	Name           string
	NewValue       string
}

func (pp *ParamProposal) GetProposalID() common.Hash {
	return pp.ProposalID
}

func (pp *ParamProposal) GetProposalType() ProposalType {
	return pp.ProposalType
}

func (pp *ParamProposal) GetPIPID() string {
	return pp.PIPID
}

func (pp *ParamProposal) GetSubmitBlock() uint64 {
	return pp.SubmitBlock
}

func (pp *ParamProposal) GetEndVotingBlock() uint64 {
	return pp.EndVotingBlock
}

func (pp *ParamProposal) GetProposer() discover.NodeID {
	return pp.Proposer
}

func (pp *ParamProposal) GetTallyResult() TallyResult {
	return pp.Result
}

func (pp *ParamProposal) Verify(submitBlock uint64, blockHash common.Hash, state xcom.StateDB, chainID *big.Int) error {
	if pp.ProposalType != Param {
		return ProposalTypeError
	}
	if err := verifyBasic(pp, blockHash, state); err != nil {
		return err
	}

	param, err := FindGovernParam(pp.Module, pp.Name, blockHash)
	if err != nil {
		log.Error("find govern parameter error", "err", err)
		return err
	} else if param == nil {
		return UnsupportedGovernParam
	} else if param.ParamValue.Value == pp.NewValue {
		return ParamProposalIsSameValue
	}

	if paramVerifier, ok := ParamVerifierMap[pp.Module+"/"+pp.Name]; ok {
		if err := paramVerifier(submitBlock, blockHash, pp.NewValue); err != nil {
			return err
		}
	} else {
		return UnsupportedGovernParam
	}

	if exist, err := FindVotingProposal(blockHash, state, Param, Version); err != nil {
		log.Error("find voting param proposal error", "err", err)
		return err
	} else if exist != nil {
		if exist.GetProposalType() == Param {
			return VotingParamProposalExist
		} else {
			return VotingVersionProposalExist
		}
	}

	//another VersionProposal in Pre-active process，exit
	proposalID, err := GetPreActiveProposalID(blockHash)
	if err != nil {
		log.Error("check pre-active version proposal error", "blockNumber", submitBlock, "blockHash", blockHash)
		return err
	}
	if proposalID != common.ZeroHash {
		return PreActiveVersionProposalExist
	}

	var voteDuration = xcom.ParamProposalVote_DurationSeconds()

	endVotingBlock := xutil.EstimateEndVotingBlockForParaProposal(submitBlock, voteDuration)
	if endVotingBlock <= submitBlock {
		log.Error("the end-voting-block is lower than submit-block. Please check configuration")
		return common.InternalError
	}
	pp.EndVotingBlock = endVotingBlock
	log.Debug("verify Parameter Proposal", "PIPID", pp.PIPID, "voteDuration", voteDuration, "endVotingBlock", endVotingBlock, "blockNumber", submitBlock, "blockHash", blockHash)

	return nil
}

func (pp *ParamProposal) String() string {
	return fmt.Sprintf(`Proposal %x: 
  Type:               	%x
  PIPID:			    %s
  Proposer:            	%x
  SubmitBlock:        	%d
  EndVotingBlock:   	%d
  Module:   			%s
  Name:   				%s
  NewValue:   			%s`,
		pp.ProposalID, pp.ProposalType, pp.PIPID, pp.Proposer, pp.SubmitBlock, pp.EndVotingBlock, pp.Module, pp.Name, pp.NewValue)
}

func verifyBasic(p Proposal, blockHash common.Hash, state xcom.StateDB) error {
	log.Debug("verify proposal basic parameters", "proposalID", p.GetProposalID(), "proposer", p.GetProposer(), "pipID", p.GetPIPID(), "endVotingBlock", p.GetEndVotingBlock(), "submitBlock", p.GetSubmitBlock())

	if p.GetProposalID() != common.ZeroHash {
		p, err := GetProposal(p.GetProposalID(), state)
		if err != nil {
			return err
		}
		if nil != p {
			return ProposalIDExist
		}
	} else {
		return ProposalIDEmpty
	}

	if p.GetProposer() == discover.ZeroNodeID {
		return ProposerEmpty
	}

	//if a PIPID is used in a proposal which is passed, this PIPID cannot be used in another proposal
	if len(p.GetPIPID()) == 0 {
		return PIPIDEmpty
	} else if pipIdList, err := ListPIPID(state); err != nil {
		log.Error("list PIPID error", "err", err)
		return err
	} else if isPIPIDExist(p.GetPIPID(), pipIdList) {
		return PIPIDExist
	}

	//if a PIPID is used in a proposal which is at voting stage, this PIPID cannot be used in another proposal
	if votingPIDList, err := ListVotingProposalID(blockHash); err != nil {
		log.Error("list voting proposal ID error", "err", err)
		return err
	} else {
		for _, votingPID := range votingPIDList {
			if exist, err := GetExistProposal(votingPID, state); err != nil {
				log.Error("get existing proposal error", "err", err)
				return err
			} else if exist.GetPIPID() == p.GetPIPID() {
				return PIPIDExist
			}
		}
	}

	return nil
}

func isPIPIDExist(pipID string, pipIDList []string) bool {
	for _, id := range pipIDList {
		if pipID == id {
			return true
		}
	}
	return false
}
