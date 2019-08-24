package plugin

import (
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

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

const (
	LessThanFoundationYearDeveloperRate    = 100
	AfterFoundationYearDeveloperRewardRate = 50
	AfterFoundationYearFoundRewardRate     = 50

	IncreaseIssue = 40

	RewardPoolIncreaseRate = 80 // 80% of fixed-issued tokens are allocated to reward pool each year
)

var (
	rewardOnce sync.Once
	rm         *rewardMgrPlugin = nil
)

func RewardMgrInstance() *rewardMgrPlugin {
	rewardOnce.Do(func() {
		log.Info("Init Reward plugin ...")
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

	blockNumber := head.Number.Uint64()

	log.Debug("Call EndBlock on reward_plugin start", "blockNumber", blockNumber, "blockHash", blockHash.Hex())

	thisYear := xutil.CalculateYear(blockNumber)
	var lastYear uint32
	if thisYear != 0 {
		lastYear = thisYear - 1
	}

	stakingReward, newBlockReward := rmp.calculateExpectReward(thisYear, lastYear, state)

	log.Debug("Call EndBlock on reward_plugin: after call calculateExpectReward", "blockNumber", blockNumber,
		"blockHash", blockHash.Hex(), "stakingReward", stakingReward, "packageBlockReward", newBlockReward)

	if xutil.IsSettlementPeriod(blockNumber) {
		if err := rmp.allocateStakingReward(blockNumber, blockHash, stakingReward, state); err != nil {
			return err
		}
	}

	rmp.allocatePackageBlock(blockNumber, blockHash, head.Coinbase, newBlockReward, state)

	// the block at the end of each year, additional issuance
	if xutil.IsYearEnd(blockNumber) {
		rmp.increaseIssuance(thisYear, lastYear, state)
	}

	log.Debug("Call EndBlock on reward_plugin End")

	return nil
}

// Confirmed does nothing
func (rmp *rewardMgrPlugin) Confirmed(block *types.Block) error {
	return nil
}

func (rmp *rewardMgrPlugin) isLessThanFoundationYear(thisYear uint32) bool {
	if thisYear < xcom.PlatONFoundationYear() {
		return true
	}
	return false
}

func (rmp *rewardMgrPlugin) addPlatONFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	platonFoundationIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(vm.PlatONFoundationAddress, platonFoundationIncr)
}

func (rmp *rewardMgrPlugin) addCommunityDeveloperFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	developerFoundationIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(vm.CommunityDeveloperFoundation, developerFoundationIncr)
}
func (rmp *rewardMgrPlugin) addRewardPoolIncreaseIssuance(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	rewardpoolIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(vm.RewardManagerPoolAddr, rewardpoolIncr)
}

// increaseIssuance used for increase issuance at the end of each year
func (rmp *rewardMgrPlugin) increaseIssuance(thisYear, lastYear uint32, state xcom.StateDB) {
	var currIssuance *big.Int
	//issuance increase
	{
		histIssuance := GetHistoryCumulativeIssue(state, lastYear)
		currIssuance = new(big.Int).Div(histIssuance, big.NewInt(IncreaseIssue)) // 2.5% increase in the previous year
		// Restore the cumulative issue at this year end
		histIssuance.Add(histIssuance, currIssuance)
		SetYearEndCumulativeIssue(state, thisYear, histIssuance)
		log.Debug("Call EndBlock on reward_plugin: increase issuance", "thisYear", thisYear, "addIssuance", currIssuance, "hit", histIssuance)

	}
	rmp.addRewardPoolIncreaseIssuance(state, currIssuance, RewardPoolIncreaseRate)

	if rmp.isLessThanFoundationYear(thisYear) {
		rmp.addCommunityDeveloperFoundation(state, currIssuance, LessThanFoundationYearDeveloperRate)
	} else {
		rmp.addCommunityDeveloperFoundation(state, currIssuance, AfterFoundationYearDeveloperRewardRate)
		rmp.addPlatONFoundation(state, currIssuance, AfterFoundationYearFoundRewardRate)
	}
	balance := state.GetBalance(vm.RewardManagerPoolAddr)
	SetYearEndBalance(state, thisYear, balance)

}

// allocateStakingReward used for reward staking at the settle block
func (rmp *rewardMgrPlugin) allocateStakingReward(blockNumber uint64, blockHash common.Hash, reward *big.Int, state xcom.StateDB) error {

	log.Info("Allocate staking reward start", "blockNumber", blockNumber, "hash", blockHash,
		"epoch", xutil.CalculateEpoch(blockNumber), "reward", reward)

	list, err := StakingInstance().GetVerifierList(blockHash, blockNumber, false)
	if err != nil {
		log.Error("Failed to allocateStakingReward: call GetVerifierList is failed", "blockNumber", blockNumber, "hash", blockHash, "err", err)
		return err
	}
	rmp.rewardStakingByValidatorList(state, list, reward)
	return nil
}

func (rmp *rewardMgrPlugin) rewardStakingByValidatorList(state xcom.StateDB, list staking.ValidatorExQueue, reward *big.Int) {
	validatorNum := int64(len(list))
	everyValidatorReward := new(big.Int).Div(reward, big.NewInt(validatorNum))

	log.Debug("calculate validator staking reward", "validator length", len(list), "everyOneReward", everyValidatorReward)
	totalValidatorReward := new(big.Int)
	for _, value := range list {
		addr := value.BenefitAddress
		if addr != vm.RewardManagerPoolAddr {

			log.Debug("allocate staking reward one-by-one", "nodeId", value.NodeId.String(),
				"benefitAddress", addr.String(), "staking reward", everyValidatorReward)

			state.AddBalance(addr, everyValidatorReward)
			totalValidatorReward.Add(totalValidatorReward, everyValidatorReward)
		}
	}
	state.SubBalance(vm.RewardManagerPoolAddr, totalValidatorReward)
}

// allocatePackageBlock used for reward new block. it returns coinbase and error
func (rmp *rewardMgrPlugin) allocatePackageBlock(blockNumber uint64, blockHash common.Hash, coinBase common.Address, reward *big.Int, state xcom.StateDB) {

	if coinBase != vm.RewardManagerPoolAddr {

		log.Debug("allocate package reward", "blockNumber", blockNumber, "blockHash", blockHash.Hex(),
			"coinBase", coinBase.String(), "reward", reward)

		state.SubBalance(vm.RewardManagerPoolAddr, reward)
		state.AddBalance(coinBase, reward)
	}
}

//  Calculation percentage ,  input 100,10    cal:  100*10/100 = 10
func percentageCalculation(mount *big.Int, rate uint64) *big.Int {
	ratio := new(big.Int).Mul(mount, big.NewInt(int64(rate)))
	return new(big.Int).Div(ratio, common.Big100)
}

// calculateExpectReward used for calculate the stakingReward and newBlockReward that should be send in each corresponding period
func (rmp *rewardMgrPlugin) calculateExpectReward(thisYear, lastYear uint32, state xcom.StateDB) (*big.Int, *big.Int) {
	// get expected settlement epochs and new blocks per year first
	epochs := xutil.EpochsPerYear()
	blocks := xutil.CalcBlocksEachYear()
	lastYearBalance := GetYearEndBalance(state, lastYear)

	totalNewBlockReward := percentageCalculation(lastYearBalance, xcom.NewBlockRewardRate())
	totalStakingReward := new(big.Int).Sub(lastYearBalance, totalNewBlockReward)

	newBlockReward := new(big.Int).Div(totalNewBlockReward, big.NewInt(int64(blocks)))
	stakingReward := new(big.Int).Div(totalStakingReward, big.NewInt(int64(epochs)))

	log.Debug("Call calculateExpectReward", "thisYear", thisYear, "lastYear", lastYear,
		"lastYearBalance", lastYearBalance, "totalNewBlockReward", totalNewBlockReward,
		"totalStakingReward", totalStakingReward, "epochs of this year", epochs,
		"blocks of this year", blocks, "newBlockReward", newBlockReward, "stakingReward", stakingReward)

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
	log.Trace("show history cumulative issue", "lastYear", year, "amount", issue.SetBytes(bIssue))
	return issue.SetBytes(bIssue)
}

func SetYearEndBalance(state xcom.StateDB, year uint32, balance *big.Int) {
	yearEndBalanceKey := reward.HistoryBalancePrefix(year)
	state.SetState(vm.RewardManagerPoolAddr, yearEndBalanceKey, balance.Bytes())
}

func GetYearEndBalance(state xcom.StateDB, year uint32) *big.Int {
	var balance = new(big.Int)
	yearEndBalanceKey := reward.HistoryBalancePrefix(year)
	bBalance := state.GetState(vm.RewardManagerPoolAddr, yearEndBalanceKey)
	log.Trace("show balance of reward pool at last year end", "lastYear", year, "amount", balance.SetBytes(bBalance))
	return balance.SetBytes(bBalance)
}
