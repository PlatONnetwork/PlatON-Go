package vm

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	_ "fmt"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

// Custom func
func create_staking(blockNumber *big.Int, blockHash common.Hash, state *mock.MockStateDB, index int, t *testing.T) *StakingContract {

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

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

	xcom.GetCryptoHandler().SetPrivateKey(priKeyArr[index])
	//xcom.GetCryptoHandler().SetPrivateKey(crypto.HexMustToECDSA("30999a0591d350be57ef38ff17babdc73aa8e091532c0a8aafb76375dcb9f0c9"))

	versionSign := common.VersionSign{}
	versionSign.SetBytes(xcom.GetCryptoHandler().MustSign(initProgramVersion))
	//versionSign.SetBytes(xcom.GetCryptoHandler().MustSign(1793))
	sign, _ := rlp.EncodeToBytes(versionSign)

	//fmt.Println("The versionSign", versionSign.Hex())
	//return nil

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

	return contract
}

func create_delegate(contract *StakingContract, index int, t *testing.T) {
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1004))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10)
	amount, _ := rlp.EncodeToBytes(StakeThreshold)

	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, nodeId)
	params = append(params, amount)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("delegate encode rlp data fail", err)
		return
	} else {
		t.Log("delegate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the delegate result Msg:", r.ErrMsg)

}

func getCandidate(contract *StakingContract, index int, t *testing.T) {
	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1105))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])

	params = append(params, fnType)
	params = append(params, nodeId)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getCandidate encode rlp data fail: %v", err)
		return
	} else {
		t.Log("getCandidate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the Candidate info:", r.Data)

}

/**
Standard test cases
*/

func TestStakingContract_createStaking(t *testing.T) {

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

	create_staking(blockNumber, blockHash, state, 1, t)
}

func TestStakingContract_editCandidate(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}
	state.Prepare(txHashArr[0], blockHash, 0)
	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 error: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	state.Prepare(txHashArr[1], blockHash2, 1)

	// edit
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1001))

	benefitAddress, _ := rlp.EncodeToBytes(addrArr[0])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("I am Xu !?")
	nodeName, _ := rlp.EncodeToBytes("Xu, China")
	website, _ := rlp.EncodeToBytes("https://www.Xu.net")
	details, _ := rlp.EncodeToBytes("Xu super node")

	params = append(params, fnType)
	params = append(params, benefitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("edit candidate encode rlp data fail: %v", err)
		return
	} else {
		t.Log("edit candidate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the editStaking result Msg:", r.ErrMsg)

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 error: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

}

func TestStakingContract_increaseStaking(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}
	state.Prepare(txHashArr[0], blockHash, 0)
	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 error: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	state.Prepare(txHashArr[1], blockHash2, 1)

	// increase

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1002))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	typ, _ := rlp.EncodeToBytes(uint16(0))
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index-1], 10) // equal or more than "1000000000000000000000000"
	amount, _ := rlp.EncodeToBytes(StakeThreshold)

	params = append(params, fnType)
	params = append(params, nodeId)
	params = append(params, typ)
	params = append(params, amount)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("increaseStaking encode rlp data fail", err)
		return
	} else {
		t.Log("increaseStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the increaseStaking result Msg:", r.ErrMsg)

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 error: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

}

func TestStakingContract_withdrewCandidate(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	state.Prepare(txHashArr[0], blockHash, 0)
	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	state.Prepare(txHashArr[1], blockHash2, 1)

	// withdrewStaking

	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1003))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])

	params = append(params, fnType)
	params = append(params, nodeId)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("withdrewStaking encode rlp data fail", err)
		return
	} else {
		t.Log("withdrewStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the withdrew candidate result Msg:", r.ErrMsg)

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get CandidateInfo
	//getCandidate(contract2, index, t)

}

func TestStakingContract_delegate(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	state.Prepare(txHashArr[0], blockHash, 0)
	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, delegate_sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	state.Prepare(txHashArr[1], blockHash2, 1)
	// delegate
	create_delegate(contract2, index, t)

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

}

func TestStakingContract_withdrewDelegate(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	state.Prepare(txHashArr[0], blockHash, 0)
	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, delegate_sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	state.Prepare(txHashArr[1], blockHash, 1)
	// delegate
	create_delegate(contract, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, delegate_sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	state.Prepare(txHashArr[2], blockHash2, 0)

	// withdrewDelegate
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1005))
	stakingBlockNum, _ := rlp.EncodeToBytes(blockNumber.Uint64())
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	withdrewAmount, _ := new(big.Int).SetString(balanceStr[index], 10)
	amount, _ := rlp.EncodeToBytes(withdrewAmount)

	params = append(params, fnType)
	params = append(params, stakingBlockNum)
	params = append(params, nodeId)
	params = append(params, amount)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("delegate encode rlp data fail", err)
		return
	} else {
		t.Log("delegate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the withdelegate result Msg:", r.ErrMsg)

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)
}

func TestStakingContract_getVerifierList(t *testing.T) {

	state, genesis, _ := newChainState()
	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}
	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// init staking data into block 1
	build_staking_data(genesis.Hash())

	if err := sndb.NewBlock(blockNumber2, genesis.Hash(), blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
		return
	}

	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1100))

	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getVerifierList encode rlp data fail:%v", err)
		return
	} else {
		t.Log("getVerifierList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the getVerifierList result Data:", r.Data)

}

func TestStakingContract_getValidatorList(t *testing.T) {

	state, genesis, _ := newChainState()
	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}
	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	// init staking data into block 1
	build_staking_data(genesis.Hash())

	if err := sndb.NewBlock(blockNumber2, genesis.Hash(), blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber2, err)
		return
	}

	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1101))

	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getValidatorList encode rlp data fail:%v", err)
		return
	} else {
		t.Log("getValidatorList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the getValidatorList result Data:", r.Data)

}

func TestStakingContract_getCandidateList(t *testing.T) {

	state, genesis, _ := newChainState()

	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
		return
	}

	for i := 0; i < 2; i++ {
		state.Prepare(txHashArr[i], blockHash, i)
		create_staking(blockNumber, blockHash, state, i, t)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}
	//sndb.Compaction()

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	for i := 2; i < 4; i++ {
		state.Prepare(txHashArr[i], blockHash2, i)
		create_staking(blockNumber2, blockHash2, state, i, t)
	}

	// getCandidate List
	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}
	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1102))

	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("getCandidateList encode rlp data fail:%v", err)
		return
	} else {
		t.Log("getCandidateList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the getCandidateList result Data:", r.Data)

}

func TestStakingContract_getRelatedListByDelAddr(t *testing.T) {

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

	for i := 0; i < 4; i++ {
		state.Prepare(txHashArr[i], blockHash, i)
		create_staking(blockNumber, blockHash, state, i, t)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, delegate_sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// delegate
	for i := 0; i < 3; i++ {
		state.Prepare(txHashArr[i], blockHash2, i)
		create_delegate(contract2, i, t)
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
		return
	}

	// get RelatedListByDelAddr
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1103))
	delAddr, _ := rlp.EncodeToBytes(sender)

	params = append(params, fnType)
	params = append(params, delAddr)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("getRelatedListByDelAddr encode rlp data fail", err)
		return
	} else {
		t.Log("getRelatedListByDelAddr data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the getRelatedListByDelAddr result Data:", r.Data)
}

func TestStakingContract_getDelegateInfo(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	index := 1

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("newBlock err", err)
		return
	}

	state.Prepare(txHashArr[0], blockHash, 0)
	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	contract := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, delegate_sender),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	state.Prepare(txHashArr[1], blockHash, 1)
	// delegate
	create_delegate(contract, index, t)

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Failed to commit block 1, err: %v", err)
		return
	}
	//sndb.Compaction()

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
		return
	}

	contract2 := &StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0, sender),
		Evm:      newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	state.Prepare(txHashArr[2], blockHash2, 2)
	// get DelegateInfo
	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1104))
	stakingBlockNum, _ := rlp.EncodeToBytes(blockNumber.Uint64())
	delAddr, _ := rlp.EncodeToBytes(delegate_sender)
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])

	params = append(params, fnType)
	params = append(params, stakingBlockNum)
	params = append(params, delAddr)
	params = append(params, nodeId)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("getDelegateInfo encode rlp data fail", err)
		return
	} else {
		t.Log("getDelegateInfo data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())

	assert.True(t, nil == err)

	var r xcom.Result
	err = json.Unmarshal(res, &r)
	assert.True(t, nil == err)
	assert.Equal(t, true, r.Status)
	t.Log("the getDelegateInfo result Data:", r.Data)
}

func TestStakingContract_getCandidateInfo(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("Failed to newBlock", err)
		return
	}

	state.Prepare(txHashArr[0], blockHash, 0)
	contract := create_staking(blockNumber, blockHash, state, 1, t)
	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
		return
	}
	//sndb.Compaction()

	// get candidate Info
	getCandidate(contract, 1, t)
}

/**
Expand test cases
*/

func TestStakingContract_batchCreateStaking(t *testing.T) {

	state, genesis, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()
	defer func() {
		sndb.Clear()
	}()

	if err := sndb.NewBlock(blockNumber, genesis.Hash(), blockHash); nil != err {
		t.Error("Failed to newBlock", err)
		return
	}

	for i := 0; i < 4; i++ {
		state.Prepare(txHashArr[i], blockHash, i)
		create_staking(blockNumber, blockHash, state, i, t)
	}

}

func TestStakingContract_cleanSnapshotDB(t *testing.T) {
	sndb := snapshotdb.Instance()
	sndb.Clear()
}
