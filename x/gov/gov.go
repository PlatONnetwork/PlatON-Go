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
func (gov *Gov) Submit(proposal Proposal, state xcom.StateDB) common.Hash {
	return state.TxHash()
}

//投票，只有验证人能投票
func (gov *Gov) Vote(vote Vote, state *xcom.StateDB) bool {
	return true
}

//版本声明，验证人/候选人可以声明
func (gov *Gov) DeclareVersion(declaredNodeID *discover.NodeID, version uint, state *xcom.StateDB) (bool, error) {
	return true, nil
}

//查询提案
func (gov *Gov) GetProposal(proposalID common.Hash, state *xcom.StateDB) *Proposal {
	return nil
}

//查询提案结果
func (gov *Gov) GetTallyResult(proposalID common.Hash, state *xcom.StateDB) *TallyResult {
	return nil
}

//查询提案列表
func (gov *Gov) ListProposal(state *xcom.StateDB) []*Proposal {
	return nil
}

//投票结束时，进行投票计算
func (gov *Gov) tally(proposalID common.Hash, state *xcom.StateDB) bool {
	return true
}
