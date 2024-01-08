// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package state

import (
	"bytes"
	"github.com/PlatONnetwork/PlatON-Go/core/types"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func BenchmarkCutOriginal(b *testing.B) {
	value := common.HexToHash("0x01")
	for i := 0; i < b.N; i++ {
		bytes.TrimLeft(value[:], "\x00")
	}
}

func BenchmarkCutsetterFn(b *testing.B) {
	value := common.HexToHash("0x01")
	cutSetFn := func(r rune) bool { return r == 0 }
	for i := 0; i < b.N; i++ {
		bytes.TrimLeftFunc(value[:], cutSetFn)
	}
}

func BenchmarkCutCustomTrim(b *testing.B) {
	value := common.HexToHash("0x01")
	for i := 0; i < b.N; i++ {
		common.TrimLeftZeroes(value[:])
	}
}

func TestStateObject(t *testing.T) {
	x := types.StateAccount{
		Root: common.HexToHash("0x1000000000000000000000000000000000000001"),
	}
	x2 := newObject(nil, common.HexToAddress("0x1000000000000000000000000000000000000001"), x)
	x3 := x2.deepCopy(nil)
	x2.data.Root = common.HexToHash("0x1000000000000000000000000000000000000012")
	t.Log(x2.data.Root.String())
	t.Log(x3.data.Root.String())
}

func TestStateObjectValuePrefix(t *testing.T) {
	hash := common.HexToHash("0x1000000000000000000000000000000000000001")
	addr := common.HexToAddress("0x1000000000000000000000000000000000000001")
	key := []byte("key")
	value := []byte("value")
	x2 := newObject(nil, addr, types.StateAccount{
		Root:             hash,
		StorageKeyPrefix: addr.Bytes(),
	})

	dbValue := x2.getPrefixValue(hash.Bytes(), key, value)
	if !bytes.Equal(value, x2.removePrefixValue(dbValue)) {
		t.Fatal("prefix error")
	}
}
