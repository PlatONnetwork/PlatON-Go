package gov

import (
	"encoding/json"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"sync"
)

var (
	ValueDelimiter = []byte(":")
)



var dbOnce sync.Once
var govDB *GovDB

type GovDB struct {
	govdbErr error
	snapdb   GovSnapshotDB
}

func GovDBInstance() *GovDB {
	dbOnce.Do(func() {
		govDB = &GovDB{snapdb: GovSnapshotDB{}}
	})
	return govDB
}

func tobytes(data interface{}) []byte {
	if bytes, err := json.Marshal(data); err != nil {
		return bytes
	} else {
		log.Error("govdb, marshal value to bytes error..")
		panic(err)
	}
}

// 保存提案记录，value编码规则:
//  value 为[]byte，其中byte[byte.len -2] 为type,byte[0:byte.len-2]为proposal
func (self *GovDB) SetProposal(proposal Proposal, state xcom.StateDB) error {

	bytes, e := json.Marshal(proposal)
	if e != nil {
		return common.NewSysError(e.Error())
	}

	value := append(bytes, byte(proposal.GetProposalType()))
	state.SetState(vm.GovContractAddr, KeyProposal(proposal.GetProposalID()), value)

	return nil
}

func (self *GovDB) setError(err error) {
	if err != nil {
		self.govdbErr = err
		panic(err)
	}
}

// 查询提案记录，获取value后，解码
func (self *GovDB) GetProposal(proposalID common.Hash, state xcom.StateDB) (Proposal, error) {
	value := state.GetState(vm.GovContractAddr, KeyProposal(proposalID))
	if len(value) == 0 {
		return nil, nil
	}
	var p Proposal
	pData := value[0 : len(value)-1]
	pType := value[len(value)-1]
	if pType == byte(Text) {
		var proposal TextProposal
		if e := json.Unmarshal(pData, &proposal); e != nil {
			return nil, common.NewSysError(e.Error())
		}
		p = proposal
	} else if pType == byte(Version) {
		var proposal VersionProposal
		//proposal = VersionProposal{TextProposal{},0,common.Big0}
		if e := json.Unmarshal(pData, &proposal); e != nil {
			return nil, common.NewSysError(e.Error())
		}
		p = proposal
	} else {
		return nil, common.NewSysError("Incorrect proposal type.")
	}

	return p, nil
}

func (self *GovDB) GetExistProposal(proposalID common.Hash, state xcom.StateDB) (Proposal, error) {
	p, err := self.GetProposal(proposalID, state)
	if err != nil {
		return nil, err
	}else if p == nil {
		log.Error("Cannot find proposal.", "proposalID", proposalID)
		return nil, common.NewSysError("Cannot find proposal.")
	}else{
		return p, nil
	}
}


// 从snapdb查询各个列表id,然后从逐条从statedb查询
func (self *GovDB) GetProposalList(blockHash common.Hash, state xcom.StateDB) ([]Proposal, error) {
	proposalIds, err := self.snapdb.getAllProposalIDList(blockHash)
	if err != nil {
		return nil, common.NewSysError(err.Error())
	}
	var proposls []Proposal
	for _, proposalId := range proposalIds {
		proposal, err := self.GetExistProposal(proposalId, state)
		if err != nil {
			return nil, err
		}else {
			proposls = append(proposls, proposal)
		}
	}
	return proposls, nil
}

//保存投票记录
func (self *GovDB) SetVote(proposalID common.Hash, voter discover.NodeID, option VoteOption, state xcom.StateDB) (error) {
	voteValueList, err := self.ListVoteValue(proposalID, state)
	if err != nil {
		return common.NewSysError(err.Error())
	}
	voteValueList = append(voteValueList, VoteValue{voter, option})

	voteListBytes, _ := json.Marshal(voteValueList)

	state.SetState(vm.GovContractAddr, KeyVote(proposalID), voteListBytes)
	return nil
}

// 查询投票记录
func (self *GovDB) ListVoteValue(proposalID common.Hash, state xcom.StateDB) ([]VoteValue, error) {
	voteListBytes := state.GetState(vm.GovContractAddr, KeyVote(proposalID))
	if len(voteListBytes) == 0 {
		return nil, nil
	}
	var voteList []VoteValue
	if err := json.Unmarshal(voteListBytes, &voteList); err != nil {
		return nil, common.NewSysError(err.Error()) //errors.New("Unmarshal VoteValue error")
	}
	return voteList,nil
}

func (self *GovDB) ListVotedVerifier(proposalID common.Hash, state xcom.StateDB) ([]discover.NodeID, error) {
	var voterList []discover.NodeID
	valueList, err := self.ListVoteValue(proposalID, state)
	if err != nil {
		return nil, common.NewSysError(err.Error())
	}
	for _, value := range valueList {
		voterList = append(voterList, value.VoteNodeID)
	}

	return voterList, nil
}

// 保存投票结果
func (self *GovDB) SetTallyResult(tallyResult TallyResult, state xcom.StateDB) error {
	value, err := json.Marshal(tallyResult)
	if err != nil {
		return common.NewSysError(err.Error())
	}
	state.SetState(vm.GovContractAddr, KeyTallyResult(tallyResult.ProposalID), value)
	return nil
}

// 查询投票结果
func (self *GovDB) GetTallyResult(proposalID common.Hash, state xcom.StateDB) (*TallyResult, error) {
	value := state.GetState(vm.GovContractAddr, KeyTallyResult(proposalID))

	if len(value) == 0 {
		return nil, nil
	}

	var tallyResult TallyResult
	if err := json.Unmarshal(value, &tallyResult); err != nil {
		return nil, common.NewSysError(err.Error())
	}
	return &tallyResult, nil

}

// 保存生效版本记录
func (self *GovDB) SetPreActiveVersion(preActiveVersion uint32, state xcom.StateDB) error {
	state.SetState(vm.GovContractAddr, KeyPreActiveVersion(), common.Uint32ToBytes(preActiveVersion))
	return nil
}

// 查询生效版本记录
func (self *GovDB) GetPreActiveVersion(state xcom.StateDB) uint32 {
	value := state.GetState(vm.GovContractAddr, KeyPreActiveVersion())
	return common.BytesToUint32(value)
}

// 保存生效版本记录
func (self *GovDB) SetActiveVersion(activeVersion uint32, state xcom.StateDB) error {
	state.SetState(vm.GovContractAddr, KeyActiveVersion(), common.Uint32ToBytes(activeVersion))
	return nil
}

// 查询生效版本记录
func (self *GovDB) GetActiveVersion(state xcom.StateDB) uint32 {
	value := state.GetState(vm.GovContractAddr, KeyActiveVersion())
	return common.BytesToUint32(value)
}

// 查询正在投票的提案
func (self *GovDB) ListVotingProposal(blockHash common.Hash, state xcom.StateDB) ([]common.Hash, error) {
	value, err := govDB.snapdb.getVotingIDList(blockHash)
	if err != nil {
		log.Error("List voting proposal ID error")
		return nil, common.NewSysError(err.Error())
	}
	return value, common.NewSysError(err.Error())
}

// 获取投票结束的提案
func (self *GovDB) ListEndProposalID(blockHash common.Hash, state xcom.StateDB) ([]common.Hash, error) {
	value, err := govDB.snapdb.getEndIDList(blockHash)
	if err != nil {
		return nil, common.NewSysError(err.Error())
	}

	return value, nil
}

// 查询预生效的升级提案
func (self *GovDB) GetPreActiveProposalID(blockHash common.Hash, state xcom.StateDB) (common.Hash, error) {
	value, err := govDB.snapdb.getPreActiveIDList(blockHash)
	if err != nil {
		//log.Error("Get pre-active proposal ID error")
		return common.Hash{}, common.NewSysError(err.Error())
	}
	if len(value) > 0 {
		return value[0], nil
	}else{
		return common.Hash{}, nil
	}
}

// 把新增提案的ID增加到正在投票的提案队列中
func (self *GovDB) AddVotingProposalID(blockHash common.Hash, proposalID common.Hash, state xcom.StateDB) error {
	if err := govDB.snapdb.addProposalByKey(blockHash, KeyVotingProposals(), proposalID); err != nil {
		//log.Error("add voting proposal to snapshot db error:%s", err)
		return common.NewSysError(err.Error())
	}

	return nil
}

// 把提案的ID从正在投票的提案队列中移动到预激活中
func (self *GovDB) MoveVotingProposalIDToPreActive(blockHash common.Hash, proposalID common.Hash, state xcom.StateDB) error {

	voting, err := self.snapdb.getVotingIDList(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}
	voting = remove(voting, proposalID)

	pre, err := self.snapdb.getPreActiveIDList(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	pre = append(pre, proposalID)

	//重新写入
	err = self.snapdb.addProposalByKey(blockHash, KeyVotingProposals(), proposalID)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	err = self.snapdb.addProposalByKey(blockHash, KeyPreActiveProposals(), proposalID)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	return nil
}

func remove(list []common.Hash, item common.Hash) []common.Hash {
	for i, id := range list {
		if id == item {
			list = append(list[:i], list[i+1:]...)
		}
	}
	return list
}

// 把提案的ID从正在投票的提案队列中移动到投票结束的提案队列中
func (self *GovDB) MoveVotingProposalIDToEnd(blockHash common.Hash, proposalID common.Hash, state xcom.StateDB) error {

	voting, err := self.snapdb.getVotingIDList(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	voting = remove(voting, proposalID)

	end, err := self.snapdb.getEndIDList(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	end = append(end, proposalID)

	//重新写入
	err = self.snapdb.addProposalByKey(blockHash, KeyVotingProposals(), proposalID)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	err = self.snapdb.addProposalByKey(blockHash, KeyEndProposals(), proposalID)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	return nil
}

// 把提案的ID从预激活的提案队列中移动到投票结束的提案队列中
func (self *GovDB) MovePreActiveProposalIDToEnd(blockHash common.Hash, proposalID common.Hash, state xcom.StateDB) error {

	pre, err := self.snapdb.getPreActiveIDList(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	pre = remove(pre, proposalID)

	end, err := self.snapdb.getEndIDList(blockHash)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	end = append(end, proposalID)

	//重新写入
	err = self.snapdb.addProposalByKey(blockHash, KeyPreActiveProposals(), proposalID)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	err = self.snapdb.addProposalByKey(blockHash, KeyEndProposals(), proposalID)
	if err != nil {
		return common.NewSysError(err.Error())
	}

	return nil
}

// 增加升级提案投票期间版本声明的验证人/候选人记录
func (self *GovDB) AddActiveNode(blockHash common.Hash, proposalID common.Hash, nodeID discover.NodeID) error {
	if err := self.snapdb.addActiveNode(blockHash, nodeID, proposalID); err != nil {
		log.Error("add declared node to snapshot db error,", err)
		return common.NewSysError(err.Error())
	}
	return nil
}

// 获取升级提案投票期间版本升声明的节点列表
func (self *GovDB) GetActiveNodeList(blockHash common.Hash, proposalID common.Hash) ([]discover.NodeID, error) {
	nodes, err := self.snapdb.getActiveNodeList(blockHash, proposalID)
	if err != nil {
		log.Error("get declared node list from snapshot db error,", err)
		return nil, common.NewSysError(err.Error())
	}
	return nodes, nil
}

// 升级后，清除做过版本声明的节点
func (self *GovDB) ClearActiveNodes(blockHash common.Hash, proposalID common.Hash) error {
	err := self.snapdb.deleteActiveNodeList(blockHash, proposalID)
	if err != nil {
		log.Error("delete declared node list from snapshot db error,", err)
		return common.NewSysError(err.Error())
	}
	return nil
}

// 增加已投票验证人记录
func (self *GovDB) AddVotedVerifier(blockHash common.Hash, proposalID common.Hash, voter discover.NodeID) error {
	if err := self.snapdb.addVotedVerifier(blockHash, voter, proposalID); err != nil {
		log.Error("add voted node to snapshot db error,", err)
		return common.NewSysError(err.Error())
	}
	return nil
}

// 累计在结算周期内可投票的所有验证人
func (self *GovDB) AccuVerifiers(blockHash common.Hash, proposalID common.Hash, verifierList []discover.NodeID) error {
	if err := self.snapdb.addTotalVerifiers(blockHash, proposalID, verifierList); err != nil {
		log.Error("add total verifier to snapshot db error,", err)
		return common.NewSysError(err.Error())
	}
	return nil
}

// 获取所有可投票验证人总数
func (self *GovDB) AccuVerifiersLength(blockHash common.Hash, proposalID common.Hash) (uint16, error) {
	if l, err := self.snapdb.getAccuVerifiersLength(blockHash, proposalID); err != nil {
		log.Error("add total verifier to  snapshot db error,", err)
		return 0, common.NewSysError(err.Error())
	} else {
		return l, nil
	}
}
