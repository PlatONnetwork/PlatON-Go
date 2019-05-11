package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"math/rand"
)

// Returns a random offset between 0 and n
func randomOffset(n int) int {
	if n == 0 {
		return 0
	}
	return int(rand.Uint32() % uint32(n))
}

func produceHash(msgType byte, hash common.Hash) common.Hash {
	hashByt := hash.Bytes()
	hashByt[0] = msgType
	hashByt[1] = 0
	hashByt[2] = 0
	hashByt[3] = 0
	result := common.Hash{}
	result.SetBytes(hashByt)
	return result
}