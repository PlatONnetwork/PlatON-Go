package vm_test

import (
	"fmt"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	commonvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var (
	snapdb    snapshotdb.DB
	govPlugin *plugin.GovPlugin
	gc        *vm.GovContract
)

func buildSubmitTextInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2000))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode("githubID"))
	input = append(input, common.MustRlpEncode("textTopic"))
	input = append(input, common.MustRlpEncode("textDesc"))
	input = append(input, common.MustRlpEncode("textUrl"))
	input = append(input, common.MustRlpEncode(uint64(21480)))

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildSubmitVersionInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2001))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode("githubID"))
	input = append(input, common.MustRlpEncode("versionTopic"))
	input = append(input, common.MustRlpEncode("versionDesc"))
	input = append(input, common.MustRlpEncode("versionUrl"))
	input = append(input, common.MustRlpEncode(uint32(1<<16|1<<8|1))) //new version : 1.1.1
	input = append(input, common.MustRlpEncode(uint64(21480)))
	input = append(input, common.MustRlpEncode(uint64(22500)))

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildSubmitParamInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2002))) // func type code
	input = append(input, common.MustRlpEncode("githubID"))
	input = append(input, common.MustRlpEncode("paramTopic"))
	input = append(input, common.MustRlpEncode("paramDesc"))
	input = append(input, common.MustRlpEncode("paramUrl"))
	input = append(input, common.MustRlpEncode("param1"))
	input = append(input, common.MustRlpEncode(""))
	input = append(input, common.MustRlpEncode("newValue"))

	input = append(input, common.MustRlpEncode(uint64(21500)))

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildVoteInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2003))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode(txHashArr[0]))
	input = append(input, common.MustRlpEncode(uint8(1)))
	input = append(input, common.MustRlpEncode(uint32(2<<16|0<<8|0)))

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildDeclareInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2004)))         // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0]))         // param 1 ...
	input = append(input, common.MustRlpEncode(uint32(1<<16|1<<8|1))) //new version : 1.1.1

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildGetProposalInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2100))) // func type code
	input = append(input, common.MustRlpEncode(txHashArr[0])) // param 1 ...

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildGetTallyResultInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2101))) // func type code
	input = append(input, common.MustRlpEncode(txHashArr[0])) // param 1 ...

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildListProposalInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2102))) // func type code
	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildGetActiveVersionInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2103))) // func type code
	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildGetCodeVersionInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2104))) // func type code
	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildListParamInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2105))) // func type code
	return common.Bytes2Hex(common.MustRlpEncode(input))
}

var successExpected = hexutil.Encode(common.MustRlpEncode(xcom.Result{true, "", ""}))

// each element means a call. we can reorder these elements to test different scenarios
var govContractCombinedTests = []vm.PrecompiledTest{
	{
		Input:    buildSubmitTextInput(),
		Expected: successExpected,
		Name:     "submitText1",
	},
	{
		Input:    buildSubmitVersionInput(),
		Expected: successExpected,
		Name:     "submitVersion1",
	},
	{
		Input:    buildVoteInput(),
		Expected: successExpected,
		Name:     "vote1",
	},
	{
		Input:    buildDeclareInput(),
		Expected: successExpected,
		Name:     "declare1",
	},
	{
		Input:    buildGetProposalInput(),
		Expected: successExpected,
		Name:     "getProposal1",
	},
	/*
		{
			Input:		buildGetTallyResultInput(),
			Expected:	successExpected,
			Name:		"getTallyResult1",
		},
	*/
	{
		Input:    buildListProposalInput(),
		Expected: successExpected,
		Name:     "listProposal1",
	},
}

func setup(t *testing.T) func() {
	t.Log("setup()......")

	precompiledContract := vm.PlatONPrecompiledContracts[commonvm.GovContractAddr]
	gc, _ = precompiledContract.(*vm.GovContract)
	state, genesis, _ := newChainState()
	gc.Evm = newEvm(blockNumber, blockHash, state)
	gc.Contract = newContract(common.Big0)

	newPlugins()

	govPlugin = plugin.GovPluginInstance()
	gc.Plugin = govPlugin

	build_staking_data(genesis.Hash())

	snapdb = snapshotdb.Instance()

	return func() {
		t.Log("tear down()......")
		snapdb.Clear()
	}
}

func testPlatONPrecompiled(idx int, t *testing.T) {

	test := govContractCombinedTests[idx]

	//in := common.Hex2Bytes(test.Input)
	//gc.Contract.Gas = gc.RequiredGas(in)

	state := gc.Evm.StateDB.(*state.StateDB)

	state.Prepare(txHashArr[idx], blockHash, idx)

	t.Run(fmt.Sprintf("%s-Gas=%d", test.Name, gc.Contract.Gas), func(t *testing.T) {
		if res, err := vm.RunPlatONPrecompiledContract(gc, common.Hex2Bytes(test.Input), gc.Contract); err != nil {
			t.Error(err)
		} else if common.Bytes2Hex0x(res) != test.Expected {

			t.Log(string(res))
			/*var r xcom.Result
			if err = rlp.DecodeBytes(res, &r); err != nil {
				t.Error(err)
			} else {
				t.Log(r.Data)
			}*/
		}
	})
}

// Tests the sample inputs from the elliptic curve pairing check EIP 197.
func TestPrecompiledGovContract(t *testing.T) {
	defer setup(t)()

	for i := 0; i < len(govContractCombinedTests); i++ {
		testPlatONPrecompiled(i, t)
	}
}
