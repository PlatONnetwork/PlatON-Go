// Copyright 2018-2020 The PlatON Network Authors
// This file is part of the PlatON-Go library.
//
// The PlatON-Go library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The PlatON-Go library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the PlatON-Go library. If not, see <http://www.gnu.org/licenses/>.

package utils

import (
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_RandomOffset(t *testing.T) {
	testCase := []struct {
		min int
		max int
	}{
		{min: 1, max: 10},
		{min: 2, max: 10},
		{min: 3, max: 10},
		{min: 4, max: 10},
		{min: 0, max: 5},
		{min: 5, max: 300},
	}
	for _, data := range testCase {
		offset := RandomOffset(data.min)
		t.Logf("offset: %d", offset)
		if data.min == 0 && offset != 0 {
			t.Fatalf("bad offset")
		}
		if offset < data.min && offset > data.max {
			t.Errorf("RandomOffset has incorrect value. offset:{%v}", offset)
		}
	}
}

func Test_BuildHash(t *testing.T) {
	var testCase = []struct {
		mType byte
		bytes []byte
		want  string
	}{
		{
			mType: 0x0,
			bytes: []byte("This is test data in 1"),
			want:  "0xcabbb3ea7b964fb678accab3051cd0893f0e94bca1d34304e9129c7c339bbcb4",
		},
		{
			mType: 0x1,
			bytes: []byte("This is test data in 2"),
			want:  "0xb4d9ca8710397e752c344724ea4d733473bb3f88afb94a7095c4dd2e4b61487a",
		},
		{
			mType: 0x2,
			bytes: []byte("This is test data in 3"),
			want:  "0xf449775f29162c3c63c740f93ab298a418145bac26ce120f5a16b55b0f7cb7d4",
		},
	}
	for _, v := range testCase {
		hash := BuildHash(v.mType, v.bytes)
		if hash.String() != v.want {
			t.Error("error")
		}
	}
}

// Detect performance of function execution.
func Benchmark_BuildHash(b *testing.B) {
	b.StopTimer()
	b.StartTimer()
	for i := 0; i < b.N; i++ {
		BuildHash(byte(i), []byte(fmt.Sprintf("%d", i)))
	}
}

func Test_MergeBytes(t *testing.T) {
	var testCase = []struct {
		flag bool
		el1  []byte
		el2  []byte
		want []byte
	}{
		{
			flag: true, el1: []byte("1"), el2: []byte("24"), want: []byte("124"),
		},
		{
			flag: true, el1: []byte("Hello"), el2: []byte(" World"), want: []byte("Hello World"),
		},
		{
			flag: false, el1: []byte("Hello "), el2: []byte(" World"), want: []byte("Hello World"),
		},
	}
	for _, v := range testCase {
		result := MergeBytes(v.el1, v.el2)
		// Simulate an exception using flag.
		if v.flag {
			assert.Equal(t, result, v.want)
		} else {
			assert.NotEqual(t, result, v.want)
		}
	}
}

func Test_SortMap(t *testing.T) {
	testCase := []struct {
		key   string
		value int64
	}{
		{"a", 1},
		{"c", 3},
		{"b", 2},
		{"d", 5},
		{"e", 2},
	}
	m := make(map[string]int64, len(testCase))
	for _, v := range testCase {
		m[v.key] = v.value
	}
	result := SortMap(m)
	t.Log(result)
	t.Log(result[:3])
	assert.Equal(t, "a", result[0].Key)
	assert.Equal(t, int64(1), result[0].Value)
}

func Test_Push(t *testing.T) {
	testCase := []struct {
		key   string
		value int64
	}{
		{"a", 1},
		{"c", 3},
		{"b", 2},
		{"d", 5},
		{"e", 2},
	}
	var pair KeyValuePairList
	for _, v := range testCase {
		pair.Push(KeyValuePair{v.key, v.value})
	}
	sort.Sort(pair)
	t.Log(pair)
	t.Log(pair[:3])
	assert.Equal(t, "a", pair[0].Key)
	assert.Equal(t, int64(1), pair[0].Value)
	assert.Equal(t, 5, pair.Len())
	oldPair := pair
	value := pair.Pop()
	v, ok := value.(KeyValuePair)
	assert.True(t, ok)
	assert.Equal(t, oldPair[oldPair.Len()-1].Key, v.Key)
	assert.Equal(t, oldPair.Len()-1, pair.Len())
}

func Test_Atomic(t *testing.T) {
	var at int32
	//go SetFalse(&at)
	SetFalse(&at)
	assert.True(t, False(&at))
	SetTrue(&at)
	assert.True(t, True(&at))
}
