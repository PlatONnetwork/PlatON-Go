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
            "block_hash": "0x504fc256d64711833c5e9ab5968ef3ae9129af90a6f48ea6125c9a98bf0643a2",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0xb1950823ca8fcd02283e18abd28a8b7d5e1951f3",
             "NodeID": "f58de166211ed50e510f9bb0453bc6c93fa6a2f83ab5e10155fb1f52ecb3d8c1a79a406ebca6b4171a03c0a5169cde60e406852c31627924d4f2b1f7d889f7a9",
             "blsPubKey": "cabac737d66770861eba0bc233af9a1ebdee32a21bedfed37f3ab1f8f493a9009b6d3f1a96c96da6492f2547dfc39374e6de25805db601dc66748a1aad8c740c"
            },
            "cannibalize": "POWEymMIHesJDiWxCp5BIhWFukqtiWKjJRnFO6dv00k=",
            "signature": "0xa50083eb0bac47298aa7094f5babc75376f332aebdb9781b8166b813cc1dfa81"
           },
           "PrepareB": {
            "epoch": 1,
            "view_number": 1,
            "block_hash": "0xc34c83e56b31e40f4960460c8d70fcff27a8a8b14a69205d33e70f066e5e291c",
            "block_number": 1,
            "block_index": 1,
            "validate_node": {
             "index": 0,
             "address": "0xb1950823ca8fcd02283e18abd28a8b7d5e1951f3",
             "NodeID": "f58de166211ed50e510f9bb0453bc6c93fa6a2f83ab5e10155fb1f52ecb3d8c1a79a406ebca6b4171a03c0a5169cde60e406852c31627924d4f2b1f7d889f7a9",
             "blsPubKey": "cabac737d66770861eba0bc233af9a1ebdee32a21bedfed37f3ab1f8f493a9009b6d3f1a96c96da6492f2547dfc39374e6de25805db601dc66748a1aad8c740c"
            },
            "cannibalize": "0ubn5EnNGfC08PAxrsaU30JyKfIStQBpDecQPqV1Gsw=",
            "signature": "0xf462d8b58b5fd6282f1da21287283baba225fdffecbea4c4cabee88f3868209b"
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

	addr := common.HexToAddress("0xb1950823ca8fcd02283e18abd28a8b7d5e1951f3")
	nodeId, err := discover.HexID("0xf58de166211ed50e510f9bb0453bc6c93fa6a2f83ab5e10155fb1f52ecb3d8c1a79a406ebca6b4171a03c0a5169cde60e406852c31627924d4f2b1f7d889f7a9")
	if nil != err {
		t.Fatal(err)
	}
	var blsKey bls.SecretKey
	skbyte, err := hex.DecodeString("8f7358f97aec6eccb400f878357e0ae87c93b3d1e8f6da68fe77438b9f7ec01d")
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
	addr, _ := rlp.EncodeToBytes(common.HexToAddress("0xb1950823ca8fcd02283e18abd28a8b7d5e1951f3"))
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
