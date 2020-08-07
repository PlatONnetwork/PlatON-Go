package statsdb

import (
	"math/big"
	"testing"

	"gotest.tools/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	address     = common.MustBech32ToAddress("lax1e8su9veseal8t8eyj0zuw49nfkvtqlun2sy6wj")
	nodeAddress = common.NodeAddress(address)
	nodeId      = common.NodeID(discover.MustHexID("0x362003c50ed3a523cdede37a001803b8f0fed27cb402b3d6127a1a96661ec202318f68f4c76d9b0bfbabfd551a178d4335eaeaa9b7981a4df30dfc8c0bfe3384"))
	blockNo     = big.NewInt(int64(234))
)

func buildExeBlockData() *common.ExeBlockData {
	blockNumber := blockNo.Uint64()
	common.InitExeBlockData(blockNumber)

	candidate := &common.CandidateInfo{nodeId, address}
	candidateInfoList := []*common.CandidateInfo{candidate}

	common.CollectRestrictingReleaseItem(blockNumber, address, big.NewInt(111), big.NewInt(0))
	common.CollectStakingFrozenItem(blockNumber, nodeId, nodeAddress, 222, false)
	common.CollectDuplicatedSignSlashingSetting(blockNumber, 2000, 60)

	rewardData := &common.RewardData{BlockRewardAmount: big.NewInt(111), StakingRewardAmount: big.NewInt(111), CandidateInfoList: candidateInfoList}
	common.CollectRewardData(blockNumber, rewardData)

	return common.GetExeBlockData(blockNumber)
}

func Test_DB(t *testing.T) {
	SetDBPath("d:\\swap\\statsdb")

	blockData := buildExeBlockData()
	Instance().WriteExeBlockData(blockNo, blockData)

	data := Instance().ReadExeBlockData(blockNo)
	assert.Equal(t, blockData.RewardData.CandidateInfoList[0].NodeID, data.RewardData.CandidateInfoList[0].NodeID)

	Instance().DeleteExeBlockData(blockNo)

	data2 := Instance().ReadExeBlockData(blockNo)
	if data2 != nil {
		t.Fatal("db item not deleted")
	} else {
		t.Log("db item deleted")
	}
}
