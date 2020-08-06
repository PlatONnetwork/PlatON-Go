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

package miner

import (
	"math/big"
	"testing"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"

	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/core/state"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/event"
	"github.com/PlatONnetwork/PlatON-Go/params"
	_ "github.com/PlatONnetwork/PlatON-Go/x/xcom"
)

var (
	// Test chain configurations
	testTxPoolConfig core.TxPoolConfig
	chainConfig      *params.ChainConfig

	// Test accounts
	testBankKey, _  = crypto.GenerateKey()
	testBankAddress = crypto.PubkeyToAddress(testBankKey.PublicKey)
	testBankFunds   = big.NewInt(1000000000000000000)

	testUserKey, _  = crypto.GenerateKey()
	testUserAddress = crypto.PubkeyToAddress(testUserKey.PublicKey)

	// Test transactions
	pendingTxs []*types.Transaction
	newTxs     []*types.Transaction
)

func init() {
	testTxPoolConfig = core.DefaultTxPoolConfig
	testTxPoolConfig.Journal = ""
	chainConfig = params.TestChainConfig
	chainConfig.Clique = &params.CliqueConfig{
		Period: 10,
		Epoch:  30000,
	}
	tx1, _ := types.SignTx(types.NewTransaction(0, testUserAddress, big.NewInt(1000), params.TxGas, nil, nil), types.NewEIP155Signer(chainConfig.ChainID), testBankKey)
	pendingTxs = append(pendingTxs, tx1)
	tx2, _ := types.SignTx(types.NewTransaction(1, testUserAddress, big.NewInt(1000), params.TxGas, nil, nil), types.NewEIP155Signer(chainConfig.ChainID), testBankKey)
	newTxs = append(newTxs, tx2)

}

// testWorkerBackend implements worker.Backend interfaces and wraps all information needed during the testing.
type testWorkerBackend struct {
	db         ethdb.Database
	txPool     *core.TxPool
	chain      *core.BlockChain
	testTxFeed event.Feed
	chainCache *core.BlockChainCache
	engine     consensus.Engine
}

func newTestWorkerBackend(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine, n int, mux *event.TypeMux) *testWorkerBackend {
	var (
		db    = rawdb.NewMemoryDatabase()
		gspec = core.Genesis{
			Config: chainConfig,
			Alloc:  core.GenesisAlloc{testBankAddress: {Balance: testBankFunds}},
		}
	)

	switch engine.(type) {
	case consensus.Bft:
		gspec.ExtraData = make([]byte, 32+common.AddressLength+65)
		copy(gspec.ExtraData[32:], testBankAddress[:])
	default:
		t.Fatalf("unexpected consensus engine type: %T", engine)
	}
	genesis := gspec.MustCommit(db)

	engine.InsertChain(genesis)
	bft := engine.(*consensus.BftMock)
	bft.EventMux = mux
	chain, _ := core.NewBlockChain(db, nil, gspec.Config, engine, vm.Config{}, nil)
	blockChainCache := core.NewBlockChainCache(chain)

	stateDB, _ := state.New(genesis.Root(), state.NewDatabase(db))

	blockChainCache.WriteStateDB(genesis.Header().SealHash(), stateDB, 0)

	txpool := core.NewTxPool(testTxPoolConfig, chainConfig, blockChainCache)

	// Generate a small n-block chain and an uncle block for it
	if n > 0 {
		blocks, _ := core.GenerateChain(chainConfig, genesis, engine, db, n, func(i int, gen *core.BlockGen) {
			gen.SetCoinbase(testBankAddress)
		})
		if _, err := chain.InsertChain(blocks); err != nil {
			t.Fatalf("failed to insert origin chain: %v", err)
		}
	}
	parent := genesis
	if n > 0 {
		parent = chain.GetBlockByHash(chain.CurrentBlock().ParentHash())
	}
	core.GenerateChain(chainConfig, parent, engine, db, 1, func(i int, gen *core.BlockGen) {
		gen.SetCoinbase(testUserAddress)
	})

	return &testWorkerBackend{
		db:         db,
		chain:      chain,
		txPool:     txpool,
		chainCache: blockChainCache,
		engine:     engine,
	}
}

func (b *testWorkerBackend) BlockChain() *core.BlockChain { return b.chain }
func (b *testWorkerBackend) TxPool() *core.TxPool         { return b.txPool }
func (b *testWorkerBackend) PostChainEvents(events []interface{}) {
	b.chain.PostChainEvents(events, nil)
}

func newTestWorker(t *testing.T, chainConfig *params.ChainConfig, miningConfig *core.MiningConfig, engine consensus.Engine, blocks int) (*worker, *testWorkerBackend) {

	event := new(event.TypeMux)
	backend := newTestWorkerBackend(t, chainConfig, engine, blocks, event)
	core.NewExecutor(chainConfig, backend.chain, vm.Config{}, nil)

	bftResultSub := event.Subscribe(cbfttypes.CbftResult{})
	core.NewBlockChainReactor(event, chainConfig.ChainID)
	w := newWorker(chainConfig, miningConfig, &vm.Config{}, engine, backend, event, time.Second, params.GenesisGasLimit, nil, backend.chainCache, 0)
	go func() {

		for {
			select {
			case obj := <-bftResultSub.Chan():

				if obj == nil {
					continue
				}
				cbftResult, ok := obj.Data.(cbfttypes.CbftResult)
				if !ok {
					log.Error("blockchain_reactor receive bft result type error")
					continue
				}

				stateDB, err := w.blockChainCache.MakeStateDB(cbftResult.Block)
				if nil != err {
					panic(err)
				}

				// block write to real chain
				_, err = w.chain.WriteBlockWithState(cbftResult.Block, nil, stateDB)
				if nil != err {
					panic(err)
				}

				// block write to BftMock engine chain
				backend.engine.InsertChain(cbftResult.Block)
			}
		}

	}()
	return w, backend
}

//func TestEmptyWorkCbft(t *testing.T) {
//	testEmptyWork(t, chainConfig, consensus.NewFaker())
//}

//func testEmptyWork(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine) {
//	defer engine.Close()
//
//	minningConfig := &core.MiningConfig{
//		MiningLogAtDepth:       7,
//		TxChanSize:             4096,
//		ChainHeadChanSize:      10,
//		ChainSideChanSize:      10,
//		ResultQueueSize:        10,
//		ResubmitAdjustChanSize: 10,
//		MinRecommitInterval:    1 * time.Second,
//		MaxRecommitInterval:    15 * time.Second,
//		IntervalAdjustRatio:    0.1,
//		IntervalAdjustBias:     200 * 1000.0 * 1000.0,
//		StaleThreshold:         7,
//		DefaultCommitRatio:     0.95,
//	}
//	w, _ := newTestWorker(t, chainConfig, minningConfig, engine, 0)
//
//	defer w.close()
//
//	var (
//		taskCh    = make(chan struct{}, 2)
//		taskIndex int
//	)
//
//	checkEqual := func(t *testing.T, task *task, index int) {
//		receiptLen, balance := 0, big.NewInt(0)
//		if index == 1 {
//			receiptLen, balance = 1, big.NewInt(1000)
//		}
//		if len(task.receipts) != receiptLen {
//			t.Errorf("receipt number mismatch: have %d, want %d", len(task.receipts), receiptLen)
//		}
//		if task.state.GetBalance(testUserAddress).Cmp(balance) != 0 {
//			t.Errorf("account balance mismatch: have %d, want %d", task.state.GetBalance(testUserAddress), balance)
//		}
//	}
//
//	w.newTaskHook = func(task *task) {
//		if task.block.NumberU64() == 1 {
//			checkEqual(t, task, taskIndex)
//			taskIndex += 1
//			taskCh <- struct{}{}
//		}
//	}
//	w.fullTaskHook = func() {
//		time.Sleep(1000 * time.Millisecond)
//	}
//
//	// Ensure worker has finished initialization
//	go func() {
//		for {
//			b := w.pendingBlock()
//			if b != nil && b.NumberU64() == 1 {
//				break
//			}
//		}
//	}()
//
//	go w.start()
//	go func() {
//		for i := 0; i < 2; i += 1 {
//			select {
//			case <-taskCh:
//			case <-time.NewTimer(2 * time.Second).C:
//				t.Error("new task timeout")
//			}
//		}
//	}()
//}

func TestPendingStateAndBlockCbft(t *testing.T) {
	testPendingStateAndBlock(t, chainConfig, consensus.NewFaker())
}

func testPendingStateAndBlock(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine) {
	defer engine.Close()

	minningConfig := &core.MiningConfig{
		MiningLogAtDepth:       7,
		TxChanSize:             4096,
		ChainHeadChanSize:      10,
		ChainSideChanSize:      10,
		ResultQueueSize:        10,
		ResubmitAdjustChanSize: 10,
		MinRecommitInterval:    1 * time.Second,
		MaxRecommitInterval:    15 * time.Second,
		IntervalAdjustRatio:    0.1,
		IntervalAdjustBias:     200 * 1000.0 * 1000.0,
		StaleThreshold:         7,
		DefaultCommitRatio:     0.95,
	}

	w, b := newTestWorker(t, chainConfig, minningConfig, engine, 0)
	go w.start()
	defer w.close()

	// First add pendingTxs to the trading pool
	b.txPool.AddLocals(pendingTxs)

	// Ensure snapshot has been updated.
	time.Sleep(100 * time.Millisecond)
	block, state := w.pending()
	if block.NumberU64() < 1 {
		t.Errorf("block number mismatch: have %d, want >= %d", block.NumberU64(), 1)
	}
	if balance := state.GetBalance(testUserAddress); balance.Cmp(big.NewInt(1000)) != 0 {
		t.Errorf("account balance mismatch: have %d, want %d", balance, 1000)
	}

	// Add newTxs to the trading pool
	b.txPool.AddLocals(newTxs)

	// Ensure the new tx events has been processed
	time.Sleep(100 * time.Millisecond)
	block, state = w.pending()
	if balance := state.GetBalance(testUserAddress); balance.Cmp(big.NewInt(2000)) != 0 {
		t.Errorf("account balance mismatch: have %d, want %d", balance, 2000)
	}
}

func TestRegenerateMiningBlockCbft(t *testing.T) {
	testRegenerateMiningBlock(t, chainConfig, consensus.NewFaker())
}

func testRegenerateMiningBlock(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine) {
	defer engine.Close()

	minningConfig := &core.MiningConfig{
		MiningLogAtDepth:       7,
		TxChanSize:             4096,
		ChainHeadChanSize:      10,
		ChainSideChanSize:      10,
		ResultQueueSize:        10,
		ResubmitAdjustChanSize: 10,
		MinRecommitInterval:    1 * time.Second,
		MaxRecommitInterval:    15 * time.Second,
		IntervalAdjustRatio:    0.1,
		IntervalAdjustBias:     200 * 1000.0 * 1000.0,
		StaleThreshold:         7,
		DefaultCommitRatio:     0.95,
	}
	w, b := newTestWorker(t, chainConfig, minningConfig, engine, 0)
	b.txPool.AddLocals(pendingTxs)
	b.txPool.AddLocals(newTxs)
	defer w.close()

	var taskCh = make(chan struct{})

	taskIndex := 0
	w.newTaskHook = func(task *task) {
		if task.block.NumberU64() == 1 {
			if taskIndex == 2 {
				receiptLen, balance := 2, big.NewInt(2000)
				if len(task.receipts) != receiptLen {
					t.Errorf("receipt number mismatch: have %d, want %d", len(task.receipts), receiptLen)
				}
				if task.state.GetBalance(testUserAddress).Cmp(balance) != 0 {
					t.Errorf("account balance mismatch: have %d, want %d", task.state.GetBalance(testUserAddress), balance)
				}
			}
			taskCh <- struct{}{}
			taskIndex += 1
		}
	}
	w.skipSealHook = func(task *task) bool {
		return true
	}
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}
	// Ensure worker has finished initialization
	go func() {
		for {
			b := w.pendingBlock()
			if b != nil && b.NumberU64() == 1 {
				break
			}
		}
	}()

	go w.start()
	// Ignore the first two works
	go func() {
		for i := 0; i < 2; i += 1 {
			select {
			case <-taskCh:
			case <-time.NewTimer(time.Second).C:
				t.Error("new task timeout")
			}
		}
	}()
	b.txPool.AddLocals(newTxs)
	time.Sleep(time.Second)

	go func() {
		select {
		case <-taskCh:
		case <-time.NewTimer(time.Second).C:
			t.Error("new task timeout")
		}
	}()
}

func TestAdjustIntervalCbft(t *testing.T) {
	testAdjustInterval(t, chainConfig, consensus.NewFaker())
}

func testAdjustInterval(t *testing.T, chainConfig *params.ChainConfig, engine consensus.Engine) {
	defer engine.Close()

	minningConfig := &core.MiningConfig{
		MiningLogAtDepth:       7,
		TxChanSize:             4096,
		ChainHeadChanSize:      10,
		ChainSideChanSize:      10,
		ResultQueueSize:        10,
		ResubmitAdjustChanSize: 10,
		MinRecommitInterval:    1 * time.Second,
		MaxRecommitInterval:    15 * time.Second,
		IntervalAdjustRatio:    0.1,
		IntervalAdjustBias:     200 * 1000.0 * 1000.0,
		StaleThreshold:         7,
		DefaultCommitRatio:     0.95,
	}

	w, _ := newTestWorker(t, chainConfig, minningConfig, engine, 0)
	defer w.close()

	w.skipSealHook = func(task *task) bool {
		return true
	}
	w.fullTaskHook = func() {
		time.Sleep(100 * time.Millisecond)
	}
	var (
		progress = make(chan struct{}, 10)
		result   = make([]float64, 0, 10)
		index    = 0
		start    = false
	)
	w.resubmitHook = func(minInterval time.Duration, recommitInterval time.Duration) {
		// Short circuit if interval checking hasn't started.
		if !start {
			return
		}
		var wantMinInterval, wantRecommitInterval time.Duration

		switch index {
		case 0:
			wantMinInterval, wantRecommitInterval = 3*time.Second, 3*time.Second
		case 1:
			origin := float64(3 * time.Second.Nanoseconds())
			estimate := origin*(1-minningConfig.IntervalAdjustRatio) + minningConfig.IntervalAdjustRatio*(origin/0.8+minningConfig.IntervalAdjustBias)
			wantMinInterval, wantRecommitInterval = 3*time.Second, time.Duration(int(estimate))*time.Nanosecond
		case 2:
			estimate := result[index-1]
			min := float64(3 * time.Second.Nanoseconds())
			estimate = estimate*(1-minningConfig.IntervalAdjustRatio) + minningConfig.IntervalAdjustRatio*(min-minningConfig.IntervalAdjustBias)
			wantMinInterval, wantRecommitInterval = 3*time.Second, time.Duration(int(estimate))*time.Nanosecond
		case 3:
			wantMinInterval, wantRecommitInterval = time.Second, time.Second
		}

		// Check interval
		if minInterval != wantMinInterval {
			t.Errorf("resubmit min interval mismatch: have %v, want %v ", minInterval, wantMinInterval)
		}
		if recommitInterval != wantRecommitInterval {
			t.Errorf("resubmit interval mismatch: have %v, want %v", recommitInterval, wantRecommitInterval)
		}
		result = append(result, float64(recommitInterval.Nanoseconds()))
		index += 1
		progress <- struct{}{}
	}
	// Ensure worker has finished initialization
	go func() {
		for {
			b := w.pendingBlock()
			if b != nil && b.NumberU64() == 1 {
				break
			}
		}
	}()

	w.start()

	time.Sleep(time.Second)

	start = true
	w.setRecommitInterval(3 * time.Second)
	go func() {
		select {
		case <-progress:
		case <-time.NewTimer(time.Second).C:
			t.Error("interval reset timeout")
		}
	}()

	w.resubmitAdjustCh <- &intervalAdjust{inc: true, ratio: 0.8}
	go func() {

		select {
		case <-progress:
		case <-time.NewTimer(time.Second).C:
			t.Error("interval reset timeout")
		}
	}()

	w.resubmitAdjustCh <- &intervalAdjust{inc: false}
	go func() {
		select {
		case <-progress:
		case <-time.NewTimer(time.Second).C:
			t.Error("interval reset timeout")
		}
	}()

	w.setRecommitInterval(500 * time.Millisecond)
	go func() {
		select {
		case <-progress:
		case <-time.NewTimer(time.Second).C:
			t.Error("interval reset timeout")
		}
	}()
}
