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
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasCallDataCopy(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 113 {
		t.Errorf("Expected: 113, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}

func TestGasReturnDataCopy(t *testing.T) {
	gasTable := params.GasTableConstantinople
	stack := newstack()
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	stack.push(new(big.Int).SetUint64(100))
	gas, err := gasReturnDataCopy(gasTable, &EVM{}, &Contract{}, stack, NewMemory(), 1024)
	if gas != 113 {
		t.Errorf("Expected: 113, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
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
		t.Errorf("Expected: 2773, got %d", gas)
	}
	if err != nil {
		t.Error("not expected error")
	}
}
