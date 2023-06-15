package eth

import (
	"encoding/hex"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls/types"
	"github.com/PlatONnetwork/PlatON-Go/crypto/blseth/blst"
	"github.com/PlatONnetwork/PlatON-Go/crypto/blseth/common"
	"strings"
)

type SecretKey struct {
	key common.SecretKey
}

func (s *SecretKey) Sign(m string) types.ISign {
	sign := s.key.Sign([]byte(m))
	return &Sign{
		sign,
	}
}

func (s *SecretKey) Serialize() []byte {
	return s.key.Marshal()
}

func (s *SecretKey) Deserialize(buf []byte) error {
	key, err := blst.SecretKeyFromBytes(buf)
	if err != nil {
		return err
	}
	s.key = key
	return nil
}
func (s *SecretKey) SetLittleEndian(buf []byte) error {
	s.key = blst.NewEmptyKey()
	return s.key.SetLittleEndian(buf)
}
func (s *SecretKey) GetLittleEndian() []byte {
	return s.key.GetLittleEndian()
}
func (s *SecretKey) SetByCSPRNG() {
	key, _ := blst.RandKey()
	s.key = key
}

func (s SecretKey) GetPublicKey() types.IPublicKey {
	pubKey := s.key.PublicKey()
	return &PublicKey{
		pubKey: pubKey,
	}
}

func (s SecretKey) MakeSchnorrNIZKP() (types.ISchnorrProof, error) {
	proof := s.key.MakeSchnorrNIZKP()
	return &SchnorrProof{
		proof: proof,
	}, nil
}

type PublicKey struct {
	pubKey common.PublicKey
}

func (p *PublicKey) Add(rhs types.IPublicKey) {
	pub := rhs.(*PublicKey)
	p.pubKey.Aggregate(pub.pubKey)
}

func (p *PublicKey) UnmarshalText(text []byte) error {
	key, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	pubKey, err := blst.PublicKeyFromBytes(key)
	if err != nil {
		return err
	}
	p.pubKey = pubKey
	return nil
}

func (p PublicKey) Serialize() []byte {
	return p.pubKey.Marshal()
}

func (p PublicKey) Bytes() []byte {
	return p.pubKey.Marshal()
}

func (p PublicKey) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(p.pubKey.Marshal())), nil
}

type PublicKeyHex struct {
	pubKey [48]byte
}

func (p PublicKeyHex) ParseBlsPubKey() (types.IPublicKey, error) {
	pubKey, err := blst.PublicKeyFromBytes(p.pubKey[:])
	if err != nil {
		return nil, err
	}
	return &PublicKey{
		pubKey: pubKey,
	}, nil
}

func (p PublicKeyHex) Bytes() []byte {
	return p.pubKey[:]
}

func (p *PublicKeyHex) UnmarshalText(text []byte) error {
	key, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	} else if len(key) != len(p.pubKey) {
		return fmt.Errorf("wrong length, want %d hex chars", len(p.pubKey)*2)
	}
	copy(p.pubKey[:], key)
	return nil
}

type Sign struct {
	sign common.Signature
}

//func (s Sign) Verify(pubKey common.PublicKey, msg []byte) bool {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s Sign) AggregateVerify(pubKeys []common.PublicKey, msgs [][32]byte) bool {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s Sign) FastAggregateVerify(pubKeys []common.PublicKey, msg [32]byte) bool {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s Sign) Eth2FastAggregateVerify(pubKeys []common.PublicKey, msg [32]byte) bool {
//	//TODO implement me
//	panic("implement me")
//}

func (s Sign) Marshal() []byte {
	return s.sign.Marshal()
}

func (s Sign) Copy() common.Signature {
	return s.sign.Copy()
}

func (s *Sign) Add(rhs types.ISign) {
	sign := rhs.(*Sign)
	s.sign = blst.AggregateSignatures([]common.Signature{s.sign, sign.sign})
}

func (s Sign) Verify(pub types.IPublicKey, m string) bool {
	pubkey := pub.(*PublicKey)
	return s.sign.Verify(pubkey.pubKey, []byte(m))
}

func (s Sign) Serialize() []byte {
	return s.sign.Marshal()
}

func (s Sign) Deserialize(buf []byte) error {
	sign, err := blst.SignatureFromBytes(buf)
	if err != nil {
		return err
	}
	s.sign = sign
	return nil
}

type SchnorrProof struct {
	proof common.SchnorrProof
}

func (s SchnorrProof) VerifySchnorrNIZK(pk types.IPublicKey) error {
	return s.proof.Verify(pk.(*PublicKey).pubKey)
}

func (s *SchnorrProof) UnmarshalText(text []byte) error {
	key, err := hex.DecodeString(string(text))
	if err != nil {
		return err
	}
	return s.proof.Unmarshal(key)
}

func (s SchnorrProof) MarshalText() ([]byte, error) {
	return []byte(fmt.Sprintf("%x", s.proof.Marshal())), nil
}

type SchnorrProofHex struct {
	proof [64]byte
}

func (s SchnorrProofHex) MarshalText() ([]byte, error) {
	return []byte(hex.EncodeToString(s.proof[:])), nil
}

func (s *SchnorrProofHex) UnmarshalText(text []byte) error {
	b, err := hex.DecodeString(strings.TrimPrefix(string(text), "0x"))
	if err != nil {
		return err
	} else if len(b) != len(s.proof) {
		return fmt.Errorf("wrong length, want %d hex chars", len(s.proof)*2)
	}
	copy(s.proof[:], b)
	return nil
}
