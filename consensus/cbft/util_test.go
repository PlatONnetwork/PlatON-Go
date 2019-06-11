package cbft

import (
	"github.com/PlatONnetwork/PlatON-Go/common"
	"github.com/PlatONnetwork/PlatON-Go/crypto/sha3"
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

func TestProduceHash(t *testing.T) {
	var data = []struct{
		mType byte
		bytes []byte
		want string
	}{
		{
			mType: 0x0,
			bytes: []byte("This is test data in 1"),
			want: "0xcabbb3ea7b964fb678accab3051cd0893f0e94bca1d34304e9129c7c339bbcb4",
		},
		{
			mType: 0x1,
			bytes: []byte("This is test data in 2"),
			want: "0xb4d9ca8710397e752c344724ea4d733473bb3f88afb94a7095c4dd2e4b61487a",
		},
		{
			mType: 0x2,
			bytes: []byte("This is test data in 3"),
			want: "0xf449775f29162c3c63c740f93ab298a418145bac26ce120f5a16b55b0f7cb7d4",
		},
	}
	for _, v := range data {
		hash := produceHash(v.mType, v.bytes)
		if hash.String() != v.want {
			t.Error("error")
		}
	}

}

func TestUint64ToBytes(t *testing.T) {
	var wants = []struct{
		src uint64
		want string
	}{
		{
			src: 1558679713,
			want: "5fcb2251f5b31c73534c57718f0d60b23bc99898a0c4c4e69ae97b4a09f17205",
		},
		{
			src: 1558679714,
			want: "2b83fe25cd31f504192d7e9fa725f8d4d334d724feaf35cd26d225050b825683",
		},
		{
			src: 1558679715,
			want: "8ce02e5594c6da16f9c6d3958119c7ed6f0d25d3b45aeec10bcfdfe258aaf83f",
		},
	}
	for _, v := range wants {
		target := sha3.Sum256(uint64ToBytes(v.src))
		t_hex := common.Bytes2Hex(target[:])
		if t_hex != v.want {
			t.Error("error")
		}
	}
}