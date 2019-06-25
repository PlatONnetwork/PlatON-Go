package gov

import (
	"bytes"
	"encoding/binary"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"sync"
)

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
	state.SetState(vm.GovContractAddr, KeyProposal(proposal.GetProposalID()), MustEncoded(proposal))
	return true
}

// 查询提案记录，获取value后，第一个byte作为type，后续的才是proposal对象的编码
func (govDB *GovDB) GetProposal(proposalID common.Hash, state xcom.StateDB) Proposal {
	data := state.GetState(vm.GovContractAddr, KeyProposal(proposalID))
	var p Proposal
	if len(data) > 0 {
		pType := data[0]
		pData := data[1:]
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
	return true
}

// 查询投票记录
func (govDB *GovDB) ListVote(proposalID common.Hash, state xcom.StateDB) []Vote {
	return nil
}

// 保存投票结果
func (govDB *GovDB) SetTallyResult(tallyResult *TallyResult, state xcom.StateDB) bool {
	return true
}

// 查询投票结果
func (govDB *GovDB) GetTallyResult(proposalID common.Hash, state xcom.StateDB) *TallyResult {
	return nil
}

// 保存生效版本记录
func (govDB *GovDB) SetPreActiveVersion(preActiveVersion uint, state xcom.StateDB) bool {
	return true
}

// 查询生效版本记录
func (govDB *GovDB) GetPreActiveVersion(state xcom.StateDB) uint {

	data, err := govDB.local.Get(state.TxHash(), KeyPreActiveVersion())
	if err != nil {
		log.Error("Get PreActiveVersion error")
		return 0
	}

	b_buf := bytes.NewBuffer(data)
	var x uint
	binary.Read(b_buf, binary.BigEndian, &x)

	return x
}

// 保存生效版本记录
func (govDB *GovDB) SetActiveVersion(activeVersion uint, state xcom.StateDB) bool {
	state.SetState(vm.GovContractAddr, KeyActiveVersion(), MustEncoded(activeVersion))
	return true
}

// 查询生效版本记录
func (govDB *GovDB) GetActiveVersion(state xcom.StateDB) uint {
	data := state.GetState(vm.GovContractAddr, KeyActiveVersion())
	var ver uint = 0
	if len(data) > 0 {

		MustDecoded(data, &ver)
	}
	return 0
}

// 查询正在投票的提案
func (govDB *GovDB) ListVotingProposalID(state xcom.StateDB) []common.Hash {
	return nil
}

// 获取投票结束的提案
func (govDB *GovDB) ListEndProsposalID(state xcom.StateDB) []common.Hash {
	return nil
}

// 查询预生效的升级提案
func (govDB *GovDB) GetPreActiveProsposalID(state xcom.StateDB) common.Hash {
	return common.Hash{}
}

// 把新增提案的ID增加到正在投票的提案队列中
func (govDB *GovDB) AddVotingProposalID(proposalID common.Hash, state xcom.StateDB) bool {
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

// 累计验证人记录
func (govDB *GovDB) AccuVerifier(proposalID common.Hash, verifierList []*discover.NodeID, state xcom.StateDB) bool {
	return true
}
