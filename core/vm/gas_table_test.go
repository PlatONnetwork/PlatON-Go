// Copyright 2017 The go-ethereum Authors
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

package vm

import (
	"math/big"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/common/math"

	"github.com/PlatONnetwork/PlatON-Go/common"

	"github.com/PlatONnetwork/PlatON-Go/params"
	"github.com/stretchr/testify/assert"
)

func TestMemoryGasCost(t *testing.T) {
	//size := uint64(math.MaxUint64 - 64)
	size := uint64(0xffffffffe0)
	v, err := memoryGasCost(&Memory{}, size)
	if err != nil {
		t.Error("didn't expect error:", err)
	}
	if v != 36028899963961341 {
		t.Errorf("Expected: 36028899963961341, got %d", v)
	}

	_, err = memoryGasCost(&Memory{}, size+1)
	if err == nil {
		t.Error("expected error")
	}
}

func TestConstGasFunc(t *testing.T) {
	gas := uint64(100)
	gasFunc := constGasFunc(gas)
	gasRes, err := gasFunc(params.GasTableHomestead, &EVM{}, &Contract{}, &Stack{}, &Memory{}, 10)
	assert.Nil(t, err)
	assert.Equal(t, gas, gasRes)
}

func TestGasCallDataCopy(t *testing.T) {
	gasTable := params.GasTableConstantinople
	overUint := overUint64()
	testCases := []struct {
		elements   []*big.Int
		memorySize uint64
		expected   uint64
		isNil      bool
	}{
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(100), uint2BigInt(100)}, expected: 15, memorySize: 0, isNil: true},
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(100), uint2BigInt(100)}, expected: 113, memorySize: 1024, isNil: true},
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(100), uint2BigInt(100)}, expected: 0, memorySize: 0xffffffffe1, isNil: false},
		{elements: []*big.Int{uint2BigInt(2), overUint, uint2BigInt(3), uint2BigInt(4)}, expected: 0, memorySize: 1024, isNil: false},
	}
	for i, v := range testCases {
		stack := mockStack(v.elements...)
		gas, err := gasCallDataCopy(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), v.memorySize)
		if gas != v.expected {
			t.Errorf("Testcase %d - Expected: %d, got %d", i, v.expected, gas)
		}
		if v.isNil && err != nil {
			t.Error("not expected error")
		}
	}

}

func overUint64() *big.Int {
	data1 := new(big.Int).SetUint64(math.MaxUint64)
	data2 := new(big.Int).SetUint64(math.MaxUint64)
	res := new(big.Int)
	res.Mul(data1, data2)
	return res
}

func TestGasReturnDataCopy(t *testing.T) {
	gasTable := params.GasTableConstantinople
	overUint := overUint64()
	testCases := []struct {
		elements   []*big.Int
		memorySize uint64
		expected   uint64
		isNil      bool
	}{
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(100), uint2BigInt(100)}, expected: 15, memorySize: 0, isNil: true},
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(100), uint2BigInt(100)}, expected: 113, memorySize: 1024, isNil: true},
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(100), uint2BigInt(100)}, expected: 0, memorySize: 0xffffffffe1, isNil: false},
		{elements: []*big.Int{uint2BigInt(2), overUint, uint2BigInt(3), uint2BigInt(4)}, expected: 0, memorySize: 1024, isNil: false},
	}
	for i, v := range testCases {
		stack := mockStack(v.elements...)
		gas, err := gasReturnDataCopy(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), v.memorySize)
		if gas != v.expected {
			t.Errorf("Testcase %d - Expected: %d, got %d", i, v.expected, gas)
		}
		if v.isNil && err != nil {
			t.Error("not expected error")
		}
	}

}

func TestGasSStore(t *testing.T) {
	gasTable := params.GasTableConstantinople
	overUint := overUint64()
	contract := newContract(new(big.Int).SetUint64(0), common.BytesToAddress([]byte("a")))
	evm := &EVM{
		StateDB: createMockState(),
	}
	testCases := []struct {
		elements   []*big.Int
		memorySize uint64
		expected   uint64
		isNil      bool
	}{
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(100), uint2BigInt(100)}, expected: 20000, memorySize: 0, isNil: true},
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(0), uint2BigInt(1)}, expected: 5000, memorySize: 1024, isNil: true},
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(100), uint2BigInt(100)}, expected: 20000, memorySize: 1024, isNil: true},
		{elements: []*big.Int{uint2BigInt(100), uint2BigInt(100), uint2BigInt(100), uint2BigInt(100)}, expected: 20000, memorySize: 0xffffffffe1, isNil: false},
		{elements: []*big.Int{uint2BigInt(2), overUint, uint2BigInt(3), uint2BigInt(4)}, expected: 20000, memorySize: 1024, isNil: false},
	}
	for i, v := range testCases {
		stack := mockStack(v.elements...)
		gas, err := gasSStore(gasTable, evm, contract, stack, NewMemory(), v.memorySize)
		if gas != v.expected {
			t.Errorf("Testcase %d - Expected: %d, got %d", i, v.expected, gas)
		}
		if v.isNil && err != nil {
			t.Error("not expected error")
		}
	}
}

func uint2BigInt(u uint64) *big.Int {
	return new(big.Int).SetUint64(u)
}

func mockStack(b ...*big.Int) *Stack {
	stack := newstack()
	for _, v := range b {
		stack.push(v)
	}
	return stack
}

func TestMakeGasLog(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gasLogFunc := makeGasLog(4)
	gas, err := gasLogFunc(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 2773 {
		t.Errorf("Expected: 2773, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasSha3(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasSha3(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 152 {
		t.Errorf("Expected: 152, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasCodeCopy(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasCodeCopy(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 113 {
		t.Errorf("Expected: 113, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasExtCodeCopy(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasExtCodeCopy(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 810 {
		t.Errorf("Expected: 810, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasMLoad(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasMLoad(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 101 {
		t.Errorf("Expected: 101, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasMStore8(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasMStore8(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 101 {
		t.Errorf("Expected: 101, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasMStore(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasMStore(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 101 {
		t.Errorf("Expected: 101, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasCreate(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasCreate(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 32098 {
		t.Errorf("Expected: 32098, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasCreate2(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasCreate2(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 32122 {
		t.Errorf("Expected: 32122, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasBalance(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasBalance(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 400 {
		t.Errorf("Expected: 400, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasExtCodeSize(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasExtCodeSize(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 700 {
		t.Errorf("Expected: 700, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasSLoad(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasSLoad(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 200 {
		t.Errorf("Expected: 200, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasExp(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasExp(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 60 {
		t.Errorf("Expected: 60, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasCall(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasCall(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000}, stack, NewMemory(), 1024)
	if gas != 34898 {
		t.Errorf("Expected: 34898, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasCallCode(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasCallCode(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000}, stack, NewMemory(), 1024)
	if gas != 9898 {
		t.Errorf("Expected: 9898, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasReturn(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasReturn(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000}, stack, NewMemory(), 1024)
	if gas != 98 {
		t.Errorf("Expected: 98, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasRevert(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasRevert(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000}, stack, NewMemory(), 1024)
	if gas != 98 {
		t.Errorf("Expected: 98, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

type MockAddressRef struct{}

func (ref *MockAddressRef) Address() common.Address {
	return common.BytesToAddress([]byte("aaa"))
}

func TestGasSuicide(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasSuicide(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000, self: &MockAddressRef{}}, stack, NewMemory(), 1024)
	if gas != 5000 {
		t.Errorf("Expected: 5000, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasDelegateCall(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasDelegateCall(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000, self: &MockAddressRef{}}, stack, NewMemory(), 1024)
	if gas != 898 {
		t.Errorf("Expected: 898, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasStaticCall(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasStaticCall(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000, self: &MockAddressRef{}}, stack, NewMemory(), 1024)
	if gas != 898 {
		t.Errorf("Expected: 898, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasPush(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasPush(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000}, stack, NewMemory(), 1024)
	if gas != 3 {
		t.Errorf("Expected: 3, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasSwap(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasSwap(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000}, stack, NewMemory(), 1024)
	if gas != 3 {
		t.Errorf("Expected: 3, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasDup(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stateDB, _, _ := newChainState()

	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasDup(gasTable, &EVM{StateDB: stateDB}, &Contract{Gas: 1000}, stack, NewMemory(), 1024)
	if gas != 3 {
		t.Errorf("Expected: 3, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}
