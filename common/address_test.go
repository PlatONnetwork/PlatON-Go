package common

import (
	"database/sql/driver"
	"encoding/json"
	"math/big"
	"reflect"
	"testing"
)

//func TestEmptyAddress(t *testing.T) {
//	add := MustBech32ToAddress("")
//	if add != ZeroAddr {
//		t.Error("ZeroAddr not compare")
//	}
//}

func TestIsStringAddress(t *testing.T) {
	tests := []struct {
		str string
		exp bool
	}{
		{"lat1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4", true},
		{"lat1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4", true},
		{"lao1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4", false},
		{"lam1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4", false},
		{"lat1x4w7852dxs70sy2mgf8w0s7tmvqx3cz2ydaxq4", false},
		{"lat1x4w7852dxs69sy2mgf8w0s7tmv", false},
		{"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beae", false},
		{"5aaeb6053f3e94c9b9a09f33669435e7ef1beaed11", false},
		{"0xxaaeb6053f3e94c9b9a09f33669435e7ef1beaed", false},
	}

	for _, test := range tests {
		if result := IsBech32Address(test.str); result != test.exp {
			t.Errorf("IsHexAddress(%s) == %v; expected %v",
				test.str, result, test.exp)
		}
	}
}

/*func TestIsHexAddress(t *testing.T) {
	tests := []struct {
		str string
		exp bool
	}{
		{"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", true},
		{"5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", true},
		{"0X5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", true},
		{"0XAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", true},
		{"0xAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA", true},
		{"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed1", false},
		{"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beae", false},
		{"5aaeb6053f3e94c9b9a09f33669435e7ef1beaed11", false},
		{"0xxaaeb6053f3e94c9b9a09f33669435e7ef1beaed", false},
	}

	for _, test := range tests {
		if result := IsHexAddress(test.str); result != test.exp {
			t.Errorf("IsHexAddress(%s) == %v; expected %v",
				test.str, result, test.exp)
		}
	}
}*/

func TestAddressUnmarshalJSON(t *testing.T) {
	byteInt, _ := new(big.Int).SetString("1447849688758881245826040560403612658163865659982", 10)
	var tests = []struct {
		Input     string
		ShouldErr bool
		Output    *big.Int
	}{
		{"", true, nil},
		{`""`, true, nil},
		{`"0x"`, true, nil},
		{`"0x00"`, true, nil},
		{`"0xG000000000000000000000000000000000000000"`, true, nil},
		{`"lac1flzyluu23zjknw70duwd00z6u9jgx7vug9n7t4"`, true, nil},
		{`"lax1lkdax58s3m3upsvmsk5wzcg55ydxp2jwqpvpf2"`, false, byteInt},
	}

	for i, test := range tests {
		var v Address
		err := json.Unmarshal([]byte(test.Input), &v)
		if err != nil && !test.ShouldErr {
			t.Errorf("test #%d: unexpected error: %v", i, err)
		}
		if err == nil {
			if test.ShouldErr {
				t.Errorf("test #%d: expected error, got none", i)
			}
			if v.Big().Cmp(test.Output) != 0 {
				t.Errorf("test #%d: address mismatch: have %v, want %v", i, v.Big(), test.Output)
			}
		}
	}
}

func TestAddressHexChecksum(t *testing.T) {
	var tests = []struct {
		Input  string
		Output Address
	}{
		//{"0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed", ZeroAddr},
		//{"0xfb6916095ca1df60bb79ce92ce3ea74c37c5d359", ZeroAddr},
		//{"0xdbf03b407c01e7cd3cbea99509d93f8dddc8c6fb", ZeroAddr},
		//{"0xd1220a0cf47c7b9be7a2e6ba89f429762e7b9adb", ZeroAddr},
		// Ensure that non-standard length input values are handled correctly
		//{"0xa", ZeroAddr},
		//{"0x0a", ZeroAddr},
		//{"0x00a", ZeroAddr},
		//{"0x000000000000000000000000000000000000000a", ZeroAddr},
		{"lat1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4", MustBech32ToAddress("lat1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4")},
	}
	for i, test := range tests {
		output := MustBech32ToAddress(test.Input)
		if output != test.Output {
			t.Errorf("test #%d: failed to match when it should (%s != %s)", i, output, test.Output)
		}
	}
}

//
//func BenchmarkAddressHex(b *testing.B) {
//	testAddr := HexToAddress("0x5aaeb6053f3e94c9b9a09f33669435e7ef1beaed")
//	for n := 0; n < b.N; n++ {
//		testAddr.Hex()
//	}
//}

func BenchmarkAddressString(b *testing.B) {
	testAddr := MustBech32ToAddress("lat1x4w7852dxs69sy2mgf8w0s7tmvqx3cz2ydaxq4")
	for n := 0; n < b.N; n++ {
		testAddr.String()
	}
}

func TestMixedcaseAccount_Address(t *testing.T) {

	// https://github.com/ethereum/EIPs/blob/master/EIPS/eip-55.md
	// Note: 0X{checksum_addr} is not valid according to spec above

	var res []struct {
		A     MixedcaseAddress
		Valid bool
	}
	if err := json.Unmarshal([]byte(`[
		{"A" : "0xae967917c465db8578ca9024c205720b1a3651A9", "Valid": false},
		{"A" : "0xAe967917c465db8578ca9024c205720b1a3651A9", "Valid": true},
		{"A" : "0XAe967917c465db8578ca9024c205720b1a3651A9", "Valid": false},
		{"A" : "0x1111111111111111111112222222222223333323", "Valid": true}
		]`), &res); err != nil {
		t.Fatal(err)
	}

	for _, r := range res {
		if got := r.A.ValidChecksum(); got != r.Valid {
			t.Errorf("Expected checksum %v, got checksum %v, input %v", r.Valid, got, r.A.String())
		}
	}

	//These should throw exceptions:
	var r2 []MixedcaseAddress
	for _, r := range []string{
		`["0x11111111111111111111122222222222233333"]`,     // Too short
		`["0x111111111111111111111222222222222333332"]`,    // Too short
		`["0x11111111111111111111122222222222233333234"]`,  // Too long
		`["0x111111111111111111111222222222222333332344"]`, // Too long
		`["1111111111111111111112222222222223333323"]`,     // Missing 0x
		`["x1111111111111111111112222222222223333323"]`,    // Missing 0
		`["0xG111111111111111111112222222222223333323"]`,   //Non-hex
	} {
		if err := json.Unmarshal([]byte(r), &r2); err == nil {
			t.Errorf("Expected failure, input %v", r)
		}

	}

}

func TestAddress_Scan(t *testing.T) {
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "working scan",
			args: args{src: []byte{
				0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
				0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
			}},
			wantErr: false,
		},
		{
			name:    "non working scan",
			args:    args{src: int64(1234567890)},
			wantErr: true,
		},
		{
			name: "invalid length scan",
			args: args{src: []byte{
				0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
				0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a,
			}},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &Address{}
			if err := a.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Address.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}

			if !tt.wantErr {
				for i := range a {
					if a[i] != tt.args.src.([]byte)[i] {
						t.Errorf(
							"Address.Scan() didn't scan the %d src correctly (have %X, want %X)",
							i, a[i], tt.args.src.([]byte)[i],
						)
					}
				}
			}
		})
	}
}

func TestAddress_Value(t *testing.T) {
	b := []byte{
		0xb2, 0x6f, 0x2b, 0x34, 0x2a, 0xab, 0x24, 0xbc, 0xf6, 0x3e,
		0xa2, 0x18, 0xc6, 0xa9, 0x27, 0x4d, 0x30, 0xab, 0x9a, 0x15,
	}
	var usedA Address
	usedA.SetBytes(b)
	tests := []struct {
		name    string
		a       Address
		want    driver.Value
		wantErr bool
	}{
		{
			name:    "Working value",
			a:       usedA,
			want:    b,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.a.Value()
			if (err != nil) != tt.wantErr {
				t.Errorf("Address.Value() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Address.Value() = %v, want %v", got, tt.want)
			}
		})
	}
}
