package blst

import (
	"crypto/sha256"
	"errors"
	"github.com/PlatONnetwork/PlatON-Go/crypto/blseth/common"
	blst "github.com/supranational/blst/bindings/go"
)

type SchnorrProof struct {
	C, R *bls12SecretKey
}

func (s *SchnorrProof) Marshal() []byte {
	return append(s.C.Marshal(), s.R.Marshal()...)
}
func (s *SchnorrProof) Unmarshal(buf []byte) error {
	if len(buf)%2 != 0 {
		return errors.New("the length of C and R not equal in proof")
	}

	pivot := len(buf) / 2
	c, err := SecretKeyFromBytes(buf[:pivot])
	if err != nil {
		return err
	}
	r, err := SecretKeyFromBytes(buf[pivot:])
	if err != nil {
		return err
	}
	s.C = c.(*bls12SecretKey)
	s.R = r.(*bls12SecretKey)
	return nil
}

/*
$P_r = G*r$
$P_c = P*c$
$V_1 = P_r + P_c$
$h = H(G2, P, V_1)$
$c_1 = MapToScalar(h)$
$c_1 == c$
*/
func (s *SchnorrProof) Verify(pk common.PublicKey) error {
	g1 := blst.P1Generator()
	P := pk.(*PublicKey).p
	Pr := g1.Mult(s.R.p)
	Pc := P1AffineToP1(P).Mult(s.C.p)
	V := Pr.Add(Pc)

	hash := sha256.New()
	hash.Write(P1ToP1Affine(g1).Compress())
	hash.Write(P.Compress())
	hash.Write(P1ToP1Affine(V).Compress())
	h := hash.Sum(nil)
	c1 := blst.HashToScalar(h, dst)
	if !c1.Equals(s.C.p) {
		return errors.New("verify failed")
	}
	return nil
}
