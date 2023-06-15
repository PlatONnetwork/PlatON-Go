package bls

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestBls(t *testing.T) {
	sk := &SecretKey{}
	sk.SetByCSPRNG()
	sk = sk.ToSwapKey()
	sign := sk.Sign("hello")

	bytes := sk.Serialize()
	require.True(t, sign.Verify(sk.GetPublicKey(), "hello"))
	var sk2 SecretKey
	require.Nil(t, sk2.Deserialize(bytes))
	sign2 := sk2.Sign("hello")
	require.Equal(t, sign.Serialize(), sign2.Serialize())
	require.Nil(t, sign2.Deserialize(sign.Serialize()))
	require.Equal(t, sk.GetPublicKey().Serialize(), sk2.GetPublicKey().Serialize())
	bytes, err := sk.GetPublicKey().MarshalText()
	require.Nil(t, err)

	var pub PublicKey
	pub.UnmarshalText(bytes)
	require.Equal(t, pub.Serialize(), sk.GetPublicKey().Serialize())

	ethKey := sk.ToEthKey()
	require.Equal(t, ethKey.GetLittleEndian(), sk2.GetLittleEndian())
	require.Equal(t, ethKey.Serialize(), sk2.ToEthKey().Serialize())
	sign = ethKey.Sign("hello")
	require.True(t, sign.Verify(ethKey.GetPublicKey(), "hello"))

	proof, err := sk.MakeSchnorrNIZKP()
	require.Nil(t, err)
	bytes, err = proof.MarshalText()
	require.Nil(t, err)
	require.True(t, len(bytes) > 0)
	require.Nil(t, proof.UnmarshalText(bytes))
	proof.VerifySchnorrNIZK(sk.GetPublicKey())

	proof, err = ethKey.MakeSchnorrNIZKP()
	require.Nil(t, err)
	bytes, err = proof.MarshalText()
	require.Nil(t, err)
	require.True(t, len(bytes) > 0)
	require.Nil(t, proof.UnmarshalText(bytes))
	proof.VerifySchnorrNIZK(ethKey.GetPublicKey())
}

func TestBlsHex(t *testing.T) {
	sk := &SecretKey{}
	sk.SetByCSPRNG()
	sk = sk.ToSwapKey()
	var pubHex PublicKeyHex
	pubText, err := sk.GetPublicKey().MarshalText()
	require.Nil(t, err)
	require.Nil(t, pubHex.UnmarshalText(pubText))
	pubKey, err := pubHex.ParseBlsPubKey()
	require.Nil(t, err)
	require.Equal(t, pubKey.Serialize(), sk.GetPublicKey().Serialize())
	require.Equal(t, pubKey.Bytes(), sk.GetPublicKey().Bytes())

	var proofHex SchnorrProofHex
	proof, err := sk.MakeSchnorrNIZKP()
	require.Nil(t, err)
	proofText, err := proof.MarshalText()
	require.Nil(t, proofHex.UnmarshalText(proofText))
	require.Nil(t, proof.UnmarshalText(proofText))
	require.Nil(t, proof.VerifySchnorrNIZK(sk.GetPublicKey()))
}

func TestBlsAgg(t *testing.T) {
	var sks []*SecretKey
	var pubs []*PublicKey
	for i := 0; i < 5; i++ {
		sk := &SecretKey{}
		sk.SetByCSPRNG()
		sk = sk.ToSwapKey()
		sks = append(sks, sk)
		pubs = append(pubs, sk.GetPublicKey())
	}

	var sign *Sign
	for _, sk := range sks {
		if sign == nil {
			sign = sk.Sign("hello")
		} else {
			sign.Add(sk.Sign("hello"))
		}
	}

	var pub *PublicKey
	for _, sk := range sks {
		if pub == nil {
			pub = sk.GetPublicKey()
		} else {
			pub.Add(sk.GetPublicKey())
		}
	}
	require.True(t, sign.Verify(pub, "hello"))
}
