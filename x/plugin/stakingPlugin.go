package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type StakingPlugin struct {
	skDB *StakingDB

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


func (sk *StakingPlugin) BeginBlock (header *types.Header, state *state.StateDB) (bool, error) {

	return false, nil
}

func (sk *StakingPlugin) EndBlock(header *types.Header, state *state.StateDB) (bool, error) {


	return false, nil
}

func (sk *StakingPlugin) Confirmed(block *types.Block) error {

	return nil
}
