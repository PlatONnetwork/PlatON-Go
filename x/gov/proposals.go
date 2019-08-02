package gov

import (
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

var (
	GovParamMap = map[string]interface{}{"param1": nil, "param2": nil, "param3": nil}
)

type ProposalType uint8

const (
	Text    ProposalType = 0x01
	Version ProposalType = 0x02
	Param   ProposalType = 0x03
)

type ProposalStatus uint8

const (
	Voting    ProposalStatus = 0x01
	Pass      ProposalStatus = 0x02
	Failed    ProposalStatus = 0x03
	PreActive ProposalStatus = 0x04
	Active    ProposalStatus = 0x05
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

type ParamValue struct {
	Name  string      `json:"Name"`
	Value interface{} `json:"Value"`
}

type Proposal interface {
	//SetProposalID(proposalID common.Hash)
	GetProposalID() common.Hash

	//SetGithubID(githubID string)
	//GetGithubID() string

	//SetTopic(topic string)
	//GetTopic() string

	//SetDesc(desc string)
	//GetDesc() string

	//SetProposalType(proposalType ProposalType)
	GetProposalType() ProposalType

	//SetUrl(url string)
	GetUrl() string

	//SetSubmitBlock(blockNumber uint64)
	GetSubmitBlock() uint64

	//SetEndVotingBlock(blockNumber uint64)
	GetEndVotingBlock() uint64

	//SetProposer(proposer discover.NodeID)
	GetProposer() discover.NodeID

	//SetTallyResult(tallyResult TallyResult)
	GetTallyResult() TallyResult

	Verify(submitBlock uint64, state xcom.StateDB) error

	String() string
}

type TextProposal struct {
	ProposalID common.Hash
	//GithubID     string
	ProposalType ProposalType
	//Topic          string
	//Desc           string
	Url            string
	SubmitBlock    uint64
	EndVotingBlock uint64
	Proposer       discover.NodeID
	Result         TallyResult
}

/*func (tp *TextProposal) SetProposalID(proposalID common.Hash) {
	tp.ProposalID = proposalID
}
*/
func (tp TextProposal) GetProposalID() common.Hash {
	return tp.ProposalID
}

/*func (tp *TextProposal) SetGithubID(githubID string) {
	tp.GithubID = githubID
}*/

/*func (tp TextProposal) GetGithubID() string {
	return tp.GithubID
}*/

/*func (tp *TextProposal) SetProposalType(proposalType ProposalType) {
	tp.ProposalType = proposalType
}*/

func (tp TextProposal) GetProposalType() ProposalType {
	return tp.ProposalType
}

/*func (tp *TextProposal) SetTopic(topic string) {
	tp.Topic = topic
}*/

/*func (tp TextProposal) GetTopic() string {
	return tp.Topic
}
*/
/*func (tp *TextProposal) SetDesc(desc string) {
	tp.Desc = desc
}
*/
/*func (tp TextProposal) GetDesc() string {
	return tp.Desc
}*/

/*func (tp *TextProposal) SetUrl(url string) {
	tp.Url = url
}*/

func (tp TextProposal) GetUrl() string {
	return tp.Url
}

/*func (tp *TextProposal) SetSubmitBlock(blockNumber uint64) {
	tp.SubmitBlock = blockNumber
}*/

func (tp TextProposal) GetSubmitBlock() uint64 {
	return tp.SubmitBlock
}

/*func (tp *TextProposal) SetEndVotingBlock(blockNumber uint64) {
	tp.EndVotingBlock = blockNumber
}*/

func (tp TextProposal) GetEndVotingBlock() uint64 {
	return tp.EndVotingBlock
}

/*func (tp *TextProposal) SetProposer(proposer discover.NodeID) {
	tp.Proposer = proposer
}*/

func (tp TextProposal) GetProposer() discover.NodeID {
	return tp.Proposer
}

/*func (tp *TextProposal) SetTallyResult(result TallyResult) {
	tp.Result = result
}*/

func (tp TextProposal) GetTallyResult() TallyResult {
	return tp.Result
}

func (tp TextProposal) Verify(submitBlock uint64, state xcom.StateDB) error {
	if tp.ProposalType != Text {
		return common.NewBizError("Proposal Type error.")
	}
	return verifyBasic(tp.ProposalID, tp.Proposer, tp.Url, tp.EndVotingBlock, submitBlock, state)
	//return verifyBasic(tp.ProposalID, tp.Proposer, tp.Topic, tp.Desc, tp.GithubID, tp.Url, tp.EndVotingBlock, submitBlock, state)
}

func (tp TextProposal) String() string {
	return fmt.Sprintf(`Proposal %x: 
  Type:               	%x
  Proposer:            	%x
  SubmitBlock:        	%d
  EndVotingBlock:   	%d`, tp.ProposalID, tp.ProposalType, tp.Proposer, tp.SubmitBlock, tp.EndVotingBlock)
}

type VersionProposal struct {
	ProposalID common.Hash
	//GithubID       string
	ProposalType ProposalType
	//Topic          string
	//Desc           string
	Url            string
	SubmitBlock    uint64
	EndVotingBlock uint64
	Proposer       discover.NodeID
	Result         TallyResult
	NewVersion     uint32
	ActiveBlock    uint64
}

/*func (vp VersionProposal) SetProposalID(proposalID common.Hash) {
	vp.ProposalID = proposalID
}
*/
func (vp VersionProposal) GetProposalID() common.Hash {
	return vp.ProposalID
}

/*func (vp VersionProposal) SetGithubID(githubID string) {
	vp.GithubID = githubID
}*/

/*func (vp VersionProposal) GetGithubID() string {
	return vp.GithubID
}*/

/*func (vp VersionProposal) SetProposalType(proposalType ProposalType) {
	vp.ProposalType = proposalType
}*/

func (vp VersionProposal) GetProposalType() ProposalType {
	return vp.ProposalType
}

/*func (tp *TextProposal) SetTopic(topic string) {
	vp.Topic = topic
}*/

/*func (vp VersionProposal) GetTopic() string {
	return vp.Topic
}*/

/*func (tp *TextProposal) SetDesc(desc string) {
	vp.Desc = desc
}
*/
/*func (vp VersionProposal) GetDesc() string {
	return vp.Desc
}
*/
/*func (vp VersionProposal) SetUrl(url string) {
	vp.Url = url
}*/

func (vp VersionProposal) GetUrl() string {
	return vp.Url
}

/*func (tp *TextProposal) SetSubmitBlock(blockNumber uint64) {
	vp.SubmitBlock = blockNumber
}*/

func (vp VersionProposal) GetSubmitBlock() uint64 {
	return vp.SubmitBlock
}

/*func (vp VersionProposal) SetEndVotingBlock(blockNumber uint64) {
	vp.EndVotingBlock = blockNumber
}*/

func (vp VersionProposal) GetEndVotingBlock() uint64 {
	return vp.EndVotingBlock
}

/*func (vp *VersionProposal) SetProposer(proposer discover.NodeID) {
	vp.Proposer = proposer
}*/

func (vp VersionProposal) GetProposer() discover.NodeID {
	return vp.Proposer
}

/*func (vp VersionProposal) SetTallyResult(result TallyResult) {
	vp.Result = result
}*/

func (vp VersionProposal) GetTallyResult() TallyResult {
	return vp.Result
}

/*func (vp *VersionProposal) SetNewVersion(newVersion uint32) {
	vp.NewVersion = newVersion
}*/

func (vp VersionProposal) GetNewVersion() uint32 {
	return vp.NewVersion
}

/*func (vp *VersionProposal) SetActiveBlock(activeBlock uint64) {
	vp.ActiveBlock = activeBlock
}*/

func (vp VersionProposal) GetActiveBlock() uint64 {
	return vp.ActiveBlock
}

func (vp VersionProposal) Verify(submitBlock uint64, state xcom.StateDB) error {

	if vp.ProposalType != Version {
		return common.NewBizError("Proposal Type error.")
	}

	if err := verifyBasic(vp.ProposalID, vp.Proposer, vp.Url, vp.EndVotingBlock, submitBlock, state); err != nil {
		return err
	}

	if vp.NewVersion>>8 <= uint32(GovDBInstance().GetActiveVersion(state))>>8 {
		return common.NewBizError("New version should larger than current version.")
	}

	if vp.ActiveBlock <= vp.EndVotingBlock {
		log.Warn("active-block should greater than end-voting-block")
		return common.NewBizError("active-block invalid.")
	} else {
		difference := vp.ActiveBlock - (vp.EndVotingBlock + 20)

		remainder := difference % xutil.ConsensusSize()
		if remainder != 0 {
			log.Warn("active-block should be multi-consensus-rounds greater than end-voting-block.")
			return common.NewBizError("active-block invalid.")
		} else {
			quotient := difference / xutil.ConsensusSize()
			if quotient <= 4 || quotient > 10 {
				log.Warn("active-block should be (4,10] consensus-rounds greater than end-voting-block.")
				return common.NewBizError("active-block invalid.")
			}
		}
	}
	return nil
}

func (vp VersionProposal) String() string {
	return fmt.Sprintf(`Proposal %x: 
  Type:               	%x
  Proposer:            	%x
  SubmitBlock:        	%d
  EndVotingBlock:   	%d
  ActiveBlock:   		%d
  NewVersion:   		%d`,
		vp.ProposalID, vp.ProposalType, vp.Proposer, vp.SubmitBlock, vp.EndVotingBlock, vp.ActiveBlock, vp.NewVersion)
}

type ParamProposal struct {
	ProposalID common.Hash
	//GithubID       string
	ProposalType ProposalType
	//Topic          string
	//Desc           string
	Url            string
	SubmitBlock    uint64
	EndVotingBlock uint64
	Proposer       discover.NodeID
	Result         TallyResult

	ParamName    string
	CurrentValue string
	NewValue     string
}

func (pp ParamProposal) GetProposalID() common.Hash {
	return pp.ProposalID
}

/*func (pp ParamProposal) GetGithubID() string {
	return pp.GithubID
}
*/
func (pp ParamProposal) GetProposalType() ProposalType {
	return pp.ProposalType
}

/*func (pp ParamProposal) GetTopic() string {
	return pp.Topic
}

func (pp ParamProposal) GetDesc() string {
	return pp.Desc
}*/

func (pp ParamProposal) GetUrl() string {
	return pp.Url
}

func (pp ParamProposal) GetSubmitBlock() uint64 {
	return pp.SubmitBlock
}

func (pp ParamProposal) GetEndVotingBlock() uint64 {
	return pp.EndVotingBlock
}

func (pp ParamProposal) GetProposer() discover.NodeID {
	return pp.Proposer
}

func (pp ParamProposal) GetTallyResult() TallyResult {
	return pp.Result
}

func (pp ParamProposal) GetParamName() string {
	return pp.ParamName
}

func (pp ParamProposal) GetCurrentValue() string {
	return pp.CurrentValue
}

func (pp ParamProposal) GetNewValue() string {
	return pp.NewValue
}

func (pp ParamProposal) Verify(submitBlock uint64, state xcom.StateDB) error {

	if pp.ProposalType != Param {
		return common.NewBizError("Proposal Type error.")
	}

	if err := verifyBasic(pp.ProposalID, pp.Proposer, pp.Url, pp.EndVotingBlock, submitBlock, state); err != nil {
		return err
	}

	if _, exist := GovParamMap[pp.ParamName]; !exist {
		return common.NewBizError("unsupported parameter.")
	}

	return nil

}

func (pp ParamProposal) String() string {
	return fmt.Sprintf(`Proposal %x: 
  Type:               	%x
  Proposer:            	%x
  SubmitBlock:        	%d
  ParamName:        	%s
  CurrentValue:        	%s
  NewValue:   			%s`,
		pp.ProposalID, pp.ProposalType, pp.Proposer, pp.SubmitBlock, pp.ParamName, pp.CurrentValue, pp.NewValue)
}

func verifyBasic(proposalID common.Hash, proposer discover.NodeID, url string, endVotingBlock uint64, submitBlock uint64, state xcom.StateDB) error {
	if len(proposalID) > 0 {
		p, err := GovDBInstance().GetProposal(proposalID, state)
		if err != nil {
			return err
		}
		if nil != p {
			return common.NewBizError("ProposalID is already used.")
		}
	} else {
		return common.NewBizError("ProposalID is empty.")
	}

	if len(proposer) == 0 {
		return common.NewBizError("Proposer is empty.")
	}

	/*if len(topic) == 0 || len(topic) > 128 {
		return common.NewBizError("Topic is empty or the size is bigger than 128.")
	}
	if len(desc) > 512 {
		return common.NewBizError("description's size is bigger than 512.")
	}*/
	/*if len(vp.GithubID) == 0 || vp.GithubID == gov.govDB.GetProposal(vp.ProposalID, state).GetGithubID() {
		var err error = errors.New("[GOV] Verify(): GithubID empty or duplicated.")
		return false, err
	}
	if len(vp.Url) == 0 || vp.GithubID == gov.govDB.GetProposal(vp.ProposalID, state).GetUrl() {
		var err error = errors.New("[GOV] Verify(): Github URL empty or duplicated.")
		return false, err
	}*/

	if (endVotingBlock+20)%xutil.ConsensusSize() != 0 {
		log.Warn("proposal's end-voting-block should be a particular block that less 20 than a certain consensus round")
		return common.NewBizError("end-voting-block invalid.")
	}

	submitRound := xutil.CalculateRound(submitBlock)
	endVotingRound := xutil.CalculateRound(endVotingBlock)

	if endVotingRound <= submitRound {
		log.Warn("end-voting-block's consensus round should be greater than submit-block's")
		return common.NewBizError("end-voting-block invalid.")
	}

	if endVotingRound > (submitRound + xutil.MaxVotingConsensusRounds()) {
		log.Warn("proposal's end-voting-block is too greater than the max consensus rounds")
		return common.NewBizError("end-voting-block invalid.")
	}

	return nil
}
