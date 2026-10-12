// Copyright 2026 The PlatON Network Authors
// This file is part of the PlatON-Go library.

package consensus

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/params"
)

func TestIsHawkingEnabledUncertainWithoutStateAt(t *testing.T) {
	cfg := &params.ChainConfig{
		GenesisVersion: params.FORKVERSION_1_5_0,
		PauliBlock:     big.NewInt(0),
		// HawkingBlock unset: activation depends on governance state.
	}
	parent := &types.Header{
		Number: big.NewInt(100),
		Root:   common.HexToHash("0x01"),
	}
	number := big.NewInt(101)

	enabled, certain := IsHawkingEnabled(cfg, number, parent, nil)
	if enabled {
		t.Fatal("expected enabled=false without state")
	}
	if certain {
		t.Fatal("expected certain=false when stateAt is nil (snap header sync)")
	}
}

func TestIsHawkingEnabledCertainWithConfigFork(t *testing.T) {
	cfg := &params.ChainConfig{
		GenesisVersion: params.FORKVERSION_1_5_0,
		PauliBlock:     big.NewInt(0),
		HawkingBlock:   big.NewInt(0),
	}
	enabled, certain := IsHawkingEnabled(cfg, big.NewInt(1), nil, nil)
	if !enabled || !certain {
		t.Fatalf("config HawkingBlock should be enabled&certain, got enabled=%v certain=%v", enabled, certain)
	}
}

func TestIsHawkingEnabledCertainBeforePauli(t *testing.T) {
	cfg := &params.ChainConfig{
		GenesisVersion: params.FORKVERSION_1_4_0,
	}
	enabled, certain := IsHawkingEnabled(cfg, big.NewInt(1), nil, nil)
	if enabled || !certain {
		t.Fatalf("before Pauli should be disabled&certain, got enabled=%v certain=%v", enabled, certain)
	}
}
