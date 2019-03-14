package pposm

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/core/vm"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"math/big"
)

type TicketPoolContext struct {
	Configs *params.PposConfig
}

var tContext *TicketPoolContext

// Initialize the global ticket pool context object
func NewTicketPoolContext(configs *params.PposConfig) *TicketPoolContext {
	tContext = &TicketPoolContext{
		Configs: configs,
	}
	return tContext
}

func GetTicketPoolContextPtr() *TicketPoolContext {
	return tContext
}

func (c *TicketPoolContext) initTicketPool() *TicketPool {
	return NewTicketPool(c.Configs)
}

func (c *TicketPoolContext) GetPoolNumber (state vm.StateDB) uint32 {
	return c.initTicketPool().GetPoolNumber(state)
}

func (c *TicketPoolContext) VoteTicket (state vm.StateDB, owner common.Address, voteNumber uint32, deposit *big.Int, nodeId discover.NodeID, blockNumber *big.Int) (uint32, error) {
	return c.initTicketPool().VoteTicket(state, owner, voteNumber, deposit, nodeId, blockNumber)
}

func (c *TicketPoolContext) GetTicket(state vm.StateDB, ticketId common.Hash) *types.Ticket {
	return c.initTicketPool().GetTicket(state, ticketId)
}

func (c *TicketPoolContext) GetExpireTicketIds(state vm.StateDB, blockNumber *big.Int) []common.Hash {
	return c.initTicketPool().GetExpireTicketIds(state, blockNumber)
}

func (c *TicketPoolContext) GetTicketList (state vm.StateDB, ticketIds []common.Hash) []*types.Ticket {
	return c.initTicketPool().GetTicketList(state, ticketIds)
}

func (c *TicketPoolContext) GetCandidateTicketIds (state vm.StateDB, nodeId discover.NodeID) []common.Hash {
	return c.initTicketPool().GetCandidateTicketIds(state, nodeId)
}

func (c *TicketPoolContext) GetCandidateEpoch (state vm.StateDB, nodeId discover.NodeID) uint64 {
	return c.initTicketPool().GetCandidateEpoch(state, nodeId)
}

func (c *TicketPoolContext) GetTicketPrice (state vm.StateDB) *big.Int {
	return c.initTicketPool().GetTicketPrice(state)
}

func (c *TicketPoolContext) Notify (state vm.StateDB, blockNumber *big.Int) error {
	return c.initTicketPool().Notify(state, blockNumber)
}

func (c *TicketPoolContext) StoreHash (state vm.StateDB) error {
	return c.initTicketPool().CommitHash(state)
}

func (c *TicketPoolContext) GetCandidateTicketCount (state vm.StateDB, nodeId discover.NodeID) uint32 {
	return c.initTicketPool().GetCandidateTicketCount(state, nodeId)
}

func (c *TicketPoolContext) GetCandidatesTicketCount (state vm.StateDB, nodeIds []discover.NodeID) map[discover.NodeID]int {
	return c.initTicketPool().GetCandidatesTicketCount(state, nodeIds)
}

func (c *TicketPoolContext) GetCandidatesTicketIds (state vm.StateDB, nodeIds []discover.NodeID) map[discover.NodeID][]common.Hash {
	return c.initTicketPool().GetCandidatesTicketIds(state, nodeIds)
}

func (c *TicketPoolContext) DropReturnTicket(stateDB vm.StateDB, blockNumber *big.Int, nodeIds ...discover.NodeID) error {
	return c.initTicketPool().DropReturnTicket(stateDB, blockNumber, nodeIds...)
}

func (c *TicketPoolContext) ReturnTicket(stateDB vm.StateDB, nodeId discover.NodeID, ticketId common.Hash, blockNumber *big.Int) error {
	return c.initTicketPool().ReturnTicket(stateDB, nodeId, ticketId, blockNumber)
}

func (c *TicketPoolContext) SelectionLuckyTicket(stateDB vm.StateDB, nodeId discover.NodeID, blockHash common.Hash) (common.Hash, error) {
	return c.initTicketPool().SelectionLuckyTicket(stateDB, nodeId, blockHash)
}





