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
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
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



func Test_CleanSnapshotDB (t *testing.T) {
	sndb := snapshotdb.Instance()
	sndb.Clear()
}


func build_vrf_Nonce() ([]byte, [][]byte) {
	preNonces := make([][]byte, 0)
	curentNonce := crypto.Keccak256([]byte(string("nonce")))
	for i := 0; i < int(xcom.EpochValidatorNum); i++ {
		preNonces = append(preNonces, crypto.Keccak256([]byte(string(time.Now().UnixNano() + int64(i))))[:])
		time.Sleep(time.Microsecond * 10)
	}
	return curentNonce, preNonces
}

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
			ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+1)] + "balabalala" + chaList[index],
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
		t.Log("Get Delegate Info is:", string(delByte))
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

	/**
	Start Create Staking
	 */
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

	/**
	Start Get Candidate Info
	 */
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

	/**
	Start GetCandidateInfoByIrr

	Get Candidate Info
	 */
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

	/**
	Start GetCandidateList
	 */
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

	/**
	Start Edit Candidate
	 */
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

	/**
	Start IncreaseStaking
	 */
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

	/**
	Start WithdrewCandidate
	 */
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



	// Add UNStakingItems
	stakingDB := staking.NewStakingDB ()

	epoch := xutil.CalculateEpoch(blockNumber.Uint64())
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	if err := stakingDB.AddUnStakeItemStore(blockHash, epoch, addr); nil != err {
		t.Error("Failed to AddUnStakeItemStore:", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Errorf("Commit 1 err: %v", err)
	}

	// Get Candidate Info
	getCandidate(blockHash, index, t)

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
	}

	/**
	Start HandleUnCandidateItem
	 */
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

	/**
	Start Delegate
	 */
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

	/**
	Start Withdrew Delegate
	 */
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
	/**
	Start get Delegate info
	 */
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


	t.Log("First delegate ~~")
	for i := 0; i < 2; i++ {
		// 0, 1
		c := getCandidate(blockHash, i, t)
		// Delegate  0, 1
		_, err := delegate(state, blockHash,  blockNumber, c, 0, i, t)
		if nil != err {
			t.Errorf("Failed to Delegate: Num: %d, error: %v", i, err)
		}

		//t.Log("First: Del => Can:", addrArr[i+1].Hex(), c.NodeId.String(), c.StakingBlockNum)

	}


	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}



	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}


	t.Log("Second delegate ~~")
	for i := 1; i < 3; i++ {
		// 0, 1
		c := getCandidate(blockHash2, i-1, t)
		// Delegate
		_, err := delegate(state, blockHash2,  blockNumber2, c, 0, i, t)
		if nil != err {
			t.Errorf("Failed to Delegate: Num: %d, error: %v", i, err)
		}

		//t.Log("Second: Del => Can:", addrArr[i+1].Hex(), c.NodeId.String(), c.StakingBlockNum)

	}



	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}

	t.Log("Finished Delegate ~~")
	/**
	Start get RelatedList
	  */
	rel, err := plugin.StakingInstance().GetRelatedListByDelAddr(blockHash2, addrArr[1+1])
	if nil != err {
		t.Error("Failed to GetRelatedListByDelAddr:", err)
		return
	}

	relByte, _ := json.Marshal(rel)
	t.Log("Get RelatedList is:", string(relByte))
}


func TestStakingPlugin_HandleUnDelegateItem(t *testing.T) {


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


	c := getCandidate(blockHash, index, t)
	// Delegate
	_, err = delegate(state, blockHash,  blockNumber, c, 0, index, t)
	if nil != err {
		t.Error("Failed to Delegate:", err)
		return
	}


	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}

	t.Log("Finished Delegate ~~")
	// get Delegate info
	del := getDelegate(blockHash2, blockNumber.Uint64(), index, t)


	// Add UnDelegateItem
	stakingDB := staking.NewStakingDB()

	epoch := xutil.CalculateEpoch(blockNumber2.Uint64())

	amount := common.Big256

	delAddr := addrArr[index+1]

	err = stakingDB.AddUnDelegateItemStore(blockHash2, delAddr, c.NodeId, epoch, c.StakingBlockNum, amount)
	if nil != err {
		t.Error("Failed to AddUnDelegateItemStore:", err)
		return
	}
	del.Reduction = new(big.Int).Add(del.Reduction, amount)
	// update del
	if err := stakingDB.SetDelegateStore(blockHash2, delAddr, c.NodeId, c.StakingBlockNum, del); nil != err {
		t.Error("Failed to Update Delegate When AddUnDelegateItemStore:", err)
		return
	}

	if err := stakingDB.DelCanPowerStore(blockHash2, c); nil != err {
		t.Error("Failed to DelCanPowerStore:", err)
		return
	}

	// change candidate shares
	c.Shares = new(big.Int).Sub(c.Shares, amount)

	canAddr, _ := xutil.NodeId2Addr(c.NodeId)

	if err := stakingDB.SetCandidateStore(blockHash2, canAddr, c); nil != err {
		t.Error("Failed to SetCandidateStore:", err)
		return
	}

	if err := stakingDB.SetCanPowerStore(blockHash2, canAddr, c); nil != err {
		t.Error("Failed to SetCanPowerStore:", err)
		return
	}

	/**
	Start HandleUnDelegateItem
	 */
	err = plugin.StakingInstance().HandleUnDelegateItem(state, blockHash2, epoch)
	if nil != err {
		t.Error("Failed to HandleUnDelegateItem:", err)
		return
	}


	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}

	t.Log("Finished HandleUnDelegateItem ~~")

	// get Candiddate
	c = getCandidate(blockHash2, index, t)

	// get Delegate
	getDelegate(blockHash2, c.StakingBlockNum, index, t)


}

func TestStakingPlugin_ElectNextVerifierList(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList
	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue



	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}


	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}


	/*
	Start ElectNextVerifierList
	*/
	targetNum := xcom.EpochSize*xcom.ConsensusSize
	fmt.Println("targetNum:", targetNum)

	targetNumInt := big.NewInt(int64(targetNum))

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}

	err = plugin.StakingInstance().ElectNextVerifierList(blockHash2, targetNumInt.Uint64())
	if nil != err {
		t.Errorf("Failed to ElectNextVerifierList, err: %v", err)
	}

}

func TestStakingPlugin_Election(t *testing.T) {

	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue


	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	// build gensis current validatorList
	new_validatorArr := &staking.Validator_array{
		Start: start,
		End:   xcom.ConsensusSize,
	}

	new_validatorArr.Arr = queue[:int(xcom.ConsValidatorNum)]

	err = stakingDB.SetCurrentValidatorList(blockHash, new_validatorArr)
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
		sndb.Put(blockHash, xcom.NonceStorageKey, enValue)
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	/*
	Start Election
	*/
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}

	header := &types.Header{
		ParentHash:  blockHash,
		Number: big.NewInt(230),
		Nonce: types.EncodeNonce(currNonce),
	}

	err = plugin.StakingInstance().Election(blockHash2, header)
	if nil != err {
		t.Errorf("Failed to Election, err: %v", err)
	}

}

func TestStakingPlugin_Switch(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue


	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	// build gensis current validatorList
	new_validatorArr := &staking.Validator_array{
		Start: start,
		End:   xcom.ConsensusSize,
	}

	new_validatorArr.Arr = queue[:int(xcom.ConsValidatorNum)]

	err = stakingDB.SetCurrentValidatorList(blockHash, new_validatorArr)
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
		sndb.Put(blockHash, xcom.NonceStorageKey, enValue)
	}

	// Commit Block 1
	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}



	header := &types.Header{
		ParentHash:  blockHash,
		Number: big.NewInt(230),
		Nonce: types.EncodeNonce(currNonce),
	}


	err = plugin.StakingInstance().Election(blockHash2, header)
	if nil != err {
		t.Errorf("Failed to Election, err: %v", err)
	}


	// Commit Block 2
	if err := sndb.Commit(blockHash2); nil != err {
		t.Error("Commit 2 err", err)
	}

	if err := sndb.NewBlock(blockNumber3, blockHash2, blockHash3); nil != err {
		t.Error("newBlock 3 err", err)
	}

	/**
	Start Switch
	 */
	err = plugin.StakingInstance().Switch(blockHash3, big.NewInt(250).Uint64())
	if nil != err {
		t.Errorf("Failed to Switch, err: %v", err)
	}

}

func TestStakingPlugin_SlashCandidates(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	// Will be Slashing candidate
	slashQueue := make(staking.CandidateQueue, 2)


	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

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
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue


	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}


	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
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

	err = plugin.StakingInstance().SlashCandidates(state, blockHash2, blockNumber2.Uint64(), slash1.NodeId, slash1.Shares, false, staking.LowRatio)
	if nil != err {
		t.Errorf("Failed to SlashCandidates first can (LowRatio), err: %v", err)
		return
	}

	sla := new(big.Int).Div(slash2.Shares, big.NewInt(10))

	err = plugin.StakingInstance().SlashCandidates(state, blockHash2, blockNumber2.Uint64(), slash2.NodeId, sla, true, staking.DoubleSign)
	if nil != err {
		t.Errorf("Failed to SlashCandidates Second can (DoubleSign), err: %v", err)
		return
	}

}

func TestStakingPlugin_DeclarePromoteNotify(t *testing.T) {

	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}


	queue := make(staking.CandidateQueue, 0)
	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

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
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}

	/**
	Start DeclarePromoteNotify
	 */
	for i, can := range queue {
		err = plugin.StakingInstance().DeclarePromoteNotify(blockHash2, blockNumber2.Uint64(), can.NodeId, promoteVersion)
		if nil != err {
			t.Errorf("Failed to DeclarePromoteNotify, index: %d, err: %v", i, err)
			return
		}
	}

}

func TestStakingPlugin_ProposalPassedNotify(t *testing.T) {


	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}


	nodeIdArr := make([]discover.NodeID, 0)
	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}

		if i < 20 {
			nodeIdArr = append(nodeIdArr, canTmp.NodeId)
		}
	}

	// Commit Block 1
	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}

	/**
	Start ProposalPassedNotify
	 */
	err = plugin.StakingInstance().ProposalPassedNotify(blockHash2, blockNumber2.Uint64(), nodeIdArr, promoteVersion)
	if nil != err {
		t.Errorf("Failed to ProposalPassedNotify, err: %v", err)
		return
	}
}


func TestStakingPlugin_GetCandidateONEpoch(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList
	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue



	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}


	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}


	/**
	Start GetCandidateONEpoch
	 */
	canQueue, err := plugin.StakingInstance().GetCandidateONEpoch(blockHash2, blockNumber2.Uint64(), plugin.QueryStartNotIrr)
	if nil != err {
		t.Errorf("Failed to GetCandidateONEpoch by QueryStartNotIrr, err: %v", err)
		return
	}

	canArr, _ := json.Marshal(canQueue)
	t.Log("GetCandidateONEpoch by QueryStartNotIrr:", string(canArr))

	canQueue, err = plugin.StakingInstance().GetCandidateONEpoch(blockHash2, blockNumber2.Uint64(), plugin.QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetCandidateONEpoch by QueryStartIrr, err: %v", err)
		return
	}
	canArr, _ = json.Marshal(canQueue)
	t.Log("GetCandidateONEpoch by QueryStartIrr:", string(canArr))
}

func TestStakingPlugin_GetCandidateONRound(t *testing.T) {


	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue


	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	// build gensis current validatorList
	new_validatorArr := &staking.Validator_array{
		Start: start,
		End:   xcom.ConsensusSize,
	}

	new_validatorArr.Arr = queue[:int(xcom.ConsValidatorNum)]

	err = stakingDB.SetCurrentValidatorList(blockHash, new_validatorArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis current round validatorList, err: %v", err)
		return
	}



	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}




	/**
	Start GetCandidateONRound
	 */
	canQueue, err := plugin.StakingInstance().GetCandidateONRound(blockHash2, blockNumber2.Uint64(), plugin.CurrentRound,  plugin.QueryStartNotIrr)
	if nil != err {
		t.Errorf("Failed to GetCandidateONRound by QueryStartNotIrr, err: %v", err)
		return
	}

	canArr, _ := json.Marshal(canQueue)
	t.Log("GetCandidateONRound by QueryStartNotIrr:", string(canArr))

	canQueue, err = plugin.StakingInstance().GetCandidateONRound(blockHash2, blockNumber2.Uint64(), plugin.CurrentRound, plugin.QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetCandidateONRound by QueryStartIrr, err: %v", err)
		return
	}
	canArr, _ = json.Marshal(canQueue)
	t.Log("GetCandidateONRound by QueryStartIrr:", string(canArr))

}


func TestStakingPlugin_GetValidatorList(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue


	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	// build gensis current validatorList
	new_validatorArr := &staking.Validator_array{
		Start: start,
		End:   xcom.ConsensusSize,
	}

	new_validatorArr.Arr = queue[:int(xcom.ConsValidatorNum)]

	err = stakingDB.SetCurrentValidatorList(blockHash, new_validatorArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis current round validatorList, err: %v", err)
		return
	}



	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}


	/**
	Start  GetValidatorList
	 */
	validatorExQueue, err := plugin.StakingInstance().GetValidatorList(blockHash2, blockNumber2.Uint64(), plugin.CurrentRound,  plugin.QueryStartNotIrr)
	if nil != err {
		t.Errorf("Failed to GetValidatorList by QueryStartNotIrr, err: %v", err)
		return
	}

	validatorExArr, _ := json.Marshal(validatorExQueue)
	t.Log("GetValidatorList by QueryStartNotIrr:", string(validatorExArr))

	validatorExQueue, err = plugin.StakingInstance().GetValidatorList(blockHash2, blockNumber2.Uint64(), plugin.CurrentRound, plugin.QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetValidatorList by QueryStartIrr, err: %v", err)
		return
	}
	validatorExArr, _ = json.Marshal(validatorExQueue)
	t.Log("GetValidatorList by QueryStartIrr:", string(validatorExArr))

}


func TestStakingPlugin_GetVerifierList(t *testing.T) {

	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue


	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}

	/**
	Start GetVerifierList
	 */
	validatorExQueue, err := plugin.StakingInstance().GetVerifierList(blockHash2, blockNumber2.Uint64(),  plugin.QueryStartNotIrr)
	if nil != err {
		t.Errorf("Failed to GetVerifierList by QueryStartNotIrr, err: %v", err)
		return
	}

	validatorExArr, _ := json.Marshal(validatorExQueue)
	t.Log("GetVerifierList by QueryStartNotIrr:", string(validatorExArr))

	validatorExQueue, err = plugin.StakingInstance().GetVerifierList(blockHash2, blockNumber2.Uint64(), plugin.QueryStartIrr)
	if nil != err {
		t.Errorf("Failed to GetVerifierList by QueryStartIrr, err: %v", err)
		return
	}
	validatorExArr, _ = json.Marshal(validatorExQueue)
	t.Log("GetVerifierList by QueryStartIrr:", string(validatorExArr))

}


func TestStakingPlugin_ListCurrentValidatorID(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue


	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	// build gensis current validatorList
	new_validatorArr := &staking.Validator_array{
		Start: start,
		End:   xcom.ConsensusSize,
	}

	new_validatorArr.Arr = queue[:int(xcom.ConsValidatorNum)]

	err = stakingDB.SetCurrentValidatorList(blockHash, new_validatorArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis current round validatorList, err: %v", err)
		return
	}



	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}


	/**
	Start  ListCurrentValidatorID
	 */
	validatorIdQueue, err := plugin.StakingInstance().ListCurrentValidatorID(blockHash2, blockNumber2.Uint64())
	if nil != err {
		t.Errorf("Failed to ListCurrentValidatorID, err: %v", err)
		return
	}

	validatorIdArr, _ := json.Marshal(validatorIdQueue)
	t.Log("ListCurrentValidatorID:", string(validatorIdArr))


}

func TestStakingPlugin_ListVerifierNodeID(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}

	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue

	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}


	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}


	/**
	Start  ListVerifierNodeID
	 */
	validatorIdQueue, err := plugin.StakingInstance().ListVerifierNodeID(blockHash2, blockNumber2.Uint64())
	if nil != err {
		t.Errorf("Failed to ListVerifierNodeID, err: %v", err)
		return
	}

	validatorIdArr, _ := json.Marshal(validatorIdQueue)
	t.Log("ListVerifierNodeID:", string(validatorIdArr))
}

func TestStakingPlugin_IsCandidate(t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	nodeIdArr := make([]discover.NodeID, 0)

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

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
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}


	/**
	Start  IsCandidate
	 */
	 for i, nodeId := range nodeIdArr {
		 yes, err := plugin.StakingInstance().IsCandidate(blockHash2, nodeId, plugin.QueryStartNotIrr)
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



	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	nodeIdArr := make([]discover.NodeID, 0)

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}

		if i < 20 {
			nodeIdArr = append(nodeIdArr, canTmp.NodeId)
		}
	}



	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue


	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}

	// build gensis current validatorList
	new_validatorArr := &staking.Validator_array{
		Start: start,
		End:   xcom.ConsensusSize,
	}

	new_validatorArr.Arr = queue[:int(xcom.ConsValidatorNum)]

	err = stakingDB.SetCurrentValidatorList(blockHash, new_validatorArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis current round validatorList, err: %v", err)
		return
	}

	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}




	/**
	Start  IsCurrValidator
	 */
	for i, nodeId := range nodeIdArr {
		yes, err := plugin.StakingInstance().IsCurrValidator(blockHash2, nodeId, plugin.QueryStartNotIrr)
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



	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()
	xcom.NewVrfHandler(common.ZeroHash.Bytes())

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	nodeIdArr := make([]discover.NodeID, 0)

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}

		if i < 20 {
			nodeIdArr = append(nodeIdArr, canTmp.NodeId)
		}
	}



	stakingDB := staking.NewStakingDB()

	// build genesis VerifierList

	start := uint64(1)
	end := xcom.EpochSize*xcom.ConsensusSize

	new_verifierArr := &staking.Validator_array{
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
		if uint64(count) == xcom.EpochValidatorNum {
			break
		}
		addrSuffix := iter.Value()
		var can *staking.Candidate

		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash, addrSuffix)
		if nil != err {
			t.Error("Failed to ElectNextVerifierList", "canAddr", common.BytesToAddress(addrSuffix).Hex(), "err", err)
			return
		}

		addr := common.BytesToAddress(addrSuffix)

		powerStr := [staking.SWeightItem]string{fmt.Sprint(can.ProcessVersion), can.Shares.String(),
			fmt.Sprint(can.StakingBlockNum), fmt.Sprint(can.StakingTxIndex)}

		val := &staking.Validator{
			NodeAddress:   addr,
			NodeId:        can.NodeId,
			StakingWeight: powerStr,
			ValidatorTerm: 0,
		}
		queue = append(queue, val)
		count ++
	}

	new_verifierArr.Arr = queue


	err = stakingDB.SetVerfierList(blockHash, new_verifierArr)
	if nil != err {
		t.Errorf("Failed to Set Genesis VerfierList, err: %v", err)
		return
	}



	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}
	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}




	/**
	Start  IsCurrVerifier
	 */
	for i, nodeId := range nodeIdArr {
		yes, err := plugin.StakingInstance().IsCurrVerifier(blockHash2, nodeId, plugin.QueryStartNotIrr)
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
		stakingWeight := [staking.SWeightItem]string{}
		stakingWeight[0] = strconv.Itoa(mrand.Intn(5)+1)
		v1 := new(big.Int).SetInt64(time.Now().UnixNano())
		v1.Mul(v1, new(big.Int).SetInt64(1e18))
		v1.Add(v1, new(big.Int).SetInt64(int64(mrand.Intn(1000))))
		stakingWeight[1] = v1.Text(10)
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

func Test_IteratorCandidate (t *testing.T) {
	defer func() {
		sndb.Clear()
	}()

	state, err := newChainState()
	if nil != err {
		t.Error("Failed to build the state", err)
	}
	newPlugins()

	build_gov_data(state)

	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
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


		canTmp := &staking.Candidate{
			NodeId:          nodeId,
			StakingAddress:  sender,
			BenifitAddress:  addr,
			StakingBlockNum: uint64(i),
			StakingTxIndex:  uint32(index),
			Shares:          balance,

			// Prevent null pointer initialization
			Released: common.Big0,
			ReleasedHes: common.Big0,
			RestrictingPlan: common.Big0,
			RestrictingPlanHes: common.Big0,

			Description: staking.Description{
				NodeName:   nodeNameArr[index] + "_" + fmt.Sprint(i),
				ExternalId: nodeNameArr[index] + chaList[(len(chaList)-1)%(index+ii+1)] + "balabalala" + chaList[index],
				Website:    "www." + nodeNameArr[index] + "_" + fmt.Sprint(i) + ".org",
				Details:    "This is " + nodeNameArr[index] + "_" + fmt.Sprint(i) + " Super Node",
			},
		}

		canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)
		err = plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance,
			initProcessVersion, 0, canAddr, canTmp)

		if nil != err {
			t.Errorf("Failed to Create Staking, num: %d, err: %v", i, err)
			return
		}
	}


	// commit
	if err := sndb.Commit(blockHash); nil != err {
		t.Error("Commit 1 err", err)
	}


	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock 2 err", err)
	}


	stakingDB := staking.NewStakingDB()

	iter := stakingDB.IteratorCandidatePowerByBlockHash(blockHash2, 0)

	queue := make(staking.CandidateQueue, 0)

	for iter.Valid(); iter.Next(); {
		addrSuffix := iter.Value()
		can, err := stakingDB.GetCandidateStoreWithSuffix(blockHash2, addrSuffix)
		if nil != err {
			t.Errorf("Failed to Iterator Candidate info, err: %v", err)
			return
		}

		val := fmt.Sprint(initProcessVersion) + "_" + can.Shares.String() + "_" + fmt.Sprint(can.StakingBlockNum) + "_" + fmt.Sprint(can.StakingTxIndex)
		t.Log("Val:", val)

		queue = append(queue, can)
	}

	arrJson, _ := json.Marshal(queue)
	t.Log("CandidateList:", string(arrJson))
	t.Log("Candidate queue length:", len(queue))
}


func Test_Iterator (t *testing.T) {

	defer func() {
		sndb.Clear()
	}()


	sndb := snapshotdb.Instance()

	if err := sndb.NewBlock(blockNumber, common.ZeroHash, blockHash); nil != err {
		t.Error("newBlock err", err)
	}

	initProcessVersion := uint32(1<<16 | 0<<8 | 0) // 65536

	for i := 0; i < 1000; i ++ {

		var index int
		if i >= len(balanceStr) {
			index = i%(len(balanceStr)-1)
		}


		balance, _ := new(big.Int).SetString(balanceStr[index], 10)

		mrand.Seed(time.Now().UnixNano())

		weight := mrand.Intn(1000000000)


		balance = new(big.Int).Add(balance, big.NewInt(int64(weight)))


		key := staking.TallyPowerKey(balance, uint64(i), uint32(index), initProcessVersion)
		val := fmt.Sprint(initProcessVersion) + "_" + balance.String() + "_" + fmt.Sprint(i) + "_" + fmt.Sprint(index)
		sndb.Put(blockHash, key, []byte(val))
	}


	// iter
	iter := sndb.Ranking(blockHash, staking.CanPowerKeyPrefix, 0)
	if err := iter.Error(); nil != err {
		t.Errorf("Failed to interator, err: %v", err)
		return
	}
	defer iter.Release()
	for iter.Valid(); iter.Next(); {
		t.Log("Value:=", string(iter.Value()))
	}

}