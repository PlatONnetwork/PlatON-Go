package xcom

import "testing"

func TestBinomialDistribution(t *testing.T) {
	var n int64 = 10
	p := 0.6
	bd := NewBinomialDistribution(n, p)
	x, err := bd.InverseCumulativeProbability(0.46765656)
	if nil != err {
		t.Error(err)
	}
	t.Log("x", x)
	tp, err := bd.CumulativeProbability(3)
	if nil != err {
		t.Error(err)
	}
	t.Log("tp", tp)
}
