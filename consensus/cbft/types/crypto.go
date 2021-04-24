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

package types

import (
	"fmt"
	"reflect"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

const (
	SignatureLength = 64
)

type Signature [SignatureLength]byte

func (sig *Signature) String() string {
	return fmt.Sprintf("%x", sig[:])
}

func (sig *Signature) SetBytes(signSlice []byte) {
	copy(sig[:], signSlice[:])
}

func (sig *Signature) Bytes() []byte {
	target := make([]byte, len(sig))
	copy(target[:], sig[:])
	return target
}

// MarshalText returns the hex representation of a.
func (sig Signature) MarshalText() ([]byte, error) {
	return hexutil.Bytes(sig[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (sig *Signature) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockConfirmSign", input, sig[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (sig *Signature) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(reflect.TypeOf(Signature{}), input, sig[:])
}

func BytesToSignature(signSlice []byte) Signature {
	var sign Signature
	copy(sign[:], signSlice[:])
	return sign
}

type QuorumCert struct {
	Epoch        uint64          `json:"epoch"`
	ViewNumber   uint64          `json:"viewNumber"`
	BlockHash    common.Hash     `json:"blockHash"`
	BlockNumber  uint64          `json:"blockNumber"`
	BlockIndex   uint32          `json:"blockIndex"`
	Signature    Signature       `json:"signature"`
	ValidatorSet *utils.BitArray `json:"validatorSet"`
}

func (q QuorumCert) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		q.Epoch,
		q.ViewNumber,
		q.BlockHash,
		q.BlockNumber,
		q.BlockIndex,
	})
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (q QuorumCert) Len() int {
	length := 0
	for i := uint32(0); i < q.ValidatorSet.Size(); i++ {
		if q.ValidatorSet.GetIndex(i) {
			length++
		}
	}
	return length
}

func (q *QuorumCert) String() string {
	if q == nil {
		return ""
	}
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d,Index:%d,ValidatorSet:%s}", q.Epoch, q.ViewNumber, q.BlockHash.TerminalString(), q.BlockNumber, q.BlockIndex, q.ValidatorSet.String())
}

// if the two quorumCert have the same blockNumber
func (q *QuorumCert) HigherBlockView(blockEpoch, blockView uint64) bool {
	return q.Epoch > blockEpoch || (q.Epoch == blockEpoch && q.ViewNumber > blockView)
}

func (q *QuorumCert) HigherQuorumCert(blockNumber uint64, blockEpoch, blockView uint64) bool {
	if q.BlockNumber > blockNumber {
		return true
	} else if q.BlockNumber == blockNumber {
		return q.HigherBlockView(blockEpoch, blockView)
	}
	return false
}

type ViewChangeQuorumCert struct {
	Epoch           uint64          `json:"epoch"`
	ViewNumber      uint64          `json:"viewNumber"`
	BlockHash       common.Hash     `json:"blockHash"`
	BlockNumber     uint64          `json:"blockNumber"`
	BlockEpoch      uint64          `json:"blockEpoch"`
	BlockViewNumber uint64          `json:"blockViewNumber"`
	Signature       Signature       `json:"signature"`
	ValidatorSet    *utils.BitArray `json:"validatorSet"`
}

func (q ViewChangeQuorumCert) CannibalizeBytes() ([]byte, error) {
	buf, err := rlp.EncodeToBytes([]interface{}{
		q.Epoch,
		q.ViewNumber,
		q.BlockHash,
		q.BlockNumber,
		q.BlockEpoch,
		q.BlockViewNumber,
	})
	if err != nil {
		return nil, err
	}
	return crypto.Keccak256(buf), nil
}

func (q ViewChangeQuorumCert) Len() int {
	length := 0
	for i := uint32(0); i < q.ValidatorSet.Size(); i++ {
		if q.ValidatorSet.GetIndex(i) {
			length++
		}
	}
	return length
}

func (q ViewChangeQuorumCert) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d,BlockEpoch:%d,BlockViewNumber:%d:ValidatorSet:%s}", q.Epoch, q.ViewNumber, q.BlockHash.TerminalString(), q.BlockNumber, q.BlockEpoch, q.BlockViewNumber, q.ValidatorSet.String())
}

// if the two quorumCert have the same blockNumber
func (q *ViewChangeQuorumCert) HigherBlockView(blockEpoch, blockView uint64) bool {
	return q.BlockEpoch > blockEpoch || (q.BlockEpoch == blockEpoch && q.BlockViewNumber > blockView)
}

func (q *ViewChangeQuorumCert) HigherQuorumCert(c *ViewChangeQuorumCert) bool {
	if q.BlockNumber > c.BlockNumber {
		return true
	} else if q.BlockNumber == c.BlockNumber {
		return q.HigherBlockView(c.BlockEpoch, c.BlockViewNumber)
	}
	return false
}

func (q *ViewChangeQuorumCert) Copy() *ViewChangeQuorumCert {
	return &ViewChangeQuorumCert{
		Epoch:        q.Epoch,
		ViewNumber:   q.ViewNumber,
		BlockHash:    q.BlockHash,
		BlockNumber:  q.BlockNumber,
		Signature:    q.Signature,
		ValidatorSet: q.ValidatorSet.Copy(),
	}
}

func (v ViewChangeQC) EqualAll(epoch uint64, viewNumber uint64) error {
	for _, v := range v.QCs {
		if v.ViewNumber != viewNumber || v.Epoch != epoch {
			return fmt.Errorf("not equal, local:{%d}, want{%d}", v.ViewNumber, viewNumber)
		}
	}
	return nil
}

type ViewChangeQC struct {
	QCs []*ViewChangeQuorumCert `json:"qcs"`
}

func (v ViewChangeQC) MaxBlock() (uint64, uint64, uint64, uint64, common.Hash, uint64) {
	if len(v.QCs) == 0 {
		return 0, 0, 0, 0, common.Hash{}, 0
	}

	maxQC := v.QCs[0]
	for _, qc := range v.QCs {
		if qc.HigherQuorumCert(maxQC) {
			maxQC = qc
		}
	}
	return maxQC.Epoch, maxQC.ViewNumber, maxQC.BlockEpoch, maxQC.BlockViewNumber, maxQC.BlockHash, maxQC.BlockNumber
}

func (v ViewChangeQC) Len() int {
	length := 0
	for _, qc := range v.QCs {
		length += qc.Len()
	}
	return length
}

func (v ViewChangeQC) String() string {
	epoch, view, blockEpoch, blockViewNumber, hash, number := v.MaxBlock()
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,BlockEpoch:%d,BlockViewNumber:%d,Hash:%s,Number:%d}", epoch, view, blockEpoch, blockViewNumber, hash.TerminalString(), number)
}

func (v ViewChangeQC) ExistViewChange(epoch, viewNumber uint64, blockHash common.Hash) bool {
	for _, vc := range v.QCs {
		if vc.Epoch == epoch && vc.ViewNumber == viewNumber && vc.BlockHash == blockHash {
			return true
		}
	}
	return false
}

func (v *ViewChangeQC) AppendQuorumCert(viewChangeQC *ViewChangeQuorumCert) {
	v.QCs = append(v.QCs, viewChangeQC)
}
