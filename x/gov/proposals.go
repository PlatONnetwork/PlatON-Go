package gov

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
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

	Verify(curBlockNum *big.Int, state xcom.StateDB) (bool, error)

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

func (tp TextProposal) Verify(curBlockNum *big.Int, state xcom.StateDB) (bool, error) {

	return true, nil
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

func (vp VersionProposal) Verify(curBlockNum *big.Int, state xcom.StateDB) (bool, error) {
	/*if len(vp.ProposalID) == 0 || nil != gov.govDB.GetProposal(vp.ProposalID, state) {
		var err error = errors.New("[GOV] Verify(): ProposalID is empty or ProposalID already used.")
		return false, err
	}*/
	if len(vp.Proposer) == 0 {
		var err error = errors.New("[GOV] Verify(): Proposer is empty.")
		return false, err
	}
	if vp.ProposalType != 0x02 {
		var err error = errors.New("[GOV] Verify(): Proposal Type error.")
		return false, err
	}
	if len(vp.Topic) == 0 || len(vp.Topic) > 128 {
		var err error = errors.New("[GOV] Verify(): Topic is empty or larger than 128.")
		return false, err
	}
	if len(vp.Desc) > 512 {
		var err error = errors.New("[GOV] Verify(): Description too long.")
		return false, err
	}
	/*if len(vp.GithubID) == 0 || vp.GithubID == gov.govDB.GetProposal(vp.ProposalID, state).GetGithubID() {
		var err error = errors.New("[GOV] Verify(): GithubID empty or duplicated.")
		return false, err
	}
	if len(vp.Url) == 0 || vp.GithubID == gov.govDB.GetProposal(vp.ProposalID, state).GetUrl() {
		var err error = errors.New("[GOV] Verify(): Github URL empty or duplicated.")
		return false, err
	}*/
	//TODO
	/*if vp.EndVotingBlock == big.NewInt(0) || vp.EndVotingBlock.Cmp(curBlockNum.Add(curBlockNum, twoWeek)) > 0 {
		var err error = errors.New("[GOV] Verify(): Github URL empty or duplicated.")
		return false, err
	}*/

	/*if vp.NewVersion>>8 <= uint(gov.govDB.GetActiveVersion(state))>>8 {
		var err error = errors.New("[GOV] Verify(): NewVersion should larger than current version.")
		return false, err
	}*/
	//TODO
	/*if vp.ActiveBlock == big.NewInt(0) || vp.ActiveBlock.Cmp(fourRoundConsensus) <= 4 || vp.ActiveBlock.Cmp(fourRoundConsensus) >= 10 {
		var err error = errors.New("[GOV] Verify(): invalid ActiveBlock.")
		return false, err
	}*/

	return true, nil
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
