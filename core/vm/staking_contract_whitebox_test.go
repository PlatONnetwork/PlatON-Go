package vm_test

import (
	"bytes"
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/params"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
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

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
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

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}

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

	balance, _ := new(big.Int).SetString("20000000000000000000000000", 10)
	buildDbRestrictingPlan(t, sender, balance, 1, state)

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
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

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}

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

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:      newEvm(blockNumber, blockHash, state),
	}

	index := 1

	state.Prepare(txHashArr[index], blockHash, index+1)

	var args [][]byte
	args = make([][]byte, 0)

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

	version := uint32(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch)
	programVersion, _ := rlp.EncodeToBytes(version)

	args = append(args, fnType)
	args = append(args, typ)
	args = append(args, benefitAddress)
	args = append(args, nodeId)
	args = append(args, externalId)
	args = append(args, nodeName)
	args = append(args, website)
	args = append(args, details)
	args = append(args, amount)
	args = append(args, programVersion)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, args)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}
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

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
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

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}

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

	buf2 := new(bytes.Buffer)
	err = rlp.Encode(buf2, args)
	if err != nil {
		t.Errorf("createStaking2 encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking2 data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err = contract.Run(buf2.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}

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

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
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

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}

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

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
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

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}

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

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
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

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}

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

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
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

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	} else {
		t.Log(string(res))
	}

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
	}

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
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
	nid := discover.MustHexID("0x01234567890121345678901123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345")

	nodeId, _ := rlp.EncodeToBytes(nid)
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www." + nodeNameArr[index] + ".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index] + " super node")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	programVersion, _ := rlp.EncodeToBytes(initProgramVersion)

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

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Errorf("createStaking encode rlp data fail: %v", err)
	} else {
		t.Log("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		assert.Equal(t, err != err, "err is nil", err)

	} else {
		t.Log(string(res))
	}

}
