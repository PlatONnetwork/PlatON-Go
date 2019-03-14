package common

import "time"

type Timer struct {
	start time.Time
}

func NewTimer () *Timer {
	return new(Timer)
}

func (t *Timer) Begin() {
	t.start = time.Now()
}

func (t *Timer) End() float64 {
	tns := time.Since(t.start).Nanoseconds()
	tms := float64(tns) / float64(1e6)
	return tms

}
