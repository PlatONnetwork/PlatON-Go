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

package evidence

import (
	"bytes"
	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/crypto/bls"
	"github.com/PlatONnetwork/PlatON-Go/p2p/discover"

	"github.com/PlatONnetwork/PlatON-Go/common/consensus"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

const (
	DuplicatePrepareBlockType consensus.EvidenceType = 1
	DuplicatePrepareVoteType  consensus.EvidenceType = 2
	DuplicateViewChangeType   consensus.EvidenceType = 3
)

// DuplicatePrepareBlockEvidence recording duplicate blocks
type DuplicatePrepareBlockEvidence struct {
	PrepareA *EvidencePrepare `json:"prepareA"`
	PrepareB *EvidencePrepare `json:"prepareB"`
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
			buf, _ = rlp.EncodeToBytes([]interface{}{
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
	return fmt.Sprintf("DuplicatePrepareBlockEvidence, epoch:%d, viewNumber:%d, blockNumber:%d, blockHashA:%s, blockHashB:%s",
		d.PrepareA.Epoch, d.PrepareA.ViewNumber, d.PrepareA.BlockNumber, d.PrepareA.BlockHash.String(), d.PrepareB.BlockHash.String())
}

// Validate verify the validity of the DuplicatePrepareBlockEvidence
// the same epoch,viewNumber,blockNumber,nodeId and different blockHash
func (d DuplicatePrepareBlockEvidence) Validate() error {
	if d.PrepareA.Epoch != d.PrepareB.Epoch {
		return fmt.Errorf("DuplicatePrepareBlockEvidence, epoch is different, prepareA:%d, prepareB:%d", d.PrepareA.Epoch, d.PrepareB.Epoch)
	}
	if d.PrepareA.ViewNumber != d.PrepareB.ViewNumber {
		return fmt.Errorf("DuplicatePrepareBlockEvidence, viewNumber is different, prepareA:%d, prepareB:%d", d.PrepareA.ViewNumber, d.PrepareB.ViewNumber)
	}
	if d.PrepareA.BlockNumber != d.PrepareB.BlockNumber {
		return fmt.Errorf("DuplicatePrepareBlockEvidence, blockNumber is different, prepareA:%d, prepareB:%d", d.PrepareA.BlockNumber, d.PrepareB.BlockNumber)
	}
	validateNodeA, validateNodeB := d.PrepareA.ValidateNode, d.PrepareB.ValidateNode
	if validateNodeA.Index != validateNodeB.Index || validateNodeA.NodeID != validateNodeB.NodeID ||
		!bytes.Equal(validateNodeA.BlsPubKey.Serialize(), validateNodeB.BlsPubKey.Serialize()) {
		return fmt.Errorf("DuplicatePrepareBlockEvidence, validator do not match, prepareA:%s, prepareB:%s", validateNodeA.NodeID.TerminalString(), validateNodeB.NodeID.TerminalString())
	}
	if d.PrepareA.BlockHash == d.PrepareB.BlockHash {
		return fmt.Errorf("DuplicatePrepareBlockEvidence, blockHash is equal, prepareA:%s, prepareB:%s", d.PrepareA.BlockHash.String(), d.PrepareB.BlockHash.String())
	}
	// Verify consensus msg signature
	if err := d.PrepareA.Verify(); err != nil {
		return fmt.Errorf("DuplicatePrepareBlockEvidence, prepareA verify failed")
	}
	if err := d.PrepareB.Verify(); err != nil {
		return fmt.Errorf("DuplicatePrepareBlockEvidence, prepareB verify failed")
	}
	return nil
}

func (d DuplicatePrepareBlockEvidence) NodeID() discover.NodeID {
	return d.PrepareA.ValidateNode.NodeID
}

func (d DuplicatePrepareBlockEvidence) BlsPubKey() *bls.PublicKey {
	return d.PrepareA.ValidateNode.BlsPubKey
}

func (d DuplicatePrepareBlockEvidence) Type() consensus.EvidenceType {
	return DuplicatePrepareBlockType
}

func (d DuplicatePrepareBlockEvidence) ValidateMsg() bool {
	return d.PrepareA != nil && d.PrepareA.ValidateNode != nil && d.PrepareA.ValidateNode.BlsPubKey != nil &&
		d.PrepareB != nil && d.PrepareB.ValidateNode != nil && d.PrepareB.ValidateNode.BlsPubKey != nil
}

// DuplicatePrepareVoteEvidence recording duplicate vote
type DuplicatePrepareVoteEvidence struct {
	VoteA *EvidenceVote `json:"voteA"`
	VoteB *EvidenceVote `json:"voteB"`
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
			buf, _ = rlp.EncodeToBytes([]interface{}{
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
	return fmt.Sprintf("DuplicatePrepareVoteEvidence, epoch:%d, viewNumber:%d, blockNumber:%d, blockHashA:%s, blockHashB:%s",
		d.VoteA.Epoch, d.VoteA.ViewNumber, d.VoteA.BlockNumber, d.VoteA.BlockHash.String(), d.VoteB.BlockHash.String())
}

// Validate verify the validity of the duplicatePrepareVoteEvidence
// the same epoch,viewNumber,blockNumber,nodeId and different blockHash
func (d DuplicatePrepareVoteEvidence) Validate() error {
	if d.VoteA.Epoch != d.VoteB.Epoch {
		return fmt.Errorf("DuplicatePrepareVoteEvidence, epoch is different, voteA:%d, voteB:%d", d.VoteA.Epoch, d.VoteB.Epoch)
	}
	if d.VoteA.ViewNumber != d.VoteB.ViewNumber {
		return fmt.Errorf("DuplicatePrepareVoteEvidence, viewNumber is different, voteA:%d, voteB:%d", d.VoteA.ViewNumber, d.VoteB.ViewNumber)
	}
	if d.VoteA.BlockNumber != d.VoteB.BlockNumber {
		return fmt.Errorf("DuplicatePrepareVoteEvidence, blockNumber is different, voteA:%d, voteB:%d", d.VoteA.BlockNumber, d.VoteB.BlockNumber)
	}
	validateNodeA, validateNodeB := d.VoteA.ValidateNode, d.VoteB.ValidateNode
	if validateNodeA.Index != validateNodeB.Index || validateNodeA.NodeID != validateNodeB.NodeID ||
		!bytes.Equal(validateNodeA.BlsPubKey.Serialize(), validateNodeB.BlsPubKey.Serialize()) {
		return fmt.Errorf("DuplicatePrepareVoteEvidence, validator do not match, voteA:%s, voteB:%s", validateNodeA.NodeID.TerminalString(), validateNodeB.NodeID.TerminalString())
	}
	if d.VoteA.BlockHash == d.VoteB.BlockHash {
		return fmt.Errorf("DuplicatePrepareVoteEvidence, blockHash is equal, voteA:%s, voteB:%s", d.VoteA.BlockHash.String(), d.VoteB.BlockHash.String())
	}
	// Verify consensus msg signature
	if err := d.VoteA.Verify(); err != nil {
		return fmt.Errorf("DuplicatePrepareVoteEvidence, voteA verify failed")
	}
	if err := d.VoteB.Verify(); err != nil {
		return fmt.Errorf("DuplicatePrepareVoteEvidence, voteB verify failed")
	}
	return nil
}

func (d DuplicatePrepareVoteEvidence) NodeID() discover.NodeID {
	return d.VoteA.ValidateNode.NodeID
}

func (d DuplicatePrepareVoteEvidence) BlsPubKey() *bls.PublicKey {
	return d.VoteA.ValidateNode.BlsPubKey
}

func (d DuplicatePrepareVoteEvidence) Type() consensus.EvidenceType {
	return DuplicatePrepareVoteType
}

func (d DuplicatePrepareVoteEvidence) ValidateMsg() bool {
	return d.VoteA != nil && d.VoteA.ValidateNode != nil && d.VoteA.ValidateNode.BlsPubKey != nil &&
		d.VoteB != nil && d.VoteB.ValidateNode != nil && d.VoteB.ValidateNode.BlsPubKey != nil
}

// DuplicateViewChangeEvidence recording duplicate viewChange
type DuplicateViewChangeEvidence struct {
	ViewA *EvidenceView `json:"viewA"`
	ViewB *EvidenceView `json:"viewB"`
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
			buf, _ = rlp.EncodeToBytes([]interface{}{
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
	return fmt.Sprintf("DuplicateViewChangeEvidence, epoch:%d, viewNumber:%d, blockNumber:%d, blockHashA:%s, blockHashB:%s",
		d.ViewA.Epoch, d.ViewA.ViewNumber, d.ViewA.BlockNumber, d.ViewA.BlockHash.String(), d.ViewB.BlockHash.String())
}

// Validate verify the validity of the duplicateViewChangeEvidence
// the same epoch,viewNumber,nodeId and different block
func (d DuplicateViewChangeEvidence) Validate() error {
	if d.ViewA.Epoch != d.ViewB.Epoch {
		return fmt.Errorf("DuplicateViewChangeEvidence, epoch is different, viewA:%d, viewB:%d", d.ViewA.Epoch, d.ViewB.Epoch)
	}
	if d.ViewA.ViewNumber != d.ViewB.ViewNumber {
		return fmt.Errorf("DuplicateViewChangeEvidence, viewNumber is different, viewA:%d, viewB:%d", d.ViewA.ViewNumber, d.ViewB.ViewNumber)
	}
	validateNodeA, validateNodeB := d.ViewA.ValidateNode, d.ViewB.ValidateNode
	if validateNodeA.Index != validateNodeB.Index || validateNodeA.NodeID != validateNodeB.NodeID ||
		!bytes.Equal(validateNodeA.BlsPubKey.Serialize(), validateNodeB.BlsPubKey.Serialize()) {
		return fmt.Errorf("DuplicateViewChangeEvidence, validator do not match, viewA:%s, viewB:%s", validateNodeA.NodeID.TerminalString(), validateNodeB.NodeID.TerminalString())
	}
	if d.ViewA.BlockNumber == d.ViewB.BlockNumber && d.ViewA.BlockHash == d.ViewB.BlockHash {
		return fmt.Errorf("DuplicateViewChangeEvidence, blockNumber and blockHash is equal, viewANumber:%d, viewAHash:%s, viewANumber:%d, viewBHash:%s", d.ViewA.BlockNumber, d.ViewA.BlockHash.String(), d.ViewB.BlockNumber, d.ViewB.BlockHash.String())
	}
	// Verify consensus msg signature
	if err := d.ViewA.Verify(); err != nil {
		return fmt.Errorf("DuplicateViewChangeEvidence, viewA verify failed")
	}
	if err := d.ViewB.Verify(); err != nil {
		return fmt.Errorf("DuplicateViewChangeEvidence, ViewB verify failed")
	}
	return nil
}

func (d DuplicateViewChangeEvidence) NodeID() discover.NodeID {
	return d.ViewA.ValidateNode.NodeID
}

func (d DuplicateViewChangeEvidence) BlsPubKey() *bls.PublicKey {
	return d.ViewA.ValidateNode.BlsPubKey
}

func (d DuplicateViewChangeEvidence) Type() consensus.EvidenceType {
	return DuplicateViewChangeType
}

func (d DuplicateViewChangeEvidence) ValidateMsg() bool {
	return d.ViewA != nil && d.ViewA.ValidateNode != nil && d.ViewA.ValidateNode.BlsPubKey != nil &&
		d.ViewB != nil && d.ViewB.ValidateNode != nil && d.ViewB.ValidateNode.BlsPubKey != nil
}

// EvidenceData encapsulate externally visible duplicate data
type EvidenceData struct {
	DP []*DuplicatePrepareBlockEvidence `json:"duplicatePrepare"`
	DV []*DuplicatePrepareVoteEvidence  `json:"duplicateVote"`
	DC []*DuplicateViewChangeEvidence   `json:"duplicateViewchange"`
}

func NewEvidenceData() *EvidenceData {
	return &EvidenceData{
		DP: make([]*DuplicatePrepareBlockEvidence, 0),
		DV: make([]*DuplicatePrepareVoteEvidence, 0),
		DC: make([]*DuplicateViewChangeEvidence, 0),
	}
}

// ClassifyEvidence tries to convert evidence list to evidenceData
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
