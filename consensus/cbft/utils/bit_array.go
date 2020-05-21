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
	"encoding/binary"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
)

// BitArray is a thread-unsafe implementation of a bit array.
type BitArray struct {
	Bits  uint32   `json:"bits"`  // NOTE: persisted via reflect, must be exported
	Elems []uint64 `json:"elems"` // NOTE: persisted via reflect, must be exported
}

// NewBitArray returns a new bit array.
// It returns nil if the number of bits is zero.
func NewBitArray(bits uint32) *BitArray {
	if bits <= 0 {
		return nil
	}
	return &BitArray{
		Bits:  bits,
		Elems: make([]uint64, (bits+63)/64),
	}
}

// Size returns the number of bits in the bitarray
func (bA *BitArray) Size() uint32 {
	if bA == nil {
		return 0
	}
	return bA.Bits
}

// GetIndex returns the bit at index i within the bit array.
// The behavior is undefined if i >= bA.Bits
func (bA *BitArray) GetIndex(i uint32) bool {
	if bA == nil {
		return false
	}
	return bA.getIndex(i)
}

func (bA *BitArray) getIndex(i uint32) bool {
	if i >= bA.Bits {
		return false
	}
	return bA.Elems[i/64]&(uint64(1)<<uint(i%64)) > 0
}

// SetIndex sets the bit at index i within the bit array.
// The behavior is undefined if i >= bA.Bits
func (bA *BitArray) SetIndex(i uint32, v bool) bool {
	if bA == nil {
		return false
	}
	return bA.setIndex(i, v)
}

func (bA *BitArray) setIndex(i uint32, v bool) bool {
	if i >= bA.Bits {
		return false
	}
	if v {
		bA.Elems[i/64] |= (uint64(1) << uint(i%64))
	} else {
		bA.Elems[i/64] &= ^(uint64(1) << uint(i%64))
	}
	return true
}

// Copy returns a copy of the provided bit array.
func (bA *BitArray) Copy() *BitArray {
	if bA == nil {
		return nil
	}
	return bA.copy()
}

func (bA *BitArray) copy() *BitArray {
	c := make([]uint64, len(bA.Elems))
	copy(c, bA.Elems)
	return &BitArray{
		Bits:  bA.Bits,
		Elems: c,
	}
}

func (bA *BitArray) copyBits(bits uint32) *BitArray {
	c := make([]uint64, (bits+63)/64)
	copy(c, bA.Elems)
	return &BitArray{
		Bits:  bits,
		Elems: c,
	}
}

// Or returns a bit array resulting from a bitwise OR of the two bit arrays.
// If the two bit-arrys have different lengths, Or right-pads the smaller of the two bit-arrays with zeroes.
// Thus the size of the return value is the maximum of the two provided bit arrays.
func (bA *BitArray) Or(o *BitArray) *BitArray {
	if bA == nil && o == nil {
		return nil
	}
	if bA == nil && o != nil {
		return o.Copy()
	}
	if o == nil {
		return bA.Copy()
	}
	c := bA.copyBits(MaxUInt(bA.Bits, o.Bits))
	smaller := MinInt(len(bA.Elems), len(o.Elems))
	for i := 0; i < smaller; i++ {
		c.Elems[i] |= o.Elems[i]
	}
	return c
}

// And returns a bit array resulting from a bitwise AND of the two bit arrays.
// If the two bit-arrys have different lengths, this truncates the larger of the two bit-arrays from the right.
// Thus the size of the return value is the minimum of the two provided bit arrays.
func (bA *BitArray) And(o *BitArray) *BitArray {
	if bA == nil || o == nil {
		return nil
	}
	return bA.and(o)
}

func (bA *BitArray) and(o *BitArray) *BitArray {
	c := bA.copyBits(MinUInt(bA.Bits, o.Bits))
	for i := 0; i < len(c.Elems); i++ {
		c.Elems[i] &= o.Elems[i]
	}
	return c
}

// Not returns a bit array resulting from a bitwise Not of the provided bit array.
func (bA *BitArray) Not() *BitArray {
	if bA == nil {
		return nil // Degenerate
	}
	return bA.not()
}

func (bA *BitArray) not() *BitArray {
	c := bA.copy()
	for i := 0; i < len(c.Elems); i++ {
		c.Elems[i] = ^c.Elems[i]
	}
	return c
}

// Sub subtracts the two bit-arrays bitwise, without carrying the bits.
// Note that carryless subtraction of a - b is (a and not b).
// The output is the same as bA, regardless of o's size.
// If bA is longer than o, o is right padded with zeroes
func (bA *BitArray) Sub(o *BitArray) *BitArray {
	if bA == nil || o == nil {
		// TODO: Decide if we should do 1's complement here?
		return nil
	}
	// output is the same size as bA
	c := bA.copyBits(bA.Bits)
	// Only iterate to the minimum size between the two.
	// If o is longer, those bits are ignored.
	// If bA is longer, then skipping those iterations is equivalent
	// to right padding with 0's
	smaller := MinInt(len(bA.Elems), len(o.Elems))
	for i := 0; i < smaller; i++ {
		// &^ is and not in golang
		c.Elems[i] &^= o.Elems[i]
	}
	return c
}

// IsEmpty returns true iff all bits in the bit array are 0
func (bA *BitArray) IsEmpty() bool {
	if bA == nil {
		return true // should this be opposite?
	}
	for _, e := range bA.Elems {
		if e > 0 {
			return false
		}
	}
	return true
}

// IsFull returns true iff all bits in the bit array are 1.
func (bA *BitArray) IsFull() bool {
	if bA == nil {
		return true
	}

	// Check all elements except the last
	for _, elem := range bA.Elems[:len(bA.Elems)-1] {
		if (^elem) != 0 {
			return false
		}
	}

	// Check that the last element has (lastElemBits) 1's
	lastElemBits := (bA.Bits+63)%64 + 1
	lastElem := bA.Elems[len(bA.Elems)-1]
	return (lastElem+1)&((uint64(1)<<uint(lastElemBits))-1) == 0
}

// PickRandom returns a random index for a set bit in the bit array.
// If there is no such value, it returns 0, false.
// It uses the global randomness in `random.go` to get this index.
func (bA *BitArray) PickRandom() (uint32, bool) {
	if bA == nil {
		return 0, false
	}

	trueIndices := bA.getTrueIndices()

	if len(trueIndices) == 0 { // no bits set to true
		return 0, false
	}

	return trueIndices[RandIntn(int(len(trueIndices)))], true
}

func (bA *BitArray) getTrueIndices() []uint32 {
	trueIndices := make([]uint32, 0, bA.Bits)
	curBit := uint32(0)
	numElems := len(bA.Elems)
	// set all true indices
	for i := 0; i < numElems-1; i++ {
		elem := bA.Elems[i]
		if elem == 0 {
			curBit += 64
			continue
		}
		for j := 0; j < 64; j++ {
			if (elem & (uint64(1) << uint64(j))) > 0 {
				trueIndices = append(trueIndices, curBit)
			}
			curBit++
		}
	}
	// handle last element
	lastElem := bA.Elems[numElems-1]
	numFinalBits := bA.Bits - curBit
	for i := uint32(0); i < numFinalBits; i++ {
		if (lastElem & (uint64(1) << uint64(i))) > 0 {
			trueIndices = append(trueIndices, curBit)
		}
		curBit++
	}
	return trueIndices
}

// String returns a string representation of BitArray: BA{<bit-string>},
// where <bit-string> is a sequence of 'x' (1) and '_' (0).
// The <bit-string> includes spaces and newlines to help people.
// For a simple sequence of 'x' and '_' characters with no spaces or newlines,
// see the MarshalJSON() method.
// Example: "BA{_x_}" or "nil-BitArray" for nil.
func (bA *BitArray) String() string {
	return bA.StringIndented("")
}

// StringIndented returns the same thing as String(), but applies the indent
// at every 10th bit, and twice at every 50th bit.
func (bA *BitArray) StringIndented(indent string) string {
	if bA == nil {
		return "nil-BitArray"
	}
	return bA.stringIndented(indent)
}

func (bA *BitArray) stringIndented(indent string) string {
	lines := []string{}
	bits := ""
	for i := uint32(0); i < bA.Bits; i++ {
		if bA.getIndex(i) {
			bits += "x"
		} else {
			bits += "_"
		}
		if i%100 == 99 {
			lines = append(lines, bits)
			bits = ""
		}
		if i%10 == 9 {
			bits += indent
		}
		if i%50 == 49 {
			bits += indent
		}
	}
	if len(bits) > 0 {
		lines = append(lines, bits)
	}
	return fmt.Sprintf("BA{%v:%v}", bA.Bits, strings.Join(lines, indent))
}

// Bytes returns the byte representation of the bits within the bitarray.
func (bA *BitArray) Bytes() []byte {

	numBytes := (bA.Bits + 7) / 8
	bytes := make([]byte, numBytes)
	for i := 0; i < len(bA.Elems); i++ {
		elemBytes := [8]byte{}
		binary.LittleEndian.PutUint64(elemBytes[:], bA.Elems[i])
		copy(bytes[i*8:], elemBytes[:])
	}
	return bytes
}

// Update sets the bA's bits to be that of the other bit array.
// The copying begins from the begin of both bit arrays.
func (bA *BitArray) Update(o *BitArray) {
	if bA == nil || o == nil {
		return
	}

	copy(bA.Elems, o.Elems)
}

// MarshalJSON implements json.Marshaler interface by marshaling bit array
// using a custom format: a string of '-' or 'x' where 'x' denotes the 1 bit.
func (bA *BitArray) MarshalJSON() ([]byte, error) {
	if bA == nil {
		return []byte("null"), nil
	}

	bits := `"`
	for i := uint32(0); i < bA.Bits; i++ {
		if bA.getIndex(i) {
			bits += `x`
		} else {
			bits += `_`
		}
	}
	bits += `"`
	return []byte(bits), nil
}

var bitArrayJSONRegexp = regexp.MustCompile(`\A"([_x]*)"\z`)

// UnmarshalJSON implements json.Unmarshaler interface by unmarshaling a custom
// JSON description.
func (bA *BitArray) UnmarshalJSON(bz []byte) error {
	b := string(bz)
	if b == "null" {
		// This is required e.g. for encoding/json when decoding
		// into a pointer with pre-allocated BitArray.
		bA.Bits = 0
		bA.Elems = nil
		return nil
	}

	// Validate 'b'.
	match := bitArrayJSONRegexp.FindStringSubmatch(b)
	if match == nil {
		return fmt.Errorf("BitArray in JSON should be a string of format %q but got %s", bitArrayJSONRegexp.String(), b)
	}
	bits := match[1]

	// Construct new BitArray and copy over.
	numBits := uint32(len(bits))
	bA2 := NewBitArray(numBits)
	for i := uint32(0); i < numBits; i++ {
		if bits[i] == 'x' {
			bA2.SetIndex(i, true)
		}
	}
	*bA = *bA2
	return nil
}

func MaxInt64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}

func MaxUInt(a, b uint32) uint32 {
	if a > b {
		return a
	}
	return b
}

//-----------------------------------------------------------------------------

func MinInt64(a, b int64) int64 {
	if a < b {
		return a
	}
	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func MinUInt(a, b uint32) uint32 {
	if a < b {
		return a
	}
	return b
}
func RandBytes(n int) []byte {
	bs := make([]byte, n)
	for i := 0; i < len(bs); i++ {
		bs[i] = byte(rand.Intn(n) & 0xFF)
	}
	return bs
}

func Rand32Bytes(n uint32) []byte {
	bs := make([]byte, n)
	for i := 0; i < len(bs); i++ {
		bs[i] = byte(rand.Int31n(int32(n)) & 0xFF)
	}
	return bs
}

func RandIntn(n int) int {
	return rand.Intn(n)
}

func RandInt31n(n int32) uint32 {
	return uint32(rand.Int31n(n))
}
