package gov

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"sync"
)

var (
	ValueDelimiter               = []byte(":")
)

type VoteValue struct {
	voter discover.NodeID
	option VoteOption
}

var dbOnce sync.Once
var govDB *GovDB

type GovDB struct {
	local xcom.SnapshotDB
}

func NewGovDB(local xcom.SnapshotDB) *GovDB {
	dbOnce.Do(func() {
		govDB = &GovDB{local: local}
	})
	return govDB
}

// Panics if error.
func MustEncoded(obj interface{}) []byte {
	bz, err := rlp.EncodeToBytes(obj)
	if err != nil {
		panic(err)
	}
	return bz
}

// Panics if error.
func MustDecoded(data []byte, dist interface{}) {
	err := rlp.DecodeBytes(data, dist)
	if err != nil {
		panic(err)
	}
}

// 保存提案记录，编码后，再在编码的前加上type，作为value存储
func (govDB *GovDB) SetProposal(proposal Proposal, state xcom.StateDB) bool {
	value :=	bytes.Join([][]byte{
		[]byte {byte(proposal.GetProposalType())},
		MustEncoded(proposal),
	}, []byte(""))
	state.SetState(vm.GovContractAddr, KeyProposal(proposal.GetProposalID()), value)

	return true
}

// 查询提案记录，获取value后，第一个byte作为type，后续的才是proposal对象的编码
func (govDB *GovDB) GetProposal(proposalID common.Hash, state xcom.StateDB) Proposal {
	value := state.GetState(vm.GovContractAddr, KeyProposal(proposalID))
	var p Proposal
	if len(value) > 0 {
		pType := value[0]
		pData := value[1:]
		if pType == byte(Text) {
			p = &TextProposal{}
			MustDecoded(pData, p)
		} else if pType == byte(Version) {
			p = &VersionProposal{}
			MustDecoded(pData, p)
		} else {
			log.Error("Proposal type error", "proposalID", proposalID, "proposalType", pType)
			return nil
		}
	}
	return p
}

// 保存投票记录
func (govDB *GovDB) SetVote(proposalID common.Hash, voter *discover.NodeID, option VoteOption, state xcom.StateDB) bool {
	value := state.GetState(vm.GovContractAddr, KeyVote(proposalID))
	var vvList []VoteValue
	MustDecoded(value, vvList)

	vv := VoteValue{*voter, option}

	vvList = append(vvList, vv)

	state.SetState(vm.GovContractAddr, KeyVote(proposalID), MustEncoded(vvList))
	return true
}

// 查询投票记录
func (govDB *GovDB) ListVote(proposalID common.Hash, state xcom.StateDB) []*Vote {
	value := state.GetState(vm.GovContractAddr, KeyVote(proposalID))
	var vv []VoteValue

	if len(value) > 0 {
		MustDecoded(value, vv)
	}

	voteList := []*Vote{}
	if len(vv) > 0 {
		for _, v := range  vv {
			vote := &Vote{proposalID, v.voter, v.option}
			voteList = append(voteList, vote)
		}
	}
	return voteList
}

// 保存投票结果
func (govDB *GovDB) SetTallyResult(tallyResult *TallyResult, state xcom.StateDB) bool {
	value := MustEncoded(*tallyResult)
	state.SetState(vm.GovContractAddr, KeyTallyResult(tallyResult.ProposalID), value)
	return true
}

// 查询投票结果
func (govDB *GovDB) GetTallyResult(proposalID common.Hash, state xcom.StateDB) *TallyResult {
	value := state.GetState(vm.GovContractAddr, KeyTallyResult(proposalID))

	var tallyResult TallyResult
	if len(value) > 0 {
		MustDecoded(value, tallyResult)
	}

	return &tallyResult
}

// 保存生效版本记录
func (govDB *GovDB) SetPreActiveVersion(preActiveVersion uint, state xcom.StateDB) bool {
	value := MustEncoded(preActiveVersion)
	state.SetState(vm.GovContractAddr, KeyPreActiveVersion(), value)
	return true
}

// 查询生效版本记录
func (govDB *GovDB) GetPreActiveVersion(state xcom.StateDB) uint32 {

	value := state.GetState(vm.GovContractAddr, KeyPreActiveVersion())

	var preActiveVersion uint32 = 0

	if len(value) > 0 {
		MustDecoded(value, &preActiveVersion)
	}
	return preActiveVersion

/*
	data, err := govDB.local.Get(state.TxHash(), KeyPreActiveVersion())
	if err != nil {
		log.Error("Get PreActiveVersion error")
		return 0
	}

	b_buf := bytes.NewBuffer(data)
	var x uint
	binary.Read(b_buf, binary.BigEndian, &x)

	return x
*/
}

// 保存生效版本记录
func (govDB *GovDB) SetActiveVersion(activeVersion uint, state xcom.StateDB) bool {
	state.SetState(vm.GovContractAddr, KeyActiveVersion(), MustEncoded(activeVersion))
	return true
}

// 查询生效版本记录
func (govDB *GovDB) GetActiveVersion(state xcom.StateDB) uint32 {
	value := state.GetState(vm.GovContractAddr, KeyActiveVersion())
	var version uint32 = 0
	if len(value) > 0 {
		MustDecoded(value, &version)
	}
	return version
}

// 查询正在投票的提案
func (govDB *GovDB) ListVotingProposalID(state xcom.StateDB) []common.Hash {
	value, err := govDB.local.Get(state.TxHash(), KeyVotingProposals())
	if err != nil {
		log.Error("List voting proposal ID error")
		return nil
	}

	var idList []common.Hash
	if len(value) > 0 {
		MustDecoded(value, &idList)
	}

	return idList
}

// 获取投票结束的提案
func (govDB *GovDB) ListEndProposalID(state xcom.StateDB) []common.Hash {
	value, err := govDB.local.Get(state.TxHash(), KeyEndProposals())
	if err != nil {
		log.Error("List end proposal ID error")
		return nil
	}

	var idList []common.Hash
	if len(value) > 0 {
		MustDecoded(value, &idList)
	}

	return idList
}

// 查询预生效的升级提案
func (govDB *GovDB) GetPreActiveProposalID(state xcom.StateDB) common.Hash {
	value, err := govDB.local.Get(state.TxHash(), KeyPreActiveProposalID())
	if err != nil {
		log.Error("Get pre-active proposal ID error")
		return common.Hash{}
	}

	var id common.Hash
	if len(value) > 0 {
		MustDecoded(value, &id)
	}
	return id
}

// 把新增提案的ID增加到正在投票的提案队列中
func (govDB *GovDB) AddVotingProposalID(proposalID common.Hash, state xcom.StateDB) bool {
	idList := govDB.ListVotingProposalID(state)

	idList = append(idList, proposalID)

	govDB.local.Put(state.TxHash(), KeyVotingProposals(), MustEncoded(idList))

	return true

}

// 把提案的ID从正在投票的提案队列中移动到预激活中
func (govDB *GovDB) MoveVotingProposalIDToPreActive(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}

// 把提案的ID从正在投票的提案队列中移动到投票结束的提案队列中
func (govDB *GovDB) MoveVotingProposalIDToEnd(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}

// 把提案的ID从预激活的提案队列中移动到投票结束的提案队列中
func (govDB *GovDB) MovePreActiveProposalIDToEnd(proposalID common.Hash, state xcom.StateDB) bool {
	return true
}

// 增加已投票验证人记录
func (govDB *GovDB) AddVotedVerifier(proposalID common.Hash, voter *discover.NodeID, state xcom.StateDB) bool {
	return true
}

// 增加升级提案投票期间版本声明的验证人/候选人记录
func (govDB *GovDB) AddActiveNode(nodeID *discover.NodeID, state xcom.StateDB) bool {
	return true
}

// 返回升级提案投票期间版本声明的验证人/候选人记录，并清除原记录
func (govDB *GovDB) ClearActiveNodes(nodeID *discover.NodeID, state xcom.StateDB) []*discover.NodeID {
	return nil
}

// 累积验证人记录
func (govDB *GovDB) AccuVerifiers(proposalID common.Hash, verifierList []*discover.NodeID, state xcom.StateDB) bool {
	return true
}


// 累积验证人记录
func (govDB *GovDB) AccuVerifiersLength(proposalID common.Hash, verifierList []*discover.NodeID, state xcom.StateDB) uint16 {
	return 0
}