package types

import (
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/common"
)

type Signature struct {
}

func (s *Signature) Bytes() []byte {
	return nil
}

func (s *Signature) SetBytes(buf []byte) {
}

type QuorumCert struct {
	Epoch       uint64      `json:"epoch"`
	ViewNumber  uint64      `json:"view_number"`
	BlockHash   common.Hash `json:"block_hash"`
	BlockNumber uint64      `json:"block_number"`
	BlockIndex  uint32      `json:"block_index"`
	Signature   Signature   `json:"signature"`
}

type ViewChangeQuorumCert struct {
	Epoch       uint64      `json:"epoch"`
	ViewNumber  uint64      `json:"view_number"`
	BlockHash   common.Hash `json:"block_hash"`
	BlockNumber uint64      `json:"block_number"`
	Signature   Signature   `json:"signature"`
	PrepareQC   *QuorumCert `json:"prepare_qc"`
}

type ViewChangeQC struct {
	QCs []*ViewChangeQuorumCert
}

func (v ViewChangeQC) Verify() error {
	//todo implement verify
	return nil
}

func (v ViewChangeQC) MaxBlock() (uint64, uint64, common.Hash, uint64) {
	if len(v.QCs) == 0 {
		return 0, 0, common.Hash{}, 0
	}
	epoch, view, hash, number := v.QCs[0].Epoch, v.QCs[0].ViewNumber, v.QCs[0].BlockHash, v.QCs[0].BlockNumber

	for _, qc := range v.QCs {
		if view < qc.ViewNumber {
			epoch, view, hash, number = qc.Epoch, qc.ViewNumber, qc.BlockHash, qc.BlockNumber
		} else if view == qc.ViewNumber && number < qc.BlockNumber {
			hash, number = qc.BlockHash, qc.BlockNumber
		}
	}
	return epoch, view, hash, number
}

func (v ViewChangeQC) String() string {
	epoch, view, hash, number := v.MaxBlock()
	return fmt.Sprintf("{Epoch:%d,ViewNumber:%d,Hash:%s,Number:%d}", epoch, view, hash, number)
}
