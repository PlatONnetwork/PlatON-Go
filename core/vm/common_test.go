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

	"github.com/PlatONnetwork/PlatON-Go/common/math"

	"github.com/PlatONnetwork/PlatON-Go/common"

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

func TestGetData(t *testing.T) {
	testCases := []struct {
		b     []byte
		start uint64
		want  string
	}{
		{b: []byte{0x01, 0x02, 0x03}, start: 2, want: "03000000000000000000"},
		{b: []byte{0x01, 0x02, 0x02}, start: 2, want: "02000000000000000000"},
		{b: []byte{0x01, 0x02, 0x04}, start: 2, want: "04000000000000000000"},
		{b: []byte{0x01, 0x02, 0x03}, start: 1, want: "02030000000000000000"},
		{b: []byte{0x01, 0x02, 0x03}, start: 0, want: "01020300000000000000"},
	}
	for _, v := range testCases {
		r := getData(v.b, v.start, 10)
		assert.Equal(t, v.want, common.Bytes2Hex(r))
	}
}

func TestBigUint64(t *testing.T) {
	v1, _ := new(big.Int).SetString("1000000000000000000000", 10)
	v2, _ := new(big.Int).SetString("1000000000000000000000", 10)
	testCases := []struct {
		v     *big.Int
		wantv uint64
		wantr bool
	}{
		{v: buildBigInt(10), wantv: 10, wantr: false},
		{v: v1, wantv: v2.Uint64(), wantr: true},
		{v: buildBigInt(0), wantv: 0, wantr: false},
	}
	for _, v := range testCases {
		resv, resr := bigUint64(v.v)
		assert.Equal(t, v.wantv, resv)
		assert.Equal(t, v.wantr, resr)
	}
}

func TestToWordSize(t *testing.T) {
	testCases := []struct {
		v      uint64
		expect uint64
	}{
		{100, 4},
		{10000, 313},
		{math.MaxUint64 - 1, 576460752303423488},
		{math.MaxUint64 - 2, 576460752303423488},
	}
	for _, v := range testCases {
		resv := toWordSize(v.v)
		assert.Equal(t, v.expect, resv)
	}
}

func TestAllZero(t *testing.T) {
	testCases := []struct {
		v      []byte
		expect bool
	}{
		{[]byte{0x00, 0x00}, true},
		{[]byte{0x00, 0x01}, false},
	}
	for _, v := range testCases {
		resv := allZero(v.v)
		assert.Equal(t, v.expect, resv)
	}
}
