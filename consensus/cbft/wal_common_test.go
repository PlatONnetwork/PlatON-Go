package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
	"time"
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

func buildPrepareBlock() *prepareBlock {
	viewChangeVotes := make([]*viewChangeVote, 0)
	viewChangeVotes = append(viewChangeVotes, &viewChangeVote{
		Timestamp:      uint64(time.Now().UnixNano()),
		BlockNum:       111,
		BlockHash:      common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065747"),
		ProposalIndex:  1111,
		ProposalAddr:   common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185822"),
		ValidatorIndex: 11111,
		ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185822"),
	})
	viewChangeVotes = append(viewChangeVotes, &viewChangeVote{
		Timestamp:      uint64(time.Now().UnixNano()),
		BlockNum:       222,
		BlockHash:      common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065749"),
		ProposalIndex:  2222,
		ProposalAddr:   common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185829"),
		ValidatorIndex: 22222,
		ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185829"),
	})
	return &prepareBlock{
		Timestamp:     uint64(time.Now().UnixNano()),
		Block:         block,
		ProposalIndex: 666,
		View: &viewChange{
			Timestamp:     uint64(time.Now().UnixNano()),
			ProposalIndex: 12,
			ProposalAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185822"),
			BaseBlockNum:  10086,
			BaseBlockHash: common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		},
		ViewChangeVotes: viewChangeVotes,
	}
}

func buildHighestPrepareBlock() *highestPrepareBlock {
	pvs := make([]*prepareVote, 0)
	pvs = append(pvs, &prepareVote{
		Timestamp: uint64(time.Now().UnixNano()),
		Number:    7777,
	})
	votes := make([]*prepareVotes, 0)
	votes = append(votes, &prepareVotes{
		Hash:   common.HexToHash("0x76fded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		Number: 5678,
		Votes:  pvs,
	})
	votes = append(votes, &prepareVotes{
		Hash:   common.HexToHash("0x76fded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		Number: 6789,
		Votes:  pvs,
	})
	return &highestPrepareBlock{
		Votes: votes,
	}
}

func buildPrepareVote() *prepareVote {
	return &prepareVote{
		Timestamp:      uint64(time.Now().UnixNano()),
		Hash:           common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066326"),
		Number:         16666,
		ValidatorIndex: 1,
		ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185821"),
	}
}

func buildPrepareVotes() *prepareVotes {
	votes := make([]*prepareVote, 0)
	votes = append(votes, &prepareVote{
		Timestamp:      uint64(time.Now().UnixNano()),
		Hash:           common.HexToHash("0x76fded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		Number:         8877,
		ValidatorIndex: 9900,
		ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185827"),
	})
	votes = append(votes, &prepareVote{
		Timestamp:      uint64(time.Now().UnixNano()),
		Hash:           common.HexToHash("0x76fded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		Number:         8878,
		ValidatorIndex: 9901,
		ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185828"),
	})
	return &prepareVotes{
		Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		Number: 7788,
		Votes:  votes,
	}
}

func buildViewChange() *viewChange {
	return &viewChange{
		Timestamp:     uint64(time.Now().UnixNano()),
		ProposalIndex: 12,
		ProposalAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185822"),
		BaseBlockNum:  10086,
		BaseBlockHash: common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
	}
}

func buildConfirmedPrepareBlock() *confirmedPrepareBlock {
	return &confirmedPrepareBlock{
		Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		Number: 7788,
		VoteBits: NewBitArray(110),
	}
}

func buildGetPrepareVote() *getPrepareVote {
	return &getPrepareVote{
		Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		Number: 7788,
		VoteBits: NewBitArray(110),
	}
}

func buildGetPrepareBlock() *getPrepareBlock {
	return &getPrepareBlock{
		Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		Number: 7788,
	}
}

func buildGetHighestPrepareBlock() *getHighestPrepareBlock {
	return &getHighestPrepareBlock{
		Lowest: 1,
	}
}

func buildCbftStatusData() *cbftStatusData {
	return &cbftStatusData{
		BN: big.NewInt(999),
		CurrentBlock: common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
	}
}

func buildPrepareBlockHash() *prepareBlockHash {
	return &prepareBlockHash{
		Hash:   common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df066329"),
		Number: 7788,
	}
}

func buildviewChangeVote() *viewChangeVote {
	return &viewChangeVote{
		Timestamp:      uint64(time.Now().UnixNano()),
		BlockNum:       222,
		BlockHash:      common.HexToHash("0x8bfded8b3ccdd1d31bf049b4abf72415a0cc829cdcc0b750a73e0da5df065749"),
		ProposalIndex:  2222,
		ProposalAddr:   common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185829"),
		ValidatorIndex: 22222,
		ValidatorAddr:  common.HexToAddress("0x493301712671ada506ba6ca7891f436d29185829"),
	}
}

func ordinalMessages() int {
	if ordinal == len(messages) {
		ordinal = 0
	}

	current := ordinal
	ordinal = ordinal + 1
	return current
}
