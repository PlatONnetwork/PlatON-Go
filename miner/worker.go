// Copyright 2015 The go-ethereum Authors
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

package miner

import (
	"encoding/hex"
	"fmt"
	"math/big"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

// environment is the worker's current environment and holds all of the current state information.
type environment struct {
	signer types.Signer

	state      *state.StateDB // apply state changes here
	snapshotDB snapshotdb.DB

	tcount  int           // tx count in cycle
	gasPool *core.GasPool // available gas used to pack transactions

	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt
}

func (evm *environment) RevertToDBSnapshot(snapshotDBID, stateDBID int) {
	evm.snapshotDB.RevertToSnapshot(common.ZeroHash, snapshotDBID)
	evm.state.RevertToSnapshot(stateDBID)
}

func (evm *environment) DBSnapshot() (snapshotID int, StateDBID int) {
	return evm.snapshotDB.Snapshot(common.ZeroHash), evm.state.Snapshot()
}

// task contains all information for consensus engine sealing and result submitting.
type task struct {
	receipts  []*types.Receipt
	state     *state.StateDB
	block     *types.Block
	createdAt time.Time
}

const (
	commitInterruptNone int32 = iota
	commitInterruptNewHead
	commitInterruptResubmit
)

const (
	commitStatusIdle int32 = iota
	commitStatusCommitting
)

// newWorkReq represents a request for new sealing work submitting with relative interrupt notifier.
type newWorkReq struct {
	interrupt     *int32
	noempty       bool
	timestamp     time.Time
	blockDeadline time.Time
	commitBlock   *types.Block
}

// intervalAdjust represents a resubmitting interval adjustment.
type intervalAdjust struct {
	ratio float64
	inc   bool
}

/*type commitWorkEnv struct {
	baseLock            sync.RWMutex
	currentBaseBlock     *types.Block
	commitTime          int64
	highestLock         sync.RWMutex
	nextBaseBlock *types.Block
}*/

type commitWorkEnv struct {
	currentBaseBlock atomic.Value
	commitStatus     int32
	nextBlockTime    atomic.Value //time.Time
}

func (e *commitWorkEnv) getCurrentBaseBlock() *types.Block {
	v := e.currentBaseBlock.Load()
	if v == nil {
		return nil
	} else {
		return v.(*types.Block)
	}
}

func (e *commitWorkEnv) setCommitStatusIdle() {
	atomic.StoreInt32(&e.commitStatus, commitStatusIdle)
}

// worker is the main object which takes care of submitting new work to consensus engine
// and gathering the sealing result.
type worker struct {
	EmptyBlock   string
	config       *params.ChainConfig
	miningConfig *core.MiningConfig
	vmConfig     *vm.Config
	engine       consensus.Engine
	eth          Backend
	chain        *core.BlockChain

	gasFloor uint64
	//gasCeil  uint64

	// Subscriptions
	mux          *event.TypeMux
	txsCh        chan core.NewTxsEvent
	txsSub       event.Subscription
	chainHeadCh  chan core.ChainHeadEvent
	chainHeadSub event.Subscription
	chainSideCh  chan core.ChainSideEvent
	chainSideSub event.Subscription

	// Channels
	newWorkCh             chan *newWorkReq
	taskCh                chan *task
	resultCh              chan *types.Block
	prepareResultCh       chan *types.Block
	prepareCompleteCh     chan struct{}
	highestLogicalBlockCh chan *types.Block
	startCh               chan struct{}
	exitCh                chan struct{}
	resubmitIntervalCh    chan time.Duration
	resubmitAdjustCh      chan *intervalAdjust

	current     *environment       // An environment for current running cycle.
	unconfirmed *unconfirmedBlocks // A set of locally mined blocks pending canonicalness confirmations.

	mu       sync.RWMutex // The lock used to protect the coinbase and extra fields
	coinbase common.Address
	//extra    []byte

	pendingMu    sync.RWMutex
	pendingTasks map[common.Hash]*task

	snapshotMu    sync.RWMutex // The lock used to protect the block snapshot and state snapshot
	snapshotBlock *types.Block
	snapshotState *state.StateDB

	// atomic status counters
	running int32 // The indicator whether the consensus engine is running or not.
	newTxs  int32 // New arrival transaction count since last sealing work submitting.

	// External functions
	isLocalBlock func(block *types.Block) bool // Function used to determine whether the specified block is mined by local miner.

	blockChainCache *core.BlockChainCache
	commitWorkEnv   *commitWorkEnv
	recommit        time.Duration
	commitDuration  int64 //in Millisecond

	bftResultSub *event.TypeMuxSubscription
	// Test hooks
	newTaskHook  func(*task)                        // Method to call upon receiving a new sealing task.
	skipSealHook func(*task) bool                   // Method to decide whether skipping the sealing.
	fullTaskHook func()                             // Method to call before pushing the full sealing task.
	resubmitHook func(time.Duration, time.Duration) // Method to call upon updating resubmitting interval.

	committer core.Committer

	vmTimeout uint64
}

func newWorker(config *params.ChainConfig, miningConfig *core.MiningConfig, vmConfig *vm.Config, engine consensus.Engine,
	eth Backend, mux *event.TypeMux, recommit time.Duration, gasFloor uint64, isLocalBlock func(*types.Block) bool,
	blockChainCache *core.BlockChainCache, vmTimeout uint64) *worker {

	worker := &worker{
		config:       config,
		miningConfig: miningConfig,
		vmConfig:     vmConfig,
		engine:       engine,
		eth:          eth,
		mux:          mux,
		chain:        eth.BlockChain(),
		gasFloor:     gasFloor,
		//gasCeil:            gasCeil,
		isLocalBlock:       isLocalBlock,
		unconfirmed:        newUnconfirmedBlocks(eth.BlockChain(), miningConfig.MiningLogAtDepth),
		pendingTasks:       make(map[common.Hash]*task),
		txsCh:              make(chan core.NewTxsEvent, miningConfig.TxChanSize),
		chainHeadCh:        make(chan core.ChainHeadEvent, miningConfig.ChainHeadChanSize),
		chainSideCh:        make(chan core.ChainSideEvent, miningConfig.ChainSideChanSize),
		newWorkCh:          make(chan *newWorkReq),
		taskCh:             make(chan *task),
		resultCh:           make(chan *types.Block, miningConfig.ResultQueueSize),
		prepareResultCh:    make(chan *types.Block, miningConfig.ResultQueueSize),
		prepareCompleteCh:  make(chan struct{}, 1),
		exitCh:             make(chan struct{}),
		startCh:            make(chan struct{}, 1),
		resubmitIntervalCh: make(chan time.Duration),
		resubmitAdjustCh:   make(chan *intervalAdjust, miningConfig.ResubmitAdjustChanSize),
		blockChainCache:    blockChainCache,
		commitWorkEnv:      &commitWorkEnv{},
		vmTimeout:          vmTimeout,
	}
	// Subscribe NewTxsEvent for tx pool
	// worker.txsSub = eth.TxPool().SubscribeNewTxsEvent(worker.txsCh)
	// Subscribe events for blockchain
	worker.chainHeadSub = eth.BlockChain().SubscribeChainHeadEvent(worker.chainHeadCh)
	worker.chainSideSub = eth.BlockChain().SubscribeChainSideEvent(worker.chainSideCh)
	worker.bftResultSub = worker.mux.Subscribe(cbfttypes.CbftResult{})
	// Sanitize recommit interval if the user-specified one is too short.
	if config.Cbft != nil {
		recommit = time.Duration(config.Cbft.Period) * time.Second
	}
	if recommit < miningConfig.MinRecommitInterval {
		log.Warn("Sanitizing miner recommit interval", "provided", recommit, "updated", miningConfig.MinRecommitInterval)
		recommit = miningConfig.MinRecommitInterval
	}

	worker.EmptyBlock = config.EmptyBlock

	worker.recommit = recommit
	worker.commitDuration = int64((float64)(recommit.Nanoseconds()/1e6) * miningConfig.DefaultCommitRatio)
	log.Info("CommitDuration in Millisecond", "commitDuration", worker.commitDuration)

	worker.commitWorkEnv.nextBlockTime.Store(time.Now())

	worker.setCommitter(NewParallelTxsCommitter(worker))
	//worker.setCommitter(NewTxsCommitter(worker))

	go worker.mainLoop()
	go worker.newWorkLoop(recommit)
	go worker.resultLoop()
	go worker.taskLoop()

	// Submit first work to initialize pending state.
	worker.startCh <- struct{}{}

	return worker
}

// setEtherbase sets the etherbase used to initialize the block coinbase field.
//func (w *worker) setEtherbase(addr common.Address) {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//	w.coinbase = addr
//}

//// setExtra sets the content used to initialize the block extra field.
//func (w *worker) setExtra(extra []byte) {
//	w.mu.Lock()
//	defer w.mu.Unlock()
//	w.extra = extra
//}

// setRecommitInterval updates the interval for miner sealing work recommitting.
func (w *worker) setRecommitInterval(interval time.Duration) {
	w.resubmitIntervalCh <- interval
}

func (w *worker) setCommitter(committer core.Committer) {
	w.committer = committer
}

// pending returns the pending state and corresponding block.
func (w *worker) pending() (*types.Block, *state.StateDB) {
	// return a snapshot to avoid contention on currentMu mutex
	if _, ok := w.engine.(consensus.Bft); ok {
		return w.makePending()
	} else {
		w.snapshotMu.RLock()
		defer w.snapshotMu.RUnlock()
		if w.snapshotState == nil {
			return nil, nil
		}
		return w.snapshotBlock, w.snapshotState.Copy()
	}
}

// pendingBlock returns pending block.
func (w *worker) pendingBlock() *types.Block {
	// return a snapshot to avoid contention on currentMu mutex
	if _, ok := w.engine.(consensus.Bft); ok {
		pendingBlock, st := w.makePending()
		st.ClearParentReference()
		return pendingBlock
	} else {
		w.snapshotMu.RLock()
		defer w.snapshotMu.RUnlock()
		return w.snapshotBlock
	}
}

// start sets the running status as 1 and triggers new work submitting.
func (w *worker) start() {
	atomic.StoreInt32(&w.running, 1)
	w.startCh <- struct{}{}
}

// stop sets the running status as 0.
func (w *worker) stop() {
	atomic.StoreInt32(&w.running, 0)
}

// isRunning returns an indicator whether worker is running or not.
func (w *worker) isRunning() bool {
	return atomic.LoadInt32(&w.running) == 1
}

// close terminates all background threads maintained by the worker.
// Note the worker does not support being closed multiple times.
func (w *worker) close() {
	close(w.exitCh)
}

// newWorkLoop is a standalone goroutine to submit new mining work upon received events.
func (w *worker) newWorkLoop(recommit time.Duration) {
	var (
		interrupt   *int32
		minRecommit = recommit // minimal resubmit interval specified by user.
		timestamp   time.Time  // timestamp for each round of mining.
	)

	vdEvent := w.mux.Subscribe(cbfttypes.UpdateValidatorEvent{})
	defer vdEvent.Unsubscribe()

	timer := time.NewTimer(0)
	<-timer.C // discard the initial tick

	// commit aborts in-flight transaction execution with given signal and resubmits a new one.
	commit := func(noempty bool, s int32, baseBlock *types.Block, blockDeadline time.Time) {
		if interrupt != nil {
			atomic.StoreInt32(interrupt, s)
		}
		if baseBlock == nil {
			// Just abort pending block
			return
		}
		interrupt = new(int32)
		log.Info("Begin to commit new worker", "baseBlockHash", baseBlock.Hash(), "baseBlockNumber", baseBlock.Number(), "timestamp", common.Millis(timestamp), "deadline", common.Millis(blockDeadline), "deadlineDuration", blockDeadline.Sub(timestamp))
		w.newWorkCh <- &newWorkReq{interrupt: interrupt, noempty: noempty, timestamp: timestamp, blockDeadline: blockDeadline, commitBlock: baseBlock}
		timer.Reset(blockDeadline.Sub(timestamp))
		atomic.StoreInt32(&w.newTxs, 0)
	}
	// recalcRecommit recalculates the resubmitting interval upon feedback.
	recalcRecommit := func(target float64, inc bool) {
		var (
			prev = float64(recommit.Nanoseconds())
			next float64
		)
		if inc {
			next = prev*(1-w.miningConfig.IntervalAdjustRatio) + w.miningConfig.IntervalAdjustRatio*(target+w.miningConfig.IntervalAdjustBias)
			// Recap if interval is larger than the maximum time interval
			if next > float64(w.miningConfig.MaxRecommitInterval.Nanoseconds()) {
				next = float64(w.miningConfig.MaxRecommitInterval.Nanoseconds())
			}
		} else {
			next = prev*(1-w.miningConfig.IntervalAdjustRatio) + w.miningConfig.IntervalAdjustRatio*(target-w.miningConfig.IntervalAdjustBias)
			// Recap if interval is less than the user specified minimum
			if next < float64(minRecommit.Nanoseconds()) {
				next = float64(minRecommit.Nanoseconds())
			}
		}
		recommit = time.Duration(int64(next))
	}
	// clearPending cleans the stale pending tasks.
	clearPending := func(number uint64) {
		w.pendingMu.Lock()
		for h, t := range w.pendingTasks {
			if t.block.NumberU64()+w.miningConfig.StaleThreshold <= number {
				delete(w.pendingTasks, h)
			}
		}
		w.pendingMu.Unlock()
	}

	for {
		select {
		case <-w.startCh:
			timestamp = time.Now()
			log.Debug("Clear Pending", "number", w.chain.CurrentBlock().NumberU64())
			clearPending(w.chain.CurrentBlock().NumberU64())
			if _, ok := w.engine.(consensus.Bft); ok {
				//w.makePending()
				timer.Reset(50 * time.Millisecond)
			} else {
				commit(false, commitInterruptNewHead, nil, time.Now().Add(50*time.Millisecond))
			}

		case head := <-w.chainHeadCh:
			timestamp = time.Now()
			log.Debug("Clear Pending", "number", head.Block.NumberU64())
			clearPending(head.Block.NumberU64())
			//commit(false, commitInterruptNewHead)
			// clear consensus cache
			log.Debug("received a event of ChainHeadEvent", "hash", head.Block.Hash(), "number", head.Block.NumberU64(), "parentHash", head.Block.ParentHash())

			status := atomic.LoadInt32(&w.commitWorkEnv.commitStatus)
			current := w.commitWorkEnv.getCurrentBaseBlock()
			isNewHead := current == nil || (head.Block.Number().Cmp(current.Number()) >= 0 && head.Block.Hash() != current.Hash())
			if status == commitStatusCommitting && isNewHead {
				// Interrupt committing.
				if current != nil {
					log.Debug("Interrupt committing while received a event of ChainHeadEvent",
						"hash", head.Block.Hash(), "number", head.Block.NumberU64(), "parentHash", head.Block.ParentHash(),
						"currentNumber", current.Number(), "currentHash", current.Hash(), "currentParent", current.ParentHash())
				} else {
					log.Debug("Interrupt committing while received a event of ChainHeadEvent",
						"hash", head.Block.Hash(), "number", head.Block.NumberU64(), "parentHash", head.Block.ParentHash())
				}
				commit(false, commitInterruptNewHead, nil, time.Now().Add(50*time.Millisecond))
			}

		case <-timer.C:
			// If mining is running resubmit a new work cycle periodically to pull in
			// higher priced transactions. Disable this overhead for pending blocks.
			timestamp = time.Now()
			status := atomic.LoadInt32(&w.commitWorkEnv.commitStatus)
			if w.isRunning() {
				if cbftEngine, ok := w.engine.(consensus.Bft); ok {
					if status == commitStatusIdle {
						if shouldSeal, err := cbftEngine.ShouldSeal(timestamp); err == nil {
							if shouldSeal {
								if shouldCommit, commitBlock := w.shouldCommit(timestamp); shouldCommit {
									log.Debug("Begin to package new block regularly")
									blockDeadline := w.engine.(consensus.Bft).CalcBlockDeadline(timestamp)
									commit(false, commitInterruptResubmit, commitBlock, blockDeadline)
									continue
								}
							}
						}
					}
					timer.Reset(50 * time.Millisecond)
				} else if w.config.Clique == nil || w.config.Clique.Period > 0 {
					// Short circuit if no new transaction arrives.
					if atomic.LoadInt32(&w.newTxs) == 0 {
						timer.Reset(recommit)
						continue
					}
					//timestamp = time.Now().UnixNano() / 1e6
					commit(true, commitInterruptResubmit, nil, time.Now().Add(50*time.Millisecond))
				}
			}

		case interval := <-w.resubmitIntervalCh:
			timestamp = time.Now()
			if _, ok := w.engine.(consensus.Bft); !ok {
				// Adjust resubmit interval explicitly by user.
				if interval < w.miningConfig.MinRecommitInterval {
					log.Warn("Sanitizing miner recommit interval", "provided", interval, "updated", w.miningConfig.MinRecommitInterval)
					interval = w.miningConfig.MinRecommitInterval
				}
				log.Info("Miner recommit interval update", "from", minRecommit, "to", interval)
				minRecommit, recommit = interval, interval

				if w.resubmitHook != nil {
					w.resubmitHook(minRecommit, recommit)
				}
			}

		case adjust := <-w.resubmitAdjustCh:
			timestamp = time.Now()

			if _, ok := w.engine.(consensus.Bft); !ok {
				// Adjust resubmit interval by feedback.
				if adjust.inc {
					before := recommit
					recalcRecommit(float64(recommit.Nanoseconds())/adjust.ratio, true)
					log.Trace("Increase miner recommit interval", "from", before, "to", recommit)
				} else {
					before := recommit
					recalcRecommit(float64(minRecommit.Nanoseconds()), false)
					log.Trace("Decrease miner recommit interval", "from", before, "to", recommit)
				}

				if w.resubmitHook != nil {
					w.resubmitHook(minRecommit, recommit)
				}
			}

		case <-w.exitCh:
			return
		}
	}
}

// mainLoop is a standalone goroutine to regenerate the sealing task based on the received event.
func (w *worker) mainLoop() {
	// defer w.txsSub.Unsubscribe()
	defer w.chainHeadSub.Unsubscribe()
	defer w.chainSideSub.Unsubscribe()

	for {
		select {
		case req := <-w.newWorkCh:
			if err := w.commitNewWork(req.interrupt, req.noempty, common.Millis(req.timestamp), req.commitBlock, req.blockDeadline); err != nil {
				// If error during this committing, the task ends and change the CommitStatus to idle to allow the next commiting to be triggered
				log.Warn("Failed to commitNewWork", "baseBlockNumber", req.commitBlock.NumberU64(), "baseBlockHash", req.commitBlock.Hash(), "error", err)
				w.commitWorkEnv.setCommitStatusIdle()
			}

		case <-w.chainSideCh:

		case <-w.exitCh:
			// System stopped
			return

		case <-w.chainHeadSub.Err():
			return
		case <-w.chainSideSub.Err():
			return

		case <-w.prepareCompleteCh:
			// Indicates that a seal operation has completed, change the CommitStatus to idle regardless of success or failure
			w.commitWorkEnv.setCommitStatusIdle()

		case block := <-w.prepareResultCh:
			// Short circuit when receiving empty result.
			if block == nil {
				continue
			}
			// Short circuit when receiving duplicate result caused by resubmitting.
			if w.chain.HasBlock(block.Hash(), block.NumberU64()) {
				continue
			}
			var (
				sealhash = w.engine.SealHash(block.Header())
				hash     = block.Hash()
			)

			w.pendingMu.RLock()
			_, exist := w.pendingTasks[sealhash]
			w.pendingMu.RUnlock()
			if !exist {
				log.Error("Block found but no relative pending task", "number", block.Number(), "sealhash", sealhash, "hash", hash)
				continue
			}
		}
	}
}

// taskLoop is a standalone goroutine to fetch sealing task from the generator and
// push them to consensus engine.
func (w *worker) taskLoop() {
	var (
		stopCh chan struct{}
		prev   common.Hash
	)

	// interrupt aborts the in-flight sealing task.
	interrupt := func() {
		if stopCh != nil {
			close(stopCh)
			stopCh = nil
		}
	}
	for {
		select {
		case task := <-w.taskCh:
			if w.newTaskHook != nil {
				w.newTaskHook(task)
			}
			// Reject duplicate sealing work due to resubmitting.
			sealHash := w.engine.SealHash(task.block.Header())
			if sealHash == prev {
				w.commitWorkEnv.setCommitStatusIdle()
				continue
			}
			// Interrupt previous sealing operation
			interrupt()
			stopCh, prev = make(chan struct{}), sealHash

			if w.skipSealHook != nil && w.skipSealHook(task) {
				w.commitWorkEnv.setCommitStatusIdle()
				continue
			}
			w.pendingMu.Lock()
			w.pendingTasks[sealHash] = task
			w.pendingMu.Unlock()

			if cbftEngine, ok := w.engine.(consensus.Bft); ok {

				// Save stateDB to cache, receipts to cache
				w.blockChainCache.WriteStateDB(sealHash, task.state, task.block.NumberU64())
				w.blockChainCache.WriteReceipts(sealHash, task.receipts, task.block.NumberU64())
				w.blockChainCache.AddSealBlock(sealHash, task.block.NumberU64())
				log.Debug("Add seal block to blockchain cache", "sealHash", sealHash, "number", task.block.NumberU64())
				if err := cbftEngine.Seal(w.chain, task.block, w.prepareResultCh, stopCh, w.prepareCompleteCh); err != nil {
					log.Warn("Block sealing failed on bft engine", "err", err)
					w.commitWorkEnv.setCommitStatusIdle()
				}
				continue
			}

			if err := w.engine.Seal(w.chain, task.block, w.resultCh, stopCh, w.prepareCompleteCh); err != nil {
				log.Warn("Block sealing failed", "err", err)
				w.commitWorkEnv.setCommitStatusIdle()
			}

		case <-w.exitCh:
			interrupt()
			return
		}
	}
}

// resultLoop is a standalone goroutine to handle sealing result submitting
// and flush relative data to the database.
func (w *worker) resultLoop() {
	for {
		select {
		case obj := <-w.bftResultSub.Chan():
			if obj == nil {
				//log.Error("receive nil maybe channel is closed")
				continue
			}
			cbftResult, ok := obj.Data.(cbfttypes.CbftResult)
			if !ok {
				log.Error("Receive bft result type error")
				continue
			}
			block := cbftResult.Block
			// Short circuit when receiving empty result.
			if block == nil {
				log.Error("Cbft result error, block is nil")
				continue
			}

			var (
				hash     = block.Hash()
				sealhash = w.engine.SealHash(block.Header())
				number   = block.NumberU64()
			)
			// Short circuit when receiving duplicated cbft result caused by resubmitting or P2P sync.
			if w.chain.HasBlock(block.Hash(), block.NumberU64()) {
				log.Warn("Duplicated cbft result caused by resubmitting or P2P sync.", "hash", hash, "number", number)
				continue
			}
			w.pendingMu.RLock()
			task, exist := w.pendingTasks[sealhash]
			w.pendingMu.RUnlock()

			log.Debug("Pending task", "exist", exist)
			var _receipts []*types.Receipt
			var _state *state.StateDB
			//todo remove extra magic number
			if exist && w.engine.(consensus.Bft).IsSignedBySelf(sealhash, block.Header()) {
				_receipts = task.receipts
				_state = task.state
				stateIsNil := _state == nil
				log.Debug("Block is packaged by local", "hash", hash, "number", number, "len(Receipts)", len(_receipts), "stateIsNil", stateIsNil)
			} else {
				_receipts = w.blockChainCache.ReadReceipts(sealhash)
				_state = w.blockChainCache.ReadStateDB(sealhash)
				stateIsNil := _state == nil
				log.Debug("Block is packaged by other", "hash", hash, "number", number, "len(Receipts)", len(_receipts), "blockRoot", block.Root(), "stateIsNil", stateIsNil)
			}
			if _state == nil {
				log.Warn("Handle cbft result error, state is nil, maybe block is synced from other peer", "hash", hash, "number", number)
				continue
			} else if len(block.Transactions()) > 0 && len(_receipts) == 0 {
				log.Warn("Handle cbft result error, block has transactions but receipts is nil, maybe block is synced from other peer", "hash", hash, "number", number)
				continue
			}
			log.Debug("Cbft consensus successful", "hash", hash, "number", number, "timestamp", time.Now().UnixNano()/1e6)

			// Different block could share same sealhash, deep copy here to prevent write-write conflict.
			var (
				receipts = make([]*types.Receipt, len(_receipts))
				logs     []*types.Log
			)
			for i, receipt := range _receipts {
				// add block location fields
				receipt.BlockHash = hash
				receipt.BlockNumber = block.Number()
				receipt.TransactionIndex = uint(i)

				receipts[i] = new(types.Receipt)
				*receipts[i] = *receipt
				// Update the block hash in all logs since it is now available and not when the
				// receipt/log of individual transactions were created.
				for _, log := range receipt.Logs {
					log.BlockHash = hash
				}
				logs = append(logs, receipt.Logs...)
			}
			// Commit block and state to database.
			block.SetExtraData(cbftResult.ExtraData)
			log.Debug("Write extra data", "txs", len(block.Transactions()), "extra", len(block.ExtraData()))
			// update 3-chain state
			cbftResult.ChainStateUpdateCB()
			stat, err := w.chain.WriteBlockWithState(block, receipts, _state)
			if err != nil {
				if cbftResult.SyncState != nil {
					cbftResult.SyncState <- err
				}
				log.Error("Failed writing block to chain", "hash", block.Hash(), "number", block.NumberU64(), "err", err)
				continue
			}

			//cbftResult.SyncState <- err
			log.Info("Successfully write new block", "hash", block.Hash(), "number", block.NumberU64(), "coinbase", block.Coinbase(), "time", block.Time(), "root", block.Root())

			// Broadcast the block and announce chain insertion event
			w.mux.Post(core.NewMinedBlockEvent{Block: block})

			var events []interface{}
			switch stat {
			case core.CanonStatTy:
				events = append(events, core.ChainEvent{Block: block, Hash: block.Hash(), Logs: logs})
				events = append(events, core.ChainHeadEvent{Block: block})
			case core.SideStatTy:
				events = append(events, core.ChainSideEvent{Block: block})
			}
			w.chain.PostChainEvents(events, logs)

		case <-w.exitCh:
			return
		}
	}
}

// makeCurrent creates a new environment for the current cycle.
func (w *worker) makeCurrent(parent *types.Block, header *types.Header) error {
	var (
		state *state.StateDB
		err   error
	)
	if _, ok := w.engine.(consensus.Bft); ok {
		state, err = w.blockChainCache.MakeStateDB(parent)
	} else {
		state, err = w.chain.StateAt(parent.Root())
	}
	if err != nil {
		return err
	}
	env := &environment{
		signer:     types.NewEIP155Signer(w.config.ChainID),
		snapshotDB: snapshotdb.Instance(),
		state:      state,
		header:     header,
	}

	// Keep track of transactions which return errors so they can be removed
	env.tcount = 0
	w.current = env
	return nil
}

// updateSnapshot updates pending snapshot block and state.
// Note this function assumes the current variable is thread safe.
func (w *worker) updateSnapshot() {
	w.snapshotMu.Lock()
	defer w.snapshotMu.Unlock()

	w.snapshotBlock = types.NewBlock(
		w.current.header,
		w.current.txs,
		w.current.receipts,
	)

	w.snapshotState = w.current.state.Copy()
}

func (w *worker) commitTransaction(tx *types.Transaction) ([]*types.Log, error) {
	snapForSnap, snapForState := w.current.DBSnapshot()

	receipt, _, err := core.ApplyTransaction(w.config, w.chain, w.current.gasPool, w.current.state,
		w.current.header, tx, &w.current.header.GasUsed, *w.chain.GetVMConfig())
	if err != nil {
		log.Error("Failed to commitTransaction on worker", "blockNumer", w.current.header.Number.Uint64(), "txHash", tx.Hash().String(), "err", err)
		w.current.RevertToDBSnapshot(snapForSnap, snapForState)
		return nil, err
	}
	w.current.txs = append(w.current.txs, tx)
	w.current.receipts = append(w.current.receipts, receipt)

	return receipt.Logs, nil
}

func (w *worker) commitTransactionsWithHeader(header *types.Header, txs *types.TransactionsByPriceAndNonce, interrupt *int32, timestamp int64, blockDeadline time.Time) (bool, bool) {
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
			log.Warn("Gas limit exceeded for current block", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "tx.hash", tx.Hash(), "sender", from, "senderCurNonce", w.current.state.GetNonce(from), "txNonce", tx.Nonce())
			txs.Pop()

		case core.ErrNonceTooLow:
			// New head notification data race between the transaction pool and miner, shift
			log.Warn("Skipping transaction with low nonce", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "tx.hash", tx.Hash(), "sender", from, "senderCurNonce", w.current.state.GetNonce(from), "tx.nonce", tx.Nonce())
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
			log.Warn("Skipping account with exec timeout tx", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "tx.hash", tx.Hash(), "sender", from, "senderCurNonce", w.current.state.GetNonce(from), "txNonce", tx.Nonce())
			txs.Pop()

		case vm.ErrWASMUndefinedPanic:
			log.Warn("Skipping account with wasm vm exec undefined err", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "tx.hash", tx.Hash(), "sender", from, "senderCurNonce", w.current.state.GetNonce(from), "txNonce", tx.Nonce())
			txs.Pop()

		default:
			// Strange error, discard the transaction and get the next in line (note, the
			// nonce-too-high clause will prevent us from executing in vain).
			log.Warn("Transaction failed, account skipped", "blockNumber", header.Number, "blockParentHash", header.ParentHash, "tx.hash", tx.Hash(), "err", err)
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

// commitNewWork generates several new sealing tasks based on the parent block.
func (w *worker) commitNewWork(interrupt *int32, noempty bool, timestamp int64, commitBlock *types.Block, blockDeadline time.Time) error {
	w.mu.RLock()
	defer w.mu.RUnlock()

	atomic.StoreInt32(&w.commitWorkEnv.commitStatus, commitStatusCommitting)

	defer func() {
		if engine, ok := w.engine.(consensus.Bft); ok {
			w.commitWorkEnv.nextBlockTime.Store(engine.CalcNextBlockTime(common.MillisToTime(timestamp)))
			log.Debug("Next block time", "time", common.Beautiful(w.commitWorkEnv.nextBlockTime.Load().(time.Time)))
		}
	}()

	tstart := time.Now()

	var parent *types.Block
	if _, ok := w.engine.(consensus.Bft); ok {
		parent = commitBlock
		//timestamp = time.Now().UnixNano() / 1e6
	} else {
		parent = w.chain.CurrentBlock()
		if parent.Time().Cmp(new(big.Int).SetInt64(timestamp)) >= 0 {
			timestamp = parent.Time().Int64() + 1
		}
		// this will ensure we're not going off too far in the future
		if now := time.Now().Unix(); timestamp > now+1 {
			wait := time.Duration(timestamp-now) * time.Second
			log.Info("Mining too far in the future", "wait", common.PrettyDuration(wait))
			time.Sleep(wait)
		}
	}

	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		GasLimit:   core.CalcGasLimit(parent, w.gasFloor),
		Time:       big.NewInt(timestamp),
	}

	log.Info("Cbft begin to consensus for new block", "number", header.Number, "nonce", hexutil.Encode(header.Nonce[:]), "gasLimit", header.GasLimit, "parentHash", parent.Hash(), "parentNumber", parent.NumberU64(), "parentStateRoot", parent.Root(), "timestamp", common.MillisToString(timestamp))
	// Initialize the header extra in Prepare function of engine
	if err := w.engine.Prepare(w.chain, header); err != nil {
		log.Error("Failed to prepare header for mining", "err", err)
		return err
	}

	// Could potentially happen if starting to mine in an odd state.
	err := w.makeCurrent(parent, header)
	if err != nil {
		log.Error("Failed to create mining context", "err", err)
		return err
	}
	//make header extra after w.current and it's state initialized
	extraData := w.makeExtraData()
	copy(header.Extra[:len(extraData)], extraData)

	// BeginBlocker()
	if err := core.GetReactorInstance().BeginBlocker(header, w.current.state); nil != err {
		log.Error("Failed to GetReactorInstance BeginBlocker on worker", "blockNumber", header.Number, "err", err)
		return err
	}

	// Only set the coinbase if our consensus engine is running (avoid spurious block rewards)
	if w.isRunning() {
		if b, ok := w.engine.(consensus.Bft); ok {
			core.GetReactorInstance().SetWorkerCoinBase(header, b.NodeID())
		}
	}

	if !noempty && "on" == w.EmptyBlock {
		// Create an empty block based on temporary copied state for sealing in advance without waiting block
		// execution finished.
		if _, ok := w.engine.(consensus.Bft); !ok {
			if err := w.commit(nil, false, tstart); nil != err {
				log.Error("Failed to commitNewWork on worker: call commit is failed", "blockNumber", header.Number, "err", err)
			}
		}
	}

	log.Trace("Validator mode", "mode", w.config.Cbft.ValidatorMode)
	if w.config.Cbft.ValidatorMode == "inner" {
		// Check if need to switch validators.
		// If needed, make a inner contract transaction
		// and pack into pending block.
		if w.shouldSwitch() && w.commitInnerTransaction(timestamp, blockDeadline) != nil {
			return fmt.Errorf("commit inner transaction error")
		}
	}

	// Fill the block with all available pending transactions.
	startTime := time.Now()
	var pending map[common.Address]types.Transactions

	pending, err = w.eth.TxPool().PendingLimited()
	if err != nil {
		log.Error("Failed to fetch pending transactions", "time", common.PrettyDuration(time.Since(startTime)), "err", err)
		return err
	}

	log.Debug("Fetch pending transactions success", "number", header.Number, "pendingLength", len(pending), "time", common.PrettyDuration(time.Since(startTime)))

	// Short circuit if there is no available pending transactions
	if len(pending) == 0 {
		//// No empty block
		//if "off" == w.EmptyBlock {
		//	return
		//}
		if _, ok := w.engine.(consensus.Bft); ok {
			if err := w.commit(nil, true, tstart); nil != err {
				log.Error("Failed to commitNewWork on worker: call commit is failed", "blockNumber", header.Number, "err", err)
				return err
			}
		} else {
			w.updateSnapshot()
		}
		return nil
	}

	txsCount := 0
	for _, accTxs := range pending {
		txsCount = txsCount + len(accTxs)
	}
	// Split the pending transactions into locals and remotes
	localTxs, remoteTxs := make(map[common.Address]types.Transactions), pending
	for _, account := range w.eth.TxPool().Locals() {
		if txs := remoteTxs[account]; len(txs) > 0 {
			delete(remoteTxs, account)
			localTxs[account] = txs
		}
	}
	localTxsCount := 0
	remoteTxsCount := 0
	for _, laccTxs := range localTxs {
		localTxsCount = localTxsCount + len(laccTxs)
	}
	for _, raccTxs := range remoteTxs {
		remoteTxsCount = remoteTxsCount + len(raccTxs)
	}
	log.Debug("Execute pending transactions", "number", header.Number, "localTxCount", localTxsCount, "remoteTxCount", remoteTxsCount, "txsCount", txsCount)

	startTime = time.Now()
	var localTimeout = false
	if len(localTxs) > 0 {
		txs := types.NewTransactionsByPriceAndNonce(w.current.signer, localTxs)
		if failed, timeout := w.committer.CommitTransactions(header, txs, interrupt, timestamp, blockDeadline); failed {
			return fmt.Errorf("commit transactions error")
		} else {
			localTimeout = timeout
		}
	}

	commitLocalTxCount := w.current.tcount
	log.Debug("Local transactions executing stat", "number", header.Number, "involvedTxCount", commitLocalTxCount, "time", time.Since(startTime))

	startTime = time.Now()
	if !localTimeout && len(remoteTxs) > 0 {
		txs := types.NewTransactionsByPriceAndNonce(w.current.signer, remoteTxs)

		if failed, _ := w.committer.CommitTransactions(header, txs, interrupt, timestamp, blockDeadline); failed {
			return fmt.Errorf("commit transactions error")
		}
	}
	commitRemoteTxCount := w.current.tcount - commitLocalTxCount
	log.Debug("Remote transactions executing stat", "number", header.Number, "involvedTxCount", commitRemoteTxCount, "time", time.Since(startTime))

	if err := w.commit(w.fullTaskHook, true, tstart); nil != err {
		log.Error("Failed to commitNewWork on worker: call commit is failed", "blockNumber", header.Number, "err", err)
		return err
	}
	log.Info("Commit new work", "number", header.Number, "pending", txsCount, "txs", w.current.tcount, "diff", txsCount-w.current.tcount, "duration", time.Since(tstart))
	return nil
}

// commit runs any post-transaction state modifications, assembles the final block
// and commits new work if consensus engine is running.
func (w *worker) commit(interval func(), update bool, start time.Time) error {
	//if "off" == w.EmptyBlock && 0 == len(w.current.txs) {
	//	return nil
	//}

	// Deep copy receipts here to avoid interaction between different tasks.
	receipts := make([]*types.Receipt, len(w.current.receipts))
	for i, l := range w.current.receipts {
		receipts[i] = new(types.Receipt)
		*receipts[i] = *l
	}

	s := w.current.state.Copy()

	// EndBlocker()
	if err := core.GetReactorInstance().EndBlocker(w.current.header, s); nil != err {
		log.Error("Failed to GetReactorInstance EndBlocker on worker", "blockNumber",
			w.current.header.Number.Uint64(), "err", err)
		return err
	}

	block, err := w.engine.Finalize(w.chain, w.current.header, s, w.current.txs, w.current.receipts)

	if err != nil {
		return err
	}
	if update {
		w.updateSnapshot()
	}
	if w.isRunning() {
		if interval != nil {
			interval()
		}
		select {
		case w.taskCh <- &task{receipts: receipts, state: s, block: block, createdAt: time.Now()}:
			w.unconfirmed.Shift(block.NumberU64() - 1)

			feesWei := new(big.Int)
			for i, tx := range block.Transactions() {
				feesWei.Add(feesWei, new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), tx.GasPrice()))
			}
			feesEth := new(big.Float).Quo(new(big.Float).SetInt(feesWei), new(big.Float).SetInt(big.NewInt(params.LAT)))

			log.Info("Commit new mining work", "number", block.Number(), "sealhash", w.engine.SealHash(block.Header()), "receiptHash", block.ReceiptHash(),
				"txs", w.current.tcount, "gas", block.GasUsed(), "fees", feesEth, "elapsed", common.PrettyDuration(time.Since(start)))
		case <-w.exitCh:
			log.Info("Worker has exited")
		}
		return nil
	}
	w.commitWorkEnv.setCommitStatusIdle()
	return nil
}

func (w *worker) makePending() (*types.Block, *state.StateDB) {
	var parent = w.engine.NextBaseBlock()
	var parentChain = w.chain.CurrentBlock()

	if parentChain.NumberU64() >= parent.NumberU64() {
		parent = parentChain
	}
	log.Debug("Parent in makePending", "number", parent.NumberU64(), "hash", parent.Hash())

	if parent != nil {
		state, err := w.blockChainCache.MakeStateDB(parent)
		if err == nil {
			block := types.NewBlock(
				parent.Header(),
				parent.Transactions(),
				nil,
			)

			return block, state
		}
	}
	return nil, nil
}

func (w *worker) shouldCommit(timestamp time.Time) (bool, *types.Block) {
	currentBaseBlock := w.commitWorkEnv.getCurrentBaseBlock()
	nextBaseBlock := w.engine.NextBaseBlock()
	nextBaseBlockTime := common.MillisToTime(nextBaseBlock.Time().Int64())

	if timestamp.Before(nextBaseBlockTime) {
		log.Warn("Invalid packing timestamp,current timestamp is lower than the parent timestamp", "parentBlockTime", common.Beautiful(nextBaseBlockTime), "currentBlockTime", common.Beautiful(timestamp))
		return false, nil
	}

	nextBlockTime := w.commitWorkEnv.nextBlockTime.Load().(time.Time)
	blockTime := w.engine.(consensus.Bft).CalcNextBlockTime(nextBaseBlockTime)
	if nextBlockTime.Before(blockTime) && time.Now().Before(blockTime) {
		log.Debug("Invalid nextBlockTime,recalc it", "nextBlockTime", common.Beautiful(nextBlockTime), "blockTime", common.Beautiful(blockTime))
		w.commitWorkEnv.nextBlockTime.Store(blockTime)
		return false, nil
	}

	status := atomic.LoadInt32(&w.commitWorkEnv.commitStatus)
	shouldCommit := nextBlockTime.Before(time.Now()) && status == commitStatusIdle
	log.Trace("Check should commit", "shouldCommit", shouldCommit, "status", status, "timestamp", timestamp, "nextBlockTime", nextBlockTime)

	if shouldCommit && nextBaseBlock != nil {
		var err error
		w.commitWorkEnv.currentBaseBlock.Store(nextBaseBlock)
		if err != nil {
			log.Error("Calc next block time failed", "err", err)
			return false, nil
		}
		if currentBaseBlock == nil {
			log.Debug("Check if time's up in shouldCommit()", "result", shouldCommit,
				"next.number", nextBaseBlock.Number(),
				"next.hash", nextBaseBlock.Hash(),
				"next.timestamp", common.MillisToString(nextBaseBlock.Time().Int64()),
				"nextBlockTime", nextBlockTime,
				"timestamp", timestamp)
		} else {
			log.Debug("Check if time's up in shouldCommit()", "result", shouldCommit,
				"current.number", currentBaseBlock.Number(),
				"current.hash", currentBaseBlock.Hash(),
				"current.timestamp", common.MillisToString(currentBaseBlock.Time().Int64()),
				"next.number", nextBaseBlock.Number(),
				"next.hash", nextBaseBlock.Hash(),
				"next.timestamp", common.MillisToString(nextBaseBlock.Time().Int64()),
				"nextBlockTime", nextBlockTime,
				"timestamp", timestamp)
		}
	}
	return shouldCommit, nextBaseBlock
}

// make default extra data when preparing new block
func (w *worker) makeExtraData() []byte {
	// create default extradata
	extra, _ := rlp.EncodeToBytes([]interface{}{
		//uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
		gov.GetCurrentActiveVersion(w.current.state),
		"platon",
		runtime.Version(),
		runtime.GOOS,
	})
	if uint64(len(extra)) > params.MaximumExtraDataSize {
		log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
		extra = nil
	}

	log.Debug("Prepare header extra", "data", hex.EncodeToString(extra))
	return extra
}
