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
	"os"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/xcom"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/common/vm"

	"github.com/PlatONnetwork/PlatON-Go/x/staking"

	"github.com/PlatONnetwork/PlatON-Go/x/xutil"
)

type testBlockChain struct {
	currentBlock *types.Block
	items        map[uint64]*types.Block
	blockNumbers map[common.Hash]uint64
	accounts     map[common.Address]*big.Int
}

func (chain *testBlockChain) CurrentHeader() *types.Header {
	return chain.currentBlock.Header()
}
func (chain *testBlockChain) GetHeaderByHash(hash common.Hash) *types.Header {
	return chain.items[chain.blockNumbers[hash]].Header()
}
func (chain *testBlockChain) GetHeaderByNumber(number uint64) *types.Header {
	return chain.items[number].Header()
}

func (chain *testBlockChain) insertBlock() *types.Block {
	parentBlock := chain.currentBlock
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		panic(err)
	}
	addr := crypto.PubkeyToAddress(privateKey.PublicKey)
	header := &types.Header{
		ParentHash: parentBlock.Hash(),
		Number:     new(big.Int).Add(parentBlock.Number(), common.Big1),
		GasLimit:   1000000,
		Time:       new(big.Int).SetInt64(time.Now().UnixNano() / 1e6),
		Coinbase:   addr,
	}
	if _, ok := chain.accounts[header.Coinbase]; !ok {
		chain.accounts[header.Coinbase] = new(big.Int)
	}
	block := types.NewBlock(header, nil, nil)
	chain.currentBlock = block
	chain.items[header.Number.Uint64()] = block
	chain.blockNumbers[header.Hash()] = header.Number.Uint64()
	return block
}

func newTestBlockChain() *testBlockChain {
	header := &types.Header{
		ParentHash: common.ZeroHash,
		Number:     new(big.Int),
		GasLimit:   1000000,
		Time:       new(big.Int).SetInt64(time.Now().UnixNano() / 1e6),
	}
	block := types.NewBlock(header, nil, nil)
	chain := &testBlockChain{
		currentBlock: block,
		items:        make(map[uint64]*types.Block),
		blockNumbers: make(map[common.Hash]uint64),
		accounts:     make(map[common.Address]*big.Int),
	}
	chain.items[header.Number.Uint64()] = block
	chain.blockNumbers[header.Hash()] = header.Number.Uint64()
	return chain
}

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
	chain := newTestBlockChain()
	snapshotdb.SetDBBlockChain(chain)
	mockDB := buildStateDB(t)

	t.Run("CalcEpochReward", func(t *testing.T) {
		log.Root().SetHandler(log.CallerFileHandler(log.LvlFilterHandler(log.Lvl(4), log.StreamHandler(os.Stderr, log.TerminalFormat(true)))))

		yearBalance := big.NewInt(1e18)
		//rate := xcom.NewBlockRewardRate()
		//epochs := xutil.EpochsPerYear()
		//blocks := xutil.CalcBlocksEachYear()
		//thisYear, lastYear := uint32(2), uint32(1)
		//expectNewBlockReward := percentageCalculation(yearBalance, rate)
		SetYearEndBalance(mockDB, 0, yearBalance)
		mockDB.AddBalance(vm.RewardManagerPoolAddr, yearBalance)

		validatorQueueList, err := buildTestStakingData(1, xutil.CalcBlocksEachEpoch())
		if nil != err {
			t.Fatalf("buildTestStakingData fail: %v", err)
		}
		var packageReward *big.Int
		var stakingReward *big.Int
		for i := 0; i < int(xutil.CalcBlocksEachEpoch()*2); i++ {
			parentBlock := chain.currentBlock
			block := chain.insertBlock()

			if err := snapshotdb.Instance().NewBlock(block.Number(), parentBlock.Hash(), common.ZeroHash); nil != err {
				t.Fatal(err)
			}
			if err := plugin.EndBlock(common.ZeroHash, block.Header(), mockDB); nil != err {
				t.Fatalf("call endBlock fail, err：%v", err)
			}
			rand.Seed(int64(time.Now().Nanosecond()))
			time.Sleep(time.Duration(int(time.Millisecond) * rand.Intn(5)))

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
			if _, ok := chain.accounts[block.Coinbase()]; ok {
				balance := chain.accounts[block.Coinbase()]
				balance.Add(balance, packageReward)
			}
			if err := snapshotdb.Instance().Flush(block.Hash(), block.Number()); nil != err {
				t.Fatal(err)
			}
			if err := snapshotdb.Instance().Commit(block.Hash()); nil != err {
				t.Fatal(err)
			}

			assert.Equal(t, chain.accounts[block.Coinbase()], mockDB.GetBalance(block.Coinbase()))

			if xutil.IsEndOfEpoch(block.NumberU64()) {
				everyValidatorReward := new(big.Int).Div(stakingReward, big.NewInt(int64(len(validatorQueueList))))
				for _, value := range validatorQueueList {
					balance := chain.accounts[value.NodeAddress]
					if balance == nil {
						balance = new(big.Int)
						chain.accounts[value.NodeAddress] = balance
					}
					balance.Add(balance, everyValidatorReward)
					assert.Equal(t, balance, mockDB.GetBalance(value.NodeAddress))
				}

				validatorQueueList, err = buildTestStakingData(block.NumberU64()+1, block.NumberU64()+xutil.CalcBlocksEachEpoch())
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

		/*plugin.rewardStakingByValidatorList(mockDB, list, stakingReward)
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

		assert.Equal(t, new(big.Int).Sub(thisYearIssue, lastYearIssue), new(big.Int).Div(lastYearIssue, big.NewInt(IncreaseIssue)))*/

	})

}
