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

package plugin

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/big"
	mrand "math/rand"
	"os"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/crypto/vrf"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/x/reward"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/x/handler"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

/**
test tool
*/
func Test_CleanSnapshotDB(t *testing.T) {
	sndb := snapshotdb.Instance()
	sndb.Clear()
}

func PrintObject(s string, obj interface{}) {
	objs, _ := json.Marshal(obj)
	log.Debug(s + " == " + string(objs))
}

func watching(eventMux *event.TypeMux, t *testing.T) {
	events := eventMux.Subscribe(cbfttypes.AddValidatorEvent{})
	defer events.Unsubscribe()

	for {
		select {
		case ev := <-events.Chan():
			if ev == nil {
				t.Error("ev is nil, may be Server closing")
				continue
			}

			switch ev.Data.(type) {
			case cbfttypes.AddValidatorEvent:
				addEv, ok := ev.Data.(cbfttypes.AddValidatorEvent)
				if !ok {
					t.Error("Received add validator event type error")
					continue
				}

				str, _ := json.Marshal(addEv)
				t.Log("P2P Received the add validator is:", string(str))
			default:
				t.Error("Received unexcepted event")
			}

		}
	}
}

func build_vrf_Nonce() ([]byte, [][]byte) {
	preNonces := make([][]byte, 0)
	curentNonce := crypto.Keccak256([]byte(string("nonce")))
	for i := 0; i < int(xcom.MaxValidators()); i++ {
		preNonces = append(preNonces, crypto.Keccak256([]byte(time.Now().Add(time.Duration(i)).String())[:]))
		time.Sleep(time.Microsecond * 10)
	}
	return curentNonce, preNonces
}

func buildPrepareData(genesis *types.Block, t *testing.T) (*types.Header, error) {
	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if nil != err {
		t.Error("Failed to rlp vrf nonces", "err", err)
		return nil, err
	}

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.MaxValidators())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return nil, err
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return nil, err
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return nil, err
		}

		canTmp := &staking.Candidate{
			CandidateBase: &staking.CandidateBase{
				NodeId:          nodeId,
				BlsPubKey:       blsKeyHex,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(j),
				StakingTxIndex:  uint32(index),
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
				},
			},
			CandidateMutable: &staking.CandidateMutable{
				Shares: balance,

				// Prevent null pointer initialization
				Released:           common.Big0,
				ReleasedHes:        common.Big0,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canTmp.ProgramVersion, canTmp.Shares, canTmp.StakingBlockNum, canTmp.StakingTxIndex, canTmp.NodeId)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return nil, err
		}

		// Store Candidate Base info
		canBaseKey := staking.CanBaseKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp.CandidateBase); nil != err {
			t.Errorf("Failed to Store CandidateBase info: PutBaseDB failed. error:%s", err.Error())
			return nil, err
		} else {

			if err := sndb.PutBaseDB(canBaseKey, val); nil != err {
				t.Errorf("Failed to Store CandidateBase info: PutBaseDB failed. error:%s", err.Error())
				return nil, err
			}
		}

		// Store Candidate Mutable info
		canMutableKey := staking.CanMutableKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canTmp.CandidateMutable); nil != err {
			t.Errorf("Failed to Store CandidateMutable info: PutBaseDB failed. error:%s", err.Error())
			return nil, err
		} else {

			if err := sndb.PutBaseDB(canMutableKey, val); nil != err {
				t.Errorf("Failed to Store CandidateMutable info: PutBaseDB failed. error:%s", err.Error())
				return nil, err
			}
		}

		if j < int(xcom.MaxValidators()) {
			v := &staking.Validator{
				NodeAddress:     canAddr,
				NodeId:          canTmp.NodeId,
				BlsPubKey:       canTmp.BlsPubKey,
				ProgramVersion:  canTmp.ProgramVersion,
				Shares:          canTmp.Shares,
				StakingBlockNum: canTmp.StakingBlockNum,
				StakingTxIndex:  canTmp.StakingTxIndex,
				ValidatorTerm:   0,
			}
			validatorQueue[j] = v
		}
	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return nil, err
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return nil, err
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return nil, err
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return nil, err
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error:%s", err.Error())
		return nil, err
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return nil, err
	}

	PrintObject("Test round", validatorQueue[:xcom.MaxConsensusVals()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.MaxConsensusVals()])
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return nil, err
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return nil, err
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return nil, err
	}

	// SetCurrent to snapshotDB
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return nil, err
	}
	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	newNumber := big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(newNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	// new block
	nonce := crypto.Keccak256([]byte(time.Now().Add(time.Duration(1)).String()))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      newNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(newNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return nil, err
	}

	return header, err
}

func create_staking(state xcom.StateDB, blockNumber *big.Int, blockHash common.Hash, index int, typ uint16, t *testing.T) error {

	balance, _ := new(big.Int).SetString(balanceStr[index], 10)
	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	canTmp := &staking.Candidate{}

	var blsKeyHex bls.PublicKeyHex

	b, _ := blsKey.GetPublicKey().MarshalText()
	err := blsKeyHex.UnmarshalText(b)
	if nil != err {
		log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
		return err
	}

	canBase := &staking.CandidateBase{
		NodeId:          nodeIdArr[index],
		BlsPubKey:       blsKeyHex,
		StakingAddress:  sender,
		BenefitAddress:  addrArr[index],
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  uint32(index),
		ProgramVersion:  xutil.CalcVersion(initProgramVersion),

		// Prevent null pointer initialization

		Description: staking.Description{
			NodeName:   nodeNameArr[index],
			ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+1)] + "balabalala" + chaList[index],
			Website:    "www." + nodeNameArr[index] + ".org",
			Details:    "This is " + nodeNameArr[index] + " Super Node",
		},
	}

	canMutable := &staking.CandidateMutable{
		Shares: balance,
		// Prevent null pointer initialization
		Released:               common.Big0,
		ReleasedHes:            common.Big0,
		RestrictingPlan:        common.Big0,
		RestrictingPlanHes:     common.Big0,
		WithdrewDelegateEpoch:  0,
		WithdrewDelegateAmount: common.Big0,
	}

	canTmp.CandidateBase = canBase
	canTmp.CandidateMutable = canMutable

	canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

	return StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, typ, canAddr, canTmp)
}

func getCandidate(blockHash common.Hash, index int) (*staking.Candidate, error) {
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	if can, err := StakingInstance().GetCandidateInfo(blockHash, addr); nil != err {
		return nil, err
	} else {

		return can, nil
	}
}

func delegate(state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	can *staking.Candidate, typ uint16, index int, t *testing.T) (*staking.Delegation, error) {

	delAddr := addrArr[index+1]

	// build delegate
	del := new(staking.Delegation)

	// Prevent null pointer initialization
	del.Released = common.Big0
	del.RestrictingPlan = common.Big0
	del.ReleasedHes = common.Big0
	del.RestrictingPlanHes = common.Big0
	del.WithdrewEpoch = 0
	del.WithdrewAmount = new(big.Int).Set(common.Big0)
	del.UnLockEpoch = 0
	del.CumulativeIncome = new(big.Int).Set(common.Big0)
	//amount := common.Big257  // FAIL
	amount, _ := new(big.Int).SetString(balanceStr[index+1], 10) // PASS

	canAddr, _ := xutil.NodeId2Addr(can.NodeId)

	delegateRewardPerList := make([]*reward.DelegateRewardPer, 0)

	return del, StakingInstance().Delegate(state, blockHash, blockNumber, delAddr, del, canAddr, can, 0, amount, delegateRewardPerList)
}

func getDelegate(blockHash common.Hash, stakingNum uint64, index int, t *testing.T) *staking.Delegation {

	del, err := StakingInstance().GetDelegateInfo(blockHash, addrArr[index+1], nodeIdArr[index], stakingNum)
	if nil != err {
		t.Log("Failed to GetDelegateInfo:", err)
	} else {
		delByte, _ := json.Marshal(del)
		t.Log("Get Delegate Info is:", string(delByte))
	}
	return del
}

/**
Standard test cases
*/

func TestStakingPlugin_BeginBlock(t *testing.T) {
	// nothings in that
}

func TestStakingPlugin_EndBlock(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to rlp vrf nonces: %v", err)) {
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()

	if !assert.Nil(t, err, fmt.Sprintf("Failed to generate random Address private key: %v", err)) {
		return
	}

	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.MaxValidators())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if !assert.Nil(t, err, fmt.Sprintf("Failed to generate random NodeId private key: %v", err)) {
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if !assert.Nil(t, err, fmt.Sprintf("Failed to generate random Address private key: %v", err)) {
			return
		}
		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		var blsKeyHex bls.PublicKeyHex
		err = blsKeyHex.UnmarshalText(blsKey.Serialize())
		if nil != err {
			return
		}

		canBase := &staking.CandidateBase{
			NodeId:          nodeId,
			BlsPubKey:       blsKeyHex,
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}
		canMutable := &staking.CandidateMutable{
			Shares: balance,
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,
		}

		canAddr, _ := xutil.NodeId2Addr(canBase.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canBase.ProgramVersion, canMutable.Shares, canBase.StakingBlockNum, canBase.StakingTxIndex, canBase.NodeId)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store CandidateBase info
		canBaseKey := staking.CanBaseKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canBase); nil != err {
			t.Errorf("Failed to Store Candidate Base info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canBaseKey, val); nil != err {
				t.Errorf("Failed to Store Candidate Base info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		// Store CandidateMutable info
		canMutableKey := staking.CanMutableKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canMutable); nil != err {
			t.Errorf("Failed to Store Candidate Mutable info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canMutableKey, val); nil != err {
				t.Errorf("Failed to Store Candidate Mutable info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.MaxValidators()) {
			v := &staking.Validator{
				NodeAddress:     canAddr,
				NodeId:          canBase.NodeId,
				BlsPubKey:       canBase.BlsPubKey,
				ProgramVersion:  canBase.ProgramVersion,
				Shares:          canMutable.Shares,
				StakingBlockNum: canBase.StakingBlockNum,
				StakingTxIndex:  canBase.StakingTxIndex,

				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}

	}

	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error:%s", err.Error())
		return
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		t.Errorf("Failed to rlp encodeing genesis validators. error:%s", err.Error())
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error: %v", err)) {
		return
	}
	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	PrintObject("Test round", validatorQueue[:xcom.MaxConsensusVals()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.MaxConsensusVals()])
	if !assert.Nil(t, err, fmt.Sprintf("Failed to rlp encodeing genesis validators. error: %v", err)) {
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		t.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error())
		return
	}

	/**
	EndBlock to Election()
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(time.Now().Add(time.Duration(1)).String()))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	err = StakingInstance().EndBlock(currentHash, header, state)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to EndBlock, blockNumber: %d, err: %v", currentNumber, err)) {
		return
	}

	if err := sndb.Commit(currentHash); nil != err {
		t.Errorf("Failed to Commit, blockNumber: %d, blockHHash: %s, err: %v", currentNumber, currentHash.Hex(), err)
		return
	}

	if err := sndb.Compaction(); nil != err {
		t.Errorf("Failed to Compaction, blockNumber: %d, blockHHash: %s, err: %v", currentNumber, currentHash.Hex(), err)
		return
	}

	// new block
	privateKey2, err := crypto.GenerateKey()
	if nil != err {
		t.Errorf("Failed to generate random Address private key: %v", err)
		return
	}
	nodeId2 := discover.PubkeyID(&privateKey2.PublicKey)
	currentHash = crypto.Keccak256Hash([]byte(nodeId2.String()))

	/**
	Elect Epoch validator list  == ElectionNextList()
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() * xutil.EpochSize())) // 600

	preNum := new(big.Int).Sub(currentNumber, big.NewInt(1)) // 599

	if err := sndb.SetCurrent(currentHash, *preNum, *preNum); nil != err {
		panic(fmt.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error()))
	}

	nonce = crypto.Keccak256([]byte(time.Now().Add(time.Duration(1)).String()))[:]
	header = &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	err = StakingInstance().EndBlock(currentHash, header, state)
	assert.Nil(t, err, fmt.Sprintf("Failed to Election, blockNumber: %d, err: %v", currentNumber, err))
}

func TestStakingPlugin_Confirmed(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// New VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	// build vrf proof
	// build ancestor nonces
	_, nonces := build_vrf_Nonce()
	enValue, err := rlp.EncodeToBytes(nonces)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to rlp vrf nonces: %v", err)) {
		return
	}

	// new block
	privateKey, err := crypto.GenerateKey()
	if !assert.Nil(t, err, fmt.Sprintf("Failed to generate random Address private key: %v", err)) {
		return
	}

	nodeId := discover.PubkeyID(&privateKey.PublicKey)
	currentHash := crypto.Keccak256Hash([]byte(nodeId.String()))
	currentNumber := big.NewInt(1)

	// build genesis veriferList and validatorList
	validatorQueue := make(staking.ValidatorQueue, xcom.MaxValidators())

	for j := 0; j < 1000; j++ {
		var index int = j % 25

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if !assert.Nil(t, err, fmt.Sprintf("Failed to generate random Address private key: %v", err)) {
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		var blsKeyHex bls.PublicKeyHex
		err = blsKeyHex.UnmarshalText(blsKey.Serialize())
		if nil != err {
			return
		}

		canBase := &staking.CandidateBase{
			NodeId:          nodeId,
			BlsPubKey:       blsKeyHex,
			StakingAddress:  sender,
			BenefitAddress:  addr,
			StakingBlockNum: uint64(1),
			StakingTxIndex:  uint32(index),
			ProgramVersion:  xutil.CalcVersion(initProgramVersion),

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(j),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(j) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(j) + " Super Node",
			},
		}
		canMutable := &staking.CandidateMutable{
			Shares: balance,
			// Prevent null pointer initialization
			Released:           common.Big0,
			ReleasedHes:        common.Big0,
			RestrictingPlan:    common.Big0,
			RestrictingPlanHes: common.Big0,
		}

		canAddr, _ := xutil.NodeId2Addr(canBase.NodeId)

		// Store Candidate power
		powerKey := staking.TallyPowerKey(canBase.ProgramVersion, canMutable.Shares, canBase.StakingBlockNum, canBase.StakingTxIndex, canBase.NodeId)
		if err := sndb.PutBaseDB(powerKey, canAddr.Bytes()); nil != err {
			t.Errorf("Failed to Store Candidate Power: PutBaseDB failed. error:%s", err.Error())
			return
		}

		// Store CandidateBase info
		canBaseKey := staking.CanBaseKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canBase); nil != err {
			t.Errorf("Failed to Store Candidate Base info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canBaseKey, val); nil != err {
				t.Errorf("Failed to Store Candidate Base info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		// Store CandidateMutable info
		canMutableKey := staking.CanMutableKeyByAddr(canAddr)
		if val, err := rlp.EncodeToBytes(canMutable); nil != err {
			t.Errorf("Failed to Store Candidate Mutable info: PutBaseDB failed. error:%s", err.Error())
			return
		} else {

			if err := sndb.PutBaseDB(canMutableKey, val); nil != err {
				t.Errorf("Failed to Store Candidate Mutable info: PutBaseDB failed. error:%s", err.Error())
				return
			}
		}

		if j < int(xcom.MaxValidators()) {
			v := &staking.Validator{
				NodeAddress:     canAddr,
				NodeId:          canBase.NodeId,
				BlsPubKey:       canBase.BlsPubKey,
				ProgramVersion:  canBase.ProgramVersion,
				Shares:          canMutable.Shares,
				StakingBlockNum: canBase.StakingBlockNum,
				StakingTxIndex:  canBase.StakingTxIndex,

				ValidatorTerm: 0,
			}
			validatorQueue[j] = v
		}
	}
	/**
	*******
	build genesis epoch validators
	*******
	*/
	verifierIndex := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
	}

	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)

	// current epoch start and end indexs
	epoch_index, err := rlp.EncodeToBytes(epochIndexArr)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to Store Epoch Validators start and end index: rlp encodeing failed. error: %v", err)) {
		return
	}

	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epoch_index); nil != err {
		t.Errorf("Failed to Store Epoch Validators start and end index: PutBaseDB failed. error:%s", err.Error())
		return
	}

	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to rlp encodeing genesis validators. error: %v", err)) {
		return
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		t.Errorf("Failed to Store Epoch Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	/**
	*******
	build genesis curr round validators
	*******
	*/
	curr_indexInfo := &staking.ValArrIndex{
		Start: 1,
		End:   xutil.ConsensusSize(),
	}
	roundIndexArr := make(staking.ValArrIndexQueue, 0)
	roundIndexArr = append(roundIndexArr, curr_indexInfo)

	// round index
	round_index, err := rlp.EncodeToBytes(roundIndexArr)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to Store Round Validators start and end indexs: rlp encodeing failed. error: %v", err)) {
		return
	}

	if err := sndb.PutBaseDB(staking.GetRoundIndexKey(), round_index); nil != err {
		t.Errorf("Failed to Store Round Validators start and end indexs: PutBaseDB failed. error:%s", err.Error())
		return
	}

	PrintObject("Test round", validatorQueue[:xcom.MaxConsensusVals()])
	roundArr, err := rlp.EncodeToBytes(validatorQueue[:xcom.MaxConsensusVals()])
	if !assert.Nil(t, err, fmt.Sprintf("Failed to rlp encodeing genesis validators. error: %v", err)) {
		return
	}
	// Store Current Round validator
	if err := sndb.PutBaseDB(staking.GetRoundValArrKey(curr_indexInfo.Start, curr_indexInfo.End), roundArr); nil != err {
		t.Errorf("Failed to Store Current Round Validators: PutBaseDB failed. error:%s", err.Error())
		return
	}

	// Store vrf nonces
	if err := sndb.PutBaseDB(handler.NonceStorageKey, enValue); nil != err {
		t.Errorf("Failed to Store Current Vrf nonces : PutBaseDB failed. error:%s", err.Error())
		return
	}

	// SetCurrent to snapshotDB
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50
	preNum1 := new(big.Int).Sub(currentNumber, big.NewInt(1))
	if err := sndb.SetCurrent(currentHash, *preNum1, *preNum1); nil != err {
		t.Errorf("Failed to SetCurrent by snapshotdb. error:%s", err.Error())
		return
	}

	/**
	EndBlock to Election()
	*/
	// new block
	currentNumber = big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())) // 50

	nonce := crypto.Keccak256([]byte(time.Now().Add(time.Duration(1)).String()))[:]
	header := &types.Header{
		ParentHash:  currentHash,
		Coinbase:    sender,
		Root:        common.ZeroHash,
		TxHash:      types.EmptyRootHash,
		ReceiptHash: types.EmptyRootHash,
		Number:      currentNumber,
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		Nonce:       types.EncodeNonce(nonce),
	}
	currentHash = header.Hash()

	if err := sndb.NewBlock(currentNumber, header.ParentHash, currentHash); nil != err {
		t.Errorf("Failed to snapshotDB New Block, err: %v", err)
		return
	}

	err = StakingInstance().EndBlock(currentHash, header, state)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to EndBlock, blockNumber: %d, err: %v", currentNumber, err)) {
		return
	}

	/**
	Start Confirmed
	*/

	eventMux := &event.TypeMux{}
	StakingInstance().SetEventMux(eventMux)
	go watching(eventMux, t)

	blockElection := types.NewBlock(header, nil, nil)

	next, err := StakingInstance().getNextValList(blockElection.Hash(), blockElection.Number().Uint64(), QueryStartNotIrr)

	assert.Nil(t, err, fmt.Sprintf("Failed to getNextValList, blockNumber: %d, err: %v", blockElection.Number().Uint64(), err))

	err = StakingInstance().Confirmed(next.Arr[0].NodeId, blockElection)
	assert.Nil(t, err, fmt.Sprintf("Failed to Confirmed, blockNumber: %d, err: %v", blockElection.Number().Uint64(), err))

}

func TestStakingPlugin_CreateCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}

	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	/**
	Start Create Staking
	*/
	err = create_staking(state, blockNumber, blockHash, 1, 0, t)
	assert.Nil(t, err, fmt.Sprintf("Failed to Create Staking: %v", err))
}

func TestStakingPlugin_GetCandidateInfo(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/**
	Start Get Candidate Info
	*/
	can, err := getCandidate(blockHash, index)
	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

}

func TestStakingPlugin_GetCandidateInfoByIrr(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/**
	Start GetCandidateInfoByIrr

	Get Candidate Info
	*/
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	can, err := StakingInstance().GetCandidateInfoByIrr(addr)

	assert.Nil(t, err, fmt.Sprintf("Failed to GetCandidateInfoByIrr: %v", err))
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

}

func TestStakingPlugin_GetCandidateList(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	count := 0
	for i := 0; i < 4; i++ {
		if err := create_staking(state, blockNumber, blockHash, i, 0, t); nil != err {
			t.Error("Failed to Create num: "+fmt.Sprint(i)+" Staking", err)
			return
		}
		count++
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/**
	Start GetCandidateList
	*/

	queue, err := StakingInstance().GetCandidateList(blockHash, blockNumber.Uint64())
	assert.Nil(t, err, fmt.Sprintf("Failed to GetCandidateList: %v", err))
	assert.Equal(t, count, len(queue))
	queueByte, _ := json.Marshal(queue)
	t.Log("Get CandidateList Info is:", string(queueByte))
}

func TestStakingPlugin_EditorCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	var c *staking.Candidate
	// Get Candidate Info
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
		return
	}

	/**
	Start Edit Candidate
	*/
	c.NodeName = nodeNameArr[index+1]
	c.ExternalId = "What is this ?"
	c.Website = "www.baidu.com"
	c.Details = "This is buidu website ?"

	canAddr, _ := xutil.NodeId2Addr(c.NodeId)

	err = StakingInstance().EditCandidate(blockHash2, blockNumber2, canAddr, c)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to EditCandidate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get Candidate info after edit
	can, err := getCandidate(blockHash2, index)

	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

}

func TestStakingPlugin_IncreaseStaking(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	var c *staking.Candidate
	// Get Candidate Info
	if can, err := getCandidate(blockHash, index); nil != err {
		t.Errorf("Failed to Get candidate info, err: %v", err)
		return
	} else {
		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
		return
	}

	/**
	Start IncreaseStaking
	*/
	canAddr, _ := xutil.NodeId2Addr(c.NodeId)
	err = StakingInstance().IncreaseStaking(state, blockHash2, blockNumber2, common.Big256, uint16(0), canAddr, c)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to IncreaseStaking: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get Candidate info
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])
	can, err := StakingInstance().GetCandidateInfoByIrr(addr)

	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

}

func TestStakingPlugin_WithdrewCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	// Get Candidate Info
	can, err := getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
		return
	}

	/**
	Start WithdrewStaking
	*/
	canAddr, _ := xutil.NodeId2Addr(can.NodeId)
	err = StakingInstance().WithdrewStaking(state, blockHash2, blockNumber2, canAddr, can)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewStaking: %v", err)) {
		return
	}

	t.Log("Finish WithdrewStaking ~~")
	// get Candidate info
	if _, err := getCandidate(blockHash2, index); snapshotdb.IsDbNotFoundErr(err) {
		t.Logf("expect candidate info is no found, err: %v", err)
		return
	} else {
		t.Error("It is not expect~")
		return
	}

}

func TestStakingPlugin_HandleUnCandidateItem(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	// Add UNStakingItems
	//stakingDB := staking.NewStakingDB()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())
	canAddr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	if err := StakingInstance().addUnStakeItem(state, blockNumber.Uint64(), blockHash, epoch, nodeIdArr[index], canAddr, blockNumber.Uint64()); nil != err {
		t.Error("Failed to AddUnStakeItemStore:", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	// Get Candidate Info
	can, err := getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
		return
	}

	// get Candidate info
	if _, err := getCandidate(blockHash2, index); snapshotdb.IsDbNotFoundErr(err) {
		t.Fatal(fmt.Sprintf("expect candidate info is no found, err: %v", err))
	}
	/**
	Start HandleUnCandidateItem
	*/
	err = StakingInstance().HandleUnCandidateItem(state, blockNumber2.Uint64(), blockHash2, epoch+xcom.UnStakeFreezeDuration())

	if !assert.Nil(t, err, fmt.Sprintf("Failed to HandleUnCandidateItem: %v", err)) {
		return
	}

	// get Candidate info
	_, err = getCandidate(blockHash2, index)
	assert.True(t, snapshotdb.IsDbNotFoundErr(err))

	// Penalty for simulating low block rate, and then return to normal state
	index++
	if err := create_staking(state, blockNumber2, blockHash2, index, 0, t); nil != err {
		t.Fatal(err)
	}
	canAddr, _ = xutil.NodeId2Addr(nodeIdArr[index])
	if err := StakingInstance().addRecoveryUnStakeItem(blockNumber2.Uint64(), blockHash2, nodeIdArr[index], canAddr, blockNumber2.Uint64()); nil != err {
		t.Error("Failed to AddUnStakeItemStore:", err)
		return
	}
	epoch = xutil.CalculateEpoch(blockNumber2.Uint64())
	err = StakingInstance().HandleUnCandidateItem(state, blockNumber2.Uint64(), blockHash2, epoch+xcom.ZeroProduceFreezeDuration())
	assert.Nil(t, err)

	recoveryCan, err := getCandidate(blockHash2, index)
	assert.Nil(t, err)
	assert.NotNil(t, recoveryCan)
	assert.True(t, recoveryCan.IsValid())

	// The simulation first punishes the low block rate, and then the double sign punishment.
	// After the lock-up period of the low block rate penalty expires, the double-signing pledge freeze
	index++
	if err := create_staking(state, blockNumber2, blockHash2, index, 0, t); nil != err {
		t.Fatal(err)
	}
	canAddr, _ = xutil.NodeId2Addr(nodeIdArr[index])
	recoveryCan2, err := getCandidate(blockHash2, index)
	assert.Nil(t, err)
	recoveryCan2.AppendStatus(staking.Invalided)
	recoveryCan2.AppendStatus(staking.LowRatio)
	recoveryCan2.AppendStatus(staking.DuplicateSign)
	assert.True(t, recoveryCan2.IsInvalidLowRatio())
	assert.Nil(t, StakingInstance().EditCandidate(blockHash2, blockNumber2, canAddr, recoveryCan2))

	// Handle the lock period of low block rate, and increase the double sign freeze operation
	if err := StakingInstance().addRecoveryUnStakeItem(blockNumber2.Uint64(), blockHash2, nodeIdArr[index], canAddr, blockNumber2.Uint64()); nil != err {
		t.Error("Failed to AddUnStakeItemStore:", err)
		return
	}
	newBlockNumber := new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch()*xcom.ZeroProduceFreezeDuration() + blockNumber2.Uint64())
	epoch = xutil.CalculateEpoch(newBlockNumber.Uint64())
	err = StakingInstance().HandleUnCandidateItem(state, newBlockNumber.Uint64(), blockHash2, epoch)
	assert.Nil(t, err)

	recoveryCan2, err = getCandidate(blockHash2, index)
	assert.Nil(t, err)

	assert.NotNil(t, recoveryCan2)
	assert.True(t, recoveryCan2.IsInvalidDuplicateSign())
	assert.False(t, recoveryCan2.IsInvalidLowRatio())

	// Handle double-signature freeze and release pledge, delete nodes
	newBlockNumber.Add(newBlockNumber, new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch()*xcom.UnStakeFreezeDuration()))
	err = StakingInstance().HandleUnCandidateItem(state, newBlockNumber.Uint64(), blockHash2, xcom.UnStakeFreezeDuration()+epoch)
	assert.Nil(t, err)
	recoveryCan2, err = getCandidate(blockHash2, index)
	assert.True(t, snapshotdb.IsDbNotFoundErr(err))
}

func TestStakingPlugin_Delegate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()

	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	// Get Candidate Info
	can, err := getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start Delegate
	*/
	del, err := delegate(state, blockHash2, blockNumber2, can, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
		return
	}
	t.Log("Finish Delegate ~~, Info is:", del)
	can, err = getCandidate(blockHash2, index)

	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	assert.True(t, can.DelegateTotalHes.Cmp(del.ReleasedHes) == 0)
	assert.True(t, can.DelegateEpoch == del.DelegateEpoch)
	assert.True(t, del.CumulativeIncome.Cmp(common.Big0) == 0)

	delegateRewardPerList := make([]*reward.DelegateRewardPer, 0)
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    1,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    2,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})

	canAddr, _ := xutil.NodeId2Addr(can.NodeId)

	curBlockNumber := new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch() * 3)
	if err := sndb.NewBlock(curBlockNumber, blockHash2, blockHash3); nil != err {
		t.Error("newBlock 3 err", err)
		return
	}

	expectedCumulativeIncome := delegateRewardPerList[1].CalDelegateReward(del.ReleasedHes)
	delegateAmount := new(big.Int).Mul(new(big.Int).SetInt64(10), new(big.Int).SetInt64(params.LAT))
	if err := StakingInstance().Delegate(state, blockHash3, curBlockNumber, addrArr[index+1], del, canAddr, can, 0, delegateAmount, delegateRewardPerList); nil != err {
		t.Fatal("Failed to Delegate:", err)
	}

	assert.True(t, del.CumulativeIncome.Cmp(expectedCumulativeIncome) == 0)
	assert.True(t, del.DelegateEpoch == 3)
	assert.True(t, del.ReleasedHes.Cmp(delegateAmount) == 0)
	assert.True(t, can.DelegateEpoch == del.DelegateEpoch)
	assert.True(t, can.DelegateTotal.Cmp(del.Released) == 0)

	t.Log("Finish Delegate ~~, Info is:", del)

	t.Log("Get Candidate Info is:", can)

}

func TestStakingPlugin_WithdrewDelegation(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	// Get Candidate Info
	can, err := getCandidate(blockHash, index)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	// Delegate
	del, err := delegate(state, blockHash, blockNumber, can, 0, index, t)

	delegateRewardPoolBalance, _ := new(big.Int).SetString(balanceStr[index+1], 10) // PASS
	state.AddBalance(vm.DelegateRewardPoolAddr, new(big.Int).Mul(new(big.Int).Set(delegateRewardPoolBalance), new(big.Int).Set(delegateRewardPoolBalance)))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to delegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	t.Log("Finish delegate ~~")
	can, err = getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start Withdrew Delegate
	*/
	amount := common.Big257
	delegateTotalHes := can.DelegateTotalHes
	err = StakingInstance().WithdrewDelegation(state, blockHash2, blockNumber2, amount, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, make([]*reward.DelegateRewardPer, 0))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}
	t.Log("Finish WithdrewDelegate ~~", del)
	can, err = getCandidate(blockHash2, index)

	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	assert.True(t, new(big.Int).Sub(delegateTotalHes, amount).Cmp(can.DelegateTotalHes) == 0)
	assert.True(t, new(big.Int).Sub(delegateTotalHes, amount).Cmp(del.ReleasedHes) == 0)
	assert.True(t, del.WithdrewEpoch == 0)
	assert.True(t, del.WithdrewAmount.Cmp(common.Big0) == 0)
	assert.True(t, del.UnLockEpoch == 0)
	assert.True(t, del.CumulativeIncome.Cmp(common.Big0) == 0)

	curBlockNumber := new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch() * 3)
	if err := sndb.NewBlock(curBlockNumber, blockHash2, blockHash3); nil != err {
		t.Error("newBlock 3 err", err)
		return
	}

	delegateRewardPerList := make([]*reward.DelegateRewardPer, 0)
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    1,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[0], sndb); nil != err {
		t.Fatal(err)
	}
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    2,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[1], sndb); nil != err {
		t.Fatal(err)
	}

	expectedIssueIncome := delegateRewardPerList[1].CalDelegateReward(del.ReleasedHes)
	err = StakingInstance().WithdrewDelegation(state, blockHash3, curBlockNumber, del.ReleasedHes, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, delegateRewardPerList)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	can, err = getCandidate(blockHash3, index)

	assert.True(t, expectedIssueIncome.Cmp(del.CumulativeIncome) == 0)
	assert.True(t, del.WithdrewEpoch > 0)
	assert.True(t, del.WithdrewAmount.Cmp(common.Big0) > 0)
	assert.True(t, del.UnLockEpoch > 0)
	assert.True(t, del.ReleasedHes.Cmp(common.Big0) == 0)
	assert.True(t, new(big.Int).Sub(del.Released, del.WithdrewAmount).Cmp(common.Big0) == 0)
	t.Log("Get Candidate Info is:", can)
}

func TestStakingPlugin_WithdrewDelegation_AllHes(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	// Get Candidate Info
	can, err := getCandidate(blockHash, index)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	// Delegate
	del, err := delegate(state, blockHash, blockNumber, can, 0, index, t)

	delegateRewardPoolBalance, _ := new(big.Int).SetString(balanceStr[index+1], 10) // PASS
	state.AddBalance(vm.DelegateRewardPoolAddr, new(big.Int).Mul(new(big.Int).Set(delegateRewardPoolBalance), new(big.Int).Set(delegateRewardPoolBalance)))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to delegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	t.Log("Finish delegate ~~")
	can, err = getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start Withdrew Delegate
	*/
	err = StakingInstance().WithdrewDelegation(state, blockHash2, blockNumber2, can.DelegateTotalHes, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, make([]*reward.DelegateRewardPer, 0))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}
	can, err = getCandidate(blockHash2, index)

	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	newDel := getDelegate(blockHash2, blockNumber.Uint64(), index, t)
	assert.True(t, newDel == nil)
}

func TestStakingPlugin_WithdrewDelegation_SlashCandidates(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	// Get Candidate Info
	can, err := getCandidate(blockHash, index)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	// Delegate
	del, err := delegate(state, blockHash, blockNumber, can, 0, index, t)

	delegateRewardPoolBalance, _ := new(big.Int).SetString(balanceStr[index+1], 10) // PASS
	state.AddBalance(vm.DelegateRewardPoolAddr, new(big.Int).Mul(new(big.Int).Set(delegateRewardPoolBalance), new(big.Int).Set(delegateRewardPoolBalance)))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to delegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	can, err = getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	curBlockNumber := new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch() * 2)
	if err := sndb.NewBlock(curBlockNumber, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}
	slashItem := &staking.SlashNodeItem{NodeId: can.NodeId, Amount: new(big.Int).Set(common.Big1), SlashType: staking.LowRatio, BenefitAddr: can.StakingAddress}
	StakingInstance().SlashCandidates(state, blockHash2, curBlockNumber.Uint64(), slashItem)

	/**
	Start Withdrew Delegate
	*/
	err = StakingInstance().WithdrewDelegation(state, blockHash2, curBlockNumber, new(big.Int).Set(common.Big1), addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, make([]*reward.DelegateRewardPer, 0))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	can, err = getCandidate(blockHash2, index)

	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	newDel := getDelegate(blockHash2, blockNumber.Uint64(), index, t)
	assert.True(t, newDel != nil)
	assert.True(t, new(big.Int).Sub(newDel.Released, common.Big1).Cmp(new(big.Int).Sub(can.Shares, can.Released)) == 0)
	assert.True(t, can.DelegateTotal.Cmp(new(big.Int).Sub(newDel.Released, common.Big1)) == 0)
	assert.True(t, can.WithdrewDelegateAmount.Cmp(del.WithdrewAmount) == 0)
}

func TestStakingPlugin_WithdrewDelegation_Repeat(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	// Get Candidate Info
	can, err := getCandidate(blockHash, index)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	// Delegate
	del, err := delegate(state, blockHash, blockNumber, can, 0, index, t)

	delegateRewardPoolBalance, _ := new(big.Int).SetString(balanceStr[index+1], 10) // PASS
	state.AddBalance(vm.DelegateRewardPoolAddr, new(big.Int).Mul(new(big.Int).Set(delegateRewardPoolBalance), new(big.Int).Set(delegateRewardPoolBalance)))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to delegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	t.Log("Finish delegate ~~")
	can, err = getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start Withdrew Delegate
	*/
	amount := common.Big257
	delegateTotalHes := can.DelegateTotalHes
	err = StakingInstance().WithdrewDelegation(state, blockHash2, blockNumber2, amount, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, make([]*reward.DelegateRewardPer, 0))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}
	t.Log("Finish WithdrewDelegate ~~", del)
	can, err = getCandidate(blockHash2, index)

	newDel := getDelegate(blockHash2, blockNumber.Uint64(), index, t)
	assert.True(t, newDel != nil)
	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	assert.True(t, new(big.Int).Sub(delegateTotalHes, amount).Cmp(can.DelegateTotalHes) == 0)
	assert.True(t, new(big.Int).Sub(delegateTotalHes, amount).Cmp(del.ReleasedHes) == 0)
	assert.True(t, del.WithdrewEpoch == 0)
	assert.True(t, del.WithdrewAmount.Cmp(common.Big0) == 0)
	assert.True(t, del.UnLockEpoch == 0)
	assert.True(t, del.CumulativeIncome.Cmp(common.Big0) == 0)

	curBlockNumber := new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch() * 3)
	if err := sndb.NewBlock(curBlockNumber, blockHash2, blockHash3); nil != err {
		t.Error("newBlock 3 err", err)
		return
	}

	delegateRewardPerList := make([]*reward.DelegateRewardPer, 0)
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    1,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[0], sndb); nil != err {
		t.Fatal(err)
	}
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    2,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[1], sndb); nil != err {
		t.Fatal(err)
	}

	expectedIssueIncome := delegateRewardPerList[1].CalDelegateReward(del.ReleasedHes)
	err = StakingInstance().WithdrewDelegation(state, blockHash3, curBlockNumber, del.ReleasedHes, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, delegateRewardPerList)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	can, err = getCandidate(blockHash3, index)

	assert.True(t, expectedIssueIncome.Cmp(del.CumulativeIncome) == 0)
	assert.True(t, del.WithdrewEpoch > 0)
	assert.True(t, del.WithdrewAmount.Cmp(common.Big0) > 0)
	assert.True(t, del.UnLockEpoch > 0)
	assert.True(t, del.ReleasedHes.Cmp(common.Big0) == 0)
	assert.True(t, new(big.Int).Sub(del.Released, del.WithdrewAmount).Cmp(common.Big0) == 0)
	t.Log("Get Candidate Info is:", can)

	err = StakingInstance().WithdrewDelegation(state, blockHash3, curBlockNumber, del.ReleasedHes, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, delegateRewardPerList)
	if err == nil {
		t.Fatal("Failed to execute WithdrewDelegation")
	}
	berr, _ := err.(*common.BizError)
	assert.True(t, berr.Code == staking.ErrAlreadyWithdrewDelegation.Code)

	canAddr, _ := xutil.NodeId2Addr(can.NodeId)
	if err := StakingInstance().Delegate(state, blockHash3, curBlockNumber, addrArr[index+1], del, canAddr, can, 0, amount, nil); err != nil {
		t.Fatal(err)
	}
	oldWithdrewAmount := del.WithdrewAmount
	err = StakingInstance().WithdrewDelegation(state, blockHash3, curBlockNumber, del.ReleasedHes, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, nil)
	if err != nil {
		t.Fatal("Failed to execute WithdrewDelegation")
	}
	assert.True(t, del.ReleasedHes.Cmp(common.Big0) == 0)
	assert.True(t, oldWithdrewAmount.Cmp(del.WithdrewAmount) == 0)

	if err := StakingInstance().Delegate(state, blockHash3, curBlockNumber, addrArr[index+1], del, canAddr, can, 0, amount, nil); err != nil {
		t.Fatal(err)
	}
	err = StakingInstance().WithdrewDelegation(state, blockHash3, curBlockNumber, new(big.Int).Mul(amount, common.Big2), addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, delegateRewardPerList)
	if err == nil {
		t.Fatal("Failed to execute WithdrewDelegation")
	}

	berr, _ = err.(*common.BizError)
	assert.True(t, berr.Code == staking.ErrHesNotEnough.Code)
}

func TestStakingPlugin_WithdrewDelegateAndRedeemDelegationAll(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	// Get Candidate Info
	can, err := getCandidate(blockHash, index)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	// Delegate
	del, err := delegate(state, blockHash, blockNumber, can, 0, index, t)

	delegateRewardPoolBalance, _ := new(big.Int).SetString(balanceStr[index+1], 10) // PASS
	state.AddBalance(vm.DelegateRewardPoolAddr, new(big.Int).Mul(new(big.Int).Set(delegateRewardPoolBalance), new(big.Int).Set(delegateRewardPoolBalance)))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to delegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	t.Log("Finish delegate ~~")
	can, err = getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start Withdrew Delegate
	*/
	amount := common.Big257
	delegateTotalHes := can.DelegateTotalHes
	err = StakingInstance().WithdrewDelegation(state, blockHash2, blockNumber2, amount, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, make([]*reward.DelegateRewardPer, 0))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}
	t.Log("Finish WithdrewDelegate ~~", del)
	can, err = getCandidate(blockHash2, index)

	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	assert.True(t, new(big.Int).Sub(delegateTotalHes, amount).Cmp(can.DelegateTotalHes) == 0)
	assert.True(t, new(big.Int).Sub(delegateTotalHes, amount).Cmp(del.ReleasedHes) == 0)
	assert.True(t, del.WithdrewEpoch == 0)
	assert.True(t, del.WithdrewAmount.Cmp(common.Big0) == 0)
	assert.True(t, del.UnLockEpoch == 0)
	assert.True(t, del.CumulativeIncome.Cmp(common.Big0) == 0)

	curBlockNumber := new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch() * 3)
	if err := sndb.NewBlock(curBlockNumber, blockHash2, blockHash3); nil != err {
		t.Error("newBlock 3 err", err)
		return
	}

	delegateRewardPerList := make([]*reward.DelegateRewardPer, 0)
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    1,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[0], sndb); nil != err {
		t.Fatal(err)
	}
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    2,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[1], sndb); nil != err {
		t.Fatal(err)
	}

	expectedIssueIncome := delegateRewardPerList[1].CalDelegateReward(del.ReleasedHes)
	expectedBalance := new(big.Int).Add(state.GetBalance(addrArr[index+1]), expectedIssueIncome)
	expectedBalance = new(big.Int).Add(expectedBalance, del.ReleasedHes)
	err = StakingInstance().WithdrewDelegation(state, blockHash3, curBlockNumber, del.ReleasedHes, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, delegateRewardPerList)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	can, err = getCandidate(blockHash3, index)
	assert.True(t, expectedIssueIncome.Cmp(del.CumulativeIncome) == 0)
	assert.True(t, del.WithdrewEpoch > 0)
	assert.True(t, del.WithdrewAmount.Cmp(common.Big0) > 0)
	assert.True(t, del.UnLockEpoch > 0)
	assert.True(t, del.ReleasedHes.Cmp(common.Big0) == 0)
	assert.True(t, new(big.Int).Sub(del.Released, del.WithdrewAmount).Cmp(common.Big0) == 0)
	t.Log("Get Candidate Info is:", can)

	_, err = StakingInstance().RedeemDelegation(state, blockHash3, curBlockNumber, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, nil)
	if err == nil {
		t.Fatal("Failed to execute Redeem Delegation")
	}
	berr, _ := err.(*common.BizError)
	assert.True(t, berr.Code == staking.ErrWithdrewDelegationLocking.Code)

	duration, err := gov.GovernUnDelegateFreezeDuration(blockNumber.Uint64(), blockHash)
	if nil != err {
		t.Fatal(err)
	}
	curBlockNumber = new(big.Int).Add(curBlockNumber, new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch()*duration+1))

	income, err := StakingInstance().RedeemDelegation(state, blockHash3, curBlockNumber, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, nil)
	if err != nil {
		t.Fatal("Failed to execute Redeem Delegation")
	}

	can, err = getCandidate(blockHash3, index)
	assert.True(t, expectedBalance.Cmp(state.GetBalance(addrArr[index+1])) == 0)
	assert.True(t, expectedIssueIncome.Cmp(income) == 0)
	newDel := getDelegate(blockHash3, blockNumber.Uint64(), index, t)
	assert.True(t, newDel == nil)
}

func TestStakingPlugin_WithdrewDelegateAndRedeemDelegationPart(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	// Get Candidate Info
	can, err := getCandidate(blockHash, index)
	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	// Delegate
	del, err := delegate(state, blockHash, blockNumber, can, 0, index, t)

	delegateRewardPoolBalance, _ := new(big.Int).SetString(balanceStr[index+1], 10) // PASS
	state.AddBalance(vm.DelegateRewardPoolAddr, new(big.Int).Mul(new(big.Int).Set(delegateRewardPoolBalance), new(big.Int).Set(delegateRewardPoolBalance)))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to delegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	t.Log("Finish delegate ~~")
	can, err = getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start Withdrew Delegate
	*/
	amount := common.Big257
	delegateTotalHes := can.DelegateTotalHes
	err = StakingInstance().WithdrewDelegation(state, blockHash2, blockNumber2, amount, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, make([]*reward.DelegateRewardPer, 0))

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}
	t.Log("Finish WithdrewDelegate ~~", del)
	can, err = getCandidate(blockHash2, index)

	assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err))
	assert.True(t, nil != can)
	assert.True(t, new(big.Int).Sub(delegateTotalHes, amount).Cmp(can.DelegateTotalHes) == 0)
	assert.True(t, new(big.Int).Sub(delegateTotalHes, amount).Cmp(del.ReleasedHes) == 0)
	assert.True(t, del.WithdrewEpoch == 0)
	assert.True(t, del.WithdrewAmount.Cmp(common.Big0) == 0)
	assert.True(t, del.UnLockEpoch == 0)
	assert.True(t, del.CumulativeIncome.Cmp(common.Big0) == 0)

	curBlockNumber := new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch() * 3)
	if err := sndb.NewBlock(curBlockNumber, blockHash2, blockHash3); nil != err {
		t.Error("newBlock 3 err", err)
		return
	}
	_, err = StakingInstance().RedeemDelegation(state, blockHash3, curBlockNumber, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, nil)
	if err == nil {
		t.Fatal("Failed to execute Redeem Delegation")
	}
	berr, _ := err.(*common.BizError)
	assert.True(t, berr.Code == staking.ErrNotRedeemDelegation.Code)

	delegateRewardPerList := make([]*reward.DelegateRewardPer, 0)
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    1,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[0], sndb); nil != err {
		t.Fatal(err)
	}
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    2,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[1], sndb); nil != err {
		t.Fatal(err)
	}

	oldReleased := new(big.Int).Set(del.ReleasedHes)
	expectedIssueIncome := delegateRewardPerList[1].CalDelegateReward(del.ReleasedHes)
	// Revocation of part of the delegation
	err = StakingInstance().WithdrewDelegation(state, blockHash3, curBlockNumber, amount, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, delegateRewardPerList)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	can, err = getCandidate(blockHash3, index)
	assert.True(t, expectedIssueIncome.Cmp(del.CumulativeIncome) == 0)
	assert.True(t, del.WithdrewEpoch > 0)
	assert.True(t, del.WithdrewAmount.Cmp(common.Big0) > 0)
	assert.True(t, del.UnLockEpoch > 0)
	assert.True(t, del.ReleasedHes.Cmp(common.Big0) == 0)
	assert.True(t, del.Released.Cmp(oldReleased) == 0)
	assert.True(t, del.WithdrewAmount.Cmp(amount) == 0)
	t.Log("Get Candidate Info is:", can)

	duration, err := gov.GovernUnDelegateFreezeDuration(blockNumber.Uint64(), blockHash)
	if nil != err {
		t.Fatal(err)
	}
	curBlockNumber = new(big.Int).Add(curBlockNumber, new(big.Int).SetUint64(xutil.CalcBlocksEachEpoch()*duration+1))

	delegateRewardPerList = make([]*reward.DelegateRewardPer, 0)
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    3,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[0], sndb); nil != err {
		t.Fatal(err)
	}
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    4,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[1], sndb); nil != err {
		t.Fatal(err)
	}
	delegateRewardPerList = append(delegateRewardPerList, &reward.DelegateRewardPer{
		Epoch:    5,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	if err := AppendDelegateRewardPer(blockHash3, can.NodeId, can.StakingBlockNum, delegateRewardPerList[2], sndb); nil != err {
		t.Fatal(err)
	}

	expectedIssueIncome = new(big.Int).Add(expectedIssueIncome, delegateRewardPerList[0].CalDelegateReward(del.Released))
	expectedIssueIncome = new(big.Int).Add(expectedIssueIncome, delegateRewardPerList[1].CalDelegateReward(new(big.Int).Sub(del.Released, del.WithdrewAmount)))
	expectedIssueIncome = new(big.Int).Add(expectedIssueIncome, delegateRewardPerList[2].CalDelegateReward(new(big.Int).Sub(del.Released, del.WithdrewAmount)))
	oldReleased = new(big.Int).Set(del.Released)
	income, err := StakingInstance().RedeemDelegation(state, blockHash3, curBlockNumber, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, delegateRewardPerList)
	if err != nil {
		t.Fatal("Failed to execute Redeem Delegation")
	}

	can, err = getCandidate(blockHash3, index)
	assert.True(t, income.Cmp(common.Big0) == 0)
	newDel := getDelegate(blockHash3, blockNumber.Uint64(), index, t)
	assert.True(t, newDel != nil)
	assert.True(t, newDel.CumulativeIncome.Cmp(expectedIssueIncome) == 0)
	assert.True(t, newDel.WithdrewAmount.Cmp(common.Big0) == 0)
	assert.True(t, newDel.WithdrewEpoch == 0)
	assert.True(t, newDel.UnLockEpoch == 0)
	assert.True(t, new(big.Int).Sub(oldReleased, newDel.Released).Cmp(amount) == 0)
}

func TestStakingPlugin_GetDelegateInfo(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	t.Log("Finish delegate ~~")

	can, err := getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	// Delegate
	_, err = delegate(state, blockHash2, blockNumber2, can, 0, index, t)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to delegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
		return
	}

	t.Log("Finished Delegate ~~")
	// get Delegate info
	del := getDelegate(blockHash2, blockNumber.Uint64(), index, t)
	assert.True(t, nil != del)
	t.Log("Get Delegate Info is:", del)
}

func TestStakingPlugin_GetDelegateInfoByIrr(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	can, err := getCandidate(blockHash, index)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to getCandidate: %v", err)) {
		return
	}
	assert.True(t, nil != can)
	t.Log("Get Candidate Info is:", can)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	t.Log("Start delegate ~~")
	// Delegate
	_, err = delegate(state, blockHash2, blockNumber2, can, 0, index, t)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to delegate: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
		return
	}

	t.Log("Finished Delegate ~~")
	/**
	Start get Delegate info
	*/
	del, err := StakingInstance().GetDelegateInfoByIrr(addrArr[index+1], nodeIdArr[index], blockNumber.Uint64())
	if nil != err {
		t.Error("Failed to GetDelegateInfoByIrr:", err)
		return
	}

	assert.Nil(t, err, fmt.Sprintf("Failed to GetDelegateInfoByIrr: %v", err))
	assert.True(t, nil != del)
	t.Log("Get Delegate Info is:", del)

}

func TestStakingPlugin_GetRelatedListByDelAddr(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()
	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	// staking 0, 1, 2, 3
	for i := 0; i < 4; i++ {
		if err := create_staking(state, blockNumber, blockHash, i, 0, t); nil != err {
			t.Error("Failed to Create Staking", err)
			return
		}
	}

	t.Log("First delegate ~~")
	for i := 0; i < 2; i++ {
		// 0, 1
		var c *staking.Candidate

		if can, err := getCandidate(blockHash, i); nil != err {
			t.Errorf("Failed to Get candidate info, err: %v", err)
			return
		} else {
			canByte, _ := json.Marshal(can)
			t.Log("Get Candidate Info is:", string(canByte))
			c = can
		}
		// Delegate  0, 1
		_, err := delegate(state, blockHash, blockNumber, c, 0, i, t)
		if nil != err {
			t.Errorf("Failed to Delegate: Num: %d, error: %v", i, err)
			return
		}
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	t.Log("Second delegate ~~")
	for i := 1; i < 3; i++ {
		// 0, 1
		var c *staking.Candidate
		if can, err := getCandidate(blockHash2, i-1); nil != err {
			t.Errorf("Failed to Get candidate info, err: %v", err)
			return
		} else {
			canByte, _ := json.Marshal(can)
			t.Log("Get Candidate Info is:", string(canByte))
			c = can
		}

		// Delegate
		_, err := delegate(state, blockHash2, blockNumber2, c, 0, i, t)
		if nil != err {
			t.Errorf("Failed to Delegate: Num: %d, error: %v", i, err)
			return
		}
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
		return
	}

	t.Log("Finished Delegate ~~")
	/**
	Start get RelatedList
	*/
	rel, err := StakingInstance().GetRelatedListByDelAddr(blockHash2, addrArr[1+1])

	assert.Nil(t, err, fmt.Sprintf("Failed to GetRelatedListByDelAddr: %v", err))
	assert.True(t, nil != rel)
	t.Log("Get RelateList Info is:", rel)
}

func TestStakingPlugin_ElectNextVerifierList(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		canTmp := &staking.Candidate{
			CandidateBase: &staking.CandidateBase{
				NodeId:          nodeId,
				BlsPubKey:       blsKeyHex,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(i),
				StakingTxIndex:  uint32(index),
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
				},
			},
			CandidateMutable: &staking.CandidateMutable{
				Shares: balance,

				// Prevent null pointer initialization
				Released:           common.Big0,
				ReleasedHes:        common.Big0,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList
	start := uint64(1)
	end := xutil.EpochSize() * xutil.ConsensusSize()

	new_verifierArr := &staking.ValidatorArray{
		Start: start,
		End:   end,
	}

	queue := make(staking.ValidatorQueue, 0)

	iter := sndb.Ranking(blockHash, staking.CanPowerKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		t.Errorf("Failed to build genesis VerifierList, the iter is  err: %v", err)
		return
	}

	defer iter.Release()

	// for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {

	count := 0
	for iter.Valid(); iter.Next(); {
		if uint64(count) == xcom.MaxValidators() {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToNodeAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToNodeAddress(addrSuffix)

		val := &staking.Validator{
			NodeAddress:     addr,
			NodeId:          can.NodeId,
			BlsPubKey:       can.BlsPubKey,
			ProgramVersion:  can.ProgramVersion,
			Shares:          can.Shares,
			StakingBlockNum: can.StakingBlockNum,
			StakingTxIndex:  can.StakingTxIndex,
			ValidatorTerm:   0,
		}
		queue = append(queue, val)
		count++
	}

	new_verifierArr.Arr = queue

	err = setVerifierList(blockHash, new_verifierArr)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to VerifierList: %v", err)) {
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/*
		Start ElectNextVerifierList
	*/
	targetNum := xutil.EpochSize() * xutil.ConsensusSize()

	targetNumInt := big.NewInt(int64(targetNum))

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	err = StakingInstance().ElectNextVerifierList(blockHash2, targetNumInt.Uint64(), state)

	assert.Nil(t, err, fmt.Sprintf("Failed to ElectNextVerifierList: %v", err))

}

func TestStakingPlugin_Election(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// Must new VrfHandler instance by genesis block Hash
	handler.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		canTmp := &staking.Candidate{
			CandidateBase: &staking.CandidateBase{
				NodeId:          nodeId,
				BlsPubKey:       blsKeyHex,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(i),
				StakingTxIndex:  uint32(index),
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
				},
			},
			CandidateMutable: &staking.CandidateMutable{
				Shares: balance,

				// Prevent null pointer initialization
				Released:           common.Big0,
				ReleasedHes:        common.Big0,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xutil.EpochSize() * xutil.ConsensusSize()

	new_verifierArr := &staking.ValidatorArray{
		Start: start,
		End:   end,
	}

	queue := make(staking.ValidatorQueue, 0)

	iter := sndb.Ranking(blockHash, staking.CanPowerKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		t.Errorf("Failed to build genesis VerifierList, the iter is  err: %v", err)
		return
	}

	defer iter.Release()

	count := 0
	for iter.Valid(); iter.Next(); {
		if uint64(count) == xcom.MaxValidators() {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToNodeAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToNodeAddress(addrSuffix)

		val := &staking.Validator{
			NodeAddress:     addr,
			NodeId:          can.NodeId,
			BlsPubKey:       can.BlsPubKey,
			ProgramVersion:  can.ProgramVersion,
			Shares:          can.Shares,
			StakingBlockNum: can.StakingBlockNum,
			StakingTxIndex:  can.StakingTxIndex,
			ValidatorTerm:   0,
		}
		queue = append(queue, val)
		count++
	}

	new_verifierArr.Arr = queue

	err = setVerifierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	// build gensis current validatorList
	new_validatorArr := &staking.ValidatorArray{
		Start: start,
		End:   xutil.ConsensusSize(),
	}

	new_validatorArr.Arr = queue[:int(xcom.MaxConsensusVals())]

	err = setRoundValList(blockHash, new_validatorArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis current round validatorList, err: %v", err)
		return
	}

	// build ancestor nonces
	currNonce, nonces := build_vrf_Nonce()
	if enValue, err := rlp.EncodeToBytes(nonces); nil != err {
		t.Error("Storage previous nonce failed", "err", err)
		return
	} else {
		sndb.Put(blockHash, handler.NonceStorageKey, enValue)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/*
		Start Election
	*/
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	header := &types.Header{
		ParentHash: blockHash,
		Number:     big.NewInt(int64(xutil.ConsensusSize() - xcom.ElectionDistance())),
		Nonce:      types.EncodeNonce(currNonce),
	}

	err = StakingInstance().Election(blockHash2, header, state)

	assert.Nil(t, err, fmt.Sprintf("Failed to Election: %v", err))

}

func TestStakingPlugin_SlashCandidates(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	// Will be Slashing candidate
	slashQueue := make(staking.CandidateQueue, 5)

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		releasedHes := new(big.Int).SetUint64(10000)
		canTmp := &staking.Candidate{
			CandidateBase: &staking.CandidateBase{
				NodeId:          nodeId,
				BlsPubKey:       blsKeyHex,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(i),
				StakingTxIndex:  uint32(index),
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
				},
			},
			CandidateMutable: &staking.CandidateMutable{
				Shares: new(big.Int).Add(balance, releasedHes),

				// Prevent null pointer initialization
				Released:           new(big.Int).Set(balance),
				ReleasedHes:        releasedHes,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, canTmp.Shares, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
		if i < len(slashQueue) {
			slashQueue[i] = canTmp
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xutil.EpochSize() * xutil.ConsensusSize()

	new_verifierArr := &staking.ValidatorArray{
		Start: start,
		End:   end,
	}

	queue := make(staking.ValidatorQueue, 0)

	iter := sndb.Ranking(blockHash, staking.CanPowerKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		t.Errorf("Failed to build genesis VerifierList, the iter is  err: %v", err)
		return
	}

	defer iter.Release()

	// for count := 0; iterator.Valid() && count < int(maxValidators); iterator.Next() {

	count := 0
	for iter.Valid(); iter.Next(); {
		if uint64(count) == xcom.MaxValidators() {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToNodeAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToNodeAddress(addrSuffix)

		val := &staking.Validator{
			NodeAddress:     addr,
			NodeId:          can.NodeId,
			BlsPubKey:       can.BlsPubKey,
			ProgramVersion:  can.ProgramVersion,
			Shares:          can.Shares,
			StakingBlockNum: can.StakingBlockNum,
			StakingTxIndex:  can.StakingTxIndex,
			ValidatorTerm:   0,
		}
		queue = append(queue, val)
		count++
	}

	new_verifierArr.Arr = queue

	err = setVerifierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	/**
	Start SlashCandidates
	*/
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock err", err)
		return
	}

	slash1 := slashQueue[0]
	slash2 := slashQueue[1]

	slashItemQueue := make(staking.SlashQueue, 0)

	// Be punished for less than the quality deposit
	slashItem1 := &staking.SlashNodeItem{
		NodeId:      slash1.NodeId,
		Amount:      slash1.Released,
		SlashType:   staking.LowRatio,
		BenefitAddr: vm.RewardManagerPoolAddr,
	}

	// Double sign penalty
	sla := new(big.Int).Div(slash2.Released, big.NewInt(10))
	caller := common.MustBech32ToAddress("lax1uj3zd9yz00axz7ls88ynwsp3jprhjd9ldx9qpm")
	slashItem2 := &staking.SlashNodeItem{
		NodeId:      slash2.NodeId,
		Amount:      sla,
		SlashType:   staking.DuplicateSign,
		BenefitAddr: caller,
	}
	slashItemQueue = append(slashItemQueue, slashItem1)
	slashItemQueue = append(slashItemQueue, slashItem2)

	// Penalty for two low block rates
	slash3 := slashQueue[2]
	slashAmount3 := new(big.Int).Div(slash3.Released, big.NewInt(10))
	slashItem3_1 := &staking.SlashNodeItem{
		NodeId:      slash3.NodeId,
		Amount:      slashAmount3,
		SlashType:   staking.LowRatio,
		BenefitAddr: vm.RewardManagerPoolAddr,
	}
	slashItem3_2 := &staking.SlashNodeItem{
		NodeId:      slash3.NodeId,
		Amount:      slashAmount3,
		SlashType:   staking.LowRatio,
		BenefitAddr: vm.RewardManagerPoolAddr,
	}
	slashItemQueue = append(slashItemQueue, slashItem3_1)
	slashItemQueue = append(slashItemQueue, slashItem3_2)

	// Penalty for low block rate first, and then trigger double sign penalty
	slash4 := slashQueue[3]
	slashAmount4 := new(big.Int).Div(slash4.Released, big.NewInt(10))
	slashItem4_1 := &staking.SlashNodeItem{
		NodeId:      slash4.NodeId,
		Amount:      slashAmount4,
		SlashType:   staking.LowRatio,
		BenefitAddr: vm.RewardManagerPoolAddr,
	}
	slashItem4_2 := &staking.SlashNodeItem{
		NodeId:      slash4.NodeId,
		Amount:      slashAmount4,
		SlashType:   staking.DuplicateSign,
		BenefitAddr: caller,
	}
	slashItemQueue = append(slashItemQueue, slashItem4_1)
	slashItemQueue = append(slashItemQueue, slashItem4_2)

	// Double signing penalty first, and then triggering low block rate penalty
	slash5 := slashQueue[4]
	slashAmount5 := new(big.Int).Div(slash5.Released, big.NewInt(10))
	slashItem5_1 := &staking.SlashNodeItem{
		NodeId:      slash5.NodeId,
		Amount:      slashAmount5,
		SlashType:   staking.DuplicateSign,
		BenefitAddr: caller,
	}
	slashItem5_2 := &staking.SlashNodeItem{
		NodeId:      slash5.NodeId,
		Amount:      slashAmount5,
		SlashType:   staking.LowRatio,
		BenefitAddr: vm.RewardManagerPoolAddr,
	}
	slashItemQueue = append(slashItemQueue, slashItem5_1)
	slashItemQueue = append(slashItemQueue, slashItem5_2)

	err = StakingInstance().SlashCandidates(state, blockHash2, blockNumber2.Uint64(), slashItemQueue...)
	assert.Nil(t, err, fmt.Sprintf("Failed to SlashCandidates Second can (DuplicateSign), err: %v", err))

	canAddr1, _ := xutil.NodeId2Addr(slash1.NodeId)
	can1, err := StakingInstance().GetCandidateInfo(blockHash2, canAddr1)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, can1.Released.Cmp(new(big.Int).Sub(slash1.Released, slashItem1.Amount)) == 0)
	assert.True(t, can1.ReleasedHes.Cmp(slash1.ReleasedHes) == 0)
	assert.True(t, can1.Shares.Cmp(common.Big0) > 0)

	canAddr2, _ := xutil.NodeId2Addr(slash2.NodeId)
	can2, err := StakingInstance().GetCandidateInfo(blockHash2, canAddr2)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, can2.Released.Cmp(new(big.Int).Sub(slash2.Released, slashItem2.Amount)) == 0)
	assert.True(t, can2.ReleasedHes.Cmp(common.Big0) == 0)
	assert.True(t, can2.Shares.Cmp(common.Big0) == 0)

	canAddr3, _ := xutil.NodeId2Addr(slash3.NodeId)
	can3, err := StakingInstance().GetCandidateInfo(blockHash2, canAddr3)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, can3.Released.Cmp(new(big.Int).Sub(slash3.Released, slashAmount3)) == 0)
	assert.True(t, can3.ReleasedHes.Cmp(slash3.ReleasedHes) == 0)
	assert.True(t, can3.Shares.Cmp(common.Big0) > 0)
	assert.True(t, can3.IsInvalidLowRatio())

	canAddr4, _ := xutil.NodeId2Addr(slash4.NodeId)
	can4, err := StakingInstance().GetCandidateInfo(blockHash2, canAddr4)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, can4.Released.Cmp(new(big.Int).Sub(slash4.Released, new(big.Int).Add(slashAmount4, slashAmount4))) == 0)
	assert.True(t, can4.ReleasedHes.Cmp(common.Big0) == 0)
	assert.True(t, can4.Shares.Cmp(common.Big0) == 0)
	assert.True(t, can4.IsInvalidLowRatio() && can4.IsInvalidDuplicateSign())

	canAddr5, _ := xutil.NodeId2Addr(slash5.NodeId)
	can5, err := StakingInstance().GetCandidateInfo(blockHash2, canAddr5)
	if err != nil {
		t.Fatal(err)
	}
	assert.True(t, can5.Released.Cmp(new(big.Int).Sub(slash5.Released, new(big.Int).Add(slashAmount5, slashAmount5))) == 0)
	assert.True(t, can5.ReleasedHes.Cmp(common.Big0) == 0)
	assert.True(t, can5.Shares.Cmp(common.Big0) == 0)
	assert.True(t, can5.IsInvalidLowRatio() && can5.IsInvalidDuplicateSign())
}

func TestStakingPlugin_DeclarePromoteNotify(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	handler.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	queue := make(staking.CandidateQueue, 0)
	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		canTmp := &staking.Candidate{
			CandidateBase: &staking.CandidateBase{
				NodeId:          nodeId,
				BlsPubKey:       blsKeyHex,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(i),
				StakingTxIndex:  uint32(index),
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
				},
			},
			CandidateMutable: &staking.CandidateMutable{
				Shares: balance,

				// Prevent null pointer initialization
				Released:           common.Big0,
				ReleasedHes:        common.Big0,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}

		if i < 20 {
			queue = append(queue, canTmp)
		}
	}

	// Commit Block 1
	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start DeclarePromoteNotify
	*/
	for i, can := range queue {
		err = StakingInstance().DeclarePromoteNotify(blockHash2, blockNumber2.Uint64(), can.NodeId, promoteVersion)

		assert.Nil(t, err, fmt.Sprintf("Failed to DeclarePromoteNotify, index: %d, err: %v", i, err))
	}

}

func TestStakingPlugin_ProposalPassedNotify(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}

	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	handler.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	validatorQueue := make(staking.ValidatorQueue, 0)

	nodeIdArr := make([]discover.NodeID, 0)
	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		canTmp := &staking.Candidate{
			CandidateBase: &staking.CandidateBase{
				NodeId:          nodeId,
				BlsPubKey:       blsKeyHex,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(i),
				StakingTxIndex:  uint32(index),
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
				},
			},
			CandidateMutable: &staking.CandidateMutable{
				Shares: balance,

				// Prevent null pointer initialization
				Released:           common.Big0,
				ReleasedHes:        common.Big0,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if !assert.Nil(t, err, fmt.Sprintf("Failed to Create Staking, num: %d, err: %v", i, err)) {
			return
		}

		if i < 20 {
			nodeIdArr = append(nodeIdArr, canTmp.NodeId)
		}

		v := &staking.Validator{
			NodeAddress:     canAddr,
			NodeId:          canTmp.NodeId,
			BlsPubKey:       canTmp.BlsPubKey,
			ProgramVersion:  canTmp.ProgramVersion,
			Shares:          canTmp.Shares,
			StakingBlockNum: canTmp.StakingBlockNum,
			StakingTxIndex:  canTmp.StakingTxIndex,
			ValidatorTerm:   0,
		}

		validatorQueue = append(validatorQueue, v)
	}

	epoch_Arr := &staking.ValidatorArray{
		Start: 1,
		End:   xutil.CalcBlocksEachEpoch(),
		Arr:   validatorQueue,
	}

	curr_Arr := &staking.ValidatorArray{
		Start: 1,
		End:   xutil.ConsensusSize(),
		Arr:   validatorQueue,
	}

	t.Log("Store Curr Epoch VerifierList", "len", len(epoch_Arr.Arr))
	if err := setVerifierList(blockHash, epoch_Arr); nil != err {
		log.Error("Failed to setVerifierList", err)
		return
	}

	t.Log("Store CuRR Round Validator", "len", len(epoch_Arr.Arr))
	if err := setRoundValList(blockHash, curr_Arr); nil != err {
		log.Error("Failed to setVerifierList", err)
		return
	}

	// Commit Block 1
	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start ProposalPassedNotify
	*/
	err = StakingInstance().ProposalPassedNotify(blockHash2, blockNumber2.Uint64(), nodeIdArr, promoteVersion)

	assert.Nil(t, err, fmt.Sprintf("Failed to ProposalPassedNotify, err: %v", err))
}

func TestStakingPlugin_GetCandidateONEpoch(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)

	/**
	Start GetCandidateONEpoch
	*/
	canNotIrrQueue, err := StakingInstance().GetCandidateONEpoch(header.Hash(), header.Number.Uint64(), QueryStartNotIrr)

	assert.Nil(t, err, fmt.Sprintf("Failed to GetCandidateONEpoch by QueryStartNotIrr, err: %v", err))
	assert.True(t, 0 != len(canNotIrrQueue))
	t.Log("GetCandidateONEpoch by QueryStartNotIrr:", canNotIrrQueue)

	canQueue, err := StakingInstance().GetCandidateONEpoch(header.Hash(), header.Number.Uint64(), QueryStartIrr)

	assert.Nil(t, err, fmt.Sprintf("Failed to GetCandidateONEpoch by QueryStartIrr, err: %v", err))
	assert.True(t, 0 != len(canQueue))
	t.Log("GetCandidateONEpoch by QueryStartIrr:", canQueue)
}

func TestStakingPlugin_GetCandidateONRound(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)
	if nil != err {
		return
	}

	/**
	Start GetCandidateONRound
	*/
	canNotIrrQueue, err := StakingInstance().GetCandidateONRound(header.Hash(), header.Number.Uint64(), CurrentRound, QueryStartNotIrr)

	assert.Nil(t, err, fmt.Sprintf("Failed to GetCandidateONRound by QueryStartNotIrr, err: %v", err))
	assert.True(t, 0 != len(canNotIrrQueue))
	t.Log("GetCandidateONRound by QueryStartNotIrr:", canNotIrrQueue)

	canQueue, err := StakingInstance().GetCandidateONRound(header.Hash(), header.Number.Uint64(), CurrentRound, QueryStartIrr)

	assert.Nil(t, err, fmt.Sprintf("Failed to GetCandidateONRound by QueryStartIrr, err: %v", err))

	assert.True(t, 0 != len(canQueue))
	t.Log("GetCandidateONRound by QueryStartIrr:", canQueue)

}

func TestStakingPlugin_GetValidatorList(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)
	if nil != err {
		return
	}

	/**
	Start  GetValidatorList
	*/
	validatorNotIrrExQueue, err := StakingInstance().GetValidatorList(header.Hash(), header.Number.Uint64(), CurrentRound, QueryStartNotIrr)

	assert.Nil(t, err, fmt.Sprintf("Failed to GetValidatorList by QueryStartNotIrr, err: %v", err))
	assert.True(t, 0 != len(validatorNotIrrExQueue))
	t.Log("GetValidatorList by QueryStartNotIrr:", validatorNotIrrExQueue)

	validatorExQueue, err := StakingInstance().GetValidatorList(header.Hash(), header.Number.Uint64(), CurrentRound, QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetValidatorList by QueryStartIrr, err: %v", err)
		return
	}

	assert.Nil(t, err, fmt.Sprintf("Failed to GetValidatorList by QueryStartIrr, err: %v", err))
	assert.True(t, 0 != len(validatorExQueue))
	t.Log("GetValidatorList by QueryStartIrr:", validatorExQueue)

}

func TestStakingPlugin_GetVerifierList(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)
	if nil != err {
		return
	}

	/**
	Start GetVerifierList
	*/
	validatorNotIrrExQueue, err := StakingInstance().GetVerifierList(header.Hash(), header.Number.Uint64(), QueryStartNotIrr)

	assert.Nil(t, err, fmt.Sprintf("Failed to GetVerifierList by QueryStartNotIrr, err: %v", err))
	assert.True(t, 0 != len(validatorNotIrrExQueue))
	t.Log("GetVerifierList by QueryStartNotIrr:", validatorNotIrrExQueue)

	validatorExQueue, err := StakingInstance().GetVerifierList(header.Hash(), header.Number.Uint64(), QueryStartIrr)

	assert.Nil(t, err, fmt.Sprintf("Failed to GetVerifierList by QueryStartIrr, err: %v", err))
	assert.True(t, 0 != len(validatorExQueue))
	t.Log("GetVerifierList by QueryStartIrr:", validatorExQueue)

}

func TestStakingPlugin_ListCurrentValidatorID(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)
	if nil != err {
		return
	}

	/**
	Start  ListCurrentValidatorID
	*/
	validatorIdQueue, err := StakingInstance().ListCurrentValidatorID(header.Hash(), header.Number.Uint64())

	assert.Nil(t, err, fmt.Sprintf("Failed to ListCurrentValidatorID, err: %v", err))
	assert.True(t, 0 != len(validatorIdQueue))
	t.Log("ListCurrentValidatorID:", validatorIdQueue)

}

func TestStakingPlugin_ListVerifierNodeID(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)
	if nil != err {
		return
	}

	/**
	Start ListVerifierNodeId
	*/

	/**
	Start  ListVerifierNodeID
	*/
	validatorIdQueue, err := StakingInstance().ListVerifierNodeID(header.Hash(), header.Number.Uint64())

	assert.Nil(t, err, fmt.Sprintf("Failed to ListVerifierNodeID, err: %v", err))
	assert.True(t, 0 != len(validatorIdQueue))
	t.Log("ListVerifierNodeID:", validatorIdQueue)
}

func TestStakingPlugin_IsCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	handler.NewVrfHandler(genesis.Hash().Bytes())

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	nodeIdArr := make([]discover.NodeID, 0)

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()

		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		canTmp := &staking.Candidate{
			CandidateBase: &staking.CandidateBase{
				NodeId:          nodeId,
				BlsPubKey:       blsKeyHex,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(i),
				StakingTxIndex:  uint32(index),
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
				},
			},
			CandidateMutable: &staking.CandidateMutable{
				Shares: balance,

				// Prevent null pointer initialization
				Released:           common.Big0,
				ReleasedHes:        common.Big0,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}

		if i < 20 {
			nodeIdArr = append(nodeIdArr, canTmp.NodeId)
		}
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	/**
	Start  IsCandidate
	*/
	for i, nodeId := range nodeIdArr {
		yes, err := StakingInstance().IsCandidate(blockHash2, nodeId, QueryStartNotIrr)
		if nil != err {
			t.Errorf("Failed to IsCandidate, index: %d, err: %v", i, err)
			return
		}
		if !yes {
			t.Logf("The NodeId is not a Id of Candidate, nodeId: %s", nodeId.String())
		}
	}
}

func TestStakingPlugin_IsCurrValidator(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)
	if nil != err {
		return
	}
	/**
	Start  IsCurrValidator
	*/
	for i, nodeId := range nodeIdArr {
		yes, err := StakingInstance().IsCurrValidator(header.Hash(), header.Number.Uint64(), nodeId, QueryStartNotIrr)
		if nil != err {
			t.Errorf("Failed to IsCurrValidator, index: %d, err: %v", i, err)
			return
		}
		if !yes {
			t.Logf("The NodeId is not a Id of current round validator, nodeId: %s", nodeId.String())
		}
	}

}

func TestStakingPlugin_IsCurrVerifier(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)
	if nil != err {
		return
	}

	/**
	Start  IsCurrVerifier
	*/
	for i, nodeId := range nodeIdArr {
		yes, err := StakingInstance().IsCurrVerifier(header.Hash(), header.Number.Uint64(), nodeId, QueryStartNotIrr)
		if nil != err {
			t.Errorf("Failed to IsCurrVerifier, index: %d, err: %v", i, err)
			return
		}
		if !yes {
			t.Logf("The NodeId is not a Id of Epoch validator, nodeId: %s", nodeId.String())
		}
	}
}

// for consensus
func TestStakingPlugin_GetLastNumber(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)
	if nil != err {
		return
	}
	/**
	Start  GetLastNumber
	*/
	endNumber := StakingInstance().GetLastNumber(header.Number.Uint64())

	round := xutil.CalculateRound(header.Number.Uint64())
	blockNum := round * xutil.ConsensusSize()
	assert.True(t, endNumber == blockNum, fmt.Sprintf("currentNumber: %d, currentRound: %d endNumber: %d, targetNumber: %d", header.Number, round, endNumber, blockNum))

}

func TestStakingPlugin_GetValidator(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	header, err := buildPrepareData(genesis, t)
	if nil != err {
		return
	}

	/**
	Start  GetValidator
	*/
	valArr, err := StakingInstance().GetValidator(header.Number.Uint64())

	assert.Nil(t, err, fmt.Sprintf("Failed to GetValidator, err: %v", err))
	assert.True(t, nil != valArr)
	t.Log("GetValidator the validators is:", valArr)

}

func TestStakingPlugin_IsCandidateNode(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if _, err := buildPrepareData(genesis, t); nil != err {
		return
	}
	/**
	Start  IsCandidateNode
	*/
	yes := StakingInstance().IsCandidateNode(nodeIdArr[0])

	t.Log("IsCandidateNode the flag is:", yes)

}

func TestStakingPlugin_ProbabilityElection(t *testing.T) {

	newChainState()

	curve := elliptic.P256()
	vqList := make(staking.ValidatorQueue, 0)
	preNonces := make([][]byte, 0)
	currentNonce := crypto.Keccak256([]byte(string("nonce")))
	for i := 0; i < int(xcom.MaxValidators()); i++ {

		mrand.Seed(time.Now().UnixNano())
		v1 := new(big.Int).SetInt64(time.Now().UnixNano())
		v1.Mul(v1, new(big.Int).SetInt64(1e18))
		v1.Add(v1, new(big.Int).SetInt64(int64(mrand.Intn(1000))))

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		privKey, _ := ecdsa.GenerateKey(curve, rand.Reader)
		nodeId := discover.PubkeyID(&privKey.PublicKey)
		addr := crypto.PubkeyToNodeAddress(privKey.PublicKey)

		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		v := &staking.Validator{
			NodeAddress: addr,
			NodeId:      nodeId,
			BlsPubKey:   blsKeyHex,

			ProgramVersion:  uint32(mrand.Intn(5) + 1),
			Shares:          v1.SetInt64(10),
			StakingBlockNum: uint64(mrand.Intn(230)),
			StakingTxIndex:  uint32(mrand.Intn(1000)),
			ValidatorTerm:   1,
		}
		vqList = append(vqList, v)
		preNonces = append(preNonces, crypto.Keccak256(common.Int64ToBytes(time.Now().UnixNano() + int64(i)))[:])
		time.Sleep(time.Microsecond * 10)
	}

	result, err := probabilityElection(vqList, int(xcom.ShiftValidatorNum()), currentNonce, preNonces, 1, params.GenesisVersion)
	assert.Nil(t, err, fmt.Sprintf("Failed to probabilityElection, err: %v", err))
	assert.True(t, nil != result, "the result is nil")

}

func TestVerElection(t *testing.T) {
	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))
	verList := "{\"Start\": 2956251, \"End\": 2967000, \"Arr\": [{\"NodeId\": \"cc77a9f28ae8a89acc48ec969c34fad9ed877d636eccf22725af7bd89b274f702d40073d2dda1b9fbfc572746979a7876c1dca93fc617918586a53747d211d43\",\"NodeAddress\": \"0xe3fa4c4d652cdf451784e2b8e095aa2d3bab8d6e\",\"BlsPubKey\": \"7af22c4d167dc3ae266d1bfffdd4807f97f6dd8a14623e1a903bce1df60da3cf3d7021ddbdb1f17cf18475f747d74608b350df6f07351480740590b1a56de15ddb4a1a561c76af1e7b6f6944ad7ac165f510303f84f87559f7e8126d2bc6b60d\",\"ProgramVersion\": 3330,\"Shares\": 34518546279650000000000,\"StakingBlockNum\": 501228,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"3a30e0ca53a5ad43e92d6df887468d9f735b756c153cf5bfadbf70f5e3851cecbe0d88f48b1b26a9a9abe7cd2e0ba45758ab710be68b345525a19f6782798ef6\",\"NodeAddress\": \"0x55d50f604564d4adbcdb1f1ebd9e41598ae17061\",\"BlsPubKey\": \"c72328d711cb3c55fa147c128bd28ffb9df3bac1f072b558cb369e4a66ff9695846ccdf013b7df574f3afca916fb0b19d365a4a5e8b35cae6df6af7150e00a455fdc0f3659d1eb3b15bf86bdc80c151389cc44a14dec0bd649755094fb8d390f\",\"ProgramVersion\": 3330,\"Shares\": 32460963287070000000000,\"StakingBlockNum\": 504218,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"e628b3a40b18d02203099ab74bcc4d7b94350001c08a2f0b9060348962656aa6d8f021a46494c545a352978c6f837a67618a9b4a7f59bc7e90764d2f4006ae7f\",\"NodeAddress\": \"0x9d71aa21e8fbb12d297d6322465e2d47d500926e\",\"BlsPubKey\": \"7127bbe5dc06a1ff17d763b2e34aca3c488c0d7065a8e5160550410ecfca05551ab75981efc8ad7a3af8fe113d54f903ec0533312cdfd3b2607499cf43eefbac2a12502f65dbc66cbefd299059b0916278876e73e83d4249a3e4e1e842e60f89\",\"ProgramVersion\": 3330,\"Shares\": 23574467265790000000000,\"StakingBlockNum\": 1362243,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"5f2bf1ada8117f9fca7117c3d402f375d6b97296f2e40cac4eba9f3cf137d1d046468d11da8a670dab6ca468bf0ee770c4ab6246887fd78c85bf244bcecbd255\",\"NodeAddress\": \"0x00b505f612efab92722496efbbb6a1d85cff8ba9\",\"BlsPubKey\": \"5ceab320ae03fa246d0b1cf5558adf6b061bb2e626e303655009d58ebe22ddee43350e493dabeba1c8d9678e74a33d01e6b3d4c33998a5d9e91d21f3edf7021d77e0a18097e7ef6c827e2d2ac8a2831761ebcd21155f73cdfa082723eb492291\",\"ProgramVersion\": 3330,\"Shares\": 22805711411757777777779,\"StakingBlockNum\": 504491,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"8f1c8333053cad81c76baf4df547eab49d1d1e1f7602b470ed3cd36e14461042731b20904252bff629d1fe60b05c3223653410c72a30fe90ef66ed95e6f4289a\",\"NodeAddress\": \"0xbc01d8e6451059ebd7b395352f1c7e16f52d630e\",\"BlsPubKey\": \"48f9db95995c1dbe724446b693fef640d226abe9a9c6c473310e387a99f43fb79b02a533f805dfe0d69c6bbe5cee3207e9e09f3d7d6a09ab2f51eee62050c6260d4e9cebfcc5daf0eb11e54aa6aa85440297016c34315168cb971ba330932198\",\"ProgramVersion\": 3330,\"Shares\": 22125894495482071394430,\"StakingBlockNum\": 2346506,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"89baa83d0447aaabdf21f784532060bbcd0ca18742c13b09553c233e42225b36b0d52096d8e2a9db18a370a1b01099fa8eda654536c9847105ad8f6431ad5e55\",\"NodeAddress\": \"0xf677bb3671babf020081ed553b35ae0c0750b8bb\",\"BlsPubKey\": \"d4dc2bcb9de271da287ead8a28a17c967a85f6b4dff9b2e0a110b2500d77fbb0f3e272de81fe0b40cd11daab1eec5e09da6383eb634e7ff206d1a0a8784dc7f0d34c3f0e0fd9f6b9ef2944fe89d34f0f4183ce8f64c92cddf45e8e163ea0f084\",\"ProgramVersion\": 3330,\"Shares\": 21530654151080000000000,\"StakingBlockNum\": 504350,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"0c2bf9b0d73ebd5728e810deb55a6f656cc44b92a8399aa0898eaa26bad3232707441e381f7a83eb1c0b38bde813bf84e7311b8d9fe75f7c5eea6c219c3103cc\",\"NodeAddress\": \"0x34c7c9dd1836d559b230b66fe1c0e7a3d5442a42\",\"BlsPubKey\": \"17080c4ded2600a27cf20209640b621e0591a9601b7c1f0e8e2648df1e83c9ef286e4490e3b752b9ee4e7a54f0a44212c8dc6b6d5e4df2057c566bb1c54f99595ca6666b81422b7707d39ce4ae12318184ffd509e87b246a1421235a709f7403\",\"ProgramVersion\": 3330,\"Shares\": 20451680760985649790440,\"StakingBlockNum\": 502949,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"c67c30ced9155115d3fe25c61a586a0597362b8c9cf5d63bc64f3ada2b7943e918a6d77f3e8c2e3f998f13262718b547b6cfe51bb5808d0e1e82c6019c223eb7\",\"NodeAddress\": \"0x9a15f542f24366cebccef9e116de2d73e590f64e\",\"BlsPubKey\": \"33c58ca22ed942a679c5542e833f22a59cfac05a0dc1e2ba3f11901cb24f5aa948ef85e10b889284338c2d10c6f1060559e90c090e76b561ba4f9c433fe1984aefd74b43f04fefa3849c10fb89a4dd4bff2aee32d80b95a793487ab0796cdf00\",\"ProgramVersion\": 3330,\"Shares\": 20287661049100000000000,\"StakingBlockNum\": 545395,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"55d6d96974e14b96aa4d6a5f629eff081fa7b9f47f4cd9dc700fc07cf8532d0336a118612f7c085a96bbaa775ab7db03dadce4759f3168d9dac162e1955d6167\",\"NodeAddress\": \"0x76c60cd52a526b936961d53236e9ffddabc9adde\",\"BlsPubKey\": \"5eaaaa17dd1b974e66148999598392835a58687b64668d424889efc1046383db87283b1cf2cb7356c95b475749714b0a9a73356011783cbafb2817d4172555bf27e0bd030756be6615bca1a62f0cc4030d404fb2e05075222df35696d6c6ca11\",\"ProgramVersion\": 3330,\"Shares\": 18996171244328749999100,\"StakingBlockNum\": 519282,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"5c7c9dd985ddf54ddd865c4c26d0bfc291568f421100eb96fd6c69593bf226a4d18257f00234155d3e8a8b12a24d89525711f58e47e643dc1c198448dbce9dd1\",\"NodeAddress\": \"0xe6fffd1fc6e3c8abce6ca5846a6baa87ecfdadec\",\"BlsPubKey\": \"27d389e976f835428c56ea9911ffb0aefce19fb7a36be72e43530323e129c595aad2ec95e9d3524455bb04fd687f5316f83c757e26b24cbced61571635fb9f10867975115570cfdf234227864074d463cc8883f17a1a9f3831eb8bac8f0ed08a\",\"ProgramVersion\": 3330,\"Shares\": 17837436721380000000000,\"StakingBlockNum\": 2579769,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"d0f648262b54c16a0c4f82d1c2b08b638620537fc26321a47ace94b6691e1fa289fc4e06e1d94113b33ad1c33e3726655448f6f6a8bd1ea1c99a9f84da92c3d9\",\"NodeAddress\": \"0xc1c3d8fbedb1f5c66a2bf212cd3846ac043f4f2d\",\"BlsPubKey\": \"fde87f9f32308139e75917ada4e8f80468be72c039778bfb1cc3d9e4fa6dca6d13159901ba793d0c0e940dd7afa8f40c16d95e8520b157fcce83ace8e3eb914aa8648ceea9a21aa353bd205831a613c87aad0fdb86a4533a061d85901e510311\",\"ProgramVersion\": 3330,\"Shares\": 17590569142200000000000,\"StakingBlockNum\": 511342,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"4e44d7995f1b2b01bdbc65e5448f9fe3c7082f46f85fe980bccc3418836cec8ffab248dc0b54c0d8974e1327b8992643ef9db222ec9a44512a092d33b274bc1d\",\"NodeAddress\": \"0x0519b81a5aab47fcc17625297c135be8cb4f947a\",\"BlsPubKey\": \"4c15c75a224236dd21ffdf45584e4e498f59f5d4cc7c8f4d33122e063733b2b5aba9eb0bb174c2c0b055856fd646ad0e49ec2115d6a1497aa3c87417e50231082b1869ce7d6ccf89b0dad06a1d97033ccbc92c7147a77b683a50cb890fcc0082\",\"ProgramVersion\": 3330,\"Shares\": 17578743408780000000000,\"StakingBlockNum\": 502981,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"68ca2e833cdb044017d5f2bc35c2c61ecd5b340cac3384b3b33116c7caf951afdfef6fbf2542322de06a485beeb7b0e5c198b819c3dd991958eaa36c148bb356\",\"NodeAddress\": \"0xe14d1d48a3fe6581be663769c07958dce34b814e\",\"BlsPubKey\": \"1585198354f31d5ea60e898783ff9683e22ec70c923918fcd7ecaf06eddd3835b3d33b03490f78588c672ed2e7ed8b133b55178e44c38efaf9374b061fa5b36a2ee0aee5561d006b6cd4f1014b5a645e71a0bcb874d1e83c9e85fdf4e1faa216\",\"ProgramVersion\": 3330,\"Shares\": 17513708048380000000000,\"StakingBlockNum\": 2494715,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"1dbe057f33d9748e1d396d624f4c2554f67742f18247e6be6c615c56c70a6e18a6604dd887fd1e9ffdf9708486fb76b711cb5d8e66ccc69d2cee09428832aa98\",\"NodeAddress\": \"0xc68734e38851cd77b13a2eaee4873f866ffcde10\",\"BlsPubKey\": \"1f36e0553f73583909d8cca8ed9bf3df9e86abf3be0c4bcc812c57f16646f39a9e89bdcb3be897cb17eed9c244336207aca11ad7255087f45e131f51d32f583716e0e152a9f569b72a129959d970f1fcbf2f7241d1be4de6358db09ef4aaba8d\",\"ProgramVersion\": 3330,\"Shares\": 17436913713494264227100,\"StakingBlockNum\": 504281,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"1a85eaa54fd49ee4c79bb5ae72e9ce94bf826b2aa0cedd8d8d0eba41a634027f5c1fbe5de0b0b4a74837dd4c260d887c602674a11c5e611fe0b93b16a9e1e64b\",\"NodeAddress\": \"0xa5e23ada52f3b72ace9ed09095b47967ad3e2065\",\"BlsPubKey\": \"3955af34a4b501060ebeaddcee53a7cf380e500ce0d471e84141ed5ade9610158b01332b3fc993102e6a8eeab1e01509a2e39dcdaf741eee558eaeb9adfe2d4fd40ba91ca1e3485c2b0b0e0a5e2134bc91a80742abf2eefcf3c48c4268b91085\",\"ProgramVersion\": 3330,\"Shares\": 17237822748218677755418,\"StakingBlockNum\": 501312,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"0eb6b43a9945a062e67b45248084ec4b5da5f22d35a58991c8f508666253fbd1b679b633728f4c3384ee878ca5efca7623786fdf623b4e5288ace830dc237614\",\"NodeAddress\": \"0xfdc018925d09d7e8d30f2ba8d22ec5c38895b10d\",\"BlsPubKey\": \"f455a871bdf294b297f2b3a1abf284c335ba7aa78a8f7a822cc2504d1bdf1780819d1ec8633dc70c2d0904a2ed530715b93a748c4236cdfc0c4e4660f7acf550f13cad1eeb94e415235a2128cb8b947dda4b057fad7ed6d64349d106d5485100\",\"ProgramVersion\": 3330,\"Shares\": 17146065929655883087462,\"StakingBlockNum\": 501216,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"fa24468c3cd252da1ba8872665bfd8250db1017173ccb4fb8c42296bd2af5114bb4d575dad05ee2228e755a24665131c4b8ce0d8512c90c49bda82a3468624e3\",\"NodeAddress\": \"0x82daa0f33673aa7d603acec578ab6456b377cf83\",\"BlsPubKey\": \"6c63872a0495975998fbed9b93c7ddf0e987b30d8f837c92d675a665ec8f9b745a9e2c022b1713e761b7b06fba199106c2d2be3358a8ee4bc60e33de47628f959982cc2c7319cc50ce56bb2d39b0f8e94016b658be9dcf2ecfead790f1a5b18c\",\"ProgramVersion\": 3330,\"Shares\": 17128651657010773912267,\"StakingBlockNum\": 501285,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"79f86478381b2472b009f790dd0b8b32f65169cd3a447e43eafb0f359f9edb16ef4d15e45cec54adc5e691ddf43d6ec29066b8dc90457482cf9189caebc9d99f\",\"NodeAddress\": \"0x9f7d2586809d1918ae8e00f5753c969574bcd1f2\",\"BlsPubKey\": \"1874e9bde7d474e70a0699d451d5fbbf9ac7b7fba4fdfc86cd3a0a6543ecdf39fc8c0c8ae1c5053db2ee354366bb3302fff6789ed37f1bedd5dd00169a177711ebdafacab439c633d11096e3ba16a6a882226c11ef3785b1542ad9ca42d08d03\",\"ProgramVersion\": 3330,\"Shares\": 17119642829111484850576,\"StakingBlockNum\": 500834,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"32b7ca6dec2f4e96187d6dbbed02a224dc91b605302689ac921f069c9bd314d431ffaa4ecd220164513cba63361a78e672ba2c9557e44b76c7e6433ed5fbee94\",\"NodeAddress\": \"0x620a35919a2091c10b7a009335ce018fe5a8389b\",\"BlsPubKey\": \"8257f6a0a7d0f968be4f659129ffcb60300b6e6b921a167b3ced96610d4da7c1c9ebc414b23f14c93e88ef7221f7a301099901913819ba527e7dcc8ee8ebf6c51d6b727d3bc8368250fe743019faefd540e48e0313f0b01e9c01131974d5258f\",\"ProgramVersion\": 3330,\"Shares\": 17109574638651925151293,\"StakingBlockNum\": 501504,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"a3f400758b678e15d7db188f2448411b8b472cb94da7a712dd7842f03c879478e156fd7a549a782a12de5f5dd9dc979dfd2946f8396b9328dd4fbdead37e49ba\",\"NodeAddress\": \"0x5329eac67cc1e5f006cfa7409143c4d0dd27de0e\",\"BlsPubKey\": \"ecbf9c5c303114b2d64274d31814a6d4084857a228c9938f473c775ca89972ef30de5792fad4bd7ee55d00eb5c8d5519f52a62d283ee4942197aea24085ec869aa2038c24e9bda607ba11392d7cceb0c508ec0c39cae16198024a0074fdce109\",\"ProgramVersion\": 3330,\"Shares\": 17095203760016677231399,\"StakingBlockNum\": 501191,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"95c3fbd07041e78483ccc11a598ceb6b7e9bde2dcae68b5a7c4ea6876a423801d1fadb0de79bdcf52376236bd4ce2c1131a1e2640bd26856c967e17415184924\",\"NodeAddress\": \"0xbeb95341e03b6840b8062c92cf8eae04b7adddb1\",\"BlsPubKey\": \"60354ac94671d8cd86b137198986e08e1f71dd66d1f5f2734aeb04ee79440bb334dcddeff249aee35ac8cba9d31bbb1602a708cbe071f76a4252547530c6acad8a97984c38884db1cc9eb4cb15a2aca208da04e01cb2d412013aa77e9ba5fe80\",\"ProgramVersion\": 3330,\"Shares\": 16795434182540000000000,\"StakingBlockNum\": 501689,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"f01b7b56c23a67376f0b1578feefed152952bbc3326fd2c412432a6f7e95ebe4d8adc231c8f3a4e3ab7d35c05278073766f89c5565a1f5267cd73a07c1d2fb04\",\"NodeAddress\": \"0xf4fa244af04b37757f68fc6ee3a030c8e8c7dc47\",\"BlsPubKey\": \"7729c6673584cddf6ba5d10431b2cbfa32ccd7ded7024dc4e27937076d27427b12dac770e99f3936795aa96b0dff5d147ee9d6a53549bb27d759c0c491164e8c67a792b41ac967cd6b2655f325c6b0652fe04cfe4544a61b94c6dd801d9ee50c\",\"ProgramVersion\": 3330,\"Shares\": 16782651890870000000000,\"StakingBlockNum\": 504821,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"239a70bf6909fe8f63bdb298799fe81d29c45d9a35952ae0cee48f77d79c678f9e9c4e2ba63c0691c0a5d5adcedc159b361bf0d20d6a3be534adf080e6d2205f\",\"NodeAddress\": \"0xc274b46edff09e27dcdc69958055bad82942ab9e\",\"BlsPubKey\": \"b67df19cec0e4ba905d6e83fa952f3a01a7eaa6eb9ca260c3739527e39b5740ed2c4234eaf72a7452c5c5112b16cab135b7b0b983aeeb3e00efeb5b3baecfdcdd1d9c6c0a7eaf187404a5a95852dce88e4f927858fd2b981092b82c52d0ffd96\",\"ProgramVersion\": 3330,\"Shares\": 16768031913180000000000,\"StakingBlockNum\": 503116,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"168fff55f6a28c3f10802e90cefb4daf58eb9a60c55eaaeda6588b20eece884ab9104da866633929e3c0ea413de2624b361939cf1513e1b67460686b7192f528\",\"NodeAddress\": \"0xfbb0006f5ab5b690bb11a6aa2a149de4e03f7c06\",\"BlsPubKey\": \"625e5d23f40999be0baa0c2a2778ee8c94c9a16a1dbb0e273949dbc337b0a5f10653c12f556d21312906e15468831a054cbebf02ca1454da76fb0d3e1a606b31163c1366970ef3287e432d8d51d884c2e614422caadf24e46408f574fb71d690\",\"ProgramVersion\": 3330,\"Shares\": 16640439222705189967460,\"StakingBlockNum\": 504425,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"9460fce5beea98e4d56c62a920bb041f45e48a5a7b96d12d02a16cbb20863be9c76491127533d9cefa5b4cec48ae6595b7ba347ef7dc8277cfb343eebde4646b\",\"NodeAddress\": \"0xd2d114926a5d2772e17b6a6d21915a5dd541014b\",\"BlsPubKey\": \"b7f3cca67a2177b7503c4cbca39d97962ce3d0cdcfde33e0028ca646800380c1a5a717ee15708a0b83ff7990a279a9052133c1c6fe676deadfe7a0c0df7d4882927f6f3c38f264842655d0037e772fca65ba145c8c74cafca54cbb87a6881386\",\"ProgramVersion\": 3330,\"Shares\": 16378755752354505353500,\"StakingBlockNum\": 500468,\"StakingTxIndex\": 1,\"ValidatorTerm\": 0},{\"NodeId\": \"853667f088dbad823cfb9fc956ddf343cacf4cffc9d5de1e7227a2261bd2c5df356cb0a914f9404f21cc502d46a9d87bf90c5863286d07ebb0b06fd71f9c1192\",\"NodeAddress\": \"0xa5608dac040324ab8a399db4134e42eb971b6f6a\",\"BlsPubKey\": \"9a22fac3d51645540bcbb457dd21a8a55da7eb2432678f0328a567ec52069d10f7d81f2f245817b4ce01caedc5cdbf095ef327d6893067850df7301c85f8106d75670f23a420e647c0ca0c49fa232652442ab538232df3652f71137130ecc88b\",\"ProgramVersion\": 3330,\"Shares\": 16070247319030000000000,\"StakingBlockNum\": 502165,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"fff1010bbf1762d13bf13828142c612a7d287f0f1367f8104a78f001145fd788fb44b87e9eac404bc2e880602450405850ff286658781dce130aee981394551d\",\"NodeAddress\": \"0xc7b552cfaacf3bb7a987a17895620308a827fcf6\",\"BlsPubKey\": \"31500f22231733201efe29b22c58d4b2d135c3cb29ce39ff75a63e180b2bf2530b87b36fb7a00523bc63d23b3ed6780c297b9581328a3d42275b87c9bfd3d11994fc1f440184368e0cc044ead7b9bf42a04273d403925b8b7c9c73ce7574df0e\",\"ProgramVersion\": 3330,\"Shares\": 15811633221920000000000,\"StakingBlockNum\": 902037,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"f6cb212b47105dbbb5b42e52492c330ba92b9a40706de88d7f3de02e1fc5507e43bc25db9b474a4b75ad4c2091d79f297e1fe5d3e09f8e4a12aba7c26db33717\",\"NodeAddress\": \"0xbf787567d3cf5b8aa8a15319fcbb5034972e7909\",\"BlsPubKey\": \"dae3743c8f90d48d330706a7aa79415ab858eaeede837db929df13f4d596353293dbb0480c693b2659906ebfd4dd4a057190d1a142b69506302cd028e0dff287ed71e4de93b0d56d660d871dc7b9a31f3e5628e58b6c12b946d9ee5cf4620098\",\"ProgramVersion\": 3330,\"Shares\": 15811133161830000000000,\"StakingBlockNum\": 2416508,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"511ab9921b1ecf4bd8c76193a1c281f57a03190eae5418a23a7920a8064f89b9022bfae56d7fd2e740c2bc90c07e7aa78f201fa19c8b2b6e0dd15d8c97bec8c6\",\"NodeAddress\": \"0x497b41ce07248b61af71f58fed4e07257767bb80\",\"BlsPubKey\": \"7cd9c47600b8f4b8974eea5002d4fa83d6980d7298101697b6ddb66e99a284a168db66130161b66ee52b23f1aeb35712136e7d0b9106ff1e4d97095470d3deb9e5482129e1bd523ea41711bb808f6d140ba7f2220ce21f9c2862d16f996c0f8a\",\"ProgramVersion\": 3330,\"Shares\": 15738052896570000000000,\"StakingBlockNum\": 502098,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"4922ef029b4bd93a1bca252d3ce4ed4d50192179e1c06e3ef452a18fe17d21e8c42eadf7d6aec47de2b90b07befa6c448d82f55599be6141c370dfac0b56f547\",\"NodeAddress\": \"0x29776a7df0e0af65079d114601300fc5248b49c9\",\"BlsPubKey\": \"951e973bd15e1cf6ca363f50f189e82c05a6fdecd792cf31c09975da5454c7468b47225563c148e2171b99fe654a4502ebf550ae7d7555b748425b8382eae45d1cd1fe6925d2c8dece1e286292b41b0426081f75c5d926c82ca1b34792ddc50a\",\"ProgramVersion\": 3330,\"Shares\": 15721622597090000000000,\"StakingBlockNum\": 500969,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"580ca2b64f4ea1ae5f266b62e475b456d168d84debed84185719809c8b0e35f8c03271e097b5a80cfde67c29ce11bd3a76cd95befceab5a050565240250b74e5\",\"NodeAddress\": \"0xcfac8dad806ade6f709819536271b7d3f0b361a7\",\"BlsPubKey\": \"b912fb0ba265f59e7be75000ca445bfa1a42388ae2acfc95d969110ec2b6917a24339b83e37ab453c958b6b04ed444001b646ce49f60c6d5a1788df6cd09611ed10036de6dbe1efa180a3980180a646076827b06a2046f131b200e9b51d49811\",\"ProgramVersion\": 3330,\"Shares\": 15668552452290000000000,\"StakingBlockNum\": 983559,\"StakingTxIndex\": 1,\"ValidatorTerm\": 0},{\"NodeId\": \"919ef95c3ffbc36a458fb3f7d11ac2734ef9507eda3c1cb13b69ce950a0ddc57c29655a0141a6e68a5d42d0f1c5d7bc695d71b440b458b42a7e3a06487b3a2f5\",\"NodeAddress\": \"0xaca4356ad0ad2cac0dea4f50a9701b855b1286cd\",\"BlsPubKey\": \"ebd3b8ecb79dcf4381ada8daee1aaf593d37574f6f2782b95976c015c66bc0de030e9930a08d6aa4ea361394bc9d6c0f8aa7ff184a8556d024cef121e5e9de644619680ef8069e8435f8e593ad33a6b7c92a84f325408e5d355ea22b11e1dc15\",\"ProgramVersion\": 3330,\"Shares\": 15644735097337401731584,\"StakingBlockNum\": 503012,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"ed552a64f708696ac53962b88e927181688c8bc260787c82e1c9c21a62da4ce59c31fc594e48249e89392ce2e6e2a0320d6688b38ad7884ff6fe664faf4b12d9\",\"NodeAddress\": \"0x72e49c76d69fcfef1d83c7de3e173e41644c9e06\",\"BlsPubKey\": \"c70b6490377a1b769e8c0d3f94d50e443f560119a63d4fdbc6a76235291197a405b680e9c1a366445b4edb52849076069544ed42c0ec58ec35dab5521720095392e3e5eb0d57acca9f066cf4424a11c1d716e344c1f7d939ce22d030e561ae8a\",\"ProgramVersion\": 3330,\"Shares\": 15622205109200000000000,\"StakingBlockNum\": 2405590,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"4be36c1bda3d630c2448b646f2cb96a44867f558c32b3445f474cc17cb26639f80bf3a5f8bb993261cf7862e3d43e6306eec2a5c8c712f4e1ead3bb126899abd\",\"NodeAddress\": \"0xc7e6d08e51376c23f84e933d1b85a5419982230d\",\"BlsPubKey\": \"5ff1a7d5a44243ed7ee8fdba53f2cdd9d3b4b51f4272bdb41d55d5d107b2d8874432268a91b0a80380a81dd3e969f10d9f1fb081c683dead84ae0d6669b3d3244f278eba7e76ece2f3d177927ae99dec9ce25e81a809010bf7141ab21fea7d8d\",\"ProgramVersion\": 3330,\"Shares\": 15443602497100000000000,\"StakingBlockNum\": 1911772,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"680b23be9f9b1fa7684086ebd465bbd5503305738dae58146043632aa553c79c6a22330605a74f03e84379e9706d7644a8cbe5e638a70a58d20e7eb5042ec3ca\",\"NodeAddress\": \"0x75a6a6ad63a15f7bac83c223c46dee2a1b048f28\",\"BlsPubKey\": \"fa95e5661585ac75826c7f9695573667fbd49ca32796e9586092c0dc9ed350f57258e0379d098ac60db2efe8474cd206c5519b64097da582ffb29cb445fe8b3e053df48ee9df52ef7951ebce63859381ab55845b4cbedecb7636fa3bb7b69794\",\"ProgramVersion\": 3330,\"Shares\": 15405732675040000000000,\"StakingBlockNum\": 2070321,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"71fa06ef88ecc022451c3aa82414009ec856ccc0288fc13ae1c7a0e1bb7e914c8f419238acd2ad4bbccc15ed8f9b17e50aa7f6d6e1c8a6b53f95bf8f0817b6ca\",\"NodeAddress\": \"0x89a14254f6e3d19cb7804f922b05ca96374c6275\",\"BlsPubKey\": \"41def228353b05e9d6f74f5c7226c6ffede15dd6c153fdcd0ff35b212d3a39ab13b176000c27917a039ca17038bd2a09ef89ad3c6193b9160782cb0f7c845ce6f530679c1767ef0e76bd39630a60e7489c3adf70fd473ae42ae6224d136eb817\",\"ProgramVersion\": 3330,\"Shares\": 15329573047885822000000,\"StakingBlockNum\": 503991,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"00aa2e361168765c0df1ebe981455dd64f4c091adb01e162567101f8f1b5fa31ac960c198406a00fc558ebb8c14843ef72ed47d5931d0605a0e6caab41c8a86e\",\"NodeAddress\": \"0xd0b6726fd3c8381f70353cf35d7ccf2e2754a488\",\"BlsPubKey\": \"566329278be248579ead55b38d02dc7ad0cf334939d55cfa265301da9731d0b7cebabd62a788f8dd4cd241416aa87e10f035c9563f4bb3570dbf856f3797a8d9418db571d1b1e9b7c9fecb25f5f7f308bf4c47da793eb428be6c69da36aad187\",\"ProgramVersion\": 3330,\"Shares\": 15319269848140000000000,\"StakingBlockNum\": 2716479,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"73686446756b83458b910015775081acaa83d6727f4dcf71b732a827ebd563bc4553e1a3171029899ca276f285d7b63b36609ab8aaa0f138eeee38724d89a15e\",\"NodeAddress\": \"0x759823cd8b4ce98850f01858a3b2bef8edd6c254\",\"BlsPubKey\": \"6874380b357448a418ee072962361c2eb94f2a5295bd1e219dab29bfd36ff0339c2263b4bf69590f012bbd1465282f06e306429c546c6c78207f754a30d757512f7896743e06b2c21288e9e7c319b233c1a6de944cb95313948de8b53d6eb785\",\"ProgramVersion\": 3330,\"Shares\": 15306629151480000000000,\"StakingBlockNum\": 2311143,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"6eeddc7bea67b08cb2bff68b4c5c6cda0b4234779a21b78b244203acab504b801ca299dcbf5b50c4650e502196b05f4fd6f3582fd0973528fd3047c36ab2198d\",\"NodeAddress\": \"0x0722802e626858984834298cf878880f68389d62\",\"BlsPubKey\": \"6551128e8153bfd0ac148195b0b13f06c58732ee66518d2f6cb12179b3752d3d6fb7a1e78ffd358f2a7a59d4fb6b2017e5804bd71b68ef8771fd355e32a9e4cf5f3f1986ea8b0e7cc6db6561a59da5e89a4ec8cf5d8ed86772cc284e3964d416\",\"ProgramVersion\": 3330,\"Shares\": 15295000000000000000000,\"StakingBlockNum\": 500819,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"a6bbf49d6f1df8e3de68123802c182719fc2b17d40d655a0233a30f5131f35f58d48c87a462bdc58ec1a9014c340baebb9220d1b7c965c4f45e6cc7ca2ed6e81\",\"NodeAddress\": \"0x4598115f7bbde1a0b9d3ccc051b23ca15910a5e6\",\"BlsPubKey\": \"44f216a39546d3d0aa883588e0b1eee82b82c5c104398c34b7186bfc52d7be2a21fd3ae59348a6a1810699b362bfea02a248c778d95591170d2bdbbf7ad6d47a3b4b0eca5e02334511a5881e7ef3a1d3cb610cb51a054ce0cf7b620dd2f49781\",\"ProgramVersion\": 3330,\"Shares\": 15294354411592738953000,\"StakingBlockNum\": 500470,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"1d13161894d676fa124655820b92c722e04194ebc82bd2df5487dce684c21af80bd1e95230748d5aa65f2c28176fb231297f0cc5257a49051a7e33f3dd114b93\",\"NodeAddress\": \"0x6049ae08f8814324566ed7359a10de5cfff0fed5\",\"BlsPubKey\": \"fb46f3f276dc62ed0f433ad2f3705eff8e3f33cac7d629db3fb170a537a9be3c9f5ad37afdefbf1982ae86c97215e2008d4b60cfb72acfff4aa2eb467e331315e7fb8c0d3cd7c6fb4fa36b6b8b754ef3de55454a479c5cbe0951ebb69ac93904\",\"ProgramVersion\": 3330,\"Shares\": 15293631295150000000000,\"StakingBlockNum\": 1936836,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"11daedba2520a87da234f48bdaa4373e536b85367b863c218e6638c558214708470260830c2f17feccf8997187c70c96469bcca9ea0f5b522aae6fadea8e9ddf\",\"NodeAddress\": \"0x3a677f8c307e37879ca804ca8fb53e71f5712f2f\",\"BlsPubKey\": \"7ba318f2212a4c71175de4821d4fd2c4cbbbbb573e0b64482d373389ed5771a0e1e07b3d4cb07c7301227e63b3ef230d27f7f99b7f924c27e742e4d871cc14d452b5bfb94865826a79b7ab579811fe41bfcedae09687f10042e03e176082850a\",\"ProgramVersion\": 3330,\"Shares\": 15285260093600000000000,\"StakingBlockNum\": 1985569,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"73601a21e758e4052b43e68d9c09f2750dbc82b1736cce6a52258d42a18cde6c006d1eae39d44dda08987362d0ba8f8df43df05b3330281e866207117de25885\",\"NodeAddress\": \"0x6ea662bc535279d3d58943d91e9a83068b520943\",\"BlsPubKey\": \"8feb6b9d2eac780aea2ff48b1f18e30bbdba1e932f0d914f164fc72eaf6675330d2761f7ef6b6a8248afb9b5454d9c00e8f69e6174e43f512664062ddbe7e16826b0ace506426a7d22ce567b6314babaf3fe5572cb5177a9387a3de4048e1d82\",\"ProgramVersion\": 3330,\"Shares\": 15261433074070000000000,\"StakingBlockNum\": 500464,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"196e461fe3fea40d260bcef18eb0907fe61fb55322d3b446fe2dfa37127e1e6850435352c8a51635a755de9a25acb85a0efda56634e69be49a30d02908faecc7\",\"NodeAddress\": \"0xe277ea2ea1da21d564d8f3123bbdef15429fecb5\",\"BlsPubKey\": \"088034c9c6854654f74df863d091100bb9de792077b732f4c2c1e6cf812141dc1b2ab90ad3353e4b0c7c226bab5b0c06d39de7f52d043f0aea185d494ea9a84f657390c225133b2261ca66dc0c706d79c2ff3c84d5da37d8d8d8f55af546cf0d\",\"ProgramVersion\": 3330,\"Shares\": 15252235082420000000000,\"StakingBlockNum\": 2490123,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"efbd089457980cb916e6cd196591c673720882930f30a2bd0c11c1a86478c47bd719ad565de36f717989c9722ba5aa98534062daa300cb66ef6469cf52e2a955\",\"NodeAddress\": \"0x6edef90f9d18eac981a7d3718f17efa7d874c105\",\"BlsPubKey\": \"182bea09fe67e57eec605dbec12412ebd79c29e0a8d35e3a58400879589484c52293c854da8ecc1ae9acfe27e6898c02a2a500de1f0e48964ae6190e67b039a78f51f0fa5ed12e8f445fa8f76b4e6c6a81ba3c4ad2c0c8e3f37473d29d199a12\",\"ProgramVersion\": 3330,\"Shares\": 15220000000000000000000,\"StakingBlockNum\": 504836,\"StakingTxIndex\": 2,\"ValidatorTerm\": 0},{\"NodeId\": \"88d14bc287a42fec8a490b90258ddfbc18e98a32f42f9dbe27a4369b8e64a41e22ca84aba764ef7d1326b5b9a8e7b54d6d10f23dc665ff7e7c40c04a5f0a3703\",\"NodeAddress\": \"0x79ff4d0254310b0f8df61596388883c9eea59bce\",\"BlsPubKey\": \"d823cd10f3276e3438948f44ea05c1e0834847ec126532201eba4f5cb8103d2422ee7ff63e24a67c8bafd732ae1655131fa6b592ab7b3f3f7ae81b76822a738a915c50e146fea8bd68472861168816ee7e2de5ae8ec13f1130a0a77ebdb51989\",\"ProgramVersion\": 3330,\"Shares\": 15171661310620000000000,\"StakingBlockNum\": 500760,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"de2dd79123ef3d9c69e4bd9e331a9d9d8f693860e3a77b1c4ddc1fd4b40e2824977dea27a2311d5d8bc3fe8be70abfcc814b2498a17c008bc0ba13150ad962e1\",\"NodeAddress\": \"0x1ca296abf6bd1a3f1b1f6456057f10d9c755c28d\",\"BlsPubKey\": \"b9b999d372a9c3116ffefa9aae1cd91d1413b0bceaba912634d4ad1d2d2d8fed603170b5df3dc4d9ff87f95fe9a3b609f592f9feaf0c8149341b3f89c17ec285d647808db45fbb7375c5e9da014d7fe2af8f5511c7b5328663309012b733580b\",\"ProgramVersion\": 3330,\"Shares\": 15148979159857629376000,\"StakingBlockNum\": 613734,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"c537d6e1393521608e1dea636cf2f0e095d1e73e1b6151bfd0159b563fc9106f9c80118fce0b3d211cdcad2df0359a1596c3983e79bc26b8d19c8325e1869aaa\",\"NodeAddress\": \"0x028fde5b753ffca406fc2b4e3b60582d28911483\",\"BlsPubKey\": \"c1b6a3e3128024dcdfbcde9868d2b9d65a774ee28b50d1a843e1051f4f2189e13d702fedec9d8742191e71b0ed1fd8165ae0b5ec2769e31458c05c9772d419a86906cdb51fcfb576c4edc79cb914174cfecccf1d27c6e4f91ffd8794cd1e9292\",\"ProgramVersion\": 3330,\"Shares\": 15125016095720000000000,\"StakingBlockNum\": 1469141,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"218b20ee9cd80ea631c771996262a1d4b74b6db5c14321d55f42810a812b128c8c1e3bca0877e19e4ca12db6875d1922f11cfaeb9b2ec78da5157e523c7db8d3\",\"NodeAddress\": \"0x75ed15714ebe7ba6ea1c387c4c281b849f1ba840\",\"BlsPubKey\": \"f55a19110b073eb60124e1b240503546fc556b27d2ac58a5e752ec8b1045f6e7a28377f8736001b19fda7b33ad218004d750eee38f665668e3922ad69f2b1fa90ce5383413044208003e7a93e102a1f910b6be5e76e15f1ce09a4014c57e2989\",\"ProgramVersion\": 3330,\"Shares\": 15118907262572999991000,\"StakingBlockNum\": 503336,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"3015687818ccdeac78b66f2b15dfa4534da10abf7aab713267431114b4ede6fb6bc27f7d0cf8e792c356da6e916d81c80b65205679d2046beb9a69cfe8c374f7\",\"NodeAddress\": \"0x74f3a7e9f1884dea069b2047c9226d43c0de1ea3\",\"BlsPubKey\": \"31ed35c17d92ecf6466a75a3d18c92da5b33a5e31aac5f4a082c7d94dca53f248a67b19c56bec1f1deef100cf952b91908c8f6ab71659ea762f35417204d9aff1427064bae7d57d7f7e0ec7f32017f6f6d5076b8481eaf3fa88559018130de07\",\"ProgramVersion\": 3330,\"Shares\": 15092905507790000000000,\"StakingBlockNum\": 983558,\"StakingTxIndex\": 1,\"ValidatorTerm\": 0},{\"NodeId\": \"f7a07f6ff8c282a4a1c1febaf68531bbb557fa262c641d022f249e926e222c9190602b502d2772fa4a2834ffb4bf6298f4c80467178e90728c5a852504f571ad\",\"NodeAddress\": \"0x33b15e1c1e6b41d94a73a7a70154c3486e1387b6\",\"BlsPubKey\": \"c227209ce310a0c32eff51f32aab35642e494593d8803cfd23862b8e0b1761e16df2ee54414853dab5965af27e5fde0ee5695a0b585e67994e2d3eba66bfe000daf873c8cb817d2545dd959973912170f314d62e6d4c887f678174a2cba28614\",\"ProgramVersion\": 3330,\"Shares\": 15077886624870000000000,\"StakingBlockNum\": 502229,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"f2ec2830850a4e9dd48b358f908e1f22448cf5b0314750363acfd6a531edd2056237d39d182b92891a58e7d9862a43ee143a049167b7700914c41c726fad1399\",\"NodeAddress\": \"0xb2a8d0c2dcb25c97e008d9be84bfd2330f6d2e25\",\"BlsPubKey\": \"be0973f0e98091eb68e4a4251c4346e69a9116a05688efa44995beac8ea7208aedccf95cb7cac5ded9253ccdc5db280ad90bec9f2e4ce5c51f4fd79dbdede0c91c15b586ea899746fc5d5b4da9caab8c73ffff17f9e3eeda217f779cd14c0e11\",\"ProgramVersion\": 3330,\"Shares\": 15039115745970000000000,\"StakingBlockNum\": 502569,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"14c2b76bd5945f8f77da071e0d46b253adf765551919ef3e7f80885a7a5506bb70e34926c546d27490b5a58cf01cc3f079296b0b97a01daa694dcb94c97548bc\",\"NodeAddress\": \"0xc3e4b571cba4717c03e08d4472e5de5d289984f1\",\"BlsPubKey\": \"03c4108b21e026011b8a2ef47ca7e451f1e8d1453888f0139897094ccd24b6f0438040b0c3694548983489ce2ae5220d1072b01dd0e072c792839aae5ac2323338af3be58ccbdfcdb1ed34b812893e98046d7a55c0196b865ad05d28dbcf6c13\",\"ProgramVersion\": 3330,\"Shares\": 15033349194663807465856,\"StakingBlockNum\": 500966,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"46e97c2f1774a2e505c8629d4e5d654f7faca48dfe7c23e4ebbbf86655feeae502096da83ab1d6d3cc47da5ff852bcb45490cceda4bbe7055f6d7cd91d5d8640\",\"NodeAddress\": \"0x22ca0cf3d7e98ae5c7f936572e4c8c2c5aa98a58\",\"BlsPubKey\": \"272019eb1fa85f1e6523b6dfa5c2e9de3b2a0920bd1069717e5098d0cb14b70300463a32bdd890dde706bbb82ef5460b568f8ae86eaa4092f61ec54ba4aeaf06ab897c348edb94be5d43bb320a9c824240c014a2e1954a5ec63f4c91cf83bc8c\",\"ProgramVersion\": 3330,\"Shares\": 15003000000000000000000,\"StakingBlockNum\": 503384,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"0525a2e651701ed46764b4648fe88a64c70ecbf69d6d897adf48f549772d9d67840288a1396f82d6d6fc7b3ca6adf0778995935e445a3e7b682cc3c41ccd83ff\",\"NodeAddress\": \"0xb6ea326a2960b7dfb3b21e209fd7eb0557aeecc3\",\"BlsPubKey\": \"2713752711e4ef95b6ffb5df76efe8f576665cb457af62d820f983edad92292c88dfd3f83484fb67ec2decb7294427060bbdaa9dc96ac6de576138ad1b145b39ca2ef34d6c4c327fa33da7daa89dab03e49b2de2b045e50179f10a029d41978d\",\"ProgramVersion\": 3330,\"Shares\": 15000000000000000000000,\"StakingBlockNum\": 503026,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"25f4caa5654e4b132de6ada9a1b8a6c474536566f481a9dc015c2ff75b1dd79e99c81a420ae5243a68b68eb62b468f2d1c5349eef9f123e6ed4f186efd8d1d09\",\"NodeAddress\": \"0xa551760aaa02c02114f074b8b99cbcb985cf9ce6\",\"BlsPubKey\": \"83a6dec9fdb06ab195501a19b3513b04a5e129e6b84152ae130db1a9792c30f6455c840f5bce15fbe927181b086a9914b53a79dea39b720e760a34dee8b4be9c86e94618cebc1119b6b2cea91f74e491e8953a9d6785d811301eaeecf48dab13\",\"ProgramVersion\": 3330,\"Shares\": 14972069781100000000000,\"StakingBlockNum\": 554287,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"9839ced934daefca276eb8945c978e97774e9bfe8a101104b3202e3b38229fa4c71a06145366b4b89a135d89002ffd5bbbd12f9da33bc7c20026b290e45fa5c8\",\"NodeAddress\": \"0x3476f78a1191b80f650f7eba5d7da9afd0d379af\",\"BlsPubKey\": \"ad4dd3d0283628f2ab20e83a3ea35bbee4eac240849ea3d44655a8b3c0ac5aeb7dac751c3902c44a4d0a14a1e04abb0a99e66f5587607a50a557f0be47bac44229c1d470e9e483445f472c75a096605305a6bbbdac6caa8c60ac3520148a4d8f\",\"ProgramVersion\": 3330,\"Shares\": 14958456447540000000000,\"StakingBlockNum\": 1383401,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"86c63b28fedc3364b25c4864dc46e158f3948a87abaccc1fea686a7405ca17ed90a33dd73bdc0a0585d19de5855f3281d06b8f4f594b1cabd97efcc7039522fd\",\"NodeAddress\": \"0x9829fdf89238c1a524723f6293eff6ab822aecd3\",\"BlsPubKey\": \"4cd9c85c24ddaf02674a9be8012b511f8b3973ac8b7029d646d8a168a7540582559cd37750b851445ad8ff0358ca7f15dac7b82ce3b2873849be7800afffbd06db08e833c29883e3bdd54794c2d62422ab1217213040d160f86663334ed10a13\",\"ProgramVersion\": 3330,\"Shares\": 14944878393530000000000,\"StakingBlockNum\": 504837,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"3c8ed26005e0128371bdbd6e681e1e668eeef08a92a196690675d19774b2ca21adde544c4af740db6baaac50dd0a77aa48f2e99198d7ca807562aa3de7737ca7\",\"NodeAddress\": \"0x4dd6284df3d576ef586992c281bb45887c1ba4c7\",\"BlsPubKey\": \"8d0c44d713852bbd7c54d673f0448e71d319ab2ecd0f20459462df23a3a25d7c1f99a8aabc2207c8f62c5fb26a91c9054cac3c47b0ab4c64bc91f996805f0dc7902544589d76cff58e0b5b01790e3e51d635bcd4e525d3a1b46a2138c605b598\",\"ProgramVersion\": 3330,\"Shares\": 14933604050620000000000,\"StakingBlockNum\": 502144,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"bcb2be51df673ae15820d8ac083ab565c78408f79e372a5858e11fc23f1f711890e1402459fab088bd2beb95729fad5c6dbe681b3ebfc3af143b25afc14fcb59\",\"NodeAddress\": \"0x9729196bf494c30e77b2e5604c2334682092c10e\",\"BlsPubKey\": \"8d8c502969f02a9adbd3ab83532297f844e73353226944f2bb938729876e231d77abd8672455b20aab88799088f3540f4d42a2999837da8d63527bfe47ac47d5bcab72e2a68ef22efb7dcb7a95d17a49592e675f0bc589ee8a44aa972733e392\",\"ProgramVersion\": 3330,\"Shares\": 14908918887860000000000,\"StakingBlockNum\": 504069,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"370248c1d479da63476c52c742decc6b04b92908cfeccc53bf3a3367971f8de08a270a7e756f4714e99d0d10521426a4270335233ce5ebe8029828f35707c33a\",\"NodeAddress\": \"0xc05dd74cf6fe0226c6b99be0b7d711c6cf38e30b\",\"BlsPubKey\": \"7f3e41b7bf2eb00d6fd04ec708570408adc92f622a1b0481b60e0cc1b6f6b1d3b333382873166b6e2bd2142e09288e05d5e9a5f6f7ab9d58b0c8144a2c60b5f7a40a0ee5f3eb01d2d8a2d9e6a4a23bc8bbde483fd8dabca8b5f714ce631b8f99\",\"ProgramVersion\": 3330,\"Shares\": 14905149739680000000000,\"StakingBlockNum\": 499903,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"1d4b0709c6bc493af2086c3d38daf22ab64a0e992704a571af7479b0b58292d0bbab33552095082e3aa9297a8f338bd2679345e3f02bdefbc68c80993e0800e4\",\"NodeAddress\": \"0x7dd99e0a434b899d24f98a68a5ed7bc3529e00eb\",\"BlsPubKey\": \"9a8a2aacd60cc5e5d5f5db1eaddeabbfc1e14facd0fbc31e62d4a055217bf62e2ee82bf8f512a23156a7f3940aa40f0c5da085de11e7743b419d9986ba18ead44650209da84d4485471c6cdc1d8802f16a6c2d240078dd8ffe0a0e3637b90a8c\",\"ProgramVersion\": 3330,\"Shares\": 14881211330400000000000,\"StakingBlockNum\": 504836,\"StakingTxIndex\": 14,\"ValidatorTerm\": 0},{\"NodeId\": \"2f65f95d81071a43578a8c73a293952af05228505adcfb6645822af027c8f6a0d54a54dd8b51fe901fd774fd94f0c9bc666c02841af9986cb694834294f89494\",\"NodeAddress\": \"0x4a9e6af6a50944625127be66ad044362a694897e\",\"BlsPubKey\": \"afb0444bec39e3d268dc4c11ba08096e7ade2546d8cb5328182d862ec93b4dcff16b1c3ff728700b13cf80cebfce3611a8d7bda744f800e8414cdbea53619479d2f7780bd8a2a24bc07bcbd36c693305eb06009e9f140f1a806fabbecab99f8d\",\"ProgramVersion\": 3330,\"Shares\": 14866200650459999999694,\"StakingBlockNum\": 972709,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"23d1bee8744c6d63bbd02c2b0a1685b215762dc934d1fe5e3b88432e04edfaaaac78f8e5c52d46adee9707c14380971ccd30a5c4e544f8056f793ac328927939\",\"NodeAddress\": \"0x367653f065eeb853fe65e52fd0ca6dcbbc85928f\",\"BlsPubKey\": \"5b2267892ebcbac7a4fb47a8ff17fbe6cc304df0839ca8f7bd3a36662013917d6d64223c8f01d7ce430fd1789303bd09f6d52b256f73b14205fd22061f6cde1278beaf82f643bf8989b92b214589bd75c4ba941f53cfcd892a426d977df03501\",\"ProgramVersion\": 3330,\"Shares\": 14856252405540000000000,\"StakingBlockNum\": 502248,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"ade53a03334eda53f0bed8f3e1572a7d5faf59f40e464b683e9160946a738b0d72b333b93c07b4b75341de633784b4a2d00aacf78ce082949470046a4875fbee\",\"NodeAddress\": \"0xff7514255dbc45c374b63580ce49c966e9d5ca6e\",\"BlsPubKey\": \"53a4a81b7a70e6a81fb57c9d929e48a32118b767fd1fdb9d411ab64626d60573af7417ec63fe9bfdbf232ab08f128707857e39f181a0429d13489777e0d511a9e123d24c0357b85660473679f0897f51b5dd139f96fa678cd92d1ecc6c9ec783\",\"ProgramVersion\": 3330,\"Shares\": 14835851767100000000000,\"StakingBlockNum\": 502590,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"4112b18659726c932cd724ff53965fae2d6a68f990f8784250bd608981464025b368df7115cc403d88e1f973c6c12b52903d2c278b97887650b1740e755825cf\",\"NodeAddress\": \"0xb39d60a19103cca815864aea01ffb23083c5a352\",\"BlsPubKey\": \"ed37aba54abfa5afe5326e3d8f5434d47ef612cf17b17bbecaec8d53af518c3a43b9af3d05cd5cbb1bdd7a8d3bf1860534e8c5d69d2c4b6b2ab67159f572dde799c530d2d5d37bd4c28a2bf4034cabc9ee62454a95f887f8ea25bcaaded28111\",\"ProgramVersion\": 3330,\"Shares\": 14831989410729999995000,\"StakingBlockNum\": 575054,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"80b7f1c6b67a79079948fa270b340f59ec0c346623fabb44db7a62794ca0ecdac06b07721b8c94a3c6f34a4b83a9010c686ed64925e23e98f1c49455651644cd\",\"NodeAddress\": \"0xafb203068da4046e85cf4098dd1d2dbb0bd82a2e\",\"BlsPubKey\": \"a19b6a1be79d99cb763bc344daa30e89cd06f9dea3ec269ce85fe0455899105230b4b2b9d8a472088160d9f6e0fae207d2c14f62eae0a1eb315b3ca26fa59e12b20f37508635f8c1f9aacd381bc67b167ebe33fcedc0e3b3c48401ca01cad60f\",\"ProgramVersion\": 3330,\"Shares\": 14830252405540000000000,\"StakingBlockNum\": 503131,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"cbf8b3b0bd51bc4cd3edb55a5cbd297e7dafcdfa17ffe82d1f4d2d5f1be1782030f6ac26b9f3910337acb8fcc356f79233851a35836f4197b3d4c8038c27cbfa\",\"NodeAddress\": \"0xd5a7742661e629d5110f7061f6cb9ef6f986596f\",\"BlsPubKey\": \"c89bde1c013248fc775bc61de6df53587b1f0061595240029ca77e53cb064e2931c94b7bf2ce073ca16c932604f11b1339b5933e9fe18c46bba17b06732e21b94658a1190c5a812ae13fe94e322c6a2d0bff0c7e9bdfe671358e2afe6dd09d95\",\"ProgramVersion\": 3330,\"Shares\": 14825914273900000000000,\"StakingBlockNum\": 504836,\"StakingTxIndex\": 3,\"ValidatorTerm\": 0},{\"NodeId\": \"eb1a3c276744580b50d478e7910debcc7af7b3807e58ca23b7542a8265d996ebd5037ec9e6ba36ee8528aec0465bf375f32acf3538f75d79a11ab501bfff11b1\",\"NodeAddress\": \"0xea6096fc1d7dc4c8863c006c83d104eeae466572\",\"BlsPubKey\": \"ef948881dee485518539adc51a6c79042657d1ebaf8471af45efe331634358d4744c2591e8f811ee8c18ec41f9334f08169da43e618199fa1aef3d752178150223bfde126a1c7aa7a5c9ddf1570208cd405600d3a6ea884147f90d19fc23c102\",\"ProgramVersion\": 3330,\"Shares\": 14820000000000000000000,\"StakingBlockNum\": 519599,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"d839188ce67070dbc2bdd2a93a8a9abacf1ee89432ab18feb8de668bd595312f2483be56197d15e977770f59560cee5c03301a6c3cb90415951d0de26c2bbde8\",\"NodeAddress\": \"0xbb30b39be75bad60318cef7fbc6ba7146008440d\",\"BlsPubKey\": \"2b099ba28403f8bc2878252db4c2e7a08699f4d01cf0ccd10b5e351c8fee12692988d871166d10db613605ec19bbf00b59a98a29518a28b5a37261c730181837e21fe84a7d2fd69fc1883113c12e342eca05074aa45abe16220b2a564e195e11\",\"ProgramVersion\": 3330,\"Shares\": 14808976886010000000000,\"StakingBlockNum\": 504836,\"StakingTxIndex\": 1,\"ValidatorTerm\": 0},{\"NodeId\": \"94795e21b23fa28473fd93d17edfaf24d5a0b9ec0c5dc7626c4e8882d43548d02e11aeaa16191ddf5128cddc6d0bb352f461570049a765c3a19e86d9771b3173\",\"NodeAddress\": \"0x666915b57c88ecda79a5602d1ecba4af84278dc1\",\"BlsPubKey\": \"473596c65a32b46a8b2ef0b509edfa335e9e836ada4d0150150ad4ceb86d7da45a8ac78ff4e11e53bd1f30fe1beb5011a4ca2d523339693c08b6d43710b8e5439ff0e5c8b568fbd9bea9a041cceb7a0c6505e1c0ab49fbe32095b2c376a23e0a\",\"ProgramVersion\": 3330,\"Shares\": 14795015642285693531400,\"StakingBlockNum\": 501309,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"0178bdb413b9acc82b80fc1d31ef4e66065708c2ef8758bf5121100049273dd9b41fb16c64bc747d732575b959a45400db757a417ac0d6dffe220ce68827eee5\",\"NodeAddress\": \"0x41ee089643a082255f448f9bfeacc7345e908c39\",\"BlsPubKey\": \"bd56ed0587790922250e57f049236b0eee22ebce130b8eac32a383c38892e9a237536ed8fef3f2e3be518f038256a9174d99c7681488b1b242c8e7f387c97b4b92ba351f30e6fab776cc474b02d8894db7959ef2ab16ba12c600371b70119f09\",\"ProgramVersion\": 3330,\"Shares\": 14795000000000000000000,\"StakingBlockNum\": 984400,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"24b9584429dcb52ef6da575a0246d59158ff04f37118e9ce580f100e9da4a99064db252648f78497cf4b27f53eeaace7ca795ff75734e0e95386a5e3282f5fff\",\"NodeAddress\": \"0xfc5208d29bf4ba3c7ae50b55072e01c178c02e49\",\"BlsPubKey\": \"ac7a1ee46b906e7f7456d16a707b14b9118e7ede78bbe99bccad6b9bb06fb750fe41a3a3ba2447f734b9b5aaff7b9f0f10d0f1ea1a90e0f45716b2b9f0a673658e8553230e6958161f919db29b13cfb46d34576ae2bace71a0446486eab8920d\",\"ProgramVersion\": 3330,\"Shares\": 14790000000000000000000,\"StakingBlockNum\": 504066,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"a98f15ecc908e6ce68b7fe29ea56c8c552f09b658352d0d0fb0fc2f08aa50b186cac9468a2742b9be9ec5d6ed0840d157afa88cb51ed2efacfd51cea26a7aa84\",\"NodeAddress\": \"0x5d4a4cbf68f65c1eea4d18804eb59c9661a09959\",\"BlsPubKey\": \"08f278dc94e3a4988f4c136a72ac20d9a51bec60b831e73a45cac791cd6f975fbb2ba2dda5996f1457ad3570310ce60f8752fe84d0c70d53b8fe290ea0f09fa0193cacb10e36d7d10129d33d696cbbd22e3e64128cb2e7e418a40ef9d2df5687\",\"ProgramVersion\": 3330,\"Shares\": 14776346005420000000000,\"StakingBlockNum\": 502133,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"32b67023a9057c66ff9b9e0f9ca5a23a09e686a707c303695219ef0aa82ca2c7bd965da5cb1262232376584ed4402877fe786b58c3b1adc8fafbb3174b3f24f3\",\"NodeAddress\": \"0x00d2abf51b8fe6991324ef046fc559db1a586a0f\",\"BlsPubKey\": \"69671b96dc6e165e196deedc21c222d90b2f9c260dadf4e1057b7f55760e8e75bb2332b8c04fb1f03c4b3b68bfb3cf0f902f686a7d9e06f1e302ea054805b70d4101271cfc0d3976dead9edf95804329284c871c464cae6be361d90201e43b8a\",\"ProgramVersion\": 3330,\"Shares\": 14776314253260000000000,\"StakingBlockNum\": 502244,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"1fd9fd7d9c31dad117384c7cc2f223a9e76f7aa81e30f38f030f24212be3fa20ca1f067d878a8ae97deb89b81efbd0542b3880cbd428b4ffae494fcd2c31834b\",\"NodeAddress\": \"0xf1f36cdf690165d8e46f0b63cf90724df44b8627\",\"BlsPubKey\": \"d54a2a1d9c0b7906aa39c0c96652017290f2bcd1a6e2fd0599b9d0c02e7cfefb529119da2cf80232aed58b3e426398072f2d72dc16b29bef3b07a2495b60f1d84f8daf443b0844a9239f35e9bc65ef69afd4ac891ab5bebc45cdd401cbf21f94\",\"ProgramVersion\": 3330,\"Shares\": 14763820267986000000000,\"StakingBlockNum\": 518839,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"949e3793be67761380ae1536832cd5d310a555ded62d8f049636c6292815d27340188a679d5bbcf66d48fcbe459c05cf61cf475d168ec059b7e4a5ee22b64ba7\",\"NodeAddress\": \"0xaed50917444c4ab908cc9104af3ae1951c182aff\",\"BlsPubKey\": \"66baa611920643bd26cdb2cf8b77bc0584e85083ca540db316170362e7d16e8100732c128bdead8eb31afe759d49700bafc183b61ce3fcb079ef7792715013b1f94ae082e8a85e4c70525397b9d706b62e1e7b8299e8c040b3a6998d7815b398\",\"ProgramVersion\": 3330,\"Shares\": 14749000000000000000000,\"StakingBlockNum\": 504054,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"8f2dc504635ff5b394c3bf7debd885db3f4108772d8e51c26a3860a97d34295ec5c15cffe9fdccd88052408f3a84208718291b077e257898e603ec7b195d3684\",\"NodeAddress\": \"0xad327a189acc7be561ef2dc02a4a60fce741c362\",\"BlsPubKey\": \"5180f5d8c6975613e3f9b1d3fb23c9a982d208936d941591acedb47fddbd301ba5563b2b3e632ef50a3751dac5a06e0ec89233ca403c56d7095f740ad0df4a5bf82e5dee0fa51178d931b7b55b875302d3c79471f9ee522fa4b1439852a29d87\",\"ProgramVersion\": 3330,\"Shares\": 14745056426090000000000,\"StakingBlockNum\": 972640,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"a71e63927d7c96421448a9e79833d3c4a3d9de323b9d09c93a33da4ec2e7239efac29735e5298b21c23f7cfaacdec949060bd34c6afe04fc75b6501504058756\",\"NodeAddress\": \"0x9904290822ecb241c66e412049b6dee69eccf378\",\"BlsPubKey\": \"ffee961e1fef9122aa275413ed12cdf2a772836a4bacc37929f8c71fbb5e0c305dbfb05e6e2585209c17e58fbe6ff7024430c3a675c8dceaa9f2524476868640f1d24f8d58a5371763e779a378d0c59851aa6220395f6ef7ba19f8ed59454d8b\",\"ProgramVersion\": 3330,\"Shares\": 14706417656450000000000,\"StakingBlockNum\": 538769,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"89ca7ccb7fab8e4c8b1b24c747670757b9ef1b3b7631f64e6ea6b469c5936c501fcdcfa7fef2a77521072162c1fc0f8a1663899d31ebb1bc7d00678634ef746c\",\"NodeAddress\": \"0x800dc7e382ea1266ffd1389b3e2868ef91257b76\",\"BlsPubKey\": \"26a9c78012452132ae0fabd3a0972e5faeb1b6eda24b809cf4f14f8f6d316732538617a90cb88baada66303d98a0360f2e4a8d35678ad1c4628bff41708a89df4d7dd3b28076cb20bd3e5488fe4992083a27be19a1d67a21196ce961f7c37685\",\"ProgramVersion\": 3330,\"Shares\": 14703608433819819996000,\"StakingBlockNum\": 522944,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"c92676ec311d571b92a7184237e8e7b7acc46c7b66c922d11a3547cda42693c76ae183ced5dcb6259eca20fe45fd732e0c8d7ad14e834f94efba54bd48de0627\",\"NodeAddress\": \"0x8024ae1213389d17bf56cc9c2025256b969cd46d\",\"BlsPubKey\": \"4590ab46498291af88072f614562c137503faa0c8eb6fd61dc92f8d0b4a30916fd9856cea05923b924cc1f5fc6d0d3008b6b029e2c2775b552b4209c312e04e00390c9720226afd4a0757ea24e10e9a690793c1bb71cba849ae19c38f7fad691\",\"ProgramVersion\": 3330,\"Shares\": 14701224012945624999200,\"StakingBlockNum\": 810539,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"c0e1be37b517f3bc7a946b2fe6b01caf815f4f6d52fadadd47b98aaf7cb19f974f7d5a7e05180f88442dad9db1b413e045df0d94ccda8eeb6873151f056ca0a1\",\"NodeAddress\": \"0x15aabb2a881de0f7ccf4f7c9c995f52d9f5d03ea\",\"BlsPubKey\": \"4af66a69714602897e4026f38c976be639cca5adc2c5256e8a71271e6f24dc0575afea7f8f2e008e26adaebccc82cd095712558975915d93e923c0040f96413be1b3e3cfb967d6adb802134c24ee8378104786a47568a5614bcc3f3fd5d51b8f\",\"ProgramVersion\": 3330,\"Shares\": 14688577477580000000000,\"StakingBlockNum\": 502183,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"6e409559d6b93b01e330400c8ca56a26ce979fef23edf4460e450975d0755dda5d44bcb4632871bffd2447d1de2b43f7639b8742985ab2cb52d984518c7cfefd\",\"NodeAddress\": \"0x3ae46598de92272f48fa3e71466e48afbd245af1\",\"BlsPubKey\": \"fbc2aeee173fd5af5ee1324ac501177248218f7a7ce01bcabd4b3a0b9574a451528b7e33873110791deb9156ba804207d9d5c79f8cdeba18e1a89f38151a1f92b887de1d21d90c0feb71f6c395886b0444e8c6136088128ea08996fca4ab7891\",\"ProgramVersion\": 3330,\"Shares\": 14685427091500000000000,\"StakingBlockNum\": 2068278,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"3e05f80d931922509c7f65ab647b39efec56c95360fe293413368b8e5db7cfb385ab7f8a77acd6a8c39ae6a0d108e0723cc84c19c8797ecd0af1c5c6e94a541b\",\"NodeAddress\": \"0x16af3a7c125307b8b1e7277b1ca413a6982aa67b\",\"BlsPubKey\": \"376726520bd478921a4320ed9a61909ae6899ef8da3f38ff306f4aac2721eed6c615e516c131e355cb08fd08d2d3f1187597264c01d8ae808493ed48f7cb39e24661fb83c18885e552e453caf01306ecbd98e1503cb21ec93caa2f10bb35d918\",\"ProgramVersion\": 3330,\"Shares\": 14681282295000000000000,\"StakingBlockNum\": 503501,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"8bc8734315acf2af4c92a458f077f1f8c96f0530fb43510c11361f1d6469631423206ef76cd879ade849ee15fbcaeb042e3721168614b4fad4eecd60a6aa3e94\",\"NodeAddress\": \"0xd2928cb4baa9356550c945a7f4d3a4858163afb3\",\"BlsPubKey\": \"8a3406a8e447337cdf5aebdd0a6600dd54142713f3adc735821ed8bf59e0f040973ed2e22d4448a574270bae3b3a9f116ff10e8a4d55b99c4e2c5449a82a94294359f6d93a8c81b5845cf210c5a91332f9ba51554db170b4bb817ea907173988\",\"ProgramVersion\": 3330,\"Shares\": 14676475832270000000000,\"StakingBlockNum\": 618133,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"ab74f5500dd35497ce09b2dc92a3da26ea371dd9f6d438559b6e19c8f1622ee630951b510cb370aca8267f9bb9a9108bc532ec48dd077474cb79a48122f2ab03\",\"NodeAddress\": \"0x6267901551f74a72ad3f7465c6c676a0ced7ba50\",\"BlsPubKey\": \"a74daa9361921a6d42eb8c150976efea2dab4a56d5177c9e4834985ce6e85270870c158bd47a67381d832ce2fe2b5c035c807d06a9468ed4cd7fef3cce3efbec0a0f40019cbe541266281c4fc86cbde3e7ae52932987328ea1313f754d33d504\",\"ProgramVersion\": 3330,\"Shares\": 14673785808460000000000,\"StakingBlockNum\": 507203,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"a55b37d341ec17a34081e45cf7289ca9d6ab1152f78ef911adb4a2e4f2700dd935175b6193183c9bc8dd3bbe352f5b5013cbefe692cb046c70dae551683be27d\",\"NodeAddress\": \"0x7f8a3f386ac3e7fcd284552bcfc84ce15939523b\",\"BlsPubKey\": \"2b96ee52937f8c2c798fb7fa8477338b6a774e62a9546ebb12b9f3ce9da4222c32c082488af284c20e1ec3ea32b01e057e722669136f1a51d03b24ddc0659fd79aa225be28ce230f7836260a1aa8a97f79f10095a0ce839c414f9ebc2e01d508\",\"ProgramVersion\": 3330,\"Shares\": 14662986158310000000000,\"StakingBlockNum\": 504836,\"StakingTxIndex\": 10,\"ValidatorTerm\": 0},{\"NodeId\": \"76dbe6b89833676301c3840a866e85e773a20da92bb2b02cfb9eae1ba00f38b898d708b97875677cbc979b8ad39bd76ddd9a8562a7fd498fe1def5520698f0ce\",\"NodeAddress\": \"0xafb3ca19a77fd89c1d7b27b2d83d90f6987144f7\",\"BlsPubKey\": \"9aad3d261150750165725e722ff02a8871b33547976476695ad340a1a23c439686be4bb7749725425defda505e2552150a3b48069f7926cfab830f6feacacd8196a6f9e0c2b650a9669b1742f492ed50c542dae4f1bfb597aa19ee4cb5830902\",\"ProgramVersion\": 3330,\"Shares\": 14662311488510000000000,\"StakingBlockNum\": 504188,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"88e70a87f6acc8edf3b381c02f3c3317392e458af688920bbfe04e3694979847e25d59fb7fe2c1d3487f1ae5a7876fbcefabe06f722dfa28a83f3ca4853c4254\",\"NodeAddress\": \"0xd3ea2a9cb22cb174638721e1fe3a8881a92605d6\",\"BlsPubKey\": \"8fb4214835d0f408cf6849963b3e939b7a91f8419b3d0dc9ec98c566fe30ff12c4f8b8ce34a3241817f54c1bfdf19f0102bbbb4ebd6ba3303104e7fd88865d705be176161be8806d41411046792d8b5ca995137c173e8c716a37dbeab60e9819\",\"ProgramVersion\": 3330,\"Shares\": 14660989458190459243400,\"StakingBlockNum\": 500469,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"94bdbf207f6390354debfc2b3ff30ed101bc52d339e8310b2e2a1dd235cb9d40d27d65013c04030e05bbaceeea635dfdbfbbb47c683d29c6751f2bb3159e6abd\",\"NodeAddress\": \"0x11e0347cb54198c1f95edf17987272da112065c2\",\"BlsPubKey\": \"71831b8ce8836f3e1cda27e7ec4e39310008a3c1e68e3cc2f54085839e7b792b5f74ea5758b9e7e97fbc9b320179c910eaf0f041c4e43358292cd31710b7add54ce6d5996a8143e24f7f681b994b335df83ae91ca3dd53d271ff77ee409f1387\",\"ProgramVersion\": 3330,\"Shares\": 14653595699100000000000,\"StakingBlockNum\": 557309,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"fad2c7f917eb3057d85031eae8bbda52541b527dd1d24a25e7e9b40d7329570a85dc45ec61b189a9cc30047ae906a08dc375558828e1c76dc853ce99b42b91e4\",\"NodeAddress\": \"0xbbab9a89fd953205b90ffe3c0e6e47e84eecf6b1\",\"BlsPubKey\": \"795286e5120519bb718851b0a608addadebb6adc24424281a02e5058f43728a372fe3ab0a288f59d5670c553787f0912b562349a01fc66e884a17b828315f9f0d55db551a7d249d02b4691581580f8fbf790b3103f13676fdaa3eac3339efe99\",\"ProgramVersion\": 3330,\"Shares\": 14589247081246790319232,\"StakingBlockNum\": 509677,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"c45a914bba5ad5c735489a100954eb08c49c3f718879e483d0debbe3e83da409cb1f9100ac73a2951bdf38886db419d8648534698cad05ab99c27d6869f62582\",\"NodeAddress\": \"0x582140b4c1e369a6ed553fa0bf7aae816c8a36f4\",\"BlsPubKey\": \"7c6345ff519fa260d343d62dfcf33c1e0702526a804d629622d66b38748ff157e4b7a204e8d39d6687ef5d72943e2419a9c084dcac9264f71e53b81850a1d8a69377dfd359f55978877d00cdaafb2ec3207178e95164840eea69b719cd0b8585\",\"ProgramVersion\": 3330,\"Shares\": 14529194116773044337152,\"StakingBlockNum\": 502421,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"8fbc190c20707087eac8dbcbffa28a489e5a6c555c3b9542caa60307af5f08a2faa5db23b8ede13691d52cad2a1075432f9a754673f34f60e274d66742475a3a\",\"NodeAddress\": \"0x374c022dd6104b23eda6b520698abca63cda52f2\",\"BlsPubKey\": \"dfdf58335e4cbc177cc86a4a776d67e555a9ebf366585901a41f9aea02ba10cd4ed79f6c1d3e84d6624f07c103a87b0b11977e50513a33d59d234d508879a5d7011111d4b1661c3606fb348e1744060a95b772b464d6d3271e73f87ef4227e00\",\"ProgramVersion\": 3330,\"Shares\": 14430436831746666666219,\"StakingBlockNum\": 553698,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"3d5b3b555e498b9a67d2c81de054263c7e9d4be7f1adc325b3413c5b59b1c165598379ca51b15c48b4203761c69c986de13d782bc12599dcf748e63f88f481f9\",\"NodeAddress\": \"0xfdd70ba9fce8eeb450dfa578861fa13719441497\",\"BlsPubKey\": \"1d7d8356ee3f11994b92d397bc119c77262574dd858344c61b249dd2e17807ede8b59a46bb51a22062165f35842f4b0cabfb7c1e62d649c01b1da65162676714c70bc59e5fb38efc4ccc43ca5975638db8882b3788cf8c04ae1219b79836a80a\",\"ProgramVersion\": 3330,\"Shares\": 14344304953605202663424,\"StakingBlockNum\": 502582,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"b5dfaeec9aa9114a182bdb3e2805f89a9e30fe5e3fb942ff0ea76dd92a72ecc70df74ce9c52e18f2a5ef696a4676516b13c69f154e91a6cc1b177c5eb8968750\",\"NodeAddress\": \"0x32d2b4dc4873c6ba657e8f233a14ac3321af00db\",\"BlsPubKey\": \"5f78661c73e25a7139e929f59290e92b1721bec465004a65c84b1e44bbd4a662172a439d3a7d0627207c598cfca6021703e82dd9ae03025f18de4b349ec35ba87828e863f92133002ad43f247b58ea2f7b2cb75be561d66b25c351342290f188\",\"ProgramVersion\": 3330,\"Shares\": 14318898655590000000000,\"StakingBlockNum\": 613769,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"505ef930162c1b0736d7b7a44a52e0e179597bd3a43e13b057734abffb46978311ef10bf16d565aad1ed4f208493b030c173c7bcfbff989d7d4bf2c37925f621\",\"NodeAddress\": \"0x63e459e814e8cf13c27c47bf41ba6655c1f0f4ee\",\"BlsPubKey\": \"b45055fe5d39191e09988a4154e48b8733ef95a9c6d9f90aacb1c91c2e98966948c9f550a15050cd87a7cf45a39e3301714c7e5df95aaa3d594c1a2043cc6323cb17069f03c1a8652b65d3c2eeca82f800782cbe6fc4df59f11af07af933f887\",\"ProgramVersion\": 3330,\"Shares\": 14246809530340000000000,\"StakingBlockNum\": 515799,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"3a8b3b9ee95e46f63025236e3c46eb9c87e2a8169cd681233f98b2646ac873a8ac3f6c5082f033ce1103936ca50fb98ef990582e85f3b95004dff51a2d6f2e76\",\"NodeAddress\": \"0xfd83afe15bd57d5a95cfef79d1a3de1f65388d5f\",\"BlsPubKey\": \"79ab4e31516a929ae5b09126633854aa37e6a36e14005a81703fc9d3559ab096913ec4bbb4ef838e8ce9d0b57846eb07b5424ebe46f11b987dd19b8d4b327abe65d93a860228aca2d1b821977ee0bd2d5628c8a13e61db9e032b44597a403382\",\"ProgramVersion\": 3330,\"Shares\": 13947737500400000000000,\"StakingBlockNum\": 515647,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"d94b1914a4e187acd6b6038d99a48fea21cd8208d2aff92fdca28be9731772aefb5ef469d9f5b8a2221161be761337fa54f7ea8edb2f675ea22c9f278827fc39\",\"NodeAddress\": \"0x99796ffa7091d72cd08ea7d3e7c92b8f38495f6d\",\"BlsPubKey\": \"0aaa6134088f848e5122225a6a8693430884348ab02b65b10dedee62971022dc0933c97126249875462016850a5e4805f6bfdb275c7f844229f63c3348954633d5a59c734af41b005822f98602d5c3b4e5b7910b05846f72d6d2ec1a2a906f85\",\"ProgramVersion\": 3330,\"Shares\": 13698864283620000000000,\"StakingBlockNum\": 501500,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"854c6554239251092bc7de885b586042760a34584a686a4b360e171e7e27fab23ab5d6c8faa2519b1e5af6d9a2009e234b0f6261a0daf4a6987942a72c23783e\",\"NodeAddress\": \"0xfff0a254288be5379a5d9c8d122e25993eedaafb\",\"BlsPubKey\": \"e42a2d6e21a905c74a54e7d9a25be4ce750da5bda62619b1de14e3d203cdae2fce1f89ba407f72399797170bd2dead07dd0279e420ee66546cadf85ab0a5ad4994c66a4ed7628844cdee4306d0814f52fc8fd72b373b90eb20708aa8ef87070e\",\"ProgramVersion\": 3330,\"Shares\": 13382891431110000000000,\"StakingBlockNum\": 506622,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"e2053e04f95afa5c8378677de62212edb972e21b40421786c53de57141853ce870481a80b68a449903479751da114a30d27a568812e178b919bd1e9b2c82f92e\",\"NodeAddress\": \"0xd6277094ee34fb92ac8d98fe4dbf265234e70441\",\"BlsPubKey\": \"5d536d08895d92b17bfe5b7fcaf97a7ed92cbd26ee945d62023e6fb941a3ba9379ebb38f6c55eddcb4939b1565e7230838404eb48127193bdbf9ef4102dc6ccff7917da658a5c8ebc3e6473ee42a6c60a6622b162441b7ff4a314d2e39929c09\",\"ProgramVersion\": 3330,\"Shares\": 12048528240959778459250,\"StakingBlockNum\": 508822,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0},{\"NodeId\": \"29c5d9c702b1dbd69bb5ae76aa396b0b4b60fd1f650717c01fe677494f3625ceb14df406624f46dccd5f788f890e4b85e51a529068877e5c5b772ef1c0c93d7f\",\"NodeAddress\": \"0x464c91ec2b2612d687a8c80ca6e96f4b805f1b03\",\"BlsPubKey\": \"ba2d324d7f5e0ff94002d5d12d6f6b5b632203cc011cbf426dae4bb76ef3823f6e970b81897abfcf10845dfdab53d90003a36945ad5ecebc4e6f3dca6f44e70907127e83698ab17be24274b822ece6dbd07c90798d80954024d36775fc060f19\",\"ProgramVersion\": 3330,\"Shares\": 11838817196570000000000,\"StakingBlockNum\": 502474,\"StakingTxIndex\": 0,\"ValidatorTerm\": 0}]}"
	var validatorArr *staking.ValidatorArray
	err := json.Unmarshal([]byte(verList), &validatorArr)
	if err != nil {
		panic(err)
	}
	preNonces := make([][]byte, 0)
	for i, v := range validatorArr.Arr {
		//fmt.Println("nodeId", v.NodeId.TerminalString(), "weight", v.Shares, "stakingNumber", v.StakingBlockNum)
		v.Shares = new(big.Int).Mul(v.Shares, new(big.Int).SetUint64(10))
		fmt.Println(fmt.Sprintf("%v\t%v\t%v", hex.EncodeToString(v.NodeAddress.Bytes()), new(big.Int).Div(v.Shares, new(big.Int).SetInt64(1e18)).Text(10), 0))
		preNonces = append(preNonces, crypto.Keccak256(common.Int64ToBytes(time.Now().UnixNano() + int64(i)))[:])
	}
	probabilityElection(validatorArr.Arr, 8, crypto.Keccak256([]byte(string("nonce"))), preNonces, 1, params.GenesisVersion)
}

func TestStakingPlugin_RandomOrderValidatorQueue(t *testing.T) {
	newPlugins()
	handler.NewVrfHandler(make([]byte, 0))
	defer func() {
		slash.db.Clear()
	}()

	gov.InitGenesisGovernParam(common.ZeroHash, slash.db, 2048)

	privateKey, _ := crypto.GenerateKey()
	vqList := make(staking.ValidatorQueue, 0)
	dataList := make([][]byte, 0)
	data := common.Int64ToBytes(time.Now().UnixNano())
	if err := slash.db.NewBlock(new(big.Int).SetUint64(1), blockHash, common.ZeroHash); nil != err {
		t.Fatal(err)
	}
	for i := 0; i < int(xcom.MaxConsensusVals()); i++ {
		vrfData, err := vrf.Prove(privateKey, data)
		if nil != err {
			t.Fatal(err)
		}
		data = vrf.ProofToHash(vrfData)
		dataList = append(dataList, data)

		tempPrivateKey, _ := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		nodeId := discover.PubkeyID(&tempPrivateKey.PublicKey)
		addr := crypto.PubkeyToNodeAddress(tempPrivateKey.PublicKey)
		v := &staking.Validator{
			NodeAddress: addr,
			NodeId:      nodeId,
		}
		vqList = append(vqList, v)
	}
	if enValue, err := rlp.EncodeToBytes(dataList); nil != err {
		t.Fatal(err)
	} else {
		if err := slash.db.Put(common.ZeroHash, handler.NonceStorageKey, enValue); nil != err {
			t.Fatal(err)
		}
	}
	resultQueue, err := randomOrderValidatorQueue(1, common.ZeroHash, vqList)
	if nil != err {
		t.Fatal(err)
	}
	assert.True(t, len(resultQueue) == len(vqList))
}

/**
Expand test cases
*/

func Test_IteratorCandidate(t *testing.T) {

	state, genesis, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
		return
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	for i := 0; i < 1000; i++ {

		var index int
		if i >= len(balanceStr) {
			index = i % (len(balanceStr) - 1)
		}

		//t.Log("Create Staking num:", index)

		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)

		ii := mrand.Intn(len(chaList))

		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))

		privateKey, err := crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random NodeId private key: %v", err)
			return
		}

		nodeId := discover.PubkeyID(&privateKey.PublicKey)

		privateKey, err = crypto.GenerateKey()
		if nil != err {
			t.Errorf("Failed to generate random Address private key: %v", err)
			return
		}

		addr := crypto.PubkeyToAddress(privateKey.PublicKey)

		var blsKey bls.SecretKey
		blsKey.SetByCSPRNG()
		var blsKeyHex bls.PublicKeyHex
		b, _ := blsKey.GetPublicKey().MarshalText()
		if err := blsKeyHex.UnmarshalText(b); nil != err {
			log.Error("Failed to blsKeyHex.UnmarshalText", "err", err)
			return
		}

		canTmp := &staking.Candidate{
			CandidateBase: &staking.CandidateBase{
				NodeId:          nodeId,
				BlsPubKey:       blsKeyHex,
				StakingAddress:  sender,
				BenefitAddress:  addr,
				StakingBlockNum: uint64(i),
				StakingTxIndex:  uint32(index),
				ProgramVersion:  xutil.CalcVersion(initProgramVersion),

				Description: staking.Description{
					NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
					ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
					Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
					Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
				},
			},
			CandidateMutable: &staking.CandidateMutable{
				Shares: balance,

				// Prevent null pointer initialization
				Released:           common.Big0,
				ReleasedHes:        common.Big0,
				RestrictingPlan:    common.Big0,
				RestrictingPlanHes: common.Big0,
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

		err = StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	// commit
	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
		return
	}

	stakingDB := staking.NewStakingDB()

	iter := stakingDB.IteratorCandidatePowerByBlockHash(blockHash2, 0)
	if err := iter.Error(); nil != err {
		t.Error("Get iter err", err)
		return
	}
	defer iter.Release()

	queue := make(staking.CandidateQueue, 0)

	for iter.Valid(); iter.Next(); {
		addrSuffix := iter.Value()
		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash2, addrSuffix)
		if nil != err {
			t.Errorf("Failed to Iterator Candidate info, err: %v", err)
			return
		}

		val := fmt.Sprint(can.ProgramVersion) + "_" + can.Shares.String() + "_" + fmt.Sprint(can.StakingBlockNum) + "_" + fmt.Sprint(can.StakingTxIndex)
		t.Log("Val:", val)

		queue = append(queue, can)
	}

	arrJson, _ := json.Marshal(queue)
	t.Log("CandidateList:", string(arrJson))
	t.Log("Candidate queue length:", len(queue))
}

func TestStakingPlugin_CalcDelegateIncome(t *testing.T) {
	del := new(staking.Delegation)
	del.Released = new(big.Int).SetInt64(0)
	del.RestrictingPlan = new(big.Int).SetInt64(0)
	del.ReleasedHes = new(big.Int).Mul(new(big.Int).SetInt64(100), new(big.Int).SetInt64(params.LAT))
	del.RestrictingPlanHes = new(big.Int).SetInt64(0)
	del.DelegateEpoch = 1
	del.CumulativeIncome = new(big.Int)
	per := make([]*reward.DelegateRewardPer, 0)
	per = append(per, &reward.DelegateRewardPer{
		Epoch:    1,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	per = append(per, &reward.DelegateRewardPer{
		Epoch:    2,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(200),
	})
	expectedCumulativeIncome := per[1].CalDelegateReward(del.ReleasedHes)
	calcDelegateIncome(3, del, per)
	assert.True(t, del.CumulativeIncome.Cmp(expectedCumulativeIncome) == 0)

	del = new(staking.Delegation)
	del.Released = new(big.Int).Mul(new(big.Int).SetInt64(100), new(big.Int).SetInt64(params.LAT))
	del.RestrictingPlan = new(big.Int).SetInt64(0)
	del.ReleasedHes = new(big.Int).Mul(new(big.Int).SetInt64(100), new(big.Int).SetInt64(params.LAT))
	del.RestrictingPlanHes = new(big.Int).SetInt64(0)
	del.DelegateEpoch = 2
	del.CumulativeIncome = new(big.Int)
	per = make([]*reward.DelegateRewardPer, 0)
	per = append(per, &reward.DelegateRewardPer{
		Epoch:    2,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})
	per = append(per, &reward.DelegateRewardPer{
		Epoch:    3,
		Delegate: new(big.Int).SetUint64(10),
		Reward:   new(big.Int).SetUint64(100),
	})

	expectedCumulativeIncome = per[0].CalDelegateReward(del.Released)
	expectedCumulativeIncome = expectedCumulativeIncome.Add(expectedCumulativeIncome, per[1].CalDelegateReward(new(big.Int).Add(del.Released, del.ReleasedHes)))
	calcDelegateIncome(4, del, per)
	assert.True(t, del.CumulativeIncome.Cmp(expectedCumulativeIncome) == 0)
}
