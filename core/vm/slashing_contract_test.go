// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package vm

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
	addr := common.MustBech32ToAddress("lax1r9tx0n00etv5c5smmlctlpg8jas7p78n8x3n9x")
	nodeId, err := discover.HexID("51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483")
	if nil != err {
		t.Fatal(err)
	}
	build_staking_data(genesis.Hash())
	newKey := staking.GetRoundValAddrArrKey(1)
	newValue := make([]common.NodeAddress, 0, 1)
	newValue = append(newValue, common.NodeAddress(addr))
	if err := staking.NewStakingDB().StoreRoundValidatorAddrs(blockHash, newKey, newValue); nil != err {
		t.Fatal(err)
	}
	contract := &SlashingContract{
		Plugin:   plugin.SlashInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, common.ZeroHash, state),
	}
	plugin.SlashInstance().SetDecodeEvidenceFun(evidence.NewEvidence)
	plugin.StakingInstance()
	plugin.GovPluginInstance()

	state.Prepare(txHashArr[1], common.ZeroHash, 2)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(3000))
	dupType, _ := rlp.EncodeToBytes(uint8(1))
	dataStr := `{
          "prepareA": {
           "epoch": 1,
           "viewNumber": 1,
           "blockHash": "0xbb6d4b83af8667929a9cb4918bbf790a97bb136775353765388d0add3437cde6",
           "blockNumber": 1,
           "blockIndex": 1,
           "blockData": "0x45b20c5ba595be254943aa57cc80562e84f1fb3bafbf4a414e30570c93a39579",
           "validateNode": {
            "index": 0,
            "nodeId": "51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483",
            "blsPubKey": "752fe419bbdc2d2222009e450f2932657bbc2370028d396ba556a49439fe1cc11903354dcb6dac552a124e0b3db0d90edcd334d7aabda0c3f1ade12ca22372f876212ac456d549dbbd04d2c8c8fb3e33760215e114b4d60313c142f7b8bbfd87"
           },
           "signature": "0x36015fee15253487e8125b86505377d8540b1a95d1a6b13f714baa55b12bd06ec7d5755a98230cdc88858470afa8cb0000000000000000000000000000000000"
          },
          "prepareB": {
           "epoch": 1,
           "viewNumber": 1,
           "blockHash": "0xf46c45f7ebb4a999dd030b9f799198b785654293dbe41aa7e909223af0e8c4ba",
           "blockNumber": 1,
           "blockIndex": 1,
           "blockData": "0xd630e96d127f55319392f20d4fd917e3e7cba19ad366c031b9dff05e056d9420",
           "validateNode": {
            "index": 0,
            "nodeId": "51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483",
            "blsPubKey": "752fe419bbdc2d2222009e450f2932657bbc2370028d396ba556a49439fe1cc11903354dcb6dac552a124e0b3db0d90edcd334d7aabda0c3f1ade12ca22372f876212ac456d549dbbd04d2c8c8fb3e33760215e114b4d60313c142f7b8bbfd87"
           },
           "signature": "0x783892b9b766f9f4c2a1d45b1fd53ca9ea56a82e38a998939edc17bc7fd756267d3c145c03bc6c1412302cf590645d8200000000000000000000000000000000"
          }
         }`
	data, _ := rlp.EncodeToBytes(dataStr)

	params = append(params, fnType)
	params = append(params, dupType)
	params = append(params, data)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		t.Fatalf("reportDuplicateSign encode rlp data fail: %v", err)
	} else {
		t.Log("reportDuplicateSign data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	var blsKey bls.SecretKey
	skbyte, err := hex.DecodeString("b36d4c3c3e8ee7fba3fbedcda4e0493e699cd95b68594093a8498c618680480a")
	if nil != err {
		t.Fatalf("reportDuplicateSign DecodeString byte data fail: %v", err)
	}
	blsKey.SetLittleEndian(skbyte)
	var blsKeyHex bls.PublicKeyHex
	b, _ := blsKey.GetPublicKey().MarshalText()
	if err := blsKeyHex.UnmarshalText(b); nil != err {
		t.Fatalf("Failed to blsKeyHex.UnmarshalText: %v", err)
		return
	}

	can := &staking.Candidate{
		CandidateBase: &staking.CandidateBase{
			NodeId:          nodeId,
			BlsPubKey:       blsKeyHex,
			StakingAddress:  addr,
			BenefitAddress:  addr,
			StakingBlockNum: blockNumber.Uint64(),
			StakingTxIndex:  1,
			ProgramVersion:  initProgramVersion,
		},
		CandidateMutable: &staking.CandidateMutable{
			Shares: new(big.Int).SetUint64(1000),

			Released:           common.Big256,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,
		},
	}
	state.CreateAccount(addr)
	state.AddBalance(addr, new(big.Int).SetUint64(1000000000000000000))
	if err := snapshotdb.Instance().NewBlock(blockNumber2, blockHash, common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	if err := plugin.StakingInstance().CreateCandidate(state, common.ZeroHash, blockNumber2, can.Shares, 0, common.NodeAddress(addr), can); nil != err {
		t.Fatal(err)
	}
	runContract(contract, buf.Bytes(), t)
}

func TestSlashingContract_CheckMutiSign(t *testing.T) {
	state, _, err := newChainState()
	if nil != err {
		t.Fatal(err)
	}
	contract := &SlashingContract{
		Plugin:   plugin.SlashInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}
	state.Prepare(txHashArr[1], blockHash, 2)

	var params [][]byte
	params = make([][]byte, 0)

	nodeId, err := discover.HexID("51c0559c065400151377d71acd7a17282a7c8abcfefdb11992dcecafde15e100b8e31e1a5e74834a04792d016f166c80b9923423fe280570e8131debf591d483")
	if nil != err {
		t.Fatal(err)
	}
	fnType, _ := rlp.EncodeToBytes(uint16(3001))
	typ, _ := rlp.EncodeToBytes(uint8(1))
	enNodeId, _ := rlp.EncodeToBytes(nodeId)
	blockNumber, _ := rlp.EncodeToBytes(uint16(1))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, enNodeId)
	params = append(params, blockNumber)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		t.Fatalf("checkDuplicateSign encode rlp data fail: %v", err)
	} else {
		t.Log("checkDuplicateSign data rlp: ", hexutil.Encode(buf.Bytes()))
	}
	runContract(contract, buf.Bytes(), t)
}

func runContract(contract *SlashingContract, buf []byte, t *testing.T) {
	res, err := contract.Run(buf)
	if nil != err {
		t.Fatal(err)
	} else {
		t.Log(string(res))
	}
}
