package gov

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
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
	return govDB.GetPreActiveVersion(state)
}

//获取当前生效版本，不会返回nil
func (gov *Gov) GetActiveVersion(state xcom.StateDB) uint32 {
	return govDB.GetActiveVersion(state)
}

//实现BasePlugin
func (gov *Gov) BeginBlock(blockHash common.Hash, state xcom.StateDB) (bool, error) {
	//是否当前结算的结束
	//TODO
	if plugin.StakingInstance.isEndofSettleCycle(blockHash) {
		curVerifierList := plugin.StakingInstance(nil).GetVerifierList(state)
		votingProposalIDs := gov.govDB.ListVotingProposalID(state)
		for _, votingProposalID := range votingProposalIDs {
			length := gov.govDB.AddAccuVerifiersLength(votingProposalID, curVerifierList, state)
			if 0 != length {
				var err error = errors.New("[GOV] BeginBlock(): add accuVerifier failed.")
				return false, err
			}
		}

	}
	return true, nil
}

func (gov *Gov) EndBlock(blockHash common.Hash, state xcom.StateDB, curBlockNum *big.Int) (bool, error) {

	votingProposalIDs := gov.govDB.ListVotingProposalID(state)
	for _, votingProposalID := range votingProposalIDs {
		votingProposal := gov.govDB.GetProposal(votingProposalID, state)
		if votingProposal.GetEndVotingBlock() == curBlockNum {
			if !staking.isEndofSettleCycle(blockHash) {
				//TODO:
				curVerifierList := plugin.StakingInstance(nil).GetVerifierList(state)
				length := gov.govDB.AddAccuVerifiersLength(votingProposalID, curVerifierList, state)
				if 0 != length {
					var err error = errors.New("[GOV] EndBlock(): add accuVerifier failed.")
					return false, err
				}
			}
			ok, status := gov.tally(votingProposalID, blockHash, &state)
			if !ok {
				var err error = errors.New("[GOV] EndBlock(): tally failed.")
				return false, err
			}
			if status == Pass {
				_, ok := votingProposal.(VersionProposal)
				if ok {
					gov.govDB.MoveVotingProposalIDToPreActive(votingProposalID, state)
				}
				_, ok = votingProposal.(TextProposal)
				if ok {
					gov.govDB.MoveVotingProposalIDToEnd(votingProposalID, state)
				}
			}
		}
	}
	preActiveProposalID := gov.govDB.GetPreActiveProposalID(state)
	proposal := gov.govDB.GetProposal(preActiveProposalID, state)
	versionProposal, ok := proposal.(VersionProposal)
	if ok {
		if versionProposal.GetActiveBlock() == evm.BlockNumber {
			//TODO
			curValidatorList := plugin.StakingInstance(nil).GetValidatorList(state)
			var updatedNodes uint8 = 0
			for val := range curValidatorList {
				//TODO:
				if val.version == versionProposal.NewVersion {
					updatedNodes++
				}
			}
			if updatedNodes == MinConsensusNodes {
				tallyResult := gov.govDB.GetTallyResult(preActiveProposalID, state)
				if nil != tallyResult {
					tallyResult.Status = Active
				}
				if !gov.govDB.SetTallyResult(tallyResult, state) {
					var err error = errors.New("[GOV] EndBlock(): Unable to set tallyResult.")
					return false, err
				}
				gov.govDB.MovePreActiveProposalIDToEnd(preActiveProposalID, state)
				staking.NotifyUpdated(versionProposal.NewVersion)
				//TODO 把ActiveNode中的节点拿出来通知staking
				staking.Notify()

			}
		}
	}

	return true, nil
}

func isVerifier(proposer discover.NodeID, vList []discover.NodeID) bool {
	for _, v := range vList {
		if proposer == v {
			return true
		}
	}
	return false
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
	if !isVerifier(proposal.GetProposer(), verifierList) {
		var err error = errors.New("[GOV] Submit(): proposer is not verifier.")
		return false, err
	}

	//1.文本提案处理
	_, ok := proposal.(TextProposal)
	if ok {
		return true, nil
	}

	//2. 版本提案处理
	//判断是否有VersionProposal正在投票中，有则退出
	votingProposalIDs := gov.govDB.ListVotingProposalID(state)
	for _, votingProposalID := range votingProposalIDs {
		votingProposal := gov.govDB.GetProposal(votingProposalID, state)
		_, ok := votingProposal.(VersionProposal)
		if ok {
			var err error = errors.New("[GOV] Submit(): existing a voting version proposal.")
			return false, err
		}
	}
	//判断是否有VersionProposal正在Pre-active阶段，有则退出
	if len(gov.govDB.GetPreActiveProposalID(state)) > 0 {
		var err error = errors.New("[GOV] Submit(): existing a pre-active version proposal.")
		return false, err
	}

	//持久化相关
	ok = gov.govDB.SetProposal(proposal, state)
	if !ok {
		var err error = errors.New("[GOV] Submit(): set proposal failed.")
		return false, err
	}
	ok = gov.govDB.AddVotingProposalID(proposal.GetProposalID(), state)
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
	proposal := gov.govDB.GetProposal(vote.ProposalID, state)
	if !plugin.StakingInstance(nil).isAdressCorrespondingToNodeID(from, proposal.GetProposer()) {
		var err error = errors.New("[GOV] Vote(): tx sender is not the declare proposer.")
		return false, err
	}

	//判断vote.voteNodeID是否为Verifier
	proposer := proposal.GetProposer()
	//TODO: Staking Plugin
	if !isVerifier(proposer, verifierList) {
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
	votingProposalIDs := gov.govDB.ListVotingProposalID(state)
	if !isVoting(vote.ProposalID, votingProposalIDs) {
		var err error = errors.New("[GOV] Vote(): vote.proposalID is not in voting.")
		return false, err
	}

	//持久化相关
	if !gov.govDB.SetVote(vote.ProposalID, &vote.VoteNodeID, vote.VoteOption, state) {
		var err error = errors.New("[GOV] Vote(): Set vote failed.")
		return false, err
	}
	if !gov.govDB.AddActiveNode(&vote.VoteNodeID, state) {
		var err error = errors.New("[GOV] Vote(): Add activeNode failed.")
		return false, err
	}
	if !gov.govDB.AddVotedVerifier(vote.ProposalID, &vote.VoteNodeID, state) {
		var err error = errors.New("[GOV] Vote(): Add VotedVerifier failed.")
		return false, err
	}
	//存入AddActiveNode，等预生效再通知Staking
	if !gov.govDB.AddActiveNode(&proposer, state) {
		var err error = errors.New("[GOV] DeclareVersion(): add active node failed.")
		return false, err
	}
	return true, nil
}

func getLargeVersion(version uint) uint {
	return version >> 8
}

//版本声明，验证人/候选人可以声明
func (gov *Gov) DeclareVersion(from common.Address, declaredNodeID *discover.NodeID, version uint, blockHash common.Hash, state *xcom.StateDB) (bool, error) {

	activeVersion := uint(gov.govDB.GetActiveVersion(state))
	if activeVersion <= 0 {
		var err error = errors.New("[GOV] DeclareVersion(): add active version failed.")
		return false, err
	}

	if getLargeVersion(version) == getLargeVersion(activeVersion) {
		//TODO 通知staking
	}

	votingProposalIDs := gov.govDB.ListVotingProposalID(state)

	for _, votingProposalID := range votingProposalIDs {
		votingProposal := gov.govDB.GetProposal(votingProposalID, state)
		if nil == votingProposal {
			continue
		}
		versionProposal, ok := votingProposal.(VersionProposal)
		//在版本提案的投票周期内
		if ok {
			if getLargeVersion(versionProposal.GetNewVersion()) == getLargeVersion(version) {
				proposer := versionProposal.GetProposer()
				//存入AddActiveNode，等预生效再通知Staking
				if !gov.govDB.AddActiveNode(&proposer, state) {
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
	proposal := gov.govDB.GetProposal(proposalID, state)
	if nil != proposal {
		return &proposal
	}
	log.Error("[GOV] GetProposal(): Unable to get proposal.")
	return nil
}

//查询提案结果
func (gov *Gov) GetTallyResult(proposalID common.Hash, state *xcom.StateDB) *TallyResult {
	tallyResult := gov.govDB.GetTallyResult(proposalID, state)
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
	accuVerifiersData := (*state).GetState(vm.GovContractAddr, KeyAccuVerifiers(proposalID))
	accuVerifiers := []discover.NodeID{}
	MustDecoded(accuVerifiersData, accuVerifiers)
	accuVerifiersCnt := uint16(len(accuVerifiers))

	votedVerifierData := (*state).GetState(vm.GovContractAddr, KeyVotedVerifier(proposalID))
	votedVerifier := []discover.NodeID{}
	MustDecoded(votedVerifierData, votedVerifier)
	votedVerifierCnt := uint16(len(votedVerifier))

	status := Voting
	supportRate := uint8(votedVerifierCnt / accuVerifiersCnt)
	if supportRate > MinConsensusNodes {
		status = PreActive
	} else {
		status = Failed
	}

	tallyResult := &TallyResult{
		ProposalID:    proposalID,
		Yeas:          votedVerifierCnt,
		AccuVerifiers: accuVerifiersCnt,
		Status:        status,
	}

	if !gov.govDB.SetTallyResult(tallyResult, state) {
		log.Error("[GOV] tally(): Unable to set tallyResult.")
		return false, status
	}
	return true, status

}
