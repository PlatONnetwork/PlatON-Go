// Copyright 2018-2019 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package plugin

import (
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/x/reward"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type RewardMgrPlugin struct {
	currentYear    uint32
	stakingReward  *big.Int
	newBlockReward *big.Int
}

const (
	LessThanFoundationYearDeveloperRate    = 100
	AfterFoundationYearDeveloperRewardRate = 50
	AfterFoundationYearFoundRewardRate     = 50
	IncreaseIssue                          = 40
	RewardPoolIncreaseRate                 = 80 // 80% of fixed-issued tokens are allocated to reward pool each year
)

var (
	rewardOnce sync.Once
	rm         *RewardMgrPlugin = nil
)

func RewardMgrInstance() *RewardMgrPlugin {
	rewardOnce.Do(func() {
		log.Info("Init Reward plugin ...")
		rm = &RewardMgrPlugin{}
	})
	return rm
}

// BeginBlock does something like check input params before execute transactions,
// in RewardMgrPlugin it does nothing.
func (rmp *RewardMgrPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	return nil
}

// EndBlock will handle reward work, if it's time to settle, reward staking. Then reward worker
// for create new block, this is necessary. At last if current block is the last block at the end
// of year, increasing issuance.
func (rmp *RewardMgrPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	blockNumber := head.Number.Uint64()

	thisYear := xutil.CalculateYear(blockNumber)
	var lastYear uint32
	if thisYear != 0 {
		lastYear = thisYear - 1
	}

	if thisYear != rmp.currentYear {
		rmp.stakingReward, rmp.newBlockReward = rmp.calculateExpectReward(thisYear, lastYear, state)
		rmp.currentYear = thisYear
	}
	stakingReward := new(big.Int).Set(rmp.stakingReward)
	packageReward := new(big.Int).Set(rmp.newBlockReward)

	if xutil.IsEndOfEpoch(blockNumber) {
		if err := rmp.allocateStakingReward(blockNumber, blockHash, stakingReward, state); err != nil {
			return err
		}
	}

	rmp.allocatePackageBlock(blockNumber, blockHash, head.Coinbase, packageReward, state)

	// the block at the end of each year, additional issuance
	if xutil.IsYearEnd(blockNumber) {
		rmp.increaseIssuance(thisYear, lastYear, state)
	}

	return nil
}

// Confirmed does nothing
func (rmp *RewardMgrPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}

func (rmp *RewardMgrPlugin) isLessThanFoundationYear(thisYear uint32) bool {
	if thisYear < xcom.PlatONFoundationYear()-1 {
		return true
	}
	return false
}

func (rmp *RewardMgrPlugin) addPlatONFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	platonFoundationIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(xcom.PlatONFundAccount(), platonFoundationIncr)
}

func (rmp *RewardMgrPlugin) addCommunityDeveloperFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	developerFoundationIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(xcom.CDFAccount(), developerFoundationIncr)
}
func (rmp *RewardMgrPlugin) addRewardPoolIncreaseIssuance(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	rewardpoolIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(vm.RewardManagerPoolAddr, rewardpoolIncr)
}

// increaseIssuance used for increase issuance at the end of each year
func (rmp *RewardMgrPlugin) increaseIssuance(thisYear, lastYear uint32, state xcom.StateDB) {
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
	rewardpoolIncr := percentageCalculation(currIssuance, uint64(RewardPoolIncreaseRate))
	state.AddBalance(vm.RewardManagerPoolAddr, rewardpoolIncr)
	lessBalance := new(big.Int).Sub(currIssuance, rewardpoolIncr)
	if rmp.isLessThanFoundationYear(thisYear) {
		log.Debug("Call EndBlock on reward_plugin: increase issuance to developer", "thisYear", thisYear, "developBalance", lessBalance)
		rmp.addCommunityDeveloperFoundation(state, lessBalance, LessThanFoundationYearDeveloperRate)
	} else {
		log.Debug("Call EndBlock on reward_plugin: increase issuance to developer and platon", "thisYear", thisYear, "develop and platon Balance", lessBalance)
		rmp.addCommunityDeveloperFoundation(state, lessBalance, AfterFoundationYearDeveloperRewardRate)
		rmp.addPlatONFoundation(state, lessBalance, AfterFoundationYearFoundRewardRate)
	}
	balance := state.GetBalance(vm.RewardManagerPoolAddr)
	SetYearEndBalance(state, thisYear, balance)

}

// allocateStakingReward used for reward staking at the settle block
func (rmp *RewardMgrPlugin) allocateStakingReward(blockNumber uint64, blockHash common.Hash, reward *big.Int, state xcom.StateDB) error {

	log.Info("Allocate staking reward start", "blockNumber", blockNumber, "hash", blockHash,
		"epoch", xutil.CalculateEpoch(blockNumber), "reward", reward)

	verifierList, err := stk.GetVerifierList(blockHash, blockNumber, false)
	if err != nil {
		log.Error("Failed to allocateStakingReward: call GetVerifierList is failed", "blockNumber", blockNumber, "hash", blockHash, "err", err)
		return err
	}
	rmp.rewardStakingByValidatorList(state, verifierList, reward)
	return nil
}

func (rmp *RewardMgrPlugin) rewardStakingByValidatorList(state xcom.StateDB, list staking.ValidatorExQueue, reward *big.Int) {
	validatorNum := int64(len(list))
	everyValidatorReward := new(big.Int).Div(reward, big.NewInt(validatorNum))

	log.Debug("calculate validator staking reward", "validator length", validatorNum, "everyOneReward", everyValidatorReward)
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
func (rmp *RewardMgrPlugin) allocatePackageBlock(blockNumber uint64, blockHash common.Hash, coinBase common.Address, reward *big.Int, state xcom.StateDB) {

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
func (rmp *RewardMgrPlugin) calculateExpectReward(thisYear, lastYear uint32, state xcom.StateDB) (*big.Int, *big.Int) {
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
