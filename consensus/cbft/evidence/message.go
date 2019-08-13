package evidence

import (
	"errors"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"

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
	Signature    ctypes.Signature `json:"signature"`
}

func NewEvidencePrepare(pb *protocols.PrepareBlock, node *cbfttypes.ValidateNode) (*EvidencePrepare, error) {
	return &EvidencePrepare{
		Epoch:        pb.Epoch,
		ViewNumber:   pb.ViewNumber,
		BlockHash:    pb.Block.Hash(),
		BlockNumber:  pb.Block.NumberU64(),
		BlockIndex:   pb.BlockIndex,
		ValidateNode: NewEvidenceNode(node),
		Signature:    pb.Signature,
	}, nil
}

func (ep *EvidencePrepare) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		ep.Epoch,
		ep.ViewNumber,
		ep.BlockHash,
		ep.BlockIndex,
		ep.ValidateNode.Index,
	})
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (ep *EvidencePrepare) Verify() error {
	data, err := ep.CannibalizeBytes()
	if err != nil {
		return err
	}
	return ep.ValidateNode.Verify(data, ep.Signature.Bytes())
}

type EvidenceVote struct {
	Epoch        uint64           `json:"epoch"`
	ViewNumber   uint64           `json:"view_number"`
	BlockHash    common.Hash      `json:"block_hash"`
	BlockNumber  uint64           `json:"block_number"`
	BlockIndex   uint32           `json:"block_index"` // The block number of the current ViewNumber proposal, 0....10
	ValidateNode *EvidenceNode    `json:"validate_node"`
	Signature    ctypes.Signature `json:"signature"`
}

func NewEvidenceVote(pv *protocols.PrepareVote, node *cbfttypes.ValidateNode) (*EvidenceVote, error) {
	return &EvidenceVote{
		Epoch:        pv.Epoch,
		ViewNumber:   pv.ViewNumber,
		BlockHash:    pv.BlockHash,
		BlockNumber:  pv.BlockNumber,
		BlockIndex:   pv.BlockIndex,
		ValidateNode: NewEvidenceNode(node),
		Signature:    pv.Signature,
	}, nil
}

func (ev *EvidenceVote) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		ev.Epoch,
		ev.ViewNumber,
		ev.BlockHash,
		ev.BlockNumber,
		ev.BlockIndex,
	})

	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (ev *EvidenceVote) Verify() error {
	data, err := ev.CannibalizeBytes()
	if err != nil {
		return err
	}
	return ev.ValidateNode.Verify(data, ev.Signature.Bytes())
}

type EvidenceView struct {
	Epoch        uint64           `json:"epoch"`
	ViewNumber   uint64           `json:"view_number"`
	BlockHash    common.Hash      `json:"block_hash"`
	BlockNumber  uint64           `json:"block_number"`
	ValidateNode *EvidenceNode    `json:"validate_node"`
	Signature    ctypes.Signature `json:"signature"`
	BlockEpoch   uint64           `json:"block_epoch"`
	BlockView    uint64           `json:"block_view"`
}

func NewEvidenceView(vc *protocols.ViewChange, node *cbfttypes.ValidateNode) (*EvidenceView, error) {
	return &EvidenceView{
		Epoch:        vc.Epoch,
		ViewNumber:   vc.ViewNumber,
		BlockHash:    vc.BlockHash,
		BlockNumber:  vc.BlockNumber,
		ValidateNode: NewEvidenceNode(node),
		Signature:    vc.Signature,
	}, nil
}

func (ev *EvidenceView) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		ev.Epoch,
		ev.ViewNumber,
		ev.BlockHash,
		ev.BlockNumber,
		ev.BlockEpoch,
		ev.BlockView,
	})

	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (ev *EvidenceView) Verify() error {
	data, err := ev.CannibalizeBytes()
	if err != nil {
		return err
	}
	return ev.ValidateNode.Verify(data, ev.Signature.Bytes())
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
