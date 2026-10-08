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

package utils

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// newBitmap returns a BitArray of the given size with every index in set enabled.
func newBitmap(bits uint32, set ...uint32) *BitArray {
	bA := NewBitArray(bits)
	for _, i := range set {
		bA.SetIndex(i, true)
	}
	return bA
}

// referenceOr is an independent implementation of the semantics documented on
// Or: the result has max(a, o) bits and each bit is a[i] || o[i], with the
// shorter operand right-padded by zeroes. It never calls Or, so it can serve as
// an oracle for the real implementation.
func referenceOr(a, o *BitArray) *BitArray {
	bits := MaxUInt(a.Bits, o.Bits)
	want := NewBitArray(bits)
	for i := uint32(0); i < bits; i++ {
		if a.GetIndex(i) || o.GetIndex(i) {
			want.SetIndex(i, true)
		}
	}
	return want
}

// TestOrDifferentBitmapLengths covers Or when the two operands do not carry the
// same number of bits, which is what the ValidatorSets of two sub-certificates
// of a ViewChangeQC look like when they were built for differently sized
// validator sets. Or must right-pad the shorter operand with zeroes, return the
// larger of the two sizes, be commutative, and leave both operands untouched.
func TestOrDifferentBitmapLengths(t *testing.T) {
	cases := []struct {
		name  string
		aBits uint32
		oBits uint32
	}{
		{"same length", 65, 65},
		{"receiver longer within one word", 130, 64},
		{"receiver longer across words", 130, 65},
		{"receiver shorter crossing a word boundary", 63, 65},
		{"receiver one bit shorter", 64, 65},
		{"receiver much shorter", 1, 256},
		{"receiver much longer", 256, 1},
		{"exactly one word each", 64, 64},
		{"one word versus two", 64, 128},
		{"two words versus three", 128, 192},
		{"validator sets of 4 versus 200", 4, 200},
		{"validator sets of 200 versus 4", 200, 4},
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			a := newBitmap(tc.aBits, 0, tc.aBits/2, tc.aBits-1)
			o := newBitmap(tc.oBits, 1, tc.oBits/2, tc.oBits-1)
			want := referenceOr(a, o)

			forward := a.Or(o)
			require.Equal(t, want.Bits, forward.Bits, "result must carry the max of the two sizes")
			require.Equal(t, want.Elems, forward.Elems, "a.Or(o) must OR in every word of o")

			reverse := o.Or(a)
			require.Equal(t, want.Bits, reverse.Bits, "result must carry the max of the two sizes")
			require.Equal(t, want.Elems, reverse.Elems, "Or must be commutative")

			// Neither operand may be modified.
			require.Equal(t, newBitmap(tc.aBits, 0, tc.aBits/2, tc.aBits-1).Elems, a.Elems,
				"Or must not mutate the receiver")
			require.Equal(t, newBitmap(tc.oBits, 1, tc.oBits/2, tc.oBits-1).Elems, o.Elems,
				"Or must not mutate the argument")
		})
	}
}

// TestOrKeepsTrailingWordsOfLongerArgument pins the case where the receiver
// spans fewer 64-bit words than the argument: the trailing words of the
// argument belong to the result and must not be dropped.
func TestOrKeepsTrailingWordsOfLongerArgument(t *testing.T) {
	short := newBitmap(63, 62, 5) // one word
	long := newBitmap(65, 64, 5)  // two words; bit 64 lives in word 1

	got := short.Or(long)

	require.Equal(t, uint32(65), got.Bits, "result must be as large as the larger operand")
	require.True(t, got.GetIndex(64), "bit in the trailing word of the longer argument must survive Or")
	require.True(t, got.GetIndex(62), "receiver bit must survive Or")
	require.True(t, got.GetIndex(5), "overlapping bit must survive Or")
	require.False(t, got.GetIndex(50), "unset bit must stay unset")
}

// TestOrValidatorBitmapLengths mirrors how countDistinctViewChangeSigners unions
// sub-certificate ValidatorSets: a smaller set accumulated first, then a larger
// one, and the other way round.
func TestOrValidatorBitmapLengths(t *testing.T) {
	small := newBitmap(4, 2)     // validator 2 of 4
	large := newBitmap(200, 100) // validator 100 of 200

	countSetBits := func(bA *BitArray) int {
		count := 0
		for i := uint32(0); i < bA.Size(); i++ {
			if bA.GetIndex(i) {
				count++
			}
		}
		return count
	}

	accumulated := small.Or(large)
	require.Equal(t, uint32(200), accumulated.Bits)
	require.Equal(t, 2, countSetBits(accumulated),
		"union must count both validators regardless of the order they are OR-ed in")

	reversed := large.Or(small)
	require.Equal(t, uint32(200), reversed.Bits)
	require.Equal(t, 2, countSetBits(reversed),
		"union must count both validators regardless of the order they are OR-ed in")

	require.Equal(t, accumulated.Elems, reversed.Elems, "union must not depend on accumulation order")
}

// TestOrWithNilAcrossLengths covers the nil degenerate cases for a
// multi-word operand, where Or falls back to copying the non-nil side.
func TestOrWithNilAcrossLengths(t *testing.T) {
	var nilBA *BitArray
	bA := newBitmap(130, 64, 129)

	require.Equal(t, bA, nilBA.Or(bA), "nil receiver must return a copy of the argument")
	require.Equal(t, bA, bA.Or(nilBA), "nil argument must return a copy of the receiver")
	require.Equal(t, (*BitArray)(nil), nilBA.Or(nilBA), "two nil operands must give nil")
}
