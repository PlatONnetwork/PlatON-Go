package wal

import (
	"math/big"

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
	header         = &types.Header{
		Number: big.NewInt(100),
	}
	block   = types.NewBlock(header, nil, nil)
	ordinal = 0
)

func buildPrepareBlock() *protocols.PrepareBlock {
	return &protocols.PrepareBlock{
		Epoch:         epoch,
		ViewNumber:    viewNumber,
		Block:         block,
		BlockIndex:    blockIndex,
		ProposalIndex: proposalIndex,
		PrepareQC:     buildQuorumCert(),
		ViewChangeQC:  buildViewChangeQC(),
		Signature:     ctypes.BytesToSignature(utils.Rand32Bytes(32)),
	}
}

func buildQuorumCert() *ctypes.QuorumCert {
	return &ctypes.QuorumCert{
		Epoch:        epoch,
		ViewNumber:   viewNumber,
		BlockHash:    common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:  blockNumber,
		BlockIndex:   blockIndex,
		Signature:    ctypes.BytesToSignature(utils.Rand32Bytes(32)),
		ValidatorSet: utils.NewBitArray(110),
	}
}

func buildViewChangeQC() *ctypes.ViewChangeQC {
	return &ctypes.ViewChangeQC{
		QCs: []*ctypes.ViewChangeQuorumCert{{
			Epoch:        epoch,
			ViewNumber:   viewNumber,
			BlockHash:    common.BytesToHash(utils.Rand32Bytes(32)),
			BlockNumber:  blockNumber,
			Signature:    ctypes.BytesToSignature(utils.Rand32Bytes(32)),
			ValidatorSet: utils.NewBitArray(110),
		}},
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
		Signature:      ctypes.BytesToSignature(utils.Rand32Bytes(32)),
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
		Signature:      ctypes.BytesToSignature(utils.Rand32Bytes(32)),
	}
}

func buildSendPrepareBlock() *protocols.SendPrepareBlock {
	return &protocols.SendPrepareBlock{
		Prepare: buildPrepareBlock(),
	}
}

func buildSendPrepareVote() *protocols.SendPrepareVote {
	return &protocols.SendPrepareVote{
		Block: block,
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
		Block:        block,
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
