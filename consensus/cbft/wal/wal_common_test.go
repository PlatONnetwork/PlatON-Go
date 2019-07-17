package wal

import (
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft"
	"math/big"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

var (
	header = &types.Header{
		Number: big.NewInt(1),
	}
	block = types.NewBlock(header, nil, nil)
	ordinal = 0
)

func buildPeerId() discover.NodeID {
	peerId, _ := discover.HexID("b6c8c9f99bfebfa4fb174df720b9385dbd398de699ec36750af3f38f8e310d4f0b90447acbef64bdf924c4b59280f3d42bb256e6123b53e9a7e99e4c432549d6")
	return peerId
}

func buildPrepareBlock() *protocols.PrepareBlock {
	return &protocols.PrepareBlock{
		Epoch: 1,
		ViewNumber: 1,
		Block: block,
		BlockIndex: 1,
		ProposalIndex:1,
		ProposalAddr:common.BytesToAddress(cbft.Rand32Bytes(20)),
		PrepareQC: &ctypes.QuorumCert{
			ViewNumber:1,
			BlockHash:common.BytesToHash(cbft.Rand32Bytes(32)),
			BlockNumber:1,
			Signature: ctypes.Signature{},
		},
		ViewChangeQC  []*ctypes.QuorumCert `json:"viewchange_qc"` //viewchange aggregate signature
		Signature: ctypes.Signature{},
	}
}

func buildPrepareVote() *protocols.PrepareVote {
	return &protocols.PrepareVote{
		Epoch:1,
		ViewNumber:1,
		BlockHash:common.BytesToHash(cbft.Rand32Bytes(32)),
		BlockNumber:1,
		BlockIndex:1,
		ParentQC: ctypes.QuorumCert{
			ViewNumber:1,
			BlockHash:common.BytesToHash(cbft.Rand32Bytes(32)),
			BlockNumber:1,
			Signature: ctypes.Signature{},
		},
		Signature: ctypes.Signature{},
	}
}

func buildViewChange() *viewChange {
	return &viewChange{
		Timestamp:     uint64(time.Now().UnixNano()),
		ProposalIndex: 12,
		ProposalAddr:  common.BytesToAddress(Rand32Bytes(20)),
		BaseBlockNum:  10086,
		BaseBlockHash: common.BytesToHash(Rand32Bytes(32)),
	}
}

func buildSendPrepareBlock() *sendPrepareBlock {
	return &sendPrepareBlock{
		PrepareBlock: buildPrepareBlock(),
	}
}

func buildSendPrepareVote() *sendPrepareVote {
	return &sendPrepareVote{
		PrepareBlock: buildPrepareBlock(),
		PrepareVote:  nil,
	}
}

func buildSendViewChange() *sendViewChange {
	return &sendViewChange{
		ViewChange: buildViewChange(),
	}
}

func buildConfirmedViewChange() *confirmedViewChange {
	votes := make([]*viewChangeVote, 0, 2)
	votes = append(votes, buildviewChangeVote())
	votes = append(votes, buildviewChangeVote())
	return &confirmedViewChange{
		ViewChange: buildViewChange(),
		//ViewChangeResp:  buildviewChangeVote(),
		ViewChangeResp:  nil,
		ViewChangeVotes: votes,
		Master:          true,
	}
}

func ordinalMessages() int {
	if ordinal == len(wal_messages) {
		ordinal = 0
	}

	current := ordinal
	ordinal = ordinal + 1
	return current
}
