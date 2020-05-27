package miner

import (
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/vm"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

type TxsCommitter struct {
	worker *worker
}

func NewTxsCommitter(w *worker) *TxsCommitter {
	return &TxsCommitter{
		worker: w,
	}
}

func (c *TxsCommitter) CommitTransactions(header *types.Header, txs *types.TransactionsByPriceAndNonce, interrupt *int32, timestamp int64, blockDeadline time.Time) (bool, bool) {
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
	var bftEngine = w.config.Cbft != nil

	for {
		now := time.Now()
		if bftEngine && (blockDeadline.Equal(now) || blockDeadline.Before(now)) {

			log.Warn("interrupt current tx-executing", "now", time.Now().UnixNano()/1e6, "timestamp", timestamp, "commitDuration", w.commitDuration, "deadlineDuration", common.Millis(blockDeadline)-timestamp)
			//log.Warn("interrupt current tx-executing cause timeout, and continue the remainder package process", "timeout", w.commitDuration, "txCount", w.current.tcount)
			timeout = true
			break
		}
		// In the following three cases, we will interrupt the execution of the transaction.
		// (1) new head block event arrival, the interrupt signal is 1
		// (2) worker start or restart, the interrupt signal is 1
		// (3) worker recreate the mining block with any newly arrived transactions, the interrupt signal is 2.
		// For the first two cases, the semi-finished work will be discarded.
		// For the third case, the semi-finished work will be submitted to the consensus engine.
		if interrupt != nil && atomic.LoadInt32(interrupt) != commitInterruptNone {
			// Notify resubmit loop to increase resubmitting interval due to too frequent commits.
			if atomic.LoadInt32(interrupt) == commitInterruptResubmit {
				ratio := float64(w.current.header.GasLimit-w.current.gasPool.Gas()) / float64(w.current.header.GasLimit)
				if ratio < 0.1 {
					ratio = 0.1
				}
				w.resubmitAdjustCh <- &intervalAdjust{
					ratio: ratio,
					inc:   true,
				}
			}
			return atomic.LoadInt32(interrupt) == commitInterruptNewHead, timeout
		}
		// If we don't have enough gas for any further transactions then we're done
		if w.current.gasPool.Gas() < params.TxGas {
			log.Trace("Not enough gas for further transactions", "have", w.current.gasPool, "want", params.TxGas)
			break
		}
		// Retrieve the next transaction and abort if all done
		tx := txs.Peek()
		if tx == nil {
			break
		}
		// Error may be ignored here. The error has already been checked
		// during transaction acceptance is the transaction pool.
		//
		// We use the eip155 signer regardless of the current hf.
		from, _ := types.Sender(w.current.signer, tx)
		// Check whether the tx is replay protected. If we're not in the EIP155 hf
		// phase, start ignoring the sender until we do.
		if !w.config.IsEIP155(w.current.header.Number) {
			log.Trace("Ignoring reply protected transaction", "hash", tx.Hash(), "eip155", w.config.EIP155Block)

			txs.Pop()
			continue
		}
		// Start executing the transaction
		w.current.state.Prepare(tx.Hash(), common.Hash{}, w.current.tcount)

		logs, err := w.commitTransaction(tx)

		switch err {
		case core.ErrGasLimitReached:
			// Pop the current out-of-gas transaction without shifting in the next from the account
			log.Warn("Gas limit exceeded for current block", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "tx.hash", tx.Hash(), "sender", from, "state", w.current.state)
			txs.Pop()

		case core.ErrNonceTooLow:
			// New head notification data race between the transaction pool and miner, shift
			//log.Warn("Skipping transaction with low nonce", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "tx.hash", tx.Hash(), "sender", from, "senderCurNonce", w.current.state.GetNonce(from), "tx.nonce", tx.Nonce())
			txs.Shift()

		case core.ErrNonceTooHigh:
			// Reorg notification data race between the transaction pool and miner, skip account =
			log.Warn("Skipping account with high nonce", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "tx.hash", tx.Hash(), "sender", from, "senderCurNonce", w.current.state.GetNonce(from), "tx.nonce", tx.Nonce())
			txs.Pop()

		case nil:
			// Everything ok, collect the logs and shift in the next transaction from the same account
			coalescedLogs = append(coalescedLogs, logs...)
			w.current.tcount++
			txs.Shift()

		case vm.ErrAbort:
			log.Warn("Skipping account with exec timeout tx", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "hash", tx.Hash(), "sender", from, "senderCurNonce", w.current.state.GetNonce(from), "txNonce", tx.Nonce())
			txs.Pop()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			log.Warn("Transaction failed, account skipped", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "hash", tx.Hash(), "hash", tx.Hash(), "err", err)
			txs.Shift()
		}
	}

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
	return false, timeout
}
