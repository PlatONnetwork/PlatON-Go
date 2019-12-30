// Copyright 2018-2019 The PlatON Network Authors
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
