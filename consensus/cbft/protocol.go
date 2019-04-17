package cbft

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
	"reflect"
)

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

type msgInfo struct {
	msg    interface{}
	peerID discover.NodeID
}

// CBFT consensus message
type prepareBlock struct {
	Timestamp       uint64
	Block           *types.Block
	ProposalIndex   uint32            `json:"proposal_index"`
	ProposalAddr    common.Address    `json:"proposal_address"`
	View            *viewChange       `json:"view"`
	ViewChangeVotes []*viewChangeVote `json:"viewchange_votes"`
}

func (pb *prepareBlock) String() string {
	if pb == nil {
		return ""
	}
	return fmt.Sprintf("[Timestamp:%d Hash:%s Number:%d ProposalAddr:%s, ProposalIndex:%d]", pb.Timestamp, pb.Block.Hash().TerminalString(), pb.Block.NumberU64(), pb.ProposalAddr.String(), pb.ProposalIndex)
}

type prepareVote struct {
	Timestamp      uint64
	Hash           common.Hash
	Number         uint64
	ValidatorIndex uint32
	ValidatorAddr  common.Address
	Signature      common.BlockConfirmSign
}

func (pv *prepareVote) String() string {
	if pv == nil {
		return ""
	}
	return fmt.Sprintf("[Timestamp:%d Hash:%s Number:%d ValidatorAddr:%s ValidatorIndex:%d]", pv.Timestamp, pv.Hash.TerminalString(), pv.Number, pv.ValidatorAddr.String(), pv.ValidatorIndex)
}

type viewChange struct {
	Timestamp                uint64                  `json:"timestamp"`
	ProposalIndex            uint32                  `json:"proposal_index"`
	ProposalAddr             common.Address          `json:"proposal_address"`
	HighestBlockNum          uint64                  `json:"highest_block_number"`
	HighestBlockHash         common.Hash             `json:"highest_block_hash"`
	HighestBlockPrepareVotes []*prepareVote          `json:"highest_block_prepare_votes"`
	Signature                common.BlockConfirmSign `json:"-"`
}

func (v *viewChange) String() string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("[Timestamp:%d ProposalIndex:%d ProposalAddr:%s HighestBlockNum:%d HighestBlockHash:%s]",
		v.Timestamp, v.ProposalIndex, v.ProposalAddr.String(), v.HighestBlockNum, v.HighestBlockHash.TerminalString())
}

func (v *viewChange) Equal(view *viewChange) bool {
	return v.Timestamp == view.Timestamp &&
		v.HighestBlockNum == view.HighestBlockNum &&
		v.HighestBlockHash == view.HighestBlockHash
}

func (v *viewChange) CopyWithoutVotes() *viewChange {
	return &viewChange{
		Timestamp:        v.Timestamp,
		ProposalIndex:    v.ProposalIndex,
		ProposalAddr:     v.ProposalAddr,
		HighestBlockNum:  v.HighestBlockNum,
		HighestBlockHash: v.HighestBlockHash,
		Signature:        v.Signature,
	}
}
func (v *viewChange) Copy() *viewChange {
	view := &viewChange{
		Timestamp:                v.Timestamp,
		ProposalIndex:            v.ProposalIndex,
		ProposalAddr:             v.ProposalAddr,
		HighestBlockNum:          v.HighestBlockNum,
		HighestBlockHash:         v.HighestBlockHash,
		HighestBlockPrepareVotes: make([]*prepareVote, len(v.HighestBlockPrepareVotes)),
		Signature:                v.Signature,
	}
	for i, pv := range v.HighestBlockPrepareVotes {
		view.HighestBlockPrepareVotes[i] = pv
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
	ValidatorAddr  common.Address          `json:"-"`
	Signature      common.BlockConfirmSign `json:"-"`
}

func (v *viewChangeVote) String() string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("[Timestamp:%d BlockNum:%d BlockHash:%s ValidatorIndex:%d ValidatorAddr:%s]",
		v.Timestamp, v.BlockNum, v.BlockHash.TerminalString(), v.ValidatorIndex, v.ValidatorAddr.String())
}

func (v *viewChangeVote) EqualViewChange(vote *viewChange) bool {
	return v.Timestamp == vote.Timestamp &&
		v.BlockNum == vote.HighestBlockNum &&
		v.BlockHash == vote.HighestBlockHash &&
		v.ProposalIndex == vote.ProposalIndex &&
		v.ProposalAddr == vote.ProposalAddr
}

func (v *viewChangeVote) ViewChangeWithSignature() *viewChange {
	return &viewChange{
		Timestamp:        v.Timestamp,
		HighestBlockNum:  v.BlockNum,
		HighestBlockHash: v.BlockHash,
		ProposalIndex:    v.ProposalIndex,
		ProposalAddr:     v.ProposalAddr,
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

type getHighestPrepareBlock struct {
	Lowest uint64
}

func (gpb *getHighestPrepareBlock) String() string {
	if gpb == nil {
		return ""
	}
	return fmt.Sprintf("[Lowest:%d]", gpb.Lowest)
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

type getPrepareBlock struct {
	Hash   common.Hash
	Number uint64
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
	}
	panic(fmt.Sprintf("invalid msg type %v", reflect.TypeOf(msg)))
}
