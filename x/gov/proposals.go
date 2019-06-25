package gov

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

type ProposalType uint8

const (
	Text    ProposalType = 0x01
	Version ProposalType = 0x02
)

type ProposalStatus uint8

const (
	Voting    ProposalStatus = 0x01
	Pass      ProposalStatus = 0x02
	Failed    ProposalStatus = 0x03
	PreActive ProposalStatus = 0x04
	Active    ProposalStatus = 0x05
)

type VoteOption uint8

const (
	Yeas VoteOption = iota + 1
	Nays
	Abstentions
)

type TallyResult struct {
	proposalID    common.Hash       `json:"proposalID"`
	yeas          uint16            `json:"yeas"`
	nays          uint64            `json:"nays"`
	abstentions   uint16            `json:"abstentions"`
	accuVerifiers []discover.NodeID `json:"accuVerifiers"`
	status        ProposalStatus    `json:"status"`
}

type Vote struct {
	proposalID common.Hash     `json:"proposalID"`
	voteNodeID discover.NodeID `json:"voteNodeID"`
	voteOption VoteOption      `json:"voteOption"`
}

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

	Verify() bool

	String() string
}

type TextProposal struct {
	proposalID     common.Hash
	githubID       string
	proposalType   ProposalType
	topic          string
	desc           string
	url            string
	submitBlock    uint64
	endVotingBlock uint64
	proposer       discover.NodeID
	result         TallyResult
}

func (tp TextProposal) SetProposalID(proposalID common.Hash) {
	tp.proposalID = proposalID
}

func (tp TextProposal) GetProposalID() common.Hash {
	return tp.proposalID
}

func (tp TextProposal) SetGithubID(githubID string) {
	tp.githubID = githubID
}

func (tp TextProposal) GetGithubID() string {
	return tp.githubID
}

func (tp TextProposal) SetProposalType(proposalType ProposalType) {
	tp.proposalType = proposalType
}

func (tp TextProposal) GetProposalType() ProposalType {
	return tp.proposalType
}

func (tp TextProposal) SetTopic(topic string) {
	tp.topic = topic
}

func (tp TextProposal) GetTopic() string {
	return tp.topic
}

func (tp TextProposal) SetDesc(desc string) {
	tp.desc = desc
}

func (tp TextProposal) GetDesc() string {
	return tp.desc
}

func (tp TextProposal) SetUrl(url string) {
	tp.url = url
}

func (tp TextProposal) GetUrl() string {
	return tp.url
}

func (tp TextProposal) SetSubmitBlock(blockNumber uint64) {
	tp.submitBlock = blockNumber
}

func (tp TextProposal) GetSubmitBlock() uint64 {
	return tp.submitBlock
}

func (tp TextProposal) SetEndVotingBlock(blockNumber uint64) {
	tp.endVotingBlock = blockNumber
}

func (tp TextProposal) GetEndVotingBlock() uint64 {
	return tp.endVotingBlock
}

func (tp TextProposal) SetProposer(proposer discover.NodeID) {
	tp.proposer = proposer
}

func (tp TextProposal) GetProposer() discover.NodeID {
	return tp.proposer
}

func (tp TextProposal) SetTallyResult(result TallyResult) {
	tp.result = result
}

func (tp TextProposal) GetTallyResult() TallyResult {
	return tp.result
}

func (tp TextProposal) Verify() bool {
	return true
}

func (tp TextProposal) String() string {
	return fmt.Sprintf(`Proposal %d: 
  GithubID:            	%s
  Topic:              	%s
  Type:               	%s
  Proposer:            	%s
  SubmitBlock:        	%s
  EndVotingBlock:   	%s`, tp.proposalID, tp.githubID, tp.topic, tp.proposalType, tp.proposer,
		tp.submitBlock, tp.GetEndVotingBlock())
}

type VersionProposal struct {
	TextProposal
	newVersion  uint
	activeBlock uint64
}

func (vp VersionProposal) SetNewVersion(newVersion uint) {
	vp.newVersion = newVersion
}

func (vp VersionProposal) GetNewVersion() uint {
	return vp.newVersion
}

func (vp VersionProposal) SetActiveBlock(activeBlock uint64) {
	vp.activeBlock = activeBlock
}

func (vp VersionProposal) GetActiveBlock() uint64 {
	return vp.activeBlock
}

func (vp VersionProposal) Verify() bool {
	return true
}

func (vp VersionProposal) String() string {
	return fmt.Sprintf(`Proposal %d: 
  GithubID:            	%s
  Topic:              	%s
  Type:               	%s
  Proposer:            	%s
  SubmitBlock:        	%s
  EndVotingBlock:   	%s,
  ActiveBlock:   		%s,
  NewVersion:   		%s`,
		vp.proposalID, vp.githubID, vp.topic, vp.proposalType, vp.proposer, vp.submitBlock, vp.GetEndVotingBlock(), vp.GetActiveBlock(), vp.GetNewVersion())
}
