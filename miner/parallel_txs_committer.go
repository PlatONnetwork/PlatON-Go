package miner

import (
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
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

func (c *ParallelTxsCommitter) CommitTransactions(env *environment, txs *types.TransactionsByPriceAndNonce, interrupt *int32, timestamp int64, blockDeadline time.Time, tempContractCache map[common.Address]struct{}) (bool, bool) {

	w := c.worker

	// Short circuit if current is nil
	timeout := false

	if env == nil {
		return true, timeout
	}

	if env.gasPool == nil {
		env.gasPool = new(core.GasPool).AddGas(env.header.GasLimit)
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
	signer := types.MakeSigner(c.worker.chainConfig, env.header.Number, gov.Gte150VersionState(env.state))
	ctx := core.NewParallelContext(env.state, env.header, common.Hash{}, env.gasPool, true, signer, tempContractCache)
	ctx.SetBlockDeadline(blockDeadline)
	ctx.SetBlockGasUsedHolder(&(env.header.GasUsed))
	ctx.SetTxList(parallelTxs)
	log.Trace("Begin to execute transactions", "number", env.header.Number)
	if err := core.GetExecutor().ExecuteTransactions(ctx); err != nil {
		log.Debug("pack txs err", "err", err)
		return true, ctx.IsTimeout()
	}
	log.Trace("End to execute transactions", "number", env.header.Number)

	env.txs = append(env.txs, ctx.GetPackedTxList()...)
	env.tcount = len(env.txs)
	env.receipts = append(env.receipts, ctx.GetReceipts()...)
	//w.current.header.GasUsed = ctx.GetBlockGasUsed()
	coalescedLogs = append(coalescedLogs, ctx.GetLogs()...)

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
		w.pendingLogsFeed.Send(cpy)
	}
	// Notify resubmit loop to decrease resubmitting interval if current interval is larger
	// than the user-specified one.
	if interrupt != nil {
		w.resubmitAdjustCh <- &intervalAdjust{inc: false}
	}
	log.Debug("End to commit transactions", "number", env.header.Number)
	return false, ctx.IsTimeout()
}
