package protocols

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
)

// Maximum cap on the size of a cbft protocol message
const CbftProtocolMaxMsgSize = 10 * 1024 * 1024

const (
	CBFTStatusMsg        = 0x00 // Protocol messages belonging to cbft
	PrepareBlockMsg      = 0x01
	PrepareVoteMsg       = 0x02
	ViewChangeMsg        = 0x03
	GetPrepareBlockMsg   = 0x04
	GetQuorumCertMsg     = 0x05
	QuorumCertMsg        = 0x06
	GetQCPrepareBlockMsg = 0x07
	QCPrepareBlockMsg    = 0x08
	GetPrepareVoteMsg    = 0x09
	PrepareBlockHashMsg  = 0x0a
	PingMsg              = 0x0b
	PongMsg              = 0x0c
)

// A is used to convert specific message types according to the message body.
// The program is forcibly terminated if there is an unmatched message type and
// all types must exist in the match list.
func MessageType(msg interface{}) uint64 {
	// todo: need to process depending on mmessageType.
	switch msg.(type) {
	default:
		return PrepareBlockHashMsg
	}
	panic(fmt.Sprintf("unknown message type [%v]", reflect.TypeOf(msg)))
}

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

//
type GetPrepareBlock struct {
}

func (s *GetPrepareBlock) String() string {
	panic("implement me")
}

func (s *GetPrepareBlock) MsgHash() common.Hash {
	panic("implement me")
}

func (s *GetPrepareBlock) BHash() common.Hash {
	panic("implement me")
}

type GetQuorumCert struct {
}

func (s *GetQuorumCert) String() string {
	panic("implement me")
}

func (s *GetQuorumCert) MsgHash() common.Hash {
	panic("implement me")
}

func (s *GetQuorumCert) BHash() common.Hash {
	panic("implement me")
}

type QuorumCert struct {
}

func (s *QuorumCert) String() string {
	panic("implement me")
}

func (s *QuorumCert) MsgHash() common.Hash {
	panic("implement me")
}

func (s *QuorumCert) BHash() common.Hash {
	panic("implement me")
}

type GetQCPrepareBlock struct {
}

func (s *GetQCPrepareBlock) String() string {
	panic("implement me")
}

func (s *GetQCPrepareBlock) MsgHash() common.Hash {
	panic("implement me")
}

func (s *GetQCPrepareBlock) BHash() common.Hash {
	panic("implement me")
}

type QCPrepareBlock struct {
}

func (s *QCPrepareBlock) String() string {
	panic("implement me")
}

func (s *QCPrepareBlock) MsgHash() common.Hash {
	panic("implement me")
}

func (s *QCPrepareBlock) BHash() common.Hash {
	panic("implement me")
}

type GetPrepareVote struct {
}

func (s *GetPrepareVote) String() string {
	panic("implement me")
}

func (s *GetPrepareVote) MsgHash() common.Hash {
	panic("implement me")
}

func (s *GetPrepareVote) BHash() common.Hash {
	panic("implement me")
}

type PrepareBlockHash struct {
}

func (s *PrepareBlockHash) String() string {
	panic("implement me")
}

func (s *PrepareBlockHash) MsgHash() common.Hash {
	panic("implement me")
}

func (s *PrepareBlockHash) BHash() common.Hash {
	panic("implement me")
}

type Ping [1]string

func (s *Ping) String() string {
	panic("implement me")
}

func (s *Ping) MsgHash() common.Hash {
	panic("implement me")
}

func (s *Ping) BHash() common.Hash {
	panic("implement me")
}

type Pong [1]string

func (s *Pong) String() string {
	panic("implement me")
}

func (s *Pong) MsgHash() common.Hash {
	panic("implement me")
}

func (s *Pong) BHash() common.Hash {
	panic("implement me")
}
