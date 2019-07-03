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
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
	"testing"
)

var snapdb = snapshotdb.Instance()
var nodeID = discover.MustHexID("bf5317e3a60a55e9c7d09cb20d4381f579c4318eb1031426612959ab5fa7a9d3f3e362b58887e83df8048115501f0b0390b4cdab4548b2728b6633ab692f9ca1")

func buildSubmitTextInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2000)))			// func type code
	input = append(input, common.MustRlpEncode(nodeID))				// param 1 ...
	input = append(input, common.MustRlpEncode("githubID"))
	input = append(input, common.MustRlpEncode("textTopic"))
	input = append(input, common.MustRlpEncode("textDesc"))
	input = append(input, common.MustRlpEncode("textUrl"))
	input = append(input, common.MustRlpEncode(uint64(1000)))

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

var successExpected = hexutil.Encode(common.MustRlpEncode(xcom.Result{true, "", ""}))

var govContractTests = []vm.PrecompiledTest{
	{
		  	Input:		buildSubmitTextInput(),
		 	Expected:	successExpected,
		   	Name:		"submitText1",
	},
}

func testPlatONPrecompiled(addr common.Address, test vm.PrecompiledTest, t *testing.T) {
	p := vm.PlatONPrecompiledContracts[addr]
	gc, _ := p.(*vm.GovContract)
	gc.Evm = newEvm()
	gc.Contract = newContract()

	govPlugin := plugin.GovPluginInstance()
	govPlugin.TestMode = true
	gc.Plugin = govPlugin

	defer snapdb.Clear()


	in := common.Hex2Bytes(test.Input)
	contract := vm.NewContract(vm.AccountRef(common.HexToAddress("0x12")),nil, new(big.Int), p.RequiredGas(in))
	t.Run(fmt.Sprintf("%s-Gas=%d", test.Name, contract.Gas), func(t *testing.T) {
		if res, err := vm.RunPlatONPrecompiledContract(p, in, contract); err != nil {
			t.Error(err)
		} else if common.Bytes2Hex(res) != test.Expected {
			t.Errorf("Expected %v, got %v", test.Expected, common.Bytes2Hex(res))
		}
	})
}


// Tests the sample inputs from the elliptic curve pairing check EIP 197.
func TestPrecompiledGovContract(t *testing.T) {
	testPlatONPrecompiled(common.HexToAddress("0x1000000000000000000000000000000000000005"), govContractTests[0], t)

	/*for _, test := range govContractTests {
		testPlatONPrecompiled(vm.GovContractAddr, test, t)
	}*/
}

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

	hash, _ := newblock(snapdb, big.NewInt(1))

	context := vm.Context{
		BlockNumber: big.NewInt(1),
		BlockHash:hash,
	}
	evm.Context = context
	return evm
}

func newContract() *vm.Contract {
	callerAddress := vm.AccountRef(common.HexToAddress("0x12"))
	contract := vm.NewContract(callerAddress, callerAddress, big.NewInt(1000), uint64(1))
	return contract
}

func newblock(snapdb snapshotdb.DB, blockNumber *big.Int) (common.Hash, error) {

	recognizedHash := generateHash("recognizedHash")

	commitHash := recognizedHash
	if err := snapdb.NewBlock(blockNumber, common.Hash{}, commitHash); err != nil {
		return common.Hash{}, err
	}

	if err := snapdb.Put(commitHash, []byte("wu"), []byte("wei")); err != nil {
		return common.Hash{}, err
	}

	get, err := snapdb.Get(commitHash, []byte("wu"))
	if err != nil {
		return common.Hash{}, err
	}
	fmt.Printf("get result :%s", get)

	return commitHash, nil
}

func generateHash(n string) common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(n))
	return rlpHash(buf.Bytes())
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
