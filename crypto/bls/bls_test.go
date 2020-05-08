package bls

import (
	"testing"

	"github.com/PlatONnetwork/PlatON-Go/rlp"
)
import (
	"bytes"
	"fmt"
	"strconv"

	"github.com/PlatONnetwork/PlatON-Go/crypto"
)

var unitN = 0

// Tests (for Benchmarks see below)

func testPre(t *testing.T) {
	t.Log("init")
	{
		var id ID
		err := id.SetLittleEndian([]byte{6, 5, 4, 3, 2, 1})
		if err != nil {
			t.Error(err)
		}
		t.Log("id :", id.GetHexString())
		var id2 ID
		err = id2.SetHexString(id.GetHexString())
		if err != nil {
			t.Fatal(err)
		}
		if !id.IsEqual(&id2) {
			t.Errorf("not same id\n%s\n%s", id.GetHexString(), id2.GetHexString())
		}
		err = id2.SetDecString(id.GetDecString())
		if err != nil {
			t.Fatal(err)
		}
		if !id.IsEqual(&id2) {
			t.Errorf("not same id\n%s\n%s", id.GetDecString(), id2.GetDecString())
		}
	}
	{
		var sec SecretKey
		err := sec.SetLittleEndian([]byte{1, 2, 3, 4, 5, 6})
		if err != nil {
			t.Error(err)
		}
		t.Log("sec=", sec.GetHexString())
	}

	t.Log("create secret key")
	m := "this is a bls sample for go"
	var sec SecretKey
	sec.SetByCSPRNG()
	t.Log("sec:", sec.GetHexString())
	t.Log("create public key")
	pub := sec.GetPublicKey()
	t.Log("pub:", pub.GetHexString())
	sign := sec.Sign(m)
	t.Log("sign:", sign.GetHexString())
	if !sign.Verify(pub, m) {
		t.Error("Signature does not verify")
	}

	// How to make array of SecretKey
	{
		sec := make([]SecretKey, 3)
		for i := 0; i < len(sec); i++ {
			sec[i].SetByCSPRNG()
			t.Log("sec=", sec[i].GetHexString())
		}
	}
}

func testByCs(t *testing.T) {
	/*err := Init(CurveFp254BNb)
	if err != nil {
		t.Fatal(err)
	}*/
	sec := make([]SecretKey, 3)
	for i := 0; i < len(sec); i++ {
		sec[i].SetByCSPRNG()
		t.Log("sec=", sec[i].GetHexString())
	}
}

func testStringConversion(t *testing.T) {
	t.Log("testRecoverSecretKey")
	var sec SecretKey
	var s string
	if unitN == 6 {
		s = "16798108731015832284940804142231733909759579603404752749028378864165570215949"
	} else {
		s = "40804142231733909759579603404752749028378864165570215949"
	}
	err := sec.SetDecString(s)
	if err != nil {
		t.Fatal(err)
	}
	if s != sec.GetDecString() {
		t.Error("not equal")
	}
	s = sec.GetHexString()
	var sec2 SecretKey
	err = sec2.SetHexString(s)
	if err != nil {
		t.Fatal(err)
	}
	if !sec.IsEqual(&sec2) {
		t.Error("not equal")
	}
}

func testRecoverSecretKey(t *testing.T) {
	t.Log("testRecoverSecretKey")
	k := 3000
	var sec SecretKey
	sec.SetByCSPRNG()
	t.Logf("sec=%s\n", sec.GetHexString())

	// make master secret key
	msk := sec.GetMasterSecretKey(k)

	n := k
	secVec := make([]SecretKey, n)
	idVec := make([]ID, n)
	for i := 0; i < n; i++ {
		err := idVec[i].SetLittleEndian([]byte{byte(i & 255), byte(i >> 8), 2, 3, 4, 5})
		if err != nil {
			t.Error(err)
		}
		err = secVec[i].Set(msk, &idVec[i])
		if err != nil {
			t.Error(err)
		}
		// t.Logf("idVec[%d]=%s\n", i, idVec[i].GetHexString())
	}
	// recover sec2 from secVec and idVec
	var sec2 SecretKey
	err := sec2.Recover(secVec, idVec)
	if err != nil {
		t.Error(err)
	}
	if !sec.IsEqual(&sec2) {
		t.Errorf("Mismatch in recovered secret key:\n  %s\n  %s.", sec.GetHexString(), sec2.GetHexString())
	}
}

func testEachSign(t *testing.T, m string, msk []SecretKey, mpk []PublicKey) ([]ID, []SecretKey, []PublicKey, []Sign) {
	idTbl := []byte{3, 5, 193, 22, 15}
	n := len(idTbl)

	secVec := make([]SecretKey, n)
	pubVec := make([]PublicKey, n)
	signVec := make([]Sign, n)
	idVec := make([]ID, n)

	for i := 0; i < n; i++ {
		err := idVec[i].SetLittleEndian([]byte{idTbl[i], 0, 0, 0, 0, 0})
		if err != nil {
			t.Error(err)
		}
		t.Logf("idVec[%d]=%s\n", i, idVec[i].GetHexString())

		err = secVec[i].Set(msk, &idVec[i])
		if err != nil {
			t.Error(err)
		}

		err = pubVec[i].Set(mpk, &idVec[i])
		if err != nil {
			t.Error(err)
		}
		t.Logf("pubVec[%d]=%s\n", i, pubVec[i].GetHexString())

		if !pubVec[i].IsEqual(secVec[i].GetPublicKey()) {
			t.Errorf("Pubkey derivation does not match\n%s\n%s", pubVec[i].GetHexString(), secVec[i].GetPublicKey().GetHexString())
		}

		signVec[i] = *secVec[i].Sign(m)
		if !signVec[i].Verify(&pubVec[i], m) {
			t.Error("Pubkey derivation does not match")
		}
	}
	return idVec, secVec, pubVec, signVec
}
func testSign(t *testing.T) {
	m := "testSign"
	t.Log(m)

	var sec0 SecretKey
	sec0.SetByCSPRNG()
	pub0 := sec0.GetPublicKey()
	s0 := sec0.Sign(m)
	if !s0.Verify(pub0, m) {
		t.Error("Signature does not verify")
	}

	k := 3
	msk := sec0.GetMasterSecretKey(k)
	mpk := GetMasterPublicKey(msk)
	idVec, secVec, pubVec, signVec := testEachSign(t, m, msk, mpk)

	var sec1 SecretKey
	err := sec1.Recover(secVec, idVec)
	if err != nil {
		t.Error(err)
	}
	if !sec0.IsEqual(&sec1) {
		t.Error("Mismatch in recovered seckey.")
	}
	var pub1 PublicKey
	err = pub1.Recover(pubVec, idVec)
	if err != nil {
		t.Error(err)
	}
	if !pub0.IsEqual(&pub1) {
		t.Error("Mismatch in recovered pubkey.")
	}
	var s1 Sign
	err = s1.Recover(signVec, idVec)
	if err != nil {
		t.Error(err)
	}
	if !s0.IsEqual(&s1) {
		t.Error("Mismatch in recovered signature.")
	}
}

func testAdd(t *testing.T) {
	t.Log("testAdd")
	var sec1 SecretKey
	var sec2 SecretKey
	sec1.SetByCSPRNG()
	sec2.SetByCSPRNG()

	pub1 := sec1.GetPublicKey()
	pub2 := sec2.GetPublicKey()

	m := "test test"
	sign1 := sec1.Sign(m)
	sign2 := sec2.Sign(m)

	t.Log("sign1    :", sign1.GetHexString())
	sign1.Add(sign2)
	t.Log("sign1 add:", sign1.GetHexString())
	pub1.Add(pub2)
	if !sign1.Verify(pub1, m) {
		t.Fail()
	}
}

func testPop(t *testing.T) {
	t.Log("testPop")
	var sec SecretKey
	sec.SetByCSPRNG()
	pop := sec.GetPop()
	if !pop.VerifyPop(sec.GetPublicKey()) {
		t.Errorf("Valid Pop does not verify")
	}
	sec.SetByCSPRNG()
	if pop.VerifyPop(sec.GetPublicKey()) {
		t.Errorf("Invalid Pop verifies")
	}
}

func testData(t *testing.T) {
	t.Log("testData")
	var sec1, sec2 SecretKey
	sec1.SetByCSPRNG()
	b := sec1.GetLittleEndian()
	err := sec2.SetLittleEndian(b)
	if err != nil {
		t.Fatal(err)
	}
	if !sec1.IsEqual(&sec2) {
		t.Error("SecretKey not same")
	}
	pub1 := sec1.GetPublicKey()
	b = pub1.Serialize()
	var pub2 PublicKey
	err = pub2.Deserialize(b)
	if err != nil {
		t.Fatal(err)
	}
	if !pub1.IsEqual(&pub2) {
		t.Error("PublicKey not same")
	}
	m := "doremi"
	sign1 := sec1.Sign(m)
	b = sign1.Serialize()
	var sign2 Sign
	err = sign2.Deserialize(b)
	if err != nil {
		t.Fatal(err)
	}
	if !sign1.IsEqual(&sign2) {
		t.Error("Sign not same")
	}
}

func testOrder(t *testing.T, c int) {
	var curve string
	var field string
	if c == CurveFp254BNb {
		curve = "16798108731015832284940804142231733909759579603404752749028378864165570215949"
		field = "16798108731015832284940804142231733909889187121439069848933715426072753864723"
	} else if c == CurveFp382_1 {
		curve = "5540996953667913971058039301942914304734176495422447785042938606876043190415948413757785063597439175372845535461389"
		field = "5540996953667913971058039301942914304734176495422447785045292539108217242186829586959562222833658991069414454984723"
	} else if c == CurveFp382_2 {
		curve = "5541245505022739011583672869577435255026888277144126952448297309161979278754528049907713682488818304329661351460877"
		field = "5541245505022739011583672869577435255026888277144126952450651294188487038640194767986566260919128250811286032482323"
	} else if c == BLS12_381 {
		curve = "52435875175126190479447740508185965837690552500527637822603658699938581184513"
		field = "4002409555221667393417789825735904156556882819939007885332058136124031650490837864442687629129015664037894272559787"
	} else {
		t.Fatal("bad c", c)
	}
	s := GetCurveOrder()
	if s != curve {
		t.Errorf("bad curve order\n%s\n%s\n", s, curve)
	}
	s = GetFieldOrder()
	if s != field {
		t.Errorf("bad field order\n%s\n%s\n", s, field)
	}
}

func testDHKeyExchange(t *testing.T) {
	var sec1, sec2 SecretKey
	sec1.SetByCSPRNG()
	sec2.SetByCSPRNG()
	pub1 := sec1.GetPublicKey()
	pub2 := sec2.GetPublicKey()
	out1 := DHKeyExchange(&sec1, pub2)
	out2 := DHKeyExchange(&sec2, pub1)
	if !out1.IsEqual(&out2) {
		t.Errorf("DH key is not equal")
	}
}

func test(t *testing.T, c int) {
	err := Init(c)
	if err != nil {
		t.Fatal(err)
	}
	unitN = GetOpUnitSize()
	t.Logf("unitN=%d\n", unitN)
	testPairing(t)
	testPre(t)
	testRecoverSecretKey(t)
	testAdd(t)
	testSign(t)
	testPop(t)
	testData(t)
	testStringConversion(t)
	testOrder(t, c)
	testDHKeyExchange(t)
	//add
	testByCs(t)
	testAggregateSign(t, c)
	testSchnorr_test(t, c)
	testSchnorrNIZk(t, c)
	//testProofText(t, c)

}

func TestNmain(t *testing.T) {
	t.Logf("GetMaxOpUnitSize() = %d\n", GetMaxOpUnitSize())
	t.Logf("GetFrUnitSize() = %d\n", GetFrUnitSize())
	t.Log("CurveFp254BNb")
	test(t, CurveFp254BNb)
	if GetMaxOpUnitSize() == 6 {
		if GetFrUnitSize() == 6 {
			t.Log("CurveFp382_1")
			test(t, CurveFp382_1)
			t.Log("CurveFp382_2")
			test(t, CurveFp382_2)
		} else {
			t.Log("BLS12_381")
			test(t, BLS12_381)
		}
	}
}

//add @20190716
func testGetGOfG2(t *testing.T, c int) {
	err := Init(c)
	if err != nil {
		t.Fatal(err)
	}
	P := GetGeneratorOfG2()
	fmt.Printf("P=%s\n", P.GetHexString())
}

func TestForGetG(t *testing.T) {
	testGetGOfG2(t, CurveFp254BNb)
	if GetMaxOpUnitSize() == 6 {
		if GetFrUnitSize() == 6 {
			testGetGOfG2(t, CurveFp382_1)
			testGetGOfG2(t, CurveFp382_2)
		}
		testGetGOfG2(t, BLS12_381)
	}
}

func testAggregateSign(t *testing.T, c int) {
	err := Init(c)
	if err != nil {
		t.Fatal(err)
	}
	m := "test test"
	var k int = 3
	msk := make([]SecretKey, k)
	mpk := make([]PublicKey, k)
	msig := make([]Sign, k)
	for i := 0; i < k; i++ {
		msk[i].SetByCSPRNG()
		mpk[i] = *msk[i].GetPublicKey()
		msig[i] = *msk[i].Sign(m)
		fmt.Printf("msk[%d]=%s\n", i, msk[i].GetHexString())
		fmt.Printf("mpk[%d]=%s\n", i, mpk[i].GetHexString())
		fmt.Printf("msig[%d]=%s\n", i, msig[i].GetHexString())
		if !msig[i].Verify(&mpk[i], m) {
			fmt.Println("verify fail")
		}
	}
	var sig Sign
	for i := 0; i < len(msig); i++ {
		sig.Add(&msig[i])
	}
	fmt.Printf("sig=%s\n", sig.GetHexString())

	err = BatchVerifySameMsg(c, m, mpk, sig)
	if err != nil {
		t.Fatal(err)
	}
	input := []string{"test test", "test test add one", "test test add two"}
	p_hm, err1 := MsgsToHashToG1(input)
	if err1 != nil {
		t.Fatal(err1)
	}
	for j := 0; j < k; j++ {
		msig[j] = *msk[j].Sign(input[j])
		fmt.Printf("msig[%d]=%s\n", j, msig[j].GetHexString())
		if !msig[j].Verify(&mpk[j], input[j]) {
			fmt.Println("verify fail")
		}
	}

	var sig1 Sign
	for i := 0; i < len(msig); i++ {
		sig1.Add(&msig[i])
	}
	fmt.Printf("sig1=%s\n", sig1.GetHexString())

	err = BatchVerifyDistinctMsg(c, mpk, p_hm, sig1)
	if err != nil {
		t.Fatal(err)
	}
}

func testSchnorr_test(t *testing.T, cu int) {
	err := Init(cu)
	if err != nil {
		t.Fatal(err)
	}
	var a SecretKey
	var v SecretKey
	a.SetByCSPRNG()
	v.SetByCSPRNG()
	P := a.GetPublicKey()
	fmt.Printf("P=%s\n", P.GetHexString())
	V := v.GetPublicKey()
	fmt.Printf("V=%s\n", V.GetHexString())
	G := GetGeneratorOfG2()
	fmt.Printf("G=%s\n", G.GetHexString())
	input1 := G.Serialize()
	input2 := V.Serialize()
	input3 := P.Serialize()
	var buffer bytes.Buffer
	buffer.Write(input1)
	buffer.Write(input2)
	buffer.Write(input3)
	output := buffer.Bytes()
	h := crypto.Keccak256(output)
	var c SecretKey
	err = c.SetLittleEndian(h)
	if err != nil {
		fmt.Println("Deserialize fail")
		t.Error(err)
	}
	temp := a
	temp.Mul(&c)
	r := v
	r.Sub(&temp)
	Pr := r.GetPublicKey()
	var Pc PublicKey
	Pc = *P
	Pc.Mul(&c)
	Psum := Pr
	Psum.Add(&Pc)
	fmt.Printf("r=%s\n", r.GetHexString())
	fmt.Printf("c=%s\n", c.GetHexString())
	fmt.Printf("Psum=%s\n", Psum.GetHexString())
	fmt.Printf("Pc=%s\n", Pc.GetHexString())
	if !V.IsEqual(Psum) {
		t.Errorf("V = G*[r] + A*[c] not equal")
	}

	fmt.Printf("G is valid=%v\n", G2IsValid(G))
	fmt.Printf("P is valid=%v\n", G2IsValid(P))
	fmt.Printf("V is valid=%v\n", G2IsValid(V))

	err = Schnorr_test(cu, r, c, *G, *V, *P)
	if err != nil {
		fmt.Println("test fail")
		t.Error(err)
	}
}

func testSchnorrNIZk(t *testing.T, c int) {
	err := Init(c)
	if err != nil {
		t.Fatal(err)
	}
	var a SecretKey
	a.SetByCSPRNG()
	proof, err := SchnorrNIZKProve(c, a)
	if err != nil {
		t.Fatal(err)
	}

	P := a.GetPublicKey()
	err = SchnorrNIZKVerify(c, *proof, *P)
	if err != nil {
		t.Fatal(err)
	}
}

func testSynthSameMsg(t *testing.T, c int) {
	err := Init(c)
	if err != nil {
		t.Fatal(err)
	}
	var k int = 10
	var sk SecretKey
	mpk := make([]PublicKey, k)
	input := make([]string, k)
	for i := 0; i < k; i++ {
		input[i] = fmt.Sprintf("test%d", i)
		sk.SetByCSPRNG()
		mpk[i] = *sk.GetPublicKey()
		fmt.Printf("input[%d]=%s\n", i, input[i])
		fmt.Printf("mpk[%d]=%s\n", i, mpk[i].GetHexString())
	}

	input[3] = "test1"
	input[6] = "test8"
	input[4] = "test5"
	input[7] = "test5"

	pk, msg, index, err := SynthSameMsg(c, mpk, input)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("pklen:%d\n", len(pk))
	fmt.Printf("msglen:%d\n", len(msg))
	for j := 0; j < len(pk); j++ {
		fmt.Printf("msg[%d]=%s\n", j, msg[j])
		fmt.Printf("pk[%d]=%s\n", j, pk[j].GetHexString())
	}
	for m := 0; m < len(index); m++ {
		fmt.Printf("index[%d]=%s\n", m, index[m])
	}

}

func TestBlsInit(t *testing.T) {
	fmt.Println("bls MCLBN_FP_UNIT_SIZE:", GetMaxOpUnitSize())
	fmt.Println("bls MCLBN_FR_UNIT_SIZE:", GetFrUnitSize())
	err := Init(CurveFp254BNb)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("exc bls CurveFp254BNb")
	if GetMaxOpUnitSize() == 6 {
		if GetFrUnitSize() == 6 {
			err := Init(CurveFp382_1)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("exc bls CurveFp382_1")
		}
		err := Init(BLS12_381)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("exc bls BLS12_381")
	}
}

func TestProof_RLP(t *testing.T) {
	err := Init(CurveFp254BNb)
	//err := Init(BLS12_381)
	if err != nil {
		t.Fatal(err)
	}
	var a SecretKey
	a.SetByCSPRNG()
	proof, err := SchnorrNIZKProve(CurveFp254BNb, a)
	//proof, err := SchnorrNIZKProve(BLS12_381, a)
	if err != nil {
		t.Fatal(err)
	}

	// proof rlp encode
	proofRlp, err := rlp.EncodeToBytes(proof)
	if err != nil {
		t.Fatal(err)
	}

	t.Log("proof rlp:", string(proofRlp))
	t.Log("\n")

	var proof2 Proof
	err = rlp.DecodeBytes(proofRlp, &proof2)
	if err != nil {
		t.Fatal(err)
	}

	P := a.GetPublicKey()
	err = SchnorrNIZKVerify(CurveFp254BNb, proof2, *P)
	//err = SchnorrNIZKVerify(BLS12_381, proof2, *P)
	if err != nil {
		t.Fatal(err)
	}
}

// Benchmarks

var curve = CurveFp382_1

//var curve = CurveFp254BNb

func BenchmarkBatchVerifyDistinctMsg(b *testing.B) {
	b.StopTimer()
	err := Init(CurveFp254BNb)
	if err != nil {
		b.Fatal(err)
	}
	var k int = 2
	msk := make([]SecretKey, k)
	mpk := make([]PublicKey, k)
	msig := make([]Sign, k)
	input := make([]string, k)
	for i := 0; i < k; i++ {
		input[i] = fmt.Sprintf("test%d", i)
		msk[i].SetByCSPRNG()
		mpk[i] = *msk[i].GetPublicKey()
		msig[i] = *msk[i].Sign(input[i])
		//fmt.Printf("msk[%d]=%s\n", i, msk[i].GetHexString())
		//fmt.Printf("mpk[%d]=%s\n", i, mpk[i].GetHexString())
		//fmt.Printf("msig[%d]=%s\n", i, msig[i].GetHexString())
		if !msig[i].Verify(&mpk[i], input[i]) {
			fmt.Println("verify fail")
		}
	}
	p_hm, err1 := MsgsToHashToG1(input)
	if err1 != nil {
		b.Fatal(err1)
	}
	var sig Sign
	for i := 0; i < len(msig); i++ {
		sig.Add(&msig[i])
	}
	//fmt.Printf("sig=%s\n", sig.GetHexString())

	for n := 0; n < b.N; n++ {
		b.StartTimer()
		err = BatchVerifyDistinctMsg(CurveFp254BNb, mpk, p_hm, sig)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
	}

}

func BenchmarkBatchVerifySameMsg(b *testing.B) {
	b.StopTimer()
	err := Init(CurveFp254BNb)
	if err != nil {
		b.Fatal(err)
	}
	m := "test test"
	var k int = 2
	msk := make([]SecretKey, k)
	mpk := make([]PublicKey, k)
	msig := make([]Sign, k)
	for i := 0; i < k; i++ {
		msk[i].SetByCSPRNG()
		mpk[i] = *msk[i].GetPublicKey()
		msig[i] = *msk[i].Sign(m)
		fmt.Printf("msk[%d]=%s\n", i, msk[i].GetHexString())
		fmt.Printf("mpk[%d]=%s\n", i, mpk[i].GetHexString())
		fmt.Printf("msig[%d]=%s\n", i, msig[i].GetHexString())
		if !msig[i].Verify(&mpk[i], m) {
			fmt.Println("verify fail")
		}
	}
	var sig Sign
	for i := 0; i < len(msig); i++ {
		sig.Add(&msig[i])
	}
	fmt.Printf("sig=%s\n", sig.GetHexString())

	for n := 0; n < b.N; n++ {
		b.StartTimer()
		err = BatchVerifySameMsg(CurveFp254BNb, m, mpk, sig)
		b.StopTimer()
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPubkeyFromSeckey(b *testing.B) {
	b.StopTimer()
	err := Init(curve)
	if err != nil {
		b.Fatal(err)
	}
	var sec SecretKey
	for n := 0; n < b.N; n++ {
		sec.SetByCSPRNG()
		b.StartTimer()
		sec.GetPublicKey()
		b.StopTimer()
	}
}

func BenchmarkSigning(b *testing.B) {
	b.StopTimer()
	err := Init(curve)
	if err != nil {
		b.Fatal(err)
	}
	var sec SecretKey
	for n := 0; n < b.N; n++ {
		sec.SetByCSPRNG()
		b.StartTimer()
		sec.Sign(strconv.Itoa(n))
		b.StopTimer()
	}
}

func BenchmarkValidation(b *testing.B) {
	b.StopTimer()
	err := Init(curve)
	if err != nil {
		b.Fatal(err)
	}
	var sec SecretKey
	for n := 0; n < b.N; n++ {
		sec.SetByCSPRNG()
		pub := sec.GetPublicKey()
		m := strconv.Itoa(n)
		sig := sec.Sign(m)
		b.StartTimer()
		sig.Verify(pub, m)
		b.StopTimer()
	}
}

func benchmarkDeriveSeckeyShare(k int, b *testing.B) {
	b.StopTimer()
	err := Init(curve)
	if err != nil {
		b.Fatal(err)
	}
	var sec SecretKey
	sec.SetByCSPRNG()
	msk := sec.GetMasterSecretKey(k)
	var id ID
	for n := 0; n < b.N; n++ {
		err = id.SetLittleEndian([]byte{1, 2, 3, 4, 5, byte(n)})
		if err != nil {
			b.Error(err)
		}
		b.StartTimer()
		err := sec.Set(msk, &id)
		b.StopTimer()
		if err != nil {
			b.Error(err)
		}
	}
}

//func BenchmarkDeriveSeckeyShare100(b *testing.B)  { benchmarkDeriveSeckeyShare(100, b) }
//func BenchmarkDeriveSeckeyShare200(b *testing.B)  { benchmarkDeriveSeckeyShare(200, b) }
func BenchmarkDeriveSeckeyShare500(b *testing.B) { benchmarkDeriveSeckeyShare(500, b) }

//func BenchmarkDeriveSeckeyShare1000(b *testing.B) { benchmarkDeriveSeckeyShare(1000, b) }

func benchmarkRecoverSeckey(k int, b *testing.B) {
	b.StopTimer()
	err := Init(curve)
	if err != nil {
		b.Fatal(err)
	}
	var sec SecretKey
	sec.SetByCSPRNG()
	msk := sec.GetMasterSecretKey(k)

	// derive n shares
	n := k
	secVec := make([]SecretKey, n)
	idVec := make([]ID, n)
	for i := 0; i < n; i++ {
		err := idVec[i].SetLittleEndian([]byte{1, 2, 3, 4, 5, byte(i)})
		if err != nil {
			b.Error(err)
		}
		err = secVec[i].Set(msk, &idVec[i])
		if err != nil {
			b.Error(err)
		}
	}

	// recover from secVec and idVec
	var sec2 SecretKey
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		err := sec2.Recover(secVec, idVec)
		if err != nil {
			b.Errorf("%s\n", err)
		}
	}
}

func BenchmarkRecoverSeckey100(b *testing.B)  { benchmarkRecoverSeckey(100, b) }
func BenchmarkRecoverSeckey200(b *testing.B)  { benchmarkRecoverSeckey(200, b) }
func BenchmarkRecoverSeckey500(b *testing.B)  { benchmarkRecoverSeckey(500, b) }
func BenchmarkRecoverSeckey1000(b *testing.B) { benchmarkRecoverSeckey(1000, b) }

func benchmarkRecoverSignature(k int, b *testing.B) {
	b.StopTimer()
	err := Init(curve)
	if err != nil {
		b.Fatal(err)
	}
	var sec SecretKey
	sec.SetByCSPRNG()
	msk := sec.GetMasterSecretKey(k)

	// derive n shares
	n := k
	idVec := make([]ID, n)
	secVec := make([]SecretKey, n)
	signVec := make([]Sign, n)
	for i := 0; i < n; i++ {
		err := idVec[i].SetLittleEndian([]byte{1, 2, 3, 4, 5, byte(i)})
		if err != nil {
			b.Error(err)
		}
		err = secVec[i].Set(msk, &idVec[i])
		if err != nil {
			b.Error(err)
		}
		signVec[i] = *secVec[i].Sign("test message")
	}

	// recover signature
	var sig Sign
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		err := sig.Recover(signVec, idVec)
		if err != nil {
			b.Error(err)
		}
	}
}

func BenchmarkRecoverSignature100(b *testing.B)  { benchmarkRecoverSignature(100, b) }
func BenchmarkRecoverSignature200(b *testing.B)  { benchmarkRecoverSignature(200, b) }
func BenchmarkRecoverSignature500(b *testing.B)  { benchmarkRecoverSignature(500, b) }
func BenchmarkRecoverSignature1000(b *testing.B) { benchmarkRecoverSignature(1000, b) }
