package gov

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	vm2 "github.com/simonhsj/PlatON-Go/core/vm"
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
func (gov *Gov) GetPreActiveVersion(state xcom.StateDB) uint {
	return govDB.GetPreActiveVersion(state)
}

//获取当前生效版本，不会返回nil
func (gov *Gov) GetActiveVersion(state xcom.StateDB) uint {
	return govDB.GetActiveVersion(state)
}

//实现BasePlugin
func (gov *Gov) BeginBlock(blockHash common.Hash, state xcom.StateDB) (bool, error) {
	//是否当前结算的结束
	if plugin.StakingInstance.isEndofSettleCycle(blockHash) {
		curVerifierList := plugin.StakingInstance(nil).GetVerifierList(state)
		votingProposalIDs := gov.govDB.ListVotingProposalID(state)
		if len(votingProposalIDs) > 0 {
			for _, votingProposalID := range votingProposalIDs {
				length := gov.govDB.AddAccuVerifiersLength(votingProposalID, curVerifierList, state)
				if 0 != length {
					var err error = errors.New("[GOV] BeginBlock(): add accuVerifier failed.")
					return false, err
				}
			}
		}
	}
	return true, nil
}

func (gov *Gov) EndBlock(blockHash common.Hash, state xcom.StateDB, evm *vm2.EVM) (bool, error) {

	votingProposalIDs := gov.govDB.ListVotingProposalID(state)
	if len(votingProposalIDs) > 0 {
		for _, votingProposalID := range votingProposalIDs {
			votingProposal := gov.govDB.GetProposal(votingProposalID, state)
			if votingProposal.GetEndVotingBlock() == evm.BlockNumber {
				if !staking.isEndofSettleCycle(blockHash) {
					//TODO:
					curVerifierList := plugin.StakingInstance(nil).GetVerifierList(state)
					length := gov.govDB.AddAccuVerifiersLength(votingProposalID, curVerifierList, state)
					if 0 != length {
						var err error = errors.New("[GOV] EndBlock(): add accuVerifier failed.")
						return false, err
					}
				}
				ok, status := gov.tally(votingProposalID, &state)
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
		if len(preActiveProposalID) > 0 {
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
func (gov *Gov) Submit(from common.Address, proposal Proposal, state xcom.StateDB) common.Hash {

	//参数校验
	if !proposal.Verify() {
		log.Error("[GOV] Submit(): param error.")
		return state.TxHash()
	}

	//判断proposer是否为Verifier
	//TODO
	verifierList, err := plugin.StakingInstance(nil).GetVerifierList(state)
	if err != nil {
		return state.TxHash()
	}

	if !isVerifier(proposal.GetProposer(), verifierList) {
		log.Error("[GOV] Submit(): proposer is not verifier.")
		return state.TxHash()
	}

	//1.文本提案处理
	_, ok := proposal.(TextProposal)
	if ok {
		return state.TxHash()
	}

	//2. 版本提案处理
	//判断是否有VersionProposal正在投票中，有则退出
	votingProposalIDs := gov.govDB.ListVotingProposalID(state)
	if len(votingProposalIDs) > 0 {
		for _, votingProposalID := range votingProposalIDs {
			votingProposal := gov.govDB.GetProposal(votingProposalID, state)
			_, ok := votingProposal.(VersionProposal)
			if ok {
				log.Error("[GOV] Submit(): existing a voting version proposal.")
				return state.TxHash()
			}
		}
	}
	//判断是否有VersionProposal正在Pre-active阶段，有则退出
	if len(gov.govDB.GetPreActiveProposalID(state)) > 0 {
		log.Error("[GOV] Submit(): existing a pre-active version proposal.")
		return state.TxHash()
	}

	//持久化相关
	ok = gov.govDB.SetProposal(proposal, state)
	if !ok {
		log.Error("[GOV] Submit(): set proposal failed.")
		return state.TxHash()
	}
	ok = gov.govDB.AddVotingProposalID(proposal.GetProposalID(), state)
	if !ok {
		log.Error("[GOV] Submit(): add VotingProposalID failed.")
		return state.TxHash()
	}
	return state.TxHash()
}

//投票，只有验证人能投票
func (gov *Gov) Vote(from common.Address, vote Vote, state *xcom.StateDB) bool {

	//判断vote.voteNodeID是否为Verifier
	proposer := gov.govDB.GetProposal(vote.ProposalID, state).GetProposer()
	//TODO
	if !isVerifier(proposer, verifierList) {
		return false

	}

	//voteOption范围检查
	if vote.VoteOption <= Yes && vote.VoteOption >= Abstention {
		return false
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
		return false
	}

	//持久化相关
	if !gov.govDB.SetVote(vote.ProposalID, &vote.VoteNodeID, vote.VoteOption, state) {
		return false
	}
	if !gov.govDB.AddActiveNode(&vote.VoteNodeID, state) {
		return false
	}
	if !gov.govDB.AddVotedVerifier(vote.ProposalID, &vote.VoteNodeID, state) {
		return false
	}

	return true
}

func getLargeVersion(version uint) uint {
	return version & 0
}

//版本声明，验证人/候选人可以声明
func (gov *Gov) DeclareVersion(from common.Address, declaredNodeID *discover.NodeID, version uint, state *xcom.StateDB) (bool, error) {

	activeVersion := gov.govDB.GetActiveVersion(state)
	if activeVersion <= 0 {
		var err error = errors.New("[GOV] DeclareVersion(): add active version failed.")
		return false, err
	}

	if getLargeVersion(version) == getLargeVersion(activeVersion) {
		//TODO 通知staking
	}

	votingProposalIDs := gov.govDB.ListVotingProposalID(state)
	if len(votingProposalIDs) > 0 {
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
func (gov *Gov) ListProposal(state *xcom.StateDB) []*Proposal {
	return nil
}

//投票结束时，进行投票计算
func (gov *Gov) tally(proposalID common.Hash, state *xcom.StateDB) (bool, ProposalStatus) {
	accuVerifiersData := (*state).GetState(vm.GovContractAddr, KeyAccuVerifiers(proposalID))
	accuVerifiers := []discover.NodeID{}
	MustDecoded(accuVerifiersData, accuVerifiers)
	accuVerifiersCnt := uint16(len(accuVerifiers))

	votedVerifierData := (*state).GetState(vm.GovContractAddr, KeyVotedVerifier(proposalID))
	votedVerifier := []discover.NodeID{}
	MustDecoded(votedVerifierData, votedVerifier)
	votedVerifierCnt := uint16(len(votedVerifier))

	status := Voting
	supportRate := votedVerifierCnt / accuVerifiersCnt
	//TODO change 25 to config number
	if supportRate > 25 {
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
