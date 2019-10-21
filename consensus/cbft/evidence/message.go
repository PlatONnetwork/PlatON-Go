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
	ViewNumber   uint64           `json:"viewNumber"`
	BlockHash    common.Hash      `json:"blockHash"`
	BlockNumber  uint64           `json:"blockNumber"`
	BlockIndex   uint32           `json:"blockIndex"` // The block number of the current ViewNumber proposal, 0....10
	BlockData    common.Hash      `json:"blockData"`
	ValidateNode *EvidenceNode    `json:"validateNode"`
	Signature    ctypes.Signature `json:"signature"`
}

func NewEvidencePrepare(pb *protocols.PrepareBlock, node *cbfttypes.ValidateNode) (*EvidencePrepare, error) {
	blockData, err := rlp.EncodeToBytes(pb.Block)
	if err != nil {
		return nil, err
	}
	return &EvidencePrepare{
		Epoch:        pb.Epoch,
		ViewNumber:   pb.ViewNumber,
		BlockHash:    pb.Block.Hash(),
		BlockNumber:  pb.Block.NumberU64(),
		BlockIndex:   pb.BlockIndex,
		BlockData:    common.BytesToHash(crypto.Keccak256(blockData)),
		ValidateNode: NewEvidenceNode(node),
		Signature:    pb.Signature,
	}, nil
}

func (ep *EvidencePrepare) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		ep.Epoch,
		ep.ViewNumber,
		ep.BlockHash,
		ep.BlockNumber,
		ep.BlockData,
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
	ViewNumber   uint64           `json:"viewNumber"`
	BlockHash    common.Hash      `json:"blockHash"`
	BlockNumber  uint64           `json:"blockNumber"`
	BlockIndex   uint32           `json:"blockIndex"` // The block number of the current ViewNumber proposal, 0....10
	ValidateNode *EvidenceNode    `json:"validateNode"`
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
	ViewNumber   uint64           `json:"viewNumber"`
	BlockHash    common.Hash      `json:"blockHash"`
	BlockNumber  uint64           `json:"blockNumber"`
	ValidateNode *EvidenceNode    `json:"validateNode"`
	Signature    ctypes.Signature `json:"signature"`
	BlockEpoch   uint64           `json:"blockEpoch"`
	BlockView    uint64           `json:"blockView"`
}

func NewEvidenceView(vc *protocols.ViewChange, node *cbfttypes.ValidateNode) (*EvidenceView, error) {
	blockEpoch, blockView := uint64(0), uint64(0)
	if vc.PrepareQC != nil {
		blockEpoch, blockView = vc.PrepareQC.Epoch, vc.PrepareQC.ViewNumber
	}
	return &EvidenceView{
		Epoch:        vc.Epoch,
		ViewNumber:   vc.ViewNumber,
		BlockHash:    vc.BlockHash,
		BlockNumber:  vc.BlockNumber,
		ValidateNode: NewEvidenceNode(node),
		Signature:    vc.Signature,
		BlockEpoch:   blockEpoch,
		BlockView:    blockView,
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

// EvidenceNode mainly used to save node BlsPubKey
type EvidenceNode struct {
	Index     uint32          `json:"index"`
	Address   common.Address  `json:"address"`
	NodeID    discover.NodeID `json:"nodeId"`
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

// bls verifies signature
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
