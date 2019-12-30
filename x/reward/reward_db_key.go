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

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

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
	delegateRewardTotalKey   = []byte("DelegateRewardTotalKey")
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
	index := epoch / DelegateRewardPerLength
	add, err := xutil.NodeId2Addr(nodeID)
	if err != nil {
		panic(err)
	}
	keyAdd := append(delegateRewardPerKey, add.Bytes()...)
	keyAdd = append(delegateRewardPerKey, common.Uint64ToBytes(stakingNum)...)
	return append(keyAdd, common.Uint64ToBytes(index)...)
}

func DelegateRewardPerKeys(nodeID discover.NodeID, stakingNum, fromEpoch, toEpoch uint64) [][]byte {
	indexFrom := fromEpoch / DelegateRewardPerLength
	indexTo := toEpoch / DelegateRewardPerLength
	add, err := xutil.NodeId2Addr(nodeID)
	if err != nil {
		panic(err)
	}
	delegateRewardPerPrefix := append(add[:], common.Uint64ToBytes(stakingNum)...)
	keys := make([][]byte, 0)
	for i := indexFrom; i <= indexTo; i++ {
		delegateRewardPerKey := append(delegateRewardPerPrefix[:], common.Uint64ToBytes(i)...)
		keys = append(keys, delegateRewardPerKey)
	}
	return keys
}

func DelegateRewardTotalKey(nodeID discover.NodeID, stakingNum uint64) []byte {
	add, err := xutil.NodeId2Addr(nodeID)
	if err != nil {
		panic(err)
	}
	keyAdd := append(delegateRewardTotalKey, add.Bytes()...)
	keyAdd = append(delegateRewardTotalKey, common.Uint64ToBytes(stakingNum)...)
	return keyAdd
}

func NewDelegateRewardPer(epoch uint64, per, total *big.Int) *DelegateRewardPer {
	return &DelegateRewardPer{
		TotalAmount: total,
		Amount:      per,
		Epoch:       epoch,
	}
}

type DelegateRewardPer struct {
	TotalAmount *big.Int
	Epoch       uint64 `rlp:"nil"`
	Amount      *big.Int
}

type DelegateRewardPerList struct {
	Pers   map[uint64]*DelegateRewardPer
	Epochs []uint64
}

func NewDelegateRewardPerList() *DelegateRewardPerList {
	del := new(DelegateRewardPerList)
	del.Pers = make(map[uint64]*DelegateRewardPer)
	del.Epochs = make([]uint64, 0)
	return del
}

func (d *DelegateRewardPerList) AppendDelegateRewardPer(per *DelegateRewardPer) {
	//index := epoch % DelegateRewardPerLength
	d.Pers[per.Epoch] = per
	d.Epochs = append(d.Epochs, per.Epoch)
}

func (d *DelegateRewardPerList) DecreaseTotalAmount(epoch uint64, amount *big.Int) {
	per, ok := d.Pers[epoch]
	if !ok {
		return
	}
	per.TotalAmount.Sub(per.TotalAmount, amount)
	if per.TotalAmount.Cmp(common.Big0) <= 0 {
		delete(d.Pers, epoch)
		delIndex := 0
		for i, v := range d.Epochs {
			if v == epoch {
				delIndex = i
				break
			}
		}
		d.Epochs = append(d.Epochs[:delIndex], d.Epochs[delIndex+1:]...)
	}
}

func (d *DelegateRewardPerList) ShouldDel() bool {
	if len(d.Epochs) == 0 {
		return true
	}
	return false
}

type NodeDelegateReward struct {
	NodeID     discover.NodeID `json:"nodeID" rlp:"nodeID"`
	Reward     *big.Int        `json:"reward" rlp:"reward"`
	StakingNum uint64          `json:"stakingNum" rlp:"stakingNum"`
}

type NodeDelegateRewardPresenter struct {
	NodeID     discover.NodeID `json:"nodeID" rlp:"nodeID"`
	Reward     *hexutil.Big    `json:"reward" rlp:"reward"`
	StakingNum uint64          `json:"stakingNum" rlp:"stakingNum"`
}

type DelegateRewardReceive struct {
	Reward *big.Int
	Epoch  uint64
}
