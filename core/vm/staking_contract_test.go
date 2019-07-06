package vm_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"math/big"
	"testing"
)





// Custom func
func create_staking (blockNumber *big.Int, blockHash common.Hash, state *state.StateDB, index int, t *testing.T) *vm.StakingContract {

	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber, blockHash, state),
	}

	state.Prepare(txHashArr[index], blockHash, index + 1)


	var params [][]byte
	params = make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1000))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	benifitAddress, _ := rlp.EncodeToBytes(addrArr[index])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("xssssddddffffggggg")
	nodeName, _ := rlp.EncodeToBytes(nodeNameArr[index] + ", China")
	website, _ := rlp.EncodeToBytes("https://www."+nodeNameArr[index]+".network")
	details, _ := rlp.EncodeToBytes(nodeNameArr[index]+" super node")
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index], 10) // equal or more than "1000000000000000000000000"
	amount, _ := rlp.EncodeToBytes(StakeThreshold)
	processVersion, _ := rlp.EncodeToBytes(initProcessVersion)


	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, benifitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)
	params = append(params, amount)
	params = append(params, processVersion)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("createStaking encode rlp data fail")
	} else {
		fmt.Println("createStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}


	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}

	return contract
}

func create_delegate (contract *vm.StakingContract, index int, t *testing.T) {
	var params [][]byte
	params = make([][]byte, 0)


	fnType, _ := rlp.EncodeToBytes(uint16(1004))
	typ, _ := rlp.EncodeToBytes(uint16(0))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	StakeThreshold, _ := new(big.Int).SetString(balanceStr[index+6], 10)
	amount, _ := rlp.EncodeToBytes(StakeThreshold)


	params = append(params, fnType)
	params = append(params, typ)
	params = append(params, nodeId)
	params = append(params, amount)


	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("delegate encode rlp data fail", err)
	} else {
		t.Log("delegate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}
}

func getCandidate (contract *vm.StakingContract, index int, t *testing.T) {
	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1105))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])

	params = append(params, fnType)
	params = append(params, nodeId)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("getCandidate encode rlp data fail")
	} else {
		fmt.Println("getCandidate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("getCandidate err", err)
	}else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			fmt.Println(err)
		}

		if r.Status {
			t.Log("the Candidate info:", r.Data)
		}else {
			t.Error("getCandidate failed", r.ErrMsg)
		}
	}
}

func TestRLP_encode (t *testing.T) {

	var params [][]byte
	params = make([][]byte, 0)

	fnType, err := rlp.EncodeToBytes(uint16(1100))
	if nil != err {
		fmt.Println("fnType err", err)
	}else {
		var num uint16
		rlp.DecodeBytes(fnType, &num)
		fmt.Println("num is ", num)
	}
	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err = rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("rlp stakingContract encode rlp data fail")
	} else {
		fmt.Println("rlp stakingContract data rlp: ", hexutil.Encode(buf.Bytes()))
	}
}


/**
Standard test cases
 */

func TestStakingContract_createStaking(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()


	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		fmt.Println("newBlock err", err)
	}
	create_staking(blockNumber, blockHash, state, 1, t)
}


func TestStakingContract_editorCandidate(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()

	index := 1

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		fmt.Println("newBlock err", err)
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)


	sndb.Commit(blockHash)


	// get CandidateInfo
	getCandidate(contract1, index, t)


	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber2, blockHash2, state),
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

	// edit
	var params [][]byte
	params = make([][]byte, 0)


	fnType, _ := rlp.EncodeToBytes(uint16(1001))

	benifitAddress, _ := rlp.EncodeToBytes(addrArr[0])
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	externalId, _ := rlp.EncodeToBytes("I am Xu !?")
	nodeName, _ := rlp.EncodeToBytes("Xu, China")
	website, _ := rlp.EncodeToBytes("https://www.Xu.net")
	details, _ := rlp.EncodeToBytes("Xu super node")



	params = append(params, fnType)
	params = append(params, benifitAddress)
	params = append(params, nodeId)
	params = append(params, externalId)
	params = append(params, nodeName)
	params = append(params, website)
	params = append(params, details)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("edit candidate encode rlp data fail")
	} else {
		fmt.Println("edit candidate data rlp: ", hexutil.Encode(buf.Bytes()))
	}


	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}

	// get CandidateInfo
	getCandidate(contract2, index, t)

}


func TestStakingContract_increaseStaking (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()

	index := 1

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	sndb.Commit(blockHash)


	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber2, blockHash2, state),
	}


	// get CandidateInfo
	getCandidate(contract2, index, t)

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
	} else {
		t.Log("increaseStaking data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}

	sndb.Commit(blockHash2)


	// get CandidateInfo
	getCandidate(contract2, index, t)

}

func TestStakingContract_withdrewCandidate (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()

	index := 1

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	sndb.Commit(blockHash)


	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber2, blockHash2, state),
	}


	// get CandidateInfo
	getCandidate(contract2, index, t)

	// withdrewCandidate

	var params [][]byte
	params = make([][]byte, 0)


	fnType, _ := rlp.EncodeToBytes(uint16(1003))
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])


	params = append(params, fnType)
	params = append(params, nodeId)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("withdrewCandidate encode rlp data fail", err)
	} else {
		t.Log("withdrewCandidate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}

	sndb.Commit(blockHash2)


	// get CandidateInfo
	getCandidate(contract2, index, t)

}


func TestStakingContract_delegate (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()

	index := 1

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	sndb.Commit(blockHash)

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber2, blockHash2, state),
	}


	// get CandidateInfo
	getCandidate(contract2, index, t)

	// delegate
	create_delegate(contract2, index, t)

	sndb.Commit(blockHash2)


	// get CandidateInfo
	getCandidate(contract2, index, t)



}

func TestStakingContract_withdrewDelegate (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()

	index := 1

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	// delegate
	create_delegate(contract1, index, t)

	sndb.Commit(blockHash)

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber2, blockHash2, state),
	}


	// get CandidateInfo
	getCandidate(contract2, index, t)

	// withdrewDelegate
	var params [][]byte
	params = make([][]byte, 0)


	fnType, _ := rlp.EncodeToBytes(uint16(1005))
	stakingBlockNum, _ := rlp.EncodeToBytes(blockNumber.Uint64())
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])
	StakeThreshold, _ := new(big.Int).SetString("4600000", 10)
	amount, _ := rlp.EncodeToBytes(StakeThreshold)


	params = append(params, fnType)
	params = append(params, stakingBlockNum)
	params = append(params, nodeId)
	params = append(params, amount)


	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("delegate encode rlp data fail", err)
	} else {
		t.Log("delegate data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}

	sndb.Commit(blockHash2)

	// get CandidateInfo
	getCandidate(contract2, index, t)
}


func TestStakingContract_getVerifierList (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()
	state, _ := newChainState()
	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber, blockHash, state),
	}
	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
	}


	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1100))

	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("getVerifierList encode rlp data fail")
	} else {
		fmt.Println("getVerifierList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("getVerifierList err", err)
	}else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			fmt.Println(err)
		}

		if r.Status {
			t.Log("the VerifierList info:", r.Data)
		}else {
			t.Error("getVerifierList failed", r.ErrMsg)
		}
	}


}


func TestStakingContract_getValidatorList (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()
	state, _ := newChainState()
	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber, blockHash, state),
	}
	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
	}


	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1101))

	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("getValidatorList encode rlp data fail")
	} else {
		fmt.Println("getValidatorList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("getValidatorList err", err)
	}else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			fmt.Println(err)
		}

		if r.Status {
			t.Log("the ValidatorList info:", r.Data)
		}else {
			t.Error("getValidatorList failed", r.ErrMsg)
		}
	}

}


func TestStakingContract_getCandidateList(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()

	//state.Prepare(txHashArr[idx], blockHash, idx)
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Errorf("newBlock failed, blockNumber1: %d, err:%v", blockNumber, err)
	}


	for i:= 0; i < 2; i++ {
		create_staking(blockNumber, blockHash, state, i, t)
	}

	sndb.Commit(blockHash)
	//sndb.Compaction()

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
	}


	for i:= 2; i < 4; i++ {
		create_staking(blockNumber2, blockHash2, state, i, t)
	}


	// getCandidate List
	contract := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber2, blockHash2, state),
	}
	params := make([][]byte, 0)

	fnType, _ := rlp.EncodeToBytes(uint16(1102))

	params = append(params, fnType)

	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		fmt.Println(err)
		t.Errorf("getCandidateList encode rlp data fail")
	} else {
		fmt.Println("getCandidateList data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract.Run(buf.Bytes())
	if nil != err {
		t.Error("getCandidateList err", err)
	}else {

		var r xcom.Result
		err = json.Unmarshal(res, &r)
		if nil != err {
			fmt.Println(err)
		}

		if r.Status {
			t.Log("the CandidateList info:", r.Data)
		}else {
			t.Error("CandidateList failed", r.ErrMsg)
		}
	}


}
 // TODO
func TestStakingContract_getRelatedListByDelAddr (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()

	index := 1

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	// delegate
	create_delegate(contract1, index, t)

	sndb.Commit(blockHash)

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber2, blockHash2, state),
	}


	// get CandidateInfo
	getCandidate(contract2, index, t)

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
	} else {
		t.Log("getRelatedListByDelAddr data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}
}

func TestStakingContract_getDelegateInfo (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()

	index := 1

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	contract1 := create_staking(blockNumber, blockHash, state, index, t)

	// delegate
	create_delegate(contract1, index, t)

	sndb.Commit(blockHash)
	//sndb.Compaction()

	// get CandidateInfo
	getCandidate(contract1, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Errorf("newBlock failed, blockNumber2: %d, err:%v", blockNumber2, err)
	}

	contract2 := &vm.StakingContract{
		Plugin:   plugin.StakingInstance(),
		Contract: newContract(common.Big0),
		Evm:	 newEvm(blockNumber2, blockHash2, state),
	}


	// get CandidateInfo
	getCandidate(contract2, index, t)

	// get DelegateInfo
	var params [][]byte
	params = make([][]byte, 0)


	fnType, _ := rlp.EncodeToBytes(uint16(1104))
	stakingBlockNum, _ := rlp.EncodeToBytes(blockNumber.Uint64())
	delAddr, _ := rlp.EncodeToBytes(sender)
	nodeId, _ := rlp.EncodeToBytes(nodeIdArr[index])



	params = append(params, fnType)
	params = append(params, stakingBlockNum)
	params = append(params, delAddr)
	params = append(params, nodeId)


	buf := new(bytes.Buffer)
	err := rlp.Encode(buf, params)
	if err != nil {
		t.Error("getDelegateInfo encode rlp data fail", err)
	} else {
		t.Log("getDelegateInfo data rlp: ", hexutil.Encode(buf.Bytes()))
	}

	res, err := contract2.Run(buf.Bytes())
	if nil != err {
		t.Error(err)
	}else {
		t.Log(string(res))
	}
}


func TestStakingContract_getCandidateInfo (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		fmt.Println("newBlock err", err)
	}
	contract := create_staking(blockNumber, blockHash, state, 1, t)
	sndb.Commit(blockHash)
	//sndb.Compaction()


	// get candidate Info
	getCandidate(contract, 1, t)
}


/**
Expand test cases
 */


func TestStakingContract_batchCreateStaking(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, _ := newChainState()
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		fmt.Println("newBlock err", err)
	}

	for i:= 0; i < 4; i++ {
		create_staking(blockNumber, blockHash, state, 1, t)
	}

}




func TestStakingContract_cleanSnapshotDB(t *testing.T) {
	sndb := snapshotdb.Instance()
	sndb.Clear()
}