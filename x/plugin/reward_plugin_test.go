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
	"math/rand"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common/mock"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/reward"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

func buildTestStakingData(epochStart, epochEnd uint64) (staking.ValidatorQueue, error) {
	validatorQueue := make(staking.ValidatorQueue, xcom.MaxValidators())
	for i := 0; i < 50; i++ {
		privateKey, err := crypto.GenerateKey()
		if nil != err {
			return nil, err
		}
		addr := crypto.PubkeyToAddress(privateKey.PublicKey)
		nodeId := discover.PubkeyID(&privateKey.PublicKey)
		canTmp := &staking.Candidate{
			CandidateBase: &staking.CandidateBase{
				NodeId:         nodeId,
				BenefitAddress: addr,
				Description:    staking.Description{},
			},
			CandidateMutable: &staking.CandidateMutable{},
		}
		// Store Candidate Base info
		canBaseKey := staking.CanBaseKeyByAddr(addr)
		if val, err := rlp.EncodeToBytes(canTmp.CandidateBase); nil != err {
			return nil, err
		} else {

			if err := sndb.PutBaseDB(canBaseKey, val); nil != err {
				return nil, err
			}
		}

		// Store Candidate Mutable info
		canMutableKey := staking.CanMutableKeyByAddr(addr)
		if val, err := rlp.EncodeToBytes(canTmp.CandidateMutable); nil != err {
			return nil, err
		} else {

			if err := sndb.PutBaseDB(canMutableKey, val); nil != err {
				return nil, err
			}
		}
		if i < int(xcom.MaxValidators()) {
			v := &staking.Validator{
				NodeAddress:     addr,
				NodeId:          canTmp.NodeId,
				BlsPubKey:       canTmp.BlsPubKey,
				ProgramVersion:  canTmp.ProgramVersion,
				Shares:          canTmp.Shares,
				StakingBlockNum: canTmp.StakingBlockNum,
				StakingTxIndex:  canTmp.StakingTxIndex,
				ValidatorTerm:   0,
			}
			validatorQueue[i] = v
		}
	}
	verifierIndex := &staking.ValArrIndex{
		Start: epochStart,
		End:   epochEnd,
	}
	epochIndexArr := make(staking.ValArrIndexQueue, 0)
	epochIndexArr = append(epochIndexArr, verifierIndex)
	// current epoch start and end indexs
	epochIndex, err := rlp.EncodeToBytes(epochIndexArr)
	if nil != err {
		return nil, err
	}
	if err := sndb.PutBaseDB(staking.GetEpochIndexKey(), epochIndex); nil != err {
		return nil, err
	}
	epochArr, err := rlp.EncodeToBytes(validatorQueue)
	if nil != err {
		return nil, err
	}
	// Store Epoch validators
	if err := sndb.PutBaseDB(staking.GetEpochValArrKey(verifierIndex.Start, verifierIndex.End), epochArr); nil != err {
		return nil, err
	}
	return validatorQueue, nil
}

func TestRewardPlugin(t *testing.T) {
	var plugin = RewardMgrInstance()
	StakingInstance()
	chain := mock.NewChain()
	mockDB := chain.StateDB
	snapshotdb.SetDBBlockChain(chain)

	defaultEc := *xcom.GetEc(xcom.DefaultTestNet)
	defer func() {
		xcom.ResetEconomicDefaultConfig(&defaultEc)
		snapshotdb.Instance().Clear()
	}()

	ec := xcom.GetEc(xcom.DefaultTestNet)
	ec.Common.AdditionalCycleTime = 1
	ec.Common.MaxEpochMinutes = 1
	ec.Common.MaxConsensusVals = 1
	ec.Common.NodeBlockTimeWindow = 1
	ec.Common.PerRoundBlocks = 1
	millisecond = 3
	minutes = 400

	t.Run("CalcEpochReward", func(t *testing.T) {
		//log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

		accounts := make(map[common.Address]*big.Int)

		yearBalance := big.NewInt(1e18)
		SetYearEndCumulativeIssue(mockDB, 0, yearBalance)
		SetYearEndBalance(mockDB, 0, yearBalance)
		mockDB.AddBalance(vm.RewardManagerPoolAddr, yearBalance)

		validatorQueueList, err := buildTestStakingData(1, xutil.CalcBlocksEachEpoch())
		if nil != err {
			t.Fatalf("buildTestStakingData fail: %v", err)
		}
		rand.Seed(int64(time.Now().Nanosecond()))
		var packageReward *big.Int
		var stakingReward *big.Int
		for i := 0; i < int(xutil.CalcBlocksEachEpoch()*5); i++ {
			parentHeader := chain.CurrentHeader()
			chain.AddBlock()
			currentHeader := chain.CurrentHeader()

			if err := snapshotdb.Instance().NewBlock(currentHeader.Number, parentHeader.Hash(), common.ZeroHash); nil != err {
				t.Fatal(err)
			}
			if err := plugin.EndBlock(common.ZeroHash, currentHeader, mockDB); nil != err {
				t.Fatalf("call endBlock fail, err：%v", err)
			}
			if currentHeader.Number.Uint64() < xutil.CalcBlocksEachEpoch() {
				time.Sleep(time.Duration(int(time.Millisecond) * 3))
			}

			if packageReward == nil {
				packageReward, err = LoadNewBlockReward(common.ZeroHash, plugin.db)
				if nil != err {
					t.Fatalf("call LoadNewBlockReward fail: %v", err)
				}
				stakingReward, err = LoadStakingReward(common.ZeroHash, plugin.db)
				if nil != err {
					t.Fatalf("call LoadStakingReward fail: %v", err)
				}
			}
			balance, ok := accounts[currentHeader.Coinbase]
			if !ok {
				balance = new(big.Int)
				accounts[currentHeader.Coinbase] = balance
			}
			balance.Add(balance, packageReward)
			if err := snapshotdb.Instance().Flush(currentHeader.Hash(), currentHeader.Number); nil != err {
				t.Fatal(err)
			}
			if err := snapshotdb.Instance().Commit(currentHeader.Hash()); nil != err {
				t.Fatal(err)
			}

			assert.Equal(t, accounts[currentHeader.Coinbase], mockDB.GetBalance(currentHeader.Coinbase))

			if xutil.IsEndOfEpoch(currentHeader.Number.Uint64()) {
				everyValidatorReward := new(big.Int).Div(stakingReward, big.NewInt(int64(len(validatorQueueList))))
				for _, value := range validatorQueueList {
					balance := accounts[value.NodeAddress]
					if balance == nil {
						balance = new(big.Int)
						accounts[value.NodeAddress] = balance
					}
					balance.Add(balance, everyValidatorReward)
					assert.Equal(t, balance, mockDB.GetBalance(value.NodeAddress))
				}

				validatorQueueList, err = buildTestStakingData(currentHeader.Number.Uint64()+1, currentHeader.Number.Uint64()+xutil.CalcBlocksEachEpoch())
				if nil != err {
					t.Fatalf("buildTestStakingData fail: %v", err)
				}

				packageReward, err = LoadNewBlockReward(common.ZeroHash, plugin.db)
				if nil != err {
					t.Fatalf("call LoadNewBlockReward fail: %v", err)
				}
				stakingReward, err = LoadStakingReward(common.ZeroHash, plugin.db)
				if nil != err {
					t.Fatalf("call LoadStakingReward fail: %v", err)
				}
			}
		}
	})
}

func TestIncreaseIssuance(t *testing.T) {
	var plugin = RewardMgrInstance()

	mockDB := buildStateDB(t)

	thisYear, lastYear := uint32(2), uint32(1)

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
}
