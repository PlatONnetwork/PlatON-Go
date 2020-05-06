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
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func TestValidJumpdest(t *testing.T) {
	code := []byte{0x00, 0x5b}
	contract := &Contract{
		Code:      code,
		CodeHash:  common.BytesToHash(code),
		jumpdests: make(map[common.Hash]bitvec),
	}
	r := contract.validJumpdest(new(big.Int).SetUint64(3))
	if r {
		t.Errorf("Expected false, got true")
	}
	r = contract.validJumpdest(new(big.Int).SetUint64(1))
	if !r {
		t.Errorf("Expected true, got false")
	}
	r = contract.validJumpdest(new(big.Int).SetUint64(2))
	if r {
		t.Errorf("Expected false, got true")
	}
}

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
		t.Errorf("Not equal, expect: %s, actual: %s", contract.caller.Address(), c.CallerAddress)
	}
}

func TestGetOp(t *testing.T) {
	code := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	testCases := []struct {
		n    uint64
		want OpCode
	}{
		{n: 0, want: STOP},
		{n: 1, want: ADD},
		{n: 2, want: MUL},
		{n: 3, want: SUB},
		{n: 4, want: DIV},
		{n: 5, want: SDIV},
		{n: 6, want: MOD},
	}
	c := &Contract{
		Code: code,
	}
	// iterate and verify.
	for _, v := range testCases {
		opCode := c.GetOp(v.n)
		assert.Equal(t, v.want, opCode)
	}
}

func TestGetByte(t *testing.T) {
	code := []byte{0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06}
	testCases := []struct {
		n    uint64
		want byte
	}{
		{n: 0, want: byte(STOP)},
		{n: 1, want: byte(ADD)},
		{n: 2, want: byte(MUL)},
		{n: 3, want: byte(SUB)},
		{n: 4, want: byte(DIV)},
		{n: 5, want: byte(SDIV)},
		{n: 6, want: byte(MOD)},
		{n: 100, want: byte(0x00)},
	}
	c := &Contract{
		Code: code,
	}
	// iterate and verify.
	for _, v := range testCases {
		r := c.GetByte(v.n)
		assert.Equal(t, v.want, r)
	}
}

func TestCaller(t *testing.T) {
	addr := common.BytesToAddress([]byte("aaa"))
	contract := &Contract{
		CallerAddress: addr,
	}
	cr := contract.Caller()
	if cr != addr {
		t.Errorf("Not equal, expect: %s, actual: %s", addr, cr)
	}
}

func TestUseGas(t *testing.T) {
	contract := &Contract{
		Gas: 1000,
	}
	cr := contract.UseGas(100)
	if !cr {
		t.Errorf("Expected: true, got false")
	}
	laveGas := contract.Gas - 100
	if laveGas != 800 {
		t.Errorf("Expected: 800, actual: %d", laveGas)
	}

	// Simulation does not hold.
	cr = contract.UseGas(1000)
	if cr {
		t.Errorf("Expected: false, got true")
	}
}

func TestAddress(t *testing.T) {
	addr := common.BytesToAddress([]byte("aaa"))
	contract := &Contract{
		self: &MockAddressRef{},
	}
	cr := contract.Address()
	if cr != addr {
		t.Errorf("Expected: %s, got: %s", addr, cr)
	}
}

func TestValue(t *testing.T) {
	contract := &Contract{
		value: buildBigInt(100),
	}
	cr := contract.Value()
	if cr.Cmp(new(big.Int).SetUint64(100)) != 0 {
		t.Errorf("Expected: 100, got: %d", cr)
	}
}

func TestSetting(t *testing.T) {
	// test SetCallCode of method.
	contract := &Contract{
		value: buildBigInt(100),
	}
	addr := common.BytesToAddress([]byte("I'm address"))
	code := []byte{0x00, 0x10}
	hash := common.BytesToHash(code)
	contract.SetCallCode(&addr, hash, code)

	if *contract.CodeAddr != addr {
		t.Errorf("Expected: %s, got %s", addr, contract.CodeAddr)
	}
	assert.Equal(t, code, contract.Code)
	if contract.CodeHash != hash {
		t.Errorf("Expected: %s, got %s", hash, contract.CodeHash)
	}

	// test SetCodeOptionalHash
	optional := codeAndHash{
		code: code,
		hash: hash,
	}
	contract.SetCodeOptionalHash(&addr, &optional)
	if *contract.CodeAddr != addr {
		t.Errorf("Expected: %s, got %s", addr, contract.CodeAddr)
	}
	assert.Equal(t, code, contract.Code)
	if contract.CodeHash != hash {
		t.Errorf("Expected: %s, got %s", hash, contract.CodeHash)
	}

	// test SetCallAbi.
	contract.SetCallAbi(&addr, hash, []byte{})

}
