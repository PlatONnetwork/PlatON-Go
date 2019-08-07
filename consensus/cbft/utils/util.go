package utils

import (
	"bytes"
	"io"
	"math/rand"
	"sort"
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

type KeyValuePair struct {
	Key   string
	Value int64
}

// A slice of Pairs that implements sort.Interface to sort by Value.
type KeyValuePairList []KeyValuePair

func (p KeyValuePairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
func (p KeyValuePairList) Len() int           { return len(p) }
func (p KeyValuePairList) Less(i, j int) bool { return p[i].Value < p[j].Value }

func (h *KeyValuePairList) Push(x interface{}) {
	*h = append(*h, x.(KeyValuePair))
}

func (h *KeyValuePairList) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func SortMap(m map[string]int64) KeyValuePairList {
	p := make(KeyValuePairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = KeyValuePair{k, v}
		i += 1
	}
	sort.Sort(p)
	return p
}
