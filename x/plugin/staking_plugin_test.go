package plugin_test

import (
	"testing"
	"math/big"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/x/plugin"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"encoding/json"
	"fmt"
)


func create_staking (blockNumber *big.Int, blockHash common.Hash, state *state.StateDB, index int, typ uint16, t *testing.T) error {

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
			ExternalId: nodeNameArr[index] + chaList[len(chaList)%index] + "balabalala" + chaList[index],
			Website:    "www." + nodeNameArr[index] + ".org",
			Details:    "This is " + nodeNameArr[index] + " Super Node",
		},
	}

	canAddr, _ := xutil.NodeId2Addr(canTmp.NodeId)

	return plugin.StakingInstance().CreateCandidate(state, blockHash, blockNumber, balance, initProcessVersion, typ, canAddr, canTmp)
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


	if err := create_staking(blockNumber, blockHash, state, 1, 1, t); nil != err {
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

	if err := create_staking(blockNumber, blockHash, state, index, 1, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	sndb.Commit(blockHash)

	// Get Candidate Info
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])
	if can, err := plugin.StakingInstance().GetCandidateInfoByIrr(addr); nil != err {
		t.Error("Failed to GetCandidateInfoByIrr", err)
	}else {

		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
	}

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

	if err := create_staking(blockNumber, blockHash, state, index, 1, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	sndb.Commit(blockHash)

	// Get Candidate Info
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])
	if can, err := plugin.StakingInstance().GetCandidateInfo(blockHash, addr); nil != err {
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

	index := 1


	for i := 0; i < 4; i++ {
		if err := create_staking(blockNumber, blockHash, state, index, 1, t); nil != err {
			t.Error("Failed to Create num: " + fmt.Sprint(i) + " Staking", err)
		}
	}

	sndb.Commit(blockHash)

	if queue, err := plugin.StakingInstance().GetCandidateList(blockHash, plugin.QueryStartIrr); nil != err {
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

	if err := create_staking(blockNumber, blockHash, state, index, 1, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	sndb.Commit(blockHash)

	// Get Candidate Info
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	var c *staking.Candidate
	if can, err := plugin.StakingInstance().GetCandidateInfo(blockHash, addr); nil != err {
		t.Error("Failed to Get Candidate info", err)
	}else {

		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

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

	if err := create_staking(blockNumber, blockHash, state, index, 1, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	sndb.Commit(blockHash)

	// Get Candidate Info
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	var c *staking.Candidate
	if can, err := plugin.StakingInstance().GetCandidateInfo(blockHash, addr); nil != err {
		t.Error("Failed to Get Candidate info", err)
	}else {

		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
	}

	// IncreaseStaking
	if err := plugin.StakingInstance().IncreaseStaking(state, blockHash2, blockNumber2, common.Big256, uint16(1), c); nil != err {
		t.Error("Failed to IncreaseStaking", err)
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

	if err := create_staking(blockNumber, blockHash, state, index, 1, t); nil != err {
		t.Error("Failed to Create Staking", err)
	}

	sndb.Commit(blockHash)

	// Get Candidate Info
	addr, _ := xutil.NodeId2Addr(nodeIdArr[index])

	var c *staking.Candidate
	if can, err := plugin.StakingInstance().GetCandidateInfo(blockHash, addr); nil != err {
		t.Error("Failed to Get Candidate info", err)
	}else {

		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

	if err := sndb.NewBlock(blockNumber2, blockHash, blockHash2); nil != err {
		t.Error("newBlock2 err", err)
	}

	// IncreaseStaking
	if err := plugin.StakingInstance().WithdrewCandidate(state, blockHash2, blockNumber2,  c); nil != err {
		t.Error("Failed to WithdrewCandidate", err)
	}

	// get Candidate info
	if can, err := plugin.StakingInstance().GetCandidateInfo(blockHash2, addr); nil != err {
		t.Error("Failed to Get Candidate info", err)
	}else {

		canByte, _ := json.Marshal(can)
		t.Log("Get Candidate Info is:", string(canByte))
		c = can
	}

}

func TestStakingPlugin_HandleUnCandidateItem(t *testing.T) {

}

func TestStakingPlugin_GetDelegateInfo(t *testing.T) {

}

func TestStakingPlugin_GetDelegateInfoByIrr(t *testing.T) {

}

func TestStakingPlugin_GetRelatedListByDelAddr(t *testing.T) {

}


func TestStakingPlugin_Delegate(t *testing.T) {

}

func TestStakingPlugin_WithdrewDelegate(t *testing.T) {

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

/**
Expand test cases
*/