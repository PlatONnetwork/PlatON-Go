package miner

import (
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

type ParallelTxsCommitter struct {
	worker *worker
}

func NewParallelTxsCommitter(w *worker) *ParallelTxsCommitter {
	return &ParallelTxsCommitter{
		worker: w,
	}
}

func (c *ParallelTxsCommitter) CommitTransactions(header *types.Header, txs *types.TransactionsByPriceAndNonce, interrupt *int32, timestamp int64, startTime, blockDeadline time.Time) (failed bool, isTimeout bool) {
	w := c.worker

	// Short circuit if current is nil
	timeout := false

	if w.current == nil {
		return true, timeout
	}

	if w.current.gasPool == nil {
		w.current.gasPool = new(core.GasPool).AddGas(w.current.header.GasLimit)
	}

	var coalescedLogs []*types.Log
	//var bftEngine = w.config.Cbft != nil

	var parallelTxs types.Transactions
	for {
		tx := txs.Peek()
		if tx == nil {
			break
		} else {
			parallelTxs = append(parallelTxs, tx)
			txs.Shift()
		}
	}

	ctx := core.NewPackBlockContext(w.current.state, header, common.Hash{}, w.current.gasPool, startTime, blockDeadline)
	ctx.SetTxList(parallelTxs)
	if err := core.GetExecutor().PackBlockTxs(ctx); err != nil {
		log.Debug("pack txs err", "err", err)
		return true, ctx.IsTimeout()
	}

	w.current.txs = append(w.current.txs, ctx.GetPackedTxList()...)
	w.current.tcount += len(w.current.txs)
	w.current.receipts = append(w.current.receipts, ctx.GetReceipts()...)
	//w.current.header.GasUsed = ctx.GetBlockGasUsed()
	coalescedLogs = append(coalescedLogs, ctx.GetLogs()...)

	/*root := w.current.state.IntermediateRoot(true)
	log.Debug("pack block", "w.current.state.Root", root.Hex())
	log.Debug("pack block", "w.current.txs", len(w.current.txs), "number", w.current.header.Number.Uint64())
	log.Debug("pack block", "w.current.receipts", len(w.current.receipts), "number", w.current.header.Number.Uint64())
	log.Debug("pack block", "coalescedLogs", len(coalescedLogs), "number", w.current.header.Number.Uint64())
	log.Debug("pack block", "w.current.gasPool", w.current.gasPool.Gas(), "number", w.current.header.Number.Uint64())
	log.Debug("pack block", "w.current.header.GasUsed", w.current.header.GasUsed, "number", w.current.header.Number.Uint64())*/

	if !w.isRunning() && len(coalescedLogs) > 0 {
		// We don't push the pendingLogsEvent while we are mining. The reason is that
		// when we are mining, the worker will regenerate a mining block every 3 seconds.
		// In order to avoid pushing the repeated pendingLog, we disable the pending log pushing.

		// make a copy, the state caches the logs and these logs get "upgraded" from pending to mined
		// logs by filling in the block hash when the block was mined by the local miner. This can
		// cause a race condition if a log was "upgraded" before the PendingLogsEvent is processed.
		cpy := make([]*types.Log, len(coalescedLogs))
		for i, l := range coalescedLogs {
			cpy[i] = new(types.Log)
			*cpy[i] = *l
		}
		go w.mux.Post(core.PendingLogsEvent{Logs: cpy})
	}
	// Notify resubmit loop to decrease resubmitting interval if current interval is larger
	// than the user-specified one.
	if interrupt != nil {
		w.resubmitAdjustCh <- &intervalAdjust{inc: false}
	}
	return false, ctx.IsTimeout()
}
