package common

import (
	"bytes"

	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func MustRlpEncode(val interface{}) []byte {
	bytes, err := rlp.EncodeToBytes(val)
	if err != nil {
		panic(err)
	}
	return bytes
}

func GenerateKVHash(k, v []byte, oldHash Hash) Hash {
	var buf bytes.Buffer
	buf.Write(k)
	buf.Write(v)
	buf.Write(oldHash.Bytes())
	return RlpHash(buf.Bytes())
}

func RlpHash(x interface{}) (h Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}
