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

package restricting

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
)

// for genesis and plugin test
type RestrictingInfo struct {
	NeedRelease     *big.Int
	AdvanceAmount   *big.Int
	CachePlanAmount *big.Int
	ReleaseList     []uint64 // ReleaseList representation which epoch will release restricting
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
