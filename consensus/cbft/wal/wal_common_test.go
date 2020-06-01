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

package wal

import (
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var (
	blockNumber    = uint64(100)
	blockIndex     = uint32(1)
	proposalIndex  = uint32(2)
	validatorIndex = uint32(6)
	ordinal        = 0
)

func newBlock() *types.Block {
	header := &types.Header{
		Number:      big.NewInt(int64(blockNumber)),
		ParentHash:  common.BytesToHash(utils.Rand32Bytes(32)),
		Time:        big.NewInt(time.Now().UnixNano()),
		Extra:       make([]byte, 97),
		ReceiptHash: common.BytesToHash(utils.Rand32Bytes(32)),
		Root:        common.BytesToHash(utils.Rand32Bytes(32)),
	}
	block := types.NewBlockWithHeader(header)
	return block
}

func buildPrepareBlock() *protocols.PrepareBlock {
	return &protocols.PrepareBlock{
		Epoch:         epoch,
		ViewNumber:    viewNumber,
		Block:         newBlock(),
		BlockIndex:    blockIndex,
		ProposalIndex: proposalIndex,
		PrepareQC:     buildQuorumCert(),
		ViewChangeQC:  buildViewChangeQC(),
		Signature:     ctypes.BytesToSignature(utils.Rand32Bytes(64)),
	}
}

func buildQuorumCert() *ctypes.QuorumCert {
	return &ctypes.QuorumCert{
		Epoch:        epoch,
		ViewNumber:   viewNumber,
		BlockHash:    common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:  blockNumber,
		BlockIndex:   blockIndex,
		Signature:    ctypes.BytesToSignature(utils.Rand32Bytes(64)),
		ValidatorSet: utils.NewBitArray(25),
	}
}

func buildViewChangeQuorumCert(epoch, viewNumber uint64, blockNumber uint64) *ctypes.ViewChangeQuorumCert {
	return &ctypes.ViewChangeQuorumCert{
		Epoch:        epoch,
		ViewNumber:   viewNumber,
		BlockHash:    common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:  blockNumber,
		Signature:    ctypes.BytesToSignature(utils.Rand32Bytes(64)),
		ValidatorSet: utils.NewBitArray(25),
	}
}

func buildViewChangeQC() *ctypes.ViewChangeQC {
	return &ctypes.ViewChangeQC{
		QCs: []*ctypes.ViewChangeQuorumCert{
			buildViewChangeQuorumCert(epoch, viewNumber, blockNumber),
			buildViewChangeQuorumCert(epoch, viewNumber, blockNumber),
			buildViewChangeQuorumCert(epoch, viewNumber, blockNumber),
		},
	}
}

func buildPrepareVote() *protocols.PrepareVote {
	return &protocols.PrepareVote{
		Epoch:          epoch,
		ViewNumber:     viewNumber,
		BlockHash:      common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:    blockNumber,
		BlockIndex:     blockIndex,
		ValidatorIndex: validatorIndex,
		ParentQC:       buildQuorumCert(),
		Signature:      ctypes.BytesToSignature(utils.Rand32Bytes(64)),
	}
}

func buildViewChange() *protocols.ViewChange {
	return &protocols.ViewChange{
		Epoch:          epoch,
		ViewNumber:     viewNumber,
		BlockHash:      common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:    blockNumber,
		ValidatorIndex: validatorIndex,
		PrepareQC:      buildQuorumCert(),
		Signature:      ctypes.BytesToSignature(utils.Rand32Bytes(64)),
	}
}

func buildSendPrepareBlock() *protocols.SendPrepareBlock {
	return &protocols.SendPrepareBlock{
		Prepare: buildPrepareBlock(),
	}
}

func buildSendPrepareVote() *protocols.SendPrepareVote {
	return &protocols.SendPrepareVote{
		Block: newBlock(),
		Vote:  buildPrepareVote(),
	}
}

func buildSendViewChange() *protocols.SendViewChange {
	return &protocols.SendViewChange{
		ViewChange: buildViewChange(),
	}
}

func buildConfirmedViewChange() *protocols.ConfirmedViewChange {
	return &protocols.ConfirmedViewChange{
		Epoch:        epoch,
		ViewNumber:   viewNumber,
		Block:        newBlock(),
		QC:           buildQuorumCert(),
		ViewChangeQC: buildViewChangeQC(),
	}
}

func ordinalMessages() int {
	if ordinal == len(protocols.WalMessages) {
		ordinal = 0
	}

	current := ordinal
	ordinal = ordinal + 1
	return current
}
