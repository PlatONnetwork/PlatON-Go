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

package math

import "testing"

func TestBinomialDistribution(t *testing.T) {
	var n int64 = 10
	p := 0.6
	bd := NewBinomialDistribution(n, p)
	tp := 0.0
	for i := 0; i < 11; i++ {
		x, err := bd.InverseCumulativeProbability(tp)
		if nil != err {
			t.Error(err)
		}
		t.Log("tp", tp, "x", x)
		tp += 0.1
	}
	bd = NewBinomialDistribution(1000000, p)
	for i := 0; i < 1000000; i++ {
		tp, err := bd.CumulativeProbability(int64(i))
		if nil != err {
			t.Error(err)
		}
		if tp > 0 && tp < 1 {
			t.Log("x", i, "tp", tp)
		}
	}
	bd = NewBinomialDistribution((1<<63)-1, p)
	tp, err := bd.CumulativeProbability(int64(10000000))
	if nil != err {
		t.Error(err)
	}
	tp, err = bd.CumulativeProbability(int64(100000000))
	if nil != err {
		t.Error(err)
	}
	tp, err = bd.CumulativeProbability(int64(1000000000))
	if nil != err {
		t.Error(err)
	}
	tp, err = bd.CumulativeProbability(int64(10000000000))
	if nil != err {
		t.Error(err)
	}
	tp, err = bd.CumulativeProbability(int64(100000000000))
	if nil != err {
		t.Error(err)
	}
	tp, err = bd.CumulativeProbability(int64(1000000000000))
	if nil != err {
		t.Error(err)
	}
	tp = 0
	for i := 0; i < 11; i++ {
		x, err := bd.InverseCumulativeProbability(tp)
		if nil != err {
			t.Error(err)
		}
		t.Log("tp", tp, "x", x)
		tp += 0.1
	}
}
