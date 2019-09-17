package bls

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

// Match only 128 hex char length proof
type SchnorrProofHex [64]byte

func (pfe SchnorrProofHex) String() string {
	return hex.EncodeToString(pfe[:])
}

// MarshalText implements the encoding.TextMarshaler interface.
func (pfe SchnorrProofHex) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(pfe[:])), nil
}

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (pfe *SchnorrProofHex) UnmarshalText(text []byte) error {

	var p SchnorrProofHex
	b, err := hex.DecodeString(strings.TrimPrefix(string(text), "0x"))
	if err != nil {
		return err
	} else if len(b) != len(p) {
		return fmt.Errorf("wrong length, want %d hex chars", len(p)*2)
	}
	copy(p[:], b)

	*pfe = p
	return nil
}

type SchnorrProof struct {
	C, R SecretKey
}

// Serialize --
func (pf *SchnorrProof) Serialize() []byte {
	return append(pf.C.Serialize(), (pf.R.Serialize())...)

}

// Deserialize --
func (pf *SchnorrProof) Deserialize(buf []byte) error {
	if len(buf)%2 != 0 {
		return errors.New("the length of C and R not equal in proof")
	}

	pivot := len(buf) / 2

	pf.C.Deserialize(buf[:pivot])
	pf.R.Deserialize(buf[pivot:])
	return nil
}

func (pf *SchnorrProof) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%x", pf.Serialize())), nil
}

func (pf *SchnorrProof) UnmarshalText(text []byte) error {
	key, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	return pf.Deserialize(key)
}

func (pf *SchnorrProof) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, pf.Serialize())
}

func (pf *SchnorrProof) DecodeRLP(s *rlp.Stream) error {
	buf, err := s.Bytes()
	if err != nil {
		return err
	}
	return pf.Deserialize(buf)
}

func (sec *SecretKey) MakeSchnorrNIZKP() (*SchnorrProof, error) {

	P := sec.GetPublicKey()
	var sk SecretKey
	sk.SetByCSPRNG()
	V := sk.GetPublicKey()
	G := GetGeneratorOfG2()
	input1 := G.Serialize()
	input2 := P.Serialize()
	input3 := V.Serialize()
	var buffer bytes.Buffer
	buffer.Write(input1)
	buffer.Write(input2)
	buffer.Write(input3)
	output := buffer.Bytes()
	h := crypto.Keccak256(output)
	var c SecretKey
	err := c.SetLittleEndian(h)
	if err != nil {
		return nil, err
	}
	temp := *sec
	temp.Mul(&c)
	r := sk
	r.Sub(&temp)
	sig := new(SchnorrProof)
	sig.C = c
	sig.R = r
	return sig, nil
}

func (sig *SchnorrProof) VerifySchnorrNIZK(pk PublicKey) error {

	if !G2IsValid(&pk) {
		return errors.New("P isnot valid")
	}
	c := sig.C
	r := sig.R
	G := GetGeneratorOfG2()
	//V1 = G * r + A * c     c = H(G || pk || V’)
	var Pr PublicKey
	Pr = *G
	Pr.Mul(&r)
	Pc := pk
	Pc.Mul(&c)
	V1 := Pr
	V1.Add(&Pc)
	input1 := G.Serialize()
	input2 := pk.Serialize()
	input3 := V1.Serialize()
	var buffer bytes.Buffer
	buffer.Write(input1)
	buffer.Write(input2)
	buffer.Write(input3)
	output := buffer.Bytes()
	h := crypto.Keccak256(output)
	var c1 SecretKey
	err := c1.SetLittleEndian(h)
	if err != nil {
		return err
	}
	if !c.IsEqual(&c1) {
		return errors.New("not same c = H(G || pk || V’)")
	}
	return nil

}
