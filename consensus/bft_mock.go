package consensus

import (
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rpc"
)

// BftMock represents a simulated consensus structure.
type BftMock struct {
}

// InsertChain is a fake interface, no need to implement.
func (bm *BftMock) InsertChain(block *types.Block) error {
	panic("implement me")
}

// FastSyncCommitHead is a fake interface, no need to implement.
func (bm *BftMock) FastSyncCommitHead(block *types.Block) error {
	panic("implement me")
}

// Start is a fake interface, no need to implement.
func (bm *BftMock) Start(chain ChainReader, blockCacheWriter BlockCacheWriter, pool TxPoolReset, agency Agency) error {
	panic("implement me")
}

// CalcBlockDeadline is a fake interface, no need to implement.
func (bm *BftMock) CalcBlockDeadline(timePoint time.Time) time.Time {
	panic("implement me")
}

// CalcNextBlockTime is a fake interface, no need to implement.
func (bm *BftMock) CalcNextBlockTime(timePoint time.Time) time.Time {
	panic("implement me")
}

// GetBlockWithoutLock is a fake interface, no need to implement.
func (bm *BftMock) GetBlockWithoutLock(hash common.Hash, number uint64) *types.Block {
	panic("implement me")
}

// IsSignedBySelf is a fake interface, no need to implement.
func (bm *BftMock) IsSignedBySelf(sealHash common.Hash, header *types.Header) bool {
	panic("implement me")
}

// Evidences is a fake interface, no need to implement.
func (bm *BftMock) Evidences() string {
	panic("implement me")
}

// UnmarshalEvidence is a fake interface, no need to implement.
func (bm *BftMock) UnmarshalEvidence(data []byte) (consensus.Evidences, error) {
	panic("implement me")
}

// Author retrieves the Ethereum address of the account that minted the given
// block, which may be different from the header's coinbase if a consensus
// engine is based on signatures.
func (bm *BftMock) Author(header *types.Header) (common.Address, error) {
	return common.Address{}, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (bm *BftMock) VerifyHeader(chain ChainReader, header *types.Header, seal bool) error {
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (bm *BftMock) VerifyHeaders(chain ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	return nil, nil
}

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (bm *BftMock) VerifyUncles(chain ChainReader, block *types.Block) error {
	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (bm *BftMock) VerifySeal(chain ChainReader, header *types.Header) error {
	return nil
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (bm *BftMock) Prepare(chain ChainReader, header *types.Header) error {
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (bm *BftMock) Finalize(chain ChainReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
	receipts []*types.Receipt) (*types.Block, error) {
	header.Root = state.IntermediateRoot(true)

	// Header seems complete, assemble into a block and return
	return types.NewBlock(header, txs, receipts), nil
}

// Seal generates a new sealing request for the given input block and pushes
// the result into the given channel.
//
// Note, the method returns immediately and will send the result async. More
// than one result may also be returned depending on the consensus algorithm.
func (bm *BftMock) Seal(chain ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error {
	return nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (bm *BftMock) SealHash(header *types.Header) common.Hash {
	return common.Hash{}
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have.
func (bm *BftMock) CalcDifficulty(chain ChainReader, time uint64, parent *types.Header) *big.Int {
	return nil
}

// APIs returns the RPC APIs this consensus engine provides.
func (bm *BftMock) APIs(chain ChainReader) []rpc.API {
	return nil
}

// Protocols is a fake interface, no need to implement.
func (bm *BftMock) Protocols() []p2p.Protocol {
	return []p2p.Protocol{}
}

// Close terminates any background threads maintained by the consensus engine.
func (bm *BftMock) Close() error {
	return nil
}

// ConsensusNodes returns the current consensus node address list.
func (bm *BftMock) ConsensusNodes() ([]discover.NodeID, error) {
	return nil, nil
}

// ShouldSeal returns whether the current node is out of the block
func (bm *BftMock) ShouldSeal(curTime int64) (bool, error) {
	return true, nil
}

// OnBlockSignature received a new block signature
// Need to verify if the signature is signed by nodeID
func (bm *BftMock) OnBlockSignature(chain ChainReader, nodeID discover.NodeID, sig *cbfttypes.BlockSignature) error {
	return nil
}

// OnNewBlock processes the BFT signatures
func (bm *BftMock) OnNewBlock(chain ChainReader, block *types.Block) error {
	return nil
}

// OnPong processes the BFT signatures
func (bm *BftMock) OnPong(nodeID discover.NodeID, netLatency int64) error {
	return nil

}

// OnBlockSynced sends a signal if a block synced from other peer.
func (bm *BftMock) OnBlockSynced() {

}

// CheckConsensusNode is a fake interface, no need to implement.
func (bm *BftMock) CheckConsensusNode(nodeID discover.NodeID) (bool, error) {
	return true, nil
}

// IsConsensusNode is a fake interface, no need to implement.
func (bm *BftMock) IsConsensusNode() (bool, error) {
	return true, nil
}

// HighestLogicalBlock is a fake interface, no need to implement.
func (bm *BftMock) HighestLogicalBlock() *types.Block {
	return nil
}

// HighestConfirmedBlock is a fake interface, no need to implement.
func (bm *BftMock) HighestConfirmedBlock() *types.Block {
	return nil
}

// GetBlock is a fake interface, no need to implement.
func (bm *BftMock) GetBlock(hash common.Hash, number uint64) *types.Block {
	return nil
}

// NextBaseBlock is a fake interface, no need to implement.
func (bm *BftMock) NextBaseBlock() *types.Block {
	return nil
}

// HasBlock is a fake interface, no need to implement.
func (bm *BftMock) HasBlock(hash common.Hash, number uint64) bool {
	return true
}

// GetBlockByHash is a fake interface, no need to implement.
func (bm *BftMock) GetBlockByHash(hash common.Hash) *types.Block {
	return nil
}

// Status is a fake interface, no need to implement.
func (bm *BftMock) Status() string {
	return ""
}

// CurrentBlock is a fake interface, no need to implement.
func (bm *BftMock) CurrentBlock() *types.Block {
	return nil
}

// TracingSwitch is a fake interface, no need to implement.
func (bm *BftMock) TracingSwitch(flag int8) {

}
