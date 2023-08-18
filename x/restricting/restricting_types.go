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

package restricting

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
)

// for genesis and plugin test
// 每个账户，作为锁仓计划的资金释放目标对象，都可以有且只有一个这样的对象，记录当前的锁仓计划状态
type RestrictingInfo struct {
	//欠释放金额，到了结算周期需要释放却因为质押而无法释放的金额
	NeedRelease *big.Int
	// 用于质押和委托的金额
	AdvanceAmount *big.Int
	// 可用的锁仓金额 = 可用锁仓金额 - 已释放的（需要释放的） - 被惩罚的(用于质押而被处罚)
	CachePlanAmount *big.Int

	ReleaseList []uint64 // ReleaseList representation which epoch will release restricting
}

func (r *RestrictingInfo) RemoveEpoch(epoch uint64) {
	for i, target := range r.ReleaseList {
		if target == epoch {
			r.ReleaseList = append(r.ReleaseList[:i], r.ReleaseList[i+1:]...)
			break
		}
	}
}

// for contract, plugin test, byte util
type RestrictingPlan struct {
	Epoch  uint64   `json:"epoch"`  // epoch representation of the released epoch at the target blockNumber
	Amount *big.Int `json:"amount"` // amount representation of the released amount
}

// for plugin test
type ReleaseAmountInfo struct {
	Height uint64       `json:"blockNumber"` // blockNumber representation of the block number at the released epoch
	Amount *hexutil.Big `json:"amount"`      // amount representation of the released amount
}

// for plugin test
type Result struct {
	Balance *hexutil.Big        `json:"balance"`
	Debt    *hexutil.Big        `json:"debt"`
	Entry   []ReleaseAmountInfo `json:"plans"`
	Pledge  *hexutil.Big        `json:"Pledge"`
}
