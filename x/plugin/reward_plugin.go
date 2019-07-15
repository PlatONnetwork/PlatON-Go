package plugin

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/reward"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type rewardMgrPlugin struct {

}

var (
	rm *rewardMgrPlugin = nil
)

func RewardMgrInstance() *rewardMgrPlugin {
	if rm == nil {
		rm = & rewardMgrPlugin {}
	}
	return rm
}

//func ClearRewardPlugin() error {
//	if nil == rm {
//		return common.NewSysError("the RewardPlugin already be nil")
//	}
//	rm = nil
//	return nil
//}

// BeginBlock does something like check input params before execute transactions,
// in rewardMgrPlugin it does nothing.
func (rmp *rewardMgrPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	return nil
}

// EndBlock will handle reward work, if it's time to settle, reward staking. Then reward worker
// for create new block, this is necessary. At last if current block is the last block at the end
// of year, increasing issuance.
func (rmp *rewardMgrPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	blockNumber := head.Number.Uint64()
	year := xutil.CalculateYear(blockNumber)
	stakingReward, newBlockReward, _ := rmp.calculateExpectReward(uint32(year), state)

	if xutil.IsSettlementPeriod(blockNumber) {
		log.Info("ready to reward staking", "period", xutil.CalculateRound(blockNumber))
		if err:= rmp.rewardStaking(head, stakingReward, state); err != nil {
			return err
		}
	}

	log.Debug("ready to reward new block", "blockNumber", blockNumber, "hash", head.Hash())
	if err := rmp.rewardNewBlock(head, newBlockReward, state); err != nil {
		return err
	}

	if xutil.IsYearEnd(blockNumber) {
		log.Info("ready to increase issuance", "blockNumber", blockNumber, "hash", head.Hash(), "year", year)
		if err := rmp.increaseIssuance(uint32(year), state); err != nil {
			return err
		}
	}

	return nil
}

// Confirmed does nothing
func (rmp *rewardMgrPlugin) Confirmed(block *types.Block) error {
	return nil
}

// increaseIssuance used for increase issuance at the end of each year
func (rmp *rewardMgrPlugin) increaseIssuance(year uint32, state xcom.StateDB) error {
	var (
		rate = new(big.Int)
		temp = new(big.Int)
	)

	// get historical issuance
	histIssuance := GetLatestCumulativeIssue(state)

	// every year will increase 2.5 percent of historical issuance, and 1/5 send to
	// community developer foundation, the left send to reward manage pool
	currIssuance := temp.Div(histIssuance, rate.SetUint64(40))
	develop := temp.Div(currIssuance, rate.SetUint64(5))
	rewards := temp.Sub(currIssuance, develop)

	histIssuance = temp.Add(histIssuance, currIssuance)
	SetYearEndCumulativeIssue(state, year, histIssuance)

	state.AddBalance(vm.CommunityDeveloperFoundation, develop)
	state.AddBalance(vm.RewardManagerPoolAddr, rewards)

	return nil
}

// rewardStaking used for reward staking at the settle block
func (rmp *rewardMgrPlugin) rewardStaking(head *types.Header, reward *big.Int, state xcom.StateDB) error {
	blockHash := head.Hash()
	blockNumber := head.Number.Uint64()

	list, err := StakingInstance().GetVerifierList(blockHash, blockNumber, false)
	if err != nil {
		log.Debug("get verifier list failed in rewardStaking", "hash", blockHash, "blockNumber", blockNumber)
		return err
	}

	for index := 0; index < len(list); index++ {
		addr := list[index].BenifitAddress

		if addr != vm.RewardManagerPoolAddr {
			state.SubBalance(vm.RewardManagerPoolAddr, reward)
			state.AddBalance(addr, reward)
		}
	}

	return nil
}

// rewardNewBlock used for reward new block. it returns coinbase and error
func (rmp *rewardMgrPlugin) rewardNewBlock(head *types.Header, reward *big.Int, state xcom.StateDB) error {
	rewardAddr := head.Coinbase
	if rewardAddr != vm.RewardManagerPoolAddr {
		state.SubBalance(vm.RewardManagerPoolAddr, reward)
		state.AddBalance(rewardAddr, reward)
	}
	return nil
}

// calculateExpectReward used for calculate the stakingReward and newBlockReward that should be send in each corresponding period
func (rmp *rewardMgrPlugin) calculateExpectReward(year uint32, state xcom.StateDB) (*big.Int, *big.Int, error) {
	var (
		stakingReward  = new(big.Int)
		newBlockReward = new(big.Int)
		temp           = new(big.Int)
	)

	expectNewBlocks := int64(365) * 24 * 3600 / 1
	expectEpochs := int64(365) * 24 * 3600 / int64(xcom.ConsensusSize() * xcom.EpochSize())

	issuance := GetLatestCumulativeIssue(state)
	totalNewBlockReward := temp.Div(issuance, big.NewInt(5))
	totalStakingReward := temp.Sub(issuance, totalNewBlockReward)

	newBlockReward = temp.Div(totalNewBlockReward, big.NewInt(expectNewBlocks))
	stakingReward = temp.Div(totalStakingReward, big.NewInt(expectEpochs))

	return stakingReward, newBlockReward, nil
}

// SetYearEndCumulativeIssue used for set historical cumulative increase at the end of the year
func SetYearEndCumulativeIssue(state xcom.StateDB, year uint32, total *big.Int) {
	yearEndIncreaseKey := reward.GetHistoryIncreaseKey(year)
	state.SetState(vm.RewardManagerPoolAddr, yearEndIncreaseKey, total.Bytes())
}

// GetLatestCumulativeIssue used for get the cumulative issuance in the most recent year
func GetLatestCumulativeIssue(state xcom.StateDB) *big.Int {
	var issue = new(big.Int)
	// !!!
	// latestYear := getLastYear()
	// !!!
	latestYear := uint32(0)
	LastYearIncreaseKey := reward.GetHistoryIncreaseKey(latestYear)
	bIssue := state.GetState(vm.RewardManagerPoolAddr, LastYearIncreaseKey)
	return issue.SetBytes(bIssue)
}

// GetHistoryCumulativeIssue used for get the cumulative issuance of a certain year in history
func GetHistoryCumulativeIssue(state xcom.StateDB, year uint32) *big.Int {
	var issue = new(big.Int)
	histIncreaseKey := reward.GetHistoryIncreaseKey(year)
	bIssue := state.GetState(vm.RewardManagerPoolAddr, histIncreaseKey)
	return issue.SetBytes(bIssue)
}
