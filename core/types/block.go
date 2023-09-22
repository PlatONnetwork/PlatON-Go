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
	"crypto/ecdsa"
	"encoding/binary"
	"fmt"
	"io"
	"math/big"
	"sync/atomic"
	"time"
	"unsafe"

	"golang.org/x/crypto/sha3"

	json2 "github.com/PlatONnetwork/PlatON-Go/common/json"

	"github.com/PlatONnetwork/PlatON-Go/log"

	"github.com/PlatONnetwork/PlatON-Go/crypto"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

var (
	EmptyRootHash  = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")
	EmptyUncleHash = rlpHash([]*Header(nil))

	// Extra field in the block header, maximum length
	ExtraMaxSize      = 97
	HttpEthCompatible = false
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

func (n *BlockNonce) ETHBlockNonce() ETHBlockNonce {
	var a ETHBlockNonce
	for i := 0; i < len(a); i++ {
		a[i] = n[i]
	}
	return a
}

// A ETHBlockNonce is a 64-bit hash which proves (combined with the
// mix-hash) that a sufficient amount of computation has been carried
// out on a block.
type ETHBlockNonce [8]byte

// ETHBlockNonce converts the given integer to a block nonce.
func EncodeETHNonce(i uint64) ETHBlockNonce {
	var n ETHBlockNonce
	binary.BigEndian.PutUint64(n[:], i)
	return n
}

// Uint64 returns the integer value of a block nonce.
func (n ETHBlockNonce) Uint64() uint64 {
	return binary.BigEndian.Uint64(n[:])
}

// MarshalText encodes n as a hex string with 0x prefix.
func (n ETHBlockNonce) MarshalText() ([]byte, error) {
	return hexutil.Bytes(n[:]).MarshalText()
}

// UnmarshalText implements encoding.TextUnmarshaler.
func (n *ETHBlockNonce) UnmarshalText(input []byte) error {
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
	Time        uint64         `json:"timestamp"        gencodec:"required"`
	Extra       []byte         `json:"extraData"        gencodec:"required"`
	Nonce       BlockNonce     `json:"nonce"            gencodec:"required"`

	// caches
	sealHash  atomic.Value `json:"-" rlp:"-"`
	hash      atomic.Value `json:"-" rlp:"-"`
	publicKey atomic.Value `json:"-" rlp:"-"`
}

// MarshalJSON2 marshals as JSON.
func (h Header) MarshalJSON2() ([]byte, error) {
	if HttpEthCompatible {
		type Header struct {
			ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
			Coinbase    common.Address `json:"miner"            gencodec:"required"`
			Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
			TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
			ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
			Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
			Number      *hexutil.Big   `json:"number"           gencodec:"required"`
			GasLimit    hexutil.Uint64 `json:"gasLimit"         gencodec:"required"`
			GasUsed     hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
			Time        hexutil.Uint64 `json:"timestamp"        gencodec:"required"`
			Extra       hexutil.Bytes  `json:"extraData"        gencodec:"required"`
			Nonce       ETHBlockNonce  `json:"nonce"            gencodec:"required"`
			Hash        common.Hash    `json:"hash"`

			UncleHash  common.Hash  `json:"sha3Uncles"       gencodec:"required"`
			Difficulty *hexutil.Big `json:"difficulty"       gencodec:"required"`
		}
		var enc Header
		enc.ParentHash = h.ParentHash
		enc.Coinbase = h.Coinbase
		enc.Root = h.Root
		enc.TxHash = h.TxHash
		enc.ReceiptHash = h.ReceiptHash
		enc.Bloom = h.Bloom
		enc.Number = (*hexutil.Big)(h.Number)
		enc.GasLimit = hexutil.Uint64(h.GasLimit)
		enc.GasUsed = hexutil.Uint64(h.GasUsed)
		enc.Time = hexutil.Uint64(h.Time / 1000)
		enc.Extra = h.Extra
		enc.Nonce = h.Nonce.ETHBlockNonce()
		enc.Hash = h.Hash()
		enc.UncleHash = common.ZeroHash
		enc.Difficulty = (*hexutil.Big)(h.Number)
		return json2.Marshal(&enc)
	}
	type Header struct {
		ParentHash  common.Hash    `json:"parentHash"       gencodec:"required"`
		Coinbase    common.Address `json:"miner"            gencodec:"required"`
		Root        common.Hash    `json:"stateRoot"        gencodec:"required"`
		TxHash      common.Hash    `json:"transactionsRoot" gencodec:"required"`
		ReceiptHash common.Hash    `json:"receiptsRoot"     gencodec:"required"`
		Bloom       Bloom          `json:"logsBloom"        gencodec:"required"`
		Number      *hexutil.Big   `json:"number"           gencodec:"required"`
		GasLimit    hexutil.Uint64 `json:"gasLimit"         gencodec:"required"`
		GasUsed     hexutil.Uint64 `json:"gasUsed"          gencodec:"required"`
		Time        hexutil.Uint64 `json:"timestamp"        gencodec:"required"`
		Extra       hexutil.Bytes  `json:"extraData"        gencodec:"required"`
		Nonce       BlockNonce     `json:"nonce"            gencodec:"required"`
		Hash        common.Hash    `json:"hash"`
	}
	var enc Header
	enc.ParentHash = h.ParentHash
	enc.Coinbase = h.Coinbase
	enc.Root = h.Root
	enc.TxHash = h.TxHash
	enc.ReceiptHash = h.ReceiptHash
	enc.Bloom = h.Bloom
	enc.Number = (*hexutil.Big)(h.Number)
	enc.GasLimit = hexutil.Uint64(h.GasLimit)
	enc.GasUsed = hexutil.Uint64(h.GasUsed)
	enc.Time = hexutil.Uint64(h.Time)
	enc.Extra = h.Extra
	enc.Nonce = h.Nonce
	enc.Hash = h.Hash()
	return json2.Marshal(&enc)
}

// field type overrides for gencodec
type headerMarshaling struct {
	Number   *hexutil.Big
	GasLimit hexutil.Uint64
	GasUsed  hexutil.Uint64
	Time     hexutil.Uint64
	Extra    hexutil.Bytes
	Hash     common.Hash `json:"hash"` // adds call to Hash() in MarshalJSON
}

// Hash returns the block hash of the header, which is simply the keccak256 hash of its
// RLP encoding.
func (h *Header) Hash() common.Hash {
	return rlpHash(h)
}

// SanityCheck checks a few basic things -- these checks are way beyond what
// any 'sane' production values should hold, and can mainly be used to prevent
// that the unbounded fields are stuffed with junk data to add processing
// overhead
func (h *Header) SanityCheck() error {
	if h.Number != nil && !h.Number.IsUint64() {
		return fmt.Errorf("too large block number: bitlen %d", h.Number.BitLen())
	}
	if eLen := len(h.Extra); eLen > ExtraMaxSize {
		return fmt.Errorf("too large block extradata: size %d", eLen)
	}
	return nil
}

func (h *Header) CacheHash() common.Hash {
	if hash := h.hash.Load(); hash != nil {
		return hash.(common.Hash)
	}
	v := rlpHash(h)
	h.hash.Store(v)
	return v
}

func (h *Header) CachePublicKey() *ecdsa.PublicKey {
	if pk := h.publicKey.Load(); pk != nil {
		return pk.(*ecdsa.PublicKey)
	}

	sign := h.Extra[32:97]
	sealhash := h.SealHash().Bytes()

	pk, err := crypto.SigToPub(sealhash, sign)
	if err != nil {
		log.Error("cache publicKey fail,sigToPub fail", "err", err)
		return nil
	}
	h.publicKey.Store(pk)
	return pk
}

// SealHash returns the keccak256 seal hash of b's header.
// The seal hash is computed on the first call and cached thereafter.
func (h *Header) SealHash() (hash common.Hash) {
	if sealHash := h.sealHash.Load(); sealHash != nil {
		return sealHash.(common.Hash)
	}
	v := h._sealHash()
	h.sealHash.Store(v)
	return v
}

func (h *Header) _sealHash() (hash common.Hash) {
	extra := h.Extra

	hasher := sha3.NewLegacyKeccak256()
	if len(h.Extra) > 32 {
		extra = h.Extra[0:32]
	}
	rlp.Encode(hasher, []interface{}{
		h.ParentHash,
		h.Coinbase,
		h.Root,
		h.TxHash,
		h.ReceiptHash,
		h.Bloom,
		h.Number,
		h.GasLimit,
		h.GasUsed,
		h.Time,
		extra,
		h.Nonce,
	})

	hasher.Sum(hash[:0])
	return hash
}

// Size returns the approximate memory used by all internal contents. It is used
// to approximate and limit the memory consumption of various caches.
func (h *Header) Size() common.StorageSize {
	return common.StorageSize(unsafe.Sizeof(*h)) + common.StorageSize(len(h.Extra)+(h.Number.BitLen())/8)
}

// Signature returns the signature of seal hash from extra.
func (h *Header) Signature() []byte {
	if len(h.Extra) < 32 {
		return []byte{}
	}
	return h.Extra[32:97]
}

func (h *Header) ExtraData() []byte {
	if len(h.Extra) < 32 {
		return []byte{}
	}
	return h.Extra[:32]
}

// EmptyBody returns true if there is no additional 'body' to complete the header
// that is: no transactions and no uncles.
func (h *Header) EmptyBody() bool {
	return h.TxHash == EmptyRootHash
}

// EmptyReceipts returns true if there are no receipts for this header/block.
func (h *Header) EmptyReceipts() bool {
	return h.ReceiptHash == EmptyRootHash
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

	CalTxFromCH chan int
}

// "external" block encoding. used for eth protocol, etc.
type extblock struct {
	Header    *Header
	Txs       []*Transaction
	ExtraData []byte
}

// NewBlock creates a new block. The input data is copied,
// changes to header and to the field values will not affect the
// block.
//
// The values of TxHash, ReceiptHash and Bloom in header
// are ignored and set to values derived from the given txs
// and receipts.
func NewBlock(header *Header, txs []*Transaction, receipts []*Receipt, hasher TrieHasher) *Block {
	b := &Block{header: CopyHeader(header)}

	// TODO: panic if len(txs) != len(receipts)
	if len(txs) == 0 {
		b.header.TxHash = EmptyRootHash
	} else {
		b.header.TxHash = DeriveSha(Transactions(txs), hasher)
		b.transactions = make(Transactions, len(txs))
		copy(b.transactions, txs)
	}

	if len(receipts) == 0 {
		b.header.ReceiptHash = EmptyRootHash
	} else {
		b.header.ReceiptHash = DeriveSha(Receipts(receipts), hasher)
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
func (b *Block) ExtraData() []byte             { return common.CopyBytes(b.extraData) }
func (b *Block) Number() *big.Int              { return new(big.Int).Set(b.header.Number) }
func (b *Block) GasLimit() uint64              { return b.header.GasLimit }
func (b *Block) GasUsed() uint64               { return b.header.GasUsed }
func (b *Block) Time() uint64                  { return b.header.Time }

func (b *Block) NumberU64() uint64        { return b.header.Number.Uint64() }
func (b *Block) Nonce() []byte            { return common.CopyBytes(b.header.Nonce.Bytes()) }
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

// SanityCheck can be used to prevent that unbounded fields are
// stuffed with junk data to add processing overhead
func (b *Block) SanityCheck() error {
	return b.header.SanityCheck()
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
