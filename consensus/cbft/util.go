package cbft

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
	"io"
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

func combineBytes(bts ...[]byte) []byte {
	buffer := bytes.NewBuffer(make([]byte, 0, 128))
	for _, v := range bts {
		io.Copy(buffer, bytes.NewReader(v))
	}
	temp := buffer.Bytes()
	length := len(temp)
	var response []byte
	//are we wasting more than 10% space?
	if cap(temp) > (length + length / 10) {
		response = make([]byte, length)
		copy(response, temp)
	} else {
		response = temp
	}
	return response
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

/*func recoverAddr(msg ConsensusMsg) (common.Address, error) {
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
}*/