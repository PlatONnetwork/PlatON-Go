package plugin

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/reward"
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type rewardMgrPlugin struct {

}

var (
	RewardMgrPlugin *rewardMgrPlugin = nil
)

func RewardMgrInstance() *rewardMgrPlugin {
	if RewardMgrPlugin == nil {
		RewardMgrPlugin = & rewardMgrPlugin {}
	}
	return RewardMgrPlugin
}

// BeginBlock does something like check input params before execute transactions,
// in rewardMgrPlugin it does nothing.
func (rmp *rewardMgrPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	return nil
}

// EndBlock will handle reward work, if current block is end of year, increase issuance; if it is time to settle,
// reward staking. At last reward worker for create new block, this is necessary.
func (rmp *rewardMgrPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {

	blockNumber := head.Number.Uint64()
	year := xutil.CalculateYear(blockNumber)

	if xutil.IsYearEnd(blockNumber) {
		log.Info("last block at end of year", "year", year)
		if err := rmp.increaseIssuance(uint32(year), state); err != nil {
			return err
		}
	}

	stakingReward, newBlockReward, _ := rmp.calculateExpectReward(uint32(year), state)

	if xutil.IsSettlementPeriod(blockNumber) {
		log.Info("settle block", "period", xutil.CalculateRound(blockNumber))

		_ = rmp.rewardStaking(head, stakingReward, state)
	}

	log.Info("rewardNewBlock at last", "blockNumber", blockNumber, "hash", head.Hash())
	if addr, err := rmp.rewardNewBlock(head, newBlockReward, state); err != nil {
		return err
	} else {
		head.Coinbase = addr
	}

	return nil
}

// Confirmed does nothing
func (rmp *rewardMgrPlugin) Confirmed(block *types.Block) error {
	return nil
}

// increaseIssuance used for increase issuance at the end of each year
func (rmp *rewardMgrPlugin) increaseIssuance(year uint32, state xcom.StateDB) error {
	var rate *big.Int

	hisIssuance := GetLatestCumulativeIssue(state)

	currIssuance := hisIssuance.Div(hisIssuance, rate.SetUint64(40))
	develop := currIssuance.Div(currIssuance, rate.SetUint64(5))
	rewards := currIssuance.Sub(currIssuance, develop)

	hisIssuance = hisIssuance.Add(hisIssuance, currIssuance)
	SetYearEndCumulativeIssue(state, year, hisIssuance)

	state.AddBalance(vm.CommunityDeveloperFoundation, develop)
	state.AddBalance(vm.RewardManagerPoolAddr, rewards)

	return nil
}

// rewardStaking used for reward staking at the settle block
func (rmp *rewardMgrPlugin) rewardStaking(head *types.Header, reward *big.Int, state xcom.StateDB) error {

	// stakingPlugin.GetVerifierList()

	var list []*staking.Candidate

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
func (rmp *rewardMgrPlugin) rewardNewBlock(head *types.Header, reward *big.Int, state xcom.StateDB) (common.Address, error) {

	sign := head.Extra[32:97]

	pubKey, err := crypto.SigToPub(head.Hash().Bytes(), sign)
	if err != nil {
		log.Error("failed to ecrecover sign to public key", "blockNumber", head.Number, "hash", head.Hash(), "sign", sign)
		return  common.Address{0}, common.NewSysError(err.Error())
	}

	nodeID := discover.PubkeyID(pubKey)
	log.Debug("node", "NodeID", nodeID)

	var can staking.Candidate // stakingPlugin.GetCandidate(nodeID)
	rewardAddr := can.BenifitAddress
	if rewardAddr != vm.RewardManagerPoolAddr {
		state.SubBalance(vm.RewardManagerPoolAddr, reward)
		state.AddBalance(rewardAddr, reward)
	}

	return rewardAddr, nil
}

// calculateExpectReward used for calculate the stakingReward and newBlockReward that should be send in each corresponding period
func (rmp *rewardMgrPlugin) calculateExpectReward(year uint32, state xcom.StateDB) (*big.Int, *big.Int, error) {
	var (
		stakingReward  *big.Int
		newBlockReward *big.Int
	)

	expectNewBlocks := int64(365) * 24 * 3600 / 1
	expectEpochs := int64(365) * 24 * 3600 / int64(xcom.ConsensusSize * xcom.EpochSize)

	issuance := GetLatestCumulativeIssue(state)
	totalNewBlockReward := issuance.Div(issuance, big.NewInt(4))
	totalStakingReward := issuance.Sub(issuance, totalNewBlockReward)

	newBlockReward = totalNewBlockReward.Div(totalNewBlockReward, big.NewInt(expectNewBlocks))
	stakingReward = totalStakingReward.Div(totalStakingReward, big.NewInt(expectEpochs))

	return stakingReward, newBlockReward, nil
}

// SetYearEndCumulativeIssue used for set historical cumulative increase at the end of the year
func SetYearEndCumulativeIssue(state xcom.StateDB, year uint32, total *big.Int) {
	yearEndIncreaseKey := reward.GetHistoryIncreaseKey(year)
	state.SetState(vm.RewardManagerPoolAddr, yearEndIncreaseKey, total.Bytes())
}

// GetLatestCumulativeIssue used for get the cumulative issuance in the most recent year
func GetLatestCumulativeIssue(state xcom.StateDB) *big.Int {
	var issue *big.Int

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
	var issue *big.Int

	hisIncreaseKey := reward.GetHistoryIncreaseKey(year)
	bIssue := state.GetState(vm.RewardManagerPoolAddr, hisIncreaseKey)

	return issue.SetBytes(bIssue)
}