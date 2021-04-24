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

package core

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

var (
	lastNumberKey = []byte("last-clean-number")

	minCleanTimeout = time.Minute

	cleanDistance uint64 = 1
)

type CleanupEvent struct{}

type CleanBatch struct {
	lock  sync.Mutex
	batch ethdb.Batch
}

func (cb *CleanBatch) Delete(key []byte) error {
	cb.lock.Lock()
	defer cb.lock.Unlock()
	return cb.batch.Delete(key)
}

func (cb *CleanBatch) ValueSize() int {
	cb.lock.Lock()
	defer cb.lock.Unlock()
	return cb.batch.ValueSize()
}

func (cb *CleanBatch) WriteAndRest() error {
	cb.lock.Lock()
	defer cb.lock.Unlock()
	if err := cb.batch.Write(); err != nil {
		cb.batch.Reset()
		return err
	}
	cb.batch.Reset()
	return nil
}

type Cleaner struct {
	stopped      common.AtomicBool
	cleaning     common.AtomicBool
	interval     uint64
	lastNumber   uint64
	cleanTimeout time.Duration
	gcMpt        bool

	wg        sync.WaitGroup
	exit      chan struct{}
	cleanFeed event.Feed
	scope     event.SubscriptionScope
	cleanCh   chan *CleanupEvent

	batch      CleanBatch
	lock       sync.RWMutex
	blockchain *BlockChain
}

func NewCleaner(blockchain *BlockChain, interval uint64, cleanTimeout time.Duration, gcMpt bool) *Cleaner {
	c := &Cleaner{
		interval:     interval,
		lastNumber:   0,
		cleanTimeout: cleanTimeout,
		gcMpt:        gcMpt,
		exit:         make(chan struct{}),
		cleanCh:      make(chan *CleanupEvent, 1),
		batch: CleanBatch{
			batch: blockchain.db.NewBatch(),
		},
		blockchain: blockchain,
	}

	if c.cleanTimeout < minCleanTimeout {
		c.cleanTimeout = minCleanTimeout
	}

	buf, err := c.blockchain.db.Get(lastNumberKey)
	if err == nil && len(buf) > 0 {
		lastNumber := common.BytesToUint64(buf)
		atomic.StoreUint64(&c.lastNumber, lastNumber)
	}

	c.scope.Track(c.cleanFeed.Subscribe(c.cleanCh))
	c.wg.Add(1)
	go c.loop()
	return c
}

func (c *Cleaner) Stop() {
	if c.stopped.IsSet() {
		return
	}

	c.scope.Close()
	close(c.exit)

	c.stopped.Set(true)
	c.wg.Wait()
}

func (c *Cleaner) Cleanup() {
	if c.cleaning.IsSet() {
		return
	}
	c.cleaning.Set(true)
	c.cleanFeed.Send(&CleanupEvent{})
}

func (c *Cleaner) NeedCleanup() bool {
	lastNumber := atomic.LoadUint64(&c.lastNumber)
	return c.blockchain.CurrentBlock().NumberU64()-lastNumber >= 2*c.interval && !c.cleaning.IsSet()
}

func (c *Cleaner) loop() {
	defer c.wg.Done()

	for {
		select {
		case <-c.cleanCh:
			c.cleanup()
		case <-c.exit:
			return
		}
	}
}

func (c *Cleaner) cleanup() {
	defer c.cleaning.Set(false)

	db := c.blockchain.db

	lastNumber := atomic.LoadUint64(&c.lastNumber)
	currentBlock := c.blockchain.CurrentBlock()
	if currentBlock.NumberU64()-lastNumber <= cleanDistance {
		return
	}

	var (
		receipts = 0
		txs      = 0
		keys     = 0
	)

	t := time.Now()
	log.Info("Start cleanup database", "interval", c.interval, "cleanTimeout", c.cleanTimeout, "gcMpt", c.gcMpt, "lastNumber", atomic.LoadUint64(&c.lastNumber), "number", currentBlock.NumberU64(), "hash", currentBlock.Hash())
	defer func() {
		log.Info("Finish cleanup database", "lastNumber", atomic.LoadUint64(&c.lastNumber), "receipts", receipts, "txs", txs, "keys", keys, "elapsed", time.Since(t))
	}()

	if currentBlock.NumberU64()-c.lastNumber >= 2*c.interval {
		number := lastNumber + 1
		for ; number <= currentBlock.NumberU64()-c.interval; number++ {
			block := c.blockchain.GetBlockByNumber(number)
			if block == nil {
				log.Error("Found bad header", "number", number)
				return
			}

			rawdb.DeleteReceipts(db, block.Hash(), block.NumberU64())

			//batch := c.blockchain.db.NewBatch()
			//for _, tx := range block.Transactions() {
			//	txs++
			//	rawdb.DeleteTxLookupEntry(batch, tx.Hash())
			//}
			//batch.Write()

			receipts++

			if time.Since(t) >= c.cleanTimeout || c.stopped.IsSet() {
				atomic.StoreUint64(&c.lastNumber, number)
				db.Put(lastNumberKey, common.Uint64ToBytes(number))
				return
			}
		}
		atomic.StoreUint64(&c.lastNumber, number-1)
		db.Put(lastNumberKey, common.Uint64ToBytes(number-1))
	}

}
