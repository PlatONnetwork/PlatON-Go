package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"math/rand"
)

// Returns a random offset between 0 and n
func randomOffset(n int) int {
	if n == 0 {
		return 0
	}
	return int(rand.Uint32() % uint32(n))
}

func produceHash(msgType byte, bytes []byte) common.Hash {
	bytes[0] = msgType
	hashBytes := sha3.Sum256(bytes)
	result := common.Hash{}
	result.SetBytes(hashBytes[:])
	return result
}

func uint64ToBytes(n uint64) []byte {
	return []byte{
		byte(n),
		byte(n >> 8),
		byte(n >> 16),
		byte(n >> 24),
		byte(n >> 32),
		byte(n >> 40),
		byte(n >> 48),
		byte(n >> 56),
	}
}

func recoverAddr(msg ConsensusMsg) (common.Address, error) {
	data, err := msg.CannibalizeBytes()
	recPubKey, err := crypto.Ecrecover(data, msg.Sign())
	if err != nil {
		return common.Address{}, err
	}
	pub, err := crypto.UnmarshalPubkey(recPubKey)
	if err != nil {
		return common.Address{}, err
	}
	recAddr := crypto.PubkeyToAddress(*pub)
	return recAddr, nil
}