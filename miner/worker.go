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

	"github.com/PlatONnetwork/PlatON-Go/consensus/misc"

	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/trie"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/log"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

// environment is the worker's current environment and holds all
// information of the sealing block generation.
type environment struct {
	signer types.Signer

	state      *state.StateDB // apply state changes here
	snapshotDB snapshotdb.DB

	tcount   int           // tx count in cycle
	gasPool  *core.GasPool // available gas used to pack transactions
	coinbase common.Address

	header   *types.Header
	txs      []*types.Transaction
	receipts []*types.Receipt
}

// copy creates a deep copy of environment.
func (env *environment) copy() *environment {
	cpy := &environment{
		signer:     env.signer,
		state:      env.state.Copy(),
		snapshotDB: env.snapshotDB,
		tcount:     env.tcount,
		coinbase:   env.coinbase,
		header:     types.CopyHeader(env.header),
		receipts:   copyReceipts(env.receipts),
	}
	if env.gasPool != nil {
		gasPool := *env.gasPool
		cpy.gasPool = &gasPool
	}
	// The content of txs and uncles are immutable, unnecessary
	// to do the expensive deep copy for them.
	cpy.txs = make([]*types.Transaction, len(env.txs))
	copy(cpy.txs, env.txs)

	return cpy
}

// discard terminates the background prefetcher go-routine. It should
// always be called for all created environment instances otherwise
// the go-routine leak can happen.
func (env *environment) discard() {
	if env.state == nil {
		return
	}
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

	// maxRecommitInterval is the maximum time interval to recreate the mining block with
	// any newly arrived transactions.
	maxRecommitInterval = 15 * time.Second

	// intervalAdjustRatio is the impact a single interval adjustment has on sealing work
	// resubmitting interval.
	intervalAdjustRatio = 0.1

	// intervalAdjustBias is applied during the new resubmit interval calculation in favor of
	// increasing upper limit or decreasing lower limit so that the limit can be reachable.
	intervalAdjustBias = 200 * 1000.0 * 1000.0
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

// getWorkReq represents a request for getting a new sealing work with provided parameters.
type getWorkReq struct {
	params *generateParams
	err    error
	result chan *types.Block
}

// intervalAdjust represents a resubmitting interval adjustment.
type intervalAdjust struct {
	ratio float64
	inc   bool
}

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
	config       *Config
	chainConfig  *params.ChainConfig
	miningConfig *core.MiningConfig
	engine       consensus.Engine
	eth          Backend
	chain        *core.BlockChain

	// Feeds
	pendingLogsFeed event.Feed

	// Subscriptions
	mux          *event.TypeMux
	chainHeadCh  chan core.ChainHeadEvent
	chainHeadSub event.Subscription
	chainSideCh  chan core.ChainSideEvent
	chainSideSub event.Subscription

	// Channels
	newWorkCh             chan *newWorkReq
	getWorkCh             chan *getWorkReq
	taskCh                chan *task
	resultCh              chan *types.Block
	prepareResultCh       chan *types.Block
	prepareCompleteCh     chan struct{}
	highestLogicalBlockCh chan *types.Block
	startCh               chan struct{}
	exitCh                chan struct{}
	resubmitIntervalCh    chan time.Duration
	resubmitAdjustCh      chan *intervalAdjust

	wg sync.WaitGroup

	current     *environment       // An environment for current running cycle.
	unconfirmed *unconfirmedBlocks // A set of locally mined blocks pending canonicalness confirmations.

	mu       sync.RWMutex // The lock used to protect the coinbase and extra fields
	coinbase common.Address
	//extra    []byte

	pendingMu    sync.RWMutex
	pendingTasks map[common.Hash]*task

	snapshotMu       sync.RWMutex // The lock used to protect the snapshots below
	snapshotBlock    *types.Block
	snapshotReceipts types.Receipts
	snapshotState    *state.StateDB

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

	committer Committer

	vmTimeout uint64
}

func newWorker(config *Config, chainConfig *params.ChainConfig, miningConfig *core.MiningConfig, engine consensus.Engine,
	eth Backend, mux *event.TypeMux, isLocalBlock func(*types.Block) bool,
	blockChainCache *core.BlockChainCache, vmTimeout uint64) *worker {

	worker := &worker{
		config:             config,
		chainConfig:        chainConfig,
		miningConfig:       miningConfig,
		engine:             engine,
		eth:                eth,
		mux:                mux,
		chain:              eth.BlockChain(),
		isLocalBlock:       isLocalBlock,
		unconfirmed:        newUnconfirmedBlocks(eth.BlockChain(), miningConfig.MiningLogAtDepth),
		pendingTasks:       make(map[common.Hash]*task),
		chainHeadCh:        make(chan core.ChainHeadEvent, miningConfig.ChainHeadChanSize),
		chainSideCh:        make(chan core.ChainSideEvent, miningConfig.ChainSideChanSize),
		newWorkCh:          make(chan *newWorkReq),
		getWorkCh:          make(chan *getWorkReq),
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
	if chainConfig.Cbft != nil {
		worker.config.Recommit = time.Duration(chainConfig.Cbft.Period) * time.Second
	}
	if worker.config.Recommit < miningConfig.MinRecommitInterval {
		log.Warn("Sanitizing miner recommit interval", "provided", worker.config.Recommit, "updated", miningConfig.MinRecommitInterval)
		worker.config.Recommit = miningConfig.MinRecommitInterval
	}

	worker.EmptyBlock = chainConfig.EmptyBlock

	worker.recommit = worker.config.Recommit
	worker.commitDuration = int64((float64)(worker.config.Recommit.Nanoseconds()/1e6) * miningConfig.DefaultCommitRatio)
	log.Info("CommitDuration in Millisecond", "commitDuration", worker.commitDuration)

	worker.commitWorkEnv.nextBlockTime.Store(time.Now())

	worker.setCommitter(NewParallelTxsCommitter(worker))
	//worker.setCommitter(NewTxsCommitter(worker))

	worker.wg.Add(4)
	go worker.mainLoop()
	go worker.newWorkLoop(worker.config.Recommit)
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
	select {
	case w.resubmitIntervalCh <- interval:
	case <-w.exitCh:
	}
}

func (w *worker) setCommitter(committer Committer) {
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

// pendingBlockAndReceipts returns pending block and corresponding receipts.
func (w *worker) pendingBlockAndReceipts() (*types.Block, types.Receipts) {
	// return a snapshot to avoid contention on currentMu mutex
	w.snapshotMu.RLock()
	defer w.snapshotMu.RUnlock()
	return w.snapshotBlock, w.snapshotReceipts
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
	atomic.StoreInt32(&w.running, 0)
	close(w.exitCh)
	w.wg.Wait()
}

// recalcRecommit recalculates the resubmitting interval upon feedback.
func recalcRecommit(minRecommit, prev time.Duration, target float64, inc bool) time.Duration {
	var (
		prevF = float64(prev.Nanoseconds())
		next  float64
	)
	if inc {
		next = prevF*(1-intervalAdjustRatio) + intervalAdjustRatio*(target+intervalAdjustBias)
		max := float64(maxRecommitInterval.Nanoseconds())
		if next > max {
			next = max
		}
	} else {
		next = prevF*(1-intervalAdjustRatio) + intervalAdjustRatio*(target-intervalAdjustBias)
		min := float64(minRecommit.Nanoseconds())
		if next < min {
			next = min
		}
	}
	return time.Duration(int64(next))
}

// newWorkLoop is a standalone goroutine to submit new mining work upon received events.
func (w *worker) newWorkLoop(recommit time.Duration) {
	defer w.wg.Done()
	var (
		interrupt   *int32
		minRecommit = recommit // minimal resubmit interval specified by user.
		timestamp   time.Time  // timestamp for each round of mining.
	)

	vdEvent := w.mux.Subscribe(cbfttypes.UpdateValidatorEvent{})
	defer vdEvent.Unsubscribe()

	timer := time.NewTimer(0)
	defer timer.Stop()
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
		select {
		case w.newWorkCh <- &newWorkReq{interrupt: interrupt, noempty: noempty, timestamp: timestamp, blockDeadline: blockDeadline, commitBlock: baseBlock}:
		case <-w.exitCh:
			return
		}
		timer.Reset(blockDeadline.Sub(timestamp))
		atomic.StoreInt32(&w.newTxs, 0)
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
				} else {
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
					target := float64(recommit.Nanoseconds()) / adjust.ratio
					recommit = recalcRecommit(minRecommit, recommit, target, true)
					log.Trace("Increase miner recommit interval", "from", before, "to", recommit)
				} else {
					before := recommit
					recommit = recalcRecommit(minRecommit, recommit, float64(minRecommit.Nanoseconds()), false)
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
	defer w.wg.Done()
	// defer w.txsSub.Unsubscribe()
	defer w.chainHeadSub.Unsubscribe()
	defer w.chainSideSub.Unsubscribe()
	defer func() {
		if w.current != nil {
			w.current.discard()
		}
	}()

	for {
		select {
		case req := <-w.newWorkCh:
			if err := w.commitWork(req.interrupt, req.noempty, common.Millis(req.timestamp), req.commitBlock, req.blockDeadline); err != nil {
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
	defer w.wg.Done()
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
				task.state.UpdateSnaps()
				log.Debug("Add seal block to blockchain cache", "sealHash", sealHash, "number", task.block.NumberU64())
				if err := cbftEngine.Seal(w.chain, task.block, w.prepareResultCh, stopCh, w.prepareCompleteCh); err != nil {
					log.Warn("Block sealing failed on bft engine", "err", err)
					w.commitWorkEnv.setCommitStatusIdle()
					w.pendingMu.Lock()
					delete(w.pendingTasks, sealHash)
					w.pendingMu.Unlock()
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
	defer w.wg.Done()
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
			for i, taskReceipt := range _receipts {
				receipt := new(types.Receipt)
				receipts[i] = receipt
				*receipt = *taskReceipt

				// add block location fields
				receipt.BlockHash = hash
				receipt.BlockNumber = block.Number()
				receipt.TransactionIndex = uint(i)

				// Update the block hash in all logs since it is now available and not when the
				// receipt/log of individual transactions were created.
				receipt.Logs = make([]*types.Log, len(taskReceipt.Logs))
				for i, taskLog := range taskReceipt.Logs {
					log := new(types.Log)
					receipt.Logs[i] = log
					*log = *taskLog
					log.BlockHash = hash
				}
				logs = append(logs, receipt.Logs...)
			}
			// Commit block and state to database.
			block.SetExtraData(cbftResult.ExtraData)
			log.Debug("Write extra data", "txs", len(block.Transactions()), "extra", len(block.ExtraData()))
			// update 3-chain state
			cbftResult.ChainStateUpdateCB()
			_, err := w.chain.WriteBlockWithState(block, receipts, logs, _state, true)
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
			if !w.engine.Syncing() {
				log.Trace("Broadcast the block and announce chain insertion event", "hash", block.Hash(), "number", block.NumberU64())
				w.mux.Post(core.NewMinedBlockEvent{Block: block})
			}

		case <-w.exitCh:
			return
		}
	}
}

// makeEnv creates a new environment for the sealing block.
func (w *worker) makeEnv(parent *types.Block, header *types.Header, generateExtra func(*state.StateDB) []byte) (*environment, error) {
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
		return nil, err
	}

	extra := generateExtra(state)
	copy(header.Extra[:len(extra)], extra)
	log.Debug("Prepare header extra", "data", hex.EncodeToString(extra), "num", header.Number)

	env := &environment{
		signer:     types.MakeSigner(w.chainConfig, header.Number, gov.Gte150VersionState(state)),
		snapshotDB: snapshotdb.Instance(),
		state:      state,
		header:     header,
		coinbase:   header.Coinbase,
	}

	// Keep track of transactions which return errors so they can be removed
	env.tcount = 0
	return env, nil
}

// updateSnapshot updates pending snapshot block, receipts and state.
func (w *worker) updateSnapshot(env *environment) {
	w.snapshotMu.Lock()
	defer w.snapshotMu.Unlock()

	w.snapshotBlock = types.NewBlock(
		env.header,
		env.txs,
		env.receipts,
		trie.NewStackTrie(nil),
	)
	w.snapshotReceipts = copyReceipts(env.receipts)
	w.snapshotState = env.state.Copy()
}

func (w *worker) commitTransaction(env *environment, tx *types.Transaction) ([]*types.Log, error) {
	snapForSnap, snapForState := env.DBSnapshot()

	vmCfg := *w.chain.GetVMConfig()       // value copy
	vmCfg.VmTimeoutDuration = w.vmTimeout // set vm execution smart contract timeout duration
	receipt, err := core.ApplyTransaction(w.chainConfig, w.chain, env.gasPool, env.state,
		env.header, tx, &env.header.GasUsed, vmCfg)
	if err != nil {
		log.Error("Failed to commitTransaction on worker", "blockNumer", env.header.Number.Uint64(), "txHash", tx.Hash().String(), "err", err)
		env.RevertToDBSnapshot(snapForSnap, snapForState)
		return nil, err
	}
	env.txs = append(env.txs, tx)
	env.receipts = append(env.receipts, receipt)

	return receipt.Logs, nil
}

// generateParams wraps various of settings for generating sealing task.
type generateParams struct {
	timestamp uint64       // The timstamp for sealing task
	parent    *types.Block // Parent block hash, empty means the latest chain head
}

// prepareWork constructs the sealing task according to the given parameters,
// either based on the last chain head or specified parent. In this function
// the pending transactions are not filled yet, only the empty task returned.
func (w *worker) prepareWork(genParams *generateParams) (*environment, error) {
	// Find the parent block for sealing task
	var (
		parent    *types.Block
		timestamp = genParams.timestamp
	)
	if _, ok := w.engine.(consensus.Bft); ok {
		parent = genParams.parent
	} else {
		parent = w.chain.CurrentBlock()

		// Sanity check the timestamp correctness, recap the timestamp
		// to parent+1 if the mutation is allowed.
		if parent.Time() >= timestamp {
			timestamp = parent.Time() + 1
		}
	}

	// Construct the sealing block header, set the extra field if it's allowed
	num := parent.Number()
	header := &types.Header{
		ParentHash: parent.Hash(),
		Number:     num.Add(num, common.Big1),
		GasLimit:   core.CalcGasLimit(parent, w.config.GasFloor, snapshotdb.Instance()),
		Time:       timestamp,
	}

	log.Info("Cbft begin to consensus for new block", "number", header.Number, "nonce", hexutil.Encode(header.Nonce[:]), "gasLimit", header.GasLimit, "parentHash", parent.Hash(), "parentNumber", parent.NumberU64(), "parentStateRoot", parent.Root(), "timestamp", common.MillisToString(int64(timestamp)))

	// Initialize the header extra in Prepare function of engine
	if err := w.engine.Prepare(w.chain, header); err != nil {
		log.Error("Failed to prepare header for mining", "err", err)
		return nil, err
	}

	// Set baseFee and GasLimit if we are on an EIP-1559 chain
	if w.chainConfig.IsPauli(header.Number) {
		header.BaseFee = misc.CalcBaseFee(w.chainConfig, parent.Header())
		parentGasLimit := parent.GasLimit()
		if !w.chainConfig.IsPauli(parent.Number()) {
			// Bump by 2x
			parentGasLimit = parent.GasLimit() * params.ElasticityMultiplier
		}
		gasCeil := core.CalcGasCeil(parent, snapshotdb.Instance())
		header.GasLimit = core.CalcGasLimit1559(parentGasLimit, gasCeil)
	}

	// BeginBlocker()
	reactor := core.GetReactorInstance()

	// Only set the coinbase if our consensus engine is running (avoid spurious block rewards)
	if w.isRunning() {
		if b, ok := w.engine.(consensus.Bft); ok {
			reactor.SetWorkerCoinBase(header, b.Node().IDv0())
		}
	}

	if err := reactor.PrepareHeaderNonce(header); err != nil {
		return nil, err
	}

	env, err := w.makeEnv(parent, header, func(state *state.StateDB) []byte {
		// create default extradata
		extra, _ := rlp.EncodeToBytes([]interface{}{
			//uint(params.VersionMajor<<16 | params.VersionMinor<<8 | params.VersionPatch),
			gov.GetCurrentActiveVersion(state),
			"platon",
			runtime.Version(),
			runtime.GOOS,
		})
		if uint64(len(extra)) > params.MaximumExtraDataSize {
			log.Warn("Miner extra data exceed limit", "extra", hexutil.Bytes(extra), "limit", params.MaximumExtraDataSize)
			extra = nil
		}

		return extra
	})
	if err != nil {
		log.Error("Failed to create sealing context", "err", err)
		return nil, err
	}

	return env, nil
}

// fillTransactions retrieves the pending transactions from the txpool and fills them
// into the given sealing block. The transaction selection and ordering strategy can
// be customized with the plugin in the future.
func (w *worker) fillTransactions(interrupt *int32, env *environment, timestamp int64, blockDeadline time.Time) error {
	// Split the pending transactions into locals and remotes
	// Fill the block with all available pending transactions.
	pending := w.eth.TxPool().Pending(true, true)

	localTxs, remoteTxs := make(map[common.Address]types.Transactions), pending
	for _, account := range w.eth.TxPool().Locals() {
		if txs := remoteTxs[account]; len(txs) > 0 {
			delete(remoteTxs, account)
			localTxs[account] = txs
		}
	}

	var (
		localTimeout      = false
		tempContractCache = make(map[common.Address]struct{})
	)

	if len(localTxs) > 0 {
		txs := types.NewTransactionsByPriceAndNonce(env.signer, localTxs, env.header.BaseFee)
		if failed, timeout := w.committer.CommitTransactions(env, txs, interrupt, timestamp, blockDeadline, tempContractCache); failed {
			return fmt.Errorf("commit transactions error")
		} else {
			localTimeout = timeout
		}
	}

	if !localTimeout && len(remoteTxs) > 0 {
		txs := types.NewTransactionsByPriceAndNonce(env.signer, remoteTxs, env.header.BaseFee)
		if failed, _ := w.committer.CommitTransactions(env, txs, interrupt, timestamp, blockDeadline, tempContractCache); failed {
			return fmt.Errorf("commit transactions error")
		}
	}

	return nil
}

// commitWork generates several new sealing tasks based on the parent block
// and submit them to the sealer.
func (w *worker) commitWork(interrupt *int32, noempty bool, timestamp int64, commitBlock *types.Block, blockDeadline time.Time) error {
	w.mu.RLock()
	defer w.mu.RUnlock()
	atomic.StoreInt32(&w.commitWorkEnv.commitStatus, commitStatusCommitting)
	defer func() {
		if engine, ok := w.engine.(consensus.Bft); ok {
			w.commitWorkEnv.nextBlockTime.Store(engine.CalcNextBlockTime(common.MillisToTime(timestamp)))
			log.Debug("Next block time", "time", common.Beautiful(w.commitWorkEnv.nextBlockTime.Load().(time.Time)))
		}
	}()
	start := time.Now()
	work, err := w.prepareWork(&generateParams{
		timestamp: uint64(timestamp),
		parent:    commitBlock,
	})
	if err != nil {
		return err
	}
	if err := core.GetReactorInstance().NewBlock(work.header, work.state, common.ZeroHash); err != nil {
		log.Error("Failed to GetReactorInstance BeginBlocker on worker", "blockNumber", work.header.Number, "err", err)
		return err
	}

	// Create an empty block based on temporary copied state for
	// sealing in advance without waiting block execution finished.
	if !noempty && "on" == w.EmptyBlock {
		// Create an empty block based on temporary copied state for sealing in advance without waiting block
		// execution finished.
		if _, ok := w.engine.(consensus.Bft); !ok {
			if err := w.commit(work.copy(), nil, false, start); nil != err {
				log.Error("Failed to commitNewWork on worker: call commit is failed", "blockNumber", work.header.Number, "err", err)
			}
		}
	}
	log.Trace("Validator mode", "mode", w.chainConfig.Cbft.ValidatorMode)
	if w.chainConfig.Cbft.ValidatorMode == "inner" {
		// Check if need to switch validators.
		// If needed, make a inner contract transaction
		// and pack into pending block.
		if w.shouldSwitch(work) && w.commitInnerTransaction(work, timestamp, blockDeadline) != nil {
			return fmt.Errorf("commit inner transaction error")
		}
	}

	// Fill pending transactions from the txpool
	if err := w.fillTransactions(interrupt, work, timestamp, blockDeadline); err != nil {
		return err
	}
	if err := w.commit(work.copy(), w.fullTaskHook, true, start); err != nil {
		log.Error("Failed to commitNewWork on worker: call commit is failed", "blockNumber", work.header.Number, "err", err)
		return err
	}

	// Swap out the old work with the new one, terminating any leftover
	// prefetcher processes in the mean time and starting a new one.
	if w.current != nil {
		w.current.discard()
	}
	w.current = work
	log.Info("Commit new work", "number", work.header.Number, "txs", work.tcount, "duration", time.Since(start))
	return nil
}

// commit runs any post-transaction state modifications, assembles the final block
// and commits new work if consensus engine is running.
func (w *worker) commit(env *environment, interval func(), update bool, start time.Time) error {

	// EndBlocker()
	if err := core.GetReactorInstance().EndBlocker(env.header, env.state); nil != err {
		log.Error("Failed to GetReactorInstance EndBlocker on worker", "blockNumber",
			env.header.Number.Uint64(), "err", err)
		return err
	}

	block, err := w.engine.Finalize(w.chain, env.header, env.state, env.txs, env.receipts)

	if err != nil {
		return err
	}
	if update {
		w.updateSnapshot(env)
	}
	if w.isRunning() {
		if interval != nil {
			interval()
		}
		select {
		case w.taskCh <- &task{receipts: env.receipts, state: env.state, block: block, createdAt: time.Now()}:
			w.unconfirmed.Shift(block.NumberU64() - 1)
			log.Info("Commit new mining work", "number", block.Number(), "sealhash", w.engine.SealHash(block.Header()), "receiptHash", block.ReceiptHash(),
				"txs", env.tcount, "gas", block.GasUsed(), "fees", totalFees(block, env.receipts), "elapsed", common.PrettyDuration(time.Since(start)))
		case <-w.exitCh:
			log.Info("Worker has exited")
		}
		return nil
	}
	w.commitWorkEnv.setCommitStatusIdle()
	return nil
}

// copyReceipts makes a deep copy of the given receipts.
func copyReceipts(receipts []*types.Receipt) []*types.Receipt {
	result := make([]*types.Receipt, len(receipts))
	for i, l := range receipts {
		cpy := *l
		result[i] = &cpy
	}
	return result
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
			block := types.NewBlock(parent.Header(), parent.Transactions(), nil, new(trie.Trie))

			return block, state
		}
	}
	return nil, nil
}

// totalFees computes total consumed fees in ETH. Block transactions and receipts have to have the same order.
func totalFees(block *types.Block, receipts []*types.Receipt) *big.Float {
	feesWei := new(big.Int)
	for i, tx := range block.Transactions() {
		minerFee, _ := tx.EffectiveGasTip(block.BaseFee())
		feesWei.Add(feesWei, new(big.Int).Mul(new(big.Int).SetUint64(receipts[i].GasUsed), minerFee))
	}
	return new(big.Float).Quo(new(big.Float).SetInt(feesWei), new(big.Float).SetInt(big.NewInt(params.LAT)))
}

func (w *worker) shouldCommit(timestamp time.Time) (bool, *types.Block) {
	currentBaseBlock := w.commitWorkEnv.getCurrentBaseBlock()
	nextBaseBlock := w.engine.NextBaseBlock()
	nextBaseBlockTime := common.MillisToTime(int64(nextBaseBlock.Time()))

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
				"next.timestamp", common.MillisToString(int64(nextBaseBlock.Time())),
				"nextBlockTime", nextBlockTime,
				"timestamp", timestamp)
		} else {
			log.Debug("Check if time's up in shouldCommit()", "result", shouldCommit,
				"current.number", currentBaseBlock.Number(),
				"current.hash", currentBaseBlock.Hash(),
				"current.timestamp", common.MillisToString(int64(currentBaseBlock.Time())),
				"next.number", nextBaseBlock.Number(),
				"next.hash", nextBaseBlock.Hash(),
				"next.timestamp", common.MillisToString(int64(nextBaseBlock.Time())),
				"nextBlockTime", nextBlockTime,
				"timestamp", timestamp)
		}
	}
	return shouldCommit, nextBaseBlock
}
