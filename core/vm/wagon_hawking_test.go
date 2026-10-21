// Copyright 2026 The PlatON Network Authors
// This file is part of the PlatON-Go library.

package vm

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

func wasmTestContext(caller, self common.Address, gas uint64, hawking bool) *VMContext {
	contract := NewContract(AccountRef(caller), AccountRef(self), big.NewInt(0), gas)
	return &VMContext{
		evm:      &EVM{chainRules: params.Rules{IsHawking: hawking}},
		contract: contract,
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
