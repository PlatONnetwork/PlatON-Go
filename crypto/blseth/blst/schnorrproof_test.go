package blst

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestProof(t *testing.T) {
	key, _ := RandKey()
	proof := key.(*bls12SecretKey).MakeSchnorrNIZKP()
	assert.Nil(t, proof.Verify(key.PublicKey()))
}
