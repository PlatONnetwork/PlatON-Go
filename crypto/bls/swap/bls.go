package swap

import (
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls/types"
	blsswap "github.com/PlatONnetwork/PlatON-Go/crypto/blsswap"
)

type SecretKey struct {
	key blsswap.SecretKey
}

func (s *SecretKey) Sign(m string) types.ISign {
	sign := s.key.Sign(m)
	return &Sign{
		sign,
	}
}

func (s *SecretKey) Serialize() []byte {
	return s.key.Serialize()
}

func (s *SecretKey) Deserialize(buf []byte) error {
	return s.key.Deserialize(buf)
}
func (s *SecretKey) SetLittleEndian(buf []byte) error {
	return s.key.SetLittleEndian(buf)
}
func (s *SecretKey) GetLittleEndian() []byte {
	return s.key.GetLittleEndian()
}
func (s *SecretKey) SetByCSPRNG() {
	s.key.SetByCSPRNG()
}

func (s SecretKey) GetPublicKey() types.IPublicKey {
	pubKey := s.key.GetPublicKey()
	return &PublicKey{
		pubKey: pubKey,
	}
}

func (s SecretKey) MakeSchnorrNIZKP() (types.ISchnorrProof, error) {
	proof, err := s.key.MakeSchnorrNIZKP()
	if err != nil {
		return nil, err
	}
	return &SchnorrProof{
		proof: proof,
	}, nil
}

type PublicKey struct {
	pubKey *blsswap.PublicKey
}

func NewPublicKey() *PublicKey {
	return &PublicKey{
		pubKey: &blsswap.PublicKey{},
	}
}
func (p *PublicKey) Add(rhs types.IPublicKey) {
	pub := rhs.(*PublicKey)
	p.pubKey.Add(pub.pubKey)
}

func (p *PublicKey) UnmarshalText(text []byte) error {
	p.pubKey.UnmarshalText(text)
	return nil
}

func (p PublicKey) Serialize() []byte {
	return p.pubKey.Serialize()
}

func (p PublicKey) Bytes() []byte {
	return p.pubKey.Serialize()
}

func (p PublicKey) MarshalText() ([]byte, error) {
	return p.pubKey.MarshalText()
}

type PublicKeyHex struct {
	pubKey blsswap.PublicKeyHex
}

func (p *PublicKeyHex) ParseBlsPubKey() (types.IPublicKey, error) {
	pubKey, err := p.pubKey.ParseBlsPubKey()
	if err != nil {
		return nil, err
	}
	return &PublicKey{
		pubKey: pubKey,
	}, nil
}

func (p PublicKeyHex) Bytes() []byte {
	return p.pubKey.Bytes()
}

func (p *PublicKeyHex) UnmarshalText(text []byte) error {
	return p.pubKey.UnmarshalText(text)
}

type Sign struct {
	sign *blsswap.Sign
}

//func (s Sign) Verify(pubKey blsswap.PublicKey, msg []byte) bool {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s Sign) AggregateVerify(pubKeys []blsswap.PublicKey, msgs [][32]byte) bool {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s Sign) FastAggregateVerify(pubKeys []blsswap.PublicKey, msg [32]byte) bool {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (s Sign) Eth2FastAggregateVerify(pubKeys []blsswap.PublicKey, msg [32]byte) bool {
//	//TODO implement me
//	panic("implement me")
//}

func (s *Sign) Add(sign types.ISign) {
	s.sign.Add(sign.(*Sign).sign)
}

func (s Sign) Verify(pub types.IPublicKey, m string) bool {
	pubkey := pub.(*PublicKey)
	return s.sign.Verify(pubkey.pubKey, m)
}

func (s Sign) Serialize() []byte {
	return s.sign.Serialize()
}

func (s *Sign) Deserialize(buf []byte) error {
	s.sign = &blsswap.Sign{}
	return s.sign.Deserialize(buf)
}

type SchnorrProof struct {
	proof *blsswap.SchnorrProof
}

func (s SchnorrProof) VerifySchnorrNIZK(pk types.IPublicKey) error {
	return s.proof.VerifySchnorrNIZK(*pk.(*PublicKey).pubKey)
}

func (s *SchnorrProof) UnmarshalText(text []byte) error {
	return s.proof.UnmarshalText(text)
}

func (s SchnorrProof) MarshalText() ([]byte, error) {
	return s.proof.MarshalText()
}
func (s *SchnorrProof) Serialize() []byte {
	return s.proof.Serialize()
}
func (s *SchnorrProof) Deserialize(buf []byte) error {
	return s.proof.Deserialize(buf)
}

type SchnorrProofHex struct {
	proof blsswap.SchnorrProofHex
}

func (s SchnorrProofHex) MarshalText() ([]byte, error) {
	return s.proof.MarshalText()
}

func (s *SchnorrProofHex) UnmarshalText(text []byte) error {
	return s.UnmarshalText(text)
}
