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

package snapshotdb

import (
	"math/big"
	"testing"
)

func TestEncodeJournalKey(t *testing.T) {
	num := new(big.Int).SetUint64(111)
	key := EncodeWalKey(num)
	num2 := DecodeWalKey(key)
	if num.Cmp(num2) != 0 {
		t.Error("num should same")
	}
}
