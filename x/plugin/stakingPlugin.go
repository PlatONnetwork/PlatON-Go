package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"sync"
)

type StakingPlugin struct {
	skDB *StakingDB
	once  sync.Once

}



var stk *StakingPlugin



// Instance a global StakingPlugin
func  StakingInstance (db interface{}) *StakingPlugin {
	if nil == stk {
		stk = &StakingPlugin{
			skDB: NewStakingDB(db),

		}
	}
	return stk
}




func (sk *StakingPlugin) BeginBlock (header *types.Header, state xcom.StateDB) (bool, error) {

	return false, nil
}

func (sk *StakingPlugin) EndBlock(header *types.Header, state xcom.StateDB) (bool, error) {


	return false, nil
}

func (sk *StakingPlugin) Confirmed(block *types.Block) error {

	return nil
}






func (sk *StakingPlugin) GetCandidateInfo(state xcom.StateDB, blockHash common.Hash,  nodeId discover.NodeID) (*xcom.Candidate, error) {

	canByte, err := sk.skDB.Get(blockHash, xcom.CandidateKey(nodeId))
	if nil != err {
		return nil, err
	}

	var can xcom.Candidate

	if err := rlp.DecodeBytes(canByte, &can); nil != err {
		return nil, err
	}

	return &can, nil
}

func (sk *StakingPlugin) CreateCandidate(can *xcom.Candidate) (bool, error) {




	return false, nil
}



