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
		discover.MustHexID("0xa6ef31a2006f55f5039e23ccccef343e735d56699bde947cfe253d441f5f291561640a8e2bbaf8a85a8a367b939efcef6f80ae28d2bd3d0b21bdac01c3aa6f2f"),
		discover.MustHexID("0xc7fc34d6d8b3d894a35895aaf2f788ed445e03b7673f7ce820aa6fdc02908eeab6982b7eb97e983cc708bcec093b3bc512b0b1fbf668e6ab94cd91f2d642e591"),
		discover.MustHexID("0x97e424be5e58bfd4533303f8f515211599fd4ffe208646f7bfdf27885e50b6dd85d957587180988e76ae77b4b6563820a27b16885419e5ba6f575f19f6cb36b0"),
	}
	addrArr = []common.Address{
		common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80aadd24"),
		common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80444555"),
		common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80eeda12"),
		common.HexToAddress("0x740ce31b3fac20dac379db243021a51e8052234"),
	}


	blockNumer = big.NewInt(1)
	blockHash = common.HexToHash("9d4fb5346abcf593ad80a0d3d5a371b22c962418ad34189d5b1b39065668d663")

	sender = common.HexToAddress("0xeef233120ce31b3fac20dac379db243021a5234")

	sndb = snapshotdb.Instance()

	sender_balance, _ = new(big.Int).SetString("9999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999999", 10)






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
	state.AddBalance(sender, sender_balance)
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

func newContract(value *big.Int) *vm.Contract {
	callerAddress := vm.AccountRef(sender)
	fmt.Println("newContract sender :", callerAddress.Address().Hex())
	contract := vm.NewContract(callerAddress, callerAddress, value, uint64(1))
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


	//canArr := make(staking.CandidateQueue, 0)


	c1 := &staking.Candidate{
		NodeId: nodeId_A,
		StakingAddress: sender,
		BenifitAddress: addrArr[1],
		StakingTxIndex: uint32(2),
		ProcessVersion: uint32(1),
		Status: staking.Valided,
		StakingEpoch: uint32(1),
		StakingBlockNum: uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "xxccccdddddddd",
			NodeName: "I Am " +fmt.Sprint(1),
			Website: "www.baidu.com",
			Details: "this is  baidu ~~",
		},
	}

	c2 := &staking.Candidate{
		NodeId: nodeId_B,
		StakingAddress: sender,
		BenifitAddress: addrArr[2],
		StakingTxIndex: uint32(3),
		ProcessVersion: uint32(1),
		Status: staking.Valided,
		StakingEpoch: uint32(1),
		StakingBlockNum: uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "SFSFSFSFSFSFSSFS",
			NodeName: "I Am " +fmt.Sprint(2),
			Website: "www.JD.com",
			Details: "this is  JD ~~",
		},
	}



	c3 := &staking.Candidate{
		NodeId: nodeId_C,
		StakingAddress: sender,
		BenifitAddress: addrArr[3],
		StakingTxIndex: uint32(4),
		ProcessVersion: uint32(1),
		Status: staking.Valided,
		StakingEpoch: uint32(1),
		StakingBlockNum: uint64(1),
		Shares:             common.Big256,
		Released:           common.Big2,
		ReleasedHes:        common.Big32,
		RestrictingPlan:    common.Big1,
		RestrictingPlanHes: common.Big257,
		Description: staking.Description{
			ExternalId: "FWAGGDGDGG",
			NodeName: "I Am " +fmt.Sprint(3),
			Website: "www.alibaba.com",
			Details: "this is  alibaba ~~",
		},
	}


	//canArr = append(canArr, c1)
	//canArr = append(canArr, c2)
	//canArr = append(canArr, c3)

	stakingDB.SetCanPowerStore(blockHash, addr_A, c1)
	stakingDB.SetCanPowerStore(blockHash, addr_B, c2)
	stakingDB.SetCanPowerStore(blockHash, addr_C, c3)


	stakingDB.SetCandidateStore(blockHash, addr_A, c1)
	stakingDB.SetCandidateStore(blockHash, addr_B, c2)
	stakingDB.SetCandidateStore(blockHash, addr_C, c3)


	queue := make(staking.ValidatorQueue, 0)

	v1 := &staking.Validator{
		NodeAddress: addr_A,
		NodeId: c1.NodeId,
		StakingWeight: [4]string{"1", common.Big256.String(), fmt.Sprint(c1.StakingBlockNum), fmt.Sprint(c1.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v2 := &staking.Validator{
		NodeAddress: addr_B,
		NodeId: c2.NodeId,
		StakingWeight: [4]string{"1", common.Big256.String(), fmt.Sprint(c2.StakingBlockNum), fmt.Sprint(c2.StakingTxIndex)},
		ValidatorTerm: 0,
	}

	v3 := &staking.Validator{
		NodeAddress: addr_C,
		NodeId: c3.NodeId,
		StakingWeight: [4]string{"1", common.Big256.String(), fmt.Sprint(c3.StakingBlockNum), fmt.Sprint(c3.StakingTxIndex)},
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



