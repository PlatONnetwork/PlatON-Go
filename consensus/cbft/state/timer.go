package state

import (
	"math"
	"time"
)

const (
	exponentBase = float64(1.5)
	maxExponent  = 2
)

type viewTimer struct {
	//Timer last timeout
	deadline time.Time
	timer    *time.Timer

	//Time window length calculation module
	timeInterval    viewTimeInterval
	preViewInterval uint64
}

func newViewTimer(period uint64) *viewTimer {
	timer := time.NewTimer(0)
	timer.Stop()
	return &viewTimer{timer: timer,
		timeInterval:    viewTimeInterval{baseMs: period * uint64(time.Millisecond), exponentBase: exponentBase, maxExponent: maxExponent},
		preViewInterval: 1,
	}
}

// Ensure that the timeout period is adjusted smoothly.
// Each time the adjustment is compared with the previous one, it is gradually lower than the previous one, and then gradually decreases.
func (t *viewTimer) calViewInterval(viewInterval uint64) uint64 {
	if t.preViewInterval > 1 && viewInterval == 1 || t.preViewInterval > viewInterval {
		viewInterval = t.preViewInterval - 1
	} else if t.preViewInterval != 1 && t.preViewInterval == viewInterval {
		viewInterval += 1
	}
	t.preViewInterval = viewInterval
	return viewInterval
}

func (t *viewTimer) setupTimer(viewInterval uint64) {
	viewInterval = t.calViewInterval(viewInterval)
	duration := t.timeInterval.getViewTimeInterval(viewInterval)
	t.deadline = time.Now().Add(duration)
	t.stopTimer()
	t.timer.Reset(duration)
}

func (t *viewTimer) stopTimer() {
	if !t.timer.Stop() {
		select {
		case <-t.timer.C:
		default:
		}
	}
}
func (t *viewTimer) timerChan() <-chan time.Time {
	return t.timer.C
}

func (t viewTimer) isDeadline() bool {
	return t.deadline.Before(time.Now())
}

// Calculate the time window of each view，time=b*e^m
type viewTimeInterval struct {
	baseMs       uint64
	exponentBase float64
	maxExponent  uint64
}

func (vt viewTimeInterval) getViewTimeInterval(viewInterval uint64) time.Duration {
	pow := viewInterval - 1
	if pow > vt.maxExponent {
		pow = vt.maxExponent
	}
	mul := math.Pow(vt.exponentBase, float64(pow))
	return time.Duration(uint64(float64(vt.baseMs) * mul))
}
