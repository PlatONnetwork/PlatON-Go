package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func makeViewChangeQuorumCert(epoch, viewNumber uint64, blockHash common.Hash, blockNumber uint64, blockEpoch, blockViewNumber uint64) *ViewChangeQuorumCert {
	return &ViewChangeQuorumCert{
		Epoch:           epoch,
		ViewNumber:      viewNumber,
		BlockHash:       blockHash,
		BlockNumber:     blockNumber,
		BlockEpoch:      blockEpoch,
		BlockViewNumber: blockViewNumber,
	}
}

func TestViewChangeQC_MaxBlock(t *testing.T) {
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
}
