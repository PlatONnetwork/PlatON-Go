package downloader

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

// BlockChainWrapper wrap functions required to sync a (full or fast) blockchain.
type BlockChainWrapper struct {
	blockchain BlockChain
	engine     consensus.Engine
}

// NewBlockChainWrapper create a BlockChainWrapper.
func NewBlockChainWrapper(chain BlockChain, engine consensus.Engine) *BlockChainWrapper {
	return &BlockChainWrapper{
		blockchain: chain,
		engine:     engine,
	}
}

// HasHeader verifies a header's presence in the local chain.
func (bw *BlockChainWrapper) HasHeader(hash common.Hash, number uint64) bool {
	return bw.blockchain.HasHeader(hash, number)
}

// GetHeaderByHash retrieves the head header from the local chain.
func (bw *BlockChainWrapper) GetHeaderByHash(hash common.Hash) *types.Header {
	return bw.blockchain.GetHeaderByHash(hash)
}

// CurrentHeader retrieves the head header from the local chain.
func (bw *BlockChainWrapper) CurrentHeader() *types.Header {
	return bw.blockchain.CurrentHeader()
}

// InsertHeaderChain inserts a batch of headers into the local chain.
func (bw *BlockChainWrapper) InsertHeaderChain(headers []*types.Header, checkFreq int) (int, error) {
	return bw.blockchain.InsertHeaderChain(headers, checkFreq)
}

// Rollback removes a few recently added elements from the local chain.
func (bw *BlockChainWrapper) Rollback(chain []common.Hash) {
	bw.blockchain.Rollback(chain)
}

// HasBlock verifies a block's presence in the local chain.
func (bw *BlockChainWrapper) HasBlock(hash common.Hash, number uint64) bool {
	if bw.engine != nil {
		return bw.engine.HasBlock(hash, number)
	}
	return bw.blockchain.HasBlock(hash, number)
}

// GetBlockByHash retrieves a block from the local chain.
func (bw *BlockChainWrapper) GetBlockByHash(hash common.Hash) *types.Block {
	if bw.engine != nil {
		return bw.engine.GetBlockByHash(hash)
	}
	return bw.blockchain.GetBlockByHash(hash)
}

// CurrentBlock retrieves the head block from the local chain.
func (bw *BlockChainWrapper) CurrentBlock() *types.Block {
	if bw.engine != nil {
		return bw.engine.CurrentBlock()
	}
	return bw.blockchain.CurrentBlock()
}

// CurrentFastBlock retrieves the head fast block from the local chain.
func (bw *BlockChainWrapper) CurrentFastBlock() *types.Block {
	return bw.blockchain.CurrentFastBlock()
}

// FastSyncCommitHead directly commits the head block to a certain entity.
func (bw *BlockChainWrapper) FastSyncCommitHead(hash common.Hash) error {
	return bw.blockchain.FastSyncCommitHead(hash)
}

// InsertChain inserts a batch of blocks into the local chain.
func (bw *BlockChainWrapper) InsertChain(blocks types.Blocks) (int, error) {
	return bw.blockchain.InsertChain(blocks)
}

// InsertReceiptChain inserts a batch of receipts into the local chain.
func (bw *BlockChainWrapper) InsertReceiptChain(blocks types.Blocks, receipts []types.Receipts) (int, error) {
	return bw.blockchain.InsertReceiptChain(blocks, receipts)
}
