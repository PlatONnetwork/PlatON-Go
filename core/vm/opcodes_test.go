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
)

func TestStringToOp(t *testing.T) {
	testCases := []struct {
		strOpCode  string
		wantOpCode OpCode
	}{
		{
			strOpCode:  "CALLCODE",
			wantOpCode: CALLCODE,
		}, {
			strOpCode:  "PUSH1",
			wantOpCode: PUSH1,
		}, {
			strOpCode:  "MLOAD",
			wantOpCode: MLOAD,
		},
	}
	for _, v := range testCases {
		opCode := StringToOp(v.strOpCode)
		assert.Equal(t, v.wantOpCode, opCode)
		assert.Equal(t, v.strOpCode, v.wantOpCode.String())
	}

	// test -> string
	str := OpCode(0x00).String()
	t.Log(str)

}

func TestOpCode_IsPush(t *testing.T) {
	testCases := []struct {
		opCode OpCode
		want   bool
	}{
		{opCode: CALLCODE, want: false},
		{opCode: CALLDATALOAD, want: false},
		{opCode: MLOAD, want: false},
		{opCode: SUB, want: false},
		{opCode: PUSH, want: false},
		{opCode: PUSH1, want: true},
	}
	for _, v := range testCases {
		assert.Equal(t, v.want, v.opCode.IsPush())
		assert.Equal(t, false, v.opCode.IsStaticJump())
	}
}
