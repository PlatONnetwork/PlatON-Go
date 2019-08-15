package evidence

import (
	"bytes"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

const (
	DuplicatePrepareBlockType = 1
	DuplicatePrepareVoteType  = 2
	DuplicateViewChangeType   = 3
)

type DuplicatePrepareBlockEvidence struct {
	PrepareA *EvidencePrepare
	PrepareB *EvidencePrepare
}

func (d DuplicatePrepareBlockEvidence) BlockNumber() uint64 {
	return d.PrepareA.BlockNumber
}

func (d DuplicatePrepareBlockEvidence) Epoch() uint64 {
	return d.PrepareA.Epoch
}

func (d DuplicatePrepareBlockEvidence) ViewNumber() uint64 {
	return d.PrepareA.ViewNumber
}

func (d DuplicatePrepareBlockEvidence) Hash() []byte {
	var buf []byte
	if ac, err := d.PrepareA.CannibalizeBytes(); err == nil {
		if bc, err := d.PrepareB.CannibalizeBytes(); err == nil {
			buf, err = rlp.EncodeToBytes([]interface{}{
				ac,
				d.PrepareA.Signature.Bytes(),
				bc,
				d.PrepareB.Signature.Bytes(),
			})
		}
	}
	return crypto.Keccak256(buf)
}

func (d DuplicatePrepareBlockEvidence) Equal(ev consensus.Evidence) bool {
	_, ok := ev.(*DuplicatePrepareBlockEvidence)
	if !ok {
		return false
	}
	dh := d.Hash()
	eh := ev.Hash()
	return bytes.Equal(dh, eh)
}

func (d DuplicatePrepareBlockEvidence) Error() string {
	return fmt.Sprintf("DuplicatePrepareBlockEvidence epoch:%d, viewNumber:%d, blockNumber:%d blockHashA:%s, blockHashB:%s",
		d.PrepareA.Epoch, d.PrepareA.ViewNumber, d.PrepareA.BlockNumber, d.PrepareA.BlockHash.String(), d.PrepareB.BlockHash.String())
}

func (d DuplicatePrepareBlockEvidence) Validate() error {
	if d.PrepareA.Epoch != d.PrepareB.Epoch {
		return fmt.Errorf("DuplicatePrepareBlockEvidence Epoch is different, PrepareA:%d, PrepareB:%d", d.PrepareA.Epoch, d.PrepareB.Epoch)
	}
	if d.PrepareA.ViewNumber != d.PrepareB.ViewNumber {
		return fmt.Errorf("DuplicatePrepareBlockEvidence ViewNumber is different, PrepareA:%d, PrepareB:%d", d.PrepareA.ViewNumber, d.PrepareB.ViewNumber)
	}
	if d.PrepareA.BlockNumber != d.PrepareB.BlockNumber {
		return fmt.Errorf("DuplicatePrepareBlockEvidence BlockNumber is different, PrepareA:%d, PrepareB:%d", d.PrepareA.BlockNumber, d.PrepareB.BlockNumber)
	}
	validateNodeA, validateNodeB := d.PrepareA.ValidateNode, d.PrepareB.ValidateNode
	if validateNodeA.Index != validateNodeB.Index || validateNodeA.Address != validateNodeB.Address {
		return fmt.Errorf("DuplicatePrepareBlockEvidence Validator do not match, PrepareA:%s, PrepareB:%s", validateNodeA.Address, validateNodeB.Address)
	}
	if d.PrepareA.BlockHash == d.PrepareB.BlockHash {
		return fmt.Errorf("DuplicatePrepareBlockEvidence BlockHash is equal, PrepareA:%s, PrepareB:%s", d.PrepareA.BlockHash, d.PrepareB.BlockHash)
	}
	// Verify consensus msg signature
	if err := d.PrepareA.Verify(); err != nil {
		return fmt.Errorf("DuplicatePrepareBlockEvidence prepareA verify failed")
	}
	if err := d.PrepareB.Verify(); err != nil {
		return fmt.Errorf("DuplicatePrepareBlockEvidence prepareB verify failed")
	}
	return nil
}

func (d DuplicatePrepareBlockEvidence) Address() common.Address {
	return d.PrepareA.ValidateNode.Address
}

func (d DuplicatePrepareBlockEvidence) Type() consensus.EvidenceType {
	return DuplicatePrepareBlockType
}

type DuplicatePrepareVoteEvidence struct {
	VoteA *EvidenceVote
	VoteB *EvidenceVote
}

func (d DuplicatePrepareVoteEvidence) BlockNumber() uint64 {
	return d.VoteA.BlockNumber
}

func (d DuplicatePrepareVoteEvidence) Epoch() uint64 {
	return d.VoteA.Epoch
}

func (d DuplicatePrepareVoteEvidence) ViewNumber() uint64 {
	return d.VoteA.ViewNumber
}

func (d DuplicatePrepareVoteEvidence) Hash() []byte {
	var buf []byte
	if ac, err := d.VoteA.CannibalizeBytes(); err == nil {
		if bc, err := d.VoteB.CannibalizeBytes(); err == nil {
			buf, err = rlp.EncodeToBytes([]interface{}{
				ac,
				d.VoteA.Signature.Bytes(),
				bc,
				d.VoteB.Signature.Bytes(),
			})
		}
	}
	return crypto.Keccak256(buf)
}

func (d DuplicatePrepareVoteEvidence) Equal(ev consensus.Evidence) bool {
	_, ok := ev.(*DuplicatePrepareVoteEvidence)
	if !ok {
		return false
	}
	dh := d.Hash()
	eh := ev.Hash()
	return bytes.Equal(dh, eh)
}

func (d DuplicatePrepareVoteEvidence) Error() string {
	return fmt.Sprintf("DuplicatePrepareVoteEvidence epoch:%d, viewNumber:%d, blockNumber:%d blockHashA:%s, blockHashB:%s",
		d.VoteA.Epoch, d.VoteA.ViewNumber, d.VoteA.BlockNumber, d.VoteA.BlockHash.String(), d.VoteB.BlockHash.String())
}

func (d DuplicatePrepareVoteEvidence) Validate() error {
	if d.VoteA.Epoch != d.VoteB.Epoch {
		return fmt.Errorf("DuplicatePrepareVoteEvidence Epoch is different, VoteA:%d, VoteB:%d", d.VoteA.Epoch, d.VoteB.Epoch)
	}
	if d.VoteA.ViewNumber != d.VoteB.ViewNumber {
		return fmt.Errorf("DuplicatePrepareVoteEvidence ViewNumber is different, VoteA:%d, VoteB:%d", d.VoteA.ViewNumber, d.VoteB.ViewNumber)
	}
	if d.VoteA.BlockNumber != d.VoteB.BlockNumber {
		return fmt.Errorf("DuplicatePrepareVoteEvidence BlockNumber is different, VoteA:%d, VoteB:%d", d.VoteA.BlockNumber, d.VoteB.BlockNumber)
	}
	validateNodeA, validateNodeB := d.VoteA.ValidateNode, d.VoteB.ValidateNode
	if validateNodeA.Index != validateNodeB.Index || validateNodeA.Address != validateNodeB.Address {
		return fmt.Errorf("DuplicatePrepareVoteEvidence Validator do not match, VoteA:%s, VoteB:%s", validateNodeA.Address, validateNodeB.Address)
	}
	if d.VoteA.BlockHash == d.VoteB.BlockHash {
		return fmt.Errorf("DuplicatePrepareVoteEvidence BlockHash is equal, VoteA:%s, VoteB:%s", d.VoteA.BlockHash, d.VoteB.BlockHash)
	}
	// Verify consensus msg signature
	if err := d.VoteA.Verify(); err != nil {
		return fmt.Errorf("DuplicatePrepareVoteEvidence voteA verify failed")
	}
	if err := d.VoteB.Verify(); err != nil {
		return fmt.Errorf("DuplicatePrepareVoteEvidence voteB verify failed")
	}
	return nil
}

func (d DuplicatePrepareVoteEvidence) Address() common.Address {
	return d.VoteA.ValidateNode.Address
}

func (d DuplicatePrepareVoteEvidence) Type() consensus.EvidenceType {
	return DuplicatePrepareVoteType
}

type DuplicateViewChangeEvidence struct {
	ViewA *EvidenceView
	ViewB *EvidenceView
}

func (d DuplicateViewChangeEvidence) BlockNumber() uint64 {
	return d.ViewA.BlockNumber
}

func (d DuplicateViewChangeEvidence) Epoch() uint64 {
	return d.ViewA.Epoch
}

func (d DuplicateViewChangeEvidence) ViewNumber() uint64 {
	return d.ViewA.ViewNumber
}

func (d DuplicateViewChangeEvidence) Hash() []byte {
	var buf []byte
	if ac, err := d.ViewA.CannibalizeBytes(); err == nil {
		if bc, err := d.ViewB.CannibalizeBytes(); err == nil {
			buf, err = rlp.EncodeToBytes([]interface{}{
				ac,
				d.ViewA.Signature.Bytes(),
				bc,
				d.ViewB.Signature.Bytes(),
			})
		}
	}
	return crypto.Keccak256(buf)
}

func (d DuplicateViewChangeEvidence) Equal(ev consensus.Evidence) bool {
	_, ok := ev.(*DuplicateViewChangeEvidence)
	if !ok {
		return false
	}
	dh := d.Hash()
	eh := ev.Hash()
	return bytes.Equal(dh, eh)
}

func (d DuplicateViewChangeEvidence) Error() string {
	return fmt.Sprintf("DuplicateViewChangeEvidence epoch:%d, viewNumber:%d, blockNumber:%d blockHashA:%s, blockHashB:%s",
		d.ViewA.Epoch, d.ViewA.ViewNumber, d.ViewA.BlockNumber, d.ViewA.BlockHash.String(), d.ViewB.BlockHash.String())
}

func (d DuplicateViewChangeEvidence) Validate() error {
	if d.ViewA.Epoch != d.ViewB.Epoch {
		return fmt.Errorf("DuplicateViewChangeEvidence Epoch is different, ViewA:%d, ViewB:%d", d.ViewA.Epoch, d.ViewB.Epoch)
	}
	if d.ViewA.ViewNumber != d.ViewB.ViewNumber {
		return fmt.Errorf("DuplicateViewChangeEvidence ViewNumber is different, ViewA:%d, ViewB:%d", d.ViewA.ViewNumber, d.ViewB.ViewNumber)
	}
	if d.ViewA.BlockNumber != d.ViewB.BlockNumber {
		return fmt.Errorf("DuplicateViewChangeEvidence BlockNumber is different, ViewA:%d, ViewB:%d", d.ViewA.BlockNumber, d.ViewB.BlockNumber)
	}
	validateNodeA, validateNodeB := d.ViewA.ValidateNode, d.ViewB.ValidateNode
	if validateNodeA.Index != validateNodeB.Index || validateNodeA.Address != validateNodeB.Address {
		return fmt.Errorf("DuplicateViewChangeEvidence Validator do not match, ViewA:%s, ViewB:%s", validateNodeA.Address, validateNodeB.Address)
	}
	if d.ViewA.BlockHash == d.ViewB.BlockHash {
		return fmt.Errorf("DuplicateViewChangeEvidence BlockHash is equal, ViewA:%s, ViewB:%s", d.ViewA.BlockHash, d.ViewB.BlockHash)
	}
	// Verify consensus msg signature
	if err := d.ViewA.Verify(); err != nil {
		return fmt.Errorf("DuplicateViewChangeEvidence ViewA verify failed")
	}
	if err := d.ViewB.Verify(); err != nil {
		return fmt.Errorf("DuplicateViewChangeEvidence ViewB verify failed")
	}
	return nil
}

func (d DuplicateViewChangeEvidence) Address() common.Address {
	return d.ViewA.ValidateNode.Address
}

func (d DuplicateViewChangeEvidence) Type() consensus.EvidenceType {
	return DuplicateViewChangeType
}

type EvidenceData struct {
	DP []*DuplicatePrepareBlockEvidence `json:"duplicate_prepare"`
	DV []*DuplicatePrepareVoteEvidence  `json:"duplicate_vote"`
	DC []*DuplicateViewChangeEvidence   `json:"duplicate_viewchange"`
}

func NewEvidenceData() *EvidenceData {
	return &EvidenceData{
		DP: make([]*DuplicatePrepareBlockEvidence, 0),
		DV: make([]*DuplicatePrepareVoteEvidence, 0),
		DC: make([]*DuplicateViewChangeEvidence, 0),
	}
}

func ClassifyEvidence(evds consensus.Evidences) *EvidenceData {
	ed := NewEvidenceData()
	for _, e := range evds {
		switch e.(type) {
		case *DuplicatePrepareBlockEvidence:
			ed.DP = append(ed.DP, e.(*DuplicatePrepareBlockEvidence))
		case *DuplicatePrepareVoteEvidence:
			ed.DV = append(ed.DV, e.(*DuplicatePrepareVoteEvidence))
		case *DuplicateViewChangeEvidence:
			ed.DC = append(ed.DC, e.(*DuplicateViewChangeEvidence))
		}
	}
	return ed
}
