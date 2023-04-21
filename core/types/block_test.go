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

package types

import (
	"bytes"
	"fmt"
	"hash"
	"math/big"
	"reflect"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/math"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/params"

	"golang.org/x/crypto/sha3"

	"github.com/PlatONnetwork/PlatON-Go/common/hexutil"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// testHasher is the helper tool for transaction/receipt list hashing.
// The original hasher is trie, in order to get rid of import cycle,
// use the testing hasher instead.
type testHasher struct {
	hasher hash.Hash
}

func newHasher() *testHasher {
	return &testHasher{hasher: sha3.NewLegacyKeccak256()}
}

func (h *testHasher) Reset() {
	h.hasher.Reset()
}

func (h *testHasher) Update(key, val []byte) {
	h.hasher.Write(key)
	h.hasher.Write(val)
}

func (h *testHasher) Hash() common.Hash {
	return common.BytesToHash(h.hasher.Sum(nil))
}

// from bcValidBlockTest.json, "SimpleTx"
func TestBlockEncoding(t *testing.T) {
	blockEnc := common.FromHex("f90264f901fda00000000000000000000000000000000000000000000000000000000000000000948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a00000000000000000000000000000000000000000000000000000000000000000a00000000000000000000000000000000000000000000000000000000000000000b901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000080832fefd8825208845506eb0780b8510376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000f861f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a8023a09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b180")
	var block Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		t.Fatal("decode error: ", err)
	}

	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
		}
	}
	check("GasLimit", block.GasLimit(), uint64(3141592))
	check("GasUsed", block.GasUsed(), uint64(21000))
	check("Coinbase", block.Coinbase(), common.MustBech32ToAddress("lax13zy0ruv447se9nlwscrfskzvqv85e8d35gau40"))
	check("Root", block.Root(), common.HexToHash("ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017"))
	check("Hash", block.Hash(), common.HexToHash("499987a73fa100f582328c92c1239262edf5c0a3479face652c89f60314aa805"))
	check("Nonce", block.Nonce(), EncodeNonce(hexutil.MustDecode("0x0376e56dffd12ab53bb149bda4e0cbce2b6aabe4cccc0df0b5a39e12977a2fcd23")).Bytes())
	check("Time", block.Time(), uint64(1426516743))
	check("Size", block.Size(), common.StorageSize(len(blockEnc)))

	tx1 := NewTransaction(0, common.MustBech32ToAddress("lax1p908ht4x5mrufsklawtha7kry6h42tv8sxxrdc"), big.NewInt(10), 50000, big.NewInt(10), nil)

	tx1, _ = tx1.WithSignature(NewEIP155Signer(new(big.Int)), common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100"))
	fmt.Println(block.Transactions()[0].Hash())
	fmt.Println(tx1.Hash())
	check("len(Transactions)", len(block.Transactions()), 1)
	check("Transactions[0].Hash", block.Transactions()[0].Hash(), tx1.Hash())

	ourBlockEnc, err := rlp.EncodeToBytes(&block)
	if err != nil {
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourBlockEnc, blockEnc) {
		t.Errorf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockEnc, blockEnc)
	}
}

func TestEIP2718BlockEncoding(t *testing.T) {
	blockEnc := common.FromHex("f90319f90211a00000000000000000000000000000000000000000000000000000000000000000a01dcc4de8dec75d7aab85b567b6ccd41ad312451b948a7413f0a142fd40d49347948888f1f195afa192cfee860698584c030f4c9db1a0ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017a0e6e49996c7ec59f7a23d22b83239a60151512c65613bf84a0d7da336399ebc4aa0cafe75574d59780665a97fbfd11365c7545aa8f1abf4e5e12e8243334ef7286bb901000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000083020000820200832fefd882a410845506eb0796636f6f6c65737420626c6f636b206f6e20636861696ea0bd4472abb6659ebe3ee06ee4d7b72a00a9f4d001caca51342001075469aff49888a13a5a8c8f2bb1c4f90101f85f800a82c35094095e7baea6a6c7c4c2dfeb977efac326af552d870a801ba09bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094fa08a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b1b89e01f89b01800a8301e24194095e7baea6a6c7c4c2dfeb977efac326af552d878080f838f7940000000000000000000000000000000000000001e1a0000000000000000000000000000000000000000000000000000000000000000001a03dbacc8d0259f2508625e97fdfc57cd85fdd16e5821bc2c10bdd1a52649e8335a0476e10695b183a87b0aa292a7f4b78ef0c3fbe62aa2c42c84e1d9c3da159ef14c0")
	var block Block
	if err := rlp.DecodeBytes(blockEnc, &block); err != nil {
		t.Fatal("decode error: ", err)
	}

	check := func(f string, got, want interface{}) {
		if !reflect.DeepEqual(got, want) {
			t.Errorf("%s mismatch: got %v, want %v", f, got, want)
		}
	}
	check("GasLimit", block.GasLimit(), uint64(3141592))
	check("GasUsed", block.GasUsed(), uint64(42000))
	check("Coinbase", block.Coinbase(), common.HexToAddress("8888f1f195afa192cfee860698584c030f4c9db1"))
	check("Root", block.Root(), common.HexToHash("ef1552a40b7165c3cd773806b9e0c165b75356e0314bf0706f279c729f51e017"))
	check("Nonce", block.Nonce(), uint64(0xa13a5a8c8f2bb1c4))
	check("Time", block.Time(), uint64(1426516743))
	check("Size", block.Size(), common.StorageSize(len(blockEnc)))

	// Create legacy tx.
	to := common.HexToAddress("095e7baea6a6c7c4c2dfeb977efac326af552d87")
	tx1 := NewTx(&LegacyTx{
		Nonce:    0,
		To:       &to,
		Value:    big.NewInt(10),
		Gas:      50000,
		GasPrice: big.NewInt(10),
	})
	sig := common.Hex2Bytes("9bea4c4daac7c7c52e093e6a4c35dbbcf8856f1af7b059ba20253e70848d094f8a8fae537ce25ed8cb5af9adac3f141af69bd515bd2ba031522df09b97dd72b100")
	tx1, _ = tx1.WithSignature(HomesteadSigner{}, sig)

	// Create ACL tx.
	addr := common.HexToAddress("0x0000000000000000000000000000000000000001")
	tx2 := NewTx(&AccessListTx{
		ChainID:    big.NewInt(1),
		Nonce:      0,
		To:         &to,
		Gas:        123457,
		GasPrice:   big.NewInt(10),
		AccessList: AccessList{{Address: addr, StorageKeys: []common.Hash{{0}}}},
	})
	sig2 := common.Hex2Bytes("3dbacc8d0259f2508625e97fdfc57cd85fdd16e5821bc2c10bdd1a52649e8335476e10695b183a87b0aa292a7f4b78ef0c3fbe62aa2c42c84e1d9c3da159ef1401")
	tx2, _ = tx2.WithSignature(NewEIP2930Signer(big.NewInt(1)), sig2)

	check("len(Transactions)", len(block.Transactions()), 2)
	check("Transactions[0].Hash", block.Transactions()[0].Hash(), tx1.Hash())
	check("Transactions[1].Hash", block.Transactions()[1].Hash(), tx2.Hash())
	check("Transactions[1].Type()", block.Transactions()[1].Type(), uint8(AccessListTxType))

	ourBlockEnc, err := rlp.EncodeToBytes(&block)
	if err != nil {
		t.Fatal("encode error: ", err)
	}
	if !bytes.Equal(ourBlockEnc, blockEnc) {
		t.Errorf("encoded block mismatch:\ngot:  %x\nwant: %x", ourBlockEnc, blockEnc)
	}
}

var benchBuffer = bytes.NewBuffer(make([]byte, 0, 32000))

func BenchmarkEncodeBlock(b *testing.B) {
	block := makeBenchBlock()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		benchBuffer.Reset()
		if err := rlp.Encode(benchBuffer, block); err != nil {
			b.Fatal(err)
		}
	}
}

func makeBenchBlock() *Block {
	var (
		key, _   = crypto.GenerateKey()
		txs      = make([]*Transaction, 70)
		receipts = make([]*Receipt, len(txs))
		signer   = NewEIP155Signer(params.TestChainConfig.ChainID)
		uncles   = make([]*Header, 3)
	)
	header := &Header{
		Number:   math.BigPow(2, 9),
		GasLimit: 12345678,
		GasUsed:  1476322,
		Time:     9876543,
		Extra:    []byte("coolest block on chain"),
	}
	for i := range txs {
		amount := math.BigPow(2, int64(i))
		price := big.NewInt(300000)
		data := make([]byte, 100)
		tx := NewTransaction(uint64(i), common.Address{}, amount, 123457, price, data)
		signedTx, err := SignTx(tx, signer, key)
		if err != nil {
			panic(err)
		}
		txs[i] = signedTx
		receipts[i] = NewReceipt(make([]byte, 32), false, tx.Gas())
	}
	for i := range uncles {
		uncles[i] = &Header{
			Number:   math.BigPow(2, 9),
			GasLimit: 12345678,
			GasUsed:  1476322,
			Time:     9876543,
			Extra:    []byte("benchmark uncle"),
		}
	}
	return NewBlock(header, txs, receipts, newHasher())
}
