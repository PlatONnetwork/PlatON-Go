package cbft

import (
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
