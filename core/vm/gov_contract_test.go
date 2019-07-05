package vm_test

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	commonvm "github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"reflect"
	"testing"
)

var (
	snapdb 		snapshotdb.DB
	govPlugin	*plugin.GovPlugin
	gc			*vm.GovContract
)


func buildSubmitTextInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2000)))				// func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0]))				// param 1 ...
	input = append(input, common.MustRlpEncode("githubID"))
	input = append(input, common.MustRlpEncode("textTopic"))
	input = append(input, common.MustRlpEncode("textDesc"))
	input = append(input, common.MustRlpEncode("textUrl"))
	input = append(input, common.MustRlpEncode(uint64(1000)))

	return common.Bytes2Hex(common.MustRlpEncode(input))
}


func buildSubmitVersionInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2001)))				// func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0]))				// param 1 ...
	input = append(input, common.MustRlpEncode("githubID"))
	input = append(input, common.MustRlpEncode("versionTopic"))
	input = append(input, common.MustRlpEncode("versionDesc"))
	input = append(input, common.MustRlpEncode("versionUrl"))
	input = append(input, common.MustRlpEncode(uint32(1<<16 | 1<<8 | 1)))	//new version : 1.1.1
	input = append(input, common.MustRlpEncode(uint64(1000)))
	input = append(input, common.MustRlpEncode(uint64(2000)))

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildVoteInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2002)))				// func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0]))				// param 1 ...
	input = append(input, common.MustRlpEncode(txHashArr[0]))
	input = append(input, common.MustRlpEncode(uint8(1)))

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildDeclareInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2003)))				// func type code
	input = append(input, common.MustRlpEncode(nodeIdArr[0]))				// param 1 ...
	input = append(input, common.MustRlpEncode(uint32(1<<16 | 1<<8 | 1)))	//new version : 1.1.1

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildGetProposalInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2100)))				// func type code
	input = append(input, common.MustRlpEncode(txHashArr[0]))				// param 1 ...

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildGetTallyResultInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2101)))				// func type code
	input = append(input, common.MustRlpEncode(txHashArr[0]))				// param 1 ...

	return common.Bytes2Hex(common.MustRlpEncode(input))
}

func buildListProposalInput() string {
	var input [][]byte
	input = make([][]byte, 0)
	input = append(input, common.MustRlpEncode(uint16(2102)))				// func type code
	return common.Bytes2Hex(common.MustRlpEncode(input))
}


var successExpected = hexutil.Encode(common.MustRlpEncode(xcom.Result{true, "", ""}))

// each element means a call. we can reorder these elements to test different scenarios
var govContractCombinedTests = []vm.PrecompiledTest{
	{
		Input:		buildSubmitTextInput(),
		Expected:	successExpected,
		Name:		"submitText1",
	},
	{
		Input:		buildSubmitVersionInput(),
		Expected:	successExpected,
		Name:		"submitVersion1",
	},
	{
		Input:		buildVoteInput(),
		Expected:	successExpected,
		Name:		"vote1",
	},
	{
		Input:		buildDeclareInput(),
		Expected:	successExpected,
		Name:		"declare1",
	},
	{
		Input:		buildGetProposalInput(),
		Expected:	successExpected,
		Name:		"getProposal1",
	},
	/*
	{
		Input:		buildGetTallyResultInput(),
		Expected:	successExpected,
		Name:		"getTallyResult1",
	},
	*/
	{
		Input:		buildListProposalInput(),
		Expected:	successExpected,
		Name:		"listProposal1",
	},
}

func setup(t *testing.T) func() {
	t.Log("setup()......")

	precompiledContract := vm.PlatONPrecompiledContracts[commonvm.GovContractAddr]
	gc, _ = precompiledContract.(*vm.GovContract)
	gc.Evm = newEvm()
	gc.Contract = newContract(common.Big0)

	govPlugin = plugin.GovPluginInstance()
	gc.Plugin = govPlugin

	build_staking_data()

	snapdb = snapshotdb.Instance()

	plugin.StakingInstance()

	return func() {
		t.Log("tear down()......")
		snapdb.Clear()
	}
}

func testPlatONPrecompiled(idx int, t *testing.T) {

	test := govContractCombinedTests[idx]

	in := common.Hex2Bytes(test.Input)
	gc.Contract.Gas = gc.RequiredGas(in)

	state :=gc.Evm.StateDB.(*state.StateDB)

	state.Prepare(txHashArr[idx], blockHash, idx)

	t.Run(fmt.Sprintf("%s-Gas=%d", test.Name, gc.Contract.Gas), func(t *testing.T) {
		if res, err := vm.RunPlatONPrecompiledContract(gc, common.Hex2Bytes(test.Input), gc.Contract); err != nil {
			t.Error(err)
		} else if  common.Bytes2Hex0x(res) != test.Expected {

			t.Log(res)

			var r xcom.Result
			if err = rlp.DecodeBytes(res, &r); err != nil {
				t.Error(err)
			}else{
				t.Log(r.Data)

				fmt.Println("------ r.Data type--------")
				fmt.Println(reflect.TypeOf(r.Data).String())
				if reflect.TypeOf(r.Data).String() == "gov.TallyResult" {

					coded, _ := rlp.EncodeToBytes(r.Data)

					var tallyResulst gov.TallyResult
					if err = rlp.DecodeBytes(coded, &tallyResulst); nil!= err {
						fmt.Println("decode to TallyResulst error")
					}else {
						fmt.Println("decode to TallyResulst OK", tallyResulst)
					}
				}else if reflect.TypeOf(r.Data).String() == "[]uint8" {
					data, ok := r.Data.([]uint8)
					if ok {
						if data[0]== uint8(gov.Text) {
							coded, _ := rlp.EncodeToBytes(data[1:])
							var rlpData []byte
							if err = rlp.DecodeBytes(coded, &rlpData); nil!= err {
								fmt.Println("decode transfered data to []byte error")
							}else {
								var text gov.TextProposal
								if err = rlp.DecodeBytes(rlpData, &text); err != nil {
									fmt.Println("decode to text proposal failed", err)
								}else {
									fmt.Println("decode to text proposal OK", text)
								}
							}
						}
					}
				}else if reflect.TypeOf(r.Data).String() == "[]interface {}" {
					pDataList, ok := r.Data.([]interface {})
					if ok {
						for _, eachP := range pDataList{
							pType := eachP.([]uint8)[0]
							pByte := eachP.([]uint8)[1:]
							coded, _ := rlp.EncodeToBytes(pByte)
							var rlpData []byte
							if err = rlp.DecodeBytes(coded, &rlpData); nil!= err {
								fmt.Println("decode transfered data to []byte error")
							}else {
								if pType == uint8(gov.Text) {
									var text gov.TextProposal
									if err = rlp.DecodeBytes(rlpData, &text); err != nil {
										fmt.Println("decode to text proposal failed", err)
									}else {
										fmt.Println("decode to text proposal OK", text)
									}
								}else if pType == uint8(gov.Version) {
									var version gov.VersionProposal
									if err = rlp.DecodeBytes(rlpData, &version); err != nil {
										fmt.Println("decode to version proposal failed", err)
									}else {
										fmt.Println("decode to version proposal OK", version)
									}
								}
							}
						}
					}
				}
			}
		}
	})
}


// Tests the sample inputs from the elliptic curve pairing check EIP 197.
func TestPrecompiledGovContract(t *testing.T) {
	defer setup(t)()
	for i := 0; i < len(govContractCombinedTests); i++ {
		testPlatONPrecompiled( i, t)
	}
}
