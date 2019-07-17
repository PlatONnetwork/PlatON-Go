package wal

import (
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

var (
	header = &types.Header{
		Number: big.NewInt(1),
	}
	block   = types.NewBlock(header, nil, nil)
	ordinal = 0
)

func buildPrepareBlock() *protocols.PrepareBlock {
	viewChangeQC := make([]*ctypes.QuorumCert, 0)
	viewChangeQC = append(viewChangeQC, &ctypes.QuorumCert{
		ViewNumber:  1,
		BlockHash:   common.BytesToHash(cbft.Rand32Bytes(32)),
		BlockNumber: 1,
		Signature:   ctypes.Signature{},
	})
	return &protocols.PrepareBlock{
		Epoch:         1,
		ViewNumber:    1,
		Block:         block,
		BlockIndex:    1,
		ProposalIndex: 1,
		ProposalAddr:  common.BytesToAddress(cbft.Rand32Bytes(20)),
		PrepareQC: &ctypes.QuorumCert{
			ViewNumber:  1,
			BlockHash:   common.BytesToHash(cbft.Rand32Bytes(32)),
			BlockNumber: 1,
			Signature:   ctypes.Signature{},
		},
		ViewChangeQC: viewChangeQC,
		Signature:    ctypes.Signature{},
	}
}

func buildPrepareVote() *protocols.PrepareVote {
	return &protocols.PrepareVote{
		Epoch:       1,
		ViewNumber:  1,
		BlockHash:   common.BytesToHash(cbft.Rand32Bytes(32)),
		BlockNumber: 1,
		BlockIndex:  1,
		ParentQC: &ctypes.QuorumCert{
			ViewNumber:  1,
			BlockHash:   common.BytesToHash(cbft.Rand32Bytes(32)),
			BlockNumber: 1,
			Signature:   ctypes.Signature{},
		},
		Signature: ctypes.Signature{},
	}
}

func buildViewChange() *protocols.ViewChange {
	return &protocols.ViewChange{
		Epoch:      1,
		ViewNumber: 1,
		BlockHash:  common.BytesToHash(cbft.Rand32Bytes(32)),
		BlockNum:   1,
		PrepareQC: &ctypes.QuorumCert{
			ViewNumber:  1,
			BlockHash:   common.BytesToHash(cbft.Rand32Bytes(32)),
			BlockNumber: 1,
			Signature:   ctypes.Signature{},
		},
		Signature: ctypes.Signature{},
	}
}

func buildSendPrepareBlock() *protocols.SendPrepareBlock {
	return &protocols.SendPrepareBlock{
		PrepareBlock: buildPrepareBlock(),
	}
}

func buildSendPrepareVote() *protocols.SendPrepareVote {
	return &protocols.SendPrepareVote{
		PrepareBlock: buildPrepareBlock(),
		PrepareVote:  buildPrepareVote(),
	}
}

func buildSendViewChange() *protocols.SendViewChange {
	return &protocols.SendViewChange{
		ViewChange: buildViewChange(),
	}
}

func buildConfirmedViewChange() *protocols.ConfirmedViewChange {
	return &protocols.ConfirmedViewChange{
		ViewChange: buildViewChange(),
		Master:     true,
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
