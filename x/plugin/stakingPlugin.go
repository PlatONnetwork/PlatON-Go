package plugin

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	gerr "github.com/go-errors/errors"
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


func (sk *StakingPlugin) BeginBlock (header *types.Header, state StateDB) (bool, error) {

	return false, nil
}

func (sk *StakingPlugin) EndBlock(header *types.Header, state StateDB) (bool, error) {


	return false, nil
}

func (sk *StakingPlugin) Confirmed(block *types.Block) error {

	return nil
}



func (sk *StakingPlugin) Verify_tx_data(source [][]byte) (err error)  {

	defer func() {
		if errInr := recover(); nil != errInr {
			err = errors.New(fmt.Sprintf("parse tx data is panic: %s, txHash: %s", gerr.Wrap(err, 2).ErrorStack()))
		}
	}()


	return nil
}


func (sk *StakingPlugin) GetVal(state StateDB) {
	
}