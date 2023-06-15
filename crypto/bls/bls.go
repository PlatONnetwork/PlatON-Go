package bls

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls/eth"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls/swap"
	"github.com/PlatONnetwork/PlatON-Go/crypto/bls/types"
	blsswap "github.com/PlatONnetwork/PlatON-Go/crypto/blsswap"
	"io"
	"os"
	"strings"
)

const (
	Bls12381     = 0
	Bls12381Swap = 1

	EthPubKeyLength     = 48
	SwapPubKeyLength    = 96
	EthPubKeyHexLength  = EthPubKeyLength * 2
	SwapPubKeyHexLength = SwapPubKeyLength * 2

	EthSignLength         = 96
	SwapSignLength        = 48
	EthSignHexLength      = EthSignLength * 2
	SwapSignHExLength     = SwapSignLength * 2
	SchnorrProofHexLength = 128
)

var (
	BlsVersion = Bls12381Swap
)

func init() {
	Init(Bls12381Swap)
}
func Init(version int) {
	if version == Bls12381Swap {
		blsswap.Init(blsswap.BLS12_381)
	}
	BlsVersion = version
}

type SecretKey struct {
	key types.ISecretKey
}

func (s SecretKey) Sign(m string) (sign *Sign) {
	return &Sign{s.key.Sign(m)}
}

func (s SecretKey) Serialize() []byte {
	return s.key.Serialize()
}

func (s *SecretKey) createKey() {
	if s.key == nil {
		switch BlsVersion {
		case Bls12381:
			s.key = &eth.SecretKey{}
		case Bls12381Swap:
			s.key = &swap.SecretKey{}
		}
	}
}

func (s *SecretKey) ToSwapKey() *SecretKey {
	key := &swap.SecretKey{}
	key.SetLittleEndian(s.GetLittleEndian())
	return &SecretKey{
		key: key,
	}
}

func (s *SecretKey) ToEthKey() *SecretKey {
	key := &eth.SecretKey{}
	key.SetLittleEndian(s.GetLittleEndian())
	return &SecretKey{
		key: key,
	}
}

func (s *SecretKey) Deserialize(buf []byte) error {
	s.createKey()
	return s.key.Deserialize(buf)
}

func (s *SecretKey) SetLittleEndian(buf []byte) error {
	s.createKey()
	return s.key.SetLittleEndian(buf)
}

func (s *SecretKey) GetLittleEndian() []byte {
	return s.key.GetLittleEndian()
}

func (s *SecretKey) SetByCSPRNG() {
	s.createKey()
	s.key.SetByCSPRNG()
}

func (s SecretKey) GetPublicKey() *PublicKey {
	return &PublicKey{
		s.key.GetPublicKey(),
	}
}

func (s SecretKey) MakeSchnorrNIZKP() (*SchnorrProof, error) {
	proof, err := s.key.MakeSchnorrNIZKP()
	if err != nil {
		return nil, err
	}
	return &SchnorrProof{
		proof: proof,
	}, nil
}

type PublicKey struct {
	pubKey types.IPublicKey
}

func (p *PublicKey) createKeyByText(length int) error {
	if p.pubKey == nil {
		switch length {
		case EthPubKeyHexLength:
			p.pubKey = &eth.PublicKey{}
		case SwapPubKeyHexLength:
			p.pubKey = swap.NewPublicKey()
		default:
			return errors.New("illegal length")
		}
	}
	return nil
}
func (p *PublicKey) Add(rhs *PublicKey) {
	p.pubKey.Add(rhs.pubKey)
}

func (p *PublicKey) UnmarshalText(text []byte) error {
	if err := p.createKeyByText(len(text)); err != nil {
		return err
	}
	return p.pubKey.UnmarshalText(text)
}

func (p PublicKey) Serialize() []byte {
	return p.pubKey.Serialize()
}

func (p PublicKey) Bytes() []byte {
	return p.pubKey.Bytes()
}

func (p PublicKey) MarshalText() ([]byte, error) {
	return p.pubKey.MarshalText()
}

type PublicKeyHex struct {
	pubKey types.IPublicKeyHex
}

func (p *PublicKeyHex) createKeyByText(length int) error {
	if p.pubKey == nil {
		switch length {
		case EthPubKeyHexLength:
			p.pubKey = &eth.PublicKeyHex{}
		case SwapPubKeyHexLength:
			p.pubKey = &swap.PublicKeyHex{}
		default:
			return errors.New("illegal length")
		}
	}
	return nil
}
func (p PublicKeyHex) ParseBlsPubKey() (*PublicKey, error) {
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
	if err := p.createKeyByText(len(text)); err != nil {
		return err
	}
	return p.pubKey.UnmarshalText(text)
}

type Sign struct {
	sign types.ISign
}

func (s *Sign) createKey(length int) error {
	if s.sign == nil {
		switch length {
		case EthSignLength:
			s.sign = &eth.Sign{}
		case SwapSignLength:
			s.sign = &swap.Sign{}
		default:
			return errors.New("illegal length")
		}
	}
	return nil
}
func (s *Sign) Add(rhs *Sign) {
	s.sign.Add(rhs.sign)
}
func (s Sign) Verify(pub *PublicKey, m string) bool {
	return s.sign.Verify(pub.pubKey, m)
}

func (s Sign) Serialize() []byte {
	return s.sign.Serialize()
}

func (s *Sign) Deserialize(buf []byte) error {
	if err := s.createKey(len(buf)); err != nil {
		return err
	}
	return s.sign.Deserialize(buf)
}

type SchnorrProof struct {
	proof types.ISchnorrProof
	text  []byte
}

func (s *SchnorrProof) VerifySchnorrNIZK(pk *PublicKey) error {
	if s.proof == nil && s.text == nil {
		return errors.New("proof is nil")
	}
	if s.proof == nil {
		switch len(pk.Serialize()) {
		case EthPubKeyLength:
			s.proof = &eth.SchnorrProof{}
			s.proof.UnmarshalText(s.text)
		case SwapPubKeyLength:
			s.proof = &swap.SchnorrProof{}
			s.proof.UnmarshalText(s.text)
		default:
			return errors.New("illegal public key")
		}
	}
	return s.proof.VerifySchnorrNIZK(pk.pubKey)
}

func (s *SchnorrProof) UnmarshalText(text []byte) error {
	if len(text) != SchnorrProofHexLength {
		return errors.New("illegal length")
	}
	s.text = text
	return nil
}

func (s SchnorrProof) MarshalText() ([]byte, error) {
	return s.proof.MarshalText()
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

func LoadBLS(file string) (*SecretKey, error) {
	buf := make([]byte, 64)
	fd, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer fd.Close()
	if _, err := io.ReadFull(fd, buf); err != nil {
		return nil, err
	}
	var sec SecretKey
	key, err := hex.DecodeString(string(buf))
	if err != nil {
		return nil, err
	}
	err = sec.SetLittleEndian(key)
	return &sec, err
}

func SaveBLS(file string, key *SecretKey) error {
	k := hex.EncodeToString(key.GetLittleEndian())
	return os.WriteFile(file, []byte(k), 0600)
}

func GenerateKey() *SecretKey {
	var privateKey SecretKey
	privateKey.SetByCSPRNG()
	return &privateKey
}
