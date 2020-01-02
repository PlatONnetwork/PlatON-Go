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

package reward

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

var (
	HistoryIncreasePrefix         = []byte("RewardHistory")
	LastYearEndBalancePrefix      = []byte("RewardBalance")
	YearStartBlockNumberKey       = []byte("YearStartBlockNumberKey")
	YearStartTimeKey              = []byte("YearStartTimeKey")
	RemainingRewardKey            = []byte("RemainingRewardKey")
	NewBlockRewardKey             = []byte("NewBlockRewardKey")
	StakingRewardKey              = []byte("StakingRewardKey")
	ChainYearNumberKey            = []byte("ChainYearNumberKey")
	delegateRewardPerKey          = []byte("DelegateRewardPerKey")
	currentEpochDelegateRewardKey = []byte("currentEpochDelegateRewardKey")
)

// GetHistoryIncreaseKey used for search the balance of reward pool at last year
func GetHistoryIncreaseKey(year uint32) []byte {
	return append(HistoryIncreasePrefix, common.Uint32ToBytes(year)...)
}

//
func HistoryBalancePrefix(year uint32) []byte {
	return append(LastYearEndBalancePrefix, common.Uint32ToBytes(year)...)
}

func DelegateRewardPerKey(nodeID discover.NodeID, epoch uint32) []byte {
	return nil
}

func CurrentEpochDelegateRewardKey(nodeID discover.NodeID) []byte {
	return nil
}

type DelegateRewardPer struct {
	Epoch     uint64 `rlp:"nil"`
	Amount    *big.Int
	NodeCount uint
}

type DelegateRewardPerList []DelegateRewardPer

//this is use for NodeCounts--
func (d *DelegateRewardPerList) HandleNodeCount(epoch uint) {

}

func (d *DelegateRewardPerList) SetDelegateRewardPer(epoch uint, amount *big.Int) {

}

type NodeDelegateReward struct {
	NodeID discover.NodeID
	Reward *big.Int
}
