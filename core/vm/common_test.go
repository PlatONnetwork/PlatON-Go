// Copyright 2018-2019 The PlatON Network Authors
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
)

func buildBigInt(v int64) *big.Int {
	return new(big.Int).SetInt64(v)
}

func TestCalcMemSize(t *testing.T) {
	testCases := []struct {
		off  *big.Int
		l    *big.Int
		want *big.Int
	}{
		{off: buildBigInt(10), l: buildBigInt(10), want: buildBigInt(20)},
		{off: buildBigInt(0), l: buildBigInt(0), want: buildBigInt(0)},
		{off: buildBigInt(1), l: buildBigInt(-1), want: buildBigInt(0)},
		{off: buildBigInt(10), l: buildBigInt(-10), want: buildBigInt(0)},
		{off: buildBigInt(3), l: buildBigInt(4), want: buildBigInt(7)},
	}
	for _, v := range testCases {
		res := calcMemSize(v.off, v.l)
		assert.Equal(t, v.want.Uint64(), res.Uint64())
	}
}
