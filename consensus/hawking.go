// Copyright 2021 The PlatON Network Authors
// This file is part of the PlatON-Go library.

package consensus

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/snapshotdb"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/x/gov"
)

// StateAtFunc loads state at a given root hash.
type StateAtFunc func(root common.Hash) (*state.StateDB, error)

// StateAtFromChain returns a StateAtFunc when chain implements StateAt.
func StateAtFromChain(chain interface{}) StateAtFunc {
	type stateAtReader interface {
		StateAt(root common.Hash) (*state.StateDB, error)
	}
	if sr, ok := chain.(stateAtReader); ok {
		return sr.StateAt
	}
	return nil
}

// IsHawkingEnabled reports whether Hawking (FORKVERSION_1_6_0) rules apply to the
// block identified by number and parent header. It mirrors block_validator and
// state_processor checks: config fork height, active governance version from
// parent state, or a version proposal activating on this block.
//
// certain is false when parent state or its governance snapshot is unavailable.
// Callers defer rejection of post-fork fields to full execution in that case.
func IsHawkingEnabled(config *params.ChainConfig, number *big.Int, parent *types.Header, stateAt StateAtFunc) (enabled, certain bool) {
	if !config.IsPauli(number) && config.GetPauliBlock() != nil {
		return false, true
	}
	if config.IsHawking(number) {
		return true, true
	}
	if number == nil || number.Sign() == 0 {
		return false, true
	}
	// Without parent state we cannot evaluate governance activation.
	if parent == nil || stateAt == nil {
		return false, false
	}
	statedb, err := stateAt(parent.Root)
	if err != nil || statedb == nil {
		return false, false
	}
	versions, err := gov.NewGov(nil).GetCurrentActiveVersionList(statedb)
	if err != nil || statedb.Error() != nil {
		return false, false
	}
	if len(versions) > 0 && gov.Gte160Version(versions[0].ActiveVersion) {
		return true, true
	}
	return gov.VersionActivationAt(snapshotdb.Instance(), parent.Hash(), number, statedb)
}
