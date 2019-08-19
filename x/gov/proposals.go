package gov

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

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
	Verify(blockNumber uint64, blockHash common.Hash, state xcom.StateDB) error
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

func (tp *TextProposal) Verify(submitBlock uint64, blockHash common.Hash, state xcom.StateDB) error {
	if tp.ProposalType != Text {
		return common.NewBizError("Proposal Type error.")
	}

	if err := verifyBasic(tp, state); err != nil {
		return err
	}

	endVotingBlock := xutil.CalEndVotingBlock(submitBlock, xcom.TextProposalVote_ConsensusRounds())
	tp.EndVotingBlock = endVotingBlock

	log.Debug("text proposal", "endVotingBlock", tp.EndVotingBlock, "consensusSize", xutil.ConsensusSize(), "xcom.ElectionDistance()", xcom.ElectionDistance())
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

func (vp *VersionProposal) Verify(submitBlock uint64, blockHash common.Hash, state xcom.StateDB) error {

	if vp.ProposalType != Version {
		return common.NewBizError("Proposal Type error.")
	}

	if vp.EndVotingRounds > xcom.VersionProposalVote_ConsensusRounds() {
		return common.NewBizError("voting consensus rounds too large.")
	}

	if err := verifyBasic(vp, state); err != nil {
		return err
	}

	endVotingBlock := xutil.CalEndVotingBlock(submitBlock, vp.EndVotingRounds)

	activeBlock := xutil.CalActiveBlock(endVotingBlock)

	vp.EndVotingBlock = endVotingBlock
	vp.ActiveBlock = activeBlock

	if vp.NewVersion>>8 <= uint32(GetCurrentActiveVersion(state))>>8 {
		return common.NewBizError("New version should larger than current active version.")
	}

	if exist, err := FindVotingVersionProposal(blockHash, submitBlock, state); err != nil {
		return err
	} else if exist != nil {
		log.Error("there is another version proposal at voting stage", "proposalID", exist.ProposalID)
		return common.NewBizError("there is another version proposal at voting stage")
	}

	//another VersionProposal in Pre-active process，exit
	proposalID, err := GetPreActiveProposalID(blockHash)
	if err != nil {
		log.Error("to check if there's a pre-active version proposal failed.", "blockNumber", submitBlock, "blockHash", blockHash)
		return err
	}
	if proposalID != common.ZeroHash {
		return common.NewBizError("there is another pre-active version proposal")
	}

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

func (cp *CancelProposal) Verify(submitBlock uint64, blockHash common.Hash, state xcom.StateDB) error {
	if cp.ProposalType != Cancel {
		return common.NewBizError("Proposal Type error.")
	}

	if err := verifyBasic(cp, state); err != nil {
		return err
	}

	endVotingBlock := xutil.CalEndVotingBlock(submitBlock, cp.EndVotingRounds)
	cp.EndVotingBlock = endVotingBlock

	if exist, err := FindVotingCancelProposal(blockHash, submitBlock, state); err != nil {
		return err
	} else if exist != nil {
		log.Error("there is another cancel proposal at voting stage", "proposalID", exist.ProposalID)
		return common.NewBizError("there is another cancel proposal at voting stage")
	}

	if tobeCanceled, err := GetExistProposal(cp.TobeCanceled, state); err != nil {
		log.Error("find to be canceled version proposal error", "err", err)
		return common.NewBizError("find to be canceled version proposal error")
	} else if tobeCanceled.GetProposalType() != Version {
		return common.NewBizError("to be canceled proposal should be version proposal")
	} else if votingList, err := ListVotingProposal(blockHash); err != nil {
		return err
	} else if !xutil.InHashList(cp.TobeCanceled, votingList) {
		return common.NewBizError("to be canceled version proposal should be at voting stage")
	} else if cp.EndVotingBlock >= tobeCanceled.GetEndVotingBlock() {
		return common.NewBizError("voting consensus rounds too large.")
	}
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

func verifyBasic(p Proposal, state xcom.StateDB) error {
	log.Debug("verify proposal basic parameters", "proposalID", p.GetProposalID(), "proposer", p.GetProposer(), "pipID", p.GetPIPID(), "endVotingBlock", p.GetEndVotingBlock(), "submitBlock", p.GetSubmitBlock())

	if p.GetProposalID() != common.ZeroHash {
		p, err := GetProposal(p.GetProposalID(), state)
		if err != nil {
			return err
		}
		if nil != p {
			return common.NewBizError("ProposalID is already used.")
		}
	} else {
		return common.NewBizError("ProposalID is empty.")
	}

	if p.GetProposer() == discover.ZeroNodeID {
		return common.NewBizError("Proposer is empty.")
	}

	if len(p.GetPIPID()) == 0 {
		return common.NewBizError("PIPID is empty.")
	} else if pipIdList, err := ListPIPID(state); err != nil {
		return err
	} else if isPIPIDExist(p.GetPIPID(), pipIdList) {
		return common.NewBizError("PIPID is existing.")
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
