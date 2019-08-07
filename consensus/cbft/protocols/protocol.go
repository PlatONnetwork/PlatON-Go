package protocols

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// Maximum cap on the size of a cbft protocol message
const CbftProtocolMaxMsgSize = 10 * 1024 * 1024

// Default average delay time (in milliseconds).
const DEFAULT_AVG_LATENCY = 100

const (
	CBFTStatusMsg           = 0x00 // Protocol messages belonging to cbft
	PrepareBlockMsg         = 0x01
	PrepareVoteMsg          = 0x02
	ViewChangeMsg           = 0x03
	GetPrepareBlockMsg      = 0x04
	GetBlockQuorumCertMsg   = 0x05
	BlockQuorumCertMsg      = 0x06
	GetPrepareVoteMsg       = 0x07
	PrepareVotesMsg         = 0x08
	GetQCBlockListMsg       = 0x09
	QCBlockListMsg          = 0x0a
	GetLatestStatusMsg      = 0x0b
	LatestStatusMsg         = 0x0c
	PrepareBlockHashMsg     = 0x0d
	GetViewChangeMsg        = 0x0e
	PingMsg                 = 0x0f
	PongMsg                 = 0x10
	ViewChangeQuorumCertMsg = 0x11
)

// A is used to convert specific message types according to the message body.
// The program is forcibly terminated if there is an unmatched message type and
// all types must exist in the match list.
func MessageType(msg interface{}) uint64 {
	// todo: need to process depending on mmessageType.
	switch msg.(type) {
	case *CbftStatusData:
		return CBFTStatusMsg
	case *PrepareBlock:
		return PrepareBlockMsg
	case *PrepareVote:
		return PrepareVoteMsg
	case *ViewChange:
		return ViewChangeMsg
	case *GetPrepareBlock:
		return GetPrepareBlockMsg
	case *GetBlockQuorumCert:
		return GetBlockQuorumCertMsg
	case *BlockQuorumCert:
		return BlockQuorumCertMsg
	case *GetQCBlockList:
		return GetQCBlockListMsg
	case *GetPrepareVote:
		return GetPrepareVoteMsg
	case *PrepareBlockHash:
		return PrepareBlockHashMsg
	case *PrepareVotes:
		return PrepareVotesMsg
	case *QCBlockList:
		return QCBlockListMsg
	case *GetLatestStatus:
		return GetLatestStatusMsg
	case *LatestStatus:
		return LatestStatusMsg
	case *GetViewChange:
		return GetViewChangeMsg
	case *Ping:
		return PingMsg
	case *Pong:
		return PongMsg
	case *ViewChangeQuorumCert:
		return ViewChangeQuorumCertMsg
	default:
	}
	panic(fmt.Sprintf("unknown message type [%v}", reflect.TypeOf(msg)))
}

// Proposed block carrier.
type PrepareBlock struct {
	Epoch         uint64               `json:"epoch"`
	ViewNumber    uint64               `json:"view_number"`
	Block         *types.Block         `json:"block_hash"`
	BlockIndex    uint32               `json:"block_index"` // The block number of the current ViewNumber proposal, 0....10
	ProposalIndex uint32               `json:"proposal_index"`
	PrepareQC     *ctypes.QuorumCert   `json:"prepare_qc"rlp:"nil"`    // N-f aggregate signature
	ViewChangeQC  *ctypes.ViewChangeQC `json:"viewchange_qc"rlp:"nil"` // viewChange aggregate signature
	Signature     ctypes.Signature     `json:"signature"`
}

func (pb *PrepareBlock) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d,BlockIndex:%d}",
		pb.Epoch, pb.ViewNumber, pb.Block.Hash().TerminalString(), pb.Block.NumberU64(), pb.BlockIndex)
}

func (pb *PrepareBlock) MsgHash() common.Hash {
	return utils.BuildHash(PrepareBlockMsg,
		utils.MergeBytes(common.Uint64ToBytes(pb.ViewNumber), pb.Block.Hash().Bytes(), pb.Signature.Bytes()))
}

func (pb *PrepareBlock) BHash() common.Hash {
	return pb.Block.Hash()
}

func (pb *PrepareBlock) BlockNum() uint64 {
	return pb.Block.NumberU64()
}

func (pb *PrepareBlock) NodeIndex() uint32 {
	return pb.ProposalIndex
}

func (pb *PrepareBlock) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		pb.Epoch,
		pb.ViewNumber,
		pb.Block.Hash(),
		pb.BlockIndex,
		pb.ProposalIndex,
	})
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (pb *PrepareBlock) Sign() []byte {
	return pb.Signature.Bytes()
}

func (pb *PrepareBlock) SetSign(sign []byte) {
	pb.Signature.SetBytes(sign)
}

// Removed the validator address, index. Mainly to ensure that the signature hash of the aggregate signature is consistent
type PrepareVote struct {
	Epoch          uint64             `json:"epoch"`
	ViewNumber     uint64             `json:"view_number"`
	BlockHash      common.Hash        `json:"block_hash"`
	BlockNumber    uint64             `json:"block_number"`
	BlockIndex     uint32             `json:"block_index"` // The block number of the current ViewNumber proposal, 0....10
	ValidatorIndex uint32             `json:"validator_index"`
	ParentQC       *ctypes.QuorumCert `json:"parent_qc"rlp:"nil"`
	Signature      ctypes.Signature   `json:"signature"`
}

func (pv *PrepareVote) String() string {
	return fmt.Sprintf("{Epoch:%d,VN:%d,BlockHash:%s,BlockNumber:%d,BlockIndex:%d}",
		pv.Epoch, pv.ViewNumber, pv.BlockHash.TerminalString(), pv.BlockNumber, pv.BlockIndex)
}

func (pv *PrepareVote) MsgHash() common.Hash {
	return utils.BuildHash(PrepareVoteMsg,
		utils.MergeBytes(common.Uint64ToBytes(pv.ViewNumber), pv.BlockHash.Bytes(), common.Uint32ToBytes(pv.BlockIndex), pv.Signature.Bytes()))
}

func (pv *PrepareVote) BHash() common.Hash {
	return pv.BlockHash
}

func (pv *PrepareVote) BlockNum() uint64 {
	return pv.BlockNumber
}

func (pv *PrepareVote) NodeIndex() uint32 {
	return pv.ValidatorIndex
}

func (pv *PrepareVote) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		pv.Epoch,
		pv.ViewNumber,
		pv.BlockHash,
		pv.BlockNumber,
		pv.BlockIndex,
	})

	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (pv *PrepareVote) Sign() []byte {
	return pv.Signature.Bytes()
}

func (pv *PrepareVote) SetSign(sign []byte) {
	pv.Signature.SetBytes(sign)
}

// Message structure for view switching.
type ViewChange struct {
	Epoch          uint64             `json:"epoch"`
	ViewNumber     uint64             `json:"view_number"`
	BlockHash      common.Hash        `json:"block_hash"`
	BlockNumber    uint64             `json:"block_number"`
	ValidatorIndex uint32             `json:"validator_index"`
	PrepareQC      *ctypes.QuorumCert `json:"prepare_qc"rlp:"nil"`
	Signature      ctypes.Signature   `json:"signature"`
}

func (vc *ViewChange) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockHash:%vc,BlockNumber:%d}",
		vc.Epoch, vc.ViewNumber, vc.BlockHash.TerminalString(), vc.BlockNumber)
}

func (vc *ViewChange) MsgHash() common.Hash {
	return utils.BuildHash(ViewChangeMsg, utils.MergeBytes(common.Uint64ToBytes(vc.ViewNumber),
		vc.BlockHash.Bytes(), common.Uint64ToBytes(vc.BlockNumber)))
}

func (vc *ViewChange) BHash() common.Hash {
	return vc.BlockHash
}

func (vc *ViewChange) BlockNum() uint64 {
	return vc.BlockNumber
}

func (vc *ViewChange) NodeIndex() uint32 {
	return vc.ValidatorIndex
}

func (vc *ViewChange) CannibalizeBytes() ([]byte, error) {
	blockEpoch, blockView := uint64(0), uint64(0)
	if vc.PrepareQC != nil {
		blockEpoch, blockView = vc.PrepareQC.Epoch, vc.ViewNumber
	}
	buf, err := rlp.EncodeToBytes([]interface{}{
		vc.Epoch,
		vc.ViewNumber,
		vc.BlockHash,
		vc.BlockNumber,
		blockEpoch,
		blockView,
	})

	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (vc *ViewChange) Sign() []byte {
	return vc.Signature.Bytes()
}

func (vc *ViewChange) SetSign(sign []byte) {
	vc.Signature.SetBytes(sign)
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
	return fmt.Sprintf("{ProtocolVersion:%d,QCBn:%d,LockBn:%d,CmtBn:%d}",
		s.ProtocolVersion, s.QCBn.Uint64(), s.LockBn.Uint64(), s.CmtBn.Uint64())
}

func (s *CbftStatusData) MsgHash() common.Hash {
	return utils.BuildHash(CBFTStatusMsg, utils.MergeBytes(s.QCBlock.Bytes(),
		s.LockBlock.Bytes(), s.CmtBlock.Bytes()))
}

func (s *CbftStatusData) BHash() common.Hash {
	return s.QCBlock
}

// CBFT protocol message - used to get the
// proposed block information.
type GetPrepareBlock struct {
	Epoch      uint64
	ViewNumber uint64
	BlockIndex uint32
}

func (s *GetPrepareBlock) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockIndex:%d}", s.Epoch, s.ViewNumber, s.BlockIndex)
}

func (s *GetPrepareBlock) MsgHash() common.Hash {
	return utils.BuildHash(GetPrepareBlockMsg, utils.MergeBytes(common.Uint64ToBytes(s.ViewNumber), common.Uint32ToBytes(s.BlockIndex)))
}

func (s *GetPrepareBlock) BHash() common.Hash {
	return common.Hash{}
}

// Protocol message for obtaining an aggregated signature.
// todo: Need to determine the attribute field - ParentQC.
type GetBlockQuorumCert struct {
	BlockHash   common.Hash `json:"block_hash"`   // The hash of the block to be acquired.
	BlockNumber uint64      `json:"block_number"` // The number of the block to be acquired.
}

func (s *GetBlockQuorumCert) String() string {
	return fmt.Sprintf("{Hash:%s,Number:%d}", s.BlockHash.TerminalString(), s.BlockNumber)
}

func (s *GetBlockQuorumCert) MsgHash() common.Hash {
	return utils.BuildHash(GetBlockQuorumCertMsg, utils.MergeBytes(s.BlockHash.Bytes(), common.Uint64ToBytes(s.BlockNumber)))
}

func (s *GetBlockQuorumCert) BHash() common.Hash {
	return s.BlockHash
}

// Aggregate signature response message, representing
// aggregated signature information for a block.
type BlockQuorumCert struct {
	BlockQC *ctypes.QuorumCert `json:"qc"` // Block aggregation signature information
}

func (s *BlockQuorumCert) String() string {
	return fmt.Sprintf("{ViewNumber:%d,Hash:%s,Number:%d}",
		s.BlockQC.ViewNumber, s.BlockQC.BlockHash.TerminalString(), s.BlockQC.BlockNumber)
}

func (s *BlockQuorumCert) MsgHash() common.Hash {
	return utils.BuildHash(BlockQuorumCertMsg, utils.MergeBytes(
		s.BlockQC.BlockHash.Bytes(),
		common.Uint64ToBytes(s.BlockQC.BlockNumber), s.BlockQC.Signature.Bytes()))
}

func (s *BlockQuorumCert) BHash() common.Hash {
	return s.BlockQC.BlockHash
}

// Used to get block information that has reached QC.
// Note: Get up to 3 blocks of data at a time.
type GetQCBlockList struct {
	BlockHash   common.Hash `json:"block_hash"`
	BlockNumber uint64      `json:"block_number"` // The number corresponding to the block.
}

func (s *GetQCBlockList) String() string {
	return fmt.Sprintf("{Number:%d}", s.BlockNumber)
}

func (s *GetQCBlockList) MsgHash() common.Hash {
	return utils.BuildHash(GetQCBlockListMsg, utils.MergeBytes(
		common.Uint64ToBytes(s.BlockNumber)))
}

func (s *GetQCBlockList) BHash() common.Hash {
	return common.Hash{}
}

// Message used to get block voting.
type GetPrepareVote struct {
	ViewNumber  uint32
	BlockHash   common.Hash
	BlockNumber uint64
	BlockIndex  uint32
	VoteBits    *utils.BitArray
}

func (s *GetPrepareVote) String() string {
	return fmt.Sprintf("{Hash:%s,Number:%d,ViewNumber:%d}", s.BlockHash.TerminalString(), s.BlockNumber, s.ViewNumber)
}

func (s *GetPrepareVote) MsgHash() common.Hash {
	return utils.BuildHash(GetPrepareVoteMsg, utils.MergeBytes(
		s.BlockHash.Bytes(), common.Uint64ToBytes(s.BlockNumber),
		s.VoteBits.Bytes()))
}

func (s *GetPrepareVote) BHash() common.Hash {
	return s.BlockHash
}

// Message used to respond to the number of block votes.
type PrepareVotes struct {
	BlockHash   common.Hash
	BlockNumber uint64
	Votes       []*PrepareVote // Block voting set.
}

func (s *PrepareVotes) String() string {
	return fmt.Sprintf("{Hash:%s,Number:%d,Votes:%d}", s.BlockHash.TerminalString(), s.BlockNumber, len(s.Votes))
}

func (s *PrepareVotes) MsgHash() common.Hash {
	return utils.BuildHash(PrepareVotesMsg, utils.MergeBytes(s.BlockHash.Bytes(), common.Uint64ToBytes(s.BlockNumber)))
}

func (s *PrepareVotes) BHash() common.Hash {
	return s.BlockHash
}

// Represents the hash of the proposed block for secondary propagation.
type PrepareBlockHash struct {
	Epoch       uint64
	ViewNumber  uint64
	BlockIndex  uint32
	BlockHash   common.Hash
	BlockNumber uint64
}

func (s *PrepareBlockHash) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockIndex:%d,Hash:%s,Number:%d}", s.Epoch, s.ViewNumber, s.BlockIndex, s.BlockHash.TerminalString(), s.BlockNumber)
}

func (s *PrepareBlockHash) MsgHash() common.Hash {
	return utils.BuildHash(PrepareBlockHashMsg, utils.MergeBytes(s.BlockHash.Bytes(), common.Uint64ToBytes(s.BlockNumber)))
}

func (s *PrepareBlockHash) BHash() common.Hash {
	return s.BlockHash
}

// For time detection.
type Ping [1]string

func (s *Ping) String() string {
	return fmt.Sprintf("{pingTime:%s}", s[0])
}

func (s *Ping) MsgHash() common.Hash {
	return utils.BuildHash(PingMsg, utils.MergeBytes([]byte(s[0])))
}

func (s *Ping) BHash() common.Hash {
	return common.Hash{}
}

// Response to ping.
type Pong [1]string

func (s *Pong) String() string {
	return fmt.Sprintf("{pongTime:%s}", s[0])
}

func (s *Pong) MsgHash() common.Hash {
	return utils.BuildHash(PongMsg, utils.MergeBytes([]byte(s[0])))
}

func (s *Pong) BHash() common.Hash {
	return common.Hash{}
}

// CBFT synchronize blocks that have reached qc.
type QCBlockList struct {
	QC     []*ctypes.QuorumCert
	Blocks []*types.Block
}

func (s *QCBlockList) String() string {
	return fmt.Sprintf("{QC.Len:%d,Blocks.Len:%d}", len(s.QC), len(s.Blocks))
}

func (s *QCBlockList) MsgHash() common.Hash {
	if len(s.QC) != 0 {
		return utils.BuildHash(QCBlockListMsg, utils.MergeBytes(s.QC[0].BlockHash.Bytes(),
			s.QC[0].Signature.Bytes()))
	}
	if len(s.Blocks) != 0 {
		return utils.BuildHash(QCBlockListMsg, utils.MergeBytes(s.Blocks[0].Hash().Bytes(),
			s.Blocks[0].Number().Bytes()))
	}
	return common.Hash{}
}

func (s *QCBlockList) BHash() common.Hash {
	// No explicit hash value and return empty hash.
	return common.Hash{}
}

// State synchronization for nodes.
type GetLatestStatus struct {
	BlockNumber uint64 // Block height sent by the requester
	LogicType   uint64 // LogicType: 1 QCBn, 2 LockedBn, 3 CommitBn
}

func (s *GetLatestStatus) String() string {
	return fmt.Sprintf("{BlockNumber:%d,LogicType:%d}", s.BlockNumber, s.LogicType)
}

func (s *GetLatestStatus) MsgHash() common.Hash {
	return utils.BuildHash(GetLatestStatusMsg,
		utils.MergeBytes(common.Uint64ToBytes(s.BlockNumber), common.Uint64ToBytes(s.LogicType)))
}

func (s *GetLatestStatus) BHash() common.Hash {
	return common.Hash{}
}

// Response message to GetLatestStatus request.
type LatestStatus struct {
	BlockNumber uint64 // Block height sent by responder.
	LogicType   uint64 // LogicType: 1 QCBn, 2 LockedBn, 3 CommitBn
}

func (s *LatestStatus) String() string {
	return fmt.Sprintf("{BlockNumber:%d,LogicType:%d}", s.BlockNumber, s.LogicType)
}

func (s *LatestStatus) MsgHash() common.Hash {
	return utils.BuildHash(LatestStatusMsg,
		utils.MergeBytes(common.Uint64ToBytes(s.BlockNumber), common.Uint64ToBytes(s.LogicType)))
}

func (s *LatestStatus) BHash() common.Hash {
	return common.Hash{}
}

// Used to actively request to get viewChange.
type GetViewChange struct {
	Epoch       uint64   `json:"epoch"`
	ViewNumber  uint64   `json:"view_number"`
	NodeIndexes []uint32 `json:"node_indexes"`
}

func (s *GetViewChange) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,NodeIndexesLen:%d}", s.Epoch, s.ViewNumber, len(s.NodeIndexes))
}

func (s *GetViewChange) MsgHash() common.Hash {
	return utils.BuildHash(GetViewChangeMsg,
		utils.MergeBytes(common.Uint64ToBytes(s.Epoch), common.Uint64ToBytes(s.ViewNumber)))
}

func (s *GetViewChange) BHash() common.Hash {
	return common.Hash{}
}

type ViewChangeQuorumCert struct {
	ViewChangeQC *ctypes.ViewChangeQC `json:"viewchange_qc"` // viewChange aggregate signature
}

func (v *ViewChangeQuorumCert) String() string {
	epoch, viewNumber, hash, number := v.ViewChangeQC.MaxBlock()
	return fmt.Sprintf("{Epoch:%d,VN:%d,Hash:%s,Number:%d}",
		epoch, viewNumber, hash.TerminalString(), number)
}

func (v *ViewChangeQuorumCert) MsgHash() common.Hash {
	epoch, viewNumber, hash, number := v.ViewChangeQC.MaxBlock()
	return utils.BuildHash(ViewChangeQuorumCertMsg, utils.MergeBytes(
		common.Uint64ToBytes(epoch),
		common.Uint64ToBytes(viewNumber),
		hash.Bytes(),
		common.Uint64ToBytes(number)))
}

func (v *ViewChangeQuorumCert) BHash() common.Hash {
	_, _, hash, _ := v.ViewChangeQC.MaxBlock()
	return hash
}
