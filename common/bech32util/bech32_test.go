package bech32util

import (
	"bytes"
	"crypto/sha256"
	"testing"
)

func TestEncodeAndDecode(t *testing.T) {

	sum := sha256.Sum256([]byte("hello world\n"))

	bech, err := ConvertAndEncode("shasum", sum[:])

	if err != nil {
		t.Error(err)
	}
	hrp, data, err := DecodeAndConvert(bech)

	if err != nil {
		t.Error(err)
	}
	if hrp != "shasum" {
		t.Error("Invalid hrp")
	}
	if !bytes.Equal(data, sum[:]) {
		t.Error("Invalid decode")
	}
}
