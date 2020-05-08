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
