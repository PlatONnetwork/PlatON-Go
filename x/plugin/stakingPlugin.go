package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/core/types"
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




func (sk *StakingPlugin) BeginBlock (header *types.Header, state StateDB) (bool, error) {

	return false, nil
}

func (sk *StakingPlugin) EndBlock(header *types.Header, state StateDB) (bool, error) {


	return false, nil
}

func (sk *StakingPlugin) Confirmed(block *types.Block) error {

	return nil
}






func (sk *StakingPlugin) GetVal(state StateDB) {
	
}