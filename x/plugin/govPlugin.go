package plugin

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
	"sync"
)

type GovPlugin struct {
	govDB   *gov.GovDB
	once sync.Once
}

var govPlugin *GovPlugin



func (govPlugin *GovPlugin) Confirmed(block *types.Block) error {
	return nil
}


//实现BasePlugin
func (govPlugin *GovPlugin) BeginBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error) {

	//TODO: Staking Plugin
	//是否当前结算的结束
	//if stk.IsEndofSettleCycle(state) {
	if true {
		verifierList, err := stk.ListVerifierNodeID(blockHash, header.Number.Uint64())

		if err != nil {
			err := errors.New("[GOV] BeginBlock(): ListVerifierNodeID failed.")
			return false, err
		}

		votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)
		for _, votingProposalID := range votingProposalIDs {
			ok := govPlugin.govDB.AccuVerifiers(blockHash, votingProposalID, verifierList)
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

func (govPlugin *GovPlugin) EndBlock(blockHash common.Hash, header *types.Header, state xcom.StateDB) (bool, error) {

	votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)
	for _, votingProposalID := range votingProposalIDs {
		votingProposal, err := govPlugin.govDB.GetProposal(votingProposalID, state)
		if nil != err {
			msg := fmt.Sprintf("[GOV] EndBlock(): Unable to get proposal: %s", votingProposalID)
			err = errors.New(msg)
			return false, err
		}
		if votingProposal.GetEndVotingBlock().Uint64() == header.Number.Uint64() {
			//if !stk.isEndofSettleCycle(blockHash) {
			if true {

				verifierList, err := stk.ListVerifierNodeID(blockHash, header.Number.Uint64())
				if err != nil {
					err := errors.New("[GOV] BeginBlock(): ListVerifierNodeID failed.")
					return false, err
				}

				ok := govPlugin.govDB.AccuVerifiers(blockHash, votingProposalID, verifierList)
				if !ok {
					err = errors.New("[GOV] EndBlock(): add Verifiers failed.")
					return false, err
				}
			}
			ok, status := govPlugin.tally(votingProposalID, blockHash, state)
			if !ok {
				err = errors.New("[GOV] EndBlock(): tally failed.")
				return false, err
			}
			if status == gov.Pass {
				_, ok := votingProposal.(gov.VersionProposal)
				if ok {
					govPlugin.govDB.MoveVotingProposalIDToPreActive(blockHash, votingProposalID, state)
				}
				_, ok = votingProposal.(gov.TextProposal)
				if ok {
					govPlugin.govDB.MoveVotingProposalIDToEnd(blockHash, votingProposalID, state)
				}
			}
		}
	}
	preActiveProposalID := govPlugin.govDB.GetPreActiveProposalID(blockHash, state)
	proposal, err := govPlugin.govDB.GetProposal(preActiveProposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] EndBlock(): Unable to get proposal: %s", preActiveProposalID)
		err = errors.New(msg)
		return false, err
	}
	versionProposal, ok := proposal.(gov.VersionProposal)
	if !ok {
		return true, nil
	}
	if versionProposal.GetActiveBlock().Uint64() == header.Number.Uint64() {
		verifierList, err := stk.ListVerifierNodeID(blockHash, header.Number.Uint64())
		if err != nil {
			err := errors.New("[GOV] BeginBlock(): ListVerifierNodeID failed.")
			return false, err
		}

		var updatedNodes uint8 = 0
		declareList := govPlugin.govDB.GetActiveNodeList(blockHash, preActiveProposalID)
		for _, val := range verifierList {
			if inNodeList(val, declareList) {
				updatedNodes++
			}
		}
		if updatedNodes == 25 {
			tallyResult, err := govPlugin.govDB.GetTallyResult(preActiveProposalID, state)
			if err != nil {
				err = errors.New("[GOV] EndBlock(): get tallyResult failed.")
				return false, err
			}
			tallyResult.Status = gov.Active
			if !govPlugin.govDB.SetTallyResult(*tallyResult, state) {
				err = errors.New("[GOV] EndBlock(): Unable to set tallyResult.")
				return false, err
			}
			govPlugin.govDB.MovePreActiveProposalIDToEnd(blockHash, preActiveProposalID, state)

			govPlugin.govDB.ClearActiveNodes(blockHash, preActiveProposalID)

			//stk.NotifyActive(versionProposal.NewVersion)
		}
	}
	return true, nil
}




//获取预生效版本，可以返回nil
func (govPlugin *GovPlugin) GetPreActiveVersion(state xcom.StateDB) uint32 {
	return govPlugin.govDB.GetPreActiveVersion(state)
}

//获取当前生效版本，不会返回nil
func (govPlugin *GovPlugin) GetActiveVersion(state xcom.StateDB) uint32 {
	return govPlugin.govDB.GetActiveVersion(state)
}


//提交提案，只有验证人才能提交提案
func (govPlugin *GovPlugin) Submit(curBlockNum *big.Int, from common.Address, proposal gov.Proposal, blockHash common.Hash, state xcom.StateDB) (bool, error) {

	//TODO 检查交易发起人的Address和NodeID是否对应；
	/*if !stk.isAdressCorrespondingToNodeID(from, proposal.GetProposer()) {
		var err error = errors.New("[GOV] Submit(): tx sender is not the declare proposer.")
		return false, err
	}*/

	//参数校验
	success, err := proposal.Verify(curBlockNum, state)
	if !success || err != nil {
		var err error = errors.New("[GOV] Submit(): param error.")
		return false, err
	}

	//判断proposer是否为Verifier
	verifierList, err := stk.ListVerifierNodeID(blockHash, curBlockNum.Uint64())
	if err != nil {
		err := errors.New("[GOV] BeginBlock(): ListVerifierNodeID failed.")
		return false, err
	}

	if err != nil {
		var err error = errors.New("[GOV] Submit(): get verifier list failed.")
		return false, err
	}
	if !inNodeList(proposal.GetProposer(), verifierList) {
		var err error = errors.New("[GOV] Submit(): proposer is not verifier.")
		return false, err
	}

	//升级提案的额外处理
	_, ok := proposal.(gov.VersionProposal)
	if ok {
		//判断是否有VersionProposal正在投票中，有则退出
		votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)
		for _, votingProposalID := range votingProposalIDs {
			votingProposal, err := govPlugin.govDB.GetProposal(votingProposalID, state)
			if err != nil {
				msg := fmt.Sprintf("[GOV] Submit(): Unable to get proposal: %s", votingProposalID)
				err = errors.New(msg)
				return false, err
			}
			_, ok := votingProposal.(gov.VersionProposal)
			if ok {
				var err error = errors.New("[GOV] Submit(): existing a voting version proposal.")
				return false, err
			}
		}
		//判断是否有VersionProposal正在Pre-active阶段，有则退出
		if len(govPlugin.govDB.GetPreActiveProposalID(blockHash, state)) > 0 {
			var err error = errors.New("[GOV] Submit(): existing a pre-active version proposal.")
			return false, err
		}
	}
	//持久化相关
	ok, err = govPlugin.govDB.SetProposal(proposal, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Submit(): Unable to set proposal: %s", proposal.GetProposalID())
		err = errors.New(msg)
		return false, err
	}
	if !ok {
		var err error = errors.New("[GOV] Submit(): set proposal failed.")
		return false, err
	}
	ok = govPlugin.govDB.AddVotingProposalID(blockHash, proposal.GetProposalID(), state)
	if !ok {
		var err error = errors.New("[GOV] Submit(): add VotingProposalID failed.")
		return false, err
	}
	return true, nil
}

//投票，只有验证人能投票
func (govPlugin *GovPlugin) Vote(from common.Address, vote gov.Vote, blockHash common.Hash, state xcom.StateDB) (bool, error) {

	if len(vote.ProposalID) == 0 || len(vote.VoteNodeID) == 0 || vote.VoteOption == 0 {
		var err error = errors.New("[GOV] Vote(): empty parameter detected.")
		return false, err
	}
	//检查交易发起人的Address和NodeID是否对应；
	/*proposal, err := govPlugin.govDB.GetProposal(vote.ProposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Vote(): Unable to get proposal: %s", vote.ProposalID)
		err = errors.New(msg)
		return false, err
	}

	if !stk.isAdressCorrespondingToNodeID(from, proposal.GetProposer()) {
		var err error = errors.New("[GOV] Vote(): tx sender is not the declare proposer.")
		return false, err
	}*/

	//判断vote.voteNodeID是否为Verifier
	//TODO: Staking Plugin
	success, err := stk.IsCurrVerifier(blockHash, vote.VoteNodeID)

	if !success || err !=nil {
		var err error = errors.New("[GOV] Vote(): proposer is not verifier.")
		return false, err
	}

	//voteOption范围检查
	if vote.VoteOption <= gov.Yes && vote.VoteOption >= gov.Abstention {
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
	votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)
	if !isVoting(vote.ProposalID, votingProposalIDs) {
		var err error = errors.New("[GOV] Vote(): vote.proposalID is not in voting.")
		return false, err
	}

	//持久化相关
	if !govPlugin.govDB.SetVote(vote.ProposalID, vote.VoteNodeID, vote.VoteOption, state) {
		var err error = errors.New("[GOV] Vote(): Set vote failed.")
		return false, err
	}
	if !govPlugin.govDB.AddActiveNode(blockHash, vote.ProposalID, vote.VoteNodeID) {
		var err error = errors.New("[GOV] Vote(): Add activeNode failed.")
		return false, err
	}
	if !govPlugin.govDB.AddVotedVerifier(vote.ProposalID, vote.ProposalID, vote.VoteNodeID) {
		var err error = errors.New("[GOV] Vote(): Add VotedVerifier failed.")
		return false, err
	}
	//存入AddActiveNode，等预生效再通知Staking
	//if !govPlugin.govDB.addActiveNode(&proposer, state) {
	//	var err error = errors.New("[GOV] DeclareVersion(): add active node failed.")
	//	return false, err
	//}
	return true, nil
}

func getLargeVersion(version uint) uint {
	return version >> 8
}

//版本声明，验证人/候选人可以声明
func (govPlugin *GovPlugin) DeclareVersion(from common.Address, declaredNodeID *discover.NodeID, version uint, blockHash common.Hash, state xcom.StateDB) (bool, error) {

	activeVersion := uint(govPlugin.govDB.GetActiveVersion(state))
	if activeVersion <= 0 {
		var err error = errors.New("[GOV] DeclareVersion(): add active version failed.")
		return false, err
	}

	if getLargeVersion(version) == getLargeVersion(activeVersion) {
		//TODO 通知staking
	}

	votingProposalIDs := govPlugin.govDB.ListVotingProposal(blockHash, state)

	for _, votingProposalID := range votingProposalIDs {
		votingProposal, err := govPlugin.govDB.GetProposal(votingProposalID, state)
		if err != nil {
			msg := fmt.Sprintf("[GOV] Submit(): Unable to set proposal: %s", votingProposalID)
			err = errors.New(msg)
			return false, err
		}
		if nil == votingProposal {
			continue
		}
		versionProposal, ok := votingProposal.(gov.VersionProposal)
		//在版本提案的投票周期内
		if ok {
			if getLargeVersion(versionProposal.GetNewVersion()) == getLargeVersion(version) {
				proposer := versionProposal.GetProposer()
				//存入AddActiveNode，等预生效再通知Staking
				if !govPlugin.govDB.AddActiveNode(blockHash,votingProposalID, proposer) {
					var err error = errors.New("[GOV] DeclareVersion(): add active node failed.")
					return false, err
				}
			}
		}
	}

	return true, nil
}

//查询提案
func (govPlugin *GovPlugin) GetProposal(proposalID common.Hash, state xcom.StateDB) gov.Proposal {
	proposal, err := govPlugin.govDB.GetProposal(proposalID, state)
	if err != nil {
		msg := fmt.Sprintf("[GOV] Submit(): Unable to set proposal: %s", proposalID)
		err = errors.New(msg)
		return nil
	}
	if proposal != nil {
		return proposal
	}
	log.Error("[GOV] GetProposal(): Unable to get proposal.")
	return nil
}

//查询提案结果
func (govPlugin *GovPlugin) GetTallyResult(proposalID common.Hash, state xcom.StateDB) *gov.TallyResult {
	tallyResult, err := govPlugin.govDB.GetTallyResult(proposalID, state)
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
func (govPlugin *GovPlugin) ListProposal(blockHash common.Hash, state xcom.StateDB) []gov.Proposal {
	return nil
}

//投票结束时，进行投票计算

func (govPlugin *GovPlugin) tally(proposalID common.Hash, blockHash common.Hash, state xcom.StateDB) (bool, gov.ProposalStatus) {

	accuVerifiersCnt := uint16(govPlugin.govDB.AccuVerifiersLength(blockHash, proposalID))
	voteCnt := uint16(len(govPlugin.govDB.ListVote(proposalID, state)))

	status := gov.Voting
	supportRate := voteCnt / accuVerifiersCnt
	//TODO
	if supportRate > accuVerifiersCnt {
		status = gov.PreActive
	} else {
		status = gov.Failed
	}

	tallyResult := &gov.TallyResult{
		ProposalID:    proposalID,
		Yeas:          voteCnt,
		AccuVerifiers: accuVerifiersCnt,
		Status:        status,
	}

	if !govPlugin.govDB.SetTallyResult(*tallyResult, state) {
		log.Error("[GOV] tally(): Unable to set tallyResult.")
		return false, status
	}
	return true, status

}