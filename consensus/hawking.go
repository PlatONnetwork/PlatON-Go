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
// certain is false when the caller cannot load parent state (e.g. HeaderChain
// during snap header sync). Callers must not treat uncertain results as a
// definitive "before Hawking" when rejecting post-fork header fields.
func IsHawkingEnabled(config *params.ChainConfig, number *big.Int, parent *types.Header, stateAt StateAtFunc) (enabled, certain bool) {
	if !config.IsPauli(number) {
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
	if err != nil {
		return false, false
	}
	govInst := gov.NewGov(snapshotdb.Instance())
	if govInst.Gte160VersionState(statedb) {
		return true, true
	}
	gdb := gov.NewGovDB(snapshotdb.Instance())
	if gdb.GetPreActiveVersion(parent.Hash()) < params.FORKVERSION_1_6_0 {
		return false, true
	}
	preActiveID, err := gdb.GetPreActiveProposalID(parent.Hash())
	if err != nil || preActiveID == common.ZeroHash {
		return false, true
	}
	prop, err := gdb.GetExistProposal(preActiveID, statedb)
	if err != nil {
		return false, true
	}
	vp, ok := prop.(*gov.VersionProposal)
	return ok && vp.NewVersion >= params.FORKVERSION_1_6_0 && vp.GetActiveBlock() == number.Uint64(), true
}
