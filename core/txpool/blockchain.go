package txpool

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/event"
)

// txPoolBlockChain provides chain state to the transaction pool.
type txPoolBlockChain interface {
	CurrentBlock() *types.Header
	GetBlock(hash common.Hash, number uint64) *types.Block
	GetState(header *types.Header) (*state.StateDB, error)
}

// TxPoolBlockChain differs from BlockChain in that CurrentBlock returns the
// highest logical block when the engine is BFT.
type TxPoolBlockChain struct {
	chain *core.BlockChainCache
}

func NewTxPoolBlockChain(chain *core.BlockChainCache) *TxPoolBlockChain {
	return &TxPoolBlockChain{chain: chain}
}

func (tx *TxPoolBlockChain) CurrentBlock() *types.Header {
	block := tx.chain.Engine().CurrentBlock()
	if block != nil {
		return block.Header()
	}
	return tx.chain.BlockChain.CurrentBlock()
}

func (tx *TxPoolBlockChain) GetBlock(hash common.Hash, number uint64) *types.Block {
	return tx.chain.GetBlockInMemory(hash, number)
}

func (tx *TxPoolBlockChain) GetState(header *types.Header) (*state.StateDB, error) {
	return tx.chain.GetState(header)
}

func (tx *TxPoolBlockChain) SubscribeChainHeadEvent(ch chan<- core.ChainHeadEvent) event.Subscription {
	return tx.chain.BlockChain.SubscribeChainHeadEvent(ch)
}
