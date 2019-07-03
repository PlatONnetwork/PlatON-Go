package vm_test

import (
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
)



var (
	nodeIdArr = []discover.NodeID{
		discover.MustHexID("0x1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee2840e429"),
		discover.MustHexID("0x2f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee28435466"),
		discover.MustHexID("0x3f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee28544878"),
		discover.MustHexID("0x3f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee28564646"),
	}
	addrArr = []common.Address{
		common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80qeqqee"),
		common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80444555"),
		common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80wrwwwd"),
		common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80vvbbbb"),
	}


	blockNumer = big.NewInt(1)
	blockHash = common.HexToHash("9d4fb5346abcf593ad80a0d3d5a371b22c962418ad34189d5b1b39065668d663")


	sndb = snapshotdb.Instance()
)


func newChainState() (*state.StateDB, error) {
	var (
		db      = ethdb.NewMemDatabase()
		genesis = new(core.Genesis).MustCommit(db)
	)
	fmt.Println("genesis", genesis)
	// Initialize a fresh chain with only a genesis block
	blockchain, _ := core.NewBlockChain(db, nil, params.AllEthashProtocolChanges, nil, vm.Config{}, nil)

	var state *state.StateDB
	if statedb, err := blockchain.State(); nil != err {
		return nil, errors.New("reference statedb failed" + err.Error())
	} else {
		state = statedb
	}
	return state, nil
}


func newEvm() *vm.EVM {
	state, _ := newChainState()
	evm := &vm.EVM{
		StateDB:              state,
	}
	context := vm.Context{
		BlockNumber: blockNumer,
		BlockHash: blockHash,
	}
	evm.Context = context
	return evm
}

func newContract() *vm.Contract {
	callerAddress := vm.AccountRef(common.HexToAddress("0x120ce31b3fac20dac379db243021a52234vvbbbb"))
	contract := vm.NewContract(callerAddress, callerAddress, big.NewInt(1000), uint64(1))
	return contract
}


func build_staking_data (){


	stakingDB := staking.NewStakingDB ()

	// MOCK

	nodeId_A := nodeIdArr[0]
	addr_A, _ := xutil.NodeId2Addr(nodeId_A)

	nodeId_B := nodeIdArr[1]
	addr_B, _ := xutil.NodeId2Addr(nodeId_B)

	nodeId_C := nodeIdArr[2]
	addr_C, _ := xutil.NodeId2Addr(nodeId_C)

	queue := make(staking.ValidatorQueue, 0)

	v1 := &staking.Validator{
		NodeAddress: addr_A,
		NodeId: nodeId_A,
		StakingWeight: [4]string{"", "", "", ""},
		ValidatorTerm: 0,
	}

	v2 := &staking.Validator{
		NodeAddress: addr_B,
		NodeId: nodeId_B,
		StakingWeight: [4]string{"", "", "", ""},
		ValidatorTerm: 0,
	}

	v3 := &staking.Validator{
		NodeAddress: addr_C,
		NodeId: nodeId_C,
		StakingWeight: [4]string{"", "", "", ""},
		ValidatorTerm: 0,
	}

	queue = append(queue, v1)
	queue = append(queue, v2)
	queue = append(queue, v3)


	val_Arr :=  &staking.Validator_array{
		Start: 1,
		End: 22000,
		Arr: queue,
	}

	stakingDB.SetVerfierList(blockHash, val_Arr)

	stakingDB.SetPreValidatorList(blockHash, val_Arr)
	stakingDB.SetCurrentValidatorList(blockHash, val_Arr)
}



