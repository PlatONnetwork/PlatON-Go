package core

import (
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type BasePlugin interface {
	GetInstance() *BasePlugin
	BeginBlock(header *types.Header, state *state.StateDB) (bool, error)
	EndBlock(header *types.Header, state *state.StateDB) (bool, error)
	Confirmed(block *types.Block) error
}


