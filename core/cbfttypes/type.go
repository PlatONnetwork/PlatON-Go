package cbfttypes

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// Block's Signature info
type BlockSignature struct {
	SignHash  common.Hash // Signature hash，header[0:32]
	Hash      common.Hash // Block hash，header[:]
	Number    *big.Int
	Signature *common.BlockConfirmSign
}

func (bs *BlockSignature) Copy() *BlockSignature {
	sign := *bs.Signature
	return &BlockSignature{
		SignHash:  bs.SignHash,
		Hash:      bs.Hash,
		Number:    new(big.Int).Set(bs.Number),
		Signature: &sign,
	}
}

type CbftResult struct {
	Block     *types.Block
	ExtraData []byte
	SyncState chan error
}

type ProducerState struct {
	count int
	miner common.Address
}

func (ps *ProducerState) Add(miner common.Address) {
	if ps.miner == miner {
		ps.count++
	} else {
		ps.miner = miner
		ps.count = 1
	}
}

func (ps *ProducerState) Get() (common.Address, int) {
	return ps.miner, ps.count
}

func (ps *ProducerState) Validate(period int) bool {
	return ps.count < period
}

type AddValidatorEvent struct {
	NodeID discover.NodeID
}

type RemoveValidatorEvent struct {
	NodeID discover.NodeID
}

type UpdateValidatorEvent struct{}
