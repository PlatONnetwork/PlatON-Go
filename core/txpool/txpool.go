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
	"errors"
	"fmt"
	"math/big"
	"sync"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/metrics"
)

// TxStatus is the current status of a transaction as seen by the pool.
type TxStatus uint

const (
	TxStatusUnknown TxStatus = iota
	TxStatusQueued
	TxStatusPending
	TxStatusIncluded
)

var reservationsGaugeName = "txpool/reservations"

// BlockChain defines the minimal chain interface used by the coordinator loop.
type BlockChain interface {
	CurrentBlock() *types.Header
	SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription
}

// legacySubPool extends SubPool with PlatON-specific legacy pool APIs.
type legacySubPool interface {
	SubPool
	PendingPlain(enforceTips, limited bool) map[common.Address]types.Transactions
	ResetBlock(newBlock *types.Block)
	ForkedReset(newHeader *types.Header, rollback []*types.Block)
	Count() int
	GasPrice() *big.Int
}

// TxPool is an aggregator for various transaction specific pools.
type TxPool struct {
	subpools []SubPool
	legacy   legacySubPool

	reservations map[common.Address]SubPool
	reserveLock  sync.Mutex

	subs event.SubscriptionScope
	quit chan chan error
}

// New creates a new transaction pool from subpools.
func New(gasTip *big.Int, chain BlockChain, subpools []SubPool) (*TxPool, error) {
	head := chain.CurrentBlock()

	pool := &TxPool{
		subpools:     subpools,
		reservations: make(map[common.Address]SubPool),
		quit:         make(chan chan error),
	}
	for _, subpool := range subpools {
		if lp, ok := subpool.(legacySubPool); ok {
			pool.legacy = lp
			break
		}
	}
	for i, subpool := range subpools {
		if err := subpool.Init(gasTip, head, pool.reserver(i, subpool)); err != nil {
			for j := i - 1; j >= 0; j-- {
				subpools[j].Close()
			}
			return nil, err
		}
	}
	go pool.loop(head, chain)
	return pool, nil
}

func (p *TxPool) reserver(id int, subpool SubPool) AddressReserver {
	return func(addr common.Address, reserve bool) error {
		p.reserveLock.Lock()
		defer p.reserveLock.Unlock()

		owner, exists := p.reservations[addr]
		if reserve {
			if exists {
				if owner == subpool {
					log.Error("pool attempted to reserve already-owned address", "address", addr)
					return nil
				}
				return errors.New("address already reserved")
			}
			p.reservations[addr] = subpool
			if metrics.Enabled {
				m := fmt.Sprintf("%s/%d", reservationsGaugeName, id)
				metrics.GetOrRegisterGauge(m, nil).Inc(1)
			}
			return nil
		}
		if !exists {
			log.Error("pool attempted to unreserve non-reserved address", "address", addr)
			return errors.New("address not reserved")
		}
		if subpool != owner {
			log.Error("pool attempted to unreserve non-owned address", "address", addr)
			return errors.New("address not owned")
		}
		delete(p.reservations, addr)
		if metrics.Enabled {
			m := fmt.Sprintf("%s/%d", reservationsGaugeName, id)
			metrics.GetOrRegisterGauge(m, nil).Dec(1)
		}
		return nil
	}
}

func (p *TxPool) loop(head *types.Header, chain BlockChain) {
	var (
		newHeadCh  = make(chan core.ChainHeadEvent)
		newHeadSub = chain.SubscribeChainHeadEvent(newHeadCh)
	)
	defer newHeadSub.Unsubscribe()

	var (
		oldHead = head
		newHead = oldHead
	)
	var (
		resetBusy = make(chan struct{}, 1)
		resetDone = make(chan *types.Header)
	)
	var errc chan error
	for errc == nil {
		if newHead != oldHead {
			select {
			case resetBusy <- struct{}{}:
				go func(oldHead, newHead *types.Header) {
					for _, subpool := range p.subpools {
						subpool.Reset(oldHead, newHead)
					}
					resetDone <- newHead
				}(oldHead, newHead)
			default:
			}
		}
		select {
		case event := <-newHeadCh:
			newHead = event.Block.Header()
		case head := <-resetDone:
			oldHead = head
			<-resetBusy
		case errc = <-p.quit:
		}
	}
	errc <- nil
}

// Close terminates the transaction pool and all its subpools.
func (p *TxPool) Close() error { return p.Stop() }

// Stop terminates the transaction pool (PlatON compatibility alias).
func (p *TxPool) Stop() error {
	var errs []error
	errc := make(chan error)
	p.quit <- errc
	if err := <-errc; err != nil {
		errs = append(errs, err)
	}
	for _, subpool := range p.subpools {
		if err := subpool.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("subpool close errors: %v", errs)
	}
	return nil
}

func (p *TxPool) SetGasTip(tip *big.Int) {
	for _, subpool := range p.subpools {
		subpool.SetGasTip(tip)
	}
}

func (p *TxPool) GasPrice() *big.Int {
	if p.legacy != nil {
		return p.legacy.GasPrice()
	}
	return new(big.Int)
}

func (p *TxPool) SetGasPrice(price *big.Int) { p.SetGasTip(price) }

func (p *TxPool) Has(hash common.Hash) bool {
	for _, subpool := range p.subpools {
		if subpool.Has(hash) {
			return true
		}
	}
	return false
}

func (p *TxPool) Get(hash common.Hash) *types.Transaction {
	for _, subpool := range p.subpools {
		if tx := subpool.Get(hash); tx != nil {
			return tx.Tx
		}
	}
	return nil
}

func (p *TxPool) Add(txs []*Transaction, local bool, sync bool) []error {
	txsets := make([][]*Transaction, len(p.subpools))
	splits := make([]int, len(txs))
	for i, tx := range txs {
		splits[i] = -1
		for j, subpool := range p.subpools {
			if subpool.Filter(tx.Tx) {
				txsets[j] = append(txsets[j], tx)
				splits[i] = j
				break
			}
		}
	}
	errsets := make([][]error, len(p.subpools))
	for i := range p.subpools {
		errsets[i] = p.subpools[i].Add(txsets[i], local, sync)
	}
	errs := make([]error, len(txs))
	for i, split := range splits {
		if split == -1 {
			errs[i] = core.ErrTxTypeNotSupported
			continue
		}
		errs[i] = errsets[split][0]
		errsets[split] = errsets[split][1:]
	}
	return errs
}

func (p *TxPool) AddRemotes(txs []*types.Transaction) []error {
	return p.addPlain(txs, false, false)
}

func (p *TxPool) AddRemotesSync(txs []*types.Transaction) []error {
	return p.addPlain(txs, false, true)
}

func (p *TxPool) AddRemote(tx *types.Transaction) error {
	return p.addPlain([]*types.Transaction{tx}, false, false)[0]
}

func (p *TxPool) AddLocals(txs []*types.Transaction) []error {
	return p.addPlain(txs, true, true)
}

func (p *TxPool) AddLocal(tx *types.Transaction) error {
	return p.addPlain([]*types.Transaction{tx}, true, true)[0]
}

func (p *TxPool) addPlain(txs []*types.Transaction, local, sync bool) []error {
	wrapped := make([]*Transaction, len(txs))
	for i, tx := range txs {
		wrapped[i] = &Transaction{Tx: tx}
	}
	return p.Add(wrapped, local, sync)
}

func (p *TxPool) Pending(enforceTips, limited bool) map[common.Address]types.Transactions {
	if p.legacy != nil {
		return p.legacy.PendingPlain(enforceTips, limited)
	}
	return nil
}

// PendingPlain is an alias for Pending for callers that use the explicit name.
func (p *TxPool) PendingPlain(enforceTips, limited bool) map[common.Address]types.Transactions {
	return p.Pending(enforceTips, limited)
}

func (p *TxPool) SubscribeNewTxsEvent(ch chan<- core.NewTxsEvent) event.Subscription {
	subs := make([]event.Subscription, len(p.subpools))
	for i, subpool := range p.subpools {
		subs[i] = subpool.SubscribeTransactions(ch)
	}
	return p.subs.Track(event.JoinSubscriptions(subs...))
}

func (p *TxPool) Nonce(addr common.Address) uint64 {
	var nonce uint64
	for _, subpool := range p.subpools {
		if next := subpool.Nonce(addr); nonce < next {
			nonce = next
		}
	}
	return nonce
}

func (p *TxPool) Stats() (int, int) {
	var runnable, blocked int
	for _, subpool := range p.subpools {
		run, block := subpool.Stats()
		runnable += run
		blocked += block
	}
	return runnable, blocked
}

func (p *TxPool) Content() (map[common.Address][]*types.Transaction, map[common.Address][]*types.Transaction) {
	runnable := make(map[common.Address][]*types.Transaction)
	blocked := make(map[common.Address][]*types.Transaction)
	for _, subpool := range p.subpools {
		run, block := subpool.Content()
		for addr, txs := range run {
			runnable[addr] = txs
		}
		for addr, txs := range block {
			blocked[addr] = txs
		}
	}
	return runnable, blocked
}

func (p *TxPool) ContentFrom(addr common.Address) ([]*types.Transaction, []*types.Transaction) {
	for _, subpool := range p.subpools {
		run, block := subpool.ContentFrom(addr)
		if len(run) != 0 || len(block) != 0 {
			return run, block
		}
	}
	return []*types.Transaction{}, []*types.Transaction{}
}

func (p *TxPool) Locals() []common.Address {
	locals := make(map[common.Address]struct{})
	for _, subpool := range p.subpools {
		for _, local := range subpool.Locals() {
			locals[local] = struct{}{}
		}
	}
	flat := make([]common.Address, 0, len(locals))
	for local := range locals {
		flat = append(flat, local)
	}
	return flat
}

func (p *TxPool) Status(hash common.Hash) TxStatus {
	for _, subpool := range p.subpools {
		if status := subpool.Status(hash); status != TxStatusUnknown {
			return status
		}
	}
	return TxStatusUnknown
}

func (p *TxPool) Count() int {
	if p.legacy != nil {
		return p.legacy.Count()
	}
	return 0
}

// Reset is used by CBFT to reset the pool on new block.
func (p *TxPool) Reset(newBlock *types.Block) {
	if p.legacy != nil {
		p.legacy.ResetBlock(newBlock)
	}
}

// ForkedReset handles rollback reorgs in CBFT.
func (p *TxPool) ForkedReset(newHeader *types.Header, rollback []*types.Block) {
	if p.legacy != nil {
		p.legacy.ForkedReset(newHeader, rollback)
	}
}
