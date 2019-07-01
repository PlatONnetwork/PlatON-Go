package plugin

import (
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


func GetAwardMgrInstance() *rewardMgrPlugin {
	if RewardMgrPlugin == nil {
		RewardMgrPlugin = & rewardMgrPlugin {}
	}
	return RewardMgrPlugin
}

// BeginBlock does something like check input params before execute transactions,
// in LockRepoPlugin it does nothing.
func (rmp *rewardMgrPlugin) BeginBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) (bool, error) {
	return false, nil
}

// EndBlock invoke releaseLockRepo
func (rmp *rewardMgrPlugin) EndBlock(blockHash common.Hash, head *types.Header, state xcom.StateDB) (bool, error) {

	blockNumber := head.Number.Uint64()

	// check is time to increaseIssuance
	if xutil.IsYearEnd(blockNumber) {
		log.Info("last block at end of year", "year", blockNumber / 12 * 30 * 4 * xcom.ConsensusSize * xcom.EpochSize)

		if success, err := rmp.increaseIssuance(head, state); err != nil {
			return success, err
		}
	}

	pledgeReward, newBlockReward, err := rmp.computePeriodAward(head)
	if err != nil {
		return false, err
	}

	// check is time to reward pledge
	if xutil.IsSettlementPeriod(blockNumber) {
		log.Info("settle block", "period", xutil.CalculateRound(blockNumber))

		if success, err := rmp.rewardPledge(head, pledgeReward, state); err != nil {
			return success, err
		}
	}

	// every block need rewardNewBLock
	log.Debug("time to rewardNewBlock at last", "blockNumber", blockNumber, "hash", head.Hash())
	addr, success, err := rmp.rewardNewBlock(head, newBlockReward, state)
	if err != nil {
		return success, err
	} else {
		head.Coinbase = addr
	}

	return true, nil
}

// Comfired does nothing
func (rmp *rewardMgrPlugin) Confirmed(block *types.Block) (bool, error) {
	return true, nil
}

// increaseIssuance does nothing
func (rmp *rewardMgrPlugin) increaseIssuance(head *types.Header, state xcom.StateDB) (bool, error) {

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

	return true, nil
}

// rewardPledge does nothing
func (rmp *rewardMgrPlugin) rewardPledge(head *types.Header, reward *big.Int, state xcom.StateDB) (bool, error) {

	// stakingPlugin.GetVerifierList()  获取列表

	var list []*xcom.Candidate

	for index := 0; index < len(list); index++ {
		addr := list[index].BenifitAddress

		if addr != vm.RewardManagerPoolAddr {
			state.SubBalance(vm.RewardManagerPoolAddr, reward)
			state.AddBalance(addr, reward)
		}
	}

	return true, nil
}

// rewardNewBlock does nothing, return coinbase address and error
func (rmp *rewardMgrPlugin) rewardNewBlock(head *types.Header, reward *big.Int, state xcom.StateDB) (common.Address, bool, error) {

	sign := head.Extra[32:97]

	pubKey, err := crypto.SigToPub(head.Hash().Bytes(), sign)
	if err != nil {
		log.Error("failed to ecrecover sign to public key", "blockNumber", head.Number, "hash", head.Hash(), "sign", sign)
		return  common.Address{0}, false, err
	}

	nodeID := discover.PubkeyID(pubKey)

	log.Debug("node", "NodeID", nodeID)

	var can xcom.Candidate // stakingPlugin.GetCandidate(nodeID)

	rewardAddr := can.BenifitAddress

	if rewardAddr != vm.RewardManagerPoolAddr {
		state.SubBalance(vm.RewardManagerPoolAddr, reward)
		state.AddBalance(rewardAddr, reward)
	}

	return rewardAddr, true, nil
}

// computePeriodAward does nothing
func (rmp *rewardMgrPlugin) computePeriodAward(head *types.Header) (*big.Int, *big.Int, error) {
	var (
		pledgeReward  *big.Int
		newBlockReward *big.Int
	)


	return pledgeReward, newBlockReward, nil
}

