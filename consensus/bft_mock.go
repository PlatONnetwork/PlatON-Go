// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package consensus

import (
	"bytes"
	"fmt"
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/event"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"
	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/core/state"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"github.com/PlatONnetwork/PlatON-Go/rpc"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
)

func NewFaker() *BftMock {
	c := new(BftMock)
	c.Blocks = make([]*types.Block, 0)
	c.blockIndexs = make(map[common.Hash]int, 0)
	return c
}

func NewFailFaker(number uint64) *BftMock {
	c := NewFaker()
	c.fakeFail = number
	return c
}

//func NewFakeDelayer(delay time.Duration) *BftMock {
//	c := new(BftMock)
//	c.Blocks = make([]*types.Block, 0)
//	c.fakeDelay = delay
//	return c
//}

// BftMock represents a simulated consensus structure.
type BftMock struct {
	EventMux    *event.TypeMux
	Blocks      []*types.Block
	blockIndexs map[common.Hash]int
	Next        uint32
	Current     *types.Block
	Base        *types.Block
	fakeFail    uint64 // Block number which fails BFT check even in fake mode
	//fakeDelay time.Duration // Time delay to sleep for before returning from verify
}

// InsertChain is a fake interface, no need to implement.
func (bm *BftMock) InsertChain(block *types.Block) error {
	if _, ok := bm.blockIndexs[block.Hash()]; ok {
		return nil
	}

	if len(bm.Blocks) != 0 && bm.Blocks[len(bm.Blocks)-1].Hash() != block.ParentHash() {
		return nil
	}
	bm.Blocks = append(bm.Blocks, block)
	bm.blockIndexs[block.Hash()] = len(bm.Blocks) - 1
	bm.Current = block
	bm.Base = block

	return nil
}

func (bm *BftMock) GetPrepareQC(number uint64) *ctypes.QuorumCert {
	panic("implement me")
}

// FastSyncCommitHead is a fake interface, no need to implement.
func (bm *BftMock) FastSyncCommitHead(block *types.Block) error {
	// todo implement me
	return nil
}

// Start is a fake interface, no need to implement.
func (bm *BftMock) Start(chain ChainReader, blockCacheWriter BlockCacheWriter, pool TxPoolReset, agency Agency) error {
	// todo implement me
	return nil
}

// CalcBlockDeadline is a fake interface, no need to implement.
func (bm *BftMock) CalcBlockDeadline(timePoint time.Time) time.Time {

	now := time.Now()

	if timePoint.Equal(now) || timePoint.Before(now) {
		return now.Add(now.Sub(timePoint)).Add(10 * time.Millisecond)
	}
	return timePoint.Add(10 * time.Millisecond)
}

// CalcNextBlockTime is a fake interface, no need to implement.
func (bm *BftMock) CalcNextBlockTime(timePoint time.Time) time.Time {
	// todo implement me
	return time.Now()
}

// GetBlockWithoutLock is a fake interface, no need to implement.
func (bm *BftMock) GetBlockWithoutLock(hash common.Hash, number uint64) *types.Block {
	// todo implement me
	return nil
}

// IsSignedBySelf is a fake interface, no need to implement.
func (bm *BftMock) IsSignedBySelf(sealHash common.Hash, header *types.Header) bool {
	// todo implement me
	return true
}

// Evidences is a fake interface, no need to implement.
func (bm *BftMock) Evidences() string {
	// todo implement me
	return ""
}

// UnmarshalEvidence is a fake interface, no need to implement.
func (bm *BftMock) UnmarshalEvidence(data []byte) (consensus.Evidences, error) {
	// todo implement me
	return nil, nil
}

func (bm *BftMock) NodeID() discover.NodeID {
	privateKey, err := crypto.GenerateKey()
	if nil != err {
		panic(fmt.Sprintf("Failed to generate random NodeId private key: %v", err))
	}
	return discover.PubkeyID(&privateKey.PublicKey)
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
	if bm.fakeFail == header.Number.Uint64() {
		return fmt.Errorf("failed verifyHeader on bftMock")
	}
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (bm *BftMock) VerifyHeaders(chain ChainReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error) {
	results := make(chan error, len(headers))
	c := make(chan<- struct{})

	//time.Sleep(bm.fakeDelay)
	go func() {
		for i := range headers {
			if bm.fakeFail == headers[i].Number.Uint64() {
				results <- fmt.Errorf("failed verifyHeader on bftMock")
			} else {
				results <- nil
			}

		}
	}()
	return c, results
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (bm *BftMock) VerifySeal(chain ChainReader, header *types.Header) error {
	return nil
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (bm *BftMock) Prepare(chain ChainReader, header *types.Header) error {
	//header.Extra[0:31] to store block's version info etc. and right pad with 0x00;
	//header.Extra[32:] to store block's sign of producer, the length of sign is 65.
	if len(header.Extra) < 32 {
		header.Extra = append(header.Extra, bytes.Repeat([]byte{0x00}, 32-len(header.Extra))...)
	}
	header.Extra = header.Extra[:32]
	//init header.Extra[32: 32+65]
	header.Extra = append(header.Extra, make([]byte, ExtraSeal)...)
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
func (bm *BftMock) Seal(chain ChainReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}, complete chan<- struct{}) error {

	header := block.Header()
	if block.NumberU64() == 0 {
		return fmt.Errorf("unknown block")
	}
	sign := header.SealHash().Bytes()
	copy(header.Extra[len(header.Extra)-ExtraSeal:], sign[:])
	sealBlock := block.WithSeal(header)
	complete <- struct{}{}
	results <- sealBlock
	bm.EventMux.Post(cbfttypes.CbftResult{
		Block: sealBlock,
		//ExtraData:          extra,
		//SyncState:          cbft.commitErrCh,
		ChainStateUpdateCB: func() {
			// Do nothings
			//fmt.Println("result the block", "Number", sealBlock.NumberU64(), "Hash", sealBlock.Hash().Hex())
		},
	})
	return nil
}

// SealHash returns the hash of a block prior to it being sealed.
func (bm *BftMock) SealHash(header *types.Header) common.Hash {
	return header.SealHash()
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
func (bm *BftMock) ShouldSeal(curTime time.Time) (bool, error) {
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
func (bm *BftMock) IsConsensusNode() bool {
	return true
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
	return bm.Base
}

// HasBlock is a fake interface, no need to implement.
func (bm *BftMock) HasBlock(hash common.Hash, number uint64) bool {
	return true
}

// GetBlockByHash is a fake interface, no need to implement.
func (bm *BftMock) GetBlockByHash(hash common.Hash) *types.Block {

	if index, ok := bm.blockIndexs[hash]; ok {
		return bm.Blocks[index]
	}

	return nil
}

// GetBlockByHash get the specified block by hash and number.
func (bm *BftMock) GetBlockByHashAndNum(hash common.Hash, number uint64) *types.Block {
	return nil
}

// Status is a fake interface, no need to implement.
func (bm *BftMock) Status() string {
	return ""
}

// CurrentBlock is a fake interface, no need to implement.
func (bm *BftMock) CurrentBlock() *types.Block {
	//if len(bm.Blocks) == 0 {
	//	h := types.Header{Number: big.NewInt(0)}
	//	return types.NewBlockWithHeader(&h)
	//}
	//return bm.Blocks[len(bm.Blocks)-1]
	return bm.Current
}

// TracingSwitch is a fake interface, no need to implement.
func (bm *BftMock) TracingSwitch(flag int8) {

}

func (bm *BftMock) Pause() {

}
func (bm *BftMock) Resume() {

}

func (bm *BftMock) DecodeExtra(extra []byte) (common.Hash, uint64, error) {
	return common.Hash{}, 0, nil
}
