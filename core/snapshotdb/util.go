package snapshotdb

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
	"io"
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
