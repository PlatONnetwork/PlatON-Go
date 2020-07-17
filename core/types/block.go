// Copyright 2014 The go-ethereum Authors
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

// Package types contains data types related to Ethereum consensus.
package types

import (
	"fmt"
	"io"
	"math/big"
	"sort"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var (
	EmptyRootHash = DeriveSha(Transactions{})
)

// BlockNonce is an 81-byte vrf proof containing random numbers
// Used to verify the block when receiving the block
type BlockNonce [81]byte

// EncodeNonce converts the given byte to a block nonce.
func EncodeNonce(v []byte) BlockNonce {
	var n BlockNonce
	copy(n[:], v)
	return n
}

func (n BlockNonce) Bytes() []byte {
	return n[:]
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n BlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *BlockNonce) UnmarshalText(input []byte) error {
	return hexutil.UnmarshalFixedText("BlockNonce", input, n[:])
}

//go:generate gencodec -type Header -field-override headerMarshaling -out gen_header_json.go

// Header represents a block header in the Ethereum blockchain.
type Header struct {
	ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
	Coinbase    common.Address `json:"miner"            gencodec:"required"`
	Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
	TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
	ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
	Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
	Number      *big.Int       `json:"number"           gencodec:"required"`
	GasLimit    uint64         `json:"gasLimit"         gencodec:"required"`
	GasUsed     uint64         `json:"gasUsed"          gencodec:"required"`
	Time        *big.Int       `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	Nonce       BlockNonce     `json:"nonce"            gencodec:"required"`

	// caches
	sealHash atomic.Value `json:"-" rlp:"-"`
}

// field type overrides for gencodec
type headerMarshaling struct {
	Number   *hexutil.Big
	GasLimit hexutil.Uint64
	GasUsed  hexutil.Uint64
	Time     *hexutil.Big
	Extra    hexutil.Bytes
	Hash     common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	return rlpHash(h)
}

// SealHash returns the keccak256 seal hash of b's header.
// The seal hash is computed on the first call and cached thereafter.
func (header *Header) SealHash() (hash common.Hash) {
	if sealHash := header.sealHash.Load(); sealHash != nil {
		return sealHash.(common.Hash)
	}
	v := header._sealHash()
	header.sealHash.Store(v)
	return v
}

func (header *Header) _sealHash() (hash common.Hash) {
	extra := header.Extra

	hasher := sha3.NewKeccak256()
	if len(header.Extra) > 32 {
		extra = header.Extra[0:32]
	}
	rlp.Encode(hasher, []interface{}{
		header.ParentHash,
		header.Coinbase,
		header.Root,
		header.TxHash,
		header.ReceiptHash,
		header.Bloom,
		header.Number,
		header.GasLimit,
		header.GasUsed,
		header.Time,
		extra,
		header.Nonce,
	})

	hasher.Sum(hash[:0])
	return hash
}

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (h *Header) Size() common.StorageSize {
	return common.StorageSize(unsafe.Sizeof(*h)) + common.StorageSize(len(h.Extra)+(h.Number.BitLen()+h.Time.BitLen())/8)
}

// Signature returns the signature of seal hash from extra.
func (h *Header) Signature() []byte {
	if len(h.Extra) < 32 {
		return []byte{}
	}
	return h.Extra[32:]
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

// Body is a simple (mutable, non-safe) data container for storing and moving
// a block's data contents (transactions) together.
type Body struct {
	Transactions []*Transaction
	ExtraData    []byte
}

// Block represents an entire block in the Ethereum blockchain.
type Block struct {
	header       *Header
	transactions Transactions

	// caches
	hash atomic.Value
	size atomic.Value

	// These fields are used by package eth to track
	// inter-peer block relay.
	ReceivedAt   time.Time
	ReceivedFrom interface{}
	extraData    []byte

	CalTxFromCH chan struct{}
}

// [deprecated by eth/63]
// StorageBlock defines the RLP encoding of a Block stored in the
// state database. The StorageBlock encoding contains fields that
// would otherwise need to be recomputed.
type StorageBlock Block

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header    *Header
	Txs       []*Transaction
	ExtraData []byte
}

// [deprecated by eth/63]
// "storage" block encoding. used for database.
type storageblock struct {
	Header *Header
	Txs    []*Transaction
}

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs
// and receipts.
func NewBlock(header *Header, txs []*Transaction, receipts []*Receipt) *Block {
	b := &Block{header: CopyHeader(header)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.header.TxHash = EmptyRootHash
	} else {
		b.header.TxHash = DeriveSha(Transactions(txs))
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)
	}

	if len(receipts) == 0 {
		b.header.ReceiptHash = EmptyRootHash
	} else {
		b.header.ReceiptHash = DeriveSha(Receipts(receipts))
		b.header.Bloom = CreateBloom(receipts)
	}

	return b
}

// NewBlockWithHeader creates a block with the given header data. The
// header data is copied, changes to header and to the field values
// will not affect the block.
func NewBlockWithHeader(header *Header) *Block {
	return &Block{header: CopyHeader(header)}
}

// NewSimplifiedBlock creates a block with the given number and hash data.
func NewSimplifiedBlock(number uint64, hash common.Hash) *Block {
	header := &Header{
		Number: big.NewInt(int64(number)),
	}
	block := NewBlockWithHeader(header)
	block.hash.Store(hash)
	return block
}

// CopyHeader creates a deep copy of a block header to prevent side effects from
// modifying a header variable.
func CopyHeader(h *Header) *Header {
	cpy := *h
	if cpy.Time = new(big.Int); h.Time != nil {
		cpy.Time.Set(h.Time)
	}
	if cpy.Number = new(big.Int); h.Number != nil {
		cpy.Number.Set(h.Number)
	}
	if len(h.Extra) > 0 {
		cpy.Extra = make([]byte, len(h.Extra))
		copy(cpy.Extra, h.Extra)
	}
	return &cpy
}

// DecodeRLP decodes the Ethereum
func (b *Block) DecodeRLP(s *rlp.Stream) error {
	var eb extblock
	_, size, _ := s.Kind()
	if err := s.Decode(&eb); err != nil {
		return err
	}
	b.header, b.transactions, b.extraData = eb.Header, eb.Txs, eb.ExtraData
	b.size.Store(common.StorageSize(rlp.ListSize(size)))
	return nil
}

// EncodeRLP serializes b into the Ethereum RLP block format.
func (b *Block) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, extblock{
		Header:    b.header,
		Txs:       b.transactions,
		ExtraData: b.extraData,
	})
}

// [deprecated by eth/63]
func (b *StorageBlock) DecodeRLP(s *rlp.Stream) error {
	var sb storageblock
	if err := s.Decode(&sb); err != nil {
		return err
	}
	b.header, b.transactions = sb.Header, sb.Txs
	return nil
}

// TODO: copies

func (b *Block) Transactions() Transactions { return b.transactions }

func (b *Block) Transaction(hash common.Hash) *Transaction {
	for _, transaction := range b.transactions {
		if transaction.Hash() == hash {
			return transaction
		}
	}
	return nil
}
func (b *Block) SetExtraData(extraData []byte) { b.extraData = extraData }
func (b *Block) ExtraData() []byte             { return b.extraData }
func (b *Block) Number() *big.Int              { return new(big.Int).Set(b.header.Number) }
func (b *Block) GasLimit() uint64              { return b.header.GasLimit }
func (b *Block) GasUsed() uint64               { return b.header.GasUsed }
func (b *Block) Time() *big.Int                { return new(big.Int).Set(b.header.Time) }

func (b *Block) NumberU64() uint64        { return b.header.Number.Uint64() }
func (b *Block) Nonce() []byte            { return b.header.Nonce[:] }
func (b *Block) Bloom() Bloom             { return b.header.Bloom }
func (b *Block) Coinbase() common.Address { return b.header.Coinbase }
func (b *Block) Root() common.Hash        { return b.header.Root }
func (b *Block) ParentHash() common.Hash  { return b.header.ParentHash }
func (b *Block) TxHash() common.Hash      { return b.header.TxHash }
func (b *Block) ReceiptHash() common.Hash { return b.header.ReceiptHash }
func (b *Block) Extra() []byte            { return common.CopyBytes(b.header.Extra) }

func (b *Block) Header() *Header { return CopyHeader(b.header) }

// Body returns the non-header content of the block.
func (b *Block) Body() *Body { return &Body{b.transactions, b.extraData} }

// Size returns the true RLP encoded storage size of the block, either by encoding
// and returning it, or returning a previsouly cached value.
func (b *Block) Size() common.StorageSize {
	if size := b.size.Load(); size != nil {
		return size.(common.StorageSize)
	}
	c := writeCounter(0)
	rlp.Encode(&c, b)
	b.size.Store(common.StorageSize(c))
	return common.StorageSize(c)
}

type writeCounter common.StorageSize

func (c *writeCounter) Write(b []byte) (int, error) {
	*c += writeCounter(len(b))
	return len(b), nil
}

// WithSeal returns a new block with the data from b but the header replaced with
// the sealed one.
func (b *Block) WithSeal(header *Header) *Block {
	cpy := *header

	return &Block{
		header:       &cpy,
		transactions: b.transactions,
	}
}

// WithBody returns a new block with the given transaction.
func (b *Block) WithBody(transactions []*Transaction, extraData []byte) *Block {
	block := &Block{
		header:       CopyHeader(b.header),
		transactions: make([]*Transaction, len(transactions)),
		extraData:    make([]byte, len(extraData)),
	}
	copy(block.transactions, transactions)
	copy(block.extraData, extraData)
	return block
}

// Hash returns the keccak256 hash of b's header.
// The hash is computed on the first call and cached thereafter.
func (b *Block) Hash() common.Hash {
	if hash := b.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := b.header.Hash()
	b.hash.Store(v)
	return v
}

type Blocks []*Block

func (b Blocks) String() string {
	s := "["
	for _, v := range b {
		s += fmt.Sprintf("[hash:%s, number:%d]", v.Hash().TerminalString(), v.NumberU64())
	}
	s += "]"
	return s
}

type BlockBy func(b1, b2 *Block) bool

func (self BlockBy) Sort(blocks Blocks) {
	bs := blockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

type blockSorter struct {
	blocks Blocks
	by     func(b1, b2 *Block) bool
}

func (self blockSorter) Len() int { return len(self.blocks) }
func (self blockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self blockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

func Number(b1, b2 *Block) bool { return b1.header.Number.Cmp(b2.header.Number) < 0 }
