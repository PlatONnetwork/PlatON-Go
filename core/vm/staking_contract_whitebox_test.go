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
	"encoding/json"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/node"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
)

/**
This is a white test cases for staking_contract
*/

/**
susccess test case
*/
func Test_CreateStake_HighThreshold_by_freeVon(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	index := 1

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	runContractSendTransaction(contract, params, "createStaking", t)

}

func Test_CreateStake_HighThreshold_by_restrictplanVon(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	balance, _ := new(big.Int).SetString(balanceStr[index], 10)
	buildDbRestrictingPlan(t, sender, balance, 1, state)

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(1))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	runContractSendTransaction(contract, params, "createStaking", t)

}

func Test_CreateStake_RightVersion(t *testing.T) {
	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	index := 1

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	runContractSendTransaction(contract, params, "createStaking", t)
}

/**
failure test case
*/
func Test_CreateStake_RepeatStake(t *testing.T) {
	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	index := 1

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	runContractSendTransaction(contract, params, "createStaking", t)

	// repeat stake
	var args [][]byte
	args = make([][]byte, 0)

	fnType2, _ := rlp.EncodeToBytes(uint16(1000))
	typ2, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress2, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId2, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId2, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName2, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website2, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details2, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold2, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"
	amount2, _ := rlp.EncodeToBytes(StakeThreshold2)
	rewardPer2, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion2, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign2 := common.VersionSign{}
	versionSign2.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign2, _ := rlp.EncodeToBytes(versionSign2)

	var blsKey2 bls.SecretKey
	blsKey2.SetByCSPRNG()

	var keyEntries2 bls.PublicKeyHex
	blsHex2 := hex.EncodeToString(blsKey2.GetPublicKey().Serialize())
	keyEntries2.UnmarshalText([]byte(blsHex2))

	blsPkm2, _ := rlp.EncodeToBytes(keyEntries2)

	// generate the bls proof
	proof2, _ := blsKey2.MakeSchnorrNIZKP()
	proofByte2, _ := proof2.MarshalText()
	var proofHex2 bls.SchnorrProofHex
	proofHex2.UnmarshalText(proofByte2)
	proofRlp2, _ := rlp.EncodeToBytes(proofHex2)

	args = append(args, fnType2)
	args = append(args, typ2)
	args = append(args, benefitAddress2)
	args = append(args, nodeId2)
	args = append(args, externalId2)
	args = append(args, nodeName2)
	args = append(args, website2)
	args = append(args, details2)
	args = append(args, amount2)
	args = append(args, rewardPer2)
	args = append(args, programVersion2)
	args = append(args, sign2)
	args = append(args, blsPkm2)
	args = append(args, proofRlp2)

	buf2 := new(bytes.Buffer)
	err := rlp.Encode(buf2, args)
	if err != nil {
		t.Errorf("createStaking2 encode rlp data fail: %v", err)
		return
	} else {
		t.Log("createStaking2 data rlp: ", hexutil.Encode(buf2.Bytes()))
	}

	res, err := contract.Run(buf2.Bytes())

	assert.True(t, nil == err)

	var r2 uint32
	err = json.Unmarshal(res, &r2)
	assert.True(t, nil == err)
	assert.NotEqual(t, common.OkCode, r2)

}

func Test_CreateStake_LowBalance_by_freeVon(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	index := 1

	state.Prepare(txHashArr[index], blockHash, index+1)

	// reset sender balance
	state.SubBalance(sender, state.GetBalance(sender))

	StakeThreshold := xcom.StakeThreshold()
	initBalance := new(big.Int).Sub(xcom.StakeThreshold(), common.Big1) // equal or more than "1000000000000000000000000"
	state.AddBalance(sender, initBalance)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")

	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersionBytes))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
		return
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r uint32
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.NotEqual(t, common.OkCode, r)

}

func Test_CreateStake_LowThreshold_by_freeVon(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	index := 1

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold := new(big.Int).Sub(xcom.StakeThreshold(), common.Big1) // equal or more than "1000000000000000000000000"

	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
		return
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r uint32
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.NotEqual(t, common.OkCode, r)

}

func Test_CreateStake_LowBalance_by_restrictplanVon(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	StakeThreshold := xcom.StakeThreshold()
	initBalance := new(big.Int).Sub(xcom.StakeThreshold(), common.Big1) // equal or more than "1000000000000000000000000"

	buildDbRestrictingPlan(t, sender, initBalance, 1, state)

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(1))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")

	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
		return
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r uint32
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.NotEqual(t, common.OkCode, r)

}

func Test_CreateStake_LowThreshold_by_restrictplanVon(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	index := 1

	StakeThreshold := xcom.StakeThreshold()
	initBalance := new(big.Int).Sub(xcom.StakeThreshold(), common.Big1) // equal or more than "1000000000000000000000000"

	buildDbRestrictingPlan(t, sender, StakeThreshold, 1, state)

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(1))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")

	amount, _ := rlp.EncodeToBytes(initBalance)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
		return
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r uint32
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.NotEqual(t, common.OkCode, r)

}

func Test_CreateStake_by_InvalidNodeId(t *testing.T) {
	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	index := 1

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])

	// build a invalid nodeId
	//
	//0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
	//0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000
	nid := discover.MustHexID("0x00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000")

	nodeId, _ := rlp.EncodeToBytes(nid)
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"

	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
		return
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r uint32
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.NotEqual(t, common.OkCode, r)

}

func Test_CreateStake_by_FlowDescLen(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	index := 1

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("sssxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxsfsfsfsfsfsfsfsfsfsfsfsfsfADADADADADs")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China adadadadafsafsdfsdfsdfsdfsdfsdADADADADADADADAf")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".networkdadadadasdwdqwdqwdADADADADADADADADADAqwdqwdqwdqwdqwdQWDQwdQWD.com")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super nodeFFFAADADDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDDD")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"

	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
		return
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r uint32
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.NotEqual(t, common.OkCode, r)

}

func Test_CreateStake_by_LowVersionSign(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	index := 1

	state.Prepare(txHashArr[index], blockHash, index+1)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("test low version")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", Low version")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".lowVersion.com")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node low version")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"

	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))

	version := uint32(0<<16 | 9<<8 | 0)

	programVersion, _ := rlp.EncodeToBytes(version)

	node.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(node.GetCryptoHandler().MustSign(version))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()

	var keyEntries bls.PublicKeyHex
	blsHex := hex.EncodeToString(blsKey.GetPublicKey().Serialize())
	keyEntries.UnmarshalText([]byte(blsHex))

	blsPkm, _ := rlp.EncodeToBytes(keyEntries)

	// generate the bls proof
	proof, _ := blsKey.MakeSchnorrNIZKP()
	proofByte, _ := proof.MarshalText()
	var proofHex bls.SchnorrProofHex
	proofHex.UnmarshalText(proofByte)
	proofRlp, _ := rlp.EncodeToBytes(proofHex)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, rewardPer)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)
	params = append(params, proofRlp)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
		return
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r uint32
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.NotEqual(t, common.OkCode, r)

}

func Test_EditStake_by_RightParams(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
	}
	state.Prepare(txHashArr[0], blockHash, 0)

	contract := create_staking(blockNumber, blockHash, state, 1, t)

	index := 1

	state.Prepare(txHashArr[index+1], blockHash, index+2)

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1001))

	benefitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	rewardPer, _ := rlp.EncodeToBytes(uint64(5000))
	externalId, _ := rlp.EncodeToBytes("test low version")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", Low version")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".lowVersion.com")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node low version")

	params = append(params, fnType)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, rewardPer)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("editStaking encode rlp data fail: %v", err)
		return
	} else {
		t.Log("editStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r uint32
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, common.OkCode, r)

}
