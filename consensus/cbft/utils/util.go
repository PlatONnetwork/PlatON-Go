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

// RandomOffset returns a random offset between 0 and n
func RandomOffset(n int) int {
	if n == 0 {
		return 0
	}
	return int(rand.Uint32() % uint32(n))
}

// BuildHash converts byte array to hash. Use sha256 to
// generate a unique message hash.
func BuildHash(msgType byte, bytes []byte) common.Hash {
	bytes[0] = msgType
	hashBytes := sha3.Sum256(bytes)
	result := common.Hash{}
	result.SetBytes(hashBytes[:])
	return result
}

// MergeBytes merges multiple bytes of data and
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

// Returns whether the specified value is equal to 1.
func True(atm *int32) bool {
	return atomic.LoadInt32(atm) == 1
}

// Returns whether the specified value is equal to 0.
func False(atm *int32) bool {
	return atomic.LoadInt32(atm) == 0
}

// Set the specified variable to 0.
func SetFalse(atm *int32) {
	atomic.StoreInt32(atm, 0)
}

// Set the specified variable to 1.
func SetTrue(atm *int32) {
	atomic.StoreInt32(atm, 1)
}

// Represents a k-v key-value pair.
type KeyValuePair struct {
	Key   string
	Value int64
}

// KeyValuePairList is a slice of Pairs that implements
// sort.Interface to sort by Value.
type KeyValuePairList []KeyValuePair

func (kvp KeyValuePairList) Swap(i, j int)      { kvp[i], kvp[j] = kvp[j], kvp[i] }
func (kvp KeyValuePairList) Len() int           { return len(kvp) }
func (kvp KeyValuePairList) Less(i, j int) bool { return kvp[i].Value < kvp[j].Value }

// Add an element.
func (kvp *KeyValuePairList) Push(x interface{}) {
	*kvp = append(*kvp, x.(KeyValuePair))
}

// Pop up an element.
func (kvp *KeyValuePairList) Pop() interface{} {
	old := *kvp
	n := len(old)
	x := old[n-1]
	*kvp = old[0 : n-1]
	return x
}

// Sort the Map according to the key,
// return the list of sorted key-value pairs.
func SortMap(m map[string]int64) KeyValuePairList {
	p := make(KeyValuePairList, len(m))
	i := 0
	for k, v := range m {
		p[i] = KeyValuePair{k, v}
		i++
	}
	sort.Sort(p)
	return p
}
