package bls

import (
	"testing"

	"fmt"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)

func TestSchnorrNIZK(t *testing.T) {
	err := Init(CurveFp254BNb)
	if err != nil {
		t.Fatal(err)
	}
	var sk SecretKey
	sk.SetByCSPRNG()
	fmt.Printf("sk1=%s\n", sk.GetHexString())
	proof, err := sk.MakeSchnorrNIZKP()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("sk2=%s\n", sk.GetHexString())
	pk := sk.GetPublicKey()
	err = proof.VerifySchnorrNIZK(*pk)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSchnorrProofText(t *testing.T) {

	// not support：
	//		CurveFp382_1
	//		CurveFp382_2
	// support:
	// 		CurveFp254BNb
	// 		BLS12_381

	err := Init(CurveFp254BNb)
	//err := Init(BLS12_381)
	if err != nil {
		t.Fatal(err)
	}
	var a SecretKey
	a.SetByCSPRNG()
	schnorrProof, err := a.MakeSchnorrNIZKP()
	if err != nil {
		t.Fatal(err)
	}
	// to text
	textByte, err := schnorrProof.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("schnorrProof: ", string(textByte))
	t.Log("\n")

	// proof to SchnorrProofHex
	var schnorrProofhex SchnorrProofHex
	schnorrProofhex.UnmarshalText(textByte)

	t.Log("schnorrProofhex: ", string(schnorrProofhex[:]))
	t.Log("\n")
	// schnorrProofhex rlp encode
	pentriesRlp, err := rlp.EncodeToBytes(schnorrProofhex)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("schnorrProofhex rlp:", string(pentriesRlp))
	t.Log("\n")

	// proof to SchnorrProofHex2
	var schnorrProofhex2 SchnorrProofHex
	err = rlp.DecodeBytes(pentriesRlp, &schnorrProofhex2)
	if err != nil {
		t.Fatal(err)
	}

	textByte2, err := schnorrProofhex2.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	t.Log("schnorrProofhex2: ", string(textByte2))

	proof2 := new(SchnorrProof)
	err = proof2.UnmarshalText(textByte2)
	if err != nil {
		t.Fatal(err)
	}

	P := a.GetPublicKey()

	err = proof2.VerifySchnorrNIZK(*P)
	if err != nil {
		t.Fatal(err)
	}
}
