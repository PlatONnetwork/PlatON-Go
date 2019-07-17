package protocols

import (
	"fmt"
	"math/big"

	"github.com/PlatONnetwork/PlatON-Go/common"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

// Maximum cap on the size of a cbft protocol message
const CbftProtocolMaxMsgSize = 10 * 1024 * 1024

const (
	CBFTStatusMsg = 0x00 // Protocol messages belonging to cbft
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

// Removed the validator address, index. Mainly to ensure that the signature hash of the aggregate signature is consistent
type PrepareVote struct {
	Epoch       uint64             `json:"epoch"`
	ViewNumber  uint64             `json:"view_number"`
	BlockHash   common.Hash        `json:"block_hash"`
	BlockNumber uint64             `json:"block_number"`
	BlockIndex  uint32             `json:"block_index"` //The block number of the current ViewNumber proposal, 0....10
	ParentQC    *ctypes.QuorumCert `json:"parent_qc"`
	Signature   ctypes.Signature   `json:"signature"`
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
	Epoch      uint64             `json:"epoch"`
	ViewNumber uint64             `json:"view_number"`
	BlockHash  common.Hash        `json:"block_hash"`
	BlockNum   uint64             `json:"block_number"`
	PrepareQC  *ctypes.QuorumCert `json:"prepare_qc"`
	Signature  ctypes.Signature   `json:"signature"`
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

// cbftStatusData implement Message and including status information about peer.
type CbftStatusData struct {
	ProtocolVersion uint32      `json:"protocol_version"` // CBFT protocol version number.
	QCBn            *big.Int    `json:"qc_bn"`            // The highest local block number for collecting block signatures.
	QCBlock         common.Hash `json:"qc_block"`         // The highest local block hash for collecting block signatures.
	LockBn          *big.Int    `json:"lock_bn"`          // Locally locked block number.
	LockBlock       common.Hash `json:"lock_block"`       // Locally locked block hash.
	CmtBn           *big.Int    `json:"cmt_bn"`           // Locally submitted block number.
	CmtBlock        common.Hash `json:"cmt_block"`        // Locally submitted block hash.
}

func (s *CbftStatusData) String() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("[ProtocolVersion:%d, QCBn:%d, LockBn:%d, CmtBn:%d]", s.QCBn.Uint64(), s.LockBn.Uint64(), s.CmtBn.Uint64())
}

func (s *CbftStatusData) MsgHash() common.Hash {
	if s == nil {
		return common.Hash{}
	}
	return utils.BuildHash(CBFTStatusMsg, utils.MergeBytes(s.QCBlock.Bytes(), s.LockBlock.Bytes(), s.CmtBlock.Bytes()))
}

func (s *CbftStatusData) BHash() common.Hash {
	if s == nil {
		return common.Hash{}
	}
	return s.QCBlock
}
