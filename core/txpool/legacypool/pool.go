// Copyright 2023 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package legacypool

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/core/txpool"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

// NewTxPool creates the default PlatON transaction pool backed by the legacy pool.
func NewTxPool(config txpool.Config, chainconfig *params.ChainConfig, chain *txpool.TxPoolBlockChain) *txpool.TxPool {
	legacy := New(config, chainconfig, chain)
	pool, err := txpool.New(new(big.Int).SetUint64(config.PriceLimit), chain, []txpool.SubPool{legacy})
	if err != nil {
		panic(err)
	}
	return pool
}
