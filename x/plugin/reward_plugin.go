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
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"math"
	"math/big"
	"sync"

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
func (rmp *RewardMgrPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	blockNumber := head.Number.Uint64()

	packageReward := new(big.Int)
	stakingReward := new(big.Int)
	var err error

	if head.Number.Uint64() == common.Big1.Uint64() {
		packageReward, stakingReward, err = rmp.CalcEpochReward(blockHash, head, state)
		if nil != err {
			log.Error("Execute CalcEpochReward fail", "blockNumber", head.Number.Uint64(), "blockHash", blockHash.TerminalString(), "err", err)
			return err
		}
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

	if err := rmp.AllocatePackageBlock(blockHash, head, packageReward, state); err != nil {
		return err
	}

	if xutil.IsEndOfEpoch(blockNumber) {
		verifierList, err := rmp.AllocateStakingReward(blockNumber, blockHash, stakingReward, state)
		if err != nil {
			return err
		}

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
func (rmp *RewardMgrPlugin) increaseIssuance(thisYear, lastYear uint32, state xcom.StateDB, blockNumber uint64, blockHash common.Hash) error {
	var currIssuance *big.Int
	//issuance increase
	{
		histIssuance := GetHistoryCumulativeIssue(state, lastYear)
		increaseIssuanceRatio, err := gov.GovernIncreaseIssuanceRatio(blockNumber, blockHash)
		if nil != err {
			log.Error("Failed to increaseIssuance, call GovernIncreaseIssuanceRatio is failed", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
				"histIssuance", histIssuance, "err", err)
			return err
		}

		tmp := new(big.Int).Mul(histIssuance, big.NewInt(int64(increaseIssuanceRatio)))
		currIssuance = tmp.Div(tmp, big.NewInt(10000))

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
	if err := rmp.rewardStakingByValidatorList(state, verifierList, sreward); err != nil {
		log.Error("reward staking by validator list fail", "err", err, "bn", blockNumber, "bh", blockHash)
		return nil, err
	}

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

			per := reward.NewDelegateRewardPer(currentEpoch, verifier.CurrentEpochDelegateReward, verifier.DelegateTotal)
			if err := AppendDelegateRewardPer(blockHash, verifier.NodeId, verifier.StakingBlockNum, per, rmp.db); err != nil {
				log.Error("call handleDelegatePerReward fail AppendDelegateRewardPer", "blockNumber", blockNumber, "blockHash", blockHash.TerminalString(),
					"nodeId", verifier.NodeId.TerminalString(), "err", err, "CurrentEpochDelegateReward", verifier.CurrentEpochDelegateReward, "delegateTotal", verifier.DelegateTotal)
				return err
			}
			currentEpochDelegateReward := new(big.Int).Set(verifier.CurrentEpochDelegateReward)

			verifier.PrepareNextEpoch()
			canAddr, err := xutil.NodeId2Addr(verifier.NodeId)
			if nil != err {
				log.Error("Failed to handleDelegatePerReward on rewardMgrPlugin: nodeId parse addr failed",
					"blockNumber", blockNumber, "blockHash", blockHash, "nodeID", verifier.NodeId.String(), "err", err)
				return err
			}
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

func (rmp *RewardMgrPlugin) rewardStakingByValidatorList(state xcom.StateDB, list []*staking.Candidate, reward *big.Int) error {
	validatorNum := int64(len(list))
	everyValidatorReward := new(big.Int).Div(reward, big.NewInt(validatorNum))

	log.Debug("calculate validator staking reward", "validator length", validatorNum, "everyOneReward", everyValidatorReward)
	totalValidatorReward, totalValidatorDelegateReward := new(big.Int), new(big.Int)

	for _, value := range list {
		delegateReward, stakingReward := new(big.Int), new(big.Int).Set(everyValidatorReward)
		if value.ShouldGiveDelegateReward() {
			delegateReward, stakingReward = rmp.CalDelegateRewardAndNodeReward(everyValidatorReward, value.RewardPer)
			totalValidatorDelegateReward.Add(totalValidatorDelegateReward, delegateReward)
			log.Debug("allocate delegate reward of staking one-by-one", "nodeId", value.NodeId.TerminalString(), "staking reward", stakingReward, "per", value.RewardPer, "delegateReward", delegateReward)
			//the  CurrentEpochDelegateReward will use by cal delegate reward Per
			value.CurrentEpochDelegateReward.Add(value.CurrentEpochDelegateReward, delegateReward)
		}
		if value.BenefitAddress != vm.RewardManagerPoolAddr {
			log.Debug("allocate staking reward one-by-one", "nodeId", value.NodeId.String(),
				"benefitAddress", value.BenefitAddress.String(), "staking reward", stakingReward)
			state.AddBalance(value.BenefitAddress, stakingReward)
			totalValidatorReward.Add(totalValidatorReward, stakingReward)
		}
	}
	state.AddBalance(vm.DelegateRewardPoolAddr, totalValidatorDelegateReward)
	state.SubBalance(vm.RewardManagerPoolAddr, new(big.Int).Add(totalValidatorDelegateReward, totalValidatorReward))
	return nil
}

func (rmp *RewardMgrPlugin) getBlockMinderAddress(blockHash common.Hash, head *types.Header) (discover.NodeID, common.NodeAddress, error) {
	if blockHash == common.ZeroHash {
		return rmp.nodeID, rmp.nodeADD, nil
	}
	sign := head.Extra[32:97]
	sealhash := head.SealHash().Bytes()
	pk, err := crypto.SigToPub(sealhash, sign)
	if err != nil {
		return discover.ZeroNodeID, common.ZeroNodeAddr, err
	}
	return discover.PubkeyID(pk), crypto.PubkeyToNodeAddress(*pk), nil
}

// AllocatePackageBlock used for reward new block. it returns coinbase and error
func (rmp *RewardMgrPlugin) AllocatePackageBlock(blockHash common.Hash, head *types.Header, reward *big.Int, state xcom.StateDB) error {
	nodeID, add, err := rmp.getBlockMinderAddress(blockHash, head)
	if err != nil {
		log.Error("AllocatePackageBlock getBlockMinderAddress fail", "err", err, "blockNumber", head.Number, "blockHash", blockHash)
		return err
	}

	currVerifier, err := rmp.stakingPlugin.IsCurrVerifier(blockHash, head.Number.Uint64(), nodeID, false)
	if err != nil {
		log.Error("AllocatePackageBlock IsCurrVerifier fail", "err", err, "blockNumber", head.Number, "blockHash", blockHash)
		return err
	}
	if currVerifier {
		cm, err := rmp.stakingPlugin.GetCanMutable(blockHash, add)
		if err != nil {
			log.Error("AllocatePackageBlock GetCanMutable fail", "err", err, "blockNumber", head.Number, "blockHash", blockHash, "add", add)
			return err
		}
		if cm.ShouldGiveDelegateReward() {
			delegateReward := new(big.Int).SetUint64(0)
			delegateReward, reward = rmp.CalDelegateRewardAndNodeReward(reward, cm.RewardPer)

			state.SubBalance(vm.RewardManagerPoolAddr, delegateReward)
			state.AddBalance(vm.DelegateRewardPoolAddr, delegateReward)
			cm.CurrentEpochDelegateReward.Add(cm.CurrentEpochDelegateReward, delegateReward)
			log.Debug("allocate package reward, delegate reward", "blockNumber", head.Number, "blockHash", blockHash, "delegateReward", delegateReward, "epochDelegateReward", cm.CurrentEpochDelegateReward)

			if err := rmp.stakingPlugin.db.SetCanMutableStore(blockHash, add, cm); err != nil {
				log.Error("AllocatePackageBlock SetCanMutableStore fail", "err", err, "blockNumber", head.Number, "blockHash", blockHash)
				return err
			}
		}
	}

	if head.Coinbase != vm.RewardManagerPoolAddr {

		log.Debug("allocate package reward,block reward", "blockNumber", head.Number, "blockHash", blockHash,
			"coinBase", head.Coinbase.String(), "reward", reward)

		state.SubBalance(vm.RewardManagerPoolAddr, reward)
		state.AddBalance(head.Coinbase, reward)
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

func (rmp *RewardMgrPlugin) runIncreaseIssuance(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	yes, err := xcom.IsYearEnd(blockHash, head.Number.Uint64())
	if nil != err {
		log.Error("Failed to execute runIncreaseIssuance function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return err
	}
	if yes {
		yearNumber, err := LoadChainYearNumber(blockHash, rmp.db)
		if nil != err {
			return err
		}
		yearNumber++
		if err := StorageChainYearNumber(blockHash, rmp.db, yearNumber); nil != err {
			log.Error("Failed to execute runIncreaseIssuance function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return err
		}
		incIssuanceTime, err := xcom.LoadIncIssuanceTime(blockHash, rmp.db)
		if nil != err {
			return err
		}
		if err := rmp.increaseIssuance(yearNumber, yearNumber-1, state, head.Number.Uint64(), blockHash); nil != err {
			return err
		}
		// After the increase issue is completed, update the number of rewards to be issued
		if err := StorageRemainingReward(blockHash, rmp.db, GetYearEndBalance(state, yearNumber)); nil != err {
			log.Error("Failed to execute runIncreaseIssuance function", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
			return err
		}
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
	incIssuanceTime, err := xcom.LoadIncIssuanceTime(blockHash, rmp.db)
	if nil != err {
		log.Error("load incIssuanceTime fail", "currentBlockNumber", head.Number, "currentBlockHash", blockHash.TerminalString(), "err", err)
		return nil, nil, err
	}
	if yearStartTime == 0 {
		yearStartBlockNumber = head.Number.Uint64()
		yearStartTime = head.Time.Int64()
		incIssuanceTime = yearStartTime + int64(xcom.AdditionalCycleTime()*uint64(minutes))
		if err := xcom.StorageIncIssuanceTime(blockHash, rmp.db, incIssuanceTime); nil != err {
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
		incIssuanceNumber, err := xcom.LoadIncIssuanceNumber(blockHash, rmp.db)
		if nil != err {
			return nil, nil, err
		}
		addition := true
		if yearStartBlockNumber == 1 {
			if head.Number.Uint64() <= incIssuanceNumber {
				addition = false
			}
		}
		if addition {
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
		epochTotalReward.Add(epochTotalReward, remainReward)
		remainReward = new(big.Int)
		log.Info("Call CalcEpochReward, The current time has exceeded the expected additional issue time", "currBlockNumber", head.Number, "currBlockHash", blockHash,
			"currBlockTime", head.Time.Int64(), "incIssuanceTime", incIssuanceTime, "epochTotalReward", epochTotalReward)
	} else {
		remainTime := incIssuanceTime - head.Time.Int64()

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
	// If the last settlement cycle is left, record the increaseIssuance block height
	if remainEpoch == 1 {
		incIssuanceNumber := new(big.Int).Add(head.Number, new(big.Int).SetUint64(epochBlocks)).Uint64()
		if err := xcom.StorageIncIssuanceNumber(blockHash, rmp.db, incIssuanceNumber); nil != err {
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
