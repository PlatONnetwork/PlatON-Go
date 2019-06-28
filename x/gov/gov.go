package gov

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
	"sync"
)

var govOnce sync.Once
var gov *Gov

const MinConsensusNodes uint8 = 25

type Gov struct {
	govDB *GovDB
}

func NewGov(govDB *GovDB) *Gov {
	govOnce.Do(func() {
		gov = &Gov{govDB: govDB}
	})
	return gov
}

func GovInstance() *Gov {
	if gov == nil {
		panic("Gov not initialized correctly")
	}
	return gov
}

//获取预生效版本，可以返回nil
func (gov *Gov) GetPreActiveVersion(state xcom.StateDB) uint32 {
	return govDB.getPreActiveVersion(state)
}

//获取当前生效版本，不会返回nil
func (gov *Gov) GetActiveVersion(state xcom.StateDB) uint32 {
	return govDB.getActiveVersion(state)
}

//实现BasePlugin
func (gov *Gov) BeginBlock(blockHash common.Hash, state xcom.StateDB) (bool, error) {

	//TODO: Staking Plugin
	//是否当前结算的结束
	if plugin.StakingInstance.isEndofSettleCycle(state) {
		curVerifierList := plugin.StakingInstance(nil).GetVerifierList(state)
		votingProposalIDs := gov.govDB.getVotingProposalIdList(blockHash, state)
		for _, votingProposalID := range votingProposalIDs {
			ok := gov.govDB.addVerifiers(blockHash, votingProposalID, curVerifierList)
			if !ok {
				err := errors.New("[GOV] BeginBlock(): add Verifiers failed.")
				return false, err
			}
		}
	}
	return true, nil
}

func inNodeList(proposer discover.NodeID, vList []discover.NodeID) bool {
	for _, v := range vList {
		if proposer == v {
			return true
		}
	}
	return false
}

func (gov *Gov) EndBlock(blockHash common.Hash, state xcom.StateDB, curBlockNum *big.Int) (bool, error) {

	votingProposalIDs := gov.govDB.getVotingProposalIdList(blockHash, state)
	for _, votingProposalID := range votingProposalIDs {
		votingProposal, err := gov.govDB.getProposal(votingProposalID, state)
		if nil != err {
			msg := fmt.Sprintf("[GOV] EndBlock(): Unable to get proposal: %s", votingProposalID)
			err = errors.New(msg)
			return false, err
		}
		if votingProposal.GetEndVotingBlock() == curBlockNum {
			if !staking.isEndofSettleCycle(blockHash) {
				//TODO:
				curVerifierList := plugin.StakingInstance(nil).GetVerifierList(state)
				ok := gov.govDB.addVerifiers(blockHash, votingProposalID, curVerifierList)
				if !ok {
					err = errors.New("[GOV] EndBlock(): add Verifiers failed.")
					return false, err
				}
			}
			ok, status := gov.tally(votingProposalID, blockHash, &state)
			if !ok {
				err = errors.New("[GOV] EndBlock(): tally failed.")
				return false, err
			}
			if status == Pass {
				_, ok := votingProposal.(VersionProposal)
				if ok {
					gov.govDB.moveVotingProposalIDToPreActive(blockHash, votingProposalID, state)
				}
				_, ok = votingProposal.(TextProposal)
				if ok {
					gov.govDB.moveVotingProposalIDToEnd(blockHash, votingProposalID, state)
				}
			}
		}
	}
	preActiveProposalID := gov.govDB.getPreActiveProposalId(blockHash, state)
	proposal, err := gov.govDB.getProposal(preActiveProposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] EndBlock(): Unable to get proposal: %s", preActiveProposalID)
		err = errors.New(msg)
		return false, err
	}
	versionProposal, ok := proposal.(VersionProposal)
	if !ok {
		return true, nil
	}
	if versionProposal.GetActiveBlock() == curBlockNum {
		//TODO
		curValidatorList := plugin.StakingInstance(nil).GetValidatorList(state)
		var updatedNodes uint8 = 0
		declareList := gov.govDB.getDeclaredNodeList(preActiveProposalID, state)
		for val := range curValidatorList {
			if inNodeList(val, declareList) {
				updatedNodes++
			}
		}
		if updatedNodes == MinConsensusNodes {
			tallyResult, err := gov.govDB.getTallyResult(preActiveProposalID, state)
			if err != nil {
				err = errors.New("[GOV] EndBlock(): get tallyResult failed.")
				return false, err
			}
			tallyResult.Status = Active
			if !gov.govDB.setTallyResult(*tallyResult, state) {
				err = errors.New("[GOV] EndBlock(): Unable to set tallyResult.")
				return false, err
			}
			gov.govDB.MovePreActiveProposalIDToEnd(blockHash, preActiveProposalID, state)

			for val := range declareList {
				gov.govDB.clearActiveNodes(val, state)
			}

			staking.NotifyUpdated(versionProposal.NewVersion)
			//TODO 把ActiveNode中的节点拿出来通知staking
			staking.Notify()
		}
	}
	return true, nil
}

//提交提案，只有验证人才能提交提案
func (gov *Gov) Submit(curBlockNum *big.Int, from common.Address, proposal Proposal, blockHash common.Hash, state xcom.StateDB) (bool, error) {

	//TODO 检查交易发起人的Address和NodeID是否对应；
	if !plugin.StakingInstance(nil).isAdressCorrespondingToNodeID(from, proposal.GetProposer()) {
		var err error = errors.New("[GOV] Submit(): tx sender is not the declare proposer.")
		return false, err
	}

	//参数校验
	if !proposal.Verify(curBlockNum, state) {
		var err error = errors.New("[GOV] Submit(): param error.")
		return false, err
	}

	//判断proposer是否为Verifier
	//TODO
	verifierList, err := plugin.StakingInstance(nil).GetVerifierList(state)
	if err != nil {
		var err error = errors.New("[GOV] Submit(): get verifier list failed.")
		return false, err
	}
	if !inNodeList(proposal.GetProposer(), verifierList) {
		var err error = errors.New("[GOV] Submit(): proposer is not verifier.")
		return false, err
	}

	//升级提案的额外处理
	_, ok := proposal.(VersionProposal)
	if ok {
		//判断是否有VersionProposal正在投票中，有则退出
		votingProposalIDs := gov.govDB.getVotingProposalIdList(blockHash, state)
		for _, votingProposalID := range votingProposalIDs {
			votingProposal, err := gov.govDB.getProposal(votingProposalID, state)
			if err != nil {
				msg := fmt.Sprintf("[GOV] Submit(): Unable to get proposal: %s", votingProposalID)
				err = errors.New(msg)
				return false, err
			}
			_, ok := votingProposal.(VersionProposal)
			if ok {
				var err error = errors.New("[GOV] Submit(): existing a voting version proposal.")
				return false, err
			}
		}
		//判断是否有VersionProposal正在Pre-active阶段，有则退出
		if len(gov.govDB.getPreActiveProposalId(blockHash, state)) > 0 {
			var err error = errors.New("[GOV] Submit(): existing a pre-active version proposal.")
			return false, err
		}
	}
	//持久化相关
	ok, err = gov.govDB.setProposal(proposal, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Submit(): Unable to set proposal: %s", proposal.GetProposalID())
		err = errors.New(msg)
		return false, err
	}
	if !ok {
		var err error = errors.New("[GOV] Submit(): set proposal failed.")
		return false, err
	}
	ok = gov.govDB.addVotingProposalID(blockHash, proposal.GetProposalID(), state)
	if !ok {
		var err error = errors.New("[GOV] Submit(): add VotingProposalID failed.")
		return false, err
	}
	return true, nil
}

//投票，只有验证人能投票
func (gov *Gov) Vote(from common.Address, vote Vote, blockHash common.Hash, state *xcom.StateDB) (bool, error) {

	if len(vote.ProposalID) == 0 || len(vote.VoteNodeID) == 0 || vote.VoteOption == 0 {
		var err error = errors.New("[GOV] Vote(): empty parameter detected.")
		return false, err
	}
	//TODO: Staking Plugin
	//检查交易发起人的Address和NodeID是否对应；
	proposal, err := gov.govDB.getProposal(vote.ProposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Vote(): Unable to get proposal: %s", vote.ProposalID)
		err = errors.New(msg)
		return false, err
	}
	if !plugin.StakingInstance(nil).isAdressCorrespondingToNodeID(from, proposal.GetProposer()) {
		var err error = errors.New("[GOV] Vote(): tx sender is not the declare proposer.")
		return false, err
	}

	//判断vote.voteNodeID是否为Verifier
	proposer := proposal.GetProposer()
	//TODO: Staking Plugin
	if !inNodeList(proposer, verifierList) {
		var err error = errors.New("[GOV] Vote(): proposer is not verifier.")
		return false, err
	}

	//voteOption范围检查
	if vote.VoteOption <= Yes && vote.VoteOption >= Abstention {
		var err error = errors.New("[GOV] Vote(): VoteOption invalid.")
		return false, err
	}

	//判断vote.proposalID是否存在voting中
	isVoting := func(proposalID common.Hash, votingProposalList []common.Hash) bool {
		for _, votingProposal := range votingProposalList {
			if proposalID == votingProposal {
				return true
			}
		}
		return false
	}
	votingProposalIDs := gov.govDB.getVotingProposalIdList(blockHash, state)
	if !isVoting(vote.ProposalID, votingProposalIDs) {
		var err error = errors.New("[GOV] Vote(): vote.proposalID is not in voting.")
		return false, err
	}

	//持久化相关
	if !gov.govDB.setVote(vote.ProposalID, vote.VoteNodeID, vote.VoteOption, state) {
		var err error = errors.New("[GOV] Vote(): Set vote failed.")
		return false, err
	}
	if !gov.govDB.addDeclaredNode(&vote.VoteNodeID, state) {
		var err error = errors.New("[GOV] Vote(): Add activeNode failed.")
		return false, err
	}
	if !gov.govDB.addVerifiers(vote.ProposalID, &vote.VoteNodeID, state) {
		var err error = errors.New("[GOV] Vote(): Add VotedVerifier failed.")
		return false, err
	}
	//存入AddActiveNode，等预生效再通知Staking
	//if !gov.govDB.addActiveNode(&proposer, state) {
	//	var err error = errors.New("[GOV] DeclareVersion(): add active node failed.")
	//	return false, err
	//}
	return true, nil
}

func getLargeVersion(version uint) uint {
	return version >> 8
}

//版本声明，验证人/候选人可以声明
func (gov *Gov) DeclareVersion(from common.Address, declaredNodeID *discover.NodeID, version uint, blockHash common.Hash, state *xcom.StateDB) (bool, error) {

	activeVersion := uint(gov.govDB.getActiveVersion(state))
	if activeVersion <= 0 {
		var err error = errors.New("[GOV] DeclareVersion(): add active version failed.")
		return false, err
	}

	if getLargeVersion(version) == getLargeVersion(activeVersion) {
		//TODO 通知staking
	}

	votingProposalIDs := gov.govDB.getVotingProposalIdList(blockHash, state)

	for _, votingProposalID := range votingProposalIDs {
		votingProposal, err := gov.govDB.getProposal(votingProposalID, state)
		if err != nil {
			msg := fmt.Sprintf("[GOV] Submit(): Unable to set proposal: %s", votingProposalID)
			err = errors.New(msg)
			return false, err
		}
		if nil == votingProposal {
			continue
		}
		versionProposal, ok := votingProposal.(VersionProposal)
		//在版本提案的投票周期内
		if ok {
			if getLargeVersion(versionProposal.GetNewVersion()) == getLargeVersion(version) {
				proposer := versionProposal.GetProposer()
				//存入AddActiveNode，等预生效再通知Staking
				if !gov.govDB.addDeclaredNode(blockHash, votingProposalID, proposer) {
					var err error = errors.New("[GOV] DeclareVersion(): add active node failed.")
					return false, err
				}
			}
		}
	}

	return true, nil
}

//查询提案
func (gov *Gov) GetProposal(proposalID common.Hash, state *xcom.StateDB) *Proposal {
	proposal, err := gov.govDB.getProposal(proposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Submit(): Unable to set proposal: %s", proposalID)
		err = errors.New(msg)
		return nil
	}
	if proposal != nil {
		return &proposal
	}
	log.Error("[GOV] GetProposal(): Unable to get proposal.")
	return nil
}

//查询提案结果
func (gov *Gov) GetTallyResult(proposalID common.Hash, state *xcom.StateDB) *TallyResult {
	tallyResult, err := gov.govDB.getTallyResult(proposalID, state)
	if err != nil {
		log.Error("[GOV] GetTallyResult(): Unable to get tallyResult from db.")
	}
	if nil != tallyResult {
		return tallyResult
	}
	log.Error("[GOV] GetTallyResult(): Unable to get tallyResult.")
	return nil
}

//查询提案列表
func (gov *Gov) ListProposal(blockHash common.Hash, state xcom.StateDB) []*Proposal {
	return nil
}

//投票结束时，进行投票计算

func (gov *Gov) tally(proposalID common.Hash, blockHash common.Hash, state *xcom.StateDB) (bool, ProposalStatus) {

	accuVerifiersCnt := uint16(gov.govDB.getVerifiersLength(proposalID, state))
	voteCnt := uint16(len(gov.govDB.getProposalVoteList(proposalID, state)))

	status := Voting
	supportRate := voteCnt / accuVerifiersCnt
	//TODO
	if supportRate > n {
		status = PreActive
	} else {
		status = Failed
	}

	tallyResult := &TallyResult{
		ProposalID:    proposalID,
		Yeas:          voteCnt,
		AccuVerifiers: accuVerifiersCnt,
		Status:        status,
	}

	if !gov.govDB.setTallyResult(*tallyResult, state) {
		log.Error("[GOV] tally(): Unable to set tallyResult.")
		return false, status
	}
	return true, status

}
