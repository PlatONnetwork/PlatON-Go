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

func TestAsDelegate(t *testing.T) {
	contract := &Contract{
		caller: &Contract{
			CallerAddress: common.BytesToAddress([]byte("aaa")),
			self:          &MockAddressRef{},
			value:         buildBigInt(1),
		},
	}
	c := contract.AsDelegate()
	if c.CallerAddress != contract.caller.Address() {
		t.Logf("Not equal, expect: %s, actual: %s", contract.caller.Address(), c.CallerAddress)
	}
}
