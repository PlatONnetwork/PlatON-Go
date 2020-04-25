// Copyright 2015 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"database/sql/driver"
	"encoding/hex"
	"fmt"
	"math/big"
	"math/rand"
	"reflect"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
)

// Lengths of hashes and addresses in bytes.
const (
	// HashLength is the expected length of the hash
	HashLength        = 32
	VersionSignLength = 65
	// AddressLength is the expected length of the address
	AddressLength          = 20
	BlockConfirmSignLength = 65
	ExtraSeal              = 65
)

var (
	hashT    = reflect.TypeOf(Hash{})
	addressT = reflect.TypeOf(Address{})

	ZeroHash     = HexToHash(Hash{}.String())
	ZeroAddr     = Address{}
	ZeroNodeAddr = NodeAddress{}
)

// Hash represents the 32 byte Keccak256 hash of arbitrary data.
type Hash [HashLength]byte

// BytesToHash sets b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BytesToHash(b []byte) Hash {
	var h Hash
	h.SetBytes(b)
	return h
}

// BigToHash sets byte representation of b to hash.
// If b is larger than len(h), b will be cropped from the left.
func BigToHash(b *big.Int) Hash { return BytesToHash(b.Bytes()) }

// HexToHash sets byte representation of s to hash.
// If b is larger than len(h), b will be cropped from the left.
func HexToHash(s string) Hash { return BytesToHash(FromHex(s)) }

// Bytes gets the byte representation of the underlying hash.
func (h Hash) Bytes() []byte { return h[:] }

// Big converts a hash to a big integer.
func (h Hash) Big() *big.Int { return new(big.Int).SetBytes(h[:]) }

// Hex converts a hash to a hex string.
func (h Hash) Hex() string { return hexutil.Encode(h[:]) }

// Hex converts a hash to a hex string with no prefix of 0x.
func (h Hash) HexWithNoPrefix() string {
	hex := hexutil.Encode(h[:])
	return strings.TrimPrefix(hex, "0x")
}

// TerminalString implements log.TerminalStringer, formatting a string for console
// output during logging.
func (h Hash) TerminalString() string {
	return fmt.Sprintf("%x…%x", h[:3], h[29:])
}

// String implements the stringer interface and is used also by the logger when
// doing full logging into a file.
func (h Hash) String() string {
	return h.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (h Hash) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), h[:])
}

// UnmarshalText parses a hash in hex syntax.
func (h *Hash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Hash", input, h[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (h *Hash) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(hashT, input, h[:])
}

// MarshalText returns the hex representation of h.
func (h Hash) MarshalText() ([]byte, error) {
	return hexutil.Bytes(h[:]).MarshalText()
}

// SetBytes sets the hash to the value of b.
// If b is larger than len(h), b will be cropped from the left.
func (h *Hash) SetBytes(b []byte) {
	if len(b) > len(h) {
		b = b[len(b)-HashLength:]
	}

	copy(h[HashLength-len(b):], b)
}

// Generate implements testing/quick.Generator.
func (h Hash) Generate(rand *rand.Rand, size int) reflect.Value {
	m := rand.Intn(len(h))
	for i := len(h) - 1; i > m; i-- {
		h[i] = byte(rand.Uint32())
	}
	return reflect.ValueOf(h)
}

// Scan implements Scanner for database/sql.
func (h *Hash) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Hash", src)
	}
	if len(srcB) != HashLength {
		return fmt.Errorf("can't scan []byte of len %d into Hash, want %d", len(srcB), HashLength)
	}
	copy(h[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (h Hash) Value() (driver.Value, error) {
	return h[:], nil
}

// UnprefixedHash allows marshaling a Hash without 0x prefix.
type UnprefixedHash Hash

// UnmarshalText decodes the hash from hex. The 0x prefix is optional.
func (h *UnprefixedHash) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedHash", input, h[:])
}

// MarshalText encodes the hash as hex.
func (h UnprefixedHash) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(h[:])), nil
}

type BlockConfirmSign [BlockConfirmSignLength]byte

func (sig *BlockConfirmSign) String() string {
	return fmt.Sprintf("%x", sig[:])
}

func (sig *BlockConfirmSign) SetBytes(signSlice []byte) {
	copy(sig[:], signSlice[:])
}

func (sig *BlockConfirmSign) Bytes() []byte {
	target := make([]byte, len(sig))
	copy(target[:], sig[:])
	return target
}

// MarshalText returns the hex representation of a.
func (a BlockConfirmSign) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *BlockConfirmSign) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockConfirmSign", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *BlockConfirmSign) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

func NewBlockConfirmSign(signSlice []byte) *BlockConfirmSign {
	var sign BlockConfirmSign
	copy(sign[:], signSlice[:])
	return &sign
}

type VersionSign [VersionSignLength]byte

func BytesToVersionSign(b []byte) VersionSign {
	var h VersionSign
	h.SetBytes(b)
	return h
}

func (s VersionSign) Bytes() []byte { return s[:] }

func (s VersionSign) Hex() string { return hexutil.Encode(s[:]) }

func (s VersionSign) HexWithNoPrefix() string {
	hex := hexutil.Encode(s[:])
	return strings.TrimPrefix(hex, "0x")
}

func (s VersionSign) TerminalString() string {
	return fmt.Sprintf("%x…%x", s[:3], s[61:])
}

func (s VersionSign) String() string {
	return s.Hex()
}

func (s VersionSign) Format(st fmt.State, c rune) {
	fmt.Fprintf(st, "%"+string(c), s[:])
}

func (s *VersionSign) SetBytes(b []byte) {
	if len(b) > len(s) {
		b = b[len(b)-VersionSignLength:]
	}
	copy(s[VersionSignLength-len(b):], b)
}

// MarshalText implements the encoding.TextMarshaler interface.
func (s VersionSign) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(s[:])), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (s *VersionSign) UnmarshalText(text []byte) error {
	id, err := HexSign(string(text))
	if err != nil {
		return err
	}
	*s = id
	return nil
}

// HexID converts a hex string to a NodeID.
// The string may be prefixed with 0x.
func HexSign(in string) (VersionSign, error) {
	var vs VersionSign
	b, err := hex.DecodeString(strings.TrimPrefix(in, "0x"))
	if err != nil {
		return vs, err
	} else if len(b) != len(vs) {
		return vs, fmt.Errorf("wrong length, want %d hex chars", len(vs)*2)
	}
	copy(vs[:], b)
	return vs, nil
}

/*// MustHexID converts a hex string to a NodeID.
// It panics if the string is not a valid NodeID.
func MustHexSign(in string) VersionSign {
	vs, err := HexSign(in)
	if err != nil {
		panic(err)
	}
	return vs
}*/
