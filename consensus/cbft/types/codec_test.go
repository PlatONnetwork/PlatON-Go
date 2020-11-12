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
