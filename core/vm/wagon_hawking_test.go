// Copyright 2026 The PlatON Network Authors
// This file is part of the PlatON-Go library.

package vm

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/rawdb"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
	"github.com/PlatONnetwork/wagon/exec"
)

func TestWasmMigrationBufferGas(t *testing.T) {
	for _, active := range []bool{false, true} {
		ctx := wasmTestContext(common.Address{1}, common.Address{2}, 100, active)
		ctx.evm.depth = int(params.CallCreateDepth) + 1
		machine := newTestVM(ctx.evm)
		machine.SetHostCtx(ctx)
		if status := MigrateContract(exec.NewProcess(machine), 0, 0, 3, 0, 5, 0, 7); status != -1 {
			t.Fatalf("unexpected migration status %d", status)
		}
		want := uint64(100) - ctx.gasTable.SLoad
		if active {
			want -= 3 + 5 + 7
		}
		if ctx.contract.Gas != want {
			t.Fatalf("Hawking=%t: remaining gas %d, want %d", active, ctx.contract.Gas, want)
		}
	}
}

func wasmTestContext(caller, self common.Address, gas uint64, hawking bool) *VMContext {
	contract := NewContract(AccountRef(caller), AccountRef(self), big.NewInt(0), gas)
	return &VMContext{
		evm:      &EVM{chainRules: params.Rules{IsHawking: hawking}},
		contract: contract,
	}
}

func TestWasmHawkingFromStateOncePerFrame(t *testing.T) {
	db, _ := state.New(common.Hash{}, state.NewDatabase(rawdb.NewMemoryDatabase()), nil)
	gdb := gov.NewGovDB(snapshotdb.Instance())
	if err := gdb.AddActiveVersion(params.FORKVERSION_1_5_0, 0, db); err != nil {
		t.Fatal(err)
	}
	evm := &EVM{StateDB: db, SnapshotDB: snapshotdb.Instance()}

	ctx := &VMContext{evm: evm}
	if wasmHawking(ctx) {
		t.Fatal("Hawking reported at 1.5.0")
	}
	if err := gdb.AddActiveVersion(params.FORKVERSION_1_6_0, 1, db); err != nil {
		t.Fatal(err)
	}
	if wasmHawking(ctx) {
		t.Fatal("frame result changed after its first lookup")
	}
	if !wasmHawking(&VMContext{evm: evm}) {
		t.Fatal("new frame did not see 1.6.0")
	}
}

func TestWasmDeployPayerHawking(t *testing.T) {
	caller := common.HexToAddress("0x11")
	self := common.HexToAddress("0x22")

	pre := wasmTestContext(caller, self, 1000, false)
	if got := wasmDeployPayer(pre); got != caller {
		t.Fatalf("pre-Hawking payer = %s, want caller %s", got.Hex(), caller.Hex())
	}

	post := wasmTestContext(caller, self, 1000, true)
	if got := wasmDeployPayer(post); got != self {
		t.Fatalf("Hawking payer = %s, want contract %s", got.Hex(), self.Hex())
	}
}

func TestWasmReadOnlyCallRejectedAfterHawking(t *testing.T) {
	ctx := wasmTestContext(common.HexToAddress("0x11"), common.HexToAddress("0x22"), 1000, true)
	ctx.readOnly = true
	defer func() {
		err, ok := recover().(error)
		if !ok || err != ErrWASMWriteProtection {
			t.Fatalf("recover = %v, want ErrWASMWriteProtection", err)
		}
	}()
	rejectWasmWriteInReadOnly(ctx)
	t.Fatal("expected panic")
}

func TestWasmReadOnlyCallAllowedBeforeHawking(t *testing.T) {
	ctx := wasmTestContext(common.HexToAddress("0x11"), common.HexToAddress("0x22"), 1000, false)
	ctx.readOnly = true
	rejectWasmWriteInReadOnly(ctx)
}

func TestLimitWasmBuffer(t *testing.T) {
	ctx := wasmTestContext(common.HexToAddress("0x11"), common.HexToAddress("0x22"), 100, true)
	defer func() {
		if recover() != ErrOutOfGas {
			t.Fatal("expected ErrOutOfGas")
		}
	}()
	limitWasmBuffer(ctx, 101)
	t.Fatal("expected panic")
}

func TestLimitWasmBufferPreHawking(t *testing.T) {
	ctx := wasmTestContext(common.HexToAddress("0x11"), common.HexToAddress("0x22"), 100, false)
	limitWasmBuffer(ctx, 1<<31)
}

func TestLimitWasmBufferCapsLinearMemory(t *testing.T) {
	ctx := wasmTestContext(common.HexToAddress("0x11"), common.HexToAddress("0x22"), uint64(memoryLimit)+10, true)
	defer func() {
		if recover() != ErrOutOfGas {
			t.Fatal("expected ErrOutOfGas")
		}
	}()
	limitWasmBuffer(ctx, uint32(memoryLimit)+1)
	t.Fatal("expected panic")
}
