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

package vm

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func TestSet(t *testing.T) {
	m := NewMemory()
	testCases := []struct {
		value  []byte
		offset uint64
		size   uint64
		want   string
	}{
		{[]byte{0x00}, 0, 1, "00"},
		{[]byte{0x01}, 0, 1, "01"},
		{[]byte{0x00, 0x01, 0x02}, 0, 3, "000102"},
	}
	for _, v := range testCases {
		m.Resize(v.size)
		m.Set(v.offset, v.size, v.value)
		actual := common.Bytes2Hex(m.store)
		if actual != v.want {
			t.Errorf("Expected: %s, got: %s", v.want, actual)
		}
	}
}

func TestSet32(t *testing.T) {
	m := NewMemory()
	testCases := []struct {
		val    common.Hash
		offset uint64
		want   common.Hash
	}{
		{common.BytesToHash([]byte{0x11}), 0, common.BytesToHash([]byte{0x11})},
		{common.BytesToHash([]byte{0x00}), 0, common.BytesToHash([]byte{0x00})},
		{common.BytesToHash([]byte{0x11, 0xab}), 0, common.BytesToHash([]byte{0x11, 0xab})},
	}
	for _, v := range testCases {
		m.Resize(32)
		m.Set32(v.offset, v.val.Big())
		actual := common.Bytes2Hex(m.Data())
		if actual != v.want.HexWithNoPrefix() {
			t.Errorf("Expected: %s, got: %s", v.want.Hex(), actual)
		}
	}
}

func TestResize(t *testing.T) {
	m := NewMemory()
	testCases := []struct {
		size int64
	}{
		{size: 10},
		{size: 1000},
		{size: 2000},
	}
	for _, v := range testCases {
		m.Resize(uint64(v.size))
		assert.Equal(t, int(v.size), m.Len())
	}
}

func TestGet(t *testing.T) {
	m := NewMemory()
	testCases := []struct {
		value  []byte
		offset uint64
		size   uint64
		want   string
	}{
		{[]byte{0x00}, 0, 0, ""},
		{[]byte{0x00}, 0, 1, "00"},
		{[]byte{0x01}, 0, 1, "01"},
		{[]byte{0x00, 0x01, 0x02}, 0, 3, "000102"},
	}
	for _, v := range testCases {
		m.Resize(v.size)
		m.Set(v.offset, v.size, v.value)
		actual := common.Bytes2Hex(m.Get(int64(v.offset), int64(v.size)))
		if actual != v.want {
			t.Errorf("Expected: %s, got: %s", v.want, actual)
		}
	}
}

func TestGetPtr(t *testing.T) {
	m := NewMemory()
	testCases := []struct {
		value  []byte
		offset uint64
		size   uint64
		want   string
	}{
		{[]byte{0x00}, 0, 0, ""},
		{[]byte{0x00}, 0, 1, "00"},
		{[]byte{0x01}, 0, 1, "01"},
		{[]byte{0x00, 0x01, 0x02}, 0, 3, "000102"},
	}
	for _, v := range testCases {
		m.Resize(v.size)
		m.Set(v.offset, v.size, v.value)
		actual := common.Bytes2Hex(m.GetPtr(int64(v.offset), int64(v.size)))
		if actual != v.want {
			t.Errorf("Expected: %s, got: %s", v.want, actual)
		}
		m.Print()
	}
}
