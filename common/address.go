package common

import (
	"database/sql/driver"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"golang.org/x/crypto/sha3"
	"math/big"
	"strings"

	"github.com/btcsuite/btcutil/bech32"

	"github.com/PlatONnetwork/PlatON-Go/common/bech32util"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/log"
)

const (
	DefaultAddressHRP = "lat"
)

var currentAddressHRP string

func GetAddressHRP() string {
	if currentAddressHRP == "" {
		return DefaultAddressHRP
	}
	return currentAddressHRP
}

func SetAddressHRP(s string) error {
	if s == "" {
		s = DefaultAddressHRP
	}
	if len(s) != 3 {
		return errors.New("the length of addressHRP must be 3")
	}
	log.Info("the address hrp has been set", "hrp", s)
	currentAddressHRP = s
	return nil
}

func CheckAddressHRP(s string) bool {
	if currentAddressHRP != "" && s != currentAddressHRP {
		return false
	}
	return true
}

/////////// Address

// Address represents the 20 byte address of an Ethereum account.
type Address [AddressLength]byte

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToAddress(b []byte) Address {
	var a Address
	a.SetBytes(b)
	return a
}

// BigToAddress returns Address with byte values of b.
// If b is larger than len(h), b will be cropped from the left.
func BigToAddress(b *big.Int) Address { return BytesToAddress(b.Bytes()) }

// Deprecated: address to string is use bech32 now
// HexToAddress returns Address with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToAddress(s string) Address { return BytesToAddress(FromHex(s)) }

// MustBech32ToAddress returns Address with byte values of s.
// If s is Decode fail, it will return zero address.
func MustBech32ToAddress(s string) Address {
	add, err := Bech32ToAddress(s)
	if err != nil {
		log.Error("must Bech32ToAddress fail", "err", err)
		panic(err)
	}
	return add
}

// MustBech32ToAddress returns Address with byte values of s.
// If s is Decode fail, it will return zero address.
func Bech32ToAddress(s string) (Address, error) {
	hrpDecode, converted, err := bech32util.DecodeAndConvert(s)
	if err != nil {
		return Address{}, err
	}
	if !CheckAddressHRP(hrpDecode) {
		return Address{}, fmt.Errorf("the address hrp not compare right,input:%s", s)
	}

	if currentAddressHRP == "" {
		log.Warn("the address hrp not set yet", "input", s)
	} else if currentAddressHRP != hrpDecode {
		log.Warn("the address not compare current net", "want", currentAddressHRP, "input", s)
	}
	var a Address
	a.SetBytes(converted)
	return a, nil
}

// Bech32ToAddressWithoutCheckHrp returns Address with byte values of s.
// If s is Decode fail, it will return zero address.
func Bech32ToAddressWithoutCheckHrp(s string) Address {
	_, converted, err := bech32util.DecodeAndConvert(s)
	if err != nil {
		log.Error(" hrp address decode  fail", "err", err)
		panic(err)
	}
	var a Address
	a.SetBytes(converted)
	return a
}

// Deprecated: address to string is use bech32 now
// IsHexAddress verifies whether a string can represent a valid hex-encoded
// Ethereum address or not.
func IsHexAddress(s string) bool {
	if hasHexPrefix(s) {
		s = s[2:]
	}
	return len(s) == 2*AddressLength && isHex(s)
}

func IsBech32Address(s string) bool {
	hrp, _, err := bech32.Decode(s)
	if err != nil {
		return false
	}
	if !CheckAddressHRP(hrp) {
		return false
	}
	return true
}

// Bytes gets the string representation of the underlying address.
func (a Address) Bytes() []byte { return a[:] }

// Big converts an address to a big integer.
func (a Address) Big() *big.Int { return new(big.Int).SetBytes(a[:]) }

// Hash converts an address to a hash by left-padding it with zeros.
func (a Address) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.it's use for node address
func (a Address) Hex() string {
	return "0x" + a.HexWithNoPrefix()
}

func (a Address) HexWithNoPrefix() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return string(result)
}

// String implements fmt.Stringer.
func (a Address) String() string {
	return a.Bech32()
}

func (a Address) Bech32() string {
	return a.Bech32WithHRP(GetAddressHRP())
}

func (a Address) Bech32WithHRP(hrp string) string {
	if v, err := bech32util.ConvertAndEncode(hrp, a.Bytes()); err != nil {
		log.Error("address can't ConvertAndEncode to string", "err", err, "add", a.Bytes())
		return ""
	} else {
		return v
	}
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a Address) Format(s fmt.State, c rune) {
	switch string(c) {
	case "s":
		fmt.Fprintf(s, "%"+string(c), a.String())
	default:
		fmt.Fprintf(s, "%"+string(c), a[:])
	}
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *Address) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a Address) MarshalText() ([]byte, error) {
	v, err := bech32util.ConvertAndEncode(GetAddressHRP(), a.Bytes())
	if err != nil {
		return nil, err
	}
	return []byte(v), nil
}

// UnmarshalText parses a hash in hex syntax.
func (a *Address) UnmarshalText(input []byte) error {
	if hexutil.BytesHave0xPrefix(input) {
		return hexutil.UnmarshalFixedText(addressT.String(), input, a[:])
	}
	hrpDecode, converted, err := bech32util.DecodeAndConvert(string(input))
	if err != nil {
		return err
	}
	if !CheckAddressHRP(hrpDecode) {
		return fmt.Errorf("the address not compare current net,want %v,have %v", GetAddressHRP(), string(input))
	}
	a.SetBytes(converted)
	return nil
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *Address) UnmarshalJSON(input []byte) error {
	if !isString(input) {
		return hexutil.ErrNonString(addressT)
	}
	if hexutil.BytesHave0xPrefix(input[1 : len(input)-1]) {
		return hexutil.WrapTypeError(hexutil.UnmarshalFixedText(addressT.String(), input[1:len(input)-1], a[:]), addressT)
	}
	hrpDecode, v, err := bech32util.DecodeAndConvert(string(input[1 : len(input)-1]))
	if err != nil {
		return &json.UnmarshalTypeError{Value: err.Error(), Type: addressT}
	}
	if !CheckAddressHRP(hrpDecode) {
		return &json.UnmarshalTypeError{Value: fmt.Sprintf("hrpDecode not compare the current net,want %v,have %v", GetAddressHRP(), hrpDecode), Type: addressT}
	}
	a.SetBytes(v)
	return nil
}

func isString(input []byte) bool {
	return len(input) >= 2 && input[0] == '"' && input[len(input)-1] == '"'
}

// Scan implements Scanner for database/sql.
func (a *Address) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a Address) Value() (driver.Value, error) {
	return a[:], nil
}

// UnprefixedAddress allows marshaling an Address without 0x prefix.
type UnprefixedAddress Address

// UnmarshalText decodes the address from hex. The 0x prefix is optional.
func (a *UnprefixedAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedUnprefixedText("UnprefixedAddress", input, a[:])
}

// MarshalText encodes the address as hex.
func (a UnprefixedAddress) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(a[:])), nil
}

// MixedcaseAddress retains the original string, which may or may not be
// correctly checksummed
type MixedcaseAddress struct {
	addr     Address
	original string
}

// NewMixedcaseAddress constructor (mainly for testing)
func NewMixedcaseAddress(addr Address) MixedcaseAddress {
	return MixedcaseAddress{addr: addr, original: addr.Hex()}
}

// NewMixedcaseAddressFromString is mainly meant for unit-testing
func NewMixedcaseAddressFromString(hexaddr string) (*MixedcaseAddress, error) {
	if !IsHexAddress(hexaddr) {
		return nil, fmt.Errorf("Invalid address")
	}
	a := FromHex(hexaddr)
	return &MixedcaseAddress{addr: BytesToAddress(a), original: hexaddr}, nil
}

// UnmarshalJSON parses MixedcaseAddress
func (ma *MixedcaseAddress) UnmarshalJSON(input []byte) error {
	if err := hexutil.UnmarshalFixedJSON(addressT, input, ma.addr[:]); err != nil {
		return err
	}
	return json.Unmarshal(input, &ma.original)
}

// MarshalJSON marshals the original value
func (ma *MixedcaseAddress) MarshalJSON() ([]byte, error) {
	if strings.HasPrefix(ma.original, "0x") || strings.HasPrefix(ma.original, "0X") {
		return json.Marshal(fmt.Sprintf("0x%s", ma.original[2:]))
	}
	return json.Marshal(fmt.Sprintf("0x%s", ma.original))
}

// Address returns the address
func (ma *MixedcaseAddress) Address() Address {
	return ma.addr
}

// String implements fmt.Stringer
func (ma *MixedcaseAddress) String() string {
	if ma.ValidChecksum() {
		return fmt.Sprintf("%s [chksum ok]", ma.original)
	}
	return fmt.Sprintf("%s [chksum INVALID]", ma.original)
}

// ValidChecksum returns true if the address has valid checksum
func (ma *MixedcaseAddress) ValidChecksum() bool {
	return ma.original == ma.addr.Hex()
}

// Original returns the mixed-case input string
func (ma *MixedcaseAddress) Original() string {
	return ma.original
}

// BytesToAddress returns Address with value b.
// If b is larger than len(h), b will be cropped from the left.
func BytesToNodeAddress(b []byte) NodeAddress {
	var a NodeAddress
	a.SetBytes(b)
	return a
}

// HexToNodeAddress returns NodeAddress with byte values of s.
// If s is larger than len(h), s will be cropped from the left.
func HexToNodeAddress(s string) NodeAddress { return NodeAddress(BytesToAddress(FromHex(s))) }

type NodeAddress Address

// Bytes gets the string representation of the underlying address.
func (a NodeAddress) Bytes() []byte { return a[:] }

// Hash converts an address to a hash by left-padding it with zeros.
func (a NodeAddress) Hash() Hash { return BytesToHash(a[:]) }

// Hex returns an EIP55-compliant hex string representation of the address.
func (a NodeAddress) Hex() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return "0x" + string(result)
}

func (a NodeAddress) HexWithNoPrefix() string {
	unchecksummed := hex.EncodeToString(a[:])
	sha := sha3.NewLegacyKeccak256()
	sha.Write([]byte(unchecksummed))
	hash := sha.Sum(nil)

	result := []byte(unchecksummed)
	for i := 0; i < len(result); i++ {
		hashByte := hash[i/2]
		if i%2 == 0 {
			hashByte = hashByte >> 4
		} else {
			hashByte &= 0xf
		}
		if result[i] > '9' && hashByte > 7 {
			result[i] -= 32
		}
	}
	return string(result)
}

// String implements fmt.Stringer.
func (a NodeAddress) String() string {
	return a.Hex()
}

// Format implements fmt.Formatter, forcing the byte slice to be formatted as is,
// without going through the stringer interface used for logging.
func (a NodeAddress) Format(s fmt.State, c rune) {
	fmt.Fprintf(s, "%"+string(c), a[:])
}

// SetBytes sets the address to the value of b.
// If b is larger than len(a) it will panic.
func (a *NodeAddress) SetBytes(b []byte) {
	if len(b) > len(a) {
		b = b[len(b)-AddressLength:]
	}
	copy(a[AddressLength-len(b):], b)
}

// MarshalText returns the hex representation of a.
func (a NodeAddress) MarshalText() ([]byte, error) {
	return hexutil.Bytes(a[:]).MarshalText()
}

// UnmarshalText parses a hash in hex syntax.
func (a *NodeAddress) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("Address", input, a[:])
}

// UnmarshalJSON parses a hash in hex syntax.
func (a *NodeAddress) UnmarshalJSON(input []byte) error {
	return hexutil.UnmarshalFixedJSON(addressT, input, a[:])
}

// Scan implements Scanner for database/sql.
func (a *NodeAddress) Scan(src interface{}) error {
	srcB, ok := src.([]byte)
	if !ok {
		return fmt.Errorf("can't scan %T into Address", src)
	}
	if len(srcB) != AddressLength {
		return fmt.Errorf("can't scan []byte of len %d into Address, want %d", len(srcB), AddressLength)
	}
	copy(a[:], srcB)
	return nil
}

// Value implements valuer for database/sql.
func (a NodeAddress) Value() (driver.Value, error) {
	return a[:], nil
}
