// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package protocols

import (
	"fmt"
	"math/big"
	"reflect"
	"sync/atomic"

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
const DefaultAvgLatency = 100

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
	ViewChangesMsg          = 0x12
)

// A is used to convert specific message types according to the message body.
// The program is forcibly terminated if there is an unmatched message type and
// all types must exist in the match list.
func MessageType(msg interface{}) uint64 {
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
	case *ViewChanges:
		return ViewChangesMsg
	default:
	}
	panic(fmt.Sprintf("unknown message type [%v}", reflect.TypeOf(msg)))
}

// Proposed block carrier.
type PrepareBlock struct {
	Epoch         uint64               `json:"epoch"`
	ViewNumber    uint64               `json:"viewNumber"`
	Block         *types.Block         `json:"blockHash"`
	BlockIndex    uint32               `json:"blockIndex"`             // The block number of the current ViewNumber proposal, 0....10
	ProposalIndex uint32               `json:"proposalIndex"`          // Proposer's index.
	PrepareQC     *ctypes.QuorumCert   `json:"prepareQC" rlp:"nil"`    // N-f aggregate signature
	ViewChangeQC  *ctypes.ViewChangeQC `json:"viewchangeQC" rlp:"nil"` // viewChange aggregate signature
	Signature     ctypes.Signature     `json:"signature"`              // PrepareBlock signature information
	messageHash   atomic.Value         `rlp:"-"`
}

func (pb *PrepareBlock) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d,BlockIndex:%d,ProposalIndex:%d,ParentHash:%s}",
		pb.Epoch, pb.ViewNumber, pb.Block.Hash().TerminalString(), pb.Block.NumberU64(), pb.BlockIndex, pb.ProposalIndex, pb.Block.ParentHash().TerminalString())
}

func (pb *PrepareBlock) MsgHash() common.Hash {
	if mhash := pb.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(PrepareBlockMsg,
		utils.MergeBytes(common.Uint64ToBytes(pb.ViewNumber), pb.Block.Hash().Bytes(), pb.Signature.Bytes()))
	pb.messageHash.Store(v)
	return v
}

func (pb *PrepareBlock) BHash() common.Hash {
	return pb.Block.Hash()
}
func (pb *PrepareBlock) EpochNum() uint64 {
	return pb.Epoch
}
func (pb *PrepareBlock) ViewNum() uint64 {
	return pb.ViewNumber
}
func (pb *PrepareBlock) BlockNum() uint64 {
	return pb.Block.NumberU64()
}

func (pb *PrepareBlock) NodeIndex() uint32 {
	return pb.ProposalIndex
}

func (pb *PrepareBlock) CannibalizeBytes() ([]byte, error) {
	blockData, err := rlp.EncodeToBytes(pb.Block)
	if err != nil {
		return nil, err
	}
	buf, err := rlp.EncodeToBytes([]interface{}{
		pb.Epoch,
		pb.ViewNumber,
		pb.Block.Hash(),
		pb.Block.NumberU64(),
		crypto.Keccak256(blockData),
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

// Removed the validator address, index. Mainly to ensure
// that the signature hash of the aggregate signature is consistent.
type PrepareVote struct {
	Epoch          uint64             `json:"epoch"`
	ViewNumber     uint64             `json:"viewNumber"`
	BlockHash      common.Hash        `json:"blockHash"`
	BlockNumber    uint64             `json:"blockNumber"`
	BlockIndex     uint32             `json:"blockIndex"` // The block number of the current ViewNumber proposal, 0....10
	ValidatorIndex uint32             `json:"validatorIndex"`
	ParentQC       *ctypes.QuorumCert `json:"parentQC" rlp:"nil"`
	Signature      ctypes.Signature   `json:"signature"`
	messageHash    atomic.Value       `rlp:"-"`
}

func (pv *PrepareVote) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d,BlockIndex:%d,ValidatorIndex:%d}",
		pv.Epoch, pv.ViewNumber, pv.BlockHash.TerminalString(), pv.BlockNumber, pv.BlockIndex, pv.ValidatorIndex)
}

func (pv *PrepareVote) MsgHash() common.Hash {
	if mhash := pv.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(PrepareVoteMsg,
		utils.MergeBytes(common.Uint64ToBytes(pv.ViewNumber), pv.BlockHash.Bytes(), common.Uint32ToBytes(pv.BlockIndex), pv.Signature.Bytes()))
	pv.messageHash.Store(v)
	return v
}

func (pv *PrepareVote) BHash() common.Hash {
	return pv.BlockHash
}
func (pv *PrepareVote) EpochNum() uint64 {
	return pv.Epoch
}
func (pv *PrepareVote) ViewNum() uint64 {
	return pv.ViewNumber
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

func (pv *PrepareVote) EqualState(vote *PrepareVote) bool {
	return pv.Epoch == vote.Epoch &&
		pv.ViewNumber == vote.ViewNumber &&
		pv.BlockHash == vote.BlockHash &&
		pv.BlockNumber == vote.BlockNumber &&
		pv.BlockIndex == vote.BlockIndex
}

// ViewChange is message structure for view switching.
type ViewChange struct {
	Epoch          uint64             `json:"epoch"`
	ViewNumber     uint64             `json:"viewNumber"`
	BlockHash      common.Hash        `json:"blockHash"`
	BlockNumber    uint64             `json:"blockNumber"`
	ValidatorIndex uint32             `json:"validatorIndex"`
	PrepareQC      *ctypes.QuorumCert `json:"prepareQC" rlp:"nil"`
	Signature      ctypes.Signature   `json:"signature"`
	messageHash    atomic.Value       `rlp:"-"`
}

func (vc *ViewChange) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockHash:%s,BlockNumber:%d,ValidatorIndex:%d}",
		vc.Epoch, vc.ViewNumber, vc.BlockHash.TerminalString(), vc.BlockNumber, vc.ValidatorIndex)
}

func (vc *ViewChange) MsgHash() common.Hash {
	if mhash := vc.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(ViewChangeMsg, utils.MergeBytes(common.Uint64ToBytes(vc.ViewNumber),
		vc.BlockHash.Bytes(), common.Uint64ToBytes(vc.BlockNumber), common.Uint32ToBytes(vc.ValidatorIndex),
		vc.Signature.Bytes()))
	vc.messageHash.Store(v)
	return v
}

func (vc *ViewChange) BHash() common.Hash {
	return vc.BlockHash
}

func (vc *ViewChange) EpochNum() uint64 {
	return vc.Epoch
}

func (vc *ViewChange) ViewNum() uint64 {
	return vc.ViewNumber
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
		blockEpoch, blockView = vc.PrepareQC.Epoch, vc.PrepareQC.ViewNumber
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

type ViewChanges struct {
	VCs         []*ViewChange
	messageHash atomic.Value `rlp:"-"`
}

func (v ViewChanges) String() string {
	if len(v.VCs) != 0 {
		epoch, viewNumber := v.VCs[0].Epoch, v.VCs[0].ViewNumber
		return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Len:%d}", epoch, viewNumber, len(v.VCs))
	}
	return fmt.Sprintf("{Len:%d}", len(v.VCs))
}

func (v ViewChanges) MsgHash() common.Hash {
	if mhash := v.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	var mv common.Hash
	if len(v.VCs) != 0 {
		epoch, viewNumber := v.VCs[0].Epoch, v.VCs[0].ViewNumber
		mv = utils.BuildHash(ViewChangesMsg, utils.MergeBytes(common.Uint64ToBytes(epoch),
			common.Uint64ToBytes(viewNumber)))
	} else {
		mv = utils.BuildHash(ViewChangesMsg, common.Hash{}.Bytes())
	}
	v.messageHash.Store(mv)
	return mv
}

func (ViewChanges) BHash() common.Hash {
	return common.Hash{}
}

// CbftStatusData implement Message and including status information about peer.
type CbftStatusData struct {
	ProtocolVersion uint32       `json:"protocolVersion"` // CBFT protocol version number.
	QCBn            *big.Int     `json:"qcBn"`            // The highest local block number for collecting block signatures.
	QCBlock         common.Hash  `json:"qcBlock"`         // The highest local block hash for collecting block signatures.
	LockBn          *big.Int     `json:"lockBn"`          // Locally locked block number.
	LockBlock       common.Hash  `json:"lockBlock"`       // Locally locked block hash.
	CmtBn           *big.Int     `json:"cmtBn"`           // Locally submitted block number.
	CmtBlock        common.Hash  `json:"cmtBlock"`        // Locally submitted block hash.
	messageHash     atomic.Value `rlp:"-"`
}

func (s *CbftStatusData) String() string {
	return fmt.Sprintf("{ProtocolVersion:%d,QCBn:%d,LockBn:%d,CmtBn:%d}",
		s.ProtocolVersion, s.QCBn.Uint64(), s.LockBn.Uint64(), s.CmtBn.Uint64())
}

func (s *CbftStatusData) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(CBFTStatusMsg, utils.MergeBytes(s.QCBlock.Bytes(),
		s.LockBlock.Bytes(), s.CmtBlock.Bytes()))
	s.messageHash.Store(v)
	return v
}

func (s *CbftStatusData) BHash() common.Hash {
	return s.QCBlock
}

// GetPrepareBlock is used to get the
// proposed block information.
type GetPrepareBlock struct {
	Epoch       uint64
	ViewNumber  uint64
	BlockIndex  uint32
	messageHash atomic.Value `rlp:"-"`
}

func (s *GetPrepareBlock) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockIndex:%d}", s.Epoch, s.ViewNumber, s.BlockIndex)
}

func (s *GetPrepareBlock) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(GetPrepareBlockMsg, utils.MergeBytes(common.Uint64ToBytes(s.Epoch), common.Uint64ToBytes(s.ViewNumber), common.Uint32ToBytes(s.BlockIndex)))
	s.messageHash.Store(v)
	return v
}

func (s *GetPrepareBlock) BHash() common.Hash {
	return common.Hash{}
}

// GetBlockQuorumCert is the protocol message for obtaining an aggregated signature.
// todo: Need to determine the attribute field - ParentQC.
type GetBlockQuorumCert struct {
	BlockHash   common.Hash  `json:"blockHash"`   // The hash of the block to be acquired.
	BlockNumber uint64       `json:"blockNumber"` // The number of the block to be acquired.
	messageHash atomic.Value `json:"-" rlp:"-"`
}

func (s *GetBlockQuorumCert) String() string {
	return fmt.Sprintf("{Hash:%s,Number:%d}", s.BlockHash.TerminalString(), s.BlockNumber)
}

func (s *GetBlockQuorumCert) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(GetBlockQuorumCertMsg, utils.MergeBytes(s.BlockHash.Bytes(), common.Uint64ToBytes(s.BlockNumber)))
	s.messageHash.Store(v)
	return v
}

func (s *GetBlockQuorumCert) BHash() common.Hash {
	return s.BlockHash
}

// Aggregate signature response message, representing
// aggregated signature information for a block.
type BlockQuorumCert struct {
	BlockQC     *ctypes.QuorumCert `json:"qc"`        // Block aggregation signature information.
	messageHash atomic.Value       `json:"-" rlp:"-"` // BlockQuorumCert hash value.
}

func (s *BlockQuorumCert) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockIndex:%d,Hash:%s,Number:%d}",
		s.BlockQC.Epoch, s.BlockQC.ViewNumber, s.BlockQC.BlockIndex, s.BlockQC.BlockHash.TerminalString(), s.BlockQC.BlockNumber)
}

func (s *BlockQuorumCert) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(BlockQuorumCertMsg, utils.MergeBytes(
		s.BlockQC.BlockHash.Bytes(),
		common.Uint64ToBytes(s.BlockQC.BlockNumber), s.BlockQC.Signature.Bytes()))
	s.messageHash.Store(v)
	return v
}

func (s *BlockQuorumCert) BHash() common.Hash {
	return s.BlockQC.BlockHash
}

// Used to get block information that has reached QC.
// Note: Get up to 3 blocks of data at a time.
type GetQCBlockList struct {
	BlockHash   common.Hash  `json:"blockHash"`   // The hash to the block.
	BlockNumber uint64       `json:"blockNumber"` // The number corresponding to the block.
	messageHash atomic.Value `json:"-" rlp:"-"`
}

func (s *GetQCBlockList) String() string {
	return fmt.Sprintf("{BlockkNumber:%d,BlockHash:%s}", s.BlockNumber, s.BlockHash.TerminalString())
}

func (s *GetQCBlockList) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(GetQCBlockListMsg, utils.MergeBytes(
		common.Uint64ToBytes(s.BlockNumber)))
	s.messageHash.Store(v)
	return v
}

func (s *GetQCBlockList) BHash() common.Hash {
	return common.Hash{}
}

// Message used to get block voting.
type GetPrepareVote struct {
	Epoch       uint64
	ViewNumber  uint64
	BlockIndex  uint32
	UnKnownSet  *utils.BitArray
	messageHash atomic.Value `json:"-" rlp:"-"`
}

func (s *GetPrepareVote) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockIndex:%d,UnKnownSet:%s}", s.Epoch, s.ViewNumber, s.BlockIndex, s.UnKnownSet.String())
}

func (s *GetPrepareVote) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(GetPrepareVoteMsg, utils.MergeBytes(common.Uint64ToBytes(s.Epoch), common.Uint64ToBytes(s.ViewNumber),
		common.Uint32ToBytes(s.BlockIndex)))
	s.messageHash.Store(v)
	return v
}

func (s *GetPrepareVote) BHash() common.Hash {
	return common.Hash{}
}

// Message used to respond to the number of block votes.
type PrepareVotes struct {
	Epoch       uint64
	ViewNumber  uint64
	BlockIndex  uint32
	Votes       []*PrepareVote // Block voting set.
	messageHash atomic.Value   `json:"-" rlp:"-"`
}

func (s *PrepareVotes) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockIndex:%d}", s.Epoch, s.ViewNumber, s.BlockIndex)
}

func (s *PrepareVotes) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(PrepareVotesMsg, utils.MergeBytes(common.Uint64ToBytes(s.Epoch), common.Uint64ToBytes(s.ViewNumber), common.Uint32ToBytes(s.BlockIndex)))
	s.messageHash.Store(v)
	return v
}

func (s *PrepareVotes) BHash() common.Hash {
	return common.Hash{}
}

// Represents the hash of the proposed block for secondary propagation.
type PrepareBlockHash struct {
	Epoch       uint64
	ViewNumber  uint64
	BlockIndex  uint32
	BlockHash   common.Hash
	BlockNumber uint64
	messageHash atomic.Value `json:"-" rlp:"-"`
}

func (s *PrepareBlockHash) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockIndex:%d,Hash:%s,Number:%d}", s.Epoch, s.ViewNumber, s.BlockIndex, s.BlockHash.TerminalString(), s.BlockNumber)
}

func (s *PrepareBlockHash) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(PrepareBlockHashMsg, utils.MergeBytes(s.BlockHash.Bytes(), common.Uint64ToBytes(s.BlockNumber),
		common.Uint32ToBytes(s.BlockIndex), common.Uint64ToBytes(s.ViewNumber)))
	s.messageHash.Store(v)
	return v

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
	QC           []*ctypes.QuorumCert
	Blocks       []*types.Block
	ForkedQC     []*ctypes.QuorumCert
	ForkedBlocks []*types.Block
	messageHash  atomic.Value `rlp:"-"`
}

func (s *QCBlockList) String() string {
	return fmt.Sprintf("{QC.Len:%d,Blocks.Len:%d,ForkedQC.Len:%d,ForkedBlocks.Len:%d}", len(s.QC), len(s.Blocks), len(s.ForkedQC), len(s.ForkedBlocks))
}

func (s *QCBlockList) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := common.Hash{}
	if len(s.QC) != 0 {
		v = utils.BuildHash(QCBlockListMsg, utils.MergeBytes(s.QC[0].BlockHash.Bytes(),
			s.QC[0].Signature.Bytes()))
	}
	if len(s.Blocks) != 0 {
		v = utils.BuildHash(QCBlockListMsg, utils.MergeBytes(s.Blocks[0].Hash().Bytes(),
			s.Blocks[0].Number().Bytes()))
	}
	s.messageHash.Store(v)
	return v
}

func (s *QCBlockList) BHash() common.Hash {
	// No explicit hash value and return empty hash.
	return common.Hash{}
}

// State synchronization for nodes.
type GetLatestStatus struct {
	BlockNumber  uint64             `json:"blockNumber"`           // QC Block height
	BlockHash    common.Hash        `json:"blockHash"`             // QC block hash
	QuorumCert   *ctypes.QuorumCert `json:"quorumCert" rlp:"nil"`  // QC quorumCert
	LBlockNumber uint64             `json:"lBlockNumber"`          // Locked block height
	LBlockHash   common.Hash        `json:"lBlockHash"`            // Locked block hash
	LQuorumCert  *ctypes.QuorumCert `json:"lQuorumCert" rlp:"nil"` // Locked quorumCert
	LogicType    uint64             `json:"logicType"`             // LogicType: 1 QCBn, 2 LockedBn, 3 CommitBn
	messageHash  atomic.Value       `rlp:"-"`
}

func (s *GetLatestStatus) String() string {
	return fmt.Sprintf("{BlockNumber:%d,BlockHash:%s,QuorumCert:%s,LBlockNumber:%d,LBlockHash:%s,LQuorumCert:%s,LogicType:%d}",
		s.BlockNumber, s.BlockHash.String(), s.QuorumCert.String(), s.LBlockNumber, s.LBlockHash.String(), s.LQuorumCert.String(), s.LogicType)
}

func (s *GetLatestStatus) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(GetLatestStatusMsg,
		utils.MergeBytes(common.Uint64ToBytes(s.BlockNumber), common.Uint64ToBytes(s.LogicType),
			common.Uint64ToBytes(s.LBlockNumber), s.BlockHash.Bytes(), s.LBlockHash.Bytes()))
	s.messageHash.Store(v)
	return v

}

func (s *GetLatestStatus) BHash() common.Hash {
	return s.BlockHash
}

// Response message to GetLatestStatus request.
type LatestStatus struct {
	BlockNumber  uint64             `json:"blockNumber"`           // QC Block height
	BlockHash    common.Hash        `json:"blockHash"`             // QC block hash
	QuorumCert   *ctypes.QuorumCert `json:"quorumCert" rlp:"nil"`  // QC quorumCert
	LBlockNumber uint64             `json:"lBlockNumber"`          // Locked block height
	LBlockHash   common.Hash        `json:"lBlockHash"`            // Locked block hash
	LQuorumCert  *ctypes.QuorumCert `json:"lQuorumCert" rlp:"nil"` // Locked quorumCert
	LogicType    uint64             `json:"logicType"`             // LogicType: 1 QCBn, 2 LockedBn, 3 CommitBn
	messageHash  atomic.Value       `rlp:"-"`
}

func (s *LatestStatus) String() string {
	return fmt.Sprintf("{BlockNumber:%d,BlockHash:%s,QuorumCert:%s,LBlockNumber:%d,LBlockHash:%s,LQuorumCert:%s,LogicType:%d}",
		s.BlockNumber, s.BlockHash.String(), s.QuorumCert.String(), s.LBlockNumber, s.LBlockHash.String(), s.LQuorumCert.String(), s.LogicType)
}

func (s *LatestStatus) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(LatestStatusMsg,
		utils.MergeBytes(common.Uint64ToBytes(s.BlockNumber), common.Uint64ToBytes(s.LogicType),
			common.Uint64ToBytes(s.LBlockNumber), s.BlockHash.Bytes(), s.LBlockHash.Bytes()))
	s.messageHash.Store(v)
	return v
}

func (s *LatestStatus) BHash() common.Hash {
	return s.BlockHash
}

// Used to actively request to get viewChange.
type GetViewChange struct {
	Epoch          uint64          `json:"epoch"`
	ViewNumber     uint64          `json:"viewNumber"`
	ViewChangeBits *utils.BitArray `json:"nodeIndexes"`
	messageHash    atomic.Value    `rlp:"-"`
}

func (s *GetViewChange) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,NodeIndexesLen:%s}", s.Epoch, s.ViewNumber, s.ViewChangeBits.String())
}

func (s *GetViewChange) MsgHash() common.Hash {
	if mhash := s.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	v := utils.BuildHash(GetViewChangeMsg,
		utils.MergeBytes(common.Uint64ToBytes(s.Epoch), common.Uint64ToBytes(s.ViewNumber)))
	s.messageHash.Store(v)
	return v
}

func (s *GetViewChange) BHash() common.Hash {
	return common.Hash{}
}

type ViewChangeQuorumCert struct {
	ViewChangeQC        *ctypes.ViewChangeQC `json:"viewchangeQC"`                  // viewChange aggregate signature
	HighestViewChangeQC *ctypes.ViewChangeQC `json:"highestViewChangeQC" rlp:"nil"` // the highest viewChangeQC of current epoch
	messageHash         atomic.Value         `rlp:"-"`
}

func (v *ViewChangeQuorumCert) String() string {
	epoch, viewNumber, blockEpoch, blockViewNumber, hash, number := v.ViewChangeQC.MaxBlock()
	if v.HighestViewChangeQC == nil {
		return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockEpoch:%d,BlockViewNumber:%d,Hash:%s,Number:%d}",
			epoch, viewNumber, blockEpoch, blockViewNumber, hash.TerminalString(), number)
	}
	highestEpoch, highestViewNumber, highestBlockEpoch, highestBlockViewNumber, highestHash, highestNumber := v.HighestViewChangeQC.MaxBlock()
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockEpoch:%d,BlockViewNumber:%d,Hash:%s,Number:%d,"+
		"HighestEpoch:%d,HighestViewNumber:%d,HighestBlockEpoch:%d,HighestBlockViewNumber:%d,HighestHash:%s,HighestNumber:%d}",
		epoch, viewNumber, blockEpoch, blockViewNumber, hash.TerminalString(), number,
		highestEpoch, highestViewNumber, highestBlockEpoch, highestBlockViewNumber, highestHash.TerminalString(), highestNumber)
}

func (v *ViewChangeQuorumCert) MsgHash() common.Hash {
	if mhash := v.messageHash.Load(); mhash != nil {
		return mhash.(common.Hash)
	}
	epoch, viewNumber, blockEpoch, blockViewNumber, hash, number := v.ViewChangeQC.MaxBlock()
	mv := utils.BuildHash(ViewChangeQuorumCertMsg, utils.MergeBytes(
		common.Uint64ToBytes(epoch),
		common.Uint64ToBytes(viewNumber),
		common.Uint64ToBytes(blockEpoch),
		common.Uint64ToBytes(blockViewNumber),
		hash.Bytes(),
		common.Uint64ToBytes(number)))
	v.messageHash.Store(mv)
	return mv
}

func (v *ViewChangeQuorumCert) BHash() common.Hash {
	_, _, _, _, hash, _ := v.ViewChangeQC.MaxBlock()
	return hash
}
