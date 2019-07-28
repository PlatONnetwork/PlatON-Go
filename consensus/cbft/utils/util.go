package utils

import (
	"bytes"
	"io"
	"math/rand"
	"sync/atomic"

	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
)

// Returns a random offset between 0 and n
func RandomOffset(n int) int {
	if n == 0 {
		return 0
	}
	return int(rand.Uint32() % uint32(n))
}

// Convert byte array to hash. Use sha256 to
// generate a unique message hash.
func BuildHash(msgType byte, bytes []byte) common.Hash {
	bytes[0] = msgType
	hashBytes := sha3.Sum256(bytes)
	result := common.Hash{}
	result.SetBytes(hashBytes[:])
	return result
}

// A merges multiple bytes of data and
// returns the merged byte array.
func MergeBytes(bts ...[]byte) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 128))
	for _, v := range bts {
		io.Copy(buffer, bytes.NewReader(v))
	}
	temp := buffer.Bytes()
	length := len(temp)
	var response []byte
	if cap(temp) > (length + length/10) {
		response = make([]byte, length)
		copy(response, temp)
	} else {
		response = temp
	}
	return response
}

func True(atm *int32) bool {
	return atomic.LoadInt32(atm) == 1
}

func False(atm *int32) bool {
	return atomic.LoadInt32(atm) == 0
}

func SetFalse(atm *int32) {
	atomic.StoreInt32(atm, 0)
}

func SetTrue(atm *int32) {
	atomic.StoreInt32(atm, 1)
}
