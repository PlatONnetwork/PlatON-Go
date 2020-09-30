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

	"github.com/PlatONnetwork/PlatON-Go/common/byteutil"
)

func TestMemorySha3(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x01}))
	stack.push(byteutil.BytesToBigInt([]byte{0x02}))
	r := memorySha3(stack)
	if r.Uint64() != 3 {
		t.Errorf("Expected: 3, got %d", r.Uint64())
	}
}

func TestMemoryCallDataCopy(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x01}))
	stack.push(byteutil.BytesToBigInt([]byte{0x02}))
	stack.push(byteutil.BytesToBigInt([]byte{0x03}))
	r := memoryCallDataCopy(stack)
	if r.Uint64() != 4 {
		t.Errorf("Expected: 4, got %d", r.Uint64())
	}
}

func TestMemoryReturnDataCopy(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x01}))
	stack.push(byteutil.BytesToBigInt([]byte{0x02}))
	stack.push(byteutil.BytesToBigInt([]byte{0x03}))
	stack.push(byteutil.BytesToBigInt([]byte{0x05}))
	r := memoryReturnDataCopy(stack)
	if r.Uint64() != 7 {
		t.Errorf("Expected: 7, got %d", r.Uint64())
	}
}

func TestMemoryCodeCopy(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x01}))
	stack.push(byteutil.BytesToBigInt([]byte{0x02}))
	stack.push(byteutil.BytesToBigInt([]byte{0x03}))
	stack.push(byteutil.BytesToBigInt([]byte{0x05}))
	stack.push(byteutil.BytesToBigInt([]byte{0x06}))
	stack.push(byteutil.BytesToBigInt([]byte{0x07}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	r := memoryCodeCopy(stack)
	if r.Uint64() != 14 {
		t.Errorf("Expected: 14, got %d", r.Uint64())
	}
}

func TestMemoryExtCodeCopy(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x01}))
	stack.push(byteutil.BytesToBigInt([]byte{0x02}))
	stack.push(byteutil.BytesToBigInt([]byte{0x03}))
	stack.push(byteutil.BytesToBigInt([]byte{0x05}))
	stack.push(byteutil.BytesToBigInt([]byte{0x06}))
	stack.push(byteutil.BytesToBigInt([]byte{0x07}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	r := memoryExtCodeCopy(stack)
	if r.Uint64() != 12 {
		t.Errorf("Expected: 12, got %d", r.Uint64())
	}
}

func TestMemoryMLoad(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	r := memoryMLoad(stack)
	if r.Uint64() != 40 {
		t.Errorf("Expected: 40, got %d", r.Uint64())
	}
}

func TestMemoryMStore8(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	r := memoryMStore8(stack)
	if r.Uint64() != 9 {
		t.Errorf("Expected: 9, got %d", r.Uint64())
	}
}

func TestMemoryMStore(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	r := memoryMStore(stack)
	if r.Uint64() != 40 {
		t.Errorf("Expected: 40, got %d", r.Uint64())
	}
}

func TestMemoryCreate(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	stack.push(byteutil.BytesToBigInt([]byte{0x03}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	r := memoryCreate(stack)
	if r.Uint64() != 11 {
		t.Errorf("Expected: 11, got %d", r.Uint64())
	}
	r = memoryCreate(stack)
	if r.Uint64() != 11 {
		t.Errorf("Expected: 11, got %d", r.Uint64())
	}
}

func TestMemoryCall(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	stack.push(byteutil.BytesToBigInt([]byte{0x06}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	stack.push(byteutil.BytesToBigInt([]byte{0x03}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	r := memoryCall(stack)
	if r.Uint64() != 16 {
		t.Errorf("Expected: 16, got %d", r.Uint64())
	}

	// memoryDelegateCall verify.
	r = memoryDelegateCall(stack)
	if r.Uint64() != 16 {
		t.Errorf("Expected: 16, got %d", r.Uint64())
	}

	// memoryStaticCall verify.
	r = memoryDelegateCall(stack)
	if r.Uint64() != 16 {
		t.Errorf("Expected: 16, got %d", r.Uint64())
	}
}

func TestMemoryReturn(t *testing.T) {
	stack := newstack()
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	stack.push(byteutil.BytesToBigInt([]byte{0x03}))
	stack.push(byteutil.BytesToBigInt([]byte{0x08}))
	r := memoryReturn(stack)
	if r.Uint64() != 11 {
		t.Errorf("Expected: 11, got %d", r.Uint64())
	}

	// for memoryRevert.
	r = memoryRevert(stack)
	if r.Uint64() != 11 {
		t.Errorf("Expected: 11, got %d", r.Uint64())
	}

	// for memoryLog.
	r = memoryLog(stack)
	if r.Uint64() != 11 {
		t.Errorf("Expected: 11, got %d", r.Uint64())
	}
}
