package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/x/staking"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/vm"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type rewardMgrPlugin struct {

}

var (
	RewardMgrPlugin *rewardMgrPlugin = nil
)

func GetRewardMgrInstance() *rewardMgrPlugin {
	if RewardMgrPlugin == nil {
		RewardMgrPlugin = & rewardMgrPlugin {}
	}
	return RewardMgrPlugin
}

// BeginBlock does something like check input params before execute transactions,
// in LockRepoPlugin it does nothing.
func (rmp *rewardMgrPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {
	return nil
}

// EndBlock invoke releaseLockRepo
func (rmp *rewardMgrPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) error {

	blockNumber := head.Number.Uint64()

	// check is time to increaseIssuance
	if xutil.IsYearEnd(blockNumber) {
		log.Info("last block at end of year", "year", xutil.CalculateYears(blockNumber))

		if err := rmp.increaseIssuance(head, state); err != nil {
			return err
		}
	}

	stakingReward, newBlockReward, err := rmp.calculateExpectReward(head)
	if err != nil {
		return err
	}

	// check is time to reward staking
	if xutil.IsSettlementPeriod(blockNumber) {
		log.Info("settle block", "period", xutil.CalculateRound(blockNumber))

		if err := rmp.rewardStaking(head, stakingReward, state); err != nil {
			return err
		}
	}

	// every block need reward NewBLock
	log.Debug("time to rewardNewBlock at last", "blockNumber", blockNumber, "hash", head.Hash())
	addr, err := rmp.rewardNewBlock(head, newBlockReward, state)
	if err != nil {
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

// increaseIssuance does nothing
func (rmp *rewardMgrPlugin) increaseIssuance(head *types.Header, state xcom.StateDB) error {

	year := xutil.CalculateYears(head.Number.Uint64())

	switch {
	case year == 1:
		// 第一年末开始第二年的增发
		// 获取创世块上的余额, 增发比例为2.5

		state.AddBalance(vm.RewardManagerPoolAddr, big.NewInt(100))   // 4/5
		state.AddBalance(vm.GovContractAddr, big.NewInt(20))     // 1/5

	case year > 1:
		// 从第三年开始,每年获取上一年末已经发行的总量，这个有点困难，可能要记录下到目前为止发行的总量

		state.AddBalance(vm.RewardManagerPoolAddr, big.NewInt(100))   // 4/5
		state.AddBalance(vm.GovContractAddr, big.NewInt(20))     // 1/5

	}

	return nil
}

// rewardStaking does nothing
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

// rewardNewBlock does nothing, return coinbase address and error
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

// getLastBalance used for get balance vm.RewardManagerPoolAddr of last year end
func getLastBalance() *big.Int {
	last, _ := new(big.Int).SetString("10000000000000000000000000000", 10)
	return last
}

// calculateExpectReward does nothing
func (rmp *rewardMgrPlugin) calculateExpectReward(head *types.Header) (*big.Int, *big.Int, error) {
	var (
		stakingReward  *big.Int
		newBlockReward *big.Int
	)

	expectNewBlocks := int64(365) * 24 * 3600 / 1
	expectEpochs := int64(365) * 24 * 3600 / int64(xcom.ConsensusSize * xcom.EpochSize)

	balance := getLastBalance()
	totalNewBlockReward := balance.Div(balance, big.NewInt(4))
	totalStakingReward := balance.Sub(balance, totalNewBlockReward)

	newBlockReward = totalNewBlockReward.Div(totalNewBlockReward, big.NewInt(expectNewBlocks))
	stakingReward = totalStakingReward.Div(totalNewBlockReward, big.NewInt(expectEpochs))

	return stakingReward, newBlockReward, nil
}

