package utils

import "unsafe"

func String2bytes(str string) []byte {
	x := (*[2]uintptr)(unsafe.Pointer(&str))
	h := [3]uintptr{x[0], x[1], x[1]}
	return *(*[]byte)(unsafe.Pointer(&h))
}

func Bytes2string(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}