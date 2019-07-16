package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

type ConsensusMsg interface {
	CannibalizeBytes() ([]byte, error)
	Sign() []byte
}

type Message interface {
	String() string
	MsgHash() common.Hash
	BHash() common.Hash
}

type MsgInfo struct {
	Msg    Message
	PeerID discover.NodeID
}

type prepareBlock struct {
	Epoch         uint64         `json:"epoch"`
	ViewNumber    uint64         `json:"view_number"`
	Block         *types.Block   `json:"block_hash"`
	BlockIndex    uint32         `json:"block_index"` //The block number of the current ViewNumber proposal, 0....10
	ProposalIndex uint32         `json:"proposal_index"`
	ProposalAddr  common.Address `json:"proposal_address"`
	PrepareQC     *quorumCert    `json:"prepare_qc"`    //N-f aggregate signature
	ViewChangeQC  []*quorumCert  `json:"viewchange_qc"` //viewchange aggregate signature
	Signature     Signature      `json:"signature"`
}

func (prepareBlock) String() string {
	panic("implement me")
}

func (prepareBlock) MsgHash() common.Hash {
	panic("implement me")
}

func (prepareBlock) BHash() common.Hash {
	panic("implement me")
}

//Removed the validator address, index. Mainly to ensure that the signature hash of the aggregate signature is consistent
type prepareVote struct {
	Epoch       uint64       `json:"epoch"`
	ViewNumber  uint64       `json:"view_number"`
	BlockHash   *types.Block `json:"block_hash"`
	BlockNumber uint64       `json:"block_number"`
	BlockIndex  uint32       `json:"block_index"` //The block number of the current ViewNumber proposal, 0....10
	ParentQC    quorumCert   `json:"parent_qc"`
	Signature   Signature    `json:"signature"`
}

func (prepareVote) String() string {
	panic("implement me")
}

func (prepareVote) MsgHash() common.Hash {
	panic("implement me")
}

func (prepareVote) BHash() common.Hash {
	panic("implement me")
}

type viewChange struct {
	Epoch      uint64      `json:"epoch"`
	ViewNumber uint64      `json:"view_number"`
	BlockHash  common.Hash `json:"block_hash"`
	BlockNum   uint64      `json:"block_number"`
	PrepareQC  quorumCert  `json:"prepare_qc"`
	Signature  Signature   `json:"signature"`
}

func (viewChange) String() string {
	panic("implement me")
}

func (viewChange) MsgHash() common.Hash {
	panic("implement me")
}

func (viewChange) BHash() common.Hash {
	panic("implement me")
}
