package gov

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
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
	Nays          uint64         `json:"nays"`
	Abstentions   uint16         `json:"abstentions"`
	AccuVerifiers uint16         `json:"accuVerifiers"`
	Status        ProposalStatus `json:"status"`
}

type Vote struct {
	ProposalID common.Hash     `json:"proposalID"`
	VoteNodeID discover.NodeID `json:"voteNodeID"`
	VoteOption VoteOption      `json:"voteOption"`
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

	SetSubmitBlock(blockNumber *big.Int)
	GetSubmitBlock() *big.Int

	SetEndVotingBlock(blockNumber *big.Int)
	GetEndVotingBlock() *big.Int

	SetProposer(proposer discover.NodeID)
	GetProposer() discover.NodeID

	SetTallyResult(tallyResult TallyResult)
	GetTallyResult() TallyResult

	Verify() bool

	String() string
}

type TextProposal struct {
	ProposalID     common.Hash
	GithubID       string
	ProposalType   ProposalType
	Topic          string
	Desc           string
	Url            string
	SubmitBlock    *big.Int
	EndVotingBlock *big.Int
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

func (tp TextProposal) SetSubmitBlock(blockNumber *big.Int) {
	tp.SubmitBlock = blockNumber
}

func (tp TextProposal) GetSubmitBlock() *big.Int {
	return tp.SubmitBlock
}

func (tp TextProposal) SetEndVotingBlock(blockNumber *big.Int) {
	tp.EndVotingBlock = blockNumber
}

func (tp TextProposal) GetEndVotingBlock() *big.Int {
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

func (tp TextProposal) Verify() bool {
	return true
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
	NewVersion  uint
	ActiveBlock *big.Int
}

func (vp VersionProposal) SetNewVersion(newVersion uint) {
	vp.NewVersion = newVersion
}

func (vp VersionProposal) GetNewVersion() uint {
	return vp.NewVersion
}

func (vp VersionProposal) SetActiveBlock(activeBlock *big.Int) {
	vp.ActiveBlock = activeBlock
}

func (vp VersionProposal) GetActiveBlock() *big.Int {
	return vp.ActiveBlock
}

func (vp VersionProposal) Verify() bool {
	return true
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
