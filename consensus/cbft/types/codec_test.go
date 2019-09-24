package types

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/consensus/cbft/utils"
)

func TestCodec(t *testing.T) {
	// EncodeExtra
	cbftVersion := 1
	qc := &QuorumCert{
		Epoch:        1,
		ViewNumber:   0,
		BlockHash:    common.BytesToHash(utils.Rand32Bytes(32)),
		BlockNumber:  1,
		BlockIndex:   0,
		Signature:    Signature{},
		ValidatorSet: utils.NewBitArray(25),
	}
	data, err := EncodeExtra(byte(cbftVersion), qc)
	assert.Nil(t, err)
	assert.True(t, len(data) > 0)

	// DecodeExtra
	version, cert, err := DecodeExtra(data)
	assert.Nil(t, err)
	assert.Equal(t, byte(cbftVersion), version)
	assert.Equal(t, qc.BlockHash, cert.BlockHash)
}
