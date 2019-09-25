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
	addr := common.HexToAddress("0x076c72c53c569df9998448832a61371ac76d0d05")
	nodeId, err := discover.HexID("b68b23496b820f4133e42b747f1d4f17b7fd1cb6b065c613254a5717d856f7a56dabdb0e30657f18fb9074c7cb60eb62a6b35ad61898da407dae2cb8efe68511")
	if nil != err {
		t.Fatal(err)
	}
	build_staking_data(genesis.Hash())
	newKey := staking.GetRoundValAddrArrKey(1)
	newValue := make([]common.Address, 0, 1)
	newValue = append(newValue, addr)
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
            "blockHash": "0x86c86e7ddb977fbd2f1d0b5cb92510c230775deef02b60d161c3912244473b54",
            "blockNumber": 1,
            "blockIndex": 1,
            "validateNode": {
             "index": 0,
             "address": "0x076c72c53c569df9998448832a61371ac76d0d05",
             "nodeId": "b68b23496b820f4133e42b747f1d4f17b7fd1cb6b065c613254a5717d856f7a56dabdb0e30657f18fb9074c7cb60eb62a6b35ad61898da407dae2cb8efe68511",
             "blsPubKey": "6021741b867202a3e60b91452d80e98f148aefadbb5ff1860f1fec5a8af14be20ca81fd73c231d6f67d4c9d2d516ac1297c8126ed7c441e476c0623c157638ea3b5b2189f3a20a78b2fd5fb32e5d7de055e4d2a0c181d05892be59cf01f8ab88"
            },
            "signature": "0x8c77b2178239fd525b774845cc7437ecdf5e6175ab4cc49dcb93eae6df288fd978e5290f59420f93bba22effd768f38900000000000000000000000000000000"
           },
           "prepareB": {
            "epoch": 1,
            "viewNumber": 1,
            "blockHash": "0xeccd7a0b7793a74615721e883ab5223de30c5cf4d2ced9ab9dfc782e8604d416",
            "blockNumber": 1,
            "blockIndex": 1,
            "validateNode": {
             "index": 0,
             "address": "0x076c72c53c569df9998448832a61371ac76d0d05",
             "nodeId": "b68b23496b820f4133e42b747f1d4f17b7fd1cb6b065c613254a5717d856f7a56dabdb0e30657f18fb9074c7cb60eb62a6b35ad61898da407dae2cb8efe68511",
             "blsPubKey": "6021741b867202a3e60b91452d80e98f148aefadbb5ff1860f1fec5a8af14be20ca81fd73c231d6f67d4c9d2d516ac1297c8126ed7c441e476c0623c157638ea3b5b2189f3a20a78b2fd5fb32e5d7de055e4d2a0c181d05892be59cf01f8ab88"
            },
            "signature": "0x5213b4122f8f86874f537fa9eda702bba2e47a7b8ecc0ff997101d675a174ee5884ec85e8ea5155c5a6ad6b55326670d00000000000000000000000000000000"
           }
          }`
	data, _ := rlp.EncodeToBytes(dataStr)

	params = append(params, fnType)
	params = append(params, dupType)
	params = append(params, data)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		t.Fatalf("ReportDuplicateSign encode rlp data fail: %v", err)
	} else {
		t.Log("ReportDuplicateSign data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	var blsKey bls.SecretKey
	skbyte, err := hex.DecodeString("155b9a6f5575b9b5a4d8658f616660a549674b36c858e6c606d08ec5c20c4637")
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

		Released:           common.Big256,
		ReleasedHes:        common.Big0,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
	}
	state.CreateAccount(addr)
	state.AddBalance(addr, new(big.Int).SetUint64(1000000000000000000))
	if err := snapshotdb.Instance().NewBlock(blockNumber2, blockHash, common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	if err := plugin.StakingInstance().CreateCandidate(state, common.ZeroHash, blockNumber2, can.Shares, 0, addr, can); nil != err {
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

	fnType, _ := rlp.EncodeToBytes(uint16(3001))
	typ, _ := rlp.EncodeToBytes(uint8(1))
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

func runContract(contract *SlashingContract, buf []byte, t *testing.T) {
	res, err := contract.Run(buf)
	if nil != err {
		t.Fatal(err)
	} else {
		t.Log(string(res))
	}
}
