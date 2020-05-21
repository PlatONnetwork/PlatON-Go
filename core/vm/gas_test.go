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

package vm

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/params"
)

func TestCallGas(t *testing.T) {

	gasTable := params.GasTableHomestead
	_, err := callGas(gasTable, 100, 2, new(big.Int).SetUint64(10))
	assert.Nil(t, err)

	gasTable = params.GasTableConstantinople
	_, err = callGas(gasTable, 100, 2, new(big.Int).SetUint64(1000000000000000))
	assert.Nil(t, err)
}
