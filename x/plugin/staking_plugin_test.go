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
	"encoding/json"
	"fmt"
	"math/big"
	mrand "math/rand"
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
		preNonces = append(preNonces, crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(i))))[:])
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
		powerKey := staking.TallyPowerKey(canTmp.ProgramVersion, canTmp.Shares, canTmp.NodeId, canTmp.StakingBlockNum, canTmp.StakingTxIndex)
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
	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
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
		Released:           common.Big0,
		ReleasedHes:        common.Big0,
		RestrictingPlan:    common.Big0,
		RestrictingPlanHes: common.Big0,
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
		powerKey := staking.TallyPowerKey(canBase.ProgramVersion, canMutable.Shares, canBase.NodeId, canBase.StakingBlockNum, canBase.StakingTxIndex)
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

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
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

	nonce = crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
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
		powerKey := staking.TallyPowerKey(canBase.ProgramVersion, canMutable.Shares, canBase.NodeId, canBase.StakingBlockNum, canBase.StakingTxIndex)
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

	nonce := crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(1))))[:]
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

	/**
	Start HandleUnCandidateItem
	*/
	err = StakingInstance().HandleUnCandidateItem(state, blockNumber2.Uint64(), blockHash2, epoch+xcom.UnStakeFreezeDuration())

	if !assert.Nil(t, err, fmt.Sprintf("Failed to HandleUnCandidateItem: %v", err)) {
		return
	}

	// get Candidate info
	if _, err := getCandidate(blockHash2, index); snapshotdb.IsDbNotFoundErr(err) {
		t.Logf("expect candidate info is no found, err: %v", err)
		return
	} else {
		t.Error("It is not expect~")
		return
	}

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
	assert.True(t, del.CumulativeIncome == nil)

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

func TestStakingPlugin_WithdrewDelegate(t *testing.T) {

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
	_, err = StakingInstance().WithdrewDelegate(state, blockHash2, blockNumber2, amount, addrArr[index+1],
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
	issueIncome, err := StakingInstance().WithdrewDelegate(state, blockHash3, curBlockNumber, del.ReleasedHes, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del, delegateRewardPerList)

	if !assert.Nil(t, err, fmt.Sprintf("Failed to WithdrewDelegate: %v", err)) {
		return
	}

	can, err = getCandidate(blockHash3, index)

	assert.True(t, expectedIssueIncome.Cmp(issueIncome) == 0)
	assert.True(t, expectedBalance.Cmp(state.GetBalance(addrArr[index+1])) == 0)
	t.Log("Get Candidate Info is:", can)
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
	slashQueue := make(staking.CandidateQueue, 2)

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

	slashItem1 := &staking.SlashNodeItem{
		NodeId:      slash1.NodeId,
		Amount:      slash1.Released,
		SlashType:   staking.LowRatio,
		BenefitAddr: vm.RewardManagerPoolAddr,
	}

	sla := new(big.Int).Div(slash2.Released, big.NewInt(10))
	caller := common.MustBech32ToAddress("lax1uj3zd9yz00axz7ls88ynwsp3jprhjd9ldx9qpm")

	slashItem2 := &staking.SlashNodeItem{
		NodeId:      slash2.NodeId,
		Amount:      sla,
		SlashType:   staking.DuplicateSign,
		BenefitAddr: caller,
	}
	slashItemQueue := make(staking.SlashQueue, 0)
	slashItemQueue = append(slashItemQueue, slashItem1)
	slashItemQueue = append(slashItemQueue, slashItem2)

	err = StakingInstance().SlashCandidates(state, blockHash2, blockNumber2.Uint64(), slashItemQueue...)

	assert.Nil(t, err, fmt.Sprintf("Failed to SlashCandidates Second can (DuplicateSign), err: %v", err))

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
