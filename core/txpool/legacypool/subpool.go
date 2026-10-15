package legacypool

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/txpool"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
)

// Filter implements txpool.SubPool.
func (pool *LegacyPool) Filter(tx *types.Transaction) bool {
	switch tx.Type() {
	case types.LegacyTxType, types.AccessListTxType, types.DynamicFeeTxType:
		return true
	default:
		return false
	}
}

// Init implements txpool.SubPool.
func (pool *LegacyPool) Init(gasTip *big.Int, head *types.Header, reserve txpool.AddressReserver) error {
	pool.reserve = reserve
	if gasTip != nil {
		pool.gasTip.Store(new(big.Int).Set(gasTip))
	}
	return nil
}

// Reset implements txpool.SubPool.
func (pool *LegacyPool) Reset(oldHead, newHead *types.Header) {
	wait := pool.requestReset(oldHead, newHead)
	<-wait
}

// SubscribeTransactions implements txpool.SubPool.
func (pool *LegacyPool) SubscribeTransactions(ch chan<- core.NewTxsEvent) event.Subscription {
	return pool.SubscribeNewTxsEvent(ch)
}

// Add implements txpool.SubPool.
func (pool *LegacyPool) Add(txs []*txpool.Transaction, local bool, sync bool) []error {
	unwrapped := make([]*types.Transaction, len(txs))
	for i, tx := range txs {
		unwrapped[i] = tx.Tx
	}
	return pool.addTxs(unwrapped, local, sync)
}

// Get implements txpool.SubPool.
func (pool *LegacyPool) Get(hash common.Hash) *txpool.Transaction {
	tx := pool.get(hash)
	if tx == nil {
		return nil
	}
	return &txpool.Transaction{Tx: tx}
}

// GetTransaction returns a plain transaction if present in the pool.
func (pool *LegacyPool) GetTransaction(hash common.Hash) *types.Transaction {
	return pool.get(hash)
}

// Pending implements txpool.SubPool.
func (pool *LegacyPool) Pending(enforceTips bool) map[common.Address][]*txpool.LazyTransaction {
	return pool.pendingLazy(enforceTips, false)
}

func (pool *LegacyPool) pendingLazy(enforceTips, limited bool) map[common.Address][]*txpool.LazyTransaction {
	pending := pool.PendingPlain(enforceTips, limited)
	lazy := make(map[common.Address][]*txpool.LazyTransaction, len(pending))
	for addr, txs := range pending {
		set := make([]*txpool.LazyTransaction, len(txs))
		for i, tx := range txs {
			set[i] = &txpool.LazyTransaction{
				Pool:      pool,
				Hash:      tx.Hash(),
				Tx:        &txpool.Transaction{Tx: tx},
				GasFeeCap: tx.GasFeeCap(),
				GasTipCap: tx.GasTipCap(),
			}
		}
		lazy[addr] = set
	}
	return lazy
}

// Status implements txpool.SubPool.
func (pool *LegacyPool) Status(hash common.Hash) txpool.TxStatus {
	statuses := pool.Statuses([]common.Hash{hash})
	if len(statuses) == 0 {
		return txpool.TxStatusUnknown
	}
	return statuses[0]
}
