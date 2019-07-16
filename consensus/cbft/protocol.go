package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"math/big"
)

const CbftProtocolMaxMsgSize = 10 * 1024 * 1024 // Maximum cap on the size of a cbft protocol message

const (
	CBFTStatusMsg = 0x00 // Protocol messages belonging to cbft
)

type errCode int

const (
	ErrMsgTooLarge = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrCbftProtocolVersionMismatch
	ErrNoStatusMsg
	ErrForkedBlock
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

var errorToString = map[int]string{
	ErrMsgTooLarge:                 "Message too long",
	ErrDecode:                      "Invalid message",
	ErrInvalidMsgCode:              "Invalid message code",
	ErrCbftProtocolVersionMismatch: "CBFT Protocol version mismatch",
	ErrNoStatusMsg:                 "No status message",
	ErrForkedBlock:                 "Forked block",
}

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

type viewChange struct {
	Epoch      uint64      `json:"epoch"`
	ViewNumber uint64      `json:"view_number"`
	BlockHash  common.Hash `json:"block_hash"`
	BlockNum   uint64      `json:"block_number"`
	PrepareQC  quorumCert  `json:"prepare_qc"`
	Signature  Signature   `json:"signature"`
}

// Implement Message and including status information about peer.
type cbftStatusData struct {
	ProtocolVersion uint32      `json:"protocol_version"` // CBFT protocol version number.
	QCBn            *big.Int    `json:"qc_bn"`            // The highest local block number for collecting block signatures.
	QCBlock         common.Hash `json:"qc_block"`         // The highest local block hash for collecting block signatures.
	LockBn          *big.Int    `json:"lock_bn"`          // Locally locked block number.
	LockBlock       common.Hash `json:"lock_block"`       // Locally locked block hash.
	CmtBn           *big.Int    `json:"cmt_bn"`           // Locally submitted block number.
	CmtBlock        common.Hash `json:"cmt_block"`        // Locally submitted block hash.
}

func (s *cbftStatusData) String() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("[ProtocolVersion:%d, QCBn:%d, LockBn:%d, CmtBn:%d]", s.QCBn.Uint64(), s.LockBn.Uint64(), s.CmtBn.Uint64())
}

func (s *cbftStatusData) MsgHash() common.Hash {
	if s == nil {
		return common.Hash{}
	}
	return buildHash(CBFTStatusMsg, mergeBytes(s.QCBlock.Bytes(), s.LockBlock.Bytes(), s.CmtBlock.Bytes()))
}

func (s *cbftStatusData) BHash() common.Hash {
	if s == nil {
		return common.Hash{}
	}
	return s.QCBlock
}
