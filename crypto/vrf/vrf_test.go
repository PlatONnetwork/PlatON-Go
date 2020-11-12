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
