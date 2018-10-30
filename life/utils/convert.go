package utils

import (
	"bytes"
	"encoding/binary"
	"unsafe"
)

func String2bytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}

func Int64ToBytes(i int64) []byte {
	buf := make([]byte, 8)
	binary.BigEndian.PutUint64(buf,uint64(i))
	return buf
}

func BytesToInt64(bys []byte) int64 {
	buf := bytes.NewBuffer(bys)
	var res int64
	binary.Read(buf, binary.BigEndian, &res)
	return res
}

