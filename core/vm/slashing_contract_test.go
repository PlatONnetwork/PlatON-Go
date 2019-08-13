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
		Contract: newContract(common.Big0),
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
            "block_hash": "0x00a452c6116ac9df049016437f8a35b4e29c17d63632314f0266df2b0dcd4bef",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257",
             "NodeID": "bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c",
             "blsPubKey": "f3ee4cb60b04358c21460b9dd0832028959a6d0052218d796c96a5eac01b541f88595d62ee52880e0a77ecf8ffde63966a5d0d70028c08dfca622827563df99e"
            },
            "signature": "0x554a2a2f1b0d197730c707b595016b5f735ce0df0a5e9efd28a77764f295af1f"
           },
           "PrepareB": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0x3f643f315a72d54815e3b638e53a1f293834e6d9109c4c0f3e5d9c7171bf1cf2",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257",
             "NodeID": "bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c",
             "blsPubKey": "f3ee4cb60b04358c21460b9dd0832028959a6d0052218d796c96a5eac01b541f88595d62ee52880e0a77ecf8ffde63966a5d0d70028c08dfca622827563df99e"
            },
            "signature": "0x9e626bd0fd19290c7ff23a605259735de216f9c26df51ddaf51f66f0aade4097"
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

	addr := common.HexToAddress("0x9e3e0f0f366b26b965f3aa3ed67603fb480b1257")
	nodeId, err := discover.HexID("bf1c6f0159513755be9bbb12da983c0743f0e8553c07f40e5e3c07eba84c6584aec141ed2e87e94ababee483e7d4809e85f9e2d043d0cb73bd46149fbc2f2f8c")
	if nil != err {
		t.Fatal(err)
	}
	var blsKey bls.SecretKey
	skbyte, err := hex.DecodeString("d6ba381339988d7984393cd1892969d78eae735c588a9528c834676faf333507")
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
		Contract: newContract(common.Big0),
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
