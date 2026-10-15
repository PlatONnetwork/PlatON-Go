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

package txpool

import (
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto/kzg4844"
	"github.com/PlatONnetwork/PlatON-Go/event"
)

// Transaction is a helper struct to group together a canonical transaction with
// satellite data items that are needed by the pool but are not part of the chain.
type Transaction struct {
	Tx *types.Transaction

	BlobTxBlobs   []kzg4844.Blob
	BlobTxCommits []kzg4844.Commitment
	BlobTxProofs  []kzg4844.Proof
}

// LazyTransaction contains a small subset of the transaction properties that is
// enough for the miner and other APIs to handle large batches of transactions;
// and supports pulling up the entire transaction when really needed.
type LazyTransaction struct {
	Pool SubPool
	Hash common.Hash
	Tx   *Transaction

	Time      time.Time
	GasFeeCap *big.Int
	GasTipCap *big.Int
}

// Resolve retrieves the full transaction belonging to a lazy handle if it is still
// maintained by the transaction pool.
func (ltx *LazyTransaction) Resolve() *Transaction {
	if ltx.Tx == nil {
		ltx.Tx = ltx.Pool.Get(ltx.Hash)
	}
	return ltx.Tx
}

// AddressReserver is passed by the main transaction pool to subpools, so they
// may request (and relinquish) exclusive access to certain addresses.
type AddressReserver func(addr common.Address, reserve bool) error

// SubPool represents a specialized transaction pool that lives on its own (e.g.
// blob pool). Since independent of how many specialized pools we have, they do
// need to be updated in lockstep and assemble into one coherent view for block
// production, this interface defines the common methods that allow the primary
// transaction pool to manage the subpools.
type SubPool interface {
	Filter(tx *types.Transaction) bool

	Init(gasTip *big.Int, head *types.Header, reserve AddressReserver) error

	Close() error

	Reset(oldHead, newHead *types.Header)

	SetGasTip(tip *big.Int)

	Has(hash common.Hash) bool

	Get(hash common.Hash) *Transaction

	Add(txs []*Transaction, local bool, sync bool) []error

	Pending(enforceTips bool) map[common.Address][]*LazyTransaction

	SubscribeTransactions(ch chan<- core.NewTxsEvent) event.Subscription

	Nonce(addr common.Address) uint64

	Stats() (int, int)

	Content() (map[common.Address][]*types.Transaction, map[common.Address][]*types.Transaction)

	ContentFrom(addr common.Address) ([]*types.Transaction, []*types.Transaction)

	Locals() []common.Address

	Status(hash common.Hash) TxStatus
}
