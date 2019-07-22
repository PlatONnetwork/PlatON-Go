package state

import "time"

const (
	baseMs       = uint64(10 * time.Second)
	exponentBase = float64(1.5)
	maxExponent  = 2
)

type viewTimer struct {
	//Timer last timeout
	deadline time.Time
	timer    *time.Timer

	//Time window length calculation module
	timeInterval viewTimeInterval
}

func newViewTimer() *viewTimer {
	return &viewTimer{timer: time.NewTimer(0), timeInterval: viewTimeInterval{baseMs: baseMs, exponentBase: exponentBase, maxExponent: maxExponent}}
}

func (t *viewTimer) setupTimer(viewInterval uint64) {
	duration := t.timeInterval.getViewTimeInterval(viewInterval)
	t.deadline = time.Now().Add(duration)
	t.timer.Reset(duration)
}

func (t viewTimer) isDeadline() bool {
	return time.Now().Sub(t.deadline) <= 0
}

// Calculate the time window of each view，time=b*e^m
type viewTimeInterval struct {
	baseMs       uint64
	exponentBase float64
	maxExponent  uint64
}

func (vt viewTimeInterval) getViewTimeInterval(viewInterval uint64) time.Duration {
	return 0
}
