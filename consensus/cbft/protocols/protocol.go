package protocols

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

type PrepareBlock struct {
	Epoch         uint64               `json:"epoch"`
	ViewNumber    uint64               `json:"view_number"`
	Block         *types.Block         `json:"block_hash"`
	BlockIndex    uint32               `json:"block_index"` //The block number of the current ViewNumber proposal, 0....10
	ProposalIndex uint32               `json:"proposal_index"`
	ProposalAddr  common.Address       `json:"proposal_address"`
	PrepareQC     *ctypes.QuorumCert   `json:"prepare_qc"`    //N-f aggregate signature
	ViewChangeQC  []*ctypes.QuorumCert `json:"viewchange_qc"` //viewchange aggregate signature
	Signature     ctypes.Signature     `json:"signature"`
}

func (PrepareBlock) String() string {
	panic("implement me")
}

func (PrepareBlock) MsgHash() common.Hash {
	panic("implement me")
}

func (PrepareBlock) BHash() common.Hash {
	panic("implement me")
}

//Removed the validator address, index. Mainly to ensure that the signature hash of the aggregate signature is consistent
type PrepareVote struct {
	Epoch       uint64            `json:"epoch"`
	ViewNumber  uint64            `json:"view_number"`
	BlockHash   *types.Block      `json:"block_hash"`
	BlockNumber uint64            `json:"block_number"`
	BlockIndex  uint32            `json:"block_index"` //The block number of the current ViewNumber proposal, 0....10
	ParentQC    ctypes.QuorumCert `json:"parent_qc"`
	Signature   ctypes.Signature  `json:"signature"`
}

func (PrepareVote) String() string {
	panic("implement me")
}

func (PrepareVote) MsgHash() common.Hash {
	panic("implement me")
}

func (PrepareVote) BHash() common.Hash {
	panic("implement me")
}

type ViewChange struct {
	Epoch      uint64            `json:"epoch"`
	ViewNumber uint64            `json:"view_number"`
	BlockHash  common.Hash       `json:"block_hash"`
	BlockNum   uint64            `json:"block_number"`
	PrepareQC  ctypes.QuorumCert `json:"prepare_qc"`
	Signature  ctypes.Signature  `json:"signature"`
}

func (ViewChange) String() string {
	panic("implement me")
}

func (ViewChange) MsgHash() common.Hash {
	panic("implement me")
}

func (ViewChange) BHash() common.Hash {
	panic("implement me")
}
