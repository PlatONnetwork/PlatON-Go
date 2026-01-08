// Copyright 2021 The go-ethereum Authors
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

package tracers

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"reflect"
	"sort"
	"sync/atomic"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/eth/tracers/logger"
	"github.com/PlatONnetwork/PlatON-Go/ethdb"
	"github.com/PlatONnetwork/PlatON-Go/internal/ethapi"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
)

func init() {
	params.GetEc(params.DefaultUnitTestNet)
}

var (
	errStateNotFound = errors.New("state not found")
	errBlockNotFound = errors.New("block not found")
)

type Result struct {
	Result interface{} `json:"result,omitempty"`
}
type testBackend struct {
	chainConfig *params.ChainConfig
	engine      consensus.Engine
	chaindb     ethdb.Database
	chain       *core.BlockChain

	refHook func() // Hook is invoked when the requested state is referenced
	relHook func() // Hook is invoked when the requested state is released
}

// testBackend creates a new test backend. OBS: After test is done, teardown must be
// invoked in order to release associated resources.
func newTestBackend(t *testing.T, n int, gspec *core.Genesis, generator func(i int, b *core.BlockGen)) *testBackend {
	var gendb = rawdb.NewMemoryDatabase()
	var genesis = gspec.MustCommit(gendb)
	backend := &testBackend{
		chainConfig: gspec.Config,
		engine:      consensus.NewFakerWithDataBase(gendb, genesis),
		chaindb:     gendb,
	}
	// Generate blocks for testing
	//gspec.Config = backend.chainConfig
	//blocks, _ := core.GenerateChain(backend.chainConfig, genesis, backend.engine, gendb, n, generator)

	// Import the canonical chain
	//gspec.MustCommit(backend.chaindb)
	//cacheConfig := &core.CacheConfig{
	//	TrieCleanLimit: 256,
	//	TrieDirtyLimit: 256,
	//	TrieTimeLimit:  5 * time.Minute,
	//	SnapshotLimit:  0,
	//}
	/*chain, err := core.NewBlockChain(backend.chaindb, nil, backend.chainConfig, backend.engine, vm.Config{}, nil, nil)
	if err != nil {
		t.Fatalf("failed to create tester chain: %v", err)
	}
	if n, err := chain.InsertChain(blocks); err != nil {
		t.Fatalf("block %d: failed to insert into chain: %v", n, err)
	}*/
	backend.chain, _ = core.GenerateBlockChain2(gspec, genesis, backend.engine, gendb, n, generator)
	return backend
}

func (b *testBackend) HeaderByHash(ctx context.Context, hash common.Hash) (*types.Header, error) {
	return b.chain.GetHeaderByHash(hash), nil
}

func (b *testBackend) HeaderByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Header, error) {
	if number == rpc.PendingBlockNumber || number == rpc.LatestBlockNumber {
		return b.chain.CurrentHeader(), nil
	}
	return b.chain.GetHeaderByNumber(uint64(number)), nil
}

func (b *testBackend) BlockByHash(ctx context.Context, hash common.Hash) (*types.Block, error) {
	return b.chain.GetBlockByHash(hash), nil
}

func (b *testBackend) BlockByNumber(ctx context.Context, number rpc.BlockNumber) (*types.Block, error) {
	if number == rpc.PendingBlockNumber || number == rpc.LatestBlockNumber {
		return b.chain.CurrentBlock(), nil
	}
	return b.chain.GetBlockByNumber(uint64(number)), nil
}

func (b *testBackend) GetTransaction(ctx context.Context, txHash common.Hash) (*types.Transaction, common.Hash, uint64, uint64, error) {
	tx, hash, blockNumber, index := rawdb.ReadTransaction(b.chaindb, txHash)
	return tx, hash, blockNumber, index, nil
}

func (b *testBackend) RPCGasCap() uint64 {
	return 25000000
}

func (b *testBackend) ChainConfig() *params.ChainConfig {
	return b.chainConfig
}

func (b *testBackend) Engine() consensus.Engine {
	return b.engine
}

func (b *testBackend) ChainDb() ethdb.Database {
	return b.chaindb
}

// teardown releases the associated resources.
func (b *testBackend) teardown() {
	b.chain.Stop()
}

func (b *testBackend) StateAtBlock(ctx context.Context, block *types.Block, reexec uint64, base *state.StateDB, readOnly bool, preferDisk bool) (*state.StateDB, snapshotdb.DB, StateReleaseFunc, error) {
	statedb, err := b.chain.StateAt(block.Root())
	if err != nil {
		return nil, nil, nil, errStateNotFound
	}
	if b.refHook != nil {
		b.refHook()
	}
	release := func() {
		if b.relHook != nil {
			b.relHook()
		}
	}
	return statedb, snapshotdb.Instance(), release, nil
}

func (b *testBackend) StateAtTransaction(ctx context.Context, block *types.Block, txIndex int, reexec uint64) (core.Message, vm.BlockContext, *state.StateDB, snapshotdb.DB, StateReleaseFunc, error) {
	parent := b.chain.GetBlock(block.ParentHash(), block.NumberU64()-1)
	if parent == nil {
		return nil, vm.BlockContext{}, nil, nil, nil, errBlockNotFound
	}
	statedb, _, release, err := b.StateAtBlock(ctx, parent, reexec, nil, true, false)
	if err != nil {
		return nil, vm.BlockContext{}, nil, nil, nil, errStateNotFound
	}
	if txIndex == 0 && len(block.Transactions()) == 0 {
		return nil, vm.BlockContext{}, statedb, snapshotdb.Instance(), release, nil
	}
	// Recompute transactions up to the target index.
	signer := types.MakeSigner(b.chainConfig, block.Number(), true)
	for idx, tx := range block.Transactions() {
		msg, _ := tx.AsMessage(signer, block.BaseFee())
		txContext := core.NewEVMTxContext(msg)
		context := core.NewEVMBlockContext(block.Header(), b.chain)
		if idx == txIndex {
			return msg, context, statedb, snapshotdb.Instance(), release, nil
		}
		vmenv := vm.NewEVM(context, txContext, snapshotdb.Instance(), statedb, b.chainConfig, vm.Config{})
		if _, err := core.ApplyMessage(vmenv, msg, new(core.GasPool).AddGas(tx.Gas())); err != nil {
			return nil, vm.BlockContext{}, nil, snapshotdb.Instance(), nil, fmt.Errorf("transaction %#x failed: %v", tx.Hash(), err)
		}
		statedb.Finalise(false)
	}
	return nil, vm.BlockContext{}, nil, snapshotdb.Instance(), nil, fmt.Errorf("transaction index %d out of range for block %#x", txIndex, block.Hash())
}

func TestTraceCall(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &core.Genesis{
		Config: params.TestChainConfig,
		Alloc: core.GenesisAlloc{
			accounts[0].addr: {Balance: big.NewInt(params.LAT)},
			accounts[1].addr: {Balance: big.NewInt(params.LAT)},
			accounts[2].addr: {Balance: big.NewInt(params.LAT)},
		},
		BaseFee: big.NewInt(params.InitialBaseFee),
	}
	genBlocks := 10
	signer := types.MakeSigner(params.TestChainConfig, new(big.Int).SetUint64(1), true)
	backend := newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, b.BaseFee(), nil), signer, accounts[0].key)
		b.AddTx(tx)
	})
	defer backend.teardown()
	api := NewAPI(backend)
	var testSuite = []struct {
		blockNumber rpc.BlockNumber
		call        ethapi.TransactionArgs
		config      *TraceCallConfig
		expectErr   error
		expect      string
	}{
		// Standard JSON trace upon the genesis, plain transfer.
		{
			blockNumber: rpc.BlockNumber(0),
			call: ethapi.TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect:    `{"gas":21000,"failed":false,"returnValue":"","structLogs":[]}`,
		},
		// Standard JSON trace upon the head, plain transfer.
		{
			blockNumber: rpc.BlockNumber(genBlocks),
			call: ethapi.TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect:    `{"gas":21000,"failed":false,"returnValue":"","structLogs":[]}`,
		},
		// Standard JSON trace upon the non-existent block, error expects
		{
			blockNumber: rpc.BlockNumber(genBlocks + 1),
			call: ethapi.TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: fmt.Errorf("block #%d not found", genBlocks+1),
			//expect:    nil,
		},
		// Standard JSON trace upon the latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: nil,
			expect:    `{"gas":21000,"failed":false,"returnValue":"","structLogs":[]}`,
		},
		// Tracing on 'pending' should fail:
		{
			blockNumber: rpc.PendingBlockNumber,
			call: ethapi.TransactionArgs{
				From:  &accounts[0].addr,
				To:    &accounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    nil,
			expectErr: errors.New("tracing on top of pending is not supported"),
		},
		{
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From:  &accounts[0].addr,
				Input: &hexutil.Bytes{0x60, 0x80, 0x60, 0x40, 0x43}, // blocknumber
			},
			config: &TraceCallConfig{
				BlockOverrides: &ethapi.BlockOverrides{Number: (*hexutil.Big)(big.NewInt(0x1337))},
			},
			expectErr: nil,
			expect:    `{"gas":53090,"failed":false,"returnValue":"","structLogs":[{"pc":0,"op":"PUSH1","gas":24946918,"gasCost":3,"depth":1,"stack":[]},{"pc":2,"op":"PUSH1","gas":24946915,"gasCost":3,"depth":1,"stack":["0x80"]},{"pc":4,"op":"NUMBER","gas":24946912,"gasCost":2,"depth":1,"stack":["0x80","0x40"]},{"pc":5,"op":"STOP","gas":24946910,"gasCost":0,"depth":1,"stack":["0x80","0x40","0x1337"]}]}`,
		},
	}
	for i, testspec := range testSuite {
		result, err := api.TraceCall(context.Background(), testspec.call, rpc.BlockNumberOrHash{BlockNumber: &testspec.blockNumber}, testspec.config)
		if testspec.expectErr != nil {
			if err == nil {
				t.Errorf("test %d: expect error %v, got nothing", i, testspec.expectErr)
				continue
			}
			if !reflect.DeepEqual(err, testspec.expectErr) {
				t.Errorf("test %d: error mismatch, want %v, git %v", i, testspec.expectErr, err)
			}
		} else {
			if err != nil {
				t.Errorf("test %d: expect no error, got %v", i, err)
				continue
			}
			var have *logger.ExecutionResult
			if err := json.Unmarshal(result.(json.RawMessage), &have); err != nil {
				t.Errorf("test %d: failed to unmarshal result %v", i, err)
			}
			var want *logger.ExecutionResult
			if err := json.Unmarshal([]byte(testspec.expect), &want); err != nil {
				t.Errorf("test %d: failed to unmarshal result %v", i, err)
			}
			if !reflect.DeepEqual(have, want) {
				t.Errorf("test %d: result mismatch, want %v, got %v", i, testspec.expect, string(result.(json.RawMessage)))
			}
		}
	}
}

func TestTraceTransaction(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(2)
	genesis := &core.Genesis{Alloc: core.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.LAT)},
		accounts[1].addr: {Balance: big.NewInt(params.LAT)},
	},
		BaseFee: big.NewInt(params.InitialBaseFee),
		Config:  params.TestChainConfig,
	}
	target := common.Hash{}
	signer := types.MakeSigner(params.TestChainConfig, new(big.Int).SetUint64(1), true)
	backend := newTestBackend(t, 1, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, b.BaseFee(), nil), signer, accounts[0].key)
		b.AddTx(tx)
		target = tx.Hash()
	})
	defer backend.chain.Stop()
	api := NewAPI(backend)
	result, err := api.TraceTransaction(context.Background(), target, nil)
	if err != nil {
		t.Errorf("Failed to trace transaction %v", err)
	}
	var have *logger.ExecutionResult
	if err := json.Unmarshal(result.(json.RawMessage), &have); err != nil {
		t.Errorf("failed to unmarshal result %v", err)
	}
	if !reflect.DeepEqual(have, &logger.ExecutionResult{
		Gas:         params.TxGas,
		Failed:      false,
		ReturnValue: "",
		StructLogs:  []logger.StructLogRes{},
	}) {
		t.Error("Transaction tracing result is different")
	}

	// Test non-existent transaction
	_, err = api.TraceTransaction(context.Background(), common.Hash{42}, nil)
	if !errors.Is(err, errTxNotFound) {
		t.Fatalf("want %v, have %v", errTxNotFound, err)
	}
}

func TestTraceBlock(t *testing.T) {
	t.Parallel()

	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &core.Genesis{
		Alloc: core.GenesisAlloc{
			accounts[0].addr: {Balance: big.NewInt(params.LAT)},
			accounts[1].addr: {Balance: big.NewInt(params.LAT)},
			accounts[2].addr: {Balance: big.NewInt(params.LAT)},
		},
		BaseFee: big.NewInt(params.InitialBaseFee),
		Config:  params.TestChainConfig,
	}
	genBlocks := 10
	signer := types.MakeSigner(params.TestChainConfig, new(big.Int).SetUint64(1), true)
	backend := newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, b.BaseFee(), nil), signer, accounts[0].key)
		b.AddTx(tx)
	})
	defer backend.chain.Stop()
	api := NewAPI(backend)

	var testSuite = []struct {
		blockNumber rpc.BlockNumber
		config      *TraceConfig
		want        string
		expectErr   error
	}{
		// Trace genesis block, expect error
		{
			blockNumber: rpc.BlockNumber(0),
			expectErr:   errors.New("genesis is not traceable"),
		},
		// Trace head block
		{
			blockNumber: rpc.BlockNumber(genBlocks),
			want:        `[{"result":{"gas":21000,"failed":false,"returnValue":"","structLogs":[]}}]`,
		},
		// Trace non-existent block
		{
			blockNumber: rpc.BlockNumber(genBlocks + 1),
			expectErr:   fmt.Errorf("block #%d not found", genBlocks+1),
		},
		// Trace latest block
		{
			blockNumber: rpc.LatestBlockNumber,
			want:        `[{"result":{"gas":21000,"failed":false,"returnValue":"","structLogs":[]}}]`,
		},
		// Trace pending block
		{
			blockNumber: rpc.PendingBlockNumber,
			want:        `[{"result":{"gas":21000,"failed":false,"returnValue":"","structLogs":[]}}]`,
		},
	}
	for i, tc := range testSuite {
		result, err := api.TraceBlockByNumber(context.Background(), tc.blockNumber, tc.config)
		if tc.expectErr != nil {
			if err == nil {
				t.Errorf("test %d, want error %v", i, tc.expectErr)
				continue
			}
			if !reflect.DeepEqual(err, tc.expectErr) {
				t.Errorf("test %d: error mismatch, want %v, get %v", i, tc.expectErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("test %d, want no error, have %v", i, err)
			continue
		}
		have, _ := json.Marshal([]*Result{{Result: result[0].Result}})
		want := tc.want
		if string(have) != want {
			t.Errorf("test %d, result mismatch, have\n%v\n, want\n%v\n", i, string(have), want)
		}
	}
}

func TestTracingWithOverrides(t *testing.T) {
	t.Parallel()
	// Initialize test accounts
	accounts := newAccounts(3)
	storageAccount := common.Address{0x13, 37}
	genesis := &core.Genesis{Alloc: core.GenesisAlloc{
		accounts[0].addr: {Balance: big.NewInt(params.LAT)},
		accounts[1].addr: {Balance: big.NewInt(params.LAT)},
		accounts[2].addr: {Balance: big.NewInt(params.LAT)},
		// An account with existing storage
		storageAccount: {
			Balance: new(big.Int),
			Storage: map[common.Hash]common.Hash{
				common.HexToHash("0x03"): common.HexToHash("0x33"),
				common.HexToHash("0x04"): common.HexToHash("0x44"),
			},
		},
	},
		BaseFee: big.NewInt(params.InitialBaseFee),
		Config:  params.TestChainConfig,
	}
	genBlocks := 10
	signer := types.HomesteadSigner{}
	backend := newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		tx, _ := types.SignTx(types.NewTransaction(uint64(i), accounts[1].addr, big.NewInt(1000), params.TxGas, b.BaseFee(), nil), signer, accounts[0].key)
		b.AddTx(tx)
	})
	defer backend.chain.Stop()
	api := NewAPI(backend)
	randomAccounts := newAccounts(3)
	type res struct {
		Gas         int
		Failed      bool
		ReturnValue string
	}
	var testSuite = []struct {
		blockNumber rpc.BlockNumber
		call        ethapi.TransactionArgs
		config      *TraceCallConfig
		expectErr   error
		want        string
	}{
		// Call which can only succeed if state is state overridden
		{
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From:  &randomAccounts[0].addr,
				To:    &randomAccounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config: &TraceCallConfig{
				StateOverrides: &ethapi.StateOverride{
					randomAccounts[0].addr: ethapi.OverrideAccount{Balance: newRPCBalance(new(big.Int).Mul(big.NewInt(1), big.NewInt(params.LAT)))},
				},
			},
			want: `{"gas":21000,"failed":false,"returnValue":""}`,
		},
		// Invalid call without state overriding
		{
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From:  &randomAccounts[0].addr,
				To:    &randomAccounts[1].addr,
				Value: (*hexutil.Big)(big.NewInt(1000)),
			},
			config:    &TraceCallConfig{},
			expectErr: core.ErrInsufficientFundsForTransfer,
		},
		// Successful simple contract call
		//
		// // SPDX-License-Identifier: GPL-3.0
		//
		//  pragma solidity >=0.7.0 <0.8.0;
		//
		//  /**
		//   * @title Storage
		//   * @dev Store & retrieve value in a variable
		//   */
		//  contract Storage {
		//      uint256 public number;
		//      constructor() {
		//          number = block.number;
		//      }
		//  }
		{
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From: &randomAccounts[0].addr,
				To:   &randomAccounts[2].addr,
				Data: newRPCBytes(common.Hex2Bytes("8381f58a")), // call number()
			},
			config: &TraceCallConfig{
				//Tracer: &tracer,
				StateOverrides: &ethapi.StateOverride{
					randomAccounts[2].addr: ethapi.OverrideAccount{
						Code:      newRPCBytes(common.Hex2Bytes("6080604052348015600f57600080fd5b506004361060285760003560e01c80638381f58a14602d575b600080fd5b60336049565b6040518082815260200191505060405180910390f35b6000548156fea2646970667358221220eab35ffa6ab2adfe380772a48b8ba78e82a1b820a18fcb6f59aa4efb20a5f60064736f6c63430007040033")),
						StateDiff: newStates([]common.Hash{{}}, []common.Hash{common.BigToHash(big.NewInt(123))}),
					},
				},
			},
			want: `{"gas":23347,"failed":false,"returnValue":"000000000000000000000000000000000000000000000000000000000000007b"}`,
		},
		{ // Override blocknumber
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From: &accounts[0].addr,
				// BLOCKNUMBER PUSH1 MSTORE
				Input: newRPCBytes(common.Hex2Bytes("4360005260206000f3")),
				//&hexutil.Bytes{0x43}, // blocknumber
			},
			config: &TraceCallConfig{
				BlockOverrides: &ethapi.BlockOverrides{Number: (*hexutil.Big)(big.NewInt(0x1337))},
			},
			want: `{"gas":59539,"failed":false,"returnValue":"0000000000000000000000000000000000000000000000000000000000001337","structLogs":[{"pc":0,"op":"NUMBER","gas":24946878,"gasCost":2,"depth":1,"stack":[]},{"pc":1,"op":"PUSH1","gas":24946876,"gasCost":3,"depth":1,"stack":["0x1337"]},{"pc":3,"op":"MSTORE","gas":24946873,"gasCost":6,"depth":1,"stack":["0x1337","0x0"]},{"pc":4,"op":"PUSH1","gas":24946867,"gasCost":3,"depth":1,"stack":[]},{"pc":6,"op":"PUSH1","gas":24946864,"gasCost":3,"depth":1,"stack":["0x20"]},{"pc":8,"op":"RETURN","gas":24946861,"gasCost":0,"depth":1,"stack":["0x20","0x0"]}]}`,
		},
		{ // Override blocknumber, and query a blockhash
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From: &accounts[0].addr,
				Input: &hexutil.Bytes{
					0x60, 0x00, 0x40, // BLOCKHASH(0)
					0x60, 0x00, 0x52, // STORE memory offset 0
					0x61, 0x13, 0x36, 0x40, // BLOCKHASH(0x1336)
					0x60, 0x20, 0x52, // STORE memory offset 32
					0x61, 0x13, 0x37, 0x40, // BLOCKHASH(0x1337)
					0x60, 0x40, 0x52, // STORE memory offset 64
					0x60, 0x60, 0x60, 0x00, 0xf3, // RETURN (0-96)

				}, // blocknumber
			},
			config: &TraceCallConfig{
				BlockOverrides: &ethapi.BlockOverrides{Number: (*hexutil.Big)(big.NewInt(0x1337))},
			},
			want: `{"gas":72668,"failed":false,"returnValue":"000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000","structLogs":[{"pc":0,"op":"PUSH1","gas":24946634,"gasCost":3,"depth":1,"stack":[]},{"pc":2,"op":"BLOCKHASH","gas":24946631,"gasCost":20,"depth":1,"stack":["0x0"]},{"pc":3,"op":"PUSH1","gas":24946611,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":5,"op":"MSTORE","gas":24946608,"gasCost":6,"depth":1,"stack":["0x0","0x0"]},{"pc":6,"op":"PUSH2","gas":24946602,"gasCost":3,"depth":1,"stack":[]},{"pc":9,"op":"BLOCKHASH","gas":24946599,"gasCost":20,"depth":1,"stack":["0x1336"]},{"pc":10,"op":"PUSH1","gas":24946579,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":12,"op":"MSTORE","gas":24946576,"gasCost":6,"depth":1,"stack":["0x0","0x20"]},{"pc":13,"op":"PUSH2","gas":24946570,"gasCost":3,"depth":1,"stack":[]},{"pc":16,"op":"BLOCKHASH","gas":24946567,"gasCost":20,"depth":1,"stack":["0x1337"]},{"pc":17,"op":"PUSH1","gas":24946547,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":19,"op":"MSTORE","gas":24946544,"gasCost":6,"depth":1,"stack":["0x0","0x40"]},{"pc":20,"op":"PUSH1","gas":24946538,"gasCost":3,"depth":1,"stack":[]},{"pc":22,"op":"PUSH1","gas":24946535,"gasCost":3,"depth":1,"stack":["0x60"]},{"pc":24,"op":"RETURN","gas":24946532,"gasCost":0,"depth":1,"stack":["0x60","0x0"]}]}`,
		},
		{ // First with only code override, not storage override
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From: &randomAccounts[0].addr,
				To:   &randomAccounts[2].addr,
				Data: newRPCBytes(common.Hex2Bytes("f8a8fd6d")), //
			},
			config: &TraceCallConfig{
				StateOverrides: &ethapi.StateOverride{
					randomAccounts[2].addr: ethapi.OverrideAccount{
						Code: newRPCBytes(common.Hex2Bytes("6080604052348015600f57600080fd5b506004361060325760003560e01c806366e41cb7146037578063f8a8fd6d14603f575b600080fd5b603d6057565b005b60456062565b60405190815260200160405180910390f35b610539600090815580fd5b60006001600081905550306001600160a01b03166366e41cb76040518163ffffffff1660e01b8152600401600060405180830381600087803b15801560a657600080fd5b505af192505050801560b6575060015b60e9573d80801560e1576040519150601f19603f3d011682016040523d82523d6000602084013e60e6565b606091505b50505b506000549056fea26469706673582212205ce45de745a5308f713cb2f448589177ba5a442d1a2eff945afaa8915961b4d064736f6c634300080c0033")),
					},
				},
			},
			want: `{"gas":27000,"failed":false,"returnValue":"0000000000000000000000000000000000000000000000000000000000000001","structLogs":[{"pc":0,"op":"PUSH1","gas":24978936,"gasCost":3,"depth":1,"stack":[]},{"pc":2,"op":"PUSH1","gas":24978933,"gasCost":3,"depth":1,"stack":["0x80"]},{"pc":4,"op":"MSTORE","gas":24978930,"gasCost":12,"depth":1,"stack":["0x80","0x40"]},{"pc":5,"op":"CALLVALUE","gas":24978918,"gasCost":2,"depth":1,"stack":[]},{"pc":6,"op":"DUP1","gas":24978916,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":7,"op":"ISZERO","gas":24978913,"gasCost":3,"depth":1,"stack":["0x0","0x0"]},{"pc":8,"op":"PUSH1","gas":24978910,"gasCost":3,"depth":1,"stack":["0x0","0x1"]},{"pc":10,"op":"JUMPI","gas":24978907,"gasCost":10,"depth":1,"stack":["0x0","0x1","0xf"]},{"pc":15,"op":"JUMPDEST","gas":24978897,"gasCost":1,"depth":1,"stack":["0x0"]},{"pc":16,"op":"POP","gas":24978896,"gasCost":2,"depth":1,"stack":["0x0"]},{"pc":17,"op":"PUSH1","gas":24978894,"gasCost":3,"depth":1,"stack":[]},{"pc":19,"op":"CALLDATASIZE","gas":24978891,"gasCost":2,"depth":1,"stack":["0x4"]},{"pc":20,"op":"LT","gas":24978889,"gasCost":3,"depth":1,"stack":["0x4","0x4"]},{"pc":21,"op":"PUSH1","gas":24978886,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":23,"op":"JUMPI","gas":24978883,"gasCost":10,"depth":1,"stack":["0x0","0x32"]},{"pc":24,"op":"PUSH1","gas":24978873,"gasCost":3,"depth":1,"stack":[]},{"pc":26,"op":"CALLDATALOAD","gas":24978870,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":27,"op":"PUSH1","gas":24978867,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d00000000000000000000000000000000000000000000000000000000"]},{"pc":29,"op":"SHR","gas":24978864,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d00000000000000000000000000000000000000000000000000000000","0xe0"]},{"pc":30,"op":"DUP1","gas":24978861,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d"]},{"pc":31,"op":"PUSH4","gas":24978858,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xf8a8fd6d"]},{"pc":36,"op":"EQ","gas":24978855,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xf8a8fd6d","0x66e41cb7"]},{"pc":37,"op":"PUSH1","gas":24978852,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x0"]},{"pc":39,"op":"JUMPI","gas":24978849,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x0","0x37"]},{"pc":40,"op":"DUP1","gas":24978839,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d"]},{"pc":41,"op":"PUSH4","gas":24978836,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xf8a8fd6d"]},{"pc":46,"op":"EQ","gas":24978833,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xf8a8fd6d","0xf8a8fd6d"]},{"pc":47,"op":"PUSH1","gas":24978830,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x1"]},{"pc":49,"op":"JUMPI","gas":24978827,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x1","0x3f"]},{"pc":63,"op":"JUMPDEST","gas":24978817,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d"]},{"pc":64,"op":"PUSH1","gas":24978816,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d"]},{"pc":66,"op":"PUSH1","gas":24978813,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45"]},{"pc":68,"op":"JUMP","gas":24978810,"gasCost":8,"depth":1,"stack":["0xf8a8fd6d","0x45","0x62"]},{"pc":98,"op":"JUMPDEST","gas":24978802,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45"]},{"pc":99,"op":"PUSH1","gas":24978801,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45"]},{"pc":101,"op":"PUSH1","gas":24978798,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":103,"op":"PUSH1","gas":24978795,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1"]},{"pc":105,"op":"DUP2","gas":24978792,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1","0x0"]},{"pc":106,"op":"SWAP1","gas":24978789,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1","0x0","0x1"]},{"pc":107,"op":"SSTORE","gas":24978786,"gasCost":5000,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1","0x1","0x0"],"storage":{"0000000000000000000000000000000000000000000000000000000000000000":"0000000000000000000000000000000000000000000000000000000000000001"}},{"pc":108,"op":"POP","gas":24973786,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1"]},{"pc":109,"op":"ADDRESS","gas":24973784,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":110,"op":"PUSH1","gas":24973782,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":112,"op":"PUSH1","gas":24973779,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x1"]},{"pc":114,"op":"PUSH1","gas":24973776,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x1","0x1"]},{"pc":116,"op":"SHL","gas":24973773,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x1","0x1","0xa0"]},{"pc":117,"op":"SUB","gas":24973770,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x1","0x10000000000000000000000000000000000000000"]},{"pc":118,"op":"AND","gas":24973767,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0xffffffffffffffffffffffffffffffffffffffff"]},{"pc":119,"op":"PUSH4","gas":24973764,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":124,"op":"PUSH1","gas":24973761,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7"]},{"pc":126,"op":"MLOAD","gas":24973758,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x40"]},{"pc":127,"op":"DUP2","gas":24973755,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80"]},{"pc":128,"op":"PUSH4","gas":24973752,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb7"]},{"pc":133,"op":"AND","gas":24973749,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb7","0xffffffff"]},{"pc":134,"op":"PUSH1","gas":24973746,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb7"]},{"pc":136,"op":"SHL","gas":24973743,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb7","0xe0"]},{"pc":137,"op":"DUP2","gas":24973740,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb700000000000000000000000000000000000000000000000000000000"]},{"pc":138,"op":"MSTORE","gas":24973737,"gasCost":9,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb700000000000000000000000000000000000000000000000000000000","0x80"]},{"pc":139,"op":"PUSH1","gas":24973728,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80"]},{"pc":141,"op":"ADD","gas":24973725,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x4"]},{"pc":142,"op":"PUSH1","gas":24973722,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84"]},{"pc":144,"op":"PUSH1","gas":24973719,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0"]},{"pc":146,"op":"MLOAD","gas":24973716,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x40"]},{"pc":147,"op":"DUP1","gas":24973713,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80"]},{"pc":148,"op":"DUP4","gas":24973710,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x80"]},{"pc":149,"op":"SUB","gas":24973707,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x80","0x84"]},{"pc":150,"op":"DUP2","gas":24973704,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4"]},{"pc":151,"op":"PUSH1","gas":24973701,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80"]},{"pc":153,"op":"DUP8","gas":24973698,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0"]},{"pc":154,"op":"DUP1","gas":24973695,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":155,"op":"EXTCODESIZE","gas":24973692,"gasCost":100,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":156,"op":"ISZERO","gas":24973592,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x126"]},{"pc":157,"op":"DUP1","gas":24973589,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0"]},{"pc":158,"op":"ISZERO","gas":24973586,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0","0x0"]},{"pc":159,"op":"PUSH1","gas":24973583,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0","0x1"]},{"pc":161,"op":"JUMPI","gas":24973580,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0","0x1","0xa6"]},{"pc":166,"op":"JUMPDEST","gas":24973570,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0"]},{"pc":167,"op":"POP","gas":24973569,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0"]},{"pc":168,"op":"GAS","gas":24973567,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":169,"op":"CALL","gas":24973565,"gasCost":24583355,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x17d10fd"]},{"pc":0,"op":"PUSH1","gas":24583255,"gasCost":3,"depth":2,"stack":[]},{"pc":2,"op":"PUSH1","gas":24583252,"gasCost":3,"depth":2,"stack":["0x80"]},{"pc":4,"op":"MSTORE","gas":24583249,"gasCost":12,"depth":2,"stack":["0x80","0x40"]},{"pc":5,"op":"CALLVALUE","gas":24583237,"gasCost":2,"depth":2,"stack":[]},{"pc":6,"op":"DUP1","gas":24583235,"gasCost":3,"depth":2,"stack":["0x0"]},{"pc":7,"op":"ISZERO","gas":24583232,"gasCost":3,"depth":2,"stack":["0x0","0x0"]},{"pc":8,"op":"PUSH1","gas":24583229,"gasCost":3,"depth":2,"stack":["0x0","0x1"]},{"pc":10,"op":"JUMPI","gas":24583226,"gasCost":10,"depth":2,"stack":["0x0","0x1","0xf"]},{"pc":15,"op":"JUMPDEST","gas":24583216,"gasCost":1,"depth":2,"stack":["0x0"]},{"pc":16,"op":"POP","gas":24583215,"gasCost":2,"depth":2,"stack":["0x0"]},{"pc":17,"op":"PUSH1","gas":24583213,"gasCost":3,"depth":2,"stack":[]},{"pc":19,"op":"CALLDATASIZE","gas":24583210,"gasCost":2,"depth":2,"stack":["0x4"]},{"pc":20,"op":"LT","gas":24583208,"gasCost":3,"depth":2,"stack":["0x4","0x4"]},{"pc":21,"op":"PUSH1","gas":24583205,"gasCost":3,"depth":2,"stack":["0x0"]},{"pc":23,"op":"JUMPI","gas":24583202,"gasCost":10,"depth":2,"stack":["0x0","0x32"]},{"pc":24,"op":"PUSH1","gas":24583192,"gasCost":3,"depth":2,"stack":[]},{"pc":26,"op":"CALLDATALOAD","gas":24583189,"gasCost":3,"depth":2,"stack":["0x0"]},{"pc":27,"op":"PUSH1","gas":24583186,"gasCost":3,"depth":2,"stack":["0x66e41cb700000000000000000000000000000000000000000000000000000000"]},{"pc":29,"op":"SHR","gas":24583183,"gasCost":3,"depth":2,"stack":["0x66e41cb700000000000000000000000000000000000000000000000000000000","0xe0"]},{"pc":30,"op":"DUP1","gas":24583180,"gasCost":3,"depth":2,"stack":["0x66e41cb7"]},{"pc":31,"op":"PUSH4","gas":24583177,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x66e41cb7"]},{"pc":36,"op":"EQ","gas":24583174,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x66e41cb7","0x66e41cb7"]},{"pc":37,"op":"PUSH1","gas":24583171,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x1"]},{"pc":39,"op":"JUMPI","gas":24583168,"gasCost":10,"depth":2,"stack":["0x66e41cb7","0x1","0x37"]},{"pc":55,"op":"JUMPDEST","gas":24583158,"gasCost":1,"depth":2,"stack":["0x66e41cb7"]},{"pc":56,"op":"PUSH1","gas":24583157,"gasCost":3,"depth":2,"stack":["0x66e41cb7"]},{"pc":58,"op":"PUSH1","gas":24583154,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d"]},{"pc":60,"op":"JUMP","gas":24583151,"gasCost":8,"depth":2,"stack":["0x66e41cb7","0x3d","0x57"]},{"pc":87,"op":"JUMPDEST","gas":24583143,"gasCost":1,"depth":2,"stack":["0x66e41cb7","0x3d"]},{"pc":88,"op":"PUSH2","gas":24583142,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d"]},{"pc":91,"op":"PUSH1","gas":24583139,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d","0x539"]},{"pc":93,"op":"SWAP1","gas":24583136,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d","0x539","0x0"]},{"pc":94,"op":"DUP2","gas":24583133,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d","0x0","0x539"]},{"pc":95,"op":"SSTORE","gas":24583130,"gasCost":100,"depth":2,"stack":["0x66e41cb7","0x3d","0x0","0x539","0x0"],"storage":{"0000000000000000000000000000000000000000000000000000000000000000":"0000000000000000000000000000000000000000000000000000000000000539"}},{"pc":96,"op":"DUP1","gas":24583030,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d","0x0"]},{"pc":97,"op":"REVERT","gas":24583027,"gasCost":0,"depth":2,"stack":["0x66e41cb7","0x3d","0x0","0x0"]},{"pc":170,"op":"SWAP3","gas":24973237,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0"]},{"pc":171,"op":"POP","gas":24973234,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x66e41cb7","0x84","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":172,"op":"POP","gas":24973232,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x66e41cb7","0x84"]},{"pc":173,"op":"POP","gas":24973230,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x66e41cb7"]},{"pc":174,"op":"DUP1","gas":24973228,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0"]},{"pc":175,"op":"ISZERO","gas":24973225,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0"]},{"pc":176,"op":"PUSH1","gas":24973222,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x1"]},{"pc":178,"op":"JUMPI","gas":24973219,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x1","0xb6"]},{"pc":182,"op":"JUMPDEST","gas":24973209,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0"]},{"pc":183,"op":"PUSH1","gas":24973208,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0"]},{"pc":185,"op":"JUMPI","gas":24973205,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0xe9"]},{"pc":186,"op":"RETURNDATASIZE","gas":24973195,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":187,"op":"DUP1","gas":24973193,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0"]},{"pc":188,"op":"DUP1","gas":24973190,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0"]},{"pc":189,"op":"ISZERO","gas":24973187,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0","0x0"]},{"pc":190,"op":"PUSH1","gas":24973184,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0","0x1"]},{"pc":192,"op":"JUMPI","gas":24973181,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0","0x1","0xe1"]},{"pc":225,"op":"JUMPDEST","gas":24973171,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0"]},{"pc":226,"op":"PUSH1","gas":24973170,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0"]},{"pc":228,"op":"SWAP2","gas":24973167,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0","0x60"]},{"pc":229,"op":"POP","gas":24973164,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x60","0x0","0x0"]},{"pc":230,"op":"JUMPDEST","gas":24973162,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x60","0x0"]},{"pc":231,"op":"POP","gas":24973161,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x60","0x0"]},{"pc":232,"op":"POP","gas":24973159,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x60"]},{"pc":233,"op":"JUMPDEST","gas":24973157,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":234,"op":"POP","gas":24973156,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":235,"op":"PUSH1","gas":24973154,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45"]},{"pc":237,"op":"SLOAD","gas":24973151,"gasCost":100,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"],"storage":{"0000000000000000000000000000000000000000000000000000000000000000":"0000000000000000000000000000000000000000000000000000000000000001"}},{"pc":238,"op":"SWAP1","gas":24973051,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x1"]},{"pc":239,"op":"JUMP","gas":24973048,"gasCost":8,"depth":1,"stack":["0xf8a8fd6d","0x1","0x45"]},{"pc":69,"op":"JUMPDEST","gas":24973040,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x1"]},{"pc":70,"op":"PUSH1","gas":24973039,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x1"]},{"pc":72,"op":"MLOAD","gas":24973036,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x1","0x40"]},{"pc":73,"op":"SWAP1","gas":24973033,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x1","0x80"]},{"pc":74,"op":"DUP2","gas":24973030,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x1"]},{"pc":75,"op":"MSTORE","gas":24973027,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x1","0x80"]},{"pc":76,"op":"PUSH1","gas":24973024,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80"]},{"pc":78,"op":"ADD","gas":24973021,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x20"]},{"pc":79,"op":"PUSH1","gas":24973018,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xa0"]},{"pc":81,"op":"MLOAD","gas":24973015,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xa0","0x40"]},{"pc":82,"op":"DUP1","gas":24973012,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xa0","0x80"]},{"pc":83,"op":"SWAP2","gas":24973009,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xa0","0x80","0x80"]},{"pc":84,"op":"SUB","gas":24973006,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x80","0xa0"]},{"pc":85,"op":"SWAP1","gas":24973003,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x20"]},{"pc":86,"op":"RETURN","gas":24973000,"gasCost":0,"depth":1,"stack":["0xf8a8fd6d","0x20","0x80"]}]}`,
		},
		{ // Same again, this time with storage override
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From: &randomAccounts[0].addr,
				To:   &randomAccounts[2].addr,
				Data: newRPCBytes(common.Hex2Bytes("f8a8fd6d")), //
			},
			config: &TraceCallConfig{
				StateOverrides: &ethapi.StateOverride{
					randomAccounts[2].addr: ethapi.OverrideAccount{
						Code:  newRPCBytes(common.Hex2Bytes("6080604052348015600f57600080fd5b506004361060325760003560e01c806366e41cb7146037578063f8a8fd6d14603f575b600080fd5b603d6057565b005b60456062565b60405190815260200160405180910390f35b610539600090815580fd5b60006001600081905550306001600160a01b03166366e41cb76040518163ffffffff1660e01b8152600401600060405180830381600087803b15801560a657600080fd5b505af192505050801560b6575060015b60e9573d80801560e1576040519150601f19603f3d011682016040523d82523d6000602084013e60e6565b606091505b50505b506000549056fea26469706673582212205ce45de745a5308f713cb2f448589177ba5a442d1a2eff945afaa8915961b4d064736f6c634300080c0033")),
						State: newStates([]common.Hash{{}}, []common.Hash{{}}),
					},
				},
			},
			//want: `{"gas":46900,"failed":false,"returnValue":"0000000000000000000000000000000000000000000000000000000000000539"}`,
			want: `{"gas":27000,"failed":false,"returnValue":"0000000000000000000000000000000000000000000000000000000000000001","structLogs":[{"pc":0,"op":"PUSH1","gas":24978936,"gasCost":3,"depth":1,"stack":[]},{"pc":2,"op":"PUSH1","gas":24978933,"gasCost":3,"depth":1,"stack":["0x80"]},{"pc":4,"op":"MSTORE","gas":24978930,"gasCost":12,"depth":1,"stack":["0x80","0x40"]},{"pc":5,"op":"CALLVALUE","gas":24978918,"gasCost":2,"depth":1,"stack":[]},{"pc":6,"op":"DUP1","gas":24978916,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":7,"op":"ISZERO","gas":24978913,"gasCost":3,"depth":1,"stack":["0x0","0x0"]},{"pc":8,"op":"PUSH1","gas":24978910,"gasCost":3,"depth":1,"stack":["0x0","0x1"]},{"pc":10,"op":"JUMPI","gas":24978907,"gasCost":10,"depth":1,"stack":["0x0","0x1","0xf"]},{"pc":15,"op":"JUMPDEST","gas":24978897,"gasCost":1,"depth":1,"stack":["0x0"]},{"pc":16,"op":"POP","gas":24978896,"gasCost":2,"depth":1,"stack":["0x0"]},{"pc":17,"op":"PUSH1","gas":24978894,"gasCost":3,"depth":1,"stack":[]},{"pc":19,"op":"CALLDATASIZE","gas":24978891,"gasCost":2,"depth":1,"stack":["0x4"]},{"pc":20,"op":"LT","gas":24978889,"gasCost":3,"depth":1,"stack":["0x4","0x4"]},{"pc":21,"op":"PUSH1","gas":24978886,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":23,"op":"JUMPI","gas":24978883,"gasCost":10,"depth":1,"stack":["0x0","0x32"]},{"pc":24,"op":"PUSH1","gas":24978873,"gasCost":3,"depth":1,"stack":[]},{"pc":26,"op":"CALLDATALOAD","gas":24978870,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":27,"op":"PUSH1","gas":24978867,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d00000000000000000000000000000000000000000000000000000000"]},{"pc":29,"op":"SHR","gas":24978864,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d00000000000000000000000000000000000000000000000000000000","0xe0"]},{"pc":30,"op":"DUP1","gas":24978861,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d"]},{"pc":31,"op":"PUSH4","gas":24978858,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xf8a8fd6d"]},{"pc":36,"op":"EQ","gas":24978855,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xf8a8fd6d","0x66e41cb7"]},{"pc":37,"op":"PUSH1","gas":24978852,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x0"]},{"pc":39,"op":"JUMPI","gas":24978849,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x0","0x37"]},{"pc":40,"op":"DUP1","gas":24978839,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d"]},{"pc":41,"op":"PUSH4","gas":24978836,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xf8a8fd6d"]},{"pc":46,"op":"EQ","gas":24978833,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xf8a8fd6d","0xf8a8fd6d"]},{"pc":47,"op":"PUSH1","gas":24978830,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x1"]},{"pc":49,"op":"JUMPI","gas":24978827,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x1","0x3f"]},{"pc":63,"op":"JUMPDEST","gas":24978817,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d"]},{"pc":64,"op":"PUSH1","gas":24978816,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d"]},{"pc":66,"op":"PUSH1","gas":24978813,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45"]},{"pc":68,"op":"JUMP","gas":24978810,"gasCost":8,"depth":1,"stack":["0xf8a8fd6d","0x45","0x62"]},{"pc":98,"op":"JUMPDEST","gas":24978802,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45"]},{"pc":99,"op":"PUSH1","gas":24978801,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45"]},{"pc":101,"op":"PUSH1","gas":24978798,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":103,"op":"PUSH1","gas":24978795,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1"]},{"pc":105,"op":"DUP2","gas":24978792,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1","0x0"]},{"pc":106,"op":"SWAP1","gas":24978789,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1","0x0","0x1"]},{"pc":107,"op":"SSTORE","gas":24978786,"gasCost":5000,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1","0x1","0x0"],"storage":{"0000000000000000000000000000000000000000000000000000000000000000":"0000000000000000000000000000000000000000000000000000000000000001"}},{"pc":108,"op":"POP","gas":24973786,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x1"]},{"pc":109,"op":"ADDRESS","gas":24973784,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":110,"op":"PUSH1","gas":24973782,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":112,"op":"PUSH1","gas":24973779,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x1"]},{"pc":114,"op":"PUSH1","gas":24973776,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x1","0x1"]},{"pc":116,"op":"SHL","gas":24973773,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x1","0x1","0xa0"]},{"pc":117,"op":"SUB","gas":24973770,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x1","0x10000000000000000000000000000000000000000"]},{"pc":118,"op":"AND","gas":24973767,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0xffffffffffffffffffffffffffffffffffffffff"]},{"pc":119,"op":"PUSH4","gas":24973764,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":124,"op":"PUSH1","gas":24973761,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7"]},{"pc":126,"op":"MLOAD","gas":24973758,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x40"]},{"pc":127,"op":"DUP2","gas":24973755,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80"]},{"pc":128,"op":"PUSH4","gas":24973752,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb7"]},{"pc":133,"op":"AND","gas":24973749,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb7","0xffffffff"]},{"pc":134,"op":"PUSH1","gas":24973746,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb7"]},{"pc":136,"op":"SHL","gas":24973743,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb7","0xe0"]},{"pc":137,"op":"DUP2","gas":24973740,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb700000000000000000000000000000000000000000000000000000000"]},{"pc":138,"op":"MSTORE","gas":24973737,"gasCost":9,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x66e41cb700000000000000000000000000000000000000000000000000000000","0x80"]},{"pc":139,"op":"PUSH1","gas":24973728,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80"]},{"pc":141,"op":"ADD","gas":24973725,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x80","0x4"]},{"pc":142,"op":"PUSH1","gas":24973722,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84"]},{"pc":144,"op":"PUSH1","gas":24973719,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0"]},{"pc":146,"op":"MLOAD","gas":24973716,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x40"]},{"pc":147,"op":"DUP1","gas":24973713,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80"]},{"pc":148,"op":"DUP4","gas":24973710,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x80"]},{"pc":149,"op":"SUB","gas":24973707,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x80","0x84"]},{"pc":150,"op":"DUP2","gas":24973704,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4"]},{"pc":151,"op":"PUSH1","gas":24973701,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80"]},{"pc":153,"op":"DUP8","gas":24973698,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0"]},{"pc":154,"op":"DUP1","gas":24973695,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":155,"op":"EXTCODESIZE","gas":24973692,"gasCost":100,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":156,"op":"ISZERO","gas":24973592,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x126"]},{"pc":157,"op":"DUP1","gas":24973589,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0"]},{"pc":158,"op":"ISZERO","gas":24973586,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0","0x0"]},{"pc":159,"op":"PUSH1","gas":24973583,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0","0x1"]},{"pc":161,"op":"JUMPI","gas":24973580,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0","0x1","0xa6"]},{"pc":166,"op":"JUMPDEST","gas":24973570,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0"]},{"pc":167,"op":"POP","gas":24973569,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x0"]},{"pc":168,"op":"GAS","gas":24973567,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":169,"op":"CALL","gas":24973565,"gasCost":24583355,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0","0x80","0x4","0x80","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x17d10fd"]},{"pc":0,"op":"PUSH1","gas":24583255,"gasCost":3,"depth":2,"stack":[]},{"pc":2,"op":"PUSH1","gas":24583252,"gasCost":3,"depth":2,"stack":["0x80"]},{"pc":4,"op":"MSTORE","gas":24583249,"gasCost":12,"depth":2,"stack":["0x80","0x40"]},{"pc":5,"op":"CALLVALUE","gas":24583237,"gasCost":2,"depth":2,"stack":[]},{"pc":6,"op":"DUP1","gas":24583235,"gasCost":3,"depth":2,"stack":["0x0"]},{"pc":7,"op":"ISZERO","gas":24583232,"gasCost":3,"depth":2,"stack":["0x0","0x0"]},{"pc":8,"op":"PUSH1","gas":24583229,"gasCost":3,"depth":2,"stack":["0x0","0x1"]},{"pc":10,"op":"JUMPI","gas":24583226,"gasCost":10,"depth":2,"stack":["0x0","0x1","0xf"]},{"pc":15,"op":"JUMPDEST","gas":24583216,"gasCost":1,"depth":2,"stack":["0x0"]},{"pc":16,"op":"POP","gas":24583215,"gasCost":2,"depth":2,"stack":["0x0"]},{"pc":17,"op":"PUSH1","gas":24583213,"gasCost":3,"depth":2,"stack":[]},{"pc":19,"op":"CALLDATASIZE","gas":24583210,"gasCost":2,"depth":2,"stack":["0x4"]},{"pc":20,"op":"LT","gas":24583208,"gasCost":3,"depth":2,"stack":["0x4","0x4"]},{"pc":21,"op":"PUSH1","gas":24583205,"gasCost":3,"depth":2,"stack":["0x0"]},{"pc":23,"op":"JUMPI","gas":24583202,"gasCost":10,"depth":2,"stack":["0x0","0x32"]},{"pc":24,"op":"PUSH1","gas":24583192,"gasCost":3,"depth":2,"stack":[]},{"pc":26,"op":"CALLDATALOAD","gas":24583189,"gasCost":3,"depth":2,"stack":["0x0"]},{"pc":27,"op":"PUSH1","gas":24583186,"gasCost":3,"depth":2,"stack":["0x66e41cb700000000000000000000000000000000000000000000000000000000"]},{"pc":29,"op":"SHR","gas":24583183,"gasCost":3,"depth":2,"stack":["0x66e41cb700000000000000000000000000000000000000000000000000000000","0xe0"]},{"pc":30,"op":"DUP1","gas":24583180,"gasCost":3,"depth":2,"stack":["0x66e41cb7"]},{"pc":31,"op":"PUSH4","gas":24583177,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x66e41cb7"]},{"pc":36,"op":"EQ","gas":24583174,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x66e41cb7","0x66e41cb7"]},{"pc":37,"op":"PUSH1","gas":24583171,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x1"]},{"pc":39,"op":"JUMPI","gas":24583168,"gasCost":10,"depth":2,"stack":["0x66e41cb7","0x1","0x37"]},{"pc":55,"op":"JUMPDEST","gas":24583158,"gasCost":1,"depth":2,"stack":["0x66e41cb7"]},{"pc":56,"op":"PUSH1","gas":24583157,"gasCost":3,"depth":2,"stack":["0x66e41cb7"]},{"pc":58,"op":"PUSH1","gas":24583154,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d"]},{"pc":60,"op":"JUMP","gas":24583151,"gasCost":8,"depth":2,"stack":["0x66e41cb7","0x3d","0x57"]},{"pc":87,"op":"JUMPDEST","gas":24583143,"gasCost":1,"depth":2,"stack":["0x66e41cb7","0x3d"]},{"pc":88,"op":"PUSH2","gas":24583142,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d"]},{"pc":91,"op":"PUSH1","gas":24583139,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d","0x539"]},{"pc":93,"op":"SWAP1","gas":24583136,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d","0x539","0x0"]},{"pc":94,"op":"DUP2","gas":24583133,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d","0x0","0x539"]},{"pc":95,"op":"SSTORE","gas":24583130,"gasCost":100,"depth":2,"stack":["0x66e41cb7","0x3d","0x0","0x539","0x0"],"storage":{"0000000000000000000000000000000000000000000000000000000000000000":"0000000000000000000000000000000000000000000000000000000000000539"}},{"pc":96,"op":"DUP1","gas":24583030,"gasCost":3,"depth":2,"stack":["0x66e41cb7","0x3d","0x0"]},{"pc":97,"op":"REVERT","gas":24583027,"gasCost":0,"depth":2,"stack":["0x66e41cb7","0x3d","0x0","0x0"]},{"pc":170,"op":"SWAP3","gas":24973237,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba","0x66e41cb7","0x84","0x0"]},{"pc":171,"op":"POP","gas":24973234,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x66e41cb7","0x84","0xec154c779028f8ef5c1fa9b03e4cb82671a2caba"]},{"pc":172,"op":"POP","gas":24973232,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x66e41cb7","0x84"]},{"pc":173,"op":"POP","gas":24973230,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x66e41cb7"]},{"pc":174,"op":"DUP1","gas":24973228,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0"]},{"pc":175,"op":"ISZERO","gas":24973225,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0"]},{"pc":176,"op":"PUSH1","gas":24973222,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x1"]},{"pc":178,"op":"JUMPI","gas":24973219,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x1","0xb6"]},{"pc":182,"op":"JUMPDEST","gas":24973209,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0"]},{"pc":183,"op":"PUSH1","gas":24973208,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0"]},{"pc":185,"op":"JUMPI","gas":24973205,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0xe9"]},{"pc":186,"op":"RETURNDATASIZE","gas":24973195,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":187,"op":"DUP1","gas":24973193,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0"]},{"pc":188,"op":"DUP1","gas":24973190,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0"]},{"pc":189,"op":"ISZERO","gas":24973187,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0","0x0"]},{"pc":190,"op":"PUSH1","gas":24973184,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0","0x1"]},{"pc":192,"op":"JUMPI","gas":24973181,"gasCost":10,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0","0x1","0xe1"]},{"pc":225,"op":"JUMPDEST","gas":24973171,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0"]},{"pc":226,"op":"PUSH1","gas":24973170,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0"]},{"pc":228,"op":"SWAP2","gas":24973167,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x0","0x0","0x60"]},{"pc":229,"op":"POP","gas":24973164,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x60","0x0","0x0"]},{"pc":230,"op":"JUMPDEST","gas":24973162,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x60","0x0"]},{"pc":231,"op":"POP","gas":24973161,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x60","0x0"]},{"pc":232,"op":"POP","gas":24973159,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0","0x60"]},{"pc":233,"op":"JUMPDEST","gas":24973157,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":234,"op":"POP","gas":24973156,"gasCost":2,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"]},{"pc":235,"op":"PUSH1","gas":24973154,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45"]},{"pc":237,"op":"SLOAD","gas":24973151,"gasCost":100,"depth":1,"stack":["0xf8a8fd6d","0x45","0x0"],"storage":{"0000000000000000000000000000000000000000000000000000000000000000":"0000000000000000000000000000000000000000000000000000000000000001"}},{"pc":238,"op":"SWAP1","gas":24973051,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x45","0x1"]},{"pc":239,"op":"JUMP","gas":24973048,"gasCost":8,"depth":1,"stack":["0xf8a8fd6d","0x1","0x45"]},{"pc":69,"op":"JUMPDEST","gas":24973040,"gasCost":1,"depth":1,"stack":["0xf8a8fd6d","0x1"]},{"pc":70,"op":"PUSH1","gas":24973039,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x1"]},{"pc":72,"op":"MLOAD","gas":24973036,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x1","0x40"]},{"pc":73,"op":"SWAP1","gas":24973033,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x1","0x80"]},{"pc":74,"op":"DUP2","gas":24973030,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x1"]},{"pc":75,"op":"MSTORE","gas":24973027,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x1","0x80"]},{"pc":76,"op":"PUSH1","gas":24973024,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80"]},{"pc":78,"op":"ADD","gas":24973021,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x20"]},{"pc":79,"op":"PUSH1","gas":24973018,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xa0"]},{"pc":81,"op":"MLOAD","gas":24973015,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xa0","0x40"]},{"pc":82,"op":"DUP1","gas":24973012,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xa0","0x80"]},{"pc":83,"op":"SWAP2","gas":24973009,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0xa0","0x80","0x80"]},{"pc":84,"op":"SUB","gas":24973006,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x80","0xa0"]},{"pc":85,"op":"SWAP1","gas":24973003,"gasCost":3,"depth":1,"stack":["0xf8a8fd6d","0x80","0x20"]},{"pc":86,"op":"RETURN","gas":24973000,"gasCost":0,"depth":1,"stack":["0xf8a8fd6d","0x20","0x80"]}]}`,
		},
		{ // No state override
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From: &randomAccounts[0].addr,
				To:   &storageAccount,
				Data: newRPCBytes(common.Hex2Bytes("f8a8fd6d")), //
			},
			config: &TraceCallConfig{
				StateOverrides: &ethapi.StateOverride{
					storageAccount: ethapi.OverrideAccount{
						Code: newRPCBytes([]byte{
							// SLOAD(3) + SLOAD(4) (which is 0x77)
							byte(vm.PUSH1), 0x04,
							byte(vm.SLOAD),
							byte(vm.PUSH1), 0x03,
							byte(vm.SLOAD),
							byte(vm.ADD),
							// 0x77 -> MSTORE(0)
							byte(vm.PUSH1), 0x00,
							byte(vm.MSTORE),
							// RETURN (0, 32)
							byte(vm.PUSH1), 32,
							byte(vm.PUSH1), 00,
							byte(vm.RETURN),
						}),
					},
				},
			},
			want: `{"gas":25288,"failed":false,"returnValue":"0000000000000000000000000000000000000000000000000000000000000077"}`,
		},
		{ // Full state override
			// The original storage is
			// 3: 0x33
			// 4: 0x44
			// With a full override, where we set 3:0x11, the slot 4 should be
			// removed. So SLOT(3)+SLOT(4) should be 0x11.
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From: &randomAccounts[0].addr,
				To:   &storageAccount,
				Data: newRPCBytes(common.Hex2Bytes("f8a8fd6d")), //
			},
			config: &TraceCallConfig{
				StateOverrides: &ethapi.StateOverride{
					storageAccount: ethapi.OverrideAccount{
						Code: newRPCBytes([]byte{
							// SLOAD(3) + SLOAD(4) (which is now 0x11 + 0x00)
							byte(vm.PUSH1), 0x04,
							byte(vm.SLOAD),
							byte(vm.PUSH1), 0x03,
							byte(vm.SLOAD),
							byte(vm.ADD),
							// 0x11 -> MSTORE(0)
							byte(vm.PUSH1), 0x00,
							byte(vm.MSTORE),
							// RETURN (0, 32)
							byte(vm.PUSH1), 32,
							byte(vm.PUSH1), 00,
							byte(vm.RETURN),
						}),
						State: newStates(
							[]common.Hash{common.HexToHash("0x03")},
							[]common.Hash{common.HexToHash("0x11")}),
					},
				},
			},
			want: `{"gas":25288,"failed":false,"returnValue":"0000000000000000000000000000000000000000000000000000000000000000","structLogs":[{"pc":0,"op":"PUSH1","gas":24978936,"gasCost":3,"depth":1,"stack":[]},{"pc":2,"op":"SLOAD","gas":24978933,"gasCost":2100,"depth":1,"stack":["0x4"],"storage":{"0000000000000000000000000000000000000000000000000000000000000004":"0000000000000000000000000000000000000000000000000000000000000000"}},{"pc":3,"op":"PUSH1","gas":24976833,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":5,"op":"SLOAD","gas":24976830,"gasCost":2100,"depth":1,"stack":["0x0","0x3"],"storage":{"0000000000000000000000000000000000000000000000000000000000000003":"0000000000000000000000000000000000000000000000000000000000000000","0000000000000000000000000000000000000000000000000000000000000004":"0000000000000000000000000000000000000000000000000000000000000000"}},{"pc":6,"op":"ADD","gas":24974730,"gasCost":3,"depth":1,"stack":["0x0","0x0"]},{"pc":7,"op":"PUSH1","gas":24974727,"gasCost":3,"depth":1,"stack":["0x0"]},{"pc":9,"op":"MSTORE","gas":24974724,"gasCost":6,"depth":1,"stack":["0x0","0x0"]},{"pc":10,"op":"PUSH1","gas":24974718,"gasCost":3,"depth":1,"stack":[]},{"pc":12,"op":"PUSH1","gas":24974715,"gasCost":3,"depth":1,"stack":["0x20"]},{"pc":14,"op":"RETURN","gas":24974712,"gasCost":0,"depth":1,"stack":["0x20","0x0"]}]}`,
		},
		{ // Partial state override
			// The original storage is
			// 3: 0x33
			// 4: 0x44
			// With a partial override, where we set 3:0x11, the slot 4 as before.
			// So SLOT(3)+SLOT(4) should be 0x55.
			blockNumber: rpc.LatestBlockNumber,
			call: ethapi.TransactionArgs{
				From: &randomAccounts[0].addr,
				To:   &storageAccount,
				Data: newRPCBytes(common.Hex2Bytes("f8a8fd6d")), //
			},
			config: &TraceCallConfig{
				StateOverrides: &ethapi.StateOverride{
					storageAccount: ethapi.OverrideAccount{
						Code: newRPCBytes([]byte{
							// SLOAD(3) + SLOAD(4) (which is now 0x11 + 0x44)
							byte(vm.PUSH1), 0x04,
							byte(vm.SLOAD),
							byte(vm.PUSH1), 0x03,
							byte(vm.SLOAD),
							byte(vm.ADD),
							// 0x55 -> MSTORE(0)
							byte(vm.PUSH1), 0x00,
							byte(vm.MSTORE),
							// RETURN (0, 32)
							byte(vm.PUSH1), 32,
							byte(vm.PUSH1), 00,
							byte(vm.RETURN),
						}),
						StateDiff: &map[common.Hash]common.Hash{
							common.HexToHash("0x03"): common.HexToHash("0x11"),
						},
					},
				},
			},
			want: `{"gas":25288,"failed":false,"returnValue":"0000000000000000000000000000000000000000000000000000000000000055"}`,
		},
	}
	for i, tc := range testSuite {
		t.Log("index", i, "tc", tc.want)
		result, err := api.TraceCall(context.Background(), tc.call, rpc.BlockNumberOrHash{BlockNumber: &tc.blockNumber}, tc.config)
		if tc.expectErr != nil {
			if err == nil {
				t.Errorf("test %d: want error %v, have nothing", i, tc.expectErr)
				continue
			}
			if !errors.Is(err, tc.expectErr) {
				t.Errorf("test %d: error mismatch, want %v, have %v", i, tc.expectErr, err)
			}
			continue
		}
		if err != nil {
			t.Errorf("test %d: want no error, have %v", i, err)
			continue
		}
		// Turn result into res-struct
		var (
			have res
			want res
		)
		resBytes, _ := json.Marshal(result)
		json.Unmarshal(resBytes, &have)
		json.Unmarshal([]byte(tc.want), &want)
		if !reflect.DeepEqual(have, want) {
			t.Errorf("test %d, result mismatch, have\n%v\n, want\n%v\n", i, string(resBytes), want)
		}
	}
}

type Account struct {
	key  *ecdsa.PrivateKey
	addr common.Address
}

type Accounts []Account

func (a Accounts) Len() int           { return len(a) }
func (a Accounts) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Accounts) Less(i, j int) bool { return bytes.Compare(a[i].addr.Bytes(), a[j].addr.Bytes()) < 0 }

func newAccounts(n int) (accounts Accounts) {
	for i := 0; i < n; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		accounts = append(accounts, Account{key: key, addr: addr})
	}
	sort.Sort(accounts)
	return accounts
}

func newRPCBalance(balance *big.Int) **hexutil.Big {
	rpcBalance := (*hexutil.Big)(balance)
	return &rpcBalance
}

func newRPCBytes(bytes []byte) *hexutil.Bytes {
	rpcBytes := hexutil.Bytes(bytes)
	return &rpcBytes
}

func newStates(keys []common.Hash, vals []common.Hash) *map[common.Hash]common.Hash {
	if len(keys) != len(vals) {
		panic("invalid input")
	}
	m := make(map[common.Hash]common.Hash)
	for i := 0; i < len(keys); i++ {
		m[keys[i]] = vals[i]
	}
	return &m
}

func TestTraceChain(t *testing.T) {
	// Initialize test accounts
	accounts := newAccounts(3)
	genesis := &core.Genesis{
		Config: params.TestChainConfig,
		Alloc: core.GenesisAlloc{
			accounts[0].addr: {Balance: big.NewInt(params.LAT)},
			accounts[1].addr: {Balance: big.NewInt(params.LAT)},
			accounts[2].addr: {Balance: big.NewInt(params.LAT)},
		},
		EconomicModel: params.GetEc(params.DefaultUnitTestNet),
	}
	genBlocks := 50
	signer := types.HomesteadSigner{}

	var (
		ref   uint32 // total refs has made
		rel   uint32 // total rels has made
		nonce uint64
	)
	backend := newTestBackend(t, genBlocks, genesis, func(i int, b *core.BlockGen) {
		// Transfer from account[0] to account[1]
		//    value: 1000 wei
		//    fee:   0 wei
		for j := 0; j < i+1; j++ {
			tx, _ := types.SignTx(types.NewTransaction(nonce, accounts[1].addr, big.NewInt(1000), params.TxGas, b.BaseFee(), nil), signer, accounts[0].key)
			b.AddTx(tx)
			nonce += 1
		}
	})
	backend.refHook = func() { atomic.AddUint32(&ref, 1) }
	backend.relHook = func() { atomic.AddUint32(&rel, 1) }
	api := NewAPI(backend)

	single := `{"result":{"gas":21000,"failed":false,"returnValue":"","structLogs":[]}}`
	var cases = []struct {
		start  uint64
		end    uint64
		config *TraceConfig
	}{
		{0, 50, nil},  // the entire chain range, blocks [1, 50]
		{10, 20, nil}, // the middle chain range, blocks [11, 20]
	}
	for _, c := range cases {
		ref, rel = 0, 0 // clean up the counters

		from, _ := api.blockByNumber(context.Background(), rpc.BlockNumber(c.start))
		to, _ := api.blockByNumber(context.Background(), rpc.BlockNumber(c.end))
		resCh := api.traceChain(context.Background(), from, to, c.config, nil)

		next := c.start + 1
		for result := range resCh {
			if next != uint64(result.Block) {
				t.Error("Unexpected tracing block")
			}
			if len(result.Traces) != int(next) {
				t.Error("Unexpected tracing result")
			}
			for _, trace := range result.Traces {

				blob, _ := json.Marshal(&Result{
					Result: trace.Result,
				})
				if string(blob) != single {
					t.Error("Unexpected tracing result", string(blob))
				}
			}
			next += 1
		}
		if next != c.end+1 {
			t.Error("Missing tracing block")
		}
		if ref != rel {
			t.Errorf("Ref and deref actions are not equal, ref %d rel %d", ref, rel)
		}
	}
}
