package vm_test

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

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

func buildSubmitTextInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2000))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode("textUrl"))
	input = append(input, common.MustRlpEncode(xutil.ConsensusSize()*5-xcom.ElectionDistance()))

	return common.MustRlpEncode(input)
}

func buildSubmitVersionInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2001))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode("versionUrl"))
	input = append(input, common.MustRlpEncode(uint32(1<<16|1<<8|1))) //new version : 1.1.1
	input = append(input, common.MustRlpEncode(xutil.ConsensusSize()*5-xcom.ElectionDistance()))
	input = append(input, common.MustRlpEncode(xutil.ConsensusSize()*10+1))

	return common.MustRlpEncode(input)
}

func buildSubmitParamInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2002))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ..
	input = append(input, common.MustRlpEncode("paramUrl"))
	input = append(input, common.MustRlpEncode("param1"))
	input = append(input, common.MustRlpEncode(""))
	input = append(input, common.MustRlpEncode("newValue"))
	input = append(input, common.MustRlpEncode(xutil.ConsensusSize()*5-xcom.ElectionDistance()))

	return common.MustRlpEncode(input)
}

func buildVoteInput(nodeIdx, txIdx int) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2003)))       // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[nodeIdx])) // param 1 ...
	input = append(input, common.MustRlpEncode(txHashArr[txIdx]))
	input = append(input, common.MustRlpEncode(uint8(1)))
	ver := uint32(2<<16 | 0<<8 | 0)
	verBytes := common.Uint32ToBytes(ver)
	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[nodeIdx])
	sign, _ := chandler.Sign(verBytes)
	input = append(input, common.MustRlpEncode(sign))

	return common.MustRlpEncode(input)
}

func buildDeclareInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2004))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...

	ver := uint32(1<<16 | 1<<8 | 1)
	verBytes := common.Uint32ToBytes(ver)
	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])
	sign, _ := chandler.Sign(verBytes)

	input = append(input, common.MustRlpEncode(ver)) //new version : 1.1.1
	input = append(input, common.MustRlpEncode(sign))
	return common.MustRlpEncode(input)
}

func buildGetProposalInput(idx int) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2100)))   // func type code
	input = append(input, common.MustRlpEncode(txHashArr[idx])) // param 1 ...

	return common.MustRlpEncode(input)
}

func buildGetTallyResultInput(idx int) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2101))) // func type code
	input = append(input, common.MustRlpEncode(txHashArr[0])) // param 1 ...

	return common.MustRlpEncode(input)
}

func buildListProposalInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2102))) // func type code
	return common.MustRlpEncode(input)
}

func buildGetActiveVersionInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2103))) // func type code
	return common.MustRlpEncode(input)
}

func buildGetProgramVersionInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2104))) // func type code
	return common.MustRlpEncode(input)
}

func buildListParamInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2105))) // func type code
	return common.MustRlpEncode(input)
}

var successExpected = hexutil.Encode(common.MustRlpEncode(xcom.Result{true, "", ""}))

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

func TestGovContract_SubmitText(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], blockHash, 0)

	runGovContract(gc, buildSubmitTextInput(), t)
}

func TestGovContract_GetTextProposal(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], blockHash, 0)

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitTextInput(), t)

	runGovContract(gc, buildGetProposalInput(0), t)
}

func TestGovContract_SubmitVersion(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[1], blockHash, 1)

	runGovContract(gc, buildSubmitTextInput(), t)
}

func TestGovContract_GetVersionProposal(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[1], blockHash, 1)

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildGetProposalInput(1), t)
}

func TestGovContract_DeclareVersion(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[1], blockHash, 1)

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildDeclareInput(), t)
}

func TestGovContract_SubmitParam(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[1], blockHash, 1)

	runGovContract(gc, buildSubmitParamInput(), t)
}

func TestGovContract_GetParamProposal(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[2], blockHash, 2)

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitParamInput(), t)

	runGovContract(gc, buildGetProposalInput(2), t)
}

func TestGovContract_OneNodeVoteVersionProposal(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[1], blockHash, 1)

	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildVoteInput(0, 1), t)
}

func TestGovContract_AllNodeVoteVersionProposal(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[1], blockHash, 1)
	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)
	for i := 0; i < 3; i++ {
		runGovContract(gc, buildVoteInput(i, 1), t)
	}
}

func TestGovContract_ListProposal(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], blockHash, 0)
	//submit a proposal
	runGovContract(gc, buildSubmitTextInput(), t)

	state.Prepare(txHashArr[1], blockHash, 1)
	runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildListProposalInput(), t)

}

func TestGovContract_GetActiveVersion(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], blockHash, 0)
	runGovContract(gc, buildGetActiveVersionInput(), t)
}

func TestGovContract_GetProgramVersion(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], blockHash, 0)
	runGovContract(gc, buildGetProgramVersionInput(), t)
}
func TestGovContract_ListParam(t *testing.T) {
	defer setup(t)()
	state := gc.Evm.StateDB.(*state.StateDB)
	state.Prepare(txHashArr[0], blockHash, 0)
	runGovContract(gc, buildListParamInput(), t)
}

func runGovContract(contract *vm.GovContract, buf []byte, t *testing.T) {
	res, err := contract.Run(buf)
	if nil != err {
		t.Fatal(err)
	} else {
		t.Log(string(res))
	}
}
