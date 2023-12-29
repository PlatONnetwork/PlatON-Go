// Copyright 2021 The PlatON Network Authors
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
	"encoding/json"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/p2p/enode"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func NewDelegateRewardPer(epoch uint64, totalReward, totalDelegate *big.Int) *DelegateRewardPer {
	return &DelegateRewardPer{
		Left:     totalDelegate,
		Epoch:    epoch,
		Delegate: totalDelegate,
		Reward:   totalReward,
	}
}

//todo：这英文每看到啊？
type DelegateRewardPer struct {
	//节点在该周期的剩余委托金额，账户每次领取收益都会扣除对应委托金额，当该值为0时该记录删除
	Left  *big.Int
	Epoch uint64
	//节点在该周期的总委托金额
	Delegate *big.Int
	//节点在该周期的总委托收益
	Reward *big.Int
}

func (d *DelegateRewardPer) CalDelegateReward(delegate *big.Int) *big.Int {
	tmp := new(big.Int).Mul(delegate, d.Reward)
	return new(big.Int).Div(tmp, d.Delegate)
}

func NewDelegateRewardPerList() *DelegateRewardPerList {
	del := new(DelegateRewardPerList)
	del.Pers = make([]*DelegateRewardPer, 0)
	return del
}

type DelegateRewardPerList struct {
	Pers    []*DelegateRewardPer
	changed bool
}

func (d DelegateRewardPerList) String() string {
	v, err := json.Marshal(d)
	if err != nil {
		panic(err)
	}
	return string(v)
}

func (d *DelegateRewardPerList) AppendDelegateRewardPer(per *DelegateRewardPer) {
	d.Pers = append(d.Pers, per)
}

func (d *DelegateRewardPerList) DecreaseTotalAmount(receipt []DelegateRewardReceipt) int {
	var indexOfReceipt int
	for indexOfList := 0; indexOfList < len(d.Pers) && indexOfReceipt < len(receipt); {
		if d.Pers[indexOfList].Epoch < receipt[indexOfReceipt].Epoch {
			indexOfList++
		} else if d.Pers[indexOfList].Epoch > receipt[indexOfReceipt].Epoch {
			indexOfReceipt++
		} else {
			d.Pers[indexOfList].Left.Sub(d.Pers[indexOfList].Left, receipt[indexOfReceipt].Delegate)
			if d.Pers[indexOfList].Left.Cmp(common.Big0) <= 0 {
				d.Pers = append(d.Pers[:indexOfList], d.Pers[indexOfList+1:]...)
			} else {
				indexOfList++
			}
			indexOfReceipt++
			d.changed = true
		}
	}
	return indexOfReceipt
}

func (d *DelegateRewardPerList) ShouldDel() bool {
	if len(d.Pers) == 0 {
		return true
	}
	return false
}

func (d *DelegateRewardPerList) IsChange() bool {
	return d.changed
}

type NodeDelegateReward struct {
	NodeID     enode.IDv0 `json:"nodeID"`
	StakingNum uint64     `json:"stakingNum"`
	Reward     *big.Int   `json:"reward" rlp:"nil"`
}

type NodeDelegateRewardPresenter struct {
	NodeID     enode.IDv0   `json:"nodeID" `
	Reward     *hexutil.Big `json:"reward" `
	StakingNum uint64       `json:"stakingNum"`
}

type DelegateRewardReceipt struct {
	//this is the account  total effective delegate amount with the node  on this epoch
	Delegate *big.Int
	Epoch    uint64
}
