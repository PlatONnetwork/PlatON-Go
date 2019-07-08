package plugin_test

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"math/big"
	mrand "math/rand"
	"strconv"
	"testing"
	"time"
)


func create_staking (state *state.StateDB, blockNumber *big.Int, blockHash common.Hash, index int, typ uint16, t *testing.T) error {

	balance, _ := new(big.Int).SetString(balanceStr[index], 10)
	canTmp := &staking.Candidate{
		NodeId:          nodeIdArr[index],
		StakingAddress:  sender,
		BenifitAddress:  addrArr[index],
		StakingBlockNum: blockNumber.Uint64(),
		StakingTxIndex:  uint32(index),
		Shares:          balance,

		// Prevent null pointer initialization
		Released: common.Big0,
		ReleasedHes: common.Big0,
		RestrictingPlan: common.Big0,
		RestrictingPlanHes: common.Big0,

		Description: staking.Description{
			NodeName:   nodeNameArr[index],
			ExternalId: nodeNameArr[index] + chaList[len(chaList)%(index+1)] + "balabalala" + chaList[index],
			Website:    "www." + nodeNameArr[index] + ".org",
			Details:    "This is " + nodeNameArr[index] + " Super Node",
		},
	}

	canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

	return plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, initProcessVersion, typ, canAddr, canTmp)
}

func getCandidate (blockHash common.Hash, index int, t *testing.T) *staking.Candidate {
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	var c *staking.Candidate
	if can, err := plugin.StakingInstance().GetCandidateInfo(blockHash, addr); nil != err {
		t.Log("Failed to Get Candidate info", err)
	}else {

		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}
	return c
}


func delegate (state xcom.StateDB, blockHash common.Hash, blockNumber *big.Int,
	 can *staking.Candidate, typ uint16, index int, t *testing.T) (*staking.Delegation, error) {

	delAddr := addrArr[index+1]

	// build delegate
	del := new(staking.Delegation)

	// Prevent null pointer initialization
	del.Released = common.Big0
	del.RestrictingPlan = common.Big0
	del.ReleasedHes = common.Big0
	del.RestrictingPlanHes = common.Big0
	del.Reduction = common.Big0

	//amount := common.Big257  // FAIL
	amount, _ := new(big.Int).SetString(balanceStr[index+1], 10)  // PASS

	return del, plugin.StakingInstance().Delegate(state, blockHash, blockNumber, delAddr, del, can, 0, amount)
}


func getDelegate (blockHash common.Hash, stakingNum uint64, index int, t *testing.T) *staking.Delegation {


	del, err := plugin.StakingInstance().GetDelegateInfo(blockHash, addrArr[index+1], nodeIdArr[index], stakingNum)
	if nil != err {
		t.Log("Failed to GetDelegateInfo:", err)
	}else {
		delByte, _ := json.Marshal(del)
		t.Log("Get Candidate Info is:", string(delByte))
	}
	return del
}



/**
Standard test cases
*/

func TestStakingPlugin_BeginBlock(t *testing.T) {

}

func TestStakingPlugin_EndBlock(t *testing.T) {

}

func TestStakingPlugin_Confirmed(t *testing.T) {

}


func TestStakingPlugin_CreateCandidate(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()


	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}


	if err := create_staking(state, blockNumber, blockHash, 1, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}
}

func TestStakingPlugin_GetCandidateInfo(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	// Get Candidate Info
	getCandidate(blockHash, index, t)

}

func TestStakingPlugin_GetCandidateInfoByIrr(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	// Get Candidate Info
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])
	if can, err := plugin.StakingInstance().GetCandidateInfoByIrr(addr); nil != err {
		t.Error("Failed to Get Candidate info", err)
	}else {

		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
	}
}

func TestStakingPlugin_GetCandidateList(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}


	for i := 0; i < 4; i++ {
		if err := create_staking(state, blockNumber, blockHash, i, 0, t); nil != err {
			t.Error("Failed to Create num: " + fmt.Sprint(i) + " Staking", err)
		}
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	if queue, err := plugin.StakingInstance().GetCandidateList(blockHash); nil != err {
		t.Error("Failed to GetCandidateList", err)
	}else {
		queueByte, _ := json.Marshal(queue)
		t.Log("GetCandidateList is:", string(queueByte))
	}
}

func TestStakingPlugin_EditorCandidate(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
	}

	// Get Candidate Info
	c := getCandidate(blockHash, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
	}

	// Edit Candidate
	c.NodeName = nodeNameArr[index+1]
	c.ExternalId = "What is this ?"
	c.Website = "www.baidu.com"
	c.Details = "This is buidu website ?"
	if err := plugin.StakingInstance().EditorCandidate(blockHash2, blockNumber2, c); nil != err {
		t.Error("Failed to EditorCandidate", err)
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
	}

	// get Candidate info after edit
	getCandidate(blockHash2, index, t)

}

func TestStakingPlugin_IncreaseStaking(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
	}

	// Get Candidate Info
	c := getCandidate(blockHash, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
	}

	// IncreaseStaking
	if err := plugin.StakingInstance().IncreaseStaking(state, blockHash2, blockNumber2, common.Big256, uint16(0), c); nil != err {
		t.Error("Failed to IncreaseStaking", err)
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Errorf("Commit 2 err: %v", err)
	}

	// get Candidate info
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])
	if can, err := plugin.StakingInstance().GetCandidateInfoByIrr(addr); nil != err {
	//if can, err := plugin.StakingInstance().GetCandidateInfo(blockHash2, addr); nil != err {
		t.Error("Failed to Get Candidate info After Increase", err)
	}else {

		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info After Increase is:", string(canByte))
		c = can
	}

}

func TestStakingPlugin_WithdrewCandidate(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	// Get Candidate Info
	c := getCandidate(blockHash, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
	}

	// IncreaseStaking
	if err := plugin.StakingInstance().WithdrewCandidate(state, blockHash2, blockNumber2,  c); nil != err {
		t.Error("Failed to WithdrewCandidate", err)
	}

	t.Log("Finish WithdrewCandidate ~~")
	// get Candidate info
	getCandidate(blockHash2, index, t)

}

func TestStakingPlugin_HandleUnCandidateItem(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
	}

	// Get Candidate Info
	getCandidate(blockHash, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
	}

    // Add UNStakingItems
	stakingDB := staking.NewStakingDB ()

	epoch := xutil.CalculateEpoch(blockNumber2.Uint64())
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	if err := stakingDB.AddUnStakeItemStore(blockHash2, epoch, addr); nil != err {
		t.Error("Failed to AddUnStakeItemStore:", err)
		return
	}

	err = plugin.StakingInstance().HandleUnCandidateItem(state, blockHash2,  uint64(2))
	if nil != err {
		t.Error("Failed to HandleUnCandidateItem:", err)
		return
	}

	t.Log("Finish HandleUnCandidateItem ~~")

	// get Candidate
	getCandidate(blockHash2, index, t)


}

func TestStakingPlugin_Delegate(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	// Get Candidate Info
	c := getCandidate(blockHash, index, t)


	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}

	// Delegate
	_, err = delegate(state, blockHash2,  blockNumber2, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}
	t.Log("Finish Delegate ~~")
	getCandidate(blockHash2, index, t)


}


func TestStakingPlugin_WithdrewDelegate(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	// Get Candidate Info
	c := getCandidate(blockHash, index, t)

	// Delegate
	del, err := delegate(state, blockHash,  blockNumber, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	t.Log("Finish delegate ~~")
	getCandidate(blockHash, index, t)


	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}

	// Withdrew Delegate
	err = plugin.StakingInstance().WithdrewDelegate(state, blockHash2, blockNumber2, common.Big257, addrArr[index+1],
		nodeIdArr[index], blockNumber.Uint64(), del)
	if nil != err {
		t.Error("Failed to WithdrewDelegate:", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}
	t.Log("Finish WithdrewDelegate ~~")
	getCandidate(blockHash2, index, t)
}

func TestStakingPlugin_GetDelegateInfo(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}


	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	t.Log("Finish delegate ~~")
	c := getCandidate(blockHash, index, t)


	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}


	// Delegate
	_, err = delegate(state, blockHash2,  blockNumber2, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}

	t.Log("Finished Delegate ~~")
	// get Delegate info
	getDelegate(blockHash2, blockNumber.Uint64(), index, t)
}

func TestStakingPlugin_GetDelegateInfoByIrr(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	index := 1

	if err := create_staking(state, blockNumber, blockHash, index, 0, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}


	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}


	c := getCandidate(blockHash, index, t)


	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}

	t.Log("Start delegate ~~")
	// Delegate
	_, err = delegate(state, blockHash2,  blockNumber2, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}

	t.Log("Finished Delegate ~~")
	// get Delegate info
	del, err := plugin.StakingInstance().GetDelegateInfoByIrr(addrArr[index+1], nodeIdArr[index], blockNumber.Uint64())
	if nil != err {
		t.Error("Failed to GetDelegateInfoByIrr:", err)
		return
	}

	delByte, _ := json.Marshal(del)
	t.Log("Get Delegate is:", string(delByte))

}

func TestStakingPlugin_GetRelatedListByDelAddr(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}


	// staking 0, 1, 2, 3
	for i := 0; i < 4; i ++ {
		if err := create_staking(state, blockNumber, blockHash, i, 0, t); nil != err {
			t.Error("Failed to Create Staking", err)
		}
	}


	//canArr1 := make(staking.CandidateQueue, 0)
	//delArr1 := make([]*staking.Delegation, 0)

	t.Log("First delegate ~~")
	for i := 0; i < 2; i++ {
		// 0, 1
		c := getCandidate(blockHash, i, t)
		// Delegate  0, 1
		_, err := delegate(state, blockHash,  blockNumber, c, 0, i, t)
		if nil != err {
			t.Errorf("Failed to Delegate: Num: %d, error: %v", i, err)
		}

		t.Log("First: Del => Can:", addrArr[i+1].Hex(), c.NodeId.String(), c.StakingBlockNum)

		//canArr1 = append(canArr1, c)
		//delArr1 = append(delArr1, d)
	}

	//can1, _ := json.Marshal(canArr1)
	//del1, _ := json.Marshal(delArr1)
	//
	//t.Log("First Delegate of Candidate :", string(can1))
	//t.Log("First Delegate of Delegate :", string(del1))

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}



	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}


	//canArr2 := make(staking.CandidateQueue, 0)
	//delArr2 := make([]*staking.Delegation, 0)

	t.Log("Second delegate ~~")
	for i := 1; i < 3; i++ {
		// 0, 1
		c := getCandidate(blockHash2, i-1, t)
		// Delegate
		_, err := delegate(state, blockHash2,  blockNumber2, c, 0, i, t)
		if nil != err {
			t.Errorf("Failed to Delegate: Num: %d, error: %v", i, err)
		}

		t.Log("Second: Del => Can:", addrArr[i+1].Hex(), c.NodeId.String(), c.StakingBlockNum)


		//canArr2 = append(canArr2, c)
		//delArr2 = append(delArr2, d)
	}

	//can2, _ := json.Marshal(canArr2)
	//del2, _ := json.Marshal(delArr2)

	//t.Log("Second Delegate of Candidate :", string(can2))
	//t.Log("Second Delegate of Delegate :", string(del2))

	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}

	//t.Log("Finished Delegate ~~")
	//// get Delegate info
	//rel, err := plugin.StakingInstance().GetRelatedListByDelAddr(blockHash2, addrArr[1+1])
	//if nil != err {
	//	t.Error("Failed to GetRelatedListByDelAddr:", err)
	//	return
	//}
	//
	//relByte, _ := json.Marshal(rel)d
	//t.Log("Get RelatedList is:", string(relByte))




	//  test  blockHash
	prefix := staking.DelegateKeyPrefix
	iter := sndb.Ranking(blockHash, prefix, 0)

	for iter.Valid(); iter.Next(); {
		key := iter.Key()

		prefixLen := len(staking.DelegateKeyPrefix)

		nodeIdLen := discover.NodeIDBits / 8

		// delAddr
		delAddrByte := key[prefixLen: prefixLen+common.AddressLength]
		delAddr := common.BytesToAddress(delAddrByte)

		// nodeId
		nodeIdByte := key[prefixLen+common.AddressLength: prefixLen+common.AddressLength+nodeIdLen]
		nodeId := discover.MustBytesID(nodeIdByte)

		// stakenum
		stakeNumByte := key[prefixLen+common.AddressLength+nodeIdLen:]

		num := common.BytesToUint64(stakeNumByte)

		// related
		t.Log("blockHash Iter:", delAddr.Hex(), nodeId.String(), num)
	}



	//  test  blockHash2
 	prefix = staking.DelegateKeyPrefix
	iter = sndb.Ranking(blockHash2, prefix, 0)

	for iter.Valid(); iter.Next(); {
		key := iter.Key()

		prefixLen := len(staking.DelegateKeyPrefix)

		nodeIdLen := discover.NodeIDBits / 8

		// delAddr
		delAddrByte := key[prefixLen: prefixLen+common.AddressLength]
		delAddr := common.BytesToAddress(delAddrByte)

		// nodeId
		nodeIdByte := key[prefixLen+common.AddressLength: prefixLen+common.AddressLength+nodeIdLen]
		nodeId := discover.MustBytesID(nodeIdByte)

		// stakenum
		stakeNumByte := key[prefixLen+common.AddressLength+nodeIdLen:]

		num := common.BytesToUint64(stakeNumByte)

		// related
		t.Log("blockHash2 Iter:", delAddr.Hex(), nodeId.String(), num)
	}


}



func TestStakingPlugin_HandleUnDelegateItem(t *testing.T) {

}

func TestStakingPlugin_ElectNextVerifierList(t *testing.T) {

}

func TestStakingPlugin_Election(t *testing.T) {

}

func TestStakingPlugin_Switch(t *testing.T) {

}

func TestStakingPlugin_SlashCandidates(t *testing.T) {

}

func TestStakingPlugin_DeclarePromoteNotify(t *testing.T) {

}

func TestStakingPlugin_ProposalPassedNotify(t *testing.T) {

}


func TestStakingPlugin_GetCandidateONEpoch(t *testing.T) {

}

func TestStakingPlugin_GetCandidateONRound(t *testing.T) {

}


func TestStakingPlugin_GetValidatorList(t *testing.T) {

}


func TestStakingPlugin_GetVerifierList(t *testing.T) {

}


func TestStakingPlugin_ListCurrentValidatorID(t *testing.T) {

}

func TestStakingPlugin_ListVerifierNodeID(t *testing.T) {

}

func TestStakingPlugin_IsCandidate(t *testing.T) {

}

func TestStakingPlugin_IsCurrValidator(t *testing.T) {

}

func TestStakingPlugin_IsCurrVerifier(t *testing.T) {

}



// for consensus
func TestStakingPlugin_GetLastNumber(t *testing.T) {

}

func TestStakingPlugin_GetValidator(t *testing.T) {

}

func TestStakingPlugin_IsCandidateNode(t *testing.T) {

}

func TestStakingPlugin_ProbabilityElection(t *testing.T) {
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
		w := new(big.Int).Mul(new(big.Int).SetUint64(uint64(mrand.Intn(1000000000))), xcom.StakeThreshold)
		stakingWeight[1] = w.Text(10)
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

/**
Expand test cases
*/

func Test_CleanSnapshotDB (t *testing.T) {
	sndb := snapshotdb.Instance()
	sndb.Clear()
}