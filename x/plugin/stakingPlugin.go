package plugin

import (
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"fmt"
	"math/big"
	"sync"
)

type StakingPlugin struct {
	skDB *StakingDB
	once  sync.Once

}



var stk *StakingPlugin



var (
	AccountVonNotEnough = errors.New("The Account von is not Enough")
)


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






func (sk *StakingPlugin) GetCandidateInfo(blockHash common.Hash,  nodeId discover.NodeID) (*xcom.Candidate, error) {

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

func (sk *StakingPlugin) CreateCandidate(state xcom.StateDB, blockHash common.Hash, typ uint16, can *xcom.Candidate) (bool, error) {

	// from account free von
	if typ == 0 {
		origin := state.GetBalance(can.StakingAddress)
		if origin.Cmp(can.ReleasedTmp) < 0 {
			log.Error("Failed to CreateCandidate on stakingPlugin: the account free von is not Enough", "originVon", origin, "stakingVon", can.ReleasedTmp)
			return false, AccountVonNotEnough
		}
		state.SubBalance(can.StakingAddress, can.ReleasedTmp)
		state.AddBalance(vm.StakingContractAddr, can.ReleasedTmp)

	}else if typ == 1 {  //  from account lockRepo von
		 // TODO call lockRepoPlugin



	}

	// build power queue
	//powerKey := tallyPowerKey(can.Shares, can.StakingBlockNum, can.StakingTxIndex)

	// TODO sk.skDB.Put(blockHash, powerKey, )

	// TODO










	return false, nil
}



func tallyPowerKey(shares *big.Int, stakeBlockNum uint64, stakeTxIndex uint32) []byte {


	priority := new(big.Int).Sub(math.MaxBig256, shares)
	prio := priority.String()
	num := fmt.Sprint(stakeBlockNum)
	index := fmt.Sprint(stakeTxIndex)
	return append(xcom.CanPowerKeyPrefix, append([]byte(prio), append([]byte(num), []byte(index)...)...)...)
}

