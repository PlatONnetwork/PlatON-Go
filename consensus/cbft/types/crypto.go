package types

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

type Signature struct {
}

type QuorumCert struct {
	ViewNumber  uint64      `json:"view_number"`
	BlockHash   common.Hash `json:"block_hash"`
	BlockNumber uint64      `json:"block_number"`
	Signature   Signature   `json:"signature"`
}
