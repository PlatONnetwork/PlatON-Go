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

package reward

import (
	"github.com/PlatONnetwork/PlatON-Go/x/xutil"

	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

const DelegateRewardPerLength = 1000

var (
	HistoryIncreasePrefix    = []byte("RewardHistory")
	LastYearEndBalancePrefix = []byte("RewardBalance")
	YearStartBlockNumberKey  = []byte("YearStartBlockNumberKey")
	YearStartTimeKey         = []byte("YearStartTimeKey")
	RemainingRewardKey       = []byte("RemainingRewardKey")
	NewBlockRewardKey        = []byte("NewBlockRewardKey")
	StakingRewardKey         = []byte("StakingRewardKey")
	ChainYearNumberKey       = []byte("ChainYearNumberKey")
	delegateRewardPerKey     = []byte("DelegateRewardPerKey")
)

// GetHistoryIncreaseKey used for search the balance of reward pool at last year
func GetHistoryIncreaseKey(year uint32) []byte {
	return append(HistoryIncreasePrefix, common.Uint32ToBytes(year)...)
}

//
func HistoryBalancePrefix(year uint32) []byte {
	return append(LastYearEndBalancePrefix, common.Uint32ToBytes(year)...)
}

func DelegateRewardPerKey(nodeID discover.NodeID, stakingNum, epoch uint64) []byte {
	index := uint32(epoch / DelegateRewardPerLength)
	add, err := xutil.NodeId2Addr(nodeID)
	if err != nil {
		panic(err)
	}
	perKeyLength := len(delegateRewardPerKey)
	lengthUint32, lengthUint64 := 4, 8
	keyAdd := make([]byte, perKeyLength+common.AddressLength+lengthUint64+lengthUint32)
	n := copy(keyAdd[:perKeyLength], delegateRewardPerKey)
	n += copy(keyAdd[n:n+common.AddressLength], add.Bytes())
	n += copy(keyAdd[n:n+lengthUint64], common.Uint64ToBytes(stakingNum))
	copy(keyAdd[n:n+lengthUint32], common.Uint32ToBytes(index))

	return keyAdd
}

func DelegateRewardPerKeys(nodeID discover.NodeID, stakingNum, fromEpoch, toEpoch uint64) [][]byte {
	indexFrom := uint32(fromEpoch / DelegateRewardPerLength)
	indexTo := uint32(toEpoch / DelegateRewardPerLength)
	add, err := xutil.NodeId2Addr(nodeID)
	if err != nil {
		panic(err)
	}
	perKeyLength := len(delegateRewardPerKey)
	lengthUint64 := 8

	delegateRewardPerPrefix := make([]byte, perKeyLength+common.AddressLength+lengthUint64)
	n := copy(delegateRewardPerPrefix[:perKeyLength], delegateRewardPerKey)
	n += copy(delegateRewardPerPrefix[n:n+common.AddressLength], add.Bytes())
	n += copy(delegateRewardPerPrefix[n:n+lengthUint64], common.Uint64ToBytes(stakingNum))

	keys := make([][]byte, 0)
	for i := indexFrom; i <= indexTo; i++ {
		delegateRewardPerKey := append(delegateRewardPerPrefix[:], common.Uint32ToBytes(i)...)
		keys = append(keys, delegateRewardPerKey)
	}
	return keys
}
