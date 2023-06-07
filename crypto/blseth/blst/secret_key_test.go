//go:build ((linux && amd64) || (linux && arm64) || (darwin && amd64) || (darwin && arm64) || (windows && amd64)) && !blst_disabled
// +build linux,amd64 linux,arm64 darwin,amd64 darwin,arm64 windows,amd64
// +build !blst_disabled

package blst_test

import (
	"bytes"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/crypto/blseth/blst"
	"github.com/PlatONnetwork/PlatON-Go/crypto/blseth/common"
)

func ToBytes32(x []byte) [32]byte {
	var y [32]byte
	copy(y[:], x)
	return y
}
func TestMarshalUnmarshal(t *testing.T) {
	priv, err := blst.RandKey()
	require.NoError(t, err)
	b := priv.Marshal()
	b32 := ToBytes32(b)
	pk, err := blst.SecretKeyFromBytes(b32[:])
	require.NoError(t, err)
	pk2, err := blst.SecretKeyFromBytes(b32[:])
	require.NoError(t, err)
	assert.Equal(t, pk.Marshal(), pk2.Marshal(), "Keys not equal")
}

func TestSecretKeyFromBytes(t *testing.T) {
	tests := []struct {
		name  string
		input []byte
		err   error
	}{
		{
			name: "Nil",
			err:  errors.New("secret key must be 32 bytes"),
		},
		{
			name:  "Empty",
			input: []byte{},
			err:   errors.New("secret key must be 32 bytes"),
		},
		{
			name:  "Short",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("secret key must be 32 bytes"),
		},
		{
			name:  "Long",
			input: []byte{0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00},
			err:   errors.New("secret key must be 32 bytes"),
		},
		{
			name:  "Bad",
			input: []byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
			err:   common.ErrSecretUnmarshal,
		},
		{
			name:  "Good",
			input: []byte{0x25, 0x29, 0x5f, 0x0d, 0x1d, 0x59, 0x2a, 0x90, 0xb3, 0x33, 0xe2, 0x6e, 0x85, 0x14, 0x97, 0x08, 0x20, 0x8e, 0x9f, 0x8e, 0x8b, 0xc1, 0x8f, 0x6c, 0x77, 0xbd, 0x62, 0xf8, 0xad, 0x7a, 0x68, 0x66},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			res, err := blst.SecretKeyFromBytes(test.input)
			if test.err != nil {
				assert.NotEqual(t, nil, err, "No error returned")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, 0, bytes.Compare(res.Marshal(), test.input))
			}
		})
	}
}

//sk:1ba530af467ed12ca19a85e237f3d87e8c2a2082ce27b0d3151a163007821991
//vk:8065cbede15639eb0bf6449d37d1b2a2532f5beea3b4bef882f3798af39205229a44a574500382e07e29b630f23a1085
//hash:8e92e763fe002fc6ce0c4733bdaf4c7ec1b31971117023606fc521defbf0234783924914b338399e9b177fb272d31cd715806c34539b6746f097c92b68083696eae15f765b93b3b778da0774ba758451bd73f4e2ea4afe0792e80d13361f9a5c
//sig:92d9ac6f8de08004776e6d99f79b161f054ecaf30fa8f160e2028a993b23829a11dd733659cea875ad3d907a03a42ea7097d617d68755ba3a0451b524820f9056618273bb709951ecc53c4fd979db9dcc41c77c76f15a97a64bd1726bc9f3c19

func TestSign(t *testing.T) {
	rk, err := blst.RandKey()
	require.NoError(t, err)
	msg := []byte("hello")
	sig := rk.Sign(msg)
	fmt.Println("sk:", hex.EncodeToString(rk.Marshal()))
	fmt.Println("pk:", hex.EncodeToString(rk.PublicKey().Marshal()))
	fmt.Println("sig:", hex.EncodeToString(sig.Marshal()))
	fmt.Println(sig.Verify(rk.PublicKey(), msg))
	bytes, _ := hex.DecodeString("1ba530af467ed12ca19a85e237f3d87e8c2a2082ce27b0d3151a163007821991")
	sk, err := blst.SecretKeyFromBytes(bytes)
	if err != nil {
		t.Fatal(err)
	}
	sig = sk.Sign(msg)
	fmt.Println("sig:", hex.EncodeToString(sig.Marshal()))

}
func TestSerialize(t *testing.T) {
	rk, err := blst.RandKey()
	require.NoError(t, err)
	b := rk.Marshal()

	_, err = blst.SecretKeyFromBytes(b)
	assert.NoError(t, err)
}

func TestZeroKey(t *testing.T) {
	// Is Zero
	var zKey [32]byte
	assert.Equal(t, true, blst.IsZero(zKey[:]))

	// Is Not Zero
	_, err := rand.Read(zKey[:])
	assert.NoError(t, err)
	assert.Equal(t, false, blst.IsZero(zKey[:]))
}
