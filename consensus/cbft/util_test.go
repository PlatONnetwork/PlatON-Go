package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/deckarep/golang-set"
	"math/big"
	"testing"
)

func TestRandomOffset(t *testing.T) {
	var expect = []struct{
		end 	 int
		wanted   bool
	}{
		{end: 10, wanted: true},
		{end: 0,  wanted: true},
	}
	for _, data := range expect {
		res := randomOffset(data.end)
		if !data.wanted && (0 > res || res > data.end) {
			t.Errorf("randomOffset has incorrect value. result:{%v}", res)
		}
	}
}

func TestRandomOffset_Collision(t *testing.T) {
	vals := make(map[int]struct{})
	for i := 0; i < 100; i++ {
		offset := randomOffset(2 << 30)
		if _, ok := vals[offset]; ok {
			t.Fatalf("got collision")
		}
		vals[offset] = struct{}{}
	}
}

func TestRandomOffset_Zero(t *testing.T) {
	offset := randomOffset(0)
	if offset != 0 {
		t.Fatalf("bad offset")
	}
}

func TestSetHash(t *testing.T) {
	var hashes mapset.Set
	hashes = mapset.NewSet()
	hashes.Add(common.BigToHash(big.NewInt(10)))
	hashes.Add(common.BigToHash(big.NewInt(11)))
	hashes.Add(common.BigToHash(big.NewInt(12)))

	var con interface{} = common.BigToHash(big.NewInt(10))
	if hashes.Contains(con) {
		t.Log("exists")
	} else {
		t.Error("not exists")
	}
}
