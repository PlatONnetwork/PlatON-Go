package plugin

import (
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/reward"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

const RewardNewBlockRate = 4 //

type rewardMgrPlugin struct {
}

type issuanceData struct {
	history *big.Int // the counts of historical issuance
	this    *big.Int // the counts of issuance at this year end
	develop *big.Int // the counts of issuance at this year end used for developer foundation
	reward  *big.Int // the counts of issuance at this year end used for reward pool
}

var (
	rewardOnce sync.Once
	rm         *rewardMgrPlugin = nil
)

func RewardMgrInstance() *rewardMgrPlugin {
	rewardOnce.Do(func() {
		rm = &rewardMgrPlugin{}
	})
	return rm
}

/*func ClearRewardPlugin() error {
	if nil == rm {
		return common.NewSysError("the RewardPlugin already be nil")
	}
	rm = nil
	return nil
}*/

// BeginBlock does something like check input params before execute transactions,
// in rewardMgrPlugin it does nothing.
func (rmp *rewardMgrPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	return nil
}

// EndBlock will handle reward work, if it's time to settle, reward staking. Then reward worker
// for create new block, this is necessary. At last if current block is the last block at the end
// of year, increasing issuance.
func (rmp *rewardMgrPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {

	var lastYear uint32
	blockNumber := head.Number.Uint64()
	thisYear := xutil.CalculateYear(blockNumber)

	if thisYear == 0 {
		lastYear = 0
	} else {
		lastYear = thisYear - 1
	}
	// every year will increase 2.5 percent of historical issuance, and 1/5 send to
	// community developer foundation, the left send to reward manage pool
	histIssuance := GetHistoryCumulativeIssue(state, lastYear)
	currIssuance := new(big.Int).Div(histIssuance, big.NewInt(40))
	devIssuance := new(big.Int).Div(currIssuance, big.NewInt(5))
	rewardIssuance := new(big.Int).Sub(currIssuance, devIssuance)

	stakingReward, newBlockReward := rmp.calculateExpectReward(rewardIssuance, state)

	log.Trace("show calculated data", "year", thisYear, "staking", stakingReward, "newBlock", newBlockReward)

	if xutil.IsSettlementPeriod(blockNumber) {
		log.Info("ready to reward staking", "period", xutil.CalculateRound(blockNumber))
		if err := rmp.rewardStaking(head, stakingReward, state); err != nil {
			return err
		}
	}

	log.Debug("ready to reward new block", "blockNumber", blockNumber, "hash", head.Hash())
	if err := rmp.rewardNewBlock(head, newBlockReward, state); err != nil {
		return err
	}

	if xutil.IsYearEnd(blockNumber) {
		log.Info("ready to increase issuance", "blockNumber", blockNumber, "hash", head.Hash(), "year", thisYear)

		issuance := &issuanceData{
			history: histIssuance,
			this:    currIssuance,
			develop: devIssuance,
			reward:  rewardIssuance,
		}
		rmp.increaseIssuance(issuance, thisYear, state)
	}

	return nil
}

// Confirmed does nothing
func (rmp *rewardMgrPlugin) Confirmed(block *types.Block) error {
	return nil
}

// increaseIssuance used for increase issuance at the end of each year
func (rmp *rewardMgrPlugin) increaseIssuance(issuance *issuanceData, year uint32, state xcom.StateDB) {
	state.AddBalance(vm.CommunityDeveloperFoundation, issuance.develop)
	state.AddBalance(vm.RewardManagerPoolAddr, issuance.reward)

	// restore the cumulative issue at this year end
	issuance.history.Add(issuance.history, issuance.this)
	SetYearEndCumulativeIssue(state, year, issuance.history)
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
func (rmp *rewardMgrPlugin) calculateExpectReward(reward *big.Int, state xcom.StateDB) (*big.Int, *big.Int) {
	var (
		totalReward         = common.Big0
		totalNewBlockReward = common.Big0
		totalStakingReward  = common.Big0
	)

	// get expected settlement epochs and new blocks per year first
	epochs := xutil.EpochsPerYear()
	blocks := xutil.CalcBlocksEachYear()

	// total rewards are the balance of reward pool and the new reward
	balance := state.GetBalance(vm.RewardManagerPoolAddr)
	totalReward = totalReward.Add(reward, balance)

	// 1/4 of total reward is used for reward create new block, the left is used for staking
	totalNewBlockReward = totalNewBlockReward.Div(totalReward, big.NewInt(RewardNewBlockRate))
	totalStakingReward = totalStakingReward.Sub(totalReward, totalNewBlockReward)

	newBlockReward := totalNewBlockReward.Div(totalNewBlockReward, big.NewInt(int64(blocks)))
	stakingReward := totalStakingReward.Div(totalStakingReward, big.NewInt(int64(epochs)))

	return stakingReward, newBlockReward
}

// SetYearEndCumulativeIssue used for set historical cumulative increase at the end of the year
func SetYearEndCumulativeIssue(state xcom.StateDB, year uint32, total *big.Int) {
	yearEndIncreaseKey := reward.GetHistoryIncreaseKey(year)
	state.SetState(vm.RewardManagerPoolAddr, yearEndIncreaseKey, total.Bytes())
}

// GetHistoryCumulativeIssue used for get the cumulative issuance of a certain year in history
func GetHistoryCumulativeIssue(state xcom.StateDB, year uint32) *big.Int {
	var issue = new(big.Int)
	histIncreaseKey := reward.GetHistoryIncreaseKey(year)
	bIssue := state.GetState(vm.RewardManagerPoolAddr, histIncreaseKey)
	return issue.SetBytes(bIssue)
}
