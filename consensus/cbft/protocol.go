package cbft

import (
	"encoding/json"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"math/big"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

const CbftProtocolMaxMsgSize = 10 * 1024 * 1024

const (
	PrepareBlockMsg           = 0x00
	PrepareVoteMsg            = 0x01
	ViewChangeMsg             = 0x02
	ViewChangeVoteMsg         = 0x03
	ConfirmedPrepareBlockMsg  = 0x04
	GetPrepareVoteMsg         = 0x05
	PrepareVotesMsg           = 0x06
	GetPrepareBlockMsg        = 0x07
	GetHighestPrepareBlockMsg = 0x08
	HighestPrepareBlockMsg    = 0x09

	CBFTStatusMsg       = 0x0a
	PrepareBlockHashMsg = 0x0b
)

type errCode int

const (
	ErrMsgTooLarge = iota
	ErrDecode
	ErrInvalidMsgCode
	ErrProtocolVersionMismatch
	ErrNetworkIdMismatch
	ErrGenesisBlockMismatch
	ErrNoStatusMsg
	ErrExtraStatusMsg
	ErrSuspendedPeer
)

func (e errCode) String() string {
	return errorToString[int(e)]
}

// XXX change once legacy code is out
var errorToString = map[int]string{
	ErrMsgTooLarge:             "Message too long",
	ErrDecode:                  "Invalid message",
	ErrInvalidMsgCode:          "Invalid message code",
	ErrProtocolVersionMismatch: "Protocol version mismatch",
	ErrNetworkIdMismatch:       "NetworkId mismatch",
	ErrGenesisBlockMismatch:    "Genesis block mismatch",
	ErrNoStatusMsg:             "No status message",
	ErrExtraStatusMsg:          "Extra status message",
	ErrSuspendedPeer:           "Suspended peer",
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

// CBFT consensus message
type prepareBlock struct {
	Timestamp       uint64 `json:"timestamp"`
	Block           *types.Block
	ProposalIndex   uint32            `json:"proposal_index"`
	ProposalAddr    common.Address    `json:"proposal_address"`
	View            *viewChange       `json:"view"`
	ViewChangeVotes []*viewChangeVote `json:"viewchange_votes"`
	Extra           []byte
}

func (pb prepareBlock) MarshalJSON() ([]byte, error) {
	type prepareBlock struct {
		Timestamp       uint64            `json:"timestamp"`
		BlockHash       common.Hash       `json:"block_hash"`
		BlockNumber     uint64            `json:"block_number"`
		ProposalIndex   uint32            `json:"proposal_index"`
		ProposalAddr    common.Address    `json:"proposal_address"`
		View            *viewChange       `json:"view"`
		ViewChangeVotes []*viewChangeVote `json:"viewchange_votes"`
	}

	var p prepareBlock

	p.Timestamp = pb.Timestamp
	p.BlockHash = pb.Block.Hash()
	p.BlockNumber = pb.Block.NumberU64()
	p.ProposalIndex = pb.ProposalIndex
	p.ProposalAddr = pb.ProposalAddr
	p.View = pb.View
	p.ViewChangeVotes = pb.ViewChangeVotes

	return json.Marshal(&p)
}

func (pb *prepareBlock) CannibalizeBytes() ([]byte, error) {
	return pb.Block.Header().SealHash().Bytes(), nil
}

func (pb *prepareBlock) Sign() []byte {
	return pb.Block.Extra()[len(pb.Block.Extra())-extraSeal:]
}

func (pb *prepareBlock) String() string {
	if pb == nil {
		return ""
	}
	return fmt.Sprintf("[Timestamp:%d Hash:%s Number:%d ProposalAddr:%s, ProposalIndex:%d]", pb.Timestamp, pb.Block.Hash().TerminalString(), pb.Block.NumberU64(), pb.ProposalAddr.String(), pb.ProposalIndex)
}

func (pb *prepareBlock) MsgHash() common.Hash {
	if pb == nil {
		return common.Hash{}
	}
	bytes := make([]byte, 0)
	bytes = append(bytes, pb.Block.Hash().Bytes()...)
	bytes = append(bytes, pb.ProposalAddr.Bytes()...)
	bytes = append(bytes, uint64ToBytes(pb.Timestamp)...)
	return produceHash(PrepareBlockMsg, bytes)
}

func (pb *prepareBlock) BHash() common.Hash {
	if pb == nil {
		return common.Hash{}
	}
	return pb.Block.Hash()
}

type prepareBlockHash struct {
	Hash   common.Hash
	Number uint64
}

func (pbh *prepareBlockHash) String() string {
	if pbh == nil {
		return ""
	}
	return fmt.Sprintf("[Hash:%s, Number:%d]", pbh.Hash.TerminalString(), pbh.Number)
}

func (pbh *prepareBlockHash) MsgHash() common.Hash {
	if pbh == nil {
		return common.Hash{}
	}
	return produceHash(PrepareBlockHashMsg, pbh.Hash.Bytes())
}

func (pbh *prepareBlockHash) BHash() common.Hash {
	if pbh == nil {
		return common.Hash{}
	}
	return pbh.Hash
}

type prepareVote struct {
	Timestamp      uint64                  `json:"timestamp"`
	Hash           common.Hash             `json:"block_hash"`
	Number         uint64                  `json:"block_number"`
	ValidatorIndex uint32                  `json:"validator_index"`
	ValidatorAddr  common.Address          `json:"validator_address"`
	Signature      common.BlockConfirmSign `json:"signature"`
	Extra          []byte                  `json:"-"`
}

func (pv *prepareVote) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		pv.Timestamp,
		pv.Hash,
		pv.Number,
		pv.ValidatorIndex,
		pv.ValidatorAddr,
	})

	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(buf), nil
}

func (pv *prepareVote) Sign() []byte {
	return pv.Signature.Bytes()
}

func (pv *prepareVote) String() string {
	if pv == nil {
		return ""
	}
	return fmt.Sprintf("[Timestamp:%d Hash:%s Number:%d ValidatorAddr:%s ValidatorIndex:%d]", pv.Timestamp, pv.Hash.TerminalString(), pv.Number, pv.ValidatorAddr.String(), pv.ValidatorIndex)
}

func (pv *prepareVote) MsgHash() common.Hash {
	if pv == nil {
		return common.Hash{}
	}
	bytes := make([]byte, 0)
	bytes = append(bytes, pv.Hash.Bytes()...)
	bytes = append(bytes, pv.ValidatorAddr.Bytes()...)
	bytes = append(bytes, uint64ToBytes(pv.Timestamp)...)
	return produceHash(PrepareVoteMsg, bytes)
}

func (pv *prepareVote) BHash() common.Hash {
	if pv == nil {
		return common.Hash{}
	}
	return pv.Hash
}

type viewChange struct {
	Timestamp            uint64                  `json:"timestamp"`
	ProposalIndex        uint32                  `json:"proposal_index"`
	ProposalAddr         common.Address          `json:"proposal_address"`
	BaseBlockNum         uint64                  `json:"base_block_number"`
	BaseBlockHash        common.Hash             `json:"base_block_hash"`
	BaseBlockPrepareVote []*prepareVote          `json:"base_block_prepare_votes"`
	Signature            common.BlockConfirmSign `json:"signature"`
	Extra                []byte                  `json:"-"`
}

func (v *viewChange) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		v.Timestamp,
		v.ProposalIndex,
		v.ProposalAddr,
		v.BaseBlockNum,
		v.BaseBlockHash,
	})

	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(buf), nil
}

func (v *viewChange) Sign() []byte {
	return v.Signature.Bytes()
}

func (v *viewChange) String() string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("[Timestamp:%d ProposalIndex:%d ProposalAddr:%s BaseBlockNum:%d BaseBlockHash:%s]",
		v.Timestamp, v.ProposalIndex, v.ProposalAddr.String(), v.BaseBlockNum, v.BaseBlockHash.TerminalString())
}

func (v *viewChange) MsgHash() common.Hash {
	if v == nil {
		return common.Hash{}
	}
	bytes := make([]byte, 0)
	bytes = append(bytes, v.Signature.Bytes()...)
	bytes = append(bytes, v.ProposalAddr.Bytes()...)
	bytes = append(bytes, uint64ToBytes(v.Timestamp)...)
	return produceHash(ViewChangeMsg, bytes)
}

func (v *viewChange) BHash() common.Hash {
	if v == nil {
		return common.Hash{}
	}
	return common.BytesToHash(v.Signature.Bytes())
}

func (v *viewChange) Equal(view *viewChange) bool {
	return v.Timestamp == view.Timestamp &&
		v.BaseBlockNum == view.BaseBlockNum &&
		v.BaseBlockHash == view.BaseBlockHash
}

func (v *viewChange) CopyWithoutVotes() *viewChange {
	return &viewChange{
		Timestamp:     v.Timestamp,
		ProposalIndex: v.ProposalIndex,
		ProposalAddr:  v.ProposalAddr,
		BaseBlockNum:  v.BaseBlockNum,
		BaseBlockHash: v.BaseBlockHash,
		Signature:     v.Signature,
	}
}
func (v *viewChange) Copy() *viewChange {
	view := &viewChange{
		Timestamp:            v.Timestamp,
		ProposalIndex:        v.ProposalIndex,
		ProposalAddr:         v.ProposalAddr,
		BaseBlockNum:         v.BaseBlockNum,
		BaseBlockHash:        v.BaseBlockHash,
		BaseBlockPrepareVote: make([]*prepareVote, len(v.BaseBlockPrepareVote)),
		Signature:            v.Signature,
	}
	for i, pv := range v.BaseBlockPrepareVote {
		view.BaseBlockPrepareVote[i] = pv
	}
	return view
}

type viewChangeVote struct {
	Timestamp      uint64                  `json:"timestamp"`
	BlockNum       uint64                  `json:"block_number"`
	BlockHash      common.Hash             `json:"block_hash"`
	ProposalIndex  uint32                  `json:"proposal_index"`
	ProposalAddr   common.Address          `json:"proposal_address"`
	ValidatorIndex uint32                  `json:"validator_index"`
	ValidatorAddr  common.Address          `json:"validator_address"`
	Signature      common.BlockConfirmSign `json:"signature"`
	Extra          []byte                  `json:"-"`
}

func (v *viewChangeVote) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		v.Timestamp,
		v.BlockNum,
		v.BlockHash,
		v.ProposalIndex,
		v.ProposalAddr,
		v.ValidatorIndex,
		v.ValidatorAddr,
	})

	if err != nil {
		return nil, err
	}

	return crypto.Keccak256(buf), nil
}

func (v *viewChangeVote) Sign() []byte {
	return v.Signature.Bytes()
}

func (v *viewChangeVote) String() string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("[Timestamp:%d BlockNum:%d BlockHash:%s ValidatorIndex:%d ValidatorAddr:%s]",
		v.Timestamp, v.BlockNum, v.BlockHash.TerminalString(), v.ValidatorIndex, v.ValidatorAddr.String())
}

func (v *viewChangeVote) MsgHash() common.Hash {
	if v == nil {
		return common.Hash{}
	}
	bytes := make([]byte, 0)
	bytes = append(bytes, v.Signature.Bytes()...)
	bytes = append(bytes, v.ValidatorAddr.Bytes()[:5]...)
	bytes = append(bytes, uint64ToBytes(v.Timestamp)...)

	return produceHash(ViewChangeVoteMsg, bytes)
}

func (v *viewChangeVote) BHash() common.Hash {
	if v == nil {
		return common.Hash{}
	}
	return common.BytesToHash(v.Signature.Bytes())
}

func (v *viewChangeVote) EqualViewChange(vote *viewChange) bool {
	return v.Timestamp == vote.Timestamp &&
		v.BlockNum == vote.BaseBlockNum &&
		v.BlockHash == vote.BaseBlockHash &&
		v.ProposalIndex == vote.ProposalIndex &&
		v.ProposalAddr == vote.ProposalAddr
}

func (v *viewChangeVote) ViewChangeWithSignature() *viewChange {
	return &viewChange{
		Timestamp:     v.Timestamp,
		BaseBlockNum:  v.BlockNum,
		BaseBlockHash: v.BlockHash,
		ProposalIndex: v.ProposalIndex,
		ProposalAddr:  v.ProposalAddr,
	}
}

type confirmedPrepareBlock struct {
	Hash     common.Hash
	Number   uint64
	VoteBits *BitArray
}

func (cpb *confirmedPrepareBlock) String() string {
	if cpb == nil {
		return ""
	}
	return fmt.Sprintf("[Hash:%s Number:%d VoteBits:%s]", cpb.Hash.String(), cpb.Number, cpb.VoteBits.String())
}

func (cpb *confirmedPrepareBlock) MsgHash() common.Hash {
	if cpb == nil {
		return common.Hash{}
	}
	return produceHash(ConfirmedPrepareBlockMsg, cpb.Hash.Bytes())
}

func (cpb *confirmedPrepareBlock) BHash() common.Hash {
	if cpb == nil {
		return common.Hash{}
	}
	return cpb.Hash
}

type getHighestPrepareBlock struct {
	Lowest uint64
}

func (gpb *getHighestPrepareBlock) String() string {
	if gpb == nil {
		return ""
	}
	return fmt.Sprintf("[Lowest:%d]", gpb.Lowest)
}

func (gpb *getHighestPrepareBlock) MsgHash() common.Hash {
	if gpb == nil {
		return common.Hash{}
	}
	return produceHash(GetHighestPrepareBlockMsg, common.BigToHash(new(big.Int).SetUint64(gpb.Lowest)).Bytes())
}

func (gpb *getHighestPrepareBlock) BHash() common.Hash {
	if gpb == nil {
		return common.Hash{}
	}
	return common.BigToHash(new(big.Int).SetUint64(gpb.Lowest))
}

type highestPrepareBlock struct {
	CommitedBlock    []*types.Block
	UnconfirmedBlock []*prepareBlock
	Votes            []*prepareVotes
}

func (pb *highestPrepareBlock) String() string {
	if pb == nil {
		return ""
	}
	return fmt.Sprintf("[CommitedBlock:%d UnconfirmedBlock:%d, Votes:%d]", len(pb.CommitedBlock), len(pb.UnconfirmedBlock), len(pb.Votes))
}

func (pb *highestPrepareBlock) MsgHash() common.Hash {
	if pb == nil {
		return common.Hash{}
	}
	return produceHash(HighestPrepareBlockMsg, common.Hash{}.Bytes())
}

func (pb *highestPrepareBlock) BHash() common.Hash {
	if pb == nil {
		return common.Hash{}
	}
	return common.Hash{}
}

type getPrepareBlock struct {
	Hash   common.Hash
	Number uint64
}

func (gpb *getPrepareBlock) String() string {
	if gpb == nil {
		return ""
	}
	return fmt.Sprintf("[Hash:%s] Number:%d", gpb.Hash.String(), gpb.Number)
}

func (gpb *getPrepareBlock) MsgHash() common.Hash {
	if gpb == nil {
		return common.Hash{}
	}
	return produceHash(GetPrepareBlockMsg, gpb.Hash.Bytes())
}

func (gpb *getPrepareBlock) BHash() common.Hash {
	if gpb == nil {
		return common.Hash{}
	}
	return gpb.Hash
}

type getPrepareVote struct {
	Hash     common.Hash
	Number   uint64
	VoteBits *BitArray
}

func (pv *getPrepareVote) String() string {
	if pv == nil {
		return ""
	}
	return fmt.Sprintf("[Hash:%s Number:%d VoteBits:%s]", pv.Hash.String(), pv.Number, pv.VoteBits.String())
}

func (pv *getPrepareVote) MsgHash() common.Hash {
	if pv == nil {
		return common.Hash{}
	}
	return produceHash(GetPrepareVoteMsg, pv.Hash.Bytes())
}

func (pv *getPrepareVote) BHash() common.Hash {
	if pv == nil {
		return common.Hash{}
	}
	return pv.Hash
}

type prepareVotes struct {
	Hash   common.Hash
	Number uint64
	Votes  []*prepareVote
}

func (pv *prepareVotes) String() string {
	if pv == nil {
		return ""
	}
	return fmt.Sprintf("[Hash:%s Number:%d prepareVotes:%d]", pv.Hash.String(), pv.Number, len(pv.Votes))
}

func (pv *prepareVotes) MsgHash() common.Hash {
	if pv == nil {
		return common.Hash{}
	}
	return produceHash(PrepareVotesMsg, pv.Hash.Bytes())
}

func (pv *prepareVotes) BHash() common.Hash {
	if pv == nil {
		return common.Hash{}
	}
	return pv.Hash
}

//p2p sync message
type signBitArray struct {
	BlockHash common.Hash `json:"block_hash"`
	BlockNum  uint64      `json:"block_number"`
	SignBits  *BitArray
}

func (v *signBitArray) Copy() *signBitArray {
	return &signBitArray{
		BlockHash: v.BlockHash,
		BlockNum:  v.BlockNum,
		SignBits:  v.SignBits.Copy(),
	}
}

func (v *signBitArray) String() string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("[BlockHash:%s BlockNum:%d]", v.BlockHash.String(), v.BlockNum)
}

func (v *signBitArray) MsgHash() common.Hash {
	if v == nil {
		return common.Hash{}
	}
	return v.BlockHash
}

func (v *signBitArray) BHash() common.Hash {
	if v == nil {
		return common.Hash{}
	}
	return v.BlockHash
}

type cbftStatusData struct {
	BN           *big.Int
	CurrentBlock common.Hash
}

func (s *cbftStatusData) String() string {
	if s == nil {
		return ""
	}
	return fmt.Sprintf("[BlockNumber:%d, BlockHash:%s]", s.BN.Int64(), s.CurrentBlock.String())
}

func (s *cbftStatusData) MsgHash() common.Hash {
	if s == nil {
		return common.Hash{}
	}
	return produceHash(CBFTStatusMsg, s.CurrentBlock.Bytes())
}

func (s *cbftStatusData) BHash() common.Hash {
	if s == nil {
		return common.Hash{}
	}
	return s.CurrentBlock
}

var (
	messages = []interface{}{
		prepareBlock{},
		prepareVote{},
		viewChange{},
		viewChangeVote{},
		confirmedPrepareBlock{},
		getPrepareVote{},
		prepareVotes{},
		getPrepareBlock{},
		getHighestPrepareBlock{},
		highestPrepareBlock{},
		cbftStatusData{},
		prepareBlockHash{},
	}
)

func MessageType(msg interface{}) uint64 {
	switch msg.(type) {
	case *prepareBlock:
		return PrepareBlockMsg
	case *prepareVote:
		return PrepareVoteMsg
	case *viewChange:
		return ViewChangeMsg
	case *viewChangeVote:
		return ViewChangeVoteMsg
	case *confirmedPrepareBlock:
		return ConfirmedPrepareBlockMsg
	case *getPrepareVote:
		return GetPrepareVoteMsg
	case *prepareVotes:
		return PrepareVotesMsg
	case *getPrepareBlock:
		return GetPrepareBlockMsg
	case *getHighestPrepareBlock:
		return GetHighestPrepareBlockMsg
	case *highestPrepareBlock:
		return HighestPrepareBlockMsg
	case *cbftStatusData:
		return CBFTStatusMsg
	case *prepareBlockHash:
		return PrepareBlockHashMsg
	}
	panic(fmt.Sprintf("invalid msg type %v", reflect.TypeOf(msg)))
}
