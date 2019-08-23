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
	SignatureLength = 32
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
	ViewNumber   uint64          `json:"view_number"`
	BlockHash    common.Hash     `json:"block_hash"`
	BlockNumber  uint64          `json:"block_number"`
	BlockIndex   uint32          `json:"block_index"`
	Signature    Signature       `json:"signature"`
	ValidatorSet *utils.BitArray `json:"validator_set"`
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

func (q QuorumCert) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d,Index:%d}", q.Epoch, q.ViewNumber, q.BlockHash.TerminalString(), q.BlockNumber, q.BlockIndex)
}

type ViewChangeQuorumCert struct {
	Epoch           uint64          `json:"epoch"`
	ViewNumber      uint64          `json:"view_number"`
	BlockHash       common.Hash     `json:"block_hash"`
	BlockNumber     uint64          `json:"block_number"`
	BlockEpoch      uint64          `json:"block_epoch"`
	BlockViewNumber uint64          `json:"block_view_number"`
	Signature       Signature       `json:"signature"`
	ValidatorSet    *utils.BitArray `json:"validator_set"`
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
	for i := uint32(0); i < q.ValidatorSet.Bits; i++ {
		if q.ValidatorSet.GetIndex(i) {
			length++
		}
	}
	return length
}

func (q ViewChangeQuorumCert) String() string {
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d,BlockEpoch:%d,BlockViewNumber:%d}", q.Epoch, q.ViewNumber, q.BlockHash.TerminalString(), q.BlockNumber, q.BlockEpoch, q.BlockViewNumber)
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
	epoch, view, blockEpoch, blockView, hash, number := v.QCs[0].Epoch, v.QCs[0].ViewNumber, v.QCs[0].BlockEpoch, v.QCs[0].BlockViewNumber, v.QCs[0].BlockHash, v.QCs[0].BlockNumber

	for _, qc := range v.QCs {
		if blockEpoch < qc.BlockEpoch || blockView < qc.BlockViewNumber {
			blockEpoch, blockView, hash, number = qc.BlockEpoch, qc.BlockViewNumber, qc.BlockHash, qc.BlockNumber
		} else if number < qc.BlockNumber {
			hash, number = qc.BlockHash, qc.BlockNumber
		}
	}
	return epoch, view, blockEpoch, blockView, hash, number
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
