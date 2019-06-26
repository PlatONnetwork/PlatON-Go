package gov

import (
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
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

// 保存提案记录，value编码规则:
//  value 为[]byte，其中byte[byte.len -2] 为type,byte[0:byte.len-2]为proposal
func (self *GovDB) SetProposal(proposal Proposal, state xcom.StateDB) (bool, error) {

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
//func (self *GovDB) SetVote(proposalID common.Hash, voter *discover.NodeID, option VoteOption, state xcom.StateDB) bool {
//value := state.GetState(vm.GovContractAddr, KeyVote(proposalID))
//var vvList []VoteValue
//MustDecoded(value, vvList)

//vv := VoteValue{*voter, option}
//
//vvList = append(vvList, vv)
//
//state.SetState(vm.GovContractAddr, KeyVote(proposalID), vvList)
//return true
//}

//// 查询投票记录
//func (self *GovDB) ListVote(proposalID common.Hash, state xcom.StateDB) []*Vote {
//	value := state.GetState(vm.GovContractAddr, KeyVote(proposalID))
//	var vv []VoteValue
//
//	if len(value) > 0 {
//		MustDecoded(value, vv)
//	}
//
//	voteList := []*Vote{}
//	if len(vv) > 0 {
//		for _, v := range  vv {
//			vote := &Vote{proposalID, v.voter, v.option}
//			voteList = append(voteList, vote)
//		}
//	}
//	return voteList
//}
//
//// 保存投票结果
//func (self *GovDB) SetTallyResult(tallyResult *TallyResult, state xcom.StateDB) bool {
//	value := MustEncoded(*tallyResult)
//	state.SetState(vm.GovContractAddr, KeyTallyResult(tallyResult.ProposalID), value)
//	return true
//}
//
//// 查询投票结果
//func (self *GovDB) GetTallyResult(proposalID common.Hash, state xcom.StateDB) *TallyResult {
//	value := state.GetState(vm.GovContractAddr, KeyTallyResult(proposalID))
//
//	var tallyResult TallyResult
//	if len(value) > 0 {
//		MustDecoded(value, tallyResult)
//	}
//
//	return &tallyResult
//}
//
// 保存生效版本记录
func (self *GovDB) SetPreActiveVersion(preActiveVersion uint, state xcom.StateDB) bool {
	state.SetState(vm.GovContractAddr, KeyPreActiveVersion(), common.Int64ToBytes(int64(preActiveVersion)))
	return true
}

// 查询生效版本记录
func (self *GovDB) GetPreActiveVersion(state xcom.StateDB) uint {

	//value := state.GetState(vm.GovContractAddr, KeyPreActiveVersion())
	//
	//if len(value) > 0 {
	//	preActiveVersion := common.BytesToInt64(value)
	//}
	//return preActiveVersio
	return 0
}

///*
//	data, err := govDB.local.Get(state.TxHash(), KeyPreActiveVersion())
//	if err != nil {
//		log.Error("Get PreActiveVersion error")
//		return 0
//	}
//
//	b_buf := bytes.NewBuffer(data)
//	var x uint
//	binary.Read(b_buf, binary.BigEndian, &x)
//
//	return x
//*/
//}
//
//// 保存生效版本记录
func (self *GovDB) SetActiveVersion(activeVersion uint, state xcom.StateDB) bool {
	//state.SetState(vm.GovContractAddr, KeyActiveVersion(), MustEncoded(activeVersion))
	return true
}

// 查询生效版本记录
func (self *GovDB) GetActiveVersion(state xcom.StateDB) uint {
	//value := state.GetState(vm.GovContractAddr, KeyActiveVersion())
	//var version uint = 0
	//if len(value) > 0 {
	//	MustDecoded(value, &version)
	//}
	return 0
}

//
//// 查询正在投票的提案
//func (self *GovDB) ListVotingProposalID(state xcom.StateDB) []common.Hash {
//	value, err := govDB.local.Get(state.TxHash(), KeyVotingProposals())
//	if err != nil {
//		log.Error("List voting proposal ID error")
//		return nil
//	}
//
//	var idList []common.Hash
//	if len(value) > 0 {
//		MustDecoded(value, &idList)
//	}
//
//	return idList
//}
//
//// 获取投票结束的提案
//func (self *GovDB) ListEndProposalID(state xcom.StateDB) []common.Hash {
//	value, err := govDB.local.Get(state.TxHash(), KeyEndProposals())
//	if err != nil {
//		log.Error("List end proposal ID error")
//		return nil
//	}
//
//	var idList []common.Hash
//	if len(value) > 0 {
//		MustDecoded(value, &idList)
//	}
//
//	return idList
//}
//
//// 查询预生效的升级提案
//func (self *GovDB) GetPreActiveProposalID(state xcom.StateDB) common.Hash {
//	value, err := govDB.local.Get(state.TxHash(), KeyPreActiveProposals())
//	if err != nil {
//		log.Error("Get pre-active proposal ID error")
//		return common.Hash{}
//	}
//
//	var id common.Hash
//	if len(value) > 0 {
//		MustDecoded(value, &id)
//	}
//	return id
//}
//
//// 把新增提案的ID增加到正在投票的提案队列中
//func (self *GovDB) AddVotingProposalID(proposalID common.Hash, state xcom.StateDB) bool {
//	idList := govDB.ListVotingProposalID(state)
//
//	idList = append(idList, proposalID)
//
//	govDB.local.Put(state.TxHash(), KeyVotingProposals(), MustEncoded(idList))
//
//	return true
//
//}

// 把提案的ID从正在投票的提案队列中移动到预激活中
func (self *GovDB) MoveVotingProposalIDToPreActive(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}

// 把提案的ID从正在投票的提案队列中移动到投票结束的提案队列中
func (self *GovDB) MoveVotingProposalIDToEnd(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}

// 把提案的ID从预激活的提案队列中移动到投票结束的提案队列中
func (self *GovDB) MovePreActiveProposalIDToEnd(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}

// 增加已投票验证人记录
func (self *GovDB) AddVotedVerifier(proposalID common.Hash, voter *discover.NodeID, state xcom.StateDB) bool {
	return true
}

// 增加升级提案投票期间版本声明的验证人/候选人记录
func (self *GovDB) AddActiveNode(nodeID *discover.NodeID, state xcom.StateDB) bool {
	return true
}

// 返回升级提案投票期间版本声明的验证人/候选人记录，并清除原记录
func (self *GovDB) ClearActiveNodes(nodeID *discover.NodeID, state xcom.StateDB) []*discover.NodeID {
	return nil
}

// 累积验证人记录
func (self *GovDB) AccuVerifiers(proposalID common.Hash, verifierList []*discover.NodeID, state xcom.StateDB) bool {
	return true
}

// 累积验证人记录
func (self *GovDB) AccuVerifiersLength(proposalID common.Hash, verifierList []*discover.NodeID, state xcom.StateDB) uint16 {
	return 0
}
