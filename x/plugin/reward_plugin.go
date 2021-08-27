// Copyright 2018-2020 The PlatON Network Authors
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
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"sort"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/rlp"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

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
	db            snapshotdb.DB
	nodeID        discover.NodeID
	nodeADD       common.NodeAddress
	stakingPlugin *StakingPlugin
}

const (
	LessThanFoundationYearDeveloperRate    = 100
	AfterFoundationYearDeveloperRewardRate = 50
	AfterFoundationYearFoundRewardRate     = 50
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
		sdb := snapshotdb.Instance()
		rm = &RewardMgrPlugin{
			db:            sdb,
			stakingPlugin: StakingInstance(),
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
func (rmp *RewardMgrPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB, downloading Downloading) error {
	blockNumber := head.Number.Uint64()

	// 待分配的出块奖励金额，每个结算周期可能不一样
	packageReward := new(big.Int)
	stakingReward := new(big.Int)
	var err error

	if head.Number.Uint64() == common.Big1.Uint64() {
		//第一个块，也就是第一个EPOCH，所以首先要计算第一个EPOCH的出块奖励、质押奖励
		packageReward, stakingReward, err = rmp.CalcEpochReward(blockHash, head, state)
		if nil != err {
			log.Error("Execute CalcEpochReward fail", "blockNumber", head.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
			return err
		}
		log.Info("Block 1 reward", "packageReward", packageReward.Uint64())
	} else {
		packageReward, err = LoadNewBlockReward(blockHash, rmp.db)
		if nil != err {
			log.Error("Load NewBlockReward fail", "blockNumber", head.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
			return err
		}
		stakingReward, err = LoadStakingReward(blockHash, rmp.db)
		if nil != err {
			log.Error("Load StakingReward fail", "blockNumber", head.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
			return err
		}
	}
	//委托用户奖励，包括节点（质押奖励+出块奖励）*委托分红比例
	//分配出块奖励
	if err := rmp.AllocatePackageBlock(blockHash, head, packageReward, state); err != nil {
		return err
	}

	if xutil.IsEndOfEpoch(blockNumber) {
		//结算周期末，分配质押奖励。
		//质押节点能直接收到质押奖励，委托用户的质押奖励，会从激励池转入委托奖励合约。
		verifierList, err := rmp.AllocateStakingReward(blockNumber, blockHash, stakingReward, state)
		if err != nil {
			return err
		}
		// 保存质押节点，给委托用户的奖励信息。
		if err := rmp.HandleDelegatePerReward(blockHash, blockNumber, verifierList, state); err != nil {
			return err
		}

		if err := rmp.runIncreaseIssuance(blockHash, head, state); nil != err {
			return err
		}
		if _, _, err := rmp.CalcEpochReward(blockHash, head, state); nil != err {
			log.Error("Execute CalcEpochReward fail", "blockNumber", head.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
			return err
		}
	}
	return nil
}

//stats
func convertVerifier(verifierList []*staking.Candidate) []*common.CandidateInfo {
	candidateInfoList := make([]*common.CandidateInfo, len(verifierList))
	for idx, verifier := range verifierList {
		candidateInfo := &common.CandidateInfo{
			NodeID: common.NodeID(verifier.NodeId), MinerAddress: verifier.BenefitAddress,
		}
		candidateInfoList[idx] = candidateInfo
	}
	return candidateInfoList
}

// Confirmed does nothing
func (rmp *RewardMgrPlugin) Confirmed(nodeId discover.NodeID, block *types.Block) error {
	return nil
}

func (rmp *RewardMgrPlugin) SetCurrentNodeID(nodeId discover.NodeID) {
	rmp.nodeID = nodeId
	add, err := xutil.NodeId2Addr(rmp.nodeID)
	if err != nil {
		panic(err)
	}
	rmp.nodeADD = add
}

//platonFoundationYear这个配置值，表示从这次增发开始，需要分配一部分给PlatONFundation
//创世块已经发现了1次；所以第一年增发时(year=1)，实际上时第二次增发了；所以判断条件需要 - 1
func (rmp *RewardMgrPlugin) isLessThanFoundationYear(thisYear uint32) bool {
	if thisYear < xcom.PlatONFoundationYear()-1 {
		return true
	}
	return false
}

//stats
//func (rmp *RewardMgrPlugin) addPlatONFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
func (rmp *RewardMgrPlugin) addPlatONFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) (common.Address, *big.Int) {
	platonFoundationIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(xcom.PlatONFundAccount(), platonFoundationIncr)
	return xcom.PlatONFundAccount(), platonFoundationIncr
}

//stats
//func (rmp *RewardMgrPlugin) addCommunityDeveloperFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
func (rmp *RewardMgrPlugin) addCommunityDeveloperFoundation(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) (common.Address, *big.Int) {
	developerFoundationIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(xcom.CDFAccount(), developerFoundationIncr)
	return xcom.CDFAccount(), developerFoundationIncr
}
func (rmp *RewardMgrPlugin) addRewardPoolIncreaseIssuance(state xcom.StateDB, currIssuance *big.Int, allocateRate uint32) {
	rewardpoolIncr := percentageCalculation(currIssuance, uint64(allocateRate))
	state.AddBalance(vm.RewardManagerPoolAddr, rewardpoolIncr)
}

// increaseIssuance used for increase issuance at the end of each year
func (rmp *RewardMgrPlugin) increaseIssuance(thisYear, lastYear uint32, state xcom.StateDB, blockNumber uint64, blockHash common.Hash) error {
	var currIssuance *big.Int

	//stats: 收集增发数据
	additionalIssuance := new(common.AdditionalIssuanceData)
	//issuance increase
	{
		//todo: ppos_config.issue_total，目前累计增发多少（不包含本次增发）
		histIssuance := GetHistoryCumulativeIssue(state, lastYear)

		//todo: ppos_config.issue_ratio，读取发行比例（缺省是genesis.json中定义，参数提案通过后，修改相应参数值）
		increaseIssuanceRatio, err := gov.GovernIncreaseIssuanceRatio(blockNumber, blockHash)
		if nil != err {
			log.Error("Failed to increaseIssuance, call GovernIncreaseIssuanceRatio is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"histIssuance", histIssuance, "err", err)
			return err
		}

		//计算此次发行金额 = issue_total * increaseIssuanceRatio / 10000
		tmp := new(big.Int).Mul(histIssuance, big.NewInt(int64(increaseIssuanceRatio)))
		currIssuance = tmp.Div(tmp, big.NewInt(10000))

		// Restore the cumulative issue at this year end
		/*histIssuance.Add(histIssuance, currIssuance)
		SetYearEndCumulativeIssue(state, thisYear, histIssuance)
		log.Debug("Call EndBlock on reward_plugin: increase issuance", "thisYear", thisYear, "addIssuance", currIssuance, "hit", histIssuance)
		*/
		//stats
		//计算总发行金额
		newTotalIssuance := new(big.Int).Add(histIssuance, currIssuance)
		//todo: chain_env.issue_amount, chain_env.total_issue_amount，更新发行金额，总发行金额(整个操作在最后做)。
		//todo:增加一个表issue_history表，block/hash/time/issue/issue_total/reward_pool/开发者基金增加值/Platon基金增加值
		//todo:chain_env.reward_pool_available, chain_env.reward_pool_next_year_available，
		//todo: 重新设置 reward_pool_available = reward_pool_available + reward_pool_next_year_available + currIssuance, reward_pool_next_year_available=0
		SetYearEndCumulativeIssue(state, thisYear, newTotalIssuance)
		log.Debug("Call EndBlock on reward_plugin: increase issuance", "thisYear", thisYear, "origTotalIssuance", histIssuance, "increment", currIssuance, "newTotalIssuance", newTotalIssuance)

		//stats: 收集增发数据
		additionalIssuance.AdditionalNo = thisYear
		additionalIssuance.AdditionalBase = histIssuance          //上年发行量
		additionalIssuance.AdditionalAmount = currIssuance        //今年增发量 = 上年发行量 * 今年增发率
		additionalIssuance.AdditionalRate = increaseIssuanceRatio //今年增发率
	}
	//今年的增发量，需要转入一部分到激励池中，以及其它基金会账户
	rewardpoolIncr := percentageCalculation(currIssuance, uint64(RewardPoolIncreaseRate))
	state.AddBalance(vm.RewardManagerPoolAddr, rewardpoolIncr)

	//stats: 收集增发数据
	additionalIssuance.AddIssuanceItem(vm.RewardManagerPoolAddr, rewardpoolIncr)

	//todo: 增发剩余金额
	lessBalance := new(big.Int).Sub(currIssuance, rewardpoolIncr)
	if rmp.isLessThanFoundationYear(thisYear) {
		//如果增发次数小于配置值（注意，配置值包含了创世块中已经分配的一次），则增发剩余金额都转入社区开发者金额
		log.Debug("Call EndBlock on reward_plugin: increase issuance to developer", "thisYear", thisYear, "developBalance", lessBalance)
		address, amount := rmp.addCommunityDeveloperFoundation(state, lessBalance, LessThanFoundationYearDeveloperRate)
		//stats: 收集增发数据
		additionalIssuance.AddIssuanceItem(address, amount)
	} else {
		//否则，增发剩余金额的50%，分配给社区开发者基金，50%，分配给PlatON基金会
		log.Debug("Call EndBlock on reward_plugin: increase issuance to developer and platon", "thisYear", thisYear, "develop and platon Balance", lessBalance)

		address, amount := rmp.addCommunityDeveloperFoundation(state, lessBalance, AfterFoundationYearDeveloperRewardRate)
		//stats: 收集增发数据
		additionalIssuance.AddIssuanceItem(address, amount)

		address, amount = rmp.addPlatONFoundation(state, lessBalance, AfterFoundationYearFoundRewardRate)
		//stats: 收集增发数据
		additionalIssuance.AddIssuanceItem(address, amount)
	}

	common.CollectAdditionalIssuance(blockNumber, additionalIssuance)

	balance := state.GetBalance(vm.RewardManagerPoolAddr)
	SetYearEndBalance(state, thisYear, balance)
	return nil
}

// AllocateStakingReward used for reward staking at the settle block
func (rmp *RewardMgrPlugin) AllocateStakingReward(blockNumber uint64, blockHash common.Hash, sreward *big.Int, state xcom.StateDB) ([]*staking.Candidate, error) {

	log.Info("Allocate staking reward start", "blockNumber", blockNumber, "hash", blockHash,
		"epoch", xutil.CalculateEpoch(blockNumber), "reward", sreward)
	verifierList, err := rmp.stakingPlugin.GetVerifierCandidateInfo(blockHash, blockNumber)
	if err != nil {
		log.Error("Failed to AllocateStakingReward: call GetVerifierList is failed", "blockNumber", blockNumber, "hash", blockHash, "err", err)
		return nil, err
	}

	//把这个周期每个质押节点的质押奖励，进行分配
	if err := rmp.rewardStakingByValidatorList(state, verifierList, sreward); err != nil {
		log.Error("reward staking by validator list fail", "err", err, "bn", blockNumber, "bh", blockHash)
		return nil, err
	}

	//stats: 收集待分配的质押奖励金额，每个结算周期可能不一样
	common.CollectStakingRewardData(blockNumber, sreward, convertVerifier(verifierList))
	return verifierList, nil
}

func (rmp *RewardMgrPlugin) ReturnDelegateReward(address common.Address, amount *big.Int, state xcom.StateDB) error {
	if amount.Cmp(common.Big0) > 0 {

		DelegateRewardPool := state.GetBalance(vm.DelegateRewardPoolAddr)

		if DelegateRewardPool.Cmp(amount) < 0 {
			return fmt.Errorf("DelegateRewardPool balance is not enougth,want %v have %v", amount, DelegateRewardPool)
		}

		state.SubBalance(vm.DelegateRewardPoolAddr, amount)
		state.AddBalance(address, amount)
	}
	return nil
}

//  保存质押节点，给委托用户的奖励的信息。
func (rmp *RewardMgrPlugin) HandleDelegatePerReward(blockHash common.Hash, blockNumber uint64, list []*staking.Candidate, state xcom.StateDB) error {
	currentEpoch := xutil.CalculateEpoch(blockNumber)
	for _, verifier := range list {
		if verifier.CurrentEpochDelegateReward.Cmp(common.Big0) == 0 {
			continue
		}
		if verifier.DelegateTotal.Cmp(common.Big0) == 0 {
			log.Debug("handleDelegatePerReward return delegateReward", "epoch", currentEpoch, "reward", verifier.CurrentEpochDelegateReward, "add", verifier.BenefitAddress)
			if err := rmp.ReturnDelegateReward(verifier.BenefitAddress, verifier.CurrentEpochDelegateReward, state); err != nil {
				log.Error("HandleDelegatePerReward ReturnDelegateReward fail", "err", err, "blockNumber", blockNumber)
			}
		} else {
			//质押节点给有效委托的奖励信息。（总的，没有算每个有效委托的奖励）
			per := reward.NewDelegateRewardPer(currentEpoch, verifier.CurrentEpochDelegateReward, verifier.DelegateTotal)
			//把奖励信息保存起来。
			//把每个节点，按key=节点ID+质押块高，来保存应该分配的委托奖励信息
			if err := AppendDelegateRewardPer(blockHash, verifier.NodeId, verifier.StakingBlockNum, per, rmp.db); err != nil {
				log.Error("call handleDelegatePerReward fail AppendDelegateRewardPer", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
					"nodeId", verifier.NodeId.TerminalString(), "err", err, "CurrentEpochDelegateReward", verifier.CurrentEpochDelegateReward, "delegateTotal", verifier.DelegateTotal)
				return err
			}
			currentEpochDelegateReward := new(big.Int).Set(verifier.CurrentEpochDelegateReward)

			//查看是否有新的委托分红比例生效。
			//清除节点的当前结算周期的委托分红金额；如果有新的委托分红比例生效，就切换到新的委托分红比例。
			changed := verifier.PrepareNextEpoch()
			canAddr, err := xutil.NodeId2Addr(verifier.NodeId)
			if nil != err {
				log.Error("Failed to handleDelegatePerReward on rewardMgrPlugin: nodeId parse addr failed",
					"blockNumber", blockNumber, "blockHash", blockHash, "nodeID", verifier.NodeId.String(), "err", err)
				return err
			}

			if changed {
				//stats
				common.CollectCandidateChanged(blockNumber, canAddr)
			}

			//为下个结算周期保存节点新的新信息（累计委托分红，新周期累计分红，新的分红比例）
			//todo:lvxiaoyi，这个逻辑放到PrepareNextEpoch()中，作为一个整体逻辑
			if err := rmp.stakingPlugin.db.SetCanMutableStore(blockHash, canAddr, verifier.CandidateMutable); err != nil {
				log.Error("Failed to handleDelegatePerReward on rewardMgrPlugin: setCanMutableStore  failed",
					"blockNumber", blockNumber, "blockHash", blockHash, "err", err, "mutable", verifier.CandidateMutable)
				return err
			}
			log.Debug("handleDelegatePerReward add newDelegateRewardPer", "blockNum", blockNumber, "node_id", verifier.NodeId.TerminalString(), "stakingNum", verifier.StakingBlockNum,
				"cu_epoch_delegate_reward", currentEpochDelegateReward, "total_delegate_reward", verifier.DelegateRewardTotal, "total_delegate", verifier.DelegateTotal,
				"epoch", currentEpoch)
		}
	}
	return nil
}

func (rmp *RewardMgrPlugin) WithdrawDelegateReward(blockHash common.Hash, blockNum uint64, account common.Address, list []*DelegationInfoWithRewardPerList, state xcom.StateDB) ([]reward.NodeDelegateReward, error) {
	log.Debug("Call withdraw delegate reward: begin", "account", account, "list", list, "blockNum", blockNum, "blockHash", blockHash, "epoch", xutil.CalculateEpoch(blockNum))

	rewards := make([]reward.NodeDelegateReward, 0)
	if len(list) == 0 {
		return rewards, nil
	}
	currentEpoch := xutil.CalculateEpoch(blockNum)
	receiveReward := new(big.Int)
	for _, delWithPer := range list {
		rewardsReceive := calcDelegateIncome(currentEpoch, delWithPer.DelegationInfo.Delegation, delWithPer.RewardPerList)
		rewards = append(rewards, reward.NodeDelegateReward{
			NodeID:     delWithPer.DelegationInfo.NodeID,
			StakingNum: delWithPer.DelegationInfo.StakeBlockNumber,
			Reward:     new(big.Int).Set(delWithPer.DelegationInfo.Delegation.CumulativeIncome),
		})
		if len(rewardsReceive) != 0 {
			if err := UpdateDelegateRewardPer(blockHash, delWithPer.DelegationInfo.NodeID, delWithPer.DelegationInfo.StakeBlockNumber, rewardsReceive, rmp.db); err != nil {
				log.Error("call WithdrawDelegateReward UpdateDelegateRewardPer fail", "err", err)
				return nil, err
			}
		}
		// Execute new logic after this version.
		// Update the delegation information only when there is delegation income available.
		if delWithPer.DelegationInfo.Delegation.CumulativeIncome.Cmp(common.Big0) > 0 {
			receiveReward.Add(receiveReward, delWithPer.DelegationInfo.Delegation.CumulativeIncome)
			delWithPer.DelegationInfo.Delegation.CleanCumulativeIncome(uint32(currentEpoch))
			if err := rmp.stakingPlugin.db.SetDelegateStore(blockHash, account, delWithPer.DelegationInfo.NodeID, delWithPer.DelegationInfo.StakeBlockNumber, delWithPer.DelegationInfo.Delegation); err != nil {
				return nil, err
			}
		}

		log.Debug("WithdrawDelegateReward rewardsReceive", "rewardsReceive", rewardsReceive, "blockNum", blockNum)
	}
	if receiveReward.Cmp(common.Big0) > 0 {
		if err := rmp.ReturnDelegateReward(account, receiveReward, state); err != nil {
			log.Error("Call withdraw delegate reward ReturnDelegateReward fail", "err", err, "blockNum", blockNum)
			return nil, common.InternalError
		}
	}
	log.Debug("Call withdraw delegate reward: end", "account", account, "rewards", rewards, "blockNum", blockNum, "blockHash", blockHash, "receiveReward", receiveReward)

	return rewards, nil
}

func (rmp *RewardMgrPlugin) GetDelegateReward(blockHash common.Hash, blockNum uint64, account common.Address, nodes []discover.NodeID, state xcom.StateDB) ([]reward.NodeDelegateRewardPresenter, error) {
	log.Debug("Call RewardMgrPlugin: query delegate reward result begin", "account", account, "nodes", nodes, "num", blockNum)

	dls, err := rmp.stakingPlugin.db.GetDelegatesInfo(blockHash, account)
	if err != nil {
		log.Error("Call GetDelegateReward GetDelegatesInfo fail", "err", err, "account", account)
		return nil, err
	}
	if len(dls) == 0 {
		return nil, reward.ErrDelegationNotFound
	}
	if len(nodes) > 0 {
		nodeMap := make(map[discover.NodeID]struct{})
		for _, node := range nodes {
			nodeMap[node] = struct{}{}
		}

		for i := 0; i < len(dls); {
			if _, ok := nodeMap[dls[i].NodeID]; !ok {
				dls = append(dls[:i], dls[i+1:]...)
			} else {
				i++
			}
		}
		if len(dls) == 0 {
			return nil, reward.ErrDelegationNotFound
		}
	} else {
		if len(dls) > int(xcom.TheNumberOfDelegationsReward()) {
			sort.Sort(staking.DelByDelegateEpoch(dls))
		}
	}

	currentEpoch := xutil.CalculateEpoch(blockNum)
	delegationInfoWithRewardPerList := make([]*DelegationInfoWithRewardPerList, 0)
	for _, stakingNode := range dls {
		delegateRewardPerList, err := rmp.GetDelegateRewardPerList(blockHash, stakingNode.NodeID, stakingNode.StakeBlockNumber, uint64(stakingNode.Delegation.DelegateEpoch), currentEpoch-1)
		if err != nil {
			log.Error("Call GetDelegateReward GetDelegateRewardPerList fail", "err", err, "account", account)
			return nil, err
		}
		delegationInfoWithRewardPerList = append(delegationInfoWithRewardPerList, NewDelegationInfoWithRewardPerList(stakingNode, delegateRewardPerList))
	}
	rewards := make([]reward.NodeDelegateRewardPresenter, 0)

	for _, delWithPer := range delegationInfoWithRewardPerList {
		calcDelegateIncome(currentEpoch, delWithPer.DelegationInfo.Delegation, delWithPer.RewardPerList)

		rewards = append(rewards, reward.NodeDelegateRewardPresenter{
			NodeID:     delWithPer.DelegationInfo.NodeID,
			StakingNum: delWithPer.DelegationInfo.StakeBlockNumber,
			Reward:     (*hexutil.Big)(new(big.Int).Set(delWithPer.DelegationInfo.Delegation.CumulativeIncome)),
		})

	}
	log.Debug("Call RewardMgrPlugin: query delegate reward result end", "num", blockNum, "account", account, "nodes", nodes, "rewards", rewards, "perList", delegationInfoWithRewardPerList)

	return rewards, nil
}

func (rmp *RewardMgrPlugin) CalDelegateRewardAndNodeReward(totalReward *big.Int, per uint16) (*big.Int, *big.Int) {
	tmp := new(big.Int).Mul(totalReward, big.NewInt(int64(per)))
	tmp.Div(tmp, big.NewInt(10000))
	return tmp, new(big.Int).Sub(totalReward, tmp)
}

// 把这个周期每个质押节点的质押奖励，进行分配。
// 1. 质押节点，直接从激励池拿到质押奖励。
// 2. 委托用户，质押奖励从激励池发放到委托激励合约。
// 3. 为每个质押节点，记录应该分配给委托用户的所有奖励。
func (rmp *RewardMgrPlugin) rewardStakingByValidatorList(state xcom.StateDB, list []*staking.Candidate, reward *big.Int) error {
	validatorNum := int64(len(list))
	//每个结算周期的质押奖励时一定的，给所有质押节点来分。这里求每个质押节点能分到的质押奖励。
	everyValidatorReward := new(big.Int).Div(reward, big.NewInt(validatorNum))

	log.Debug("calculate validator staking reward", "validator length", validatorNum, "everyOneReward", everyValidatorReward)
	totalValidatorReward, totalValidatorDelegateReward := new(big.Int), new(big.Int)

	for _, value := range list {
		delegateReward, stakingReward := new(big.Int), new(big.Int).Set(everyValidatorReward)
		if value.ShouldGiveDelegateReward() {
			// 计算质押奖励，在质押节点，和委托用户之间的分配。
			delegateReward, stakingReward = rmp.CalDelegateRewardAndNodeReward(everyValidatorReward, value.RewardPer)
			totalValidatorDelegateReward.Add(totalValidatorDelegateReward, delegateReward)
			log.Debug("allocate delegate reward of staking one-by-one", "nodeId", value.NodeId.TerminalString(), "staking reward", stakingReward, "per", value.RewardPer, "delegateReward", delegateReward)
			//the  CurrentEpochDelegateReward will use by cal delegate reward Per
			//把当前要分配给委托用户的质押奖励，累计到需要分配给委托用户的总奖励中
			value.CurrentEpochDelegateReward.Add(value.CurrentEpochDelegateReward, delegateReward)
		}
		if value.BenefitAddress != vm.RewardManagerPoolAddr {
			log.Debug("allocate staking reward one-by-one", "nodeId", value.NodeId.String(),
				"benefitAddress", value.BenefitAddress.String(), "staking reward", stakingReward)

			//给节质押点发放质押奖励
			state.AddBalance(value.BenefitAddress, stakingReward)
			//记录总发放的质押奖励，后面需要从激励池中扣除
			totalValidatorReward.Add(totalValidatorReward, stakingReward)
		}
	}
	// 把这个结算周期，分配给所有用户的质押奖励，都转入委托奖励合约中。
	state.AddBalance(vm.DelegateRewardPoolAddr, totalValidatorDelegateReward)
	//从激励池中扣除已经发放给所有质押节点的质押奖励。
	state.SubBalance(vm.RewardManagerPoolAddr, new(big.Int).Add(totalValidatorDelegateReward, totalValidatorReward))
	return nil
}

func (rmp *RewardMgrPlugin) getBlockMinderAddress(blockHash common.Hash, head *types.Header) (discover.NodeID, common.NodeAddress, error) {
	if blockHash == common.ZeroHash {
		return rmp.nodeID, rmp.nodeADD, nil
	}
	pk := head.CachePublicKey()
	if pk == nil {
		return discover.ZeroNodeID, common.ZeroNodeAddr, errors.New("failed to get the public key of the block producer")
	}
	return discover.PubkeyID(pk), crypto.PubkeyToNodeAddress(*pk), nil
}

// AllocatePackageBlock used for reward new block. it returns coinbase and error
// 每出一个块，马上分配出块奖励
// 1. 出块节点，直接从激励池拿到质押奖励。
// 2. 委托用户，出块奖励从激励池发放到委托激励合约。
// 3. 为每个质押节点，记录应该分配给委托用户的所有奖励。
func (rmp *RewardMgrPlugin) AllocatePackageBlock(blockHash common.Hash, head *types.Header, reward *big.Int, state xcom.StateDB) error {

	//由header中pubkey得出块节点id
	nodeID, add, err := rmp.getBlockMinderAddress(blockHash, head)
	if err != nil {
		log.Error("AllocatePackageBlock getBlockMinderAddress fail", "err", err, "blockNumber", head.Number, "blockHash", blockHash)
		return err
	}

	log.Debug("Alloc block reward", "blockNumber", head.Number.Uint64(), "blockHash", blockHash, "nodeID", nodeID.String(), "nodeAddr", add.String(), "coinBase", head.Coinbase.Bech32())

	currVerifier, err := rmp.stakingPlugin.IsCurrVerifier(blockHash, head.Number.Uint64(), nodeID, false)
	if err != nil {
		log.Error("AllocatePackageBlock IsCurrVerifier fail", "err", err, "blockNumber", head.Number, "blockHash", blockHash)
		return err
	}
	//stats,当前结算周期的最后一个选举块上，确定了下一个结算周期开始的第一个共识轮的23个出块节点；
	//但是到了当前结算周期末，才确定下一个结算周期的101人；那么下一个结算周期第一轮出块时，出块节点可能不再在下一个结算周期的101里了。此时，它出的块，出块奖励不分红，都是它的出块奖励。
	//tddo:跟踪系统需要知道coinBase/minerAddress和nodeId的对应关系
	blockReward := big.NewInt(0).Set(reward)
	if currVerifier {
		cm, err := rmp.stakingPlugin.GetCanMutable(blockHash, add)
		if err != nil {
			log.Error("AllocatePackageBlock GetCanMutable fail", "err", err, "blockNumber", head.Number, "blockHash", blockHash, "add", add)
			return err
		}
		if cm.ShouldGiveDelegateReward() {
			delegateReward := new(big.Int).SetUint64(0)
			delegateReward, reward = rmp.CalDelegateRewardAndNodeReward(reward, cm.RewardPer)
			//2. 委托用户，出块奖励从激励池发放到委托激励合约。
			state.SubBalance(vm.RewardManagerPoolAddr, delegateReward)
			state.AddBalance(vm.DelegateRewardPoolAddr, delegateReward)

			//为每个质押节点，记录应该分配给委托用户的所有奖励。
			cm.CurrentEpochDelegateReward.Add(cm.CurrentEpochDelegateReward, delegateReward)
			log.Debug("allocate package reward, delegate reward", "blockNumber", head.Number, "blockHash", blockHash, "delegateReward", delegateReward, "epochDelegateReward", cm.CurrentEpochDelegateReward)

			if err := rmp.stakingPlugin.db.SetCanMutableStore(blockHash, add, cm); err != nil {
				log.Error("AllocatePackageBlock SetCanMutableStore fail", "err", err, "blockNumber", head.Number, "blockHash", blockHash)
				return err
			}
		}

		//stats: 收集待分配的出块奖励金额，每个结算周期可能不一样，当前节点已经不在101备选人列表中，委托用户不能参与本结算周期的委托分红。
		common.CollectBlockRewardData(head.Number.Uint64(), blockReward, true)
	} else {
		log.Warn("nodeID is not in verify list now, delegator has no block reward", "blockNumber", head.Number, "blockHash", blockHash, "nodeID", nodeID.String())
		//stats: 收集待分配的出块奖励金额，每个结算周期可能不一样
		common.CollectBlockRewardData(head.Number.Uint64(), blockReward, false)
	}

	if head.Coinbase != vm.RewardManagerPoolAddr {
		log.Debug("allocate package reward,block reward", "blockNumber", head.Number, "blockHash", blockHash, "nodeID", nodeID.String(),
			"coinBase", head.Coinbase.String(), "reward", reward)
		// 1. 出块节点，直接从激励池拿到出块奖励。
		state.SubBalance(vm.RewardManagerPoolAddr, reward)
		state.AddBalance(head.Coinbase, reward)
	} else {
		log.Warn("Coinbase equals RewardManagerPool address", "blockNumber", head.Number, "blockHash", blockHash, "nodeID", nodeID.String())
	}
	return nil
}

type DelegationInfoWithRewardPerList struct {
	DelegationInfo *staking.DelegationInfo
	RewardPerList  []*reward.DelegateRewardPer
}

func (d *DelegationInfoWithRewardPerList) String() string {
	v, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(v)
}

func NewDelegationInfoWithRewardPerList(delegationInfo *staking.DelegationInfo, rewardPerList []*reward.DelegateRewardPer) *DelegationInfoWithRewardPerList {
	return &DelegationInfoWithRewardPerList{delegationInfo, rewardPerList}
}

func (rmp *RewardMgrPlugin) GetDelegateRewardPerList(blockHash common.Hash, nodeID discover.NodeID, stakingNum, fromEpoch, toEpoch uint64) ([]*reward.DelegateRewardPer, error) {
	return getDelegateRewardPerList(blockHash, nodeID, stakingNum, fromEpoch, toEpoch, rmp.db)
}

func getDelegateRewardPerList(blockHash common.Hash, nodeID discover.NodeID, stakingNum, fromEpoch, toEpoch uint64, db snapshotdb.DB) ([]*reward.DelegateRewardPer, error) {
	keys := reward.DelegateRewardPerKeys(nodeID, stakingNum, fromEpoch, toEpoch)
	pers := make([]*reward.DelegateRewardPer, 0)
	for _, key := range keys {
		val, err := db.Get(blockHash, key)
		if err != nil {
			if err == snapshotdb.ErrNotFound {
				continue
			}
			return nil, err
		}
		list := new(reward.DelegateRewardPerList)
		if err := rlp.DecodeBytes(val, list); err != nil {
			return nil, err
		}
		for _, per := range list.Pers {
			if per.Epoch >= fromEpoch && per.Epoch <= toEpoch {
				pers = append(pers, per)
			}
		}
	}
	return pers, nil
}

func AppendDelegateRewardPer(blockHash common.Hash, nodeID discover.NodeID, stakingNum uint64, per *reward.DelegateRewardPer, db snapshotdb.DB) error {
	key := reward.DelegateRewardPerKey(nodeID, stakingNum, per.Epoch)
	list := reward.NewDelegateRewardPerList()
	val, err := db.Get(blockHash, key)
	if err != nil {
		if err != snapshotdb.ErrNotFound {
			return err
		}
	} else {
		if err := rlp.DecodeBytes(val, list); err != nil {
			return err
		}
	}

	list.AppendDelegateRewardPer(per)
	v, err := rlp.EncodeToBytes(list)
	if err != nil {
		return err
	}
	if err := db.Put(blockHash, key, v); err != nil {
		return err
	}
	log.Debug("append delegate rewardPer", "nodeID", nodeID.TerminalString(), "stkNum", stakingNum, "per", per)
	return nil
}

func UpdateDelegateRewardPer(blockHash common.Hash, nodeID discover.NodeID, stakingNum uint64, receives []reward.DelegateRewardReceipt, db snapshotdb.DB) error {
	if len(receives) == 0 {
		return nil
	}
	keys := reward.DelegateRewardPerKeys(nodeID, stakingNum, receives[0].Epoch, receives[len(receives)-1].Epoch)
	var reIndex int
	for _, key := range keys {
		val, err := db.Get(blockHash, key)
		if err != nil {
			if err == snapshotdb.ErrNotFound {
				continue
			}
			return err
		}
		list := new(reward.DelegateRewardPerList)
		if err := rlp.DecodeBytes(val, list); err != nil {
			return err
		}
		if len(receives)-reIndex < reward.DelegateRewardPerLength {
			reIndex += list.DecreaseTotalAmount(receives[reIndex:])
		} else {
			reIndex += list.DecreaseTotalAmount(receives[reIndex : reIndex+reward.DelegateRewardPerLength])
		}
		if list.IsChange() {
			log.Debug("update delegate reward per list", "nodeID", nodeID.TerminalString(), "stkNum", stakingNum, "list", list)
			if list.ShouldDel() {
				if err := db.Del(blockHash, key); err != nil {
					return err
				}
			} else {
				v, err := rlp.EncodeToBytes(list)
				if err != nil {
					return err
				}
				if err := db.Put(blockHash, key, v); err != nil {
					return err
				}
			}
		}
	}
	return nil
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

// 设置第year年（从0开始），激励池的初始值
func SetYearEndBalance(state xcom.StateDB, year uint32, balance *big.Int) {
	yearEndBalanceKey := reward.HistoryBalancePrefix(year)
	state.SetState(vm.RewardManagerPoolAddr, yearEndBalanceKey, balance.Bytes())
	log.Info("SetYearEndBalance", "address", vm.RewardManagerPoolAddr.Bech32(), "balance", balance)
}

// 查询第year年（从0开始），激励池的初始值
func GetYearEndBalance(state xcom.StateDB, year uint32) *big.Int {
	var balance = new(big.Int)
	yearEndBalanceKey := reward.HistoryBalancePrefix(year)
	bBalance := state.GetState(vm.RewardManagerPoolAddr, yearEndBalanceKey)
	log.Trace("show balance of reward pool at last year end", "lastYear", year, "amount", balance.SetBytes(bBalance))
	return balance.SetBytes(bBalance)
}

//在每个epoch之后，检查是否要增发
func (rmp *RewardMgrPlugin) runIncreaseIssuance(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	yes, err := xcom.IsYearEnd(blockHash, head.Number.Uint64())
	if nil != err {
		log.Error("Failed to execute runIncreaseIssuance function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return err
	}
	if yes {
		//最后一个epoch的最后一块，是增发块
		yearNumber, err := LoadChainYearNumber(blockHash, rmp.db)
		if nil != err {
			return err
		}
		yearNumber++
		//todo:chain_current_variables.chain_age，链年龄+1
		if err := StorageChainYearNumber(blockHash, rmp.db, yearNumber); nil != err {
			log.Error("Failed to execute runIncreaseIssuance function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return err
		}
		//todo: ppos_config.issue_cycle，发行周期（分钟）
		incIssuanceTime, err := xcom.LoadIncIssuanceTime(blockHash, rmp.db)
		if nil != err {
			return err
		}
		//todo: 执行增发逻辑
		if err := rmp.increaseIssuance(yearNumber, yearNumber-1, state, head.Number.Uint64(), blockHash); nil != err {
			return err
		}
		// After the increase issue is completed, update the number of rewards to be issued
		// todo: 这个不需要了
		if err := StorageRemainingReward(blockHash, rmp.db, GetYearEndBalance(state, yearNumber)); nil != err {
			log.Error("Failed to execute runIncreaseIssuance function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return err
		}
		// todo: chain_current_variables.issue_time，更新增发时间
		if err := xcom.StorageIncIssuanceTime(blockHash, rmp.db, incIssuanceTime+int64(xcom.AdditionalCycleTime()*uint64(minutes))); nil != err {
			log.Error("storage incIssuanceTime fail", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return err
		}
		remainReward, err := LoadRemainingReward(blockHash, rmp.db)
		if nil != err {
			return err
		}
		incIssuanceTime, err = xcom.LoadIncIssuanceTime(blockHash, rmp.db)
		if nil != err {
			return err
		}
		log.Info("Call CalcEpochReward, increaseIssuance successful", "currBlockNumber", head.Number, "currBlockHash", blockHash, "currBlockTime", head.Time.Int64(),
			"yearNumber", yearNumber, "incIssuanceTime", incIssuanceTime, "yearEndBalance", GetYearEndBalance(state, yearNumber), "remainingReward", remainReward)
	}
	return nil
}

//每个结算周期末，计算下一个epoch激励金；
//如果在这个结算周期末需要增发，则先增发，在执行计算下一个epoch的激励金
//特殊点：第一个epoch的激励金，是第一个epoch的第一个块来计算的。
func (rmp *RewardMgrPlugin) CalcEpochReward(blockHash common.Hash, head *types.Header, state xcom.StateDB) (*big.Int, *big.Int, error) {
	//获取当年剩余的激励池可用金额数,
	//重放难度系数：1， desc：用blockHash, reward.RemainingRewardKey作为UNIKEY存储到mysql中
	//注意：创世块中，要把第0年的激励池可用金额数写入
	//todo:replay: mysql 保存当前激励池在当前年度的可用金额，当前年度惩罚金额（激惩罚金额只能在下个年度使用）
	//todo:凡是链上保存的当前的某个变量的值，这些变量，可以放到一个表中，表示当前值，如chain_env表
	//todo:chain_env.reward_pool_available, chain_env.reward_pool_next_year_available
	remainReward, err := LoadRemainingReward(blockHash, rmp.db)
	if nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}

	//滑动窗口，用于计算平均出块时间的。窗口尺寸是第一个增发周期的块高数量（从0到第1年最后一个块高）。
	//第1次增发过后，每过1个结算周期，窗口起始块高就增加1个结算周期的块高数（注意从0块开始计算）
	//滑动窗口起始块，起始时间。链的创世块的块高=0，创世块的开始时间=0
	//重放难度系数：1， desc：用blockHash, reward.YearStartBlockNumberKey 作为UNIKEY存储到mysql中
	//重放难度系数：1， desc：用blockHash, reward.YearStartTimeKey 作为UNIKEY存储到mysql中
	//todo:2rd: mysql 保存滑动窗口的开始时间，开始块高.(起始可以都是0)
	//todo:chain_env.sliding_window_start_block,chain_env.sliding_window_start_time
	yearStartBlockNumber, yearStartTime, err := LoadYearStartTime(blockHash, rmp.db)
	if nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	//获取下一次增发的预期时间
	//todo:2rd: mysql 保存下个增发点的预期时间(起始可以是0)
	//todo:chain_env.issue_time，增发时间
	incIssuanceTime, err := xcom.LoadIncIssuanceTime(blockHash, rmp.db)
	if nil != err {
		log.Error("load incIssuanceTime fail", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	if yearStartTime == 0 { //特殊处理：说明是链的第1块，则需要计算1岁时的增发时间
		yearStartBlockNumber = head.Number.Uint64() //此时滑动窗口起始块高 yearStartBlockNumber=1
		yearStartTime = head.Time.Int64()           //此时滑动窗口起始时间 yearStartTime=区块1的时间
		//计算链下一年的预期增发时间点（当前区块时间 + 增发周期长度）
		//todo:2rd: mysql 保存gengeis.json的配置，增发周期（注意时间单位：分钟)，设置新表：ppos_env
		//todo: ppos_env.issue_cycle，发行周期（分钟）
		incIssuanceTime = yearStartTime + int64(xcom.AdditionalCycleTime()*uint64(minutes))
		if err := xcom.StorageIncIssuanceTime(blockHash, rmp.db, incIssuanceTime); nil != err {
			log.Error("storage incIssuanceTime fail", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return nil, nil, err
		}
		//保存滑动窗口开始时间，开始区块
		if err := StorageYearStartTime(blockHash, rmp.db, yearStartBlockNumber, yearStartTime); nil != err {
			log.Error("Storage year start time and block height failed", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return nil, nil, err
		}
		//获取上年末激励池余额（由于是第1快，所以也就是创世块中设置的、第0年的激励池金额）
		//todo: 这个究竟有什么作用？
		//每年开始的remainReward=增发进激励池的+上年的balance余额（惩罚的+上年没用完的（如果有））
		remainReward = GetYearEndBalance(state, 0)
		log.Info("Call CalcEpochReward, First calculation", "currBlockNumber", head.Number, "currBlockHash", blockHash, "currBlockTime", head.Time.Int64(),
			"yearStartBlockNumber", yearStartBlockNumber, "yearStartTime", yearStartTime, "incIssuanceTime", incIssuanceTime, "remainReward", remainReward)
	}

	//计算链年龄，也是链的增发年度数。满1年，就是1；未满1年，就是0
	//todo:chain_env.age，链年龄
	yearNumber, err := LoadChainYearNumber(blockHash, rmp.db)
	if nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	// When the first issuance is completed
	// Each settlement cycle needs to update the year start time,
	// which is used to calculate the average annual block production rate
	//计算Epoch有多少个区块（这是固定数量的）
	//todo: ppos_env.epoch_size，Epoch有多少个区块
	epochBlocks := xutil.CalcBlocksEachEpoch()
	if yearNumber > 0 { //满1年（说明增发过了）
		//增发过，则取最近的增发块高（这个值，只有在增发前一个结算周期的最后一个块，才会修改）
		//todo:chain_env.issue_block，增发区块
		incIssuanceNumber, err := xcom.LoadIncIssuanceNumber(blockHash, rmp.db)
		if nil != err {
			return nil, nil, err
		}
		addition := true
		//滑动窗口起始块高，说明窗口起始块还没调整过，刚完成第一次增发, yearNumber刚刚变成 1
		//滑动窗口开始时，size是不确定的，size的调整是以epochBlocks为单位的。窗口起始块高是0（或者1），每过1个结算周期，结束块高增长一个epochBlocks；
		//当第一次增发后，增发完成的那个块高，就是滑动窗口的结束块高，此时，滑动窗口的size确定。
		//以后，每过一个结算周期，起始块高+一个epochBlocks；结算块高就是当前块高。
		if yearStartBlockNumber == 1 {
			if head.Number.Uint64() <= incIssuanceNumber {
				//由于CalcEpochReward只有在结算周期末执行，所以，这里出现 < 的情况应该没有
				//此时刚到第一次增发点，不用调整滑动窗口起始块
				addition = false //
			}
		}
		//默认每过一个结算周期，都要调整滑动窗口的起始块高
		if addition {
			if yearStartBlockNumber == 1 {
				//第一次调整时，是第一次增发后的第一个结算周期末，调整后起始块位于0岁的第1个结算周期末（因为yearStartBlockNumber=1，这个起点不是0岁0周期末，而是0岁1周期起始块）。
				yearStartBlockNumber += epochBlocks - 1
			} else {
				yearStartBlockNumber += epochBlocks
			}
			//调整滑动窗口的起始时间
			yearStartTime = snapshotdb.GetDBBlockChain().GetHeaderByNumber(yearStartBlockNumber).Time.Int64()
			if err := StorageYearStartTime(blockHash, rmp.db, yearStartBlockNumber, yearStartTime); nil != err {
				log.Error("Storage year start time and block height failed", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
				return nil, nil, err
			}
			log.Debug("Call CalcEpochReward, Adjust the sampling range of the block time", "currBlockNumber", head.Number, "currBlockHash", blockHash, "currBlockTime", head.Time.Int64(),
				"epochBlocks", epochBlocks, "yearNumber", yearNumber, "yearStartBlockNumber", yearStartBlockNumber, "yearStartTime", yearStartTime)
		}
	}

	// First calculation, calculated according to the default block interval.
	// In each subsequent settlement cycle, an average block generation interval needs to be calculated.
	avgPackTime := xcom.Interval() * uint64(millisecond)
	if head.Number.Uint64() > yearStartBlockNumber {
		diffNumber := head.Number.Uint64() - yearStartBlockNumber
		diffTime := head.Time.Int64() - yearStartTime

		avgPackTime = uint64(diffTime) / diffNumber
		log.Debug("Call CalcEpochReward, Calculate the average block production time in the previous year", "currBlockNumber", head.Number, "currBlockHash", blockHash,
			"currBlockTime", head.Time.Int64(), "yearStartBlockNumber", yearStartBlockNumber, "yearStartTime", yearStartTime, "diffNumber", diffNumber, "diffTime", diffTime,
			"avgPackTime", avgPackTime)
	}
	if err := xcom.StorageAvgPackTime(blockHash, rmp.db, avgPackTime); nil != err {
		log.Error("Failed to execute StorageAvgPackTime function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "avgPackTime", avgPackTime, "err", err)
		return nil, nil, err
	}

	epochTotalReward := new(big.Int)
	// If the expected increase issue time is exceeded,
	// the increase issue time will be postponed for one settlement cycle,
	// and the remaining rewards will all be issued in the next settlement cycle
	remainEpoch := 1
	if head.Time.Int64() >= incIssuanceTime {
		//如果当前区块时间>=下个增发时间点，则下个增发被推迟一个epoch，激励池中本年度可用资金，将在接下来的epoch中分配完。
		epochTotalReward.Add(epochTotalReward, remainReward)
		remainReward = new(big.Int)
		log.Info("Call CalcEpochReward, The current time has exceeded the expected additional issue time", "currBlockNumber", head.Number, "currBlockHash", blockHash,
			"currBlockTime", head.Time.Int64(), "incIssuanceTime", incIssuanceTime, "epochTotalReward", epochTotalReward)
	} else {

		//本年度剩余时间（下次增发的预期时间-当前区块时间）
		remainTime := incIssuanceTime - head.Time.Int64()

		//本年度剩余区块（剩余时间/出块平均间隔）
		remainBlocks := math.Ceil(float64(remainTime) / float64(avgPackTime))
		if remainBlocks > float64(epochBlocks) { //剩余区块 > epoch区块数，则剩余epoch>1
			remainEpoch = int(math.Ceil(remainBlocks / float64(epochBlocks))) //向上取整
		}
		//计算剩余epoch的每个epoch的激励金额
		epochTotalReward = new(big.Int).Div(remainReward, new(big.Int).SetInt64(int64(remainEpoch)))
		// Subtract the total reward for the next cycle to calculate the remaining rewards to be issued
		//如果剩余epoch=1，则激励池剩余的本年度可用金额，将全部被用来分配给剩余区块，
		//remainReward将等于0
		remainReward = remainReward.Sub(remainReward, epochTotalReward)
		log.Debug("Call CalcEpochReward, Calculation of rewards for the next settlement cycle", "currBlockNumber", head.Number, "currBlockHash", blockHash,
			"currBlockTime", head.Time.Int64(), "incIssuanceTime", incIssuanceTime, "remainTime", remainTime, "remainBlocks", remainBlocks, "epochBlocks", epochBlocks,
			"remainEpoch", remainEpoch, "remainReward", remainReward, "epochTotalReward", epochTotalReward)
	}
	// If the last settlement cycle is left, record the increaseIssuance block height
	if remainEpoch == 1 { //本年度还有最后一个epoch（当前块是倒数第二个epoch的最后一块）
		//此时，可以确定下一个增发的区块高度了
		incIssuanceNumber := new(big.Int).Add(head.Number, new(big.Int).SetUint64(epochBlocks)).Uint64()
		if err := xcom.StorageIncIssuanceNumber(blockHash, rmp.db, incIssuanceNumber); nil != err {
			return nil, nil, err
		}
		log.Info("Call CalcEpochReward, IncIssuanceNumber stored successfully", "currBlockNumber", head.Number, "currBlockHash", blockHash,
			"epochBlocks", epochBlocks, "incIssuanceNumber", incIssuanceNumber)
	}
	// Get the total block reward and Staking reward for each settlement cycle
	// 下一个epoch的总出块激励金，以及质押激励金。
	epochTotalNewBlockReward := percentageCalculation(epochTotalReward, xcom.NewBlockRewardRate())
	//下一个epoch的总质押奖励，总质押激励 = 总激励 - 出块激励。
	epochTotalStakingReward := new(big.Int).Sub(epochTotalReward, epochTotalNewBlockReward)
	if err := StorageRemainingReward(blockHash, rmp.db, remainReward); nil != err {
		log.Error("Failed to execute CalcEpochReward function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	//下个epoch，每个块的出块奖励
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
		if err == snapshotdb.ErrNotFound {
			return new(big.Int).SetUint64(0), nil
		}
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
		if err == snapshotdb.ErrNotFound {
			return new(big.Int).SetUint64(0), nil
		}
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
