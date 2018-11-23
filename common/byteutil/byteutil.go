// create by platon
package byteutil

import (
	"Platon-go/common"
	"bytes"
	"encoding/binary"
	"math/big"
)

var Command = map[string] interface{} {
	"string" : BytesToString,
	"[]uint8" : OriginBytes,
	"[64]uint8" : BytesTo64Bytes,
	"[32]uint8" : BytesTo32Bytes,
	"int" : BytesToInt,
	"*big.Int" : BytesToBigInt,
	"uint32" : binary.BigEndian.Uint32,
	"uint64" : binary.BigEndian.Uint64,
	"int32" : common.BytesToInt32,
	"int64" : common.BytesToInt64,
	"float32" : common.BytesToFloat32,
	"float64" : common.BytesToFloat64,
	"discover.NodeID" : BytesTo64Bytes,
	"common.Hash": common.BytesToHash,
	"common.Address" : common.BytesToAddress,
}

func BytesTo32Bytes(curByte []byte) [32]byte {
	var arr [32]byte
	copy(arr[:], curByte)
	return arr
}

func BytesTo64Bytes(curByte []byte) [64]byte {
	var arr [64]byte
	copy(arr[:], curByte)
	return arr
}

func OriginBytes(curByte []byte) []byte {
	return curByte
}

func BytesToBigInt(curByte []byte) *big.Int {
	big1 := new(big.Int).SetInt64(common.BytesToInt64(curByte))
	return big1
}

func BytesToInt(curByte []byte) int {
	bytesBuffer := bytes.NewBuffer(curByte)
	var x int32
	binary.Read(bytesBuffer, binary.BigEndian, &x)
	b := int(x)
	return b
}

func BytesToString(curByte []byte) string {
	return string(curByte)
}

func StringToBytes(curStr string) []byte {
	return []byte(curStr)
}

func BoolToBytes(val bool) []byte {
	buf := bytes.NewBuffer([]byte{})
	binary.Write(buf, binary.BigEndian, true)
	return buf.Bytes()
}

func IntToBytes(curInt int) []byte {
	x := int32(curInt)
	bytesBuffer := bytes.NewBuffer([]byte{})
	binary.Write(bytesBuffer, binary.BigEndian, &x)
	return bytesBuffer.Bytes()
}

func Uint64ToBytes(val uint64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf, val)
	return buf[:]
}