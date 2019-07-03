package vm_test

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/eth"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/node"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"math/big"
	"testing"
)



var (
	nodeIdArr = []discover.NodeID{
		discover.MustHexID("0x1f3a8672348ff6b789e416762ad53e69063138b8eb4d8780101658f24b2369f1a8e09499226b467d8bc0c4e03e1dc903df857eeb3c67733d21b6aaee28422334"),
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
		BlockNumber: big.NewInt(7),
	}
	evm.Context = context
	return evm
}

func newContract() *vm.Contract {
	callerAddress := vm.AccountRef(common.HexToAddress("0x12"))
	contract := vm.NewContract(callerAddress, callerAddress, big.NewInt(1000), uint64(1))
	return contract
}



func TestRLP_encode (t *testing.T) {

	var params [][]byte
	params = make([][]byte, 0)

	fnType, err := rlp.EncodeToBytes(uint16(1102))
	if nil != err {
		fmt.Println("fnType err", err)
	}else {
		var num uint16
		rlp.DecodeBytes(fnType, &num)
		fmt.Println("num is ", num)
	}
	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateDeposit encode rlp data fail")
	} else {
		fmt.Println("CandidateDeposit data rlp: ", hexutil.Encode(buf.Bytes()))
	}
}


func TestStakingContract_createStaking(t *testing.T) {
	stakingContract := vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(),
		Evm:	  newEvm(),
	}

	node_config := &node.Config{
		DataDir: "./",
	}

	n, err := node.New(node_config)
	if nil != err {

	}
	snapshotdb.SetDBPath(n.S)

	//typ uint16, benifitAddress common.Address, nodeId discover.NodeID,
	//	externalId, nodeName, website, details string, amount *big.Int, processVersion uint32
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	benifitAddress, _ := rlp.EncodeToBytes(addrArr[1])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[0])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes("PlatON, China")
	website, _ := rlp.EncodeToBytes("https://www.platon.network")
	details, _ := rlp.EncodeToBytes("platon super node")
	amount, _ := rlp.EncodeToBytes(big.NewInt(1213))
	processVersion, _ := rlp.EncodeToBytes(uint32(456))




	params = append(params, fnType)
	params = append(params, benifitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, processVersion)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("CandidateDeposit encode rlp data fail")
	} else {
		fmt.Println("CandidateDeposit data rlp: ", hexutil.Encode(buf.Bytes()))
	}


	res, err := stakingContract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}
}
