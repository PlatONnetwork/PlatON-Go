package vm

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	commonvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var (
	govPlugin   *plugin.GovPlugin
	gc          *GovContract
	versionSign common.VersionSign
	chandler    *xcom.CryptoHandler
)

func init() {
	chandler = xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])
	versionSign.SetBytes(chandler.MustSign(promoteVersion))
}

func buildSubmitTextInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2000))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[1])) // param 1 ...
	input = append(input, common.MustRlpEncode("textUrl"))

	return common.MustRlpEncode(input)
}

func buildSubmitVersionInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2001))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode("versionUrl"))
	input = append(input, common.MustRlpEncode(promoteVersion)) //new version : 1.1.1
	input = append(input, common.MustRlpEncode(uint64(5)))

	return common.MustRlpEncode(input)
}

func buildSubmitCancelInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2005))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ..
	input = append(input, common.MustRlpEncode("cancelPIPID"))
	input = append(input, common.MustRlpEncode(uint64(5)))
	input = append(input, common.MustRlpEncode(txHashArr[0]))
	return common.MustRlpEncode(input)
}

func buildVoteInput(nodeIdx, txIdx int) []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2003)))       // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[nodeIdx])) // param 1 ...
	input = append(input, common.MustRlpEncode(txHashArr[txIdx]))
	input = append(input, common.MustRlpEncode(uint8(1)))
	input = append(input, common.MustRlpEncode(promoteVersion))
	input = append(input, common.MustRlpEncode(versionSign))

	return common.MustRlpEncode(input)
}

func buildDeclareInput() []byte {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2004))) // func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0])) // param 1 ...
	input = append(input, common.MustRlpEncode(promoteVersion))
	input = append(input, common.MustRlpEncode(versionSign))
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

var successExpected = hexutil.Encode(common.MustRlpEncode(xcom.Result{true, "", ""}))

func setup(t *testing.T) {
	t.Log("setup()......")

	precompiledContract := PlatONPrecompiledContracts[commonvm.GovContractAddr]
	gc, _ = precompiledContract.(*GovContract)
	gc.Contract = newContract(common.Big0, sender)

	state, genesis, _ := newChainState()
	gc.Evm = newEvm(blockNumber, blockHash, state)

	newPlugins()
	govPlugin = plugin.GovPluginInstance()
	gc.Plugin = govPlugin

	build_staking_data(genesis.Hash())
}

func clear(t *testing.T) {
	t.Log("tear down()......")
	sndb.Clear()
}

func TestGovContract_SubmitText(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 2)
	sndb.NewBlock(blockNumber2, blockHash, blockHash2)

	runGovContract(gc, buildSubmitTextInput(), t)
}

func TestGovContract_GetTextProposal(t *testing.T) {
	setup(t)
	defer clear(t)
	//state := gc.Evm.StateDB.(*mock.MockStateDB)
	//state.Prepare(txHashArr[0], blockHash2, 0)

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitTextInput(), t)
	//state.Prepare(txHashArr[1], blockHash2, 0)
	runGovContract(gc, buildGetProposalInput(0), t)
}

func TestGovContract_SubmitVersion(t *testing.T) {
	setup(t)
	//defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash2, 2)
	sndb.NewBlock(blockNumber2, blockHash, blockHash2)

	runGovContract(gc, buildSubmitVersionInput(), t)
}

func TestGovContract_GetVersionProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], blockHash, 1)
	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildGetProposalInput(0), t)
}

func TestGovContract_DeclareVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], blockHash, 1)

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	chandler := xcom.GetCryptoHandler()
	chandler.SetPrivateKey(priKeyArr[0])
	runGovContract(gc, buildDeclareInput(), t)

	if nodeList, err := gov.GetActiveNodeList(blockHash, txHashArr[0]); err != nil {
		t.Error("cannot list ActiveNode")
	} else if len(nodeList) == 1 {
		t.Log("nodeList", nodeList[0])
	} else {
		t.Error("cannot list ActiveNode")
	}

}

func TestGovContract_SubmitCancel(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], blockHash, 1)
	//runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildSubmitCancelInput(), t)
}

func TestGovContract_GetCancelProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[2], blockHash, 2)

	//submit a proposal and get it.
	runGovContract(gc, buildSubmitCancelInput(), t)

	runGovContract(gc, buildGetProposalInput(2), t)
}

func TestGovContract_OneNodeVoteVersionProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[1], blockHash, 1)

	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildVoteInput(0, 1), t)
}

func TestGovContract_AllNodeVoteVersionProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[1], blockHash, 1)
	//submit a proposal and vote for it.
	runGovContract(gc, buildSubmitVersionInput(), t)
	for i := 0; i < 3; i++ {
		runGovContract(gc, buildVoteInput(i, 1), t)
	}
}

func TestGovContract_ListProposal(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], blockHash, 0)
	//submit a proposal
	runGovContract(gc, buildSubmitTextInput(), t)

	state.Prepare(txHashArr[1], blockHash, 1)
	runGovContract(gc, buildSubmitVersionInput(), t)

	runGovContract(gc, buildListProposalInput(), t)

}

func TestGovContract_GetActiveVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], blockHash, 0)
	runGovContract(gc, buildGetActiveVersionInput(), t)
}

func TestGovContract_GetProgramVersion(t *testing.T) {
	setup(t)
	defer clear(t)

	state := gc.Evm.StateDB.(*mock.MockStateDB)
	state.Prepare(txHashArr[0], blockHash, 0)
	runGovContract(gc, buildGetProgramVersionInput(), t)
}

func runGovContract(contract *GovContract, buf []byte, t *testing.T) {
	res, err := contract.Run(buf)
	if nil != err {
		t.Fatal(err)
	} else {
		t.Log(string(res))
	}
}

func Test_ResetVoteOption(t *testing.T) {
	v := gov.VoteInfo{}
	v.ProposalID = common.ZeroHash
	v.VoteNodeID = discover.NodeID{}
	v.VoteOption = gov.Abstention
	t.Log(v)

	v.VoteOption = gov.Yes
	t.Log(v)
}
