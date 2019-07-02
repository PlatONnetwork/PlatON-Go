package gov

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type ProposalType byte

const (
	Text    ProposalType = 0x01
	Version ProposalType = 0x02
)

type ProposalStatus byte

const (
	Voting    ProposalStatus = 0x01
	Pass      ProposalStatus = 0x02
	Failed    ProposalStatus = 0x03
	PreActive ProposalStatus = 0x04
	Active    ProposalStatus = 0x05
)

type VoteOption byte

const (
	Yes VoteOption = iota + 1
	No
	Abstention
)

type TallyResult struct {
	ProposalID    common.Hash    `json:"proposalID"`
	Yeas          uint16         `json:"yeas"`
	Nays          uint16         `json:"nays"`
	Abstentions   uint16         `json:"abstentions"`
	AccuVerifiers uint16         `json:"accuVerifiers"`
	Status        ProposalStatus `json:"status"`
}

type Vote struct {
	ProposalID common.Hash     `json:"proposalID"`
	VoteNodeID discover.NodeID `json:"voteNodeID"`
	VoteOption VoteOption      `json:"voteOption"`
}


type VoteValue struct {
	VoteNodeID discover.NodeID `json:"voteNodeID"`
	VoteOption VoteOption      `json:"voteOption"`
}

var MaxVotingDuration = uint64(14 * 24 * 60 * 60) / xcom.ConsensusSize * xcom.ConsensusSize

type Proposal interface {
	SetProposalID(proposalID common.Hash)
	GetProposalID() common.Hash

	SetGithubID(githubID string)
	GetGithubID() string

	SetTopic(topic string)
	GetTopic() string

	SetDesc(desc string)
	GetDesc() string

	SetProposalType(proposalType ProposalType)
	GetProposalType() ProposalType

	SetUrl(url string)
	GetUrl() string

	SetSubmitBlock(blockNumber uint64)
	GetSubmitBlock() uint64

	SetEndVotingBlock(blockNumber uint64)
	GetEndVotingBlock() uint64

	SetProposer(proposer discover.NodeID)
	GetProposer() discover.NodeID

	SetTallyResult(tallyResult TallyResult)
	GetTallyResult() TallyResult

	Verify(curBlockNum uint64, state xcom.StateDB) (error)

	String() string
}

type TextProposal struct {
	ProposalID     common.Hash
	GithubID       string
	ProposalType   ProposalType
	Topic          string
	Desc           string
	Url            string
	SubmitBlock    uint64
	EndVotingBlock uint64
	Proposer       discover.NodeID
	Result         TallyResult
}

func (tp TextProposal) SetProposalID(proposalID common.Hash) {
	tp.ProposalID = proposalID
}

func (tp TextProposal) GetProposalID() common.Hash {
	return tp.ProposalID
}

func (tp TextProposal) SetGithubID(githubID string) {
	tp.GithubID = githubID
}

func (tp TextProposal) GetGithubID() string {
	return tp.GithubID
}

func (tp TextProposal) SetProposalType(proposalType ProposalType) {
	tp.ProposalType = proposalType
}

func (tp TextProposal) GetProposalType() ProposalType {
	return tp.ProposalType
}

func (tp TextProposal) SetTopic(topic string) {
	tp.Topic = topic
}

func (tp TextProposal) GetTopic() string {
	return tp.Topic
}

func (tp TextProposal) SetDesc(desc string) {
	tp.Desc = desc
}

func (tp TextProposal) GetDesc() string {
	return tp.Desc
}

func (tp TextProposal) SetUrl(url string) {
	tp.Url = url
}

func (tp TextProposal) GetUrl() string {
	return tp.Url
}

func (tp TextProposal) SetSubmitBlock(blockNumber uint64) {
	tp.SubmitBlock = blockNumber
}

func (tp TextProposal) GetSubmitBlock() uint64 {
	return tp.SubmitBlock
}

func (tp TextProposal) SetEndVotingBlock(blockNumber uint64) {
	tp.EndVotingBlock = blockNumber
}

func (tp TextProposal) GetEndVotingBlock() uint64 {
	return tp.EndVotingBlock
}

func (tp TextProposal) SetProposer(proposer discover.NodeID) {
	tp.Proposer = proposer
}

func (tp TextProposal) GetProposer() discover.NodeID {
	return tp.Proposer
}

func (tp TextProposal) SetTallyResult(result TallyResult) {
	tp.Result = result
}

func (tp TextProposal) GetTallyResult() TallyResult {
	return tp.Result
}

func (tp TextProposal) Verify(curBlockNum uint64, state xcom.StateDB) (error) {
	return verifyBasic(tp.ProposalID, tp.Proposer, tp.ProposalType, tp.Topic, tp.Desc, tp.GithubID, tp.Url, tp.EndVotingBlock, curBlockNum, state)
}

func (tp TextProposal) String() string {
	return fmt.Sprintf(`Proposal %x: 
  GithubID:            	%s
  Topic:              	%s
  Type:               	%x
  Proposer:            	%x
  SubmitBlock:        	%d
  EndVotingBlock:   	%d`, tp.ProposalID, tp.GithubID, tp.Topic, tp.ProposalType, tp.Proposer, tp.SubmitBlock, tp.EndVotingBlock)
}

type VersionProposal struct {
	TextProposal
	NewVersion  uint32
	ActiveBlock uint64
}

func (vp VersionProposal) SetNewVersion(newVersion uint32) {
	vp.NewVersion = newVersion
}

func (vp VersionProposal) GetNewVersion() uint32 {
	return vp.NewVersion
}

func (vp VersionProposal) SetActiveBlock(activeBlock uint64) {
	vp.ActiveBlock = activeBlock
}

func (vp VersionProposal) GetActiveBlock() uint64 {
	return vp.ActiveBlock
}

func verifyBasic(proposalID common.Hash, proposer discover.NodeID, proposalType ProposalType, topic, desc, githubID, url  string, endVotingBlock uint64, curBlockNum uint64, state xcom.StateDB) (error){
	if len(proposalID) >0 {
		p, err := GovDBInstance().GetProposal(proposalID, state);
		if p == nil {
			return err
		}
		if nil != p {
			return common.NewBizError("ProposalID is already used.")
		}
	}else{
		return common.NewBizError("ProposalID is empty.")
	}

	if len(proposer) == 0 {
		return common.NewBizError("Proposer is empty.")
	}
	if proposalType != Version {
		return common.NewBizError("Proposal Type error.")
	}
	if len(topic) == 0 || len(topic) > 128 {
		return common.NewBizError("Topic is empty or the size is bigger than 128.")
	}
	if len(desc) > 512 {
		return common.NewBizError("description's size is bigger than 512.")
	}
	/*if len(vp.GithubID) == 0 || vp.GithubID == gov.govDB.GetProposal(vp.ProposalID, state).GetGithubID() {
		var err error = errors.New("[GOV] Verify(): GithubID empty or duplicated.")
		return false, err
	}
	if len(vp.Url) == 0 || vp.GithubID == gov.govDB.GetProposal(vp.ProposalID, state).GetUrl() {
		var err error = errors.New("[GOV] Verify(): Github URL empty or duplicated.")
		return false, err
	}*/

	if xutil.CalculateRound(endVotingBlock) - xutil.CalculateRound(curBlockNum) <= 0 || endVotingBlock > curBlockNum + MaxVotingDuration {
		return common.NewBizError("end voting block number invalid.")
	}

	return nil
}

func (vp VersionProposal) Verify(curBlockNum uint64, state xcom.StateDB) (error) {

	if err := verifyBasic(vp.ProposalID, vp.Proposer, vp.ProposalType, vp.Topic, vp.Desc, vp.GithubID, vp.Url, vp.EndVotingBlock, curBlockNum, state); err != nil {
		return err
	}

	if vp.NewVersion>>8 <= uint32(GovDBInstance().GetActiveVersion(state))>>8 {
		return common.NewBizError("New version should larger than current version.")
	}

	difference := vp.ActiveBlock - vp.EndVotingBlock
	quotient := difference / xcom.ConsensusSize
	remainder := difference % xcom.ConsensusSize

	if difference <= 0 || remainder != 0 || quotient < 4 || quotient > 10 {
		return common.NewBizError("active block number invalid.")
	}

	return nil
}

func (vp VersionProposal) String() string {
	return fmt.Sprintf(`Proposal %x: 
  GithubID:            	%s
  Topic:              	%s
  Type:               	%x
  Proposer:            	%x
  SubmitBlock:        	%d
  EndVotingBlock:   	%d
  ActiveBlock:   		%d
  NewVersion:   		%d`,
		vp.ProposalID, vp.GithubID, vp.Topic, vp.ProposalType, vp.Proposer, vp.SubmitBlock, vp.EndVotingBlock, vp.ActiveBlock, vp.NewVersion)
}
