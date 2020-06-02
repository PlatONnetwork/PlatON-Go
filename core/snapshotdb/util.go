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

package snapshotdb

import (
	"bytes"
	"io"
	"math/big"
	"math/rand"
	"sort"
	"time"

	"github.com/PlatONnetwork/PlatON-Go/core/types"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func generateKVHash(k, v []byte, hash common.Hash) common.Hash {
	var buf bytes.Buffer
	buf.Write(k)
	buf.Write(v)
	buf.Write(hash.Bytes())
	return rlpHash(buf.Bytes())
}

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func encode(x interface{}) ([]byte, error) {
	return rlp.EncodeToBytes(x)
}

func decode(r io.Reader, val interface{}) error {
	return rlp.Decode(r, val)
}

func generateHeader(num *big.Int, parentHash common.Hash) *types.Header {
	h := new(types.Header)
	h.Number = num
	h.ParentHash = parentHash
	return h
}

func generateHash(n string) common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(n))
	return rlpHash(buf.Bytes())
}

func randomString2(s string) []byte {
	b := new(bytes.Buffer)
	if s != "" {
		b.Write([]byte(s))
	}
	for i := 0; i < 8; i++ {
		b.WriteByte(' ' + byte(rand.Uint64()))
	}
	return b.Bytes()
}

func generatekv(n int) kvs {
	rand.Seed(time.Now().UnixNano())
	kvs := make(kvs, n)
	for i := 0; i < n; i++ {
		kvs[i] = kv{
			key:   randomString2(""),
			value: randomString2(""),
		}
	}
	sort.Sort(kvs)
	return kvs
}

func generatekvWithPrefix(n int, p string) kvs {
	rand.Seed(time.Now().UnixNano())
	kvs := make(kvs, n)
	for i := 0; i < n; i++ {
		kvs[i] = kv{
			key:   randomString2(p),
			value: randomString2(p),
		}
	}
	sort.Sort(kvs)
	return kvs
}

type kv struct {
	key   []byte
	value []byte
}

type kvs []kv

func (k kvs) Len() int {
	return len(k)
}

func (k kvs) Less(i, j int) bool {
	n := bytes.Compare(k[i].key, k[j].key)
	if n == -1 {
		return true
	}
	return false
}

func (k kvs) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

type kvsMaxToMin []kv

func (k kvsMaxToMin) Len() int {
	return len(k)
}

func (k kvsMaxToMin) Less(i, j int) bool {
	if bytes.Compare(k[i].key, k[j].key) >= 0 {
		return true
	}
	return false
}

func (k kvsMaxToMin) Swap(i, j int) {
	k[i], k[j] = k[j], k[i]
}

func (k *kvsMaxToMin) Push(x interface{}) {
	*k = append(*k, x.(kv))
}

func (k *kvsMaxToMin) Pop() interface{} {
	n := len(*k)
	x := (*k)[n-1]
	*k = (*k)[:n-1]
	return x
}
