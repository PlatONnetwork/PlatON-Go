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

package rfc6979

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/crypto/secp256k1"
	"math/big"
	"testing"
)

// https://tools.ietf.org/html/rfc6979#appendix-A.1
func TestGenerate_k(t *testing.T) {
	q, _ := new(big.Int).SetString("4000000000000000000020108A2E0CC0D99F8A5EF", 16)
	x, _ := new(big.Int).SetString("09A4D6792295A7F730FC3F2B49CBC0F62E862272F", 16)
	hash, _ := hex.DecodeString("AF2BDBE1AA9B6EC1E2ADE1D694F41FC71A831D0268E9891562113D8A62ADD1BF")
	expected, _ := new(big.Int).SetString("23AF4074C90A02B3FE61D286D5C87F425E6BDD81B", 16)
	var actual *big.Int
	generate_k(q, x, sha256.New, hash, func(k *big.Int) bool {
		actual = k
		return true
	})

	if actual.Cmp(expected) != 0 {
		t.Errorf("Expected %x, got %x", expected, actual)
	}
}

func TestDeterministicNonce(t *testing.T) {
	curve := secp256k1.S256()
	sk, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		panic(err)
	}
	fmt.Printf("sk:%x\n", sk.D.Bytes())
	msg := "hello"
	k, err := ECVRF_nonce_generation(sk.D.Bytes(), []byte(msg))
	fmt.Printf("k:%x\n", k.D.Bytes())

}
