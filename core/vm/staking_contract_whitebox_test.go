package vm

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/x/handler"

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
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)

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
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)

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
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)
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
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)

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
	programVersion2, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign2 := common.VersionSign{}
	versionSign2.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersionBytes))
	sign2, _ := rlp.EncodeToBytes(versionSign2)

	var blsKey2 bls.SecretKey
	blsKey2.SetByCSPRNG()
	blsPkm2, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey2.GetPublicKey().Serialize()))

	args = append(args, fnType2)
	args = append(args, typ2)
	args = append(args, benefitAddress2)
	args = append(args, nodeId2)
	args = append(args, externalId2)
	args = append(args, nodeName2)
	args = append(args, website2)
	args = append(args, details2)
	args = append(args, amount2)
	args = append(args, programVersion2)
	args = append(args, sign2)
	args = append(args, blsPkm2)

	buf2 := new(bytes.Buffer)
	err = rlp.Encode(buf2, args)
	if err != nil {
		t.Errorf("createStaking2 encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking2 data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err = contract.Run(buf2.Bytes())

	assert.True(t, nil == err)

	var r2 xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, false, r2.Status)
	t.Log("the staking result Msg:", r2.ErrMsg)

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
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersionBytes))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, false, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)

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
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, false, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)

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
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, false, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)

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
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, false, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)

}

// todo
func Test_CreateStake_by_InvalidNodeId(t *testing.T) {
	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
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
	//0xf0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000010
	//0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff
	nid := discover.MustHexID("0xffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")

	nodeId, _ := rlp.EncodeToBytes(nid)
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"

	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, false, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)

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
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

	handler.GetCryptoHandler().SetPrivateKey(priKeyArr[index])

	versionSign := common.VersionSign{}
	versionSign.SetBytes(handler.GetCryptoHandler().MustSign(initProgramVersion))
	sign, _ := rlp.EncodeToBytes(versionSign)

	var blsKey bls.SecretKey
	blsKey.SetByCSPRNG()
	blsPkm, _ := rlp.EncodeToBytes(hex.EncodeToString(blsKey.GetPublicKey().Serialize()))

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, programVersion)
	params = append(params, sign)
	params = append(params, blsPkm)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, false, r.Status)
	t.Log("the staking result Msg:", r.ErrMsg)

}

func Test_CreateStake_by_WrongVersionSign(t *testing.T) {

}
