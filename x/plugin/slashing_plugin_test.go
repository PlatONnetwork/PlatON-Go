package plugin_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
	mrand "math/rand"
	"strconv"
	"testing"
	"time"
)

func initInfo(t *testing.T) (*plugin.SlashingPlugin, xcom.StateDB) {
	si := plugin.SlashInstance()
	plugin.StakingInstance()
	db := ethdb.NewMemDatabase()
	stateDB, err := state.New(common.Hash{}, state.NewDatabase(db))
	if nil != err {
		t.Error(err)
	}
	return si, stateDB
}

func TestSlashingPlugin_BeginBlock(t *testing.T) {
	/*
	si, stateDB := initInfo(t)
	if _, err := si.EndBlock(hash, header, stateDB); nil != err {
		t.Error(err)
	}*/
}

func TestSlashingPlugin_Confirmed(t *testing.T) {
	si, _ := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	blockNumber := new(big.Int).SetUint64(1)
	pri, err := crypto.GenerateKey()
	if err != nil {
		panic(err)
	}
	for i := 0; i < 251; i++ {
		header := &types.Header{
			Number:blockNumber,
			Extra:make([]byte, 97),
		}
		sign, err := crypto.Sign(header.SealHash().Bytes(), pri)
		if nil != err {
			t.Error(err)
		}
		copy(header.Extra[len(header.Extra)-common.ExtraSeal:], sign[:])
		block := types.NewBlock(header, nil, nil)
		if err := si.Confirmed(block); nil != err {
			t.Error(err)
		}
		blockNumber.Add(blockNumber, new(big.Int).SetUint64(1))
	}
	result, err := si.GetPreNodeAmount()
	if nil != err {
		t.Error(err)
	}
	fmt.Println(result)
}
func TestSlashingPlugin_Slash(t *testing.T) {
	si, stateDB := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	si.SetDecodeEvidenceFun(cbft.NewEvidences)
	data := `{
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
	blockNumber := new(big.Int).SetUint64(1)
	if err := si.Slash(data, common.ZeroHash, blockNumber.Uint64(), stateDB); nil != err {
		t.Error(err)
	}
}

func TestSlashingPlugin_CheckMutiSign(t *testing.T) {
	si, stateDB := initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	addr := common.ZeroAddr
	if _,_, err := si.CheckMutiSign(addr, 1, 1, stateDB); nil != err {
		t.Error(err)
	}
}

func TestStakingPlugin_ProbabilityElection(t *testing.T) {
	initInfo(t)
	defer func() {
		snapshotdb.Instance().Clear()
	}()
	curve := elliptic.P256()
	vqList := make(staking.ValidatorQueue, 0)
	preNonces := make([][]byte, 0)
	curentNonce := crypto.Keccak256([]byte(string("nonce")))
	for i := 0; i < int(xcom.EpochValidatorNum); i++ {
		privKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
		nodeId := discover.PubkeyID(&privKey.PublicKey)
		addr := crypto.PubkeyToAddress(privKey.PublicKey)
		mrand.Seed(time.Now().UnixNano())
		stakingWeight := [4]string{}
		stakingWeight[0] = "1"
		stakingWeight[1] = strconv.Itoa(mrand.Intn(100000000))
		stakingWeight[2] = strconv.Itoa(mrand.Intn(230))
		stakingWeight[3] = strconv.Itoa(mrand.Intn(1000))
		v := &staking.Validator{
			NodeAddress:addr,
			NodeId:nodeId,
			StakingWeight:stakingWeight,
			ValidatorTerm:1,
		}
		vqList = append(vqList, v)
		preNonces = append(preNonces, crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(i))))[:])
		time.Sleep(time.Microsecond * 10)
	}
	for _, v := range vqList {
		t.Log("Generate Validator", "addr", hex.EncodeToString(v.NodeAddress.Bytes()), "stakingWeight", v.StakingWeight)
	}
	result, err := plugin.StakingInstance().ProbabilityElection(vqList, curentNonce, preNonces)
	if nil != err {
		t.Error(err)
	}
	t.Log("election success", result)
	for _, v := range result {
		t.Log("Validator", "addr", hex.EncodeToString(v.NodeAddress.Bytes()), "stakingWeight", v.StakingWeight)
	}
}