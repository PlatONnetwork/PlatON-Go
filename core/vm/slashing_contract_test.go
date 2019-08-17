package vm_test

import (
	"bytes"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/evidence"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
)

func TestSlashingContract_ReportMutiSign(t *testing.T) {
	state, genesis, err := newChainState()
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	if nil != err {
		t.Fatal(err)
	}
	build_staking_data(genesis.Hash())
	contract := &vm.SlashingContract{
		Plugin:   plugin.SlashInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}
	plugin.SlashInstance().SetDecodeEvidenceFun(evidence.NewEvidences)
	plugin.StakingInstance()
	plugin.GovPluginInstance()

	state.Prepare(txHashArr[1], blockHash, 2)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(3000))
	dataStr := `{
         "duplicate_prepare": [
          {
           "PrepareA": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0x09c94e00f687891f5de80146d906b55a249408dfd27afcad5a87bdad6fc28957",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x27383a8d350139588daba349dcd6ef1d745da004",
             "NodeID": "2560887689ce96e8a8361684c6b54061b6e4c667357284e8e301f8f51ff26efe4d7202708fda6fe4d5593188dacb5ce7114087d4c6840b529c48f617c6dff270",
             "blsPubKey": "8d0638bb1e58c33c12ea5735a5635cd51a26305ffda44f99f2190a28fa3ebd175db279caefd5f3b385d1fa04e7094499e78355efdd9fd96a08bc817963b42486"
            },
            "signature": "0x13eb58303156f63d8961a916694d3f659d58804ef1d783ee5c8c7fc3ca393b8a"
           },
           "PrepareB": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0xd1fc79053b8e9fd6a7d9061b4e12a282110429bd0e643aa477083f221a8cba8c",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x27383a8d350139588daba349dcd6ef1d745da004",
             "NodeID": "2560887689ce96e8a8361684c6b54061b6e4c667357284e8e301f8f51ff26efe4d7202708fda6fe4d5593188dacb5ce7114087d4c6840b529c48f617c6dff270",
             "blsPubKey": "8d0638bb1e58c33c12ea5735a5635cd51a26305ffda44f99f2190a28fa3ebd175db279caefd5f3b385d1fa04e7094499e78355efdd9fd96a08bc817963b42486"
            },
            "signature": "0x8de9fbb57edf75934b4caf40c95d569d03a75c762343066db737fdb2b818c313"
           }
          }
         ],
         "duplicate_vote": [],
         "duplicate_viewchange": []
        }`
	data, _ := rlp.EncodeToBytes(dataStr)

	params = append(params, fnType)
	params = append(params, data)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		t.Fatalf("ReportDuplicateSign encode rlp data fail: %v", err)
	} else {
		t.Log("ReportDuplicateSign data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	addr := common.HexToAddress("0x27383a8d350139588daba349dcd6ef1d745da004")
	nodeId, err := discover.HexID("2560887689ce96e8a8361684c6b54061b6e4c667357284e8e301f8f51ff26efe4d7202708fda6fe4d5593188dacb5ce7114087d4c6840b529c48f617c6dff270")
	if nil != err {
		t.Fatal(err)
	}
	var blsKey bls.SecretKey
	skbyte, err := hex.DecodeString("328146a789365c846b1da23f89bc178a32febf217fb496153efe9e34ac4fa621")
	if nil != err {
		t.Fatalf("ReportDuplicateSign DecodeString byte data fail: %v", err)
	}
	blsKey.SetLittleEndian(skbyte)
	can := &staking.Candidate{
		NodeId:          nodeId,
		BlsPubKey:       *blsKey.GetPublicKey(),
		StakingAddress:  addr,
		BenefitAddress:  addr,
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  1,
		ProgramVersion:  initProgramVersion,
		Shares:          new(big.Int).SetUint64(1000),

		Released:           common.Big0,
		ReleasedHes:        common.Big0,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
	}
	state.CreateAccount(addr)
	state.AddBalance(addr, new(big.Int).SetUint64(1000000000000000000))
	if err := plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, can.Shares, 0, addr, can); nil != err {
		t.Fatal(err)
	}
	runContract(contract, buf.Bytes(), t)
}

func TestSlashingContract_CheckMutiSign(t *testing.T) {
	state, _, err := newChainState()
	if nil != err {
		t.Fatal(err)
	}
	contract := &vm.SlashingContract{
		Plugin:   plugin.SlashInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}
	state.Prepare(txHashArr[1], blockHash, 2)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(3001))
	typ, _ := rlp.EncodeToBytes(uint32(1))
	addr, _ := rlp.EncodeToBytes(common.HexToAddress("0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257"))
	blockNumber, _ := rlp.EncodeToBytes(uint16(1))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, addr)
	params = append(params, blockNumber)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		t.Fatalf("CheckDuplicateSign encode rlp data fail: %v", err)
	} else {
		t.Log("CheckDuplicateSign data rlp: ", hexutil.Encode(buf.Bytes()))
	}
	runContract(contract, buf.Bytes(), t)
}

func runContract(contract *vm.SlashingContract, buf []byte, t *testing.T) {
	res, err := contract.Run(buf)
	if nil != err {
		t.Fatal(err)
	} else {
		t.Log(string(res))
	}
}
