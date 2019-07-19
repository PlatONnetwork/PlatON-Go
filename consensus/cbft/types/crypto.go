package types

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

type Signature struct {
}

func (s *Signature) Bytes() []byte {
	return nil
}

type QuorumCert struct {
	ViewNumber  uint64      `json:"view_number"`
	BlockHash   common.Hash `json:"block_hash"`
	BlockNumber uint64      `json:"block_number"`
	BlockIndex  uint32      `json:"block_index"`
	Signature   Signature   `json:"signature"`
}
