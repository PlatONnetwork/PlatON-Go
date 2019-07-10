package vm_test

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"testing"
)

func TestSlashingContract_ReportMutiSign(t *testing.T) {
	state, err := newChainState()
	if nil != err {
		t.Error(err)
	}
	contract := &vm.SlashingContract{
		Plugin:   plugin.SlashInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber, blockHash, state),
	}

	plugin.SlashInstance().SetDecodeEvidenceFun(cbft.NewEvidences)

	state.Prepare(txHashArr[1], blockHash, 2)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(3000))
	dataStr := `{
          "duplicate_prepare": [
            {
              "VoteA": {
                "timestamp": 0,
                "block_hash": "0x0a0409021f020b080a16070609071c141f19011d090b091303121e1802130407",
                "block_number": 1,
                "validator_index": 1,
                "validator_address": "0x120b77ab712589ebd42d69003893ef962cc52832",
                "signature": "0xa65e16b3bc4862fdd893eaaaaecf1e415cdc2c8a08e4bbb1f6b2a1f4bf4e2d0c0ec27857da86a5f3150b32bee75322073cec320e51fe0a123cc4238ee4155bf001"
              },
              "VoteB": {
                "timestamp": 0,
                "block_hash": "0x18030d1e01071b1d071a12151e100a091f060801031917161e0a0d0f02161d0e",
                "block_number": 1,
                "validator_index": 1,
                "validator_address": "0x120b77ab712589ebd42d69003893ef962cc52832",
                "signature": "0x9126f9a339c8c4a873efc397062d67e9e9109895cd9da0d09a010d5f5ebbc6e76d285f7d87f801850c8552234101b651c8b7601b4ea077328c27e4f86d66a1bf00"
              }
            }
          ],
          "duplicate_viewchange": [],
          "timestamp_viewchange": []
        }`
	data, _ := rlp.EncodeToBytes(dataStr)

	params = append(params, fnType)
	params = append(params, data)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("ReportMutiSign encode rlp data fail: %v", err)
	} else {
		t.Log("ReportMutiSign data rlp: ", hexutil.Encode(buf.Bytes()))
	}
	runContract(contract, buf.Bytes(), t)
}

func TestSlashingContract_CheckMutiSign(t *testing.T) {
	state, err := newChainState()
	if nil != err {
		t.Error(err)
	}
	contract := &vm.SlashingContract{
		Plugin:   plugin.SlashInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber, blockHash, state),
	}

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(3001))
	typ, _ := rlp.EncodeToBytes(uint16(1))
	addr, _ := rlp.EncodeToBytes(common.HexToAddress("0x740ce31b3fac20dac379db243021a51e80111111"))
	blockNumber, _ := rlp.EncodeToBytes(uint64(10))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, addr)
	params = append(params, blockNumber)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("CheckMutiSign encode rlp data fail: %v", err)
	} else {
		t.Log("CheckMutiSign data rlp: ", hexutil.Encode(buf.Bytes()))
	}
	runContract(contract, buf.Bytes(), t)
}

func runContract(contract *vm.SlashingContract, buf []byte, t *testing.T) {
	res, err := contract.Run(buf)
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}
}