package rfc6979

import (
	"bytes"
	"crypto/hmac"
	"hash"
	"math/big"
	"crypto/sha256"
	"github.com/PlatONnetwork/PlatON-Go/crypto/secp256k1"
	"crypto/ecdsa"
)

const (
	tryLimit = 10000
)
// mac returns an HMAC of the given key and message.
func hmac_k(alg func() hash.Hash, k, m[]byte) []byte {
	h := hmac.New(alg, k)
	h.Write(m)
	return h.Sum(nil)
}

// https://tools.ietf.org/html/rfc6979#section-2.3.2
func bits2int(in []byte, qlen int) *big.Int {
	vlen := len(in) * 8
	v := new(big.Int).SetBytes(in)
	if vlen > qlen {
		v = new(big.Int).Rsh(v, uint(vlen-qlen))
	}
	return v
}

// https://tools.ietf.org/html/rfc6979#section-2.3.3
func int2octets(v *big.Int, rlen int) []byte {
	out := v.Bytes()

	// pad with zeros if it's too short
	if len(out) < rlen {
		out2 := make([]byte, rlen)
		copy(out2[rlen-len(out):], out)
		return out2
	}
	// drop most significant bytes if it's too long
	if len(out) > rlen {
		out2 := make([]byte, rlen)
		copy(out2, out[len(out)-rlen:])
		return out2
	}

	return out
}

// https://tools.ietf.org/html/rfc6979#section-2.3.4
func bits2octets(in []byte, q *big.Int, qlen, rlen int) []byte {
	z1 := bits2int(in, qlen)
	z2 := new(big.Int).Sub(z1, q)
	if z2.Sign() < 0 {
		return int2octets(z1, rlen)
	}
	return int2octets(z2, rlen)
}

var one = big.NewInt(1)

// https://tools.ietf.org/html/rfc6979#section-3.2
func generate_k(q, x *big.Int, alg func() hash.Hash, hash []byte, test func(*big.Int) bool) {
	qlen := q.BitLen()
	hlen := alg().Size()
	rlen := (qlen + 7) >> 3
	bx := append(int2octets(x, rlen), bits2octets(hash, q, qlen, rlen)...)
	// Step B
	v := bytes.Repeat([]byte{0x01}, hlen)
	// Step C
	k := bytes.Repeat([]byte{0x00}, hlen)
	// Step D
	k = hmac_k(alg, k, append(append(v, 0x00), bx...))
	// Step E
	v = hmac_k(alg, k, v)
	// Step F
	k = hmac_k(alg, k, append(append(v, 0x01), bx...))
	// Step G
	v = hmac_k(alg, k, v)
	// Step H
	for i := int64(0); i < tryLimit; i++ {
		// Step H1
		var t []byte
		// Step H2
		for len(t) < qlen/8 {
			v = hmac_k(alg, k, v)
			t = append(t, v...)
		}
		// Step H3
		secret := bits2int(t, qlen)
		if secret.Cmp(one) >= 0 && secret.Cmp(q) < 0 && test(secret) {
			return
		}
		k = hmac_k(alg, k, append(v, 0x00))
		v = hmac_k(alg, k, v)
	}
	panic("generate_k: couldn't generate a new k")
}

func ECVRF_nonce_generation(sk []byte,m []byte)(*ecdsa.PrivateKey, error){
	curve := secp256k1.S256()

	hash := sha256.New()
	hash.Write(m)
	h := hash.Sum(nil)

	var sec *big.Int
	generate_k(curve.N,new(big.Int).SetBytes(sk), sha256.New, h, func(k *big.Int) bool {
		sec = k
		return true
	})
	priv := new(ecdsa.PrivateKey)
	priv.PublicKey.Curve = curve
	priv.D = sec
	priv.PublicKey.X, priv.PublicKey.Y = curve.ScalarBaseMult(sec.Bytes())

	return priv,nil
}

