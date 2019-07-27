package wal

import (
	"math/big"
	"math/rand"

	"github.com/PlatONnetwork/PlatON-Go/common"
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
	return &protocols.PrepareBlock{
		Epoch:         1,
		ViewNumber:    1,
		Block:         block,
		BlockIndex:    1,
		ProposalIndex: 1,
		PrepareQC: &ctypes.QuorumCert{
			ViewNumber:  1,
			BlockHash:   common.BytesToHash(Rand32Bytes(32)),
			BlockNumber: 1,
			Signature:   ctypes.Signature{},
		},
		ViewChangeQC: &ctypes.ViewChangeQC{},
		Signature:    ctypes.Signature{},
	}
}

func buildQuorumCert() *ctypes.QuorumCert {
	return &ctypes.QuorumCert{
		ViewNumber:  viewNumber,
		BlockHash:   common.BytesToHash(Rand32Bytes(32)),
		BlockNumber: block.NumberU64(),
	}
}

func buildPrepareVote() *protocols.PrepareVote {
	return &protocols.PrepareVote{
		Epoch:       1,
		ViewNumber:  1,
		BlockHash:   common.BytesToHash(Rand32Bytes(32)),
		BlockNumber: 1,
		BlockIndex:  1,
		ParentQC: &ctypes.QuorumCert{
			ViewNumber:  1,
			BlockHash:   common.BytesToHash(Rand32Bytes(32)),
			BlockNumber: 1,
			Signature:   ctypes.Signature{},
		},
		Signature: ctypes.Signature{},
	}
}

func buildViewChange() *protocols.ViewChange {
	return &protocols.ViewChange{
		Epoch:       1,
		ViewNumber:  1,
		BlockHash:   common.BytesToHash(Rand32Bytes(32)),
		BlockNumber: 1,
		PrepareQC: &ctypes.QuorumCert{
			ViewNumber:  1,
			BlockHash:   common.BytesToHash(Rand32Bytes(32)),
			BlockNumber: 1,
			Signature:   ctypes.Signature{},
		},
		Signature: ctypes.Signature{},
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
		Epoch:      epoch,
		ViewNumber: viewNumber,
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

func Rand32Bytes(n uint32) []byte {
	bs := make([]byte, n)
	for i := 0; i < len(bs); i++ {
		bs[i] = byte(rand.Int31n(int32(n)) & 0xFF)
	}
	return bs
}
