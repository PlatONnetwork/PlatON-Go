package plugin

import (
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type BasePlugin interface {
	BeginBlock(header *types.Header, state *state.StateDB) (bool, error)
	EndBlock(header *types.Header, state *state.StateDB) (bool, error)
	Confirmed(block *types.Block) error
}


