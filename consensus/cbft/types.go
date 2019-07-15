package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

type Signature struct {
}

type quorumCert struct {
	ViewNumber  uint64      `json:"view_number"`
	BlockHash   common.Hash `json:"block_hash"`
	BlockNumber uint64      `json:"block_number"`
	Signature   Signature   `json:"signature"`
}

type prepareVotes struct {
	votes []*prepareVote
}

type viewBlocks struct {
	blocks map[uint32]*viewBlock
}

type prepareVoteSet struct {
	votes map[uint32]*prepareVote
}

type viewChanges struct {
	viewChanges map[common.Address]*viewChange
}
