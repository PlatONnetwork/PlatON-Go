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

package core

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/PlatONnetwork/PlatON-Go/common"
)

func TestInt32ToBytes(t *testing.T) {

	var i int32 = 12 // 00000000, 00000000, 00000000, 00001100
	b := Int32ToBytes(i)
	arr := []byte{0x00, 0x00, 0x00, 0x0c}
	assert.Equal(t, b, arr, fmt.Sprintf("Expect: %v", arr))

}

func TestBytesToInt32(t *testing.T) {
	arr := []byte{0x00, 0x00, 0x00, 0x0c}
	i := BytesToInt32(arr)
	assert.Equal(t, i, int32(12), fmt.Sprintf("Expect: %d", int32(12)))
}

func TestInt64ToBytes(t *testing.T) {

	var i int64 = 12
	b := Int64ToBytes(i)
	arr := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c}
	assert.Equal(t, b, arr, fmt.Sprintf("Expect: %v", arr))
}

func TestBytesToInt64(t *testing.T) {
	arr := []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0c}
	i := BytesToInt64(arr)
	assert.Equal(t, i, int64(12), fmt.Sprintf("Expect: %d", int64(12)))
}

func TestByteConvert(t *testing.T) {
	//bytes, _ := hexutil.Decode("0x0c55699c")
	hash := common.BytesToHash(Int32ToBytes(121))

	result := BytesConverter(hash.Bytes(), "int32")
	fmt.Printf("\nresult: %v\n", result)

}

func TestStringConverter(t *testing.T) {
	result, err := StringConverter("false", "bool")
	fmt.Printf("\nresult: %v\n", result)
	if err != nil {
		fmt.Printf("\nerr: %v\n", err.Error())
	}
	//buf := bytes.NewBuffer([]byte{})
	//binary.Write(buf, binary.BigEndian, "true")
	//fmt.Println(buf.Bytes())
	//fmt.Println(len(buf.Bytes()))

	//fmt.Printf("%v",i)
}
