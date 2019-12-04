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
	"math"
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

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
	db snapshotdb.DB
}

const (
	LessThanFoundationYearDeveloperRate    = 100
	AfterFoundationYearDeveloperRewardRate = 50
	AfterFoundationYearFoundRewardRate     = 50
	IncreaseIssue                          = 40
	RewardPoolIncreaseRate                 = 80 // 80% of fixed-issued tokens are allocated to reward pool each year
)

var (
	rewardOnce  sync.Once
	rm          *RewardMgrPlugin = nil
	millisecond                  = 1000
	minutes                      = 60 * millisecond
)

func RewardMgrInstance() *RewardMgrPlugin {
	rewardOnce.Do(func() {
		log.Info("Init Reward plugin ...")
		rm = &RewardMgrPlugin{
			db: snapshotdb.Instance(),
		}
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

	packageReward := new(big.Int)
	stakingReward := new(big.Int)
	var err error

	if head.Number.Uint64() == common.Big1.Uint64() {
		packageReward, stakingReward, err = rmp.CalcEpochReward(blockHash, head, state)
		if nil != err {
			log.Error("Execute CalcEpochReward fail", "blockNumber", head.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
		}
	} else {
		packageReward, err = LoadNewBlockReward(blockHash, rmp.db)
		if nil != err {
			log.Error("Load NewBlockReward fail", "blockNumber", head.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
		}
		stakingReward, err = LoadStakingReward(blockHash, rmp.db)
		if nil != err {
			log.Error("Load StakingReward fail", "blockNumber", head.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
		}
	}

	rmp.allocatePackageBlock(blockNumber, blockHash, head.Coinbase, packageReward, state)

	if xutil.IsEndOfEpoch(blockNumber) {
		if err := rmp.allocateStakingReward(blockNumber, blockHash, stakingReward, state); err != nil {
			return err
		}
		if _, _, err := rmp.CalcEpochReward(blockHash, head, state); nil != err {
			log.Error("Execute CalcEpochReward fail", "blockNumber", head.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
		}
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

func (rmp *RewardMgrPlugin) CalcEpochReward(blockHash common.Hash, head *types.Header, state xcom.StateDB) (*big.Int, *big.Int, error) {
	remainReward, err := LoadRemainingReward(blockHash, rmp.db)
	if nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	yearStartBlockNumber, yearStartTime, err := LoadYearStartTime(blockHash, rmp.db)
	if nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	incIssuanceTime, err := LoadIncIssuanceTime(blockHash, rmp.db)
	if nil != err {
		log.Error("load incIssuanceTime fail", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	if yearStartTime == 0 {
		yearStartBlockNumber = head.Number.Uint64()
		yearStartTime = head.Time.Int64()
		incIssuanceTime = yearStartTime + int64(xcom.AdditionalCycleTime()*uint64(minutes))
		if err := StorageIncIssuanceTime(blockHash, rmp.db, incIssuanceTime); nil != err {
			log.Error("storage incIssuanceTime fail", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return nil, nil, err
		}
		if err := StorageYearStartTime(blockHash, rmp.db, yearStartBlockNumber, yearStartTime); nil != err {
			log.Error("Storage year start time and block height failed", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return nil, nil, err
		}
		remainReward = GetYearEndBalance(state, 0)
		log.Info("Call CalcEpochReward, First calculation", "currBlockNumber", head.Number, "currBlockHash", blockHash, "currBlockTime", head.Time.Int64(),
			"yearStartBlockNumber", yearStartBlockNumber, "yearStartTime", yearStartTime, "incIssuanceTime", incIssuanceTime, "remainReward", remainReward)
	}
	yearNumber, err := LoadChainYearNumber(blockHash, rmp.db)
	if nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	// When the first issuance is completed
	// Each settlement cycle needs to update the year start time,
	// which is used to calculate the average annual block production rate
	epochBlocks := xutil.CalcBlocksEachEpoch()
	if yearNumber > 0 {
		if yearStartBlockNumber == 1 {
			yearStartBlockNumber += epochBlocks - 1
		} else {
			yearStartBlockNumber += epochBlocks
		}
		yearStartTime = snapshotdb.GetDBBlockChain().GetHeaderByNumber(yearStartBlockNumber).Time.Int64()
		if err := StorageYearStartTime(blockHash, rmp.db, yearStartBlockNumber, yearStartTime); nil != err {
			log.Error("Storage year start time and block height failed", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return nil, nil, err
		}
		log.Debug("Call CalcEpochReward, Adjust the sampling range of the block time", "currBlockNumber", head.Number, "currBlockHash", blockHash, "currBlockTime", head.Time.Int64(),
			"epochBlocks", epochBlocks, "yearNumber", yearNumber, "yearStartBlockNumber", yearStartBlockNumber, "yearStartTime", yearStartTime)
	}

	avgPackTime := xcom.Interval() * uint64(millisecond)
	if head.Number.Uint64() > yearStartBlockNumber {
		diffNumber := head.Number.Uint64() - yearStartBlockNumber
		diffTime := head.Time.Int64() - yearStartTime
		// If it is less than or equal to the block difference,
		// it is calculated according to the default average block production time
		if uint64(diffTime) > diffNumber {
			avgPackTime = uint64(diffTime) / diffNumber
		}
		log.Debug("Call CalcEpochReward, Calculate the average block production time in the previous year", "currBlockNumber", head.Number, "currBlockHash", blockHash,
			"currBlockTime", head.Time.Int64(), "yearStartBlockNumber", yearStartBlockNumber, "yearStartTime", yearStartTime, "diffNumber", diffNumber, "diffTime", diffTime,
			"avgPackTime", avgPackTime)
	}
	if err := StorageAvgPackTime(blockHash, rmp.db, avgPackTime); nil != err {
		log.Error("Failed to execute StorageAvgPackTime function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "avgPackTime", avgPackTime, "err", err)
		return nil, nil, err
	}

	if remainReward.Cmp(common.Big0) == 0 {
		yearNumber++
		if err := StorageChainYearNumber(blockHash, rmp.db, yearNumber); nil != err {
			log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return nil, nil, err
		}
		rmp.increaseIssuance(yearNumber, yearNumber-1, state)
		// After the increase issue is completed, update the number of rewards to be issued
		if err := StorageRemainingReward(blockHash, rmp.db, GetYearEndBalance(state, yearNumber)); nil != err {
			log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return nil, nil, err
		}
		if err := StorageIncIssuanceTime(blockHash, rmp.db, incIssuanceTime+int64(xcom.AdditionalCycleTime()*uint64(minutes))); nil != err {
			log.Error("storage incIssuanceTime fail", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return nil, nil, err
		}
		remainReward, err = LoadRemainingReward(blockHash, rmp.db)
		if nil != err {
			return nil, nil, err
		}
		incIssuanceTime, err = LoadIncIssuanceTime(blockHash, rmp.db)
		if nil != err {
			return nil, nil, err
		}
		log.Info("Call CalcEpochReward, increaseIssuance successful", "currBlockNumber", head.Number, "currBlockHash", blockHash, "currBlockTime", head.Time.Int64(),
			"yearNumber", yearNumber, "incIssuanceTime", incIssuanceTime, "yearEndBalance", GetYearEndBalance(state, yearNumber), "remainingReward", remainReward)
	}
	epochTotalReward := new(big.Int)
	// If the expected increase issue time is exceeded,
	// the increase issue time will be postponed for one settlement cycle,
	// and the remaining rewards will all be issued in the next settlement cycle
	if head.Time.Int64() >= incIssuanceTime {
		epochTotalReward.Add(epochTotalReward, remainReward)
		remainReward = new(big.Int)
		log.Info("Call CalcEpochReward, The current time has exceeded the expected additional issue time", "currBlockNumber", head.Number, "currBlockHash", blockHash,
			"currBlockTime", head.Time.Int64(), "incIssuanceTime", incIssuanceTime, "epochTotalReward", epochTotalReward)
	} else {
		remainTime := incIssuanceTime - head.Time.Int64()
		remainEpoch := 1
		remainBlocks := math.Ceil(float64(remainTime) / float64(avgPackTime))
		if remainBlocks > float64(epochBlocks) {
			remainEpoch = int(math.Ceil(remainBlocks / float64(epochBlocks)))
		}
		epochTotalReward = new(big.Int).Div(remainReward, new(big.Int).SetInt64(int64(remainEpoch)))
		// Subtract the total reward for the next cycle to calculate the remaining rewards to be issued
		remainReward = remainReward.Sub(remainReward, epochTotalReward)
		log.Debug("Call CalcEpochReward, Calculation of rewards for the next settlement cycle", "currBlockNumber", head.Number, "currBlockHash", blockHash,
			"currBlockTime", head.Time.Int64(), "incIssuanceTime", incIssuanceTime, "remainTime", remainTime, "remainBlocks", remainBlocks, "epochBlocks", epochBlocks,
			"remainEpoch", remainEpoch, "remainReward", remainReward, "epochTotalReward", epochTotalReward)
	}
	if remainReward.Cmp(common.Big0) == 0 {
		incIssuanceNumber := new(big.Int).Add(head.Number, new(big.Int).SetUint64(epochBlocks)).Uint64()
		if err := StorageIncIssuanceNumber(blockHash, rmp.db, incIssuanceNumber); nil != err {
			return nil, nil, err
		}
		log.Info("Call CalcEpochReward, IncIssuanceNumber stored successfully", "currBlockNumber", head.Number, "currBlockHash", blockHash,
			"epochBlocks", epochBlocks, "incIssuanceNumber", incIssuanceNumber)
	}
	// Get the total block reward and pledge reward for each settlement cycle
	epochTotalNewBlockReward := percentageCalculation(epochTotalReward, xcom.NewBlockRewardRate())
	epochTotalStakingReward := new(big.Int).Sub(epochTotalReward, epochTotalNewBlockReward)
	if err := StorageRemainingReward(blockHash, rmp.db, remainReward); nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	newBlockReward := new(big.Int).Div(epochTotalNewBlockReward, new(big.Int).SetInt64(int64(epochBlocks)))
	if err := StorageNewBlockReward(blockHash, rmp.db, newBlockReward); nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	if err := StorageStakingReward(blockHash, rmp.db, epochTotalStakingReward); nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	log.Debug("Call CalcEpochReward, Cycle reward", "currBlockNumber", head.Number, "currBlockHash", blockHash, "currBlockTime", head.Time.Int64(),
		"epochTotalReward", epochTotalReward, "newBlockRewardRate", xcom.NewBlockRewardRate(), "epochTotalNewBlockReward", epochTotalNewBlockReward,
		"epochTotalStakingReward", epochTotalStakingReward, "epochBlocks", epochBlocks, "newBlockReward", newBlockReward)
	return newBlockReward, epochTotalStakingReward, nil
}

func StorageYearStartTime(hash common.Hash, snapshotDB snapshotdb.DB, blockNumber uint64, yearStartTime int64) error {
	if blockNumber > 0 {
		if err := snapshotDB.Put(hash, reward.YearStartBlockNumberKey, common.Uint64ToBytes(blockNumber)); nil != err {
			log.Error("Failed to execute StorageYearStartEndTime function", "hash", hash.TerminalString(), "key", string(reward.YearStartBlockNumberKey),
				"value", blockNumber, "err", err)
			return err
		}
	}
	if yearStartTime > 0 {
		if err := snapshotDB.Put(hash, reward.YearStartTimeKey, common.Int64ToBytes(yearStartTime)); nil != err {
			log.Error("Failed to execute StorageYearStartEndTime function", "hash", hash.TerminalString(), "key", string(reward.YearStartTimeKey),
				"value", yearStartTime, "err", err)
			return err
		}
	}
	return nil
}

func LoadYearStartTime(hash common.Hash, snapshotDB snapshotdb.DB) (yearStartBlockNumber uint64, yearStartTime int64, error error) {
	yearStartBlockNumberByte, err := snapshotDB.Get(hash, reward.YearStartBlockNumberKey)
	if nil != err {
		if err != snapshotdb.ErrNotFound {
			log.Error("Failed to execute LoadYearStartEndTime function", "hash", hash.TerminalString(), "key", string(reward.YearStartBlockNumberKey), "err", err)
			return 0, 0, err
		} else {
			yearStartBlockNumber = 0
		}
	} else {
		yearStartBlockNumber = common.BytesToUint64(yearStartBlockNumberByte)
	}
	yearStartTimeByte, err := snapshotDB.Get(hash, reward.YearStartTimeKey)
	if nil != err {
		if err != snapshotdb.ErrNotFound {
			log.Error("Failed to execute LoadYearStartEndTime function", "hash", hash.TerminalString(), "key", string(reward.YearStartTimeKey), "err", err)
			return 0, 0, err
		} else {
			yearStartTime = 0
		}
	} else {
		yearStartTime = common.BytesToInt64(yearStartTimeByte)
	}
	return
}

func StorageIncIssuanceTime(hash common.Hash, snapshotDB snapshotdb.DB, incTime int64) error {
	if err := snapshotDB.Put(hash, reward.IncIssuanceTimeKey, common.Int64ToBytes(incTime)); nil != err {
		log.Error("Failed to execute StorageIncIssuanceTime function", "hash", hash.TerminalString(), "key", string(reward.IncIssuanceTimeKey),
			"value", incTime, "err", err)
		return err
	}
	return nil
}

func LoadIncIssuanceTime(hash common.Hash, snapshotDB snapshotdb.DB) (int64, error) {
	incTimeByte, err := snapshotDB.Get(hash, reward.IncIssuanceTimeKey)
	if nil != err {
		if err != snapshotdb.ErrNotFound {
			log.Error("Failed to execute LoadIncIssuanceTime function", "hash", hash.TerminalString(), "key", string(reward.IncIssuanceTimeKey), "err", err)
			return 0, err
		} else {
			return 0, nil
		}
	}
	return common.BytesToInt64(incTimeByte), nil
}

func StorageRemainingReward(hash common.Hash, snapshotDB snapshotdb.DB, remainReward *big.Int) error {
	if err := snapshotDB.Put(hash, reward.RemainingRewardKey, remainReward.Bytes()); nil != err {
		log.Error("Failed to execute StorageRemainingReward function", "hash", hash.TerminalString(), "remainReward", remainReward, "err", err)
		return err
	}
	return nil
}

func LoadRemainingReward(hash common.Hash, snapshotDB snapshotdb.DB) (*big.Int, error) {
	remainRewardByte, err := snapshotDB.Get(hash, reward.RemainingRewardKey)
	if nil != err {
		if err != snapshotdb.ErrNotFound {
			log.Error("Failed to execute LoadRemainingReward function", "hash", hash.TerminalString(), "key", string(reward.RemainingRewardKey), "err", err)
			return nil, err
		}
	}
	return new(big.Int).SetBytes(remainRewardByte), nil
}

func StorageNewBlockReward(hash common.Hash, snapshotDB snapshotdb.DB, newBlockReward *big.Int) error {
	if err := snapshotDB.Put(hash, reward.NewBlockRewardKey, newBlockReward.Bytes()); nil != err {
		log.Error("Failed to execute StorageNewBlockReward function", "hash", hash.TerminalString(), "newBlockReward", newBlockReward, "err", err)
		return err
	}
	return nil
}

func LoadNewBlockReward(hash common.Hash, snapshotDB snapshotdb.DB) (*big.Int, error) {
	newBlockRewardByte, err := snapshotDB.Get(hash, reward.NewBlockRewardKey)
	if nil != err {
		log.Error("Failed to execute LoadRemainingReward function", "hash", hash.TerminalString(), "key", string(reward.NewBlockRewardKey), "err", err)
		return nil, err
	}
	return new(big.Int).SetBytes(newBlockRewardByte), nil
}

func StorageStakingReward(hash common.Hash, snapshotDB snapshotdb.DB, stakingReward *big.Int) error {
	if err := snapshotDB.Put(hash, reward.StakingRewardKey, stakingReward.Bytes()); nil != err {
		log.Error("Failed to execute StorageStakingReward function", "hash", hash.TerminalString(), "stakingReward", stakingReward, "err", err)
		return err
	}
	return nil
}

func LoadStakingReward(hash common.Hash, snapshotDB snapshotdb.DB) (*big.Int, error) {
	stakingRewardByte, err := snapshotDB.Get(hash, reward.StakingRewardKey)
	if nil != err {
		log.Error("Failed to execute LoadStakingReward function", "hash", hash.TerminalString(), "key", string(reward.StakingRewardKey), "err", err)
		return nil, err
	}
	return new(big.Int).SetBytes(stakingRewardByte), nil
}

func StorageChainYearNumber(hash common.Hash, snapshotDB snapshotdb.DB, yearNumber uint32) error {
	if err := snapshotDB.Put(hash, reward.ChainYearNumberKey, common.Uint32ToBytes(yearNumber)); nil != err {
		log.Error("Failed to execute StorageChainYearNumber function", "hash", hash.TerminalString(), "yearNumber", yearNumber, "err", err)
		return err
	}
	return nil
}

func LoadChainYearNumber(hash common.Hash, snapshotDB snapshotdb.DB) (uint32, error) {
	chainYearNumberByte, err := snapshotDB.Get(hash, reward.ChainYearNumberKey)
	if nil != err {
		if err == snapshotdb.ErrNotFound {
			log.Info("Data obtained for the first year", "hash", hash.TerminalString())
			return 0, nil
		}
		log.Error("Failed to execute LoadChainYearNumber function", "hash", hash.TerminalString(), "key", string(reward.ChainYearNumberKey), "err", err)
		return 0, err
	}
	return common.BytesToUint32(chainYearNumberByte), nil
}

func StorageIncIssuanceNumber(hash common.Hash, snapshotDB snapshotdb.DB, incIssuanceNumber uint64) error {
	if err := snapshotDB.Put(hash, reward.IncIssuanceNumberKey, common.Uint64ToBytes(incIssuanceNumber)); nil != err {
		log.Error("Failed to execute StorageIncIssuanceNumber function", "hash", hash.TerminalString(), "incIssuanceNumber", incIssuanceNumber, "err", err)
		return err
	}
	return nil
}

func LoadIncIssuanceNumber(hash common.Hash, snapshotDB snapshotdb.DB) (uint64, error) {
	incIssuanceNumberByte, err := snapshotDB.Get(hash, reward.IncIssuanceNumberKey)
	if nil != err {
		if err == snapshotdb.ErrNotFound {
			return 0, nil
		}
		log.Error("Failed to execute LoadIncIssuanceNumber function", "hash", hash.TerminalString(), "key", string(reward.IncIssuanceNumberKey), "err", err)
		return 0, err
	}
	return common.BytesToUint64(incIssuanceNumberByte), nil
}

func IsYearEnd(hash common.Hash, blockNumber uint64, db snapshotdb.DB) (bool, error) {
	number, err := LoadIncIssuanceNumber(hash, db)
	if nil != err {
		return false, err
	}
	if number == blockNumber {
		return true, nil
	}
	return false, nil
}

func StorageAvgPackTime(hash common.Hash, snapshotDB snapshotdb.DB, avgPackTime uint64) error {
	if err := snapshotDB.Put(hash, reward.AvgPackTimeKey, common.Uint64ToBytes(avgPackTime)); nil != err {
		log.Error("Failed to execute StorageAvgPackTime function", "hash", hash.TerminalString(), "avgPackTime", avgPackTime, "err", err)
		return err
	}
	return nil
}

func LoadAvgPackTime(hash common.Hash, snapshotDB snapshotdb.DB) (uint64, error) {
	avgPackTimeByte, err := snapshotDB.Get(hash, reward.AvgPackTimeKey)
	if nil != err {
		if err == snapshotdb.ErrNotFound {
			return 0, nil
		}
		log.Error("Failed to execute LoadAvgPackTime function", "hash", hash.TerminalString(), "key", string(reward.AvgPackTimeKey), "err", err)
		return 0, err
	}
	return common.BytesToUint64(avgPackTimeByte), nil
}
