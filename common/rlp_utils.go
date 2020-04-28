package common

import "github.com/PlatONnetwork/PlatON-Go/rlp"

func MustRlpEncode(val interface{}) []byte {
	bytes, err := rlp.EncodeToBytes(val)
	if err != nil {
		panic(err)
	}
	return bytes
}

/*
func MustRlpDecode(bytes []byte, val interface{}) {
	err := rlp.DecodeBytes(bytes, val)
	if err != nil {
		panic(err)
	}
}*/
