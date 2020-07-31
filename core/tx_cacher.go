// Copyright 2018 The go-ethereum Authors
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

package core

import (
	"runtime"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

// senderCacher is a concurrent transaction sender recoverer anc cacher.
var SenderCacher = NewTxSenderCacher(runtime.NumCPU())

// txSenderCacherRequest is a request for recovering transaction senders with a
// specific signature scheme and caching it into the transactions themselves.
//
// The inc field defines the number of transactions to skip after each recovery,
// which is used to feed the same underlying input array to different threads but
// ensure they process the early transactions fast.
type txSenderCacherRequest struct {
	signer types.Signer
	txs    []*types.Transaction
	inc    int
	doneCh chan int
	starts int
}

// txSenderCacher is a helper structure to concurrently ecrecover transaction
// senders from digital signatures on background threads.
type txSenderCacher struct {
	threads int
	tasks   chan *txSenderCacherRequest
	txPool  *TxPool
}

//todoewTxSenderCacher creates a new transaction sender background cacher and starts
// as many processing goroutines as allowed by the GOMAXPROCS on construction.
func NewTxSenderCacher(threads int) *txSenderCacher {
	cacher := &txSenderCacher{
		tasks:   make(chan *txSenderCacherRequest, threads),
		threads: threads,
	}
	for i := 0; i < threads; i++ {
		go cacher.cache()
	}
	return cacher
}

// if set txpool ,will find from txpool first,if txpool have the tx,will not cal from any more
func (cacher *txSenderCacher) SetTxPool(txPool *TxPool) {
	cacher.txPool = txPool
}

// cache is an infinite loop, caching transaction senders from various forms of
// data structures.
func (cacher *txSenderCacher) cache() {
	for task := range cacher.tasks {
		txCal := 0
		for i := task.starts; i < len(task.txs); i += task.inc {
			types.Sender(task.signer, task.txs[i])
			txCal++
		}
		if task.doneCh != nil {
			task.doneCh <- txCal
		}
	}
}

// recover recovers the senders from a batch of transactions and caches them
// back into the same data structures. There is no validation being done, nor
// any reaction to invalid signatures. That is up to calling code later.
func (cacher *txSenderCacher) recover(signer types.Signer, txs []*types.Transaction) {
	// If there's nothing to recover, abort
	if len(txs) == 0 {
		return
	}
	// Ensure we have meaningful task sizes and schedule the recoveries
	tasks := cacher.threads
	if len(txs) < tasks*4 {
		tasks = (len(txs) + 3) / 4
	}
	for i := 0; i < tasks; i++ {
		cacher.tasks <- &txSenderCacherRequest{
			signer: signer,
			txs:    txs,
			inc:    tasks,
			starts: i,
		}
	}
}

// recoverFromBlocks recovers the senders from a batch of blocks and caches them
// back into the same data structures. There is no validation being done, nor
// any reaction to invalid signatures. That is up to calling code later.
/*func (cacher *txSenderCacher) recoverFromBlocks(signer types.Signer, blocks []*types.Block) {
	count := 0
	for _, block := range blocks {
		count += len(block.Transactions())
	}
	txs := make([]*types.Transaction, 0, count)
	if cacher.txPool != nil {
		for _, block := range blocks {
			for _, tx := range block.Transactions() {
				if txInPool := cacher.txPool.Get(tx.Hash()); txInPool != nil {
					tx = txInPool
				} else {
					txs = append(txs, tx)
				}
			}
		}
	} else {
		for _, block := range blocks {
			for _, tx := range block.Transactions() {
				txs = append(txs, tx)
			}
		}
	}

	if len(txs) > 0 {
		cacher.recover(signer, txs)
	}
}*/

/*func (cacher *txSenderCacher) RecoverTxsFromPool(signer types.Signer, txs []*types.Transaction) chan struct{} {
	// Ensure we have meaningful task sizes and schedule the recoveries
	tasks := cacher.threads
	if len(txs) < tasks*4 {
		tasks = (len(txs) + 3) / 4
	}

	CalTx := make(chan struct{}, tasks)
	for i := 0; i < tasks; i++ {
		cacher.tasks <- &txSenderCacherRequest{
			signer: signer,
			txs:    txs,
			inc:    tasks,
			done: CalTx,
			starts: i,
		}
	}
	return CalTx
}*/

// recoverFromBlock recovers the senders from  block and caches them
// back into the same data structures. There is no validation being done, nor
// any reaction to invalid signatures. That is up to calling code later.
func (cacher *txSenderCacher) RecoverFromBlock(signer types.Signer, block *types.Block) {
	count := len(block.Transactions())
	if count == 0 {
		return
	}
	txs := make([]*types.Transaction, 0, count)

	if cacher.txPool != nil && cacher.txPool.count() >= 200 {
		for i, tx := range block.Transactions() {
			if txInPool := cacher.txPool.Get(tx.Hash()); txInPool != nil {
				block.Transactions()[i].CacheFromAddr(signer, txInPool.FromAddr(signer))
			} else {
				txs = append(txs, block.Transactions()[i])
			}
		}
	} else {
		txs = block.Transactions()
	}
	if len(txs) == 0 {
		return
	}
	// Ensure we have meaningful task sizes and schedule the recoveries
	tasks := cacher.threads
	if len(txs) < tasks*4 {
		tasks = (len(txs) + 3) / 4
	}
	log.Trace("Start recover tx FromBlock", "number", block.Number(), "txs", len(txs), "tasks", tasks)
	block.CalTxFromCH = make(chan int, tasks)
	for i := 0; i < tasks; i++ {
		cacher.tasks <- &txSenderCacherRequest{
			signer: signer,
			txs:    txs,
			inc:    tasks,
			doneCh: block.CalTxFromCH,
			starts: i,
		}
	}
}
