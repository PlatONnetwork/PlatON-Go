package snapshotdb

import (
	"bytes"
	"io"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

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

func generateHash(n string) common.Hash {
	var buf bytes.Buffer
	buf.Write([]byte(n))
	return rlpHash(buf.Bytes())
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
