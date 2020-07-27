// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header verification.
type ChainReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *types.Header

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash common.Hash) *types.Header

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *types.Block

	// CurrentBlock retrieves the current head block of the canonical chain.
	CurrentBlock() *types.Block
}

// TxPoolReset stands for transaction pool.
type TxPoolReset interface {
	ForkedReset(newHeader *types.Header, rollback []*types.Block)
	Reset(newBlock *types.Block)
}

// BlockCacheWriter executions block, you need to pass in the parent
// block to find the parent block state
type BlockCacheWriter interface {
	Execute(block *types.Block, parent *types.Block) error
	ClearCache(block *types.Block)
	WriteBlock(block *types.Block) error
}

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainReader, header *types.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(chain ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifySeal checks whether the crypto seal on a header is valid according to
	// the consensus rules of the given engine.
	VerifySeal(chain ChainReader, header *types.Header) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(chain ChainReader, header *types.Header) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// and assembles the final block.
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(chain ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		receipts []*types.Receipt) (*types.Block, error)

	// Seal generates a new sealing request for the given input block and pushes
	// the result into the given channel.
	//
	// Note, the method returns immediately and will send the result async. More
	// than one result may also be returned depending on the consensus algorithm.
	Seal(chain ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}, complete chan<- struct{}) error

	// SealHash returns the hash of a block prior to it being sealed.
	SealHash(header *types.Header) common.Hash

	// APIs returns the RPC APIs this consensus engine provides.
	APIs(chain ChainReader) []rpc.API

	Protocols() []p2p.Protocol

	NextBaseBlock() *types.Block

	InsertChain(block *types.Block) error

	HasBlock(hash common.Hash, number uint64) bool

	GetBlockByHash(hash common.Hash) *types.Block

	GetBlockByHashAndNum(hash common.Hash, number uint64) *types.Block

	CurrentBlock() *types.Block

	FastSyncCommitHead(block *types.Block) error

	// Close terminates any background threads maintained by the consensus engine.
	Close() error

	// Pause consensus
	Pause()
	// Resume consensus
	Resume()

	DecodeExtra(extra []byte) (common.Hash, uint64, error)
}

// PoW is a consensus engine based on proof-of-work.
type PoW interface {
	Engine

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	Hashrate() float64
}

// Agency defines the interface that the authentication
// mechanism must implement.
type Agency interface {
	Sign(msg interface{}) error
	VerifySign(msg interface{}) error
	Flush(header *types.Header) error
	VerifyHeader(header *types.Header, stateDB *state.StateDB) error
	GetLastNumber(blockNumber uint64) uint64
	GetValidator(blockNumber uint64) (*cbfttypes.Validators, error)
	IsCandidateNode(nodeID discover.NodeID) bool
	OnCommit(block *types.Block) error
}

// Bft defines the functions that BFT consensus
// must implement.
type Bft interface {
	Engine

	Start(chain ChainReader, blockCacheWriter BlockCacheWriter, pool TxPoolReset, agency Agency) error

	// Returns the current consensus node address list.
	ConsensusNodes() ([]discover.NodeID, error)

	// Returns whether the current node is out of the block
	ShouldSeal(curTime time.Time) (bool, error)

	CalcBlockDeadline(timePoint time.Time) time.Time

	CalcNextBlockTime(timePoint time.Time) time.Time

	IsConsensusNode() bool

	GetBlock(hash common.Hash, number uint64) *types.Block

	GetBlockWithoutLock(hash common.Hash, number uint64) *types.Block

	IsSignedBySelf(sealHash common.Hash, header *types.Header) bool

	Evidences() string

	TracingSwitch(flag int8)

	// NodeID is temporary.
	NodeID() discover.NodeID
}
