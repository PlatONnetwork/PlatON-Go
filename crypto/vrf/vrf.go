package vrf

import (
	"crypto/ecdsa"
	"errors"
)

var (
	NotSupportKey = errors.New("Unsupported key type")
)

func Prove(key *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	pk := ECP2OS(key.X, key.Y)
	sk := make([]byte, 32)
	blob := key.D.Bytes()
	copy(sk[32-len(blob):], blob)
	if pi, err := eCVRF_prove(pk, sk, data[:]); err != nil {
		return nil, err
	} else {
		return pi, nil
	}
}

func Verify(key *ecdsa.PublicKey, pi []byte, data []byte) (bool, error) {
	if res, err := eCVRF_verify(ECP2OS(key.X, key.Y), pi, data[:]); err != nil {
		return false, err
	} else {
		return res, nil
	}
}

func ProofToHash(pi []byte) []byte {
	return eCVRF_proof2hash(pi)
}