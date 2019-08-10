package evidence

import (
	"errors"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/protocols"
	ctypes "github.com/PlatONnetwork/PlatON-Go/consensus/cbft/types"
	"github.com/PlatONnetwork/PlatON-Go/core/cbfttypes"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"
)

// Proposed block carrier.
type EvidencePrepare struct {
	Epoch        uint64           `json:"epoch"`
	ViewNumber   uint64           `json:"view_number"`
	BlockHash    common.Hash      `json:"block_hash"`
	BlockNumber  uint64           `json:"block_number"`
	BlockIndex   uint32           `json:"block_index"` // The block number of the current ViewNumber proposal, 0....10
	ValidateNode *EvidenceNode    `json:"validate_node"`
	Cannibalize  []byte           `json:"cannibalize"`
	Signature    ctypes.Signature `json:"signature"`
}

func NewEvidencePrepare(pb *protocols.PrepareBlock, node *cbfttypes.ValidateNode) (*EvidencePrepare, error) {
	cannibalize, err := pb.CannibalizeBytes()
	if err != nil {
		return nil, err
	}
	return &EvidencePrepare{
		Epoch:        pb.Epoch,
		ViewNumber:   pb.ViewNumber,
		BlockHash:    pb.Block.Hash(),
		BlockNumber:  pb.Block.NumberU64(),
		BlockIndex:   pb.BlockIndex,
		ValidateNode: NewEvidenceNode(node),
		Cannibalize:  cannibalize,
		Signature:    pb.Signature,
	}, nil
}

type EvidenceVote struct {
	Epoch        uint64           `json:"epoch"`
	ViewNumber   uint64           `json:"view_number"`
	BlockHash    common.Hash      `json:"block_hash"`
	BlockNumber  uint64           `json:"block_number"`
	BlockIndex   uint32           `json:"block_index"` // The block number of the current ViewNumber proposal, 0....10
	ValidateNode *EvidenceNode    `json:"validate_node"`
	Cannibalize  []byte           `json:"cannibalize"`
	Signature    ctypes.Signature `json:"signature"`
}

func NewEvidenceVote(pv *protocols.PrepareVote, node *cbfttypes.ValidateNode) (*EvidenceVote, error) {
	cannibalize, err := pv.CannibalizeBytes()
	if err != nil {
		return nil, err
	}
	return &EvidenceVote{
		Epoch:        pv.Epoch,
		ViewNumber:   pv.ViewNumber,
		BlockHash:    pv.BlockHash,
		BlockNumber:  pv.BlockNumber,
		BlockIndex:   pv.BlockIndex,
		ValidateNode: NewEvidenceNode(node),
		Cannibalize:  cannibalize,
		Signature:    pv.Signature,
	}, nil
}

type EvidenceView struct {
	Epoch        uint64           `json:"epoch"`
	ViewNumber   uint64           `json:"view_number"`
	BlockHash    common.Hash      `json:"block_hash"`
	BlockNumber  uint64           `json:"block_number"`
	ValidateNode *EvidenceNode    `json:"validate_node"`
	Cannibalize  []byte           `json:"cannibalize"`
	Signature    ctypes.Signature `json:"signature"`
}

func NewEvidenceView(vc *protocols.ViewChange, node *cbfttypes.ValidateNode) (*EvidenceView, error) {
	cannibalize, err := vc.CannibalizeBytes()
	if err != nil {
		return nil, err
	}
	return &EvidenceView{
		Epoch:        vc.Epoch,
		ViewNumber:   vc.ViewNumber,
		BlockHash:    vc.BlockHash,
		BlockNumber:  vc.BlockNumber,
		ValidateNode: NewEvidenceNode(node),
		Cannibalize:  cannibalize,
		Signature:    vc.Signature,
	}, nil
}

type EvidenceNode struct {
	Index     uint32          `json:"index"`
	Address   common.Address  `json:"address"`
	NodeID    discover.NodeID `json:"NodeID"`
	BlsPubKey *bls.PublicKey  `json:"blsPubKey"`
}

func NewEvidenceNode(node *cbfttypes.ValidateNode) *EvidenceNode {
	return &EvidenceNode{
		Index:     node.Index,
		Address:   node.Address,
		NodeID:    node.NodeID,
		BlsPubKey: node.BlsPubKey,
	}
}

func (vn *EvidenceNode) Verify(data, sign []byte) error {
	var sig bls.Sign
	err := sig.Deserialize(sign)
	if err != nil {
		return err
	}

	if !sig.Verify(vn.BlsPubKey, string(data)) {
		return errors.New("bls verifies signature fail")
	}
	return nil
}
