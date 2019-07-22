package types

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
)

type Signature struct {
}

func (s *Signature) Bytes() []byte {
	return nil
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

func (v ViewChangeQC) MaxBlock() (uint64, common.Hash, uint64) {
	if len(v.QCs) == 0 {
		return 0, common.Hash{}, 0
	}
	view, hash, number := v.QCs[0].ViewNumber, v.QCs[0].BlockHash, v.QCs[0].BlockNumber

	for _, qc := range v.QCs {
		if view < qc.ViewNumber {
			view, hash, number = qc.ViewNumber, qc.BlockHash, qc.BlockNumber
		} else if view == qc.ViewNumber && number < qc.BlockNumber {
			hash, number = qc.BlockHash, qc.BlockNumber
		}
	}
	return view, hash, number
}
