// create by platon
package byteutil

import (
	"Platon-go/common"
	"bytes"
	"encoding/binary"
	"fmt"
	"math/big"
	"reflect"
)

// map 封装所有的方法
var Command = map[string] interface{} {
	"string" : BytesToString,
	"int" : BytesToInt,
	"Int" : BytesToBigInt,
	"uint32" : binary.LittleEndian.Uint32,
	"uint64" : binary.LittleEndian.Uint64,
	"int32" : common.BytesToInt32,
	"int64" : common.BytesToInt64,
	"float32" : common.BytesToFloat32,
	"float64" : common.BytesToFloat64,
	"Hash": common.BytesToHash,
	"Address" : common.BytesToAddress,
}

func Uint64ToBytes(val uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return buf[:]
}

func BoolToBytes(val bool) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, true)
	return buf.Bytes()
}

func BytesToBigInt(curByte []byte) interface{} {
	big1 := new(big.Int).SetInt64(BytesToInt64(curByte))
	fmt.Println(reflect.TypeOf(big1))
	return big1
}

func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

func BytesToInt(curByte []byte) int {
	bytesBuffer := bytes.NewBuffer(curByte)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	b := int(x)
	return b
}

func IntToBytes(curInt int) []byte {
	x := int32(curInt)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &x)
	return bytesBuffer.Bytes()
}

func BytesToString(curByte []byte) string {
	return string(curByte)
}

func StringToBytes(curStr string) []byte {
	return []byte(curStr)
}