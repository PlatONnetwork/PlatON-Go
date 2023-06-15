//go:build ((linux && amd64) || (linux && arm64) || (darwin && amd64) || (darwin && arm64) || (windows && amd64)) && !blst_disabled
// +build linux,amd64 linux,arm64 darwin,amd64 darwin,arm64 windows,amd64
// +build !blst_disabled

package blst

import "C"
import (
	"crypto/sha256"
	"crypto/subtle"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/crypto/blseth/params"

	"crypto/rand"
	"github.com/PlatONnetwork/PlatON-Go/crypto/blseth/common"
	blst "github.com/supranational/blst/bindings/go"
)

func G1() *blst.P1 {
	return blst.P1Generator()
}

// bls12SecretKey used in the BLS signature scheme.
type bls12SecretKey struct {
	p *blst.SecretKey
}

func NewEmptyKey() common.SecretKey {
	var p blst.SecretKey
	return &bls12SecretKey{
		p: &p,
	}
}

// RandKey creates a new private key using a random method provided as an io.Reader.
func RandKey() (common.SecretKey, error) {
	// Generate 32 bytes of randomness
	var ikm [32]byte
	_, err := rand.Read(ikm[:])
	if err != nil {
		return nil, err
	}
	// Defensive check, that we have not generated a secret key,
	secKey := &bls12SecretKey{blst.KeyGen(ikm[:])}
	if IsZero(secKey.Marshal()) {
		return nil, common.ErrZeroKey
	}
	return secKey, nil
}

// SecretKeyFromBytes creates a BLS private key from a BigEndian byte slice.
func SecretKeyFromBytes(privKey []byte) (common.SecretKey, error) {
	if len(privKey) != params.BLSSecretKeyLength {
		return nil, fmt.Errorf("secret key must be %d bytes", params.BLSSecretKeyLength)
	}
	secKey := new(blst.SecretKey).Deserialize(privKey)
	if secKey == nil {
		return nil, common.ErrSecretUnmarshal
	}
	wrappedKey := &bls12SecretKey{p: secKey}
	if IsZero(privKey) {
		return nil, common.ErrZeroKey
	}
	return wrappedKey, nil
}

// PublicKey obtains the public key corresponding to the BLS secret key.
func (s *bls12SecretKey) PublicKey() common.PublicKey {
	return &PublicKey{p: new(blstPublicKey).From(s.p)}
}

// IsZero checks if the secret key is a zero key.
func IsZero(sKey []byte) bool {
	b := byte(0)
	for _, s := range sKey {
		b |= s
	}
	return subtle.ConstantTimeByteEq(b, 0) == 1
}

// Sign a message using a secret key - in a beacon/validator client.
//
// In IETF draft BLS specification:
// Sign(SK, message) -> signature: a signing algorithm that generates
//
//	a deterministic signature given a secret key SK and a message.
//
// In Ethereum proof of stake specification:
// def Sign(SK: int, message: Bytes) -> BLSSignature
func (s *bls12SecretKey) Sign(msg []byte) common.Signature {
	signature := new(blstSignature).Sign(s.p, msg, dst)
	return &Signature{s: signature}
}

// Marshal a secret key into a LittleEndian byte slice.
func (s *bls12SecretKey) Marshal() []byte {
	keyBytes := s.p.Serialize()
	return keyBytes
}
func (s *bls12SecretKey) SetLittleEndian(buf []byte) error {
	if sec := s.p.DeserializeLittleEndian(buf); sec != nil {
		return errors.New("deserialize failed")
	}
	return nil
}
func (s *bls12SecretKey) GetLittleEndian() []byte {
	return s.p.SerializeLittleEndian()
}

/*
$s$ $P$ $k$ $V$ $G2$ $c$
$h = H(G2, P, V)$
$c = MapToCurve(h)$
$r = k - s * c$
$C = c, R= r$l
*/
func (s *bls12SecretKey) MakeSchnorrNIZKP() common.SchnorrProof {
	P := P1AffineToP1(new(blstPublicKey).From(s.p))

	k, _ := RandKey()
	V := P1AffineToP1(new(blstPublicKey).From(k.(*bls12SecretKey).p))
	g1 := blst.P1Generator()

	hash := sha256.New()
	hash.Write(P1ToP1Affine(g1).Compress())
	hash.Write(P1ToP1Affine(P).Compress())
	hash.Write(P1ToP1Affine(V).Compress())
	h := hash.Sum(nil)
	c := blst.HashToScalar(h, dst)
	r, _ := s.p.Mul(c)
	r, _ = k.(*bls12SecretKey).p.Sub(r)
	return &SchnorrProof{
		C: &bls12SecretKey{c},
		R: &bls12SecretKey{r},
	}
}

func P1AffineToP1(p *blst.P1Affine) *blst.P1 {
	var p1 blst.P1
	p1.FromAffine(p)
	return &p1
}

func P1ToP1Affine(p *blst.P1) *blst.P1Affine {
	return p.ToAffine()
}
