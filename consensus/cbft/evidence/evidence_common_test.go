package evidence

import (
	"math/big"
	"math/rand"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
)

var (
	header = &types.Header{
		Number: big.NewInt(1),
	}
	block = types.NewBlock(header, nil, nil)
)

func buildPrepareBlock() *protocols.PrepareBlock {
	viewChangeQC := make([]*ctypes.QuorumCert, 0)
	viewChangeQC = append(viewChangeQC, &ctypes.QuorumCert{
		ViewNumber:  1,
		BlockHash:   common.BytesToHash(Rand32Bytes(32)),
		BlockNumber: 1,
		Signature:   ctypes.Signature{},
	})
	return &protocols.PrepareBlock{
		Epoch:         1,
		ViewNumber:    1,
		Block:         block,
		BlockIndex:    1,
		ProposalIndex: 1,
		ProposalAddr:  common.BytesToAddress(Rand32Bytes(20)),
		PrepareQC: &ctypes.QuorumCert{
			ViewNumber:  1,
			BlockHash:   common.BytesToHash(Rand32Bytes(32)),
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

func Rand32Bytes(n uint32) []byte {
	bs := make([]byte, n)
	for i := 0; i < len(bs); i++ {
		bs[i] = byte(rand.Int31n(int32(n)) & 0xFF)
	}
	return bs
}
