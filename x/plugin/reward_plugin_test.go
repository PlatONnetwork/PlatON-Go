package plugin

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/x/reward"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

	"github.com/PlatONnetwork/PlatON-Go/x/xcom"
	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

func TestRewardPlugin(t *testing.T) {
	var plugin = new(rewardMgrPlugin)
	mockDB := buildStateDB(t)

	t.Run("CalculateExpectReward", func(t *testing.T) {
		//	log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

		yearBalance := big.NewInt(1E18)
		rate := xcom.NewBlockRewardRate()
		epochs := xutil.EpochsPerYear()
		blocks := xutil.CalcBlocksEachYear()
		thisYear, lastYear := uint32(2), uint32(1)
		expectNewBlockReward := percentageCalculation(yearBalance, rate)
		SetYearEndBalance(mockDB, lastYear, yearBalance)
		mockDB.AddBalance(vm.RewardManagerPoolAddr, yearBalance)

		stakingReward, newBlockReward := plugin.calculateExpectReward(thisYear, lastYear, mockDB)

		expectStakingReward := new(big.Int).Sub(yearBalance, expectNewBlockReward)

		assert.Equal(t, expectNewBlockReward.Div(expectNewBlockReward, big.NewInt(int64(blocks))), newBlockReward)
		assert.Equal(t, expectStakingReward.Div(expectStakingReward, big.NewInt(int64(epochs))), stakingReward)

		list := make(staking.ValidatorExQueue, 0)
		for _, value := range addrArr {
			list = append(list, &staking.ValidatorEx{
				BenefitAddress: value,
			})
		}

		plugin.rewardStakingByValidatorList(mockDB, list, stakingReward)
		everyValidatorReward := new(big.Int).Div(stakingReward, big.NewInt(int64(len(list))))
		for _, value := range list {
			assert.Equal(t, everyValidatorReward, mockDB.GetBalance(value.BenefitAddress))
		}

		account := common.HexToAddress("0xeef233120ce31b3fac20dac379db243021a5234")
		plugin.allocatePackageBlock(10, common.ZeroHash, account, newBlockReward, mockDB)

		assert.Equal(t, newBlockReward, mockDB.GetBalance(account))

		lastIssue := GetHistoryCumulativeIssue(mockDB, lastYear)

		plugin.increaseIssuance(thisYear, lastYear, mockDB)

		newIssue := GetHistoryCumulativeIssue(mockDB, thisYear)

		tmp := new(big.Int).Sub(newIssue, lastIssue)
		assert.Equal(t, lastIssue, tmp.Mul(tmp, big.NewInt(IncreaseIssue)))

		lastYearIssue := new(big.Int).SetBytes(mockDB.GetState(vm.RewardManagerPoolAddr, reward.GetHistoryIncreaseKey(lastYear)))

		if plugin.isLessThanFoundationYear(thisYear) {
			mockDB.GetBalance(xcom.CDFAccount())

		} else {
			mockDB.GetBalance(xcom.CDFAccount())
			mockDB.GetBalance(xcom.PlatONFundAccount())
		}
		mockDB.GetBalance(vm.RewardManagerPoolAddr)

		thisYearIssue := new(big.Int).SetBytes(mockDB.GetState(vm.RewardManagerPoolAddr, reward.GetHistoryIncreaseKey(thisYear)))

		assert.Equal(t, new(big.Int).Sub(thisYearIssue, lastYearIssue), new(big.Int).Div(lastYearIssue, big.NewInt(IncreaseIssue)))

	})

}
