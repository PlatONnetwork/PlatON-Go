// Copyright 2021 The PlatON Network Authors
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

package xcom

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func TestDefaultEMConfig(t *testing.T) {
	t.Run("DefaultMainNet", func(t *testing.T) {
		if getDefaultEMConfig(DefaultMainNet) == nil {
			t.Error("DefaultMainNet can't be nil config")
		}
		if err := CheckEconomicModel(params.FORKVERSION_1_3_0); nil != err {
			t.Error(err)
		}
	})
	t.Run("DefaultTestNet", func(t *testing.T) {
		if getDefaultEMConfig(DefaultTestNet) == nil {
			t.Error("DefaultTestNet can't be nil config")
		}
		if err := CheckEconomicModel(params.FORKVERSION_1_3_0); nil != err {
			t.Error(err)
		}
	})
	t.Run("DefaultUnitTestNet", func(t *testing.T) {
		if getDefaultEMConfig(DefaultUnitTestNet) == nil {
			t.Error("DefaultUnitTestNet can't be nil config")
		}
		if err := CheckEconomicModel(params.FORKVERSION_1_3_0); nil != err {
			t.Error(err)
		}
	})
	if getDefaultEMConfig(10) != nil {
		t.Error("the chain config not support")
	}
}

func TestMainNetHash(t *testing.T) {
	tempEc := getDefaultEMConfig(DefaultMainNet)
	bytes, err := rlp.EncodeToBytes(tempEc)
	if err != nil {
		t.Error(err)
	}
	assert.True(t, common.RlpHash(bytes).Hex() == MainNetECHash)
}
