package bls

import "testing"
import (
	"fmt"
)

func testBadPointOfG2(t *testing.T) {
/*	err := Init(CurveFp382_2)
	if err != nil {
		t.Fatal(err)
	}*/
	var Q G2
	// this value is not in G2 so should return an error
	err := Q.SetString("1 18d3d8c085a5a5e7553c3a4eb628e88b8465bf4de2612e35a0a4eb018fb0c82e9698896031e62fd7633ffd824a859474 1dc6edfcf33e29575d4791faed8e7203832217423bf7f7fbf1f6b36625b12e7132c15fbc15562ce93362a322fb83dd0d 65836963b1f7b6959030ddfa15ab38ce056097e91dedffd996c1808624fa7e2644a77be606290aa555cda8481cfb3cb 1b77b708d3d4f65aeedf54b58393463a42f0dc5856baadb5ce608036baeca398c5d9e6b169473a8838098fd72fd28b50", 16)
	if err == nil {
		t.Error(err)
	}
}

func testGT(t *testing.T) {
/*	err := Init(CurveFp382_2)
	if err != nil {
		t.Fatal(err)
	}*/
	var x GT
	x.Clear()
	if !x.IsZero() {
		t.Errorf("not zero")
	}
	x.SetInt64(1)
	if !x.IsOne() {
		t.Errorf("not one")
	}
}

func testNegAdd(t *testing.T) {
	var x Fr
	var P1, P2, P3 G1
	var Q1, Q2, Q3 G2
	err := P1.HashAndMapTo([]byte("this"))
	if err != nil {
		t.Error(err)
	}
	err = Q1.HashAndMapTo([]byte("this"))
	if err != nil {
		t.Error(err)
	}
	fmt.Printf("P1=%s\n", P1.GetString(16))
	fmt.Printf("Q1=%s\n", Q1.GetString(16))
	G1Neg(&P2, &P1)
	G2Neg(&Q2, &Q1)
	fmt.Printf("P2=%s\n", P2.GetString(16))
	fmt.Printf("Q2=%s\n", Q2.GetString(16))

	x.SetInt64(-1)
	G1Mul(&P3, &P1, &x)
	G2Mul(&Q3, &Q1, &x)
	if !P2.IsEqual(&P3) {
		t.Errorf("P2 != P3 %s\n", P3.GetString(16))
	}
	if !Q2.IsEqual(&Q3) {
		t.Errorf("Q2 != Q3 %s\n", Q3.GetString(16))
	}

	G1Add(&P2, &P2, &P1)
	G2Add(&Q2, &Q2, &Q1)
	if !P2.IsZero() {
		t.Errorf("P2 is not zero %s\n", P2.GetString(16))
	}
	if !Q2.IsZero() {
		t.Errorf("Q2 is not zero %s\n", Q2.GetString(16))
	}
}

func testPairing(t *testing.T) {
	var a, b, ab Fr
	err := a.SetString("123", 10)
	if err != nil {
		t.Error(err)
		return
	}
	err = b.SetString("456", 10)
	if err != nil {
		t.Error(err)
		return
	}
	FrMul(&ab, &a, &b)
	var P, aP G1
	var Q, bQ G2
	err = P.HashAndMapTo([]byte("this"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("P=%s\n", P.GetString(16))
	G1Mul(&aP, &P, &a)
	fmt.Printf("aP=%s\n", aP.GetString(16))
	err = Q.HashAndMapTo([]byte("that"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("Q=%s\n", Q.GetString(16))
	G2Mul(&bQ, &Q, &b)
	fmt.Printf("bQ=%s\n", bQ.GetString(16))
	var e1, e2 GT
	Pairing(&e1, &P, &Q)
	fmt.Printf("e1=%s\n", e1.GetString(16))
	Pairing(&e2, &aP, &bQ)
	fmt.Printf("e2=%s\n", e1.GetString(16))
	GTPow(&e1, &e1, &ab)
	fmt.Printf("e1=%s\n", e1.GetString(16))
	if !e1.IsEqual(&e2) {
		t.Errorf("not equal pairing\n%s\n%s", e1.GetString(16), e2.GetString(16))
	}

}

func testMclFor(t *testing.T) {
/*	err := Init(BLS12_381)
	if err != nil {
		t.Fatal(err)
	}*/
	//add for test
	var x1, x2  Fr
	err := x1.SetString("987", 10)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("x1=%s\n", x1.GetString(16))
	err = x2.SetString("654", 10)
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("x2=%s\n", x2.GetString(16))

	var m1, m2,sig1,sig2,sig G1
	var g12,h1,h2  G2
	err = m1.HashAndMapTo([]byte("this"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("m1=%s\n", m1.GetString(16))
	err = m2.HashAndMapTo([]byte("abcd"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("m2=%s\n", m2.GetString(16))

	G1Mul(&sig1, &m1, &x1)
	fmt.Printf("sig1=%s\n", sig1.GetString(16))
	G1Mul(&sig2, &m2, &x2)
	fmt.Printf("sig2=%s\n", sig2.GetString(16))
	G1Add(&sig, &sig1, &sig2)
	fmt.Printf("sig=%s\n", sig.GetString(16))

	err = g12.HashAndMapTo([]byte("defg"))
	if err != nil {
		t.Error(err)
		return
	}
	fmt.Printf("G=%s\n", g12.GetString(16))

	G2Mul(&h1, &g12, &x1)
	fmt.Printf("h1=%s\n", h1.GetString(16))
	G2Mul(&h2, &g12, &x2)
	fmt.Printf("h2=%s\n", h2.GetString(16))

	var e,e1,e2,e3,e11,e22,e33 GT
	Pairing(&e1, &m1, &h1)
	fmt.Printf("e1=%s\n", e1.GetString(16))
	Pairing(&e11, &sig1, &g12)
	fmt.Printf("e11=%s\n", e11.GetString(16))

	Pairing(&e2, &m2, &h2)
	fmt.Printf("e2=%s\n", e2.GetString(16))
	Pairing(&e22, &sig2, &g12)
	fmt.Printf("e22=%s\n", e22.GetString(16))

	Pairing(&e, &sig, &g12)
	fmt.Printf("e=%s\n", e.GetString(16))

	GTAdd(&e3, &e1, &e2)
	fmt.Printf("e3=%s\n", e3.GetString(16))
	GTAdd(&e33, &e11, &e22)
	fmt.Printf("e33=%s\n", e33.GetString(16))
	var e5 GT
	GTMul(&e5, &e1, &e2)
	fmt.Printf("e5=%s\n", e5.GetString(16))

	if !e.IsEqual(&e5) {
		t.Errorf("not equal pairing\n%s\n%s", e.GetString(16), e5.GetString(16))
	}

}


func testMcl(t *testing.T, c int) {
	err := Init(c)
	if err != nil {
		t.Fatal(err)
	}
	testNegAdd(t)
	testPairing(t)
	//add
	testMclFor(t)
	testGT(t)
	testBadPointOfG2(t)
}

func TestMclMain(t *testing.T) {
	t.Logf("GetMaxOpUnitSize() = %d\n", GetMaxOpUnitSize())
	t.Logf("GetFrUnitSize() = %d\n", GetFrUnitSize())
	t.Log("CurveFp254BNb")
	testMcl(t, CurveFp254BNb)
	if GetMaxOpUnitSize() == 6 {
		if GetFrUnitSize() == 6 {
			t.Log("CurveFp382_1")
			testMcl(t, CurveFp382_1)
			t.Log("CurveFp382_2")
			testMcl(t, CurveFp382_2)
		}
		t.Log("BLS12_381")
		testMcl(t, BLS12_381)
	}
}

func TestMclInit(t *testing.T) {
	fmt.Println("mcl MCLBN_FP_UNIT_SIZE:",GetMaxOpUnitSize())
	fmt.Println("mcl MCLBN_FR_UNIT_SIZE:",GetFrUnitSize())
	err := Init(CurveFp254BNb)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println("exc mcl CurveFp254BNb")
	if GetMaxOpUnitSize() == 6 {
		if GetFrUnitSize() == 6 {
			err := Init(CurveFp382_1)
			if err != nil {
				t.Fatal(err)
			}
			fmt.Println("exc mcl CurveFp382_1")
		}
		err := Init(BLS12_381)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println("exc mcl BLS12_381")
	}
}

