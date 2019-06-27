package gov

import (
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"sync"
)

var (
	ValueDelimiter = []byte(":")
)

type VoteValue struct {
	voter  discover.NodeID
	option VoteOption
}

var dbOnce sync.Once
var govDB *GovDB

type GovDB struct {
	govdbErr error
	snapdb   GovSnapshotDB
}

func NewGovDB(snapdb xcom.SnapshotDB) *GovDB {
	dbOnce.Do(func() {
		govDB = &GovDB{snapdb: GovSnapshotDB{snapdb}}
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
func (self *GovDB) setProposal(proposal Proposal, state xcom.StateDB) (bool, error) {

	bytes, e := json.Marshal(proposal)
	if e != nil {
		return false, e
	}

	value := append(bytes, byte(proposal.GetProposalType()))
	state.SetState(vm.GovContractAddr, KeyProposal(proposal.GetProposalID()), value)

	return true, nil
}

func (self *GovDB) setError(err error) {
	if err != nil {
		self.govdbErr = err
		panic(err)
	}
}

// 查询提案记录，获取value后，解码
func (self *GovDB) getProposal(proposalID common.Hash, state xcom.StateDB) (Proposal, error) {
	value := state.GetState(vm.GovContractAddr, KeyProposal(proposalID))
	if len(value) == 0 {
		return nil, fmt.Errorf("no value found!")
	}
	var proposal Proposal
	pData := value[0 : len(value)-1]
	pType := value[len(value)-1]
	if pType == byte(Text) {
		proposal = &TextProposal{}
		if e := json.Unmarshal(pData, proposal); e != nil {
			return nil, e
		}
	} else if pType == byte(Version) {
		proposal = &VersionProposal{}
		if e := json.Unmarshal(pData, proposal); e != nil {
			return nil, e
		}
	} else {
		return nil, fmt.Errorf("incorrect propsal type:%b!", pType)
	}

	return proposal, nil
}

// 从snapdb查询各个列表id,然后从逐条从statedb查询
func (self *GovDB) getProposalList(blockHash common.Hash, state xcom.StateDB) ([]Proposal, error) {
	proposalIds, err := self.snapdb.getAllProposalIDList(blockHash)
	if err != nil {
		return nil, err
	}
	var proposls []Proposal
	for _, hash := range proposalIds {
		proposal, _ := self.getProposal(hash, state)
		if proposal != nil {
			proposls = append(proposls, proposal)
		}
	}
	return proposls, nil
}

//保存投票记录
func (self *GovDB) setVote(proposalID common.Hash, voter discover.NodeID, option VoteOption, state xcom.StateDB) bool {
	voteList := self.listVote(proposalID, state)
	voteList = append(voteList, VoteValue{voter, option})

	voteListBytes, _ := json.Marshal(voteList)

	state.SetState(vm.GovContractAddr, KeyVote(proposalID), voteListBytes)
	return true
}

// 查询投票记录
func (self *GovDB) listVote(proposalID common.Hash, state xcom.StateDB) []VoteValue {
	voteListBytes := state.GetState(vm.GovContractAddr, KeyVote(proposalID))

	var voteList []VoteValue
	if err := json.Unmarshal(voteListBytes, voteList); err != nil {
		return nil
	}
	return voteList
}

// 保存投票结果
func (self *GovDB) setTallyResult(tallyResult TallyResult, state xcom.StateDB) bool {
	value, _ := json.Marshal(tallyResult)
	state.SetState(vm.GovContractAddr, KeyTallyResult(tallyResult.ProposalID), value)
	return true
}

// 查询投票结果
func (self *GovDB) getTallyResult(proposalID common.Hash, state xcom.StateDB) (*TallyResult, error) {
	value := state.GetState(vm.GovContractAddr, KeyTallyResult(proposalID))

	var tallyResult TallyResult
	if err := json.Unmarshal(value, &tallyResult); err != nil {
		return nil, err
	}

	return &tallyResult, nil
}

// 保存生效版本记录
func (self *GovDB) setPreActiveVersion(preActiveVersion uint32, state xcom.StateDB) bool {
	state.SetState(vm.GovContractAddr, KeyPreActiveVersion(), byteutil.Uint32ToBytes(preActiveVersion))
	return true
}

// 查询生效版本记录
func (self *GovDB) getPreActiveVersion(state xcom.StateDB) uint32 {

	value := state.GetState(vm.GovContractAddr, KeyPreActiveVersion())
	return byteutil.BytesToUint32(value)
}

// 保存生效版本记录
func (self *GovDB) setActiveVersion(activeVersion uint, state xcom.StateDB) bool {

	state.SetState(vm.GovContractAddr, KeyActiveVersion(), tobytes(activeVersion))
	return true
}

// 查询生效版本记录
func (self *GovDB) getActiveVersion(state xcom.StateDB) uint32 {
	value := state.GetState(vm.GovContractAddr, KeyActiveVersion())
	return byteutil.BytesToUint32(value)
}

// 查询正在投票的提案
func (self *GovDB) listVotingProposal(blockHash common.Hash, state xcom.StateDB) []common.Hash {
	value, err := govDB.snapdb.getVotingIDList(blockHash)
	if err != nil {
		log.Error("List voting proposal ID error")
		return nil
	}
	return value
}

// 获取投票结束的提案
func (self *GovDB) listEndProposalID(blockHash common.Hash, state xcom.StateDB) []common.Hash {
	value, err := govDB.snapdb.getEndIDList(blockHash)
	if err != nil {
		log.Error("List end proposal ID error")
		return nil
	}

	return value
}

// 查询预生效的升级提案
func (self *GovDB) getPreActiveProposalID(blockHash common.Hash, state xcom.StateDB) []common.Hash {
	value, err := govDB.snapdb.getPreActiveIDList(blockHash)
	if err != nil {
		log.Error("Get pre-active proposal ID error")
		return nil
	}
	return value
}

// 把新增提案的ID增加到正在投票的提案队列中
func (self *GovDB) addVotingProposalID(blockHash common.Hash, proposalID common.Hash, state xcom.StateDB) bool {
	if err := govDB.snapdb.addProposalByKey(blockHash, KeyVotingProposals(), proposalID); err != nil {
		log.Error("add voting proposal to snapshot db error:%s", err)
		return false
	}

	return true
}

// 把提案的ID从正在投票的提案队列中移动到预激活中
func (self *GovDB) moveVotingProposalIDToPreActive(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}

// 把提案的ID从正在投票的提案队列中移动到投票结束的提案队列中
func (self *GovDB) moveVotingProposalIDToEnd(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}

// 把提案的ID从预激活的提案队列中移动到投票结束的提案队列中
func (self *GovDB) MovePreActiveProposalIDToEnd(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}

// 增加已投票验证人记录
func (self *GovDB) addVotedVerifier(proposalID common.Hash, voter *discover.NodeID, state xcom.StateDB) bool {
	return true
}

// 增加升级提案投票期间版本声明的验证人/候选人记录
func (self *GovDB) addActiveNode(nodeID *discover.NodeID, state xcom.StateDB) bool {
	return true
}

// 返回升级提案投票期间版本声明的验证人/候选人记录，并清除原记录
func (self *GovDB) clearActiveNodes(nodeID *discover.NodeID, state xcom.StateDB) []*discover.NodeID {
	return nil
}

// 累积验证人记录
func (self *GovDB) accuVerifiers(proposalID common.Hash, verifierList []*discover.NodeID, state xcom.StateDB) bool {
	return true
}

// 累积验证人记录
func (self *GovDB) accuVerifiersLength(proposalID common.Hash, verifierList []*discover.NodeID, state xcom.StateDB) uint16 {
	return 0
}
