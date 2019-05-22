package pposm

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"math/big"
)

type CandidatePoolContext struct {
	Configs *params.PposConfig
}

var cContext *CandidatePoolContext

// Initialize the global candidate pool context object
func NewCandidatePoolContext(configs *params.PposConfig) *CandidatePoolContext {
	cContext = &CandidatePoolContext{
		Configs: configs,
	}
	return cContext
}

func GetCandidateContextPtr() *CandidatePoolContext {
	return cContext
}

func (c *CandidatePoolContext) initCandidatePool() *CandidatePool {
	return NewCandidatePool(c.Configs)
}

func (c *CandidatePoolContext) SetCandidate(state vm.StateDB, nodeId discover.NodeID, can *types.Candidate) error {
	return c.initCandidatePool().SetCandidate(state, nodeId, can)
}

func (c *CandidatePoolContext) GetCandidate(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) *types.Candidate {
	return c.initCandidatePool().GetCandidate(state, nodeId, blockNumber)
}

func (c *CandidatePoolContext) GetCandidateArr(state vm.StateDB, blockNumber *big.Int, nodeIds ...discover.NodeID) types.CandidateQueue {
	return c.initCandidatePool().GetCandidateArr(state, blockNumber, nodeIds...)
}

func (c *CandidatePoolContext) GetWitnessCandidate(state vm.StateDB, nodeId discover.NodeID, flag int, blockNumber *big.Int) *types.Candidate {
	return c.initCandidatePool().GetWitnessCandidate(state, nodeId, flag, blockNumber)
}

func (c *CandidatePoolContext) WithdrawCandidate(state vm.StateDB, nodeId discover.NodeID, price, blockNumber *big.Int) error {
	return c.initCandidatePool().WithdrawCandidate(state, nodeId, price, blockNumber)
}

func (c *CandidatePoolContext) GetChosens(state vm.StateDB, flag int, blockNumber *big.Int) types.KindCanQueue {
	return c.initCandidatePool().GetChosens(state, flag, blockNumber)
}

func (c *CandidatePoolContext) GetCandidatePendArr (state vm.StateDB, flag int, blockNumber *big.Int) types.CandidateQueue {
	return c.initCandidatePool().GetCandidatePendArr(state, flag, blockNumber)
}

func (c *CandidatePoolContext) GetChairpersons(state vm.StateDB, blockNumber *big.Int) types.CandidateQueue {
	return c.initCandidatePool().GetChairpersons(state, blockNumber)
}

func (c *CandidatePoolContext) GetDefeat(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) types.RefundQueue {
	return c.initCandidatePool().GetDefeat(state, nodeId, blockNumber)
}

func (c *CandidatePoolContext) IsDefeat(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) bool {
	return c.initCandidatePool().IsDefeat(state, nodeId, blockNumber)
}

func (c *CandidatePoolContext) IsChosens(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) bool {
	return c.initCandidatePool().IsChosens(state, nodeId, blockNumber)
}

func (c *CandidatePoolContext) RefundBalance(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) error {
	return c.initCandidatePool().RefundBalance(state, nodeId, blockNumber)
}

func (c *CandidatePoolContext) GetOwner(state vm.StateDB, nodeId discover.NodeID, blockNumber *big.Int) common.Address {
	return c.initCandidatePool().GetOwner(state, nodeId, blockNumber)
}

func (c *CandidatePoolContext) GetRefundInterval(blockNumber *big.Int) uint32 {
	return c.initCandidatePool().GetRefundInterval(blockNumber)
}

func (c *CandidatePoolContext) MaxCount() uint32 {
	return c.initCandidatePool().MaxCount()
}

func (c *CandidatePoolContext) MaxChair() uint32 {
	return c.initCandidatePool().MaxChair()
}

func (c *CandidatePoolContext) Election(state *state.StateDB, parentHash common.Hash, blocknumber *big.Int) ([]*discover.Node, error) {
	return c.initCandidatePool().Election(state, parentHash, blocknumber)
}

func (c *CandidatePoolContext) Switch(state *state.StateDB, blockNumber *big.Int) bool {
	return c.initCandidatePool().Switch(state, blockNumber)
}

func (c *CandidatePoolContext) GetWitness(state *state.StateDB, flag int, blockNumber *big.Int) ([]*discover.Node, error) {
	return c.initCandidatePool().GetWitness(state, flag, blockNumber)
}

func (c *CandidatePoolContext) GetAllWitness(state *state.StateDB, blockNumber *big.Int) ([]*discover.Node, []*discover.Node, []*discover.Node, error) {
	return c.initCandidatePool().GetAllWitness(state, blockNumber)
}

func (c *CandidatePoolContext) SetCandidateExtra(state vm.StateDB, nodeId discover.NodeID, extra string) error {
	return c.initCandidatePool().SetCandidateExtra(state, nodeId, extra)
}

func (c *CandidatePoolContext) UpdateElectedQueue(state vm.StateDB, currBlockNumber *big.Int, nodeIds ...discover.NodeID) error {
	return c.initCandidatePool().UpdateElectedQueue(state, currBlockNumber, nodeIds...)
}
