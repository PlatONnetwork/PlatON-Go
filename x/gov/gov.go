package gov

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"sync"
)

var govOnce sync.Once
var gov *Gov

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
	return true, nil
}
func (gov *Gov) EndBlock(blockHash common.Hash, state xcom.StateDB) (bool, error) {
	return true, nil
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
	if len(gov.govDB.GetPreActiveProsposalID(state)) > 0 {
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
func (gov *Gov) Vote(from common.Address, vote Vote, state xcom.StateDB) bool {
	return true
}

//版本声明，验证人/候选人可以声明
func (gov *Gov) DeclareVersion(from common.Address, declaredNodeID *discover.NodeID, version uint, state xcom.StateDB) (bool, error) {
	return true, nil
}

//查询提案
func (gov *Gov) GetProposal(proposalID common.Hash, state xcom.StateDB) *Proposal {
	return nil
}

//查询提案结果
func (gov *Gov) GetTallyResult(proposalID common.Hash, state xcom.StateDB) *TallyResult {
	return nil
}

//查询提案列表
func (gov *Gov) ListProposal(state xcom.StateDB) []*Proposal {
	return nil
}

//投票结束时，进行投票计算
func (gov *Gov) tally(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}
