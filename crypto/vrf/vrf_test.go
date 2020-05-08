package vrf

import (
	"crypto/ecdsa"
	"crypto/rand"
	"io"
	"testing"
)

func TestVrf(t *testing.T) {
	for i := 0; i < 10; i++ {
		sk, err := ecdsa.GenerateKey(curve, rand.Reader)
		if nil != err {
			t.Fatal("GenerateKey fail", err)
		}
		sk2, err := ecdsa.GenerateKey(curve, rand.Reader)
		if nil != err {
			t.Fatal("GenerateKey fail", err)
		}
		data := make([]byte, 32)
		io.ReadFull(rand.Reader, data)
		pi, err := Prove(sk, data)
		if nil != err {
			t.Fatal("Generate vrf proof failed", err)
		}
		ok, err := Verify(&sk.PublicKey, pi, data)
		if nil != err || !ok {
			t.Fatal("verification failed", err)
		}
		ok2, err := Verify(&sk2.PublicKey, pi, data)
		if nil != err || ok2 {
			t.Fatal("verification failed", err)
		}
		data = append(data, []byte("message")...)
		ok3, err := Verify(&sk.PublicKey, pi, data)
		if nil != err || ok3 {
			t.Fatal("verification failed", err)
		}
	}
}
