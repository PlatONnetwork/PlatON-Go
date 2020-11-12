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

package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func randBitArray(bits uint32) (*BitArray, []byte) {
	src := Rand32Bytes((bits + 7) / 8)
	bA := NewBitArray(bits)
	for i := uint32(0); i < uint32(len(src)); i++ {
		for j := uint32(0); j < 8; j++ {
			if i*8+j >= bits {
				return bA, src
			}
			setBit := src[i]&(1<<uint(j)) > 0
			bA.SetIndex(i*8+j, setBit)
		}
	}
	return bA, src
}

func TestRLP(t *testing.T) {
	bA1, _ := randBitArray(51)
	if buf, err := rlp.EncodeToBytes(bA1); err != nil {
		t.Error(err)
	} else {
		var bA2 BitArray
		if err := rlp.DecodeBytes(buf, &bA2); err != nil {
			assert.Error(t, err)
		}
		require.Equal(t, bA1, &bA2)
	}

}

func TestAnd(t *testing.T) {

	bA1, _ := randBitArray(51)
	bA2, _ := randBitArray(31)
	bA3 := bA1.And(bA2)

	var bNil *BitArray
	require.Equal(t, bNil.And(bA1), (*BitArray)(nil))
	require.Equal(t, bA1.And(nil), (*BitArray)(nil))
	require.Equal(t, bNil.And(nil), (*BitArray)(nil))

	if bA3.Bits != 31 {
		t.Error("Expected min bits", bA3.Bits)
	}
	if len(bA3.Elems) != len(bA2.Elems) {
		t.Error("Expected min elems length")
	}
	for i := uint32(0); i < bA3.Bits; i++ {
		expected := bA1.GetIndex(i) && bA2.GetIndex(i)
		if bA3.GetIndex(i) != expected {
			t.Error("Wrong bit from bA3", i, bA1.GetIndex(i), bA2.GetIndex(i), bA3.GetIndex(i))
		}
	}
}

func TestOr(t *testing.T) {

	bA1, _ := randBitArray(51)
	bA2, _ := randBitArray(31)
	bA3 := bA1.Or(bA2)

	bNil := (*BitArray)(nil)
	require.Equal(t, bNil.Or(bA1), bA1)
	require.Equal(t, bA1.Or(nil), bA1)
	require.Equal(t, bNil.Or(nil), (*BitArray)(nil))

	if bA3.Bits != 51 {
		t.Error("Expected max bits")
	}
	if len(bA3.Elems) != len(bA1.Elems) {
		t.Error("Expected max elems length")
	}
	for i := uint32(0); i < bA3.Bits; i++ {
		expected := bA1.GetIndex(i) || bA2.GetIndex(i)
		if bA3.GetIndex(i) != expected {
			t.Error("Wrong bit from bA3", i, bA1.GetIndex(i), bA2.GetIndex(i), bA3.GetIndex(i))
		}
	}
}

func TestSub(t *testing.T) {
	testCases := []struct {
		initBA        string
		subtractingBA string
		expectedBA    string
	}{
		{"null", `null`, `null`},
		{`null`, `null`, `null`},
		{`"x"`, `null`, `null`},
		{`null`, `"x"`, `null`},
		{`"x"`, `"x"`, `"_"`},
		{`"xxxxxx"`, `"x_x_x_"`, `"_x_x_x"`},
		{`"x_x_x_"`, `"xxxxxx"`, `"______"`},
		{`"xxxxxx"`, `"x_x_x_xxxx"`, `"_x_x_x"`},
		{`"x_x_x_xxxx"`, `"xxxxxx"`, `"______xxxx"`},
		{`"xxxxxxxxxx"`, `"x_x_x_"`, `"_x_x_xxxxx"`},
		{`"x_x_x_"`, `"xxxxxxxxxx"`, `"______"`},
	}
	for _, tc := range testCases {
		var bA *BitArray
		err := json.Unmarshal([]byte(tc.initBA), &bA)
		require.Nil(t, err)

		var o *BitArray
		err = json.Unmarshal([]byte(tc.subtractingBA), &o)
		require.Nil(t, err)

		got, _ := json.Marshal(bA.Sub(o))
		require.Equal(t, tc.expectedBA, string(got), "%s minus %s doesn't equal %s", tc.initBA, tc.subtractingBA, tc.expectedBA)
	}
}

func TestBytes(t *testing.T) {
	bA := NewBitArray(4)
	bA.SetIndex(0, true)
	check := func(bA *BitArray, bz []byte) {
		if !bytes.Equal(bA.Bytes(), bz) {
			panic(fmt.Sprintf("Expected %X but got %X", bz, bA.Bytes()))
		}
	}
	check(bA, []byte{0x01})
	bA.SetIndex(3, true)
	check(bA, []byte{0x09})

	bA = NewBitArray(9)
	check(bA, []byte{0x00, 0x00})
	bA.SetIndex(7, true)
	check(bA, []byte{0x80, 0x00})
	bA.SetIndex(8, true)
	check(bA, []byte{0x80, 0x01})

	bA = NewBitArray(16)
	check(bA, []byte{0x00, 0x00})
	bA.SetIndex(7, true)
	check(bA, []byte{0x80, 0x00})
	bA.SetIndex(8, true)
	check(bA, []byte{0x80, 0x01})
	bA.SetIndex(9, true)
	check(bA, []byte{0x80, 0x03})
	assert.Equal(t, uint32(16), bA.Size())
}

func TestEmptyFull(t *testing.T) {
	ns := []uint32{47, 123}
	for _, n := range ns {
		bA := NewBitArray(n)
		if !bA.IsEmpty() {
			t.Fatal("Expected bit array to be empty")
		}
		for i := uint32(0); i < n; i++ {
			bA.SetIndex(i, true)
		}
		if !bA.IsFull() {
			t.Fatal("Expected bit array to be full")
		}
	}
}

func TestUpdateNeverPanics(t *testing.T) {
	newRandBitArray := func(n uint32) *BitArray {
		ba, _ := randBitArray(n)
		return ba
	}
	pairs := []struct {
		a, b *BitArray
	}{
		{nil, nil},
		{newRandBitArray(10), newRandBitArray(12)},
		{newRandBitArray(23), newRandBitArray(23)},
		{newRandBitArray(37), nil},
		{nil, NewBitArray(10)},
	}

	for _, pair := range pairs {
		a, b := pair.a, pair.b
		a.Update(b)
		b.Update(a)
	}
}

func TestJSONMarshalUnmarshal(t *testing.T) {

	bA1 := NewBitArray(0)

	bA2 := NewBitArray(1)

	bA3 := NewBitArray(1)
	bA3.SetIndex(0, true)

	bA4 := NewBitArray(5)
	bA4.SetIndex(0, true)
	bA4.SetIndex(1, true)

	testCases := []struct {
		bA           *BitArray
		marshalledBA string
	}{
		{nil, `null`},
		{bA1, `null`},
		{bA2, `"_"`},
		{bA3, `"x"`},
		{bA4, `"xx___"`},
	}

	for _, tc := range testCases {
		t.Run(tc.bA.String(), func(t *testing.T) {
			bz, err := json.Marshal(tc.bA)
			require.NoError(t, err)

			assert.Equal(t, tc.marshalledBA, string(bz))

			var unmarshalledBA *BitArray
			err = json.Unmarshal(bz, &unmarshalledBA)
			require.NoError(t, err)

			if tc.bA == nil {
				require.Nil(t, unmarshalledBA)
			} else {
				require.NotNil(t, unmarshalledBA)
				assert.EqualValues(t, tc.bA.Bits, unmarshalledBA.Bits)
				if assert.EqualValues(t, tc.bA.String(), unmarshalledBA.String()) {
					assert.EqualValues(t, tc.bA.Elems, unmarshalledBA.Elems)
				}
			}
		})
	}
}

func TestGetIndex(t *testing.T) {
	var bit *BitArray
	isExists := bit.GetIndex(0)
	assert.True(t, !isExists)

	// for: getIndex
	bitArr := NewBitArray(3)
	isExists = bitArr.GetIndex(100)
	assert.False(t, isExists)
}

func TestNot(t *testing.T) {
	bA1, _ := randBitArray(51)
	//bA2, _ := randBitArray(31)

	bNil := (*BitArray)(nil)
	bA3 := bNil.Not()
	assert.Nil(t, bA3)

	bA4 := bA1.not()
	assert.NotEqual(t, bA1, bA4)
}

func TestPickRandom(t *testing.T) {
	bNil := (*BitArray)(nil)
	i, isExists := bNil.PickRandom()
	assert.Equal(t, i, uint32(0))
	assert.Equal(t, isExists, false)

	bA1, _ := randBitArray(51)
	i, _ = bA1.PickRandom()
	assert.NotEqual(t, i, uint32(0))
	assert.NotEqual(t, i, true)
}

func TestGetTrueIndices(t *testing.T) {
	bA1, _ := randBitArray(51)
	res := bA1.getTrueIndices()
	assert.True(t, len(res) != 0)
}

func TestMinInt64(t *testing.T) {
	testCase := []struct {
		num1 int64
		num2 int64
		want int64
	}{
		{4, 5, 4},
		{45, 30, 30},
	}
	for _, v := range testCase {
		r := MinInt64(v.num1, v.num2)
		if r != v.want {
			t.Errorf("MinInt64 error, want: %v, current:%v", v.want, r)
		}
	}
}

func TestMinInt(t *testing.T) {
	testCase := []struct {
		num1 int
		num2 int
		want int
	}{
		{-4, 5, -4},
		{45, 30, 30},
	}
	for _, v := range testCase {
		r := MinInt(v.num1, v.num2)
		if r != v.want {
			t.Errorf("MinInt error, want: %v, current:%v", v.want, r)
		}
	}
	randByt := RandBytes(10)
	assert.True(t, len(randByt) == 10)
}
