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
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func Test_QuorumCert(t *testing.T) {
	qc := &QuorumCert{
		Epoch:        1,
		ViewNumber:   1,
		BlockHash:    common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:  2,
		BlockIndex:   0,
		Signature:    BytesToSignature(utils.Rand32Bytes(64)),
		ValidatorSet: utils.NewBitArray(25),
	}
	qc.ValidatorSet.SetIndex(0, true)
	qc.ValidatorSet.SetIndex(24, true)
	_, err := qc.CannibalizeBytes()
	assert.Nil(t, err)
	assert.Equal(t, 2, qc.Len())
	assert.NotEmpty(t, qc.String())
	assert.True(t, qc.HigherQuorumCert(1, 1, 0))
	assert.True(t, qc.HigherQuorumCert(2, 1, 0))
}

func Test_ViewChangeQC(t *testing.T) {
	viewChangeQC := new(ViewChangeQC)
	hash1 := common.BytesToHash(utils.Rand32Bytes(32))
	hash2 := common.BytesToHash(utils.Rand32Bytes(32))
	hash3 := common.BytesToHash(utils.Rand32Bytes(32))
	viewChangeQC.AppendQuorumCert(makeViewChangeQuorumCert(2, 3, hash1, 9, 2, 2))
	viewChangeQC.AppendQuorumCert(makeViewChangeQuorumCert(2, 3, hash2, 9, 2, 3))
	assert.True(t, viewChangeQC.ExistViewChange(2, 3, hash2))
	assert.NotEmpty(t, viewChangeQC.String())
	assert.Equal(t, 2, viewChangeQC.Len())
	viewChangeQC.AppendQuorumCert(makeViewChangeQuorumCert(2, 4, hash3, 9, 2, 3))
	assert.NotNil(t, viewChangeQC.EqualAll(2, 3))
	last := viewChangeQC.QCs[len(viewChangeQC.QCs)-1]
	copy := last.Copy()
	assert.NotEmpty(t, copy.String())
	_, err := copy.CannibalizeBytes()
	assert.Nil(t, err)
	assert.Equal(t, hash3, copy.BlockHash)
}

func makeViewChangeQuorumCert(epoch, viewNumber uint64, blockHash common.Hash, blockNumber uint64, blockEpoch, blockViewNumber uint64) *ViewChangeQuorumCert {
	cert := &ViewChangeQuorumCert{
		Epoch:           epoch,
		ViewNumber:      viewNumber,
		BlockHash:       blockHash,
		BlockNumber:     blockNumber,
		BlockEpoch:      blockEpoch,
		BlockViewNumber: blockViewNumber,
		Signature:       BytesToSignature(utils.Rand32Bytes(64)),
		ValidatorSet:    utils.NewBitArray(25),
	}
	cert.ValidatorSet.SetIndex(0, true)
	return cert
}

func Test_ViewChangeQC_MaxBlock(t *testing.T) {
	certs := []*ViewChangeQuorumCert{
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 9, 2, 1),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 9, 2, 3),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 10, 2, 1),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 10, 2, 1),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 10, 2, 2),
		makeViewChangeQuorumCert(2, 3, common.BytesToHash(utils.Rand32Bytes(32)), 10, 1, 25),
	}
	viewChangeQC := &ViewChangeQC{
		QCs: certs,
	}
	epoch, viewNumber, blockEpoch, blockViewNumber, blockHash, blockNumber := viewChangeQC.MaxBlock()
	assert.Equal(t, certs[4].Epoch, epoch)
	assert.Equal(t, certs[4].ViewNumber, viewNumber)
	assert.Equal(t, certs[4].BlockEpoch, blockEpoch)
	assert.Equal(t, certs[4].BlockViewNumber, blockViewNumber)
	assert.Equal(t, certs[4].BlockHash, blockHash)
	assert.Equal(t, certs[4].BlockNumber, blockNumber)

	viewChangeQC.QCs = nil
	epoch, viewNumber, blockEpoch, blockViewNumber, blockHash, blockNumber = viewChangeQC.MaxBlock()
	assert.Equal(t, uint64(0), epoch)
}
